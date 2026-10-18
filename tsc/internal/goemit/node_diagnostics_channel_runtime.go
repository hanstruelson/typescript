package goemit

// nodeDiagnosticsChannelRuntime maps diagnostics channels through weak Go
// references and adapts Bun's subscription, store, and tracing wrappers.
const nodeDiagnosticsChannelRuntime = `
type tsNodeChannel struct { object *tsObject; name tsValue }
type tsNodeChannelRef struct { id uint64; state weak.Pointer[tsNodeChannel] }
type tsNodeChannelCleanup struct { key tsCollectionKey; id uint64 }
type tsNodeDiagnosticsState struct { sync.Once; sync.Mutex; value *tsObject; channels map[tsCollectionKey]tsNodeChannelRef; nextChannelID atomic.Uint64; prototype,activePrototype *tsObject; constructor tsValue }
var tsNodeDiagnosticsSingleton tsNodeDiagnosticsState

func tsNodeDiagnosticsInit(loop *tsLoop) {
	tsNodeDiagnosticsSingleton.Do(func(){
		tsNodeDiagnosticsSingleton.channels=map[tsCollectionKey]tsNodeChannelRef{}
		tsNodeDiagnosticsSingleton.prototype=tsNewObject()
		tsNodeDiagnosticsSingleton.activePrototype=tsNewObject()
		tsNodeDiagnosticsSingleton.activePrototype.prototype=tsObjectValue(tsNodeDiagnosticsSingleton.prototype)
		tsNodeChannelInstallPrototypes(loop)
	})
}

func tsNodeChannelFinalized(cleanup tsNodeChannelCleanup){tsNodeDiagnosticsSingleton.Lock();if current,ok:=tsNodeDiagnosticsSingleton.channels[cleanup.key];ok&&current.id==cleanup.id{delete(tsNodeDiagnosticsSingleton.channels,cleanup.key)};tsNodeDiagnosticsSingleton.Unlock()}
func tsNodeDiagnosticsLookup(key tsCollectionKey)*tsNodeChannel{tsNodeDiagnosticsSingleton.Lock();defer tsNodeDiagnosticsSingleton.Unlock();reference,exists:=tsNodeDiagnosticsSingleton.channels[key];if !exists{return nil};state:=reference.state.Value();if state==nil{delete(tsNodeDiagnosticsSingleton.channels,key)};return state}

func tsNodeChannelNameKey(loop *tsLoop,name tsValue) tsCollectionKey { return tsCollectionHash(loop,name) }

func tsNodeChannelState(value tsValue) *tsNodeChannel {
	if value.kind!=tsObjectKind{return nil}
	state,_:=(*tsObject)(value.ref).nativeData.(*tsNodeChannel);return state
}

func tsNodeChannelIsInstance(value tsValue)bool{if value.kind!=tsObjectKind{return false};prototype:=tsPrototypeOf(value);return tsStrictEqual(prototype,tsObjectValue(tsNodeDiagnosticsSingleton.prototype))||tsStrictEqual(prototype,tsObjectValue(tsNodeDiagnosticsSingleton.activePrototype))}

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
	stores:=state.object.get("_stores")
	storeCount:=0
	if stores.kind==tsCollectionKind{storeCount=len((*tsCollection)(stores.ref).index)}
	if len(tsNodeChannelSubscribers(state))==0&&storeCount==0 {
		state.object.prototype=tsObjectValue(tsNodeDiagnosticsSingleton.prototype)
		state.object.set("_subscribers",tsU)
		state.object.set("_stores",tsU)
	}
}

func tsNodeChannelActivate(state *tsNodeChannel){if tsIsUndefined(state.object.get("_subscribers")){state.object.prototype=tsObjectValue(tsNodeDiagnosticsSingleton.activePrototype);state.object.set("_subscribers",tsArrayValue(&tsArray{}));state.object.set("_stores",tsNewCollection(false))}}

func tsNodeChannelInstallPrototypes(loop *tsLoop){
	base:=tsNodeDiagnosticsSingleton.prototype
	active:=tsNodeDiagnosticsSingleton.activePrototype
	base.descriptors=map[string]*tsDescriptor{}
	active.descriptors=map[string]*tsDescriptor{}
	base.set("subscribe",tsNamedNativeMethod("subscribe",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeChannelState(receiver);if state==nil{tsPropertyFailure("Channel method called on incompatible receiver")};tsNodeChannelActivate(state);method:=tsFSProperty(receiver,"subscribe");return tsCallReceiver(loop,method,receiver,tsArg(args,0))}))
	base.set("unsubscribe",tsNamedNativeMethod("unsubscribe",1,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsBooleanValue(false)}))
	base.set("bindStore",tsNamedNativeMethod("bindStore",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeChannelState(receiver);if state==nil{tsPropertyFailure("Channel method called on incompatible receiver")};tsNodeChannelActivate(state);method:=tsFSProperty(receiver,"bindStore");return tsCallReceiver(loop,method,receiver,tsArg(args,0),tsArg(args,1))}))
	base.set("unbindStore",tsNamedNativeMethod("unbindStore",1,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsBooleanValue(false)}))
	base.set("publish",tsNamedNativeMethod("publish",1,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsU}))
	base.set("runStores",tsNamedNativeMethod("runStores",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{tail:=[]tsValue{};if len(args)>3{tail=args[3:]};return tsCallReceiver(loop,tsArg(args,1),tsArg(args,2),tail...)}))
	base.order=append(base.order,"hasSubscribers")
	base.descriptors["hasSubscribers"]=&tsDescriptor{accessor:true,configurable:true,getter:tsNamedNativeMethod("get hasSubscribers",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsBooleanValue(false)})}
	active.set("subscribe",tsNamedNativeMethod("subscribe",1,func(_ *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeChannelState(receiver);if state==nil{tsPropertyFailure("Channel method called on incompatible receiver")};callback:=tsArg(args,0);if callback.kind!=tsFunctionKind&&callback.kind!=tsClassKind{tsInternalValidatorInvalidType("subscription","a function",callback)};array:=(*tsArray)(state.object.get("_subscribers").ref);array.values=append(array.values,callback);return tsU}))
	active.set("unsubscribe",tsNamedNativeMethod("unsubscribe",1,func(_ *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeChannelState(receiver);if state==nil{tsPropertyFailure("Channel method called on incompatible receiver")};callback:=tsArg(args,0);array:=(*tsArray)(state.object.get("_subscribers").ref);for index,subscription:=range array.values{if tsStrictEqual(subscription,callback){copy(array.values[index:],array.values[index+1:]);array.values[len(array.values)-1]=tsU;array.values=array.values[:len(array.values)-1];tsNodeChannelMaybeMarkInactive(state);return tsBooleanValue(true)}};return tsBooleanValue(false)}))
	active.set("bindStore",tsNamedNativeMethod("bindStore",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeChannelState(receiver);if state==nil{tsPropertyFailure("Channel method called on incompatible receiver")};stores:=(*tsCollection)(state.object.get("_stores").ref);stores.put(loop,tsArg(args,0),tsArg(args,1));return tsU}))
	active.set("unbindStore",tsNamedNativeMethod("unbindStore",1,func(_ *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeChannelState(receiver);if state==nil{tsPropertyFailure("Channel method called on incompatible receiver")};stores:=state.object.get("_stores");if stores.kind!=tsCollectionKind{tsPropertyFailure("this._stores.has is not a function")};collection:=(*tsCollection)(stores.ref);hash:=tsCollectionHash(tsArg(args,0));entry,exists:=collection.index[hash];if !exists{return tsBooleanValue(false)};delete(collection.index,hash);collection.entries[entry]=tsCollectionEntry{};tsNodeChannelMaybeMarkInactive(state);return tsBooleanValue(true)}))
	active.set("publish",tsNamedNativeMethod("publish",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeChannelState(receiver);if state==nil{tsPropertyFailure("Channel method called on incompatible receiver")};tsNodeChannelPublish(loop,state,tsArg(args,0));return tsU}))
	active.set("runStores",tsNamedNativeMethod("runStores",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeChannelState(receiver);if state==nil{tsPropertyFailure("Channel method called on incompatible receiver")};tail:=[]tsValue{};if len(args)>3{tail=args[3:]};return tsNodeChannelRunStores(loop,state,tsArg(args,0),tsArg(args,1),tsArg(args,2),tail)}))
	active.order=append(active.order,"hasSubscribers")
	active.descriptors["hasSubscribers"]=&tsDescriptor{accessor:true,configurable:true,getter:tsNamedNativeMethod("get hasSubscribers",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsBooleanValue(true)})}
	for _,prototype:=range []*tsObject{base,active}{for _,name:=range prototype.order{if prototype.descriptors[name]==nil{prototype.descriptors[name]=&tsDescriptor{writable:true,configurable:true}}}}
}

func tsNodeChannelCreate(loop *tsLoop,name tsValue) tsValue {
		key:=tsNodeChannelNameKey(loop,name)
	object:=tsNewObject();object.prototype=tsObjectValue(tsNodeDiagnosticsSingleton.prototype)
	state:=&tsNodeChannel{object:object,name:name};object.nativeData=state;id:=tsNodeDiagnosticsSingleton.nextChannelID.Add(1);tsNodeDiagnosticsSingleton.Lock();tsNodeDiagnosticsSingleton.channels[key]=tsNodeChannelRef{id:id,state:weak.Make(state)};tsNodeDiagnosticsSingleton.Unlock();runtime.AddCleanup(state,tsNodeChannelFinalized,tsNodeChannelCleanup{key:key,id:id})
	object.set("_subscribers",tsU)
	object.set("_stores",tsU)
	object.set("name",name)
	return tsObjectValue(object)
}

func tsNodeDiagnosticsChannel(loop *tsLoop,name tsValue) tsValue {
	tsNodeDiagnosticsInit(loop);key:=tsNodeChannelNameKey(loop,name)
	if state:=tsNodeDiagnosticsLookup(key);state!=nil{return tsObjectValue(state.object)}
	if name.kind!=tsStringKind&&name.kind!=tsSymbolKind{tsFSInvalid("channel","a string or symbol")}
	return tsNodeChannelCreate(loop,name)
}

func tsNodeChannelPublish(loop *tsLoop,state *tsNodeChannel,data tsValue) {
	for index:=0;index<tsNodeChannelSubscriberLength(tsNodeChannelSubscriberList(state));index++{
		callback:=tsFSProperty(tsNodeChannelSubscriberList(state),strconv.Itoa(index))
		func(){defer func(){if failure:=recover();failure!=nil{failureValue:=tsUnwrap(failure);process:=tsNodeProcessModule(loop);report:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{tsNodeNetReportError(loop,failureValue);return tsU}));tsNodeProcessNextTick(loop,(*tsObject)(process.ref),[]tsValue{report})}}();tsCall(loop,callback,data,state.name)}()
	}
}

func tsNodeChannelWrapStore(loop *tsLoop,data tsValue,next func()tsValue,store,transform tsValue)func()tsValue{
	return func()tsValue{
		context:=data
		if !tsIsUndefined(transform){failed:=false;var failure tsValue;func(){defer func(){if caught:=recover();caught!=nil{failed=true;failure=tsUnwrap(caught)}}();context=tsCall(loop,transform,data)}();if failed{process:=tsNodeProcessModule(loop);report:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{tsNodeNetReportError(loop,failure);return tsU}));tsNodeProcessNextTick(loop,(*tsObject)(process.ref),[]tsValue{report});return next()}}
		continuation:=tsFunctionValue(tsFunc(func(_ *tsLoop,_ ...tsValue)tsValue{return next()}))
		method:=tsFSProperty(store,"run")
		return tsCallReceiver(loop,method,store,context,continuation)
	}
}

func tsNodeChannelRunStores(loop *tsLoop,state *tsNodeChannel,data,callback,thisArg tsValue,args []tsValue) tsValue {
	run:=func()tsValue{tsNodeChannelPublish(loop,state,data);return tsCallReceiver(loop,callback,thisArg,args...)}
	storesValue:=state.object.get("_stores")
	if storesValue.kind!=tsCollectionKind{if tsStrictEqual(state.object.prototype,tsObjectValue(tsNodeDiagnosticsSingleton.activePrototype)){tsPropertyFailure("this._stores.entries is not a function")};return run()}
	stores:=(*tsCollection)(storesValue.ref)
	for index:=0;index<len(stores.entries);index++{
		entry:=stores.entries[index];if !entry.live{continue}
		run=tsNodeChannelWrapStore(loop,data,run,entry.key,entry.value)
	}
	return run()
}

func tsNodeDiagnosticsPromiseResolve(loop *tsLoop,value tsValue)tsValue{
	if value.kind==tsPromiseKind{return value}
	promise:=loop.promise();settled:=false
	var resolve func(tsValue)
	reject:=func(reason tsValue){if settled{return};settled=true;promise.settle(tsResult{reason,true})}
	resolve=func(resolution tsValue){
		if settled{return}
		if resolution.kind==tsPromiseKind{
			if resolution.ref==unsafe.Pointer(promise){reject(tsErrorValue(&tsRuntimeError{name:"TypeError",message:"Chaining cycle detected for promise"}));return}
			loop.await(resolution,func(result tsResult){if result.rejected{reject(result.value)}else{resolve(result.value)}});return
		}
		if resolution.kind==tsObjectKind||resolution.kind==tsInstanceKind||resolution.kind==tsArrayKind||resolution.kind==tsFunctionKind{
			var then tsValue;failed:=false;var failure any
			func(){defer func(){if caught:=recover();caught!=nil{failed=true;failure=caught}}();then=tsFSProperty(resolution,"then")}()
			if failed{reject(tsUnwrap(failure));return}
			if then.kind==tsFunctionKind{
				called:=false
				resolveFunction:=tsFunctionValue(tsFunc(func(_ *tsLoop,args ...tsValue)tsValue{if !called{called=true;resolve(tsArg(args,0))};return tsU}))
				rejectFunction:=tsFunctionValue(tsFunc(func(_ *tsLoop,args ...tsValue)tsValue{if !called{called=true;reject(tsArg(args,0))};return tsU}))
				func(){defer func(){if caught:=recover();caught!=nil&&!called{called=true;reject(tsUnwrap(caught))}}();tsCallReceiver(loop,then,resolution,resolveFunction,rejectFunction)}()
				return
			}
		}
		if !settled{settled=true;promise.settle(tsResult{resolution,false})}
	}
	resolve(value)
	return tsPromiseValue(promise)
}

func tsNodeTraceExecute(loop *tsLoop,context,errorChannel,endChannel tsValue,storeResult bool,call func()tsValue)tsValue{
	var result tsValue;var thrown any;failed:=false
	func(){defer func(){if failure:=recover();failure!=nil{failed=true;thrown=failure}}();result=call();if storeResult{tsSet(loop,context,tsStringReference(tsStringUTF8("result")),result)}}()
	if failed{
		reason:=tsUnwrap(thrown);var catchFailure any
		func(){defer func(){catchFailure=recover()}();tsSet(loop,context,tsStringReference(tsStringUTF8("error")),reason);tsNodeTracePublish(loop,errorChannel,context)}()
		tsNodeTracePublish(loop,endChannel,context)
		if catchFailure!=nil{panic(catchFailure)};panic(thrown)
	}
	tsNodeTracePublish(loop,endChannel,context)
	return result
}

func tsNodeTraceSettle(loop *tsLoop,promise *tsPromise,context,errorChannel,asyncStart,asyncEnd tsValue,result tsResult){
	defer func(){if failure:=recover();failure!=nil{promise.settle(tsResult{tsUnwrap(failure),true})}}()
	if result.rejected{tsSet(loop,context,tsStringReference(tsStringUTF8("error")),result.value);tsNodeTracePublish(loop,errorChannel,context)}else{tsSet(loop,context,tsStringReference(tsStringUTF8("result")),result.value)}
	tsNodeTracePublish(loop,asyncStart,context);tsNodeTracePublish(loop,asyncEnd,context);promise.settle(result)
}

func tsNodeDiagnosticsChannelConstructor() tsValue {
	tsNodeDiagnosticsInit(loop)
	if !tsIsUndefined(tsNodeDiagnosticsSingleton.constructor){return tsNodeDiagnosticsSingleton.constructor}
	function:=tsFunc(func(_ *tsLoop,_ ...tsValue)tsValue{tsPropertyFailure("Channel must be constructed with new");return tsU})
	function.name="Channel";function.length=1;function.constructible=true
	function.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeChannelCreate(loop,tsArg(args,0))}
	function.properties=tsNewObject();function.properties.set("prototype",tsObjectValue(tsNodeDiagnosticsSingleton.prototype))
	hasInstanceKey:=tsPropertyKey(tsWellKnownSymbol("hasInstance"));function.properties.set(hasInstanceKey,tsNamedNativeMethod("[Symbol.hasInstance]",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsBooleanValue(tsNodeChannelIsInstance(tsArg(args,0)))}))
	function.properties.descriptors=map[string]*tsDescriptor{"prototype":{},hasInstanceKey:{writable:true,configurable:true}}
	tsNodeDiagnosticsSingleton.prototype.set("constructor",tsFunctionValue(function));tsNodeDiagnosticsSingleton.activePrototype.set("constructor",tsFunctionValue(function))
	tsNodeDiagnosticsSingleton.prototype.descriptors["constructor"]=&tsDescriptor{writable:true,configurable:true};tsNodeDiagnosticsSingleton.activePrototype.descriptors["constructor"]=&tsDescriptor{writable:true,configurable:true}
	tsNodeDiagnosticsSingleton.constructor=tsFunctionValue(function)
	return tsNodeDiagnosticsSingleton.constructor
}

func tsNodeTraceChannel(loop *tsLoop,nameOrChannels tsValue) tsValue {
	trace:=tsNewObject()
	keys:=[]string{"start","end","asyncStart","asyncEnd","error"}
	if nameOrChannels.kind==tsStringKind {
		for _,name:=range keys{trace.set(name,tsNodeDiagnosticsChannel(loop,tsStringReference(tsStringUTF8("tracing:"+tsText(nameOrChannels)+":"+name))))}
	} else if nameOrChannels.kind==tsObjectKind||nameOrChannels.kind==tsInstanceKind {
		for _,name:=range keys{channel:=tsFSProperty(nameOrChannels,name);if !tsNodeChannelIsInstance(channel){tsFSInvalid("nameOrChannels."+name,"a Channel")};trace.set(name,channel)}
	} else {tsFSInvalid("nameOrChannels","a string, object, or Channel")}
	trace.set("subscribe",tsNamedNativeMethod("subscribe",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{
		handlers:=tsArg(args,0);if tsNullish(handlers){tsPropertyFailure("Cannot read properties of undefined")};for _,name:=range keys{callback:=tsFSProperty(handlers,name);if !tsTruthy(callback){continue};channel:=tsFSProperty(receiver,name);if tsNullish(channel){continue};tsCallReceiver(loop,tsFSProperty(channel,"subscribe"),channel,callback)};return tsU
	}))
	trace.set("unsubscribe",tsNamedNativeMethod("unsubscribe",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{
		handlers:=tsArg(args,0);if tsNullish(handlers){tsPropertyFailure("Cannot read properties of undefined")};done:=true;for _,name:=range keys{callback:=tsFSProperty(handlers,name);if !tsTruthy(callback){continue};channel:=tsFSProperty(receiver,name);if tsNullish(channel)||!tsTruthy(tsCallReceiver(loop,tsFSProperty(channel,"unsubscribe"),channel,callback)){done=false}};return tsBooleanValue(done)
	}))
	trace.set("traceSync",tsNamedNativeMethod("traceSync",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{
		fn:=tsArg(args,0);context:=tsArg(args,1);if tsIsUndefined(context){context=tsObjectValue(tsNewObject())};thisArg:=tsArg(args,2);callArgs:=[]tsValue{};if len(args)>3{callArgs=args[3:]}
		start,end,fail:=tsFSProperty(receiver,"start"),tsFSProperty(receiver,"end"),tsFSProperty(receiver,"error")
		body:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{return tsNodeTraceExecute(loop,context,fail,end,true,func()tsValue{return tsCallReceiver(loop,fn,thisArg,callArgs...)})}))
		return tsCallReceiver(loop,tsFSProperty(start,"runStores"),start,context,body)
	}))
	trace.set("tracePromise",tsNamedNativeMethod("tracePromise",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{
		fn:=tsArg(args,0);context:=tsArg(args,1);if tsIsUndefined(context){context=tsObjectValue(tsNewObject())};thisArg:=tsArg(args,2);callArgs:=[]tsValue{};if len(args)>3{callArgs=args[3:]}
		start,end,asyncStart,asyncEnd,fail:=tsFSProperty(receiver,"start"),tsFSProperty(receiver,"end"),tsFSProperty(receiver,"asyncStart"),tsFSProperty(receiver,"asyncEnd"),tsFSProperty(receiver,"error")
		body:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{result:=tsNodeTraceExecute(loop,context,fail,end,false,func()tsValue{return tsCallReceiver(loop,fn,thisArg,callArgs...)});promise:=loop.promise();loop.await(tsNodeDiagnosticsPromiseResolve(loop,result),func(settled tsResult){tsNodeTraceSettle(loop,promise,context,fail,asyncStart,asyncEnd,settled)});return tsPromiseValue(promise)}))
		return tsCallReceiver(loop,tsFSProperty(start,"runStores"),start,context,body)
	}))
	trace.set("traceCallback",tsNamedNativeMethod("traceCallback",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{
		fn:=tsArg(args,0);position:=-1;if !tsIsUndefined(tsArg(args,1)){positionNumber:=math.Trunc(tsNumber(tsArg(args,1)));if math.IsNaN(positionNumber){position=0}else if math.IsInf(positionNumber,1){position=1<<30}else if math.IsInf(positionNumber,-1){position=-(1<<30)}else{position=int(positionNumber)}};context:=tsArg(args,2);if tsIsUndefined(context){context=tsObjectValue(tsNewObject())};thisArg:=tsArg(args,3);callArgs:=[]tsValue{};if len(args)>4{callArgs=append(callArgs,args[4:]...)}
		callbackIndex:=position;if callbackIndex<0{callbackIndex+=len(callArgs)};if callbackIndex<0||callbackIndex>=len(callArgs){tsInternalValidatorInvalidType("callback","a function",tsU)};callback:=callArgs[callbackIndex];if callback.kind!=tsFunctionKind&&callback.kind!=tsClassKind{tsInternalValidatorInvalidType("callback","a function",callback)}
		start,end,asyncStart,asyncEnd,fail:=tsFSProperty(receiver,"start"),tsFSProperty(receiver,"end"),tsFSProperty(receiver,"asyncStart"),tsFSProperty(receiver,"asyncEnd"),tsFSProperty(receiver,"error")
		wrappedCall:=func(loop *tsLoop,this tsValue,values []tsValue)tsValue{err:=tsArg(values,0);if tsTruthy(err){tsSet(loop,context,tsStringReference(tsStringUTF8("error")),err);tsNodeTracePublish(loop,fail,context)}else{tsSet(loop,context,tsStringReference(tsStringUTF8("result")),tsArg(values,1))};inner:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)(result tsValue){defer tsNodeTracePublish(loop,asyncEnd,context);return tsCallReceiver(loop,callback,this,values...)}));return tsCallReceiver(loop,tsFSProperty(asyncStart,"runStores"),asyncStart,context,inner)}
		wrapped:=tsFunctionValue(tsFunc(func(loop *tsLoop,values ...tsValue)tsValue{return wrappedCall(loop,tsU,values)}));(*tsFunction)(wrapped.ref).receiverCall=func(loop *tsLoop,this tsValue,values ...tsValue)tsValue{return wrappedCall(loop,this,values)}
		spliceIndex:=position;if spliceIndex<0{spliceIndex+=len(callArgs)};if spliceIndex<0{spliceIndex=0};if spliceIndex>len(callArgs){spliceIndex=len(callArgs)};callArgs=append(callArgs[:spliceIndex],append([]tsValue{wrapped},callArgs[spliceIndex+1:]...)...)
		tracedFn:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{return tsNodeTraceExecute(loop,context,fail,end,false,func()tsValue{return tsCallReceiver(loop,fn,thisArg,callArgs...)})}))
		return tsCallReceiver(loop,tsFSProperty(start,"runStores"),start,context,tracedFn)
	}))
	return tsObjectValue(trace)
}

func tsNodeTracePublish(loop *tsLoop,channel,data tsValue) { tsCallReceiver(loop,tsFSProperty(channel,"publish"),channel,data) }

func tsNodeDiagnosticsModule(loop *tsLoop) tsValue {
	tsNodeDiagnosticsInit(loop)
	if tsNodeDiagnosticsSingleton.value!=nil{return tsObjectValue(tsNodeDiagnosticsSingleton.value)}
	module:=tsNewObject()
	module.set("channel",tsNamedNativeMethod("channel",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeDiagnosticsChannel(loop,tsArg(args,0))}))
	module.set("hasSubscribers",tsNamedNativeMethod("hasSubscribers",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{name:=tsArg(args,0);key:=tsNodeChannelNameKey(loop,name);state:=tsNodeDiagnosticsLookup(key);return tsBooleanValue(state!=nil&&tsTruthy(tsFSProperty(tsObjectValue(state.object),"hasSubscribers")))}))
	module.set("subscribe",tsNamedNativeMethod("subscribe",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{channel:=tsNodeDiagnosticsChannel(loop,tsArg(args,0));tsCallReceiver(loop,tsFSProperty(channel,"subscribe"),channel,tsArg(args,1));return tsU}))
	module.set("unsubscribe",tsNamedNativeMethod("unsubscribe",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{channel:=tsNodeDiagnosticsChannel(loop,tsArg(args,0));return tsCallReceiver(loop,tsFSProperty(channel,"unsubscribe"),channel,tsArg(args,1))}))
	module.set("tracingChannel",tsNamedNativeMethod("tracingChannel",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeTraceChannel(loop,tsArg(args,0))}))
	module.set("Channel",tsNodeDiagnosticsChannelConstructor())
	tsNodeDiagnosticsSingleton.value=module
	return tsObjectValue(module)
}
`
