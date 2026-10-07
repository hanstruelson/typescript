package goemit

const fsBufferMethodsRuntime = `
func tsNodeBufferNumericMethod(value tsValue,name string)tsValue{
    write:=strings.HasPrefix(name,"write");if !write&&!strings.HasPrefix(name,"read"){return tsU};spec:=strings.TrimPrefix(strings.TrimPrefix(name,"read"),"write")
    bigNumber:=strings.HasPrefix(spec,"Big");spec=strings.TrimPrefix(spec,"Big");spec=strings.ReplaceAll(spec,"Uint","UInt")
    little:=strings.HasSuffix(spec,"LE");spec=strings.TrimSuffix(strings.TrimSuffix(spec,"LE"),"BE")
    width:=0;floating:=false;signed:=strings.HasPrefix(spec,"Int")
    switch spec{case "UInt","Int":width=-1;case "UInt8","Int8":width=1;case "UInt16","Int16":width=2;case "UInt32","Int32":width=4;case "UInt64","Int64":if !bigNumber{return tsU};width=8;case "Float":width=4;floating=true;case "Double":width=8;floating=true;default:return tsU}
    if bigNumber&&width!=8{return tsU};if width!=1&&!strings.HasSuffix(name,"LE")&&!strings.HasSuffix(name,"BE"){return tsU}
    return tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{
        data:=tsFSBytes(value);offsetIndex:=0;if write{offsetIndex=1};offset:=int(tsFSInteger(tsArg(args,offsetIndex),"offset",0,0));size:=width
        if size<0{size=int(tsFSInteger(tsArg(args,offsetIndex+1),"byteLength",0,1));if size>6{panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","byteLength must be between 1 and 6")})}}
        if offset>len(data)||size>len(data)-offset{panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","offset is outside the bounds of the Buffer")})}
        if !write{
            bits:=uint64(0);for i:=0;i<size;i++{index:=i;if !little{index=size-1-i};bits|=uint64(data[offset+index])<<uint(i*8)}
            if floating{number:=float64(0);if size==4{number=float64(math.Float32frombits(uint32(bits)))}else{number=math.Float64frombits(bits)};return tsNumberValue(number)}
            if bigNumber{if signed{return tsBigIntSigned(int64(bits))};return tsBigIntUnsigned(bits)}
            if signed{shift:=uint(64-size*8);return tsNumberValue(float64(int64(bits<<shift)>>shift))};return tsNumberValue(float64(bits))
        }
        item:=tsArg(args,0);bits:=uint64(0)
        if bigNumber{if item.kind!=tsBigIntKind{tsFSInvalid("value","a bigint")};number:=(*big.Int)(item.ref);if signed{if !number.IsInt64(){panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","value is outside int64 bounds")})};bits=uint64(number.Int64())}else{if !number.IsUint64(){panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","value is outside uint64 bounds")})};bits=number.Uint64()}}
        if !bigNumber{number:=tsNumber(tsToPrimitive(item));if floating{if size==4{bits=uint64(math.Float32bits(float32(number)))}else{bits=math.Float64bits(number)}}else{
            maximum:=math.Ldexp(1,size*8)-1;minimum:=float64(0);if signed{maximum=math.Ldexp(1,size*8-1)-1;minimum=-math.Ldexp(1,size*8-1)}
            if number<minimum||number>maximum{panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","value is outside the integer bounds")})};if !math.IsNaN(number){bits=uint64(int64(math.Trunc(number)))}
        }}
        for i:=0;i<size;i++{index:=i;if !little{index=size-1-i};data[offset+index]=byte(bits>>uint(i*8))};return tsNumberValue(float64(offset+size))
    }))
}
func tsNodeBufferExtraGet(value tsValue,name string)tsValue{
    if method:=tsNodeBufferNumericMethod(value,name);!tsIsUndefined(method){return method}
    data:=tsFSBytes(value)
    switch name{
    case "parent":return tsArrayBufferValue((*tsTypedArray)(value.ref).buffer)
    case "offset":return tsNumberValue(float64((*tsTypedArray)(value.ref).offset))
    case "toLocaleString":return tsNodeBufferGet(value,"toString")
    case "toJSON":return tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{out:=tsNewObject();out.set("type",tsStringReference(tsStringUTF8("Buffer")));array:=&tsArray{values:make([]tsValue,len(data))};for i,b:=range data{array.values[i]=tsNumberValue(float64(b))};out.set("data",tsArrayValue(array));return tsObjectValue(out)}))
    case "swap16","swap32","swap64":return tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{width:=2;if name=="swap32"{width=4}else if name=="swap64"{width=8};if len(data)%width!=0{panic(tsThrown{tsFSException("RangeError","ERR_INVALID_BUFFER_SIZE","Buffer size must be a multiple of the element width")})};for first:=0;first<len(data);first+=width{for i:=0;i<width/2;i++{data[first+i],data[first+width-1-i]=data[first+width-1-i],data[first+i]}};return value}))
    case "includes","indexOf","lastIndexOf":return tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{
        item,position,encoding:=tsArg(args,0),tsArg(args,1),tsArg(args,2);if position.kind==tsStringKind{encoding=position;position=tsU}
        needle:=[]byte{};if tsIsNumeric(item){needle=[]byte{byte(tsNumberUint32(tsNumber(item)))}}else{needle=tsFSData(item,encoding)}
        offset:=float64(0);if name=="lastIndexOf"{offset=float64(len(data))};if !tsIsUndefined(position){offset=tsNumber(position);if math.IsNaN(offset){offset=0;if name=="lastIndexOf"{offset=float64(len(data))}}};if offset<0{offset+=float64(len(data))}
        result:=-1
        if name=="lastIndexOf"{if offset>=0{if offset>float64(len(data)){offset=float64(len(data))};limit:=int(offset)+len(needle);if limit>len(data){limit=len(data)};result=bytes.LastIndex(data[:limit],needle);if len(needle)==0{result=int(offset)}}}else{if offset<0{offset=0};if offset<=float64(len(data)){index:=bytes.Index(data[int(offset):],needle);if index>=0{result=int(offset)+index}}}
        if name=="includes"{return tsBooleanValue(result>=0)};return tsNumberValue(float64(result))
    }))
    }
    // Bun and Node expose these encoding primitives on Buffer.prototype too.
    for _,encoding:=range []string{"ascii","base64","base64url","latin1","hex","ucs2","utf8"}{
        if name==encoding+"Slice"{codec:=encoding;return tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsCall(loop,tsNodeBufferGet(value,"toString"),tsStringReference(tsStringUTF8(codec)),tsArg(args,0),tsArg(args,1))}))}
        if name==encoding+"Write"{codec:=encoding;return tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsCall(loop,tsNodeBufferGet(value,"write"),tsArg(args,0),tsArg(args,1),tsArg(args,2),tsStringReference(tsStringUTF8(codec)))}))}
    }
    return tsU
}
func tsNodeBufferInstallStatics(out *tsObject){
    out.set("poolSize",tsNumberValue(8192))
    out.set("compare",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNumberValue(float64(bytes.Compare(tsFSBytes(tsArg(args,0)),tsFSBytes(tsArg(args,1)))))})))
    out.set("isEncoding",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)(result tsValue){result=tsBooleanValue(false);defer func(){if recover()!=nil{result=tsBooleanValue(false)}}();if tsArg(args,0).kind!=tsStringKind{return result};encoding:=tsFSEncoding(tsArg(args,0),"");return tsBooleanValue(encoding!="buffer")})))
    out.set("copyBytesFrom",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{value:=tsArg(args,0);if value.kind!=tsTypedArrayKind{tsFSInvalid("view","a TypedArray")};array:=(*tsTypedArray)(value.ref);offset:=int(tsFSInteger(tsArg(args,1),"offset",0,0));if offset>array.length{offset=array.length};length:=int(tsFSInteger(tsArg(args,2),"length",int64(array.length-offset),0));if length>array.length-offset{length=array.length-offset};source:=tsFSBytes(value);return tsNodeBuffer(append([]byte{},source[offset*array.bytes:(offset+length)*array.bytes]...))})))
    out.set("of",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{data:=make([]byte,len(args));for i,item:=range args{data[i]=byte(tsNumberUint32(tsNumber(item)))};return tsNodeBuffer(data)})))
}
func tsNodeBufferStringIndex(value tsValue,length,fallback int)int{
    if tsIsUndefined(value){return fallback};number:=tsNumber(value);if math.IsNaN(number)||number<0{return 0};if number>float64(length){return length};return int(math.Trunc(number))
}
`
