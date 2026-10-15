package goemit

// internalStreamReadableRuntime translates the stateful Readable core. State
// stays on the JS side of the adapter as both the public stream API and
// Node's internal stream helpers inspect these fields.
const internalStreamReadableRuntime = `
type tsNodeReadableState struct {
	stream, object tsValue
	buffer []tsValue
	length, highWaterMark int
	objectMode, ended, endEmitted, reading, needReadable bool
	flowing, flowingSet, paused, destroyed, closed, errored bool
	constructed bool
	errorValue, decoder tsValue
	encoding, defaultEncoding string
	readHook tsValue
	readable bool
	close, autoDestroy bool
	sync bool
	pipes []tsValue
	awaitDrain int
	dataEmitted bool
}

var tsNodeReadableSingleton struct {
	sync.Once
	constructor *tsFunction
	prototype *tsObject
	states map[uintptr]*tsNodeReadableState
}

func tsNodeReadableKey(name string) tsValue { return tsStringReference(tsStringUTF8(name)) }
func tsNodeReadableGet(loop *tsLoop, object tsValue, name string) tsValue { return tsGet(loop, object, tsNodeReadableKey(name)) }
func tsNodeReadableSet(loop *tsLoop, object tsValue, name string, value tsValue) { tsSet(loop, object, tsNodeReadableKey(name), value) }
func tsNodeReadableStateFor(stream tsValue) *tsNodeReadableState {
	if stream.kind != tsObjectKind { return nil }
	return tsNodeReadableSingleton.states[uintptr(stream.ref)]
}
func tsNodeReadableSyncState(state *tsNodeReadableState) {
	if state == nil || state.object.kind != tsObjectKind { return }
	object := (*tsObject)(state.object.ref)
	object.set("buffer", tsArrayValue(&tsArray{values:append([]tsValue{},state.buffer...)}))
	object.set("bufferIndex", tsNumberValue(0))
	object.set("length", tsNumberValue(float64(state.length)))
	object.set("highWaterMark", tsNumberValue(float64(state.highWaterMark)))
	object.set("objectMode", tsBooleanValue(state.objectMode))
	object.set("ended", tsBooleanValue(state.ended))
	object.set("endEmitted", tsBooleanValue(state.endEmitted))
	object.set("reading", tsBooleanValue(state.reading))
	object.set("needReadable", tsBooleanValue(state.needReadable))
	object.set("flowing", tsBooleanValue(state.flowing))
	object.set("destroyed", tsBooleanValue(state.destroyed))
	object.set("closed", tsBooleanValue(state.closed))
	object.set("errored", state.errorValue)
	object.set("readable", tsBooleanValue(state.readable))
	object.set("constructed",tsBooleanValue(state.constructed))
	object.set("sync",tsBooleanValue(state.sync))
	object.set("autoDestroy", tsBooleanValue(state.autoDestroy))
	object.set("emitClose", tsBooleanValue(state.close))
	object.set("pipes", tsArrayValue(&tsArray{values:append([]tsValue{},state.pipes...)}))
	object.set("pipesCount", tsNumberValue(float64(len(state.pipes))))
	object.set("dataEmitted", tsBooleanValue(state.dataEmitted))
	bits := uint32(0)
	if state.objectMode { bits |= 1 << 0 }; if state.autoDestroy { bits |= 1 << 2 }; if state.close { bits |= 1 << 3 }
	if state.destroyed { bits |= 1 << 4 }; if state.closed { bits |= 1 << 5 }; if state.errorValue.kind != tsNullKind && !tsIsUndefined(state.errorValue) { bits |= 1 << 7 }; if state.constructed {bits |= 1 << 8}
	object.set(tsPropertyKey(tsInternalStreamUtilsKey(loop,"kState")),tsNumberValue(float64(bits)))
}

func tsNodeReadableEmit(loop *tsLoop, state *tsNodeReadableState, event string, args ...tsValue) tsValue {
	if state == nil { return tsBooleanValue(false) }
	return tsCallReceiver(loop,tsNodeReadableGet(loop,state.stream,"emit"),state.stream,append([]tsValue{tsNodeReadableKey(event)},args...)...)
}
func tsNodeReadableDefer(loop *tsLoop, callback func()) { context:=loop.asyncContext; loop.ready=append(loop.ready,func(){tsWithAsyncContext(loop,context,callback)}) }
func tsNodeReadableError(loop *tsLoop,state *tsNodeReadableState,errorValue tsValue) {
	if state == nil { return }
	if tsNullish(errorValue) { errorValue=tsFSException("TypeError","ERR_STREAM_NULL_VALUES","May not write null values to stream") }
	state.errorValue=errorValue; state.reading=false; tsNodeReadableSyncState(state)
	if state.autoDestroy { tsNodeReadableDestroy(loop,state,errorValue,tsU); return }
	tsNodeReadableDefer(loop,func(){tsNodeReadableEmit(loop,state,"error",errorValue)})
}
func tsNodeReadableEnd(loop *tsLoop,state *tsNodeReadableState) {
	if state==nil || state.endEmitted || !state.ended || state.length!=0 { return }
	tsNodeReadableDefer(loop,func(){
		if state.endEmitted || state.length!=0 { return }
		state.endEmitted=true; state.readable=false; tsNodeReadableSyncState(state); tsNodeReadableEmit(loop,state,"end")
		if state.autoDestroy { tsNodeReadableDestroy(loop,state,tsU,tsU) }
	})
}
func tsNodeReadableScheduleRead(loop *tsLoop,state *tsNodeReadableState) {
	if state==nil || state.reading || state.ended || state.destroyed || !state.constructed || !state.readable { return }
	state.reading=true; tsNodeReadableSyncState(state)
	readHook:=state.readHook
	if readHook.kind!=tsFunctionKind { readHook=tsNodeReadableGet(loop,state.stream,"_read") }
	if readHook.kind==tsFunctionKind {
		func(){ defer func(){ if failure:=recover(); failure!=nil { state.reading=false; tsNodeReadableError(loop,state,tsUnwrap(failure)) } }(); tsCallReceiver(loop,readHook,state.stream,tsNumberValue(float64(state.highWaterMark))) }()
	} else { state.reading=false }
}
func tsNodeReadableHWM(loop *tsLoop,options tsValue) int {
	objectMode:=tsTruthy(tsNodeReadableGet(loop,options,"objectMode"))
	state:=tsObjectValue(tsNewObject()); tsSet(loop,state,tsNodeReadableKey("objectMode"),tsBooleanValue(objectMode))
	hwm:=tsCallReceiver(loop,tsGet(loop,tsInternalStreamStateModule(loop),tsNodeReadableKey("getHighWaterMark")),tsInternalStreamStateModule(loop),state,options,tsNodeReadableKey("readableHighWaterMark"),tsBooleanValue(false))
	return int(tsNumber(hwm))
}
func tsNodeReadableCreate(loop *tsLoop, receiver, options tsValue) tsValue {
	tsNodeReadableInit(loop)
	if tsNullish(options) { options = tsObjectValue(tsNewObject()) }
	var stream tsValue
	if receiver.kind==tsObjectKind { stream=tsLegacyStreamCreate(loop,receiver,options) } else { stream=tsLegacyStreamCreate(loop,tsU,options) }
	object:=(*tsObject)(stream.ref); object.prototype=tsObjectValue(tsNodeReadableSingleton.prototype)
	stateObject:=tsNewObject(); object.set("_readableState",tsObjectValue(stateObject))
	defaultEncoding:="utf8"; encoding:=tsNodeReadableGet(loop,options,"defaultEncoding"); if !tsNullish(encoding) { defaultEncoding=tsText(encoding) }
	if tsNodeBufferEncoding(loop,tsStringReference(tsStringUTF8(defaultEncoding)))=="unsupported" { panic(tsThrown{tsFSException("TypeError","ERR_UNKNOWN_ENCODING","Unknown encoding: "+defaultEncoding)}) }
	state:=&tsNodeReadableState{stream:stream,object:tsObjectValue(stateObject),highWaterMark:tsNodeReadableHWM(loop,options),objectMode:tsTruthy(tsNodeReadableGet(loop,options,"objectMode")),defaultEncoding:defaultEncoding,readable:true,constructed:true,sync:true,close:!tsStrictEqual(tsNodeReadableGet(loop,options,"emitClose"),tsBooleanValue(false)),autoDestroy:!tsStrictEqual(tsNodeReadableGet(loop,options,"autoDestroy"),tsBooleanValue(false)),errorValue:tsNull}
	state.readHook=tsNodeReadableGet(loop,options,"read"); if state.readHook.kind==tsFunctionKind { tsNodeReadableSet(loop,stream,"_read",state.readHook) }
	state.decoder=tsU
	if enc:=tsNodeReadableGet(loop,options,"encoding"); !tsNullish(enc) { state.encoding=tsText(enc); decoderModule:=tsNodeStringDecoderModule(); state.decoder=tsConstruct(tsGet(loop,decoderModule,tsNodeReadableKey("StringDecoder")),enc) }
	tsNodeReadableSingleton.states[uintptr(stream.ref)]=state
	stateObject.set("stream",stream); stateObject.set("buffer",tsArrayValue(&tsArray{})); stateObject.set("bufferIndex",tsNumberValue(0)); stateObject.set("pipes",tsArrayValue(&tsArray{})); stateObject.set("awaitDrainWriters",tsNull);stateObject.set("defaultEncoding",tsStringReference(tsStringUTF8(defaultEncoding)))
	if state.encoding!="" { stateObject.set("encoding",tsStringReference(tsStringUTF8(state.encoding)));stateObject.set("decoder",state.decoder) }
	tsNodeReadableSyncState(state)
	if tsNodeReadableGet(loop,options,"destroy").kind==tsFunctionKind { tsNodeReadableSet(loop,stream,"_destroy",tsNodeReadableGet(loop,options,"destroy")) }
	if construct:=tsNodeReadableGet(loop,options,"construct");construct.kind==tsFunctionKind { tsNodeReadableSet(loop,stream,"_construct",construct);state.constructed=false;tsNodeReadableSyncState(state);callback:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{state.constructed=true;if tsTruthy(tsArg(args,0)){tsNodeReadableError(loop,state,tsArg(args,0))}else{tsNodeReadableSyncState(state);if state.needReadable{tsNodeReadableScheduleRead(loop,state)}};return tsU}));tsInternalStreamConstruct(loop,stream,callback) }
	if signal:=tsNodeReadableGet(loop,options,"signal"); !tsNullish(signal) { tsInternalStreamAddAbortSignal(loop,signal,stream) }
	return stream
}
func tsNodeReadableInit(loop *tsLoop) {
	tsNodeReadableSingleton.Do(func(){
		tsNodeReadableSingleton.states=map[uintptr]*tsNodeReadableState{}
		legacy:=tsInternalStreamLegacyModule(loop); base:=tsGet(loop,legacy,tsNodeReadableKey("Stream")); baseProto:=tsGet(loop,base,tsNodeReadableKey("prototype"))
		prototype:=tsNewObject(); prototype.prototype=baseProto; tsNodeReadableSingleton.prototype=prototype
		constructor:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeReadableCreate(loop,tsU,tsArg(args,0))}); constructor.name="Readable"; constructor.length=1; constructor.constructible=true
		constructor.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeReadableCreate(loop,tsU,tsArg(args,0))}
		constructor.receiverCall=func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeReadableCreate(loop,receiver,tsArg(args,0))}
		constructor.properties=tsNewObject(); constructor.properties.set("prototype",tsObjectValue(prototype)); constructor.properties.descriptors=map[string]*tsDescriptor{"prototype":{writable:false}}
		prototype.set("constructor",tsFunctionValue(constructor)); prototype.descriptors=map[string]*tsDescriptor{"constructor":{writable:true,configurable:true}}
		tsNodeReadableSingleton.constructor=constructor
		tsNodeReadableInstallPrototype(loop)
	})
}

func tsNodeReadablePush(loop *tsLoop,state *tsNodeReadableState,chunk,encoding tsValue,front bool) tsValue {
	if state==nil { return tsBooleanValue(false) }
	if tsIsNull(chunk) { state.reading=false; if !front&&!state.ended { state.ended=true; if state.decoder.kind==tsObjectKind { tail:=tsCallReceiver(loop,tsGet(loop,state.decoder,tsNodeReadableKey("end")),state.decoder); if !tsIsUndefined(tail)&&tsText(tail)!="" { state.buffer=append(state.buffer,tail); state.length+=len((*tsString)(tail.ref).units) } }; tsNodeReadableSyncState(state); if state.length>0 { tsNodeReadableEmit(loop,state,"readable") } else { tsNodeReadableEnd(loop,state) } }; return tsBooleanValue(false) }
	if tsIsUndefined(chunk) { return tsBooleanValue(false) }
	if state.destroyed || state.ended { return tsBooleanValue(false) }
	if !state.objectMode {
		if chunk.kind==tsStringKind { enc:=state.defaultEncoding; if !tsIsUndefined(encoding)&&!tsIsNull(encoding){enc=tsText(encoding)}; if state.encoding==""||state.encoding!=enc { chunk=tsNodeBufferFrom(loop,[]tsValue{chunk,tsStringReference(tsStringUTF8(enc))}); encoding=tsU } }
		if chunk.kind==tsDataViewKind { chunk=tsNodeBufferFrom(loop,[]tsValue{chunk}) }
		if chunk.kind!=tsTypedArrayKind && !(chunk.kind==tsStringKind&&state.encoding!="") { tsNodeReadableError(loop,state,tsFSException("TypeError","ERR_INVALID_ARG_TYPE","The \"chunk\" argument must be of type string or an instance of Buffer, TypedArray, or DataView")); return tsBooleanValue(false) }
		if chunk.kind==tsTypedArrayKind && (*tsTypedArray)(chunk.ref).length==0 { state.reading=false; return tsBooleanValue(state.length<state.highWaterMark||state.length==0) }
		if state.decoder.kind==tsObjectKind && tsIsUndefined(encoding) { chunk=tsCallReceiver(loop,tsGet(loop,state.decoder,tsNodeReadableKey("write")),state.decoder,chunk); if tsText(chunk)=="" {return tsBooleanValue(state.length<state.highWaterMark||state.length==0)} }
	}
	length:=1; if !state.objectMode { if chunk.kind==tsStringKind {length=len((*tsString)(chunk.ref).units)} else {length=len(tsFSBytes(chunk))} }
	if front {state.buffer=append([]tsValue{chunk},state.buffer...)}else{state.buffer=append(state.buffer,chunk)}; state.length+=length;state.reading=false
	if state.flowing && !front { for state.length>0 { value:=tsNodeReadableTake(loop,state,0);if tsNullish(value){break};state.dataEmitted=true;tsNodeReadableEmit(loop,state,"data",value) } } else if state.needReadable {state.needReadable=false;tsNodeReadableEmit(loop,state,"readable")}
	tsNodeReadableSyncState(state); if !front&&state.length<state.highWaterMark {tsNodeReadableScheduleRead(loop,state)}
	return tsBooleanValue(!state.ended&&(state.length<state.highWaterMark||state.length==0))
}
func tsNodeReadableTake(loop *tsLoop,state *tsNodeReadableState,n int) tsValue {
	if state.length==0||len(state.buffer)==0{return tsNull}
	if state.objectMode {value:=state.buffer[0];state.buffer=state.buffer[1:];state.length--;return value}
	if n<=0||n>state.length {n=state.length}; if state.encoding!="" { if n>state.length {n=state.length}; remaining:=n;var out strings.Builder;for remaining>0&&len(state.buffer)>0 {chunk:=state.buffer[0];text:="";if chunk.kind==tsStringKind{text=tsText(chunk)}else{text=tsText(tsNodeBufferToString(loop,chunk,[]tsValue{tsStringReference(tsStringUTF8("utf8"))}))};if len(text)<=remaining {out.WriteString(text);remaining-=len(text);state.length-=len(text);state.buffer=state.buffer[1:]} else {out.WriteString(text[:remaining]);state.buffer[0]=tsStringReference(tsStringUTF8(text[remaining:]));state.length-=remaining;remaining=0} };return tsStringReference(tsStringUTF8(out.String())) }
	if len(state.buffer)==1 {chunk:=state.buffer[0];size:=len(tsFSBytes(chunk));if n==size {state.buffer=nil;state.length=0;return chunk};data:=tsFSBytes(chunk);result:=tsNodeBuffer(append([]byte(nil),data[:n]...));state.buffer[0]=tsNodeBuffer(append([]byte(nil),data[n:]...));state.length-=n;return result}
	data:=make([]byte,0,n);for len(data)<n&&len(state.buffer)>0 {chunk:=state.buffer[0];chunkData:=tsFSBytes(chunk);take:=n-len(data);if take>=len(chunkData){data=append(data,chunkData...);state.buffer=state.buffer[1:]}else{data=append(data,chunkData[:take]...);state.buffer[0]=tsNodeBuffer(append([]byte(nil),chunkData[take:]...))};state.length-=take};return tsNodeBuffer(data)
}
func tsNodeReadableRead(loop *tsLoop,state *tsNodeReadableState,nv tsValue) tsValue {
	if state==nil{return tsNull};n:=0
	if !tsIsUndefined(nv) { number:=tsNumber(nv);if math.IsNaN(number){n=0}else{n=int(math.Trunc(number))} }
	if n>state.highWaterMark {state.highWaterMark=int(math.Ceil(math.Log2(float64(n)))) ;state.highWaterMark=1<<state.highWaterMark}
	if len(state.buffer)==0&&!state.ended {state.needReadable=true;tsNodeReadableScheduleRead(loop,state)}
	if n<=0 {n=state.length};if state.objectMode&&state.length>0 {n=1};if n>state.length {n=state.length}
	if n==0 {if state.ended {tsNodeReadableEnd(loop,state)}else{state.needReadable=true};return tsNull}
	result:=tsNodeReadableTake(loop,state,n);if !tsNullish(result){state.dataEmitted=true;tsNodeReadableEmit(loop,state,"data",result)}
	if state.length==0&&state.ended {tsNodeReadableEnd(loop,state)};tsNodeReadableSyncState(state);return result
}
func tsNodeReadableDestroy(loop *tsLoop,state *tsNodeReadableState,err,callback tsValue) {
	if state==nil{return};if state.destroyed {if callback.kind==tsFunctionKind {tsCall(loop,callback,err)};return};state.destroyed=true;state.readable=false;if tsTruthy(err){state.errorValue=err};tsNodeReadableSyncState(state)
	destroy:=tsNodeReadableGet(loop,state.stream,"_destroy");if destroy.kind==tsFunctionKind {done:=false;cb:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{if done{return tsU};done=true;failure:=tsArg(args,0);state.closed=true;if tsTruthy(failure){state.errorValue=failure};tsNodeReadableSyncState(state);if callback.kind==tsFunctionKind{tsCall(loop,callback,failure)};if tsTruthy(failure){tsNodeReadableDefer(loop,func(){tsNodeReadableEmit(loop,state,"error",failure)})};if state.close{tsNodeReadableDefer(loop,func(){tsNodeReadableEmit(loop,state,"close")})};return tsU}));tsCallReceiver(loop,destroy,state.stream,err,cb);return}
	state.closed=true;tsNodeReadableSyncState(state);if tsTruthy(err){tsNodeReadableDefer(loop,func(){tsNodeReadableEmit(loop,state,"error",err)})};if state.close{tsNodeReadableDefer(loop,func(){tsNodeReadableEmit(loop,state,"close")})};if callback.kind==tsFunctionKind{tsCall(loop,callback,err)}
}

func tsNodeReadableStateCompat(loop *tsLoop,options tsValue)tsValue{
	state:=tsNewObject();state.set("objectMode",tsBooleanValue(tsTruthy(tsNodeReadableGet(loop,options,"objectMode"))));state.set("highWaterMark",tsNumberValue(float64(tsNodeReadableHWM(loop,options))));state.set("buffer",tsArrayValue(&tsArray{}));state.set("bufferIndex",tsNumberValue(0));state.set("length",tsNumberValue(0));state.set("pipes",tsArrayValue(&tsArray{}));state.set("ended",tsBooleanValue(false));state.set("endEmitted",tsBooleanValue(false));state.set("destroyed",tsBooleanValue(false));state.set("errored",tsNull);state.set("closed",tsBooleanValue(false));state.set("reading",tsBooleanValue(false));state.set("needReadable",tsBooleanValue(false));state.set("sync",tsBooleanValue(true))
	stateKey:=tsInternalStreamUtilsKey(loop,"kState");state.set(tsPropertyKey(stateKey),tsNumberValue(0));state.descriptors=map[string]*tsDescriptor{}
	state.descriptors["pipesCount"]=&tsDescriptor{accessor:true,configurable:true,getter:tsNamedNativeMethod("get pipesCount",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{pipes:=tsGet(loop,receiver,tsNodeReadableKey("pipes"));if pipes.kind!=tsArrayKind{return tsNumberValue(0)};return tsNumberValue(float64((*tsArray)(pipes.ref).length()))})}
	state.descriptors["paused"]=&tsDescriptor{accessor:true,configurable:true,getter:tsNamedNativeMethod("get paused",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{bits:=uint32(tsNumber(tsGet(loop,receiver,stateKey)));return tsBooleanValue(bits&(1<<26)!=0)})};state.descriptors["paused"].setter=tsNamedNativeMethod("set paused",1,func(_ *tsLoop,receiver tsValue,args ...tsValue)tsValue{bits:=uint32(tsNumber(tsGet(loop,receiver,stateKey)));bits|=1<<25;if tsTruthy(tsArg(args,0)){bits|=1<<26}else{bits &^=1<<26};tsSet(loop,receiver,stateKey,tsNumberValue(float64(bits)));return tsU})
	return tsObjectValue(state)
}

func tsNodeReadableInstallPrototype(loop *tsLoop) {
	p:=tsNodeReadableSingleton.prototype
	define:=func(name string,length int,fn func(*tsLoop,tsValue,...tsValue)tsValue){p.set(name,tsNamedNativeMethod(name,length,fn));p.descriptors[name]=&tsDescriptor{writable:true,configurable:true}}
	define("_read",1,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{panic(tsThrown{tsFSException("Error","ERR_METHOD_NOT_IMPLEMENTED","The _read() method is not implemented")})})
	define("_destroy",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{if tsArg(args,1).kind==tsFunctionKind{tsCall(loop,tsArg(args,1),tsArg(args,0))};return tsU})
	define("push",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeReadablePush(loop,tsNodeReadableStateFor(receiver),tsArg(args,0),tsArg(args,1),false)})
	define("unshift",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeReadablePush(loop,tsNodeReadableStateFor(receiver),tsArg(args,0),tsArg(args,1),true)})
	define("read",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeReadableRead(loop,tsNodeReadableStateFor(receiver),tsArg(args,0))})
	define("destroy",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{tsNodeReadableDestroy(loop,tsNodeReadableStateFor(receiver),tsArg(args,0),tsArg(args,1));return receiver})
	define("_undestroy",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{if state:=tsNodeReadableStateFor(receiver);state!=nil{state.destroyed=false;state.closed=false;state.ended=false;state.endEmitted=false;state.errorValue=tsNull;tsNodeReadableSyncState(state)};return receiver})
	define("isPaused",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeReadableStateFor(receiver);return tsBooleanValue(state==nil||state.paused||(!state.flowingSet||!state.flowing))})
	define("setEncoding",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeReadableStateFor(receiver);if state==nil{return receiver};enc:=tsText(tsArg(args,0));state.encoding=enc;module:=tsNodeStringDecoderModule();state.decoder=tsConstruct(tsGet(loop,module,tsNodeReadableKey("StringDecoder")),tsArg(args,0));var all strings.Builder;for _,item:=range state.buffer{if item.kind==tsStringKind{all.WriteString(tsText(item))}else{all.WriteString(tsText(tsCallReceiver(loop,tsGet(loop,state.decoder,tsNodeReadableKey("write")),state.decoder,item)))}};state.buffer=nil;state.length=0;if all.Len()>0{chunk:=tsStringReference(tsStringUTF8(all.String()));state.buffer=append(state.buffer,chunk);state.length=len(all.String())};tsNodeReadableSyncState(state);return receiver})
	define("pipe",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{dest:=tsLegacyPipe(loop,receiver,args);if state:=tsNodeReadableStateFor(receiver);state!=nil{state.pipes=append(state.pipes,dest);tsNodeReadableSyncState(state)};return dest})
	define("unpipe",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeReadableStateFor(receiver);if state==nil{return receiver};dest:=tsArg(args,0);if tsIsUndefined(dest){state.pipes=nil}else{for i:=0;i<len(state.pipes);i++{if tsStrictEqual(state.pipes[i],dest){state.pipes=append(state.pipes[:i],state.pipes[i+1:]...);break}}};tsNodeReadableSyncState(state);return receiver})
	// Event methods delegate to the legacy Stream prototype without recursing.
	for _,entry:=range []struct{name string;length int}{{"on",2},{"addListener",2},{"removeListener",2},{"removeAllListeners",1},{"off",2}}{name:=entry.name;define(name,entry.length,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{legacy:=tsInternalStreamLegacyModule(loop);streamClass:=tsGet(loop,legacy,tsNodeReadableKey("Stream"));proto:=tsGet(loop,streamClass,tsNodeReadableKey("prototype"));method:=tsGet(loop,proto,tsNodeReadableKey(name));if name=="addListener"{method=tsGet(loop,proto,tsNodeReadableKey("addListener"))};result:=tsCallReceiver(loop,method,receiver,args...);if (name=="on"||name=="addListener")&&tsText(tsArg(args,0))=="data"{if state:=tsNodeReadableStateFor(receiver);state!=nil{state.flowing=true;state.flowingSet=true;tsNodeReadableSyncState(state);tsNodeReadableResume(loop,state)}};return result})}
	define("resume",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{if state:=tsNodeReadableStateFor(receiver);state!=nil{state.paused=false;state.flowing=true;state.flowingSet=true;tsNodeReadableSyncState(state);tsNodeReadableResume(loop,state)};return receiver})
	define("pause",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{if state:=tsNodeReadableStateFor(receiver);state!=nil{state.paused=true;state.flowing=false;state.flowingSet=true;tsNodeReadableSyncState(state);tsNodeReadableEmit(loop,state,"pause")};return receiver})
	define("wrap",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeReadableStateFor(receiver);source:=tsArg(args,0);if state==nil{return receiver};onData:=tsFunctionValue(tsFunc(func(loop *tsLoop,values ...tsValue)tsValue{chunk:=tsArg(values,0);if !tsTruthy(tsNodeReadablePush(loop,state,chunk,tsU,false)){pause:=tsNodeReadableGet(loop,source,"pause");if pause.kind==tsFunctionKind{tsCallReceiver(loop,pause,source)}};return tsU}));onEnd:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{tsNodeReadablePush(loop,state,tsNull,tsU,false);return tsU}));tsCallReceiver(loop,tsNodeReadableGet(loop,source,"on"),source,tsNodeReadableKey("data"),onData);tsCallReceiver(loop,tsNodeReadableGet(loop,source,"on"),source,tsNodeReadableKey("end"),onEnd);state.readHook=tsNamedNativeMethod("_read",1,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{resume:=tsNodeReadableGet(loop,source,"resume");if resume.kind==tsFunctionKind{tsCallReceiver(loop,resume,source)};return tsU});return receiver})
	asyncKey:=tsPropertyKey(tsWellKnownSymbol("asyncIterator"));p.set(asyncKey,tsNamedNativeMethod("[Symbol.asyncIterator]",0,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeReadableAsyncIterator(loop,receiver,tsU)}));p.descriptors[asyncKey]=&tsDescriptor{writable:true,configurable:true}
	define("iterator",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{options:=tsArg(args,0);if !tsIsUndefined(options){tsInternalValidatorObject(options,tsNodeReadableKey("options"),tsU)};iterator:=tsNodeReadableAsyncIterator(loop,receiver,options);tsSet(loop,iterator,tsNodeReadableKey("stream"),receiver);return iterator})
	asyncDisposeKey:=tsPropertyKey(tsWellKnownSymbol("asyncDispose"));p.set(asyncDisposeKey,tsNamedNativeMethod("[Symbol.asyncDispose]",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeReadableStateFor(receiver);if state==nil{return tsPromiseValue(loop.resolved(tsNull,false))};if !state.destroyed{tsNodeReadableDestroy(loop,state,tsFSException("AbortError","ABORT_ERR","The operation was aborted"),tsU)};promise:=loop.promise();callback:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{err:=tsArg(args,0);if tsTruthy(err){promise.settle(tsResult{err,true})}else{promise.settle(tsResult{tsNull,false})};return tsU}));options:=tsNewObject();options.set("writable",tsBooleanValue(false));tsCall(loop,tsInternalStreamEndOfStreamModule(loop),receiver,tsObjectValue(options),callback);return tsPromiseValue(promise)}));p.descriptors[asyncDisposeKey]=&tsDescriptor{writable:true,configurable:true}
	types:=tsInternalStreamIterTypesModule(loop);toAsync:=tsGet(loop,types,tsNodeReadableKey("toAsyncStreamable"));toAsyncKey:=tsPropertyKey(toAsync);p.set(toAsyncKey,tsNamedNativeMethod("[Symbol.for(Stream.toAsyncStreamable)]",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{module:=tsNodeModuleRequire(loop,[]tsValue{tsStringReference(tsStringUTF8("internal/streams/iter/classic"))});create:=tsGet(loop,module,tsNodeReadableKey("createBatchedAsyncIterator"));normalize:=tsNull;if state:=tsNodeReadableStateFor(receiver);state!=nil&&(state.objectMode||state.encoding!=""){normalize=tsGet(loop,module,tsNodeReadableKey("normalizeBatch"))};iter:=tsCall(loop,create,receiver,normalize);tsSet(loop,iter,tsGet(loop,types,tsNodeReadableKey("kValidatedSource")),tsBooleanValue(true));tsSet(loop,iter,tsNodeReadableKey("stream"),receiver);return iter}));p.descriptors[toAsyncKey]=&tsDescriptor{writable:true,configurable:true}
	define("compose",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{options:=tsArg(args,1);if !tsNullish(options){tsInternalValidatorObject(options,tsNodeReadableKey("options"),tsU);signal:=tsNodeReadableGet(loop,options,"signal");if !tsNullish(signal){tsInternalValidatorAbortSignal(signal,tsNodeReadableKey("options.signal"))}};module:=tsNodeModuleRequire(loop,[]tsValue{tsStringReference(tsStringUTF8("internal/streams/compose"))});compose:=tsGet(loop,module,tsNodeReadableKey("default"));result:=tsCall(loop,compose,receiver,tsArg(args,0));signal:=tsU;if !tsNullish(options){signal=tsNodeReadableGet(loop,options,"signal")};if !tsNullish(signal){tsInternalStreamAddAbortSignalNoValidate(loop,signal,result)};return result})
	for _,entry:=range []struct{name string;getter func(*tsLoop,*tsNodeReadableState)tsValue}{{"readable",func(_ *tsLoop,s *tsNodeReadableState)tsValue{return tsBooleanValue(s!=nil&&s.readable&&!s.destroyed&&!tsTruthy(s.errorValue)&&!s.endEmitted)}},{"readableDidRead",func(_ *tsLoop,s *tsNodeReadableState)tsValue{return tsBooleanValue(s!=nil&&s.dataEmitted)}},{"readableAborted",func(_ *tsLoop,s *tsNodeReadableState)tsValue{return tsBooleanValue(s!=nil&&s.readable&&(s.destroyed||tsTruthy(s.errorValue))&&!s.endEmitted)}},{"readableHighWaterMark",func(_ *tsLoop,s *tsNodeReadableState)tsValue{if s==nil{return tsNumberValue(0)};return tsNumberValue(float64(s.highWaterMark))}},{"readableBuffer",func(_ *tsLoop,s *tsNodeReadableState)tsValue{if s==nil{return tsU};return tsArrayValue(&tsArray{values:append([]tsValue{},s.buffer...)})}},{"readableFlowing",func(_ *tsLoop,s *tsNodeReadableState)tsValue{if s==nil||!s.flowingSet{return tsNull};return tsBooleanValue(s.flowing)}},{"readableLength",func(_ *tsLoop,s *tsNodeReadableState)tsValue{if s==nil{return tsNumberValue(0)};return tsNumberValue(float64(s.length))}},{"readableObjectMode",func(_ *tsLoop,s *tsNodeReadableState)tsValue{return tsBooleanValue(s!=nil&&s.objectMode)}},{"readableEncoding",func(_ *tsLoop,s *tsNodeReadableState)tsValue{if s==nil||s.encoding==""{return tsNull};return tsStringReference(tsStringUTF8(s.encoding))}},{"errored",func(_ *tsLoop,s *tsNodeReadableState)tsValue{if s==nil{return tsNull};return s.errorValue}},{"closed",func(_ *tsLoop,s *tsNodeReadableState)tsValue{return tsBooleanValue(s!=nil&&s.closed)}},{"destroyed",func(_ *tsLoop,s *tsNodeReadableState)tsValue{return tsBooleanValue(s!=nil&&s.destroyed)}},{"readableEnded",func(_ *tsLoop,s *tsNodeReadableState)tsValue{return tsBooleanValue(s!=nil&&s.endEmitted)}}}{name:=entry.name;getter:=entry.getter;p.descriptors[name]=&tsDescriptor{accessor:true,configurable:true,getter:tsNamedNativeMethod("get "+name,0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{return getter(loop,tsNodeReadableStateFor(receiver))})}}
	p.descriptors["readable"].setter=tsNamedNativeMethod("set readable",1,func(_ *tsLoop,receiver tsValue,args ...tsValue)tsValue{if state:=tsNodeReadableStateFor(receiver);state!=nil{state.readable=tsTruthy(tsArg(args,0));tsNodeReadableSyncState(state)};return tsU})
	p.descriptors["readableFlowing"].setter=tsNamedNativeMethod("set readableFlowing",1,func(_ *tsLoop,receiver tsValue,args ...tsValue)tsValue{if state:=tsNodeReadableStateFor(receiver);state!=nil{value:=tsArg(args,0);state.flowingSet=!tsNullish(value);state.flowing=tsTruthy(value);tsNodeReadableSyncState(state)};return tsU})
	p.descriptors["destroyed"].setter=tsNamedNativeMethod("set destroyed",1,func(_ *tsLoop,receiver tsValue,args ...tsValue)tsValue{if state:=tsNodeReadableStateFor(receiver);state!=nil{state.destroyed=tsTruthy(tsArg(args,0));tsNodeReadableSyncState(state)};return tsU})
	constructor:=tsNodeReadableSingleton.constructor;constructor.properties.set("ReadableState",tsNamedNativeMethod("ReadableState",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeReadableStateCompat(loop,tsArg(args,0))}));constructor.properties.set("from",tsNamedNativeMethod("from",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsInternalStreamFrom(loop,[]tsValue{tsFunctionValue(constructor),tsArg(args,0),tsArg(args,1)})}));constructor.properties.set("fromWeb",tsNamedNativeMethod("fromWeb",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{module:=tsNodeModuleRequire(loop,[]tsValue{tsStringReference(tsStringUTF8("internal/webstreams_adapters"))});return tsCall(loop,tsGet(loop,module,tsNodeReadableKey("newStreamReadableFromReadableStream")),tsArg(args,0),tsArg(args,1))}));constructor.properties.set("toWeb",tsNamedNativeMethod("toWeb",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{module:=tsNodeModuleRequire(loop,[]tsValue{tsStringReference(tsStringUTF8("internal/webstreams_adapters"))});return tsCall(loop,tsGet(loop,module,tsNodeReadableKey("newReadableStreamFromStreamReadable")),tsArg(args,0),tsArg(args,1))}));constructor.properties.set("_fromList",tsNamedNativeMethod("_fromList",2,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{state:=tsNodeReadableGet(loop,tsArg(args,1),"stream");return tsNodeReadableRead(loop,tsNodeReadableStateFor(state),tsArg(args,0))}));constructor.properties.set("wrap",tsNamedNativeMethod("wrap",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{stream:=tsConstruct(tsFunctionValue(constructor),tsArg(args,1));return tsCallReceiver(loop,tsGet(loop,stream,tsNodeReadableKey("wrap")),stream,tsArg(args,0))}))
}

func tsNodeReadableResume(loop *tsLoop,state *tsNodeReadableState) {
	if state==nil||state.destroyed{return};tsNodeReadableDefer(loop,func(){tsNodeReadableEmit(loop,state,"resume");for {if state.destroyed||!state.flowing{break};chunk:=tsNodeReadableRead(loop,state,tsU);if tsNullish(chunk){if !state.ended{tsNodeReadableScheduleRead(loop,state)};break};if state.length==0&&state.ended{tsNodeReadableEnd(loop,state)}}})
}
func tsNodeReadableAsyncIterator(loop *tsLoop,stream,options tsValue) tsValue {
	state:=tsNodeReadableStateFor(stream);if state==nil{return tsAsyncIteratorObject(tsAsyncIterate(tsU))};iterator:=&tsAsyncIterator{}
	iterator.next=func(loop *tsLoop)*tsPromise{
		if value:=tsNodeReadableRead(loop,state,tsU);!tsNullish(value){return loop.resolved(tsIteratorResult(value,false),false)}
		if tsTruthy(state.errorValue){promise:=loop.promise();promise.settle(tsResult{state.errorValue,true});return promise};if state.endEmitted||state.destroyed{return loop.resolved(tsIteratorResult(tsU,true),false)}
		promise:=loop.promise();listeners:=map[string]tsValue{};settled:=false
		cleanup:=func(){if settled{return};settled=true;for event,listener:=range listeners{tsCallReceiver(loop,tsNodeReadableGet(loop,stream,"removeListener"),stream,tsNodeReadableKey(event),listener)}}
		wake:=func(loop *tsLoop,event string,args []tsValue)tsValue{if settled{return tsU};if event=="error"{err:=tsArg(args,0);state.errorValue=err;cleanup();promise.settle(tsResult{err,true});return tsU};if event=="close"&&!state.endEmitted&&!tsTruthy(state.errorValue){err:=tsFSException("Error","ERR_STREAM_PREMATURE_CLOSE","Premature close");cleanup();promise.settle(tsResult{err,true});return tsU};if value:=tsNodeReadableRead(loop,state,tsU);!tsNullish(value){cleanup();promise.settle(tsResult{tsIteratorResult(value,false),false});return tsU};if tsTruthy(state.errorValue){cleanup();promise.settle(tsResult{state.errorValue,true});return tsU};if state.endEmitted||state.destroyed{cleanup();promise.settle(tsResult{tsIteratorResult(tsU,true),false})};return tsU}
		for _,event:=range []string{"readable","end","error","close"}{name:=event;listener:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return wake(loop,name,args)}));listeners[event]=listener;tsCallReceiver(loop,tsNodeReadableGet(loop,stream,"on"),stream,tsNodeReadableKey(event),listener)}
		return promise
	}
	iterator.close=func(loop *tsLoop)*tsPromise{if !tsStrictEqual(tsNodeReadableGet(loop,options,"destroyOnReturn"),tsBooleanValue(false)){tsNodeReadableDestroy(loop,state,tsU,tsU)};return loop.resolved(tsIteratorResult(tsU,true),false)}
	iterator.resume=func(loop *tsLoop,action string,value tsValue)*tsPromise{if action=="return"{return iterator.close(loop)};return iterator.next(loop)}
	result:=tsAsyncIteratorObject(iterator);tsSet(loop,result,tsNodeReadableKey("stream"),stream);return result
}

func tsInternalStreamReadableModule(loop *tsLoop) tsValue {
	tsNodeReadableInit(loop);module:=tsNewObject();module.set("default",tsFunctionValue(tsNodeReadableSingleton.constructor));module.set("Readable",tsFunctionValue(tsNodeReadableSingleton.constructor));return tsObjectValue(module)
}
`
