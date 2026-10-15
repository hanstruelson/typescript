package goemit

// internalReplModeRuntime preserves the two identity-based mode sentinels.
const internalReplModeRuntime = `
var tsInternalReplMode struct { sync.Once; sloppy,strict,module tsValue }
func tsInternalReplModeModule(loop *tsLoop) tsValue {
	tsInternalReplMode.Do(func(){module:=tsNewObject();tsInternalReplMode.sloppy=tsNewSymbol(tsStringUTF8("repl-sloppy"));tsInternalReplMode.strict=tsNewSymbol(tsStringUTF8("repl-strict"));module.set("REPL_MODE_SLOPPY",tsInternalReplMode.sloppy);module.set("REPL_MODE_STRICT",tsInternalReplMode.strict);module.set("default",tsObjectValue(module));tsInternalReplMode.module=tsObjectValue(module)})
	return tsInternalReplMode.module
}
`
