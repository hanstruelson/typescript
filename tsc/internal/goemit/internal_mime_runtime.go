package goemit

// internalMIMERuntime translates Bun's MIMEParams and MIMEType native classes
// into stateful runtime objects with the same public methods and accessors.
const internalMIMERuntime = `
type tsInternalMIMEParam struct { name,value string }
type tsInternalMIMEParamsState struct { entries []tsInternalMIMEParam; index map[string]int }
type tsInternalMIMETypeState struct { typeName,subtype string; params tsValue }
var tsInternalMIME struct {
    sync.Once
    module,paramsPrototype,mimeTypePrototype *tsObject
    paramsConstructor,mimeTypeConstructor *tsFunction
    params map[unsafe.Pointer]*tsInternalMIMEParamsState
    mimeTypes map[unsafe.Pointer]*tsInternalMIMETypeState
}

func tsInternalMIMEString(loop *tsLoop,value tsValue) string { return tsStringValueWithLoop(loop,value).String() }
func tsInternalMIMEWhitespace(char byte) bool { return char==' '||char=='\t'||char=='\r'||char=='\n' }
func tsInternalMIMETokenChar(char rune) bool {
    if char>='0'&&char<='9'||char>='A'&&char<='Z'||char>='a'&&char<='z'{return true}
    return char==96||strings.ContainsRune("!#$%&'*+-.^_|~",char)
}
func tsInternalMIMEInvalidToken(text string) int {
    index:=0
    for _,char:=range text { if !tsInternalMIMETokenChar(char){return index};index++ }
    return -1
}
func tsInternalMIMEQuotedChar(char rune) bool { return char=='\t'||char>=' '&&char<='~'||char>=0x80&&char<=0xff }
func tsInternalMIMEInvalidQuoted(text string) int {
    index:=0
    for _,char:=range text { if !tsInternalMIMEQuotedChar(char){return index};index++ }
    return -1
}
func tsInternalMIMEError(kind,text string,index int) {
    where:=""
    if index>=0 {where=" at index "+strconv.Itoa(index)}
    panic(tsThrown{tsFSException("TypeError","ERR_INVALID_MIME_SYNTAX","The MIME "+kind+" "+strconv.Quote(text)+" is invalid"+where)})
}
func tsInternalMIMEEncodeParam(value string) string {
    if value==""{return "\"\""}
    if tsInternalMIMEInvalidToken(value)<0{return value}
    return "\""+strings.NewReplacer("\\","\\\\","\"","\\\"").Replace(value)+"\""
}
func tsInternalMIMEParamsNew(paramsString string) *tsInternalMIMEParamsState {
    state:=&tsInternalMIMEParamsState{index:map[string]int{}}
    position:=0
    for position<len(paramsString) {
        for position<len(paramsString)&&(tsInternalMIMEWhitespace(paramsString[position])||paramsString[position]==';'){position++}
        if position>=len(paramsString){break}
        nameStart:=position
        for position<len(paramsString)&&paramsString[position]!=';'&&paramsString[position]!='='{position++}
        name:=strings.ToLower(paramsString[nameStart:position])
        if position>=len(paramsString)||paramsString[position]!='=' {
            for position<len(paramsString)&&paramsString[position]!=';'{position++};if position<len(paramsString){position++};continue
        }
        position++
        value:=""
        if position<len(paramsString)&&paramsString[position]=='"' {
            position++;start:=position;escaped:=false
            for position<len(paramsString) {char:=paramsString[position];if escaped{escaped=false}else if char=='\\'{escaped=true}else if char=='"'{break};position++}
            raw:=paramsString[start:position];value=strings.NewReplacer("\\\\","\\","\\\"","\"").Replace(raw)
            if position<len(paramsString)&&paramsString[position]=='"'{position++}else{for position<len(paramsString)&&paramsString[position]!=';'{position++}}
        } else {
            start:=position;for position<len(paramsString)&&paramsString[position]!=';'{position++};value=strings.TrimRight(paramsString[start:position]," \t\r\n")
            if value==""{if position<len(paramsString){position++};continue}
        }
        if name==""||tsInternalMIMEInvalidToken(name)>=0||tsInternalMIMEInvalidQuoted(value)>=0 {if position<len(paramsString)&&paramsString[position]==';'{position++};continue}
        if _,exists:=state.index[name];!exists {state.index[name]=len(state.entries);state.entries=append(state.entries,tsInternalMIMEParam{name,value})}
        if position<len(paramsString)&&paramsString[position]==';'{position++}
    }
    return state
}
func tsInternalMIMEParamsString(state *tsInternalMIMEParamsState) string {
    var out strings.Builder
    for index,item:=range state.entries {if index>0{out.WriteByte(';')};out.WriteString(item.name);out.WriteByte('=');out.WriteString(tsInternalMIMEEncodeParam(item.value))}
    return out.String()
}
func tsInternalMIMEParamsFor(value tsValue) *tsInternalMIMEParamsState {
    if value.kind!=tsObjectKind||tsInternalMIME.params[value.ref]==nil{tsPropertyFailure("Value of 'this' must be a MIMEParams instance")}
    return tsInternalMIME.params[value.ref]
}
func tsInternalMIMETypeFor(value tsValue) *tsInternalMIMETypeState {
    if value.kind!=tsObjectKind||tsInternalMIME.mimeTypes[value.ref]==nil{tsPropertyFailure("Value of 'this' must be a MIMEType instance")}
    return tsInternalMIME.mimeTypes[value.ref]
}
func tsInternalMIMEParamsObject(state *tsInternalMIMEParamsState) tsValue {
    object:=tsNewObject();object.nativeClass=&tsClass{builtin:"MIMEParams"};object.prototype=tsObjectValue(tsInternalMIME.paramsPrototype);tsInternalMIME.params[unsafe.Pointer(object)]=state
    return tsObjectValue(object)
}
func tsInternalMIMEParamsIterator(receiver tsValue,mode string) tsValue {
    state:=tsInternalMIMEParamsFor(receiver);index:=0
    iterator:=&tsIterator{pull:func()(tsValue,bool){if index>=len(state.entries){return tsU,false};item:=state.entries[index];index++;switch mode{case "keys":return tsStringReference(tsStringUTF8(item.name)),true;case "values":return tsStringReference(tsStringUTF8(item.value)),true;default:return tsArrayValue(&tsArray{values:[]tsValue{tsStringReference(tsStringUTF8(item.name)),tsStringReference(tsStringUTF8(item.value))}}),true}}}
    return tsIteratorValueReference(iterator)
}
func tsInternalMIMEParseType(loop *tsLoop,input string)(string,string,string) {
    position:=0;for position<len(input)&&tsInternalMIMEWhitespace(input[position]){position++}
    slash:=strings.IndexByte(input[position:],'/')
    if slash<0{tsInternalMIMEError("type",input,tsInternalMIMEInvalidToken(input[position:]));return "","",""}
    slash+=position;typeName:=input[position:slash]
    if typeName==""||tsInternalMIMEInvalidToken(typeName)>=0{tsInternalMIMEError("type",input,tsInternalMIMEInvalidToken(typeName));return "","",""}
    typeName=strings.ToLower(typeName);position=slash+1
    semi:=strings.IndexByte(input[position:],';');params:="";rawSubtype:=input[position:]
    if semi>=0{semi+=position;rawSubtype=input[position:semi];params=input[semi+1:]}
    rawSubtype=strings.TrimRight(rawSubtype," \t\r\n")
    if rawSubtype==""||tsInternalMIMEInvalidToken(rawSubtype)>=0{tsInternalMIMEError("subtype",input,tsInternalMIMEInvalidToken(rawSubtype));return "","",""}
    return typeName,strings.ToLower(rawSubtype),params
}
func tsInternalMIMETypeObject(typeName,subtype string,params tsValue) tsValue {
    object:=tsNewObject();object.nativeClass=&tsClass{builtin:"MIMEType"};object.prototype=tsObjectValue(tsInternalMIME.mimeTypePrototype);tsInternalMIME.mimeTypes[unsafe.Pointer(object)]=&tsInternalMIMETypeState{typeName,subtype,params};return tsObjectValue(object)
}
func tsInternalMIMEInit(loop *tsLoop) {
    tsInternalMIME.Do(func(){
        tsInternalMIME.params=map[unsafe.Pointer]*tsInternalMIMEParamsState{};tsInternalMIME.mimeTypes=map[unsafe.Pointer]*tsInternalMIMETypeState{}
        paramsPrototype:=tsNewObject();tsInternalMIME.paramsPrototype=paramsPrototype
        paramsConstructor:=tsFunc(func(_ *tsLoop,_ ...tsValue)tsValue{tsPropertyFailure("Class constructor MIMEParams cannot be invoked without 'new'");return tsU});paramsConstructor.name="MIMEParams";paramsConstructor.length=0;paramsConstructor.constructible=true
        paramsConstructor.construct=func(_ *tsLoop,_ ...tsValue)tsValue{return tsInternalMIMEParamsObject(tsInternalMIMEParamsNew(""))};paramsConstructor.properties=tsNewObject();paramsConstructor.properties.set("prototype",tsObjectValue(paramsPrototype));paramsConstructor.properties.descriptors=map[string]*tsDescriptor{"prototype":{}}
        paramsPrototype.set("constructor",tsFunctionValue(paramsConstructor));paramsPrototype.descriptors=map[string]*tsDescriptor{"constructor":{writable:true,configurable:true}}
        for _,method:=range []string{"get","has","set","delete","toString","entries","keys","values"}{name:=method;arity:=1;if name=="set"{arity=2};if name=="toString"||name=="entries"||name=="keys"||name=="values"{arity=0};paramsPrototype.set(name,tsNamedNativeMethod(name,arity,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{
            state:=tsInternalMIMEParamsFor(receiver);first:=tsInternalMIMEString(loop,tsArg(args,0));slot,found:=state.index[first]
            switch name {
            case "get":if !found{return tsNull};return tsStringReference(tsStringUTF8(state.entries[slot].value))
            case "has":return tsBooleanValue(found)
            case "set":value:=tsInternalMIMEString(loop,tsArg(args,1));if invalid:=tsInternalMIMEInvalidToken(first);first==""||invalid>=0{tsInternalMIMEError("parameter name",first,invalid)};if invalid:=tsInternalMIMEInvalidQuoted(value);invalid>=0{tsInternalMIMEError("parameter value",value,invalid)};if found{state.entries[slot].value=value}else{state.index[first]=len(state.entries);state.entries=append(state.entries,tsInternalMIMEParam{first,value})};return tsU
            case "delete":if found{state.entries=append(state.entries[:slot],state.entries[slot+1:]...);state.index=map[string]int{};for index,item:=range state.entries{state.index[item.name]=index}};return tsU
            case "toString":return tsStringReference(tsStringUTF8(tsInternalMIMEParamsString(state)))
            case "entries":return tsInternalMIMEParamsIterator(receiver,"entries")
            case "keys":return tsInternalMIMEParamsIterator(receiver,"keys")
            default:return tsInternalMIMEParamsIterator(receiver,"values")
            }
        }));paramsPrototype.descriptors[name]=&tsDescriptor{writable:true,configurable:true}}
        entries:=paramsPrototype.get("entries");paramsPrototype.set(tsPropertyKey(tsWellKnownSymbol("iterator")),entries);paramsPrototype.descriptors[tsPropertyKey(tsWellKnownSymbol("iterator"))]=&tsDescriptor{writable:true,configurable:true};paramsPrototype.set("toJSON",paramsPrototype.get("toString"));paramsPrototype.descriptors["toJSON"]=&tsDescriptor{writable:true,configurable:true}
        tag:=tsPropertyKey(tsWellKnownSymbol("toStringTag"));paramsPrototype.set(tag,tsStringReference(tsStringUTF8("MIMEParams")));paramsPrototype.descriptors[tag]=&tsDescriptor{configurable:true}
        tsInternalMIME.paramsConstructor=paramsConstructor

        mimePrototype:=tsNewObject();tsInternalMIME.mimeTypePrototype=mimePrototype
        mimeConstructor:=tsFunc(func(_ *tsLoop,_ ...tsValue)tsValue{tsPropertyFailure("Class constructor MIMEType cannot be invoked without 'new'");return tsU});mimeConstructor.name="MIMEType";mimeConstructor.length=1;mimeConstructor.constructible=true
        mimeConstructor.construct=func(loop *tsLoop,args ...tsValue)tsValue{input:=tsInternalMIMEString(loop,tsArg(args,0));typeName,subtype,params:=tsInternalMIMEParseType(loop,input);return tsInternalMIMETypeObject(typeName,subtype,tsInternalMIMEParamsObject(tsInternalMIMEParamsNew(params)))}
        mimeConstructor.properties=tsNewObject();mimeConstructor.properties.set("prototype",tsObjectValue(mimePrototype));mimeConstructor.properties.descriptors=map[string]*tsDescriptor{"prototype":{}}
        mimePrototype.set("constructor",tsFunctionValue(mimeConstructor));mimePrototype.descriptors=map[string]*tsDescriptor{"constructor":{writable:true,configurable:true}}
        mimePrototype.order=append(mimePrototype.order,"type","subtype","essence","params")
        mimePrototype.descriptors["type"]=&tsDescriptor{accessor:true,enumerable:true,configurable:true,getter:tsNativeMethod(func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{return tsStringReference(tsStringUTF8(tsInternalMIMETypeFor(receiver).typeName))}),setter:tsNativeMethod(func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsInternalMIMETypeFor(receiver);value:=tsInternalMIMEString(loop,tsArg(args,0));if invalid:=tsInternalMIMEInvalidToken(value);value==""||invalid>=0{tsInternalMIMEError("type",value,invalid)};state.typeName=strings.ToLower(value);return tsU})}
        mimePrototype.descriptors["subtype"]=&tsDescriptor{accessor:true,enumerable:true,configurable:true,getter:tsNativeMethod(func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{return tsStringReference(tsStringUTF8(tsInternalMIMETypeFor(receiver).subtype))}),setter:tsNativeMethod(func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsInternalMIMETypeFor(receiver);value:=tsInternalMIMEString(loop,tsArg(args,0));if invalid:=tsInternalMIMEInvalidToken(value);value==""||invalid>=0{tsInternalMIMEError("subtype",value,invalid)};state.subtype=strings.ToLower(value);return tsU})}
        mimePrototype.descriptors["essence"]=&tsDescriptor{accessor:true,enumerable:true,configurable:true,getter:tsNativeMethod(func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsInternalMIMETypeFor(receiver);return tsStringReference(tsStringUTF8(state.typeName+"/"+state.subtype))})}
        mimePrototype.descriptors["params"]=&tsDescriptor{accessor:true,enumerable:true,configurable:true,getter:tsNativeMethod(func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{return tsInternalMIMETypeFor(receiver).params})}
        mimePrototype.set("toString",tsNamedNativeMethod("toString",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsInternalMIMETypeFor(receiver);params:=tsInternalMIMEParamsString(tsInternalMIMEParamsFor(state.params));value:=state.typeName+"/"+state.subtype;if params!=""{value+=";"+params};return tsStringReference(tsStringUTF8(value))}));mimePrototype.descriptors["toString"]=&tsDescriptor{writable:true,configurable:true};mimePrototype.set("toJSON",mimePrototype.get("toString"));mimePrototype.descriptors["toJSON"]=&tsDescriptor{writable:true,configurable:true}
        tag=tsPropertyKey(tsWellKnownSymbol("toStringTag"));mimePrototype.set(tag,tsStringReference(tsStringUTF8("MIMEType")));mimePrototype.descriptors[tag]=&tsDescriptor{configurable:true}
        tsInternalMIME.mimeTypeConstructor=mimeConstructor
        module:=tsNewObject();module.set("MIMEParams",tsFunctionValue(paramsConstructor));module.set("MIMEType",tsFunctionValue(mimeConstructor));tsInternalMIME.module=module
    })
}
func tsInternalMIMEModule(loop *tsLoop) tsValue {tsInternalMIMEInit(loop);return tsObjectValue(tsInternalMIME.module)}
`
