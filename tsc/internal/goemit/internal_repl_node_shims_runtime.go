package goemit

// internalReplNodeShimsRuntime collects the small internal modules consumed
// by Bun's Node REPL sources and routes them through the existing Go adapters.
const internalReplNodeShimsRuntime = `
var tsInternalReplNodeShims struct { sync.Once; value tsValue; builtinLibs tsValue; captureCallbacks []tsValue; captureInstalled bool }

func tsInternalReplNodeShimsModule(loop *tsLoop) tsValue {
	tsInternalReplNodeShims.Do(func(){
		module:=tsNewObject();util:=tsNodeUtilModule(loop);moduleModule:=tsNodeModuleModule(loop);colors:=tsInternalUtilColorsModule(loop);primordials:=tsInternalReplPrimordialsModule(loop)
		wrapIntrinsic:=func(name string,args ...tsValue)tsValue{fn:=tsGet(loop,primordials,tsStringReference(tsStringUTF8(name)));return tsCall(loop,fn,args...)}
		module.set("SideEffectFreeRegExpPrototypeSymbolReplace",tsNamedNativeMethod("SideEffectFreeRegExpPrototypeSymbolReplace",3,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return wrapIntrinsic("RegExpPrototypeSymbolReplace",tsArg(args,0),tsArg(args,1),tsArg(args,2))}))
		module.set("SideEffectFreeRegExpPrototypeSymbolSplit",tsNamedNativeMethod("SideEffectFreeRegExpPrototypeSymbolSplit",3,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return wrapIntrinsic("RegExpPrototypeSymbolSplit",tsArg(args,0),tsArg(args,1),tsArg(args,2))}))
		module.set("decorateErrorStack",tsNamedNativeMethod("decorateErrorStack",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsInternalReplNodeShimsDecorateStack(loop,tsArg(args,0))}))
		module.set("isError",tsNamedNativeMethod("isError",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsBooleanValue(tsInternalReplNodeShimsIsError(loop,tsArg(args,0)))}))
		module.set("deprecate",tsGet(loop,util,tsStringReference(tsStringUTF8("deprecate"))))
		module.set("shouldColorize",tsGet(loop,colors,tsStringReference(tsStringUTF8("shouldColorize"))))
		module.set("debuglog",tsNamedNativeMethod("debuglog",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{debuglog:=tsGet(loop,util,tsStringReference(tsStringUTF8("debuglog")));fn:=tsCallReceiver(loop,debuglog,util,tsArg(args,0));callback:=tsArg(args,1);if callback.kind==tsFunctionKind{tsCall(loop,callback,fn)};return fn}))
		module.set("sendInspectorCommand",tsNamedNativeMethod("sendInspectorCommand",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{onError:=tsArg(args,1);if onError.kind==tsFunctionKind{return tsCall(loop,onError)};return tsU}))
		types:=tsGet(loop,util,tsStringReference(tsStringUTF8("types")));module.set("isProxy",tsGet(loop,types,tsStringReference(tsStringUTF8("isProxy"))))
		module.set("getOptionValue",tsNamedNativeMethod("getOptionValue",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsInternalReplNodeShimsOption(loop,tsArg(args,0))}))
		module.set("isEnabled",tsNamedNativeMethod("isEnabled",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsBooleanValue(false)}));module.set("has",tsNamedNativeMethod("has",1,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsBooleanValue(true)}))
		builtinModule:=tsNewObject();builtinModule.set("getSchemeOnlyModuleNames",tsNamedNativeMethod("getSchemeOnlyModuleNames",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsArrayValue(&tsArray{values:[]tsValue{tsStringReference(tsStringUTF8("test"))}})}));module.set("BuiltinModule",tsObjectValue(builtinModule))
		formatMap:=tsNewObject();for _,entry:=range [][2]string{{".cjs","commonjs"},{".js","module"},{".json","json"},{".mjs","module"},{".node","addon"},{".wasm","wasm"}}{formatMap.set(entry[0],tsStringReference(tsStringUTF8(entry[1])))};module.set("extensionFormatMap",tsObjectValue(formatMap))
		loader:=tsNewObject();loader.set("kEvaluationPhase",tsStringReference(tsStringUTF8("evaluation")));loader.set("kSourcePhase",tsStringReference(tsStringUTF8("source")));loader.set("import",tsNamedNativeMethod("import",4,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsPromiseValue(loop.resolved(tsFSException("Error","ERR_METHOD_NOT_IMPLEMENTED","Dynamic import is not available in the Go runtime"),true))}));module.set("getOrInitializeCascadedLoader",tsNamedNativeMethod("getOrInitializeCascadedLoader",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsObjectValue(loader)}))
		module.set("getBuiltinLibs",tsNamedNativeMethod("getBuiltinLibs",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsInternalReplNodeShimsBuiltinLibs(loop)}));module.set("addBuiltinLibsToObject",tsNamedNativeMethod("addBuiltinLibsToObject",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsInternalReplNodeShimsAddBuiltins(loop,tsArg(args,0))}))
		module.set("makeRequireFunction",tsNamedNativeMethod("makeRequireFunction",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsInternalReplNodeShimsMakeRequire(loop,tsArg(args,0))}))
		module.set("makeContextifyScript",tsNamedNativeMethod("makeContextifyScript",9,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{vm:=tsNodeVMModule(loop);constructor:=tsGet(loop,vm,tsStringReference(tsStringUTF8("Script")));options:=tsNewObject();for index,name:=range []string{"filename","lineOffset","columnOffset","cachedData","produceCachedData","parsingContext","hostDefinedOptionId","importModuleDynamically"}{value:=tsArg(args,index+1);if !tsIsUndefined(value){options.set(name,value)}};return tsConstruct(constructor,tsArg(args,0),tsObjectValue(options))}))
		module.set("startSigintWatchdog",tsNamedNativeMethod("startSigintWatchdog",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsBooleanValue(true)}));module.set("stopSigintWatchdog",tsNamedNativeMethod("stopSigintWatchdog",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsBooleanValue(false)}))
		constants:=tsNewObject();constants.set("ALL_PROPERTIES",tsNumberValue(0));constants.set("SKIP_SYMBOLS",tsNumberValue(16));module.set("constants",tsObjectValue(constants));inspect:=tsInternalInspectModule(loop);module.set("getOwnNonIndexProperties",tsGet(loop,inspect,tsStringReference(tsStringUTF8("getOwnNonIndexProperties"))))
		module.set("addUncaughtExceptionCaptureCallback",tsNamedNativeMethod("addUncaughtExceptionCaptureCallback",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsInternalReplNodeShimsAddCapture(loop,tsArg(args,0))}))
		moduleConstructor:=tsInternalReplNodeShimsModuleShim(loop,moduleModule);module.set("Module",moduleConstructor)
		module.set("default",tsObjectValue(module));tsInternalReplNodeShims.value=tsObjectValue(module)
	})
	return tsInternalReplNodeShims.value
}
func tsInternalReplNodeShimsDecorateStack(loop *tsLoop, errorValue tsValue) tsValue {
	stack:=tsGet(loop,errorValue,tsStringReference(tsStringUTF8("stack")));if stack.kind!=tsStringKind{return errorValue};text:=tsText(stack);lines:=strings.Split(text,"\n");last:=len(lines)-1
	for index,line:=range lines{line=strings.ReplaceAll(line,"<anonymous> (","");if strings.HasPrefix(line,"  at REPL"){parts:=strings.Split(line,":");if len(parts)>=3{last=index}};lines[index]=line};if last>=0&&last<len(lines)-1{lines=lines[:last]};tsSet(loop,errorValue,tsStringReference(tsStringUTF8("stack")),tsStringReference(tsStringUTF8(strings.Join(lines,"\n"))));return errorValue
}
func tsInternalReplNodeShimsIsError(loop *tsLoop, value tsValue) bool {
	if value.kind==tsErrorKind{return true};if value.kind==tsObjectKind&&(*tsObject)(value.ref).nativeError{return true};if value.kind==tsInstanceKind&&(*tsProperties)(value.ref).nativeError{return true};util:=tsNodeUtilModule(loop);types:=tsGet(loop,util,tsStringReference(tsStringUTF8("types")));test:=tsGet(loop,types,tsStringReference(tsStringUTF8("isNativeError")));return test.kind==tsFunctionKind&&tsTruthy(tsCall(loop,test,value))
}
func tsInternalReplNodeShimsOption(loop *tsLoop, name tsValue) tsValue {
	value:=tsText(name);switch value{case "--pending-deprecation","--experimental-repl-await","--use-strict":
		process:=tsNodeProcessModule(loop);args:=tsGet(loop,process,tsStringReference(tsStringUTF8("execArgv")));found:=false;if args.kind==tsArrayKind{for _,item:=range (*tsArray)(args.ref).values{if tsText(item)==value{found=true}}};if value=="--experimental-repl-await"{return tsBooleanValue(true)};if value=="--use-strict"{return tsBooleanValue(false)};return tsBooleanValue(found)
	};return tsU
}
func tsInternalReplNodeShimsBuiltinLibs(loop *tsLoop) tsValue {
	if !tsIsUndefined(tsInternalReplNodeShims.builtinLibs){return tsInternalReplNodeShims.builtinLibs};module:=tsNodeModuleModule(loop);all:=tsGet(loop,module,tsStringReference(tsStringUTF8("builtinModules")));out:=&tsArray{};if all.kind==tsArrayKind{for _,item:=range (*tsArray)(all.ref).values{name:=tsText(item);if strings.HasPrefix(name,"_")||strings.HasPrefix(name,"node:")||strings.HasPrefix(name,"bun")||name=="undici"||name=="ws"{continue};out.values=append(out.values,item)}};tsInternalReplNodeShims.builtinLibs=tsArrayValue(out);return tsInternalReplNodeShims.builtinLibs
}
func tsInternalReplNodeShimsAddBuiltins(loop *tsLoop, object tsValue) tsValue {
	if object.kind!=tsObjectKind&&object.kind!=tsInstanceKind{return tsU};for _,entry:=range (*tsArray)(tsInternalReplNodeShimsBuiltinLibs(loop).ref).values{name:=tsText(entry);if strings.Contains(name,"/"){continue};key:=tsStringReference(tsStringUTF8(name));if tsObjectHasOwn(object,name){continue};moduleObject:=object
		getter:=tsNamedNativeMethod("get "+name,0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{lib:=tsNodeModuleRequire(loop,[]tsValue{entry});return lib});setter:=tsNamedNativeMethod("set "+name,1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{tsSet(loop,moduleObject,key,tsArg(args,0));return tsU})
		properties:=(*tsObject)(object.ref);if properties.descriptors==nil{properties.descriptors=map[string]*tsDescriptor{}};properties.order=append(properties.order,name);properties.descriptors[name]=&tsDescriptor{accessor:true,enumerable:false,configurable:true,getter:getter,setter:setter}
	};return object
}
func tsInternalReplNodeShimsMakeRequire(loop *tsLoop, moduleValue tsValue) tsValue {
	module:=tsNodeModuleModule(loop);createRequire:=tsGet(loop,module,tsStringReference(tsStringUTF8("createRequire")));process:=tsNodeProcessModule(loop);execPath:=tsGet(loop,process,tsStringReference(tsStringUTF8("execPath")));cwd:=tsCallReceiver(loop,tsGet(loop,process,tsStringReference(tsStringUTF8("cwd"))),process)
	if tsIsUndefined(cwd)||cwd.kind!=tsStringKind{path:=tsPathModule();cwd=tsCallReceiver(loop,tsGet(loop,path,tsStringReference(tsStringUTF8("dirname"))),path,execPath)}
	path:=tsPathModule();join:=tsGet(loop,path,tsStringReference(tsStringUTF8("join")));filename:=tsCallReceiver(loop,join,path,cwd,tsStringReference(tsStringUTF8("<repl>")));return tsCallReceiver(loop,createRequire,module,filename)
}
func tsInternalReplNodeShimsModuleShim(loop *tsLoop, moduleValue tsValue) tsValue {
	prototype:=tsNewObject();constructor:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{object:=tsNewObject();object.prototype=tsObjectValue(prototype);object.set("id",tsArg(args,0));object.set("path",tsStringReference(tsStringUTF8("")));object.set("exports",tsObjectValue(tsNewObject()));object.set("filename",tsNull);object.set("loaded",tsBooleanValue(false));object.set("children",tsArrayValue(&tsArray{}));object.set("paths",tsArrayValue(&tsArray{}));object.set("parent",tsArg(args,1));return tsObjectValue(object)});constructor.name="Module";constructor.length=2;constructor.constructible=true;constructor.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsCall(loop,tsFunctionValue(constructor),args...)};constructor.properties=tsNewObject();constructor.properties.set("prototype",tsObjectValue(prototype));prototype.set("constructor",tsFunctionValue(constructor))
	for _,name:=range []string{"builtinModules","globalPaths","_extensions"}{constructor.properties.set(name,tsGet(loop,moduleValue,tsStringReference(tsStringUTF8(name))))}
	for _,name:=range []string{"_nodeModulePaths","_resolveLookupPaths","_resolveFilename"}{method:=tsGet(loop,moduleValue,tsStringReference(tsStringUTF8(name)));if !tsIsUndefined(method){constructor.properties.set(name,method)}}
	return tsFunctionValue(constructor)
}
func tsInternalReplNodeShimsAddCapture(loop *tsLoop, callback tsValue) tsValue {
	if callback.kind!=tsFunctionKind{tsFSInvalid("cb","a function")};tsInternalReplNodeShims.captureCallbacks=append(tsInternalReplNodeShims.captureCallbacks,callback)
	if tsInternalReplNodeShims.captureInstalled{return tsU};process:=tsNodeProcessModule(loop);set:=tsGet(loop,process,tsStringReference(tsStringUTF8("setUncaughtExceptionCaptureCallback")));if set.kind!=tsFunctionKind{return tsU}
	dispatch:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{errorValue:=tsArg(args,0);for _,fn:=range tsInternalReplNodeShims.captureCallbacks{if tsTruthy(tsCall(loop,fn,errorValue)){return tsU}};emit:=tsGet(loop,process,tsStringReference(tsStringUTF8("emit")));if emit.kind==tsFunctionKind&&tsTruthy(tsCallReceiver(loop,emit,process,tsStringReference(tsStringUTF8("uncaughtException")),errorValue,tsStringReference(tsStringUTF8("uncaughtException")))){return tsU};return tsU}))
	func(){defer func(){if recover()!=nil{}}();tsCallReceiver(loop,set,process,dispatch);tsInternalReplNodeShims.captureInstalled=true}();return tsU
}
`
