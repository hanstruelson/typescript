package goemit

import "github.com/microsoft/TypeScript/tsc/internal/ast"

func (e *emitter) analyzeDeleteEffects() {
	e.reflectiveFields = e.dynamicPrototypes
	e.deleteEffects = map[*classInfo]map[string]bool{}
	var visit func(*ast.Node)
	visit = func(node *ast.Node) {
		if node.Kind == ast.KindDeleteExpression {
			e.dynamicPrototypes = true
			operand := node.AsDeleteExpression().Expression
			for operand.Kind == ast.KindParenthesizedExpression {
				operand = operand.AsParenthesizedExpression().Expression
			}
			var receiver *ast.Node
			name := "*"
			if operand.Kind == ast.KindPropertyAccessExpression {
				receiver = operand.AsPropertyAccessExpression().Expression
				name = operand.Name().Text()
			}
			if operand.Kind == ast.KindElementAccessExpression {
				p := operand.AsElementAccessExpression()
				receiver = p.Expression
				if p.ArgumentExpression.Kind == ast.KindStringLiteral {
					name = p.ArgumentExpression.Text()
				}
			}
			if receiver != nil {
				builder := &machineBuilder{e: e}
				for at := node.Parent; at != nil; at = at.Parent {
					if at.Kind == ast.KindClassDeclaration || at.Kind == ast.KindClassExpression {
						builder.concrete = e.classes[at]
						break
					}
				}
				known := builder.undecoratedClassOf(receiver)
				for _, c := range e.classOrder {
					for _, target := range []*classInfo{c, c.statics} {
						if target == nil {
							continue
						}
						affects := known == nil
						for at := target; at != nil; at = at.parent {
							if at == known {
								affects = true
							}
						}
						if affects {
							if e.deleteEffects[target] == nil {
								e.deleteEffects[target] = map[string]bool{}
							}
							e.deleteEffects[target][name] = true
						}
					}
				}
			}
		}
		node.ForEachChild(func(child *ast.Node) bool { visit(child); return false })
	}
	visit(e.file.AsNode())
}

// A proof of initialized, unchanged own storage permits a plain native load.
// Unknown aliases/reflection invalidate the proof; spelling alone is insufficient.
func (b *machineBuilder) directOwnField(receiver *ast.Node, name string) (*classInfo, bool) {
	if b.e.reflectiveFields {
		return nil, false
	}
	c := b.undecoratedClassOf(receiver)
	if c == nil || ast.HasDecorators(c.node) {
		return nil, false
	}
	effects := b.e.deleteEffects[c]
	if effects["*"] || effects[name] {
		return nil, false
	}
	field := c.fields[name]
	if field == nil || field.Kind != ast.KindPropertyDeclaration || field.AsPropertyDeclaration().Initializer == nil {
		return nil, false
	}
	shape := b.e.primitive(field)
	if shape.kind == "" || shape.nulls != 0 {
		return nil, false
	}
	return c, true
}

func (e *emitter) dynamicFieldAccess(node *ast.Node) bool {
	if node == nil || e.deleteEffects == nil {
		return false
	}
	var receiver *ast.Node
	name := "*"
	switch node.Kind {
	case ast.KindPropertyAccessExpression:
		receiver = node.AsPropertyAccessExpression().Expression
		name = node.Name().Text()
	case ast.KindElementAccessExpression:
		receiver = node.AsElementAccessExpression().Expression
		key := node.AsElementAccessExpression().ArgumentExpression
		if key.Kind == ast.KindStringLiteral {
			name = key.Text()
		}
	default:
		return false
	}
	builder := &machineBuilder{e: e}
	for at := node.Parent; at != nil; at = at.Parent {
		if at.Kind == ast.KindClassDeclaration {
			builder.concrete = e.classes[at]
			break
		}
	}
	c := builder.undecoratedClassOf(receiver)
	if c == nil {
		return false
	}
	if e.reflectiveFields {
		return true
	}
	effects := e.deleteEffects[c]
	return effects["*"] || effects[name]
}
