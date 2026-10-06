package goemit

import (
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"

	"github.com/microsoft/TypeScript/tsc/internal/ast"
)

// A specialization retains the declaring class of every inherited body. This
// makes super lexical while ordinary this calls dispatch on the concrete child.
type classInfo struct {
	node                         *ast.Node
	name                         string
	parent                       *classInfo
	fields                       map[string]*ast.Node
	fieldOrder                   []string
	methods                      map[string]*ast.Node
	methodOrder                  []string
	constructor                  *ast.Node
	prepared, preparing, emitted bool
	environment                  []*binding
	static                       bool
	statics                      *classInfo
}

func memberName(name string) string { return fmt.Sprintf("P%x", []byte(name)) }
func (e *emitter) classReference(node *ast.Node) *classInfo {
	if node == nil || node.Kind != ast.KindIdentifier {
		return nil
	}
	if declaration := e.reference(node); declaration != nil {
		return e.classes[declaration]
	}
	if node.Kind == ast.KindIdentifier {
		for _, c := range e.classOrder {
			if c.node.Name() != nil && c.node.Name().Text() == node.Text() {
				return c
			}
		}
	}
	return nil
}
func (e *emitter) prepareClasses() {
	for _, c := range e.classOrder {
		e.prepareClass(c)
	}
}
func (e *emitter) prepareClass(c *classInfo) {
	if c.prepared {
		return
	}
	if c.preparing {
		e.fail(c.node, "cyclic class inheritance")
		return
	}
	c.preparing = true
	// File identity keeps generated declarations distinct in a multi-file bundle.
	h := fnv.New64a()
	h.Write([]byte(e.file.FileName()))
	c.name = fmt.Sprintf("%s_%x", c.name, h.Sum64())
	if clauses := c.node.AsClassDeclaration().HeritageClauses; clauses != nil {
		for _, clause := range clauses.Nodes {
			hc := clause.AsHeritageClause()
			if hc.Token != ast.KindExtendsKeyword {
				continue
			}
			if len(hc.Types.Nodes) != 1 {
				e.fail(clause, "a class must have one base class")
				continue
			}
			expression := hc.Types.Nodes[0].AsExpressionWithTypeArguments().Expression
			c.parent = e.classReference(expression)
			if c.parent == nil {
				e.fail(expression, "base class must be a statically resolved class in the same source file")
				continue
			}
			e.prepareClass(c.parent)
			if !c.parent.prepared {
				c.parent = nil
				continue
			}
			for _, name := range c.parent.fieldOrder {
				c.fields[name] = c.parent.fields[name]
				c.fieldOrder = append(c.fieldOrder, name)
			}
			for _, name := range c.parent.methodOrder {
				c.methods[name] = c.parent.methods[name]
				c.methodOrder = append(c.methodOrder, name)
			}
		}
	}
	for _, member := range c.node.AsClassDeclaration().Members.Nodes {
		if member.Kind == ast.KindClassStaticBlockDeclaration || ast.HasSyntacticModifier(member, ast.ModifierFlagsStatic) {
			continue
		}
		switch member.Kind {
		case ast.KindConstructor:
			if member.Body() != nil {
				c.constructor = member
				for _, param := range member.Parameters() {
					if ast.GetCombinedModifierFlags(param)&ast.ModifierFlagsParameterPropertyModifier == 0 {
						continue
					}
					if param.Name() == nil || param.Name().Kind != ast.KindIdentifier {
						e.fail(param, "parameter properties require identifier names")
						continue
					}
					name := param.Name().Text()
					if c.fields[name] == nil {
						c.fieldOrder = append(c.fieldOrder, name)
					}
					c.fields[name] = param
				}

			}
		case ast.KindPropertyDeclaration:
			if member.Name() == nil || (member.Name().Kind != ast.KindIdentifier && member.Name().Kind != ast.KindStringLiteral) {
				e.fail(member, "class fields require ordinary literal names")
				continue
			}
			shape := e.primitive(member)
			if member.Type() != nil && member.Type().Kind == ast.KindBigIntKeyword {
				e.fail(member, "bigint fields are not supported yet")
				continue
			}
			_ = shape
			name := member.Name().Text()
			if c.fields[name] == nil {
				c.fieldOrder = append(c.fieldOrder, name)
			}
			c.fields[name] = member
		case ast.KindMethodDeclaration:
			if ast.GetFunctionFlags(member)&ast.FunctionFlagsGenerator != 0 {
				e.fail(member, "generator methods are not supported")
				continue
			}
			if member.Name() == nil || (member.Name().Kind != ast.KindIdentifier && member.Name().Kind != ast.KindStringLiteral) {
				e.fail(member, "methods require ordinary literal names")
				continue
			}
			if member.Body() == nil {
				continue
			}
			name := member.Name().Text()
			if c.methods[name] == nil {
				c.methodOrder = append(c.methodOrder, name)
			}
			c.methods[name] = member
		case ast.KindSemicolonClassElement:
		default:
			e.fail(member, "unsupported class member "+member.Kind.String())
		}
	}
	for _, name := range c.fieldOrder {
		if c.methods[name] != nil {
			e.fail(c.fields[name], "a field and method cannot share a name")
		}
	}
	c.statics = &classInfo{node: c.node, name: c.name + "Static", static: true, fields: map[string]*ast.Node{}, methods: map[string]*ast.Node{}, prepared: true}
	st := c.statics
	if c.parent != nil {
		st.parent = c.parent.statics
		for _, name := range st.parent.fieldOrder {
			st.fieldOrder = append(st.fieldOrder, name)
			st.fields[name] = st.parent.fields[name]
		}
		for _, name := range st.parent.methodOrder {
			st.methodOrder = append(st.methodOrder, name)
			st.methods[name] = st.parent.methods[name]
		}
	}
	for _, member := range c.node.AsClassDeclaration().Members.Nodes {
		if member.Kind == ast.KindClassStaticBlockDeclaration {
			continue
		}
		if !ast.HasSyntacticModifier(member, ast.ModifierFlagsStatic) {
			continue
		}
		if member.Name() == nil || (member.Name().Kind != ast.KindIdentifier && member.Name().Kind != ast.KindStringLiteral) {
			e.fail(member, "static members require ordinary literal names")
			continue
		}
		name := member.Name().Text()
		switch member.Kind {
		case ast.KindPropertyDeclaration:
			if st.fields[name] == nil {
				st.fieldOrder = append(st.fieldOrder, name)
			}
			st.fields[name] = member
		case ast.KindMethodDeclaration:
			if member.Body() == nil {
				continue
			}
			if ast.GetFunctionFlags(member)&ast.FunctionFlagsGenerator != 0 {
				e.fail(member, "generator methods are not supported")
				continue
			}
			if st.methods[name] == nil {
				st.methodOrder = append(st.methodOrder, name)
			}
			st.methods[name] = member
		default:
			e.fail(member, "unsupported static member "+member.Kind.String())
		}
	}
	for _, name := range st.fieldOrder {
		if st.methods[name] != nil {
			e.fail(st.fields[name], "a static field and method cannot share a name")
		}
	}
	c.preparing = false
	c.prepared = true
}
func (c *classInfo) ownerOf(method *ast.Node) *classInfo {
	for at := c; at != nil; at = at.parent {
		for _, member := range at.node.AsClassDeclaration().Members.Nodes {
			if member == method {
				return at
			}
		}
	}
	return nil
}
func implName(c *classInfo, name string) string { return "impl_" + c.name + "_" + memberName(name) }
func ctorName(c *classInfo) string              { return "ctor_" + c.name }
func (b *machineBuilder) classDeclaration(node *ast.Node) {
	c := b.e.classes[node]
	if c == nil {
		return
	}
	if !c.emitted {
		b.emitClass(c)
		b.emitClass(c.statics)
	}
	cell := b.e.bindings[node]
	if cell == nil {
		return
	}
	var make strings.Builder
	fmt.Fprintf(&make, "func() *tsClass {class:=&tsClass{layout:&tsLayout_%s};", c.name)
	if c.parent != nil {
		if parent := b.e.bindings[c.parent.node]; parent != nil {
			make.WriteString("class.parent=tsClassPointer(" + parent.name + ".get());")
		}
	}
	fmt.Fprintf(&make, "class.construct=func(args ...tsValue) tsValue {self:=&%s{loop:loop, properties:tsNewProperties()};self.properties.class=class", c.name)
	for _, capture := range c.environment {
		fmt.Fprintf(&make, ";self.%s=%s", capture.name, capture.name)
	}
	make.WriteString(";self.initProperties()")
	fmt.Fprintf(&make, ";self.%s(args...);if !self.properties.initialized {panic(\"Derived constructor did not call super\")};return self};return class}()", ctorName(c))
	b.emit(cell.name + ".init(" + make.String() + ")")
	st := c.statics
	singleton := b.typedTemp("&"+st.name+"{loop:loop,properties:tsNewProperties()}", "*"+st.name)
	target := singleton
	b.emit("tsClassPointer(" + cell.name + ".get()).static=tsInstanceValue(" + target + ".properties)")
	b.emit(target + ".properties.class=tsClassPointer(" + cell.name + ".get())")
	if st.parent != nil {
		b.emit(target + ".parent=tsAs_" + st.parent.name + "(tsClassPointer(" + cell.name + ".get()).parent.static)")
		b.emit(target + ".properties.prototype=" + target + ".parent.properties")
	}
	for _, capture := range st.environment {
		b.emit(target + "." + capture.name + "=" + capture.name)
	}
	b.emit(target + ".initProperties()")
	b.emit(target + "." + ctorName(st) + "()")
	if b.module && ast.HasSyntacticModifier(node, ast.ModifierFlagsExport) {
		name := "default"
		if node.Name() != nil {
			name = node.Name().Text()
		}
		if ast.HasSyntacticModifier(node, ast.ModifierFlagsDefault) {
			name = "default"
		}
		b.publish(name, cell.name+".get()", cell.name)
	}
}
func (b *machineBuilder) emitClass(c *classInfo) {
	c.emitted = true
	candidates := append(append([]*binding{}, b.captures...), b.locals...)
	needed := map[*binding]bool{}
	var visit func(*ast.Node)
	visit = func(node *ast.Node) {
		if node.Kind == ast.KindIdentifier && (node.Parent == nil || node.Parent.Name() != node) {
			if cell := b.e.binding(node); cell != nil {
				needed[cell] = true
			}
		}
		node.ForEachChild(func(child *ast.Node) bool { visit(child); return false })
	}
	for at := c; at != nil; at = at.parent {
		visit(at.node)
	}
	for _, cell := range candidates {
		if needed[cell] {
			c.environment = append(c.environment, cell)
		}
	}
	var out strings.Builder
	fmt.Fprintf(&out, "\ntype %s struct {loop *tsLoop; properties *tsProperties\n", c.name)
	if c.static && c.parent != nil {
		fmt.Fprintf(&out, "parent *%s\n", c.parent.name)
	}
	for _, cell := range c.environment {
		fmt.Fprintf(&out, "%s %s\n", cell.name, cell.pointerType())
	}
	for _, name := range c.fieldOrder {
		id := memberName(name)
		shape := b.e.primitive(c.fields[name])
		typ := shape.goType()
		if shape.nulls != 0 && shape.kind != "" {
			typ = "tsOptional[" + typ + "]"
		}
		fmt.Fprintf(&out, "%s %s;has%s bool\n", id, typ, id)
	}
	out.WriteString("}\n")
	// Concrete layouts are checked before casting; parent-typed references use
	// views of callbacks, so descendants need no Go interface boxing.
	fmt.Fprintf(&out, "var tsLayout_%s=tsClassLayout{name:%q}\n", c.name, c.name)
	fmt.Fprintf(&out, "func tsAs_%s(value tsValue)*%s{p:=tsInstanceProperties(value);if p.layout!=&tsLayout_%s{panic(\"Invalid concrete class layout\")};return (*%s)(p.self)}\n", c.name, c.name, c.name, c.name)
	fmt.Fprintf(&out, "type %sView struct {\n", c.name)
	for _, name := range c.fieldOrder {
		id := memberName(name)
		fmt.Fprintf(&out, "Get%s func()tsValue;Set%s func(tsValue)tsValue\n", id, id)
	}
	for _, name := range c.methodOrder {
		fmt.Fprintf(&out, "Call%s func(...tsValue)tsValue\n", memberName(name))
	}
	out.WriteString("}\n")
	fmt.Fprintf(&out, "func tsView_%s(value tsValue)%sView{p:=tsInstanceProperties(value);for class:=p.class;class!=nil;class=class.parent{if class.layout==&tsLayout_%s{return %sView{", c.name, c.name, c.name, c.name)
	for _, name := range c.fieldOrder {
		id := memberName(name)
		fmt.Fprintf(&out, "Get%s:p.declared[%q].get,Set%s:p.declared[%q].set,", id, name, id, name)
	}
	for _, name := range c.methodOrder {
		fmt.Fprintf(&out, "Call%s:(*tsFunction)(p.methods[%q].ref).call,", memberName(name), name)
	}
	out.WriteString("}}};panic(\"Invalid class view\")}\n")
	for _, name := range c.fieldOrder {
		id := memberName(name)
		shape := b.e.primitive(c.fields[name])
		value := "self." + id
		storage := shape.goType()
		checked := "value"
		if shape.kind != "" {
			checked = fmt.Sprintf("tsBoundary(value,%q,%d,%t)", shape.kind, shape.nulls, b.e.coerce)
		}
		if shape.nulls != 0 && shape.kind != "" {
			value += ".raw()"
			storage = "tsOptional[" + storage + "]"
			checked = "tsOptionalFrom[" + shape.goType() + "](" + checked + ")"
		} else if shape.kind != "" {
			checked = "tsNative[" + storage + "](" + checked + ")"
		}
		fallback := ""
		if c.static && c.parent != nil && c.parent.fields[name] != nil {
			fallback = fmt.Sprintf("if !self.properties.own[%q] {return self.parent.Get%s()};", name, id)
		}
		fmt.Fprintf(&out, "func(self *%s) Get%s() tsValue {%sif !self.has%s {return tsU};return %s}\n", c.name, id, fallback, id, value)
		fmt.Fprintf(&out, "func(self *%s) Set%s(value tsValue) tsValue {self.%s=%s;self.has%s=true;self.properties.define(%q);return self.Get%s()}\n", c.name, id, id, checked, id, name, id)

	}
	fmt.Fprintf(&out, "func(self *%s) initProperties() {self.properties.self=unsafe.Pointer(self);self.properties.layout=&tsLayout_%s;\n", c.name, c.name)
	for _, name := range c.fieldOrder {
		id := memberName(name)
		fmt.Fprintf(&out, "self.properties.declared[%q]=tsProperty{get:func()tsValue{return self.Get%s()},set:func(value tsValue)tsValue{return self.Set%s(value)}}\n", name, id, id)
	}
	for _, name := range c.methodOrder {
		fmt.Fprintf(&out, "self.properties.methods[%q]=tsFunc(self.Call%s)\n", name, memberName(name))
	}
	out.WriteString("}\n")
	fmt.Fprintf(&out, "func(self *%s) newBlank() *%s {other:=&%s{loop:self.loop,properties:tsNewProperties()};other.properties.class=self.properties.class\n", c.name, c.name, c.name)
	for _, cell := range c.environment {
		fmt.Fprintf(&out, "other.%s=self.%s\n", cell.name, cell.name)
	}
	out.WriteString("other.initProperties();return other}\n")
	for _, name := range c.methodOrder {
		method := c.methods[name]
		owner := c.ownerOf(method)
		fmt.Fprintf(&out, "func(self *%s) Call%s(args ...tsValue) tsValue {return self.%s(args...)}\n", c.name, memberName(name), implName(owner, name))
	}
	for at := c; at != nil; at = at.parent {
		for index, member := range at.node.AsClassDeclaration().Members.Nodes {
			if c.static && member.Kind == ast.KindClassStaticBlockDeclaration {
				fmt.Fprintf(&out, "func(self *%s) staticBlock_%s_%d(args ...tsValue) tsValue {\n%s\n}\n", c.name, at.name, index, b.classBody(c, at, member, false))
			}
			if member.Kind == ast.KindMethodDeclaration && member.Body() != nil && member.Name() != nil && ast.HasSyntacticModifier(member, ast.ModifierFlagsStatic) == c.static {
				fmt.Fprintf(&out, "func(self *%s) %s(args ...tsValue) tsValue {\n%s\n}\n", c.name, implName(at, member.Name().Text()), b.classBody(c, at, member, false))
			}
		}
		fmt.Fprintf(&out, "func(self *%s) %s(args ...tsValue) tsValue {\n%s\n}\n", c.name, ctorName(at), b.classBody(c, at, at.constructor, true))
	}
	b.e.classText.WriteString(out.String())
}
func (b *machineBuilder) classBody(concrete, declaring *classInfo, node *ast.Node, constructor bool) string {
	owner := node
	if owner == nil {
		owner = declaring.node
	}
	child := b.e.newMachine(owner, nil, node != nil && ast.HasSyntacticModifier(node, ast.ModifierFlagsAsync))
	child.concrete = concrete
	child.declaring = declaring
	child.constructor = constructor
	child.direct = directClassBody(node)
	child.receiver = "self"
	child.captures = concrete.environment
	// A default constructor has no function-local declarations.
	if node == nil {
		child.locals = nil
	}
	if constructor && (declaring.parent == nil || declaring.static) {
		child.emit("self.properties.initialized=true")
		child.initializeFields()
	}
	if node != nil && node.Kind != ast.KindClassStaticBlockDeclaration {
		for index, param := range node.Parameters() {
			child.initializeParameter(param, index)
		}
		if constructor && (declaring.parent == nil || declaring.static) {
			child.initializeParameterProperties()
		}
	}
	if constructor && node == nil && declaring.parent != nil && !declaring.static {
		child.emit("self." + ctorName(declaring.parent) + "(args...)")
		child.initializeFields()
	}
	if node != nil {
		if child.direct {
			body := node.Body()
			if node.Kind == ast.KindClassStaticBlockDeclaration {
				body = node.AsClassStaticBlockDeclaration().Body
			}
			child.directStatements(body.AsBlock().Statements.Nodes)
		} else {
			body := node.Body()
			if node.Kind == ast.KindClassStaticBlockDeclaration {
				body = node.AsClassStaticBlockDeclaration().Body
			}
			child.statements(body.AsBlock().Statements.Nodes)
		}
	}
	if !child.direct {
		child.abrupt("return", "tsU", 0, 0)
	}
	var prefix strings.Builder
	prefix.WriteString("loop:=self.loop\n")
	for _, cell := range concrete.environment {
		fmt.Fprintf(&prefix, "%s:=self.%s;_=%s\n", cell.name, cell.name, cell.name)
	}
	if child.direct {
		prefix.WriteString("_ = loop\n")
		prefix.WriteString(child.finishDirect())
	} else {
		prefix.WriteString("return tsCall(" + child.finish() + ",args...)")
	}
	return prefix.String()
}
func (b *machineBuilder) initializeFields() {
	if !b.declaring.static && b.declaring.constructor != nil {
		for _, param := range b.declaring.constructor.Parameters() {
			if ast.GetCombinedModifierFlags(param)&ast.ModifierFlagsParameterPropertyModifier != 0 {
				b.initializeEmptyField(param)
			}
		}
	}

	for index, member := range b.declaring.node.AsClassDeclaration().Members.Nodes {
		if b.declaring.static && member.Kind == ast.KindClassStaticBlockDeclaration {
			b.emit(fmt.Sprintf("self.staticBlock_%s_%d()", b.declaring.name, index))
			continue
		}
		if member.Kind != ast.KindPropertyDeclaration || (ast.HasSyntacticModifier(member, ast.ModifierFlagsStatic) != b.declaring.static || ast.HasSyntacticModifier(member, ast.ModifierFlagsAmbient|ast.ModifierFlagsAbstract)) {
			continue
		}
		p := member.AsPropertyDeclaration()
		if member.Name() == nil {
			continue
		}
		id := memberName(member.Name().Text())
		if p.Initializer == nil {
			b.initializeEmptyField(member)
		} else {
			value := ""
			shape := b.e.primitive(member)
			if shape.numeric() && shape.nulls == 0 {
				value = b.numericExpression(p.Initializer, shape)
			} else {
				value = b.expression(p.Initializer)
			}
			b.emit("self.Set" + id + "(" + value + ")")
		}
	}
}

// classOf uses declaration identity and explicit annotations. Unknown receivers
// retain dynamic bracket access; concrete class dot access bypasses lookup.
func (b *machineBuilder) classOf(node *ast.Node) *classInfo {
	if node == nil {
		return nil
	}
	switch node.Kind {
	case ast.KindThisKeyword:
		return b.concrete
	case ast.KindParenthesizedExpression:
		return b.classOf(node.AsParenthesizedExpression().Expression)
	case ast.KindAsExpression:
		return b.classType(node.AsAsExpression().Type)
	case ast.KindNewExpression:
		return b.e.classReference(node.AsNewExpression().Expression)
	case ast.KindIdentifier:
		if c := b.e.classReference(node); c != nil {
			return c.statics
		}
		decl := b.e.reference(node)
		if decl == nil {
			return nil
		}
		if decl.Kind == ast.KindVariableDeclaration {
			d := decl.AsVariableDeclaration()
			if c := b.classType(d.Type); c != nil {
				return c
			}
			if d.Initializer != nil && d.Initializer.Kind == ast.KindNewExpression {
				return b.classOf(d.Initializer)
			}
		}
		if decl.Kind == ast.KindParameter {
			return b.classType(decl.AsParameterDeclaration().Type)
		}
	case ast.KindElementAccessExpression:
		expr := node.AsElementAccessExpression().Expression
		if expr.Kind == ast.KindIdentifier {
			decl := b.e.reference(expr)
			if decl != nil && decl.Kind == ast.KindVariableDeclaration {
				typ := decl.AsVariableDeclaration().Type
				if typ != nil && typ.Kind == ast.KindArrayType {
					return b.classType(typ.AsArrayTypeNode().ElementType)
				}
			}
		}
	}
	return nil
}
func (b *machineBuilder) classType(node *ast.Node) *classInfo {
	if node != nil && node.Kind == ast.KindTypeReference {
		return b.e.classReference(node.AsTypeReferenceNode().TypeName)
	}
	return nil
}
func (b *machineBuilder) classProperty(node *ast.Node, receiver, name string) (string, bool) {
	if node.Kind == ast.KindSuperKeyword {
		if b.declaring == nil || b.declaring.parent == nil {
			b.e.fail(node, "super requires a derived class")
			return "tsU", true
		}
		if b.declaring.static && b.declaring.parent.fields[name] != nil {
			target := b.receiver
			for at := b.concrete; at != nil && at != b.declaring.parent; at = at.parent {
				target += ".parent"
			}
			return target + ".Get" + memberName(name) + "()", true
		}
		method := b.declaring.parent.methods[name]
		if method == nil {
			b.e.fail(node, "super property must name a parent method")
			return "tsU", true
		}
		owner := b.declaring.parent.ownerOf(method)
		return "tsFunc(" + b.receiver + "." + implName(owner, name) + ")", true
	}
	c := b.classOf(node)
	if c == nil {
		return "", false
	}
	id := memberName(name)
	target := b.classTarget(node, receiver, c)
	if node.Kind == ast.KindThisKeyword {
		target = b.receiver
	}
	if c.fields[name] != nil {
		return target + ".Get" + id + "()", true
	}
	if c.methods[name] != nil {
		return "tsFunc(" + target + ".Call" + id + ")", true
	}
	b.e.fail(node, "unknown class member "+strconv.Quote(name))
	return "tsU", true
}

// A concrete allocation needs neither a property hash nor interface dispatch.
// Explicit parent annotations dispatch through concrete callback views.
func (b *machineBuilder) classTarget(node *ast.Node, receiver string, c *classInfo) string {
	if node.Kind == ast.KindThisKeyword {
		return b.receiver
	}
	if c.static {
		return "tsAs_" + c.name + "(tsClassPointer(" + receiver + ").static)"
	}
	concrete := node.Kind == ast.KindNewExpression
	if node.Kind == ast.KindIdentifier {
		if decl := b.e.reference(node); decl != nil && decl.Kind == ast.KindVariableDeclaration {
			d := decl.AsVariableDeclaration()
			concrete = d.Type == nil && d.Initializer != nil && d.Initializer.Kind == ast.KindNewExpression
		}
	}
	if concrete {
		return "tsAs_" + c.name + "(" + receiver + ")"
	}
	return "tsView_" + c.name + "(" + receiver + ")"
}

func numericFieldInitializer(node *ast.Node) bool {
	if node.Kind == ast.KindNumericLiteral {
		return true
	}
	if node.Kind == ast.KindPrefixUnaryExpression {
		n := node.AsPrefixUnaryExpression()
		return (n.Operator == ast.KindMinusToken || n.Operator == ast.KindPlusToken) && numericFieldInitializer(n.Operand)
	}
	return false
}

func (b *machineBuilder) initializeEmptyField(member *ast.Node) {
	id := memberName(member.Name().Text())
	shape := b.e.primitive(member)
	zero := "tsU"
	switch {
	case shape.numeric():
		zero = "0"
	case shape.kind == "boolean":
		zero = "false"
	case shape.kind == "string":
		zero = "nil"
	}
	if shape.nulls != 0 && shape.kind != "" {
		zero = "tsOptional[" + shape.goType() + "]{tag:2}"
	}
	b.emit(fmt.Sprintf("self.properties.define(%q)", member.Name().Text()))
	b.emit("self.has" + id + "=false;self." + id + "=" + zero)
}
func (b *machineBuilder) initializeParameterProperties() {
	if !b.constructor || b.owner == nil || b.owner.Kind != ast.KindConstructor {
		return
	}
	for _, param := range b.owner.Parameters() {
		if ast.GetCombinedModifierFlags(param)&ast.ModifierFlagsParameterPropertyModifier == 0 {
			continue
		}
		if cell := b.e.bindings[param]; cell != nil {
			name := param.Name().Text()
			id := memberName(name)
			field := b.e.primitive(b.concrete.fields[name])
			if cell.primitive.kind != "" && cell.primitive == field {
				read := "read"
				if field.nulls != 0 {
					read = "optional"
				}
				b.emit("self." + id + "=" + cell.name + "." + read + "();self.has" + id + "=true")
				b.emit(fmt.Sprintf("self.properties.define(%q)", name))
			} else {
				b.emit("self.Set" + id + "(" + cell.name + ".get())")
			}
		}
	}
}
