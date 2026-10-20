package goemit

// nodeHTTPCommonRuntime adapts Node's _http_common helpers and the HTTPParser
// binding boundary with a compact incremental HTTP/1 parser.
const nodeHTTPCommonRuntime = `
type tsNodeHTTPParserState struct {
	buffer []byte
	current []byte
	started time.Time
	lastMessageStart time.Time
	maxHeaderSize uint64
	connections tsValue
	parserType int
	lenientFlags uint32
	chunkState uint8
	chunkRemaining uint64
	chunkTrailers []tsValue
	paused, closed, initialized, headerDone, messageBegun, consumed, skipBody, upgrade, closeDelimited, closedByMessage, inExecute bool
	bodyStart, bodyDelivered int
	bodyReceived int64
}

func tsNodeHTTPCommonWordByte(value byte) bool {
	return value >= 'a' && value <= 'z' || value >= 'A' && value <= 'Z' || value >= '0' && value <= '9' || value == '_'
}

func tsNodeHTTPCommonHasBoundaryToken(value, token string) bool {
	lowerValue := strings.ToLower(value)
	lowerToken := strings.ToLower(token)
	for offset := 0; offset <= len(lowerValue)-len(lowerToken); {
		index := strings.Index(lowerValue[offset:], lowerToken)
		if index < 0 {
			return false
		}
		start := offset + index
		end := start + len(lowerToken)
		beforeWord := start > 0 && tsNodeHTTPCommonWordByte(lowerValue[start-1])
		afterWord := end < len(lowerValue) && tsNodeHTTPCommonWordByte(lowerValue[end])
		if !beforeWord && !afterWord {
			return true
		}
		offset = start + 1
	}
	return false
}

type tsNodeHTTPConnectionListState struct { all, active []tsValue }

var tsNodeHTTPCommonSingleton struct {
	sync.Once
	module, constructor, parsers, connectionsConstructor tsValue
	prototype, connectionsPrototype *tsObject
	states map[uintptr]*tsNodeHTTPParserState
	connections map[uintptr]*tsNodeHTTPConnectionListState
	kIncomingMessage, kSkipPendingData tsValue
	methods, allMethods []string
	warnedLenient bool
}

func tsNodeHTTPCommonKey(name string) tsValue { return tsStringReference(tsStringUTF8(name)) }
func tsNodeHTTPCommonGet(loop *tsLoop, object tsValue, name string) tsValue { return tsGet(loop, object, tsNodeHTTPCommonKey(name)) }
func tsNodeHTTPCommonSet(loop *tsLoop, object tsValue, name string, value tsValue) { tsSet(loop, object, tsNodeHTTPCommonKey(name), value) }
func tsNodeHTTPCommonStateFor(parser tsValue) *tsNodeHTTPParserState {
	if parser.kind != tsObjectKind { return nil }
	return tsNodeHTTPCommonSingleton.states[uintptr(parser.ref)]
}
func tsNodeHTTPCommonIndex(index int) tsValue { return tsNodeHTTPCommonKey(strconv.Itoa(index)) }
func tsNodeHTTPCommonCall(loop *tsLoop, receiver tsValue, name string, args ...tsValue) tsValue {
	method := tsNodeHTTPCommonGet(loop, receiver, name)
	if method.kind == tsFunctionKind { return tsCallReceiver(loop, method, receiver, args...) }
	return tsU
}
func tsNodeHTTPCommonCallback(loop *tsLoop, parser tsValue, index int, args ...tsValue) tsValue {
	callback := tsGet(loop, parser, tsNodeHTTPCommonIndex(index))
	if callback.kind == tsFunctionKind { return tsCallReceiver(loop, callback, parser, args...) }
	return tsU
}
func tsNodeHTTPCommonHeaderError(name, code, message string) tsValue {
	errorValue := tsFSException("TypeError", code, message)
	return errorValue
}
func tsNodeHTTPCommonInvalidHeaderValue(loop *tsLoop, value tsValue, name tsValue) bool {
	text := tsText(tsStringValueWithLoop(loop, value))
	for _, char := range text {
		if char == 0 || char == '\n' || char == '\r' || char > 0xff { return true }
	}
	return false
}
func tsNodeHTTPCommonStrictHeaderValue(loop *tsLoop, value tsValue) bool {
	text := tsText(tsStringValueWithLoop(loop, value))
	for _, char := range text {
		if char == '\t' || char >= 0x20 && char <= 0x7e || char >= 0x80 && char <= 0xff { continue }
		return true
	}
	return false
}
func tsNodeHTTPCommonValidateHeaderName(loop *tsLoop, name, label tsValue) tsValue {
	validator := tsInternalValidatorsModule(loop)
	check := tsNodeHTTPCommonGet(loop, validator, "checkIsHttpToken")
	if name.kind != tsStringKind || tsText(name) == "" || !tsTruthy(tsCall(loop, check, name)) {
		if !tsTruthy(label) { label = tsNodeHTTPCommonKey("Header name") }
		message := tsText(label) + " must be a valid HTTP token [" + tsText(tsStringValueWithLoop(loop, name)) + "]"
		panic(tsThrown{tsNodeHTTPCommonHeaderError("Header name", "ERR_INVALID_HTTP_TOKEN", message)})
	}
	return tsU
}
func tsNodeHTTPCommonValidateHeaderValue(loop *tsLoop, name, value tsValue) tsValue {
	if tsIsUndefined(value) {
		message := "Invalid value \"undefined\" for header \"" + tsText(name) + "\""
		panic(tsThrown{tsNodeHTTPCommonHeaderError("Header value", "ERR_HTTP_INVALID_HEADER_VALUE", message)})
	}
	if tsNodeHTTPCommonStrictHeaderValue(loop, value) {
		message := "Invalid character in header content [" + tsText(name) + "]"
		panic(tsThrown{tsNodeHTTPCommonHeaderError("Header value", "ERR_INVALID_CHAR", message)})
	}
	return tsU
}
func tsNodeHTTPCommonParseError(loop *tsLoop, parser tsValue, bytesParsed int, code, reason string) tsValue {
	errorValue := tsFSException("Error", code, "Parse Error: "+reason)
	tsSet(loop, errorValue, tsNodeHTTPCommonKey("bytesParsed"), tsNumberValue(float64(bytesParsed)))
	tsSet(loop, errorValue, tsNodeHTTPCommonKey("reason"), tsNodeHTTPCommonKey(reason))
	tsSet(loop, errorValue, tsNodeHTTPCommonKey("rawPacket"), tsNodeHTTPCommonCall(loop, parser, "getCurrentBuffer"))
	return errorValue
}
func tsNodeHTTPCommonParserBufferValue(data []byte) tsValue { return tsNodeBuffer(append([]byte(nil), data...)) }
func tsNodeHTTPCommonSetHeaders(loop *tsLoop, parser tsValue, headers []tsValue) {
	tsNodeHTTPCommonSet(loop, parser, "_headers", tsArrayValue(&tsArray{values: headers}))
}
func tsNodeHTTPCommonHTTPToken(value string) bool {
	if value == "" { return false }
	for _, ch := range []byte(value) {
		if ch >= '0' && ch <= '9' || ch >= 'A' && ch <= 'Z' || ch >= 'a' && ch <= 'z' { continue }
		switch ch { case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '\x60', '|', '~': continue }
		return false
	}
	return true
}
func tsNodeHTTPCommonHeaderValueInvalid(value string, relaxed bool) bool {
	for _, ch := range []byte(value) {
		if relaxed {
			if ch == 0 || ch == '\r' || ch == '\n' { return true }
			continue
		}
		if ch == '\t' || ch >= 0x20 && ch <= 0x7e || ch >= 0x80 { continue }
		return true
	}
	return false
}
func tsNodeHTTPCommonHeaderHasToken(value, token string) bool {
	for _, item := range strings.Split(value, ",") { if strings.EqualFold(strings.TrimSpace(item), token) { return true } }
	return false
}
func tsNodeHTTPCommonTransferChunked(value string) (bool, bool, bool) {
	if strings.TrimSpace(value) == "" { return false, false, false }
	items := strings.Split(strings.ToLower(value), ",")
	hasChunked := false
	duplicateChunked := false
	for _, item := range items {
		if strings.TrimSpace(item) == "chunked" {
			if hasChunked { duplicateChunked = true }
			hasChunked = true
		}
	}
	return hasChunked, strings.TrimSpace(items[len(items)-1]) == "chunked", duplicateChunked
}
func tsNodeHTTPCommonParseContentLength(value string) (int64, bool) {
	value = strings.Trim(value, " \t")
	if value == "" { return 0, false }
	for _, digit := range []byte(value) { if digit < '0' || digit > '9' { return 0, false } }
	length, err := strconv.ParseInt(value, 10, 64)
	if err != nil { return 0, false }
	return length, true
}
func tsNodeHTTPCommonProtocolVersion(value string) (string, int, int, bool) {
	protocol := ""
	for _, candidate := range []string{"HTTP/", "RTSP/", "ICE/"} {
		if strings.HasPrefix(value, candidate) { protocol = candidate[:len(candidate)-1]; value = value[len(candidate):]; break }
	}
	if protocol == "" || len(value) != 3 || value[0] < '0' || value[0] > '9' || value[1] != '.' || value[2] < '0' || value[2] > '9' { return "", 0, 0, false }
	return protocol, int(value[0]-'0'), int(value[2]-'0'), true
}
func tsNodeHTTPCommonProtocolVersionSupported(major, minor int) bool {
	return major == 0 && minor == 9 || major == 1 && (minor == 0 || minor == 1) || major == 2 && minor == 0
}
func tsNodeHTTPCommonHeaderPairs(lines []string, max int, lenientFlags uint32) ([]tsValue, string) {
	pairs := make([]tsValue, 0, len(lines)*2)
	for _, line := range lines {
		if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {
			if lenientFlags&1 == 0 || len(pairs) < 2 { return nil, line }
			value := strings.TrimLeft(line, " \t")
			if tsNodeHTTPCommonHeaderValueInvalid(value, lenientFlags&1024 != 0) { return nil, line }
			pairs[len(pairs)-1] = tsNodeHTTPCommonKey(tsText(pairs[len(pairs)-1]) + " " + strings.TrimRight(value, " \t"))
			continue
		}
		colon := strings.IndexByte(line, ':')
		if colon <= 0 { return nil, line }
		field := line[:colon]; if lenientFlags&1 != 0 { field = strings.TrimRight(field, " \t") }
		if !tsNodeHTTPCommonHTTPToken(field) { return nil, line }
		if max > 0 && len(pairs)/2 >= max { break }
	value := strings.TrimRight(strings.TrimLeft(line[colon+1:], " \t"), " \t")
		if tsNodeHTTPCommonHeaderValueInvalid(value, lenientFlags&1024 != 0) { return nil, line }
		pairs = append(pairs, tsNodeHTTPCommonKey(field), tsNodeHTTPCommonKey(value))
	}
	return pairs, ""
}
func tsNodeHTTPCommonParserHeadersCallback(loop *tsLoop, parser tsValue, headers []tsValue, url tsValue) tsValue {
	max := int(tsNumber(tsNodeHTTPCommonGet(loop, parser, "maxHeaderPairs")))
	current := tsNodeHTTPCommonArray(tsNodeHTTPCommonGet(loop, parser, "_headers"))
	capacity := max - len(current.values)
	if max <= 0 || capacity >= len(headers) { current.values = append(current.values, headers...) } else if capacity > 0 { current.values = append(current.values, headers[:capacity]...) }
	tsNodeHTTPCommonSet(loop, parser, "_headers", tsArrayValue(current))
	old := tsNodeHTTPCommonGet(loop, parser, "_url")
	tsNodeHTTPCommonSet(loop, parser, "_url", tsNodeHTTPCommonKey(tsText(old)+tsText(url)))
	return tsU
}
func tsNodeHTTPCommonArray(value tsValue) *tsArray {
	if value.kind == tsArrayKind { return (*tsArray)(value.ref) }
	return &tsArray{}
}
func tsNodeHTTPCommonConnectionListState(connections tsValue) *tsNodeHTTPConnectionListState {
	if connections.kind != tsObjectKind { return nil }
	return tsNodeHTTPCommonSingleton.connections[uintptr(connections.ref)]
}
func tsNodeHTTPCommonConnectionListPush(_ *tsLoop, connections, parser tsValue, active bool) {
	state := tsNodeHTTPCommonConnectionListState(connections)
	if state == nil || parser.kind != tsObjectKind { return }
	items := &state.all; if active { items = &state.active }
	for _, existing := range *items { if tsStrictEqual(existing, parser) { return } }
	*items = append(*items, parser)
}
func tsNodeHTTPCommonConnectionListPop(_ *tsLoop, connections, parser tsValue, active bool) {
	state := tsNodeHTTPCommonConnectionListState(connections)
	if state == nil { return }
	items := &state.all; if active { items = &state.active }
	for index, existing := range *items {
		if tsStrictEqual(existing, parser) { copy((*items)[index:], (*items)[index+1:]); *items = (*items)[:len(*items)-1]; return }
	}
}
func tsNodeHTTPCommonConnectionListValues(_ *tsLoop, receiver tsValue, active bool) []tsValue {
	state := tsNodeHTTPCommonConnectionListState(receiver); if state == nil { return nil }
	items := state.all; if active { items = state.active }
	return append([]tsValue(nil), items...)
}
func tsNodeHTTPCommonConnectionListCreate(loop *tsLoop) tsValue {
	object := tsNewObject()
	object.prototype = tsObjectValue(tsNodeHTTPCommonSingleton.connectionsPrototype)
	tsNodeHTTPCommonSingleton.connections[uintptr(unsafe.Pointer(object))] = &tsNodeHTTPConnectionListState{}
	return tsObjectValue(object)
}
func tsNodeHTTPCommonConnectionListMethod(loop *tsLoop, receiver tsValue, name string, args []tsValue) tsValue {
	if tsNodeHTTPCommonConnectionListState(receiver) == nil { return tsU }
	switch name {
	case "all": return tsArrayValue(&tsArray{values:tsNodeHTTPCommonConnectionListValues(loop, receiver, false)})
	case "active": return tsArrayValue(&tsArray{values:tsNodeHTTPCommonConnectionListValues(loop, receiver, true)})
	case "idle":
		idle := &tsArray{}
		for _, parser := range tsNodeHTTPCommonConnectionListValues(loop, receiver, false) { if state := tsNodeHTTPCommonStateFor(parser); state != nil && state.lastMessageStart.IsZero() { idle.values = append(idle.values, parser) } }
		return tsArrayValue(idle)
	case "expired":
		headersTimeout := uint64(tsToUint32(tsArg(args, 0))) * 1_000_000
		requestTimeout := uint64(tsToUint32(tsArg(args, 1))) * 1_000_000
		if headersTimeout == 0 && requestTimeout == 0 { return tsArrayValue(&tsArray{}) }
		if requestTimeout > 0 && headersTimeout > requestTimeout { headersTimeout, requestTimeout = requestTimeout, headersTimeout }
		uptime := time.Since(tsNodeProcessStarted)
		now := time.Now()
		var headersDeadline, requestDeadline time.Time
		if headersTimeout > 0 && uptime > time.Duration(headersTimeout) { headersDeadline = now.Add(-time.Duration(headersTimeout)) }
		if requestTimeout > 0 && uptime > time.Duration(requestTimeout) { requestDeadline = now.Add(-time.Duration(requestTimeout)) }
		expired := &tsArray{}
		for _, parser := range tsNodeHTTPCommonConnectionListValues(loop, receiver, true) {
			state := tsNodeHTTPCommonStateFor(parser)
			if state == nil || state.lastMessageStart.IsZero() { continue }
			tooOldForHeaders := !state.headerDone && !headersDeadline.IsZero() && state.lastMessageStart.Before(headersDeadline)
			tooOldForRequest := !requestDeadline.IsZero() && state.lastMessageStart.Before(requestDeadline)
			if tooOldForHeaders || tooOldForRequest { expired.values = append(expired.values, parser); tsNodeHTTPCommonConnectionListPop(loop, receiver, parser, true) }
		}
		return tsArrayValue(expired)
	}
	return tsU
}
func tsNodeHTTPCommonHeadersComplete(loop *tsLoop, parser tsValue, args []tsValue) tsValue {
	major, minor := tsArg(args, 0), tsArg(args, 1)
	headers, method, rawURL := tsArg(args, 2), tsArg(args, 3), tsArg(args, 4)
	status, statusMessage, upgrade, keepAlive := tsArg(args, 5), tsArg(args, 6), tsArg(args, 7), tsArg(args, 8)
	socket := tsNodeHTTPCommonGet(loop, parser, "socket")
	if tsIsUndefined(headers) { headers = tsNodeHTTPCommonGet(loop, parser, "_headers"); tsNodeHTTPCommonSet(loop, parser, "_headers", tsArrayValue(&tsArray{})) }
	if tsIsUndefined(rawURL) { rawURL = tsNodeHTTPCommonGet(loop, parser, "_url"); tsNodeHTTPCommonSet(loop, parser, "_url", tsNodeHTTPCommonKey("")) }
	constructor := tsGet(loop, tsNodeHTTPIncomingModule(loop), tsNodeHTTPCommonKey("IncomingMessage"))
	server := tsNodeHTTPCommonGet(loop, socket, "server")
	customConstructor := tsGet(loop, server, tsNodeHTTPCommonSingleton.kIncomingMessage)
	if !tsNullish(customConstructor) { constructor = customConstructor }
	incoming := tsConstruct(constructor, socket)
	tsNodeHTTPCommonSet(loop, parser, "incoming", incoming)
	tsNodeHTTPCommonSet(loop, incoming, "httpVersionMajor", major)
	tsNodeHTTPCommonSet(loop, incoming, "httpVersionMinor", minor)
	tsNodeHTTPCommonSet(loop, incoming, "httpVersion", tsNodeHTTPCommonKey(tsText(major)+"."+tsText(minor)))
	join := tsNodeHTTPCommonGet(loop, server, "joinDuplicateHeaders")
	if !tsTruthy(join) { join = tsNodeHTTPCommonGet(loop, parser, "joinDuplicateHeaders") }
	tsNodeHTTPCommonSet(loop, incoming, "joinDuplicateHeaders", join)
	tsNodeHTTPCommonSet(loop, incoming, "url", rawURL)
	tsNodeHTTPCommonSet(loop, incoming, "upgrade", upgrade)
	count := len(tsNodeHTTPCommonArray(headers).values)
	max := int(tsNumber(tsNodeHTTPCommonGet(loop, parser, "maxHeaderPairs")))
	if max > 0 && count > max { count = max }
	tsNodeHTTPCommonCall(loop, incoming, "_addHeaderLines", headers, tsNumberValue(float64(count)))
	if tsIsNumeric(method) {
		index := int(tsNumber(method)); if index >= 0 && index < len(tsNodeHTTPCommonSingleton.allMethods) { tsNodeHTTPCommonSet(loop, incoming, "method", tsNodeHTTPCommonKey(tsNodeHTTPCommonSingleton.allMethods[index])) }
	} else {
		tsNodeHTTPCommonSet(loop, incoming, "statusCode", status)
		tsNodeHTTPCommonSet(loop, incoming, "statusMessage", statusMessage)
	}
	return tsNodeHTTPCommonCall(loop, parser, "onIncoming", incoming, keepAlive)
}
func tsNodeHTTPCommonParserOnBody(loop *tsLoop, parser, body tsValue) tsValue {
	incoming := tsNodeHTTPCommonGet(loop, parser, "incoming")
	if tsNullish(incoming) || tsTruthy(tsGet(loop, incoming, tsNodeHTTPCommonSingleton.kSkipPendingData)) { return tsU }
	if !tsTruthy(tsNodeHTTPCommonGet(loop, incoming, "_dumped")) {
		ret := tsNodeHTTPCommonCall(loop, incoming, "push", body)
		if !tsTruthy(ret) { tsNodeHTTPIncomingReadStop(loop, tsNodeHTTPCommonGet(loop, parser, "socket")) }
	}
	return tsU
}
func tsNodeHTTPCommonParserOnMessageComplete(loop *tsLoop, parser tsValue) tsValue {
	state := tsNodeHTTPCommonStateFor(parser)
	if state != nil {
		// Node moves completed parsers out of active tracking, clears their
		// deadline, and keeps them in the all-connections set while idle.
		tsNodeHTTPCommonConnectionListPop(loop, state.connections, parser, false)
		tsNodeHTTPCommonConnectionListPop(loop, state.connections, parser, true)
		state.lastMessageStart = time.Time{}
		state.messageBegun = false
		tsNodeHTTPCommonConnectionListPush(loop, state.connections, parser, false)
	}
	incoming := tsNodeHTTPCommonGet(loop, parser, "incoming")
	if !tsNullish(incoming) && !tsTruthy(tsGet(loop, incoming, tsNodeHTTPCommonSingleton.kSkipPendingData)) {
		tsNodeHTTPCommonSet(loop, incoming, "complete", tsBooleanValue(true))
		headers := tsNodeHTTPCommonGet(loop, parser, "_headers")
		if tsNodeHTTPCommonArray(headers).length() > 0 {
			tsNodeHTTPCommonCall(loop, incoming, "_addHeaderLines", headers, tsNumberValue(float64(tsNodeHTTPCommonArray(headers).length())))
			tsNodeHTTPCommonSet(loop, parser, "_headers", tsArrayValue(&tsArray{}))
			tsNodeHTTPCommonSet(loop, parser, "_url", tsNodeHTTPCommonKey(""))
		}
		tsNodeHTTPCommonCall(loop, incoming, "push", tsNull)
	}
	tsNodeHTTPIncomingReadStart(loop, tsNodeHTTPCommonGet(loop, parser, "socket"))
	return tsU
}
func tsNodeHTTPCommonCloseParser(parser tsValue) tsValue { if state := tsNodeHTTPCommonStateFor(parser); state != nil { state.closed = true }; return tsU }
func tsNodeHTTPCommonParserMessageBegin(loop *tsLoop, parser tsValue, state *tsNodeHTTPParserState) {
	if state == nil || state.messageBegun { return }
	// llhttp starts a fresh deadline for each message on a persistent socket.
	// Reinsert to retain ordering by the parser's current request-start time.
	tsNodeHTTPCommonConnectionListPop(loop, state.connections, parser, false)
	tsNodeHTTPCommonConnectionListPop(loop, state.connections, parser, true)
	state.lastMessageStart = time.Now()
	state.messageBegun = true
	tsNodeHTTPCommonConnectionListPush(loop, state.connections, parser, false)
	tsNodeHTTPCommonConnectionListPush(loop, state.connections, parser, true)
	if callback := tsGet(loop, parser, tsNodeHTTPCommonIndex(0)); callback.kind == tsFunctionKind { tsCallReceiver(loop, callback, parser) }
}
func tsNodeHTTPCommonCleanParser(loop *tsLoop, parser tsValue) {
	tsNodeHTTPCommonSet(loop, parser, "_headers", tsArrayValue(&tsArray{}))
	tsNodeHTTPCommonSet(loop, parser, "_url", tsNodeHTTPCommonKey(""))
	for _, name := range []string{"socket", "incoming", "outgoing", "onIncoming", "joinDuplicateHeaders"} { tsNodeHTTPCommonSet(loop, parser, name, tsNull) }
	tsNodeHTTPCommonSet(loop, parser, "maxHeaderPairs", tsNumberValue(2000))
	tsNodeHTTPCommonSet(loop, parser, "_consumed", tsBooleanValue(false))
	for _, index := range []int{0, 5, 6} { tsSet(loop, parser, tsNodeHTTPCommonIndex(index), tsNull) }
	if state := tsNodeHTTPCommonStateFor(parser); state != nil { state.buffer = nil; state.current = nil; state.paused = false; state.initialized = false; state.headerDone = false; state.messageBegun = false; state.consumed = false; state.skipBody = false; state.upgrade = false; state.closeDelimited = false; state.closedByMessage = false; state.bodyStart = 0; state.bodyDelivered = 0; state.bodyReceived = 0; state.chunkState = 0; state.chunkRemaining = 0; state.chunkTrailers = nil; state.started = time.Time{}; state.lastMessageStart = time.Time{}; state.connections = tsU; state.maxHeaderSize = 0; state.lenientFlags = 0 }
}
func tsNodeHTTPCommonChunkLine(data []byte, offset int, lenientFlags uint32) (int, int, bool, bool) {
	for index := offset; index < len(data); index++ {
		switch data[index] {
		case '\r':
			if index+1 < len(data) && data[index+1] == '\n' { return index, index + 2, true, false }
			if lenientFlags&64 != 0 { return index, index + 1, true, false }
			if index+1 == len(data) { return 0, 0, false, false }
			return 0, 0, false, true
		case '\n':
			if lenientFlags&256 != 0 { return index, index + 1, true, false }
			return 0, 0, false, true
		}
	}
	return 0, 0, false, false
}

func tsNodeHTTPCommonChunkExtensionsValid(extension []byte) bool {
	for offset := 0; offset < len(extension); {
		nameStart := offset
		for offset < len(extension) && tsNodeHTTPCommonHTTPToken(string(extension[offset:offset+1])) { offset++ }
		if offset == nameStart { return false }
		if offset < len(extension) && extension[offset] == '=' {
			offset++
			if offset == len(extension) { return false }
			if extension[offset] == '"' {
				offset++
				closed := false
				for offset < len(extension) {
					ch := extension[offset]
					if ch == '"' { offset++; closed = true; break }
					if ch == '\\' {
						offset++
						if offset == len(extension) { return false }
						ch = extension[offset]
						if ch != '\t' && (ch < 0x20 || ch == 0x7f) { return false }
					} else if ch != '\t' && (ch < 0x20 || ch == 0x7f) { return false }
					offset++
				}
				if !closed { return false }
			} else {
				valueStart := offset
				for offset < len(extension) && extension[offset] != ';' && extension[offset] != '"' {
					ch := extension[offset]
					if ch < 0x20 || ch > 0x7e { return false }
					offset++
				}
				if offset == valueStart { return false }
			}
		}
		if offset == len(extension) { return true }
		if extension[offset] != ';' { return false }
		offset++
		if offset == len(extension) { return false }
	}
	return true
}

func tsNodeHTTPCommonChunkedStep(data []byte, state *tsNodeHTTPParserState, maxTrailerPairs int) ([]byte, int, []tsValue, bool, string) {
	var body []byte
	offset := 0
	for {
		switch state.chunkState {
		case 0:
			lineEnd, lineNext, complete, invalidLineEnding := tsNodeHTTPCommonChunkLine(data, offset, state.lenientFlags)
			if !complete {
				if invalidLineEnding { return body, offset, nil, false, "Invalid chunk line ending" }
				if len(data)-offset > 16_384 { return body, offset, nil, false, "HPE_CHUNK_EXTENSIONS_OVERFLOW" }
				return body, offset, nil, false, ""
			}
			rawSizeLine := data[offset:lineEnd]
			sizeEnd := len(rawSizeLine)
			extensionStart := -1
			if semi := bytes.IndexByte(rawSizeLine, ';'); semi >= 0 { sizeEnd = semi; extensionStart = semi + 1 }
			sizeBytes := rawSizeLine[:sizeEnd]
			digitsEnd := 0
			for digitsEnd < len(sizeBytes) && (sizeBytes[digitsEnd] >= '0' && sizeBytes[digitsEnd] <= '9' || sizeBytes[digitsEnd] >= 'a' && sizeBytes[digitsEnd] <= 'f' || sizeBytes[digitsEnd] >= 'A' && sizeBytes[digitsEnd] <= 'F') { digitsEnd++ }
			if digitsEnd == 0 { return body, offset, nil, false, "Invalid chunk size" }
			for _, ch := range sizeBytes[digitsEnd:] {
				if state.lenientFlags&512 == 0 || ch != ' ' && ch != '\t' { return body, offset, nil, false, "Invalid chunk size" }
			}
			if extensionStart >= 0 {
				extension := rawSizeLine[extensionStart:]
				if len(extension) > 16_384 { return body, offset, nil, false, "HPE_CHUNK_EXTENSIONS_OVERFLOW" }
				if !tsNodeHTTPCommonChunkExtensionsValid(extension) { return body, offset, nil, false, "Invalid chunk extension" }
			}
			size, err := strconv.ParseUint(string(sizeBytes[:digitsEnd]), 16, 64)
			if err != nil { return body, offset, nil, false, "Invalid chunk size" }
			offset = lineNext
			if size == 0 { state.chunkState = 3; continue }
			state.chunkRemaining = size
			state.chunkState = 1
		case 1:
			available := len(data) - offset
			if available == 0 { return body, offset, nil, false, "" }
			count := uint64(available)
			if count > state.chunkRemaining { count = state.chunkRemaining }
			end := offset + int(count)
			body = append(body, data[offset:end]...)
			offset = end
			state.chunkRemaining -= count
			if state.chunkRemaining > 0 { return body, offset, nil, false, "" }
			state.chunkState = 2
		case 2:
			if offset == len(data) { return body, offset, nil, false, "" }
			switch data[offset] {
			case '\r':
				if offset+1 == len(data) {
					if state.lenientFlags&128 == 0 { return body, offset, nil, false, "" }
					offset++
				} else if data[offset+1] == '\n' { offset += 2
				} else if state.lenientFlags&128 != 0 { offset++
				} else { return body, offset, nil, false, "Invalid chunk terminator" }
			case '\n':
				if state.lenientFlags&256 == 0 { return body, offset, nil, false, "Invalid chunk terminator" }
				offset++
			default:
				if state.lenientFlags&128 == 0 { return body, offset, nil, false, "Invalid chunk terminator" }
			}
			state.chunkState = 0
		case 3:
			lineEnd, lineNext, complete, invalidLineEnding := tsNodeHTTPCommonChunkLine(data, offset, state.lenientFlags)
			if !complete {
				if invalidLineEnding { return body, offset, nil, false, "Invalid chunk line ending" }
				if len(data)-offset > 16_384 { return body, offset, nil, false, "Invalid chunk trailer" }
				return body, offset, nil, false, ""
			}
			if lineEnd == offset {
				offset = lineNext
				trailers := append([]tsValue(nil), state.chunkTrailers...)
				state.chunkTrailers = nil
				state.chunkState = 0
				return body, offset, trailers, true, ""
			}
			pairs, badLine := tsNodeHTTPCommonHeaderPairs([]string{string(data[offset:lineEnd])}, 0, state.lenientFlags)
			if badLine != "" || len(pairs) != 2 { return body, offset, nil, false, "Invalid chunk trailer" }
			if maxTrailerPairs <= 0 || len(state.chunkTrailers)/2 < maxTrailerPairs {
				state.chunkTrailers = append(state.chunkTrailers, pairs...)
			}
			offset = lineNext
		default:
			return body, offset, nil, false, "Invalid chunk size"
		}
	}
}

func tsNodeHTTPCommonParserExecute(loop *tsLoop, parser tsValue, data tsValue) tsValue {
	state := tsNodeHTTPCommonStateFor(parser)
	if state == nil { panic(tsThrown{tsFSException("TypeError", "ERR_INVALID_THIS", "HTTPParser.execute called on an incompatible receiver")}) }
	if state.closed { return tsU }
	if data.kind != tsTypedArrayKind && data.kind != tsDataViewKind { return tsU }
	if !state.initialized { panic(tsThrown{tsFSException("Error", "ERR_INVALID_STATE", "HTTPParser is not initialized")}) }
	if state.inExecute { panic(tsThrown{tsFSException("TypeError", "", "HTTPParser.execute is not reentrant")}) }
	state.inExecute = true
	defer func() { state.inExecute = false }()
	chunk := tsFSBytes(data)
	bufferedBefore := len(state.buffer)
	state.current = append(state.current[:0], chunk...); state.buffer = append(state.buffer, chunk...)
	parsedFromCurrent := func(total int) int {
		parsed := total - bufferedBefore
		if parsed < 0 { return 0 }
		if parsed > len(chunk) { return len(chunk) }
		return parsed
	}
	if state.paused { return tsNodeHTTPCommonParseError(loop, parser, 0, "HPE_PAUSED", "Paused") }
	totalParsed := 0
parseNextMessage:
	if state.closedByMessage {
		offset := 0
		for offset < len(state.buffer) && (state.buffer[offset] == '\r' || state.buffer[offset] == '\n') { offset++ }
		if offset < len(state.buffer) {
			if state.lenientFlags&32 == 0 { return tsNodeHTTPCommonParseError(loop, parser, offset, "HPE_CLOSED_CONNECTION", "Data after \x60Connection: close\x60") }
			state.buffer = nil
			return tsNumberValue(float64(len(chunk)))
		}
		state.buffer = nil
		return tsNumberValue(float64(parsedFromCurrent(totalParsed)))
	}
	if len(state.buffer) > 0 { tsNodeHTTPCommonParserMessageBegin(loop, parser, state); if state.paused { return tsNodeHTTPCommonParseError(loop, parser, 0, "HPE_PAUSED", "Paused") } }
	boundary := bytes.Index(state.buffer, []byte("\r\n\r\n"))
	http09Request := false
	if boundary < 0 && state.parserType == 1 {
		// HTTP/0.9 has no header section: its request line is followed directly
		// by the end of the request. Recognize that form only after the CRLF is
		// complete, so a fragmented request line remains buffered as usual.
		requestLineEnd := bytes.Index(state.buffer, []byte("\r\n"))
		if requestLineEnd >= 0 && len(strings.Fields(string(state.buffer[:requestLineEnd]))) == 2 {
			boundary = requestLineEnd
			http09Request = true
		}
	}
	if boundary < 0 {
		if state.maxHeaderSize > 0 && uint64(len(state.buffer)) > state.maxHeaderSize { return tsNodeHTTPCommonParseError(loop, parser, 0, "HPE_HEADER_OVERFLOW", "Header overflow") }
		return tsNumberValue(float64(len(chunk)))
	}
	headerEnd := boundary + 4
	if http09Request { headerEnd = boundary + 2 }
	if state.maxHeaderSize > 0 && uint64(headerEnd) > state.maxHeaderSize { return tsNodeHTTPCommonParseError(loop, parser, boundary, "HPE_HEADER_OVERFLOW", "Header overflow") }
	alreadyParsedHeaders := state.headerDone
	state.headerDone = true
	lines := strings.Split(string(state.buffer[:boundary]), "\r\n")
	if len(lines) == 0 || lines[0] == "" { state.headerDone = false; return tsNodeHTTPCommonParseError(loop, parser, 0, "HPE_INVALID_CONSTANT", "Invalid HTTP message") }
	major, minor := 1, 1; methodValue, urlValue, status, statusMessage := tsU, tsNodeHTTPCommonKey(""), tsU, tsU
	parts := strings.Fields(lines[0])
	if len(parts) < 2 { state.headerDone = false; return tsNodeHTTPCommonParseError(loop, parser, 0, "HPE_INVALID_CONSTANT", "Invalid start line") }
	protocolStart := strings.HasPrefix(parts[0], "HTTP/") || strings.HasPrefix(parts[0], "RTSP/") || strings.HasPrefix(parts[0], "ICE/")
	if protocolStart && state.parserType == 1 { state.headerDone = false; return tsNodeHTTPCommonParseError(loop, parser, 0, "HPE_INVALID_CONSTANT", "Expected HTTP method") }
	if protocolStart {
		_, parsedMajor, parsedMinor, ok := tsNodeHTTPCommonProtocolVersion(parts[0])
		major, minor = parsedMajor, parsedMinor
		if !ok { state.headerDone = false; return tsNodeHTTPCommonParseError(loop, parser, 0, "HPE_INVALID_VERSION", "Invalid HTTP version") }
		if !tsNodeHTTPCommonProtocolVersionSupported(major, minor) && state.lenientFlags&16 == 0 { state.headerDone = false; return tsNodeHTTPCommonParseError(loop, parser, 0, "HPE_INVALID_VERSION", "Invalid HTTP version") }
		line := lines[0]
		versionEnd := strings.IndexByte(line, ' ')
		if versionEnd < 0 || versionEnd+1 >= len(line) || line[versionEnd+1] == ' ' { state.headerDone = false; return tsNodeHTTPCommonParseError(loop, parser, 0, "HPE_INVALID_STATUS", "Invalid status code") }
		rest := line[versionEnd+1:]
		statusEnd := strings.IndexByte(rest, ' ')
		statusToken := rest
		if statusEnd >= 0 { statusToken = rest[:statusEnd] }
		if len(statusToken) != 3 { state.headerDone = false; return tsNodeHTTPCommonParseError(loop, parser, 0, "HPE_INVALID_STATUS", "Invalid status code") }
		for _, digit := range []byte(statusToken) { if digit < '0' || digit > '9' { state.headerDone = false; return tsNodeHTTPCommonParseError(loop, parser, 0, "HPE_INVALID_STATUS", "Invalid status code") } }
		code, _ := strconv.Atoi(statusToken); status = tsNumberValue(float64(code))
		statusMessage = tsNodeHTTPCommonKey("")
		if statusEnd >= 0 { statusMessage = tsNodeHTTPCommonKey(rest[statusEnd+1:]) }
	} else {
		if strings.ContainsAny(lines[0], "\t\v\f") || strings.HasPrefix(lines[0], " ") || strings.HasSuffix(lines[0], " ") || strings.Contains(lines[0], "  ") { state.headerDone = false; return tsNodeHTTPCommonParseError(loop, parser, 0, "HPE_INVALID_CONSTANT", "Invalid start line") }
		if len(parts) > 3 { state.headerDone = false; return tsNodeHTTPCommonParseError(loop, parser, 0, "HPE_INVALID_CONSTANT", "Invalid start line") }
		methodValue = tsU
		for index, method := range tsNodeHTTPCommonSingleton.allMethods { if method == parts[0] { methodValue = tsNumberValue(float64(index)); break } }
		if tsIsUndefined(methodValue) { state.headerDone = false; return tsNodeHTTPCommonParseError(loop, parser, 0, "HPE_INVALID_METHOD", "Invalid method") }
		for _, ch := range []byte(parts[1]) { if ch <= 0x20 || ch == 0x7f { state.headerDone = false; return tsNodeHTTPCommonParseError(loop, parser, 0, "HPE_INVALID_URL", "Invalid URL") } }
		urlValue = tsNodeHTTPCommonKey(parts[1])
		protocol := "HTTP"
		if len(parts) == 2 { major, minor = 0, 9 } else {
			var ok bool
			protocol, major, minor, ok = tsNodeHTTPCommonProtocolVersion(parts[2])
			if !ok || !tsNodeHTTPCommonProtocolVersionSupported(major, minor) && state.lenientFlags&16 == 0 { state.headerDone = false; return tsNodeHTTPCommonParseError(loop, parser, 0, "HPE_INVALID_VERSION", "Invalid HTTP version") }
		}
		methodName := tsNodeHTTPCommonSingleton.allMethods[int(tsNumber(methodValue))]
		if major == 0 && minor == 9 && methodName != "GET" { state.headerDone = false; return tsNodeHTTPCommonParseError(loop, parser, 0, "HPE_INVALID_METHOD", "HTTP/0.9 only supports GET") }
		if protocol == "ICE" && methodName != "SOURCE" { state.headerDone = false; return tsNodeHTTPCommonParseError(loop, parser, 0, "HPE_INVALID_CONSTANT", "Expected SOURCE method for ICE/x.x request") }
		if protocol == "HTTP" {
			switch methodName { case "DESCRIBE", "ANNOUNCE", "SETUP", "PLAY", "PAUSE", "TEARDOWN", "GET_PARAMETER", "SET_PARAMETER", "REDIRECT", "RECORD", "FLUSH": state.headerDone = false; return tsNodeHTTPCommonParseError(loop, parser, 0, "HPE_INVALID_METHOD", "Invalid method for HTTP/x.x request") }
		}
		if protocol == "RTSP" {
			switch methodName { case "GET", "POST", "OPTIONS", "DESCRIBE", "ANNOUNCE", "SETUP", "PLAY", "PAUSE", "TEARDOWN", "GET_PARAMETER", "SET_PARAMETER", "REDIRECT", "RECORD", "FLUSH": default: state.headerDone = false; return tsNodeHTTPCommonParseError(loop, parser, 0, "HPE_INVALID_CONSTANT", "Invalid method for RTSP/x.x request") }
		}
	}
	pairs, badLine := tsNodeHTTPCommonHeaderPairs(lines[1:], 0, state.lenientFlags)
	if badLine != "" { state.headerDone = false; return tsNodeHTTPCommonParseError(loop, parser, 0, "HPE_INVALID_HEADER_TOKEN", "Invalid header") }
	if state.parserType == 1 && major == 0 && minor == 9 && len(pairs) != 0 { state.headerDone = false; return tsNodeHTTPCommonParseError(loop, parser, boundary, "HPE_INVALID_HEADER_TOKEN", "HTTP/0.9 requests cannot include headers") }
	if state.parserType == 1 {
		maxPairs := tsNumber(tsNodeHTTPCommonGet(loop, parser, "maxHeaderPairs"))
		if maxPairs > 0 && float64(len(pairs)) > maxPairs { state.headerDone = false; return tsNodeHTTPCommonParseError(loop, parser, boundary, "HPE_HEADER_OVERFLOW", "Header overflow") }
	}
	flat := tsArrayValue(&tsArray{values:pairs})
	connection, transfer, lengthValue, upgradeValue := "", "", "", ""
	contentLengthSeen, transferEncodingSeen := false, false
	for index := 0; index+1 < len(pairs); index += 2 {
		key := strings.ToLower(tsText(pairs[index])); value := tsText(pairs[index+1])
		switch key {
		case "connection": if connection != "" { connection += "," }; connection += strings.ToLower(value)
		case "transfer-encoding": transferEncodingSeen = true; if transfer != "" { transfer += "," }; transfer += strings.ToLower(value)
		case "content-length":
			if contentLengthSeen { state.headerDone = false; return tsNodeHTTPCommonParseError(loop, parser, boundary, "HPE_UNEXPECTED_CONTENT_LENGTH", "Duplicate Content-Length") }
			contentLengthSeen = true
			lengthValue = value
		case "upgrade": upgradeValue = value
		}
	}
	hasChunked, chunkedLast, duplicateChunked := tsNodeHTTPCommonTransferChunked(transfer)
	if transferEncodingSeen && strings.TrimSpace(transfer) == "" { state.headerDone = false; return tsNodeHTTPCommonParseError(loop, parser, boundary, "HPE_INVALID_TRANSFER_ENCODING", "Invalid \x60Transfer-Encoding\x60 header") }
	if duplicateChunked && state.lenientFlags&8 == 0 { state.headerDone = false; return tsNodeHTTPCommonParseError(loop, parser, boundary, "HPE_INVALID_TRANSFER_ENCODING", "Chunked encoding must not be applied more than once") }
	if state.parserType == 1 && transferEncodingSeen && !chunkedLast && state.lenientFlags&8 == 0 { state.headerDone = false; return tsNodeHTTPCommonParseError(loop, parser, boundary, "HPE_INVALID_TRANSFER_ENCODING", "Request has invalid \x60Transfer-Encoding\x60") }
	if transferEncodingSeen && contentLengthSeen && state.lenientFlags&2 == 0 { state.headerDone = false; return tsNodeHTTPCommonParseError(loop, parser, boundary, "HPE_UNEXPECTED_CONTENT_LENGTH", "Content-Length can't be present with Transfer-Encoding") }
	keepAlive := major > 0 && minor > 0
	if tsNodeHTTPCommonHeaderHasToken(connection, "close") { keepAlive = false } else if tsNodeHTTPCommonHeaderHasToken(connection, "keep-alive") { keepAlive = true }
	noBodyStatus := status.kind == tsNumberKind && (int(tsNumber(status)) < 200 || int(tsNumber(status)) == 204 || int(tsNumber(status)) == 304)
	needsEOF := state.parserType == 2 && !noBodyStatus && (transferEncodingSeen && !chunkedLast || !transferEncodingSeen && !contentLengthSeen)
	if needsEOF { keepAlive = false }
	upgrade := tsNodeHTTPCommonHeaderHasToken(connection, "upgrade") && upgradeValue != "" && (state.parserType == 1 || int(tsNumber(status)) == 101)
	if state.parserType == 1 && methodValue.kind == tsNumberKind && tsNodeHTTPCommonSingleton.allMethods[int(tsNumber(methodValue))] == "CONNECT" { upgrade = true }
	shouldKeepAlive := tsBooleanValue(keepAlive)
	if !alreadyParsedHeaders {
		tsNodeHTTPCommonParserMessageBegin(loop, parser, state)
		headersValue, callbackURL := flat, urlValue
		if headerCount := len(pairs)/2; headerCount > 31 {
			// The native binding may flush its fixed-size header buffer more than
			// once. Keep the callback sequence and original field order intact.
			for start := 0; start < len(pairs); {
				end := start + 31*2; if end > len(pairs) { end = len(pairs) }
				urlChunk := callbackURL; if start > 0 { urlChunk = tsNodeHTTPCommonKey("") }
				tsNodeHTTPCommonCallback(loop, parser, 1, tsArrayValue(&tsArray{values:append([]tsValue(nil), pairs[start:end]...)}), urlChunk)
				start = end
			}
			headersValue, callbackURL = tsU, tsU
		}
		completeArgs := []tsValue{tsNumberValue(float64(major)), tsNumberValue(float64(minor)), headersValue, methodValue, callbackURL, status, statusMessage, upgrade, shouldKeepAlive}
		callbackResult := tsNodeHTTPCommonCallback(loop, parser, 2, completeArgs...)
		state.skipBody = tsTruthy(callbackResult)
		state.upgrade = tsNumber(callbackResult) == 2
		if state.paused { return tsNodeHTTPCommonParseError(loop, parser, boundary, "HPE_PAUSED", "Paused") }
	}
	bodyStart := headerEnd
	if !alreadyParsedHeaders {
		state.bodyStart = bodyStart
		state.bodyDelivered = 0
		state.bodyReceived = 0
		state.chunkState = 0
		state.chunkRemaining = 0
		state.chunkTrailers = nil
		noBodyStatus := status.kind == tsNumberKind && (int(tsNumber(status)) < 200 || int(tsNumber(status)) == 204 || int(tsNumber(status)) == 304)
		state.closeDelimited = state.parserType == 2 && !contentLengthSeen && !chunkedLast && !noBodyStatus
	}
	consumedBody := 0; var body []byte; var trailers []tsValue; var chunkError string; bodyComplete := false
	if state.upgrade {
		// Node hands accepted upgrades and CONNECT tunnels to the application at
		// the end of the headers. Bytes after that boundary belong to the new
		// protocol even when the HTTP request advertised a body length.
		bodyComplete = true
	} else if chunkedLast {
		maxTrailerPairs := int(tsNumber(tsNodeHTTPCommonGet(loop, parser, "maxHeaderPairs")))
		body, consumedBody, trailers, bodyComplete, chunkError = tsNodeHTTPCommonChunkedStep(state.buffer[bodyStart:], state, maxTrailerPairs)
	} else if contentLengthSeen {
		length, validLength := tsNodeHTTPCommonParseContentLength(lengthValue); if !validLength { state.headerDone = false; return tsNodeHTTPCommonParseError(loop, parser, boundary, "HPE_INVALID_CONTENT_LENGTH", "Invalid content length") }
		available := len(state.buffer)-bodyStart
		remaining := length - state.bodyReceived
		if remaining < 0 { state.headerDone = false; return tsNodeHTTPCommonParseError(loop, parser, boundary, "HPE_INVALID_CONTENT_LENGTH", "Invalid content length") }
		bodyLength := available; if int64(bodyLength) > remaining { bodyLength = int(remaining) }
		if bodyLength > 0 { body = append([]byte(nil), state.buffer[bodyStart:bodyStart+bodyLength]...) }
		consumedBody = bodyLength
		state.bodyReceived += int64(bodyLength)
		bodyComplete = int64(available) >= remaining
	} else if state.closeDelimited {
		if bodyStart < len(state.buffer) { body = append([]byte(nil), state.buffer[bodyStart:]...) }
	} else if state.parserType == 1 || status.kind == tsNumberKind && (int(tsNumber(status)) < 200 || int(tsNumber(status)) == 204 || int(tsNumber(status)) == 304) {
		bodyComplete = true
	}
		if !state.skipBody && len(body) > state.bodyDelivered {
			part := body[state.bodyDelivered:]
			tsNodeHTTPCommonCallback(loop, parser, 3, tsNodeHTTPCommonParserBufferValue(part), tsNumberValue(0))
			state.bodyDelivered = len(body)
			if state.paused { return tsNodeHTTPCommonParseError(loop, parser, boundary, "HPE_PAUSED", "Paused") }
		}
		if chunkError != "" {
			state.headerDone = false
			code, reason := "HPE_INVALID_CHUNK_SIZE", chunkError
			switch chunkError {
			case "Invalid chunk trailer": code = "HPE_INVALID_HEADER_TOKEN"
			case "Invalid chunk extension", "Invalid chunk terminator", "Invalid chunk line ending": code = "HPE_STRICT"
			case "HPE_CHUNK_EXTENSIONS_OVERFLOW": code, reason = chunkError, "Chunk extensions overflow"
			}
			return tsNodeHTTPCommonParseError(loop, parser, boundary, code, reason)
		}
	if !bodyComplete {
		if chunkedLast || contentLengthSeen && !state.closeDelimited {
			// Once body bytes have been delivered, keep the request headers and
			// unread bytes only. Chunk framing state and parsed trailers live in
			// the parser state, so consumed wire bytes can be discarded here too.
			remaining := bodyStart + consumedBody
			if remaining > len(state.buffer) { remaining = len(state.buffer) }
			buffer := append([]byte(nil), state.buffer[:bodyStart]...)
			state.buffer = append(buffer, state.buffer[remaining:]...)
			state.bodyDelivered = 0
		}
		return tsNumberValue(float64(len(chunk)))
	}
	if len(trailers) > 0 { tsNodeHTTPCommonSet(loop, parser, "_headers", tsArrayValue(&tsArray{values:trailers})) }
	messageComplete := tsGet(loop, parser, tsNodeHTTPCommonIndex(4)); if messageComplete.kind == tsFunctionKind { tsCallReceiver(loop, messageComplete, parser) }
	if state.paused { return tsNodeHTTPCommonParseError(loop, parser, boundary, "HPE_PAUSED", "Paused") }
	parsed := bodyStart + consumedBody
	if hasChunked { parsed = bodyStart + consumedBody }
	if parsed > len(state.buffer) { parsed = len(state.buffer) }
	upgraded := state.upgrade
	state.buffer = append([]byte(nil), state.buffer[parsed:]...); state.headerDone = false; state.skipBody = false; state.upgrade = false; state.closeDelimited = false; state.bodyStart = 0; state.bodyDelivered = 0; state.bodyReceived = 0; state.chunkState = 0; state.chunkRemaining = 0; state.chunkTrailers = nil
	if execute := tsGet(loop, parser, tsNodeHTTPCommonIndex(5)); execute.kind == tsFunctionKind { tsCallReceiver(loop, execute, parser, tsNumberValue(float64(parsedFromCurrent(parsed)))) }
	totalParsed += parsed
	if upgraded { return tsNumberValue(float64(parsedFromCurrent(totalParsed))) }
	if !keepAlive { state.closedByMessage = true }
	if len(state.buffer) == 0 { return tsNumberValue(float64(parsedFromCurrent(totalParsed))) }
	goto parseNextMessage
}
func tsNodeHTTPCommonParserCreate(loop *tsLoop) tsValue {
	parser := tsNewObject(); parser.prototype = tsObjectValue(tsNodeHTTPCommonSingleton.prototype)
	state := &tsNodeHTTPParserState{}; tsNodeHTTPCommonSingleton.states[uintptr(unsafe.Pointer(parser))] = state
	tsNodeHTTPCommonCleanParser(loop, tsObjectValue(parser))
	tsNodeHTTPCommonSet(loop, tsObjectValue(parser), "_type", tsNumberValue(0))
	return tsObjectValue(parser)
}
func tsNodeHTTPCommonParserInit(loop *tsLoop) {
	tsNodeHTTPCommonSingleton.Do(func() {
		tsNodeHTTPCommonSingleton.states = map[uintptr]*tsNodeHTTPParserState{}
		tsNodeHTTPCommonSingleton.connections = map[uintptr]*tsNodeHTTPConnectionListState{}
		tsNodeHTTPCommonSingleton.kIncomingMessage = tsNewSymbol(tsStringUTF8("IncomingMessage"))
		tsNodeHTTPCommonSingleton.kSkipPendingData = tsNewSymbol(tsStringUTF8("SkipPendingData"))
		methods := tsNodeHTTPCommonArray(tsInternalHTTPGet(loop, tsInternalHTTPModule(loop), "METHODS"))
		for _, value := range methods.values { tsNodeHTTPCommonSingleton.methods = append(tsNodeHTTPCommonSingleton.methods, tsText(value)) }
		for _, name := range strings.Fields("DELETE GET HEAD POST PUT CONNECT OPTIONS TRACE COPY LOCK MKCOL MOVE PROPFIND PROPPATCH SEARCH UNLOCK BIND REBIND UNBIND ACL REPORT MKACTIVITY CHECKOUT MERGE M-SEARCH NOTIFY SUBSCRIBE UNSUBSCRIBE PATCH PURGE MKCALENDAR LINK UNLINK SOURCE PRI DESCRIBE ANNOUNCE SETUP PLAY PAUSE TEARDOWN GET_PARAMETER SET_PARAMETER REDIRECT RECORD FLUSH QUERY") { tsNodeHTTPCommonSingleton.allMethods = append(tsNodeHTTPCommonSingleton.allMethods, name) }
		prototype := tsNewObject(); tsNodeHTTPCommonSingleton.prototype = prototype
		for name, length := range map[string]int{"close":0,"free":0,"remove":0,"finish":0,"pause":0,"resume":0,"unconsume":0,"getCurrentBuffer":0,"duration":0,"headersCompleted":0,"consume":0,"execute":0,"initialize":0} {
			methodName := name
			prototype.set(name, tsNamedNativeMethod(name, length, func(loop *tsLoop, receiver tsValue, args ...tsValue) tsValue { return tsNodeHTTPCommonParserMethod(loop, receiver, methodName, args) }))
		}
		tag := tsPropertyKey(tsWellKnownSymbol("toStringTag")); prototype.set(tag, tsNodeHTTPCommonKey("HTTPParser")); prototype.descriptors = map[string]*tsDescriptor{tag:{configurable:true}}
		constructor := tsFunc(func(_ *tsLoop, _ ...tsValue) tsValue { return tsU })
		constructor.name = "HTTPParser"; constructor.length = 0; constructor.constructible = true
		constructor.construct = func(loop *tsLoop, _ ...tsValue) tsValue { return tsNodeHTTPCommonParserCreate(loop) }
		constructor.properties = tsNewObject(); constructor.properties.set("prototype", tsObjectValue(prototype))
		prototype.set("constructor", tsFunctionValue(constructor))
		constants := map[string]float64{"REQUEST":1,"RESPONSE":2,"kOnMessageBegin":0,"kOnHeaders":1,"kOnHeadersComplete":2,"kOnBody":3,"kOnMessageComplete":4,"kOnExecute":5,"kOnTimeout":6,"kLenientNone":0,"kLenientHeaders":1,"kLenientChunkedLength":2,"kLenientKeepAlive":4,"kLenientTransferEncoding":8,"kLenientVersion":16,"kLenientDataAfterClose":32,"kLenientOptionalLFAfterCR":64,"kLenientOptionalCRLFAfterChunk":128,"kLenientOptionalCRBeforeLF":256,"kLenientSpacesAfterChunkSize":512,"kLenientHeaderValueRelaxed":1024,"kLenientAll":2047}
		for name, number := range constants { constructor.properties.set(name, tsNumberValue(number)) }
		tsNodeHTTPCommonSingleton.constructor = tsFunctionValue(constructor)
		connectionsPrototype := tsNewObject(); tsNodeHTTPCommonSingleton.connectionsPrototype = connectionsPrototype
		for name, length := range map[string]int{"all":0,"idle":0,"active":0,"expired":2} {
			methodName := name
			connectionsPrototype.set(name, tsNamedNativeMethod(name, length, func(loop *tsLoop, receiver tsValue, args ...tsValue) tsValue { return tsNodeHTTPCommonConnectionListMethod(loop, receiver, methodName, args) }))
		}
		connectionsTag := tsPropertyKey(tsWellKnownSymbol("toStringTag")); connectionsPrototype.set(connectionsTag, tsNodeHTTPCommonKey("ConnectionsList")); connectionsPrototype.descriptors = map[string]*tsDescriptor{connectionsTag:{configurable:true}}
		connectionsPrototype.set("constructor", tsU)
		connectionsConstructor := tsFunc(func(_ *tsLoop, _ ...tsValue) tsValue { return tsU })
		connectionsConstructor.name = "ConnectionsList"; connectionsConstructor.length = 2; connectionsConstructor.constructible = true
		connectionsConstructor.construct = func(loop *tsLoop, _ ...tsValue) tsValue { return tsNodeHTTPCommonConnectionListCreate(loop) }
		connectionsConstructor.properties = tsNewObject(); connectionsConstructor.properties.set("prototype", tsObjectValue(connectionsPrototype))
		connectionsPrototype.set("constructor", tsFunctionValue(connectionsConstructor))
		tsNodeHTTPCommonSingleton.connectionsConstructor = tsFunctionValue(connectionsConstructor)
	})
}
func tsNodeHTTPCommonParserMethod(loop *tsLoop, parser tsValue, name string, args []tsValue) tsValue {
	state := tsNodeHTTPCommonStateFor(parser)
	if state == nil { panic(tsThrown{tsFSException("TypeError", "ERR_INVALID_THIS", "HTTPParser method called on an incompatible receiver")}) }
	if state.closed { return tsU }
	switch name {
	case "close": state.closed = true; return tsU
	case "free": return tsU
	case "remove":
		if state != nil { tsNodeHTTPCommonConnectionListPop(loop,state.connections,parser,false);tsNodeHTTPCommonConnectionListPop(loop,state.connections,parser,true) }
		socket := tsNodeHTTPCommonGet(loop, parser, "socket"); if !tsNullish(socket) && tsStrictEqual(tsNodeHTTPCommonGet(loop, socket, "parser"), parser) { tsNodeHTTPCommonSet(loop, socket, "parser", tsNull) }; return tsU
	case "execute": return tsNodeHTTPCommonParserExecute(loop, parser, tsArg(args, 0))
	case "finish":
		if state.closed { return tsU }
		if state.paused { return tsU }
		if !state.initialized { panic(tsThrown{tsFSException("Error", "ERR_INVALID_STATE", "HTTPParser is not initialized")}) }
		if state != nil && state.closeDelimited {
			if !state.skipBody && len(state.buffer)-state.bodyStart > state.bodyDelivered {
				body := state.buffer[state.bodyStart+state.bodyDelivered:]
				tsNodeHTTPCommonCallback(loop, parser, 3, tsNodeHTTPCommonParserBufferValue(body), tsNumberValue(0))
				state.bodyDelivered = len(state.buffer)-state.bodyStart
				if state.paused { return tsNodeHTTPCommonParseError(loop, parser, len(state.buffer), "HPE_PAUSED", "Paused") }
			}
			tsNodeHTTPCommonCallback(loop, parser, 4)
			state.buffer = nil; state.headerDone = false; state.closeDelimited = false; state.closedByMessage = true; state.bodyStart = 0; state.bodyDelivered = 0; state.skipBody = false
			return tsU
		}
		if state != nil && len(state.buffer) == 0 { return tsU }
		return tsNodeHTTPCommonParseError(loop, parser, len(state.buffer), "HPE_INVALID_EOF_STATE", "Invalid EOF state")
	case "initialize":
		if state != nil {
			connections := tsArg(args, 4)
			if !tsNullish(connections) && tsNodeHTTPCommonConnectionListState(connections) == nil { return tsU }
			tsNodeHTTPCommonConnectionListPop(loop,state.connections,parser,false);tsNodeHTTPCommonConnectionListPop(loop,state.connections,parser,true)
			state.parserType=int(tsNumber(tsArg(args,0)));state.lenientFlags=uint32(tsToUint32(tsArg(args,3)));state.initialized=true;state.closed=false;state.started=time.Now();state.lastMessageStart=time.Time{};state.connections=connections
			state.maxHeaderSize=0
			if maxValue:=tsArg(args,2);maxValue.kind==tsNumberKind&&tsNumber(maxValue)>0{state.maxHeaderSize=uint64(tsNumber(maxValue))}
			if state.maxHeaderSize==0{if maxValue:=tsCall(loop,tsInternalHTTPGet(loop,tsInternalHTTPModule(loop),"getMaxHTTPHeaderSize"));tsNumber(maxValue)>0{state.maxHeaderSize=uint64(tsNumber(maxValue))}}
			state.buffer=nil;state.current=nil;state.paused=false;state.headerDone=false;state.messageBegun=false;state.consumed=false;state.skipBody=false;state.upgrade=false;state.closeDelimited=false;state.closedByMessage=false;state.bodyStart=0;state.bodyDelivered=0;state.bodyReceived=0;state.chunkState=0;state.chunkRemaining=0;state.chunkTrailers=nil
			if !tsNullish(state.connections){state.lastMessageStart=state.started;tsNodeHTTPCommonConnectionListPush(loop,state.connections,parser,false);tsNodeHTTPCommonConnectionListPush(loop,state.connections,parser,true)}else{state.connections=tsU}
			tsNodeHTTPCommonSet(loop, parser, "_type", tsArg(args, 0)); tsNodeHTTPCommonSet(loop, parser, "socket", tsArg(args, 1)); tsNodeHTTPCommonSet(loop, parser, "maxHeaderSize", tsArg(args, 2)); tsNodeHTTPCommonSet(loop, parser, "lenientFlags", tsArg(args, 3))
		}
		return tsU
	case "pause": state.paused = true; return tsU
	case "resume": state.paused = false; return tsU
	case "consume": state.consumed = true; tsNodeHTTPCommonSet(loop, parser, "_consumed", tsBooleanValue(true)); return tsU
	case "unconsume": state.consumed = false; tsNodeHTTPCommonSet(loop, parser, "_consumed", tsBooleanValue(false)); return tsU
	case "getCurrentBuffer": return tsNodeHTTPCommonParserBufferValue(state.current)
	case "duration": if state.lastMessageStart.IsZero() { return tsNumberValue(0) }; return tsNumberValue(float64(time.Since(state.lastMessageStart).Nanoseconds())/1e6)
	case "headersCompleted": return tsBooleanValue(state.headerDone)
	}
	return tsU
}
func tsNodeHTTPCommonParserOnHeaders(loop *tsLoop, parser tsValue, args ...tsValue) tsValue {
	headers := tsNodeHTTPCommonArray(tsArg(args, 0)).values
	return tsNodeHTTPCommonParserHeadersCallback(loop, parser, headers, tsArg(args, 1))
}
func tsNodeHTTPCommonParserOnHeadersComplete(loop *tsLoop, parser tsValue, args ...tsValue) tsValue {
	return tsNodeHTTPCommonHeadersComplete(loop, parser, args)
}
func tsNodeHTTPCommonParserOnBodyCallback(loop *tsLoop, parser tsValue, args ...tsValue) tsValue {
	return tsNodeHTTPCommonParserOnBody(loop, parser, tsArg(args, 0))
}
func tsNodeHTTPCommonParserOnMessageCompleteCallback(loop *tsLoop, parser tsValue, _ ...tsValue) tsValue {
	return tsNodeHTTPCommonParserOnMessageComplete(loop, parser)
}
func tsNodeHTTPCommonParserConstructorCallback(loop *tsLoop, _ tsValue, _ ...tsValue) tsValue {
	parser := tsNodeHTTPCommonParserCreate(loop)
	for index, callback := range map[int]tsValue{1:tsNamedNativeMethod("parserOnHeaders", 2, tsNodeHTTPCommonParserOnHeaders),2:tsNamedNativeMethod("parserOnHeadersComplete", 9, tsNodeHTTPCommonParserOnHeadersComplete),3:tsNamedNativeMethod("parserOnBody", 2, tsNodeHTTPCommonParserOnBodyCallback),4:tsNamedNativeMethod("parserOnMessageComplete", 0, tsNodeHTTPCommonParserOnMessageCompleteCallback)} { tsSet(loop, parser, tsNodeHTTPCommonIndex(index), callback) }
	return parser
}
func tsNodeHTTPCommonCloseParserInstance(loop *tsLoop, parser tsValue) { tsNodeHTTPCommonCall(loop, parser, "close") }
func tsNodeHTTPCommonFreeParser(loop *tsLoop, parser, req, socket tsValue) tsValue {
	if !tsNullish(parser) {
		if tsTruthy(tsNodeHTTPCommonGet(loop, parser, "_consumed")) { tsNodeHTTPCommonCall(loop, parser, "unconsume") }
		tsNodeHTTPCommonCall(loop, parser, "remove")
		tsNodeHTTPCommonCleanParser(loop, parser)
		freed := tsNodeHTTPCommonCall(loop, tsNodeHTTPCommonSingleton.parsers, "free", parser)
		if !tsTruthy(freed) { cb := tsFunctionValue(tsFunc(func(loop *tsLoop, _ ...tsValue) tsValue { tsNodeHTTPCommonCloseParserInstance(loop, parser); return tsU })); tsInternalHTTPNextTick(loop, cb) } else { tsNodeHTTPCommonCall(loop, parser, "free") }
	}
	if !tsNullish(req) { tsNodeHTTPCommonSet(loop, req, "parser", tsNull) }
	if !tsNullish(socket) { tsNodeHTTPCommonSet(loop, socket, "parser", tsNull) }
	return tsU
}
func tsNodeHTTPCommonPrepareError(loop *tsLoop, err, parser, rawPacket tsValue) tsValue {
	if !tsTruthy(rawPacket) { rawPacket = tsNodeHTTPCommonCall(loop, parser, "getCurrentBuffer") }
	tsSet(loop, err, tsNodeHTTPCommonKey("rawPacket"), rawPacket)
	reason := tsNodeHTTPCommonGet(loop, err, "reason")
	if reason.kind == tsStringKind { tsNodeHTTPCommonSet(loop, err, "message", tsNodeHTTPCommonKey("Parse Error: "+tsText(reason))) }
	return tsU
}
func tsNodeHTTPCommonInsecureParser(loop *tsLoop) bool {
	process := tsNodeProcessModule(loop); args := tsNodeHTTPCommonArray(tsNodeHTTPCommonGet(loop, process, "execArgv"))
	for _, arg := range args.values { if tsText(arg) == "--insecure-http-parser" { return true } }
	return false
}
func tsNodeHTTPCommonIsLenient(loop *tsLoop) tsValue {
	lenient := tsNodeHTTPCommonInsecureParser(loop)
	if lenient && !tsNodeHTTPCommonSingleton.warnedLenient {
		tsNodeHTTPCommonSingleton.warnedLenient = true
		process := tsNodeProcessModule(loop); tsNodeHTTPCommonCall(loop, process, "emitWarning", tsNodeHTTPCommonKey("Using insecure HTTP parsing"))
	}
	return tsBooleanValue(lenient)
}
func tsNodeHTTPCommonCalculateLenientFlags(loop *tsLoop, validation, insecureOption tsValue) tsValue {
	constructor := tsNodeHTTPCommonSingleton.constructor; constants := (*tsFunction)(constructor.ref).properties
	if tsText(validation) == "strict" { return constants.get("kLenientNone") }
	if tsText(validation) == "relaxed" { return constants.get("kLenientHeaderValueRelaxed") }
	if tsText(validation) == "insecure" { return constants.get("kLenientAll") }
	lenient := insecureOption
	if tsIsUndefined(lenient) { lenient = tsNodeHTTPCommonIsLenient(loop) }
	if tsTruthy(lenient) { return constants.get("kLenientAll") }
	return constants.get("kLenientNone")
}
func tsNodeHTTPCommonModule(loop *tsLoop) tsValue {
	tsNodeHTTPCommonParserInit(loop)
	if !tsIsUndefined(tsNodeHTTPCommonSingleton.module) { return tsNodeHTTPCommonSingleton.module }
	internal := tsInternalHTTPModule(loop)
	methods := tsNodeHTTPCommonGet(loop, internal, "METHODS")
	freeList := tsNodeModuleRequire(loop, []tsValue{tsNodeHTTPCommonKey("internal/freelist")})
	// FreeList invokes its parser factory with a constructor callback.
	callback := tsFunc(func(_ *tsLoop, _ ...tsValue) tsValue { return tsU })
	callback.constructible = true
	callback.construct = func(loop *tsLoop, _ ...tsValue) tsValue { return tsNodeHTTPCommonParserConstructorCallback(loop, tsU) }
	freeListValue := tsConstruct(freeList, tsNodeHTTPCommonKey("parsers"), tsNumberValue(1000), tsFunctionValue(callback))
	tsNodeHTTPCommonSingleton.parsers = freeListValue
	module := tsNewObject()
	allMethods := &tsArray{}; for _, name := range tsNodeHTTPCommonSingleton.allMethods { allMethods.values = append(allMethods.values, tsNodeHTTPCommonKey(name)) }
	module.set("methods", methods); module.set("allMethods", tsArrayValue(allMethods)); module.set("HTTPParser", tsNodeHTTPCommonSingleton.constructor); module.set("ConnectionsList", tsNodeHTTPCommonSingleton.connectionsConstructor); module.set("parsers", freeListValue)
	module.set("kIncomingMessage", tsNodeHTTPCommonSingleton.kIncomingMessage); module.set("kSkipPendingData", tsNodeHTTPCommonSingleton.kSkipPendingData); module.set("MAX_HEADER_PAIRS", tsNumberValue(2000))
	module.set("validateHeaderName", tsNamedNativeMethod("validateHeaderName", 2, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue { return tsNodeHTTPCommonValidateHeaderName(loop, tsArg(args, 0), tsArg(args, 1)) }))
	module.set("validateHeaderValue", tsNamedNativeMethod("validateHeaderValue", 2, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue { return tsNodeHTTPCommonValidateHeaderValue(loop, tsArg(args, 0), tsArg(args, 1)) }))
	module.set("_checkIsHttpToken", tsNamedNativeMethod("_checkIsHttpToken", 1, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue { check := tsGet(loop, tsInternalValidatorsModule(loop), tsNodeHTTPCommonKey("checkIsHttpToken")); return tsCall(loop, check, tsArg(args, 0)) }))
	module.set("_checkInvalidHeaderChar", tsNamedNativeMethod("_checkInvalidHeaderChar", 2, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue { if tsTruthy(tsArg(args, 1)) { return tsBooleanValue(tsNodeHTTPCommonInvalidHeaderValue(loop, tsArg(args, 0), tsU)) }; return tsBooleanValue(tsNodeHTTPCommonStrictHeaderValue(loop, tsArg(args, 0))) }))
	module.set("CRLF", tsNodeHTTPCommonKey("\r\n"))
	module.set("freeParser", tsNamedNativeMethod("freeParser", 3, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue { return tsNodeHTTPCommonFreeParser(loop, tsArg(args, 0), tsArg(args, 1), tsArg(args, 2)) }))
	module.set("isLenient", tsNamedNativeMethod("isLenient", 0, func(loop *tsLoop, _ tsValue, _ ...tsValue) tsValue { return tsNodeHTTPCommonIsLenient(loop) }))
	module.set("calculateLenientFlags", tsNamedNativeMethod("calculateLenientFlags", 2, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue { return tsNodeHTTPCommonCalculateLenientFlags(loop, tsArg(args, 0), tsArg(args, 1)) }))
	module.set("prepareError", tsNamedNativeMethod("prepareError", 3, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue { return tsNodeHTTPCommonPrepareError(loop, tsArg(args, 0), tsArg(args, 1), tsArg(args, 2)) }))
	module.set("default", tsObjectValue(module)); tsNodeHTTPCommonSingleton.module = tsObjectValue(module)
	return tsObjectValue(module)
}
`
