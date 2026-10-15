package goemit

// internalFSWatchRuntime maps Bun's watch modules onto the filesystem watcher
// runtime, which owns OS subscriptions and dispatches events on each JS loop.
const internalFSWatchRuntime = `
var tsInternalFSWatchSingleton struct { sync.Once; value tsValue }

func tsInternalFSWatchIgnoreMatcher(loop *tsLoop, value tsValue) tsValue {
	if tsNullish(value) { return tsNull }
	matchers := []tsValue{}
	if value.kind == tsArrayKind { it:=tsIterate(value); for it.next(loop) { matchers=append(matchers,it.value) } } else { matchers=append(matchers,value) }
	for _, matcher := range matchers {
		if matcher.kind == tsStringKind { if tsText(matcher)=="" { panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_VALUE","options.ignore must not be empty")}) }; continue }
		if matcher.kind == tsRegExpKind || matcher.kind == tsFunctionKind { continue }
		panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_TYPE","The options.ignore argument must be a string, RegExp, or Function")})
	}
	return tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{
		filename:=tsStringValueWithLoop(loop,tsArg(args,0)).String()
		base:=filepath.Base(filename)
		for _,matcher:=range matchers {
			matched:=false
			switch matcher.kind {
			case tsStringKind:
				pattern:=tsText(matcher);matched=tsFSGlobMatch(pattern,filename);if !matched&&!strings.Contains(pattern,"/"){matched=tsFSGlobMatch(pattern,base)}
			case tsRegExpKind:
				exec:=tsGet(loop,matcher,tsStringReference(tsStringUTF8("exec")));result:=tsCallReceiver(loop,exec,matcher,tsStringReference(tsStringUTF8(filename)));matched=!tsIsNull(result)&&!tsIsUndefined(result)
			case tsFunctionKind:
				matched=tsTruthy(tsCall(loop,matcher,tsStringReference(tsStringUTF8(filename))))
			}
			if matched{return tsBooleanValue(true)}
		}
		return tsBooleanValue(false)
	}))
}

func tsInternalFSWatchModule(loop *tsLoop) tsValue {
	tsInternalFSWatchSingleton.Do(func(){
		module:=tsNewObject()
		module.set("watch",tsNamedNativeMethod("watch",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsFSWatch(loop,args)}))
		module.set("createIgnoreMatcher",tsNamedNativeMethod("createIgnoreMatcher",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsInternalFSWatchIgnoreMatcher(loop,tsArg(args,0))}))
		module.set("default",tsObjectValue(module));tsInternalFSWatchSingleton.value=tsObjectValue(module)
	})
	return tsInternalFSWatchSingleton.value
}

`
