package regexpengine

// v8GoBoyerMooreSearchPlan contains the conservative character filter and
// shift selected by V8's Boyer-Moore lookahead analysis. The Go emitter uses
// it to move the initial search position past starts that cannot match.
type v8GoBoyerMooreSearchPlan struct {
	enabled        bool
	offset         int
	advanceBy      int
	characterMap   [2]uint64
	unicodeAdvance bool
}

func v8GoBuildBoyerMooreSearchPlan(body v8NodesNode, compiler *v8CompilerCompiler) v8GoBoyerMooreSearchPlan {
	if body == nil || compiler == nil || !compiler.HasSearchPrefix() || IsEitherUnicode(compiler.Flags()) {
		return v8GoBoyerMooreSearchPlan{}
	}
	// V8's search analysis is built around forward input offsets. A lookaround
	// can inspect input behind the candidate, so translating its quick-check
	// offset directly into a candidate-start offset is unsound. Keep the normal
	// graph search for these patterns until the plan carries signed offsets.
	if v8GoGraphContainsLookaround(body) {
		return v8GoBoyerMooreSearchPlan{}
	}

	length := int(body.EatsAtLeast(false))
	if length == 0 {
		return v8GoBoyerMooreSearchPlan{}
	}
	if length > v8CompilerMaxLookaheadForBoyerMoore {
		length = v8CompilerMaxLookaheadForBoyerMoore
	}
	lookahead := v8NodesNewBoyerMooreLookahead(length, compiler)
	lookahead.SetCachesNodeInfo(false)
	v8CompilerFillInBMInfo(body, 0, v8NodesRecursionBudget, lookahead, false)
	offset, advanceBy, table, _, found := lookahead.BuildSkipTable()
	if !found || advanceBy <= 0 {
		return v8GoBoyerMooreSearchPlan{}
	}
	minLookahead := offset + 1 - advanceBy
	for index := minLookahead; index <= offset; index++ {
		if lookahead.Count(index) == 0 {
			// V8 emits an unconditional backtrack for an empty position map.
			// The Go prefilter only represents a union table, so keep the normal
			// graph search when this special case occurs.
			return v8GoBoyerMooreSearchPlan{}
		}
	}

	var characterMap [2]uint64
	for character, possible := range table {
		if possible != 0 {
			characterMap[character>>6] |= uint64(1) << uint(character&63)
		}
	}
	if characterMap[0] == 0 && characterMap[1] == 0 {
		return v8GoBoyerMooreSearchPlan{}
	}
	return v8GoBoyerMooreSearchPlan{
		enabled:      true,
		offset:       offset,
		advanceBy:    advanceBy,
		characterMap: characterMap,
	}
}

func v8GoGraphContainsLookaround(root v8NodesNode) bool {
	if root == nil {
		return false
	}
	seen := make(map[v8NodesNode]bool)
	stack := []v8NodesNode{root}
	for len(stack) != 0 {
		index := len(stack) - 1
		node := stack[index]
		stack = stack[:index]
		if node == nil || seen[node] {
			continue
		}
		seen[node] = true
		switch value := node.(type) {
		case *v8NodesNegativeLookaroundChoiceNode:
			return true
		case *v8NodesActionNode:
			switch value.ActionType() {
			case v8NodesBeginPositiveSubmatch, v8NodesBeginNegativeSubmatch:
				return true
			}
			if successor := value.OnSuccess(); successor != nil {
				stack = append(stack, successor)
			}
		case *v8NodesChoiceNode:
			if alternatives := value.Alternatives(); alternatives != nil {
				for child := 0; child < alternatives.Len(); child++ {
					stack = append(stack, alternatives.At(child).Node())
				}
			}
		case *v8NodesLoopChoiceNode:
			if alternatives := value.Alternatives(); alternatives != nil {
				for child := 0; child < alternatives.Len(); child++ {
					stack = append(stack, alternatives.At(child).Node())
				}
			}
		default:
			if seq := node.AsSeqNode(); seq != nil {
				stack = append(stack, seq.OnSuccess())
			}
		}
	}
	return false
}
