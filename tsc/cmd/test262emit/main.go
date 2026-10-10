// test262emit emits JavaScript fixtures without suppressing unsupported diagnostics.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/microsoft/TypeScript/tsc/internal/ast"
	"github.com/microsoft/TypeScript/tsc/internal/bundled"
	"github.com/microsoft/TypeScript/tsc/internal/compiler"
	"github.com/microsoft/TypeScript/tsc/internal/core"
	"github.com/microsoft/TypeScript/tsc/internal/tsoptions"
	"github.com/microsoft/TypeScript/tsc/internal/tspath"
	"github.com/microsoft/TypeScript/tsc/internal/vfs/vfstest"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: test262emit input.js output.go")
		os.Exit(2)
	}
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	fs := bundled.WrapFS(vfstest.FromMap(map[string]string{"/src/input.js": string(data)}, tspath.CaseSensitive))
	options := core.CompilerOptions{Target: core.ScriptTargetGo, OutDir: "/out", AllowJs: core.TSTrue}
	p := compiler.NewProgram(compiler.ProgramOptions{Config: tsoptions.NewParsedCommandLine(&options, []tspath.RootedFilePath{"/src/input.js"}, nil, "/src", fs.CaseSensitivity()), Host: compiler.NewCompilerHost(fs, bundled.LibPath(), nil, nil, nil)})
	diagnostics := p.GetSyntacticDiagnostics(context.Background(), nil)
	if len(diagnostics) > 0 {
		report(diagnostics)
		os.Exit(3)
	}
	result := p.Emit(context.Background(), compiler.EmitOptions{WriteFile: func(path tspath.RootedFilePath, text string, _ *compiler.WriteFileData) error {
		return os.WriteFile(os.Args[2], []byte(text), 0600)
	}})
	if result.EmitSkipped || len(result.Diagnostics) > 0 {
		report(result.Diagnostics)
		os.Exit(4)
	}
}
func report(diagnostics []*ast.Diagnostic) {
	for _, d := range diagnostics {
		b, _ := json.Marshal(map[string]string{"diagnostic": d.String()})
		fmt.Fprintln(os.Stderr, string(b))
	}
}
