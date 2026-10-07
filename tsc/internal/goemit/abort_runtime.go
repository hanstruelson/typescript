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
    signal:=state.events.object;signal.nativeClass=tsNativeClass("AbortSignal");signal.abort=state;signal.set("aborted",tsBooleanValue(false));signal.set("reason",tsU);signal.set("onabort",tsNull)
    signal.set("throwIfAborted",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{reason,aborted:=tsAbortReason(state);if aborted{panic(tsThrown{reason})};return tsU})))
    signal.set("addEventListener",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{event:=tsText(tsArg(args,0));listener:=tsArg(args,1);if tsNullish(listener){return tsU};if listener.kind!=tsFunctionKind{tsFSInvalid("listener","a function")};once:=false;options:=tsArg(args,2);if options.kind==tsObjectKind||options.kind==tsInstanceKind{once=tsTruthy(tsFSProperty(options,"once"))};state.events.listeners[event]=append(state.events.listeners[event],tsFSListener{callback:listener,once:once});return tsU})))
    signal.set("removeEventListener",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{tsCall(loop,signal.get("removeListener"),args...);return tsU})))
    controller:=tsNewObject();controller.nativeClass=tsNativeClass("AbortController");controller.set("signal",tsObjectValue(signal));controller.set("abort",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{reason:=tsArg(args,0);if tsIsUndefined(reason){reason=tsAbortDefault()};tsAbortTrigger(state,reason);return tsU})))
    return tsObjectValue(controller)
}
func tsAbortTrigger(state *tsAbortState,reason tsValue){
    state.Lock();if state.aborted{state.Unlock();return};state.aborted=true;state.reason=reason;followers:=state.followers;state.followers=nil;close(state.done);state.Unlock()
    deliver:=func(){loop:=state.events.loop;signal:=tsObjectValue(state.events.object);state.events.object.set("aborted",tsBooleanValue(true));state.events.object.set("reason",reason);event:=tsNewObject();event.set("type",tsStringReference(tsStringUTF8("abort")));event.set("target",signal);event.set("currentTarget",signal);eventValue:=tsObjectValue(event);state.events.emitEvent("abort",eventValue);handler:=state.events.object.get("onabort");if handler.kind==tsFunctionKind{tsCall(loop,handler,eventValue)}}
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
