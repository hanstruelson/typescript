package goemit

// internalStreamDuplexRuntime combines the mapped Readable and Writable
// adapters using Node's readable prototype plus copied writable members.
const internalStreamDuplexRuntime = `
var tsNodeDuplexSingleton struct { sync.Once; constructor *tsFunction; prototype *tsObject }
func tsNodeDuplexKey(name string) tsValue { return tsStringReference(tsStringUTF8(name)) }
func tsNodeDuplexInit(loop *tsLoop) {
	tsNodeDuplexSingleton.Do(func() {
		readable:=tsValue(tsNodeReadableModuleValue(loop))
		writable:=tsInternalStreamWritableModule(loop)
		readableCtor:=tsGet(loop,readable,tsNodeDuplexKey("Readable"))
		writableCtor:=tsGet(loop,writable,tsNodeDuplexKey("Writable"))
		readableProto:=tsGet(loop,readableCtor,tsNodeDuplexKey("prototype"))
		writableProto:=tsGet(loop,writableCtor,tsNodeDuplexKey("prototype"))
		prototype:=tsNewObject();prototype.prototype=readableProto
		if writableProto.kind==tsObjectKind {
			from:=(*tsObject)(writableProto.ref)
			for _,name:=range from.order {
				if name=="constructor" { continue }
				prototype.set(name,from.values[name])
				if descriptor:=from.descriptors[name];descriptor!=nil { copy:=*descriptor;if prototype.descriptors==nil{prototype.descriptors=map[string]*tsDescriptor{}};prototype.descriptors[name]=&copy }
			}
		}
		constructor:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeDuplexCreate(loop,tsU,tsArg(args,0))})
		constructor.name="Duplex";constructor.length=1;constructor.constructible=true
		constructor.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeDuplexCreate(loop,tsU,tsArg(args,0))}
		constructor.receiverCall=func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeDuplexCreate(loop,receiver,tsArg(args,0))}
		constructor.properties=tsNewObject();constructor.properties.set("prototype",tsObjectValue(prototype));constructor.properties.descriptors=map[string]*tsDescriptor{"prototype":{writable:false}}
		constructor.properties.set("fromWeb",tsNamedNativeMethod("fromWeb",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{module:=tsNodeModuleRequire(loop,[]tsValue{tsStringReference(tsStringUTF8("internal/webstreams_adapters"))});return tsCall(loop,tsGet(loop,module,tsNodeDuplexKey("newStreamDuplexFromReadableWritablePair")),tsArg(args,0),tsArg(args,1))}))
		constructor.properties.set("toWeb",tsNamedNativeMethod("toWeb",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{module:=tsNodeModuleRequire(loop,[]tsValue{tsStringReference(tsStringUTF8("internal/webstreams_adapters"))});return tsCall(loop,tsGet(loop,module,tsNodeDuplexKey("newReadableWritablePairFromDuplex")),tsArg(args,0),tsArg(args,1))}))
		constructor.properties.set("from",tsNamedNativeMethod("from",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{module:=tsNodeModuleRequire(loop,[]tsValue{tsStringReference(tsStringUTF8("internal/streams/duplexify"))});return tsCall(loop,tsGet(loop,module,tsNodeDuplexKey("default")),tsArg(args,0),tsNodeDuplexKey("body"))}))
		prototype.set("constructor",tsFunctionValue(constructor));if prototype.descriptors==nil{prototype.descriptors=map[string]*tsDescriptor{}};prototype.descriptors["constructor"]=&tsDescriptor{writable:true,configurable:true}
		tsNodeDuplexSingleton.constructor=constructor;tsNodeDuplexSingleton.prototype=prototype
		tsNodeDuplexInstall(loop)
	})
}
func tsNodeDuplexOptions(loop *tsLoop,options tsValue) tsValue {
	if options.kind!=tsObjectKind { return options }
	source:=(*tsObject)(options.ref);copy:=tsNewObject();copy.prototype=source.prototype
	for _,name:=range source.order { if name=="construct"||name=="signal" { continue };copy.set(name,source.values[name]);if descriptor:=source.descriptors[name];descriptor!=nil{item:=*descriptor;if copy.descriptors==nil{copy.descriptors=map[string]*tsDescriptor{}};copy.descriptors[name]=&item} }
	return tsObjectValue(copy)
}
func tsNodeDuplexCreate(loop *tsLoop,receiver,options tsValue) tsValue {
	tsNodeDuplexInit(loop);if tsNullish(options){options=tsObjectValue(tsNewObject())}
	constructionOptions:=tsNodeDuplexOptions(loop,options)
	stream:=tsNodeReadableCreate(loop,receiver,constructionOptions)
	tsNodeWritableCreate(loop,stream,constructionOptions)
	object:=(*tsObject)(stream.ref);object.prototype=tsObjectValue(tsNodeDuplexSingleton.prototype)
	readable:=tsNodeReadableStateFor(stream);writable:=tsNodeWritableStateFor(stream)
	if readable!=nil {
		readable.objectMode=tsTruthy(tsNodeReadableGet(loop,options,"objectMode"))||tsTruthy(tsNodeReadableGet(loop,options,"readableObjectMode"))
		readable.highWaterMark=int(tsNumber(tsInternalStreamStateGetHighWaterMark(loop,readable.object,options,tsNodeDuplexKey("readableHighWaterMark"),true)))
		tsNodeReadableSyncState(readable)
		if tsStrictEqual(tsNodeReadableGet(loop,options,"readable"),tsBooleanValue(false)){readable.readable=false;readable.ended=true;readable.endEmitted=true;tsNodeReadableSyncState(readable)}
	}
	if writable!=nil {
		writable.objectMode=tsTruthy(tsNodeWritableGet(loop,options,"objectMode"))||tsTruthy(tsNodeWritableGet(loop,options,"writableObjectMode"))
		writable.highWaterMark=int(tsNumber(tsInternalStreamStateGetHighWaterMark(loop,writable.object,options,tsNodeDuplexKey("writableHighWaterMark"),true)))
		tsNodeWritableSync(writable)
		if tsStrictEqual(tsNodeWritableGet(loop,options,"writable"),tsBooleanValue(false)){writable.writableEnabled=false;writable.ending=true;writable.ended=true;writable.finished=true;tsNodeWritableSync(writable)}
	}
	allowHalfOpen:=!tsStrictEqual(tsNodeDuplexGet(loop,options,"allowHalfOpen"),tsBooleanValue(false));tsNodeDuplexSet(loop,stream,"allowHalfOpen",tsBooleanValue(allowHalfOpen))
	if !allowHalfOpen { listener:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{if state:=tsNodeWritableStateFor(stream);state!=nil&&!state.ending{tsNodeWritableEnd(loop,state,nil)};return tsU}));tsCallReceiver(loop,tsNodeDuplexGet(loop,stream,"once"),stream,tsNodeDuplexKey("end"),listener) }
	if signal:=tsNodeDuplexGet(loop,options,"signal");!tsNullish(signal){tsInternalStreamAddAbortSignal(loop,signal,stream)}
	if construct:=tsNodeDuplexGet(loop,options,"construct");construct.kind==tsFunctionKind {
		tsNodeDuplexSet(loop,stream,"_construct",construct)
		if readable!=nil{readable.constructed=false;tsNodeReadableSyncState(readable)};if writable!=nil{writable.constructed=false;tsNodeWritableSync(writable)}
		called:=false;callback:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{if called{return tsU};called=true;err:=tsArg(args,0);if readable!=nil{readable.constructed=true;tsNodeReadableSyncState(readable)};if writable!=nil{writable.constructed=true;tsNodeWritableSync(writable)};if tsTruthy(err){tsInternalStreamDestroy(loop,stream,err,tsU)}else{if readable!=nil&&readable.needReadable{tsNodeReadableScheduleRead(loop,readable)};if writable!=nil{tsNodeWritableFlush(loop,writable)}};return tsU}))
		tsCallReceiver(loop,construct,stream,callback)
	}
	return stream
}
func tsNodeDuplexGet(loop *tsLoop,value tsValue,name string) tsValue { return tsGet(loop,value,tsNodeDuplexKey(name)) }
func tsNodeDuplexSet(loop *tsLoop,value tsValue,name string,item tsValue) { tsSet(loop,value,tsNodeDuplexKey(name),item) }
func tsNodeDuplexInstall(loop *tsLoop) {
	p:=tsNodeDuplexSingleton.prototype
	if p.descriptors==nil { p.descriptors=map[string]*tsDescriptor{} }
	p.descriptors["destroyed"]=&tsDescriptor{accessor:true,configurable:true,getter:tsNamedNativeMethod("get destroyed",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{readable:=tsNodeReadableStateFor(receiver);writable:=tsNodeWritableStateFor(receiver);if readable==nil||writable==nil{return tsBooleanValue(false)};return tsBooleanValue(readable.destroyed&&writable.destroyed)})}
	p.descriptors["destroyed"].setter=tsNamedNativeMethod("set destroyed",1,func(_ *tsLoop,receiver tsValue,args ...tsValue)tsValue{value:=tsTruthy(tsArg(args,0));if readable:=tsNodeReadableStateFor(receiver);readable!=nil{readable.destroyed=value;tsNodeReadableSyncState(readable)};if writable:=tsNodeWritableStateFor(receiver);writable!=nil{writable.destroyed=value;tsNodeWritableSync(writable)};return tsU})
}
func tsNodeDuplexModuleValue(loop *tsLoop) tsValue { tsNodeReadableInit(loop);module:=tsNewObject();module.set("default",tsFunctionValue(tsNodeReadableSingleton.constructor));module.set("Readable",tsFunctionValue(tsNodeReadableSingleton.constructor));return tsObjectValue(module) }
func tsInternalStreamDuplexModule(loop *tsLoop) tsValue { tsNodeDuplexInit(loop);module:=tsNewObject();module.set("default",tsFunctionValue(tsNodeDuplexSingleton.constructor));module.set("Duplex",tsFunctionValue(tsNodeDuplexSingleton.constructor));return tsObjectValue(module) }
`
