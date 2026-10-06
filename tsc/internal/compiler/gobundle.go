package compiler

import (
	"context"
	"fmt"

	"github.com/microsoft/TypeScript/tsc/internal/ast"
	"github.com/microsoft/TypeScript/tsc/internal/diagnostics"
	"github.com/microsoft/TypeScript/tsc/internal/goemit"
	"github.com/microsoft/TypeScript/tsc/internal/outputpaths"
)

func needsGoBundle(files []*ast.SourceFile) bool {
	if len(files) > 1 {
		return true
	}
	for _, file := range files {
		for _, node := range file.Statements.Nodes {
			if ast.HasSyntacticModifier(node, ast.ModifierFlagsExport) || node.Kind == ast.KindExportDeclaration || node.Kind == ast.KindExportAssignment {
				return true
			}
			if node.Kind == ast.KindImportDeclaration {
				name := node.AsImportDeclaration().ModuleSpecifier.Text()
				if name != "node:fs/promises" && name != "fs/promises" {
					return true
				}
			}
		}
	}
	return false
}

// emitGoBundle borrows each file's normal emit resolver sequentially. It resolves
// the complete graph with the Program and writes a single Go executable, rather
// than independently executing and emitting every dependency as a root.
func (p *Program) emitGoBundle(ctx context.Context, options EmitOptions, files []*ast.SourceFile) *EmitResult {
	result := &EmitResult{}
	if len(files) == 0 {
		return result
	}
	names := make([]string, len(files))
	fragments := make([]string, len(files))
	byName := make(map[string]*ast.SourceFile, len(files))
	targets := make(map[*ast.Node]string)
	imported := make(map[string]bool)
	for index, file := range files {
		names[index] = file.FileName().AsString()
		byName[names[index]] = file
	}
	for _, file := range files {
		var visit func(*ast.Node)
		visit = func(node *ast.Node) {
			var specifier *ast.Node
			switch node.Kind {
			case ast.KindImportDeclaration:
				specifier = node.AsImportDeclaration().ModuleSpecifier
			case ast.KindExportDeclaration:
				specifier = node.AsExportDeclaration().ModuleSpecifier
			case ast.KindCallExpression:
				call := node.AsCallExpression()
				if call.Expression.Kind == ast.KindImportKeyword && len(call.Arguments.Nodes) == 1 && call.Arguments.Nodes[0].Kind == ast.KindStringLiteral {
					specifier = call.Arguments.Nodes[0]
				}
			}
			if specifier != nil {
				resolved := p.GetResolvedModuleFromModuleSpecifier(file, specifier)
				if resolved != nil && resolved.IsResolved() {
					target := p.GetSourceFileForResolvedModule(resolved)
					if target != nil && byName[target.FileName().AsString()] != nil {
						name := target.FileName().AsString()
						targets[node] = name
						imported[name] = true
					}
				}
			}
			node.ForEachChild(func(child *ast.Node) bool { visit(child); return false })
		}
		visit(file.AsNode())
	}

	var candidates []*ast.SourceFile
	for _, name := range p.opts.Config.FileNames() {
		if file := byName[name.AsString()]; file != nil {
			candidates = append(candidates, file)
		}
	}
	if len(candidates) == 0 {
		candidates = options.TargetSourceFiles
	}
	if len(candidates) == 0 {
		candidates = files[:1]
	}
	roots := []string{}
	for _, file := range candidates {
		name := file.FileName().AsString()
		if byName[name] != nil && !imported[name] {
			roots = append(roots, name)
		}
	}
	// An entirely circular entry graph still needs a starting module. Use the
	// first configured entry, not dependency traversal order.
	if len(roots) == 0 {
		roots = append(roots, candidates[0].FileName().AsString())
	}
	entry := byName[roots[0]]
	path := outputpaths.GetOutputPathsFor(entry, p.Options(), p, outputpaths.ForceEmitPaths{Js: options.ForceEmit}).JsFilePath()
	if path == "" {
		return result
	}
	if !options.ForceEmit && p.IsEmitBlocked(path) {
		result.EmitSkipped = true
		return result
	}
	for index, file := range files {
		if ctx.Err() != nil {
			result.EmitSkipped = true
			return result
		}
		host, done := newEmitHost(ctx, p, file)
		fragment, diags := goemit.EmitModule(goemit.ModuleInput{File: file, Resolver: host.GetEmitResolver(), Targets: targets}, p.Options())
		done()
		fragments[index] = fragment
		result.Diagnostics = append(result.Diagnostics, diags...)
	}
	if len(result.Diagnostics) > 0 {
		result.EmitSkipped = true
		return result
	}
	text, err := goemit.Bundle(names, fragments, roots)
	if err != nil {
		result.Diagnostics = append(result.Diagnostics, ast.NewCompilerDiagnostic(diagnostics.Could_not_write_file_0_Colon_1, path, fmt.Sprintf("Go bundle generation: %v", err)))
		result.EmitSkipped = true
		return result
	}
	if options.WriteFile != nil {
		err = options.WriteFile(path, text, &WriteFileData{SourceFile: entry})
	} else {
		err = p.Host().FS().WriteFile(path, text)
	}
	if err != nil {
		result.Diagnostics = append(result.Diagnostics, ast.NewCompilerDiagnostic(diagnostics.Could_not_write_file_0_Colon_1, path, err.Error()))
		result.EmitSkipped = true
		return result
	}
	result.EmittedFiles = append(result.EmittedFiles, path)
	if p.Options().GetEmitDeclarations() && options.EmitOnly == EmitAll {
		declarationOptions := options
		declarationOptions.EmitOnly = EmitOnlyDts
		return CombineEmitResults([]*EmitResult{result, p.Emit(ctx, declarationOptions)})
	}
	return result
}
