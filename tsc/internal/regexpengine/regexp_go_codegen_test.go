package regexpengine

import (
	"fmt"
	"go/ast"
	"go/importer"
	goparser "go/parser"
	"go/token"
	"go/types"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func TestCompileToGoLargeClassUsesUTF16Bitmap(t *testing.T) {
	var pattern strings.Builder
	pattern.WriteByte('[')
	for unit := 0; unit < 140; unit += 2 {
		fmt.Fprintf(&pattern, `\u%04x`, unit)
	}
	pattern.WriteByte(']')
	source, err := CompileToGo(pattern.String(), "", "MatchLargeClass")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(source, "MatchLargeClassIntoRangeBits") {
		t.Fatal("large character class was expanded into a Go boolean expression instead of a bitmap")
	}
}

func TestCompileToGoSmoke(t *testing.T) {
	for _, pattern := range []string{"a", "ab", "a|b", "[a-z]", "\\d+", "\\w+", ".", "[^a]", "^a$", "(a)b", "a*", "a+", "a?", "a*?", "a{2,3}", "a(b|c)", "(?<name>a)"} {
		source, err := CompileToGo(pattern, "", "Match")
		if err != nil {
			t.Errorf("CompileToGo(%q): %v", pattern, err)
			continue
		}
		fileSet := token.NewFileSet()
		file, err := goparser.ParseFile(fileSet, "generated.go", "package generated\n"+source, goparser.AllErrors)
		if err != nil {
			t.Errorf("generated source for %q does not parse: %v", pattern, err)
			continue
		}
		if _, err := (&types.Config{Importer: importer.Default()}).Check("generated", fileSet, []*ast.File{file}, nil); err != nil {
			t.Errorf("generated source for %q does not typecheck: %v\n%s", pattern, err, strings.TrimSpace(source))
		}
	}
}

func TestCompileToGoFactorsBinaryCommonPrefix(t *testing.T) {
	source, err := CompileToGo("abcde|abcxy", "", "MatchPrefix")
	if err != nil {
		t.Fatal(err)
	}
	for _, character := range []int{'a', 'b', 'c'} {
		if got := strings.Count(source, fmt.Sprintf("!= %d", character)); got != 1 {
			t.Errorf("shared prefix character %q emitted %d times; want one direct check", rune(character), got)
		}
	}
}

func TestCompileToGoLargeFiniteStringSetUsesCompactMatcherAndBacktracks(t *testing.T) {
	var pattern strings.Builder
	pattern.WriteString(`[\q{a|ab`)
	for index := 0; index < 64; index++ {
		fmt.Fprintf(&pattern, `|x%03d`, index)
	}
	pattern.WriteString(`}]b`)
	matcher, err := CompileToGo(pattern.String(), "v", "MatchLargeStringSet")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(matcher, "MatchLargeStringSetIntoStringSet0Nodes") {
		t.Fatal("large finite string set was not lowered to a compact automaton table")
	}
	if len(matcher) >= 32*1024 {
		t.Fatalf("finite string-set matcher stayed unexpectedly large: %d bytes", len(matcher))
	}

	var program strings.Builder
	program.WriteString("package main\nfunc units(s string) []uint16 { out := make([]uint16, len(s)); for i := range s { out[i] = uint16(s[i]) }; return out }\n")
	program.WriteString(matcher)
	program.WriteString("func main() { _, end, _, ok := MatchLargeStringSet(units(\"ab\"), 0); if !ok || end != 2 { panic(\"finite string-set backtracking failed\") } }\n")
	directory := t.TempDir()
	file := filepath.Join(directory, "main.go")
	if err := os.WriteFile(file, []byte(program.String()), 0o600); err != nil {
		t.Fatal(err)
	}
	goBin := filepath.Join(runtime.GOROOT(), "bin", "go")
	command := exec.Command(goBin, "run", file)
	command.Env = append(os.Environ(), "GO111MODULE=off", "GOCACHE=/workspace/go-build-cache")
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("generated finite string-set matcher failed: %v\n%s", err, output)
	}
}

func TestCompileToGoQuantifiedRGIEmojiStringProperty(t *testing.T) {
	const handshake = "🫱🏻‍🫲🏼"
	var program strings.Builder
	program.WriteString("package main\nimport \"unicode/utf16\"\n")
	patterns := []struct{ pattern, name string }{
		{`\p{RGI_Emoji}`, "MatchUnanchored"},
		{`^\p{RGI_Emoji}$`, "MatchAnchored"},
		{`^\p{RGI_Emoji}+$`, "MatchRGIEmoji"},
	}
	for _, test := range patterns {
		matcher, err := CompileToGo(test.pattern, "v", test.name)
		if err != nil {
			t.Fatalf("CompileToGo(%q): %v", test.pattern, err)
		}
		program.WriteString(matcher)
	}
	fmt.Fprintf(&program, "func main() { input := utf16.Encode([]rune(%s)); ", strconv.Quote(handshake))
	for _, test := range patterns {
		fmt.Fprintf(&program, "{ start, end, _, ok := %s(input, 0); if !ok || start != 0 || end != len(input) { panic(%s) } }; ", test.name, strconv.Quote(test.pattern+" failed to match full RGI sequence"))
	}
	program.WriteString("}\n")
	directory := t.TempDir()
	file := filepath.Join(directory, "main.go")
	if err := os.WriteFile(file, []byte(program.String()), 0o600); err != nil {
		t.Fatal(err)
	}
	goBin := filepath.Join(runtime.GOROOT(), "bin", "go")
	command := exec.Command(goBin, "run", file)
	command.Env = append(os.Environ(), "GO111MODULE=off", "GOCACHE=/workspace/go-build-cache")
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("generated quantified RGI matcher failed: %v\n%s", err, output)
	}
}

func TestCompileToGoExecutesGeneratedMatcher(t *testing.T) {
	cases := []struct {
		pattern string
		input   string
		start   int
		wantOK  bool
		want    int
		end     int
	}{
		{"a", "ba", 0, true, 1, 2},
		{"a|b", "b", 0, true, 0, 1},
		{"a|ab", "ab", 0, true, 0, 1},
		{"abc|abcdef", "zabcdef", 0, true, 1, 4},
		{"abcdef|abcxyz", "zabcxyz", 0, true, 1, 7},
		{"[a-z]", "7x", 0, true, 1, 2},
		{"^a$", "a", 0, true, 0, 1},
		{"^a$", "ba", 0, false, 0, 0},
		{"a$", "a\n", 0, false, 0, 0},
		{"a$", "a\r\n", 0, false, 0, 0},
		{"(a)b", "zab", 0, true, 1, 3},
		{"a*", "aaab", 0, true, 0, 3},
		{"a+", "baa", 0, true, 1, 3},
		{"a?", "b", 0, true, 0, 0},
		{"a*?", "aaa", 0, true, 0, 0},
		{"a{2,3}", "baaa", 0, true, 1, 4},
		{"a(b|c)", "zac", 0, true, 1, 3},
		{"(a)\\1", "zaa", 0, true, 1, 3},
	}
	var source strings.Builder
	source.WriteString("package main\nfunc units(s string) []uint16 { out := make([]uint16, len(s)); for i := range s { out[i] = uint16(s[i]) }; return out }\n")
	for index, test := range cases {
		matcher, err := CompileToGo(test.pattern, "", "Match"+itoa(index))
		if err != nil {
			t.Fatalf("CompileToGo(%q): %v", test.pattern, err)
		}
		source.WriteString(matcher)
	}
	source.WriteString("func main() {\n")
	for index, test := range cases {
		start, end, ok := "matchStart"+itoa(index), "matchEnd"+itoa(index), "matched"+itoa(index)
		call := start + ", " + end + ", _, " + ok + " := Match" + itoa(index) + "(units(" + strconv.Quote(test.input) + "), " + itoa(test.start) + "); "
		call += "if " + ok + " != " + strconv.FormatBool(test.wantOK) + " { panic(\"unexpected match status\") }; "
		if test.wantOK {
			call += "if " + start + " != " + itoa(test.want) + " || " + end + " != " + itoa(test.end) + " { panic(\"unexpected match span\") }; "
		}
		call += "_ = " + start + "; _ = " + end + "; "
		source.WriteString(call + "\n")
	}
	source.WriteString("}\n")
	dir := t.TempDir()
	file := filepath.Join(dir, "main.go")
	if err := os.WriteFile(file, []byte(source.String()), 0o600); err != nil {
		t.Fatal(err)
	}
	goBin := filepath.Join(runtime.GOROOT(), "bin", "go")
	command := exec.Command(goBin, "run", file)
	command.Env = append(os.Environ(), "GO111MODULE=off", "GOCACHE=/workspace/go-build-cache")
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("generated matcher program failed: %v\n%s", err, output)
	}
}

func itoa(value int) string { return strconv.Itoa(value) }
