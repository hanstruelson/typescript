package goemit

// nodeClusterRuntime ports the primary-process cluster facade and its
// CommonJS surface. Socket-handle scheduling is held at the Worker boundary;
// the primary still exposes settings, worker events, and process forking.
const nodeClusterRuntime = `
type tsNodeClusterState struct {
	sync.Once
	module, workers, settings, intercom, workerConstructor tsValue
	handles map[string]tsValue
	ids int
	initialized bool
	schedulingPolicy int
}
var tsNodeClusterSingleton tsNodeClusterState

func tsNodeClusterKey(name string) tsValue { return tsStringReference(tsStringUTF8(name)) }
func tsNodeClusterObject(loop *tsLoop, values map[string]tsValue) tsValue { object:=tsNewObject();for key,value:=range values{object.set(key,value)};return tsObjectValue(object) }
func tsNodeClusterProcess(loop *tsLoop) tsValue { return tsNodeProcessModule(loop) }
func tsNodeClusterHandleCall(loop *tsLoop, handle tsValue, name string, args ...tsValue) tsValue {
	return tsCallReceiver(loop, tsGet(loop, handle, tsNodeClusterKey(name)), handle, args...)
}
func tsNodeClusterSend(loop *tsLoop, worker, message, handle, callback tsValue) tsValue {
	process := tsGet(loop, worker, tsNodeClusterKey("process"))
	send := tsGet(loop, process, tsNodeClusterKey("send"))
	if send.kind != tsFunctionKind { return tsBooleanValue(false) }
	return tsCallReceiver(loop, send, process, message, handle, tsU, callback)
}
func tsNodeClusterRemoveHandlesForWorker(loop *tsLoop, worker tsValue, channelGone bool) {
	for key, handle := range tsNodeClusterSingleton.handles {
		removed := tsNodeClusterHandleCall(loop, handle, "remove", worker, tsBooleanValue(channelGone))
		if tsTruthy(removed) { delete(tsNodeClusterSingleton.handles, key) }
	}
}
func tsNodeClusterRemoveWorker(loop *tsLoop, worker tsValue) {
	if tsNullish(worker) { tsPropertyFailure("ERR_INTERNAL_ASSERTION") }
	id := int(tsNumber(tsGet(loop, worker, tsNodeClusterKey("id"))))
	delete((*tsObject)(tsNodeClusterSingleton.workers.ref).values, strconv.Itoa(id))
	if len((*tsObject)(tsNodeClusterSingleton.workers.ref).values) == 0 {
		if len(tsNodeClusterSingleton.handles) != 0 { tsPropertyFailure("Resource leak detected.") }
		tsNodeEventEmit(loop, tsNodeEventState(tsNodeClusterSingleton.intercom), tsNodeClusterKey("disconnect"), nil)
	}
}
func tsNodeClusterSendReply(loop *tsLoop, worker, message, handle tsValue) tsValue {
	return tsNodeClusterSend(loop, worker, message, handle, tsU)
}
func tsNodeClusterQueryServer(loop *tsLoop, worker, message tsValue) {
	if tsTruthy(tsGet(loop, worker, tsNodeClusterKey("exitedAfterDisconnect"))) { return }
	address, port := tsGet(loop, message, tsNodeClusterKey("address")), tsGet(loop, message, tsNodeClusterKey("port"))
	addressType, fd := tsGet(loop, message, tsNodeClusterKey("addressType")), tsGet(loop, message, tsNodeClusterKey("fd"))
	if tsIsNumeric(port) && tsNumber(port)<0 && address.kind==tsStringKind && runtime.GOOS!="windows" {
		process:=tsNodeClusterProcess(loop);cwd:=tsCallReceiver(loop,tsGet(loop,process,tsNodeClusterKey("cwd")),process);path:=tsPathModule();relative:=tsCallReceiver(loop,tsGet(loop,path,tsNodeClusterKey("relative")),path,cwd,address);if len(tsText(relative))<len(tsText(address)){address=relative}
	}
	key := tsText(address)+":"+tsText(port)+":"+tsText(addressType)+":"+tsText(fd)
	if tsIsNumeric(port) && tsNumber(port)==0 { key += ":"+tsText(tsGet(loop,message,tsNodeClusterKey("index"))) }
	cachedHandle, hadCached := tsNodeClusterSingleton.handles[key]
	handle, reused := tsU, false
	if hadCached && !tsTruthy(tsNodeClusterHandleCall(loop,cachedHandle,"has",worker)) { handle, reused = cachedHandle, true }
	sharedOnly:=tsTruthy(tsGet(loop,message,tsNodeClusterKey("sharedOnly")))
	udp:=strings.Contains(tsText(addressType),"udp")
	hint:="TLS and non-TLS cluster workers cannot share the same address:port under SCHED_RR (Go TLS sockets do not transfer round-robin descriptors)"
	if reused && sharedOnly && tsInternalClusterRRFor(handle)!=nil { failure:=tsNewObject();failure.set("errno",tsNumberValue(-22));failure.set("key",tsStringReference(tsStringUTF8(key)));failure.set("ack",tsGet(loop,message,tsNodeClusterKey("seq")));failure.set("data",tsGet(loop,handle,tsNodeClusterKey("data")));failure.set("bunHint",tsStringReference(tsStringUTF8(hint)));tsNodeClusterSendReply(loop,worker,tsObjectValue(failure),tsNull);return }
	if reused && tsNodeClusterSingleton.schedulingPolicy==2 && !sharedOnly && !udp { if state:=tsInternalClusterSharedFor(handle);state!=nil&&state.sharedOnly { failure:=tsNewObject();failure.set("errno",tsNumberValue(-22));failure.set("key",tsStringReference(tsStringUTF8(key)));failure.set("ack",tsGet(loop,message,tsNodeClusterKey("seq")));failure.set("data",tsGet(loop,handle,tsNodeClusterKey("data")));failure.set("bunHint",tsStringReference(tsStringUTF8(hint)));tsNodeClusterSendReply(loop,worker,tsObjectValue(failure),tsNull);return } }
	if tsIsUndefined(handle) || tsIsNull(handle) {
		if runtime.GOOS=="windows" && udp { error:=tsNodeCPError("Error","ENOTSUP","cannot share a dgram socket with a worker on Windows");tsNodeClusterHandleCall(loop,worker,"emit",tsNodeClusterKey("error"),error);return }
		if tsNodeClusterSingleton.schedulingPolicy!=2 || sharedOnly || udp { handle=tsConstruct(loop,tsInternalClusterSharedHandleModule(loop),tsStringReference(tsStringUTF8(key)),address,message) } else { handle=tsConstruct(loop,tsInternalClusterRoundRobinModule(loop),tsStringReference(tsStringUTF8(key)),address,message) }
		if !hadCached { tsNodeClusterSingleton.handles[key]=handle }
	}
	if !tsTruthy(tsGet(loop,handle,tsNodeClusterKey("data"))) { tsSet(loop,handle,tsNodeClusterKey("data"),tsGet(loop,message,tsNodeClusterKey("data"))) }
	callback:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{
		errno,reply,serverHandle:=tsArg(args,0),tsArg(args,1),tsArg(args,2)
		data:=tsGet(loop,handle,tsNodeClusterKey("data"));out:=tsNewObject();out.set("errno",errno);out.set("key",tsStringReference(tsStringUTF8(key)));out.set("ack",tsGet(loop,message,tsNodeClusterKey("seq")));out.set("data",data)
		if reply.kind==tsObjectKind{for _,name:=range (*tsObject)(reply.ref).order{out.set(name,(*tsObject)(reply.ref).values[name])}}
		result:=tsNodeClusterSend(loop,worker,tsObjectValue(out),serverHandle,tsU)
		if !tsTruthy(result)&&!tsNullish(serverHandle){failure:=tsNewObject();failure.set("errno",tsNumberValue(-105));failure.set("key",tsStringReference(tsStringUTF8(key)));failure.set("ack",tsGet(loop,message,tsNodeClusterKey("seq")));failure.set("data",data);tsNodeClusterSendReply(loop,worker,tsObjectValue(failure),tsNull);if tsTruthy(tsNodeClusterHandleCall(loop,handle,"remove",worker)){delete(tsNodeClusterSingleton.handles,key)}}
		if hadCached && !reused && !tsTruthy(errno) { tsNodeClusterHandleCall(loop,handle,"remove",worker) }
		return tsU
	}))
	tsNodeClusterHandleCall(loop,handle,"add",worker,callback)
}
func tsNodeClusterOnMessage(loop *tsLoop, worker, message tsValue) {
	switch tsText(tsGet(loop,message,tsNodeClusterKey("act"))) {
	case "online":
		tsSet(loop,worker,tsNodeClusterKey("state"),tsNodeClusterKey("online"));tsNodeClusterHandleCall(loop,worker,"emit",tsNodeClusterKey("online"));tsCallReceiver(loop,tsGet(loop,tsNodeClusterSingleton.module,tsNodeClusterKey("emit")),tsNodeClusterSingleton.module,tsNodeClusterKey("online"),worker)
	case "exitedAfterDisconnect":
		tsSet(loop,worker,tsNodeClusterKey("exitedAfterDisconnect"),tsBooleanValue(true));ack:=tsNewObject();ack.set("ack",tsGet(loop,message,tsNodeClusterKey("seq")));tsNodeClusterSendReply(loop,worker,tsObjectValue(ack),tsNull)
	case "queryServer": tsNodeClusterQueryServer(loop,worker,message)
	case "listening":
		info:=tsNodeClusterObject(loop,map[string]tsValue{"addressType":tsGet(loop,message,tsNodeClusterKey("addressType")),"address":tsGet(loop,message,tsNodeClusterKey("address")),"port":tsGet(loop,message,tsNodeClusterKey("port")),"fd":tsGet(loop,message,tsNodeClusterKey("fd"))})
		tsSet(loop,worker,tsNodeClusterKey("state"),tsNodeClusterKey("listening"));tsNodeClusterHandleCall(loop,worker,"emit",tsNodeClusterKey("listening"),info);tsCallReceiver(loop,tsGet(loop,tsNodeClusterSingleton.module,tsNodeClusterKey("emit")),tsNodeClusterSingleton.module,tsNodeClusterKey("listening"),worker,info)
	case "close":
		key:=tsText(tsGet(loop,message,tsNodeClusterKey("key")));if handle,ok:=tsNodeClusterSingleton.handles[key];ok&&tsTruthy(tsNodeClusterHandleCall(loop,handle,"remove",worker)){delete(tsNodeClusterSingleton.handles,key)}
	}
}
func tsNodeClusterNextTick(loop *tsLoop, callback tsValue, args ...tsValue) {
	process:=tsNodeClusterProcess(loop);nextTick:=tsGet(loop,process,tsNodeClusterKey("nextTick"));if nextTick.kind==tsFunctionKind{tsCallReceiver(loop,nextTick,process,append([]tsValue{callback},args...)...)}else{loop.ready=append(loop.ready,func(){tsCall(loop,callback,args...)})}
}
func tsNodeClusterMergeSettings(loop *tsLoop, base, override tsValue) tsValue {
	merged:=tsNewObject()
	for _,source:=range []tsValue{base,override}{if source.kind!=tsObjectKind&&source.kind!=tsInstanceKind{continue};object:=(*tsObject)(source.ref);for _,key:=range object.order{merged.set(key,object.values[key])}}
	return tsObjectValue(merged)
}
func tsNodeClusterSetupPrimary(loop *tsLoop, options tsValue) tsValue {
	process:=tsNodeClusterProcess(loop);argv:=tsGet(loop,process,tsNodeClusterKey("argv"));execArgv:=tsGet(loop,process,tsNodeClusterKey("execArgv"));args:=&tsArray{}
	if argv.kind==tsArrayKind{source:=(*tsArray)(argv.ref);for index:=2;index<source.length();index++{args.values=append(args.values,source.at(index))}}
	defaults:=tsNodeClusterObject(loop,map[string]tsValue{"args":tsArrayValue(args),"exec":tsArg((*tsArray)(argv.ref).values,1),"execArgv":execArgv,"silent":tsBooleanValue(false)})
	tsNodeClusterSingleton.settings=tsNodeClusterMergeSettings(loop,defaults,tsNodeClusterSingleton.settings)
	tsNodeClusterSingleton.settings=tsNodeClusterMergeSettings(loop,tsNodeClusterSingleton.settings,options)
	if tsNodeClusterSingleton.initialized{tsNodeClusterNextTick(loop,tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{tsNodeEventEmit(loop,tsNodeEventState(tsNodeClusterSingleton.module),tsNodeClusterKey("setup"),[]tsValue{tsNodeClusterSingleton.settings});return tsU})),tsNodeClusterSingleton.settings);return tsU}
	tsNodeClusterSingleton.initialized=true
	policy:=tsGet(loop,tsNodeClusterSingleton.module,tsNodeClusterKey("schedulingPolicy"));if tsIsUndefined(policy){policy=tsNumberValue(2)};tsNodeClusterSingleton.schedulingPolicy=int(tsNumber(policy))
	callback:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{tsNodeEventEmit(loop,tsNodeEventState(tsNodeClusterSingleton.module),tsNodeClusterKey("setup"),[]tsValue{tsNodeClusterSingleton.settings});return tsU}));tsNodeClusterNextTick(loop,callback,tsNodeClusterSingleton.settings);return tsU
}
func tsNodeClusterCreateWorkerProcess(loop *tsLoop,id int,environment tsValue) tsValue {
	process:=tsNodeClusterProcess(loop);env:=tsNodeClusterMergeSettings(loop,tsGet(loop,process,tsNodeClusterKey("env")),environment);envObject:=(*tsObject)(env.ref);envObject.set("NODE_UNIQUE_ID",tsNodeClusterKey(strconv.Itoa(id)))
	options:=tsNewObject();for _,name:=range []string{"cwd","serialization","silent","windowsHide","stdio","gid","uid"}{value:=tsGet(loop,tsNodeClusterSingleton.settings,tsNodeClusterKey(name));if !tsIsUndefined(value){options.set(name,value)}}
	options.set("env",env);options.set("execArgv",tsGet(loop,tsNodeClusterSingleton.settings,tsNodeClusterKey("execArgv")))
	childProcess:=tsNodeChildProcessModule(loop);fork:=tsGet(loop,childProcess,tsNodeClusterKey("fork"));if fork.kind!=tsFunctionKind{panic(tsThrown{tsFSException("Error","ERR_METHOD_NOT_IMPLEMENTED","child_process.fork is unavailable")})}
	exec:=tsGet(loop,tsNodeClusterSingleton.settings,tsNodeClusterKey("exec"));args:=tsGet(loop,tsNodeClusterSingleton.settings,tsNodeClusterKey("args"));return tsCallReceiver(loop,fork,childProcess,exec,args,options)
}
func tsNodeClusterFork(loop *tsLoop,environment tsValue) tsValue {
	tsNodeClusterSetupPrimary(loop,tsU);tsNodeClusterSingleton.ids++;id:=tsNodeClusterSingleton.ids;child:=tsNodeClusterCreateWorkerProcess(loop,id,environment)
	options:=tsNodeClusterObject(loop,map[string]tsValue{"id":tsNumberValue(float64(id)),"process":child})
	worker:=tsConstruct(loop,tsNodeClusterSingleton.workerConstructor,options)
	workerState:=tsInternalClusterWorkerFor(worker)
	if workerState!=nil{
		messageListener:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsCallReceiver(loop,tsGet(loop,tsNodeClusterSingleton.module,tsNodeClusterKey("emit")),tsNodeClusterSingleton.module,append([]tsValue{tsNodeClusterKey("message"),worker},args...)...)}))
		on:=tsGet(loop,worker,tsNodeClusterKey("on"));tsCallReceiver(loop,on,worker,tsNodeClusterKey("message"),messageListener)
		protocolListener:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{tsNodeClusterOnMessage(loop,worker,tsArg(args,0));return tsU}));if procOn:=tsGet(loop,child,tsNodeClusterKey("on"));procOn.kind==tsFunctionKind{tsCallReceiver(loop,procOn,child,tsNodeClusterKey("message"),protocolListener)}
		for _,eventName:=range []string{"exit","disconnect"}{eventName:=eventName;listener:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{if eventName=="exit"{connected:=tsCallReceiver(loop,tsGet(loop,worker,tsNodeClusterKey("isConnected")),worker);if !tsTruthy(connected){tsNodeClusterRemoveHandlesForWorker(loop,worker,true);tsNodeClusterRemoveWorker(loop,worker)};tsSet(loop,worker,tsNodeClusterKey("state"),tsNodeClusterKey("dead"));tsCallReceiver(loop,tsGet(loop,worker,tsNodeClusterKey("emit")),worker,append([]tsValue{tsNodeClusterKey("exit")},args...)...);tsCallReceiver(loop,tsGet(loop,tsNodeClusterSingleton.module,tsNodeClusterKey("emit")),tsNodeClusterSingleton.module,append([]tsValue{tsNodeClusterKey("exit"),worker},args...)...)}else{tsNodeClusterRemoveHandlesForWorker(loop,worker,true);dead:=tsCallReceiver(loop,tsGet(loop,worker,tsNodeClusterKey("isDead")),worker);if tsTruthy(dead){tsNodeClusterRemoveWorker(loop,worker)};tsCallReceiver(loop,tsGet(loop,worker,tsNodeClusterKey("emit")),worker,tsNodeClusterKey("disconnect"));tsCallReceiver(loop,tsGet(loop,tsNodeClusterSingleton.module,tsNodeClusterKey("emit")),tsNodeClusterSingleton.module,tsNodeClusterKey("disconnect"),worker)};return tsU}))
			procOn:=tsGet(loop,child,tsNodeClusterKey("once"));if procOn.kind==tsFunctionKind{tsCallReceiver(loop,procOn,child,tsNodeClusterKey(eventName),listener)}
		}
	}
	(*tsObject)(tsNodeClusterSingleton.workers.ref).set(strconv.Itoa(id),worker)
	callback:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{tsCallReceiver(loop,tsGet(loop,tsNodeClusterSingleton.module,tsNodeClusterKey("emit")),tsNodeClusterSingleton.module,tsNodeClusterKey("fork"),worker);return tsU}));tsNodeClusterNextTick(loop,callback)
	return worker
}
func tsNodeClusterDisconnect(loop *tsLoop, callback tsValue) tsValue {
	workers:=(*tsObject)(tsNodeClusterSingleton.workers.ref);if len(workers.values)==0{tsNodeClusterNextTick(loop,tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{tsNodeEventEmit(loop,tsNodeEventState(tsNodeClusterSingleton.intercom),tsNodeClusterKey("disconnect"),nil);return tsU}))) }else{for _,worker:=range workers.values{if tsTruthy(tsCallReceiver(loop,tsGet(loop,worker,tsNodeClusterKey("isConnected")),worker)){tsCallReceiver(loop,tsGet(loop,worker,tsNodeClusterKey("disconnect")),worker)}}}
	if callback.kind==tsFunctionKind{tsCallReceiver(loop,tsGet(loop,tsNodeClusterSingleton.intercom,tsNodeClusterKey("once")),tsNodeClusterSingleton.intercom,tsNodeClusterKey("disconnect"),callback)};return tsU
}
func tsNodeClusterModule(loop *tsLoop) tsValue {
	if !tsTruthy(tsGet(loop, tsInternalClusterIsPrimaryModule(loop), tsNodeClusterKey("isPrimary"))) { return tsInternalClusterChildModule(loop) }
	tsNodeClusterSingleton.Do(func(){
		tsNodeEventsInit(loop);cluster:=tsNodeEventCreate(loop,tsU,tsU);tsNodeClusterSingleton.module=cluster;tsNodeClusterSingleton.workers=tsObjectValue(tsNewObject());tsNodeClusterSingleton.settings=tsObjectValue(tsNewObject());tsNodeClusterSingleton.handles=map[string]tsValue{};tsNodeClusterSingleton.intercom=tsNodeEventCreate(loop,tsU,tsU);tsNodeClusterSingleton.workerConstructor=tsInternalClusterWorkerModule(loop)
		clusterObject:=(*tsObject)(cluster.ref);clusterObject.set("isWorker",tsBooleanValue(false));clusterObject.set("isMaster",tsBooleanValue(true));clusterObject.set("isPrimary",tsBooleanValue(true));clusterObject.set("Worker",tsNodeClusterSingleton.workerConstructor);clusterObject.set("workers",tsNodeClusterSingleton.workers);clusterObject.set("settings",tsNodeClusterSingleton.settings);clusterObject.set("SCHED_NONE",tsNumberValue(1));clusterObject.set("SCHED_RR",tsNumberValue(2))
		workerPrototype:=tsInternalClusterWorker.prototype
		workerPrototype.set("disconnect",tsNamedNativeMethod("disconnect",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsInternalClusterWorkerFor(receiver);if state==nil{return receiver};state.object.set("exitedAfterDisconnect",tsBooleanValue(true));message:=tsNewObject();message.set("act",tsNodeClusterKey("disconnect"));tsNodeClusterSend(loop,receiver,tsObjectValue(message),tsNull,tsU);process:=state.process;disconnect:=tsGet(loop,process,tsNodeClusterKey("disconnect"));if disconnect.kind==tsFunctionKind{tsCallReceiver(loop,disconnect,process)};tsNodeClusterRemoveHandlesForWorker(loop,receiver,false);tsNodeClusterRemoveWorker(loop,receiver);return receiver}))
		workerPrototype.set("destroy",tsNamedNativeMethod("destroy",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsInternalClusterWorkerFor(receiver);if state==nil{return receiver};signal:=tsArg(args,0);if tsIsUndefined(signal)||tsIsNull(signal)||tsText(signal)==""{signal=tsNodeClusterKey("SIGTERM")};kill:=tsGet(loop,state.process,tsNodeClusterKey("kill"));if kill.kind==tsFunctionKind{tsCallReceiver(loop,kill,state.process,signal)};return receiver}))
		process:=tsNodeClusterProcess(loop);policy:=tsGet(loop,tsGet(loop,process,tsNodeClusterKey("env")),tsNodeClusterKey("NODE_CLUSTER_SCHED_POLICY"));selected:=2;if tsText(policy)=="none"{selected=1};clusterObject.set("schedulingPolicy",tsNumberValue(float64(selected)))
		setup:=tsNamedNativeMethod("setupPrimary",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeClusterSetupPrimary(loop,tsArg(args,0))});clusterObject.set("setupPrimary",setup);clusterObject.set("setupMaster",setup)
		clusterObject.set("fork",tsNamedNativeMethod("fork",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeClusterFork(loop,tsArg(args,0))}))
		clusterObject.set("disconnect",tsNamedNativeMethod("disconnect",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeClusterDisconnect(loop,tsArg(args,0))}))
		clusterObject.set("default",cluster)
	})
	return tsNodeClusterSingleton.module
}
`
