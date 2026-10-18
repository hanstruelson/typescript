package goemit

// internalClusterWorkerRuntime ports the common Cluster Worker object. The
// primary and child adapters attach their process-specific disconnect methods
// to this shared prototype.
const internalClusterWorkerRuntime = `
type tsInternalClusterWorkerState struct { object *tsObject; id int; process tsValue }
var tsInternalClusterWorker struct { sync.Once; constructor tsValue; prototype *tsObject; workers map[*tsObject]*tsInternalClusterWorkerState }

func tsInternalClusterWorkerFor(receiver tsValue) *tsInternalClusterWorkerState {
	if receiver.kind != tsObjectKind { return nil }
	return tsInternalClusterWorker.workers[(*tsObject)(receiver.ref)]
}

func tsInternalClusterWorkerInit(loop *tsLoop) {
	tsInternalClusterWorker.Do(func() {
		tsNodeEventsInit(loop)
		tsInternalClusterWorker.workers = map[*tsObject]*tsInternalClusterWorkerState{}
		prototype := tsNewObject()
		prototype.prototype = tsObjectValue(tsNodeEventsSingleton.prototype)
		prototype.set("kill", tsNamedNativeMethod("kill", 1, func(loop *tsLoop, receiver tsValue, args ...tsValue) tsValue {
			process := tsU
			if state := tsInternalClusterWorkerFor(receiver); state != nil { process = state.process }
			return tsCallReceiver(loop, tsGet(loop, process, tsStringReference(tsStringUTF8("kill"))), process, tsArg(args, 0))
		}))
		prototype.set("send", tsNamedNativeMethod("send", 1, func(loop *tsLoop, receiver tsValue, args ...tsValue) tsValue {
			process := tsU
			if state := tsInternalClusterWorkerFor(receiver); state != nil { process = state.process }
			return tsCallReceiver(loop, tsGet(loop, process, tsStringReference(tsStringUTF8("send"))), process, args...)
		}))
		prototype.set("isDead", tsNamedNativeMethod("isDead", 0, func(loop *tsLoop, receiver tsValue, _ ...tsValue) tsValue {
			state := tsInternalClusterWorkerFor(receiver)
			if state == nil { return tsBooleanValue(true) }
			process := state.process
			exitCode := tsGet(loop, process, tsStringReference(tsStringUTF8("exitCode")))
			signalCode := tsGet(loop, process, tsStringReference(tsStringUTF8("signalCode")))
			return tsBooleanValue(!tsNullish(exitCode) || !tsNullish(signalCode))
		}))
		prototype.set("isConnected", tsNamedNativeMethod("isConnected", 0, func(loop *tsLoop, receiver tsValue, _ ...tsValue) tsValue {
			state := tsInternalClusterWorkerFor(receiver)
			return tsBooleanValue(state != nil && tsTruthy(tsGet(loop, state.process, tsStringReference(tsStringUTF8("connected")))))
		}))
		prototype.set("disconnect", tsNamedNativeMethod("disconnect", 0, func(loop *tsLoop, receiver tsValue, _ ...tsValue) tsValue {
			state := tsInternalClusterWorkerFor(receiver)
			if state != nil { disconnect := tsGet(loop, state.process, tsStringReference(tsStringUTF8("disconnect"))); if disconnect.kind == tsFunctionKind { tsCallReceiver(loop, disconnect, state.process) } }
			return receiver
		}))
		prototype.set("destroy", tsNamedNativeMethod("destroy", 1, func(loop *tsLoop, receiver tsValue, args ...tsValue) tsValue {
			state := tsInternalClusterWorkerFor(receiver)
			if state != nil { kill := tsGet(loop, state.process, tsStringReference(tsStringUTF8("kill"))); if kill.kind == tsFunctionKind { tsCallReceiver(loop, kill, state.process, tsArg(args, 0)) } }
			return receiver
		}))
		constructor := tsFunc(func(_ *tsLoop, _ ...tsValue) tsValue { tsPropertyFailure("Class constructor Worker cannot be invoked without 'new'"); return tsU })
		constructor.name = "Worker"
		constructor.constructible = true
		constructor.properties = tsNewObject()
		constructor.properties.set("prototype", tsObjectValue(prototype))
		constructor.construct = func(loop *tsLoop, args ...tsValue) tsValue { return tsInternalClusterWorkerCreate(loop, args) }
		tsInternalClusterWorker.prototype = prototype
		tsInternalClusterWorker.constructor = tsFunctionValue(constructor)
	})
}

func tsInternalClusterWorkerCreate(loop *tsLoop, args []tsValue) tsValue {
	options := tsArg(args, 0)
	if tsIsUndefined(options) || tsIsNull(options) { options = tsObjectValue(tsNewObject()) }
	if options.kind != tsObjectKind && options.kind != tsInstanceKind { tsFSInvalid("options", "an object") }
	process := tsGet(loop, options, tsStringReference(tsStringUTF8("process")))
	stateValue := tsGet(loop, options, tsStringReference(tsStringUTF8("state")))
	stateName := "none"
	if !tsIsUndefined(stateValue) && !tsIsNull(stateValue) && tsText(stateValue) != "" { stateName = tsText(stateValue) }
	id := int(tsNumber(tsGet(loop, options, tsStringReference(tsStringUTF8("id")))))
	value := tsNodeEventCreate(loop, tsU, tsU)
	object := (*tsObject)(value.ref)
	object.prototype = tsObjectValue(tsInternalClusterWorker.prototype)
	worker := &tsInternalClusterWorkerState{object: object, id: id, process: process}
	tsInternalClusterWorker.workers[object] = worker
	object.set("exitedAfterDisconnect", tsU)
	object.set("state", tsStringReference(tsStringUTF8(stateName)))
	object.set("id", tsNumberValue(float64(id)))
	if !tsIsUndefined(process) && !tsIsNull(process) {
		for _, eventName := range []string{"error", "message"} {
			eventName := eventName
			listener := tsFunctionValue(tsFunc(func(loop *tsLoop, values ...tsValue) tsValue {
				return tsCallReceiver(loop, tsGet(loop, value, tsStringReference(tsStringUTF8("emit"))), value, append([]tsValue{tsStringReference(tsStringUTF8(eventName))}, values...)...)
			}))
			on := tsGet(loop, process, tsStringReference(tsStringUTF8("on")))
			if on.kind == tsFunctionKind { tsCallReceiver(loop, on, process, tsStringReference(tsStringUTF8(eventName)), listener) }
		}
	}
	return value
}

func tsInternalClusterWorkerModule(loop *tsLoop) tsValue {
	tsInternalClusterWorkerInit(loop)
	return tsInternalClusterWorker.constructor
}
`
