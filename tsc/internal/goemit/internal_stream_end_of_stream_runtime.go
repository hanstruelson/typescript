package goemit

// internalStreamEndOfStreamRuntime maps Node's EOS listener state machine onto
// the local stream predicates, event methods, promises, and AbortSignal values.
const internalStreamEndOfStreamRuntime = `
var tsInternalEosState struct {
    sync.Once
    module tsValue
    syncCallback tsValue
}

func tsInternalEosError(code, message string) tsValue {
    return tsFSException("Error", code, message)
}

func tsInternalEosAbortError(loop *tsLoop, signal tsValue) tsValue {
    reason := tsGet(loop, signal, tsStringReference(tsStringUTF8("reason")))
    errorValue := tsFSException("AbortError", "ABORT_ERR", "The operation was aborted")
    if errorValue.kind == tsErrorKind {
        (*tsRuntimeError)(errorValue.ref).fields["cause"] = reason
    }
    return errorValue
}

func tsInternalEosGetErrored(loop *tsLoop, stream tsValue) tsValue {
    writable := tsInternalStreamErrored(loop, stream, true)
    if !tsNullish(writable) && writable.kind != tsBooleanKind && tsTruthy(writable) {
        return writable
    }
    readable := tsInternalStreamErrored(loop, stream, false)
    if !tsNullish(readable) && readable.kind != tsBooleanKind && tsTruthy(readable) {
        return readable
    }
    return tsNull
}

func tsInternalEosOnCloseError(loop *tsLoop, stream tsValue, readable, readableFinished, writable, writableFinished bool) tsValue {
    if errorValue := tsInternalEosGetErrored(loop, stream); !tsNullish(errorValue) {
        return errorValue
    }
    if readable && !readableFinished && tsInternalStreamIsReadableNode(loop, stream, true) && !tsTruthy(tsInternalStreamIsReadableFinished(loop, stream, tsBooleanValue(false))) {
        return tsInternalEosError("ERR_STREAM_PREMATURE_CLOSE", "Premature close")
    }
    if writable && !writableFinished && !tsTruthy(tsInternalStreamIsWritableFinished(loop, stream, tsBooleanValue(false))) {
        return tsInternalEosError("ERR_STREAM_PREMATURE_CLOSE", "Premature close")
    }
    return tsNull
}

func tsInternalEosHasAsyncContext(loop *tsLoop) bool {
    if loop.asyncContext != nil {
        return true
    }
    hooks := tsInternalAsyncHooksModule()
    enabled := tsGet(loop, hooks, tsStringReference(tsStringUTF8("enabledHooksExist")))
    return tsTruthy(tsCallReceiver(loop, enabled, hooks))
}

func tsInternalEosBindAsyncResource(loop *tsLoop, callback tsValue) tsValue {
    hooks := tsNodeAsyncHooksModule()
    constructor := tsGet(loop, hooks, tsStringReference(tsStringUTF8("AsyncResource")))
    resource := tsConstruct(constructor, tsStringReference(tsStringUTF8("STREAM_END_OF_STREAM")))
    run := tsGet(loop, resource, tsStringReference(tsStringUTF8("runInAsyncScope")))
    return tsNativeMethod(func(loop *tsLoop, receiver tsValue, args ...tsValue) tsValue {
        callArgs := append([]tsValue{callback, receiver}, args...)
        return tsCallReceiver(loop, run, resource, callArgs...)
    })
}

func tsInternalEosNextTick(loop *tsLoop, callback tsValue) {
    process := tsNodeProcessModule(loop)
    nextTick := tsGet(loop, process, tsStringReference(tsStringUTF8("nextTick")))
    tsCallReceiver(loop, nextTick, process, callback)
}

func tsInternalEosCall(loop *tsLoop, stream tsValue, name string, args ...tsValue) {
    method := tsGet(loop, stream, tsStringReference(tsStringUTF8(name)))
    if method.kind == tsFunctionKind {
        tsCallReceiver(loop, method, stream, args...)
    }
}

func tsInternalEosOnce(callback tsValue) tsValue {
    called := false
    return tsNativeMethod(func(loop *tsLoop, receiver tsValue, args ...tsValue) tsValue {
        if called {
            return tsU
        }
        called = true
        return tsCallReceiver(loop, callback, receiver, args...)
    })
}

func tsInternalEosAttachAbortListener(loop *tsLoop, signal, callback tsValue) tsValue {
    module := tsInternalAbortListenerModule(loop)
    add := tsGet(loop, module, tsStringReference(tsStringUTF8("addAbortListener")))
    return tsCallReceiver(loop, add, module, signal, callback)
}

func tsInternalEosDispose(loop *tsLoop, disposable tsValue) {
    dispose := tsGet(loop, disposable, tsWellKnownSymbol("dispose"))
    if dispose.kind == tsFunctionKind {
        tsCallReceiver(loop, dispose, disposable)
    }
}

func tsInternalEosWebClosedPromise(loop *tsLoop, stream tsValue) tsValue {
    closedKey := tsInternalStreamUtilsKey(loop, "nodejs.webstream.isClosedPromise")
    closedState := tsGet(loop, stream, closedKey)
    if !tsNullish(closedState) {
        promise := tsGet(loop, closedState, tsStringReference(tsStringUTF8("promise")))
        if promise.kind == tsPromiseKind {
            return promise
        }
    }
    if stream.kind == tsObjectKind && (*tsObject)(stream.ref).webStream != nil {
        return tsPromiseValue(tsWebStreamClosedPromise(loop, (*tsObject)(stream.ref).webStream))
    }
    promise := tsGet(loop, stream, tsStringReference(tsStringUTF8("closed")))
    if promise.kind == tsPromiseKind {
        return promise
    }
    return tsPromiseValue(loop.promise())
}

func tsInternalEosWeb(loop *tsLoop, stream, options, callback tsValue) tsValue {
    if tsInternalEosHasAsyncContext(loop) {
        callback = tsInternalEosOnce(tsInternalEosBindAsyncResource(loop, callback))
    } else {
        callback = tsInternalEosOnce(callback)
    }
    isAborted := false
    signal := tsGet(loop, options, tsStringReference(tsStringUTF8("signal")))
    var disposable tsValue
    if !tsNullish(signal) {
        abort := tsNativeMethod(func(loop *tsLoop, _ tsValue, _ ...tsValue) tsValue {
            isAborted = true
            tsCallReceiver(loop, callback, stream, tsInternalEosAbortError(loop, signal))
            return tsU
        })
        if tsTruthy(tsGet(loop, signal, tsStringReference(tsStringUTF8("aborted")))) {
            abortValue := abort
            tsInternalEosNextTick(loop, tsFunctionValue(tsFunc(func(loop *tsLoop, _ ...tsValue) tsValue {
                return tsCallReceiver(loop, abortValue, tsU)
            })))
        } else {
            disposable = tsInternalEosAttachAbortListener(loop, signal, abort)
            original := callback
            callback = tsInternalEosOnce(tsNativeMethod(func(loop *tsLoop, receiver tsValue, args ...tsValue) tsValue {
                tsInternalEosDispose(loop, disposable)
                return tsCallReceiver(loop, original, receiver, args...)
            }))
        }
    }
    resolver := func(rejected bool, value tsValue) {
        if isAborted {
            return
        }
        var invoke tsValue
        if rejected {
            invoke = tsFunctionValue(tsFunc(func(loop *tsLoop, _ ...tsValue) tsValue {
                return tsCallReceiver(loop, callback, stream, value)
            }))
        } else {
            invoke = tsFunctionValue(tsFunc(func(loop *tsLoop, _ ...tsValue) tsValue {
                return tsCallReceiver(loop, callback, stream)
            }))
        }
        tsInternalEosNextTick(loop, invoke)
    }
    closed := tsInternalEosWebClosedPromise(loop, stream)
    if closed.kind == tsPromiseKind {
        (*tsPromise)(closed.ref).then(func(result tsResult) { resolver(result.rejected, result.value) })
    } else {
        then := tsGet(loop, closed, tsStringReference(tsStringUTF8("then")))
        fulfilled := tsNamedNativeMethod("eosWebFulfilled", 1, func(_ *tsLoop, _ tsValue, args ...tsValue) tsValue { resolver(false, tsArg(args, 0)); return tsU })
        rejected := tsNamedNativeMethod("eosWebRejected", 1, func(_ *tsLoop, _ tsValue, args ...tsValue) tsValue { resolver(true, tsArg(args, 0)); return tsU })
        if then.kind == tsFunctionKind {
            tsCallReceiver(loop, then, closed, fulfilled, rejected)
        }
    }
    return tsFunctionValue(tsFunc(func(_ *tsLoop, _ ...tsValue) tsValue { return tsU }))
}

func tsInternalEos(loop *tsLoop, args []tsValue) tsValue {
    stream := tsArg(args, 0)
    options := tsArg(args, 1)
    callback := tsArg(args, 2)
    if len(args) == 2 {
        callback = options
        options = tsInternalSharedSingleton.emptyObject
    } else if tsNullish(options) {
        options = tsInternalSharedSingleton.emptyObject
    } else {
        tsInternalValidatorObject(options, tsStringReference(tsStringUTF8("options")), tsNumberValue(0))
    }
    if callback.kind != tsFunctionKind {
        tsInternalValidatorInvalidType("callback", "of type function", callback)
    }
    signal := tsGet(loop, options, tsStringReference(tsStringUTF8("signal")))
    tsInternalValidatorAbortSignal(signal, tsStringReference(tsStringUTF8("options.signal")))

    if tsInternalStreamIsWebClass(stream, "ReadableStream") || tsInternalStreamIsWebClass(stream, "WritableStream") {
        return tsInternalEosWeb(loop, stream, options, callback)
    }
    if !tsInternalStreamIsNode(loop, stream) {
        tsInternalValidatorInvalidType("stream", "an instance of ReadableStream, WritableStream, or Stream", stream)
    }

    readableValue := tsGet(loop, options, tsStringReference(tsStringUTF8("readable")))
    writableValue := tsGet(loop, options, tsStringReference(tsStringUTF8("writable")))
    if tsIsUndefined(readableValue) || tsIsNull(readableValue) { readableValue = tsBooleanValue(tsInternalStreamIsReadableNode(loop, stream, false)) }
    if tsIsUndefined(writableValue) || tsIsNull(writableValue) { writableValue = tsBooleanValue(tsInternalStreamIsWritableNode(loop, stream)) }
    readable := tsTruthy(readableValue)
    writable := tsTruthy(writableValue)
    willEmitClose := tsTruthy(tsInternalStreamWillEmitClose(loop, stream)) && tsStrictEqual(tsBooleanValue(tsInternalStreamIsReadableNode(loop, stream, false)), readableValue) && tsStrictEqual(tsBooleanValue(tsInternalStreamIsWritableNode(loop, stream)), writableValue)
    writableFinished := tsTruthy(tsInternalStreamIsWritableFinished(loop, stream, tsBooleanValue(false)))
    readableFinished := tsTruthy(tsInternalStreamIsReadableFinished(loop, stream, tsBooleanValue(false)))
    wState := tsGet(loop, stream, tsStringReference(tsStringUTF8("_writableState")))
    rState := tsGet(loop, stream, tsStringReference(tsStringUTF8("_readableState")))

    var immediateResult tsValue
    immediateSet := false
    if tsTruthy(tsInternalStreamIsClosed(loop, stream)) {
        immediateResult = tsInternalEosOnCloseError(loop, stream, readable, readableFinished, writable, writableFinished)
        immediateSet = true
    } else if tsTruthy(tsGet(loop, wState, tsStringReference(tsStringUTF8("errorEmitted")))) || tsTruthy(tsGet(loop, rState, tsStringReference(tsStringUTF8("errorEmitted")))) {
        if !willEmitClose { immediateResult = tsInternalEosGetErrored(loop, stream); immediateSet = true }
    } else if !readable && (!willEmitClose || tsTruthy(tsInternalStreamIsReadable(loop, stream))) && (writableFinished || tsStrictEqual(tsInternalStreamIsWritable(loop, stream), tsBooleanValue(false))) && (tsNullish(wState) || tsIsUndefined(tsGet(loop, wState, tsStringReference(tsStringUTF8("pendingcb")))) || tsNumber(tsGet(loop, wState, tsStringReference(tsStringUTF8("pendingcb")))) == 0) {
        immediateResult = tsInternalEosGetErrored(loop, stream); immediateSet = true
    } else if !writable && (!willEmitClose || tsTruthy(tsInternalStreamIsWritable(loop, stream))) && (readableFinished || tsStrictEqual(tsInternalStreamIsReadable(loop, stream), tsBooleanValue(false))) {
        immediateResult = tsInternalEosGetErrored(loop, stream); immediateSet = true
    } else if !tsNullish(rState) && tsTruthy(tsGet(loop, stream, tsStringReference(tsStringUTF8("req")))) && tsTruthy(tsGet(loop, stream, tsStringReference(tsStringUTF8("aborted")))) {
        immediateResult = tsInternalEosGetErrored(loop, stream); immediateSet = true
    }

    nop := tsFunctionValue(tsFunc(func(_ *tsLoop, _ ...tsValue) tsValue { return tsU }))
    callbackValue := callback
    cleanup := tsNamedNativeMethod("cleanup", 0, func(_ *tsLoop, _ tsValue, _ ...tsValue) tsValue { callbackValue = nop; return tsU })
    if immediateSet {
        errorOption := tsGet(loop, options, tsStringReference(tsStringUTF8("error")))
        if !tsStrictEqual(errorOption, tsBooleanValue(false)) {
            tsInternalEosCall(loop, stream, "on", tsStringReference(tsStringUTF8("error")), nop)
            cleanup = tsNamedNativeMethod("cleanup", 0, func(loop *tsLoop, _ tsValue, _ ...tsValue) tsValue {
                callbackValue = nop
                tsInternalEosCall(loop, stream, "removeListener", tsStringReference(tsStringUTF8("error")), nop)
                return tsU
            })
        }
    } else if !tsNullish(signal) && tsTruthy(tsGet(loop, signal, tsStringReference(tsStringUTF8("aborted")))) {
        immediateResult = tsInternalEosAbortError(loop, signal)
        immediateSet = true
    }

    invokeImmediate := tsFunctionValue(tsFunc(func(loop *tsLoop, _ ...tsValue) tsValue {
        if tsIsNull(immediateResult) { tsCallReceiver(loop, callbackValue, stream) } else { tsCallReceiver(loop, callbackValue, stream, immediateResult) }
        return tsU
    }))
    if immediateSet && tsTruthy(tsGet(loop, options, tsInternalEosState.syncCallback)) {
        tsCallReceiver(loop, invokeImmediate, tsU)
        return cleanup
    }
    if tsInternalEosHasAsyncContext(loop) { callbackValue = tsInternalEosBindAsyncResource(loop, callbackValue) }
    if immediateSet { tsInternalEosNextTick(loop, invokeImmediate); return cleanup }
    callbackValue = tsInternalEosOnce(callbackValue)

    var onFinish tsValue
    onLegacyFinish := tsFunctionValue(tsFunc(func(loop *tsLoop, _ ...tsValue) tsValue {
        if !tsTruthy(tsGet(loop, stream, tsStringReference(tsStringUTF8("writable")))) { tsCallReceiver(loop, onFinish, stream) }
        return tsU
    }))
    onFinish = tsFunctionValue(tsFunc(func(loop *tsLoop, _ ...tsValue) tsValue {
        writableFinished = true
        if tsTruthy(tsGet(loop, stream, tsStringReference(tsStringUTF8("destroyed")))) { willEmitClose = false }
        if willEmitClose && (!tsTruthy(tsGet(loop, stream, tsStringReference(tsStringUTF8("readable")))) || readable) { return tsU }
        if !readable || readableFinished { tsCallReceiver(loop, callbackValue, stream) }
        return tsU
    }))
    onEnd := tsFunctionValue(tsFunc(func(loop *tsLoop, _ ...tsValue) tsValue {
        readableFinished = true
        if tsTruthy(tsGet(loop, stream, tsStringReference(tsStringUTF8("destroyed")))) { willEmitClose = false }
        if willEmitClose && (!tsTruthy(tsGet(loop, stream, tsStringReference(tsStringUTF8("writable")))) || writable) { return tsU }
        if !writable || writableFinished { tsCallReceiver(loop, callbackValue, stream) }
        return tsU
    }))
    onError := tsFunctionValue(tsFunc(func(loop *tsLoop, _ ...tsValue) tsValue { tsCallReceiver(loop, callbackValue, stream, tsArg(args, 0)); return tsU }))
    onClose := tsFunctionValue(tsFunc(func(loop *tsLoop, _ ...tsValue) tsValue {
        errorValue := tsInternalEosOnCloseError(loop, stream, readable, readableFinished, writable, writableFinished)
        if tsIsNull(errorValue) { tsCallReceiver(loop, callbackValue, stream) } else { tsCallReceiver(loop, callbackValue, stream, errorValue) }
        return tsU
    }))
    onRequest := tsFunctionValue(tsFunc(func(loop *tsLoop, _ ...tsValue) tsValue {
        request := tsGet(loop, stream, tsStringReference(tsStringUTF8("req")))
        method := tsGet(loop, request, tsStringReference(tsStringUTF8("on")))
        if method.kind == tsFunctionKind { tsCallReceiver(loop, method, request, tsStringReference(tsStringUTF8("finish")), onFinish) }
        return tsU
    }))
    isRequest := tsTruthy(tsGet(loop, stream, tsStringReference(tsStringUTF8("setHeader")))) && tsGet(loop, stream, tsStringReference(tsStringUTF8("abort"))).kind == tsFunctionKind
    if isRequest {
        tsInternalEosCall(loop, stream, "on", tsStringReference(tsStringUTF8("complete")), onFinish)
        if !willEmitClose { tsInternalEosCall(loop, stream, "on", tsStringReference(tsStringUTF8("abort")), onClose) }
        if tsTruthy(tsGet(loop, stream, tsStringReference(tsStringUTF8("req")))) { tsCallReceiver(loop, onRequest, tsU) } else { tsInternalEosCall(loop, stream, "on", tsStringReference(tsStringUTF8("request")), onRequest) }
    } else if writable && tsNullish(wState) {
        tsInternalEosCall(loop, stream, "on", tsStringReference(tsStringUTF8("end")), onLegacyFinish)
        tsInternalEosCall(loop, stream, "on", tsStringReference(tsStringUTF8("close")), onLegacyFinish)
    }
    aborted := tsGet(loop, stream, tsStringReference(tsStringUTF8("aborted")))
    if !willEmitClose && aborted.kind == tsBooleanKind { tsInternalEosCall(loop, stream, "on", tsStringReference(tsStringUTF8("aborted")), onClose) }
    tsInternalEosCall(loop, stream, "on", tsStringReference(tsStringUTF8("end")), onEnd)
    tsInternalEosCall(loop, stream, "on", tsStringReference(tsStringUTF8("finish")), onFinish)
    if !tsStrictEqual(tsGet(loop, options, tsStringReference(tsStringUTF8("error"))), tsBooleanValue(false)) { tsInternalEosCall(loop, stream, "on", tsStringReference(tsStringUTF8("error")), onError) }
    tsInternalEosCall(loop, stream, "on", tsStringReference(tsStringUTF8("close")), onClose)

    cleanup = tsNamedNativeMethod("cleanup", 0, func(loop *tsLoop, _ tsValue, _ ...tsValue) tsValue {
        callbackValue = nop
        for _, item := range []struct{ event string; listener tsValue }{{"aborted", onClose}, {"complete", onFinish}, {"abort", onClose}, {"request", onRequest}, {"end", onLegacyFinish}, {"close", onLegacyFinish}, {"finish", onFinish}, {"end", onEnd}, {"error", onError}, {"close", onClose}} {
            tsInternalEosCall(loop, stream, "removeListener", tsStringReference(tsStringUTF8(item.event)), item.listener)
        }
        request := tsGet(loop, stream, tsStringReference(tsStringUTF8("req")))
        if !tsNullish(request) {
            method := tsGet(loop, request, tsStringReference(tsStringUTF8("removeListener")))
            if method.kind == tsFunctionKind { tsCallReceiver(loop, method, request, tsStringReference(tsStringUTF8("finish")), onFinish) }
        }
        return tsU
    })
    if !tsNullish(signal) {
        abort := tsFunctionValue(tsFunc(func(loop *tsLoop, _ ...tsValue) tsValue {
            endCallback := callbackValue
            tsCallReceiver(loop, cleanup, tsU)
            tsCallReceiver(loop, endCallback, stream, tsInternalEosAbortError(loop, signal))
            return tsU
        }))
        disposable := tsInternalEosAttachAbortListener(loop, signal, abort)
        original := callbackValue
        callbackValue = tsInternalEosOnce(tsNativeMethod(func(loop *tsLoop, receiver tsValue, args ...tsValue) tsValue {
            tsInternalEosDispose(loop, disposable)
            return tsCallReceiver(loop, original, receiver, args...)
        }))
    }
    return cleanup
}

func tsInternalEosFinished(loop *tsLoop, args []tsValue) tsValue {
    stream := tsArg(args, 0)
    options := tsArg(args, 1)
    autoCleanup := false
    if tsNullish(options) { options = tsInternalSharedSingleton.emptyObject }
    cleanupOption := tsGet(loop, options, tsStringReference(tsStringUTF8("cleanup")))
    if tsTruthy(cleanupOption) {
        if cleanupOption.kind != tsBooleanKind { tsInternalValidatorInvalidType("cleanup", "of type boolean", cleanupOption) }
        autoCleanup = tsTruthy(cleanupOption)
    }
    promise := loop.promise()
    resolve := tsNamedNativeMethod("resolve", 1, func(_ *tsLoop, _ tsValue, values ...tsValue) tsValue { promise.settle(tsResult{tsArg(values, 0), false}); return tsU })
    reject := tsNamedNativeMethod("reject", 1, func(_ *tsLoop, _ tsValue, values ...tsValue) tsValue { promise.settle(tsResult{tsArg(values, 0), true}); return tsU })
    var cleanup tsValue
    callback := tsFunctionValue(tsFunc(func(loop *tsLoop, values ...tsValue) tsValue {
        if autoCleanup && cleanup.kind == tsFunctionKind { tsCallReceiver(loop, cleanup, tsU) }
        err := tsArg(values, 0)
        if tsTruthy(err) { tsCallReceiver(loop, reject, tsU, err) } else { tsCallReceiver(loop, resolve, tsU) }
        return tsU
    }))
    cleanup = tsInternalEos(loop, []tsValue{stream, options, callback})
    return tsPromiseValue(promise)
}

func tsInternalStreamEndOfStreamModule(loop *tsLoop) tsValue {
    tsInternalEosState.Do(func() {
        tsInternalSharedModule(loop)
        tsInternalStreamUtilsInit(loop)
        tsInternalEosState.syncCallback = tsNewSymbol(tsStringUTF8("kEosNodeSynchronousCallback"))
        eos := tsFunctionValue(tsFunc(func(loop *tsLoop, args ...tsValue) tsValue { return tsInternalEos(loop, args) }))
        function := (*tsFunction)(eos.ref)
        function.properties = tsNewObject()
        function.properties.set("finished", tsNamedNativeMethod("finished", 2, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue { return tsInternalEosFinished(loop, args) }))
        function.properties.set("kEosNodeSynchronousCallback", tsInternalEosState.syncCallback)
        tsInternalEosState.module = eos
    })
    return tsInternalEosState.module
}
`
