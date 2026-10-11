package goemit

import "github.com/microsoft/TypeScript/tsc/internal/ast"

// analyzeEarlyErrors performs the small set of grammar-time restrictions that
// materially change the generated program. These are checked before lowering;
// they never become runtime branches in valid code.
func (e *emitter) analyzeEarlyErrors() {
	fileStrict := false
	if e.file != nil {
		for _, statement := range e.file.Statements.Nodes {
			if statement.Kind == ast.KindImportDeclaration || statement.Kind == ast.KindExportDeclaration || statement.Kind == ast.KindExportAssignment {
				fileStrict = true
			}
			if statement.Kind == ast.KindExpressionStatement && statement.AsExpressionStatement().Expression.Kind == ast.KindStringLiteral && statement.AsExpressionStatement().Expression.Text() == "use strict" {
				fileStrict = true
				break
			}
		}
	}
	var visit func(*ast.Node, bool, bool)
	visit = func(node *ast.Node, strict, inClass bool) {
		if node == nil {
			return
		}
		if node.Kind == ast.KindClassDeclaration || node.Kind == ast.KindClassExpression {
			inClass = true
		}
		if ast.IsFunctionLike(node) {
			strict = strict || inClass
			if body := node.Body(); body != nil && body.Kind == ast.KindBlock {
				for _, statement := range body.AsBlock().Statements.Nodes {
					if statement.Kind == ast.KindExpressionStatement && statement.AsExpressionStatement().Expression.Kind == ast.KindStringLiteral && statement.AsExpressionStatement().Expression.Text() == "use strict" {
						strict = true
						break
					}
					break
				}
			}
			if ast.GetFunctionFlags(node)&ast.FunctionFlagsGenerator != 0 {
				for _, parameter := range node.Parameters() {
					if initializer := parameter.Initializer(); initializer != nil && containsYield(initializer) {
						e.fail(initializer, "yield is not allowed in a generator parameter initializer")
					}
				}
			}
		}
		switch node.Kind {
		case ast.KindBinaryExpression:
			binary := node.AsBinaryExpression()
			if strict && ast.IsAssignmentOperator(binary.OperatorToken.Kind) && binary.Left.Kind == ast.KindIdentifier && binary.Left.Text() == "arguments" {
				e.fail(binary.Left, "cannot assign to arguments in strict mode")
			}
		case ast.KindPrefixUnaryExpression:
			if strict && (node.AsPrefixUnaryExpression().Operator == ast.KindPlusPlusToken || node.AsPrefixUnaryExpression().Operator == ast.KindMinusMinusToken) && isArguments(node.AsPrefixUnaryExpression().Operand) {
				e.fail(node, "cannot update arguments in strict mode")
			}
		case ast.KindPostfixUnaryExpression:
			if strict && isArguments(node.AsPostfixUnaryExpression().Operand) {
				e.fail(node, "cannot update arguments in strict mode")
			}
		case ast.KindDeleteExpression:
			if strict && isArguments(node.AsDeleteExpression().Expression) {
				e.fail(node, "cannot delete arguments in strict mode")
			}
		case ast.KindPropertyAccessExpression:
			if node.Name() != nil && ast.IsPrivateIdentifier(node.Name()) && !inClass {
				e.fail(node, "private identifiers are only accessible inside their declaring class")
			}
		case ast.KindElementAccessExpression:
			// Element access cannot name a private identifier, so it remains dynamic.
		}
		node.ForEachChild(func(child *ast.Node) bool { visit(child, strict, inClass); return false })
	}
	visit(e.file.AsNode(), fileStrict, false)
}

func isArguments(node *ast.Node) bool {
	return node != nil && node.Kind == ast.KindIdentifier && node.Text() == "arguments"
}

func containsYield(node *ast.Node) bool {
	found := false
	var visit func(*ast.Node)
	visit = func(current *ast.Node) {
		if current == nil || found {
			return
		}
		if current.Kind == ast.KindYieldExpression {
			found = true
			return
		}
		if current != node && ast.IsFunctionLike(current) {
			return
		}
		current.ForEachChild(func(child *ast.Node) bool { visit(child); return false })
	}
	visit(node)
	return found
}
