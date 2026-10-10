package goemit

import "github.com/microsoft/TypeScript/tsc/internal/ast"

func observesThis(node *ast.Node) bool {
	found := false
	var visit func(*ast.Node)
	visit = func(n *ast.Node) {
		if n != node && ast.IsFunctionLike(n) && n.Kind != ast.KindArrowFunction {
			return
		}
		if n.Kind == ast.KindThisKeyword {
			found = true
		}
		n.ForEachChild(func(c *ast.Node) bool { visit(c); return false })
	}
	visit(node)
	return found
}

func strictFunction(node *ast.Node) bool {
	for at := node; at != nil; at = at.Parent {
		if at.Kind == ast.KindClassDeclaration || at.Kind == ast.KindClassExpression {
			return true
		}
		var statements []*ast.Node
		if at.Kind == ast.KindSourceFile {
			statements = at.AsSourceFile().Statements.Nodes
		}
		if ast.IsFunctionLike(at) && at.Body() != nil && at.Body().Kind == ast.KindBlock {
			statements = at.Body().AsBlock().Statements.Nodes
		}
		for _, statement := range statements {
			if statement.Kind != ast.KindExpressionStatement {
				break
			}
			expression := statement.AsExpressionStatement().Expression
			if expression.Kind != ast.KindStringLiteral {
				break
			}
			if expression.Text() == "use strict" {
				return true
			}
		}
		for _, statement := range statements {
			if statement.Kind == ast.KindImportDeclaration || statement.Kind == ast.KindExportDeclaration || ast.HasSyntacticModifier(statement, ast.ModifierFlagsExport) {
				return true
			}
		}
	}
	return false
}
