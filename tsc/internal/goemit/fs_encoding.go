package goemit

const fsEncodingRuntime = `
func tsFSEncoding(options tsValue, fallback string) string {
    value := options
    if options.kind != tsStringKind { value = tsFSProperty(options,"encoding") }
    if tsNullish(value) { return fallback }
    if value.kind != tsStringKind { tsFSInvalid("encoding","a string") }
    name := strings.ToLower((*tsString)(value.ref).String())
    switch name {
    case "utf8","utf-8": return "utf8"
    case "utf16le","utf-16le","ucs2","ucs-2": return "utf16le"
    case "latin1","binary": return "latin1"
    case "ascii","base64","base64url","buffer": return name
    case "hex","hex\x00": return "hex"
    }
    panic(tsThrown{tsFSException("TypeError","ERR_UNKNOWN_ENCODING","Unknown encoding: "+name)})
}

func tsFSDecode(data []byte, encoding string) tsValue {
    switch encoding {
    case "buffer": return tsNodeBuffer(data)
    case "hex": return tsStringReference(tsStringUTF8(hex.EncodeToString(data)))
    case "base64": return tsStringReference(tsStringUTF8(base64.StdEncoding.EncodeToString(data)))
    case "base64url": return tsStringReference(tsStringUTF8(base64.RawURLEncoding.EncodeToString(data)))
    case "utf16le":
        units := make([]uint16,len(data)/2)
        for i:=range units {units[i]=binary.LittleEndian.Uint16(data[i*2:])}
        return tsStringReference(tsStringUnits(units))
    case "ascii","latin1":
        units:=make([]uint16,len(data))
        for i, c:=range data {if encoding=="ascii"{c&=127};units[i]=uint16(c)}
        return tsStringReference(tsStringUnits(units))
    default:
        var out strings.Builder
        for len(data)>0 {r,size:=utf8.DecodeRune(data);out.WriteRune(r);data=data[size:]}
        return tsStringReference(tsStringUTF8(out.String()))
    }
}

func tsFSEncode(text *tsString, encoding string) []byte {
    switch encoding {
    case "utf16le":
        data:=make([]byte,len(text.units)*2)
        for i,c:=range text.units {binary.LittleEndian.PutUint16(data[i*2:],c)}
        return data
    case "latin1","ascii":
        data:=make([]byte,len(text.units));for i,c:=range text.units{data[i]=byte(c)};return data
    case "hex":
        source:=text.String();data:=make([]byte,0,len(source)/2)
        for i:=0;i+1<len(source);i+=2 {pair,err:=hex.DecodeString(source[i:i+2]);if err!=nil{break};data=append(data,pair[0])};return data
    case "base64","base64url":
        source:=strings.Map(func(r rune)rune{if r=='-'{return '+'};if r=='_'{return '/'};if r>='A'&&r<='Z'||r>='a'&&r<='z'||r>='0'&&r<='9'||r=='+'||r=='/'||r=='='{return r};return -1},text.String())
        if at:=strings.IndexByte(source,'=');at>=0{source=source[:at]}
        if len(source)%4==1{source=source[:len(source)-1]}
        data,_:=base64.RawStdEncoding.DecodeString(source);return data
    default:return []byte(text.String())
    }
}

func tsFSBytes(value tsValue) []byte {
    if value.kind==tsTypedArrayKind {
        array:=(*tsTypedArray)(value.ref)
        return array.buffer.data[array.offset:array.offset+array.length*array.bytes]
    }
    if value.kind==tsArrayBufferKind {return (*tsArrayBuffer)(value.ref).data}
    tsFSInvalid("buffer","a Buffer or typed array")
    return nil
}

func tsFSData(value, options tsValue) []byte {
    if value.kind==tsStringKind{return tsFSEncode((*tsString)(value.ref),tsFSEncoding(options,"utf8"))}
    return tsFSBytes(value)
}

func tsNodeBuffer(data []byte) tsValue {
    array:=tsNativeArray(data,"Uint8Array",func(n float64)uint8{return uint8(tsNumberUint32(n))})
    array.nodeBuffer=true
    return tsTypedArrayValue(array)
}

func tsNodeBufferGet(value tsValue, name string) tsValue {
    array:=(*tsTypedArray)(value.ref)
    data:=tsNodeBufferVisibleBytes(value)
    switch name {
    case "toString":return tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeBufferToString(loop,value,args)}))
    case "equals":return tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{
        other:=tsArg(args,0);if other.kind!=tsTypedArrayKind&&other.kind!=tsDataViewKind {tsFSInvalid("otherBuffer","a Buffer or Uint8Array")}
        return tsBooleanValue(bytes.Equal(data,tsNodeBufferVisibleBytes(other)))
    }))
    case "compare":return tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeBufferCompare(loop,value,args)}))
    case "slice","subarray":return tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{
        first,last:=tsArrayBound(tsArg(args,0),len(data),0),tsArrayBound(tsArg(args,1),len(data),len(data))
        if last<first{last=first};view:=tsNodeBuffer(data[first:last]);target:=(*tsTypedArray)(view.ref)
        target.buffer=array.buffer;target.offset=array.offset+first;tsArrayBufferRegister(target.buffer,target);return view
    }))
    case "fill":return tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeBufferFill(loop,value,args)}))
    case "write":return tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeBufferWrite(loop,value,args)}))
    case "copy":return tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeBufferCopy(loop,value,args)}))
    }
    return tsNodeBufferExtraGet(value,name)
}

var tsNodeBufferSingleton struct{sync.Once;value *tsObject;class *tsClass}
func tsNodeBufferModule() tsValue {
    tsNodeBufferSingleton.Do(func(){
        out:=tsNewObject();tsNodeBufferSingleton.value=out;tsNodeBufferInstallStatics(out)
        from:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeBufferFrom(loop,args)});from.name="from";from.length=3;from.constructible=true;from.construct=from.call
        out.set("from",tsFunctionValue(from))
        out.set("alloc",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeBufferAllocate(loop,args,true)})))
        out.set("allocUnsafe",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeBufferAllocate(loop,args,false)})))
        out.set("allocUnsafeSlow",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeBufferAllocate(loop,args,false)})))
        out.set("isBuffer",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{value:=tsArg(args,0);return tsBooleanValue(value.kind==tsTypedArrayKind&&(*tsTypedArray)(value.ref).nodeBuffer)})))
        out.set("byteLength",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeBufferByteLength(loop,args)})))
        out.set("concat",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeBufferConcat(loop,args)})))
        class:=tsNativeClass("Buffer");tsNodeBufferSingleton.class=class;properties:=tsInstanceProperties(class.static);for _,name:=range out.order{properties.define(name);properties.extra[name]=out.get(name)}
    })
    return tsClassValue(tsNodeBufferSingleton.class)
}
`
