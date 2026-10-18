package goemit

// nodeHTTPRuntime maps the public node:http facade onto the lower-level HTTP
// adapters. The https import remains lazy in setGlobalProxyFromEnv, matching
// Bun's wrapper and allowing the http/https module cycle.
const nodeHTTPRuntime = `
var tsNodeHTTPModuleSingleton struct { module tsValue }

func tsNodeHTTPKey(name string) tsValue { return tsStringReference(tsStringUTF8(name)) }
func tsNodeHTTPGet(loop *tsLoop,value tsValue,name string) tsValue { return tsGet(loop,value,tsNodeHTTPKey(name)) }
func tsNodeHTTPSet(loop *tsLoop,value tsValue,name string,next tsValue) { tsSet(loop,value,tsNodeHTTPKey(name),next) }
func tsNodeHTTPCall(loop *tsLoop,receiver tsValue,name string,args ...tsValue) tsValue { fn:=tsNodeHTTPGet(loop,receiver,name);if fn.kind!=tsFunctionKind{return tsU};return tsCallReceiver(loop,fn,receiver,args...) }

func tsNodeHTTPProxyConfigError(loop *tsLoop,internal tsValue,redact,proxy tsValue) {
	message:="Invalid proxy configuration: "+tsText(tsCall(loop,redact,proxy))
	err:=tsFSException("TypeError","ERR_PROXY_INVALID_CONFIG",message)
	panic(tsThrown{err})
}

func tsNodeHTTPNoopRestore(loop *tsLoop) tsValue {
	return tsNamedNativeMethod("restore",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsU})
}

func tsNodeHTTPSetGlobalProxyFromEnv(loop *tsLoop,args []tsValue) tsValue {
	env:=tsArg(args,0)
	if tsIsUndefined(env) { env=tsNodeHTTPGet(loop,tsNodeProcessModule(loop),"env") }
	validators:=tsInternalValidatorsModule(loop)
	tsCall(loop,tsNodeHTTPGet(loop,validators,"validateObject"),env,tsNodeHTTPKey("proxyEnv"))
	internal:=tsInternalHTTPModule(loop)
	parse:=tsNodeHTTPGet(loop,internal,"parseProxyUrl")
	redact:=tsNodeHTTPGet(loop,internal,"redactInvalidProxyUrl")
	httpProxy:=tsCall(loop,parse,env,tsNodeHTTPKey("http:"))
	httpsProxy:=tsCall(loop,parse,env,tsNodeHTTPKey("https:"))
	if tsIsNull(httpProxy)&&tsIsNull(httpsProxy) { return tsNodeHTTPNoopRestore(loop) }
	global:=tsGlobalObject()
	urlClass:=tsNodeHTTPGet(loop,global,"URL")
	canParse:=tsNodeHTTPGet(loop,urlClass,"canParse")
	if !tsIsNull(httpProxy)&&canParse.kind==tsFunctionKind&&!tsTruthy(tsCallReceiver(loop,canParse,urlClass,httpProxy)) { tsNodeHTTPProxyConfigError(loop,internal,redact,httpProxy) }
	if !tsIsNull(httpsProxy)&&canParse.kind==tsFunctionKind&&!tsTruthy(tsCallReceiver(loop,canParse,urlClass,httpsProxy)) { tsNodeHTTPProxyConfigError(loop,internal,redact,httpsProxy) }
	httpAgent:=tsNodeHTTPAgentModule(loop)
	originalHTTP,originalHTTPS:=tsU,tsU
	if !tsIsNull(httpProxy) {
		originalHTTP=tsNodeHTTPGet(loop,httpAgent,"globalAgent")
		agentOptions:=tsNewObject();agentOptions.set("keepAlive",tsBooleanValue(true));agentOptions.set("scheduling",tsNodeHTTPKey("lifo"));agentOptions.set("timeout",tsNumberValue(5000));agentOptions.set("proxyEnv",env)
		tsNodeHTTPSet(loop,httpAgent,"globalAgent",tsConstruct(tsNodeHTTPGet(loop,httpAgent,"Agent"),tsObjectValue(agentOptions)))
	}
	var httpsModule tsValue
	if !tsIsNull(httpsProxy) {
		httpsModule=tsNodeModuleRequire(loop,[]tsValue{tsNodeHTTPKey("https")})
		originalHTTPS=tsNodeHTTPGet(loop,httpsModule,"globalAgent")
		agentOptions:=tsNewObject();agentOptions.set("keepAlive",tsBooleanValue(true));agentOptions.set("scheduling",tsNodeHTTPKey("lifo"));agentOptions.set("timeout",tsNumberValue(5000));agentOptions.set("proxyEnv",env)
		tsNodeHTTPSet(loop,httpsModule,"globalAgent",tsConstruct(tsNodeHTTPGet(loop,httpsModule,"Agent"),tsObjectValue(agentOptions)))
	}
	return tsNamedNativeMethod("restore",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{
		if !tsIsUndefined(originalHTTP) { tsNodeHTTPSet(loop,httpAgent,"globalAgent",originalHTTP) }
		if !tsIsUndefined(originalHTTPS) { module:=tsNodeModuleRequire(loop,[]tsValue{tsNodeHTTPKey("https")});tsNodeHTTPSet(loop,module,"globalAgent",originalHTTPS) }
		return tsU
	})
}

func tsNodeHTTPModule(loop *tsLoop) tsValue {
	if !tsIsUndefined(tsNodeHTTPModuleSingleton.module) { return tsNodeHTTPModuleSingleton.module }
	validators:=tsInternalValidatorsModule(loop)
	agent:=tsNodeHTTPAgentModule(loop)
	client:=tsNodeHTTPClientModule(loop)
	common:=tsNodeHTTPCommonModule(loop)
	incoming:=tsNodeHTTPIncomingModule(loop)
	outgoing:=tsNodeHTTPOutgoingModule(loop)
	server:=tsNodeHTTPServerAPIModule(loop)
	internal:=tsInternalHTTPModule(loop)
	process:=tsNodeProcessModule(loop)
	if tsTruthy(tsNodeHTTPGet(loop,tsNodeHTTPGet(loop,process,"env"),"NODE_DEBUG")) {
		util:=tsNodeUtilModule(loop);debuglog:=tsNodeHTTPGet(loop,util,"debuglog");if debuglog.kind==tsFunctionKind { tsCall(loop,debuglog,tsNodeHTTPKey("http")) }
	}
	global:=tsGlobalObject()
	module:=tsNewObject()
	for _,entry:=range []struct{name string;value tsValue}{
		{"Agent",tsNodeHTTPGet(loop,agent,"Agent")},{"Server",tsNodeHTTPGet(loop,server,"Server")},{"METHODS",tsNodeHTTPGet(loop,internal,"METHODS")},{"STATUS_CODES",tsNodeHTTPGet(loop,internal,"STATUS_CODES")},
		{"ServerResponse",tsNodeHTTPGet(loop,server,"ServerResponse")},{"IncomingMessage",tsNodeHTTPGet(loop,incoming,"IncomingMessage")},{"ClientRequest",tsNodeHTTPGet(loop,client,"ClientRequest")},{"OutgoingMessage",tsNodeHTTPGet(loop,outgoing,"OutgoingMessage")},
		{"validateHeaderName",tsNodeHTTPGet(loop,common,"validateHeaderName")},{"validateHeaderValue",tsNodeHTTPGet(loop,common,"validateHeaderValue")},
		{"WebSocket",tsNodeHTTPGet(loop,global,"WebSocket")},{"CloseEvent",tsNodeHTTPGet(loop,global,"CloseEvent")},{"MessageEvent",tsNodeHTTPGet(loop,global,"MessageEvent")},
	} { module.set(entry.name,entry.value) }
	createServer:=tsNamedNativeMethod("createServer",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsConstruct(tsNodeHTTPGet(loop,server,"Server"),args...)})
	request:=tsNamedNativeMethod("request",3,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsConstruct(tsNodeHTTPGet(loop,client,"ClientRequest"),args...)})
	get:=tsNamedNativeMethod("get",3,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{req:=tsCall(loop,request,args...);tsNodeHTTPCall(loop,req,"end");return req})
	module.set("createServer",createServer);module.set("request",request);module.set("get",get)
	module.set("setMaxIdleHTTPParsers",tsNamedNativeMethod("setMaxIdleHTTPParsers",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{
		max:=tsArg(args,0);tsCall(loop,tsNodeHTTPGet(loop,validators,"validateInteger"),max,tsNodeHTTPKey("max"),tsNumberValue(1));tsNodeHTTPSet(loop,tsNodeHTTPGet(loop,common,"parsers"),"max",max);return tsU
	}))
	module.set("setGlobalProxyFromEnv",tsNamedNativeMethod("setGlobalProxyFromEnv",0,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeHTTPSetGlobalProxyFromEnv(loop,args)}))
	module.descriptors=map[string]*tsDescriptor{
		"maxHeaderSize":{accessor:true,enumerable:true,configurable:true,getter:tsNamedNativeMethod("get maxHeaderSize",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNodeHTTPCall(loop,internal,"getMaxHTTPHeaderSize")}),setter:tsNamedNativeMethod("set maxHeaderSize",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{tsNodeHTTPCall(loop,internal,"setMaxHTTPHeaderSize",tsArg(args,0));return tsU})},
		"globalAgent":{accessor:true,enumerable:true,configurable:true,getter:tsNamedNativeMethod("get globalAgent",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNodeHTTPGet(loop,tsNodeHTTPAgentModule(loop),"globalAgent")}),setter:tsNamedNativeMethod("set globalAgent",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{tsNodeHTTPSet(loop,tsNodeHTTPAgentModule(loop),"globalAgent",tsArg(args,0));return tsU})},
	}
	module.set("default",tsObjectValue(module));tsNodeHTTPModuleSingleton.module=tsObjectValue(module);return tsNodeHTTPModuleSingleton.module
}
`
