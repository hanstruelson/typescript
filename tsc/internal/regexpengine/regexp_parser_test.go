package regexpengine

import (
	"reflect"
	"strings"
	"testing"
	"unicode/utf16"
)

func TestVerifySyntaxParsesV8FrontendWithoutExecutingPatterns(t *testing.T) {
	for _, test := range []struct {
		pattern, flags string
		wantError      bool
	}{
		{pattern: `(?<name>a)\k<name>`, flags: `u`},
		{pattern: `(?<=a)b`, flags: `u`},
		{pattern: `(?=a)a`, flags: `l`},
		{pattern: `(?<name>a)\k<missing>`, flags: `u`, wantError: true},
		{pattern: `[`, wantError: true},
		{pattern: `(?<name>a)(?<name>b)`, wantError: true},
	} {
		err := VerifySyntax(test.pattern, test.flags)
		if (err != nil) != test.wantError {
			t.Errorf("VerifySyntax(%q, %q) error = %v, wantError %v", test.pattern, test.flags, err, test.wantError)
		}
	}
}

func TestCaptureNamesUsePinnedUnicodeIdentifierTables(t *testing.T) {
	for _, name := range []string{
		"$value",
		"_value",
		"℘",        // Other_ID_Start.
		"a\u200cb", // ECMAScript adds ZWNJ to IdentifierPart.
		"a\u200db", // ECMAScript adds ZWJ to IdentifierPart.
		"𮯰",        // CJK Extension I, present in the pinned Unicode data.
	} {
		if err := VerifySyntax("(?<"+name+">a)", "u"); err != nil {
			t.Errorf("capture name %q rejected: %v", name, err)
		}
	}
	for _, name := range []string{"\u0301a", "😀", "a-b"} {
		if err := VerifySyntax("(?<"+name+">a)", "u"); err == nil {
			t.Errorf("invalid capture name %q accepted", name)
		}
	}
}

func TestCaptureNamesExposeDuplicateNamesInDisjointAlternatives(t *testing.T) {
	names, err := CaptureNamesUTF16(utf16.Encode([]rune(`(?:(?<x>a)|(?<x>b))\k<x>`)), "")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := names["x"], []int{1, 2}; !reflect.DeepEqual(got, want) {
		t.Fatalf("capture indices for x = %v; want %v", got, want)
	}
	if len(names) != 1 {
		t.Fatalf("named capture table = %#v; want only x", names)
	}
}

func TestScopedRegExpModifierSyntax(t *testing.T) {
	for _, test := range []struct{ pattern, flags string }{
		{`(?i:ab)c`, ""},
		{`(?-i:a)b`, "i"},
		{`(?i:(?-i:a)A)`, ""},
		{`(?i:[a-z])`, ""},
		{`(?-i:[a-z])`, "i"},
		{`(?s:.)`, ""},
		{`(?m:^b)`, ""},
		{`(?i-ms:a)`, "ms"},
	} {
		if err := VerifySyntax(test.pattern, test.flags); err != nil {
			t.Errorf("VerifySyntax(%q, %q): %v", test.pattern, test.flags, err)
		}
	}
	for _, pattern := range []string{`(?i)`, `(?ii:a)`, `(?i-i:a)`, `(?i--m:a)`, `(?-:a)`, `(?i-d:a)`} {
		if err := VerifySyntax(pattern, ""); err == nil {
			t.Errorf("VerifySyntax(%q) unexpectedly succeeded", pattern)
		}
	}
}

func TestAnnexBSyntaxAndUnicodeEarlyErrors(t *testing.T) {
	for _, pattern := range []string{`[a-]`, `\x1G`, `\xGG`, `\u00G0`, `\uG000`, `\c1`, `\k<foo`} {
		if err := VerifySyntax(pattern, ""); err != nil {
			t.Errorf("legacy syntax %q rejected: %v", pattern, err)
		}
	}
	for _, pattern := range []string{`\x1G`, `\xGG`, `\u00G0`, `\uG000`, `\c1`, `(?<x>a)\k<foo`, `\k<foo`} {
		if err := VerifySyntax(pattern, "u"); err == nil {
			t.Errorf("Unicode syntax %q unexpectedly succeeded", pattern)
		}
	}
	for _, pattern := range []string{`(?<\u0030x>a)`, `(?<\u{110000}>a)`, `(?<\u{D800}>a)`} {
		if err := VerifySyntax(pattern, ""); err == nil {
			t.Errorf("invalid escaped capture name %q accepted", pattern)
		}
	}
}

func TestDuplicateNamedCapturesMustBeMutuallyExclusive(t *testing.T) {
	for _, pattern := range []string{
		`(?<x>a)(?<x>b)`,
		`(?:(?<x>a)|b)(?<x>c)`,
		`(?:(?<x>a)(?<x>b)|c)`,
	} {
		if err := VerifySyntax(pattern, ""); err == nil {
			t.Errorf("VerifySyntax(%q) accepted duplicate named captures that can coexist", pattern)
		}
	}
}

func TestParserNestingAndCaptureLimitsRemainSyntaxErrors(t *testing.T) {
	groups := strings.Repeat("(", maxParserNesting+1) + "a" + strings.Repeat(")", maxParserNesting+1)
	if err := VerifySyntax(groups, ""); err == nil {
		t.Fatal("excessive group nesting unexpectedly parsed")
	}
	classes := strings.Repeat("[", maxParserNesting+1) + "a" + strings.Repeat("]", maxParserNesting+1)
	if err := VerifySyntax(classes, "v"); err == nil {
		t.Fatal("excessive nested class depth unexpectedly parsed")
	}
	captures := strings.Repeat("()", maxCaptureGroups+1)
	if err := VerifySyntax(captures, ""); err == nil || !strings.Contains(err.Error(), errorString(errorTooManyCaptures)) {
		t.Fatalf("excessive capture count error = %v; want V8's too-many-captures error", err)
	}
}
