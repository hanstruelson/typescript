package goemit

// internalStreamTransformRuntime builds Transform on the shared Duplex
// runtime and preserves the one-chunk transform/backpressure contract.
const internalStreamTransformRuntime = `
type tsNodeTransformState struct { callback, flush, transform, final tsValue }
var tsNodeTransformSingleton struct { sync.Once; constructor *tsFunction; prototype *tsObject; final tsValue; states map[uintptr]*tsNodeTransformState }
func tsNodeTransformKey(name string) tsValue { return tsStringReference(tsStringUTF8(name)) }
func tsNodeTransformGet(loop *tsLoop,value tsValue,name string) tsValue { return tsGet(loop,value,tsNodeTransformKey(name)) }
func tsNodeTransformSet(loop *tsLoop,value tsValue,name string,item tsValue) { tsSet(loop,value,tsNodeTransformKey(name),item) }
func tsNodeTransformStateFor(value tsValue)*tsNodeTransformState{if value.kind!=tsObjectKind{return nil};return tsNodeTransformSingleton.states[uintptr(value.ref)]}
func tsNodeTransformDefer(loop *tsLoop,callback tsValue,args ...tsValue){values:=append([]tsValue{},args...);tsNodeWritableDefer(loop,func(){if callback.kind==tsFunctionKind{tsCall(loop,callback,values...)}})}
func tsNodeTransformFinal(loop *tsLoop,stream,callback tsValue) tsValue {
	flush:=tsNodeTransformGet(loop,stream,"_flush")
	finish:=func(err,data tsValue){if tsTruthy(err){if callback.kind==tsFunctionKind{tsCall(loop,callback,err)}else{tsNodeReadableDestroy(loop,tsNodeReadableStateFor(stream),err,tsU)};return};if !tsNullish(data){tsNodeClassicMethod(loop,stream,"push",data)};tsNodeClassicMethod(loop,stream,"push",tsNull);if callback.kind==tsFunctionKind{tsCall(loop,callback,tsU)}}
	if flush.kind==tsFunctionKind&&!tsTruthy(tsNodeTransformGet(loop,stream,"destroyed")){called:=false;done:=tsFunctionValue(tsFunc(func(_ *tsLoop,args ...tsValue)tsValue{if called{return tsU};called=true;finish(tsArg(args,0),tsArg(args,1));return tsU}));return tsCallReceiver(loop,flush,stream,done)}
	finish(tsU,tsU);return tsU
}
func tsNodeTransformPrefinish(loop *tsLoop,stream tsValue) tsValue {state:=tsNodeTransformStateFor(stream);if state!=nil&&!tsStrictEqual(tsNodeTransformGet(loop,stream,"_final"),state.final){return tsNodeTransformFinal(loop,stream,tsU)};return tsU}
func tsNodeTransformCreate(loop *tsLoop,receiver,options tsValue) tsValue {
	tsNodeTransformInit(loop);if tsNullish(options){options=tsObjectValue(tsNewObject())}
	stateObject:=tsObjectValue(tsNewObject());tsSet(loop,stateObject,tsNodeTransformKey("objectMode"),tsBooleanValue(tsTruthy(tsNodeTransformGet(loop,options,"objectMode"))))
	readableHWM:=tsInternalStreamStateGetHighWaterMark(loop,stateObject,options,tsNodeTransformKey("readableHighWaterMark"),true)
	if tsNumber(readableHWM)==0 {copy:=tsNewObject();tsObjectSpread(copy,options);copy.set("highWaterMark",tsNull);writableHWM:=tsNodeTransformGet(loop,options,"writableHighWaterMark");if !tsTruthy(writableHWM){writableHWM=tsNumberValue(0)};copy.set("writableHighWaterMark",writableHWM);options=tsObjectValue(copy)}
	stream:=tsNodeDuplexCreate(loop,receiver,options);object:=(*tsObject)(stream.ref);object.prototype=tsObjectValue(tsNodeTransformSingleton.prototype)
	state:=&tsNodeTransformState{callback:tsNull,transform:tsNodeTransformGet(loop,options,"transform"),flush:tsNodeTransformGet(loop,options,"flush"),final:tsNodeTransformSingleton.final};tsNodeTransformSingleton.states[uintptr(stream.ref)]=state
	if state.transform.kind==tsFunctionKind{tsNodeTransformSet(loop,stream,"_transform",state.transform)};if state.flush.kind==tsFunctionKind{tsNodeTransformSet(loop,stream,"_flush",state.flush)}
	if readable:=tsNodeReadableStateFor(stream);readable!=nil{readable.sync=false;tsNodeReadableSyncState(readable)}
	prefinish:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{return tsNodeTransformPrefinish(loop,stream)}));tsNodeClassicMethod(loop,stream,"on",tsNodeTransformKey("prefinish"),prefinish)
	return stream
}
func tsNodeTransformInit(loop *tsLoop){tsNodeTransformSingleton.Do(func(){duplex:=tsInternalStreamDuplexModule(loop);duplexCtor:=tsGet(loop,duplex,tsNodeTransformKey("Duplex"));baseProto:=tsGet(loop,duplexCtor,tsNodeTransformKey("prototype"));prototype:=tsNewObject();prototype.prototype=baseProto
	constructor:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeTransformCreate(loop,tsU,tsArg(args,0))});constructor.name="Transform";constructor.length=1;constructor.constructible=true;constructor.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeTransformCreate(loop,tsU,tsArg(args,0))};constructor.receiverCall=func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeTransformCreate(loop,receiver,tsArg(args,0))};constructor.properties=tsNewObject();constructor.properties.set("prototype",tsObjectValue(prototype));constructor.properties.descriptors=map[string]*tsDescriptor{"prototype":{writable:false}};prototype.set("constructor",tsFunctionValue(constructor));prototype.descriptors=map[string]*tsDescriptor{"constructor":{writable:true,configurable:true}}
	finalMethod:=tsNamedNativeMethod("_final",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeTransformFinal(loop,receiver,tsArg(args,0))});prototype.set("_final",finalMethod);tsNodeTransformSingleton.final=finalMethod;prototype.set("_transform",tsNamedNativeMethod("_transform",3,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{panic(tsThrown{tsFSException("Error","ERR_METHOD_NOT_IMPLEMENTED","_transform() is not implemented")})}));prototype.set("_write",tsNamedNativeMethod("_write",3,tsNodeTransformWrite));prototype.set("_read",tsNamedNativeMethod("_read",0,tsNodeTransformRead))
	tsNodeTransformSingleton.constructor=constructor;tsNodeTransformSingleton.prototype=prototype;tsNodeTransformSingleton.states=map[uintptr]*tsNodeTransformState{}
	})}
func tsNodeTransformWrite(loop *tsLoop,stream tsValue,args ...tsValue)tsValue{
	chunk,encoding,callback:=tsArg(args,0),tsArg(args,1),tsArg(args,2);readable:=tsNodeReadableStateFor(stream);writable:=tsNodeWritableStateFor(stream);if readable==nil||writable==nil{if callback.kind==tsFunctionKind{tsCall(loop,callback,tsFSException("TypeError","ERR_INVALID_STATE","Invalid Transform stream state"))};return tsU};length:=readable.length;transform:=tsNodeTransformGet(loop,stream,"_transform");done:=tsFunctionValue(tsFunc(func(loop *tsLoop,doneArgs ...tsValue)tsValue{err,value:=tsArg(doneArgs,0),tsArg(doneArgs,1);if tsTruthy(err){if callback.kind==tsFunctionKind{tsCall(loop,callback,err)};return tsU};if !tsNullish(value){tsNodeReadablePush(loop,readable,value,tsU,false)};if readable.ended{tsNodeWritableDefer(loop,func(){if callback.kind==tsFunctionKind{tsCall(loop,callback,tsU)}})}else if writable.ending||length==readable.length||readable.length<readable.highWaterMark{if callback.kind==tsFunctionKind{tsCall(loop,callback,tsU)}}else{state:=tsNodeTransformStateFor(stream);if state!=nil{state.callback=callback}};return tsU}))
	if transform.kind!=tsFunctionKind{panic(tsThrown{tsFSException("Error","ERR_METHOD_NOT_IMPLEMENTED","_transform() is not implemented")})};tsCallReceiver(loop,transform,stream,chunk,encoding,done);return tsU
}
func tsNodeTransformRead(loop *tsLoop,stream tsValue,_ ...tsValue)tsValue{state:=tsNodeTransformStateFor(stream);if state!=nil&&!tsNullish(state.callback){callback:=state.callback;state.callback=tsNull;tsCall(loop,callback,tsU)};return tsU}
func tsInternalStreamTransformModule(loop *tsLoop)tsValue{tsNodeTransformInit(loop);module:=tsNewObject();module.set("default",tsFunctionValue(tsNodeTransformSingleton.constructor));module.set("Transform",tsFunctionValue(tsNodeTransformSingleton.constructor));return tsObjectValue(module)}
`
