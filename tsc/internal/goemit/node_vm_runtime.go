package goemit

// nodeVMRuntime provides the public vm constructors and module lifecycle over
// a lightweight local context model. Bun's JSC binding is replaced by this
// in-process adapter; arbitrary engine compilation remains outside Go's runtime.
const nodeVMRuntime = `
type tsNodeVMModuleState struct {
	kind, identifier, status, source string
	context, namespace, failure, evaluateCallback tsValue
	exportNames, requests []string
	exports map[string]tsValue
	dependencies []tsValue
}
var tsNodeVMSingleton struct {
	sync.Once
	module, constants, scriptConstructor, moduleConstructor, sourceConstructor, syntheticConstructor tsValue
	contextKey, scriptKey, moduleKey, noErrorKey, perContextIDKey tsValue
	contexts map[uintptr]bool
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
	return tsTruthy(tsGet(loop,value,tsNodeVMSingleton.contextKey))
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
func tsNodeVMCreateContext(loop *tsLoop, context,options tsValue) tsValue {
	if tsNullish(context){context=tsObjectValue(tsNewObject())}
	if context.kind!=tsObjectKind&&context.kind!=tsInstanceKind{tsFSInvalid("contextObject","an object")}
	if !tsTruthy(tsGet(loop,context,tsNodeVMSingleton.contextKey)){tsSet(loop,context,tsNodeVMSingleton.contextKey,tsBooleanValue(true));tsNodeVMSingleton.contexts[uintptr(context.ref)]=true}
	if !tsIsUndefined(options){tsNodeVMOptions(options)}
	return context
}
func tsNodeVMScriptState(value tsValue) tsValue { if value.kind!=tsObjectKind{return tsU};return tsNodeVMSingleton.scripts[uintptr(value.ref)] }
func tsNodeVMNewScript(loop *tsLoop, source, options tsValue) tsValue {
	if !tsIsUndefined(source)&&source.kind!=tsStringKind{tsFSInvalid("code","a string")};opts:=tsNodeVMOptions(options)
	object:=tsNewObject();object.prototype=tsGet(loop,tsNodeVMSingleton.scriptConstructor,tsStringReference(tsStringUTF8("prototype")));value:=tsObjectValue(object)
	tsNodeVMSingleton.scripts[uintptr(value.ref)]=tsObjectValue(tsNewObject());state:=(*tsObject)(tsNodeVMSingleton.scripts[uintptr(value.ref)].ref);state.set("source",source);state.set("options",opts)
	return value
}
func tsNodeVMCreateModule(loop *tsLoop,kind string,source tsValue,exportNames []string,callback,options tsValue) tsValue {
	opts:=tsNodeVMOptions(options);context:=tsGet(loop,opts,tsStringReference(tsStringUTF8("context")));if !tsIsUndefined(context){tsInternalValidatorObject(context,tsStringReference(tsStringUTF8("options.context")),tsU);if !tsNodeVMIsContext(context){panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_TYPE","The \"options.context\" argument must be of type vm.Context")})}}
	identifier:=tsGet(loop,opts,tsStringReference(tsStringUTF8("identifier")));if !tsIsUndefined(identifier)&&identifier.kind!=tsStringKind{tsFSInvalid("options.identifier","a string")}
	if tsIsUndefined(identifier){if kind=="source"{tsNodeVMSingleton.globalID++;identifier=tsStringReference(tsStringUTF8("vm:module("+strconv.FormatUint(tsNodeVMSingleton.globalID-1,10)+")"))}else{identifier=tsStringReference(tsStringUTF8("vm:module(synthetic)"))}}
	if tsIsUndefined(context){context=tsGlobalObject()}
	prototype:=tsGet(loop,tsNodeVMSingleton.moduleConstructor,tsStringReference(tsStringUTF8("prototype")));if kind=="source"{prototype=tsGet(loop,tsNodeVMSingleton.sourceConstructor,tsStringReference(tsStringUTF8("prototype")))}else if kind=="synthetic"{prototype=tsGet(loop,tsNodeVMSingleton.syntheticConstructor,tsStringReference(tsStringUTF8("prototype")))}
	object:=tsNewObject();object.prototype=prototype;value:=tsObjectValue(object)
	state:=&tsNodeVMModuleState{kind:kind,identifier:tsText(identifier),context:context,source:tsText(source),status:"unlinked",failure:tsNodeVMSingleton.noErrorKey,exports:map[string]tsValue{},evaluateCallback:callback}
	if kind=="synthetic"{state.status="linked";state.exportNames=append([]string{},exportNames...)}else{state.requests=tsNodeVMImports(state.source)}
	tsNodeVMSingleton.modules[uintptr(value.ref)]=state
	if kind=="synthetic"{namespace:=tsNewObject();for _,name:=range state.exportNames{namespace.set(name,tsU)};state.namespace=tsObjectValue(namespace);tsNodeVMSingleton.namespaces[uintptr(namespace)]=true}
	return value
}
func tsNodeVMImports(source string) []string {
	imports:=[]string{};for offset:=0;offset<len(source);{index:=strings.Index(source[offset:],"from");if index<0{break};index+=offset+4;for index<len(source)&&(source[index]==' '||source[index]=='\t'||source[index]=='\n'||source[index]=='\r'){index++};if index>=len(source)||(source[index]!='\''&&source[index]!='"'){offset=index;continue};quote:=source[index];end:=strings.IndexByte(source[index+1:],quote);if end<0{break};specifier:=source[index+1:index+1+end];imports=append(imports,specifier);offset=index+end+2};return imports
}
func tsNodeVMInstallModuleMethods(loop *tsLoop, prototype *tsObject) {
	prototype.descriptors=map[string]*tsDescriptor{}
	for _,name:=range []string{"identifier","context","status","namespace","error","moduleRequests","dependencySpecifiers"}{property:=name;prototype.descriptors[property]=&tsDescriptor{accessor:true,configurable:true,getter:tsNamedNativeMethod("get "+property,0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeVMValidateModule(receiver);switch property{case "identifier":return tsStringReference(tsStringUTF8(state.identifier));case "context":return state.context;case "status":return tsStringReference(tsStringUTF8(state.status));case "namespace":if state.status=="unlinked"||state.status=="linking"{tsNodeVMInvalid("ERR_VM_MODULE_STATUS","Module namespace is not available before linking")};return state.namespace;case "error":if state.status!="errored"{tsNodeVMInvalid("ERR_VM_MODULE_STATUS","Module is not errored")};return state.failure;case "moduleRequests":items:=&tsArray{};for _,specifier:=range state.requests{item:=tsNewObject();item.set("specifier",tsStringReference(tsStringUTF8(specifier)));item.set("attributes",tsObjectValue(tsNewObject()));item.set("phase",tsStringReference(tsStringUTF8("evaluation")));items.values=append(items.values,tsObjectValue(item))};return tsArrayValue(items);case "dependencySpecifiers":items:=&tsArray{};for _,specifier:=range state.requests{items.values=append(items.values,tsStringReference(tsStringUTF8(specifier)))};return tsArrayValue(items)};return tsU})}
	}
	prototype.set("link",tsNamedNativeMethod("link",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeVMValidateModule(receiver);linker:=tsArg(args,0);if state.kind=="synthetic"{return tsU};if linker.kind!=tsFunctionKind{tsFSInvalid("linker","a function")};if state.status=="linked"||state.status=="evaluated"{return tsPromiseValue(loop.resolved(tsFSException("Error","ERR_VM_MODULE_ALREADY_LINKED","Module has already been linked"),true))};if state.status!="unlinked"{return tsPromiseValue(loop.resolved(tsFSException("Error","ERR_VM_MODULE_STATUS","Module must be unlinked"),true))};for _,specifier:=range state.requests{resolved:=tsCall(loop,linker,tsStringReference(tsStringUTF8(specifier)),receiver,tsObjectValue(tsNewObject()));if !tsNodeVMIsModule(resolved){return tsPromiseValue(loop.resolved(tsFSException("TypeError","ERR_VM_MODULE_NOT_MODULE","Linker must return a Module"),true))};dependency:=tsNodeVMValidateModule(resolved);if !tsStrictEqual(dependency.context,state.context){return tsPromiseValue(loop.resolved(tsFSException("Error","ERR_VM_MODULE_DIFFERENT_CONTEXT","Linked modules must use the same context"),true))};state.dependencies=append(state.dependencies,resolved)};state.status="linked";return tsPromiseValue(loop.resolved(tsU,false))}))
	prototype.set("instantiate",tsNamedNativeMethod("instantiate",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeVMValidateModule(receiver);if state.status=="unlinked"{state.status="linked"};return tsU}))
	prototype.set("evaluate",tsNamedNativeMethod("evaluate",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeVMValidateModule(receiver);if state.status!="linked"&&state.status!="evaluated"&&state.status!="errored"{return tsPromiseValue(loop.resolved(tsFSException("Error","ERR_VM_MODULE_STATUS","Module must be linked before evaluation"),true))};if state.status=="errored"{return tsPromiseValue(loop.resolved(state.failure,true))};if state.status=="evaluated"{return tsPromiseValue(loop.resolved(tsU,false))};if state.kind=="synthetic"&&state.evaluateCallback.kind==tsFunctionKind{func(){defer func(){if failure:=recover();failure!=nil{state.failure=tsUnwrap(failure);state.status="errored"}}();tsCall(loop,state.evaluateCallback,receiver)}()};if state.kind=="source"{tsNodeVMEvaluateSource(loop,state.source,state.context)};if state.status!="errored"{state.status="evaluated"};if state.status=="errored"{return tsPromiseValue(loop.resolved(state.failure,true))};return tsPromiseValue(loop.resolved(tsU,false))}))
	prototype.set("linkRequests",tsNamedNativeMethod("linkRequests",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeVMValidateModule(receiver);modules:=tsArg(args,0);if modules.kind!=tsArrayKind{tsFSInvalid("modules","an array")};items:=(*tsArray)(modules.ref);if len(items.values)!=len(state.requests){tsNodeVMInvalid("ERR_MODULE_LINK_MISMATCH","Module count does not match requests")};for _,item:=range items.values{if !tsNodeVMIsModule(item){tsNodeVMInvalid("ERR_VM_MODULE_NOT_MODULE","Expected a Module")};dependency:=tsNodeVMValidateModule(item);if !tsStrictEqual(dependency.context,state.context){tsNodeVMInvalid("ERR_VM_MODULE_DIFFERENT_CONTEXT","Linked modules must use the same context")};state.dependencies=append(state.dependencies,item)};return tsU}))
	prototype.set("hasAsyncGraph",tsNamedNativeMethod("hasAsyncGraph",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeVMValidateModule(receiver);if state.status=="unlinked"{tsNodeVMInvalid("ERR_VM_MODULE_STATUS","Module must be instantiated")};return tsBooleanValue(false)}))
	prototype.set("hasTopLevelAwait",tsNamedNativeMethod("hasTopLevelAwait",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{return tsBooleanValue(strings.Contains(tsNodeVMValidateModule(receiver).source,"await "))}))
	prototype.set("createCachedData",tsNamedNativeMethod("createCachedData",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeVMValidateModule(receiver);if state.status=="evaluating"||state.status=="evaluated"||state.status=="errored"{tsNodeVMInvalid("ERR_VM_MODULE_CANNOT_CREATE_CACHED_DATA","Cannot create cached data after evaluation")};return tsNodeBuffer(nil)}))
	prototype.set("setExport",tsNamedNativeMethod("setExport",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeVMValidateModule(receiver);name:=tsArg(args,0);if name.kind!=tsStringKind{tsFSInvalid("name","a string")};if state.kind!="synthetic"{tsNodeVMInvalid("ERR_INVALID_ARG_TYPE","setExport is only available on SyntheticModule")};found:=false;for _,item:=range state.exportNames{if item==tsText(name){found=true;break}};if !found{tsNodeVMInvalid("ERR_INVALID_ARG_VALUE","Export name is not declared")};state.exports[tsText(name)]=tsArg(args,1);tsSet(loop,state.namespace,name,tsArg(args,1));return tsU}))
}
func tsNodeVMIsModule(value tsValue) bool { if value.kind!=tsObjectKind{return false};return tsNodeVMSingleton.modules[uintptr(value.ref)]!=nil }
func tsNodeVMModule(loop *tsLoop) tsValue {
	tsNodeVMSingleton.Do(func(){
		tsNodeVMSingleton.contexts=map[uintptr]bool{};tsNodeVMSingleton.scripts=map[uintptr]tsValue{};tsNodeVMSingleton.modules=map[uintptr]*tsNodeVMModuleState{};tsNodeVMSingleton.namespaces=map[uintptr]bool{}
		tsNodeVMSingleton.contextKey=tsNewSymbol(tsStringUTF8("vm-context"));tsNodeVMSingleton.scriptKey=tsNewSymbol(tsStringUTF8("vm-script"));tsNodeVMSingleton.moduleKey=tsNewSymbol(tsStringUTF8("vm-module"));tsNodeVMSingleton.noErrorKey=tsNewSymbol(tsStringUTF8("vm-no-error"));tsNodeVMSingleton.perContextIDKey=tsNewSymbol(tsStringUTF8("vm-module-id"))
		scriptPrototype:=tsNewObject();scriptPrototype.prototype=tsObjectPrototype()
		script:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeVMNewScript(loop,tsArg(args,0),tsArg(args,1))});script.name="Script";script.length=2;script.constructible=true;script.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeVMNewScript(loop,tsArg(args,0),tsArg(args,1))};script.properties=tsNewObject();script.properties.set("prototype",tsObjectValue(scriptPrototype))
		for _,entry:=range []struct{name string;length int}{{"runInThisContext",1},{"runInContext",2},{"runInNewContext",2},{"createCachedData",0}}{name:=entry.name;scriptPrototype.set(name,tsNamedNativeMethod(name,entry.length,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeVMScriptState(receiver);if state.kind!=tsObjectKind{tsNodeNodeVMInvalidReceiver(name)};source:=tsText(tsGet(loop,state,tsStringReference(tsStringUTF8("source"))));switch name{case "runInThisContext":return tsNodeVMEvaluateSource(loop,source,tsGlobalObject());case "runInContext":context:=tsArg(args,0);tsNodeVMValidateContext(context);return tsNodeVMEvaluateSource(loop,source,context);case "runInNewContext":context:=tsNodeVMCreateContext(loop,tsArg(args,0),tsArg(args,1));return tsNodeVMEvaluateSource(loop,source,context);default:return tsNodeBuffer(nil)}}))}
		scriptPrototype.set("constructor",tsFunctionValue(script));tsNodeVMSingleton.scriptConstructor=tsFunctionValue(script)
		modulePrototype:=tsNewObject();modulePrototype.prototype=tsObjectPrototype();tsNodeVMInstallModuleMethods(loop,modulePrototype)
		base:=tsFunc(func(_ *tsLoop,_ ...tsValue)tsValue{tsNodeVMInvalid("ERR_INVALID_THIS","Module is not a constructor");return tsU});base.name="Module";base.constructible=true;base.construct=func(_ *tsLoop,_ ...tsValue)tsValue{tsNodeVMInvalid("ERR_INVALID_THIS","Module is not a constructor");return tsU};base.properties=tsNewObject();base.properties.set("prototype",tsObjectValue(modulePrototype));tsNodeVMSingleton.moduleConstructor=tsFunctionValue(base)
		sourcePrototype:=tsNewObject();sourcePrototype.prototype=tsObjectValue(modulePrototype);sourceCtor:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeVMCreateModule(loop,"source",tsArg(args,0),nil,tsU,tsArg(args,1))});sourceCtor.name="SourceTextModule";sourceCtor.length=2;sourceCtor.constructible=true;sourceCtor.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeVMCreateModule(loop,"source",tsArg(args,0),nil,tsU,tsArg(args,1))};sourceCtor.properties=tsNewObject();sourceCtor.properties.set("prototype",tsObjectValue(sourcePrototype));sourcePrototype.set("constructor",tsFunctionValue(sourceCtor));tsNodeVMSingleton.sourceConstructor=tsFunctionValue(sourceCtor)
		syntheticPrototype:=tsNewObject();syntheticPrototype.prototype=tsObjectValue(modulePrototype);syntheticCtor:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeVMSyntheticCreate(loop,args,syntheticPrototype)});syntheticCtor.name="SyntheticModule";syntheticCtor.length=2;syntheticCtor.constructible=true;syntheticCtor.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeVMSyntheticCreate(loop,args,syntheticPrototype)};syntheticCtor.properties=tsNewObject();syntheticCtor.properties.set("prototype",tsObjectValue(syntheticPrototype));syntheticPrototype.set("constructor",tsFunctionValue(syntheticCtor));tsNodeVMSingleton.syntheticConstructor=tsFunctionValue(syntheticCtor)
		constants:=tsNewObject();dont:=tsNewSymbol(tsStringUTF8("DONT_CONTEXTIFY"));loader:=tsNewSymbol(tsStringUTF8("USE_MAIN_CONTEXT_DEFAULT_LOADER"));constants.set("DONT_CONTEXTIFY",dont);constants.set("USE_MAIN_CONTEXT_DEFAULT_LOADER",loader);tsNodeVMSingleton.constants=tsObjectValue(constants)
		module:=tsNewObject();module.set("createContext",tsNamedNativeMethod("createContext",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeVMCreateContext(loop,tsArg(args,0),tsArg(args,1))}));module.set("isContext",tsNamedNativeMethod("isContext",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsBooleanValue(tsNodeVMIsContext(tsArg(args,0)))}))
		module.set("runInContext",tsNamedNativeMethod("runInContext",3,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{context:=tsArg(args,1);tsNodeVMValidateContext(context);return tsNodeVMEvaluateSource(loop,tsText(tsArg(args,0)),context)}));module.set("runInThisContext",tsNamedNativeMethod("runInThisContext",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeVMEvaluateSource(loop,tsText(tsArg(args,0)),tsGlobalObject())}));module.set("runInNewContext",tsNamedNativeMethod("runInNewContext",3,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{context:=tsNodeVMCreateContext(loop,tsArg(args,1),tsArg(args,2));return tsNodeVMEvaluateSource(loop,tsText(tsArg(args,0)),context)}));module.set("createScript",tsNamedNativeMethod("createScript",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeVMNewScript(loop,tsArg(args,0),tsArg(args,1))}))
		module.set("Script",tsNodeVMSingleton.scriptConstructor);module.set("Module",tsNodeVMSingleton.moduleConstructor);module.set("SourceTextModule",tsNodeVMSingleton.sourceConstructor);module.set("SyntheticModule",tsNodeVMSingleton.syntheticConstructor);module.set("constants",tsObjectValue(constants))
		module.set("compileFunction",tsNamedNativeMethod("compileFunction",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{body:=tsArg(args,0);if body.kind!=tsStringKind{tsFSInvalid("code","a string")};params:=tsArg(args,1);names:=[]string{};if tsIsUndefined(params){params=tsArrayValue(&tsArray{})};if params.kind!=tsArrayKind{tsFSInvalid("params","an array")};for _,value:=range (*tsArray)(params.ref).values{if value.kind!=tsStringKind{tsFSInvalid("params","an array of strings")};names=append(names,tsText(value))};source:=tsText(body);return tsFunctionValue(tsFunc(func(loop *tsLoop,values ...tsValue)tsValue{context:=tsObjectValue(tsNewObject());for index,name:=range names{tsSet(loop,context,tsStringReference(tsStringUTF8(name)),tsArg(values,index))};return tsNodeVMEvaluateSource(loop,source,context)}))}))
		module.set("isModuleNamespaceObject",tsNamedNativeMethod("isModuleNamespaceObject",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{value:=tsArg(args,0);if value.kind!=tsObjectKind{return tsBooleanValue(false)};return tsBooleanValue(tsNodeVMSingleton.namespaces[uintptr(value.ref)])}))
		module.set("measureMemory",tsNamedNativeMethod("measureMemory",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{options:=tsArg(args,0);if !tsIsUndefined(options)&&options.kind!=tsObjectKind{tsFSInvalid("options","an object")};jsc:=tsNodeBunJSCModule(loop);memory:=tsCallReceiver(loop,tsGet(loop,jsc,tsStringReference(tsStringUTF8("memoryUsage"))),jsc);current:=tsGet(loop,memory,tsStringReference(tsStringUTF8("current")));peak:=tsGet(loop,memory,tsStringReference(tsStringUTF8("peak")));measurement:=tsNewObject();measurement.set("jsMemoryEstimate",current);measurement.set("jsMemoryRange",tsArrayValue(&tsArray{values:[]tsValue{current,peak}}));out:=tsNewObject();out.set("total",tsObjectValue(measurement));return tsPromiseValue(loop.resolved(tsObjectValue(out),false))}))
		module.set("default",tsObjectValue(module));tsNodeVMSingleton.module=tsObjectValue(module)
	})
	return tsNodeVMSingleton.module
}
func tsNodeNodeVMInvalidReceiver(method string) { panic(tsThrown{tsFSException("TypeError","ERR_INVALID_THIS","Script.prototype."+method+" called on an incompatible receiver")}) }
func tsNodeVMSyntheticCreate(loop *tsLoop,args []tsValue,prototype *tsObject) tsValue {
	names:=tsArg(args,0);callback:=tsArg(args,1);options:=tsArg(args,2);if names.kind!=tsArrayKind{tsFSInvalid("exportNames","an array of strings")};if callback.kind!=tsFunctionKind{tsFSInvalid("evaluateCallback","a function")};exports:=[]string{};seen:=map[string]bool{};for _,name:=range (*tsArray)(names.ref).values{if name.kind!=tsStringKind{tsFSInvalid("exportNames","an array of strings")};text:=tsText(name);if seen[text]{tsNodeVMInvalid("ERR_INVALID_ARG_VALUE","Duplicate export name")};seen[text]=true;exports=append(exports,text)};value:=tsNodeVMCreateModule(loop,"synthetic",tsU,exports,callback,options);(*tsObject)(value.ref).prototype=tsObjectValue(prototype);return value
}
`
