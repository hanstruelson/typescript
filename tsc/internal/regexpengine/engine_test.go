package regexpengine

import (
	"reflect"
	"strings"
	"testing"
	"unicode/utf16"
)

func TestCorePatterns(t *testing.T) {
	tests := []struct {
		pattern, flags, input     string
		start, wantStart, wantEnd int
	}{
		{"a+", "", "xxaa!", 0, 2, 4},
		{"a|bc", "", "xxbc", 0, 2, 4},
		{"(a)(b)", "", "zab", 0, 1, 3},
		{"a*", "g", "ba", 0, 0, 0},
		{"^a$", "", "a", 0, 0, 1},
		{`\ba\b`, "", "-a", 0, 1, 2},
		{"^a$", "m", "x\na\ny", 0, 2, 3},
		{`\p{Letter}+`, "u", "123é", 0, 3, 4},
		{"a", "y", "ba", 1, 1, 2},
		{"[a-c]+", "i", "xxBC", 0, 2, 4},
	}
	for _, test := range tests {
		program, err := Compile(test.pattern, test.flags)
		if err != nil {
			t.Fatalf("Compile(%q): %v", test.pattern, err)
		}
		match, ok := program.Find(test.input, test.start)
		if !ok || match.Start != test.wantStart || match.End != test.wantEnd {
			t.Fatalf("/%s/%s on %q = %#v, %v; want [%d,%d)", test.pattern, test.flags, test.input, match, ok, test.wantStart, test.wantEnd)
		}
	}
}

func TestBoundedExecutionReportsResourceLimits(t *testing.T) {
	defaultLimits := DefaultExecutionLimits()
	for _, test := range []struct {
		name    string
		pattern string
		input   string
		limits  ExecutionLimits
		want    string
	}{
		{
			name:    "step budget",
			pattern: `a+`,
			input:   `aaaa`,
			limits:  ExecutionLimits{MaxSteps: 1, MaxBacktrackBytes: defaultLimits.MaxBacktrackBytes, MaxAssertionDepth: defaultLimits.MaxAssertionDepth},
			want:    "step",
		},
		{
			name:    "backtracking snapshot memory",
			pattern: `([ab]|c)c`,
			input:   `ac`,
			limits:  ExecutionLimits{MaxSteps: defaultLimits.MaxSteps, MaxBacktrackBytes: 300, MaxAssertionDepth: defaultLimits.MaxAssertionDepth},
			want:    "backtracking memory",
		},
		{
			name:    "nested assertions",
			pattern: `(?=(?=a)a)a`,
			input:   `a`,
			limits:  ExecutionLimits{MaxSteps: defaultLimits.MaxSteps, MaxBacktrackBytes: defaultLimits.MaxBacktrackBytes, MaxAssertionDepth: 1},
			want:    "assertion depth",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			program, err := Compile(test.pattern, "")
			if err != nil {
				t.Fatal(err)
			}
			_, matched, err := program.FindWithLimits(test.input, 0, test.limits)
			if err == nil || matched {
				t.Fatalf("bounded match = %v, %v; want a resource-limit error", matched, err)
			}
			limit, ok := err.(*ExecutionLimitError)
			if !ok || limit.Resource != test.want {
				t.Fatalf("bounded match error = %T %v; want %q resource limit", err, err, test.want)
			}
		})
	}
}

func TestBoundedExecutionCanRetryWithLinearSubset(t *testing.T) {
	program, err := Compile(`^(a|aa)+$`, ``)
	if err != nil {
		t.Fatal(err)
	}
	limits := DefaultExecutionLimits()
	limits.MaxSteps = 10_000
	input := strings.Repeat(`a`, 40) + `!`
	if _, matched, err := program.FindWithLimits(input, 0, limits); err == nil || matched {
		t.Fatalf("ordinary bounded search = matched %v, err %v; want resource limit", matched, err)
	}
	if _, matched, err := program.FindWithLinearFallback(input, 0, limits); err != nil || matched {
		t.Fatalf("linear retry = matched %v, err %v; want bounded no-match", matched, err)
	}

	canMatch, err := Compile(`^(a|aa)+b|a+$`, ``)
	if err != nil {
		t.Fatal(err)
	}
	input = strings.Repeat(`a`, 40)
	if _, matched, err := canMatch.FindWithLimits(input, 0, limits); err == nil || matched {
		t.Fatalf("ordinary bounded successful-pattern search = matched %v, err %v; want resource limit", matched, err)
	}
	match, matched, err := canMatch.FindWithLinearFallback(input, 0, limits)
	if err != nil || !matched || match.Start != 0 || match.End != len(input) || match.Captures[1] != (Span{-1, -1}) {
		t.Fatalf("linear retry successful match = %#v, %v, %v; want [0,%d) with first-branch capture unset", match, matched, err, len(input))
	}

	// A backreference is outside the linear subset, so fallback must preserve
	// the original resource error instead of pretending it was a no-match.
	unsupported, err := Compile(`^(a)\1$`, ``)
	if err != nil {
		t.Fatal(err)
	}
	limits.MaxSteps = 1
	_, matched, err = unsupported.FindWithLinearFallback(`aa`, 0, limits)
	if matched {
		t.Fatal("unsupported fallback unexpectedly matched")
	}
	if limit, ok := err.(*ExecutionLimitError); !ok || limit.Resource != "step" {
		t.Fatalf("unsupported fallback error = %T %v; want original step limit", err, err)
	}
}

func TestBoundedFindAllSharesBudgetAcrossGlobalMatches(t *testing.T) {
	program, err := Compile(`a`, `g`)
	if err != nil {
		t.Fatal(err)
	}
	limits := DefaultExecutionLimits()
	limits.MaxSteps = 3
	if _, err := program.FindAllWithLimits(`aaaa`, limits); err == nil {
		t.Fatal("global find-all reset its step budget for each match")
	}
	limits = DefaultExecutionLimits()
	matches, err := program.FindAllWithLimits(`aaaa`, limits)
	if err != nil || len(matches) != 4 {
		t.Fatalf("bounded global matches = %#v, %v", matches, err)
	}

	linear, err := Compile(`a`, `gl`)
	if err != nil {
		t.Fatal(err)
	}
	limits = DefaultExecutionLimits()
	limits.MaxSteps = 30
	if _, err := linear.FindAllWithLimits(`aaaa`, limits); err == nil {
		t.Fatal("linear global find-all reset its conservative work budget per match")
	}
}

func TestBoundedSearchPreservesOrdinaryMatchResults(t *testing.T) {
	for _, test := range []struct {
		pattern, flags, input string
		start                 int
	}{
		{pattern: `a|ab`, input: `ab`},
		{pattern: `(?<=a)(b+)`, flags: `u`, input: `zabbb`},
		{pattern: `([a-z]+)-(\d+)`, input: `xabc-123y`},
		{pattern: `[\q{a|ab}]b`, flags: `v`, input: `ab`},
		{pattern: `a*`, flags: `gu`, input: `😀a`},
		{pattern: `needle`, flags: `g`, input: `xxneedle`, start: 1},
	} {
		program, err := Compile(test.pattern, test.flags)
		if err != nil {
			t.Fatalf("Compile(%q, %q): %v", test.pattern, test.flags, err)
		}
		want, wantOK := program.Find(test.input, test.start)
		got, gotOK, err := program.FindWithLimits(test.input, test.start, DefaultExecutionLimits())
		if err != nil || wantOK != gotOK || !reflect.DeepEqual(want, got) {
			t.Errorf("bounded /%s/%s on %q = %#v, %v, %v; want %#v, %v", test.pattern, test.flags, test.input, got, gotOK, err, want, wantOK)
		}
		if program.Flags.Global {
			wantAll := program.FindAll(test.input)
			gotAll, err := program.FindAllWithLimits(test.input, DefaultExecutionLimits())
			if err != nil || !reflect.DeepEqual(wantAll, gotAll) {
				t.Errorf("bounded global /%s/ on %q = %#v, %v; want %#v", test.pattern, test.input, gotAll, err, wantAll)
			}
		}
	}
}

func TestFirstCharacterFilterPreservesUnfilteredResults(t *testing.T) {
	for _, test := range []struct {
		pattern, flags, input string
		wantFilter            bool
	}{
		{pattern: `abc|def`, input: `xxdef`, wantFilter: true},
		{pattern: `(?=a)a|b`, input: `xb`, wantFilter: true},
		{pattern: `(?<=a)b`, input: `zab`, wantFilter: true},
		{pattern: `(?:a)?b`, input: `xxab`, wantFilter: true},
		{pattern: `a*b`, input: `xxaaab`, wantFilter: true},
		{pattern: `(?i:[a-z])x`, input: `Qx`, wantFilter: true},
		{pattern: `\p{Letter}x`, flags: `u`, input: `!éx`, wantFilter: true},
		{pattern: `😀x`, flags: `u`, input: `x😀x`, wantFilter: true},
		{pattern: `a*`, input: `bbb`, wantFilter: false},
		{pattern: `.*z`, input: `abz`, wantFilter: false},
		{pattern: `(?<x>a)?\k<x>b`, input: `ab`, wantFilter: false},
		{pattern: `[\q{ab|cd}]x`, flags: `v`, input: `cdx`, wantFilter: false},
	} {
		program, err := Compile(test.pattern, test.flags)
		if err != nil {
			t.Fatalf("Compile(%q, %q): %v", test.pattern, test.flags, err)
		}
		if got := len(program.firstChars) != 0; got != test.wantFilter {
			t.Errorf("/%s/%s filter present = %v, want %v", test.pattern, test.flags, got, test.wantFilter)
		}
		withoutFilter := *program
		withoutFilter.firstChars = nil
		want, wantOK := withoutFilter.Find(test.input, 0)
		got, gotOK := program.Find(test.input, 0)
		if gotOK != wantOK || !reflect.DeepEqual(got, want) {
			t.Errorf("filtered /%s/%s on %q = %#v, %v; unfiltered = %#v, %v", test.pattern, test.flags, test.input, got, gotOK, want, wantOK)
		}
		bounded, boundedOK, err := program.FindWithLimits(test.input, 0, DefaultExecutionLimits())
		if err != nil || boundedOK != wantOK || !reflect.DeepEqual(bounded, want) {
			t.Errorf("bounded filtered /%s/%s on %q = %#v, %v, %v; unfiltered = %#v, %v", test.pattern, test.flags, test.input, bounded, boundedOK, err, want, wantOK)
		}
	}
}

func TestGreedySingleCharacterRepeatsUseOneRetryFrame(t *testing.T) {
	for _, test := range []struct {
		pattern, flags, input string
		want                  Match
	}{
		{pattern: `a+`, input: `zaaa!`, want: Match{Start: 1, End: 4, Captures: []Span{{1, 4}}}},
		{pattern: `a{2,}b`, input: `zaaab!`, want: Match{Start: 1, End: 5, Captures: []Span{{1, 5}}}},
		{pattern: `([a-z]+)([0-9])`, input: `xabc7!`, want: Match{Start: 0, End: 5, Captures: []Span{{0, 5}, {0, 4}, {4, 5}}}},
		{pattern: `(?<x>a*)a`, input: `aaa`, want: Match{Start: 0, End: 3, Captures: []Span{{0, 3}, {0, 2}}}},
		{pattern: `(?<=a+)b`, input: `aaab`, want: Match{Start: 3, End: 4, Captures: []Span{{3, 4}}}},
		{pattern: `(?:😀)+x`, flags: `u`, input: `x😀😀x`, want: Match{Start: 1, End: 6, Captures: []Span{{1, 6}}}},
	} {
		program, err := Compile(test.pattern, test.flags)
		if err != nil {
			t.Fatalf("Compile(%q, %q): %v", test.pattern, test.flags, err)
		}
		got, ok := program.Find(test.input, 0)
		if !ok || !reflect.DeepEqual(got, test.want) {
			t.Errorf("/%s/%s on %q = %#v, %v; want %#v", test.pattern, test.flags, test.input, got, ok, test.want)
		}
	}

	program, err := Compile(`a*b`, ``)
	if err != nil {
		t.Fatal(err)
	}
	optimized := false
	for _, ins := range program.code {
		optimized = optimized || ins.op == opGreedyLoop
	}
	if !optimized {
		t.Fatal("greedy single-character repetition did not use the compact loop instruction")
	}

	limits := DefaultExecutionLimits()
	limits.MaxBacktrackBytes = 512
	input := strings.Repeat(`a`, 50000) + `b`
	match, ok, err := program.FindWithLimits(input, 0, limits)
	if err != nil || !ok || match.End != len(input) {
		t.Fatalf("bounded long greedy repeat = %#v, %v, %v; expected one compact retry frame", match, ok, err)
	}
}

func TestLiteralAlternationLoweringPreservesPriority(t *testing.T) {
	single, err := Compile(`a|b|c`, ``)
	if err != nil {
		t.Fatal(err)
	}
	if len(single.code) == 0 || single.code[0].op != opClass {
		t.Fatalf("single-character alternatives were not lowered to a class: %#v", single.code)
	}
	for _, input := range []string{`a`, `b`, `c`} {
		if match, ok := single.Find(input, 0); !ok || match.End != 1 {
			t.Errorf("/a|b|c/ on %q = %#v, %v", input, match, ok)
		}
	}

	shared, err := Compile(`ab|ac`, ``)
	if err != nil {
		t.Fatal(err)
	}
	if len(shared.code) < 2 || shared.code[0].op != opChar || shared.code[0].rune != 'a' || shared.code[1].op != opClass {
		t.Fatalf("shared-prefix alternatives were not factored: %#v", shared.code)
	}
	if match, ok := shared.Find(`ac`, 0); !ok || match.End != 2 {
		t.Errorf("/ab|ac/ on ac = %#v, %v", match, ok)
	}

	priority, err := Compile(`a|ab`, ``)
	if err != nil {
		t.Fatal(err)
	}
	if match, ok := priority.Find(`ab`, 0); !ok || match.End != 1 {
		t.Errorf("factored /a|ab/ changed left-alternative priority: %#v, %v", match, ok)
	}

	captured, err := Compile(`(a|ab)c`, ``)
	if err != nil {
		t.Fatal(err)
	}
	match, ok := captured.Find(`abc`, 0)
	if !ok || match.End != 3 || match.Captures[1] != (Span{0, 2}) {
		t.Errorf("factored captured alternative = %#v, %v; want group [0,2)", match, ok)
	}
}

func TestLiteralAtomFastPathMatchesUTF16SearchSemantics(t *testing.T) {
	program, err := Compile(`needle`, `g`)
	if err != nil {
		t.Fatal(err)
	}
	if len(program.atom) == 0 || len(program.code) != 0 {
		t.Fatal("capture-free literal did not select the atom-search path")
	}
	match, ok := program.Find("xxneedle needle", 1)
	if !ok || match.Start != 2 || match.End != 8 || len(match.Captures) != 1 || match.Captures[0] != (Span{2, 8}) {
		t.Fatalf("literal fast-path match = %#v, %v", match, ok)
	}
	matches := program.FindAll("xxneedle needle")
	if len(matches) != 2 || matches[0].Start != 2 || matches[1].Start != 9 {
		t.Fatalf("global literal fast-path matches = %#v", matches)
	}

	pattern := []uint16{0xde00}
	subject := []uint16{0xd83d, 0xde00}
	legacy, err := CompileUTF16(pattern, "")
	if err != nil {
		t.Fatal(err)
	}
	if match, ok := legacy.FindUTF16(subject, 0); !ok || match.Start != 1 {
		t.Errorf("legacy literal did not match a UTF-16 code unit inside pair: %#v, %v", match, ok)
	}
	unicodeProgram, err := CompileUTF16(pattern, "u")
	if err != nil {
		t.Fatal(err)
	}
	if len(unicodeProgram.atom) != 0 {
		t.Fatal("Unicode literal incorrectly selected a code-unit-only atom path")
	}
	if _, ok := unicodeProgram.FindUTF16(subject, 0); ok {
		t.Error("Unicode literal matched a low surrogate inside a pair")
	}
}

func TestUnicodeCharacterClassMatchesLoneSurrogatesNotPairHalves(t *testing.T) {
	pattern := []uint16{'[', 0xd800, ']'}
	program, err := CompileUTF16(pattern, `u`)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := program.FindUTF16([]uint16{0xd800}, 0); !ok {
		t.Fatal("Unicode character class did not match a lone lead surrogate")
	}
	if _, ok := program.FindUTF16([]uint16{0xd800, 0xdc00}, 0); ok {
		t.Fatal("Unicode character class matched the lead half of a surrogate pair")
	}
}

func TestEscapedRegExpSourceAndToStringFormatting(t *testing.T) {
	for _, test := range []struct {
		source, flags, want string
	}{
		{`a/b`, `gi`, `/a\/b/gi`},
		{`[/]`, ``, `/[/]/`},
		{"a\nb", ``, `/a\nb/`},
		{"a\u2028b", ``, `/a\u2028b/`},
		{`a\/b`, ``, `/a\/b/`},
		{``, ``, `/(?:)/`},
	} {
		program, err := Compile(test.source, test.flags)
		if err != nil {
			t.Fatalf("Compile(%q, %q): %v", test.source, test.flags, err)
		}
		if got := string(utf16.Decode(program.RegExpStringUTF16())); got != test.want {
			t.Errorf("RegExp string for %q = %q; want %q", test.source, got, test.want)
		}
	}
	program, err := CompileUTF16([]uint16{0xd800}, "")
	if err != nil {
		t.Fatal(err)
	}
	if got := program.EscapedSourceUTF16(); len(got) != 1 || got[0] != 0xd800 {
		t.Errorf("escaped lone-surrogate source = %#v; want preserved code unit", got)
	}
}

func TestUnicodePropertyIgnoreCaseAndNegativeEscapeOrder(t *testing.T) {
	for _, test := range []struct {
		pattern, flags, input string
		want                  bool
	}{
		{`\p{Lu}`, "iu", "a", true},
		{`\p{Ll}`, "iu", "A", true},
		{`\P{Lowercase_Letter}`, "iu", "a", true},
		{`\P{Lowercase_Letter}`, "iv", "a", false},
	} {
		program, err := Compile(test.pattern, test.flags)
		if err != nil {
			t.Fatalf("Compile(%q, %q): %v", test.pattern, test.flags, err)
		}
		_, matched := program.Find(test.input, 0)
		if matched != test.want {
			t.Errorf("/%s/%s on %q matched %v; want %v", test.pattern, test.flags, test.input, matched, test.want)
		}
	}
}

func TestUnicodePropertyAliasesAreExact(t *testing.T) {
	for _, test := range []struct {
		pattern, input string
		want           bool
	}{
		{`\p{Letter}`, "α", true},
		{`\p{General_Category=Uppercase_Letter}`, "A", true},
		{`\p{gc=Lu}`, "A", true},
		{`\p{Script=Greek}`, "α", true},
		{`\p{sc=Grek}`, "α", true},
		{`\p{Script=Hiragana}`, "ー", false},
		{`\p{Script_Extensions=Hiragana}`, "ー", true},
		{`\p{scx=Kana}`, "ー", true},
		{`\p{scx=Greek}`, "α", true},
		{`\p{Emoji}`, "©", true},
		{`\p{Emoji_Presentation}`, "©", false},
		{`\p{EPres}`, "😀", true},
		{`\p{Emoji_Component}`, "️", true},
		{`\p{Emoji_Modifier}`, "🏽", true},
		{`\p{Emoji_Modifier_Base}`, "👍", true},
		{`\p{Extended_Pictographic}`, "©", true},
		{`\p{White_Space}`, " ", true},
		{`\p{White_Space}`, "\ufeff", false},
		{`\s`, "\ufeff", true},
		{`\p{Assigned}`, "\x00", true},
		{`\p{Assigned}`, "\U0010ffff", false},
		{`\p{General_Category=Unassigned}`, "\u0378", true},
		{`\p{Script=Unknown}`, "\u0378", true},
		{`\p{ID_Start}`, "_", false},
		{`\p{ID_Continue}`, "_", true},
	} {
		program, err := Compile(test.pattern, "u")
		if err != nil {
			t.Fatalf("Compile(%q): %v", test.pattern, err)
		}
		_, matched := program.Find(test.input, 0)
		if matched != test.want {
			t.Errorf("/%s/u on %q matched %v; want %v", test.pattern, test.input, matched, test.want)
		}
	}

	for _, pattern := range []string{
		`\p{letter}`,
		`\p{general_category=Lu}`,
		`\p{GeneralCategory=Lu}`,
		`\p{white_space}`,
		`\p{sc=grek}`,
		`\p{Script=Greek-Script}`,
	} {
		if _, err := Compile(pattern, "u"); err == nil {
			t.Errorf("Compile(%q) accepted a non-alias spelling", pattern)
		}
	}
}

func TestUnicodePropertiesOfStrings(t *testing.T) {
	family := "👨‍👩‍👧‍👦"
	for _, test := range []struct {
		pattern, input string
		want           bool
	}{
		{`^\p{Basic_Emoji}$`, "😀", true},
		{`^\p{Basic_Emoji}$`, "©️", true},
		{`^\p{Emoji_Keycap_Sequence}$`, "1️⃣", true},
		{`^\p{RGI_Emoji_Flag_Sequence}$`, "🇺🇸", true},
		{`^\p{RGI_Emoji_ZWJ_Sequence}$`, family, true},
		{`^\p{RGI_Emoji}$`, family, true},
		{`^\p{RGI_Emoji}$`, "x", false},
		{`\p{RGI_Emoji}X`, family + "X", true},
	} {
		program, err := Compile(test.pattern, "v")
		if err != nil {
			t.Fatalf("Compile(%q, v): %v", test.pattern, err)
		}
		_, matched := program.Find(test.input, 0)
		if matched != test.want {
			t.Errorf("/%s/v on %q matched %v; want %v", test.pattern, test.input, matched, test.want)
		}
	}

	for _, pattern := range []struct{ source, flags string }{
		{`\p{RGI_Emoji}`, "u"},
		{`\P{RGI_Emoji}`, "v"},
		{`[^\p{RGI_Emoji}]`, "v"},
	} {
		if _, err := Compile(pattern.source, pattern.flags); err == nil {
			t.Errorf("Compile(%q, %q) accepted an invalid string-property use", pattern.source, pattern.flags)
		}
	}
}

func TestV8BinaryUnicodePropertyAllowlistParses(t *testing.T) {
	properties := []string{
		"Alphabetic", "ASCII_Hex_Digit", "Bidi_Control", "Bidi_Mirrored",
		"Case_Ignorable", "Cased", "Changes_When_Casefolded", "Changes_When_Casemapped",
		"Changes_When_Lowercased", "Changes_When_NFKC_Casefolded", "Changes_When_Titlecased",
		"Changes_When_Uppercased", "Dash", "Default_Ignorable_Code_Point", "Deprecated",
		"Diacritic", "Emoji", "Emoji_Component", "Emoji_Modifier_Base", "Emoji_Modifier",
		"Emoji_Presentation", "Extended_Pictographic", "Extender", "Grapheme_Base",
		"Grapheme_Extend", "Hex_Digit", "ID_Continue", "ID_Start", "Ideographic",
		"IDS_Binary_Operator", "IDS_Trinary_Operator", "Join_Control",
		"Logical_Order_Exception", "Lowercase", "Math", "Noncharacter_Code_Point",
		"Pattern_Syntax", "Pattern_White_Space", "Quotation_Mark", "Radical",
		"Regional_Indicator", "Sentence_Terminal", "Soft_Dotted", "Terminal_Punctuation",
		"Unified_Ideograph", "Uppercase", "Variation_Selector", "White_Space",
		"XID_Continue", "XID_Start",
	}
	for _, property := range properties {
		if _, err := Compile(`\p{`+property+`}`, "u"); err != nil {
			t.Errorf("V8 binary property %q was rejected: %v", property, err)
		}
	}
	for _, property := range []string{`\p{Bidi_Class=R}`, `\p{Script}`, `\p{Greek}`} {
		if _, err := Compile(property, "u"); err == nil {
			t.Errorf("unsupported Unicode property %q was accepted", property)
		}
	}
}

func TestV8CaseFoldingModesAndUnicode18SupplementaryPair(t *testing.T) {
	for _, test := range []struct {
		pattern, flags, input string
		want                  bool
	}{
		{`s`, `i`, `ſ`, false},
		{`s`, `iu`, `ſ`, true},
		{`k`, `i`, `K`, false},
		{`k`, `iu`, `K`, true},
		{`ß`, `i`, `ẞ`, false},
		{`ß`, `iu`, `ẞ`, true},
		{`[s]`, `i`, `ſ`, false},
		{`[s]`, `iu`, `ſ`, true},
	} {
		program, err := Compile(test.pattern, test.flags)
		if err != nil {
			t.Fatalf("Compile(%q, %q): %v", test.pattern, test.flags, err)
		}
		_, matched := program.Find(test.input, 0)
		if matched != test.want {
			t.Errorf("/%s/%s on %q matched %v; want %v", test.pattern, test.flags, test.input, matched, test.want)
		}
	}

	// V8's special-case generator documents this Unicode 18 simple-fold
	// relationship. The Go standard library tables in this build are Unicode 17.
	longSharpS := []uint16{0xd837, 0xdf95} // U+1DF95
	for _, pattern := range []string{`ß`, `[ß]`, `\p{Ll}`} {
		program, err := Compile(pattern, `iu`)
		if err != nil {
			t.Fatalf("Compile(%q, iu): %v", pattern, err)
		}
		if _, matched := program.FindUTF16(longSharpS, 0); !matched {
			t.Errorf("/%s/iu did not match U+1DF95", pattern)
		}
	}
	for _, test := range []struct {
		pattern, flags string
		want           bool
	}{
		{`[^\u00df]`, `iu`, false},
		{`[^\u00df]`, `iv`, false},
		{`\P{Ll}`, `iu`, true},
		{`\P{Ll}`, `iv`, false},
	} {
		program, err := Compile(test.pattern, test.flags)
		if err != nil {
			t.Fatalf("Compile(%q, %s): %v", test.pattern, test.flags, err)
		}
		_, matched := program.FindUTF16(longSharpS, 0)
		if matched != test.want {
			t.Errorf("/%s/%s on U+1DF95 matched %v; want %v", test.pattern, test.flags, matched, test.want)
		}
	}
}

func TestCapturesAndGlobal(t *testing.T) {
	program, err := Compile(`(?<word>a+)(b)`, "g")
	if err != nil {
		t.Fatal(err)
	}
	matches := program.FindAll("zaaabb")
	if len(matches) != 1 || matches[0].Start != 1 || matches[0].End != 5 {
		t.Fatalf("unexpected matches: %#v", matches)
	}
	if got := matches[0].Captures[1]; got.Start != 1 || got.End != 4 {
		t.Fatalf("named capture span = %#v", got)
	}
}

func TestCaptureUTF16DistinguishesUnmatchedAndEmpty(t *testing.T) {
	program, err := Compile(`(?:(a)|())`, ``)
	if err != nil {
		t.Fatal(err)
	}
	subject := utf16.Encode([]rune(``))
	match, ok := program.FindUTF16(subject, 0)
	if !ok {
		t.Fatal("empty alternative did not match")
	}
	if _, matched := match.CaptureUTF16(subject, 1); matched {
		t.Fatal("unmatched capture was reported as participating")
	}
	empty, matched := match.CaptureUTF16(subject, 2)
	if !matched || len(empty) != 0 {
		t.Fatalf("participating empty capture = %#v, %v; want empty slice, true", empty, matched)
	}
	if _, matched := match.CaptureUTF16(subject, -1); matched {
		t.Fatal("negative capture index was accepted")
	}
	if _, matched := match.CaptureUTF16(subject, len(match.Captures)); matched {
		t.Fatal("out-of-range capture index was accepted")
	}
}

func TestBacktrackingInterpreterPriorityAndCaptureRestoration(t *testing.T) {
	for _, test := range []struct {
		pattern string
		input   string
		start   int
		end     int
		group1  Span
	}{
		{`(a|ab)c`, `abc`, 0, 3, Span{0, 2}},
		{`(a)?a`, `a`, 0, 1, Span{-1, -1}},
		{`(a+?)`, `aaa`, 0, 1, Span{0, 1}},
		{`(a+)`, `aaa`, 0, 3, Span{0, 3}},
	} {
		program, err := Compile(test.pattern, "")
		if err != nil {
			t.Fatalf("Compile(%q): %v", test.pattern, err)
		}
		match, ok := program.Find(test.input, 0)
		if !ok || match.Start != test.start || match.End != test.end {
			t.Errorf("/%s/ on %q = %#v, %v; want [%d,%d)", test.pattern, test.input, match, ok, test.start, test.end)
			continue
		}
		if got := match.Captures[1]; got != test.group1 {
			t.Errorf("/%s/ capture 1 = %#v; want %#v", test.pattern, got, test.group1)
		}
	}
}

func TestBacktrackingInterpreterAssertionsAndBackreferences(t *testing.T) {
	for _, test := range []struct {
		pattern string
		input   string
		want    bool
	}{
		{`(?<=([ab]+)([bc]+))$`, `abc`, true},
		{`(?=(a|ab))\1c`, `abc`, false}, // Successful lookahead is atomic.
		{`(a)?\1b`, `b`, true},          // An unmatched backreference is empty.
	} {
		program, err := Compile(test.pattern, "")
		if err != nil {
			t.Fatalf("Compile(%q): %v", test.pattern, err)
		}
		match, ok := program.Find(test.input, 0)
		if ok != test.want {
			t.Errorf("/%s/ on %q matched %v; want %v (%#v)", test.pattern, test.input, ok, test.want, match)
			continue
		}
		if test.pattern == `(?<=([ab]+)([bc]+))$` && (match.Captures[1] != (Span{0, 1}) || match.Captures[2] != (Span{1, 3})) {
			t.Errorf("lookbehind captures = %#v; want [0,1) and [1,3)", match.Captures)
		}
	}
	program, err := Compile(`(\u{1DF95})\1`, `iu`)
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		input string
		want  bool
	}{
		{input: string([]rune{0x1df95, 0x1df95}), want: true},
		{input: string([]rune{0x1df95, 0x00df}), want: false},
	} {
		if _, matched := program.Find(test.input, 0); matched != test.want {
			t.Errorf("Unicode-folded backreference on %U matched %v; want %v", []rune(test.input), matched, test.want)
		}
	}
}

func TestNullableRepeatBacktrackingMatchesV8(t *testing.T) {
	for _, test := range []struct {
		pattern string
		input   string
		end     int
		groups  []Span
	}{
		{`(a*)*`, `aaa`, 3, []Span{{0, 3}}},
		{`(a*)*`, ``, 0, []Span{{0, 0}}},
		{`((?:)|a)*`, `aaa`, 3, []Span{{2, 3}}},
		{`(a?)*`, `aaa`, 3, []Span{{2, 3}}},
		{`((a)?)*`, `aaa`, 3, []Span{{2, 3}, {2, 3}}},
	} {
		program, err := Compile(test.pattern, "")
		if err != nil {
			t.Fatalf("Compile(%q): %v", test.pattern, err)
		}
		match, ok := program.Find(test.input, 0)
		if !ok || match.Start != 0 || match.End != test.end {
			t.Errorf("/%s/ on %q = %#v, %v; want [0,%d)", test.pattern, test.input, match, ok, test.end)
			continue
		}
		for i, want := range test.groups {
			if got := match.Captures[i+1]; got != want {
				t.Errorf("/%s/ capture %d = %#v; want %#v", test.pattern, i+1, got, want)
			}
		}
	}
}

func TestClassStringAlternativesRestoreLoopState(t *testing.T) {
	program, err := Compile(`([\q{a|ab}])*b`, `v`)
	if err != nil {
		t.Fatal(err)
	}
	match, ok := program.Find(`ab`, 0)
	if !ok || match.Start != 0 || match.End != 2 {
		t.Fatalf("string-set alternative backtrack = %#v, %v; want [0,2)", match, ok)
	}
}

func TestCharacterClassRangesAreSortedAndMerged(t *testing.T) {
	program, err := Compile(`[zafb-c]`, ``)
	if err != nil {
		t.Fatal(err)
	}
	if len(program.code) < 1 || len(program.code[0].class) != 1 {
		t.Fatalf("class program = %#v", program.code)
	}
	ranges := program.code[0].class[0].ranges
	if len(ranges) != 3 || ranges[0] != (classRange{'a', 'c'}) || ranges[1] != (classRange{'f', 'f'}) || ranges[2] != (classRange{'z', 'z'}) {
		t.Fatalf("normalized class ranges = %#v", ranges)
	}
	for _, input := range []string{"a", "b", "c", "f", "z"} {
		if _, ok := program.Find(input, 0); !ok {
			t.Errorf("class did not match %q", input)
		}
	}
	if _, ok := program.Find("d", 0); ok {
		t.Error("class unexpectedly matched d")
	}
}

func TestGlobalEmptyMatchAdvancesByUnicodeCodePoint(t *testing.T) {
	program, err := Compile(`(?:)`, "gu")
	if err != nil {
		t.Fatal(err)
	}
	matches := program.FindAll("😀")
	if len(matches) != 2 || matches[0].Start != 0 || matches[0].End != 0 || matches[1].Start != 2 || matches[1].End != 2 {
		t.Fatalf("zero-width global matches = %#v; want empty matches at UTF-16 offsets 0 and 2", matches)
	}
}

func TestMultilineAnchorsAndStrictEndAnchor(t *testing.T) {
	for _, flags := range []string{`gm`, `gml`} {
		for _, pattern := range []string{`^`, `$`} {
			program, err := Compile(pattern, flags)
			if err != nil {
				t.Fatal(err)
			}
			matches := program.FindAll("\r\n")
			if len(matches) != 3 || matches[0].Start != 0 || matches[1].Start != 1 || matches[2].Start != 2 {
				t.Errorf("/%s/%s matches in CRLF = %#v; want zero-width matches at UTF-16 offsets 0, 1, 2", pattern, flags, matches)
			}
		}
	}
	for _, flags := range []string{``, `l`} {
		if program, err := Compile(`a$`, flags); err != nil {
			t.Fatal(err)
		} else if _, matched := program.Find("a\n", 0); matched {
			t.Errorf("/%s/ matched before a final line terminator", flags)
		}
	}
	for _, flags := range []string{`m`, `ml`} {
		if program, err := Compile(`a$`, flags); err != nil {
			t.Fatal(err)
		} else if match, matched := program.Find("a\n", 0); !matched || match.End != 1 {
			t.Errorf("/%s/ end anchor before line terminator = %#v, %v; want [0,1)", flags, match, matched)
		}
	}
}

func TestUTF16RegExpUtilities(t *testing.T) {
	input := []uint16{'a', 0xd83d, 0xde00, 'b'}
	if got := AdvanceStringIndexUTF16(input, 1, true); got != 3 {
		t.Errorf("Unicode AdvanceStringIndex at lead surrogate = %d; want 3", got)
	}
	if got := AdvanceStringIndexUTF16(input, 1, false); got != 2 {
		t.Errorf("non-Unicode AdvanceStringIndex = %d; want 2", got)
	}
	if got := AdvanceStringIndexUTF16(input, len(input), true); got != len(input)+1 {
		t.Errorf("AdvanceStringIndex at end = %d; want %d", got, len(input)+1)
	}

	match := Match{Captures: []Span{{0, 4}, {1, 3}, {2, 2}, {-1, -1}}}
	if !IsMatchedCapture(match, 2) || IsMatchedCapture(match, 3) || IsMatchedCapture(match, -1) || IsMatchedCapture(match, 4) {
		t.Fatal("capture presence test did not distinguish matched, unmatched, and invalid groups")
	}
	capture, ok := CaptureUTF16(input, match, 1)
	if !ok || len(capture) != 2 || capture[0] != 0xd83d || capture[1] != 0xde00 {
		t.Fatalf("UTF-16 capture = %#v, %v", capture, ok)
	}
	capture[0] = 'x'
	if input[1] != 0xd83d {
		t.Fatal("CaptureUTF16 returned an alias to the input buffer")
	}
	if capture, ok := CaptureUTF16(input, match, 2); !ok || len(capture) != 0 {
		t.Fatalf("empty matched capture = %#v, %v; want empty and matched", capture, ok)
	}
	if capture, ok := CaptureUTF16(input, match, 3); ok || capture != nil {
		t.Fatalf("unmatched capture = %#v, %v; want nil and false", capture, ok)
	}
}

func TestRepeatedCapturesClearWhenFinalIterationSkipsGroup(t *testing.T) {
	for _, pattern := range []string{`((a)?b)+`, `((a)?b){2}`, `((a)?b){1,2}`} {
		program, err := Compile(pattern, "")
		if err != nil {
			t.Fatalf("Compile(%q): %v", pattern, err)
		}
		match, ok := program.Find("abb", 0)
		if !ok || match.End != 3 {
			t.Fatalf("/%s/ on abb = %#v, %v", pattern, match, ok)
		}
		if got := match.Captures[1]; got != (Span{Start: 2, End: 3}) {
			t.Errorf("/%s/ outer capture = %#v; want [2,3)", pattern, got)
		}
		if got := match.Captures[2]; got != (Span{Start: -1, End: -1}) {
			t.Errorf("/%s/ skipped final-iteration capture = %#v; want unmatched", pattern, got)
		}
	}
}

func TestUnsupportedSyntaxIsExplicit(t *testing.T) {
	for _, pattern := range []string{`(?<name>a)\k<missing>`, `(?<name>a)(?<name>b)`} {
		if _, err := Compile(pattern, ""); err == nil {
			t.Fatalf("Compile(%q) unexpectedly succeeded", pattern)
		}
	}
}

func TestScopedRegExpModifiers(t *testing.T) {
	tests := []struct {
		pattern, flags, input string
		want                  bool
	}{
		{`(?i:ab)c`, "", "ABc", true},
		{`(?i:ab)c`, "", "ABC", false},
		{`(?-i:a)b`, "i", "aB", true},
		{`(?-i:a)b`, "i", "Ab", false},
		{`(?i:(?-i:a)A)`, "", "aA", true},
		{`(?i:(?-i:a)A)`, "", "AA", false},
		{`(?i:(a)\1)`, "", "aA", true},
		{`(?i:[a-z])`, "", "G", true},
		{`(?-i:[a-z])`, "i", "G", false},
		{`(?s:.)`, "", "\n", true},
		{`(?-s:.)`, "s", "\n", false},
		{`(?m:^b)`, "", "x\nb", true},
		{`(?-m:^b)`, "m", "x\nb", false},
		{`(?i-ms:a)`, "ms", "A", true},
	}
	for _, test := range tests {
		program, err := Compile(test.pattern, test.flags)
		if err != nil {
			t.Fatalf("Compile(%q, %q): %v", test.pattern, test.flags, err)
		}
		_, matched := program.Find(test.input, 0)
		if matched != test.want {
			t.Errorf("/%s/%s on %q matched %v; want %v", test.pattern, test.flags, test.input, matched, test.want)
		}
	}
}

func TestInvalidScopedRegExpModifiers(t *testing.T) {
	for _, pattern := range []string{`(?i)`, `(?ii:a)`, `(?i-i:a)`, `(?i--m:a)`, `(?-:a)`, `(?i-d:a)`} {
		if _, err := Compile(pattern, ""); err == nil {
			t.Errorf("Compile(%q) unexpectedly succeeded", pattern)
		}
	}
}

func TestLegacyCharacterClassTrailingHyphen(t *testing.T) {
	program, err := Compile(`[a-]`, "")
	if err != nil {
		t.Fatal(err)
	}
	for _, input := range []string{"a", "-"} {
		if _, matched := program.Find(input, 0); !matched {
			t.Errorf("/[a-]/ unexpectedly failed to match %q", input)
		}
	}
	if _, matched := program.Find("]", 0); matched {
		t.Fatal("/[a-]/ incorrectly matched a closing bracket")
	}
}

func TestNamedCaptureAndBackreferenceUnicodeEscapes(t *testing.T) {
	for _, test := range []struct {
		pattern, input string
	}{
		{`(?<\u0061>x)\k<a>`, "xx"},
		{`(?<a>x)\k<\u0061>`, "xx"},
		{`(?<\u{10400}>x)\k<𐐀>`, "xx"},
		{`(?<\uD835\uDC00>x)\k<𝐀>`, "xx"},
	} {
		program, err := Compile(test.pattern, "")
		if err != nil {
			t.Fatalf("Compile(%q): %v", test.pattern, err)
		}
		if _, matched := program.Find(test.input, 0); !matched {
			t.Errorf("/%s/ did not match %q", test.pattern, test.input)
		}
	}
}

func TestDuplicateNamedCapturesInDisjointAlternatives(t *testing.T) {
	program, err := Compile(`(?:(?<x>a)|(?<x>b))\k<x>`, ``)
	if err != nil {
		t.Fatalf("Compile duplicate name in disjoint alternatives: %v", err)
	}
	if got, want := program.GroupCaptureIndices["x"], []int{1, 2}; !reflect.DeepEqual(got, want) {
		t.Fatalf("capture indices for x = %v; want %v", got, want)
	}
	if got := program.GroupIndices["x"]; got != 1 {
		t.Fatalf("representative capture index for x = %d; want first capture 1", got)
	}
	for _, test := range []struct {
		input  string
		first  Span
		second Span
	}{
		{input: `aa`, first: Span{0, 1}, second: Span{-1, -1}},
		{input: `bb`, first: Span{-1, -1}, second: Span{0, 1}},
	} {
		match, ok := program.Find(test.input, 0)
		if !ok || match.End != 2 {
			t.Errorf("duplicate-name backreference on %q = %#v, %v; want match [0,2)", test.input, match, ok)
			continue
		}
		if match.Captures[1] != test.first || match.Captures[2] != test.second {
			t.Errorf("duplicate-name captures on %q = %#v; want group spans %#v and %#v", test.input, match.Captures, test.first, test.second)
		}
	}

	// A named backreference is empty when none of the duplicate alternatives
	// containing that name participated.
	empty, err := Compile(`(?:(?<x>a)|b)\k<x>`, ``)
	if err != nil {
		t.Fatal(err)
	}
	if match, ok := empty.Find(`b`, 0); !ok || match.End != 1 {
		t.Errorf("unmatched duplicate-name backreference on b = %#v, %v; want [0,1)", match, ok)
	}
}

func TestDuplicateNamedCapturesMustBeMutuallyExclusive(t *testing.T) {
	for _, pattern := range []string{
		`(?<x>a)(?<x>b)`,
		`(?:(?<x>a)|b)(?<x>c)`,
		`(?:(?<x>a)(?<x>b)|c)`,
	} {
		if _, err := Compile(pattern, ``); err == nil {
			t.Errorf("Compile(%q) accepted duplicate named captures that can coexist", pattern)
		}
	}
}

func TestInvalidEscapedCaptureNames(t *testing.T) {
	for _, pattern := range []string{`(?<\u0030x>a)`, `(?<\u{110000}>a)`, `(?<\u{D800}>a)`} {
		if _, err := Compile(pattern, ""); err == nil {
			t.Errorf("Compile(%q) unexpectedly succeeded", pattern)
		}
	}
}

func TestLegacyNamedBackreferenceLookalikeWithoutNamedCaptures(t *testing.T) {
	for _, test := range []struct {
		pattern, input string
	}{
		{`\k<foo`, `k<foo`},
		{`\k<\u0061>`, `k<a>`},
	} {
		program, err := Compile(test.pattern, "")
		if err != nil {
			t.Fatalf("Compile(%q): %v", test.pattern, err)
		}
		if _, matched := program.Find(test.input, 0); !matched {
			t.Errorf("legacy /%s/ did not match %q", test.pattern, test.input)
		}
	}
	if _, err := Compile(`(?<x>a)\k<foo`, ""); err == nil {
		t.Fatal("unterminated named backreference accepted when a named capture exists")
	}
	if _, err := Compile(`\k<foo`, "u"); err == nil {
		t.Fatal("invalid named backreference accepted in Unicode mode")
	}
	program, err := Compile(`(?:[]|(?<x>a)\k<x>)`, "")
	if err != nil {
		t.Fatal(err)
	}
	if _, matched := program.Find("aa", 0); !matched {
		t.Fatal("named capture after an empty class was missed by the pre-scan")
	}
}

func TestParserNestingLimitIsReportedAsSyntaxError(t *testing.T) {
	groups := strings.Repeat("(", maxParserNesting+1) + "a" + strings.Repeat(")", maxParserNesting+1)
	if _, err := Compile(groups, ""); err == nil {
		t.Fatal("excessive group nesting unexpectedly compiled")
	}
	classes := strings.Repeat("[", maxParserNesting+1) + "a" + strings.Repeat("]", maxParserNesting+1)
	if _, err := Compile(classes, "v"); err == nil {
		t.Fatal("excessive nested class depth unexpectedly compiled")
	}

	captures := strings.Repeat("()", maxCaptureGroups+1)
	if _, err := Compile(captures, ""); err == nil || !strings.Contains(err.Error(), errorString(errorTooManyCaptures)) {
		t.Fatalf("excessive capture count error = %v; want V8's too-many-captures error", err)
	}
}

func TestFlagValidation(t *testing.T) {
	if _, err := ParseFlags("gg"); err == nil {
		t.Fatal("duplicate flags accepted")
	}
	if _, err := ParseFlags("q"); err == nil {
		t.Fatal("unknown flag accepted")
	}
	flags, err := ParseFlags("gimsuy")
	if err != nil || !flags.Global || !flags.IgnoreCase || !flags.Multiline || !flags.DotAll || !flags.Unicode || !flags.Sticky {
		t.Fatalf("flags = %#v, %v", flags, err)
	}
	if _, err := ParseFlags("uv"); err == nil {
		t.Fatal("mutually exclusive Unicode flags accepted")
	}
	linear, err := ParseFlags("l")
	if err != nil || !IsLinear(linear) {
		t.Fatalf("linear flag = %#v, %v", linear, err)
	}
}

func TestV8RegExpErrorCatalog(t *testing.T) {
	if got, want := len(regexpErrorMessages), int(errorNumErrors); got != want {
		t.Fatalf("error catalog has %d entries; want %d", got, want)
	}
	for code := regexpError(0); code < errorNumErrors; code++ {
		_ = errorString(code)
	}
	if got := errorString(errorNothingToRepeat); got != "Nothing to repeat" {
		t.Errorf("nothing-to-repeat error text = %q", got)
	}
	if got := errorString(errorEscapeAtEndOfPattern); got != `\ at end of pattern` {
		t.Errorf("trailing-escape error text = %q", got)
	}
	for _, code := range []regexpError{errorStackOverflow, errorAnalysisStackOverflow, errorTooLarge} {
		if !errorIsStackOverflow(code) {
			t.Errorf("error %d was not classified as a stack/resource overflow", code)
		}
	}
	if errorIsStackOverflow(errorInvalidEscape) {
		t.Fatal("invalid escape was classified as a stack/resource overflow")
	}
}

func TestV8FlagTableAndUnicodePredicates(t *testing.T) {
	flagChars := []rune("dgilmsuvy")
	for _, char := range flagChars {
		flag, ok := TryFlagFromChar(char)
		if !ok {
			t.Fatalf("TryFlagFromChar(%q) failed", char)
		}
		var flags Flags
		flags.Set(flag, true)
		if got := flags.String(); got != string(char) {
			t.Errorf("flag %q formats as %q", char, got)
		}
		if !flags.Has(flag) {
			t.Errorf("flag %q was not set", char)
		}
	}
	if _, ok := TryFlagFromChar('x'); ok {
		t.Fatal("unknown flag character was accepted")
	}

	unicodeSets, err := ParseFlags("v")
	if err != nil || !IsUnicodeSets(unicodeSets) || IsUnicode(unicodeSets) || !IsEitherUnicode(unicodeSets) || unicodeSets.String() != "v" {
		t.Fatalf("v-mode flag predicates/format = %#v, %v", unicodeSets, err)
	}
	for _, flags := range []string{"uy", "ug", "vy", "vg"} {
		parsed, err := ParseFlags(flags)
		if err != nil || !ShouldOptionallyStepBackToLeadSurrogate(parsed) {
			t.Errorf("surrogate rewind predicate for %q = %v, %v", flags, ShouldOptionallyStepBackToLeadSurrogate(parsed), err)
		}
	}
	for _, flags := range []string{"u", "g", "y", ""} {
		parsed, err := ParseFlags(flags)
		if err != nil || ShouldOptionallyStepBackToLeadSurrogate(parsed) {
			t.Errorf("unexpected surrogate rewind predicate for %q: %v, %v", flags, ShouldOptionallyStepBackToLeadSurrogate(parsed), err)
		}
	}
}
