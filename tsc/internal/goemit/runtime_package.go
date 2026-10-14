package goemit

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"path"
	"sort"
	"strconv"
	"strings"
)

const runtimePackageImport = "github.com/microsoft/TypeScript/tsc/runtime"

type runtimeExportPlan struct {
	globals    map[string]string
	members    map[string]string
	fieldTypes map[string]map[string]bool
}

func runtimeTemplate() string {
	return Runtime + ValueRuntime + ModuleRuntime + ClassRuntime + TypeRuntime + ConversionRuntime +
		GrowableArrayRuntime + StringRuntime + RegexRuntime + EqualityRuntime + ObjectRuntime +
		CollectionRuntime + SetRuntime + ArrayRuntime + ArrayFlattenRuntime + ArrayLikeRuntime +
		ArrayBuiltinRuntime + BufferRuntime + TypedArrayBuiltinRuntime + NativeArrayAccessRuntime +
		DecoratorRuntime
}

func exportName(name string) string {
	if name == "" || ast.IsExported(name) {
		return name
	}
	return strings.ToUpper(name[:1]) + name[1:]
}

func runtimePlan(file *ast.File) (*runtimeExportPlan, error) {
	plan := &runtimeExportPlan{globals: map[string]string{}, members: map[string]string{}, fieldTypes: map[string]map[string]bool{}}
	methodNames := map[string]bool{}
	ast.Inspect(file, func(node ast.Node) bool {
		switch value := node.(type) {
		case *ast.FuncDecl:
			if value.Recv != nil {
				methodNames[value.Name.Name] = true
			}
		case *ast.InterfaceType:
			for _, method := range value.Methods.List {
				for _, name := range method.Names {
					methodNames[name.Name] = true
				}
			}
		}
		return true
	})
	addGlobal := func(name string) error {
		if name == "init" || name == "main" || ast.IsExported(name) {
			return nil
		}
		public := exportName(name)
		for old, other := range plan.globals {
			if old != name && other == public {
				return fmt.Errorf("runtime exports %q and %q as the same name %q", old, name, public)
			}
		}
		plan.globals[name] = public
		return nil
	}
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			if d.Recv == nil {
				if err := addGlobal(d.Name.Name); err != nil {
					return nil, err
				}
			} else if !ast.IsExported(d.Name.Name) {
				plan.members[d.Name.Name] = exportName(d.Name.Name)
			}
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				switch s := spec.(type) {
				case *ast.TypeSpec:
					if err := addGlobal(s.Name.Name); err != nil {
						return nil, err
					}
					if iface, ok := s.Type.(*ast.InterfaceType); ok {
						for _, method := range iface.Methods.List {
							for _, name := range method.Names {
								methodNames[name.Name] = true
							}
						}
					}
				case *ast.ValueSpec:
					for _, name := range s.Names {
						if err := addGlobal(name.Name); err != nil {
							return nil, err
						}
					}
				}
			}
		}
	}
	ast.Inspect(file, func(node ast.Node) bool {
		structure, ok := node.(*ast.StructType)
		if !ok {
			return true
		}
		for _, field := range structure.Fields.List {
			for _, name := range field.Names {
				if ast.IsExported(name.Name) {
					continue
				}
				public := exportName(name.Name)
				// A field and method cannot share the same exported name. Give
				// fields a disjoint prefix where the runtime has that collision.
				if methodNames[public] {
					public = "RuntimeField" + public
				}
				plan.members[name.Name] = public
			}
		}
		return true
	})
	for name := range methodNames {
		if !ast.IsExported(name) {
			if plan.members[name] == "" || plan.members[name] == exportName(name) {
				plan.members[name] = exportName(name)
			}
		}
	}
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range gen.Specs {
			typ, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			if structure, ok := typ.Type.(*ast.StructType); ok {
				fields := map[string]bool{}
				for _, field := range structure.Fields.List {
					for _, name := range field.Names {
						fields[name.Name] = true
					}
				}
				plan.fieldTypes[typ.Name.Name] = fields
			}
		}
	}
	return plan, nil
}

func compositeTypeName(expr ast.Expr) string {
	switch value := expr.(type) {
	case *ast.Ident:
		return value.Name
	case *ast.StarExpr:
		return compositeTypeName(value.X)
	case *ast.IndexExpr:
		return compositeTypeName(value.X)
	case *ast.IndexListExpr:
		return compositeTypeName(value.X)
	}
	return ""
}

func rewriteRuntimeCompositeKeys(file *ast.File, plan *runtimeExportPlan, fset *token.FileSet, edits map[int]string) error {
	var visit func(*ast.CompositeLit, ast.Expr)
	visit = func(literal *ast.CompositeLit, inferred ast.Expr) {
		typeExpr := literal.Type
		if typeExpr == nil {
			typeExpr = inferred
		}
		fields := plan.fieldTypes[compositeTypeName(typeExpr)]
		if structure, ok := typeExpr.(*ast.StructType); ok {
			fields = map[string]bool{}
			for _, field := range structure.Fields.List {
				for _, name := range field.Names {
					fields[name.Name] = true
				}
			}
		}
		if fields != nil {
			for _, element := range literal.Elts {
				item, ok := element.(*ast.KeyValueExpr)
				if !ok {
					continue
				}
				key, ok := item.Key.(*ast.Ident)
				if !ok || !fields[key.Name] {
					continue
				}
				if public := plan.members[key.Name]; public != "" {
					if edits == nil {
						key.Name = public
					} else {
						edits[fset.Position(key.Pos()).Offset] = public
					}
				}
			}
		}
		var elementType ast.Expr
		switch typed := typeExpr.(type) {
		case *ast.MapType:
			elementType = typed.Value
		case *ast.ArrayType:
			elementType = typed.Elt
		case *ast.Ellipsis:
			elementType = typed.Elt
		}
		if elementType != nil {
			for _, element := range literal.Elts {
				if item, ok := element.(*ast.KeyValueExpr); ok {
					if nested, ok := item.Value.(*ast.CompositeLit); ok {
						visit(nested, elementType)
					}
				} else if nested, ok := element.(*ast.CompositeLit); ok {
					visit(nested, elementType)
				}
			}
		}
	}
	ast.Inspect(file, func(node ast.Node) bool {
		if literal, ok := node.(*ast.CompositeLit); ok {
			visit(literal, nil)
		}
		return true
	})
	return nil
}

func runtimeSourceAndPlan(source []byte, packageName string) ([]byte, *runtimeExportPlan, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "runtime.go", source, parser.ParseComments)
	if err != nil {
		return nil, nil, err
	}
	plan, err := runtimePlan(file)
	if err != nil {
		return nil, nil, err
	}

	// Keep composite literal field labels distinct from ordinary identifiers.
	if err := rewriteRuntimeCompositeKeys(file, plan, fset, nil); err != nil {
		return nil, nil, err
	}

	// Export data fields and methods that generated programs access directly.
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			if d.Recv != nil {
				if name := plan.members[d.Name.Name]; name != "" {
					d.Name.Name = name
				}
			}
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				typ, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				if iface, ok := typ.Type.(*ast.InterfaceType); ok {
					for _, method := range iface.Methods.List {
						for _, name := range method.Names {
							if public := plan.members[name.Name]; public != "" {
								name.Name = public
							}
						}
					}
				}
			}
		}
	}
	ast.Inspect(file, func(node ast.Node) bool {
		structure, ok := node.(*ast.StructType)
		if !ok {
			return true
		}
		for _, field := range structure.Fields.List {
			for _, name := range field.Names {
				if public := plan.members[name.Name]; public != "" {
					name.Name = public
				}
			}
		}
		return true
	})
	selectorNames := map[*ast.Ident]bool{}
	keyNames := map[*ast.Ident]bool{}
	ast.Inspect(file, func(node ast.Node) bool {
		switch value := node.(type) {
		case *ast.SelectorExpr:
			if public := plan.members[value.Sel.Name]; public != "" {
				value.Sel.Name = public
			}
			selectorNames[value.Sel] = true
		case *ast.KeyValueExpr:
			if key, ok := value.Key.(*ast.Ident); ok {
				keyNames[key] = true
			}
		}
		return true
	})
	ast.Inspect(file, func(node ast.Node) bool {
		id, ok := node.(*ast.Ident)
		if !ok || selectorNames[id] || keyNames[id] {
			return true
		}
		if public := plan.globals[id.Name]; public != "" {
			id.Name = public
		}
		return true
	})
	file.Name.Name = packageName
	var output bytes.Buffer
	if err := format.Node(&output, fset, file); err != nil {
		return nil, nil, err
	}
	formatted, err := format.Source(output.Bytes())
	return formatted, plan, err
}

func runtimeImportSpecs(source []byte) (map[string]string, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "runtime.go", source, 0)
	if err != nil {
		return nil, err
	}
	imports := map[string]string{}
	for _, spec := range file.Imports {
		pathName, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			return nil, err
		}
		name := path.Base(pathName)
		if spec.Name != nil {
			name = spec.Name.Name
		}
		imports[name] = spec.Path.Value
	}
	return imports, nil
}

func runtimeExportPlanForEmitter() (*runtimeExportPlan, error) {
	file, err := parser.ParseFile(token.NewFileSet(), "runtime.go", []byte("package main\n"+runtimeTemplate()), 0)
	if err != nil {
		return nil, err
	}
	return runtimePlan(file)
}

func applySourceEdits(source []byte, edits map[int]string) []byte {
	source = append([]byte(nil), source...)
	positions := make([]int, 0, len(edits))
	for position := range edits {
		positions = append(positions, position)
	}
	for i := 0; i < len(positions); i++ {
		for j := i + 1; j < len(positions); j++ {
			if positions[j] > positions[i] {
				positions[i], positions[j] = positions[j], positions[i]
			}
		}
	}
	for _, position := range positions {
		end := position
		for end < len(source) && (source[end] == '_' || source[end] >= '0' && source[end] <= '9' || source[end] >= 'a' && source[end] <= 'z' || source[end] >= 'A' && source[end] <= 'Z') {
			end++
		}
		source = append(source[:position], append([]byte(edits[position]), source[end:]...)...)
	}
	return source
}

func transformGeneratedBody(body []byte, plan *runtimeExportPlan, imports map[string]string) ([]byte, error) {
	wrapped := append([]byte("package main\n"), body...)
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "generated.go", wrapped, 0)
	if err != nil {
		line := 0
		_, _ = fmt.Sscanf(err.Error(), "generated.go:%d:", &line)
		lines := bytes.Split(wrapped, []byte("\n"))
		start, end := line-2, line+1
		if start < 1 {
			start = 1
		}
		if end > len(lines) {
			end = len(lines)
		}
		var excerpt strings.Builder
		for index := start; index <= end; index++ {
			fmt.Fprintf(&excerpt, "%d: %s\n", index, lines[index-1])
		}
		return nil, fmt.Errorf("%w\n%s", err, excerpt.String())
	}
	base := len("package main\n")
	fieldTypes := map[string]map[string]bool{}
	for name, fields := range plan.fieldTypes {
		fieldTypes[name] = fields
	}
	for _, decl := range file.Decls {
		if gen, ok := decl.(*ast.GenDecl); ok {
			for _, spec := range gen.Specs {
				if typ, ok := spec.(*ast.TypeSpec); ok {
					if structure, ok := typ.Type.(*ast.StructType); ok {
						fields := map[string]bool{}
						for _, field := range structure.Fields.List {
							for _, name := range field.Names {
								fields[name.Name] = true
							}
						}
						fieldTypes[typ.Name.Name] = fields
					}
				}
			}
		}
	}
	bodyPlan := *plan
	bodyPlan.fieldTypes = fieldTypes
	memberNames := map[string]string{}
	for name, public := range plan.members {
		memberNames[name] = public
	}
	edits := map[int]string{}
	keyEdits := map[int]string{}
	if err := rewriteRuntimeCompositeKeys(file, &bodyPlan, fset, keyEdits); err != nil {
		return nil, err
	}
	for offset, replacement := range keyEdits {
		edits[offset-base] = replacement
	}
	usedImports := map[string]bool{}
	selectorSels := map[token.Pos]bool{}
	keyPositions := map[token.Pos]bool{}
	fieldPositions := map[token.Pos]bool{}
	ast.Inspect(file, func(node ast.Node) bool {
		switch value := node.(type) {
		case *ast.SelectorExpr:
			selectorSels[value.Sel.Pos()] = true
			if public := memberNames[value.Sel.Name]; public != "" {
				edits[fset.Position(value.Sel.Pos()).Offset-base] = public
			}
			if id, ok := value.X.(*ast.Ident); ok && imports[id.Name] != "" {
				usedImports[id.Name] = true
			}
		case *ast.KeyValueExpr:
			if key, ok := value.Key.(*ast.Ident); ok {
				keyPositions[key.Pos()] = true
			}
		case *ast.CompositeLit:
			if fields := fieldTypes[compositeTypeName(value.Type)]; fields != nil {
				for _, element := range value.Elts {
					item, ok := element.(*ast.KeyValueExpr)
					if !ok {
						continue
					}
					key, ok := item.Key.(*ast.Ident)
					if !ok || !fields[key.Name] {
						continue
					}
					if public := memberNames[key.Name]; public != "" {
						edits[fset.Position(key.Pos()).Offset-base] = public
					}
				}
			}
		}
		return true
	})
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			if d.Recv != nil {
				if public := memberNames[d.Name.Name]; public != "" {
					edits[fset.Position(d.Name.Pos()).Offset-base] = public
				}
			}
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				typ, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				if structure, ok := typ.Type.(*ast.StructType); ok {
					for _, field := range structure.Fields.List {
						for _, name := range field.Names {
							fieldPositions[name.Pos()] = true
							if public := memberNames[name.Name]; public != "" {
								edits[fset.Position(name.Pos()).Offset-base] = public
							}
						}
					}
				}
			}
		}
	}
	ast.Inspect(file, func(node ast.Node) bool {
		id, ok := node.(*ast.Ident)
		if !ok {
			return true
		}
		if id.Obj == nil && !selectorSels[id.Pos()] && !keyPositions[id.Pos()] && !fieldPositions[id.Pos()] {
			if public := plan.globals[id.Name]; public != "" {
				offset := fset.Position(id.Pos()).Offset - base
				edits[offset] = "tsruntime." + public
			}
		}
		if public := memberNames[id.Name]; public != "" && fieldPositions[id.Pos()] {
			edits[fset.Position(id.Pos()).Offset-base] = public
		}
		return true
	})
	result := applySourceEdits(body, edits)
	formatted, err := format.Source(append([]byte("package main\n"), result...))
	if err != nil {
		return nil, err
	}
	return bytes.TrimPrefix(formatted, []byte("package main\n")), nil
}

func packageMainSource(fullSource []byte, imports map[string]string, plan *runtimeExportPlan) ([]byte, error) {
	marker := []byte("// TSC_GOEMIT_RUNTIME_END")
	index := bytes.Index(fullSource, marker)
	if index < 0 {
		return nil, fmt.Errorf("generated source has no runtime marker")
	}
	body := fullSource[index+len(marker):]
	transformed, err := transformGeneratedBody(body, plan, imports)
	if err != nil {
		return nil, fmt.Errorf("transform generated Go body: %w", err)
	}
	wrapped := append([]byte("package main\n"), body...)
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "generated.go", wrapped, 0)
	if err != nil {
		return nil, err
	}
	used := map[string]bool{}
	ast.Inspect(file, func(node ast.Node) bool {
		if selector, ok := node.(*ast.SelectorExpr); ok {
			if id, ok := selector.X.(*ast.Ident); ok && imports[id.Name] != "" {
				used[id.Name] = true
			}
		}
		return true
	})
	var importLines []string
	importLines = append(importLines, "tsruntime \""+runtimePackageImport+"\"")
	names := make([]string, 0, len(imports))
	for name := range imports {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		importPath := imports[name]
		if used[name] {
			if name == path.Base(strings.Trim(importPath, "\"")) {
				importLines = append(importLines, importPath)
			} else {
				importLines = append(importLines, name+" "+importPath)
			}
		}
	}
	result := []byte("package main\n\nimport (\n\t" + strings.Join(importLines, "\n\t") + "\n)\n\n// TSC_GOEMIT_RUNTIME_END\n")
	result = append(result, transformed...)
	formatted, err := format.Source(result)
	if err == nil {
		return formatted, nil
	}
	line := 0
	_, _ = fmt.Sscanf(err.Error(), "generated.go:%d:", &line)
	if line == 0 {
		_, _ = fmt.Sscanf(err.Error(), "format: %d:", &line)
	}
	lines := bytes.Split(result, []byte("\n"))
	start, end := line-2, line+1
	if start < 1 {
		start = 1
	}
	if end > len(lines) {
		end = len(lines)
	}
	var excerpt strings.Builder
	for index := start; index <= end; index++ {
		fmt.Fprintf(&excerpt, "%d: %s\n", index, lines[index-1])
	}
	return nil, fmt.Errorf("%w\n%s", err, excerpt.String())
}

// RuntimePackageSource emits the prebuilt runtime package consumed by generated
// programs. It is called by cmd/generate-runtime and checked by goemit tests.
func RuntimePackageSource() ([]byte, error) {
	raw := []byte("package main\n" + runtimeTemplate() + "\n// TSC_GOEMIT_RUNTIME_END\n")
	formatted, err := formatValueSource(raw)
	if err != nil {
		return nil, err
	}
	exported, _, err := runtimeSourceAndPlan(formatted, "tsruntime")
	if err != nil {
		return nil, err
	}
	return append([]byte("// Code generated by cmd/generate-runtime; DO NOT EDIT.\n"), exported...), nil
}
