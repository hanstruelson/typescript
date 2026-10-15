package goemit

// internalValidatorsRuntime translates the JavaScript validator facade and
// its NodeValidator.cpp native entry points into local runtime checks.
const internalValidatorsRuntime = `
var tsInternalValidators struct { sync.Once; module *tsObject }

func tsInternalValidatorName(value tsValue) string {
    if tsIsUndefined(value) { return "undefined" }
    return tsText(value)
}

func tsInternalValidatorInvalidType(name,expected string,value tsValue) {
    panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_TYPE","The \""+name+"\" argument must be "+expected)})
}

func tsInternalValidatorOutOfRange(name,rangeText string,value tsValue) {
    panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","The \""+name+"\" argument must be "+rangeText+". Received "+tsText(value))})
}

func tsInternalValidatorInvalidValue(name,message string,value tsValue) {
    panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_VALUE","The argument '"+name+"' "+message)})
}

func tsInternalValidatorIsURL(value tsValue) bool {
    if value.kind!=tsObjectKind{return false}
    object:=(*tsObject)(value.ref)
    return object.nativeClass!=nil&&object.nativeClass.builtin=="URL"
}

func tsInternalValidatorPath(loop *tsLoop,value tsValue) tsValue {
    if tsInternalValidatorIsURL(value) {
        return tsStringReference(tsStringUTF8(tsFSPath(value)))
    }
    if value.kind!=tsStringKind {tsInternalValidatorInvalidType("path","a string or URL",value)}
    text:=tsText(value)
    if strings.HasPrefix(text,"file:") {return tsStringReference(tsStringUTF8(tsInternalValidatorFileURLPath(text)))}
    return tsStringReference(tsStringUTF8(tsPathPosixResolve([]tsValue{value})))
}

func tsInternalValidatorFileURLPath(text string) string {
    parsed,err:=url.Parse(text)
    if err!=nil||parsed.Scheme!="file" {panic(tsThrown{tsFSException("TypeError","ERR_INVALID_URL","Invalid file URL")})}
    if parsed.Host!=""&&parsed.Host!="localhost" {panic(tsThrown{tsFSException("TypeError","ERR_INVALID_FILE_URL_HOST","File URL host must be empty or localhost")})}
    if strings.Contains(strings.ToLower(parsed.EscapedPath()),"%2f") {panic(tsThrown{tsFSException("TypeError","ERR_INVALID_FILE_URL_PATH","File URL path must not contain encoded slashes")})}
    path,err:=url.PathUnescape(parsed.EscapedPath())
    if err!=nil {panic(tsThrown{tsFSException("TypeError","ERR_INVALID_FILE_URL_PATH",err.Error())})}
    if strings.ContainsRune(path,0) {panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_VALUE","Path must not contain null bytes")})}
    return path
}

func tsInternalValidatorFsPath(loop *tsLoop,value tsValue,name string) tsValue {
    if tsInternalValidatorIsURL(value) { return tsStringReference(tsStringUTF8(tsFSPath(value))) }
    if value.kind==tsStringKind {
        if strings.IndexByte(tsText(value),0)>=0 {tsInternalValidatorInvalidValue(name,"must be a string, Uint8Array, or URL without null bytes",value)}
        return value
    }
    if value.kind==tsTypedArrayKind {
        array:=(*tsTypedArray)(value.ref)
        if array.buffer.detached{tsArrayTypeFailure("Cannot use a detached TypedArray")}
        data:=array.buffer.data[array.offset:array.offset+array.length*array.bytes]
        for _,item:=range data{if item==0{tsInternalValidatorInvalidValue(name,"must be a string, Uint8Array, or URL without null bytes",value)}}
        if array.nodeBuffer||array.name=="Uint8Array"{return value}
    }
    tsInternalValidatorInvalidType(name,"a string, Buffer, or URL",value)
    return tsU
}

func tsInternalValidatorHttpToken(value tsValue) bool {
    if value.kind!=tsStringKind{return false}
    text:=tsText(value);if text==""{return false}
    for _,char:=range text {
        if char>127||char<33{return false}
        if strings.ContainsRune("()<>@,;:\\\"/[]?={} \t",char){return false}
    }
    return true
}

func tsInternalValidatorLinkFormat(value tsValue,name string) {
    if value.kind!=tsStringKind {tsInternalValidatorInvalidValue(name,"must be an array or string of format \"</styles.css>; rel=preload; as=style\"",value)}
    text:=tsText(value)
    if strings.ContainsAny(text,"\r\n")||!strings.HasPrefix(text,"<")||strings.IndexByte(text,'>')<1 {tsInternalValidatorInvalidValue(name,"must be an array or string of format \"</styles.css>; rel=preload; as=style\"",value)}
    fail:=func(){tsInternalValidatorInvalidValue(name,"must be an array or string of format \"</styles.css>; rel=preload; as=style\"",value)}
    rest:=text[strings.IndexByte(text,'>')+1:]
    isSpace:=func(char byte)bool{return strings.ContainsRune(" \t\n\v\f\r",rune(char))}
    for len(rest)>0 {
        for len(rest)>0&&isSpace(rest[0]){rest=rest[1:]}
        if len(rest)==0||rest[0]!=';'{fail()}
        rest=rest[1:]
        for len(rest)>0&&isSpace(rest[0]){rest=rest[1:]}
        nameLength:=0
        for nameLength<len(rest)&&rest[nameLength]!=';'&&rest[nameLength]!='='&&rest[nameLength]!='"'&&!isSpace(rest[nameLength]){nameLength++}
        if nameLength==0{fail()};rest=rest[nameLength:]
        if len(rest)>0&&rest[0]=='=' {
            rest=rest[1:];quoted:=len(rest)>0&&rest[0]=='"';if quoted{rest=rest[1:]}
            valueLength:=0
            for valueLength<len(rest)&&rest[valueLength]!=';'&&rest[valueLength]!='"'&&!isSpace(rest[valueLength]){valueLength++}
            rest=rest[valueLength:]
            if quoted {if len(rest)==0||rest[0]!='"'{fail()};rest=rest[1:]} else if len(rest)>0&&rest[0]=='"'{fail()}
        }
    }
}

func tsInternalValidatorCheckLink(loop *tsLoop,value tsValue) tsValue {
    if value.kind==tsStringKind {tsInternalValidatorLinkFormat(value,"hints");return value}
    if value.kind==tsArrayKind {
        array:=(*tsArray)(value.ref);var out strings.Builder
        for index:=0;index<array.length();index++{item:=array.at(index);tsInternalValidatorLinkFormat(item,"hints");if index>0{out.WriteString(", ")};out.WriteString(tsText(item))}
        return tsStringReference(tsStringUTF8(out.String()))
    }
    tsInternalValidatorInvalidValue("hints","must be an array or string of format \"</styles.css>; rel=preload; as=style\"",value)
    return tsU
}

func tsInternalValidatorObject(value,name,options tsValue) {
    flags:=int32(0);if tsIsNumeric(options){flags=int32(tsNumber(options))}
    if value.kind==tsNullKind {if flags&1!=0{return};tsInternalValidatorInvalidType(tsInternalValidatorName(name),"an object",value)}
    if value.kind==tsFunctionKind||value.kind==tsClassKind {if flags&4!=0{return};tsInternalValidatorInvalidType(tsInternalValidatorName(name),"an object",value)}
    if value.kind==tsArrayKind&&flags&2==0 {tsInternalValidatorInvalidType(tsInternalValidatorName(name),"an object",value)}
    if value.kind!=tsObjectKind&&value.kind!=tsInstanceKind&&value.kind!=tsArrayKind&&value.kind!=tsFunctionKind&&value.kind!=tsClassKind {tsInternalValidatorInvalidType(tsInternalValidatorName(name),"an object",value)}
}

func tsInternalValidatorNumber(value,name,min,max tsValue,integer bool) {
    if !tsIsNumeric(value) {tsInternalValidatorInvalidType(tsInternalValidatorName(name),"a number",value)}
    number:=tsNumber(value)
    lower,upper:=math.Inf(-1),math.Inf(1)
    if !tsNullish(min){lower=min.ToNumber()};if !tsNullish(max){upper=max.ToNumber()}
    if integer {
        if tsNullish(min){lower=-9007199254740991};if tsNullish(max){upper=9007199254740991}
        if math.IsNaN(number)||math.IsInf(number,0)||math.Trunc(number)!=number {tsInternalValidatorOutOfRange(tsInternalValidatorName(name),"an integer",value)}
        if math.IsNaN(lower){lower=-9007199254740991};if math.IsNaN(upper){upper=9007199254740991};if upper<lower{upper=lower}
        if number<lower||number>upper {tsInternalValidatorOutOfRange(tsInternalValidatorName(name),"a number "+tsInternalValidatorRange(lower,upper),value)}
        return
    }
    if (number<lower||number>upper||((!tsNullish(min)||!tsNullish(max))&&math.IsNaN(number))) {
        if !tsNullish(min)&&!tsNullish(max){tsInternalValidatorOutOfRange(tsInternalValidatorName(name),"a number "+tsInternalValidatorRange(lower,upper),value)}
        if !tsNullish(min){tsInternalValidatorOutOfRange(tsInternalValidatorName(name),">= "+tsNumberText(lower),value)}
        if !tsNullish(max){tsInternalValidatorOutOfRange(tsInternalValidatorName(name),"<= "+tsNumberText(upper),value)}
        tsInternalValidatorOutOfRange(tsInternalValidatorName(name),"a number",value)
    }
}

func tsInternalValidatorRange(lower,upper float64) string {
    return ">= "+tsNumberText(lower)+" && <= "+tsNumberText(upper)
}

func tsInternalValidatorFinite(value,name tsValue) bool {
    if tsIsUndefined(value)||value.kind==tsNumberKind&&math.IsNaN(tsNumber(value)){return false}
    if !tsIsNumeric(value){tsInternalValidatorInvalidType(tsInternalValidatorName(name),"a number",value)}
    if math.IsInf(tsNumber(value),0){tsInternalValidatorOutOfRange(tsInternalValidatorName(name),"a finite number",value)}
    return true
}

func tsInternalValidatorCheckRanges(number,name,lower,upper,def tsValue) tsValue {
    if !tsInternalValidatorFinite(number,name){return def}
    value:=number.ToNumber();min,max:=lower.ToNumber(),upper.ToNumber()
    if value<min||value>max{tsInternalValidatorOutOfRange(tsInternalValidatorName(name),tsInternalValidatorRange(min,max),number)}
    return number
}

func tsInternalValidatorPort(port,name,allowZero tsValue) tsValue {
    if tsIsUndefined(name){name=tsStringReference(tsStringUTF8("Port"))};if tsIsUndefined(allowZero){allowZero=tsBooleanValue(true)}
    validType:=tsIsNumeric(port)||port.kind==tsStringKind
    value:=math.NaN();if validType {
        if port.kind==tsStringKind&&strings.TrimSpace(tsText(port))=="" {validType=false} else {value=port.ToNumber()}
    }
    if !validType||math.IsNaN(value)||math.IsInf(value,0)||math.Trunc(value)!=value||value<0||value>65535||value==0&&!tsTruthy(allowZero) {
        tsInternalValidatorOutOfRange(tsInternalValidatorName(name),"a valid port",port)
    }
    return port
}

func tsInternalValidatorAbortSignal(value,name tsValue) {
    if tsIsUndefined(value){return}
    if (value.kind!=tsObjectKind&&value.kind!=tsInstanceKind)||!tsHas(value,tsStringReference(tsStringUTF8("aborted"))) {tsInternalValidatorInvalidType(tsInternalValidatorName(name),"an AbortSignal",value)}
}

func tsInternalValidatorArray(value,name,minLength tsValue) {
    if tsIsUndefined(minLength){minLength=tsNumberValue(0)}
    if value.kind!=tsArrayKind{panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_TYPE","The \""+tsInternalValidatorName(name)+"\" argument must be an instance of Array")})}
    if float64((*tsArray)(value.ref).length())<minLength.ToNumber(){tsInternalValidatorInvalidValue(tsInternalValidatorName(name),"must be longer than "+tsNumberText(minLength.ToNumber()),value)}
}

func tsInternalValidatorInt32(value,name,min,max tsValue) {
    if !tsIsNumeric(value){tsInternalValidatorInvalidType(tsInternalValidatorName(name),"a number",value)}
    lower,upper:=float64(-2147483648),float64(2147483647);if !tsIsUndefined(min){lower=min.ToNumber()};if !tsIsUndefined(max){upper=max.ToNumber()}
    number:=tsNumber(value);if math.Trunc(number)!=number{tsInternalValidatorOutOfRange(tsInternalValidatorName(name),"an integer",value)};if number<lower||number>upper{tsInternalValidatorOutOfRange(tsInternalValidatorName(name),tsInternalValidatorRange(lower,upper),value)}
}

func tsInternalValidatorUint32(value,name,positive tsValue) {
    if !tsIsNumeric(value){tsInternalValidatorInvalidType(tsInternalValidatorName(name),"a number",value)}
    minimum:=float64(0);if tsTruthy(positive){minimum=1};number:=tsNumber(value)
    if math.Trunc(number)!=number{tsInternalValidatorOutOfRange(tsInternalValidatorName(name),"an integer",value)}
    if number<minimum||number>4294967295{tsInternalValidatorOutOfRange(tsInternalValidatorName(name),tsInternalValidatorRange(minimum,4294967295),value)}
}

func tsInternalValidatorEncoding(loop *tsLoop,data,encoding tsValue) {
    name:=tsNodeBufferParseEncoding(loop,encoding,"utf8",true)
    if name=="hex" {
        length:=0
        switch data.kind {case tsStringKind:length=len((*tsString)(data.ref).units);case tsTypedArrayKind:length=(*tsTypedArray)(data.ref).length;case tsDataViewKind:length=(*tsDataView)(data.ref).length;case tsArrayBufferKind:length=len((*tsArrayBuffer)(data.ref).data);default:length=int(tsInternalStreamUtilsGet(loop,data,"length").ToNumber())}
        if length%2!=0{tsInternalValidatorInvalidValue("encoding","is invalid for data of length "+strconv.Itoa(length),encoding)}
    }
}

func tsInternalValidatorBuffer(value,name tsValue) {
    if tsIsUndefined(value){return}
    if value.kind!=tsTypedArrayKind&&value.kind!=tsDataViewKind{if tsIsUndefined(name){name=tsStringReference(tsStringUTF8("buffer"))};tsInternalValidatorInvalidType(tsInternalValidatorName(name),"an instance of Buffer, TypedArray, or DataView",value)}
}

func tsInternalValidatorOneOf(value,name,choices tsValue) {
    if choices.kind!=tsArrayKind{tsInternalValidatorInvalidType("values","an Array",choices)}
    array:=(*tsArray)(choices.ref)
    for index:=0;index<array.length();index++ {candidate:=array.at(index);if tsSameValue(value,candidate)||tsIsNumeric(value)&&tsIsNumeric(candidate)&&tsNumber(value)==tsNumber(candidate){return}}
    tsInternalValidatorInvalidValue(tsInternalValidatorName(name),"must be one of the allowed values",value)
}

func tsInternalValidatorsModule(loop *tsLoop) tsValue {
    tsInternalValidators.Do(func(){
        module:=tsNewObject()
        module.set("kValidateObjectNone",tsNumberValue(0));module.set("kValidateObjectAllowNullable",tsNumberValue(1));module.set("kValidateObjectAllowArray",tsNumberValue(2));module.set("kValidateObjectAllowFunction",tsNumberValue(4))
        methods:=[]string{"validateObject","validateLinkHeaderValue","checkIsHttpToken","validateInteger","validateNumber","validateString","validateFiniteNumber","checkRangesOrGetDefault","validateFunction","validateBoolean","validatePort","validateAbortSignal","validateArray","validateInt32","validateUint32","validateEncoding","validateBuffer","validateOneOf","isUint8Array","getValidatedPath","getValidatedFsPath","throwIfNullBytesInFileName"}
        for _,name:=range methods {method:=name;module.set(name,tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{
            value,nameValue:=tsArg(args,0),tsArg(args,1)
            switch method {
            case "validateObject":tsInternalValidatorObject(value,nameValue,tsArg(args,2))
            case "validateLinkHeaderValue":return tsInternalValidatorCheckLink(loop,value)
            case "checkIsHttpToken":return tsBooleanValue(tsInternalValidatorHttpToken(value))
            case "validateInteger":tsInternalValidatorNumber(value,nameValue,tsArg(args,2),tsArg(args,3),true)
            case "validateNumber":tsInternalValidatorNumber(value,nameValue,tsArg(args,2),tsArg(args,3),false)
            case "validateString":if value.kind!=tsStringKind{tsInternalValidatorInvalidType(tsInternalValidatorName(nameValue),"a string",value)}
            case "validateFiniteNumber":return tsBooleanValue(tsInternalValidatorFinite(value,nameValue))
            case "checkRangesOrGetDefault":return tsInternalValidatorCheckRanges(value,nameValue,tsArg(args,2),tsArg(args,3),tsArg(args,4))
            case "validateFunction":if value.kind!=tsFunctionKind&&value.kind!=tsClassKind{tsInternalValidatorInvalidType(tsInternalValidatorName(nameValue),"a function",value)}
            case "validateBoolean":if value.kind!=tsBooleanKind{tsInternalValidatorInvalidType(tsInternalValidatorName(nameValue),"a boolean",value)}
            case "validatePort":return tsInternalValidatorPort(value,nameValue,tsArg(args,2))
            case "validateAbortSignal":tsInternalValidatorAbortSignal(value,nameValue)
            case "validateArray":tsInternalValidatorArray(value,nameValue,tsArg(args,2))
            case "validateInt32":tsInternalValidatorInt32(value,nameValue,tsArg(args,2),tsArg(args,3))
            case "validateUint32":tsInternalValidatorUint32(value,nameValue,tsArg(args,2))
            case "validateEncoding":tsInternalValidatorEncoding(loop,value,nameValue)
            case "validateBuffer":tsInternalValidatorBuffer(value,nameValue)
            case "validateOneOf":tsInternalValidatorOneOf(value,nameValue,tsArg(args,2))
            case "isUint8Array":return tsBooleanValue(value.kind==tsTypedArrayKind&&((*tsTypedArray)(value.ref).name=="Uint8Array"||(*tsTypedArray)(value.ref).nodeBuffer))
            case "getValidatedPath":return tsInternalValidatorPath(loop,value)
            case "getValidatedFsPath":prop:=tsArg(args,1);if tsIsUndefined(prop){prop=tsStringReference(tsStringUTF8("path"))};return tsInternalValidatorFsPath(loop,value,tsText(prop))
            case "throwIfNullBytesInFileName":if value.kind==tsStringKind&&strings.ContainsRune(tsText(value),0){tsInternalValidatorInvalidValue("path","must be a string without null bytes",value)}
            }
            return tsU
        }))}
        for _,name:=range []string{"validateLinkHeaderValue","validateString","validateFunction","validateBoolean","getValidatedPath","getValidatedFsPath","throwIfNullBytesInFileName"}{tsInternalSharedHideFromStack([]tsValue{module.get(name)})}
        tsInternalValidators.module=module
    })
    return tsObjectValue(tsInternalValidators.module)
}
`
