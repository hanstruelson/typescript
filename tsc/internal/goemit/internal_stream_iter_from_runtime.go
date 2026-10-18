package goemit

// internalStreamIterFromRuntime translates Bun's byte stream normalization
// helpers. Sources are accepted through the same protocol symbols used by the
// Readable and iterator adapters.
const internalStreamIterFromRuntime = `
const tsNodeIterFromBatchSize = 128

func tsNodeIterFromTypes(loop *tsLoop) tsValue { return tsNodeUtilTypesModule() }
func tsNodeIterFromIsType(value tsValue, name string) bool { return tsNodeUtilType(name, value) }
func tsNodeIterFromIsPrimitive(value tsValue) bool {
	return value.kind == tsStringKind || value.kind == tsArrayBufferKind || value.kind == tsTypedArrayKind || value.kind == tsDataViewKind
}
func tsNodeIterFromIsSyncIterable(loop *tsLoop, value tsValue) bool {
	if tsNullish(value) || value.kind == tsStringKind { return false }
	return tsGet(loop, value, tsWellKnownSymbol("iterator")).kind == tsFunctionKind
}
func tsNodeIterFromIsAsyncIterable(loop *tsLoop, value tsValue) bool {
	if tsNullish(value) { return false }
	return tsGet(loop, value, tsWellKnownSymbol("asyncIterator")).kind == tsFunctionKind
}
func tsNodeIterFromInvalid(name, expected string, value tsValue) {
	_ = value
	panic(tsThrown{tsFSException("TypeError", "ERR_INVALID_ARG_TYPE", "The \""+name+"\" argument must be "+expected)})
}
func tsNodeIterFromTypeList(names ...string) string {
	if len(names) == 0 { return "a valid streamable value" }
	if len(names) == 1 { return "of type "+names[0] }
	if len(names) == 2 { return "of type "+names[0]+" or "+names[1] }
	return "of type "+strings.Join(names[:len(names)-1], ", ")+", or "+names[len(names)-1]
}
func tsNodeIterFromPrimitiveToUint8Array(loop *tsLoop, value tsValue) tsValue {
	if value.kind == tsStringKind { return tsInternalStreamIterUtilsToUint8Array(loop, value) }
	if value.kind == tsArrayBufferKind { return tsNewTypedArray("Uint8Array", value) }
	if value.kind == tsTypedArrayKind {
		array := (*tsTypedArray)(value.ref)
		if array.name == "Uint8Array" { return value }
		return tsTypedArrayBufferExtended("Uint8Array", []tsValue{tsArrayBufferValue(array.buffer), tsNumberValue(float64(array.offset)), tsNumberValue(float64(array.length*array.bytes))})
	}
	if value.kind == tsDataViewKind {
		view := (*tsDataView)(value.ref)
		return tsTypedArrayBufferExtended("Uint8Array", []tsValue{tsArrayBufferValue(view.buffer), tsNumberValue(float64(view.offset)), tsNumberValue(float64(view.length))})
	}
	tsNodeIterFromInvalid("chunk", "a string, ArrayBuffer, or ArrayBufferView", value)
	return tsU
}
func tsNodeIterFromArrayBufferViewToUint8Array(loop *tsLoop, value tsValue) tsValue {
	if value.kind == tsTypedArrayKind || value.kind == tsDataViewKind { return tsNodeIterFromPrimitiveToUint8Array(loop, value) }
	tsNodeIterFromInvalid("chunk", "an ArrayBufferView", value)
	return tsU
}
func tsNodeIterFromIsUint8ArrayBatch(value tsValue) bool {
	if value.kind != tsArrayKind { return false }
	array := (*tsArray)(value.ref)
	if array.length() == 0 { return true }
	if !tsNodeIterFromIsType(array.at(0), "isUint8Array") { return false }
	if array.length() == 1 { return true }
	if !tsNodeIterFromIsType(array.at(array.length()-1), "isUint8Array") { return false }
	for i := 1; i < array.length()-1; i++ { if !tsNodeIterFromIsType(array.at(i), "isUint8Array") { return false } }
	return true
}
func tsNodeIterFromArray(values []tsValue) tsValue { return tsArrayValue(&tsArray{values: append([]tsValue{}, values...)}) }
func tsNodeIterFromBatchSlices(values []tsValue, bound bool) []tsValue {
	if len(values) == 0 { return nil }
	if !bound || len(values) <= tsNodeIterFromBatchSize { return []tsValue{tsNodeIterFromArray(values)} }
	out := make([]tsValue, 0, (len(values)+tsNodeIterFromBatchSize-1)/tsNodeIterFromBatchSize)
	for i := 0; i < len(values); i += tsNodeIterFromBatchSize { end := i+tsNodeIterFromBatchSize; if end > len(values) { end = len(values) }; out = append(out, tsNodeIterFromArray(values[i:end])) }
	return out
}
func tsNodeIterFromHasProtocol(loop *tsLoop, value, symbol tsValue) bool {
	return !tsNullish(value) && (value.kind == tsObjectKind || value.kind == tsInstanceKind) && tsGet(loop, value, symbol).kind == tsFunctionKind
}
func tsNodeIterFromSyncValue(loop *tsLoop, value tsValue, chunks *[]tsValue) {
	if tsNodeIterFromIsPrimitive(value) { *chunks = append(*chunks, tsNodeIterFromPrimitiveToUint8Array(loop, value)); return }
	types := tsInternalStreamIterTypesModule(loop)
	toStreamable := tsGet(loop, types, tsStringReference(tsStringUTF8("toStreamable")))
	if tsNodeIterFromHasProtocol(loop, value, toStreamable) { method := tsGet(loop, value, toStreamable); tsNodeIterFromSyncValue(loop, tsCallReceiver(loop, method, value), chunks); return }
	if value.kind == tsArrayKind { array := (*tsArray)(value.ref); for i := 0; i < array.length(); i++ { tsNodeIterFromSyncValue(loop, array.at(i), chunks) }; return }
	if tsNodeIterFromIsSyncIterable(loop, value) { iterator := tsIterate(value); for iterator.next() { tsNodeIterFromSyncValue(loop, iterator.value, chunks) }; return }
	tsNodeIterFromInvalid("value", "a string, ArrayBuffer, ArrayBufferView, Iterable, or toStreamable value", value)
}
func tsNodeIterFromSyncSource(loop *tsLoop, source tsValue) []tsValue {
	if !tsNodeIterFromIsSyncIterable(loop, source) { tsNodeIterFromInvalid("source", "an Iterable", source) }
	iterator := tsIterate(source)
	output := []tsValue{}
	batch := []tsValue{}
	flush := func() { if len(batch) > 0 { output = append(output, tsNodeIterFromArray(batch)); batch = nil } }
	for iterator.next() {
		value := iterator.value
		if tsNodeIterFromIsUint8ArrayBatch(value) {
			flush()
			items := (*tsArray)(value.ref).values
			output = append(output, tsNodeIterFromBatchSlices(items, true)...)
			continue
		}
		if tsNodeIterFromIsType(value, "isUint8Array") {
			batch = append(batch, value)
			if len(batch) == tsNodeIterFromBatchSize { flush() }
			continue
		}
		flush()
		chunks := []tsValue{}
		tsNodeIterFromSyncValue(loop, value, &chunks)
		output = append(output, tsNodeIterFromBatchSlices(chunks, true)...)
	}
	flush()
	return output
}
func tsNodeIterFromIteratorObject(loop *tsLoop,values []tsValue, async bool) tsValue {
	index := 0
	object := tsNewObject()
	if async {
		iterator := &tsAsyncIterator{}
		iterator.next = func(loop *tsLoop) *tsPromise { promise := loop.promise(); if index >= len(values) { promise.settle(tsResult{tsIteratorResult(tsU, true), false}); return promise }; value := values[index]; index++; promise.settle(tsResult{tsIteratorResult(value, false), false}); return promise }
		iterator.close = func(loop *tsLoop) *tsPromise { index = len(values); return loop.resolved(tsIteratorResult(tsU, true), false) }
		iterator.resume = func(loop *tsLoop, action string, value tsValue) *tsPromise { if action == "return" { return iterator.close(loop) }; return iterator.next(loop) }
		return tsAsyncIteratorObject(iterator)
	}
	next := tsNamedNativeMethod("next", 0, func(_ *tsLoop, _ tsValue, _ ...tsValue) tsValue { if index >= len(values) { return tsIteratorResult(tsU, true) }; value := values[index]; index++; return tsIteratorResult(value, false) })
	object.set("next", next)
	object.set(tsPropertyKey(tsWellKnownSymbol("iterator")), tsNamedNativeMethod("[Symbol.iterator]", 0, func(_ *tsLoop, receiver tsValue, _ ...tsValue) tsValue { return receiver }))
	return tsObjectValue(object)
}
func tsNodeIterFromIterable(loop *tsLoop,values []tsValue, async bool) tsValue {
	result := tsNodeIterFromIteratorObject(loop,values, async)
	if async { return result }
	return result
}
func tsNodeIterFromEmpty(loop *tsLoop,async bool) tsValue { return tsNodeIterFromIterable(loop,nil, async) }
func tsNodeIterFromArrayBatch(loop *tsLoop,input tsValue, async bool) tsValue {
	items := (*tsArray)(input.ref).values
	return tsNodeIterFromIterable(loop,tsNodeIterFromBatchSlices(items, true), async)
}
func tsNodeIterFromAsyncValue(loop *tsLoop, value tsValue, chunks []tsValue, done func([]tsValue, tsValue)) {
	if tsNodeUtilType("isPromise", value) { loop.await(value, func(result tsResult) { if result.rejected { done(chunks, result.value) } else { tsNodeIterFromAsyncValue(loop, result.value, chunks, done) } }); return }
	if tsNodeIterFromIsPrimitive(value) { chunks = append(chunks, tsNodeIterFromPrimitiveToUint8Array(loop, value)); done(chunks, tsNull); return }
	types := tsInternalStreamIterTypesModule(loop)
	toAsync := tsGet(loop, types, tsStringReference(tsStringUTF8("toAsyncStreamable")))
	if tsNodeIterFromHasProtocol(loop, value, toAsync) { result := tsCallReceiver(loop, tsGet(loop, value, toAsync), value); tsNodeIterFromAsyncValue(loop, result, chunks, done); return }
	toSync := tsGet(loop, types, tsStringReference(tsStringUTF8("toStreamable")))
	if tsNodeIterFromHasProtocol(loop, value, toSync) { tsNodeIterFromAsyncValue(loop, tsCallReceiver(loop, tsGet(loop, value, toSync), value), chunks, done); return }
	if value.kind == tsArrayKind {
		array := (*tsArray)(value.ref); var visit func(int); visit = func(index int) { if index >= array.length() { done(chunks, tsNull); return }; tsNodeIterFromAsyncValue(loop, array.at(index), chunks, func(next []tsValue, err tsValue) { if tsTruthy(err) { done(next, err) } else { chunks = next; visit(index+1) } }) }; visit(0); return
	}
	if tsNodeIterFromIsAsyncIterable(loop, value) {
		iterator := tsAsyncIterate(value); var visit func(); visit = func() { loop.await(tsPromiseValue(iterator.next(loop)), func(result tsResult) { if result.rejected { done(chunks, result.value); return }; if tsTruthy(tsGet(loop, result.value, tsStringReference(tsStringUTF8("done")))) { done(chunks, tsNull); return }; item := tsGet(loop, result.value, tsStringReference(tsStringUTF8("value"))); tsNodeIterFromAsyncValue(loop, item, chunks, func(next []tsValue, err tsValue) { if tsTruthy(err) { done(next, err) } else { chunks = next; visit() } }) }) }; visit(); return
	}
	if tsNodeIterFromIsSyncIterable(loop, value) { iterator := tsIterate(value); var visit func(); visit = func() { if !iterator.next() { done(chunks, tsNull); return }; tsNodeIterFromAsyncValue(loop, iterator.value, chunks, func(next []tsValue, err tsValue) { if tsTruthy(err) { done(next, err) } else { chunks = next; visit() } }) }; visit(); return }
	done(chunks, tsFSException("TypeError", "ERR_INVALID_ARG_TYPE", "The \"value\" argument must be a string, ArrayBuffer, ArrayBufferView, Iterable, AsyncIterable, toStreamable, or toAsyncStreamable value"))
}
type tsNodeIterFromAsyncSource struct { source tsValue; iterator *tsAsyncIterator; direct bool; index int; batches []tsValue; pending tsValue; pendingResolve bool; closed bool }
func tsNodeIterFromAsyncSourceValue(source tsValue) *tsNodeIterFromAsyncSource { return &tsNodeIterFromAsyncSource{source:source} }
func tsNodeIterFromAsyncSourceNext(loop *tsLoop, state *tsNodeIterFromAsyncSource) *tsPromise {
	promise := loop.promise()
	var start func()
	start = func() {
		if state.closed { promise.settle(tsResult{tsIteratorResult(tsU, true), false}); return }
		if state.pendingResolve {
			loop.await(state.pending, func(result tsResult) { if result.rejected { promise.settle(result); return }; state.pendingResolve = false; resolved := result.value; validatedKey := tsGet(loop, tsInternalStreamIterTypesModule(loop), tsStringReference(tsStringUTF8("kValidatedSource"))); if !tsNullish(resolved) && tsTruthy(tsGet(loop, resolved, validatedKey)) { state.source = resolved; state.direct = true; state.iterator = nil } else { nextSource := tsNodeIterFromFrom(loop, resolved); state.source = nextSource; state.direct = true; state.iterator = nil }; start() })
			return
		}
		if state.batches != nil { if state.index < len(state.batches) { value := state.batches[state.index]; state.index++; promise.settle(tsResult{tsIteratorResult(value, false), false}); return }; state.closed = true; promise.settle(tsResult{tsIteratorResult(tsU, true), false}); return }
		if state.iterator == nil { state.iterator = tsAsyncIterate(state.source) }
		loop.await(tsPromiseValue(state.iterator.next(loop)), func(result tsResult) {
			if result.rejected { promise.settle(result); return }
			if tsTruthy(tsGet(loop, result.value, tsStringReference(tsStringUTF8("done")))) { state.closed = true; promise.settle(tsResult{tsIteratorResult(tsU, true), false}); return }
			value := tsGet(loop, result.value, tsStringReference(tsStringUTF8("value")))
			if state.direct { promise.settle(tsResult{tsIteratorResult(value, false), false}); return }
			if tsNodeIterFromIsUint8ArrayBatch(value) { if (*tsArray)(value.ref).length() == 0 { start(); return }; promise.settle(tsResult{tsIteratorResult(value, false), false}); return }
			if tsNodeIterFromIsType(value, "isUint8Array") { promise.settle(tsResult{tsIteratorResult(tsNodeIterFromArray([]tsValue{value}), false), false}); return }
			tsNodeIterFromAsyncValue(loop, value, nil, func(chunks []tsValue, err tsValue) { if tsTruthy(err) { promise.settle(tsResult{err, true}); return }; if len(chunks) == 0 { start(); return }; promise.settle(tsResult{tsIteratorResult(tsNodeIterFromArray(chunks), false), false}) })
		})
	}
	start()
	return promise
}
func tsNodeIterFromAsyncIterator(source tsValue) tsValue {
	state := tsNodeIterFromAsyncSourceValue(source)
	iterator := &tsAsyncIterator{}
	iterator.next = func(loop *tsLoop) *tsPromise { return tsNodeIterFromAsyncSourceNext(loop, state) }
	iterator.close = func(loop *tsLoop) *tsPromise { state.closed = true; if state.iterator != nil { return state.iterator.close(loop) }; return loop.resolved(tsIteratorResult(tsU, true), false) }
	iterator.resume = func(loop *tsLoop, action string, value tsValue) *tsPromise { if action == "return" { return iterator.close(loop) }; return iterator.next(loop) }
	return tsAsyncIteratorObject(iterator)
}
func tsNodeIterFromNormalizeAsyncValue(value tsValue) tsValue {
	iterator := &tsAsyncIterator{}
	done := false
	loaded := false
	chunks := []tsValue{}
	index := 0
	iterator.next = func(loop *tsLoop) *tsPromise {
		promise := loop.promise()
		if done { promise.settle(tsResult{tsIteratorResult(tsU, true), false}); return promise }
		if loaded { if index >= len(chunks) { done = true; promise.settle(tsResult{tsIteratorResult(tsU, true), false}); return promise }; item := chunks[index]; index++; promise.settle(tsResult{tsIteratorResult(item, false), false}); return promise }
		loaded = true
		tsNodeIterFromAsyncValue(loop, value, nil, func(values []tsValue, err tsValue) { if tsTruthy(err) { done = true; promise.settle(tsResult{err, true}); return }; chunks = values; if len(chunks) == 0 { done = true; promise.settle(tsResult{tsIteratorResult(tsU, true), false}); return }; item := chunks[0]; index = 1; promise.settle(tsResult{tsIteratorResult(item, false), false}) })
		return promise
	}
	iterator.close = func(loop *tsLoop) *tsPromise { done = true; return loop.resolved(tsIteratorResult(tsU, true), false) }
	iterator.resume = func(loop *tsLoop, action string, value tsValue) *tsPromise { if action == "return" { return iterator.close(loop) }; return iterator.next(loop) }
	return tsAsyncIteratorObject(iterator)
}
func tsNodeIterFromFromSync(loop *tsLoop, input tsValue) tsValue {
	if tsNullish(input) { tsNodeIterFromInvalid("input", "a non-null value", input) }
	if tsNodeIterFromIsPrimitive(input) { return tsNodeIterFromIterable(loop,[]tsValue{tsNodeIterFromArray([]tsValue{tsNodeIterFromPrimitiveToUint8Array(loop, input)})}, false) }
	if input.kind == tsArrayKind { array := (*tsArray)(input.ref); if array.length() == 0 { return tsNodeIterFromEmpty(loop,false) }; if tsNodeIterFromIsType(array.at(0), "isUint8Array") { all := true; for _, item := range array.values { if !tsNodeIterFromIsType(item, "isUint8Array") { all = false; break } }; if all { return tsNodeIterFromArrayBatch(loop,input, false) } } }
	types := tsInternalStreamIterTypesModule(loop); toStreamable := tsGet(loop, types, tsStringReference(tsStringUTF8("toStreamable"))); if tsNodeIterFromHasProtocol(loop, input, toStreamable) { return tsNodeIterFromFromSync(loop, tsCallReceiver(loop, tsGet(loop, input, toStreamable), input)) }
	if tsNodeIterFromIsAsyncIterable(loop, input) { tsNodeIterFromInvalid("input", "a synchronous input (not AsyncIterable)", input) }
	if tsNodeUtilType("isPromise", input) || (input.kind == tsObjectKind && tsGet(loop, input, tsStringReference(tsStringUTF8("then"))).kind == tsFunctionKind) { tsNodeIterFromInvalid("input", "a synchronous input (not Promise)", input) }
	if !tsNodeIterFromIsSyncIterable(loop, input) { tsNodeIterFromInvalid("input", "a string, ArrayBuffer, ArrayBufferView, Iterable, or toStreamable value", input) }
	return tsNodeIterFromIterable(loop,tsNodeIterFromSyncSource(loop, input), false)
}
func tsNodeIterFromFrom(loop *tsLoop, input tsValue) tsValue {
	if tsNullish(input) { tsNodeIterFromInvalid("input", "a non-null value", input) }
	types := tsInternalStreamIterTypesModule(loop); validated := tsGet(loop, types, tsStringReference(tsStringUTF8("kValidatedSource"))); if tsTruthy(tsGet(loop, input, validated)) { return input }
	if tsNodeIterFromIsPrimitive(input) { return tsNodeIterFromIterable(loop,[]tsValue{tsNodeIterFromArray([]tsValue{tsNodeIterFromPrimitiveToUint8Array(loop, input)})}, true) }
	if input.kind == tsArrayKind { array := (*tsArray)(input.ref); if array.length() == 0 { return tsNodeIterFromEmpty(loop,true) }; if tsNodeIterFromIsType(array.at(0), "isUint8Array") { all := true; for _, item := range array.values { if !tsNodeIterFromIsType(item, "isUint8Array") { all = false; break } }; if all { return tsNodeIterFromArrayBatch(loop,input, true) } } }
	toAsync := tsGet(loop, types, tsStringReference(tsStringUTF8("toAsyncStreamable"))); if tsNodeIterFromHasProtocol(loop, input, toAsync) { result := tsCallReceiver(loop, tsGet(loop, input, toAsync), input); if !tsNullish(result) && tsTruthy(tsGet(loop, result, validated)) { return result }; return tsNodeIterFromDeferredAsyncResult(loop, result) }
	toSync := tsGet(loop, types, tsStringReference(tsStringUTF8("toStreamable"))); if tsNodeIterFromHasProtocol(loop, input, toSync) { return tsNodeIterFromFrom(loop, tsCallReceiver(loop, tsGet(loop, input, toSync), input)) }
	if !tsNodeIterFromIsSyncIterable(loop, input) && !tsNodeIterFromIsAsyncIterable(loop, input) { tsNodeIterFromInvalid("input", "a string, ArrayBuffer, ArrayBufferView, Iterable, AsyncIterable, toStreamable, or toAsyncStreamable value", input) }
	return tsNodeIterFromAsyncIterator(input)
}
func tsNodeIterFromDeferredAsyncResult(loop *tsLoop, result tsValue) tsValue {
	state := &tsNodeIterFromAsyncSource{pending: result, pendingResolve: true}
	iterator := &tsAsyncIterator{}
	iterator.next = func(loop *tsLoop) *tsPromise { return tsNodeIterFromAsyncSourceNext(loop, state) }
	iterator.close = func(loop *tsLoop) *tsPromise { state.closed = true; if state.iterator != nil { return state.iterator.close(loop) }; return loop.resolved(tsIteratorResult(tsU, true), false) }
	iterator.resume = func(loop *tsLoop, action string, value tsValue) *tsPromise { if action == "return" { return iterator.close(loop) }; return iterator.next(loop) }
	return tsAsyncIteratorObject(iterator)
}
func tsInternalStreamIterFromModule(loop *tsLoop) tsValue {
	module := tsNewObject()
	module.set("arrayBufferViewToUint8Array", tsNamedNativeMethod("arrayBufferViewToUint8Array", 1, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue { return tsNodeIterFromArrayBufferViewToUint8Array(loop, tsArg(args, 0)) }))
	module.set("from", tsNamedNativeMethod("from", 1, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue { return tsNodeIterFromFrom(loop, tsArg(args, 0)) }))
	module.set("fromSync", tsNamedNativeMethod("fromSync", 1, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue { return tsNodeIterFromFromSync(loop, tsArg(args, 0)) }))
	module.set("isAsyncIterable", tsNamedNativeMethod("isAsyncIterable", 1, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue { return tsBooleanValue(tsNodeIterFromIsAsyncIterable(loop, tsArg(args, 0))) }))
	module.set("isPrimitiveChunk", tsNamedNativeMethod("isPrimitiveChunk", 1, func(_ *tsLoop, _ tsValue, args ...tsValue) tsValue { return tsBooleanValue(tsNodeIterFromIsPrimitive(tsArg(args, 0))) }))
	module.set("isSyncIterable", tsNamedNativeMethod("isSyncIterable", 1, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue { return tsBooleanValue(tsNodeIterFromIsSyncIterable(loop, tsArg(args, 0))) }))
	module.set("isUint8ArrayBatch", tsNamedNativeMethod("isUint8ArrayBatch", 1, func(_ *tsLoop, _ tsValue, args ...tsValue) tsValue { return tsBooleanValue(tsNodeIterFromIsUint8ArrayBatch(tsArg(args, 0))) }))
	module.set("normalizeAsyncValue", tsNamedNativeMethod("normalizeAsyncValue", 1, func(_ *tsLoop, _ tsValue, args ...tsValue) tsValue { return tsNodeIterFromNormalizeAsyncValue(tsArg(args, 0)) }))
	module.set("default", tsObjectValue(module))
	return tsObjectValue(module)
}
`
