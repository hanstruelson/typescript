package goemit

// internalPrimordialsRuntime maps Bun's primordials bundle onto the generated
// runtime's isolated collection, iterator, typed-array, and Promise values.
const internalPrimordialsRuntime = `
var tsInternalPrimordials struct { sync.Once; module *tsObject }

func tsInternalPrimordialsUncurry(loop *tsLoop, target tsValue) tsValue {
	if target.kind != tsFunctionKind { tsPropertyFailure("uncurryThis target must be callable") }
	function := tsFunc(func(loop *tsLoop, args ...tsValue) tsValue {
		receiver := tsArg(args, 0)
		return tsCallReceiver(loop, target, receiver, tsRestArgs(args, 1).values...)
	})
	function.name = "uncurried " + (*tsFunction)(target.ref).name
	function.length = max(0, (*tsFunction)(target.ref).length-1)
	return tsFunctionValue(function)
}

func tsInternalPrimordialsIterator(loop *tsLoop, args ...tsValue) tsValue {
	value := tsArg(args, 0)
	return tsIteratorValueReference(tsIterate(value))
}

func tsInternalPrimordialsGetter(name string) tsValue {
	return tsNamedNativeMethod("get "+name, 0, func(_ *tsLoop, _ tsValue, args ...tsValue) tsValue {
		value := tsArg(args, 0)
		switch name {
		case "size":
			if value.kind == tsCollectionKind { return tsNumberValue(float64(len((*tsCollection)(value.ref).index))) }
		case "buffer":
			if value.kind == tsTypedArrayKind { return tsArrayBufferValue((*tsTypedArray)(value.ref).buffer) }
		case "byteLength":
			if value.kind == tsTypedArrayKind { array := (*tsTypedArray)(value.ref); return tsNumberValue(float64(array.length*array.bytes)) }
		case "byteOffset":
			if value.kind == tsTypedArrayKind { return tsNumberValue(float64((*tsTypedArray)(value.ref).offset)) }
		case "length":
			if value.kind == tsTypedArrayKind { return tsNumberValue(float64((*tsTypedArray)(value.ref).length)) }
		case "Symbol.toStringTag":
			if value.kind == tsTypedArrayKind { return tsStringReference(tsStringUTF8((*tsTypedArray)(value.ref).name)) }
		}
		tsPropertyFailure("Primordial getter called on incompatible receiver")
		return tsU
	})
}

func tsInternalPrimordialsPromiseAll(loop *tsLoop, values tsValue, mapper tsValue, returnArray bool) *tsPromise {
	length := int(tsToLength(tsGet(loop, values, tsStringReference(tsStringUTF8("length")))))
	result := loop.promise()
	var output *tsArray
	if returnArray { output = &tsArray{values: make([]tsValue, length)} }
	if length == 0 {
		if returnArray { result.settle(tsResult{tsArrayValue(output), false}) } else { result.settle(tsResult{tsU, false}) }
		return result
	}
	remaining := length
	for index := 0; index < length; index++ {
		item := tsGet(loop, values, tsStringReference(tsStringUTF8(strconv.Itoa(index))))
		if !tsNullish(mapper) { item = tsCall(loop, mapper, item, tsNumberValue(float64(index))) }
		at := index
		loop.await(item, func(settled tsResult) {
			if settled.rejected { result.settle(settled); return }
			if returnArray { output.values[at] = settled.value }
			remaining--
			if remaining == 0 {
				if returnArray { result.settle(tsResult{tsArrayValue(output), false}) } else { result.settle(tsResult{tsU, false}) }
			}
		})
	}
	return result
}

func tsInternalPrimordialsModule(loop *tsLoop) tsValue {
	tsInternalPrimordials.Do(func() {
		module := tsNewObject()
		module.set("uncurryThis", tsNamedNativeMethod("uncurryThis", 1, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
			return tsInternalPrimordialsUncurry(loop, tsArg(args, 0))
		}))
		module.set("SafeMap", tsBuiltinClass("Map"))
		module.set("SafeSet", tsBuiltinClass("Set"))
		// Weak collections share the generated runtime's collection storage;
		// object-key identity is retained, while weak reachability is not modeled.
		module.set("SafeWeakMap", tsBuiltinClass("Map"))
		module.set("SafeArrayIterator", tsFunc(func(loop *tsLoop, args ...tsValue) tsValue { return tsInternalPrimordialsIterator(loop, args...) }))
		module.set("SafeStringIterator", tsFunc(func(loop *tsLoop, args ...tsValue) tsValue { return tsInternalPrimordialsIterator(loop, args...) }))
		module.set("MapPrototypeGetSize", tsInternalPrimordialsGetter("size"))
		module.set("SetPrototypeGetSize", tsInternalPrimordialsGetter("size"))
		module.set("TypedArrayPrototypeGetBuffer", tsInternalPrimordialsGetter("buffer"))
		module.set("TypedArrayPrototypeGetByteLength", tsInternalPrimordialsGetter("byteLength"))
		module.set("TypedArrayPrototypeGetByteOffset", tsInternalPrimordialsGetter("byteOffset"))
		module.set("TypedArrayPrototypeGetLength", tsInternalPrimordialsGetter("length"))
		module.set("TypedArrayPrototypeGetSymbolToStringTag", tsInternalPrimordialsGetter("Symbol.toStringTag"))
		module.set("SafePromiseAllReturnVoid", tsNamedNativeMethod("SafePromiseAllReturnVoid", 2, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
			return tsPromiseValue(tsInternalPrimordialsPromiseAll(loop, tsArg(args, 0), tsArg(args, 1), false))
		}))
		module.set("SafePromiseAllReturnArrayLike", tsNamedNativeMethod("SafePromiseAllReturnArrayLike", 2, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
			return tsPromiseValue(tsInternalPrimordialsPromiseAll(loop, tsArg(args, 0), tsArg(args, 1), true))
		}))
		for _, name := range []string{"Uint8ClampedArray", "Uint8Array", "Uint16Array", "Uint32Array", "Int8Array", "Int16Array", "Int32Array", "Float16Array", "Float32Array", "Float64Array", "BigUint64Array", "BigInt64Array"} {
			module.set(name, tsBuiltinClass(name))
		}
		tsInternalPrimordials.module = module
	})
	return tsObjectValue(tsInternalPrimordials.module)
}
`
