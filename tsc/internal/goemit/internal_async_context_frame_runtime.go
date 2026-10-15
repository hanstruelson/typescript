package goemit

// internalAsyncContextFrameRuntime maps Bun's JS shim onto loop-owned async
// context frames. ModuleGraph frame queries are unavailable in this runtime.
const internalAsyncContextFrameRuntime = `
var tsInternalAsyncContextFrameState struct {
    sync.Mutex
    module *tsObject
    values map[*tsAsyncContextFrame]tsValue
    frames map[*tsObject]*tsAsyncContextFrame
}

func tsInternalAsyncContextFrameValue(frame *tsAsyncContextFrame) tsValue {
    if frame==nil{return tsU}
    tsInternalAsyncContextFrameState.Lock();defer tsInternalAsyncContextFrameState.Unlock()
    if value,ok:=tsInternalAsyncContextFrameState.values[frame];ok{return value}
    object:=tsNewObject();object.prototype=tsNull;object.nonExtensible=true
    value:=tsObjectValue(object);tsInternalAsyncContextFrameState.values[frame]=value;tsInternalAsyncContextFrameState.frames[object]=frame
    return value
}

func tsInternalAsyncContextFramePointer(value tsValue) *tsAsyncContextFrame {
    if value.kind!=tsObjectKind{return nil}
    tsInternalAsyncContextFrameState.Lock();defer tsInternalAsyncContextFrameState.Unlock()
    return tsInternalAsyncContextFrameState.frames[(*tsObject)(value.ref)]
}

func tsInternalAsyncContextFrameModule(loop *tsLoop) tsValue {
    tsInternalAsyncContextFrameState.Lock()
    if tsInternalAsyncContextFrameState.values==nil{tsInternalAsyncContextFrameState.values=map[*tsAsyncContextFrame]tsValue{};tsInternalAsyncContextFrameState.frames=map[*tsObject]*tsAsyncContextFrame{}}
    if tsInternalAsyncContextFrameState.module!=nil{module:=tsInternalAsyncContextFrameState.module;tsInternalAsyncContextFrameState.Unlock();return tsObjectValue(module)}
    out:=tsNewObject();tsInternalAsyncContextFrameState.module=out;tsInternalAsyncContextFrameState.Unlock()
    out.set("enabled",tsBooleanValue(true))
    out.set("current",tsNamedNativeMethod("current",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsInternalAsyncContextFrameValue(loop.asyncContext)}))
    out.set("currentGraphFrame",tsNamedNativeMethod("currentGraphFrame",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsU}))
    out.set("currentGraph",tsNamedNativeMethod("currentGraph",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsU}))
    out.set("graphOf",tsNamedNativeMethod("graphOf",1,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsU}))
    out.set("exchange",tsNamedNativeMethod("exchange",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{previous:=tsInternalAsyncContextFrameValue(loop.asyncContext);loop.asyncContext=tsInternalAsyncContextFramePointer(tsArg(args,0));return previous}))
    out.set("run",tsNamedNativeMethod("run",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{
        frame:=tsInternalAsyncContextFramePointer(tsArg(args,0));callback:=tsArg(args,1);thisArg:=tsArg(args,2);if len(args)<3{thisArg=tsU}
        if callback.kind!=tsFunctionKind{tsPropertyFailure("AsyncContextFrame.run requires a function")}
        callArgs:=[]tsValue{};if len(args)>3{callArgs=args[3:]}
        if frame==loop.asyncContext{return tsCallReceiver(loop,callback,thisArg,callArgs...)}
        var result tsValue;tsWithAsyncContext(loop,frame,func(){result=tsCallReceiver(loop,callback,thisArg,callArgs...)});return result
    }))
    out.set("default",tsObjectValue(out))
    return tsObjectValue(out)
}
`
