package goemit

import (
	"fmt"
	"github.com/microsoft/TypeScript/tsc/internal/ast"
	"strconv"
)

func classMemberName(member *ast.Node) string {
	if member.Name().Kind == ast.KindComputedPropertyName {
		return fmt.Sprintf("@computed:%d", member.Pos())
	}
	return member.Name().Text()
}
func classKeyExpression(class, name string) string {
	if len(name) > 10 && name[:10] == "@computed:" {
		return "tsClassKey(" + class + "," + strconv.Quote(name) + ")"
	}
	return strconv.Quote(name)
}

const computedClassRuntime = `
func tsClassKey(class *tsClass,name string)string{for at:=class;at!=nil;at=at.parent{if key,ok:=at.computedKeys[name];ok{return key}};panic("Missing computed class key")}
`
