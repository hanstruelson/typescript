package goemit

// internalStreamDestroyRuntime maps the Node stream destroy lifecycle, state
// flags, callback scheduling, and legacy destroyer helper onto local emitters.
const internalStreamDestroyRuntime = `
var tsInternalStreamDestroyState struct {
    sync.Once
    destroySymbol tsValue
    constructSymbol tsValue
}

func tsInternalStreamDestroyInit() {
    tsInternalStreamDestroyState.Do(func() {
        tsInternalStreamDestroyState.destroySymbol = tsNewSymbol(tsStringUTF8("kDestroy"))
        tsInternalStreamDestroyState.constructSymbol = tsNewSymbol(tsStringUTF8("kConstruct"))
    })
}

func tsInternalStreamDestroyStateValue(loop *tsLoop, stream tsValue, writable bool) tsValue {
    name := "_readableState"
    if writable { name = "_writableState" }
    return tsGet(loop, stream, tsStringReference(tsStringUTF8(name)))
}

func tsInternalStreamDestroyBits(loop *tsLoop, stream, state tsValue, bitName string) uint32 {
    if tsNullish(state) { return 0 }
    bits := tsGet(loop, state, tsInternalStreamUtilsKey(loop, "kState"))
    if tsIsNumeric(bits) { return uint32(tsNumber(bits)) }
    property := map[string]string{"kObjectMode":"objectMode", "kErrorEmitted":"errorEmitted", "kAutoDestroy":"autoDestroy", "kEmitClose":"emitClose", "kDestroyed":"destroyed", "kClosed":"closed", "kCloseEmitted":"closeEmitted", "kConstructed":"constructed"}[bitName]
    if property != "" && tsTruthy(tsGet(loop, state, tsStringReference(tsStringUTF8(property)))) { return uint32(tsNumber(tsInternalStreamUtilsKey(loop, bitName))) }
    if bitName == "kErrored" && !tsNullish(tsGet(loop, state, tsStringReference(tsStringUTF8("errored")))) { return uint32(tsNumber(tsInternalStreamUtilsKey(loop, bitName))) }
    if tsTruthy(tsGet(loop, stream, tsStringReference(tsStringUTF8("destroyed")))) && bitName == "kDestroyed" { return uint32(tsNumber(tsInternalStreamUtilsKey(loop, bitName))) }
    return 0
}

func tsInternalStreamDestroySetBits(loop *tsLoop, state tsValue, bitNames ...string) {
    if tsNullish(state) { return }
    key := tsInternalStreamUtilsKey(loop, "kState")
    bitsValue := tsGet(loop, state, key)
    bits := uint32(0)
    if tsIsNumeric(bitsValue) { bits = uint32(tsNumber(bitsValue)) }
    for _, name := range bitNames {
        bits |= uint32(tsNumber(tsInternalStreamUtilsKey(loop, name)))
        property := map[string]string{"kObjectMode":"objectMode", "kErrorEmitted":"errorEmitted", "kAutoDestroy":"autoDestroy", "kEmitClose":"emitClose", "kDestroyed":"destroyed", "kClosed":"closed", "kCloseEmitted":"closeEmitted", "kConstructed":"constructed"}[name]
        if property != "" { tsSet(loop, state, tsStringReference(tsStringUTF8(property)), tsBooleanValue(true)) }
    }
    tsSet(loop, state, key, tsNumberValue(float64(bits)))
}

func tsInternalStreamDestroyClearBits(loop *tsLoop, state tsValue, bitNames ...string) {
    if tsNullish(state) { return }
    key := tsInternalStreamUtilsKey(loop, "kState")
    bitsValue := tsGet(loop, state, key)
    bits := uint32(0)
    if tsIsNumeric(bitsValue) { bits = uint32(tsNumber(bitsValue)) }
    for _, name := range bitNames {
        bits &^= uint32(tsNumber(tsInternalStreamUtilsKey(loop, name)))
        property := map[string]string{"kObjectMode":"objectMode", "kErrorEmitted":"errorEmitted", "kAutoDestroy":"autoDestroy", "kEmitClose":"emitClose", "kDestroyed":"destroyed", "kClosed":"closed", "kCloseEmitted":"closeEmitted", "kConstructed":"constructed"}[name]
        if property != "" { tsSet(loop, state, tsStringReference(tsStringUTF8(property)), tsBooleanValue(false)) }
    }
    tsSet(loop, state, key, tsNumberValue(float64(bits)))
}

func tsInternalStreamDestroyCheckError(loop *tsLoop, err, writable, readable tsValue) {
    if !tsTruthy(err) { return }
    for _, state := range []tsValue{writable, readable} {
        if tsNullish(state) { continue }
        if !tsTruthy(tsGet(loop, state, tsStringReference(tsStringUTF8("errored")))) {
            tsSet(loop, state, tsStringReference(tsStringUTF8("errored")), err)
            tsInternalStreamDestroySetBits(loop, state, "kErrored")
        }
    }
}

func tsInternalStreamDestroyNextTick(loop *tsLoop, callback tsValue, args ...tsValue) {
    process := tsNodeProcessModule(loop)
    nextTick := tsGet(loop, process, tsStringReference(tsStringUTF8("nextTick")))
    values := append([]tsValue{callback}, args...)
    tsCallReceiver(loop, nextTick, process, values...)
}

func tsInternalStreamDestroyEmit(loop *tsLoop, stream tsValue, event tsValue, args ...tsValue) {
    emit := tsGet(loop, stream, tsStringReference(tsStringUTF8("emit")))
    if emit.kind == tsFunctionKind { tsCallReceiver(loop, emit, stream, append([]tsValue{event}, args...)...) }
}

func tsInternalStreamDestroyEmitErrorNT(loop *tsLoop, stream, err tsValue) tsValue {
    readable := tsInternalStreamDestroyStateValue(loop, stream, false)
    writable := tsInternalStreamDestroyStateValue(loop, stream, true)
    mask := uint32(tsNumber(tsInternalStreamUtilsKey(loop, "kErrorEmitted")))
    if tsInternalStreamDestroyBits(loop, stream, readable, "kErrorEmitted")&mask != 0 || tsInternalStreamDestroyBits(loop, stream, writable, "kErrorEmitted")&mask != 0 { return tsU }
    tsInternalStreamDestroySetBits(loop, readable, "kErrorEmitted")
    tsInternalStreamDestroySetBits(loop, writable, "kErrorEmitted")
    tsInternalStreamDestroyEmit(loop, stream, tsStringReference(tsStringUTF8("error")), err)
    return tsU
}

func tsInternalStreamDestroyEmitCloseNT(loop *tsLoop, stream tsValue) tsValue {
    readable := tsInternalStreamDestroyStateValue(loop, stream, false)
    writable := tsInternalStreamDestroyStateValue(loop, stream, true)
    tsInternalStreamDestroySetBits(loop, writable, "kCloseEmitted")
    tsInternalStreamDestroySetBits(loop, readable, "kCloseEmitted")
    emitMask := uint32(tsNumber(tsInternalStreamUtilsKey(loop, "kEmitClose")))
    emit := tsInternalStreamDestroyBits(loop, stream, writable, "kEmitClose")&emitMask != 0 || tsInternalStreamDestroyBits(loop, stream, readable, "kEmitClose")&emitMask != 0
    if emit { tsInternalStreamDestroyEmit(loop, stream, tsStringReference(tsStringUTF8("close"))) }
    return tsU
}

func tsInternalStreamDestroyEmitErrorCloseNT(loop *tsLoop, stream, err tsValue) tsValue {
    tsInternalStreamDestroyEmitErrorNT(loop, stream, err)
    tsInternalStreamDestroyEmitCloseNT(loop, stream)
    return tsU
}

func tsInternalStreamDestroy(loop *tsLoop, receiver, err, callback tsValue) tsValue {
    writable := tsInternalStreamDestroyStateValue(loop, receiver, true)
    readable := tsInternalStreamDestroyStateValue(loop, receiver, false)
    destroyedMask := uint32(tsNumber(tsInternalStreamUtilsKey(loop, "kDestroyed")))
    if (!tsNullish(writable) && tsInternalStreamDestroyBits(loop, receiver, writable, "kDestroyed")&destroyedMask != 0) || (!tsNullish(readable) && tsInternalStreamDestroyBits(loop, receiver, readable, "kDestroyed")&destroyedMask != 0) {
        if callback.kind == tsFunctionKind { tsCall(loop, callback) }
        return receiver
    }
    tsInternalStreamDestroyCheckError(loop, err, writable, readable)
    tsInternalStreamDestroySetBits(loop, writable, "kDestroyed")
    tsInternalStreamDestroySetBits(loop, readable, "kDestroyed")
    state := writable
    if tsNullish(state) { state = readable }
    constructedMask := uint32(tsNumber(tsInternalStreamUtilsKey(loop, "kConstructed")))
    if !tsNullish(state) && tsInternalStreamDestroyBits(loop, receiver, state, "kConstructed")&constructedMask == 0 {
        tsInternalStreamDestroyInit()
        listener := tsFunctionValue(tsFunc(func(loop *tsLoop, args ...tsValue) tsValue {
            errors := tsInternalErrorsModule(loop)
            aggregate := tsGet(loop, errors, tsStringReference(tsStringUTF8("aggregateTwoErrors")))
            combined := tsCallReceiver(loop, aggregate, errors, tsArg(args, 0), err)
            tsInternalStreamDestroyDo(loop, receiver, combined, callback)
            return tsU
        }))
        once := tsGet(loop, receiver, tsStringReference(tsStringUTF8("once")))
        if once.kind == tsFunctionKind { tsCallReceiver(loop, once, receiver, tsInternalStreamDestroyState.destroySymbol, listener) }
    } else {
        tsInternalStreamDestroyDo(loop, receiver, err, callback)
    }
    return receiver
}

func tsInternalStreamDestroyDo(loop *tsLoop, stream, err, callback tsValue) {
    called := false
    onDestroy := tsFunctionValue(tsFunc(func(loop *tsLoop, args ...tsValue) tsValue {
        if called { return tsU }
        called = true
        errorValue := tsArg(args, 0)
        readable := tsInternalStreamDestroyStateValue(loop, stream, false)
        writable := tsInternalStreamDestroyStateValue(loop, stream, true)
        tsInternalStreamDestroyCheckError(loop, errorValue, writable, readable)
        tsInternalStreamDestroySetBits(loop, writable, "kClosed")
        tsInternalStreamDestroySetBits(loop, readable, "kClosed")
        if callback.kind == tsFunctionKind { tsCall(loop, callback, errorValue) }
        if tsTruthy(errorValue) {
            eventCallback := tsFunctionValue(tsFunc(func(loop *tsLoop, _ ...tsValue) tsValue { return tsInternalStreamDestroyEmitErrorCloseNT(loop, stream, errorValue) }))
            tsInternalStreamDestroyNextTick(loop, eventCallback)
        } else {
            eventCallback := tsFunctionValue(tsFunc(func(loop *tsLoop, _ ...tsValue) tsValue { return tsInternalStreamDestroyEmitCloseNT(loop, stream) }))
            tsInternalStreamDestroyNextTick(loop, eventCallback)
        }
        return tsU
    }))
    destroy := tsGet(loop, stream, tsStringReference(tsStringUTF8("_destroy")))
    if destroy.kind != tsFunctionKind { tsInternalStreamDestroyEmitCloseNT(loop, stream); return }
    destroyError := err
    if !tsTruthy(destroyError) { destroyError = tsNull }
    defer func() { if failure := recover(); failure != nil { tsCall(loop, onDestroy, tsUnwrap(failure)) } }()
    tsCallReceiver(loop, destroy, stream, destroyError, onDestroy)
}

func tsInternalStreamUndestroy(loop *tsLoop, stream tsValue) tsValue {
    readable := tsInternalStreamDestroyStateValue(loop, stream, false)
    writable := tsInternalStreamDestroyStateValue(loop, stream, true)
    if !tsNullish(readable) {
        for name, value := range map[string]tsValue{"constructed": tsBooleanValue(true), "closed": tsBooleanValue(false), "closeEmitted": tsBooleanValue(false), "destroyed": tsBooleanValue(false), "errored": tsNull, "errorEmitted": tsBooleanValue(false), "reading": tsBooleanValue(false)} { tsSet(loop, readable, tsStringReference(tsStringUTF8(name)), value) }
        ended := tsStrictEqual(tsGet(loop, readable, tsStringReference(tsStringUTF8("readable"))), tsBooleanValue(false))
        tsSet(loop, readable, tsStringReference(tsStringUTF8("ended")), tsBooleanValue(ended))
        tsSet(loop, readable, tsStringReference(tsStringUTF8("endEmitted")), tsBooleanValue(ended))
    }
    if !tsNullish(writable) {
        for name, value := range map[string]tsValue{"constructed": tsBooleanValue(true), "destroyed": tsBooleanValue(false), "closed": tsBooleanValue(false), "closeEmitted": tsBooleanValue(false), "errored": tsNull, "errorEmitted": tsBooleanValue(false), "finalCalled": tsBooleanValue(false), "prefinished": tsBooleanValue(false)} { tsSet(loop, writable, tsStringReference(tsStringUTF8(name)), value) }
        ended := tsStrictEqual(tsGet(loop, writable, tsStringReference(tsStringUTF8("writable"))), tsBooleanValue(false))
        tsSet(loop, writable, tsStringReference(tsStringUTF8("ended")), tsBooleanValue(ended))
        tsSet(loop, writable, tsStringReference(tsStringUTF8("ending")), tsBooleanValue(ended))
        tsSet(loop, writable, tsStringReference(tsStringUTF8("finished")), tsBooleanValue(ended))
    }
    tsInternalStreamDestroyClearBits(loop, readable, "kDestroyed", "kClosed", "kCloseEmitted", "kErrored", "kErrorEmitted", "kConstructed")
    tsInternalStreamDestroyClearBits(loop, writable, "kDestroyed", "kClosed", "kCloseEmitted", "kErrored", "kErrorEmitted", "kConstructed")
    tsInternalStreamDestroySetBits(loop, readable, "kConstructed")
    tsInternalStreamDestroySetBits(loop, writable, "kConstructed")
    return tsU
}

func tsInternalStreamErrorOrDestroy(loop *tsLoop, stream, err, syncValue tsValue) tsValue {
    readable := tsInternalStreamDestroyStateValue(loop, stream, false)
    writable := tsInternalStreamDestroyStateValue(loop, stream, true)
    destroyedMask := uint32(tsNumber(tsInternalStreamUtilsKey(loop, "kDestroyed")))
    if (!tsNullish(writable) && tsInternalStreamDestroyBits(loop, stream, writable, "kDestroyed")&destroyedMask != 0) || (!tsNullish(readable) && tsInternalStreamDestroyBits(loop, stream, readable, "kDestroyed")&destroyedMask != 0) { return stream }
    autoMask := uint32(tsNumber(tsInternalStreamUtilsKey(loop, "kAutoDestroy")))
    if (!tsNullish(readable) && tsInternalStreamDestroyBits(loop, stream, readable, "kAutoDestroy")&autoMask != 0) || (!tsNullish(writable) && tsInternalStreamDestroyBits(loop, stream, writable, "kAutoDestroy")&autoMask != 0) {
        destroy := tsGet(loop, stream, tsStringReference(tsStringUTF8("destroy")))
        if destroy.kind == tsFunctionKind { tsCallReceiver(loop, destroy, stream, err) }
    } else if tsTruthy(err) {
        tsInternalStreamDestroyCheckError(loop, err, writable, readable)
        if tsTruthy(syncValue) {
            callback := tsFunctionValue(tsFunc(func(loop *tsLoop, _ ...tsValue) tsValue { return tsInternalStreamDestroyEmitErrorNT(loop, stream, err) }))
            tsInternalStreamDestroyNextTick(loop, callback)
        } else { tsInternalStreamDestroyEmitErrorNT(loop, stream, err) }
    }
    return stream
}

func tsInternalStreamConstruct(loop *tsLoop, stream, callback tsValue) tsValue {
    constructMethod := tsGet(loop, stream, tsStringReference(tsStringUTF8("_construct")))
    if constructMethod.kind != tsFunctionKind { return tsU }
    readable := tsInternalStreamDestroyStateValue(loop, stream, false)
    writable := tsInternalStreamDestroyStateValue(loop, stream, true)
    tsInternalStreamDestroyClearBits(loop, readable, "kConstructed")
    tsInternalStreamDestroyClearBits(loop, writable, "kConstructed")
    once := tsGet(loop, stream, tsStringReference(tsStringUTF8("once")))
    tsInternalStreamDestroyInit()
    if once.kind == tsFunctionKind { tsCallReceiver(loop, once, stream, tsInternalStreamDestroyState.constructSymbol, callback) }
    count := tsGet(loop, stream, tsStringReference(tsStringUTF8("listenerCount")))
    if count.kind == tsFunctionKind && tsNumber(tsCallReceiver(loop, count, stream, tsInternalStreamDestroyState.constructSymbol)) > 1 { return tsU }
    task := tsFunctionValue(tsFunc(func(loop *tsLoop, _ ...tsValue) tsValue { return tsInternalStreamConstructNT(loop, stream) }))
    tsInternalStreamDestroyNextTick(loop, task)
    return tsU
}

func tsInternalStreamConstructNT(loop *tsLoop, stream tsValue) tsValue {
    called := false
    onConstruct := tsFunctionValue(tsFunc(func(loop *tsLoop, args ...tsValue) tsValue {
        err := tsArg(args, 0)
        if called {
            if tsNullish(err) { err = tsFSException("Error", "ERR_MULTIPLE_CALLBACK", "Callback called multiple times") }
            return tsInternalStreamErrorOrDestroy(loop, stream, err, tsU)
        }
        called = true
        readable := tsInternalStreamDestroyStateValue(loop, stream, false)
        writable := tsInternalStreamDestroyStateValue(loop, stream, true)
        tsInternalStreamDestroySetBits(loop, readable, "kConstructed")
        tsInternalStreamDestroySetBits(loop, writable, "kConstructed")
        state := writable
        if tsNullish(state) { state = readable }
        if tsTruthy(tsGet(loop, state, tsStringReference(tsStringUTF8("destroyed")))) {
            tsInternalStreamDestroyEmit(loop, stream, tsInternalStreamDestroyState.constructSymbol, err)
        } else if tsTruthy(err) {
            tsInternalStreamErrorOrDestroy(loop, stream, err, tsBooleanValue(true))
        } else {
            tsInternalStreamDestroyEmit(loop, stream, tsInternalStreamDestroyState.constructSymbol)
        }
        return tsU
    }))
    constructMethod := tsGet(loop, stream, tsStringReference(tsStringUTF8("_construct")))
    if constructMethod.kind != tsFunctionKind { return tsU }
    callback := tsFunctionValue(tsFunc(func(loop *tsLoop, args ...tsValue) tsValue {
        err := tsArg(args, 0)
        task := tsFunctionValue(tsFunc(func(loop *tsLoop, _ ...tsValue) tsValue { return tsCall(loop, onConstruct, err) }))
        tsInternalStreamDestroyNextTick(loop, task)
        return tsU
    }))
    defer func() {
        if failure := recover(); failure != nil {
            err := tsUnwrap(failure)
            task := tsFunctionValue(tsFunc(func(loop *tsLoop, _ ...tsValue) tsValue { return tsCall(loop, onConstruct, err) }))
            tsInternalStreamDestroyNextTick(loop, task)
        }
    }()
    tsCallReceiver(loop, constructMethod, stream, callback)
    return tsU
}

func tsInternalStreamIsRequest(loop *tsLoop, stream tsValue) bool {
    return !tsNullish(stream) && tsTruthy(tsGet(loop, stream, tsStringReference(tsStringUTF8("setHeader")))) && tsGet(loop, stream, tsStringReference(tsStringUTF8("abort"))).kind == tsFunctionKind
}

func tsInternalStreamEmitCloseLegacy(loop *tsLoop, stream tsValue) tsValue {
    tsInternalStreamDestroyEmit(loop, stream, tsStringReference(tsStringUTF8("close")))
    return tsU
}

func tsInternalStreamEmitErrorCloseLegacy(loop *tsLoop, stream, err tsValue) tsValue {
    tsInternalStreamDestroyEmit(loop, stream, tsStringReference(tsStringUTF8("error")), err)
    task := tsFunctionValue(tsFunc(func(loop *tsLoop, _ ...tsValue) tsValue { return tsInternalStreamEmitCloseLegacy(loop, stream) }))
    tsInternalStreamDestroyNextTick(loop, task)
    return tsU
}

func tsInternalStreamDestroyer(loop *tsLoop, stream, err tsValue) tsValue {
    if tsNullish(stream) || tsTruthy(tsInternalStreamIsDestroyed(loop, stream)) { return tsU }
    if !tsTruthy(err) && !tsTruthy(tsInternalStreamIsFinished(loop, stream, tsU)) {
        err = tsFSException("AbortError", "ABORT_ERR", "The operation was aborted")
    }
    if tsTruthy(tsInternalStreamIsServerRequest(loop, stream)) {
        tsSet(loop, stream, tsStringReference(tsStringUTF8("socket")), tsNull)
        destroy := tsGet(loop, stream, tsStringReference(tsStringUTF8("destroy")))
        if destroy.kind == tsFunctionKind { tsCallReceiver(loop, destroy, stream, err) }
    } else if tsInternalStreamIsRequest(loop, stream) {
        abort := tsGet(loop, stream, tsStringReference(tsStringUTF8("abort")))
        tsCallReceiver(loop, abort, stream)
    } else {
        request := tsGet(loop, stream, tsStringReference(tsStringUTF8("req")))
        if tsInternalStreamIsRequest(loop, request) {
            abort := tsGet(loop, request, tsStringReference(tsStringUTF8("abort")))
            tsCallReceiver(loop, abort, request)
        } else if destroy := tsGet(loop, stream, tsStringReference(tsStringUTF8("destroy"))); destroy.kind == tsFunctionKind {
            tsCallReceiver(loop, destroy, stream, err)
        } else if close := tsGet(loop, stream, tsStringReference(tsStringUTF8("close"))); close.kind == tsFunctionKind {
            tsCallReceiver(loop, close, stream)
        } else if tsTruthy(err) {
            task := tsFunctionValue(tsFunc(func(loop *tsLoop, _ ...tsValue) tsValue { return tsInternalStreamEmitErrorCloseLegacy(loop, stream, err) }))
            tsInternalStreamDestroyNextTick(loop, task)
        } else {
            task := tsFunctionValue(tsFunc(func(loop *tsLoop, _ ...tsValue) tsValue { return tsInternalStreamEmitCloseLegacy(loop, stream) }))
            tsInternalStreamDestroyNextTick(loop, task)
        }
    }
    if !tsTruthy(tsGet(loop, stream, tsStringReference(tsStringUTF8("destroyed")))) {
        destroyedKey := tsInternalStreamUtilsKey(loop, "nodejs.stream.destroyed")
        tsSet(loop, stream, destroyedKey, tsBooleanValue(true))
    }
    return tsU
}

func tsInternalStreamDestroyModule(loop *tsLoop) tsValue {
    tsInternalStreamDestroyInit()
    module := tsNewObject()
    module.set("construct", tsNamedNativeMethod("construct", 2, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue { return tsInternalStreamConstruct(loop, tsArg(args, 0), tsArg(args, 1)) }))
    module.set("destroyer", tsNamedNativeMethod("destroyer", 2, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue { return tsInternalStreamDestroyer(loop, tsArg(args, 0), tsArg(args, 1)) }))
    module.set("destroy", tsNamedNativeMethod("destroy", 2, func(loop *tsLoop, receiver tsValue, args ...tsValue) tsValue { return tsInternalStreamDestroy(loop, receiver, tsArg(args, 0), tsArg(args, 1)) }))
    module.set("undestroy", tsNamedNativeMethod("undestroy", 0, func(loop *tsLoop, receiver tsValue, _ ...tsValue) tsValue { return tsInternalStreamUndestroy(loop, receiver) }))
    module.set("errorOrDestroy", tsNamedNativeMethod("errorOrDestroy", 3, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue { return tsInternalStreamErrorOrDestroy(loop, tsArg(args, 0), tsArg(args, 1), tsArg(args, 2)) }))
    module.set("default", tsObjectValue(module))
    return tsObjectValue(module)
}
`
