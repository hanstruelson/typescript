package goemit

import (
	"go/ast"
	"go/token"
	"strings"
)

// threadLoopContext gives every JavaScript call an explicit execution context.
// Native adapters and Go host callbacks retain their own ABI and close over it.
func threadLoopContext(file *ast.File) {
	required := map[string]bool{"tsCall": true, "tsEngine": true, "call": true, "construct": true, "pull": true}
	functions := []*ast.FuncDecl{}
	addParam := func(fn *ast.FuncType) {
		if fn.Params == nil {
			fn.Params = &ast.FieldList{}
		}
		if len(fn.Params.List) > 0 && valueType(fn.Params.List[0].Type) == "*tsLoop" {
			return
		}
		names := []*ast.Ident{ast.NewIdent("loop")}
		if len(fn.Params.List) > 0 && len(fn.Params.List[0].Names) == 0 {
			names = nil
		}
		fn.Params.List = append([]*ast.Field{{Names: names, Type: &ast.StarExpr{X: ast.NewIdent("tsLoop")}}}, fn.Params.List...)
	}
	// Function-value callbacks must receive the caller's context, even when their
	// body happens to use only native arithmetic today.
	ast.Inspect(file, func(node ast.Node) bool {
		if field, ok := node.(*ast.Field); ok {
			for _, name := range field.Names {
				if strings.HasPrefix(name.Name, "Set") || strings.HasPrefix(name.Name, "Call") {
					if _, ok := field.Type.(*ast.FuncType); ok {
						required[name.Name] = true
					}
				}
			}
		}
		if fn, ok := node.(*ast.FuncDecl); ok {
			functions = append(functions, fn)
		}
		if call, ok := node.(*ast.CallExpr); ok && valueCallName(call.Fun) == "tsFunc" && len(call.Args) == 1 {
			if fn, ok := call.Args[0].(*ast.FuncLit); ok {
				addParam(fn.Type)
			} else {
				required[valueCallName(call.Args[0])] = true
			}
		}
		return true
	})
	callsRequired := func(body *ast.BlockStmt) bool {
		found := false
		ast.Inspect(body, func(node ast.Node) bool {
			if lit, ok := node.(*ast.FuncLit); ok && len(lit.Type.Params.List) > 0 && valueType(lit.Type.Params.List[0].Type) == "*tsLoop" {
				return false
			}
			if call, ok := node.(*ast.CallExpr); ok && required[valueCallName(call.Fun)] {
				found = true
			}
			return !found
		})
		return found
	}
	for changed := true; changed; {
		changed = false
		for _, fn := range functions {
			name := fn.Name.Name
			if name == "tsFunc" || name == "String" || name == "main" || required[name] {
				continue
			}
			if callsRequired(fn.Body) {
				required[name] = true
				changed = true
			}
		}
	}
	// A loop-owned receiver already supplies context. Other methods, including
	// class methods and numeric boundary cells, receive it explicitly.
	receiverContext := map[string]bool{}
	for _, fn := range functions {
		if !required[fn.Name.Name] {
			continue
		}
		if fn.Recv != nil {
			typ := valueType(fn.Recv.List[0].Type)
			receiver := fn.Recv.List[0].Names[0].Name
			var context ast.Expr
			switch typ {
			case "*tsLoop":
				context = ast.NewIdent(receiver)
			case "*tsMachine", "*tsPromise", "*tsModule":
				context = &ast.SelectorExpr{X: ast.NewIdent(receiver), Sel: ast.NewIdent("loop")}
			}
			if context != nil && !(fn.Type.Params != nil && len(fn.Type.Params.List) > 0 && valueType(fn.Type.Params.List[0].Type) == "*tsLoop") {
				receiverContext[fn.Name.Name] = true
				if receiver != "loop" {
					fn.Body.List = append([]ast.Stmt{&ast.AssignStmt{Lhs: []ast.Expr{ast.NewIdent("loop")}, Tok: token.DEFINE, Rhs: []ast.Expr{context}}}, fn.Body.List...)
					fn.Body.List = append([]ast.Stmt{fn.Body.List[0], &ast.AssignStmt{Lhs: []ast.Expr{ast.NewIdent("_")}, Tok: token.ASSIGN, Rhs: []ast.Expr{ast.NewIdent("loop")}}}, fn.Body.List[1:]...)
				}
				continue
			}
		}
		addParam(fn.Type)
	}
	// Constructor closures carry the executing loop, just like functions.
	ast.Inspect(file, func(node ast.Node) bool {
		if assign, ok := node.(*ast.AssignStmt); ok && len(assign.Lhs) == 1 && len(assign.Rhs) == 1 {
			if sel, ok := assign.Lhs[0].(*ast.SelectorExpr); ok && required[sel.Sel.Name] && !receiverContext[sel.Sel.Name] {
				if lit, ok := assign.Rhs[0].(*ast.FuncLit); ok {
					addParam(lit.Type)
				}
			}
		}
		return true
	})
	// Stored callbacks for class constructors, fields and views use the same ABI
	// as their implementations. Callback literals already naming loop are stable.
	ast.Inspect(file, func(node ast.Node) bool {
		if field, ok := node.(*ast.Field); ok {
			for _, name := range field.Names {
				if required[name.Name] && !receiverContext[name.Name] {
					if fn, ok := field.Type.(*ast.FuncType); ok {
						addParam(fn)
					}
				}
			}
		}
		if pair, ok := node.(*ast.KeyValueExpr); ok {
			if key, ok := pair.Key.(*ast.Ident); ok && required[key.Name] && !receiverContext[key.Name] {
				if lit, ok := pair.Value.(*ast.FuncLit); ok {
					addParam(lit.Type)
				}
			}
		}
		return true
	})
	ast.Inspect(file, func(node ast.Node) bool {
		if call, ok := node.(*ast.CallExpr); ok {
			name := valueCallName(call.Fun)
			_, selector := call.Fun.(*ast.SelectorExpr)
			if required[name] && !receiverContext[name] && name != "tsFunc" && (name != "call" || selector) {
				if len(call.Args) == 0 || valueType(call.Args[0]) != "loop" {
					call.Args = append([]ast.Expr{ast.NewIdent("loop")}, call.Args...)
				}
			}
		}
		return true
	})
	// Older class lowering captured the creation loop in each method. Explicit
	// parameters now carry the executing worker's loop.
	for _, fn := range functions {
		if fn.Recv == nil || !required[fn.Name.Name] {
			continue
		}
		list := fn.Body.List[:0]
		for _, stmt := range fn.Body.List {
			assign, ok := stmt.(*ast.AssignStmt)
			if ok && len(assign.Lhs) == 1 && len(assign.Rhs) == 1 && valueType(assign.Lhs[0]) == "loop" && strings.HasSuffix(valueType(assign.Rhs[0]), ".loop") && valueType(fn.Recv.List[0].Type) != "*tsMachine" && valueType(fn.Recv.List[0].Type) != "*tsPromise" && valueType(fn.Recv.List[0].Type) != "*tsModule" {
				continue
			}
			list = append(list, stmt)
		}
		fn.Body.List = list
	}
}
