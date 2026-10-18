package goemit

import (
	"go/ast"
	"go/token"
	"strings"
)

func loopTypeKey(name string) string {
	name = strings.TrimPrefix(strings.Trim(name, "()"), "*")
	if at := strings.IndexByte(name, '['); at >= 0 {
		name = name[:at]
	}
	return name
}

func loopFunctionResult(fn *ast.FuncType) string {
	if fn == nil || fn.Results == nil || len(fn.Results.List) == 0 {
		return ""
	}
	return valueType(fn.Results.List[0].Type)
}

func loopExpressionType(expr ast.Expr, fields map[string]map[string]string, functions map[string]string, methods map[string]map[string]string) string {
	if expr == nil {
		return ""
	}
	switch value := expr.(type) {
	case *ast.Ident:
		if value.Obj == nil {
			return ""
		}
		switch declaration := value.Obj.Decl.(type) {
		case *ast.Field:
			return valueType(declaration.Type)
		case *ast.ValueSpec:
			if declaration.Type != nil {
				return valueType(declaration.Type)
			}
			for index, name := range declaration.Names {
				if name == value && index < len(declaration.Values) {
					return loopExpressionType(declaration.Values[index], fields, functions, methods)
				}
			}
		case *ast.AssignStmt:
			for index, left := range declaration.Lhs {
				if name, ok := left.(*ast.Ident); ok && name == value && index < len(declaration.Rhs) {
					return loopExpressionType(declaration.Rhs[index], fields, functions, methods)
				}
			}
		case *ast.FuncDecl:
			return loopFunctionResult(declaration.Type)
		}
	case *ast.SelectorExpr:
		base := loopTypeKey(loopExpressionType(value.X, fields, functions, methods))
		if fieldType := fields[base][value.Sel.Name]; fieldType != "" {
			return fieldType
		}
		if methodType := methods[base][value.Sel.Name]; methodType != "" {
			return methodType
		}
	case *ast.CallExpr:
		if selector, ok := value.Fun.(*ast.SelectorExpr); ok {
			base := loopTypeKey(loopExpressionType(selector.X, fields, functions, methods))
			return methods[base][selector.Sel.Name]
		}
		typeName := strings.Trim(valueType(value.Fun), "()")
		base := loopTypeKey(typeName)
		if fields[base] != nil || methods[base] != nil {
			return typeName
		}
		return functions[valueCallName(value.Fun)]
	case *ast.ParenExpr:
		return loopExpressionType(value.X, fields, functions, methods)
	case *ast.StarExpr:
		return "*" + loopExpressionType(value.X, fields, functions, methods)
	case *ast.UnaryExpr:
		if value.Op == token.AND {
			return "*" + loopExpressionType(value.X, fields, functions, methods)
		}
		return loopExpressionType(value.X, fields, functions, methods)
	case *ast.CompositeLit:
		return valueType(value.Type)
	case *ast.TypeAssertExpr:
		return valueType(value.Type)
	}
	return ""
}

func insertLoopArgument(arguments []ast.Expr, position int) []ast.Expr {
	if position < 0 {
		position = 0
	}
	if position > len(arguments) {
		position = len(arguments)
	}
	arguments = append(arguments, nil)
	copy(arguments[position+1:], arguments[position:])
	arguments[position] = ast.NewIdent("loop")
	return arguments
}

// threadLoopContext gives every JavaScript call an explicit execution context.
// Native adapters and Go host callbacks retain their own ABI and close over it.
func threadLoopContext(file *ast.File) {
	// Keep the runtime ABI independent of the generated program. These roots are
	// the loop-aware entry points shared by every emitted program. Console.log is
	// also passed as a named callback by generated code, so include it here even
	// though the runtime template itself does not take its address.
	required := map[string]bool{"tsCall": true, "tsEngine": true, "call": true, "construct": true, "pull": true, "growPut": true, "growSplice": true, "writeField": true, "tsConsoleLog": true}
	requiredMethods := map[string]bool{}
	functions := []*ast.FuncDecl{}
	methodTypes := map[string]map[string]string{}
	methodLoopPositions := map[string]map[string][]int{}
	structFields := map[string]map[string]string{}
	functionResults := map[string]string{}
	for _, declaration := range file.Decls {
		switch value := declaration.(type) {
		case *ast.FuncDecl:
			functionResults[value.Name.Name] = loopFunctionResult(value.Type)
		case *ast.GenDecl:
			for _, spec := range value.Specs {
				typ, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				structure, ok := typ.Type.(*ast.StructType)
				if !ok {
					continue
				}
				fields := map[string]string{}
				for _, field := range structure.Fields.List {
					for _, name := range field.Names {
						fields[name.Name] = valueType(field.Type)
					}
				}
				structFields[typ.Name.Name] = fields
			}
		}
	}
	localFunctionValues := map[string]bool{}
	loopParameterPositions := map[string][]int{}
	loopIdentifiers := map[string]bool{"loop": true}
	// Keep call-site insertion idempotent by checking for an explicit loop
	// argument on each call. A name-only "already explicit" set is incorrect
	// for overloaded methods such as `get`, where unrelated receiver methods
	// can have different loop signatures.
	addParam := func(fn *ast.FuncType) {
		if fn.Params == nil {
			fn.Params = &ast.FieldList{}
		}
		for _, parameter := range fn.Params.List {
			if valueType(parameter.Type) == "*tsLoop" {
				return
			}
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
		// Calls through local function values must retain their declared
		// signature. A name-only lookup here can mistake a local callback such
		// as `get := func(value tsValue) ...` for an unrelated loop-aware helper
		// with the same name and inject an invalid execution-context argument.
		switch declaration := node.(type) {
		case *ast.AssignStmt:
			for index, left := range declaration.Lhs {
				if index >= len(declaration.Rhs) {
					continue
				}
				if name, ok := left.(*ast.Ident); ok {
					if _, ok := declaration.Rhs[index].(*ast.FuncLit); ok {
						localFunctionValues[name.Name] = true
					}
				}
			}
		case *ast.ValueSpec:
			for index, name := range declaration.Names {
				if index < len(declaration.Values) {
					if _, ok := declaration.Values[index].(*ast.FuncLit); ok {
						localFunctionValues[name.Name] = true
					}
				}
			}
		}
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
			if fn.Recv != nil {
				receiverType := loopTypeKey(valueType(fn.Recv.List[0].Type))
				if methodTypes[receiverType] == nil {
					methodTypes[receiverType] = map[string]string{}
					methodLoopPositions[receiverType] = map[string][]int{}
				}
				methodTypes[receiverType][fn.Name.Name] = loopFunctionResult(fn.Type)
			}
			if fn.Recv != nil && valueType(fn.Recv.List[0].Type) == "*tsLoop" {
				for _, name := range fn.Recv.List[0].Names {
					loopIdentifiers[name.Name] = true
				}
			}
			position := 0
			for _, parameter := range fn.Type.Params.List {
				if valueType(parameter.Type) == "*tsLoop" {
					for _, name := range parameter.Names {
						loopIdentifiers[name.Name] = true
					}
					if fn.Recv != nil {
						receiverType := loopTypeKey(valueType(fn.Recv.List[0].Type))
						methodLoopPositions[receiverType][fn.Name.Name] = append(methodLoopPositions[receiverType][fn.Name.Name], position)
						requiredMethods[receiverType+"."+fn.Name.Name] = true
					} else {
						loopParameterPositions[fn.Name.Name] = append(loopParameterPositions[fn.Name.Name], position)
						required[fn.Name.Name] = true
					}
					break
				}
				count := len(parameter.Names)
				if count == 0 {
					count = 1
				}
				position += count
			}
		}
		if call, ok := node.(*ast.CallExpr); ok && valueCallName(call.Fun) == "tsFunc" && len(call.Args) == 1 {
			if fn, ok := call.Args[0].(*ast.FuncLit); ok {
				addParam(fn.Type)
			} else {
				if selector, ok := call.Args[0].(*ast.SelectorExpr); ok {
					receiverType := loopTypeKey(loopExpressionType(selector.X, structFields, functionResults, methodTypes))
					requiredMethods[receiverType+"."+selector.Sel.Name] = true
				} else {
					required[valueCallName(call.Args[0])] = true
				}
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
			if call, ok := node.(*ast.CallExpr); ok {
				if selector, ok := call.Fun.(*ast.SelectorExpr); ok {
					receiverType := loopTypeKey(loopExpressionType(selector.X, structFields, functionResults, methodTypes))
					if requiredMethods[receiverType+"."+selector.Sel.Name] || (strings.HasPrefix(selector.Sel.Name, "Set") || strings.HasPrefix(selector.Sel.Name, "Call")) && required[selector.Sel.Name] {
						found = true
					}
				} else if required[valueCallName(call.Fun)] {
					found = true
				}
			}
			return !found
		})
		return found
	}
	for changed := true; changed; {
		changed = false
		for _, fn := range functions {
			name := fn.Name.Name
			key := ""
			isRequired := required[name]
			if fn.Recv != nil {
				key = loopTypeKey(valueType(fn.Recv.List[0].Type)) + "." + name
				isRequired = requiredMethods[key]
			}
			if name == "tsFunc" || name == "String" || name == "main" || isRequired {
				continue
			}
			if callsRequired(fn.Body) {
				if key != "" {
					requiredMethods[key] = true
				} else {
					required[name] = true
				}
				changed = true
			}
		}
	}
	// A loop-owned receiver already supplies context. Other methods, including
	// class methods and numeric boundary cells, receive it explicitly.
	receiverContext := map[string]bool{}
	for _, fn := range functions {
		name := fn.Name.Name
		key := name
		isRequired := required[name]
		if fn.Recv != nil {
			receiverType := loopTypeKey(valueType(fn.Recv.List[0].Type))
			key = receiverType + "." + name
			isRequired = requiredMethods[key]
		}
		if !isRequired {
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
			default:
				if structFields[loopTypeKey(typ)]["loop"] == "*tsLoop" {
					context = &ast.SelectorExpr{X: ast.NewIdent(receiver), Sel: ast.NewIdent("loop")}
				}
			}
			hasLoopParameter := false
			if fn.Type.Params != nil {
				for _, parameter := range fn.Type.Params.List {
					if valueType(parameter.Type) == "*tsLoop" {
						hasLoopParameter = true
						break
					}
				}
			}
			if context != nil && !hasLoopParameter {
				receiverContext[key] = true
				if receiver != "loop" {
					fn.Body.List = append([]ast.Stmt{&ast.AssignStmt{Lhs: []ast.Expr{ast.NewIdent("loop")}, Tok: token.DEFINE, Rhs: []ast.Expr{context}}}, fn.Body.List...)
					fn.Body.List = append([]ast.Stmt{fn.Body.List[0], &ast.AssignStmt{Lhs: []ast.Expr{ast.NewIdent("_")}, Tok: token.ASSIGN, Rhs: []ast.Expr{ast.NewIdent("loop")}}}, fn.Body.List[1:]...)
				}
				continue
			}
		}
		hasLoopParameter := false
		if fn.Recv != nil {
			hasLoopParameter = len(methodLoopPositions[loopTypeKey(valueType(fn.Recv.List[0].Type))][fn.Name.Name]) > 0
		} else {
			hasLoopParameter = len(loopParameterPositions[fn.Name.Name]) > 0
		}
		addParam(fn.Type)
		if !hasLoopParameter && fn.Recv != nil {
			receiverType := loopTypeKey(valueType(fn.Recv.List[0].Type))
			methodLoopPositions[receiverType][fn.Name.Name] = []int{0}
		} else if !hasLoopParameter {
			loopParameterPositions[fn.Name.Name] = []int{0}
		}
	}
	isLoopArgument := func(expr ast.Expr) bool {
		switch value := expr.(type) {
		case *ast.Ident:
			return loopIdentifiers[value.Name]
		case *ast.SelectorExpr:
			return value.Sel.Name == "loop"
		}
		return false
	}
	goBuiltinCalls := map[string]bool{
		"append": true, "cap": true, "clear": true, "close": true, "complex": true,
		"copy": true, "delete": true, "imag": true, "len": true, "make": true,
		"max": true, "min": true, "new": true, "panic": true, "print": true,
		"println": true, "real": true, "recover": true,
	}
	// Constructor closures carry the executing loop, just like functions.
	ast.Inspect(file, func(node ast.Node) bool {
		if assign, ok := node.(*ast.AssignStmt); ok && len(assign.Lhs) == 1 && len(assign.Rhs) == 1 {
			if sel, ok := assign.Lhs[0].(*ast.SelectorExpr); ok && required[sel.Sel.Name] && (strings.HasPrefix(sel.Sel.Name, "Set") || strings.HasPrefix(sel.Sel.Name, "Call")) {
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
				if required[name.Name] && (strings.HasPrefix(name.Name, "Set") || strings.HasPrefix(name.Name, "Call")) {
					if fn, ok := field.Type.(*ast.FuncType); ok {
						addParam(fn)
					}
				}
			}
		}
		if pair, ok := node.(*ast.KeyValueExpr); ok {
			if key, ok := pair.Key.(*ast.Ident); ok && required[key.Name] && (strings.HasPrefix(key.Name, "Set") || strings.HasPrefix(key.Name, "Call")) {
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
			selector, isSelector := call.Fun.(*ast.SelectorExpr)
			if !isSelector && goBuiltinCalls[name] {
				return true
			}
			if isSelector {
				// Generated binding cells expose a zero-argument Get method, while
				// value objects expose loop-aware RuntimeFieldGet(key). The imported
				// receiver type is unknown here, so preserve the zero-argument cell ABI.
				if (selector.Sel.Name == "get" || selector.Sel.Name == "RuntimeFieldGet") && (len(call.Args) == 0 || len(call.Args) == 1 && isLoopArgument(call.Args[0])) {
					return true
				}
				receiverType := loopTypeKey(loopExpressionType(selector.X, structFields, functionResults, methodTypes))
				if receiverType == "" {
					candidate := ""
					ambiguous := false
					for typeName, members := range methodTypes {
						if _, exists := members[name]; !exists || !requiredMethods[typeName+"."+name] {
							continue
						}
						if candidate != "" {
							ambiguous = true
							break
						}
						candidate = typeName
					}
					if !ambiguous {
						receiverType = candidate
					}
				}
				key := receiverType + "." + name
				fieldCall := (strings.HasPrefix(name, "Set") || strings.HasPrefix(name, "Call")) && required[name]
				if !requiredMethods[key] && !fieldCall || receiverContext[key] {
					return true
				}
				explicitLoop := false
				positions := methodLoopPositions[receiverType][name]
				position := 0
				if len(positions) > 0 {
					position = positions[0]
				}
				for _, position := range positions {
					if position < len(call.Args) && (isLoopArgument(call.Args[position]) || valueType(call.Args[position]) == "nil") {
						explicitLoop = true
						break
					}
				}
				if len(positions) == 0 && len(call.Args) > 0 && isLoopArgument(call.Args[0]) {
					explicitLoop = true
				}
				if !explicitLoop {
					call.Args = insertLoopArgument(call.Args, position)
				}
				return true
			}
			if required[name] && !localFunctionValues[name] && name != "tsFunc" && name != "call" {
				explicitLoop := false
				positions := loopParameterPositions[name]
				position := 0
				if len(positions) > 0 {
					position = positions[0]
				}
				for _, position := range positions {
					if position < len(call.Args) && (isLoopArgument(call.Args[position]) || valueType(call.Args[position]) == "nil") {
						explicitLoop = true
						break
					}
				}
				if !explicitLoop && (len(positions) > 0 || len(call.Args) == 0 || !isLoopArgument(call.Args[0])) {
					call.Args = insertLoopArgument(call.Args, position)
				}
			}
		}
		return true
	})
	// Older class lowering captured the creation loop in each method. Explicit
	// parameters now carry the executing worker's loop.
	for _, fn := range functions {
		if fn.Recv == nil || !requiredMethods[loopTypeKey(valueType(fn.Recv.List[0].Type))+"."+fn.Name.Name] {
			continue
		}
		contextKey := loopTypeKey(valueType(fn.Recv.List[0].Type)) + "." + fn.Name.Name
		list := fn.Body.List[:0]
		keptContextAlias := false
		for _, stmt := range fn.Body.List {
			assign, ok := stmt.(*ast.AssignStmt)
			if ok && len(assign.Lhs) == 1 && len(assign.Rhs) == 1 && valueType(assign.Lhs[0]) == "loop" && strings.HasSuffix(valueType(assign.Rhs[0]), ".loop") && valueType(fn.Recv.List[0].Type) != "*tsMachine" && valueType(fn.Recv.List[0].Type) != "*tsPromise" && valueType(fn.Recv.List[0].Type) != "*tsModule" {
				if receiverContext[contextKey] && !keptContextAlias {
					keptContextAlias = true
				} else {
					continue
				}
			}
			list = append(list, stmt)
		}
		fn.Body.List = list
	}
}
