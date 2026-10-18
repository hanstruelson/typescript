package goemit

// nodeUtilRuntime supplies node:util's already-unblocked base APIs and keeps
// the module object ready for subsequent source-order additions.
const nodeUtilRuntime = `
var tsNodeUtilSingleton struct {
    sync.Once
    value *tsObject
    regExpPrototypeExec, safeMap, isDeepStrictEqual tsValue
    debugPattern string
    debuggers map[string]tsValue
    styles map[string]tsNodeUtilStyle
    hexStyles map[string]tsNodeUtilStyle
    hexStyleOrder []string
    aborted map[tsValue][]func()
}

type tsNodeUtilStyle struct { open, close string; keepClose bool }

func tsNodeUtilString(value string) tsValue { return tsStringReference(tsStringUTF8(value)) }
func tsNodeUtilText(loop *tsLoop, value tsValue) string { return tsStringValueWithLoop(loop, value).String() }
func tsNodeUtilProperty(loop *tsLoop, value tsValue, name string) tsValue { return tsGet(loop, value, tsNodeUtilString(name)) }
func tsNodeUtilRequire(loop *tsLoop, name string) tsValue { return tsNodeModuleRequire(loop, []tsValue{tsNodeUtilString(name)}) }
func tsNodeUtilIsObject(value tsValue) bool {
    switch value.kind { case tsObjectKind,tsInstanceKind,tsArrayKind,tsTypedArrayKind,tsDataViewKind,tsArrayBufferKind,tsCollectionKind,tsRegExpKind,tsErrorKind,tsPromiseKind,tsIteratorKind,tsNamespaceKind: return true }
    return false
}

func tsNodeUtilDebugMatches(pattern, set string) bool {
    set = strings.ToUpper(set)
    for _, part := range strings.Split(strings.ToUpper(pattern), ",") {
        if part == "" { continue }
        p, s := 0, 0
        star, mark := -1, 0
        for s < len(set) {
            if p < len(part) && (part[p] == set[s] || part[p] == '?') { p++; s++; continue }
            if p < len(part) && part[p] == '*' { star = p; p++; mark = s; continue }
            if star >= 0 { p = star + 1; mark++; s = mark; continue }
            break
        }
        for p < len(part) && part[p] == '*' { p++ }
        if p == len(part) && s == len(set) { return true }
    }
    return false
}

func tsNodeUtilInspectFormat(loop *tsLoop, args []tsValue) tsValue {
    inspect := tsInternalInspectModule(loop)
    format := tsNodeUtilProperty(loop, inspect, "format")
    return tsCallReceiver(loop, format, inspect, args...)
}

func tsNodeUtilConsoleWrite(loop *tsLoop, method string, args []tsValue) {
    console := tsNodeConsoleModule()
    fn := tsNodeUtilProperty(loop, console, method)
    if fn.kind == tsFunctionKind { tsCallReceiver(loop, fn, console, args...) }
}

func tsNodeUtilGetStyleCache(loop *tsLoop) map[string]tsNodeUtilStyle {
    if tsNodeUtilSingleton.styles != nil { return tsNodeUtilSingleton.styles }
    tsNodeUtilSingleton.styles = make(map[string]tsNodeUtilStyle)
    inspect := tsNodeUtilProperty(loop, tsInternalInspectModule(loop), "inspect")
    colors := tsNodeUtilProperty(loop, inspect, "colors")
    for _, name := range tsObjectKeys(tsObjectValue((*tsObject)(colors.ref))).values {
        key := tsText(name); codes := tsGet(loop, colors, name)
        if codes.kind != tsArrayKind || (*tsArray)(codes.ref).length() < 2 { continue }
        open, close := int(tsNumber((*tsArray)(codes.ref).at(0))), int(tsNumber((*tsArray)(codes.ref).at(1)))
        tsNodeUtilSingleton.styles[key] = tsNodeUtilStyle{open:"\x1b["+strconv.Itoa(open)+"m",close:"\x1b["+strconv.Itoa(close)+"m",keepClose:open==1||open==2}
    }
    return tsNodeUtilSingleton.styles
}

func tsNodeUtilReplaceClose(text string, style tsNodeUtilStyle) string {
    if style.close == "" { return text }
    var out strings.Builder
    start := 0
    for {
        index := strings.Index(text[start:], style.close)
        if index < 0 { break }
        index += start
        after := index + len(style.close)
        if after == len(text) { break }
        out.WriteString(text[start:index])
        if style.keepClose { out.WriteString(style.close) }
        out.WriteString(style.open)
        start = after
    }
    out.WriteString(text[start:])
    return out.String()
}

func tsNodeUtilHexStyle(loop *tsLoop, hexText string) (tsNodeUtilStyle, bool) {
    if style, ok := tsNodeUtilSingleton.hexStyles[hexText]; ok { return style, true }
    if len(hexText)!=4 && len(hexText)!=7 { return tsNodeUtilStyle{}, false }
    if hexText[0]!='#' { return tsNodeUtilStyle{}, false }
    raw:=hexText[1:]
    if len(raw)==3 { raw=string([]byte{raw[0],raw[0],raw[1],raw[1],raw[2],raw[2]}) }
    rgb,err:=hex.DecodeString(raw); if err!=nil || len(rgb)!=3 { return tsNodeUtilStyle{}, false }
    style:=tsNodeUtilStyle{open:fmt.Sprintf("\x1b[38;2;%d;%d;%dm",rgb[0],rgb[1],rgb[2]),close:"\x1b[39m"}
    if tsNodeUtilSingleton.hexStyles==nil { tsNodeUtilSingleton.hexStyles=make(map[string]tsNodeUtilStyle) }
    if len(tsNodeUtilSingleton.hexStyleOrder)>=256 { old:=tsNodeUtilSingleton.hexStyleOrder[0];tsNodeUtilSingleton.hexStyleOrder=tsNodeUtilSingleton.hexStyleOrder[1:];delete(tsNodeUtilSingleton.hexStyles,old) }
    tsNodeUtilSingleton.hexStyles[hexText]=style;tsNodeUtilSingleton.hexStyleOrder=append(tsNodeUtilSingleton.hexStyleOrder,hexText)
    return style,true
}

func tsNodeUtilStyleText(loop *tsLoop, format,text,options tsValue) tsValue {
    validateStream:=true
    if !tsIsUndefined(options) {
        if options.kind!=tsObjectKind && options.kind!=tsInstanceKind { tsInternalValidatorInvalidType("options","an object",options) }
        if value:=tsNodeUtilProperty(loop,options,"validateStream");!tsIsUndefined(value)&&!tsIsNull(value){if value.kind!=tsBooleanKind{tsInternalValidatorInvalidType("options.validateStream","a boolean",value)};validateStream=tsTruthy(value)}
    }
    if !validateStream && format.kind==tsStringKind && text.kind==tsStringKind {
        key:=tsText(format);plain:=tsText(text)
        if key=="none" { return text }
        if style,ok:=tsNodeUtilGetStyleCache(loop)[key];ok { return tsNodeUtilString(style.open+tsNodeUtilReplaceClose(plain,style)+style.close) }
        if strings.HasPrefix(key,"#") { if style,ok:=tsNodeUtilHexStyle(loop,key);ok{return tsNodeUtilString(style.open+tsNodeUtilReplaceClose(plain,style)+style.close)} }
    }
    if text.kind!=tsStringKind { tsInternalValidatorInvalidType("text","a string",text) }
    if !tsIsUndefined(options) && options.kind!=tsObjectKind && options.kind!=tsInstanceKind { tsInternalValidatorInvalidType("options","an object",options) }
    validateOption:=tsU
    if !tsIsUndefined(options) { validateOption=tsNodeUtilProperty(loop,options,"validateStream") }
    if !tsIsUndefined(validateOption)&&!tsIsNull(validateOption)&&validateOption.kind!=tsBooleanKind { tsInternalValidatorInvalidType("options.validateStream","a boolean",validateOption) }
    skip:=false
    if validateStream {
        stream:=tsU;if !tsIsUndefined(options){stream=tsNodeUtilProperty(loop,options,"stream")};if tsNullish(stream){stream=tsNodeUtilProperty(loop,tsNodeProcessModule(loop),"stdout")}
        utils:=tsNodeUtilRequire(loop,"internal/streams/utils")
        valid:=false
        for _,name:=range []string{"isReadableStream","isWritableStream","isNodeStream"}{fn:=tsNodeUtilProperty(loop,utils,name);if fn.kind==tsFunctionKind && tsTruthy(tsCallReceiver(loop,fn,utils,stream)){valid=true;break}}
        if !valid { panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_TYPE","The \"stream\" argument must be an instance of ReadableStream, WritableStream, or Stream")}) }
        colors:=tsNodeUtilRequire(loop,"internal/util/colors");should:=tsNodeUtilProperty(loop,colors,"shouldColorize");skip=!tsTruthy(tsCallReceiver(loop,should,colors,stream))
    }
    keys:=[]tsValue{format};if format.kind==tsArrayKind{a:=(*tsArray)(format.ref);keys=a.values}
    cache:=tsNodeUtilGetStyleCache(loop);opens,closes,plain:="","",tsText(text)
    for _,keyValue:=range keys {
        if keyValue.kind!=tsStringKind {
            inspect:=tsNodeUtilProperty(loop,tsInternalInspectModule(loop),"inspect")
            colors:=tsNodeUtilProperty(loop,inspect,"colors")
            tsInternalValidatorOneOf(keyValue,tsNodeUtilString("format"),tsArrayValue(tsObjectKeys(colors)))
            continue
        }
        key:=tsText(keyValue);if key=="none" {continue}
        var style tsNodeUtilStyle;var ok bool
        if strings.HasPrefix(key,"#") { style,ok=tsNodeUtilHexStyle(loop,key);if !ok{tsInternalValidatorInvalidValue("format","must be a valid hex color (#RGB or #RRGGBB)",keyValue)} } else { style,ok=cache[key] }
        if !ok { tsInternalValidatorInvalidValue("format","must be one of the available color styles",keyValue) }
        opens+=style.open;closes=style.close+closes;plain=tsNodeUtilReplaceClose(plain,style)
    }
    if skip {return text};return tsNodeUtilString(opens+plain+closes)
}

func tsNodeUtilSystemMessage(name string) string {
    messages:=map[string]string{"E2BIG":"argument list too long","EACCES":"permission denied","EADDRINUSE":"address already in use","EADDRNOTAVAIL":"address not available","EAFNOSUPPORT":"address family not supported","EAGAIN":"resource temporarily unavailable","EAI_ADDRFAMILY":"address family not supported","EAI_AGAIN":"temporary failure","EAI_BADFLAGS":"bad ai_flags value","EAI_BADHINTS":"invalid value for hints","EAI_CANCELED":"request canceled","EAI_FAIL":"permanent failure","EAI_FAMILY":"ai_family not supported","EAI_MEMORY":"out of memory","EAI_NODATA":"no address","EAI_NONAME":"unknown node or service","EAI_OVERFLOW":"argument buffer overflow","EAI_PROTOCOL":"resolved protocol is unknown","EAI_SERVICE":"service not available for socket type","EAI_SOCKTYPE":"socket type not supported","EALREADY":"connection already in progress","EBADF":"bad file descriptor","EBUSY":"resource busy or locked","ECANCELED":"operation canceled","ECHARSET":"invalid Unicode character","ECONNABORTED":"software caused connection abort","ECONNREFUSED":"connection refused","ECONNRESET":"connection reset by peer","EDESTADDRREQ":"destination address required","EEXIST":"file already exists","EFAULT":"bad address in system call argument","EFBIG":"file too large","EHOSTUNREACH":"host is unreachable","EINTR":"interrupted system call","EINVAL":"invalid argument","EIO":"i/o error","EISCONN":"socket is already connected","EISDIR":"illegal operation on a directory","ELOOP":"too many symbolic links encountered","EMFILE":"too many open files","EMSGSIZE":"message too long","ENAMETOOLONG":"name too long","ENETDOWN":"network is down","ENETUNREACH":"network is unreachable","ENFILE":"file table overflow","ENOBUFS":"no buffer space available","ENODEV":"no such device","ENOENT":"no such file or directory","ENOMEM":"not enough memory","ENONET":"machine is not on the network","ENOPROTOOPT":"protocol not available","ENOSPC":"no space left on device","ENOSYS":"function not implemented","ENOTCONN":"socket is not connected","ENOTDIR":"not a directory","ENOTEMPTY":"directory not empty","ENOTSOCK":"socket operation on non-socket","ENOTSUP":"operation not supported on socket","EOVERFLOW":"value too large for defined data type","EPERM":"operation not permitted","EPIPE":"broken pipe","EPROTO":"protocol error","EPROTONOSUPPORT":"protocol not supported","EPROTOTYPE":"protocol wrong type for socket","ERANGE":"result too large","EROFS":"read-only file system","ESHUTDOWN":"cannot send after transport endpoint shutdown","ESPIPE":"invalid seek","ESRCH":"no such process","ETIMEDOUT":"connection timed out","ETXTBSY":"text file is busy","EXDEV":"cross-device link not permitted","UNKNOWN":"unknown error","EOF":"end of file","ENXIO":"no such device or address","EMLINK":"too many links","EHOSTDOWN":"host is down","EREMOTEIO":"remote I/O error","ENOTTY":"inappropriate ioctl for device","EFTYPE":"inappropriate file type or format","EILSEQ":"illegal byte sequence","ESOCKTNOSUPPORT":"socket type not supported","ENODATA":"no data available","EUNATCH":"protocol driver not attached","ENOEXEC":"exec format error"}
    if text,ok:=messages[name];ok{return text};return name
}

func tsNodeUtilDebugLog(loop *tsLoop, set tsValue) tsValue {
    name:=strings.ToUpper(tsNodeUtilText(loop,set))
    if tsNodeUtilSingleton.debuggers==nil { tsNodeUtilSingleton.debuggers=make(map[string]tsValue) }
    if logger,ok:=tsNodeUtilSingleton.debuggers[name];ok{return logger}
    enabled:=tsNodeUtilDebugMatches(tsNodeUtilSingleton.debugPattern,name)
    if !enabled { logger:=tsNamedNativeMethod("debuglog."+name,0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsU});tsNodeUtilSingleton.debuggers[name]=logger;return logger }
    if name=="HTTP"||name=="HTTP2" { process:=tsNodeProcessModule(loop);warn:=tsNodeUtilProperty(loop,process,"emitWarning");tsCall(loop,warn,tsNodeUtilString("Setting the NODE_DEBUG environment variable to '"+strings.ToLower(name)+"' can expose sensitive data (such as passwords, tokens and authentication headers) in the resulting log.")) }
    pid:=tsNodeUtilProperty(loop,tsNodeProcessModule(loop),"pid")
    logger:=tsNamedNativeMethod("debuglog."+name,0,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{
        message:=tsNodeUtilInspectFormat(loop,args)
        text:=name+" "+tsText(pid)+": "+tsText(message)
        tsNodeUtilConsoleWrite(loop,"error",[]tsValue{tsNodeUtilString(text)})
        return tsU
    })
    tsNodeUtilSingleton.debuggers[name]=logger
    return logger
}

func tsNodeUtilCallbackifyError(loop *tsLoop, reason tsValue) tsValue {
    if tsTruthy(reason){return reason}
    failure:=tsNewError(loop,"Error",[]tsValue{tsNodeUtilString("Promise was rejected with a falsy value")})
    tsSet(loop,failure,tsNodeUtilString("reason"),reason)
    tsSet(loop,failure,tsNodeUtilString("code"),tsNodeUtilString("ERR_FALSY_VALUE_REJECTION"))
    return failure
}

func tsNodeUtilCallbackify(loop *tsLoop, original tsValue) tsValue {
    if original.kind!=tsFunctionKind { tsInternalValidatorInvalidType("original","a function",original) }
    source:=(*tsFunction)(original.ref)
    callbackified:=tsNamedNativeMethod(source.name+"Callbackified",source.length+1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{
        if len(args)==0||tsArg(args,len(args)-1).kind!=tsFunctionKind { tsInternalValidatorInvalidType("last argument","a function",tsArg(args,len(args)-1)) }
        callback:=tsArg(args,len(args)-1); callArgs:=append([]tsValue{},args[:len(args)-1]...)
        var result tsValue;var thrown any
        func(){defer func(){if recovered:=recover();recovered!=nil{thrown=recovered}}();result=tsCallReceiver(loop,original,receiver,callArgs...)}()
        promise:=(*tsPromise)(nil)
        if thrown!=nil { promise=loop.resolved(tsUnwrap(thrown),true) } else if result.kind==tsPromiseKind { promise=(*tsPromise)(result.ref) } else { promise=loop.resolved(result,false) }
        promise.then(func(settled tsResult){
            process:=tsNodeProcessModule(loop);nextTick:=tsNodeUtilProperty(loop,process,"nextTick")
            if settled.rejected { tsCall(loop,nextTick,tsNamedNativeMethod("callbackifyOnRejected",2,func(loop *tsLoop,_ tsValue,values ...tsValue)tsValue{reason:=tsNodeUtilCallbackifyError(loop,tsArg(values,0));cb:=tsArg(values,1);return tsCallReceiver(loop,cb,receiver,reason)}),settled.value,callback) } else { tsCall(loop,nextTick,callback,tsNull,settled.value) }
        })
        return tsU
    })
    fn:=(*tsFunction)(callbackified.ref)
    if fn.properties==nil{fn.properties=tsNewObject()}
    if source.name!=""{fn.properties.set("name",tsNodeUtilString(source.name+"Callbackified"))}
    if source.properties!=nil { for _,key:=range source.properties.order { if key=="name"||key=="length" {continue};fn.properties.set(key,source.properties.values[key]) } }
    tsSetPrototype(callbackified,tsPrototypeOf(original))
    return callbackified
}

func tsNodeUtilInherits(loop *tsLoop, ctor, superCtor tsValue) tsValue {
    if ctor.kind!=tsFunctionKind&&ctor.kind!=tsClassKind {tsInternalValidatorInvalidType("ctor","a function",ctor)}
    if superCtor.kind!=tsFunctionKind&&superCtor.kind!=tsClassKind {tsInternalValidatorInvalidType("superCtor","a function",superCtor)}
    prototype:=tsNodeUtilProperty(loop,superCtor,"prototype")
    if tsIsUndefined(prototype){tsInternalValidatorInvalidType("superCtor.prototype","an object",prototype)}
    tsSet(loop,ctor,tsNodeUtilString("super_"),superCtor)
    ownPrototype:=tsNodeUtilProperty(loop,ctor,"prototype")
    tsSetPrototype(ownPrototype,prototype)
    return tsU
}

func tsNodeUtilExtend(loop *tsLoop, origin, add tsValue) tsValue {
    if !tsTruthy(add)||!tsNodeUtilIsObject(add){return origin}
    keys:=tsObjectKeys(add)
    for index:=keys.length()-1;index>=0;index-- { key:=keys.at(index);tsSet(loop,origin,key,tsGet(loop,add,key)) }
    return origin
}

func tsNodeUtilGetSystemErrorMessage(loop *tsLoop,args []tsValue) tsValue {
    value:=tsArg(args,0);tsInternalValidatorNumber(value,tsNodeUtilString("err"),tsU,tsU,true)
    if tsNumber(value)>=0 {tsInternalValidatorOutOfRange("err","a negative integer",value)}
    name:=tsText(tsNodeUtilBindingInternalErrorName(value))
    if strings.HasPrefix(name,"Unknown system error"){return tsNodeUtilString(name)}
    return tsNodeUtilString(tsNodeUtilSystemMessage(name))
}

func tsNodeUtilErrnoException(loop *tsLoop,args []tsValue) tsValue {
    shared:=tsInternalSharedModule(loop);constructor:=tsNodeUtilProperty(loop,shared,"ErrnoException")
    return tsConstruct(constructor,tsArg(args,0),tsArg(args,1),tsArg(args,2))
}

func tsNodeUtilGetCallSites(loop *tsLoop,args []tsValue) tsValue {
    count:=tsArg(args,0);options:=tsArg(args,1)
    if options.kind==tsUndefinedKind { optionsArg:=tsArg(args,0);if optionsArg.kind==tsObjectKind||optionsArg.kind==tsInstanceKind{options=optionsArg;count=tsNumberValue(10)} }
    if tsIsUndefined(count){count=tsNumberValue(10)}
    if !tsIsUndefined(options) {
        if options.kind!=tsObjectKind&&options.kind!=tsInstanceKind{tsInternalValidatorInvalidType("options","an object",options)}
        sourceMap:=tsNodeUtilProperty(loop,options,"sourceMap");if !tsIsUndefined(sourceMap)&&sourceMap.kind!=tsBooleanKind{tsInternalValidatorInvalidType("options.sourceMap","a boolean",sourceMap)}
    }
    tsInternalValidatorNumber(count,tsNodeUtilString("frameCount"),tsNumberValue(1),tsNumberValue(200),true)
    pcs:=make([]uintptr,int(tsNumber(count))+8);n:=runtime.Callers(2,pcs);frames:=runtime.CallersFrames(pcs[:n]);out:=&tsArray{}
    for index:=0;index<int(tsNumber(count));index++ {frame,more:=frames.Next();if !more{break};record:=tsNewObject();record.set("functionName",tsNodeUtilString(frame.Function));record.set("scriptId",tsNodeUtilString("0"));record.set("scriptName",tsNodeUtilString(filepath.Base(frame.File)));record.set("lineNumber",tsNumberValue(float64(frame.Line)));record.set("columnNumber",tsNumberValue(1));record.set("column",tsNumberValue(1));out.values=append(out.values,tsObjectValue(record))}
    return tsArrayValue(out)
}

func tsNodeUtilConvertSignal(loop *tsLoop, signal tsValue) tsValue {
    osModule:=tsNodeOSModule();constants:=tsNodeUtilProperty(loop,osModule,"constants");signals:=tsNodeUtilProperty(loop,constants,"signals")
    choices:=tsArrayValue(tsObjectKeys(signals));tsInternalValidatorOneOf(signal,tsNodeUtilString("signalCode"),choices)
    return tsNumberValue(128+tsNumber(tsNodeUtilProperty(loop,signals,tsText(signal))))
}

func tsNodeUtilAborted(loop *tsLoop, signal, resource tsValue) tsValue {
    if signal.kind!=tsObjectKind||(*tsObject)(signal.ref).abort==nil {tsInternalValidatorInvalidType("signal","an AbortSignal",signal)}
    if tsIsPrimitiveValue(resource)||tsNullish(resource) {tsInternalValidatorInvalidType("resource","an object",resource)}
    state:=(*tsObject)(signal.ref).abort
    promise:=loop.promise()
    state.Lock();if state.aborted{state.Unlock();return tsPromiseValue(loop.resolved(tsU,false))};state.followers=append(state.followers,func(tsValue){promise.settle(tsResult{tsU,false})});state.Unlock()
    return tsPromiseValue(promise)
}

func tsNodeUtilTextCodec(loop *tsLoop, encoder bool) tsValue {
    name:="TextDecoder";if encoder{name="TextEncoder"}
    constructor:=&tsFunction{constructible:true,name:name,length:0,properties:tsNewObject()}
    constructor.construct=func(loop *tsLoop,args ...tsValue)tsValue{
        out:=tsNewObject()
        if encoder {
            out.set("encoding",tsNodeUtilString("utf-8"))
            out.set("encode",tsNamedNativeMethod("encode",1,func(loop *tsLoop,_ tsValue,values ...tsValue)tsValue{input:=tsArg(values,0);text:="";if !tsIsUndefined(input){text=tsNodeUtilText(loop,input)};return tsTypedArrayValue(tsNativeArray([]byte(text),"Uint8Array",func(value float64)uint8{return uint8(value)}))}))
            out.set("encodeInto",tsNamedNativeMethod("encodeInto",2,func(loop *tsLoop,_ tsValue,values ...tsValue)tsValue{
                input,destination:=tsArg(values,0),tsArg(values,1);text:="";if !tsIsUndefined(input){text=tsNodeUtilText(loop,input)}
                if destination.kind!=tsTypedArrayKind||(*tsTypedArray)(destination.ref).name!="Uint8Array"{tsInternalValidatorInvalidType("destination","a Uint8Array",destination)}
                bytes:=[]byte(text);array:=(*tsTypedArray)(destination.ref);written:=len(bytes);if written>array.length{written=array.length};for i:=0;i<written;i++{array.write(i,float64(bytes[i]))}
                out:=tsNewObject();out.set("read",tsNumberValue(float64(written)));out.set("written",tsNumberValue(float64(written)));return tsObjectValue(out)
            }))
        } else {
            label:=tsArg(args,0);encoding:="utf-8";if !tsIsUndefined(label){encoding=strings.ToLower(tsNodeUtilText(loop,label))}
            switch encoding {case "utf-8","utf8","unicode-1-1-utf-8":encoding="utf-8";case "utf-16le","utf-16","ucs-2","ucs2":encoding="utf-16le";case "latin1","iso-8859-1","ascii","us-ascii","windows-1252":encoding="windows-1252";default:panic(tsThrown{tsFSException("RangeError","ERR_ENCODING_NOT_SUPPORTED","The "+encoding+" encoding is not supported")})}
            fatal:=tsTruthy(tsArg(args,1));ignoreBOM:=tsTruthy(tsArg(args,2));out.set("encoding",tsNodeUtilString(encoding));out.set("fatal",tsBooleanValue(fatal));out.set("ignoreBOM",tsBooleanValue(ignoreBOM))
            out.set("decode",tsNamedNativeMethod("decode",1,func(loop *tsLoop,_ tsValue,values ...tsValue)tsValue{
                input:=tsArg(values,0);data:=[]byte{}
                if tsIsUndefined(input)||input.kind==tsNullKind{return tsNodeUtilString("")}
                switch input.kind {case tsArrayBufferKind:data=(*tsArrayBuffer)(input.ref).data;case tsTypedArrayKind:array:=(*tsTypedArray)(input.ref);if array.buffer.detached{tsArrayTypeFailure("Cannot decode a detached ArrayBufferView")};data=array.buffer.data[array.offset:array.offset+array.length*array.bytes];case tsDataViewKind:view:=(*tsDataView)(input.ref);if view.buffer.detached{tsArrayTypeFailure("Cannot decode a detached ArrayBufferView")};data=view.buffer.data[view.offset:view.offset+view.length];default:tsInternalValidatorInvalidType("input","an ArrayBuffer or view",input)}
                text:=""
                if encoding=="utf-16le"{units:=make([]uint16,0,len(data)/2);for i:=0;i+1<len(data);i+=2{units=append(units,uint16(data[i])|uint16(data[i+1])<<8)};text=string(utf16.Decode(units))}else if encoding=="windows-1252"{runes:=make([]rune,len(data));for i,b:=range data{runes[i]=rune(b);if encoding=="windows-1252"&&strings.Contains("\u0080\u0082\u0083\u0084\u0085\u0086\u0087\u0088\u0089\u008a\u008b\u008c\u008e\u0091\u0092\u0093\u0094\u0095\u0096\u0097\u0098\u0099\u009a\u009b\u009c\u009e\u009f",string(runes[i])){runes[i]=[]rune("€‚ƒ„…†‡ˆ‰Š‹ŒŽ‘’“”•–—˜™š›œžŸ")[strings.IndexRune("\u0080\u0082\u0083\u0084\u0085\u0086\u0087\u0088\u0089\u008a\u008b\u008c\u008e\u0091\u0092\u0093\u0094\u0095\u0096\u0097\u0098\u0099\u009a\u009b\u009c\u009e\u009f",runes[i])]}};text=string(runes)}else{if fatal&&!utf8.Valid(data){panic(tsThrown{tsFSException("TypeError","ERR_ENCODING_INVALID_ENCODED_DATA","The encoded data was not valid for encoding utf-8")})};text=string(data)}
                if !ignoreBOM&&strings.HasPrefix(text,"\ufeff"){text=strings.TrimPrefix(text,"\ufeff")};return tsNodeUtilString(text)
            }))
        }
        return tsObjectValue(out)
    }
    constructor.properties.set("name",tsNodeUtilString(name));constructor.properties.set("length",tsNumberValue(0));return tsFunctionValue(constructor)
}

func tsNodeUtilSetTraceSigInt(_ *tsLoop, enable tsValue) tsValue {
    if enable.kind!=tsBooleanKind{tsInternalValidatorInvalidType("enable","a boolean",enable)}
    return tsU
}

func tsInternalUtilComparisonsModule(loop *tsLoop) tsValue {
	module := tsNewObject()
	module.set("isDeepStrictEqual", tsNamedNativeMethod("isDeepStrictEqual", 3, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
		left, right := tsArg(args, 0), tsArg(args, 1)
		skipPrototype := tsTruthy(tsArg(args, 2))
		if !skipPrototype && !tsIsPrimitiveValue(left) && !tsIsPrimitiveValue(right) && !tsSameValue(tsPrototypeOf(left), tsPrototypeOf(right)) {
			return tsBooleanValue(false)
		}
		return tsBooleanValue(tsNodeAssertDeep(loop, left, right, true, map[tsNodeAssertPair]bool{}))
	}))
	return tsObjectValue(module)
}

func tsNodeUtilSystemErrorName(loop *tsLoop,args []tsValue) tsValue {
    value:=tsArg(args,0)
    tsInternalValidatorNumber(value,tsStringReference(tsStringUTF8("err")),tsU,tsU,true)
    if tsNumber(value)>=0{tsInternalValidatorOutOfRange("err","a negative integer",value)}
    name:=tsNodeUtilBindingInternalErrorName(value)
    if strings.HasPrefix(tsText(name),"Unknown system error"){panic(tsThrown{tsFSException("RangeError","ERR_UNKNOWN_SYSTEM_ERROR",tsText(name))})}
    return name
}
func tsNodeUtilSystemErrorMap(loop *tsLoop) tsValue {
    result:=tsNewCollection(false)
    for _,entry:=range tsNodeUtilSystemErrors(){pair:=&tsArray{};pair.values=append(pair.values,tsNodeUtilString(entry.name),tsNodeUtilString(tsNodeUtilSystemMessage(entry.name)));tsSet(loop,result,tsNumberValue(float64(entry.code)),tsArrayValue(pair))}
    return result
}
func tsNodeUtilModule(loop *tsLoop) tsValue {
    tsNodeUtilSingleton.Do(func(){
        types:=tsNodeUtilTypesModule()
        promisifyModule:=tsInternalPromisifyModule(loop)
        validators:=tsInternalValidatorsModule(loop)
        shared:=tsInternalSharedModule(loop)
        mime:=tsInternalMIMEModule(loop)
        deprecation:=tsInternalDeprecateModule(loop)
        bindings:=tsNodeUtilBindingModule(loop)
        primordials:=tsInternalPrimordialsModule(loop)
        uncurryThis:=tsGet(loop,primordials,tsStringReference(tsStringUTF8("uncurryThis")))
        tsNodeUtilSingleton.regExpPrototypeExec=tsCall(loop,uncurryThis,tsGet(loop,tsRegExpPrototype(),tsStringReference(tsStringUTF8("exec"))))
        tsNodeUtilSingleton.safeMap=tsGet(loop,primordials,tsStringReference(tsStringUTF8("SafeMap")))
        comparisons:=tsInternalUtilComparisonsModule(loop)
        tsNodeUtilSingleton.isDeepStrictEqual=tsGet(loop,comparisons,tsStringReference(tsStringUTF8("isDeepStrictEqual")))
        inspection:=tsInternalInspectModule(loop)
        module:=tsNewObject();module.set("types",types)
        module.set("promisify",tsGet(loop,promisifyModule,tsStringReference(tsStringUTF8("promisify"))))
        module.set("parseArgs",tsNodeParseArgsFunction())
        for _,name:=range []string{"inspect","format","formatWithOptions","getOwnNonIndexProperties","stripVTControlCharacters"}{module.set(name,tsGet(loop,inspection,tsStringReference(tsStringUTF8(name))))}
        tsNodeUtilSingleton.debuggers=make(map[string]tsValue)
        env:=tsGet(loop,tsNodeProcessModule(loop),tsStringReference(tsStringUTF8("env")))
        tsNodeUtilSingleton.debugPattern=tsText(tsGet(loop,env,tsStringReference(tsStringUTF8("NODE_DEBUG"))))
        module.set("debuglog",tsNamedNativeMethod("debuglog",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeUtilDebugLog(loop,tsArg(args,0))}))
        module.set("debug",module.values["debuglog"])
        module.set("styleText",tsNamedNativeMethod("styleText",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeUtilStyleText(loop,tsArg(args,0),tsArg(args,1),tsArg(args,2))}))
        module.set("_errnoException",tsNamedNativeMethod("_errnoException",3,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeUtilErrnoException(loop,args)}))
        module.set("_extend",tsNamedNativeMethod("_extend",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeUtilExtend(loop,tsArg(args,0),tsArg(args,1))}))
        module.set("callbackify",tsNamedNativeMethod("callbackify",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeUtilCallbackify(loop,tsArg(args,0))}))
        module.set("inherits",tsNamedNativeMethod("inherits",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeUtilInherits(loop,tsArg(args,0),tsArg(args,1))}))
        module.set("log",tsNamedNativeMethod("log",0,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{stamp:=time.Now().Format("02 Jan 15:04:05");message:=tsNodeUtilInspectFormat(loop,args);tsNodeUtilConsoleWrite(loop,"log",[]tsValue{tsNodeUtilString(stamp+" - "+tsText(message))});return tsU}))
        module.set("toUSVString",tsNamedNativeMethod("toUSVString",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{value:=tsArg(args,0);if value.kind==tsSymbolKind{tsPropertyFailure("Cannot convert a Symbol value to a string")};return tsStringWellFormed(tsStringValueWithLoop(loop,value),true)}))
        module.set("getCallSites",tsNamedNativeMethod("getCallSites",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeUtilGetCallSites(loop,args)}))
        module.set("convertProcessSignalToExitCode",tsNamedNativeMethod("convertProcessSignalToExitCode",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeUtilConvertSignal(loop,tsArg(args,0))}))
        module.set("getSystemErrorMessage",tsNamedNativeMethod("getSystemErrorMessage",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeUtilGetSystemErrorMessage(loop,args)}))
        module.set("aborted",tsNamedNativeMethod("aborted",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeUtilAborted(loop,tsArg(args,0),tsArg(args,1))}))
        module.set("setTraceSigInt",tsNamedNativeMethod("setTraceSigInt",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeUtilSetTraceSigInt(loop,tsArg(args,0))}))
        module.set("isDeepStrictEqual",tsNodeUtilSingleton.isDeepStrictEqual)
        for name,check:=range map[string]func(tsValue)bool{
            "isArray":func(v tsValue)bool{return v.kind==tsArrayKind},
            "isBoolean":func(v tsValue)bool{return v.kind==tsBooleanKind},
            "isBuffer":func(v tsValue)bool{return v.kind==tsTypedArrayKind&&(*tsTypedArray)(v.ref).nodeBuffer},
            "isNull":func(v tsValue)bool{return v.kind==tsNullKind},
            "isNullOrUndefined":func(v tsValue)bool{return tsNullish(v)},
            "isNumber":func(v tsValue)bool{return tsIsNumeric(v)},
            "isString":func(v tsValue)bool{return v.kind==tsStringKind},
            "isSymbol":func(v tsValue)bool{return v.kind==tsSymbolKind},
            "isUndefined":func(v tsValue)bool{return v.kind==tsUndefinedKind},
            "isRegExp":func(v tsValue)bool{return tsNodeUtilType("isRegExp",v)},
            "isObject":func(v tsValue)bool{return tsNodeUtilIsObject(v)},
            "isDate":func(v tsValue)bool{return tsNodeUtilType("isDate",v)},
            "isError":func(v tsValue)bool{return tsNodeUtilType("isNativeError",v)},
            "isFunction":func(v tsValue)bool{return v.kind==tsFunctionKind||v.kind==tsClassKind},
            "isPrimitive":func(v tsValue)bool{return v.kind==tsNullKind||v.kind==tsUndefinedKind||v.kind==tsBooleanKind||tsIsNumeric(v)||v.kind==tsStringKind||v.kind==tsSymbolKind},
        } {predicateName:=name;predicate:=check;module.set(predicateName,tsNamedNativeMethod(predicateName,1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsBooleanValue(predicate(tsArg(args,0)))}))}
        module.set("TextDecoder",tsNodeUtilTextCodec(loop,false));module.set("TextEncoder",tsNodeUtilTextCodec(loop,true))
        module.set("getSystemErrorName",tsNamedNativeMethod("getSystemErrorName",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeUtilSystemErrorName(loop,args)}))
        module.set("getSystemErrorMap",tsNamedNativeMethod("getSystemErrorMap",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNodeUtilSystemErrorMap()}))
        module.set("MIMEType",tsGet(loop,mime,tsStringReference(tsStringUTF8("MIMEType"))))
        module.set("MIMEParams",tsGet(loop,mime,tsStringReference(tsStringUTF8("MIMEParams"))))
        module.set("deprecate",tsGet(loop,deprecation,tsStringReference(tsStringUTF8("deprecate"))))
        module.set("validateString",tsGet(loop,validators,tsStringReference(tsStringUTF8("validateString"))))
        module.set("validateOneOf",tsGet(loop,validators,tsStringReference(tsStringUTF8("validateOneOf"))))
        module.set("validateBoolean",tsGet(loop,validators,tsStringReference(tsStringUTF8("validateBoolean"))))
        module.set("validateObject",tsGet(loop,validators,tsStringReference(tsStringUTF8("validateObject"))))
        module.set("validateInteger",tsGet(loop,validators,tsStringReference(tsStringUTF8("validateInteger"))))
        module.set("resistStopPropagation",tsGet(loop,shared,tsStringReference(tsStringUTF8("resistStopPropagation"))))
        module.set("ErrnoException",tsGet(loop,shared,tsStringReference(tsStringUTF8("ErrnoException"))))
        for _,name:=range []string{"internalErrorName","internalErrorEntries","etimedoutErrorCode","enobufsErrorCode","ecanceledErrorCode","uvTranslateSysError","extractedSplitNewLinesFastPathStringsOnly","normalizeEncoding","parseEnv"}{module.set(name,tsGet(loop,bindings,tsStringReference(tsStringUTF8(name))))}
        tsNodeUtilSingleton.value=module
    })
    return tsObjectValue(tsNodeUtilSingleton.value)
}
`
