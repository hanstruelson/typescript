package goemit

// Filesystem Web streams serialize pulls and deliver all reactions on the
// consumer's event loop. Their source is a native async iterator of byte views.
const webStreamRuntime = `
type tsWebReadRequest struct { promise *tsPromise; loop *tsLoop; token uint64; canceled bool;view tsValue;minimum,capacity,filled int }
type tsWebStream struct {
    iterator *tsAsyncIterator;object *tsObject;locked,ended,busy bool;failure tsValue
    token uint64;queue []*tsWebReadRequest;buffered []tsValue;readerClosed *tsPromise
}
func tsWebStreamValue(iterator *tsAsyncIterator)tsValue{
    stream:=&tsWebStream{iterator:iterator,object:tsNewObject()};stream.object.webStream=stream;stream.object.nativeClass=tsNativeClass("ReadableStream");stream.object.set("locked",tsBooleanValue(false))
    stream.object.set("getReader",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return stream.acquireReader(loop,tsArg(args,0))})))
    stream.object.set("cancel",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{if stream.locked{return tsPromiseValue(loop.resolved(tsFSException("TypeError","ERR_INVALID_STATE","ReadableStream is locked"),true))};return tsPromiseValue(stream.cancelStream(loop))})))
    factory:=func(preventCancel bool)*tsAsyncIterator{
        if stream.locked{panic(tsThrown{tsFSException("TypeError","ERR_INVALID_STATE","ReadableStream is locked")})};stream.locked=true;stream.token++;token:=stream.token;stream.object.set("locked",tsBooleanValue(true))
        return &tsAsyncIterator{
            next:func(loop *tsLoop)*tsPromise{promise:=stream.readNext(loop,token);promise.then(func(result tsResult){if result.rejected||tsTruthy(tsFSProperty(result.value,"done")){stream.releaseReader(token)}});return promise},
            close:func(loop *tsLoop)*tsPromise{stream.releaseReader(token);if preventCancel{return loop.resolved(tsIteratorResult(tsU,true),false)};promise:=loop.promise();loop.await(tsPromiseValue(stream.cancelStream(loop)),func(result tsResult){if !result.rejected{result.value=tsIteratorResult(tsU,true)};promise.settle(result)});return promise},
        }
    }
    stream.object.asyncIterator=func()*tsAsyncIterator{return factory(false)}
    stream.object.set("values",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsAsyncIteratorObject(factory(tsTruthy(tsFSProperty(tsArg(args,0),"preventCancel"))))})))
    return tsObjectValue(stream.object)
}
func(stream *tsWebStream)acquireReader(loop *tsLoop,options tsValue)tsValue{
    if stream.locked{panic(tsThrown{tsFSException("TypeError","ERR_INVALID_STATE","ReadableStream is locked")})};mode:=tsFSProperty(options,"mode");byob:=false;if !tsIsUndefined(mode){if mode.kind!=tsStringKind||tsText(mode)!="byob"{tsFSInvalid("mode","undefined or byob")};byob=true}
    stream.locked=true;stream.token++;token:=stream.token;stream.object.set("locked",tsBooleanValue(true));out:=tsNewObject();class:="ReadableStreamDefaultReader";if byob{class="ReadableStreamBYOBReader"};out.nativeClass=tsNativeClass(class)
    closed:=loop.promise();stream.readerClosed=closed;if !tsIsUndefined(stream.failure){closed.settle(tsResult{stream.failure,true})}else if stream.ended{closed.settle(tsResult{tsU,false})};out.set("closed",tsPromiseValue(closed))
    out.set("read",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{if !byob{return tsPromiseValue(stream.readNext(loop,token))};var result tsValue;func(){defer func(){if failure:=recover();failure!=nil{result=tsPromiseValue(loop.resolved(tsUnwrap(failure),true))}}();view:=tsArg(args,0);if view.kind!=tsTypedArrayKind{tsFSInvalid("view","a TypedArray")};array:=(*tsTypedArray)(view.ref);minimum:=int(tsFSInteger(tsFSProperty(tsArg(args,1),"min"),"min",1,1));if minimum>array.length{panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","min exceeds the BYOB view length")})};transferred:=tsBYOBTransfer(view);result=tsPromiseValue(stream.readBYOB(loop,token,transferred,minimum))}();return result})))
    out.set("cancel",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{if !stream.locked||stream.token!=token{return tsPromiseValue(loop.resolved(tsFSException("TypeError","ERR_INVALID_STATE","Reader has been released"),true))};return tsPromiseValue(stream.cancelStream(loop))})))
    out.set("releaseLock",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{if stream.token==token&&stream.locked{stream.releaseReader(token);failure:=tsFSException("TypeError","ERR_INVALID_STATE","Reader has been released");closed.mu.Lock();closed.handled=true;closed.mu.Unlock();closed.settle(tsResult{failure,true});for _,request:=range stream.queue{if request.token==token{request.canceled=true;request.promise.settle(tsResult{failure,true})}}};return tsU})))
    return tsObjectValue(out)
}
func(stream *tsWebStream)releaseReader(token uint64){if stream.token==token{stream.locked=false;stream.object.set("locked",tsBooleanValue(false));stream.readerClosed=nil}}
func(stream *tsWebStream)readNext(loop *tsLoop,token uint64)*tsPromise{
    if !stream.locked||stream.token!=token{return loop.resolved(tsFSException("TypeError","ERR_INVALID_STATE","Reader has been released"),true)}
    if !tsIsUndefined(stream.failure){return loop.resolved(stream.failure,true)};if stream.ended&&len(stream.buffered)==0{return loop.resolved(tsIteratorResult(tsU,true),false)}
    promise:=loop.promise();stream.queue=append(stream.queue,&tsWebReadRequest{promise:promise,loop:loop,token:token});stream.pumpStream();return promise
}
func(stream *tsWebStream)readBYOB(loop *tsLoop,token uint64,view tsValue,minimum int)*tsPromise{
    if !stream.locked||stream.token!=token{return loop.resolved(tsFSException("TypeError","ERR_INVALID_STATE","Reader has been released"),true)}
    if !tsIsUndefined(stream.failure){return loop.resolved(stream.failure,true)}
    array:=(*tsTypedArray)(view.ref);promise:=loop.promise();stream.queue=append(stream.queue,&tsWebReadRequest{promise:promise,loop:loop,token:token,view:view,minimum:minimum*array.bytes,capacity:array.length*array.bytes});stream.pumpStream();return promise
}
func(stream *tsWebStream)commitBYOB(request *tsWebReadRequest,done bool)bool{
    array:=(*tsTypedArray)(request.view.ref)
    if request.filled==0&&done{array.length=0;request.promise.settle(tsResult{tsIteratorResult(request.view,true),false});return true}
    aligned:=request.filled-request.filled%array.bytes
    if aligned<request.minimum&&!done{return false}
    if done&&request.filled%array.bytes!=0{stream.failure=tsFSException("TypeError","ERR_INVALID_STATE","The stream ended with a partial BYOB element");request.promise.settle(tsResult{stream.failure,true});return true}
    if aligned<request.filled{data:=tsFSBytes(request.view);tail:=tsByteView(append([]byte{},data[aligned:request.filled]...));stream.buffered=append([]tsValue{tail},stream.buffered...)}
    array.length=aligned/array.bytes;request.promise.settle(tsResult{tsIteratorResult(request.view,false),false});return true
}
func(stream *tsWebStream)pumpStream(){
    if stream.busy{return};for len(stream.queue)>0{
        request:=stream.queue[0];if request.canceled{stream.queue=stream.queue[1:];continue};loop:=request.loop;byob:=request.view.kind==tsTypedArrayKind
        if !tsIsUndefined(stream.failure){stream.queue=stream.queue[1:];request.promise.settle(tsResult{stream.failure,true});continue}
        if len(stream.buffered)>0{
            value:=stream.buffered[0];stream.buffered=stream.buffered[1:]
            if !byob{stream.queue=stream.queue[1:];request.promise.settle(tsResult{tsIteratorResult(value,false),false});continue}
            source:=tsFSBytes(value);target:=tsFSBytes(request.view);n:=copy(target[request.filled:],source);request.filled+=n;if n<len(source){stream.buffered=append([]tsValue{tsByteView(source[n:])},stream.buffered...)}
            if stream.commitBYOB(request,false){stream.queue=stream.queue[1:]};continue
        }
        if stream.ended{stream.queue=stream.queue[1:];if byob{stream.commitBYOB(request,true)}else{request.promise.settle(tsResult{tsIteratorResult(tsU,true),false})};continue}
        stream.busy=true;direct:=byob&&stream.iterator.nextInto!=nil;var promise *tsPromise
        if direct{array:=(*tsTypedArray)(request.view.ref);view:=tsNativeBufferView[uint8](array.buffer,array.offset+request.filled,request.capacity-request.filled,"Uint8Array",func(number float64)uint8{return uint8(tsNumberUint32(number))});promise=stream.iterator.nextInto(loop,tsTypedArrayValue(view))}else{promise=stream.iterator.next(loop)}
        loop.await(tsPromiseValue(promise),func(result tsResult){stream.busy=false
            if result.rejected{stream.failure=result.value;stream.ended=true;if stream.readerClosed!=nil{stream.readerClosed.settle(result)}}else if tsTruthy(tsFSProperty(result.value,"done")){stream.ended=true;if stream.readerClosed!=nil{stream.readerClosed.settle(tsResult{tsU,false})}}
            if request.canceled{
                stream.queue=stream.queue[1:];if !stream.ended&&!result.rejected{stream.buffered=append(stream.buffered,tsFSProperty(result.value,"value"))}
            }else if result.rejected{stream.queue=stream.queue[1:];request.promise.settle(result)
            }else if !byob{stream.queue=stream.queue[1:];request.promise.settle(result)
            }else{
                if !tsTruthy(tsFSProperty(result.value,"done")){source:=tsFSBytes(tsFSProperty(result.value,"value"));if direct{request.filled+=len(source)}else{target:=tsFSBytes(request.view);n:=copy(target[request.filled:],source);request.filled+=n;if n<len(source){stream.buffered=append(stream.buffered,tsByteView(source[n:]))}}}
                if stream.commitBYOB(request,stream.ended){stream.queue=stream.queue[1:]}
            }
            stream.pumpStream()
        });return
    }
}

func(stream *tsWebStream)cancelStream(loop *tsLoop)*tsPromise{
    stream.ended=true;stream.buffered=nil;for _,request:=range stream.queue{request.promise.settle(tsResult{tsIteratorResult(tsU,true),false});request.canceled=true};if stream.readerClosed!=nil{stream.readerClosed.settle(tsResult{tsU,false})};return stream.iterator.close(loop)
}
`
