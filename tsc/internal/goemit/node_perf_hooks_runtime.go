package goemit

// nodePerfHooksRuntime implements the Node wrapper with Go's monotonic clock.
// Entry observers and histogram buckets are local prototype implementations.
const nodePerfHooksRuntime = `
type tsNodePerfObserver struct {
	callback tsValue
	object *tsObject
	types map[string]bool
	entries []tsValue
	queued bool
	active bool
}

type tsNodePerfState struct {
	sync.Once
	value *tsObject
	performance *tsObject
	started time.Time
	timeOrigin float64
	marks map[string]float64
	entries []tsValue
	observers []*tsNodePerfObserver
	nodeTiming tsValue
}

var tsNodePerfSingleton tsNodePerfState

func tsNodePerfNow() float64 {
	if tsNodePerfSingleton.started.IsZero() { return 0 }
	return float64(time.Since(tsNodePerfSingleton.started).Nanoseconds()) / 1e6
}

func tsNodePerfEntry(name, entryType string, start, duration float64, detail tsValue) tsValue {
	out := tsNewObject()
	out.nativeClass = tsNativeClass("PerformanceEntry")
	out.set("name", tsStringReference(tsStringUTF8(name)))
	out.set("entryType", tsStringReference(tsStringUTF8(entryType)))
	out.set("startTime", tsNumberValue(start))
	out.set("duration", tsNumberValue(duration))
	if tsIsUndefined(detail) { detail = tsNull }
	out.set("detail", detail)
	out.set("toJSON", tsNamedNativeMethod("toJSON", 0, func(_ *tsLoop, receiver tsValue, _ ...tsValue) tsValue {
		out := tsNewObject()
		for _, key := range []string{"name", "entryType", "startTime", "duration", "detail"} { out.set(key, tsGet(receiver, tsStringReference(tsStringUTF8(key)))) }
		return tsObjectValue(out)
	}))
	return tsObjectValue(out)
}

func tsNodePerfEntryList(entries []tsValue) tsValue {
	list := entries
	out := tsNewObject()
	filter := func(args []tsValue, byName bool) tsValue {
		name, kind := tsArg(args, 0), tsArg(args, 1)
		result := &tsArray{}
		for index := 0; index < len(list); index++ {
			entry := list[index]
			if byName && !tsIsUndefined(name) && !tsStrictEqual(tsGet(entry, tsStringReference(tsStringUTF8("name"))), name) { continue }
			if !tsIsUndefined(kind) && !tsStrictEqual(tsGet(entry, tsStringReference(tsStringUTF8("entryType"))), kind) { continue }
			result.values = append(result.values, entry)
		}
		return tsArrayValue(result)
	}
	out.set("getEntries", tsNamedNativeMethod("getEntries", 0, func(_ *tsLoop, _ tsValue, _ ...tsValue) tsValue { return tsArrayValue(&tsArray{values:list}) }))
	out.set("getEntriesByName", tsNamedNativeMethod("getEntriesByName", 1, func(_ *tsLoop, _ tsValue, args ...tsValue) tsValue { return filter(args, true) }))
	out.set("getEntriesByType", tsNamedNativeMethod("getEntriesByType", 1, func(_ *tsLoop, _ tsValue, args ...tsValue) tsValue { return filter([]tsValue{tsU, tsArg(args, 0)}, false) }))
	return tsObjectValue(out)
}

func tsNodePerfNotify(loop *tsLoop, entry tsValue) {
	state := &tsNodePerfSingleton
	state.entries = append(state.entries, entry)
	typeName := tsText(tsGet(entry, tsStringReference(tsStringUTF8("entryType"))))
	for _, observer := range state.observers {
		if !observer.active || !observer.types[typeName] { continue }
		observer.entries = append(observer.entries, entry)
		if observer.queued { continue }
		observer.queued = true
		loop.post(func() {
			if !observer.active { observer.entries = nil; observer.queued = false; return }
			entries := observer.entries
			observer.entries = nil
			observer.queued = false
			if len(entries) > 0 { tsCall(loop, observer.callback, tsNodePerfEntryList(entries), tsObjectValue(observer.object)) }
		})
	}
}

func tsNodePerfGetEntries(name, entryType string) tsValue {
	result := &tsArray{}
	for _, entry := range tsNodePerfSingleton.entries {
		if name != "" && tsText(tsGet(entry, tsStringReference(tsStringUTF8("name")))) != name { continue }
		if entryType != "" && tsText(tsGet(entry, tsStringReference(tsStringUTF8("entryType")))) != entryType { continue }
		result.values = append(result.values, entry)
	}
	sort.SliceStable(result.values, func(i, j int) bool {
		return tsNumber(tsGet(result.values[i], tsStringReference(tsStringUTF8("startTime")))) < tsNumber(tsGet(result.values[j], tsStringReference(tsStringUTF8("startTime"))))
	})
	return tsArrayValue(result)
}

func tsNodePerfMark(loop *tsLoop, args []tsValue) tsValue {
	nameValue := tsArg(args, 0)
	if nameValue.kind != tsStringKind { tsFSInvalid("name", "a string") }
	name := tsStringValue(nameValue).String()
	options := tsArg(args, 1)
	if !tsIsUndefined(options) && (options.kind != tsObjectKind && options.kind != tsInstanceKind) { tsFSInvalid("options", "an object") }
	start := tsNodePerfNow()
	if explicit := tsFSProperty(options, "startTime"); !tsIsUndefined(explicit) {
		if !tsIsNumeric(explicit) { tsFSInvalid("options.startTime", "a number") }
		start = tsNumber(explicit)
	}
	detail := tsFSProperty(options, "detail")
	tsNodePerfSingleton.marks[name] = start
	entry := tsNodePerfEntry(name, "mark", start, 0, detail)
	tsNodePerfNotify(loop, entry)
	return entry
}

func tsNodePerfMeasure(loop *tsLoop, args []tsValue) tsValue {
	nameValue := tsArg(args, 0)
	if nameValue.kind != tsStringKind { tsFSInvalid("name", "a string") }
	name := tsStringValue(nameValue).String()
	start, end := 0.0, tsNodePerfNow()
	options := tsArg(args, 1)
	if options.kind == tsStringKind {
		mark, ok := tsNodePerfSingleton.marks[tsText(options)]
		if !ok { panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"SyntaxError",message:"The mark '"+tsText(options)+"' does not exist"})}) }
		start = mark
	} else if options.kind == tsObjectKind || options.kind == tsInstanceKind {
		if value := tsFSProperty(options, "start"); !tsIsUndefined(value) {
			if value.kind == tsStringKind { mark, ok := tsNodePerfSingleton.marks[tsText(value)]; if !ok { tsPropertyFailure("The start mark does not exist") }; start = mark
			} else if tsIsNumeric(value) { start = tsNumber(value) } else { tsFSInvalid("options.start", "a string or number") }
		}
		if value := tsFSProperty(options, "end"); !tsIsUndefined(value) {
			if value.kind == tsStringKind { mark, ok := tsNodePerfSingleton.marks[tsText(value)]; if !ok { tsPropertyFailure("The end mark does not exist") }; end = mark
			} else if tsIsNumeric(value) { end = tsNumber(value) } else { tsFSInvalid("options.end", "a string or number") }
		}
		if value := tsFSProperty(options, "startTime"); !tsIsUndefined(value) { if !tsIsNumeric(value) { tsFSInvalid("options.startTime", "a number") }; start = tsNumber(value) }
		if value := tsFSProperty(options, "duration"); !tsIsUndefined(value) { if !tsIsNumeric(value) { tsFSInvalid("options.duration", "a number") }; end = start + tsNumber(value) }
	} else if !tsIsUndefined(options) {
		tsFSInvalid("options", "a string or object")
	}
	if !tsIsUndefined(tsArg(args, 2)) && tsArg(args, 2).kind == tsStringKind {
		mark, ok := tsNodePerfSingleton.marks[tsText(tsArg(args, 2))]
		if !ok { tsPropertyFailure("The end mark does not exist") }
		end = mark
	}
	entry := tsNodePerfEntry(name, "measure", start, end-start, tsFSProperty(options, "detail"))
	tsNodePerfNotify(loop, entry)
	return entry
}

func tsNodePerfUtilization() tsValue {
	out := tsNewObject()
	out.set("idle", tsNumberValue(0))
	out.set("active", tsNumberValue(0))
	out.set("utilization", tsNumberValue(0))
	return tsObjectValue(out)
}

func tsNodePerfObserverConstructor() tsValue {
	function := tsFunc(func(_ *tsLoop, _ ...tsValue) tsValue { tsPropertyFailure("PerformanceObserver must be constructed with new"); return tsU })
	function.name = "PerformanceObserver"
	function.length = 1
	function.constructible = true
	function.construct = func(loop *tsLoop, args ...tsValue) tsValue {
		callback := tsArg(args, 0)
		if callback.kind != tsFunctionKind { tsFSInvalid("callback", "a function") }
		object := tsNewObject()
		object.nativeClass = tsNativeClass("PerformanceObserver")
		observer := &tsNodePerfObserver{callback:callback,object:object,types:map[string]bool{},active:true}
		tsNodePerfSingleton.observers = append(tsNodePerfSingleton.observers, observer)
		object.set("observe", tsNamedNativeMethod("observe", 1, func(_ *tsLoop, _ tsValue, args ...tsValue) tsValue {
			options := tsArg(args, 0)
			if options.kind != tsObjectKind && options.kind != tsInstanceKind { tsFSInvalid("options", "an object") }
			entryTypes := tsFSProperty(options, "entryTypes")
			if entryTypes.kind == tsArrayKind {
				for _, entryType := range (*tsArray)(entryTypes.ref).values { observer.types[tsText(entryType)] = true }
			} else if typeValue := tsFSProperty(options, "type"); !tsIsUndefined(typeValue) {
				observer.types[tsText(typeValue)] = true
			} else {
				tsPropertyFailure("The options.entryTypes property must be an array")
			}
			observer.active = true
			return tsU
		}))
		object.set("disconnect", tsNamedNativeMethod("disconnect", 0, func(_ *tsLoop, _ tsValue, _ ...tsValue) tsValue { observer.active = false; observer.entries = nil; observer.types = map[string]bool{}; return tsU }))
		object.set("takeRecords", tsNamedNativeMethod("takeRecords", 0, func(_ *tsLoop, _ tsValue, _ ...tsValue) tsValue { entries := observer.entries; observer.entries = nil; return tsArrayValue(&tsArray{values:entries}) }))
		return tsObjectValue(object)
	}
	function.properties = tsNewObject()
	function.properties.set("supportedEntryTypes", tsArrayValue(&tsArray{values:[]tsValue{tsStringReference(tsStringUTF8("function")),tsStringReference(tsStringUTF8("mark")),tsStringReference(tsStringUTF8("measure"))}}))
	return tsFunctionValue(function)
}

func tsNodePerfTimerify(loop *tsLoop, callback tsValue) tsValue {
	if callback.kind != tsFunctionKind { tsFSInvalid("fn", "a function") }
	name := "timerified " + (*tsFunction)(callback.ref).name
	return tsNamedNativeMethod(name, (*tsFunction)(callback.ref).length, func(loop *tsLoop, receiver tsValue, args ...tsValue) tsValue {
		start := tsNodePerfNow()
		result := tsCallReceiver(loop, callback, receiver, args...)
		record := func() { entry := tsNodePerfEntry((*tsFunction)(callback.ref).name, "function", start, tsNodePerfNow()-start, tsArrayValue(&tsArray{values:append([]tsValue(nil),args...)})); tsNodePerfNotify(loop, entry) }
		if result.kind == tsPromiseKind { (*tsPromise)(result.ref).then(func(_ tsResult) { record() }) } else { record() }
		return result
	})
}

type tsNodeHistogram struct { values []uint64; previous uint64; exceeds uint64 }

func tsNodeHistogramValue(histogram *tsNodeHistogram) tsValue {
	out := tsNewObject()
	percentile := func(p float64) uint64 {
		if len(histogram.values) == 0 { return 0 }
		values := append([]uint64(nil), histogram.values...)
		sort.Slice(values, func(i,j int) bool { return values[i] < values[j] })
		if p < 0 { p = 0 }; if p > 100 { p = 100 }
		index := int(math.Ceil(p/100*float64(len(values))))-1; if index < 0 { index = 0 }; return values[index]
	}
	stats := func() (uint64,uint64,float64,float64) {
		if len(histogram.values)==0 { return 0,0,0,0 }
		min,max := histogram.values[0],histogram.values[0]; mean:=float64(0)
		for _,value:=range histogram.values { if value<min {min=value}; if value>max {max=value}; mean+=float64(value) }
		mean/=float64(len(histogram.values)); variance:=float64(0); for _,value:=range histogram.values { delta:=float64(value)-mean; variance+=delta*delta }; variance/=float64(len(histogram.values))
		return min,max,mean,math.Sqrt(variance)
	}
	record := tsNamedNativeMethod("record", 1, func(_ *tsLoop, _ tsValue, args ...tsValue) tsValue {
		value:=tsArg(args,0); var number uint64
		if value.kind==tsBigIntKind { number=(*big.Int)(value.ref).Uint64() } else if tsIsNumeric(value) { n:=tsNumber(value); if n<0||math.IsNaN(n)||math.IsInf(n,0) { tsPropertyFailure("Histogram value must be a finite non-negative number") }; number=uint64(n) } else { tsFSInvalid("val", "a number or bigint") }
		histogram.values=append(histogram.values,number); return tsU
	})
	out.set("record",record)
	out.set("recordDelta",tsNamedNativeMethod("recordDelta",0,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{value:=tsArg(args,0);number:=uint64(0);if value.kind==tsBigIntKind{number=(*big.Int)(value.ref).Uint64()}else if tsIsNumeric(value){n:=tsNumber(value);if n>0&&!math.IsInf(n,0)&&!math.IsNaN(n){number=uint64(n)}};delta:=number-histogram.previous;histogram.previous=number;histogram.values=append(histogram.values,delta);return tsU}))
	out.set("reset",tsNamedNativeMethod("reset",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{histogram.values=nil;histogram.previous=0;histogram.exceeds=0;return tsU}))
	out.set("percentile",tsNamedNativeMethod("percentile",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNumberValue(float64(percentile(tsNumber(tsArg(args,0)))))}))
	out.set("percentileBigInt",tsNamedNativeMethod("percentileBigInt",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsBigIntValue(new(big.Int).SetUint64(percentile(tsNumber(tsArg(args,0)))))}))
	out.descriptors=map[string]*tsDescriptor{}
	for _,name:=range []string{"count","min","max","mean","stddev","exceeds","minBigInt","maxBigInt","meanBigInt","percentiles"} {
		key:=name
		getter:=tsNamedNativeMethod("get "+key,0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{
			min,max,mean,stddev:=stats()
			switch key {
			case "count":return tsNumberValue(float64(len(histogram.values)))
			case "min":return tsNumberValue(float64(min))
			case "max":return tsNumberValue(float64(max))
			case "mean":return tsNumberValue(mean)
			case "stddev":return tsNumberValue(stddev)
			case "exceeds":return tsNumberValue(float64(histogram.exceeds))
			case "minBigInt":return tsBigIntValue(new(big.Int).SetUint64(min))
			case "maxBigInt":return tsBigIntValue(new(big.Int).SetUint64(max))
			case "meanBigInt":return tsBigIntValue(new(big.Int).SetUint64(uint64(mean)))
			default:out:=tsNewCollection(false);for _,p:=range []float64{0,50,75,90,99,100}{out.put(tsNumberValue(p),tsNumberValue(float64(percentile(p))))};return tsCollectionValue(out)
			}
		})
		out.order=append(out.order,key);out.descriptors[key]=&tsDescriptor{accessor:true,enumerable:true,configurable:true,getter:getter}
	}
	return tsObjectValue(out)
}

func tsNodePerfHistogram(loop *tsLoop,args []tsValue) tsValue {
	options:=tsArg(args,0)
	if !tsIsUndefined(options)&&(options.kind!=tsObjectKind&&options.kind!=tsInstanceKind){tsFSInvalid("options","an object")}
	figures:=tsFSProperty(options,"figures")
	if !tsIsUndefined(figures)&&(!tsIsNumeric(figures)||tsNumber(figures)<1||tsNumber(figures)>5){tsPropertyFailure("options.figures must be between 1 and 5")}
	return tsNodeHistogramValue(&tsNodeHistogram{})
}

func tsNodePerfMonitorEventLoopDelay(loop *tsLoop,args []tsValue) tsValue {
	value:=tsNodePerfHistogram(loop,args)
	object:=(*tsObject)(value.ref)
	enabled:=false
	object.set("enable",tsNamedNativeMethod("enable",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{if enabled{return tsBooleanValue(false)};enabled=true;return tsBooleanValue(true)}))
	object.set("disable",tsNamedNativeMethod("disable",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{if !enabled{return tsBooleanValue(false)};enabled=false;return tsBooleanValue(true)}))
	dispose:=tsWellKnownSymbol("dispose")
	if !tsIsUndefined(dispose){object.set(tsPropertyKey(dispose),object.get("disable"))}
	return value
}

func tsNodePerfConstructor(name string, construct func(*tsLoop,...tsValue)tsValue) tsValue {
	function:=tsFunc(func(_ *tsLoop,_ ...tsValue)tsValue{tsPropertyFailure(name+" must be constructed with new");return tsU})
	function.name=name;function.constructible=true;function.construct=construct
	return tsFunctionValue(function)
}

func tsNodePerfHooksModule(loop *tsLoop) tsValue {
	tsNodePerfSingleton.Do(func(){
		wallStart:=time.Now()
		tsNodePerfSingleton.started=wallStart
		tsNodePerfSingleton.timeOrigin=float64(wallStart.UnixNano())/1e6
		tsNodePerfSingleton.marks=map[string]float64{}
		performance:=tsNewObject()
		performance.set("now",tsNamedNativeMethod("now",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNumberValue(tsNodePerfNow())}))
		performance.set("timeOrigin",tsNumberValue(tsNodePerfSingleton.timeOrigin))
		performance.set("toJSON",tsNamedNativeMethod("toJSON",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{out:=tsNewObject();out.set("timeOrigin",tsNumberValue(tsNodePerfSingleton.timeOrigin));return tsObjectValue(out)}))
		performance.set("mark",tsNamedNativeMethod("mark",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodePerfMark(loop,args)}))
		performance.set("measure",tsNamedNativeMethod("measure",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodePerfMeasure(loop,args)}))
		for _,method:=range []string{"getEntries","getEntriesByName","getEntriesByType"}{name:=method;performance.set(name,tsNamedNativeMethod(name,0,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{switch name{case "getEntries":return tsNodePerfGetEntries("","");case "getEntriesByName":return tsNodePerfGetEntries(tsStringValue(tsArg(args,0)).String(),"");default:return tsNodePerfGetEntries("",tsStringValue(tsArg(args,0)).String())}}))}
		performance.set("clearMarks",tsNamedNativeMethod("clearMarks",0,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{name:=tsArg(args,0);kept:=[]tsValue{};for _,entry:=range tsNodePerfSingleton.entries{entryName:=tsGet(entry,tsStringReference(tsStringUTF8("name")));entryType:=tsText(tsGet(entry,tsStringReference(tsStringUTF8("entryType"))));if entryType!="mark"||!tsIsUndefined(name)&&!tsStrictEqual(entryName,name){kept=append(kept,entry)}};tsNodePerfSingleton.entries=kept;if tsIsUndefined(name){tsNodePerfSingleton.marks=map[string]float64{}}else{delete(tsNodePerfSingleton.marks,tsText(name))};return tsU}))
		performance.set("clearMeasures",tsNamedNativeMethod("clearMeasures",0,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{name:=tsArg(args,0);kept:=[]tsValue{};for _,entry:=range tsNodePerfSingleton.entries{if tsText(tsGet(entry,tsStringReference(tsStringUTF8("entryType"))))!="measure"||!tsIsUndefined(name)&&!tsStrictEqual(tsGet(entry,tsStringReference(tsStringUTF8("name"))),name){kept=append(kept,entry)}};tsNodePerfSingleton.entries=kept;return tsU}))
		performance.set("eventLoopUtilization",tsNamedNativeMethod("eventLoopUtilization",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNodePerfUtilization()}))
		performance.set("timerify",tsNamedNativeMethod("timerify",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodePerfTimerify(loop,tsArg(args,0))}))
		nodeTiming:=tsNewObject();nodeTiming.nativeClass=tsNativeClass("PerformanceNodeTiming")
		for _,key:=range []string{"name","entryType","startTime","duration","bootstrapComplete","environment","idleTime","loopExit","loopStart","nodeStart","v8Start"}{value:=tsNumberValue(0);switch key{case "name","entryType":value=tsStringReference(tsStringUTF8("node"));case "duration":value=tsNumberValue(tsNodePerfNow());case "startTime","bootstrapComplete","environment","nodeStart","v8Start":value=tsNumberValue(tsNodePerfSingleton.timeOrigin);case "idleTime","loopStart":value=tsNumberValue(1);case "loopExit":value=tsNumberValue(-1)};nodeTiming.set(key,value)}
		nodeTiming.set("toJSON",tsNamedNativeMethod("toJSON",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsObjectValue(nodeTiming)}))
		performance.set("nodeTiming",tsObjectValue(nodeTiming));tsNodePerfSingleton.nodeTiming=tsObjectValue(nodeTiming);tsNodePerfSingleton.performance=performance
		module:=tsNewObject()
		module.set("performance",tsObjectValue(performance))
		constants:=tsNewObject();for name,value:=range map[string]float64{"NODE_PERFORMANCE_ENTRY_TYPE_DNS":4,"NODE_PERFORMANCE_ENTRY_TYPE_GC":0,"NODE_PERFORMANCE_ENTRY_TYPE_HTTP":1,"NODE_PERFORMANCE_ENTRY_TYPE_HTTP2":2,"NODE_PERFORMANCE_ENTRY_TYPE_NET":3,"NODE_PERFORMANCE_GC_FLAGS_ALL_AVAILABLE_GARBAGE":16,"NODE_PERFORMANCE_GC_FLAGS_ALL_EXTERNAL_MEMORY":32,"NODE_PERFORMANCE_GC_FLAGS_CONSTRUCT_RETAINED":2,"NODE_PERFORMANCE_GC_FLAGS_FORCED":4,"NODE_PERFORMANCE_GC_FLAGS_NO":0,"NODE_PERFORMANCE_GC_FLAGS_SCHEDULE_IDLE":64,"NODE_PERFORMANCE_GC_FLAGS_SYNCHRONOUS_PHANTOM_PROCESSING":8,"NODE_PERFORMANCE_GC_INCREMENTAL":8,"NODE_PERFORMANCE_GC_MAJOR":4,"NODE_PERFORMANCE_GC_MINOR":1,"NODE_PERFORMANCE_GC_WEAKCB":16,"NODE_PERFORMANCE_MILESTONE_BOOTSTRAP_COMPLETE":7,"NODE_PERFORMANCE_MILESTONE_ENVIRONMENT":2,"NODE_PERFORMANCE_MILESTONE_LOOP_EXIT":6,"NODE_PERFORMANCE_MILESTONE_LOOP_START":5,"NODE_PERFORMANCE_MILESTONE_NODE_START":3,"NODE_PERFORMANCE_MILESTONE_TIME_ORIGIN_TIMESTAMP":0,"NODE_PERFORMANCE_MILESTONE_TIME_ORIGIN":1,"NODE_PERFORMANCE_MILESTONE_V8_START":4}{constants.set(name,tsNumberValue(value))};module.set("constants",tsObjectValue(constants))
		module.set("Performance",tsNodePerfConstructor("Performance",func(_ *tsLoop,_ ...tsValue)tsValue{return tsObjectValue(performance)}))
		module.set("PerformanceEntry",tsNodePerfConstructor("PerformanceEntry",func(_ *tsLoop,args ...tsValue)tsValue{return tsNodePerfEntry(tsStringValue(tsArg(args,0)).String(),tsStringValue(tsArg(args,1)).String(),tsNumber(tsArg(args,2)),tsNumber(tsArg(args,3)),tsU)}))
		module.set("PerformanceMark",tsNodePerfConstructor("PerformanceMark",func(_ *tsLoop,args ...tsValue)tsValue{return tsNodePerfEntry(tsStringValue(tsArg(args,0)).String(),"mark",tsNodePerfNow(),0,tsU)}))
		module.set("PerformanceMeasure",tsNodePerfConstructor("PerformanceMeasure",func(_ *tsLoop,args ...tsValue)tsValue{return tsNodePerfEntry(tsStringValue(tsArg(args,0)).String(),"measure",0,0,tsU)}))
		module.set("PerformanceObserver",tsNodePerfObserverConstructor())
		module.set("PerformanceObserverEntryList",tsNodePerfConstructor("PerformanceObserverEntryList",func(_ *tsLoop,_ ...tsValue)tsValue{return tsNodePerfEntryList(nil)}))
		module.set("PerformanceNodeTiming",tsNodePerfConstructor("PerformanceNodeTiming",func(_ *tsLoop,_ ...tsValue)tsValue{return tsObjectValue(nodeTiming)}))
		module.set("PerformanceResourceTiming",tsNodePerfConstructor("PerformanceResourceTiming",func(_ *tsLoop,_ ...tsValue)tsValue{return tsNodePerfEntry("","resource",0,0,tsU)}))
		module.set("eventLoopUtilization",tsNamedNativeMethod("eventLoopUtilization",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNodePerfUtilization()}))
		module.set("timerify",tsNamedNativeMethod("timerify",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodePerfTimerify(loop,tsArg(args,0))}))
		module.set("monitorEventLoopDelay",tsNamedNativeMethod("monitorEventLoopDelay",0,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodePerfMonitorEventLoopDelay(loop,args)}))
		module.set("createHistogram",tsNamedNativeMethod("createHistogram",0,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodePerfHistogram(loop,args)}))
		tsNodePerfSingleton.value=module
		global:=tsGlobalObject();if tsIsUndefined(tsGet(global,tsStringReference(tsStringUTF8("performance")))){tsSet(loop,global,tsStringReference(tsStringUTF8("performance")),tsObjectValue(performance))}
	})
	return tsObjectValue(tsNodePerfSingleton.value)
}
`
