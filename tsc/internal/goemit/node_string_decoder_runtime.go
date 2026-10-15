package goemit

// nodeStringDecoderRuntime translates the stateful portion of Bun's
// JSStringDecoder binding into a per-instance Go decoder state.
const nodeStringDecoderRuntime = `
type tsStringDecoderState struct {
    encoding string
    pending []byte
    lastNeed, lastTotal int
    lastChar [4]byte
}

var tsNodeStringDecoderSingleton struct {
    sync.Once
    module *tsObject
    constructor *tsFunction
    prototype *tsObject
}

func tsStringDecoderBytes(value tsValue) []byte {
    switch value.kind {
    case tsTypedArrayKind:
        array := (*tsTypedArray)(value.ref)
        if array.buffer.detached { tsArrayTypeFailure("Cannot use a detached TypedArray") }
        return array.buffer.data[array.offset:array.offset+array.length*array.bytes]
    case tsDataViewKind:
        view := (*tsDataView)(value.ref)
        if view.buffer.detached { tsArrayTypeFailure("Cannot use a detached DataView") }
        return view.buffer.data[view.offset:view.offset+view.length]
    default:
        tsFSInvalid("buf", "a Buffer, TypedArray, or DataView")
    }
    return nil
}

func tsStringDecoderStateFor(receiver tsValue, method string) *tsStringDecoderState {
    if receiver.kind == tsObjectKind {
        if state := (*tsObject)(receiver.ref).nativeStringDecoder; state != nil { return state }
    }
    panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"TypeError", message:"StringDecoder.prototype."+method+" called on an incompatible receiver"})})
}

func tsStringDecoderSetPending(state *tsStringDecoderState, data []byte, total int) {
    state.pending = append(state.pending[:0], data...)
    state.lastTotal = total
    state.lastNeed = total-len(data)
    copy(state.lastChar[:], data)
}

func tsStringDecoderUTF8Tail(data []byte) (int, int) {
    if len(data) == 0 { return 0, 0 }
    start := len(data)-1
    for start >= 0 && data[start]&0xc0 == 0x80 { start-- }
    if start < 0 { return 0, 0 }
    first := data[start]
    total := 0
    switch {
    case first >= 0xc2 && first <= 0xdf: total = 2
    case first >= 0xe0 && first <= 0xef: total = 3
    case first >= 0xf0 && first <= 0xf4: total = 4
    default: return 0, 0
    }
    available := len(data)-start
    if available >= total { return 0, 0 }
    for i:=start+1; i<len(data); i++ {
        if data[i]&0xc0 != 0x80 { return 0, 0 }
    }
    if available > 1 {
        second := data[start+1]
        if first == 0xe0 && second < 0xa0 || first == 0xed && second >= 0xa0 || first == 0xf0 && second < 0x90 || first == 0xf4 && second >= 0x90 { return 0, 0 }
    }
    return start, total
}

func tsStringDecoderUTF16Tail(data []byte) (int, int) {
    if len(data)%2 != 0 { return len(data)-1, 2 }
    if len(data) < 2 { return 0, 0 }
    last := binary.LittleEndian.Uint16(data[len(data)-2:])
    if last >= 0xd800 && last <= 0xdbff { return len(data)-2, 4 }
    return 0, 0
}

func tsStringDecoderDecoded(data []byte, encoding string) *tsString {
    value := tsFSDecode(data, encoding)
    if value.kind != tsStringKind { return tsStringUnits(nil) }
    return (*tsString)(value.ref)
}

func tsStringDecoderWriteBytes(state *tsStringDecoderState, input []byte) *tsString {
    if len(input) == 0 { return tsStringUnits(nil) }
    data := make([]byte, 0, len(state.pending)+len(input))
    data = append(data, state.pending...)
    data = append(data, input...)
    state.pending = state.pending[:0]
    state.lastNeed, state.lastTotal = 0, 0

    switch state.encoding {
    case "utf8":
        start, total := tsStringDecoderUTF8Tail(data)
        if total != 0 { tsStringDecoderSetPending(state, data[start:], total); data = data[:start] }
        return tsStringDecoderDecoded(data, "utf8")
    case "utf16le":
        start, total := tsStringDecoderUTF16Tail(data)
        if total != 0 { tsStringDecoderSetPending(state, data[start:], total); data = data[:start] }
        return tsStringDecoderDecoded(data, "utf16le")
    case "base64", "base64url":
        tail := len(data)%3
        if tail != 0 { tsStringDecoderSetPending(state, data[len(data)-tail:], 3); data = data[:len(data)-tail] }
        return tsStringDecoderDecoded(data, state.encoding)
    default:
        return tsStringDecoderDecoded(data, state.encoding)
    }
}

func tsStringDecoderResult(first, second *tsString) tsValue {
    if len(first.units) == 0 { return tsStringReference(second) }
    if len(second.units) == 0 { return tsStringReference(first) }
    return tsStringReference(tsStringConcat(first, second))
}

func tsStringDecoderWrite(state *tsStringDecoderState, input []byte) *tsString {
    return tsStringDecoderWriteBytes(state, input)
}

func tsStringDecoderEnd(state *tsStringDecoderState, input []byte) tsValue {
    defer func() { state.pending = nil; state.lastNeed = 0; state.lastTotal = 0 }()
    first := tsStringDecoderWriteBytes(state, input)
    tail := tsStringDecoderDecoded(state.pending, state.encoding)
    return tsStringDecoderResult(first, tail)
}

func tsStringDecoderCreate(loop *tsLoop, encodingValue tsValue) tsValue {
    encoding := "utf8"
    if !tsNullish(encodingValue) { encoding = tsNodeBufferParseEncoding(loop, encodingValue, "utf8", true) }
    state := &tsStringDecoderState{encoding:encoding}
    object := tsNewObject()
    object.prototype = tsObjectValue(tsNodeStringDecoderSingleton.prototype)
    object.nativeStringDecoder = state
    return tsObjectValue(object)
}

func tsNodeStringDecoderModule() tsValue {
    tsNodeStringDecoderSingleton.Do(func() {
        prototype := tsNewObject()
        prototype.prototype = tsObjectPrototype()
        constructor := tsFunc(func(loop *tsLoop, args ...tsValue) tsValue {
            return tsStringDecoderCreate(loop, tsArg(args, 0))
        })
        constructor.name = "StringDecoder"
        constructor.length = 1
        constructor.constructible = true
        constructor.construct = func(loop *tsLoop, args ...tsValue) tsValue {
            return tsStringDecoderCreate(loop, tsArg(args, 0))
        }
        constructor.receiverCall = func(loop *tsLoop, receiver tsValue, args ...tsValue) tsValue {
            created := tsStringDecoderCreate(loop, tsArg(args,0))
            if receiver.kind == tsObjectKind {
                object := (*tsObject)(receiver.ref)
                object.nativeStringDecoder = (*tsObject)((*tsObject)(created.ref)).nativeStringDecoder
                object.set("encoding",tsStringReference(tsStringUTF8(object.nativeStringDecoder.encoding)))
                if object.descriptors == nil { object.descriptors = map[string]*tsDescriptor{} }
                object.descriptors["encoding"] = &tsDescriptor{writable:true, configurable:true}
                return receiver
            }
            return created
        }
        constructor.properties = tsNewObject()
        constructor.properties.set("prototype", tsObjectValue(prototype))
        constructor.properties.descriptors = map[string]*tsDescriptor{"prototype":{writable:false}}
        prototype.set("constructor", tsFunctionValue(constructor))
        prototype.descriptors = map[string]*tsDescriptor{"constructor":{writable:true, configurable:true}}
        for _, entry := range []struct{name string; length int}{{"write",1},{"end",1},{"text",2}} {
            name := entry.name
            prototype.set(name, tsNamedNativeMethod(name, entry.length, func(loop *tsLoop, receiver tsValue, args ...tsValue) tsValue {
                state := tsStringDecoderStateFor(receiver, name)
                switch name {
                case "write":
                    if len(args) == 0 { tsFSInvalid("buf", "a Buffer, TypedArray, or DataView") }
                    value := tsArg(args,0)
                    if value.kind == tsStringKind { return value }
                    return tsStringReference(tsStringDecoderWrite(state, tsStringDecoderBytes(value)))
                case "end":
                    if len(args) == 0 { return tsStringDecoderEnd(state,nil) }
                    return tsStringDecoderEnd(state,tsStringDecoderBytes(tsArg(args,0)))
                default:
                    if len(args) < 2 { tsFSInvalid("offset", "a number") }
                    data := tsStringDecoderBytes(tsArg(args,0))
                    offset := int(tsArg(args,1).ToInt32())
                    if offset < 0 || offset > len(data) { return tsStringReference(tsStringUnits(nil)) }
                    return tsStringReference(tsStringDecoderWrite(state,data[offset:]))
                }
            }))
            prototype.descriptors[name] = &tsDescriptor{writable:true, configurable:true}
        }
        for _, property := range []string{"lastChar","lastNeed","lastTotal","encoding"} {
            name := property
            prototype.descriptors[name] = &tsDescriptor{accessor:true, configurable:true, getter:tsNativeMethod(func(_ *tsLoop, receiver tsValue, _ ...tsValue) tsValue {
                state := tsStringDecoderStateFor(receiver,name)
                switch name {
                case "lastChar": return tsNodeBuffer(append([]byte(nil),state.lastChar[:]...))
                case "lastNeed": return tsNumberValue(float64(state.lastNeed))
                case "lastTotal": return tsNumberValue(float64(state.lastTotal))
                default: return tsStringReference(tsStringUTF8(state.encoding))
                }
            })}
        }
        tag := tsPropertyKey(tsWellKnownSymbol("toStringTag"))
        prototype.set(tag, tsStringReference(tsStringUTF8("StringDecoder")))
        prototype.descriptors[tag] = &tsDescriptor{configurable:true}
        module := tsNewObject()
        module.set("StringDecoder",tsFunctionValue(constructor))
        tsNodeStringDecoderSingleton.prototype = prototype
        tsNodeStringDecoderSingleton.constructor = constructor
        tsNodeStringDecoderSingleton.module = module
    })
    return tsObjectValue(tsNodeStringDecoderSingleton.module)
}
`
