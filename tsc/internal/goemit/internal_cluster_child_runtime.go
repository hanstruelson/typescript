package goemit

// internalClusterChildRuntime maps Bun's worker-side cluster protocol onto
// the generated process/event values. Native descriptor adoption is modeled
// by lightweight handle objects until the runtime has process IPC handles.
const internalClusterChildRuntime = `
type tsInternalClusterChildIndex struct { next int; values map[int]bool }
type tsInternalClusterChildState struct {
	sync.Once
	module, worker, process, ownerKey tsValue
	handles map[string]tsValue
	indexes map[string]*tsInternalClusterChildIndex
	callbacks map[uint64]tsValue
	sequence uint64
}
var tsInternalClusterChild tsInternalClusterChildState

func tsInternalClusterChildKey(name string) tsValue { return tsStringReference(tsStringUTF8(name)) }
func tsInternalClusterChildGet(loop *tsLoop, value tsValue, name string) tsValue { return tsGet(loop, value, tsInternalClusterChildKey(name)) }
func tsInternalClusterChildSet(loop *tsLoop, value tsValue, name string, item tsValue) { tsSet(loop, value, tsInternalClusterChildKey(name), item) }
func tsInternalClusterChildEvent(loop *tsLoop, target tsValue, event string, args ...tsValue) tsValue { return tsCallReceiver(loop, tsInternalClusterChildGet(loop, target, "emit"), target, append([]tsValue{tsInternalClusterChildKey(event)}, args...)...) }
func tsInternalClusterChildDeleteHandle(key string) { delete(tsInternalClusterChild.handles, key) }
func tsInternalClusterChildRemoveIndex(key string, index int) {
	set := tsInternalClusterChild.indexes[key]
	if set == nil { return }
	delete(set.values, index)
	if len(set.values) == 0 { delete(tsInternalClusterChild.indexes, key) }
}
func tsInternalClusterChildSend(loop *tsLoop, message, callback tsValue) tsValue {
	process := tsInternalClusterChild.process
	if !tsTruthy(tsInternalClusterChildGet(loop, process, "connected")) { return tsBooleanValue(false) }
	sequence := tsInternalClusterChild.sequence
	tsInternalClusterChild.sequence++
	wire := tsNewObject()
	wire.set("cmd", tsInternalClusterChildKey("NODE_CLUSTER"))
	if message.kind == tsObjectKind || message.kind == tsInstanceKind {
		for _, key := range (*tsObject)(message.ref).order { wire.set(key, (*tsObject)(message.ref).values[key]) }
	}
	wire.set("seq", tsNumberValue(float64(sequence)))
	if callback.kind == tsFunctionKind { tsInternalClusterChild.callbacks[sequence] = callback }
	send := tsInternalClusterChildGet(loop, process, "send")
	if send.kind != tsFunctionKind { return tsBooleanValue(false) }
	return tsCallReceiver(loop, send, process, tsObjectValue(wire))
}
func tsInternalClusterChildMakeConnectionHandle(loop *tsLoop, fd tsValue) tsValue {
	handle := tsNewObject()
	handle.set("fd", fd)
	closed := false
	handle.set("close", tsNamedNativeMethod("close", 1, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
		if !closed { closed = true }
		callback := tsArg(args, 0)
		if callback.kind == tsFunctionKind { tsNodeClusterNextTick(loop, callback) }
		return tsU
	}))
	return tsObjectValue(handle)
}
func tsInternalClusterChildMakeSharedHandle(loop *tsLoop, fd tsValue) tsValue {
	handle := tsNewObject()
	handle.set("sharedFd", fd)
	handle.set("adopted", tsBooleanValue(false))
	fdOpen := true
	handle.set("close", tsNamedNativeMethod("close", 1, func(loop *tsLoop, receiver tsValue, args ...tsValue) tsValue {
		if fdOpen && !tsTruthy(tsInternalClusterChildGet(loop, receiver, "adopted")) { fdOpen = false }
		callback := tsArg(args, 0)
		if callback.kind == tsFunctionKind { tsNodeClusterNextTick(loop, callback) }
		return tsU
	}))
	return tsObjectValue(handle)
}
func tsInternalClusterChildOnMessage(loop *tsLoop, message, handle tsValue) {
	ack := tsInternalClusterChildGet(loop, message, "ack")
	if !tsIsUndefined(ack) {
		sequence := uint64(tsNumber(ack))
		if callback := tsInternalClusterChild.callbacks[sequence]; callback.kind == tsFunctionKind {
			delete(tsInternalClusterChild.callbacks, sequence)
			tsCallReceiver(loop, callback, tsU, message, handle)
			return
		}
	}
	if tsInternalClusterChildGet(loop, message, "act").kind == tsStringKind && tsText(tsInternalClusterChildGet(loop, message, "act")) == "newconn" && tsIsNumeric(handle) && tsNumber(handle) >= 0 {
		handle = tsInternalClusterChildMakeConnectionHandle(loop, handle)
	}
	process := tsInternalClusterChild.process
	tsInternalClusterChildEvent(loop, process, "internalMessage", message, handle)
	act := tsText(tsInternalClusterChildGet(loop, message, "act"))
	if act == "newconn" { tsInternalClusterChildOnConnection(loop, message, handle) }
	if act == "disconnect" {
		disconnect := tsInternalClusterChildGet(loop, tsInternalClusterChild.worker, "_disconnect")
		if disconnect.kind == tsFunctionKind { tsCallReceiver(loop, disconnect, tsInternalClusterChild.worker, tsBooleanValue(true)) }
	}
}
func tsInternalClusterChildOnConnection(loop *tsLoop, message, handle tsValue) {
	key := tsText(tsInternalClusterChildGet(loop, message, "key"))
	server, accepted := tsInternalClusterChild.handles[key], false
	if !tsIsUndefined(server) && !tsIsNull(server) {
		accepted = true
		if !tsIsUndefined(tsInternalClusterChild.ownerKey) {
			owner := tsGet(loop, server, tsInternalClusterChild.ownerKey)
			if !tsNullish(owner) {
				maximum, current := tsInternalClusterChildGet(loop, owner, "maxConnections"), tsInternalClusterChildGet(loop, owner, "_connections")
				if !tsNullish(maximum) && tsNumber(current) >= tsNumber(maximum) && !tsTruthy(tsInternalClusterChildGet(loop, owner, "dropMaxConnection")) { accepted = false }
			}
		}
	}
	ack := tsNewObject(); ack.set("ack", tsInternalClusterChildGet(loop, message, "seq")); ack.set("accepted", tsBooleanValue(accepted))
	tsInternalClusterChildSend(loop, tsObjectValue(ack), tsU)
	if accepted {
		onconnection := tsInternalClusterChildGet(loop, server, "onconnection")
		if onconnection.kind == tsFunctionKind { tsCallReceiver(loop, onconnection, server, tsNumberValue(0), handle) }
	} else {
		close := tsInternalClusterChildGet(loop, handle, "close")
		if close.kind == tsFunctionKind { tsCallReceiver(loop, close, handle) }
	}
}
func tsInternalClusterChildShared(loop *tsLoop, message, handle tsValue, indexesKey string, index int, callback tsValue) {
	key := tsText(tsInternalClusterChildGet(loop, message, "key"))
	original := tsInternalClusterChildGet(loop, handle, "close")
	released := false
	close := tsNamedNativeMethod("close", 1, func(loop *tsLoop, receiver tsValue, args ...tsValue) tsValue {
		if !released {
			released = true
			closeMessage := tsNewObject(); closeMessage.set("act", tsInternalClusterChildKey("close")); closeMessage.set("key", tsInternalClusterChildKey(key))
			tsInternalClusterChildSend(loop, tsObjectValue(closeMessage), tsU)
			tsInternalClusterChildDeleteHandle(key); tsInternalClusterChildRemoveIndex(indexesKey, index)
		}
		if original.kind == tsFunctionKind { return tsCallReceiver(loop, original, handle, args...) }
		return tsU
	})
	tsInternalClusterChildSet(loop, handle, "close", close)
	tsInternalClusterChild.handles[key] = handle
	tsCall(loop, callback, tsInternalClusterChildGet(loop, message, "errno"), handle, message)
}
func tsInternalClusterChildRoundRobin(loop *tsLoop, message tsValue, indexesKey string, index int, callback tsValue) {
	if errno := tsInternalClusterChildGet(loop, message, "errno"); tsTruthy(errno) { tsCall(loop, callback, errno, tsNull, message); return }
	key := tsText(tsInternalClusterChildGet(loop, message, "key"))
	handle := tsNewObject(); handle.set("refed", tsBooleanValue(true))
	handle.set("ref", tsNamedNativeMethod("ref", 0, func(loop *tsLoop, receiver tsValue, _ ...tsValue) tsValue { tsSet(loop, receiver, tsInternalClusterChildKey("refed"), tsBooleanValue(true)); return receiver }))
	handle.set("unref", tsNamedNativeMethod("unref", 0, func(loop *tsLoop, receiver tsValue, _ ...tsValue) tsValue { tsSet(loop, receiver, tsInternalClusterChildKey("refed"), tsBooleanValue(false)); return receiver }))
	handle.set("listen", tsNamedNativeMethod("listen", 1, func(_ *tsLoop, _ tsValue, _ ...tsValue) tsValue { return tsNumberValue(0) }))
	handle.set("close", tsNamedNativeMethod("close", 0, func(loop *tsLoop, receiver tsValue, _ ...tsValue) tsValue {
		if tsIsUndefined(tsGet(loop, receiver, tsInternalClusterChildKey("key"))) { return tsU }
		closeMessage := tsNewObject(); closeMessage.set("act", tsInternalClusterChildKey("close")); closeMessage.set("key", tsInternalClusterChildKey(key))
		tsInternalClusterChildSend(loop, tsObjectValue(closeMessage), tsU); tsInternalClusterChildDeleteHandle(key); tsInternalClusterChildRemoveIndex(indexesKey, index); tsInternalClusterChildSet(loop, receiver, "key", tsU); return tsU
	}))
	handle.set("getsockname", tsNamedNativeMethod("getsockname", 1, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
		out := tsArg(args, 0); sockname := tsInternalClusterChildGet(loop, message, "sockname")
		if out.kind == tsObjectKind && sockname.kind == tsObjectKind { for _, field := range (*tsObject)(sockname.ref).order { tsSet(loop, out, tsInternalClusterChildKey(field), (*tsObject)(sockname.ref).values[field]) } }
		return tsNumberValue(0)
	}))
	handle.set("key", tsInternalClusterChildKey(key)); value := tsObjectValue(handle)
	tsInternalClusterChild.handles[key] = value
	tsCall(loop, callback, tsNumberValue(0), value, message)
}
func tsInternalClusterChildGetServer(loop *tsLoop, obj, options, callback tsValue) tsValue {
	address := tsInternalClusterChildGet(loop, options, "address")
	port := tsInternalClusterChildGet(loop, options, "port")
	if tsIsNumeric(port) && tsNumber(port) < 0 && tsText(tsInternalClusterChildGet(loop, tsNodeProcessModule(loop), "platform")) != "win32" {
		path := tsPathModuleFor(runtime.GOOS == "windows"); resolve := tsGet(loop, path, tsInternalClusterChildKey("resolve")); address = tsCallReceiver(loop, resolve, path, address)
	}
	indexesKey := strings.Join([]string{tsText(address), tsText(port), tsText(tsInternalClusterChildGet(loop, options, "addressType")), tsText(tsInternalClusterChildGet(loop, options, "fd"))}, ":")
	indexSet := tsInternalClusterChild.indexes[indexesKey]
	if indexSet == nil { indexSet = &tsInternalClusterChildIndex{values: map[int]bool{}}; tsInternalClusterChild.indexes[indexesKey] = indexSet }
	index := indexSet.next; indexSet.next++; indexSet.values[index] = true
	message := tsNewObject()
	if options.kind == tsObjectKind || options.kind == tsInstanceKind { for _, name := range (*tsObject)(options.ref).order { message.set(name, (*tsObject)(options.ref).values[name]) } }
	message.set("act", tsInternalClusterChildKey("queryServer")); message.set("index", tsNumberValue(float64(index))); message.set("data", tsNull); message.set("address", address)
	getData := tsInternalClusterChildGet(loop, obj, "_getServerData"); if getData.kind == tsFunctionKind { message.set("data", tsCallReceiver(loop, getData, obj)) }
	complete := tsFunctionValue(tsFunc(func(loop *tsLoop, values ...tsValue) tsValue {
		reply, handle := tsArg(values, 0), tsArg(values, 1)
		setData := tsInternalClusterChildGet(loop, obj, "_setServerData"); if setData.kind == tsFunctionKind { tsCallReceiver(loop, setData, obj, tsInternalClusterChildGet(loop, reply, "data")) }
		if tsIsNumeric(handle) && tsNumber(handle) >= 0 { handle = tsInternalClusterChildMakeSharedHandle(loop, handle) }
		if tsTruthy(handle) { tsInternalClusterChildShared(loop, reply, handle, indexesKey, index, callback) } else { tsInternalClusterChildRoundRobin(loop, reply, indexesKey, index, callback) }
		return tsU
	}))
	tsInternalClusterChildSend(loop, tsObjectValue(message), complete)
	onListening := tsFunctionValue(tsFunc(func(loop *tsLoop, _ ...tsValue) tsValue {
		if !indexSet.values[index] { return tsU }
		state := tsInternalClusterWorkerFor(tsInternalClusterChild.worker); if state != nil { state.object.set("state", tsInternalClusterChildKey("listening")) }
		addressMethod := tsInternalClusterChildGet(loop, obj, "address"); info := tsU; if addressMethod.kind == tsFunctionKind { info = tsCallReceiver(loop, addressMethod, obj) }
		message.set("act", tsInternalClusterChildKey("listening")); if tsNullish(info) { message.set("port", port) } else { actualPort := tsInternalClusterChildGet(loop, info, "port"); if !tsNullish(actualPort) && tsTruthy(actualPort) { message.set("port", actualPort) } }
		tsInternalClusterChildSend(loop, tsObjectValue(message), tsU); return tsU
	}))
	once := tsInternalClusterChildGet(loop, obj, "once"); if once.kind == tsFunctionKind { tsCallReceiver(loop, once, obj, tsInternalClusterChildKey("listening"), onListening) }
	return tsU
}
func tsInternalClusterChildSetup(loop *tsLoop) tsValue {
	process := tsInternalClusterChild.process
	env := tsInternalClusterChildGet(loop, process, "env")
	id, _ := strconv.Atoi(tsText(tsInternalClusterChildGet(loop, env, "NODE_UNIQUE_ID")))
	options := tsNewObject(); options.set("id", tsNumberValue(float64(id))); options.set("process", process); options.set("state", tsInternalClusterChildKey("online"))
	worker := tsConstruct(loop, tsInternalClusterWorkerModule(loop), tsObjectValue(options)); tsInternalClusterChild.worker = worker
	tsInternalClusterChildSet(loop, tsInternalClusterChild.module, "worker", worker)
	disconnectListener := tsFunctionValue(tsFunc(func(loop *tsLoop, _ ...tsValue) tsValue {
		tsInternalClusterChildSet(loop, process, "channel", tsNull); tsInternalClusterChildEvent(loop, worker, "disconnect")
		if !tsTruthy(tsInternalClusterChildGet(loop, worker, "exitedAfterDisconnect")) { exit := tsInternalClusterChildGet(loop, process, "exit"); if exit.kind == tsFunctionKind { tsCallReceiver(loop, exit, process, tsNumberValue(0)) } }
		return tsU
	}))
	once := tsInternalClusterChildGet(loop, process, "once"); if once.kind == tsFunctionKind { tsCallReceiver(loop, once, process, tsInternalClusterChildKey("disconnect"), disconnectListener) }
	messageListener := tsFunctionValue(tsFunc(func(loop *tsLoop, args ...tsValue) tsValue { tsInternalClusterChildOnMessage(loop, tsArg(args, 0), tsArg(args, 1)); return tsU }))
	on := tsInternalClusterChildGet(loop, process, "on"); if on.kind == tsFunctionKind { tsCallReceiver(loop, on, process, tsInternalClusterChildKey("internalMessage"), messageListener) }
	online := tsNewObject(); online.set("act", tsInternalClusterChildKey("online")); tsInternalClusterChildSend(loop, tsObjectValue(online), tsU)
	return worker
}
func tsInternalClusterChildDisconnect(loop *tsLoop, worker tsValue, primaryInitiated bool) tsValue {
	tsInternalClusterChildSet(loop, worker, "exitedAfterDisconnect", tsBooleanValue(true))
	waiting := 1
	finish := tsFunctionValue(tsFunc(func(loop *tsLoop, _ ...tsValue) tsValue { waiting--; if waiting == 0 { if primaryInitiated { if tsTruthy(tsInternalClusterChildGet(loop, tsInternalClusterChild.process, "connected")) { tsCallReceiver(loop, tsInternalClusterChildGet(loop, tsInternalClusterChild.process, "disconnect"), tsInternalClusterChild.process) } } else { done := tsFunctionValue(tsFunc(func(loop *tsLoop, _ ...tsValue) tsValue { if tsTruthy(tsInternalClusterChildGet(loop, tsInternalClusterChild.process, "connected")) { tsCallReceiver(loop, tsInternalClusterChildGet(loop, tsInternalClusterChild.process, "disconnect"), tsInternalClusterChild.process) }; return tsU })); message := tsNewObject(); message.set("act", tsInternalClusterChildKey("exitedAfterDisconnect")); tsInternalClusterChildSend(loop, tsObjectValue(message), done) } }; return tsU }))
	for key, handle := range tsInternalClusterChild.handles { waiting++; owner := tsU; if !tsIsUndefined(tsInternalClusterChild.ownerKey) { owner = tsGet(loop, handle, tsInternalClusterChild.ownerKey) }; closeTarget := handle; if !tsNullish(owner) { closeTarget = owner }; close := tsInternalClusterChildGet(loop, closeTarget, "close"); if close.kind == tsFunctionKind { tsCallReceiver(loop, close, closeTarget, finish) } else { tsCall(loop, finish) }; delete(tsInternalClusterChild.handles, key) }
	tsCall(loop, finish)
	return worker
}
func tsInternalClusterChildInit(loop *tsLoop) {
	tsInternalClusterChild.Do(func() {
		tsNodeEventsInit(loop); tsInternalClusterWorkerInit(loop)
		tsInternalClusterChild.process = tsNodeProcessModule(loop); tsInternalClusterChild.handles = map[string]tsValue{}; tsInternalClusterChild.indexes = map[string]*tsInternalClusterChildIndex{}; tsInternalClusterChild.callbacks = map[uint64]tsValue{}
		module := tsNodeEventCreate(loop, tsU, tsU); tsInternalClusterChild.module = module
		shared := tsInternalSharedModule(loop); tsInternalClusterChild.ownerKey = tsInternalClusterChildGet(loop, shared, "kClusterOwner")
		(*tsObject)(module.ref).set("isWorker", tsBooleanValue(true)); (*tsObject)(module.ref).set("isMaster", tsBooleanValue(false)); (*tsObject)(module.ref).set("isPrimary", tsBooleanValue(false)); (*tsObject)(module.ref).set("worker", tsNull); (*tsObject)(module.ref).set("Worker", tsInternalClusterWorkerModule(loop))
		(*tsObject)(module.ref).set("_setupWorker", tsNamedNativeMethod("_setupWorker", 0, func(loop *tsLoop, _ tsValue, _ ...tsValue) tsValue { return tsInternalClusterChildSetup(loop) }))
		(*tsObject)(module.ref).set("_getServer", tsNamedNativeMethod("_getServer", 3, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue { return tsInternalClusterChildGetServer(loop, tsArg(args, 0), tsArg(args, 1), tsArg(args, 2)) }))
		workerPrototype := tsInternalClusterWorker.prototype
		workerPrototype.set("disconnect", tsNamedNativeMethod("disconnect", 0, func(loop *tsLoop, receiver tsValue, _ ...tsValue) tsValue { state := tsInternalClusterWorkerFor(receiver); if state != nil && tsText(state.object.values["state"]) != "disconnecting" && tsText(state.object.values["state"]) != "destroying" { state.object.set("state", tsInternalClusterChildKey("disconnecting")); tsInternalClusterChildDisconnect(loop, receiver, false) }; return receiver }))
		workerPrototype.set("_disconnect", tsNamedNativeMethod("_disconnect", 1, func(loop *tsLoop, receiver tsValue, args ...tsValue) tsValue { return tsInternalClusterChildDisconnect(loop, receiver, tsTruthy(tsArg(args, 0))) }))
		workerPrototype.set("destroy", tsNamedNativeMethod("destroy", 0, func(loop *tsLoop, receiver tsValue, _ ...tsValue) tsValue { if tsText(tsInternalClusterChildGet(loop, receiver, "state")) == "destroying" { return tsU }; tsInternalClusterChildSet(loop, receiver, "exitedAfterDisconnect", tsBooleanValue(true)); connected := tsCallReceiver(loop, tsInternalClusterChildGet(loop, receiver, "isConnected"), receiver); process := tsInternalClusterChild.process; if !tsTruthy(connected) { tsCallReceiver(loop, tsInternalClusterChildGet(loop, process, "exit"), process, tsNumberValue(0)) } else { tsInternalClusterChildSet(loop, receiver, "state", tsInternalClusterChildKey("destroying")); message := tsNewObject(); message.set("act", tsInternalClusterChildKey("exitedAfterDisconnect")); done := tsFunctionValue(tsFunc(func(loop *tsLoop, _ ...tsValue) tsValue { tsCallReceiver(loop, tsInternalClusterChildGet(loop, process, "disconnect"), process); return tsU })); tsInternalClusterChildSend(loop, tsObjectValue(message), done) }; return tsU }))
	})
}
func tsInternalClusterChildModule(loop *tsLoop) tsValue {
	tsInternalClusterChildInit(loop)
	env := tsInternalClusterChildGet(loop, tsNodeProcessModule(loop), "env")
	if !tsIsUndefined(tsInternalClusterChildGet(loop, env, "NODE_UNIQUE_ID")) && tsIsNull(tsInternalClusterChildGet(loop, tsInternalClusterChild.module, "worker")) { tsInternalClusterChildSetup(loop) }
	return tsInternalClusterChild.module
}
`
