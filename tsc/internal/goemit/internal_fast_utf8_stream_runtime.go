package goemit

// internalFastUTF8StreamRuntime is a Go filesystem-backed implementation of
// Bun's internal fast UTF-8/file writer surface. It keeps the public buffering,
// flush, reopen, and lifecycle methods while using the shared FS descriptor API.
const internalFastUTF8StreamRuntime = `
type tsInternalFastUTF8StreamState struct {
	loop *tsLoop
	object *tsObject
	events *tsFSEvents
	fd tsValue
	file string
	mode tsValue
	buffer []byte
	current []byte
	minLength,maxLength,maxWrite,highWaterMark int
	periodicFlush int
	sync,fsync,append,mkdir,bufferMode bool
	opening,writing,ending,destroyed,closed bool
	flushCallbacks []tsValue
	timerCancel chan struct{}
}

func tsInternalFastUTF8StreamSet(state *tsInternalFastUTF8StreamState,name string,value tsValue) { state.object.set(name,value) }
func tsInternalFastUTF8StreamUpdate(state *tsInternalFastUTF8StreamState) {
	state.object.set("fd",state.fd);state.object.set("file",tsStringReference(tsStringUTF8(state.file)))
	state.object.set("writing",tsBooleanValue(state.writing));state.object.set("destroyed",tsBooleanValue(state.destroyed))
}
func tsInternalFastUTF8StreamWriteSync(state *tsInternalFastUTF8StreamState,data []byte) int {
	if len(data)==0{return 0}
	written:=int(tsNumber(tsFSInvoke(state.loop,"write","sync",[]tsValue{state.fd,tsNodeBuffer(data),tsNumberValue(0),tsNumberValue(float64(len(data))),tsNull})))
	if state.fsync{syncResult:=tsFSExecute(state.loop,"fsync",[]tsValue{state.fd},false);if syncResult.rejected{panic(tsThrown{syncResult.value})}}
	return written
}
func tsInternalFastUTF8StreamFlushCallbacks(state *tsInternalFastUTF8StreamState,err tsValue) {
	callbacks:=state.flushCallbacks;state.flushCallbacks=nil
	for _,callback:=range callbacks{if callback.kind==tsFunctionKind{state.loop.ready=append(state.loop.ready,func(){tsCall(state.loop,callback,err)})}}
}
func tsInternalFastUTF8StreamClose(state *tsInternalFastUTF8StreamState,emitFinish bool) {
	if state.closed{return};if state.timerCancel!=nil{close(state.timerCancel);state.timerCancel=nil}
	finish:=func(result tsResult){state.closed=true;state.destroyed=true;state.fd=tsNumberValue(-1);tsInternalFastUTF8StreamUpdate(state);if result.rejected{state.events.emitEvent("error",result.value);return};if emitFinish{state.events.emitEvent("finish")};state.events.emitEvent("close")}
	if tsIsUndefined(state.fd)||tsIsNull(state.fd)||tsNumber(state.fd)<0{finish(tsResult{tsU,false});return}
	if tsIsNumeric(state.fd)&&(tsNumber(state.fd)==1||tsNumber(state.fd)==2){finish(tsResult{tsU,false});return}
	fd:=state.fd;state.fd=tsNumberValue(-1);tsInternalFastUTF8StreamUpdate(state)
	state.loop.submit(func()tsResult{return tsFSExecute(state.loop,"close",[]tsValue{fd},true)},finish)
}
func tsInternalFastUTF8StreamFinish(state *tsInternalFastUTF8StreamState) {
	if state.destroyed||state.writing||state.opening{return}
	if len(state.buffer)>0{tsInternalFastUTF8StreamStartWrite(state);return}
	if state.ending{tsInternalFastUTF8StreamClose(state,true);return}
	tsInternalFastUTF8StreamFlushCallbacks(state,tsNull);state.events.emitEvent("drain")
}
func tsInternalFastUTF8StreamStartWrite(state *tsInternalFastUTF8StreamState) {
	if state.destroyed||state.opening||state.writing||tsNullish(state.fd){return}
	if len(state.buffer)==0{tsInternalFastUTF8StreamFinish(state);return}
	size:=len(state.buffer);if state.maxWrite>0&&size>state.maxWrite{size=state.maxWrite}
	state.current=append([]byte{},state.buffer[:size]...);state.writing=true;tsInternalFastUTF8StreamUpdate(state)
	if state.sync {
		result:=tsFSDescriptorIO("write",[]tsValue{state.fd,tsNodeBuffer(state.current),tsNumberValue(0),tsNumberValue(float64(len(state.current))),tsNull})
		if result.rejected{state.writing=false;tsInternalFastUTF8StreamUpdate(state);state.events.emitEvent("error",result.value);return}
		n:=int(tsNumber(tsFSProperty(result.value,"bytesWritten")));if n<0{n=0};if n>len(state.buffer){n=len(state.buffer)};state.buffer=append([]byte{},state.buffer[n:]...);state.events.emitEvent("write",tsNumberValue(float64(n)))
		if state.fsync{syncResult:=tsFSExecute(state.loop,"fsync",[]tsValue{state.fd},false);if syncResult.rejected{state.events.emitEvent("error",syncResult.value);return}}
		state.writing=false;state.current=nil;tsInternalFastUTF8StreamUpdate(state);tsInternalFastUTF8StreamStartWrite(state);return
	}
	fd,current:=state.fd,append([]byte{},state.current...)
	state.loop.submit(func()tsResult{return tsFSDescriptorIO("write",[]tsValue{fd,tsNodeBuffer(current),tsNumberValue(0),tsNumberValue(float64(len(current))),tsNull})},func(result tsResult){
		state.writing=false;state.current=nil
		if result.rejected{tsInternalFastUTF8StreamUpdate(state);state.events.emitEvent("error",result.value);tsInternalFastUTF8StreamFlushCallbacks(state,result.value);return}
		n:=int(tsNumber(tsFSProperty(result.value,"bytesWritten")));if n<0{n=0};if n>len(state.buffer){n=len(state.buffer)};state.buffer=append([]byte{},state.buffer[n:]...);state.events.emitEvent("write",tsNumberValue(float64(n)));tsInternalFastUTF8StreamUpdate(state)
		if state.fsync{syncResult:=tsFSExecute(state.loop,"fsync",[]tsValue{state.fd},true);if syncResult.rejected{state.events.emitEvent("error",syncResult.value);return};state.loop.submit(func()tsResult{return syncResult},func(_ tsResult){tsInternalFastUTF8StreamStartWrite(state)})}else{tsInternalFastUTF8StreamStartWrite(state)}
	})
}
func tsInternalFastUTF8StreamOpen(loop *tsLoop,state *tsInternalFastUTF8StreamState,path string) {
	state.file=path;state.opening=true;tsInternalFastUTF8StreamUpdate(state)
	if state.mkdir{if err:=os.MkdirAll(filepath.Dir(path),0777);err!=nil{state.opening=false;state.events.emitEvent("error",tsFSError(err,"mkdir",filepath.Dir(path),""));return}}
	flags:="a";if !state.append{flags="w"};pathValue:=tsStringReference(tsStringUTF8(path));flagValue:=tsStringReference(tsStringUTF8(flags))
	if state.sync{
		result:=tsFSExecute(loop,"open",[]tsValue{pathValue,flagValue,state.mode},false);state.opening=false;if result.rejected{panic(tsThrown{result.value})};state.fd=result.value;tsInternalFastUTF8StreamUpdate(state);loop.ready=append(loop.ready,func(){state.events.emitEvent("ready")});return
	}
	loop.submit(func()tsResult{return tsFSExecute(loop,"open",[]tsValue{pathValue,flagValue,state.mode},true)},func(result tsResult){state.opening=false;if result.rejected{tsInternalFastUTF8StreamUpdate(state);state.events.emitEvent("error",result.value);return};state.fd=result.value;tsInternalFastUTF8StreamUpdate(state);state.events.emitEvent("ready");if state.ending||len(state.buffer)>state.minLength{tsInternalFastUTF8StreamStartWrite(state)}})
}
func tsInternalFastUTF8StreamConstruct(loop *tsLoop,args []tsValue) tsValue {
	options:=tsArg(args,0);if tsIsUndefined(options){options=tsObjectValue(tsNewObject())};if options.kind!=tsObjectKind{tsFSInvalid("options","an object")}
	fd:=tsFSProperty(options,"fd");if tsIsNull(fd)||tsIsUndefined(fd){fd=tsFSProperty(options,"dest")}
	if !tsIsNumeric(fd)&&fd.kind!=tsStringKind{tsFSInvalid("fd","a number or string")}
	contentMode:=tsFSOption(options,"contentMode",tsStringReference(tsStringUTF8("utf8")));if contentMode.kind!=tsStringKind||(tsText(contentMode)!="utf8"&&tsText(contentMode)!="buffer"){panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_VALUE","options.contentMode must be 'utf8' or 'buffer'")})}
	state:=&tsInternalFastUTF8StreamState{loop:loop,fd:tsU,events:tsFSEventObject(loop),minLength:int(tsFSInteger(tsFSProperty(options,"minLength"),"options.minLength",0,0)),maxLength:int(tsFSInteger(tsFSProperty(options,"maxLength"),"options.maxLength",0,0)),maxWrite:int(tsFSInteger(tsFSProperty(options,"maxWrite"),"options.maxWrite",16384,1)),highWaterMark:16387,periodicFlush:int(tsFSInteger(tsFSProperty(options,"periodicFlush"),"options.periodicFlush",0,0)),sync:tsTruthy(tsFSProperty(options,"sync")),fsync:tsTruthy(tsFSProperty(options,"fsync")),append:tsTruthy(tsFSOption(options,"append",tsBooleanValue(true))),mkdir:tsTruthy(tsFSProperty(options,"mkdir")),bufferMode:tsText(contentMode)=="buffer",mode:tsFSProperty(options,"mode")}
	for _,key:=range []string{"sync","fsync","append","mkdir"}{value:=tsFSProperty(options,key);if !tsIsUndefined(value)&&value.kind!=tsBooleanKind{tsFSInvalid("options."+key,"a boolean")}}
	if state.minLength>state.highWaterMark{state.highWaterMark=state.minLength}
	if state.minLength>=state.maxWrite{panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","minLength should be smaller than maxWrite")})}
	if !tsIsUndefined(tsFSProperty(options,"retryEAGAIN"))&&tsFSProperty(options,"retryEAGAIN").kind!=tsFunctionKind{tsFSInvalid("options.retryEAGAIN","a function")}
	state.object=state.events.object;state.object.set("fd",tsNumberValue(-1));state.object.set("file",tsNull);state.object.set("mode",state.mode);state.object.set("minLength",tsNumberValue(float64(state.minLength)));state.object.set("maxLength",tsNumberValue(float64(state.maxLength)));state.object.set("writing",tsBooleanValue(false));state.object.set("sync",tsBooleanValue(state.sync));state.object.set("fsync",tsBooleanValue(state.fsync));state.object.set("append",tsBooleanValue(state.append));state.object.set("periodicFlush",tsNumberValue(float64(state.periodicFlush)));state.object.set("contentMode",contentMode);state.object.set("mkdir",tsBooleanValue(state.mkdir));state.object.set("destroyed",tsBooleanValue(false))
	state.object.set("write",tsNamedNativeMethod("write",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{
		if state.destroyed{panic(tsThrown{tsFSException("Error","ERR_INVALID_STATE","Utf8Stream is destroyed")})}
		data:=tsArg(args,0);var bytes []byte
		if state.bufferMode{if data.kind!=tsTypedArrayKind||!(*tsTypedArray)(data.ref).nodeBuffer{tsFSInvalid("data","a Buffer")};bytes=append([]byte{},tsFSBytes(data)...)}else{if data.kind!=tsStringKind{tsFSInvalid("data","a string")};bytes=[]byte(tsText(data))}
		if state.maxLength>0&&len(state.buffer)+len(bytes)>state.maxLength{state.events.emitEvent("drop",data);return tsBooleanValue(len(state.buffer)<state.highWaterMark)}
		state.buffer=append(state.buffer,bytes...);if !state.writing&&len(state.buffer)>=state.minLength{tsInternalFastUTF8StreamStartWrite(state)};return tsBooleanValue(len(state.buffer)<state.highWaterMark)
	}))
	state.object.set("flush",tsNamedNativeMethod("flush",0,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{callback:=tsArg(args,0);if tsIsUndefined(callback)||tsIsNull(callback){callback=tsFunctionValue(tsFunc(func(_ *tsLoop,_ ...tsValue)tsValue{return tsU}))};if callback.kind!=tsFunctionKind{tsFSInvalid("cb","a function")};state.flushCallbacks=append(state.flushCallbacks,callback);if len(state.buffer)>0{tsInternalFastUTF8StreamStartWrite(state)}else{tsInternalFastUTF8StreamFlushCallbacks(state,tsNull)};return tsU}))
	state.object.set("flushSync",tsNamedNativeMethod("flushSync",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{if state.destroyed{panic(tsThrown{tsFSException("Error","ERR_INVALID_STATE","Utf8Stream is destroyed")})};if tsNullish(state.fd){panic(tsThrown{tsFSException("Error","ERR_INVALID_STATE","Invalid file descriptor")})};for len(state.buffer)>0{n:=tsInternalFastUTF8StreamWriteSync(state,state.buffer);if n<=0{break};state.buffer=append([]byte{},state.buffer[n:]...)};return tsU}))
	state.object.set("reopen",tsNamedNativeMethod("reopen",0,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{if state.destroyed{panic(tsThrown{tsFSException("Error","ERR_INVALID_STATE","Utf8Stream is destroyed")})};if state.file==""{panic(tsThrown{tsFSException("Error","ERR_OPERATION_FAILED","Unable to reopen a file descriptor; pass a file path")})};path:=state.file;if !tsIsUndefined(tsArg(args,0)){path=tsText(tsArg(args,0))};state.file=path;state.opening=false;state.fd=tsU;tsInternalFastUTF8StreamOpen(loop,state,path);return tsU}))
	state.object.set("end",tsNamedNativeMethod("end",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{if state.destroyed{panic(tsThrown{tsFSException("Error","ERR_INVALID_STATE","Utf8Stream is destroyed")})};state.ending=true;if !state.writing&&!state.opening{tsInternalFastUTF8StreamFinish(state)};return tsU}))
	state.object.set("destroy",tsNamedNativeMethod("destroy",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{if !state.destroyed{state.buffer=nil;state.destroyed=true;tsInternalFastUTF8StreamClose(state,false)};return tsU}))
	state.object.set(tsPropertyKey(tsWellKnownSymbol("dispose")),state.object.get("destroy"))
	if tsIsNumeric(fd){state.fd=fd;state.object.set("fd",fd);state.object.set("file",tsNull);loop.ready=append(loop.ready,func(){state.events.emitEvent("ready")})}else{tsInternalFastUTF8StreamOpen(loop,state,tsText(fd))}
	if state.periodicFlush>0&&!state.sync{state.timerCancel=make(chan struct{});cancel:=state.timerCancel;task:=loop.submit(func()tsResult{ticker:=time.NewTicker(time.Duration(state.periodicFlush)*time.Millisecond);defer ticker.Stop();for{select{case <-cancel:return tsResult{tsU,false};case <-ticker.C:loop.post(func(){if !state.destroyed&&len(state.buffer)>0{tsInternalFastUTF8StreamStartWrite(state)}})}}},func(tsResult){});task.unref=true;loop.resourceClosers=append(loop.resourceClosers,func(){select{case <-cancel:default:close(cancel)}})}
	return tsObjectValue(state.object)
}
var tsInternalFastUTF8StreamSingleton struct { sync.Once; value tsValue }
func tsInternalFastUTF8StreamModule(loop *tsLoop) tsValue {
	tsInternalFastUTF8StreamSingleton.Do(func(){
		constructor:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{tsPropertyFailure("Utf8Stream must be constructed with new");return tsU});constructor.name="Utf8Stream";constructor.length=1;constructor.constructible=true
		constructor.properties=tsNewObject();prototype:=tsNewObject();events:=tsNodeEventsModule(loop);eventEmitter:=tsFSProperty(events,"EventEmitter");if !tsNullish(eventEmitter){prototype.prototype=tsFunctionGet(loop,eventEmitter,"prototype")};constructor.properties.set("prototype",tsObjectValue(prototype))
		constructor.construct=func(loop *tsLoop,args ...tsValue)tsValue{value:=tsInternalFastUTF8StreamConstruct(loop,args);if value.kind==tsObjectKind{(*tsObject)(value.ref).prototype=tsObjectValue(prototype)};return value}
		tsInternalFastUTF8StreamSingleton.value=tsFunctionValue(constructor)
	})
	return tsInternalFastUTF8StreamSingleton.value
}
`
