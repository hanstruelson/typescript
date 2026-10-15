package goemit

// NodeBufferRuntime adds the module-level APIs that sit beside Buffer's
// byte-array implementation in fs_encoding.go and fs_buffer_methods.go.
const NodeBufferRuntime = `
func tsNodeBufferModuleTypeError(message string) {
	panic(tsThrown{tsFSException("TypeError", "ERR_INVALID_ARG_TYPE", message)})
}

func tsNodeBufferByteView(value tsValue) []byte {
	switch value.kind {
	case tsArrayBufferKind:
		buffer := (*tsArrayBuffer)(value.ref)
		if buffer.detached { panic(tsThrown{tsFSException("TypeError", "ERR_INVALID_STATE", "Cannot validate on a detached buffer")}) }
		return buffer.data
	case tsTypedArrayKind:
		array := (*tsTypedArray)(value.ref)
		if array.buffer.detached { panic(tsThrown{tsFSException("TypeError", "ERR_INVALID_STATE", "Cannot validate on a detached buffer")}) }
		return array.buffer.data[array.offset:array.offset+array.length*array.bytes]
	case tsDataViewKind:
		view := (*tsDataView)(value.ref)
		if view.buffer.detached { panic(tsThrown{tsFSException("TypeError", "ERR_INVALID_STATE", "Cannot validate on a detached buffer")}) }
		return view.buffer.data[view.offset:view.offset+view.length]
	default:
		tsNodeBufferModuleTypeError("First argument must be an ArrayBufferView")
	}
	return nil
}

func tsNodeBufferVisibleBytes(value tsValue) []byte {
	if value.kind==tsTypedArrayKind&&(*tsTypedArray)(value.ref).buffer.detached {return nil}
	if value.kind==tsDataViewKind&&(*tsDataView)(value.ref).buffer.detached {return nil}
	return tsFSBytes(value)
}

func tsNodeBufferIsAscii(value tsValue) tsValue {
	data := tsNodeBufferByteView(value)
	for _, item := range data { if item > 0x7f { return tsBooleanValue(false) } }
	return tsBooleanValue(true)
}

func tsNodeBufferIsUtf8(value tsValue) tsValue {
	return tsBooleanValue(utf8.Valid(tsNodeBufferByteView(value)))
}

func tsNodeBufferByteLength(_ *tsLoop, args []tsValue) tsValue {
	value := tsArg(args, 0)
	if value.kind == tsStringKind {
		text := (*tsString)(value.ref)
		encoding := "utf8"
		if candidate := tsArg(args, 1); candidate.kind == tsStringKind {
			encoding = strings.ToLower(tsText(candidate))
		}
		switch encoding {
		case "hex", "hex\x00":
			return tsNumberValue(float64(len(text.units) >> 1))
		case "base64", "base64url":
			units := len(text.units)
			if units > 0 && text.units[units-1] == '=' {
				units--
				if units > 0 && text.units[units-1] == '=' {
					units--
				}
			}
			return tsNumberValue(float64((units * 3) >> 2))
		case "latin1", "binary", "ascii":
			return tsNumberValue(float64(len(text.units)))
		case "ucs2", "ucs-2", "utf16le", "utf-16le":
			return tsNumberValue(float64(len(text.units) * 2))
		default:
			// Node treats an unknown encoding in byteLength as UTF-8.
			return tsNumberValue(float64(len([]byte(text.String()))))
		}
	}
	if value.kind == tsArrayBufferKind || value.kind == tsTypedArrayKind || value.kind == tsDataViewKind {
		return tsNumberValue(float64(len(tsNodeBufferVisibleBytes(value))))
	}
	tsFSInvalid("string", "a string, Buffer, or ArrayBuffer")
	return tsU
}

func tsNodeBufferParseEncoding(loop *tsLoop, value tsValue, fallback string, emptyDefaults bool) string {
	if tsNullish(value) { if emptyDefaults { return fallback };value=tsStringReference(tsStringValueWithLoop(loop,value)) }
	if emptyDefaults&&value.kind==tsStringKind&&len((*tsString)(value.ref).units)==0 { return fallback }
	name:=strings.ToLower(tsText(tsStringReference(tsStringValueWithLoop(loop,value))))
	switch name {
	case "utf8","utf-8":return "utf8"
	case "ucs2","ucs-2","utf16le","utf-16le":return "utf16le"
	case "latin1","binary":return "latin1"
	case "ascii","base64","base64url":return name
	case "hex","hex\x00":return "hex"
	default:panic(tsThrown{tsFSException("TypeError","ERR_UNKNOWN_ENCODING","Unknown encoding: "+name)})
	}
}

func tsNodeBufferFill(loop *tsLoop, receiver tsValue, args []tsValue) tsValue {
	data:=tsNodeBufferVisibleBytes(receiver);value:=tsArg(args,0);limit:=len(data)
	offsetValue,endValue,encodingValue:=tsArg(args,1),tsArg(args,2),tsArg(args,3)
	if value.kind==tsStringKind {
		if tsIsUndefined(offsetValue)||offsetValue.kind==tsStringKind { encodingValue=offsetValue;offsetValue=tsU;endValue=tsU
		} else if endValue.kind==tsStringKind { encodingValue=endValue;endValue=tsU }
	}
	encoding:="utf8"
	if value.kind==tsStringKind { encoding=tsNodeBufferParseEncoding(loop,encodingValue,"utf8",true) }
	offset,end:=0,limit
	if !tsIsUndefined(offsetValue) {
		validate:=func(candidate tsValue,name string,max int)int {
			if !tsIsNumeric(candidate) { tsFSInvalid(name,"a number") }
			n:=tsNumber(candidate)
			if math.IsNaN(n)||math.IsInf(n,0)||math.Trunc(n)!=n||n<0||n>float64(max) { panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE",name+" is out of range")}) }
			return int(n)
		}
		offset=validate(offsetValue,"offset",4294967296)
		if !tsIsUndefined(endValue) { end=validate(endValue,"end",limit) }
	}
	if offset>=end { return receiver }
	var pattern []byte
	switch value.kind {
	case tsStringKind:
		pattern=tsFSEncode((*tsString)(value.ref),encoding)
	case tsTypedArrayKind,tsDataViewKind:
		pattern=tsNodeBufferByteView(value)
		if len(pattern)==0 { panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_VALUE","Buffer cannot be empty")}) }
	default:
		byteValue:=byte(tsNumberUint32(tsNumber(tsPrimitiveHint(loop,value,false))))
		pattern=[]byte{byteValue}
	}
	if offset>len(data) { offset=len(data) };if end>len(data) { end=len(data) };if offset>=end { return receiver }
	if len(pattern)==0 { clear(data[offset:end]);return receiver }
	for i:=offset;i<end;i++ { data[i]=pattern[(i-offset)%len(pattern)] }
	return receiver
}

func tsNodeBufferCompare(loop *tsLoop, receiver tsValue, args []tsValue) tsValue {
	if !tsIsUint8BufferLike(receiver) { tsFSInvalid("target","a Buffer or Uint8Array") }
	other:=tsArg(args,0);if !tsIsUint8BufferLike(other) { tsFSInvalid("target","a Buffer or Uint8Array") }
	source,target:=tsNodeBufferVisibleBytes(receiver),tsNodeBufferVisibleBytes(other)
	sourceStart,sourceEnd,targetStart,targetEnd:=0,len(source),0,len(target)
	parse:=func(value tsValue,name string,fallback,max int)int {
		if tsIsUndefined(value) {return fallback};if !tsIsNumeric(value) {tsFSInvalid(name,"a number")};n:=tsNumber(value)
		if math.IsNaN(n)||math.IsInf(n,0)||math.Trunc(n)!=n||n<0||n>4294967296 {panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE",name+" is out of range")})}
		return int(n)
	}
	// The native argument dispatcher validates sourceEnd, sourceStart, targetEnd,
	// then targetStart, even though these values appear in the opposite order.
	sourceEnd=parse(tsArg(args,4),"sourceEnd",sourceEnd,len(source))
	sourceStart=parse(tsArg(args,3),"sourceStart",0,len(source))
	targetEnd=parse(tsArg(args,2),"targetEnd",targetEnd,len(target))
	targetStart=parse(tsArg(args,1),"targetStart",0,len(target))
	if sourceEnd>len(source) {panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","sourceEnd is outside the Buffer")})}
	if targetEnd>len(target) {panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","targetEnd is outside the Buffer")})}
	if sourceStart>=sourceEnd {if targetStart>=targetEnd{return tsNumberValue(0)};return tsNumberValue(-1)}
	if targetStart>=targetEnd {return tsNumberValue(1)}
	sourceRange,targetRange:=source[sourceStart:sourceEnd],target[targetStart:targetEnd]
	return tsNumberValue(float64(bytes.Compare(sourceRange,targetRange)))
}

func tsNodeBufferCopy(loop *tsLoop, receiver tsValue, args []tsValue) tsValue {
	source:=tsNodeBufferVisibleBytes(receiver);targetValue:=tsArg(args,0)
	if targetValue.kind!=tsTypedArrayKind&&targetValue.kind!=tsDataViewKind {tsFSInvalid("target","a Buffer or Uint8Array")}
	target:=tsNodeBufferVisibleBytes(targetValue)
	toOffset:=func(value tsValue)float64 {
		if tsIsUndefined(value) {return 0}
		if tsIsNumeric(value) {n:=tsNumber(value);if !math.IsNaN(n)&&!math.IsInf(n,0)&&math.Trunc(n)==n{return n}}
		n:=tsNumber(tsPrimitiveHint(loop,value,false))
		if math.IsNaN(n)||n < -9007199254740991||n > 9007199254740991 {return 0}
		return math.Floor(n)
	}
	targetStart:=toOffset(tsArg(args,1))
	if targetStart<0 {panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","targetStart must be >= 0")})}
	sourceStart:=toOffset(tsArg(args,2))
	if sourceStart<0 {panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","sourceStart must be >= 0")})}
	if sourceStart>float64(len(source)) {panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","sourceStart is outside the Buffer")})}
	sourceEnd:=float64(0)
	if !tsIsUndefined(tsArg(args,3)) {sourceEnd=toOffset(tsArg(args,3))}
	if !tsIsUndefined(tsArg(args,3))&&sourceEnd<0 {panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","sourceEnd must be >= 0")})}
	if tsIsUndefined(tsArg(args,3))||sourceEnd>float64(len(source)) {sourceEnd=float64(len(source))}
	if targetStart>=float64(len(target))||sourceStart>=sourceEnd {return tsNumberValue(0)}
	start,end:=int(sourceStart),int(sourceEnd);targetOffset:=int(targetStart)
	count:=end-start;if count>len(target)-targetOffset {count=len(target)-targetOffset}
	return tsNumberValue(float64(copy(target[targetOffset:targetOffset+count],source[start:start+count])))
}

func tsIsUint8BufferLike(value tsValue) bool {
	if value.kind!=tsTypedArrayKind { return false }
	array:=(*tsTypedArray)(value.ref)
	return array.nodeBuffer||array.name=="Uint8Array"
}

func tsNodeBufferCompareStatic(args []tsValue) tsValue {
	left,right:=tsArg(args,0),tsArg(args,1)
	view:=func(value tsValue)bool{return value.kind==tsTypedArrayKind||value.kind==tsDataViewKind}
	if !view(left){tsFSInvalid("buf1","a Buffer or Uint8Array")};if !view(right){tsFSInvalid("buf2","a Buffer or Uint8Array")}
	return tsNumberValue(float64(bytes.Compare(tsNodeBufferVisibleBytes(left),tsNodeBufferVisibleBytes(right))))
}

func tsNodeBufferConcat(loop *tsLoop, args []tsValue) tsValue {
	if len(args)==0 {return tsNodeBuffer([]byte{})}
	list:=tsArg(args,0);if list.kind!=tsArrayKind {tsFSInvalid("list","an instance of Array")}
	array:=(*tsArray)(list.ref);if array.length()==0 {return tsNodeBuffer([]byte{})}
	items:=make([]tsValue,array.length());available:=0
	for i:=range items {
		item:=array.at(i);if !tsIsUint8BufferLike(item) {tsFSInvalid("list["+strconv.Itoa(i)+"]","a Buffer or Uint8Array")}
		typed:=(*tsTypedArray)(item.ref);if typed.buffer.detached {panic(tsThrown{tsFSException("TypeError","ERR_INVALID_STATE","ArrayBufferView is detached")})}
		items[i]=item;available+=typed.length*typed.bytes
	}
	length:=available
	if value:=tsArg(args,1);!tsIsUndefined(value) {
		if !tsIsNumeric(value) {tsFSInvalid("length","a number")};n:=tsNumber(value)
		if math.IsNaN(n)||math.IsInf(n,0)||math.Trunc(n)!=n||n<0||n>4294967296 {panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","length is out of range")})}
		length=int(n)
	}
	out:=make([]byte,length);offset:=0
	for _,item:=range items {if offset>=length {break};source:=tsFSBytes(item);offset+=copy(out[offset:],source)}
	return tsNodeBuffer(out)
}

func tsNodeBufferCopyBytesFrom(loop *tsLoop, args []tsValue) tsValue {
	value:=tsArg(args,0);if value.kind!=tsTypedArrayKind {tsFSInvalid("view","a TypedArray")}
	array:=(*tsTypedArray)(value.ref);if array.buffer.detached {panic(tsThrown{tsFSException("TypeError","ERR_INVALID_STATE","Cannot copy from a detached TypedArray")})}
	if array.length==0 {return tsNodeBuffer([]byte{})}
	offset:=0
	if !tsIsUndefined(tsArg(args,1)) {offset=tsNodeBufferValidatedIndex(loop,tsArg(args,1),"offset")}
	if offset>=array.length {return tsNodeBuffer([]byte{})}
	end:=array.length
	if !tsIsUndefined(tsArg(args,2)) {length:=tsNodeBufferValidatedIndex(loop,tsArg(args,2),"length");if length>array.length-offset {length=array.length-offset};end=offset+length}
	data:=tsFSBytes(value);first,last:=offset*array.bytes,end*array.bytes
	return tsNodeBuffer(append([]byte(nil),data[first:last]...))
}

func tsNodeBufferValidatedIndex(loop *tsLoop, value tsValue, name string) int {
	if !tsIsNumeric(value) {tsFSInvalid(name,"a number")}
	n:=tsNumber(value);max:=float64(int(^uint(0)>>1))
	if math.IsNaN(n)||math.IsInf(n,0)||math.Trunc(n)!=n||n<0||n>max {panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE",name+" is out of range")})}
	return int(n)
}

func tsNodeBufferSearchOffset(length int, offset float64, needleLength int, forward bool) int {
	length64 := int64(length)
	offset64 := int64(offset)
	if offset64 < 0 {
		if offset64+length64 >= 0 {
			return int(length64 + offset64)
		}
		if forward || needleLength == 0 {
			return 0
		}
		return -1
	}
	if offset64+needle64 <= length64 {
		return int(offset64)
	}
	if needleLength == 0 {
		return length
	}
	if forward {
		return -1
	}
	return length - 1
}

func tsNodeBufferFind(data, needle []byte, offset, end int, last, utf16Search bool) int {
	if end > len(data) { end = len(data) }
	if end < 0 { end = 0 }
	if utf16Search {
		end &= ^1
		if len(needle)&^1 == 0 { needle = needle[:len(needle)&^1] }
	}
	if len(needle) == 0 {
		if offset < 0 { offset = 0 }; if offset > end { offset = end }; return offset
	}
	if end == 0 || offset < 0 { return -1 }
	step := 1
	if utf16Search { step = 2; offset &= ^1 }
	if last {
		lastStart := offset
		if lastStart+len(needle) > end { lastStart = end-len(needle) }
		lastStart -= lastStart % step
		for i:=lastStart; i>=0; i-=step { if i+len(needle)<=end && bytes.Equal(data[i:i+len(needle)],needle) { return i } }
		return -1
	}
	if offset < 0 { offset = 0 }
	offset += (step-offset%step)%step
	if offset+len(needle)>end { return -1 }
	for i:=offset; i+len(needle)<=end; i+=step { if bytes.Equal(data[i:i+len(needle)],needle) { return i } }
	return -1
}

func tsNodeBufferSearch(loop *tsLoop, receiver tsValue, name string, args []tsValue) tsValue {
	last := name == "lastIndexOf"
	data := tsNodeBufferVisibleBytes(receiver)
	needleValue := tsArg(args,0)
	offsetValue := tsArg(args,1)
	endValue := tsArg(args,2)
	encodingValue := tsArg(args,3)
	if endValue.kind == tsStringKind { encodingValue=endValue; endValue=tsU }
	if offsetValue.kind == tsStringKind { encodingValue=offsetValue; offsetValue=tsU }
	encoding := "utf8"
	encodingValid:=true
	var encodingString *tsString
	isNumber:=tsIsNumeric(needleValue)
	end:=len(data)
	if !tsIsUndefined(endValue) { n:=tsNumber(tsPrimitiveHint(loop,endValue,false));if math.IsNaN(n){n=math.Inf(1)};if n<=0 {end=0} else if n<float64(end) {end=int(math.Trunc(n))} }
	offset:=math.NaN()
	if !tsIsUndefined(offsetValue) { offset=tsNumber(tsPrimitiveHint(loop,offsetValue,false));if offset>2147483647 {offset=2147483647};if offset< -2147483648 {offset= -2147483648} }
	if math.IsNaN(offset) { if last {offset=float64(len(data))} else {offset=0} }
	if !isNumber&&!tsIsUndefined(encodingValue) {
		encodingString=tsStringValueWithLoop(loop,encodingValue)
		name:=strings.ToLower(tsText(tsStringReference(encodingString)))
		switch name {
		case "utf8","utf-8":encoding="utf8"
		case "utf16le","utf-16le","ucs2","ucs-2":encoding="utf16le"
		case "latin1","binary":encoding="latin1"
		case "ascii","base64","base64url":encoding=name
		case "hex","hex\x00":encoding="hex"
		default:encodingValid=false
		}
	}
	var needle []byte
	switch needleValue.kind {
	case tsNumberKind,tsFloat32Kind,tsIntKind,tsInt8Kind,tsInt16Kind,tsInt32Kind,tsInt64Kind,tsUintKind,tsUint8Kind,tsUint16Kind,tsUint32Kind,tsUint64Kind:
		needle=[]byte{byte(tsNumberUint32(tsNumber(needleValue)))}
	case tsStringKind:
		if !encodingValid { _=tsFSEncoding(tsStringReference(encodingString),"") }
		needle=tsFSEncode((*tsString)(needleValue.ref),encoding)
	case tsTypedArrayKind:
		array:=(*tsTypedArray)(needleValue.ref)
		if array.name!="Uint8Array"&&!array.nodeBuffer { tsFSInvalid("value","one of type number or string or an instance of Buffer or Uint8Array") }
		if !encodingValid {encoding="utf8"}
		needle=append([]byte(nil),tsNodeBufferVisibleBytes(needleValue)...)
	default:
		tsFSInvalid("value","one of type number or string or an instance of Buffer or Uint8Array")
	}
	forward:=!last
	searchLength:=len(data)
	if encoding=="utf16le" {searchLength &^= 1}
	searchOffset:=tsNodeBufferSearchOffset(searchLength,offset,len(needle),forward)
	if len(needle)==0 { if searchOffset<0 {searchOffset=0};if searchOffset>end {searchOffset=end};return tsNumberValue(float64(searchOffset)) }
	if searchOffset<0 {return tsNumberValue(-1)}
	result:=tsNodeBufferFind(data,needle,searchOffset,end,last,encoding=="utf16le")
	return tsNumberValue(float64(result))
}

func tsNodeBufferToString(loop *tsLoop, receiver tsValue, args []tsValue) tsValue {
	data := tsNodeBufferVisibleBytes(receiver)
	if len(data) == 0 { return tsStringReference(tsStringUTF8("")) }
	encoding := "utf8"
	if candidate:=tsArg(args,0); !tsIsUndefined(candidate) {
		encoding=tsNodeBufferParseEncoding(loop,candidate,"utf8",false)
	}
	start:=0
	if len(args)>0 {
		n:=tsNumber(tsPrimitiveHint(loop,tsArg(args,1),false))
		if math.IsNaN(n)||n<0 {n=0}
		if n>float64(len(data)) {return tsStringReference(tsStringUTF8(""))}
		start=int(n)
	}
	end:=len(data)
	if candidate:=tsArg(args,2); !tsIsUndefined(candidate) {
		n:=tsNumber(tsPrimitiveHint(loop,candidate,false))
		if math.IsNaN(n)||n<=0 {end=0} else if n<float64(end) {end=int(n)}
	}
	if end<=start {return tsStringReference(tsStringUTF8(""))}
	return tsFSDecode(data[start:end],encoding)
}

func tsNodeBufferIndexArgument(loop *tsLoop, value tsValue, fallback int) int {
	if tsIsUndefined(value) { return fallback }
	n:=tsNumber(tsPrimitiveHint(loop,value,false))
	if math.IsNaN(n)||n==0 { return 0 }
	if n<0 { panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","Index out of range")}) }
	if n>=float64(int(^uint(0)>>1)) { return int(^uint(0)>>1) }
	return int(math.Trunc(n))
}

func tsNodeBufferEncodingSlice(loop *tsLoop, receiver tsValue, encoding string, args []tsValue) tsValue {
	data:=tsNodeBufferVisibleBytes(receiver)
	if len(data)==0 { return tsStringReference(tsStringUTF8("")) }
	start:=tsNodeBufferIndexArgument(loop,tsArg(args,0),0)
	end:=tsNodeBufferIndexArgument(loop,tsArg(args,1),len(data))
	if start>=end { return tsStringReference(tsStringUTF8("")) }
	if end>len(data) { panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","Index out of range")}) }
	return tsFSDecode(data[start:end],encoding)
}

func tsNodeBufferRawWrite(loop *tsLoop, receiver tsValue, encoding string, args []tsValue) tsValue {
	if tsArg(args,0).kind!=tsStringKind { tsFSInvalid("argument","a string") }
	data:=tsNodeBufferVisibleBytes(receiver)
	offset:=tsNodeBufferIndexArgument(loop,tsArg(args,1),0)
	if offset>len(data) { panic(tsThrown{tsFSException("RangeError","ERR_BUFFER_OUT_OF_BOUNDS","offset is outside the bounds of the Buffer")}) }
	length:=tsNodeBufferIndexArgument(loop,tsArg(args,2),len(data)-offset)
	if length>len(data)-offset { length=len(data)-offset }
	if length==0 { return tsNumberValue(0) }
	encoded:=tsFSEncode((*tsString)(tsArg(args,0).ref),encoding)
	return tsNumberValue(float64(copy(data[offset:offset+length],encoded)))
}

func tsNodeBufferStrictWrite(loop *tsLoop, receiver tsValue, encoding string, args []tsValue) tsValue {
	data:=tsNodeBufferVisibleBytes(receiver)
	offsetValue,lengthValue:=tsArg(args,1),tsArg(args,2)
	offset:=float64(0)
	if !tsIsUndefined(offsetValue) { offset=tsNumber(tsPrimitiveHint(loop,offsetValue,false)) }
	length:=float64(0)
	lengthUndefined:=tsIsUndefined(lengthValue)
	if lengthUndefined { length=math.NaN() } else { length=tsNumber(tsPrimitiveHint(loop,lengthValue,false)) }
	if !math.IsNaN(offset)&&(offset<0||offset>float64(len(data))) { panic(tsThrown{tsFSException("RangeError","ERR_BUFFER_OUT_OF_BOUNDS","offset is outside the bounds of the Buffer")}) }
	safeOffset:=0;if !math.IsNaN(offset){safeOffset=int(offset)}
	if lengthUndefined { if math.IsNaN(offset){length=0}else{length=float64(len(data)-safeOffset)} }
	if !math.IsNaN(offset)&&(length<0||length>float64(len(data)-safeOffset)) { panic(tsThrown{tsFSException("RangeError","ERR_BUFFER_OUT_OF_BOUNDS","length is outside the bounds of the Buffer")}) }
	if math.IsNaN(length)||length<0 { length=0 }
	if length>float64(len(data)-safeOffset){length=float64(len(data)-safeOffset)}
	if tsArg(args,0).kind!=tsStringKind { tsFSInvalid("argument","a string") }
	encoded:=tsFSEncode((*tsString)(tsArg(args,0).ref),encoding)
	limit:=int(length);if limit>len(encoded){limit=len(encoded)}
	return tsNumberValue(float64(copy(data[safeOffset:safeOffset+int(length)],encoded[:limit])))
}

func tsNodeBufferWrite(loop *tsLoop, receiver tsValue, args []tsValue) tsValue {
	data:=tsNodeBufferVisibleBytes(receiver)
	stringValue,offsetValue,lengthValue,encodingValue:=tsArg(args,0),tsArg(args,1),tsArg(args,2),tsArg(args,3)
	if tsIsUndefined(offsetValue) {
		if stringValue.kind!=tsStringKind { tsFSInvalid("string","a string") }
		encoded:=tsFSEncode((*tsString)(stringValue.ref),"utf8")
		return tsNumberValue(float64(copy(data,encoded)))
	}
	if tsIsUndefined(lengthValue)&&offsetValue.kind==tsStringKind {
		encodingValue=offsetValue
		offsetValue=tsNumberValue(0)
		stringValue=tsStringReference(tsStringValueWithLoop(loop,stringValue))
	}
	validateOffset:=func(value tsValue,name string)int {
		if !tsIsNumeric(value) { tsFSInvalid(name,"a number") }
		n:=tsNumber(value)
		if math.IsNaN(n)||math.IsInf(n,0)||math.Trunc(n)!=n { panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE",name+" must be an integer")}) }
		if n<0||n>float64(len(data)) { panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE",name+" is out of range")}) }
		return int(n)
	}
	offset:=validateOffset(offsetValue,"offset")
	remaining:=len(data)-offset
	length:=remaining
	if lengthValue.kind==tsStringKind { encodingValue=lengthValue } else if !tsIsUndefined(lengthValue) { length=validateOffset(lengthValue,"length");if length>remaining {length=remaining} }
	if stringValue.kind!=tsStringKind { tsFSInvalid("string","a string") }
	encoding:="utf8"
	if tsTruthy(encodingValue) { encoding=tsNodeBufferParseEncoding(loop,encodingValue,"utf8",false) }
	encoded:=tsFSEncode((*tsString)(stringValue.ref),encoding)
	return tsNumberValue(float64(copy(data[offset:offset+length],encoded)))
}

func tsNodeBufferArrayLike(loop *tsLoop, value tsValue) tsValue {
	data:=make([]byte,0)
	it:=tsArrayLikeIterator(value)
	for it.next() { data=append(data,byte(tsNumberUint32(tsNumber(tsPrimitiveHint(loop,it.value,false))))) }
	return tsNodeBuffer(data)
}

func tsNodeBufferFrom(loop *tsLoop, args []tsValue) tsValue {
	source:=tsArg(args,0)
	encodingOrOffset,lengthValue:=tsArg(args,1),tsArg(args,2)
	if source.kind==tsStringKind {
		encoding:="utf8";if encodingOrOffset.kind==tsStringKind {encoding=tsNodeBufferParseEncoding(loop,encodingOrOffset,"utf8",false)}
		return tsNodeBuffer(tsFSEncode((*tsString)(source.ref),encoding))
	}
	if source.kind==tsArrayBufferKind {
		buffer:=(*tsArrayBuffer)(source.ref)
		if buffer.detached { panic(tsThrown{tsFSException("TypeError","ERR_INVALID_STATE","Buffer is detached")}) }
		byteLength:=len(buffer.data);offsetD:=float64(0)
		if !tsIsUndefined(encodingOrOffset) { offsetD=tsNumber(tsPrimitiveHint(loop,encodingOrOffset,false));if math.IsNaN(offsetD){offsetD=0};if offsetD>float64(byteLength)||math.Trunc(offsetD)<0 {panic(tsThrown{tsFSException("RangeError","ERR_BUFFER_OUT_OF_BOUNDS","offset is outside the bounds of the Buffer")})} }
		offset:=int(offsetD);length:=byteLength-offset
		if !tsIsUndefined(lengthValue) { n:=tsNumber(tsPrimitiveHint(loop,lengthValue,false));if n>0 {if n>float64(byteLength)-offsetD {panic(tsThrown{tsFSException("RangeError","ERR_BUFFER_OUT_OF_BOUNDS","length is outside the bounds of the Buffer")})};length=int(math.Trunc(n))} else {length=0} }
		view:=tsNodeBuffer(buffer.data[offset:offset+length]);array:=(*tsTypedArray)(view.ref);array.buffer=buffer;array.offset=offset;array.length=length;tsArrayBufferRegister(buffer,array);return view
	}
	if source.kind==tsTypedArrayKind {
		array:=(*tsTypedArray)(source.ref)
		if array.buffer.detached { panic(tsThrown{tsFSException("RangeError","ERR_BUFFER_OUT_OF_BOUNDS","Buffer is detached")}) }
		if array.name=="Uint8Array"||array.name=="Uint8ClampedArray" { return tsNodeBuffer(append([]byte(nil),tsFSBytes(source)...)) }
		return tsNodeBufferArrayLike(loop,source)
	}
	if source.kind==tsDataViewKind { return tsNodeBuffer(append([]byte(nil),tsNodeBufferByteView(source)...)) }
	if source.kind==tsArrayKind { return tsNodeBufferArrayLike(loop,source) }
	if source.kind==tsObjectKind {
		if method:=tsFSProperty(source,"valueOf");method.kind==tsFunctionKind {
			converted:=tsCallReceiver(loop,method,source)
			if !tsNullish(converted)&&!tsStrictEqual(converted,source)&&(converted.kind==tsStringKind||converted.kind==tsObjectKind||converted.kind==tsArrayKind||converted.kind==tsTypedArrayKind||converted.kind==tsArrayBufferKind) {
				return tsNodeBufferFrom(loop,[]tsValue{converted,encodingOrOffset,lengthValue})
			}
		}
		valueLength:=tsFSProperty(source,"length")
		needsArrayLike:=!tsIsUndefined(valueLength)
		if !needsArrayLike { needsArrayLike=tsFSProperty(source,"buffer").kind==tsArrayBufferKind }
		if needsArrayLike {
			if !tsIsNumeric(valueLength) { return tsNodeBuffer([]byte{}) }
			n:=tsNumber(valueLength);if n<=0||math.IsNaN(n) { return tsNodeBuffer([]byte{}) }
			return tsNodeBufferArrayLike(loop,source)
		}
		typeValue,dataValue:=tsFSProperty(source,"type"),tsFSProperty(source,"data")
		if tsText(typeValue)=="Buffer"&&dataValue.kind==tsArrayKind { return tsNodeBufferArrayLike(loop,dataValue) }
		toPrimitive:=tsGet(source,tsWellKnownSymbol("toPrimitive"))
		if toPrimitive.kind==tsFunctionKind {
			primitive:=tsCallReceiver(loop,toPrimitive,source,tsStringReference(tsStringUTF8("string")))
			if primitive.kind==tsStringKind { encoding:="utf8";if encodingOrOffset.kind==tsStringKind {encoding=tsNodeBufferParseEncoding(loop,encodingOrOffset,"utf8",false)};return tsNodeBuffer(tsFSEncode((*tsString)(primitive.ref),encoding)) }
		}
	}
	tsFSInvalid("first argument","a string, Buffer, ArrayBuffer, Array, or array-like object")
	return tsU
}

func tsNodeBufferAllocate(loop *tsLoop, args []tsValue, fill bool) tsValue {
	value:=tsArg(args,0)
	if !tsIsNumeric(value) { tsFSInvalid("size","a number") }
	n:=tsNumber(value)
	if math.IsNaN(n)||math.IsInf(n,0)||n<0||n>4294967296 { panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","size is outside the Buffer bounds")}) }
	length:=int(math.Trunc(n))
	buffer:=tsNodeBuffer(make([]byte,length))
	if fill&&length>0&&len(args)>1&&!tsIsUndefined(tsArg(args,1)) { tsNodeBufferFill(loop,buffer,[]tsValue{tsArg(args,1),tsNumberValue(0),tsNumberValue(float64(length)),tsArg(args,2)}) }
	return buffer
}

func tsNodeBufferTranscodeFailure(code string, errno int) {
	error := tsErrorValue(&tsRuntimeError{
		name: "Error",
		message: "Unable to transcode Buffer ["+code+"]",
		fields: map[string]tsValue{"code": tsStringReference(tsStringUTF8(code)), "errno": tsNumberValue(float64(errno))},
	})
	panic(tsThrown{error})
}

func tsNodeBufferEncoding(loop *tsLoop, value tsValue) string {
	if tsNullish(value) { return "utf8" }
	name := strings.ToLower(tsText(tsStringValueWithLoop(loop, value)))
	if name == "" { return "utf8" }
	switch name {
	case "utf8", "utf-8": return "utf8"
	case "ucs2", "ucs-2", "utf16le", "utf-16le": return "ucs2"
	case "latin1", "binary": return "latin1"
	case "ascii": return "ascii"
	default: return "unsupported"
	}
}

func tsNodeBufferWellFormedUTF16(units []uint16) []uint16 {
	out := make([]uint16, 0, len(units))
	for i := 0; i < len(units); i++ {
		unit := units[i]
		if unit >= 0xd800 && unit <= 0xdbff {
			if i+1 < len(units) && units[i+1] >= 0xdc00 && units[i+1] <= 0xdfff {
				out = append(out, unit, units[i+1]); i++; continue
			}
			out = append(out, 0xfffd); continue
		}
		if unit >= 0xdc00 && unit <= 0xdfff { out = append(out, 0xfffd); continue }
		out = append(out, unit)
	}
	return out
}

func tsNodeBufferStrictUTF16(units []uint16) ([]rune, bool) {
	runes := make([]rune, 0, len(units))
	for i := 0; i < len(units); i++ {
		unit := units[i]
		if unit >= 0xd800 && unit <= 0xdbff {
			if i+1 >= len(units) || units[i+1] < 0xdc00 || units[i+1] > 0xdfff { return nil, false }
			runes = append(runes, utf16.DecodeRune(unit, units[i+1])); i++; continue
		}
		if unit >= 0xdc00 && unit <= 0xdfff { return nil, false }
		runes = append(runes, rune(unit))
	}
	return runes, true
}

func tsNodeBufferDecodeUnits(data []byte, encoding string, replaceOdd bool) []uint16 {
	switch encoding {
	case "latin1":
		units := make([]uint16, len(data)); for i, item := range data { units[i] = uint16(item) }; return units
	case "ascii":
		units := make([]uint16, len(data)); for i, item := range data { if item > 0x7f { units[i] = 0xfffd } else { units[i] = uint16(item) } }; return units
	case "utf8":
		units := make([]uint16, 0, len(data)); text := string(data)
		for len(text) > 0 { r, size := utf8.DecodeRuneInString(text); if r == utf8.RuneError && size == 1 { units = append(units, 0xfffd); text = text[1:]; continue }; units = append(units, utf16.Encode([]rune{r})...); text = text[size:] }
		return units
	case "ucs2":
		count := len(data)/2; if replaceOdd && len(data)%2 != 0 { count++ }
		units := make([]uint16, count); for i := 0; i < len(data)/2; i++ { units[i] = binary.LittleEndian.Uint16(data[i*2:]) }
		if count > len(data)/2 { units[count-1] = 0xfffd }
		return tsNodeBufferWellFormedUTF16(units)
	}
	return nil
}

func tsNodeBufferEncodeUnits(units []uint16, encoding string) []byte {
	switch encoding {
	case "ucs2":
		out := make([]byte, len(units)*2); for i, unit := range units { binary.LittleEndian.PutUint16(out[i*2:], unit) }; return out
	case "latin1", "ascii":
		limit := uint16(0xff); if encoding == "ascii" { limit = 0x7f }
		out := make([]byte, 0, len(units))
		for i, unit := range units { if unit >= 0xdc00 && unit <= 0xdfff && i > 0 && units[i-1] >= 0xd800 && units[i-1] <= 0xdbff { continue }; if unit > limit { out = append(out, '?') } else { out = append(out, byte(unit)) } }
		return out
	default:
		runes, _ := tsNodeBufferStrictUTF16(units); return []byte(string(runes))
	}
}

func tsNodeBufferTranscode(loop *tsLoop, args []tsValue) tsValue {
	source := tsArg(args, 0)
	if source.kind != tsTypedArrayKind { tsNodeBufferModuleTypeError("The \"source\" argument must be of type Buffer or Uint8Array") }
	array := (*tsTypedArray)(source.ref)
	if !array.nodeBuffer && array.name != "Uint8Array" { tsNodeBufferModuleTypeError("The \"source\" argument must be of type Buffer or Uint8Array") }
	if array.buffer.detached || array.length == 0 { return tsNodeBuffer([]byte{}) }
	from, to := tsNodeBufferEncoding(loop, tsArg(args, 1)), tsNodeBufferEncoding(loop, tsArg(args, 2))
	if from == "unsupported" || to == "unsupported" { tsNodeBufferTranscodeFailure("U_ILLEGAL_ARGUMENT_ERROR", 1) }
	if array.buffer.detached || array.length == 0 { return tsNodeBuffer([]byte{}) }
	data := array.buffer.data[array.offset:array.offset+array.length*array.bytes]
	if (from == "ascii" || from == "latin1") && to == "ucs2" {
		units := make([]uint16, len(data)); for i, item := range data { units[i] = uint16(item) }; return tsNodeBuffer(tsNodeBufferEncodeUnits(units, to))
	}
	if from == "utf8" && to == "ucs2" {
		if !utf8.Valid(data) { tsNodeBufferTranscodeFailure("U_INVALID_CHAR_FOUND", 10) }
		units := utf16.Encode([]rune(string(data))); return tsNodeBuffer(tsNodeBufferEncodeUnits(units, to))
	}
	if from == "ucs2" && to == "utf8" {
		count := len(data)/2; units := make([]uint16, count); for i := range units { units[i] = binary.LittleEndian.Uint16(data[i*2:]) }
		if count == 0 { tsNodeBufferTranscodeFailure("U_INVALID_CHAR_FOUND", 10) }
		runes, valid := tsNodeBufferStrictUTF16(units); if !valid { tsNodeBufferTranscodeFailure("U_INVALID_CHAR_FOUND", 10) }
		return tsNodeBuffer([]byte(string(runes)))
	}
	units := tsNodeBufferDecodeUnits(data, from, to == "ucs2")
	return tsNodeBuffer(tsNodeBufferEncodeUnits(units, to))
}

func tsNodeBufferAtob(loop *tsLoop, args []tsValue) tsValue {
	text := tsText(tsStringValueWithLoop(loop, tsArg(args, 0)))
	text = strings.Map(func(r rune) rune { if r == ' ' || r == '\t' || r == '\n' || r == '\r' || r == '\f' { return -1 }; return r }, text)
	for i, r := range text { if !(r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '+' || r == '/' || r == '=') { _ = i; panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"InvalidCharacterError",message:"The string to be decoded is not correctly encoded."})}) } }
	if at := strings.IndexByte(text, '='); at >= 0 { if len(text)-at > 2 || strings.Trim(text[at:], "=") != "" { panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"InvalidCharacterError",message:"The string to be decoded is not correctly encoded."})}) }; text = text[:at] }
	if len(text)%4 == 1 { panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"InvalidCharacterError",message:"The string to be decoded is not correctly encoded."})}) }
	data, err := base64.RawStdEncoding.DecodeString(text); if err != nil { panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"InvalidCharacterError",message:"The string to be decoded is not correctly encoded."})}) }
	units := make([]uint16, len(data)); for i, item := range data { units[i] = uint16(item) }; return tsStringReference(tsStringUnits(units))
}

func tsNodeBufferBtoa(loop *tsLoop, args []tsValue) tsValue {
	text := tsStringValueWithLoop(loop, tsArg(args, 0)); data := make([]byte, len(text.units))
	for i, unit := range text.units { if unit > 0xff { panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"InvalidCharacterError",message:"The string to be encoded contains characters outside of the Latin1 range."})}) }; data[i] = byte(unit) }
	return tsStringReference(tsStringUTF8(base64.StdEncoding.EncodeToString(data)))
}

func tsNodeBufferFileConstruct(loop *tsLoop, args []tsValue) tsValue {
	blob := tsBlobConstruct(loop, []tsValue{tsArg(args, 0), tsArg(args, 2)})
	if blob.kind != tsObjectKind { return blob }
	object := (*tsObject)(blob.ref); object.nativeClass = tsNativeClass("File")
	object.set("name", tsStringReference(tsStringValueWithLoop(loop, tsArg(args, 1))))
	lastModified := float64(time.Now().UnixMilli()); if options := tsArg(args, 2); !tsIsUndefined(options) { candidate := tsGet(options, tsStringReference(tsStringUTF8("lastModified"))); if !tsIsUndefined(candidate) { lastModified = tsNumber(tsToPrimitive(candidate)) } }
	object.set("lastModified", tsNumberValue(lastModified))
	object.descriptors = map[string]*tsDescriptor{"name": {writable:false, configurable:true}, "lastModified": {writable:false, configurable:true}}
	return blob
}

func tsNodeBufferSlowBuffer() tsValue {
	function := tsFunc(func(loop *tsLoop, args ...tsValue) tsValue { return tsNodeBufferAllocate(loop,args,false) })
	function.name = "SlowBuffer"; function.length = 1; function.constructible = true; function.construct = function.call
	return tsFunctionValue(function)
}

func tsNodeBufferPopulateNamespace(out *tsObject, loop *tsLoop) {
	out.set("SlowBuffer", tsNodeBufferSlowBuffer())
	out.set("Blob", tsClassValue(tsNativeClass("Blob")))
	out.set("File", tsClassValue(tsNativeClass("File")))
	out.set("INSPECT_MAX_BYTES", tsNumberValue(50))
	out.set("kMaxLength", tsNumberValue(4294967296))
	out.set("kStringMaxLength", tsNumberValue(536870888))
	constants := tsNewObject(); constants.set("MAX_LENGTH", tsNumberValue(4294967296)); constants.set("MAX_STRING_LENGTH", tsNumberValue(536870888)); out.set("constants", tsObjectValue(constants))
	out.set("atob", tsFunctionValue(tsFunc(func(loop *tsLoop, args ...tsValue) tsValue { return tsNodeBufferAtob(loop, args) })))
	out.set("btoa", tsFunctionValue(tsFunc(func(loop *tsLoop, args ...tsValue) tsValue { return tsNodeBufferBtoa(loop, args) })))
	out.set("transcode", tsFunctionValue(tsFunc(func(loop *tsLoop, args ...tsValue) tsValue { return tsNodeBufferTranscode(loop, args) })))
	out.set("resolveObjectURL", tsFunctionValue(tsFunc(func(_ *tsLoop, _ ...tsValue) tsValue { return tsU })))
	out.set("isAscii", tsFunctionValue(tsFunc(func(_ *tsLoop, args ...tsValue) tsValue { return tsNodeBufferIsAscii(tsArg(args, 0)) })))
	out.set("isUtf8", tsFunctionValue(tsFunc(func(_ *tsLoop, args ...tsValue) tsValue { return tsNodeBufferIsUtf8(tsArg(args, 0)) })))
}
`
