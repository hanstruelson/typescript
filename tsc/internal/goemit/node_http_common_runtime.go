package goemit

// nodeHTTPCommonRuntime adapts Node's _http_common helpers and the HTTPParser
// binding boundary with a compact incremental HTTP/1 parser.
const nodeHTTPCommonRuntime = `
type tsNodeHTTPParserState struct {
	buffer []byte
	current []byte
	started time.Time
	parserType int
	paused, closed, initialized, headerDone, consumed, skipBody, closeDelimited bool
	bodyStart int
}

var tsNodeHTTPCommonSingleton struct {
	sync.Once
	module, constructor, parsers tsValue
	prototype *tsObject
	states map[uintptr]*tsNodeHTTPParserState
	kIncomingMessage, kSkipPendingData tsValue
	methods []string
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
func tsNodeHTTPCommonHeaderPairs(lines []string, max int) ([]tsValue, string) {
	pairs := make([]tsValue, 0, len(lines)*2)
	for _, line := range lines {
		colon := strings.IndexByte(line, ':')
		if colon <= 0 { return nil, line }
		if max > 0 && len(pairs)/2 >= max { break }
		value := strings.TrimLeft(line[colon+1:], " \t")
		pairs = append(pairs, tsNodeHTTPCommonKey(line[:colon]), tsNodeHTTPCommonKey(value))
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
		index := int(tsNumber(method)); if index >= 0 && index < len(tsNodeHTTPCommonSingleton.methods) { tsNodeHTTPCommonSet(loop, incoming, "method", tsNodeHTTPCommonKey(tsNodeHTTPCommonSingleton.methods[index])) }
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
func tsNodeHTTPCommonCleanParser(loop *tsLoop, parser tsValue) {
	tsNodeHTTPCommonSet(loop, parser, "_headers", tsArrayValue(&tsArray{}))
	tsNodeHTTPCommonSet(loop, parser, "_url", tsNodeHTTPCommonKey(""))
	for _, name := range []string{"socket", "incoming", "outgoing", "onIncoming", "joinDuplicateHeaders"} { tsNodeHTTPCommonSet(loop, parser, name, tsNull) }
	tsNodeHTTPCommonSet(loop, parser, "maxHeaderPairs", tsNumberValue(2000))
	tsNodeHTTPCommonSet(loop, parser, "_consumed", tsBooleanValue(false))
	for _, index := range []int{0, 5, 6} { tsSet(loop, parser, tsNodeHTTPCommonIndex(index), tsNull) }
	if state := tsNodeHTTPCommonStateFor(parser); state != nil { state.buffer = nil; state.current = nil; state.paused = false; state.initialized = false; state.headerDone = false; state.consumed = false; state.skipBody = false; state.closeDelimited = false; state.bodyStart = 0; state.started = time.Time{} }
}
func tsNodeHTTPCommonChunked(data []byte) ([]byte, int, []tsValue, bool, string) {
	var body []byte
	trailers := []tsValue{}
	for offset := 0; ; {
		lineEnd := bytes.Index(data[offset:], []byte("\r\n")); if lineEnd < 0 { return body, 0, nil, false, "Invalid chunk size" }
		lineEnd += offset; sizeText := string(data[offset:lineEnd]); if semi := strings.IndexByte(sizeText, ';'); semi >= 0 { sizeText = sizeText[:semi] }
		size, err := strconv.ParseInt(strings.TrimSpace(sizeText), 16, 64); if err != nil || size < 0 { return body, 0, nil, false, "Invalid chunk size" }
		offset = lineEnd + 2
		if size == 0 {
			for {
				end := bytes.Index(data[offset:], []byte("\r\n")); if end < 0 { return body, 0, nil, false, "Invalid chunk trailer" }
				end += offset
				if end == offset { return body, end + 2, trailers, true, "" }
				name, value, ok := strings.Cut(string(data[offset:end]), ":"); if !ok { return body, 0, nil, false, "Invalid chunk trailer" }
				trailers = append(trailers, tsNodeHTTPCommonKey(name), tsNodeHTTPCommonKey(strings.TrimLeft(value, " \t")))
				offset = end + 2
			}
		}
		if size > int64(len(data)-offset-2) { return body, 0, nil, false, "Incomplete chunk" }
		body = append(body, data[offset:offset+int(size)]...)
		offset += int(size)
		if offset+2 > len(data) || data[offset] != '\r' || data[offset+1] != '\n' { return body, 0, nil, false, "Invalid chunk terminator" }
		offset += 2
	}
}
func tsNodeHTTPCommonParserExecute(loop *tsLoop, parser tsValue, data tsValue) tsValue {
	state := tsNodeHTTPCommonStateFor(parser); if state == nil { return tsNodeHTTPCommonParseError(loop, parser, 0, "HPE_INVALID_INTERNAL_STATE", "Parser is not initialized") }
	chunk := tsFSBytes(data); state.current = append(state.current[:0], chunk...); state.buffer = append(state.buffer, chunk...)
	boundary := bytes.Index(state.buffer, []byte("\r\n\r\n"))
	if boundary < 0 {
		max := tsInternalHTTPGet(loop, tsInternalHTTPModule(loop), "getMaxHTTPHeaderSize"); if tsNumber(max) > 0 && len(state.buffer) > int(tsNumber(max)) { return tsNodeHTTPCommonParseError(loop, parser, 0, "HPE_HEADER_OVERFLOW", "Header overflow") }
		return tsNumberValue(float64(len(chunk)))
	}
	headerEnd := boundary + 4
	if max := tsInternalHTTPGet(loop, tsInternalHTTPModule(loop), "getMaxHTTPHeaderSize"); tsNumber(max) > 0 && headerEnd > int(tsNumber(max)) { return tsNodeHTTPCommonParseError(loop, parser, boundary, "HPE_HEADER_OVERFLOW", "Header overflow") }
	alreadyParsedHeaders := state.headerDone
	state.headerDone = true
	lines := strings.Split(string(state.buffer[:boundary]), "\r\n")
	if len(lines) == 0 || lines[0] == "" { state.headerDone = false; return tsNodeHTTPCommonParseError(loop, parser, 0, "HPE_INVALID_CONSTANT", "Invalid HTTP message") }
	major, minor := 1, 1; methodValue, urlValue, status, statusMessage := tsU, tsNodeHTTPCommonKey(""), tsU, tsU
	parts := strings.Fields(lines[0])
	if len(parts) < 2 { state.headerDone = false; return tsNodeHTTPCommonParseError(loop, parser, 0, "HPE_INVALID_CONSTANT", "Invalid start line") }
	if strings.HasPrefix(parts[0], "HTTP/") {
		version := strings.TrimPrefix(parts[0], "HTTP/"); _, _ = fmt.Sscanf(version, "%d.%d", &major, &minor)
		code, _ := strconv.Atoi(parts[1]); status = tsNumberValue(float64(code)); if len(parts) > 2 { statusMessage = tsNodeHTTPCommonKey(strings.Join(parts[2:], " ")) }
	} else {
		methodValue = tsU
		for index, method := range tsNodeHTTPCommonSingleton.methods { if method == parts[0] { methodValue = tsNumberValue(float64(index)); break } }
		if tsIsUndefined(methodValue) { state.headerDone = false; return tsNodeHTTPCommonParseError(loop, parser, 0, "HPE_INVALID_METHOD", "Invalid method") }
		urlValue = tsNodeHTTPCommonKey(parts[1]); if len(parts) > 2 && strings.HasPrefix(parts[2], "HTTP/") { _, _ = fmt.Sscanf(strings.TrimPrefix(parts[2], "HTTP/"), "%d.%d", &major, &minor) }
	}
	pairs, badLine := tsNodeHTTPCommonHeaderPairs(lines[1:], 0)
	if badLine != "" { state.headerDone = false; return tsNodeHTTPCommonParseError(loop, parser, 0, "HPE_INVALID_HEADER_TOKEN", "Invalid header") }
	flat := tsArrayValue(&tsArray{values:pairs})
	connection, transfer, lengthValue, upgradeValue := "", "", "", ""
	for index := 0; index+1 < len(pairs); index += 2 {
		key := strings.ToLower(tsText(pairs[index])); value := tsText(pairs[index+1])
		switch key { case "connection": connection = strings.ToLower(value); case "transfer-encoding": transfer = strings.ToLower(value); case "content-length": lengthValue = value; case "upgrade": upgradeValue = value }
	}
	shouldKeepAlive := tsBooleanValue(major > 1 || major == 1 && minor >= 1)
	if strings.Contains(connection, "close") { shouldKeepAlive = tsBooleanValue(false) } else if strings.Contains(connection, "keep-alive") { shouldKeepAlive = tsBooleanValue(true) }
	upgrade := tsBooleanValue(strings.Contains(connection, "upgrade") && upgradeValue != "")
	if !alreadyParsedHeaders {
		beginCallback := tsGet(loop, parser, tsNodeHTTPCommonIndex(0)); if beginCallback.kind == tsFunctionKind { tsCallReceiver(loop, beginCallback, parser) }
		completeArgs := []tsValue{tsNumberValue(float64(major)), tsNumberValue(float64(minor)), flat, methodValue, urlValue, status, statusMessage, upgrade, shouldKeepAlive}
		callbackResult := tsNodeHTTPCommonCallback(loop, parser, 2, completeArgs...)
		state.skipBody = tsTruthy(callbackResult)
	}
	bodyStart := headerEnd
	if !alreadyParsedHeaders {
		state.bodyStart = bodyStart
		noBodyStatus := status.kind == tsNumberKind && (int(tsNumber(status)) < 200 || int(tsNumber(status)) == 204 || int(tsNumber(status)) == 304)
		state.closeDelimited = state.parserType == 2 && lengthValue == "" && !strings.Contains(transfer, "chunked") && !noBodyStatus
	}
	consumedBody := 0; var body []byte; var trailers []tsValue; bodyComplete := state.skipBody
	if state.skipBody {
		bodyComplete = true
	} else if strings.Contains(transfer, "chunked") {
		body, consumedBody, trailers, bodyComplete, _ = tsNodeHTTPCommonChunked(state.buffer[bodyStart:])
	} else if lengthValue != "" {
		length, err := strconv.Atoi(strings.TrimSpace(lengthValue)); if err != nil || length < 0 { state.headerDone = false; return tsNodeHTTPCommonParseError(loop, parser, boundary, "HPE_INVALID_CONTENT_LENGTH", "Invalid content length") }
		if len(state.buffer)-bodyStart >= length { body = append([]byte(nil), state.buffer[bodyStart:bodyStart+length]...); consumedBody = length; bodyComplete = true }
	} else if state.parserType == 1 || status.kind == tsNumberKind && (int(tsNumber(status)) < 200 || int(tsNumber(status)) == 204 || int(tsNumber(status)) == 304) {
		bodyComplete = true
	}
	if len(body) > 0 { tsNodeHTTPCommonCallback(loop, parser, 3, tsNodeHTTPCommonParserBufferValue(body), tsNumberValue(0)) }
	if !bodyComplete { return tsNumberValue(float64(len(chunk))) }
	if len(trailers) > 0 { tsNodeHTTPCommonSet(loop, parser, "_headers", tsArrayValue(&tsArray{values:trailers})) }
	messageComplete := tsGet(loop, parser, tsNodeHTTPCommonIndex(4)); if messageComplete.kind == tsFunctionKind { tsCallReceiver(loop, messageComplete, parser) }
	parsed := bodyStart + consumedBody
	if strings.Contains(transfer, "chunked") { parsed = bodyStart + consumedBody }
	if parsed > len(state.buffer) { parsed = len(state.buffer) }
	state.buffer = append([]byte(nil), state.buffer[parsed:]...); state.headerDone = false; state.skipBody = false; state.closeDelimited = false; state.bodyStart = 0
	if execute := tsGet(loop, parser, tsNodeHTTPCommonIndex(5)); execute.kind == tsFunctionKind { tsCallReceiver(loop, execute, parser, tsNumberValue(float64(parsed))) }
	return tsNumberValue(float64(parsed))
}
func tsNodeHTTPCommonParserCreate(loop *tsLoop) tsValue {
	parser := tsNewObject(); parser.prototype = tsObjectValue(tsNodeHTTPCommonSingleton.prototype)
	state := &tsNodeHTTPParserState{started:time.Now()}; tsNodeHTTPCommonSingleton.states[uintptr(parser.ref)] = state
	tsNodeHTTPCommonCleanParser(loop, tsObjectValue(parser))
	tsNodeHTTPCommonSet(loop, tsObjectValue(parser), "_type", tsNumberValue(0))
	return tsObjectValue(parser)
}
func tsNodeHTTPCommonParserInit(loop *tsLoop) {
	tsNodeHTTPCommonSingleton.Do(func() {
		tsNodeHTTPCommonSingleton.states = map[uintptr]*tsNodeHTTPParserState{}
		tsNodeHTTPCommonSingleton.kIncomingMessage = tsNewSymbol(tsStringUTF8("IncomingMessage"))
		tsNodeHTTPCommonSingleton.kSkipPendingData = tsNewSymbol(tsStringUTF8("SkipPendingData"))
		methods := tsNodeHTTPCommonArray(tsInternalHTTPGet(loop, tsInternalHTTPModule(loop), "METHODS"))
		for _, value := range methods.values { tsNodeHTTPCommonSingleton.methods = append(tsNodeHTTPCommonSingleton.methods, tsText(value)) }
		prototype := tsNewObject(); tsNodeHTTPCommonSingleton.prototype = prototype
		for name, length := range map[string]int{"close":0,"free":0,"remove":0,"finish":0,"pause":0,"resume":0,"unconsume":0,"getCurrentBuffer":0,"duration":0,"headersCompleted":0,"consume":1,"execute":1,"initialize":4} {
			methodName := name
			prototype.set(name, tsNamedNativeMethod(name, length, func(loop *tsLoop, receiver tsValue, args ...tsValue) tsValue { return tsNodeHTTPCommonParserMethod(loop, receiver, methodName, args) }))
		}
		constructor := tsFunc(func(_ *tsLoop, _ ...tsValue) tsValue { tsPropertyFailure("Class constructor HTTPParser cannot be invoked without new"); return tsU })
		constructor.name = "HTTPParser"; constructor.length = 0; constructor.constructible = true
		constructor.construct = func(loop *tsLoop, _ ...tsValue) tsValue { return tsNodeHTTPCommonParserCreate(loop) }
		constructor.properties = tsNewObject(); constructor.properties.set("prototype", tsObjectValue(prototype))
		prototype.set("constructor", tsFunctionValue(constructor))
		constants := map[string]float64{"REQUEST":1,"RESPONSE":2,"kOnMessageBegin":0,"kOnHeaders":1,"kOnHeadersComplete":2,"kOnBody":3,"kOnMessageComplete":4,"kOnExecute":5,"kOnTimeout":6,"kLenientNone":0,"kLenientHeaders":1,"kLenientChunkedLength":2,"kLenientKeepAlive":4,"kLenientTransferEncoding":8,"kLenientVersion":16,"kLenientDataAfterClose":32,"kLenientOptionalLFAfterCR":64,"kLenientOptionalCRLFAfterChunk":128,"kLenientOptionalCRBeforeLF":256,"kLenientSpacesAfterChunkSize":512,"kLenientHeaderValueRelaxed":1024,"kLenientAll":2047}
		for name, number := range constants { constructor.properties.set(name, tsNumberValue(number)) }
		tsNodeHTTPCommonSingleton.constructor = tsFunctionValue(constructor)
	})
}
func tsNodeHTTPCommonParserMethod(loop *tsLoop, parser tsValue, name string, args []tsValue) tsValue {
	state := tsNodeHTTPCommonStateFor(parser)
	switch name {
	case "close": if state != nil { state.closed = true }; return tsU
	case "free": if state != nil { state.buffer = nil; state.current = nil }; return tsU
	case "remove": socket := tsNodeHTTPCommonGet(loop, parser, "socket"); if !tsNullish(socket) && tsStrictEqual(tsNodeHTTPCommonGet(loop, socket, "parser"), parser) { tsNodeHTTPCommonSet(loop, socket, "parser", tsNull) }; return tsU
	case "execute": return tsNodeHTTPCommonParserExecute(loop, parser, tsArg(args, 0))
	case "finish":
		if state != nil && state.closeDelimited {
			if state.bodyStart < len(state.buffer) { tsNodeHTTPCommonCallback(loop, parser, 3, tsNodeHTTPCommonParserBufferValue(state.buffer[state.bodyStart:]), tsNumberValue(0)) }
			tsNodeHTTPCommonCallback(loop, parser, 4)
			state.buffer = nil; state.headerDone = false; state.closeDelimited = false; state.bodyStart = 0
			return tsU
		}
		if state != nil && len(state.buffer) == 0 { return tsU }
		return tsNodeHTTPCommonParseError(loop, parser, len(state.buffer), "HPE_INVALID_EOF_STATE", "Invalid EOF state")
	case "initialize": if state != nil { state.parserType = int(tsNumber(tsArg(args, 0))); state.initialized = true; state.started = time.Now(); tsNodeHTTPCommonSet(loop, parser, "_type", tsArg(args, 0)); tsNodeHTTPCommonSet(loop, parser, "socket", tsArg(args, 1)); tsNodeHTTPCommonSet(loop, parser, "maxHeaderSize", tsArg(args, 2)); tsNodeHTTPCommonSet(loop, parser, "lenientFlags", tsArg(args, 3)) }; return tsU
	case "pause": if state != nil { state.paused = true }; return tsU
	case "resume": if state != nil { state.paused = false }; return tsU
	case "consume": if state != nil { state.consumed = true }; tsNodeHTTPCommonSet(loop, parser, "_consumed", tsBooleanValue(true)); return tsU
	case "unconsume": if state != nil { state.consumed = false }; tsNodeHTTPCommonSet(loop, parser, "_consumed", tsBooleanValue(false)); return tsU
	case "getCurrentBuffer": if state == nil { return tsU }; return tsNodeHTTPCommonParserBufferValue(state.current)
	case "duration": if state == nil || state.started.IsZero() { return tsU }; return tsNumberValue(float64(time.Since(state.started).Nanoseconds())/1e6)
	case "headersCompleted": if state == nil { return tsU }; return tsBooleanValue(state.headerDone)
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
	for index, callback := range map[int]tsValue{1:tsFunctionValue(tsNamedNativeMethod("parserOnHeaders", 2, tsNodeHTTPCommonParserOnHeaders)),2:tsFunctionValue(tsNamedNativeMethod("parserOnHeadersComplete", 9, tsNodeHTTPCommonParserOnHeadersComplete)),3:tsFunctionValue(tsNamedNativeMethod("parserOnBody", 2, tsNodeHTTPCommonParserOnBodyCallback)),4:tsFunctionValue(tsNamedNativeMethod("parserOnMessageComplete", 0, tsNodeHTTPCommonParserOnMessageCompleteCallback))} { tsSet(loop, parser, tsNodeHTTPCommonIndex(index), callback) }
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
	module.set("methods", methods); module.set("allMethods", methods); module.set("HTTPParser", tsNodeHTTPCommonSingleton.constructor); module.set("ConnectionsList", tsU); module.set("parsers", freeListValue)
	module.set("kIncomingMessage", tsNodeHTTPCommonSingleton.kIncomingMessage); module.set("kSkipPendingData", tsNodeHTTPCommonSingleton.kSkipPendingData); module.set("MAX_HEADER_PAIRS", tsNumberValue(2000))
	module.set("validateHeaderName", tsNamedNativeMethod("validateHeaderName", 2, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue { return tsNodeHTTPCommonValidateHeaderName(loop, tsArg(args, 0), tsArg(args, 1)) }))
	module.set("validateHeaderValue", tsNamedNativeMethod("validateHeaderValue", 2, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue { return tsNodeHTTPCommonValidateHeaderValue(loop, tsArg(args, 0), tsArg(args, 1)) }))
	module.set("_checkIsHttpToken", tsNamedNativeMethod("_checkIsHttpToken", 1, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue { check := tsGet(loop, tsInternalValidatorsModule(loop), tsNodeHTTPCommonKey("checkIsHttpToken")); return tsCall(loop, check, tsArg(args, 0)) }))
	module.set("_checkInvalidHeaderChar", tsNamedNativeMethod("_checkInvalidHeaderChar", 2, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue { if tsTruthy(tsArg(args, 1)) { return tsBooleanValue(tsNodeHTTPCommonInvalidHeaderValue(loop, tsArg(args, 0), tsU)) }; return tsBooleanValue(tsNodeHTTPCommonStrictHeaderValue(loop, tsArg(args, 0))) }))
	module.set("chunkExpression", tsConstruct(tsRegExpConstructor(), tsNodeHTTPCommonKey("(?:^|\\W)chunked(?:$|\\W)"), tsNodeHTTPCommonKey("i")))
	module.set("continueExpression", tsConstruct(tsRegExpConstructor(), tsNodeHTTPCommonKey("(?:^|\\W)100-continue(?:$|\\W)"), tsNodeHTTPCommonKey("i")))
	module.set("CRLF", tsNodeHTTPCommonKey("\r\n"))
	module.set("freeParser", tsNamedNativeMethod("freeParser", 3, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue { return tsNodeHTTPCommonFreeParser(loop, tsArg(args, 0), tsArg(args, 1), tsArg(args, 2)) }))
	module.set("isLenient", tsNamedNativeMethod("isLenient", 0, func(loop *tsLoop, _ tsValue, _ ...tsValue) tsValue { return tsNodeHTTPCommonIsLenient(loop) }))
	module.set("calculateLenientFlags", tsNamedNativeMethod("calculateLenientFlags", 2, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue { return tsNodeHTTPCommonCalculateLenientFlags(loop, tsArg(args, 0), tsArg(args, 1)) }))
	module.set("prepareError", tsNamedNativeMethod("prepareError", 3, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue { return tsNodeHTTPCommonPrepareError(loop, tsArg(args, 0), tsArg(args, 1), tsArg(args, 2)) }))
	module.set("default", tsObjectValue(module)); tsNodeHTTPCommonSingleton.module = tsObjectValue(module)
	return tsObjectValue(module)
}
`
