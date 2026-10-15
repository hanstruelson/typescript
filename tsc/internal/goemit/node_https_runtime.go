package goemit

// nodeHTTPSRuntime maps Bun's HTTPS wrapper to the public HTTP adapter and
// Go-backed TLS/network modules. The HTTP CONNECT path is represented with the
// same socket event sequence and per-agent TLS session cache as the source.
const nodeHTTPSRuntime = `
var tsNodeHTTPSingleton struct {
	module, constructor, globalAgent tsValue
	prototype, serverPrototype *tsObject
	sessions map[uintptr]*tsNodeHTTPSessionCache
	poolObjectIDs map[uintptr]uint64
	poolObjectCount uint64
}
type tsNodeHTTPSessionCache struct { entries map[string]tsValue; order []string }

func tsNodeHTTPSKey(name string) tsValue { return tsStringReference(tsStringUTF8(name)) }
func tsNodeHTTPSGet(loop *tsLoop,value tsValue,name string) tsValue { return tsGet(loop,value,tsNodeHTTPSKey(name)) }
func tsNodeHTTPSSet(loop *tsLoop,value tsValue,name string,next tsValue) { tsSet(loop,value,tsNodeHTTPSKey(name),next) }
func tsNodeHTTPSCall(loop *tsLoop,receiver tsValue,name string,args ...tsValue) tsValue { fn:=tsNodeHTTPSGet(loop,receiver,name);if fn.kind!=tsFunctionKind{return tsU};return tsCallReceiver(loop,fn,receiver,args...) }
func tsNodeHTTPSPointer(value tsValue) uintptr { if value.kind==tsObjectKind{return uintptr(unsafe.Pointer(value.ref))};return 0 }
func tsNodeHTTPSCacheFor(agent tsValue)*tsNodeHTTPSessionCache { if tsNodeHTTPSingleton.sessions==nil{tsNodeHTTPSingleton.sessions=map[uintptr]*tsNodeHTTPSessionCache{}};key:=tsNodeHTTPSPointer(agent);cache:=tsNodeHTTPSingleton.sessions[key];if cache==nil{cache=&tsNodeHTTPSessionCache{entries:map[string]tsValue{}};tsNodeHTTPSingleton.sessions[key]=cache};return cache }

func tsNodeHTTPSJSON(loop *tsLoop,value tsValue) string { global:=tsObjectValue(tsGlobalObject());json:=tsNodeHTTPSGet(loop,global,"JSON");stringify:=tsNodeHTTPSGet(loop,json,"stringify");if stringify.kind==tsFunctionKind{result:=tsCallReceiver(loop,stringify,json,value);if result.kind==tsStringKind{return tsText(result)}};return "undefined" }
func tsNodeHTTPSPoolPart(loop *tsLoop,value tsValue) string {
	if tsNullish(value){return tsText(value)}
	if value.kind==tsArrayKind{parts:=[]string{};for _,entry:=range (*tsArray)(value.ref).values{parts=append(parts,tsNodeHTTPSPoolPart(loop,entry))};return strings.Join(parts,",")}
	if value.kind==tsArrayBufferKind||value.kind==tsTypedArrayKind||value.kind==tsDataViewKind{return string(tsFSBytes(value))}
	if value.kind==tsObjectKind||value.kind==tsInstanceKind {
		buf:=tsNodeHTTPSGet(loop,value,"buf");if !tsIsUndefined(buf)&&!tsNullish(buf){return tsNodeHTTPSPoolPart(loop,buf)}
		ptr:=tsNodeHTTPSPointer(value);if ptr==0&&value.kind==tsInstanceKind{ptr=uintptr(unsafe.Pointer(value.ref))};if tsNodeHTTPSingleton.poolObjectIDs==nil{tsNodeHTTPSingleton.poolObjectIDs=map[uintptr]uint64{}};id:=tsNodeHTTPSingleton.poolObjectIDs[ptr];if id==0{tsNodeHTTPSingleton.poolObjectCount++;id=tsNodeHTTPSingleton.poolObjectCount;tsNodeHTTPSingleton.poolObjectIDs[ptr]=id};return "[object #"+strconv.FormatUint(id,10)+"]"
	}
	return tsText(value)
}
func tsNodeHTTPSPFXPoolKey(loop *tsLoop,pfx,passphrase tsValue)string{
	if pfx.kind!=tsArrayKind&&(pfx.kind!=tsObjectKind&&pfx.kind!=tsInstanceKind||tsIsUndefined(tsNodeHTTPSGet(loop,pfx,"buf"))){return tsNodeHTTPSPoolPart(loop,pfx)}
	entries:=[]tsValue{pfx};if pfx.kind==tsArrayKind{entries=(*tsArray)(pfx.ref).values};key:="";for _,entry:=range entries{raw:=entry;pass:=passphrase;if (entry.kind==tsObjectKind||entry.kind==tsInstanceKind){buf:=tsNodeHTTPSGet(loop,entry,"buf");if !tsIsUndefined(buf){raw=buf};entryPass:=tsNodeHTTPSGet(loop,entry,"passphrase");if tsTruthy(entryPass){pass=entryPass}};key+=":"+tsNodeHTTPSPoolPart(loop,raw)+":"+tsText(pass)};return key
}
func tsNodeHTTPSAgentName(loop *tsLoop,agent,options tsValue)tsValue{
	base:=tsNodeHTTPAgentModule(loop);ctor:=tsNodeHTTPSGet(loop,base,"Agent");prototype:=tsNodeHTTPSGet(loop,ctor,"prototype");baseName:=tsNodeHTTPSGet(loop,prototype,"getName");name:=tsCallReceiver(loop,baseName,agent,options);text:=tsText(name)+":"
	for _,field:=range []string{"ca","cert","clientCertEngine","ciphers","key","pfx","rejectUnauthorized","servername","host","minVersion","maxVersion","secureProtocol","crl","honorCipherOrder","ecdhCurve","dhparam","secureOptions","sessionIdContext","sigalgs","privateKeyIdentifier","privateKeyEngine"}{value:=tsNodeHTTPSGet(loop,options,field);if !tsTruthy(value){if field=="rejectUnauthorized"||field=="honorCipherOrder"||field=="secureOptions"{if !tsIsUndefined(value){text+=":"+tsText(value)}};continue};if field=="clientCertEngine"||field=="ciphers"||field=="servername"||field=="minVersion"||field=="maxVersion"||field=="secureProtocol"||field=="honorCipherOrder"||field=="ecdhCurve"||field=="sessionIdContext"||field=="privateKeyIdentifier"||field=="privateKeyEngine"{text+=":"+tsText(value)}else if field=="pfx"{text+=":"+tsNodeHTTPSPFXPoolKey(loop,value,tsNodeHTTPSGet(loop,options,"passphrase"))}else if field=="sigalgs"{text+=":"+tsNodeHTTPSJSON(loop,value)}else{text+=":"+tsNodeHTTPSPoolPart(loop,value)}}
	for _,field:=range []string{"certFile","keyFile","caFile"}{value:=tsNodeHTTPSGet(loop,options,field);if tsTruthy(value){text+=":"+field+"="+tsNodeHTTPSJSON(loop,value)}}
	if tsTruthy(tsNodeHTTPSGet(loop,options,"allowPartialTrustChain")){text+=":allowPartialTrustChain"};if secure:=tsNodeHTTPSGet(loop,options,"secureContext");tsTruthy(secure){text+=":secureContext="+tsNodeHTTPSPoolPart(loop,secure)}
	key:=tsNodeHTTPGet(loop,tsInternalHTTPModule(loop),"kPerRequestCheckServerIdentity");if perRequest:=tsGet(loop,options,key);tsTruthy(perRequest){text+=":"+tsText(perRequest)}
	return tsNodeHTTPSKey(text)
}
func tsNodeHTTPSGetSession(loop *tsLoop,agent,key tsValue)tsValue{return tsNodeHTTPSCacheFor(agent).entries[tsText(key)]}
func tsNodeHTTPSCacheSession(loop *tsLoop,agent,key,session tsValue)tsValue{
	cache:=tsNodeHTTPSCacheFor(agent);name:=tsText(key);maximum:=int(tsNumber(tsNodeHTTPSGet(loop,agent,"maxCachedSessions")));if maximum==0{return tsU};if _,exists:=cache.entries[name];exists{cache.entries[name]=session;return tsU};if len(cache.order)>=maximum{old:=cache.order[0];cache.order=cache.order[1:];delete(cache.entries,old)};cache.order=append(cache.order,name);cache.entries[name]=session;return tsU
}
func tsNodeHTTPSEvictSession(loop *tsLoop,agent,key tsValue){cache:=tsNodeHTTPSCacheFor(agent);name:=tsText(key);for index,value:=range cache.order{if value==name{cache.order=append(cache.order[:index],cache.order[index+1:]...);delete(cache.entries,name);return}}}

func tsNodeHTTPSTunnelConfig(loop *tsLoop,agent,options tsValue)tsValue{
	internal:=tsInternalHTTPModule(loop);key:=tsNodeHTTPSGet(loop,internal,"kProxyConfig");config:=tsGet(loop,agent,key);if tsNullish(config){return tsNull};protocol:=tsNodeHTTPSGet(loop,options,"protocol");if tsIsUndefined(protocol){protocol=tsNodeHTTPSGet(loop,agent,"protocol")};if tsText(protocol)!="https:"{return tsNull}
	check:=tsNodeHTTPSGet(loop,internal,"checkShouldUseProxy");if check.kind!=tsFunctionKind||!tsTruthy(tsCall(loop,check,config,options)){return tsNull}
	host:=tsNodeHTTPSGet(loop,options,"host");ip:=tsNodeHTTPSGet(loop,tsInternalNetIPModule(),"isIP");ipType:=int(tsNumber(tsCall(loop,ip,host)));requestHost:=tsText(host);if ipType==6{requestHost="["+requestHost+"]"};port:=tsNodeHTTPSGet(loop,options,"port");if !tsTruthy(port){port=tsNodeHTTPSGet(loop,agent,"defaultPort")};endpoint:=requestHost+":"+tsText(port)
	common:=tsNodeHTTPCommonModule(loop);tsNodeHTTPSCall(loop,common,"validateHeaderValue",tsNodeHTTPSKey("host"),tsNodeHTTPSKey(endpoint))
	payload:="CONNECT "+endpoint+" HTTP/1.1\r\n";auth:=tsNodeHTTPSGet(loop,config,"auth");if tsTruthy(auth){payload+="proxy-authorization: "+tsText(auth)+"\r\n"};if tsTruthy(tsNodeHTTPSGet(loop,agent,"keepAlive"))||!math.IsInf(tsNumber(tsNodeHTTPSGet(loop,agent,"maxSockets")),1){payload+="proxy-connection: keep-alive\r\n"};payload+="host: "+endpoint+"\r\n\r\n"
	requestOptions:=tsNewObject();tsObjectSpread(requestOptions,options);if !tsTruthy(tsNodeHTTPSGet(loop,requestOptions,"servername"))&&ipType==0{requestOptions.set("servername",host)}
	out:=tsNewObject();out.set("proxyTunnelPayload",tsNodeHTTPSKey(payload));out.set("requestOptions",tsObjectValue(requestOptions));return tsObjectValue(out)
}
func tsNodeHTTPSStatusLine(data string)string{if index:=strings.Index(data,"\r\n");index>=0{return data[:index]};return data}
func tsNodeHTTPSAfterTunnel(loop *tsLoop,agent,proxySocket,options,tunnel,callback tsValue){
	payload:=tsText(tsNodeHTTPSGet(loop,tunnel,"proxyTunnelPayload"));done:=false;buffer:=""
	var readFn tsValue
	readFn=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{
		if done{return tsU};for {chunk:=tsNodeHTTPSCall(loop,proxySocket,"read");if tsNullish(chunk){break};buffer+=string(tsFSBytes(chunk));end:=strings.Index(buffer,"\r\n\r\n");if end<0{continue};status:=tsNodeHTTPSStatusLine(buffer);fields:=strings.Fields(status);code:="";if len(fields)>1{code=fields[1]};done=true;tsNodeHTTPSCall(loop,proxySocket,"removeListener",tsNodeHTTPSKey("readable"),readFn)
			if code!="200"{target:="";if fields:=strings.Fields(strings.Split(payload,"\r")[0]);len(fields)>1{target=fields[1]};config:=tsGet(loop,agent,tsNodeHTTPSGet(loop,tsInternalHTTPModule(loop),"kProxyConfig"));href:=tsNodeHTTPSGet(loop,config,"href");err:=tsFSException("Error","ERR_PROXY_TUNNEL","Failed to establish tunnel to "+target+" via "+tsText(href)+": "+status);statusCode,_:=strconv.Atoi(code);errObj:=tsNodeHTTPSAddField(loop,err,"statusCode",tsNumberValue(float64(statusCode)));tsNodeHTTPSCall(loop,proxySocket,"destroy");if callback.kind==tsFunctionKind{tsCall(loop,callback,errObj,proxySocket)};return tsU}
			requestOptions:=tsNodeHTTPSGet(loop,tunnel,"requestOptions");tsNodeHTTPSSet(loop,tunnel,"requestOptions",tsNull);tsNodeHTTPSSet(loop,requestOptions,"socket",proxySocket);tls:=tsNodeTLSModule(loop);var secure tsValue
			connected:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{if callback.kind==tsFunctionKind{tsCall(loop,callback,tsNull,secure)};return tsU}));secure=tsCall(loop,tsNodeHTTPSGet(loop,tls,"connect"),requestOptions,connected);if secure.kind==tsObjectKind{tsNodeHTTPSCall(loop,secure,"on",tsNodeHTTPSKey("free"),tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{tsNodeHTTPSCall(loop,proxySocket,"emit",tsNodeHTTPSKey("free"));return tsU})));tsNodeHTTPSCall(loop,secure,"on",tsNodeHTTPSKey("error"),tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{if callback.kind==tsFunctionKind{tsCall(loop,callback,tsArg(args,0),secure)};return tsU}))) }
			return tsU
		};return tsU
	}))
	endFn:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{if done{return tsU};done=true;err:=tsFSException("Error","ERR_PROXY_TUNNEL","Connection to establish proxy tunnel ended unexpectedly");if callback.kind==tsFunctionKind{tsCall(loop,callback,err,proxySocket)};return tsU}))
	proxyErr:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{if done{return tsU};done=true;if callback.kind==tsFunctionKind{tsCall(loop,callback,tsArg(args,0),proxySocket)};return tsU}))
	for _,entry:=range []struct{event string;fn tsValue}{{"readable",readFn},{"end",endFn},{"error",proxyErr}}{tsNodeHTTPSCall(loop,proxySocket,"on",tsNodeHTTPSKey(entry.event),entry.fn)}
	timeout:=tsNodeHTTPSGet(loop,options,"timeout");if tsTruthy(timeout){timeoutFn:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{if done{return tsU};done=true;err:=tsFSException("Error","ERR_PROXY_TUNNEL","Connection to establish proxy tunnel timed out after "+tsText(timeout)+"ms");tsNodeHTTPSAddField(loop,err,"proxyTunnelTimeout",timeout);if callback.kind==tsFunctionKind{tsCall(loop,callback,err,proxySocket)};return tsU}));tsNodeHTTPSCall(loop,proxySocket,"setTimeout",timeout,timeoutFn)}
	tsNodeHTTPSCall(loop,proxySocket,"write",tsNodeHTTPSKey(payload));
}
func tsNodeHTTPSAddField(loop *tsLoop,value tsValue,name string,field tsValue)tsValue{if value.kind==tsObjectKind{(*tsObject)(value.ref).set(name,field)};return value}
func tsNodeHTTPSCreateConnection(loop *tsLoop,agent tsValue,args []tsValue)tsValue{
	options:=tsNewObject();var source tsValue
	if len(args)>0&&(args[0].kind==tsObjectKind||args[0].kind==tsInstanceKind){source=args[0]}else if len(args)>1&&(args[1].kind==tsObjectKind||args[1].kind==tsInstanceKind){source=args[1]}else if len(args)>2&&(args[2].kind==tsObjectKind||args[2].kind==tsInstanceKind){source=args[2]}
	if !tsIsUndefined(source){tsObjectSpread(options,source)};if len(args)>0&&args[0].kind==tsNumberKind{options.set("port",args[0])};if len(args)>1&&args[1].kind==tsStringKind{options.set("host",args[1])}
	callback:=tsU;if len(args)>0&&args[len(args)-1].kind==tsFunctionKind{callback=args[len(args)-1]}
	internal:=tsInternalHTTPModule(loop);cacheKey:=tsNodeHTTPSGet(loop,options,"_agentKey");perRequest:=tsGet(loop,options,tsNodeHTTPSGet(loop,internal,"kPerRequestCheckServerIdentity"));reuse:=tsTruthy(cacheKey)&&!tsTruthy(perRequest);if reuse{session:=tsNodeHTTPSGetSession(loop,agent,cacheKey);if !tsIsUndefined(session){copy:=tsNewObject();copy.set("session",session);tsObjectSpread(copy,tsObjectValue(options));options=copy}}
	tunnel:=tsNodeHTTPSTunnelConfig(loop,agent,tsObjectValue(options));tls:=tsNodeTLSModule(loop)
	if tsIsNull(tunnel){secure:=tsCall(loop,tsNodeHTTPSGet(loop,tls,"connect"),tsObjectValue(options),callback);if reuse{sessionFn:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{tsNodeHTTPSCacheSession(loop,agent,cacheKey,tsArg(args,0));return tsU}));closeFn:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{if tsTruthy(tsArg(args,0)){tsNodeHTTPSEvictSession(loop,agent,cacheKey)};return tsU}));tsNodeHTTPSCall(loop,secure,"on",tsNodeHTTPSKey("session"),sessionFn);tsNodeHTTPSCall(loop,secure,"once",tsNodeHTTPSKey("close"),closeFn)};return secure}
	config:=tsGet(loop,agent,tsNodeHTTPSGet(loop,internal,"kProxyConfig"));connectOptions:=tsNewObject();tsObjectSpread(connectOptions,tsNodeHTTPSGet(loop,config,"proxyConnectionOptions"));net:=tsNodeNetModule(loop);protocol:=tsText(tsNodeHTTPSGet(loop,config,"protocol"));connect:=tsNodeHTTPSGet(loop,net,"connect");if protocol=="https:"{connect=tsNodeHTTPSGet(loop,tls,"connect")}
	connected:=tsFunctionValue(tsFunc(func(loop *tsLoop,socketArgs ...tsValue)tsValue{socket:=tsArg(socketArgs,0);tsSet(loop,socket,tsNodeHTTPSGet(loop,internal,"kWaitForProxyTunnel"),tsBooleanValue(true));tsNodeHTTPSAfterTunnel(loop,agent,socket,tsObjectValue(options),tunnel,callback);return tsU}));socket:=tsCall(loop,connect,tsObjectValue(connectOptions),connected);if socket.kind==tsObjectKind{tsSet(loop,socket,tsNodeHTTPSGet(loop,internal,"kWaitForProxyTunnel"),tsBooleanValue(true))};return socket
}
func tsNodeHTTPSAgentCreate(loop *tsLoop,options tsValue)tsValue{
	baseModule:=tsNodeHTTPAgentModule(loop);baseCtor:=tsNodeHTTPSGet(loop,baseModule,"Agent");if tsNodeHTTPSingleton.prototype==nil{proto:=tsNewObject();proto.prototype=tsObjectValue(tsNodeHTTPAgentSingleton.prototype);tsNodeHTTPSingleton.prototype=proto}
	copy:=tsNewObject();tsObjectSpread(copy,options);if tsIsUndefined(tsNodeHTTPSGet(loop,tsObjectValue(copy),"defaultPort"))||tsNullish(tsNodeHTTPSGet(loop,tsObjectValue(copy),"defaultPort")){copy.set("defaultPort",tsNumberValue(443))};if tsIsUndefined(tsNodeHTTPSGet(loop,tsObjectValue(copy),"protocol"))||tsNullish(tsNodeHTTPSGet(loop,tsObjectValue(copy),"protocol")){copy.set("protocol",tsNodeHTTPSKey("https:"))}
	agent:=tsConstruct(baseCtor,tsObjectValue(copy));if agent.kind==tsObjectKind{(*tsObject)(agent.ref).prototype=tsObjectValue(tsNodeHTTPSingleton.prototype)};maximum:=tsNodeHTTPSGet(loop,tsObjectValue(copy),"maxCachedSessions");if tsIsUndefined(maximum){maximum=tsNumberValue(100)};tsNodeHTTPSSet(loop,agent,"maxCachedSessions",maximum);tsNodeHTTPSCacheFor(agent)
	return agent
}
func tsNodeHTTPSAgentInit(loop *tsLoop){
	if tsNodeHTTPSingleton.constructor.kind!=tsUndefinedKind{return};base:=tsNodeHTTPAgentModule(loop);baseCtor:=tsNodeHTTPSGet(loop,base,"Agent");baseProto:=tsNodeHTTPSGet(loop,baseCtor,"prototype");proto:=tsNewObject();proto.prototype=baseProto;tsNodeHTTPSingleton.prototype=proto
	proto.set("getName",tsNamedNativeMethod("getName",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeHTTPSAgentName(loop,receiver,tsArg(args,0))}))
	proto.set("createConnection",tsNamedNativeMethod("createConnection",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeHTTPSCreateConnection(loop,receiver,args)}))
	proto.set("_getSession",tsNamedNativeMethod("_getSession",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeHTTPSGetSession(loop,receiver,tsArg(args,0))}))
	proto.set("_cacheSession",tsNamedNativeMethod("_cacheSession",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeHTTPSCacheSession(loop,receiver,tsArg(args,0),tsArg(args,1))}))
	proto.set("_evictSession",tsNamedNativeMethod("_evictSession",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{tsNodeHTTPSEvictSession(loop,receiver,tsArg(args,0));return tsU}))
	ctor:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeHTTPSAgentCreate(loop,tsArg(args,0))});ctor.name="Agent";ctor.length=1;ctor.constructible=true;ctor.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeHTTPSAgentCreate(loop,tsArg(args,0))};ctor.receiverCall=func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeHTTPSAgentCreate(loop,tsArg(args,0))};ctor.properties=tsNewObject();ctor.properties.set("prototype",tsObjectValue(proto));tsNodeHTTPSingleton.constructor=tsFunctionValue(ctor)
}

func tsNodeHTTPSRequest(loop *tsLoop,args []tsValue)tsValue{
	values:=append([]tsValue{},args...);options:=tsNewObject();global:=tsObjectValue(tsGlobalObject());urlCtor:=tsNodeHTTPSGet(loop,global,"URL")
	if len(values)>0&&values[0].kind==tsStringKind{url:=tsConstruct(urlCtor,values[0]);options=(*tsObject)(tsInternalURLToHTTPOptions(loop,[]tsValue{url}).ref);values=values[1:]}else if len(values)>0&&tsInstanceOf(values[0],urlCtor){options=(*tsObject)(tsInternalURLToHTTPOptions(loop,[]tsValue{values[0]}).ref);values=values[1:]}
	if len(values)>0&&!tsNullish(values[0])&&values[0].kind!=tsFunctionKind{tsObjectSpread(options,tsObjectValue(options));tsObjectSpread(options,values[0]);values=values[1:]}
	options.set("_defaultAgent",tsNodeHTTPSingleton.globalAgent);callArgs:=[]tsValue{tsObjectValue(options)};callArgs=append(callArgs,values...);http:=tsNodeHTTPModule(loop);return tsConstruct(tsNodeHTTPSGet(loop,http,"ClientRequest"),callArgs...)
}
func tsNodeHTTPSGetRequest(loop *tsLoop,args []tsValue)tsValue{request:=tsNodeHTTPSRequest(loop,args);tsNodeHTTPSCall(loop,request,"end");return request}

func tsNodeHTTPSServerCreate(loop *tsLoop,args []tsValue)tsValue{
	options,listener:=tsArg(args,0),tsArg(args,1);if options.kind==tsFunctionKind{listener=options;options=tsObjectValue(tsNewObject())}else if tsNullish(options){options=tsObjectValue(tsNewObject())}else{tsCall(loop,tsNodeHTTPSGet(loop,tsInternalValidatorsModule(loop),"validateObject"),options,tsNodeHTTPSKey("options"));copy:=tsNewObject();tsObjectSpread(copy,options);options=tsObjectValue(copy)}
	if !tsTruthy(tsNodeHTTPSGet(loop,options,"ALPNProtocols"))&&!tsTruthy(tsNodeHTTPSGet(loop,options,"ALPNCallback")){protocols:=&tsArray{values:[]tsValue{tsNodeHTTPSKey("http/1.1")}};tsNodeHTTPSSet(loop,options,"ALPNProtocols",tsArrayValue(protocols))}
	server:=tsNodeHTTPServerCreate(loop,[]tsValue{options,listener});if server.kind==tsObjectKind{if tsNodeHTTPSingleton.serverPrototype==nil{proto:=tsNewObject();proto.prototype=tsNodeHTTPSGet(loop,tsNodeHTTPSGet(loop,tsNodeHTTPModule(loop),"Server"),"prototype");tsNodeHTTPSingleton.serverPrototype=proto};(*tsObject)(server.ref).prototype=tsObjectValue(tsNodeHTTPSingleton.serverPrototype)}
	alpn:=tsNodeHTTPSGet(loop,options,"ALPNProtocols");if tsTruthy(alpn){tls:=tsNodeTLSModule(loop);converted:=tsNodeHTTPSCall(loop,tls,"convertALPNProtocols",alpn,server);tsNodeHTTPSSet(loop,server,"ALPNProtocols",converted)}
	blockList:=tsNodeHTTPSGet(loop,options,"blockList");if tsTruthy(blockList){net:=tsNodeNetModule(loop);ctor:=tsNodeHTTPSGet(loop,net,"BlockList");check:=tsNodeHTTPSGet(loop,ctor,"isBlockList");if !tsTruthy(tsCall(loop,check,blockList)){panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_TYPE","The \"options.blockList\" property must be an instance of net.BlockList")})};tsNodeHTTPSSet(loop,server,"blockList",blockList)}
	internal:=tsInternalHTTPModule(loop);tlsSymbol:=tsNodeHTTPSGet(loop,internal,"isTlsSymbol");tsSet(loop,server,tlsSymbol,tsBooleanValue(true));tsNodeHTTPSSet(loop,server,"ALPNCallback",tsNodeHTTPSGet(loop,options,"ALPNCallback"));return server
}
func tsNodeHTTPSModule(loop *tsLoop)tsValue{
	if !tsIsUndefined(tsNodeHTTPSingleton.module){return tsNodeHTTPSingleton.module}
	tsNodeHTTPSAgentInit(loop);process:=tsNodeProcessModule(loop);env:=tsNodeHTTPSGet(loop,process,"env");options:=tsNewObject();options.set("keepAlive",tsBooleanValue(true));options.set("scheduling",tsNodeHTTPSKey("lifo"));options.set("timeout",tsNumberValue(5000));if tsTruthy(tsNodeHTTPSCall(loop,tsNodeHTTPAgentModule(loop),"shouldUseEnvProxy")){options.set("proxyEnv",env)}
	tsNodeHTTPSingleton.globalAgent=tsNodeHTTPSAgentCreate(loop,tsObjectValue(options))
	http:=tsNodeHTTPModule(loop);module:=tsNewObject();module.set("Agent",tsNodeHTTPSingleton.constructor);module.set("globalAgent",tsNodeHTTPSingleton.globalAgent)
	serverCtor:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeHTTPSServerCreate(loop,args)});serverCtor.name="Server";serverCtor.length=2;serverCtor.constructible=true;serverCtor.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeHTTPSServerCreate(loop,args)};serverCtor.receiverCall=func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeHTTPSServerCreate(loop,args)};serverCtor.properties=tsNewObject();if tsNodeHTTPSingleton.serverPrototype==nil{proto:=tsNewObject();proto.prototype=tsNodeHTTPSGet(loop,tsNodeHTTPSGet(loop,http,"Server"),"prototype");tsNodeHTTPSingleton.serverPrototype=proto};serverCtor.properties.set("prototype",tsObjectValue(tsNodeHTTPSingleton.serverPrototype));tsNodeHTTPSingleton.serverPrototype.set("constructor",tsFunctionValue(serverCtor));tsNodeHTTPSingleton.serverPrototype.prototype=tsNodeHTTPSGet(loop,tsNodeHTTPSGet(loop,http,"Server"),"prototype")
	module.set("Server",tsFunctionValue(serverCtor));module.set("createServer",tsNamedNativeMethod("createServer",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeHTTPSServerCreate(loop,args)}));module.set("request",tsNamedNativeMethod("request",3,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeHTTPSRequest(loop,args)}));module.set("get",tsNamedNativeMethod("get",3,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeHTTPSGetRequest(loop,args)}));module.set("default",tsObjectValue(module));tsNodeHTTPSingleton.module=tsObjectValue(module);return tsNodeHTTPSingleton.module
}
`
