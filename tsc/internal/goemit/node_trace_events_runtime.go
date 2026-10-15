package goemit

// nodeTraceEventsRuntime implements Bun's public Tracing wrapper over the
// shared internal trace-event agent.
const nodeTraceEventsRuntime = `
type tsNodeTracing struct { object *tsObject; categories []string; enabled bool }
type tsNodeTraceEventsState struct { sync.Once; value *tsObject; tracings map[*tsObject]*tsNodeTracing; counts map[string]int }
var tsNodeTraceEventsSingleton tsNodeTraceEventsState

func tsNodeTraceEventsInit() {
	tsNodeTraceEventsSingleton.Do(func(){
		tsNodeTraceEventsSingleton.tracings=map[*tsObject]*tsNodeTracing{}
		tsNodeTraceEventsSingleton.counts=map[string]int{}
	})
}

func tsNodeTraceSetEnabled(loop *tsLoop,tracing *tsNodeTracing,enabled bool) {
	if tracing.enabled==enabled{return}
	tracing.enabled=enabled
	tracing.object.set("enabled",tsBooleanValue(enabled))
	categories:=make([]tsValue,len(tracing.categories));for index,category:=range tracing.categories{categories[index]=tsStringReference(tsStringUTF8(category))}
	if enabled{tsInternalTraceEnablePublic(loop,categories)}else{tsInternalTraceDisable(categories)}
}

func tsNodeTraceCreate(loop *tsLoop,options tsValue) tsValue {
	if options.kind!=tsObjectKind&&options.kind!=tsInstanceKind{tsFSInvalid("options","an object")}
	value:=tsFSProperty(options,"categories");if value.kind!=tsArrayKind{tsFSInvalid("options.categories","a string array")}
	array:=(*tsArray)(value.ref);if array.length()==0{tsPropertyFailure("At least one category is required")}
	categories:=make([]string,array.length())
	for index:=0;index<array.length();index++{item:=array.at(index);if item.kind!=tsStringKind{tsFSInvalid("options.categories","a string array")};categories[index]=tsText(item)}
	object:=tsNewObject();state:=&tsNodeTracing{object:object,categories:categories};tsNodeTraceEventsSingleton.tracings[object]=state
	joined:=strings.Join(categories,",")
	object.set("enabled",tsBooleanValue(false));object.set("categories",tsStringReference(tsStringUTF8(joined)))
	object.descriptors=map[string]*tsDescriptor{"enabled":{writable:false,enumerable:true,configurable:true},"categories":{writable:false,enumerable:true,configurable:true}}
	object.set("enable",tsNamedNativeMethod("enable",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{if tracing:=tsNodeTraceEventsSingleton.tracings[(*tsObject)(receiver.ref)];tracing!=nil{tsNodeTraceSetEnabled(loop,tracing,true)};return receiver}))
	object.set("disable",tsNamedNativeMethod("disable",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{if tracing:=tsNodeTraceEventsSingleton.tracings[(*tsObject)(receiver.ref)];tracing!=nil{tsNodeTraceSetEnabled(loop,tracing,false)};return receiver}))
	return tsObjectValue(object)
}

func tsNodeTraceEventsModule(loop *tsLoop) tsValue {
	tsNodeTraceEventsInit()
	tsInternalTraceEventsInit(loop)
	if tsNodeTraceEventsSingleton.value!=nil{return tsObjectValue(tsNodeTraceEventsSingleton.value)}
	module:=tsNewObject()
	module.set("createTracing",tsNamedNativeMethod("createTracing",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeTraceCreate(loop,tsArg(args,0))}))
	module.set("getEnabledCategories",tsNamedNativeMethod("getEnabledCategories",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsInternalTraceEnabledString()}))
	tsNodeTraceEventsSingleton.value=module
	return tsObjectValue(module)
}
`
