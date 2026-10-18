package goemit

// URI encoding and decoding operate on the runtime's UTF-16 strings. Dynamic
// source evaluation is not involved; these are the standard ECMAScript URI
// conversion functions used by JavaScript packages.
const uriRuntime = `
func tsURIMalformed() {
	panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"URIError",message:"URI malformed"})})
}

func tsURIEncodeAllowed(value rune, component bool) bool {
	if value >= 'a' && value <= 'z' || value >= 'A' && value <= 'Z' || value >= '0' && value <= '9' || strings.ContainsRune("-_.!~*'()", value) {
		return true
	}
	return !component && strings.ContainsRune(";/?:@&=+$,#", value)
}

func tsURIEncode(loop *tsLoop, args []tsValue, component bool) tsValue {
	text := tsStringValueWithLoop(loop, tsArg(args,0))
	const hex = "0123456789ABCDEF"
	var out strings.Builder
	for index:=0; index<len(text.units); index++ {
		unit:=text.units[index]
		value:=rune(unit)
		if unit>=0xD800 && unit<=0xDBFF {
			if index+1>=len(text.units) || text.units[index+1]<0xDC00 || text.units[index+1]>0xDFFF { tsURIMalformed() }
			value=utf16.DecodeRune(rune(unit),rune(text.units[index+1]));index++
		} else if unit>=0xDC00 && unit<=0xDFFF { tsURIMalformed() }
		if tsURIEncodeAllowed(value,component) { out.WriteRune(value);continue }
		for _,octet:=range []byte(string(value)) { out.WriteByte('%');out.WriteByte(hex[octet>>4]);out.WriteByte(hex[octet&15]) }
	}
	return tsStringReference(tsStringUTF8(out.String()))
}

func tsURIHex(units []uint16,index int)(byte,bool){
	if index+2>=len(units)||units[index]!='%'{return 0,false}
	value:=byte(0)
	for _,unit:=range units[index+1:index+3]{value<<=4;switch{case unit>='0'&&unit<='9':value|=byte(unit-'0');case unit>='a'&&unit<='f':value|=byte(unit-'a'+10);case unit>='A'&&unit<='F':value|=byte(unit-'A'+10);default:return 0,false}}
	return value,true
}

func tsURIReserved(value byte)bool{return strings.ContainsRune(";/?:@&=+$,#",rune(value))}

func tsURIDecode(loop *tsLoop,args []tsValue,component bool)tsValue{
	text:=tsStringValueWithLoop(loop,tsArg(args,0));units:=text.units;out:=make([]uint16,0,len(units))
	for index:=0;index<len(units);{
		if units[index]!='%'{out=append(out,units[index]);index++;continue}
		first,ok:=tsURIHex(units,index);if !ok{tsURIMalformed()}
		width:=1
		switch{case first<0x80:width=1;case first>=0xC2&&first<=0xDF:width=2;case first>=0xE0&&first<=0xEF:width=3;case first>=0xF0&&first<=0xF4:width=4;default:tsURIMalformed()}
		if !component&&width==1&&tsURIReserved(first){out=append(out,units[index:index+3]...);index+=3;continue}
		encoded:=make([]byte,width);encoded[0]=first
		for offset:=1;offset<width;offset++{next,valid:=tsURIHex(units,index+offset*3);if !valid||next<0x80||next>0xBF{tsURIMalformed()};encoded[offset]=next}
		value,size:=utf8.DecodeRune(encoded);if size!=width{tsURIMalformed()}
		if value<=0xFFFF{out=append(out,uint16(value))}else{high,low:=utf16.EncodeRune(value);out=append(out,uint16(high),uint16(low))}
		index+=width*3
	}
	return tsStringReference(tsStringUnits(out))
}
`
