package goemit

const TypeRuntime = `
type tsBindingCell interface {get() tsValue;raw() tsValue;initializedState() bool}
func(c *tsCell)raw() tsValue{return c.value}
func(c *tsCell)initializedState() bool{return c.initialized}
// A nullable value has a native payload and a small presence tag. The tag
// distinguishes number, null and undefined without boxing the payload.
type tsOptional[T any] struct {value T;tag uint8}
func tsOptionalFrom[T any](value tsValue)tsOptional[T]{if value==nil {return tsOptional[T]{tag:1}};if tsIsUndefined(value){return tsOptional[T]{tag:2}};return tsOptional[T]{value:value.(T)}}
func(v tsOptional[T])raw()tsValue {if v.tag==1{return nil};if v.tag==2{return tsU};return v.value}
func tsNullableNumber(value tsOptional[float64])float64 {switch value.tag {case 1:return 0;case 2:return math.NaN();default:return value.value}}
type tsTypedCell[T any] struct {value T;tag uint8;initialized,constant bool;kind string;nulls uint8;coerce bool}
func tsTypedBinding[T any](initialized,constant bool,kind string,nulls uint8,coerce bool)*tsTypedCell[T]{return &tsTypedCell[T]{tag:2,initialized:initialized,constant:constant,kind:kind,nulls:nulls,coerce:coerce}}
func(c *tsTypedCell[T])raw() tsValue {if c.tag==1{return nil};if c.tag==2{return tsU};return c.value}
func(c *tsTypedCell[T])initializedState()bool{return c.initialized}
func(c *tsTypedCell[T])get()tsValue{if !c.initialized {panic("Cannot access binding before initialization")};return c.raw()}
func(c *tsTypedCell[T])read()T{if !c.initialized {panic("Cannot access binding before initialization")};if c.tag!=0 {panic("Non-nullable binding contains a nullish value")};return c.value}
func(c *tsTypedCell[T])optional()tsOptional[T]{if !c.initialized {panic("Cannot access binding before initialization")};return tsOptional[T]{c.value,c.tag}}
func(c *tsTypedCell[T])init(value tsValue)tsValue{c.initialized=true;if tsIsUndefined(value){c.tag=2;return value};if value==nil {c.tag=1;return value};checked:=tsBoundary(value,c.kind,c.nulls,c.coerce);c.value=checked.(T);c.tag=0;return checked}
func(c *tsTypedCell[T])initNative(value T)T{c.value=value;c.tag=0;c.initialized=true;return value}
func(c *tsTypedCell[T])setNative(value T)T{if !c.initialized {panic("Cannot access binding before initialization")};if c.constant {panic("Assignment to constant variable")};return c.initNative(value)}
func(c *tsTypedCell[T])set(value tsValue)tsValue{if !c.initialized{panic("Cannot access binding before initialization")};if c.constant {panic("Assignment to constant variable")};checked:=tsBoundary(value,c.kind,c.nulls,c.coerce);return c.init(checked)}
func(c *tsTypedCell[T])require(){c.init(tsBoundary(c.raw(),c.kind,c.nulls,c.coerce))}
func tsCloneTyped[T any](cell *tsTypedCell[T])*tsTypedCell[T]{copy:=*cell;return &copy}
type tsRuntimeError struct {name,message string}
func(e *tsRuntimeError)String()string{return e.name+": "+e.message}
func tsTypeFailure(kind string){panic(tsThrown{&tsRuntimeError{"TypeError","This value only accepts a "+kind+". Set coerceAny to true to enable automatic conversion."}})}
func tsRestArgs(args []tsValue,index int)*tsArray {if index>=len(args){return &tsArray{}};return &tsArray{values:append([]tsValue{},args[index:]...)}}
func tsTypeOf(value tsValue)*tsString {name:="object";switch value.(type){case tsUndefined:name="undefined";case *tsString:name="string";case float64:name="number";case bool:name="boolean";case *tsFunction,*tsClass:name="function"};return tsStringUTF8(name)}
func tsBoundary(value tsValue,kind string,nulls uint8,coerce bool) tsValue {
 if value==nil {if nulls&1!=0 {return nil};if coerce {switch kind {case "number":return float64(0);case "string":return tsStringUTF8("null");case "boolean":return false}};tsTypeFailure(kind)}
 if tsIsUndefined(value){if nulls&2!=0{return tsU};if coerce {switch kind {case "number":return math.NaN();case "string":return tsStringUTF8("undefined");case "boolean":return false}};tsTypeFailure(kind)}
 switch kind {case "number":if _,ok:=value.(float64);ok{return value};if coerce{return tsNumber(value)}
 case "string":if _,ok:=value.(*tsString);ok{return value};if coerce{return tsStringValue(value)}
 case "boolean":if _,ok:=value.(bool);ok{return value};if coerce{return tsTruthy(value)}}
 tsTypeFailure(kind);return tsU
}
`
