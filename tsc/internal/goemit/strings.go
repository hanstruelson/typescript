package goemit

import (
	"fmt"
	"strings"
	"unicode/utf16"

	"github.com/microsoft/TypeScript/tsc/internal/stringutil"
)

// Scanner text is WTF-8: ordinary rune iteration would corrupt lone surrogates.
func stringLiteral(text string) string {
	var units []string
	for len(text) > 0 {
		r, n := stringutil.DecodeJSStringRune(text)
		text = text[n:]
		if r > 0xffff {
			a, b := utf16.EncodeRune(r)
			units = append(units, fmt.Sprint(uint16(a)), fmt.Sprint(uint16(b)))
		} else {
			units = append(units, fmt.Sprint(uint16(r)))
		}
	}
	return "tsStringUnits([]uint16{" + strings.Join(units, ",") + "})"
}

const StringRuntime = `
// Strings are immutable UTF-16 sequences. Native UTF-8 is used only at I/O and
// library boundaries; indexing, comparison and slicing use code units directly.
type tsString struct {units []uint16}
func tsStringUnits(units []uint16) *tsString {return &tsString{units:units}}
func tsStringUTF8(text string) *tsString {return tsStringUnits(utf16.Encode([]rune(text)))}
func(s *tsString) String() string {return string(utf16.Decode(s.units))}
func tsStringValue(value tsValue)*tsString{switch value.kind{case tsStringKind:return (*tsString)(value.ref);case tsNumberKind:return tsStringUTF8(tsEngine().ToValue(value.number).String());case tsArrayKind:v:=(*tsArray)(value.ref);parts:=[]*tsString{};for i,item:=range v.values{if i>0{parts=append(parts,tsStringUTF8(","))};if !tsNullish(item){parts=append(parts,tsStringValue(item))}};return tsStringConcat(parts...);case tsObjectKind:v:=(*tsObject)(value.ref);if method:=v.values["toString"];method.kind==tsFunctionKind{return tsStringValue(tsCall(method))};return tsStringUTF8("[object Object]");case tsInstanceKind:return tsStringUTF8("[object Object]");case tsECMAKind:return tsStringUTF8((*tsECMAObject)(value.ref).object.String());default:return tsStringUTF8(tsText(value))}}
func tsStringEqual(a,b *tsString) bool {if len(a.units)!=len(b.units){return false};for i,c:=range a.units {if c!=b.units[i]{return false}};return true}
func tsStringCompare(a,b *tsString) int {for i,c:=range a.units {if i>=len(b.units){return 1};if c<b.units[i]{return -1};if c>b.units[i]{return 1}};if len(a.units)<len(b.units){return -1};return 0}
func tsStringConcat(values ...*tsString) *tsString {size:=0;for _,s:=range values {size+=len(s.units)};result:=make([]uint16,0,size);for _,s:=range values {result=append(result,s.units...)};return tsStringUnits(result)}
func tsUint32(value tsValue)uint32 {number:=tsNumber(value);if math.IsNaN(number)||math.IsInf(number,0)||number==0{return 0};number=math.Mod(math.Trunc(number),4294967296);if number<0 {number+=4294967296};return uint32(number)}
func tsInteger(value tsValue) float64 {number:=tsNumber(value);if math.IsNaN(number)||number==0 {return 0};return math.Trunc(number)}
func tsPosition(value tsValue,length int,defaultValue int) int {if tsIsUndefined(value){return defaultValue};number:=tsInteger(value);if number<0{return 0};if number>float64(length){return length};return int(number)}
func tsSlicePosition(value tsValue,length int,defaultValue int) int {if tsIsUndefined(value){return defaultValue};number:=tsInteger(value);if number<0{number+=float64(length);if number<0 {number=0}};if number>float64(length){number=float64(length)};return int(number)}
func(s *tsString) slice(start,end int) *tsString {return tsStringUnits(s.units[start:end])}
func tsStringIndex(s,search *tsString,start int) int {for i:=start;i+len(search.units)<=len(s.units);i++ {match:=true;for j,c:=range search.units {if s.units[i+j]!=c {match=false;break}};if match{return i}};return -1}
func tsStringWhitespace(c uint16) bool {switch c {case 0x9,0xa,0xb,0xc,0xd,0x20,0xa0,0x1680,0x2028,0x2029,0x202f,0x205f,0x3000,0xfeff:return true};return c>=0x2000&&c<=0x200a}
func tsStringWellFormed(s *tsString,replace bool) tsValue {var result []uint16;if replace {result=append([]uint16{},s.units...)};for i:=0;i<len(s.units);i++ {c:=s.units[i];if c>=0xd800&&c<=0xdbff {if i+1<len(s.units)&&s.units[i+1]>=0xdc00&&s.units[i+1]<=0xdfff {i++;continue};if !replace{return false};result[i]=0xfffd}else if c>=0xdc00&&c<=0xdfff {if !replace{return false};result[i]=0xfffd}};if replace{return tsStringUnits(result)};return true}
func tsRejectRegExp(value tsValue){if value.kind==tsRegExpKind{panic("String search argument must not be a RegExp")}}
func tsStringGet(s *tsString,key tsValue) tsValue {
 name:=tsText(key);if name=="length" {return float64(len(s.units))}
 if index,err:=strconv.Atoi(name);err==nil && strconv.Itoa(index)==name {if index<0||index>=len(s.units){return tsU};return s.slice(index,index+1)}
 switch name {case "at","charAt","charCodeAt","codePointAt","concat","endsWith","includes","indexOf","lastIndexOf","localeCompare","match","matchAll","normalize","padEnd","padStart","repeat","replace","replaceAll","search","slice","split","startsWith","substring","substr","toLowerCase","toUpperCase","toLocaleLowerCase","toLocaleUpperCase","trim","trimStart","trimEnd","trimLeft","trimRight","toString","valueOf","isWellFormed","toWellFormed","anchor","big","blink","bold","fixed","fontcolor","fontsize","italics","link","small","strike","sub","sup":return tsFunc(func(args ...tsValue)tsValue{return tsStringMethod(s,name,args)})}
 return tsU
}
func tsStringMethod(s *tsString,name string,args []tsValue) tsValue {
 length:=len(s.units);a,b:=tsArg(args,0),tsArg(args,1)
 switch name {
 case "toString","valueOf":return s
 case "at","charAt","charCodeAt","codePointAt":
  number:=tsInteger(a);if name=="at"&&number<0{number+=float64(length)}
  if number<0||number>=float64(length) {switch name {case "charAt":return tsStringUnits(nil);case "charCodeAt":return math.NaN();default:return tsU}}
  i:=int(number);c:=s.units[i];if name=="charCodeAt" {return float64(c)}
  if name=="codePointAt" {if c>=0xd800&&c<=0xdbff&&i+1<length&&s.units[i+1]>=0xdc00&&s.units[i+1]<=0xdfff {return float64(utf16.DecodeRune(rune(c),rune(s.units[i+1])))};return float64(c)}
  return s.slice(i,i+1)
 case "concat":values:=[]*tsString{s};for _,value:=range args {values=append(values,tsStringValue(value))};return tsStringConcat(values...)
 case "slice":start,end:=tsSlicePosition(a,length,0),tsSlicePosition(b,length,length);if end<start {end=start};return s.slice(start,end)
 case "substring":start,end:=tsPosition(a,length,0),tsPosition(b,length,length);if end<start {start,end=end,start};return s.slice(start,end)
 case "substr":start:=tsSlicePosition(a,length,0);count:=tsPosition(b,length-start,length-start);return s.slice(start,start+count)
 case "indexOf","includes","startsWith","endsWith","lastIndexOf":
  if name=="includes"||name=="startsWith"||name=="endsWith" {tsRejectRegExp(a)}
  search:=tsStringValue(a);start:=tsPosition(b,length,0)
  if name=="lastIndexOf" {if tsIsUndefined(b)||math.IsNaN(tsNumber(b)){start=length};if start>length-len(search.units){start=length-len(search.units)};for i:=start;i>=0;i-- {if tsStringIndex(s.slice(i,i+len(search.units)),search,0)==0 {return float64(i)}};return float64(-1)}
  if name=="endsWith" {end:=tsPosition(b,length,length);start=end-len(search.units);return start>=0&&tsStringIndex(s.slice(start,end),search,0)==0}
  index:=tsStringIndex(s,search,start);switch name {case "includes":return index>=0;case "startsWith":return index==start;default:return float64(index)}
 case "repeat":count:=tsInteger(a);if count<0||math.IsInf(count,0){panic("Invalid string repeat count")};if length==0{return s};if count>float64(1<<28)/float64(length){panic("Invalid string length")};values:=make([]uint16,0,int(count)*length);for i:=0;i<int(count);i++ {values=append(values,s.units...)};return tsStringUnits(values)
 case "padStart","padEnd":target:=tsInteger(a);if target<=float64(length){return s};if target>1<<28 {panic("Invalid string length")};fill:=tsStringValue(b);if tsIsUndefined(b){fill=tsStringUTF8(" ")};if len(fill.units)==0{return s};padding:=make([]uint16,int(target)-length);for i:=range padding {padding[i]=fill.units[i%len(fill.units)]};if name=="padStart" {return tsStringConcat(tsStringUnits(padding),s)};return tsStringConcat(s,tsStringUnits(padding))
 case "trim","trimStart","trimEnd","trimLeft","trimRight":start,end:=0,length;if name!="trimEnd"&&name!="trimRight" {for start<end&&tsStringWhitespace(s.units[start]){start++}};if name!="trimStart"&&name!="trimLeft" {for end>start&&tsStringWhitespace(s.units[end-1]){end--}};return s.slice(start,end)
 case "isWellFormed":return tsStringWellFormed(s,false)
 case "toWellFormed":return tsStringWellFormed(s,true)
 case "split":if a.kind==tsRegExpKind{return tsRegexString(s,name,args)};limit:=uint32(0xffffffff);if !tsIsUndefined(b){limit=tsUint32(b)};out:= &tsArray{};if limit==0{return out};if tsIsUndefined(a){out.values=[]tsValue{s};return out};separator:=tsStringValue(a);if len(separator.units)==0 {for i:=0;i<length&&uint32(len(out.values))<limit;i++ {out.values=append(out.values,s.slice(i,i+1))};return out};start:=0;for uint32(len(out.values))<limit {index:=tsStringIndex(s,separator,start);if index<0 {out.values=append(out.values,s.slice(start,length));break};out.values=append(out.values,s.slice(start,index));start=index+len(separator.units)};return out
 case "localeCompare","normalize","toLowerCase","toUpperCase","toLocaleLowerCase","toLocaleUpperCase":return tsUnicodeString(s,name,args)
 case "match","matchAll","replace","replaceAll","search":return tsRegexString(s,name,args)
 }
 tags:=map[string]string{"anchor":"a","big":"big","blink":"blink","bold":"b","fixed":"tt","fontcolor":"font","fontsize":"font","italics":"i","link":"a","small":"small","strike":"strike","sub":"sub","sup":"sup"};tag,ok:=tags[name];if !ok {panic("Unsupported String method")};prefix:="<"+tag;attribute:="";switch name {case "anchor":attribute="name";case "link":attribute="href";case "fontcolor":attribute="color";case "fontsize":attribute="size"};if attribute!="" {value:=tsStringValue(a);escaped:=[]uint16{};for _,c:=range value.units {if c==34 {escaped=append(escaped,[]uint16{'&','q','u','o','t',';'}...)}else {escaped=append(escaped,c)}};return tsStringConcat(tsStringUTF8(prefix+" "+attribute+"=\""),tsStringUnits(escaped),tsStringUTF8("\">"),s,tsStringUTF8("</"+tag+">"))};return tsStringConcat(tsStringUTF8(prefix+">"),s,tsStringUTF8("</"+tag+">"))
}
func tsStringStatic(name string,args []tsValue) tsValue {units:=[]uint16{};switch name {case "fromCharCode":for _,value:=range args {units=append(units,uint16(tsUint32(value)))};case "fromCodePoint":for _,value:=range args {number:=tsNumber(value);if math.IsNaN(number)||number<0||number>0x10ffff||math.Trunc(number)!=number {panic("Invalid code point")};if number<=0xffff {units=append(units,uint16(number))}else {a,b:=utf16.EncodeRune(rune(number));units=append(units,uint16(a),uint16(b))}};case "raw":raw:=tsGet(tsArg(args,0),tsStringUTF8("raw"));length:=int(tsNumber(tsGet(raw,tsStringUTF8("length"))));values:=[]*tsString{};for i:=0;i<length;i++ {values=append(values,tsStringValue(tsGet(raw,float64(i))));if i+1<length && i+1<len(args) {values=append(values,tsStringValue(args[i+1]))}};return tsStringConcat(values...);default:panic("Unknown String static method")};return tsStringUnits(units)}

// Transform valid scalar segments while preserving isolated UTF-16 surrogates.
func tsStringTransform(s *tsString, transform func(string)string)*tsString {
 out:=[]uint16{};start:=0
 flush:=func(end int){if end>start {out=append(out,utf16.Encode([]rune(transform(string(utf16.Decode(s.units[start:end])))))...)}}
 for i:=0;i<len(s.units);i++ {c:=s.units[i];if c>=0xd800&&c<=0xdbff&&i+1<len(s.units)&&s.units[i+1]>=0xdc00&&s.units[i+1]<=0xdfff {i++;continue};if c>=0xd800&&c<=0xdfff {flush(i);out=append(out,c);start=i+1}}
 flush(len(s.units));return tsStringUnits(out)
}
func tsStringLocale(value tsValue)language.Tag {
 if tsIsUndefined(value){return language.Und};if value.kind==tsArrayKind {array:=(*tsArray)(value.ref);if len(array.values)==0{return language.Und};value=array.values[0]}
 tag,err:=language.Parse(tsStringValue(value).String());if err!=nil {panic(tsThrown{&tsRuntimeError{name:"RangeError",message:"Invalid language tag"}})};return tag
}
func tsUnicodeString(s *tsString,name string,args []tsValue)tsValue {
 switch name {
 case "normalize":form:=tsArg(args,0);label:="NFC";if !tsIsUndefined(form){label=tsStringValue(form).String()};var normalization norm.Form;switch label {case "NFC":normalization=norm.NFC;case "NFD":normalization=norm.NFD;case "NFKC":normalization=norm.NFKC;case "NFKD":normalization=norm.NFKD;default:panic(tsThrown{&tsRuntimeError{name:"RangeError",message:"Invalid normalization form"}})};return tsStringTransform(s,normalization.String)
 case "toLowerCase","toUpperCase","toLocaleLowerCase","toLocaleUpperCase":tag:=language.Und;if strings.HasPrefix(name,"toLocale"){tag=tsStringLocale(tsArg(args,0))};if name=="toLowerCase"||name=="toLocaleLowerCase" {return tsStringTransform(s,cases.Lower(tag).String)};return tsStringTransform(s,cases.Upper(tag).String)
 case "localeCompare":tag:=tsStringLocale(tsArg(args,1));options:=[]collate.Option{};config:=tsArg(args,2);if !tsNullish(config){sensitivity:=tsGet(config,tsStringUTF8("sensitivity"));if !tsIsUndefined(sensitivity){switch tsStringValue(sensitivity).String(){case "base":options=append(options,collate.IgnoreCase,collate.IgnoreDiacritics);case "accent":options=append(options,collate.IgnoreCase);case "case":options=append(options,collate.IgnoreDiacritics);case "variant":default:panic(tsThrown{&tsRuntimeError{name:"RangeError",message:"Invalid collation sensitivity"}})}};if tsTruthy(tsGet(config,tsStringUTF8("numeric"))){options=append(options,collate.Numeric)}};return float64(collate.New(tag,options...).CompareString(s.String(),tsStringValue(tsArg(args,0)).String()))
 };panic("Unsupported Unicode operation")
}
`
