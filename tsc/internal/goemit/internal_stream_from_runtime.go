package goemit

// internalStreamFromRuntime adapts Node's iterable-to-Readable stream bridge.
const internalStreamFromRuntime = `
func tsInternalStreamFromNullValuesError() tsValue {
    return tsFSException("TypeError", "ERR_STREAM_NULL_VALUES", "May not write null values to stream")
}

func tsInternalStreamFromOptions(loop *tsLoop, options tsValue, highWaterMark int) tsValue {
    out := tsNewObject()
    out.set("objectMode", tsBooleanValue(true))
    if highWaterMark >= 0 { out.set("highWaterMark", tsNumberValue(float64(highWaterMark))) }
    if options.kind == tsObjectKind || options.kind == tsInstanceKind || options.kind == tsArrayKind {
        for _, key := range tsObjectKeys(options).values { tsSet(loop, out, key, tsGet(loop, options, key)) }
    }
    return tsObjectValue(out)
}

func tsInternalStreamFromConstruct(loop *tsLoop, readable tsValue, options tsValue) tsValue {
    return tsConstruct(readable, options)
}

func tsInternalStreamFrom(loop *tsLoop, args []tsValue) tsValue {
    readableClass, iterable, options := tsArg(args, 0), tsArg(args, 1), tsArg(args, 2)
    byteBuffer := iterable.kind == tsTypedArrayKind && (*tsTypedArray)(iterable.ref).nodeBuffer
    if iterable.kind == tsStringKind || byteBuffer {
        streamOptions := tsInternalStreamFromOptions(loop, options, -1)
        pushed := false
        read := tsNamedNativeMethod("read", 0, func(loop *tsLoop, receiver tsValue, _ ...tsValue) tsValue {
            if !pushed {
                pushed = true
                push := tsGet(loop, receiver, tsStringReference(tsStringUTF8("push")))
                if push.kind == tsFunctionKind { tsCallReceiver(loop, push, receiver, iterable); tsCallReceiver(loop, push, receiver, tsNull) }
            }
            return tsU
        })
        tsSet(loop, streamOptions, tsStringReference(tsStringUTF8("read")), read)
        return tsInternalStreamFromConstruct(loop, readableClass, streamOptions)
    }

    asyncMethod := tsGet(loop, iterable, tsWellKnownSymbol("asyncIterator"))
    syncMethod := tsGet(loop, iterable, tsWellKnownSymbol("iterator"))
    if tsNullish(asyncMethod) && tsNullish(syncMethod) {
        tsInternalValidatorInvalidType("iterable", "an instance of Iterable", iterable)
    }
    iterator := tsAsyncIterate(iterable)
    streamOptions := tsInternalStreamFromOptions(loop, options, 1)
    stream := tsInternalStreamFromConstruct(loop, readableClass, streamOptions)
    reading := false
    var next func(*tsLoop)
    next = func(loop *tsLoop) {
        if reading { return }
        reading = true
        promise := iterator.next(loop)
        promise.then(func(result tsResult) {
            if result.rejected {
                reading = false
                destroy := tsGet(loop, stream, tsStringReference(tsStringUTF8("destroy")))
                if destroy.kind == tsFunctionKind { tsCallReceiver(loop, destroy, stream, result.value) }
                return
            }
            if tsTruthy(tsGet(loop, result.value, tsStringReference(tsStringUTF8("done")))) {
                push := tsGet(loop, stream, tsStringReference(tsStringUTF8("push")))
                if push.kind == tsFunctionKind { tsCallReceiver(loop, push, stream, tsNull) }
                return
            }
            value := tsGet(loop, result.value, tsStringReference(tsStringUTF8("value")))
            if tsIsNull(value) {
                reading = false
                destroy := tsGet(loop, stream, tsStringReference(tsStringUTF8("destroy")))
                if destroy.kind == tsFunctionKind { tsCallReceiver(loop, destroy, stream, tsInternalStreamFromNullValuesError()) }
                return
            }
            push := tsGet(loop, stream, tsStringReference(tsStringUTF8("push")))
            accepted := tsTruthy(tsCallReceiver(loop, push, stream, value))
            if accepted { next(loop) } else { reading = false }
        })
    }
    read := tsNamedNativeMethod("read", 0, func(loop *tsLoop, _ tsValue, _ ...tsValue) tsValue { next(loop); return tsU })
    tsSet(loop, stream, tsStringReference(tsStringUTF8("_read")), read)
    originalDestroy := tsGet(loop, stream, tsStringReference(tsStringUTF8("_destroy")))
    replacementDestroy := tsNativeMethod(func(loop *tsLoop, receiver tsValue, destroyArgs ...tsValue) tsValue {
        err, callback := tsArg(destroyArgs, 0), tsArg(destroyArgs, 1)
        finish := func(destroyErr tsValue) {
            combined := destroyErr
            if !tsTruthy(combined) { combined = err }
            closePromise := iterator.close(loop)
            closePromise.then(func(result tsResult) {
                finalErr := combined
                if result.rejected {
                    errors := tsInternalErrorsModule(loop)
                    aggregate := tsGet(loop, errors, tsStringReference(tsStringUTF8("aggregateTwoErrors")))
                    finalErr = tsCallReceiver(loop, aggregate, errors, combined, result.value)
                }
                if callback.kind == tsFunctionKind {
                    callbackValue := callback
                    finalErrValue := finalErr
                    tick := tsFunctionValue(tsFunc(func(loop *tsLoop, _ ...tsValue) tsValue { tsCall(loop, callbackValue, finalErrValue); return tsU }))
                    tsInternalStreamDestroyNextTick(loop, tick)
                }
            })
        }
        if originalDestroy.kind == tsFunctionKind {
            done := tsFunctionValue(tsFunc(func(loop *tsLoop, values ...tsValue) tsValue { finish(tsArg(values, 0)); return tsU }))
            tsCallReceiver(loop, originalDestroy, receiver, err, done)
        } else {
            finish(tsU)
        }
        return tsU
    })
    tsSet(loop, stream, tsStringReference(tsStringUTF8("_destroy")), replacementDestroy)
    return stream
}

func tsInternalStreamFromModule(loop *tsLoop) tsValue {
    return tsNamedNativeMethod("from", 3, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue { return tsInternalStreamFrom(loop, args) })
}
`
