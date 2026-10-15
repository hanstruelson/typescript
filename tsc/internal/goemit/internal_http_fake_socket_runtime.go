package goemit

// internalHTTPFakeSocketRuntime maps Bun's in-process HTTP socket facade onto
// the shared Duplex stream implementation.
const internalHTTPFakeSocketRuntime = `
var tsInternalHTTPFakeSocketSingleton struct { sync.Once; constructor *tsFunction; prototype *tsObject; module tsValue; socketData,serverSymbol,autoDestroyed tsValue }
func tsInternalHTTPFakeSocketKey(name string)tsValue{return tsStringReference(tsStringUTF8(name))}
func tsInternalHTTPFakeSocketGet(loop *tsLoop,value tsValue,name string)tsValue{return tsGet(loop,value,tsInternalHTTPFakeSocketKey(name))}
func tsInternalHTTPFakeSocketData(loop *tsLoop,receiver tsValue)tsValue{return tsGet(loop,receiver,tsInternalHTTPFakeSocketSingleton.socketData)}
func tsInternalHTTPFakeSocketAddress(loop *tsLoop,receiver tsValue)tsValue{
	cached:=tsInternalHTTPFakeSocketGet(loop,receiver,"_fakeSocketAddress")
	if !tsIsUndefined(cached){return cached}
	address:=tsNewObject();data:=tsInternalHTTPFakeSocketData(loop,receiver)
	if data.kind==tsArrayKind&&len((*tsArray)(data.ref).values)>2{
		server:=tsGet(loop,data,tsStringReference(tsStringUTF8("0")));request:=tsGet(loop,data,tsStringReference(tsStringUTF8("2")))
		handle:=tsGet(loop,server,tsInternalHTTPFakeSocketSingleton.serverSymbol)
		requestIP:=tsInternalHTTPFakeSocketGet(loop,handle,"requestIP")
		if requestIP.kind==tsFunctionKind{result:=tsCallReceiver(loop,requestIP,handle,request);if !tsNullish(result){return tsSet(loop,receiver,tsInternalHTTPFakeSocketKey("_fakeSocketAddress"),result)}}
	}
	return tsSet(loop,receiver,tsInternalHTTPFakeSocketKey("_fakeSocketAddress"),tsObjectValue(address))
}
func tsInternalHTTPFakeSocketCreate(loop *tsLoop,httpMessage tsValue)tsValue{
	tsInternalHTTPFakeSocketInit(loop)
	duplex:=tsInternalStreamDuplexModule(loop);constructor:=tsInternalHTTPFakeSocketGet(loop,duplex,"Duplex")
	stream:=tsConstruct(loop,constructor);object:=(*tsObject)(stream.ref);object.prototype=tsObjectValue(tsInternalHTTPFakeSocketSingleton.prototype)
	tsSet(loop,stream,tsInternalHTTPFakeSocketKey("_httpMessage"),httpMessage)
	tsSet(loop,stream,tsInternalHTTPFakeSocketKey("bytesRead"),tsNumberValue(0));tsSet(loop,stream,tsInternalHTTPFakeSocketKey("bytesWritten"),tsNumberValue(0))
	tsSet(loop,stream,tsInternalHTTPFakeSocketKey("connecting"),tsBooleanValue(false));tsSet(loop,stream,tsInternalHTTPFakeSocketKey("timeout"),tsNumberValue(0));tsSet(loop,stream,tsInternalHTTPFakeSocketKey("isServer"),tsBooleanValue(false))
	tsSet(loop,stream,tsInternalHTTPFakeSocketKey("_onTimeout"),tsNamedNativeMethod("_onTimeout",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{emit:=tsInternalHTTPFakeSocketGet(loop,receiver,"emit");if emit.kind==tsFunctionKind{tsCallReceiver(loop,emit,receiver,tsInternalHTTPFakeSocketKey("timeout"))};return tsU}))
	return stream
}
func tsInternalHTTPFakeSocketInit(loop *tsLoop){
	tsInternalHTTPFakeSocketSingleton.Do(func(){
		http:=tsInternalHTTPModule(loop);shared:=tsInternalSharedModule(loop)
		tsInternalHTTPFakeSocketSingleton.socketData=tsInternalHTTPFakeSocketGet(loop,http,"kInternalSocketData")
		tsInternalHTTPFakeSocketSingleton.serverSymbol=tsInternalHTTPFakeSocketGet(loop,http,"serverSymbol")
		tsInternalHTTPFakeSocketSingleton.autoDestroyed=tsInternalHTTPFakeSocketGet(loop,shared,"kAutoDestroyed")
		duplex:=tsInternalStreamDuplexModule(loop);duplexCtor:=tsInternalHTTPFakeSocketGet(loop,duplex,"Duplex");duplexProto:=tsInternalHTTPFakeSocketGet(loop,duplexCtor,"prototype")
		prototype:=tsNewObject();prototype.prototype=duplexProto
		constructor:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsInternalHTTPFakeSocketCreate(loop,tsArg(args,0))})
		constructor.name="Socket";constructor.length=1;constructor.constructible=true
		constructor.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsInternalHTTPFakeSocketCreate(loop,tsArg(args,0))}
		constructor.receiverCall=func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsInternalHTTPFakeSocketCreate(loop,tsArg(args,0))}
		constructor.properties=tsNewObject();constructor.properties.set("prototype",tsObjectValue(prototype));constructor.properties.set("name",tsInternalHTTPFakeSocketKey("Socket"));constructor.properties.descriptors=map[string]*tsDescriptor{"prototype":{writable:false},"name":{configurable:true}}
		prototype.set("constructor",tsFunctionValue(constructor));prototype.descriptors=map[string]*tsDescriptor{"constructor":{writable:true,configurable:true}}
		define:=func(name string,length int,call func(*tsLoop,tsValue,...tsValue)tsValue){prototype.set(name,tsNamedNativeMethod(name,length,call));prototype.descriptors[name]=&tsDescriptor{writable:true,configurable:true}}
		define("address",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{return tsInternalHTTPFakeSocketAddress(loop,receiver)})
		define("connect",3,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{return receiver})
		define("_destroy",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{
			data:=tsInternalHTTPFakeSocketData(loop,receiver)
			if data.kind==tsArrayKind&&len((*tsArray)(data.ref).values)>1{
				response:=tsGet(loop,data,tsInternalHTTPFakeSocketKey("1"));request:=tsInternalHTTPFakeSocketGet(loop,response,"req")
				if tsNullish(request)||!tsTruthy(tsGet(loop,request,tsInternalHTTPFakeSocketSingleton.autoDestroyed)){end:=tsInternalHTTPFakeSocketGet(loop,response,"end");if end.kind==tsFunctionKind{tsCallReceiver(loop,end,response)}}
			}
			callback:=tsArg(args,1);if callback.kind==tsFunctionKind{tsCall(loop,callback,tsArg(args,0))};return tsU
		})
		define("_final",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{callback:=tsArg(args,0);if callback.kind==tsFunctionKind{tsCall(loop,callback,tsU)};return tsU})
		define("_read",1,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsU})
		define("ref",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{return receiver})
		define("resetAndDestroy",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{return receiver})
		define("setKeepAlive",2,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{return receiver})
		define("setNoDelay",1,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{return receiver})
		define("setTimeout",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{
			data:=tsInternalHTTPFakeSocketData(loop,receiver);if data.kind!=tsArrayKind||len((*tsArray)(data.ref).values)<2{return tsU}
			response:=tsGet(loop,data,tsInternalHTTPFakeSocketKey("1"));request:=tsInternalHTTPFakeSocketGet(loop,response,"req");if tsNullish(request){return receiver};setTimeout:=tsInternalHTTPFakeSocketGet(loop,request,"setTimeout")
			if setTimeout.kind==tsFunctionKind{tsCallReceiver(loop,setTimeout,request,tsArg(args,0),tsArg(args,1))};return receiver
		})
		define("unref",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{return receiver})
		define("_write",3,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{callback:=tsArg(args,2);if callback.kind==tsFunctionKind{tsCall(loop,callback,tsU)};return tsU})
		define("destroy",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{
			message:=tsInternalHTTPFakeSocketGet(loop,receiver,"_httpMessage");destroy:=tsU;if !tsNullish(message){destroy=tsInternalHTTPFakeSocketGet(loop,message,"destroy")};if destroy.kind==tsFunctionKind{tsCallReceiver(loop,destroy,message)}
			baseDestroy:=tsInternalHTTPFakeSocketGet(loop,duplexProto,"destroy");if baseDestroy.kind==tsFunctionKind{return tsCallReceiver(loop,baseDestroy,receiver)};return receiver
		})
		accessor:=func(name string,getter func(*tsLoop,tsValue)tsValue,setter func(*tsLoop,tsValue,tsValue)tsValue){descriptor:=&tsDescriptor{accessor:true,configurable:true,getter:tsNamedNativeMethod("get "+name,0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{return getter(loop,receiver)})};if setter!=nil{descriptor.setter=tsNamedNativeMethod("set "+name,1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return setter(loop,receiver,tsArg(args,0))})};prototype.descriptors[name]=descriptor}
		accessor("bufferSize",func(loop *tsLoop,receiver tsValue)tsValue{return tsInternalHTTPFakeSocketGet(loop,receiver,"writableLength")},nil)
		accessor("localAddress",func(loop *tsLoop,receiver tsValue)tsValue{address:=tsInternalHTTPFakeSocketAddress(loop,receiver);if tsTruthy(address){return tsInternalHTTPFakeSocketKey("127.0.0.1")};return tsU},nil)
		accessor("localFamily",func(_ *tsLoop,_ tsValue)tsValue{return tsInternalHTTPFakeSocketKey("IPv4")},nil)
		accessor("localPort",func(_ *tsLoop,_ tsValue)tsValue{return tsNumberValue(80)},nil)
		accessor("pending",func(loop *tsLoop,receiver tsValue)tsValue{return tsInternalHTTPFakeSocketGet(loop,receiver,"connecting")},nil)
		accessor("readyState",func(loop *tsLoop,receiver tsValue)tsValue{if tsTruthy(tsInternalHTTPFakeSocketGet(loop,receiver,"connecting")){return tsInternalHTTPFakeSocketKey("opening")};if tsTruthy(tsInternalHTTPFakeSocketGet(loop,receiver,"readable")){if tsTruthy(tsInternalHTTPFakeSocketGet(loop,receiver,"writable")){return tsInternalHTTPFakeSocketKey("open")};return tsInternalHTTPFakeSocketKey("readOnly")};if tsTruthy(tsInternalHTTPFakeSocketGet(loop,receiver,"writable")){return tsInternalHTTPFakeSocketKey("writeOnly")};return tsInternalHTTPFakeSocketKey("closed")},nil)
		for _,field:=range []string{"address","port","family"}{property:=field;accessor("remote"+strings.ToUpper(property[:1])+property[1:],func(loop *tsLoop,receiver tsValue)tsValue{return tsInternalHTTPFakeSocketGet(loop,tsInternalHTTPFakeSocketAddress(loop,receiver),property)},func(loop *tsLoop,receiver,value tsValue)tsValue{return tsSet(loop,tsInternalHTTPFakeSocketAddress(loop,receiver),tsInternalHTTPFakeSocketKey(property),value)})}
		tsInternalHTTPFakeSocketSingleton.prototype=prototype;tsInternalHTTPFakeSocketSingleton.constructor=constructor
		module:=tsNewObject();module.set("FakeSocket",tsFunctionValue(constructor));module.set("default",tsObjectValue(module));tsInternalHTTPFakeSocketSingleton.module=tsObjectValue(module)
	})
}
func tsInternalHTTPFakeSocketModule(loop *tsLoop)tsValue{tsInternalHTTPFakeSocketInit(loop);return tsInternalHTTPFakeSocketSingleton.module}
`
