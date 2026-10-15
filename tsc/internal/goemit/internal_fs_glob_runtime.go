package goemit

// internalFSGlobRuntime adapts Bun's internal/fs/glob exports to the runtime's
// filesystem glob walker. Bun's minimatch parser is represented by the local
// glob matcher; the public entry points and async iterator shape are retained.
const internalFSGlobRuntime = `
var tsInternalFSGlobSingleton struct { sync.Once; value tsValue; class *tsClass }

func tsInternalFSGlobModule(loop *tsLoop) tsValue {
	tsInternalFSGlobSingleton.Do(func() {
		module := tsNewObject()
		glob := tsNamedNativeMethod("glob", 1, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
			return tsFSGlobIterator(loop, args)
		})
		globSync := tsNamedNativeMethod("globSync", 1, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
			result := tsFSGlob(loop, args, false)
			if result.rejected { panic(tsThrown{result.value}) }
			return result.value
		})
		class := &tsClass{builtin: "Glob", nativePrototype: tsNewObject()}
		class.nativePrototype.set("constructor", tsClassValue(class))
		class.construct = func(args ...tsValue) tsValue {
			pattern, options := tsArg(args, 0), tsArg(args, 1)
			instance := tsNewObject()
			instance.nativeClass = class
			instance.prototype = tsObjectValue(class.nativePrototype)
			instance.set("matchers", tsArrayValue(&tsArray{}))
			instance.set("glob", tsNamedNativeMethod("glob", 0, func(loop *tsLoop, _ tsValue, _ ...tsValue) tsValue {
				return tsFSGlobIterator(loop, []tsValue{pattern, options})
			}))
			instance.set("globSync", tsNamedNativeMethod("globSync", 0, func(loop *tsLoop, _ tsValue, _ ...tsValue) tsValue {
				result := tsFSGlob(loop, []tsValue{pattern, options}, false)
				if result.rejected { panic(tsThrown{result.value}) }
				return result.value
			}))
			return tsObjectValue(instance)
		}
		module.set("glob", glob)
		module.set("globSync", globSync)
		module.set("Glob", tsClassValue(class))
		module.set("default", tsObjectValue(module))
		tsInternalFSGlobSingleton.value = tsObjectValue(module)
		tsInternalFSGlobSingleton.class = class
	})
	return tsInternalFSGlobSingleton.value
}
`
