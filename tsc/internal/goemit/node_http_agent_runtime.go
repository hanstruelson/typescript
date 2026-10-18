package goemit

// nodeHTTPAgentRuntime adapts Bun's lib/_http_agent.js connection pool to the
// runtime's EventEmitter, net.Socket, and TLS socket values.
const nodeHTTPAgentRuntime = `
var tsNodeHTTPAgentSingleton struct {
	module, constructor, globalAgent tsValue
	prototype *tsObject
	kOnKeylog, kRequestOptions, kRequestAsyncResource, kOwnerFrame, kFreeSocketError, kDestroyOnRead, kMaybeKeylog tsValue
}
func tsNodeHTTPAgentKey(name string)tsValue{return tsStringReference(tsStringUTF8(name))}
func tsNodeHTTPAgentGet(loop *tsLoop,value tsValue,name string)tsValue{return tsGet(loop,value,tsNodeHTTPAgentKey(name))}
func tsNodeHTTPAgentSet(loop *tsLoop,value tsValue,name string,next tsValue){tsSet(loop,value,tsNodeHTTPAgentKey(name),next)}
func tsNodeHTTPAgentCall(loop *tsLoop,receiver tsValue,name string,args ...tsValue)tsValue{fn:=tsNodeHTTPAgentGet(loop,receiver,name);if fn.kind!=tsFunctionKind{return tsU};return tsCallReceiver(loop,fn,receiver,args...)}
func tsNodeHTTPAgentArray(value tsValue)*tsArray{if value.kind==tsArrayKind{return (*tsArray)(value.ref)};return &tsArray{}}
func tsNodeHTTPAgentArrayValue(array *tsArray)tsValue{return tsArrayValue(array)}
func tsNodeHTTPAgentMerge(loop *tsLoop,first,second tsValue)tsValue{out:=tsNewObject();tsObjectSpread(out,first);tsObjectSpread(out,second);return tsObjectValue(out)}
func tsNodeHTTPAgentPush(array tsValue,values ...tsValue){if array.kind!=tsArrayKind{return};target:=(*tsArray)(array.ref);target.values=append(target.values,values...)}
func tsNodeHTTPAgentRemove(array tsValue,index int){if array.kind!=tsArrayKind{return};target:=(*tsArray)(array.ref);if index>=0&&index<len(target.values){target.values=append(target.values[:index],target.values[index+1:]...)}}
func tsNodeHTTPAgentEmit(loop *tsLoop,value tsValue,event string,args ...tsValue){emit:=tsNodeHTTPAgentGet(loop,value,"emit");if emit.kind==tsFunctionKind{tsCallReceiver(loop,emit,value,append([]tsValue{tsNodeHTTPAgentKey(event)},args...)...)}}
func tsNodeHTTPAgentDropEmpty(loop *tsLoop,agent tsValue,name string){sockets:=tsNodeHTTPAgentGet(loop,agent,"sockets");entry:=tsGet(loop,sockets,tsNodeHTTPAgentKey(name));if entry.kind==tsArrayKind&&len((*tsArray)(entry.ref).values)==0{tsSet(loop,sockets,tsNodeHTTPAgentKey(name),tsU)}}
func tsNodeHTTPAgentReqSocket(loop *tsLoop,agent,req,socket tsValue){
	onSocket:=tsNodeHTTPAgentGet(loop,req,"onSocket");if onSocket.kind==tsFunctionKind{tsCallReceiver(loop,onSocket,req,socket)}
	timeout:=tsNodeHTTPAgentGet(loop,req,"timeout");agentTimeout:=tsNodeHTTPAgentGet(loop,tsNodeHTTPAgentGet(loop,agent,"options"),"timeout")
	if !tsIsUndefined(timeout)&&!tsStrictEqual(timeout,agentTimeout){setTimeout:=tsNodeHTTPAgentGet(loop,socket,"setTimeout");if setTimeout.kind==tsFunctionKind{tsCallReceiver(loop,setTimeout,socket,timeout)}}
}
func tsNodeHTTPAgentName(loop *tsLoop,options tsValue)string{
	host:=tsNodeHTTPAgentGet(loop,options,"host");if tsNullish(host)||tsText(host)==""{host=tsNodeHTTPAgentKey("localhost")}
	name:=tsText(host)+":";port:=tsNodeHTTPAgentGet(loop,options,"port");if tsTruthy(port){name+=tsText(port)};name+=":"
	local:=tsNodeHTTPAgentGet(loop,options,"localAddress");if tsTruthy(local){name+=tsText(local)}
	family:=tsNodeHTTPAgentGet(loop,options,"family");if tsStrictEqual(family,tsNumberValue(4))||tsStrictEqual(family,tsNumberValue(6)){name+=":"+tsText(family)}
	socketPath:=tsNodeHTTPAgentGet(loop,options,"socketPath");if tsTruthy(socketPath){name+=":"+tsText(socketPath)}
	return name
}
func tsNodeHTTPAgentNormalizeServerName(loop *tsLoop,options,req tsValue){
	servername:=tsNodeHTTPAgentGet(loop,options,"servername");if !tsNullish(servername)&&(!tsIsUndefined(servername)||tsText(servername)==""){return}
	host:=tsNodeHTTPAgentGet(loop,options,"host");header:=tsNodeHTTPAgentCall(loop,req,"getHeader",tsNodeHTTPAgentKey("host"))
	if tsTruthy(header){text:=tsText(header);if len(text)>0&&text[0]=='['{if index:=strings.IndexByte(text,']');index>=0{text=text[1:index]}}else if index:=strings.IndexByte(text,':');index>=0{text=text[:index]};host=tsNodeHTTPAgentKey(text)}
	if host.kind==tsStringKind&&(net.ParseIP(tsText(host))!=nil){host=tsNodeHTTPAgentKey("")}
	if tsIsUndefined(servername)||tsNullish(servername){tsNodeHTTPAgentSet(loop,options,"servername",host)}
}
func tsNodeHTTPAgentSetRequestSocket(loop *tsLoop,agent,req,socket tsValue){tsNodeHTTPAgentReqSocket(loop,agent,req,socket)}
func tsNodeHTTPAgentOnFree(loop *tsLoop,agent,socket,options tsValue)tsValue{
	if !tsTruthy(tsNodeHTTPAgentGet(loop,socket,"writable"))||tsNumber(tsNodeHTTPAgentGet(loop,socket,"readableLength"))>0{tsNodeHTTPAgentCall(loop,socket,"destroy");return tsU}
	name:=tsStringReference(tsStringUTF8(tsNodeHTTPAgentName(loop,options)));requests:=tsNodeHTTPAgentGet(loop,agent,"requests");queue:=tsGet(loop,requests,name)
	if queue.kind==tsArrayKind&&len((*tsArray)(queue.ref).values)>0{req:=(*tsArray)(queue.ref).values[0];tsNodeHTTPAgentRemove(queue,0);if len((*tsArray)(queue.ref).values)==0{tsSet(loop,requests,name,tsU)};tsNodeHTTPAgentReqSocket(loop,agent,req,socket);return tsU}
	request:=tsNodeHTTPAgentGet(loop,socket,"_httpMessage");identity:=tsGet(loop,options,tsInternalHTTPGet(loop,tsInternalHTTPModule(loop),"kPerRequestCheckServerIdentity"))
	if tsNullish(request)||!tsTruthy(tsNodeHTTPAgentGet(loop,request,"shouldKeepAlive"))||!tsTruthy(tsNodeHTTPAgentGet(loop,agent,"keepAlive"))||tsTruthy(identity){tsNodeHTTPAgentCall(loop,socket,"destroy");return tsU}
	freePools:=tsNodeHTTPAgentGet(loop,agent,"freeSockets");free:=tsGet(loop,freePools,name);if free.kind!=tsArrayKind{free=tsArrayValue(&tsArray{});tsSet(loop,freePools,name,free)}
	sockets:=tsGet(loop,tsNodeHTTPAgentGet(loop,agent,"sockets"),name);count:=len(tsNodeHTTPAgentArray(free).values)+len(tsNodeHTTPAgentArray(sockets).values)
	if tsNumber(tsNodeHTTPAgentGet(loop,agent,"totalSocketCount"))>tsNumber(tsNodeHTTPAgentGet(loop,agent,"maxTotalSockets"))||float64(count)>tsNumber(tsNodeHTTPAgentGet(loop,agent,"maxSockets"))||float64(len(tsNodeHTTPAgentArray(free).values))>=tsNumber(tsNodeHTTPAgentGet(loop,agent,"maxFreeSockets"))||!tsTruthy(tsNodeHTTPAgentCall(loop,agent,"keepSocketAlive",socket)){tsNodeHTTPAgentCall(loop,socket,"destroy");return tsU}
	tsNodeHTTPAgentSet(loop,socket,"_httpMessage",tsNull);tsNodeHTTPAgentCall(loop,agent,"removeSocket",socket,options);on:=tsNodeHTTPAgentGet(loop,socket,"once");if on.kind==tsFunctionKind{tsCallReceiver(loop,on,socket,tsNodeHTTPAgentKey("error"),tsNodeHTTPAgentSingleton.kFreeSocketError)};tsSet(loop,socket,tsNodeHTTPAgentSingleton.kDestroyOnRead,tsBooleanValue(true));tsNodeHTTPAgentPush(free,socket);return tsU
}
func tsNodeHTTPAgentDestroyOneFree(loop *tsLoop,agent tsValue){free:=tsNodeHTTPAgentGet(loop,agent,"freeSockets");for _,key:=range tsObjectKeys(free).values{for _,socket:=range tsNodeHTTPAgentArray(tsGet(loop,free,key)).values{if !tsTruthy(tsNodeHTTPAgentGet(loop,socket,"destroyed")){tsNodeHTTPAgentCall(loop,socket,"destroy");return}}}}
func tsNodeHTTPAgentEnableKeylog(loop *tsLoop,agent tsValue,event tsValue)tsValue{
	if tsText(event)!="keylog"||!tsIsUndefined(tsGet(loop,agent,tsNodeHTTPAgentSingleton.kOnKeylog)){return tsU};maybe:=tsGet(loop,agent,tsNodeHTTPAgentSingleton.kMaybeKeylog);tsNodeHTTPAgentCall(loop,agent,"removeListener",tsNodeHTTPAgentKey("newListener"),maybe)
	keylog:=tsNamedNativeMethod("onkeylog",1,func(loop *tsLoop,socket tsValue,args ...tsValue)tsValue{tsNodeHTTPAgentEmit(loop,agent,"keylog",tsArg(args,0),socket);return tsU});tsSet(loop,agent,tsNodeHTTPAgentSingleton.kOnKeylog,keylog)
	sockets:=tsNodeHTTPAgentGet(loop,agent,"sockets");for _,key:=range tsObjectKeys(sockets).values{for _,socket:=range tsNodeHTTPAgentArray(tsGet(loop,sockets,key)).values{tsNodeHTTPAgentCall(loop,socket,"on",tsNodeHTTPAgentKey("keylog"),keylog)}};return tsU
}
func tsNodeHTTPAgentInstallListeners(loop *tsLoop,agent,socket,options tsValue){
	on:=tsNodeHTTPAgentGet(loop,socket,"on");if on.kind!=tsFunctionKind{return}
	free:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{tsNodeHTTPAgentEmit(loop,agent,"free",socket,options);return tsU}))
	close:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{count:=tsNumber(tsNodeHTTPAgentGet(loop,agent,"totalSocketCount"))-1;tsNodeHTTPAgentSet(loop,agent,"totalSocketCount",tsNumberValue(count));tsNodeHTTPAgentCall(loop,agent,"removeSocket",socket,options);return tsU}))
	timeout:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{freePools:=tsNodeHTTPAgentGet(loop,agent,"freeSockets");for _,key:=range tsObjectKeys(freePools).values{pool:=tsGet(loop,freePools,key);for _,candidate:=range tsNodeHTTPAgentArray(pool).values{if tsStrictEqual(candidate,socket){tsNodeHTTPAgentCall(loop,socket,"destroy");return tsU}}};return tsU}))
	remove:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{tsNodeHTTPAgentSet(loop,agent,"totalSocketCount",tsNumberValue(tsNumber(tsNodeHTTPAgentGet(loop,agent,"totalSocketCount"))-1));tsNodeHTTPAgentCall(loop,agent,"removeSocket",socket,options);for _,event:=range []string{"close","free","timeout","agentRemove"}{tsNodeHTTPAgentCall(loop,socket,"removeListener",tsNodeHTTPAgentKey(event),tsU)};return tsU}))
	for _,entry:=range []struct{event string;fn tsValue}{{"free",free},{"close",close},{"timeout",timeout},{"agentRemove",remove}}{tsCallReceiver(loop,on,socket,tsNodeHTTPAgentKey(entry.event),entry.fn)}
	keylog:=tsGet(loop,agent,tsNodeHTTPAgentSingleton.kOnKeylog);if keylog.kind==tsFunctionKind{tsCallReceiver(loop,on,socket,tsNodeHTTPAgentKey("keylog"),keylog)}
}
func tsNodeHTTPAgentSocketCreated(loop *tsLoop,agent,req,name,err,socket,callback tsValue){
	if !tsNullish(err){tsNodeHTTPAgentDropEmpty(loop,agent,tsText(name));if callback.kind==tsFunctionKind{tsCall(loop,callback,err,tsNull)};return}
	pools:=tsNodeHTTPAgentGet(loop,agent,"sockets");pool:=tsGet(loop,pools,name);if pool.kind!=tsArrayKind{pool=tsArrayValue(&tsArray{});tsSet(loop,pools,name,pool)}
	tsNodeHTTPAgentPush(pool,socket);tsNodeHTTPAgentSet(loop,agent,"totalSocketCount",tsNumberValue(tsNumber(tsNodeHTTPAgentGet(loop,agent,"totalSocketCount"))+1))
	tsNodeHTTPAgentInstallListeners(loop,agent,socket,tsNodeHTTPAgentGet(loop,req,"_agentOptions"))
	if callback.kind==tsFunctionKind{tsCall(loop,callback,tsNull,socket)}
}
func tsNodeHTTPAgentCreateSocket(loop *tsLoop,agent,req,options,callback tsValue)tsValue{
	options=tsNodeHTTPAgentMerge(loop,options,tsNodeHTTPAgentGet(loop,agent,"options"));socketPath:=tsNodeHTTPAgentGet(loop,options,"socketPath");if tsTruthy(socketPath){tsNodeHTTPAgentSet(loop,options,"path",socketPath)}
	tsNodeHTTPAgentNormalizeServerName(loop,options,req)
	timeout:=tsNodeHTTPAgentGet(loop,req,"timeout");if !tsTruthy(timeout){timeout=tsNodeHTTPAgentGet(loop,tsNodeHTTPAgentGet(loop,agent,"options"),"timeout")};if tsTruthy(timeout){tsNodeHTTPAgentSet(loop,options,"timeout",timeout)}
	name:=tsStringReference(tsStringUTF8(tsNodeHTTPAgentName(loop,options)));tsNodeHTTPAgentSet(loop,options,"_agentKey",name);tsNodeHTTPAgentSet(loop,options,"encoding",tsNull);tsNodeHTTPAgentSet(loop,req,"_agentOptions",options)
	if tsTruthy(tsNodeHTTPAgentGet(loop,agent,"keepAlive")){tsNodeHTTPAgentSet(loop,options,"keepAlive",tsBooleanValue(true));tsNodeHTTPAgentSet(loop,options,"keepAliveInitialDelay",tsNodeHTTPAgentGet(loop,agent,"keepAliveMsecs"))}
	once:=false;finish:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{if once{return tsU};once=true;tsNodeHTTPAgentSocketCreated(loop,agent,req,name,tsArg(args,0),tsArg(args,1),callback);return tsU}))
	socket:=tsNodeHTTPAgentCall(loop,agent,"createConnection",options,finish)
	if !tsNullish(socket)&&!tsTruthy(tsGet(loop,socket,tsInternalHTTPGet(loop,tsInternalHTTPModule(loop),"kWaitForProxyTunnel"))){tsCall(loop,finish,tsNull,socket)}
	return socket
}
func tsNodeHTTPAgentCreateConnection(loop *tsLoop,agent tsValue,args ...tsValue)tsValue{
	netModule:=tsNodeNetModule(loop);normalized:=tsNodeNetNormalizeArgs(args);parts:=tsNodeHTTPAgentArray(normalized).values;options:=tsArg(parts,0);callback:=tsArg(parts,1)
	config:=tsGet(loop,agent,tsInternalHTTPGet(loop,tsInternalHTTPModule(loop),"kProxyConfig"))
	if tsNullish(config)||!tsTruthy(tsInternalHTTPGet(loop,tsInternalHTTPModule(loop),"checkShouldUseProxy")){return tsCallReceiver(loop,tsNodeHTTPAgentGet(loop,netModule,"createConnection"),netModule,args...)}
	check:=tsInternalHTTPGet(loop,tsInternalHTTPModule(loop),"checkShouldUseProxy");if !tsTruthy(tsCall(loop,check,config,options)){return tsCallReceiver(loop,tsNodeHTTPAgentGet(loop,netModule,"createConnection"),netModule,args...)}
	proxyOptions:=tsInternalHTTPGet(loop,config,"proxyConnectionOptions");protocol:=tsText(tsInternalHTTPGet(loop,config,"protocol"));if protocol=="https:"{tls:=tsNodeTLSModule(loop);return tsCallReceiver(loop,tsNodeHTTPAgentGet(loop,tls,"connect"),tls,proxyOptions,callback)}
	return tsCallReceiver(loop,tsNodeHTTPAgentGet(loop,netModule,"connect"),netModule,proxyOptions,callback)
}
func tsNodeHTTPAgentRemoveSocket(loop *tsLoop,agent,socket,options tsValue)tsValue{
	name:=tsStringReference(tsStringUTF8(tsNodeHTTPAgentName(loop,options)));for _,setName:=range []string{"sockets","freeSockets"}{set:=tsNodeHTTPAgentGet(loop,agent,setName);if setName=="freeSockets"&&tsTruthy(tsNodeHTTPAgentGet(loop,socket,"writable")){continue};pool:=tsGet(loop,set,name);if pool.kind==tsArrayKind{array:=(*tsArray)(pool.ref);for index,candidate:=range array.values{if tsStrictEqual(candidate,socket){tsNodeHTTPAgentRemove(pool,index);break}};if len(array.values)==0{tsSet(loop,set,name,tsU)}}}
	requests:=tsNodeHTTPAgentGet(loop,agent,"requests");queue:=tsGet(loop,requests,name);var req tsValue;queueName:=name;requestOptions:=options
	if queue.kind==tsArrayKind&&len((*tsArray)(queue.ref).values)>0{req=(*tsArray)(queue.ref).values[0]}else{for _,key:=range tsObjectKeys(requests).values{otherQueue:=tsGet(loop,requests,key);if otherQueue.kind==tsArrayKind&&len((*tsArray)(otherQueue.ref).values)>0{req=(*tsArray)(otherQueue.ref).values[0];queueName=key;requestOptions=tsGet(loop,req,tsNodeHTTPAgentSingleton.kRequestOptions);break}}}
	if !tsIsUndefined(req)&&!tsIsUndefined(requestOptions){cb:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{created:=tsArg(args,1);if !tsNullish(tsArg(args,0)){tsNodeHTTPAgentCall(loop,req,"onSocket",tsNull,tsArg(args,0));return tsU};pool:=tsGet(loop,requests,queueName);if pool.kind==tsArrayKind{for index,item:=range (*tsArray)(pool.ref).values{if tsStrictEqual(item,req){tsNodeHTTPAgentRemove(pool,index);break}}};tsNodeHTTPAgentReqSocket(loop,agent,req,created);return tsU}));tsNodeHTTPAgentCreateSocket(loop,agent,req,requestOptions,cb)}
	return tsU
}
func tsNodeHTTPAgentInit(loop *tsLoop){
	if tsNodeHTTPAgentSingleton.prototype!=nil{return}
	events:=tsNodeEventsModule(loop);eventProto:=tsGet(loop,events,tsNodeHTTPAgentKey("prototype"))
	tsNodeHTTPAgentSingleton.prototype=tsNewObject();tsNodeHTTPAgentSingleton.prototype.prototype=eventProto
	tsNodeHTTPAgentSingleton.kOnKeylog=tsNewSymbol(tsStringUTF8("onkeylog"));tsNodeHTTPAgentSingleton.kRequestOptions=tsNewSymbol(tsStringUTF8("requestOptions"));tsNodeHTTPAgentSingleton.kRequestAsyncResource=tsNewSymbol(tsStringUTF8("requestAsyncResource"));tsNodeHTTPAgentSingleton.kOwnerFrame=tsNewSymbol(tsStringUTF8("ownerFrame"));tsNodeHTTPAgentSingleton.kMaybeKeylog=tsNewSymbol(tsStringUTF8("maybeEnableKeylog"));tsNodeHTTPAgentSingleton.kDestroyOnRead=tsInternalHTTPGet(loop,tsInternalNetSymbolsModule(),"kDestroyOnRead");tsNodeHTTPAgentSingleton.kFreeSocketError=tsNamedNativeMethod("freeSocketErrorListener",1,func(loop *tsLoop,socket tsValue,_ ...tsValue)tsValue{tsNodeHTTPAgentCall(loop,socket,"destroy");tsNodeHTTPAgentEmit(loop,socket,"agentRemove");return tsU})
	p:=tsNodeHTTPAgentSingleton.prototype
	p.set("getName",tsNamedNativeMethod("getName",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsStringReference(tsStringUTF8(tsNodeHTTPAgentName(loop,tsArg(args,0))))}))
	p.set("createConnection",tsNamedNativeMethod("createConnection",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeHTTPAgentCreateConnection(loop,receiver,args...)}))
	p.set("addRequest",tsNamedNativeMethod("addRequest",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{
		req,options:=tsArg(args,0),tsArg(args,1);if options.kind==tsStringKind{legacy:=tsNewObject();legacy.set("host",options);legacy.set("port",tsArg(args,2));legacy.set("localAddress",tsArg(args,3));options=tsObjectValue(legacy)}
		options=tsNodeHTTPAgentMerge(loop,options,tsNodeHTTPAgentGet(loop,receiver,"options"));socketPath:=tsNodeHTTPAgentGet(loop,options,"socketPath");if tsTruthy(socketPath){tsNodeHTTPAgentSet(loop,options,"path",socketPath)};tsNodeHTTPAgentNormalizeServerName(loop,options,req)
		name:=tsStringReference(tsStringUTF8(tsNodeHTTPAgentName(loop,options)));sockets:=tsNodeHTTPAgentGet(loop,receiver,"sockets");pool:=tsGet(loop,sockets,name);if pool.kind!=tsArrayKind{pool=tsArrayValue(&tsArray{});tsSet(loop,sockets,name,pool)}
		freePools:=tsNodeHTTPAgentGet(loop,receiver,"freeSockets");free:=tsGet(loop,freePools,name);var socket tsValue
		if free.kind==tsArrayKind{arr:=(*tsArray)(free.ref);for len(arr.values)>0&&tsTruthy(tsNodeHTTPAgentGet(loop,arr.values[0],"destroyed")){tsNodeHTTPAgentRemove(free,0)};if len(arr.values)>0{index:=0;if tsText(tsNodeHTTPAgentGet(loop,receiver,"scheduling"))!="fifo"{index=len(arr.values)-1};socket=arr.values[index];tsNodeHTTPAgentRemove(free,index)}}
		if !tsIsUndefined(socket){tsNodeHTTPAgentCall(loop,receiver,"reuseSocket",socket,req);tsNodeHTTPAgentReqSocket(loop,receiver,req,socket);tsNodeHTTPAgentPush(pool,socket);return tsU}
		capacity:=len(tsNodeHTTPAgentArray(pool).values)+len(tsNodeHTTPAgentArray(free).values);maxSockets:=tsNumber(tsNodeHTTPAgentGet(loop,receiver,"maxSockets"));maxTotal:=tsNumber(tsNodeHTTPAgentGet(loop,receiver,"maxTotalSockets"))
		if float64(capacity)<maxSockets&&tsNumber(tsNodeHTTPAgentGet(loop,receiver,"totalSocketCount"))<maxTotal{cb:=tsFunctionValue(tsFunc(func(loop *tsLoop,values ...tsValue)tsValue{socket:=tsArg(values,1);if !tsNullish(tsArg(values,0)){tsNodeHTTPAgentCall(loop,req,"onSocket",tsNull,tsArg(values,0));return tsU};tsNodeHTTPAgentReqSocket(loop,receiver,req,socket);return tsU}));tsNodeHTTPAgentCreateSocket(loop,receiver,req,options,cb);return tsU}
		if tsTruthy(tsGet(loop,options,tsInternalHTTPGet(loop,tsInternalHTTPModule(loop),"kPerRequestCheckServerIdentity"))){tsNodeHTTPAgentDestroyOneFree(loop,receiver)}
		requests:=tsNodeHTTPAgentGet(loop,receiver,"requests");queue:=tsGet(loop,requests,name);if queue.kind!=tsArrayKind{queue=tsArrayValue(&tsArray{});tsSet(loop,requests,name,queue)};tsSet(loop,req,tsNodeHTTPAgentSingleton.kRequestOptions,options);tsNodeHTTPAgentPush(queue,req);return tsU
	}))
	p.set("createSocket",tsNamedNativeMethod("createSocket",3,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeHTTPAgentCreateSocket(loop,receiver,tsArg(args,0),tsArg(args,1),tsArg(args,2))}))
	p.set("removeSocket",tsNamedNativeMethod("removeSocket",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeHTTPAgentRemoveSocket(loop,receiver,tsArg(args,0),tsArg(args,1))}))
	p.set("keepSocketAlive",tsNamedNativeMethod("keepSocketAlive",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{socket:=tsArg(args,0);delay:=tsNodeHTTPAgentGet(loop,receiver,"keepAliveMsecs");tsNodeHTTPAgentCall(loop,socket,"setKeepAlive",tsBooleanValue(true),delay);tsNodeHTTPAgentCall(loop,socket,"unref");timeout:=tsNodeHTTPAgentGet(loop,tsNodeHTTPAgentGet(loop,receiver,"options"),"timeout");if !tsTruthy(timeout){timeout=tsNumberValue(0)};res:=tsNodeHTTPAgentGet(loop,tsNodeHTTPAgentGet(loop,socket,"_httpMessage"),"res");headers:=tsNodeHTTPAgentGet(loop,res,"headers");hint:=tsNodeHTTPAgentGet(loop,headers,"keep-alive");canKeep:=true;if !tsNullish(hint){text:=tsText(hint);if at:=strings.Index(text,"timeout=");at>=0{digits:="";for _,r:=range text[at+8:]{if r<'0'||r>'9'{break};digits+=string(r)};if seconds,err:=strconv.Atoi(digits);err==nil&&digits!=""{limit:=seconds*1000-int(tsNumber(tsNodeHTTPAgentGet(loop,receiver,"agentKeepAliveTimeoutBuffer")));if limit<=0{canKeep=false}else if tsNumber(timeout)==0||float64(limit)<tsNumber(timeout){timeout=tsNumberValue(float64(limit))}}}};if !tsStrictEqual(tsNodeHTTPAgentGet(loop,socket,"timeout"),timeout){tsNodeHTTPAgentCall(loop,socket,"setTimeout",timeout)};return tsBooleanValue(canKeep)}))
	p.set("reuseSocket",tsNamedNativeMethod("reuseSocket",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{socket,req:=tsArg(args,0),tsArg(args,1);tsNodeHTTPAgentCall(loop,socket,"removeListener",tsNodeHTTPAgentKey("error"),tsNodeHTTPAgentSingleton.kFreeSocketError);tsSet(loop,socket,tsNodeHTTPAgentSingleton.kDestroyOnRead,tsBooleanValue(false));tsNodeHTTPAgentSet(loop,req,"reusedSocket",tsBooleanValue(true));tsNodeHTTPAgentCall(loop,socket,"ref");return tsU}))
	p.set("destroy",tsNamedNativeMethod("destroy",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{for _,setName:=range []string{"freeSockets","sockets"}{set:=tsNodeHTTPAgentGet(loop,receiver,setName);for _,key:=range tsObjectKeys(set).values{pool:=tsGet(loop,set,key);for _,socket:=range tsNodeHTTPAgentArray(pool).values{tsNodeHTTPAgentCall(loop,socket,"destroy")}}};return tsU}))
	ctor:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeHTTPAgentCreate(loop,tsArg(args,0))});ctor.name="Agent";ctor.length=1;ctor.constructible=true;ctor.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeHTTPAgentCreate(loop,tsArg(args,0))};ctor.receiverCall=func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeHTTPAgentCreate(loop,tsArg(args,0))};ctor.properties=tsNewObject();ctor.properties.set("prototype",tsObjectValue(p));ctor.properties.set("defaultMaxSockets",tsNumberValue(math.Inf(1)));tsNodeHTTPAgentSingleton.constructor=tsFunctionValue(ctor)
}
func tsNodeHTTPAgentCreate(loop *tsLoop,options tsValue)tsValue{
	tsNodeHTTPAgentInit(loop);if tsNullish(options){options=tsObjectValue(tsNewObject())}
	object:=tsNewObject();object.prototype=tsObjectValue(tsNodeHTTPAgentSingleton.prototype);receiver:=tsNodeEventCreate(loop,tsObjectValue(object),tsU);freeListener:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{tsNodeHTTPAgentOnFree(loop,receiver,tsArg(args,0),tsArg(args,1));return tsU}));tsNodeHTTPAgentCall(loop,receiver,"on",tsNodeHTTPAgentKey("free"),freeListener);maybe:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeHTTPAgentEnableKeylog(loop,receiver,tsArg(args,0))}));tsSet(loop,receiver,tsNodeHTTPAgentSingleton.kMaybeKeylog,maybe);tsNodeHTTPAgentCall(loop,receiver,"on",tsNodeHTTPAgentKey("newListener"),maybe);settings:=tsNewObject();tsObjectSpread(settings,options);opts:=tsObjectValue(settings);tsNodeHTTPAgentSet(loop,receiver,"options",opts)
	defaultPort:=tsNodeHTTPAgentGet(loop,opts,"defaultPort");if !tsTruthy(defaultPort){defaultPort=tsNumberValue(80)};tsNodeHTTPAgentSet(loop,receiver,"defaultPort",defaultPort)
	protocol:=tsNodeHTTPAgentGet(loop,opts,"protocol");if !tsTruthy(protocol){protocol=tsNodeHTTPAgentKey("http:")};tsNodeHTTPAgentSet(loop,receiver,"protocol",protocol)
	if tsIsUndefined(tsNodeHTTPAgentGet(loop,opts,"noDelay")){settings.set("noDelay",tsBooleanValue(true))};settings.set("path",tsNull)
	for _,name:=range []string{"requests","sockets","freeSockets"}{tsNodeHTTPAgentSet(loop,receiver,name,tsObjectValue(tsNewObject()))}
	for name,value:=range map[string]tsValue{"keepAliveMsecs":tsNumberValue(1000),"keepAlive":tsBooleanValue(false),"maxSockets":tsNumberValue(math.Inf(1)),"maxFreeSockets":tsNumberValue(256),"scheduling":tsNodeHTTPAgentKey("lifo"),"maxTotalSockets":tsNumberValue(math.Inf(1)),"totalSocketCount":tsNumberValue(0),"agentKeepAliveTimeoutBuffer":tsNumberValue(1000)}{
		override:=tsNodeHTTPAgentGet(loop,opts,name);if !tsIsUndefined(override)&&!tsNullish(override){value=override};tsNodeHTTPAgentSet(loop,receiver,name,value)
	}
	scheduling:=tsText(tsNodeHTTPAgentGet(loop,receiver,"scheduling"));if scheduling!="fifo"&&scheduling!="lifo"{panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_VALUE","The value of scheduling is invalid")})}
	proxyEnv:=tsNodeHTTPAgentGet(loop,opts,"proxyEnv");if proxyEnv.kind==tsObjectKind||proxyEnv.kind==tsInstanceKind{internal:=tsInternalHTTPModule(loop);parse:=tsInternalHTTPGet(loop,internal,"parseProxyConfigFromEnv");config:=tsCall(loop,parse,proxyEnv,protocol,tsNodeHTTPAgentGet(loop,receiver,"keepAlive"));tsSet(loop,receiver,tsInternalHTTPGet(loop,internal,"kProxyConfig"),config)}
	frame:=tsCall(loop,tsInternalHTTPGet(loop,tsInternalAsyncContextFrameModule(loop),"currentGraphFrame"));if !tsIsUndefined(frame){tsSet(loop,receiver,tsNodeHTTPAgentSingleton.kOwnerFrame,frame)}
	return receiver
}
func tsNodeHTTPAgentEnvProxy(loop *tsLoop)tsValue{
	process:=tsNodeProcessModule(loop);env:=tsNodeHTTPAgentGet(loop,process,"env");flag:=tsNodeHTTPAgentGet(loop,env,"NODE_USE_ENV_PROXY")
	if !tsTruthy(flag)||tsText(flag)=="0"{return tsBooleanValue(false)}
	for _,name:=range []string{"http_proxy","HTTP_PROXY","https_proxy","HTTPS_PROXY"}{value:=tsNodeHTTPAgentGet(loop,env,name);if !tsNullish(value)&&tsText(value)!=""{parsed,err:=url.Parse(tsText(value));if err!=nil||parsed.Scheme==""||parsed.Host==""{panic(tsThrown{tsFSException("TypeError","ERR_INVALID_URL","Invalid URL")})}}}
	return tsBooleanValue(true)
}
func tsNodeHTTPAgentModule(loop *tsLoop)tsValue{
	if !tsIsUndefined(tsNodeHTTPAgentSingleton.module){return tsNodeHTTPAgentSingleton.module};tsNodeHTTPAgentInit(loop)
	out:=tsNewObject();out.set("Agent",tsNodeHTTPAgentSingleton.constructor);out.set("shouldUseEnvProxy",tsNamedNativeMethod("shouldUseEnvProxy",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNodeHTTPAgentEnvProxy(loop)}))
	options:=tsNewObject();options.set("keepAlive",tsBooleanValue(true));options.set("scheduling",tsNodeHTTPAgentKey("lifo"));options.set("timeout",tsNumberValue(5000))
	if tsTruthy(tsNodeHTTPAgentEnvProxy(loop)){options.set("proxyEnv",tsNodeHTTPAgentGet(loop,tsNodeProcessModule(loop),"env"))}
	tsNodeHTTPAgentSingleton.globalAgent=tsNodeHTTPAgentCreate(loop,tsObjectValue(options));out.set("globalAgent",tsNodeHTTPAgentSingleton.globalAgent);out.set("default",tsObjectValue(out));tsNodeHTTPAgentSingleton.module=tsObjectValue(out);return tsNodeHTTPAgentSingleton.module
}
`
