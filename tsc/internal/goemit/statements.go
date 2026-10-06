package goemit

import (
	"fmt"
	"strconv"

	"github.com/microsoft/TypeScript/tsc/internal/ast"
)

func (b *machineBuilder) declarations(list *ast.Node) {
	for _, decl := range list.AsVariableDeclarationList().Declarations.Nodes {
		if ast.IsBindingPattern(decl.Name()) {
			value := b.expression(decl.AsVariableDeclaration().Initializer)
			b.bindPattern(decl.Name(), value, true)
			continue
		}
		cell := b.e.bindings[decl]
		if cell == nil {
			continue
		}
		d := decl.AsVariableDeclaration()
		value := "tsU"
		if d.Initializer != nil {
			if cell.primitive.numeric() && cell.primitive.nulls == 0 {
				value = b.numericExpression(d.Initializer, cell.primitive)
			} else {
				value = b.expression(d.Initializer)
			}
		}
		native := cell.primitive.kind != "" && b.tempType(value) == cell.primitive.goType()
		init, set := "init", "set"
		if native {
			init, set = "initNative", "setNative"
		}
		if cell.lexical {
			if d.Initializer != nil && cell.primitive.kind != "" && !native {
				value = fmt.Sprintf("tsBoundary(%s,%q,%d,%t)", value, cell.primitive.kind, cell.primitive.nulls, b.e.coerce)
			}
			b.emit(cell.name + "." + init + "(" + value + ")")
		} else if d.Initializer != nil {
			b.emit(cell.name + "." + set + "(" + value + ")")
		}
		if b.module && ast.HasSyntacticModifier(list.Parent, ast.ModifierFlagsExport) {
			b.publish(decl.Name().Text(), cell.name+".get()", cell.name)
		}
	}
}
func (b *machineBuilder) statement(node *ast.Node) { b.statementLabel(node, "") }
func (b *machineBuilder) statementLabel(node *ast.Node, label string) {
	if node == nil || ast.HasSyntacticModifier(node, ast.ModifierFlagsAmbient) {
		return
	}
	switch node.Kind {
	case ast.KindEmptyStatement, ast.KindInterfaceDeclaration, ast.KindTypeAliasDeclaration:
	case ast.KindEnumDeclaration:
		b.enumDeclaration(node)
	case ast.KindClassDeclaration:
		b.classDeclaration(node)
	case ast.KindFunctionDeclaration:
		if b.module && node.Body() != nil && ast.HasSyntacticModifier(node, ast.ModifierFlagsExport) {
			name := "default"
			if !ast.HasSyntacticModifier(node, ast.ModifierFlagsDefault) {
				name = node.Name().Text()
			}
			cell := b.e.bindings[node]
			if cell != nil {
				b.publish(name, cell.name+".get()", cell.name)
			}
		}
	case ast.KindExportDeclaration:
		b.exportDeclaration(node)
	case ast.KindExportAssignment:
		if !b.module {
			b.e.fail(node, "exports require module bundling")
			return
		}
		n := node.AsExportAssignment()
		if n.IsExportEquals {
			b.e.fail(node, "export equals is not supported")
			return
		}
		value := b.expression(n.Expression)
		b.publish("default", value, "nil")
	case ast.KindImportDeclaration:
		b.importDeclaration(node)
	case ast.KindBlock:
		b.statements(node.AsBlock().Statements.Nodes)
	case ast.KindVariableStatement:
		b.declarations(node.AsVariableStatement().DeclarationList)
	case ast.KindExpressionStatement:
		b.expression(node.AsExpressionStatement().Expression)
	case ast.KindReturnStatement:
		value := "tsU"
		if node.AsReturnStatement().Expression != nil {
			value = b.expression(node.AsReturnStatement().Expression)
		}
		b.abrupt("return", b.returnValue(value), 0, 0)
		b.current = b.block()
	case ast.KindThrowStatement:
		value := b.expression(node.AsThrowStatement().Expression)
		b.abrupt("throw", value, 0, 0)
		b.current = b.block()
	case ast.KindIfStatement:
		n := node.AsIfStatement()
		value := b.expression(n.Expression)
		yes, no, end := b.block(), b.block(), b.block()
		b.emit(fmt.Sprintf("if tsTruthy(%s) {m.pc=%d} else {m.pc=%d}; return", value, yes, no))
		b.current = yes
		b.statement(n.ThenStatement)
		b.jump(end)
		b.current = no
		b.statement(n.ElseStatement)
		b.jump(end)
		b.current = end
	case ast.KindWhileStatement, ast.KindDoStatement:
		var condition, body *ast.Node
		if node.Kind == ast.KindWhileStatement {
			n := node.AsWhileStatement()
			condition = n.Expression
			body = n.Statement
		} else {
			n := node.AsDoStatement()
			condition = n.Expression
			body = n.Statement
		}
		test, start, end := b.block(), b.block(), b.block()
		if node.Kind == ast.KindDoStatement {
			b.jump(start)
		} else {
			b.jump(test)
		}
		b.loops = append(b.loops, loopTarget{end, test, b.depth, label})
		b.current = start
		b.statement(body)
		b.jump(test)
		b.current = test
		value := b.expression(condition)
		b.emit(fmt.Sprintf("if tsTruthy(%s) {m.pc=%d} else {m.pc=%d}; return", value, start, end))
		b.loops = b.loops[:len(b.loops)-1]
		b.current = end
	case ast.KindForStatement:
		b.forStatement(node, label)
	case ast.KindForInStatement:
		b.forInStatement(node, label)
	case ast.KindSwitchStatement:
		b.switchStatement(node, label)
	case ast.KindForOfStatement:
		b.forOfStatement(node, label)
	case ast.KindBreakStatement, ast.KindContinueStatement:
		targetLabel := ""
		if node.Kind == ast.KindBreakStatement {
			if n := node.AsBreakStatement().Label; n != nil {
				targetLabel = n.Text()
			}
		} else {
			if n := node.AsContinueStatement().Label; n != nil {
				targetLabel = n.Text()
			}
		}
		found := false
		for i := len(b.loops) - 1; i >= 0; i-- {
			target := b.loops[i]
			if node.Kind == ast.KindContinueStatement && target.continuePC < 0 {
				continue
			}
			if targetLabel != "" && target.label != targetLabel {
				continue
			}
			pc := target.breakPC
			if node.Kind == ast.KindContinueStatement {
				pc = target.continuePC
			}
			b.abrupt("jump", "tsU", pc, target.depth)
			found = true
			break
		}
		if !found {
			b.e.fail(node, "break/continue requires a supported enclosing loop")
		}
		b.current = b.block()
	case ast.KindLabeledStatement:
		n := node.AsLabeledStatement()
		b.statementLabel(n.Statement, n.Label.Text())
	case ast.KindTryStatement:
		b.tryStatement(node)
	default:
		b.e.fail(node, "unsupported statement "+node.Kind.String())
	}
}
func (b *machineBuilder) forStatement(node *ast.Node, label string) {
	n := node.AsForStatement()
	var iteration []*binding
	if n.Initializer != nil {
		if n.Initializer.Kind == ast.KindVariableDeclarationList {
			for _, decl := range n.Initializer.AsVariableDeclarationList().Declarations.Nodes {
				cell := b.e.bindings[decl]
				if cell != nil && cell.lexical {
					b.allocate(cell)
					iteration = append(iteration, cell)
				}
			}
			b.declarations(n.Initializer)
		} else {
			b.expression(n.Initializer)
		}
	}
	clone := func() {
		for _, cell := range iteration {
			clone := "tsClone"
			if cell.primitive.kind != "" {
				clone = "tsCloneTyped"
			}
			b.emit(cell.name + "=" + clone + "(" + cell.name + ")")
		}
	}
	// The initializer environment is distinct from the first iteration's bindings.
	clone()
	test, body, increment, end := b.block(), b.block(), b.block(), b.block()
	b.jump(test)
	b.loops = append(b.loops, loopTarget{end, increment, b.depth, label})
	b.current = test
	value := "true"
	if n.Condition != nil {
		value = b.expression(n.Condition)
	}
	b.emit(fmt.Sprintf("if tsTruthy(%s) {m.pc=%d} else {m.pc=%d}; return", value, body, end))
	b.current = body
	b.statement(n.Statement)
	b.jump(increment)
	b.current = increment
	clone()
	if n.Incrementor != nil {
		b.expression(n.Incrementor)
	}
	b.jump(test)
	b.loops = b.loops[:len(b.loops)-1]
	b.current = end
}
func (b *machineBuilder) forOfStatement(node *ast.Node, label string) {
	n := node.AsForInOrOfStatement()
	if n.AwaitModifier != nil {
		b.e.fail(node, "for-await-of is not supported yet")
		return
	}
	if n.Initializer.Kind != ast.KindVariableDeclarationList || len(n.Initializer.AsVariableDeclarationList().Declarations.Nodes) != 1 {
		b.e.fail(node, "for-of requires one identifier declaration")
		return
	}
	decl := n.Initializer.AsVariableDeclarationList().Declarations.Nodes[0]
	cell := b.e.bindings[decl]
	pattern := ast.IsBindingPattern(decl.Name())
	if cell == nil && !pattern {
		return
	}
	lexical := n.Initializer.AsVariableDeclarationList().Flags&ast.NodeFlagsBlockScoped != 0
	// The loop declaration is in scope (uninitialized) while evaluating its RHS.
	if pattern {
		b.declarePatternCells(decl.Name())
	} else if cell.lexical {
		b.allocate(cell)
	}
	iterable := b.expression(n.Expression)
	iterator := b.typedTemp("tsIterate("+iterable+")", "*tsIterator")
	test, body, end := b.block(), b.block(), b.block()
	b.jump(test)
	b.loops = append(b.loops, loopTarget{end, test, b.depth, label})
	b.current = test
	b.emit(fmt.Sprintf("if %s.next() {m.pc=%d} else {m.pc=%d}; return", iterator, body, end))
	b.current = body
	if pattern {
		if lexical {
			b.declarePatternCells(decl.Name())
		}
		b.bindPattern(decl.Name(), iterator+".value", lexical)
	} else if cell.lexical {
		b.allocate(cell)
		b.emit(cell.name + ".init(" + iterator + ".value)")
	} else {
		b.emit(cell.name + ".set(" + iterator + ".value)")
	}
	b.statement(n.Statement)
	b.jump(test)
	b.loops = b.loops[:len(b.loops)-1]
	b.current = end
}
func (b *machineBuilder) tryStatement(node *ast.Node) {
	n := node.AsTryStatement()
	start, end := b.block(), b.block()
	catch, finally := -1, -1
	if n.CatchClause != nil {
		catch = b.block()
	}
	if n.FinallyBlock != nil {
		finally = b.block()
	}
	depth := b.depth
	b.emit(fmt.Sprintf("m.handlers=append(m.handlers,tsHandler{catch:%d,finally:%d,end:%d}); m.pc=%d; return", catch, finally, end, start))
	b.depth++
	b.current = start
	b.statement(n.TryBlock)
	b.abrupt("jump", "tsU", end, depth)
	if catch >= 0 {
		b.current = catch
		c := n.CatchClause.AsCatchClause()
		if c.VariableDeclaration != nil {
			cell := b.e.bindings[c.VariableDeclaration]
			b.allocate(cell)
			if cell != nil {
				b.emit(cell.name + ".init(m.thrown)")
			}
		}
		b.statement(c.Block)
		b.abrupt("jump", "tsU", end, depth)
	}
	if finally >= 0 {
		b.current = finally
		b.statement(n.FinallyBlock)
		b.emit("m.endFinally(); return")
	}
	b.depth--
	b.current = end
}
func (b *machineBuilder) importDeclaration(node *ast.Node) {
	n := node.AsImportDeclaration()
	module := n.ModuleSpecifier.Text()
	if b.e.targets[node] != "" {
		b.localImport(node)
		return
	}
	if n.ImportClause == nil {
		b.e.fail(node, "unresolved side-effect import "+module)
		return
	}
	c := n.ImportClause.AsImportClause()
	if c.PhaseModifier == ast.KindTypeKeyword {
		return
	}
	if module != "node:fs/promises" && module != "fs/promises" {
		b.e.fail(node, "only readFile from node:fs/promises is supported")
		return
	}
	if c.Name() != nil || c.NamedBindings == nil || c.NamedBindings.Kind != ast.KindNamedImports {
		b.e.fail(node, "readFile requires a named import")
		return
	}
	for _, specifier := range c.NamedBindings.AsNamedImports().Elements.Nodes {
		s := specifier.AsImportSpecifier()
		if s.IsTypeOnly {
			continue
		}
		name := specifier.Name().Text()
		if s.PropertyName != nil {
			name = s.PropertyName.Text()
		}
		if name != "readFile" {
			b.e.fail(specifier, "unsupported fs/promises import "+strconv.Quote(name))
			continue
		}
		b.e.imports[specifier] = "readFile"
	}
}

func (b *machineBuilder) switchStatement(node *ast.Node, label string) {
	n := node.AsSwitchStatement()
	value := b.expression(n.Expression)
	clauses := n.CaseBlock.AsCaseBlock().Clauses.Nodes
	end := b.block()
	states := []int{}
	defaultPC := end
	all := []*ast.Node{}
	for _, clause := range clauses {
		states = append(states, b.block())
		if clause.Kind == ast.KindDefaultClause {
			defaultPC = states[len(states)-1]
			all = append(all, clause.AsCaseOrDefaultClause().Statements.Nodes...)
		} else {
			all = append(all, clause.AsCaseOrDefaultClause().Statements.Nodes...)
		}
	}
	b.enter(all)
	for i, clause := range clauses {
		if clause.Kind == ast.KindDefaultClause {
			continue
		}
		test := b.expression(clause.AsCaseOrDefaultClause().Expression)
		next := b.block()
		b.emit(fmt.Sprintf("if tsStrictEqual(%s,%s){m.pc=%d}else{m.pc=%d};return", value, test, states[i], next))
		b.current = next
	}
	b.jump(defaultPC)
	b.loops = append(b.loops, loopTarget{end, -1, b.depth, label})
	for i, clause := range clauses {
		b.current = states[i]
		var statements []*ast.Node
		if clause.Kind == ast.KindDefaultClause {
			statements = clause.AsCaseOrDefaultClause().Statements.Nodes
		} else {
			statements = clause.AsCaseOrDefaultClause().Statements.Nodes
		}
		for _, statement := range statements {
			b.statement(statement)
		}
		next := end
		if i+1 < len(states) {
			next = states[i+1]
		}
		b.jump(next)
	}
	b.loops = b.loops[:len(b.loops)-1]
	b.current = end
}
