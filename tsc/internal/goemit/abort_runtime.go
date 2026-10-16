package goemit

// Cancellation is a native shared channel. Host I/O reads the locked reason;
// application event listeners remain on the signal's creating event loop.
const abortRuntime = `
type tsAbortState struct {
    sync.Mutex
    done chan struct{}
    reason tsValue
    aborted bool
    events *tsFSEvents
    followers []func(tsValue)
}
func tsAbortReason(state *tsAbortState)(tsValue,bool){state.Lock();defer state.Unlock();return state.reason,state.aborted}
func tsAbortDefault()tsValue{return tsErrorValue(&tsRuntimeError{name:"AbortError",message:"This operation was aborted",fields:map[string]tsValue{"code":tsNumberValue(20)}})}
func tsAbortController(loop *tsLoop)tsValue{
    state:=&tsAbortState{done:make(chan struct{}),events:tsFSEventObject(loop)}
    tsEventTargetInitialize();signal:=state.events.object;signal.nativeClass=tsNativeClass("AbortSignal");signal.prototype=tsObjectValue(tsEventTargetIntrinsics.eventTargetPrototype);signal.eventTarget=&tsEventTargetState{object:signal,loop:loop,listeners:map[string][]*tsEventTargetListener{}};signal.abort=state;signal.set("aborted",tsBooleanValue(false));signal.set("reason",tsU);signal.set("onabort",tsNull)
    signal.set("throwIfAborted",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{reason,aborted:=tsAbortReason(state);if aborted{panic(tsThrown{reason})};return tsU})))
    signal.set("addEventListener",tsNamedNativeMethod("addEventListener",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{if len(args)<2{tsPropertyFailure("EventTarget.addEventListener requires 2 arguments")};tsEventTargetAdd(loop,receiver,tsArg(args,0),tsArg(args,1),tsArg(args,2));return tsU}))
    signal.set("removeEventListener",tsNamedNativeMethod("removeEventListener",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{if len(args)<2{tsPropertyFailure("EventTarget.removeEventListener requires 2 arguments")};tsEventTargetRemove(loop,receiver,tsArg(args,0),tsArg(args,1),tsArg(args,2));return tsU}))
    controller:=tsNewObject();controller.nativeClass=tsNativeClass("AbortController");controller.set("signal",tsObjectValue(signal));controller.set("abort",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{reason:=tsArg(args,0);if tsIsUndefined(reason){reason=tsAbortDefault()};tsAbortTrigger(state,reason);return tsU})))
    return tsObjectValue(controller)
}
func tsAbortTrigger(state *tsAbortState,reason tsValue){
    state.Lock();if state.aborted{state.Unlock();return};state.aborted=true;state.reason=reason;followers:=state.followers;state.followers=nil;close(state.done);state.Unlock()
    deliver:=func(){loop:=state.events.loop;signal:=tsObjectValue(state.events.object);state.events.object.set("aborted",tsBooleanValue(true));state.events.object.set("reason",reason);eventValue:=tsConstruct(tsEventConstructor(),tsStringReference(tsStringUTF8("abort")));tsEventTargetDispatch(loop,signal,eventValue);handler:=state.events.object.get("onabort");if handler.kind==tsFunctionKind{tsCall(loop,handler,eventValue)}}
    // abort() dispatch is synchronous on the signal's loop. Native cancellation
    // is published before handlers run, matching AbortSignal's ordering.
    deliver();for _,follow:=range followers{follow(reason)}
}
func tsFSAbortSignal(options tsValue)*tsAbortState{
    value:=tsFSProperty(options,"signal");if tsIsUndefined(value){return nil}
    if value.kind!=tsObjectKind||(*tsObject)(value.ref).abort==nil{tsFSInvalid("options.signal","an AbortSignal")}
    return (*tsObject)(value.ref).abort
}
func tsFSAbortError(signal *tsAbortState)tsValue{
    reason,_:=tsAbortReason(signal);value:=tsFSException("AbortError","ABORT_ERR","The operation was aborted");(*tsRuntimeError)(value.ref).fields["cause"]=reason;return value
}
func tsFSCheckAbort(signal *tsAbortState){if signal==nil{return};_,aborted:=tsAbortReason(signal);if aborted{panic(tsThrown{tsFSAbortError(signal)})}}
func tsAbortSignalModule(loop *tsLoop)tsValue{
    out:=tsNewObject()
    out.set("abort",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{controller:=tsAbortController(loop);tsCall(loop,tsFSProperty(controller,"abort"),args...);return tsFSProperty(controller,"signal")})))
    out.set("timeout",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{
        milliseconds:=tsFSInteger(tsArg(args,0),"delay",0,0);if milliseconds>4294967295{panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","delay is out of range")})};controller:=tsAbortController(loop);signal:=tsFSProperty(controller,"signal");state:=(*tsObject)(signal.ref).abort
        cancel:=make(chan struct{});var once sync.Once;stop:=func(){once.Do(func(){close(cancel)})}
        task:=loop.submit(func()tsResult{timer:=time.NewTimer(time.Duration(milliseconds)*time.Millisecond);defer timer.Stop();select{case <-timer.C:return tsResult{tsBooleanValue(true),false};case <-cancel:return tsResult{tsBooleanValue(false),false}}},func(result tsResult){if tsTruthy(result.value){reason:=tsErrorValue(&tsRuntimeError{name:"TimeoutError",message:"The operation was aborted due to timeout",fields:map[string]tsValue{"code":tsNumberValue(23)}});tsAbortTrigger(state,reason)}});task.unref=true;loop.resourceClosers=append(loop.resourceClosers,stop);return signal
    })))
    return tsObjectValue(out)
}
`
