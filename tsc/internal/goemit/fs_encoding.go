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
    case "ascii","hex","base64","base64url","buffer": return name
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
    data:=tsFSBytes(value)
    switch name {
    case "toString":return tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{
        encoding:=tsFSEncoding(tsArg(args,0),"utf8")
        first,last:=tsNodeBufferStringIndex(tsArg(args,1),len(data),0),tsNodeBufferStringIndex(tsArg(args,2),len(data),len(data))
        if last<first{last=first};return tsFSDecode(data[first:last],encoding)
    }))
    case "equals","compare":return tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{
        compared:=bytes.Compare(data,tsFSBytes(tsArg(args,0)))
        if name=="equals"{return tsBooleanValue(compared==0)}
        return tsNumberValue(float64(compared))
    }))
    case "slice","subarray":return tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{
        first,last:=tsArrayBound(tsArg(args,0),len(data),0),tsArrayBound(tsArg(args,1),len(data),len(data))
        if last<first{last=first};view:=tsNodeBuffer(data[first:last]);target:=(*tsTypedArray)(view.ref)
        target.buffer=array.buffer;target.offset=array.offset+first;tsArrayBufferRegister(target.buffer,target);return view
    }))
    case "fill":return tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{
        first,last:=tsArrayBound(tsArg(args,1),len(data),0),tsArrayBound(tsArg(args,2),len(data),len(data))
        item:=tsArg(args,0);pattern:=[]byte{}
        if tsIsNumeric(item){pattern=[]byte{byte(tsNumberUint32(tsNumber(item)))}}else{pattern=tsFSData(item,tsArg(args,3))}
        if len(pattern)==0{if item.kind==tsStringKind&&len((*tsString)(item.ref).units)==0{pattern=[]byte{0}}else{tsFSInvalid("value","a nonempty string, number, or Buffer")}}
        for i:=first;i<last;i++{data[i]=pattern[(i-first)%len(pattern)]};return value
    }))
    case "write":return tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{
        offset:=0;if tsArg(args,1).kind!=tsStringKind{offset=int(tsFSInteger(tsArg(args,1),"offset",0,0))};if offset>len(data){tsFSInteger(tsNumberValue(-1),"offset",0,0)}
        length:=len(data)-offset;encoding:=tsArg(args,3)
        if tsArg(args,1).kind==tsStringKind{offset=0;encoding=args[1]}
        if tsArg(args,2).kind==tsStringKind{encoding=args[2]}else if !tsIsUndefined(tsArg(args,2)){length=int(tsFSInteger(args[2],"length",int64(length),0))}
        if length>len(data)-offset{length=len(data)-offset}
        item:=tsArg(args,0);if item.kind!=tsStringKind{tsFSInvalid("string","a string")}
        source:=tsFSEncode((*tsString)(item.ref),tsFSEncoding(encoding,"utf8"));return tsNumberValue(float64(copy(data[offset:offset+length],source)))
    }))
    case "copy":return tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{
        target:=tsFSBytes(tsArg(args,0));offset:=int(tsFSInteger(tsArg(args,1),"targetStart",0,0));first:=int(tsFSInteger(tsArg(args,2),"sourceStart",0,0));last:=int(tsFSInteger(tsArg(args,3),"sourceEnd",int64(len(data)),0))
        if offset>len(target){offset=len(target)};if first>len(data){first=len(data)};if last>len(data){last=len(data)};if last<first{last=first}
        return tsNumberValue(float64(copy(target[offset:],data[first:last])))
    }))
    }
    return tsNodeBufferExtraGet(value,name)
}

var tsNodeBufferSingleton struct{sync.Once;value *tsObject;class *tsClass}
func tsNodeBufferModule() tsValue {
    tsNodeBufferSingleton.Do(func(){
        out:=tsNewObject();tsNodeBufferSingleton.value=out;tsNodeBufferInstallStatics(out)
        out.set("from",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{
            source:=tsArg(args,0)
            if source.kind==tsObjectKind{if tsFSProperty(source,"type").kind==tsStringKind&&tsText(tsFSProperty(source,"type"))=="Buffer"{source=tsFSProperty(source,"data")}else if method:=tsFSProperty(source,"valueOf");method.kind==tsFunctionKind{converted:=tsCall(loop,method);if !tsStrictEqual(converted,source){return tsCall(loop,out.get("from"),converted,tsArg(args,1))}}}
            if source.kind==tsObjectKind&&tsIsUndefined(tsFSProperty(source,"length")){tsFSInvalid("first argument","a string, Buffer, ArrayBuffer, or array-like object")}
            if source.kind==tsStringKind{return tsNodeBuffer(tsFSEncode((*tsString)(source.ref),tsFSEncoding(tsArg(args,1),"utf8")))}
            if source.kind==tsArrayBufferKind{
                data:=(*tsArrayBuffer)(source.ref).data;offset:=0;if tsArg(args,1).kind!=tsStringKind{offset=int(tsFSInteger(tsArg(args,1),"offset",0,0))};if offset>len(data){tsFSInvalid("offset","within the buffer")}
                length:=int(tsFSInteger(tsArg(args,2),"length",int64(len(data)-offset),0));if length>len(data)-offset{tsFSInvalid("length","within the buffer")}
                value:=tsNodeBuffer(data[offset:offset+length]);array:=(*tsTypedArray)(value.ref);array.buffer=(*tsArrayBuffer)(source.ref);array.offset=offset;tsArrayBufferRegister(array.buffer,array);return value
            }
            data:=[]byte{};it:=tsArrayLikeIterator(source);for it.next(){data=append(data,byte(tsNumberUint32(tsNumber(it.value))))};return tsNodeBuffer(data)
        })))
        allocate:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{if !tsIsNumeric(tsArg(args,0)){tsFSInvalid("size","a number")};length:=int(tsFSInteger(tsArg(args,0),"size",0,0));value:=tsNodeBuffer(make([]byte,length));if !tsIsUndefined(tsArg(args,1)){tsCall(loop,tsNodeBufferGet(value,"fill"),tsArg(args,1),tsNumberValue(0),tsNumberValue(float64(length)),tsArg(args,2))};return value}))
        out.set("alloc",allocate);out.set("allocUnsafe",allocate);out.set("allocUnsafeSlow",allocate)
        out.set("isBuffer",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{value:=tsArg(args,0);return tsBooleanValue(value.kind==tsTypedArrayKind&&(*tsTypedArray)(value.ref).nodeBuffer)})))
        out.set("byteLength",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNumberValue(float64(len(tsFSData(tsArg(args,0),tsArg(args,1)))))})))
        out.set("concat",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{data:=[]byte{};it:=tsIterate(tsArg(args,0));for it.next(){data=append(data,tsFSBytes(it.value)...)};if !tsIsUndefined(tsArg(args,1)){length:=int(tsFSInteger(args[1],"length",0,0));if length<len(data){data=data[:length]}else{data=append(data,make([]byte,length-len(data))...)}};return tsNodeBuffer(data)})))
        class:=tsNativeClass("Buffer");tsNodeBufferSingleton.class=class;properties:=tsInstanceProperties(class.static);for _,name:=range out.order{properties.define(name);properties.extra[name]=out.get(name)}
    })
    return tsClassValue(tsNodeBufferSingleton.class)
}
`
