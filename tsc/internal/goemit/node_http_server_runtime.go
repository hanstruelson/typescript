package goemit

// nodeHTTPServerRuntime implements the Node HTTP server API on the dedicated
// HTTP/1 connection and response engine in nodeHTTP1ServerRuntime.
const nodeHTTPServerRuntime = `
var tsNodeHTTPServerSingleton struct {
	sync.Once
	symbols map[string]tsValue
	parserTimeoutKey tsValue
	traceEvents tsValue
	socketConstructor tsValue
}

func tsNodeHTTPServerKey(name string) tsValue { return tsStringReference(tsStringUTF8(name)) }
func tsNodeHTTPServerGet(loop *tsLoop, value tsValue, name string) tsValue { return tsGet(loop, value, tsNodeHTTPServerKey(name)) }
func tsNodeHTTPServerSet(loop *tsLoop, value tsValue, name string, next tsValue) { tsSet(loop, value, tsNodeHTTPServerKey(name), next) }
func tsNodeHTTPServerSym(loop *tsLoop, value tsValue, name string) tsValue { return tsGet(loop, value, tsNodeHTTPServerSingleton.symbols[name]) }
func tsNodeHTTPServerPut(loop *tsLoop, value tsValue, name string, next tsValue) { tsSet(loop, value, tsNodeHTTPServerSingleton.symbols[name], next) }
func tsNodeHTTPServerCall(loop *tsLoop, value tsValue, name string, args ...tsValue) tsValue {
	method := tsNodeHTTPServerGet(loop, value, name)
	if method.kind == tsFunctionKind { return tsCallReceiver(loop, method, value, args...) }
	return tsU
}
func tsNodeHTTPServerEmit(loop *tsLoop, value tsValue, event string, args ...tsValue) tsValue {
	return tsNodeHTTPServerCall(loop, value, "emit", append([]tsValue{tsNodeHTTPServerKey(event)}, args...)...)
}
func tsNodeHTTPServerInit(loop *tsLoop) {
	tsNodeHTTPServerSingleton.Do(func() {
		internal := tsInternalHTTPModule(loop)
		tsNodeHTTPServerSingleton.symbols = map[string]tsValue{}
		for _, name := range []string{"kServerResponse", "kChunkedEncoding", "kShouldKeepAlive", "kOptimizeEmptyRequests", "kConnectionsCheckingInterval", "kTrackedConnections", "kPendingDrainClose", "kPendingCloseGenerations", "kListenerGeneration", "kHandedOff", "kHttpAllowHalfOpen", "kMaxHeadersCount", "kPendingFinish", "kBytesBuffered", "kAccountingFlushScheduled"} {
			tsNodeHTTPServerSingleton.symbols[name] = tsNewSymbol(tsStringUTF8("http.server." + name))
		}
		for _, name := range []string{"kInternalSocketData", "serverSymbol", "kHandle", "kOnReadParsed", "kHandoffResponse", "kRealListen", "tlsSymbol", "optionsSymbol", "headerStateSymbol", "kPendingCallbacks", "kRequest", "kCloseCallback", "fakeSocketSymbol", "kOutHeaders", "kNeedDrain"} {
			tsNodeHTTPServerSingleton.symbols[name] = tsInternalHTTPGet(loop, internal, name)
		}
		common := tsNodeHTTPCommonModule(loop)
		tsNodeHTTPServerSingleton.symbols["kIncomingMessage"] = tsNodeHTTPCommonGet(loop, common, "kIncomingMessage")
		parser := tsNodeHTTPCommonGet(loop, common, "HTTPParser")
		tsNodeHTTPServerSingleton.parserTimeoutKey = tsNodeHTTPServerGet(loop, parser, "kOnTimeout")
	})
}
func tsNodeHTTPServerCloseHandle(loop *tsLoop, handle tsValue) {
	if !tsTruthy(tsNodeHTTPServerGet(loop, handle, "closed")) { tsNodeHTTPServerCall(loop, handle, "close") }
}
func tsNodeHTTPServerTraceStart(loop *tsLoop, response tsValue) {
	if tsIsUndefined(tsNodeHTTPServerSingleton.traceEvents) { tsNodeHTTPServerSingleton.traceEvents = tsNodeModuleRequire(loop, []tsValue{tsNodeHTTPServerKey("internal/trace_events")}) }
	if !tsTruthy(tsNodeHTTPServerCall(loop, tsNodeHTTPServerSingleton.traceEvents, "isCategoryGroupEnabled", tsNodeHTTPServerKey("node,node.http"))) { return }
	tsNodeHTTPServerCall(loop, tsNodeHTTPServerSingleton.traceEvents, "emitEvent", tsNodeHTTPServerKey("b"), tsNodeHTTPServerKey("node,node.http"), tsNodeHTTPServerKey("http.server.request"))
	tsNodeHTTPServerCall(loop, response, "once", tsNodeHTTPServerKey("finish"), tsFunctionValue(tsFunc(func(loop *tsLoop, _ ...tsValue) tsValue { tsNodeHTTPServerTraceEnd(loop); return tsU })))
}
func tsNodeHTTPServerTraceEnd(loop *tsLoop) {
	if tsIsUndefined(tsNodeHTTPServerSingleton.traceEvents) { tsNodeHTTPServerSingleton.traceEvents = tsNodeModuleRequire(loop, []tsValue{tsNodeHTTPServerKey("internal/trace_events")}) }
	tsNodeHTTPServerCall(loop, tsNodeHTTPServerSingleton.traceEvents, "emitEvent", tsNodeHTTPServerKey("e"), tsNodeHTTPServerKey("node,node.http"), tsNodeHTTPServerKey("http.server.request"))
}
func tsNodeHTTPServerCollectionEntries(value tsValue) []tsCollectionEntry {
	if value.kind != tsCollectionKind { return nil }
	return (*tsCollection)(value.ref).entries
}
func tsNodeHTTPServerEmitClose(loop *tsLoop, server tsValue) {
	if !tsTruthy(tsNodeHTTPServerSym(loop, server, "kPendingDrainClose")) { return }
	pending := tsNodeHTTPServerSym(loop, server, "kPendingCloseGenerations")
	drainedAtClose := true
	for _, entry := range tsNodeHTTPServerCollectionEntries(pending) {
		if !entry.live { continue }
		generation := entry.key
		if !tsTruthy(tsNodeHTTPServerGet(loop, generation, "nativeClosed")) { return }
		drainedAtClose = drainedAtClose && tsTruthy(tsNodeHTTPServerGet(loop, generation, "drainedAtClose"))
	}
	relistened := !tsNullish(tsNodeHTTPServerSym(loop, server, "serverSymbol"))
	tracked := tsNodeHTTPServerSym(loop, server, "kTrackedConnections")
	connectionCount := 0
	if tracked.kind == tsCollectionKind { connectionCount = len((*tsCollection)(tracked.ref).index) }
	if (relistened || connectionCount > 0) && !drainedAtClose { return }
	tsNodeHTTPServerPut(loop, server, "kPendingDrainClose", tsBooleanValue(false))
	clear := tsNodeHTTPServerGet(loop, pending, "clear")
	if clear.kind == tsFunctionKind { tsCallReceiver(loop, clear, pending) }
	if !relistened { tsNodeHTTPServerPut(loop, server, "kListenerGeneration", tsU) }
	tsNodeHTTPServerEmit(loop, server, "close")
}
func tsNodeHTTPServerMarkListenerClosed(loop *tsLoop, server, generation tsValue) {
	tsNodeHTTPServerSet(loop, generation, "nativeClosed", tsBooleanValue(true))
	tsNodeHTTPServerEmitClose(loop, server)
}
func tsNodeHTTPServerSetCloseCallback(loop *tsLoop, self, callback tsValue) {
	current := tsNodeHTTPServerSym(loop, self, "kCloseCallback")
	if tsIsUndefined(callback) {
		tsNodeHTTPServerPut(loop, self, "kCloseCallback", tsU)
		return
	}
	if tsStrictEqual(callback, current) { return }
	if tsTruthy(current) { panic(tsThrown{tsErrorValue(&tsRuntimeError{name: "Error", message: "Close callback already set"})}) }
	tsNodeHTTPServerPut(loop, self, "kCloseCallback", callback)
}
func tsNodeHTTPServerAssignSocket(loop *tsLoop, self, socket tsValue) {
	if tsTruthy(tsNodeHTTPServerGet(loop, socket, "_httpMessage")) { panic(tsThrown{tsFSException("Error", "ERR_HTTP_SOCKET_ASSIGNED", "Socket is already assigned to a request")}) }
	tsNodeHTTPServerSet(loop, socket, "_httpMessage", self)
	tsNodeHTTPServerSetCloseCallback(loop, socket, tsNodeHTTPServerFn(loop, "onServerResponseClose"))
	tsNodeHTTPServerSet(loop, self, "socket", socket)
	tsNodeHTTPServerEmit(loop, self, "socket", socket)
}
func tsNodeHTTPServerFn(loop *tsLoop, name string) tsValue {
	return tsNamedNativeMethod(name, 2, func(loop *tsLoop, receiver tsValue, args ...tsValue) tsValue {
		return tsNodeHTTPServerDispatch(loop, name, receiver, args)
	})
}
func tsNodeHTTPServerDispatch(loop *tsLoop, name string, receiver tsValue, args []tsValue) tsValue {
	arg := func(index int) tsValue { return tsArg(args, index) }
	switch name {
	case "emitCloseServer": tsNodeHTTPServerEmitClose(loop, receiver); return tsU
	case "markListenerGenerationClosed": tsNodeHTTPServerMarkListenerClosed(loop, receiver, arg(0)); return tsU
	case "emitCloseNTServer":
		server, generation := receiver, arg(0)
		callback := tsFunctionValue(tsFunc(func(loop *tsLoop, _ ...tsValue) tsValue { tsNodeHTTPServerDispatch(loop, "markListenerGenerationClosed", server, []tsValue{generation}); return tsU }))
		tsInternalHTTPNextTick(loop, callback); return tsU
	case "onServerResponseClose":
		message := tsNodeHTTPServerGet(loop, receiver, "_httpMessage")
		if !tsNullish(message) { tsNodeHTTPServerCall(loop, tsInternalHTTPModule(loop), "emitCloseNT", message) }
		return tsU
	case "onNodeHTTPServerSocketTimeout":
		request := tsNodeHTTPServerSym(loop, receiver, "kRequest")
		requestTimedOut := !tsNullish(request) && !tsTruthy(tsNodeHTTPServerGet(loop, request, "complete")) && tsTruthy(tsNodeHTTPServerEmit(loop, request, "timeout", receiver))
		response := tsNodeHTTPServerGet(loop, receiver, "_httpMessage")
		responseTimedOut := !tsNullish(response) && tsTruthy(tsNodeHTTPServerEmit(loop, response, "timeout", receiver))
		server := tsNodeHTTPServerGet(loop, receiver, "server")
		serverTimedOut := tsTruthy(tsNodeHTTPServerEmit(loop, server, "timeout", receiver))
		if !requestTimedOut && !responseTimedOut && !serverTimedOut { tsNodeHTTPServerCall(loop, receiver, "destroy") }
		return tsU
	case "emitListeningNextTick":
		server := arg(0)
		if tsNullish(tsNodeHTTPServerSym(loop, server, "serverSymbol")) { return tsU }
		shared := tsInternalSharedModule(loop)
		stopped := tsNodeHTTPServerCall(loop, shared, "isStoppedModuleGraphRunning")
		if tsTruthy(stopped) { return tsU }
		return tsNodeHTTPServerEmit(loop, server, "listening", tsNull, arg(1), arg(2))
	case "emitListenErrorNextTick": return tsNodeHTTPServerEmit(loop, arg(0), "error", arg(1))
	case "normalizeServerTls":
		tls := receiver
		tsNodeHTTPServerSet(loop, tls, "requestCert", tsBooleanValue(tsStrictEqual(tsNodeHTTPServerGet(loop, tls, "requestCert"), tsBooleanValue(true))))
		tsNodeHTTPServerSet(loop, tls, "rejectUnauthorized", tsBooleanValue(!tsStrictEqual(tsNodeHTTPServerGet(loop, tls, "rejectUnauthorized"), tsBooleanValue(false))))
		return tls
	case "serverParserShimFree", "serverParserShimClose", "serverParserShimOnTimeout": return tsU
	case "createServerParserShim":
		object := tsNewObject()
		object.set("incoming", tsNull); object.set("outgoing", tsNull); object.set("socket", arg(0))
		object.set("free", tsNodeHTTPServerFn(loop, "serverParserShimFree")); object.set("close", tsNodeHTTPServerFn(loop, "serverParserShimClose"))
		object.set(tsPropertyKey(tsNodeHTTPServerSingleton.parserTimeoutKey), tsNodeHTTPServerFn(loop, "serverParserShimOnTimeout"))
		return tsObjectValue(object)
	case "releaseServerParserShim":
		socket, request := receiver, arg(0); parser := tsNodeHTTPServerGet(loop, socket, "parser")
		if tsNullish(parser) { return tsU }
		tsNodeHTTPServerSet(loop, socket, "parser", tsNull); tsNodeHTTPServerSet(loop, parser, "incoming", tsNull); tsNodeHTTPServerSet(loop, parser, "socket", tsNull)
		tsSet(loop, parser, tsNodeHTTPServerSingleton.parserTimeoutKey, tsNull)
		tsNodeHTTPServerCall(loop, parser, "free")
		if !tsNullish(request) && !tsNullish(tsNodeHTTPServerGet(loop, request, "parser")) { tsNodeHTTPServerSet(loop, request, "parser", tsNull) }
		return tsU
	case "strictContentLength":
		response, headerState, fromEnd := receiver, arg(0), tsTruthy(arg(1))
		if !tsTruthy(tsNodeHTTPServerGet(loop, response, "strictContentLength")) { return tsU }
		none := tsInternalHTTPGet(loop, tsInternalHTTPModule(loop), "NodeHTTPHeaderState")
		if !fromEnd && tsStrictEqual(headerState, tsNodeHTTPServerGet(loop, none, "none")) { return tsU }
		contentLength := tsNodeHTTPServerGet(loop, response, "_contentLength"); if tsNullish(contentLength) { contentLength = tsNodeHTTPOutgoingHeaderGet(loop, response, tsNodeHTTPServerKey("content-length")) }
		if tsNullish(contentLength) || !tsTruthy(tsNodeHTTPServerGet(loop, response, "_hasBody")) || tsTruthy(tsNodeHTTPServerGet(loop, response, "chunkedEncoding")) || tsNodeHTTPOutgoingHeaderGet(loop, response, tsNodeHTTPServerKey("transfer-encoding")).kind != tsUndefinedKind { return tsU }
		if !tsIsNumeric(contentLength) { number := tsNumber(tsToPrimitive(contentLength)); if math.IsNaN(number) { return tsU }; return tsNumberValue(number) }
		return contentLength
	case "emitCloseServerEvent": tsNodeHTTPServerEmitClose(loop, receiver); return tsU
	case "closeHandle": tsNodeHTTPServerCloseHandle(loop, receiver); return tsU
	case "assignSocketInternal": tsNodeHTTPServerAssignSocket(loop, receiver, arg(0)); return tsU
	case "connectionListenerPrelude":
		socket := arg(0)
		if tsTruthy(tsNodeHTTPServerGet(loop, socket, "_nodeHTTPServerSocket")) && !tsTruthy(tsNodeHTTPServerSym(loop, socket, "kHandedOff")) { return tsBooleanValue(true) }
		return tsBooleanValue(false)
	}
	return tsU
}
`

// nodeHTTPServerAPIRuntime maps Server, ServerResponse, and the standard HTTP
// server lifecycle onto the local net and crypto/tls runtime adapters.
const nodeHTTPServerAPIRuntime = `
var tsNodeHTTPServerAPI struct {
	sync.Once
	module, serverConstructor, responseConstructor tsValue
	serverPrototype, responsePrototype *tsObject
}
func tsNodeHTTPServerAPIServerProto() *tsObject { return tsNodeHTTPServerAPI.serverPrototype }
func tsNodeHTTPServerAPIResponseProto() *tsObject { return tsNodeHTTPServerAPI.responsePrototype }
func tsNodeHTTPServerAPIServerClass() tsValue { return tsNodeHTTPServerAPI.serverConstructor }
func tsNodeHTTPServerAPIResponseClass() tsValue { return tsNodeHTTPServerAPI.responseConstructor }
func tsNodeHTTPServerAPIKey(name string) tsValue { return tsStringReference(tsStringUTF8(name)) }
func tsNodeHTTPServerAPIGet(loop *tsLoop, value tsValue, name string) tsValue { return tsGet(loop, value, tsNodeHTTPServerAPIKey(name)) }
func tsNodeHTTPServerAPISet(loop *tsLoop, value tsValue, name string, next tsValue) { tsSet(loop, value, tsNodeHTTPServerAPIKey(name), next) }
func tsNodeHTTPServerAPIInit(loop *tsLoop) {
	tsNodeHTTPServerAPI.Do(func() {
		tsNodeHTTPServerInit(loop)
		tsNodeEventsInit(loop)
		tsNodeHTTPOutgoingInit(loop)
		tsNodeHTTPIncomingInit(loop)
		tsNodeHTTPServerAPI.serverPrototype = tsNewObject()
		tsNodeHTTPServerAPI.serverPrototype.prototype = tsObjectValue(tsNodeEventsSingleton.prototype)
		tsNodeHTTPServerAPI.responsePrototype = tsNewObject()
		tsNodeHTTPServerAPI.responsePrototype.prototype = tsObjectValue(tsNodeHTTPOutgoingSingleton.prototype)
	})
}
func tsNodeHTTPServerValidate(loop *tsLoop, name string, value tsValue, extra ...tsValue) {
	validator := tsNodeHTTPServerAPIGet(loop, tsInternalValidatorsModule(loop), name)
	if validator.kind == tsFunctionKind { tsCall(loop, validator, append([]tsValue{value}, extra...)...) }
}
func tsNodeHTTPServerStoreOptions(loop *tsLoop, server, options tsValue) {
	for _, entry := range []struct{name string; http1 float64}{{"requestTimeout",300000},{"keepAliveTimeout",5000},{"keepAliveTimeoutBuffer",1000},{"connectionsCheckingInterval",30000}} {
		value := tsNodeHTTPServerAPIGet(loop, options, entry.name)
		if !tsIsUndefined(value) { tsNodeHTTPServerValidate(loop,"validateInteger",value,tsNodeHTTPServerAPIKey(entry.name),tsNumberValue(0)) } else { value=tsNumberValue(entry.http1) }
		tsNodeHTTPServerAPISet(loop,server,entry.name,value)
	}
	requestTimeout:=tsNodeHTTPServerAPIGet(loop,server,"requestTimeout"); headersTimeout:=tsNodeHTTPServerAPIGet(loop,options,"headersTimeout")
	if tsIsUndefined(headersTimeout) { headersTimeout=tsNumberValue(math.Min(60000,tsNumber(requestTimeout))) } else { tsNodeHTTPServerValidate(loop,"validateInteger",headersTimeout,tsNodeHTTPServerAPIKey("headersTimeout"),tsNumberValue(0)) }
	if tsNumber(requestTimeout)>0&&tsNumber(headersTimeout)>tsNumber(requestTimeout){panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","The value of \"headersTimeout\" must be less than or equal to \"requestTimeout\".")})};tsNodeHTTPServerAPISet(loop,server,"headersTimeout",headersTimeout)
	for _, entry:=range []struct{name string; validator string; arg tsValue}{{"maxHeaderSize","validateInteger",tsNumberValue(0)},{"highWaterMark","validateInteger",tsNumberValue(0)},{"maxRequestsPerSocket","validateInteger",tsNumberValue(0)},{"timeout","validateInteger",tsNumberValue(0)},{"maxConnections","validateInteger",tsNumberValue(0)}} { value:=tsNodeHTTPServerAPIGet(loop,options,entry.name);if !tsIsUndefined(value){tsNodeHTTPServerValidate(loop,entry.validator,value,tsNodeHTTPServerAPIKey(entry.name),entry.arg);tsNodeHTTPServerAPISet(loop,server,entry.name,value)} }
	for _, name:=range []string{"requireHostHeader","insecureHTTPParser","httpAllowHalfOpen","joinDuplicateHeaders","rejectNonStandardBodyWrites","optimizeEmptyRequests"} { value:=tsNodeHTTPServerAPIGet(loop,options,name);if !tsIsUndefined(value){tsNodeHTTPServerValidate(loop,"validateBoolean",value,tsNodeHTTPServerAPIKey("options."+name));tsNodeHTTPServerAPISet(loop,server,name,value)} }
	if tsIsUndefined(tsNodeHTTPServerAPIGet(loop,server,"requireHostHeader")){tsNodeHTTPServerAPISet(loop,server,"requireHostHeader",tsBooleanValue(true))}
	if tsIsUndefined(tsNodeHTTPServerAPIGet(loop,server,"httpAllowHalfOpen")){tsNodeHTTPServerAPISet(loop,server,"httpAllowHalfOpen",tsBooleanValue(false))}
	if tsIsUndefined(tsNodeHTTPServerAPIGet(loop,server,"shouldUpgradeCallback")){tsNodeHTTPServerAPISet(loop,server,"shouldUpgradeCallback",tsNamedNativeMethod("defaultShouldUpgradeCallback",0,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsBooleanValue(tsNumber(tsNodeHTTPServerCall(loop,tsArg(args,0),"listenerCount",tsNodeHTTPServerAPIKey("upgrade")))>0)}))}
	for _, name:=range []string{"httpValidation","shouldUpgradeCallback","IncomingMessage","ServerResponse","uniqueHeaders"} { value:=tsNodeHTTPServerAPIGet(loop,options,name);if !tsIsUndefined(value){tsNodeHTTPServerAPISet(loop,server,name,value)} }
	if value:=tsNodeHTTPServerAPIGet(loop,options,"httpValidation");!tsIsUndefined(value){tsNodeHTTPServerValidate(loop,"validateOneOf",value,tsNodeHTTPServerAPIKey("options.httpValidation"),tsArrayValue(&tsArray{values:[]tsValue{tsNodeHTTPServerAPIKey("strict"),tsNodeHTTPServerAPIKey("relaxed"),tsNodeHTTPServerAPIKey("insecure")}}))}
	if value:=tsNodeHTTPServerAPIGet(loop,options,"uniqueHeaders");!tsIsUndefined(value){tsNodeHTTPOutgoingPut(loop,server,tsNodeHTTPOutgoingSingleton.kUniqueHeaders,tsNodeHTTPOutgoingParseUniqueHeaders(loop,value))}
	tsNodeHTTPServerAPISet(loop,server,"timeout",tsNumberValue(0));tsNodeHTTPServerAPISet(loop,server,"maxRequestsPerSocket",tsNumberValue(0))
}
func tsNodeHTTPServerCreate(loop *tsLoop,args []tsValue) tsValue {
	tsNodeHTTPServerAPIInit(loop)
	options,callback:=tsArg(args,0),tsArg(args,1)
	if options.kind==tsFunctionKind { callback=options;options=tsObjectValue(tsNewObject()) } else if tsNullish(options) { options=tsObjectValue(tsNewObject()) } else if options.kind!=tsObjectKind&&options.kind!=tsInstanceKind { panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_TYPE","The \"options\" argument must be of type object.")}) } else { copy:=tsNewObject();tsObjectSpread(copy,options);options=tsObjectValue(copy) }
	object:=tsNewObject();object.prototype=tsObjectValue(tsNodeHTTPServerAPIServerProto());server:=tsNodeEventCreate(loop,tsObjectValue(object),tsU);object.nativeClass=tsNativeClass("Server")
	for _,name:=range []string{"_unref","noDelay"}{tsNodeHTTPServerAPISet(loop,server,name,tsBooleanValue(name=="noDelay"))}
	for _,name:=range []string{"timeout","maxRequestsPerSocket"}{tsNodeHTTPServerAPISet(loop,server,name,tsNumberValue(0))}
	for _,name:=range []string{"kTrackedConnections","kPendingCloseGenerations"}{tsNodeHTTPServerPut(loop,server,name,tsNewCollection(true,tsU))}
	for _,name:=range []string{"kPendingDrainClose"}{tsNodeHTTPServerPut(loop,server,name,tsBooleanValue(false))}
	tsNodeHTTPServerPut(loop,server,"serverSymbol",tsNull);tsNodeHTTPServerPut(loop,server,"tlsSymbol",tsNull);tsNodeHTTPServerPut(loop,server,"optionsSymbol",options)
	tsNodeHTTPServerPut(loop,server,"kIncomingMessage",tsNodeHTTPServerAPIGet(loop,options,"IncomingMessage"));tsNodeHTTPServerPut(loop,server,"kServerResponse",tsNodeHTTPServerAPIGet(loop,options,"ServerResponse"))
	tsNodeHTTPServerStoreOptions(loop,server,options)
	object.descriptors=map[string]*tsDescriptor{}
	object.descriptors["httpAllowHalfOpen"]=&tsDescriptor{accessor:true,enumerable:true,configurable:true,getter:tsNamedNativeMethod("get httpAllowHalfOpen",0,func(loop *tsLoop,r tsValue,_ ...tsValue)tsValue{return tsNodeHTTPServerSym(loop,r,"kHttpAllowHalfOpen")}),setter:tsNamedNativeMethod("set httpAllowHalfOpen",1,func(loop *tsLoop,r tsValue,a ...tsValue)tsValue{value:=tsArg(a,0);tsNodeHTTPServerPut(loop,r,"kHttpAllowHalfOpen",value);return tsU})}
	object.descriptors["maxHeadersCount"]=&tsDescriptor{accessor:true,enumerable:true,configurable:true,getter:tsNamedNativeMethod("get maxHeadersCount",0,func(loop *tsLoop,r tsValue,_ ...tsValue)tsValue{return tsNodeHTTPServerSym(loop,r,"kMaxHeadersCount")}),setter:tsNamedNativeMethod("set maxHeadersCount",1,func(loop *tsLoop,r tsValue,a ...tsValue)tsValue{tsNodeHTTPServerPut(loop,r,"kMaxHeadersCount",tsArg(a,0));return tsU})}
	connection:=tsNodeHTTPServerMethodValue(loop,"connectionListener",1);tsNodeHTTPServerCall(loop,server,"on",tsNodeHTTPServerAPIKey("connection"),connection)
	prepend:=tsFunctionValue(tsFunc(func(loop *tsLoop,values ...tsValue)tsValue{socket:=tsArg(values,0);if !tsNullish(socket){tsNodeHTTPServerAPISet(loop,socket,"server",server)};return tsU}));tsNodeHTTPServerCall(loop,server,"prependListener",tsNodeHTTPServerAPIKey("connection"),prepend)
	listening:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{tsNodeHTTPServerMethod(loop,"serverSetupTracking",server,nil);return tsU}));tsNodeHTTPServerCall(loop,server,"on",tsNodeHTTPServerAPIKey("listening"),listening)
	if callback.kind==tsFunctionKind{tsNodeHTTPServerCall(loop,server,"on",tsNodeHTTPServerAPIKey("request"),callback)}
	if !tsIsUndefined(tsNodeHTTPServerAPIGet(loop,options,"pfx"))||!tsIsUndefined(tsNodeHTTPServerAPIGet(loop,options,"cert"))||!tsIsUndefined(tsNodeHTTPServerAPIGet(loop,options,"key")){tlsOptions:=tsNewObject();tsObjectSpread(tlsOptions,options);tsNodeHTTPServerDispatch(loop,"normalizeServerTls",tsObjectValue(tlsOptions),nil);tsNodeHTTPServerPut(loop,server,"tlsSymbol",tsObjectValue(tlsOptions))}
	return server
}
func tsNodeHTTPServerResponseCreate(loop *tsLoop,req,options tsValue) tsValue {
	tsNodeHTTPServerAPIInit(loop);if tsNullish(options){options=tsObjectValue(tsNewObject())}
	object:=tsNewObject();object.prototype=tsObjectValue(tsNodeHTTPServerAPIResponseProto());response:=tsNodeHTTPOutgoingCreate(loop,tsObjectValue(object),options)
	tsNodeHTTPServerAPISet(loop,response,"req",req);tsNodeHTTPServerAPISet(loop,response,"useChunkedEncodingByDefault",tsBooleanValue(true));tsNodeHTTPServerAPISet(loop,response,"sendDate",tsBooleanValue(true));tsNodeHTTPServerAPISet(loop,response,"_sent100",tsBooleanValue(false));tsNodeHTTPServerAPISet(loop,response,"statusCode",tsNumberValue(200));tsNodeHTTPServerAPISet(loop,response,"statusMessage",tsU)
	tsNodeHTTPServerAPISet(loop,response,"_removedConnection",tsBooleanValue(false));tsNodeHTTPServerAPISet(loop,response,"_removedContLen",tsBooleanValue(false));tsNodeHTTPServerAPISet(loop,response,"_hasBody",tsBooleanValue(true));tsNodeHTTPServerAPISet(loop,response,"_ended",tsBooleanValue(false));tsNodeHTTPServerAPISet(loop,response,"_closed",tsBooleanValue(false));tsNodeHTTPServerAPISet(loop,response,"finished",tsBooleanValue(false));tsNodeHTTPServerPut(loop,response,"kPendingFinish",tsU);tsNodeHTTPServerPut(loop,response,"kBytesBuffered",tsNumberValue(0));tsNodeHTTPServerPut(loop,response,"kAccountingFlushScheduled",tsBooleanValue(false));tsNodeHTTPServerPut(loop,response,"kPendingCallbacks",tsArrayValue(&tsArray{}));tsNodeHTTPOutgoingPut(loop,response,tsNodeHTTPOutgoingSingleton.kNeedDrain,tsBooleanValue(false))
	tsNodeHTTPServerSetHeaderState(loop,response,"none");tsNodeHTTPServerAPISet(loop,response,"shouldKeepAlive",tsBooleanValue(tsNodeHTTPServerResponseKeepAlive(loop,req)))
	server:=tsNodeHTTPServerGet(loop,tsNodeHTTPIncomingGet(loop,req,"socket"),"server");reject:=tsNodeHTTPServerAPIGet(loop,options,"rejectNonStandardBodyWrites");if tsIsUndefined(reject){reject=tsNodeHTTPServerAPIGet(loop,server,"rejectNonStandardBodyWrites")};if !tsIsUndefined(reject){tsNodeHTTPServerAPISet(loop,response,"rejectNonStandardBodyWrites",reject);tsNodeHTTPOutgoingPut(loop,response,tsNodeHTTPOutgoingSingleton.kRejectBody,reject)}
	water:=tsNodeHTTPServerAPIGet(loop,server,"highWaterMark");if tsIsNumeric(water){tsNodeHTTPOutgoingPut(loop,response,tsNodeHTTPOutgoingSingleton.kHighWaterMark,water)}
	if tsText(tsNodeHTTPIncomingGet(loop,req,"method"))=="HEAD"{tsNodeHTTPServerAPISet(loop,response,"_hasBody",tsBooleanValue(false))}
	if tsNumber(tsNodeHTTPIncomingGet(loop,req,"httpVersionMajor"))<1||tsNumber(tsNodeHTTPIncomingGet(loop,req,"httpVersionMinor"))<1{tsNodeHTTPServerAPISet(loop,response,"useChunkedEncodingByDefault",tsBooleanValue(false));tsNodeHTTPServerAPISet(loop,response,"shouldKeepAlive",tsBooleanValue(false))}
	tsNodeHTTPServerPut(loop,response,"kHandle",tsNodeHTTPServerAPIGet(loop,req,"_handle"));tsNodeHTTPOutgoingPut(loop,response,tsNodeHTTPOutgoingSingleton.kSocket,tsNull)
	return response
}
func tsNodeHTTPServerHeaderState(loop *tsLoop,response tsValue) tsValue {
	return tsGet(loop,response,tsNodeHTTPServerSingleton.symbols["headerStateSymbol"])
}
func tsNodeHTTPServerSetHeaderState(loop *tsLoop,response tsValue,name string) {
	state:=tsInternalHTTPGet(loop,tsInternalHTTPModule(loop),"NodeHTTPHeaderState")
	tsSet(loop,response,tsNodeHTTPServerSingleton.symbols["headerStateSymbol"],tsGet(loop,state,tsNodeHTTPServerAPIKey(name)))
}
func tsNodeHTTPServerHeaderStateAssignedOrSent(loop *tsLoop,response tsValue) bool {
	state:=tsNodeHTTPServerHeaderState(loop,response);states:=tsInternalHTTPGet(loop,tsInternalHTTPModule(loop),"NodeHTTPHeaderState")
	return tsStrictEqual(state,tsGet(loop,states,tsNodeHTTPServerAPIKey("assigned")))||tsStrictEqual(state,tsGet(loop,states,tsNodeHTTPServerAPIKey("sent")))
}
func tsNodeHTTPServerResponseKeepAlive(loop *tsLoop,request tsValue) bool {
	if tsNullish(request){return true}
	if tsNumber(tsNodeHTTPIncomingGet(loop,request,"httpVersionMajor"))==1&&tsNumber(tsNodeHTTPIncomingGet(loop,request,"httpVersionMinor"))==0{return false}
	key:=tsNodeHTTPIncomingGet(loop,tsNodeHTTPIncomingModule(loop),"kReqShouldKeepAlive");cached:=tsGet(loop,request,key)
	if !tsIsUndefined(cached){return tsTruthy(cached)}
	headers:=tsNodeHTTPIncomingGet(loop,request,"headers");connection:=tsNodeHTTPServerGet(loop,headers,"connection")
	if connection.kind==tsStringKind&&strings.Contains(strings.ToLower(tsText(connection)),"close"){return false}
	return true
}
func tsNodeHTTPServerResponseCloseOnFinish(loop *tsLoop,response,socket tsValue) {
	if tsTruthy(tsNodeHTTP1ServerGet(loop,response,"kMustCloseConnection")) { return }
	if tsNullish(socket)||tsTruthy(tsNodeHTTPServerGet(loop,socket,"destroyed")){return}
	current:=tsNodeHTTPServerGet(loop,socket,"_httpMessage");if !tsNullish(current)&&!tsStrictEqual(current,response){return};queue:=tsNodeHTTP1ServerSym(loop,socket,"kQueue");if queue.kind==tsArrayKind&&len((*tsArray)(queue.ref).values)>0{return}
	server:=tsNodeHTTPServerGet(loop,socket,"server");timeout:=tsNodeHTTPServerGet(loop,server,"keepAliveTimeout");buffer:=tsNodeHTTPServerGet(loop,server,"keepAliveTimeoutBuffer")
	if tsIsNumeric(timeout)&&tsNumber(timeout)>0{extra:=0.0;if tsIsNumeric(buffer){extra=tsNumber(buffer)};tsNodeHTTPServerCall(loop,socket,"setTimeout",tsNumberValue(tsNumber(timeout)+extra))}
}
func tsNodeHTTPServerResponseQueueFinish(loop *tsLoop,response,callback tsValue) {
	tsNodeHTTPServerMethod(loop,"responseCallPendingCallbacks",response,nil)
	cb:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{
		internal:=tsInternalHTTPModule(loop);closeFn:=tsInternalHTTPGet(loop,internal,"emitCloseNT");if closeFn.kind==tsFunctionKind{tsCall(loop,closeFn,response)}
		tsNodeHTTPServerCall(loop,response,"emit",tsNodeHTTPServerAPIKey("finish"))
		if callback.kind==tsFunctionKind{tsCall(loop,callback)}
		return tsU
	}))
	tsInternalHTTPNextTick(loop,cb)
}
func tsNodeHTTPServerQueuePipelinedOp(loop *tsLoop,response,queue tsValue,kind string,chunk,encoding,callback tsValue) bool {
	ops:=tsNodeHTTPServerGet(loop,queue,"ops");if ops.kind!=tsArrayKind{ops=tsArrayValue(&tsArray{});tsNodeHTTPServerSet(loop,queue,"ops",ops)}
	(*tsArray)(ops.ref).values=append((*tsArray)(ops.ref).values,tsArrayValue(&tsArray{values:[]tsValue{tsNodeHTTPServerAPIKey(kind),chunk,encoding,callback}}))
	bytes:=0;if !tsNullish(chunk){bytes=tsNodeHTTPOutgoingByteLength(loop,chunk,encoding)};total:=tsNumber(tsNodeHTTPServerGet(loop,queue,"bytes"))+float64(bytes);tsNodeHTTPServerSet(loop,queue,"bytes",tsNumberValue(total))
	socket:=tsNodeHTTPServerGet(loop,queue,"socket");if !tsNullish(socket){queuedBytes:=tsNumber(tsNodeHTTPServerGet(loop,socket,"_httpQueuedBytes"))+float64(bytes);tsNodeHTTPServerSet(loop,socket,"_httpQueuedBytes",tsNumberValue(queuedBytes))}
	mark:=tsNumber(tsNodeHTTPServerGet(loop,response,"writableHighWaterMark"));needDrain:=total>=mark;tsNodeHTTPServerSet(loop,queue,"needDrain",tsBooleanValue(needDrain));if needDrain&&!tsNullish(socket){tsNodeHTTP1ServerCallPipeline(loop,tsNodeHTTP1ServerPipeline(loop),"maybePauseHTTP1Reads",socket)};return !needDrain
}
func tsNodeHTTPServerResponseWrite(loop *tsLoop,response tsValue,args []tsValue) tsValue {
	chunk,encoding,callback:=tsArg(args,0),tsArg(args,1),tsArg(args,2)
	if chunk.kind==tsFunctionKind{callback=chunk;chunk=tsU;encoding=tsU}else if encoding.kind==tsFunctionKind{callback=encoding;encoding=tsU}
	if callback.kind!=tsFunctionKind{callback=tsU}
	finished:=tsInternalHTTPGet(loop,tsInternalHTTPModule(loop),"hasServerResponseFinished")
	if finished.kind==tsFunctionKind&&tsTruthy(tsCall(loop,finished,response,chunk,callback,tsBooleanValue(false))){return tsBooleanValue(false)}
	queued:=tsNodeHTTP1ServerSym(loop,response,"kQueueState");if !tsNullish(queued){if !tsTruthy(tsNodeHTTPServerGet(loop,response,"_header")){tsNodeHTTPServerCall(loop,response,"_implicitHeader")};if !tsTruthy(tsNodeHTTPServerGet(loop,response,"_hasBody")){if tsTruthy(tsNodeHTTPServerGet(loop,response,"rejectNonStandardBodyWrites"))&&tsTruthy(chunk){panic(tsThrown{tsFSException("TypeError","ERR_HTTP_BODY_NOT_ALLOWED","Adding content for this request method or response status is not allowed")})};if callback.kind==tsFunctionKind{tsInternalHTTPNextTick(loop,callback)};return tsBooleanValue(true)};return tsBooleanValue(tsNodeHTTPServerQueuePipelinedOp(loop,response,queued,"write",chunk,encoding,callback))}
	handle:=tsNodeHTTPServerSym(loop,response,"kHandle")
	if tsNullish(handle){return tsNodeHTTPOutgoingMethod(loop,response,"write",[]tsValue{chunk,encoding,callback})}
	if !tsTruthy(tsNodeHTTPServerGet(loop,response,"_header")){tsNodeHTTPServerCall(loop,response,"_implicitHeader")}
	if !tsTruthy(tsNodeHTTPServerGet(loop,response,"_hasBody")){
		if tsTruthy(tsNodeHTTPServerGet(loop,response,"rejectNonStandardBodyWrites")){panic(tsThrown{tsFSException("TypeError","ERR_HTTP_BODY_NOT_ALLOWED","Adding content for this request method or response status is not allowed")})}
		if callback.kind==tsFunctionKind{tsInternalHTTPNextTick(loop,callback)};return tsBooleanValue(true)
	}
	onWritable:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{tsNodeHTTPServerMethod(loop,"responseCallPendingCallbacks",response,nil);tsNodeHTTPServerCall(loop,response,"emit",tsNodeHTTPServerAPIKey("drain"));return tsU}))
	strict:=tsNodeHTTPServerDispatch(loop,"strictContentLength",response,[]tsValue{tsNodeHTTPServerHeaderState(loop,response),tsBooleanValue(false)})
	result:=tsNodeHTTPServerCall(loop,handle,"write",chunk,encoding,onWritable,strict)
	tsNodeHTTPServerSetHeaderState(loop,response,"sent")
	if tsIsNumeric(result)&&tsNumber(result)<0{if callback.kind==tsFunctionKind{pending:=tsNodeHTTPServerSym(loop,response,"kPendingCallbacks");if pending.kind==tsArrayKind{(*tsArray)(pending.ref).values=append((*tsArray)(pending.ref).values,callback)}};tsNodeHTTPServerSet(loop,handle,"onwritable",onWritable);return tsBooleanValue(false)}
	if callback.kind==tsFunctionKind{tsInternalHTTPNextTick(loop,callback)}
	bytes:=float64(tsNodeHTTPOutgoingByteLength(loop,chunk,encoding));old:=tsNodeHTTPServerSym(loop,response,"kBytesBuffered");if !tsIsNumeric(old){old=tsNumberValue(0)};total:=tsNumber(old)+bytes;tsNodeHTTPServerPut(loop,response,"kBytesBuffered",tsNumberValue(total))
	mark:=tsNumber(tsNodeHTTPServerGet(loop,response,"writableHighWaterMark"));if total>=mark{tsNodeHTTPOutgoingPut(loop,response,tsNodeHTTPOutgoingSingleton.kNeedDrain,tsBooleanValue(true))}
	if !tsTruthy(tsNodeHTTPServerSym(loop,response,"kAccountingFlushScheduled")){tsNodeHTTPServerPut(loop,response,"kAccountingFlushScheduled",tsBooleanValue(true));flush:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{return tsNodeHTTPServerMethod(loop,"responseFlushWriteAccounting",response,nil)}));tsInternalHTTPNextTick(loop,flush)}
	return tsBooleanValue(total<mark)
}
func tsNodeHTTPServerResponseEnd(loop *tsLoop,response tsValue,args []tsValue) tsValue {
	chunk,encoding,callback:=tsArg(args,0),tsArg(args,1),tsArg(args,2)
	if chunk.kind==tsFunctionKind{callback=chunk;chunk=tsU;encoding=tsU}else if encoding.kind==tsFunctionKind{callback=encoding;encoding=tsU}
	if callback.kind!=tsFunctionKind{callback=tsU}
	finished:=tsInternalHTTPGet(loop,tsInternalHTTPModule(loop),"hasServerResponseFinished")
	if finished.kind==tsFunctionKind&&tsTruthy(tsCall(loop,finished,response,chunk,callback,tsBooleanValue(true))){return response}
	queued:=tsNodeHTTP1ServerSym(loop,response,"kQueueState");if !tsNullish(queued){if !tsTruthy(tsNodeHTTPServerGet(loop,response,"_header")){tsNodeHTTPServerCall(loop,response,"_implicitHeader")};if !tsTruthy(tsNodeHTTPServerGet(loop,response,"_hasBody"))&&tsTruthy(chunk){if tsTruthy(tsNodeHTTPServerGet(loop,response,"rejectNonStandardBodyWrites")){panic(tsThrown{tsFSException("TypeError","ERR_HTTP_BODY_NOT_ALLOWED","Adding content for this request method or response status is not allowed")})};chunk=tsU};tsNodeHTTPServerQueuePipelinedOp(loop,response,queued,"end",chunk,encoding,callback);tsNodeHTTPServerSet(loop,queued,"ended",tsBooleanValue(true));tsNodeHTTPServerAPISet(loop,response,"finished",tsBooleanValue(true));return response}
	if tsNullish(tsNodeHTTPServerSym(loop,response,"kHandle")){return tsNodeHTTPOutgoingEnd(loop,response,[]tsValue{chunk,encoding,callback})}
	if !tsTruthy(tsNodeHTTPServerGet(loop,response,"_header")){if !tsNullish(chunk){tsNodeHTTPServerAPISet(loop,response,"_contentLength",tsNumberValue(float64(tsNodeHTTPOutgoingByteLength(loop,chunk,encoding))))};tsNodeHTTPServerCall(loop,response,"_implicitHeader")}
	if (tsNumber(tsNodeHTTPServerGet(loop,response,"statusCode"))==204||tsNumber(tsNodeHTTPServerGet(loop,response,"statusCode"))==304||tsNumber(tsNodeHTTPServerGet(loop,response,"statusCode"))>=100&&tsNumber(tsNodeHTTPServerGet(loop,response,"statusCode"))<200)&&tsTruthy(chunk){if tsTruthy(tsNodeHTTPServerGet(loop,response,"rejectNonStandardBodyWrites")){panic(tsThrown{tsFSException("TypeError","ERR_HTTP_BODY_NOT_ALLOWED","Adding content for this request method or response status is not allowed")})};chunk=tsU}
	if tsTruthy(tsNodeHTTPServerGet(loop,response,"_trailer"))&&tsTruthy(tsNodeHTTPServerGet(loop,response,"_hasBody")){tsNodeHTTPServerCall(loop,tsNodeHTTPOutgoingSymbol(loop,response,tsNodeHTTPOutgoingSingleton.kSocket),"setResponseTrailers",tsNodeHTTPServerGet(loop,response,"_trailer"))}
	handle:=tsNodeHTTPServerSym(loop,response,"kHandle");strict:=tsNodeHTTPServerDispatch(loop,"strictContentLength",response,[]tsValue{tsNodeHTTPServerHeaderState(loop,response),tsBooleanValue(true)})
	result:=tsNodeHTTPServerCall(loop,handle,"end",chunk,encoding,tsU,strict)
	tsNodeHTTPServerAPISet(loop,response,"_header",tsNodeHTTPServerAPIKey(" "));tsNodeHTTPServerAPISet(loop,response,"finished",tsBooleanValue(true));tsNodeHTTPServerSet(loop,response,"_ended",tsBooleanValue(true));tsNodeHTTPServerCall(loop,response,"emit",tsNodeHTTPServerAPIKey("prefinish"))
	tsNodeHTTPServerSetHeaderState(loop,response,"sent")
	if tsIsNumeric(result)&&tsNumber(result)<0{
		pending:=callback;if pending.kind!=tsFunctionKind{pending=tsNull};tsNodeHTTPServerPut(loop,response,"kPendingFinish",pending)
		writable:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{cb:=tsNodeHTTPServerSym(loop,response,"kPendingFinish");if tsIsUndefined(cb){return tsU};tsNodeHTTPServerPut(loop,response,"kPendingFinish",tsU);tsNodeHTTPServerResponseQueueFinish(loop,response,cb);return tsU}))
		tsNodeHTTPServerSet(loop,handle,"onwritable",writable)
	}else{tsNodeHTTPServerResponseQueueFinish(loop,response,callback)}
	return response
}
func tsNodeHTTPServerWriteHead(loop *tsLoop,response tsValue,args []tsValue) tsValue {
	if tsNodeHTTPServerHeaderStateAssignedOrSent(loop,response){panic(tsThrown{tsFSException("Error","ERR_HTTP_HEADERS_SENT","Cannot write headers after they are sent to the client")})}
	status:=tsArg(args,0);tsNodeHTTPServerValidate(loop,"validateInteger",status,tsNodeHTTPServerAPIKey("statusCode"),tsNumberValue(100),tsNumberValue(999));code:=int(tsNumber(status));if code<100||code>999{panic(tsThrown{tsFSException("RangeError","ERR_HTTP_INVALID_STATUS_CODE","Invalid status code: "+tsText(status))})}
	reason,headers:=tsArg(args,1),tsArg(args,2);if reason.kind==tsStringKind{tsNodeHTTPServerAPISet(loop,response,"statusMessage",reason)}else{headers=reason;codes:=tsInternalHTTPGet(loop,tsInternalHTTPModule(loop),"STATUS_CODES");message:=tsGet(loop,codes,tsNodeHTTPServerAPIKey(strconv.Itoa(code)));if !tsTruthy(tsNodeHTTPServerAPIGet(loop,response,"statusMessage")){if tsNullish(message){message=tsNodeHTTPServerAPIKey("unknown")};tsNodeHTTPServerAPISet(loop,response,"statusMessage",message)}}
	check:=tsNodeHTTPCommonGet(loop,tsNodeHTTPCommonModule(loop),"_checkInvalidHeaderChar");if tsTruthy(tsCall(loop,check,tsNodeHTTPServerAPIGet(loop,response,"statusMessage"))){panic(tsThrown{tsFSException("TypeError","ERR_INVALID_CHAR","Invalid character in statusMessage")})}
	tsNodeHTTPServerAPISet(loop,response,"statusCode",tsNumberValue(float64(code)));if code==204||code==304||code>=100&&code<200{tsNodeHTTPServerAPISet(loop,response,"_hasBody",tsBooleanValue(false))}
	if headers.kind==tsArrayKind{items:=(*tsArray)(headers.ref).values;for i:=0;i+1<len(items);i+=2{tsNodeHTTPOutgoingCall(loop,response,"setHeader",items[i],items[i+1])}}else if headers.kind==tsObjectKind||headers.kind==tsInstanceKind{for _,key:=range tsObjectKeys(headers).values{tsNodeHTTPOutgoingCall(loop,response,"setHeader",key,tsGet(loop,headers,key))}}
	connectionHeaderExplicit := tsTruthy(tsNodeHTTPOutgoingCall(loop, response, "hasHeader", tsNodeHTTPServerAPIKey("connection")))
	keepAliveHeaderExplicit := tsTruthy(tsNodeHTTPOutgoingCall(loop, response, "hasHeader", tsNodeHTTPServerAPIKey("keep-alive")))
	if !connectionHeaderExplicit {
		if tsTruthy(tsNodeHTTPServerAPIGet(loop, response, "shouldKeepAlive")) {
			tsNodeHTTPOutgoingCall(loop, response, "setHeader", tsNodeHTTPServerAPIKey("Connection"), tsNodeHTTPServerAPIKey("keep-alive"))
			timeout := tsNodeHTTPServerAPIGet(loop, response, "_keepAliveTimeout")
			if tsIsNumeric(timeout) && tsNumber(timeout) > 0 {
				socket := tsNodeHTTPIncomingGet(loop, tsNodeHTTPServerAPIGet(loop, response, "req"), "socket")
				httpServer := tsNodeHTTPServerGet(loop, socket, "server")
				buffer := tsNodeHTTPServerAPIGet(loop, httpServer, "keepAliveTimeoutBuffer")
				total := tsNumber(timeout) + tsNumber(buffer)
				keepAlive := "timeout=" + strconv.Itoa(int(math.Floor(total/1000)))
				maxRequests := tsNodeHTTPServerGet(loop, response, "_maxRequestsPerSocket")
				if !keepAliveHeaderExplicit && tsIsNumeric(maxRequests) && tsNumber(maxRequests) > 0 { keepAlive += ", max=" + strconv.Itoa(int(tsNumber(maxRequests))) }
				if !keepAliveHeaderExplicit { tsNodeHTTPOutgoingCall(loop, response, "setHeader", tsNodeHTTPServerAPIKey("Keep-Alive"), tsNodeHTTPServerAPIKey(keepAlive)) }
			}
		} else {
			tsNodeHTTPOutgoingCall(loop, response, "setHeader", tsNodeHTTPServerAPIKey("Connection"), tsNodeHTTPServerAPIKey("close"))
			pipeline := tsNodeHTTP1ServerPipeline(loop)
			closeKey := tsNodeHTTP1ServerGet(loop, pipeline, "kMustCloseConnection")
			tsSet(loop, response, closeKey, tsBooleanValue(true))
		}
	}
	transferEncoding := tsNodeHTTPOutgoingHeaderGet(loop, response, tsNodeHTTPServerAPIKey("transfer-encoding"))
	chunkedTransferEncoding, _, _ := tsNodeHTTPCommonTransferChunked(tsNodeHTTPOutgoingText(loop, transferEncoding))
	forceClose := tsTruthy(tsNodeHTTPServerGet(loop, response, "maxRequestsOnConnectionReached")) && !connectionHeaderExplicit || (code == 204 || code == 304) && chunkedTransferEncoding
	if forceClose {
		if !connectionHeaderExplicit {
			tsNodeHTTPOutgoingCall(loop, response, "removeHeader", tsNodeHTTPServerAPIKey("connection"))
			if !keepAliveHeaderExplicit { tsNodeHTTPOutgoingCall(loop, response, "removeHeader", tsNodeHTTPServerAPIKey("keep-alive")) }
			tsNodeHTTPOutgoingCall(loop, response, "setHeader", tsNodeHTTPServerAPIKey("Connection"), tsNodeHTTPServerAPIKey("close"))
		}
		pipeline := tsNodeHTTP1ServerPipeline(loop)
		closeKey := tsNodeHTTP1ServerGet(loop, pipeline, "kMustCloseConnection")
		tsSet(loop, response, closeKey, tsBooleanValue(true))
	}
	statusLine:="HTTP/1.1 "+strconv.Itoa(code)+" "+tsText(tsNodeHTTPServerAPIGet(loop,response,"statusMessage"))+"\r\n";tsNodeHTTPOutgoingSet(loop,response,"_header",tsNodeHTTPServerAPIKey(statusLine));tsNodeHTTPOutgoingCall(loop,response,"_storeHeader",tsNodeHTTPServerAPIKey(statusLine),tsNodeHTTPOutgoingSymbol(loop,response,tsNodeHTTPOutgoingSingleton.kOutHeaders))
	connectionHeader := tsNodeHTTPOutgoingHeaderGet(loop, response, tsNodeHTTPServerAPIKey("connection"))
	closeConnection := forceClose || !tsTruthy(tsNodeHTTPServerAPIGet(loop, response, "shouldKeepAlive")) || (connectionHeader.kind == tsStringKind && tsNodeHTTPCommonHasBoundaryToken(tsText(connectionHeader), "close"))
	if closeConnection{pipeline:=tsNodeHTTP1ServerPipeline(loop);closeKey:=tsNodeHTTP1ServerGet(loop,pipeline,"kMustCloseConnection");tsSet(loop,response,closeKey,tsBooleanValue(true))}
	handle := tsNodeHTTPServerSym(loop, response, "kHandle")
	if !tsNullish(handle) {
		autoBits := 0
		if tsTruthy(tsNodeHTTPServerAPIGet(loop, response, "sendDate")) { autoBits |= 1 }
		if !tsTruthy(tsNodeHTTPServerAPIGet(loop, response, "_hasBody")) { autoBits |= 32 }
		if closeConnection { autoBits |= 4 } else { autoBits |= 2 }
		socket := tsNodeHTTPIncomingGet(loop, tsNodeHTTPServerAPIGet(loop, response, "req"), "socket")
		server := tsNodeHTTPServerGet(loop, socket, "server")
		timeout := tsNodeHTTPServerAPIGet(loop, server, "keepAliveTimeout")
		keepSecs := tsNumberValue(0)
		if !closeConnection && tsIsNumeric(timeout) && tsNumber(timeout) > 0 {
			autoBits |= 8
			buffer := tsNodeHTTPServerAPIGet(loop, server, "keepAliveTimeoutBuffer")
			total := tsNumber(timeout) + tsNumber(buffer)
			keepSecs = tsNumberValue(math.Floor(total/1000))
		}
		tsNodeHTTPServerCall(loop, handle, "writeHead", tsNumberValue(float64(code)), tsNodeHTTPServerAPIGet(loop, response, "statusMessage"), tsNodeHTTPOutgoingHeaders(loop, response), tsNumberValue(float64(autoBits)), keepSecs)
	}
	tsNodeHTTPServerSetHeaderState(loop,response,"assigned")
	return response
}
func tsNodeHTTPServerMethod(loop *tsLoop,name string,receiver tsValue,args []tsValue) tsValue {
	arg:=func(i int)tsValue{return tsArg(args,i)}
	switch name {
	case "serverRef","serverUnref": unref:=name=="serverUnref";tsNodeHTTPServerAPISet(loop,receiver,"_unref",tsBooleanValue(unref));server:=tsNodeHTTPServerSym(loop,receiver,"serverSymbol");if !tsNullish(server){method:="ref";if unref{method="unref"};tsNodeHTTPServerCall(loop,server,method)};return receiver
	case "serverAddress":server:=tsNodeHTTPServerSym(loop,receiver,"serverSymbol");if tsNullish(server){return tsNull};return tsNodeHTTPServerCall(loop,server,"address")
	case "serverGetConnections":tracked:=tsNodeHTTPServerSym(loop,receiver,"kTrackedConnections");count:=0;if tracked.kind==tsCollectionKind{count=len((*tsCollection)(tracked.ref).index)};cb:=arg(0);if cb.kind==tsFunctionKind{tsInternalHTTPNextTick(loop,cb,tsNull,tsNumberValue(float64(count)))};return receiver
	case "serverCloseAll","serverCloseIdle":
		if name=="serverCloseIdle"{tsNodeHTTP1ServerCloseIdle(loop,receiver)}else{tsNodeHTTP1ServerCloseAll(loop,receiver)}
		tracked:=tsNodeHTTPServerSym(loop,receiver,"kTrackedConnections");if tracked.kind==tsCollectionKind{for _,entry:=range (*tsCollection)(tracked.ref).entries{if !entry.live||tsTruthy(tsNodeHTTPServerSym(loop,entry.key,"kHandedOff")){continue};if name=="serverCloseIdle"&&tsTruthy(tsNodeHTTPServerGet(loop,entry.key,"_httpMessage")){continue};tsNodeHTTPServerCall(loop,entry.key,"destroy")}};return tsU
	case "serverClose":
		timers:=tsNodeTimersModule();interval:=tsNodeHTTPServerSym(loop,receiver,"kConnectionsCheckingInterval");if !tsNullish(interval){tsNodeHTTPServerAPIGet(loop,timers,"clearInterval");tsCall(loop,tsNodeHTTPServerAPIGet(loop,timers,"clearInterval"),interval)}
		server:=tsNodeHTTPServerSym(loop,receiver,"serverSymbol");if tsNullish(server){cb:=arg(0);if cb.kind==tsFunctionKind{err:=tsFSException("Error","ERR_SERVER_NOT_RUNNING","Server is not running.");tsInternalHTTPNextTick(loop,cb,err)};return receiver}
		cb:=arg(0);if cb.kind==tsFunctionKind{tsNodeHTTPServerCall(loop,receiver,"once",tsNodeHTTPServerAPIKey("close"),cb)};tsNodeHTTPServerPut(loop,receiver,"serverSymbol",tsU);tsNodeHTTPServerPut(loop,receiver,"kPendingDrainClose",tsBooleanValue(true));tsNodeHTTPServerCall(loop,server,"close");return receiver
	case "serverSetupTracking":
		timers:=tsNodeTimersModule();clear:=tsNodeHTTPServerAPIGet(loop,timers,"clearInterval");old:=tsNodeHTTPServerSym(loop,receiver,"kConnectionsCheckingInterval");if !tsNullish(old){tsCall(loop,clear,old)};delay:=tsNodeHTTPServerAPIGet(loop,receiver,"connectionsCheckingInterval");if !tsIsNumeric(delay){delay=tsNumberValue(30000)};callback:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{return tsNodeHTTPServerMethod(loop,"serverCheckConnections",receiver,nil)}));interval:=tsCall(loop,tsNodeHTTPServerAPIGet(loop,timers,"setInterval"),callback,delay);tsNodeHTTPServerPut(loop,receiver,"kConnectionsCheckingInterval",interval);tsNodeHTTPServerCall(loop,interval,"unref");return tsU
	case "serverCheckConnections":
		if tsNumber(tsNodeHTTPServerAPIGet(loop,receiver,"headersTimeout"))==0&&tsNumber(tsNodeHTTPServerAPIGet(loop,receiver,"requestTimeout"))==0{return tsU};tracked:=tsNodeHTTPServerSym(loop,receiver,"kTrackedConnections");if tracked.kind==tsCollectionKind{for _,entry:=range (*tsCollection)(tracked.ref).entries{if !entry.live{continue};socket:=entry.key;handle:=tsNodeHTTPServerSym(loop,socket,"kHandle");timeout:=tsNodeHTTPServerAPIGet(loop,handle,"isRequestTimedOut");if timeout.kind==tsFunctionKind&&tsTruthy(tsCall(loop,timeout,tsNodeHTTPServerAPIGet(loop,receiver,"headersTimeout"),tsNodeHTTPServerAPIGet(loop,receiver,"requestTimeout"))){tsNodeHTTPServerMethod(loop,"serverRequestTimeout",socket,nil)}}};return tsU
	case "serverRequestTimeout":
		err:=tsFSException("Error","ERR_HTTP_REQUEST_TIMEOUT","Request timeout");tsNodeHTTPServerEmit(loop,receiver,"error",err);return tsU
	case "serverCaptureRejection":
		err,event,argsValue:=arg(0),arg(1),arg(2);if tsText(event)=="request"&&argsValue.kind==tsArrayKind&&len((*tsArray)(argsValue.ref).values)>1{response:=(*tsArray)(argsValue.ref).values[1];if !tsTruthy(tsNodeHTTPServerAPIGet(loop,response,"headersSent"))&&!tsTruthy(tsNodeHTTPServerAPIGet(loop,response,"writableEnded")){names:=tsNodeHTTPOutgoingCall(loop,response,"getHeaderNames");if names.kind==tsArrayKind{for _,header:=range (*tsArray)(names.ref).values{tsNodeHTTPOutgoingCall(loop,response,"removeHeader",header)}};tsNodeHTTPServerAPISet(loop,response,"statusCode",tsNumberValue(500));tsNodeHTTPOutgoingCall(loop,response,"end",tsNodeHTTPServerAPIKey("Internal Server Error"))}else{tsNodeHTTPServerCall(loop,response,"destroy")}}else{_ = err};return tsU
	case "serverAsyncDispose":
		promise:=loop.promise();callback:=tsFunctionValue(tsFunc(func(loop *tsLoop,values ...tsValue)tsValue{err:=tsArg(values,0);if tsTruthy(err){promise.settle(tsResult{err,true})}else{promise.settle(tsResult{tsU,false})};return tsU}));tsNodeHTTPServerCall(loop,receiver,"close",callback);return tsPromiseValue(promise)
	case "serverSetTimeout":ms:=arg(0);if tsIsUndefined(ms){ms=tsNumberValue(0)};tsNodeHTTPServerAPISet(loop,receiver,"timeout",ms);if arg(1).kind==tsFunctionKind{tsNodeHTTPServerCall(loop,receiver,"on",tsNodeHTTPServerAPIKey("timeout"),arg(1))};return receiver
	case "serverListen":return tsNodeHTTPServerListen(loop,receiver,args)
	case "connectionListener":
		socket:=arg(0);if timeout:=tsNodeHTTPServerAPIGet(loop,receiver,"timeout");tsTruthy(timeout){tsNodeHTTPServerCall(loop,socket,"setTimeout",timeout)};tsNodeHTTPServerAPISet(loop,socket,"allowHalfOpen",tsNodeHTTPServerSym(loop,receiver,"kHttpAllowHalfOpen"));tsNodeHTTPServerCall(loop,socket,"on",tsNodeHTTPServerAPIKey("timeout"),tsNodeHTTPServerFn(loop,"onNodeHTTPServerSocketTimeout"));if tsTruthy(tsNodeHTTPServerAPIGet(loop,socket,"_nodeHTTPServerSocket"))&&!tsTruthy(tsNodeHTTPServerSym(loop,socket,"kHandedOff")){return tsU};tracked:=tsNodeHTTPServerSym(loop,receiver,"kTrackedConnections");if tracked.kind==tsCollectionKind{add:=tsGet(loop,tracked,tsNodeHTTPServerAPIKey("add"));if add.kind==tsFunctionKind{tsCallReceiver(loop,add,tracked,socket)}};maxConnections:=tsNodeHTTPServerAPIGet(loop,receiver,"maxConnections");count:=0;if tracked.kind==tsCollectionKind{count=len((*tsCollection)(tracked.ref).index)};if tsIsNumeric(maxConnections)&&tsNumber(maxConnections)>0&&float64(count)>tsNumber(maxConnections){if tracked.kind==tsCollectionKind{del:=tsGet(loop,tracked,tsNodeHTTPServerAPIKey("delete"));if del.kind==tsFunctionKind{tsCallReceiver(loop,del,tracked,socket)}};data:=tsNewObject();for _,name:=range []string{"localAddress","localPort","localFamily","remoteAddress","remotePort","remoteFamily"}{tsSet(loop,data,tsNodeHTTPServerAPIKey(name),tsNodeHTTPServerAPIGet(loop,socket,name))};tsNodeHTTPServerCall(loop,socket,"destroy");tsNodeHTTPServerEmit(loop,receiver,"drop",tsObjectValue(data));return tsU};blockList:=tsNodeHTTPServerAPIGet(loop,receiver,"blockList");remote:=tsNodeHTTPServerAPIGet(loop,socket,"remoteAddress");ipType:=tsNumber(tsCall(loop,tsNodeHTTPServerAPIGet(loop,tsInternalNetIPModule(),"isIP"),remote));if !tsNullish(blockList)&&ipType!=0{check:=tsNodeHTTPServerAPIGet(loop,blockList,"check");family:="ipv4";if ipType==6{family="ipv6"};if check.kind==tsFunctionKind&&tsTruthy(tsCall(loop,check,remote,tsNodeHTTPServerAPIKey(family))){if tracked.kind==tsCollectionKind{del:=tsGet(loop,tracked,tsNodeHTTPServerAPIKey("delete"));if del.kind==tsFunctionKind{tsCallReceiver(loop,del,tracked,socket)}};tsNodeHTTPServerCall(loop,socket,"destroy");return tsU}};closed:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{current:=tsNodeHTTPServerSym(loop,receiver,"kTrackedConnections");del:=tsGet(loop,current,tsNodeHTTPServerAPIKey("delete"));if del.kind==tsFunctionKind{tsCallReceiver(loop,del,current,socket)};return tsU}));tsNodeHTTPServerCall(loop,socket,"once",tsNodeHTTPServerAPIKey("close"),closed);tsNodeHTTPServerAPISet(loop,socket,"server",receiver);tsNodeHTTPServerAPISet(loop,socket,"parser",tsNodeHTTPServerDispatch(loop,"createServerParserShim",socket,nil));options:=tsNewObject();http1Options:=tsNewObject();incoming:=tsNodeHTTPServerSym(loop,receiver,"kIncomingMessage");response:=tsNodeHTTPServerSym(loop,receiver,"kServerResponse");if tsNullish(incoming){incoming=tsNodeHTTPIncomingModule(loop);incoming=tsNodeHTTPIncomingGet(loop,incoming,"IncomingMessage")};if tsNullish(response){response=tsNodeHTTPServerAPIGet(loop,tsNodeHTTPServerAPIModule(loop),"ServerResponse")};http1Options.set("IncomingMessage",incoming);http1Options.set("ServerResponse",response);options.set("http1Options",tsObjectValue(http1Options));http1:=tsNodeHTTP1ServerModule(loop);fn:=tsNodeHTTP1ServerGet(loop,http1,"connectionListenerHTTP1");if fn.kind==tsFunctionKind{tsCall(loop,fn,receiver,socket,tsObjectValue(options))};return tsU
	case "responseAssignSocket":
		socket:=arg(0);tsNodeHTTPServerAssignSocket(loop,receiver,socket);tsNodeHTTPOutgoingPut(loop,receiver,tsNodeHTTPOutgoingSingleton.kSocket,socket);tsNodeHTTPServerCall(loop,socket,"once",tsNodeHTTPServerAPIKey("close"),tsNodeHTTPServerFn(loop,"onServerResponseClose"));tsNodeHTTPOutgoingCall(loop,receiver,"_flush");return receiver
	case "responseDetachSocket":
		socket:=arg(0);if tsNullish(socket){socket=tsNodeHTTPOutgoingSymbol(loop,receiver,tsNodeHTTPOutgoingSingleton.kSocket)};if !tsNullish(socket){tsNodeHTTPServerCall(loop,socket,"removeListener",tsNodeHTTPServerAPIKey("close"),tsNodeHTTPServerFn(loop,"onServerResponseClose"));tsNodeHTTPServerSet(loop,socket,"_httpMessage",tsNull);tsNodeHTTPServerSetCloseCallback(loop,socket,tsU);request:=tsNodeHTTPServerAPIGet(loop,receiver,"req");requestKey:=tsNodeHTTPServerSingleton.symbols["kRequest"];if tsStrictEqual(tsGet(loop,socket,requestKey),request){tsSet(loop,socket,requestKey,tsU)};tsNodeHTTPServerEmit(loop,receiver,"socket",tsNull)};tsNodeHTTPServerAPISet(loop,receiver,"socket",tsNull);tsNodeHTTPOutgoingPut(loop,receiver,tsNodeHTTPOutgoingSingleton.kSocket,tsNull);return tsU
	case "responseWriteHead":return tsNodeHTTPServerWriteHead(loop,receiver,args)
	case "responseImplicitHeader":if tsNodeHTTPServerHeaderStateAssignedOrSent(loop,receiver){return receiver};return tsNodeHTTPServerWriteHead(loop,receiver,[]tsValue{tsNodeHTTPServerAPIGet(loop,receiver,"statusCode")})
	case "responseWrite":return tsNodeHTTPServerResponseWrite(loop,receiver,args)
	case "responseEnd":return tsNodeHTTPServerResponseEnd(loop,receiver,args)
	case "responseSetHeader","responseAppendHeader","responseSetHeaders","responseRemoveHeader":
		if tsNodeHTTPServerHeaderStateAssignedOrSent(loop,receiver){operation:="set";if name=="responseAppendHeader"{operation="append"};panic(tsThrown{tsFSException("Error","ERR_HTTP_HEADERS_SENT","Cannot "+operation+" headers after they are sent to the client")})}
		switch name{case "responseSetHeader":return tsNodeHTTPOutgoingMethod(loop,receiver,"setHeader",args);case "responseAppendHeader":return tsNodeHTTPOutgoingMethod(loop,receiver,"appendHeader",args);case "responseSetHeaders":return tsNodeHTTPOutgoingMethod(loop,receiver,"setHeaders",args);default:return tsNodeHTTPOutgoingMethod(loop,receiver,"removeHeader",args)}
	case "responseWriteRaw":
		chunk,encoding,callback:=arg(0),arg(1),arg(2);handle:=tsNodeHTTPServerSym(loop,receiver,"kHandle")
		if tsNullish(handle)||!tsNullish(tsNodeHTTP1ServerSym(loop,receiver,"kQueueState")){return tsNodeHTTPOutgoingWriteRaw(loop,receiver,chunk,encoding,callback)}
		tsNodeHTTPServerCall(loop,handle,"writeInformational",chunk,encoding);if callback.kind==tsFunctionKind{tsInternalHTTPNextTick(loop,callback)};return tsBooleanValue(true)
	case "responseSend":
		data,encoding,callback:=arg(0),arg(1),arg(2);handle:=tsNodeHTTPServerSym(loop,receiver,"kHandle")
		if tsNullish(handle)||!tsNullish(tsNodeHTTP1ServerSym(loop,receiver,"kQueueState")){return tsNodeHTTPOutgoingMethod(loop,receiver,"_send",args)}
		if !tsTruthy(tsNodeHTTPServerGet(loop,receiver,"_header")){tsNodeHTTPServerCall(loop,receiver,"_implicitHeader")};if !tsStrictEqual(tsNodeHTTPServerHeaderState(loop,receiver),tsGet(loop,tsInternalHTTPGet(loop,tsInternalHTTPModule(loop),"NodeHTTPHeaderState"),tsNodeHTTPServerAPIKey("sent"))){tsNodeHTTPServerSetHeaderState(loop,receiver,"sent")}
		strict:=tsNodeHTTPServerDispatch(loop,"strictContentLength",receiver,[]tsValue{tsNodeHTTPServerHeaderState(loop,receiver),tsBooleanValue(false)});return tsNodeHTTPServerCall(loop,handle,"write",data,encoding,callback,strict)
	case "responseWriteContinue":
		if tsNodeHTTPServerHeaderStateAssignedOrSent(loop,receiver){panic(tsThrown{tsFSException("Error","ERR_HTTP_HEADERS_SENT","Cannot write headers after they are sent to the client")})}
		handle:=tsNodeHTTPServerSym(loop,receiver,"kHandle");if tsNullish(handle)||!tsNullish(tsNodeHTTP1ServerSym(loop,receiver,"kQueueState")){tsNodeHTTPServerMethod(loop,"responseWriteInformation",receiver,[]tsValue{tsNumberValue(100),tsNull,arg(0)})}else{tsNodeHTTPServerCall(loop,handle,"writeContinue");if arg(0).kind==tsFunctionKind{tsInternalHTTPNextTick(loop,arg(0))}}
		tsNodeHTTPServerAPISet(loop,receiver,"_sent100",tsBooleanValue(true));return tsU
	case "responseWriteProcessing":return tsNodeHTTPServerMethod(loop,"responseWriteInformation",receiver,[]tsValue{tsNumberValue(102),tsNull,arg(0)})
	case "responseWriteInformation":
		if tsNodeHTTPServerHeaderStateAssignedOrSent(loop,receiver){panic(tsThrown{tsFSException("Error","ERR_HTTP_HEADERS_SENT","Cannot write headers after they are sent to the client")})}
		status:=int(tsNumber(arg(0)));if status<100||status>199||status==101{panic(tsThrown{tsFSException("RangeError","ERR_HTTP_INVALID_STATUS_CODE","Invalid informational status code: "+strconv.Itoa(status))})}
		if !tsIsUndefined(arg(0))&&tsNumber(arg(0))!=float64(status){panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","The value of \"statusCode\" is out of range")})}
		reason:=tsGet(loop,tsInternalHTTPGet(loop,tsInternalHTTPModule(loop),"STATUS_CODES"),tsNodeHTTPServerAPIKey(strconv.Itoa(status)));if !tsTruthy(reason){reason=tsNodeHTTPServerAPIKey("unknown")};output:="HTTP/1.1 "+strconv.Itoa(status)+" "+tsText(reason)+"\r\n"
		appendHeader:=func(key,value tsValue){tsNodeHTTPCommonValidateHeaderName(loop,key,tsNodeHTTPServerAPIKey("name"));validate:=tsNodeHTTPCommonGet(loop,tsNodeHTTPCommonModule(loop),"validateHeaderValue");if validate.kind==tsFunctionKind{tsCall(loop,validate,key,value)};output+=tsText(key)+": "+tsNodeHTTPOutgoingText(loop,value)+"\r\n"}
		headers:=arg(1);if headers.kind==tsArrayKind{items:=(*tsArray)(headers.ref).values;if len(items)>0&&items[0].kind==tsArrayKind{for _,entry:=range items{appendHeader(tsArg((*tsArray)(entry.ref).values,0),tsArg((*tsArray)(entry.ref).values,1))}}else{if len(items)%2!=0{panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_VALUE","The \"headers\" argument must be an even-length array")})};for i:=0;i+1<len(items);i+=2{appendHeader(items[i],items[i+1])}}}else if !tsNullish(headers){if headers.kind!=tsObjectKind&&headers.kind!=tsInstanceKind{panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_TYPE","The \"headers\" argument must be of type object")})};for _,key:=range tsObjectKeys(headers).values{appendHeader(key,tsGet(loop,headers,key))}}
		output+="\r\n";return tsNodeHTTPServerMethod(loop,"responseWriteRaw",receiver,[]tsValue{tsNodeHTTPServerAPIKey(output),tsNodeHTTPServerAPIKey("ascii"),arg(2)})
	case "responseDestroy":
		if tsTruthy(tsNodeHTTPServerAPIGet(loop,receiver,"destroyed")){return receiver};handle:=tsNodeHTTPServerSym(loop,receiver,"kHandle");if !tsNullish(handle)&&tsNullish(tsNodeHTTP1ServerSym(loop,receiver,"kQueueState")){tsNodeHTTPServerCall(loop,handle,"abort")};callbacks:=tsNodeHTTPServerSym(loop,receiver,"kPendingCallbacks");if callbacks.kind==tsArrayKind{callbackArray:=(*tsArray)(callbacks.ref);pending:=make([]tsValue,len(callbackArray.values));copy(pending,callbackArray.values);callbackArray.resize(0);err:=arg(0);if tsNullish(err){err=tsFSException("Error","ERR_STREAM_DESTROYED","Cannot call write after a stream was destroyed")};for _,callback:=range pending{if callback.kind==tsFunctionKind{tsInternalHTTPNextTick(loop,callback,err)}}};return tsNodeHTTPOutgoingMethod(loop,receiver,"destroy",args)
	case "responseFlushHeaders":
		state:=tsNodeHTTPServerHeaderState(loop,receiver);states:=tsInternalHTTPGet(loop,tsInternalHTTPModule(loop),"NodeHTTPHeaderState");if tsStrictEqual(state,tsGet(loop,states,tsNodeHTTPServerAPIKey("sent"))){return tsU};if !tsStrictEqual(state,tsGet(loop,states,tsNodeHTTPServerAPIKey("assigned"))){tsNodeHTTPServerCall(loop,receiver,"_implicitHeader")};handle:=tsNodeHTTPServerSym(loop,receiver,"kHandle");if !tsNullish(handle){tsNodeHTTPServerCall(loop,handle,"flushHeaders")}else{tsNodeHTTPOutgoingCall(loop,receiver,"_send",tsNodeHTTPServerAPIKey(""))};tsNodeHTTPServerSetHeaderState(loop,receiver,"sent");return tsU
	case "responseSetTimeout":
		var callback tsValue=arg(1);if callback.kind==tsFunctionKind{tsNodeHTTPServerCall(loop,receiver,"on",tsNodeHTTPServerAPIKey("timeout"),callback)};socket:=tsNodeHTTPOutgoingSymbol(loop,receiver,tsNodeHTTPOutgoingSingleton.kSocket)
		if tsNullish(socket){listener:=tsFunctionValue(tsFunc(func(loop *tsLoop,values ...tsValue)tsValue{tsNodeHTTPServerCall(loop,tsArg(values,0),"setTimeout",arg(0));return tsU}));tsNodeHTTPServerCall(loop,receiver,"once",tsNodeHTTPServerAPIKey("socket"),listener)}else{tsNodeHTTPServerCall(loop,socket,"setTimeout",arg(0))};return receiver
	case "responseWriteEarlyHints":
		hints:=arg(0);if hints.kind!=tsObjectKind&&hints.kind!=tsInstanceKind{panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_TYPE","The \"hints\" argument must be of type object")})};link:=tsNodeHTTPServerAPIGet(loop,hints,"link");if tsNullish(link){return tsU};headers:=tsNewObject();tsSet(loop,headers,tsNodeHTTPServerAPIKey("Link"),link);for _,key:=range tsObjectKeys(hints).values{if tsText(key)!="link"{tsSet(loop,headers,key,tsGet(loop,hints,key))}};return tsNodeHTTPServerMethod(loop,"responseWriteInformation",receiver,[]tsValue{tsNumberValue(103),tsObjectValue(headers),arg(1)})
	case "responseCallPendingCallbacks":
		pending:=tsNodeHTTPServerSym(loop,receiver,"kPendingCallbacks");if pending.kind!=tsArrayKind{return tsU};array:=(*tsArray)(pending.ref);count:=len(array.values);values:=make([]tsValue,count);copy(values,array.values[:count]);for _,callback:=range values{if callback.kind==tsFunctionKind{tsInternalHTTPNextTick(loop,callback)}};if len(array.values)==count{array.resize(0)}else{array.values=append([]tsValue(nil),array.values[count:]...)};return tsU
	case "responseFlushWriteAccounting":
		tsNodeHTTPServerPut(loop,receiver,"kAccountingFlushScheduled",tsBooleanValue(false));bytes:=tsNumber(tsNodeHTTPServerSym(loop,receiver,"kBytesBuffered"));tsNodeHTTPServerPut(loop,receiver,"kBytesBuffered",tsNumberValue(0));if bytes>=tsNumber(tsNodeHTTPServerGet(loop,receiver,"writableHighWaterMark"))&&!tsTruthy(tsNodeHTTPServerAPIGet(loop,receiver,"destroyed"))&&!tsTruthy(tsNodeHTTPServerAPIGet(loop,receiver,"finished")){tsNodeHTTPOutgoingPut(loop,receiver,tsNodeHTTPOutgoingSingleton.kNeedDrain,tsBooleanValue(false));tsNodeHTTPServerCall(loop,receiver,"emit",tsNodeHTTPServerAPIKey("drain"))};return tsU
	case "responseFinish":tsNodeHTTPServerCall(loop,receiver,"emit",tsNodeHTTPServerAPIKey("prefinish"));return tsU
	case "responseCork","responseUncork":socket:=tsNodeHTTPServerMethod(loop,"responseSocket",receiver,nil);if !tsNullish(socket){method:="uncork";if name=="responseCork"{method="cork"};tsNodeHTTPServerCall(loop,socket,method)};return tsU
	case "responseSocket":
		queued:=tsNodeHTTP1ServerSym(loop,receiver,"kQueueState");if !tsNullish(queued){return tsNull};socket:=tsNodeHTTPOutgoingSymbol(loop,receiver,tsNodeHTTPOutgoingSingleton.kSocket);if !tsNullish(socket){return socket};fake:=tsNodeHTTPServerSym(loop,receiver,"fakeSocketSymbol");if !tsNullish(fake){return fake};fake=tsInternalHTTPFakeSocketCreate(loop,receiver);tsNodeHTTPServerPut(loop,receiver,"fakeSocketSymbol",fake);return fake
	case "responseSetSocket":socket:=arg(0);tsNodeHTTPServerPut(loop,receiver,"fakeSocketSymbol",socket);tsNodeHTTPOutgoingPut(loop,receiver,tsNodeHTTPOutgoingSingleton.kSocket,socket);return tsU
	case "responseHeadersProperty":return tsNodeHTTPOutgoingHeaders(loop,receiver)
	case "responseSetHeadersProperty":
		value:=arg(0);if tsNodeHTTPServerHeaderStateAssignedOrSent(loop,receiver){panic(tsThrown{tsFSException("Error","ERR_HTTP_HEADERS_SENT","Cannot set headers after they are sent to the client")})};tsNodeHTTPOutgoingPut(loop,receiver,tsNodeHTTPOutgoingSingleton.kOutHeaders,tsNull);if tsNullish(value){return tsU};if value.kind==tsArrayKind{for _,entry:=range (*tsArray)(value.ref).values{pair:=tsNodeHTTPIncomingArray(entry).values;if len(pair)>=2{tsNodeHTTPOutgoingMethod(loop,receiver,"appendHeader",[]tsValue{pair[0],pair[1]})}}}else if value.kind==tsObjectKind||value.kind==tsInstanceKind{for _,key:=range tsObjectKeys(value).values{tsNodeHTTPOutgoingMethod(loop,receiver,"setHeader",[]tsValue{key,tsGet(loop,value,key)})}};return tsU
	case "responseWritable":return tsBooleanValue(true)
	case "responseWritableNeedDrain":handle:=tsNodeHTTPServerSym(loop,receiver,"kHandle");queued:=tsNodeHTTP1ServerSym(loop,receiver,"kQueueState");return tsBooleanValue(!tsTruthy(tsNodeHTTPServerAPIGet(loop,receiver,"destroyed"))&&!tsTruthy(tsNodeHTTPServerAPIGet(loop,receiver,"finished"))&&(tsNumber(tsNodeHTTPServerGet(loop,handle,"bufferedAmount"))!=0||tsTruthy(tsNodeHTTPServerGet(loop,handle,"onwritable"))||tsTruthy(tsNodeHTTPServerSym(loop,receiver,"kNeedDrain"))||tsNumber(tsNodeHTTPServerSym(loop,receiver,"kBytesBuffered"))>=tsNumber(tsNodeHTTPServerGet(loop,receiver,"writableHighWaterMark"))||tsTruthy(tsNodeHTTPServerGet(loop,queued,"needDrain"))))
	case "responseWritableFinished":handle:=tsNodeHTTPServerSym(loop,receiver,"kHandle");return tsBooleanValue(tsTruthy(tsNodeHTTPServerAPIGet(loop,receiver,"finished"))&&(tsNullish(handle)||tsTruthy(tsNodeHTTPServerGet(loop,handle,"finished"))||(int(tsNumber(tsNodeHTTPServerGet(loop,handle,"flags")))&3)!=0))
	case "responseWritableLength":if tsTruthy(tsNodeHTTPServerMethod(loop,"responseWritableFinished",receiver,nil)){return tsNumberValue(0)};length:=tsNumber(tsNodeHTTPServerSym(loop,receiver,"kBytesBuffered"));handle:=tsNodeHTTPServerSym(loop,receiver,"kHandle");length+=tsNumber(tsNodeHTTPServerGet(loop,handle,"bufferedAmount"));queued:=tsNodeHTTP1ServerSym(loop,receiver,"kQueueState");length+=tsNumber(tsNodeHTTPServerGet(loop,queued,"bytes"));return tsNumberValue(length)
	case "responseHighWaterMark":socket:=tsNodeHTTPServerMethod(loop,"responseSocket",receiver,nil);mark:=tsNodeHTTPServerGet(loop,socket,"writableHighWaterMark");if tsIsNumeric(mark){return mark};mark=tsNodeHTTPServerSym(loop,receiver,"kHighWaterMark");if tsIsNumeric(mark){return mark};return tsNumberValue(65536)
	case "responseHeadersSent":return tsBooleanValue(tsNodeHTTPServerHeaderStateAssignedOrSent(loop,receiver))
	case "responseHeaderSet":return tsNodeHTTPServerMethod(loop,"responseSetSocket",receiver,args)
	case "responseHeaderGet":return tsNodeHTTPOutgoingHeaderGet(loop,receiver,arg(0))
	case "responseShouldKeepAliveGet":value:=tsNodeHTTPServerSym(loop,receiver,"kShouldKeepAlive");if tsIsUndefined(value){return tsBooleanValue(true)};return value
	case "responseShouldKeepAliveSet":tsNodeHTTPServerPut(loop,receiver,"kShouldKeepAlive",arg(0));return tsU
	case "responseClosed":return tsNodeHTTPServerAPIGet(loop,receiver,"_closed")
	case "responseWritableCorked":socket:=tsNodeHTTPServerMethod(loop,"responseSocket",receiver,nil);value:=tsNodeHTTPServerGet(loop,socket,"writableCorked");if tsIsNumeric(value){return value};return tsNumberValue(0)
	case "responseChunkedEncoding":return tsNodeHTTPServerSym(loop,receiver,"kChunkedEncoding")
	case "responseSetChunkedEncoding":tsNodeHTTPServerPut(loop,receiver,"kChunkedEncoding",arg(0));return tsU
	case "responseSetHeadersSent":if tsTruthy(arg(0)){tsNodeHTTPServerSetHeaderState(loop,receiver,"sent")}else{tsNodeHTTPServerSetHeaderState(loop,receiver,"none")};return tsU
	case "responseNoopSetter":return tsU
	case "responseFinishDrained":tsNodeHTTPServerResponseQueueFinish(loop,receiver,arg(0));return tsU
	case "responseEmit":
		event:=arg(0);if tsText(event)=="close"{var pending tsValue=tsNodeHTTPServerSym(loop,receiver,"kPendingFinish");if !tsIsUndefined(pending){tsNodeHTTPServerPut(loop,receiver,"kPendingFinish",tsU);tsNodeHTTPServerAPISet(loop,receiver,"_closed",tsBooleanValue(true));tsNodeHTTPServerMethod(loop,"responseCallPendingCallbacks",receiver,nil);finishClose:=tsFunctionValue(tsFunc(func(loop *tsLoop,values ...tsValue)tsValue{response,finishedCallback:=tsArg(values,0),tsArg(values,1);tsNodeHTTPServerMethod(loop,"responseEmitResponseFinishedThenClose",response,[]tsValue{finishedCallback});return tsU}));tsInternalHTTPNextTick(loop,finishClose,receiver,pending);return tsBooleanValue(false)};var closeCallback tsValue=tsInternalHTTPGet(loop,tsInternalHTTPModule(loop),"callCloseCallback");if closeCallback.kind==tsFunctionKind{tsCall(loop,closeCallback,receiver)}}
		state:=tsNodeEventState(receiver);if state==nil{return tsBooleanValue(false)};return tsBooleanValue(tsNodeEventEmit(loop,state,event,args[1:]))
	case "responseEmitResponseFinishedThenClose":
		var callback tsValue=arg(0);closeNT:=tsInternalHTTPGet(loop,tsInternalHTTPModule(loop),"emitCloseNT");if closeNT.kind==tsFunctionKind{tsCall(loop,closeNT,receiver)};tsNodeHTTPServerCall(loop,receiver,"emit",tsNodeHTTPServerAPIKey("finish"));if callback.kind==tsFunctionKind{tsCall(loop,callback)};state:=tsNodeEventState(receiver);if state!=nil{tsNodeEventEmit(loop,state,tsNodeHTTPServerAPIKey("close"),nil)};return tsU
	case "responseCaptureRejection":return tsNodeHTTPServerMethod(loop,"responseDestroy",receiver,args)
	}
	return tsU
}
func tsNodeHTTPServerMethodValue(loop *tsLoop,name string,length int) tsValue {method:=name;return tsNamedNativeMethod(method,length,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeHTTPServerMethod(loop,method,receiver,args)})}
func tsNodeHTTPServerListen(loop *tsLoop,server tsValue,args []tsValue) tsValue {
	if !tsNullish(tsNodeHTTPServerSym(loop,server,"serverSymbol")){panic(tsThrown{tsFSException("Error","ERR_SERVER_ALREADY_LISTEN","Listen method has been called more than once without closing.")})}
	options,callback:=tsArg(args,0),tsArg(args,len(args)-1);if options.kind==tsFunctionKind{callback=options;options=tsObjectValue(tsNewObject())};if options.kind==tsStringKind {text:=tsText(options);if port,err:=strconv.Atoi(text);err==nil{optionObject:=tsNewObject();optionObject.set("port",tsNumberValue(float64(port)));options=tsObjectValue(optionObject)}else{path:=tsNewObject();path.set("path",tsNodeHTTPServerAPIKey(text));options=tsObjectValue(path)}}else if options.kind!=tsObjectKind&&options.kind!=tsInstanceKind {normalized:=tsNodeNetNormalizeArgs(args);options=tsArg((*tsArray)(normalized.ref).values,0)}
	if callback.kind==tsFunctionKind{tsNodeHTTPServerCall(loop,server,"once",tsNodeHTTPServerAPIKey("listening"),callback)}
	tls:=tsNodeHTTPServerSym(loop,server,"tlsSymbol");var inner tsValue
	if !tsNullish(tls){tlsOptions:=tsNewObject();tsObjectSpread(tlsOptions,tls);tsObjectSpread(tlsOptions,options);tlsOptionsValue:=tsObjectValue(tlsOptions);inner=tsNodeTLSCreateServer(loop,[]tsValue{tlsOptionsValue})}else{inner=tsNodeNetServerCreate(loop,options)}
	tsNodeHTTPServerPut(loop,server,"serverSymbol",inner)
	connection:=tsFunctionValue(tsFunc(func(loop *tsLoop,values ...tsValue)tsValue{tsNodeHTTPServerEmit(loop,server,"connection",tsArg(values,0));return tsU}));tsNodeHTTPServerCall(loop,inner,"on",tsNodeHTTPServerAPIKey("connection"),connection)
	listening:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{tsNodeHTTPServerDispatch(loop,"emitListeningNextTick",server,[]tsValue{server,tsNodeHTTPServerAPIGet(loop,inner,"hostname"),tsNodeHTTPServerAPIGet(loop,inner,"port")});return tsU}));tsNodeHTTPServerCall(loop,inner,"once",tsNodeHTTPServerAPIKey("listening"),listening)
	secure:=tsFunctionValue(tsFunc(func(loop *tsLoop,values ...tsValue)tsValue{tsNodeHTTPServerEmit(loop,server,"secureConnection",tsArg(values,0));return tsU}));if !tsNullish(tls){tsNodeHTTPServerCall(loop,inner,"on",tsNodeHTTPServerAPIKey("secureConnection"),secure)}
	errorListener:=tsFunctionValue(tsFunc(func(loop *tsLoop,values ...tsValue)tsValue{if !tsNodeNetServerListening(loop,inner)&&tsStrictEqual(tsNodeHTTPServerSym(loop,server,"serverSymbol"),inner){tsNodeHTTPServerPut(loop,server,"serverSymbol",tsNull)};tsNodeHTTPServerDispatch(loop,"emitListenErrorNextTick",server,[]tsValue{server,tsArg(values,0)});return tsU}));tsNodeHTTPServerCall(loop,inner,"once",tsNodeHTTPServerAPIKey("error"),errorListener)
	closeListener:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{if tsNullish(tsNodeHTTPServerSym(loop,server,"serverSymbol")){tsNodeHTTPServerPut(loop,server,"kPendingDrainClose",tsBooleanValue(false));tsNodeHTTPServerEmit(loop,server,"close")};return tsU}));tsNodeHTTPServerCall(loop,inner,"once",tsNodeHTTPServerAPIKey("close"),closeListener)
	if !tsNullish(tls){tsNodeTLSServerListen(loop,inner,[]tsValue{options})}else{tsNodeNetListen(loop,inner,options)}
	return server
}
func tsNodeHTTPServerAPIModule(loop *tsLoop) tsValue {
	tsNodeHTTPServerAPIInit(loop);if !tsIsUndefined(tsNodeHTTPServerAPI.module){return tsNodeHTTPServerAPI.module}
	serverFn:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeHTTPServerCreate(loop,args)});serverFn.name="Server";serverFn.length=2;serverFn.constructible=true;serverFn.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeHTTPServerCreate(loop,args)};serverFn.receiverCall=func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeHTTPServerCreate(loop,args)};serverFn.properties=tsNewObject();serverFn.properties.set("prototype",tsObjectValue(tsNodeHTTPServerAPIServerProto()));tsNodeHTTPServerAPI.serverConstructor=tsFunctionValue(serverFn)
	responseFn:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeHTTPServerResponseCreate(loop,tsArg(args,0),tsArg(args,1))});responseFn.name="ServerResponse";responseFn.length=2;responseFn.constructible=true;responseFn.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeHTTPServerResponseCreate(loop,tsArg(args,0),tsArg(args,1))};responseFn.receiverCall=func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeHTTPServerResponseCreate(loop,tsArg(args,0),tsArg(args,1))};responseFn.properties=tsNewObject();responseFn.properties.set("prototype",tsObjectValue(tsNodeHTTPServerAPIResponseProto()));tsNodeHTTPServerAPI.responseConstructor=tsFunctionValue(responseFn)
	sp,rp:=tsNodeHTTPServerAPI.serverPrototype,tsNodeHTTPServerAPI.responsePrototype;sp.set("constructor",tsNodeHTTPServerAPI.serverConstructor);rp.set("constructor",tsNodeHTTPServerAPI.responseConstructor);rp.descriptors=map[string]*tsDescriptor{}
	serverMethods:=map[string]string{"ref":"serverRef","unref":"serverUnref","closeAllConnections":"serverCloseAll","closeIdleConnections":"serverCloseIdle","getConnections":"serverGetConnections","close":"serverClose","address":"serverAddress","listen":"serverListen","setTimeout":"serverSetTimeout","[captureRejection]":"serverCaptureRejection"};for name,dispatch:=range serverMethods{sp.set(name,tsNodeHTTPServerMethodValue(loop,dispatch,2))}
	responseMethods:=map[string]string{
		"assignSocket":"responseAssignSocket","detachSocket":"responseDetachSocket","writeHead":"responseWriteHead","_implicitHeader":"responseImplicitHeader",
		"writeContinue":"responseWriteContinue","writeProcessing":"responseWriteProcessing","destroy":"responseDestroy","flushHeaders":"responseFlushHeaders","setTimeout":"responseSetTimeout",
		"writeEarlyHints":"responseWriteEarlyHints","writeInformation":"responseWriteInformation","setHeader":"responseSetHeader","appendHeader":"responseAppendHeader",
		"setHeaders":"responseSetHeaders","removeHeader":"responseRemoveHeader","write":"responseWrite","end":"responseEnd","_writeRaw":"responseWriteRaw","_send":"responseSend",
		"_finish":"responseFinish","_callPendingCallbacks":"responseCallPendingCallbacks","cork":"responseCork","uncork":"responseUncork","emit":"responseEmit",
	}
	for name,dispatch:=range responseMethods{length:=2;if name=="detachSocket"||name=="_implicitHeader"||name=="writeContinue"||name=="writeProcessing"||name=="destroy"||name=="flushHeaders"||name=="_finish"||name=="_callPendingCallbacks"||name=="cork"||name=="uncork"{length=0};rp.set(name,tsNodeHTTPServerMethodValue(loop,dispatch,length));rp.descriptors[name]=&tsDescriptor{writable:true,configurable:true}}
	for _,name:=range []string{"getHeader","getHeaderNames","getRawHeaderNames","getHeaders","hasHeader","addTrailers"}{rp.set(name,tsNodeHTTPOutgoingGet(loop,tsNodeHTTPOutgoingSingleton.prototype,name));rp.descriptors[name]=&tsDescriptor{writable:true,configurable:true}}
	rk:=tsNodeHTTPServerAPIServerProto()
	rk.descriptors=map[string]*tsDescriptor{}
	rk.descriptors["listening"]=&tsDescriptor{accessor:true,configurable:true,enumerable:true,getter:tsNamedNativeMethod("get listening",0,func(loop *tsLoop,r tsValue,_ ...tsValue)tsValue{return tsBooleanValue(!tsNullish(tsNodeHTTPServerSym(loop,r,"serverSymbol")))} )}
	asyncKey:=tsPropertyKey(tsWellKnownSymbol("asyncDispose"));sp.set(asyncKey,tsNodeHTTPServerMethodValue(loop,"serverAsyncDispose",0));rp.descriptors["headersSent"]=&tsDescriptor{accessor:true,configurable:true,getter:tsNodeHTTPServerMethodValue(loop,"responseHeadersSent",0),setter:tsNodeHTTPServerMethodValue(loop,"responseSetHeadersSent",1)}
	responseAccessor:=func(name,getter,setter string,enumerable bool){descriptor:=&tsDescriptor{accessor:true,configurable:true,enumerable:enumerable};if getter!=""{descriptor.getter=tsNodeHTTPServerMethodValue(loop,getter,0)};if setter!=""{descriptor.setter=tsNodeHTTPServerMethodValue(loop,setter,1)};rp.descriptors[name]=descriptor}
	responseAccessor("headers","responseHeadersProperty","responseSetHeadersProperty",false)
	responseAccessor("socket","responseSocket","responseSetSocket",false);responseAccessor("connection","responseSocket","responseSetSocket",false)
	responseAccessor("chunkedEncoding","responseChunkedEncoding","responseSetChunkedEncoding",false)
	responseAccessor("writableCorked","responseWritableCorked","",false);responseAccessor("writable","responseWritable","responseNoopSetter",false)
	responseAccessor("writableNeedDrain","responseWritableNeedDrain","",false);responseAccessor("writableFinished","responseWritableFinished","",false);responseAccessor("writableLength","responseWritableLength","",false);responseAccessor("writableHighWaterMark","responseHighWaterMark","",false)
	responseAccessor("closed","responseClosed","",false);responseAccessor("shouldKeepAlive","responseShouldKeepAliveGet","responseShouldKeepAliveSet",false)
	module:=tsNewObject();module.set("Server",tsNodeHTTPServerAPIServerClass());module.set("ServerResponse",tsNodeHTTPServerAPIResponseClass());module.set("kConnectionsCheckingInterval",tsNodeHTTPServerSingleton.symbols["kConnectionsCheckingInterval"]);module.set("default",tsObjectValue(module));tsNodeHTTPServerAPI.module=tsObjectValue(module);return tsNodeHTTPServerAPI.module
}
func tsNodeHTTPServerAddon(loop *tsLoop,name string)(tsValue,bool){if name=="_http_server"{return tsNodeHTTPServerAPIModule(loop),true};return tsU,false}
`
