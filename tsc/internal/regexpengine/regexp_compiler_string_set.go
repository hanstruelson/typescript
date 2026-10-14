package regexpengine

import (
	"encoding/binary"
	"sort"
)

// v8NodesFiniteStringSet is a compact acyclic automaton for a Unicode-set
// operand made from a finite collection of strings. It is not a regex VM: the
// automaton recognizes only the literal strings in one class-set atom, while
// the surrounding expression remains the normal V8-shaped graph.
type v8NodesFiniteStringSet struct {
	root   int
	nodes  []v8NodesFiniteStringSetNode
	edges  []v8NodesFiniteStringSetEdge
	ranges []v8ASTCharacterRange
}

type v8NodesFiniteStringSetNode struct {
	firstEdge int
	edgeCount int
	terminal  bool
}

type v8NodesFiniteStringSetEdge struct {
	codePoint rune
	target    int
}

type v8TonodeMutableStringSetNode struct {
	terminal bool
	children map[rune]*v8TonodeMutableStringSetNode
}

type v8TonodeCanonicalStringSetNode struct {
	terminal bool
	edges    []v8NodesFiniteStringSetEdge
}

// v8TonodeFiniteStringSetSourceCost estimates the amount of direct matcher
// source that would be emitted for a finite string set. Large sets use a
// minimized trie table so Go's compiler sees a small matcher body instead of
// thousands of repeated graph labels and backtracking blocks.
func v8TonodeFiniteStringSetSourceCost(stringsIn v8ASTCharacterClassStrings) int {
	cost := 0
	for _, item := range stringsIn {
		// One literal code point needs a bounds check, a comparison, position
		// movement, and usually a backtracking frame at its surrounding fork.
		cost += 160 + 48*len(item.codePoints)
	}
	return cost
}

func v8TonodeShouldUseFiniteStringSet(stringsIn v8ASTCharacterClassStrings) bool {
	// Keep small sets as ordinary graph nodes so the Go compiler can inline
	// their direct comparisons. Switch only when the rough direct-source cost
	// is at least 16 KiB; this is a representation choice, not a regex-size
	// limit, and does not reject any valid pattern.
	return len(stringsIn) > 0 && v8TonodeFiniteStringSetSourceCost(stringsIn) >= 16*1024
}

func v8TonodeBuildFiniteStringSet(stringsIn v8ASTCharacterClassStrings, ranges []v8ASTCharacterRange) *v8NodesFiniteStringSet {
	root := &v8TonodeMutableStringSetNode{}
	for _, item := range stringsIn {
		current := root
		for _, codePoint := range item.codePoints {
			if current.children == nil {
				current.children = make(map[rune]*v8TonodeMutableStringSetNode)
			}
			next := current.children[codePoint]
			if next == nil {
				next = &v8TonodeMutableStringSetNode{}
				current.children[codePoint] = next
			}
			current = next
		}
		current.terminal = true
	}

	// Intern equivalent suffix states. This turns the prefix trie into a
	// minimal acyclic deterministic automaton for the finite language, sharing
	// both common prefixes and identical suffix subtrees without recursion.
	type frame struct {
		node   *v8TonodeMutableStringSetNode
		expand bool
	}
	stack := []frame{{node: root}}
	canonicalID := make(map[*v8TonodeMutableStringSetNode]int)
	canonical := make([]v8TonodeCanonicalStringSetNode, 0, len(stringsIn))
	intern := make(map[string]int)
	for len(stack) != 0 {
		last := len(stack) - 1
		current := stack[last]
		stack = stack[:last]
		if !current.expand {
			stack = append(stack, frame{node: current.node, expand: true})
			keys := v8TonodeSortedStringSetKeys(current.node.children)
			for index := len(keys) - 1; index >= 0; index-- {
				stack = append(stack, frame{node: current.node.children[keys[index]]})
			}
			continue
		}

		keys := v8TonodeSortedStringSetKeys(current.node.children)
		edges := make([]v8NodesFiniteStringSetEdge, 0, len(keys))
		for _, key := range keys {
			edges = append(edges, v8NodesFiniteStringSetEdge{codePoint: key, target: canonicalID[current.node.children[key]]})
		}
		signature := make([]byte, 1+8*len(edges))
		if current.node.terminal {
			signature[0] = 1
		}
		for index, edge := range edges {
			binary.LittleEndian.PutUint32(signature[1+8*index:], uint32(edge.codePoint))
			binary.LittleEndian.PutUint32(signature[1+8*index+4:], uint32(edge.target))
		}
		key := string(signature)
		if existing, ok := intern[key]; ok {
			canonicalID[current.node] = existing
			continue
		}
		id := len(canonical)
		canonicalID[current.node] = id
		intern[key] = id
		canonical = append(canonical, v8TonodeCanonicalStringSetNode{terminal: current.node.terminal, edges: edges})
	}

	nodes := make([]v8NodesFiniteStringSetNode, len(canonical))
	edges := make([]v8NodesFiniteStringSetEdge, 0)
	for index, item := range canonical {
		nodes[index] = v8NodesFiniteStringSetNode{
			firstEdge: len(edges),
			edgeCount: len(item.edges),
			terminal:  item.terminal,
		}
		edges = append(edges, item.edges...)
	}
	return &v8NodesFiniteStringSet{
		root:   canonicalID[root],
		nodes:  nodes,
		edges:  edges,
		ranges: v8TonodeCanonicalizeRanges(append([]v8ASTCharacterRange(nil), ranges...)),
	}
}

func v8TonodeSortedStringSetKeys(children map[rune]*v8TonodeMutableStringSetNode) []rune {
	keys := make([]rune, 0, len(children))
	for key := range children {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}

func (set *v8NodesFiniteStringSet) firstCodePointRanges() []v8ASTCharacterRange {
	if set == nil || set.root < 0 || set.root >= len(set.nodes) {
		return nil
	}
	root := set.nodes[set.root]
	ranges := append([]v8ASTCharacterRange(nil), set.ranges...)
	for index := root.firstEdge; index < root.firstEdge+root.edgeCount; index++ {
		value := set.edges[index].codePoint
		ranges = append(ranges, v8ASTRangeSingleton(value))
	}
	return v8TonodeCanonicalizeRanges(ranges)
}

func (set *v8NodesFiniteStringSet) hasEmptyString() bool {
	return set != nil && set.root >= 0 && set.root < len(set.nodes) && set.nodes[set.root].terminal
}
