package goemit

// internalStreamLegacyRuntime supplies the EventEmitter based stream base
// class used by Bun's internal readable, writable and duplex implementations.
const internalStreamLegacyRuntime = `
var tsLegacyStreamSingleton struct {
    sync.Once
    module *tsObject
    constructor *tsFunction
    prototype *tsObject
}

func tsLegacyEventCall(loop *tsLoop, receiver tsValue, name string, args ...tsValue) tsValue {
    return tsCallReceiver(loop,tsGet(loop,receiver,tsStringReference(tsStringUTF8(name))),receiver,args...)
}

func tsLegacyStreamCreate(loop *tsLoop, receiver, options tsValue) tsValue {
    tsNodeEventsModule(loop)
    var object *tsObject
    if receiver.kind == tsObjectKind { object = (*tsObject)(receiver.ref) } else { object = tsNewObject() }
    value := tsNodeEventCreate(loop,tsObjectValue(object),options)
    object.prototype = tsObjectValue(tsLegacyStreamSingleton.prototype)
    return value
}

func tsLegacyPrependListener(loop *tsLoop, emitter,event,listener tsValue) {
    prepend := tsGet(loop,emitter,tsStringReference(tsStringUTF8("prependListener")))
    if prepend.kind == tsFunctionKind {
        tsCallReceiver(loop,prepend,emitter,event,listener)
        return
    }
    tsCallReceiver(loop,tsGet(loop,emitter,tsStringReference(tsStringUTF8("on"))),emitter,event,listener)
}

func tsLegacyCleanup(loop *tsLoop, source,dest tsValue, ondata,ondrain,onend,onclose,onerror,cleanup tsValue) {
    for _,entry := range []struct{target tsValue;event string;listener tsValue}{
        {source,"data",ondata},{dest,"drain",ondrain},{source,"end",onend},{source,"close",onclose},
        {source,"error",onerror},{dest,"error",onerror},{source,"end",cleanup},{source,"close",cleanup},{dest,"close",cleanup},
    } {
        tsLegacyEventCall(loop,entry.target,"removeListener",tsStringReference(tsStringUTF8(entry.event)),entry.listener)
    }
}

func tsLegacyPipe(loop *tsLoop, source tsValue, args []tsValue) tsValue {
    dest := tsArg(args,0)
    if tsNullish(dest) { tsFSInvalid("destination", "a writable stream") }
    options := tsArg(args,1)
    ondata := tsFunctionValue(tsFunc(func(loop *tsLoop,chunks ...tsValue) tsValue {
        chunk := tsArg(chunks,0)
        writable := tsGet(loop,dest,tsStringReference(tsStringUTF8("writable")))
        result := tsLegacyEventCall(loop,dest,"write",chunk)
        pause := tsGet(loop,source,tsStringReference(tsStringUTF8("pause")))
        if tsTruthy(writable) && result.kind == tsBooleanKind && !tsTruthy(result) && pause.kind == tsFunctionKind { tsCallReceiver(loop,pause,source) }
        return tsU
    }))
    ondrain := tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue) tsValue {
        readable := tsGet(loop,source,tsStringReference(tsStringUTF8("readable")))
        resume := tsGet(loop,source,tsStringReference(tsStringUTF8("resume")))
        if tsTruthy(readable) && resume.kind == tsFunctionKind { tsCallReceiver(loop,resume,source) }
        return tsU
    }))
    didEnd := false
    onend := tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue) tsValue {
        if didEnd { return tsU }; didEnd = true
        tsLegacyEventCall(loop,dest,"end")
        return tsU
    }))
    onclose := tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue) tsValue {
        if didEnd { return tsU }; didEnd = true
        destroy := tsGet(loop,dest,tsStringReference(tsStringUTF8("destroy")))
        if destroy.kind == tsFunctionKind { tsCallReceiver(loop,destroy,dest) }
        return tsU
    }))
    var onerror, cleanup tsValue
    cleanup = tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue) tsValue {
        tsLegacyCleanup(loop,source,dest,ondata,ondrain,onend,onclose,onerror,cleanup)
        return tsU
    }))
    onerror = tsNamedNativeMethod("onerror",1,func(loop *tsLoop,emitter tsValue,errors ...tsValue) tsValue {
        tsLegacyCleanup(loop,source,dest,ondata,ondrain,onend,onclose,onerror,cleanup)
        event := tsStringReference(tsStringUTF8("error"))
        count := tsLegacyEventCall(loop,emitter,"listenerCount",event)
        if tsIsNumeric(count) && tsNumber(count) == 0 { tsLegacyEventCall(loop,emitter,"emit",event,tsArg(errors,0)) }
        return tsU
    })

    tsLegacyEventCall(loop,source,"on",tsStringReference(tsStringUTF8("data")),ondata)
    tsLegacyEventCall(loop,dest,"on",tsStringReference(tsStringUTF8("drain")),ondrain)
    shouldEnd := !tsTruthy(tsGet(loop,dest,tsStringReference(tsStringUTF8("_isStdio")))) && (tsNullish(options) || !tsStrictEqual(tsGet(loop,options,tsStringReference(tsStringUTF8("end"))),tsBooleanValue(false)))
    if shouldEnd {
        tsLegacyEventCall(loop,source,"on",tsStringReference(tsStringUTF8("end")),onend)
        tsLegacyEventCall(loop,source,"on",tsStringReference(tsStringUTF8("close")),onclose)
    }
    tsLegacyPrependListener(loop,source,tsStringReference(tsStringUTF8("error")),onerror)
    tsLegacyPrependListener(loop,dest,tsStringReference(tsStringUTF8("error")),onerror)
    tsLegacyEventCall(loop,source,"on",tsStringReference(tsStringUTF8("end")),cleanup)
    tsLegacyEventCall(loop,source,"on",tsStringReference(tsStringUTF8("close")),cleanup)
    tsLegacyEventCall(loop,dest,"on",tsStringReference(tsStringUTF8("close")),cleanup)
    tsLegacyEventCall(loop,dest,"emit",tsStringReference(tsStringUTF8("pipe")),source)
    return dest
}

func tsLegacyEventNames(loop *tsLoop, receiver tsValue) tsValue {
    state := tsNodeEventState(receiver)
    out := &tsArray{}
    if state == nil { return tsArrayValue(out) }
    for key,list := range state.events { if len(list)>0 { out.values=append(out.values,state.names[key]) } }
    return tsArrayValue(out)
}

func tsInternalStreamLegacyModule(loop *tsLoop) tsValue {
    tsNodeEventsModule(loop)
    tsLegacyStreamSingleton.Do(func() {
        prototype := tsNewObject()
        prototype.prototype = tsObjectValue(tsNodeEventsSingleton.prototype)
        function := tsFunc(func(loop *tsLoop,args ...tsValue) tsValue { return tsLegacyStreamCreate(loop,tsU,tsArg(args,0)) })
        function.name = "Stream"
        function.length = 1
        function.constructible = true
        function.construct = func(loop *tsLoop,args ...tsValue) tsValue { return tsLegacyStreamCreate(loop,tsU,tsArg(args,0)) }
        function.receiverCall = func(loop *tsLoop,receiver tsValue,args ...tsValue) tsValue {
            if receiver.kind == tsObjectKind && tsPrototypeContains(receiver,tsObjectValue(prototype)) { return tsLegacyStreamCreate(loop,receiver,tsArg(args,0)) }
            return tsLegacyStreamCreate(loop,tsU,tsArg(args,0))
        }
        function.properties = tsNewObject()
        function.properties.set("prototype",tsObjectValue(prototype))
        function.properties.descriptors = map[string]*tsDescriptor{"prototype":{}}
        prototype.set("constructor",tsFunctionValue(function))
        prototype.descriptors = map[string]*tsDescriptor{"constructor":{writable:true,configurable:true}}
        prototype.set("pipe",tsNamedNativeMethod("pipe",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsLegacyPipe(loop,receiver,args)}))
        prototype.set("eventNames",tsNamedNativeMethod("eventNames",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{return tsLegacyEventNames(loop,receiver)}))
        for _,name := range []string{"pipe","eventNames"} { prototype.descriptors[name]=&tsDescriptor{writable:true,configurable:true} }
        isView := tsNamedNativeMethod("isArrayBufferView",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{value:=tsArg(args,0);return tsBooleanValue(value.kind==tsTypedArrayKind||value.kind==tsDataViewKind)})
        isUint8 := tsNamedNativeMethod("isUint8Array",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{value:=tsArg(args,0);if value.kind!=tsTypedArrayKind{return tsBooleanValue(false)};array:=(*tsTypedArray)(value.ref);return tsBooleanValue(array.name=="Uint8Array"||array.nodeBuffer)})
        toBuffer := tsNamedNativeMethod("_uint8ArrayToBuffer",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{
            chunk:=tsArg(args,0);if chunk.kind!=tsTypedArrayKind||(*tsTypedArray)(chunk.ref).name!="Uint8Array" {tsFSInvalid("chunk","a Uint8Array")}
            array:=(*tsTypedArray)(chunk.ref)
            return tsNodeBufferFrom(loop,[]tsValue{tsArrayBufferValue(array.buffer),tsNumberValue(float64(array.offset)),tsNumberValue(float64(array.length*array.bytes))})
        })
        function.properties.set("_isArrayBufferView",isView)
        function.properties.set("_isUint8Array",isUint8)
        function.properties.set("_uint8ArrayToBuffer",toBuffer)
        for _,name := range []string{"_isArrayBufferView","_isUint8Array","_uint8ArrayToBuffer"} { function.properties.descriptors[name]=&tsDescriptor{writable:true,configurable:true} }
        module:=tsNewObject()
        module.set("Stream",tsFunctionValue(function))
        module.set("prependListener",tsNamedNativeMethod("prependListener",3,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{tsLegacyPrependListener(loop,tsArg(args,0),tsArg(args,1),tsArg(args,2));return tsU}))
        tsLegacyStreamSingleton.prototype=prototype
        tsLegacyStreamSingleton.constructor=function
        tsLegacyStreamSingleton.module=module
    })
    return tsObjectValue(tsLegacyStreamSingleton.module)
}
`
