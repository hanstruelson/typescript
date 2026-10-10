package goemit

const errorRuntime = `
func tsErrorObject(loop *tsLoop,error *tsRuntimeError)*tsObject{if error.object!=nil{return error.object};object:=tsNewObject();object.nativeError=true;object.errorStack=error.stack;object.prototype=tsErrorPrototype(error.name);object.descriptors=map[string]*tsDescriptor{};object.set("message",tsStringReference(tsStringUTF8(error.message)));object.descriptors["message"]=&tsDescriptor{writable:true,configurable:true};for name,value:=range error.fields{object.set(name,value)};error.object=object;return object}
func tsErrorSuper(loop *tsLoop,properties *tsProperties,name string,args []tsValue){
 error:=(*tsObject)(tsNewError(loop,name,args).ref)
 if properties.initialized{panic(tsErrorValue(&tsRuntimeError{name:"ReferenceError",message:"Super constructor may only be called once"}))}
 for _,key:=range error.order{properties.define(key);properties.extra[key]=error.values[key];if properties.accessors==nil{properties.accessors=map[string]*tsDescriptor{}};properties.accessors[key]=error.descriptors[key]}
 properties.errorStack=error.errorStack;properties.nativeError=true;properties.initialized=true
}
var tsErrors struct{sync.Once;constructors map[string]*tsFunction;prototypes map[string]*tsObject}
func tsInitErrors(){tsErrors.Do(func(){
 tsErrors.constructors=map[string]*tsFunction{};tsErrors.prototypes=map[string]*tsObject{}
 for _,name:=range []string{"Error","TypeError","RangeError","ReferenceError","SyntaxError","URIError","EvalError","AggregateError"}{
  kind:=name;prototype:=tsNewObject();if name!="Error"{prototype.prototype=tsObjectValue(tsErrors.prototypes["Error"])}
  prototype.set("name",tsStringReference(tsStringUTF8(name)));prototype.set("message",tsStringReference(tsStringUTF8("")));prototype.descriptors=map[string]*tsDescriptor{"name":{writable:true,configurable:true},"message":{writable:true,configurable:true}}
  constructor:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNewError(loop,kind,args)});constructor.name=name;constructor.length=1;if name=="AggregateError"{constructor.length=2};constructor.constructible=true;constructor.properties=tsNewObject();constructor.properties.set("prototype",tsObjectValue(prototype));constructor.properties.descriptors=map[string]*tsDescriptor{"prototype":{}}
  if name!="Error"{constructor.internalPrototype=tsFunctionValue(tsErrors.constructors["Error"])}
  prototype.set("constructor",tsFunctionValue(constructor));prototype.descriptors["constructor"]=&tsDescriptor{writable:true,configurable:true};tsErrors.constructors[name]=constructor;tsErrors.prototypes[name]=prototype
 }
 prototype:=tsErrors.prototypes["Error"];prototype.set("toString",tsNamedNativeMethod("toString",0,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{
  if tsIsPrimitiveValue(receiver){tsPropertyFailure("Error.toString requires an object")};name:=tsGet(loop,receiver,"name");message:=tsGet(loop,receiver,"message");n:="Error";m:="";if !tsIsUndefined(name){n=tsStringValue(loop,name).String()};if !tsIsUndefined(message){m=tsStringValue(loop,message).String()};if n==""{return tsStringReference(tsStringUTF8(m))};if m==""{return tsStringReference(tsStringUTF8(n))};return tsStringReference(tsStringUTF8(n+": "+m))
 }));prototype.descriptors["toString"]=&tsDescriptor{writable:true,configurable:true};prototype.order=append(prototype.order,"stack");prototype.descriptors["stack"]=&tsDescriptor{accessor:true,configurable:true,getter:tsNamedNativeMethod("get stack",0,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsErrorStackGet(loop,receiver)}),setter:tsNamedNativeMethod("set stack",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsErrorStackSet(loop,receiver,tsArg(args,0))})}
 })}
func tsErrorConstructor(name string)tsValue{tsInitErrors();return tsFunctionValue(tsErrors.constructors[name])}
func tsErrorPrototype(name string)tsValue{tsInitErrors();if prototype:=tsErrors.prototypes[name];prototype!=nil{return tsObjectValue(prototype)};return tsObjectValue(tsErrors.prototypes["Error"])}
func tsNewError(loop *tsLoop,name string,args []tsValue)tsValue{
 object:=tsNewObject();object.nativeError=true;object.prototype=tsErrorPrototype(name);object.descriptors=map[string]*tsDescriptor{}
 index:=0;if name=="AggregateError"{index=1};message:=tsArg(args,index)
 if !tsIsUndefined(message){object.set("message",tsStringReference(tsStringValue(loop,message)));object.descriptors["message"]=&tsDescriptor{writable:true,configurable:true}}
 options:=tsArg(args,index+1);if !tsIsPrimitiveValue(options)&&tsHas(options,tsStringReference(tsStringUTF8("cause"))){object.set("cause",tsGet(loop,options,"cause"));object.descriptors["cause"]=&tsDescriptor{writable:true,configurable:true}}
 if name=="AggregateError"{values:=&tsArray{};iterator:=tsIterate(tsArg(args,0));for iterator.next(loop){values.values=append(values.values,iterator.value)};object.set("errors",tsArrayValue(values));object.descriptors["errors"]=&tsDescriptor{writable:true,configurable:true}}
 title:=name;if value,ok:=object.values["message"];ok{title+=": "+tsText(value)};object.errorStack=tsCaptureErrorStack(title)
 return tsObjectValue(object)
}
`
