package goemit

// internalStreamDuplexifyRuntime normalizes the Node stream, Web stream,
// iterable, and promise inputs accepted by Bun's internal duplexify helper.
const internalStreamDuplexifyRuntime = `
type tsNodeDuplexifyGenState struct { loop *tsLoop; waiter *tsPromise; queued []tsValue; controller,signal tsValue; stopped bool }
func tsNodeDuplexifyKey(name string) tsValue { return tsStringReference(tsStringUTF8(name)) }
func tsNodeDuplexifyGet(loop *tsLoop,value tsValue,name string) tsValue { return tsGet(loop,value,tsNodeDuplexifyKey(name)) }
func tsNodeDuplexifySet(loop *tsLoop,value tsValue,name string,item tsValue) { tsSet(loop,value,tsNodeDuplexifyKey(name),item) }
func tsNodeDuplexifyCall(loop *tsLoop,value tsValue,name string,args ...tsValue) tsValue { method:=tsNodeDuplexifyGet(loop,value,name);if method.kind!=tsFunctionKind{return tsU};return tsCallReceiver(loop,method,value,args...) }
func tsNodeDuplexifyAwait(loop *tsLoop,value tsValue,callback func(tsResult)) { if value.kind==tsPromiseKind{loop.await(value,callback)}else{callback(tsResult{value,false})} }
func tsNodeDuplexifyIsIterable(loop *tsLoop,value tsValue) bool { return !tsNullish(tsGet(loop,value,tsWellKnownSymbol("iterator")))||!tsNullish(tsGet(loop,value,tsWellKnownSymbol("asyncIterator"))) }
func tsNodeDuplexifyIsReadable(loop *tsLoop,value tsValue) bool { return !tsNullish(tsNodeDuplexifyGet(loop,value,"_readableState"))&&tsTruthy(tsNodeDuplexifyGet(loop,value,"readable")) }
func tsNodeDuplexifyIsWritable(loop *tsLoop,value tsValue) bool { return !tsNullish(tsNodeDuplexifyGet(loop,value,"_writableState"))&&tsTruthy(tsNodeDuplexifyGet(loop,value,"writable")) }
func tsNodeDuplexifyIsNode(loop *tsLoop,value tsValue) bool { return tsTruthy(value)&&(tsNodeDuplexifyGet(loop,value,"on").kind==tsFunctionKind&&(tsNodeDuplexifyGet(loop,value,"pipe").kind==tsFunctionKind||tsNodeDuplexifyGet(loop,value,"write").kind==tsFunctionKind||!tsNullish(tsNodeDuplexifyGet(loop,value,"_readableState"))||!tsNullish(tsNodeDuplexifyGet(loop,value,"_writableState")))) }
func tsNodeDuplexifyOptions(values map[string]tsValue) tsValue { object:=tsNewObject();for key,value:=range values{object.set(key,value)};return tsObjectValue(object) }
func tsNodeDuplexifyAsyncGen(loop *tsLoop,fn tsValue) (tsValue,tsValue,tsValue,tsValue) {
	controller:=tsAbortController(loop);state:=&tsNodeDuplexifyGenState{loop:loop,controller:controller,signal:tsFSProperty(controller,"signal")}
	iterator:=&tsAsyncIterator{}
	iterator.next=func(loop *tsLoop)*tsPromise{if len(state.queued)>0{item:=state.queued[0];state.queued=state.queued[1:];if callback:=tsNodeDuplexifyGet(loop,item,"callback");callback.kind==tsFunctionKind{tsInternalStreamDestroyNextTick(loop,callback)};return loop.resolved(tsIteratorResult(tsNodeDuplexifyGet(loop,item,"chunk"),false),false)};if state.stopped{return loop.resolved(tsIteratorResult(tsU,true),false)};state.waiter=loop.promise();return state.waiter}
	iterator.close=func(loop *tsLoop)*tsPromise{state.stopped=true;return loop.resolved(tsIteratorResult(tsU,true),false)}
	input:=tsAsyncIteratorObject(iterator);options:=tsNewObject();options.set("signal",state.signal);result:=tsCall(loop,fn,input,tsObjectValue(options))
	write:=tsNamedNativeMethod("write",3,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{chunk,callback:=tsArg(args,0),tsArg(args,2);item:=tsNewObject();item.set("chunk",chunk);item.set("done",tsBooleanValue(false));item.set("callback",callback);value:=tsObjectValue(item);if state.waiter!=nil{if callback.kind==tsFunctionKind{tsInternalStreamDestroyNextTick(loop,callback)};state.waiter.settle(tsResult{tsIteratorResult(chunk,false),false});state.waiter=nil}else{state.queued=append(state.queued,value)};return tsU})
	final:=tsNamedNativeMethod("final",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{callback:=tsArg(args,0);state.stopped=true;if state.waiter!=nil{state.waiter.settle(tsResult{tsIteratorResult(tsU,true),false});state.waiter=nil};if result.kind==tsPromiseKind{tsNodeDuplexifyAwait(loop,result,func(settled tsResult){if settled.rejected{tsCall(loop,callback,settled.value)}else{tsCall(loop,callback,tsU)}})}else{tsCall(loop,callback,tsU)};return tsU})
	destroy:=tsNamedNativeMethod("destroy",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{err,callback:=tsArg(args,0),tsArg(args,1);tsNodeDuplexifyCall(loop,state.controller,"abort",err);state.stopped=true;if state.waiter!=nil{state.waiter.settle(tsResult{tsIteratorResult(tsU,true),false});state.waiter=nil};if callback.kind==tsFunctionKind{tsCall(loop,callback,err)};return tsU})
	return result,write,final,destroy
}
func tsNodeDuplexifyPair(loop *tsLoop,pair tsValue) tsValue {
	r,w:=tsNodeDuplexifyGet(loop,pair,"readable"),tsNodeDuplexifyGet(loop,pair,"writable")
	if !tsNullish(r)&&tsNodeDuplexifyGet(loop,r,"read").kind!=tsFunctionKind { module:=tsInternalStreamReadableModule(loop);constructor:=tsGet(loop,module,tsNodeDuplexifyKey("Readable"));r=tsCall(loop,tsGet(loop,constructor,tsNodeDuplexifyKey("wrap")),r) }
	readable,writable:=tsNodeDuplexifyIsReadable(loop,r),tsNodeDuplexifyIsWritable(loop,w)
	options:=tsNodeDuplexifyOptions(map[string]tsValue{"readable":tsBooleanValue(readable),"writable":tsBooleanValue(writable),"readableObjectMode":tsBooleanValue(tsTruthy(tsNodeDuplexifyGet(loop,r,"readableObjectMode"))),"writableObjectMode":tsBooleanValue(tsTruthy(tsNodeDuplexifyGet(loop,w,"writableObjectMode")))})
	var d tsValue;var ondrain,onfinish,onreadable,onclose tsValue
	finished:=func(err tsValue){callback:=onclose;onclose=tsNull;if callback.kind==tsFunctionKind{tsCall(loop,callback,err)}else if tsTruthy(err){tsNodeDuplexifyCall(loop,d,"destroy",err)}}
	if writable {
		write:=tsNamedNativeMethod("_write",3,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{chunk,encoding,callback:=tsArg(args,0),tsArg(args,1),tsArg(args,2);result:=tsCallReceiver(loop,tsNodeDuplexifyGet(loop,w,"write"),w,chunk,encoding);if tsTruthy(result){tsCall(loop,callback,tsU)}else{ondrain=callback};return tsU})
		final:=tsNamedNativeMethod("_final",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{callback:=tsArg(args,0);onfinish=callback;tsNodeDuplexifyCall(loop,w,"end");return tsU})
		(*tsObject)(options.ref).set("write",write);(*tsObject)(options.ref).set("final",final)
	}
	if readable {
		read:=tsNamedNativeMethod("_read",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{for {value:=tsNodeDuplexifyCall(loop,r,"read");if tsNullish(value){onreadable=tsNamedNativeMethod("resume-readable",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{if onreadable.kind==tsFunctionKind{onreadable=tsNull};tsNodeDuplexifyCall(loop,d,"_read");return tsU});tsCallReceiver(loop,tsNodeDuplexifyGet(loop,r,"once"),r,tsNodeDuplexifyKey("readable"),onreadable);return tsU};accepted:=tsNodeDuplexifyCall(loop,d,"push",value);if !tsTruthy(accepted){return tsU}}})
		(*tsObject)(options.ref).set("read",read)
	}
	destroy:=tsNamedNativeMethod("_destroy",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{err,callback:=tsArg(args,0),tsArg(args,1);if tsNullish(err)&&!tsIsNull(onclose){err=tsFSException("AbortError","ABORT_ERR","The operation was aborted")};onreadable=tsNull;ondrain=tsNull;onfinish=tsNull;if tsIsNull(onclose){tsCall(loop,callback,err)}else{onclose=callback;tsInternalStreamDestroyer(loop,w,err);tsInternalStreamDestroyer(loop,r,err)};return tsU})
	(*tsObject)(options.ref).set("destroy",destroy);d=tsNodeDuplexCreate(loop,tsU,options)
	if writable {
		closed:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{writable=false;err:=tsArg(args,0);if tsTruthy(err){tsInternalStreamDestroyer(loop,r,err)};finished(err);return tsU}));eosOptions:=tsNewObject();eosOptions.set("writable",tsBooleanValue(true));eosOptions.set("readable",tsBooleanValue(false));tsCall(loop,tsInternalStreamEndOfStreamModule(loop),w,tsObjectValue(eosOptions),closed)
		ondrain=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{callback:=ondrain;ondrain=tsNull;if callback.kind==tsFunctionKind{tsCall(loop,callback,tsU)};return tsU}));tsCallReceiver(loop,tsNodeDuplexGet(loop,w,"on"),w,tsNodeDuplexKey("drain"),ondrain)
		onfinish=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{callback:=onfinish;onfinish=tsNull;if callback.kind==tsFunctionKind{tsCall(loop,callback,tsU)};return tsU}));tsCallReceiver(loop,tsNodeDuplexGet(loop,w,"on"),w,tsNodeDuplexKey("finish"),onfinish)
	}
	if readable {
		closed:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{readable=false;err:=tsArg(args,0);if tsTruthy(err){tsInternalStreamDestroyer(loop,w,err)};finished(err);return tsU}));eosOptions:=tsNewObject();eosOptions.set("readable",tsBooleanValue(true));eosOptions.set("writable",tsBooleanValue(false));tsCall(loop,tsInternalStreamEndOfStreamModule(loop),r,tsObjectValue(eosOptions),closed)
		onend:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{tsNodeDuplexifyCall(loop,d,"push",tsNull);return tsU}));tsCallReceiver(loop,tsNodeDuplexGet(loop,r,"on"),r,tsNodeDuplexKey("end"),onend)
	}
	return d
}
func tsNodeDuplexifyBody(loop *tsLoop,body,name tsValue) tsValue {
	if tsInternalStreamIsDuplex(loop,body){return body}
	if tsInternalStreamIsReadableNode(loop,body,false){pair:=tsNewObject();pair.set("readable",body);return tsNodeDuplexifyPair(loop,tsObjectValue(pair))}
	if tsInternalStreamIsWritableNode(loop,body){pair:=tsNewObject();pair.set("writable",body);return tsNodeDuplexifyPair(loop,tsObjectValue(pair))}
	if tsNodeDuplexifyIsNode(loop,body){return tsNodeDuplexCreate(loop,tsU,tsNodeDuplexifyOptions(map[string]tsValue{"readable":tsBooleanValue(false),"writable":tsBooleanValue(false)}))}
	if tsInternalStreamIsWebClass(body,"ReadableStream"){module:=tsInternalStreamReadableModule(loop);return tsCall(loop,tsGet(loop,tsGet(loop,module,tsNodeDuplexifyKey("Readable")),tsNodeDuplexifyKey("fromWeb")),body)}
	if tsInternalStreamIsWebClass(body,"WritableStream"){module:=tsInternalStreamWritableModule(loop);return tsCall(loop,tsGet(loop,tsGet(loop,module,tsNodeDuplexifyKey("Writable")),tsNodeDuplexifyKey("fromWeb")),body)}
	if body.kind==tsFunctionKind {
		value,write,final,destroy:=tsNodeDuplexifyAsyncGen(loop,body);if tsInternalStreamIsDuplex(loop,value){return value}
		if tsNodeDuplexifyIsIterable(loop,value){options:=tsNodeDuplexifyOptions(map[string]tsValue{"objectMode":tsBooleanValue(true),"write":write,"final":final,"destroy":destroy});return tsInternalStreamFrom(loop,[]tsValue{tsFunctionValue(tsNodeDuplexSingleton.constructor),value,options})}
		if value.kind==tsPromiseKind {duplex:=tsNodeDuplexCreate(loop,tsU,tsNodeDuplexifyOptions(map[string]tsValue{"objectMode":tsBooleanValue(true),"readable":tsBooleanValue(false),"write":write,"destroy":destroy,"final":final}));promise:=(*tsPromise)(value.ref);promise.then(func(result tsResult){if result.rejected{tsInternalStreamDestroyer(loop,duplex,result.value)}else if !tsNullish(result.value){tsNodeDuplexifyCall(loop,duplex,"push",result.value)}});return duplex}
		panic(tsThrown{tsFSException("TypeError","ERR_INVALID_RETURN_VALUE","Expected an Iterable, AsyncIterable or AsyncFunction return value")})
	}
	if body.kind==tsObjectKind&&(*tsObject)(body.ref).nativeClass!=nil&&(*tsObject)(body.ref).nativeClass.builtin=="Blob" {arrayBuffer:=tsNodeDuplexifyCall(loop,body,"arrayBuffer");return tsNodeDuplexifyBody(loop,arrayBuffer,name)}
	if tsNodeDuplexifyIsIterable(loop,body){options:=tsNodeDuplexifyOptions(map[string]tsValue{"objectMode":tsBooleanValue(true),"writable":tsBooleanValue(false)});return tsInternalStreamFrom(loop,[]tsValue{tsFunctionValue(tsNodeDuplexSingleton.constructor),body,options})}
	readable,writable:=tsNodeDuplexifyGet(loop,body,"readable"),tsNodeDuplexifyGet(loop,body,"writable")
	if tsInternalStreamIsWebClass(readable,"ReadableStream")&&tsInternalStreamIsWebClass(writable,"WritableStream"){module:=tsInternalStreamDuplexModule(loop);return tsCall(loop,tsGet(loop,tsGet(loop,module,tsNodeDuplexifyKey("Duplex")),tsNodeDuplexifyKey("fromWeb")),body)}
	if tsGet(loop,body,tsNodeDuplexifyKey("writable")).kind==tsObjectKind||tsGet(loop,body,tsNodeDuplexifyKey("readable")).kind==tsObjectKind {
		pair:=tsNewObject();if tsTruthy(readable){if !tsInternalStreamIsReadableNode(loop,readable,false){readable=tsNodeDuplexifyBody(loop,readable,tsNodeDuplexifyKey("readable"))};pair.set("readable",readable)};if tsTruthy(writable){if !tsInternalStreamIsWritableNode(loop,writable){writable=tsNodeDuplexifyBody(loop,writable,tsNodeDuplexifyKey("writable"))};pair.set("writable",writable)};return tsNodeDuplexifyPair(loop,tsObjectValue(pair))
	}
	if body.kind==tsPromiseKind {duplex:=tsNodeDuplexCreate(loop,tsU,tsNodeDuplexifyOptions(map[string]tsValue{"objectMode":tsBooleanValue(true),"writable":tsBooleanValue(false)}));(*tsPromise)(body.ref).then(func(result tsResult){if result.rejected{tsInternalStreamDestroyer(loop,duplex,result.value)}else{if !tsNullish(result.value){tsNodeDuplexifyCall(loop,duplex,"push",result.value)};tsNodeDuplexifyCall(loop,duplex,"push",tsNull)}});return duplex}
	panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_TYPE","The "+tsText(name)+" argument must be a Blob, ReadableStream, WritableStream, Stream, Iterable, AsyncIterable, Function, readable/writable pair, or Promise")})
}
func tsInternalStreamDuplexifyModule(loop *tsLoop) tsValue {module:=tsNewObject();module.set("default",tsNamedNativeMethod("duplexify",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{name:=tsArg(args,1);if tsIsUndefined(name){name=tsStringReference(tsStringUTF8("body"))};return tsNodeDuplexifyBody(loop,tsArg(args,0),name)}));return tsObjectValue(module)}
`
