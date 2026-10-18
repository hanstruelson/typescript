package regexpengine

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"strconv"
	"unicode/utf16"
)

// CompileUTF16ToGoMatcherLiteral emits an expression that evaluates to a
// standalone UTF-16 matcher closure. Matcher tables are placed in a factory
// closure so the result can be compiled as a runtime statement without
// package-level declarations or a temporary source file.
func CompileUTF16ToGoMatcherLiteral(source []uint16, flagText string) (expression string, captureSlots int, err error) {
	const functionName = "tsJITDynamicRegex"
	generated, err := CompileUTF16ToGo(source, flagText, functionName)
	if err != nil {
		return "", 0, err
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "matcher.go", "package jitmatcher\n"+generated, parser.AllErrors)
	if err != nil {
		return "", 0, fmt.Errorf("parse generated Go matcher: %w", err)
	}

	entryName := functionName + "Into"
	var entry *ast.FuncDecl
	var captureSlotsName string = functionName + "CaptureSlots"
	var captureSlotsValue int
	// Generated matcher tables are immutable package-level var/const specs.
	// The selected entry may refer to only some of them, so retain exactly the
	// declarations reachable by identifier from that function.
	type candidate struct {
		statement *ast.DeclStmt
		names     []string
	}
	var candidates []candidate
	byName := make(map[string]int)
	for _, declaration := range file.Decls {
		switch declaration := declaration.(type) {
		case *ast.FuncDecl:
			if declaration.Name.Name == entryName {
				entry = declaration
			}
		case *ast.GenDecl:
			if declaration.Tok != token.VAR && declaration.Tok != token.CONST {
				continue
			}
			for _, rawSpec := range declaration.Specs {
				spec, ok := rawSpec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				single := *declaration
				single.Specs = []ast.Spec{spec}
				item := candidate{statement: &ast.DeclStmt{Decl: &single}}
				for _, name := range spec.Names {
					item.names = append(item.names, name.Name)
					if name.Name == captureSlotsName && len(spec.Values) == 1 {
						if literal, ok := spec.Values[0].(*ast.BasicLit); ok && literal.Kind == token.INT {
							captureSlotsValue, _ = strconv.Atoi(literal.Value)
						}
					}
				}
				index := len(candidates)
				candidates = append(candidates, item)
				for _, name := range item.names {
					byName[name] = index
				}
			}
		}
	}
	if entry == nil || entry.Body == nil {
		return "", 0, fmt.Errorf("generated Go matcher has no %s function", entryName)
	}
	if captureSlotsValue <= 0 {
		return "", 0, fmt.Errorf("generated Go matcher has no valid %s constant", captureSlotsName)
	}

	needed := make(map[int]bool)
	var collectNames func(ast.Node)
	collectNames = func(node ast.Node) {
		ast.Inspect(node, func(node ast.Node) bool {
			identifier, ok := node.(*ast.Ident)
			if !ok {
				return true
			}
			if index, ok := byName[identifier.Name]; ok && !needed[index] {
				needed[index] = true
				collectNames(candidates[index].statement)
			}
			return true
		})
	}
	collectNames(entry.Body)

	statements := make([]ast.Stmt, 0, len(candidates)+1)
	for index, item := range candidates {
		if needed[index] {
			statements = append(statements, item.statement)
		}
	}
	matcher := &ast.FuncLit{Type: entry.Type, Body: entry.Body}
	statements = append(statements, &ast.ReturnStmt{Results: []ast.Expr{matcher}})
	factory := &ast.FuncLit{
		Type: &ast.FuncType{
			Params:  &ast.FieldList{},
			Results: &ast.FieldList{List: []*ast.Field{{Type: entry.Type}}},
		},
		Body: &ast.BlockStmt{List: statements},
	}
	call := &ast.CallExpr{Fun: factory}
	var output bytes.Buffer
	if err := format.Node(&output, fset, call); err != nil {
		return "", 0, fmt.Errorf("format generated matcher closure: %w", err)
	}
	return output.String(), captureSlotsValue, nil
}

// CompileToGoMatcherLiteral is the UTF-8 convenience form of
// CompileUTF16ToGoMatcherLiteral.
func CompileToGoMatcherLiteral(source, flagText string) (expression string, captureSlots int, err error) {
	return CompileUTF16ToGoMatcherLiteral(utf16.Encode([]rune(source)), flagText)
}
