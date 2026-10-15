package goemit

// internalStreamUtilsRuntime maps Node's internal stream predicates and
// shared state symbols onto the generated runtime's objects and Web streams.
const internalStreamUtilsRuntime = `
var tsInternalStreamUtils struct {
    sync.Once
    module *tsObject
    keys map[string]tsValue
}

func tsInternalStreamUtilsInit(loop *tsLoop) {
    tsInternalStreamUtils.Do(func() {
        tsInternalStreamUtils.keys=map[string]tsValue{}
        symbolFor:=tsGet(loop,tsSymbolFunction(),tsStringReference(tsStringUTF8("for")))
        for _,name:=range []string{"nodejs.stream.destroyed","nodejs.stream.errored","nodejs.stream.readable","nodejs.stream.writable","nodejs.stream.disturbed","nodejs.webstream.isClosedPromise","nodejs.webstream.controllerErrorFunction"} {
            tsInternalStreamUtils.keys[name]=tsCall(loop,symbolFor,tsStringReference(tsStringUTF8(name)))
        }
        tsInternalStreamUtils.keys["kOnConstructed"]=tsNewSymbol(tsStringUTF8("kOnConstructed"))
        tsInternalStreamUtils.keys["kState"]=tsNewSymbol(tsStringUTF8("kState"))
        tsInternalStreamUtils.keys["kObjectMode"]=tsNumberValue(1<<0)
        tsInternalStreamUtils.keys["kErrorEmitted"]=tsNumberValue(1<<1)
        tsInternalStreamUtils.keys["kAutoDestroy"]=tsNumberValue(1<<2)
        tsInternalStreamUtils.keys["kEmitClose"]=tsNumberValue(1<<3)
        tsInternalStreamUtils.keys["kDestroyed"]=tsNumberValue(1<<4)
        tsInternalStreamUtils.keys["kClosed"]=tsNumberValue(1<<5)
        tsInternalStreamUtils.keys["kCloseEmitted"]=tsNumberValue(1<<6)
        tsInternalStreamUtils.keys["kErrored"]=tsNumberValue(1<<7)
        tsInternalStreamUtils.keys["kConstructed"]=tsNumberValue(1<<8)
    })
}

func tsInternalStreamUtilsKey(loop *tsLoop,name string) tsValue {
    tsInternalStreamUtilsInit(loop)
    return tsInternalStreamUtils.keys[name]
}

func tsInternalStreamUtilsGet(loop *tsLoop,receiver tsValue,name string) tsValue {
    if tsNullish(receiver) { return tsU }
    return tsGet(loop,receiver,tsStringReference(tsStringUTF8(name)))
}

func tsInternalStreamUtilsGetKey(loop *tsLoop,receiver,key tsValue) tsValue {
    if tsNullish(receiver) { return tsU }
    return tsGet(loop,receiver,key)
}

func tsInternalStreamUtilsFunction(loop *tsLoop,receiver tsValue,name string) bool {
    return tsInternalStreamUtilsGet(loop,receiver,name).kind==tsFunctionKind
}

func tsInternalStreamUtilsBool(loop *tsLoop,receiver tsValue,name string) bool {
    value:=tsInternalStreamUtilsGet(loop,receiver,name)
    return value.kind==tsBooleanKind
}

func tsInternalStreamIsReadableNode(loop *tsLoop,stream tsValue,strict bool) bool {
    if !tsTruthy(stream)||!tsInternalStreamUtilsFunction(loop,stream,"pipe")||!tsInternalStreamUtilsFunction(loop,stream,"on") { return false }
    if strict&&(!tsInternalStreamUtilsFunction(loop,stream,"pause")||!tsInternalStreamUtilsFunction(loop,stream,"resume")){return false}
    writable:=tsInternalStreamUtilsGet(loop,stream,"_writableState")
    readable:=tsInternalStreamUtilsGet(loop,stream,"_readableState")
    if !tsNullish(writable)&&tsStrictEqual(tsInternalStreamUtilsGet(loop,readable,"readable"),tsBooleanValue(false)){return false}
    if !tsNullish(writable)&&tsNullish(readable){return false}
    return true
}

func tsInternalStreamIsWritableNode(loop *tsLoop,stream tsValue) bool {
    if !tsTruthy(stream)||!tsInternalStreamUtilsFunction(loop,stream,"write")||!tsInternalStreamUtilsFunction(loop,stream,"on"){return false}
    readable:=tsInternalStreamUtilsGet(loop,stream,"_readableState")
    writable:=tsInternalStreamUtilsGet(loop,stream,"_writableState")
    return tsNullish(readable)||!tsStrictEqual(tsInternalStreamUtilsGet(loop,writable,"writable"),tsBooleanValue(false))
}

func tsInternalStreamIsNode(loop *tsLoop,stream tsValue) bool {
    if !tsTruthy(stream){return false}
    if !tsNullish(tsInternalStreamUtilsGet(loop,stream,"_readableState"))||!tsNullish(tsInternalStreamUtilsGet(loop,stream,"_writableState")){return true}
    return tsInternalStreamUtilsFunction(loop,stream,"write")&&tsInternalStreamUtilsFunction(loop,stream,"on")||tsInternalStreamUtilsFunction(loop,stream,"pipe")&&tsInternalStreamUtilsFunction(loop,stream,"on")
}

func tsInternalStreamIsDuplex(loop *tsLoop,stream tsValue) bool {
    return tsTruthy(stream)&&tsInternalStreamUtilsFunction(loop,stream,"pipe")&&!tsNullish(tsInternalStreamUtilsGet(loop,stream,"_readableState"))&&tsInternalStreamUtilsFunction(loop,stream,"on")&&tsInternalStreamUtilsFunction(loop,stream,"write")
}

func tsInternalStreamIsWebClass(stream tsValue,name string) bool {
    if stream.kind!=tsObjectKind{return false}
    object:=(*tsObject)(stream.ref)
    if object.nativeClass!=nil&&object.nativeClass.builtin==name{return true}
    if object.webStream!=nil&&name=="ReadableStream"{return true}
    return false
}

func tsInternalStreamIsWeb(loop *tsLoop,stream tsValue) bool {
    _=loop
    return tsInternalStreamIsWebClass(stream,"ReadableStream")||tsInternalStreamIsWebClass(stream,"WritableStream")||tsInternalStreamIsWebClass(stream,"TransformStream")
}

func tsInternalStreamState(loop *tsLoop,stream tsValue,name string) tsValue {
    return tsInternalStreamUtilsGet(loop,stream,name)
}

func tsInternalStreamStateBit(loop *tsLoop,stream tsValue,bitName string) bool {
    key:=tsInternalStreamUtilsKey(loop,"kState")
    bits:=tsInternalStreamUtilsGetKey(loop,stream,key)
    bit:=tsNumber(tsInternalStreamUtilsKey(loop,bitName))
    return tsIsNumeric(bits)&&uint32(tsNumber(bits))&uint32(bit)!=0
}

func tsInternalStreamIsDestroyed(loop *tsLoop,stream tsValue) tsValue {
    if !tsInternalStreamIsNode(loop,stream){return tsNull}
    writable:=tsInternalStreamState(loop,stream,"_writableState")
    readable:=tsInternalStreamState(loop,stream,"_readableState")
    state:=writable;if tsNullish(state){state=readable}
    shared:=tsInternalStreamUtilsKey(loop,"nodejs.stream.destroyed")
    return tsBooleanValue(tsTruthy(tsInternalStreamUtilsGet(loop,stream,"destroyed"))||tsTruthy(tsInternalStreamUtilsGetKey(loop,stream,shared))||tsTruthy(tsInternalStreamUtilsGet(loop,state,"destroyed"))||tsInternalStreamStateBit(loop,stream,"kDestroyed"))
}

func tsInternalStreamIsWritableEnded(loop *tsLoop,stream tsValue) tsValue {
    if !tsInternalStreamIsWritableNode(loop,stream){return tsNull}
    if tsStrictEqual(tsInternalStreamUtilsGet(loop,stream,"writableEnded"),tsBooleanValue(true)){return tsBooleanValue(true)}
    state:=tsInternalStreamState(loop,stream,"_writableState")
    if !tsNullish(tsInternalStreamUtilsGet(loop,state,"errored")){return tsBooleanValue(false)}
    ended:=tsInternalStreamUtilsGet(loop,state,"ended")
    if ended.kind!=tsBooleanKind{return tsNull}
    return ended
}

func tsInternalStreamIsWritableFinished(loop *tsLoop,stream,strict tsValue) tsValue {
    if !tsInternalStreamIsWritableNode(loop,stream){return tsNull}
    if tsStrictEqual(tsInternalStreamUtilsGet(loop,stream,"writableFinished"),tsBooleanValue(true)){return tsBooleanValue(true)}
    state:=tsInternalStreamState(loop,stream,"_writableState")
    if !tsNullish(tsInternalStreamUtilsGet(loop,state,"errored")){return tsBooleanValue(false)}
    finished:=tsInternalStreamUtilsGet(loop,state,"finished")
    if finished.kind!=tsBooleanKind{return tsNull}
    relaxed:=tsStrictEqual(strict,tsBooleanValue(false))&&tsStrictEqual(tsInternalStreamUtilsGet(loop,state,"ended"),tsBooleanValue(true))&&tsStrictEqual(tsInternalStreamUtilsGet(loop,state,"length"),tsNumberValue(0))
    return tsBooleanValue(tsTruthy(finished)||relaxed)
}

func tsInternalStreamIsReadableFinished(loop *tsLoop,stream,strict tsValue) tsValue {
    if !tsInternalStreamIsReadableNode(loop,stream,false){return tsNull}
    state:=tsInternalStreamState(loop,stream,"_readableState")
    if !tsNullish(tsInternalStreamUtilsGet(loop,state,"errored")){return tsBooleanValue(false)}
    ended:=tsInternalStreamUtilsGet(loop,state,"endEmitted")
    if ended.kind!=tsBooleanKind{return tsNull}
    relaxed:=tsStrictEqual(strict,tsBooleanValue(false))&&tsStrictEqual(tsInternalStreamUtilsGet(loop,state,"ended"),tsBooleanValue(true))&&tsStrictEqual(tsInternalStreamUtilsGet(loop,state,"length"),tsNumberValue(0))
    return tsBooleanValue(tsTruthy(ended)||relaxed)
}

func tsInternalStreamIsReadable(loop *tsLoop,stream tsValue) tsValue {
    shared:=tsInternalStreamUtilsKey(loop,"nodejs.stream.readable")
    value:=tsInternalStreamUtilsGetKey(loop,stream,shared)
    if !tsNullish(value){return value}
    value=tsInternalStreamUtilsGet(loop,stream,"readable")
    if value.kind!=tsBooleanKind{return tsNull}
    if tsTruthy(tsInternalStreamIsDestroyed(loop,stream)){return tsBooleanValue(false)}
    return tsBooleanValue(tsInternalStreamIsReadableNode(loop,stream,false)&&tsTruthy(value)&&!tsTruthy(tsInternalStreamIsReadableFinished(loop,stream,tsU)))
}

func tsInternalStreamIsWritable(loop *tsLoop,stream tsValue) tsValue {
    shared:=tsInternalStreamUtilsKey(loop,"nodejs.stream.writable")
    value:=tsInternalStreamUtilsGetKey(loop,stream,shared)
    if !tsNullish(value){return value}
    value=tsInternalStreamUtilsGet(loop,stream,"writable")
    if value.kind!=tsBooleanKind{return tsNull}
    if tsTruthy(tsInternalStreamIsDestroyed(loop,stream)){return tsBooleanValue(false)}
    return tsBooleanValue(tsInternalStreamIsWritableNode(loop,stream)&&tsTruthy(value)&&!tsTruthy(tsInternalStreamIsWritableEnded(loop,stream)))
}

func tsInternalStreamIsIterable(loop *tsLoop,stream,async tsValue) tsValue {
    if tsNullish(stream){return tsBooleanValue(false)}
    iterator:=tsGet(loop,stream,tsWellKnownSymbol("iterator"))
    asyncIterator:=tsGet(loop,stream,tsWellKnownSymbol("asyncIterator"))
    if tsStrictEqual(async,tsBooleanValue(true)){return tsBooleanValue(asyncIterator.kind==tsFunctionKind)}
    if tsStrictEqual(async,tsBooleanValue(false)){return tsBooleanValue(iterator.kind==tsFunctionKind)}
    return tsBooleanValue(iterator.kind==tsFunctionKind||asyncIterator.kind==tsFunctionKind)
}

func tsInternalStreamIsFinished(loop *tsLoop,stream,options tsValue) tsValue {
    if !tsInternalStreamIsNode(loop,stream){return tsNull}
    if tsTruthy(tsInternalStreamIsDestroyed(loop,stream)){return tsBooleanValue(true)}
    readable:=tsInternalStreamUtilsGet(loop,options,"readable")
    if !tsStrictEqual(readable,tsBooleanValue(false))&&tsTruthy(tsInternalStreamIsReadable(loop,stream)){return tsBooleanValue(false)}
    writable:=tsInternalStreamUtilsGet(loop,options,"writable")
    if !tsStrictEqual(writable,tsBooleanValue(false))&&tsTruthy(tsInternalStreamIsWritable(loop,stream)){return tsBooleanValue(false)}
    return tsBooleanValue(true)
}

func tsInternalStreamErrored(loop *tsLoop,stream,writable bool) tsValue {
    if !tsInternalStreamIsNode(loop,stream){return tsNull}
    name:="readableErrored";stateName:="_readableState";if writable{name="writableErrored";stateName="_writableState"}
    value:=tsInternalStreamUtilsGet(loop,stream,name);if tsTruthy(value){return value}
    state:=tsInternalStreamState(loop,stream,stateName);value=tsInternalStreamUtilsGet(loop,state,"errored");if tsNullish(value){return tsNull};return value
}

func tsInternalStreamIsOutgoing(loop *tsLoop,stream tsValue) bool {
    for _,name:=range []string{"_closed","_defaultKeepAlive","_removedConnection","_removedContLen"}{if !tsInternalStreamUtilsBool(loop,stream,name){return false}}
    return true
}

func tsInternalStreamIsClosed(loop *tsLoop,stream tsValue) tsValue {
    if !tsInternalStreamIsNode(loop,stream){return tsNull}
    closed:=tsInternalStreamUtilsGet(loop,stream,"closed");if closed.kind==tsBooleanKind{return closed}
    writable:=tsInternalStreamState(loop,stream,"_writableState");readable:=tsInternalStreamState(loop,stream,"_readableState")
    wc,rc:=tsInternalStreamUtilsGet(loop,writable,"closed"),tsInternalStreamUtilsGet(loop,readable,"closed")
    if wc.kind==tsBooleanKind||rc.kind==tsBooleanKind{return tsBooleanValue(tsTruthy(wc)||tsTruthy(rc))}
    if tsInternalStreamUtilsBool(loop,stream,"_closed")&&tsInternalStreamIsOutgoing(loop,stream){return tsInternalStreamUtilsGet(loop,stream,"_closed")}
    return tsNull
}

func tsInternalStreamIsServerRequest(loop *tsLoop,stream tsValue) bool {
    if !tsInternalStreamUtilsBool(loop,stream,"_consuming")||!tsInternalStreamUtilsBool(loop,stream,"_dumped"){return false}
    request:=tsInternalStreamUtilsGet(loop,stream,"req")
    return tsIsUndefined(tsInternalStreamUtilsGet(loop,request,"upgradeOrConnect"))
}

func tsInternalStreamWillEmitClose(loop *tsLoop,stream tsValue) tsValue {
    if !tsInternalStreamIsNode(loop,stream){return tsNull}
    writable:=tsInternalStreamState(loop,stream,"_writableState");readable:=tsInternalStreamState(loop,stream,"_readableState");state:=writable;if tsNullish(state){state=readable}
    if tsNullish(state)&&tsInternalStreamIsOutgoing(loop,stream){return tsBooleanValue(true)}
    return tsBooleanValue(tsTruthy(tsInternalStreamUtilsGet(loop,state,"autoDestroy"))&&tsTruthy(tsInternalStreamUtilsGet(loop,state,"emitClose"))&&tsStrictEqual(tsInternalStreamUtilsGet(loop,state,"closed"),tsBooleanValue(false)))
}

func tsInternalStreamIsDisturbed(loop *tsLoop,stream tsValue) tsValue {
    shared:=tsInternalStreamUtilsKey(loop,"nodejs.stream.disturbed")
    value:=tsInternalStreamUtilsGetKey(loop,stream,shared)
    if tsNullish(value){value=tsInternalStreamUtilsGet(loop,stream,"readableDidRead");if !tsTruthy(value){value=tsInternalStreamUtilsGet(loop,stream,"readableAborted")}}
    return tsBooleanValue(tsTruthy(value))
}

func tsInternalStreamIsErrored(loop *tsLoop,stream tsValue) tsValue {
    if tsNullish(stream){return tsBooleanValue(false)}
    value:=tsInternalStreamUtilsGetKey(loop,stream,tsInternalStreamUtilsKey(loop,"nodejs.stream.errored"))
    for _,name:=range []string{"readableErrored","writableErrored"}{if tsNullish(value){value=tsInternalStreamUtilsGet(loop,stream,name)}}
    for _,stateName:=range []string{"_readableState","_writableState"}{state:=tsInternalStreamUtilsGet(loop,stream,stateName);for _,name:=range []string{"errorEmitted","errored"}{if tsNullish(value){value=tsInternalStreamUtilsGet(loop,state,name)}}}
    return tsBooleanValue(tsTruthy(value))
}

func tsInternalStreamUtilsModule(loop *tsLoop) tsValue {
    tsInternalStreamUtilsInit(loop)
    if tsInternalStreamUtils.module!=nil{return tsObjectValue(tsInternalStreamUtils.module)}
    module:=tsNewObject()
    module.set("kOnConstructed",tsInternalStreamUtilsKey(loop,"kOnConstructed"))
    for _,name:=range []string{"kIsDestroyed","kIsErrored","kIsReadable","kIsWritable","kIsDisturbed","kIsClosedPromise","kControllerErrorFunction","kState","kObjectMode","kErrorEmitted","kAutoDestroy","kEmitClose","kDestroyed","kClosed","kCloseEmitted","kErrored","kConstructed"}{module.set(name,tsInternalStreamUtilsKey(loop,name))}
    predicates:=[]string{"isDestroyed","isDisturbed","isErrored","isReadable","isClosed","isDuplexNodeStream","isFinished","isIterable","isReadableNodeStream","isReadableStream","isReadableFinished","isReadableErrored","isNodeStream","isWebStream","isWritable","isWritableNodeStream","isWritableStream","isWritableEnded","isWritableFinished","isWritableErrored","isServerRequest","willEmitClose","isTransformStream"}
    for _,name:=range predicates {
        method:=name
        module.set(name,tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{
            stream:=tsArg(args,0)
            switch method {
            case "isDestroyed":return tsInternalStreamIsDestroyed(loop,stream)
            case "isDisturbed":return tsInternalStreamIsDisturbed(loop,stream)
            case "isErrored":return tsInternalStreamIsErrored(loop,stream)
            case "isReadable":return tsInternalStreamIsReadable(loop,stream)
            case "isClosed":return tsInternalStreamIsClosed(loop,stream)
            case "isDuplexNodeStream":return tsBooleanValue(tsInternalStreamIsDuplex(loop,stream))
            case "isFinished":return tsInternalStreamIsFinished(loop,stream,tsArg(args,1))
            case "isIterable":return tsInternalStreamIsIterable(loop,stream,tsArg(args,1))
            case "isReadableNodeStream":return tsBooleanValue(tsInternalStreamIsReadableNode(loop,stream,tsTruthy(tsArg(args,1))))
            case "isReadableStream":return tsBooleanValue(tsInternalStreamIsWebClass(stream,"ReadableStream"))
            case "isReadableFinished":return tsInternalStreamIsReadableFinished(loop,stream,tsArg(args,1))
            case "isReadableErrored":return tsInternalStreamErrored(loop,stream,false)
            case "isNodeStream":return tsBooleanValue(tsInternalStreamIsNode(loop,stream))
            case "isWebStream":return tsBooleanValue(tsInternalStreamIsWeb(loop,stream))
            case "isWritable":return tsInternalStreamIsWritable(loop,stream)
            case "isWritableNodeStream":return tsBooleanValue(tsInternalStreamIsWritableNode(loop,stream))
            case "isWritableStream":return tsBooleanValue(tsInternalStreamIsWebClass(stream,"WritableStream"))
            case "isWritableEnded":return tsInternalStreamIsWritableEnded(loop,stream)
            case "isWritableFinished":return tsInternalStreamIsWritableFinished(loop,stream,tsArg(args,1))
            case "isWritableErrored":return tsInternalStreamErrored(loop,stream,true)
            case "isServerRequest":return tsBooleanValue(tsInternalStreamIsServerRequest(loop,stream))
            case "willEmitClose":return tsInternalStreamWillEmitClose(loop,stream)
            case "isTransformStream":return tsBooleanValue(tsInternalStreamIsWebClass(stream,"TransformStream"))
            }
            return tsU
        }))
    }
    tsInternalStreamUtils.module=module
    return tsObjectValue(module)
}
`
