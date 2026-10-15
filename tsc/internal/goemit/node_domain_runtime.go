package goemit

// nodeDomainRuntime layers Node's legacy domain API over the local EventEmitter.
const nodeDomainRuntime = `
type tsNodeDomainState struct { sync.Once; value *tsObject; stack *tsArray; active tsValue }
var tsNodeDomainSingleton tsNodeDomainState

func tsNodeDomainSetActive(value tsValue) {
	tsNodeDomainSingleton.active=value
	if tsNodeDomainSingleton.value!=nil { tsNodeDomainSingleton.value.set("active",value) }
}

func tsNodeDomainEmitError(loop *tsLoop,domain tsValue,err tsValue) {
	if tsNullish(err){err=tsErrorValue(&tsRuntimeError{name:"Error",message:"No error supplied to domain"})}
	if err.kind==tsErrorKind||err.kind==tsObjectKind {
		tsSet(loop,err,tsStringReference(tsStringUTF8("domainEmitter")),domain)
		tsSet(loop,err,tsStringReference(tsStringUTF8("domain")),domain)
		tsSet(loop,err,tsStringReference(tsStringUTF8("domainThrown")),tsBooleanValue(false))
	}
	tsCallReceiver(loop,tsGet(loop,domain,tsStringReference(tsStringUTF8("emit"))),domain,tsStringReference(tsStringUTF8("error")),err)
}

func tsNodeDomainEnter(domain tsValue) tsValue {
	tsNodeDomainSingleton.stack.values=append(tsNodeDomainSingleton.stack.values,domain)
	tsNodeDomainSetActive(domain)
	return domain
}

func tsNodeDomainExit(domain tsValue) tsValue {
	stack:=tsNodeDomainSingleton.stack.values
	index:=-1
	for at:=len(stack)-1;at>=0;at--{if tsStrictEqual(stack[at],domain){index=at;break}}
	if index<0{return domain}
	for len(tsNodeDomainSingleton.stack.values)>index{last:=len(tsNodeDomainSingleton.stack.values)-1;tsNodeDomainSingleton.stack.values[last]=tsU;tsNodeDomainSingleton.stack.values=tsNodeDomainSingleton.stack.values[:last]}
	if len(tsNodeDomainSingleton.stack.values)>0{tsNodeDomainSetActive(tsNodeDomainSingleton.stack.values[len(tsNodeDomainSingleton.stack.values)-1])}else{tsNodeDomainSetActive(tsNull)}
	return domain
}

func tsNodeDomainCreate(loop *tsLoop) tsValue {
	domain:=tsNodeEventCreate(loop,tsU,tsU)
	object:=(*tsObject)(domain.ref)
	attached:=map[*tsObject]tsValue{}
	object.set("add",tsNamedNativeMethod("add",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{
		emitter:=tsArg(args,0);method:=tsGet(loop,emitter,tsStringReference(tsStringUTF8("on")))
		handler:=tsNamedNativeMethod("emitError",1,func(loop *tsLoop,_ tsValue,values ...tsValue)tsValue{err:=tsArg(values,0);func(){defer func(){if failure:=recover();failure!=nil{tsNodeDomainEmitError(loop,domain,tsUnwrap(failure))}}();tsNodeDomainEmitError(loop,domain,err)}();return tsU})
		if emitter.kind==tsObjectKind{attached[(*tsObject)(emitter.ref)]=handler}
		tsCallReceiver(loop,method,emitter,tsStringReference(tsStringUTF8("error")),handler);return receiver
	}))
	object.set("remove",tsNamedNativeMethod("remove",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{
		emitter:=tsArg(args,0);method:=tsGet(loop,emitter,tsStringReference(tsStringUTF8("removeListener")))
		if method.kind!=tsFunctionKind{return receiver}
		if emitter.kind==tsObjectKind{if handler,ok:=attached[(*tsObject)(emitter.ref)];ok{tsCallReceiver(loop,method,emitter,tsStringReference(tsStringUTF8("error")),handler);delete(attached,(*tsObject)(emitter.ref))}}
		return receiver
	}))
	object.set("bind",tsNamedNativeMethod("bind",1,func(_ *tsLoop,receiver tsValue,args ...tsValue)tsValue{
		callback:=tsArg(args,0);if callback.kind!=tsFunctionKind{tsFSInvalid("callback","a function")}
		return tsFunctionValue(tsFunc(func(loop *tsLoop,values ...tsValue)tsValue{var result tsValue;func(){defer func(){if failure:=recover();failure!=nil{tsNodeDomainEmitError(loop,domain,tsUnwrap(failure))}}();result=tsCall(loop,callback,values...)}();return result}))
	}))
	object.set("intercept",tsNamedNativeMethod("intercept",1,func(_ *tsLoop,receiver tsValue,args ...tsValue)tsValue{
		callback:=tsArg(args,0);if callback.kind!=tsFunctionKind{tsFSInvalid("callback","a function")}
		return tsFunctionValue(tsFunc(func(loop *tsLoop,values ...tsValue)tsValue{err:=tsArg(values,0);if tsTruthy(err){tsNodeDomainEmitError(loop,domain,err);return tsU};rest:=[]tsValue{};if len(values)>1{rest=values[1:]};var result tsValue;func(){defer func(){if failure:=recover();failure!=nil{tsNodeDomainEmitError(loop,domain,tsUnwrap(failure))}}();result=tsCall(loop,callback,rest...)}();return result}))
	}))
	object.set("run",tsNamedNativeMethod("run",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{
		callback:=tsArg(args,0);if callback.kind!=tsFunctionKind{tsFSInvalid("fn","a function")};values:=[]tsValue{};if len(args)>1{values=args[1:]}
		tsNodeDomainEnter(domain);var result tsValue
		func(){defer func(){failure:=recover();tsNodeDomainExit(domain);if failure!=nil{tsNodeDomainEmitError(loop,domain,tsUnwrap(failure))}}();result=tsCallReceiver(loop,callback,domain,values...)}()
		return result
	}))
	object.set("dispose",tsNamedNativeMethod("dispose",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{tsCallReceiver(loop,tsGet(loop,receiver,tsStringReference(tsStringUTF8("removeAllListeners"))),receiver);return receiver}))
	object.set("enter",tsNamedNativeMethod("enter",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{return tsNodeDomainEnter(domain)}))
	object.set("exit",tsNamedNativeMethod("exit",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{return tsNodeDomainExit(domain)}))
	return domain
}

func tsNodeDomainModule(loop *tsLoop) tsValue {
	tsNodeDomainSingleton.Do(func(){
		tsNodeEventsInit(loop)
		tsNodeDomainSingleton.stack=&tsArray{}
		out:=tsNewObject()
		create:=tsNamedNativeMethod("createDomain",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNodeDomainCreate(loop)})
		out.set("createDomain",create);out.set("create",create)
		out.set("_stack",tsArrayValue(tsNodeDomainSingleton.stack));out.set("active",tsNull)
		tsNodeDomainSingleton.value=out;tsNodeDomainSingleton.active=tsNull
	})
	return tsObjectValue(tsNodeDomainSingleton.value)
}
`
