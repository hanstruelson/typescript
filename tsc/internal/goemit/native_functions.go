package goemit

import (
	"fmt"
	"strings"

	"github.com/microsoft/TypeScript/tsc/internal/ast"
)

type nativeFunction struct {
	name            string
	node            *ast.Node
	parameters      []primitive
	arrayParameters []primitive
	result          primitive
}

func (e *emitter) planNativeFunctions() {
	e.nativeFunctions = map[*ast.Node]*nativeFunction{}
	e.stableCallbacks = map[*ast.Node]*ast.Node{}
	for _, node := range e.file.Statements.Nodes {
		if node.Kind != ast.KindFunctionDeclaration || node.Body() == nil || len(node.TypeParameters()) != 0 || !nativeCallbackBody(node) {
			continue
		}
		result := e.primitiveFunctionReturn(node)
		if result.kind == "" && result.nulls != 2 || result.kind != "" && result.nulls != 0 {
			continue
		}
		fn := &nativeFunction{name: e.unique("Native"), node: node, result: result}
		valid := true
		for _, param := range node.Parameters() {
			p := e.primitive(param)
			a := e.arrayElementPrimitive(param)
			decl := param.AsParameterDeclaration()
			if (p.kind == "" && (a.kind == "")) || p.nulls != 0 || decl.Initializer != nil || decl.DotDotDotToken != nil || param.QuestionToken() != nil || param.Name().Kind != ast.KindIdentifier {
				valid = false
				break
			}
			fn.parameters = append(fn.parameters, p)
			fn.arrayParameters = append(fn.arrayParameters, a)
		}
		if valid {
			e.nativeFunctions[node] = fn
		}
	}

	// Immutable function-valued bindings retain their concrete callback ABI.
	var callbacks func(*ast.Node)
	callbacks = func(node *ast.Node) {
		if node.Kind == ast.KindFunctionDeclaration && node.Body() != nil && len(node.TypeParameters()) == 0 && nativeCallbackBody(node) {
			e.stableCallbacks[node] = node
		}
		if node.Kind == ast.KindVariableDeclaration {
			cell := e.bindings[node]
			init := node.AsVariableDeclaration().Initializer
			if cell != nil && nativeStableCallbackBinding(node, cell.owner) && cell.constant && init != nil && (init.Kind == ast.KindArrowFunction || init.Kind == ast.KindFunctionExpression) && init.Name() == nil && nativeCallbackBody(init) && len(init.TypeParameters()) == 0 {
				e.stableCallbacks[node] = init
				result := e.primitiveFunctionReturn(init)
				valid := result.kind != "" && result.nulls == 0 || result.kind == "" && result.nulls == 2
				fn := &nativeFunction{name: e.unique("NativeCallback"), node: init, result: result}
				for _, param := range init.Parameters() {
					p := e.primitive(param)
					a := e.arrayElementPrimitive(param)
					decl := param.AsParameterDeclaration()
					if p.kind == "" && a.kind == "" || p.nulls != 0 || decl.Initializer != nil || decl.DotDotDotToken != nil || param.QuestionToken() != nil || param.Name().Kind != ast.KindIdentifier {
						valid = false
					}
					fn.parameters = append(fn.parameters, p)
					fn.arrayParameters = append(fn.arrayParameters, a)
				}
				if valid {
					e.nativeFunctions[init] = fn
					e.nativeFunctions[node] = fn
				}
			}
		}
		node.ForEachChild(func(child *ast.Node) bool { callbacks(child); return false })
	}
	callbacks(e.file.AsNode())
	// A reassigned function must continue to dispatch through its binding cell.
	var visit func(*ast.Node)
	visit = func(node *ast.Node) {
		if node.Kind == ast.KindBinaryExpression {
			n := node.AsBinaryExpression()
			if ast.IsAssignmentOperator(n.OperatorToken.Kind) && n.Left.Kind == ast.KindIdentifier {
				delete(e.nativeFunctions, e.reference(n.Left))
				delete(e.stableCallbacks, e.reference(n.Left))
			}
		}
		node.ForEachChild(func(child *ast.Node) bool { visit(child); return false })
	}
	visit(e.file.AsNode())
}
func (f *nativeFunction) signature() string {
	parts := []string{"*tsLoop"}
	for i, p := range f.parameters {
		if a := f.arrayParameters[i]; a.kind != "" {
			parts = append(parts, "*tsGrowableStorage["+a.goType()+"]")
			continue
		}
		parts = append(parts, p.goType())
	}
	return "func(" + strings.Join(parts, ",") + ") " + f.result.goType()
}
func (b *machineBuilder) nativeFunction(node *ast.Node, fn *nativeFunction) string {
	child := b.e.newMachine(node, b, false)
	child.direct = directClassBody(node)
	child.nativeReturn = fn.result.goType()
	child.nativeResult = fn.result
	if !child.direct {
		child.nativeOutput = child.e.unique("NativeReturn")
		child.temps = append(child.temps, child.nativeOutput)
		child.tempTypes = map[string]string{child.nativeOutput: fn.result.elementGoType()}
	}
	if node.Kind == ast.KindArrowFunction {
		child.concrete = b.concrete
		child.declaring = b.declaring
		child.receiver = b.receiver
	}
	parameters := []string{"loop *tsLoop"}
	arguments := []string{"loop"}
	for i, param := range node.Parameters() {
		name := fmt.Sprintf("nativeArg%d", i)
		cell := b.e.bindings[param]
		if a := fn.arrayParameters[i]; a.kind != "" {
			parameters = append(parameters, name+" *tsGrowableStorage["+a.goType()+"]")
			child.emit(cell.name + ".initStorage(" + name + ")")
			arguments = append(arguments, fmt.Sprintf("tsDenseStorage[%s](tsArrayBoundary(tsArg(args,%d),%q,%d,%t),%q,%d)", a.goType(), i, a.kind, a.nulls, b.e.coerce, a.kind, a.nulls))
			continue
		}
		parameters = append(parameters, name+" "+fn.parameters[i].goType())
		child.emit(cell.name + ".initNative(" + name + ")")
		arguments = append(arguments, fmt.Sprintf("tsNative[%s](tsBoundary(tsArg(args,%d),%q,0,%t))", fn.parameters[i].goType(), i, fn.parameters[i].kind, b.e.coerce))
	}
	if node.Body().Kind == ast.KindBlock {
		if child.direct {
			child.directStatements(node.Body().AsBlock().Statements.Nodes)
		} else {
			child.statements(node.Body().AsBlock().Statements.Nodes)
			child.endNativeMachine()
		}
	} else {
		value := child.denseElement(node.Body(), fn.result)
		child.emit("return " + value)
	}
	b.temps = append(b.temps, fn.name)
	if b.tempTypes == nil {
		b.tempTypes = map[string]string{}
	}
	b.tempTypes[fn.name] = fn.signature()
	body := child.finishDirect()
	if !child.direct {
		body = child.finishNativeMachine()
	}
	b.emit(fn.name + "=func(" + strings.Join(parameters, ",") + ") " + fn.result.goType() + " {\n" + body + "\n}")
	return fmt.Sprintf("func()*tsFunction{function:=tsFunc(func(args ...tsValue)tsValue{return %s(%s)});function.constructible=%t;function.name=%q;function.length=%d;return function}()", fn.name, strings.Join(arguments, ","), node.Kind == ast.KindFunctionDeclaration || node.Kind == ast.KindFunctionExpression, functionName(node), functionLength(node))
}
func (b *machineBuilder) nativeCall(node *ast.Node, fn *nativeFunction) string {
	call := node.AsCallExpression()
	arguments := []string{"loop"}
	values := []string{}
	for i, arg := range call.Arguments.Nodes {
		if i < len(fn.parameters) && fn.arrayParameters[i].kind != "" {
			a := fn.arrayParameters[i]
			if source := b.e.arrayElementPrimitive(arg); source.kind == a.kind && source.nulls == a.nulls {
				values = append(values, b.denseStorage(arg, a))
			} else {
				value := b.expression(arg)
				values = append(values, b.typedTemp(fmt.Sprintf("tsDenseStorage[%s](tsArrayBoundary(%s,%q,%d,%t),%q,%d)", a.goType(), value, a.kind, a.nulls, b.e.coerce, a.kind, a.nulls), "*tsGrowableStorage["+a.goType()+"]"))
			}
		} else if i < len(fn.parameters) && fn.parameters[i].numeric() {
			values = append(values, b.numericExpression(arg, fn.parameters[i]))
		} else {
			values = append(values, b.expression(arg))
		}
	}
	for i, param := range fn.parameters {
		if a := fn.arrayParameters[i]; a.kind != "" {
			if i < len(values) {
				arguments = append(arguments, values[i])
			} else {
				arguments = append(arguments, fmt.Sprintf("tsDenseStorage[%s](tsArrayBoundary(tsU,%q,%d,%t),%q,%d)", a.goType(), a.kind, a.nulls, b.e.coerce, a.kind, a.nulls))
			}
			continue
		}
		value := "tsU"
		if i < len(values) {
			value = values[i]
		}
		if b.tempType(value) != param.goType() {
			value = fmt.Sprintf("tsNative[%s](tsBoundary(%s,%q,0,%t))", param.goType(), value, param.kind, b.e.coerce)
		}
		arguments = append(arguments, value)
	}
	return b.typedTemp(fn.name+"("+strings.Join(arguments, ",")+")", fn.result.goType())
}

func nativeStableCallbackBinding(node, owner *ast.Node) bool {
	for current := node.Parent; current != nil && current != owner; current = current.Parent {
		switch current.Kind {
		case ast.KindForStatement, ast.KindForOfStatement, ast.KindForInStatement, ast.KindWhileStatement, ast.KindDoStatement:
			return false
		}
	}
	return true
}
