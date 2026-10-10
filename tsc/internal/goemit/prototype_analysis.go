package goemit

import "github.com/microsoft/TypeScript/tsc/internal/ast"

// Direct class operations are valid only in a closed, non-reflective source.
// This deliberately sacrifices an optimization when an alias could escape;
// finding no literal '.prototype' alone is not sufficient evidence of safety.
func (e *emitter) analyzePrototypes() {
	var visit func(*ast.Node)
	visit = func(node *ast.Node) {
		if ast.HasSyntacticModifier(node, ast.ModifierFlagsExport) {
			e.dynamicPrototypes = true
		}
		switch node.Kind {
		case ast.KindGetAccessor, ast.KindSetAccessor:
			e.dynamicPrototypes = true
		case ast.KindComputedPropertyName:
			if node.Parent != nil && node.Parent.Kind != ast.KindObjectLiteralExpression {
				e.dynamicPrototypes = true
			}
		case ast.KindImportDeclaration, ast.KindExportDeclaration, ast.KindElementAccessExpression:
			if node.Kind != ast.KindElementAccessExpression || node.Parent == nil || node.Parent.Kind != ast.KindDeleteExpression {
				e.dynamicPrototypes = true
			}
		case ast.KindPropertyAccessExpression:
			name := node.Name().Text()
			property := node.AsPropertyAccessExpression()
			if e.knownClassReceiver(property.Expression) && !e.closedOwnField(node) {
				parent := node.Parent
				if parent == nil || parent.Kind != ast.KindCallExpression || parent.AsCallExpression().Expression != node {
					e.dynamicPrototypes = true
				}
			}
			switch name {
			case "prototype", "__proto__", "getPrototypeOf", "setPrototypeOf", "defineProperty", "defineProperties", "create":
				e.dynamicPrototypes = true
			}
		case ast.KindNewExpression:
			if e.classReference(node.AsNewExpression().Expression) == nil {
				e.dynamicPrototypes = true
			}
		case ast.KindBinaryExpression:
			expression := node.AsBinaryExpression()
			if ast.IsAssignmentOperator(expression.OperatorToken.Kind) && expression.Left.Kind == ast.KindPropertyAccessExpression {
				if !e.closedOwnField(expression.Left) {
					e.dynamicPrototypes = true
				}
			}
		case ast.KindPrefixUnaryExpression:
			expression := node.AsPrefixUnaryExpression()
			if (expression.Operator == ast.KindPlusPlusToken || expression.Operator == ast.KindMinusMinusToken) && expression.Operand.Kind == ast.KindPropertyAccessExpression && !e.closedOwnField(expression.Operand) {
				e.dynamicPrototypes = true
			}
		case ast.KindPostfixUnaryExpression:
			expression := node.AsPostfixUnaryExpression()
			if expression.Operand.Kind == ast.KindPropertyAccessExpression && !e.closedOwnField(expression.Operand) {
				e.dynamicPrototypes = true
			}
		case ast.KindCallExpression:
			expression := node.AsCallExpression().Expression
			if expression.Kind == ast.KindPropertyAccessExpression {
				receiver := expression.AsPropertyAccessExpression().Expression
				safe := receiver.Kind == ast.KindThisKeyword || receiver.Kind == ast.KindSuperKeyword || e.classReference(receiver) != nil
				if receiver.Kind == ast.KindIdentifier {
					switch receiver.Text() {
					case "console", "Math", "Number", "String", "BigInt", "JSON", "Object":
						safe = true
					}
					if declaration := e.reference(receiver); declaration != nil && declaration.Kind == ast.KindVariableDeclaration {
						initial := declaration.AsVariableDeclaration().Initializer
						if initial != nil && initial.Kind == ast.KindNewExpression && e.classReference(initial.AsNewExpression().Expression) != nil {
							safe = true
						}
					}
				}
				if !safe {
					e.dynamicPrototypes = true
				}
			}
			if expression.Kind == ast.KindIdentifier {
				// Unknown calls can receive a constructor or instance and mutate it.
				// Arithmetic intrinsics cannot modify their object arguments.
				switch expression.Text() {
				case "Number", "String", "BigInt":
				default:
					e.dynamicPrototypes = true
				}
			}
		}
		node.ForEachChild(func(child *ast.Node) bool { visit(child); return false })
	}
	visit(e.file.AsNode())
}

// A declared field is an own property and cannot be shadowed by a prototype.
// Method writes and unknown receivers invalidate direct dispatch, including
// writes through `this`; spelling only '.prototype' would miss these changes.
func (e *emitter) closedOwnField(node *ast.Node) bool {
	p := node.AsPropertyAccessExpression()
	var class *classInfo
	static := false
	if p.Expression.Kind == ast.KindThisKeyword {
		for at := node.Parent; at != nil; at = at.Parent {
			if at.Kind == ast.KindMethodDeclaration {
				static = ast.HasSyntacticModifier(at, ast.ModifierFlagsStatic)
			}
			if at.Kind == ast.KindClassDeclaration || at.Kind == ast.KindClassExpression {
				class = e.classes[at]
				break
			}
		}
	} else if p.Expression.Kind == ast.KindIdentifier {
		if c := e.classReference(p.Expression); c != nil {
			class = c
			static = true
		} else if declaration := e.reference(p.Expression); declaration != nil && declaration.Kind == ast.KindVariableDeclaration {
			initial := declaration.AsVariableDeclaration().Initializer
			if initial != nil && initial.Kind == ast.KindNewExpression {
				class = e.classReference(initial.AsNewExpression().Expression)
			}
		}
	}
	if class == nil {
		return false
	}
	for _, member := range classMembers(class.node).Nodes {
		if member.Kind == ast.KindPropertyDeclaration && member.Name() != nil && classMemberName(member) == p.Name().Text() && ast.HasSyntacticModifier(member, ast.ModifierFlagsStatic) == static {
			return true
		}
	}
	return false
}

func (e *emitter) knownClassReceiver(node *ast.Node) bool {
	if node.Kind == ast.KindThisKeyword {
		for at := node.Parent; at != nil; at = at.Parent {
			if at.Kind == ast.KindClassDeclaration || at.Kind == ast.KindClassExpression {
				return true
			}
		}
	}
	if e.classReference(node) != nil {
		return true
	}
	if node.Kind != ast.KindIdentifier {
		return false
	}
	declaration := e.reference(node)
	if declaration == nil {
		return false
	}
	if declaration.Kind == ast.KindVariableDeclaration {
		variable := declaration.AsVariableDeclaration()
		if variable.Initializer != nil && variable.Initializer.Kind == ast.KindNewExpression && e.classReference(variable.Initializer.AsNewExpression().Expression) != nil {
			return true
		}
		if variable.Type != nil && variable.Type.Kind == ast.KindTypeReference {
			return e.classReference(variable.Type.AsTypeReferenceNode().TypeName) != nil
		}
	}
	if declaration.Kind == ast.KindParameter {
		parameter := declaration.AsParameterDeclaration()
		if parameter.Type != nil && parameter.Type.Kind == ast.KindTypeReference {
			return e.classReference(parameter.Type.AsTypeReferenceNode().TypeName) != nil
		}
	}
	return false
}
