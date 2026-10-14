package regexpengine

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

type v8GoStringSetTableNames struct {
	nodes  string
	edges  string
	ranges string
}

func (emitter *v8GoGraphEmitter) registerFiniteStringSet(set *v8NodesFiniteStringSet) v8GoStringSetTableNames {
	if emitter.stringSetNames == nil {
		emitter.stringSetNames = make(map[*v8NodesFiniteStringSet]v8GoStringSetTableNames)
	}
	if names, exists := emitter.stringSetNames[set]; exists {
		return names
	}
	base := emitter.name + "StringSet" + strconv.Itoa(len(emitter.stringSetNames))
	names := v8GoStringSetTableNames{nodes: base + "Nodes", edges: base + "Edges", ranges: base + "Ranges"}
	var table strings.Builder
	if len(set.nodes)+len(set.edges)+len(set.ranges) > 4096 {
		emitter.emitCompactFiniteStringSetTable(&table, base, names, set)
	} else {
		fmt.Fprintf(&table, "var %s = [...]struct { rxEdgeBase, rxEdgeLength int; rxAccepting bool }{\n", names.nodes)
		for _, node := range set.nodes {
			fmt.Fprintf(&table, "{rxEdgeBase:%d,rxEdgeLength:%d,rxAccepting:%t},\n", node.firstEdge, node.edgeCount, node.terminal)
		}
		table.WriteString("}\n")
		fmt.Fprintf(&table, "var %s = [...]struct { rxRuneValue uint32; nextNode int }{\n", names.edges)
		for _, edge := range set.edges {
			fmt.Fprintf(&table, "{rxRuneValue:%d,nextNode:%d},\n", uint32(edge.codePoint), edge.target)
		}
		table.WriteString("}\n")
		fmt.Fprintf(&table, "var %s = [...]struct { rxRangeStart, rxRangeEnd uint32 }{\n", names.ranges)
		for _, interval := range set.ranges {
			fmt.Fprintf(&table, "{rxRangeStart:%d,rxRangeEnd:%d},\n", uint32(interval.from), uint32(interval.to))
		}
		table.WriteString("}\n")
	}
	emitter.stringSetNames[set] = names
	emitter.stringSetTables = append(emitter.stringSetTables, table.String())
	return names
}

// Large Unicode string properties can have tens of thousands of trie nodes.
// Listing each Go struct literal creates a huge compiler syntax tree and
// drives Go's compile memory far beyond the size of the generated matcher.
// Encode those static tables compactly and decode them once at package init;
// matcher execution still walks the same direct tables without interpretation.
func (emitter *v8GoGraphEmitter) emitCompactFiniteStringSetTable(table *strings.Builder, base string, names v8GoStringSetTableNames, set *v8NodesFiniteStringSet) {
	var data bytes.Buffer
	writeVarint := func(value uint64) {
		var buffer [binary.MaxVarintLen64]byte
		length := binary.PutUvarint(buffer[:], value)
		_, _ = data.Write(buffer[:length])
	}
	for _, node := range set.nodes {
		writeVarint(uint64(node.firstEdge))
		writeVarint(uint64(node.edgeCount))
		if node.terminal {
			writeVarint(1)
		} else {
			writeVarint(0)
		}
	}
	for _, edge := range set.edges {
		writeVarint(uint64(uint32(edge.codePoint)))
		writeVarint(uint64(edge.target))
	}
	for _, interval := range set.ranges {
		writeVarint(uint64(uint32(interval.from)))
		writeVarint(uint64(uint32(interval.to)))
	}

	decoder := emitter.name + "DecodeRegexTable"
	if !emitter.stringSetDecoderEmitted {
		emitter.stringSetDecoderEmitted = true
		fmt.Fprintf(table, `func %s(input string) []byte {
value := func(character byte) uint32 {
switch { case character >= 'A' && character <= 'Z': return uint32(character-'A'); case character >= 'a' && character <= 'z': return uint32(character-'a')+26; case character >= '0' && character <= '9': return uint32(character-'0')+52; case character == '+': return 62; default: return 63 }
}
output := make([]byte, len(input)*3/4+2); written := 0
for index := 0; index < len(input); index += 4 {
a, b := value(input[index]), value(input[index+1]); output[written] = byte(a<<2 | b>>4); written++
if index+2 < len(input) { c := value(input[index+2]); output[written] = byte(b<<4 | c>>2); written++; if index+3 < len(input) { d := value(input[index+3]); output[written] = byte(c<<6 | d); written++ } }
}
return output[:written]
}
`, decoder)
	}
	fmt.Fprintf(table, "const %sData = %q\n", base, base64.RawStdEncoding.EncodeToString(data.Bytes()))
	fmt.Fprintf(table, "var %s []struct { rxEdgeBase, rxEdgeLength int; rxAccepting bool }\n", names.nodes)
	fmt.Fprintf(table, "var %s []struct { rxRuneValue uint32; nextNode int }\n", names.edges)
	fmt.Fprintf(table, "var %s []struct { rxRangeStart, rxRangeEnd uint32 }\n", names.ranges)
	table.WriteString("func init() {\n")
	fmt.Fprintf(table, "data := %s(%sData); cursor := 0\n", decoder, base)
	table.WriteString("read := func() uint32 { var value uint32; var shift uint; for { part := data[cursor]; cursor++; value |= uint32(part&0x7f) << shift; if part < 0x80 { return value }; shift += 7 } }\n")
	fmt.Fprintf(table, "%s = make([]struct { rxEdgeBase, rxEdgeLength int; rxAccepting bool }, %d)\n", names.nodes, len(set.nodes))
	fmt.Fprintf(table, "for index := range %s { %s[index].rxEdgeBase = int(read()); %s[index].rxEdgeLength = int(read()); %s[index].rxAccepting = read() != 0 }\n", names.nodes, names.nodes, names.nodes, names.nodes)
	fmt.Fprintf(table, "%s = make([]struct { rxRuneValue uint32; nextNode int }, %d)\n", names.edges, len(set.edges))
	fmt.Fprintf(table, "for index := range %s { %s[index].rxRuneValue = read(); %s[index].nextNode = int(read()) }\n", names.edges, names.edges, names.edges)
	fmt.Fprintf(table, "%s = make([]struct { rxRangeStart, rxRangeEnd uint32 }, %d)\n", names.ranges, len(set.ranges))
	fmt.Fprintf(table, "for index := range %s { %s[index].rxRangeStart = read(); %s[index].rxRangeEnd = read() }\n", names.ranges, names.ranges, names.ranges)
	table.WriteString("}\n")
}

func (emitter *v8GoGraphEmitter) emitFiniteStringSetNode(node *v8NodesTextNode) error {
	set := node.finiteStringSet
	if set == nil || len(set.nodes) == 0 || set.root < 0 || set.root >= len(set.nodes) {
		return fmt.Errorf("invalid finite string-set automaton")
	}
	names := emitter.registerFiniteStringSet(set)
	seq := node.AsSeqNode()
	if seq == nil || seq.OnSuccess() == nil {
		return fmt.Errorf("finite string-set graph node has no successor")
	}
	target, exists := emitter.nodeIDs[seq.OnSuccess()]
	if !exists {
		return fmt.Errorf("finite string-set successor is absent from graph")
	}
	choiceKey := len(emitter.choiceCases) + 1
	choiceLabel := emitter.nodeIDs[node]
	emitter.choiceCases = append(emitter.choiceCases, v8GoChoiceCase{
		key: choiceKey, target: target, choiceLabel: choiceLabel,
	})
	nameID := strconv.Itoa(choiceLabel)
	emitter.text.WriteString("{\n")
	emitter.text.WriteString("rxSetLast")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" := -1\n")
	if set.hasEmptyString() {
		emitter.text.WriteString("rxSetLast")
		emitter.text.WriteString(nameID)
		emitter.text.WriteString(" = pos\n")
	}

	if len(set.ranges) != 0 {
		emitter.emitFiniteStringSetRangeMatch(nameID, names, choiceKey, node.Flags())
	}
	emitter.emitFiniteStringSetTrieWalk(nameID, names, choiceKey, set.root, node.Flags())
	emitter.text.WriteString("if rxSetLast")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" < 0 { goto rxBacktrack }; pos = rxSetLast")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString("\n}\n")
	return nil
}

func (emitter *v8GoGraphEmitter) emitFiniteStringSetRangeMatch(nameID string, names v8GoStringSetTableNames, choiceKey int, flags Flags) {
	emitter.text.WriteString("if pos >= 0 && pos < len(input) { rxSetRangeRune")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" := uint32(input[pos]); rxSetRangeWidth")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" := 1; if rxSetRangeRune")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" >= 0xd800 && rxSetRangeRune")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" <= 0xdbff && pos+1 < len(input) && input[pos+1] >= 0xdc00 && input[pos+1] <= 0xdfff { rxSetRangeRune")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" = 0x10000 + ((rxSetRangeRune")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString("-0xd800)<<10) + (uint32(input[pos+1])-0xdc00); rxSetRangeWidth")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" = 2 }; rxSetRangeLow")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(", rxSetRangeHigh")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" := 0, len(")
	emitter.text.WriteString(names.ranges)
	emitter.text.WriteString("); for rxSetRangeLow")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" < rxSetRangeHigh")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" { rxSetRangeMiddle")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" := rxSetRangeLow")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString("+(rxSetRangeHigh")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString("-rxSetRangeLow")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(")/2; if rxSetRangeRune")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" > ")
	emitter.text.WriteString(names.ranges)
	emitter.text.WriteString("[rxSetRangeMiddle")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString("].rxRangeEnd { rxSetRangeLow")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" = rxSetRangeMiddle")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString("+1 } else { rxSetRangeHigh")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" = rxSetRangeMiddle")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" } }; if rxSetRangeLow")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" < len(")
	emitter.text.WriteString(names.ranges)
	emitter.text.WriteString(") && ")
	emitter.text.WriteString(names.ranges)
	emitter.text.WriteString("[rxSetRangeLow")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString("].rxRangeStart <= rxSetRangeRune")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" {\n")
	if flags.IgnoreCase {
		// Ranges in Unicode-set operands are already case-fold closed by the
		// parser; the raw input code point is therefore tested here.
	}
	emitter.emitFiniteStringSetPush(choiceKey, "rxSetLast"+nameID, "Range"+nameID)
	emitter.text.WriteString("rxSetLast")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" = pos+rxSetRangeWidth")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString("\n}\n}\n")
}

func (emitter *v8GoGraphEmitter) emitFiniteStringSetTrieWalk(nameID string, names v8GoStringSetTableNames, choiceKey, root int, flags Flags) {
	emitter.text.WriteString("rxSetState")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" := ")
	emitter.text.WriteString(strconv.Itoa(root))
	emitter.text.WriteString("; rxSetPos")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" := pos\nfor rxSetPos")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" >= 0 && rxSetPos")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" < len(input) {\nrxSetRune")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" := uint32(input[rxSetPos")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString("]); rxSetWidth")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" := 1; if rxSetRune")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" >= 0xd800 && rxSetRune")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" <= 0xdbff && rxSetPos")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString("+1 < len(input) && input[rxSetPos")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString("+1] >= 0xdc00 && input[rxSetPos")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString("+1] <= 0xdfff { rxSetRune")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" = 0x10000 + ((rxSetRune")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString("-0xd800)<<10) + (uint32(input[rxSetPos")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString("+1])-0xdc00); rxSetWidth")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" = 2 }; ")
	if flags.IgnoreCase {
		emitter.text.WriteString("rxSetRune")
		emitter.text.WriteString(nameID)
		emitter.text.WriteString(" = uint32(rxUnicodeFoldKey(rune(rxSetRune")
		emitter.text.WriteString(nameID)
		emitter.text.WriteString("))) ")
	}
	emitter.text.WriteString("rxSetNode")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" := ")
	emitter.text.WriteString(names.nodes)
	emitter.text.WriteString("[rxSetState")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString("]; rxSetEdgeLow")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" := rxSetNode")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(".rxEdgeBase; rxSetEdgeHigh")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" := rxSetEdgeLow")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString("+rxSetNode")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(".rxEdgeLength; rxSetEdgeEnd")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" := rxSetEdgeHigh")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString("; for rxSetEdgeLow")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" < rxSetEdgeHigh")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" { rxSetEdgeMiddle")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" := rxSetEdgeLow")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString("+(rxSetEdgeHigh")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString("-rxSetEdgeLow")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(")/2; if rxSetRune")
	emitter.text.WriteString(nameID)
	// Find the first edge whose value is greater than or equal to the input
	// code point; retain the original end because the search converges low/high.
	emitter.text.WriteString(" <= ")
	emitter.text.WriteString(names.edges)
	emitter.text.WriteString("[rxSetEdgeMiddle")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString("].rxRuneValue { rxSetEdgeHigh")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" = rxSetEdgeMiddle")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" } else { rxSetEdgeLow")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" = rxSetEdgeMiddle")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString("+1 } }; if rxSetEdgeLow")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" >= rxSetEdgeEnd")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" || ")
	emitter.text.WriteString(names.edges)
	emitter.text.WriteString("[rxSetEdgeLow")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString("].rxRuneValue != rxSetRune")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" { break }; rxSetState")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" = ")
	emitter.text.WriteString(names.edges)
	emitter.text.WriteString("[rxSetEdgeLow")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString("].nextNode; rxSetPos")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" += rxSetWidth")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString("; if ")
	emitter.text.WriteString(names.nodes)
	emitter.text.WriteString("[rxSetState")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString("].rxAccepting {\n")
	emitter.emitFiniteStringSetPush(choiceKey, "rxSetLast"+nameID, "Terminal"+nameID)
	emitter.text.WriteString("rxSetLast")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString(" = rxSetPos")
	emitter.text.WriteString(nameID)
	emitter.text.WriteString("\n}\n}\n")
}

func (emitter *v8GoGraphEmitter) emitFiniteStringSetPush(choiceKey int, position, suffix string) {
	frameSlots := max(emitter.graph.registerCount, 2) + 2
	key := strconv.Itoa(choiceKey)
	emitter.text.WriteString("if ")
	emitter.text.WriteString(position)
	emitter.text.WriteString(" >= 0 { if len(stack)+rxStackFrameSlots > rxMaxStackSlots { return -1, -1, false }; if len(stack)+rxStackFrameSlots > cap(stack) { rxSetStackCap")
	emitter.text.WriteString(suffix)
	emitter.text.WriteString(" := cap(stack)*2; if rxSetStackCap")
	emitter.text.WriteString(suffix)
	emitter.text.WriteString(" < len(stack)+rxStackFrameSlots { rxSetStackCap")
	emitter.text.WriteString(suffix)
	emitter.text.WriteString(" = len(stack)+rxStackFrameSlots }; if rxSetStackCap")
	emitter.text.WriteString(suffix)
	emitter.text.WriteString(" > rxMaxStackSlots { rxSetStackCap")
	emitter.text.WriteString(suffix)
	emitter.text.WriteString(" = rxMaxStackSlots }; rxSetGrownStack")
	emitter.text.WriteString(suffix)
	emitter.text.WriteString(" := make([]int32, len(stack), rxSetStackCap")
	emitter.text.WriteString(suffix)
	emitter.text.WriteString("); copy(rxSetGrownStack")
	emitter.text.WriteString(suffix)
	emitter.text.WriteString(", stack); stack = rxSetGrownStack")
	emitter.text.WriteString(suffix)
	emitter.text.WriteString(" }; stack = append(stack, int32(")
	emitter.text.WriteString(key)
	emitter.text.WriteString("), int32(")
	emitter.text.WriteString(position)
	emitter.text.WriteString(")")
	for register := 0; register < frameSlots-2; register++ {
		emitter.text.WriteString(", int32(regs[")
		emitter.text.WriteString(strconv.Itoa(register))
		emitter.text.WriteString("]) ")
	}
	emitter.text.WriteString(") }\n")
}
