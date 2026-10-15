package goemit

const fsBufferMethodsRuntime = `
func tsNodeBufferNumericMethodLegacy(value tsValue,name string)tsValue{
    write:=strings.HasPrefix(name,"write");if !write&&!strings.HasPrefix(name,"read"){return tsU};spec:=strings.TrimPrefix(strings.TrimPrefix(name,"read"),"write")
    bigNumber:=strings.HasPrefix(spec,"Big");spec=strings.TrimPrefix(spec,"Big");spec=strings.ReplaceAll(spec,"Uint","UInt")
    little:=strings.HasSuffix(spec,"LE")||!strings.HasSuffix(spec,"BE");spec=strings.TrimSuffix(strings.TrimSuffix(spec,"LE"),"BE")
    width:=0;floating:=false;signed:=strings.HasPrefix(spec,"Int")
    switch spec{case "UInt","Int":width=-1;case "UInt8","Int8":width=1;case "UInt16","Int16":width=2;case "UInt32","Int32":width=4;case "UInt64","Int64":if !bigNumber{return tsU};width=8;case "Float":width=4;floating=true;case "Double":width=8;floating=true;default:return tsU}
    if bigNumber&&width!=8{return tsU};if width<0&&!strings.HasSuffix(name,"LE")&&!strings.HasSuffix(name,"BE"){return tsU}
    return tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{
        data:=tsNodeBufferVisibleBytes(value);offsetIndex:=0;if write{offsetIndex=1};offset:=int(tsFSInteger(tsArg(args,offsetIndex),"offset",0,0));size:=width
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
        if !bigNumber{number:=tsNumber(tsPrimitiveHint(loop,item,false));if floating{if size==4{bits=uint64(math.Float32bits(float32(number)))}else{bits=math.Float64bits(number)}}else{
            maximum:=math.Ldexp(1,size*8)-1;minimum:=float64(0);if signed{maximum=math.Ldexp(1,size*8-1)-1;minimum=-math.Ldexp(1,size*8-1)}
            if number<minimum||number>maximum{panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","value is outside the integer bounds")})};if !math.IsNaN(number){bits=uint64(int64(math.Trunc(number)))}
        }}
        for i:=0;i<size;i++{index:=i;if !little{index=size-1-i};data[offset+index]=byte(bits>>uint(i*8))};return tsNumberValue(float64(offset+size))
    }))
}
func tsNodeBufferNumericMethod(value tsValue,name string)tsValue{
    write:=strings.HasPrefix(name,"write")
    if !write&&!strings.HasPrefix(name,"read"){return tsU}
    spec:=strings.TrimPrefix(strings.TrimPrefix(name,"read"),"write")
    bigNumber:=strings.HasPrefix(spec,"Big");spec=strings.TrimPrefix(spec,"Big");spec=strings.ReplaceAll(spec,"Uint","UInt")
    little:=strings.HasSuffix(spec,"LE")||!strings.HasSuffix(spec,"BE")
    spec=strings.TrimSuffix(strings.TrimSuffix(spec,"LE"),"BE")
    variableWidth:=spec=="UInt"||spec=="Int"
    width:=0;floating:=false;signed:=strings.HasPrefix(spec,"Int")
    switch spec{case "UInt","Int":width=-1;case "UInt8","Int8":width=1;case "UInt16","Int16":width=2;case "UInt32","Int32":width=4;case "UInt64","Int64":if !bigNumber{return tsU};width=8;case "Float":width=4;floating=true;case "Double":width=8;floating=true;default:return tsU}
    if bigNumber&&width!=8{return tsU};if variableWidth&&!strings.HasSuffix(name,"LE")&&!strings.HasSuffix(name,"BE"){return tsU}
    return tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{
        data:=tsNodeBufferVisibleBytes(value);offsetIndex:=0;if write{offsetIndex=1};size:=width
        if variableWidth{byteLengthValue:=tsArg(args,offsetIndex+1);if !tsIsNumeric(byteLengthValue){tsFSInvalid("byteLength","a number")};n:=tsNumber(byteLengthValue);if math.IsNaN(n)||math.IsInf(n,0)||math.Trunc(n)!=n||n<1||n>6{panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","byteLength must be between 1 and 6")})};size=int(n)}
        item:=tsArg(args,0);number:=float64(0);bits:=uint64(0)
        if write{if bigNumber{if item.kind!=tsBigIntKind{tsFSInvalid("value","a bigint")};integer:=(*big.Int)(item.ref);if signed{if !integer.IsInt64(){panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","value is outside int64 bounds")})};bits=uint64(integer.Int64())}else{if !integer.IsUint64(){panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","value is outside uint64 bounds")})};bits=integer.Uint64()}}else{number=tsNumber(tsPrimitiveHint(loop,item,false))}}
        offsetValue:=tsArg(args,offsetIndex);if !variableWidth&&tsIsUndefined(offsetValue){offsetValue=tsNumberValue(0)}
        checkOffsetType:=func(){if !tsIsNumeric(offsetValue){tsFSInvalid("offset","a number")}}
        valueInRange:=true
        if write&&!bigNumber&&!floating{maximum:=math.Ldexp(1,size*8)-1;minimum:=float64(0);if signed{maximum=math.Ldexp(1,size*8-1)-1;minimum=-math.Ldexp(1,size*8-1)};valueInRange=!(number<minimum||number>maximum)}
        if write&&size==1{checkOffsetType()};if write&&!valueInRange{panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","value is outside the integer bounds")})};if write&&size!=1{checkOffsetType()};if !write{checkOffsetType()}
        offset:=int(tsFSInteger(offsetValue,"offset",0,0));if offset>len(data)||size>len(data)-offset{panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","offset is outside the bounds of the Buffer")})}
        if !write{for i:=0;i<size;i++{index:=i;if !little{index=size-1-i};bits|=uint64(data[offset+index])<<uint(i*8)};if floating{if size==4{return tsNumberValue(float64(math.Float32frombits(uint32(bits))))};return tsNumberValue(math.Float64frombits(bits))};if bigNumber{if signed{return tsBigIntSigned(int64(bits))};return tsBigIntUnsigned(bits)};if signed{shift:=uint(64-size*8);return tsNumberValue(float64(int64(bits<<shift)>>shift))};return tsNumberValue(float64(bits))}
        if !bigNumber{if floating{if size==4{bits=uint64(math.Float32bits(float32(number)))}else{bits=math.Float64bits(number)}}else if variableWidth&&size>4{bits=uint64(tsNumberUint32(number))|uint64(tsNumberUint32(math.Floor(number/4294967296)))<<32}else{bits=uint64(tsNumberUint32(number))}}
        for i:=0;i<size;i++{index:=i;if !little{index=size-1-i};data[offset+index]=byte(bits>>uint(i*8))};return tsNumberValue(float64(offset+size))
    }))
}
func tsNodeBufferExtraGet(value tsValue,name string)tsValue{
    if method:=tsNodeBufferNumericMethod(value,name);!tsIsUndefined(method){return method}
    data:=tsNodeBufferVisibleBytes(value)
    switch name{
    case "parent":return tsArrayBufferValue((*tsTypedArray)(value.ref).buffer)
    case "offset":return tsNumberValue(float64((*tsTypedArray)(value.ref).offset))
    case "toLocaleString":return tsNodeBufferGet(value,"toString")
    case "toJSON":return tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{out:=tsNewObject();out.set("type",tsStringReference(tsStringUTF8("Buffer")));array:=&tsArray{values:make([]tsValue,len(data))};for i,b:=range data{array.values[i]=tsNumberValue(float64(b))};out.set("data",tsArrayValue(array));return tsObjectValue(out)}))
    case "inspect":return tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{
        limit:=len(data);if limit>50{limit=50};var out strings.Builder;out.WriteString("<Buffer")
        for _,item:=range data[:limit]{out.WriteByte(' ');out.WriteByte("0123456789abcdef"[item>>4]);out.WriteByte("0123456789abcdef"[item&15])}
        if len(data)>limit{remaining:=len(data)-limit;out.WriteString(" ... ");out.WriteString(strconv.Itoa(remaining));out.WriteString(" more byte");if remaining!=1{out.WriteByte('s')}}
        if len(data)==0{out.WriteByte(' ')};out.WriteByte('>');return tsStringReference(tsStringUTF8(out.String()))
    }))
    case "swap16","swap32","swap64":return tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{width:=2;if name=="swap32"{width=4}else if name=="swap64"{width=8};if len(data)%width!=0{panic(tsThrown{tsFSException("RangeError","ERR_INVALID_BUFFER_SIZE","Buffer size must be a multiple of the element width")})};for first:=0;first<len(data);first+=width{for i:=0;i<width/2;i++{data[first+i],data[first+width-1-i]=data[first+width-1-i],data[first+i]}};return value}))
    case "includes","indexOf","lastIndexOf":return tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{
        result:=tsNodeBufferSearch(loop,value,name,args);if name=="includes"{return tsBooleanValue(tsNumber(result)!=-1)};return result
    }))
    }
    // Bun and Node expose these encoding primitives on Buffer.prototype too.
    for _,encoding:=range []string{"ascii","base64","base64url","latin1","hex","ucs2","utf8"}{
        if name==encoding+"Slice"{codec:=encoding;return tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{resolved:=tsFSEncoding(tsStringReference(tsStringUTF8(codec)),"utf8");return tsNodeBufferEncodingSlice(loop,value,resolved,args)}))}
        if name==encoding+"Write"{codec:=encoding;return tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{
            resolved:=tsFSEncoding(tsStringReference(tsStringUTF8(codec)),"utf8")
            if codec=="utf8"||codec=="latin1"||codec=="ascii"{return tsNodeBufferStrictWrite(loop,value,resolved,args)}
            return tsNodeBufferRawWrite(loop,value,resolved,args)
        }))}
    }
    return tsU
}
func tsNodeBufferInstallStatics(out *tsObject){
    out.set("poolSize",tsNumberValue(8192))
    out.set("compare",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeBufferCompareStatic(args)})))
    out.set("isEncoding",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)(result tsValue){result=tsBooleanValue(false);defer func(){if recover()!=nil{result=tsBooleanValue(false)}}();if tsArg(args,0).kind!=tsStringKind{return result};encoding:=tsFSEncoding(tsArg(args,0),"");return tsBooleanValue(encoding!="buffer")})))
    out.set("copyBytesFrom",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeBufferCopyBytesFrom(loop,args)})))
    out.set("of",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{data:=make([]byte,len(args));for i,item:=range args{data[i]=byte(tsNumberUint32(tsNumber(tsPrimitiveHint(loop,item,false))))};return tsNodeBuffer(data)})))
    out.set("isAscii",tsFunctionValue(tsFunc(func(_ *tsLoop,args ...tsValue)tsValue{return tsNodeBufferIsAscii(tsArg(args,0))})))
    out.set("isUtf8",tsFunctionValue(tsFunc(func(_ *tsLoop,args ...tsValue)tsValue{return tsNodeBufferIsUtf8(tsArg(args,0))})))
    out.set("transcode",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeBufferTranscode(loop,args)})))
}
func tsNodeBufferStringIndex(value tsValue,length,fallback int)int{
    if tsIsUndefined(value){return fallback};number:=tsNumber(value);if math.IsNaN(number)||number<0{return 0};if number>float64(length){return length};return int(math.Trunc(number))
}
`
