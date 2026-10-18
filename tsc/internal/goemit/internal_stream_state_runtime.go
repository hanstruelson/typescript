package goemit

// internalStreamStateRuntime translates Node's high-water-mark defaults and
// per-stream option validation.
const internalStreamStateRuntime = `
var tsInternalStreamStateSingleton struct {
    sync.Mutex
    bytes int
    objects int
    initialized bool
    module *tsObject
}

func tsInternalStreamStateInit() {
    tsInternalStreamStateSingleton.Lock()
    if !tsInternalStreamStateSingleton.initialized {
        tsInternalStreamStateSingleton.bytes = 64 * 1024
        if runtime.GOOS == "windows" { tsInternalStreamStateSingleton.bytes = 16 * 1024 }
        tsInternalStreamStateSingleton.objects = 16
        tsInternalStreamStateSingleton.initialized = true
    }
    tsInternalStreamStateSingleton.Unlock()
}

func tsInternalStreamStateHighWaterMarkFrom(loop *tsLoop, options, duplexKey tsValue, isDuplex bool) tsValue {
    highWaterMark := tsGet(loop, options, tsStringReference(tsStringUTF8("highWaterMark")))
    if !tsNullish(highWaterMark) { return highWaterMark }
    if isDuplex { return tsGet(loop, options, duplexKey) }
    return tsNull
}

func tsInternalStreamStateDefault(objectMode bool) int {
    tsInternalStreamStateInit()
    tsInternalStreamStateSingleton.Lock()
    defer tsInternalStreamStateSingleton.Unlock()
    if objectMode { return tsInternalStreamStateSingleton.objects }
    return tsInternalStreamStateSingleton.bytes
}

func tsInternalStreamStateSetDefault(loop *tsLoop, objectMode, value tsValue) tsValue {
    if !tsIsNumeric(value) || math.IsNaN(tsNumber(value)) || math.IsInf(tsNumber(value), 0) || math.Trunc(tsNumber(value)) != tsNumber(value) || tsNumber(value) < 0 {
        tsInternalValidatorInvalidType("value", "a non-negative integer", value)
    }
    tsInternalStreamStateInit()
    tsInternalStreamStateSingleton.Lock()
    if tsTruthy(objectMode) { tsInternalStreamStateSingleton.objects = int(tsNumber(value)) } else { tsInternalStreamStateSingleton.bytes = int(tsNumber(value)) }
    tsInternalStreamStateSingleton.Unlock()
    return tsU
}

func tsInternalStreamStateGetHighWaterMark(loop *tsLoop, state, options, duplexKey tsValue, duplex bool) tsValue {
    highWaterMark := tsInternalStreamStateHighWaterMarkFrom(loop, options, duplexKey, duplex)
    if !tsNullish(highWaterMark) {
        number := tsNumber(highWaterMark)
        if !tsIsNumeric(highWaterMark) || math.IsNaN(number) || math.IsInf(number, 0) || math.Trunc(number) != number || number < 0 {
            name := "options.highWaterMark"
            if duplex { name = "options." + tsText(duplexKey) }
            panic(tsThrown{tsFSException("TypeError", "ERR_INVALID_ARG_VALUE", "The \""+name+"\" argument is invalid. Received "+tsText(highWaterMark))})
        }
        return tsNumberValue(math.Floor(number))
    }
    return tsNumberValue(float64(tsInternalStreamStateDefault(tsTruthy(tsGet(loop, state, tsStringReference(tsStringUTF8("objectMode")))))))
}

func tsInternalStreamStateModule(loop *tsLoop) tsValue {
    tsInternalStreamStateInit()
    tsInternalStreamStateSingleton.Lock()
    if tsInternalStreamStateSingleton.module != nil { module := tsInternalStreamStateSingleton.module; tsInternalStreamStateSingleton.Unlock(); return tsObjectValue(module) }
    module := tsNewObject()
    tsInternalStreamStateSingleton.module = module
    tsInternalStreamStateSingleton.Unlock()
    module.set("getHighWaterMark", tsNamedNativeMethod("getHighWaterMark", 4, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
        return tsInternalStreamStateGetHighWaterMark(loop, tsArg(args, 0), tsArg(args, 1), tsArg(args, 2), tsTruthy(tsArg(args, 3)))
    }))
    module.set("getDefaultHighWaterMark", tsNamedNativeMethod("getDefaultHighWaterMark", 1, func(_ *tsLoop, _ tsValue, args ...tsValue) tsValue {
        return tsNumberValue(float64(tsInternalStreamStateDefault(tsTruthy(tsArg(args, 0)))))
    }))
    module.set("setDefaultHighWaterMark", tsNamedNativeMethod("setDefaultHighWaterMark", 2, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
        return tsInternalStreamStateSetDefault(loop, tsArg(args, 0), tsArg(args, 1))
    }))
    module.set("default", tsObjectValue(module))
    return tsObjectValue(module)
}
`
