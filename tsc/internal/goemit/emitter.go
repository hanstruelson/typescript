// Package goemit lowers TypeScript into standalone Go programs. Async bodies use
// persistent bindings and numbered resume blocks, following ES5 generator lowering.
package goemit

import (
	"fmt"
	"strings"

	"github.com/microsoft/TypeScript/tsc/internal/ast"
	"github.com/microsoft/TypeScript/tsc/internal/core"
	"github.com/microsoft/TypeScript/tsc/internal/diagnostics"
)

// ReferenceResolver is the existing compiler binding resolver. Binding identity,
// rather than identifier spelling, preserves lexical shadowing after hoisting.
type ReferenceResolver interface {
	GetReferencedValueDeclarationUnsafe(*ast.IdentifierNode) *ast.Declaration
	GetReferencedImportDeclaration(*ast.IdentifierNode) *ast.Declaration
}
type binding struct {
	name               string
	declaration, owner *ast.Node
	lexical, constant  bool
	primitive          primitive
	maybeUndefined     bool
}
type emitter struct {
	file            *ast.SourceFile
	resolver        ReferenceResolver
	bindings        map[*ast.Node]*binding
	imports         map[*ast.Node]string
	diags           []*ast.Diagnostic
	next            int
	targets         map[*ast.Node]string
	classes         map[*ast.Node]*classInfo
	classOrder      []*classInfo
	classText       strings.Builder
	coerce          bool
	nativeFunctions map[*ast.Node]*nativeFunction
}

func (e *emitter) fail(node *ast.Node, message string) {
	e.diags = append(e.diags, ast.NewDiagnosticFromText(e.file, node.Loc, 95001, diagnostics.CategoryError, "Go target: "+message, nil, nil, false, false))
}
func (e *emitter) unique(prefix string) string {
	e.next++
	return fmt.Sprintf("ts%s%d", prefix, e.next)
}
func (e *emitter) reference(node *ast.Node) *ast.Node {
	if e.resolver == nil {
		return nil
	}
	if imported := e.resolver.GetReferencedImportDeclaration(node); imported != nil {
		return imported
	}
	return e.resolver.GetReferencedValueDeclarationUnsafe(node)
}
func (e *emitter) binding(node *ast.Node) *binding { return e.bindings[e.reference(node)] }
func (e *emitter) declare(node, owner *ast.Node, lexical, constant bool) {
	if node.Kind == ast.KindBindingElement && (node.Name() == nil || node.Name().Kind == ast.KindOmittedExpression) {
		return
	}
	if node.Name() == nil && node.Kind == ast.KindFunctionDeclaration && ast.HasSyntacticModifier(node, ast.ModifierFlagsDefault) {
		e.bindings[node] = &binding{e.unique("Binding"), node, owner, true, false, primitive{}, false}
		return
	}
	if node.Name() != nil && ast.IsBindingPattern(node.Name()) {
		for _, element := range node.Name().AsBindingPattern().Elements.Nodes {
			if element.Kind == ast.KindBindingElement {
				e.declare(element, owner, lexical, constant)
				if init := element.AsBindingElement().Initializer; init != nil {
					e.collect(init, owner)
				}
			}
		}
		return
	}
	if node.Name() == nil || node.Name().Kind != ast.KindIdentifier {
		e.fail(node, "unsupported binding form "+node.Kind.String()+fmt.Sprintf(" at %d", node.Loc.Pos()))
		return
	}
	if symbol := node.Symbol(); symbol != nil && symbol.ValueDeclaration != nil && (!lexical || node.Kind == ast.KindFunctionDeclaration) {
		if b := e.bindings[symbol.ValueDeclaration]; b != nil {
			e.bindings[node] = b
			return
		}
	}
	p := primitive{}
	if node.Kind == ast.KindVariableDeclaration || node.Kind == ast.KindParameter || node.Kind == ast.KindBindingElement {
		p = e.primitive(node)

	}
	e.bindings[node] = &binding{e.unique("Binding"), node, owner, lexical, constant, p, node.Kind == ast.KindVariableDeclaration && !lexical && p.kind != ""}
}
func (e *emitter) collect(node, owner *ast.Node) {
	switch node.Kind {
	case ast.KindClassDeclaration:
		e.declare(node, owner, true, true)
		if e.classes == nil {
			e.classes = map[*ast.Node]*classInfo{}
		}
		c := &classInfo{node: node, name: e.unique("Class"), methods: map[string]*ast.Node{}, fields: map[string]*ast.Node{}}
		e.classes[node] = c
		e.classOrder = append(e.classOrder, c)
		for _, member := range node.AsClassDeclaration().Members.Nodes {
			if member.Kind == ast.KindMethodDeclaration || member.Kind == ast.KindConstructor {
				for _, param := range member.Parameters() {
					e.declare(param, member, false, false)
					if init := param.AsParameterDeclaration().Initializer; init != nil {
						e.collect(init, member)
					}
				}
				if member.Body() != nil {
					e.collect(member.Body(), member)
				}
			} else if member.Kind == ast.KindClassStaticBlockDeclaration {
				e.collect(member.AsClassStaticBlockDeclaration().Body, member)
			} else if member.Kind == ast.KindPropertyDeclaration && member.AsPropertyDeclaration().Initializer != nil {
				e.collect(member.AsPropertyDeclaration().Initializer, owner)
			}
		}
		return
	case ast.KindImportDeclaration:
		if e.targets[node] != "" {
			clause := node.AsImportDeclaration().ImportClause
			if clause != nil {
				c := clause.AsImportClause()
				if c.PhaseModifier != ast.KindTypeKeyword {
					if clause.Name() != nil {
						e.declare(clause, owner, true, true)
					}
					if c.NamedBindings != nil {
						if c.NamedBindings.Kind == ast.KindNamedImports {
							for _, specifier := range c.NamedBindings.AsNamedImports().Elements.Nodes {
								if !specifier.AsImportSpecifier().IsTypeOnly {
									e.declare(specifier, owner, true, true)
								}
							}
						} else {
							e.declare(c.NamedBindings, owner, true, true)
						}
					}
				}
			}
		}
		return
	case ast.KindFunctionDeclaration:
		if node.Body() == nil {
			return
		}
		e.declare(node, owner, true, false)
		if ast.GetFunctionFlags(node)&ast.FunctionFlagsGenerator != 0 {
			e.fail(node, "generator functions are not supported")
		}
		for _, param := range node.Parameters() {
			e.declare(param, node, false, false)
			if init := param.AsParameterDeclaration().Initializer; init != nil {
				e.collect(init, node)
			}
		}
		e.collect(node.Body(), node)
		return
	case ast.KindArrowFunction, ast.KindFunctionExpression:
		if ast.GetFunctionFlags(node)&ast.FunctionFlagsGenerator != 0 {
			e.fail(node, "generator functions are not supported")
		}
		if node.Kind == ast.KindFunctionExpression && node.Name() != nil {
			e.fail(node, "named function expressions are not supported")
		}
		for _, param := range node.Parameters() {
			e.declare(param, node, false, false)
			if init := param.AsParameterDeclaration().Initializer; init != nil {
				e.collect(init, node)
			}
		}
		e.collect(node.Body(), node)
		return
	case ast.KindVariableDeclarationList:
		list := node.AsVariableDeclarationList()
		if node.Flags&ast.NodeFlagsUsing != 0 {
			e.fail(node, "using declarations are not supported")
		}
		for _, decl := range list.Declarations.Nodes {
			e.declare(decl, owner, node.Flags&ast.NodeFlagsBlockScoped != 0, node.Flags&ast.NodeFlagsConst != 0)
			if init := decl.AsVariableDeclaration().Initializer; init != nil {
				e.collect(init, owner)
			}
		}
		return
	case ast.KindCatchClause:
		c := node.AsCatchClause()
		if c.VariableDeclaration != nil {
			e.declare(c.VariableDeclaration, owner, true, false)
		}
	}
	node.ForEachChild(func(child *ast.Node) bool { e.collect(child, owner); return false })
}

// Emit retains normal compiler checking and output handling. Unsupported syntax
// is diagnosed before tsValue output is written.
func Emit(file *ast.SourceFile, options *core.CompilerOptions, resolver ReferenceResolver) (string, []*ast.Diagnostic) {
	e := &emitter{file: file, resolver: resolver, bindings: make(map[*ast.Node]*binding), imports: make(map[*ast.Node]string)}
	e.coerce = options.CoerceAny.IsTrue()
	if options.SourceMap.IsTrue() || options.InlineSourceMap.IsTrue() {
		e.fail(file.AsNode(), "source maps are not supported")
	}
	if options.Jsx != core.JsxEmitNone {
		e.fail(file.AsNode(), "JSX is not supported")
	}
	if len(file.Diagnostics()) != 0 {
		return "", file.Diagnostics()
	}
	e.collect(file.AsNode(), file.AsNode())
	e.planNativeFunctions()
	e.prepareClasses()
	if len(e.diags) != 0 {
		return "", e.diags
	}
	builder := e.newMachine(file.AsNode(), nil, false)
	for _, node := range file.Statements.Nodes {
		if node.Kind == ast.KindImportDeclaration {
			builder.importDeclaration(node)
		}
	}
	builder.statements(file.Statements.Nodes)
	builder.emit("m.complete(tsU)")
	body := builder.finish()
	if len(e.diags) != 0 {
		return "", e.diags
	}
	text, err := formatValueSource([]byte("package main\n\n" + Runtime + ValueRuntime + ModuleRuntime + ClassRuntime + TypeRuntime + StringRuntime + RegexRuntime + EqualityRuntime + ObjectRuntime + e.classText.String() + "\nfunc main() {\nloop := tsNewLoop()\nloop.invoke(func(){ _ = tsCall(" + body + ") })\nif err := loop.run(); err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }\n}\n"))
	if err != nil {
		e.fail(file.AsNode(), "could not format output: "+err.Error())
		return "", e.diags
	}
	return string(text), nil
}

type resumeBlock struct{ lines []string }
type loopTarget struct {
	breakPC, continuePC, depth int
	label                      string
}
type machineBuilder struct {
	e                   *emitter
	owner               *ast.Node
	parent              *machineBuilder
	async               bool
	blocks              []resumeBlock
	current             int
	temps               []string
	tempTypes           map[string]string
	locals, captures    []*binding
	loops               []loopTarget
	depth               int
	module              bool
	concrete, declaring *classInfo
	direct              bool
	nativeReturn        string
	constructor         bool
	receiver            string
}

func (e *emitter) newMachine(owner *ast.Node, parent *machineBuilder, async bool) *machineBuilder {
	b := &machineBuilder{e: e, owner: owner, parent: parent, async: async, blocks: []resumeBlock{{}}}
	// Stable declaration order keeps generated programs and baselines deterministic.
	seen := map[*binding]bool{}
	var visit func(*ast.Node)
	visit = func(node *ast.Node) {
		if cell := e.bindings[node]; cell != nil && cell.owner == owner && !seen[cell] {
			b.locals = append(b.locals, cell)
			seen[cell] = true
		}
		node.ForEachChild(func(child *ast.Node) bool { visit(child); return false })
	}
	visit(owner)
	if parent != nil {
		b.captures = append(b.captures, parent.captures...)
		b.captures = append(b.captures, parent.locals...)
	}
	return b
}
func (b *machineBuilder) emit(line string) {
	b.blocks[b.current].lines = append(b.blocks[b.current].lines, line)
}
func (b *machineBuilder) block() int {
	b.blocks = append(b.blocks, resumeBlock{})
	return len(b.blocks) - 1
}
func (b *machineBuilder) jump(pc int) { b.emit(fmt.Sprintf("m.pc = %d; return", pc)) }
func (b *machineBuilder) temp(value string) string {
	name := b.e.unique("Temp")
	b.temps = append(b.temps, name)
	b.emit(name + " = " + value)
	return name
}
func (b *machineBuilder) abrupt(kind, value string, target, depth int) {
	b.emit(fmt.Sprintf("m.transfer(tsAbrupt{kind:%q, value:%s, target:%d, depth:%d}); return", kind, value, target, depth))
}
func (b *machineBuilder) statements(nodes []*ast.Node) {
	b.enter(nodes)
	for _, node := range nodes {
		b.statement(node)
	}
}
func (b *machineBuilder) allocate(cell *binding) {
	if cell != nil {
		b.emit(cell.name + " = " + b.bindingFactory(cell, false))
	}
}
func (b *machineBuilder) enter(nodes []*ast.Node) {
	// Allocate all lexical cells before creating closures, but initialize values at
	// their original execution points. Re-entering a block allocates fresh cells.
	for _, node := range nodes {
		if node.Kind == ast.KindImportDeclaration {
			if clause := node.AsImportDeclaration().ImportClause; clause != nil {
				b.allocate(b.e.bindings[clause])
				c := clause.AsImportClause()
				if c.NamedBindings != nil {
					if c.NamedBindings.Kind == ast.KindNamedImports {
						for _, specifier := range c.NamedBindings.AsNamedImports().Elements.Nodes {
							b.allocate(b.e.bindings[specifier])
						}
					} else {
						b.allocate(b.e.bindings[c.NamedBindings])
					}
				}
			}
		}
		if node.Kind == ast.KindVariableStatement {
			for _, decl := range node.AsVariableStatement().DeclarationList.AsVariableDeclarationList().Declarations.Nodes {
				if ast.IsBindingPattern(decl.Name()) {
					b.declarePatternCells(decl.Name())
					continue
				}
				cell := b.e.bindings[decl]
				if cell != nil && cell.lexical {
					b.allocate(cell)
				}
			}
		}
		if node.Kind == ast.KindClassDeclaration {
			b.allocate(b.e.bindings[node])
		}
		if node.Kind == ast.KindFunctionDeclaration && node.Body() != nil {
			b.allocate(b.e.bindings[node])
		}
	}
	for _, node := range nodes {
		if node.Kind == ast.KindFunctionDeclaration && node.Body() != nil {
			cell := b.e.bindings[node]
			if cell != nil {
				b.emit(cell.name + ".init(" + b.function(node) + ")")
			}
		}
	}
}
func (b *machineBuilder) function(node *ast.Node) string {
	if fn := b.e.nativeFunctions[node]; fn != nil {
		return b.nativeFunction(node, fn)
	}
	child := b.e.newMachine(node, b, ast.HasSyntacticModifier(node, ast.ModifierFlagsAsync))
	if node.Kind == ast.KindArrowFunction {
		child.concrete = b.concrete
		child.declaring = b.declaring
		child.receiver = b.receiver
	}
	for index, param := range node.Parameters() {
		p := param.AsParameterDeclaration()
		if param.Name().Kind != ast.KindIdentifier && !ast.IsBindingPattern(param.Name()) {
			b.e.fail(param, "rest and destructuring parameters are not supported")
			continue
		}
		if ast.IsBindingPattern(param.Name()) {
			value := child.temp(fmt.Sprintf("tsArg(args,%d)", index))
			if p.Initializer != nil {
				yes, end := child.block(), child.block()
				child.emit(fmt.Sprintf("if tsIsUndefined(%s){m.pc=%d}else{m.pc=%d};return", value, yes, end))
				child.current = yes
				fallback := child.expression(p.Initializer)
				child.emit(value + "=" + fallback)
				child.jump(end)
				child.current = end
			}
			child.bindPattern(param.Name(), value, true)
			continue
		}
		cell := b.e.bindings[param]
		if cell == nil {
			continue
		}
		if p.DotDotDotToken != nil {
			child.emit(fmt.Sprintf("%s.init(tsRestArgs(args,%d))", cell.name, index))
		} else {
			child.emit(fmt.Sprintf("%s.init(tsArg(args,%d))", cell.name, index))
		}
		if p.Initializer != nil {
			initialize, end := child.block(), child.block()
			child.emit(fmt.Sprintf("if tsIsUndefined(%s.get()) {m.pc=%d} else {m.pc=%d}; return", cell.name, initialize, end))
			child.current = initialize
			value := child.expression(p.Initializer)
			child.emit(cell.name + ".init(" + value + ")")
			child.jump(end)
			child.current = end
		}
		child.validateParameter(param)
	}
	if node.Body().Kind == ast.KindBlock {
		child.statements(node.Body().AsBlock().Statements.Nodes)
		child.abrupt("return", "tsU", 0, 0)
	} else {
		value := child.expression(node.Body())
		child.abrupt("return", value, 0, 0)
	}
	return child.finish()
}
func (b *machineBuilder) finish() string {
	var out strings.Builder
	captures := []string{}
	args := []string{}
	for _, cell := range b.captures {
		captures = append(captures, cell.name+" "+cell.pointerType())
		args = append(args, cell.name)
	}
	if b.module {
		out.WriteString("func(module *tsModule) {\n")
	} else {
		fmt.Fprintf(&out, "func(%s) *tsFunction { return tsFunc(func(args ...tsValue) tsValue {\n", strings.Join(captures, ","))
	}
	for _, cell := range b.locals {
		fmt.Fprintf(&out, "var %s %s\n_ = %s\n", cell.name, cell.pointerType(), cell.name)
		if !cell.lexical {
			initialized := !cell.isParameter()
			fmt.Fprintf(&out, "%s = %s\n", cell.name, b.bindingFactory(cell, initialized))
		}
	}
	for _, temp := range b.temps {
		fmt.Fprintf(&out, "var %s %s\n_ = %s\n", temp, b.tempType(temp), temp)
	}
	fmt.Fprintf(&out, "m := &tsMachine{loop:loop, async:%t}\n", b.async && !b.module)
	if b.module {
		out.WriteString("m.module = module\nmodule.machine = m\n")
	}
	if b.async && !b.module {
		out.WriteString("m.output = loop.promise()\n")
	}
	out.WriteString("m.step = func() { switch m.pc {\n")
	for index, block := range b.blocks {
		fmt.Fprintf(&out, "case %d:\n%s\n", index, strings.Join(block.lines, "\n"))
	}
	out.WriteString("default: panic(\"Invalid resume state\")\n} }\n")
	if b.module {
		out.WriteString("}\n")
		return out.String()
	}
	out.WriteString("m.resume()\n")
	if b.async {
		out.WriteString("return m.output\n")
	} else {
		out.WriteString("return m.result\n")
	}
	fmt.Fprintf(&out, "}) }(%s)", strings.Join(args, ","))
	return out.String()
}
