package regexpengine

import (
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

func TestCompileToGoMatcherLiteral(t *testing.T) {
	for _, test := range []struct {
		name, pattern, flags string
		captureSlots         int
	}{
		{name: "atom", pattern: "users", captureSlots: 2},
		{name: "graph", pattern: `^/users/[0-9]+/?$`, flags: "i", captureSlots: 2},
		{name: "capture", pattern: `(?<id>[0-9]+)`, captureSlots: 4},
	} {
		t.Run(test.name, func(t *testing.T) {
			expression, captureSlots, err := CompileToGoMatcherLiteral(test.pattern, test.flags)
			if err != nil {
				t.Fatalf("CompileToGoMatcherLiteral: %v", err)
			}
			if captureSlots != test.captureSlots {
				t.Fatalf("capture slots = %d, want %d", captureSlots, test.captureSlots)
			}
			if _, err := parser.ParseExprFrom(token.NewFileSet(), "matcher.go", expression, parser.AllErrors); err != nil {
				t.Fatalf("generated matcher is not a Go expression: %v\n%s", err, expression)
			}
			if !strings.Contains(expression, "return func(") {
				t.Fatalf("generated matcher does not return a function literal:\n%s", expression)
			}
			if strings.Contains(expression, "func tsJITDynamicRegexInto") {
				t.Fatalf("generated matcher contains a package-level function declaration:\n%s", expression)
			}
		})
	}
}
