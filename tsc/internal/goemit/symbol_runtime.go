package goemit

// String keys cannot contain 0xff: their encoding is UTF-8/WTF-8. Symbol keys
// occupy a disjoint internal namespace and retain the symbol allocation through
// a string pointing into its key buffer. No conversion of a Symbol to text occurs.
// The registry holds weak references, so ordinary symbols are not immortal.
const symbolRuntime = `
func tsObjectTag(loop *tsLoop,value tsValue)tsValue{tag:="Object";switch value.kind{case tsUndefinedKind:return tsStringReference(tsStringUTF8("[object Undefined]"));case tsNullKind:return tsStringReference(tsStringUTF8("[object Null]"));case tsArrayKind:tag="Array";case tsFunctionKind,tsClassKind:tag="Function";case tsErrorKind:tag="Error";case tsInstanceKind:if (*tsProperties)(value.ref).nativeError{tag="Error"};case tsStringKind:tag="String";case tsBooleanKind:tag="Boolean";case tsSymbolKind:tag="Symbol";case tsBigIntKind:tag="BigInt";case tsObjectKind:object:=(*tsObject)(value.ref);if object.nativeDate{tag="Date"};if object.nativeError{tag="Error"};if object.argumentMap!=nil{tag="Arguments"}};if tsIsNumeric(value){tag="Number"};custom:=tsGet(loop,value,tsWellKnownSymbol("toStringTag"));if custom.kind==tsStringKind{tag=tsText(custom)};return tsStringReference(tsStringUTF8("[object "+tag+"]"))}
type tsSymbol struct{registered bool;registryKey string;description *tsString;id uint64;buffer [9]byte;key string}
var tsSymbolIDs atomic.Uint64
var tsSymbolKeys struct{sync.Mutex;values map[uint64]weak.Pointer[tsSymbol]}
func tsNewSymbol(description *tsString)tsValue{
 symbol:=&tsSymbol{description:description,id:tsSymbolIDs.Add(1)};symbol.buffer[0]=255;binary.LittleEndian.PutUint64(symbol.buffer[1:],symbol.id);symbol.key=unsafe.String(&symbol.buffer[0],len(symbol.buffer))
 tsSymbolKeys.Lock();if tsSymbolKeys.values==nil{tsSymbolKeys.values=map[uint64]weak.Pointer[tsSymbol]{}};tsSymbolKeys.values[symbol.id]=weak.Make(symbol);tsSymbolKeys.Unlock()
 runtime.AddCleanup(symbol,func(id uint64){tsSymbolKeys.Lock();delete(tsSymbolKeys.values,id);tsSymbolKeys.Unlock()},symbol.id)
 return tsValue{kind:tsSymbolKind,ref:unsafe.Pointer(symbol)}
}
func tsIsSymbolKey(key string)bool{return len(key)==9&&key[0]==255}
func tsKeyValue(key string)tsValue{if !tsIsSymbolKey(key){return tsStringReference(tsStringKey(key))};id:=binary.LittleEndian.Uint64([]byte(key[1:]));tsSymbolKeys.Lock();symbol:=tsSymbolKeys.values[id].Value();tsSymbolKeys.Unlock();if symbol==nil{panic("Expired symbol property key")};return tsValue{kind:tsSymbolKind,ref:unsafe.Pointer(symbol)}}
func tsSymbolText(symbol *tsSymbol)string{description:="";if symbol.description!=nil{description=symbol.description.String()};return "Symbol("+description+")"}
type tsSymbolNamespace struct{sync.Once;function *tsFunction;prototype *tsObject;known map[string]tsValue;registryMu sync.Mutex;registry map[string]tsValue}
var tsSymbols tsSymbolNamespace
func tsInitSymbols(){tsSymbols.Do(func(){
 tsSymbols.registry=map[string]tsValue{};tsSymbols.known=map[string]tsValue{};tsSymbols.prototype=tsNewObject()
 function:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{value:=tsArg(args,0);if tsIsUndefined(value){return tsNewSymbol(nil)};return tsNewSymbol(tsStringValue(loop,value))});function.name="Symbol";function.length=0;function.properties=tsNewObject();function.properties.descriptors=map[string]*tsDescriptor{};tsSymbols.function=function
 function.properties.set("prototype",tsObjectValue(tsSymbols.prototype));function.properties.descriptors["prototype"]=&tsDescriptor{}
 function.properties.set("for",tsNamedNativeMethod("for",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{description:=tsStringValue(loop,tsArg(args,0));key:=tsPropertyKey(tsStringReference(description));tsSymbols.registryMu.Lock();defer tsSymbols.registryMu.Unlock();if value,ok:=tsSymbols.registry[key];ok{return value};value:=tsNewSymbol(description);symbol:=(*tsSymbol)(value.ref);symbol.registered=true;symbol.registryKey=key;tsSymbols.registry[key]=value;return value}))
 function.properties.set("keyFor",tsNamedNativeMethod("keyFor",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{value:=tsArg(args,0);if value.kind!=tsSymbolKind{tsPropertyFailure("Symbol.keyFor requires a symbol")};symbol:=(*tsSymbol)(value.ref);if symbol.registered{return tsStringReference(tsStringKey(symbol.registryKey))};return tsU}))
 for _,name:=range []string{"for","keyFor"}{function.properties.descriptors[name]=&tsDescriptor{writable:true,configurable:true}}
 for _,name:=range []string{"iterator","asyncIterator","hasInstance","isConcatSpreadable","match","matchAll","replace","search","species","split","toPrimitive","toStringTag","unscopables","dispose","asyncDispose"}{value:=tsNewSymbol(tsStringUTF8("Symbol."+name));tsSymbols.known[name]=value;function.properties.set(name,value);function.properties.descriptors[name]=&tsDescriptor{}}
 prototype:=tsSymbols.prototype;prototype.set("constructor",tsFunctionValue(function));prototype.descriptors=map[string]*tsDescriptor{"constructor":{writable:true,configurable:true}}
 for _,name:=range []string{"toString","valueOf"}{operation:=name;prototype.set(name,tsNamedNativeMethod(name,0,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{if receiver.kind!=tsSymbolKind{tsPropertyFailure("Invalid Symbol receiver")};if operation=="valueOf"{return receiver};return tsStringReference(tsStringUTF8(tsSymbolText((*tsSymbol)(receiver.ref))))}));prototype.descriptors[name]=&tsDescriptor{writable:true,configurable:true}}
 prototype.order=append(prototype.order,"description");prototype.descriptors["description"]=&tsDescriptor{accessor:true,configurable:true,getter:tsNativeMethod(func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{if receiver.kind!=tsSymbolKind{tsPropertyFailure("Invalid Symbol receiver")};description:=(*tsSymbol)(receiver.ref).description;if description==nil{return tsU};return tsStringReference(description)})}
 key:=(*tsSymbol)(tsSymbols.known["toPrimitive"].ref).key;prototype.set(key,tsNamedNativeMethod("[Symbol.toPrimitive]",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{if receiver.kind!=tsSymbolKind{tsPropertyFailure("Invalid Symbol receiver")};return receiver}));prototype.descriptors[key]=&tsDescriptor{configurable:true}
 key=(*tsSymbol)(tsSymbols.known["toStringTag"].ref).key;prototype.set(key,tsStringReference(tsStringUTF8("Symbol")));prototype.descriptors[key]=&tsDescriptor{configurable:true}
 })}
func tsSymbolFunction()tsValue{tsInitSymbols();return tsFunctionValue(tsSymbols.function)}
func tsWellKnownSymbol(name string)tsValue{tsInitSymbols();return tsSymbols.known[name]}
func tsOwnKeys(value tsValue,symbolsOnly bool)*tsArray{
 var names []string
 switch value.kind{case tsObjectKind:names=(*tsObject)(value.ref).order;case tsFunctionKind:names=(*tsFunction)(value.ref).metadata().order;case tsInstanceKind:names=(*tsProperties)(value.ref).order;case tsClassKind:names=tsInstanceProperties((*tsClass)(value.ref).static).order;case tsArrayKind:array:=(*tsArray)(value.ref);for i:=0;i<array.length();i++{if !array.holes[i]{names=append(names,strconv.Itoa(i))}};names=append(names,"length");for name:=range array.properties{names=append(names,name)};default:tsPropertyFailure("Own keys require an object")}
 out:=&tsArray{};for _,name:=range tsOrderedKeys(names){if !tsIsSymbolKey(name)&&!symbolsOnly{out.values=append(out.values,tsKeyValue(name))}};for _,name:=range names{if tsIsSymbolKey(name){out.values=append(out.values,tsKeyValue(name))}};return out
}
`
