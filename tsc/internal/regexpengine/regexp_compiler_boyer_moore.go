package regexpengine

import "math/bits"

func v8CompilerSmearBitsRight(value uint32) uint32 {
	value |= value >> 1
	value |= value >> 2
	value |= value >> 4
	value |= value >> 8
	value |= value >> 16
	return value
}

// v8CompilerFillInBMInfo ports the node-local Boyer-Moore lookahead walk from
// regexp-compiler.cc. It only adds possible characters, so the maps remain a
// safe over-approximation even when the recursion budget is exhausted.
func v8CompilerFillInBMInfo(node v8NodesNode, offset, budget int, lookahead *v8NodesBoyerMooreLookahead, notAtStart bool) {
	if node == nil || lookahead == nil {
		panic("V8 FillInBMInfo requires a node and lookahead")
	}
	switch value := node.(type) {
	case *v8NodesNegativeSubmatchSuccess, *v8NodesEndNode, *v8NodesUnanchoredAdvanceNode:
		// End nodes add nothing; UnanchoredAdvance changes the input offset.
		return
	case *v8NodesBackReferenceNode:
		lookahead.SetRest(offset)
		value.SetBMInfo(notAtStart, lookahead)
		return
	case *v8NodesNegativeLookaroundChoiceNode:
		v8CompilerFillInBMInfo(value.ContinueNode(), offset, budget-1, lookahead, notAtStart)
		value.SetBMInfo(notAtStart, lookahead)
		return
	case *v8NodesLoopChoiceNode:
		if value.BodyCanBeZeroLength() || budget <= 0 {
			lookahead.SetRest(offset)
			value.SetBMInfo(notAtStart, lookahead)
			return
		}
		v8CompilerFillInBMChoice(value.AsChoiceNode(), offset, budget-1, lookahead, notAtStart)
		value.SetBMInfo(notAtStart, lookahead)
		return
	case *v8NodesChoiceNode:
		v8CompilerFillInBMChoice(value, offset, budget, lookahead, notAtStart)
		return
	case *v8NodesActionNode:
		if value.ActionType() == v8NodesPositiveSubmatchSuccess {
			value.SetBMInfo(notAtStart, lookahead)
			return
		}
		next := value.OnSuccess()
		if value.ActionType() == v8NodesBeginPositiveSubmatch {
			success := value.SuccessNode()
			if success == nil || success.OnSuccess() == nil {
				panic("V8 positive submatch has no BM successor")
			}
			next = success.OnSuccess()
		}
		v8CompilerFillInBMInfo(next, offset, budget-1, lookahead, notAtStart)
		value.SetBMInfo(notAtStart, lookahead)
		return
	case *v8NodesAssertionNode:
		if value.AssertionType() == v8NodesAssertAtStart && notAtStart {
			return
		}
		v8CompilerFillInBMInfo(value.OnSuccess(), offset, budget-1, lookahead, notAtStart)
		value.SetBMInfo(notAtStart, lookahead)
		return
	case *v8NodesTextNode:
		v8CompilerTextFillInBMInfo(value, offset, budget, lookahead, notAtStart)
		return
	default:
		panic("unknown V8 graph node in FillInBMInfo")
	}
}

func v8CompilerFillInBMChoice(choice *v8NodesChoiceNode, offset, budget int, lookahead *v8NodesBoyerMooreLookahead, notAtStart bool) {
	alternatives := choice.Alternatives()
	if alternatives == nil || alternatives.Len() == 0 {
		panic("V8 Boyer-Moore choice has no alternatives")
	}
	budget = (budget - 1) / alternatives.Len()
	for index := 0; index < alternatives.Len(); index++ {
		alternative := alternatives.At(index)
		if guards := alternative.Guards(); guards != nil && guards.Len() != 0 {
			lookahead.SetRest(offset)
			choice.SetBMInfo(notAtStart, lookahead)
			return
		}
		v8CompilerFillInBMInfo(alternative.Node(), offset, budget, lookahead, notAtStart)
	}
	choice.SetBMInfo(notAtStart, lookahead)
}

func v8CompilerTextFillInBMInfo(node *v8NodesTextNode, initialOffset, budget int, lookahead *v8NodesBoyerMooreLookahead, notAtStart bool) {
	if initialOffset >= lookahead.Length() || node.ReadBackward() || node.finiteStringSet != nil {
		return
	}
	offset := initialOffset
	for index := 0; index < node.Elements().Len(); index++ {
		if offset >= lookahead.Length() {
			if initialOffset == 0 {
				node.SetBMInfo(notAtStart, lookahead)
			}
			return
		}
		element := *node.Elements().At(index)
		if element.textType == v8ASTTextAtom {
			atom, ok := element.tree.(*v8ASTAtom)
			if !ok || atom == nil {
				panic("V8 TextNode BM analysis found an invalid atom")
			}
			for _, character := range atom.data {
				if offset >= lookahead.Length() {
					if initialOffset == 0 {
						node.SetBMInfo(notAtStart, lookahead)
					}
					return
				}
				if node.Flags().IgnoreCase {
					for _, equivalent := range v8CompilerCaseIndependentLetters(node.Flags(), character, lookahead.Compiler().OneByte()) {
						lookahead.Set(offset, int(equivalent))
					}
				} else {
					lookahead.Set(offset, int(character))
				}
				offset++
			}
			continue
		}
		if element.textType != v8ASTTextClassRanges {
			panic("unknown V8 text element in FillInBMInfo")
		}
		class, ok := element.tree.(*v8ASTClassRanges)
		if !ok || class == nil {
			panic("V8 TextNode BM analysis found invalid class ranges")
		}
		if class.isNegated() {
			lookahead.SetAll(offset)
		} else {
			for _, current := range class.ranges() {
				lookahead.SetInterval(offset, v8ASTNewInterval(int(current.from), int(current.to)))
			}
		}
		offset++
	}
	if offset >= lookahead.Length() {
		if initialOffset == 0 {
			node.SetBMInfo(notAtStart, lookahead)
		}
		return
	}
	v8CompilerFillInBMInfo(node.OnSuccess(), offset, budget-1, lookahead, true)
	if initialOffset == 0 {
		node.SetBMInfo(notAtStart, lookahead)
	}
}

// v8CompilerGetQuickCheckDetails ports the graph walk used by V8's quick check
// generator. The summaries are conservative: an admitted input is rechecked
// by the matcher, while impossible prefixes can be rejected immediately.
func v8CompilerGetQuickCheckDetails(node v8NodesNode, details *v8NodesQuickCheckDetails, compiler *v8CompilerCompiler, filledIn int, notAtStart bool, budget int) {
	if node == nil || details == nil || compiler == nil {
		panic("V8 GetQuickCheckDetails requires node, details, and compiler")
	}
	switch value := node.(type) {
	case *v8NodesNegativeSubmatchSuccess, *v8NodesEndNode:
		details.SetCannotMatchFrom(filledIn)
	case *v8NodesNegativeLookaroundChoiceNode:
		v8CompilerGetQuickCheckDetails(value.ContinueNode(), details, compiler, filledIn, notAtStart, budget-1)
	case *v8NodesLoopChoiceNode:
		if value.BodyCanBeZeroLength() || budget <= 0 {
			return
		}
		v8CompilerChoiceGetQuickCheckDetails(value.AsChoiceNode(), details, compiler, filledIn, notAtStart || value.NotAtStart(), budget)
	case *v8NodesChoiceNode:
		v8CompilerChoiceGetQuickCheckDetails(value, details, compiler, filledIn, notAtStart, budget)
	case *v8NodesActionNode:
		switch value.ActionType() {
		case v8NodesPositiveSubmatchSuccess:
			return
		case v8NodesBeginPositiveSubmatch:
			success := value.SuccessNode()
			if success == nil || success.OnSuccess() == nil {
				panic("V8 positive submatch has no quick-check successor")
			}
			v8CompilerGetQuickCheckDetails(success.OnSuccess(), details, compiler, filledIn, notAtStart, budget-1)
		default:
			v8CompilerGetQuickCheckDetails(value.OnSuccess(), details, compiler, filledIn, notAtStart, budget-1)
		}
	case *v8NodesAssertionNode:
		if value.AssertionType() == v8NodesAssertAtStart && notAtStart || value.AssertionType() == v8NodesAssertAtEnd {
			details.SetCannotMatchFrom(filledIn)
			return
		}
		v8CompilerGetQuickCheckDetails(value.OnSuccess(), details, compiler, filledIn, notAtStart, budget-1)
	case *v8NodesTextNode:
		v8CompilerTextGetQuickCheckDetails(value, details, compiler, filledIn, notAtStart, budget)
	case *v8NodesBackReferenceNode, *v8NodesUnanchoredAdvanceNode:
		// Their next input position cannot be summarized safely here.
		return
	default:
		panic("unknown V8 graph node in GetQuickCheckDetails")
	}
}

func v8CompilerChoiceGetQuickCheckDetails(choice *v8NodesChoiceNode, details *v8NodesQuickCheckDetails, compiler *v8CompilerCompiler, filledIn int, notAtStart bool, budget int) {
	alternatives := choice.Alternatives()
	if alternatives == nil || alternatives.Len() == 0 {
		panic("V8 quick-check choice has no alternatives")
	}
	notAtStart = notAtStart || choice.NotAtStart()
	budget /= alternatives.Len()
	v8CompilerGetQuickCheckDetails(alternatives.At(0).Node(), details, compiler, filledIn, notAtStart, budget)
	for index := 1; index < alternatives.Len(); index++ {
		other := v8NodesNewQuickCheckDetailsWithCharacters(details.Characters())
		v8CompilerGetQuickCheckDetails(alternatives.At(index).Node(), other, compiler, filledIn, notAtStart, budget)
		details.Merge(other, filledIn)
	}
}

func v8CompilerTextGetQuickCheckDetails(node *v8NodesTextNode, details *v8NodesQuickCheckDetails, compiler *v8CompilerCompiler, filledIn int, _ bool, budget int) {
	if node.ReadBackward() || node.finiteStringSet != nil {
		return
	}
	if filledIn >= details.Characters() {
		return
	}
	charMask := v8CompilerCharMask(compiler.OneByte())
	for index := 0; index < node.Elements().Len(); index++ {
		element := *node.Elements().At(index)
		if element.textType == v8ASTTextAtom {
			atom, ok := element.tree.(*v8ASTAtom)
			if !ok || atom == nil {
				panic("V8 TextNode quick-check analysis found an invalid atom")
			}
			for _, character := range atom.data {
				if filledIn == details.Characters() {
					return
				}
				position := details.Position(filledIn)
				if node.Flags().IgnoreCase {
					letters := v8CompilerCaseIndependentLetters(node.Flags(), character, compiler.OneByte())
					if len(letters) == 0 {
						details.SetCannotMatchFrom(filledIn)
						position.determinesPerfectly = false
						return
					}
					if len(letters) == 1 {
						position.mask = charMask
						position.value = uint32(letters[0])
						position.determinesPerfectly = true
					} else {
						commonBits, value := charMask, uint32(letters[0])
						for _, letter := range letters[1:] {
							differingBits := ((uint32(letter) & commonBits) ^ value)
							commonBits ^= differingBits
							value &= commonBits
						}
						oneZero := commonBits | ^charMask
						position.determinesPerfectly = len(letters) == 2 && (^oneZero)&((^oneZero)-1) == 0
						position.mask = commonBits
						position.value = value
					}
				} else {
					if uint32(character) > charMask {
						details.SetCannotMatchFrom(filledIn)
						position.determinesPerfectly = false
						return
					}
					position.mask = charMask
					position.value = uint32(character)
					position.determinesPerfectly = true
				}
				filledIn++
			}
			continue
		}
		if element.textType != v8ASTTextClassRanges {
			panic("unknown V8 text element in GetQuickCheckDetails")
		}
		if filledIn == details.Characters() {
			return
		}
		position := details.Position(filledIn)
		class, ok := element.tree.(*v8ASTClassRanges)
		if !ok || class == nil {
			panic("V8 TextNode quick-check analysis found invalid class ranges")
		}
		ranges := v8TonodeCanonicalizeRanges(class.ranges())
		if class.isNegated() || len(ranges) == 0 {
			position.mask, position.value = 0, 0
		} else {
			first := 0
			for first < len(ranges) && uint32(ranges[first].from) > charMask {
				first++
			}
			if first == len(ranges) {
				details.SetCannotMatchFrom(filledIn)
				position.determinesPerfectly = false
				return
			}
			firstFrom := uint32(ranges[first].from)
			firstTo := uint32(ranges[first].to)
			if firstTo > charMask {
				firstTo = charMask
			}
			totalCharacters := firstTo - firstFrom + 1
			differingBits := firstFrom ^ firstTo
			commonBits := ^v8CompilerSmearBitsRight(differingBits)
			value := firstFrom & commonBits
			for rangeIndex := first + 1; rangeIndex < len(ranges); rangeIndex++ {
				current := ranges[rangeIndex]
				from := uint32(current.from)
				if from > charMask {
					continue
				}
				to := uint32(current.to)
				if to > charMask {
					to = charMask
				}
				totalCharacters += to - from + 1
				newCommonBits := ^v8CompilerSmearBitsRight(from ^ to)
				commonBits &= newCommonBits
				value &= newCommonBits
				newDifferingBits := (from & commonBits) ^ value
				commonBits ^= newDifferingBits
				value &= commonBits
			}
			zeroBits := bits.OnesCount32((^commonBits) & charMask)
			position.mask, position.value = commonBits, value
			position.determinesPerfectly = totalCharacters == (1 << uint(zeroBits))
		}
		filledIn++
	}
	if details.CannotMatch() {
		return
	}
	v8CompilerGetQuickCheckDetails(node.OnSuccess(), details, compiler, filledIn, true, budget-1)
}
