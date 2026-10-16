package goemit

// nodeEventsRuntime implements EventEmitter and the common static helpers.
// Bun's validator, abort-listener and async-resource glue is represented by
// direct checks against this runtime's existing values and AbortSignal state.
const nodeEventsRuntime = `
type tsNodeEventListener struct { original, raw tsValue; once bool }
type tsNodeEventEmitter struct { object *tsObject; events map[string][]*tsNodeEventListener; names map[string]tsValue; max int; capture bool }
type tsNodeEventsState struct { sync.Once; constructor *tsFunction; prototype *tsObject; emitters map[*tsObject]*tsNodeEventEmitter; defaultMax int; capture bool; errorMonitor, rejectionSymbol, maxTargetListeners, maxTargetWarned tsValue }
var tsNodeEventsSingleton tsNodeEventsState

func tsNodeEventsInit(loop *tsLoop) {
	tsNodeEventsSingleton.Do(func() {
		tsNodeEventsSingleton.emitters=map[*tsObject]*tsNodeEventEmitter{}
		tsNodeEventsSingleton.defaultMax=10
		tsNodeEventsSingleton.prototype=tsNewObject()
		registry:=tsGet(loop,tsSymbolFunction(),tsStringReference(tsStringUTF8("for")))
		tsNodeEventsSingleton.errorMonitor=tsCall(loop,registry,tsStringReference(tsStringUTF8("events.errorMonitor")))
		tsNodeEventsSingleton.rejectionSymbol=tsCall(loop,registry,tsStringReference(tsStringUTF8("nodejs.rejection")))
		tsNodeEventsSingleton.maxTargetListeners=tsNewSymbol(tsStringUTF8("events.maxEventTargetListeners"))
		tsNodeEventsSingleton.maxTargetWarned=tsNewSymbol(tsStringUTF8("events.maxEventTargetListenersWarned"))
		tsNodeEventsInstallPrototype(loop)
	})
}

func tsNodeEventKey(loop *tsLoop, event tsValue) string {
	if event.kind==tsSymbolKind { return tsPropertyKey(loop,event) }
	if event.kind!=tsStringKind { tsFSInvalid("eventName", "a string or symbol") }
	return tsText(event)
}

func tsNodeEventState(receiver tsValue) *tsNodeEventEmitter {
	if receiver.kind!=tsObjectKind { tsPropertyFailure("EventEmitter method called on incompatible receiver") }
	return tsNodeEventsSingleton.emitters[(*tsObject)(receiver.ref)]
}

func tsNodeEventsAddTargetListener(loop *tsLoop,target,event,listener tsValue,once bool) {
	if state:=tsNodeEventState(target);state!=nil { tsNodeEventAdd(loop,target,event,listener,false,once);return }
	if tsEventTargetStateFor(target)==nil { tsFSInvalid("emitter","an EventEmitter or EventTarget") }
	options:=tsNewObject();options.set("once",tsBooleanValue(once));tsEventTargetAdd(loop,target,event,listener,tsObjectValue(options))
}
func tsNodeEventsRemoveTargetListener(loop *tsLoop,target,event,listener tsValue) {
	if state:=tsNodeEventState(target);state!=nil { tsNodeEventRemove(loop,state,event,listener);return }
	if tsEventTargetStateFor(target)==nil { tsFSInvalid("emitter","an EventEmitter or EventTarget") }
	tsEventTargetRemove(loop,target,event,listener,tsU)
}

func tsNodeEventCreate(loop *tsLoop, receiver tsValue, options tsValue) tsValue {
	var object *tsObject
	if receiver.kind==tsObjectKind { object=(*tsObject)(receiver.ref) }
	if object==nil { object=tsNewObject() }
	if object.nativeClass==nil { object.nativeClass=tsNativeClass("EventEmitter") }
	if tsIsUndefined(object.prototype) { object.prototype=tsObjectValue(tsNodeEventsSingleton.prototype) }
	capture:=tsNodeEventsSingleton.capture
	if options.kind==tsObjectKind||options.kind==tsInstanceKind {
		value:=tsFSProperty(options,"captureRejections")
		if !tsIsUndefined(value) { if value.kind!=tsBooleanKind { tsFSInvalid("options.captureRejections","a boolean") }; capture=tsTruthy(value) }
	} else if !tsIsUndefined(options) { tsFSInvalid("options","an object") }
	state:=&tsNodeEventEmitter{object:object,events:map[string][]*tsNodeEventListener{},names:map[string]tsValue{},max:tsNodeEventsSingleton.defaultMax,capture:capture}
	tsNodeEventsSingleton.emitters[object]=state
	return tsObjectValue(object)
}

func tsNodeEventsInstallPrototype(loop *tsLoop) {
	prototype:=tsNodeEventsSingleton.prototype
	prototype.set("setMaxListeners",tsNamedNativeMethod("setMaxListeners",1,func(_ *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeEventState(receiver);if state==nil{tsPropertyFailure("Invalid EventEmitter receiver")};value:=tsArg(args,0);if !tsIsNumeric(value)||tsNumber(value)<0||math.IsNaN(tsNumber(value)){tsFSInvalid("n","a non-negative number")};state.max=int(tsNumber(value));return receiver}))
	prototype.set("getMaxListeners",tsNamedNativeMethod("getMaxListeners",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeEventState(receiver);if state==nil{tsPropertyFailure("Invalid EventEmitter receiver")};return tsNumberValue(float64(state.max))}))
	prototype.set("emit",tsNamedNativeMethod("emit",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeEventState(receiver);if state==nil{tsPropertyFailure("Invalid EventEmitter receiver")};return tsBooleanValue(tsNodeEventEmit(loop,state,tsArg(args,0),args[1:]))}))
	prototype.set("addListener",tsNamedNativeMethod("addListener",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeEventAdd(loop,receiver,tsArg(args,0),tsArg(args,1),false,false)}))
	prototype.set("on",prototype.get("addListener"))
	prototype.set("prependListener",tsNamedNativeMethod("prependListener",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeEventAdd(loop,receiver,tsArg(args,0),tsArg(args,1),true,false)}))
	prototype.set("once",tsNamedNativeMethod("once",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeEventAdd(loop,receiver,tsArg(args,0),tsArg(args,1),false,true)}))
	prototype.set("prependOnceListener",tsNamedNativeMethod("prependOnceListener",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeEventAdd(loop,receiver,tsArg(args,0),tsArg(args,1),true,true)}))
	prototype.set("removeListener",tsNamedNativeMethod("removeListener",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeEventState(receiver);if state==nil{tsPropertyFailure("Invalid EventEmitter receiver")};tsNodeEventRemove(loop,state,tsArg(args,0),tsArg(args,1));return receiver}))
	prototype.set("off",prototype.get("removeListener"))
	prototype.set("removeAllListeners",tsNamedNativeMethod("removeAllListeners",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeEventState(receiver);if state==nil{tsPropertyFailure("Invalid EventEmitter receiver")};event:=tsArg(args,0);if tsIsUndefined(event){state.events=map[string][]*tsNodeEventListener{};state.names=map[string]tsValue{};tsNodeProcessMemoryPressureListenerChange(loop,state.object,"memoryPressure");tsNodeProcessSignalListenerChange(loop,state.object,"")}else{key:=tsNodeEventKey(loop,event);delete(state.events,key);delete(state.names,key);tsNodeProcessMemoryPressureListenerChange(loop,state.object,key);tsNodeProcessSignalListenerChange(loop,state.object,key)};return receiver}))
	prototype.set("listeners",tsNamedNativeMethod("listeners",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeEventListeners(loop,receiver,tsArg(args,0),false)}))
	prototype.set("rawListeners",tsNamedNativeMethod("rawListeners",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeEventListeners(loop,receiver,tsArg(args,0),true)}))
	prototype.set("listenerCount",tsNamedNativeMethod("listenerCount",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeEventCount(loop,receiver,tsArg(args,0),tsArg(args,1))}))
	prototype.set("eventNames",tsNamedNativeMethod("eventNames",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeEventState(receiver);out:=&tsArray{};if state!=nil{for key,list:=range state.events{if len(list)>0{out.values=append(out.values,state.names[key])}}};return tsArrayValue(out)}))
}

func tsNodeEventAdd(loop *tsLoop,receiver,event,listener tsValue,prepend,once bool) tsValue {
	state:=tsNodeEventState(receiver);if state==nil{tsPropertyFailure("Invalid EventEmitter receiver")}
	if listener.kind!=tsFunctionKind{tsFSInvalid("listener","a function")}
	key:=tsNodeEventKey(loop,event)
	entry:=&tsNodeEventListener{original:listener,raw:listener,once:once}
	if once {
		var wrapper tsValue
		wrapper=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{
			tsNodeEventRemove(loop,state,event,wrapper)
			return tsCallReceiver(loop,listener,tsObjectValue(state.object),args...)
		}))
		(*tsFunction)(wrapper.ref).name="bound onceWrapper"
		(*tsFunction)(wrapper.ref).properties=tsNewObject()
		(*tsFunction)(wrapper.ref).properties.set("listener",listener)
		entry.raw=wrapper
	}
	list:=state.events[key]
	if prepend { list=append([]*tsNodeEventListener{entry},list...) } else { list=append(list,entry) }
	state.events[key]=list;state.names[key]=event
	tsNodeProcessMemoryPressureListenerChange(loop,state.object,key)
	tsNodeProcessSignalListenerChange(loop,state.object,key)
	return receiver
}

func tsNodeEventRemove(loop *tsLoop,state *tsNodeEventEmitter,event,listener tsValue) {
	key:=tsNodeEventKey(loop,event);list:=state.events[key]
	for index:=len(list)-1;index>=0;index-- {
		entry:=list[index];matches:=tsStrictEqual(entry.original,listener)||tsStrictEqual(entry.raw,listener)
		if !matches&&entry.raw.kind==tsFunctionKind { matches=tsStrictEqual(tsGet(loop,entry.raw,tsStringReference(tsStringUTF8("listener"))),listener) }
		if matches { list=append(list[:index],list[index+1:]...);break }
	}
	if len(list)==0 { delete(state.events,key);delete(state.names,key) } else { state.events[key]=list }
	tsNodeProcessMemoryPressureListenerChange(loop,state.object,key)
	tsNodeProcessSignalListenerChange(loop,state.object,key)
}

func tsNodeEventEmit(loop *tsLoop,state *tsNodeEventEmitter,event tsValue,args []tsValue) bool {
	key:=tsNodeEventKey(loop,event)
	if tsText(event)=="error" {
		monitor:=state.events[tsPropertyKey(loop,tsNodeEventsSingleton.errorMonitor)]
		for index:=0;index<len(monitor);index++{var entry *tsNodeEventListener;entry=monitor[index];tsCallReceiver(loop,entry.raw,tsObjectValue(state.object),args...)}
	}
	listeners:=append([]*tsNodeEventListener(nil),state.events[key]...)
	if len(listeners)==0 {
		if tsText(event)=="error" {
			value:=tsArg(args,0);if value.kind==tsErrorKind{panic(tsThrown{value})}
			err:=tsErrorValue(&tsRuntimeError{name:"Error",message:"Unhandled 'error' event",fields:map[string]tsValue{"context":value}});panic(tsThrown{err})
		}
		return false
	}
	for index:=0;index<len(listeners);index++ {
		var entry *tsNodeEventListener;entry=listeners[index]
		if entry.once { tsNodeEventRemove(loop,state,event,entry.raw) }
		result:=tsCallReceiver(loop,entry.raw,tsObjectValue(state.object),args...)
		if state.capture&&result.kind==tsPromiseKind {
			promise:=(*tsPromise)(result.ref)
			promise.then(func(settled tsResult){if settled.rejected{tsNodeEventEmit(loop,state,tsStringReference(tsStringUTF8("error")),[]tsValue{settled.value})}})
		}
	}
	return true
}

func tsNodeEventListeners(loop *tsLoop,receiver,event tsValue,raw bool) tsValue {
	state:=tsNodeEventState(receiver);out:=&tsArray{};if state==nil{values,ok:=tsEventTargetListenerValues(loop,receiver,event);if !ok{tsFSInvalid("emitter","an EventEmitter or EventTarget")};out.values=values;return tsArrayValue(out)}
	listeners:=state.events[tsNodeEventKey(loop,event)];for index:=0;index<len(listeners);index++{var entry *tsNodeEventListener;entry=listeners[index];value:=entry.original;if raw{value=entry.raw};out.values=append(out.values,value)}
	return tsArrayValue(out)
}

func tsNodeEventCount(loop *tsLoop,receiver,event,listener tsValue) tsValue {
	state:=tsNodeEventState(receiver);if state==nil{count,ok:=tsEventTargetListenerCount(loop,receiver,event);if !ok{tsFSInvalid("emitter","an EventEmitter or EventTarget")};return tsNumberValue(float64(count))}
	count:=0;listeners:=state.events[tsNodeEventKey(loop,event)];for index:=0;index<len(listeners);index++{var entry *tsNodeEventListener;entry=listeners[index];if tsIsUndefined(listener)||tsStrictEqual(entry.original,listener){count++}}
	return tsNumberValue(float64(count))
}

func tsNodeEventsOnce(loop *tsLoop,args []tsValue) tsValue {
	emitter,event,options:=tsArg(args,0),tsArg(args,1),tsArg(args,2)
	state:=tsNodeEventState(emitter);target:=tsEventTargetStateFor(emitter);if state==nil&&target==nil{tsFSInvalid("emitter","an EventEmitter or EventTarget")}
	promise:=loop.promise();signal:=tsFSAbortSignal(options)
	if signal!=nil{if _,aborted:=tsAbortReason(signal);aborted{promise.settle(tsResult{tsFSAbortError(signal),true});return tsPromiseValue(promise)}}
	var listener tsValue
	var errorListener tsValue
	listener=tsFunctionValue(tsFunc(func(loop *tsLoop,values ...tsValue)tsValue{
			tsNodeEventsRemoveTargetListener(loop,emitter,event,listener)
			if !tsIsUndefined(errorListener){tsNodeEventsRemoveTargetListener(loop,emitter,tsStringReference(tsStringUTF8("error")),errorListener)}
		if signal!=nil{if _,aborted:=tsAbortReason(signal);aborted{promise.settle(tsResult{tsFSAbortError(signal),true});return tsU}}
		promise.settle(tsResult{tsArrayValue(&tsArray{values:append([]tsValue(nil),values...)}),false})
		return tsU
	}))
	tsNodeEventsAddTargetListener(loop,emitter,event,listener,true)
	if state!=nil&&tsText(event)!="error" {
		errorListener=tsFunctionValue(tsFunc(func(loop *tsLoop,values ...tsValue)tsValue{tsNodeEventRemove(loop,state,event,listener);tsNodeEventRemove(loop,state,tsStringReference(tsStringUTF8("error")),errorListener);promise.settle(tsResult{tsArg(values,0),true});return tsU}))
		tsNodeEventAdd(loop,emitter,tsStringReference(tsStringUTF8("error")),errorListener,false,true)
	}
	if signal!=nil{signal.followers=append(signal.followers,func(_ tsValue){tsNodeEventsRemoveTargetListener(loop,emitter,event,listener);if !tsIsUndefined(errorListener)&&state!=nil{tsNodeEventRemove(loop,state,tsStringReference(tsStringUTF8("error")),errorListener)};promise.settle(tsResult{tsFSAbortError(signal),true})})}
	return tsPromiseValue(promise)
}

func tsNodeEventsOn(loop *tsLoop,args []tsValue) tsValue {
	emitter,event,options:=tsArg(args,0),tsArg(args,1),tsArg(args,2)
	state:=tsNodeEventState(emitter);target:=tsEventTargetStateFor(emitter);if state==nil&&target==nil{tsFSInvalid("emitter","an EventEmitter or EventTarget")}
	signal:=tsFSAbortSignal(options);queue:=[]tsValue{};waiters:=[]*tsPromise{};closed:=false;aborted:=false
	deliver:=func(){for len(queue)>0&&len(waiters)>0{promise:=waiters[0];waiters=waiters[1:];item:=queue[0];queue=queue[1:];promise.settle(tsResult{tsIteratorResult(item,false),false})};if aborted&&len(queue)==0{for _,promise:=range waiters{promise.settle(tsResult{tsFSAbortError(signal),true})};waiters=nil}}
	listener:=tsFunctionValue(tsFunc(func(_ *tsLoop,args ...tsValue)tsValue{queue=append(queue,tsArrayValue(&tsArray{values:append([]tsValue(nil),args...)}));deliver();return tsU}))
	tsNodeEventsAddTargetListener(loop,emitter,event,listener,false)
	iterator:=&tsAsyncIterator{
		next:func(loop *tsLoop)*tsPromise{if len(queue)>0{value:=queue[0];queue=queue[1:];return loop.resolved(tsIteratorResult(value,false),false)};if closed{return loop.resolved(tsIteratorResult(tsU,true),false)};if aborted{return loop.resolved(tsFSAbortError(signal),true)};promise:=loop.promise();waiters=append(waiters,promise);return promise},
			close:func(loop *tsLoop)*tsPromise{if !closed{closed=true;tsNodeEventsRemoveTargetListener(loop,emitter,event,listener);for _,promise:=range waiters{promise.settle(tsResult{tsIteratorResult(tsU,true),false})};waiters=nil};return loop.resolved(tsIteratorResult(tsU,true),false)},
	}
	if signal!=nil{signal.followers=append(signal.followers,func(_ tsValue){aborted=true;tsNodeEventsRemoveTargetListener(loop,emitter,event,listener);deliver()})}
	return tsAsyncIteratorObject(iterator)
}

func tsNodeEventAsyncConstructor() tsValue {
	function:=tsFunc(func(_ *tsLoop,_ ...tsValue)tsValue{tsPropertyFailure("EventEmitterAsyncResource must be constructed with new");return tsU})
	function.name="EventEmitterAsyncResource";function.constructible=true
	function.construct=func(loop *tsLoop,args ...tsValue)tsValue{
		value:=tsNodeEventCreate(loop,tsU,tsArg(args,0));object:=(*tsObject)(value.ref)
		object.set("asyncId",tsNumberValue(0));object.set("triggerAsyncId",tsNumberValue(0));object.set("asyncResource",value)
		object.set("emitDestroy",tsNamedNativeMethod("emitDestroy",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsU}))
		return value
	}
	return tsFunctionValue(function)
}

func tsNodeEventsModule(loop *tsLoop) tsValue {
	tsNodeEventsInit(loop)
	if tsNodeEventsSingleton.constructor!=nil{return tsFunctionValue(tsNodeEventsSingleton.constructor)}
	function:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeEventCreate(loop,tsU,tsArg(args,0))})
	function.name="EventEmitter";function.length=0;function.constructible=true
	function.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeEventCreate(loop,tsU,tsArg(args,0))}
	function.receiverCall=func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeEventCreate(loop,receiver,tsArg(args,0))}
	function.properties=tsNewObject();function.properties.set("prototype",tsObjectValue(tsNodeEventsSingleton.prototype));function.properties.descriptors=map[string]*tsDescriptor{"prototype":{}}
	properties:=function.properties
	properties.order=append(properties.order,"defaultMaxListeners")
	properties.descriptors["defaultMaxListeners"]=&tsDescriptor{accessor:true,enumerable:true,configurable:true,
		getter:tsNamedNativeMethod("get defaultMaxListeners",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNumberValue(float64(tsNodeEventsSingleton.defaultMax))}),
		setter:tsNamedNativeMethod("set defaultMaxListeners",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{value:=tsArg(args,0);if !tsIsNumeric(value)||tsNumber(value)<0{tsFSInvalid("defaultMaxListeners","a non-negative number")};tsNodeEventsSingleton.defaultMax=int(tsNumber(value));return tsU}),
	}
	properties.order=append(properties.order,"captureRejections")
	properties.descriptors["captureRejections"]=&tsDescriptor{accessor:true,enumerable:true,configurable:true,
		getter:tsNamedNativeMethod("get captureRejections",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsBooleanValue(tsNodeEventsSingleton.capture)}),
		setter:tsNamedNativeMethod("set captureRejections",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{value:=tsArg(args,0);if value.kind!=tsBooleanKind{tsFSInvalid("captureRejections","a boolean")};tsNodeEventsSingleton.capture=tsTruthy(value);return tsU}),
	}
	functionValue:=tsFunctionValue(function)
	properties.set("EventEmitter",functionValue);properties.set("usingDomains",tsBooleanValue(false));properties.set("errorMonitor",tsNodeEventsSingleton.errorMonitor);properties.set("captureRejectionSymbol",tsNodeEventsSingleton.rejectionSymbol)
	properties.set("kMaxEventTargetListeners",tsNodeEventsSingleton.maxTargetListeners);properties.set("kMaxEventTargetListenersWarned",tsNodeEventsSingleton.maxTargetWarned)
	properties.set("init",functionValue)
	properties.set("once",tsNamedNativeMethod("once",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeEventsOnce(loop,args)}))
	properties.set("on",tsNamedNativeMethod("on",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeEventsOn(loop,args)}))
	properties.set("getEventListeners",tsNamedNativeMethod("getEventListeners",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeEventListeners(loop,tsArg(args,0),tsArg(args,1),false)}))
	properties.set("getMaxListeners",tsNamedNativeMethod("getMaxListeners",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{target:=tsArg(args,0);if state:=tsNodeEventState(target);state!=nil{return tsNumberValue(float64(state.max))};if tsEventTargetStateFor(target)!=nil{value:=tsGet(loop,target,tsNodeEventsSingleton.maxTargetListeners);if tsIsUndefined(value){return tsNumberValue(float64(tsNodeEventsSingleton.defaultMax))};return value};tsFSInvalid("emitter","an EventEmitter or EventTarget");return tsU}))
	properties.set("setMaxListeners",tsNamedNativeMethod("setMaxListeners",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{value:=tsArg(args,0);if tsIsUndefined(value){value=tsNumberValue(float64(tsNodeEventsSingleton.defaultMax))};if !tsIsNumeric(value)||tsNumber(value)<0||math.IsNaN(tsNumber(value)){tsFSInvalid("n","a non-negative number")};if len(args)==1{tsNodeEventsSingleton.defaultMax=int(tsNumber(value));return tsU};for _,target:=range args[1:]{if state:=tsNodeEventState(target);state!=nil{state.max=int(tsNumber(value));continue};if tsEventTargetStateFor(target)!=nil{tsSet(loop,target,tsNodeEventsSingleton.maxTargetListeners,value);tsSet(loop,target,tsNodeEventsSingleton.maxTargetWarned,tsBooleanValue(false));continue};tsFSInvalid("eventTargets","an EventEmitter or EventTarget")};return tsU}))
	properties.set("listenerCount",tsNamedNativeMethod("listenerCount",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeEventCount(loop,tsArg(args,0),tsArg(args,1),tsArg(args,2))}))
	properties.set("addAbortListener",tsNamedNativeMethod("addAbortListener",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{signal:=tsArg(args,0);listener:=tsArg(args,1);if listener.kind!=tsFunctionKind{tsFSInvalid("listener","a function")};if signal.kind!=tsObjectKind||(*tsObject)(signal.ref).abort==nil{tsFSInvalid("signal","an AbortSignal")};state:=(*tsObject)(signal.ref).abort;_,aborted:=tsAbortReason(state);if aborted{tsCall(loop,listener);return tsObjectValue(tsNewObject())};state.followers=append(state.followers,func(reason tsValue){tsCall(loop,listener,tsObjectValue(tsNewObject()))});disposable:=tsNewObject();dispose:=tsNamedNativeMethod("[Symbol.dispose]",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsU});disposable.set(tsPropertyKey(tsWellKnownSymbol("dispose")),dispose);return tsObjectValue(disposable)}))
	properties.set("EventEmitterAsyncResource",tsNodeEventAsyncConstructor())
	tsNodeEventsSingleton.constructor=function
	return functionValue
}
`
