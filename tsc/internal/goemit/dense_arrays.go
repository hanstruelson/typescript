package goemit

import (
	"fmt"
	"github.com/microsoft/TypeScript/tsc/internal/ast"
	"strconv"
	"strings"
)

func (b *machineBuilder) denseElement(node *ast.Node, p primitive) string {
	if p.nulls != 0 {
		tag := 0
		if node.Kind == ast.KindNullKeyword {
			tag = 1
		} else if node.Kind == ast.KindIdentifier && node.Text() == "undefined" && b.e.binding(node) == nil {
			tag = 2
		}
		if tag != 0 {
			return b.typedTemp(fmt.Sprintf("tsNullableCheck(tsOptional[%s]{tag:%d},%d)", p.goType(), tag, p.nulls), p.elementGoType())
		}
	}

	value := ""
	if p.numeric() && (p.nulls == 0 || node.Kind == ast.KindNumericLiteral) {
		base := p
		base.nulls = 0
		value = b.numericExpression(node, base)
	} else {
		value = b.expression(node)
	}
	if b.tempType(value) != p.elementGoType() {
		value = b.comparisonBoundary(value, p)
	}
	return value
}
func (b *machineBuilder) denseSpread(target string, node *ast.Node, p primitive) {
	if p.nulls != 0 {
		source := b.e.arrayElementPrimitive(node)
		if source.kind != "" {
			storage := b.denseStorage(node, source)
			index := b.e.unique("SpreadIndex")
			b.emit("for " + index + ":=range " + storage + ".values {")
			value := b.typedTemp("tsNullableAt["+source.goType()+"]("+storage+","+index+")", "tsOptional["+source.goType()+"]")
			converted := b.comparisonBoundary(value, p)
			b.emit(target + "=append(" + target + "," + converted + ")")
			b.emit("}")
			return
		}
	}

	source := b.e.arrayElementPrimitive(node)
	if source.kind != "" && source.nulls == 0 {
		storage := b.denseStorage(node, source)
		if source.goType() == p.goType() {
			b.emit(target + "=append(" + target + "," + storage + ".values...)")
			return
		}
		item := b.e.unique("ArrayItem")
		b.tempTypes[item] = source.goType()
		b.emit("for _," + item + ":=range " + storage + ".values {")
		converted := b.comparisonBoundary(item, p)
		b.emit(target + "=append(" + target + "," + converted + ")")
		b.emit("}")
		return
	}
	value := b.expression(node)
	iterator := b.typedTemp("tsIterate("+value+")", "*tsIterator")
	b.emit("for " + iterator + ".next() {")
	item := b.comparisonBoundary(b.temp(iterator+".value"), p)
	b.emit(target + "=append(" + target + "," + item + ")")
	b.emit("}")
}
func (b *machineBuilder) denseArguments(nodes []*ast.Node, p primitive) []string {
	spread := false
	for _, node := range nodes {
		spread = spread || node.Kind == ast.KindSpreadElement
	}
	if !spread {
		args := []string{}
		for _, node := range nodes {
			args = append(args, b.denseElement(node, p))
		}
		return args
	}
	items := b.typedTemp("[]"+p.elementGoType()+"{}", "[]"+p.elementGoType())
	for _, node := range nodes {
		if node.Kind == ast.KindSpreadElement {
			b.denseSpread(items, node.AsSpreadElement().Expression, p)
		} else {
			value := b.denseElement(node, p)
			b.emit(items + "=append(" + items + "," + value + ")")
		}
	}
	return []string{items + "..."}
}
func (b *machineBuilder) denseIndex(node *ast.Node) string {
	value := b.expression(node)
	if (primitive{kind: b.tempType(value)}).numeric() {
		return "float64(" + value + ")"
	}
	return b.comparisonBoundary(value, primitive{kind: "number"})
}
func (b *machineBuilder) denseArrayCall(call *ast.CallExpression) (string, bool) {
	if call.Expression.Kind != ast.KindPropertyAccessExpression {
		return "", false
	}
	property := call.Expression.AsPropertyAccessExpression()
	p := b.e.arrayElementPrimitive(property.Expression)
	if value, ok := b.nativeCallbackCall(call, p); ok {
		return value, true
	}
	if p.kind == "" {
		return "", false
	}
	name := property.Name().Text()
	switch name {
	case "push", "unshift", "pop", "shift", "splice", "slice", "fill", "reverse", "copyWithin":
	default:
		return "", false
	}
	for _, node := range call.Arguments.Nodes {
		if node.Kind == ast.KindSpreadElement && name != "push" && name != "unshift" {
			return "", false
		}
	}
	if name == "fill" && len(call.Arguments.Nodes) == 0 {
		return "", false
	}
	storage := b.denseStorage(property.Expression, p)
	nodes := call.Arguments.Nodes
	switch name {
	case "push", "unshift":
		args := append([]string{storage}, b.denseArguments(nodes, p)...)
		helper := "tsDensePush"
		if p.nulls != 0 {
			helper = "tsNullablePush"
		}
		if name == "unshift" {
			helper = "tsDenseUnshift"
			if p.nulls != 0 {
				helper = "tsNullableUnshift"
			}
		}
		return b.typedTemp(helper+"["+p.goType()+"]("+strings.Join(args, ",")+")", "float64"), true
	case "pop", "shift":
		for _, arg := range nodes {
			b.emit("_=" + b.expression(arg))
		}
		helper := "tsDensePop"
		if p.nulls != 0 {
			helper = "tsNullablePop"
		}
		return b.typedTemp(fmt.Sprintf(helper+"[%s](%s,%t)", p.goType(), storage, name == "shift"), "tsOptional["+p.goType()+"]"), true
	case "splice", "slice":
		start := "float64(0)"
		if len(nodes) > 0 {
			start = b.denseIndex(nodes[0])
		}
		second := "math.Inf(1)"
		if name == "splice" && len(nodes) == 0 {
			second = "float64(0)"
		}
		if len(nodes) > 1 {
			second = b.denseIndex(nodes[1])
		}
		args := []string{storage, strconv.Quote(p.kind), strconv.FormatBool(b.e.coerce), start, second}
		helper := "tsDenseSlice"
		if name == "splice" {
			helper = "tsDenseSplice"
			if p.nulls != 0 {
				helper = "tsNullableSplice"
			}
			if len(nodes) > 2 {
				args = append(args, b.denseArguments(nodes[2:], p)...)
			}
		} else {
			for _, node := range nodes[min(2, len(nodes)):] {
				b.emit("_=" + b.expression(node))
			}
		}
		return b.typedTemp(helper+"["+p.goType()+"]("+strings.Join(args, ",")+")", "*tsArray"), true
	case "fill", "copyWithin":
		args := []string{storage}
		if name == "fill" {
			if len(nodes) == 0 {
				return "", false
			}
			args = append(args, b.denseElement(nodes[0], p))
		} else {
			index := "float64(0)"
			if len(nodes) > 0 {
				index = b.denseIndex(nodes[0])
			}
			args = append(args, index)
		}
		start, end := "float64(0)", "math.Inf(1)"
		if len(nodes) > 1 {
			start = b.denseIndex(nodes[1])
		}
		if len(nodes) > 2 {
			end = b.denseIndex(nodes[2])
		}
		args = append(args, start, end)
		for _, node := range nodes[min(3, len(nodes)):] {
			b.emit("_=" + b.expression(node))
		}
		helper := "tsDenseFill"
		if p.nulls != 0 {
			helper = "tsNullableFill"
		}
		if name == "copyWithin" {
			helper = "tsDenseCopyWithin"
		}
		b.emit(helper + "[" + p.goType() + "](" + strings.Join(args, ",") + ")")
		return b.temp("tsArrayValue(tsDenseArrayWrapper(" + storage + "))"), true
	case "reverse":
		for _, node := range nodes {
			b.emit("_=" + b.expression(node))
		}
		b.emit("tsDenseReverse[" + p.goType() + "](" + storage + ")")
		return b.temp("tsArrayValue(tsDenseArrayWrapper(" + storage + "))"), true
	}
	return "", false
}
