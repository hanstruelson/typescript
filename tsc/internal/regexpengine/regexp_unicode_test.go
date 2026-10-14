package regexpengine

import "testing"

func TestUnicodePropertyAliasesAreExact(t *testing.T) {
	for _, pattern := range []string{
		`\p{Letter}`,
		`\p{General_Category=Uppercase_Letter}`,
		`\p{gc=Lu}`,
		`\p{Script=Greek}`,
		`\p{sc=Grek}`,
		`\p{Script_Extensions=Hiragana}`,
		`\p{scx=Kana}`,
		`\p{Emoji}`,
		`\p{White_Space}`,
		`\p{Assigned}`,
		`\p{ID_Start}`,
	} {
		if err := VerifySyntax(pattern, "u"); err != nil {
			t.Errorf("VerifySyntax(%q, u): %v", pattern, err)
		}
	}
	for _, pattern := range []string{
		`\p{letter}`,
		`\p{general_category=Lu}`,
		`\p{GeneralCategory=Lu}`,
		`\p{white_space}`,
		`\p{sc=grek}`,
		`\p{Script=Greek-Script}`,
		`\p{Bidi_Class=R}`,
	} {
		if err := VerifySyntax(pattern, "u"); err == nil {
			t.Errorf("VerifySyntax(%q, u) accepted a non-alias or unsupported property", pattern)
		}
	}
}

func TestUnicodePropertyTablesPreserveRepresentativeMembership(t *testing.T) {
	for _, test := range []struct {
		name string
		char rune
		want bool
	}{
		{"Letter", 'α', true},
		{"General_Category=Uppercase_Letter", 'A', true},
		{"Script=Greek", 'α', true},
		{"Script=Hiragana", 'ー', false},
		{"Script_Extensions=Hiragana", 'ー', true},
		{"Emoji", '©', true},
		{"Emoji_Presentation", '©', false},
		{"White_Space", ' ', true},
		{"White_Space", '\ufeff', false},
		{"Assigned", 0x10ffff, false},
		{"General_Category=Unassigned", 0x0378, true},
	} {
		property, ok := canonicalUnicodeProperty(test.name)
		if !ok {
			t.Errorf("canonicalUnicodeProperty(%q) failed", test.name)
			continue
		}
		if got := unicodePropertyContains(test.char, property); got != test.want {
			t.Errorf("property %q contains U+%04X = %v; want %v", test.name, test.char, got, test.want)
		}
	}
}

func TestUnicodePropertiesOfStringsAreParsedOnlyInVMode(t *testing.T) {
	for _, pattern := range []string{
		`\p{Basic_Emoji}`,
		`\p{Emoji_Keycap_Sequence}`,
		`\p{RGI_Emoji_Flag_Sequence}`,
		`\p{RGI_Emoji_ZWJ_Sequence}`,
		`\p{RGI_Emoji}`,
	} {
		if err := VerifySyntax(pattern, "v"); err != nil {
			t.Errorf("VerifySyntax(%q, v): %v", pattern, err)
		}
	}
	for _, test := range []struct{ pattern, flags string }{
		{`\p{RGI_Emoji}`, "u"},
		{`\P{RGI_Emoji}`, "v"},
		{`[^\p{RGI_Emoji}]`, "v"},
	} {
		if err := VerifySyntax(test.pattern, test.flags); err == nil {
			t.Errorf("VerifySyntax(%q, %q) accepted an invalid string-property use", test.pattern, test.flags)
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
		if err := VerifySyntax(`\p{`+property+`}`, "u"); err != nil {
			t.Errorf("V8 binary property %q was rejected: %v", property, err)
		}
	}
}

func TestV8CaseFoldingTablesKeepLegacyAndUnicodeModesDistinct(t *testing.T) {
	for _, test := range []struct {
		a, b        rune
		unicodeMode bool
		want        bool
	}{
		{'s', 'ſ', false, false},
		{'s', 'ſ', true, true},
		{'k', 'K', false, false},
		{'k', 'K', true, true},
		{'α', 'Α', false, true},
		{'i', 'ı', true, false},
		{'ß', 'ẞ', false, false},
		{'ß', 'ẞ', true, true},
		{'ß', 0x1df95, true, true},
	} {
		got := test.a == test.b || (test.unicodeMode && forEachUnicodeCaseEquivalent(test.a, func(candidate rune) bool { return candidate == test.b })) || (!test.unicodeMode && canonicalizeLegacy(test.a) == canonicalizeLegacy(test.b))
		if got != test.want {
			t.Errorf("case equivalence U+%04X / U+%04X unicode=%v = %v; want %v", test.a, test.b, test.unicodeMode, got, test.want)
		}
	}
}

func TestCharacterClassRangeCanonicalization(t *testing.T) {
	got := v8ASTCanonicalizeRanges([]v8ASTCharacterRange{
		{from: 'z', to: 'z'},
		{from: 'a', to: 'b'},
		{from: 'b', to: 'c'},
		{from: 'f', to: 'f'},
	})
	want := []v8ASTCharacterRange{{from: 'a', to: 'c'}, {from: 'f', to: 'f'}, {from: 'z', to: 'z'}}
	if len(got) != len(want) {
		t.Fatalf("canonical ranges = %#v; want %#v", got, want)
	}
	for index := range want {
		if got[index] != want[index] {
			t.Fatalf("canonical ranges = %#v; want %#v", got, want)
		}
	}
}
