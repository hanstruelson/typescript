package goemit

const TypeRuntime = `
type tsBindingCell struct {raw func() tsValue;initializedState func() bool}
func(c *tsCell)bindingRef()tsBindingCell{return tsBindingCell{raw:c.raw,initializedState:c.initializedState}}
func(c *tsTypedCell[T])bindingRef()tsBindingCell{return tsBindingCell{raw:c.raw,initializedState:c.initializedState}}
func(c *tsCell)raw() tsValue{return c.value}
func(c *tsCell)initializedState() bool{return c.initialized}
// A nullable value has a native payload and a small presence tag. The tag
// distinguishes number, null and undefined without boxing the payload.
type tsOptional[T tsPrimitive] struct {value T;tag uint8}
func tsOptionalFrom[T tsPrimitive](value tsValue)tsOptional[T]{if value.kind==tsNullKind {return tsOptional[T]{tag:1}};if value.kind==tsUndefinedKind{return tsOptional[T]{tag:2}};return tsOptional[T]{value:tsNative[T](value)}}
func(v tsOptional[T])raw()tsValue {if v.tag==1{return tsNull};if v.tag==2{return tsU};return tsPrimitiveValue(v.value)}
func tsNullableNumber(value tsOptional[float64])float64 {switch value.tag {case 1:return 0;case 2:return math.NaN();default:return value.value}}
type tsTypedCell[T tsPrimitive] struct {value T;tag uint8;initialized,constant bool;kind string;nulls uint8;coerce bool}
func tsTypedBinding[T tsPrimitive](initialized,constant bool,kind string,nulls uint8,coerce bool)*tsTypedCell[T]{return &tsTypedCell[T]{tag:2,initialized:initialized,constant:constant,kind:kind,nulls:nulls,coerce:coerce}}
func(c *tsTypedCell[T])raw() tsValue {if c.tag==1{return tsNull};if c.tag==2{return tsU};return tsPrimitiveValue(c.value)}
func(c *tsTypedCell[T])initializedState()bool{return c.initialized}
func(c *tsTypedCell[T])get()tsValue{if !c.initialized {panic("Cannot access binding before initialization")};return c.raw()}
func(c *tsTypedCell[T])read()T{if !c.initialized {panic("Cannot access binding before initialization")};if c.tag!=0 {panic("Non-nullable binding contains a nullish value")};return c.value}
func(c *tsTypedCell[T])optional()tsOptional[T]{if !c.initialized {panic("Cannot access binding before initialization")};return tsOptional[T]{c.value,c.tag}}
func(c *tsTypedCell[T])init(value tsValue)tsValue{c.initialized=true;if value.kind==tsUndefinedKind{c.tag=2;return value};if value.kind==tsNullKind {c.tag=1;return value};checked:=tsBoundary(value,c.kind,c.nulls,c.coerce);c.value=tsNative[T](checked);c.tag=0;return checked}
func(c *tsTypedCell[T])initNative(value T)T{c.value=value;c.tag=0;c.initialized=true;return value}
func(c *tsTypedCell[T])setNative(value T)T{if !c.initialized {panic("Cannot access binding before initialization")};if c.constant {panic("Assignment to constant variable")};return c.initNative(value)}
func(c *tsTypedCell[T])set(value tsValue)tsValue{if !c.initialized{panic("Cannot access binding before initialization")};if c.constant {panic("Assignment to constant variable")};checked:=tsBoundary(value,c.kind,c.nulls,c.coerce);return c.init(checked)}
func(c *tsTypedCell[T])require(){c.init(tsBoundary(c.raw(),c.kind,c.nulls,c.coerce))}
func tsCloneTyped[T tsPrimitive](cell *tsTypedCell[T])*tsTypedCell[T]{copy:=*cell;return &copy}
type tsRuntimeError struct {name,message string}
func(e *tsRuntimeError)String()string{return e.name+": "+e.message}
func tsTypeFailure(kind string){panic(tsThrown{&tsRuntimeError{"TypeError","This value only accepts a "+kind+". Set coerceAny to true to enable automatic conversion."}})}
func tsRestArgs(args []tsValue,index int)*tsArray {if index>=len(args){return &tsArray{}};return &tsArray{values:append([]tsValue{},args[index:]...)}}
func tsTypeOf(value tsValue)*tsString {name:="object";switch value.kind{case tsUndefinedKind:name="undefined";case tsStringKind:name="string";case tsNumberKind,tsFloat32Kind,tsIntKind,tsInt8Kind,tsInt16Kind,tsInt32Kind,tsInt64Kind,tsUintKind,tsUint8Kind,tsUint16Kind,tsUint32Kind,tsUint64Kind:name="number";case tsBooleanKind:name="boolean";case tsFunctionKind,tsClassKind:name="function"};return tsStringUTF8(name)}
func tsBoundary(value tsValue,kind string,nulls uint8,coerce bool) tsValue {
 if tsNumericType(kind) {if value.kind==tsNullKind && nulls&1!=0{return tsNull};if value.kind==tsUndefinedKind && nulls&2!=0{return tsU};return tsNumericBoundary(value,kind,coerce)}
 if value.kind==tsNullKind {if nulls&1!=0 {return tsNull};if coerce {switch kind {case "number":return float64(0);case "string":return tsStringUTF8("null");case "boolean":return false}};tsTypeFailure(kind)}
 if value.kind==tsUndefinedKind{if nulls&2!=0{return tsU};if coerce {switch kind {case "number":return math.NaN();case "string":return tsStringUTF8("undefined");case "boolean":return false}};tsTypeFailure(kind)}
 switch kind {case "number":if value.kind==tsNumberKind{return value};if coerce{return tsNumber(value)}
 case "string":if value.kind==tsStringKind{return value};if coerce{return tsStringValue(value)}
 case "boolean":if value.kind==tsBooleanKind{return value};if coerce{return tsTruthy(value)}}
 tsTypeFailure(kind);return tsU
}
func tsIsSigned(value tsValue)bool {switch value.kind {case tsIntKind,tsInt8Kind,tsInt16Kind,tsInt32Kind,tsInt64Kind:return true};return false}
func tsIsUnsigned(value tsValue)bool {switch value.kind {case tsUintKind,tsUint8Kind,tsUint16Kind,tsUint32Kind,tsUint64Kind:return true};return false}
func tsIsNumeric(value tsValue)bool{return value.kind==tsNumberKind||value.kind==tsFloat32Kind||tsIsSigned(value)||tsIsUnsigned(value)}
func tsNumericType(kind string)bool {switch kind {case "number","float64","float32","int","int8","int16","int32","int64","uint","uint8","uint16","uint32","uint64":return true};return false}
func tsNumericEqual(a,b tsValue)bool {
 if tsIsSigned(a)&&tsIsSigned(b){return int64(math.Float64bits(a.number))==int64(math.Float64bits(b.number))}
 if tsIsUnsigned(a)&&tsIsUnsigned(b){return math.Float64bits(a.number)==math.Float64bits(b.number)}
 if tsIsSigned(a)&&tsIsUnsigned(b){i:=int64(math.Float64bits(a.number));return i>=0&&uint64(i)==math.Float64bits(b.number)}
 if tsIsUnsigned(a)&&tsIsSigned(b){return tsNumericEqual(b,a)}
 if !tsIsSigned(a)&&!tsIsUnsigned(a) {if !tsIsSigned(b)&&!tsIsUnsigned(b){return a.number==b.number};return tsNumericEqual(b,a)}
 f:=b.number;if math.IsNaN(f)||math.IsInf(f,0)||math.Trunc(f)!=f{return false}
 if tsIsSigned(a){return f>=-9223372036854775808.0&&f<9223372036854775808.0&&int64(f)==int64(math.Float64bits(a.number))}
 return f>=0&&f<18446744073709551616.0&&uint64(f)==math.Float64bits(a.number)
}
func tsNumericRangeFailure(kind string){panic(tsThrown{tsErrorValue(&tsRuntimeError{"RangeError","Value cannot be represented as "+kind})})}
func tsNumericBoundary(value tsValue,kind string,coerce bool)tsValue {
 if !tsIsNumeric(value){if !coerce{tsTypeFailure(kind)};value=tsNumberValue(tsNumber(value))}
 if kind=="number"||kind=="float64" {return tsNumberValue(tsNumber(value))}
 if kind=="float32"{f:=tsNumber(value);r:=float32(f);if !math.IsInf(f,0)&&math.IsInf(float64(r),0){tsNumericRangeFailure(kind)};return tsFloat32Value(r)}
 bits:=uint(64);signed:=true
 switch kind {case "int","uint":bits=uint(strconv.IntSize);case "int8","uint8":bits=8;case "int16","uint16":bits=16;case "int32","uint32":bits=32}
 if strings.HasPrefix(kind,"uint"){signed=false}
 var raw uint64
 if tsIsSigned(value){i:=int64(math.Float64bits(value.number));if signed {if bits<64&&(i< -(int64(1)<<(bits-1))||i>(int64(1)<<(bits-1))-1){tsNumericRangeFailure(kind)}}else if i<0||(bits<64&&uint64(i)>uint64(1)<<bits-1){tsNumericRangeFailure(kind)};raw=uint64(i)
 }else if tsIsUnsigned(value){raw=math.Float64bits(value.number);if signed {if raw>uint64(1)<<(bits-1)-1 {tsNumericRangeFailure(kind)}}else if bits<64&&raw>uint64(1)<<bits-1{tsNumericRangeFailure(kind)}
 }else {f:=value.number;if math.IsNaN(f)||math.IsInf(f,0)||math.Trunc(f)!=f{tsNumericRangeFailure(kind)};if signed {bound:=math.Ldexp(1,int(bits)-1);if f< -bound||f>=bound{tsNumericRangeFailure(kind)};raw=uint64(int64(f))}else {if f<0||f>=math.Ldexp(1,int(bits)){tsNumericRangeFailure(kind)};raw=uint64(f)}}
 switch kind {
 case "int":return tsIntValue(int(raw))
 case "int8":return tsInt8Value(int8(raw))
 case "int16":return tsInt16Value(int16(raw))
 case "int32":return tsInt32Value(int32(raw))
 case "int64":return tsInt64Value(int64(raw))
 case "uint":return tsUintValue(uint(raw))
 case "uint8":return tsUint8Value(uint8(raw))
 case "uint16":return tsUint16Value(uint16(raw))
 case "uint32":return tsUint32Value(uint32(raw))
 case "uint64":return tsUint64Value(uint64(raw))
 };panic("Invalid numeric type")
}

type tsNativeInteger interface{int|int8|int16|int32|int64|uint|uint8|uint16|uint32|uint64}
func tsIntegerBinary[T tsNativeInteger](op string,left,right tsValue)tsValue {a,b:=tsNative[T](left),tsNative[T](right);switch op{case "+":return tsPrimitiveValue(a+b);case "-":return tsPrimitiveValue(a-b);case "*":return tsPrimitiveValue(a*b);case "/":return tsPrimitiveValue(a/b);case "%":return tsPrimitiveValue(a%b);case "&":return tsPrimitiveValue(a&b);case "|":return tsPrimitiveValue(a|b);case "^":return tsPrimitiveValue(a^b);case "**":if b<0{return tsNumberValue(tsPow(tsNumber(left),tsNumber(right)))};result:=T(1);for n:=uint64(b);n!=0;n>>=1{if n&1!=0{result*=a};a*=a};return tsPrimitiveValue(result)
 case "<<",">>",">>>":if b<0{tsNumericRangeFailure("nonnegative shift count")};shift:=uint64(b);if op=="<<"{return tsPrimitiveValue(a<<shift)};if op==">>"{return tsPrimitiveValue(a>>shift)};bits:=uint(unsafe.Sizeof(a)*8);raw:=math.Float64bits(left.number);if bits<64{raw&=uint64(1)<<bits-1};return tsPrimitiveValue(T(raw>>shift))};panic("Invalid integer operator")}
func tsNumericCompare(a,b tsValue)(int,bool){
 if (!tsIsSigned(a)&&!tsIsUnsigned(a)&&math.IsNaN(a.number))||(!tsIsSigned(b)&&!tsIsUnsigned(b)&&math.IsNaN(b.number)){return 0,false}
 if !tsIsSigned(a)&&!tsIsUnsigned(a)&&math.IsInf(a.number,0){if !tsIsSigned(b)&&!tsIsUnsigned(b)&&a.number==b.number{return 0,true};if a.number<0{return -1,true};return 1,true}
 if !tsIsSigned(b)&&!tsIsUnsigned(b)&&math.IsInf(b.number,0){if b.number<0{return 1,true};return -1,true}
 toRat:=func(value tsValue)*big.Rat{r:=new(big.Rat);if tsIsSigned(value){return r.SetInt64(int64(math.Float64bits(value.number)))};if tsIsUnsigned(value){return r.SetUint64(math.Float64bits(value.number))};return r.SetFloat64(value.number)}
 return toRat(a).Cmp(toRat(b)),true
}
func tsIntegerBits(kind tsKind,bits uint64)tsValue{switch kind{case tsIntKind:return tsIntValue(int(bits));case tsInt8Kind:return tsInt8Value(int8(bits));case tsInt16Kind:return tsInt16Value(int16(bits));case tsInt32Kind:return tsInt32Value(int32(bits));case tsInt64Kind:return tsInt64Value(int64(bits));case tsUintKind:return tsUintValue(uint(bits));case tsUint8Kind:return tsUint8Value(uint8(bits));case tsUint16Kind:return tsUint16Value(uint16(bits));case tsUint32Kind:return tsUint32Value(uint32(bits));case tsUint64Kind:return tsUint64Value(uint64(bits));};panic("Invalid integer kind")}
`
