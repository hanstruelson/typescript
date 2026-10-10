package goemit

import (
	"fmt"
	"github.com/microsoft/TypeScript/tsc/internal/ast"
)

// Async iteration uses the existing protected-completion machinery. Closing an
// iterator is awaited on break, return, and throw, including a pending finally.
func (b *machineBuilder) forAwaitStatement(node *ast.Node, label string) {
	n := node.AsForInOrOfStatement()
	if n.Initializer.Kind != ast.KindVariableDeclarationList || len(n.Initializer.AsVariableDeclarationList().Declarations.Nodes) != 1 {
		b.e.fail(node, "for-await-of requires one declaration")
		return
	}
	decl := n.Initializer.AsVariableDeclarationList().Declarations.Nodes[0]
	cell := b.e.bindings[decl]
	pattern := ast.IsBindingPattern(decl.Name())
	if cell == nil && !pattern {
		return
	}
	lexical := n.Initializer.AsVariableDeclarationList().Flags&ast.NodeFlagsBlockScoped != 0
	if pattern {
		b.declarePatternCells(decl.Name())
	} else if cell.lexical {
		b.allocate(cell)
	}
	iterable := b.expression(n.Expression)
	iterator := b.typedTemp("tsAsyncIterate("+iterable+")", "*tsAsyncIterator")
	test, received, body, finished, cleanup, closed, end := b.block(), b.block(), b.block(), b.block(), b.block(), b.block(), b.block()
	depth := b.depth
	b.emit(fmt.Sprintf("m.handlers=append(m.handlers,tsHandler{catch:-1,finally:%d,end:%d});m.pc=%d;return", cleanup, end, test))
	b.depth++
	b.loops = append(b.loops, loopTarget{end, test, depth, b.depth, label})
	b.current = test
	b.emit(fmt.Sprintf("m.await(%s.next(loop),%d);return", iterator, received))
	b.current = received
	b.emit(fmt.Sprintf("if tsTruthy(tsFSProperty(m.result,\"done\")){m.pc=%d}else{m.pc=%d};return", finished, body))
	b.current = body
	value := b.temp("tsFSProperty(m.result,\"value\")")
	if pattern {
		if lexical {
			b.declarePatternCells(decl.Name())
		}
		b.bindPattern(decl.Name(), value, lexical)
	} else if lexical {
		b.allocate(cell)
		b.emit(cell.name + ".init(" + value + ")")
	} else {
		b.emit(cell.name + ".set(" + value + ")")
	}
	b.statement(n.Statement)
	b.jump(test)
	b.loops = b.loops[:len(b.loops)-1]
	b.current = finished
	b.emit("m.handlers=m.handlers[:len(m.handlers)-1]")
	b.jump(end)
	b.current = cleanup
	b.emit(fmt.Sprintf("m.await(%s.close(loop),%d);return", iterator, closed))
	b.current = closed
	b.emit("m.endFinally();return")
	b.depth--
	b.current = end
}

const fsIteratorRuntime = `
type tsAsyncIterator struct {
    resume func(*tsLoop,string,tsValue)*tsPromise
    nextInto func(*tsLoop,tsValue)*tsPromise
    next func(*tsLoop)*tsPromise
    close func(*tsLoop)*tsPromise
}

func tsIteratorResult(value tsValue, done bool) tsValue {
    out:=tsNewObject();out.set("value",value);out.set("done",tsBooleanValue(done));return tsObjectValue(out)
}

func tsAsyncMissingThrow(loop *tsLoop,protocol tsValue)*tsPromise{
 promise:=loop.promise();failure:=tsErrorValue(&tsRuntimeError{name:"TypeError",message:"Delegated iterator has no throw method"})
 func(){defer func(){if thrown:=recover();thrown!=nil{promise.settle(tsResult{tsUnwrap(thrown),true})}}();method:=tsGet(loop,protocol,"return");if tsNullish(method){promise.settle(tsResult{failure,true});return};closed:=tsCallReceiver(loop,method,protocol);loop.await(closed,func(result tsResult){if result.rejected{promise.settle(result)}else if tsIsPrimitiveValue(result.value){promise.settle(tsResult{tsErrorValue(&tsRuntimeError{name:"TypeError",message:"Iterator close result must be an object"}),true})}else{promise.settle(tsResult{failure,true})}})}()
 return promise
}

func tsAsyncIterate(value tsValue) *tsAsyncIterator {
 if value.kind==tsIteratorKind&&(*tsIterator)(value.ref).machine!=nil&&(*tsIterator)(value.ref).machine.asyncGenerator!=nil{state:=(*tsIterator)(value.ref).machine.asyncGenerator;return &tsAsyncIterator{resume:state.request,next:func(loop *tsLoop)*tsPromise{return state.request(loop,"next",tsU)},close:func(loop *tsLoop)*tsPromise{return state.request(loop,"return",tsU)}}}
    if value.kind==tsObjectKind{object:=(*tsObject)(value.ref);if object.asyncIterator!=nil{return object.asyncIterator()}}
    if value.kind==tsObjectKind||value.kind==tsInstanceKind{method:=tsGet(loop,value,tsWellKnownSymbol("asyncIterator"));if !tsNullish(method){protocol:=tsCallReceiver(loop,method,value);if tsIsPrimitiveValue(protocol){tsPropertyFailure("Async iterator must be an object")};nextMethod:=tsGet(loop,protocol,"next");resume:=func(loop *tsLoop,action string,value tsValue)*tsPromise{method:=nextMethod;if action!="next"{method=tsGet(loop,protocol,action)};if tsNullish(method){if action=="return"{return loop.resolved(tsIteratorResult(value,true),false)};if action=="throw"{return tsAsyncMissingThrow(loop,protocol)};tsPropertyFailure("Async iterator method is missing")};return tsAsyncIteratorPromise(loop,tsCallReceiver(loop,method,protocol,value))};return &tsAsyncIterator{resume:resume,next:func(loop *tsLoop)*tsPromise{return resume(loop,"next",tsU)},close:func(loop *tsLoop)*tsPromise{return resume(loop,"return",tsU)}}}}
    it:=tsIterate(value)
    resume:=func(loop *tsLoop,action string,value tsValue)*tsPromise{
      promise:=loop.promise()
      func(){defer func(){if failure:=recover();failure!=nil{promise.settle(tsResult{tsUnwrap(failure),true})}}()
        result:=tsU
        if it.machine!=nil{result=tsGeneratorStep(loop,it,action,value)}else if !tsIsUndefined(it.protocol){method:=it.nextMethod;if action!="next"{method=tsGet(loop,it.protocol,action)};if tsNullish(method){if action=="return"{result=tsIteratorResult(value,true)}else if action=="throw"{tsIteratorClose(loop,it,false);tsPropertyFailure("Iterator has no throw method")}else{tsPropertyFailure("Iterator has no next method")}}else{result=tsCallReceiver(loop,method,it.protocol,value)}}else if action=="return"{result=tsIteratorResult(value,true)}else if action=="throw"{tsPropertyFailure("Iterator has no throw method")}else{ok:=it.next(loop);item:=tsU;if ok{item=it.value};result=tsIteratorResult(item,!ok)}
        if tsIsPrimitiveValue(result){tsPropertyFailure("Iterator result must be an object")};done:=tsTruthy(tsGet(loop,result,"done"));item:=tsGet(loop,result,"value")
        loop.await(item,func(settled tsResult){if settled.rejected{promise.settle(settled)}else{promise.settle(tsResult{tsIteratorResult(settled.value,done),false})}})
      }()
      return promise
    }
    return &tsAsyncIterator{resume:resume,next:func(loop *tsLoop)*tsPromise{return resume(loop,"next",tsU)},close:func(loop *tsLoop)*tsPromise{return resume(loop,"return",tsU)}}
}

func tsAsyncIteratorObject(iterator *tsAsyncIterator) tsValue {
    out:=tsNewObject();out.asyncIterator=func()*tsAsyncIterator{return iterator}
    out.set("next",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsPromiseValue(iterator.next(loop))})))
    out.set("return",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsPromiseValue(iterator.close(loop))})))
    return tsObjectValue(out)
}
`
