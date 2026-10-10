package goemit

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/microsoft/TypeScript/tsc/internal/ast"
)

// One chain shares a short-circuit destination. Parentheses terminate the chain,
// so (value?.field).next still performs an ordinary access on its result.
func (b *machineBuilder) optionalChain(node *ast.Node) string {
	return b.optionalChainMode(node, false)
}
func (b *machineBuilder) optionalChainMode(node *ast.Node, deleting bool) string {
	var links []*ast.Node
	base := node
	for ast.IsOptionalChain(base) {
		links = append(links, base)
		base = base.Expression()
	}
	value := b.expression(base)
	result := b.temp("tsU")
	if deleting {
		b.emit(result + "=true")
	}
	end := 0
	opens := 0
	if !b.direct {
		end = b.block()
	}
	for i := len(links) - 1; i >= 0; i-- {
		link := links[i]
		optional := false
		switch link.Kind {
		case ast.KindPropertyAccessExpression:
			optional = link.AsPropertyAccessExpression().QuestionDotToken != nil
		case ast.KindElementAccessExpression:
			optional = link.AsElementAccessExpression().QuestionDotToken != nil
		case ast.KindCallExpression:
			optional = link.AsCallExpression().QuestionDotToken != nil
		}
		if optional {
			if b.direct {
				b.emit("if !tsNullish(" + value + ") {")
				opens++
			} else {
				next := b.block()
				b.emit(fmt.Sprintf("if tsNullish(%s){m.pc=%d}else{m.pc=%d};return", value, end, next))
				b.current = next
			}
		}
		switch link.Kind {
		case ast.KindPropertyAccessExpression:
			if deleting && i == 0 {
				value = b.temp("tsDelete(" + value + "," + strconv.Quote(link.Name().Text()) + "," + fmt.Sprint(strictFunction(node)) + ")")
			} else {
				value = b.temp("tsGet(" + value + "," + strconv.Quote(link.Name().Text()) + ")")
			}
		case ast.KindElementAccessExpression:
			key := b.expression(link.AsElementAccessExpression().ArgumentExpression)
			if deleting && i == 0 {
				value = b.temp("tsDelete(" + value + "," + key + "," + fmt.Sprint(strictFunction(node)) + ")")
			} else {
				value = b.temp("tsGet(" + value + "," + key + ")")
			}
		case ast.KindCallExpression:
			call := link.AsCallExpression()
			spread := false
			for _, arg := range call.Arguments.Nodes {
				if arg.Kind == ast.KindSpreadElement {
					spread = true
				}
			}
			if spread {
				value = b.temp("tsCall(" + value + "," + b.spreadArguments(call.Arguments.Nodes) + "...)")
			} else {
				args := []string{value}
				for index := range call.Arguments.Nodes {
					args = append(args, b.callArgument(call, index))
				}
				value = b.temp("tsCall(" + strings.Join(args, ",") + ")")
			}
		}
	}
	if deleting && node.Kind == ast.KindCallExpression {
		value = "true"
	}
	b.emit(result + "=" + value)
	if b.direct {
		for range opens {
			b.emit("}")
		}
	} else {
		b.jump(end)
		b.current = end
	}
	return result
}
