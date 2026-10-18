package goemit

// internalStreamNativeReadableRuntime maps the native-readable wrapper onto
// Readable streams and a duck-typed native pointer object. Bun's C++ handle
// transfer itself has no host equivalent in this Go runtime.
const internalStreamNativeReadableRuntime = `
type tsNodeNativeReadableState struct {
	stream, pointer, closeState, remaining tsValue
	refCount, highWaterMark, debugID int
	constructed, pending, resized bool
}
var tsNodeNativeReadableSingleton struct { sync.Once; states map[uintptr]*tsNodeNativeReadableState; nextID int }
func tsNodeNativeReadableKey(name string) tsValue { return tsStringReference(tsStringUTF8(name)) }
func tsNodeNativeReadableStateFor(stream tsValue) *tsNodeNativeReadableState { if stream.kind!=tsObjectKind{return nil};return tsNodeNativeReadableSingleton.states[uintptr(stream.ref)] }
func tsNodeNativeReadablePointerCall(loop *tsLoop,state *tsNodeNativeReadableState,name string,args ...tsValue) tsValue {
	if state==nil||state.pointer.kind!=tsObjectKind{return tsU};method:=tsGet(loop,state.pointer,tsNodeNativeReadableKey(name));if method.kind!=tsFunctionKind{return tsU};return tsCallReceiver(loop,method,state.pointer,args...)
}
func tsNodeNativeReadableConstruct(loop *tsLoop,readableStream,options tsValue) tsValue {
	if readableStream.kind!=tsObjectKind||!tsInternalStreamIsWebClass(readableStream,"ReadableStream"){tsFSInvalid("readableStream","a ReadableStream")}
	pointer:=tsGet(loop,readableStream,tsNodeNativeReadableKey("$bunNativePtr"));if pointer.kind!=tsObjectKind{tsFSInvalid("readableStream.$bunNativePtr","an object")}
	stream:=tsNodeReadableCreate(loop,tsU,options);state:=&tsNodeNativeReadableState{stream:stream,pointer:pointer,closeState:tsArrayValue(&tsArray{values:[]tsValue{tsBooleanValue(false)}}),highWaterMark:256*1024}
	if value:=tsNodeNativeReadableKey("highWaterMark");!tsIsUndefined(value){candidate:=tsGet(loop,options,value);if tsIsNumeric(candidate){state.highWaterMark=int(tsNumber(candidate))}}
	state.resized=false;state.remaining=tsU;tsNodeNativeReadableSingleton.Do(func(){tsNodeNativeReadableSingleton.states=map[uintptr]*tsNodeNativeReadableState{}});tsNodeNativeReadableSingleton.states[uintptr(stream.ref)]=state
	object:=(*tsObject)(stream.ref);object.set("$bunNativePtr",pointer);tsNodeNativeReadableSingleton.nextID++;state.debugID=tsNodeNativeReadableSingleton.nextID;object.set("debugId",tsNumberValue(float64(state.debugID)))
	object.set("_read",tsNamedNativeMethod("_read",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeNativeReadableRead(loop,tsNodeNativeReadableStateFor(receiver),int(tsNumber(tsArg(args,0))))}))
	object.set("_destroy",tsNamedNativeMethod("_destroy",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{tsNodeNativeReadableDestroy(loop,tsNodeNativeReadableStateFor(receiver),tsArg(args,0),tsArg(args,1));return tsU}))
	object.set("ref",tsNamedNativeMethod("ref",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{return tsNodeNativeReadableRef(loop,tsNodeNativeReadableStateFor(receiver),true)}))
	object.set("unref",tsNamedNativeMethod("unref",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{return tsNodeNativeReadableRef(loop,tsNodeNativeReadableStateFor(receiver),false)}))
	closeState:=state.closeState;object.set("_nativeCloseState",closeState)
	return stream
}
func tsNodeNativeReadableHasClosed(state *tsNodeNativeReadableState) bool { if state==nil||state.closeState.kind!=tsArrayKind{return false};array:=(*tsArray)(state.closeState.ref);return array.length()>0&&tsTruthy(array.at(0)) }
func tsNodeNativeReadableView(value tsValue) bool { return value.kind==tsTypedArrayKind||value.kind==tsDataViewKind }
func tsNodeNativeReadablePush(loop *tsLoop,state *tsNodeNativeReadableState,value tsValue) {
	if state==nil{return};accepted:=tsNodeReadablePush(loop,tsNodeReadableStateFor(state.stream),value,tsU,false);if !tsTruthy(accepted){tsNodeNativeReadablePointerCall(loop,state,"setFlowing",tsBooleanValue(false))}
}
func tsNodeNativeReadableEOF(loop *tsLoop,state *tsNodeNativeReadableState) { if state!=nil&&!tsTruthy(tsNodeDuplexGet(loop,state.stream,"destroyed")){tsNodeNativeReadablePush(loop,state,tsNull)} }
func tsNodeNativeReadableHandleResult(loop *tsLoop,state *tsNodeNativeReadableState,result,chunk tsValue) tsValue {
	closed:=tsNodeNativeReadableHasClosed(state)
	if tsIsNumeric(result) {
		n:=int(tsNumber(result));if n>0&&chunk.kind!=tsUndefinedKind {data:=tsFSBytes(chunk);if n>len(data){n=len(data)};tsNodeNativeReadablePush(loop,state,tsNodeBuffer(append([]byte(nil),data[:n]...)));if n<len(data){state.remaining=tsNodeBuffer(append([]byte(nil),data[n:]...))}else{state.remaining=tsU}}
		if n>=state.highWaterMark&&!state.resized&&!closed{state.highWaterMark*=2;if state.highWaterMark>2*1024*1024{state.highWaterMark=2*1024*1024};state.resized=true}
	} else if result.kind==tsBooleanKind {
		if chunk.kind!=tsUndefinedKind&&len(tsFSBytes(chunk))>0{state.remaining=chunk};context:=loop.asyncContext;loop.ready=append(loop.ready,func(){tsWithAsyncContext(loop,context,func(){tsNodeNativeReadableEOF(loop,state)})})
	} else if tsNodeNativeReadableView(result) {
		data:=tsFSBytes(result);if len(data)>0{tsNodeNativeReadablePush(loop,state,result)};if len(data)>=state.highWaterMark&&!state.resized&&!closed{state.highWaterMark*=2;if state.highWaterMark>2*1024*1024{state.highWaterMark=2*1024*1024};state.resized=true}
	} else if !tsIsUndefined(result) { panic(tsThrown{tsFSException("TypeError","ERR_INVALID_RETURN_VALUE","Invalid result from native readable pull")}) }
	if closed {context:=loop.asyncContext;loop.ready=append(loop.ready,func(){tsWithAsyncContext(loop,context,func(){tsNodeNativeReadableEOF(loop,state)})})}
	return state.remaining
}
func tsNodeNativeReadableRead(loop *tsLoop,state *tsNodeNativeReadableState,maxToRead int) tsValue {
	if state==nil{return tsU};tsNodeNativeReadablePointerCall(loop,state,"setFlowing",tsBooleanValue(true));if state.pending{return tsU}
	if state.pointer.kind!=tsObjectKind {tsNodeNativeReadableEOF(loop,state);return tsU}
	if !state.constructed {
		started:=tsNodeNativeReadablePointerCall(loop,state,"start",tsNumberValue(float64(state.highWaterMark)))
		if tsIsNumeric(started)&&tsNumber(started)>1 {state.resized=true;if int(tsNumber(started))<state.highWaterMark{state.highWaterMark=int(tsNumber(started))}}
		if tsNodeNativeReadableView(started)&&len(tsFSBytes(started))>0{tsNodeNativeReadablePush(loop,state,started)}
		drain:=tsNodeNativeReadablePointerCall(loop,state,"drain");if tsNodeNativeReadableView(drain)&&len(tsFSBytes(drain))>0{tsNodeNativeReadablePush(loop,state,drain)};state.constructed=true
	}
	if state.remaining.kind!=tsUndefinedKind&&state.remaining.kind!=tsNullKind&&len(tsFSBytes(state.remaining))>0 {remaining:=state.remaining;state.remaining=tsU;tsNodeNativeReadablePush(loop,state,remaining)}
	capacity:=state.highWaterMark;if maxToRead>capacity&&maxToRead>0{capacity=maxToRead};if capacity<512{capacity=512};chunk:=tsNodeBufferAllocate(loop,[]tsValue{tsNumberValue(float64(capacity))},false)
	result:=tsNodeNativeReadablePointerCall(loop,state,"pull",chunk,state.closeState);if tsIsUndefined(result){panic(tsThrown{tsFSException("TypeError","ERR_INVALID_RETURN_VALUE","NativeReadable.pull returned undefined")})}
	if result.kind==tsPromiseKind {state.pending=true;promise:=(*tsPromise)(result.ref);loop.await(tsPromiseValue(promise),func(res tsResult){state.pending=false;if res.rejected{tsInternalStreamErrorOrDestroy(loop,state.stream,res.value,tsU);return};tsNodeNativeReadableHandleResult(loop,state,res.value,chunk)});return tsU}
	state.remaining=tsNodeNativeReadableHandleResult(loop,state,result,chunk);return tsU
}
func tsNodeNativeReadableDestroy(loop *tsLoop,state *tsNodeNativeReadableState,errorValue,callback tsValue) {
	if state!=nil {tsNodeNativeReadablePointerCall(loop,state,"cancel",errorValue);readable:=tsNodeReadableStateFor(state.stream);if readable!=nil&&readable.flowing&&!readable.ended&&!readable.sync{readable.buffer=nil;readable.length=0;tsNodeReadableSyncState(readable)}}
	if callback.kind==tsFunctionKind {cb:=callback;context:=loop.asyncContext;loop.ready=append(loop.ready,func(){tsWithAsyncContext(loop,context,func(){tsCall(loop,cb,errorValue)})})}
}
func tsNodeNativeReadableRef(loop *tsLoop,state *tsNodeNativeReadableState,ref bool) tsValue {
	if state==nil{return tsU};previous:=state.refCount;if ref{state.refCount++}else{state.refCount--};if ref&&previous==0||!ref&&previous==1{tsNodeNativeReadablePointerCall(loop,state,"updateRef",tsBooleanValue(ref))};return state.stream
}
func tsInternalStreamNativeReadableModule(loop *tsLoop) tsValue {
	module:=tsNewObject();module.set("constructNativeReadable",tsNamedNativeMethod("constructNativeReadable",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeNativeReadableConstruct(loop,tsArg(args,0),tsArg(args,1))}));module.set("default",tsObjectValue(module));return tsObjectValue(module)
}
`
