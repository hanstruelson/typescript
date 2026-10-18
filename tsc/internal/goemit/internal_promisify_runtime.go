package goemit

// internalPromisifyRuntime maps Bun's callback-to-Promise helper and custom
// promisify symbols onto the local function and Promise representations.
const internalPromisifyRuntime = `
var tsInternalPromisify struct { sync.Once; module *tsObject; customSymbol,argsSymbol tsValue }

func tsInternalPromisifyInit(loop *tsLoop) {
    tsInternalPromisify.Do(func(){
        tsInternalSharedModule(loop)
        symbolFor:=tsGet(loop,tsSymbolFunction(),tsStringReference(tsStringUTF8("for")))
        tsInternalPromisify.customSymbol=tsCall(loop,symbolFor,tsStringReference(tsStringUTF8("nodejs.util.promisify.custom")))
        tsInternalPromisify.argsSymbol=tsInternalSharedSingleton.customPromisifyArgs
    })
}

func tsInternalPromisifyDefine(loop *tsLoop,target,value,symbol tsValue,enumerable,configurable bool) tsValue {
    descriptor:=tsNewObject();descriptor.set("value",value);descriptor.set("enumerable",tsBooleanValue(enumerable));descriptor.set("configurable",tsBooleanValue(configurable))
    tsDefineProperty(loop,target,tsPropertyKey(symbol),tsObjectValue(descriptor))
    return value
}

func tsInternalPromisifyCloneProperties(loop *tsLoop,target,source tsValue) {
    keys:=tsOwnKeys(source,false)
    for _,key:=range keys.values {
        descriptor:=tsOwnDescriptor(loop,source,tsPropertyKey(key))
        if tsIsUndefined(descriptor){continue}
        if !tsTryDefineProperty(loop,target,tsPropertyKey(key),descriptor){tsPropertyFailure("Cannot define promisified function property")}
    }
}

func tsInternalPromisifyFunction(loop *tsLoop,original,callbackArgs tsValue) tsValue {
    function:=&tsFunction{name:"",length:0,constructible:true}
    function.call=func(loop *tsLoop,args ...tsValue)tsValue{return tsInternalPromisifyInvoke(loop,tsFunctionValue(function),tsU,original,callbackArgs,args)}
    function.receiverCall=func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsInternalPromisifyInvoke(loop,tsFunctionValue(function),receiver,original,callbackArgs,args)}
    wrapped:=tsFunctionValue(function)
    tsSetPrototype(wrapped,tsPrototypeOf(original))
    tsInternalPromisifyDefine(loop,wrapped,wrapped,tsInternalPromisify.customSymbol,true,true)
    tsInternalPromisifyCloneProperties(loop,wrapped,original)
    return wrapped
}

func tsInternalPromisifyInvoke(loop *tsLoop,_ tsValue,receiver,original,callbackArgs tsValue,args []tsValue) tsValue {
    promise:=loop.promise()
    callback:=tsFunctionValue(tsFunc(func(loop *tsLoop,values ...tsValue)tsValue{
        failure:=tsArg(values,0)
        if tsTruthy(failure){promise.settle(tsResult{failure,true});return tsU}
        result:=tsArg(values,1)
        if !tsIsUndefined(callbackArgs) {
            object:=tsNewObject();count:=0
            switch callbackArgs.kind {case tsArrayKind:count=(*tsArray)(callbackArgs.ref).length();case tsStringKind:count=len((*tsString)(callbackArgs.ref).units);default:count=int(tsNumber(tsGet(loop,callbackArgs,tsStringReference(tsStringUTF8("length")))))}
            for index:=0;index<count;index++ {
                key:=tsU
                if callbackArgs.kind==tsArrayKind{key=(*tsArray)(callbackArgs.ref).at(index)}else if callbackArgs.kind==tsStringKind{key=tsStringReference((*tsString)(callbackArgs.ref).slice(index,index+1))}else{key=tsGet(loop,callbackArgs,tsStringReference(tsStringUTF8(strconv.Itoa(index))))}
                item:=tsArg(values,index+1)
                tsSet(loop,tsObjectValue(object),key,item)
            }
            result=tsObjectValue(object)
        }
        promise.settle(tsResult{result,false})
        return tsU
    }))
    invocation:=append(append([]tsValue(nil),args...),callback)
    var maybePromise tsValue
    func(){defer func(){if failed:=recover();failed!=nil{promise.settle(tsResult{tsUnwrap(failed),true})}}();maybePromise=tsCallReceiver(loop,original,receiver,invocation...)}()
    if maybePromise.kind==tsPromiseKind {
        process:=tsNodeProcessModule(loop);emitWarning:=tsGet(loop,process,tsStringReference(tsStringUTF8("emitWarning")))
        tsCallReceiver(loop,emitWarning,process,tsStringReference(tsStringUTF8("Calling promisify on a function that returns a Promise is likely a mistake.")),tsStringReference(tsStringUTF8("DeprecationWarning")),tsStringReference(tsStringUTF8("DEP0174")))
    }
    return tsPromiseValue(promise)
}

func tsInternalPromisifyModule(loop *tsLoop) tsValue {
    tsInternalPromisifyInit(loop)
    if tsInternalPromisify.module!=nil{return tsObjectValue(tsInternalPromisify.module)}
    module:=tsNewObject()
    validate:=tsGet(loop,tsInternalValidatorsModule(loop),tsStringReference(tsStringUTF8("validateFunction")))
    defineArgs:=tsNamedNativeMethod("defineCustomPromisifyArgs",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{
        target,items:=tsArg(args,0),tsArg(args,1);tsInternalPromisifyDefine(loop,target,items,tsInternalPromisify.argsSymbol,false,false);return items
    })
    promisify:=tsNamedNativeMethod("promisify",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{
        original:=tsArg(args,0);tsCall(loop,validate,original,tsStringReference(tsStringUTF8("original")))
        custom:=tsGet(loop,original,tsInternalPromisify.customSymbol)
        if tsTruthy(custom){tsCall(loop,validate,custom,tsStringReference(tsStringUTF8("custom")));tsInternalPromisifyDefine(loop,custom,custom,tsInternalPromisify.customSymbol,false,true);return custom}
        callbackArgs:=tsGet(loop,original,tsInternalPromisify.argsSymbol)
        return tsInternalPromisifyFunction(loop,original,callbackArgs)
    })
    tsSet(loop,tsFunctionValue((*tsFunction)(promisify.ref)),tsStringReference(tsStringUTF8("custom")),tsInternalPromisify.customSymbol)
    module.set("defineCustomPromisifyArgs",defineArgs);module.set("promisify",promisify);tsInternalPromisify.module=module
    return tsObjectValue(module)
}
`
