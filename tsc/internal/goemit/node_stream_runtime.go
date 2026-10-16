package goemit

// nodeStreamRuntime assembles Bun's stream entry points over the shared
// readable, writable, duplex, iterator, pipeline, and Web Stream adapters.
const nodeStreamRuntime = `
type tsNodeDuplexPairState struct { peer, callback tsValue }
var tsNodeStreamSingleton struct { sync.Once; stream, promises, consumers, web, internal tsValue; pair map[uintptr]*tsNodeDuplexPairState }
func tsNodeStreamKey(name string) tsValue { return tsStringReference(tsStringUTF8(name)) }
func tsNodeStreamGet(loop *tsLoop,value tsValue,name string) tsValue { return tsGet(loop,value,tsNodeStreamKey(name)) }
func tsNodeStreamSet(loop *tsLoop,value tsValue,name string,item tsValue) { tsSet(loop,value,tsNodeStreamKey(name),item) }
func tsNodeStreamAwait(loop *tsLoop,value tsValue,callback func(tsResult)) { loop.await(value,callback) }
func tsNodeStreamAbort(loop *tsLoop,options tsValue) tsValue {
	signal:=tsNodeStreamGet(loop,options,"signal")
	if !tsNullish(signal)&&tsTruthy(tsNodeStreamGet(loop,signal,"aborted")){reason:=tsNodeStreamGet(loop,signal,"reason");if tsNullish(reason){reason=tsFSException("AbortError","ABORT_ERR","The operation was aborted")};return reason}
	return tsU
}
func tsNodeStreamMapIterator(loop *tsLoop,receiver,mapper,options tsValue,mode string,count int) tsValue {
	source:=tsAsyncIterate(receiver);var nested *tsAsyncIterator;remaining:=count;done:=mode=="take"&&count<=0
	iterator:=&tsAsyncIterator{}
	var next func(*tsLoop)*tsPromise
	next=func(loop *tsLoop)*tsPromise{
		promise:=loop.promise();if done{promise.settle(tsResult{tsIteratorResult(tsU,true),false});return promise}
		if reason:=tsNodeStreamAbort(loop,options);!tsIsUndefined(reason){done=true;promise.settle(tsResult{reason,true});return promise}
		var pull func(*tsLoop)
		pull=func(loop *tsLoop){
			if done{promise.settle(tsResult{tsIteratorResult(tsU,true),false});return}
			if nested!=nil{loop.await(tsPromiseValue(nested.next(loop)),func(result tsResult){if result.rejected{done=true;promise.settle(result);return};if tsTruthy(tsNodeStreamGet(loop,result.value,"done")){nested=nil;pull(loop);return};promise.settle(tsResult{tsIteratorResult(tsNodeStreamGet(loop,result.value,"value"),false),false})});return}
			loop.await(tsPromiseValue(source.next(loop)),func(result tsResult){
				if result.rejected{done=true;promise.settle(result);return};if tsTruthy(tsNodeStreamGet(loop,result.value,"done")){done=true;promise.settle(tsResult{tsIteratorResult(tsU,true),false});return};value:=tsNodeStreamGet(loop,result.value,"value")
				if mode=="drop"&&remaining>0{remaining--;pull(loop);return};if mode=="take"{if remaining<=0{done=true;promise.settle(tsResult{tsIteratorResult(tsU,true),false});return};remaining--}
				if mode=="drop"||mode=="take"{promise.settle(tsResult{tsIteratorResult(value,false),false});return}
				resultValue:=tsU;func(){defer func(){if failure:=recover();failure!=nil{done=true;promise.settle(tsResult{tsUnwrap(failure),true})}}();resultValue=tsCall(loop,mapper,value,options)}()
				loop.await(resultValue,func(mapped tsResult){if mapped.rejected{done=true;promise.settle(mapped);return};if mode=="filter"{if tsTruthy(mapped.value){promise.settle(tsResult{tsIteratorResult(value,false),false})}else{pull(loop)};return};if mode=="flatMap"{nested=tsAsyncIterate(mapped.value);pull(loop);return};promise.settle(tsResult{tsIteratorResult(mapped.value,false),false})})
			})
		}
		pull(loop);return promise
	}
	iterator.next=next;iterator.close=func(loop *tsLoop)*tsPromise{done=true;return source.close(loop)}
	return tsAsyncIteratorObject(iterator)
}
func tsNodeStreamOperator(loop *tsLoop,receiver tsValue,args []tsValue,name string) tsValue {
	var iterable tsValue
	if name=="drop"||name=="take"{number:=tsNumber(tsArg(args,0));if math.IsNaN(number){number=0};if number<0||math.IsInf(number,0){tsInternalValidatorOutOfRange("number",">= 0",tsArg(args,0))};iterable=tsNodeStreamMapIterator(loop,receiver,tsU,tsArg(args,1),name,int(number))}else{
	callback:=tsArg(args,0);if callback.kind!=tsFunctionKind{tsInternalValidatorInvalidType("fn","a function",callback)}
	iterable=tsNodeStreamMapIterator(loop,receiver,callback,tsArg(args,1),name,0)}
	module:=tsInternalStreamReadableModule(loop);constructor:=tsNodeStreamGet(loop,module,"Readable");return tsCallReceiver(loop,tsNodeStreamGet(loop,constructor,"from"),constructor,iterable)
}
func tsNodeStreamPromiseOperator(loop *tsLoop,receiver tsValue,args []tsValue,name string) tsValue {
	callback:=tsArg(args,0);options:=tsArg(args,1);if name=="reduce"{options=tsArg(args,2)};if name!="toArray"&&callback.kind!=tsFunctionKind{tsInternalValidatorInvalidType("fn","a function",callback)}
	promise:=loop.promise();iterator:=tsAsyncIterate(receiver);values:=[]tsValue{};hasInitial:=len(args)>1&&name=="reduce";accumulator:=tsArg(args,1);index:=0
	var visit func();visit=func(){
		if reason:=tsNodeStreamAbort(loop,options);!tsIsUndefined(reason){promise.settle(tsResult{reason,true});return}
		loop.await(tsPromiseValue(iterator.next(loop)),func(result tsResult){
			if result.rejected{promise.settle(result);return};if tsTruthy(tsNodeStreamGet(loop,result.value,"done")){
				switch name{case "toArray":promise.settle(tsResult{tsArrayValue(&tsArray{values:values}),false});case "reduce":if !hasInitial{promise.settle(tsResult{tsFSException("TypeError","ERR_MISSING_ARGS","Reduce of an empty stream requires an initial value"),true})}else{promise.settle(tsResult{accumulator,false})};default:if name=="every"{promise.settle(tsResult{tsBooleanValue(true),false})}else if name=="some"{promise.settle(tsResult{tsBooleanValue(false),false})}else{promise.settle(tsResult{tsU,false})}}
				return
			}
			value:=tsNodeStreamGet(loop,result.value,"value");if name=="toArray"{values=append(values,value);visit();return}
			if name=="reduce"&&!hasInitial{accumulator=value;hasInitial=true;visit();return}
			var callbackArgs []tsValue;if name=="reduce"{callbackArgs=[]tsValue{accumulator,value,options}}else{callbackArgs=[]tsValue{value,options}}
			var callbackResult tsValue;func(){defer func(){if failure:=recover();failure!=nil{promise.settle(tsResult{tsUnwrap(failure),true})}}();callbackResult=tsCall(loop,callback,callbackArgs...)}()
			loop.await(callbackResult,func(mapped tsResult){if mapped.rejected{promise.settle(mapped);return};switch name{case "reduce":accumulator=mapped.value;visit();case "forEach":visit();case "some":if tsTruthy(mapped.value){promise.settle(tsResult{tsBooleanValue(true),false})}else{visit()};case "every":if !tsTruthy(mapped.value){promise.settle(tsResult{tsBooleanValue(false),false})}else{visit()};case "find":if tsTruthy(mapped.value){promise.settle(tsResult{value,false})}else{visit()}}})
		})
	};visit();return tsPromiseValue(promise)
}
func tsNodeStreamPipelinePromises(loop *tsLoop,args []tsValue) tsValue {
	values:=append([]tsValue{},args...);options:=tsU;if len(values)>0{last:=values[len(values)-1];if (last.kind==tsObjectKind||last.kind==tsInstanceKind)&&!tsInternalStreamIsNode(loop,last)&&!tsInternalStreamIsIterable(loop,last,tsU)&&!tsInternalStreamIsWeb(loop,last){options=last;values=values[:len(values)-1]}}
	promise:=loop.promise();callback:=tsFunctionValue(tsFunc(func(_ *tsLoop,results ...tsValue)tsValue{err,value:=tsArg(results,0),tsArg(results,1);if tsTruthy(err){promise.settle(tsResult{err,true})}else{promise.settle(tsResult{value,false})};return tsU}))
	tsNodePipelineImpl(loop,tsArrayValue(&tsArray{values:values}),callback,options);return tsPromiseValue(promise)
}
func tsNodeStreamPromises(loop *tsLoop) tsValue {
	if !tsIsUndefined(tsNodeStreamSingleton.promises){return tsNodeStreamSingleton.promises}
	object:=tsNewObject();object.set("pipeline",tsNamedNativeMethod("pipeline",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeStreamPipelinePromises(loop,args)}));object.set("finished",tsNamedNativeMethod("finished",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsInternalEosFinished(loop,args)}));object.set("default",tsObjectValue(object));tsNodeStreamSingleton.promises=tsObjectValue(object);return tsNodeStreamSingleton.promises
}
func tsNodeStreamPair(loop *tsLoop,args []tsValue) tsValue {
	duplexModule:=tsInternalStreamDuplexModule(loop);constructor:=tsNodeStreamGet(loop,duplexModule,"Duplex");options:=tsArg(args,0);first:=tsConstruct(loop,constructor,options);second:=tsConstruct(loop,constructor,options)
	tsNodeStreamSingleton.pair[tsNodeCryptoPointer(first)]=&tsNodeDuplexPairState{peer:second};tsNodeStreamSingleton.pair[tsNodeCryptoPointer(second)]=&tsNodeDuplexPairState{peer:first}
	for _,side:=range []tsValue{first,second}{stream:=side;state:=tsNodeStreamSingleton.pair[tsNodeCryptoPointer(stream)];peer:=state.peer
		tsNodeStreamSet(loop,stream,"_read",tsNamedNativeMethod("_read",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{if current:=tsNodeStreamSingleton.pair[tsNodeCryptoPointer(receiver)];current!=nil&&current.callback.kind==tsFunctionKind{callback:=current.callback;current.callback=tsU;tsNodeStreamDefer(loop,func(){tsCall(loop,callback,tsU)})};return tsU}))
		tsNodeStreamSet(loop,stream,"_write",tsNamedNativeMethod("_write",3,func(loop *tsLoop,receiver tsValue,values ...tsValue)tsValue{chunk,callback:=tsArg(values,0),tsArg(values,2);current:=tsNodeStreamSingleton.pair[tsNodeCryptoPointer(receiver)];if current==nil{tsCall(loop,callback,tsFSException("Error","ERR_STREAM_DESTROYED","Duplex pair is closed"));return tsU};other:=tsNodeStreamSingleton.pair[tsNodeCryptoPointer(current.peer)];if other==nil{tsCall(loop,callback,tsFSException("Error","ERR_STREAM_DESTROYED","Duplex pair is closed"));return tsU};empty:=false;if chunk.kind==tsStringKind{empty=tsText(chunk)==""}else if chunk.kind==tsTypedArrayKind{empty=len(tsFSBytes(chunk))==0};if empty{tsNodeStreamDefer(loop,func(){tsCall(loop,callback,tsU)});return tsU};tsCallReceiver(loop,tsNodeStreamGet(loop,current.peer,"push"),current.peer,chunk);other.callback=callback;return tsU}))
		tsNodeStreamSet(loop,stream,"_final",tsNamedNativeMethod("_final",1,func(loop *tsLoop,receiver tsValue,values ...tsValue)tsValue{current:=tsNodeStreamSingleton.pair[tsNodeCryptoPointer(receiver)];if current==nil{return tsU};peer:=current.peer;callback:=tsArg(values,0);listener:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{tsCall(loop,callback,tsU);return tsU}));tsCallReceiver(loop,tsNodeStreamGet(loop,peer,"on"),peer,tsNodeStreamKey("end"),listener);tsCallReceiver(loop,tsNodeStreamGet(loop,peer,"push"),peer,tsNull);return tsU}))
		tsNodeStreamSet(loop,stream,"_destroy",tsNamedNativeMethod("_destroy",2,func(loop *tsLoop,receiver tsValue,values ...tsValue)tsValue{err,callback:=tsArg(values,0),tsArg(values,1);current:=tsNodeStreamSingleton.pair[tsNodeCryptoPointer(receiver)];if current!=nil{peer:=current.peer;tsNodeStreamDefer(loop,func(){if !tsTruthy(tsNodeStreamGet(loop,peer,"destroyed")){if tsTruthy(err){tsCallReceiver(loop,tsNodeStreamGet(loop,peer,"destroy"),peer)}else{tsCallReceiver(loop,tsNodeStreamGet(loop,peer,"push"),peer,tsNull)}}})};if callback.kind==tsFunctionKind{tsCall(loop,callback,err)};return tsU}))
		_ = state;_ = peer
	}
	return tsArrayValue(&tsArray{values:[]tsValue{first,second}})
}
func tsNodeStreamDefer(loop *tsLoop,callback func()){context:=loop.asyncContext;loop.ready=append(loop.ready,func(){tsWithAsyncContext(loop,context,callback)})}
func tsNodeStreamConsumerConvert(loop *tsLoop,method string,chunks []tsValue,settle func(tsResult)) {
	var fail func(any)
	fail=func(thrown any){settle(tsResult{tsUnwrap(thrown),true})}
	defer func(){if thrown:=recover();thrown!=nil{fail(thrown)}}()
	global:=tsObjectValue(tsGlobalObject());blobConstructor:=tsNodeStreamGet(loop,global,"Blob");blob:=tsConstruct(loop,blobConstructor,tsArrayValue(&tsArray{values:append([]tsValue(nil),chunks...)}))
	if method=="blob"{settle(tsResult{blob,false});return}
	methodName:=method;if method=="json"||method=="buffer"{methodName="text";if method=="buffer"{methodName="arrayBuffer"}}
	convert:=tsNodeStreamGet(loop,blob,methodName);result:=tsCallReceiver(loop,convert,blob)
	loop.await(result,func(converted tsResult){defer func(){if thrown:=recover();thrown!=nil{fail(thrown)}}();if converted.rejected{settle(converted);return};value:=converted.value
		if method=="json"{json:=tsNodeStreamGet(loop,global,"JSON");value=tsCallReceiver(loop,tsNodeStreamGet(loop,json,"parse"),json,value)}else if method=="buffer"{value=tsNodeBuffer(tsNodeBufferByteView(value))}
		settle(tsResult{value,false})
	})
}
func tsNodeStreamConsumerCollect(loop *tsLoop,input tsValue,method string,promise *tsPromise) {
	settled:=false;var iterator *tsAsyncIterator;finish:=func(result tsResult){if !settled{settled=true;promise.settle(result)}}
	var fail func(any)
	fail=func(thrown any){if settled{return};if iterator!=nil&&iterator.close!=nil{func(){defer func(){_ = recover()}();iterator.close(loop)}()};finish(tsResult{tsUnwrap(thrown),true})}
	safe:=func(callback func()){defer func(){if thrown:=recover();thrown!=nil{fail(thrown)}}();callback()}
	chunks:=[]tsValue{}
	if input.kind==tsObjectKind&&(*tsObject)(input.ref).webStream!=nil{stream:=(*tsObject)(input.ref).webStream;if stream.locked{finish(tsResult{tsFSException("TypeError","ERR_INVALID_STATE","Invalid state: ReadableStream is locked"),true});return};if stream.disturbed{finish(tsResult{tsFSException("TypeError","ERR_INVALID_STATE","Invalid state: ReadableStream has already been used"),true});return}}
	safe(func(){iterator=tsAsyncIterate(input)})
	if settled{return}
	var next func();next=func(){safe(func(){loop.await(tsPromiseValue(iterator.next(loop)),func(result tsResult){safe(func(){if result.rejected{finish(result);return};if tsTruthy(tsNodeStreamGet(loop,result.value,"done")){tsNodeStreamConsumerConvert(loop,method,chunks,finish);return};chunks=append(chunks,tsNodeStreamGet(loop,result.value,"value"));next()})})})};next()
}
func tsNodeStreamConsumers(loop *tsLoop) tsValue {
	if !tsIsUndefined(tsNodeStreamSingleton.consumers){return tsNodeStreamSingleton.consumers};module:=tsNewObject()
	for _,name:=range []string{"arrayBuffer","bytes","text","json","buffer","blob"}{method:=name;module.set(method,tsNamedNativeMethod(method,1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{promise:=loop.promise();tsNodeStreamConsumerCollect(loop,tsArg(args,0),method,promise);return tsPromiseValue(promise)}))};module.set("default",tsObjectValue(module));tsNodeStreamSingleton.consumers=tsObjectValue(module);return tsNodeStreamSingleton.consumers
}
func tsNodeStreamWeb(loop *tsLoop) tsValue {
	if !tsIsUndefined(tsNodeStreamSingleton.web){return tsNodeStreamSingleton.web};object:=tsNewObject();global:=tsObjectValue(tsGlobalObject());for _,name:=range []string{"ReadableStream","ReadableStreamDefaultReader","ReadableStreamBYOBReader","ReadableStreamBYOBRequest","ReadableByteStreamController","ReadableStreamDefaultController","TransformStream","TransformStreamDefaultController","WritableStream","WritableStreamDefaultWriter","WritableStreamDefaultController","ByteLengthQueuingStrategy","CountQueuingStrategy","TextEncoderStream","TextDecoderStream","CompressionStream","DecompressionStream"}{object.set(name,tsGet(loop,global,tsNodeStreamKey(name)))};object.set("default",tsObjectValue(object));tsNodeStreamSingleton.web=tsObjectValue(object);return tsNodeStreamSingleton.web
}
func tsNodeStreamModule(loop *tsLoop) tsValue {
	tsNodeStreamSingleton.Do(func(){
		tsNodeStreamSingleton.pair=map[uintptr]*tsNodeDuplexPairState{};legacy:=tsInternalStreamLegacyModule(loop);stream:=tsNodeStreamGet(loop,legacy,"Stream");if stream.kind!=tsFunctionKind{stream=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsLegacyStreamCreate(loop,tsU,tsArg(args,0))}))}
		properties:=(*tsFunction)(stream.ref).properties;if properties==nil{properties=tsNewObject();(*tsFunction)(stream.ref).properties=properties}
		readable:=tsInternalStreamReadableModule(loop);writable:=tsInternalStreamWritableModule(loop);duplex:=tsInternalStreamDuplexModule(loop);transform:=tsInternalStreamTransformModule(loop);pass:=tsInternalStreamPassThroughModule(loop)
		for _,entry:=range []struct{name string;value tsValue}{{"Readable",tsNodeStreamGet(loop,readable,"Readable")},{"Writable",tsNodeStreamGet(loop,writable,"Writable")},{"Duplex",tsNodeStreamGet(loop,duplex,"Duplex")},{"Transform",tsNodeStreamGet(loop,transform,"Transform")},{"PassThrough",tsNodeStreamGet(loop,pass,"PassThrough")}}{properties.set(entry.name,entry.value)}
		utils:=tsInternalStreamUtilsModule(loop);for _,name:=range []string{"isDestroyed","isDisturbed","isErrored","isReadable","isWritable"}{properties.set(name,tsNodeStreamGet(loop,utils,name))}
		pipeline:=tsNamedNativeMethod("pipeline",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodePipeline(loop,args)});properties.set("pipeline",pipeline)
		eos:=tsInternalStreamEndOfStreamModule(loop);properties.set("finished",eos)
		addAbort:=tsNamedNativeMethod("addAbortSignal",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsInternalStreamAddAbortSignal(loop,tsArg(args,0),tsArg(args,1))});properties.set("addAbortSignal",addAbort)
		properties.set("destroy",tsNamedNativeMethod("destroy",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsInternalStreamDestroyer(loop,tsArg(args,0),tsArg(args,1))}))
		properties.set("compose",tsNamedNativeMethod("compose",0,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeCompose(loop,args)}));properties.set("duplexPair",tsNamedNativeMethod("duplexPair",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeStreamPair(loop,args)}))
		state:=tsInternalStreamStateModule(loop);properties.set("setDefaultHighWaterMark",tsNodeStreamGet(loop,state,"setDefaultHighWaterMark"));properties.set("getDefaultHighWaterMark",tsNodeStreamGet(loop,state,"getDefaultHighWaterMark"));properties.set("promises",tsNodeStreamPromises(loop))
		properties.set("default",stream);properties.set("Stream",stream)
		promisifySymbol:=tsCall(loop,tsGet(loop,tsSymbolFunction(),tsNodeStreamKey("for")),tsStringReference(tsStringUTF8("nodejs.util.promisify.custom")));promises:=tsNodeStreamPromises(loop);(*tsFunction)(pipeline.ref).properties.set(tsPropertyKey(promisifySymbol),tsNodeStreamGet(loop,promises,"pipeline"));(*tsFunction)(eos.ref).properties.set(tsPropertyKey(promisifySymbol),tsNodeStreamGet(loop,promises,"finished"))
		readableProto:=tsNodeStreamGet(loop,tsNodeStreamGet(loop,stream,"Readable"),"prototype");for _,name:=range []string{"drop","filter","flatMap","map","take"}{method:=name;readableObject:=(*tsObject)(readableProto.ref);readableObject.set(method,tsNamedNativeMethod(method,2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeStreamOperator(loop,receiver,args,method)}))}
		for _,name:=range []string{"every","find","forEach","reduce","some","toArray"}{method:=name;readableObject:=(*tsObject)(readableProto.ref);readableObject.set(method,tsNamedNativeMethod(method,2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeStreamPromiseOperator(loop,receiver,args,method)}))}
		tsNodeStreamSingleton.stream=stream
	});return tsNodeStreamSingleton.stream
}
func tsInternalStreamModule(loop *tsLoop) tsValue {
	if !tsIsUndefined(tsNodeStreamSingleton.internal) { return tsNodeStreamSingleton.internal }
	module := tsNewObject()
	stream := tsNodeStreamModule(loop)
	module.set("Stream", stream)
	module.set("default", stream)
	tsNodeStreamSingleton.internal = tsObjectValue(module)
	return tsNodeStreamSingleton.internal
}
func tsNodeStreamAddon(loop *tsLoop,name string)(tsValue,bool){switch name{case "stream":return tsNodeStreamModule(loop),true;case "stream/promises":return tsNodeStreamPromises(loop),true;case "stream/consumers":return tsNodeStreamConsumers(loop),true;case "stream/web":return tsNodeStreamWeb(loop),true};return tsU,false}
`
