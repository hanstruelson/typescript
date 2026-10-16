package goemit

// internalAsyncHooksTickRuntime implements the stable TickObject init-hook
// list shared by node:async_hooks and process.nextTick.
const internalAsyncHooksTickRuntime = `
var tsInternalAsyncHooksTickState struct{sync.Once;module *tsObject;hooks *tsArray;nextAsyncID uint64}
func tsInternalAsyncHooksTickModule(loop *tsLoop)tsValue{
	tsInternalAsyncHooksTickState.Do(func(){hooks:=&tsArray{values:[]tsValue{}};module:=tsNewObject();tsInternalAsyncHooksTickState.hooks=hooks;tsInternalAsyncHooksTickState.nextAsyncID=1;module.set("tickInitHooks",tsArrayValue(hooks));module.set("newAsyncId",tsNamedNativeMethod("newAsyncId",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{tsInternalAsyncHooksTickState.nextAsyncID++;return tsNumberValue(float64(tsInternalAsyncHooksTickState.nextAsyncID))}));module.set("default",tsObjectValue(module));tsInternalAsyncHooksTickState.module=module})
	return tsObjectValue(tsInternalAsyncHooksTickState.module)
}
func tsNodeProcessNextTick(loop *tsLoop,process *tsObject,args []tsValue)tsValue{
	callback:=tsArg(args,0);if callback.kind!=tsFunctionKind{tsFSInvalid("callback","a function")};if tsTruthy(process.get("_exiting")){return tsU};extra:=append([]tsValue(nil),args[1:]...)
	module:=tsInternalAsyncHooksTickModule(loop);hooksValue:=tsGet(loop,module,tsStringReference(tsStringUTF8("tickInitHooks")));hooks:=(*tsArray)(hooksValue.ref)
	if hooks.length()>0{
		newID:=tsGet(loop,module,tsStringReference(tsStringUTF8("newAsyncId")));asyncID:=tsCallReceiver(loop,newID,module);resource:=tsNewObject();resource.set("callback",callback);if len(extra)>0{resource.set("args",tsArrayValue(&tsArray{values:append([]tsValue(nil),extra...)}))}else{resource.set("args",tsU)};resource.set("frame",tsInternalAsyncContextFrameValue(loop.asyncContext));snapshot:=append([]tsValue(nil),hooks.values...)
		for _,hook:=range snapshot{if hook.kind!=tsFunctionKind{continue};failed:=false;var failure tsValue;func(){defer func(){if caught:=recover();caught!=nil{failed=true;failure=tsUnwrap(caught)}}();tsCall(loop,hook,asyncID,tsStringReference(tsStringUTF8("TickObject")),tsNumberValue(0),tsObjectValue(resource))}();if failed{console:=tsNodeConsoleModule();errorMethod:=tsGet(loop,console,tsStringReference(tsStringUTF8("error")));if errorMethod.kind==tsFunctionKind{printValue:=failure;func(){defer func(){_ = recover()}();stack:=tsGet(loop,failure,tsStringReference(tsStringUTF8("stack")));if stack.kind==tsStringKind{printValue=stack}}();func(){defer func(){_ = recover()}();tsCallReceiver(loop,errorMethod,console,printValue)}()};return tsNodeProcessExit(loop,process,tsNumberValue(1))}}
	}
	context:=loop.asyncContext;loop.ready=append(loop.ready,func(){tsWithAsyncContext(loop,context,func(){tsCallReceiver(loop,callback,tsU,extra...)})});return tsU
}
`
