package goemit

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/microsoft/TypeScript/tsc/internal/ast"
)

func (b *machineBuilder) propertyKey(node *ast.Node) string {
	if node.Kind == ast.KindComputedPropertyName {
		return b.expression(node.AsComputedPropertyName().Expression)
	}
	return b.typedTemp(stringLiteral(node.Text()), "*tsString")
}
func (b *machineBuilder) objectLiteral(node *ast.Node) string {
	object := b.typedTemp("tsNewObject()", "*tsObject")
	for _, property := range node.AsObjectLiteralExpression().Properties.Nodes {
		switch property.Kind {
		case ast.KindPropertyAssignment:
			key := b.propertyKey(property.Name())
			value := b.expression(property.AsPropertyAssignment().Initializer)
			b.emit("tsSet(" + object + "," + key + "," + value + ")")
		case ast.KindShorthandPropertyAssignment:
			key := b.propertyKey(property.Name())
			value := b.expression(property.Name())
			b.emit("tsSet(" + object + "," + key + "," + value + ")")
		case ast.KindSpreadAssignment:
			value := b.expression(property.AsSpreadAssignment().Expression)
			b.emit("tsObjectSpread(" + object + "," + value + ")")
		default:
			b.e.fail(property, "object methods and accessors require receiver lowering")
		}
	}
	return object
}
func (b *machineBuilder) declarePatternCells(pattern *ast.Node) {
	for _, element := range pattern.AsBindingPattern().Elements.Nodes {
		if element.Kind != ast.KindBindingElement || element.Name() == nil || element.Name().Kind == ast.KindOmittedExpression {
			continue
		}
		if ast.IsBindingPattern(element.Name()) {
			b.declarePatternCells(element.Name())
		} else {
			b.allocate(b.e.bindings[element])
		}
	}
}
func (b *machineBuilder) bindPattern(pattern *ast.Node, value string, initialize bool) {
	if !ast.IsBindingPattern(pattern) {
		b.e.fail(pattern, "expected a binding pattern")
		return
	}
	excluded := []string{}
	iterator := ""
	if pattern.Kind == ast.KindArrayBindingPattern {
		iterator = b.typedTemp("tsIterate("+value+")", "*tsIterator")
	}
	for _, element := range pattern.AsBindingPattern().Elements.Nodes {
		if element.Kind == ast.KindOmittedExpression || element.Name() == nil || element.Name().Kind == ast.KindOmittedExpression {
			b.emit(iterator + ".next()")
			continue
		}
		n := element.AsBindingElement()
		item := ""
		if pattern.Kind == ast.KindArrayBindingPattern {
			if n.DotDotDotToken != nil {
				item = b.temp("tsIteratorRest(" + iterator + ")")
			} else {
				item = b.temp("tsIteratorValue(" + iterator + ")")
			}
		} else {
			if n.DotDotDotToken != nil {
				item = b.temp("tsObjectRest(" + value + ",[]tsValue{" + strings.Join(excluded, ",") + "})")
			} else {
				key := n.PropertyName
				if key == nil {
					key = element.Name()
				}
				name := b.propertyKey(key)
				excluded = append(excluded, name)
				item = b.temp("tsGet(" + value + "," + name + ")")
			}
		}
		if n.Initializer != nil {
			if b.direct {
				b.emit("if tsIsUndefined(" + item + ") {")
				defaultValue := b.bindingDefault(element)
				b.emit(item + "=" + defaultValue)
				b.emit("}")
			} else {
				yes, end := b.block(), b.block()
				b.emit(fmt.Sprintf("if tsIsUndefined(%s){m.pc=%d}else{m.pc=%d};return", item, yes, end))
				b.current = yes
				defaultValue := b.bindingDefault(element)
				b.emit(item + "=" + defaultValue)
				b.jump(end)
				b.current = end
			}
		}
		if ast.IsBindingPattern(element.Name()) {
			b.bindPattern(element.Name(), item, initialize)
		} else if cell := b.e.bindings[element]; cell != nil {
			method := "set"
			if initialize {
				method = "init"
			}
			if initialize && cell.primitive.kind != "" {
				item = fmt.Sprintf("tsBoundary(%s,%q,%d,%t)", item, cell.primitive.kind, cell.primitive.nulls, b.e.coerce)
			}
			b.emit(cell.name + "." + method + "(" + item + ")")
		}
	}
}
func (b *machineBuilder) taggedTemplate(node *ast.Node) string {
	n := node.AsTaggedTemplateExpression()
	tag := b.expression(n.Tag)
	cooked, raw := []string{}, []string{}
	substitutions := []*ast.Node{}
	add := func(literal *ast.Node) {
		cooked = append(cooked, stringLiteral(literal.Text()))
		text := b.e.file.Text()[literal.Loc.Pos():literal.Loc.End()]
		switch literal.Kind {
		case ast.KindTemplateHead, ast.KindTemplateMiddle:
			text = text[1 : len(text)-2]
		default:
			text = text[1 : len(text)-1]
		}
		text = strings.ReplaceAll(strings.ReplaceAll(text, "\r\n", "\n"), "\r", "\n")
		raw = append(raw, stringLiteral(text))
	}
	if n.Template.Kind == ast.KindNoSubstitutionTemplateLiteral {
		add(n.Template)
	} else {
		template := n.Template.AsTemplateExpression()
		add(template.Head)
		for _, span := range template.TemplateSpans.Nodes {
			s := span.AsTemplateSpan()
			substitutions = append(substitutions, s.Expression)
			add(s.Literal)
		}
	}
	array := b.temp("&tsArray{values:[]tsValue{" + strings.Join(cooked, ",") + "},properties:map[string]tsValue{\"raw\":&tsArray{values:[]tsValue{" + strings.Join(raw, ",") + "}}}}")
	args := []string{tag, array}
	for _, sub := range substitutions {
		args = append(args, b.expression(sub))
	}
	return b.temp("tsCall(" + strings.Join(args, ",") + ")")
}
func (b *machineBuilder) spreadArguments(nodes []*ast.Node) string {
	array := b.typedTemp("&tsArray{}", "*tsArray")
	for _, node := range nodes {
		if node.Kind == ast.KindSpreadElement {
			value := b.expression(node.AsSpreadElement().Expression)
			b.emit("tsArraySpread(" + array + "," + value + ")")
		} else {
			value := b.expression(node)
			b.emit(array + ".values=append(" + array + ".values," + value + ")")
		}
	}
	return array + ".values"
}
func (b *machineBuilder) forInStatement(node *ast.Node, label string) {
	n := node.AsForInOrOfStatement()
	value := b.expression(n.Expression)
	keys := b.temp("tsForInKeys(" + value + ")")
	iterator := b.typedTemp("tsIterate("+keys+")", "*tsIterator")
	test, body, end := b.block(), b.block(), b.block()
	b.jump(test)
	b.loops = append(b.loops, loopTarget{end, test, b.depth, label})
	b.current = test
	b.emit(fmt.Sprintf("if %s.next(){m.pc=%d}else{m.pc=%d};return", iterator, body, end))
	b.current = body
	if n.Initializer.Kind == ast.KindVariableDeclarationList {
		decl := n.Initializer.AsVariableDeclarationList().Declarations.Nodes[0]
		cell := b.e.bindings[decl]
		if cell == nil {
			b.e.fail(decl, "for-in requires an identifier binding")
		} else if cell.lexical {
			b.allocate(cell)
			b.emit(cell.name + ".init(" + iterator + ".value)")
		} else {
			b.emit(cell.name + ".set(" + iterator + ".value)")
		}
	} else {
		_, write := b.lvalue(n.Initializer)
		b.emit(write(iterator + ".value"))
	}
	b.statement(n.Statement)
	b.jump(test)
	b.loops = b.loops[:len(b.loops)-1]
	b.current = end
}
func (b *machineBuilder) objectBuiltin(name string) string {
	return "tsFunc(func(args ...tsValue)tsValue{return tsObjectBuiltin(" + strconv.Quote(name) + ",args)})"
}

func (b *machineBuilder) bindingDefault(element *ast.Node) string {
	initializer := element.AsBindingElement().Initializer
	if cell := b.e.bindings[element]; cell != nil && cell.primitive.numeric() && cell.primitive.nulls == 0 {
		return b.numericExpression(initializer, cell.primitive)
	}
	return b.expression(initializer)
}

func (b *machineBuilder) arrayLiteral(nodes []*ast.Node) string {
	p := primitive{}
	if len(nodes) > 0 && nodes[0].Parent != nil {
		p = b.e.arrayElementPrimitive(nodes[0].Parent)
	}
	constructor := "&tsArray{}"
	if p.kind != "" {
		constructor = fmt.Sprintf("tsNewGrowableArray(%q,%d,%t)", p.kind, p.nulls, b.e.coerce)
	}
	array := b.typedTemp(constructor, "*tsArray")
	for _, node := range nodes {
		switch node.Kind {
		case ast.KindOmittedExpression:
			b.emit("tsArrayHole(" + array + ")")
		case ast.KindSpreadElement:

			if p.kind != "" && p.nulls != 0 {
				items := b.typedTemp("[]"+p.elementGoType()+"{}", "[]"+p.elementGoType())
				b.denseSpread(items, node.AsSpreadElement().Expression, p)
				b.emit("tsNullablePush[" + p.goType() + "](tsDenseStorage[" + p.goType() + "](" + array + "," + strconv.Quote(p.kind) + "," + strconv.Itoa(int(p.nulls)) + ")," + items + "...)")
				continue
			}
			if p.kind != "" && p.nulls == 0 {
				storage := b.typedTemp("tsDenseStorage["+p.goType()+"]("+array+","+strconv.Quote(p.kind)+")", "*tsGrowableStorage["+p.goType()+"]")
				b.denseSpread(storage+".values", node.AsSpreadElement().Expression, p)
			} else {
				value := b.expression(node.AsSpreadElement().Expression)
				b.emit("tsArraySpread(" + array + "," + value + ")")
			}

		default:

			if p.kind != "" && p.nulls != 0 {
				value := b.denseElement(node, p)
				b.emit("tsNullablePush[" + p.goType() + "](tsDenseStorage[" + p.goType() + "](" + array + "," + strconv.Quote(p.kind) + "," + strconv.Itoa(int(p.nulls)) + ")," + value + ")")
				continue
			}
			value := ""
			if p.numeric() && p.nulls == 0 {
				value = b.numericExpression(node, p)
			} else {
				value = b.expression(node)
			}
			if p.kind != "" {
				if b.tempType(value) != p.elementGoType() {
					value = b.comparisonBoundary(value, p)
				}
				helper := "tsDensePush"
				if p.nulls != 0 {
					helper = "tsNullablePush"
				}
				b.emit(helper + "[" + p.goType() + "](tsDenseStorage[" + p.goType() + "](" + array + "," + strconv.Quote(p.kind) + "," + strconv.Itoa(int(p.nulls)) + ")," + value + ")")
			} else {
				b.emit(array + ".appendItems(" + value + ")")
			}
		}
	}
	return array
}
