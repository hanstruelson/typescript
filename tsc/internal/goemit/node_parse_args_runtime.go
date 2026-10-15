package goemit

// nodeParseArgsRuntime carries the option metadata and token classifier used
// by Bun's Rust parse_args and parse_args_utils sources.
const nodeParseArgsRuntime = `
type tsNodeParseOption struct {
	longName string
	shortName string
	valueType string
	multiple bool
	defaultValue tsValue
	hasDefault bool
}

type tsNodeParseTokenKind uint8
const (
	tsNodeParsePositional tsNodeParseTokenKind = iota
	tsNodeParseOptionKind
	tsNodeParseTerminator
)

type tsNodeParseTokenSubtype uint8
const (
	tsNodeParseSubtypeTerminator tsNodeParseTokenSubtype = iota
	tsNodeParseSubtypeLoneShort
	tsNodeParseSubtypeShortGroup
	tsNodeParseSubtypeShortValue
	tsNodeParseSubtypeLoneLong
	tsNodeParseSubtypeLongValue
	tsNodeParseSubtypePositional
)

func tsNodeParseFindOptionByShortName(shortName string, options []tsNodeParseOption) int {
	longNameIndex := -1
	for index, option := range options {
		if shortName == option.shortName { return index }
		if len([]rune(option.longName)) == 1 && shortName == option.longName { longNameIndex = index }
	}
	return longNameIndex
}

func tsNodeParseClassifyToken(argument string, options []tsNodeParseOption) tsNodeParseTokenSubtype {
	length := len([]rune(argument))
	if length == 2 && strings.HasPrefix(argument, "-") {
		if strings.HasPrefix(argument, "--") { return tsNodeParseSubtypeTerminator }
		return tsNodeParseSubtypeLoneShort
	}
	if length > 2 {
		if strings.HasPrefix(argument, "--") {
			if equal := strings.IndexByte(argument, '='); equal >= 3 { return tsNodeParseSubtypeLongValue }
			return tsNodeParseSubtypeLoneLong
		}
		if strings.HasPrefix(argument, "-") {
			runes := []rune(argument)
			shortName := string(runes[1])
			if index := tsNodeParseFindOptionByShortName(shortName, options); index >= 0 && options[index].valueType == "string" {
				return tsNodeParseSubtypeShortValue
			}
			return tsNodeParseSubtypeShortGroup
		}
	}
	return tsNodeParseSubtypePositional
}

func tsNodeParseOptionLikeValue(value string) bool {
	return len([]rune(value)) > 1 && strings.HasPrefix(value, "-")
}

type tsNodeParseOptionToken struct {
	index int
	name string
	raw string
	value tsValue
	inlineValue bool
	optionIndex int
	negative bool
	groupIndex int
	subtype tsNodeParseTokenSubtype
}

type tsNodeParseToken struct {
	kind tsNodeParseTokenKind
	index int
	value tsValue
	option tsNodeParseOptionToken
}

func tsNodeParseOwn(loop *tsLoop, value tsValue, name string) tsValue {
	if !tsObjectHasOwn(value, name) { return tsU }
	return tsGet(loop, value, tsStringReference(tsStringUTF8(name)))
}

func tsNodeParseFieldName(optionName, field string) string {
	if optionName == "" { return field }
	return "options." + optionName + "." + field
}

func tsNodeParseCheckString(name string, value tsValue) {
	if value.kind != tsStringKind { tsInternalValidatorInvalidType(name, "a string", value) }
}

func tsNodeParseCheckBoolean(name string, value tsValue) {
	if value.kind != tsBooleanKind { tsInternalValidatorInvalidType(name, "a boolean", value) }
}

func tsNodeParseCheckOptionDefault(option *tsNodeParseOption, value tsValue) {
	if option.multiple {
		if value.kind != tsArrayKind { tsInternalValidatorInvalidType(tsNodeParseFieldName(option.longName, "default"), "an array", value) }
		array := (*tsArray)(value.ref)
		for index := 0; index < array.length(); index++ {
			item := array.at(index)
			if option.valueType == "string" { tsNodeParseCheckString(tsNodeParseFieldName(option.longName, "default"), item) } else { tsNodeParseCheckBoolean(tsNodeParseFieldName(option.longName, "default"), item) }
		}
		return
	}
	if option.valueType == "string" { tsNodeParseCheckString(tsNodeParseFieldName(option.longName, "default"), value) } else { tsNodeParseCheckBoolean(tsNodeParseFieldName(option.longName, "default"), value) }
}

func tsNodeParseReadOptions(loop *tsLoop, value tsValue) []tsNodeParseOption {
	if tsNullish(value) { return nil }
	tsInternalValidatorObject(value, tsStringReference(tsStringUTF8("options")), tsU)
	keys := tsObjectKeys(value)
	options := make([]tsNodeParseOption, 0, keys.length())
	for _, keyValue := range keys.values {
		name := tsPropertyKey(keyValue)
		optionValue := tsGet(loop, value, tsStringReference(tsStringUTF8(name)))
		tsInternalValidatorObject(optionValue, tsStringReference(tsStringUTF8("options."+name)), tsU)
		option := tsNodeParseOption{longName:name, valueType:"boolean"}
		typeValue := tsNodeParseOwn(loop, optionValue, "type")
		tsNodeParseCheckString(tsNodeParseFieldName(name, "type"), typeValue)
		option.valueType = tsText(typeValue)
		if option.valueType != "boolean" && option.valueType != "string" {
			tsInternalValidatorInvalidValue(tsNodeParseFieldName(name, "type"), "must be either 'boolean' or 'string'", typeValue)
		}
		short := tsNodeParseOwn(loop, optionValue, "short")
		if !tsIsUndefined(short) {
			tsNodeParseCheckString(tsNodeParseFieldName(name, "short"), short)
			option.shortName = tsText(short)
			if len([]rune(option.shortName)) != 1 { tsInternalValidatorInvalidValue(tsNodeParseFieldName(name, "short"), "must be a single character", short) }
		}
		multiple := tsNodeParseOwn(loop, optionValue, "multiple")
		if !tsIsUndefined(multiple) { tsNodeParseCheckBoolean(tsNodeParseFieldName(name, "multiple"), multiple); option.multiple = tsTruthy(multiple) }
		defaultValue := tsNodeParseOwn(loop, optionValue, "default")
		if !tsIsUndefined(defaultValue) {
			tsNodeParseCheckOptionDefault(&option, defaultValue)
			option.defaultValue = defaultValue
			option.hasDefault = true
		}
		options = append(options, option)
	}
	return options
}

func tsNodeParseLongOption(name string, options []tsNodeParseOption) int {
	for index, option := range options { if name == option.longName { return index } }
	return -1
}

func tsNodeParseShortName(name string, options []tsNodeParseOption) string {
	index := tsNodeParseFindOptionByShortName(name, options)
	if index >= 0 { return options[index].longName }
	return name
}

func tsNodeParseOptionType(index int, options []tsNodeParseOption) string {
	if index < 0 { return "boolean" }
	return options[index].valueType
}

func tsNodeParseRuneSlice(value string, start, end int) string {
	runes := []rune(value)
	if start < 0 { start = 0 }; if start > len(runes) { start = len(runes) }
	if end < start { end = start }; if end > len(runes) { end = len(runes) }
	return string(runes[start:end])
}

func tsNodeParseRawName(option tsNodeParseOptionToken) string {
	if option.groupIndex >= 0 { return "-" + tsNodeParseRuneSlice(option.raw, option.groupIndex, option.groupIndex+1) }
	switch option.subtype {
	case tsNodeParseSubtypeShortValue: return tsNodeParseRuneSlice(option.raw, 0, 2)
	case tsNodeParseSubtypeLongValue:
		if equal := strings.IndexByte(option.raw, '='); equal >= 0 { return option.raw[:equal] }
	}
	return option.raw
}

func tsNodeParseError(code, message string) {
	panic(tsThrown{tsFSException("TypeError", code, message)})
}

func tsNodeParseUnknownOption(token tsNodeParseOptionToken, allowPositionals bool) {
	rawName := tsNodeParseRawName(token)
	if allowPositionals {
		tsNodeParseError("ERR_PARSE_ARGS_UNKNOWN_OPTION", "Unknown option '"+rawName+"'. To specify a positional argument starting with a '-', place it at the end of the command after '--', as in '-- \""+rawName+"\"")
	}
	tsNodeParseError("ERR_PARSE_ARGS_UNKNOWN_OPTION", "Unknown option '"+rawName+"'")
}

func tsNodeParseValidateOption(token tsNodeParseOptionToken, options []tsNodeParseOption, allowPositionals bool) {
	if token.optionIndex < 0 { tsNodeParseUnknownOption(token, allowPositionals) }
	option := options[token.optionIndex]
	if option.valueType == "string" && token.value.kind != tsStringKind {
		if token.negative { tsNodeParseUnknownOption(token, allowPositionals) }
		shortPrefix, separator := "", ""
		if option.shortName != "" { shortPrefix, separator = "-"+option.shortName, ", " }
		tsNodeParseError("ERR_PARSE_ARGS_INVALID_OPTION_VALUE", "Option '"+shortPrefix+separator+"--"+token.name+" <value>' argument missing")
	}
	if option.valueType == "boolean" && !tsIsUndefined(token.value) {
		shortPrefix, separator := "", ""
		if option.shortName != "" { shortPrefix, separator = "-"+option.shortName, ", " }
		tsNodeParseError("ERR_PARSE_ARGS_INVALID_OPTION_VALUE", "Option '"+shortPrefix+separator+"--"+token.name+"' does not take an argument")
	}
	if !token.inlineValue && token.value.kind == tsStringKind && tsNodeParseOptionLikeValue(tsText(token.value)) {
		rawName := tsNodeParseRawName(token)
		if strings.HasPrefix(rawName, "--") {
			tsNodeParseError("ERR_PARSE_ARGS_INVALID_OPTION_VALUE", "Option '"+rawName+"' argument is ambiguous.\nDid you forget to specify the option argument for '"+rawName+"'?\nTo specify an option argument starting with a dash use '"+rawName+"=-XYZ'.")
		}
		tsNodeParseError("ERR_PARSE_ARGS_INVALID_OPTION_VALUE", "Option '"+rawName+"' argument is ambiguous.\nDid you forget to specify the option argument for '"+rawName+"'?\nTo specify an option argument starting with a dash use '--"+token.name+"=-XYZ' or '"+rawName+"-XYZ'.")
	}
}
`
