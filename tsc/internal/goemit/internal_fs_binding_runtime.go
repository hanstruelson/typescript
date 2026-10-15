package goemit

// internalFSBindingRuntime adapts the runtime's fs operations to Bun's shared
// internal/fs/binding surface. The same object is returned on each require so
// instrumentation wrappers observe one shared set of methods.
const internalFSBindingRuntime = `
var tsInternalFSBindingSingleton struct { sync.Once; value tsValue }

func tsInternalFSBindingModule(loop *tsLoop) tsValue {
	tsInternalFSBindingSingleton.Do(func() {
		binding := tsNewObject()
		fs := (*tsObject)(tsFSModule("fs").ref)
		for _, name := range fs.order { binding.set(name, fs.values[name]) }
		promises := (*tsObject)(tsFSModule("fs/promises").ref)
		for _, name := range promises.order { binding.set(name, promises.values[name]) }
		for _, operation := range tsFSOperations {
			op := operation
			if op != "mkdtempDisposable" {
				binding.set(op, tsNamedNativeMethod(op, 1, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
					return tsFSInvoke(loop, op, "promise", args)
				}))
				binding.set(op+"Cb", tsNamedNativeMethod(op+"Cb", 1, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
					return tsInternalFSBindingCallback(loop, op, args)
				}))
			}
		}
		binding.set("realpathNative", binding.get("realpath"))
		if syncRealpath := binding.get("realpathSync"); syncRealpath.kind == tsFunctionKind {
			binding.set("realpathNativeSync", tsGet(loop, syncRealpath, tsStringReference(tsStringUTF8("native"))))
		}
		tsInternalFSBindingSingleton.value = tsObjectValue(binding)
	})
	return tsInternalFSBindingSingleton.value
}

func tsInternalFSBindingCallback(loop *tsLoop, operation string, args []tsValue) tsValue {
	if len(args) == 0 || args[0].kind != tsFunctionKind { tsFSInvalid("callback", "a function") }
	callback := args[0]
	callArgs := append([]tsValue(nil), args[1:]...)
	callArgs = append(callArgs, callback)
	return tsFSInvoke(loop, operation, "callback", callArgs)
}
`
