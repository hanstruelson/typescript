package goemit

// RegexRuntime contains the Go-side object and String method behavior for
// ahead-of-time regular-expression matchers. Each matcher body is emitted from
// the V8-shaped regexp graph during TypeScript compilation.
const RegexRuntime = `
type tsRegExpMatcher func([]uint16, int, []int32) (int, int, bool)
type tsRegExpNamedCapture struct { name string; indices []int }
type tsRegExp struct {
	prototype tsValue
	metadata *tsObject
	source *tsString
	flags string
	matcher tsRegExpMatcher
	captureSlots int
	lastIndex int
	lastIndexValue tsValue
	customLastIndex bool
	properties map[string]tsValue
	order []string
	global, ignoreCase, multiline, dotAll, unicode, unicodeSets, sticky, hasIndices bool
	namedCaptures []tsRegExpNamedCapture
}

func tsRegExpFromValue(value tsValue) *tsRegExp {
	if value.kind == tsRegExpKind {
		return (*tsRegExp)(value.ref)
	}
	if value.kind == tsInstanceKind {
		return (*tsProperties)(value.ref).regexp
	}
	return nil
}

func tsRegExpExecReceiver(loop *tsLoop, receiver tsValue, input *tsString) tsValue {
	regexp := tsRegExpFromValue(receiver)
	if regexp == nil {
		tsPropertyFailure("RegExp.prototype.exec called on incompatible receiver")
	}
	if receiver.kind == tsInstanceKind {
		regexp.lastIndexValue = tsGet(loop, receiver, tsStringReference(tsStringUTF8("lastIndex")))
		regexp.customLastIndex = true
	}
	result := tsRegExpExecValue(loop, regexp, input)
	if receiver.kind == tsInstanceKind && (regexp.global || regexp.sticky) {
		tsSet(loop, receiver, tsStringReference(tsStringUTF8("lastIndex")), tsNumberValue(float64(regexp.lastIndex)))
	}
	return result
}

var tsRegExpIntrinsics struct {
	sync.Once
	modified atomic.Bool
	constructor *tsFunction
	prototype *tsObject
	compile    *tsFunction
}

func tsRegExpConstructor() tsValue {
	tsInitRegExp()
	return tsFunctionValue(tsRegExpIntrinsics.constructor)
}

func tsRegExpPrototype() tsValue {
	tsInitRegExp()
	return tsObjectValue(tsRegExpIntrinsics.prototype)
}

func tsInitRegExp() {
	tsRegExpIntrinsics.Do(func() {
		tsInitPrototypes()
		prototype := tsNewObject()
		constructor := tsFunc(func(loop *tsLoop, args ...tsValue) tsValue {
			return tsDynamicRegExp(loop, args...)
		})
		constructor.name = "RegExp"
		constructor.length = 2
		constructor.constructible = true
		constructor.construct = func(loop *tsLoop, args ...tsValue) tsValue {
			pattern := tsArg(args, 0)
			flags := tsArg(args, 1)
			if pattern.kind == tsRegExpKind && tsIsUndefined(flags) {
				original := (*tsRegExp)(pattern.ref)
				clone := *original
				clone.prototype = tsRegExpPrototype()
				clone.lastIndex = 0
				clone.lastIndexValue = tsU
				clone.customLastIndex = false
				clone.properties = nil
				clone.order = nil
				clone.metadata = nil
				return tsRegExpValue(&clone)
			}
			if pattern.kind == tsRegExpKind {
				// RegExp(P, F) reads P.source and coerces F before constructing.
				// Keep those observable operations even when Go code generation
				// cannot compile a genuinely dynamic pattern.
				_ = tsStringValueWithLoop(loop, tsGet(loop, pattern, tsStringReference(tsStringUTF8("source"))))
				flagText := tsStringValueWithLoop(loop, flags).String()
				var valid bool
				flagText, valid = tsRegExpCanonicalFlags(flagText)
				if !valid { tsRegExpSyntaxFailure("Invalid regular expression flags") }
				original := (*tsRegExp)(pattern.ref)
				originalSemantic := strings.ReplaceAll(strings.ReplaceAll(original.flags, "g", ""), "y", "")
				requestedSemantic := strings.ReplaceAll(strings.ReplaceAll(flagText, "g", ""), "y", "")
				if originalSemantic == requestedSemantic {
					clone := *original
					clone.prototype = tsRegExpPrototype()
					clone.flags = ""
					for _, flag := range "dgimsuvy" { if strings.ContainsRune(flagText, flag) { clone.flags += string(flag) } }
					clone.hasIndices = strings.ContainsRune(flagText, 'd')
					clone.global = strings.ContainsRune(flagText, 'g')
					clone.sticky = strings.ContainsRune(flagText, 'y')
					clone.lastIndex = 0
					clone.lastIndexValue = tsU
					clone.customLastIndex = false
					clone.properties = nil
					clone.order = nil
					clone.metadata = nil
					return tsRegExpValue(&clone)
				}
				// A generic flags getter can change flags during @@split or
				// @@matchAll. Recompile simple literal sources with the new
				// flags while keeping dynamic-pattern compilation out of runtime.
				if compiled, ok := tsRegExpLiteralSearch(original.source, flagText); ok {
					clone := *(*tsRegExp)(compiled.ref)
					clone.prototype = tsRegExpPrototype()
					return tsRegExpValue(&clone)
				}
				tsPropertyFailure("Dynamic RegExp construction is not supported; use a regular-expression literal")
			}
			patternText := tsStringUTF8("")
			if !tsIsUndefined(pattern) { patternText = tsStringValueWithLoop(loop, pattern) }
			flagText := ""
			if !tsIsUndefined(flags) { flagText = tsStringValueWithLoop(loop, flags).String() }
			var valid bool
			flagText, valid = tsRegExpCanonicalFlags(flagText)
			if !valid { tsRegExpSyntaxFailure("Invalid regular expression flags") }
			if regexp, ok := tsRegExpLiteralSearch(patternText, flagText); ok { return regexp }
			tsPropertyFailure("Dynamic RegExp syntax is not supported; use a regular-expression literal (received /" + patternText.String() + "/" + flagText + ")")
			return tsU
		}
		constructor.properties = tsNewObject()
		constructor.properties.set("prototype", tsObjectValue(prototype))
		constructor.properties.descriptors = map[string]*tsDescriptor{
			"prototype": {writable: false},
		}
		value := tsFunctionValue(constructor)
		prototype.set("constructor", value)
		prototype.descriptors = map[string]*tsDescriptor{
			"constructor": {writable: true, configurable: true},
		}
		tsRegExpIntrinsics.constructor = constructor
		tsRegExpIntrinsics.prototype = prototype

		constructor.properties.set("escape", tsNamedNativeMethod("escape", 1, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
			value := tsArg(args, 0)
			if value.kind != tsStringKind {
				tsPropertyFailure("RegExp.escape requires a string")
			}
			return tsStringReference(tsRegExpEscape((*tsString)(value.ref)))
		}))
		constructor.properties.descriptors["escape"] = &tsDescriptor{writable: true, configurable: true}
		species := tsPropertyKey(tsWellKnownSymbol("species"))
		constructor.properties.order = append(constructor.properties.order, species)
		constructor.properties.descriptors[species] = &tsDescriptor{
			accessor: true,
			getter: tsNamedNativeMethod("get [Symbol.species]", 0, func(_ *tsLoop, receiver tsValue, _ ...tsValue) tsValue { return receiver }),
			configurable: true,
		}

		methods := []struct {
			name string
			length int
			call func(*tsLoop, tsValue, ...tsValue) tsValue
		}{
			{"exec", 1, func(loop *tsLoop, receiver tsValue, args ...tsValue) tsValue {
				return tsRegExpExecReceiver(loop, receiver, tsStringValueWithLoop(loop, tsArg(args, 0)))
			}},
			{"test", 1, func(loop *tsLoop, receiver tsValue, args ...tsValue) tsValue {
				return tsBooleanValue(!tsNullish(tsRegExpExec(loop, receiver, tsStringValueWithLoop(loop, tsArg(args, 0)), tsU)))
			}},
			{"compile", 2, func(_ *tsLoop, _ tsValue, _ ...tsValue) tsValue {
				tsPropertyFailure("Dynamic RegExp compilation is not supported; compile-time constant patterns are supported")
				return tsU
			}},
			{"toString", 0, func(loop *tsLoop, receiver tsValue, _ ...tsValue) tsValue {
				if tsIsPrimitiveValue(receiver) { tsPropertyFailure("RegExp.prototype.toString called on incompatible receiver") }
				source := tsStringValueWithLoop(loop, tsGet(receiver, tsStringReference(tsStringUTF8("source"))))
				flags := tsStringValueWithLoop(loop, tsGet(receiver, tsStringReference(tsStringUTF8("flags"))))
				return tsStringReference(tsStringConcat(tsStringUTF8("/"), source, tsStringUTF8("/"), flags))
			}},
			{"[Symbol.match]", 1, func(loop *tsLoop, receiver tsValue, args ...tsValue) tsValue { return tsRegExpSymbolMethod(loop, receiver, "match", args...) }},
			{"[Symbol.matchAll]", 1, func(loop *tsLoop, receiver tsValue, args ...tsValue) tsValue { return tsRegExpSymbolMethod(loop, receiver, "matchAll", args...) }},
			{"[Symbol.replace]", 2, func(loop *tsLoop, receiver tsValue, args ...tsValue) tsValue { return tsRegExpSymbolMethod(loop, receiver, "replace", args...) }},
			{"[Symbol.search]", 1, func(loop *tsLoop, receiver tsValue, args ...tsValue) tsValue { return tsRegExpSymbolMethod(loop, receiver, "search", args...) }},
			{"[Symbol.split]", 2, func(loop *tsLoop, receiver tsValue, args ...tsValue) tsValue { return tsRegExpSymbolMethod(loop, receiver, "split", args...) }},
		}
		for _, entry := range methods {
			key := entry.name
			if strings.HasPrefix(key, "[Symbol.") {
				key = tsPropertyKey(tsWellKnownSymbol(strings.TrimSuffix(strings.TrimPrefix(key, "[Symbol."), "]")))
			}
			method := tsNamedNativeMethod(entry.name, entry.length, entry.call)
			prototype.set(key, method)
			if entry.name == "compile" {
				tsRegExpIntrinsics.compile = (*tsFunction)(method.ref)
			}
			prototype.descriptors[key] = &tsDescriptor{writable: true, configurable: true}
		}
		for _, flag := range []string{"hasIndices", "global", "ignoreCase", "multiline", "dotAll", "unicode", "unicodeSets", "sticky", "source", "flags"} {
			name := flag
			getter := tsNamedNativeMethod("get "+name, 0, func(loop *tsLoop, receiver tsValue, _ ...tsValue) tsValue {
				if tsSameValue(receiver, tsObjectValue(tsRegExpIntrinsics.prototype)) {
					if name == "source" { return tsStringReference(tsStringUTF8("(?:)")) }
					if name == "flags" { return tsStringReference(tsStringUnits(nil)) }
					return tsU
				}
				// The flags accessor is generic: its ordered property reads are
				// observable and work on any non-null receiver.
				if name == "flags" { return tsStringReference(tsRegExpFlagsFromProperties(loop, receiver)) }
				if tsRegExpFromValue(receiver) == nil { tsPropertyFailure("RegExp.prototype getter called on incompatible receiver") }
				return tsRegExpIntrinsicProperty(receiver, name)
			})
			prototype.descriptors[name] = &tsDescriptor{accessor: true, getter: getter, configurable: true}
			prototype.order = append(prototype.order, name)
		}
		for index := 1; index <= 9; index++ {
			capture := index
			key := "$" + strconv.Itoa(index)
			constructor.properties.set(key, tsU)
			constructor.properties.descriptors[key] = &tsDescriptor{accessor: true, getter: tsRegExpLegacyGetter(key, capture), configurable: true}
		}
		for _, key := range []string{"$&", "lastMatch", "$+", "lastParen", "$\x60", "leftContext", "$'", "rightContext", "$_", "input"} {
			constructor.properties.set(key, tsU)
			descriptor := &tsDescriptor{accessor: true, getter: tsRegExpLegacyGetter(key, 0), configurable: true}
			if key == "$_" || key == "input" {
				descriptor.setter = tsNamedNativeMethod("set "+key, 1, func(loop *tsLoop, receiver tsValue, args ...tsValue) tsValue {
					if !tsSameValue(receiver, tsRegExpConstructor()) { tsPropertyFailure("RegExp legacy accessor called on incompatible receiver") }
					if loop.regexpLastMatchInfo == nil { loop.regexpLastMatchInfo = &tsRegExpMatchInfo{} }
					loop.regexpLastMatchInfo.lastInput = tsStringValueWithLoop(loop, tsArg(args, 0))
					return tsU
				})
			}
			constructor.properties.descriptors[key] = descriptor
		}
	})
}

func tsRegExpLegacyGetter(key string, capture int) tsValue {
	return tsNamedNativeMethod("get "+key, 0, func(loop *tsLoop, receiver tsValue, _ ...tsValue) tsValue {
		if !tsSameValue(receiver, tsRegExpConstructor()) { tsPropertyFailure("RegExp legacy accessor called on incompatible receiver") }
		info := loop.regexpLastMatchInfo
		if key == "$_" || key == "input" {
			if info == nil || info.lastInput == nil { return tsStringReference(tsStringUnits(nil)) }
			return tsStringReference(info.lastInput)
		}
		if info == nil || info.lastSubject == nil { return tsStringReference(tsStringUnits(nil)) }
		switch key {
		case "$&", "lastMatch": return tsStringReference(tsRegExpGenericCaptureGetter(info, 0, nil))
		case "$+", "lastParen":
			for index := len(info.captures)/2 - 1; index > 0; index-- { if tsRegExpIsMatchedCapture(info, index) { return tsStringReference(tsRegExpGenericCaptureGetter(info, index, nil)) } }
			return tsStringReference(tsStringUnits(nil))
		case "$\x60", "leftContext": start := int(info.captures[0]); if start < 0 { start = 0 }; return tsStringReference(info.lastSubject.slice(0, start))
		case "$'", "rightContext": end := int(info.captures[1]); if end < 0 || end > len(info.lastSubject.units) { end = len(info.lastSubject.units) }; return tsStringReference(info.lastSubject.slice(end, len(info.lastSubject.units)))
		default: return tsStringReference(tsRegExpGenericCaptureGetter(info, capture, nil))
		}
	})
}

func tsRegExpMetadata(regexp *tsRegExp) *tsObject {
	if regexp.metadata == nil {
		object := tsNewObject()
		object.values = make(map[string]tsValue, len(regexp.properties)+1)
		object.order = make([]string, 0, len(regexp.order)+1)
		object.order = append(object.order, "lastIndex")
		lastIndex := tsNumberValue(float64(regexp.lastIndex))
		if regexp.customLastIndex { lastIndex = regexp.lastIndexValue }
		object.values["lastIndex"] = lastIndex
		object.descriptors = map[string]*tsDescriptor{"lastIndex": {writable: true}}
		for _, name := range regexp.order {
			object.order = append(object.order, name)
			if property, ok := regexp.properties[name]; ok { object.values[name] = property }
			object.descriptors[name] = &tsDescriptor{writable: true, enumerable: true, configurable: true}
		}
		regexp.metadata = object
		regexp.properties = object.values
	}
	return regexp.metadata
}

func tsRegExpFlagsFromProperties(loop *tsLoop, receiver tsValue) *tsString {
	var flags strings.Builder
	for _, entry := range []struct { key, flag string }{{"hasIndices", "d"}, {"global", "g"}, {"ignoreCase", "i"}, {"multiline", "m"}, {"dotAll", "s"}, {"unicode", "u"}, {"unicodeSets", "v"}, {"sticky", "y"}} {
		if tsTruthy(tsGet(receiver, tsStringReference(tsStringUTF8(entry.key)))) { flags.WriteString(entry.flag) }
	}
	return tsStringUTF8(flags.String())
}

func tsRegExpEscape(input *tsString) *tsString {
	var out strings.Builder
	writeHex := func(prefix string, value uint16) { fmt.Fprintf(&out, "\\%s%04x", prefix, value) }
	for index := 0; index < len(input.units); index++ {
		unit := input.units[index]
		if index == 0 && (unit >= 'a' && unit <= 'z' || unit >= 'A' && unit <= 'Z' || unit >= '0' && unit <= '9') {
			fmt.Fprintf(&out, "\\x%02x", unit)
			continue
		}
		if strings.ContainsRune("^$\\.*+?()[]{}|/", rune(unit)) { out.WriteByte('\\'); out.WriteRune(rune(unit)); continue }
		switch unit {
		case ',', '-', '=', '<', '>', '#', '&', '!', '%', ':', ';', '@', '~', '\'', '\x60', '"':
			fmt.Fprintf(&out, "\\x%02x", unit)
		case 0x0000:
			out.WriteString("\\x00")
		case 0x0009:
			out.WriteString("\\t")
		case 0x000a:
			out.WriteString("\\n")
		case 0x000b:
			out.WriteString("\\v")
		case 0x000c:
			out.WriteString("\\f")
		case 0x000d:
			out.WriteString("\\r")
		case 0x0020:
			out.WriteString("\\x20")
		case 0x2028, 0x2029:
			writeHex("u", unit)
		default:
			if unit >= 0xd800 && unit <= 0xdbff {
				if index+1 < len(input.units) && input.units[index+1] >= 0xdc00 && input.units[index+1] <= 0xdfff { out.WriteRune(utf16.DecodeRune(rune(unit), rune(input.units[index+1]))); index++; continue }
				writeHex("u", unit)
			} else if unit >= 0xdc00 && unit <= 0xdfff {
				writeHex("u", unit)
			} else if tsStringWhitespace(unit) {
				if unit <= 0xff { fmt.Fprintf(&out, "\\x%02x", unit) } else { writeHex("u", unit) }
			} else {
				out.WriteRune(rune(unit))
			}
		}
	}
	return tsStringUTF8(out.String())
}

// Each event loop owns a lazy static vector. Nested calls and large capture
// sets use independent pooled vectors, matching V8's checked-out/dynamic split.
// Oversized vectors are released instead of being retained indefinitely.
const tsRegExpMaxPooledResultVectorSlots = 4096
const tsRegExpStaticResultOffsetsVectorSize = 32
var tsRegExpResultVectorPool = sync.Pool{New: func() any { return make([]int32, 0, 32) }}
var tsRegExpEmptyCaptureString = &tsString{}

type tsRegExpResultVectorScope struct {
	loop *tsLoop
	value []int32
	active bool
	static bool
}

func (scope *tsRegExpResultVectorScope) Initialize(loop *tsLoop, size int) []int32 {
	if scope.active {
		panic("RegExp result vector scope initialized twice")
	}
	if size < 0 {
		panic("negative RegExp result vector size")
	}
	scope.active = true
	scope.loop = loop
	if loop != nil && size <= tsRegExpStaticResultOffsetsVectorSize {
		if !loop.regexpStaticResultOffsetsVectorInitialized {
			loop.regexpStaticResultOffsetsVector = make([]int32, tsRegExpStaticResultOffsetsVectorSize)
			loop.regexpStaticResultOffsetsVectorInitialized = true
		}
		if loop.regexpStaticResultOffsetsVector != nil {
			scope.static = true
			scope.value = loop.regexpStaticResultOffsetsVector[:size]
			loop.regexpStaticResultOffsetsVector = nil
		}
	}
	if scope.value == nil {
		scope.value = tsRegExpResultVectorAllocate(size)
	}
	for index := range scope.value { scope.value[index] = -1 }
	return scope.value
}

func (scope *tsRegExpResultVectorScope) Value() []int32 {
	if !scope.active {
		panic("RegExp result vector scope is not initialized")
	}
	return scope.value
}

func (scope *tsRegExpResultVectorScope) Close() {
	if !scope.active {
		return
	}
	if scope.static {
		scope.loop.regexpStaticResultOffsetsVector = scope.value[:cap(scope.value)]
	} else {
		tsRegExpResultVectorFree(scope.value)
	}
	scope.value = nil
	scope.active = false
	scope.static = false
	scope.loop = nil
}

func tsRegExpResultVectorAllocate(size int) []int32 {
	value := tsRegExpResultVectorPool.Get().([]int32)
	if cap(value) < size {
		value = make([]int32, size)
	} else {
		value = value[:size]
	}
	return value
}

func tsRegExpResultVectorFree(value []int32) {
	if cap(value) == 0 || cap(value) > tsRegExpMaxPooledResultVectorSlots {
		return
	}
	for index := range value {
		value[index] = 0
	}
	tsRegExpResultVectorPool.Put(value[:0])
}

type tsRegExpMatchInfo struct {
	lastSubject *tsString
	lastInput *tsString
	captures []int32
}

// tsRegExpSetLastMatchInfo mirrors RegExp::SetLastMatchInfo. Captures are
// copied before the result-vector scope is released, and a loop reuses its
// match-info allocation when it already has enough room.
func tsRegExpSetLastMatchInfo(loop *tsLoop, subject *tsString, captureRegisterCount int, match []int32) {
	if loop == nil { return }
	registerCount := captureRegisterCount
	info := loop.regexpLastMatchInfo
	if info == nil || cap(info.captures) < registerCount {
		info = &tsRegExpMatchInfo{captures: make([]int32, registerCount)}
	} else {
		info.captures = info.captures[:registerCount]
	}
	if match != nil {
		copy(info.captures, match[:min(registerCount, len(match))])
	}
	info.lastSubject = subject
	info.lastInput = subject
	loop.regexpLastMatchInfo = info
}

func tsRegExpCaptureStartIndex(capture int) int { return capture * 2 }
func tsRegExpCaptureEndIndex(capture int) int { return capture*2 + 1 }

// GenericCaptureGetter mirrors V8's last-match capture accessor: an absent or
// unmatched capture produces the empty string and reports ok=false.
func tsRegExpGenericCaptureGetter(info *tsRegExpMatchInfo, capture int, ok *bool) *tsString {
	startIndex := tsRegExpCaptureStartIndex(capture)
	if info == nil || capture < 0 || startIndex >= len(info.captures) {
		if ok != nil { *ok = false }
		return tsRegExpEmptyCaptureString
	}
	endIndex := tsRegExpCaptureEndIndex(capture)
	start, end := int(info.captures[startIndex]), int(info.captures[endIndex])
	if start == -1 || end == -1 || start < 0 || end < start || end > len(info.lastSubject.units) {
		if ok != nil { *ok = false }
		return tsRegExpEmptyCaptureString
	}
	if ok != nil { *ok = true }
	return info.lastSubject.slice(start, end)
}

func tsRegExpIsMatchedCapture(info *tsRegExpMatchInfo, capture int) bool {
	if capture == -1 || capture < 0 || info == nil { return false }
	startIndex := tsRegExpCaptureStartIndex(capture)
	endIndex := tsRegExpCaptureEndIndex(capture)
	return endIndex < len(info.captures) && info.captures[startIndex] != -1 && info.captures[endIndex] != -1
}

func tsRegExpLastIndexKey() tsValue { return tsStringReference(tsStringUTF8("lastIndex")) }

func tsRegExpHasInitialMap(receiver tsValue) bool {
	if receiver.kind != tsRegExpKind { return false }
	return len((*tsRegExp)(receiver.ref).properties) == 0
}

// GetLastIndex and SetLastIndex use the internal RegExp slot for an ordinary
// instance and ordinary property access for receivers with a customized shape.
func tsRegExpGetLastIndex(loop *tsLoop, receiver tsValue) tsValue {
	if tsRegExpHasInitialMap(receiver) {
		regexp := (*tsRegExp)(receiver.ref)
		if !regexp.customLastIndex { return tsNumberValue(float64(regexp.lastIndex)) }
	}
	return tsGet(receiver, tsRegExpLastIndexKey())
}

func tsRegExpSetLastIndex(loop *tsLoop, receiver, value tsValue) tsValue {
	if tsRegExpHasInitialMap(receiver) {
		tsRegExpSet(receiver, value, "lastIndex")
		return receiver
	}
	return tsSet(receiver, tsRegExpLastIndexKey(), value)
}

// tsRegExpCompileKnown is the AOT form of the Annex B compile operation. The
// emitter supplies a matcher generated from a constant pattern; user overrides
// of compile still receive an ordinary call with the original arguments.
func tsRegExpCompileKnown(loop *tsLoop, receiver, method tsValue, args []tsValue, compiled tsValue) tsValue {
	tsInitRegExp()
	if method.kind != tsFunctionKind || (*tsFunction)(method.ref) != tsRegExpIntrinsics.compile {
		return tsCallReceiver(loop, method, receiver, args...)
	}
	if receiver.kind != tsRegExpKind {
		tsPropertyFailure("RegExp.prototype.compile called on incompatible receiver")
	}
	if compiled.kind != tsRegExpKind {
		tsPropertyFailure("RegExp compile pattern was not statically compiled")
	}
	regexp := (*tsRegExp)(receiver.ref)
	template := (*tsRegExp)(compiled.ref)
	regexp.source = template.source
	regexp.flags = template.flags
	regexp.matcher = template.matcher
	regexp.captureSlots = template.captureSlots
	regexp.global = template.global
	regexp.ignoreCase = template.ignoreCase
	regexp.multiline = template.multiline
	regexp.dotAll = template.dotAll
	regexp.unicode = template.unicode
	regexp.unicodeSets = template.unicodeSets
	regexp.sticky = template.sticky
	regexp.hasIndices = template.hasIndices
	regexp.namedCaptures = template.namedCaptures
	_ = args // Argument expressions have already run in JavaScript evaluation order.
	tsRegExpSetLastIndex(loop, receiver, tsNumberValue(0))
	return receiver
}

func tsRegExpSetInternalLastIndex(regexp *tsRegExp, value int) {
	if regexp.metadata != nil {
		if descriptor := regexp.metadata.descriptors["lastIndex"]; descriptor != nil && (descriptor.accessor || !descriptor.writable) {
			tsPropertyFailure("Cannot assign to read-only RegExp lastIndex")
		}
	}
	regexp.lastIndex = value
	regexp.lastIndexValue = tsU
	regexp.customLastIndex = false
	if regexp.metadata != nil { regexp.metadata.values["lastIndex"] = tsNumberValue(float64(value)) }
}

// IsUnmodifiedRegExp is conservative where V8 checks maps and protectors.
// This runtime has no mutable RegExp prototype object, so an unshadowed own
// property shape and a nonnegative small integer lastIndex are the fast case.
func tsRegExpIsUnmodifiedRegExp(value tsValue) bool {
	if value.kind != tsRegExpKind { return false }
	if tsRegExpIntrinsics.modified.Load() { return false }
	regexp := (*tsRegExp)(value.ref)
	if len(regexp.properties) != 0 || regexp.customLastIndex || regexp.metadata != nil { return false }
	// The initial shape stores lastIndex as a native integer. Read it directly
	// here so checking the RegExp fast-path protector does not box the value.
	number := float64(regexp.lastIndex)
	return number >= 0 && number <= 2147483647 && math.Trunc(number) == number
}

// AdvanceStringIndex operates in UTF-16 code units and skips a valid
// surrogate pair only when Unicode matching is enabled.
func tsRegExpAdvanceStringIndex(input *tsString, index uint64, unicodeMode bool) uint64 {
	if unicodeMode && index < uint64(len(input.units)) {
		first := input.units[index]
		if first >= 0xd800 && first <= 0xdbff && index+1 < uint64(len(input.units)) {
			second := input.units[index+1]
			if second >= 0xdc00 && second <= 0xdfff { return index + 2 }
		}
	}
	return index + 1
}

func tsRegExpToIntegerOrInfinity(loop *tsLoop, value tsValue) float64 {
	number := tsNumber(tsPrimitiveHint(loop, value, false))
	if math.IsNaN(number) || number == 0 { return 0 }
	return math.Trunc(number)
}

func tsRegExpSetAdvancedStringIndex(loop *tsLoop, receiver tsValue, input *tsString, unicodeMode bool) tsValue {
	index := uint64(tsToLength(loop, tsRegExpGetLastIndex(loop, receiver)))
	next := tsRegExpAdvanceStringIndex(input, index, unicodeMode)
	return tsRegExpSetLastIndex(loop, receiver, tsNumberValue(float64(next)))
}

func tsNewCompiledRegExp(source, flags *tsString, matcher tsRegExpMatcher, captureSlots int, namedCaptures []tsRegExpNamedCapture) tsValue {
	flagText := flags.String()
	canonicalFlags := make([]byte, 0, len(flagText))
	for _, flag := range "dgimsuvy" { if strings.ContainsRune(flagText, flag) { canonicalFlags = append(canonicalFlags, byte(flag)) } }
	flagText = string(canonicalFlags)
	regexp := &tsRegExp{prototype: tsRegExpPrototype(), source: source, flags: flagText, matcher: matcher, captureSlots: captureSlots, namedCaptures: namedCaptures}
	for _, flag := range flagText { switch flag {
	case 'd': regexp.hasIndices = true
	case 'g': regexp.global = true
	case 'i': regexp.ignoreCase = true
	case 'm': regexp.multiline = true
	case 's': regexp.dotAll = true
	case 'u': regexp.unicode = true
	case 'v': regexp.unicodeSets = true
	case 'y': regexp.sticky = true
	} }
	return tsRegExpValue(regexp)
}
func tsDynamicRegExp(loop *tsLoop, args ...tsValue) tsValue {
	tsInitRegExp()
	pattern, flags := tsArg(args, 0), tsArg(args, 1)
	if pattern.kind == tsRegExpKind {
		original := (*tsRegExp)(pattern.ref)
		if tsIsUndefined(flags) {
			clone := *original
			clone.prototype = tsRegExpPrototype()
			clone.lastIndex = 0
			clone.lastIndexValue = tsU
			clone.customLastIndex = false
			clone.properties = nil
			clone.order = nil
			clone.metadata = nil
			return tsRegExpValue(&clone)
		}
		// Preserve observable source and flags coercions for RegExp(pattern, flags).
		_ = tsStringValueWithLoop(loop, tsGet(loop, pattern, tsStringReference(tsStringUTF8("source"))))
		flagText := tsStringValueWithLoop(loop, flags).String()
		var valid bool
		flagText, valid = tsRegExpCanonicalFlags(flagText)
		if !valid { tsRegExpSyntaxFailure("Invalid regular expression flags") }
		originalSemantic := strings.ReplaceAll(strings.ReplaceAll(original.flags, "g", ""), "y", "")
		requestedSemantic := strings.ReplaceAll(strings.ReplaceAll(flagText, "g", ""), "y", "")
		if originalSemantic == requestedSemantic {
			clone := *original
			clone.prototype = tsRegExpPrototype()
			clone.flags = flagText
			clone.hasIndices = strings.ContainsRune(flagText, 'd')
			clone.global = strings.ContainsRune(flagText, 'g')
			clone.sticky = strings.ContainsRune(flagText, 'y')
			clone.lastIndex = 0
			clone.lastIndexValue = tsU
			clone.customLastIndex = false
			clone.properties = nil
			clone.order = nil
			clone.metadata = nil
			return tsRegExpValue(&clone)
		}
		if compiled, ok := tsRegExpLiteralSearch(original.source, flagText); ok {
			clone := *(*tsRegExp)(compiled.ref)
			clone.prototype = tsRegExpPrototype()
			return tsRegExpValue(&clone)
		}
		tsPropertyFailure("Dynamic RegExp construction is not supported; use a regular-expression literal")
	}
	patternText := tsStringUTF8("")
	if !tsIsUndefined(pattern) { patternText = tsStringValueWithLoop(loop, pattern) }
	flagText := ""
	if !tsIsUndefined(flags) { flagText = tsStringValueWithLoop(loop, flags).String() }
	flagText, valid := tsRegExpCanonicalFlags(flagText)
	if !valid { tsRegExpSyntaxFailure("Invalid regular expression flags") }
	if regexp, ok := tsRegExpLiteralSearch(patternText, flagText); ok { return regexp }
	tsPropertyFailure("Dynamic RegExp syntax is not supported; use a regular-expression literal (received /" + patternText.String() + "/" + flagText + ")")
	return tsU
}

func tsRegExpCanonicalFlags(flags string) (string, bool) {
	seen := [128]bool{}
	for _, flag := range flags {
		if flag >= 128 || !strings.ContainsRune("dgimsuvy", flag) || seen[flag] { return "", false }
		seen[flag] = true
	}
	if seen['u'] && seen['v'] { return "", false }
	var canonical strings.Builder
	for _, flag := range "dgimsuvy" { if seen[flag] { canonical.WriteRune(flag) } }
	return canonical.String(), true
}

func tsRegExpSyntaxFailure(message string) {
	panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"SyntaxError", message:message})})
}
func tsRegExpSource(regexp *tsRegExp) *tsString {
	source := regexp.source
	if len(source.units) == 0 {
		return tsStringUTF8("(?:)")
	}

	// V8's EscapeRegExpSource first sizes the result and returns the original
	// string when no escaping is needed. Keep that allocation-free path for
	// the common case where .source is read from an ordinary pattern.
	outputLength := 0
	inCharacterClass := false
	needsEscapes := false
	for index := 0; index < len(source.units); index++ {
		character := source.units[index]
		if character == '\\' {
			if index+1 < len(source.units) {
				next := source.units[index+1]
				if tsRegExpIsLineTerminator(next) {
					needsEscapes = true
					continue
				}
				outputLength += 2
				index++
			} else {
				outputLength++
			}
			continue
		}
		switch character {
		case '/':
			outputLength++
			if !inCharacterClass {
				outputLength++
				needsEscapes = true
			}
		case '[':
			inCharacterClass = true
			outputLength++
		case ']':
			inCharacterClass = false
			outputLength++
		case '\n', '\r':
			outputLength += 2
			needsEscapes = true
		case 0x2028, 0x2029:
			outputLength += 6
			needsEscapes = true
		default:
			outputLength++
		}
	}
	if !needsEscapes {
		return source
	}

	output := make([]uint16, 0, outputLength)
	inCharacterClass = false
	for index := 0; index < len(source.units); index++ {
		character := source.units[index]
		if character == '\\' {
			if index+1 < len(source.units) && tsRegExpIsLineTerminator(source.units[index+1]) {
				continue
			}
			output = append(output, character)
			if index+1 < len(source.units) {
				index++
				output = append(output, source.units[index])
			}
			continue
		}
		switch character {
		case '/':
			if !inCharacterClass {
				output = append(output, '\\')
			}
			output = append(output, character)
		case '[':
			inCharacterClass = true
			output = append(output, character)
		case ']':
			inCharacterClass = false
			output = append(output, character)
		case '\n':
			output = append(output, '\\', 'n')
		case '\r':
			output = append(output, '\\', 'r')
		case 0x2028:
			output = append(output, '\\', 'u', '2', '0', '2', '8')
		case 0x2029:
			output = append(output, '\\', 'u', '2', '0', '2', '9')
		default:
			output = append(output, character)
		}
	}
	return tsStringUnits(output)
}

func tsRegExpIsLineTerminator(character uint16) bool {
	return character == '\n' || character == '\r' || character == 0x2028 || character == 0x2029
}
func tsRegExpText(regexp *tsRegExp) *tsString { return tsStringConcat(tsStringUTF8("/"), tsRegExpSource(regexp), tsStringUTF8("/"+regexp.flags)) }
func tsRegExpFindInto(regexp *tsRegExp, input *tsString, start int, captures []int32) (int, int, bool) {
	if start < 0 || start > len(input.units) { return 0, 0, false }
	searchStart := start
	if (regexp.unicode || regexp.unicodeSets) && (regexp.global || regexp.sticky) && searchStart > 0 && searchStart < len(input.units) && input.units[searchStart] >= 0xdc00 && input.units[searchStart] <= 0xdfff && input.units[searchStart-1] >= 0xd800 && input.units[searchStart-1] <= 0xdbff {
		searchStart--
	}
	matchStart, matchEnd, ok := regexp.matcher(input.units, searchStart, captures)
	if !ok {
		if matchStart == -1 && matchEnd == -1 {
			// The generated matcher uses this sentinel when its V8-sized
			// backtracking stack reaches the native stack limit.
			panic(tsThrown{tsErrorValue(&tsRuntimeError{name: "RangeError", message: "Maximum call stack size exceeded"})})
		}
		return 0, 0, false
	}
	if regexp.sticky {
		if matchStart != searchStart { return 0, 0, false }
	}
	return matchStart, matchEnd, true
}
func tsRegExpFind(regexp *tsRegExp, input *tsString, start int) (int, int, []int, bool) {
	scope := &tsRegExpResultVectorScope{}
	captures := scope.Initialize(nil, regexp.captureSlots)
	defer scope.Close()
	matchStart, matchEnd, ok := tsRegExpFindInto(regexp, input, start, captures)
	if !ok { return 0, 0, nil, false }
	converted := make([]int, len(captures))
	for index, capture := range captures { converted[index] = int(capture) }
	return matchStart, matchEnd, converted, true
}

func tsRegExpExecValue(loop *tsLoop, regexp *tsRegExp, input *tsString) tsValue {
	start := 0
	lastIndex := regexp.lastIndex
	if regexp.customLastIndex { lastIndex = tsToLength(loop, regexp.lastIndexValue) }
	if regexp.global || regexp.sticky { start = lastIndex }
	scope := &tsRegExpResultVectorScope{}
	captures := scope.Initialize(loop, regexp.captureSlots)
	defer scope.Close()
	matchStart, matchEnd, ok := tsRegExpFindInto(regexp, input, start, captures)
	if !ok {
		if regexp.global || regexp.sticky { tsRegExpSetInternalLastIndex(regexp, 0) }
		return tsNull
	}
	if regexp.global || regexp.sticky { tsRegExpSetInternalLastIndex(regexp, matchEnd) }
	matchInfo := &tsRegExpMatchInfo{lastSubject: input, lastInput: input, captures: captures}
	tsRegExpSetLastMatchInfo(loop, input, len(captures), captures)
	array := &tsArray{values: make([]tsValue, len(captures)/2), descriptors: make(map[string]*tsDescriptor, len(captures)/2+3)}
	var indices *tsArray
	if regexp.hasIndices { indices = &tsArray{values: make([]tsValue, len(array.values)), descriptors: make(map[string]*tsDescriptor, len(array.values)+1)} }
	for capture := range array.values {
		array.descriptors[strconv.Itoa(capture)] = &tsDescriptor{writable: true, enumerable: true, configurable: true}
		if indices != nil { indices.descriptors[strconv.Itoa(capture)] = &tsDescriptor{writable: true, enumerable: true, configurable: true} }
		var matched bool
		value := tsRegExpGenericCaptureGetter(matchInfo, capture, &matched)
		if !matched { array.values[capture] = tsU; if indices != nil { indices.values[capture] = tsU }; continue }
		startIndex, endIndex := int(captures[tsRegExpCaptureStartIndex(capture)]), int(captures[tsRegExpCaptureEndIndex(capture)])
		array.values[capture] = tsStringReference(value)
		if indices != nil {
			pair := &tsArray{values: []tsValue{tsNumberValue(float64(startIndex)), tsNumberValue(float64(endIndex))}, descriptors: map[string]*tsDescriptor{
				"0": {writable: true, enumerable: true, configurable: true},
				"1": {writable: true, enumerable: true, configurable: true},
			}}
			indices.values[capture] = tsArrayValue(pair)
		}
	}
	groupsValue := tsU
	var groups *tsObject
	if len(regexp.namedCaptures) != 0 { groups = tsNewObject(); groups.prototype = tsNull; groupsValue = tsObjectValue(groups) }
	var indexGroups *tsObject
	if indices != nil && len(regexp.namedCaptures) != 0 { indexGroups = tsNewObject(); indexGroups.prototype = tsNull }
	for _, namedCapture := range regexp.namedCaptures {
		name, captureIndices := namedCapture.name, namedCapture.indices
		value := tsU
		var indexValue tsValue = tsU
		for _, captureIndex := range captureIndices {
			if tsRegExpIsMatchedCapture(matchInfo, captureIndex) {
				value = tsStringReference(tsRegExpGenericCaptureGetter(matchInfo, captureIndex, nil))
				if indices != nil { indexValue = indices.values[captureIndex] }
				break
			}
		}
		if groups != nil { groups.set(name, value) }
		if indexGroups != nil { indexGroups.set(name, indexValue) }
	}
	array.properties = map[string]tsValue{"index": tsNumberValue(float64(matchStart)), "input": tsStringReference(input), "groups": groupsValue}
	array.order = []string{"index", "input", "groups"}
	if indices != nil {
		indicesGroups := tsU
		if indexGroups != nil { indicesGroups = tsObjectValue(indexGroups) }
		indices.properties = map[string]tsValue{"groups": indicesGroups}
		indices.order = []string{"groups"}
		array.properties["indices"] = tsArrayValue(indices)
		array.order = append(array.order, "indices")
	}
	return tsArrayValue(array)
}

// Global string matching only needs the full-match text. The ordinary-RegExp
// path can therefore consume the V8-style offset vector directly instead of
// allocating a result array and named-group objects for every match. An own
// exec override or additional own properties stay on the generic path.
func tsRegExpMatchGlobalFast(loop *tsLoop, receiver tsValue, regexp *tsRegExp, input *tsString) tsValue {
	scope := &tsRegExpResultVectorScope{}
	captures := scope.Initialize(loop, regexp.captureSlots)
	defer scope.Close()
	out := &tsArray{}
	for {
		start := tsToLength(loop, tsRegExpGetLastIndex(loop, receiver))
		matchStart, matchEnd, ok := tsRegExpFindInto(regexp, input, start, captures)
		if !ok {
			tsRegExpSetInternalLastIndex(regexp, 0)
			break
		}
		tsRegExpSetLastMatchInfo(loop, input, len(captures), captures)
		tsRegExpSetInternalLastIndex(regexp, matchEnd)
		out.values = append(out.values, tsStringReference(input.slice(matchStart, matchEnd)))
		if matchStart == matchEnd {
			tsRegExpSetAdvancedStringIndex(loop, receiver, input, regexp.unicode || regexp.unicodeSets)
		}
	}
	if len(out.values) == 0 { return tsNull }
	return tsArrayValue(out)
}

func tsRegExpGet(loop *tsLoop, value tsValue, name string) tsValue {
	regexp := (*tsRegExp)(value.ref)
	if regexp.metadata != nil {
		if descriptor := regexp.metadata.descriptors[name]; descriptor != nil && descriptor.accessor {
			if tsIsUndefined(descriptor.getter) { return tsU }
			return tsCallReceiver(loop, descriptor.getter, value)
		}
	}
	if property, ok := regexp.properties[name]; ok { return property }
	// lastIndex is an own, non-configurable data property even after the
	// RegExp prototype has been modified. Keep that own slot visible on the
	// observable path as well as on the initial-map fast path.
	if name == "lastIndex" {
		if regexp.customLastIndex { return regexp.lastIndexValue }
		return tsNumberValue(float64(regexp.lastIndex))
	}
	if !tsRegExpIntrinsics.modified.Load() { switch name {
	case "lastIndex": if regexp.customLastIndex { return regexp.lastIndexValue }; return tsNumberValue(float64(regexp.lastIndex))
	case "source": return tsStringReference(tsRegExpSource(regexp))
	case "flags":
		// RegExp.prototype.flags is generic. Even on an ordinary RegExp, its
		// getter must read the individual flag properties so own accessors and
		// prototype changes remain observable in the specified order.
		return tsStringReference(tsRegExpFlagsFromProperties(loop, value))
	case "global": return tsBooleanValue(regexp.global)
	case "hasIndices": return tsBooleanValue(regexp.hasIndices)
	case "ignoreCase": return tsBooleanValue(regexp.ignoreCase)
	case "multiline": return tsBooleanValue(regexp.multiline)
	case "dotAll": return tsBooleanValue(regexp.dotAll)
	case "unicode": return tsBooleanValue(regexp.unicode)
	case "unicodeSets": return tsBooleanValue(regexp.unicodeSets)
	case "sticky": return tsBooleanValue(regexp.sticky)
	} }
	if !tsIsUndefined(regexp.prototype) { return tsPrototypeLookup(loop, value, regexp.prototype, name) }
	return tsU
}

func tsRegExpIntrinsicProperty(value tsValue, name string) tsValue {
	regexp := tsRegExpFromValue(value)
	if regexp == nil {
		tsPropertyFailure("RegExp.prototype getter called on incompatible receiver")
	}
	switch name {
	case "source": return tsStringReference(tsRegExpSource(regexp))
	case "global": return tsBooleanValue(regexp.global)
	case "hasIndices": return tsBooleanValue(regexp.hasIndices)
	case "ignoreCase": return tsBooleanValue(regexp.ignoreCase)
	case "multiline": return tsBooleanValue(regexp.multiline)
	case "dotAll": return tsBooleanValue(regexp.dotAll)
	case "unicode": return tsBooleanValue(regexp.unicode)
	case "unicodeSets": return tsBooleanValue(regexp.unicodeSets)
	case "sticky": return tsBooleanValue(regexp.sticky)
	}
	return tsU
}

// RegExpExec implements the observable exec lookup and result validation.
// Ordinary instances can bypass the property lookup; customized receivers take
// the general callable path, while non-callable exec falls back to the built-in
// algorithm only for RegExp instances.
func tsRegExpExec(loop *tsLoop, receiver tsValue, input *tsString, exec tsValue) tsValue {
	if tsIsUndefined(exec) {
		if tsRegExpIsUnmodifiedRegExp(receiver) {
			return tsRegExpExecValue(loop, (*tsRegExp)(receiver.ref), input)
		}
		exec = tsGet(receiver, tsStringReference(tsStringUTF8("exec")))
	}
	if exec.kind == tsFunctionKind {
		result := tsCallReceiver(loop, exec, receiver, tsStringReference(input))
		if result.kind != tsNullKind && tsIsPrimitiveValue(result) {
			tsPropertyFailure("RegExp exec method must return an object or null")
		}
		return result
	}
	if tsRegExpFromValue(receiver) == nil {
		tsPropertyFailure("RegExp.prototype.exec called on incompatible receiver")
	}
	return tsRegExpExecReceiver(loop, receiver, input)
}

func tsRegExpSet(value, item tsValue, name string) tsValue {
	regexp := (*tsRegExp)(value.ref)
	if regexp.metadata != nil {
		if descriptor := regexp.metadata.descriptors[name]; descriptor != nil {
			if descriptor.accessor {
				if tsIsUndefined(descriptor.setter) { tsPropertyFailure("RegExp property has no setter") }
				tsCallReceiver(loop, descriptor.setter, value, item)
				return item
			}
			if !descriptor.writable { tsPropertyFailure("Cannot assign to read-only RegExp property") }
			if name == "lastIndex" {
				regexp.lastIndexValue = item
				regexp.customLastIndex = true
			}
			if regexp.properties == nil { regexp.properties = make(map[string]tsValue) }
			regexp.properties[name] = item
			return item
		}
	}
	if name == "lastIndex" {
		regexp.lastIndexValue = item
		regexp.customLastIndex = true
		if regexp.metadata != nil { regexp.metadata.values[name] = item }
		return item
	}
	if name == "source" || name == "flags" || name == "global" || name == "hasIndices" || name == "ignoreCase" || name == "multiline" || name == "dotAll" || name == "unicode" || name == "unicodeSets" || name == "sticky" { tsPropertyFailure("Cannot assign to read-only RegExp property") }
	if regexp.properties == nil { regexp.properties = make(map[string]tsValue) }
	if _, exists := regexp.properties[name]; !exists {
		if regexp.metadata != nil {
			regexp.metadata.order = append(regexp.metadata.order, name)
			regexp.metadata.descriptors[name] = &tsDescriptor{writable: true, enumerable: true, configurable: true}
		} else { regexp.order = append(regexp.order, name) }
	}
	regexp.properties[name] = item
	return item
}
func tsRegExpHas(regexp *tsRegExp, name string) bool {
	if _, ok := regexp.properties[name]; ok { return true }
	if !tsRegExpIntrinsics.modified.Load() { switch name {
	case "lastIndex", "source", "flags", "global", "hasIndices", "ignoreCase", "multiline", "dotAll", "unicode", "unicodeSets", "sticky", "toString", "exec", "test": return true
	case (*tsSymbol)(tsWellKnownSymbol("match").ref).key,
		(*tsSymbol)(tsWellKnownSymbol("matchAll").ref).key,
		(*tsSymbol)(tsWellKnownSymbol("replace").ref).key,
		(*tsSymbol)(tsWellKnownSymbol("search").ref).key,
		(*tsSymbol)(tsWellKnownSymbol("split").ref).key:
		return true
	}
	}
	if !tsIsUndefined(regexp.prototype) { return tsPrototypeHas(regexp.prototype, name) }
	return false
}

func tsRegexStringValue(loop *tsLoop, receiver tsValue, name string, args []tsValue) tsValue {
	if tsNullish(receiver) {
		tsPropertyFailure("String method called on null or undefined")
	}
	argument := tsArg(args, 0)
	symbolName := ""
	switch name {
	case "match": symbolName = "match"
	case "matchAll": symbolName = "matchAll"
	case "replace", "replaceAll": symbolName = "replace"
	case "search": symbolName = "search"
	case "split": symbolName = "split"
	default: panic("unknown String regular-expression method")
	}
	if name == "matchAll" && tsIsRegExpValue(loop, argument) && !tsRegExpFlagsContainGlobal(loop, argument) {
		tsPropertyFailure("String.prototype.matchAll requires a global RegExp")
	}
	if name == "replaceAll" && tsIsRegExpValue(loop, argument) && !tsRegExpFlagsContainGlobal(loop, argument) {
		tsPropertyFailure("String.prototype.replaceAll requires a global RegExp")
	}
	if !tsIsPrimitiveValue(argument) {
		method := tsGet(argument, tsWellKnownSymbol(symbolName))
		if !tsNullish(method) {
			if method.kind != tsFunctionKind { tsPropertyFailure("String regular-expression method is not callable") }
			callArgs := []tsValue{tsStringReference(tsStringValueWithLoop(loop, receiver))}
			if name == "replace" || name == "replaceAll" { callArgs = append(callArgs, tsArg(args, 1)) }
			return tsCallReceiver(loop, method, argument, callArgs...)
		}
	}
	text := tsStringValueWithLoop(loop, receiver)
	if name == "replace" || name == "replaceAll" {
		return tsStringReplaceLiteral(loop, text, tsStringValueWithLoop(loop, argument), tsArg(args, 1), name == "replaceAll")
	}
	if name == "match" || name == "matchAll" || name == "search" {
		if argument.kind == tsRegExpKind {
			flags := ""
			if name == "matchAll" { flags = "g" }
			regexp := tsDynamicRegExp(loop, argument, tsStringReference(tsStringUTF8(flags)))
			method := tsGet(loop, regexp, tsWellKnownSymbol(symbolName))
			return tsCallReceiver(loop, method, regexp, tsStringReference(text))
		}
		pattern := tsStringValueWithLoop(loop, argument)
		flags := ""
		if name == "matchAll" { flags = "g" }
		if regexp, ok := tsRegExpLiteralSearch(pattern, flags); ok {
			method := tsGet(loop, regexp, tsWellKnownSymbol(symbolName))
			return tsCallReceiver(loop, method, regexp, tsStringReference(text))
		}
	}
	tsPropertyFailure("Dynamic RegExp construction is not supported by Go code generation; use a regular-expression literal")
	return tsU
}

// Plain literal patterns are common inputs to String.prototype.match, search,
// and matchAll. They can be represented directly as a matcher without parsing
// or compiling dynamic RegExp syntax.
func tsRegExpLiteralSearch(pattern *tsString, flags string) (tsValue, bool) {
	for _, character := range pattern.units {
		switch character {
		case '.', '*', '+', '?', '^', '$', '{', '}', '(', ')', '|', '[', ']', '\\':
			return tsU, false
		}
	}
	matcher := func(input []uint16, start int, captures []int32) (int, int, bool) {
		if start < 0 { start = 0 }
		if start > len(input) || len(captures) < 2 { return 0, 0, false }
		for index := range captures { captures[index] = -1 }
		last := len(input) - len(pattern.units)
		for position := start; position <= last; position++ {
			matched := true
			for offset, character := range pattern.units {
				if input[position+offset] != character && (!strings.ContainsRune(flags, 'i') || !strings.EqualFold(string(rune(input[position+offset])), string(rune(character)))) { matched = false; break }
			}
			if matched {
				captures[0], captures[1] = int32(position), int32(position+len(pattern.units))
				return position, position + len(pattern.units), true
			}
		}
		return 0, 0, false
	}
	return tsNewCompiledRegExp(pattern, tsStringUTF8(flags), matcher, 2, nil), true
}

// String replacement with a non-RegExp search value is a literal substring
// operation after ToString(searchValue). It does not need dynamic RegExp
// compilation and must still honor replacement callbacks and dollar tokens.
func tsStringReplaceLiteral(loop *tsLoop, input, search *tsString, replacement tsValue, replaceAll bool) tsValue {
	functional := replacement.kind == tsFunctionKind
	var replacementString *tsString
	if !functional { replacementString = tsStringValueWithLoop(loop, replacement) }
	positions := make([]int, 0, 4)
	if len(search.units) == 0 {
		for position := 0; position <= len(input.units); position++ {
			positions = append(positions, position)
			if !replaceAll { break }
		}
	} else {
		for from := 0; from <= len(input.units)-len(search.units); {
			position := tsStringIndex(input, search, from)
			if position < 0 { break }
			positions = append(positions, position)
			if !replaceAll { break }
			from = position + len(search.units)
		}
	}
	if len(positions) == 0 { return tsStringReference(input) }
	parts := make([]*tsString, 0, len(positions)*2+1)
	next := 0
	for _, position := range positions {
		matched := input.slice(position, position+len(search.units))
		parts = append(parts, input.slice(next, position))
		if functional {
			result := tsCall(loop, replacement, tsStringReference(matched), tsNumberValue(float64(position)), tsStringReference(input))
			parts = append(parts, tsStringValueWithLoop(loop, result))
		} else {
			parts = append(parts, tsRegExpGetSubstitution(loop, replacementString, input, position, matched, nil, tsU))
		}
		next = position + len(search.units)
	}
	parts = append(parts, input.slice(next, len(input.units)))
	return tsStringReference(tsStringConcat(parts...))
}

func tsIsRegExpValue(loop *tsLoop, value tsValue) bool {
	if tsIsPrimitiveValue(value) { return false }
	matcher := tsGet(value, tsWellKnownSymbol("match"))
	if !tsIsUndefined(matcher) { return tsTruthy(matcher) }
	return value.kind == tsRegExpKind
}

func tsRegExpFlagsContainGlobal(loop *tsLoop, value tsValue) bool {
	flags := tsStringValueWithLoop(loop, tsGet(value, tsStringReference(tsStringUTF8("flags"))))
	for _, character := range flags.units { if character == 'g' { return true } }
	return false
}

func tsRegExpSymbolMethod(loop *tsLoop, receiver tsValue, name string, args ...tsValue) tsValue {
	input := tsStringValueWithLoop(loop, tsArg(args, 0))
	// Preserve the direct path for an ordinary literal. Customized instances
	// and generic receivers use the spec algorithm and observable properties.
	if receiver.kind == tsRegExpKind && tsRegExpIsUnmodifiedRegExp(receiver) {
		if name == "matchAll" { return tsRegExpMatchAll(loop, receiver, input) }
		return tsRegexStringRegExp(loop, input, name, (*tsRegExp)(receiver.ref), tsArg(args, 1))
	}
	if tsIsPrimitiveValue(receiver) { tsPropertyFailure("RegExp symbol method called on incompatible receiver") }
	if name == "matchAll" { return tsRegExpMatchAll(loop, receiver, input) }
	return tsRegExpGenericSymbolMethod(loop, receiver, input, name, tsArg(args, 1))
}

// tsRegExpGenericSymbolMethod implements the RegExp symbol algorithms for
// receivers whose properties or exec method may be user supplied. The common
// unmodified RegExp literal case stays on the direct path above.
func tsRegExpGenericSymbolMethod(loop *tsLoop, receiver tsValue, input *tsString, name string, replacement tsValue) tsValue {
	readRegExpProperty := func(key string) tsValue { return tsGet(loop, receiver, tsStringReference(tsStringUTF8(key))) }
	writeRegExpProperty := func(key string, value tsValue) { tsSet(loop, receiver, tsStringReference(tsStringUTF8(key)), value) }
	switch name {
	case "search":
		previous := readRegExpProperty("lastIndex")
		zero := tsNumberValue(0)
		if !tsSameValue(previous, zero) { writeRegExpProperty("lastIndex", zero) }
		result := tsRegExpExec(loop, receiver, input, tsU)
		if !tsSameValue(readRegExpProperty("lastIndex"), previous) { writeRegExpProperty("lastIndex", previous) }
		if tsNullish(result) { return tsNumberValue(-1) }
		return tsNumberValue(tsRegExpToIntegerOrInfinity(loop, tsGet(loop, result, tsStringReference(tsStringUTF8("index")))))
	case "match":
		// @@match obtains global and Unicode behavior from the observable
		// flags property. Reading global directly skips user-defined flags
		// accessors and changes the required exception order.
		flags := tsStringValueWithLoop(loop, readRegExpProperty("flags")).String()
		global := strings.ContainsRune(flags, 'g')
		unicodeMode := strings.ContainsRune(flags, 'u') || strings.ContainsRune(flags, 'v')
		if !global { return tsRegExpExec(loop, receiver, input, tsU) }
		writeRegExpProperty("lastIndex", tsNumberValue(0))
		out := &tsArray{}
		for {
			result := tsRegExpExec(loop, receiver, input, tsU)
			if tsNullish(result) { break }
			matched := tsStringValueWithLoop(loop, tsGet(loop, result, tsStringReference(tsStringUTF8("0"))))
			out.values = append(out.values, tsStringReference(matched))
			if len(matched.units) == 0 { tsRegExpSetAdvancedStringIndex(loop, receiver, input, unicodeMode) }
		}
		if len(out.values) == 0 { return tsNull }
		return tsArrayValue(out)
	case "replace", "replaceAll":
		functional := replacement.kind == tsFunctionKind
		var replacementString *tsString
		if !functional { replacementString = tsStringValueWithLoop(loop, replacement) }
		// The generic algorithm obtains its flags string after coercing a
		// non-callable replacement. The inherited flags accessor performs the
		// ordered reads of hasIndices, global, ignoreCase, multiline, dotAll,
		// unicode, unicodeSets, and sticky; custom flags accessors are observable
		// and must be honored for both replacement forms.
		flags := tsStringValueWithLoop(loop, readRegExpProperty("flags"))
		global, unicodeMode := false, false
		for _, flag := range flags.units {
			switch flag {
			case 'g':
				global = true
			case 'u', 'v':
				unicodeMode = true
			}
		}
		if name == "replaceAll" && !global { tsPropertyFailure("String.prototype.replaceAll requires a global RegExp") }
		if global { writeRegExpProperty("lastIndex", tsNumberValue(0)) }
		matches := make([]tsValue, 0, 4)
		for {
			result := tsRegExpExec(loop, receiver, input, tsU)
			if tsNullish(result) { break }
			matches = append(matches, result)
			if !global { break }
			matched := tsStringValueWithLoop(loop, tsGet(loop, result, tsStringReference(tsStringUTF8("0"))))
			if len(matched.units) == 0 { tsRegExpSetAdvancedStringIndex(loop, receiver, input, unicodeMode) }
		}
		if len(matches) == 0 { return tsStringReference(input) }
		parts := make([]*tsString, 0, len(matches)*2+1)
		nextSourcePosition := 0
		for _, result := range matches {
			captureLength := tsToLength(loop, tsGet(loop, result, tsStringReference(tsStringUTF8("length"))))
			if captureLength > 0 {
				captureLength--
			}
			matched := tsStringValueWithLoop(loop, tsGet(loop, result, tsStringReference(tsStringUTF8("0"))))
			positionNumber := tsRegExpToIntegerOrInfinity(loop, tsGet(loop, result, tsStringReference(tsStringUTF8("index"))))
			position := 0
			if !math.IsNaN(positionNumber) && positionNumber > 0 {
				if math.IsInf(positionNumber, 1) || positionNumber >= float64(len(input.units)) { position = len(input.units) } else { position = int(positionNumber) }
			}
			args, groups := tsRegExpReplacementArguments(loop, result, input, position, matched, captureLength)
			if position < nextSourcePosition { continue }
			if functional {
				callArgs := append([]tsValue(nil), args...)
				if !tsIsUndefined(groups) { callArgs = append(callArgs, groups) }
				parts = append(parts, input.slice(nextSourcePosition, position), tsStringValueWithLoop(loop, tsCall(loop, replacement, callArgs...)))
			} else {
				// GetSubstitution first applies ToObject to a present groups value,
				// even when the replacement text has no named-capture token. Null must
				// therefore throw here rather than being silently ignored.
				if groups.kind == tsNullKind { tsPropertyFailure("Cannot convert undefined or null to object") }
				captures := args[1:len(args)-2]
				parts = append(parts, input.slice(nextSourcePosition, position), tsRegExpGetSubstitution(loop, replacementString, input, position, matched, captures, groups))
			}
			nextSourcePosition = min(position+len(matched.units), len(input.units))
		}
		parts = append(parts, input.slice(nextSourcePosition, len(input.units)))
		return tsStringReference(tsStringConcat(parts...))
	case "split":
		return tsRegExpGenericSplit(loop, receiver, input, replacement)
	default:
		tsPropertyFailure("Unsupported RegExp symbol method")
	}
	return tsU
}

func tsRegExpGenericSplit(loop *tsLoop, receiver tsValue, input *tsString, limitValue tsValue) tsValue {
	out := &tsArray{}
	_ = tsIsRegExpValue(loop, receiver)
	constructor := tsGet(loop, receiver, tsStringReference(tsStringUTF8("constructor")))
	if !tsIsUndefined(constructor) && tsIsPrimitiveValue(constructor) { tsPropertyFailure("RegExp constructor property must be an object") }
	var species tsValue
	if !tsIsUndefined(constructor) { species = tsGet(loop, constructor, tsWellKnownSymbol("species")) }
	if tsNullish(species) { species = tsRegExpConstructor() }
	if !tsIsConstructor(species) { tsPropertyFailure("RegExp species is not a constructor") }
	flags := tsStringValueWithLoop(loop, tsGet(loop, receiver, tsStringReference(tsStringUTF8("flags"))))
	newFlags := append([]uint16(nil), flags.units...)
	sticky := false
	for _, flag := range flags.units { if flag == 'y' { sticky = true; break } }
	if !sticky { newFlags = append(newFlags, 'y') }
	splitter := tsConstruct(species, receiver, tsStringReference(tsStringUnits(newFlags)))
	if tsIsPrimitiveValue(splitter) { tsPropertyFailure("RegExp species constructor must return an object") }
	limit := uint64(0xffffffff)
	if !tsIsUndefined(limitValue) { limit = uint64(tsToUint32(loop, limitValue)) }
	if limit == 0 { return tsArrayValue(out) }
	size := len(input.units)
	if size == 0 {
		result := tsRegExpExec(loop, splitter, input, tsU)
		if tsNullish(result) { out.values = append(out.values, tsStringReference(input)) }
		return tsArrayValue(out)
	}
	unicodeMode := false
	for _, flag := range flags.units { if flag == 'u' || flag == 'v' { unicodeMode = true } }
	p, q := 0, 0
	for q < size && uint64(len(out.values)) < limit {
		tsSet(loop, splitter, tsStringReference(tsStringUTF8("lastIndex")), tsNumberValue(float64(q)))
		result := tsRegExpExec(loop, splitter, input, tsU)
		if tsNullish(result) {
			q = int(tsRegExpAdvanceStringIndex(input, uint64(q), unicodeMode))
			continue
		}
		e := int(tsToLength(loop, tsGet(loop, splitter, tsStringReference(tsStringUTF8("lastIndex")))))
		if e > size { e = size }
		if e == p {
			q = int(tsRegExpAdvanceStringIndex(input, uint64(q), unicodeMode))
			continue
		}
		out.values = append(out.values, tsStringReference(input.slice(p, q)))
		length := tsToLength(loop, tsGet(loop, result, tsStringReference(tsStringUTF8("length"))))
		for capture := 1; capture < length && uint64(len(out.values)) < limit; capture++ {
			out.values = append(out.values, tsGet(loop, result, tsStringReference(tsStringUTF8(strconv.Itoa(capture)))))
		}
		p, q = e, e
	}
	if uint64(len(out.values)) < limit { out.values = append(out.values, tsStringReference(input.slice(p, size))) }
	return tsArrayValue(out)
}

var tsRegExpStringIteratorIntrinsics struct {
	sync.Once
	prototype *tsObject
}

func tsRegExpStringIteratorPrototype() tsValue {
	tsRegExpStringIteratorIntrinsics.Do(func() {
		tsInitGenerators()
		prototype := tsNewObject()
		prototype.prototype = tsObjectValue(tsGeneratorIntrinsics.iteratorPrototype)
		prototype.set("next", tsNamedNativeMethod("next", 0, func(loop *tsLoop, receiver tsValue, _ ...tsValue) tsValue {
			if receiver.kind != tsIteratorKind || !(*tsIterator)(receiver.ref).regexpString {
				tsPropertyFailure("RegExp String Iterator.prototype.next called on incompatible receiver")
			}
			iterator := (*tsIterator)(receiver.ref)
			if iterator.done { return tsIteratorResult(tsU, true) }
			if !iterator.next() { iterator.done = true; return tsIteratorResult(tsU, true) }
			return tsIteratorResult(iterator.value, false)
		}))
		prototype.descriptors = map[string]*tsDescriptor{
			"next": {writable: true, configurable: true},
		}
		tag := tsPropertyKey(tsWellKnownSymbol("toStringTag"))
		prototype.set(tag, tsStringReference(tsStringUTF8("RegExp String Iterator")))
		prototype.descriptors[tag] = &tsDescriptor{configurable: true}
		tsRegExpStringIteratorIntrinsics.prototype = prototype
	})
	return tsObjectValue(tsRegExpStringIteratorIntrinsics.prototype)
}

func tsRegExpMatchAll(loop *tsLoop, receiver tsValue, input *tsString) tsValue {
	// MatchAllIterator performs IsRegExp before reading flags or constructing
	// the species clone. Keep the observable @@match lookup on customized and
	// generic receivers.
	_ = tsIsRegExpValue(loop, receiver)
	flags := tsStringValueWithLoop(loop, tsGet(loop, receiver, tsStringReference(tsStringUTF8("flags"))))
	constructor := tsGet(loop, receiver, tsStringReference(tsStringUTF8("constructor")))
	var species tsValue
	if !tsIsUndefined(constructor) {
		if tsIsPrimitiveValue(constructor) {
			tsPropertyFailure("RegExp constructor property must be an object")
		}
		species = tsGet(loop, constructor, tsWellKnownSymbol("species"))
	}
	matcher := tsU
	if tsNullish(species) {
		matcher = tsConstruct(tsRegExpConstructor(), receiver, tsStringReference(flags))
	} else {
		if !tsIsConstructor(species) {
			tsPropertyFailure("RegExp species is not a constructor")
		}
		matcher = tsConstruct(species, receiver, tsStringReference(flags))
		if tsIsPrimitiveValue(matcher) {
			tsPropertyFailure("RegExp species constructor must return an object")
		}
	}
	lastIndex := tsToLength(loop, tsRegExpGetLastIndex(loop, receiver))
	if matcher.kind == tsRegExpKind {
		tsRegExpSetLastIndex(loop, matcher, tsNumberValue(float64(lastIndex)))
	} else {
		tsSet(matcher, tsStringReference(tsStringUTF8("lastIndex")), tsNumberValue(float64(lastIndex)))
	}
	global, fullUnicode := false, false
	for _, flag := range flags.units {
		switch flag {
		case 'g': global = true
		case 'u', 'v': fullUnicode = true
		}
	}
	done := false
	return tsIteratorValueReference(&tsIterator{prototype: tsRegExpStringIteratorPrototype(), regexpString: true, pull: func() (tsValue, bool) {
		if done {
			return tsU, false
		}
		result := tsRegExpExec(loop, matcher, input, tsU)
		if tsNullish(result) {
			done = true
			return result, false
		}
		if !global {
			done = true
		} else {
			matched := tsStringValueWithLoop(loop, tsGet(loop, result, tsStringReference(tsStringUTF8("0"))))
			if len(matched.units) == 0 {
				tsRegExpSetAdvancedStringIndex(loop, matcher, input, fullUnicode)
			}
		}
		return result, true
	}})
}

func tsRegExpReplacementArguments(loop *tsLoop, result tsValue, input *tsString, position int, matched *tsString, captureCount int) ([]tsValue, tsValue) {
	args := make([]tsValue, 0, captureCount+4)
	args = append(args, tsStringReference(matched))
	for capture := 1; capture <= captureCount; capture++ {
		value := tsGet(result, tsStringReference(tsStringUTF8(strconv.Itoa(capture))))
		args = append(args, value)
	}
	args = append(args, tsNumberValue(float64(position)), tsStringReference(input))
	groups := tsGet(result, tsStringReference(tsStringUTF8("groups")))
	return args, groups
}

func tsRegExpGetSubstitution(loop *tsLoop, replacement, input *tsString, position int, matched *tsString, captures []tsValue, groups tsValue) *tsString {
	parts := make([]*tsString, 0, 8)
	appendCapture := func(index int) bool {
		if index < 1 || index > len(captures) { return false }
		if !tsIsUndefined(captures[index-1]) { parts = append(parts, tsStringValueWithLoop(loop, captures[index-1])) }
		return true
	}
	for index := 0; index < len(replacement.units); {
		if replacement.units[index] != '$' || index+1 >= len(replacement.units) {
			start := index
			index++
			for index < len(replacement.units) && replacement.units[index] != '$' { index++ }
			parts = append(parts, replacement.slice(start, index))
			continue
		}
		next := replacement.units[index+1]
		switch next {
		case '$': parts = append(parts, tsStringUTF8("$")); index += 2; continue
		case '&': parts = append(parts, matched); index += 2; continue
		case 0x60: if position > 0 { parts = append(parts, input.slice(0, position)) }; index += 2; continue
		case '\'': end := min(position+len(matched.units), len(input.units)); if end < len(input.units) { parts = append(parts, input.slice(end, len(input.units))) }; index += 2; continue
		case '<':
			if !tsIsUndefined(groups) {
				nameEnd := index + 2
				for nameEnd < len(replacement.units) && replacement.units[nameEnd] != '>' { nameEnd++ }
				if nameEnd < len(replacement.units) {
					name := tsStringReference(replacement.slice(index+2, nameEnd))
					value := tsGet(groups, name)
					if !tsIsUndefined(value) { parts = append(parts, tsStringValueWithLoop(loop, value)) }
					index = nameEnd + 1
					continue
				}
			}
		}
		if next >= '0' && next <= '9' {
			one := int(next - '0')
			if one == 0 {
				if index+2 < len(replacement.units) && replacement.units[index+2] >= '1' && replacement.units[index+2] <= '9' {
					candidate := int(replacement.units[index+2] - '0')
					if candidate <= len(captures) && appendCapture(candidate) { index += 3; continue }
				}
			} else {
				two := one
				consumed := 2
				if index+2 < len(replacement.units) && replacement.units[index+2] >= '0' && replacement.units[index+2] <= '9' {
					candidate := one*10 + int(replacement.units[index+2]-'0')
					if candidate <= len(captures) { two = candidate; consumed = 3 }
				}
				if appendCapture(two) { index += consumed; continue }
				if consumed == 3 && appendCapture(one) { index += 2; continue }
			}
		}
		parts = append(parts, replacement.slice(index, index+1))
		index++
	}
	return tsStringConcat(parts...)
}

func tsRegexStringRegExp(loop *tsLoop, text *tsString, name string, regexp *tsRegExp, replacement tsValue) tsValue {
	receiver := tsRegExpValue(regexp)
	switch name {
	case "search":
		previous := tsRegExpGetLastIndex(loop, receiver)
		if !tsSameValue(previous, tsNumberValue(0)) { tsRegExpSetLastIndex(loop, receiver, tsNumberValue(0)) }
		result := tsRegExpExec(loop, receiver, text, tsU)
		if !tsSameValue(tsRegExpGetLastIndex(loop, receiver), previous) { tsRegExpSetLastIndex(loop, receiver, previous) }
		if tsNullish(result) { return tsNumberValue(-1) }
		return tsNumberValue(tsRegExpToIntegerOrInfinity(loop, tsGet(result, tsStringReference(tsStringUTF8("index")))))
	case "match":
		if !regexp.global { return tsRegExpExec(loop, receiver, text, tsU) }
		tsRegExpSetLastIndex(loop, receiver, tsNumberValue(0))
		if tsRegExpHasInitialMap(receiver) { return tsRegExpMatchGlobalFast(loop, receiver, regexp, text) }
		out := &tsArray{}
		for {
			result := tsRegExpExec(loop, receiver, text, tsU)
			if tsNullish(result) { break }
			fullMatch := tsGet(result, tsStringReference(tsStringUTF8("0")))
			out.values = append(out.values, fullMatch)
			if len(tsStringValue(fullMatch).units) == 0 { tsRegExpSetAdvancedStringIndex(loop, receiver, text, regexp.unicode || regexp.unicodeSets) }
		}
		if len(out.values) == 0 { return tsNull }
		return tsArrayValue(out)
	case "matchAll":
		if !regexp.global { tsPropertyFailure("String.prototype.matchAll requires a global RegExp") }
		clone := *regexp
		clone.lastIndex = tsToLength(loop, tsRegExpGetLastIndex(loop, receiver))
		clone.lastIndexValue = tsU
		clone.customLastIndex = false
		clone.properties = nil
		clone.order = nil
		clone.metadata = nil
		cloneValue := tsRegExpValue(&clone)
		return tsIteratorValueReference(&tsIterator{pull: func() (tsValue, bool) {
			result := tsRegExpExec(loop, cloneValue, text, tsU)
			if tsNullish(result) { return result, false }
			fullMatch := tsGet(result, tsStringReference(tsStringUTF8("0")))
			if len(tsStringValue(fullMatch).units) == 0 { tsRegExpSetAdvancedStringIndex(loop, cloneValue, text, clone.unicode || clone.unicodeSets) }
			return result, true
		}})
	case "replace", "replaceAll":
		global := regexp.global || name == "replaceAll"
		if name == "replaceAll" && !regexp.global { tsPropertyFailure("String.prototype.replaceAll requires a global RegExp") }
		functional := replacement.kind == tsFunctionKind
		var replacementString *tsString
		if !functional { replacementString = tsStringValueWithLoop(loop, replacement) }
		if global { tsRegExpSetLastIndex(loop, receiver, tsNumberValue(0)) }
		matches := make([]tsValue, 0, 4)
		for {
			result := tsRegExpExec(loop, receiver, text, tsU)
			if tsNullish(result) { break }
			matches = append(matches, result)
			if !global { break }
			matched := tsStringValueWithLoop(loop, tsGet(result, tsStringReference(tsStringUTF8("0"))))
			if len(matched.units) == 0 { tsRegExpSetAdvancedStringIndex(loop, receiver, text, regexp.unicode || regexp.unicodeSets) }
		}
		if len(matches) == 0 { return tsStringReference(text) }
		parts := make([]*tsString, 0, len(matches)*2+1)
		nextSourcePosition := 0
		for _, result := range matches {
			captureLength := tsToLength(loop, tsGet(result, tsStringReference(tsStringUTF8("length"))))
			if captureLength > 0 {
				captureLength--
			}
			matched := tsStringValueWithLoop(loop, tsGet(result, tsStringReference(tsStringUTF8("0"))))
			positionValue := tsGet(result, tsStringReference(tsStringUTF8("index")))
			positionNumber := tsNumber(tsToPrimitive(positionValue))
			position := 0
			if !math.IsNaN(positionNumber) && positionNumber > 0 {
				if math.IsInf(positionNumber, 1) || positionNumber >= float64(len(text.units)) { position = len(text.units) } else { position = int(math.Trunc(positionNumber)) }
			}
			if position < nextSourcePosition { continue }
			if position > len(text.units) { position = len(text.units) }
			args, groups := tsRegExpReplacementArguments(loop, result, text, position, matched, captureLength)
			if functional {
				callArgs := append([]tsValue(nil), args...)
				if !tsIsUndefined(groups) { callArgs = append(callArgs, groups) }
				parts = append(parts, text.slice(nextSourcePosition, position), tsStringValueWithLoop(loop, tsCall(loop, replacement, callArgs...)))
			} else {
				captureValues := args[1 : len(args)-2]
				substitution := tsRegExpGetSubstitution(loop, replacementString, text, position, matched, captureValues, groups)
				parts = append(parts, text.slice(nextSourcePosition, position), substitution)
			}
			nextSourcePosition = min(position+len(matched.units), len(text.units))
		}
		parts = append(parts, text.slice(nextSourcePosition, len(text.units)))
		return tsStringReference(tsStringConcat(parts...))
	case "split":
		// Converting an object limit can mutate RegExp.prototype or the
		// receiver itself. Use the observable species algorithm so splitter
		// construction happens before the conversion, as required.
		if !tsIsPrimitiveValue(replacement) {
			return tsRegExpGenericSplit(loop, receiver, text, replacement)
		}
		out := &tsArray{}
		limit := uint32(0xffffffff)
		if !tsIsUndefined(replacement) { limit = tsToUint32(replacement) }
		if limit == 0 { return tsArrayValue(out) }
		splitRegexp := *regexp
		splitRegexp.sticky = true
		splitRegexp.lastIndex = 0
		splitRegexp.customLastIndex = false
		splitRegexp.lastIndexValue = tsU
		if len(text.units) == 0 {
			scope := &tsRegExpResultVectorScope{}
			captures := scope.Initialize(loop, splitRegexp.captureSlots)
			_, _, matched := tsRegExpFindInto(&splitRegexp, text, 0, captures)
			scope.Close()
			if !matched { out.values = append(out.values, tsStringReference(text)) }
			return tsArrayValue(out)
		}
		scope := &tsRegExpResultVectorScope{}
		captures := scope.Initialize(loop, splitRegexp.captureSlots)
		defer scope.Close()
		p, q := 0, 0
		for q < len(text.units) && uint32(len(out.values)) < limit {
			matchStart, matchEnd, matched := tsRegExpFindInto(&splitRegexp, text, q, captures)
			if !matched || matchStart != q {
				q = int(tsRegExpAdvanceStringIndex(text, uint64(q), regexp.unicode || regexp.unicodeSets))
				continue
			}
			tsRegExpSetLastMatchInfo(loop, text, len(captures), captures)
			if matchEnd == p {
				q = int(tsRegExpAdvanceStringIndex(text, uint64(q), regexp.unicode || regexp.unicodeSets))
				continue
			}
			out.values = append(out.values, tsStringReference(text.slice(p, q)))
			for capture := 1; capture < regexp.captureSlots/2 && uint32(len(out.values)) < limit; capture++ {
				if !tsRegExpIsMatchedCapture(&tsRegExpMatchInfo{lastSubject: text, lastInput: text, captures: captures}, capture) { out.values = append(out.values, tsU); continue }
				startIndex, endIndex := int(captures[tsRegExpCaptureStartIndex(capture)]), int(captures[tsRegExpCaptureEndIndex(capture)])
				out.values = append(out.values, tsStringReference(text.slice(startIndex, endIndex)))
			}
			p = matchEnd
			q = p
		}
		if uint32(len(out.values)) < limit { out.values = append(out.values, tsStringReference(text.slice(p, len(text.units)))) }
		return tsArrayValue(out)
	}
	return tsU
}
`
