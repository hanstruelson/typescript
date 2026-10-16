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

type tsNodeHistogram struct { counts map[int64]uint64; previous uint64; exceeds uint64; total uint64; lowest,highest int64; figures int; manualMin,manualMax int64 }
var tsNodeHistogramStates sync.Map
var tsNodeHistogramPrototypeState struct { sync.Once; value *tsObject }
var tsNodePerfDelay struct { sync.Mutex; histogram tsValue; state *tsNodeHistogram; loop *tsLoop; task *tsTask; resolution int; enabled bool }

func tsNodeHistogramState(value tsValue) *tsNodeHistogram {
	if value.kind!=tsObjectKind { return nil }
	state,_:=tsNodeHistogramStates.Load(value.ref)
	if state==nil { return nil }
	return state.(*tsNodeHistogram)
}

func tsNodeHistogramSigned(value tsValue) int64 {
	if value.kind==tsBigIntKind { return (*big.Int)(value.ref).Int64() }
	n:=tsNumber(value)
	if math.IsNaN(n) { return 0 }
	if n>=9223372036854775808.0 { return 1<<63-1 }
	if n<=-9223372036854775808.0 { return -1<<63 }
	return int64(math.Trunc(n))
}

func tsNodeHistogramRecord(state *tsNodeHistogram,value tsValue) {
	if value.kind!=tsBigIntKind&&!tsIsNumeric(value) { tsInternalValidatorInvalidType("value","a number or bigint",value) }
	number:=tsNodeHistogramSigned(value)
	if number<1 { tsInternalValidatorOutOfRange("value",">= 1",value) }
	if number<state.lowest||number>state.highest { state.exceeds++; return }
	if state.counts==nil{state.counts=map[int64]uint64{}}
	state.counts[tsNodeHistogramHighestEquivalent(number,state.figures)]++
	state.total++
	if number<state.manualMin { state.manualMin=number }
	if number>state.manualMax { state.manualMax=number }
}

func tsNodeHistogramHighestEquivalent(value int64,figures int) int64 {
	if value<=0 { return value }
	unit:=math.Pow10(int(math.Floor(math.Log10(float64(value))))-figures+1)
	if unit<1 { return value }
	return int64(math.Floor(float64(value)/unit)*unit+unit-1)
}

func tsNodeHistogramPercentile(state *tsNodeHistogram,p float64) int64 {
	if len(state.counts)==0 { return 0 }
	values:=make([]int64,0,len(state.counts));for value:=range state.counts{values=append(values,value)}
	sort.Slice(values,func(i,j int)bool{return values[i]<values[j]})
	count:=uint64(0);for _,quantity:=range state.counts{count+=quantity}
	target:=uint64(math.Ceil(p/100*float64(count)));if target<1{target=1};var seen uint64
	for _,value:=range values{seen+=state.counts[value];if seen>=target{return value}}
	return values[len(values)-1]
}

func tsNodeHistogramStats(state *tsNodeHistogram) (int64,int64,float64,float64) {
	if state.total==0 { return 0,0,0,0 }
	mean:=float64(0);count:=uint64(0)
	for value,quantity:=range state.counts { mean+=float64(value)*float64(quantity);count+=quantity }
	if count>0 { mean/=float64(count) }
	variance:=float64(0)
	for value,quantity:=range state.counts { delta:=float64(value)-mean;variance+=delta*delta*float64(quantity) }
	if count>0 { variance/=float64(count) }
	return state.manualMin,state.manualMax,mean,math.Sqrt(variance)
}

func tsNodeHistogramPrototype() *tsObject {
	tsNodeHistogramPrototypeState.Do(func(){
		prototype:=tsNewObject();prototype.prototype=tsObjectValue(tsObjectPrototype());prototype.descriptors=map[string]*tsDescriptor{}
		method:=func(name string,length int,call func(*tsLoop,tsValue,[]tsValue)tsValue){prototype.set(name,tsNamedNativeMethod(name,length,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return call(loop,receiver,args)}));prototype.descriptors[name]=&tsDescriptor{writable:true,configurable:true}}
		method("record",1,func(_ *tsLoop,receiver tsValue,args []tsValue)tsValue{state:=tsNodeHistogramState(receiver);if state==nil{tsPropertyFailure("Histogram method called on incompatible receiver")};if len(args)==0{tsNodeAssertMissing("value")};tsNodeHistogramRecord(state,args[0]);return tsU})
		method("recordDelta",0,func(_ *tsLoop,receiver tsValue,_ []tsValue)tsValue{state:=tsNodeHistogramState(receiver);if state==nil{tsPropertyFailure("Histogram method called on incompatible receiver")};now:=uint64(time.Since(tsNodePerfSingleton.started).Nanoseconds());if state.previous!=0&&now>=state.previous{tsNodeHistogramRecord(state,tsNumberValue(float64(now-state.previous)))};state.previous=now;return tsU})
		method("add",1,func(_ *tsLoop,receiver tsValue,args []tsValue)tsValue{state:=tsNodeHistogramState(receiver);if state==nil{tsPropertyFailure("Histogram method called on incompatible receiver")};if len(args)==0{tsNodeAssertMissing("other")};other:=tsNodeHistogramState(args[0]);if other==nil{tsInternalValidatorInvalidType("argument","a RecordableHistogram",args[0])};wasEmpty:=state.total==0;state.total+=other.total;state.exceeds+=other.exceeds;if other.total>0{if wasEmpty{state.manualMin,state.manualMax=other.manualMin,other.manualMax}else{if other.manualMin<state.manualMin{state.manualMin=other.manualMin};if other.manualMax>state.manualMax{state.manualMax=other.manualMax}}};if state.counts==nil{state.counts=map[int64]uint64{}};dropped:=uint64(0);for value,quantity:=range other.counts{if value<state.lowest||value>state.highest{dropped+=quantity;continue};state.counts[tsNodeHistogramHighestEquivalent(value,state.figures)]+=quantity};return tsNumberValue(float64(dropped))})
		method("reset",0,func(_ *tsLoop,receiver tsValue,_ []tsValue)tsValue{state:=tsNodeHistogramState(receiver);if state==nil{tsPropertyFailure("Histogram method called on incompatible receiver")};state.counts=nil;state.previous=0;state.exceeds=0;state.total=0;state.manualMin=1<<63-1;state.manualMax=0;return tsU})
		for _,name:=range []string{"percentile","percentileBigInt"}{key:=name;method(key,1,func(_ *tsLoop,receiver tsValue,args []tsValue)tsValue{state:=tsNodeHistogramState(receiver);if state==nil{tsPropertyFailure("Histogram method called on incompatible receiver")};if len(args)==0{tsNodeAssertMissing("percentile")};value:=args[0];tsInternalValidatorNumber(value,"percentile",tsNumberValue(0),tsNumberValue(100),false);p:=tsNumber(value);if p<=0||math.IsNaN(p){tsInternalValidatorOutOfRange("percentile","> 0 && <= 100",value)};percentile:=tsNodeHistogramPercentile(state,p);if key=="percentileBigInt"{return tsBigIntSigned(percentile)};return tsNumberValue(float64(percentile))})}
		for _,name:=range []string{"count","countBigInt","min","minBigInt","max","maxBigInt","mean","stddev","exceeds","exceedsBigInt","percentiles","percentilesBigInt"}{key:=name;prototype.set(key,tsU);prototype.descriptors[key]=&tsDescriptor{accessor:true,configurable:true,getter:tsNamedNativeMethod("get "+key,0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeHistogramState(receiver);if state==nil{tsPropertyFailure("Histogram getter called on incompatible receiver")};min,max,mean,stddev:=tsNodeHistogramStats(state);switch key{case "count":return tsNumberValue(float64(state.total));case "countBigInt":return tsBigIntUnsigned(state.total);case "min":if state.total==0{return tsNumberValue(float64(uint64(1)<<63))};return tsNumberValue(float64(min));case "minBigInt":if state.total==0{return tsBigIntSigned(1<<63-1)};return tsBigIntSigned(min);case "max":return tsNumberValue(float64(max));case "maxBigInt":return tsBigIntSigned(max);case "mean":return tsNumberValue(mean);case "stddev":return tsNumberValue(stddev);case "exceeds":return tsNumberValue(float64(state.exceeds));case "exceedsBigInt":return tsBigIntUnsigned(state.exceeds);default:percentiles:=tsNewCollection(false);if len(state.counts)>0{for _,p:=range []float64{0,50,75,87.5,93.75,96.875,98.4375,99.21875,99.609375,99.8046875,100}{value:=tsNodeHistogramPercentile(state,p);result:=tsNumberValue(float64(value));if key=="percentilesBigInt"{result=tsBigIntSigned(value)};percentiles.put(tsNumberValue(p),result)}};return tsCollectionValue(percentiles)}}) }}
		tag:=tsPropertyKey(tsWellKnownSymbol("toStringTag"));prototype.set(tag,tsStringReference(tsStringUTF8("RecordableHistogram")));prototype.descriptors[tag]=&tsDescriptor{configurable:true};tsNodeHistogramPrototypeState.value=prototype
	})
	return tsNodeHistogramPrototypeState.value
}

func tsNodeHistogramValue(histogram *tsNodeHistogram) tsValue {
	if histogram.lowest==0{histogram.lowest=1};if histogram.highest==0{histogram.highest=9007199254740991};if histogram.figures==0{histogram.figures=3};if histogram.manualMin==0{histogram.manualMin=1<<63-1}
	object:=tsNewObject();object.nativeClass=tsNativeClass("RecordableHistogram");object.prototype=tsObjectValue(tsNodeHistogramPrototype());value:=tsObjectValue(object);tsNodeHistogramStates.Store(value.ref,histogram);return value
}

func tsNodePerfDelaySchedule(loop *tsLoop,histogram tsValue) {
	tsNodePerfDelay.Lock();if !tsNodePerfDelay.enabled||tsNodePerfDelay.loop!=loop{tsNodePerfDelay.Unlock();return};resolution:=tsNodePerfDelay.resolution;tsNodePerfDelay.Unlock()
	period:=time.Duration(resolution)*time.Millisecond;if period<=0{period=time.Millisecond};expected:=time.Now().Add(period)
	callback:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{
		delay:=time.Since(expected).Nanoseconds();if delay>0{if state:=tsNodeHistogramState(histogram);state!=nil{tsNodeHistogramRecord(state,tsNumberValue(float64(delay)))}}
		tsNodePerfDelay.Lock();active:=tsNodePerfDelay.enabled&&tsNodePerfDelay.loop==loop;if active{tsNodePerfDelay.task=nil};tsNodePerfDelay.Unlock();if active{tsNodePerfDelaySchedule(loop,histogram)};return tsU
	}))
	task:=loop.setTimeout(callback,tsNumberValue(float64(resolution)));task.unref=true
	tsNodePerfDelay.Lock();if tsNodePerfDelay.enabled&&tsNodePerfDelay.loop==loop{tsNodePerfDelay.task=task;tsNodePerfDelay.Unlock()}else{tsNodePerfDelay.Unlock();loop.clearTimeout(tsValue{ref:unsafe.Pointer(task),kind:tsTaskKind})}
}

func tsNodePerfHistogram(loop *tsLoop,args []tsValue) tsValue {
	options:=tsArg(args,0)
	if !tsIsUndefined(options){tsInternalValidatorObject(options,tsStringReference(tsStringUTF8("options")),tsU)}
	lowestValue,highestValue,figuresValue:=tsFSProperty(options,"lowest"),tsFSProperty(options,"highest"),tsFSProperty(options,"figures")
	lowest,highest,figures:=float64(1),float64(9007199254740991),float64(3)
	if !tsIsUndefined(lowestValue){if lowestValue.kind==tsBigIntKind{lowest,_=new(big.Float).SetInt((*big.Int)(lowestValue.ref)).Float64()}else if tsIsNumeric(lowestValue){lowest=tsNumber(lowestValue)}else{tsInternalValidatorInvalidType("options.lowest","a number or a bigint",lowestValue)}}
	if !tsIsUndefined(highestValue){if highestValue.kind==tsBigIntKind{highest,_=new(big.Float).SetInt((*big.Int)(highestValue.ref)).Float64()}else if tsIsNumeric(highestValue){highest=tsNumber(highestValue)}else{tsInternalValidatorInvalidType("options.highest","a number or a bigint",highestValue)}}
	if !tsIsUndefined(figuresValue){if !tsIsNumeric(figuresValue){tsInternalValidatorInvalidType("options.figures","a number",figuresValue)};figures=tsNumber(figuresValue);if figures<1||figures>5{tsInternalValidatorOutOfRange("options.figures",">= 1 && <= 5",figuresValue)}}
	if math.IsNaN(lowest){lowest=1};if math.IsNaN(highest){highest=float64(1<<63-1)};if math.IsNaN(figures){figures=3}
	if lowest<1{tsInternalValidatorOutOfRange("options.lowest",">= 1 && <= 9007199254740991",tsNumberValue(lowest))}
	if highest<2*lowest{tsInternalValidatorOutOfRange("options.highest",">= "+tsNumberText(2*lowest)+" && <= 9007199254740991",tsNumberValue(highest))}
	return tsNodeHistogramValue(&tsNodeHistogram{lowest:tsNodeHistogramSigned(tsNumberValue(lowest)),highest:tsNodeHistogramSigned(tsNumberValue(highest)),figures:int(math.Trunc(figures)),manualMin:1<<63-1})
}

func tsNodePerfMonitorEventLoopDelay(loop *tsLoop,args []tsValue) tsValue {
	options:=tsArg(args,0);resolution:=10
	if !tsIsUndefined(options){tsInternalValidatorObject(options,tsStringReference(tsStringUTF8("options")),tsU);value:=tsFSProperty(options,"resolution");if !tsIsUndefined(value){tsInternalValidatorNumber(value,"options.resolution",tsNumberValue(1),tsNumberValue(2147483647),true);resolution=int(tsNumber(value))}}
	tsNodePerfDelay.Lock();defer tsNodePerfDelay.Unlock();tsNodePerfDelay.resolution=resolution
	if tsIsUndefined(tsNodePerfDelay.histogram){tsNodePerfDelay.state=&tsNodeHistogram{lowest:1,highest:3600000000000,figures:3,manualMin:1<<63-1};tsNodePerfDelay.histogram=tsNodeHistogramValue(tsNodePerfDelay.state);object:=(*tsObject)(tsNodePerfDelay.histogram.ref)
		object.set("enable",tsNamedNativeMethod("enable",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{tsNodePerfDelay.Lock();if tsNodePerfDelay.enabled{tsNodePerfDelay.Unlock();return tsBooleanValue(false)};tsNodePerfDelay.enabled=true;tsNodePerfDelay.loop=loop;histogram:=tsNodePerfDelay.histogram;state:=tsNodeHistogramState(histogram);state.counts=nil;state.previous=0;state.exceeds=0;state.total=0;state.manualMin=1<<63-1;state.manualMax=0;tsNodePerfDelay.Unlock();tsNodePerfDelaySchedule(loop,histogram);return tsBooleanValue(true)}))
		object.set("disable",tsNamedNativeMethod("disable",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{tsNodePerfDelay.Lock();if !tsNodePerfDelay.enabled{tsNodePerfDelay.Unlock();return tsBooleanValue(false)};tsNodePerfDelay.enabled=false;task:=tsNodePerfDelay.task;tsNodePerfDelay.task=nil;owner:=tsNodePerfDelay.loop;tsNodePerfDelay.Unlock();if task!=nil&&owner!=nil{owner.clearTimeout(tsValue{ref:unsafe.Pointer(task),kind:tsTaskKind})};return tsBooleanValue(true)}))
		dispose:=tsWellKnownSymbol("dispose");if !tsIsUndefined(dispose){object.set(tsPropertyKey(dispose),object.get("disable"))}
	};return tsNodePerfDelay.histogram
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
