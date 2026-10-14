package regexpengine

// v8CompilerFirstCharacterSetBuilder ports the conservative FIRST(start) walk
// from regexp-compiler.cc. A true result means ranges contains a superset of
// every code unit that can begin a successful match. False means the set must
// not be used because some path could not be described statically.
type v8CompilerFirstCharacterSetBuilder struct {
	compiler   *v8CompilerCompiler
	ranges     []v8ASTCharacterRange
	budget     int
	codePoints bool
}

func v8CompilerComputeFirstCharacterSet(start v8NodesNode, compiler *v8CompilerCompiler) ([]v8ASTCharacterRange, bool) {
	if start == nil || compiler == nil {
		return nil, false
	}
	builder := v8CompilerFirstCharacterSetBuilder{compiler: compiler, budget: 500}
	if !builder.addNode(start, 0) {
		return nil, false
	}
	return v8TonodeCanonicalizeRanges(builder.ranges), true
}

// v8CompilerComputeFirstCodePointSet is used by Unicode-aware dispatch. The
// ordinary first-character set remains in UTF-16 code units because that is
// the matcher cursor representation; this variant combines a valid leading
// surrogate pair so alternatives beginning with different supplementary
// characters can be distinguished without changing the cursor model.
func v8CompilerComputeFirstCodePointSet(start v8NodesNode, compiler *v8CompilerCompiler) ([]v8ASTCharacterRange, bool) {
	if start == nil || compiler == nil {
		return nil, false
	}
	builder := v8CompilerFirstCharacterSetBuilder{compiler: compiler, budget: 500, codePoints: true}
	if !builder.addNode(start, 0) {
		return nil, false
	}
	return v8TonodeCanonicalizeRanges(builder.ranges), true
}

func (builder *v8CompilerFirstCharacterSetBuilder) addNode(node v8NodesNode, depth int) bool {
	if node == nil || depth > v8CompilerMaxRecursion {
		return false
	}
	builder.budget--
	if builder.budget < 0 {
		return false
	}
	switch value := node.(type) {
	case *v8NodesAssertionNode:
		return builder.addNode(value.OnSuccess(), depth+1)
	case *v8NodesActionNode:
		return builder.addAction(value, depth)
	case *v8NodesNegativeLookaroundChoiceNode:
		return builder.addNode(value.ContinueNode(), depth+1)
	case *v8NodesChoiceNode:
		return builder.addChoice(value, depth)
	case *v8NodesLoopChoiceNode:
		return builder.addChoice(value.AsChoiceNode(), depth)
	case *v8NodesTextNode:
		return builder.addText(value)
	case *v8NodesEndNode:
		// A backtrack node proves that this branch has no matching first code
		// unit. Accepting end nodes can match without consuming, so their first
		// code unit is not statically known.
		return value.IsBacktrack()
	case *v8NodesBackReferenceNode, *v8NodesUnanchoredAdvanceNode:
		return false
	default:
		// FIRST-set analysis is only a search optimization. A graph node added
		// by a later V8 compiler feature must disable the optimization instead
		// of preventing an otherwise valid expression from being emitted.
		return false
	}
}

func (builder *v8CompilerFirstCharacterSetBuilder) addAction(action *v8NodesActionNode, depth int) bool {
	switch action.ActionType() {
	case v8NodesStorePosition, v8NodesClearCaptures, v8NodesSetRegisterForLoop,
		v8NodesIncrementRegister, v8NodesEatsAtLeastAction,
		v8NodesPositiveSubmatchSuccess, v8NodesBeginNegativeSubmatch:
		return builder.addNode(action.OnSuccess(), depth+1)
	case v8NodesBeginPositiveSubmatch:
		// Prefer the positive-lookaround body when its first-character set is
		// known. If it reads backward or otherwise cannot be analyzed, the
		// continuation is still a safe constraint on the consumed input.
		if builder.addNode(action.OnSuccess(), depth+1) {
			return true
		}
		success := action.SuccessNode()
		return success != nil && builder.addNode(success.OnSuccess(), depth+1)
	case v8NodesEmptyMatchCheck, v8NodesRestorePosition:
		return false
	default:
		return false
	}
}

func (builder *v8CompilerFirstCharacterSetBuilder) addChoice(choice *v8NodesChoiceNode, depth int) bool {
	alternatives := choice.Alternatives()
	if alternatives == nil || alternatives.Len() == 0 {
		return false
	}
	for index := 0; index < alternatives.Len(); index++ {
		if !builder.addNode(alternatives.At(index).Node(), depth+1) {
			return false
		}
	}
	return true
}

func (builder *v8CompilerFirstCharacterSetBuilder) addText(text *v8NodesTextNode) bool {
	if text.finiteStringSet != nil {
		if text.finiteStringSet.hasEmptyString() {
			// This operand can accept without consuming input, so it cannot be
			// represented by a necessary first-character constraint.
			return false
		}
		ranges := text.finiteStringSet.firstCodePointRanges()
		if text.Flags().IgnoreCase {
			ranges = v8ASTCaseFoldingCloseOver(ranges, v8TonodeUnicodeCaseFolding)
		}
		if !builder.codePoints {
			ranges = v8GoFirstCharacterUnitRanges(ranges)
		}
		builder.ranges = append(builder.ranges, ranges...)
		return true
	}
	if text.ReadBackward() || text.Elements() == nil || text.Elements().Len() == 0 {
		return false
	}
	element := *text.Elements().At(0)
	if element.textType == v8ASTTextClassRanges {
		class, ok := element.tree.(*v8ASTClassRanges)
		if !ok || class == nil {
			return false
		}
		ranges := class.ranges()
		if !class.isNegated() {
			builder.ranges = append(builder.ranges, ranges...)
			return true
		}
		positive := v8TonodeCanonicalizeRanges(append([]v8ASTCharacterRange(nil), ranges...))
		builder.ranges = append(builder.ranges, v8TonodeNegateRanges(positive)...)
		return true
	}
	if element.textType != v8ASTTextAtom {
		return false
	}
	atom, ok := element.tree.(*v8ASTAtom)
	if !ok || atom == nil || len(atom.data) == 0 {
		return false
	}
	character := rune(atom.data[0])
	if builder.codePoints && IsEitherUnicode(text.Flags()) && character >= 0xd800 && character <= 0xdbff {
		if len(atom.data) < 2 || atom.data[1] < 0xdc00 || atom.data[1] > 0xdfff {
			// The graph may represent a supplementary literal as adjacent
			// one-unit atoms. Without seeing both units together, a code-point
			// dispatch set would be unsound; leave this choice on the normal path.
			return false
		}
		character = 0x10000 + (character-0xd800)<<10 + rune(atom.data[1]-0xdc00)
	}
	if !text.Flags().IgnoreCase {
		builder.ranges = append(builder.ranges, v8ASTRangeSingleton(character))
		return true
	}
	if character > 0xffff {
		return false
	}
	letters := v8CompilerCaseIndependentLetters(text.Flags(), uint16(character), builder.compiler.OneByte())
	if len(letters) == 0 {
		return false
	}
	for _, letter := range letters {
		builder.ranges = append(builder.ranges, v8ASTRangeSingleton(rune(letter)))
	}
	return true
}
