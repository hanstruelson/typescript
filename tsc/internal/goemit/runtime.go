package goemit

//go:generate go run ../../cmd/generate-runtime

// Runtime is emitted into the single output file; its dependencies use the compiler module.
// Each event loop owns its task registry and runs callbacks on one goroutine.
// User workers have their own loop and share application references explicitly.
const Runtime = `import (
 "context"
 "unsafe"
 "weak"
 "bytes"
	"bufio"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	zstd "github.com/klauspost/compress/zstd"
 "crypto"
 "crypto/aes"
 "crypto/cipher"
 "crypto/des"
 "crypto/ecdh"
 "crypto/ecdsa"
 "crypto/ed25519"
 "crypto/elliptic"
 "crypto/hmac"
 "crypto/md5"
 cryptorand "crypto/rand"
 "crypto/rsa"
 "crypto/sha1"
 "crypto/sha256"
 "crypto/sha512"
 "crypto/x509"
 cryptotls "crypto/tls"
 "encoding/base64"
 "encoding/binary"
 "encoding/hex"
 "encoding/json"
 "encoding/pem"
 "errors"
 "io"
 "net"
 "net/url"
 "path"
 "path/filepath"
 "reflect"
 "runtime"
 "syscall"
 "golang.org/x/sys/unix"
 "sync"
 "sync/atomic"
 "fmt"
 "hash"
	"hash/crc32"
 "math"
 "math/big"
 "os"
	exec "os/exec"
 gosignal "os/signal"
 osuser "os/user"
 "strings"
 "sort"
 "time"
 "unicode"
 "unicode/utf16"
 "unicode/utf8"
 "strconv"
 "golang.org/x/text/cases"
 "golang.org/x/text/language"
 "golang.org/x/text/unicode/norm"
 "golang.org/x/text/collate"
)

type tsCell struct { value tsValue; initialized, constant bool;arrayKind string;arrayNulls uint8;arrayCoerce bool }
func tsBinding(initialized, constant bool) *tsCell { return &tsCell{value:tsU,initialized:initialized,constant:constant} }
func (c *tsCell) get() tsValue {if !c.initialized {panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"ReferenceError",message:"Cannot access binding before initialization"})})};if c.value.kind==tsImportKind{ref:=(*tsImportRef)(c.value.ref);return ref.module.read(ref.name)};return c.value}
func (c *tsCell) init(value tsValue) tsValue {if c.arrayKind!=""{value=tsValue(tsArrayBoundary(value,c.arrayKind,c.arrayNulls,c.arrayCoerce))};c.value = value; c.initialized = true; return value }
func (c *tsCell) set(value tsValue) tsValue { if !c.initialized { panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"ReferenceError",message:"Cannot access binding before initialization"})}) }; if c.constant { panic("Assignment to constant variable") }; if c.arrayKind!=""{value=tsValue(tsArrayBoundary(value,c.arrayKind,c.arrayNulls,c.arrayCoerce))};c.value = value; return value }
func tsClone(c *tsCell) *tsCell { value := *c; return &value }
type tsFunction struct {generator,asyncGeneratorFunction bool;name string;length int;constructible bool;prototypeOnce,metadataOnce sync.Once;metadataRef atomic.Pointer[tsObject];prototype,internalPrototype,boundTarget,boundReceiver tsValue;boundArgs []tsValue;receiverCall func(*tsLoop,tsValue,...tsValue)tsValue;construct func(*tsLoop,...tsValue)tsValue;properties *tsObject;call func(*tsLoop,...tsValue) tsValue}
func tsFunc(call func(*tsLoop,...tsValue)tsValue)*tsFunction{return &tsFunction{call:call}}
func tsCall(loop *tsLoop,value tsValue,args ...tsValue)tsValue{if value.kind==tsClassKind{class:=(*tsClass)(value.ref);switch class.builtin{case "Date":return tsStringReference(tsStringUTF8(tsDateHuman(time.Now().In(tsDateLocal()),"toString")));case "Array","Buffer","ReadStream","WriteStream":return tsConstruct(value,args...);case "Stats","BigIntStats","StatFs","BigIntStatFs":return class.construct(args...);case "Dirent":panic(tsThrown{tsFSException("TypeError","ERR_ILLEGAL_CONSTRUCTOR","Dirent constructor cannot be called as a function")})}};if value.kind!=tsFunctionKind{tsPropertyFailure("Value is not callable")};function:=(*tsFunction)(value.ref);if function.receiverCall!=nil{return function.receiverCall(loop,tsU,args...)};return function.call(loop,args...)}
func tsArg(args []tsValue, index int) tsValue { if index < len(args) { return args[index] }; return tsU }
// ECMAScript StringNumericLiteral conversion, independent of the regexp engine.
func tsNumericStringText(value *tsString)string{return strings.TrimFunc(value.String(),func(r rune)bool{return r=='\ufeff'||r=='\u0009'||r=='\u000b'||r=='\u000c'||r=='\u0020'||r=='\u00a0'||r=='\u1680'||r>=0x2000&&r<=0x200a||r==0x202f||r==0x205f||r==0x3000||r=='\n'||r=='\r'||r==0x2028||r==0x2029})}
func tsStringNumber(value *tsString)float64{
 text:=tsNumericStringText(value)
 if text==""{return 0};if text=="Infinity"||text=="+Infinity"{return math.Inf(1)};if text=="-Infinity"{return math.Inf(-1)}
 if len(text)>2&&text[0]=='0'{base:=0;switch text[1]{case 'x','X':base=16;case 'o','O':base=8;case 'b','B':base=2};if base!=0{digits:=text[2:];for _,c:=range digits{if !(c>='0'&&c<='9'&&int(c-'0')<base||base==16&&(c>='a'&&c<='f'||c>='A'&&c<='F')){return math.NaN()}};integer,ok:=new(big.Int).SetString(digits,base);if !ok{return math.NaN()};number,_:=new(big.Float).SetInt(integer).Float64();return number}}
 // ParseFloat also accepts Go hex floats, underscores and Inf; reject those.
 for _,c:=range text{if !(c>='0'&&c<='9'||c=='.'||c=='+'||c=='-'||c=='e'||c=='E'){return math.NaN()}}
 number,err:=strconv.ParseFloat(text,64);if err!=nil{if !errors.Is(err,strconv.ErrRange){return math.NaN()}};return number
}
func tsNumberText(number float64)string{
 if math.IsNaN(number){return "NaN"};if math.IsInf(number,1){return "Infinity"};if math.IsInf(number,-1){return "-Infinity"};if number==0{return "0"}
 magnitude:=math.Abs(number);if magnitude>=1e-6&&magnitude<1e21{return strconv.FormatFloat(number,'f',-1,64)}
 text:=strconv.FormatFloat(number,'e',-1,64);index:=strings.LastIndexByte(text,'e');exponent,_:=strconv.Atoi(text[index+1:]);sign:="";if exponent>=0{sign="+"};return text[:index+1]+sign+strconv.Itoa(exponent)
}
func tsText(value tsValue)string{switch value.kind{
 case tsSymbolKind:return tsSymbolText((*tsSymbol)(value.ref))
 case tsBigIntKind:return (*big.Int)(value.ref).String()
 case tsUndefinedKind:return "undefined"
 case tsNullKind:return "null"
 case tsIntKind,tsInt8Kind,tsInt16Kind,tsInt32Kind,tsInt64Kind:return strconv.FormatInt(int64(math.Float64bits(value.number)),10)
 case tsUintKind,tsUint8Kind,tsUint16Kind,tsUint32Kind,tsUint64Kind:return strconv.FormatUint(math.Float64bits(value.number),10)
 case tsFloat32Kind:return strconv.FormatFloat(value.number,'g',-1,32)
 case tsNumberKind:return tsNumberText(value.number)
 case tsBooleanKind:return strconv.FormatBool(value.number!=0)
 case tsStringKind:return (*tsString)(value.ref).String()
 case tsErrorKind:return (*tsRuntimeError)(value.ref).String()
 case tsObjectKind:return fmt.Sprint((*tsObject)(value.ref))
 case tsArrayKind:return fmt.Sprint((*tsArray)(value.ref))
 case tsFunctionKind:return "function "+(*tsFunction)(value.ref).name+"() { [native code] }"
 case tsClassKind:return fmt.Sprint((*tsClass)(value.ref))
 case tsInstanceKind:return "[object Object]"
 case tsRegExpKind:return tsRegExpText((*tsRegExp)(value.ref)).String()

 case tsPromiseKind:return fmt.Sprint((*tsPromise)(value.ref))
 case tsTaskKind:return fmt.Sprint((*tsTask)(value.ref))
 case tsNamespaceKind:return fmt.Sprint(*(*tsNamespace)(value.ref))
 case tsImportKind:return fmt.Sprint(*(*tsImportRef)(value.ref))
 case tsModuleKind:return fmt.Sprint((*tsModule)(value.ref))
 case tsArrayBufferKind:return "[object ArrayBuffer]"
 case tsDataViewKind:return "[object DataView]"
 case tsTypedArrayKind:a:=(*tsTypedArray)(value.ref);parts:=make([]string,a.length);for i:=range parts{if a.bigint{parts[i]=tsText(a.readValue(i))}else{parts[i]=tsNumberText(a.read(i))}};return strings.Join(parts,",")
 case tsCollectionKind:if (*tsCollection)(value.ref).set{return "[object Set]"};return "[object Map]"
 case tsIteratorKind:return fmt.Sprint((*tsIterator)(value.ref))
 default:panic("Invalid JavaScript value tag")}}
func(value tsValue)String()string{return tsText(value)}
func tsConsoleLog(values ...tsValue) tsValue { for i, value := range values { if i != 0 { fmt.Print(" ") }; fmt.Print(tsText(value));if value.kind==tsBigIntKind{fmt.Print("n")} }; fmt.Println(); return tsU }
func tsIsUndefined(value tsValue) bool {return value.kind==tsUndefinedKind}
func tsNullish(value tsValue) bool {return value.kind==tsNullKind || value.kind==tsUndefinedKind}
type tsArray struct {prototype tsValue;descriptors map[string]*tsDescriptor;order []string;native unsafe.Pointer;elementKind string;elementNulls uint8;elementCoerce bool;growLength func()int;growGet func(int)tsValue;growClear func(int);growPut func(*tsLoop,int,tsValue);growResize func(int);growSplice func(*tsLoop,int,int,[]tsValue);values []tsValue;properties map[string]tsValue;holes map[int]bool }
type tsIterator struct {prototype tsValue;protocol tsValue;nextMethod tsValue;machine *tsMachine;pull func()(tsValue,bool);array *tsArray; text *tsString; collection *tsCollection;collectionIndex int;collectionMode string;index int;value tsValue;regexpString bool;done,collectionIterator,collectionSet bool}
func tsIterate(value tsValue)*tsIterator{if value.kind==tsObjectKind||value.kind==tsInstanceKind{method:=tsGet(loop,value,tsValue(tsWellKnownSymbol("iterator")));if !tsNullish(method){iterator:=tsCallReceiver(loop,method,value);if tsIsPrimitiveValue(iterator){tsPropertyFailure("Iterator must be an object")};if iterator.kind==tsIteratorKind{return (*tsIterator)(iterator.ref)};return &tsIterator{protocol:iterator,nextMethod:tsGet(loop,iterator,"next")}}};switch value.kind{case tsTypedArrayKind:array:=(*tsTypedArray)(value.ref);index:=0;return &tsIterator{pull:func()(tsValue,bool){if index>=array.length{return tsU,false};item:=tsNumberValue(array.read(index));if array.bigint{item=array.readValue(index)};index++;return item,true}};case tsCollectionKind:c:=(*tsCollection)(value.ref);mode:="entries";if c.set{mode="values"};return c.iterator(mode);case tsIteratorKind:return (*tsIterator)(value.ref);case tsArrayKind:return &tsIterator{array:(*tsArray)(value.ref)};case tsStringKind:return &tsIterator{text:(*tsString)(value.ref)};default:tsPropertyFailure("Value is not iterable");return nil}}
func (it *tsIterator) next() bool {if !tsIsUndefined(it.protocol){result:=tsCallReceiver(loop,it.nextMethod,it.protocol);if tsIsPrimitiveValue(result){tsPropertyFailure("Iterator result must be an object")};if tsTruthy(tsGet(loop,result,"done")){it.value=tsU;return false};it.value=tsGet(loop,result,"value");return true};if it.machine!=nil{result:=tsGeneratorStep(loop,it,"next",tsU);if tsTruthy(tsGet(loop,result,"done")){it.value=tsU;return false};it.value=tsGet(loop,result,"value");return true};if it.pull!=nil{value,ok:=it.pull();it.value=value;return ok};if it.array!=nil {if it.index>=it.array.length(){return false};it.value=it.array.at(it.index);it.index++;return true};if it.index>=len(it.text.units){return false};start:=it.index;it.index++;c:=it.text.units[start];if c>=0xd800&&c<=0xdbff&&it.index<len(it.text.units)&&it.text.units[it.index]>=0xdc00&&it.text.units[it.index]<=0xdfff {it.index++};it.value=it.text.slice(start,it.index);return true}
func tsGet(value,key tsValue) tsValue {if tsNullish(value){tsPropertyFailure("Cannot read properties of null or undefined")};if key.kind==tsSymbolKind{if property:=tsIteratorSymbolProperty(loop,value,key);!tsIsUndefined(property){return property}};if value.kind==tsBigIntKind&&tsText(key)=="toString"{return tsBigIntProperty(value,"toString")};if key.kind==tsStringKind&&tsIsPrimitiveValue(value)&&value.kind!=tsSymbolKind{if method:=tsConversionProperty(value,tsPropertyKey(key));!tsIsUndefined(method){return method}};switch value.kind {
 case tsBigIntKind:return tsBigIntProperty(value,tsPropertyKey(key))
 case tsFunctionKind:return tsFunctionGet(loop,value,tsPropertyKey(key))
 case tsArrayBufferKind:return tsArrayBufferGet(value,tsPropertyKey(key))
 case tsDataViewKind:return tsDataViewGet(value,tsPropertyKey(key))
 case tsTypedArrayKind:return tsValue(tsTypedArrayGet(value,tsPropertyKey(key)))
 case tsCollectionKind:return tsCollectionGet(value,tsPropertyKey(key))
 case tsIteratorKind:return tsValue(tsIteratorGet((*tsIterator)(value.ref),tsPropertyKey(key)))
 case tsClassKind:v:=(*tsClass)(value.ref);if tsPropertyKey(key)=="prototype"{return tsValue(tsClassPrototype(v))};if v.builtin!=""{if tsPropertyKey(key)=="name"{return tsStringReference(tsStringUTF8(v.builtin))};if tsPropertyKey(key)=="length"{switch v.builtin{case "Stats","BigIntStats","StatFs","BigIntStatFs":return tsNumberValue(0);case "Dirent":return tsNumberValue(3)}};if tsPropertyKey(key)=="BYTES_PER_ELEMENT"&&(v.nativePrototype==nil||v.builtin=="Buffer")&&v.builtin!="Map"&&v.builtin!="Set"&&v.builtin!="Array"{size:=1;switch v.builtin{case "Float64Array","BigInt64Array","BigUint64Array":size=8;case "Float32Array","Int32Array","Uint32Array":size=4;case "Float16Array","Int16Array","Uint16Array":size=2};return tsNumberValue(float64(size))}};return tsClassRead(loop,value,tsPropertyKey(key))
 case tsSymbolKind:tsInitSymbols();return tsPrototypeLookup(loop,value,tsObjectValue(tsSymbols.prototype),tsPropertyKey(key))
 case tsInstanceKind:return tsInstanceRead(loop,value,tsPropertyKey(key))
 case tsObjectKind:return tsObjectRead(loop,value,tsPropertyKey(key))
 case tsArrayKind:return tsArrayGet(loop,value,key)
 case tsStringKind:v:=(*tsString)(value.ref);if property:=tsStringGet(v,key);!tsIsUndefined(property){return property};return tsPrototypeLookup(loop,value,tsValue(tsObjectPrototype()),tsPropertyKey(key))
 case tsBooleanKind:return tsPrototypeLookup(loop,value,tsValue(tsObjectPrototype()),tsPropertyKey(key))
 case tsRegExpKind:return tsRegExpGet(loop,value,tsPropertyKey(key))
 case tsNumberKind,tsFloat32Kind,tsIntKind,tsInt8Kind,tsInt16Kind,tsInt32Kind,tsInt64Kind,tsUintKind,tsUint8Kind,tsUint16Kind,tsUint32Kind,tsUint64Kind:return tsPrototypeLookup(loop,value,tsObjectValue(tsNumberPrototype(loop)),tsPropertyKey(key))
  case tsErrorKind:v:=(*tsRuntimeError)(value.ref);if v.object!=nil{return tsObjectRead(loop,tsObjectValue(v.object),tsPropertyKey(key))};if tsText(key)=="message" {return tsStringUTF8(v.message)};if tsText(key)=="name" {return tsStringUTF8(v.name)};if result,ok:=v.fields[tsText(key)];ok{return result};return tsPrototypeLookup(loop,value,tsValue(tsErrorPrototype(v.name)),tsPropertyKey(key))
 case tsNamespaceKind:v:=(*tsNamespace)(value.ref);return v.module.read(tsPropertyKey(key))
 case tsPromiseKind:v:=(*tsPromise)(value.ref);name:=tsPropertyKey(key);if property,ok:=v.properties.values[name];ok{return property}
  if name=="then" || name=="catch" {return tsFunc(func(args ...tsValue) tsValue {
   next:=loop.promise();success,failure:=tsArg(args,0),tsArg(args,1);if tsText(key)=="catch" {failure=success;success=tsU}
   loop.await(tsPromiseValue(v),func(result tsResult){callback:=success;if result.rejected {callback=failure};if callback.kind!=tsFunctionKind {next.settle(result);return};func(){defer func(){if err:=recover();err!=nil {next.settle(tsResult{tsUnwrap(err),true})}}();next.settle(tsResult{tsCall(callback,result.value),false})}()});return next
  })}
 }
 panic("Unsupported property access")
}
func tsSet(value,key,item tsValue)tsValue{switch value.kind{case tsErrorKind:return tsObjectWrite(loop,tsObjectValue(tsErrorObject(loop,(*tsRuntimeError)(value.ref))),tsPropertyKey(key),item);case tsFunctionKind:return tsFunctionWrite(loop,value,tsPropertyKey(key),item);case tsPromiseKind:return tsObjectWrite(loop,tsObjectValue((*tsPromise)(value.ref).properties),tsPropertyKey(key),item);case tsTypedArrayKind:return tsValue(tsTypedArraySet(value,key,item));case tsDataViewKind:return tsValue(tsDataViewSet(value,key,item));case tsCollectionKind:c:=(*tsCollection)(value.ref);if c.properties==nil{c.properties=map[string]tsValue{}};name:=tsPropertyKey(key);if name=="size"{panic("Collection size is read-only")};if _,ok:=c.properties[name];!ok{c.order=append(c.order,name)};c.properties[name]=item;return item;case tsClassKind:return tsInstanceWrite(loop,value,tsPropertyKey(key),item);case tsInstanceKind:return tsInstanceWrite(loop,value,tsPropertyKey(key),item);case tsObjectKind:return tsObjectWrite(loop,value,tsPropertyKey(key),item);case tsRegExpKind:return tsValue(tsRegExpSet(value,item,tsPropertyKey(key)))};if value.kind!=tsArrayKind{panic("Expected an array")};array:=(*tsArray)(value.ref);if index,ok:=tsGrowableIndex(key);ok {if index==array.length(){array.appendItems(item);return item};if index<array.length(){array.put(index,item);return item};if array.native!=nil{tsArrayRangeFailure("Dense array writes cannot create holes")};for array.length()<=index{tsArrayHole(array)};array.put(index,item);return item};name:=tsPropertyKey(key);if name=="length"{length:=tsNumber(item);if length<0||length>4294967295||math.Trunc(length)!=length{panic("Invalid array length")};for array.length()<int(length){tsArrayHole(array)};array.resize(int(length));for index:=range array.holes{if index>=int(length){delete(array.holes,index)}};return item};index,err:=strconv.Atoi(name);if err!=nil||index<0||uint64(index)>=4294967295||strconv.Itoa(index)!=name{if array.properties==nil{array.properties=map[string]tsValue{}};if _,exists:=array.properties[name];!exists{array.order=append(array.order,name)};array.properties[name]=item;return item};if array.native!=nil&&index==array.length(){array.appendItems(item);return item};for array.length()<=index{tsArrayHole(array)};array.put(index,item);return item}
func tsNumber(value tsValue) float64 {switch value.kind {case tsSymbolKind:tsPropertyFailure("Cannot convert Symbol to number");return 0;case tsBigIntKind:panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"TypeError",message:"Cannot convert BigInt to number implicitly"})});case tsIntKind,tsInt8Kind,tsInt16Kind,tsInt32Kind,tsInt64Kind:return float64(int64(math.Float64bits(value.number)));case tsUintKind,tsUint8Kind,tsUint16Kind,tsUint32Kind,tsUint64Kind:return float64(math.Float64bits(value.number));case tsFloat32Kind,tsNumberKind,tsBooleanKind:return value.number;case tsNullKind:return 0;case tsUndefinedKind:return math.NaN();case tsStringKind:return tsStringNumber((*tsString)(value.ref));default:panic("Expected a number")}}
func tsTruthy(value tsValue) bool {switch value.kind {case tsBigIntKind:return (*big.Int)(value.ref).Sign()!=0;case tsUndefinedKind,tsNullKind:return false;case tsIntKind,tsInt8Kind,tsInt16Kind,tsInt32Kind,tsInt64Kind,tsUintKind,tsUint8Kind,tsUint16Kind,tsUint32Kind,tsUint64Kind:return math.Float64bits(value.number)!=0;case tsFloat32Kind,tsNumberKind:return value.number!=0 && !math.IsNaN(value.number);case tsBooleanKind:return value.number!=0;case tsStringKind:return len((*tsString)(value.ref).units)!=0;default:return true}}
func tsBinary(op string, left, right tsValue) tsValue {
 switch op{case "+","-","*","/","%","**","&","|","^","<<",">>",">>>","<","<=",">",">=":
    if op=="+"{left=tsValue(tsDefaultPrimitive(loop,left));right=tsValue(tsDefaultPrimitive(loop,right))}else{left=tsValue(tsToPrimitive(left));right=tsValue(tsToPrimitive(right))}
 }
 if left.kind==tsBigIntKind||right.kind==tsBigIntKind{return tsValue(tsBigIntBinary(op,left,right))}
 if left.kind==right.kind {switch left.kind {
 case tsIntKind:if op=="+"||op=="-"||op=="*"||op=="/"||op=="%"||op=="**"||op=="&"||op=="|"||op=="^"||op=="<<"||op==">>"||op==">>>"{return tsIntegerBinary[int](op,left,right)}
 case tsInt8Kind:if op=="+"||op=="-"||op=="*"||op=="/"||op=="%"||op=="**"||op=="&"||op=="|"||op=="^"||op=="<<"||op==">>"||op==">>>"{return tsIntegerBinary[int8](op,left,right)}
 case tsInt16Kind:if op=="+"||op=="-"||op=="*"||op=="/"||op=="%"||op=="**"||op=="&"||op=="|"||op=="^"||op=="<<"||op==">>"||op==">>>"{return tsIntegerBinary[int16](op,left,right)}
 case tsInt32Kind:if op=="+"||op=="-"||op=="*"||op=="/"||op=="%"||op=="**"||op=="&"||op=="|"||op=="^"||op=="<<"||op==">>"||op==">>>"{return tsIntegerBinary[int32](op,left,right)}
 case tsInt64Kind:if op=="+"||op=="-"||op=="*"||op=="/"||op=="%"||op=="**"||op=="&"||op=="|"||op=="^"||op=="<<"||op==">>"||op==">>>"{return tsIntegerBinary[int64](op,left,right)}
 case tsUintKind:if op=="+"||op=="-"||op=="*"||op=="/"||op=="%"||op=="**"||op=="&"||op=="|"||op=="^"||op=="<<"||op==">>"||op==">>>"{return tsIntegerBinary[uint](op,left,right)}
 case tsUint8Kind:if op=="+"||op=="-"||op=="*"||op=="/"||op=="%"||op=="**"||op=="&"||op=="|"||op=="^"||op=="<<"||op==">>"||op==">>>"{return tsIntegerBinary[uint8](op,left,right)}
 case tsUint16Kind:if op=="+"||op=="-"||op=="*"||op=="/"||op=="%"||op=="**"||op=="&"||op=="|"||op=="^"||op=="<<"||op==">>"||op==">>>"{return tsIntegerBinary[uint16](op,left,right)}
 case tsUint32Kind:if op=="+"||op=="-"||op=="*"||op=="/"||op=="%"||op=="**"||op=="&"||op=="|"||op=="^"||op=="<<"||op==">>"||op==">>>"{return tsIntegerBinary[uint32](op,left,right)}
 case tsUint64Kind:if op=="+"||op=="-"||op=="*"||op=="/"||op=="%"||op=="**"||op=="&"||op=="|"||op=="^"||op=="<<"||op==">>"||op==">>>"{return tsIntegerBinary[uint64](op,left,right)}
} }
 if tsIsNumeric(left)&&tsIsNumeric(right) {switch op{case "<","<=",">",">=":comparison,valid:=tsNumericCompare(left,right);if !valid{return false};switch op{case "<":return comparison<0;case "<=":return comparison<=0;case ">":return comparison>0;default:return comparison>=0}}}

 if op=="<"||op=="<="||op==">"||op==">="{
  comparison,valid:=0,false
  integerLeft:=tsIsSigned(left)||tsIsUnsigned(left);integerRight:=tsIsSigned(right)||tsIsUnsigned(right)
  if integerLeft&&right.kind==tsStringKind{comparison,valid=tsIntegerStringCompare(left,(*tsString)(right.ref))}else if integerRight&&left.kind==tsStringKind{comparison,valid=tsIntegerStringCompare(right,(*tsString)(left.ref));comparison=-comparison}
  if valid{switch op{case "<":return comparison<0;case "<=":return comparison<=0;case ">":return comparison>0;default:return comparison>=0}}
  if integerLeft&&right.kind==tsStringKind||integerRight&&left.kind==tsStringKind{return false}
 }
 switch op {
 case "in":return tsBooleanValue(tsHas(right,left))
 case "==","!=":equal:=tsLooseEqual(left,right);if op=="!="{return tsBooleanValue(!equal)};return tsBooleanValue(equal)
 case "**":return tsNumberValue(tsPow(tsNumber(tsValue(tsToPrimitive(left))),tsNumber(tsValue(tsToPrimitive(right)))))
 case "&","|","^","<<",">>",">>>":return tsNumberValue(tsUint32Bitwise(op,tsToUint32(left),tsToUint32(right)))
 case "+": if left.kind==tsStringKind{return tsStringConcat(tsStringValue(left),tsStringValue(right))};if right.kind==tsStringKind{return tsStringConcat(tsStringValue(left),tsStringValue(right))}; return tsNumber(left)+tsNumber(right)
 case "-": return tsNumber(left)-tsNumber(right)
 case "*": return tsNumber(left)*tsNumber(right)
 case "/": return tsNumber(left)/tsNumber(right)
 case "%": return math.Mod(tsNumber(left),tsNumber(right))
 case "<":if left.kind==tsStringKind && right.kind==tsStringKind{return tsStringCompare((*tsString)(left.ref),(*tsString)(right.ref))<0}; return tsNumber(left)<tsNumber(right)
 case "<=":if left.kind==tsStringKind && right.kind==tsStringKind{return tsStringCompare((*tsString)(left.ref),(*tsString)(right.ref))<=0}; return tsNumber(left)<=tsNumber(right)
 case ">":if left.kind==tsStringKind && right.kind==tsStringKind{return tsStringCompare((*tsString)(left.ref),(*tsString)(right.ref))>0}; return tsNumber(left)>tsNumber(right)
 case ">=":if left.kind==tsStringKind && right.kind==tsStringKind{return tsStringCompare((*tsString)(left.ref),(*tsString)(right.ref))>=0}; return tsNumber(left)>=tsNumber(right)
 case "===", "!==": equal:=tsStrictEqual(left,right); if op=="!==" { return tsBooleanValue(!equal) }; return tsBooleanValue(equal)
 default: panic("Unsupported operator")
 }
}

type tsThrown struct {value tsValue}
// recover() returns Go's any; translate it once at the panic boundary.
func tsUnwrap(value any)tsValue{switch v:=value.(type){case tsThrown:return v.value;case tsValue:return v;case error:return tsStringReference(tsStringUTF8(v.Error()));case string:return tsStringReference(tsStringUTF8(v));case nil:return tsNull;default:return tsStringReference(tsStringUTF8(fmt.Sprint(v)))}}
type tsResult struct { value tsValue; rejected bool }
type tsAsyncContextFrame struct { storage unsafe.Pointer; value tsValue; prev *tsAsyncContextFrame; masked map[unsafe.Pointer]bool }
func tsWithAsyncContext(loop *tsLoop, context *tsAsyncContextFrame, callback func()) { previous:=loop.asyncContext; loop.asyncContext=context; defer func(){loop.asyncContext=previous}(); callback() }
type tsPromise struct {mu sync.Mutex; loop *tsLoop; properties *tsObject; settled, handled bool; result tsResult; listeners []func(tsResult) }
// A task is registered and finished by the loop. The worker exclusively writes
// result until publishing this pointer through the completion channel.
type tsTask struct { work func() tsResult; finish func(tsResult); result tsResult; context *tsAsyncContextFrame; canceled bool; cancel func(); timer, completed, unref bool; deadline time.Time }
func tsDiagnosticText(loop *tsLoop,value tsValue)(text string){text=tsText(value);defer func(){if recover()!=nil{return}}();return tsStringValueWithLoop(loop,value).String()}
type tsLoop struct { stop chan struct{};resourceClosers []func(); queueMu sync.Mutex;incoming []func();wake chan struct{}; pending map[*tsTask]struct{}; completions chan *tsTask; ready []func(); promises []*tsPromise; errors []tsValue; timers []*tsTask; modules []*tsModule; asyncContext *tsAsyncContextFrame; regexpStaticResultOffsetsVector []int32; regexpStaticResultOffsetsVectorInitialized bool; regexpLastMatchInfo *tsRegExpMatchInfo }
func tsNewLoop() *tsLoop { return &tsLoop{stop:make(chan struct{}),pending:make(map[*tsTask]struct{}), completions:make(chan *tsTask),wake:make(chan struct{},1)} }
func (l *tsLoop) promise() *tsPromise { p:= &tsPromise{loop:l,properties:tsNewObject()}; l.promises=append(l.promises,p); return p }
func(l *tsLoop)post(callback func()){l.queueMu.Lock();l.incoming=append(l.incoming,callback);l.queueMu.Unlock();select{case l.wake<-struct{}{}:default:}}
func(l *tsLoop)drain(){l.queueMu.Lock();l.ready=append(l.ready,l.incoming...);l.incoming=nil;l.queueMu.Unlock()}
func (p *tsPromise) settle(result tsResult) {
 if result.value.kind==tsPromiseKind&&!result.rejected {
  other:=(*tsPromise)(result.value.ref)
  if other==p{p.settle(tsResult{"Promise cannot resolve to itself",true});return}
  other.then(p.settle);return
 }
 p.mu.Lock();if p.settled {p.mu.Unlock();return};p.settled=true;p.result=result;listeners:=p.listeners;p.listeners=nil;p.mu.Unlock()
 for _,callback:=range listeners{cb:=callback;p.loop.post(func(){cb(result)})}
}
func(p *tsPromise)then(callback func(tsResult)){context:=p.loop.asyncContext;invoke:=func(result tsResult){tsWithAsyncContext(p.loop,context,func(){callback(result)})};p.mu.Lock();p.handled=true;if !p.settled{p.listeners=append(p.listeners,invoke);p.mu.Unlock();return};result:=p.result;p.mu.Unlock();p.loop.post(func(){invoke(result)})}
func (l *tsLoop) resolved(value tsValue, rejected bool) *tsPromise { p:=l.promise(); p.settle(tsResult{value,rejected}); return p }
func (l *tsLoop) constructPromise(executor tsValue) *tsPromise {
 p:=l.promise();finished:=false
 resolve:=tsFunc(func(args ...tsValue) tsValue {if !finished {finished=true;p.settle(tsResult{tsArg(args,0),false})};return tsU})
 reject:=tsFunc(func(args ...tsValue) tsValue {if !finished {finished=true;p.settle(tsResult{tsArg(args,0),true})};return tsU})
 func(){defer func(){if err:=recover();err!=nil && !finished {finished=true;p.settle(tsResult{tsUnwrap(err),true})}}();tsCall(executor,resolve,reject)}()
 return p
}
func (l *tsLoop) all(values tsValue) *tsPromise {
 p:=l.promise();if values.kind!=tsArrayKind{p.settle(tsResult{"Promise.all expects an array",true});return p};array:=(*tsArray)(values.ref)
 output:= &tsArray{values:make([]tsValue,array.length())};remaining:=array.length()
 if remaining==0 {p.settle(tsResult{output,false});return p}
 for index:=0;index<array.length();index++ {value:=array.at(index);i:=index;l.await(value,func(result tsResult){if result.rejected {p.settle(result);return};output.values[i]=result.value;remaining--;if remaining==0 {p.settle(tsResult{output,false})}})}
 return p
}
func(l *tsLoop)await(value tsValue,callback func(tsResult)){
 if value.kind!=tsPromiseKind {context:=l.asyncContext;l.ready=append(l.ready,func(){tsWithAsyncContext(l,context,func(){callback(tsResult{value,false})})});return}
 p:=(*tsPromise)(value.ref);if p.loop==l{p.then(callback);return}
 p.mu.Lock();p.handled=true;if p.settled{result:=p.result;p.mu.Unlock();l.ready=append(l.ready,func(){callback(result)});return};p.mu.Unlock()
 done:=make(chan tsResult,1);p.then(func(result tsResult){done<-result});l.submit(func()tsResult{return <-done},callback)
}
// submit is called only on the owning loop. Native adapters snapshot inputs;
// user workers run application code with their own loop and shared references.
func (l *tsLoop) submit(work func() tsResult, finish func(tsResult)) *tsTask {
 task:= &tsTask{work:work,finish:finish,context:l.asyncContext}
 l.pending[task]=struct{}{}
 go func(task *tsTask,done chan<- *tsTask) {
  result:=tsResult{"Worker terminated without a result",true}
  defer func(){if failure:=recover();failure!=nil {result=tsResult{tsUnwrap(failure),true}};task.result=result;select{case done<-task:case <-l.stop:}}()
  result=task.work()
 }(task,l.completions)
 return task
}
// spawn reuses task registration and publishes only completion to the caller.
// The tagged argument slice copies scalar slots; references remain shared.
func(l *tsLoop)spawn(function tsValue,args ...tsValue)*tsPromise{
 if function.kind!=tsFunctionKind{panic("go requires a function call")}
 return l.start(func()tsResult{
  child:=tsNewLoop();loop:=child;_ = loop;result:=tsResult{tsU,false};completed:=false
  func(){defer func(){if failure:=recover();failure!=nil{result=tsResult{tsUnwrap(failure),true};completed=true}}();value:=tsCall(function,args...);child.await(value,func(r tsResult){result=r;completed=true})}()
  if err:=child.run();err!=nil{return tsResult{tsErrorValue(&tsRuntimeError{name:"Error",message:err.Error()}),true}}
  if !completed{return tsResult{"Worker did not complete",true}};return result
 })
}
func (l *tsLoop) start(work func() tsResult) *tsPromise {p:=l.promise();l.submit(work,p.settle);return p}
// Browser-style timers return a handle, not a Promise. The worker only sleeps;
// callback invocation and captured-variable access happen on the event loop.
func tsTimerDuration(value tsValue) time.Duration {
 milliseconds:=float64(0)
 if !tsIsUndefined(value) {milliseconds=tsNumber(value)}
 if math.IsNaN(milliseconds) || math.IsInf(milliseconds,0) || milliseconds<0 || milliseconds>2147483647 {milliseconds=0}
 return time.Duration(math.Trunc(milliseconds))*time.Millisecond
}
func (l *tsLoop) setTimeout(callback, milliseconds tsValue, args ...tsValue) *tsTask {
 if task, mocked := tsInternalMockTimerCreate(l, callback, milliseconds, args...); mocked { return task }
 if callback.kind!=tsFunctionKind {panic("setTimeout expects a function")}
 duration:=tsTimerDuration(milliseconds)
 deadline:=time.Now().Add(duration)
 canceled:=make(chan struct{})
 // Only the finish function captures the callback and its arguments.
 var task *tsTask
 task=l.submit(func()tsResult {timer:=time.NewTimer(time.Until(deadline));defer timer.Stop();select {case <-timer.C:case <-canceled:};return tsResult{tsU,false}},func(result tsResult){
  if task.canceled {return}
  if result.rejected {l.errors=append(l.errors,result.value);return}
  tsCall(callback,args...)
 })
 task.cancel=func(){close(canceled)}
 task.timer=true;task.deadline=deadline
 l.timers=append(l.timers,task)
 sort.SliceStable(l.timers,func(i,j int)bool{return l.timers[i].deadline.Before(l.timers[j].deadline)})
 return task
}
func(l *tsLoop)clearTimeout(handle tsValue)tsValue{if handle.kind==tsTaskKind{task:=(*tsTask)(handle.ref);if !task.canceled{task.canceled=true;if task.cancel!=nil{task.cancel()};if task.completed{l.removeTimer(task);delete(l.pending,task)}}};return tsU}
func(l *tsLoop)readFile(path tsValue)*tsPromise{if path.kind!=tsStringKind{return l.resolved("readFile expects a string",true)};text:=(*tsString)(path.ref);return l.start(func()tsResult{bytes,err:=os.ReadFile(text.String());if err!=nil{return tsResult{err.Error(),true}};return tsResult{tsStringUTF8(string(bytes)),false}})}
func (l *tsLoop) delay(milliseconds tsValue) *tsPromise { duration:=tsTimerDuration(milliseconds); return l.start(func() tsResult {time.Sleep(duration); return tsResult{tsU,false}}) }
func (l *tsLoop) invoke(callback func()) {
 defer func(){if failure:=recover();failure!=nil {l.errors=append(l.errors,tsUnwrap(failure))}}()
 callback()
}
func (l *tsLoop) invokeContext(context *tsAsyncContextFrame,callback func()) { tsWithAsyncContext(l,context,func(){l.invoke(callback)}) }
func (l *tsLoop) removeTimer(task *tsTask) {for i,timer:=range l.timers {if timer==task {copy(l.timers[i:],l.timers[i+1:]);l.timers[len(l.timers)-1]=nil;l.timers=l.timers[:len(l.timers)-1];return}}}
func (l *tsLoop) run() error {
 defer func(){for _,closeResource:=range l.resourceClosers{closeResource()};close(l.stop)}()
 for {
  l.drain()
  for len(l.ready)>0 {callback:=l.ready[0];l.ready[0]=nil;l.ready=l.ready[1:];l.invoke(callback);l.drain()}
  // Process one timer at a time, with a reaction checkpoint between callbacks.
  // Deadline ordering prevents goroutine scheduling from reordering equal timers.
  if len(l.timers)>0 && l.timers[0].completed {
   task:=l.timers[0];l.removeTimer(task);delete(l.pending,task);l.invokeContext(task.context,func(){task.finish(task.result)});continue
  }
  l.drain();if len(l.ready)>0{continue}
  active:=false;for pending:=range l.pending{if !pending.unref{active=true;break}};if !active {break}
  var task *tsTask;select{case task=<-l.completions:case <-l.wake:continue}
  if _,ok:=l.pending[task];!ok {l.errors=append(l.errors,"Unknown completed task");continue}
  task.completed=true
  if task.timer {
   if task.canceled {l.removeTimer(task);delete(l.pending,task)}
   continue
  }
  delete(l.pending,task)
  l.invokeContext(task.context,func(){task.finish(task.result)})
 }
 for _,module:=range l.modules {if len(module.requests)>0 && !module.failed {missing:=[]string{};for _,request:=range module.requests {if request.full {missing=append(missing,"<module completion>")};for name:=range request.missing {missing=append(missing,name)}};sort.Strings(missing);l.errors=append(l.errors,"Unresolved module dependency: "+module.name+" needs "+strings.Join(missing,", "))}}
 for _,p:=range l.promises {p.mu.Lock();unhandled:=p.settled&&p.result.rejected&&!p.handled;value:=p.result.value;p.mu.Unlock();if unhandled{l.errors=append(l.errors,"Unhandled rejection: "+tsDiagnosticText(loop,value))}}
 if len(l.errors)>0 {messages:=make([]string,len(l.errors));for i,value:=range l.errors {messages[i]=tsDiagnosticText(loop,value)};return fmt.Errorf("Async runtime errors: %s",strings.Join(messages,"; "))}
 return nil
}

// Protected regions and pending abrupt completions mirror ES5's trys/ops stacks.
type tsHandler struct { catch, finally, end int; phase int }
type tsAbrupt struct { kind string; value tsValue; target, depth int;nativeSlot int }
type tsMachine struct {generatorPrototype tsValue;asyncDelegate *tsAsyncIterator;asyncGenerator *tsAsyncGenerator;generator,started,running bool;yielded,delegatedResult tsValue;delegate *tsIterator; loop *tsLoop; pc int; step func(); output *tsPromise; async, done, suspended bool; result tsValue; handlers []tsHandler; abrupt []tsAbrupt; thrown tsValue; module *tsModule; blocked bool;nativeSlot int }
func (m *tsMachine) resume() {
 for !m.done && !m.suspended {
  func(){ defer func(){if err:=recover(); err!=nil {m.raise(tsUnwrap(err))}}(); m.step() }()
 }
 if m.asyncGenerator!=nil{m.asyncGenerator.notify(loop)}
}
func (m *tsMachine) await(value tsValue, next int) {
 m.suspended=true;m.blocked=true
 m.loop.await(value,func(result tsResult){m.blocked=false;m.suspended=false; if result.rejected {m.raise(result.value)} else {m.result=result.value; m.pc=next}; m.resume()})
}
func (m *tsMachine) complete(value tsValue) { m.done=true; m.result=value; if m.module!=nil {m.module.finish(tsResult{value,false})}; if m.async {m.output.settle(tsResult{value,false})} }
func (m *tsMachine) raise(value tsValue) {m.transfer(tsAbrupt{kind:"throw",value:value})}
func (m *tsMachine) transfer(op tsAbrupt) {
 for len(m.handlers)>op.depth {
  index:=len(m.handlers)-1; handler:= &m.handlers[index]
  if op.kind=="throw" && handler.phase==0 && handler.catch>=0 { handler.phase=1; m.thrown=op.value; m.pc=handler.catch; return }
  if handler.phase<2 && handler.finally>=0 { handler.phase=2; m.abrupt=append(m.abrupt,op); m.pc=handler.finally; return }
  if handler.phase==2 { m.abrupt=m.abrupt[:len(m.abrupt)-1] }
  m.handlers=m.handlers[:index]
 }
 switch op.kind {
 case "jump": m.pc=op.target
 case "return": m.nativeSlot=op.nativeSlot;m.complete(op.value)
 case "throw": m.done=true; if m.asyncGenerator!=nil{m.asyncGenerator.rejected=true;m.asyncGenerator.failure=op.value;return};if m.module!=nil {m.module.finish(tsResult{op.value,true})} else if m.async {m.output.settle(tsResult{op.value,true})} else {panic(tsThrown{op.value})}
 }
}
func (m *tsMachine) endFinally() { index:=len(m.handlers)-1; m.handlers=m.handlers[:index]; last:=len(m.abrupt)-1; op:=m.abrupt[last]; m.abrupt=m.abrupt[:last]; m.transfer(op) }
func tsPow(base,exponent float64)float64{if math.IsNaN(exponent){return math.NaN()};if exponent==0{return 1};if math.IsInf(exponent,0)&&math.Abs(base)==1{return math.NaN()};return math.Pow(base,exponent)}
func tsNumberUint32(number float64)uint32{if number==0||math.IsNaN(number)||math.IsInf(number,0){return 0};number=math.Mod(math.Trunc(number),4294967296);if number<0{number+=4294967296};return uint32(number)}
func tsToUint32(value tsValue)uint32{return tsNumberUint32(tsNumber(tsValue(tsToPrimitive(value))))}
func tsBooleanNumber(value bool)float64{if value{return 1};return 0}
func tsNumberBitwise(op string,left,right float64)float64{return tsUint32Bitwise(op,tsNumberUint32(left),tsNumberUint32(right))}
func tsUint32Bitwise(op string,a,b uint32)float64{switch op{case "&":return float64(int32(a&b));case "|":return float64(int32(a|b));case "^":return float64(int32(a^b));case "<<":return float64(int32(a<<(b&31)));case ">>":return float64(int32(a)>>(b&31));default:return float64(a>>(b&31))}}
func tsNumberBitwiseNot(value float64)float64{return float64(^int32(tsNumberUint32(value)))}
func tsBitwiseNot(value tsValue)tsValue{
 if value.kind==tsBigIntKind{return tsBigIntValue(new(big.Int).Not((*big.Int)(value.ref)))}
 if tsIsSigned(value)||tsIsUnsigned(value){bits:=^math.Float64bits(value.number);return tsIntegerBits(value.kind,bits)}
 return tsNumberValue(float64(^int32(tsToUint32(value))))
}
func tsIsPrimitiveValue(value tsValue)bool{return value.kind==tsSymbolKind||value.kind==tsBigIntKind||tsIsNumeric(value)||value.kind==tsStringKind||value.kind==tsBooleanKind||tsNullish(value)}
func tsToPrimitive(value tsValue)tsValue{return tsPrimitiveHint(loop,value,false)}
func tsDefaultPrimitive(loop *tsLoop,value tsValue)tsValue{return tsPrimitiveHintName(loop,value,"default")}
func tsPrimitiveHint(loop *tsLoop,value tsValue,stringHint bool)tsValue{hint:="number";if stringHint{hint="string"};return tsPrimitiveHintName(loop,value,hint)}
func tsPrimitiveHintName(loop *tsLoop,value tsValue,hint string)tsValue{
 stringHint:=hint=="string"||hint=="default"&&value.kind==tsObjectKind&&(*tsObject)(value.ref).nativeDate
 if tsIsPrimitiveValue(value){return value}
 hook:=tsGet(loop,value,tsValue(tsWellKnownSymbol("toPrimitive")));if !tsNullish(hook){if hook.kind!=tsFunctionKind{tsPropertyFailure("Symbol.toPrimitive is not callable")};result:=tsCallReceiver(loop,hook,value,tsStringReference(tsStringUTF8(hint)));if !tsIsPrimitiveValue(result){tsPropertyFailure("Symbol.toPrimitive must return a primitive")};return result}
 names:=[]string{"valueOf","toString"};if stringHint{names=[]string{"toString","valueOf"}};called:=false
 for _,name:=range names{method:=tsGet(value,tsStringReference(tsStringUTF8(name)));if method.kind==tsFunctionKind{called=true;result:=tsValue(tsCallReceiver(method,value));if tsIsPrimitiveValue(result){return result}}}
 // OrdinaryToPrimitive throws when neither method produces a primitive.
 // Falling back to an internal string here incorrectly succeeds when an
 // object has shadowed both valueOf and toString with non-callable values.
 _=called
 tsPropertyFailure("Cannot convert object to primitive value")
 return tsU
}
func tsLooseEqual(a,b tsValue)bool{
 if a.kind==tsNumberKind||a.kind==tsFloat32Kind{return tsLooseEqualNumberAny(loop,a.number,b)}
 if a.kind==tsBigIntKind||b.kind==tsBigIntKind{if tsNullish(a)||tsNullish(b){return false};if (!tsIsPrimitiveValue(a)||!tsIsPrimitiveValue(b)){return tsLooseEqual(tsDefaultPrimitive(loop,a),tsDefaultPrimitive(loop,b))};comparison,valid:=tsBigIntCompare(a,b);return valid&&comparison==0}
 if tsStrictEqual(a,b){return true};if tsNullish(a)&&tsNullish(b){return true}
 if a.kind==tsBooleanKind{return tsLooseEqual(tsNumberValue(a.number),b)};if b.kind==tsBooleanKind{return tsLooseEqual(a,tsNumberValue(b.number))}
 if (tsIsSigned(a)||tsIsUnsigned(a))&&b.kind==tsStringKind{comparison,valid:=tsIntegerStringCompare(a,(*tsString)(b.ref));return valid&&comparison==0};if (tsIsSigned(b)||tsIsUnsigned(b))&&a.kind==tsStringKind{comparison,valid:=tsIntegerStringCompare(b,(*tsString)(a.ref));return valid&&comparison==0}
 if tsIsNumeric(a)&&b.kind==tsStringKind{return tsStrictEqual(a,tsNumberValue(tsNumber(b)))};if tsIsNumeric(b)&&a.kind==tsStringKind{return tsLooseEqual(b,a)}
 if tsIsPrimitiveValue(a)&&!tsNullish(a)&&!tsIsPrimitiveValue(b){return tsLooseEqual(a,tsDefaultPrimitive(loop,b))};if tsIsPrimitiveValue(b)&&!tsNullish(b)&&!tsIsPrimitiveValue(a){return tsLooseEqual(tsDefaultPrimitive(loop,a),b)}
 return false
}
func tsHas(value,key tsValue)bool{
 if tsIsPrimitiveValue(value){tsPropertyFailure("Right operand of in must be an object")}
 if value.kind==tsErrorKind{return tsHas(tsObjectValue(tsErrorObject(loop,(*tsRuntimeError)(value.ref))),key)}
 name:=tsPropertyKey(key)
 switch value.kind{
 case tsArrayBufferKind:return !tsIsUndefined(tsValue(tsArrayBufferGet(value,name)))
 case tsDataViewKind:return !tsIsUndefined(tsValue(tsDataViewGet(value,name)))
 case tsTypedArrayKind:return !tsIsUndefined(tsValue(tsTypedArrayGet(value,name)))
 case tsCollectionKind:return !tsIsUndefined(tsValue(tsCollectionGet(value,name)))
 case tsObjectKind:object:=(*tsObject)(value.ref);if object.proxy{return tsProxyHas(loop,value,key,object.proxyTarget,object.proxyHandler)};return tsPrototypeHas(value,name)
 case tsFunctionKind,tsInstanceKind:return tsPrototypeHas(value,name)
 case tsClassKind:return tsPrototypeOwn(value,name)||tsPrototypeHas(tsPrototypeOf(value),name)
 case tsPromiseKind:p:=(*tsPromise)(value.ref);return tsOwnObjectProperty(p.properties,name)||name=="then"||name=="catch"
 case tsArrayKind:array:=(*tsArray)(value.ref);if name=="length"||name=="push"{return true};if _,ok:=array.properties[name];ok{return true};if !tsIsUndefined(tsSequenceMethod(value,name)){return true};index,err:=strconv.Atoi(name);return err==nil&&index>=0&&index<array.length()&&strconv.Itoa(index)==name&&!array.holes[index]
 case tsRegExpKind:return tsRegExpHas((*tsRegExp)(value.ref),name)
 }
 return tsPrototypeHas(value,name)
}
`

const EqualityRuntime = `
func tsLooseEqualNumberAny(loop *tsLoop,number float64,value tsValue)bool{
 if tsIsNumeric(value){return tsNumericEqual(tsNumberValue(number),value)}
 switch value.kind{
 case tsBooleanKind:return number==value.number
 case tsStringKind:return number==tsStringNumber((*tsString)(value.ref))
 case tsBigIntKind:comparison,valid:=tsBigIntCompare(value,tsNumberValue(number));return valid&&comparison==0
 }
 if !tsIsPrimitiveValue(value){return tsLooseEqualNumberAny(loop,number,tsDefaultPrimitive(loop,value))}
 return false
}
func tsPolicyEqualNumberAny(loop *tsLoop,number float64,value tsValue,strict,coerce bool)bool{
 value=tsComparisonOperand(value,"number",0,coerce)
 if value.kind==tsNumberKind{return number==value.number}
 if tsIsNumeric(value){return tsNumericEqual(tsNumberValue(number),value)}
 if value.kind==tsBigIntKind{if strict{return false};comparison,valid:=tsBigIntCompare(value,tsNumberValue(number));return valid&&comparison==0}
 if strict{return false}
 return tsLooseEqualNumberAny(loop,number,value)
}
func tsPolicyEqualOptionalNumberAny(loop *tsLoop,number tsOptional[float64],value tsValue,strict,coerce bool)bool{
 if number.tag!=0{
  if !strict{return tsNullish(value)}
  return number.tag==1&&value.kind==tsNullKind||number.tag==2&&value.kind==tsUndefinedKind
 }
 return tsPolicyEqualNumberAny(loop,number.value,value,strict,coerce)
}
func tsPolicyCompareNative[T tsPrimitive](loop *tsLoop,op string,native T,other tsValue,kind string,nulls uint8,coerce,nativeLeft bool)bool{
 known:=tsPrimitiveValue(native);checked:=tsComparisonOperand(other,kind,nulls,coerce);left,right:=checked,known;if nativeLeft{left,right=known,checked};return tsTruthy(tsPolicyCompare(op,left,right,coerce))
}
func tsPolicyCompareOptionalNative[T tsPrimitive](loop *tsLoop,op string,native tsOptional[T],other tsValue,kind string,nulls uint8,coerce,nativeLeft bool)bool{
 if native.tag!=0&&(op=="=="||op=="!="||op=="==="||op=="!=="){
  equal:=false
  if op=="=="||op=="!="{equal=tsNullish(other)}else{equal=native.tag==1&&other.kind==tsNullKind||native.tag==2&&other.kind==tsUndefinedKind}
  if op=="!="||op=="!=="{equal=!equal};return equal
 }
 known:=native.raw();checked:=tsComparisonOperand(other,kind,nulls,coerce);left,right:=checked,known;if nativeLeft{left,right=known,checked};return tsTruthy(tsPolicyCompare(op,left,right,coerce))
}
func tsStrictEqual(a,b tsValue)bool {if tsIsNumeric(a)&&tsIsNumeric(b){return tsNumericEqual(a,b)};if a.kind!=b.kind{return false};switch a.kind {case tsBigIntKind:return (*big.Int)(a.ref).Cmp((*big.Int)(b.ref))==0;case tsUndefinedKind,tsNullKind:return true;case tsNumberKind,tsBooleanKind:return a.number==b.number;case tsStringKind:return tsStringEqual((*tsString)(a.ref),(*tsString)(b.ref));case tsRegExpKind:return a.ref==b.ref;case tsTaskKind,tsModuleKind:return a.ref==b.ref;case tsNamespaceKind:return (*tsNamespace)(a.ref).module==(*tsNamespace)(b.ref).module;case tsImportKind:x,y:=(*tsImportRef)(a.ref),(*tsImportRef)(b.ref);return x.module==y.module&&x.name==y.name;case tsSymbolKind,tsErrorKind,tsFunctionKind,tsObjectKind,tsArrayKind,tsClassKind,tsPromiseKind,tsInstanceKind,tsCollectionKind,tsIteratorKind,tsTypedArrayKind,tsArrayBufferKind,tsDataViewKind:return a.ref==b.ref;default:return false}}

`
