package goemit

// internalReadlineKeypressRuntime adapts readline's UTF-8 keypress event
// bridge and shares the interface symbol used to mark completed characters.
const internalReadlineKeypressRuntime = `
var tsInternalReadlineKeys struct { sync.Once; sawKeyPress, keypressDecoder, escapeDecoder tsValue; interfaceKeys map[string]tsValue }

func tsInternalReadlineKey(loop *tsLoop, name string) tsValue {
	tsInternalReadlineKeys.Do(func(){
		tsInternalReadlineKeys.interfaceKeys=map[string]tsValue{}
		tsInternalReadlineKeys.sawKeyPress=tsNewSymbol(tsStringUTF8("kSawKeyPress"))
		tsInternalReadlineKeys.keypressDecoder=tsNewSymbol(tsStringUTF8("keypress-decoder"))
		tsInternalReadlineKeys.escapeDecoder=tsNewSymbol(tsStringUTF8("escape-decoder"))
		tsInternalReadlineKeys.interfaceKeys["kSawKeyPress"]=tsInternalReadlineKeys.sawKeyPress
		tsInternalReadlineKeys.interfaceKeys["keypress-decoder"]=tsInternalReadlineKeys.keypressDecoder
		tsInternalReadlineKeys.interfaceKeys["escape-decoder"]=tsInternalReadlineKeys.escapeDecoder
	})
	if value:=tsInternalReadlineKeys.interfaceKeys[name];!tsIsUndefined(value){return value}
	value:=tsNewSymbol(tsStringUTF8(name));tsInternalReadlineKeys.interfaceKeys[name]=value;return value
}
func tsInternalReadlineCallMethod(loop *tsLoop, receiver tsValue, name string, args ...tsValue) tsValue {
	method:=tsGet(loop,receiver,tsStringReference(tsStringUTF8(name)))
	if method.kind!=tsFunctionKind { panic(tsThrown{tsFSException("TypeError","ERR_METHOD_NOT_IMPLEMENTED","stream."+name+" is not a function")}) }
	return tsCallReceiver(loop,method,receiver,args...)
}
func tsInternalReadlineKeypressModule(loop *tsLoop) tsValue {
	keys:=tsInternalReadlineKey(loop,"kSawKeyPress");keypressKey:=tsInternalReadlineKey(loop,"keypress-decoder");escapeKey:=tsInternalReadlineKey(loop,"escape-decoder")
	function:=tsFunc(func(loop *tsLoop,args ...tsValue) tsValue {
		stream:=tsArg(args,0);iface:=tsArg(args,1);if tsIsUndefined(iface){iface=tsObjectValue(tsNewObject())}
		if !tsIsUndefined(tsGet(loop,stream,keypressKey)){return tsU}
		decoderModule:=tsNodeStringDecoderModule();decoderConstructor:=tsGet(loop,decoderModule,tsStringReference(tsStringUTF8("StringDecoder")))
		decoder:=tsConstruct(decoderConstructor,tsStringReference(tsStringUTF8("utf8")))
		tsSet(loop,stream,keypressKey,decoder)
		utils:=tsInternalReadlineUtilsModule(loop);emitKeys:=tsGet(loop,utils,tsStringReference(tsStringUTF8("emitKeys")))
		escapeDecoder:=tsCall(loop,emitKeys,stream);next:=tsGet(loop,escapeDecoder,tsStringReference(tsStringUTF8("next")));tsCallReceiver(loop,next,escapeDecoder)
		tsSet(loop,stream,escapeKey,escapeDecoder)
		escapeTimeout:=tsGet(loop,iface,tsStringReference(tsStringUTF8("escapeCodeTimeout")));if tsIsUndefined(escapeTimeout){escapeTimeout=tsNumberValue(500)}
		var timeout *tsTask
		triggerEscape:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{return tsCallReceiver(loop,tsGet(loop,escapeDecoder,tsStringReference(tsStringUTF8("next"))),escapeDecoder,tsStringReference(tsStringUTF8("")))}))
		var onData,onNewListener tsValue
		reset:=func(){current:=tsCall(loop,emitKeys,stream);tsCallReceiver(loop,tsGet(loop,current,tsStringReference(tsStringUTF8("next"))),current);tsSet(loop,stream,escapeKey,current);escapeDecoder=current}
		onData=tsFunctionValue(tsFunc(func(loop *tsLoop,data ...tsValue)tsValue{
			listeners:=tsInternalReadlineCallMethod(loop,stream,"listenerCount",tsStringReference(tsStringUTF8("keypress")))
			if tsNumber(listeners)>0 {
				if timeout!=nil{loop.clearTimeout(tsTaskValue(timeout));timeout=nil}
				input:=tsInternalReadlineCallMethod(loop,decoder,"write",tsArg(data,0));text:=tsText(input)
				if text!="" {
					whole:=tsInternalReadlineCharLengthAt(loop,tsStringValue(input),0)==len(tsStringValue(input).units)
					tsSet(loop,iface,keys,tsBooleanValue(whole));tsSet(loop,iface,tsStringReference(tsStringUTF8("isCompletionEnabled")),tsBooleanValue(false))
					length:=0
					for _,character:=range text {
						piece:=string(character);length+=len(tsStringValue(tsStringReference(tsStringUTF8(piece))).units);if length==len(tsStringValue(input).units){tsSet(loop,iface,tsStringReference(tsStringUTF8("isCompletionEnabled")),tsBooleanValue(true))}
						var failure any
						func(){defer func(){failure=recover()}();tsCallReceiver(loop,next,escapeDecoder,tsStringReference(tsStringUTF8(piece)))}()
						if failure!=nil{reset();panic(failure)}
						if length==len(tsStringValue(input).units)&&piece=="\x1b"{timeout=loop.setTimeout(triggerEscape,escapeTimeout)}
					}
				}
			} else {
				tsInternalReadlineCallMethod(loop,stream,"removeListener",tsStringReference(tsStringUTF8("data")),onData)
				tsInternalReadlineCallMethod(loop,stream,"on",tsStringReference(tsStringUTF8("newListener")),onNewListener)
			}
			return tsU
		}))
		onNewListener=tsFunctionValue(tsFunc(func(loop *tsLoop,values ...tsValue)tsValue{
			if tsText(tsArg(values,0))=="keypress"{tsInternalReadlineCallMethod(loop,stream,"on",tsStringReference(tsStringUTF8("data")),onData);tsInternalReadlineCallMethod(loop,stream,"removeListener",tsStringReference(tsStringUTF8("newListener")),onNewListener)}
			return tsU
		}))
		if tsNumber(tsInternalReadlineCallMethod(loop,stream,"listenerCount",tsStringReference(tsStringUTF8("keypress"))))>0 {tsInternalReadlineCallMethod(loop,stream,"on",tsStringReference(tsStringUTF8("data")),onData)} else {tsInternalReadlineCallMethod(loop,stream,"on",tsStringReference(tsStringUTF8("newListener")),onNewListener)}
		return tsU
	})
	function.name="emitKeypressEvents";function.length=1
	return tsFunctionValue(function)
}
`
