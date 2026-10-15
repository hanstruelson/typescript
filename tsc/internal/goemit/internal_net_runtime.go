package goemit

// internalNetRuntime supplies the shared IP and socket-timer helpers used by
// Bun's Node net, DNS, and TLS wrappers.
const internalNetRuntime = `
var tsInternalNetSingleton struct { sync.Once; ip, timers, symbols, serverAbort tsValue; timeoutSymbol tsValue }
func tsInternalNetKey(name string) tsValue { return tsStringReference(tsStringUTF8(name)) }
func tsInternalNetGet(loop *tsLoop,value tsValue,name string) tsValue { return tsGet(loop,value,tsInternalNetKey(name)) }
func tsInternalNetIPv4(value tsValue) bool {
	if value.kind!=tsStringKind{return false};text:=tsText(value);if strings.Contains(text,":"){return false};ip:=net.ParseIP(text);return ip!=nil&&ip.To4()!=nil
}
func tsInternalNetIPv6(value tsValue) bool {
	if value.kind!=tsStringKind{return false};text:=tsText(value);if !strings.Contains(text,":"){return false}
	if i:=strings.LastIndexByte(text,'%');i>=0{zone:=text[i+1:];if zone==""{return false};for _,r:=range zone{if !(r>='0'&&r<='9'||r>='a'&&r<='z'||r>='A'&&r<='Z'||strings.ContainsRune("-.:",r)){return false}};text=text[:i]}
	return net.ParseIP(text)!=nil
}
func tsInternalNetIPModule() tsValue {
	if !tsIsUndefined(tsInternalNetSingleton.ip){return tsInternalNetSingleton.ip};module:=tsNewObject()
	module.set("isIPv4",tsNamedNativeMethod("isIPv4",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsBooleanValue(tsInternalNetIPv4(tsArg(args,0)))}))
	module.set("isIPv6",tsNamedNativeMethod("isIPv6",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsBooleanValue(tsInternalNetIPv6(tsArg(args,0)))}))
	module.set("isIP",tsNamedNativeMethod("isIP",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{if tsInternalNetIPv4(tsArg(args,0)){return tsNumberValue(4)};if tsInternalNetIPv6(tsArg(args,0)){return tsNumberValue(6)};return tsNumberValue(0)}))
	module.set("default",tsObjectValue(module));tsInternalNetSingleton.ip=tsObjectValue(module);return tsInternalNetSingleton.ip
}
func tsInternalNetTimersModule(loop *tsLoop) tsValue {
	if !tsIsUndefined(tsInternalNetSingleton.timers){return tsInternalNetSingleton.timers};module:=tsNewObject();symbolFor:=tsGet(loop,tsSymbolFunction(),tsInternalNetKey("for"));tsInternalNetSingleton.timeoutSymbol=tsCall(loop,symbolFor,tsStringReference(tsStringUTF8("::buntimeout::")))
	module.set("kTimeout",tsInternalNetSingleton.timeoutSymbol);module.set("getTimerDuration",tsNamedNativeMethod("getTimerDuration",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{value,name:=tsArg(args,0),tsText(tsArg(args,1));if !tsIsNumeric(value){tsInternalValidatorInvalidType(name,"a number",value)};number:=tsNumber(value);if number<0||math.IsNaN(number)||math.IsInf(number,0){tsInternalValidatorOutOfRange(name,"a non-negative finite number",value)};const max=2147483647;if number>max{process:=tsNodeProcessModule(loop);emit:=tsGet(loop,process,tsInternalNetKey("emitWarning"));if emit.kind==tsFunctionKind{message:=tsStringReference(tsStringUTF8(strconv.FormatFloat(number,'f',-1,64)+" does not fit into a 32-bit signed integer.\nTimer duration was truncated to 2147483647."));tsCallReceiver(loop,emit,process,message,tsStringReference(tsStringUTF8("TimeoutOverflowWarning")))};return tsNumberValue(max)};return value}))
	module.set("default",tsObjectValue(module));tsInternalNetSingleton.timers=tsObjectValue(module);return tsInternalNetSingleton.timers
}
func tsInternalNetSymbolsModule() tsValue {
	if !tsIsUndefined(tsInternalNetSingleton.symbols){return tsInternalNetSingleton.symbols};module:=tsNewObject();for _,name:=range []string{"kArmHandshakeTimeout","kDestroyOnRead","kPreHandshakeWrite","kSecureConnectDone","kVerifyError","kUpgradeClientTLS","kStandaloneWrap"}{module.set(name,tsNewSymbol(tsStringUTF8(name)))};module.set("default",tsObjectValue(module));tsInternalNetSingleton.symbols=tsObjectValue(module);return tsInternalNetSingleton.symbols
}
func tsInternalNetServerAbortModule(loop *tsLoop) tsValue {
	if !tsIsUndefined(tsInternalNetSingleton.serverAbort){return tsInternalNetSingleton.serverAbort};module:=tsNewObject();module.set("addServerAbortSignalOption",tsNamedNativeMethod("addServerAbortSignalOption",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{server,options:=tsArg(args,0),tsArg(args,1);signal:=tsInternalNetGet(loop,options,"signal");if tsIsUndefined(signal){return tsU};tsInternalValidatorAbortSignal(signal,tsStringReference(tsStringUTF8("options.signal")));onAborted:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{close:=tsInternalNetGet(loop,server,"close");if close.kind==tsFunctionKind{tsCallReceiver(loop,close,server)};return tsU}));if tsTruthy(tsInternalNetGet(loop,signal,"aborted")){process:=tsNodeProcessModule(loop);tsCallReceiver(loop,tsInternalNetGet(loop,process,"nextTick"),process,onAborted);return tsU};disposable:=tsInternalAbortListenerAdd(loop,signal,onAborted);dispose:=tsGet(loop,disposable,tsWellKnownSymbol("dispose"));once:=tsInternalNetGet(loop,server,"once");if once.kind==tsFunctionKind{tsCallReceiver(loop,once,server,tsStringReference(tsStringUTF8("close")),dispose)};return tsU}));module.set("default",tsObjectValue(module));tsInternalNetSingleton.serverAbort=tsObjectValue(module);return tsInternalNetSingleton.serverAbort
}
`
