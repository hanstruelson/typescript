package goemit

// internalReadlineCallbacksRuntime maps cursor movement and terminal clearing
// onto stream.write while preserving Node's callback scheduling and errors.
const internalReadlineCallbacksRuntime = `
var tsInternalReadlineCallbacks struct { sync.Once; module tsValue }

func tsInternalReadlineCallbacksValidate(loop *tsLoop, callback tsValue) {
	if tsIsUndefined(callback) { return }
	validate:=tsGet(loop,tsInternalValidatorsModule(loop),tsStringReference(tsStringUTF8("validateFunction")))
	tsCall(loop,validate,callback,tsStringReference(tsStringUTF8("callback")))
}
func tsInternalReadlineCallbacksNextTick(loop *tsLoop, callback tsValue) {
	if callback.kind!=tsFunctionKind { return }
	loop.ready=append(loop.ready,func(){tsCallReceiver(loop,callback,tsU,tsNull)})
}
func tsInternalReadlineCallbacksWrite(loop *tsLoop, stream tsValue, data tsValue, callback tsValue) tsValue {
	write:=tsGet(loop,stream,tsStringReference(tsStringUTF8("write")))
	if write.kind!=tsFunctionKind { tsPropertyFailure("stream.write is not a function") }
	if tsIsUndefined(callback) { return tsCallReceiver(loop,write,stream,data) }
	return tsCallReceiver(loop,write,stream,data,callback)
}
func tsInternalReadlineCallbacksCSI(_ *tsLoop, values ...string) tsValue {
	return tsStringReference(tsStringUTF8("\x1b["+strings.Join(values,"")))
}
func tsInternalReadlineCallbacksNumber(value tsValue) string {
	return tsNumberText(tsNumber(value))
}
func tsInternalReadlineCallbacksModule(loop *tsLoop) tsValue {
	tsInternalReadlineCallbacks.Do(func(){
		module:=tsNewObject()
		module.set("cursorTo",tsNamedNativeMethod("cursorTo",4,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{
			stream,x,y,callback:=tsArg(args,0),tsArg(args,1),tsArg(args,2),tsArg(args,3)
			if !tsIsUndefined(callback) { tsInternalReadlineCallbacksValidate(loop,callback) }
			if y.kind==tsFunctionKind { callback=y;y=tsU }
			codes:=tsGet(loop,tsInternalReplNodeErrorsModule(loop),tsStringReference(tsStringUTF8("codes")))
			if tsIsNumeric(x)&&math.IsNaN(tsNumber(x)) { panic(tsThrown{tsConstruct(tsGet(loop,codes,tsStringReference(tsStringUTF8("ERR_INVALID_ARG_VALUE"))),tsStringReference(tsStringUTF8("x")),x)}) }
			if tsIsNumeric(y)&&math.IsNaN(tsNumber(y)) { panic(tsThrown{tsConstruct(tsGet(loop,codes,tsStringReference(tsStringUTF8("ERR_INVALID_ARG_VALUE"))),tsStringReference(tsStringUTF8("y")),y)}) }
			if tsNullish(stream)||!tsIsNumeric(x)&&!tsIsNumeric(y) { tsInternalReadlineCallbacksNextTick(loop,callback);return tsBooleanValue(true) }
			if !tsIsNumeric(x) { panic(tsThrown{tsConstruct(tsGet(loop,codes,tsStringReference(tsStringUTF8("ERR_INVALID_CURSOR_POS"))))}) }
			var data tsValue
			if !tsIsNumeric(y) { data=tsInternalReadlineCallbacksCSI(loop,tsInternalReadlineCallbacksNumber(tsNumberValue(tsNumber(x)+1))+"G") } else { data=tsInternalReadlineCallbacksCSI(loop,tsInternalReadlineCallbacksNumber(tsNumberValue(tsNumber(y)+1))+";"+tsInternalReadlineCallbacksNumber(tsNumberValue(tsNumber(x)+1))+"H") }
			return tsInternalReadlineCallbacksWrite(loop,stream,data,callback)
		}))
		module.set("moveCursor",tsNamedNativeMethod("moveCursor",4,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{
			stream,dx,dy,callback:=tsArg(args,0),tsArg(args,1),tsArg(args,2),tsArg(args,3)
			if !tsIsUndefined(callback) { tsInternalReadlineCallbacksValidate(loop,callback) }
			if tsNullish(stream)||!tsTruthy(dx)&&!tsTruthy(dy) { tsInternalReadlineCallbacksNextTick(loop,callback);return tsBooleanValue(true) }
			var out strings.Builder
			if tsNumber(dx)<0 { out.WriteString(tsText(tsInternalReadlineCallbacksCSI(loop,tsNumberText(-tsNumber(dx))+"D"))) } else if tsNumber(dx)>0 { out.WriteString(tsText(tsInternalReadlineCallbacksCSI(loop,tsNumberText(tsNumber(dx))+"C"))) }
			if tsNumber(dy)<0 { out.WriteString(tsText(tsInternalReadlineCallbacksCSI(loop,tsNumberText(-tsNumber(dy))+"A"))) } else if tsNumber(dy)>0 { out.WriteString(tsText(tsInternalReadlineCallbacksCSI(loop,tsNumberText(tsNumber(dy))+"B"))) }
			return tsInternalReadlineCallbacksWrite(loop,stream,tsStringReference(tsStringUTF8(out.String())),callback)
		}))
		module.set("clearLine",tsNamedNativeMethod("clearLine",3,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{
			stream,dir,callback:=tsArg(args,0),tsArg(args,1),tsArg(args,2)
			if !tsIsUndefined(callback) { tsInternalReadlineCallbacksValidate(loop,callback) }
			if tsNullish(stream) { tsInternalReadlineCallbacksNextTick(loop,callback);return tsBooleanValue(true) }
			code:="2K";if tsNumber(dir)<0 { code="1K" } else if tsNumber(dir)>0 { code="0K" }
			return tsInternalReadlineCallbacksWrite(loop,stream,tsInternalReadlineCallbacksCSI(loop,code),callback)
		}))
		module.set("clearScreenDown",tsNamedNativeMethod("clearScreenDown",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{
			stream,callback:=tsArg(args,0),tsArg(args,1)
			if !tsIsUndefined(callback) { tsInternalReadlineCallbacksValidate(loop,callback) }
			if tsNullish(stream) { tsInternalReadlineCallbacksNextTick(loop,callback);return tsBooleanValue(true) }
			return tsInternalReadlineCallbacksWrite(loop,stream,tsInternalReadlineCallbacksCSI(loop,"0J"),callback)
		}))
		module.set("default",tsObjectValue(module));tsInternalReadlineCallbacks.module=tsObjectValue(module)
	})
	return tsInternalReadlineCallbacks.module
}
`
