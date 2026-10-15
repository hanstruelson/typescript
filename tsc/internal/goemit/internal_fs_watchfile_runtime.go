package goemit

// internalFSWatchFileRuntime adapts Node's lazy watchFile module to the shared
// per-loop stat watcher table in fs_watch.go.
const internalFSWatchFileRuntime = `
var tsInternalFSWatchFileModuleSingleton struct { sync.Once; value tsValue }

func tsInternalFSWatchFileModuleValue(loop *tsLoop) tsValue {
	tsInternalFSWatchFileModuleSingleton.Do(func(){
		module:=tsNewObject()
		module.set("watchFile",tsNamedNativeMethod("watchFile",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsFSWatchFile(loop,args)}))
		module.set("unwatchFile",tsNamedNativeMethod("unwatchFile",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsFSUnwatchFile(loop,args)}))
		module.set("default",tsObjectValue(module));tsInternalFSWatchFileModuleSingleton.value=tsObjectValue(module)
	})
	return tsInternalFSWatchFileModuleSingleton.value
}
`
