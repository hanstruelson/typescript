package goemit

// internalWebstreamsAdaptersRuntime connects the classic Node stream
// constructors to the runtime's Web ReadableStream and WritableStream values.
const internalWebstreamsAdaptersRuntime = `
type tsNodeWebWritableState struct {
	object, source, controller tsValue
	locked, closed bool
	failure tsValue
	closedPromise *tsPromise
}
type tsNodeReadableFromWebState struct { stream, readable, reader tsValue; closed bool }
var tsNodeWebAdapterSingleton struct { sync.Once; writable map[uintptr]*tsNodeWebWritableState; readable map[uintptr]*tsNodeReadableFromWebState; readableFromWeb *tsFunction; errorSentinel tsValue; dep0201 bool }
func tsNodeWebAdapterKey(name string) tsValue { return tsStringReference(tsStringUTF8(name)) }
func tsNodeWebAdapterAwait(loop *tsLoop,value tsValue,callback func(tsResult)) {
	if value.kind==tsPromiseKind { loop.await(value,callback) } else { callback(tsResult{value,false}) }
}
func tsNodeWebAdapterThrow(loop *tsLoop,value tsValue) { if value.kind==tsErrorKind||value.kind==tsObjectKind&&(*tsObject)(value.ref).nativeError { panic(tsThrown{value}) };panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_TYPE","The value has an invalid type")}) }
func tsNodeWebAdapterWritableState(value tsValue) *tsNodeWebWritableState { if value.kind!=tsObjectKind{return nil};return tsNodeWebAdapterSingleton.writable[uintptr(value.ref)] }
func tsNodeWebAdapterResolve(loop *tsLoop,value tsValue) tsValue { if value.kind==tsPromiseKind{return value};return tsPromiseValue(loop.resolved(value,false)) }
func tsNodeWebAdapterReject(loop *tsLoop,value tsValue) tsValue { return tsPromiseValue(loop.resolved(value,true)) }
func tsNodeWebAdapterSettleClosed(state *tsNodeWebWritableState,result tsResult) { if state==nil||state.closedPromise==nil{return};state.closedPromise.settle(result) }
func tsNodeWebAdapterSourceCall(loop *tsLoop,state *tsNodeWebWritableState,name string,args ...tsValue) tsValue {
	if state==nil||tsNullish(state.source){return tsU};method:=tsGet(loop,state.source,tsNodeWebAdapterKey(name));if method.kind!=tsFunctionKind{return tsU};return tsCallReceiver(loop,method,state.source,args...)
}
func tsNodeWebAdapterWritableCreate(loop *tsLoop,source,strategy tsValue) tsValue {
	tsNodeWebAdapterSingleton.Do(func(){tsNodeWebAdapterSingleton.writable=map[uintptr]*tsNodeWebWritableState{};tsNodeWebAdapterSingleton.readable=map[uintptr]*tsNodeReadableFromWebState{};tsNodeWebAdapterSingleton.errorSentinel=tsNewSymbol(tsStringUTF8("kErrorSentinelAttached"))})
	object:=tsNewObject();object.nativeClass=tsNativeClass("WritableStream");object.set("locked",tsBooleanValue(false));state:=&tsNodeWebWritableState{object:tsObjectValue(object),source:source,failure:tsU,closedPromise:loop.promise()};object.set("desiredSize",tsNumberValue(1));tsNodeWebAdapterSingleton.writable[uintptr(unsafe.Pointer(object))]=state
	controller:=tsNewObject();state.controller=tsObjectValue(controller);controller.set("desiredSize",tsNumberValue(1))
	controller.set("error",tsNamedNativeMethod("error",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{if !state.closed{state.closed=true;state.failure=tsArg(args,0);tsNodeWebAdapterSettleClosed(state,tsResult{state.failure,true})};return tsU}))
	controller.set("signal",tsU)
	callAlgorithm:=func(name string,args ...tsValue)tsValue{return tsNodeWebAdapterSourceCall(loop,state,name,append(args,state.controller)...)}
	object.set("getWriter",tsNamedNativeMethod("getWriter",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{
		if state.locked{return tsNodeWebAdapterReject(loop,tsFSException("TypeError","ERR_INVALID_STATE","WritableStream is locked"))};state.locked=true;(*tsObject)(state.object.ref).set("locked",tsBooleanValue(true))
		writer:=tsNewObject();writer.nativeClass=tsNativeClass("WritableStreamDefaultWriter");writer.set("closed",tsPromiseValue(state.closedPromise));writer.set("ready",tsPromiseValue(loop.resolved(tsU,false)))
		writer.set("write",tsNamedNativeMethod("write",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{if state.closed{if !tsNullish(state.failure){return tsNodeWebAdapterReject(loop,state.failure)};return tsNodeWebAdapterReject(loop,tsFSException("TypeError","ERR_INVALID_STATE","WritableStream is closed"))};result:=callAlgorithm("write",tsArg(args,0));return tsNodeWebAdapterResolve(loop,result)}))
		writer.set("close",tsNamedNativeMethod("close",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNodeWebAdapterWritableClose(loop,state,callAlgorithm)}))
		writer.set("abort",tsNamedNativeMethod("abort",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeWebAdapterWritableAbort(loop,state,tsArg(args,0),callAlgorithm)}))
		writer.set("releaseLock",tsNamedNativeMethod("releaseLock",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{state.locked=false;tsSet(loop,state.object,tsNodeWebAdapterKey("locked"),tsBooleanValue(false));return tsU}))
		object.set("locked",tsBooleanValue(true));return tsObjectValue(writer)
	}))
	object.set("close",tsNamedNativeMethod("close",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNodeWebAdapterWritableClose(loop,state,callAlgorithm)}))
	object.set("abort",tsNamedNativeMethod("abort",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeWebAdapterWritableAbort(loop,state,tsArg(args,0),callAlgorithm)}))
	if !tsNullish(source){start:=tsGet(loop,source,tsNodeWebAdapterKey("start"));if start.kind==tsFunctionKind{result:=tsCallReceiver(loop,start,source,state.controller);tsNodeWebAdapterAwait(loop,result,func(settled tsResult){if settled.rejected{state.closed=true;state.failure=settled.value;tsNodeWebAdapterSettleClosed(state,settled)}})}}
	_ = strategy
	return tsObjectValue(object)
}
func tsNodeWebAdapterWritableClose(loop *tsLoop,state *tsNodeWebWritableState,callAlgorithm func(string,...tsValue)tsValue) tsValue {
	if state==nil{return tsNodeWebAdapterReject(loop,tsFSException("TypeError","ERR_INVALID_STATE","WritableStream is unavailable"))};if state.closed{if !tsNullish(state.failure){return tsNodeWebAdapterReject(loop,state.failure)};return tsNodeWebAdapterReject(loop,tsFSException("TypeError","ERR_INVALID_STATE","WritableStream is closed"))}
	result:=callAlgorithm("close");promise:=loop.promise();tsNodeWebAdapterAwait(loop,result,func(settled tsResult){state.closed=true;if settled.rejected{state.failure=settled.value;tsNodeWebAdapterSettleClosed(state,settled);promise.settle(settled)}else{tsNodeWebAdapterSettleClosed(state,tsResult{tsU,false});promise.settle(tsResult{tsU,false})}});return tsPromiseValue(promise)
}
func tsNodeWebAdapterWritableAbort(loop *tsLoop,state *tsNodeWebWritableState,reason tsValue,callAlgorithm func(string,...tsValue)tsValue) tsValue {
	if state==nil{return tsNodeWebAdapterResolve(loop,tsU)};if state.closed{if !tsNullish(state.failure){return tsNodeWebAdapterReject(loop,state.failure)};return tsNodeWebAdapterResolve(loop,tsU)};result:=callAlgorithm("abort",reason);promise:=loop.promise();tsNodeWebAdapterAwait(loop,result,func(settled tsResult){state.closed=true;if settled.rejected{state.failure=settled.value;tsNodeWebAdapterSettleClosed(state,settled);promise.settle(settled)}else{state.failure=reason;tsNodeWebAdapterSettleClosed(state,tsResult{reason,true});promise.settle(tsResult{tsU,false})}});return tsPromiseValue(promise)
}
func tsNodeWebAdapterEmptyReadable(loop *tsLoop,options tsValue) tsValue { return tsWebStreamValue(&tsAsyncIterator{next:func(loop *tsLoop)*tsPromise{return loop.resolved(tsIteratorResult(tsU,true),false)},close:func(loop *tsLoop)*tsPromise{return loop.resolved(tsIteratorResult(tsU,true),false)}}) }
func tsNodeWebAdapterWritableStreamFromClassic(loop *tsLoop,stream tsValue) tsValue {
	write:=tsGet(loop,stream,tsNodeWebAdapterKey("write"));on:=tsGet(loop,stream,tsNodeWebAdapterKey("on"));if tsNullish(stream)||write.kind!=tsFunctionKind||on.kind!=tsFunctionKind{tsInternalValidatorInvalidType("streamWritable","stream.Writable",stream)}
	if tsTruthy(tsGet(loop,stream,tsNodeWebAdapterKey("destroyed")))||!tsInternalStreamIsWritableNode(loop,stream){out:=tsNodeWebAdapterWritableCreate(loop,tsNull,tsU);state:=tsNodeWebAdapterWritableState(out);state.closed=true;tsNodeWebAdapterSettleClosed(state,tsResult{tsU,false});return out}
	source:=tsNewObject()
	source.set("write",tsNamedNativeMethod("write",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{
		chunk:=tsArg(args,0);if chunk.kind==tsArrayBufferKind{buffer:=(*tsArrayBuffer)(chunk.ref);chunk=tsTypedArrayValue(&tsTypedArray{buffer:buffer,length:len(buffer.data),bytes:1,name:"Uint8Array"})}
		needDrain:=tsTruthy(tsGet(loop,stream,tsNodeWebAdapterKey("writableNeedDrain")));accepted:=tsCallReceiver(loop,write,stream,chunk);if !needDrain&&tsTruthy(accepted){return tsU}
		promise:=loop.promise();var onDrain tsValue;onDrain=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{remove:=tsGet(loop,stream,tsNodeWebAdapterKey("removeListener"));if remove.kind==tsFunctionKind{tsCallReceiver(loop,remove,stream,tsNodeWebAdapterKey("drain"),onDrain)};promise.settle(tsResult{tsU,false});return tsU}));tsCallReceiver(loop,on,stream,tsNodeWebAdapterKey("drain"),onDrain);if !tsTruthy(tsGet(loop,stream,tsNodeWebAdapterKey("writableNeedDrain"))){promise.settle(tsResult{tsU,false})};return tsPromiseValue(promise)
	}))
	source.set("abort",tsNamedNativeMethod("abort",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{tsInternalStreamDestroyer(loop,stream,tsArg(args,0));return tsU}))
	source.set("close",tsNamedNativeMethod("close",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{promise:=loop.promise();if tsTruthy(tsGet(loop,stream,tsNodeWebAdapterKey("writableEnded"))){promise.settle(tsResult{tsU,false});return tsPromiseValue(promise)};callback:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{err:=tsArg(args,0);if tsTruthy(err){promise.settle(tsResult{err,true})}else{promise.settle(tsResult{tsU,false})};return tsU}));options:=tsNewObject();options.set("writable",tsBooleanValue(true));options.set("readable",tsBooleanValue(false));tsCall(loop,tsInternalStreamEndOfStreamModule(loop),stream,tsObjectValue(options),callback);end:=tsGet(loop,stream,tsNodeWebAdapterKey("end"));if end.kind==tsFunctionKind{tsCallReceiver(loop,end,stream)};return tsPromiseValue(promise)}))
	return tsNodeWebAdapterWritableCreate(loop,tsObjectValue(source),tsU)
}
func tsNodeWebAdapterWritableFromWeb(loop *tsLoop,writableStream,options tsValue) tsValue {
	options=tsNodeWebAdapterNormalizeOptions(options);if !tsInternalStreamIsWebClass(writableStream,"WritableStream"){tsInternalValidatorInvalidType("writableStream","WritableStream",writableStream)};tsInternalValidatorObject(options,tsNodeWebAdapterKey("options"),tsU)
	objectMode:=tsNodeWebAdapterGet(loop,options,"objectMode");if tsIsUndefined(objectMode){objectMode=tsBooleanValue(false)};decodeStrings:=tsNodeWebAdapterGet(loop,options,"decodeStrings");if tsIsUndefined(decodeStrings){decodeStrings=tsBooleanValue(true)};if objectMode.kind!=tsBooleanKind{tsInternalValidatorInvalidType("options.objectMode","a boolean",objectMode)};if decodeStrings.kind!=tsBooleanKind{tsInternalValidatorInvalidType("options.decodeStrings","a boolean",decodeStrings)}
	getWriter:=tsNodeWebAdapterGet(loop,writableStream,"getWriter");writer:=tsCallReceiver(loop,getWriter,writableStream);closed:=false
	streamOptions:=tsNewObject();for _,name:=range []string{"highWaterMark"}{value:=tsNodeWebAdapterGet(loop,options,name);if !tsIsUndefined(value){streamOptions.set(name,value)}};streamOptions.set("objectMode",objectMode);streamOptions.set("decodeStrings",decodeStrings);signal:=tsNodeWebAdapterGet(loop,options,"signal");if !tsIsUndefined(signal){streamOptions.set("signal",signal)}
	write:=func(chunk tsValue,callback tsValue) { ready:=tsNodeWebAdapterGet(loop,writer,"ready");tsNodeWebAdapterAwait(loop,ready,func(result tsResult){if result.rejected{tsCall(loop,callback,result.value);return};method:=tsNodeWebAdapterGet(loop,writer,"write");value:=tsCallReceiver(loop,method,writer,chunk);tsNodeWebAdapterAwait(loop,value,func(result tsResult){if result.rejected{tsCall(loop,callback,result.value)}else{tsCall(loop,callback,tsU)}})}) }
	streamOptions.set("write",tsNamedNativeMethod("write",3,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{chunk,callback:=tsArg(args,0),tsArg(args,2);if chunk.kind==tsStringKind&&tsTruthy(decodeStrings)&&!tsTruthy(objectMode){encoding:=tsArg(args,1);parsed:=tsNodeBufferParseEncoding(loop,encoding,"utf8",true);chunk=tsNodeBuffer(tsFSEncode((*tsString)(chunk.ref),parsed))};write(chunk,callback);return tsU}))
	streamOptions.set("writev",tsNamedNativeMethod("writev",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{chunks,callback:=tsArg(args,0),tsArg(args,1);index:=0;var step func();step=func(){if chunks.kind!=tsArrayKind||index>=(*tsArray)(chunks.ref).length(){tsCall(loop,callback,tsU);return};request:=(*tsArray)(chunks.ref).at(index);index++;write(tsNodeWebAdapterGet(loop,request,"chunk"),tsFunctionValue(tsFunc(func(loop *tsLoop,values ...tsValue)tsValue{err:=tsArg(values,0);if tsTruthy(err){tsCall(loop,callback,err)}else{step()};return tsU}))) };step();return tsU}))
	streamOptions.set("final",tsNamedNativeMethod("final",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{callback:=tsArg(args,0);if closed{tsCall(loop,callback,tsU);return tsU};closed=true;result:=tsCallReceiver(loop,tsNodeWebAdapterGet(loop,writer,"close"),writer);tsNodeWebAdapterAwait(loop,result,func(settled tsResult){if settled.rejected{tsCall(loop,callback,settled.value)}else{tsCall(loop,callback,tsU)}});return tsU}))
	streamOptions.set("destroy",tsNamedNativeMethod("destroy",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{err,callback:=tsArg(args,0),tsArg(args,1);if closed{tsCall(loop,callback,err);return tsU};closed=true;result:=tsCallReceiver(loop,tsNodeWebAdapterGet(loop,writer,"abort"),writer,err);tsNodeWebAdapterAwait(loop,result,func(settled tsResult){tsCall(loop,callback,err)});return tsU}))
	stream:=tsNodeWritableCreate(loop,tsU,tsObjectValue(streamOptions));closedPromise:=tsNodeWebAdapterGet(loop,writer,"closed");tsNodeWebAdapterAwait(loop,closedPromise,func(result tsResult){if closed{return};closed=true;if result.rejected{tsInternalStreamDestroyer(loop,stream,result.value)}else if !tsTruthy(tsNodeWebAdapterGet(loop,stream,"writableEnded")){tsInternalStreamDestroyer(loop,stream,tsFSException("Error","ERR_STREAM_PREMATURE_CLOSE","Premature close"))}});return stream
}
func tsNodeWebAdapterGet(loop *tsLoop,value tsValue,name string) tsValue { return tsGet(loop,value,tsNodeWebAdapterKey(name)) }
func tsNodeWebAdapterNormalizeOptions(options tsValue) tsValue { if tsNullish(options){return tsObjectValue(tsNewObject())};return options }
func tsNodeWebAdapterReadableStreamFromClassic(loop *tsLoop,stream,options tsValue) tsValue {
	options=tsNodeWebAdapterNormalizeOptions(options);if tsNullish(tsNodeWebAdapterGet(loop,stream,"_readableState")){tsInternalValidatorInvalidType("streamReadable","stream.Readable",stream)};tsInternalValidatorObject(options,tsNodeWebAdapterKey("options"),tsU)
	streamOptions:=tsNodeWebAdapterGet(loop,options,"type");if !tsIsUndefined(streamOptions)&&!tsStrictEqual(streamOptions,tsStringReference(tsStringUTF8("bytes"))){tsInternalValidatorOneOf(streamOptions,tsNodeWebAdapterKey("options.type"),tsArrayValue(&tsArray{values:[]tsValue{tsStringReference(tsStringUTF8("bytes")),tsU}}))}
	if !tsInternalStreamIsReadableNode(loop,stream,false){return tsNodeWebAdapterEmptyReadable(loop,options)}
	iterator:=tsAsyncIterate(tsNodeReadableAsyncIterator(loop,stream,tsU));web:=tsWebStreamValue(iterator);object:=(*tsObject)(web.ref);object.set("_nodeReadableSource",stream);object.set("_nodeReadableType",streamOptions)
	return web
}
func tsNodeWebAdapterReadableFromWeb(loop *tsLoop,readableStream,options tsValue) tsValue {
	options=tsNodeWebAdapterNormalizeOptions(options);if !tsInternalStreamIsWebClass(readableStream,"ReadableStream"){tsInternalValidatorInvalidType("readableStream","ReadableStream",readableStream)};tsInternalValidatorObject(options,tsNodeWebAdapterKey("options"),tsU)
	if tsTruthy(tsNodeWebAdapterGet(loop,readableStream,"locked")){panic(tsThrown{tsFSException("TypeError","ERR_INVALID_STATE","ReadableStream is locked")})}
	if !tsNullish(tsNodeWebAdapterGet(loop,readableStream,"$bunNativePtr")){native:=tsNodeModuleRequire(loop,[]tsValue{tsStringReference(tsStringUTF8("internal/streams/native-readable"))});constructor:=tsNodeWebAdapterGet(loop,native,"constructNativeReadable");if constructor.kind==tsFunctionKind{return tsCall(loop,constructor,readableStream,options)}}
	stream:=tsNodeReadableCreate(loop,tsU,options);reader:=tsCallReceiver(loop,tsNodeWebAdapterGet(loop,readableStream,"getReader"),readableStream);state:=&tsNodeReadableFromWebState{stream:stream,readable:readableStream,reader:reader};tsNodeWebAdapterSingleton.Do(func(){tsNodeWebAdapterSingleton.readable=map[uintptr]*tsNodeReadableFromWebState{}});tsNodeWebAdapterSingleton.readable[uintptr(stream.ref)]=state
	object:=(*tsObject)(stream.ref);object.set("_read",tsNamedNativeMethod("_read",1,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeWebAdapterSingleton.readable[uintptr(receiver.ref)];if state==nil||state.closed{return tsU};result:=tsCallReceiver(loop,tsNodeWebAdapterGet(loop,state.reader,"read"),state.reader);tsNodeWebAdapterAwait(loop,result,func(settled tsResult){if state.closed{return};if settled.rejected{state.closed=true;tsInternalStreamErrorOrDestroy(loop,state.stream,settled.value,tsU);return};if tsTruthy(tsNodeWebAdapterGet(loop,settled.value,"done")){state.closed=true;release:=tsNodeWebAdapterGet(loop,state.reader,"releaseLock");if release.kind==tsFunctionKind{tsCallReceiver(loop,release,state.reader)};tsNodeReadablePush(loop,tsNodeReadableStateFor(state.stream),tsNull,tsU,false)}else{tsNodeReadablePush(loop,tsNodeReadableStateFor(state.stream),tsNodeWebAdapterGet(loop,settled.value,"value"),tsU,false)}});return tsU}))
	object.set("_destroy",tsNamedNativeMethod("_destroy",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeWebAdapterSingleton.readable[uintptr(receiver.ref)];callback:=tsArg(args,1);if state==nil||state.closed{tsCall(loop,callback,tsArg(args,0));return tsU};state.closed=true;result:=tsCallReceiver(loop,tsNodeWebAdapterGet(loop,state.reader,"cancel"),state.reader,tsArg(args,0));tsNodeWebAdapterAwait(loop,result,func(settled tsResult){err:=tsArg(args,0);if tsNullish(err)&&settled.rejected{err=settled.value};tsCall(loop,callback,err)});return tsU}))
	return stream
}
func tsNodeWebAdapterPairFromDuplex(loop *tsLoop,duplex,options tsValue) tsValue {
	options=tsNodeWebAdapterNormalizeOptions(options);if tsNullish(tsNodeWebAdapterGet(loop,duplex,"_readableState"))||tsNullish(tsNodeWebAdapterGet(loop,duplex,"_writableState")){tsInternalValidatorInvalidType("duplex","stream.Duplex",duplex)};tsInternalValidatorObject(options,tsNodeWebAdapterKey("options"),tsU)
	readable:=tsNodeWebAdapterEmptyReadable(loop,options);writable:=tsNodeWebAdapterWritableCreate(loop,tsNull,tsU)
	if tsTruthy(tsNodeWebAdapterGet(loop,duplex,"destroyed")){state:=tsNodeWebAdapterWritableState(writable);state.closed=true;tsNodeWebAdapterSettleClosed(state,tsResult{tsU,false});return tsObjectValue(&tsObject{values:map[string]tsValue{"readable":readable,"writable":writable},order:[]string{"readable","writable"}})}
	if tsInternalStreamIsWritableNode(loop,duplex){writable=tsNodeWebAdapterWritableStreamFromClassic(loop,duplex)}
	if tsInternalStreamIsReadableNode(loop,duplex,false){readable=tsNodeWebAdapterReadableStreamFromClassic(loop,duplex,tsNodeWebAdapterOptions(loop,options,"readableType"))}
	out:=tsNewObject();out.set("readable",readable);out.set("writable",writable);return tsObjectValue(out)
}
func tsNodeWebAdapterOptions(loop *tsLoop,options tsValue,omit string) tsValue { if options.kind!=tsObjectKind{return tsObjectValue(tsNewObject())};source:=(*tsObject)(options.ref);out:=tsNewObject();for _,key:=range source.order{if key==omit{continue};out.set(key,source.values[key])};return tsObjectValue(out) }
func tsNodeWebAdapterDuplexFromPair(loop *tsLoop,pair,options tsValue) tsValue {
	if tsNullish(pair){pair=tsObjectValue(tsNewObject())};options=tsNodeWebAdapterNormalizeOptions(options);tsInternalValidatorObject(pair,tsNodeWebAdapterKey("pair"),tsU);tsInternalValidatorObject(options,tsNodeWebAdapterKey("options"),tsU)
	readableStream,writableStream:=tsNodeWebAdapterGet(loop,pair,"readable"),tsNodeWebAdapterGet(loop,pair,"writable");if !tsInternalStreamIsWebClass(readableStream,"ReadableStream"){tsInternalValidatorInvalidType("pair.readable","ReadableStream",readableStream)};if !tsInternalStreamIsWebClass(writableStream,"WritableStream"){tsInternalValidatorInvalidType("pair.writable","WritableStream",writableStream)}
	reader:=tsCallReceiver(loop,tsNodeWebAdapterGet(loop,readableStream,"getReader"),readableStream);writer:=tsCallReceiver(loop,tsNodeWebAdapterGet(loop,writableStream,"getWriter"),writableStream);settings:=tsNodeWebAdapterOptions(loop,options,"")
	read:=tsNamedNativeMethod("read",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{result:=tsCallReceiver(loop,tsNodeWebAdapterGet(loop,reader,"read"),reader);tsNodeWebAdapterAwait(loop,result,func(settled tsResult){if settled.rejected{tsInternalStreamDestroyer(loop,receiver,settled.value)}else if tsTruthy(tsNodeWebAdapterGet(loop,settled.value,"done")){tsNodeReadablePush(loop,tsNodeReadableStateFor(receiver),tsNull,tsU,false)}else{tsNodeReadablePush(loop,tsNodeReadableStateFor(receiver),tsNodeWebAdapterGet(loop,settled.value,"value"),tsU,false)}});return tsU})
	write:=tsNamedNativeMethod("write",3,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{callback:=tsArg(args,2);result:=tsCallReceiver(loop,tsNodeWebAdapterGet(loop,writer,"write"),writer,tsArg(args,0));tsNodeWebAdapterAwait(loop,result,func(settled tsResult){if settled.rejected{tsCall(loop,callback,settled.value)}else{tsCall(loop,callback,tsU)}});return tsU})
	final:=tsNamedNativeMethod("final",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{callback:=tsArg(args,0);result:=tsCallReceiver(loop,tsNodeWebAdapterGet(loop,writer,"close"),writer);tsNodeWebAdapterAwait(loop,result,func(settled tsResult){if settled.rejected{tsCall(loop,callback,settled.value)}else{tsCall(loop,callback,tsU)}});return tsU})
	destroy:=tsNamedNativeMethod("destroy",2,func(loop *tsLoop,stream tsValue,args ...tsValue)tsValue{err,callback:=tsArg(args,0),tsArg(args,1);pending:=2;first:=tsU;done:=func(result tsResult){if result.rejected&&tsNullish(first){first=result.value};pending--;if pending==0{if tsTruthy(first){tsCall(loop,callback,first)}else{tsCall(loop,callback,err)}}};for _,entry:=range []struct{name string;target tsValue}{{"abort",writer},{"cancel",reader}}{method:=tsNodeWebAdapterGet(loop,entry.target,entry.name);result:=tsCallReceiver(loop,method,entry.target,err);tsNodeWebAdapterAwait(loop,result,done)};_ = stream;return tsU})
	tsSet(loop,settings,tsNodeWebAdapterKey("read"),read);tsSet(loop,settings,tsNodeWebAdapterKey("write"),write);tsSet(loop,settings,tsNodeWebAdapterKey("final"),final);tsSet(loop,settings,tsNodeWebAdapterKey("destroy"),destroy)
	if !tsIsUndefined(tsNodeWebAdapterGet(loop,options,"allowHalfOpen")){tsSet(loop,settings,tsNodeWebAdapterKey("allowHalfOpen"),tsNodeWebAdapterGet(loop,options,"allowHalfOpen"))}
	duplex:=tsNodeDuplexCreate(loop,tsU,settings);return duplex
}
func tsNodeWebAdapterKnownError(loop *tsLoop,cause tsValue) tsValue { code:=tsGet(loop,cause,tsNodeWebAdapterKey("code"));if tsText(code)=="ERR_STREAM_PREMATURE_CLOSE"{return tsFSException("AbortError","ABORT_ERR","The operation was aborted")};if code.kind==tsStringKind&&strings.HasPrefix(tsText(code),"ERR__ERROR_"){return tsFSException("TypeError",tsText(code),"")};return cause }
func tsNodeWebAdapterInit(loop *tsLoop) {
	tsNodeWebAdapterSingleton.Do(func(){tsNodeWebAdapterSingleton.writable=map[uintptr]*tsNodeWebWritableState{};tsNodeWebAdapterSingleton.readable=map[uintptr]*tsNodeReadableFromWebState{};tsNodeWebAdapterSingleton.errorSentinel=tsNewSymbol(tsStringUTF8("kErrorSentinelAttached"))})
}
func tsInternalWebstreamsAdaptersModule(loop *tsLoop) tsValue {
	tsNodeWebAdapterInit(loop);module:=tsNewObject()
	module.set("newWritableStreamFromStreamWritable",tsNamedNativeMethod("newWritableStreamFromStreamWritable",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeWebAdapterWritableStreamFromClassic(loop,tsArg(args,0))}))
	module.set("newReadableStreamFromStreamReadable",tsNamedNativeMethod("newReadableStreamFromStreamReadable",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeWebAdapterReadableStreamFromClassic(loop,tsArg(args,0),tsArg(args,1))}))
	module.set("newStreamWritableFromWritableStream",tsNamedNativeMethod("newStreamWritableFromWritableStream",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeWebAdapterWritableFromWeb(loop,tsArg(args,0),tsArg(args,1))}))
	module.set("newStreamReadableFromReadableStream",tsNamedNativeMethod("newStreamReadableFromReadableStream",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeWebAdapterReadableFromWeb(loop,tsArg(args,0),tsArg(args,1))}))
	module.set("newReadableWritablePairFromDuplex",tsNamedNativeMethod("newReadableWritablePairFromDuplex",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeWebAdapterPairFromDuplex(loop,tsArg(args,0),tsArg(args,1))}))
	module.set("newStreamDuplexFromReadableWritablePair",tsNamedNativeMethod("newStreamDuplexFromReadableWritablePair",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeWebAdapterDuplexFromPair(loop,tsArg(args,0),tsArg(args,1))}))
	constructor:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeWebAdapterReadableFromWeb(loop,tsArg(args,1),tsArg(args,0))});constructor.name="ReadableFromWeb";constructor.length=2;constructor.constructible=true;constructor.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeWebAdapterReadableFromWeb(loop,tsArg(args,1),tsArg(args,0))};constructor.properties=tsNewObject();constructor.properties.set("prototype",tsObjectValue(tsNewObject()));module.set("_ReadableFromWeb",tsFunctionValue(constructor));module.set("default",tsObjectValue(module));return tsObjectValue(module)
}
`
