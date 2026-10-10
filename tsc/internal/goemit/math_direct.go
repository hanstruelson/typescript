package goemit

import (
	"github.com/microsoft/TypeScript/tsc/internal/ast"
	"strings"
)

func (e *emitter) planMathEffects() {
	e.directMath = !e.reflectiveFields
	var visit func(*ast.Node)
	visit = func(n *ast.Node) {
		if n.Kind == ast.KindIdentifier && n.Text() == "Math" && e.binding(n) == nil {
			p := n.Parent
			if p == nil || p.Kind != ast.KindPropertyAccessExpression || p.AsPropertyAccessExpression().Expression != n {
				e.directMath = false
			} else {
				at := p.Parent
				if at != nil && (at.Kind == ast.KindDeleteExpression || at.Kind == ast.KindPrefixUnaryExpression || at.Kind == ast.KindPostfixUnaryExpression || at.Kind == ast.KindBinaryExpression && ast.IsAssignmentOperator(at.AsBinaryExpression().OperatorToken.Kind) && at.AsBinaryExpression().Left == p) {
					e.directMath = false
				}
			}
		}
		if n.Kind == ast.KindIdentifier && n.Text() == "globalThis" {
			e.directMath = false
		}
		n.ForEachChild(func(c *ast.Node) bool { visit(c); return false })
	}
	visit(e.file.AsNode())
}
func (b *machineBuilder) nativeMathCall(call *ast.CallExpression) (string, bool) {
	if !b.e.directMath || call.Expression.Kind != ast.KindPropertyAccessExpression {
		return "", false
	}
	p := call.Expression.AsPropertyAccessExpression()
	if p.Expression.Kind != ast.KindIdentifier || p.Expression.Text() != "Math" || b.e.binding(p.Expression) != nil {
		return "", false
	}
	name := p.Name().Text()
	target := ""
	switch name {
	case "abs", "acos", "acosh", "asin", "asinh", "atan", "atanh", "cbrt", "ceil", "cos", "cosh", "exp", "expm1", "floor", "log", "log1p", "log2", "log10", "sin", "sinh", "sqrt", "tan", "tanh", "trunc":
		target = "math." + strings.ToUpper(name[:1]) + name[1:]
	case "round":
		target = "tsMathRound"
	case "fround":
		target = "tsMathFround"
	case "pow":
		target = "tsPow"
	case "atan2":
		target = "math.Atan2"
	default:
		return "", false
	}
	for _, arg := range call.Arguments.Nodes {
		shape := b.e.primitive(arg)
		if arg.Kind == ast.KindSpreadElement || !shape.numeric() || shape.nulls != 0 {
			return "", false
		}
	}
	values := []string{}
	for _, arg := range call.Arguments.Nodes {
		value := b.expression(arg)
		if (primitive{kind: b.tempType(value)}).numeric() {
			values = append(values, "float64("+value+")")
		} else {
			values = append(values, "tsNumber(tsToPrimitive("+value+"))")
		}
	}
	arity := 1
	if name == "pow" || name == "atan2" {
		arity = 2
	}
	for len(values) < arity {
		values = append(values, "math.NaN()")
	}
	return b.typedTemp(target+"("+strings.Join(values[:arity], ",")+")", "float64"), true
}
