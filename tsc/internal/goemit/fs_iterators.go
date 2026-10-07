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
	b.abrupt("jump", "tsU", end, depth)
	b.current = cleanup
	b.emit(fmt.Sprintf("m.await(%s.close(loop),%d);return", iterator, closed))
	b.current = closed
	b.emit("m.endFinally();return")
	b.depth--
	b.current = end
}

const fsIteratorRuntime = `
type tsAsyncIterator struct {
    nextInto func(*tsLoop,tsValue)*tsPromise
    next func(*tsLoop)*tsPromise
    close func(*tsLoop)*tsPromise
}

func tsIteratorResult(value tsValue, done bool) tsValue {
    out:=tsNewObject();out.set("value",value);out.set("done",tsBooleanValue(done));return tsObjectValue(out)
}

func tsAsyncIterate(value tsValue) *tsAsyncIterator {
    if value.kind==tsObjectKind{object:=(*tsObject)(value.ref);if object.asyncIterator!=nil{return object.asyncIterator()}}
    it:=tsIterate(value)
    return &tsAsyncIterator{
        next:func(loop *tsLoop)*tsPromise{ok:=it.next(loop);item:=tsU;if ok{item=it.value};return loop.resolved(tsIteratorResult(item,!ok),false)},
        close:func(loop *tsLoop)*tsPromise{return loop.resolved(tsIteratorResult(tsU,true),false)},
    }
}

func tsAsyncIteratorObject(iterator *tsAsyncIterator) tsValue {
    out:=tsNewObject();out.asyncIterator=func()*tsAsyncIterator{return iterator}
    out.set("next",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsPromiseValue(iterator.next(loop))})))
    out.set("return",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsPromiseValue(iterator.close(loop))})))
    return tsObjectValue(out)
}
`
