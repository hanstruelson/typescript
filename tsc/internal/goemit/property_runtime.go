package goemit

// Metadata is allocated only when descriptors or prototypes are requested.
// Ordinary data properties retain the existing value map and concrete class lenses.
const propertyRuntime = `
func tsObjectLiteralAccessor(loop *tsLoop,object *tsObject,key tsValue,function tsValue,setter bool){name:=tsPropertyKey(key);description:=tsNewObject();description.set("enumerable",tsBooleanValue(true));description.set("configurable",tsBooleanValue(true));if setter{description.set("set",function)}else{description.set("get",function)};tsDefineProperty(loop,tsObjectValue(object),name,tsObjectValue(description))}
type tsDescriptor struct {getter,setter tsValue;accessor,writable,enumerable,configurable bool}
func tsPropertyFailure(message string){panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"TypeError",message:message})})}
func tsObjectRead(loop *tsLoop,receiver tsValue,name string)tsValue{
 object:=(*tsObject)(receiver.ref);cursor:=object
 for cursor!=nil{if descriptor:=cursor.descriptors[name];descriptor!=nil&&descriptor.accessor{if tsIsUndefined(descriptor.getter){return tsU};return tsCallReceiver(loop,descriptor.getter,receiver)};if value,ok:=cursor.values[name];ok{return value};if cursor.nativeClass!=nil{if cursor.nativeClass.builtin=="Stats"||cursor.nativeClass.builtin=="BigIntStats"{switch name{case "atime","mtime","ctime","birthtime":return cursor.get(name)}};cursor=cursor.nativeClass.nativePrototype;continue};if cursor.prototype.kind!=tsObjectKind{break};cursor=(*tsObject)(cursor.prototype.ref)};return tsU
}
func tsObjectWrite(loop *tsLoop,receiver tsValue,name string,value tsValue)tsValue{
 object:=(*tsObject)(receiver.ref);cursor:=object
 for cursor!=nil{if descriptor:=cursor.descriptors[name];descriptor!=nil{if descriptor.accessor{if tsIsUndefined(descriptor.setter){tsPropertyFailure("Property has no setter")};tsCallReceiver(loop,descriptor.setter,receiver,value);return value};if !descriptor.writable{tsPropertyFailure("Cannot assign to read-only property")};break};if _,ok:=cursor.values[name];ok{break};if cursor.nativeClass!=nil{cursor=cursor.nativeClass.nativePrototype;continue};if cursor.prototype.kind!=tsObjectKind{break};cursor=(*tsObject)(cursor.prototype.ref)}
 if _,ok:=object.values[name];!ok&&object.descriptors[name]==nil&&object.nonExtensible{tsPropertyFailure("Cannot add property to a non-extensible object")};return object.set(name,value)
}
func tsOwnObjectProperty(object *tsObject,name string)bool{_,ok:=object.values[name];return ok||object.descriptors[name]!=nil}
func tsDescriptorObject(value tsValue)*tsObject{if value.kind!=tsObjectKind{tsPropertyFailure("Property operation requires an object")};return (*tsObject)(value.ref)}
func tsSameValue(left,right tsValue)bool{if tsIsNumeric(left)&&tsIsNumeric(right){a,b:=tsNumber(left),tsNumber(right);if math.IsNaN(a)&&math.IsNaN(b){return true};if a==0&&b==0{return math.Signbit(a)==math.Signbit(b)}};return tsStrictEqual(left,right)}
func tsNormalizeDescriptor(loop *tsLoop,input tsValue)tsValue{tsDescriptorObject(input);out:=tsNewObject();for _,key:=range []string{"enumerable","configurable","value","writable","get","set"}{if tsHas(input,tsStringReference(tsStringUTF8(key))){value:=tsGet(loop,input,tsStringReference(tsStringUTF8(key)));if key=="enumerable"||key=="configurable"||key=="writable"{value=tsBooleanValue(tsTruthy(value))};if (key=="get"||key=="set")&&!tsIsUndefined(value)&&value.kind!=tsFunctionKind{tsPropertyFailure("Getter or setter must be callable")};out.set(key,value)}};_,get:=out.values["get"];_,set:=out.values["set"];_,value:=out.values["value"];_,writable:=out.values["writable"];if (get||set)&&(value||writable){tsPropertyFailure("Invalid property descriptor: cannot mix accessors and data")};return tsObjectValue(out)}
func tsDefineProperty(loop *tsLoop,target tsValue,name string,input tsValue){
 object:=tsDescriptorObject(target);input=tsNormalizeDescriptor(loop,input);description:=tsDescriptorObject(input)
 has:=func(key string)bool{return tsHas(input,tsStringReference(tsStringUTF8(key)))}
 read:=func(key string)tsValue{return tsGet(loop,input,tsStringReference(tsStringUTF8(key)))}
 accessor:=has("get")||has("set");data:=has("value")||has("writable");if accessor&&data{tsPropertyFailure("Invalid property descriptor: cannot mix accessors and data")}
 old:=object.descriptors[name];exists:=tsOwnObjectProperty(object,name);next:=tsDescriptor{getter:tsU,setter:tsU}
 if exists{next=tsDescriptor{getter:tsU,setter:tsU,writable:true,enumerable:true,configurable:true};if old!=nil{next=*old}}else if object.nonExtensible{tsPropertyFailure("Cannot define property on a non-extensible object")}
 previous:=next
 if accessor||data{if next.accessor!=accessor{next.getter=tsU;next.setter=tsU;next.writable=false};next.accessor=accessor}
 if has("enumerable"){next.enumerable=tsTruthy(read("enumerable"))};if has("configurable"){next.configurable=tsTruthy(read("configurable"))};if has("writable"){next.writable=tsTruthy(read("writable"))}
 for _,key:=range []string{"get","set"}{if has(key){function:=read(key);if !tsIsUndefined(function)&&function.kind!=tsFunctionKind{tsPropertyFailure("Getter or setter must be callable")};if key=="get"{next.getter=function}else{next.setter=function}}}
 if exists&&!previous.configurable{if next.configurable||next.enumerable!=previous.enumerable||(accessor||data)&&next.accessor!=previous.accessor{tsPropertyFailure("Cannot redefine non-configurable property")};if previous.accessor{if !tsSameValue(previous.getter,next.getter)||!tsSameValue(previous.setter,next.setter){tsPropertyFailure("Cannot redefine non-configurable accessor")}}else if !previous.writable{if next.writable||has("value")&&!tsSameValue(object.values[name],read("value")){tsPropertyFailure("Cannot redefine read-only property")}}}
 if !exists{object.order=append(object.order,name)};if object.descriptors==nil{object.descriptors=map[string]*tsDescriptor{}};object.descriptors[name]=&next
 if next.accessor{delete(object.values,name)}else if has("value"){object.values[name]=read("value")}else if _,ok:=object.values[name];!ok{object.values[name]=tsU}
 _=description
}
func tsOwnDescriptor(loop *tsLoop,value tsValue,name string)tsValue{
 object:=tsDescriptorObject(value);if !tsOwnObjectProperty(object,name){return tsU};descriptor:=object.descriptors[name];if descriptor==nil{descriptor=&tsDescriptor{writable:true,enumerable:true,configurable:true}};out:=tsNewObject()
 if descriptor.accessor{out.set("get",descriptor.getter);out.set("set",descriptor.setter)}else{out.set("value",object.values[name]);out.set("writable",tsBooleanValue(descriptor.writable))};out.set("enumerable",tsBooleanValue(descriptor.enumerable));out.set("configurable",tsBooleanValue(descriptor.configurable));return tsObjectValue(out)
}
func tsDescriptorBuiltin(loop *tsLoop,name string,args []tsValue)tsValue{
 value:=tsArg(args,0);if name=="is"{return tsBooleanValue(tsSameValue(value,tsArg(args,1)))}
 if name=="create"{if !tsNullish(value)&&value.kind!=tsObjectKind{tsPropertyFailure("Object prototype must be an object or null")};if tsIsUndefined(value){tsPropertyFailure("Object prototype must be an object or null")};out:=tsNewObject();out.prototype=value;result:=tsObjectValue(out);if !tsIsUndefined(tsArg(args,1)){tsDescriptorBuiltin(loop,"defineProperties",[]tsValue{result,args[1]})};return result}
 if tsIsPrimitiveValue(value){switch name{case "preventExtensions","seal","freeze":return value;case "isExtensible":return tsBooleanValue(false);case "isSealed","isFrozen":return tsBooleanValue(true)}}
 object:=tsDescriptorObject(value)
 switch name{
 case "defineProperty":tsDefineProperty(loop,value,tsPropertyKey(tsArg(args,1)),tsArg(args,2));return value
 case "defineProperties":source:=tsArg(args,1);keys:=tsObjectKeys(source);descriptions:=make([]tsValue,len(keys.values));for i,key:=range keys.values{descriptions[i]=tsNormalizeDescriptor(loop,tsGet(loop,source,key))};for i,key:=range keys.values{tsDefineProperty(loop,value,tsPropertyKey(key),descriptions[i])};return value
 case "getOwnPropertyDescriptor":return tsOwnDescriptor(loop,value,tsPropertyKey(tsArg(args,1)))
 case "getOwnPropertyNames":out:=&tsArray{};for _,key:=range tsOrderedKeys(object.order){out.values=append(out.values,tsStringReference(tsStringKey(key)))};return tsArrayValue(out)
 case "getOwnPropertyDescriptors":out:=tsNewObject();for _,key:=range object.order{out.set(key,tsOwnDescriptor(loop,value,key))};return tsObjectValue(out)
 case "setPrototypeOf":prototype:=tsArg(args,1);if prototype.kind!=tsObjectKind&&prototype.kind!=tsNullKind{tsPropertyFailure("Object prototype must be an object or null")};if tsSameValue(object.prototype,prototype){return value};if object.nonExtensible{tsPropertyFailure("Cannot set prototype of a non-extensible object")};for cursor:=prototype;cursor.kind==tsObjectKind;cursor=(*tsObject)(cursor.ref).prototype{if cursor.ref==value.ref{tsPropertyFailure("Cyclic prototype value")}};object.prototype=prototype;return value
 case "preventExtensions":object.nonExtensible=true;return value
 case "isExtensible":return tsBooleanValue(!object.nonExtensible)
 case "freeze","seal":for _,key:=range object.order{descriptor:=object.descriptors[key];if descriptor==nil{descriptor=&tsDescriptor{writable:true,enumerable:true,configurable:true};if object.descriptors==nil{object.descriptors=map[string]*tsDescriptor{}};object.descriptors[key]=descriptor};descriptor.configurable=false;if name=="freeze"&&!descriptor.accessor{descriptor.writable=false}};object.nonExtensible=true;return value
 case "isFrozen","isSealed":if !object.nonExtensible{return tsBooleanValue(false)};for _,key:=range object.order{descriptor:=object.descriptors[key];if descriptor==nil||descriptor.configurable||name=="isFrozen"&&!descriptor.accessor&&descriptor.writable{return tsBooleanValue(false)}};return tsBooleanValue(true)
 };panic("Unknown descriptor operation")
}
`
