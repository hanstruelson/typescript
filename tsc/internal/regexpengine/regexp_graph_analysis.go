package regexpengine

// v8CompilerAnalyzeRegExp ports AnalyzeRegExp's two successor-to-predecessor
// propagations. It uses an explicit DFS stack in place of C++ recursion and
// Isolate's StackLimitCheck; node visitation order and the per-node cycle
// markers remain the same.
func v8CompilerAnalyzeRegExp(root v8NodesNode, oneByte bool, stack v8CompilerAnalysisStackChecker) regexpError {
	if root == nil {
		panic("V8 AnalyzeRegExp requires a graph root")
	}
	type frame struct {
		node       v8NodesNode
		entered    bool
		childIndex int
		waiting    bool
	}
	frames := []frame{{node: root}}
	for len(frames) > 0 {
		current := &frames[len(frames)-1]
		if !current.entered {
			if stack != nil && stack.HasOverflowed() {
				for index := range frames {
					if frames[index].entered {
						info := frames[index].node.Info()
						info.beingAnalyzed = false
						info.beenAnalyzed = true
					}
				}
				return errorAnalysisStackOverflow
			}
			info := current.node.Info()
			if info.beenAnalyzed || info.beingAnalyzed {
				frames = frames[:len(frames)-1]
				continue
			}
			info.beingAnalyzed = true
			if text := current.node.AsTextNode(); text != nil {
				text.MakeCaseIndependent(oneByte)
			}
			current.entered = true
		}

		if current.waiting {
			if text := current.node.AsTextNode(); text != nil {
				text.CalculateOffsets()
			}
			v8CompilerPropagateNodeAnalysis(current.node, current.childIndex)
			current.childIndex++
			current.waiting = false
			continue
		}

		childCount := v8CompilerAnalysisChildCount(current.node)
		if current.childIndex == childCount {
			info := current.node.Info()
			info.beingAnalyzed = false
			info.beenAnalyzed = true
			frames = frames[:len(frames)-1]
			continue
		}

		child := v8CompilerAnalysisChildAt(current.node, current.childIndex)
		childInfo := child.Info()
		if childInfo.beenAnalyzed || childInfo.beingAnalyzed {
			if text := current.node.AsTextNode(); text != nil {
				text.CalculateOffsets()
			}
			v8CompilerPropagateNodeAnalysis(current.node, current.childIndex)
			current.childIndex++
			continue
		}
		current.waiting = true
		frames = append(frames, frame{node: child})
	}
	return errorNone
}

func v8CompilerAnalysisChildCount(node v8NodesNode) int {
	switch value := node.(type) {
	case *v8NodesNegativeSubmatchSuccess, *v8NodesEndNode:
		return 0
	case *v8NodesNegativeLookaroundChoiceNode:
		if value.Alternatives() == nil || value.Alternatives().Len() != 2 {
			panic("V8 negative lookaround must have two alternatives")
		}
		return 2
	case *v8NodesLoopChoiceNode:
		if value.Alternatives() == nil || value.Alternatives().Len() != 2 {
			panic("V8 loop choice must have loop and continue alternatives")
		}
		return 2
	case *v8NodesChoiceNode:
		if value.Alternatives() == nil || value.Alternatives().Len() == 0 {
			panic("V8 choice node must have at least one alternative")
		}
		return value.Alternatives().Len()
	default:
		if seq := node.AsSeqNode(); seq != nil {
			if seq.OnSuccess() == nil {
				panic("V8 sequence node has no successor")
			}
			return 1
		}
		panic("unknown V8 graph node in AnalyzeRegExp")
	}
}

func v8CompilerAnalysisChildAt(node v8NodesNode, index int) v8NodesNode {
	switch value := node.(type) {
	case *v8NodesNegativeSubmatchSuccess, *v8NodesEndNode:
		panic("V8 end node has no analysis children")
	case *v8NodesNegativeLookaroundChoiceNode:
		if index == v8NodesLookaroundIndex {
			return value.LookaroundNode()
		}
		if index == v8NodesContinueIndex {
			return value.ContinueNode()
		}
	case *v8NodesLoopChoiceNode:
		// V8 deliberately analyzes the continuation before the loop body.
		if index == 0 {
			return value.ContinueNode()
		}
		if index == 1 {
			return value.LoopNode()
		}
	case *v8NodesChoiceNode:
		if index >= 0 && index < value.Alternatives().Len() {
			return value.Alternatives().At(index).Node()
		}
	default:
		if index == 0 {
			if seq := node.AsSeqNode(); seq != nil {
				return seq.OnSuccess()
			}
		}
	}
	panic("V8 graph analysis child index is out of range")
}

func v8CompilerPropagateNodeAnalysis(node v8NodesNode, childIndex int) {
	// AssertionPropagator runs before EatsAtLeastPropagator for every visitor
	// edge, matching Analysis's variadic propagator order in V8.
	switch value := node.(type) {
	case *v8NodesNegativeSubmatchSuccess, *v8NodesEndNode:
		return
	case *v8NodesNegativeLookaroundChoiceNode:
		child := v8CompilerAnalysisChildAt(value, childIndex)
		value.Info().AddFromFollowing(*child.Info())
		if childIndex == v8NodesContinueIndex {
			value.SetEatsAtLeastInfo(*child.EatsAtLeastInfo())
		}
		return
	case *v8NodesLoopChoiceNode:
		child := v8CompilerAnalysisChildAt(value, childIndex)
		value.Info().AddFromFollowing(*child.Info())
		if childIndex == 0 && !value.ReadBackward() {
			value.SetEatsAtLeastInfo(*child.EatsAtLeastInfo())
		}
		return
	case *v8NodesChoiceNode:
		child := v8CompilerAnalysisChildAt(value, childIndex)
		value.Info().AddFromFollowing(*child.Info())
		eats := *value.EatsAtLeastInfo()
		if childIndex == 0 {
			eats = v8NodesNewEatsAtLeastInfo(0xff)
		}
		eats.SetMin(*child.EatsAtLeastInfo())
		value.SetEatsAtLeastInfo(eats)
		return
	}

	seq := node.AsSeqNode()
	if seq == nil {
		panic("V8 sequence analysis received a non-sequence node")
	}
	successor := seq.OnSuccess()
	switch value := node.(type) {
	case *v8NodesTextNode:
		if !value.ReadBackward() {
			eats := v8CompilerSaturateUint8(value.Length() + int(successor.EatsAtLeastInfo().fromNotStart))
			value.SetEatsAtLeastInfo(v8NodesNewEatsAtLeastInfo(eats))
		}
	case *v8NodesActionNode:
		switch value.ActionType() {
		case v8NodesBeginPositiveSubmatch:
			successNode := value.SuccessNode()
			if successNode == nil || successNode.OnSuccess() == nil {
				panic("V8 positive submatch has no success successor")
			}
			value.SetEatsAtLeastInfo(*successNode.OnSuccess().EatsAtLeastInfo())
		case v8NodesPositiveSubmatchSuccess:
			// A positive lookahead success rewinds the input position; the source
			// deliberately leaves its zero initialized lower bound unchanged.
		case v8NodesEatsAtLeastAction:
			eats := *successor.EatsAtLeastInfo()
			eats.SetMax(value.StoredEatsAtLeast())
			value.SetEatsAtLeastInfo(eats)
		default:
			value.SetEatsAtLeastInfo(*successor.EatsAtLeastInfo())
		}
	case *v8NodesUnanchoredAdvanceNode:
		eats := v8CompilerSaturateUint8(1 + int(successor.EatsAtLeastInfo().fromNotStart))
		value.SetEatsAtLeastInfo(v8NodesNewEatsAtLeastInfo(eats))
	case *v8NodesBackReferenceNode:
		if !value.ReadBackward() {
			value.SetEatsAtLeastInfo(*successor.EatsAtLeastInfo())
		}
	case *v8NodesAssertionNode:
		eats := *successor.EatsAtLeastInfo()
		if value.AssertionType() == v8NodesAssertAtStart {
			eats.fromNotStart = 0xff
		}
		value.SetEatsAtLeastInfo(eats)
	default:
		panic("unknown V8 sequence node in AnalyzeRegExp")
	}
}

func v8CompilerSaturateUint8(value int) uint8 {
	if value < 0 {
		return 0
	}
	if value > 0xff {
		return 0xff
	}
	return uint8(value)
}
