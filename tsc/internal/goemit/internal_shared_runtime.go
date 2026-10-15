package goemit

// internalSharedRuntime holds the small common helpers that unblock Bun's
// internal module graph. The remaining shared exports are filled in at their
// saved source checkpoint after node:util's lazy references are available.
const internalSharedRuntime = `
var tsInternalSharedSingleton struct {
    sync.Once
    module *tsObject
    customPromisifyArgs tsValue
    emptyObject, internalSendOptions tsValue
    observerCounts map[string]int
    observers []tsValue
    nodeEntryTypes tsValue
    warnedFeatures map[string]bool
}

type tsInternalSharedNodeObserver struct { object *tsObject; callback,owner,frame tsValue; types map[string]bool; buffer []tsValue; scheduled bool; loop *tsLoop }
var tsInternalSharedNodeObserverStates=map[*tsObject]*tsInternalSharedNodeObserver{}

func tsInternalSharedHideFromStack(args []tsValue) tsValue {
    for _,value:=range args {
        if value.kind!=tsFunctionKind { continue }
        function:=(*tsFunction)(value.ref)
        function.name="::bunternal::"
        if function.properties==nil { function.properties=tsNewObject() }
        function.properties.set("name",tsStringReference(tsStringUTF8(function.name)))
        if function.properties.descriptors==nil{function.properties.descriptors=map[string]*tsDescriptor{}}
        function.properties.descriptors["name"]=&tsDescriptor{configurable:true}
        if object:=function.metadataRef.Load();object!=nil {
            object.set("name",tsStringReference(tsStringUTF8(function.name)))
            if object.descriptors==nil{object.descriptors=map[string]*tsDescriptor{}}
            object.descriptors["name"]=&tsDescriptor{configurable:true}
        }
    }
    return tsU
}

func tsInternalSharedNotImplemented(feature tsValue, issue, extra tsValue) tsValue {
    name:=tsText(feature);message:=name+" is not yet implemented in Bun."
    if tsTruthy(issue){message+=" Track the status & thumbs up the issue: https://github.com/oven-sh/bun/issues/"+tsText(issue)}
    if tsTruthy(extra){message+=". "+tsText(extra)}
    return tsErrorValue(&tsRuntimeError{name:"NotImplementedError",message:message,fields:map[string]tsValue{"code":tsStringReference(tsStringUTF8("ERR_NOT_IMPLEMENTED"))}})
}

func tsInternalSharedWarnNotImplementedOnce(loop *tsLoop, feature,issue tsValue) tsValue {
    name:=tsText(feature);if tsInternalSharedSingleton.warnedFeatures==nil{tsInternalSharedSingleton.warnedFeatures=map[string]bool{}}
    if tsInternalSharedSingleton.warnedFeatures[name]{return tsU};tsInternalSharedSingleton.warnedFeatures[name]=true
    warning:=tsInternalSharedNotImplemented(feature,issue,tsU);console:=tsNodeConsoleModule();warn:=tsGet(loop,console,tsStringReference(tsStringUTF8("warn")));tsCallReceiver(loop,warn,console,warning);return tsU
}

func tsInternalSharedSystemErrorName(loop *tsLoop,errno tsValue) tsValue {
    util:=tsNodeUtilModule(loop)
    getName:=tsGet(loop,util,tsStringReference(tsStringUTF8("getSystemErrorName")))
    return tsCall(loop,getName,errno)
}

func tsInternalSharedIsInsideNodeModules(frameLimit int) bool {
    if frameLimit<1 {return false};pcs:=make([]uintptr,64);count:=runtime.Callers(2,pcs);if count>frameLimit{count=frameLimit};frames:=runtime.CallersFrames(pcs[:count])
    for {frame,more:=frames.Next();path:=filepath.ToSlash(frame.File);if strings.Contains(path,"/node_modules/"){return true};if !more{break}}
    return false
}

func tsInternalSharedOnce(callback tsValue, preserve bool) tsValue {
    called:=false;var returnValue tsValue
    return tsNativeMethod(func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{
        if called {if preserve{return returnValue};return tsU}
        called=true;fn:=callback;callback=tsU
        result:=tsCallReceiver(loop,fn,receiver,args...)
        if preserve{returnValue=result}
        return result
    })
}

func tsInternalSharedGetLazy(initializer tsValue) tsValue {
    initialized:=false;var value tsValue
    return tsNativeMethod(func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{
        if initialized{return value}
        next:=tsCall(loop,initializer)
        value=next;initialized=true
        return value
    })
}

func tsInternalSharedReportUncaught(loop *tsLoop,err tsValue) {
    processValue:=tsNodeProcessModule(loop);process:=(*tsObject)(processValue.ref)
    if !tsTruthy(tsNodeProcessFatalException(loop,process,[]tsValue{err})){loop.errors=append(loop.errors,err)}
}

func tsInternalSharedGuardCallback(loop *tsLoop,callback tsValue) tsValue {
    return tsNativeMethod(func(loop *tsLoop,_ tsValue,args ...tsValue)(result tsValue){
        defer func(){if recovered:=recover();recovered!=nil{tsInternalSharedReportUncaught(loop,tsUnwrap(recovered));result=tsU}}()
        return tsCallReceiver(loop,callback,tsU,args...)
    })
}

func tsInternalSharedPerformanceNow(loop *tsLoop) float64 {
    module:=tsNodePerfHooksModule(loop);performance:=tsGet(loop,module,tsStringReference(tsStringUTF8("performance")))
    now:=tsGet(loop,performance,tsStringReference(tsStringUTF8("now")))
    return tsNumber(tsCallReceiver(loop,now,performance))
}

func tsInternalSharedPerformanceEntry(name,kind string,start,duration float64,detail tsValue) tsValue {
    out:=tsNewObject();out.set("name",tsStringReference(tsStringUTF8(name)));out.set("entryType",tsStringReference(tsStringUTF8(kind)))
    out.set("startTime",tsNumberValue(start));out.set("duration",tsNumberValue(duration));if tsIsUndefined(detail){detail=tsU};out.set("detail",detail)
    out.set("toJSON",tsNamedNativeMethod("toJSON",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{copy:=tsNewObject();for _,key:=range []string{"name","entryType","startTime","duration","detail"}{copy.set(key,out.get(key))};return tsObjectValue(copy)}))
    return tsObjectValue(out)
}

func tsInternalSharedHasObserver(_ *tsLoop,kind tsValue) tsValue {
    return tsBooleanValue(tsInternalSharedSingleton.observerCounts[tsText(kind)]>0)
}

func tsInternalSharedEnqueueNodeEntry(loop *tsLoop,entry tsValue) {
    for _,observer:=range tsInternalSharedSingleton.observers {
        method:=tsGet(loop,observer,tsStringReference(tsStringUTF8("bufferEntry")))
        if method.kind==tsFunctionKind {tsCallReceiver(loop,method,observer,entry)}
    }
}

func tsInternalSharedStartPerf(loop *tsLoop,target,key,context tsValue) {
    tsSet(loop,context,tsStringReference(tsStringUTF8("startTime")),tsNumberValue(tsInternalSharedPerformanceNow(loop)))
    tsSet(loop,target,key,context)
}

func tsInternalSharedStopPerf(loop *tsLoop,target,key,additional tsValue) {
    context:=tsGet(loop,target,key);if !tsTruthy(context){return};tsSet(loop,target,key,tsU)
    start:=tsGet(loop,context,tsStringReference(tsStringUTF8("startTime")))
    detail:=tsGet(loop,context,tsStringReference(tsStringUTF8("detail")));extra:=tsU;if !tsNullish(additional){extra=tsGet(loop,additional,tsStringReference(tsStringUTF8("detail")))}
    if !tsIsUndefined(detail)||!tsIsUndefined(extra) {
        merged:=tsNewObject();for _,source:=range []tsValue{detail,extra}{if !tsNullish(source){for _,key:=range tsObjectKeys(source).values{name:=tsText(key);tsSet(loop,merged,key,tsGet(loop,source,key))}}};detail=tsObjectValue(merged)
    }
    name:=tsText(tsGet(loop,context,tsStringReference(tsStringUTF8("name"))));kind:=tsText(tsGet(loop,context,tsStringReference(tsStringUTF8("type"))))
    entry:=tsInternalSharedPerformanceEntry(name,kind,tsNumber(start),tsInternalSharedPerformanceNow(loop)-tsNumber(start),detail)
    tsInternalSharedEnqueueNodeEntry(loop,entry)
}

func tsInternalSharedPerformanceNodeEntry() tsValue {
    constructor:=&tsFunction{constructible:true,name:"PerformanceNodeEntry",length:5,properties:tsNewObject()}
    constructor.construct=func(_ *tsLoop,args ...tsValue)tsValue{return tsInternalSharedPerformanceEntry(tsText(tsArg(args,0)),tsText(tsArg(args,1)),tsNumber(tsArg(args,2)),tsNumber(tsArg(args,3)),tsArg(args,4))}
    constructor.properties.set("name",tsStringReference(tsStringUTF8("PerformanceNodeEntry")));constructor.properties.set("length",tsNumberValue(5))
    return tsFunctionValue(constructor)
}

func tsInternalSharedIsFrameStopped(_ *tsLoop, _ tsValue) bool { return false }

func tsInternalSharedNodeEntryList(entries []tsValue) tsValue {
    sorted:=append([]tsValue{},entries...)
    sort.SliceStable(sorted,func(i,j int)bool{return tsNumber(tsGet(sorted[i],tsStringReference(tsStringUTF8("startTime"))))<tsNumber(tsGet(sorted[j],tsStringReference(tsStringUTF8("startTime"))))})
    out:=tsNewObject()
    out.set("getEntries",tsNamedNativeMethod("getEntries",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsArrayValue(&tsArray{values:append([]tsValue{},sorted...)})}))
    out.set("getEntriesByType",tsNamedNativeMethod("getEntriesByType",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{result:=&tsArray{};kind:=tsArg(args,0);for _,entry:=range sorted{if tsStrictEqual(tsGet(entry,tsStringReference(tsStringUTF8("entryType"))),kind){result.values=append(result.values,entry)}};return tsArrayValue(result)}))
    out.set("getEntriesByName",tsNamedNativeMethod("getEntriesByName",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{result:=&tsArray{};name,kind:=tsArg(args,0),tsArg(args,1);for _,entry:=range sorted{if !tsStrictEqual(tsGet(entry,tsStringReference(tsStringUTF8("name"))),name){continue};if !tsIsUndefined(kind)&&!tsStrictEqual(tsGet(entry,tsStringReference(tsStringUTF8("entryType"))),kind){continue};result.values=append(result.values,entry)};return tsArrayValue(result)}))
    return tsObjectValue(out)
}

func tsInternalSharedNodeObserverObserve(loop *tsLoop,state *tsInternalSharedNodeObserver,types tsValue) tsValue {
    for kind:=range state.types{tsInternalSharedSingleton.observerCounts[kind]--}
    state.types=map[string]bool{}
    if !tsNullish(types){iterator:=tsIterate(types);for iterator.next(){state.types[tsText(iterator.value)]=true}}
    for kind:=range state.types{tsInternalSharedSingleton.observerCounts[kind]++}
    found:=false;for _,observer:=range tsInternalSharedSingleton.observers{if observer.ref==unsafe.Pointer(state.object){found=true;break}}
    if !found{tsInternalSharedSingleton.observers=append(tsInternalSharedSingleton.observers,tsObjectValue(state.object))}
    visibleTypes:=tsNewCollection(true);for kind:=range state.types{visibleTypes.put(tsStringReference(tsStringUTF8(kind)),tsStringReference(tsStringUTF8(kind)))};state.object.set("types",tsCollectionValue(visibleTypes))
    return tsU
}

func tsInternalSharedNodeObserverDisconnect(state *tsInternalSharedNodeObserver) tsValue {
    for kind:=range state.types{tsInternalSharedSingleton.observerCounts[kind]--};state.types=map[string]bool{};state.buffer=nil
    state.object.set("types",tsCollectionValue(tsNewCollection(true)));state.object.set("buffer",tsArrayValue(&tsArray{}))
    kept:=[]tsValue{};for _,observer:=range tsInternalSharedSingleton.observers{if observer.ref!=unsafe.Pointer(state.object){kept=append(kept,observer)}};tsInternalSharedSingleton.observers=kept
    return tsU
}

func tsInternalSharedNodeObserverBuffer(loop *tsLoop,state *tsInternalSharedNodeObserver,entry tsValue) tsValue {
    kind:=tsText(tsGet(loop,entry,tsStringReference(tsStringUTF8("entryType"))));if !state.types[kind]{return tsU}
    if tsInternalSharedIsFrameStopped(loop,state.frame){tsInternalSharedNodeObserverDisconnect(state);return tsU}
    state.buffer=append(state.buffer,entry);buffer:=&tsArray{};buffer.values=append(buffer.values,state.buffer...);state.object.set("buffer",tsArrayValue(buffer));if state.scheduled{return tsU};state.scheduled=true;state.object.set("scheduled",tsBooleanValue(true))
    deliver:=tsNamedNativeMethod("deliverNodeEntries",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{
        state.scheduled=false;state.object.set("scheduled",tsBooleanValue(false));entries:=state.buffer;state.buffer=nil;state.object.set("buffer",tsArrayValue(&tsArray{}));if len(entries)==0{return tsU}
        callback:=state.callback;list:=tsInternalSharedNodeEntryList(entries);tsCallReceiver(loop,callback,tsU,list,state.owner);return tsU
    })
    timers:=tsNodeTimersModule();setImmediate:=tsGet(loop,timers,tsStringReference(tsStringUTF8("setImmediate")))
    frameModule:=tsInternalAsyncContextFrameModule(loop);currentGraph:=tsGet(loop,frameModule,tsStringReference(tsStringUTF8("currentGraph")));graph:=tsCallReceiver(loop,currentGraph,frameModule)
    if tsIsUndefined(state.frame)&&tsIsUndefined(graph){tsCallReceiver(loop,setImmediate,timers,deliver)}else{run:=tsGet(loop,frameModule,tsStringReference(tsStringUTF8("run")));tsCallReceiver(loop,run,frameModule,state.frame,setImmediate,tsU,deliver)}
    return tsU
}

func tsInternalSharedNodeObserverConstructor() tsValue {
    constructor:=&tsFunction{constructible:true,name:"NodeEntryObserver",length:2,properties:tsNewObject()}
    constructor.construct=func(loop *tsLoop,args ...tsValue)tsValue{
        callback:=tsArg(args,0);if callback.kind!=tsFunctionKind{tsInternalValidatorInvalidType("callback","a function",callback)}
        object:=tsNewObject();state:=&tsInternalSharedNodeObserver{object:object,callback:callback,owner:tsArg(args,1),types:map[string]bool{},buffer:[]tsValue{},loop:loop}
        frameModule:=tsInternalAsyncContextFrameModule(loop);current:=tsGet(loop,frameModule,tsStringReference(tsStringUTF8("currentGraphFrame")));state.frame=tsCallReceiver(loop,current,frameModule)
        tsInternalSharedNodeObserverStates[object]=state
        object.set("callback",callback);object.set("owner",state.owner);object.set("types",tsCollectionValue(tsNewCollection(true)));object.set("buffer",tsArrayValue(&tsArray{}));object.set("scheduled",tsBooleanValue(false));object.set("frame",state.frame)
        object.set("observe",tsNamedNativeMethod("observe",1,func(loop *tsLoop,receiver tsValue,values ...tsValue)tsValue{return tsInternalSharedNodeObserverObserve(loop,tsInternalSharedNodeObserverStates[(*tsObject)(receiver.ref)],tsArg(values,0))}))
        object.set("disconnect",tsNamedNativeMethod("disconnect",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{return tsInternalSharedNodeObserverDisconnect(tsInternalSharedNodeObserverStates[(*tsObject)(receiver.ref)])}))
        object.set("bufferEntry",tsNamedNativeMethod("bufferEntry",1,func(loop *tsLoop,receiver tsValue,values ...tsValue)tsValue{return tsInternalSharedNodeObserverBuffer(loop,tsInternalSharedNodeObserverStates[(*tsObject)(receiver.ref)],tsArg(values,0))}))
        return tsObjectValue(object)
    }
    constructor.properties.set("name",tsStringReference(tsStringUTF8("NodeEntryObserver")));constructor.properties.set("length",tsNumberValue(2))
    return tsFunctionValue(constructor)
}

func tsInternalSharedErrorClass(name string,length int,build func(*tsLoop,[]tsValue,tsValue)tsValue) tsValue {
    prototype:=tsNewObject();prototype.prototype=tsErrorPrototype("Error")
    constructor:=tsFunc(func(_ *tsLoop,_ ...tsValue)tsValue{tsPropertyFailure("Class constructor "+name+" cannot be invoked without 'new'");return tsU})
    constructor.name=name;constructor.length=length;constructor.constructible=true
    constructor.construct=func(loop *tsLoop,args ...tsValue)tsValue{return build(loop,args,tsObjectValue(prototype))}
    constructor.properties=tsNewObject();constructor.properties.set("prototype",tsObjectValue(prototype));constructor.properties.descriptors=map[string]*tsDescriptor{"prototype":{}}
    prototype.set("constructor",tsFunctionValue(constructor));prototype.descriptors=map[string]*tsDescriptor{"constructor":{writable:true,configurable:true}}
    return tsFunctionValue(constructor)
}

func tsInternalSharedErrorValue(loop *tsLoop,message tsValue,prototype tsValue) tsValue {
    value:=tsNewError(loop,"Error",[]tsValue{message})
    (*tsObject)(value.ref).prototype=prototype
    return value
}

func tsInternalSharedConstructors(loop *tsLoop,module *tsObject) {
    hostPort:=tsInternalSharedErrorClass("ExceptionWithHostPort",5,func(loop *tsLoop,args []tsValue,prototype tsValue)tsValue{
        errno,syscall,address,port,additional:=tsArg(args,0),tsArg(args,1),tsArg(args,2),tsArg(args,3),tsArg(args,4)
        code:=tsInternalSharedSystemErrorName(loop,errno)
        details:="";if tsTruthy(port)&&tsNumber(port)>0{details=" "+tsText(address)+":"+tsText(port)}else if tsTruthy(address){details=" "+tsText(address)}
        if tsTruthy(additional){details+=" - Local ("+tsText(additional)+")"}
        value:=tsInternalSharedErrorValue(loop,tsStringReference(tsStringUTF8(tsText(syscall)+" "+tsText(code)+details)),prototype)
        object:=(*tsObject)(value.ref);object.set("errno",errno);object.set("code",code);object.set("syscall",syscall);object.set("address",address);if tsTruthy(port){object.set("port",port)}
        return value
    })
    errnoException:=tsInternalSharedErrorClass("ErrnoException",3,func(loop *tsLoop,args []tsValue,prototype tsValue)tsValue{
        errno,syscall,original:=tsArg(args,0),tsArg(args,1),tsArg(args,2);code:=tsInternalSharedSystemErrorName(loop,errno)
        message:=tsText(syscall)+" "+tsText(code);if tsTruthy(original){message+=" "+tsText(original)}
        value:=tsInternalSharedErrorValue(loop,tsStringReference(tsStringUTF8(message)),prototype);object:=(*tsObject)(value.ref);object.set("errno",errno);object.set("code",code);object.set("syscall",syscall);return value
    })
    connReset:=tsInternalSharedErrorClass("ConnResetException",1,func(loop *tsLoop,args []tsValue,prototype tsValue)tsValue{
        value:=tsInternalSharedErrorValue(loop,tsArg(args,0),prototype);(*tsObject)(value.ref).set("code",tsStringReference(tsStringUTF8("ECONNRESET")));return value
    })
    aggregate:=tsInternalSharedErrorClass("NodeAggregateError",2,func(loop *tsLoop,args []tsValue,prototype tsValue)tsValue{
        value:=tsNewError(loop,"AggregateError",args);object:=(*tsObject)(value.ref);object.prototype=prototype
        errorCode:=tsU;errors:=tsArg(args,0);if errors.kind==tsArrayKind&&(*tsArray)(errors.ref).length()>0{errorCode=tsGet(loop,(*tsArray)(errors.ref).at(0),tsStringReference(tsStringUTF8("code")))};object.set("code",errorCode)
        return value
    })
    module.set("ExceptionWithHostPort",hostPort);module.set("ErrnoException",errnoException);module.set("ConnResetException",connReset);module.set("NodeAggregateError",aggregate)
    errorConstructor:=tsErrorConstructor("Error");aggregateConstructor:=tsErrorConstructor("AggregateError")
    for _,item:=range []struct{constructor,returns tsValue}{{hostPort,errorConstructor},{errnoException,errorConstructor},{connReset,errorConstructor},{aggregate,aggregateConstructor}}{prototype:=tsGet(loop,item.constructor,tsStringReference(tsStringUTF8("prototype")));tsSet(loop,prototype,tsStringReference(tsStringUTF8("constructor")),item.returns)}
}

func tsInternalSharedModule(loop *tsLoop) tsValue {
    tsInternalSharedSingleton.Do(func(){
        module:=tsNewObject()
        tsInternalSharedSingleton.customPromisifyArgs=tsNewSymbol(tsStringUTF8("customPromisifyArgs"))
        empty:=tsNewObject();empty.prototype=tsNull;empty.nonExtensible=true;tsInternalSharedSingleton.emptyObject=tsObjectValue(empty)
        privateOptions:=tsNewObject();privateOptions.prototype=tsNull;privateOptions.set("$internal",tsBooleanValue(true));privateOptions.nonExtensible=true;tsInternalSharedSingleton.internalSendOptions=tsObjectValue(privateOptions)
        tsInternalSharedSingleton.observerCounts=map[string]int{};tsInternalSharedSingleton.observers=[]tsValue{}
        types:=[]tsValue{};for _,name:=range []string{"net","dns","http","http2","function","quic"}{types=append(types,tsStringReference(tsStringUTF8(name)))};tsInternalSharedSingleton.nodeEntryTypes=tsNewCollection(true,tsArrayValue(&tsArray{values:types}))
        module.set("hideFromStack",tsFunc(func(_ *tsLoop,args ...tsValue)tsValue{return tsInternalSharedHideFromStack(args)}))
        throwNotImplemented:=tsNamedNativeMethod("throwNotImplemented",3,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{panic(tsThrown{tsInternalSharedNotImplemented(tsArg(args,0),tsArg(args,1),tsArg(args,2))})})
        tsInternalSharedHideFromStack([]tsValue{throwNotImplemented});module.set("throwNotImplemented",throwNotImplemented)
        module.set("warnNotImplementedOnce",tsNamedNativeMethod("warnNotImplementedOnce",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsInternalSharedWarnNotImplementedOnce(loop,tsArg(args,0),tsArg(args,1))}))
        module.set("kCustomPromisifyArgsSymbol",tsInternalSharedSingleton.customPromisifyArgs)
        module.set("isInsideNodeModules",tsNamedNativeMethod("isInsideNodeModules",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{limit:=int(tsNumber(tsArg(args,0)));return tsBooleanValue(tsInternalSharedIsInsideNodeModules(limit))}))
        tsInternalSharedConstructors(loop,module)
        module.set("once",tsNamedNativeMethod("once",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{options:=tsArg(args,1);preserve:=false;if !tsIsUndefined(options)&&!tsIsNull(options){preserve=tsTruthy(tsGet(loop,options,tsStringReference(tsStringUTF8("preserveReturnValue"))))};return tsInternalSharedOnce(tsArg(args,0),preserve)}))
        module.set("getLazy",tsNamedNativeMethod("getLazy",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsInternalSharedGetLazy(tsArg(args,0))}))
        module.set("guardCallback",tsNamedNativeMethod("guardCallback",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsInternalSharedGuardCallback(loop,tsArg(args,0))}))
        module.set("resistStopPropagation",tsNamedNativeMethod("resistStopPropagation",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{options:=tsArg(args,0);tsSet(loop,options,tsStringReference(tsStringUTF8("$kResistStopPropagation")),tsBooleanValue(true));return options}))
        module.set("kEmptyObject",tsInternalSharedSingleton.emptyObject)
        module.set("kInternalSendOptions",tsInternalSharedSingleton.internalSendOptions)
        module.set("hasObserver",tsNamedNativeMethod("hasObserver",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsInternalSharedHasObserver(loop,tsArg(args,0))}))
        module.set("startPerf",tsNamedNativeMethod("startPerf",3,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{tsInternalSharedStartPerf(loop,tsArg(args,0),tsArg(args,1),tsArg(args,2));return tsU}))
        module.set("stopPerf",tsNamedNativeMethod("stopPerf",3,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{tsInternalSharedStopPerf(loop,tsArg(args,0),tsArg(args,1),tsArg(args,2));return tsU}))
        module.set("enqueueNodeEntry",tsNamedNativeMethod("enqueueNodeEntry",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{tsInternalSharedEnqueueNodeEntry(loop,tsArg(args,0));return tsU}))
        module.set("kNodeEntryTypes",tsInternalSharedSingleton.nodeEntryTypes)
        module.set("PerformanceNodeEntry",tsInternalSharedPerformanceNodeEntry())
        module.set("isStoppedModuleGraphRunning",tsNamedNativeMethod("isStoppedModuleGraphRunning",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{frames:=tsInternalAsyncContextFrameModule(loop);current:=tsGet(loop,frames,tsStringReference(tsStringUTF8("current")));frame:=tsCallReceiver(loop,current,frames);return tsBooleanValue(tsInternalSharedIsFrameStopped(loop,frame))}))
        suffix:="\nThis is caused by either a bug in Node.js or incorrect usage of Node.js internals.\nPlease open an issue with this stack trace at https://github.com/nodejs/node/issues\n"
        module.set("kInternalAssertionSuffix",tsStringReference(tsStringUTF8(suffix)))
        module.set("NodeEntryObserver",tsInternalSharedNodeObserverConstructor())
        for name,description:=range map[string]string{"kHandle":"kHandle","kClusterOwner":"kClusterOwner","kAutoDestroyed":"kAutoDestroyed","kWeakHandler":"kWeak"}{module.set(name,tsNewSymbol(tsStringUTF8(description)))}
        tsInternalSharedSingleton.module=module
    })
    return tsObjectValue(tsInternalSharedSingleton.module)
}
`
