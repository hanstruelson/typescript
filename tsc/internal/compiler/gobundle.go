package compiler

import (
	"context"
	"fmt"
	pathpkg "path"
	"sort"
	"strings"

	"github.com/microsoft/TypeScript/tsc/internal/ast"
	"github.com/microsoft/TypeScript/tsc/internal/checker"
	"github.com/microsoft/TypeScript/tsc/internal/core"
	"github.com/microsoft/TypeScript/tsc/internal/diagnostics"
	"github.com/microsoft/TypeScript/tsc/internal/goemit"
	"github.com/microsoft/TypeScript/tsc/internal/outputpaths"
	"github.com/microsoft/TypeScript/tsc/internal/parser"
	"github.com/microsoft/TypeScript/tsc/internal/tspath"
)

func needsGoBundle(files []*ast.SourceFile) bool {
	if len(files) > 1 {
		return true
	}
	for _, file := range files {
		needsBundle := false
		var visit func(*ast.Node)
		visit = func(node *ast.Node) {
			if needsBundle {
				return
			}
			switch node.Kind {
			case ast.KindImportDeclaration:
				name := node.AsImportDeclaration().ModuleSpecifier.Text()
				if !goemit.IsNativeNodeModule(name) || isGoGlobImport(name) {
					needsBundle = true
					return
				}
			case ast.KindExportDeclaration, ast.KindExportAssignment:
				needsBundle = true
				return
			case ast.KindCallExpression:
				call := node.AsCallExpression()
				if (ast.IsRequireCall(node, false) || call.Expression.Kind == ast.KindImportKeyword) && len(call.Arguments.Nodes) == 1 && call.Arguments.Nodes[0].Kind == ast.KindStringLiteral {
					name := call.Arguments.Nodes[0].Text()
					if !goemit.IsNativeNodeModule(name) || isGoGlobImport(name) {
						needsBundle = true
						return
					}
				}
			}
			if ast.HasSyntacticModifier(node, ast.ModifierFlagsExport) || node.Kind == ast.KindExportDeclaration || node.Kind == ast.KindExportAssignment {
				needsBundle = true
				return
			}
			node.ForEachChild(func(child *ast.Node) bool {
				visit(child)
				return false
			})
		}
		visit(file.AsNode())
		if needsBundle {
			return true
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
	files = append([]*ast.SourceFile(nil), files...)
	byName := make(map[string]*ast.SourceFile, len(files))
	targets := make(map[*ast.Node]string)
	globs := make(map[*ast.Node][]goemit.GlobFile)
	bundleOnly := make(map[string]bool)
	imported := make(map[string]bool)
	moduleResolver := p.newResolver()
	for _, file := range files {
		byName[file.FileName().AsString()] = file
	}
	// Package JavaScript files are normally excluded from TypeScript's emitted
	// file list. A Go bundle must include their source and their static require
	// dependencies so imports from node_modules have executable targets.
	for index := 0; index < len(files); index++ {
		file := files[index]
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
				} else if ast.IsRequireCall(node, false) && len(call.Arguments.Nodes) == 1 && call.Arguments.Nodes[0].Kind == ast.KindStringLiteral {
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
				if goemit.IsNativeNodeModule(specifier.Text()) {
					node.ForEachChild(func(child *ast.Node) bool { visit(child); return false })
					return
				}
				resolved := p.GetResolvedModuleFromModuleSpecifier(file, specifier)
				if resolved == nil || !resolved.IsResolved() {
					// JavaScript files excluded from the normal emit graph may still
					// contain static CommonJS dependencies. Resolve those here so the
					// Go bundle can include the same dependency graph.
					resolved, _, _ = moduleResolver.ResolveModuleName(
						specifier.Text(),
						file.FileName(),
						p.GetModeForUsageLocation(file, specifier),
						nil,
					)
				}
				if resolved != nil && resolved.IsResolved() {
					target := p.GetSourceFileForResolvedModule(resolved)
					if target == nil || target.IsDeclarationFile {
						target = p.getGoRuntimeSourceForDeclaration(resolved.ResolvedFileName)
					}
					if target == nil && !resolved.ResolvedFileName.ExtensionIsOneOf([]string{
						tspath.ExtensionJson,
						tspath.ExtensionDts,
						tspath.ExtensionDcts,
						tspath.ExtensionDmts,
					}) {
						target = p.getGoSourceFileForBundle(resolved.ResolvedFileName)
					}
					if (target == nil || target.IsDeclarationFile) && resolved.ResolvedFileName.ExtensionIs(tspath.ExtensionJson) {
						target = p.getGoJSONModuleSource(resolved.ResolvedFileName)
					}
					if target != nil && target.FileName().ExtensionIs(tspath.ExtensionJson) {
						target = p.getGoJSONModuleSource(target.FileName())
					}
					if target != nil && !target.IsDeclarationFile && !target.FileName().ExtensionIsOneOf([]string{tspath.ExtensionDts, tspath.ExtensionJson}) {
						name := target.FileName().AsString()
						targets[node] = name
						imported[name] = true
						if byName[name] == nil {
							bundleOnly[name] = strings.HasSuffix(name, ".goemit.js") || p.GetSourceFileByPath(target.PathKey()) == nil
							if bundleOnly[name] {
								p.registerGoBundleSourceFile(target)
							}
							byName[name] = target
							files = append(files, target)
						}
					}
				}
			}
			node.ForEachChild(func(child *ast.Node) bool { visit(child); return false })
		}
		visit(file.AsNode())
	}
	var bundleResolver goemit.ReferenceResolver
	for _, isBundleOnly := range bundleOnly {
		if isBundleOnly {
			bundleChecker, _ := checker.NewChecker(p, nil)
			bundleResolver = bundleChecker.GetEmitResolver()
			break
		}
	}
	names := make([]string, len(files))
	fragments := make([]string, len(files))
	for index, file := range files {
		names[index] = file.FileName().AsString()
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
		var resolver goemit.ReferenceResolver
		done := func() {}
		if bundleOnly[file.FileName().AsString()] {
			resolver = bundleResolver
		} else {
			host, release := newEmitHost(ctx, p, file)
			resolver = host.GetEmitResolver()
			done = release
		}
		commonJS := file.FileName().ExtensionIsOneOf([]string{tspath.ExtensionJs, tspath.ExtensionCjs}) &&
			p.GetEmitModuleFormatOfFile(file) == core.ModuleKindCommonJS
		fragment, diags := goemit.EmitModule(goemit.ModuleInput{File: file, Resolver: resolver, Targets: targets, Globs: globs, CommonJS: commonJS}, p.Options())
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

// A resolved JavaScript dependency can be missing from the Program's regular
// source set when it was discovered only through a package's CommonJS require.
// Read and parse that source for the static Go bundle without pulling in Node.
func (p *Program) getGoSourceFileForBundle(fileName tspath.RootedFilePath) *ast.SourceFile {
	path := p.PathKeyForFileName(fileName)
	if existing := p.GetSourceFileByPath(path); existing != nil {
		return existing
	}
	if !p.Host().FS().FileExists(fileName) {
		return nil
	}
	return p.Host().GetSourceFile(ast.SourceFileParseOptions{
		FileName: fileName,
		PathKey:  path,
		ExternalModuleIndicatorOptions: ast.GetExternalModuleIndicatorOptions(
			fileName,
			p.Options(),
			p.GetSourceFileMetaData(path),
		),
	})
}

// A package's type declaration can resolve successfully even when its runtime
// JavaScript is excluded from the Program's normal source set. Go bundles need
// executable module bodies, so use the matching sibling runtime file when one
// exists (for example index.d.ts -> index.js).
func (p *Program) getGoRuntimeSourceForDeclaration(declaration tspath.RootedFilePath) *ast.SourceFile {
	name := declaration.AsString()
	var candidates []string
	switch {
	case strings.HasSuffix(name, ".d.ts"):
		stem := strings.TrimSuffix(name, ".d.ts")
		if strings.HasSuffix(stem, ".js") || strings.HasSuffix(stem, ".cjs") || strings.HasSuffix(stem, ".mjs") {
			candidates = []string{stem}
		} else {
			candidates = []string{stem + ".js", stem + ".cjs", stem + ".mjs"}
		}
	case strings.HasSuffix(name, ".d.cts"):
		candidates = []string{strings.TrimSuffix(name, ".d.cts") + ".cjs", strings.TrimSuffix(name, ".d.cts") + ".js"}
	case strings.HasSuffix(name, ".d.mts"):
		candidates = []string{strings.TrimSuffix(name, ".d.mts") + ".mjs", strings.TrimSuffix(name, ".d.mts") + ".js"}
	}
	for _, candidateName := range candidates {
		candidate := tspath.RootedFilePathFromAbsolute(candidateName)
		if !p.Host().FS().FileExists(candidate) {
			continue
		}
		path := p.PathKeyForFileName(candidate)
		if existing := p.GetSourceFileByPath(path); existing != nil {
			return existing
		}
		return p.Host().GetSourceFile(ast.SourceFileParseOptions{
			FileName: candidate,
			PathKey:  path,
			ExternalModuleIndicatorOptions: ast.GetExternalModuleIndicatorOptions(
				candidate,
				p.Options(),
				p.GetSourceFileMetaData(path),
			),
		})
	}
	return nil
}

// JSON package dependencies are static data, so lower them to a tiny
// CommonJS module that parses the original JSON literal at initialization.
func (p *Program) getGoJSONModuleSource(jsonFile tspath.RootedFilePath) *ast.SourceFile {
	virtualFile := tspath.RootedFilePathFromAbsolute(jsonFile.AsString() + ".goemit.js")
	path := p.PathKeyForFileName(virtualFile)
	if existing := p.GetSourceFileByPath(path); existing != nil {
		return existing
	}
	content, ok := p.Host().FS().ReadFile(jsonFile)
	if !ok {
		return nil
	}
	source := "module.exports = (" + string(content) + ");"
	file := parser.ParseSourceFile(ast.SourceFileParseOptions{FileName: virtualFile, PathKey: path}, source, core.ScriptKindJS)
	p.registerGoBundleSourceFile(file)
	return file
}

// Register a runtime source discovered from a package declaration so an
// auxiliary checker can bind it after graph discovery and before Go emission.
func (p *Program) registerGoBundleSourceFile(file *ast.SourceFile) {
	if p.filesByPath == nil {
		p.filesByPath = make(map[tspath.PathKey]*ast.SourceFile)
	}
	if _, exists := p.filesByPath[file.PathKey()]; exists {
		return
	}
	p.filesByPath[file.PathKey()] = file
	p.files = append(p.files, file)
	if p.sourceFileMetaDatas == nil {
		p.sourceFileMetaDatas = make(map[tspath.PathKey]ast.SourceFileMetaData)
	}
	p.sourceFileMetaDatas[file.PathKey()] = ast.SourceFileMetaData{}
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
