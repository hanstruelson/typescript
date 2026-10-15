package compiler_test

import (
	"context"
	"strings"
	"testing"

	"github.com/microsoft/TypeScript/tsc/internal/bundled"
	"github.com/microsoft/TypeScript/tsc/internal/compiler"
	"github.com/microsoft/TypeScript/tsc/internal/core"
	"github.com/microsoft/TypeScript/tsc/internal/tsoptions"
	"github.com/microsoft/TypeScript/tsc/internal/tspath"
	"github.com/microsoft/TypeScript/tsc/internal/vfs/vfstest"
)

func emitGoModuleBundle(t *testing.T, sources map[string]string) string {
	t.Helper()
	fs := bundled.WrapFS(vfstest.FromMap(sources, tspath.CaseSensitive))
	options := core.CompilerOptions{Target: core.ScriptTargetGo, OutDir: "/out", Strict: core.TSTrue}
	config := tsoptions.NewParsedCommandLine(&options, []tspath.RootedFilePath{"/src/main.ts"}, nil, "/src", fs.CaseSensitivity())
	program := compiler.NewProgram(compiler.ProgramOptions{Config: config, Host: compiler.NewCompilerHost(fs, bundled.LibPath(), nil, nil, nil)})
	outputs := map[string]string{}
	result := program.Emit(context.Background(), compiler.EmitOptions{WriteFile: func(path tspath.RootedFilePath, text string, data *compiler.WriteFileData) error {
		outputs[path.AsString()] = text
		return nil
	}})
	if result.EmitSkipped || len(result.Diagnostics) != 0 {
		for _, diagnostic := range result.Diagnostics {
			t.Log(diagnostic.String())
		}
		t.Fatal("Go module bundle emit failed")
	}
	if len(outputs) != 1 || outputs["/out/main.go"] == "" {
		t.Fatalf("expected one bundled Go file at /out/main.go, got outputs %v", outputs)
	}
	return outputs["/out/main.go"]
}

func TestGoDemandDrivenModuleEvaluation(t *testing.T) {
	for _, tc := range []struct {
		name      string
		sources   map[string]string
		want      string
		wantError string
	}{
		{
			name: "selected exports suspend and resume live bindings",
			sources: map[string]string{
				"/src/main.ts": `import { first } from "./dep.js";console.log("main:first",first);import { second } from "./dep.js";console.log("main:second",first,second);import "./dep.js";console.log("main:live",first);`,
				"/src/dep.ts":  `export let first=1;console.log("dependency:middle");first=2;export const second=first+1;console.log("dependency:tail");first=4;`,
			},
			want: "main:first 1\ndependency:middle\nmain:second 2 3\ndependency:tail\nmain:live 4\n",
		},
		{
			name: "named export set waits for every requested name",
			sources: map[string]string{
				"/src/main.ts": `import { first, second } from "./dep.js";console.log("main:both",first,second);import "./dep.js";console.log("main:done");`,
				"/src/dep.ts":  `export const first=1;console.log("dependency:first");export const second=2;console.log("dependency:tail");`,
			},
			want: "dependency:first\nmain:both 1 2\ndependency:tail\nmain:done\n",
		},
		{
			name: "cycle succeeds when the requested export is ready",
			sources: map[string]string{
				"/src/main.ts": `export const name="A";import { suffix } from "./dep.js";console.log(suffix);`,
				"/src/dep.ts":  `import { name } from "./main.js";export const suffix=name+"B";`,
			},
			want: "AB\n",
		},
		{
			name: "cycle without a ready export reports instead of hanging",
			sources: map[string]string{
				"/src/main.ts": `import { value } from "./dep.js";export const ready="main";console.log(value);`,
				"/src/dep.ts":  `import { ready } from "./main.js";export const value=ready;`,
			},
			wantError: "Unresolved module dependency",
		},
		{
			name: "renamed re-export requests the source name",
			sources: map[string]string{
				"/src/main.ts":     `import { alias } from "./reexport.js";console.log("main",alias);`,
				"/src/reexport.ts": `export { original as alias } from "./dep.js";console.log("re-export remainder");`,
				"/src/dep.ts":      `export const original=7;console.log("dependency remainder");`,
			},
			want: "main 7\n",
		},
		{
			name: "namespace import requests module completion",
			sources: map[string]string{
				"/src/main.ts": `import * as dep from "./dep.js";console.log("main",dep.value);`,
				"/src/dep.ts":  `export const value=1;console.log("dependency remainder");`,
			},
			want: "dependency remainder\nmain 1\n",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			text := emitGoModuleBundle(t, tc.sources)
			output, err := runGoProgram(t, text, false)
			if tc.wantError != "" {
				if err == nil || !strings.Contains(output, tc.wantError) {
					t.Fatalf("got %q (%v), wanted an error containing %q", output, err, tc.wantError)
				}
				return
			}
			if err != nil || output != tc.want {
				t.Fatalf("got %q (%v), want %q", output, err, tc.want)
			}
		})
	}
}
