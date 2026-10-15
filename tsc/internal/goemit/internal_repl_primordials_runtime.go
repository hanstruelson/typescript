package goemit

// internalReplPrimordialsRuntime supplies the captured intrinsic wrappers used
// by Bun's Node readline/repl JavaScript sources.
const internalReplPrimordialsRuntime = `
var tsInternalReplPrimordials struct { sync.Once; value tsValue }

func tsInternalReplPrimordialsCallable(loop *tsLoop,target,receiver tsValue,args []tsValue) tsValue {
	if target.kind!=tsFunctionKind{tsPropertyFailure("Primordial target is not callable")}
	return tsCallReceiver(loop,target,receiver,args...)
}
func tsInternalReplPrimordialsUncurried(loop *tsLoop,name string,target tsValue) tsValue {
	return tsNamedNativeMethod(name,1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{
		if len(args)==0{tsFSInvalid("receiver","an object")}
		return tsInternalReplPrimordialsCallable(loop,target,args[0],args[1:])
	})
}
func tsInternalReplPrimordialsGlobal(loop *tsLoop,name string) tsValue { return tsGlobalIdentifier(loop,name) }
func tsInternalReplPrimordialsPrototypeMethod(loop *tsLoop,globalName,method string) (tsValue,tsValue) {
	constructor:=tsInternalReplPrimordialsGlobal(loop,globalName);prototype:=tsGet(loop,constructor,tsStringReference(tsStringUTF8("prototype")))
	return tsGet(loop,prototype,tsStringReference(tsStringUTF8(method))),constructor
}
func tsInternalReplPrimordialsModule(loop *tsLoop) tsValue {
	tsInternalReplPrimordials.Do(func(){
		module:=tsNewObject();set:=func(name string,value tsValue){module.set(name,value)}
		base:=tsInternalPrimordialsModule(loop)
		for _,name:=range []string{"SafeSet","SafeStringIterator","SafeArrayIterator"}{set(name,tsFSProperty(base,name))}
		for _,name:=range []string{"Boolean","Error","Number","Promise","RegExp","Symbol","SyntaxError","globalThis"}{set(name,tsInternalReplPrimordialsGlobal(loop,name))}
		set("SymbolAsyncIterator",tsWellKnownSymbol("asyncIterator"));set("SymbolDispose",tsWellKnownSymbol("dispose"))
		static:=[]struct{name,owner,method string}{
			{"ArrayFrom","Array","from"},{"JSONStringify","JSON","stringify"},{"MathCeil","Math","ceil"},{"MathFloor","Math","floor"},{"MathMax","Math","max"},{"MathMin","Math","min"},
			{"NumberIsFinite","Number","isFinite"},{"NumberIsNaN","Number","isNaN"},{"NumberParseFloat","Number","parseFloat"},{"NumberParseInt","Number","parseInt"},
			{"ObjectAssign","Object","assign"},{"ObjectDefineProperties","Object","defineProperties"},{"ObjectDefineProperty","Object","defineProperty"},{"ObjectGetOwnPropertyDescriptor","Object","getOwnPropertyDescriptor"},{"ObjectGetOwnPropertyNames","Object","getOwnPropertyNames"},{"ObjectGetPrototypeOf","Object","getPrototypeOf"},{"ObjectKeys","Object","keys"},
			{"PromiseReject","Promise","reject"},{"StringFromCharCode","String","fromCharCode"},
		}
		for _,entry:=range static{item:=entry;owner:=tsInternalReplPrimordialsGlobal(loop,item.owner);method:=tsGet(loop,owner,tsStringReference(tsStringUTF8(item.method)));set(item.name,tsNamedNativeMethod(item.name,1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsInternalReplPrimordialsCallable(loop,method,owner,args)}))}
		for _,entry:=range []struct{name,owner,method string}{
			{"ArrayPrototypeAt","Array","at"},{"ArrayPrototypeFilter","Array","filter"},{"ArrayPrototypeForEach","Array","forEach"},{"ArrayPrototypeIncludes","Array","includes"},{"ArrayPrototypeIndexOf","Array","indexOf"},{"ArrayPrototypeJoin","Array","join"},{"ArrayPrototypeMap","Array","map"},{"ArrayPrototypePop","Array","pop"},{"ArrayPrototypePush","Array","push"},{"ArrayPrototypeReverse","Array","reverse"},{"ArrayPrototypeShift","Array","shift"},{"ArrayPrototypeSlice","Array","slice"},{"ArrayPrototypeSome","Array","some"},{"ArrayPrototypeSort","Array","sort"},{"ArrayPrototypeSplice","Array","splice"},{"ArrayPrototypeToSorted","Array","toSorted"},{"ArrayPrototypeUnshift","Array","unshift"},
			{"StringPrototypeCharAt","String","charAt"},{"StringPrototypeCharCodeAt","String","charCodeAt"},{"StringPrototypeCodePointAt","String","codePointAt"},{"StringPrototypeEndsWith","String","endsWith"},{"StringPrototypeIncludes","String","includes"},{"StringPrototypeIndexOf","String","indexOf"},{"StringPrototypeLastIndexOf","String","lastIndexOf"},{"StringPrototypeRepeat","String","repeat"},{"StringPrototypeReplaceAll","String","replaceAll"},{"StringPrototypeSlice","String","slice"},{"StringPrototypeSplit","String","split"},{"StringPrototypeStartsWith","String","startsWith"},{"StringPrototypeToLocaleLowerCase","String","toLocaleLowerCase"},{"StringPrototypeToLowerCase","String","toLowerCase"},{"StringPrototypeTrim","String","trim"},{"StringPrototypeTrimStart","String","trimStart"},
			{"RegExpPrototypeExec","RegExp","exec"},
		}{method,owner:=tsInternalReplPrimordialsPrototypeMethod(loop,entry.owner,entry.method);set(entry.name,tsInternalReplPrimordialsUncurried(loop,entry.name,method));_ = owner}
		set("DateNow",tsNamedNativeMethod("DateNow",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{date:=tsInternalReplPrimordialsGlobal(loop,"Date");return tsCallReceiver(loop,tsGet(loop,date,tsStringReference(tsStringUTF8("now"))),date)}))
		set("MathMaxApply",tsNamedNativeMethod("MathMaxApply",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{values:=[]tsValue{};it:=tsIterate(tsArg(args,0));for it.next(loop){values=append(values,it.value)};math:=tsInternalReplPrimordialsGlobal(loop,"Math");return tsInternalReplPrimordialsCallable(loop,tsGet(loop,math,tsStringReference(tsStringUTF8("max"))),math,values)}))
		set("ArrayPrototypePushApply",tsNamedNativeMethod("ArrayPrototypePushApply",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{target,values:=tsArg(args,0),tsArg(args,1);items:=[]tsValue{};it:=tsIterate(values);for it.next(loop){items=append(items,it.value)};method,_:=tsInternalReplPrimordialsPrototypeMethod(loop,"Array","push");return tsInternalReplPrimordialsCallable(loop,method,target,items)}))
		set("FunctionPrototypeBind",tsNamedNativeMethod("FunctionPrototypeBind",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{fn:=tsArg(args,0);bind:=tsGet(loop,fn,tsStringReference(tsStringUTF8("bind")));return tsCallReceiver(loop,bind,fn,tsRestArgs(args,1).values...)}))
		set("FunctionPrototypeCall",tsNamedNativeMethod("FunctionPrototypeCall",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{fn,thisArg:=tsArg(args,0),tsArg(args,1);return tsCallReceiver(loop,fn,thisArg,tsRestArgs(args,2).values...)}))
		set("ReflectApply",tsNamedNativeMethod("ReflectApply",3,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{fn,thisArg,values:=tsArg(args,0),tsArg(args,1),tsArg(args,2);items:=[]tsValue{};it:=tsIterate(values);for it.next(loop){items=append(items,it.value)};return tsCallReceiver(loop,fn,thisArg,items...)}))
		set("RegExpPrototypeSymbolReplace",tsInternalReplPrimordialsUncurried(loop,"RegExpPrototypeSymbolReplace",tsGet(loop,tsInternalReplPrimordialsGlobal(loop,"RegExp"),tsWellKnownSymbol("replace"))))
		set("RegExpPrototypeSymbolSplit",tsInternalReplPrimordialsUncurried(loop,"RegExpPrototypeSymbolSplit",tsGet(loop,tsInternalReplPrimordialsGlobal(loop,"RegExp"),tsWellKnownSymbol("split"))))
		set("SafePromiseRace",tsNamedNativeMethod("SafePromiseRace",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{promise:=tsInternalReplPrimordialsGlobal(loop,"Promise");race:=tsGet(loop,promise,tsStringReference(tsStringUTF8("race")));return tsCallReceiver(loop,race,promise,tsArg(args,0))}))
		// Captured methods in the shim are kept as aliases to the same source intrinsics.
		set("default",tsObjectValue(module));tsInternalReplPrimordials.value=tsObjectValue(module)
	})
	return tsInternalReplPrimordials.value
}
`
