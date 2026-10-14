package regexpengine

import (
	"testing"
)

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
