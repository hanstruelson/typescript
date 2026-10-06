package goemit

import (
	"fmt"
	"strings"

	"github.com/microsoft/TypeScript/tsc/internal/ast"
)

// Synchronous class bodies with structured control flow do not need a resume
// machine. Keep that machinery only for awaits and protected-region lowering.
func directClassBody(node *ast.Node) bool {
	if node == nil {
		return true
	}
	if ast.HasSyntacticModifier(node, ast.ModifierFlagsAsync) {
		return false
	}
	supported := true
	var visit func(*ast.Node)
	visit = func(n *ast.Node) {
		if n != node && (n.Kind == ast.KindArrowFunction || n.Kind == ast.KindFunctionExpression || n.Kind == ast.KindFunctionDeclaration) {
			return
		}
		switch n.Kind {
		case ast.KindAwaitExpression, ast.KindTryStatement, ast.KindForStatement, ast.KindForOfStatement, ast.KindForInStatement, ast.KindWhileStatement, ast.KindDoStatement, ast.KindBreakStatement, ast.KindContinueStatement, ast.KindLabeledStatement, ast.KindSwitchStatement:
			supported = false
		}
		n.ForEachChild(func(child *ast.Node) bool { visit(child); return false })
	}
	visit(node)
	return supported
}
func (b *machineBuilder) directStatements(nodes []*ast.Node) {
	b.enter(nodes)
	for _, node := range nodes {
		b.directStatement(node)
	}
}
func (b *machineBuilder) directStatement(node *ast.Node) {
	if node == nil || ast.HasSyntacticModifier(node, ast.ModifierFlagsAmbient) {
		return
	}
	switch node.Kind {
	case ast.KindBlock:
		b.directStatements(node.AsBlock().Statements.Nodes)
	case ast.KindExpressionStatement:
		b.expression(node.AsExpressionStatement().Expression)
	case ast.KindVariableStatement:
		b.declarations(node.AsVariableStatement().DeclarationList)
	case ast.KindReturnStatement:
		if b.nativeArrayResult.kind != "" {
			expr := node.AsReturnStatement().Expression
			if expr == nil {
				b.emit("panic(\"Missing typed array return\")")
			} else {
				value := b.denseStorage(expr, b.nativeArrayResult)
				b.emit("return " + value)
			}
			return
		}
		value := "tsU"
		expr := node.AsReturnStatement().Expression
		if expr != nil {
			p := b.e.returnPrimitive(b.owner)
			if b.nativeResult.kind != "" {
				p = b.nativeResult
			}
			if p.numeric() && p.nulls == 0 {
				value = b.numericExpression(expr, p)
			} else {
				value = b.expression(expr)
			}
		}
		if b.nativeReturn != "" {
			if b.nativeReturn == "tsValue" {
				value = b.returnValue(value)
			} else if b.nativeResult.kind != "" {
				value = b.comparisonBoundary(value, b.nativeResult)
			} else if b.tempType(value) != b.nativeReturn {
				value = "tsNative[" + b.nativeReturn + "](" + b.returnValue(value) + ")"
			}
			b.emit("return " + value)
		} else {
			b.emit("return " + b.returnValue(value))
		}
	case ast.KindThrowStatement:
		value := b.expression(node.AsThrowStatement().Expression)
		b.emit("panic(tsThrown{value:" + value + "})")
	case ast.KindIfStatement:
		n := node.AsIfStatement()
		condition := b.expression(n.Expression)
		b.emit("if tsTruthy(" + condition + ") {")
		b.directStatement(n.ThenStatement)
		b.emit("} else {")
		b.directStatement(n.ElseStatement)
		b.emit("}")
	case ast.KindEmptyStatement, ast.KindFunctionDeclaration, ast.KindInterfaceDeclaration, ast.KindTypeAliasDeclaration:
	case ast.KindEnumDeclaration:
		b.enumDeclaration(node)
	case ast.KindClassDeclaration:
		b.classDeclaration(node)
	default:
		b.e.fail(node, "unsupported synchronous class statement "+node.Kind.String())
	}
}
func (b *machineBuilder) finishDirect() string {
	var out strings.Builder
	for _, cell := range b.locals {
		fmt.Fprintf(&out, "var %s %s;_=%s\n", cell.name, cell.pointerType(), cell.name)
		if !cell.lexical {
			initialized := !cell.isParameter()
			fmt.Fprintf(&out, "%s=%s\n", cell.name, b.bindingFactory(cell, initialized))
		}
	}
	for _, temp := range b.temps {
		fmt.Fprintf(&out, "var %s %s;_=%s\n", temp, b.tempType(temp), temp)
	}
	for _, line := range b.blocks[0].lines {
		out.WriteString(line + "\n")
	}
	if b.nativeReturn != "" {
		if b.nativeReturn == "tsValue" {
			out.WriteString("return tsU\n")
		} else if strings.HasPrefix(b.nativeReturn, "tsOptional[") {
			out.WriteString("return " + b.nativeReturn + "{tag:2}\n")
		} else {
			out.WriteString("panic(\"Missing typed return value\")\n")
		}
	} else {
		out.WriteString("return tsU\n")
	}
	return out.String()
}
