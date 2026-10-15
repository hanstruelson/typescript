package goemit

// nodeBunJSCRuntime provides the heap counters consumed by Node's vm shim
// through Go's runtime instead of Bun's native JSC binding.
const nodeBunJSCRuntime = `
var tsNodeBunJSC struct { sync.Once; module tsValue; peak uint64 }
func tsNodeBunJSCModule(loop *tsLoop) tsValue {
	tsNodeBunJSC.Do(func(){
		module:=tsNewObject()
		module.set("memoryUsage",tsNamedNativeMethod("memoryUsage",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{
			var memory runtime.MemStats;runtime.ReadMemStats(&memory);if memory.Sys>tsNodeBunJSC.peak{tsNodeBunJSC.peak=memory.Sys}
			value:=tsNewObject();value.set("current",tsNumberValue(float64(memory.HeapAlloc)));value.set("peak",tsNumberValue(float64(tsNodeBunJSC.peak)));return tsObjectValue(value)
		}))
		module.set("default",tsObjectValue(module));tsNodeBunJSC.module=tsObjectValue(module)
	})
	return tsNodeBunJSC.module
}
`
