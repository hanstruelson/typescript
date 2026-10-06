package goemit

import (
	"fmt"
	"strings"

	"github.com/microsoft/TypeScript/tsc/internal/ast"
)

type nativeFunction struct {
	name       string
	parameters []primitive
	result     primitive
}

func (e *emitter) planNativeFunctions() {
	e.nativeFunctions = map[*ast.Node]*nativeFunction{}
	for _, node := range e.file.Statements.Nodes {
		if node.Kind != ast.KindFunctionDeclaration || node.Body() == nil || len(node.TypeParameters()) != 0 || !directClassBody(node) {
			continue
		}
		result := e.annotationPrimitive(node.Type())
		if result.kind == "" || result.nulls != 0 {
			continue
		}
		fn := &nativeFunction{name: e.unique("Native"), result: result}
		valid := true
		for _, param := range node.Parameters() {
			p := e.primitive(param)
			decl := param.AsParameterDeclaration()
			if p.kind == "" || p.nulls != 0 || decl.Initializer != nil || decl.DotDotDotToken != nil || param.QuestionToken() != nil || param.Name().Kind != ast.KindIdentifier {
				valid = false
				break
			}
			fn.parameters = append(fn.parameters, p)
		}
		if valid {
			e.nativeFunctions[node] = fn
		}
	}
	// A reassigned function must continue to dispatch through its binding cell.
	var visit func(*ast.Node)
	visit = func(node *ast.Node) {
		if node.Kind == ast.KindBinaryExpression {
			n := node.AsBinaryExpression()
			if ast.IsAssignmentOperator(n.OperatorToken.Kind) && n.Left.Kind == ast.KindIdentifier {
				delete(e.nativeFunctions, e.reference(n.Left))
			}
		}
		node.ForEachChild(func(child *ast.Node) bool { visit(child); return false })
	}
	visit(e.file.AsNode())
}
func (f *nativeFunction) signature() string {
	parts := []string{"*tsLoop"}
	for _, p := range f.parameters {
		parts = append(parts, p.goType())
	}
	return "func(" + strings.Join(parts, ",") + ") " + f.result.goType()
}
func (b *machineBuilder) nativeFunction(node *ast.Node, fn *nativeFunction) string {
	child := b.e.newMachine(node, b, false)
	child.direct = true
	child.nativeReturn = fn.result.goType()
	parameters := []string{"loop *tsLoop"}
	arguments := []string{"loop"}
	for i, param := range node.Parameters() {
		name := fmt.Sprintf("nativeArg%d", i)
		parameters = append(parameters, name+" "+fn.parameters[i].goType())
		cell := b.e.bindings[param]
		child.emit(cell.name + ".initNative(" + name + ")")
		arguments = append(arguments, fmt.Sprintf("tsNative[%s](tsBoundary(tsArg(args,%d),%q,0,%t))", fn.parameters[i].goType(), i, fn.parameters[i].kind, b.e.coerce))
	}
	child.directStatements(node.Body().AsBlock().Statements.Nodes)
	b.temps = append(b.temps, fn.name)
	if b.tempTypes == nil {
		b.tempTypes = map[string]string{}
	}
	b.tempTypes[fn.name] = fn.signature()
	b.emit(fn.name + "=func(" + strings.Join(parameters, ",") + ") " + fn.result.goType() + " {\n" + child.finishDirect() + "\n}")
	return "tsFunc(func(args ...tsValue)tsValue{return " + fn.name + "(" + strings.Join(arguments, ",") + ")})"
}
func (b *machineBuilder) nativeCall(node *ast.Node, fn *nativeFunction) string {
	call := node.AsCallExpression()
	arguments := []string{"loop"}
	values := []string{}
	for i, arg := range call.Arguments.Nodes {
		if i < len(fn.parameters) && fn.parameters[i].numeric() {
			values = append(values, b.numericExpression(arg, fn.parameters[i]))
		} else {
			values = append(values, b.expression(arg))
		}
	}
	for i, param := range fn.parameters {
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
