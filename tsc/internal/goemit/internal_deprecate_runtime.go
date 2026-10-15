package goemit

// internalDeprecateRuntime maps Bun's deprecation wrapper and per-code warning
// suppression to the local process warning event.
const internalDeprecateRuntime = `
var tsInternalDeprecate struct { sync.Once; module *tsObject; warned map[string]bool }

func tsInternalDeprecateEmit(loop *tsLoop,code,message tsValue,deprecated bool,local *bool) {
    if *local||tsTruthy(tsGet(loop,tsNodeProcessModule(loop),tsStringReference(tsStringUTF8("noDeprecation")))){return}
    process:=tsNodeProcessModule(loop)
    if !tsIsUndefined(code) {
        key:=tsInternalMIMEString(loop,code)
        if tsInternalDeprecate.warned[key]{*local=true;return}
        tsInternalDeprecate.warned[key]=true
    }
    *local=true
    emit:=tsGet(loop,process,tsStringReference(tsStringUTF8("emitWarning")))
    args:=[]tsValue{message,tsStringReference(tsStringUTF8("DeprecationWarning"))}
    if !tsIsUndefined(code){args=append(args,code)}
    args=append(args,tsBooleanValue(deprecated))
    tsCallReceiver(loop,emit,process,args...)
}

func tsInternalDeprecateModule(loop *tsLoop) tsValue {
    tsInternalDeprecate.Do(func(){
        tsInternalDeprecate.warned=map[string]bool{}
        module:=tsNewObject()
        validate:=tsGet(loop,tsInternalValidatorsModule(loop),tsStringReference(tsStringUTF8("validateString")))
        deprecate:=tsNamedNativeMethod("deprecate",3,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{
            original,message,code:=tsArg(args,0),tsArg(args,1),tsArg(args,2)
            if !tsIsUndefined(code){tsCall(loop,validate,code,tsStringReference(tsStringUTF8("code")))}
            warned:=false
            wrapper:=&tsFunction{name:"deprecated",length:0,constructible:true}
            wrapper.call=func(loop *tsLoop,args ...tsValue)tsValue{tsInternalDeprecateEmit(loop,code,message,false,&warned);return tsCallReceiver(loop,original,tsU,args...)}
            wrapper.receiverCall=func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{tsInternalDeprecateEmit(loop,code,message,false,&warned);return tsCallReceiver(loop,original,receiver,args...)}
            wrapper.construct=func(loop *tsLoop,args ...tsValue)tsValue{tsInternalDeprecateEmit(loop,code,message,false,&warned);return tsConstruct(original,args...)}
            deprecated:=tsFunctionValue(wrapper);tsSetPrototype(deprecated,original)
            originalPrototype:=tsGet(loop,original,tsStringReference(tsStringUTF8("prototype")))
            if tsTruthy(originalPrototype){descriptor:=tsNewObject();descriptor.set("value",originalPrototype);descriptor.set("writable",tsBooleanValue(true));tsDefineProperty(loop,deprecated,"prototype",tsObjectValue(descriptor))}
            lengthDescriptor:=tsOwnDescriptor(loop,original,"length")
            if !tsIsUndefined(lengthDescriptor){tsDefineProperty(loop,deprecated,"length",lengthDescriptor)}
            return deprecated
        })
        module.set("deprecate",deprecate)
        tsInternalDeprecate.module=module
    })
    return tsObjectValue(tsInternalDeprecate.module)
}
`
