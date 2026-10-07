package goemit

import (
	"fmt"
	"github.com/microsoft/TypeScript/tsc/internal/ast"
	"strings"
)

func nativeCallbackBody(node *ast.Node) bool {
	if ast.HasSyntacticModifier(node, ast.ModifierFlagsAsync) {
		return false
	}
	valid := true
	var visit func(*ast.Node)
	visit = func(n *ast.Node) {
		if n != node && ast.IsFunctionLike(n) {
			return
		}
		if n.Kind == ast.KindAwaitExpression || n.Kind == ast.KindYieldExpression || n.Kind == ast.KindThisKeyword {
			valid = false
		}
		n.ForEachChild(func(child *ast.Node) bool { visit(child); return false })
	}
	visit(node)
	return valid
}

// Protected regions carry only the completion control in tsMachine. The return
// payload remains a concrete Go local, including nullable primitive results.
func (b *machineBuilder) finishNativeMachine() string {
	var out strings.Builder
	for _, cell := range b.locals {
		fmt.Fprintf(&out, "var %s %s;_=%s\n", cell.name, cell.pointerType(), cell.name)
		if !cell.lexical {
			fmt.Fprintf(&out, "%s=%s\n", cell.name, b.bindingFactory(cell, !cell.isParameter()))
		}
	}
	for _, temp := range b.temps {
		kind := b.tempType(temp)
		if temp == b.nativeOutput {
			kind = fmt.Sprintf("[%d]%s", b.nativeReturnCount+1, kind)
		}
		fmt.Fprintf(&out, "var %s %s;_=%s\n", temp, kind, temp)
	}
	out.WriteString("m:=&tsMachine{loop:loop};m.step=func(){switch m.pc{\n")
	for i, block := range b.blocks {
		fmt.Fprintf(&out, "case %d:\n%s\n", i, strings.Join(block.lines, "\n"))
	}
	out.WriteString("default:panic(\"Invalid native callback resume state\")}};m.resume();return " + b.nativeOutput + "[m.nativeSlot]\n")
	return out.String()
}
func (b *machineBuilder) endNativeMachine() {
	if b.nativeArrayResult.kind != "" {
		b.emit("panic(\"Missing typed array return\")")
		return
	}
	if b.nativeResult.kind == "" {
		b.emit(b.nativeOutput + "[0]=tsU")
		b.abrupt("return", "tsU", 0, 0)
	} else if b.nativeResult.nulls&2 != 0 {
		b.emit(b.nativeOutput + "[0]=" + b.nativeResult.elementGoType() + "{tag:2}")
		b.abrupt("return", "tsU", 0, 0)
	} else {
		b.emit("panic(\"Missing typed return value\")")
	}
}
