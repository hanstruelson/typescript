package goemit

import (
	"fmt"
	"github.com/microsoft/TypeScript/tsc/internal/ast"
	"strings"
)

func conversionName(kind string) string {
	switch kind {
	case "float64":
		return "Number"
	case "*tsString":
		return "String"
	case "bool":
		return "Boolean"
	case "number":
		return "Number"
	}
	return strings.ToUpper(kind[:1]) + kind[1:]
}
func (b *machineBuilder) nativeConversion(value string, p primitive, coerce bool) (string, bool) {
	source := b.tempType(value)
	if p.nulls == 0 && source == p.goType() {
		return value, true
	}
	if source == "tsValue" || source == "" || p.kind == "" {
		return "", false
	}
	if !(primitive{kind: source}).numeric() && source != "bool" && source != "*tsString" {
		return "", false
	}
	return b.typedTemp(fmt.Sprintf("tsConvert%sFrom%s(%s,%t)", conversionName(p.kind), conversionName(source), value, coerce), p.goType()), true
}
func (b *machineBuilder) explicitConversion(value string, p primitive) string {
	if result, ok := b.nativeConversion(value, p, true); ok {
		return result
	}
	method := conversionName(p.kind)
	if p.kind == "string" {
		method = "TSString"
	}
	return b.typedTemp("("+value+").To"+method+"()", p.goType())
}

func (b *machineBuilder) primitiveConversionTarget(node *ast.Node) string {
	if resolver, ok := b.e.resolver.(interface{ GetEmitPrimitiveConversion(*ast.Node) string }); ok {
		return resolver.GetEmitPrimitiveConversion(node)
	}
	return ""
}
