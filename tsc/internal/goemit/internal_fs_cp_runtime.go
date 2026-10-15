package goemit

// internalFSCopyRuntime routes Node's recursive-copy helper to the shared Go
// filesystem walker; native clone/copy fast-path probing is represented by a
// false result so the common walker handles every source type.
const internalFSCopyRuntime = `
var tsInternalFSCopySingleton struct { sync.Once; value tsValue }
func tsInternalFSCopyModule(loop *tsLoop) tsValue {
	tsInternalFSCopySingleton.Do(func(){
		module:=tsNewObject()
		module.set("tryNativeFastPath",tsNamedNativeMethod("tryNativeFastPath",3,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{
			result:=tsNewObject();result.set("ok",tsBooleanValue(false));result.set("checked",tsU);return tsPromiseValue(loop.resolved(tsObjectValue(result),false))
		}))
		module.set("cpFn",tsNamedNativeMethod("cpFn",3,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{
			return tsFSInvoke(loop,"cp","promise",[]tsValue{tsArg(args,0),tsArg(args,1),tsArg(args,2)})
		}))
		module.set("default",tsObjectValue(module));tsInternalFSCopySingleton.value=tsObjectValue(module)
	})
	return tsInternalFSCopySingleton.value
}
`
