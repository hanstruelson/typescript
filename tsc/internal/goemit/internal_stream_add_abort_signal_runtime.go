package goemit

// internalStreamAddAbortSignalRuntime adapts Node's classic and web stream
// cancellation listeners to the local destroy and web-controller paths.
const internalStreamAddAbortSignalRuntime = `
func tsInternalStreamValidateAbortSignal(signal tsValue, name string) {
    if (signal.kind != tsObjectKind && signal.kind != tsInstanceKind) || !tsHas(signal, tsStringReference(tsStringUTF8("aborted"))) {
        tsInternalValidatorInvalidType(name, "an AbortSignal", signal)
    }
}

func tsInternalStreamAddAbortSignalNoValidate(loop *tsLoop, signal, stream tsValue) tsValue {
    if (signal.kind != tsObjectKind && signal.kind != tsInstanceKind) || !tsHas(signal, tsStringReference(tsStringUTF8("aborted"))) {
        return stream
    }
    onAbort := tsFunctionValue(tsFunc(func(loop *tsLoop, _ ...tsValue) tsValue {
        errorValue := tsInternalEosAbortError(loop, signal)
        if tsInternalStreamIsNode(loop, stream) {
            destroy := tsGet(loop, stream, tsStringReference(tsStringUTF8("destroy")))
            if destroy.kind == tsFunctionKind { tsCallReceiver(loop, destroy, stream, errorValue) }
            return tsU
        }
        controllerErrorKey := tsInternalStreamUtilsKey(loop, "nodejs.webstream.controllerErrorFunction")
        controllerError := tsGet(loop, stream, controllerErrorKey)
        if controllerError.kind == tsFunctionKind {
            tsCallReceiver(loop, controllerError, stream, errorValue)
        } else {
            tsWebStreamControllerError(loop, stream, errorValue)
        }
        return tsU
    }))
    if tsTruthy(tsGet(loop, signal, tsStringReference(tsStringUTF8("aborted")))) {
        tsCallReceiver(loop, onAbort, tsU)
    } else {
        disposable := tsInternalAbortListenerAdd(loop, signal, onAbort)
        dispose := tsGet(loop, disposable, tsWellKnownSymbol("dispose"))
        eos := tsInternalStreamEndOfStreamModule(loop)
        tsCallReceiver(loop, eos, tsU, stream, dispose)
    }
    return stream
}

func tsInternalStreamAddAbortSignal(loop *tsLoop, signal, stream tsValue) tsValue {
    tsInternalStreamValidateAbortSignal(signal, "signal")
    if !tsInternalStreamIsNode(loop, stream) && !tsInternalStreamIsWeb(loop, stream) {
        tsInternalValidatorInvalidType("stream", "an instance of ReadableStream, WritableStream, or Stream", stream)
    }
    return tsInternalStreamAddAbortSignalNoValidate(loop, signal, stream)
}

func tsInternalStreamAddAbortSignalModule(loop *tsLoop) tsValue {
    module := tsNewObject()
    add := tsNamedNativeMethod("addAbortSignal", 2, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
        return tsInternalStreamAddAbortSignal(loop, tsArg(args, 0), tsArg(args, 1))
    })
    noValidate := tsNamedNativeMethod("addAbortSignalNoValidate", 2, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
        return tsInternalStreamAddAbortSignalNoValidate(loop, tsArg(args, 0), tsArg(args, 1))
    })
    module.set("addAbortSignal", add)
    module.set("addAbortSignalNoValidate", noValidate)
    module.set("default", tsObjectValue(module))
    return tsObjectValue(module)
}
`
