package goemit

// internalInspectRuntime adapts inspect.js's option handling, recursive
// formatter, circular-reference guard, and printf-style helpers.
const internalInspectRuntime = `
type tsNodeInspectOptions struct { depth,maxArrayLength,maxStringLength,breakLength,compact int; colors,sorted,numericSeparator,showHidden,showProxy,getters,customInspect bool }
var tsInternalInspectState struct { sync.Once; activeMu sync.Mutex; active map[unsafe.Pointer]bool; module *tsObject; customSymbol tsValue; defaults *tsObject }
func tsNodeInspectOptionsFrom(loop *tsLoop,args []tsValue)tsNodeInspectOptions{
	o:=tsNodeInspectOptions{depth:2,maxArrayLength:100,maxStringLength:10000,breakLength:80,compact:3,customInspect:true}
	if defaults:=tsInternalInspectState.defaults;defaults!=nil{if n:=defaults.get("depth");tsIsNumeric(n){o.depth=int(tsNumber(n))};if n:=defaults.get("maxArrayLength");tsIsNumeric(n){o.maxArrayLength=int(tsNumber(n))};if n:=defaults.get("maxStringLength");tsIsNumeric(n){o.maxStringLength=int(tsNumber(n))};if n:=defaults.get("breakLength");tsIsNumeric(n){o.breakLength=int(tsNumber(n))};if n:=defaults.get("compact");n.kind==tsBooleanKind{if tsTruthy(n){o.compact=-1}else{o.compact=0}}else if tsIsNumeric(n){o.compact=int(tsNumber(n))};o.colors=tsTruthy(defaults.get("colors"));o.sorted=tsTruthy(defaults.get("sorted"));o.numericSeparator=tsTruthy(defaults.get("numericSeparator"));o.showHidden=tsTruthy(defaults.get("showHidden"));o.showProxy=tsTruthy(defaults.get("showProxy"));o.getters=tsTruthy(defaults.get("getters"));if n:=defaults.get("customInspect");!tsIsUndefined(n){o.customInspect=tsTruthy(n)}}
	if len(args)>1 {v:=args[1];if v.kind==tsBooleanKind&&tsTruthy(v){o.showHidden=true}else if v.kind==tsObjectKind||v.kind==tsInstanceKind||v.kind==tsArrayKind{
		read:=func(name string)tsValue{return tsGet(loop,v,tsStringReference(tsStringUTF8(name)))}
		if n:=read("depth");n.kind==tsNullKind{o.depth=-1}else if tsIsNumeric(n){o.depth=int(tsNumber(n))}
		if n:=read("maxArrayLength");n.kind==tsNullKind{o.maxArrayLength=1<<30}else if tsIsNumeric(n){o.maxArrayLength=int(tsNumber(n))}
		if n:=read("maxStringLength");n.kind==tsNullKind{o.maxStringLength=1<<30}else if tsIsNumeric(n){o.maxStringLength=int(tsNumber(n))}
		if n:=read("breakLength");tsIsNumeric(n){o.breakLength=int(tsNumber(n))}
		if n:=read("compact");n.kind==tsBooleanKind{if tsTruthy(n){o.compact=-1}else{o.compact=0}}else if tsIsNumeric(n){o.compact=int(tsNumber(n))}
		o.colors=tsTruthy(read("colors"));o.sorted=tsTruthy(read("sorted"));o.numericSeparator=tsTruthy(read("numericSeparator"));o.showHidden=tsTruthy(read("showHidden"));o.showProxy=tsTruthy(read("showProxy"));o.getters=tsTruthy(read("getters"));if custom:=read("customInspect");!tsIsUndefined(custom){o.customInspect=tsTruthy(custom)}
	}}
	if len(args)>2&&tsIsNumeric(args[2]){o.depth=int(tsNumber(args[2]))};if len(args)>3&&!tsIsUndefined(args[3]){o.colors=tsTruthy(args[3])}
	if o.maxArrayLength<0{o.maxArrayLength=1<<30};if o.maxStringLength<0{o.maxStringLength=1<<30};return o
}
func tsNodeInspectGroupDigits(value string)string{sign:="";if strings.HasPrefix(value,"-"){sign="-";value=value[1:]}else if strings.HasPrefix(value,"+"){sign="+";value=value[1:]};if len(value)<=3{return sign+value};first:=len(value)%3;if first==0{first=3};parts:=[]string{value[:first]};for i:=first;i<len(value);i+=3{parts=append(parts,value[i:i+3])};return sign+strings.Join(parts,"_")}
func tsNodeInspectNumberText(number float64,numericSeparator bool)string{value:=tsNumberText(number);if !numericSeparator||math.IsNaN(number)||math.IsInf(number,0){return value};if dot:=strings.IndexByte(value,'.');dot>=0{return tsNodeInspectGroupDigits(value[:dot])+"."+tsNodeInspectGroupDigits(value[dot+1:])};if exp:=strings.IndexAny(value,"eE");exp>=0{return tsNodeInspectGroupDigits(value[:exp])+value[exp:]};return tsNodeInspectGroupDigits(value)}
func tsNodeInspectBigIntText(value *big.Int,numericSeparator bool)string{text:=value.String();if numericSeparator{text=tsNodeInspectGroupDigits(text)};return text+"n"}
func tsNodeInspectContainer(parts []string,open,close string,depth int,o tsNodeInspectOptions)string{
	if len(parts)==0{return open+close};joined:=strings.Join(parts,", ");indent:=strings.Repeat(" ",(depth+1)*2);closingIndent:=strings.Repeat(" ",depth*2);if !strings.Contains(joined,"\n")&&len(open)+len(joined)+len(close)+depth*2+2<=o.breakLength{return open+" "+joined+" "+close};return open+"\n"+indent+strings.Join(parts,",\n"+indent)+"\n"+closingIndent+close
}
func tsNodeInspectCustomValue(loop *tsLoop,value tsValue,o tsNodeInspectOptions,depth int,seen map[unsafe.Pointer]bool)(string,bool){
	if !o.customInspect||value.kind!=tsObjectKind||tsInternalInspectState.customSymbol.kind!=tsSymbolKind{return "",false};custom:=tsGet(loop,value,tsStringReference(tsStringUTF8(tsPropertyKey(tsInternalInspectState.customSymbol))));if custom.kind!=tsFunctionKind{return "",false}
	tsInternalInspectState.activeMu.Lock();if tsInternalInspectState.active[value.ref]{tsInternalInspectState.activeMu.Unlock();return "",false};if tsInternalInspectState.active==nil{tsInternalInspectState.active=map[unsafe.Pointer]bool{}};tsInternalInspectState.active[value.ref]=true;tsInternalInspectState.activeMu.Unlock();defer func(){tsInternalInspectState.activeMu.Lock();delete(tsInternalInspectState.active,value.ref);tsInternalInspectState.activeMu.Unlock()}()
	options:=tsNewObject();options.set("depth",tsNumberValue(float64(o.depth-depth)));options.set("colors",tsBooleanValue(o.colors));options.set("stylize",tsNamedNativeMethod("stylize",2,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsStringReference(tsStringUTF8(tsNodeInspectStyle(tsText(tsArg(args,0)),tsText(tsArg(args,1)),o.colors)))}))
	inspect:=tsU;if tsInternalInspectState.module!=nil{inspect=tsGet(loop,tsObjectValue(tsInternalInspectState.module),tsStringReference(tsStringUTF8("inspect")))};result:=tsCallReceiver(loop,custom,value,tsNumberValue(float64(o.depth-depth)),tsObjectValue(options),inspect);if result.kind==value.kind&&result.ref==value.ref{return "",false};if result.kind==tsStringKind{return tsText(result),true};return tsNodeInspectValue(loop,result,o,depth,seen),true
}
func tsNodeInspectStyle(text,style string,colors bool)string{
	if !colors{return text};code:="";switch style{case "special","number","boolean":code="33";case "undefined":code="90";case "null":code="1";case "string","symbol":code="32";case "regexp":code="31"};if code==""{return text};return "\x1b["+code+"m"+text+"\x1b[39m"
}
func tsNodeInspectQuote(text string,limit int)string{
	if limit>=0&&len([]rune(text))>limit{text=string([]rune(text)[:limit])+"..."};quoted:=strconv.Quote(text);if len(quoted)>1{quoted=quoted[1:len(quoted)-1]};quoted=strings.ReplaceAll(quoted,"\\\"","\"");quoted=strings.ReplaceAll(quoted,"'","\\'");return "'"+quoted+"'"
}
func tsNodeInspectError(loop *tsLoop,value tsValue,o tsNodeInspectOptions,depth int,seen map[unsafe.Pointer]bool)string{
	text:=tsText(value);stack:=tsErrorStackGet(loop,value);if stack.kind==tsStringKind&&tsText(stack)!=""{text=tsText(stack)}
	keys:=tsObjectKeys(value);extras:=[]string{};for _,key:=range keys.values{name:=tsPropertyKey(key);if name=="stack"||name=="message"||name=="name"{continue};extras=append(extras,name+": "+tsNodeInspectValue(loop,tsGet(loop,value,tsStringReference(tsStringUTF8(name))),o,depth+1,seen))}
	if len(extras)>0{text+=" { "+strings.Join(extras,", ")+" }"};return tsNodeInspectStyle(text,"special",o.colors)
}
func tsNodeInspectValue(loop *tsLoop,value tsValue,o tsNodeInspectOptions,depth int,seen map[unsafe.Pointer]bool)string{
	style:=func(s,k string)string{return tsNodeInspectStyle(s,k,o.colors)}
	switch value.kind{
	case tsUndefinedKind:return style("undefined","undefined")
	case tsNullKind:return style("null","null")
	case tsBooleanKind:if tsTruthy(value){return style("true","boolean")};return style("false","boolean")
	case tsStringKind:return style(tsNodeInspectQuote(tsText(value),o.maxStringLength),"string")
	case tsSymbolKind:return style(tsSymbolText((*tsSymbol)(value.ref)),"symbol")
	case tsBigIntKind:return style(tsNodeInspectBigIntText((*big.Int)(value.ref),o.numericSeparator),"number")
	case tsNumberKind,tsFloat32Kind,tsIntKind,tsInt8Kind,tsInt16Kind,tsInt32Kind,tsInt64Kind,tsUintKind,tsUint8Kind,tsUint16Kind,tsUint32Kind,tsUint64Kind:return style(tsNodeInspectNumberText(tsNumber(value),o.numericSeparator),"number")
	case tsFunctionKind:name:=(*tsFunction)(value.ref).name;if name!=""{return style("[Function: "+name+"]","special")};return style("[Function]","special")
	case tsClassKind:name:=(*tsClass)(value.ref).builtin;if name==""{name="Function"};return style("[class "+name+"]","special")
	case tsRegExpKind:return style(tsRegExpText((*tsRegExp)(value.ref)).String(),"regexp")
	case tsErrorKind:return tsNodeInspectError(loop,value,o,depth,seen)
	}
	if value.kind==tsObjectKind&&(*tsObject)(value.ref).proxy{proxy:=(*tsObject)(value.ref);if o.showProxy{return "Proxy ["+tsNodeInspectValue(loop,proxy.proxyTarget,o,depth+1,seen)+", "+tsNodeInspectValue(loop,proxy.proxyHandler,o,depth+1,seen)+"]"};return tsNodeInspectValue(loop,proxy.proxyTarget,o,depth,seen)}
	if custom,ok:=tsNodeInspectCustomValue(loop,value,o,depth,seen);ok{return custom}
	if depth>o.depth&&o.depth>=0{return "[Object]"}
	if value.ref!=nil{if seen[value.ref]{return style("[Circular]","special")};seen[value.ref]=true;defer delete(seen,value.ref)}
	if value.kind==tsArrayKind{
		a:=(*tsArray)(value.ref);n:=a.length();if n>o.maxArrayLength{n=o.maxArrayLength};items:=[]string{}
		for i:=0;i<n;i++{if a.holes[i]{items=append(items,style("<empty>","undefined"))}else{items=append(items,tsNodeInspectValue(loop,a.at(i),o,depth+1,seen))}}
		if n<a.length(){items=append(items,fmt.Sprintf("... %d more items",a.length()-n))};return tsNodeInspectContainer(items,"[","]",depth,o)
	}
	if value.kind==tsCollectionKind{
		c:=(*tsCollection)(value.ref);parts:=[]string{};count:=0
		for _,entry:=range c.entries{if !entry.live{continue};if count>=o.maxArrayLength{break};if c.set{parts=append(parts,tsNodeInspectValue(loop,entry.key,o,depth+1,seen))}else{parts=append(parts,tsNodeInspectValue(loop,entry.key,o,depth+1,seen)+" => "+tsNodeInspectValue(loop,entry.value,o,depth+1,seen))};count++}
		if count<len(c.index){parts=append(parts,fmt.Sprintf("... %d more items",len(c.index)-count))};if c.set{return tsNodeInspectContainer(parts,fmt.Sprintf("Set(%d) {",len(c.index)),"}",depth,o)};return tsNodeInspectContainer(parts,fmt.Sprintf("Map(%d) {",len(c.index)),"}",depth,o)
	}
	if value.kind==tsTypedArrayKind{
		a:=(*tsTypedArray)(value.ref);parts:=[]string{};n:=a.length;if n>o.maxArrayLength{n=o.maxArrayLength};for i:=0;i<n;i++{item:=a.readValue(i);if !a.bigint{item=tsNumberValue(a.read(i))};parts=append(parts,tsNodeInspectValue(loop,item,o,depth+1,seen))};if n<a.length{parts=append(parts,fmt.Sprintf("... %d more items",a.length-n))};return tsNodeInspectContainer(parts,fmt.Sprintf("%s(%d) [",a.name,a.length),"]",depth,o)
	}
	if value.kind==tsArrayBufferKind{return fmt.Sprintf("ArrayBuffer { byteLength: %d }",len((*tsArrayBuffer)(value.ref).data))}
	if value.kind==tsPromiseKind{promise:=(*tsPromise)(value.ref);if !promise.settled{return "Promise { <pending> }"};result:=tsNodeInspectValue(loop,promise.result.value,o,depth+1,seen);if promise.result.rejected{return "Promise { <rejected> "+result+" }"};return "Promise { "+result+" }"}
	if value.kind==tsObjectKind&&(*tsObject)(value.ref).nativeDate{return "Date {}"}
	keys:=tsObjectKeys(value);names:=[]string{};for _,key:=range keys.values{names=append(names,tsPropertyKey(key))}
	if o.showHidden&&value.kind==tsObjectKind{object:=(*tsObject)(value.ref);known:=map[string]bool{};for _,name:=range names{known[name]=true};for _,name:=range object.order{if !known[name]{names=append(names,name)}}}
	if o.sorted{sort.Strings(names)}
	parts:=[]string{}
	for _,name:=range names{
		key:=name;valid:=name!="";for i,r:=range name{if !(r=='_'||r=='$'||unicode.IsLetter(r)||i>0&&unicode.IsDigit(r)){valid=false;break}};if !valid{key=tsNodeInspectQuote(name,-1)}
		descriptor:=(*tsDescriptor)(nil);if value.kind==tsObjectKind{descriptor=(*tsObject)(value.ref).descriptors[name]};if o.showHidden&&descriptor!=nil&&!descriptor.enumerable{key="["+key+"]"}
		var property string
		if descriptor!=nil&&descriptor.accessor{if !tsIsUndefined(descriptor.getter){if o.getters{property=tsNodeInspectValue(loop,tsCallReceiver(loop,descriptor.getter,value),o,depth+1,seen)}else if !tsIsUndefined(descriptor.setter){property=style("[Getter/Setter]","special")}else{property=style("[Getter]","special")}}else{property=style("[Setter]","special")}}else{property=tsNodeInspectValue(loop,tsGet(loop,value,tsStringReference(tsStringUTF8(name))),o,depth+1,seen)}
		parts=append(parts,key+": "+property)
	}
	return tsNodeInspectContainer(parts,"{","}",depth,o)
}
func tsNodeInspect(loop *tsLoop,args []tsValue)tsValue{return tsStringReference(tsStringUTF8(tsNodeInspectValue(loop,tsArg(args,0),tsNodeInspectOptionsFrom(loop,args),0,map[unsafe.Pointer]bool{})))}
func tsNodeInspectParseFloat(loop *tsLoop,value tsValue)float64{
	if value.kind==tsSymbolKind{return math.NaN()};text:=strings.TrimSpace(tsStringValueWithLoop(loop,value).String());if text==""{return math.NaN()};for end:=len(text);end>0;end--{if number,err:=strconv.ParseFloat(text[:end],64);err==nil{return number}};return math.NaN()
}
func tsNodeInspectParseInt(loop *tsLoop,value tsValue)float64{
	if value.kind==tsSymbolKind{return math.NaN()};text:=strings.TrimSpace(tsStringValueWithLoop(loop,value).String());if text==""{return math.NaN()};sign:="";if text[0]=='+'||text[0]=='-'{sign=text[:1];text=text[1:]};base:=10;if strings.HasPrefix(text,"0x")||strings.HasPrefix(text,"0X"){base=16;text=text[2:]};end:=0;for end<len(text){c:=text[end];if c>='0'&&c<='9'&&int(c-'0')<base||base==16&&(c>='a'&&c<='f'||c>='A'&&c<='F'){end++}else{break}};if end==0{return math.NaN()};integer,ok:=new(big.Int).SetString(sign+text[:end],base);if !ok{return math.NaN()};number,_:=new(big.Float).SetInt(integer).Float64();return number
}
func tsNodeInspectFormatValue(loop *tsLoop,value tsValue,o tsNodeInspectOptions,spec byte)string{
	switch spec{
	case 's':
		if tsIsNumeric(value){return tsNodeInspectNumberText(tsNumber(value),o.numericSeparator)};if value.kind==tsBigIntKind{return tsNodeInspectBigIntText((*big.Int)(value.ref),o.numericSeparator)}
		if value.kind==tsObjectKind||value.kind==tsArrayKind||value.kind==tsCollectionKind||value.kind==tsTypedArrayKind{return tsNodeInspectValue(loop,value,tsNodeInspectOptions{depth:0,maxArrayLength:o.maxArrayLength,maxStringLength:o.maxStringLength,breakLength:o.breakLength,compact:3,customInspect:o.customInspect},0,map[unsafe.Pointer]bool{})}
		return tsStringValueWithLoop(loop,value).String()
	case 'd':
		if value.kind==tsBigIntKind{return tsNodeInspectBigIntText((*big.Int)(value.ref),o.numericSeparator)};if value.kind==tsSymbolKind{return "NaN"};return tsNodeInspectNumberText(tsNumber(tsPrimitiveHint(loop,value,false)),o.numericSeparator)
	case 'i':
		if value.kind==tsBigIntKind{return tsNodeInspectBigIntText((*big.Int)(value.ref),o.numericSeparator)};return tsNodeInspectNumberText(tsNodeInspectParseInt(loop,value),o.numericSeparator)
	case 'f':return tsNodeInspectNumberText(tsNodeInspectParseFloat(loop,value),o.numericSeparator)
	case 'j':
		var result tsValue;func(){defer func(){if recover()!=nil{result=tsStringReference(tsStringUTF8("[Circular]"))}}();result=tsJSONBuiltin(loop,"stringify",[]tsValue{value})}();if tsIsUndefined(result){return "undefined"};return tsText(result)
	case 'o':
		inner:=o;inner.depth=4;inner.showHidden=true;inner.showProxy=true;return tsNodeInspectValue(loop,value,inner,0,map[unsafe.Pointer]bool{})
	case 'O':return tsNodeInspectValue(loop,value,o,0,map[unsafe.Pointer]bool{})
	};return ""
}
func tsNodeInspectFormat(loop *tsLoop,o tsNodeInspectOptions,args []tsValue)tsValue{
	if len(args)==0{return tsStringReference(tsStringUTF8(""))}
	if args[0].kind!=tsStringKind{parts:=[]string{};for _,value:=range args{if value.kind==tsStringKind{parts=append(parts,tsText(value))}else{parts=append(parts,tsNodeInspectValue(loop,value,o,0,map[unsafe.Pointer]bool{}))}};return tsStringReference(tsStringUTF8(strings.Join(parts," ")))}
	format:=tsText(args[0]);if len(args)==1{return tsStringReference(tsStringUTF8(format))};var out strings.Builder;argument:=1;lastPosition:=0
	for i:=0;i<len(format)-1;i++{if format[i]!='%'{continue};specifier:=format[i+1];if specifier=='%'&&argument>=len(args){continue};if argument>=len(args){continue};switch specifier{case 's','d','i','f','j','o','O':out.WriteString(format[lastPosition:i]);out.WriteString(tsNodeInspectFormatValue(loop,args[argument],o,specifier));argument++;i++;lastPosition=i+1;case 'c':out.WriteString(format[lastPosition:i]);argument++;i++;lastPosition=i+1;case '%':out.WriteString(format[lastPosition:i]);i++;lastPosition=i+1}}
	if lastPosition>0{out.WriteString(format[lastPosition:])}else{out.WriteString(format)}
	for ;argument<len(args);argument++{out.WriteByte(' ');value:=args[argument];if value.kind==tsStringKind{out.WriteString(tsText(value))}else{out.WriteString(tsNodeInspectValue(loop,value,o,0,map[unsafe.Pointer]bool{}))}}
	return tsStringReference(tsStringUTF8(out.String()))
}
func tsNodeStripVT(text string)string{
	runes:=[]rune(text);var out strings.Builder
	for i:=0;i<len(runes);{
		start:=i;next:=i+1;complete:=false
		if runes[i]==0x1b&&next<len(runes){
			switch runes[next]{
			case '[':next++;for next<len(runes){r:=runes[next];next++;if r>=0x40&&r<=0x7e{complete=true;break}}
			case ']':next++;for next<len(runes){if runes[next]==0x07||runes[next]==0x9c{next++;complete=true;break};if runes[next]==0x1b&&next+1<len(runes)&&runes[next+1]=='\\'{next+=2;complete=true;break};next++}
			default:if runes[next]>=0x40&&runes[next]<=0x5f{next++;complete=true}
			}
		}else if runes[i]==0x9b{next++;for next<len(runes){r:=runes[next];next++;if r>=0x40&&r<=0x7e{complete=true;break}}}
		if complete{i=next;continue};out.WriteString(string(runes[start]));i=start+1
	}
	return out.String()
}
func tsNodeInspectStripVT(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{value:=tsArg(args,0);if value.kind!=tsStringKind{panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_TYPE","The \"str\" argument must be of type string")})};return tsStringReference(tsStringUTF8(tsNodeStripVT(tsText(value))))}
func tsNodeInspectStringWidth(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{
	value:=tsArg(args,0);if value.kind!=tsStringKind{panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_TYPE","The \"str\" argument must be of type string")})};text:=tsText(value);options:=tsArg(args,1);if options.kind==tsObjectKind||options.kind==tsArrayKind{if !tsTruthy(tsGet(loop,options,tsStringReference(tsStringUTF8("countAnsiEscapeCodes")))){text=tsNodeStripVT(text)}}else if len(args)<2||tsTruthy(options){text=tsNodeStripVT(text)};normalized:=tsStringMethod(loop,(*tsString)(value.ref),"normalize",[]tsValue{tsStringReference(tsStringUTF8("NFC"))});if normalized.kind==tsStringKind{text=tsText(normalized)}
	width:=0;for _,r:=range text{if r==0x00ad{width++;continue};if unicode.Is(unicode.Cc,r)||unicode.Is(unicode.Cf,r)||unicode.Is(unicode.Me,r)||unicode.Is(unicode.Mn,r)||r>=0x1f3fb&&r<=0x1f3ff{continue};if tsNodeInspectWideCodepoint(r)||tsNodeInspectEmojiPresentation(r){width+=2}else{width++}};return tsNumberValue(float64(width))
}
func tsNodeInspectWideCodepoint(r rune)bool{return r>=0x1100&&r<=0x115f||r==0x2329||r==0x232a||r>=0x2e80&&r<=0xa4cf&&r!=0x303f||r>=0xac00&&r<=0xd7a3||r>=0xf900&&r<=0xfaff||r>=0xfe10&&r<=0xfe19||r>=0xfe30&&r<=0xfe6f||r>=0xff01&&r<=0xff60||r>=0xffe0&&r<=0xffe6||r>=0x20000&&r<=0x3fffd}
func tsNodeInspectEmojiPresentation(r rune)bool{return r>=0x1f1e6&&r<=0x1f1ff||r>=0x1f300&&r<=0x1faff||r>=0x231a&&r<=0x231b||r>=0x23e9&&r<=0x23ec||r==0x23f0||r==0x23f3||r>=0x25fd&&r<=0x25fe||r>=0x2614&&r<=0x2615||r>=0x2648&&r<=0x2653||r==0x267f||r==0x2693||r==0x26a1||r>=0x26aa&&r<=0x26ab||r>=0x26bd&&r<=0x26be||r>=0x26c4&&r<=0x26c5||r==0x26ce||r==0x26d4||r==0x26ea||r>=0x26f2&&r<=0x26f3||r==0x26f5||r==0x26fa||r==0x26fd||r==0x2705||r>=0x270a&&r<=0x270b||r==0x2728||r==0x274c||r==0x274e||r>=0x2753&&r<=0x2755||r==0x2757||r>=0x2795&&r<=0x2797||r==0x27b0||r==0x27bf||r>=0x2b1b&&r<=0x2b1c||r==0x2b50||r==0x2b55}
func tsNodeInspectPromiseDetails(value tsValue)tsValue{out:=&tsArray{};if value.kind!=tsPromiseKind{out.values=append(out.values,tsNumberValue(0),tsU);return tsArrayValue(out)};promise:=(*tsPromise)(value.ref);if !promise.settled{out.values=append(out.values,tsNumberValue(0),tsU)}else if promise.result.rejected{out.values=append(out.values,tsNumberValue(2),promise.result.value)}else{out.values=append(out.values,tsNumberValue(1),promise.result.value)};return tsArrayValue(out)}
func tsNodeInspectIsIndex(name string)bool{index,err:=strconv.ParseUint(name,10,32);return err==nil&&index<0xffffffff&&strconv.FormatUint(index,10)==name}
func tsNodeInspectOwnKeys(loop *tsLoop,value tsValue,filter int)tsValue{
	if tsNullish(value){tsPropertyFailure("Cannot convert undefined or null to object")};if filter&(1<<3)!=0&&filter&(1<<4)!=0{return tsArrayValue(&tsArray{})}
	if value.kind==tsStringKind{out:=&tsArray{};if filter&(1<<3)==0&&filter&(1<<1)==0{out.values=append(out.values,tsStringReference(tsStringUTF8("length")))};return tsArrayValue(out)}
	if value.kind!=tsObjectKind&&value.kind!=tsFunctionKind&&value.kind!=tsInstanceKind&&value.kind!=tsClassKind&&value.kind!=tsArrayKind&&value.kind!=tsRegExpKind&&value.kind!=tsTypedArrayKind&&value.kind!=tsDataViewKind&&value.kind!=tsCollectionKind{return tsArrayValue(&tsArray{})}
	keys:=[]string{}
	switch value.kind{
	case tsObjectKind:object:=(*tsObject)(value.ref);if object.proxy{value=object.proxyTarget;return tsNodeInspectOwnKeys(loop,value,filter)};keys=append(keys,object.order...);if object.boxedPrimitive.kind==tsStringKind{keys=append(keys,"length")}
	case tsFunctionKind:keys=append(keys,(*tsFunction)(value.ref).metadata().order...)
	case tsInstanceKind:keys=append(keys,(*tsProperties)(value.ref).order...)
	case tsClassKind:keys=append(keys,tsInstanceProperties((*tsClass)(value.ref).static).order...)
	case tsArrayKind:array:=(*tsArray)(value.ref);for i:=0;i<array.length();i++{if !array.holes[i]{keys=append(keys,strconv.Itoa(i))}};keys=append(keys,"length");keys=append(keys,array.order...);for key:=range array.properties{keys=append(keys,key)}
	case tsRegExpKind:keys=append(keys,tsRegExpMetadata((*tsRegExp)(value.ref)).order...)
	case tsTypedArrayKind:keys=append(keys,(*tsTypedArray)(value.ref).order...)
	case tsDataViewKind:keys=append(keys,(*tsDataView)(value.ref).order...)
	case tsCollectionKind:keys=append(keys,(*tsCollection)(value.ref).order...)
	}
	out:=&tsArray{};for _,name:=range tsOrderedKeys(keys){symbol:=tsIsSymbolKey(name);if tsNodeInspectIsIndex(name){continue};if symbol&&filter&(1<<4)!=0||!symbol&&filter&(1<<3)!=0{continue};if filter&(1<<1)!=0&&value.kind!=tsTypedArrayKind&&value.kind!=tsDataViewKind&&value.kind!=tsCollectionKind{descriptor:=tsOwnDescriptor(loop,value,name);if tsIsUndefined(descriptor)||!tsTruthy(tsGet(loop,descriptor,tsStringReference(tsStringUTF8("enumerable")))){continue}};out.values=append(out.values,tsKeyValue(name))};return tsArrayValue(out)
}
func tsInternalInspectModule(loop *tsLoop)tsValue{
	tsInternalInspectState.Do(func(){module:=tsNewObject();symbolFor:=tsGet(loop,tsSymbolFunction(),tsStringReference(tsStringUTF8("for")));tsInternalInspectState.customSymbol=tsCall(loop,symbolFor,tsStringReference(tsStringUTF8("nodejs.util.inspect.custom")))
		defaults:=tsNewObject();for name,v:=range map[string]tsValue{"showHidden":tsBooleanValue(false),"depth":tsNumberValue(2),"colors":tsBooleanValue(false),"customInspect":tsBooleanValue(true),"showProxy":tsBooleanValue(false),"maxArrayLength":tsNumberValue(100),"maxStringLength":tsNumberValue(10000),"breakLength":tsNumberValue(80),"compact":tsNumberValue(3),"sorted":tsBooleanValue(false),"getters":tsBooleanValue(false),"numericSeparator":tsBooleanValue(false)}{defaults.set(name,v)};tsInternalInspectState.defaults=defaults
		inspect:=tsNamedNativeMethod("inspect",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeInspect(loop,args)});fn:=(*tsFunction)(inspect.ref);fn.properties.set("custom",tsInternalInspectState.customSymbol);fn.properties.set("defaultOptions",tsObjectValue(defaults));fn.properties.set("replDefaults",tsObjectValue(defaults))
		colors:=tsNewObject();colorPairs:=[]struct{name string;open,close int}{{"reset",0,0},{"bold",1,22},{"dim",2,22},{"italic",3,23},{"underline",4,24},{"blink",5,25},{"inverse",7,27},{"hidden",8,28},{"strikethrough",9,29},{"doubleunderline",21,24},{"black",30,39},{"red",31,39},{"green",32,39},{"yellow",33,39},{"blue",34,39},{"magenta",35,39},{"cyan",36,39},{"white",37,39},{"bgBlack",40,49},{"bgRed",41,49},{"bgGreen",42,49},{"bgYellow",43,49},{"bgBlue",44,49},{"bgMagenta",45,49},{"bgCyan",46,49},{"bgWhite",47,49},{"framed",51,54},{"overlined",53,55},{"gray",90,39},{"redBright",91,39},{"greenBright",92,39},{"yellowBright",93,39},{"blueBright",94,39},{"magentaBright",95,39},{"cyanBright",96,39},{"whiteBright",97,39},{"bgGray",100,49},{"bgRedBright",101,49},{"bgGreenBright",102,49},{"bgYellowBright",103,49},{"bgBlueBright",104,49},{"bgMagentaBright",105,49},{"bgCyanBright",106,49},{"bgWhiteBright",107,49}};for _,entry:=range colorPairs{colors.set(entry.name,tsArrayValue(&tsArray{values:[]tsValue{tsNumberValue(float64(entry.open)),tsNumberValue(float64(entry.close))}}))};for alias,target:=range map[string]string{"grey":"gray","blackBright":"gray","bgGrey":"bgGray","bgBlackBright":"bgGray","faint":"dim","crossedout":"strikethrough","strikeThrough":"strikethrough","crossedOut":"strikethrough","conceal":"hidden","swapColors":"inverse","swapcolors":"inverse","doubleUnderline":"doubleunderline"}{colors.set(alias,colors.values[target])};fn.properties.set("colors",tsObjectValue(colors))
		styles:=tsNewObject();for name,color:=range map[string]string{"special":"cyan","number":"yellow","bigint":"yellow","boolean":"yellow","undefined":"grey","null":"bold","string":"green","symbol":"green","date":"magenta","module":"underline"}{styles.set(name,tsStringReference(tsStringUTF8(color)))};styles.set("regexp",tsStringReference(tsStringUTF8("green")));fn.properties.set("styles",tsObjectValue(styles));module.set("inspect",inspect)
		module.set("format",tsNamedNativeMethod("format",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeInspectFormat(loop,tsNodeInspectOptionsFrom(loop,nil),args)}))
		module.set("formatWithOptions",tsNamedNativeMethod("formatWithOptions",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{options:=tsArg(args,0);if options.kind!=tsObjectKind&&options.kind!=tsInstanceKind&&options.kind!=tsArrayKind{panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_TYPE","The \"inspectOptions\" argument must be of type object")})};rest:=[]tsValue{};if len(args)>1{rest=args[1:]};return tsNodeInspectFormat(loop,tsNodeInspectOptionsFrom(loop,[]tsValue{tsU,options}),rest)}))
		module.set("getOwnNonIndexProperties",tsNamedNativeMethod("getOwnNonIndexProperties",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{filter:=0;if !tsIsUndefined(tsArg(args,1)){filter=int(tsNumber(tsArg(args,1)))};return tsNodeInspectOwnKeys(loop,tsArg(args,0),filter)}))
		module.set("stripVTControlCharacters",tsNamedNativeMethod("stripVTControlCharacters",1,tsNodeInspectStripVT));module.set("getStringWidth",tsNamedNativeMethod("getStringWidth",2,tsNodeInspectStringWidth));module.set("getPromiseDetails",tsNamedNativeMethod("getPromiseDetails",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeInspectPromiseDetails(tsArg(args,0))}))
		module.set("getProxyDetails",tsNamedNativeMethod("getProxyDetails",2,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{value:=tsArg(args,0);if value.kind!=tsObjectKind||!(*tsObject)(value.ref).proxy{return tsU};object:=(*tsObject)(value.ref);if len(args)>1&&!tsTruthy(args[1]){return object.proxyTarget};out:=&tsArray{};out.values=append(out.values,object.proxyTarget,object.proxyHandler);return tsArrayValue(out)}))
		module.set("previewEntries",tsNamedNativeMethod("previewEntries",2,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{value:=tsArg(args,0);if value.kind!=tsCollectionKind{return tsArrayValue(&tsArray{})};collection:=(*tsCollection)(value.ref);items:=&tsArray{};entries:=!collection.set;for _,entry:=range collection.entries{if entry.live{if collection.set{items.values=append(items.values,entry.key)}else{pair:=&tsArray{};pair.values=append(pair.values,entry.key,entry.value);items.values=append(items.values,tsArrayValue(pair))}}};out:=&tsArray{};out.values=append(out.values,tsArrayValue(items),tsBooleanValue(entries));return tsArrayValue(out)}))
		module.set("internalGetConstructorName",tsNamedNativeMethod("internalGetConstructorName",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{value:=tsArg(args,0);if value.kind==tsObjectKind&&(*tsObject)(value.ref).nativeClass!=nil{return tsStringReference(tsStringUTF8((*tsObject)(value.ref).nativeClass.builtin))};tag:=tsText(tsObjectTag(loop,value));if strings.HasPrefix(tag,"[object "){tag=strings.TrimSuffix(strings.TrimPrefix(tag,"[object "),"]")};return tsStringReference(tsStringUTF8(tag))}));tsInternalInspectState.module=module
	});return tsObjectValue(tsInternalInspectState.module)
}
`
