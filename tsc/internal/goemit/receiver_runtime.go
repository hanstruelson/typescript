package goemit

const receiverRuntime = `
var tsGlobalNamespace struct{sync.Once;object *tsObject}
func tsGlobalObject()tsValue{tsGlobalNamespace.Do(func(){object:=tsNewObject();tsGlobalNamespace.object=object;object.set("globalThis",tsObjectValue(object));object.set("Object",tsObjectConstructor());object.set("Function",tsFunctionConstructor());object.set("Math",tsMathModule());object.set("JSON",tsJSONModule());object.set("Symbol",tsSymbolFunction());object.set("Reflect",tsReflectModule());for _,name:=range []string{"Error","TypeError","RangeError","ReferenceError","SyntaxError","URIError","EvalError","AggregateError"}{object.set(name,tsErrorConstructor(name))};object.descriptors=map[string]*tsDescriptor{};for _,name:=range object.order{object.descriptors[name]=&tsDescriptor{writable:true,configurable:true}}});return tsObjectValue(tsGlobalNamespace.object)}
func tsFunctionReceiver(receiver tsValue,strict bool)tsValue{
 if strict{return receiver};if tsNullish(receiver){return tsGlobalObject()};if !tsIsPrimitiveValue(receiver){return receiver}
 object:=tsNewObject();object.boxedPrimitive=receiver
 object.set("valueOf",tsNativeMethod(func(loop *tsLoop,self tsValue,args ...tsValue)tsValue{return receiver}))
 object.set("toString",tsNativeMethod(func(loop *tsLoop,self tsValue,args ...tsValue)tsValue{if receiver.kind==tsSymbolKind{return tsStringReference(tsStringUTF8(tsSymbolText((*tsSymbol)(receiver.ref))))};return tsStringReference(tsStringValue(loop,receiver))}))
 return tsObjectValue(object)
}
`
