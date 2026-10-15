package goemit

// internalAbortListenerRuntime maps Node's AbortSignal listener helper onto
// the runtime's EventTarget-compatible AbortSignal objects.
const internalAbortListenerRuntime = `
func tsInternalAbortListenerAdd(loop *tsLoop, signal, listener tsValue) tsValue {
    if tsIsUndefined(signal) {
        tsInternalValidatorInvalidType("signal", "an AbortSignal", signal)
    }
    tsInternalValidatorAbortSignal(signal, tsStringReference(tsStringUTF8("signal")))
    if listener.kind != tsFunctionKind {
        tsInternalValidatorInvalidType("listener", "of type function", listener)
    }

    disposable := tsNewObject()
    disposable.prototype = tsNull
    var remove func()
    if tsTruthy(tsGet(loop, signal, tsStringReference(tsStringUTF8("aborted")))) {
        loop.ready = append(loop.ready, func() { tsCallReceiver(loop, listener, tsU) })
    } else {
        options := tsNewObject()
        options.prototype = tsNull
        options.set("once", tsBooleanValue(true))
        shared := tsInternalSharedModule(loop)
        resist := tsGet(loop, shared, tsStringReference(tsStringUTF8("resistStopPropagation")))
        optionsValue := tsCallReceiver(loop, resist, shared, tsObjectValue(options))
        add := tsGet(loop, signal, tsStringReference(tsStringUTF8("addEventListener")))
        tsCallReceiver(loop, add, signal, tsStringReference(tsStringUTF8("abort")), listener, optionsValue)
        remove = func() {
            removeMethod := tsGet(loop, signal, tsStringReference(tsStringUTF8("removeEventListener")))
            if removeMethod.kind == tsFunctionKind {
                tsCallReceiver(loop, removeMethod, signal, tsStringReference(tsStringUTF8("abort")), listener)
            }
        }
    }
    dispose := tsNamedNativeMethod("[Symbol.dispose]", 0, func(_ *tsLoop, _ tsValue, _ ...tsValue) tsValue {
        if remove != nil {
            remove()
            remove = nil
        }
        return tsU
    })
    disposable.set(tsPropertyKey(tsWellKnownSymbol("dispose")), dispose)
    return tsObjectValue(disposable)
}

func tsInternalAbortListenerModule(loop *tsLoop) tsValue {
    module := tsNewObject()
    add := tsNamedNativeMethod("addAbortListener", 2, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
        return tsInternalAbortListenerAdd(loop, tsArg(args, 0), tsArg(args, 1))
    })
    module.set("addAbortListener", add)
    module.set("default", tsObjectValue(module))
    return tsObjectValue(module)
}
`
