package goemit

// internalAssertUtilsRuntime keeps node:test's innerOk dependency routed
// through the assertion adapter already registered in the runtime.
const internalAssertUtilsRuntime = `
func tsInternalAssertUtilsModule(loop *tsLoop)tsValue{module:=tsNewObject();module.set("innerOk",tsNamedNativeMethod("innerOk",4,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{value:=tsArg(args,2);if !tsTruthy(value){assertions:=tsNodeAssertModule();ok:=tsGet(loop,assertions,tsStringReference(tsStringUTF8("ok")));callArgs:=[]tsValue{value};if len(args)>3{callArgs=append(callArgs,args[3])};return tsCallReceiver(loop,ok,assertions,callArgs...) };return tsU}));module.set("default",tsObjectValue(module));return tsObjectValue(module)}
`
