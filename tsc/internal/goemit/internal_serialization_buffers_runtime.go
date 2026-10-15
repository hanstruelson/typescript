package goemit

// internalSerializationBuffersRuntime preserves Buffer identity as an
// envelope around ordinary serializer payloads, matching Bun's small helper.
const internalSerializationBuffersRuntime = `
func tsInternalSerializationBuffersTag(value tsValue) tsValue {
	stack := []tsValue{value}
	visited := map[unsafe.Pointer]bool{}
	buffers := &tsArray{}
	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if current.kind == tsTypedArrayKind {
			array := (*tsTypedArray)(current.ref)
			if array.nodeBuffer { buffers.values = append(buffers.values, current) }
			continue
		}
		if current.kind != tsObjectKind && current.kind != tsArrayKind && current.kind != tsCollectionKind { continue }
		if visited[current.ref] { continue }
		visited[current.ref] = true
		switch current.kind {
		case tsCollectionKind:
			collection := (*tsCollection)(current.ref)
			for _, entry := range collection.entries { if entry.live { stack = append(stack, entry.key); if !collection.set { stack = append(stack, entry.value) } } }
		case tsArrayKind:
			array := (*tsArray)(current.ref)
			for index, entry := range array.values { if array.holes[index] { continue }; if descriptor := array.descriptors[strconv.Itoa(index)]; descriptor != nil && descriptor.accessor { continue }; stack = append(stack, entry) }
			for name, entry := range array.properties { if descriptor := array.descriptors[name]; descriptor != nil && (descriptor.accessor || !descriptor.enumerable) { continue }; stack = append(stack, entry) }
		case tsObjectKind:
			object := (*tsObject)(current.ref)
			for _, name := range object.order { descriptor := object.descriptors[name]; if descriptor != nil && (descriptor.accessor || !descriptor.enumerable) { continue }; if entry, ok := object.values[name]; ok { stack = append(stack, entry) } }
		}
	}
	if len(buffers.values) == 0 { return tsNull }
	return tsArrayValue(&tsArray{values: []tsValue{value, tsArrayValue(buffers)}})
}
func tsInternalSerializationBuffersRestore(envelope tsValue) tsValue {
	fail := func() { panic(tsThrown{tsFSException("Error", "ERR_BUFFER_ENCODING", "failed to parse serialized buffer envelope")}) }
	if envelope.kind != tsArrayKind { fail() }
	outer := (*tsArray)(envelope.ref)
	if outer.length() != 2 || outer.at(1).kind != tsArrayKind { fail() }
	buffers := (*tsArray)(outer.at(1).ref)
	for _, value := range buffers.values { if value.kind != tsTypedArrayKind { fail() }; array := (*tsTypedArray)(value.ref); if array.name != "Uint8Array" && !array.nodeBuffer { fail() } }
	for _, value := range buffers.values { (*tsTypedArray)(value.ref).nodeBuffer = true }
	return outer.at(0)
}
func tsInternalSerializationBuffersModule(_ *tsLoop) tsValue {
	module := tsNewObject()
	module.set("tagBuffers", tsNamedNativeMethod("tagBuffers", 1, func(_ *tsLoop, _ tsValue, args ...tsValue) tsValue { return tsInternalSerializationBuffersTag(tsArg(args, 0)) }))
	module.set("restoreBuffers", tsNamedNativeMethod("restoreBuffers", 1, func(_ *tsLoop, _ tsValue, args ...tsValue) tsValue { return tsInternalSerializationBuffersRestore(tsArg(args, 0)) }))
	module.set("default", tsObjectValue(module))
	return tsObjectValue(module)
}
`
