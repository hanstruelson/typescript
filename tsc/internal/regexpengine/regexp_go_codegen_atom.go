package regexpengine

import (
	"fmt"
	"strconv"
	"strings"
)

// v8EmitGoAtomMatcher emits a direct search path for an exact literal. Its
// failure table is generated once with the matcher, so searching the subject
// takes linear time without allocating matcher state.
func v8EmitGoAtomMatcher(atom []uint16, flags Flags, name string) string {
	atomUnits := "_" + name + "AtomUnits"
	atomFailure := "_" + name + "AtomFailure"
	failure := make([]int, len(atom))
	for index := 1; index < len(atom); index++ {
		matched := failure[index-1]
		for matched > 0 && atom[index] != atom[matched] {
			matched = failure[matched-1]
		}
		if atom[index] == atom[matched] {
			matched++
		}
		failure[index] = matched
	}

	var source strings.Builder
	fmt.Fprintf(&source, "var %s = [...]uint16{\n", atomUnits)
	for index, unit := range atom {
		if index%16 == 0 && index != 0 {
			source.WriteByte('\n')
		}
		source.WriteString(strconv.FormatUint(uint64(unit), 10))
		source.WriteByte(',')
	}
	source.WriteString("\n}\n")
	fmt.Fprintf(&source, "var %s = [...]int{\n", atomFailure)
	for index, value := range failure {
		if index%16 == 0 && index != 0 {
			source.WriteByte('\n')
		}
		source.WriteString(strconv.Itoa(value))
		source.WriteByte(',')
	}
	source.WriteString("\n}\n")
	fmt.Fprintf(&source, "const %sCaptureSlots = 2\n", name)
	fmt.Fprintf(&source, `func %sInto(input []uint16, start int, captures []int32) (matchStart, matchEnd int, ok bool) {
if start < 0 { start = 0 }
if start > len(input) || len(captures) < %sCaptureSlots { return 0, 0, false }
`, name, name)
	if (flags.Global || flags.Sticky) && (flags.Unicode || flags.UnicodeSets) {
		source.WriteString(`if start > 0 && start < len(input) && input[start] >= 0xdc00 && input[start] <= 0xdfff && input[start-1] >= 0xd800 && input[start-1] <= 0xdbff { start-- }
`)
	}
	fmt.Fprintf(&source, `matched := 0
for index := start; index < len(input); index++ {
for matched > 0 && input[index] != %s[matched] { matched = %s[matched-1] }
if input[index] == %s[matched] { matched++ }
if matched == len(%s) {
matchStart = index + 1 - len(%s)
matchEnd = index + 1
captures[0] = int32(matchStart)
captures[1] = int32(matchEnd)
return matchStart, matchEnd, true
}
}
return 0, 0, false
}
`, atomUnits, atomFailure, atomUnits, atomUnits, atomUnits)
	fmt.Fprintf(&source, `func %s(input []uint16, start int) (matchStart, matchEnd int, captures []int, ok bool) {
var offsets [2]int32
matchStart, matchEnd, ok = %sInto(input, start, offsets[:])
if !ok { return 0, 0, nil, false }
return matchStart, matchEnd, []int{int(offsets[0]), int(offsets[1])}, true
}
`, name, name)
	return source.String()
}
