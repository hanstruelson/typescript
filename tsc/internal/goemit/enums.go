package goemit

import (
	"strconv"

	"github.com/microsoft/TypeScript/tsc/internal/ast"
)

// Enum members share one ordinary object. Numeric members additionally expose
// TypeScript's reverse mapping; string members only expose their forward name.
func (b *machineBuilder) enumDeclaration(node *ast.Node) {
	cell := b.e.bindings[node]
	if cell == nil {
		return
	}
	b.emit("if !" + cell.name + ".initialized {" + cell.name + ".init(tsNewObject())}")
	object := b.temp(cell.name + ".get()")
	previous := ""
	for _, member := range node.AsEnumDeclaration().Members.Nodes {
		value := "float64(0)"
		if init := member.AsEnumMember().Initializer; init != nil {
			value = b.expression(init)
		} else if previous != "" {
			value = "tsBinary(\"+\"," + previous + ",float64(1))"
		}
		value = b.temp(value)
		name := strconv.Quote(member.Name().Text())
		b.emit("tsSet(" + object + "," + name + "," + value + ")")
		b.emit("if tsIsNumeric(" + value + "){tsSet(" + object + "," + value + "," + name + ")}")
		previous = value
	}
	if b.module && ast.HasSyntacticModifier(node, ast.ModifierFlagsExport) {
		b.publish(node.Name().Text(), cell.name+".get()", cell.name)
	}
}
