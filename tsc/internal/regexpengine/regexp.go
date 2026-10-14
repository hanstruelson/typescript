// Package regexpengine parses ECMAScript regular expressions and emits
// standalone Go matcher functions. It does not execute regular expressions
// inside the compiler process.
package regexpengine

import (
	"sort"
	"unicode/utf16"
)

// VerifySyntax checks flags, grammar, capture declarations, and references
// without generating matcher code.
func VerifySyntax(source, flags string) error {
	return VerifySyntaxUTF16(utf16.Encode([]rune(source)), flags)
}

// VerifySyntaxUTF16 is the UTF-16-preserving syntax-only entry point.
func VerifySyntaxUTF16(source []uint16, flags string) error {
	parsedFlags, err := ParseFlags(flags)
	if err != nil {
		return err
	}
	return v8ASTVerifySyntaxUTF16(source, parsedFlags, nil)
}

// CaptureNamesUTF16 returns capture indices grouped by capture name. Use
// CaptureNamesOrderedUTF16 when property insertion order matters.
func CaptureNamesUTF16(source []uint16, flags string) (map[string][]int, error) {
	ordered, err := CaptureNamesOrderedUTF16(source, flags)
	if err != nil {
		return nil, err
	}
	names := make(map[string][]int, len(ordered))
	for _, capture := range ordered {
		names[capture.Name] = append(names[capture.Name], capture.Indices...)
	}
	return names, nil
}

// NamedCapture preserves V8's capture-index order when building the groups
// object. A name may refer to multiple mutually exclusive captures.
type NamedCapture struct {
	Name    string
	Indices []int
}

// CaptureNamesOrderedUTF16 returns unique named captures ordered by the first
// capture index, matching RegExp::CreateCaptureNameMap's index ordering.
func CaptureNamesOrderedUTF16(source []uint16, flags string) ([]NamedCapture, error) {
	parsedFlags, err := ParseFlags(flags)
	if err != nil {
		return nil, err
	}
	parsed, err := v8ASTParseRegExpUTF16(source, parsedFlags, nil)
	if err != nil {
		return nil, err
	}
	captures := append([]*v8ASTCapture(nil), parsed.NamedCaptures...)
	sort.Slice(captures, func(left, right int) bool {
		return captures[left].index < captures[right].index
	})
	ordered := make([]NamedCapture, 0, len(captures))
	byName := make(map[string]int, len(captures))
	for _, capture := range captures {
		name := utf16String(capture.name)
		if index, ok := byName[name]; ok {
			ordered[index].Indices = append(ordered[index].Indices, capture.index)
			continue
		}
		byName[name] = len(ordered)
		ordered = append(ordered, NamedCapture{Name: name, Indices: []int{capture.index}})
	}
	return ordered, nil
}
