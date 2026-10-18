package goemit

// nodeAsyncHooksRuntime adapts AsyncLocalStorage and AsyncResource to the
// event loop's persistent async-context frames. TickObject init hooks run
// during process.nextTick setup; other hook phases remain unsupported.
const nodeAsyncHooksRuntime = `
type tsNodeALSState struct { object *tsObject; defaultValue tsValue; name string }
type tsNodeAsyncResourceState struct { object *tsObject; context *tsAsyncContextFrame; triggerID int64; typeName string }
type tsNodeRunScopeState struct { storage *tsNodeALSState; previous tsValue; disposed bool }
var tsNodeAsyncHooksSingleton struct { sync.Once; module *tsObject; storages map[unsafe.Pointer]*tsNodeALSState; resources map[unsafe.Pointer]*tsNodeAsyncResourceState; scopes map[unsafe.Pointer]*tsNodeRunScopeState; frameMutations uint64; hasEnabledCreateHook, createHookNotImplementedWarned, asyncIDWarningIssued, executionResourceWarningIssued bool }

func tsNodeAsyncHooksInit() { tsNodeAsyncHooksSingleton.Do(func(){tsNodeAsyncHooksSingleton.module=tsNewObject();tsNodeAsyncHooksSingleton.storages=map[unsafe.Pointer]*tsNodeALSState{};tsNodeAsyncHooksSingleton.resources=map[unsafe.Pointer]*tsNodeAsyncResourceState{};tsNodeAsyncHooksSingleton.scopes=map[unsafe.Pointer]*tsNodeRunScopeState{}}) }
func tsNodeAsyncHooksWarning(loop *tsLoop,message string,warned *bool){if *warned{return};*warned=true;console:=tsNodeConsoleModule();warn:=tsGet(loop,console,tsStringReference(tsStringUTF8("warn")));if warn.kind==tsFunctionKind{tsCallReceiver(loop,warn,console,tsStringReference(tsStringUTF8("[bun] Warning:")),tsStringReference(tsStringUTF8(message)))}}

func tsNodeALSStateOf(receiver tsValue) *tsNodeALSState { if receiver.kind!=tsObjectKind {tsPropertyFailure("AsyncLocalStorage method called on incompatible receiver")};state:=tsNodeAsyncHooksSingleton.storages[receiver.ref];if state==nil{tsPropertyFailure("AsyncLocalStorage method called on incompatible receiver")};return state }
func tsNodeALSStore(loop *tsLoop,state *tsNodeALSState) tsValue { key:=unsafe.Pointer(state.object);context:=loop.asyncContext;if context==nil{return state.defaultValue};if context.masked!=nil&&context.masked[key]{return state.defaultValue};for frame:=context;frame!=nil;frame=frame.prev{if frame.storage==key{return frame.value}};return state.defaultValue }
func tsNodeALSFind(frame *tsAsyncContextFrame,key unsafe.Pointer)*tsAsyncContextFrame{for current:=frame;current!=nil;current=current.prev{if current.storage==key{return current}};return nil}
func tsNodeALSMasked(frame *tsAsyncContextFrame,key unsafe.Pointer)bool{return frame!=nil&&frame.masked!=nil&&frame.masked[key]}
func tsNodeALSMasksEqual(left,right map[unsafe.Pointer]bool)bool{if (left==nil)!=(right==nil)||len(left)!=len(right){return false};for key,value:=range left{if right[key]!=value{return false}};return true}
func tsNodeALSCopyUntil(from,stop,tail *tsAsyncContextFrame)*tsAsyncContextFrame{if from==stop{if from==nil||tail==nil||tsNodeALSMasksEqual(tail.masked,from.masked){return tail};return &tsAsyncContextFrame{storage:tail.storage,value:tail.value,prev:tail.prev,masked:from.masked}};copied:=[]*tsAsyncContextFrame{};for current:=from;current!=nil&&current!=stop;current=current.prev{copied=append(copied,current)};if stop!=nil&&tail!=nil&&!tsNodeALSMasksEqual(tail.masked,stop.masked){tail=&tsAsyncContextFrame{storage:tail.storage,value:tail.value,prev:tail.prev,masked:stop.masked}};for index:=len(copied)-1;index>=0;index--{frame:=copied[index];tail=&tsAsyncContextFrame{storage:frame.storage,value:frame.value,prev:tail,masked:frame.masked}};return tail}
func tsNodeALSWithout(frame *tsAsyncContextFrame,key unsafe.Pointer)*tsAsyncContextFrame{found:=tsNodeALSFind(frame,key);if found==nil{return frame};return tsNodeALSCopyUntil(frame,found,found.prev)}
func tsNodeALSPush(loop *tsLoop,state *tsNodeALSState,value tsValue)*tsAsyncContextFrame{key:=unsafe.Pointer(state.object);var masks map[unsafe.Pointer]bool;if head:=loop.asyncContext;head!=nil{masks=head.masked;if masks!=nil&&masks[key]{filtered:=make(map[unsafe.Pointer]bool,len(masks)-1);for item,masked:=range masks{if masked&&item!=key{filtered[item]=true}};if len(filtered)==0{masks=nil}else{masks=filtered}}};frame:=&tsAsyncContextFrame{storage:key,value:value,prev:loop.asyncContext,masked:masks};loop.asyncContext=frame;return frame}
func tsNodeALSMethod(name string) tsValue { length:=0;switch name{case "enterWith","withScope":length=1;case "run":length=2;case "exit":length=1;case "_propagate":length=3};return tsNamedNativeMethod(name,length,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeALSStateOf(receiver);switch name{
case "getStore":return tsNodeALSStore(loop,state)
case "enterWith":tsNodeALSPush(loop,state,tsArg(args,0));return tsU
case "run":callback:=tsArg(args,1);if callback.kind!=tsFunctionKind&&callback.kind!=tsClassKind{tsPropertyFailure("The callback argument must be of type function")};store:=tsArg(args,0);key:=unsafe.Pointer(state.object);prior:=loop.asyncContext;bound:=tsNodeALSFind(prior,key);beforeValue:=state.defaultValue;if bound!=nil&&!tsNodeALSMasked(prior,key){beforeValue=bound.value};if tsSameValue(beforeValue,store){return tsCallReceiver(loop,callback,tsU,args[2:]...)};mutations:=tsNodeAsyncHooksSingleton.frameMutations;head:=prior;if bound!=nil{head=tsNodeALSCopyUntil(prior,bound,bound.prev)};loop.asyncContext=head;frame:=tsNodeALSPush(loop,state,store);defer func(){current:=loop.asyncContext;if mutations==tsNodeAsyncHooksSingleton.frameMutations&&(current==frame||(current!=nil&&current.prev==frame.prev&&current.storage==key&&tsNodeALSMasksEqual(current.masked,frame.masked))){if prior==nil{loop.asyncContext=nil}else{loop.asyncContext=&tsAsyncContextFrame{storage:prior.storage,value:prior.value,prev:prior.prev,masked:prior.masked}}}else{loop.asyncContext=tsNodeALSWithout(current,key);tsNodeALSPush(loop,state,beforeValue)}}();return tsCallReceiver(loop,callback,tsU,args[2:]...)
case "exit":run:=tsNodeALSMethod("run");callArgs:=append([]tsValue{tsU,tsArg(args,0)},args[1:]...);return tsCallReceiver(loop,run,receiver,callArgs...)
case "disable":key:=unsafe.Pointer(state.object);if loop.asyncContext!=nil&&!tsNodeALSMasked(loop.asyncContext,key)&&tsNodeALSFind(loop.asyncContext,key)!=nil{masks:=make(map[unsafe.Pointer]bool,len(loop.asyncContext.masked)+1);for item,masked:=range loop.asyncContext.masked{if masked{masks[item]=true}};masks[key]=true;loop.asyncContext.masked=masks;tsNodeAsyncHooksSingleton.frameMutations++};return tsU
case "_enable","_propagate":return tsU
case "withScope":scope:=tsNewObject();scopeState:=&tsNodeRunScopeState{storage:state,previous:tsNodeALSStore(loop,state)};tsNodeAsyncHooksSingleton.scopes[unsafe.Pointer(scope)]=scopeState;scope.set("dispose",tsNamedNativeMethod("dispose",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{if item:=tsNodeAsyncHooksSingleton.scopes[receiver.ref];item!=nil&&!item.disposed{item.disposed=true;tsNodeALSPush(loop,item.storage,tsValue(item.previous))};return tsU}));dispose:=scope.get("dispose");scope.set(tsPropertyKey(tsWellKnownSymbol("dispose")),dispose);tsNodeALSPush(loop,state,tsArg(args,0));return tsObjectValue(scope)
};return tsU}) }

func tsNodeAsyncResourceStateOf(receiver tsValue) *tsNodeAsyncResourceState {if receiver.kind!=tsObjectKind{tsPropertyFailure("AsyncResource method called on incompatible receiver")};state:=tsNodeAsyncHooksSingleton.resources[receiver.ref];if state==nil{tsPropertyFailure("AsyncResource method called on incompatible receiver")};return state}
func tsNodeAsyncResourceMethod(name string) tsValue {length:=0;if name=="runInAsyncScope"||name=="bind"{length=2};return tsNamedNativeMethod(name,length,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeAsyncResourceStateOf(receiver);switch name{
case "emitBefore","emitAfter":return tsBooleanValue(true)
case "asyncId":return tsNumberValue(0)
case "triggerAsyncId":return tsNumberValue(float64(state.triggerID))
case "emitDestroy":return tsU
case "runInAsyncScope":callback:=tsArg(args,0);if callback.kind!=tsFunctionKind&&callback.kind!=tsClassKind{tsPropertyFailure("The fn argument must be a function")};thisArg:=tsArg(args,1);var result tsValue;tsWithAsyncContext(loop,state.context,func(){result=tsCallReceiver(loop,callback,thisArg,args[2:]...)});return result
case "bind":callback:=tsArg(args,0);if callback.kind!=tsFunctionKind&&callback.kind!=tsClassKind{tsInternalValidatorInvalidType("fn","a function",callback)};thisArg:=tsArg(args,1);invoke:=func(loop *tsLoop,receiver tsValue,values []tsValue)tsValue{effectiveThis:=thisArg;if tsIsUndefined(effectiveThis){effectiveThis=receiver};var result tsValue;tsWithAsyncContext(loop,state.context,func(){result=tsCallReceiver(loop,callback,effectiveThis,values...)});return result};bound:=tsFunctionValue(tsFunc(func(loop *tsLoop,values ...tsValue)tsValue{return invoke(loop,tsU,values)}));boundFunction:=(*tsFunction)(bound.ref);boundFunction.receiverCall=func(loop *tsLoop,receiver tsValue,values ...tsValue)tsValue{return invoke(loop,receiver,values)};boundFunction.name=(*tsFunction)(callback.ref).name;boundFunction.length=(*tsFunction)(callback.ref).length;return bound
};return tsU})}

func tsNodeAsyncLocalStorageConstructor() tsValue {
	fn:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{tsPropertyFailure("Class constructor AsyncLocalStorage cannot be invoked without new");return tsU});fn.name="AsyncLocalStorage";fn.length=0;fn.constructible=true
	fn.construct=func(loop *tsLoop,args ...tsValue)tsValue{options:=tsArg(args,0);defaultValue:=tsU;nameValue:=tsU;if !tsIsUndefined(options){tsInternalValidatorObject(options,tsStringReference(tsStringUTF8("options")),tsU);defaultValue=tsFSProperty(options,"defaultValue");nameValue=tsFSProperty(options,"name")};name:="";if !tsIsUndefined(nameValue){name=tsText(nameValue)};tsSetAsyncContextTrackingEnabled(true);object:=tsNewObject();state:=&tsNodeALSState{object:object,defaultValue:defaultValue,name:name};tsNodeAsyncHooksSingleton.storages[unsafe.Pointer(object)]=state
		for _,method:=range []string{"getStore","enterWith","run","exit","disable","withScope","_enable","_propagate"}{object.set(method,tsNodeALSMethod(method))};object.set("name",tsStringReference(tsStringUTF8(name)));return tsObjectValue(object)}
	static:=tsNewObject();static.set("bind",tsNamedNativeMethod("bind",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{callback:=tsArg(args,0);if callback.kind!=tsFunctionKind&&callback.kind!=tsClassKind{tsInternalValidatorInvalidType("fn","a function",callback)};context:=loop.asyncContext;boundArgs:=append([]tsValue(nil),args[1:]...);invoke:=func(loop *tsLoop,values []tsValue)tsValue{combined:=append(append([]tsValue(nil),boundArgs...),values...);var result tsValue;tsWithAsyncContext(loop,context,func(){result=tsCallReceiver(loop,callback,tsU,combined...)});return result};bound:=tsFunctionValue(tsFunc(func(loop *tsLoop,values ...tsValue)tsValue{return invoke(loop,values)}));boundFunction:=(*tsFunction)(bound.ref);boundFunction.receiverCall=func(loop *tsLoop,_ tsValue,values ...tsValue)tsValue{return invoke(loop,values)};boundFunction.name="bound ";boundFunction.length=0;return bound}))
	static.set("snapshot",tsNamedNativeMethod("snapshot",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{context:=loop.asyncContext;snapshot:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{callback:=tsArg(args,0);if callback.kind!=tsFunctionKind&&callback.kind!=tsClassKind{tsInternalValidatorInvalidType("fn","a function",callback)};var result tsValue;tsWithAsyncContext(loop,context,func(){result=tsCallReceiver(loop,callback,tsU,args[1:]...)});return result}));(*tsFunction)(snapshot.ref).length=1;return snapshot}))
	fn.properties=static
	return tsFunctionValue(fn)
}

func tsNodeAsyncResourceConstructor() tsValue {
	fn:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{tsPropertyFailure("Class constructor AsyncResource cannot be invoked without new");return tsU});fn.name="AsyncResource";fn.length=1;fn.constructible=true
fn.construct=func(loop *tsLoop,args ...tsValue)tsValue{typeValue:=tsArg(args,0);if typeValue.kind!=tsStringKind{tsInternalValidatorInvalidType("type","a string",typeValue)};options:=tsArg(args,1);trigger:=options;if !tsIsNumeric(options){trigger=tsFSProperty(options,"triggerAsyncId")};if tsIsUndefined(trigger){trigger=tsNumberValue(0)};triggerNumber:=tsNumber(trigger);if !tsIsNumeric(trigger)||math.IsNaN(triggerNumber)||math.IsInf(triggerNumber,0)||math.Trunc(triggerNumber)!=triggerNumber||triggerNumber< -1||triggerNumber>9007199254740991{panic(tsThrown{tsFSException("RangeError","ERR_INVALID_ASYNC_ID","Invalid triggerAsyncId value: "+tsText(trigger))})};if tsNodeAsyncHooksSingleton.hasEnabledCreateHook&&tsText(typeValue)==""{panic(tsThrown{tsFSException("TypeError","ERR_ASYNC_TYPE","Invalid name for async \"type\": ")})};tsSetAsyncContextTrackingEnabled(true);triggerID:=int64(triggerNumber);object:=tsNewObject();state:=&tsNodeAsyncResourceState{object:object,context:loop.asyncContext,triggerID:triggerID,typeName:tsText(typeValue)};tsNodeAsyncHooksSingleton.resources[unsafe.Pointer(object)]=state;object.set("type",typeValue);for _,method:=range []string{"emitBefore","emitAfter","asyncId","triggerAsyncId","emitDestroy","runInAsyncScope","bind"}{object.set(method,tsNodeAsyncResourceMethod(method))};return tsObjectValue(object)}
	static:=tsNewObject();static.set("bind",tsNamedNativeMethod("bind",3,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{callback:=tsArg(args,0);typeValue:=tsArg(args,1);typeName:="";if tsTruthy(typeValue){typeName=tsText(typeValue)};if typeName==""&&callback.kind==tsFunctionKind{typeName=(*tsFunction)(callback.ref).name};if typeName==""{typeName="bound-anonymous-fn"};resourceValue:=tsConstruct(tsFunctionValue(fn),tsStringReference(tsStringUTF8(typeName)));resource:=tsFSProperty(resourceValue,"bind");return tsCallReceiver(loop,resource,resourceValue,callback,tsArg(args,2))}))
	fn.properties=static;return tsFunctionValue(fn)
}

func tsNodeAsyncHooksCreateHook() tsValue {
	return tsNamedNativeMethod("createHook", 1, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
		hook := tsArg(args, 0)
		tsInternalValidatorObject(hook,tsStringReference(tsStringUTF8("hook")),tsU)
		hookCallbacks:=map[string]tsValue{}
		for _, name := range []string{"init", "before", "after", "destroy", "promiseResolve"} {
			value := tsFSProperty(hook, name)
			if !tsIsUndefined(value) && value.kind != tsFunctionKind && value.kind != tsClassKind { panic(tsThrown{tsFSException("TypeError","ERR_ASYNC_CALLBACK","hook."+name+" must be a function")}) }
			hookCallbacks[name]=value
		}
		init := hookCallbacks["init"]
		unsupportedHooks:=false;for _,name:=range []string{"before","after","destroy","promiseResolve"}{if !tsIsUndefined(hookCallbacks[name]){unsupportedHooks=true}}
		var enabledInit tsValue
		enabled := false
		out := tsNewObject()
		out.set("enable", tsNamedNativeMethod("enable", 0, func(loop *tsLoop, receiver tsValue, _ ...tsValue) tsValue {
			if init.kind == tsFunctionKind && tsIsUndefined(enabledInit) {
				callback := init
				enabledInit = tsFunctionValue(tsFunc(func(loop *tsLoop, values ...tsValue) tsValue { return tsCall(loop, callback, values...) }))
				module := tsInternalAsyncHooksTickModule(loop)
				list := (*tsArray)(tsGet(loop, module, tsStringReference(tsStringUTF8("tickInitHooks"))).ref)
				list.values = append(list.values, enabledInit)
			}
			if unsupportedHooks{tsNodeAsyncHooksWarning(loop,"async_hooks.createHook is not implemented in Bun. Hooks can still be created but will never be called.",&tsNodeAsyncHooksSingleton.createHookNotImplementedWarned)}
			tsNodeAsyncHooksSingleton.hasEnabledCreateHook=true
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
	module.set("executionAsyncId",tsNamedNativeMethod("executionAsyncId",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{tsNodeAsyncHooksWarning(loop,"async_hooks.executionAsyncId/triggerAsyncId are not implemented in Bun. It will return 0 every time.",&tsNodeAsyncHooksSingleton.asyncIDWarningIssued);return tsNumberValue(0)}));module.set("triggerAsyncId",tsNamedNativeMethod("triggerAsyncId",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNumberValue(0)}));module.set("executionAsyncResource",tsNamedNativeMethod("executionAsyncResource",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{tsNodeAsyncHooksWarning(loop,"async_hooks.executionAsyncResource is not implemented in Bun. It returns a reference to process.stdin every time.",&tsNodeAsyncHooksSingleton.executionResourceWarningIssued);process:=tsNodeProcessModule(loop);return tsGet(loop,process,tsStringReference(tsStringUTF8("stdin")))}))
	providers:=tsNewObject();for index,name:=range []string{"NONE","DIRHANDLE","DNSCHANNEL","ELDHISTOGRAM","FILEHANDLE","FILEHANDLECLOSEREQ","FIXEDSIZEBLOBCOPY","FSEVENTWRAP","FSREQCALLBACK","FSREQPROMISE","GETADDRINFOREQWRAP","GETNAMEINFOREQWRAP","HEAPSNAPSHOT","HTTP2SESSION","HTTP2STREAM","HTTP2PING","HTTP2SETTINGS","HTTPINCOMINGMESSAGE","HTTPCLIENTREQUEST","JSSTREAM","JSUDPWRAP","MESSAGEPORT","PIPECONNECTWRAP","PIPESERVERWRAP","PIPEWRAP","PROCESSWRAP","PROMISE","QUERYWRAP","SHUTDOWNWRAP","SIGNALWRAP","STATWATCHER","STREAMPIPE","TCPCONNECTWRAP","TCPSERVERWRAP","TCPWRAP","TTYWRAP","UDPSENDWRAP","UDPWRAP","SIGINTWATCHDOG","WORKER","WORKERHEAPSNAPSHOT","WRITEWRAP","ZLIB","CHECKPRIMEREQUEST","PBKDF2REQUEST","KEYPAIRGENREQUEST","KEYGENREQUEST","KEYEXPORTREQUEST","CIPHERREQUEST","DERIVEBITSREQUEST","HASHREQUEST","RANDOMBYTESREQUEST","RANDOMPRIMEREQUEST","SCRYPTREQUEST","SIGNREQUEST","TLSWRAP","VERIFYREQUEST","INSPECTORJSBINDING"}{providers.set(name,tsNumberValue(float64(index)))};module.set("asyncWrapProviders",tsObjectValue(providers));return tsObjectValue(module)
}
`
