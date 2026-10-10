package goemit

// Math is one stable, lazily initialized namespace. Native numeric storage and
// ordinary class accesses are unaffected by its presence in the runtime.
const mathRuntime = `
func tsMathFround(value float64)float64{return float64(float32(value))}
func tsMathF16round(value float64)float64{if value==0||math.IsNaN(value)||math.IsInf(value,0){return value};magnitude:=math.Abs(value);shift:=-24;if magnitude>=0.00006103515625{_,exponent:=math.Frexp(magnitude);shift=exponent-11};rounded:=math.Ldexp(math.RoundToEven(math.Ldexp(magnitude,-shift)),shift);if rounded>65504{rounded=math.Inf(1)};return math.Copysign(rounded,value)}
func tsMathSumPrecise(loop *tsLoop,source tsValue)float64{
 iterator:=tsIterate(source);sum:=new(big.Int);term:=new(big.Int);negativeZero:=true;nan,positiveInfinity,negativeInfinity:=false,false,false
 for iterator.next(loop){value:=iterator.value;if !tsIsNumeric(value){tsIteratorClose(loop,iterator,true);tsPropertyFailure("Math.sumPrecise requires number elements")};number:=tsNumber(value);if number!=0||!math.Signbit(number){negativeZero=false};if math.IsNaN(number){nan=true;continue};if math.IsInf(number,1){positiveInfinity=true;continue};if math.IsInf(number,-1){negativeInfinity=true;continue};bits:=math.Float64bits(number);exponent:=int((bits>>52)&2047);mantissa:=bits&((1<<52)-1);shift:=0;if exponent!=0{mantissa|=1<<52;shift=exponent-1};term.SetUint64(mantissa);term.Lsh(term,uint(shift));if bits>>63!=0{term.Neg(term)};sum.Add(sum,term)}
 if nan||positiveInfinity&&negativeInfinity{return math.NaN()};if positiveInfinity{return math.Inf(1)};if negativeInfinity{return math.Inf(-1)};if sum.Sign()==0{if negativeZero{return math.Copysign(0,-1)};return 0};value:=new(big.Float).SetPrec(uint(sum.BitLen())).SetInt(sum);value.SetMantExp(value,-1074);number,_:=value.Float64();return number
}
var tsMathNamespace struct{sync.Once;object *tsObject}
func tsMathModule()tsValue{tsMathNamespace.Do(func(){
 object:=tsNewObject();object.descriptors=map[string]*tsDescriptor{}
 for name,value:=range map[string]float64{"E":math.E,"LN10":math.Ln10,"LN2":math.Ln2,"LOG10E":math.Log10E,"LOG2E":math.Log2E,"PI":math.Pi,"SQRT1_2":math.Sqrt(0.5),"SQRT2":math.Sqrt2}{object.set(name,tsNumberValue(value));object.descriptors[name]=&tsDescriptor{}}
 for _,name:=range []string{"abs","acos","acosh","asin","asinh","atan","atanh","atan2","cbrt","ceil","clz32","cos","cosh","exp","expm1","floor","fround","f16round","sumPrecise","hypot","imul","log","log1p","log2","log10","max","min","pow","random","round","sign","sin","sinh","sqrt","tan","tanh","trunc"}{operation:=name;arity:=1;if name=="random"{arity=0};if name=="atan2"||name=="hypot"||name=="imul"||name=="max"||name=="min"||name=="pow"{arity=2};object.set(name,tsNamedNativeMethod(name,arity,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNumberValue(tsMathCall(loop,operation,args))}));object.descriptors[name]=&tsDescriptor{writable:true,configurable:true}}
 tag:=tsPropertyKey(tsWellKnownSymbol("toStringTag"));object.set(tag,tsStringReference(tsStringUTF8("Math")));object.descriptors[tag]=&tsDescriptor{configurable:true}
 tsMathNamespace.object=object
 });return tsObjectValue(tsMathNamespace.object)}
func tsMathRound(x float64)float64{if x==0||math.IsNaN(x)||math.IsInf(x,0)||math.Abs(x)>=4503599627370496{return x};f:=math.Floor(x);if x-f>=0.5{f++};if f==0{return math.Copysign(0,x)};return f}
func tsMathCall(loop *tsLoop,name string,args []tsValue)float64{
 // Variadic operations coerce every argument, even when an earlier result is NaN.
 if name=="sumPrecise"{return tsMathSumPrecise(loop,tsArg(args,0))}
 if name=="f16round"{return tsMathF16round(tsNumber(tsToPrimitive(tsArg(args,0))))}
 if name=="max"||name=="min"{out:=math.Inf(-1);nan:=false;if name=="min"{out=math.Inf(1)};for _,arg:=range args{x:=tsNumber(tsToPrimitive(arg));nan=nan||math.IsNaN(x);if name=="max"{out=math.Max(out,x)}else{out=math.Min(out,x)}};if nan{return math.NaN()};return out}
 if name=="hypot"{out:=0.0;nan:=false;infinite:=false;for _,arg:=range args{x:=tsNumber(tsToPrimitive(arg));nan=nan||math.IsNaN(x);infinite=infinite||math.IsInf(x,0);out=math.Hypot(out,x)};if infinite{return math.Inf(1)};if nan{return math.NaN()};return out}
 if name=="random"{var data [8]byte;if _,err:=cryptorand.Read(data[:]);err!=nil{panic(err)};return float64(binary.LittleEndian.Uint64(data[:])>>11)*(1.0/9007199254740992.0)}
 x:=tsNumber(tsToPrimitive(tsArg(args,0)))
 switch name{
 case "abs":return math.Abs(x);case "acos":return math.Acos(x);case "acosh":return math.Acosh(x);case "asin":return math.Asin(x);case "asinh":return math.Asinh(x);case "atan":return math.Atan(x);case "atanh":return math.Atanh(x);case "cbrt":return math.Cbrt(x);case "ceil":return math.Ceil(x);case "cos":return math.Cos(x);case "cosh":return math.Cosh(x);case "exp":return math.Exp(x);case "expm1":return math.Expm1(x);case "floor":return math.Floor(x);case "fround":return float64(float32(x));case "log":return math.Log(x);case "log1p":return math.Log1p(x);case "log2":return math.Log2(x);case "log10":return math.Log10(x);case "round":return tsMathRound(x);case "sign":if x==0||math.IsNaN(x){return x};return math.Copysign(1,x);case "sin":return math.Sin(x);case "sinh":return math.Sinh(x);case "sqrt":return math.Sqrt(x);case "tan":return math.Tan(x);case "tanh":return math.Tanh(x);case "trunc":return math.Trunc(x)
 case "clz32":v:=uint32(tsNumberUint32(x));n:=0;for n<32&&v&0x80000000==0{n++;v<<=1};return float64(n)
 case "atan2":return math.Atan2(x,tsNumber(tsToPrimitive(tsArg(args,1))))
 case "imul":return float64(int32(uint32(tsNumberUint32(x))*uint32(tsNumberUint32(tsNumber(tsToPrimitive(tsArg(args,1)))))))
 case "pow":y:=tsNumber(tsToPrimitive(tsArg(args,1)));if math.IsNaN(y){return math.NaN()};if y==0{return 1};if math.Abs(x)==1&&math.IsInf(y,0){return math.NaN()};return math.Pow(x,y)
 };panic("Unknown Math operation")
}
`
