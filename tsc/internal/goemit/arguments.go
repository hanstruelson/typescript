package goemit

import (
	"fmt"
	"github.com/microsoft/TypeScript/tsc/internal/ast"
	"strconv"
)

func observesArguments(node *ast.Node) bool {
	if node == nil || node.Kind == ast.KindArrowFunction {
		return false
	}
	found := false
	var visit func(*ast.Node)
	visit = func(n *ast.Node) {
		if n != node && ast.IsFunctionLike(n) && n.Kind != ast.KindArrowFunction {
			return
		}
		if n.Kind == ast.KindIdentifier && n.Text() == "arguments" {
			found = true
		}
		n.ForEachChild(func(child *ast.Node) bool { visit(child); return false })
	}
	visit(node)
	return found
}
func (b *machineBuilder) prepareArguments(node *ast.Node) {
	if !observesArguments(node) {
		return
	}
	strict := strictFunction(node)
	b.argumentsObject = b.typedTemp(fmt.Sprintf("tsArgumentsObject(args,%t)", strict), "*tsObject")
	b.argumentsCell = b.typedTemp("tsBinding(true,false)", "*tsCell")
	b.emit(b.argumentsCell + ".init(" + b.argumentsObject + ")")
	if !strict && b.dynamicReceiver {
		b.emit(b.argumentsObject + ".values[\"callee\"]=tsFunctionValue(function)")
	}
}
func (b *machineBuilder) initializeArguments(node *ast.Node) {
	if b.argumentsObject == "" {
		return
	}
	strict := strictFunction(node)
	simple := !strict
	for _, p := range node.Parameters() {
		d := p.AsParameterDeclaration()
		if p.Name().Kind != ast.KindIdentifier || d.Initializer != nil || d.DotDotDotToken != nil {
			simple = false
		}
	}
	if simple {
		seen := map[string]bool{}
		parameters := node.Parameters()
		for i := len(parameters) - 1; i >= 0; i-- {
			p := parameters[i]
			name := p.Name().Text()
			if seen[name] {
				continue
			}
			seen[name] = true
			cell := b.e.bindings[p]
			if cell == nil {
				continue
			}
			key := strconv.Quote(fmt.Sprint(i))
			b.emit(fmt.Sprintf("if len(args)>%d {%s.argumentMap[%s]=tsProperty{readField:func()tsValue{return %s.get()},writeField:func(value tsValue)tsValue{return %s.set(value)}}}", i, b.argumentsObject, key, cell.name, cell.name))
		}
	}
}

const argumentsRuntime = `
func tsArgumentsObject(args []tsValue,strict bool)*tsObject{
 object:=tsNewObject();object.argumentMap=map[string]tsProperty{};object.descriptors=map[string]*tsDescriptor{}
 for i,value:=range args{object.set(strconv.Itoa(i),value)};object.set("length",tsNumberValue(float64(len(args))));object.descriptors["length"]=&tsDescriptor{writable:true,configurable:true}
 if strict{thrower:=tsNativeMethod(func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{tsPropertyFailure("Restricted arguments.callee access");return tsU});object.order=append(object.order,"callee");object.descriptors["callee"]=&tsDescriptor{accessor:true,getter:thrower,setter:thrower}}else{object.set("callee",tsU);object.descriptors["callee"]=&tsDescriptor{writable:true,configurable:true}}
 key:=tsPropertyKey(tsWellKnownSymbol("iterator"));object.set(key,tsNamedNativeMethod("values",0,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{index:=0;return tsIteratorValueObject(&tsIterator{pull:func()(tsValue,bool){length:=tsNumber(tsGet(loop,receiver,"length"));if float64(index)>=length{return tsU,false};value:=tsGet(loop,receiver,tsNumberValue(float64(index)));index++;return value,true}})}));object.descriptors[key]=&tsDescriptor{writable:true,configurable:true}
 return object
}
`
