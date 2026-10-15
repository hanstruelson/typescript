package goemit

// internalWorkerStdioRuntime adapts Bun's MessagePort-backed worker stdio
// helpers to the runtime's Readable and Writable implementations.
const internalWorkerStdioRuntime = `
type tsInternalWorkerStdioPipe struct {
	port, stream, pending, ackListener tsValue
	incrementsPortRef, ended, startedReading bool
}

var tsInternalWorkerStdioSingleton struct { sync.Once; flushKey string; module tsValue }

func tsInternalWorkerStdioKey(name string) tsValue { return tsStringReference(tsStringUTF8(name)) }
func tsInternalWorkerStdioCall(loop *tsLoop, receiver tsValue, name string, args ...tsValue) tsValue {
	method := tsGet(loop, receiver, tsInternalWorkerStdioKey(name))
	if method.kind != tsFunctionKind { return tsU }
	return tsCallReceiver(loop, method, receiver, args...)
}

func tsInternalWorkerStdioRemove(loop *tsLoop, pipe *tsInternalWorkerStdioPipe, listener tsValue) {
	tsInternalWorkerStdioCall(loop, pipe.port, "removeEventListener", tsInternalWorkerStdioKey("message"), listener)
}

func tsInternalWorkerStdioReadable(loop *tsLoop, port tsValue, incrementsPortRef bool) tsValue {
	pipe := &tsInternalWorkerStdioPipe{port: port, incrementsPortRef: incrementsPortRef}
	module := tsInternalStreamReadableModule(loop)
	constructor := tsGet(loop, module, tsInternalWorkerStdioKey("Readable"))
	options := tsNewObject()
	options.set("read", tsFunctionValue(tsFunc(func(loop *tsLoop, _ ...tsValue) tsValue {
		if !pipe.startedReading && pipe.incrementsPortRef { pipe.startedReading = true; tsInternalWorkerStdioCall(loop, pipe.port, "ref") }
		if !pipe.ended { tsInternalWorkerStdioCall(loop, pipe.port, "postMessage", tsBooleanValue(true)) }
		return tsU
	})))
	pipe.stream = tsConstruct(constructor, tsObjectValue(options))
	pipe.ackListener = tsFunctionValue(tsFunc(func(loop *tsLoop, args ...tsValue) tsValue {
		if pipe.ended { return tsU }
		event := tsArg(args, 0)
		payload := tsGet(loop, event, tsInternalWorkerStdioKey("data"))
		if tsIsNull(payload) {
			pipe.ended = true
			tsInternalWorkerStdioCall(loop, pipe.stream, "push", tsNull)
			tsInternalWorkerStdioRemove(loop, pipe, pipe.ackListener)
			return tsU
		}
		if payload.kind == tsArrayKind {
			for _, chunk := range (*tsArray)(payload.ref).values {
				buffer := tsNodeBufferFrom(loop, []tsValue{chunk})
				tsInternalWorkerStdioCall(loop, pipe.stream, "push", buffer)
			}
		}
		return tsU
	}))
	tsInternalWorkerStdioCall(loop, port, "addEventListener", tsInternalWorkerStdioKey("message"), pipe.ackListener)
	tsInternalWorkerStdioCall(loop, port, "unref")
	closeListener := tsFunctionValue(tsFunc(func(loop *tsLoop, _ ...tsValue) tsValue {
		pipe.ended = true
		tsInternalWorkerStdioRemove(loop, pipe, pipe.ackListener)
		if pipe.startedReading && pipe.incrementsPortRef { pipe.startedReading = false; tsInternalWorkerStdioCall(loop, pipe.port, "unref") }
		return tsU
	}))
	tsInternalWorkerStdioCall(loop, pipe.stream, "on", tsInternalWorkerStdioKey("close"), closeListener)
	(*tsObject)(pipe.stream.ref).set("endFromOwner", tsFunctionValue(tsFunc(func(loop *tsLoop, _ ...tsValue) tsValue {
		if !pipe.ended { pipe.ended = true; tsInternalWorkerStdioCall(loop, pipe.stream, "push", tsNull); tsInternalWorkerStdioRemove(loop, pipe, pipe.ackListener) }
		return tsU
	})))
	return pipe.stream
}

func tsInternalWorkerStdioWritable(loop *tsLoop, port tsValue) tsValue {
	pipe := &tsInternalWorkerStdioPipe{port: port}
	module := tsInternalStreamWritableModule(loop)
	constructor := tsGet(loop, module, tsInternalWorkerStdioKey("Writable"))
	options := tsNewObject()
	options.set("decodeStrings", tsBooleanValue(false))
	pipe.ackListener = tsFunctionValue(tsFunc(func(loop *tsLoop, _ ...tsValue) tsValue {
		callback := pipe.pending
		if callback.kind == tsFunctionKind {
			pipe.pending = tsU
			tsInternalWorkerStdioCall(loop, pipe.port, "unref")
			tsCall(loop, callback)
		}
		return tsU
	}))
	tsInternalWorkerStdioCall(loop, port, "addEventListener", tsInternalWorkerStdioKey("message"), pipe.ackListener)
	tsInternalWorkerStdioCall(loop, port, "unref")
	options.set("writev", tsFunctionValue(tsFunc(func(loop *tsLoop, args ...tsValue) tsValue {
		chunks, callback := tsArg(args, 0), tsArg(args, 1)
		payload := &tsArray{}
		if chunks.kind == tsArrayKind {
			for _, entry := range (*tsArray)(chunks.ref).values {
				chunk := tsGet(loop, entry, tsInternalWorkerStdioKey("chunk"))
				encoding := tsGet(loop, entry, tsInternalWorkerStdioKey("encoding"))
				if chunk.kind == tsStringKind { chunk = tsNodeBufferFrom(loop, []tsValue{chunk, encoding}) }
				payload.values = append(payload.values, chunk)
			}
		}
		tsInternalWorkerStdioCall(loop, pipe.port, "postMessage", tsArrayValue(payload))
		process := tsNodeProcessModule(loop)
		if tsTruthy(tsGet(loop, process, tsInternalWorkerStdioKey("_exiting"))) { tsCall(loop, callback); return tsU }
		pipe.pending = callback
		tsInternalWorkerStdioCall(loop, pipe.port, "ref")
		return tsU
	})))
	options.set("final", tsFunctionValue(tsFunc(func(loop *tsLoop, args ...tsValue) tsValue {
		tsInternalWorkerStdioCall(loop, pipe.port, "postMessage", tsNull)
		tsCall(loop, tsArg(args, 0))
		return tsU
	})))
	options.set("destroy", tsFunctionValue(tsFunc(func(loop *tsLoop, args ...tsValue) tsValue {
		err, callback := tsArg(args, 0), tsArg(args, 1)
		if pipe.pending.kind == tsFunctionKind {
			pending := pipe.pending
			pipe.pending = tsU
			tsInternalWorkerStdioCall(loop, pipe.port, "unref")
			tsCall(loop, pending, err)
		}
		tsInternalWorkerStdioRemove(loop, pipe, pipe.ackListener)
		tsCall(loop, callback, err)
		return tsU
	})))
	stream := tsConstruct(constructor, tsObjectValue(options))
	pipe.stream = stream
	tsInternalWorkerStdioSingleton.Do(func() { tsInternalWorkerStdioSingleton.flushKey = tsPropertyKey(tsNewSymbol(tsStringUTF8("kFlushSync"))) })
	(*tsObject)(stream.ref).set(tsInternalWorkerStdioSingleton.flushKey, pipe.ackListener)
	return stream
}

func tsInternalWorkerStdioEndedReadable(loop *tsLoop) tsValue {
	module := tsInternalStreamReadableModule(loop)
	constructor := tsGet(loop, module, tsInternalWorkerStdioKey("Readable"))
	options := tsNewObject()
	ended := false
	var stream tsValue
	options.set("read", tsFunctionValue(tsFunc(func(loop *tsLoop, _ ...tsValue) tsValue {
		if !ended { ended = true; tsInternalWorkerStdioCall(loop, stream, "push", tsNull) }
		return tsU
	})))
	stream = tsConstruct(constructor, tsObjectValue(options))
	return stream
}

func tsInternalWorkerStdioModule(loop *tsLoop) tsValue {
	tsInternalWorkerStdioSingleton.Do(func() {
		tsInternalWorkerStdioSingleton.flushKey = tsPropertyKey(tsNewSymbol(tsStringUTF8("kFlushSync")))
		module := tsNewObject()
		module.set("kFlushSync", tsKeyValue(tsInternalWorkerStdioSingleton.flushKey))
		module.set("makePortReadable", tsNamedNativeMethod("makePortReadable", 2, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue { return tsInternalWorkerStdioReadable(loop, tsArg(args, 0), tsTruthy(tsArg(args, 1))) }))
		module.set("makePortWritable", tsNamedNativeMethod("makePortWritable", 1, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue { return tsInternalWorkerStdioWritable(loop, tsArg(args, 0)) }))
		module.set("makeEndedReadable", tsNamedNativeMethod("makeEndedReadable", 0, func(loop *tsLoop, _ tsValue, _ ...tsValue) tsValue { return tsInternalWorkerStdioEndedReadable(loop) }))
		module.set("default", tsObjectValue(module))
		tsInternalWorkerStdioSingleton.module = tsObjectValue(module)
	})
	return tsInternalWorkerStdioSingleton.module
}
`
