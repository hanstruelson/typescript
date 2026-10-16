package goemit

// internalHTTP1ServerFallbackRuntime ports the standalone response handle
// used by Bun's HTTP/1 fallback before its lazily loaded connection parser.
const internalHTTP1ServerFallbackRuntime = `
type tsHTTP1FallbackResponseState struct {
	socket tsValue
	handle tsValue
	head tsValue
	headWritten bool
	chunked bool
	noBody bool
	closeDelimited bool
	shouldKeepAlive bool
	keepAliveTimeout float64
	onWritable tsValue
	onWritableFrame tsValue
	onFinished tsValue
	finished bool
	ended bool
	aborted bool
	bytesWritten uint64
}

var tsInternalHTTP1FallbackSingleton struct {
	sync.Once
	kConnections, kActiveRequests, kMessageOpen, kPendingCallbacks tsValue
	kQueue, kQueueState, kMustCloseConnection tsValue
	module tsValue
}

func tsInternalHTTP1FallbackKey(name string) tsValue { return tsStringReference(tsStringUTF8(name)) }
func tsInternalHTTP1FallbackGet(loop *tsLoop, value tsValue, name string) tsValue { return tsGet(loop, value, tsInternalHTTP1FallbackKey(name)) }
func tsInternalHTTP1FallbackSet(loop *tsLoop, value tsValue, name string, next tsValue) { tsSet(loop, value, tsInternalHTTP1FallbackKey(name), next) }
func tsInternalHTTP1FallbackSym(loop *tsLoop, value tsValue, name string) tsValue {
	tsInternalHTTP1FallbackInit(loop)
	var key tsValue
	switch name {
	case "kHttp1Connections": key = tsInternalHTTP1FallbackSingleton.kConnections
	case "kHttp1ActiveRequests": key = tsInternalHTTP1FallbackSingleton.kActiveRequests
	case "kHttp1MessageOpen": key = tsInternalHTTP1FallbackSingleton.kMessageOpen
	case "kPendingCallbacks": key = tsInternalHTTP1FallbackSingleton.kPendingCallbacks
	case "kQueue": key = tsInternalHTTP1FallbackSingleton.kQueue
	case "kQueueState": key = tsInternalHTTP1FallbackSingleton.kQueueState
	case "kMustCloseConnection": key = tsInternalHTTP1FallbackSingleton.kMustCloseConnection
	}
	return tsGet(loop, value, key)
}
func tsInternalHTTP1FallbackPut(loop *tsLoop, value tsValue, name string, next tsValue) {
	tsInternalHTTP1FallbackInit(loop)
	var key tsValue
	switch name {
	case "kHttp1Connections": key = tsInternalHTTP1FallbackSingleton.kConnections
	case "kHttp1ActiveRequests": key = tsInternalHTTP1FallbackSingleton.kActiveRequests
	case "kHttp1MessageOpen": key = tsInternalHTTP1FallbackSingleton.kMessageOpen
	case "kPendingCallbacks": key = tsInternalHTTP1FallbackSingleton.kPendingCallbacks
	case "kQueue": key = tsInternalHTTP1FallbackSingleton.kQueue
	case "kQueueState": key = tsInternalHTTP1FallbackSingleton.kQueueState
	case "kMustCloseConnection": key = tsInternalHTTP1FallbackSingleton.kMustCloseConnection
	}
	tsSet(loop, value, key, next)
}
func tsInternalHTTP1FallbackInit(loop *tsLoop) {
	tsInternalHTTP1FallbackSingleton.Do(func() {
		symbolFor := tsGet(loop, tsSymbolFunction(), tsInternalHTTP1FallbackKey("for"))
		tsInternalHTTP1FallbackSingleton.kConnections = tsCall(loop, symbolFor, tsInternalHTTP1FallbackKey("http1Connections"))
		tsInternalHTTP1FallbackSingleton.kActiveRequests = tsCall(loop, symbolFor, tsInternalHTTP1FallbackKey("http1ActiveRequests"))
		tsInternalHTTP1FallbackSingleton.kMessageOpen = tsCall(loop, symbolFor, tsInternalHTTP1FallbackKey("http1MessageOpen"))
		tsInternalHTTP1FallbackSingleton.kPendingCallbacks = tsInternalHTTPGet(loop, tsInternalHTTPModule(loop), "kPendingCallbacks")
		tsInternalHTTP1FallbackSingleton.kQueue = tsCall(loop, symbolFor, tsInternalHTTP1FallbackKey("http1ResponseQueue"))
		tsInternalHTTP1FallbackSingleton.kQueueState = tsCall(loop, symbolFor, tsInternalHTTP1FallbackKey("http1QueuedResponse"))
		tsInternalHTTP1FallbackSingleton.kMustCloseConnection = tsCall(loop, symbolFor, tsInternalHTTP1FallbackKey("kMustCloseConnection"))
	})
}
func tsInternalHTTP1FallbackRethrow(loop *tsLoop, err tsValue) {
	callback := tsFunctionValue(tsFunc(func(loop *tsLoop, args ...tsValue) tsValue { panic(tsThrown{tsArg(args, 0)}) }))
	tsInternalHTTPNextTick(loop, callback, err)
}
func tsInternalHTTP1FallbackSafeCallback(loop *tsLoop, callback tsValue, args ...tsValue) {
	if callback.kind != tsFunctionKind { return }
	var thrown any
	func() { defer func() { thrown = recover() }(); tsCall(loop, callback, args...) }()
	if thrown != nil { tsInternalHTTP1FallbackRethrow(loop, tsUnwrap(thrown)) }
}
func tsInternalHTTP1FallbackFailPending(loop *tsLoop, response, err tsValue) {
	tsInternalHTTP1FallbackInit(loop)
	callbacks := tsGet(loop, response, tsInternalHTTP1FallbackSingleton.kPendingCallbacks)
	tsSet(loop, response, tsInternalHTTP1FallbackSingleton.kPendingCallbacks, tsArrayValue(&tsArray{}))
	if callbacks.kind != tsArrayKind { return }
	for _, callback := range (*tsArray)(callbacks.ref).values { tsInternalHTTP1FallbackSafeCallback(loop, callback, err) }
}
func tsInternalHTTP1FallbackDate() string { return time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT") }
func tsInternalHTTP1FallbackToBuffer(loop *tsLoop, chunk, encoding tsValue) tsValue {
	if tsNullish(chunk) { return tsNull }
	if chunk.kind == tsStringKind { return tsNodeBufferFrom(loop, []tsValue{chunk, encoding}) }
	return chunk
}
func tsInternalHTTP1FallbackString(loop *tsLoop, value tsValue) string { return tsNodeHTTPOutgoingText(loop, value) }
func tsInternalHTTP1FallbackSocketWrite(loop *tsLoop, state *tsHTTP1FallbackResponseState, data, callback tsValue) tsValue {
	if tsTruthy(tsInternalHTTP1FallbackGet(loop, state.socket, "writableEnded")) { return tsBooleanValue(false) }
	if callback.kind == tsFunctionKind { return tsNodeHTTPServerCall(loop, state.socket, "write", data, callback) }
	return tsNodeHTTPServerCall(loop, state.socket, "write", data)
}
func tsInternalHTTP1FallbackWriteBody(loop *tsLoop, state *tsHTTP1FallbackResponseState, buffer, callback tsValue) int {
	length := 0
	if !tsNullish(buffer) { length = len(tsNodeBufferVisibleBytes(buffer)) }
	if length > 0 {
		if state.chunked {
			tsInternalHTTP1FallbackSocketWrite(loop, state, tsInternalHTTP1FallbackKey(strconv.FormatInt(int64(length), 16)+"\r\n"), tsU)
			tsInternalHTTP1FallbackSocketWrite(loop, state, buffer, tsU)
			tsInternalHTTP1FallbackSocketWrite(loop, state, tsInternalHTTP1FallbackKey("\r\n"), callback)
		} else { tsInternalHTTP1FallbackSocketWrite(loop, state, buffer, callback) }
	}
	return length
}
func tsInternalHTTP1FallbackRunFrame(loop *tsLoop, frame, callback tsValue) {
	module := tsInternalAsyncContextFrameModule(loop)
	run := tsGet(loop, module, tsInternalHTTP1FallbackKey("run"))
	if run.kind == tsFunctionKind { tsCall(loop, run, frame, callback) } else { tsCall(loop, callback) }
}
func tsInternalHTTP1FallbackSettleWritable(loop *tsLoop, state *tsHTTP1FallbackResponseState, callback tsValue) {
	frame := state.onWritableFrame
	state.onWritable = tsNull
	state.onWritableFrame = tsU
	tsInternalHTTP1FallbackRunFrame(loop, frame, tsFunctionValue(tsFunc(func(loop *tsLoop, _ ...tsValue) tsValue { tsInternalHTTP1FallbackSafeCallback(loop, callback); return tsU })))
}
func tsInternalHTTP1FallbackFlush(loop *tsLoop, state *tsHTTP1FallbackResponseState) {
	if !state.ended || state.finished { return }
	state.finished = true
	tsInternalHTTP1FallbackSafeCallback(loop, tsInternalHTTP1FallbackHandleGet(loop, state.handle, "onflushed"))
}
func tsInternalHTTP1FallbackHandleGet(loop *tsLoop, value tsValue, name string) tsValue {
	if value.kind != tsObjectKind { return tsU }
	return (*tsObject)(value.ref).get(name)
}
func tsInternalHTTP1FallbackEndWritten(loop *tsLoop, state *tsHTTP1FallbackResponseState, err tsValue) {
	if !tsTruthy(err) && !tsTruthy(tsInternalHTTP1FallbackGet(loop, state.socket, "destroyed")) { tsInternalHTTP1FallbackFlush(loop, state) }
}
func tsInternalHTTP1FallbackCountBodyBytes(loop *tsLoop, state *tsHTTP1FallbackResponseState, buffer, expected tsValue, fromEnd bool) int {
	length := 0
	if !tsNullish(buffer) { length = len(tsNodeBufferVisibleBytes(buffer)) }
	written := state.bytesWritten + uint64(length)
	if tsIsNumeric(expected) {
		want := uint64(math.Max(0, math.Trunc(tsNumber(expected))))
		if fromEnd && written != want || !fromEnd && written > want {
			message := fmt.Sprintf("Response body's content-length of %d byte(s) does not match the content-length of %d byte(s) set in header", written, want)
			panic(tsThrown{tsFSException("Error", "ERR_HTTP_CONTENT_LENGTH_MISMATCH", message)})
		}
	}
	state.bytesWritten = written
	return length
}
func tsInternalHTTP1FallbackWriteHead(loop *tsLoop, state *tsHTTP1FallbackResponseState, contentLength tsValue) {
	if state.headWritten { return }
	state.headWritten = true
	statusCode := 200
	statusMessage := tsU
	if !tsNullish(state.head) { statusCode = int(tsNumber(tsInternalHTTP1FallbackGet(loop, state.head, "statusCode"))); statusMessage = tsInternalHTTP1FallbackGet(loop, state.head, "statusMessage") }
	if statusMessage.kind != tsStringKind || tsText(statusMessage) == "" {
		codes := tsInternalHTTPGet(loop, tsInternalHTTPModule(loop), "STATUS_CODES")
		statusMessage = tsGet(loop, codes, tsInternalHTTP1FallbackKey(strconv.Itoa(statusCode)))
		if !tsTruthy(statusMessage) { statusMessage = tsInternalHTTP1FallbackKey("unknown") }
	}
	output := "HTTP/1.1 " + strconv.Itoa(statusCode) + " " + tsText(statusMessage) + "\r\n"
	hasContentLength, hasTransferEncoding, hasDate, hasConnection, hasKeepAlive := false, false, false, false, false
	headers := tsU
	if !tsNullish(state.head) { headers = tsInternalHTTP1FallbackGet(loop, state.head, "headers") }
	if headers.kind == tsArrayKind {
		values := (*tsArray)(headers.ref).values
		for i := 0; i+1 < len(values); i += 2 {
			name, value := tsText(values[i]), values[i+1]
			if name == "\x00" { if tsText(value) == "2" { state.noBody = true } else { state.closeDelimited = true }; continue }
			switch strings.ToLower(name) {
			case "content-length": hasContentLength = true
			case "transfer-encoding": hasTransferEncoding = true; if strings.Contains(strings.ToLower(tsInternalHTTP1FallbackString(loop, value)), "chunked") { state.chunked = true }
			case "date": hasDate = true
			case "connection": hasConnection = true
			case "keep-alive": hasKeepAlive = true
			}
			output += name + ": " + tsInternalHTTP1FallbackString(loop, value) + "\r\n"
		}
	}
	autoBits := 0
	if !tsNullish(state.head) { autoBits = int(tsNumber(tsInternalHTTP1FallbackGet(loop, state.head, "autoHeaderBits"))) }
	chunkedFromAutoBits := autoBits&16 != 0
	autoContentLength := tsNull
	autoChunked := false
	if !hasContentLength && !hasTransferEncoding && !state.noBody && !state.closeDelimited {
		if chunkedFromAutoBits || tsIsNull(contentLength) { state.chunked = true; autoChunked = !chunkedFromAutoBits } else { autoContentLength = contentLength }
	}
	if tsNullish(state.head) {
		if !hasDate { output += "Date: " + tsInternalHTTP1FallbackDate() + "\r\n" }
		if !hasConnection && !state.closeDelimited {
			if state.shouldKeepAlive { output += "Connection: keep-alive\r\n"; if !hasKeepAlive { timeout := state.keepAliveTimeout; if timeout == 0 { timeout = 5000 }; output += "Keep-Alive: timeout=" + strconv.Itoa(int(math.Floor(timeout/1000))) + "\r\n" } } else { output += "Connection: close\r\n" }
		}
	} else {
		if !hasDate && autoBits&1 != 0 { output += "Date: " + tsInternalHTTP1FallbackDate() + "\r\n" }
		if !hasConnection { if autoBits&2 != 0 { output += "Connection: keep-alive\r\n"; if !hasKeepAlive && autoBits&8 != 0 { output += "Keep-Alive: timeout=" + strconv.Itoa(int(tsNumber(tsInternalHTTP1FallbackGet(loop, state.head, "keepAliveTimeoutSecs")))) + "\r\n" } } else if autoBits&4 != 0 { output += "Connection: close\r\n" } }
	}
	if !tsIsNull(autoContentLength) { output += "Content-Length: " + tsInternalHTTP1FallbackString(loop, autoContentLength) + "\r\n" } else if autoChunked || chunkedFromAutoBits && state.chunked { output += "Transfer-Encoding: chunked\r\n" }
	output += "\r\n"
	tsInternalHTTP1FallbackSocketWrite(loop, state, tsInternalHTTP1FallbackKey(output), tsU)
}
func tsInternalHTTP1FallbackCreateResponseHandle(loop *tsLoop, socket, shouldKeepAlive, keepAliveTimeout tsValue) tsValue {
	state := &tsHTTP1FallbackResponseState{socket: socket, head: tsNull, shouldKeepAlive: tsTruthy(shouldKeepAlive), keepAliveTimeout: tsNumber(keepAliveTimeout), onWritable: tsNull, onWritableFrame: tsU, onFinished: tsNull}
	object := tsNewObject(); object.prototype = tsNull; object.descriptors = map[string]*tsDescriptor{}; object.set("onflushed", tsNull); result := tsObjectValue(object); state.handle = result
	method := func(name string, length int) { object.set(name, tsNamedNativeMethod(name, length, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue { return tsInternalHTTP1FallbackHandleMethod(loop, state, name, args) })) }
	for _, name := range []string{"socketDrained", "socketEnded", "writeContinue", "flushHeaders", "flushed", "abort"} { method(name, 1) }
	for _, name := range []string{"cork", "writeInformational", "writeHead", "writeHeadAndEnd", "write", "end"} { method(name, 4) }
	define := func(name string, getter func(*tsLoop) tsValue, setter func(*tsLoop, tsValue)) { object.descriptors[name] = &tsDescriptor{accessor:true, configurable:true, enumerable:true, getter:tsNamedNativeMethod("get "+name,0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{return getter(loop)}), setter:tsNamedNativeMethod("set "+name,1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{setter(loop,tsArg(args,0));return tsU})} }
	define("flags",func(loop *tsLoop)tsValue{if state.aborted||tsTruthy(tsInternalHTTP1FallbackGet(loop,state.socket,"destroyed")){return tsNumberValue(1)};return tsNumberValue(0)},func(*tsLoop,tsValue){})
	define("ended",func(*tsLoop)tsValue{return tsBooleanValue(state.ended)},func(_ *tsLoop,v tsValue){state.ended=tsTruthy(v)})
	define("finished",func(*tsLoop)tsValue{return tsBooleanValue(state.finished)},func(_ *tsLoop,v tsValue){state.finished=tsTruthy(v)})
	define("aborted",func(*tsLoop)tsValue{return tsBooleanValue(state.aborted)},func(_ *tsLoop,v tsValue){state.aborted=tsTruthy(v)})
	define("shouldKeepAlive",func(*tsLoop)tsValue{return tsBooleanValue(state.shouldKeepAlive)},func(_ *tsLoop,v tsValue){state.shouldKeepAlive=tsTruthy(v)})
	define("onfinished",func(*tsLoop)tsValue{return state.onFinished},func(_ *tsLoop,v tsValue){state.onFinished=v})
	define("onwritable",func(*tsLoop)tsValue{if tsNullish(state.onWritable){return tsU};return state.onWritable},func(loop *tsLoop,v tsValue){state.onWritable=v;if tsTruthy(v){state.onWritableFrame=tsCall(loop,tsGet(loop,tsInternalAsyncContextFrameModule(loop),tsInternalHTTP1FallbackKey("current")))}else{state.onWritableFrame=tsU}})
	define("bufferedAmount",func(loop *tsLoop)tsValue{if tsTruthy(state.onWritable)||state.ended&&!state.finished{return tsInternalHTTP1FallbackGet(loop,state.socket,"writableLength")};return tsNumberValue(0)},func(*tsLoop,tsValue){})
	return result
}
func tsInternalHTTP1FallbackHandleMethod(loop *tsLoop, state *tsHTTP1FallbackResponseState, name string, args []tsValue) tsValue {
	arg := func(index int) tsValue { return tsArg(args, index) }
	switch name {
	case "socketDrained": if tsTruthy(state.onWritable) { tsInternalHTTP1FallbackSettleWritable(loop,state,state.onWritable) }; return tsU
	case "socketEnded": if tsTruthy(state.onWritable) { tsInternalHTTP1FallbackSettleWritable(loop,state,arg(0)) }; return tsU
	case "cork": return tsCall(loop,arg(0))
	case "writeContinue": tsInternalHTTP1FallbackSocketWrite(loop,state,tsInternalHTTP1FallbackKey("HTTP/1.1 100 Continue\r\n\r\n"),tsU); return tsU
	case "writeInformational": if !tsTruthy(tsInternalHTTP1FallbackGet(loop,state.socket,"writableEnded")) { tsNodeHTTPServerCall(loop,state.socket,"write",arg(0),arg(1)) }; return tsU
	case "writeHead":
		status := int(tsNumber(arg(0))); if status < 100 || status > 999 { panic(tsThrown{tsFSException("RangeError","ERR_HTTP_INVALID_STATUS_CODE","Invalid status code: "+tsInternalHTTP1FallbackString(loop,arg(0)))}) }
		message:=arg(1); common:=tsNodeHTTPCommonModule(loop);invalid:=tsNodeHTTPCommonGet(loop,common,"_checkInvalidHeaderChar");if message.kind==tsStringKind&&tsTruthy(tsCall(loop,invalid,message)){panic(tsThrown{tsFSException("TypeError","ERR_INVALID_CHAR","Invalid character in statusMessage")})}
		head:=tsNewObject();head.set("statusCode",tsNumberValue(float64(status)));head.set("statusMessage",message);head.set("headers",arg(2));head.set("autoHeaderBits",arg(3));head.set("keepAliveTimeoutSecs",arg(4));state.head=tsObjectValue(head);return tsU
	case "flushHeaders": tsInternalHTTP1FallbackWriteHead(loop,state,tsNull);return tsU
	case "writeHeadAndEnd": tsInternalHTTP1FallbackHandleMethod(loop,state,"writeHead",[]tsValue{arg(0),arg(1),arg(2),arg(6),arg(7)});return tsInternalHTTP1FallbackHandleMethod(loop,state,"end",[]tsValue{arg(3),arg(4),tsU,arg(5)})
	case "write":
		buffer:=tsInternalHTTP1FallbackToBuffer(loop,arg(0),arg(1));length:=tsInternalHTTP1FallbackCountBodyBytes(loop,state,buffer,arg(3),false);tsInternalHTTP1FallbackWriteHead(loop,state,tsNull);tsInternalHTTP1FallbackWriteBody(loop,state,buffer,tsU)
		if tsIsNull(buffer)||!tsTruthy(tsInternalHTTP1FallbackGet(loop,state.socket,"writableNeedDrain")){return tsNumberValue(float64(length))}
		callback:=arg(2);if tsTruthy(callback){state.onWritable=callback;state.onWritableFrame=tsCall(loop,tsGet(loop,tsInternalAsyncContextFrameModule(loop),tsInternalHTTP1FallbackKey("current")))};if length>0{return tsNumberValue(float64(-length))};return tsNumberValue(-1)
	case "end":
		if state.ended{return tsNumberValue(0)};state.onWritable=tsNull;state.onWritableFrame=tsU;buffer:=tsInternalHTTP1FallbackToBuffer(loop,arg(0),arg(1));length:=tsInternalHTTP1FallbackCountBodyBytes(loop,state,buffer,arg(3),true);tsInternalHTTP1FallbackWriteHead(loop,state,tsNumberValue(float64(length)));terminated:=state.chunked&&!state.noBody
		endCallback:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{tsInternalHTTP1FallbackEndWritten(loop,state,tsArg(args,0));return tsU}))
		tsInternalHTTP1FallbackWriteBody(loop,state,buffer,func()tsValue{if terminated{return tsU};return endCallback}());if terminated{tsInternalHTTP1FallbackSocketWrite(loop,state,tsInternalHTTP1FallbackKey("0\r\n\r\n"),endCallback)};state.ended=true
		writableLength:=tsNumber(tsInternalHTTP1FallbackGet(loop,state.socket,"writableLength"));if writableLength>0&&!tsTruthy(tsInternalHTTP1FallbackGet(loop,state.socket,"destroyed")){if tsTruthy(tsInternalHTTP1FallbackGet(loop,state.socket,"writableEnded")){tsNodeHTTPServerCall(loop,state.socket,"once",tsInternalHTTP1FallbackKey("finish"),endCallback)}else if length==0&&!terminated{tsNodeHTTPServerCall(loop,state.socket,"write",tsInternalHTTP1FallbackKey(""),tsInternalHTTP1FallbackKey("latin1"),endCallback)}}else{state.finished=true}
		finishedCallback:=state.onFinished;state.onFinished=tsNull;if tsTruthy(finishedCallback){tsInternalHTTP1FallbackSafeCallback(loop,finishedCallback)};if state.closeDelimited&&!tsTruthy(tsInternalHTTP1FallbackGet(loop,state.socket,"destroyed")){tsNodeHTTPServerCall(loop,state.socket,"end")};if state.finished{return tsNumberValue(float64(length))};return tsNumberValue(float64(-(length+1)))
	case "flushed": tsInternalHTTP1FallbackFlush(loop,state);return tsU
	case "abort": state.aborted=true;if !tsTruthy(tsInternalHTTP1FallbackGet(loop,state.socket,"destroyed")){tsNodeHTTPServerCall(loop,state.socket,"destroy")};return tsU
	}
	return tsU
}

type tsHTTP1FallbackConnection struct {
	loop *tsLoop
	server, socket, options, parser, common, http, pipeline tsValue
	req, pendingUpgrade tsValue
	flushedHeaders []tsValue
	flushedURL string
	requestsCount int
	maxRequestsPerSocket tsValue
	keepAliveTimeout tsValue
	connections tsValue
}

func tsInternalHTTP1FallbackPipeline(loop *tsLoop) tsValue {
	internal := tsInternalHTTPModule(loop)
	pipeline := tsInternalHTTPGet(loop, internal, "http1ServerPipeline")
	if pipeline.kind != tsObjectKind { pipeline = tsObjectValue(tsNewObject()); tsInternalHTTPSet(loop, internal, "http1ServerPipeline", pipeline) }
	tsInternalHTTP1FallbackInit(loop)
	install := func(name string, fn tsValue) { if tsIsUndefined(tsGet(loop, pipeline, tsInternalHTTP1FallbackKey(name))) { tsSet(loop, pipeline, tsInternalHTTP1FallbackKey(name), fn) } }
	install("kMustCloseConnection", tsInternalHTTP1FallbackSingleton.kMustCloseConnection)
	install("queuePipelinedResponse", tsNamedNativeMethod("queuePipelinedResponse", 3, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
		socket := tsArg(args, 0); var response tsValue = tsArg(args, 1); queue := tsGet(loop, socket, tsInternalHTTP1FallbackSingleton.kQueue)
		if queue.kind != tsArrayKind { queue = tsArrayValue(&tsArray{}); tsSet(loop, socket, tsInternalHTTP1FallbackSingleton.kQueue, queue) }
		state := tsNewObject(); state.set("socket", socket); state.set("isAncient", tsArg(args, 2)); state.set("ops", tsArrayValue(&tsArray{})); state.set("bytes",tsNumberValue(0));state.set("ended",tsBooleanValue(false));state.set("needDrain",tsBooleanValue(false));tsSet(loop, response, tsInternalHTTP1FallbackSingleton.kQueueState, tsObjectValue(state))
		a := (*tsArray)(queue.ref); a.values = append(a.values, response); return tsU
	}))
	install("lastPipelinedResponse", tsNamedNativeMethod("lastPipelinedResponse", 1, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
		queue := tsGet(loop, tsArg(args, 0), tsInternalHTTP1FallbackSingleton.kQueue); if queue.kind != tsArrayKind || len((*tsArray)(queue.ref).values) == 0 { return tsU }; a := (*tsArray)(queue.ref); return a.values[len(a.values)-1]
	}))
	install("advanceResponsePipeline", tsNamedNativeMethod("advanceResponsePipeline", 2, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
		server, socket := tsArg(args, 0), tsArg(args, 1); if tsNullish(socket) || tsTruthy(tsInternalHTTP1FallbackGet(loop, socket, "destroyed")) || tsTruthy(tsInternalHTTP1FallbackGet(loop, socket, "writableEnded")) || !tsNullish(tsInternalHTTP1FallbackGet(loop, socket, "_httpMessage")) { return tsU }
		queue := tsGet(loop, socket, tsInternalHTTP1FallbackSingleton.kQueue); if queue.kind != tsArrayKind { return tsU }; a := (*tsArray)(queue.ref); if len(a.values) == 0 { return tsU }
		var response tsValue = a.values[0]; a.values = a.values[1:]; tsSet(loop, socket, tsInternalHTTP1FallbackSingleton.kQueue, tsArrayValue(a))
		queueState:=tsGet(loop,response,tsInternalHTTP1FallbackSingleton.kQueueState);ops:=tsInternalHTTP1FallbackGet(loop,queueState,"ops");ended:=tsTruthy(tsInternalHTTP1FallbackGet(loop,queueState,"ended"));queuedBytes:=tsNumber(tsInternalHTTP1FallbackGet(loop,queueState,"bytes"));remainingBytes:=tsNumber(tsInternalHTTP1FallbackGet(loop,socket,"_httpQueuedBytes"))-queuedBytes;if remainingBytes<0{remainingBytes=0};tsInternalHTTP1FallbackSet(loop,socket,"_httpQueuedBytes",tsNumberValue(remainingBytes));tsInternalHTTP1FallbackSet(loop,response,"outputSize",tsNumberValue(0))
		tsSet(loop,response,tsInternalHTTP1FallbackSingleton.kQueueState,tsU)
		assign := tsInternalHTTP1FallbackGet(loop, response, "assignSocket"); if assign.kind == tsFunctionKind { tsCallReceiver(loop, assign, response, socket) }
		if ended { tsInternalHTTP1FallbackSet(loop,response,"finished",tsBooleanValue(false)) }
		if ops.kind==tsArrayKind { opArray:=(*tsArray)(ops.ref);opList:=make([]tsValue,len(opArray.values));copy(opList,opArray.values);for _,op:=range opList { kind:=tsInternalHTTP1FallbackGet(loop,op,"0");chunk:=tsInternalHTTP1FallbackGet(loop,op,"1");encoding:=tsInternalHTTP1FallbackGet(loop,op,"2");callback:=tsInternalHTTP1FallbackGet(loop,op,"3");if tsText(kind)=="end"{tsNodeHTTPServerCall(loop,response,"end",chunk,encoding,callback)}else{tsNodeHTTPServerCall(loop,response,"write",chunk,encoding,callback)} } }
		flush := tsInternalHTTP1FallbackGet(loop, response, "_flush"); if flush.kind == tsFunctionKind { tsCallReceiver(loop, flush, response) }
		_ = server; return tsU
	}))
	install("abortQueuedPipelinedResponses", tsNamedNativeMethod("abortQueuedPipelinedResponses", 2, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
		socket := tsArg(args, 0); queue := tsGet(loop, socket, tsInternalHTTP1FallbackSingleton.kQueue); if queue.kind != tsArrayKind { return tsU };queueArray:=(*tsArray)(queue.ref);items:=make([]tsValue,len(queueArray.values));copy(items,queueArray.values);queueArray.resize(0);err := tsArg(args, 1);tsInternalHTTP1FallbackSet(loop,socket,"_httpQueuedBytes",tsNumberValue(0))
		for _, response := range items { request := tsInternalHTTP1FallbackGet(loop, response, "req"); state:=tsGet(loop,response,tsInternalHTTP1FallbackSingleton.kQueueState);ops:=tsInternalHTTP1FallbackGet(loop,state,"ops");if ops.kind==tsArrayKind{for _,op:=range (*tsArray)(ops.ref).values{callback:=tsInternalHTTP1FallbackGet(loop,op,"3");if callback.kind==tsFunctionKind{tsInternalHTTP1FallbackSafeCallback(loop,callback,err)}}}; if !tsTruthy(tsInternalHTTP1FallbackGet(loop, request, "destroyed")) { if tsTruthy(err) { tsNodeHTTPIncomingCall(loop, request, "destroy", err) } else { tsNodeHTTPIncomingCall(loop, request, "destroy") } }; if !tsTruthy(tsInternalHTTP1FallbackGet(loop, response, "destroyed")) { tsNodeHTTPIncomingCall(loop, response, "destroy") } }
		return tsU
	}))
	install("maybePauseFallbackReads", tsNamedNativeMethod("maybePauseFallbackReads", 1, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
		socket := tsArg(args, 0); if tsTruthy(tsInternalHTTP1FallbackGet(loop, socket, "_paused")) { return tsU }; state := tsInternalHTTP1FallbackGet(loop, socket, "_writableState");pending:=tsNumber(tsInternalHTTP1FallbackGet(loop,socket,"writableLength"))+tsNumber(tsInternalHTTP1FallbackGet(loop,socket,"_httpQueuedBytes")); if tsTruthy(tsInternalHTTP1FallbackGet(loop, state, "needDrain")) || pending >= tsNumber(tsInternalHTTP1FallbackGet(loop, socket, "writableHighWaterMark")) { tsInternalHTTP1FallbackSet(loop, socket, "_paused", tsBooleanValue(true)); tsNodeHTTPServerCall(loop, socket, "pause") }; return tsU
	}))
	install("resumeFallbackReadsOnDrain", tsNamedNativeMethod("resumeFallbackReadsOnDrain", 1, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
		socket := tsArg(args, 0); pending:=tsNumber(tsInternalHTTP1FallbackGet(loop,socket,"writableLength"))+tsNumber(tsInternalHTTP1FallbackGet(loop,socket,"_httpQueuedBytes"));if tsTruthy(tsInternalHTTP1FallbackGet(loop, socket, "_paused")) && pending <= tsNumber(tsInternalHTTP1FallbackGet(loop, socket, "writableHighWaterMark")) { tsInternalHTTP1FallbackSet(loop, socket, "_paused", tsBooleanValue(false)); tsNodeHTTPServerCall(loop, socket, "resume") }; return tsU
	}))
	install("finishDrainedResponse", tsNamedNativeMethod("finishDrainedResponse", 1, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
		var response tsValue = tsArg(args, 0); tsInternalHTTP1FallbackSet(loop, response, "writableFinished", tsBooleanValue(true)); if !tsTruthy(tsInternalHTTP1FallbackGet(loop, response, "finished")) { tsNodeHTTPServerEmit(loop, response, "finish") }; return tsU
	}))
	return pipeline
}

func tsInternalHTTP1FallbackParserSlot(loop *tsLoop, parserCtor tsValue, name string) tsValue { return tsGet(loop, parserCtor, tsInternalHTTP1FallbackKey(name)) }
func tsInternalHTTP1FallbackCallback(loop *tsLoop, state *tsHTTP1FallbackConnection, name string, length int) tsValue {
	return tsNamedNativeMethod(name, length, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue { return tsInternalHTTP1FallbackConnectionDispatch(loop, state, name, args) })
}
func tsInternalHTTP1FallbackConnectionListener(loop *tsLoop, server, socket, options tsValue) tsValue {
	state := &tsHTTP1FallbackConnection{loop:loop, server:server, socket:socket, options:options, req:tsNull, pendingUpgrade:tsNull, flushedHeaders:[]tsValue{}}
	http1Options := tsInternalHTTP1FallbackGet(loop, options, "http1Options")
	state.http = tsInternalHTTP1FallbackGet(loop, options, "httpModule")
	if tsNullish(state.http) {
		state.http = tsObjectValue(tsNewObject())
		incoming := tsInternalHTTP1FallbackGet(loop, http1Options, "IncomingMessage")
		if tsNullish(incoming) { incoming = tsNodeHTTPIncomingGet(loop, tsNodeHTTPIncomingModule(loop), "IncomingMessage") }
		response := tsInternalHTTP1FallbackGet(loop, http1Options, "ServerResponse")
		if tsNullish(response) { response = tsNodeHTTPServerAPIGet(loop, tsNodeHTTPServerAPIModule(loop), "ServerResponse") }
		tsSet(loop, state.http, tsInternalHTTP1FallbackKey("IncomingMessage"), incoming)
		tsSet(loop, state.http, tsInternalHTTP1FallbackKey("ServerResponse"), response)
	}
	state.common = tsNodeHTTPCommonModule(loop)
	state.pipeline = tsInternalHTTP1FallbackPipeline(loop)
	internal := tsInternalHTTPModule(loop); state.keepAliveTimeout = tsInternalHTTP1FallbackGet(loop, server, "keepAliveTimeout")
	if !tsIsNumeric(state.keepAliveTimeout) { state.keepAliveTimeout = tsNumberValue(5000) }
	state.connections = tsInternalHTTP1FallbackSym(loop, server, "kHttp1Connections")
	if state.connections.kind != tsCollectionKind { state.connections = tsNewCollection(true, tsU); tsInternalHTTP1FallbackPut(loop, server, "kHttp1Connections", state.connections) }
	add := tsGet(loop, state.connections, tsInternalHTTP1FallbackKey("add")); if add.kind == tsFunctionKind { tsCallReceiver(loop, add, state.connections, socket) }
	tsInternalHTTP1FallbackSet(loop, socket, "server", server); tsInternalHTTP1FallbackPut(loop, socket, "kHttp1ActiveRequests", tsNumberValue(0)); tsInternalHTTP1FallbackPut(loop, socket, "kHttp1MessageOpen", tsBooleanValue(false))
	parserCtor := tsNodeHTTPCommonGet(loop, state.common, "HTTPParser"); parser := tsConstruct(loop, parserCtor); state.parser = parser
	requestType := tsInternalHTTP1FallbackParserSlot(loop, parserCtor, "REQUEST"); maxSize := tsInternalHTTP1FallbackGet(loop, server, "maxHeaderSize"); if !tsTruthy(maxSize) { maxSize = tsNumberValue(0) }
	flagsFn := tsNodeHTTPCommonGet(loop, state.common, "calculateLenientFlags"); flags := tsCall(loop, flagsFn, tsNodeHTTPServerGet(loop, server, "httpValidation"), tsNodeHTTPServerGet(loop, server, "insecureHTTPParser"))
	tsNodeHTTPCommonCall(loop, parser, "initialize", requestType, tsObjectValue(tsNewObject()), maxSize, flags); tsNodeHTTPCommonSet(loop, parser, "socket", socket); tsInternalHTTP1FallbackSet(loop, socket, "parser", parser)
	maxHeaders := tsNodeHTTPServerGet(loop, server, "maxHeadersCount"); if tsIsNumeric(maxHeaders) { tsNodeHTTPServerSet(loop, parser, "maxHeaderPairs", tsNumberValue(tsNumber(maxHeaders)*2)) } else { tsNodeHTTPServerSet(loop, parser, "maxHeaderPairs", tsNodeHTTPCommonGet(loop, state.common, "MAX_HEADER_PAIRS")) }
	for _, entry := range []struct{slot, name string; length int}{{"kOnMessageBegin","onHttp1MessageBegin",0},{"kOnHeaders","onHttp1Headers",2},{"kOnHeadersComplete","onHttp1HeadersComplete",9},{"kOnBody","onHttp1Body",2},{"kOnMessageComplete","onHttp1MessageComplete",0}} { tsSet(loop, parser, tsInternalHTTP1FallbackParserSlot(loop, parserCtor, entry.slot), tsInternalHTTP1FallbackCallback(loop,state, entry.name, entry.length)) }
	for _, entry := range []struct{event,name string; once bool}{{"data","onHttp1SocketData",false},{"error","onHttp1SocketErrorListener",false},{"end","onHttp1SocketEnd",true},{"drain","onHttp1SocketDrain",false},{"finish","onHttp1SocketFinish",true},{"close","onHttp1SocketClose",true}} {
		method := "on"; if entry.once { method = "once" }; tsNodeHTTPServerCall(loop, socket, method, tsInternalHTTP1FallbackKey(entry.event), tsInternalHTTP1FallbackCallback(loop,state, entry.name, 1))
	}
	tsNodeHTTPServerCall(loop, socket, "prependOnceListener", tsInternalHTTP1FallbackKey("close"), tsInternalHTTP1FallbackCallback(loop,state,"onHttp1SocketCloseFlush",0))
	_ = internal
	return tsU
}

func tsInternalHTTP1FallbackConnectionDispatch(loop *tsLoop,state *tsHTTP1FallbackConnection, name string, args []tsValue) tsValue { arg := func(index int) tsValue { return tsArg(args,index) }; server,socket,parser := state.server,state.socket,state.parser
	switch name {
	case "onHttp1MessageBegin": tsInternalHTTP1FallbackPut(loop,socket,"kHttp1MessageOpen",tsBooleanValue(true)); return tsU
	case "onHttp1Headers":
		headers := tsNodeHTTPIncomingArray(arg(0)); state.flushedHeaders = append(state.flushedHeaders, headers.values...); state.flushedURL += tsText(arg(1)); return tsU
	case "onHttp1HeadersComplete":
		major,minor,headers,method,url,upgrade,keepAlive:=arg(0),arg(1),arg(2),arg(3),arg(4),arg(7),arg(8)
		if tsIsUndefined(headers) { headers=tsArrayValue(&tsArray{values:append([]tsValue(nil),state.flushedHeaders...)});state.flushedHeaders=nil }
		if tsIsUndefined(url) { url=tsInternalHTTP1FallbackKey(state.flushedURL);state.flushedURL="" }
		if tsNumber(major)==1&&tsNumber(minor)==1&&!tsStrictEqual(tsNodeHTTPServerGet(loop,server,"requireHostHeader"),tsBooleanValue(false)){hasHost:=false;if headers.kind==tsArrayKind{raw:=(*tsArray)(headers.ref).values;for index:=0;index+1<len(raw);index+=2{if strings.EqualFold(tsText(raw[index]),"host"){hasHost=true;break}}};if !hasHost{if tsTruthy(tsNodeHTTPServerGet(loop,socket,"writable")){line:="HTTP/1.1 400 Bad Request\r\nConnection: close\r\nDate: "+tsInternalHTTP1FallbackDate()+"\r\nTransfer-Encoding: chunked\r\n\r\n0\r\n\r\n";tsNodeHTTPServerCall(loop,socket,"write",tsInternalHTTP1FallbackKey(line));tsNodeHTTPServerCall(loop,socket,"destroySoon")};return tsNumberValue(0)}}
		active:=tsNumber(tsInternalHTTP1FallbackSym(loop,socket,"kHttp1ActiveRequests"));tsInternalHTTP1FallbackPut(loop,socket,"kHttp1ActiveRequests",tsNumberValue(active+1))
		incomingCtor:=tsNodeHTTPServerGet(loop,state.http,"IncomingMessage");http1Options:=tsInternalHTTP1FallbackGet(loop,state.options,"http1Options");if custom:=tsInternalHTTP1FallbackGet(loop,http1Options,"IncomingMessage");!tsNullish(custom){incomingCtor=custom};request:=tsConstruct(loop,incomingCtor,socket);state.req=request;tsNodeHTTPServerPut(loop,socket,"kRequest",request)
		for key,value:=range map[string]tsValue{"socket":socket,"httpVersionMajor":major,"httpVersionMinor":minor,"httpVersion":tsInternalHTTP1FallbackKey(tsText(major)+"."+tsText(minor)),"url":url,"upgrade":upgrade}{tsNodeHTTPIncomingSet(loop,request,key,value)}
		requestKeys:=tsNodeHTTPIncomingModule(loop);reqKeepAlive:=tsNodeHTTPIncomingGet(loop,requestKeys,"kReqShouldKeepAlive");if !tsIsUndefined(reqKeepAlive){tsSet(loop,request,reqKeepAlive,keepAlive)}
		methods:=tsNodeHTTPCommonGet(loop,state.common,"methods");if tsIsNumeric(method){index:=int(tsNumber(method));if methods.kind==tsArrayKind&&index>=0&&index<len((*tsArray)(methods.ref).values){tsNodeHTTPIncomingSet(loop,request,"method",(*tsArray)(methods.ref).values[index])}}else{tsNodeHTTPIncomingSet(loop,request,"method",method)}
		if tsTruthy(tsNodeHTTPServerGet(loop,server,"joinDuplicateHeaders")){tsNodeHTTPIncomingSet(loop,request,"joinDuplicateHeaders",tsBooleanValue(true))};tsNodeHTTPIncomingCall(loop,request,"_addHeaderLines",headers,tsNumberValue(float64(tsNodeHTTPIncomingArray(headers).length())))
		if tsTruthy(upgrade){shouldUpgrade:=tsText(tsNodeHTTPIncomingGet(loop,request,"method"))=="CONNECT";callback:=tsNodeHTTPServerGet(loop,server,"shouldUpgradeCallback");if callback.kind==tsFunctionKind{shouldUpgrade=tsTruthy(tsCall(loop,callback,request))}else{shouldUpgrade=shouldUpgrade||tsNumber(tsNodeHTTPServerCall(loop,server,"listenerCount",tsInternalHTTP1FallbackKey("upgrade")))>0};tsNodeHTTPIncomingSet(loop,request,"upgrade",tsBooleanValue(shouldUpgrade));if shouldUpgrade{state.pendingUpgrade=request;return tsNumberValue(2)}}
		readFn:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{tsNodeHTTPIncomingSet(loop,request,"_consuming",tsBooleanValue(true));if !tsTruthy(tsInternalHTTP1FallbackGet(loop,socket,"_paused"))&&tsTruthy(tsInternalHTTP1FallbackGet(loop,socket,"readable")){tsNodeHTTPServerCall(loop,socket,"resume")};return tsU}));tsNodeHTTPIncomingSet(loop,request,"_read",readFn)
		responseCtor:=tsNodeHTTPServerGet(loop,state.http,"ServerResponse");if custom:=tsInternalHTTP1FallbackGet(loop,http1Options,"ServerResponse");!tsNullish(custom){responseCtor=custom};response:=tsConstruct(loop,responseCtor,request);tsNodeHTTPServerSet(loop,response,"_keepAliveTimeout",state.keepAliveTimeout);maxPerSocket:=tsNodeHTTPServerGet(loop,server,"maxRequestsPerSocket");tsNodeHTTPServerSet(loop,response,"_maxRequestsPerSocket",maxPerSocket)
		handle:=tsInternalHTTP1FallbackCreateResponseHandle(loop,socket,keepAlive,state.keepAliveTimeout)
		finishedFn:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{n:=tsNumber(tsInternalHTTP1FallbackSym(loop,socket,"kHttp1ActiveRequests"))-1;if n<0{n=0};tsInternalHTTP1FallbackPut(loop,socket,"kHttp1ActiveRequests",tsNumberValue(n));if !tsTruthy(keepAlive)&&!tsTruthy(tsInternalHTTP1FallbackGet(loop,socket,"destroyed")){tsNodeHTTPServerCall(loop,socket,"end")};return tsU}));tsNodeHTTPServerSet(loop,handle,"onfinished",finishedFn)
		flushedFn:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{tsInternalHTTP1FallbackCallPipeline(loop,state.pipeline,"finishDrainedResponse",response);return tsU}));tsNodeHTTPServerSet(loop,handle,"onflushed",flushedFn);tsSet(loop,response,tsInternalHTTPGet(loop,tsInternalHTTPModule(loop),"kHandle"),handle)
		if !tsNullish(tsInternalHTTP1FallbackGet(loop,socket,"_httpMessage")){tsInternalHTTP1FallbackCallPipeline(loop,state.pipeline,"queuePipelinedResponse",socket,response,tsBooleanValue(tsNumber(major)<1||tsNumber(minor)<1))}else{tsNodeHTTPServerCall(loop,response,"assignSocket",socket)}
		finishListener:=tsFunctionValue(tsFunc(func(loop *tsLoop,values ...tsValue)tsValue{responseHandle:=tsGet(loop,response,tsInternalHTTPGet(loop,tsInternalHTTPModule(loop),"kHandle"));tsNodeHTTPServerSet(loop,responseHandle,"ended",tsBooleanValue(true));tsNodeHTTPServerCall(loop,responseHandle,"flushed");tsNodeHTTPServerCall(loop,responseHandle,"onfinished");tsNodeHTTPServerResponseCloseOnFinish(loop,response,socket);finishedRequest:=tsNodeHTTPServerGet(loop,response,"req");if !tsTruthy(tsNodeHTTPIncomingGet(loop,finishedRequest,"_consuming"))&&!tsTruthy(tsNodeHTTPServerGet(loop,tsNodeHTTPServerGet(loop,finishedRequest,"_readableState"),"resumeScheduled")){tsNodeHTTPIncomingCall(loop,finishedRequest,"_dump")};tsNodeHTTPServerCall(loop,response,"detachSocket",socket);closeKey:=tsInternalHTTP1FallbackGet(loop,state.pipeline,"kMustCloseConnection");if tsTruthy(tsGet(loop,response,closeKey))||tsTruthy(tsNodeHTTPServerGet(loop,response,"_last")){destroySoon:=tsNodeHTTPServerGet(loop,socket,"destroySoon");if destroySoon.kind==tsFunctionKind{tsCallReceiver(loop,destroySoon,socket)}else if !tsTruthy(tsNodeHTTPServerGet(loop,socket,"writableEnded")){tsNodeHTTPServerCall(loop,socket,"end")};return tsU};tsInternalHTTP1FallbackCallPipeline(loop,state.pipeline,"advanceResponsePipeline",server,socket);return tsU}))
		tsNodeHTTPServerCall(loop,response,"on",tsInternalHTTP1FallbackKey("finish"),finishListener);tsInternalHTTP1FallbackCallPipeline(loop,state.pipeline,"maybePauseFallbackReads",socket)
		if tsNumber(major)==1&&tsNumber(minor)==1&&tsIsNumeric(maxPerSocket)&&tsNumber(maxPerSocket)>0{state.requestsCount++;tsNodeHTTPServerSet(loop,response,"maxRequestsOnConnectionReached",tsBooleanValue(tsNumber(maxPerSocket)<=float64(state.requestsCount)));if tsNumber(maxPerSocket)<float64(state.requestsCount){tsNodeHTTPServerEmit(loop,server,"dropRequest",request,socket);tsNodeHTTPServerCall(loop,response,"writeHead",tsNumberValue(503));tsNodeHTTPServerCall(loop,response,"end");return tsNumberValue(0)}}
		expect:=tsNodeHTTPIncomingGet(loop,request,"headers");expect=tsGet(loop,expect,tsInternalHTTP1FallbackKey("expect"));if !tsIsUndefined(expect)&&tsNumber(major)==1&&tsNumber(minor)==1{expression:=tsNodeHTTPCommonGet(loop,state.common,"continueExpression");matches:=tsTruthy(tsNodeHTTPServerCall(loop,expression,"test",tsInternalHTTP1FallbackKey(tsText(expect))));if matches{if tsNumber(tsNodeHTTPServerCall(loop,server,"listenerCount",tsInternalHTTP1FallbackKey("checkContinue")))>0{tsNodeHTTPServerEmit(loop,server,"checkContinue",request,response)}else{tsNodeHTTPServerCall(loop,response,"writeContinue");tsNodeHTTPServerEmit(loop,server,"request",request,response)}}else if tsNumber(tsNodeHTTPServerCall(loop,server,"listenerCount",tsInternalHTTP1FallbackKey("checkExpectation")))>0{tsNodeHTTPServerEmit(loop,server,"checkExpectation",request,response)}else{tsNodeHTTPServerCall(loop,response,"writeHead",tsNumberValue(417));tsNodeHTTPServerCall(loop,response,"end")};return tsNumberValue(0)}
		tsNodeHTTPServerEmit(loop,server,"request",request,response);return tsNumberValue(0)
	case "onHttp1Body": if !tsNullish(state.req)&&!tsTruthy(tsNodeHTTPIncomingGet(loop,state.req,"_dumped")){tsNodeHTTPIncomingCall(loop,state.req,"push",arg(0))};return tsU
	case "onHttp1MessageComplete":
		tsInternalHTTP1FallbackPut(loop,socket,"kHttp1MessageOpen",tsBooleanValue(false));tsNodeHTTPServerPut(loop,socket,"kRequest",tsU);trailers:=state.flushedHeaders;if len(trailers)>0{state.flushedHeaders=nil};state.flushedURL="";if !tsNullish(state.req){tsNodeHTTPIncomingSet(loop,state.req,"complete",tsBooleanValue(true));tsNodeHTTPIncomingCall(loop,state.req,"_addHeaderLines",tsArrayValue(&tsArray{values:trailers}),tsNumberValue(float64(len(trailers))));tsNodeHTTPIncomingCall(loop,state.req,"push",tsNull)};return tsU
	case "onHttp1SocketError", "onHttp1SocketErrorListener":
		err:=arg(0);common:=state.common;prepare:=tsNodeHTTPCommonGet(loop,common,"prepareError");if prepare.kind==tsFunctionKind{tsCall(loop,prepare,err,parser,arg(1))};if !tsTruthy(tsNodeHTTPServerEmit(loop,server,"clientError",err,socket)){if tsTruthy(tsNodeHTTPServerGet(loop,socket,"writable"))&&!tsTruthy(tsNodeHTTPServerGet(loop,socket,"destroyed")){code:=tsText(tsNodeHTTPServerGet(loop,err,"code"));message:="HTTP/1.1 400 Bad Request\r\nConnection: close\r\n\r\n";if code=="HPE_HEADER_OVERFLOW"{message="HTTP/1.1 431 Request Header Fields Too Large\r\nConnection: close\r\n\r\n"}else if code=="HPE_CHUNK_EXTENSIONS_OVERFLOW"{message="HTTP/1.1 413 Payload Too Large\r\nConnection: close\r\n\r\n"};tsNodeHTTPServerCall(loop,socket,"write",tsInternalHTTP1FallbackKey(message),tsInternalHTTP1FallbackKey("latin1"))};tsNodeHTTPServerCall(loop,socket,"destroy",err)};return tsU
	case "onHttp1SocketData":
		data:=arg(0);ret:=tsNodeHTTPCommonCall(loop,parser,"execute",data);if ret.kind==tsErrorKind{tsInternalHTTP1FallbackConnectionDispatch(loop,state,"onHttp1SocketError",[]tsValue{ret,data});return tsU};if !tsNullish(state.pendingUpgrade){request:=state.pendingUpgrade;state.pendingUpgrade=tsNull;for _,event:=range []struct{name,callback string}{{"data","onHttp1SocketData"},{"error","onHttp1SocketErrorListener"},{"end","onHttp1SocketEnd"},{"drain","onHttp1SocketDrain"}}{tsNodeHTTPServerCall(loop,socket,"removeListener",tsInternalHTTP1FallbackKey(event.name),tsInternalHTTP1FallbackCallback(loop,state,event.callback,1))};deleteFn:=tsGet(loop,state.connections,tsInternalHTTP1FallbackKey("delete"));if deleteFn.kind==tsFunctionKind{tsCallReceiver(loop,deleteFn,state.connections,socket)};tsNodeHTTPServerPut(loop,socket,"kHandedOff",tsBooleanValue(true));tracked:=tsNodeHTTPServerSym(loop,state.server,"kTrackedConnections");removeTracked:=tsGet(loop,tracked,tsInternalHTTP1FallbackKey("delete"));if removeTracked.kind==tsFunctionKind{tsCallReceiver(loop,removeTracked,tracked,socket)};tsNodeHTTPCommonCall(loop,parser,"close");tsInternalHTTP1FallbackSet(loop,socket,"parser",tsNull);eventName:="upgrade";if tsText(tsNodeHTTPIncomingGet(loop,request,"method"))=="CONNECT"{eventName="connect"};raw:=tsNodeBufferVisibleBytes(data);parsed:=int(tsNumber(ret));if parsed<0{parsed=0};if parsed>len(raw){parsed=len(raw)};body:=tsNodeBuffer(append([]byte(nil),raw[parsed:]...));if tsNumber(tsNodeHTTPServerCall(loop,server,"listenerCount",tsInternalHTTP1FallbackKey(eventName)))>0{tsNodeHTTPServerSet(loop,socket,"readableFlowing",tsNull);tsNodeHTTPServerEmit(loop,server,eventName,request,socket,body)}else{tsNodeHTTPServerCall(loop,socket,"destroy")}};return tsU
	case "onHttp1SocketEnd":
		result:=tsNodeHTTPCommonCall(loop,parser,"finish");if result.kind==tsErrorKind{tsInternalHTTP1FallbackConnectionDispatch(loop,state,"onHttp1SocketError",[]tsValue{result,tsU});return tsU};if !tsTruthy(tsNodeHTTPServerGet(loop,server,"httpAllowHalfOpen")){if !tsNullish(state.req)&&!tsTruthy(tsNodeHTTPIncomingGet(loop,state.req,"complete")){tsNodeHTTPIncomingCall(loop,state.req,"destroy")};if tsTruthy(tsNodeHTTPServerGet(loop,socket,"writable")){tsNodeHTTPServerCall(loop,socket,"end")};return tsU};last:=tsInternalHTTP1FallbackCallPipeline(loop,state.pipeline,"lastPipelinedResponse",socket);message:=last;if tsNullish(message){message=tsNodeHTTPServerGet(loop,socket,"_httpMessage")};if !tsNullish(message){tsNodeHTTPServerSet(loop,message,"_last",tsBooleanValue(true))}else if tsTruthy(tsNodeHTTPServerGet(loop,socket,"writable")){tsNodeHTTPServerCall(loop,socket,"end")};return tsU
	case "onHttp1SocketDrain":
		tsInternalHTTP1FallbackCallPipeline(loop,state.pipeline,"resumeFallbackReadsOnDrain",socket);message:=tsNodeHTTPServerGet(loop,socket,"_httpMessage");if !tsNullish(message){handle:=tsGet(loop,message,tsInternalHTTPGet(loop,tsInternalHTTPModule(loop),"kHandle"));tsNodeHTTPServerCall(loop,handle,"socketDrained")};return tsU
	case "onHttp1SocketFinish":
		message:=tsNodeHTTPServerGet(loop,socket,"_httpMessage");if !tsNullish(message){handle:=tsGet(loop,message,tsInternalHTTPGet(loop,tsInternalHTTPModule(loop),"kHandle"));callback:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{tsNodeHTTPServerCall(loop,message,"_callPendingCallbacks");return tsU}));tsNodeHTTPServerCall(loop,handle,"socketEnded",callback)};return tsU
	case "onHttp1SocketCloseFlush": message:=tsNodeHTTPServerGet(loop,socket,"_httpMessage");if !tsNullish(message){tsNodeHTTPServerCall(loop,tsGet(loop,message,tsInternalHTTPGet(loop,tsInternalHTTPModule(loop),"kHandle")),"flushed")};return tsU
	case "onHttp1SocketClose":
		deleteFn:=tsGet(loop,state.connections,tsInternalHTTP1FallbackKey("delete"));if deleteFn.kind==tsFunctionKind{tsCallReceiver(loop,deleteFn,state.connections,socket)};message:=tsNodeHTTPServerGet(loop,socket,"_httpMessage");if !tsNullish(message){handle:=tsGet(loop,message,tsInternalHTTPGet(loop,tsInternalHTTPModule(loop),"kHandle"));err:=tsNodeHTTPServerGet(loop,socket,"errored");if tsNullish(err){err=tsNodeHTTPServerGet(loop,message,"errored")};if tsNullish(err){err=tsFSException("Error","ERR_STREAM_DESTROYED","Cannot call write after a stream was destroyed")};callback:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{tsInternalHTTP1FallbackFailPending(loop,message,err);return tsU}));tsNodeHTTPServerCall(loop,handle,"socketEnded",callback);request:=tsNodeHTTPServerGet(loop,message,"req");if !tsNullish(request)&&!tsTruthy(tsNodeHTTPIncomingGet(loop,request,"destroyed")){if tsNumber(tsNodeHTTPServerCall(loop,request,"listenerCount",tsInternalHTTP1FallbackKey("error")))>0{tsNodeHTTPIncomingCall(loop,request,"destroy",tsFSException("Error","ECONNRESET","aborted"))}else{tsNodeHTTPIncomingCall(loop,request,"destroy")}}};tsInternalHTTP1FallbackCallPipeline(loop,state.pipeline,"abortQueuedPipelinedResponses",socket);tsNodeHTTPCommonCall(loop,parser,"close");return tsU
	}
	return tsU
}

func tsInternalHTTP1FallbackCallPipeline(loop *tsLoop, pipeline tsValue, name string, args ...tsValue) tsValue { method:=tsGet(loop,pipeline,tsInternalHTTP1FallbackKey(name));if method.kind==tsFunctionKind{return tsCall(loop,method,args...)};return tsU }
func tsInternalHTTP1FallbackCloseIdle(loop *tsLoop, server tsValue) tsValue {
	tsInternalHTTP1FallbackInit(loop);connections:=tsInternalHTTP1FallbackSym(loop,server,"kHttp1Connections");if connections.kind!=tsCollectionKind{return tsU}
	for _,entry:=range (*tsCollection)(connections.ref).entries{if !entry.live{continue};socket:=entry.key;if !tsTruthy(tsInternalHTTP1FallbackSym(loop,socket,"kHttp1ActiveRequests"))&&!tsTruthy(tsInternalHTTP1FallbackSym(loop,socket,"kHttp1MessageOpen"))&&!tsTruthy(tsInternalHTTP1FallbackGet(loop,socket,"destroyed")){tsNodeHTTPServerCall(loop,socket,"destroy")}}
	return tsU
}
func tsInternalHTTP1FallbackCloseAll(loop *tsLoop, server tsValue) tsValue {
	tsInternalHTTP1FallbackInit(loop);connections:=tsInternalHTTP1FallbackSym(loop,server,"kHttp1Connections");if connections.kind!=tsCollectionKind{return tsU}
	for _,entry:=range (*tsCollection)(connections.ref).entries{if entry.live&&!tsTruthy(tsInternalHTTP1FallbackGet(loop,entry.key,"destroyed")){tsNodeHTTPServerCall(loop,entry.key,"destroy")}}
	return tsU
}
func tsInternalHTTP1FallbackModule(loop *tsLoop) tsValue {
	if !tsIsUndefined(tsInternalHTTP1FallbackSingleton.module){return tsInternalHTTP1FallbackSingleton.module}
	tsInternalHTTP1FallbackInit(loop);module:=tsNewObject()
	module.set("connectionListenerHTTP1",tsNamedNativeMethod("connectionListenerHTTP1",3,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsInternalHTTP1FallbackConnectionListener(loop,tsArg(args,0),tsArg(args,1),tsArg(args,2))}))
	module.set("closeIdleHttp1Connections",tsNamedNativeMethod("closeIdleHttp1Connections",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsInternalHTTP1FallbackCloseIdle(loop,tsArg(args,0))}))
	module.set("closeAllHttp1Connections",tsNamedNativeMethod("closeAllHttp1Connections",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsInternalHTTP1FallbackCloseAll(loop,tsArg(args,0))}))
	module.set("kHttp1Connections",tsInternalHTTP1FallbackSingleton.kConnections);module.set("default",tsObjectValue(module));tsInternalHTTP1FallbackSingleton.module=tsObjectValue(module);return tsInternalHTTP1FallbackSingleton.module
}
`
