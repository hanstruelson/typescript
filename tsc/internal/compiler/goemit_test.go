package compiler_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/microsoft/TypeScript/tsc/internal/bundled"
	"github.com/microsoft/TypeScript/tsc/internal/compiler"
	"github.com/microsoft/TypeScript/tsc/internal/core"
	"github.com/microsoft/TypeScript/tsc/internal/tsoptions"
	"github.com/microsoft/TypeScript/tsc/internal/tspath"
	"github.com/microsoft/TypeScript/tsc/internal/vfs/vfstest"
)

func TestGoTarget(t *testing.T) {
	for _, tc := range []struct {
		name, source      string
		args              []string
		wantError, noEmit bool
	}{
		{name: "numeric main", source: "var x: number = 10; console.log(x);"},
		{name: "multiple arguments and keywords", source: "let type: number = -2.5; const fmt: number = 10; console.log(type, fmt);"},
		{name: "runtime without logging", source: "var x: number = 10;"},
		{name: "unsupported", source: "class Unsupported { get x(){return 1;} }", wantError: true},
		{name: "unsupported expression", source: "console.log({get value(){return 10;}});", wantError: true},
		{name: "source map", source: "console.log(10);", args: []string{"--sourceMap"}, wantError: true},
		{name: "noEmit", source: "console.log(10);", args: []string{"--noEmit"}, noEmit: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fs := bundled.WrapFS(vfstest.FromMap(map[string]string{"/src/input.ts": tc.source}, tspath.CaseSensitive))
			args := append([]string{"--target", "go", "--outDir", "/out", "/src/input.ts"}, tc.args...)
			config := tsoptions.ParseCommandLine(args, fs, "/src")
			if len(config.Errors) != 0 || config.CompilerOptions().Target != core.ScriptTargetGo {
				t.Fatalf("target parsing failed: %v", config.Errors)
			}
			p := compiler.NewProgram(compiler.ProgramOptions{Config: config, Host: compiler.NewCompilerHost(fs, bundled.LibPath(), nil, nil, nil)})
			outputs := map[string]string{}
			result := p.Emit(context.Background(), compiler.EmitOptions{WriteFile: func(path tspath.RootedFilePath, text string, data *compiler.WriteFileData) error {
				outputs[path.AsString()] = text
				return nil
			}})
			if tc.wantError {
				if len(result.Diagnostics) == 0 || !result.EmitSkipped || len(outputs) != 0 {
					t.Fatalf("expected failed emit: %+v, %v", result, outputs)
				}
				return
			}
			if len(result.Diagnostics) != 0 {
				t.Fatalf("emit diagnostics: %v", result.Diagnostics)
			}
			if tc.noEmit {
				if len(outputs) != 0 {
					t.Fatal("noEmit wrote output")
				}
				return
			}
			text := outputs["/out/input.go"]
			if len(outputs) != 1 || !strings.Contains(text, "func main()") || !strings.Contains(text, "func tsConsoleLog") {
				t.Fatalf("bad output: %v", outputs)
			}
			path := filepath.Join(t.TempDir(), "main.go")
			if err := os.WriteFile(path, []byte(text), 0600); err != nil {
				t.Fatal(err)
			}
			cmd := exec.Command("go", "run", path)
			cmd.Env = append(os.Environ(), "GOWORK=off")
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("generated Go failed: %v\n%s", err, output)
			}
			want := "10\n"
			if tc.name == "multiple arguments and keywords" {
				want = "-2.5 10\n"
			}
			if tc.name == "runtime without logging" {
				want = ""
			}
			if string(output) != want {
				t.Fatalf("output %q, want %q", output, want)
			}
		})
	}
}

func TestGoTargetTsconfig(t *testing.T) {
	fs := vfstest.FromMap(map[string]string{"/src/input.ts": "var x: number = 10; console.log(x);"}, tspath.CaseSensitive)
	config := tsoptions.ParseJsonConfigFileContent(map[string]any{"compilerOptions": map[string]any{"target": "go", "outDir": "out"}, "files": []any{"input.ts"}}, fs, "/src", nil, "/src/tsconfig.json", nil, nil)
	if len(config.Errors) != 0 || config.CompilerOptions().Target != core.ScriptTargetGo {
		t.Fatalf("config target parsing failed: %v", config.Errors)
	}
}
