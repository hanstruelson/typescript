package goemit

import (
	"strings"
	"testing"
)

func TestCompiledRegExpEmbeddingDeduplicatesMatchersAndKeepsNames(t *testing.T) {
	e := &emitter{}
	first := e.compiledRegExp("(?<word>a)", "d", nil)
	second := e.compiledRegExp("(?<word>a)", "d", nil)
	if first != second {
		t.Fatalf("same pattern generated different constructors:\n%s\n%s", first, second)
	}
	matcher := e.regexMatchers["(?<word>a)\x00d"]
	if matcher == "" {
		t.Fatal("matcher name was not cached")
	}
	for _, declaration := range []string{"func " + matcher + "(", "func " + matcher + "Into("} {
		if got := strings.Count(e.regexText.String(), declaration); got != 1 {
			t.Fatalf("generated %d declarations for %q", got, declaration)
		}
	}
	if !strings.Contains(first, `name:"word",indices:[]int{1}`) {
		t.Fatalf("named capture metadata missing from constructor: %s", first)
	}
}
