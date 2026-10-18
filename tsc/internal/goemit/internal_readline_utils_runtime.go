package goemit

// internalReadlineUtilsRuntime adapts the cursor, string and keypress helpers
// consumed by Node's readline implementation.
const internalReadlineUtilsRuntime = `
var tsInternalReadlineUtils struct { sync.Once; module tsValue }

func tsInternalReadlineUtilsProperty(value tsValue, name string) tsValue { return tsGet(value, tsStringReference(tsStringUTF8(name))) }
func tsInternalReadlineString(value tsValue) *tsString { return tsStringValue(value) }
func tsInternalReadlineText(value tsValue) string { return tsText(value) }
func tsInternalReadlineCSIValue(_ *tsLoop, parts tsValue, args []tsValue) tsValue {
	var out strings.Builder
	out.WriteString("\x1b[")
	if parts.kind == tsArrayKind {
		array := (*tsArray)(parts.ref)
		for index := 0; index < array.length(); index++ {
			out.WriteString(tsText(array.at(index)))
			if index < len(args) { out.WriteString(tsText(args[index])) }
		}
	}
	return tsStringReference(tsStringUTF8(out.String()))
}
func tsInternalReadlineCharLengthLeft(_ *tsLoop, text *tsString, index int) int {
	if index <= 0 { return 0 }
	if index > 1 && index-2 < len(text.units) && text.units[index-2] >= 0xd800 && text.units[index-2] <= 0xdbff { return 2 }
	if index-1 < len(text.units) && text.units[index-1] >= 0xd800 && text.units[index-1] <= 0xdbff { return 2 }
	return 1
}
func tsInternalReadlineCharLengthAt(_ *tsLoop, text *tsString, index int) int {
	if index < 0 || index >= len(text.units) { return 1 }
	if text.units[index] >= 0xd800 && text.units[index] <= 0xdbff { return 2 }
	return 1
}
func tsInternalReadlinePrefix(loop *tsLoop, value tsValue) tsValue {
	items := []tsValue{}
	if value.kind == tsArrayKind { items = append(items, (*tsArray)(value.ref).values...) }
	if len(items) == 0 { return tsStringReference(tsStringUTF8("")) }
	minimum := tsStringValue(items[0])
	for _, item := range items[1:] {
		current := tsStringValue(item)
		limit := len(minimum.units); if len(current.units) < limit { limit = len(current.units) }
		index := 0; for index < limit && minimum.units[index] == current.units[index] { index++ }
		minimum = minimum.slice(0,index)
	}
	return tsStringReference(minimum)
}
func tsInternalReadlineReverse(loop *tsLoop, line, from, to tsValue) tsValue {
	text, separator, replacement := tsText(line), tsText(from), tsText(to)
	if separator == "" {
		units := tsStringValue(line).units
		parts := make([]string, len(units)); for i, unit := range units { parts[i] = string(rune(unit)) }
		for left, right := 0, len(parts)-1; left < right; left, right = left+1, right-1 { parts[left],parts[right]=parts[right],parts[left] }
		return tsStringReference(tsStringUTF8(strings.Join(parts,replacement)))
	}
	parts := strings.Split(text,separator)
	var out strings.Builder
	for index:=len(parts)-1; index>0; index-- { out.WriteString(parts[index]);out.WriteString(replacement) }
	if len(parts)>0 { out.WriteString(parts[0]) }
	return tsStringReference(tsStringUTF8(out.String()))
}

type tsInternalReadlineKeyDecoder struct { stream tsValue; pending string }

func tsInternalReadlineKeyObject(_ *tsLoop, sequence string, name tsValue, ctrl, meta, shift bool) tsValue {
	key:=tsNewObject();key.set("sequence",tsStringReference(tsStringUTF8(sequence)));key.set("name",name);key.set("ctrl",tsBooleanValue(ctrl));key.set("meta",tsBooleanValue(meta));key.set("shift",tsBooleanValue(shift));return tsObjectValue(key)
}
func tsInternalReadlineEmitKey(loop *tsLoop, stream tsValue, sequence string, name tsValue, ctrl,meta,shift,escaped bool) {
	if stream.kind==tsObjectKind || stream.kind==tsInstanceKind || stream.kind==tsFunctionKind {
		emit:=tsGet(loop,stream,tsStringReference(tsStringUTF8("emit")))
		if emit.kind==tsFunctionKind { char:=tsStringReference(tsStringUTF8(sequence));if escaped { char=tsU };tsCallReceiver(loop,emit,stream,tsStringReference(tsStringUTF8("keypress")),char,tsInternalReadlineKeyObject(loop,sequence,name,ctrl,meta,shift)) }
	}
}
func tsInternalReadlineNamedKey(code string) (string,bool,bool) {
	keys:=map[string]string{
		"[P":"f1","[Q":"f2","[R":"f3","[S":"f4","OP":"f1","OQ":"f2","OR":"f3","OS":"f4",
		"[11~":"f1","[12~":"f2","[13~":"f3","[14~":"f4","[200~":"paste-start","[201~":"paste-end",
		"[[A":"f1","[[B":"f2","[[C":"f3","[[D":"f4","[[E":"f5","[15~":"f5","[17~":"f6","[18~":"f7","[19~":"f8","[20~":"f9","[21~":"f10","[23~":"f11","[24~":"f12",
		"[A":"up","[B":"down","[C":"right","[D":"left","[E":"clear","[F":"end","[H":"home","OA":"up","OB":"down","OC":"right","OD":"left","OE":"clear","OF":"end","OH":"home",
		"[1~":"home","[2~":"insert","[3~":"delete","[4~":"end","[5~":"pageup","[6~":"pagedown","[[5~":"pageup","[[6~":"pagedown","[7~":"home","[8~":"end",
		"[a":"up","[b":"down","[c":"right","[d":"left","[e":"clear","[2$":"insert","[3$":"delete","[5$":"pageup","[6$":"pagedown","[7$":"home","[8$":"end",
		"Oa":"up","Ob":"down","Oc":"right","Od":"left","Oe":"clear","[2^":"insert","[3^":"delete","[5^":"pageup","[6^":"pagedown","[7^":"home","[8^":"end","[Z":"tab",
	}
	name:=keys[code];shift:=strings.HasSuffix(code,"$")||code=="[Z"||len(code)==2&&code[0]=='['&&code[1]>='a'&&code[1]<='e';ctrl:=len(code)==2&&code[0]=='O'&&code[1]>='a'&&code[1]<='e'
	if name=="" { return "undefined",false,false }
	return name,ctrl,shift||code=="[Z"
}
func tsInternalReadlineEmitCharacter(loop *tsLoop, decoder *tsInternalReadlineKeyDecoder, char string) {
	if char=="" { return }
	if decoder.pending=="" {
		if char=="\x1b" { decoder.pending=char;return }
		name:=tsU;ctrl,meta,shift:=false,false,false
		switch char {case "\r":name=tsStringReference(tsStringUTF8("return"));case "\n":name=tsStringReference(tsStringUTF8("enter"));case "\t":name=tsStringReference(tsStringUTF8("tab"));case "\b","\x7f":name=tsStringReference(tsStringUTF8("backspace"));case " ":name=tsStringReference(tsStringUTF8("space"));default:
			if len(char)==1&&char[0]<=0x1a { name=tsStringReference(tsStringUTF8(string(rune(char[0]+'a'-1))));ctrl=true
			} else if len(char)==1&&((char[0]>='0'&&char[0]<='9')||(char[0]>='A'&&char[0]<='Z')||(char[0]>='a'&&char[0]<='z')) { lower:=strings.ToLower(char);name=tsStringReference(tsStringUTF8(lower));shift=char>="A"&&char<="Z" }
		}
		tsInternalReadlineEmitKey(loop,decoder.stream,char,name,ctrl,meta,shift,false);return
	}
	decoder.pending+=char
	sequence:=decoder.pending
	if sequence=="\x1b\x1b" { tsInternalReadlineEmitKey(loop,decoder.stream,sequence,tsStringReference(tsStringUTF8("escape")),false,true,false,true);decoder.pending="";return }
	if sequence=="\x1b" { return }
	if len(sequence)==2 && sequence[1]!='[' && sequence[1]!='O' {
		value:=sequence[1:];if len(value)==1&&((value[0]>='A'&&value[0]<='Z')||(value[0]>='a'&&value[0]<='z')||(value[0]>='0'&&value[0]<='9')) { tsInternalReadlineEmitKey(loop,decoder.stream,sequence,tsStringReference(tsStringUTF8(strings.ToLower(value))),false,true,value>="A"&&value<="Z",true) } else { tsInternalReadlineEmitKey(loop,decoder.stream,sequence,tsU,false,true,false,true) };decoder.pending="";return
	}
	if sequence[1]!='['&&sequence[1]!='O' { return }
	if len(sequence)==3 && sequence[1]=='O' && sequence[2]>='0'&&sequence[2]<='9' { return }
	final:=sequence[len(sequence)-1]
	complete:=false
	if len(sequence)>2 {
		if final>='A'&&final<='Z'||final>='a'&&final<='z'||final=='~'||final=='^'||final=='$' { complete=true }
		if sequence[1]=='[' && (strings.HasSuffix(sequence,"~")||strings.HasSuffix(sequence,"^")||strings.HasSuffix(sequence,"$")||final>='A'&&final<='Z'||final>='a'&&final<='z') { complete=true }
	}
	if !complete { return }
	body:=sequence[1:]
	modifier:=0
	if semicolon:=strings.LastIndexByte(body,';');semicolon>=0&&semicolon+1<len(body)&&body[semicolon+1]>='0'&&body[semicolon+1]<='9' { modifier=int(body[semicolon+1]-'1');body=body[:semicolon]+body[semicolon+2:] }
	name,ctrl,shift:=tsInternalReadlineNamedKey(body)
	ctrl=ctrl||modifier&4!=0;meta:=modifier&10!=0;shift=shift||modifier&1!=0
	if body=="[a"||body=="[b"||body=="[c"||body=="[d"||body=="[e" { shift=true }
	if body=="Oa"||body=="Ob"||body=="Oc"||body=="Od"||body=="Oe" { ctrl=true }
	tsInternalReadlineEmitKey(loop,decoder.stream,sequence,tsStringReference(tsStringUTF8(name)),ctrl,meta,shift,true)
	decoder.pending=""
}
func tsInternalReadlineDecoderObject(loop *tsLoop, stream tsValue) tsValue {
	decoder:=&tsInternalReadlineKeyDecoder{stream:stream};object:=tsNewObject()
	object.set("next",tsNamedNativeMethod("next",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{
		if len(args)>0 { text:=tsText(args[0]);if text==""&&decoder.pending!="" { sequence:=decoder.pending;decoder.pending="";if sequence=="\x1b" { tsInternalReadlineEmitKey(loop,stream,sequence,tsStringReference(tsStringUTF8("escape")),false,true,false,true) } else { tsInternalReadlineEmitKey(loop,stream,sequence,tsStringReference(tsStringUTF8("undefined")),false,false,false,true) } } else if text!="" { for _,r:=range text { tsInternalReadlineEmitCharacter(loop,decoder,string(r)) } } }
		return tsIteratorResult(tsU,false)
	}))
	iteratorKey:=tsPropertyKey(tsWellKnownSymbol("iterator"));object.set(iteratorKey,tsNamedNativeMethod("[Symbol.iterator]",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsObjectValue(object)}))
	return tsObjectValue(object)
}

func tsInternalReadlineUtilsModule(loop *tsLoop) tsValue {
	tsInternalReadlineUtils.Do(func(){
		module:=tsNewObject();escape:=tsStringReference(tsStringUTF8("\x1b"));csiFn:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsInternalReadlineCSIValue(loop,tsArg(args,0),args[1:])});csiFn.name="CSI";csiFn.length=1;csiFn.properties=tsNewObject();csi:=tsFunctionValue(csiFn)
		csiFn.properties.set("kEscape",escape);for _,entry:=range []struct{name,value string}{{"kClearToLineBeginning","1K"},{"kClearToLineEnd","0K"},{"kClearLine","2K"},{"kClearScreenDown","0J"}}{csiFn.properties.set(entry.name,tsInternalReadlineCSIValue(loop,tsArrayValue(&tsArray{values:[]tsValue{tsStringReference(tsStringUTF8(entry.value))}}),nil))}
		module.set("CSI",csi);module.set("kSubstringSearch",tsNewSymbol(tsStringUTF8("kSubstringSearch")))
		module.set("charLengthLeft",tsNamedNativeMethod("charLengthLeft",2,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNumberValue(float64(tsInternalReadlineCharLengthLeft(loop,tsStringValue(tsArg(args,0)),int(tsArg(args,1).ToInt32(loop)))))}))
		module.set("charLengthAt",tsNamedNativeMethod("charLengthAt",2,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNumberValue(float64(tsInternalReadlineCharLengthAt(loop,tsStringValue(tsArg(args,0)),int(tsArg(args,1).ToInt32(loop)))))}))
		module.set("commonPrefix",tsNamedNativeMethod("commonPrefix",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsInternalReadlinePrefix(loop,tsArg(args,0))}))
		module.set("reverseString",tsNamedNativeMethod("reverseString",3,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{from,to:=tsArg(args,1),tsArg(args,2);if tsIsUndefined(from){from=tsStringReference(tsStringUTF8("\r"))};if tsIsUndefined(to){to=tsStringReference(tsStringUTF8("\r"))};return tsInternalReadlineReverse(loop,tsArg(args,0),from,to)}))
		module.set("emitKeys",tsNamedNativeMethod("emitKeys",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsInternalReadlineDecoderObject(loop,tsArg(args,0))}))
		module.set("default",tsObjectValue(module));tsInternalReadlineUtils.module=tsObjectValue(module)
	})
	return tsInternalReadlineUtils.module
}
`
