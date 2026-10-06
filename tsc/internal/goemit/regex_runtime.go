package goemit

// The ECMAScript engine supplies RegExp and ECMAScript numeric conversions.
// No application source or eval/Function input is compiled or executed here.
const RegexRuntime = `
var tsECMA *goja.Runtime
func tsEngine()*goja.Runtime {if tsECMA==nil {tsECMA=goja.New()};return tsECMA}
type tsRegExp struct {object *goja.Object}
type tsECMAObject struct {object *goja.Object}
func tsToECMA(value tsValue) goja.Value {switch v:=value.(type){case tsUndefined:return goja.Undefined();case nil:return goja.Null();case *tsString:return goja.StringFromUTF16(v.units);case *tsRegExp:return v.object;case *tsECMAObject:return v.object;case *tsObject:object:=tsEngine().NewObject();for _,key:=range v.order {if err:=object.Set(key,tsToECMA(v.values[key]));err!=nil {tsECMAFailure(err)}};return object;case *tsArray:values:=make([]interface{},len(v.values));for i,item:=range v.values {values[i]=tsToECMA(item)};return tsEngine().NewArray(values...);case *tsFunction:return tsEngine().ToValue(func(call goja.FunctionCall)goja.Value{args:=make([]tsValue,len(call.Arguments));for i,arg:=range call.Arguments {args[i]=tsFromECMA(arg)};return tsToECMA(tsCall(v,args...))});default:return tsEngine().ToValue(v)}}
func tsFromECMA(value goja.Value) tsValue {if goja.IsUndefined(value){return tsU};if goja.IsNull(value){return nil};if text,ok:=value.(goja.String);ok {units:=make([]uint16,text.Length());for i:=range units {units[i]=text.CharAt(i)};return tsStringUnits(units)};if object,ok:=value.(*goja.Object);ok {if object.ClassName()=="RegExp" {return &tsRegExp{object}};return &tsECMAObject{object}};switch v:=value.Export().(type){case int64:return float64(v);case int32:return float64(v);case int:return float64(v);default:return v}}
func tsECMAFailure(err error) {if exception,ok:=err.(*goja.Exception);ok {panic(tsThrown{tsFromECMA(exception.Value())})};panic(err.Error())}
func tsNewRegExp(args ...tsValue) tsValue {arguments:=make([]goja.Value,len(args));for i,value:=range args {arguments[i]=tsToECMA(value)};constructor,_:=goja.AssertConstructor(tsEngine().Get("RegExp"));object,err:=constructor(nil,arguments...);if err!=nil {tsECMAFailure(err)};return &tsRegExp{object}}
func tsECMAGet(object *goja.Object,key tsValue) tsValue {name:=tsText(key);value:=object.Get(name);if call,ok:=goja.AssertFunction(value);ok {return tsFunc(func(args ...tsValue)tsValue{arguments:=make([]goja.Value,len(args));for i,arg:=range args {arguments[i]=tsToECMA(arg)};result,err:=call(object,arguments...);if err!=nil {tsECMAFailure(err)};return tsFromECMA(result)})};return tsFromECMA(value)}
func tsECMASet(object *goja.Object,key,value tsValue) tsValue {if err:=object.Set(tsText(key),tsToECMA(value));err!=nil {tsECMAFailure(err)};return value}
func tsRegexString(text *tsString,name string,args []tsValue) tsValue {prototype:=tsEngine().Get("String").ToObject(tsEngine()).Get("prototype").ToObject(tsEngine());method,ok:=goja.AssertFunction(prototype.Get(name));if !ok {panic("ECMAScript String method unavailable: "+name)};arguments:=make([]goja.Value,len(args));for i,arg:=range args {arguments[i]=tsToECMA(arg)};result,err:=method(goja.StringFromUTF16(text.units),arguments...);if err!=nil {tsECMAFailure(err)};return tsFromECMA(result)}
`
