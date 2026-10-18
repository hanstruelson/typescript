package goemit

// Stack text uses native Go frames. The accessor and setter follow the pinned
// ECMAScript stack proposal; the string format is implementation-defined.
const errorStackRuntime = `
type tsNativeErrorStack struct{sync.Once;title string;pcs []uintptr;text *tsString}
func tsCaptureErrorStack(title string)*tsNativeErrorStack{pcs:=make([]uintptr,32);count:=runtime.Callers(2,pcs);return &tsNativeErrorStack{title:title,pcs:pcs[:count]}}
func(stack *tsNativeErrorStack)string()*tsString{stack.Do(func(){var text strings.Builder;text.WriteString(stack.title);frames:=runtime.CallersFrames(stack.pcs);for{frame,more:=frames.Next();fmt.Fprintf(&text,"\n\tat %s (%s:%d)",frame.Function,frame.File,frame.Line);if !more{break}};stack.text=tsStringUTF8(text.String())});return stack.text}
func tsErrorCallSite(frame runtime.Frame)tsValue{
 object:=tsNewObject();file:=frame.File;function:=frame.Function;line:=tsNumberValue(float64(frame.Line));column:=tsNumberValue(1)
 object.set("getFileName",tsNamedNativeMethod("getFileName",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{if file==""{return tsU};return tsStringReference(tsStringUTF8(file))}))
 object.set("getLineNumber",tsNamedNativeMethod("getLineNumber",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return line}))
 object.set("getColumnNumber",tsNamedNativeMethod("getColumnNumber",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return column}))
 object.set("getFunctionName",tsNamedNativeMethod("getFunctionName",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{if function==""{return tsU};return tsStringReference(tsStringUTF8(function))}))
 object.set("getMethodName",tsNamedNativeMethod("getMethodName",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{if function==""{return tsU};return tsStringReference(tsStringUTF8(function))}))
 object.set("getTypeName",tsNamedNativeMethod("getTypeName",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsU}))
 object.set("getThis",tsNamedNativeMethod("getThis",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsU}))
 object.set("getEvalOrigin",tsNamedNativeMethod("getEvalOrigin",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsU}))
 object.set("isEval",tsNamedNativeMethod("isEval",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsBooleanValue(false)}))
 object.set("isNative",tsNamedNativeMethod("isNative",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsBooleanValue(false)}))
 object.set("isConstructor",tsNamedNativeMethod("isConstructor",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsBooleanValue(false)}))
 object.set("toString",tsNamedNativeMethod("toString",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{if function==""{return tsStringReference(tsStringUTF8(fmt.Sprintf("%s:%d:%d",file,frame.Line,1)))};return tsStringReference(tsStringUTF8(fmt.Sprintf("%s (%s:%d:%d)",function,file,frame.Line,1)))}))
 return tsObjectValue(object)
}
func(stack *tsNativeErrorStack)callSites()tsValue{frames:=&tsArray{};iterator:=runtime.CallersFrames(stack.pcs);for{frame,more:=iterator.Next();frames.values=append(frames.values,tsErrorCallSite(frame));if !more{break}};return tsArrayValue(frames)}
func tsErrorStackGet(loop *tsLoop,receiver tsValue)tsValue{if tsIsPrimitiveValue(receiver){tsPropertyFailure("Error.stack requires an object")};var stack *tsNativeErrorStack;switch receiver.kind{case tsErrorKind:stack=(*tsRuntimeError)(receiver.ref).stack;case tsObjectKind:object:=(*tsObject)(receiver.ref);if object.nativeError{stack=object.errorStack};case tsInstanceKind:properties:=(*tsProperties)(receiver.ref);if properties.nativeError{stack=properties.errorStack}};if stack==nil{return tsU};if tsErrors.constructors!=nil{errorValue:=tsFunctionValue(tsErrors.constructors["Error"]);prepare:=tsFunctionGet(loop,errorValue,"prepareStackTrace");if prepare.kind==tsFunctionKind{return tsCall(loop,prepare,receiver,stack.callSites())}};return tsStringReference(stack.string())}
func tsErrorStackSet(loop *tsLoop,receiver tsValue,value tsValue)tsValue{
 if tsIsPrimitiveValue(receiver){tsPropertyFailure("Error.stack requires an object")};if value.kind!=tsStringKind{tsPropertyFailure("Error.stack must be a string")};if receiver.kind==tsObjectKind&&receiver.ref==unsafe.Pointer(tsErrors.prototypes["Error"]){tsPropertyFailure("Cannot assign Error.prototype.stack")}
 description:=tsOwnDescriptor(loop,receiver,"stack");if tsIsUndefined(description){description:=tsNewObject();description.set("value",value);description.set("writable",tsBooleanValue(true));description.set("enumerable",tsBooleanValue(true));description.set("configurable",tsBooleanValue(true));tsDefineProperty(loop,receiver,"stack",tsObjectValue(description))}else if !tsReflectSet(loop,receiver,"stack",value,receiver){tsPropertyFailure("Cannot assign stack property")};return tsU
}
`
