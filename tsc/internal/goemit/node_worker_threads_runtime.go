package goemit

// nodeWorkerThreadsRuntime exposes the Node worker_threads API surface. The
// interpreter has no second JS VM yet, so Worker is a lifecycle shell; message
// channels and the rest of the per-thread helpers run on the current loop.
const nodeWorkerThreadsRuntime = `
type tsNodeWorkerPortState struct {
	object *tsObject
	peer *tsNodeWorkerPortState
	queue []tsValue
	dom map[string][]*tsNodeWorkerDOMListener
	onmessage, onmessageerror tsValue
	name string
	closed, started, refed, scheduled bool
}
type tsNodeWorkerDOMListener struct { original, wrapper tsValue; once bool }
type tsNodeWorkerState struct { object *tsObject; id int; name string; exited, refed bool; stdin, stdout, stderr, workerData tsValue }
type tsNodeBroadcastGroup struct { ports []*tsNodeWorkerPortState }
type tsNodeWorkerEnvironmentEntry struct { key, value tsValue }
var tsNodeWorkerThreads struct {
	sync.Once
	module, workerConstructor, workerPrototype, portConstructor, portPrototype, channelConstructor, broadcastConstructor, broadcastPrototype tsValue
	ports map[*tsObject]*tsNodeWorkerPortState
	workers map[*tsObject]*tsNodeWorkerState
	channels map[string]*tsNodeBroadcastGroup
	environment []tsNodeWorkerEnvironmentEntry
	markedUntransferable map[uintptr]bool
	markedUncloneable map[uintptr]bool
	shareEnv tsValue
	nextWorkerID int
	threadID int
	isMainThread bool
	threadName tsValue
}

func tsNodeWorkerKey(name string) tsValue { return tsStringReference(tsStringUTF8(name)) }
func tsNodeWorkerMarkable(value tsValue)bool{switch value.kind{case tsObjectKind,tsInstanceKind,tsArrayKind,tsTypedArrayKind,tsDataViewKind,tsArrayBufferKind,tsCollectionKind,tsRegExpKind,tsErrorKind:return true};return false}
func tsNodeWorkerDataCloneFailure(message string){panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"DataCloneError",message:message,fields:map[string]tsValue{"code":tsNumberValue(25)}})})}
func tsNodeWorkerCheckCloneable(loop *tsLoop,value tsValue){pending:=[]tsValue{value};seen:=map[unsafe.Pointer]bool{};for len(pending)>0{current:=pending[len(pending)-1];pending=pending[:len(pending)-1];if current.ref!=nil{if tsNodeWorkerThreads.markedUncloneable[uintptr(current.ref)]{tsNodeWorkerDataCloneFailure("Cannot clone an object marked as uncloneable")};if seen[current.ref]{continue};seen[current.ref]=true};if current.kind==tsTypedArrayKind{pending=append(pending,tsArrayBufferValue((*tsTypedArray)(current.ref).buffer))}else if current.kind==tsDataViewKind{pending=append(pending,tsArrayBufferValue((*tsDataView)(current.ref).buffer))};switch current.kind{case tsArrayKind:array:=(*tsArray)(current.ref);for index,item:=range array.values{if !array.holes[index]{pending=append(pending,item)}};for _,name:=range array.order{if descriptor:=array.descriptors[name];descriptor!=nil&&(descriptor.accessor||!descriptor.enumerable){continue};if item,ok:=array.properties[name];ok{pending=append(pending,item)}};case tsCollectionKind:collection:=(*tsCollection)(current.ref);for _,entry:=range collection.entries{if entry.live{pending=append(pending,entry.key);if !collection.set{pending=append(pending,entry.value)}}};case tsObjectKind,tsInstanceKind,tsErrorKind,tsTypedArrayKind,tsDataViewKind:for _,key:=range tsEnumerableOwnKeys(loop,current).values{descriptor:=tsOwnDescriptor(loop,current,tsText(key));if descriptor.kind==tsObjectKind{details:=(*tsObject)(descriptor.ref);if _,accessor:=details.values["get"];accessor{continue};pending=append(pending,details.values["value"])}}}}}
func tsNodeWorkerCloneValue(loop *tsLoop,value tsValue)tsValue{tsNodeWorkerCheckCloneable(loop,value);return tsNodeV8Decode(loop,tsNodeV8Encode(loop,value))}
func tsNodeWorkerCloneValueWithTransfers(loop *tsLoop,value,transferList tsValue)tsValue{
	if tsNullish(transferList){return tsNodeWorkerCloneValue(loop,value)}
	if transferList.kind!=tsArrayKind{tsInternalValidatorInvalidType("transferList","an array",transferList)}
	type portTransfer struct{ Source, Destination *tsNodeWorkerPortState }
	array:=(*tsArray)(transferList.ref);transferIDs:=map[unsafe.Pointer]uint32{};transfers:=map[uint32]tsValue{};buffers:=[]*tsArrayBuffer{};ports:=[]portTransfer{};seen:=map[unsafe.Pointer]bool{};created:=[]*tsObject{};committed:=false
	defer func(){if !committed{for _,object:=range created{delete(tsNodeWorkerThreads.ports,object)}}}()
	nextID:=uint32(0)
	for index:=0;index<array.length();index++{item:=array.at(index);key:=item.ref;if key==nil||seen[key]{tsNodeWorkerDataCloneFailure("Transfer list contains a duplicate or invalid entry")};seen[key]=true
		if item.kind==tsArrayBufferKind{buffer:=(*tsArrayBuffer)(key);if buffer.detached{tsNodeWorkerDataCloneFailure("Cannot transfer a detached ArrayBuffer")};if tsNodeWorkerThreads.markedUntransferable[uintptr(key)]{continue};id:=nextID;nextID++;transferIDs[key]=id;copyBuffer:=&tsArrayBuffer{data:append([]byte(nil),buffer.data...),maxByteLength:buffer.maxByteLength,resizable:buffer.resizable};transfers[id]=tsArrayBufferValue(copyBuffer);buffers=append(buffers,buffer);continue}
		port:=tsNodeWorkerPortFor(item);if port==nil||port.closed||port.name!=""{tsNodeWorkerDataCloneFailure("Found invalid object in transferList")};if tsNodeWorkerThreads.markedUntransferable[uintptr(key)]{tsNodeWorkerDataCloneFailure("MessagePort is marked as untransferable")};id:=nextID;nextID++;transferIDs[key]=id;destinationValue:=tsNodeWorkerPortCreate(loop,"");destination:=tsNodeWorkerPortFor(destinationValue);created=append(created,destination.object);transfers[id]=destinationValue;ports=append(ports,portTransfer{Source:port,Destination:destination})
	}
	tsNodeWorkerCheckCloneable(loop,value);payload:=tsNodeV8EncodeWithTransfers(loop,value,transferIDs);cloned:=tsNodeV8DecodeWithTransfers(loop,payload,transfers)
	for _,port:=range ports{source,destination:=port.Source,port.Destination;destination.peer=source.peer;destination.queue=source.queue;destination.started=source.started;destination.refed=source.refed;if source.peer!=nil{source.peer.peer=destination};source.peer=nil;source.queue=nil;source.closed=true;source.started=false;source.scheduled=false;tsNodeWorkerPortSchedule(loop,destination)}
	for _,buffer:=range buffers{tsArrayBufferTransfer(buffer,len(buffer.data))}
	committed=true
	return cloned
}
func tsNodeWorkerGetEnvironmentData(loop *tsLoop,key tsValue)tsValue{for index:=len(tsNodeWorkerThreads.environment)-1;index>=0;index--{entry:=tsNodeWorkerThreads.environment[index];if tsSameValueZero(entry.key,key){return tsNodeWorkerCloneValue(loop,entry.value)}};return tsU}
func tsNodeWorkerSetEnvironmentData(loop *tsLoop,key,value tsValue){for index:=range tsNodeWorkerThreads.environment{if tsSameValueZero(tsNodeWorkerThreads.environment[index].key,key){if tsIsUndefined(value){tsNodeWorkerThreads.environment=append(tsNodeWorkerThreads.environment[:index],tsNodeWorkerThreads.environment[index+1:]...)}else{tsNodeWorkerThreads.environment[index].value=tsNodeWorkerCloneValue(loop,value)};return}};if !tsIsUndefined(value){tsNodeWorkerThreads.environment=append(tsNodeWorkerThreads.environment,tsNodeWorkerEnvironmentEntry{key:key,value:tsNodeWorkerCloneValue(loop,value)})}}
func tsNodeWorkerObject(loop *tsLoop,values map[string]tsValue) tsValue { object:=tsNewObject();for name,value:=range values{object.set(name,value)};return tsObjectValue(object) }
func tsNodeWorkerError(name,code,message string) tsValue { return tsErrorValue(&tsRuntimeError{name:name,message:message,fields:map[string]tsValue{"code":tsNodeWorkerKey(code)}}) }
func tsNodeWorkerPortFor(receiver tsValue) *tsNodeWorkerPortState { if receiver.kind!=tsObjectKind{return nil};return tsNodeWorkerThreads.ports[(*tsObject)(receiver.ref)] }
func tsNodeWorkerFor(receiver tsValue) *tsNodeWorkerState { if receiver.kind!=tsObjectKind{return nil};return tsNodeWorkerThreads.workers[(*tsObject)(receiver.ref)] }
func tsNodeWorkerPortCreate(loop *tsLoop,name string) tsValue {
	tsNodeEventsInit(loop)
	value:=tsNodeEventCreate(loop,tsU,tsU);object:=(*tsObject)(value.ref);state:=&tsNodeWorkerPortState{object:object,dom:map[string][]*tsNodeWorkerDOMListener{},onmessage:tsU,onmessageerror:tsU,name:name,refed:true}
	if tsNodeWorkerThreads.portPrototype.kind==tsObjectKind{object.prototype=tsNodeWorkerThreads.portPrototype}
	tsNodeWorkerThreads.ports[object]=state
	object.set("postMessage",tsNamedNativeMethod("postMessage",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{
		port:=tsNodeWorkerPortFor(receiver);if port==nil{tsPropertyFailure("Illegal invocation")};if port.closed{panic(tsThrown{tsNodeWorkerError("Error","ERR_CLOSED_MESSAGE_PORT","Cannot send data on a closed MessagePort")})}
		data:=tsNodeWorkerCloneValueWithTransfers(loop,tsArg(args,0),func()tsValue{if port.name!=""{return tsU};return tsArg(args,1)}());if port.name!=""{group:=tsNodeWorkerThreads.channels[port.name];if group!=nil{for _,other:=range group.ports{if other!=port&&!other.closed{other.queue=append(other.queue,tsNodeWorkerCloneValue(loop,data));tsNodeWorkerPortSchedule(loop,other)}}}}else if port.peer!=nil&&!port.peer.closed{port.peer.queue=append(port.peer.queue,data);tsNodeWorkerPortSchedule(loop,port.peer)}
		return tsU
	}))
	object.set("start",tsNamedNativeMethod("start",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{if port:=tsNodeWorkerPortFor(receiver);port!=nil{port.started=true;tsNodeWorkerPortSchedule(loop,port)};return tsU}))
	object.set("close",tsNamedNativeMethod("close",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{
		port:=tsNodeWorkerPortFor(receiver);if port==nil||port.closed{return tsU};port.closed=true;port.queue=nil
		if callback:=tsArg(args,0);callback.kind==tsFunctionKind{tsNodeEventAdd(loop,receiver,tsNodeWorkerKey("close"),callback,false,true)}
		tsNodeEventEmit(loop,tsNodeEventState(receiver),tsNodeWorkerKey("close"),nil);return tsU
	}))
	object.set("ref",tsNamedNativeMethod("ref",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{if port:=tsNodeWorkerPortFor(receiver);port!=nil{port.refed=true};return receiver}))
	object.set("unref",tsNamedNativeMethod("unref",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{if port:=tsNodeWorkerPortFor(receiver);port!=nil{port.refed=false};return receiver}))
	object.set("hasRef",tsNamedNativeMethod("hasRef",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{port:=tsNodeWorkerPortFor(receiver);return tsBooleanValue(port!=nil&&port.refed&&!port.closed)}))
	object.set("addEventListener",tsNamedNativeMethod("addEventListener",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeWorkerPortAddDOM(loop,receiver,tsArg(args,0),tsArg(args,1),tsArg(args,2))}))
	object.set("removeEventListener",tsNamedNativeMethod("removeEventListener",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeWorkerPortRemoveDOM(loop,receiver,tsArg(args,0),tsArg(args,1))}))
	object.set("dispatchEvent",tsNamedNativeMethod("dispatchEvent",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeWorkerPortDispatch(loop,receiver,tsArg(args,0))}))
	return value
}
func tsNodeWorkerPortEvent(loop *tsLoop,port *tsNodeWorkerPortState,name string,data tsValue) tsValue {
	if port==nil||port.closed{return tsU};target:=tsObjectValue(port.object);var event tsValue=tsNodeWorkerObject(loop,map[string]tsValue{"type":tsNodeWorkerKey(name),"data":data,"target":target,"currentTarget":target})
	handler:=port.onmessage;if name=="messageerror"{handler=port.onmessageerror};if handler.kind==tsFunctionKind{tsCallReceiver(loop,handler,target,event)}
	if emitter:=tsNodeEventState(target);emitter!=nil{tsNodeEventEmit(loop,emitter,tsNodeWorkerKey(name),[]tsValue{data})}
	return tsU
}
func tsNodeWorkerPortSchedule(loop *tsLoop,port *tsNodeWorkerPortState) {
	if port==nil||port.scheduled||!port.started||port.closed||len(port.queue)==0{return};port.scheduled=true
	loop.ready=append(loop.ready,func(){port.scheduled=false;if port.closed||!port.started||len(port.queue)==0{return};data:=port.queue[0];port.queue=port.queue[1:];tsNodeWorkerPortEvent(loop,port,"message",data);tsNodeWorkerPortSchedule(loop,port)})
}
func tsNodeWorkerPortAddDOM(loop *tsLoop,receiver,event,listener,options tsValue) tsValue {
	port:=tsNodeWorkerPortFor(receiver);if port==nil{tsPropertyFailure("Illegal invocation")};if listener.kind!=tsFunctionKind{return tsU};name:=tsText(event);once:=tsTruthy(tsGet(loop,options,tsNodeWorkerKey("once")))
	var wrapper tsValue
	entry:=&tsNodeWorkerDOMListener{original:listener,once:once}
	wrapper=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{if once{tsNodeWorkerPortRemoveDOM(loop,receiver,event,listener)};data:=tsArg(args,0);target:=tsObjectValue(port.object);message:=tsNodeWorkerObject(loop,map[string]tsValue{"type":event,"data":data,"target":target,"currentTarget":target});return tsCallReceiver(loop,listener,target,message)}));entry.wrapper=wrapper;port.dom[name]=append(port.dom[name],entry)
	if name=="message"||name=="messageerror"{port.started=true;tsNodeWorkerPortSchedule(loop,port)}
	tsNodeEventAdd(loop,receiver,event,wrapper,false,false);return tsU
}
func tsNodeWorkerPortRemoveDOM(loop *tsLoop,receiver,event,listener tsValue) tsValue {
	port:=tsNodeWorkerPortFor(receiver);if port==nil{return tsU};name:=tsText(event);entries:=port.dom[name]
	for index:=len(entries)-1;index>=0;index--{if tsStrictEqual(entries[index].original,listener){tsNodeEventRemove(loop,tsNodeEventState(receiver),event,entries[index].wrapper);entries=append(entries[:index],entries[index+1:]...);break}}
	port.dom[name]=entries;return tsU
}
func tsNodeWorkerPortDispatch(loop *tsLoop,receiver,event tsValue) tsValue {name:=tsText(tsGet(loop,event,tsNodeWorkerKey("type")));data:=tsGet(loop,event,tsNodeWorkerKey("data"));port:=tsNodeWorkerPortFor(receiver);if port==nil{return tsBooleanValue(false)};tsNodeWorkerPortEvent(loop,port,name,data);return tsBooleanValue(true)}
func tsNodeWorkerReceive(loop *tsLoop,args []tsValue) tsValue {
	port:=tsNodeWorkerPortFor(tsArg(args,0));if port==nil{panic(tsThrown{tsNodeWorkerError("TypeError","ERR_INVALID_ARG_TYPE","The \"port\" argument must be a MessagePort")})};if len(port.queue)==0{return tsU};data:=port.queue[0];port.queue=port.queue[1:];return tsNodeWorkerObject(loop,map[string]tsValue{"message":data})
}
func tsNodeWorkerMakeChannel(loop *tsLoop) tsValue {left:=tsNodeWorkerPortCreate(loop,"");right:=tsNodeWorkerPortCreate(loop,"");a,b:=tsNodeWorkerPortFor(left),tsNodeWorkerPortFor(right);a.peer=b;b.peer=a;return tsNodeWorkerObject(loop,map[string]tsValue{"port1":left,"port2":right})}

func tsNodeWorkerCreate(loop *tsLoop,args []tsValue) tsValue {
	filename,options:=tsArg(args,0),tsArg(args,1);if tsIsUndefined(filename){panic(tsThrown{tsNodeWorkerError("TypeError","ERR_MISSING_ARGS","The \"filename\" argument is required")})}
	if filename.kind==tsObjectKind&&(*tsObject)(filename.ref).nodeURL!=nil||filename.kind==tsStringKind&&strings.HasPrefix(strings.ToLower(tsText(filename)),"file:"){filename=tsStringReference(tsStringUTF8(tsNodeURLAsFilePath(loop,filename,runtime.GOOS=="windows")))};if filename.kind!=tsStringKind{panic(tsThrown{tsNodeWorkerError("TypeError","ERR_INVALID_ARG_TYPE","The \"filename\" argument must be of type string or an instance of URL")})}
	if tsNullish(options){options=tsObjectValue(tsNewObject())};if options.kind!=tsObjectKind&&options.kind!=tsInstanceKind{panic(tsThrown{tsNodeWorkerError("TypeError","ERR_INVALID_ARG_TYPE","The \"options\" argument must be an object")})}
	eval:=tsTruthy(tsGet(loop,options,tsNodeWorkerKey("eval")));file:=tsText(filename);if !eval&&!filepath.IsAbs(file)&&!strings.HasPrefix(file,"./")&&!strings.HasPrefix(file,"../"){panic(tsThrown{tsNodeWorkerError("TypeError","ERR_WORKER_PATH","The worker script or module filename must be an absolute path or a relative path starting with './' or '../'. Received \""+file+"\"")})}
	name:=tsGet(loop,options,tsNodeWorkerKey("name"));if !tsIsUndefined(name)&&name.kind!=tsStringKind{panic(tsThrown{tsNodeWorkerError("TypeError","ERR_INVALID_ARG_TYPE","The \"options.name\" argument must be of type string")})}
	tsRuntimeUnsupported("node:worker_threads Worker()","the Go program cannot launch another JavaScript runtime or execute a dynamic worker source file")
	return tsU
}
func tsNodeWorkerThreadID(receiver tsValue) tsValue {state:=tsNodeWorkerFor(receiver);if state==nil{return tsNumberValue(0)};return tsNumberValue(float64(state.id))}
func tsNodeWorkerTerminate(loop *tsLoop,receiver tsValue) tsValue {state:=tsNodeWorkerFor(receiver);promise:=loop.promise();if state==nil||state.exited{promise.settle(tsResult{tsU,false});return tsPromiseValue(promise)};state.exited=true;state.object.set("threadId",tsNumberValue(-1));state.object.set("threadName",tsNull);state.object.set("stdin",tsNull);state.object.set("stdout",tsNull);state.object.set("stderr",tsNull);tsNodeEventEmit(loop,tsNodeEventState(receiver),tsNodeWorkerKey("exit"),[]tsValue{tsNumberValue(1)});promise.settle(tsResult{tsNumberValue(1),false});return tsPromiseValue(promise)}
func tsNodeWorkerSnapshot(loop *tsLoop) tsValue {module:=tsInternalStreamReadableModule(loop);constructor:=tsGet(loop,module,tsNodeWorkerKey("Readable"));text:="{}";sent:=false;var stream tsValue;options:=tsNewObject();options.set("read",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{if !sent{sent=true;tsCallReceiver(loop,tsGet(loop,stream,tsNodeWorkerKey("push")),stream,tsStringReference(tsStringUTF8(text)));tsCallReceiver(loop,tsGet(loop,stream,tsNodeWorkerKey("push")),stream,tsNull)};return tsU})));stream=tsConstruct(constructor,tsObjectValue(options));return stream}
func tsNodeWorkerValidateOptions(loop *tsLoop,options tsValue,name string) {if tsIsUndefined(options)||tsIsNull(options){return};if options.kind!=tsObjectKind&&options.kind!=tsInstanceKind{panic(tsThrown{tsNodeWorkerError("TypeError","ERR_INVALID_ARG_TYPE","The \""+name+"\" argument must be an object")})}}
func tsNodeWorkerPromise(loop *tsLoop,value tsValue,rejected bool) tsValue {return tsPromiseValue(loop.resolved(value,rejected))}
func tsNodeWorkerPostMessageToThread(loop *tsLoop,args []tsValue) tsValue {id:=tsArg(args,0);if tsStrictEqual(id,tsNumberValue(float64(tsNodeWorkerThreads.threadID))){return tsNodeWorkerPromise(loop,tsNodeWorkerError("Error","ERR_WORKER_MESSAGING_SAME_THREAD","Cannot send a message to the same thread."),true)};return tsNodeWorkerPromise(loop,tsNodeWorkerError("Error","ERR_WORKER_MESSAGING_FAILED","The destination thread no longer exists or is not listening for workerMessage events."),true)}

func tsNodeWorkerThreadsInit(loop *tsLoop) {
	tsNodeWorkerThreads.Do(func(){
		tsNodeEventsInit(loop);tsNodeWorkerThreads.ports=map[*tsObject]*tsNodeWorkerPortState{};tsNodeWorkerThreads.workers=map[*tsObject]*tsNodeWorkerState{};tsNodeWorkerThreads.channels=map[string]*tsNodeBroadcastGroup{};tsNodeWorkerThreads.environment=[]tsNodeWorkerEnvironmentEntry{};tsNodeWorkerThreads.markedUntransferable=map[uintptr]bool{};tsNodeWorkerThreads.markedUncloneable=map[uintptr]bool{};tsNodeWorkerThreads.threadID=0;tsNodeWorkerThreads.isMainThread=true;tsNodeWorkerThreads.threadName=tsNull
		eventProto:=tsNodeEventsSingleton.prototype;workerProto:=tsNewObject();workerProto.prototype=tsObjectValue(eventProto);tsNodeWorkerThreads.workerPrototype=tsObjectValue(workerProto)
		workerProto.set("ref",tsNamedNativeMethod("ref",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{if state:=tsNodeWorkerFor(receiver);state!=nil{state.refed=true};return receiver}));workerProto.set("unref",tsNamedNativeMethod("unref",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{if state:=tsNodeWorkerFor(receiver);state!=nil{state.refed=false};return receiver}))
		workerProto.set("postMessage",tsNamedNativeMethod("postMessage",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeWorkerFor(receiver);if state==nil||state.exited{return tsU};_ = tsNodeWorkerCloneValue(loop,tsArg(args,0));return tsU}))
		workerProto.set("terminate",tsNamedNativeMethod("terminate",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeWorkerTerminate(loop,receiver)}))
		workerProto.set("getHeapSnapshot",tsNamedNativeMethod("getHeapSnapshot",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{tsNodeWorkerValidateOptions(loop,tsArg(args,0),"options");return tsNodeWorkerPromise(loop,tsNodeWorkerSnapshot(loop),false)}))
		workerProto.set("getHeapStatistics",tsNamedNativeMethod("getHeapStatistics",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNodeWorkerPromise(loop,tsNodeWorkerObject(loop,map[string]tsValue{"totalHeapSize":tsNumberValue(0),"totalHeapSizeExecutable":tsNumberValue(0),"totalPhysicalSize":tsNumberValue(0),"totalAvailableSize":tsNumberValue(0),"usedHeapSize":tsNumberValue(0),"heapSizeLimit":tsNumberValue(0),"mallocedMemory":tsNumberValue(0),"peakMallocedMemory":tsNumberValue(0),"doesZapGarbage":tsBooleanValue(false)}),false)}))
		workerProto.set("cpuUsage",tsNamedNativeMethod("cpuUsage",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{previous:=tsArg(args,0);user,system:=float64(0),float64(0);if !tsIsUndefined(previous)&&!tsIsNull(previous){tsNodeWorkerValidateOptions(loop,previous,"prevValue");user=tsNumber(tsGet(loop,previous,tsNodeWorkerKey("user")));system=tsNumber(tsGet(loop,previous,tsNodeWorkerKey("system")))};return tsNodeWorkerPromise(loop,tsNodeWorkerObject(loop,map[string]tsValue{"user":tsNumberValue(-user),"system":tsNumberValue(-system)}),false)}))
		workerProto.set("startCpuProfile",tsNamedNativeMethod("startCpuProfile",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{tsNodeWorkerValidateOptions(loop,tsArg(args,0),"options");handle:=tsNewObject();handle.set("stop",tsNamedNativeMethod("stop",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNodeWorkerPromise(loop,tsStringReference(tsStringUTF8("")),false)}));return tsNodeWorkerPromise(loop,tsObjectValue(handle),false)}))
		workerProto.set("startHeapProfile",tsNamedNativeMethod("startHeapProfile",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{tsNodeWorkerValidateOptions(loop,tsArg(args,0),"options");if state:=tsNodeWorkerFor(receiver);state!=nil&&state.exited{return tsNodeWorkerPromise(loop,tsNodeWorkerError("Error","ERR_WORKER_NOT_RUNNING","Worker instance not running"),true)};profile:="{\"head\":{\"callFrame\":{\"functionName\":\"(root)\",\"scriptId\":\"0\",\"url\":\"\",\"lineNumber\":-1,\"columnNumber\":-1},\"selfSize\":0,\"id\":1,\"children\":[]},\"samples\":[]}";handle:=tsNewObject();handle.set("stop",tsNamedNativeMethod("stop",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNodeWorkerPromise(loop,tsStringReference(tsStringUTF8(profile)),false)}));return tsNodeWorkerPromise(loop,tsObjectValue(handle),false)}))
		portProto:=tsNewObject();portProto.prototype=tsObjectValue(eventProto)
		if portProto.descriptors==nil{portProto.descriptors=map[string]*tsDescriptor{}}
		for _,property:=range []string{"onmessage","onmessageerror"}{name:=property;portProto.descriptors[name]=&tsDescriptor{accessor:true,enumerable:true,configurable:true,getter:tsNamedNativeMethod("get "+name,0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeWorkerPortFor(receiver);if state==nil{return tsU};if name=="onmessage"{return state.onmessage};return state.onmessageerror}),setter:tsNamedNativeMethod("set "+name,1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeWorkerPortFor(receiver);if state==nil{tsPropertyFailure("Illegal invocation")};handler:=tsArg(args,0);if handler.kind!=tsFunctionKind{handler=tsU};if name=="onmessage"{state.onmessage=handler}else{state.onmessageerror=handler};if handler.kind==tsFunctionKind{state.started=true;tsNodeWorkerPortSchedule(loop,state)};return tsU})}}
		for _,method:=range []string{"on","addListener","once"}{operation:=method;portProto.set(method,tsNamedNativeMethod(method,2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeWorkerPortFor(receiver);if state==nil{tsPropertyFailure("Illegal invocation")};event:=tsArg(args,0);if tsText(event)=="message"||tsText(event)=="messageerror"{state.started=true;tsNodeWorkerPortSchedule(loop,state)};return tsNodeEventAdd(loop,receiver,event,tsArg(args,1),false,operation=="once")}))}
	tsNodeWorkerThreads.portPrototype=tsObjectValue(portProto)
		portCtor:=tsFunc(func(_ *tsLoop,_ ...tsValue)tsValue{tsPropertyFailure("Illegal constructor");return tsU});portCtor.name="MessagePort";portCtor.constructible=true;portCtor.construct=func(loop *tsLoop,_ ...tsValue)tsValue{return tsNodeWorkerPortCreate(loop,"")};portCtor.properties=tsNewObject();portCtor.properties.set("prototype",tsObjectValue(portProto));tsNodeWorkerThreads.portConstructor=tsFunctionValue(portCtor)
		channelCtor:=tsFunc(func(_ *tsLoop,_ ...tsValue)tsValue{tsPropertyFailure("Class constructor MessageChannel cannot be invoked without 'new'");return tsU});channelCtor.name="MessageChannel";channelCtor.constructible=true;channelCtor.construct=func(loop *tsLoop,_ ...tsValue)tsValue{channel:=tsNodeWorkerMakeChannel(loop);object:=(*tsObject)(channel.ref);object.prototype=tsObjectValue(tsNewObject());return channel};channelCtor.properties=tsNewObject();channelCtor.properties.set("prototype",tsObjectValue(tsNewObject()));tsNodeWorkerThreads.channelConstructor=tsFunctionValue(channelCtor)
		broadcastProto:=tsNewObject();broadcastProto.prototype=tsObjectValue(eventProto);broadcastProto.descriptors=map[string]*tsDescriptor{};for _,property:=range []string{"onmessage","onmessageerror"}{broadcastProto.descriptors[property]=portProto.descriptors[property]};tsNodeWorkerThreads.broadcastPrototype=tsObjectValue(broadcastProto)
		broadcastCtor:=tsFunc(func(_ *tsLoop,_ ...tsValue)tsValue{tsPropertyFailure("Class constructor BroadcastChannel cannot be invoked without 'new'");return tsU});broadcastCtor.name="BroadcastChannel";broadcastCtor.constructible=true;broadcastCtor.construct=func(loop *tsLoop,args ...tsValue)tsValue{name:=tsText(tsArg(args,0));port:=tsNodeWorkerPortCreate(loop,name);object:=(*tsObject)(port.ref);object.prototype=tsObjectValue(broadcastProto);object.set("name",tsNodeWorkerKey(name));group:=tsNodeWorkerThreads.channels[name];if group==nil{group=&tsNodeBroadcastGroup{};tsNodeWorkerThreads.channels[name]=group};group.ports=append(group.ports,tsNodeWorkerPortFor(port));return port};broadcastCtor.properties=tsNewObject();broadcastCtor.properties.set("prototype",tsObjectValue(broadcastProto));tsNodeWorkerThreads.broadcastConstructor=tsFunctionValue(broadcastCtor)
		share:=tsCall(loop,tsGet(loop,tsSymbolFunction(),tsNodeWorkerKey("for")),tsNodeWorkerKey("nodejs.worker_threads.SHARE_ENV"));tsNodeWorkerThreads.shareEnv=share
		workerCtor:=tsFunc(func(_ *tsLoop,_ ...tsValue)tsValue{tsPropertyFailure("Class constructor Worker cannot be invoked without 'new'");return tsU});workerCtor.name="Worker";workerCtor.constructible=true;workerCtor.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeWorkerCreate(loop,args)};workerCtor.properties=tsNewObject();workerCtor.properties.set("prototype",tsObjectValue(workerProto));tsNodeWorkerThreads.workerConstructor=tsFunctionValue(workerCtor)
		module:=tsNewObject();module.set("Worker",tsNodeWorkerThreads.workerConstructor);module.set("MessagePort",tsNodeWorkerThreads.portConstructor);module.set("MessageChannel",tsNodeWorkerThreads.channelConstructor);module.set("BroadcastChannel",tsNodeWorkerThreads.broadcastConstructor);module.set("workerData",tsU);module.set("parentPort",tsNull);module.set("resourceLimits",tsObjectValue(tsNewObject()));module.set("isMainThread",tsBooleanValue(true));module.set("isInternalThread",tsBooleanValue(false));module.set("threadId",tsNumberValue(0));module.set("threadName",tsNull);module.set("SHARE_ENV",share)
		module.set("receiveMessageOnPort",tsNamedNativeMethod("receiveMessageOnPort",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeWorkerReceive(loop,args)}))
		module.set("getEnvironmentData",tsNamedNativeMethod("getEnvironmentData",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeWorkerGetEnvironmentData(loop,tsArg(args,0))}))
		module.set("setEnvironmentData",tsNamedNativeMethod("setEnvironmentData",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{tsNodeWorkerSetEnvironmentData(loop,tsArg(args,0),tsArg(args,1));return tsU}))
		module.set("markAsUntransferable",tsNamedNativeMethod("markAsUntransferable",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{value:=tsArg(args,0);if tsNodeWorkerMarkable(value){tsNodeWorkerThreads.markedUntransferable[uintptr(value.ref)]=true};return tsU}))
		module.set("markAsUncloneable",tsNamedNativeMethod("markAsUncloneable",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{value:=tsArg(args,0);if tsNodeWorkerMarkable(value){tsNodeWorkerThreads.markedUncloneable[uintptr(value.ref)]=true};return tsU}))
		module.set("isMarkedAsUntransferable",tsNamedNativeMethod("isMarkedAsUntransferable",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsBooleanValue(tsNodeWorkerThreads.markedUntransferable[uintptr(tsArg(args,0).ref)])}))
		module.set("isMarkedAsUncloneable",tsNamedNativeMethod("isMarkedAsUncloneable",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsBooleanValue(tsNodeWorkerThreads.markedUncloneable[uintptr(tsArg(args,0).ref)])}))
		module.set("moveMessagePortToContext",tsNamedNativeMethod("moveMessagePortToContext",2,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{if tsNodeWorkerPortFor(tsArg(args,0))==nil{panic(tsThrown{tsNodeWorkerError("TypeError","ERR_INVALID_ARG_TYPE","The \"port\" argument must be a MessagePort")})};return tsArg(args,0)}))
		messaging:=tsNodeModuleRequire(loop,[]tsValue{tsNodeWorkerKey("internal/worker/messaging")});tsCall(loop,tsGet(loop,messaging,tsNodeWorkerKey("initThreadInfo")),tsNumberValue(float64(tsNodeWorkerThreads.threadID)),tsBooleanValue(tsNodeWorkerThreads.isMainThread),tsNodeWorkerThreads.channelConstructor);module.set("postMessageToThread",tsGet(loop,messaging,tsNodeWorkerKey("postMessageToThread")));module.set("default",tsObjectValue(module));tsNodeWorkerThreads.module=tsObjectValue(module)
	})
}

func tsNodeWorkerThreadsModule(loop *tsLoop) tsValue {tsNodeWorkerThreadsInit(loop);return tsNodeWorkerThreads.module}
`
