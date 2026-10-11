package regexpengine

// compileLiteralAtom extracts a capture-free sequence of literal characters.
// V8 has a separate atom-search implementation for this subset; the Go port
// keeps the same fast-path boundary and uses KMP over UTF-16 code units.
func compileLiteralAtom(root *node) ([]uint16, []int, bool) {
	units := make([]uint16, 0)
	var appendNode func(*node) bool
	appendNode = func(n *node) bool {
		switch n.kind {
		case nodeEmpty:
			return true
		case nodeLiteral:
			if n.ignoreCase {
				return false
			}
			if n.rune <= 0xffff {
				units = append(units, uint16(n.rune))
			} else {
				value := n.rune - 0x10000
				units = append(units, uint16(0xd800+(value>>10)), uint16(0xdc00+(value&0x3ff)))
			}
			return true
		case nodeConcat:
			for _, child := range n.children {
				if !appendNode(child) {
					return false
				}
			}
			return true
		default:
			return false
		}
	}
	if !appendNode(root) || len(units) == 0 {
		return nil, nil, false
	}

	// Failure links let searches remain linear even for a repetitive pattern
	// such as many 'a' code units followed by 'b'.
	failure := make([]int, len(units))
	for i, matched := 1, 0; i < len(units); i++ {
		for matched > 0 && units[i] != units[matched] {
			matched = failure[matched-1]
		}
		if units[i] == units[matched] {
			matched++
			failure[i] = matched
		}
	}
	return units, failure, true
}

func (p *Program) findLiteralAtom(input []uint16, start int) int {
	return p.findLiteralAtomWithBudget(input, start, nil)
}

func (p *Program) findLiteralAtomWithBudget(input []uint16, start int, budget *executionBudget) int {
	matched := 0
	for i := start; i < len(input); i++ {
		for matched > 0 && input[i] != p.atom[matched] {
			if budget.step() != nil {
				return -1
			}
			matched = p.atomFailure[matched-1]
		}
		if budget.step() != nil {
			return -1
		}
		if input[i] == p.atom[matched] {
			matched++
			if matched == len(p.atom) {
				return i - len(p.atom) + 1
			}
		}
	}
	return -1
}
