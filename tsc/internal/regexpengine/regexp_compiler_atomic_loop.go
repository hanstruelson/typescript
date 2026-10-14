package regexpengine

// v8CompilerAtomicLoopBodyAnalysis records the facts collected while walking
// one fixed-length loop body alternative.
type v8CompilerAtomicLoopBodyAnalysis struct {
	fixedLengthEligible bool
	setKnown            bool
}

type v8CompilerAtomicLoopContinuationAnalysis struct {
	atEndAccept        bool
	alwaysSucceeds     bool
	startsWithBoundary bool
	firstSetKnown      bool
}

const v8CompilerContinuationAlwaysSucceedsBudget = 1000

func v8CompilerRangesSubsetOfSpecialClass(ranges []v8ASTCharacterRange, special []int) bool {
	if !v8TonodeRangesAreCanonical(ranges) || len(special) < 3 || special[len(special)-1] != v8TonodeRangeEndMarker {
		return false
	}
	entryCount := len(special) - 1
	entry := 0
	for _, current := range ranges {
		for entry < entryCount && rune(special[entry+1]) <= current.from {
			entry += 2
		}
		if entry >= entryCount || current.from < rune(special[entry]) || current.to >= rune(special[entry+1]) {
			return false
		}
	}
	return true
}

func v8CompilerAppendTextElementMatchSet(element v8ASTTextElement, flags Flags, out *[]v8ASTCharacterRange) bool {
	switch element.textType {
	case v8ASTTextClassRanges:
		class, ok := element.tree.(*v8ASTClassRanges)
		if !ok || class == nil {
			return false
		}
		ranges := class.ranges()
		if !class.isNegated() {
			*out = append(*out, ranges...)
			return true
		}
		positive := v8TonodeCanonicalizeRanges(append([]v8ASTCharacterRange(nil), ranges...))
		*out = append(*out, v8TonodeNegateRanges(positive)...)
		return true
	case v8ASTTextAtom:
		if flags.IgnoreCase {
			return false
		}
		atom, ok := element.tree.(*v8ASTAtom)
		if !ok || atom == nil {
			return false
		}
		for _, character := range atom.data {
			*out = append(*out, v8ASTRangeSingleton(rune(character)))
		}
		return true
	default:
		return false
	}
}

func v8CompilerAnalyzeAtomicLoopBody(alternative *v8NodesGuardedAlternative, loop *v8NodesLoopChoiceNode, bodySet *[]v8ASTCharacterRange) v8CompilerAtomicLoopBodyAnalysis {
	result := v8CompilerAtomicLoopBodyAnalysis{setKnown: true}
	if alternative == nil || loop == nil {
		return result
	}
	node := alternative.Node()
	for depth := 0; depth <= v8CompilerMaxRecursion; depth++ {
		if node == loop {
			result.fixedLengthEligible = true
			return result
		}
		if node.FixedLengthLoopLength() == v8NodesNodeIsTooComplexForFixedLengthLoops {
			return result
		}
		seq := node.AsSeqNode()
		if seq == nil {
			return result
		}
		text := node.AsTextNode()
		if text == nil || text.ReadBackward() || text.finiteStringSet != nil {
			result.setKnown = false
			if text != nil && text.finiteStringSet != nil {
				// A finite string set can consume different lengths, so the
				// fixed-width loop rewrite must leave it on the ordinary path.
				return result
			}
		} else {
			for index := 0; index < text.Elements().Len(); index++ {
				if !v8CompilerAppendTextElementMatchSet(*text.Elements().At(index), text.Flags(), bodySet) {
					result.setKnown = false
				}
			}
		}
		node = seq.OnSuccess()
	}
	return result
}

func v8CompilerContinuationAlwaysSucceeds(node v8NodesNode, depth int, budget *int) bool {
	if depth > v8CompilerMaxRecursion || node == nil {
		return false
	}
	*budget--
	if *budget < 0 {
		return false
	}
	if action := node.AsActionNode(); action != nil {
		switch action.ActionType() {
		case v8NodesStorePosition, v8NodesClearCaptures, v8NodesSetRegisterForLoop,
			v8NodesIncrementRegister, v8NodesEatsAtLeastAction,
			v8NodesRestorePosition, v8NodesPositiveSubmatchSuccess,
			v8NodesBeginPositiveSubmatch:
			return v8CompilerContinuationAlwaysSucceeds(action.OnSuccess(), depth+1, budget)
		case v8NodesEmptyMatchCheck, v8NodesBeginNegativeSubmatch:
			return false
		}
	}
	if end := node.AsEndNode(); end != nil {
		return end.Action() == v8NodesAccept
	}
	if choice := node.AsChoiceNode(); choice != nil {
		if node.AsNegativeLookaroundChoiceNode() != nil {
			return false
		}
		for index := 0; index < choice.Alternatives().Len(); index++ {
			alternative := choice.Alternatives().At(index)
			if guards := alternative.Guards(); guards != nil && guards.Len() != 0 {
				continue
			}
			if v8CompilerContinuationAlwaysSucceeds(alternative.Node(), depth+1, budget) {
				return true
			}
		}
	}
	return false
}

func v8CompilerAnalyzeAtomicLoopContinuation(alternative *v8NodesGuardedAlternative, firstSet *[]v8ASTCharacterRange) v8CompilerAtomicLoopContinuationAnalysis {
	result := v8CompilerAtomicLoopContinuationAnalysis{}
	if alternative == nil {
		return result
	}
	if guards := alternative.Guards(); guards != nil && guards.Len() != 0 {
		return result
	}
	node := alternative.Node()
	budget := v8CompilerContinuationAlwaysSucceedsBudget
	result.alwaysSucceeds = v8CompilerContinuationAlwaysSucceeds(node, 0, &budget)
	if assertion := node.AsAssertionNode(); assertion != nil {
		result.startsWithBoundary = assertion.AssertionType() == v8NodesAssertAtBoundary
	}
	atEndAlive := false
	if assertion := node.AsAssertionNode(); assertion != nil {
		atEndAlive = assertion.AssertionType() == v8NodesAssertAtEnd
	}
	firstAlive := true
	for depth := 0; depth <= v8CompilerMaxRecursion; depth++ {
		if assertion := node.AsAssertionNode(); assertion != nil {
			if depth > 0 {
				atEndAlive = false
			}
			node = assertion.OnSuccess()
			continue
		}
		if action := node.AsActionNode(); action != nil {
			switch action.ActionType() {
			case v8NodesStorePosition, v8NodesClearCaptures:
			case v8NodesSetRegisterForLoop, v8NodesIncrementRegister:
				firstAlive = false
			case v8NodesRestorePosition, v8NodesBeginPositiveSubmatch,
				v8NodesBeginNegativeSubmatch, v8NodesPositiveSubmatchSuccess,
				v8NodesEmptyMatchCheck, v8NodesEatsAtLeastAction:
				return result
			default:
				return result
			}
			node = action.OnSuccess()
			continue
		}
		if end := node.AsEndNode(); end != nil {
			result.atEndAccept = atEndAlive && end.Action() == v8NodesAccept
			return result
		}
		if firstAlive {
			text := node.AsTextNode()
			if text != nil && !text.ReadBackward() && text.Elements() != nil && text.Elements().Len() > 0 {
				element := *text.Elements().At(0)
				switch element.textType {
				case v8ASTTextClassRanges:
					class, ok := element.tree.(*v8ASTClassRanges)
					if ok && class != nil {
						if firstSet != nil {
							if class.isNegated() {
								positive := v8TonodeCanonicalizeRanges(append([]v8ASTCharacterRange(nil), class.ranges()...))
								*firstSet = append(*firstSet, v8TonodeNegateRanges(positive)...)
							} else {
								*firstSet = append(*firstSet, class.ranges()...)
							}
						}
						result.firstSetKnown = true
					}
				case v8ASTTextAtom:
					atom, ok := element.tree.(*v8ASTAtom)
					if ok && atom != nil && len(atom.data) > 0 && !text.Flags().IgnoreCase {
						if firstSet != nil {
							*firstSet = append(*firstSet, v8ASTRangeSingleton(rune(atom.data[0])))
						}
						result.firstSetKnown = true
					}
				}
			}
		}
		return result
	}
	return result
}

func v8CompilerClassifyAtomicLoop(loop *v8NodesLoopChoiceNode) v8NodesAtomicLoopKind {
	if loop == nil || loop.Alternatives() == nil || loop.Alternatives().Len() != 2 {
		return v8NodesAtomicLoopNone
	}
	body := loop.Alternatives().At(0)
	continuation := loop.Alternatives().At(1)
	// A source-level character class is one deterministic code-point test, but
	// the graph may express Unicode/surrogate handling as several internal
	// alternatives. Once an end assertion is reached, none of those completed
	// iterations needs to be revisited: a shorter repetition would leave an
	// unmatched code point before the end (or final line terminator).
	if loop.LoopNode() == body.Node() && loop.ContinueNode() == continuation.Node() &&
		!loop.BodyCanBeZeroLength() && loop.SingleCodePointClassBody() {
		continuationInfo := v8CompilerAnalyzeAtomicLoopContinuation(continuation, nil)
		if continuationInfo.atEndAccept {
			return v8NodesAtomicLoopAtEnd
		}
	}
	bodySet := make([]v8ASTCharacterRange, 0, 4)
	bodyInfo := v8CompilerAnalyzeAtomicLoopBody(body, loop, &bodySet)
	if !bodyInfo.fixedLengthEligible {
		return v8NodesAtomicLoopNone
	}
	firstSet := make([]v8ASTCharacterRange, 0, 2)
	continuationInfo := v8CompilerAnalyzeAtomicLoopContinuation(continuation, &firstSet)
	if continuationInfo.atEndAccept {
		return v8NodesAtomicLoopAtEnd
	}
	if continuationInfo.alwaysSucceeds {
		return v8NodesAtomicLoopTotal
	}
	if !bodyInfo.setKnown {
		return v8NodesAtomicLoopNone
	}
	bodySet = v8TonodeCanonicalizeRanges(bodySet)
	if continuationInfo.firstSetKnown {
		firstSet = v8TonodeCanonicalizeRanges(firstSet)
		if !v8TonodeRangesIntersect(bodySet, firstSet) {
			return v8NodesAtomicLoopDisjoint
		}
	}
	if continuationInfo.startsWithBoundary && v8CompilerRangesSubsetOfSpecialClass(bodySet, v8TonodeWordRanges) {
		return v8NodesAtomicLoopBoundary
	}
	return v8NodesAtomicLoopNone
}

func v8NodesLoopChoiceAtomicLoopKind(loop *v8NodesLoopChoiceNode) v8NodesAtomicLoopKind {
	if loop == nil {
		return v8NodesAtomicLoopNone
	}
	if kind, valid := loop.AtomicLoopKindCache(); valid {
		return kind
	}
	kind := v8CompilerClassifyAtomicLoop(loop)
	loop.SetAtomicLoopKindCache(kind)
	return kind
}
