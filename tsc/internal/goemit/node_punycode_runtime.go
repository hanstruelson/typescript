//go:build !no_node_modules

package goemit

// nodePunycodeRuntime is Bun's standalone RFC 3492 punycode module translated
// into runtime helpers. It has no dependency on other Node builtins.
const nodePunycodeRuntime = `
var tsNodePunycodeSingleton struct { sync.Once; value *tsObject }

func tsNodePunycodeError(kind string) {
	message := "Invalid input"
	switch kind {
	case "overflow": message = "Overflow: input needs wider integers to process"
	case "not-basic": message = "Illegal input >= 0x80 (not a basic code point)"
	}
	panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"RangeError", message:message})})
}

func tsNodePunycodeBasicToDigit(value uint32) int64 {
	if value >= 0x30 && value < 0x3a { return int64(26 + value - 0x30) }
	if value >= 0x41 && value < 0x5b { return int64(value - 0x41) }
	if value >= 0x61 && value < 0x7b { return int64(value - 0x61) }
	return 36
}

func tsNodePunycodeDigitToBasic(digit int64) uint16 {
	return uint16(digit + 22 + 75*tsNodePunycodeBoolInt(digit < 26))
}

func tsNodePunycodeBoolInt(value bool) int64 { if value { return 1 }; return 0 }

func tsNodePunycodeAdapt(delta, points int64, first bool) int64 {
	if first { delta /= 700 } else { delta /= 2 }
	delta += delta / points
	k := int64(0)
	for delta > 455 { delta /= 35; k += 36 }
	return k + (36*delta)/(delta+38)
}

func tsNodePunycodeCodepointsToString(points []uint32) *tsString {
	units := make([]uint16, 0, len(points))
	for _, point := range points {
		if point > 0x10ffff {
			panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"RangeError",message:"Invalid code point " + strconv.FormatUint(uint64(point), 10)})})
		}
		if point <= 0xffff { units = append(units, uint16(point)); continue }
		point -= 0x10000
		units = append(units, uint16(0xd800+(point>>10)), uint16(0xdc00+(point&0x3ff)))
	}
	return tsStringUnits(units)
}

func tsNodePunycodeUCS2Decode(value *tsString) []uint32 {
	output := make([]uint32, 0, len(value.units))
	for index := 0; index < len(value.units); index++ {
		unit := value.units[index]
		if unit >= 0xd800 && unit <= 0xdbff && index+1 < len(value.units) {
			extra := value.units[index+1]
			if extra&0xfc00 == 0xdc00 {
				output = append(output, (uint32(unit&0x3ff)<<10)+uint32(extra&0x3ff)+0x10000)
				index++
				continue
			}
		}
		output = append(output, uint32(unit))
	}
	return output
}

func tsNodePunycodeEncode(input *tsString) *tsString {
	const maxInt int64 = 2147483647
	const base int64 = 36
	points := tsNodePunycodeUCS2Decode(input)
	output := make([]uint16, 0, len(points)+8)
	for _, point := range points { if point < 0x80 { output = append(output, uint16(point)) } }
	basicLength := len(output)
	handled := basicLength
	if basicLength > 0 { output = append(output, '-') }
	n, delta, bias := int64(128), int64(0), int64(72)
	for handled < len(points) {
		m := maxInt
		for _, point := range points { value := int64(point); if value >= n && value < m { m = value } }
		plusOne := int64(handled + 1)
		if m-n > (maxInt-delta)/plusOne { tsNodePunycodeError("overflow") }
		delta += (m - n) * plusOne
		n = m
		for _, point := range points {
			value := int64(point)
			if value < n { delta++; if delta > maxInt { tsNodePunycodeError("overflow") } }
			if value != n { continue }
			q := delta
			for k := int64(36); ; k += 36 {
				threshold := int64(1)
				if k > bias+26 { threshold = 26 } else if k > bias { threshold = k-bias }
				if q < threshold { break }
				baseMinusThreshold := base-threshold
				output = append(output, tsNodePunycodeDigitToBasic(threshold+(q-threshold)%baseMinusThreshold))
				q = (q-threshold)/baseMinusThreshold
			}
			output = append(output, tsNodePunycodeDigitToBasic(q))
			bias = tsNodePunycodeAdapt(delta, plusOne, handled == basicLength)
			delta = 0
			handled++
		}
		delta++
		n++
	}
	return tsStringUnits(output)
}

func tsNodePunycodeDecode(input *tsString) *tsString {
	const maxInt int64 = 2147483647
	const base int64 = 36
	units := input.units
	basic := -1
	for index, unit := range units { if unit == '-' { basic = index } }
	output := make([]uint32, 0, len(units))
	if basic < 0 { basic = 0 } else {
		for _, unit := range units[:basic] {
			if unit >= 0x80 { tsNodePunycodeError("not-basic") }
			output = append(output, uint32(unit))
		}
	}
	index := 0
	if basic > 0 { index = basic+1 }
	i, n, bias := int64(0), int64(128), int64(72)
	for index < len(units) {
		oldI := i
		w := int64(1)
		for k := int64(36); ; k += 36 {
			if index >= len(units) { tsNodePunycodeError("invalid-input") }
			digit := tsNodePunycodeBasicToDigit(uint32(units[index]))
			index++
			if digit >= base { tsNodePunycodeError("invalid-input") }
			if digit > (maxInt-i)/w { tsNodePunycodeError("overflow") }
			i += digit*w
			threshold := int64(1)
			if k > bias+26 { threshold = 26 } else if k > bias { threshold = k-bias }
			if digit < threshold { break }
			baseMinusThreshold := base-threshold
			if w > maxInt/baseMinusThreshold { tsNodePunycodeError("overflow") }
			w *= baseMinusThreshold
		}
		outCount := int64(len(output)+1)
		bias = tsNodePunycodeAdapt(i-oldI, outCount, oldI == 0)
		if i/outCount > maxInt-n { tsNodePunycodeError("overflow") }
		n += i/outCount
		i %= outCount
		insert := int(i)
		output = append(output, 0)
		copy(output[insert+1:], output[insert:])
		output[insert] = uint32(n)
		i++
	}
	return tsNodePunycodeCodepointsToString(output)
}

func tsNodePunycodeMapDomain(input *tsString, decode bool) *tsString {
	units := input.units
	local := []uint16(nil)
	at := -1
	for index, unit := range units { if unit == '@' { at = index; break } }
	if at >= 0 {
		local = append(append([]uint16{}, units[:at]...), '@')
		units = units[at+1:]
		for next, unit := range units { if unit == '@' { units = units[:next]; break } }
	}
	normalized := make([]uint16, len(units))
	for index, unit := range units {
		switch unit { case 0x3002, 0xff0e, 0xff61: normalized[index] = '.'; default: normalized[index] = unit }
	}
	labels := make([]*tsString, 0, 1)
	start := 0
	for index := 0; index <= len(normalized); index++ {
		if index != len(normalized) && normalized[index] != '.' { continue }
		label := tsStringUnits(normalized[start:index])
		if decode {
			if len(label.units) >= 4 && label.units[0] == 'x' && label.units[1] == 'n' && label.units[2] == '-' && label.units[3] == '-' {
				suffix := strings.ToLower(tsStringUnits(label.units[4:]).String())
				label = tsNodePunycodeDecode(tsStringUTF8(suffix))
			}
		} else {
			nonASCII := false
			for _, unit := range label.units { if unit > 0x7f { nonASCII = true; break } }
			if nonASCII { label = tsStringConcat(tsStringUTF8("xn--"), tsNodePunycodeEncode(label)) }
		}
		labels = append(labels, label)
		start = index+1
	}
	joined := make([]*tsString, 0, len(labels)*2)
	for index, label := range labels { if index > 0 { joined = append(joined, tsStringUTF8(".")) }; joined = append(joined, label) }
	return tsStringConcat(tsStringUnits(local), tsStringConcat(joined...))
}

func tsNodePunycodeFromCodepoints(points []uint32) *tsString { return tsNodePunycodeCodepointsToString(points) }

func tsNodePunycodeModule() tsValue {
	tsNodePunycodeSingleton.Do(func() {
		out := tsNewObject()
		ucs2 := tsNewObject()
		ucs2.set("decode", tsNamedNativeMethod("decode", 1, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
			return tsArrayValue(&tsArray{values: func() []tsValue { points:=tsNodePunycodeUCS2Decode(tsStringValueWithLoop(loop,tsArg(args,0))); values:=make([]tsValue,len(points)); for i,point:=range points { values[i]=tsNumberValue(float64(point)) }; return values }()})
		}))
		ucs2.set("encode", tsNamedNativeMethod("encode", 1, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
			array:=tsArg(args,0); if array.kind!=tsArrayKind { tsNodePunycodeError("invalid-input") }
			values:=(*tsArray)(array.ref); points:=make([]uint32,0,values.length())
			for i:=0;i<values.length();i++ { number:=tsNumber(values.at(i)); if math.IsNaN(number)||math.IsInf(number,0)||math.Trunc(number)!=number||number<0||number>0x10ffff { tsNodePunycodeError("invalid-input") }; points=append(points,uint32(number)) }
			return tsStringReference(tsNodePunycodeFromCodepoints(points))
		}))
		out.set("version", tsStringReference(tsStringUTF8("2.1.0")))
		out.set("ucs2", tsObjectValue(ucs2))
		out.set("decode", tsNamedNativeMethod("decode",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsStringReference(tsNodePunycodeDecode(tsStringValueWithLoop(loop,tsArg(args,0))))}))
		out.set("encode", tsNamedNativeMethod("encode",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsStringReference(tsNodePunycodeEncode(tsStringValueWithLoop(loop,tsArg(args,0))))}))
		out.set("toASCII", tsNamedNativeMethod("toASCII",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsStringReference(tsNodePunycodeMapDomain(tsStringValueWithLoop(loop,tsArg(args,0)),false))}))
		out.set("toUnicode", tsNamedNativeMethod("toUnicode",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsStringReference(tsNodePunycodeMapDomain(tsStringValueWithLoop(loop,tsArg(args,0)),true))}))
		tsNodePunycodeSingleton.value=out
	})
	return tsObjectValue(tsNodePunycodeSingleton.value)
}
`
