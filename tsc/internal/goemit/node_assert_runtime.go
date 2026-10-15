package goemit

// nodeAssertRuntime provides the Node assertion entry point using the
// transpiler's own value, error, promise, and RegExp machinery.
const nodeAssertRuntime = `
type tsNodeAssertState struct { sync.Once; module, strict tsValue }
var tsNodeAssertSingleton tsNodeAssertState

type tsNodeAssertPair struct { left, right unsafe.Pointer; kind tsKind }

func tsNodeAssertError(message tsValue, actual, expected tsValue, operator string) tsValue {
	if message.kind==tsErrorKind || message.kind==tsObjectKind && (*tsObject)(message.ref).nativeError { return message }
	text:="Assertion failed"; if !tsIsUndefined(message) && !tsNullish(message) && tsText(message)!="" { text=tsText(message) }
	fields:=map[string]tsValue{"actual":actual,"expected":expected,"operator":tsStringReference(tsStringUTF8(operator)),"code":tsStringReference(tsStringUTF8("ERR_ASSERTION")),"generatedMessage":tsBooleanValue(tsIsUndefined(message)||tsNullish(message)||tsText(message)=="")}
	return tsErrorValue(&tsRuntimeError{name:"AssertionError",message:text,fields:fields})
}

func tsNodeAssertFail(message tsValue, actual, expected tsValue, operator string) { panic(tsThrown{tsNodeAssertError(message,actual,expected,operator)}) }

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
		if !settled.rejected && shouldReject {promise.settle(tsResult{tsNodeAssertError(message,settled.value,expected,"rejects"),true});return}
		if settled.rejected && !shouldReject {promise.settle(tsResult{tsNodeAssertError(message,settled.value,tsU,"doesNotReject"),true});return}
		promise.settle(tsResult{tsNodeAssertError(message,settled.value,expected,"rejects"),true})
	})
	return tsPromiseValue(promise)
}

func tsNodeAssertMethod(name string) tsValue {
	fn:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{
		actual,expected,message:=tsArg(args,0),tsArg(args,1),tsArg(args,2)
		switch name {
		case "ok": if !tsTruthy(actual){tsNodeAssertFail(expected,actual,tsBooleanValue(true),"==")}
		case "fail": tsNodeAssertFail(actual,tsU,tsU,"fail")
		case "equal","notEqual": if len(args)<2{tsNodeAssertMissing("expected")};same:=tsLooseEqual(actual,expected)||(tsIsNumeric(actual)&&tsIsNumeric(expected)&&math.IsNaN(tsNumber(actual))&&math.IsNaN(tsNumber(expected)));if (name=="equal"&&!same)||(name=="notEqual"&&same){operator:="==";if name=="notEqual"{operator="!="};tsNodeAssertFail(message,actual,expected,operator)}
		case "strictEqual","notStrictEqual": if len(args)<2{tsNodeAssertMissing("expected")};same:=tsSameValue(actual,expected);if (name=="strictEqual"&&!same)||(name=="notStrictEqual"&&same){operator:="strictEqual";if name=="notStrictEqual"{operator="notStrictEqual"};tsNodeAssertFail(message,actual,expected,operator)}
		case "deepEqual","notDeepEqual","deepStrictEqual","notDeepStrictEqual","partialDeepStrictEqual":
			if len(args)<2{tsNodeAssertMissing("expected")};strict:=name=="deepStrictEqual"||name=="notDeepStrictEqual"||name=="partialDeepStrictEqual";same:=tsNodeAssertDeep(loop,actual,expected,strict,map[tsNodeAssertPair]bool{});if name=="partialDeepStrictEqual"{same=tsNodeAssertMatch(loop,actual,expected)};if name=="notDeepEqual"||name=="notDeepStrictEqual"{same=!same};if !same{tsNodeAssertFail(message,actual,expected,name)}
		case "throws","doesNotThrow":
			if len(args)==0{tsNodeAssertMissing("fn")};callback:=actual;want:=expected;msg:=message;if len(args)==2&&want.kind==tsStringKind{msg=want;want=tsU};thrown,didThrow:=tsNodeAssertThrown(loop,callback)
			if name=="throws" {if !didThrow||!tsIsUndefined(want)&&!tsNodeAssertMatch(loop,thrown,want){tsNodeAssertFail(msg,thrown,want,"throws")}} else if didThrow {if tsIsUndefined(want)||tsNodeAssertMatch(loop,thrown,want){tsNodeAssertFail(msg,thrown,tsU,"doesNotThrow")};panic(tsThrown{thrown})}
		case "rejects","doesNotReject": if len(args)==2&&expected.kind==tsStringKind{return tsNodeAssertAsync(loop,actual,tsU,expected,name=="rejects")};return tsNodeAssertAsync(loop,actual,expected,message,name=="rejects")
		case "ifError": if tsTruthy(actual){tsNodeAssertFail(message,actual,tsNull,"ifError")}
		case "match","doesNotMatch":
			if actual.kind!=tsStringKind||expected.kind!=tsRegExpKind {tsPropertyFailure("The string and regexp arguments are required")};matched:=tsRegExpExecReceiver(loop,expected,(*tsString)(actual.ref)).kind!=tsNullKind;if (name=="match"&&!matched)||(name=="doesNotMatch"&&matched){tsNodeAssertFail(message,actual,expected,name)}
		}
		return tsU
	})
	fn.name=name;fn.length=2
	return tsFunctionValue(fn)
}

func tsNodeAssertDecorate(value tsValue, strict bool) tsValue {
	fn:=(*tsFunction)(value.ref);fn.properties=tsNewObject()
	for _,name:=range []string{"ok","fail","equal","notEqual","deepEqual","notDeepEqual","deepStrictEqual","notDeepStrictEqual","strictEqual","notStrictEqual","partialDeepStrictEqual","throws","doesNotThrow","rejects","doesNotReject","ifError","match","doesNotMatch"}{method:=name;if strict{switch name{case "equal":method="strictEqual";case "notEqual":method="notStrictEqual";case "deepEqual":method="deepStrictEqual";case "notDeepEqual":method="notDeepStrictEqual"}};fn.properties.set(name,tsNodeAssertMethod(method))}
	fn.properties.set("AssertionError",tsNodeAssertionErrorConstructor())
	fn.properties.set("CallTracker",tsNodeCallTrackerConstructor())
	fn.properties.set("Assert",tsNodeAssertConstructor())
	return value
}

func tsNodeAssertionErrorConstructor() tsValue {
	fn:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{options:=tsArg(args,0);message:=tsFSProperty(options,"message");actual:=tsFSProperty(options,"actual");expected:=tsFSProperty(options,"expected");operator:=tsText(tsFSProperty(options,"operator"));text:=tsText(message);return tsErrorValue(&tsRuntimeError{name:"AssertionError",message:text,fields:map[string]tsValue{"actual":actual,"expected":expected,"operator":tsStringReference(tsStringUTF8(operator)),"code":tsStringReference(tsStringUTF8("ERR_ASSERTION"))}})})
	fn.name="AssertionError";fn.length=1;fn.constructible=true;fn.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsCall(loop,tsFunctionValue(fn),args...)};return tsFunctionValue(fn)
}

func tsNodeAssertConstructor() tsValue {
	fn:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{tsPropertyFailure("Class constructor Assert cannot be invoked without new")});fn.name="Assert";fn.length=1;fn.constructible=true
	fn.construct=func(loop *tsLoop,args ...tsValue)tsValue{instance:=tsNewObject();options:=tsArg(args,0);strictValue:=tsFSProperty(options,"strict");strict:=tsIsUndefined(strictValue)||tsTruthy(strictValue);for _,name:=range []string{"ok","fail","equal","notEqual","deepEqual","notDeepEqual","deepStrictEqual","notDeepStrictEqual","strictEqual","notStrictEqual","partialDeepStrictEqual","throws","doesNotThrow","rejects","doesNotReject","ifError","match","doesNotMatch"}{method:=name;if strict{switch name{case "equal":method="strictEqual";case "notEqual":method="notStrictEqual";case "deepEqual":method="deepStrictEqual";case "notDeepEqual":method="notDeepStrictEqual"}};instance.set(name,tsNodeAssertMethod(method))};return tsObjectValue(instance)}
	return tsFunctionValue(fn)
}

type tsNodeCallTracker struct { counts map[unsafe.Pointer]int; expected map[unsafe.Pointer]int; functions map[unsafe.Pointer]tsValue }
func tsNodeCallTrackerConstructor() tsValue {
	fn:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{state:=&tsNodeCallTracker{counts:map[unsafe.Pointer]int{},expected:map[unsafe.Pointer]int{},functions:map[unsafe.Pointer]tsValue{}};out:=tsNewObject()
		out.set("calls",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{callback:=tsArg(args,0);if callback.kind!=tsFunctionKind{tsPropertyFailure("The fn argument must be a function")};count:=1;if len(args)>1{count=int(tsNumber(args[1]))};state.expected[callback.ref]=count;wrapped:=tsFunc(func(loop *tsLoop,values ...tsValue)tsValue{state.counts[callback.ref]++;return tsCall(loop,callback,values...)});state.functions[wrapped]=callback;return tsFunctionValue(wrapped)})))
		out.set("getCalls",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{items:=&tsArray{};for key,callback:=range state.functions{items.values=append(items.values,tsObjectValue(tsNewObject()));_ = key;_ = callback};return tsArrayValue(items)})))
		out.set("reset",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{state.counts=map[unsafe.Pointer]int{};return tsU})))
		out.set("verify",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{for key,want:=range state.expected{if state.counts[key]!=want{tsNodeAssertFail(tsU,tsNumberValue(float64(state.counts[key])),tsNumberValue(float64(want)),"CallTracker.verify")}};return tsU})))
		out.set("report",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsArrayValue(&tsArray{})})))
		return tsObjectValue(out)
	});fn.name="CallTracker";fn.length=0;fn.constructible=true;fn.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsCall(loop,tsFunctionValue(fn),args...)};return tsFunctionValue(fn)
}

func tsNodeAssertModule() tsValue {
	tsNodeAssertSingleton.Do(func(){
		base:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{if len(args)==0||!tsTruthy(args[0]){message:=tsArg(args,1);tsNodeAssertFail(message,tsArg(args,0),tsBooleanValue(true),"==")};return tsU}))
		(*tsFunction)(base.ref).name="ok";(*tsFunction)(base.ref).length=1
		tsNodeAssertSingleton.module=tsNodeAssertDecorate(base,false)
		strict:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{if len(args)==0||!tsTruthy(args[0]){tsNodeAssertFail(tsArg(args,1),tsArg(args,0),tsBooleanValue(true),"==")};return tsU}))
		(*tsFunction)(strict.ref).name="strict";(*tsFunction)(strict.ref).length=1
		tsNodeAssertSingleton.strict=tsNodeAssertDecorate(strict,true)
		(*tsFunction)(tsNodeAssertSingleton.module.ref).properties.set("strict",tsNodeAssertSingleton.strict)
		(*tsFunction)(tsNodeAssertSingleton.strict.ref).properties.set("strict",tsNodeAssertSingleton.strict)
	})
	return tsNodeAssertSingleton.module
}

`
