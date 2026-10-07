package goemit

// Arbitrary precision values retain an actual GC-visible *big.Int. Operations
// allocate their result and never mutate an operand shared with another value.
const bigintRuntime = `
func tsBigIntValue(number *big.Int)tsValue{return tsValue{kind:tsBigIntKind,ref:unsafe.Pointer(number)}}
func tsBigIntSigned(number int64)tsValue{return tsBigIntValue(big.NewInt(number))}
func tsBigIntUnsigned(number uint64)tsValue{return tsBigIntValue(new(big.Int).SetUint64(number))}
func tsBigIntParse(text string)(*big.Int,bool){
    text=strings.TrimSpace(strings.Trim(text,"\ufeff"));if text==""{return new(big.Int),true};base:=10
    unsigned:=text;if strings.HasPrefix(unsigned,"+")||strings.HasPrefix(unsigned,"-"){unsigned=unsigned[1:]}
    if len(unsigned)>2&&unsigned[0]=='0'{switch unsigned[1]{case 'x','X','o','O','b','B':if unsigned!=text{return nil,false};base=0}}
    if strings.ContainsAny(text,"._"){return nil,false};number,ok:=new(big.Int).SetString(text,base);return number,ok
}
func tsBigIntFrom(value tsValue)tsValue{
    value=tsToPrimitive(value)
    switch value.kind{
    case tsBigIntKind:return value
    case tsStringKind:number,ok:=tsBigIntParse(tsText(value));if !ok{panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"SyntaxError",message:"Cannot convert string to a BigInt"})})};return tsBigIntValue(number)
    case tsBooleanKind:if value.number!=0{return tsBigIntSigned(1)};return tsBigIntSigned(0)
    }
    if tsIsSigned(value){return tsBigIntSigned(int64(math.Float64bits(value.number)))};if tsIsUnsigned(value){return tsBigIntUnsigned(math.Float64bits(value.number))}
    if tsIsNumeric(value){number:=tsNumber(value);if math.IsNaN(number)||math.IsInf(number,0)||math.Trunc(number)!=number{panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"RangeError",message:"A non-integer number cannot be converted to BigInt"})})};integer,_:=new(big.Float).SetFloat64(number).Int(nil);return tsBigIntValue(integer)}
    panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"TypeError",message:"Cannot convert value to BigInt"})})
}
func tsBigIntLiteral(text string)tsValue{number,ok:=tsBigIntParse(strings.ReplaceAll(strings.TrimSuffix(text,"n"),"_",""));if !ok{panic("Invalid BigInt literal")};return tsBigIntValue(number)}
func tsBigIntCompare(a,b tsValue)(int,bool){
    if a.kind!=tsBigIntKind{comparison,valid:=tsBigIntCompare(b,a);return -comparison,valid}
    integer:=(*big.Int)(a.ref)
    if b.kind==tsBigIntKind{return integer.Cmp((*big.Int)(b.ref)),true}
    if b.kind==tsStringKind{number,ok:=tsBigIntParse(tsText(b));if !ok{return 0,false};return integer.Cmp(number),true}
    if b.kind==tsBooleanKind||b.kind==tsNullKind{b=tsNumberValue(tsNumber(b))}
    if tsIsSigned(b){return integer.Cmp(big.NewInt(int64(math.Float64bits(b.number)))),true};if tsIsUnsigned(b){return integer.Cmp(new(big.Int).SetUint64(math.Float64bits(b.number))),true}
    if !tsIsNumeric(b){return 0,false};number:=tsNumber(b);if math.IsNaN(number){return 0,false};if math.IsInf(number,1){return -1,true};if math.IsInf(number,-1){return 1,true}
    return new(big.Rat).SetInt(integer).Cmp(new(big.Rat).SetFloat64(number)),true
}
func tsBigIntBinary(op string,left,right tsValue)tsValue{
    switch op{
    case "===","!==":equal:=left.kind==tsBigIntKind&&right.kind==tsBigIntKind&&(*big.Int)(left.ref).Cmp((*big.Int)(right.ref))==0;if op=="!=="{equal=!equal};return tsBooleanValue(equal)
    case "==","!=":equal:=tsLooseEqual(left,right);if op=="!="{equal=!equal};return tsBooleanValue(equal)
    case "<","<=",">",">=":left=tsToPrimitive(left);right=tsToPrimitive(right);comparison,valid:=tsBigIntCompare(left,right);if !valid{return tsBooleanValue(false)};switch op{case "<":return tsBooleanValue(comparison<0);case "<=":return tsBooleanValue(comparison<=0);case ">":return tsBooleanValue(comparison>0);default:return tsBooleanValue(comparison>=0)}
    case "+":left=tsToPrimitive(left);right=tsToPrimitive(right);if left.kind==tsStringKind||right.kind==tsStringKind{return tsStringReference(tsStringUTF8(tsText(left)+tsText(right)))}
    }
    if left.kind!=tsBigIntKind||right.kind!=tsBigIntKind{panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"TypeError",message:"Cannot mix BigInt and other types"})})}
    a,b:=(*big.Int)(left.ref),(*big.Int)(right.ref);out:=new(big.Int)
    switch op{
    case "+":out.Add(a,b)
    case "-":out.Sub(a,b)
    case "*":out.Mul(a,b)
    case "/","%":if b.Sign()==0{panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"RangeError",message:"Division by zero"})})};if op=="/"{out.Quo(a,b)}else{out.Rem(a,b)}
    case "&":out.And(a,b)
    case "|":out.Or(a,b)
    case "^":out.Xor(a,b)
    case ">>>":panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"TypeError",message:"BigInts have no unsigned right shift"})})
    case "<<",">>":shift:=new(big.Int).Abs(b);rightShift:=(op==">>")!=(b.Sign()<0);if rightShift&&(!shift.IsUint64()||shift.Uint64()>uint64(a.BitLen()+1)){if a.Sign()<0{return tsBigIntSigned(-1)};return tsBigIntSigned(0)};if !shift.IsUint64()||shift.Uint64()>1<<30{panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"RangeError",message:"Maximum BigInt size exceeded"})})};if rightShift{out.Rsh(a,uint(shift.Uint64()))}else{out.Lsh(a,uint(shift.Uint64()))}
    case "**":if b.Sign()<0{panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"RangeError",message:"Exponent must be positive"})})};if !b.IsUint64()||(a.BitLen()>1&&b.Uint64()>uint64(1<<30)/uint64(a.BitLen())){panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"RangeError",message:"Maximum BigInt size exceeded"})})};out.Exp(a,b,nil)
    default:panic("Unsupported BigInt operator")
    };return tsBigIntValue(out)
}
func tsNegate(value tsValue)tsValue{if value.kind==tsBigIntKind{return tsBigIntValue(new(big.Int).Neg((*big.Int)(value.ref)))};return tsNumberValue(-tsNumber(value))}
func tsIncrementOperand(value tsValue)tsValue{value=tsToPrimitive(value);if value.kind==tsBigIntKind{return value};return tsNumberValue(tsNumber(value))}
func tsIncrement(value tsValue,op string)tsValue{if value.kind==tsBigIntKind{return tsBigIntBinary(op,value,tsBigIntSigned(1))};return tsBinary(op,value,tsNumberValue(1))}
func tsBigIntProperty(value tsValue,name string)tsValue{
    switch name{
    case "valueOf":return tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return value}))
    case "toString","toLocaleString":return tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{radix:=10;if name=="toString"&&!tsIsUndefined(tsArg(args,0)){radix=int(tsNumber(tsArg(args,0)));if radix<2||radix>36{panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"RangeError",message:"radix must be between 2 and 36"})})}};return tsStringReference(tsStringUTF8((*big.Int)(value.ref).Text(radix)))}))
    };return tsU
}
func tsBigIntFunction()tsValue{
    function:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsBigIntFrom(tsArg(args,0))});function.properties=tsNewObject()
    for _,method:=range []string{"asIntN","asUintN"}{name:=method;function.properties.set(name,tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{width:=tsToIndex(tsArg(args,0));value:=tsArg(args,1);if tsIsNumeric(value){panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"TypeError",message:"Cannot convert number to BigInt"})})};number:=(*big.Int)(tsBigIntFrom(value).ref);if width==0{return tsBigIntSigned(0)};if width>1<<30{panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"RangeError",message:"Maximum BigInt size exceeded"})})};modulus:=new(big.Int).Lsh(big.NewInt(1),uint(width));result:=new(big.Int).Mod(number,modulus);if name=="asIntN"&&result.Bit(width-1)!=0{result.Sub(result,modulus)};return tsBigIntValue(result)})))}
    return tsFunctionValue(function)
}
func tsConvertBigInt(value tsValue,kind string)tsValue{
    integer:=(*big.Int)(value.ref)
    switch kind{
    case "string":return tsStringReference(tsStringUTF8(integer.String()))
    case "boolean":return tsBooleanValue(integer.Sign()!=0)
    case "number","float32":number,_:=new(big.Float).SetInt(integer).Float64();if math.IsInf(number,0){tsNumericRangeFailure(kind)};return tsConvertValue(tsNumberValue(number),kind,true)
    }
    if strings.HasPrefix(kind,"uint"){if !integer.IsUint64(){tsNumericRangeFailure(kind)};return tsConvertValue(tsUint64Value(integer.Uint64()),kind,true)}
    if !integer.IsInt64(){tsNumericRangeFailure(kind)};return tsConvertValue(tsInt64Value(integer.Int64()),kind,true)
}
`
