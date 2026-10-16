package goemit

const receiverRuntime = `
var tsGlobalNamespace struct{sync.Once;object *tsObject}
var tsGlobalConsoleInstall sync.Once
func tsGlobalObject()tsValue{tsGlobalNamespace.Do(func(){
 object:=tsNewObject();tsGlobalNamespace.object=object
 object.set("globalThis",tsObjectValue(object));object.set("Object",tsObjectConstructor());object.set("Function",tsFunctionConstructor())
 object.set("Array",tsBuiltinClass("Array"));object.set("String",tsStringConstructor());object.set("Number",tsNumberConstructor(loop));object.set("Boolean",tsBooleanConstructor());object.set("BigInt",tsBigIntFunction());object.set("Date",tsClassValue(tsNativeClass("Date")));object.set("RegExp",tsRegExpConstructor())
 for _,name:=range []string{"Map","Set","ArrayBuffer","DataView","Float16Array","Float64Array","Float32Array","Int8Array","Uint8Array","Uint8ClampedArray","Int16Array","Uint16Array","Int32Array","Uint32Array","BigInt64Array","BigUint64Array"}{object.set(name,tsBuiltinClass(name))}
 object.set("Math",tsMathModule());object.set("JSON",tsJSONModule());object.set("Symbol",tsSymbolFunction());object.set("Reflect",tsReflectModule());object.set("isNaN",tsIsNaNFunction());object.set("print",tsPrintFunction());object.set("Proxy",tsProxyConstructor());object.set("Intl",tsIntlModule());object.set("EventTarget",tsEventTargetConstructor());object.set("Event",tsEventConstructor())
 for _,name:=range []string{"Error","TypeError","RangeError","ReferenceError","SyntaxError","URIError","EvalError","AggregateError","SuppressedError"}{object.set(name,tsErrorConstructor(name))}
 object.descriptors=map[string]*tsDescriptor{};for _,name:=range object.order{object.descriptors[name]=&tsDescriptor{writable:true,configurable:true}}
 });tsGlobalConsoleInstall.Do(func(){global:=tsGlobalNamespace.object;if global.get("console").kind==tsUndefinedKind{global.set("console",tsNodeConsoleModule())}});return tsObjectValue(tsGlobalNamespace.object)}
func tsGlobalIdentifier(loop *tsLoop,name string)tsValue{global:=tsGlobalObject();key:=tsStringReference(tsStringUTF8(name));if tsHas(global,key){return tsGet(loop,global,key)};panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"ReferenceError",message:name+" is not defined"})})}
func tsFunctionReceiver(receiver tsValue,strict bool)tsValue{
 if strict{return receiver};if tsNullish(receiver){return tsGlobalObject()};if !tsIsPrimitiveValue(receiver){return receiver}
 object:=tsNewObject();object.boxedPrimitive=receiver
 object.set("valueOf",tsNativeMethod(func(loop *tsLoop,self tsValue,args ...tsValue)tsValue{return receiver}))
 object.set("toString",tsNativeMethod(func(loop *tsLoop,self tsValue,args ...tsValue)tsValue{if receiver.kind==tsSymbolKind{return tsStringReference(tsStringUTF8(tsSymbolText((*tsSymbol)(receiver.ref))))};return tsStringReference(tsStringValueWithLoop(loop,receiver))}))
 return tsObjectValue(object)
}
`
