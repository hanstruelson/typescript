package goemit

// Stack text uses native Go frames. The accessor and setter follow the pinned
// ECMAScript stack proposal; the string format is implementation-defined.
const errorStackRuntime = `
type tsNativeErrorStack struct{sync.Once;title string;pcs []uintptr;text *tsString}
func tsCaptureErrorStack(title string)*tsNativeErrorStack{pcs:=make([]uintptr,32);count:=runtime.Callers(2,pcs);return &tsNativeErrorStack{title:title,pcs:pcs[:count]}}
func(stack *tsNativeErrorStack)string()*tsString{stack.Do(func(){var text strings.Builder;text.WriteString(stack.title);frames:=runtime.CallersFrames(stack.pcs);for{frame,more:=frames.Next();fmt.Fprintf(&text,"\n\tat %s (%s:%d)",frame.Function,frame.File,frame.Line);if !more{break}};stack.text=tsStringUTF8(text.String())});return stack.text}
func tsErrorStackGet(loop *tsLoop,receiver tsValue)tsValue{if tsIsPrimitiveValue(receiver){tsPropertyFailure("Error.stack requires an object")};var stack *tsNativeErrorStack;switch receiver.kind{case tsErrorKind:stack=(*tsRuntimeError)(receiver.ref).stack;case tsObjectKind:object:=(*tsObject)(receiver.ref);if object.nativeError{stack=object.errorStack};case tsInstanceKind:properties:=(*tsProperties)(receiver.ref);if properties.nativeError{stack=properties.errorStack}};if stack==nil{return tsU};return tsStringReference(stack.string())}
func tsErrorStackSet(loop *tsLoop,receiver tsValue,value tsValue)tsValue{
 if tsIsPrimitiveValue(receiver){tsPropertyFailure("Error.stack requires an object")};if value.kind!=tsStringKind{tsPropertyFailure("Error.stack must be a string")};if receiver.kind==tsObjectKind&&receiver.ref==unsafe.Pointer(tsErrors.prototypes["Error"]){tsPropertyFailure("Cannot assign Error.prototype.stack")}
 description:=tsOwnDescriptor(loop,receiver,"stack");if tsIsUndefined(description){description:=tsNewObject();description.set("value",value);description.set("writable",tsBooleanValue(true));description.set("enumerable",tsBooleanValue(true));description.set("configurable",tsBooleanValue(true));tsDefineProperty(loop,receiver,"stack",tsObjectValue(description))}else if !tsReflectSet(loop,receiver,"stack",value,receiver){tsPropertyFailure("Cannot assign stack property")};return tsU
}
`
