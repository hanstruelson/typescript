package compiler

import (
	"context"
	"fmt"
	pathpkg "path"
	"sort"
	"strings"

	"github.com/microsoft/TypeScript/tsc/internal/ast"
	"github.com/microsoft/TypeScript/tsc/internal/diagnostics"
	"github.com/microsoft/TypeScript/tsc/internal/goemit"
	"github.com/microsoft/TypeScript/tsc/internal/outputpaths"
	"github.com/microsoft/TypeScript/tsc/internal/tspath"
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
				if !goemit.IsNativeNodeModule(name) || isGoGlobImport(name) {
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
	globs := make(map[*ast.Node][]goemit.GlobFile)
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
				if node.Kind == ast.KindImportDeclaration && isGoGlobImport(specifier.Text()) {
					matches := expandGoGlobImport(file, specifier.Text(), byName)
					globs[node] = matches
					for _, match := range matches {
						imported[match.Name] = true
					}
					return
				}
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
		fragment, diags := goemit.EmitModule(goemit.ModuleInput{File: file, Resolver: host.GetEmitResolver(), Targets: targets, Globs: globs}, p.Options())
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

func isGoGlobImport(specifier string) bool {
	return strings.ContainsAny(specifier, "*?[")
}

func expandGoGlobImport(importer *ast.SourceFile, specifier string, sources map[string]*ast.SourceFile) []goemit.GlobFile {
	pattern := goGlobImportPattern(importer, specifier)
	root := goGlobRoot(pattern)

	names := make([]string, 0, len(sources))
	for name := range sources {
		names = append(names, name)
	}
	sort.Strings(names)
	matches := make([]goemit.GlobFile, 0)
	for _, name := range names {
		if !goGlobMatch(pattern, name) {
			continue
		}
		relative, ok := goGlobRelative(root, name)
		if !ok {
			continue
		}
		key := strings.TrimSuffix(relative, pathpkg.Ext(relative))
		keyParts := strings.Split(key, "/")
		if len(keyParts) == 0 || keyParts[len(keyParts)-1] == "" {
			continue
		}
		file := sources[name]
		hasDefault, exports := goGlobExports(file)
		matches = append(matches, goemit.GlobFile{Key: keyParts, Name: name, HasDefault: hasDefault, Exports: exports})
	}
	return matches
}

func goGlobImportPattern(importer ast.HasFileName, specifier string) string {
	return tspath.NormalizePath(tspath.CombinePaths(importer.FileName().Directory().AsString(), specifier))
}

func goGlobRoot(pattern string) string {
	wildcard := strings.IndexAny(pattern, "*?[")
	if wildcard < 0 {
		return pathpkg.Dir(pattern)
	}
	literalDirectory := pattern[:wildcard]
	if !strings.HasSuffix(literalDirectory, "/") {
		if slash := strings.LastIndexByte(literalDirectory, '/'); slash >= 0 {
			literalDirectory = literalDirectory[:slash+1]
		} else {
			literalDirectory = ""
		}
	}
	if literalDirectory == "" {
		return pathpkg.Dir(pattern)
	}
	return pathpkg.Clean(literalDirectory)
}

func goGlobstarExtensionPattern(pattern string) string {
	parts := strings.Split(pattern, "/")
	for i, part := range parts {
		if strings.HasPrefix(part, "**.") {
			parts[i] = "**/*" + strings.TrimPrefix(part, "**")
		}
	}
	return strings.Join(parts, "/")
}

func goGlobRelative(root, name string) (string, bool) {
	root = pathpkg.Clean(root)
	if root == "/" {
		return strings.TrimPrefix(name, "/"), strings.HasPrefix(name, "/")
	}
	prefix := strings.TrimSuffix(root, "/") + "/"
	if !strings.HasPrefix(name, prefix) {
		return "", false
	}
	return strings.TrimPrefix(name, prefix), true
}

func goGlobExports(file *ast.SourceFile) (bool, []string) {
	symbol := file.AsNode().Symbol()
	if symbol == nil {
		return false, nil
	}
	var exports []string
	hasDefault := goGlobRuntimeExport(symbol.Exports[ast.InternalSymbolNameDefault])
	for name := range symbol.Exports {
		if name == ast.InternalSymbolNameDefault {
			continue
		}
		if name == ast.InternalSymbolNameExportEquals || strings.HasPrefix(name, ast.InternalSymbolNamePrefix) {
			continue
		}
		if !goGlobRuntimeExport(symbol.Exports[name]) {
			continue
		}
		exports = append(exports, name)
	}
	sort.Strings(exports)
	return hasDefault, exports
}

func goGlobRuntimeExport(symbol *ast.Symbol) bool {
	if symbol == nil {
		return false
	}
	if symbol.Flags&ast.SymbolFlagsValue != 0 {
		return true
	}
	if symbol.Flags&ast.SymbolFlagsAlias == 0 {
		return false
	}
	for _, declaration := range symbol.Declarations {
		if ast.IsExportSpecifier(declaration) && declaration.AsExportSpecifier().IsTypeOnly {
			continue
		}
		if ast.IsExportDeclaration(declaration) && declaration.AsExportDeclaration().IsTypeOnly {
			continue
		}
		return true
	}
	return false
}

func goGlobMatch(pattern, name string) bool {
	patternParts := strings.Split(pattern, "/")
	nameParts := strings.Split(name, "/")
	type position struct{ pattern, name int }
	memo := make(map[position]bool)
	visited := make(map[position]bool)
	var match func(int, int) bool
	match = func(patternIndex, nameIndex int) bool {
		key := position{patternIndex, nameIndex}
		if visited[key] {
			return memo[key]
		}
		visited[key] = true
		matched := false
		switch {
		case patternIndex == len(patternParts):
			matched = nameIndex == len(nameParts)
		case patternParts[patternIndex] == "**":
			for index := nameIndex; index <= len(nameParts) && !matched; index++ {
				matched = match(patternIndex+1, index)
			}
		case strings.HasPrefix(patternParts[patternIndex], "**"):
			segmentPattern := "*" + strings.TrimPrefix(patternParts[patternIndex], "**")
			for index := nameIndex; index < len(nameParts) && !matched; index++ {
				if goGlobSegmentMatches(segmentPattern, nameParts[index]) {
					matched = match(patternIndex+1, index+1)
				}
			}
		default:
			if nameIndex < len(nameParts) && goGlobSegmentMatches(patternParts[patternIndex], nameParts[nameIndex]) {
				matched = match(patternIndex+1, nameIndex+1)
			}
		}
		memo[key] = matched
		return matched
	}
	return match(0, 0)
}

func goGlobSegmentMatches(pattern, name string) bool {
	matched, err := pathpkg.Match(pattern, name)
	return err == nil && matched
}
