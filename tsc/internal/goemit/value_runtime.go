package goemit

// ValueRuntime owns the representation of JavaScript values. Reference bits are
// never stored in numeric payloads: ref is an actual GC-visible Go pointer.
const ValueRuntime = `
type tsKind uint64
const (
 tsUndefinedKind tsKind=iota
 tsNullKind
 tsNumberKind
 tsBooleanKind
 tsStringKind
 tsObjectKind
 tsArrayKind
 tsFunctionKind
 tsClassKind
 tsInstanceKind
 tsRegExpKind
 tsECMAKind
 tsErrorKind
 tsPromiseKind
 tsTaskKind
 tsNamespaceKind
 tsImportKind
 tsModuleKind
 tsFloat32Kind
 tsIntKind
 tsInt8Kind
 tsInt16Kind
 tsInt32Kind
 tsInt64Kind
 tsUintKind
 tsUint8Kind
 tsUint16Kind
 tsUint32Kind
 tsUint64Kind
 tsIteratorKind
 tsCollectionKind
 tsTypedArrayKind
 tsArrayBufferKind
 tsBigIntKind
 tsSymbolKind
)
type tsValue struct {kind tsKind;number float64;ref unsafe.Pointer}
var tsU=tsValue{}
var tsNull=tsValue{kind:tsNullKind}
func tsNumberValue(value float64)tsValue{return tsValue{number:value,kind:tsNumberKind}}
func tsBooleanValue(value bool)tsValue {number:=float64(0);if value {number=1};return tsValue{number:number,kind:tsBooleanKind}}
func tsStringReference(value *tsString)tsValue{return tsValue{ref:unsafe.Pointer(value),kind:tsStringKind}}
func tsFloat32Value(value float32)tsValue{return tsValue{kind:tsFloat32Kind,number:float64(value)}}
func tsIntValue(value int)tsValue{return tsValue{kind:tsIntKind,number:math.Float64frombits(uint64(value))}}
func tsInt8Value(value int8)tsValue{return tsValue{kind:tsInt8Kind,number:math.Float64frombits(uint64(value))}}
func tsInt16Value(value int16)tsValue{return tsValue{kind:tsInt16Kind,number:math.Float64frombits(uint64(value))}}
func tsInt32Value(value int32)tsValue{return tsValue{kind:tsInt32Kind,number:math.Float64frombits(uint64(value))}}
func tsInt64Value(value int64)tsValue{return tsValue{kind:tsInt64Kind,number:math.Float64frombits(uint64(value))}}
func tsUintValue(value uint)tsValue{return tsValue{kind:tsUintKind,number:math.Float64frombits(uint64(value))}}
func tsUint8Value(value uint8)tsValue{return tsValue{kind:tsUint8Kind,number:math.Float64frombits(uint64(value))}}
func tsUint16Value(value uint16)tsValue{return tsValue{kind:tsUint16Kind,number:math.Float64frombits(uint64(value))}}
func tsUint32Value(value uint32)tsValue{return tsValue{kind:tsUint32Kind,number:math.Float64frombits(uint64(value))}}
func tsUint64Value(value uint64)tsValue{return tsValue{kind:tsUint64Kind,number:math.Float64frombits(uint64(value))}}
// Reference constructors retain GC-visible pointers without Go interface boxes.
func tsObjectValue(value *tsObject)tsValue{return tsValue{ref:unsafe.Pointer(value),kind:tsObjectKind}}
func tsArrayValue(value *tsArray)tsValue{return tsValue{ref:unsafe.Pointer(value),kind:tsArrayKind}}
func tsFunctionValue(value *tsFunction)tsValue{return tsValue{ref:unsafe.Pointer(value),kind:tsFunctionKind}}
func tsClassValue(value *tsClass)tsValue{return tsValue{ref:unsafe.Pointer(value),kind:tsClassKind}}
func tsRegExpValue(value *tsRegExp)tsValue{return tsValue{ref:unsafe.Pointer(value),kind:tsRegExpKind}}
func tsECMAValue(value *tsECMAObject)tsValue{return tsValue{ref:unsafe.Pointer(value),kind:tsECMAKind}}
func tsErrorValue(value *tsRuntimeError)tsValue{if value.stack==nil{value.stack=tsCaptureErrorStack(value.name+": "+value.message)};return tsValue{ref:unsafe.Pointer(value),kind:tsErrorKind}}
func tsPromiseValue(value *tsPromise)tsValue{return tsValue{ref:unsafe.Pointer(value),kind:tsPromiseKind}}
func tsTaskValue(value *tsTask)tsValue{return tsValue{ref:unsafe.Pointer(value),kind:tsTaskKind}}
func tsModuleValue(value *tsModule)tsValue{return tsValue{ref:unsafe.Pointer(value),kind:tsModuleKind}}
func tsIteratorValueReference(value *tsIterator)tsValue{return tsValue{ref:unsafe.Pointer(value),kind:tsIteratorKind}}
func tsInstanceValue(value *tsProperties)tsValue{return tsValue{ref:unsafe.Pointer(value),kind:tsInstanceKind}}
func tsNamespaceValue(value tsNamespace)tsValue{return tsValue{ref:unsafe.Pointer(&value),kind:tsNamespaceKind}}
func tsImportValue(value tsImportRef)tsValue{return tsValue{ref:unsafe.Pointer(&value),kind:tsImportKind}}
func tsClassPointer(value tsValue)*tsClass{if value.kind!=tsClassKind{panic("Expected a class")};return (*tsClass)(value.ref)}
func tsInstanceProperties(value tsValue)*tsProperties{if value.kind!=tsInstanceKind{panic("Expected a class instance")};return (*tsProperties)(value.ref)}
// Generic constraints describe compiler-selected native shapes, not JS values.
type tsPrimitive interface{float64|float32|int|int8|int16|int32|int64|uint|uint8|uint16|uint32|uint64|bool|*tsString}
func tsNative[T tsPrimitive](value tsValue)T {
 var zero T
 // T is a compiler-selected native primitive representation. The
 // boundary already checked value.kind; this switch selects the storage shape.
 switch any(zero).(type) {
 case float64:return *(*T)(unsafe.Pointer(&value.number))
 case float32:v:=float32(value.number);return *(*T)(unsafe.Pointer(&v))
 case int:v:=int(math.Float64bits(value.number));return *(*T)(unsafe.Pointer(&v))
 case int8:v:=int8(math.Float64bits(value.number));return *(*T)(unsafe.Pointer(&v))
 case int16:v:=int16(math.Float64bits(value.number));return *(*T)(unsafe.Pointer(&v))
 case int32:v:=int32(math.Float64bits(value.number));return *(*T)(unsafe.Pointer(&v))
 case int64:v:=int64(math.Float64bits(value.number));return *(*T)(unsafe.Pointer(&v))
 case uint:v:=uint(math.Float64bits(value.number));return *(*T)(unsafe.Pointer(&v))
 case uint8:v:=uint8(math.Float64bits(value.number));return *(*T)(unsafe.Pointer(&v))
 case uint16:v:=uint16(math.Float64bits(value.number));return *(*T)(unsafe.Pointer(&v))
 case uint32:v:=uint32(math.Float64bits(value.number));return *(*T)(unsafe.Pointer(&v))
 case uint64:v:=uint64(math.Float64bits(value.number));return *(*T)(unsafe.Pointer(&v))
 case bool:v:=value.number!=0;return *(*T)(unsafe.Pointer(&v))
 case *tsString:v:=(*tsString)(value.ref);return *(*T)(unsafe.Pointer(&v))
 default:panic("Invalid native primitive representation")
 }
}
func tsPrimitiveValue[T tsPrimitive](value T)tsValue {
 var zero T
 // Only the static generic type is inspected; value is never put in an interface.
 switch any(zero).(type){
 case float64:return tsNumberValue(*(*float64)(unsafe.Pointer(&value)))
 case float32:return tsFloat32Value(*(*float32)(unsafe.Pointer(&value)))
 case int:return tsIntValue(*(*int)(unsafe.Pointer(&value)))
 case int8:return tsInt8Value(*(*int8)(unsafe.Pointer(&value)))
 case int16:return tsInt16Value(*(*int16)(unsafe.Pointer(&value)))
 case int32:return tsInt32Value(*(*int32)(unsafe.Pointer(&value)))
 case int64:return tsInt64Value(*(*int64)(unsafe.Pointer(&value)))
 case uint:return tsUintValue(*(*uint)(unsafe.Pointer(&value)))
 case uint8:return tsUint8Value(*(*uint8)(unsafe.Pointer(&value)))
 case uint16:return tsUint16Value(*(*uint16)(unsafe.Pointer(&value)))
 case uint32:return tsUint32Value(*(*uint32)(unsafe.Pointer(&value)))
 case uint64:return tsUint64Value(*(*uint64)(unsafe.Pointer(&value)))
 case bool:return tsBooleanValue(*(*bool)(unsafe.Pointer(&value)))
 case *tsString:return tsStringReference(*(**tsString)(unsafe.Pointer(&value)))
 default:panic("Invalid native primitive representation")
 }
}
`
