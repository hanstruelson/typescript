package goemit

// internalUtilColorsRuntime translates internal/util/colors.ts.
const internalUtilColorsRuntime = `
var tsInternalUtilColorsState struct{sync.Once;module *tsObject}
func tsInternalUtilColorsModule(loop *tsLoop)tsValue{
	tsInternalUtilColorsState.Do(func(){out:=tsNewObject();tsInternalUtilColorsState.module=out;for _,name:=range []string{"blue","green","white","red"}{out.set(name,tsStringReference(tsStringUTF8("")))};out.set("hasColors",tsBooleanValue(false))
		out.set("shouldColorize",tsNamedNativeMethod("shouldColorize",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{stream:=tsArg(args,0);env:=tsGet(loop,tsNodeProcessModule(loop),tsStringReference(tsStringUTF8("env")));force:=tsGet(loop,env,tsStringReference(tsStringUTF8("FORCE_COLOR")));if !tsIsUndefined(force){depth:=tsInternalTTYGetColorDepth(loop,env);return tsBooleanValue(tsNumber(depth)>2)};if !tsTruthy(tsGet(loop,stream,tsStringReference(tsStringUTF8("isTTY")))){return tsBooleanValue(false)};getDepth:=tsGet(loop,stream,tsStringReference(tsStringUTF8("getColorDepth")));if getDepth.kind==tsFunctionKind{return tsBooleanValue(tsNumber(tsCallReceiver(loop,getDepth,stream))>2)};return tsBooleanValue(true)}))
		out.set("refresh",tsNamedNativeMethod("refresh",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{process:=tsNodeProcessModule(loop);stderr:=tsGet(loop,process,tsStringReference(tsStringUTF8("stderr")));should:=tsGet(loop,tsObjectValue(out),tsStringReference(tsStringUTF8("shouldColorize")));colored:=tsTruthy(tsCallReceiver(loop,should,tsObjectValue(out),stderr));codes:=map[string]string{"blue":"34","green":"32","white":"39","red":"31"};for name,code:=range codes{value:="";if colored{value="\x1b["+code+"m"};out.set(name,tsStringReference(tsStringUTF8(value)))};out.set("hasColors",tsBooleanValue(colored));return tsU}))
		refresh:=tsGet(loop,tsObjectValue(out),tsStringReference(tsStringUTF8("refresh")));tsCallReceiver(loop,refresh,tsObjectValue(out));out.set("default",tsObjectValue(out))
	});return tsObjectValue(tsInternalUtilColorsState.module)
}
`
