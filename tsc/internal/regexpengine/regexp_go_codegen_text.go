package regexpengine

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// emitTextNode lowers a fixed-width V8 TextNode to direct UTF-16 comparisons.
// V8 stores offsets from the start of the text node; reverse-reading nodes
// therefore visit both elements and atom units from right to left.
func (emitter *v8GoGraphEmitter) emitTextNode(node *v8NodesTextNode) error {
	if node.finiteStringSet != nil {
		return emitter.emitFiniteStringSetNode(node)
	}
	if node.ReadBackward() {
		for elementIndex := node.Elements().Len() - 1; elementIndex >= 0; elementIndex-- {
			if err := emitter.emitTextElement(*node.Elements().At(elementIndex), node.Flags(), true); err != nil {
				return err
			}
		}
		return nil
	}
	for elementIndex := 0; elementIndex < node.Elements().Len(); elementIndex++ {
		if err := emitter.emitTextElement(*node.Elements().At(elementIndex), node.Flags(), false); err != nil {
			return err
		}
	}
	return nil
}

func (emitter *v8GoGraphEmitter) emitTextElement(element v8NodesTextElement, flags Flags, backward bool) error {
	switch payload := element.tree.(type) {
	case *v8ASTAtom:
		if payload == nil {
			return fmt.Errorf("nil atom in V8 text graph")
		}
		length := len(payload.data)
		if backward {
			emitter.text.WriteString("if pos < ")
			emitter.text.WriteString(strconv.Itoa(length))
			emitter.text.WriteString(" { goto rxBacktrack }\n")
			for offset := length - 1; offset >= 0; offset-- {
				if err := emitter.emitCodeUnitMatch(payload.data[offset], "pos-"+strconv.Itoa(length-offset), flags); err != nil {
					return err
				}
			}
			emitter.text.WriteString("pos -= ")
			emitter.text.WriteString(strconv.Itoa(length))
			emitter.text.WriteString("\n")
			return nil
		}
		emitter.text.WriteString("if pos < 0 || pos+")
		emitter.text.WriteString(strconv.Itoa(length))
		emitter.text.WriteString(" > len(input) { goto rxBacktrack }\n")
		for offset, unit := range payload.data {
			if err := emitter.emitCodeUnitMatch(unit, "pos+"+strconv.Itoa(offset), flags); err != nil {
				return err
			}
		}
		emitter.text.WriteString("pos += ")
		emitter.text.WriteString(strconv.Itoa(length))
		emitter.text.WriteString("\n")
		return nil
	case *v8ASTClassRanges:
		if payload == nil {
			return fmt.Errorf("nil character class in V8 text graph")
		}
		if IsEitherUnicode(flags) && v8GoClassNeedsCodePoint(payload) {
			return emitter.emitUnicodeCodePointClass(payload, backward)
		}
		index := "pos"
		if backward {
			index = "pos-1"
		}
		condition, err := emitter.classConditionAt(payload, index)
		if err != nil {
			return err
		}
		if payload.set.standardSet != nil && flags.IgnoreCase && IsEitherUnicode(flags) {
			switch *payload.set.standardSet {
			case v8ASTWord:
				condition = v8GoWordConditionAt(index, flags)
			case v8ASTNotWord:
				condition = "!(" + v8GoWordConditionAt(index, flags) + ")"
			}
		}
		if backward {
			emitter.text.WriteString("if pos <= 0 || !( ")
		} else {
			emitter.text.WriteString("if pos < 0 || pos >= len(input) || !( ")
		}
		emitter.text.WriteString(condition)
		emitter.text.WriteString(" ) { goto rxBacktrack }; pos")
		if backward {
			emitter.text.WriteString("--\n")
		} else {
			emitter.text.WriteString("++\n")
		}
		return nil
	default:
		return fmt.Errorf("Go regex codegen does not yet support text element %T", element.tree)
	}
}

// Some V8 graph classes retain Unicode scalar ranges instead of being split
// into UTF-16 lead/trail ranges. Match those as code points while retaining
// UTF-16 offsets for captures and lastIndex. BMP-only graph fragments continue
// through the ordinary one-code-unit path above.
func v8GoClassNeedsCodePoint(class *v8ASTClassRanges) bool {
	if class == nil || class.set.standardSet != nil || len(class.set.ranges) == 0 {
		return false
	}
	ranges := v8TonodeCanonicalizeRanges(class.set.ranges)
	return ranges[len(ranges)-1].to > 0xffff
}

func (emitter *v8GoGraphEmitter) emitUnicodeCodePointClass(class *v8ASTClassRanges, backward bool) error {
	condition, err := emitter.classConditionAt(class, "rxClassRune")
	if err != nil {
		return err
	}
	emitter.text.WriteString("{\n")
	if backward {
		emitter.text.WriteString("if pos <= 0 { goto rxBacktrack }; rxClassRune := uint32(input[pos-1]); rxClassWidth := 1; if input[pos-1] >= 0xdc00 && input[pos-1] <= 0xdfff && pos >= 2 && input[pos-2] >= 0xd800 && input[pos-2] <= 0xdbff { rxClassRune = 0x10000 + ((uint32(input[pos-2])-0xd800)<<10) + (uint32(input[pos-1])-0xdc00); rxClassWidth = 2 }\n")
	} else {
		emitter.text.WriteString("if pos < 0 || pos >= len(input) { goto rxBacktrack }; rxClassRune := uint32(input[pos]); rxClassWidth := 1; if input[pos] >= 0xd800 && input[pos] <= 0xdbff && pos+1 < len(input) && input[pos+1] >= 0xdc00 && input[pos+1] <= 0xdfff { rxClassRune = 0x10000 + ((uint32(input[pos])-0xd800)<<10) + (uint32(input[pos+1])-0xdc00); rxClassWidth = 2 }\n")
	}
	emitter.text.WriteString("if !( ")
	emitter.text.WriteString(condition)
	emitter.text.WriteString(" ) { goto rxBacktrack }; pos")
	if backward {
		emitter.text.WriteString(" -= rxClassWidth\n")
	} else {
		emitter.text.WriteString(" += rxClassWidth\n")
	}
	emitter.text.WriteString("}\n")
	return nil
}

func (emitter *v8GoGraphEmitter) emitCodeUnitMatch(unit uint16, index string, flags Flags) error {
	if !flags.IgnoreCase {
		emitter.text.WriteString("if input[")
		emitter.text.WriteString(index)
		emitter.text.WriteString("] != ")
		emitter.text.WriteString(strconv.Itoa(int(unit)))
		emitter.text.WriteString(" { goto rxBacktrack }\n")
		return nil
	}
	letters := v8CompilerCaseIndependentLetters(flags, unit, false)
	if len(letters) == 0 {
		emitter.text.WriteString("goto rxBacktrack\n")
		return nil
	}
	parts := make([]string, 0, len(letters))
	for _, letter := range letters {
		parts = append(parts, "input["+index+"] == "+strconv.Itoa(int(letter)))
	}
	emitter.text.WriteString("if !( ")
	emitter.text.WriteString(strings.Join(parts, " || "))
	emitter.text.WriteString(" ) { goto rxBacktrack }\n")
	return nil
}

func v8GoWordConditionAt(index string, flags Flags) string {
	character := "input[" + index + "]"
	base := "((" + character + " >= '0' && " + character + " <= '9') || (" + character + " >= 'A' && " + character + " <= 'Z') || " + character + " == '_' || (" + character + " >= 'a' && " + character + " <= 'z'))"
	if !flags.IgnoreCase || !IsEitherUnicode(flags) {
		return base
	}
	extra := make(map[uint16]struct{})
	for letter := uint16('A'); letter <= uint16('Z'); letter++ {
		for _, equivalent := range v8CompilerCaseIndependentLetters(flags, letter, false) {
			if !(equivalent >= 'A' && equivalent <= 'Z') && !(equivalent >= 'a' && equivalent <= 'z') {
				extra[equivalent] = struct{}{}
			}
		}
	}
	for letter := uint16('a'); letter <= uint16('z'); letter++ {
		for _, equivalent := range v8CompilerCaseIndependentLetters(flags, letter, false) {
			if !(equivalent >= 'A' && equivalent <= 'Z') && !(equivalent >= 'a' && equivalent <= 'z') {
				extra[equivalent] = struct{}{}
			}
		}
	}
	if len(extra) == 0 {
		return base
	}
	values := make([]int, 0, len(extra))
	for unit := range extra {
		values = append(values, int(unit))
	}
	sort.Ints(values)
	parts := []string{base}
	for _, value := range values {
		parts = append(parts, character+" == "+strconv.Itoa(value))
	}
	return "(" + strings.Join(parts, " || ") + ")"
}

func (emitter *v8GoGraphEmitter) emitCaseEqualityHelper() {
	neededLegacy := false
	neededUnicode := false
	for _, node := range emitter.nodes {
		if reference, ok := node.(*v8NodesBackReferenceNode); ok && reference.Flags().IgnoreCase {
			if IsEitherUnicode(reference.Flags()) {
				neededUnicode = true
			} else {
				neededLegacy = true
			}
		}
		if text, ok := node.(*v8NodesTextNode); ok && text.finiteStringSet != nil && text.Flags().IgnoreCase {
			neededUnicode = true
		}
	}
	if !neededLegacy && !neededUnicode {
		return
	}
	if neededLegacy {
		emitter.emitLegacyCaseEqualityHelper()
	}
	if neededUnicode {
		emitter.emitUnicodeCaseEqualityHelper()
	}
}

func (emitter *v8GoGraphEmitter) emitLegacyCaseEqualityHelper() {
	groups := make(map[string][]uint16)
	for _, sourceClass := range regexpCaseFoldClasses {
		for _, character := range sourceClass {
			if character < 0 || character > 0xffff {
				continue
			}
			letters := v8CompilerCaseIndependentLetters(emitter.graph.flags, uint16(character), false)
			if len(letters) < 2 {
				continue
			}
			copyLetters := append([]uint16(nil), letters...)
			sort.Slice(copyLetters, func(i, j int) bool { return copyLetters[i] < copyLetters[j] })
			parts := make([]string, len(copyLetters))
			for i, letter := range copyLetters {
				parts[i] = strconv.Itoa(int(letter))
			}
			groups[strings.Join(parts, ",")] = copyLetters
		}
	}
	keys := make([]string, 0, len(groups))
	for key := range groups {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	emitter.text.WriteString("rxCaseEqual := func(a, b uint16) bool { if a == b { return true }; if a >= 'A' && a <= 'Z' && b == a+('a'-'A') { return true }; if a >= 'a' && a <= 'z' && b == a-('a'-'A') { return true }; switch a {\n")
	for _, key := range keys {
		letters := groups[key]
		for _, letter := range letters {
			emitter.text.WriteString("case ")
			emitter.text.WriteString(strconv.Itoa(int(letter)))
			emitter.text.WriteString(": return ")
			for i, candidate := range letters {
				if i != 0 {
					emitter.text.WriteString(" || ")
				}
				emitter.text.WriteString("b == ")
				emitter.text.WriteString(strconv.Itoa(int(candidate)))
			}
			emitter.text.WriteString("\n")
		}
	}
	emitter.text.WriteString("}; return false }\n")
}

// emitUnicodeCaseEqualityHelper ports the code-point comparison used by
// RegExpMacroAssembler::CaseInsensitiveCompareUnicode. The helper's key is
// generated from the pinned C/S simple-fold classes, so it includes
// supplementary-plane characters and does not perform full multi-character
// folding.
func (emitter *v8GoGraphEmitter) emitUnicodeCaseEqualityHelper() {
	type foldClass struct {
		key     rune
		members []rune
	}
	classes := make([]foldClass, 0, len(regexpCaseFoldClasses))
	for _, sourceClass := range regexpCaseFoldClasses {
		if len(sourceClass) < 2 {
			continue
		}
		members := append([]rune(nil), sourceClass...)
		sort.Slice(members, func(i, j int) bool { return members[i] < members[j] })
		classes = append(classes, foldClass{key: members[0], members: members})
	}
	sort.Slice(classes, func(i, j int) bool { return classes[i].key < classes[j].key })
	emitter.text.WriteString("rxUnicodeFoldKey := func(character rune) rune { switch character {\n")
	for _, current := range classes {
		emitter.text.WriteString("case ")
		for index, member := range current.members {
			if index != 0 {
				emitter.text.WriteString(",")
			}
			emitter.text.WriteString(strconv.Itoa(int(member)))
		}
		emitter.text.WriteString(": return ")
		emitter.text.WriteString(strconv.Itoa(int(current.key)))
		emitter.text.WriteString("\n")
	}
	emitter.text.WriteString("}; return character }; rxUnicodeCaseEqual := func(a, b rune) bool { return a == b || a >= 'A' && a <= 'Z' && b == a+('a'-'A') || a >= 'a' && a <= 'z' && b == a-('a'-'A') || rxUnicodeFoldKey(a) == rxUnicodeFoldKey(b) }; _ = rxUnicodeCaseEqual\n")
}

func (emitter *v8GoGraphEmitter) emitV8BackReference(node *v8NodesBackReferenceNode) {
	startRegister, endRegister := node.StartRegister(), node.EndRegister()
	emitter.text.WriteString("{ rxStart := regs[")
	emitter.text.WriteString(strconv.Itoa(startRegister))
	emitter.text.WriteString("]; rxEnd := regs[")
	emitter.text.WriteString(strconv.Itoa(endRegister))
	emitter.text.WriteString("]; if rxStart >= 0 && rxEnd >= 0 { rxLow, rxHigh := rxStart, rxEnd; if rxLow > rxHigh { rxLow, rxHigh = rxHigh, rxLow }; rxLength := rxHigh-rxLow; ")
	if node.Flags().IgnoreCase && IsEitherUnicode(node.Flags()) {
		emitter.emitUnicodeBackReferenceComparison(node.ReadBackward())
	} else if node.ReadBackward() {
		emitter.text.WriteString("if rxLength > pos { goto rxBacktrack }; for rxOffset := 0; rxOffset < rxLength; rxOffset++ { ")
		if node.Flags().IgnoreCase {
			emitter.text.WriteString("if !rxCaseEqual(input[rxHigh-1-rxOffset], input[pos-1-rxOffset]) { goto rxBacktrack } ")
		} else {
			emitter.text.WriteString("if input[rxHigh-1-rxOffset] != input[pos-1-rxOffset] { goto rxBacktrack } ")
		}
		emitter.text.WriteString("}; pos -= rxLength ")
	} else {
		emitter.text.WriteString("if rxLength > len(input)-pos { goto rxBacktrack }; for rxOffset := 0; rxOffset < rxLength; rxOffset++ { ")
		if node.Flags().IgnoreCase {
			emitter.text.WriteString("if !rxCaseEqual(input[rxLow+rxOffset], input[pos+rxOffset]) { goto rxBacktrack } ")
		} else {
			emitter.text.WriteString("if input[rxLow+rxOffset] != input[pos+rxOffset] { goto rxBacktrack } ")
		}
		emitter.text.WriteString("}; pos += rxLength ")
	}
	emitter.text.WriteString("}; ")
	if IsEitherUnicode(node.Flags()) {
		// V8 checks the resulting position even for an unmatched (empty)
		// backreference, because it may otherwise leave the cursor between the
		// two UTF-16 units of a valid surrogate pair.
		emitter.text.WriteString("if pos > 0 && pos < len(input) && input[pos-1] >= 0xd800 && input[pos-1] <= 0xdbff && input[pos] >= 0xdc00 && input[pos] <= 0xdfff { goto rxBacktrack }; ")
	}
	emitter.text.WriteString("}\n")
}

func (emitter *v8GoGraphEmitter) emitUnicodeBackReferenceComparison(backward bool) {
	if backward {
		emitter.text.WriteString("if rxLength > pos { goto rxBacktrack }; rxLeft, rxRight := rxHigh, pos; rxTargetStart := pos-rxLength; for rxLeft > rxLow { if rxRight <= rxTargetStart { goto rxBacktrack }; rxA, rxAWidth := rune(input[rxLeft-1]), 1; if rxA >= 0xdc00 && rxA <= 0xdfff && rxLeft-2 >= rxLow && input[rxLeft-2] >= 0xd800 && input[rxLeft-2] <= 0xdbff { rxA = 0x10000 + ((rune(input[rxLeft-2])-0xd800)<<10) + (rxA-0xdc00); rxAWidth = 2 }; rxB, rxBWidth := rune(input[rxRight-1]), 1; if rxB >= 0xdc00 && rxB <= 0xdfff && rxRight-2 >= rxTargetStart && input[rxRight-2] >= 0xd800 && input[rxRight-2] <= 0xdbff { rxB = 0x10000 + ((rune(input[rxRight-2])-0xd800)<<10) + (rxB-0xdc00); rxBWidth = 2 }; if rxAWidth != rxBWidth || !rxUnicodeCaseEqual(rxA, rxB) { goto rxBacktrack }; rxLeft -= rxAWidth; rxRight -= rxBWidth }; if rxLeft != rxLow || rxRight != rxTargetStart { goto rxBacktrack }; pos = rxTargetStart ")
		return
	}
	emitter.text.WriteString("if rxLength > len(input)-pos { goto rxBacktrack }; rxLeft, rxRight, rxTargetEnd := rxLow, pos, pos+rxLength; for rxLeft < rxHigh { if rxRight >= rxTargetEnd { goto rxBacktrack }; rxA, rxAWidth := rune(input[rxLeft]), 1; if rxA >= 0xd800 && rxA <= 0xdbff && rxLeft+1 < rxHigh && input[rxLeft+1] >= 0xdc00 && input[rxLeft+1] <= 0xdfff { rxA = 0x10000 + ((rxA-0xd800)<<10) + (rune(input[rxLeft+1])-0xdc00); rxAWidth = 2 }; rxB, rxBWidth := rune(input[rxRight]), 1; if rxB >= 0xd800 && rxB <= 0xdbff && rxRight+1 < rxTargetEnd && input[rxRight+1] >= 0xdc00 && input[rxRight+1] <= 0xdfff { rxB = 0x10000 + ((rxB-0xd800)<<10) + (rune(input[rxRight+1])-0xdc00); rxBWidth = 2 }; if rxAWidth != rxBWidth || !rxUnicodeCaseEqual(rxA, rxB) { goto rxBacktrack }; rxLeft += rxAWidth; rxRight += rxBWidth }; if rxLeft != rxHigh || rxRight != rxTargetEnd { goto rxBacktrack }; pos = rxTargetEnd ")
}
