package goemit

// internalStreamIterClassicRuntime is the classic Node stream bridge used by
// stream/iter. The shared setup is kept here before its lazy Writable edge.
const internalStreamIterClassicRuntime = `
var tsInternalStreamIterClassic struct {
	sync.Once
	module *tsObject
	readable tsValue
	writable tsValue
	symbols map[string]tsValue
	fromReadable map[uintptr]tsValue
	fromWritable map[uintptr]map[string]tsValue
	maxBatch int
}

func tsInternalStreamIterClassicInit(loop *tsLoop) {
	tsInternalStreamIterClassic.Do(func(){
		// Source lines 14-35: bind the already mapped shared stream utilities
		// and the async/sync iteration symbols.
		_ = tsInternalErrorsModule(loop)
		_ = tsInternalValidatorsModule(loop)
		_ = tsInternalStreamEndOfStreamModule(loop)
		_ = tsInternalStreamAddAbortSignalModule(loop)
		_ = tsInternalStreamIterTypesModule(loop)
		_ = tsInternalStreamIterUtilsModule(loop)
		_ = tsInternalStreamDestroyModule(loop)
		types:=tsInternalStreamIterTypesModule(loop)
		symbols:=map[string]tsValue{}
		for _,name:=range []string{"asyncIterator","iterator","asyncDispose","dispose"}{symbols[name]=tsWellKnownSymbol(name)}
		for _,name:=range []string{"toAsyncStreamable","kValidatedSource","kSyncWriteAccepted","drainableProtocol"}{symbols[name]=tsGet(loop,types,tsStringReference(tsStringUTF8(name)))}
		// Promise.withResolvers is represented by a local Promise plus its
		// resolving functions; callers settle it through the normal loop.
		module:=tsNewObject()
		module.set("__promiseWithResolvers",tsNamedNativeMethod("withResolvers",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{
			promise:=loop.promise();out:=tsNewObject();out.set("promise",tsPromiseValue(promise));out.set("resolve",tsFunctionValue(tsFunc(func(_ *tsLoop,args ...tsValue)tsValue{promise.settle(tsResult{tsArg(args,0),false});return tsU})));out.set("reject",tsFunctionValue(tsFunc(func(_ *tsLoop,args ...tsValue)tsValue{promise.settle(tsResult{tsArg(args,0),true});return tsU})));return tsObjectValue(out)
		}))
		tsInternalStreamIterClassic.module=module
		tsInternalStreamIterClassic.symbols=symbols
		tsInternalStreamIterClassic.fromReadable=map[uintptr]tsValue{}
		tsInternalStreamIterClassic.fromWritable=map[uintptr]map[string]tsValue{}
		tsInternalStreamIterClassic.maxBatch=128
		module.set("__lazyReadable",tsNamedNativeMethod("lazyReadable",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsInternalStreamIterClassicLazyReadable(loop)}))
		module.set("__lazyWritable",tsNamedNativeMethod("lazyWritable",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsInternalStreamIterClassicLazyWritable(loop)}))
	})
}

func tsInternalStreamIterClassicLazyReadable(loop *tsLoop) tsValue {
	if tsNullish(tsInternalStreamIterClassic.readable){tsInternalStreamIterClassic.readable=tsNodeModuleRequire(loop,[]tsValue{tsStringReference(tsStringUTF8("internal/streams/readable"))})}
	return tsInternalStreamIterClassic.readable
}
func tsInternalStreamIterClassicLazyWritable(loop *tsLoop) tsValue {
	if tsNullish(tsInternalStreamIterClassic.writable){tsInternalStreamIterClassic.writable=tsNodeModuleRequire(loop,[]tsValue{tsStringReference(tsStringUTF8("internal/streams/writable"))})}
	return tsInternalStreamIterClassic.writable
}

type tsNodeClassicReadableIteratorState struct { stream tsValue; normalize bool; ended bool; errorValue tsValue; waiter *tsPromise; cleanup tsValue; closed bool }
type tsNodeClassicDrainWaiter struct { promise *tsPromise; drainable bool }
type tsNodeClassicWriterState struct { writable tsValue; policy string; hwm,total int; waiters []tsNodeClassicDrainWaiter; listeners bool; onDrain,onError tsValue; ended bool }

func tsNodeClassicKey(name string) tsValue { return tsStringReference(tsStringUTF8(name)) }
func tsNodeClassicGet(loop *tsLoop,value tsValue,name string) tsValue { return tsGet(loop,value,tsNodeClassicKey(name)) }
func tsNodeClassicMethod(loop *tsLoop,value tsValue,name string,args ...tsValue) tsValue { method:=tsNodeClassicGet(loop,value,name);if method.kind!=tsFunctionKind{return tsU};return tsCallReceiver(loop,method,value,args...) }
func tsNodeClassicOptions(options tsValue) tsValue { if tsIsUndefined(options){return tsObjectValue(tsNewObject())};if options.kind!=tsObjectKind&&options.kind!=tsInstanceKind&&options.kind!=tsArrayKind{tsNodeClassicInvalid("options","an object",options)};return options }
func tsNodeClassicPromise(loop *tsLoop,value tsValue) tsValue { return tsPromiseValue(loop.resolved(value,false)) }
func tsNodeClassicReject(loop *tsLoop,value tsValue) tsValue { promise:=loop.promise();promise.settle(tsResult{value,true});return tsPromiseValue(promise) }
func tsNodeClassicPointer(value tsValue) uintptr { if value.kind!=tsObjectKind&&value.kind!=tsInstanceKind&&value.kind!=tsArrayKind{return 0};return uintptr(value.ref) }
func tsNodeClassicInvalid(name,expected string,value tsValue) { _=value;panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_TYPE","The \""+name+"\" argument must be "+expected)}) }
func tsNodeClassicHWM(loop *tsLoop,options tsValue,def int) int { value:=tsNodeClassicGet(loop,options,"highWaterMark");if tsIsUndefined(value){return def};if !tsIsNumeric(value)||tsNumber(value)<0||math.Trunc(tsNumber(value))!=tsNumber(value){panic(tsThrown{tsFSException("RangeError","ERR_INVALID_OPT_VALUE","The value of \"options.highWaterMark\" is out of range")})};return int(tsNumber(value)) }
func tsNodeClassicEmit(loop *tsLoop,state *tsNodeClassicReadableIteratorState) { if state.waiter!=nil{state.waiter.settle(tsResult{tsU,false});state.waiter=nil} }
func tsNodeClassicRemove(loop *tsLoop,stream tsValue,event string,listener tsValue) { if tsNodeClassicGet(loop,stream,"removeListener").kind==tsFunctionKind{tsNodeClassicMethod(loop,stream,"removeListener",tsNodeClassicKey(event),listener)}else{tsNodeClassicMethod(loop,stream,"off",tsNodeClassicKey(event),listener)} }
func tsNodeClassicNormalizeBatch(loop *tsLoop,raw tsValue) *tsPromise {
	promise:=loop.promise();if raw.kind!=tsArrayKind{promise.settle(tsResult{tsFSException("TypeError","ERR_INVALID_ARG_TYPE","The \"raw\" argument must be an Array"),true});return promise}
	items:=(*tsArray)(raw.ref).values;chunks:=[]tsValue{};var visit func(int);visit=func(index int){if index>=len(items){if len(chunks)==0{promise.settle(tsResult{tsNull,false})}else{promise.settle(tsResult{tsNodeIterFromArray(chunks),false})};return};value:=items[index];if tsNodeUtilType("isUint8Array",value){chunks=append(chunks,value);visit(index+1);return};tsNodeIterFromAsyncValue(loop,value,nil,func(next []tsValue,err tsValue){if tsTruthy(err){promise.settle(tsResult{err,true});return};chunks=append(chunks,next...);visit(index+1)})};visit(0);return promise
}
func tsNodeClassicBatchedReadable(stream,normalize tsValue) tsValue {
	state:=&tsNodeClassicReadableIteratorState{stream:stream,normalize:!tsNullish(normalize),errorValue:tsU}
	readable:=tsNamedNativeMethod("readable",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{tsNodeClassicEmit(loop,state);return tsU})
	ended:=tsNamedNativeMethod("end",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{state.ended=true;tsNodeClassicEmit(loop,state);return tsU})
	closed:=tsNamedNativeMethod("close",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{state.ended=true;tsNodeClassicEmit(loop,state);return tsU})
	errored:=tsNamedNativeMethod("error",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{state.errorValue=tsArg(args,0);state.ended=true;tsNodeClassicEmit(loop,state);return tsU})
	tsNodeClassicMethod(loop,stream,"on",tsNodeClassicKey("readable"),readable);tsNodeClassicMethod(loop,stream,"on",tsNodeClassicKey("end"),ended);tsNodeClassicMethod(loop,stream,"on",tsNodeClassicKey("close"),closed);tsNodeClassicMethod(loop,stream,"on",tsNodeClassicKey("error"),errored)
	cleanup:=tsNamedNativeMethod("cleanup",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{if !state.closed{state.closed=true;tsNodeClassicRemove(loop,stream,"readable",readable);tsNodeClassicRemove(loop,stream,"end",ended);tsNodeClassicRemove(loop,stream,"close",closed);tsNodeClassicRemove(loop,stream,"error",errored)};return tsU});state.cleanup=cleanup
	iterator:=&tsAsyncIterator{}
	iterator.next=func(loop *tsLoop)*tsPromise{promise:=loop.promise();var advance func();advance=func(){if state.errorValue.kind!=tsUndefinedKind&&state.errorValue.kind!=tsNullKind{promise.settle(tsResult{state.errorValue,true});return};if state.ended{promise.settle(tsResult{tsIteratorResult(tsU,true),false});return};if tsTruthy(tsNodeClassicGet(loop,stream,"destroyed")){state.ended=true;promise.settle(tsResult{tsIteratorResult(tsU,true),false});return};chunk:=tsNodeClassicMethod(loop,stream,"read");if !tsIsNull(chunk)&&!tsIsUndefined(chunk){batch:=[]tsValue{chunk};for len(batch)<128&&tsNumber(tsNodeClassicGet(loop,tsNodeClassicGet(loop,stream,"_readableState"),"length"))>0{next:=tsNodeClassicMethod(loop,stream,"read");if tsNullish(next){break};batch=append(batch,next)};raw:=tsNodeIterFromArray(batch);if !state.normalize{promise.settle(tsResult{tsIteratorResult(raw,false),false});return};loop.await(tsPromiseValue(tsNodeClassicNormalizeBatch(loop,raw)),func(result tsResult){if result.rejected{promise.settle(result);return};if tsNullish(result.value){advance();return};promise.settle(tsResult{tsIteratorResult(result.value,false),false})});return};state.waiter=loop.promise();loop.await(tsPromiseValue(state.waiter),func(tsResult){advance()})};advance();return promise}
	iterator.close=func(loop *tsLoop)*tsPromise{state.ended=true;tsCallReceiver(loop,state.cleanup,tsU);return loop.resolved(tsIteratorResult(tsU,true),false)}
	iterator.resume=func(loop *tsLoop,action string,value tsValue)*tsPromise{if action=="return"{return iterator.close(loop)};return iterator.next(loop)}
	return tsAsyncIteratorObject(iterator)
}
func tsNodeClassicFromReadable(loop *tsLoop,readable tsValue) tsValue {
	if tsNullish(readable)||readable.kind!=tsObjectKind&&readable.kind!=tsInstanceKind{tsNodeClassicInvalid("readable","Readable",readable)}
	key:=tsNodeClassicPointer(readable);if cached,ok:=tsInternalStreamIterClassic.fromReadable[key];ok{return cached}
	symbol:=tsInternalStreamIterClassic.symbols["toAsyncStreamable"];if !tsNullish(tsGet(loop,readable,symbol)){result:=tsCallReceiver(loop,tsGet(loop,readable,symbol),readable);tsInternalStreamIterClassic.fromReadable[key]=result;return result}
	if tsNodeClassicGet(loop,readable,"read").kind!=tsFunctionKind||tsNodeClassicGet(loop,readable,"on").kind!=tsFunctionKind{tsNodeClassicInvalid("readable","Readable",readable)}
	readableState:=tsNodeClassicGet(loop,readable,"_readableState");normalize:=tsNull;if !tsNullish(readableState)&&(tsTruthy(tsNodeClassicGet(loop,readableState,"objectMode"))||!tsIsUndefined(tsNodeClassicGet(loop,readableState,"encoding"))&&tsNodeClassicGet(loop,readableState,"encoding").kind!=tsNullKind){normalize=tsBooleanValue(true)}
	iterator:=tsNodeClassicBatchedReadable(readable,normalize);keyValidated:=tsInternalStreamIterClassic.symbols["kValidatedSource"];tsSet(loop,iterator,keyValidated,tsBooleanValue(true));tsSet(loop,iterator,tsNodeClassicKey("stream"),readable);tsInternalStreamIterClassic.fromReadable[key]=iterator;return iterator
}
func tsNodeClassicCallbackPromise(loop *tsLoop,fn func(tsValue)tsValue) tsValue { promise:=loop.promise();value:=fn(tsFunctionValue(tsFunc(func(_ *tsLoop,args ...tsValue)tsValue{err:=tsArg(args,0);if tsTruthy(err){promise.settle(tsResult{err,true})}else{promise.settle(tsResult{tsU,false})};return tsU})));return value }
func tsNodeClassicResolveWaiters(state *tsNodeClassicWriterState,err tsValue,drain bool) { pending:=state.waiters;state.waiters=nil;for _,waiter:=range pending{value:=true;if err.kind!=tsUndefinedKind&&err.kind!=tsNullKind{value=false};if waiter.drainable{waiter.promise.settle(tsResult{tsBooleanValue(value),false})}else if value{waiter.promise.settle(tsResult{tsU,false})}else{waiter.promise.settle(tsResult{err,true})}};_ = drain }
func tsNodeClassicWriterListeners(loop *tsLoop,state *tsNodeClassicWriterState) {
	if state.listeners{return};state.listeners=true
	state.onDrain=tsNamedNativeMethod("drain",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{tsNodeClassicResolveWaiters(state,tsNull,true);return tsU})
	state.onError=tsNamedNativeMethod("error",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{tsNodeClassicResolveWaiters(state,tsArg(args,0),false);return tsU})
	tsNodeClassicMethod(loop,state.writable,"on",tsNodeClassicKey("drain"),state.onDrain);tsNodeClassicMethod(loop,state.writable,"on",tsNodeClassicKey("error"),state.onError)
}
func tsNodeClassicWriterCleanup(loop *tsLoop,state *tsNodeClassicWriterState,err tsValue) {
	if tsNullish(err){err=tsFSException("AbortError","ABORT_ERR","The operation was aborted")};tsNodeClassicResolveWaiters(state,err,false);if state.listeners{state.listeners=false;tsNodeClassicRemove(loop,state.writable,"drain",state.onDrain);tsNodeClassicRemove(loop,state.writable,"error",state.onError)}
}
func tsNodeClassicWriterWait(loop *tsLoop,state *tsNodeClassicWriterState,drainable bool) tsValue { promise:=loop.promise();state.waiters=append(state.waiters,tsNodeClassicDrainWaiter{promise:promise,drainable:drainable});tsNodeClassicWriterListeners(loop,state);return tsPromiseValue(promise) }
func tsNodeClassicWriterWritable(loop *tsLoop,state *tsNodeClassicWriterState) bool {return !tsTruthy(tsNodeClassicGet(loop,state.writable,"destroyed"))&&!tsTruthy(tsNodeClassicGet(loop,state.writable,"writableFinished"))&&!tsTruthy(tsNodeClassicGet(loop,state.writable,"writableEnded"))}
func tsNodeClassicWriterFull(loop *tsLoop,state *tsNodeClassicWriterState) bool {n:=tsNodeClassicGet(loop,state.writable,"writableLength");length:=0;if !tsIsUndefined(n)&&!tsNullish(n){length=int(tsNumber(n))};return length>=state.hwm}
func tsNodeClassicWriter(loop *tsLoop,writable,options tsValue) tsValue {
	if tsNullish(writable)||tsNodeClassicGet(loop,writable,"write").kind!=tsFunctionKind||tsNodeClassicGet(loop,writable,"on").kind!=tsFunctionKind{tsNodeClassicInvalid("writable","Writable",writable)}
	options=tsNodeClassicOptions(options)
	policyValue:=tsNodeClassicGet(loop,options,"backpressure");policy:="strict";if !tsIsUndefined(policyValue){policy=tsText(policyValue)};tsInternalStreamIterUtilsValidateBackpressure(loop,tsStringReference(tsStringUTF8(policy)))
	if tsTruthy(tsNodeClassicGet(loop,writable,"writableObjectMode")){panic(tsThrown{tsFSException("TypeError","ERR_INVALID_STATE","Cannot create a stream/iter Writer from an object-mode Writable")})};if policy=="drop-oldest"{panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_VALUE","drop-oldest is not supported for classic stream.Writable")})}
	ptr:=tsNodeClassicPointer(writable);byPolicy:=tsInternalStreamIterClassic.fromWritable[ptr];if byPolicy==nil{byPolicy=map[string]tsValue{};tsInternalStreamIterClassic.fromWritable[ptr]=byPolicy};if cached,ok:=byPolicy[policy];ok{return cached}
	hwmValue:=tsNodeClassicGet(loop,writable,"writableHighWaterMark");hwm:=16384;if !tsIsUndefined(hwmValue)&&!tsNullish(hwmValue){hwm=int(tsNumber(hwmValue))};state:=&tsNodeClassicWriterState{writable:writable,policy:policy,hwm:hwm}
	writer:=tsNewObject();writer.descriptors=map[string]*tsDescriptor{}
	writer.descriptors["desiredSize"]=&tsDescriptor{accessor:true,configurable:true,getter:tsNamedNativeMethod("get desiredSize",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{if !tsNodeClassicWriterWritable(loop,state){return tsNull};n:=tsNodeClassicGet(loop,writable,"writableLength");length:=0;if !tsIsUndefined(n)&&!tsNullish(n){length=int(tsNumber(n))};if hwm-length<0{return tsNumberValue(0)};return tsNumberValue(float64(hwm-length))})}
	writer.set("writeSync",tsNamedNativeMethod("writeSync",1,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsBooleanValue(false)}));writer.set("writevSync",tsNamedNativeMethod("writevSync",1,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsBooleanValue(false)}));writer.set("endSync",tsNamedNativeMethod("endSync",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNumberValue(-1)}))
	writer.set("write",tsNamedNativeMethod("write",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{if !tsNodeClassicWriterWritable(loop,state){return tsNodeClassicReject(loop,tsFSException("Error","ERR_STREAM_WRITE_AFTER_END","write after end"))};chunk:=tsInternalStreamIterUtilsToUint8Array(loop,tsArg(args,0));length:=len(tsFSBytes(chunk));if policy=="strict"&&tsNodeClassicWriterFull(loop,state){return tsNodeClassicReject(loop,tsFSException("RangeError","ERR_INVALID_STATE","Backpressure violation: buffer is full. Await each write() call to respect backpressure."))};if policy=="drop-newest"&&tsNodeClassicWriterFull(loop,state){state.total+=length;return tsNodeClassicPromise(loop,tsU)};state.total+=length;ok:=tsNodeClassicMethod(loop,writable,"write",chunk);if tsTruthy(ok){return tsNodeClassicPromise(loop,tsU)};if policy=="block"{return tsNodeClassicWriterWait(loop,state,false)};return tsNodeClassicPromise(loop,tsU)}))
	writer.set("writev",tsNamedNativeMethod("writev",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{chunks:=tsArg(args,0);if chunks.kind!=tsArrayKind{tsNodeClassicInvalid("chunks","an Array",chunks)};if !tsNodeClassicWriterWritable(loop,state){return tsNodeClassicReject(loop,tsFSException("Error","ERR_STREAM_WRITE_AFTER_END","write after end"))};array:=(*tsArray)(chunks.ref);if policy=="strict"&&tsNodeClassicWriterFull(loop,state){return tsNodeClassicReject(loop,tsFSException("RangeError","ERR_INVALID_STATE","Backpressure violation: buffer is full. Await each write() call to respect backpressure."))};if policy=="drop-newest"&&tsNodeClassicWriterFull(loop,state){for _,chunk:=range array.values{state.total+=len(tsFSBytes(tsInternalStreamIterUtilsToUint8Array(loop,chunk)))};return tsNodeClassicPromise(loop,tsU)};cork:=tsNodeClassicGet(loop,writable,"cork");uncork:=tsNodeClassicGet(loop,writable,"uncork");if cork.kind==tsFunctionKind&&uncork.kind==tsFunctionKind{tsCallReceiver(loop,cork,writable)};ok:=true;for _,item:=range array.values{chunk:=tsInternalStreamIterUtilsToUint8Array(loop,item);state.total+=len(tsFSBytes(chunk));if !tsTruthy(tsNodeClassicMethod(loop,writable,"write",chunk)){ok=false}};if cork.kind==tsFunctionKind&&uncork.kind==tsFunctionKind{tsCallReceiver(loop,uncork,writable)};if !ok&&policy=="block"{return tsNodeClassicWriterWait(loop,state,false)};return tsNodeClassicPromise(loop,tsU)}))
	writer.set("end",tsNamedNativeMethod("end",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{if tsTruthy(tsNodeClassicGet(loop,writable,"writableFinished"))||tsTruthy(tsNodeClassicGet(loop,writable,"destroyed")){tsNodeClassicWriterCleanup(loop,state,tsNull);return tsNodeClassicPromise(loop,tsNumberValue(float64(state.total)))};promise:=loop.promise();tsNodeClassicMethod(loop,writable,"end");options:=tsNewObject();options.set("writable",tsBooleanValue(true));options.set("readable",tsBooleanValue(false));callback:=tsFunctionValue(tsFunc(func(_ *tsLoop,args ...tsValue)tsValue{err:=tsArg(args,0);tsNodeClassicWriterCleanup(loop,state,err);if tsTruthy(err){promise.settle(tsResult{err,true})}else{promise.settle(tsResult{tsNumberValue(float64(state.total)),false})};return tsU}));tsInternalEos(loop,[]tsValue{writable,tsObjectValue(options),callback});return tsPromiseValue(promise)}))
	writer.set("fail",tsNamedNativeMethod("fail",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{reason:=tsArg(args,0);tsNodeClassicWriterCleanup(loop,state,reason);if tsNodeClassicGet(loop,writable,"destroy").kind==tsFunctionKind{tsNodeClassicMethod(loop,writable,"destroy",reason)};return tsU}))
	writer.set(tsPropertyKey(tsWellKnownSymbol("asyncDispose")),tsNamedNativeMethod("[Symbol.asyncDispose]",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{if tsNodeClassicWriterWritable(loop,state){tsNodeClassicWriterCleanup(loop,state,tsU);tsNodeClassicMethod(loop,writable,"destroy")};return tsNodeClassicPromise(loop,tsU)}))
	writer.set(tsPropertyKey(tsWellKnownSymbol("dispose")),tsNamedNativeMethod("[Symbol.dispose]",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{if tsNodeClassicWriterWritable(loop,state){tsNodeClassicWriterCleanup(loop,state,tsU);tsNodeClassicMethod(loop,writable,"destroy")};return tsU}))
	writer.set(tsPropertyKey(tsInternalStreamIterClassic.symbols["drainableProtocol"]),tsNamedNativeMethod("[drainableProtocol]",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{if !tsNodeClassicWriterWritable(loop,state){return tsNull};if !tsNodeClassicWriterFull(loop,state){return tsNodeClassicPromise(loop,tsBooleanValue(true))};return tsNodeClassicWriterWait(loop,state,true)}))
	byPolicy[policy]=tsObjectValue(writer);return tsObjectValue(writer)
}
func tsNodeClassicToReadable(loop *tsLoop,source,options tsValue) tsValue {
	if !tsNodeClassicIsAsyncIterable(loop,source){tsNodeClassicInvalid("source","AsyncIterable",source)}
	options=tsNodeClassicOptions(options)
	hwm:=tsNodeClassicHWM(loop,options,65536);iterator:=tsAsyncIterate(source);done,pumping:=false,false;var backpressure *tsPromise;var readable tsValue
	var pump func(*tsLoop);pump=func(loop *tsLoop){if pumping||done{return};pumping=true;var step func();step=func(){if done{pumping=false;return};loop.await(tsPromiseValue(iterator.next(loop)),func(result tsResult){if result.rejected{done=true;pumping=false;tsNodeClassicMethod(loop,readable,"destroy",result.value);return};if tsTruthy(tsNodeClassicGet(loop,result.value,"done")){done=true;pumping=false;tsNodeClassicMethod(loop,readable,"push",tsNull);return};batch:=tsNodeClassicGet(loop,result.value,"value");if batch.kind!=tsArrayKind{done=true;pumping=false;tsNodeClassicMethod(loop,readable,"destroy",tsFSException("TypeError","ERR_INVALID_ARG_TYPE","stream/iter source must yield Uint8Array batches"));return};values:=(*tsArray)(batch.ref).values;var pushAt func(int);pushAt=func(index int){if done{pumping=false;return};if index>=len(values){step();return};accepted:=tsNodeClassicMethod(loop,readable,"push",values[index]);if !tsTruthy(accepted){pumping=false;backpressure=loop.promise();loop.await(tsPromiseValue(backpressure),func(tsResult){if !done{pumping=true;pushAt(index+1)}});return};pushAt(index+1)};pushAt(0)})};step()}
	read:=tsNamedNativeMethod("read",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{if backpressure!=nil{backpressure.settle(tsResult{tsU,false});backpressure=nil}else{pump(loop)};return tsU});destroy:=tsNamedNativeMethod("destroy",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{done=true;if backpressure!=nil{backpressure.settle(tsResult{tsU,false});backpressure=nil};callback:=tsArg(args,1);finish:=func(err tsValue){if callback.kind==tsFunctionKind{tsCall(loop,callback,err)}};loop.await(tsPromiseValue(iterator.close(loop)),func(result tsResult){err:=tsArg(args,0);if result.rejected&&tsNullish(err){err=result.value};finish(err)});return tsU})
	streamOptions:=tsNewObject();streamOptions.set("highWaterMark",tsNumberValue(float64(hwm)));streamOptions.set("read",read);streamOptions.set("destroy",destroy);readable=tsConstruct(tsGet(loop,tsInternalStreamIterClassicLazyReadable(loop),tsNodeClassicKey("Readable")),tsObjectValue(streamOptions));signal:=tsNodeClassicGet(loop,options,"signal");if !tsNullish(signal){tsInternalStreamAddAbortSignal(loop,signal,readable)};return readable
}
func tsNodeClassicIsAsyncIterable(loop *tsLoop,value tsValue)bool{return !tsNullish(tsGet(loop,value,tsWellKnownSymbol("asyncIterator")))}
func tsNodeClassicToReadableSync(loop *tsLoop,source,options tsValue) tsValue {
	if tsNullish(tsGet(loop,source,tsWellKnownSymbol("iterator"))){tsNodeClassicInvalid("source","Iterable",source)};options=tsNodeClassicOptions(options);hwm:=tsNodeClassicHWM(loop,options,65536);iterator:=tsIterate(source);var batch tsValue;batchIndex:=0;hasBatch:=false;var readable tsValue
	read:=tsNamedNativeMethod("read",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{for{if hasBatch{array:=(*tsArray)(batch.ref);for batchIndex<array.length(){item:=array.at(batchIndex);batchIndex++;if !tsTruthy(tsNodeClassicMethod(loop,receiver,"push",item)){return tsU}};hasBatch=false;batch=tsU;batchIndex=0};if !iterator.next(){tsNodeClassicMethod(loop,receiver,"push",tsNull);return tsU};batch=iterator.value;if batch.kind!=tsArrayKind{tsNodeClassicInvalid("source","an iterable yielding arrays of Uint8Array chunks",batch)};hasBatch=true}})
	destroy:=tsNamedNativeMethod("destroy",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{batch=tsU;hasBatch=false;callback:=tsArg(args,1);if callback.kind==tsFunctionKind{tsCall(loop,callback,tsArg(args,0))};return tsU});streamOptions:=tsNewObject();streamOptions.set("highWaterMark",tsNumberValue(float64(hwm)));streamOptions.set("read",read);streamOptions.set("destroy",destroy);readable=tsConstruct(tsGet(loop,tsInternalStreamIterClassicLazyReadable(loop),tsNodeClassicKey("Readable")),tsObjectValue(streamOptions));return readable
}
func tsNodeClassicToWritable(loop *tsLoop,writer tsValue) tsValue {
	if tsNodeClassicGet(loop,writer,"write").kind!=tsFunctionKind{tsNodeClassicInvalid("writer","Writer",writer)};writableModule:=tsInternalStreamIterClassicLazyWritable(loop);constructor:=tsGet(loop,writableModule,tsNodeClassicKey("Writable"));write:=tsNamedNativeMethod("write",3,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{chunk,encoding,callback:=tsArg(args,0),tsArg(args,1),tsArg(args,2);if chunk.kind==tsStringKind{chunk=tsNodeBuffer(tsFSEncode((*tsString)(chunk.ref),tsText(encoding)))};writeSync:=tsNodeClassicGet(loop,writer,"writeSync");if writeSync.kind==tsFunctionKind{syncResult:=tsCallReceiver(loop,writeSync,writer,chunk);if tsTruthy(syncResult){tsNodeWritableDefer(loop,func(){if callback.kind==tsFunctionKind{tsCall(loop,callback,tsU)}});return tsU}};promise:=tsNodeClassicMethod(loop,writer,"write",chunk);loop.await(promise,func(result tsResult){if callback.kind==tsFunctionKind{if result.rejected{tsCall(loop,callback,result.value)}else{tsCall(loop,callback,tsU)}}});return tsU})
	writev:=tsNamedNativeMethod("writev",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{entries,callback:=tsArg(args,0),tsArg(args,1);chunks:=&tsArray{};if entries.kind==tsArrayKind{for _,entry:=range (*tsArray)(entries.ref).values{chunk:=tsNodeClassicGet(loop,entry,"chunk");encoding:=tsNodeClassicGet(loop,entry,"encoding");if chunk.kind==tsStringKind{chunk=tsNodeBuffer(tsFSEncode((*tsString)(chunk.ref),tsText(encoding)))};chunks.values=append(chunks.values,chunk)}};method:=tsNodeClassicGet(loop,writer,"writev");if method.kind!=tsFunctionKind{for _,chunk:=range chunks.values{tsNodeClassicMethod(loop,writer,"write",chunk)};if callback.kind==tsFunctionKind{tsNodeWritableDefer(loop,func(){tsCall(loop,callback,tsU)})};return tsU};promise:=tsCallReceiver(loop,method,writer,tsArrayValue(chunks));loop.await(promise,func(result tsResult){if callback.kind==tsFunctionKind{if result.rejected{tsCall(loop,callback,result.value)}else{tsCall(loop,callback,tsU)}}});return tsU})
	final:=tsNamedNativeMethod("final",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{callback:=tsArg(args,0);endSync:=tsNodeClassicGet(loop,writer,"endSync");if endSync.kind==tsFunctionKind&&tsNumber(tsCallReceiver(loop,endSync,writer))>=0{tsNodeWritableDefer(loop,func(){if callback.kind==tsFunctionKind{tsCall(loop,callback,tsU)}});return tsU};end:=tsNodeClassicGet(loop,writer,"end");if end.kind!=tsFunctionKind{if callback.kind==tsFunctionKind{tsCall(loop,callback,tsU)};return tsU};loop.await(tsCallReceiver(loop,end,writer),func(result tsResult){if callback.kind==tsFunctionKind{if result.rejected{tsCall(loop,callback,result.value)}else{tsCall(loop,callback,tsU)}}});return tsU})
	destroy:=tsNamedNativeMethod("destroy",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{err,callback:=tsArg(args,0),tsArg(args,1);if tsTruthy(err)&&tsNodeClassicGet(loop,writer,"fail").kind==tsFunctionKind{tsNodeClassicMethod(loop,writer,"fail",err)};if callback.kind==tsFunctionKind{tsCall(loop,callback,tsU)};return tsU});options:=tsNewObject();options.set("highWaterMark",tsNumberValue(9007199254740991));options.set("write",write);if tsNodeClassicGet(loop,writer,"writev").kind==tsFunctionKind{options.set("writev",writev)};options.set("final",final);options.set("destroy",destroy);return tsConstruct(constructor,tsObjectValue(options))
}

func tsInternalStreamIterClassicModule(loop *tsLoop) tsValue {
	tsInternalStreamIterClassicInit(loop)
	module:=tsInternalStreamIterClassic.module
	module.set("normalizeBatch",tsNamedNativeMethod("normalizeBatch",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsPromiseValue(tsNodeClassicNormalizeBatch(loop,tsArg(args,0)))}))
	module.set("createBatchedAsyncIterator",tsNamedNativeMethod("createBatchedAsyncIterator",2,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeClassicBatchedReadable(tsArg(args,0),tsArg(args,1))}))
	module.set("fromReadable",tsNamedNativeMethod("fromReadable",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeClassicFromReadable(loop,tsArg(args,0))}))
	module.set("fromWritable",tsNamedNativeMethod("fromWritable",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeClassicWriter(loop,tsArg(args,0),tsArg(args,1))}))
	module.set("toReadable",tsNamedNativeMethod("toReadable",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeClassicToReadable(loop,tsArg(args,0),tsArg(args,1))}))
	module.set("toReadableSync",tsNamedNativeMethod("toReadableSync",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeClassicToReadableSync(loop,tsArg(args,0),tsArg(args,1))}))
	module.set("toWritable",tsNamedNativeMethod("toWritable",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeClassicToWritable(loop,tsArg(args,0))}))
	defaultValue:=tsNewObject();for _,name:=range []string{"createBatchedAsyncIterator","normalizeBatch","fromReadable","fromWritable","toReadable","toReadableSync","toWritable"}{defaultValue.set(name,tsGet(loop,module,tsNodeClassicKey(name)))};module.set("default",tsObjectValue(defaultValue))
	return tsObjectValue(module)
}
`
