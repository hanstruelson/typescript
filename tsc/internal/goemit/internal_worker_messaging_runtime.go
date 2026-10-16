package goemit

// internalWorkerMessagingRuntime routes Node's workerMessage control messages
// over the MessagePorts exposed by nodeWorkerThreadsRuntime.
const internalWorkerMessagingRuntime = `
const (
	tsInternalWorkerMessageDelivered = 0
	tsInternalWorkerMessageNoListeners = 1
	tsInternalWorkerMessageListenerError = 2
)
type tsInternalWorkerMessagePending struct { promise *tsPromise; timeout *tsTask }
var tsInternalWorkerMessagingSingleton struct {
	sync.Once
	module tsValue
	currentThreadID int
	isMainThread bool
	messageChannel tsValue
	mainThreadPort tsValue
	threadPorts map[int]tsValue
	pending map[*tsObject]*tsInternalWorkerMessagePending
	entryEvaluated bool
	pendingMessages []tsValue
}

func tsInternalWorkerMessagingKey(name string) tsValue { return tsStringReference(tsStringUTF8(name)) }
func tsInternalWorkerMessageGet(loop *tsLoop,value tsValue,name string) tsValue { return tsGet(loop,value,tsInternalWorkerMessagingKey(name)) }
func tsNodeWorkerPortCall(loop *tsLoop,port tsValue,name string,args ...tsValue)tsValue{method:=tsGet(loop,port,tsNodeWorkerKey(name));if method.kind!=tsFunctionKind{return tsU};return tsCallReceiver(loop,method,port,args...)}
func tsInternalWorkerMessagingModule(loop *tsLoop) tsValue {
	tsInternalWorkerMessagingSingleton.Do(func() {
		tsInternalWorkerMessagingSingleton.threadPorts=map[int]tsValue{}
		tsInternalWorkerMessagingSingleton.pending=map[*tsObject]*tsInternalWorkerMessagePending{}
		module := tsNewObject()
		module.set("initThreadInfo", tsNamedNativeMethod("initThreadInfo", 3, func(_ *tsLoop, _ tsValue, args ...tsValue) tsValue {
			tsInternalWorkerMessagingSingleton.currentThreadID=int(tsNumber(tsArg(args,0)))
			tsInternalWorkerMessagingSingleton.isMainThread=tsTruthy(tsArg(args,1))
			tsInternalWorkerMessagingSingleton.messageChannel=tsArg(args,2)
			tsNodeWorkerThreads.threadID=tsInternalWorkerMessagingSingleton.currentThreadID
			tsNodeWorkerThreads.isMainThread=tsInternalWorkerMessagingSingleton.isMainThread
			return tsU
		}))
		module.set("createMessagingChannel", tsNamedNativeMethod("createMessagingChannel", 0, func(loop *tsLoop, _ tsValue, _ ...tsValue) tsValue {
			channel:=tsNodeWorkerMakeChannel(loop)
			return tsNodeWorkerObject(loop,map[string]tsValue{"portToMain":tsGet(loop,channel,tsNodeWorkerKey("port1")),"portToWorker":tsGet(loop,channel,tsNodeWorkerKey("port2"))})
		}))
		module.set("registerMainThreadPort",tsNamedNativeMethod("registerMainThreadPort",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{
			threadID,port:=int(tsNumber(tsArg(args,0))),tsArg(args,1)
			registration:=tsNodeWorkerObject(loop,map[string]tsValue{"type":tsInternalWorkerMessagingKey("registerMainThreadPort"),"threadId":tsNumberValue(float64(threadID)),"port":port})
			if tsInternalWorkerMessagingSingleton.isMainThread{tsInternalWorkerHandleMessageFromThread(loop,registration)}else if !tsNullish(tsInternalWorkerMessagingSingleton.mainThreadPort){tsNodeWorkerPortCall(loop,tsInternalWorkerMessagingSingleton.mainThreadPort,"postMessage",registration,tsArrayValue(&tsArray{values:[]tsValue{port}}))}
			return tsU
		}))
		module.set("destroyMainThreadPort",tsNamedNativeMethod("destroyMainThreadPort",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{
			message:=tsNodeWorkerObject(loop,map[string]tsValue{"type":tsInternalWorkerMessagingKey("unregisterMainThreadPort"),"threadId":tsArg(args,0)})
			if tsInternalWorkerMessagingSingleton.isMainThread{tsInternalWorkerHandleMessageFromThread(loop,message)}else if !tsNullish(tsInternalWorkerMessagingSingleton.mainThreadPort){tsNodeWorkerPortCall(loop,tsInternalWorkerMessagingSingleton.mainThreadPort,"postMessage",message)}
			return tsU
		}))
		module.set("setupMainThreadPort",tsNamedNativeMethod("setupMainThreadPort",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{
			port,registerHook:=tsArg(args,0),tsArg(args,1);tsInternalWorkerMessagingSingleton.mainThreadPort=port;tsInternalWorkerMessagingSingleton.entryEvaluated=false
			listener:=tsFunctionValue(tsFunc(func(loop *tsLoop,eventArgs ...tsValue)tsValue{message:=tsArg(eventArgs,0);if !tsInternalWorkerMessagingSingleton.entryEvaluated{tsInternalWorkerMessagingSingleton.pendingMessages=append(tsInternalWorkerMessagingSingleton.pendingMessages,message)}else{tsInternalWorkerHandleMessageFromMainThread(loop,message)};return tsU}))
			tsNodeWorkerPortCall(loop,port,"on",tsInternalWorkerMessagingKey("message"),listener)
			hook:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{tsInternalWorkerMessagingSingleton.entryEvaluated=true;pending:=tsInternalWorkerMessagingSingleton.pendingMessages;tsInternalWorkerMessagingSingleton.pendingMessages=nil;for _,message:=range pending{tsInternalWorkerHandleMessageFromMainThread(loop,message)};return tsU}))
			tsCall(loop,registerHook,hook);tsNodeWorkerPortCall(loop,port,"unref");return tsU
		}))
		module.set("postMessageToThread",tsNamedNativeMethod("postMessageToThread",4,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsInternalWorkerPostMessageToThread(loop,args)}))
		module.set("default",tsObjectValue(module));tsInternalWorkerMessagingSingleton.module=tsObjectValue(module)
	})
	return tsInternalWorkerMessagingSingleton.module
}

func tsInternalWorkerMessageFailure(response int) tsValue {
	switch response {
	case tsInternalWorkerMessageNoListeners:return tsNodeWorkerError("Error","ERR_WORKER_MESSAGING_FAILED","The destination thread no longer exists or is not listening for workerMessage events.")
	case tsInternalWorkerMessageListenerError:return tsNodeWorkerError("Error","ERR_WORKER_MESSAGING_ERRORED","The destination thread threw an error while processing the message.")
	default:return tsU
	}
}
func tsInternalWorkerMessageComplete(loop *tsLoop,memory tsValue,response int) {
	if memory.kind!=tsObjectKind{return};object:=(*tsObject)(memory.ref);object.set("result",tsNumberValue(float64(response)));object.set("status",tsNumberValue(1))
	pending:=tsInternalWorkerMessagingSingleton.pending[object];if pending==nil{return};delete(tsInternalWorkerMessagingSingleton.pending,object);if pending.timeout!=nil{loop.clearTimeout(tsTaskValue(pending.timeout))};failure:=tsInternalWorkerMessageFailure(response);if tsIsUndefined(failure){pending.promise.settle(tsResult{tsU,false})}else{pending.promise.settle(tsResult{failure,true})}
}
func tsInternalWorkerReceiveMessageFromWorker(loop *tsLoop,source,value,memory tsValue) {
	response:=tsInternalWorkerMessageNoListeners;process:=tsNodeProcessModule(loop);method:=tsGet(loop,process,tsInternalWorkerMessagingKey("listeners"));if method.kind==tsFunctionKind{
		listeners:=tsCallReceiver(loop,method,process,tsInternalWorkerMessagingKey("workerMessage"));if listeners.kind==tsArrayKind&&(*tsArray)(listeners.ref).length()>0{
			response=tsInternalWorkerMessageDelivered
			func(){defer func(){if recover()!=nil{response=tsInternalWorkerMessageListenerError}}();array:=(*tsArray)(listeners.ref);for index:=0;index<array.length();index++{tsCallReceiver(loop,array.at(index),process,value,source)}}()
		}
	}
	tsInternalWorkerMessageComplete(loop,memory,response)
}
func tsInternalWorkerSendMessageToWorker(loop *tsLoop,source,destination,value,transferList,memory tsValue) {
	if tsStrictEqual(destination,tsNumberValue(0)){tsInternalWorkerReceiveMessageFromWorker(loop,source,value,memory);return}
	port,ok:=tsInternalWorkerMessagingSingleton.threadPorts[int(tsNumber(destination))];if !ok{tsInternalWorkerMessageComplete(loop,memory,tsInternalWorkerMessageNoListeners);return}
	message:=tsNodeWorkerObject(loop,map[string]tsValue{"type":tsInternalWorkerMessagingKey("receiveMessageFromWorker"),"source":source,"value":value,"memory":memory})
	tsNodeWorkerPortCall(loop,port,"postMessage",message,transferList)
}
func tsInternalWorkerHandleMessageFromThread(loop *tsLoop,message tsValue) {
	switch tsText(tsInternalWorkerMessageGet(loop,message,"type")){
	case "registerMainThreadPort":
		threadID:=int(tsNumber(tsInternalWorkerMessageGet(loop,message,"threadId")));port:=tsInternalWorkerMessageGet(loop,message,"port");tsInternalWorkerMessagingSingleton.threadPorts[threadID]=port
		onMessage:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{incoming:=tsArg(args,0);switch tsText(tsInternalWorkerMessageGet(loop,incoming,"type")){case "sendMessageToWorker":tsInternalWorkerSendMessageToWorker(loop,tsInternalWorkerMessageGet(loop,incoming,"source"),tsInternalWorkerMessageGet(loop,incoming,"destination"),tsInternalWorkerMessageGet(loop,incoming,"value"),tsInternalWorkerMessageGet(loop,incoming,"transferList"),tsInternalWorkerMessageGet(loop,incoming,"memory"))};return tsU}))
		onClose:=tsFunctionValue(tsFunc(func(_ *tsLoop,_ ...tsValue)tsValue{if current:=tsInternalWorkerMessagingSingleton.threadPorts[threadID];tsStrictEqual(current,port){delete(tsInternalWorkerMessagingSingleton.threadPorts,threadID)};return tsU}))
		tsNodeWorkerPortCall(loop,port,"on",tsInternalWorkerMessagingKey("message"),onMessage);tsNodeWorkerPortCall(loop,port,"on",tsInternalWorkerMessagingKey("close"),onClose);tsNodeWorkerPortCall(loop,port,"unref")
	case "unregisterMainThreadPort":
		threadID:=int(tsNumber(tsInternalWorkerMessageGet(loop,message,"threadId")));if port,ok:=tsInternalWorkerMessagingSingleton.threadPorts[threadID];ok{tsNodeWorkerPortCall(loop,port,"close");delete(tsInternalWorkerMessagingSingleton.threadPorts,threadID)}
	case "sendMessageToWorker":
		tsInternalWorkerSendMessageToWorker(loop,tsInternalWorkerMessageGet(loop,message,"source"),tsInternalWorkerMessageGet(loop,message,"destination"),tsInternalWorkerMessageGet(loop,message,"value"),tsInternalWorkerMessageGet(loop,message,"transferList"),tsInternalWorkerMessageGet(loop,message,"memory"))
	}
}
func tsInternalWorkerHandleMessageFromMainThread(loop *tsLoop,message tsValue) {
	if tsStrictEqual(tsInternalWorkerMessageGet(loop,message,"type"),tsInternalWorkerMessagingKey("receiveMessageFromWorker")){tsInternalWorkerReceiveMessageFromWorker(loop,tsInternalWorkerMessageGet(loop,message,"source"),tsInternalWorkerMessageGet(loop,message,"value"),tsInternalWorkerMessageGet(loop,message,"memory"))}
}
func tsInternalWorkerPostMessageToThread(loop *tsLoop,args []tsValue)(result tsValue) {
	promise:=loop.promise();result=tsPromiseValue(promise)
	defer func(){if failure:=recover();failure!=nil{promise.settle(tsResult{tsUnwrap(failure),true});result=tsPromiseValue(promise)}}()
	threadID,value,transferList,timeout:=tsArg(args,0),tsArg(args,1),tsArg(args,2),tsArg(args,3)
	if tsIsNumeric(transferList)&&tsIsUndefined(timeout){timeout=transferList;transferList=tsArrayValue(&tsArray{})};if tsIsUndefined(transferList){transferList=tsArrayValue(&tsArray{})}
	if !tsIsUndefined(timeout){tsInternalValidatorNumber(timeout,tsInternalWorkerMessagingKey("timeout"),tsNumberValue(0),tsU,false)}
	if tsStrictEqual(threadID,tsNumberValue(float64(tsInternalWorkerMessagingSingleton.currentThreadID))){promise.settle(tsResult{tsNodeWorkerError("Error","ERR_WORKER_MESSAGING_SAME_THREAD","Cannot send a message to the same thread."),true});return result}
	memoryObject:=tsNewObject();memory:=tsObjectValue(memoryObject);pending:=&tsInternalWorkerMessagePending{promise:promise};tsInternalWorkerMessagingSingleton.pending[memoryObject]=pending
	if !tsIsUndefined(timeout){pending.timeout=loop.setTimeout(tsFunctionValue(tsFunc(func(_ *tsLoop,_ ...tsValue)tsValue{if current:=tsInternalWorkerMessagingSingleton.pending[memoryObject];current==pending{delete(tsInternalWorkerMessagingSingleton.pending,memoryObject);promise.settle(tsResult{tsNodeWorkerError("Error","ERR_WORKER_MESSAGING_TIMEOUT","The operation timed out."),true})};return tsU})),timeout)}
	message:=tsNodeWorkerObject(loop,map[string]tsValue{"type":tsInternalWorkerMessagingKey("sendMessageToWorker"),"source":tsNumberValue(float64(tsInternalWorkerMessagingSingleton.currentThreadID)),"destination":threadID,"value":value,"transferList":transferList,"memory":memory})
	if tsInternalWorkerMessagingSingleton.isMainThread{tsInternalWorkerHandleMessageFromThread(loop,message)}else if !tsNullish(tsInternalWorkerMessagingSingleton.mainThreadPort){tsNodeWorkerPortCall(loop,tsInternalWorkerMessagingSingleton.mainThreadPort,"postMessage",message,transferList)}else{tsInternalWorkerMessageComplete(loop,memory,tsInternalWorkerMessageNoListeners)}
	return result
}
`
