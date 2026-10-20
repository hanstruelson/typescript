package goemit

// nodeHTTP1ServerRuntime owns the HTTP/1 parser integration, response writer,
// connection lifecycle, and response ordering for node:http servers.
const nodeHTTP1ServerRuntime = `
type tsNodeHTTP1ResponseState struct {
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

var tsNodeHTTP1ServerSingleton struct {
	sync.Once
	kConnections, kActiveRequests, kMessageOpen, kPendingCallbacks tsValue
	kQueue, kQueueState, kMustCloseConnection tsValue
	module, pipeline tsValue
}

func tsNodeHTTP1ServerKey(name string) tsValue { return tsStringReference(tsStringUTF8(name)) }
func tsNodeHTTP1ServerGet(loop *tsLoop, value tsValue, name string) tsValue { return tsGet(loop, value, tsNodeHTTP1ServerKey(name)) }
func tsNodeHTTP1ServerSet(loop *tsLoop, value tsValue, name string, next tsValue) { tsSet(loop, value, tsNodeHTTP1ServerKey(name), next) }
func tsNodeHTTP1ServerSym(loop *tsLoop, value tsValue, name string) tsValue {
	tsNodeHTTP1ServerInit(loop)
	var key tsValue
	switch name {
	case "kHttp1Connections": key = tsNodeHTTP1ServerSingleton.kConnections
	case "kHttp1ActiveRequests": key = tsNodeHTTP1ServerSingleton.kActiveRequests
	case "kHttp1MessageOpen": key = tsNodeHTTP1ServerSingleton.kMessageOpen
	case "kPendingCallbacks": key = tsNodeHTTP1ServerSingleton.kPendingCallbacks
	case "kQueue": key = tsNodeHTTP1ServerSingleton.kQueue
	case "kQueueState": key = tsNodeHTTP1ServerSingleton.kQueueState
	case "kMustCloseConnection": key = tsNodeHTTP1ServerSingleton.kMustCloseConnection
	}
	return tsGet(loop, value, key)
}
func tsNodeHTTP1ServerPut(loop *tsLoop, value tsValue, name string, next tsValue) {
	tsNodeHTTP1ServerInit(loop)
	var key tsValue
	switch name {
	case "kHttp1Connections": key = tsNodeHTTP1ServerSingleton.kConnections
	case "kHttp1ActiveRequests": key = tsNodeHTTP1ServerSingleton.kActiveRequests
	case "kHttp1MessageOpen": key = tsNodeHTTP1ServerSingleton.kMessageOpen
	case "kPendingCallbacks": key = tsNodeHTTP1ServerSingleton.kPendingCallbacks
	case "kQueue": key = tsNodeHTTP1ServerSingleton.kQueue
	case "kQueueState": key = tsNodeHTTP1ServerSingleton.kQueueState
	case "kMustCloseConnection": key = tsNodeHTTP1ServerSingleton.kMustCloseConnection
	}
	tsSet(loop, value, key, next)
}
func tsNodeHTTP1ServerInit(loop *tsLoop) {
	tsNodeHTTP1ServerSingleton.Do(func() {
		symbolFor := tsGet(loop, tsSymbolFunction(), tsNodeHTTP1ServerKey("for"))
		tsNodeHTTP1ServerSingleton.kConnections = tsCall(loop, symbolFor, tsNodeHTTP1ServerKey("http1Connections"))
		tsNodeHTTP1ServerSingleton.kActiveRequests = tsCall(loop, symbolFor, tsNodeHTTP1ServerKey("http1ActiveRequests"))
		tsNodeHTTP1ServerSingleton.kMessageOpen = tsCall(loop, symbolFor, tsNodeHTTP1ServerKey("http1MessageOpen"))
		tsNodeHTTP1ServerSingleton.kPendingCallbacks = tsInternalHTTPGet(loop, tsInternalHTTPModule(loop), "kPendingCallbacks")
		tsNodeHTTP1ServerSingleton.kQueue = tsCall(loop, symbolFor, tsNodeHTTP1ServerKey("http1ResponseQueue"))
		tsNodeHTTP1ServerSingleton.kQueueState = tsCall(loop, symbolFor, tsNodeHTTP1ServerKey("http1QueuedResponse"))
		tsNodeHTTP1ServerSingleton.kMustCloseConnection = tsCall(loop, symbolFor, tsNodeHTTP1ServerKey("kMustCloseConnection"))
	})
}
func tsNodeHTTP1ServerRethrow(loop *tsLoop, err tsValue) {
	callback := tsFunctionValue(tsFunc(func(loop *tsLoop, args ...tsValue) tsValue { panic(tsThrown{tsArg(args, 0)}) }))
	tsInternalHTTPNextTick(loop, callback, err)
}
func tsNodeHTTP1ServerSafeCallback(loop *tsLoop, callback tsValue, args ...tsValue) {
	if callback.kind != tsFunctionKind { return }
	var thrown any
	func() { defer func() { thrown = recover() }(); tsCall(loop, callback, args...) }()
	if thrown != nil { tsNodeHTTP1ServerRethrow(loop, tsUnwrap(thrown)) }
}
func tsNodeHTTP1ServerFailPending(loop *tsLoop, response, err tsValue) {
	tsNodeHTTP1ServerInit(loop)
	callbacks := tsGet(loop, response, tsNodeHTTP1ServerSingleton.kPendingCallbacks)
	tsSet(loop, response, tsNodeHTTP1ServerSingleton.kPendingCallbacks, tsArrayValue(&tsArray{}))
	if callbacks.kind != tsArrayKind { return }
	for _, callback := range (*tsArray)(callbacks.ref).values { tsNodeHTTP1ServerSafeCallback(loop, callback, err) }
}
func tsNodeHTTP1ServerDate() string { return time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT") }
func tsNodeHTTP1ServerToBuffer(loop *tsLoop, chunk, encoding tsValue) tsValue {
	if tsNullish(chunk) { return tsNull }
	if chunk.kind == tsStringKind { return tsNodeBufferFrom(loop, []tsValue{chunk, encoding}) }
	return chunk
}
func tsNodeHTTP1ServerString(loop *tsLoop, value tsValue) string { return tsNodeHTTPOutgoingText(loop, value) }
func tsNodeHTTP1ServerSocketWrite(loop *tsLoop, state *tsNodeHTTP1ResponseState, data, callback tsValue) tsValue {
	if tsTruthy(tsNodeHTTP1ServerGet(loop, state.socket, "writableEnded")) { return tsBooleanValue(false) }
	if callback.kind == tsFunctionKind { return tsNodeHTTPServerCall(loop, state.socket, "write", data, callback) }
	return tsNodeHTTPServerCall(loop, state.socket, "write", data)
}
func tsNodeHTTP1ServerWriteBody(loop *tsLoop, state *tsNodeHTTP1ResponseState, buffer, callback tsValue) int {
	if state.noBody { return 0 }
	length := 0
	if !tsNullish(buffer) { length = len(tsNodeBufferVisibleBytes(buffer)) }
	if length > 0 {
		if state.chunked {
			tsNodeHTTP1ServerSocketWrite(loop, state, tsNodeHTTP1ServerKey(strconv.FormatInt(int64(length), 16)+"\r\n"), tsU)
			tsNodeHTTP1ServerSocketWrite(loop, state, buffer, tsU)
			tsNodeHTTP1ServerSocketWrite(loop, state, tsNodeHTTP1ServerKey("\r\n"), callback)
		} else { tsNodeHTTP1ServerSocketWrite(loop, state, buffer, callback) }
	}
	return length
}
func tsNodeHTTP1ServerRunFrame(loop *tsLoop, frame, callback tsValue) {
	module := tsInternalAsyncContextFrameModule(loop)
	run := tsGet(loop, module, tsNodeHTTP1ServerKey("run"))
	if run.kind == tsFunctionKind { tsCall(loop, run, frame, callback) } else { tsCall(loop, callback) }
}
func tsNodeHTTP1ServerSettleWritable(loop *tsLoop, state *tsNodeHTTP1ResponseState, callback tsValue) {
	frame := state.onWritableFrame
	state.onWritable = tsNull
	state.onWritableFrame = tsU
	tsNodeHTTP1ServerRunFrame(loop, frame, tsFunctionValue(tsFunc(func(loop *tsLoop, _ ...tsValue) tsValue { tsNodeHTTP1ServerSafeCallback(loop, callback); return tsU })))
}
func tsNodeHTTP1ServerFlush(loop *tsLoop, state *tsNodeHTTP1ResponseState) {
	if !state.ended || state.finished { return }
	state.finished = true
	tsNodeHTTP1ServerSafeCallback(loop, tsNodeHTTP1ServerHandleGet(loop, state.handle, "onflushed"))
}
func tsNodeHTTP1ServerHandleGet(loop *tsLoop, value tsValue, name string) tsValue {
	if value.kind != tsObjectKind { return tsU }
	return (*tsObject)(value.ref).get(name)
}
func tsNodeHTTP1ServerEndWritten(loop *tsLoop, state *tsNodeHTTP1ResponseState, err tsValue) {
	if !tsTruthy(err) && !tsTruthy(tsNodeHTTP1ServerGet(loop, state.socket, "destroyed")) { tsNodeHTTP1ServerFlush(loop, state) }
}
func tsNodeHTTP1ServerCountBodyBytes(loop *tsLoop, state *tsNodeHTTP1ResponseState, buffer, expected tsValue, fromEnd bool) int {
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
func tsNodeHTTP1ServerWriteHead(loop *tsLoop, state *tsNodeHTTP1ResponseState, contentLength tsValue) {
	if state.headWritten { return }
	state.headWritten = true
	statusCode := 200
	statusMessage := tsU
	if !tsNullish(state.head) { statusCode = int(tsNumber(tsNodeHTTP1ServerGet(loop, state.head, "statusCode"))); statusMessage = tsNodeHTTP1ServerGet(loop, state.head, "statusMessage") }
	if statusMessage.kind != tsStringKind || tsText(statusMessage) == "" {
		codes := tsInternalHTTPGet(loop, tsInternalHTTPModule(loop), "STATUS_CODES")
		statusMessage = tsGet(loop, codes, tsNodeHTTP1ServerKey(strconv.Itoa(statusCode)))
		if !tsTruthy(statusMessage) { statusMessage = tsNodeHTTP1ServerKey("unknown") }
	}
	output := "HTTP/1.1 " + strconv.Itoa(statusCode) + " " + tsText(statusMessage) + "\r\n"
	hasContentLength, hasTransferEncoding, hasDate, hasConnection, hasKeepAlive := false, false, false, false, false
	headers := tsU
	if !tsNullish(state.head) { headers = tsNodeHTTP1ServerGet(loop, state.head, "headers") }
	var writeHeader func(tsValue, tsValue)
	writeHeader = func(nameValue, value tsValue) {
		name := tsText(nameValue)
		if name == "\x00" {
			if tsText(value) == "2" { state.noBody = true } else { state.closeDelimited = true }
			return
		}
		lower := strings.ToLower(name)
		if value.kind == tsArrayKind {
			items := (*tsArray)(value.ref).values
			if lower == "cookie" {
				parts := make([]string, 0, len(items))
				for _, item := range items { parts = append(parts, tsNodeHTTP1ServerString(loop, item)) }
				value = tsNodeHTTP1ServerKey(strings.Join(parts, "; "))
			} else {
				for _, item := range items { writeHeader(nameValue, item) }
				return
			}
		}
		switch lower {
		case "content-length": hasContentLength = true
		case "transfer-encoding": hasTransferEncoding = true; chunked, _, _ := tsNodeHTTPCommonTransferChunked(tsNodeHTTP1ServerString(loop, value)); if chunked { state.chunked = true }
		case "date": hasDate = true
		case "connection": hasConnection = true
		case "keep-alive": hasKeepAlive = true
		}
		output += name + ": " + tsNodeHTTP1ServerString(loop, value) + "\r\n"
	}
	if headers.kind == tsArrayKind {
		values := (*tsArray)(headers.ref).values
		if len(values) > 0 && values[0].kind == tsArrayKind {
			for _, entry := range values { writeHeader(tsGet(loop, entry, tsNodeHTTP1ServerKey("0")), tsGet(loop, entry, tsNodeHTTP1ServerKey("1"))) }
		} else {
			for i := 0; i+1 < len(values); i += 2 { writeHeader(values[i], values[i+1]) }
		}
	} else if headers.kind == tsObjectKind || headers.kind == tsInstanceKind {
		for _, name := range tsObjectKeys(headers).values { writeHeader(name, tsGet(loop, headers, name)) }
	}
	autoBits := 0
	if !tsNullish(state.head) { autoBits = int(tsNumber(tsNodeHTTP1ServerGet(loop, state.head, "autoHeaderBits"))) }
	state.noBody = statusCode < 200 || statusCode == 204 || statusCode == 304 || autoBits&32 != 0
	chunkedFromAutoBits := autoBits&16 != 0
	autoContentLength := tsNull
	autoChunked := false
	if !hasContentLength && !hasTransferEncoding && !state.noBody && !state.closeDelimited {
		if chunkedFromAutoBits || tsIsNull(contentLength) { state.chunked = true; autoChunked = !chunkedFromAutoBits } else { autoContentLength = contentLength }
	}
	if tsNullish(state.head) {
		if !hasDate { output += "Date: " + tsNodeHTTP1ServerDate() + "\r\n" }
		if !hasConnection && !state.closeDelimited {
			if state.shouldKeepAlive { output += "Connection: keep-alive\r\n"; if !hasKeepAlive { timeout := state.keepAliveTimeout; if timeout == 0 { timeout = 5000 }; keepAlive := "Keep-Alive: timeout=" + strconv.Itoa(int(math.Floor(timeout/1000))); maxRequests := tsNodeHTTP1ServerGet(loop, state.head, "maxRequestsPerSocket"); if tsIsNumeric(maxRequests) && tsNumber(maxRequests) > 0 { keepAlive += ", max=" + strconv.Itoa(int(tsNumber(maxRequests))) }; output += keepAlive + "\r\n" } } else { output += "Connection: close\r\n" }
		}
	} else {
		if !hasDate && autoBits&1 != 0 { output += "Date: " + tsNodeHTTP1ServerDate() + "\r\n" }
		if !hasConnection { if autoBits&2 != 0 { output += "Connection: keep-alive\r\n"; if !hasKeepAlive && autoBits&8 != 0 { keepAlive := "Keep-Alive: timeout=" + strconv.Itoa(int(tsNumber(tsNodeHTTP1ServerGet(loop, state.head, "keepAliveTimeoutSecs")))); maxRequests := tsNodeHTTP1ServerGet(loop, state.head, "maxRequestsPerSocket"); if tsIsNumeric(maxRequests) && tsNumber(maxRequests) > 0 { keepAlive += ", max=" + strconv.Itoa(int(tsNumber(maxRequests))) }; output += keepAlive + "\r\n" } } else if autoBits&4 != 0 { output += "Connection: close\r\n" } }
	}
	if !tsIsNull(autoContentLength) { output += "Content-Length: " + tsNodeHTTP1ServerString(loop, autoContentLength) + "\r\n" } else if autoChunked || chunkedFromAutoBits && state.chunked { output += "Transfer-Encoding: chunked\r\n" }
	output += "\r\n"
	tsNodeHTTP1ServerSocketWrite(loop, state, tsNodeHTTP1ServerKey(output), tsU)
}
func tsNodeHTTP1ServerCreateResponseHandle(loop *tsLoop, socket, shouldKeepAlive, keepAliveTimeout tsValue) tsValue {
	state := &tsNodeHTTP1ResponseState{socket: socket, head: tsNull, shouldKeepAlive: tsTruthy(shouldKeepAlive), keepAliveTimeout: tsNumber(keepAliveTimeout), onWritable: tsNull, onWritableFrame: tsU, onFinished: tsNull}
	object := tsNewObject(); object.prototype = tsNull; object.descriptors = map[string]*tsDescriptor{}; object.set("onflushed", tsNull); result := tsObjectValue(object); state.handle = result
	method := func(name string, length int) { object.set(name, tsNamedNativeMethod(name, length, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue { return tsNodeHTTP1ServerHandleMethod(loop, state, name, args) })) }
	for _, name := range []string{"socketDrained", "socketEnded", "writeContinue", "flushHeaders", "flushed", "abort"} { method(name, 1) }
	for _, name := range []string{"cork", "writeInformational", "writeHead", "writeHeadAndEnd", "write", "end"} { method(name, 4) }
	define := func(name string, getter func(*tsLoop) tsValue, setter func(*tsLoop, tsValue)) { object.descriptors[name] = &tsDescriptor{accessor:true, configurable:true, enumerable:true, getter:tsNamedNativeMethod("get "+name,0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{return getter(loop)}), setter:tsNamedNativeMethod("set "+name,1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{setter(loop,tsArg(args,0));return tsU})} }
	define("flags",func(loop *tsLoop)tsValue{if state.aborted||tsTruthy(tsNodeHTTP1ServerGet(loop,state.socket,"destroyed")){return tsNumberValue(1)};return tsNumberValue(0)},func(*tsLoop,tsValue){})
	define("ended",func(*tsLoop)tsValue{return tsBooleanValue(state.ended)},func(_ *tsLoop,v tsValue){state.ended=tsTruthy(v)})
	define("finished",func(*tsLoop)tsValue{return tsBooleanValue(state.finished)},func(_ *tsLoop,v tsValue){state.finished=tsTruthy(v)})
	define("aborted",func(*tsLoop)tsValue{return tsBooleanValue(state.aborted)},func(_ *tsLoop,v tsValue){state.aborted=tsTruthy(v)})
	define("shouldKeepAlive",func(*tsLoop)tsValue{return tsBooleanValue(state.shouldKeepAlive)},func(_ *tsLoop,v tsValue){state.shouldKeepAlive=tsTruthy(v)})
	define("onfinished",func(*tsLoop)tsValue{return state.onFinished},func(_ *tsLoop,v tsValue){state.onFinished=v})
	define("onwritable",func(*tsLoop)tsValue{if tsNullish(state.onWritable){return tsU};return state.onWritable},func(loop *tsLoop,v tsValue){state.onWritable=v;if tsTruthy(v){state.onWritableFrame=tsCall(loop,tsGet(loop,tsInternalAsyncContextFrameModule(loop),tsNodeHTTP1ServerKey("current")))}else{state.onWritableFrame=tsU}})
	define("bufferedAmount",func(loop *tsLoop)tsValue{if tsTruthy(state.onWritable)||state.ended&&!state.finished{return tsNodeHTTP1ServerGet(loop,state.socket,"writableLength")};return tsNumberValue(0)},func(*tsLoop,tsValue){})
	return result
}
func tsNodeHTTP1ServerHandleMethod(loop *tsLoop, state *tsNodeHTTP1ResponseState, name string, args []tsValue) tsValue {
	arg := func(index int) tsValue { return tsArg(args, index) }
	switch name {
	case "socketDrained": if tsTruthy(state.onWritable) { tsNodeHTTP1ServerSettleWritable(loop,state,state.onWritable) }; return tsU
	case "socketEnded": if tsTruthy(state.onWritable) { tsNodeHTTP1ServerSettleWritable(loop,state,arg(0)) }; return tsU
	case "cork": return tsCall(loop,arg(0))
	case "writeContinue": tsNodeHTTP1ServerSocketWrite(loop,state,tsNodeHTTP1ServerKey("HTTP/1.1 100 Continue\r\n\r\n"),tsU); return tsU
	case "writeInformational": if !tsTruthy(tsNodeHTTP1ServerGet(loop,state.socket,"writableEnded")) { tsNodeHTTPServerCall(loop,state.socket,"write",arg(0),arg(1)) }; return tsU
	case "writeHead":
		status := int(tsNumber(arg(0))); if status < 100 || status > 999 { panic(tsThrown{tsFSException("RangeError","ERR_HTTP_INVALID_STATUS_CODE","Invalid status code: "+tsNodeHTTP1ServerString(loop,arg(0)))}) }
		message:=arg(1); common:=tsNodeHTTPCommonModule(loop);invalid:=tsNodeHTTPCommonGet(loop,common,"_checkInvalidHeaderChar");if message.kind==tsStringKind&&tsTruthy(tsCall(loop,invalid,message)){panic(tsThrown{tsFSException("TypeError","ERR_INVALID_CHAR","Invalid character in statusMessage")})}
		head:=tsNewObject();head.set("statusCode",tsNumberValue(float64(status)));head.set("statusMessage",message);head.set("headers",arg(2));head.set("autoHeaderBits",arg(3));head.set("keepAliveTimeoutSecs",arg(4));head.set("maxRequestsPerSocket",arg(5));state.head=tsObjectValue(head);return tsU
	case "flushHeaders": tsNodeHTTP1ServerWriteHead(loop,state,tsNull);return tsU
	case "writeHeadAndEnd": tsNodeHTTP1ServerHandleMethod(loop,state,"writeHead",[]tsValue{arg(0),arg(1),arg(2),arg(6),arg(7),arg(8)});return tsNodeHTTP1ServerHandleMethod(loop,state,"end",[]tsValue{arg(3),arg(4),tsU,arg(5)})
	case "write":
		buffer:=tsNodeHTTP1ServerToBuffer(loop,arg(0),arg(1));length:=tsNodeHTTP1ServerCountBodyBytes(loop,state,buffer,arg(3),false);tsNodeHTTP1ServerWriteHead(loop,state,tsNull);tsNodeHTTP1ServerWriteBody(loop,state,buffer,tsU)
		if tsIsNull(buffer)||!tsTruthy(tsNodeHTTP1ServerGet(loop,state.socket,"writableNeedDrain")){return tsNumberValue(float64(length))}
		callback:=arg(2);if tsTruthy(callback){state.onWritable=callback;state.onWritableFrame=tsCall(loop,tsGet(loop,tsInternalAsyncContextFrameModule(loop),tsNodeHTTP1ServerKey("current")))};if length>0{return tsNumberValue(float64(-length))};return tsNumberValue(-1)
	case "end":
		if state.ended{return tsNumberValue(0)};state.onWritable=tsNull;state.onWritableFrame=tsU;buffer:=tsNodeHTTP1ServerToBuffer(loop,arg(0),arg(1));length:=tsNodeHTTP1ServerCountBodyBytes(loop,state,buffer,arg(3),true);tsNodeHTTP1ServerWriteHead(loop,state,tsNumberValue(float64(length)));terminated:=state.chunked&&!state.noBody
		endCallback:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{tsNodeHTTP1ServerEndWritten(loop,state,tsArg(args,0));return tsU}))
		tsNodeHTTP1ServerWriteBody(loop,state,buffer,func()tsValue{if terminated{return tsU};return endCallback}());if terminated{tsNodeHTTP1ServerSocketWrite(loop,state,tsNodeHTTP1ServerKey("0\r\n\r\n"),endCallback)};state.ended=true
		writableLength:=tsNumber(tsNodeHTTP1ServerGet(loop,state.socket,"writableLength"));if writableLength>0&&!tsTruthy(tsNodeHTTP1ServerGet(loop,state.socket,"destroyed")){if tsTruthy(tsNodeHTTP1ServerGet(loop,state.socket,"writableEnded")){tsNodeHTTPServerCall(loop,state.socket,"once",tsNodeHTTP1ServerKey("finish"),endCallback)}else if length==0&&!terminated{tsNodeHTTPServerCall(loop,state.socket,"write",tsNodeHTTP1ServerKey(""),tsNodeHTTP1ServerKey("latin1"),endCallback)}}else{state.finished=true}
		finishedCallback:=state.onFinished;state.onFinished=tsNull;if tsTruthy(finishedCallback){tsNodeHTTP1ServerSafeCallback(loop,finishedCallback)};if state.closeDelimited&&!tsTruthy(tsNodeHTTP1ServerGet(loop,state.socket,"destroyed")){tsNodeHTTPServerCall(loop,state.socket,"end")};if state.finished{return tsNumberValue(float64(length))};return tsNumberValue(float64(-(length+1)))
	case "flushed": tsNodeHTTP1ServerFlush(loop,state);return tsU
	case "abort": state.aborted=true;if !tsTruthy(tsNodeHTTP1ServerGet(loop,state.socket,"destroyed")){tsNodeHTTPServerCall(loop,state.socket,"destroy")};return tsU
	}
	return tsU
}

type tsNodeHTTP1Connection struct {
	loop *tsLoop
	server, socket, options, parser, common, http, pipeline tsValue
	req, pendingUpgrade tsValue
	flushedHeaders []tsValue
	flushedURL string
	requestsCount int
	closeAfterParse bool
	maxRequestsPerSocket tsValue
	keepAliveTimeout tsValue
	connections tsValue
}

func tsNodeHTTP1ServerPipeline(loop *tsLoop) tsValue {
	tsNodeHTTP1ServerInit(loop)
	if !tsIsUndefined(tsNodeHTTP1ServerSingleton.pipeline) { return tsNodeHTTP1ServerSingleton.pipeline }
	pipeline := tsObjectValue(tsNewObject())
	tsNodeHTTP1ServerSingleton.pipeline = pipeline
	install := func(name string, fn tsValue) { tsSet(loop, pipeline, tsNodeHTTP1ServerKey(name), fn) }
	install("kMustCloseConnection", tsNodeHTTP1ServerSingleton.kMustCloseConnection)
	install("queuePipelinedResponse", tsNamedNativeMethod("queuePipelinedResponse", 3, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
		socket := tsArg(args, 0); var response tsValue = tsArg(args, 1); queue := tsGet(loop, socket, tsNodeHTTP1ServerSingleton.kQueue)
		if queue.kind != tsArrayKind { queue = tsArrayValue(&tsArray{}); tsSet(loop, socket, tsNodeHTTP1ServerSingleton.kQueue, queue) }
		state := tsNewObject(); state.set("socket", socket); state.set("isAncient", tsArg(args, 2)); state.set("ops", tsArrayValue(&tsArray{})); state.set("bytes",tsNumberValue(0));state.set("ended",tsBooleanValue(false));state.set("needDrain",tsBooleanValue(false));tsSet(loop, response, tsNodeHTTP1ServerSingleton.kQueueState, tsObjectValue(state))
		a := (*tsArray)(queue.ref); a.values = append(a.values, response); return tsU
	}))
	install("lastPipelinedResponse", tsNamedNativeMethod("lastPipelinedResponse", 1, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
		queue := tsGet(loop, tsArg(args, 0), tsNodeHTTP1ServerSingleton.kQueue); if queue.kind != tsArrayKind || len((*tsArray)(queue.ref).values) == 0 { return tsU }; a := (*tsArray)(queue.ref); return a.values[len(a.values)-1]
	}))
	install("advanceResponsePipeline", tsNamedNativeMethod("advanceResponsePipeline", 2, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
		server, socket := tsArg(args, 0), tsArg(args, 1); if tsNullish(socket) || tsTruthy(tsNodeHTTP1ServerGet(loop, socket, "destroyed")) || tsTruthy(tsNodeHTTP1ServerGet(loop, socket, "writableEnded")) || !tsNullish(tsNodeHTTP1ServerGet(loop, socket, "_httpMessage")) { return tsU }
		queue := tsGet(loop, socket, tsNodeHTTP1ServerSingleton.kQueue); if queue.kind != tsArrayKind { return tsU }; a := (*tsArray)(queue.ref); if len(a.values) == 0 { return tsU }
		var response tsValue = a.values[0]; a.values = a.values[1:]; tsSet(loop, socket, tsNodeHTTP1ServerSingleton.kQueue, tsArrayValue(a))
		queueState:=tsGet(loop,response,tsNodeHTTP1ServerSingleton.kQueueState);ops:=tsNodeHTTP1ServerGet(loop,queueState,"ops");ended:=tsTruthy(tsNodeHTTP1ServerGet(loop,queueState,"ended"));queuedBytes:=tsNumber(tsNodeHTTP1ServerGet(loop,queueState,"bytes"));remainingBytes:=tsNumber(tsNodeHTTP1ServerGet(loop,socket,"_httpQueuedBytes"))-queuedBytes;if remainingBytes<0{remainingBytes=0};tsNodeHTTP1ServerSet(loop,socket,"_httpQueuedBytes",tsNumberValue(remainingBytes));tsNodeHTTP1ServerSet(loop,response,"outputSize",tsNumberValue(0))
		tsSet(loop,response,tsNodeHTTP1ServerSingleton.kQueueState,tsU)
		assign := tsNodeHTTP1ServerGet(loop, response, "assignSocket"); if assign.kind == tsFunctionKind { tsCallReceiver(loop, assign, response, socket) }
		if ended { tsNodeHTTP1ServerSet(loop,response,"finished",tsBooleanValue(false)) }
		if ops.kind==tsArrayKind { opArray:=(*tsArray)(ops.ref);opList:=make([]tsValue,len(opArray.values));copy(opList,opArray.values);for _,op:=range opList { kind:=tsNodeHTTP1ServerGet(loop,op,"0");chunk:=tsNodeHTTP1ServerGet(loop,op,"1");encoding:=tsNodeHTTP1ServerGet(loop,op,"2");callback:=tsNodeHTTP1ServerGet(loop,op,"3");if tsText(kind)=="end"{tsNodeHTTPServerCall(loop,response,"end",chunk,encoding,callback)}else{tsNodeHTTPServerCall(loop,response,"write",chunk,encoding,callback)} } }
		flush := tsNodeHTTP1ServerGet(loop, response, "_flush"); if flush.kind == tsFunctionKind { tsCallReceiver(loop, flush, response) }
		_ = server; return tsU
	}))
	install("abortQueuedPipelinedResponses", tsNamedNativeMethod("abortQueuedPipelinedResponses", 2, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
		socket := tsArg(args, 0); queue := tsGet(loop, socket, tsNodeHTTP1ServerSingleton.kQueue); if queue.kind != tsArrayKind { return tsU };queueArray:=(*tsArray)(queue.ref);items:=make([]tsValue,len(queueArray.values));copy(items,queueArray.values);queueArray.resize(0);err := tsArg(args, 1);tsNodeHTTP1ServerSet(loop,socket,"_httpQueuedBytes",tsNumberValue(0))
		for _, response := range items { request := tsNodeHTTP1ServerGet(loop, response, "req"); state:=tsGet(loop,response,tsNodeHTTP1ServerSingleton.kQueueState);ops:=tsNodeHTTP1ServerGet(loop,state,"ops");if ops.kind==tsArrayKind{for _,op:=range (*tsArray)(ops.ref).values{callback:=tsNodeHTTP1ServerGet(loop,op,"3");if callback.kind==tsFunctionKind{tsNodeHTTP1ServerSafeCallback(loop,callback,err)}}}; if !tsTruthy(tsNodeHTTP1ServerGet(loop, request, "destroyed")) { if tsTruthy(err) { tsNodeHTTPIncomingCall(loop, request, "destroy", err) } else { tsNodeHTTPIncomingCall(loop, request, "destroy") } }; if !tsTruthy(tsNodeHTTP1ServerGet(loop, response, "destroyed")) { tsNodeHTTPIncomingCall(loop, response, "destroy") } }
		return tsU
	}))
	install("maybePauseHTTP1Reads", tsNamedNativeMethod("maybePauseHTTP1Reads", 1, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
		socket := tsArg(args, 0); if tsTruthy(tsNodeHTTP1ServerGet(loop, socket, "_paused")) { return tsU }; state := tsNodeHTTP1ServerGet(loop, socket, "_writableState");pending:=tsNumber(tsNodeHTTP1ServerGet(loop,socket,"writableLength"))+tsNumber(tsNodeHTTP1ServerGet(loop,socket,"_httpQueuedBytes")); if tsTruthy(tsNodeHTTP1ServerGet(loop, state, "needDrain")) || pending >= tsNumber(tsNodeHTTP1ServerGet(loop, socket, "writableHighWaterMark")) { tsNodeHTTP1ServerSet(loop, socket, "_paused", tsBooleanValue(true)); tsNodeHTTPServerCall(loop, socket, "pause") }; return tsU
	}))
	install("resumeHTTP1ReadsOnDrain", tsNamedNativeMethod("resumeHTTP1ReadsOnDrain", 1, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
		socket := tsArg(args, 0); pending:=tsNumber(tsNodeHTTP1ServerGet(loop,socket,"writableLength"))+tsNumber(tsNodeHTTP1ServerGet(loop,socket,"_httpQueuedBytes"));if tsTruthy(tsNodeHTTP1ServerGet(loop, socket, "_paused")) && pending <= tsNumber(tsNodeHTTP1ServerGet(loop, socket, "writableHighWaterMark")) { tsNodeHTTP1ServerSet(loop, socket, "_paused", tsBooleanValue(false)); tsNodeHTTPServerCall(loop, socket, "resume") }; return tsU
	}))
	install("finishDrainedResponse", tsNamedNativeMethod("finishDrainedResponse", 1, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
		var response tsValue = tsArg(args, 0); tsNodeHTTP1ServerSet(loop, response, "writableFinished", tsBooleanValue(true)); if !tsTruthy(tsNodeHTTP1ServerGet(loop, response, "finished")) { tsNodeHTTPServerEmit(loop, response, "finish") }; return tsU
	}))
	return pipeline
}

func tsNodeHTTP1ServerParserSlot(loop *tsLoop, parserCtor tsValue, name string) tsValue { return tsGet(loop, parserCtor, tsNodeHTTP1ServerKey(name)) }
func tsNodeHTTP1ServerCallback(loop *tsLoop, state *tsNodeHTTP1Connection, name string, length int) tsValue {
	return tsNamedNativeMethod(name, length, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue { return tsNodeHTTP1ServerConnectionDispatch(loop, state, name, args) })
}
func tsNodeHTTP1ServerConnectionListener(loop *tsLoop, server, socket, options tsValue) tsValue {
	state := &tsNodeHTTP1Connection{loop:loop, server:server, socket:socket, options:options, req:tsNull, pendingUpgrade:tsNull, flushedHeaders:[]tsValue{}}
	http1Options := tsNodeHTTP1ServerGet(loop, options, "http1Options")
	state.http = tsNodeHTTP1ServerGet(loop, options, "httpModule")
	if tsNullish(state.http) {
		state.http = tsObjectValue(tsNewObject())
		incoming := tsNodeHTTP1ServerGet(loop, http1Options, "IncomingMessage")
		if tsNullish(incoming) { incoming = tsNodeHTTPIncomingGet(loop, tsNodeHTTPIncomingModule(loop), "IncomingMessage") }
		response := tsNodeHTTP1ServerGet(loop, http1Options, "ServerResponse")
		if tsNullish(response) { response = tsNodeHTTPServerAPIGet(loop, tsNodeHTTPServerAPIModule(loop), "ServerResponse") }
		tsSet(loop, state.http, tsNodeHTTP1ServerKey("IncomingMessage"), incoming)
		tsSet(loop, state.http, tsNodeHTTP1ServerKey("ServerResponse"), response)
	}
	state.common = tsNodeHTTPCommonModule(loop)
	state.pipeline = tsNodeHTTP1ServerPipeline(loop)
	internal := tsInternalHTTPModule(loop); state.keepAliveTimeout = tsNodeHTTP1ServerGet(loop, server, "keepAliveTimeout")
	if !tsIsNumeric(state.keepAliveTimeout) { state.keepAliveTimeout = tsNumberValue(5000) }
	state.connections = tsNodeHTTP1ServerSym(loop, server, "kHttp1Connections")
	if state.connections.kind != tsCollectionKind { state.connections = tsNewCollection(true, tsU); tsNodeHTTP1ServerPut(loop, server, "kHttp1Connections", state.connections) }
	add := tsGet(loop, state.connections, tsNodeHTTP1ServerKey("add")); if add.kind == tsFunctionKind { tsCallReceiver(loop, add, state.connections, socket) }
	tsNodeHTTP1ServerSet(loop, socket, "server", server); tsNodeHTTP1ServerPut(loop, socket, "kHttp1ActiveRequests", tsNumberValue(0)); tsNodeHTTP1ServerPut(loop, socket, "kHttp1MessageOpen", tsBooleanValue(false))
	parserCtor := tsNodeHTTPCommonGet(loop, state.common, "HTTPParser"); parser := tsConstruct(loop, parserCtor); state.parser = parser
	requestType := tsNodeHTTP1ServerParserSlot(loop, parserCtor, "REQUEST"); maxSize := tsNodeHTTP1ServerGet(loop, server, "maxHeaderSize"); if !tsTruthy(maxSize) { maxSize = tsNumberValue(0) }
	flagsFn := tsNodeHTTPCommonGet(loop, state.common, "calculateLenientFlags"); flags := tsCall(loop, flagsFn, tsNodeHTTPServerGet(loop, server, "httpValidation"), tsNodeHTTPServerGet(loop, server, "insecureHTTPParser"))
	tsNodeHTTPCommonCall(loop, parser, "initialize", requestType, tsObjectValue(tsNewObject()), maxSize, flags); tsNodeHTTPCommonSet(loop, parser, "socket", socket); tsNodeHTTP1ServerSet(loop, socket, "parser", parser)
	maxHeaders := tsNodeHTTPServerGet(loop, server, "maxHeadersCount"); if tsIsNumeric(maxHeaders) { tsNodeHTTPServerSet(loop, parser, "maxHeaderPairs", tsNumberValue(tsNumber(maxHeaders)*2)) } else { tsNodeHTTPServerSet(loop, parser, "maxHeaderPairs", tsNodeHTTPCommonGet(loop, state.common, "MAX_HEADER_PAIRS")) }
	for _, entry := range []struct{slot, name string; length int}{{"kOnMessageBegin","onHttp1MessageBegin",0},{"kOnHeaders","onHttp1Headers",2},{"kOnHeadersComplete","onHttp1HeadersComplete",9},{"kOnBody","onHttp1Body",2},{"kOnMessageComplete","onHttp1MessageComplete",0}} { tsSet(loop, parser, tsNodeHTTP1ServerParserSlot(loop, parserCtor, entry.slot), tsNodeHTTP1ServerCallback(loop,state, entry.name, entry.length)) }
	for _, entry := range []struct{event,name string; once bool}{{"data","onHttp1SocketData",false},{"error","onHttp1SocketErrorListener",false},{"end","onHttp1SocketEnd",true},{"drain","onHttp1SocketDrain",false},{"finish","onHttp1SocketFinish",true},{"close","onHttp1SocketClose",true}} {
		method := "on"; if entry.once { method = "once" }; tsNodeHTTPServerCall(loop, socket, method, tsNodeHTTP1ServerKey(entry.event), tsNodeHTTP1ServerCallback(loop,state, entry.name, 1))
	}
	tsNodeHTTPServerCall(loop, socket, "prependOnceListener", tsNodeHTTP1ServerKey("close"), tsNodeHTTP1ServerCallback(loop,state,"onHttp1SocketCloseFlush",0))
	_ = internal
	return tsU
}

func tsNodeHTTP1ServerConnectionDispatch(loop *tsLoop,state *tsNodeHTTP1Connection, name string, args []tsValue) tsValue { arg := func(index int) tsValue { return tsArg(args,index) }; server,socket,parser := state.server,state.socket,state.parser
	switch name {
	case "onHttp1MessageBegin": tsNodeHTTP1ServerPut(loop,socket,"kHttp1MessageOpen",tsBooleanValue(true)); return tsU
	case "onHttp1Headers":
		headers := tsNodeHTTPIncomingArray(arg(0)); state.flushedHeaders = append(state.flushedHeaders, headers.values...); state.flushedURL += tsText(arg(1)); return tsU
	case "onHttp1HeadersComplete":
		major,minor,headers,method,url,upgrade,keepAlive:=arg(0),arg(1),arg(2),arg(3),arg(4),arg(7),arg(8)
		if tsIsUndefined(headers) { headers=tsArrayValue(&tsArray{values:append([]tsValue(nil),state.flushedHeaders...)});state.flushedHeaders=nil }
		if tsIsUndefined(url) { url=tsNodeHTTP1ServerKey(state.flushedURL);state.flushedURL="" }
		active:=tsNumber(tsNodeHTTP1ServerSym(loop,socket,"kHttp1ActiveRequests"));tsNodeHTTP1ServerPut(loop,socket,"kHttp1ActiveRequests",tsNumberValue(active+1))
		incomingCtor:=tsNodeHTTPServerGet(loop,state.http,"IncomingMessage");http1Options:=tsNodeHTTP1ServerGet(loop,state.options,"http1Options");if custom:=tsNodeHTTP1ServerGet(loop,http1Options,"IncomingMessage");!tsNullish(custom){incomingCtor=custom};request:=tsConstruct(loop,incomingCtor,socket);state.req=request;tsNodeHTTPServerPut(loop,socket,"kRequest",request)
		for key,value:=range map[string]tsValue{"socket":socket,"httpVersionMajor":major,"httpVersionMinor":minor,"httpVersion":tsNodeHTTP1ServerKey(tsText(major)+"."+tsText(minor)),"url":url,"upgrade":upgrade}{tsNodeHTTPIncomingSet(loop,request,key,value)}
		requestKeys:=tsNodeHTTPIncomingModule(loop);reqKeepAlive:=tsNodeHTTPIncomingGet(loop,requestKeys,"kReqShouldKeepAlive");if !tsIsUndefined(reqKeepAlive){tsSet(loop,request,reqKeepAlive,keepAlive)}
		methods:=tsNodeHTTPCommonGet(loop,state.common,"allMethods");if tsIsNumeric(method){index:=int(tsNumber(method));if methods.kind==tsArrayKind&&index>=0&&index<len((*tsArray)(methods.ref).values){tsNodeHTTPIncomingSet(loop,request,"method",(*tsArray)(methods.ref).values[index])}}else{tsNodeHTTPIncomingSet(loop,request,"method",method)}
		if tsTruthy(tsNodeHTTPServerGet(loop,server,"joinDuplicateHeaders")){tsNodeHTTPIncomingSet(loop,request,"joinDuplicateHeaders",tsBooleanValue(true))};tsNodeHTTPIncomingCall(loop,request,"_addHeaderLines",headers,tsNumberValue(float64(tsNodeHTTPIncomingArray(headers).length())))
		if tsTruthy(upgrade){shouldUpgrade:=tsText(tsNodeHTTPIncomingGet(loop,request,"method"))=="CONNECT";callback:=tsNodeHTTPServerGet(loop,server,"shouldUpgradeCallback");if callback.kind==tsFunctionKind{shouldUpgrade=shouldUpgrade||tsTruthy(tsCall(loop,callback,request))}else{shouldUpgrade=shouldUpgrade||tsNumber(tsNodeHTTPServerCall(loop,server,"listenerCount",tsNodeHTTP1ServerKey("upgrade")))>0};tsNodeHTTPIncomingSet(loop,request,"upgrade",tsBooleanValue(shouldUpgrade));if shouldUpgrade{state.pendingUpgrade=request;return tsNumberValue(2)}}
		readFn:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{tsNodeHTTPIncomingSet(loop,request,"_consuming",tsBooleanValue(true));if !tsTruthy(tsNodeHTTP1ServerGet(loop,socket,"_paused"))&&tsTruthy(tsNodeHTTP1ServerGet(loop,socket,"readable")){tsNodeHTTPServerCall(loop,socket,"resume")};return tsU}));tsNodeHTTPIncomingSet(loop,request,"_read",readFn)
		responseCtor:=tsNodeHTTPServerGet(loop,state.http,"ServerResponse");if custom:=tsNodeHTTP1ServerGet(loop,http1Options,"ServerResponse");!tsNullish(custom){responseCtor=custom};response:=tsConstruct(loop,responseCtor,request);tsNodeHTTPServerSet(loop,response,"_keepAliveTimeout",state.keepAliveTimeout);maxPerSocket:=tsNodeHTTPServerGet(loop,server,"maxRequestsPerSocket");tsNodeHTTPServerSet(loop,response,"_maxRequestsPerSocket",maxPerSocket)
		handle:=tsNodeHTTP1ServerCreateResponseHandle(loop,socket,keepAlive,state.keepAliveTimeout)
		finishedFn:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{n:=tsNumber(tsNodeHTTP1ServerSym(loop,socket,"kHttp1ActiveRequests"))-1;if n<0{n=0};tsNodeHTTP1ServerPut(loop,socket,"kHttp1ActiveRequests",tsNumberValue(n));if !tsTruthy(keepAlive)&&!tsTruthy(tsNodeHTTP1ServerGet(loop,socket,"destroyed")){tsNodeHTTPServerCall(loop,socket,"end")};return tsU}));tsNodeHTTPServerSet(loop,handle,"onfinished",finishedFn)
		flushedFn:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{tsNodeHTTP1ServerCallPipeline(loop,state.pipeline,"finishDrainedResponse",response);return tsU}));tsNodeHTTPServerSet(loop,handle,"onflushed",flushedFn);tsSet(loop,response,tsInternalHTTPGet(loop,tsInternalHTTPModule(loop),"kHandle"),handle)
		if !tsNullish(tsNodeHTTP1ServerGet(loop,socket,"_httpMessage")){tsNodeHTTP1ServerCallPipeline(loop,state.pipeline,"queuePipelinedResponse",socket,response,tsBooleanValue(tsNumber(major)<1||tsNumber(minor)<1))}else{tsNodeHTTPServerCall(loop,response,"assignSocket",socket)}
		finishListener:=tsFunctionValue(tsFunc(func(loop *tsLoop,values ...tsValue)tsValue{responseHandle:=tsGet(loop,response,tsInternalHTTPGet(loop,tsInternalHTTPModule(loop),"kHandle"));tsNodeHTTPServerSet(loop,responseHandle,"ended",tsBooleanValue(true));tsNodeHTTPServerCall(loop,responseHandle,"flushed");tsNodeHTTPServerCall(loop,responseHandle,"onfinished");tsNodeHTTPServerResponseCloseOnFinish(loop,response,socket);finishedRequest:=tsNodeHTTPServerGet(loop,response,"req");if !tsTruthy(tsNodeHTTPIncomingGet(loop,finishedRequest,"_consuming"))&&!tsTruthy(tsNodeHTTPServerGet(loop,tsNodeHTTPServerGet(loop,finishedRequest,"_readableState"),"resumeScheduled")){tsNodeHTTPIncomingCall(loop,finishedRequest,"_dump")};tsNodeHTTPServerCall(loop,response,"detachSocket",socket);closeKey:=tsNodeHTTP1ServerGet(loop,state.pipeline,"kMustCloseConnection");if tsTruthy(tsGet(loop,response,closeKey))||tsTruthy(tsNodeHTTPServerGet(loop,response,"_last")){destroySoon:=tsNodeHTTPServerGet(loop,socket,"destroySoon");if destroySoon.kind==tsFunctionKind{tsCallReceiver(loop,destroySoon,socket)}else if !tsTruthy(tsNodeHTTPServerGet(loop,socket,"writableEnded")){tsNodeHTTPServerCall(loop,socket,"end")};return tsU};tsNodeHTTP1ServerCallPipeline(loop,state.pipeline,"advanceResponsePipeline",server,socket);return tsU}))
		tsNodeHTTPServerCall(loop,response,"on",tsNodeHTTP1ServerKey("finish"),finishListener);tsNodeHTTP1ServerCallPipeline(loop,state.pipeline,"maybePauseHTTP1Reads",socket)
		if tsNumber(major)==1&&tsNumber(minor)==1&&tsIsNumeric(maxPerSocket)&&tsNumber(maxPerSocket)>0{state.requestsCount++;tsNodeHTTPServerSet(loop,response,"maxRequestsOnConnectionReached",tsBooleanValue(tsNumber(maxPerSocket)<=float64(state.requestsCount)));if tsNumber(maxPerSocket)<float64(state.requestsCount){tsNodeHTTPServerEmit(loop,server,"dropRequest",request,socket);tsNodeHTTPServerCall(loop,response,"writeHead",tsNumberValue(503));tsNodeHTTPServerCall(loop,response,"end");return tsNumberValue(0)}}
		if tsNumber(major)==1&&tsNumber(minor)==1&&!tsStrictEqual(tsNodeHTTPServerGet(loop,server,"requireHostHeader"),tsBooleanValue(false))&&tsIsUndefined(tsNodeHTTPIncomingGet(loop,tsNodeHTTPIncomingGet(loop,request,"headers"),"host")){tsNodeHTTPServerCall(loop,response,"writeHead",tsNumberValue(400),tsArrayValue(&tsArray{values:[]tsValue{tsNodeHTTP1ServerKey("Connection"),tsNodeHTTP1ServerKey("close")}}));tsNodeHTTPServerCall(loop,response,"end");return tsNumberValue(0)}
		expect:=tsNodeHTTPIncomingGet(loop,request,"headers");expect=tsGet(loop,expect,tsNodeHTTP1ServerKey("expect"));if !tsIsUndefined(expect)&&tsNumber(major)==1&&tsNumber(minor)==1{matches:=tsNodeHTTPCommonHasBoundaryToken(tsText(expect),"100-continue");if matches{if tsNumber(tsNodeHTTPServerCall(loop,server,"listenerCount",tsNodeHTTP1ServerKey("checkContinue")))>0{tsNodeHTTPServerEmit(loop,server,"checkContinue",request,response)}else{tsNodeHTTPServerCall(loop,response,"writeContinue");tsNodeHTTPServerEmit(loop,server,"request",request,response)}}else if tsNumber(tsNodeHTTPServerCall(loop,server,"listenerCount",tsNodeHTTP1ServerKey("checkExpectation")))>0{tsNodeHTTPServerEmit(loop,server,"checkExpectation",request,response)}else{tsNodeHTTPServerCall(loop,response,"writeHead",tsNumberValue(417));tsNodeHTTPServerCall(loop,response,"end")};return tsNumberValue(0)}
		tsNodeHTTPServerEmit(loop,server,"request",request,response);if tsText(tsNodeHTTPIncomingGet(loop,request,"method"))=="PRI"{state.closeAfterParse=true;return tsNumberValue(2)};return tsNumberValue(0)
	case "onHttp1Body": if !tsNullish(state.req)&&!tsTruthy(tsNodeHTTPIncomingGet(loop,state.req,"_dumped")){tsNodeHTTPIncomingCall(loop,state.req,"push",arg(0))};return tsU
	case "onHttp1MessageComplete":
		tsNodeHTTP1ServerPut(loop,socket,"kHttp1MessageOpen",tsBooleanValue(false));tsNodeHTTPServerPut(loop,socket,"kRequest",tsU);trailers:=state.flushedHeaders;if len(trailers)>0{state.flushedHeaders=nil};state.flushedURL="";if !tsNullish(state.req){tsNodeHTTPIncomingSet(loop,state.req,"complete",tsBooleanValue(true));tsNodeHTTPIncomingCall(loop,state.req,"_addHeaderLines",tsArrayValue(&tsArray{values:trailers}),tsNumberValue(float64(len(trailers))));tsNodeHTTPIncomingCall(loop,state.req,"push",tsNull)};return tsU
	case "onHttp1SocketError", "onHttp1SocketErrorListener":
		err:=arg(0);common:=state.common;prepare:=tsNodeHTTPCommonGet(loop,common,"prepareError");if prepare.kind==tsFunctionKind{tsCall(loop,prepare,err,parser,arg(1))};if !tsTruthy(tsNodeHTTPServerEmit(loop,server,"clientError",err,socket)){inFlight:=tsNodeHTTPServerGet(loop,socket,"_httpMessage");mayWriteError:=tsNullish(inFlight)||!tsTruthy(tsNodeHTTPServerGet(loop,inFlight,"_headerSent"));if mayWriteError&&tsTruthy(tsNodeHTTPServerGet(loop,socket,"writable"))&&!tsTruthy(tsNodeHTTPServerGet(loop,socket,"destroyed")){code:=tsText(tsNodeHTTPServerGet(loop,err,"code"));message:="HTTP/1.1 400 Bad Request\r\nConnection: close\r\n\r\n";if code=="HPE_HEADER_OVERFLOW"{message="HTTP/1.1 431 Request Header Fields Too Large\r\nConnection: close\r\n\r\n"}else if code=="HPE_CHUNK_EXTENSIONS_OVERFLOW"{message="HTTP/1.1 413 Payload Too Large\r\nConnection: close\r\n\r\n"}else if code=="ERR_HTTP_REQUEST_TIMEOUT"{message="HTTP/1.1 408 Request Timeout\r\nConnection: close\r\n\r\n"};tsNodeHTTPServerCall(loop,socket,"write",tsNodeHTTP1ServerKey(message),tsNodeHTTP1ServerKey("latin1"))};tsNodeHTTPServerCall(loop,socket,"destroy",err)};return tsU
	case "onHttp1SocketData":
		data:=arg(0);ret:=tsNodeHTTPCommonCall(loop,parser,"execute",data);if ret.kind==tsErrorKind{tsNodeHTTP1ServerConnectionDispatch(loop,state,"onHttp1SocketError",[]tsValue{ret,data});return tsU};if state.closeAfterParse{state.closeAfterParse=false;tsNodeHTTPServerCall(loop,socket,"destroy");return tsU};if !tsNullish(state.pendingUpgrade){request:=state.pendingUpgrade;state.pendingUpgrade=tsNull;for _,event:=range []struct{name,callback string}{{"data","onHttp1SocketData"},{"error","onHttp1SocketErrorListener"},{"end","onHttp1SocketEnd"},{"drain","onHttp1SocketDrain"}}{tsNodeHTTPServerCall(loop,socket,"removeListener",tsNodeHTTP1ServerKey(event.name),tsNodeHTTP1ServerCallback(loop,state,event.callback,1))};deleteFn:=tsGet(loop,state.connections,tsNodeHTTP1ServerKey("delete"));if deleteFn.kind==tsFunctionKind{tsCallReceiver(loop,deleteFn,state.connections,socket)};tsNodeHTTPServerPut(loop,socket,"kHandedOff",tsBooleanValue(true));tracked:=tsNodeHTTPServerSym(loop,state.server,"kTrackedConnections");removeTracked:=tsGet(loop,tracked,tsNodeHTTP1ServerKey("delete"));if removeTracked.kind==tsFunctionKind{tsCallReceiver(loop,removeTracked,tracked,socket)};tsNodeHTTPCommonFreeParser(loop,parser,request,socket);eventName:="upgrade";if tsText(tsNodeHTTPIncomingGet(loop,request,"method"))=="CONNECT"{eventName="connect"};raw:=tsNodeBufferVisibleBytes(data);parsed:=int(tsNumber(ret));if parsed<0{parsed=0};if parsed>len(raw){parsed=len(raw)};body:=tsNodeBuffer(append([]byte(nil),raw[parsed:]...));if tsNumber(tsNodeHTTPServerCall(loop,server,"listenerCount",tsNodeHTTP1ServerKey(eventName)))>0{tsNodeHTTPServerSet(loop,socket,"readableFlowing",tsNull);tsNodeHTTPServerEmit(loop,server,eventName,request,socket,body)}else{tsNodeHTTPServerCall(loop,socket,"destroy")}};return tsU
	case "onHttp1SocketEnd":
		result:=tsNodeHTTPCommonCall(loop,parser,"finish");if result.kind==tsErrorKind{tsNodeHTTP1ServerConnectionDispatch(loop,state,"onHttp1SocketError",[]tsValue{result,tsU});return tsU};if !tsTruthy(tsNodeHTTPServerGet(loop,server,"httpAllowHalfOpen")){if !tsNullish(state.req)&&!tsTruthy(tsNodeHTTPIncomingGet(loop,state.req,"complete")){tsNodeHTTPIncomingCall(loop,state.req,"destroy")};if tsTruthy(tsNodeHTTPServerGet(loop,socket,"writable")){tsNodeHTTPServerCall(loop,socket,"end")};return tsU};last:=tsNodeHTTP1ServerCallPipeline(loop,state.pipeline,"lastPipelinedResponse",socket);message:=last;if tsNullish(message){message=tsNodeHTTPServerGet(loop,socket,"_httpMessage")};if !tsNullish(message){tsNodeHTTPServerSet(loop,message,"_last",tsBooleanValue(true))}else if tsTruthy(tsNodeHTTPServerGet(loop,socket,"writable")){tsNodeHTTPServerCall(loop,socket,"end")};return tsU
	case "onHttp1SocketDrain":
		tsNodeHTTP1ServerCallPipeline(loop,state.pipeline,"resumeHTTP1ReadsOnDrain",socket);message:=tsNodeHTTPServerGet(loop,socket,"_httpMessage");if !tsNullish(message){handle:=tsGet(loop,message,tsInternalHTTPGet(loop,tsInternalHTTPModule(loop),"kHandle"));tsNodeHTTPServerCall(loop,handle,"socketDrained")};return tsU
	case "onHttp1SocketFinish":
		message:=tsNodeHTTPServerGet(loop,socket,"_httpMessage");if !tsNullish(message){handle:=tsGet(loop,message,tsInternalHTTPGet(loop,tsInternalHTTPModule(loop),"kHandle"));callback:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{tsNodeHTTPServerCall(loop,message,"_callPendingCallbacks");return tsU}));tsNodeHTTPServerCall(loop,handle,"socketEnded",callback)};return tsU
	case "onHttp1SocketCloseFlush": message:=tsNodeHTTPServerGet(loop,socket,"_httpMessage");if !tsNullish(message){tsNodeHTTPServerCall(loop,tsGet(loop,message,tsInternalHTTPGet(loop,tsInternalHTTPModule(loop),"kHandle")),"flushed")};return tsU
	case "onHttp1SocketClose":
		deleteFn:=tsGet(loop,state.connections,tsNodeHTTP1ServerKey("delete"));if deleteFn.kind==tsFunctionKind{tsCallReceiver(loop,deleteFn,state.connections,socket)};message:=tsNodeHTTPServerGet(loop,socket,"_httpMessage");if !tsNullish(message){handle:=tsGet(loop,message,tsInternalHTTPGet(loop,tsInternalHTTPModule(loop),"kHandle"));err:=tsNodeHTTPServerGet(loop,socket,"errored");if tsNullish(err){err=tsNodeHTTPServerGet(loop,message,"errored")};if tsNullish(err){err=tsFSException("Error","ERR_STREAM_DESTROYED","Cannot call write after a stream was destroyed")};callback:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{tsNodeHTTP1ServerFailPending(loop,message,err);return tsU}));tsNodeHTTPServerCall(loop,handle,"socketEnded",callback);request:=tsNodeHTTPServerGet(loop,message,"req");if !tsNullish(request)&&!tsTruthy(tsNodeHTTPIncomingGet(loop,request,"destroyed")){if tsNumber(tsNodeHTTPServerCall(loop,request,"listenerCount",tsNodeHTTP1ServerKey("error")))>0{tsNodeHTTPIncomingCall(loop,request,"destroy",tsFSException("Error","ECONNRESET","aborted"))}else{tsNodeHTTPIncomingCall(loop,request,"destroy")}}};tsNodeHTTP1ServerCallPipeline(loop,state.pipeline,"abortQueuedPipelinedResponses",socket);tsNodeHTTPCommonFreeParser(loop,parser,state.req,socket);return tsU
	}
	return tsU
}

func tsNodeHTTP1ServerCallPipeline(loop *tsLoop, pipeline tsValue, name string, args ...tsValue) tsValue { method:=tsGet(loop,pipeline,tsNodeHTTP1ServerKey(name));if method.kind==tsFunctionKind{return tsCall(loop,method,args...)};return tsU }
func tsNodeHTTP1ServerCloseIdle(loop *tsLoop, server tsValue) tsValue {
	tsNodeHTTP1ServerInit(loop);connections:=tsNodeHTTP1ServerSym(loop,server,"kHttp1Connections");if connections.kind!=tsCollectionKind{return tsU}
	for _,entry:=range (*tsCollection)(connections.ref).entries{if !entry.live{continue};socket:=entry.key;if !tsTruthy(tsNodeHTTP1ServerSym(loop,socket,"kHttp1ActiveRequests"))&&!tsTruthy(tsNodeHTTP1ServerSym(loop,socket,"kHttp1MessageOpen"))&&!tsTruthy(tsNodeHTTP1ServerGet(loop,socket,"destroyed")){tsNodeHTTPServerCall(loop,socket,"destroy")}}
	return tsU
}
func tsNodeHTTP1ServerCloseAll(loop *tsLoop, server tsValue) tsValue {
	tsNodeHTTP1ServerInit(loop);connections:=tsNodeHTTP1ServerSym(loop,server,"kHttp1Connections");if connections.kind!=tsCollectionKind{return tsU}
	for _,entry:=range (*tsCollection)(connections.ref).entries{if entry.live&&!tsTruthy(tsNodeHTTP1ServerGet(loop,entry.key,"destroyed")){tsNodeHTTPServerCall(loop,entry.key,"destroy")}}
	return tsU
}
func tsNodeHTTP1ServerModule(loop *tsLoop) tsValue {
	if !tsIsUndefined(tsNodeHTTP1ServerSingleton.module){return tsNodeHTTP1ServerSingleton.module}
	tsNodeHTTP1ServerInit(loop);module:=tsNewObject()
	module.set("connectionListenerHTTP1",tsNamedNativeMethod("connectionListenerHTTP1",3,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeHTTP1ServerConnectionListener(loop,tsArg(args,0),tsArg(args,1),tsArg(args,2))}))
	module.set("closeIdleHttp1Connections",tsNamedNativeMethod("closeIdleHttp1Connections",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeHTTP1ServerCloseIdle(loop,tsArg(args,0))}))
	module.set("closeAllHttp1Connections",tsNamedNativeMethod("closeAllHttp1Connections",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeHTTP1ServerCloseAll(loop,tsArg(args,0))}))
	module.set("kHttp1Connections",tsNodeHTTP1ServerSingleton.kConnections);module.set("default",tsObjectValue(module));tsNodeHTTP1ServerSingleton.module=tsObjectValue(module);return tsNodeHTTP1ServerSingleton.module
}
`
