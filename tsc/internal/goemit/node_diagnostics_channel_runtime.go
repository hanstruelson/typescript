package goemit

// nodeDiagnosticsChannelRuntime maps diagnostics channels to retained Go
// channel objects. Finalization-based cache cleanup and async-context stores
// use the runtime's simpler strong references and callback adapters.
const nodeDiagnosticsChannelRuntime = `
type tsNodeChannelStore struct { key, transform tsValue }
type tsNodeChannel struct { object *tsObject; name tsValue; stores []tsNodeChannelStore }
type tsNodeDiagnosticsState struct { sync.Once; value *tsObject; channels map[string]*tsNodeChannel; objects map[*tsObject]*tsNodeChannel; prototype *tsObject }
var tsNodeDiagnosticsSingleton tsNodeDiagnosticsState

func tsNodeDiagnosticsInit() {
	tsNodeDiagnosticsSingleton.Do(func(){
		tsNodeDiagnosticsSingleton.channels=map[string]*tsNodeChannel{}
		tsNodeDiagnosticsSingleton.objects=map[*tsObject]*tsNodeChannel{}
		tsNodeDiagnosticsSingleton.prototype=tsNewObject()
	})
}

func tsNodeChannelNameKey(loop *tsLoop,name tsValue) string {
	if name.kind==tsStringKind{return "s:"+tsText(name)}
	if name.kind==tsSymbolKind{return "y:"+(*tsSymbol)(name.ref).key}
	tsFSInvalid("channel","a string or symbol")
	return ""
}

func tsNodeChannelState(value tsValue) *tsNodeChannel {
	if value.kind!=tsObjectKind{return nil}
	return tsNodeDiagnosticsSingleton.objects[(*tsObject)(value.ref)]
}

func tsNodeChannelSubscriberList(state *tsNodeChannel) tsValue {
	return tsFSProperty(tsObjectValue(state.object),"_subscribers")
}

func tsNodeChannelSubscriberLength(value tsValue) int {
	if tsNullish(value){return 0}
	length:=tsFSProperty(value,"length")
	if !tsIsNumeric(length){return 0}
	n:=int(tsNumber(length));if n<0{n=0}
	return n
}

func tsNodeChannelSubscribers(state *tsNodeChannel) []tsValue {
	value:=tsNodeChannelSubscriberList(state)
	n:=tsNodeChannelSubscriberLength(value)
	values:=make([]tsValue,n)
	for index:=range values{values[index]=tsFSProperty(value,strconv.Itoa(index))}
	return values
}

func tsNodeChannelMaybeMarkInactive(state *tsNodeChannel) {
	if len(tsNodeChannelSubscribers(state))==0&&len(state.stores)==0 {
		state.object.set("_subscribers",tsU)
	}
}

func tsNodeChannelCreate(loop *tsLoop,name tsValue) tsValue {
	key:=tsNodeChannelNameKey(loop,name)
	object:=tsNewObject();object.nativeClass=tsNativeClass("Channel");object.prototype=tsObjectValue(tsNodeDiagnosticsSingleton.prototype)
	state:=&tsNodeChannel{object:object,name:name};tsNodeDiagnosticsSingleton.objects[object]=state;tsNodeDiagnosticsSingleton.channels[key]=state
	object.set("name",name)
	object.set("_subscribers",tsU)
	object.set("subscribe",tsNamedNativeMethod("subscribe",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{
		state:=tsNodeChannelState(receiver);if state==nil{tsPropertyFailure("Channel method called on incompatible receiver")};callback:=tsArg(args,0);if callback.kind!=tsFunctionKind{tsFSInvalid("subscription","a function")};subscribers:=tsNodeChannelSubscribers(state);subscribers=append(subscribers,callback);state.object.set("_subscribers",tsArrayValue(&tsArray{values:subscribers}));return tsU
	}))
	object.set("unsubscribe",tsNamedNativeMethod("unsubscribe",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{
		state:=tsNodeChannelState(receiver);if state==nil{tsPropertyFailure("Channel method called on incompatible receiver")};callback:=tsArg(args,0);subscribers:=tsNodeChannelSubscribers(state)
		for index:=0;index<len(subscribers);index++{if tsStrictEqual(subscribers[index],callback){subscribers=append(subscribers[:index],subscribers[index+1:]...);state.object.set("_subscribers",tsArrayValue(&tsArray{values:subscribers}));tsNodeChannelMaybeMarkInactive(state);return tsBooleanValue(true)}}
		return tsBooleanValue(false)
	}))
	object.set("publish",tsNamedNativeMethod("publish",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeChannelState(receiver);if state!=nil{tsNodeChannelPublish(loop,state,tsArg(args,0))};return tsU}))
	object.set("bindStore",tsNamedNativeMethod("bindStore",2,func(_ *tsLoop,receiver tsValue,args ...tsValue)tsValue{
		state:=tsNodeChannelState(receiver);if state==nil{tsPropertyFailure("Channel method called on incompatible receiver")};store,transform:=tsArg(args,0),tsArg(args,1);if tsIsUndefined(transform){transform=tsU}else if transform.kind!=tsFunctionKind{tsFSInvalid("transform","a function")};if tsIsUndefined(tsFSProperty(receiver,"_subscribers")){state.object.set("_subscribers",tsArrayValue(&tsArray{}))}
		for index:=range state.stores{if tsStrictEqual(state.stores[index].key,store){state.stores[index].transform=transform;return tsU}}
		state.stores=append(state.stores,tsNodeChannelStore{store,transform});return tsU
	}))
	object.set("unbindStore",tsNamedNativeMethod("unbindStore",1,func(_ *tsLoop,receiver tsValue,args ...tsValue)tsValue{
		state:=tsNodeChannelState(receiver);if state==nil{tsPropertyFailure("Channel method called on incompatible receiver")};store:=tsArg(args,0)
		for index:=range state.stores{if tsStrictEqual(state.stores[index].key,store){state.stores=append(state.stores[:index],state.stores[index+1:]...);tsNodeChannelMaybeMarkInactive(state);return tsBooleanValue(true)}}
		return tsBooleanValue(false)
	}))
	object.set("runStores",tsNamedNativeMethod("runStores",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeChannelState(receiver);if state==nil{tsPropertyFailure("Channel method called on incompatible receiver")};values:=[]tsValue{};if len(args)>3{values=args[3:]};return tsNodeChannelRunStores(loop,state,tsArg(args,0),tsArg(args,1),tsArg(args,2),values)}))
	object.order=append(object.order,"hasSubscribers")
	object.descriptors=map[string]*tsDescriptor{"hasSubscribers":{accessor:true,enumerable:true,configurable:true,getter:tsNamedNativeMethod("get hasSubscribers",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeChannelState(receiver);return tsBooleanValue(state!=nil&&(tsNodeChannelSubscriberLength(tsNodeChannelSubscriberList(state))>0||len(state.stores)>0))})}}
	return tsObjectValue(object)
}

func tsNodeDiagnosticsChannel(loop *tsLoop,name tsValue) tsValue {
	tsNodeDiagnosticsInit();key:=tsNodeChannelNameKey(loop,name)
	if state:=tsNodeDiagnosticsSingleton.channels[key];state!=nil{return tsObjectValue(state.object)}
	return tsNodeChannelCreate(loop,name)
}

func tsNodeChannelPublish(loop *tsLoop,state *tsNodeChannel,data tsValue) {
	listeners:=tsNodeChannelSubscriberList(state)
	for index:=0;index<tsNodeChannelSubscriberLength(listeners);index++{
		callback:=tsFSProperty(listeners,strconv.Itoa(index))
		func(){defer func(){if failure:=recover();failure!=nil{loop.errors=append(loop.errors,tsUnwrap(failure))}}();tsCall(loop,callback,data,state.name)}()
	}
}

func tsNodeChannelRunStores(loop *tsLoop,state *tsNodeChannel,data,callback,thisArg tsValue,args []tsValue) tsValue {
	if callback.kind!=tsFunctionKind{tsFSInvalid("fn","a function")}
	var run func(int)tsValue
	run=func(index int)tsValue{
		if index<0{tsNodeChannelPublish(loop,state,data);return tsCallReceiver(loop,callback,thisArg,args...)}
		store:=state.stores[index]
		context:=data
		if !tsIsUndefined(store.transform){context=tsCall(loop,store.transform,data)}
		next:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{return run(index-1)}))
		method:=tsFSProperty(store.key,"run")
		if method.kind!=tsFunctionKind{tsPropertyFailure("A bound store must have a run method")}
		return tsCallReceiver(loop,method,store.key,context,next)
	}
	return run(len(state.stores)-1)
}

func tsNodeDiagnosticsChannelConstructor() tsValue {
	function:=tsFunc(func(_ *tsLoop,_ ...tsValue)tsValue{tsPropertyFailure("Channel must be constructed with new");return tsU})
	function.name="Channel";function.length=1;function.constructible=true
	function.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeChannelCreate(loop,tsArg(args,0))}
	function.properties=tsNewObject();function.properties.set("prototype",tsObjectValue(tsNodeDiagnosticsSingleton.prototype));function.properties.descriptors=map[string]*tsDescriptor{"prototype":{}}
	return tsFunctionValue(function)
}

func tsNodeTraceChannel(loop *tsLoop,nameOrChannels tsValue) tsValue {
	trace:=tsNewObject()
	keys:=[]string{"start","end","asyncStart","asyncEnd","error"}
	if nameOrChannels.kind==tsStringKind {
		for _,name:=range keys{trace.set(name,tsNodeDiagnosticsChannel(loop,tsStringReference(tsStringUTF8("tracing:"+tsText(nameOrChannels)+":"+name))))}
	} else if nameOrChannels.kind==tsObjectKind||nameOrChannels.kind==tsInstanceKind {
		for _,name:=range keys{channel:=tsFSProperty(nameOrChannels,name);if tsNodeChannelState(channel)==nil{tsFSInvalid("nameOrChannels."+name,"a Channel")};trace.set(name,channel)}
	} else {tsFSInvalid("nameOrChannels","a string, object, or Channel")}
	trace.set("subscribe",tsNamedNativeMethod("subscribe",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{
		handlers:=tsArg(args,0);for _,name:=range keys{callback:=tsFSProperty(handlers,name);if callback.kind==tsFunctionKind{tsCallReceiver(loop,tsFSProperty(tsFSProperty(receiver,name),"subscribe"),tsFSProperty(receiver,name),callback)}};return tsU
	}))
	trace.set("unsubscribe",tsNamedNativeMethod("unsubscribe",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{
		handlers:=tsArg(args,0);done:=true;for _,name:=range keys{callback:=tsFSProperty(handlers,name);if callback.kind==tsFunctionKind{removed:=tsCallReceiver(loop,tsFSProperty(tsFSProperty(receiver,name),"unsubscribe"),tsFSProperty(receiver,name),callback);if !tsTruthy(removed){done=false}}};return tsBooleanValue(done)
	}))
	trace.set("traceSync",tsNamedNativeMethod("traceSync",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{
		fn:=tsArg(args,0);if fn.kind!=tsFunctionKind{tsFSInvalid("fn","a function")};context:=tsArg(args,1);if tsIsUndefined(context){context=tsObjectValue(tsNewObject())};thisArg:=tsArg(args,2);callArgs:=[]tsValue{};if len(args)>3{callArgs=args[3:]}
		start,end,fail:=tsFSProperty(receiver,"start"),tsFSProperty(receiver,"end"),tsFSProperty(receiver,"error")
		return tsCallReceiver(loop,tsFSProperty(start,"runStores"),start,context,tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{result:=tsU;func(){defer func(){if failure:=recover();failure!=nil{err:=tsUnwrap(failure);tsSet(loop,context,tsStringReference(tsStringUTF8("error")),err);tsCallReceiver(loop,tsFSProperty(fail,"publish"),fail,context);tsNodeTracePublish(loop,end,context);panic(failure)}}();result=tsCallReceiver(loop,fn,thisArg,callArgs...);tsSet(loop,context,tsStringReference(tsStringUTF8("result")),result);tsNodeTracePublish(loop,end,context)}();return result})),thisArg,callArgs...)
	}))
	trace.set("tracePromise",tsNamedNativeMethod("tracePromise",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{
		fn:=tsArg(args,0);if fn.kind!=tsFunctionKind{tsFSInvalid("fn","a function")};context:=tsArg(args,1);if tsIsUndefined(context){context=tsObjectValue(tsNewObject())};thisArg:=tsArg(args,2);callArgs:=[]tsValue{};if len(args)>3{callArgs=args[3:]}
		start,end,asyncStart,asyncEnd,fail:=tsFSProperty(receiver,"start"),tsFSProperty(receiver,"end"),tsFSProperty(receiver,"asyncStart"),tsFSProperty(receiver,"asyncEnd"),tsFSProperty(receiver,"error")
		tsNodeTracePublish(loop,start,context)
		var result tsValue;func(){defer func(){if failure:=recover();failure!=nil{err:=tsUnwrap(failure);tsSet(loop,context,tsStringReference(tsStringUTF8("error")),err);tsNodeTracePublish(loop,fail,context);tsNodeTracePublish(loop,end,context);panic(failure)}}();result=tsCallReceiver(loop,fn,thisArg,callArgs...)}()
		tsNodeTracePublish(loop,end,context);promise:=loop.promise();loop.await(result,func(settled tsResult){if settled.rejected{tsSet(loop,context,tsStringReference(tsStringUTF8("error")),settled.value);tsNodeTracePublish(loop,fail,context);tsNodeTracePublish(loop,asyncStart,context);tsNodeTracePublish(loop,asyncEnd,context);promise.settle(settled)}else{tsSet(loop,context,tsStringReference(tsStringUTF8("result")),settled.value);tsNodeTracePublish(loop,asyncStart,context);tsNodeTracePublish(loop,asyncEnd,context);promise.settle(settled)}});return tsPromiseValue(promise)
	}))
	trace.set("traceCallback",tsNamedNativeMethod("traceCallback",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{
		fn:=tsArg(args,0);if fn.kind!=tsFunctionKind{tsFSInvalid("fn","a function")};position:=-1;if tsIsNumeric(tsArg(args,1)){position=int(tsNumber(tsArg(args,1)))};context:=tsArg(args,2);if tsIsUndefined(context){context=tsObjectValue(tsNewObject())};thisArg:=tsArg(args,3);callArgs:=[]tsValue{};if len(args)>4{callArgs=append(callArgs,args[4:]...)}
		if position<0{position=len(callArgs)+position};if position<0{position=0};if position>=len(callArgs){tsFSInvalid("callback","a function")};callback:=callArgs[position];if callback.kind!=tsFunctionKind{tsFSInvalid("callback","a function")}
		start,end,asyncStart,asyncEnd,fail:=tsFSProperty(receiver,"start"),tsFSProperty(receiver,"end"),tsFSProperty(receiver,"asyncStart"),tsFSProperty(receiver,"asyncEnd"),tsFSProperty(receiver,"error")
		wrapped:=tsFunctionValue(tsFunc(func(loop *tsLoop,values ...tsValue)tsValue{err:=tsArg(values,0);if tsTruthy(err){tsSet(loop,context,tsStringReference(tsStringUTF8("error")),err);tsNodeTracePublish(loop,fail,context)}else if len(values)>1{tsSet(loop,context,tsStringReference(tsStringUTF8("result")),values[1])};tsNodeTracePublish(loop,asyncStart,context);var result tsValue;func(){defer func(){if failure:=recover();failure!=nil{tsNodeTracePublish(loop,asyncEnd,context);panic(failure)}}();result=tsCallReceiver(loop,callback,thisArg,values...)}();tsNodeTracePublish(loop,asyncEnd,context);return result}))
		callArgs[position]=wrapped
		tracedFn:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{var result tsValue;func(){defer func(){if failure:=recover();failure!=nil{err:=tsUnwrap(failure);tsSet(loop,context,tsStringReference(tsStringUTF8("error")),err);tsNodeTracePublish(loop,fail,context);tsNodeTracePublish(loop,end,context);panic(failure)}}();result=tsCallReceiver(loop,fn,thisArg,callArgs...)}();tsNodeTracePublish(loop,end,context);return result}))
		return tsNodeChannelRunStores(loop,tsNodeChannelState(start),context,tracedFn,thisArg,nil)
	}))
	return tsObjectValue(trace)
}

func tsNodeTracePublish(loop *tsLoop,channel,data tsValue) { state:=tsNodeChannelState(channel);if state!=nil{tsNodeChannelPublish(loop,state,data)} }

func tsNodeDiagnosticsModule(loop *tsLoop) tsValue {
	tsNodeDiagnosticsInit()
	if tsNodeDiagnosticsSingleton.value!=nil{return tsObjectValue(tsNodeDiagnosticsSingleton.value)}
	module:=tsNewObject()
	module.set("channel",tsNamedNativeMethod("channel",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeDiagnosticsChannel(loop,tsArg(args,0))}))
	module.set("hasSubscribers",tsNamedNativeMethod("hasSubscribers",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{name:=tsArg(args,0);key:=tsNodeChannelNameKey(loop,name);state:=tsNodeDiagnosticsSingleton.channels[key];return tsBooleanValue(state!=nil&&(tsNodeChannelSubscriberLength(tsNodeChannelSubscriberList(state))>0||len(state.stores)>0))}))
	module.set("subscribe",tsNamedNativeMethod("subscribe",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{channel:=tsNodeDiagnosticsChannel(loop,tsArg(args,0));tsCallReceiver(loop,tsFSProperty(channel,"subscribe"),channel,tsArg(args,1));return tsU}))
	module.set("unsubscribe",tsNamedNativeMethod("unsubscribe",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{channel:=tsNodeDiagnosticsChannel(loop,tsArg(args,0));return tsCallReceiver(loop,tsFSProperty(channel,"unsubscribe"),channel,tsArg(args,1))}))
	module.set("tracingChannel",tsNamedNativeMethod("tracingChannel",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeTraceChannel(loop,tsArg(args,0))}))
	module.set("Channel",tsNodeDiagnosticsChannelConstructor())
	tsNodeDiagnosticsSingleton.value=module
	return tsObjectValue(module)
}
`
