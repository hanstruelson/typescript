package goemit

// ValueRuntime owns the representation of JavaScript values. Reference bits are
// never stored in numeric payloads: ref is an actual GC-visible Go pointer.
const ValueRuntime = `
type tsKind uint8
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
 tsIteratorKind
)
type tsValue struct {number float64;ref unsafe.Pointer;kind tsKind}
var tsU=tsValue{}
var tsNull=tsValue{kind:tsNullKind}
func tsNumberValue(value float64)tsValue{return tsValue{number:value,kind:tsNumberKind}}
func tsBooleanValue(value bool)tsValue {number:=float64(0);if value {number=1};return tsValue{number:number,kind:tsBooleanKind}}
func tsStringReference(value *tsString)tsValue{return tsValue{ref:unsafe.Pointer(value),kind:tsStringKind}}
// These host adapters are confined to representation boundaries. JavaScript
// storage and calls use tsValue, never Go interfaces containing primitive values.
func tsValueOf(value any)tsValue {
 switch v:=value.(type) {
 case tsValue:return v
 case nil:return tsNull
 case tsUndefined:return tsU
 case float64:return tsNumberValue(v)
 case bool:return tsBooleanValue(v)
 case string:return tsStringReference(tsStringUTF8(v))
 case *tsString:return tsStringReference(v)
 case *tsObject:return tsValue{ref:unsafe.Pointer(v),kind:tsObjectKind}
 case *tsArray:return tsValue{ref:unsafe.Pointer(v),kind:tsArrayKind}
 case *tsFunction:return tsValue{ref:unsafe.Pointer(v),kind:tsFunctionKind}
 case *tsClass:return tsValue{ref:unsafe.Pointer(v),kind:tsClassKind}
 case tsDynamicObject:return tsValue{ref:unsafe.Pointer(v.tsProperties()),kind:tsInstanceKind}
 case *tsRegExp:return tsValue{ref:unsafe.Pointer(v),kind:tsRegExpKind}
 case *tsECMAObject:return tsValue{ref:unsafe.Pointer(v),kind:tsECMAKind}
 case *tsRuntimeError:return tsValue{ref:unsafe.Pointer(v),kind:tsErrorKind}
 case *tsPromise:return tsValue{ref:unsafe.Pointer(v),kind:tsPromiseKind}
 case *tsTask:return tsValue{ref:unsafe.Pointer(v),kind:tsTaskKind}
 case tsNamespace:return tsValue{ref:unsafe.Pointer(&v),kind:tsNamespaceKind}
 case tsImportRef:return tsValue{ref:unsafe.Pointer(&v),kind:tsImportKind}
 case *tsModule:return tsValue{ref:unsafe.Pointer(v),kind:tsModuleKind}
 case *tsIterator:return tsValue{ref:unsafe.Pointer(v),kind:tsIteratorKind}
 case interface{raw()tsValue}:return v.raw()
 default:panic("Unsupported Go value at JavaScript boundary")
 }
}
func tsUnbox(value any)any {
 if optional,ok:=value.(interface{raw()tsValue});ok {value=optional.raw()}
 v,ok:=value.(tsValue);if !ok {return value}
 switch v.kind {
 case tsUndefinedKind:return tsUndefined{}
 case tsNullKind:return nil
 case tsNumberKind:return v.number
 case tsBooleanKind:return v.number!=0
 case tsStringKind:return (*tsString)(v.ref)
 case tsObjectKind:return (*tsObject)(v.ref)
 case tsArrayKind:return (*tsArray)(v.ref)
 case tsFunctionKind:return (*tsFunction)(v.ref)
 case tsClassKind:return (*tsClass)(v.ref)
 case tsInstanceKind:return (*tsProperties)(v.ref).self
 case tsRegExpKind:return (*tsRegExp)(v.ref)
 case tsECMAKind:return (*tsECMAObject)(v.ref)
 case tsErrorKind:return (*tsRuntimeError)(v.ref)
 case tsPromiseKind:return (*tsPromise)(v.ref)
 case tsTaskKind:return (*tsTask)(v.ref)
 case tsNamespaceKind:return *(*tsNamespace)(v.ref)
 case tsImportKind:return *(*tsImportRef)(v.ref)
 case tsModuleKind:return (*tsModule)(v.ref)
 case tsIteratorKind:return (*tsIterator)(v.ref)
 default:panic("Invalid JavaScript value tag")
 }
}
func tsNative[T any](value tsValue)T {
 var zero T
 // T is one of the compiler's three native primitive representations. The
 // boundary already checked value.kind; this switch selects the storage shape.
 switch any(zero).(type) {
 case float64:return *(*T)(unsafe.Pointer(&value.number))
 case bool:v:=value.number!=0;return *(*T)(unsafe.Pointer(&v))
 case *tsString:v:=(*tsString)(value.ref);return *(*T)(unsafe.Pointer(&v))
 default:panic("Invalid native primitive representation")
 }
}
func tsPrimitiveValue[T any](value T)tsValue {switch v:=any(value).(type){case float64:return tsNumberValue(v);case bool:return tsBooleanValue(v);case *tsString:return tsStringReference(v);default:panic("Invalid native primitive representation")}}
`
