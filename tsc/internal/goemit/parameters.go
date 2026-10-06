package goemit

import (
	"fmt"

	"github.com/microsoft/TypeScript/tsc/internal/ast"
)

// Parameter defaults execute before their binding initializes. This preserves
// the temporal dead zone, while making earlier parameters available to defaults.
func (b *machineBuilder) initializeParameter(param *ast.Node, index int) {
	p := param.AsParameterDeclaration()
	if param.Name().Kind != ast.KindIdentifier && !ast.IsBindingPattern(param.Name()) {
		b.e.fail(param, "unsupported parameter binding")
		return
	}
	value := ""
	if p.DotDotDotToken != nil {
		value = b.temp(fmt.Sprintf("tsRestArgs(args,%d)", index))
	} else {
		value = b.temp(fmt.Sprintf("tsArg(args,%d)", index))
	}
	if p.Initializer != nil {
		if b.direct {
			b.emit("if tsIsUndefined(" + value + ") {")
			fallback := b.parameterDefault(param)
			b.emit(value + "=" + fallback)
			b.emit("}")
		} else {
			yes, end := b.block(), b.block()
			b.emit(fmt.Sprintf("if tsIsUndefined(%s){m.pc=%d}else{m.pc=%d};return", value, yes, end))
			b.current = yes
			fallback := b.parameterDefault(param)
			b.emit(value + "=" + fallback)
			b.jump(end)
			b.current = end
		}
	}
	if ast.IsBindingPattern(param.Name()) {
		b.bindPattern(param.Name(), value, true)
		return
	}
	if cell := b.e.bindings[param]; cell != nil {
		b.emit(cell.name + ".init(" + value + ")")
		b.validateParameter(param)
	}
}

// Constructor arguments use the declared native width before numeric literals
// can lose precision through the default JavaScript number representation.
func (b *machineBuilder) constructorArgument(call *ast.Node, class *classInfo, index int, arg *ast.Node) string {
	for class != nil {
		if class.constructor != nil {
			p := b.argumentPrimitive(class.constructor, index)
			if p.numeric() && (p.nulls == 0 || numericFieldInitializer(arg)) {
				return b.numericExpression(arg, p)
			}
			break
		}
		class = class.parent
	}
	if resolver, ok := b.e.resolver.(interface {
		GetEmitGenericCallTypes(*ast.Node, []*ast.Node) (*ast.Node, []string, []string, string, string)
	}); ok {
		_, _, parameters, _, _ := resolver.GetEmitGenericCallTypes(call, b.e.specializationCalls)
		if index < len(parameters) {
			p := goPrimitive(parameters[index])
			if p.numeric() {
				return b.numericExpression(arg, p)
			}
		}
	}

	return b.expression(arg)
}

func (b *machineBuilder) parameterDefault(param *ast.Node) string {
	initializer := param.AsParameterDeclaration().Initializer
	if cell := b.e.bindings[param]; cell != nil && cell.primitive.numeric() && cell.primitive.nulls == 0 {
		return b.numericExpression(initializer, cell.primitive)
	}
	return b.expression(initializer)
}

func (b *machineBuilder) argumentPrimitive(declaration *ast.Node, index int) primitive {
	params := declaration.Parameters()
	if len(params) == 0 {
		return primitive{}
	}
	if index >= len(params) {
		index = len(params) - 1
		if params[index].AsParameterDeclaration().DotDotDotToken == nil {
			return primitive{}
		}
	}
	param := params[index]
	if param.AsParameterDeclaration().DotDotDotToken != nil {
		typ := param.Type()
		if typ == nil {
			return primitive{}
		}
		if typ.Kind == ast.KindArrayType {
			return b.e.annotationPrimitive(typ.AsArrayTypeNode().ElementType)
		}
		if typ.Kind == ast.KindTypeReference {
			ref := typ.AsTypeReferenceNode()
			if ref.TypeName.Kind == ast.KindIdentifier && (ref.TypeName.Text() == "Array" || ref.TypeName.Text() == "ReadonlyArray") && ref.TypeArguments != nil && len(ref.TypeArguments.Nodes) == 1 {
				return b.e.annotationPrimitive(ref.TypeArguments.Nodes[0])
			}
		}
		return primitive{}
	}
	return b.e.primitive(param)
}
