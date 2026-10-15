package goemit

// internalReplNodeInspectRuntime exposes readline's width helpers and defers
// the large inspect module until its lazy property is actually read.
const internalReplNodeInspectRuntime = `
var tsInternalReplNodeInspect struct { sync.Once; value tsValue }

func tsInternalReplNodeInspectModule(loop *tsLoop) tsValue {
	tsInternalReplNodeInspect.Do(func(){
		module:=tsNewObject()
		module.set("getStringWidth",tsNamedNativeMethod("getStringWidth",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeInspectStringWidth(loop,tsU,args...)}))
		module.set("stripVTControlCharacters",tsNamedNativeMethod("stripVTControlCharacters",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeInspectStripVT(loop,tsU,args...)}))
		module.order=append(module.order,"inspect")
		module.descriptors=map[string]*tsDescriptor{"inspect":{accessor:true,enumerable:true,configurable:true,getter:tsNamedNativeMethod("get inspect",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsGet(loop,tsInternalInspectModule(loop),tsStringReference(tsStringUTF8("inspect")))})}}
		module.set("default",tsObjectValue(module));tsInternalReplNodeInspect.value=tsObjectValue(module)
	})
	return tsInternalReplNodeInspect.value
}
`
