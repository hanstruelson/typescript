package goemit

// internalStreamIterTypesRuntime holds the shared public protocols and private
// sentinels used by Bun's stream/iter bridges.
const internalStreamIterTypesRuntime = `
var tsInternalStreamIterTypes struct {
	sync.Once
	module *tsObject
	symbols map[string]tsValue
}

func tsInternalStreamIterTypesModule(loop *tsLoop) tsValue {
	tsInternalStreamIterTypes.Do(func() {
		registry := tsGet(loop,tsSymbolFunction(),tsStringReference(tsStringUTF8("for")))
		shared := map[string]string{
			"toStreamable":"Stream.toStreamable",
			"toAsyncStreamable":"Stream.toAsyncStreamable",
			"broadcastProtocol":"Stream.broadcastProtocol",
			"shareProtocol":"Stream.shareProtocol",
			"shareSyncProtocol":"Stream.shareSyncProtocol",
			"drainableProtocol":"Stream.drainableProtocol",
		}
		symbols:=map[string]tsValue{}
		for name,key:=range shared { symbols[name]=tsCall(loop,registry,tsStringReference(tsStringUTF8(key))) }
		for _,name:=range []string{"kValidatedTransform","kValidatedSource","kSyncWriteAccepted","kSyncWriteAcceptedOnFalse"} { symbols[name]=tsNewSymbol(tsStringUTF8(name)) }
		module:=tsNewObject()
		for name,value:=range symbols { module.set(name,value) }
		defaultValue:=tsNewObject()
		for name,value:=range symbols { defaultValue.set(name,value) }
		module.set("default",tsObjectValue(defaultValue))
		tsInternalStreamIterTypes.symbols=symbols
		tsInternalStreamIterTypes.module=module
	})
	return tsObjectValue(tsInternalStreamIterTypes.module)
}
`
