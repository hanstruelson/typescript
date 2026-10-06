package goemit

import (
	"github.com/microsoft/TypeScript/tsc/internal/ast"
	"strconv"
	"strings"
)

func (b *machineBuilder) classDecorators(node *ast.Node) []string {
	decorators := []string{}
	if modifiers := node.Modifiers(); modifiers != nil {
		for _, modifier := range modifiers.Nodes {
			if modifier.Kind == ast.KindDecorator {
				decorators = append(decorators, b.expression(modifier.AsDecorator().Expression))
			}
		}
	}
	return decorators
}
func (b *machineBuilder) applyClassDecorators(node *ast.Node, values []string, class string) string {
	if len(values) == 0 {
		return ""
	}
	name := ""
	if node.Name() != nil {
		name = node.Name().Text()
	}
	result := b.typedTemp("tsDecorateClass("+class+",[]tsValue{"+strings.Join(values, ",")+"},"+strconv.Quote(name)+","+strconv.FormatBool(b.e.legacyDecorators)+")", "tsClassDecoration")
	return result
}

const DecoratorRuntime = `
type tsClassDecoration struct {class tsValue;initializers []tsValue}
func tsDecorateClass(value tsValue,decorators []tsValue,name string,legacy bool)tsClassDecoration{result:=tsClassDecoration{class:value};for i:=len(decorators)-1;i>=0;i--{args:=[]tsValue{result.class};active:=true;if !legacy{context:=tsNewObject();context.set("kind",tsStringReference(tsStringUTF8("class")));context.set("name",tsStringReference(tsStringUTF8(name)));context.set("addInitializer",tsFunc(func(args ...tsValue)tsValue{if !active{panic("Cannot add initializers after decoration")};fn:=tsArg(args,0);if fn.kind!=tsFunctionKind{panic("Initializer must be callable")};result.initializers=append(result.initializers,fn);return tsU}));args=append(args,tsObjectValue(context))};replacement:=tsCall(decorators[i],args...);active=false;if !tsIsUndefined(replacement){if replacement.kind!=tsClassKind{panic("Class decorator must return a constructor or undefined")};result.class=replacement}};return result}
func tsRunClassInitializers(decoration tsClassDecoration){for _,fn:=range decoration.initializers{tsCall(fn)}}
`
