//go:build !no_node_modules

package goemit

// nodeQuerystringRuntime implements Bun's node:querystring entry point in the
// generated runtime. Each Node builtin contributes one fragment in this package.
const nodeQuerystringRuntime = `
var tsNodeQuerystringSingleton struct { sync.Once; value *tsObject }

func tsNodeQuerystringError(name, message string) {
	panic(tsThrown{tsErrorValue(&tsRuntimeError{name:name, message:message})})
}

func tsNodeQuerystringString(loop *tsLoop, value tsValue) *tsString {
	if value.kind == tsSymbolKind {
		return tsStringUTF8(tsSymbolText((*tsSymbol)(value.ref)))
	}
	return tsStringValueWithLoop(loop, value)
}

func tsNodeQuerystringImplicitString(loop *tsLoop, value tsValue) *tsString {
	if value.kind == tsSymbolKind { tsNodeQuerystringError("TypeError", "Cannot convert a Symbol value to a string") }
	return tsStringValueWithLoop(loop, value)
}

func tsNodeQuerystringPrimitive(value tsValue) *tsString {
	switch value.kind {
	case tsStringKind:
		return (*tsString)(value.ref)
	case tsNumberKind, tsFloat32Kind, tsIntKind, tsInt8Kind, tsInt16Kind, tsInt32Kind, tsInt64Kind,
		tsUintKind, tsUint8Kind, tsUint16Kind, tsUint32Kind, tsUint64Kind:
		if !math.IsInf(tsNumber(value), 0) && !math.IsNaN(tsNumber(value)) {
			return tsStringUTF8(tsText(value))
		}
	case tsBigIntKind:
		return tsStringUTF8(tsText(value))
	case tsBooleanKind:
		if value.number != 0 { return tsStringUTF8("true") }
		return tsStringUTF8("false")
	}
	return tsStringUnits(nil)
}

func tsNodeQuerystringEscapeText(value *tsString) *tsString {
	var output strings.Builder
	writeByte := func(b byte) {
		const digits = "0123456789ABCDEF"
		output.WriteByte('%')
		output.WriteByte(digits[b>>4])
		output.WriteByte(digits[b&15])
	}
	for i := 0; i < len(value.units); i++ {
		unit := value.units[i]
		if unit < 0x80 {
			b := byte(unit)
			if b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z' || b >= '0' && b <= '9' || strings.ContainsRune("-_.!~*'()", rune(b)) {
				output.WriteByte(b)
			} else {
				writeByte(b)
			}
			continue
		}
		var codepoint uint32
		if unit >= 0xd800 && unit <= 0xdbff {
			if i+1 >= len(value.units) || value.units[i+1] < 0xdc00 || value.units[i+1] > 0xdfff {
				tsNodeQuerystringError("URIError", "URI malformed")
			}
			i++
			codepoint = 0x10000 + (uint32(unit-0xd800) << 10) + uint32(value.units[i]-0xdc00)
		} else if unit >= 0xdc00 && unit <= 0xdfff {
			tsNodeQuerystringError("URIError", "URI malformed")
		} else {
			codepoint = uint32(unit)
		}
		switch {
		case codepoint < 0x800:
			writeByte(byte(0xc0 | codepoint>>6))
			writeByte(byte(0x80 | codepoint&0x3f))
		case codepoint < 0x10000:
			writeByte(byte(0xe0 | codepoint>>12))
			writeByte(byte(0x80 | codepoint>>6&0x3f))
			writeByte(byte(0x80 | codepoint&0x3f))
		default:
			writeByte(byte(0xf0 | codepoint>>18))
			writeByte(byte(0x80 | codepoint>>12&0x3f))
			writeByte(byte(0x80 | codepoint>>6&0x3f))
			writeByte(byte(0x80 | codepoint&0x3f))
		}
	}
	return tsStringUTF8(output.String())
}

func tsNodeQuerystringUnescapeBufferText(value *tsString, decodeSpaces bool) []byte {
	units := value.units
	output := make([]byte, len(units))
	outIndex := 0
	for index := 0; index < len(units); {
		current := units[index]
		if current == '+' && decodeSpaces {
			output[outIndex] = ' '
			outIndex++
			index++
			continue
		}
		if current == '%' && index+2 < len(units) {
			hexValue := func(unit uint16) int {
				switch {
				case unit >= '0' && unit <= '9': return int(unit - '0')
				case unit >= 'A' && unit <= 'F': return int(unit-'A') + 10
				case unit >= 'a' && unit <= 'f': return int(unit-'a') + 10
				default: return -1
				}
			}
			index++
			high := hexValue(units[index])
			if high < 0 {
				output[outIndex] = '%'
				outIndex++
				continue
			}
			index++
			low := hexValue(units[index])
			if low < 0 {
				output[outIndex] = '%'
				outIndex++
				index--
			} else {
				current = uint16(high*16 + low)
			}
		}
		output[outIndex] = byte(current)
		outIndex++
		index++
	}
	return output[:outIndex]
}

func tsNodeQuerystringUnescapeText(value *tsString, decodeSpaces bool) *tsString {
	decoded, err := url.PathUnescape(value.String())
	if err == nil && utf8.ValidString(decoded) {
		return tsStringUTF8(decoded)
	}
	decodedValue := tsFSDecode(tsNodeQuerystringUnescapeBufferText(value, decodeSpaces), "utf8")
	return (*tsString)(decodedValue.ref)
}

func tsNodeQuerystringHasValidEscape(value []uint16) bool {
	for index := 0; index+2 < len(value); index++ {
		if value[index] != '%' { continue }
		isHex := func(unit uint16) bool {
			return unit >= '0' && unit <= '9' || unit >= 'A' && unit <= 'F' || unit >= 'a' && unit <= 'f'
		}
		if isHex(value[index+1]) && isHex(value[index+2]) { return true }
	}
	return false
}

// refString is kept local to this module adapter to preserve the runtime's UTF-16 string value.
func tsNodeQuerystringDecodeURIComponent(value *tsString) (result *tsString, ok bool) {
	decoded, err := url.PathUnescape(value.String())
	if err != nil || !utf8.ValidString(decoded) { return nil, false }
	return tsStringUTF8(decoded), true
}

func tsNodeQuerystringCallDecoder(loop *tsLoop, decoder tsValue, input *tsString, decodeSpaces bool) (result tsValue) {
	defer func() {
		if recover() != nil {
			result = tsStringReference(tsNodeQuerystringUnescapeText(input, decodeSpaces))
		}
	}()
	return tsCall(loop, decoder, tsStringReference(input))
}

func tsNodeQuerystringSplitUnits(value, separator []uint16, start int) int {
	if len(separator) == 0 { return -1 }
	for index := start; index+len(separator) <= len(value); index++ {
		match := true
		for offset, unit := range separator {
			if value[index+offset] != unit { match = false; break }
		}
		if match { return index }
	}
	return -1
}

func tsNodeQuerystringReplacePlus(value []uint16, custom bool) *tsString {
	units := make([]uint16, 0, len(value))
	for _, unit := range value {
		if unit != '+' { units = append(units, unit); continue }
		if custom { units = append(units, '%', '2', '0') } else { units = append(units, ' ') }
	}
	return tsStringUnits(units)
}

func tsNodeQuerystringAdd(loop *tsLoop, out *tsObject, key string, value tsValue) {
	if current, ok := out.values[key]; ok {
		if current.kind == tsArrayKind {
			array := (*tsArray)(current.ref)
			array.appendItems(loop, value)
		} else {
			out.set(key, tsArrayValue(&tsArray{values: []tsValue{current, value}}))
		}
		return
	}
	out.set(key, value)
}

func tsNodeQuerystringParse(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
	out := tsNewObject()
	out.prototype = tsNull
	query := tsArg(args, 0)
	if query.kind != tsStringKind || len((*tsString)(query.ref).units) == 0 { return tsObjectValue(out) }
	sep, eq := tsStringUTF8("&"), tsStringUTF8("=")
	if value := tsArg(args, 1); !tsIsUndefined(value) && tsTruthy(value) { sep = tsNodeQuerystringString(loop, value) }
	if value := tsArg(args, 2); !tsIsUndefined(value) && tsTruthy(value) { eq = tsNodeQuerystringString(loop, value) }
	options := tsArg(args, 3)
	maxPairs := 1000.0
	if value := tsFSProperty(options, "maxKeys"); tsIsNumeric(value) {
		maxPairs = tsNumber(value)
		if maxPairs <= 0 { maxPairs = -1 }
	}
	decoder := tsFunctionValue(tsFunc(func(loop *tsLoop, args ...tsValue) tsValue { return tsStringReference(tsNodeQuerystringUnescapeText(tsNodeQuerystringString(loop, tsArg(args, 0)), false)) }))
	customDecode := tsFSProperty(options, "decodeURIComponent")
	if customDecode.kind == tsFunctionKind { decoder = customDecode }
	custom := customDecode.kind == tsFunctionKind
	units := (*tsString)(query.ref).units
	separator, equality := sep.units, eq.units
	if len(separator) == 0 { separator = []uint16{'&'} }
	if len(equality) == 0 { equality = []uint16{'='} }
	pairs := 0.0
	for start := 0; start < len(units); {
		end := tsNodeQuerystringSplitUnits(units, separator, start)
		if end < 0 { end = len(units) }
		field := units[start:end]
		if len(field) > 0 {
			if maxPairs > 0 && pairs >= maxPairs { break }
			eqAt := tsNodeQuerystringSplitUnits(field, equality, 0)
			keyUnits, valueUnits := field, []uint16{}
			if eqAt >= 0 { keyUnits, valueUnits = field[:eqAt], field[eqAt+len(equality):] }
			keyText := tsNodeQuerystringReplacePlus(keyUnits, custom)
			valueText := tsNodeQuerystringReplacePlus(valueUnits, custom)
			keyValue := tsStringReference(keyText)
			valueValue := tsStringReference(valueText)
			if len(keyText.units) > 0 && (custom || tsNodeQuerystringHasValidEscape(keyText.units)) {
				keyValue = tsNodeQuerystringCallDecoder(loop, decoder, keyText, true)
			}
			if len(valueText.units) > 0 && (custom || tsNodeQuerystringHasValidEscape(valueText.units)) {
				valueValue = tsNodeQuerystringCallDecoder(loop, decoder, valueText, true)
			}
			key := tsPropertyKey(keyValue)
			tsNodeQuerystringAdd(loop, out, key, valueValue)
			pairs++
		} else if end < len(units) {
			pairs++
		}
		if end == len(units) { break }
		start = end + len(separator)
	}
	return tsObjectValue(out)
}

func tsNodeQuerystringStringify(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
	object := tsArg(args, 0)
	if object.kind != tsObjectKind && object.kind != tsArrayKind { return tsStringReference(tsStringUTF8("")) }
	sep, eq := tsStringUTF8("&"), tsStringUTF8("=")
	if value := tsArg(args, 1); !tsIsUndefined(value) && tsTruthy(value) { sep = tsNodeQuerystringImplicitString(loop, value) }
	if value := tsArg(args, 2); !tsIsUndefined(value) && tsTruthy(value) { eq = tsNodeQuerystringImplicitString(loop, value) }
	encoder := tsFSProperty(tsArg(args, 3), "encodeURIComponent")
	custom := encoder.kind == tsFunctionKind
	keys := tsObjectKeys(object).values
	fields := make([]*tsString, 0, len(keys))
	for _, keyValue := range keys {
		key := tsStringValueWithLoop(loop, keyValue)
		encodedKey := tsNodeQuerystringEscapeText(key)
		if custom {
			encodedKey = tsStringValueWithLoop(loop, tsCall(loop, encoder, tsStringReference(key)))
		}
		value := tsGet(object, keyValue)
		values := []tsValue{value}
		if value.kind == tsArrayKind {
			values = (*tsArray)(value.ref).values
			if len(values) == 0 { continue }
		}
		for _, item := range values {
			primitive := tsNodeQuerystringPrimitive(item)
			encodedValue := primitive
			if custom {
				encodedValue = tsStringValueWithLoop(loop, tsCall(loop, encoder, tsStringReference(primitive)))
			} else if item.kind == tsStringKind && len(primitive.units) > 0 || tsIsNumeric(item) && !math.IsInf(tsNumber(item), 0) && !math.IsNaN(tsNumber(item)) && math.Abs(tsNumber(item)) >= 1e21 {
				encodedValue = tsNodeQuerystringEscapeText(primitive)
			} else if item.kind == tsStringKind || tsIsNumeric(item) || item.kind == tsBigIntKind || item.kind == tsBooleanKind {
				// The default encoder leaves ordinary primitive values unescaped.
			} else {
				encodedValue = tsStringUnits(nil)
			}
			fields = append(fields, tsStringConcat(encodedKey, eq, encodedValue))
		}
	}
	return tsStringReference(tsStringConcat(tsNodeQuerystringJoin(fields, sep)...))
}

func tsNodeQuerystringJoin(fields []*tsString, separator *tsString) []*tsString {
	if len(fields) == 0 { return []*tsString{tsStringUnits(nil)} }
	parts := make([]*tsString, 0, len(fields)*2-1)
	for i, field := range fields { if i > 0 { parts = append(parts, separator) }; parts = append(parts, field) }
	return parts
}

func tsNodeQuerystringModule() tsValue {
	tsNodeQuerystringSingleton.Do(func() {
		out := tsNewObject()
		parse := tsNamedNativeMethod("parse", 1, tsNodeQuerystringParse)
		stringify := tsNamedNativeMethod("stringify", 3, tsNodeQuerystringStringify)
		escape := tsNamedNativeMethod("escape", 1, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
			return tsStringReference(tsNodeQuerystringEscapeText(tsNodeQuerystringImplicitString(loop, tsArg(args, 0))))
		})
		unescape := tsNamedNativeMethod("unescape", 1, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
			return tsStringReference(tsNodeQuerystringUnescapeText(tsNodeQuerystringImplicitString(loop, tsArg(args, 0)), tsTruthy(tsArg(args, 1))))
		})
		unescapeBuffer := tsNamedNativeMethod("unescapeBuffer", 2, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
			return tsNodeBuffer(tsNodeQuerystringUnescapeBufferText(tsNodeQuerystringString(loop, tsArg(args, 0)), tsTruthy(tsArg(args, 1))))
		})
		out.set("parse", parse)
		out.set("decode", parse)
		out.set("stringify", stringify)
		out.set("encode", stringify)
		out.set("escape", escape)
		out.set("unescape", unescape)
		out.set("unescapeBuffer", unescapeBuffer)
		tsNodeQuerystringSingleton.value = out
	})
	return tsObjectValue(tsNodeQuerystringSingleton.value)
}
`
