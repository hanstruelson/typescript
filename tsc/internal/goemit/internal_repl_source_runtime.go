package goemit

// internalReplSourceRuntime runs the already-transpiled CommonJS bodies that
// Bun vendors for REPL helpers through the local VM and module dispatcher.
const internalReplSourceRuntime = `
func tsInternalReplEvaluateSource(loop *tsLoop,source,name,tail string)tsValue{
	wrapped:="(function(require){\\n"+source+"\\n"+tail+"\\n})"
	compiled:=tsInternalReplCompileSource(loop,wrapped,name)
	require:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeModuleRequire(loop,args)}))
	return tsCall(loop,compiled,require)
}
`
