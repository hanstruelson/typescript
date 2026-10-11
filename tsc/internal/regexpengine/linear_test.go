package regexpengine

import (
	"reflect"
	"strings"
	"testing"
	"unicode/utf16"
)

func compileLinearForTest(t *testing.T, pattern, flags string) (*linearProgram, Flags) {
	t.Helper()
	parsedFlags, err := ParseFlags(flags)
	if err != nil {
		t.Fatal(err)
	}
	parsedFlags.Linear = true
	units := utf16.Encode([]rune(pattern))
	parser := &parser{text: units, flags: parsedFlags, groupNames: map[string][]int{}}
	root, err := parser.parse()
	if err != nil {
		t.Fatal(err)
	}
	root, err = resolveReferences(root, parser)
	if err != nil {
		t.Fatal(err)
	}
	program, err := compileLinear(root, parsedFlags, parser.group)
	if err != nil {
		t.Fatal(err)
	}
	return program, parsedFlags
}

func TestLinearPriorityAndRepetition(t *testing.T) {
	tests := []struct {
		pattern, input            string
		start, wantStart, wantEnd int
	}{
		{"a+", "zaaa!", 0, 1, 4},
		{"a|ab", "ab", 0, 0, 1},
		{"a+?", "aaa", 0, 0, 1},
		{"a*", "bbb", 0, 0, 0},
	}
	for _, test := range tests {
		program, _ := compileLinearForTest(t, test.pattern, "")
		match, ok := program.FindUTF16(utf16.Encode([]rune(test.input)), test.start)
		if !ok || match.Start != test.wantStart || match.End != test.wantEnd {
			t.Errorf("linear /%s/ on %q = %#v, %v; want [%d,%d)", test.pattern, test.input, match, ok, test.wantStart, test.wantEnd)
		}
	}
}

func TestLinearScopedDotAllAndMultilineModifiers(t *testing.T) {
	for _, test := range []struct {
		pattern, input string
		want           bool
	}{
		{`(?s:.)`, "\n", true},
		{`(?-s:.)`, "\n", false},
		{`(?m:^b)`, "x\nb", true},
		{`(?-m:^b)`, "x\nb", false},
	} {
		program, err := Compile(test.pattern, "l")
		if err != nil {
			t.Fatalf("Compile(%q, l): %v", test.pattern, err)
		}
		_, matched := program.Find(test.input, 0)
		if matched != test.want {
			t.Errorf("/%s/l on %q matched %v; want %v", test.pattern, test.input, matched, test.want)
		}
	}
}

func TestLinearRejectsScopedIgnoreCase(t *testing.T) {
	if _, err := Compile(`(?i:a)`, "l"); err == nil {
		t.Fatal("the V8 linear-engine subset accepted scoped ignore-case")
	}
}

func TestLinearCapturesClearAcrossIterations(t *testing.T) {
	program, _ := compileLinearForTest(t, `((a)?b)+`, "")
	match, ok := program.FindUTF16(utf16.Encode([]rune("abb")), 0)
	if !ok || match.Start != 0 || match.End != 3 {
		t.Fatalf("linear match = %#v, %v", match, ok)
	}
	if got := match.Captures[1]; got != (Span{Start: 2, End: 3}) {
		t.Fatalf("outer capture = %#v; want last repetition [2,3)", got)
	}
	if got := match.Captures[2]; got != (Span{Start: -1, End: -1}) {
		t.Fatalf("skipped inner capture = %#v; want unmatched", got)
	}
}

func TestLinearEmptyQuantifierIterationMatchesV8LinearEngine(t *testing.T) {
	program, err := Compile(`(a?)*b`, "l")
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		input   string
		capture Span
	}{
		{input: "b", capture: Span{Start: -1, End: -1}},
		{input: "ab", capture: Span{Start: 0, End: 1}},
	} {
		match, ok := program.Find(test.input, 0)
		if !ok || match.Captures[1] != test.capture {
			t.Errorf("V8-compatible /(a?)*b/l on %q = %#v, %v; capture want %#v", test.input, match, ok, test.capture)
		}
	}
}

func TestLinearNullableQuantifierCapturesMatchV8(t *testing.T) {
	undefined := Span{Start: -1, End: -1}
	tests := []struct {
		pattern, input string
		want           []Span
	}{
		{`((a?)*)b`, "b", []Span{{0, 1}, {0, 0}, undefined}},
		{`((a?)*)b`, "ab", []Span{{0, 2}, {0, 1}, {0, 1}}},
		{`((a?)*)b`, "aab", []Span{{0, 3}, {0, 2}, {1, 2}}},
		{`(a?)+b`, "b", []Span{{0, 1}, {0, 0}}},
		{`(a?)+b`, "aab", []Span{{0, 3}, {1, 2}}},
		{`(a?){2}b`, "ab", []Span{{0, 2}, {1, 1}}},
		{`(a?){2}b`, "aab", []Span{{0, 3}, {1, 2}}},
		{`((a)?b)+`, "b", []Span{{0, 1}, {0, 1}, undefined}},
		{`((a)?b)+`, "abb", []Span{{0, 3}, {2, 3}, undefined}},
		{`((a)?b)+`, "aab", []Span{{1, 3}, {1, 3}, {1, 2}}},
	}
	for _, test := range tests {
		program, err := Compile(test.pattern, "l")
		if err != nil {
			t.Fatalf("Compile(%q, l): %v", test.pattern, err)
		}
		match, ok := program.Find(test.input, 0)
		if !ok || !reflect.DeepEqual(match.Captures, test.want) {
			t.Errorf("V8 /%s/l on %q = %#v, %v; captures want %#v", test.pattern, test.input, match, ok, test.want)
		}
	}
}

func TestLinearStickyAndUnsupportedFeatures(t *testing.T) {
	program, _ := compileLinearForTest(t, "a+", "y")
	input := utf16.Encode([]rune("ba"))
	if _, ok := program.FindUTF16(input, 0); ok {
		t.Fatal("sticky linear match searched past its start")
	}
	match, ok := program.FindUTF16(input, 1)
	if !ok || match.Start != 1 || match.End != 2 {
		t.Fatalf("sticky match = %#v, %v; want [1,2)", match, ok)
	}

	for _, pattern := range []string{`(a)\1`, `(?=a)a`, `(?<=(a))b`} {
		if _, err := Compile(pattern, "l"); err == nil {
			t.Errorf("linear engine accepted unsupported /%s/l", pattern)
		}
	}
	for _, flags := range []string{"gl", "ly"} {
		if _, err := Compile(`(?<=a)b`, flags); err == nil {
			t.Errorf("linear engine accepted lookbehind with unsupported flags %q", flags)
		}
	}
	if _, err := Compile(`a{17}`, "l"); err == nil {
		t.Fatal("linear engine accepted a quantifier beyond V8's replication limit")
	}
}

func TestLinearCapturelessLookbehind(t *testing.T) {
	tests := []struct {
		pattern, input string
		wantStart      int
		wantMatch      bool
	}{
		{`(?<=a)b`, "ab", 1, true},
		{`(?<!a)b`, "cb", 1, true},
		{`(?<!a)b`, "ab", 0, false},
		{`(?<=a{1,3})b`, "xxaaab", 5, true},
		{`(?<=(?<!x)a)b`, "ab", 1, true},
		{`(?<=(?<!x)a)b`, "xab", 0, false},
	}
	for _, test := range tests {
		program, err := Compile(test.pattern, "l")
		if err != nil {
			t.Fatalf("Compile(%q, l): %v", test.pattern, err)
		}
		match, ok := program.Find(test.input, 0)
		if ok != test.wantMatch || ok && match.Start != test.wantStart {
			t.Errorf("/%s/l on %q = %#v, %v; want match=%v at %d", test.pattern, test.input, match, ok, test.wantMatch, test.wantStart)
		}
	}
}

func TestLinearUnboundedLookbehindScansLongInput(t *testing.T) {
	program, err := Compile(`(?<=a*)b`, "l")
	if err != nil {
		t.Fatal(err)
	}
	input := strings.Repeat("a", 10000) + "b"
	match, ok := program.Find(input, 0)
	if !ok || match.Start != len(input)-1 || match.End != len(input) {
		t.Fatalf("unbounded lookbehind match = %#v, %v", match, ok)
	}
}

func TestLinearMatchesExistingInterpreterOnSupportedCorpus(t *testing.T) {
	// Zero-width quantifier iterations have intentionally different capture
	// behavior in V8's experimental l engine and are covered by a dedicated test.
	patterns := []string{
		`a+`, `a+?`, `a|ab`, `ab|a`, `(a|ab)c?`,
		`([ab]+)(a?)`, `(a*)b`, `(?:a|b){1,3}c`,
		`^a+$`, `a$`, `\ba+\b`, `(?<=a)b`, `(?<!a)b`, `(?<=a{1,3})b`,
	}
	inputs := []string{"", "a", "b", "ab", "aab", "ba", "abb", "aaab", "xxabx", "aaaabc", "ababba", "\na\n"}
	for _, pattern := range patterns {
		regular, err := Compile(pattern, "")
		if err != nil {
			t.Fatalf("Compile(%q): %v", pattern, err)
		}
		linear, err := Compile(pattern, "l")
		if err != nil {
			t.Fatalf("Compile(%q, l): %v", pattern, err)
		}
		for _, input := range inputs {
			want, wantOK := regular.Find(input, 0)
			got, gotOK := linear.Find(input, 0)
			if gotOK != wantOK || gotOK && !reflect.DeepEqual(got, want) {
				t.Errorf("/%s/ on %q: linear=%#v,%v standard=%#v,%v", pattern, input, got, gotOK, want, wantOK)
			}
		}
	}
}

func TestCompileLinearFlagUsesPrioritizedNFA(t *testing.T) {
	program, err := Compile(`(a|aa)*b`, "l")
	if err != nil {
		t.Fatal(err)
	}
	if program.linear == nil {
		t.Fatal("l flag did not select the linear engine")
	}
	input := strings.Repeat("a", 5000)
	if _, ok := program.Find(input, 0); ok {
		t.Fatal("pattern unexpectedly matched without its final b")
	}
}

func TestUnicodeGlobalStickyStartRewindsInsideSurrogatePair(t *testing.T) {
	input := []uint16{0xd83d, 0xde00}
	for _, flags := range []string{"uy", "ug", "vy", "vg"} {
		program, err := Compile(".", flags)
		if err != nil {
			t.Fatalf("Compile(%q): %v", flags, err)
		}
		match, ok := program.FindUTF16(input, 1)
		if !ok || match.Start != 0 || match.End != 2 {
			t.Errorf("/%s/ from index 1 = %#v, %v; want [0,2)", flags, match, ok)
		}
	}
}
