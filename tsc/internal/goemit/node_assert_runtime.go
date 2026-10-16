package goemit

// nodeAssertRuntime provides the Node assertion entry point using the
// transpiler's own value, error, promise, and RegExp machinery.
const nodeAssertRuntime = `
type tsNodeAssertState struct { sync.Once; module, strict tsValue }
var tsNodeAssertSingleton tsNodeAssertState

type tsNodeAssertPair struct { left, right unsafe.Pointer; kind tsKind }

func tsNodeAssertError(loop *tsLoop,message tsValue, actual, expected tsValue, operator string) tsValue {
	if message.kind==tsErrorKind || message.kind==tsObjectKind && (*tsObject)(message.ref).nativeError { return message }
	options:=tsNewObject();options.set("message",message);options.set("actual",actual);options.set("expected",expected);options.set("operator",tsStringReference(tsStringUTF8(operator)));options.set("diff",tsStringReference(tsStringUTF8("simple")))
	return tsAssertErrorNew(loop,tsObjectValue(options))
}

func tsNodeAssertFail(loop *tsLoop,message tsValue, actual, expected tsValue, operator string) { panic(tsThrown{tsNodeAssertError(loop,message,actual,expected,operator)}) }

func tsNodeAssertMissing(name string) { panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"TypeError",message:"The "+name+" argument must be specified"})}) }

func tsNodeAssertDeep(loop *tsLoop,left,right tsValue,strict bool,seen map[tsNodeAssertPair]bool) bool {
	if tsIsPrimitiveValue(left)||tsIsPrimitiveValue(right) {
		if strict{return tsSameValue(left,right)}
		return tsLooseEqual(left,right)||(tsIsNumeric(left)&&tsIsNumeric(right)&&math.IsNaN(tsNumber(left))&&math.IsNaN(tsNumber(right)))
	}
	if left.kind!=right.kind{return false};if left.ref==right.ref{return true}
	if left.kind==tsArrayKind {
		a,b:=(*tsArray)(left.ref),(*tsArray)(right.ref);if a.length()!=b.length(){return false};pair:=tsNodeAssertPair{left.ref,right.ref,left.kind};if seen[pair]{return true};seen[pair]=true
		for i:=0;i<a.length();i++{if !tsNodeAssertDeep(loop,a.at(i),b.at(i),strict,seen){return false}};return true
	}
	if left.kind!=tsObjectKind && left.kind!=tsInstanceKind && left.kind!=tsErrorKind{return false}
	pair:=tsNodeAssertPair{left.ref,right.ref,left.kind};if seen[pair]{return true};seen[pair]=true
	keysLeft:=tsObjectKeys(left);keysRight:=tsObjectKeys(right);if keysLeft.length()!=keysRight.length(){return false}
	for i:=0;i<keysLeft.length();i++{key:=keysLeft.at(i);if !tsHas(right,key){return false};if !tsNodeAssertDeep(loop,tsGet(loop,left,key),tsGet(loop,right,key),strict,seen){return false}}
	return true
}

func tsNodeAssertErrorValue(value tsValue) bool {
	switch value.kind {
	case tsErrorKind:
		return true
	case tsObjectKind:
		return (*tsObject)(value.ref).nativeError
	case tsInstanceKind:
		return (*tsProperties)(value.ref).nativeError
	default:
		return false
	}
}

func tsNodeAssertArrayIndex(name string) bool {
	index, err := strconv.ParseUint(name, 10, 32)
	return err == nil && index < 0xffffffff && strconv.FormatUint(index, 10) == name
}

func tsNodeAssertPartialProperties(loop *tsLoop, actual, expected tsValue, nonIndex bool, actualSeen, expectedSeen map[unsafe.Pointer]bool) bool {
	actualKeys := map[string]bool{}
	for _, key := range tsObjectKeys(actual).values {
		name := tsPropertyKey(key)
		if nonIndex && tsNodeAssertArrayIndex(name) { continue }
		actualKeys[name] = true
	}
	for _, key := range tsObjectKeys(expected).values {
		name := tsPropertyKey(key)
		if nonIndex && tsNodeAssertArrayIndex(name) { continue }
		if !actualKeys[name] || !tsNodeAssertPartialDeep(loop, tsGet(actual, key), tsGet(expected, key), actualSeen, expectedSeen) { return false }
	}
	return true
}

func tsNodeAssertPartialBytes(actual, expected []byte) bool {
	if len(expected) > len(actual) { return false }
	pos := 0
	for i, item := range expected {
		last := len(actual) - len(expected) + i
		for pos <= last && actual[pos] != item { pos++ }
		if pos > last { return false }
		pos++
	}
	return true
}

func tsNodeAssertPartialTypedArray(loop *tsLoop, actual, expected tsValue, actualSeen, expectedSeen map[unsafe.Pointer]bool) bool {
	a, b := (*tsTypedArray)(actual.ref), (*tsTypedArray)(expected.ref)
	actualName, expectedName := a.name, b.name
	if actualName == "Buffer" { actualName = "Uint8Array" }
	if expectedName == "Buffer" { expectedName = "Uint8Array" }
	if actualName != expectedName || b.length > a.length || a.buffer.detached || b.buffer.detached { return false }
	pos := 0
	for i := 0; i < b.length; i++ {
		last := a.length - b.length + i
		expectedValue := tsNumberValue(b.read(i)); if b.readValue != nil { expectedValue = b.readValue(i) }
		matched := false
		for pos <= last {
			actualValue := tsNumberValue(a.read(pos)); if a.readValue != nil { actualValue = a.readValue(pos) }
			pos++
			if tsSameValue(actualValue, expectedValue) { matched = true; break }
		}
		if !matched { return false }
	}
	return tsNodeAssertPartialProperties(loop, actual, expected, true, actualSeen, expectedSeen)
}

func tsNodeAssertPartialCollection(loop *tsLoop, actual, expected tsValue, actualSeen, expectedSeen map[unsafe.Pointer]bool) bool {
	a, b := (*tsCollection)(actual.ref), (*tsCollection)(expected.ref)
	if a.set != b.set { return false }
	actualCount, expectedCount := 0, 0
	for _, entry := range a.entries { if entry.live { actualCount++ } }
	for _, entry := range b.entries { if entry.live { expectedCount++ } }
	if expectedCount > actualCount { return false }
	used, identityUsed := make([]bool, len(a.entries)), make([]bool, len(a.entries))
	for _, wanted := range b.entries {
		if !wanted.live { continue }
		match := -1
		if !b.set {
			for i, candidate := range a.entries {
				if candidate.live && tsSameValueZero(candidate.key, wanted.key) && !used[i] {
					if tsNodeAssertPartialDeep(loop, candidate.value, wanted.value, actualSeen, expectedSeen) { match = i; identityUsed[i] = true; break }
				}
			}
		}
		if match < 0 && !b.set && tsIsPrimitiveValue(wanted.key) { return false }
		if match < 0 {
			for i, candidate := range a.entries {
				if !candidate.live || used[i] || identityUsed[i] { continue }
				if b.set {
					if tsNodeAssertPartialDeep(loop, candidate.key, wanted.key, actualSeen, expectedSeen) { match = i; break }
				} else if tsNodeAssertPartialDeep(loop, candidate.key, wanted.key, actualSeen, expectedSeen) && tsNodeAssertPartialDeep(loop, candidate.value, wanted.value, actualSeen, expectedSeen) { match = i; break }
			}
		}
		if match < 0 { return false }
		used[match] = true
	}
	return tsNodeAssertPartialProperties(loop, actual, expected, false, actualSeen, expectedSeen)
}

func tsNodeAssertPartialDeep(loop *tsLoop, actual, expected tsValue, actualSeen, expectedSeen map[unsafe.Pointer]bool) bool {
	if tsIsPrimitiveValue(actual) || tsIsPrimitiveValue(expected) {
		return tsSameValue(actual, expected)
	}
	if actual.kind == tsFunctionKind || actual.kind == tsClassKind || expected.kind == tsFunctionKind || expected.kind == tsClassKind {
		return tsNodeAssertDeep(loop, actual, expected, true, map[tsNodeAssertPair]bool{})
	}
	if actual.kind == tsPromiseKind || expected.kind == tsPromiseKind { return false }
	if actual.kind == expected.kind && actual.ref == expected.ref { return true }
	if actual.ref == nil || expected.ref == nil { return false }
	if actualSeen[actual.ref] && expectedSeen[expected.ref] { return true }
	if actualSeen[actual.ref] || expectedSeen[expected.ref] { return false }
	actualSeen[actual.ref], expectedSeen[expected.ref] = true, true
	defer func() { delete(actualSeen, actual.ref); delete(expectedSeen, expected.ref) }()

	if actual.kind == tsCollectionKind || expected.kind == tsCollectionKind {
		if actual.kind != tsCollectionKind || expected.kind != tsCollectionKind { return false }
		return tsNodeAssertPartialCollection(loop, actual, expected, actualSeen, expectedSeen)
	}
	if actual.kind == tsTypedArrayKind || actual.kind == tsDataViewKind || expected.kind == tsTypedArrayKind || expected.kind == tsDataViewKind {
		if actual.kind != expected.kind { return false }
		if actual.kind == tsTypedArrayKind { return tsNodeAssertPartialTypedArray(loop, actual, expected, actualSeen, expectedSeen) }
		a, b := (*tsDataView)(actual.ref), (*tsDataView)(expected.ref)
		if a.buffer.detached || b.buffer.detached || b.length > a.length || !tsNodeAssertPartialBytes(a.buffer.data[a.offset:a.offset+a.length], b.buffer.data[b.offset:b.offset+b.length]) { return false }
		return tsNodeAssertPartialProperties(loop, actual, expected, false, actualSeen, expectedSeen)
	}
	if actual.kind == tsArrayBufferKind || expected.kind == tsArrayBufferKind {
		if actual.kind != expected.kind { return false }
		a, b := (*tsArrayBuffer)(actual.ref), (*tsArrayBuffer)(expected.ref)
		if a.detached || b.detached || !tsNodeAssertPartialBytes(a.data, b.data) { return false }
		return tsNodeAssertPartialProperties(loop, actual, expected, false, actualSeen, expectedSeen)
	}
	if tsNodeAssertErrorValue(actual) || tsNodeAssertErrorValue(expected) {
		if !tsNodeAssertErrorValue(actual) || !tsNodeAssertErrorValue(expected) { return false }
		for _, name := range []string{"message", "name", "errors"} {
			want := tsGet(expected, tsStringReference(tsStringUTF8(name)))
			if tsIsUndefined(want) || name == "message" && want.kind == tsStringKind && len((*tsString)(want.ref).units) == 0 { continue }
			if !tsNodeAssertPartialDeep(loop, tsGet(actual, tsStringReference(tsStringUTF8(name))), want, actualSeen, expectedSeen) { return false }
		}
		cause := tsStringReference(tsStringUTF8("cause"))
		if tsObjectHasOwn(expected, "cause") && (!tsObjectHasOwn(actual, "cause") || !tsNodeAssertPartialDeep(loop, tsGet(actual, cause), tsGet(expected, cause), actualSeen, expectedSeen)) { return false }
		return tsNodeAssertPartialProperties(loop, actual, expected, false, actualSeen, expectedSeen)
	}
	if actual.kind == tsObjectKind && expected.kind == tsObjectKind {
		a, b := (*tsObject)(actual.ref), (*tsObject)(expected.ref)
		if a.nativeDate || b.nativeDate {
			if !a.nativeDate || !b.nativeDate || !(a.nativeTime == b.nativeTime || math.IsNaN(a.nativeTime) && math.IsNaN(b.nativeTime)) { return false }
			return tsNodeAssertPartialProperties(loop, actual, expected, false, actualSeen, expectedSeen)
		}
		if a.nodeURL != nil || b.nodeURL != nil {
			if a.nodeURL == nil || b.nodeURL == nil || tsText(tsNodeURLField(a.nodeURL, "href")) != tsText(tsNodeURLField(b.nodeURL, "href")) { return false }
			return tsNodeAssertPartialProperties(loop, actual, expected, false, actualSeen, expectedSeen)
		}
		if a.boxedPrimitive.kind != tsUndefinedKind || b.boxedPrimitive.kind != tsUndefinedKind {
			if a.boxedPrimitive.kind == tsUndefinedKind || b.boxedPrimitive.kind == tsUndefinedKind || !tsSameValue(a.boxedPrimitive, b.boxedPrimitive) { return false }
			return tsNodeAssertPartialProperties(loop, actual, expected, false, actualSeen, expectedSeen)
		}
	}
	if actual.kind == tsRegExpKind || expected.kind == tsRegExpKind {
		if actual.kind != tsRegExpKind || expected.kind != tsRegExpKind { return false }
		a, b := (*tsRegExp)(actual.ref), (*tsRegExp)(expected.ref)
		if tsStringEqual(a.source, b.source) == false || a.flags != b.flags || !tsSameValue(tsRegExpGetLastIndex(loop, actual), tsRegExpGetLastIndex(loop, expected)) { return false }
		return tsNodeAssertPartialProperties(loop, actual, expected, false, actualSeen, expectedSeen)
	}
	if actual.kind == tsArrayKind || expected.kind == tsArrayKind {
		if actual.kind != tsArrayKind || expected.kind != tsArrayKind { return false }
		a, b := (*tsArray)(actual.ref), (*tsArray)(expected.ref)
		if b.length() > a.length() { return false }
		pos := 0
		for i := 0; i < b.length(); i++ {
			last := a.length() - b.length() + i
			matched := false
			for pos <= last {
				if tsNodeAssertPartialDeep(loop, a.at(pos), b.at(i), actualSeen, expectedSeen) { matched = true; pos++; break }
				pos++
			}
			if !matched { return false }
		}
		return tsNodeAssertPartialProperties(loop, actual, expected, true, actualSeen, expectedSeen)
	}
	if actual.kind != tsObjectKind && actual.kind != tsInstanceKind || expected.kind != tsObjectKind && expected.kind != tsInstanceKind { return false }
	if actual.kind == tsObjectKind && expected.kind == tsObjectKind {
		if (tsNodeCryptoKeyFor(actual) != nil || tsNodeCryptoKeyFor(expected) != nil) && (tsNodeCryptoKeyFor(actual) == nil || tsNodeCryptoKeyFor(expected) == nil || !tsNodeCryptoKeysEqual(tsNodeCryptoKeyFor(actual), tsNodeCryptoKeyFor(expected))) { return false }
	}
	return tsNodeAssertPartialProperties(loop, actual, expected, false, actualSeen, expectedSeen)
}

func tsNodeAssertMatch(loop *tsLoop,actual,expected tsValue) bool {
	if expected.kind==tsRegExpKind {
		input:=tsStringValueWithLoop(loop,actual);return tsRegExpExecReceiver(loop,expected,input).kind!=tsNullKind
	}
	if expected.kind==tsFunctionKind {
		if tsInstanceOf(actual,expected){return true}
		return tsTruthy(tsCall(loop,expected,actual))
	}
	if expected.kind==tsObjectKind || expected.kind==tsInstanceKind {
		keys:=tsObjectKeys(expected);for i:=0;i<keys.length();i++{key:=keys.at(i);if !tsHas(actual,key)||!tsNodeAssertDeep(loop,tsGet(loop,actual,key),tsGet(loop,expected,key),true,map[tsNodeAssertPair]bool{}){return false}};return true
	}
	return false
}

func tsNodeAssertThrown(loop *tsLoop,callback tsValue) (thrown tsValue, didThrow bool) {
	defer func(){if value:=recover();value!=nil{didThrow=true;thrown=tsUnwrap(value)}}()
	if callback.kind!=tsFunctionKind{tsPropertyFailure("The \"fn\" argument must be of type function")};tsCall(loop,callback);return tsU,false
}

func tsNodeAssertAsync(loop *tsLoop,value tsValue,expected,message tsValue,shouldReject bool) tsValue {
	var result tsValue
	if value.kind==tsFunctionKind {var thrown tsValue;var failed bool;func(){defer func(){if item:=recover();item!=nil{thrown=tsUnwrap(item);failed=true}}();result=tsCall(loop,value)}();if failed{result=tsPromiseValue(loop.resolved(thrown,true))}} else {result=value}
	if result.kind!=tsPromiseKind {result=tsPromiseValue(loop.resolved(result,false))}
	promise:=loop.promise();loop.await(result,func(settled tsResult){
		matched:=settled.rejected && (!shouldReject || tsNodeAssertMatch(loop,settled.value,expected))
		if !settled.rejected && !shouldReject {matched=true}
		if matched {promise.settle(tsResult{tsU,false});return}
		if !settled.rejected && shouldReject {promise.settle(tsResult{tsNodeAssertError(loop,message,settled.value,expected,"rejects"),true});return}
		if settled.rejected && !shouldReject {promise.settle(tsResult{tsNodeAssertError(loop,message,settled.value,tsU,"doesNotReject"),true});return}
		promise.settle(tsResult{tsNodeAssertError(loop,message,settled.value,expected,"rejects"),true})
	})
	return tsPromiseValue(promise)
}

func tsNodeAssertMethod(name string) tsValue {
	fn:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{
		actual,expected,message:=tsArg(args,0),tsArg(args,1),tsArg(args,2)
		switch name {
		case "ok": if !tsTruthy(actual){tsNodeAssertFail(loop,expected,actual,tsBooleanValue(true),"==")}
		case "fail": tsNodeAssertFail(loop,actual,tsU,tsU,"fail")
		case "equal","notEqual": if len(args)<2{tsNodeAssertMissing("expected")};same:=tsLooseEqual(actual,expected)||(tsIsNumeric(actual)&&tsIsNumeric(expected)&&math.IsNaN(tsNumber(actual))&&math.IsNaN(tsNumber(expected)));if (name=="equal"&&!same)||(name=="notEqual"&&same){operator:="==";if name=="notEqual"{operator="!="};tsNodeAssertFail(loop,message,actual,expected,operator)}
		case "strictEqual","notStrictEqual": if len(args)<2{tsNodeAssertMissing("expected")};same:=tsSameValue(actual,expected);if (name=="strictEqual"&&!same)||(name=="notStrictEqual"&&same){operator:="strictEqual";if name=="notStrictEqual"{operator="notStrictEqual"};tsNodeAssertFail(loop,message,actual,expected,operator)}
		case "deepEqual","notDeepEqual","deepStrictEqual","notDeepStrictEqual","partialDeepStrictEqual":
			if len(args)<2{tsNodeAssertMissing("expected")};strict:=name=="deepStrictEqual"||name=="notDeepStrictEqual"||name=="partialDeepStrictEqual";same:=false;if name=="partialDeepStrictEqual"{same=tsNodeAssertPartialDeep(loop,actual,expected,map[unsafe.Pointer]bool{},map[unsafe.Pointer]bool{})}else{same=tsNodeAssertDeep(loop,actual,expected,strict,map[tsNodeAssertPair]bool{})};if name=="notDeepEqual"||name=="notDeepStrictEqual"{same=!same};if !same{tsNodeAssertFail(loop,message,actual,expected,name)}
		case "throws","doesNotThrow":
			if len(args)==0{tsNodeAssertMissing("fn")};callback:=actual;want:=expected;msg:=message;if len(args)==2&&want.kind==tsStringKind{msg=want;want=tsU};thrown,didThrow:=tsNodeAssertThrown(loop,callback)
			if name=="throws" {if !didThrow||!tsIsUndefined(want)&&!tsNodeAssertMatch(loop,thrown,want){tsNodeAssertFail(loop,msg,thrown,want,"throws")}} else if didThrow {if tsIsUndefined(want)||tsNodeAssertMatch(loop,thrown,want){tsNodeAssertFail(loop,msg,thrown,tsU,"doesNotThrow")};panic(tsThrown{thrown})}
		case "rejects","doesNotReject": if len(args)==2&&expected.kind==tsStringKind{return tsNodeAssertAsync(loop,actual,tsU,expected,name=="rejects")};return tsNodeAssertAsync(loop,actual,expected,message,name=="rejects")
		case "ifError": if tsTruthy(actual){tsNodeAssertFail(loop,message,actual,tsNull,"ifError")}
		case "match","doesNotMatch":
			if actual.kind!=tsStringKind||expected.kind!=tsRegExpKind {tsPropertyFailure("The string and regexp arguments are required")};matched:=tsRegExpExecReceiver(loop,expected,(*tsString)(actual.ref)).kind!=tsNullKind;if (name=="match"&&!matched)||(name=="doesNotMatch"&&matched){tsNodeAssertFail(loop,message,actual,expected,name)}
		}
		return tsU
	})
	fn.name=name;fn.length=2
	return tsFunctionValue(fn)
}

func tsNodeAssertDecorate(loop *tsLoop,value tsValue, strict bool) tsValue {
	fn:=(*tsFunction)(value.ref);fn.properties=tsNewObject()
	for _,name:=range []string{"ok","fail","equal","notEqual","deepEqual","notDeepEqual","deepStrictEqual","notDeepStrictEqual","strictEqual","notStrictEqual","partialDeepStrictEqual","throws","doesNotThrow","rejects","doesNotReject","ifError","match","doesNotMatch"}{method:=name;if strict{switch name{case "equal":method="strictEqual";case "notEqual":method="notStrictEqual";case "deepEqual":method="deepStrictEqual";case "notDeepEqual":method="notDeepStrictEqual"}};fn.properties.set(name,tsNodeAssertMethod(method))}
	fn.properties.set("AssertionError",tsInternalAssertionErrorConstructor(loop))
	fn.properties.set("CallTracker",tsNodeCallTrackerConstructor(loop))
	fn.properties.set("Assert",tsNodeAssertConstructor())
	return value
}

func tsNodeAssertConstructor() tsValue {
	fn:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{tsPropertyFailure("Class constructor Assert cannot be invoked without new")});fn.name="Assert";fn.length=1;fn.constructible=true
	fn.construct=func(loop *tsLoop,args ...tsValue)tsValue{instance:=tsNewObject();options:=tsArg(args,0);strictValue:=tsFSProperty(options,"strict");strict:=tsIsUndefined(strictValue)||tsTruthy(strictValue);for _,name:=range []string{"ok","fail","equal","notEqual","deepEqual","notDeepEqual","deepStrictEqual","notDeepStrictEqual","strictEqual","notStrictEqual","partialDeepStrictEqual","throws","doesNotThrow","rejects","doesNotReject","ifError","match","doesNotMatch"}{method:=name;if strict{switch name{case "equal":method="strictEqual";case "notEqual":method="notStrictEqual";case "deepEqual":method="deepStrictEqual";case "notDeepEqual":method="notDeepStrictEqual"}};instance.set(name,tsNodeAssertMethod(method))};return tsObjectValue(instance)}
	return tsFunctionValue(fn)
}

type tsNodeCallTrackerContext struct{expected tsValue;calls []tsValue;name string;stack tsValue;tracked tsValue}
type tsNodeCallTracker struct{object *tsObject;checks []*tsNodeCallTrackerContext;tracked map[unsafe.Pointer]*tsNodeCallTrackerContext}
var tsNodeCallTrackerStates struct{sync.Once;instances map[*tsObject]*tsNodeCallTracker;prototype *tsObject;constructor tsValue}
func tsNodeCallTrackerState(receiver tsValue)*tsNodeCallTracker{if receiver.kind!=tsObjectKind{return nil};return tsNodeCallTrackerStates.instances[(*tsObject)(receiver.ref)]}
func tsNodeCallTrackerFreezeArray(values []tsValue)tsValue{array:=&tsArray{values:append([]tsValue(nil),values...),descriptors:map[string]*tsDescriptor{}};for index:=range array.values{array.descriptors[strconv.Itoa(index)]=&tsDescriptor{writable:false,configurable:false}};array.descriptors["length"]=&tsDescriptor{writable:false};return tsArrayValue(array)}
func tsNodeCallTrackerError(code,message string)tsValue{return tsFSException("TypeError",code,message)}
func tsNodeCallTrackerReport(context *tsNodeCallTrackerContext)tsValue{
	count:=len(context.calls);if tsNumber(context.expected)==float64(count){return tsU};message:=fmt.Sprintf("Expected the %s function to be executed %s time(s) but was executed %d time(s).",context.name,tsText(context.expected),count);report:=tsNewObject();report.set("message",tsStringReference(tsStringUTF8(message)));report.set("actual",tsNumberValue(float64(count)));report.set("expected",context.expected);report.set("operator",tsStringReference(tsStringUTF8(context.name)));report.set("stack",context.stack);return tsObjectValue(report)
}
func tsNodeCallTrackerConstructor(loop *tsLoop)tsValue{
	tsNodeCallTrackerStates.Do(func(){tsNodeCallTrackerStates.instances=map[*tsObject]*tsNodeCallTracker{};prototype:=tsNewObject();tsNodeCallTrackerStates.prototype=prototype
		prototype.set("calls",tsNamedNativeMethod("calls",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{
			state:=tsNodeCallTrackerState(receiver);if state==nil{tsPropertyFailure("CallTracker method called on an incompatible receiver")};process:=tsNodeProcessModule(loop);if tsTruthy(tsGet(loop,process,tsStringReference(tsStringUTF8("_exiting")))){panic(tsThrown{tsNodeCallTrackerError("ERR_UNAVAILABLE_DURING_EXIT","Cannot use CallTracker while process is exiting")})}
			callback,expected:=tsArg(args,0),tsArg(args,1);if tsIsUndefined(expected){expected=tsNumberValue(1)};if tsIsNumeric(callback){expected=callback;callback=tsFunctionValue(tsFunc(func(_ *tsLoop,_ ...tsValue)tsValue{return tsU}));(*tsFunction)(callback.ref).name="noop"}else if tsIsUndefined(callback){callback=tsFunctionValue(tsFunc(func(_ *tsLoop,_ ...tsValue)tsValue{return tsU}));(*tsFunction)(callback.ref).name="noop"};if callback.kind!=tsFunctionKind{tsFSInvalid("fn","a function or number")}
			validate:=tsGet(loop,tsInternalValidatorsModule(loop),tsStringReference(tsStringUTF8("validateUint32")));if validate.kind==tsFunctionKind{tsCall(loop,validate,expected,tsStringReference(tsStringUTF8("expected")),tsBooleanValue(true))};name:=tsGet(loop,callback,tsStringReference(tsStringUTF8("name")));if !tsTruthy(name){name=tsStringReference(tsStringUTF8("calls"))};context:=&tsNodeCallTrackerContext{expected:expected,calls:[]tsValue{},name:tsText(name),stack:tsErrorValue(&tsRuntimeError{name:"Error",message:"",fields:map[string]tsValue{}})}
			wrapped:=tsFunc(func(loop *tsLoop,values ...tsValue)tsValue{record:=tsNewObject();record.set("thisArg",tsU);record.set("arguments",tsNodeCallTrackerFreezeArray(values));record.nonExtensible=true;record.descriptors=map[string]*tsDescriptor{"thisArg":{writable:false,configurable:false},"arguments":{writable:false,configurable:false}};context.calls=append(context.calls,tsObjectValue(record));return tsCall(loop,callback,values...)});wrapped.name=tsText(name);wrapped.length=(*tsFunction)(callback.ref).length;wrapped.receiverCall=func(loop *tsLoop,thisArg tsValue,values ...tsValue)tsValue{record:=tsNewObject();record.set("thisArg",thisArg);record.set("arguments",tsNodeCallTrackerFreezeArray(values));record.nonExtensible=true;record.descriptors=map[string]*tsDescriptor{"thisArg":{writable:false,configurable:false},"arguments":{writable:false,configurable:false}};context.calls=append(context.calls,tsObjectValue(record));return tsCallReceiver(loop,callback,thisArg,values...)};tracked:=tsFunctionValue(wrapped);context.tracked=tracked;state.checks=append(state.checks,context);state.tracked[tracked.ref]=context;return tracked
		}))
		prototype.set("getCalls",tsNamedNativeMethod("getCalls",1,func(_ *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeCallTrackerState(receiver);if state==nil{tsPropertyFailure("CallTracker method called on an incompatible receiver")};tracked:=tsArg(args,0);context:=state.tracked[tracked.ref];if tracked.kind!=tsFunctionKind||context==nil{panic(tsThrown{tsNodeCallTrackerError("ERR_INVALID_ARG_VALUE","The \"tracked\" argument is not a tracked function")})};return tsNodeCallTrackerFreezeArray(context.calls)}))
		prototype.set("reset",tsNamedNativeMethod("reset",1,func(_ *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeCallTrackerState(receiver);if state==nil{tsPropertyFailure("CallTracker method called on an incompatible receiver")};tracked:=tsArg(args,0);if tsIsUndefined(tracked){for _,context:=range state.checks{context.calls=nil};return tsU};context:=state.tracked[tracked.ref];if tracked.kind!=tsFunctionKind||context==nil{panic(tsThrown{tsNodeCallTrackerError("ERR_INVALID_ARG_VALUE","The \"tracked\" argument is not a tracked function")})};context.calls=nil;return tsU}))
		prototype.set("report",tsNamedNativeMethod("report",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeCallTrackerState(receiver);if state==nil{tsPropertyFailure("CallTracker method called on an incompatible receiver")};reports:=[]tsValue{};for _,context:=range state.checks{if report:=tsNodeCallTrackerReport(context);!tsIsUndefined(report){reports=append(reports,report)}};return tsArrayValue(&tsArray{values:reports})}))
		prototype.set("verify",tsNamedNativeMethod("verify",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{reportMethod:=tsGet(loop,receiver,tsStringReference(tsStringUTF8("report")));reports:=(*tsArray)(tsCallReceiver(loop,reportMethod,receiver).ref).values;if len(reports)==0{return tsU};message:=tsGet(loop,reports[0],tsStringReference(tsStringUTF8("message")));if len(reports)>1{message=tsStringReference(tsStringUTF8("Functions were not called the expected number of times"))};options:=tsNewObject();options.set("message",message);options.set("details",tsArrayValue(&tsArray{values:reports}));panic(tsThrown{tsAssertErrorNew(loop,tsObjectValue(options))})}))
		fn:=tsFunc(func(_ *tsLoop,_ ...tsValue)tsValue{tsPropertyFailure("Class constructor CallTracker cannot be invoked without new");return tsU});fn.name="CallTracker";fn.length=0;fn.constructible=true;fn.properties=tsNewObject();fn.properties.set("prototype",tsObjectValue(prototype));fn.construct=func(_ *tsLoop,_ ...tsValue)tsValue{object:=tsNewObject();object.prototype=tsObjectValue(prototype);state:=&tsNodeCallTracker{object:object,checks:[]*tsNodeCallTrackerContext{},tracked:map[unsafe.Pointer]*tsNodeCallTrackerContext{}};tsNodeCallTrackerStates.instances[object]=state;return tsObjectValue(object)};tsNodeCallTrackerStates.constructor=tsFunctionValue(fn)
	});return tsNodeCallTrackerStates.constructor
}

func tsNodeAssertModule() tsValue {
	tsNodeAssertSingleton.Do(func(){
		base:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{if len(args)==0||!tsTruthy(args[0]){message:=tsArg(args,1);tsNodeAssertFail(loop,message,tsArg(args,0),tsBooleanValue(true),"==")};return tsU}))
		(*tsFunction)(base.ref).name="ok";(*tsFunction)(base.ref).length=1
		tsNodeAssertSingleton.module=tsNodeAssertDecorate(loop,base,false)
		strict:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{if len(args)==0||!tsTruthy(args[0]){tsNodeAssertFail(loop,tsArg(args,1),tsArg(args,0),tsBooleanValue(true),"==")};return tsU}))
		(*tsFunction)(strict.ref).name="strict";(*tsFunction)(strict.ref).length=1
		tsNodeAssertSingleton.strict=tsNodeAssertDecorate(loop,strict,true)
		(*tsFunction)(tsNodeAssertSingleton.module.ref).properties.set("strict",tsNodeAssertSingleton.strict)
		(*tsFunction)(tsNodeAssertSingleton.strict.ref).properties.set("strict",tsNodeAssertSingleton.strict)
	})
	return tsNodeAssertSingleton.module
}

`
