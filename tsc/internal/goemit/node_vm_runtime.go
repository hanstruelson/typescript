package goemit

// nodeVMRuntime provides the public vm constructors and module lifecycle over
// a lightweight local context model. Bun's JSC binding is replaced by this
// in-process adapter; arbitrary engine compilation remains outside Go's runtime.
const nodeVMRuntime = `
type tsNodeVMModuleState struct {
	kind, identifier, status, source string
	context, namespace, failure, evaluateCallback, initializeImportMeta, dynamicImport tsValue
	exportNames, requests []string
	requestAttributes []tsValue
	exports map[string]tsValue
	exportBindings map[string]string
	dependencies []tsValue
	evaluationResult tsValue
	linkCalled bool
}
type tsNodeVMContextState struct {
	object tsValue
	name, origin, microtaskMode string
	allowStrings, allowWasm, notContextified bool
	importer tsValue
	moduleID uint64
}
var tsNodeVMSingleton struct {
	sync.Once
	module, constants, scriptConstructor, moduleConstructor, sourceConstructor, syntheticConstructor tsValue
	contextKey, scriptKey, moduleKey, noErrorKey, perContextIDKey tsValue
	contexts map[uintptr]bool
	contextOptions map[uintptr]*tsNodeVMContextState
	scripts map[uintptr]tsValue
	modules map[uintptr]*tsNodeVMModuleState
	namespaces map[uintptr]bool
	globalID uint64
}
func tsNodeVMKey(loop *tsLoop, value tsValue, key tsValue) tsValue { return tsGet(loop,value,key) }
func tsNodeVMSet(loop *tsLoop, value,key,item tsValue) { tsSet(loop,value,key,item) }
func tsNodeVMError(code,message string) tsValue { return tsFSException("TypeError",code,message) }
func tsNodeVMInvalid(code,message string) { panic(tsThrown{tsNodeVMError(code,message)}) }
func tsNodeVMOptions(value tsValue) tsValue { if tsIsUndefined(value){return tsObjectValue(tsNewObject())};if value.kind!=tsObjectKind&&value.kind!=tsInstanceKind{tsFSInvalid("options","an object")};return value }
func tsNodeVMIsContext(value tsValue) bool {
	if value.kind!=tsObjectKind&&value.kind!=tsInstanceKind{return false}
	return tsNodeVMSingleton.contexts[uintptr(value.ref)]
}
func tsNodeVMValidateContext(value tsValue) {
	if !tsNodeVMIsContext(value){panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_TYPE","The \"contextifiedObject\" argument must be an vm.Context")})}
}
func tsNodeVMValidateModule(value tsValue) *tsNodeVMModuleState {
	if value.kind!=tsObjectKind{return nil}
	state:=tsNodeVMSingleton.modules[uintptr(value.ref)];if state==nil{panic(tsThrown{tsFSException("TypeError","ERR_INVALID_THIS","Value of \"this\" must be of type Module")})};return state
}
func tsNodeVMTrimSource(source string) string {
	source=strings.TrimSpace(source)
	for strings.HasSuffix(source,";"){source=strings.TrimSpace(strings.TrimSuffix(source,";"))}
	if strings.HasPrefix(source,"return "){source=strings.TrimSpace(strings.TrimPrefix(source,"return "))}
	return source
}
func tsNodeVMLiteral(loop *tsLoop, source string, context tsValue) (tsValue,bool) {
	source=tsNodeVMTrimSource(source)
	if source=="" { return tsU,true }
	switch source {case "undefined":return tsU,true;case "null":return tsNull,true;case "true":return tsBooleanValue(true),true;case "false":return tsBooleanValue(false),true;case "NaN":return tsNumberValue(math.NaN()),true;case "Infinity":return tsNumberValue(math.Inf(1)),true;case "{}":return tsObjectValue(tsNewObject()),true;case "[]":return tsArrayValue(&tsArray{}),true}
	if (strings.HasPrefix(source,"\"")&&strings.HasSuffix(source,"\"")) || (strings.HasPrefix(source,"'")&&strings.HasSuffix(source,"'")) {
		if strings.HasPrefix(source,"\""){if value,err:=strconv.Unquote(source);err==nil{return tsStringReference(tsStringUTF8(value)),true}}
		value:=strings.Trim(source,"'");value=strings.ReplaceAll(value,"\\n","\n");value=strings.ReplaceAll(value,"\\t","\t");value=strings.ReplaceAll(value,"\\'","'");return tsStringReference(tsStringUTF8(value)),true
	}
	if number,err:=strconv.ParseFloat(source,64);err==nil{return tsNumberValue(number),true}
	if strings.HasPrefix(source,"globalThis."){source=strings.TrimPrefix(source,"globalThis.")}
	if strings.HasPrefix(source,"this."){source=strings.TrimPrefix(source,"this.")}
	if source!=""&&strings.IndexFunc(source,func(r rune)bool{return !(r=='_'||r=='$'||r>='a'&&r<='z'||r>='A'&&r<='Z'||r>='0'&&r<='9')})<0 {
		if context.kind==tsObjectKind {value:=tsGet(loop,context,tsStringReference(tsStringUTF8(source)));if !tsIsUndefined(value){return value,true}}
	}
	for _,operator:=range []string{"+","-","*","/"}{if index:=strings.Index(source,operator);index>0{left,okLeft:=tsNodeVMLiteral(loop,source[:index],context);right,okRight:=tsNodeVMLiteral(loop,source[index+1:],context);if okLeft&&okRight{if operator=="+"&&(left.kind==tsStringKind||right.kind==tsStringKind){return tsStringReference(tsStringUTF8(tsText(left)+tsText(right))),true};a,b:=tsNumber(left),tsNumber(right);switch operator{case "+":return tsNumberValue(a+b),true;case "-":return tsNumberValue(a-b),true;case "*":return tsNumberValue(a*b),true;case "/":return tsNumberValue(a/b),true}}}}
	return tsU,false
}
func tsNodeVMEvaluateSource(loop *tsLoop, source string, context tsValue) tsValue {
	text:=strings.TrimSpace(source)
	for _,prefix:=range []string{"const ","let ","var "}{if strings.HasPrefix(text,prefix){text=strings.TrimSpace(strings.TrimPrefix(text,prefix));break}}
	if equal:=strings.Index(text,"=");equal>0&&!strings.Contains(text,"==") {
		name:=strings.TrimSpace(text[:equal]);if strings.HasPrefix(name,"globalThis."){name=strings.TrimPrefix(name,"globalThis.")};if strings.HasPrefix(name,"this."){name=strings.TrimPrefix(name,"this.")}
		if name!=""{if value,ok:=tsNodeVMLiteral(loop,text[equal+1:],context);ok{if context.kind==tsObjectKind{tsSet(loop,context,tsStringReference(tsStringUTF8(name)),value)};return value}}
	}
	if value,ok:=tsNodeVMLiteral(loop,text,context);ok{return value}
	return tsU
}
func tsNodeVMContextOptions(loop *tsLoop,options tsValue)*tsNodeVMContextState{
	state:=&tsNodeVMContextState{allowStrings:true,allowWasm:true,importer:tsU}
	if tsIsUndefined(options){return state}
	opts:=tsNodeVMOptions(options)
	for _,field:=range []string{"name","origin"}{value:=tsGet(loop,opts,tsStringReference(tsStringUTF8(field)));if !tsIsUndefined(value)&&value.kind!=tsStringKind{tsFSInvalid("options."+field,"a string")};if value.kind==tsStringKind{if field=="name"{state.name=tsText(value)}else{state.origin=tsText(value)}}}
	microtask:=tsGet(loop,opts,tsStringReference(tsStringUTF8("microtaskMode")));if !tsIsUndefined(microtask){if microtask.kind!=tsStringKind||tsText(microtask)!="afterEvaluate"{panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_VALUE","The value of \"options.microtaskMode\" is invalid")})};state.microtaskMode="afterEvaluate"}
	codeGeneration:=tsGet(loop,opts,tsStringReference(tsStringUTF8("codeGeneration")));if !tsIsUndefined(codeGeneration){if codeGeneration.kind!=tsObjectKind&&codeGeneration.kind!=tsInstanceKind{tsFSInvalid("options.codeGeneration","an object")};for _,field:=range []string{"strings","wasm"}{value:=tsGet(loop,codeGeneration,tsStringReference(tsStringUTF8(field)));if tsIsUndefined(value){continue};if value.kind!=tsBooleanKind{tsFSInvalid("options.codeGeneration."+field,"a boolean")};if field=="strings"{state.allowStrings=tsTruthy(value)}else{state.allowWasm=tsTruthy(value)}}}
	importer:=tsGet(loop,opts,tsStringReference(tsStringUTF8("importModuleDynamically")));if importer.kind==tsFunctionKind{state.importer=importer}
	return state
}
func tsNodeVMCreateContext(loop *tsLoop, context,options tsValue) tsValue {
	state:=tsNodeVMContextOptions(loop,options)
	if tsIsUndefined(context){context=tsObjectValue(tsNewObject())}else if context.kind==tsSymbolKind{dont:=tsGet(loop,tsNodeVMSingleton.constants,tsStringReference(tsStringUTF8("DONT_CONTEXTIFY")));if !tsStrictEqual(context,dont){tsFSInvalid("contextObject","an object")};state.notContextified=true;context=tsObjectValue(tsNewObject())}else if context.kind!=tsObjectKind&&context.kind!=tsInstanceKind{tsFSInvalid("contextObject","an object")}
	key:=uintptr(context.ref);if tsNodeVMSingleton.contexts[key]{return context}
	state.object=context;tsNodeVMSingleton.contexts[key]=true;tsNodeVMSingleton.contextOptions[key]=state
	globalThis:=tsStringReference(tsStringUTF8("globalThis"));if tsIsUndefined(tsNodeVMKey(loop,context,globalThis)){tsNodeVMSet(loop,context,globalThis,context)}
	return context
}
func tsNodeVMScriptState(value tsValue) tsValue { if value.kind!=tsObjectKind{return tsU};return tsNodeVMSingleton.scripts[uintptr(value.ref)] }
func tsNodeVMSourceMapURL(source string)tsValue{
	for _,marker:=range []string{"//# sourceMappingURL=","//@ sourceMappingURL="}{if index:=strings.LastIndex(source,marker);index>=0{line:=source[index+len(marker):];if newline:=strings.IndexAny(line,"\r\n");newline>=0{line=line[:newline]};line=strings.TrimSpace(line);if line!=""{return tsStringReference(tsStringUTF8(line))}}}
	return tsU
}
func tsNodeVMPrepareScriptOptions(loop *tsLoop,value tsValue,allowFilenameString bool)tsValue{
	if allowFilenameString&&value.kind==tsStringKind{options:=tsNewObject();options.set("filename",value);return tsObjectValue(options)}
	opts:=tsNodeVMOptions(value)
	for _,name:=range []string{"filename","contextName","contextOrigin"}{field:=tsGet(loop,opts,tsStringReference(tsStringUTF8(name)));if !tsIsUndefined(field)&&field.kind!=tsStringKind{tsFSInvalid("options."+name,"a string")}}
	for _,name:=range []string{"lineOffset","columnOffset"}{field:=tsGet(loop,opts,tsStringReference(tsStringUTF8(name)));if !tsIsUndefined(field)&&(!tsIsNumeric(field)||math.Trunc(tsNumber(field))!=tsNumber(field)||tsNumber(field)<math.MinInt32||tsNumber(field)>math.MaxInt32){tsFSInvalid("options."+name,"an integer")}}
	produce:=tsGet(loop,opts,tsStringReference(tsStringUTF8("produceCachedData")));if !tsIsUndefined(produce)&&produce.kind!=tsBooleanKind{tsFSInvalid("options.produceCachedData","a boolean")}
	if cached:=tsGet(loop,opts,tsStringReference(tsStringUTF8("cachedData")));!tsIsUndefined(cached)&&cached.kind!=tsTypedArrayKind&&cached.kind!=tsDataViewKind{tsFSInvalid("options.cachedData","Buffer, TypedArray, or DataView")}
	if timeout:=tsGet(loop,opts,tsStringReference(tsStringUTF8("timeout")));!tsIsUndefined(timeout)&&(!tsIsNumeric(timeout)||math.Trunc(tsNumber(timeout))!=tsNumber(timeout)||tsNumber(timeout)<1){tsFSInvalid("options.timeout","an integer greater than or equal to 1")}
	return opts
}
func tsNodeVMNewScript(loop *tsLoop, source, options tsValue) tsValue {
	if tsIsUndefined(source){source=tsStringReference(tsStringUTF8(""))}else if source.kind!=tsStringKind{source=tsStringReference(tsStringUTF8(tsText(source)))}
	opts:=tsNodeVMPrepareScriptOptions(loop,options,true)
	object:=tsNewObject();object.prototype=tsGet(loop,tsNodeVMSingleton.scriptConstructor,tsStringReference(tsStringUTF8("prototype")));value:=tsObjectValue(object)
	state:=tsNewObject();state.set("source",source);state.set("options",opts);state.set("sourceMapURL",tsNodeVMSourceMapURL(tsText(source)));state.set("cachedDataProduced",tsBooleanValue(false));state.set("cachedDataRejected",tsU);state.set("produceCachedData",tsGet(loop,opts,tsStringReference(tsStringUTF8("produceCachedData"))))
	if cached:=tsGet(loop,opts,tsStringReference(tsStringUTF8("cachedData")));!tsIsUndefined(cached){state.set("cachedDataRejected",tsBooleanValue(true))}
	tsNodeVMSingleton.scripts[uintptr(value.ref)]=tsObjectValue(state)
	return value
}
func tsNodeVMRunOptions(loop *tsLoop,value tsValue)tsValue{
	opts:=tsNodeVMOptions(value)
	filename:=tsGet(loop,opts,tsStringReference(tsStringUTF8("filename")));if !tsIsUndefined(filename)&&filename.kind!=tsStringKind{tsFSInvalid("options.filename","a string")}
	for _,name:=range []string{"lineOffset","columnOffset"}{field:=tsGet(loop,opts,tsStringReference(tsStringUTF8(name)));if !tsIsUndefined(field)&&(!tsIsNumeric(field)||math.Trunc(tsNumber(field))!=tsNumber(field)||tsNumber(field)<math.MinInt32||tsNumber(field)>math.MaxInt32){tsFSInvalid("options."+name,"an integer")}}
	for _,name:=range []string{"displayErrors","breakOnSigint"}{field:=tsGet(loop,opts,tsStringReference(tsStringUTF8(name)));if !tsIsUndefined(field)&&field.kind!=tsBooleanKind{tsFSInvalid("options."+name,"a boolean")}}
	if timeout:=tsGet(loop,opts,tsStringReference(tsStringUTF8("timeout")));!tsIsUndefined(timeout)&&(!tsIsNumeric(timeout)||math.Trunc(tsNumber(timeout))!=tsNumber(timeout)||tsNumber(timeout)<1){tsFSInvalid("options.timeout","an integer greater than or equal to 1")}
	return opts
}
func tsNodeVMRunScript(loop *tsLoop,source string,context,options tsValue)tsValue{
	if context.kind==tsObjectKind||context.kind==tsInstanceKind{if state:=tsNodeVMSingleton.contextOptions[uintptr(context.ref)];state!=nil{if !state.allowStrings&&(strings.Contains(source,"eval(")||strings.Contains(source,"new Function")||strings.Contains(source,"Function(")){panic(tsThrown{tsFSException("EvalError","ERR_EVAL_DISABLED","Code generation from strings disallowed for this context")})};if !state.allowWasm&&(strings.Contains(source,"WebAssembly.compile")||strings.Contains(source,"new WebAssembly.Module")){panic(tsThrown{tsFSException("CompileError","ERR_WEBASSEMBLY_COMPILE_ERROR","Wasm code generation disallowed by embedder")})}}}
	_ = options
	return tsNodeVMEvaluateSource(loop,source,context)
}
func tsNodeVMRunObjectOptions(loop *tsLoop,value tsValue)tsValue{
	if value.kind==tsStringKind{opts:=tsNewObject();opts.set("filename",value);return tsObjectValue(opts)}
	if tsNullish(value)||value.kind!=tsObjectKind&&value.kind!=tsInstanceKind{return tsObjectValue(tsNewObject())}
	return tsObjectValue(tsNodeTLSObjectCopy(value))
}
type tsNodeVMLexeme struct { text string; kind byte; line int }
type tsNodeVMRequest struct { specifier string; attributes map[string]string; attributeOrder []string }
func tsNodeVMIdentifierByte(value byte)bool{return value=='_'||value=='$'||value>='a'&&value<='z'||value>='A'&&value<='Z'||value>='0'&&value<='9'||value>=0x80}
func tsNodeVMDecodeString(raw string,quote byte)string{
	if quote=='"'{if value,err:=strconv.Unquote("\""+raw+"\"");err==nil{return value}}
	return strings.NewReplacer("\\n","\n","\\r","\r","\\t","\t","\\'","'","\\\"","\"","\\\\","\\").Replace(raw)
}
func tsNodeVMLex(source string)[]tsNodeVMLexeme{
	tokens:=[]tsNodeVMLexeme{};line:=1
	for index:=0;index<len(source);{
		current:=source[index]
		if current==' '||current=='\t'||current=='\r'||current=='\n'{if current=='\n'{line++};index++;continue}
		if current=='/'&&index+1<len(source)&&source[index+1]=='/'{index+=2;for index<len(source)&&source[index]!='\n'{index++};continue}
		if current=='/'&&index+1<len(source)&&source[index+1]=='*'{index+=2;for index+1<len(source)&&!(source[index]=='*'&&source[index+1]=='/'){if source[index]=='\n'{line++};index++};if index+1<len(source){index+=2};continue}
		start:=index;tokenLine:=line
		if current=='\''||current=='"'||current=='`'{quote:=current;index++;for index<len(source){if source[index]=='\\'{index+=2;continue};if source[index]==quote{break};if source[index]=='\n'{line++};index++};end:=index;if index<len(source){index++};tokens=append(tokens,tsNodeVMLexeme{text:tsNodeVMDecodeString(source[start+1:end],quote),kind:'s',line:tokenLine});continue}
		if tsNodeVMIdentifierByte(current){index++;for index<len(source)&&tsNodeVMIdentifierByte(source[index]){index++};tokens=append(tokens,tsNodeVMLexeme{text:source[start:index],kind:'i',line:tokenLine});continue}
		index++;tokens=append(tokens,tsNodeVMLexeme{text:source[start:index],kind:'p',line:tokenLine})
	}
	return tokens
}
func tsNodeVMStaticRequests(source string)[]tsNodeVMRequest{
	tokens:=tsNodeVMLex(source);requests:=[]tsNodeVMRequest{};depth:=0
	for index:=0;index<len(tokens);index++{
		token:=tokens[index]
		if depth==0&&(token.text=="import"||token.text=="export"){
			if token.text=="import"&&index+1<len(tokens)&&tokens[index+1].kind=='s'{requests=append(requests,tsNodeVMRequest{specifier:tokens[index+1].text,attributes:map[string]string{}});index++;continue}
			if token.text=="import"&&index+1<len(tokens)&&(tokens[index+1].text=="("||tokens[index+1].text=="."||tokens[index+1].text=="="){continue}
			line:=token.line;from:=-1
			for cursor:=index+1;cursor<len(tokens);cursor++{
				if tokens[cursor].text==";"||(tokens[cursor].line>line&&(tokens[cursor].text=="import"||tokens[cursor].text=="export")){break}
				if tokens[cursor].text=="from"&&cursor+1<len(tokens)&&tokens[cursor+1].kind=='s'{from=cursor;break}
			}
			if from>=0{
				request:=tsNodeVMRequest{specifier:tokens[from+1].text,attributes:map[string]string{}}
				cursor:=from+2
				if cursor+1<len(tokens)&&(tokens[cursor].text=="with"||tokens[cursor].text=="assert")&&tokens[cursor+1].text=="{"{
					for cursor+=2;cursor<len(tokens)&&tokens[cursor].text!="}";{
						key:=tokens[cursor].text;cursor++;if cursor<len(tokens)&&tokens[cursor].text==":"{cursor++;if cursor<len(tokens){if _,exists:=request.attributes[key];!exists{request.attributeOrder=append(request.attributeOrder,key)};request.attributes[key]=tokens[cursor].text;cursor++}};for cursor<len(tokens)&&tokens[cursor].text!=","&&tokens[cursor].text!="}"{cursor++};if cursor<len(tokens)&&tokens[cursor].text==","{cursor++}
					}
				}
				requests=append(requests,request);index=from+1
			}
		}
		switch token.text{case "{","(","[":depth++;case "}",")","]":if depth>0{depth--}}
	}
	return requests
}
func tsNodeVMImports(source string)[]string{requests:=tsNodeVMStaticRequests(source);out:=make([]string,0,len(requests));for _,request:=range requests{out=append(out,request.specifier)};return out}
func tsNodeVMExportNames(source string)[]string{
	tokens:=tsNodeVMLex(source);out:=[]string{};seen:=map[string]bool{};depth:=0;add:=func(name string){if name!=""&&!seen[name]{seen[name]=true;out=append(out,name)}}
	for index:=0;index<len(tokens);index++{token:=tokens[index];if depth==0&&token.text=="export"&&index+1<len(tokens){next:=tokens[index+1].text;switch next{case "default":add("default");case "const","let","var","function","class","async":cursor:=index+2;if next=="async"&&cursor<len(tokens)&&tokens[cursor].text=="function"{cursor++};if cursor<len(tokens)&&tokens[cursor].kind=='i'{add(tokens[cursor].text)};case "{":for cursor:=index+2;cursor<len(tokens)&&tokens[cursor].text!="}";{local:=tokens[cursor].text;exported:=local;cursor++;if cursor<len(tokens)&&tokens[cursor].text=="as"&&cursor+1<len(tokens){exported=tokens[cursor+1].text;cursor+=2};add(exported);for cursor<len(tokens)&&tokens[cursor].text!=","&&tokens[cursor].text!="}"{cursor++};if cursor<len(tokens)&&tokens[cursor].text==","{cursor++}}}}
		switch token.text{case "{","(","[":depth++;case "}",")","]":if depth>0{depth--}}
	}
	return out
}
func tsNodeVMExportBindings(source string)map[string]string{
	tokens:=tsNodeVMLex(source);bindings:=map[string]string{};depth:=0
	for index:=0;index<len(tokens);index++{token:=tokens[index];if depth==0&&token.text=="export"&&index+1<len(tokens){switch tokens[index+1].text{
		case "default":bindings["default"]="default"
		case "const","let","var","function","class":if index+2<len(tokens)&&tokens[index+2].kind=='i'{name:=tokens[index+2].text;bindings[name]=name}
		case "async":if index+3<len(tokens)&&tokens[index+2].text=="function"&&tokens[index+3].kind=='i'{name:=tokens[index+3].text;bindings[name]=name}
		case "{":for cursor:=index+2;cursor<len(tokens)&&tokens[cursor].text!="}";{local:=tokens[cursor].text;exported:=local;cursor++;if cursor<len(tokens)&&tokens[cursor].text=="as"&&cursor+1<len(tokens){exported=tokens[cursor+1].text;cursor+=2};if local!="from"{bindings[exported]=local};for cursor<len(tokens)&&tokens[cursor].text!=","&&tokens[cursor].text!="}"{cursor++};if cursor<len(tokens)&&tokens[cursor].text==","{cursor++}}
		}}
		switch token.text{case "{","(","[":depth++;case "}",")","]":if depth>0{depth--}}
	}
	return bindings
}
func tsNodeVMEvaluateModuleBody(loop *tsLoop,state *tsNodeVMModuleState)tsValue{
	for _,line:=range strings.Split(state.source,"\n"){
		for _,part:=range strings.Split(line,";"){
			statement:=strings.TrimSpace(part);if statement==""||strings.HasPrefix(statement,"import ")||strings.HasPrefix(statement,"export {")||strings.HasPrefix(statement,"export *"){continue}
			if strings.HasPrefix(statement,"export default "){value,ok:=tsNodeVMLiteral(loop,strings.TrimSpace(strings.TrimPrefix(statement,"export default ")),state.context);if ok{tsSet(loop,state.context,tsStringReference(tsStringUTF8("default")),value)};continue}
			statement=strings.TrimSpace(strings.TrimPrefix(statement,"export "))
			tsNodeVMEvaluateSource(loop,statement,state.context)
		}
	}
	tsNodeVMEnsureNamespace(state)
	for _,name:=range state.exportNames{local:=state.exportBindings[name];if local==""{local=name};value:=tsGet(loop,state.context,tsStringReference(tsStringUTF8(local)));tsSet(loop,state.namespace,tsStringReference(tsStringUTF8(name)),value)}
	return tsU
}
func tsNodeVMCreateModule(loop *tsLoop,kind string,source tsValue,exportNames []string,callback,options tsValue) tsValue {
	opts:=tsNodeVMOptions(options);context:=tsGet(loop,opts,tsStringReference(tsStringUTF8("context")));contextProvided:=!tsIsUndefined(context)
	if contextProvided{tsInternalValidatorObject(context,tsStringReference(tsStringUTF8("options.context")),tsU);if !tsNodeVMIsContext(context){panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_TYPE","The \"options.context\" argument must be of type vm.Context")})}}
	identifier:=tsGet(loop,opts,tsStringReference(tsStringUTF8("identifier")));if !tsIsUndefined(identifier)&&identifier.kind!=tsStringKind{tsFSInvalid("options.identifier","a string")}
	if kind=="source"&&source.kind!=tsStringKind{tsFSInvalid("sourceText","a string")}
	if kind=="source"{for _,name:=range []string{"lineOffset","columnOffset"}{value:=tsGet(loop,opts,tsStringReference(tsStringUTF8(name)));if !tsIsUndefined(value)&&(!tsIsNumeric(value)||math.Trunc(tsNumber(value))!=tsNumber(value)||tsNumber(value)<math.MinInt32||tsNumber(value)>math.MaxInt32){tsFSInvalid("options."+name,"an integer")}};if cached:=tsGet(loop,opts,tsStringReference(tsStringUTF8("cachedData")));!tsIsUndefined(cached)&&cached.kind!=tsTypedArrayKind&&cached.kind!=tsDataViewKind{tsFSInvalid("options.cachedData","Buffer, TypedArray, or DataView")}}
	if tsIsUndefined(context){context=tsGlobalObject()}
	if tsIsUndefined(identifier){var id uint64;if contextProvided{if contextState:=tsNodeVMSingleton.contextOptions[uintptr(context.ref)];contextState!=nil{id=contextState.moduleID;contextState.moduleID++}}else{id=tsNodeVMSingleton.globalID;tsNodeVMSingleton.globalID++};identifier=tsStringReference(tsStringUTF8("vm:module("+strconv.FormatUint(id,10)+")"))}
	initialize:=tsGet(loop,opts,tsStringReference(tsStringUTF8("initializeImportMeta")));dynamicImport:=tsGet(loop,opts,tsStringReference(tsStringUTF8("importModuleDynamically")))
	for _,entry:=range []struct{name string;value tsValue}{{"initializeImportMeta",initialize},{"importModuleDynamically",dynamicImport}}{if !tsIsUndefined(entry.value)&&entry.value.kind!=tsFunctionKind{tsFSInvalid("options."+entry.name,"a function")}}
	if kind=="synthetic"&&callback.kind!=tsFunctionKind{tsFSInvalid("evaluateCallback","a function")}
	prototype:=tsGet(loop,tsNodeVMSingleton.moduleConstructor,tsStringReference(tsStringUTF8("prototype")));if kind=="source"{prototype=tsGet(loop,tsNodeVMSingleton.sourceConstructor,tsStringReference(tsStringUTF8("prototype")))}else if kind=="synthetic"{prototype=tsGet(loop,tsNodeVMSingleton.syntheticConstructor,tsStringReference(tsStringUTF8("prototype")))}
	object:=tsNewObject();object.prototype=prototype;value:=tsObjectValue(object)
	state:=&tsNodeVMModuleState{kind:kind,identifier:tsText(identifier),context:context,source:tsText(source),status:"unlinked",failure:tsNodeVMSingleton.noErrorKey,exports:map[string]tsValue{},evaluateCallback:callback,initializeImportMeta:initialize,dynamicImport:dynamicImport}
	if kind=="synthetic"{state.status="linked";state.exportNames=append([]string{},exportNames...)}else{requests:=tsNodeVMStaticRequests(state.source);state.requests=make([]string,0,len(requests));state.requestAttributes=make([]tsValue,0,len(requests));for _,request:=range requests{state.requests=append(state.requests,request.specifier);attributes:=tsNewObject();for _,key:=range request.attributeOrder{attributes.set(key,tsStringReference(tsStringUTF8(request.attributes[key])))};state.requestAttributes=append(state.requestAttributes,tsObjectValue(attributes))};state.exportNames=tsNodeVMExportNames(state.source);state.exportBindings=tsNodeVMExportBindings(state.source)}
	tsNodeVMSingleton.modules[uintptr(value.ref)]=state
	if kind=="synthetic"{tsNodeVMEnsureNamespace(state)}
	return value
}
func tsNodeVMEnsureNamespace(state *tsNodeVMModuleState){
	if state.namespace.kind==tsObjectKind{return};namespace:=tsNewObject();namespace.prototype=tsNull;for _,name:=range state.exportNames{namespace.set(name,tsU)};state.namespace=tsObjectValue(namespace);tsNodeVMSingleton.namespaces[uintptr(namespace)]=true
}
func tsNodeVMStatusCode(status string) int {
	switch status {
	case "unlinked":
		return 0
	case "linking":
		return 1
	case "linked":
		return 2
	case "evaluating":
		return 3
	case "evaluated":
		return 4
	case "errored":
		return 5
	default:
		return 0
	}
}
func tsNodeVMModuleRequestArray(loop *tsLoop, state *tsNodeVMModuleState) tsValue {
	items := &tsArray{}
	for index, specifier := range state.requests {
		request := tsNewObject()
		request.prototype = tsNull
		request.set("specifier", tsStringReference(tsStringUTF8(specifier)))
		attributes := tsNewObject()
		if index < len(state.requestAttributes) {
			source := state.requestAttributes[index]
			if source.kind == tsObjectKind {
				object := (*tsObject)(source.ref)
				for _, key := range object.order {
					attributes.set(key, object.values[key])
				}
			}
		}
		tsDescriptorBuiltin(loop, "freeze", []tsValue{tsObjectValue(attributes)})
		request.set("attributes", tsObjectValue(attributes))
		request.set("phase", tsStringReference(tsStringUTF8("evaluation")))
		requestValue := tsObjectValue(request)
		tsDescriptorBuiltin(loop, "freeze", []tsValue{requestValue})
		items.values = append(items.values, requestValue)
	}
	value := tsArrayValue(items)
	tsDescriptorBuiltin(loop, "freeze", []tsValue{value})
	return value
}
func tsNodeVMModuleFailure(code, message string) tsValue {
	return tsFSException("Error", code, message)
}
func tsNodeVMCallLinker(loop *tsLoop, linker, specifier, parent, options tsValue) (value, failure tsValue, failed bool) {
	defer func() {
		if recovered := recover(); recovered != nil {
			failure = tsUnwrap(recovered)
			failed = true
		}
	}()
	value = tsCall(loop, linker, specifier, parent, options)
	return
}
func tsNodeVMLinkGraph(loop *tsLoop, module, linker tsValue, complete func(tsResult)) {
	state := tsNodeVMValidateModule(module)
	if state.status == "linking" || state.status == "linked" || state.status == "evaluated" {
		complete(tsResult{tsU, false})
		return
	}
	if state.status != "unlinked" {
		complete(tsResult{tsNodeVMModuleFailure("ERR_VM_MODULE_STATUS", "Module must be unlinked"), true})
		return
	}
	if linker.kind != tsFunctionKind {
	complete(tsResult{tsNodeVMError("ERR_INVALID_ARG_TYPE", "The \"linker\" argument must be of type function"), true})
		return
	}

	state.status = "linking"
	resolved := make([]tsValue, len(state.requests))
	index := 0
	var next func()
	next = func() {
		if index == len(state.requests) {
			state.dependencies = resolved
			state.linkCalled = true
			state.status = "unlinked"
			complete(tsResult{tsU, false})
			return
		}
		current := index
		specifier := tsStringReference(tsStringUTF8(state.requests[current]))
		attributes := tsObjectValue(tsNewObject())
		if current < len(state.requestAttributes) {
			attributes = state.requestAttributes[current]
		}
		linkOptions := tsNewObject()
		linkOptions.set("attributes", attributes)
		linkOptions.set("assert", attributes)
		value, failure, failed := tsNodeVMCallLinker(loop, linker, specifier, module, tsObjectValue(linkOptions))
		if failed {
			state.failure = failure
			state.status = "errored"
			complete(tsResult{failure, true})
			return
		}
		loop.await(value, func(result tsResult) {
			if result.rejected {
				state.failure = result.value
				state.status = "errored"
				complete(result)
				return
			}
			if !tsNodeVMIsModule(result.value) {
				failure := tsFSException("TypeError", "ERR_VM_MODULE_NOT_MODULE", "Linker must return a Module")
				state.failure = failure
				state.status = "errored"
				complete(tsResult{failure, true})
				return
			}
			dependency := tsNodeVMValidateModule(result.value)
			if !tsStrictEqual(dependency.context, state.context) {
				failure := tsNodeVMModuleFailure("ERR_VM_MODULE_DIFFERENT_CONTEXT", "Linked modules must use the same context")
				state.failure = failure
				state.status = "errored"
				complete(tsResult{failure, true})
				return
			}
			if dependency.status == "errored" {
				failure := tsNodeVMModuleFailure("ERR_VM_MODULE_LINK_FAILURE", "Linker resolved to an errored module")
				state.failure = failure
				state.status = "errored"
				complete(tsResult{failure, true})
				return
			}
			if dependency.status == "unlinked" {
				tsNodeVMLinkGraph(loop, result.value, linker, func(linked tsResult) {
					if linked.rejected {
						state.failure = linked.value
						state.status = "errored"
						complete(linked)
						return
					}
					resolved[current] = result.value
					index++
					next()
				})
				return
			}
			resolved[current] = result.value
			index++
			next()
		})
	}
	next()
}
func tsNodeVMInstantiate(state *tsNodeVMModuleState) tsValue {
	if state.kind == "synthetic" {
		if tsNodeVMStatusCode(state.status) >= tsNodeVMStatusCode("linked") {
			return tsNodeVMModuleFailure("ERR_VM_MODULE_STATUS", "Cannot reinstantiate a SyntheticModule")
		}
		state.status = "linked"
		tsNodeVMEnsureNamespace(state)
		return tsU
	}
	if state.status != "unlinked" {
		return tsNodeVMModuleFailure("ERR_VM_MODULE_STATUS", "Module must be unlinked")
	}
	if !state.linkCalled {
		return tsNodeVMModuleFailure("ERR_VM_MODULE_LINK_FAILURE", "module is not linked")
	}
	worklist := []tsValue{}
	seen := map[uintptr]bool{}
	for _, dependency := range state.dependencies {
		worklist = append(worklist, dependency)
	}
	state.status = "linked"
	tsNodeVMEnsureNamespace(state)
	for len(worklist) > 0 {
		last := len(worklist) - 1
		value := worklist[last]
		worklist = worklist[:last]
		if !tsNodeVMIsModule(value) {
			continue
		}
		key := uintptr(value.ref)
		if seen[key] {
			continue
		}
		seen[key] = true
		dependency := tsNodeVMValidateModule(value)
		if dependency.status == "unlinked" {
			dependency.status = "linked"
			tsNodeVMEnsureNamespace(dependency)
		}
		worklist = append(worklist, dependency.dependencies...)
	}
	return tsU
}
func tsNodeVMEvaluateGraph(loop *tsLoop,module tsValue,visited map[uintptr]bool)tsValue{
	state:=tsNodeVMValidateModule(module)
	if state.status=="evaluated"{return tsU}
	if state.status=="errored"{return state.failure}
	if state.status!="linked"&&state.status!="evaluating"{return tsNodeVMModuleFailure("ERR_VM_MODULE_STATUS","Module must be linked before evaluation")}
	key:=uintptr(module.ref);if visited[key]{return tsU};visited[key]=true
	state.status="evaluating"
	for _,dependency:=range state.dependencies{if tsNodeVMIsModule(dependency){if failure:=tsNodeVMEvaluateGraph(loop,dependency,visited);!tsIsUndefined(failure){state.failure=failure;state.status="errored";return failure}}}
	runCallback:=func(callback func())(failure tsValue,failed bool){defer func(){if recovered:=recover();recovered!=nil{failure=tsUnwrap(recovered);failed=true}}();callback();return tsU,false}
	if state.kind=="synthetic"&&state.evaluateCallback.kind==tsFunctionKind{
		if failure,failed:=runCallback(func(){tsCallReceiver(loop,state.evaluateCallback,module)});failed{state.failure=failure;state.status="errored";return failure}
	}else if state.kind=="source"{
		if state.initializeImportMeta.kind==tsFunctionKind{meta:=tsObjectValue(tsNewObject());if failure,failed:=runCallback(func(){tsCall(loop,state.initializeImportMeta,meta,module)});failed{state.failure=failure;state.status="errored";return failure}}
		if failure,failed:=runCallback(func(){tsNodeVMEvaluateModuleBody(loop,state)});failed{state.failure=failure;state.status="errored";return failure}
	}
	state.status="evaluated";state.evaluationResult=tsU
	return tsU
}
func tsNodeVMInstallModuleMethods(loop *tsLoop, prototype *tsObject) {
	prototype.descriptors = map[string]*tsDescriptor{}
	for _, name := range []string{"identifier", "context", "status", "namespace", "error", "moduleRequests", "dependencySpecifiers"} {
		property := name
		prototype.descriptors[property] = &tsDescriptor{accessor: true, configurable: true, getter: tsNamedNativeMethod("get "+property, 0, func(loop *tsLoop, receiver tsValue, _ ...tsValue) tsValue {
			state := tsNodeVMValidateModule(receiver)
			switch property {
			case "identifier":
				return tsStringReference(tsStringUTF8(state.identifier))
			case "context":
				return state.context
			case "status":
				return tsStringReference(tsStringUTF8(state.status))
			case "namespace":
				if state.status == "unlinked" || state.status == "linking" {
					tsNodeVMInvalid("ERR_VM_MODULE_STATUS", "Module namespace is not available before linking")
				}
				tsNodeVMEnsureNamespace(state)
				return state.namespace
			case "error":
				if state.status != "errored" {
					tsNodeVMInvalid("ERR_VM_MODULE_STATUS", "Module is not errored")
				}
				return state.failure
			case "moduleRequests":
				if state.kind != "source" {
					return tsU
				}
				return tsNodeVMModuleRequestArray(loop, state)
			case "dependencySpecifiers":
				items := &tsArray{}
				for _, specifier := range state.requests {
					items.values = append(items.values, tsStringReference(tsStringUTF8(specifier)))
				}
				value := tsArrayValue(items)
				tsDescriptorBuiltin(loop, "freeze", []tsValue{value})
				return value
			}
			return tsU
		})}
	}
	prototype.set("link", tsNamedNativeMethod("link", 1, func(loop *tsLoop, receiver tsValue, args ...tsValue) tsValue {
		state := tsNodeVMValidateModule(receiver)
		if state.kind == "synthetic" {
			return tsU
		}
		linker := tsArg(args, 0)
		if linker.kind != tsFunctionKind {
			tsFSInvalid("linker", "a function")
		}
		if state.status == "linked" {
			return tsPromiseValue(loop.resolved(tsNodeVMModuleFailure("ERR_VM_MODULE_ALREADY_LINKED", "Module has already been linked"), true))
		}
		if state.status != "unlinked" {
			return tsPromiseValue(loop.resolved(tsNodeVMModuleFailure("ERR_VM_MODULE_STATUS", "Module must be unlinked"), true))
		}
		promise := loop.promise()
		tsNodeVMLinkGraph(loop, receiver, linker, func(result tsResult) {
			if result.rejected {
				promise.settle(result)
				return
			}
			if failure := tsNodeVMInstantiate(state); !tsIsUndefined(failure) {
				state.failure = failure
				state.status = "errored"
				promise.settle(tsResult{failure, true})
				return
			}
			promise.settle(tsResult{tsU, false})
		})
		return tsPromiseValue(promise)
	}))
	prototype.set("instantiate", tsNamedNativeMethod("instantiate", 0, func(_ *tsLoop, receiver tsValue, _ ...tsValue) tsValue {
		state := tsNodeVMValidateModule(receiver)
		if failure := tsNodeVMInstantiate(state); !tsIsUndefined(failure) {
			tsNodeVMInvalid("ERR_VM_MODULE_STATUS", tsText(failure))
		}
		return tsU
	}))
	prototype.set("evaluate", tsNamedNativeMethod("evaluate", 2, func(loop *tsLoop, receiver tsValue, args ...tsValue) tsValue {
		state := tsNodeVMValidateModule(receiver)
		options := tsArg(args, 0)
		if tsIsUndefined(options) {
			options = tsObjectValue(tsNewObject())
		}
		if options.kind != tsObjectKind && options.kind != tsInstanceKind {
			return tsPromiseValue(loop.resolved(tsFSException("TypeError", "ERR_INVALID_ARG_TYPE", "options must be an object"), true))
		}
		timeout := tsGet(loop, options, tsStringReference(tsStringUTF8("timeout")))
		if !tsIsUndefined(timeout) && (!tsIsNumeric(timeout) || math.Trunc(tsNumber(timeout)) != tsNumber(timeout) || tsNumber(timeout) < 0 || tsNumber(timeout) > math.MaxUint32) {
			return tsPromiseValue(loop.resolved(tsFSException("TypeError", "ERR_INVALID_ARG_VALUE", "options.timeout must be a uint32"), true))
		}
		breakOnSigint := tsGet(loop, options, tsStringReference(tsStringUTF8("breakOnSigint")))
		if tsIsUndefined(breakOnSigint) {
			breakOnSigint = tsBooleanValue(false)
		}
		if breakOnSigint.kind != tsBooleanKind {
			return tsPromiseValue(loop.resolved(tsFSException("TypeError", "ERR_INVALID_ARG_TYPE", "options.breakOnSigint must be a boolean"), true))
		}
		if state.status != "linked" && state.status != "evaluated" && state.status != "errored" {
			return tsPromiseValue(loop.resolved(tsNodeVMModuleFailure("ERR_VM_MODULE_STATUS", "Module must be linked, evaluated or errored before evaluating"), true))
		}
		if state.status == "errored" {
			return tsPromiseValue(loop.resolved(state.failure, true))
		}
		if state.status == "evaluated" {
			return tsPromiseValue(loop.resolved(tsU, false))
		}
		if failure := tsNodeVMEvaluateGraph(loop, receiver, map[uintptr]bool{}); !tsIsUndefined(failure) {
			return tsPromiseValue(loop.resolved(failure, true))
		}
		return tsPromiseValue(loop.resolved(tsU, false))
	}))
	prototype.set("linkRequests", tsNamedNativeMethod("linkRequests", 1, func(_ *tsLoop, receiver tsValue, args ...tsValue) tsValue {
		state := tsNodeVMValidateModule(receiver)
		if state.kind != "source" {
			tsNodeVMInvalid("ERR_INVALID_THIS", "linkRequests is only available on SourceTextModule")
		}
		if state.status != "unlinked" {
			tsNodeVMInvalid("ERR_VM_MODULE_STATUS", "Module must be unlinked before linking")
		}
		modules := tsArg(args, 0)
		if modules.kind != tsArrayKind {
			tsFSInvalid("modules", "an array")
		}
		items := (*tsArray)(modules.ref)
		if len(items.values) != len(state.requests) {
			tsNodeVMInvalid("ERR_MODULE_LINK_MISMATCH", "Module count does not match requests")
		}
		resolved := make([]tsValue, len(items.values))
		for index, item := range items.values {
			if !tsNodeVMIsModule(item) {
				tsNodeVMInvalid("ERR_VM_MODULE_NOT_MODULE", "Expected a Module")
			}
			dependency := tsNodeVMValidateModule(item)
			if !tsStrictEqual(dependency.context, state.context) {
				tsNodeVMInvalid("ERR_VM_MODULE_DIFFERENT_CONTEXT", "Linked modules must use the same context")
			}
			resolved[index] = item
		}
		state.dependencies = resolved
		state.linkCalled = true
		return tsU
	}))
	prototype.set("hasAsyncGraph", tsNamedNativeMethod("hasAsyncGraph", 0, func(_ *tsLoop, receiver tsValue, _ ...tsValue) tsValue {
		state := tsNodeVMValidateModule(receiver)
		if tsNodeVMStatusCode(state.status) < tsNodeVMStatusCode("linked") {
			tsNodeVMInvalid("ERR_VM_MODULE_STATUS", "Module must be instantiated")
		}
		worklist := []*tsNodeVMModuleState{state}
		visited := map[*tsNodeVMModuleState]bool{}
		for len(worklist) > 0 {
			last := len(worklist) - 1
			current := worklist[last]
			worklist = worklist[:last]
			if visited[current] {
				continue
			}
			visited[current] = true
			if current.kind == "source" && strings.Contains(current.source, "await ") {
				return tsBooleanValue(true)
			}
			for _, dependencyValue := range current.dependencies {
				if tsNodeVMIsModule(dependencyValue) {
					worklist = append(worklist, tsNodeVMValidateModule(dependencyValue))
				}
			}
		}
		return tsBooleanValue(false)
	}))
	prototype.set("hasTopLevelAwait", tsNamedNativeMethod("hasTopLevelAwait", 0, func(_ *tsLoop, receiver tsValue, _ ...tsValue) tsValue {
		state := tsNodeVMValidateModule(receiver)
		if state.kind != "source" {
			tsNodeVMInvalid("ERR_INVALID_THIS", "hasTopLevelAwait is only available on SourceTextModule")
		}
		return tsBooleanValue(strings.Contains(state.source, "await "))
	}))
	prototype.set("createCachedData", tsNamedNativeMethod("createCachedData", 0, func(_ *tsLoop, receiver tsValue, _ ...tsValue) tsValue {
		state := tsNodeVMValidateModule(receiver)
		if state.kind != "source" {
			tsNodeVMInvalid("ERR_INVALID_THIS", "createCachedData is only available on SourceTextModule")
		}
		if state.status == "evaluating" || state.status == "evaluated" || state.status == "errored" {
			tsNodeVMInvalid("ERR_VM_MODULE_CANNOT_CREATE_CACHED_DATA", "Cannot create cached data after evaluation")
		}
		return tsNodeBuffer(nil)
	}))
	prototype.set("setExport", tsNamedNativeMethod("setExport", 2, func(loop *tsLoop, receiver tsValue, args ...tsValue) tsValue {
		state := tsNodeVMValidateModule(receiver)
		name := tsArg(args, 0)
		if name.kind != tsStringKind {
			tsFSInvalid("name", "a string")
		}
		if state.kind != "synthetic" {
			tsNodeVMInvalid("ERR_INVALID_ARG_TYPE", "setExport is only available on SyntheticModule")
		}
		if tsNodeVMStatusCode(state.status) < tsNodeVMStatusCode("linked") {
			tsNodeVMInvalid("ERR_VM_MODULE_STATUS", "SyntheticModule must be linked before exports can be set")
		}
		exportName := tsText(name)
		found := false
		for _, item := range state.exportNames {
			if item == exportName {
				found = true
				break
			}
		}
		if !found {
			panic(tsThrown{tsFSException("ReferenceError", "", "Export '"+exportName+"' is not defined in module")})
		}
		state.exports[exportName] = tsArg(args, 1)
		tsNodeVMEnsureNamespace(state)
		tsSet(loop, state.namespace, name, tsArg(args, 1))
		return tsU
	}))
}
func tsNodeVMIsModule(value tsValue) bool { if value.kind!=tsObjectKind{return false};return tsNodeVMSingleton.modules[uintptr(value.ref)]!=nil }
func tsNodeVMModule(loop *tsLoop) tsValue {
	tsNodeVMSingleton.Do(func(){
		tsNodeVMSingleton.contexts=map[uintptr]bool{};tsNodeVMSingleton.contextOptions=map[uintptr]*tsNodeVMContextState{};tsNodeVMSingleton.scripts=map[uintptr]tsValue{};tsNodeVMSingleton.modules=map[uintptr]*tsNodeVMModuleState{};tsNodeVMSingleton.namespaces=map[uintptr]bool{}
		tsNodeVMSingleton.contextKey=tsNewSymbol(tsStringUTF8("vm-context"));tsNodeVMSingleton.scriptKey=tsNewSymbol(tsStringUTF8("vm-script"));tsNodeVMSingleton.moduleKey=tsNewSymbol(tsStringUTF8("vm-module"));tsNodeVMSingleton.noErrorKey=tsNewSymbol(tsStringUTF8("vm-no-error"));tsNodeVMSingleton.perContextIDKey=tsNewSymbol(tsStringUTF8("vm-module-id"))
		scriptPrototype:=tsNewObject();scriptPrototype.prototype=tsObjectPrototype()
		script:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeVMNewScript(loop,tsArg(args,0),tsArg(args,1))});script.name="Script";script.length=2;script.constructible=true;script.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeVMNewScript(loop,tsArg(args,0),tsArg(args,1))};script.properties=tsNewObject();script.properties.set("prototype",tsObjectValue(scriptPrototype))
		for _,entry:=range []struct{name string;length int}{{"runInThisContext",1},{"runInContext",2},{"runInNewContext",2},{"createCachedData",0}}{name:=entry.name;scriptPrototype.set(name,tsNamedNativeMethod(name,entry.length,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{
			state:=tsNodeVMScriptState(receiver);if state.kind!=tsObjectKind{tsNodeNodeVMInvalidReceiver(name)};source:=tsText(tsGet(loop,state,tsStringReference(tsStringUTF8("source"))));runOptions:=tsU;if name!="createCachedData"{runOptions=tsNodeVMRunOptions(loop,tsArg(args,map[string]int{"runInThisContext":0,"runInContext":1,"runInNewContext":1}[name]))}
			switch name{case "runInThisContext":return tsNodeVMRunScript(loop,source,tsGlobalObject(),runOptions);case "runInContext":context:=tsArg(args,0);tsNodeVMValidateContext(context);return tsNodeVMRunScript(loop,source,context,runOptions);case "runInNewContext":context:=tsNodeVMCreateContext(loop,tsArg(args,0),tsArg(args,1));return tsNodeVMRunScript(loop,source,context,runOptions);default:return tsNodeBuffer(nil)}
		}))}
		if scriptPrototype.descriptors==nil{scriptPrototype.descriptors=map[string]*tsDescriptor{}}
		for _,name:=range []string{"sourceMapURL","cachedData","cachedDataProduced","cachedDataRejected"}{property:=name;scriptPrototype.descriptors[property]=&tsDescriptor{accessor:true,configurable:true,getter:tsNamedNativeMethod("get "+property,0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeVMScriptState(receiver);if state.kind!=tsObjectKind{tsNodeNodeVMInvalidReceiver(property)};switch property{case "sourceMapURL":return tsGet(loop,state,tsStringReference(tsStringUTF8(property)));case "cachedData":if tsTruthy(tsGet(loop,state,tsStringReference(tsStringUTF8("produceCachedData")))){return tsNodeBuffer(nil)};return tsU;case "cachedDataProduced","cachedDataRejected":return tsGet(loop,state,tsStringReference(tsStringUTF8(property)))};return tsU})}}
		scriptPrototype.set("constructor",tsFunctionValue(script));tsNodeVMSingleton.scriptConstructor=tsFunctionValue(script)
		modulePrototype:=tsNewObject();modulePrototype.prototype=tsObjectPrototype();tsNodeVMInstallModuleMethods(loop,modulePrototype)
		base:=tsFunc(func(_ *tsLoop,_ ...tsValue)tsValue{tsNodeVMInvalid("ERR_INVALID_THIS","Module is not a constructor");return tsU});base.name="Module";base.constructible=true;base.construct=func(_ *tsLoop,_ ...tsValue)tsValue{tsNodeVMInvalid("ERR_INVALID_THIS","Module is not a constructor");return tsU};base.properties=tsNewObject();base.properties.set("prototype",tsObjectValue(modulePrototype));tsNodeVMSingleton.moduleConstructor=tsFunctionValue(base)
		sourcePrototype:=tsNewObject();sourcePrototype.prototype=tsObjectValue(modulePrototype);sourceCtor:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeVMCreateModule(loop,"source",tsArg(args,0),nil,tsU,tsArg(args,1))});sourceCtor.name="SourceTextModule";sourceCtor.length=2;sourceCtor.constructible=true;sourceCtor.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeVMCreateModule(loop,"source",tsArg(args,0),nil,tsU,tsArg(args,1))};sourceCtor.properties=tsNewObject();sourceCtor.properties.set("prototype",tsObjectValue(sourcePrototype));sourcePrototype.set("constructor",tsFunctionValue(sourceCtor));tsNodeVMSingleton.sourceConstructor=tsFunctionValue(sourceCtor)
		syntheticPrototype:=tsNewObject();syntheticPrototype.prototype=tsObjectValue(modulePrototype);syntheticCtor:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeVMSyntheticCreate(loop,args,syntheticPrototype)});syntheticCtor.name="SyntheticModule";syntheticCtor.length=2;syntheticCtor.constructible=true;syntheticCtor.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeVMSyntheticCreate(loop,args,syntheticPrototype)};syntheticCtor.properties=tsNewObject();syntheticCtor.properties.set("prototype",tsObjectValue(syntheticPrototype));syntheticPrototype.set("constructor",tsFunctionValue(syntheticCtor));tsNodeVMSingleton.syntheticConstructor=tsFunctionValue(syntheticCtor)
		constants:=tsNewObject();dont:=tsNewSymbol(tsStringUTF8("DONT_CONTEXTIFY"));loader:=tsNewSymbol(tsStringUTF8("USE_MAIN_CONTEXT_DEFAULT_LOADER"));constants.set("DONT_CONTEXTIFY",dont);constants.set("USE_MAIN_CONTEXT_DEFAULT_LOADER",loader);tsNodeVMSingleton.constants=tsObjectValue(constants)
		module:=tsNewObject();module.set("createContext",tsNamedNativeMethod("createContext",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeVMCreateContext(loop,tsArg(args,0),tsArg(args,1))}));module.set("isContext",tsNamedNativeMethod("isContext",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsBooleanValue(tsNodeVMIsContext(tsArg(args,0)))}))
		module.set("runInContext",tsNamedNativeMethod("runInContext",3,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{context:=tsArg(args,1);tsNodeVMValidateContext(context);options:=tsNodeVMRunObjectOptions(loop,tsArg(args,2));script:=tsNodeVMNewScript(loop,tsArg(args,0),options);state:=tsNodeVMScriptState(script);return tsNodeVMRunScript(loop,tsText(tsGet(loop,state,tsStringReference(tsStringUTF8("source")))),context,tsNodeVMRunOptions(loop,options))}))
		module.set("runInThisContext",tsNamedNativeMethod("runInThisContext",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{options:=tsArg(args,1);if options.kind==tsStringKind{options=tsNodeVMRunObjectOptions(loop,options)};script:=tsNodeVMNewScript(loop,tsArg(args,0),options);state:=tsNodeVMScriptState(script);return tsNodeVMRunScript(loop,tsText(tsGet(loop,state,tsStringReference(tsStringUTF8("source")))),tsGlobalObject(),tsNodeVMRunOptions(loop,options))}))
		module.set("runInNewContext",tsNamedNativeMethod("runInNewContext",3,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{options:=tsNodeVMRunObjectOptions(loop,tsArg(args,2));context:=tsNodeVMCreateContext(loop,tsArg(args,1),options);script:=tsNodeVMNewScript(loop,tsArg(args,0),options);state:=tsNodeVMScriptState(script);return tsNodeVMRunScript(loop,tsText(tsGet(loop,state,tsStringReference(tsStringUTF8("source")))),context,tsNodeVMRunOptions(loop,options))}))
		module.set("createScript",tsNamedNativeMethod("createScript",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeVMNewScript(loop,tsArg(args,0),tsArg(args,1))}))
		module.set("Script",tsNodeVMSingleton.scriptConstructor);module.set("Module",tsNodeVMSingleton.moduleConstructor);module.set("SourceTextModule",tsNodeVMSingleton.sourceConstructor);module.set("SyntheticModule",tsNodeVMSingleton.syntheticConstructor);module.set("constants",tsObjectValue(constants))
		module.set("compileFunction",tsNamedNativeMethod("compileFunction",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{
			body:=tsArg(args,0);if body.kind!=tsStringKind{tsFSInvalid("code","a string")}
			params:=tsArg(args,1);if tsIsUndefined(params){params=tsArrayValue(&tsArray{})};if params.kind!=tsArrayKind{tsFSInvalid("params","an array")}
			names:=[]string{};for _,value:=range (*tsArray)(params.ref).values{if value.kind!=tsStringKind{tsFSInvalid("params","an array of strings")};names=append(names,tsText(value))}
			opts:=tsNodeVMPrepareScriptOptions(loop,tsArg(args,2),false);context:=tsGet(loop,opts,tsStringReference(tsStringUTF8("parsingContext")));if tsIsUndefined(context){context=tsGlobalObject()}else if !tsNodeVMIsContext(context){tsFSInvalid("options.parsingContext","a vm.Context")}
			base:=tsObjectValue(tsNodeTLSObjectCopy(context));extensions:=tsGet(loop,opts,tsStringReference(tsStringUTF8("contextExtensions")));if !tsIsUndefined(extensions){if extensions.kind!=tsArrayKind{tsFSInvalid("options.contextExtensions","an array")};for index,extension:=range (*tsArray)(extensions.ref).values{if extension.kind!=tsObjectKind&&extension.kind!=tsInstanceKind{tsFSInvalid("options.contextExtensions["+strconv.Itoa(index)+"]","an object")};copy:=tsNodeTLSObjectCopy(extension);for _,key:=range copy.order{base.ref.(*tsObject).set(key,copy.values[key])}}}
			if importer:=tsGet(loop,opts,tsStringReference(tsStringUTF8("importModuleDynamically")));!tsIsUndefined(importer)&&importer.kind!=tsFunctionKind&&!tsStrictEqual(importer,tsGet(loop,tsNodeVMSingleton.constants,tsStringReference(tsStringUTF8("USE_MAIN_CONTEXT_DEFAULT_LOADER")))){tsFSInvalid("options.importModuleDynamically","a function")}
			source:=tsText(body);invoke:=func(loop *tsLoop,values ...tsValue)tsValue{environment:=tsObjectValue(tsNodeTLSObjectCopy(base));for index,name:=range names{tsSet(loop,environment,tsStringReference(tsStringUTF8(name)),tsArg(values,index))};if state:=tsNodeVMSingleton.contextOptions[uintptr(context.ref)];state!=nil&&(!state.allowStrings&&(strings.Contains(source,"eval(")||strings.Contains(source,"new Function")||strings.Contains(source,"Function("))||!state.allowWasm&&(strings.Contains(source,"WebAssembly.compile")||strings.Contains(source,"new WebAssembly.Module"))){panic(tsThrown{tsFSException("EvalError","ERR_EVAL_DISABLED","Code generation from strings disallowed for this context")})};return tsNodeVMEvaluateSource(loop,source,environment)}
			function:=tsFunc(invoke);function.name="";function.length=len(names);function.constructible=true;function.construct=func(loop *tsLoop,values ...tsValue)tsValue{return invoke(loop,values...)};return tsFunctionValue(function)
		}))
		module.set("isModuleNamespaceObject",tsNamedNativeMethod("isModuleNamespaceObject",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{value:=tsArg(args,0);if value.kind!=tsObjectKind{return tsBooleanValue(false)};return tsBooleanValue(tsNodeVMSingleton.namespaces[uintptr(value.ref)])}))
		module.set("measureMemory",tsNamedNativeMethod("measureMemory",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{options:=tsArg(args,0);if !tsIsUndefined(options)&&options.kind!=tsObjectKind&&options.kind!=tsInstanceKind{tsFSInvalid("options","an object")};mode:=tsGet(loop,options,tsStringReference(tsStringUTF8("mode")));if tsIsUndefined(mode){mode=tsStringReference(tsStringUTF8("summary"))};if mode.kind!=tsStringKind||(tsText(mode)!="summary"&&tsText(mode)!="detailed"){tsNodeVMInvalid("ERR_INVALID_ARG_VALUE","options.mode must be summary or detailed")};execution:=tsGet(loop,options,tsStringReference(tsStringUTF8("execution")));if tsIsUndefined(execution){execution=tsStringReference(tsStringUTF8("default"))};if execution.kind!=tsStringKind||(tsText(execution)!="default"&&tsText(execution)!="eager"){tsNodeVMInvalid("ERR_INVALID_ARG_VALUE","options.execution must be default or eager")};jsc:=tsNodeBunJSCModule(loop);memory:=tsCallReceiver(loop,tsGet(loop,jsc,tsStringReference(tsStringUTF8("memoryUsage"))),jsc);current:=tsGet(loop,memory,tsStringReference(tsStringUTF8("current")));peak:=tsGet(loop,memory,tsStringReference(tsStringUTF8("peak")));upper:=peak;if tsNumber(peak)<tsNumber(current){upper=current};measurement:=func()tsValue{value:=tsNewObject();value.set("jsMemoryEstimate",current);value.set("jsMemoryRange",tsArrayValue(&tsArray{values:[]tsValue{current,upper}}));return tsObjectValue(value)};out:=tsNewObject();out.set("total",measurement());if tsText(mode)=="detailed"{out.set("current",measurement());other:=&tsArray{};for range tsNodeVMSingleton.contexts{other.values=append(other.values,measurement())};out.set("other",tsArrayValue(other))};return tsPromiseValue(loop.resolved(tsObjectValue(out),false))}))
		module.set("default",tsObjectValue(module));tsNodeVMSingleton.module=tsObjectValue(module)
	})
	return tsNodeVMSingleton.module
}
func tsNodeNodeVMInvalidReceiver(method string) { panic(tsThrown{tsFSException("TypeError","ERR_INVALID_THIS","Script.prototype."+method+" called on an incompatible receiver")}) }
func tsNodeVMSyntheticCreate(loop *tsLoop,args []tsValue,prototype *tsObject) tsValue {
	names:=tsArg(args,0);callback:=tsArg(args,1);options:=tsArg(args,2);if names.kind!=tsArrayKind{tsFSInvalid("exportNames","an array of strings")};if callback.kind!=tsFunctionKind{tsFSInvalid("evaluateCallback","a function")};exports:=[]string{};seen:=map[string]bool{};for _,name:=range (*tsArray)(names.ref).values{if name.kind!=tsStringKind{tsFSInvalid("exportNames","an array of strings")};text:=tsText(name);if seen[text]{tsNodeVMInvalid("ERR_INVALID_ARG_VALUE","Duplicate export name")};seen[text]=true;exports=append(exports,text)};value:=tsNodeVMCreateModule(loop,"synthetic",tsU,exports,callback,options);(*tsObject)(value.ref).prototype=tsObjectValue(prototype);return value
}
`
