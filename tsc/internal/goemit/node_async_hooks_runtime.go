package goemit

// nodeAsyncHooksRuntime adapts AsyncLocalStorage and AsyncResource to the
// event loop's persistent async-context frames. TickObject init hooks run
// during process.nextTick setup; other hook phases remain unsupported.
const nodeAsyncHooksRuntime = `
type tsNodeALSState struct { object *tsObject; defaultValue tsValue; name string }
type tsNodeAsyncResourceState struct { object *tsObject; context *tsAsyncContextFrame; triggerID int64; typeName string }
type tsNodeRunScopeState struct { storage *tsNodeALSState; previous tsValue; disposed bool }
var tsNodeAsyncHooksSingleton struct { sync.Once; module *tsObject; storages map[unsafe.Pointer]*tsNodeALSState; resources map[unsafe.Pointer]*tsNodeAsyncResourceState; scopes map[unsafe.Pointer]*tsNodeRunScopeState }

func tsNodeAsyncHooksInit() { tsNodeAsyncHooksSingleton.Do(func(){tsNodeAsyncHooksSingleton.module=tsNewObject();tsNodeAsyncHooksSingleton.storages=map[unsafe.Pointer]*tsNodeALSState{};tsNodeAsyncHooksSingleton.resources=map[unsafe.Pointer]*tsNodeAsyncResourceState{};tsNodeAsyncHooksSingleton.scopes=map[unsafe.Pointer]*tsNodeRunScopeState{}}) }

func tsNodeALSStateOf(receiver tsValue) *tsNodeALSState { if receiver.kind!=tsObjectKind {tsPropertyFailure("AsyncLocalStorage method called on incompatible receiver")};state:=tsNodeAsyncHooksSingleton.storages[receiver.ref];if state==nil{tsPropertyFailure("AsyncLocalStorage method called on incompatible receiver")};return state }
func tsNodeALSStore(loop *tsLoop,state *tsNodeALSState) tsValue { key:=unsafe.Pointer(state.object);context:=loop.asyncContext;if context==nil{return state.defaultValue};if context.masked!=nil&&context.masked[key]{return state.defaultValue};for frame:=context;frame!=nil;frame=frame.prev{if frame.storage==key{return frame.value}};return state.defaultValue }
func tsNodeALSPush(loop *tsLoop,state *tsNodeALSState,value tsValue) *tsAsyncContextFrame { key:=unsafe.Pointer(state.object);var masks map[unsafe.Pointer]bool;if head:=loop.asyncContext;head!=nil&&len(head.masked)>0{masks=make(map[unsafe.Pointer]bool,len(head.masked));for item,masked:=range head.masked{if masked&&item!=key{masks[item]=true}}};frame:=&tsAsyncContextFrame{storage:key,value:value,prev:loop.asyncContext,masked:masks};loop.asyncContext=frame;return frame }
func tsNodeALSMethod(name string) tsValue { return tsNamedNativeMethod(name,1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeALSStateOf(receiver);switch name{
case "getStore":return tsNodeALSStore(loop,state)
case "enterWith":tsNodeALSPush(loop,state,tsArg(args,0));return tsU
case "run":callback:=tsArg(args,1);if callback.kind!=tsFunctionKind{tsPropertyFailure("The callback argument must be of type function")};store:=tsArg(args,0);if tsSameValue(tsNodeALSStore(loop,state),store){return tsCallReceiver(loop,callback,tsU,args[2:]...)};tsNodeALSPush(loop,state,store);defer func(){loop.asyncContext=func()*tsAsyncContextFrame{for frame:=loop.asyncContext;frame!=nil;frame=frame.prev{if frame.storage==unsafe.Pointer(state.object){return frame.prev}};return nil}()}();return tsCallReceiver(loop,callback,tsU,args[2:]...)
case "exit":callback:=tsArg(args,0);if callback.kind!=tsFunctionKind{tsPropertyFailure("The callback argument must be of type function")};tsNodeALSPush(loop,state,tsU);return tsCallReceiver(loop,callback,tsU,args[1:]...)
case "disable":if loop.asyncContext!=nil{if loop.asyncContext.masked==nil{loop.asyncContext.masked=map[unsafe.Pointer]bool{}};loop.asyncContext.masked[unsafe.Pointer(state.object)]=true};return tsU
case "_enable","_propagate":return tsU
case "withScope":scope:=tsNewObject();scopeState:=&tsNodeRunScopeState{storage:state,previous:tsNodeALSStore(loop,state)};tsNodeAsyncHooksSingleton.scopes[unsafe.Pointer(scope)]=scopeState;scope.set("dispose",tsNamedNativeMethod("dispose",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{if item:=tsNodeAsyncHooksSingleton.scopes[receiver.ref];item!=nil&&!item.disposed{item.disposed=true;tsNodeALSPush(loop,item.storage,tsValue(item.previous))};return tsU}));dispose:=scope.get("dispose");scope.set(tsPropertyKey(tsWellKnownSymbol("dispose")),dispose);return tsObjectValue(scope)
};return tsU}) }

func tsNodeAsyncResourceStateOf(receiver tsValue) *tsNodeAsyncResourceState {if receiver.kind!=tsObjectKind{tsPropertyFailure("AsyncResource method called on incompatible receiver")};state:=tsNodeAsyncHooksSingleton.resources[receiver.ref];if state==nil{tsPropertyFailure("AsyncResource method called on incompatible receiver")};return state}
func tsNodeAsyncResourceMethod(name string) tsValue {return tsNamedNativeMethod(name,1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeAsyncResourceStateOf(receiver);switch name{
case "emitBefore","emitAfter":return tsBooleanValue(true)
case "asyncId":return tsNumberValue(0)
case "triggerAsyncId":return tsNumberValue(float64(state.triggerID))
case "emitDestroy":return tsU
case "runInAsyncScope":callback:=tsArg(args,0);if callback.kind!=tsFunctionKind{tsPropertyFailure("The fn argument must be a function")};thisArg:=tsArg(args,1);var result tsValue;tsWithAsyncContext(loop,state.context,func(){result=tsCallReceiver(loop,callback,thisArg,args[2:]...)});return result
case "bind":callback:=tsArg(args,0);if callback.kind!=tsFunctionKind{tsPropertyFailure("The fn argument must be a function")};thisArg:=tsArg(args,1);bound:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{var result tsValue;tsWithAsyncContext(loop,state.context,func(){result=tsCallReceiver(loop,callback,thisArg,args...)});return result}));(*tsFunction)(bound.ref).name=(*tsFunction)(callback.ref).name;(*tsFunction)(bound.ref).length=(*tsFunction)(callback.ref).length;return bound
};return tsU})}

func tsNodeAsyncLocalStorageConstructor() tsValue {
	fn:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{tsPropertyFailure("Class constructor AsyncLocalStorage cannot be invoked without new")});fn.name="AsyncLocalStorage";fn.length=0;fn.constructible=true
	fn.construct=func(loop *tsLoop,args ...tsValue)tsValue{options:=tsArg(args,0);defaultValue:=tsFSProperty(options,"defaultValue");nameValue:=tsFSProperty(options,"name");name:="";if !tsIsUndefined(nameValue){name=tsText(nameValue)};tsSetAsyncContextTrackingEnabled(true);object:=tsNewObject();state:=&tsNodeALSState{object:object,defaultValue:defaultValue,name:name};tsNodeAsyncHooksSingleton.storages[unsafe.Pointer(object)]=state
		for _,method:=range []string{"getStore","enterWith","run","exit","disable","withScope","_enable","_propagate"}{object.set(method,tsNodeALSMethod(method))};object.set("name",tsStringReference(tsStringUTF8(name)));return tsObjectValue(object)}
	static:=tsNewObject();static.set("bind",tsNamedNativeMethod("bind",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{callback:=tsArg(args,0);if callback.kind!=tsFunctionKind{tsPropertyFailure("The fn argument must be a function")};context:=loop.asyncContext;bound:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{var result tsValue;tsWithAsyncContext(loop,context,func(){result=tsCallReceiver(loop,callback,tsU,args...)});return result}));(*tsFunction)(bound.ref).name=(*tsFunction)(callback.ref).name;return bound}))
	static.set("snapshot",tsNamedNativeMethod("snapshot",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{context:=loop.asyncContext;return tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{callback:=tsArg(args,0);if callback.kind!=tsFunctionKind{tsPropertyFailure("The fn argument must be a function")};var result tsValue;tsWithAsyncContext(loop,context,func(){result=tsCallReceiver(loop,callback,tsU,args[1:]...)});return result}))}))
	fn.properties=static
	return tsFunctionValue(fn)
}

func tsNodeAsyncResourceConstructor() tsValue {
	fn:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{tsPropertyFailure("Class constructor AsyncResource cannot be invoked without new")});fn.name="AsyncResource";fn.length=1;fn.constructible=true
	fn.construct=func(loop *tsLoop,args ...tsValue)tsValue{typeValue:=tsArg(args,0);if typeValue.kind!=tsStringKind{tsPropertyFailure("The type argument must be a string")};tsSetAsyncContextTrackingEnabled(true);options:=tsArg(args,1);trigger:=tsFSProperty(options,"triggerAsyncId");if tsIsUndefined(trigger)&&tsIsNumeric(options){trigger=options};triggerID:=int64(0);if !tsIsUndefined(trigger){triggerID=int64(tsNumber(trigger))};object:=tsNewObject();state:=&tsNodeAsyncResourceState{object:object,context:loop.asyncContext,triggerID:triggerID,typeName:tsText(typeValue)};tsNodeAsyncHooksSingleton.resources[unsafe.Pointer(object)]=state;object.set("type",typeValue);for _,method:=range []string{"emitBefore","emitAfter","asyncId","triggerAsyncId","emitDestroy","runInAsyncScope","bind"}{object.set(method,tsNodeAsyncResourceMethod(method))};return tsObjectValue(object)}
	static:=tsNewObject();static.set("bind",tsNamedNativeMethod("bind",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{callback:=tsArg(args,0);typeValue:=tsArg(args,1);typeName:="";if !tsNullish(typeValue){typeName=tsText(typeValue)};if typeName==""&&callback.kind==tsFunctionKind{typeName=(*tsFunction)(callback.ref).name};if typeName==""{typeName="bound-anonymous-fn"};resourceValue:=tsConstruct(tsFunctionValue(fn),tsStringReference(tsStringUTF8(typeName)));resource:=tsFSProperty(resourceValue,"bind");return tsCallReceiver(loop,resource,resourceValue,callback,tsArg(args,2))}))
	fn.properties=static;return tsFunctionValue(fn)
}

func tsNodeAsyncHooksCreateHook() tsValue {
	return tsNamedNativeMethod("createHook", 1, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
		hook := tsArg(args, 0)
		if hook.kind != tsObjectKind && hook.kind != tsInstanceKind { tsPropertyFailure("The hooks argument must be an object") }
		for _, name := range []string{"init", "before", "after", "destroy", "promiseResolve"} {
			value := tsFSProperty(hook, name)
			if !tsIsUndefined(value) && value.kind != tsFunctionKind { tsPropertyFailure("Hook callback must be a function") }
		}
		init := tsFSProperty(hook, "init")
		var enabledInit tsValue
		enabled := false
		out := tsNewObject()
		out.set("enable", tsNamedNativeMethod("enable", 0, func(_ *tsLoop, receiver tsValue, _ ...tsValue) tsValue {
			if init.kind == tsFunctionKind && tsIsUndefined(enabledInit) {
				callback := init
				enabledInit = tsFunctionValue(tsFunc(func(loop *tsLoop, values ...tsValue) tsValue { return tsCall(loop, callback, values...) }))
				module := tsInternalAsyncHooksTickModule(loop)
				list := (*tsArray)(tsGet(loop, module, tsStringReference(tsStringUTF8("tickInitHooks"))).ref)
				list.values = append(list.values, enabledInit)
			}
			if !enabled { enabled = true; tsInternalAsyncHooksMark(true) }
			return receiver
		}))
		out.set("disable", tsNamedNativeMethod("disable", 0, func(loop *tsLoop, receiver tsValue, _ ...tsValue) tsValue {
			if !tsIsUndefined(enabledInit) {
				module := tsInternalAsyncHooksTickModule(loop)
				list := (*tsArray)(tsGet(loop, module, tsStringReference(tsStringUTF8("tickInitHooks"))).ref)
				for index, entry := range list.values { if tsStrictEqual(entry, enabledInit) { list.values = append(list.values[:index], list.values[index+1:]...); break } }
				enabledInit = tsU
			}
			if enabled { enabled = false; tsInternalAsyncHooksMark(false) }
			return receiver
		}))
		return tsObjectValue(out)
	})
}

func tsNodeAsyncHooksModule() tsValue {
	tsNodeAsyncHooksInit();if len(tsNodeAsyncHooksSingleton.module.values)>0{return tsObjectValue(tsNodeAsyncHooksSingleton.module)}
	module:=tsNodeAsyncHooksSingleton.module;module.set("AsyncLocalStorage",tsNodeAsyncLocalStorageConstructor());module.set("AsyncResource",tsNodeAsyncResourceConstructor());module.set("createHook",tsNodeAsyncHooksCreateHook())
	module.set("executionAsyncId",tsNamedNativeMethod("executionAsyncId",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNumberValue(0)}));module.set("triggerAsyncId",tsNamedNativeMethod("triggerAsyncId",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNumberValue(0)}));module.set("executionAsyncResource",tsNamedNativeMethod("executionAsyncResource",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsU}))
	providers:=tsNewObject();for index,name:=range []string{"NONE","DIRHANDLE","DNSCHANNEL","ELDHISTOGRAM","FILEHANDLE","FILEHANDLECLOSEREQ","FIXEDSIZEBLOBCOPY","FSEVENTWRAP","FSREQCALLBACK","FSREQPROMISE","GETADDRINFOREQWRAP","GETNAMEINFOREQWRAP","HEAPSNAPSHOT","HTTP2SESSION","HTTP2STREAM","HTTP2PING","HTTP2SETTINGS","HTTPINCOMINGMESSAGE","HTTPCLIENTREQUEST","JSSTREAM","JSUDPWRAP","MESSAGEPORT","PIPECONNECTWRAP","PIPESERVERWRAP","PIPEWRAP","PROCESSWRAP","PROMISE","QUERYWRAP","SHUTDOWNWRAP","SIGNALWRAP","STATWATCHER","STREAMPIPE","TCPCONNECTWRAP","TCPSERVERWRAP","TCPWRAP","TTYWRAP","UDPSENDWRAP","UDPWRAP","SIGINTWATCHDOG","WORKER","WORKERHEAPSNAPSHOT","WRITEWRAP","ZLIB","CHECKPRIMEREQUEST","PBKDF2REQUEST","KEYPAIRGENREQUEST","KEYGENREQUEST","KEYEXPORTREQUEST","CIPHERREQUEST","DERIVEBITSREQUEST","HASHREQUEST","RANDOMBYTESREQUEST","RANDOMPRIMEREQUEST","SCRYPTREQUEST","SIGNREQUEST","TLSWRAP","VERIFYREQUEST","INSPECTORJSBINDING"}{providers.set(name,tsNumberValue(float64(index)))};module.set("asyncWrapProviders",tsObjectValue(providers));return tsObjectValue(module)
}
`
