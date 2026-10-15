package goemit

// internalAsyncHooksRuntime translates Bun's private async_hooks state used by
// internals that need to know whether async hooks are enabled.
const internalAsyncHooksRuntime = `
var tsInternalAsyncHooksState struct {
    sync.Mutex
    module *tsObject
    active int
    ownerSymbol tsValue
}

func tsInternalAsyncHooksMark(enabled bool) {
    tsInternalAsyncHooksState.Lock()
    if enabled {
        tsInternalAsyncHooksState.active++
    } else if tsInternalAsyncHooksState.active > 0 {
        tsInternalAsyncHooksState.active--
    }
    tsInternalAsyncHooksState.Unlock()
}

func tsInternalAsyncHooksModule() tsValue {
    tsInternalAsyncHooksState.Lock()
    defer tsInternalAsyncHooksState.Unlock()
    if tsInternalAsyncHooksState.module != nil {
        return tsObjectValue(tsInternalAsyncHooksState.module)
    }
    tsInternalAsyncHooksState.ownerSymbol = tsNewSymbol(tsStringUTF8("owner_symbol"))
    symbols := tsNewObject()
    symbols.set("owner_symbol", tsInternalAsyncHooksState.ownerSymbol)
    module := tsNewObject()
    module.set("enabledHooksExist", tsNamedNativeMethod("enabledHooksExist", 0, func(_ *tsLoop, _ tsValue, _ ...tsValue) tsValue {
        tsInternalAsyncHooksState.Lock()
        active := tsInternalAsyncHooksState.active > 0
        tsInternalAsyncHooksState.Unlock()
        return tsBooleanValue(active)
    }))
    module.set("markHookEnabled", tsNamedNativeMethod("markHookEnabled", 0, func(_ *tsLoop, _ tsValue, _ ...tsValue) tsValue {
        tsInternalAsyncHooksMark(true)
        return tsU
    }))
    module.set("markHookDisabled", tsNamedNativeMethod("markHookDisabled", 0, func(_ *tsLoop, _ tsValue, _ ...tsValue) tsValue {
        tsInternalAsyncHooksMark(false)
        return tsU
    }))
    module.set("symbols", tsObjectValue(symbols))
    tsInternalAsyncHooksState.module = module
    return tsObjectValue(module)
}
`
