package goemit

// nodeHTTPIncomingRuntime implements IncomingMessage on the local Readable
// stream and preserves Node's lazy raw/header/trailer accessors.
const nodeHTTPIncomingRuntime = `
var tsNodeHTTPIncomingSingleton struct{module,constructor tsValue;prototype *tsObject;keys map[string]tsValue}
func tsNodeHTTPIncomingKey(s string)tsValue{return tsStringReference(tsStringUTF8(s))}
func tsNodeHTTPIncomingGet(loop *tsLoop,v tsValue,s string)tsValue{return tsGet(loop,v,tsNodeHTTPIncomingKey(s))}
func tsNodeHTTPIncomingSet(loop *tsLoop,v tsValue,s string,x tsValue){tsSet(loop,v,tsNodeHTTPIncomingKey(s),x)}
func tsNodeHTTPIncomingCall(loop *tsLoop,v tsValue,s string,args ...tsValue)tsValue{f:=tsNodeHTTPIncomingGet(loop,v,s);if f.kind==tsFunctionKind{return tsCallReceiver(loop,f,v,args...)};return tsU}
func tsNodeHTTPIncomingArray(v tsValue)*tsArray{if v.kind==tsArrayKind{return (*tsArray)(v.ref)};return &tsArray{}}
func tsNodeHTTPIncomingField(field string)(string,int){key:=strings.ToLower(field);switch key{case "set-cookie":return key,1;case "cookie":return key,2;case "age","authorization","content-length","content-type","etag","expires","from","host","if-modified-since","if-unmodified-since","location","max-forwards","proxy-authorization","referer","retry-after","server","status","user-agent","last-modified":return key,3;case "date","vary","origin","expect","accept","upgrade","if-match","connection","cache-control","if-none-match","accept-encoding","accept-language","x-forwarded-for","content-encoding","x-forwarded-host","transfer-encoding","x-forwarded-proto":return key,0;default:return key,0}}
func tsNodeHTTPIncomingAdd(loop *tsLoop,msg,field,value,dest tsValue){name,kind:=tsNodeHTTPIncomingField(tsText(field));old:=tsGet(loop,dest,tsNodeHTTPIncomingKey(name));if kind==1{if old.kind!=tsArrayKind{values:=[]tsValue{};if !tsIsUndefined(old){values=append(values,old)};values=append(values,value);tsSet(loop,dest,tsNodeHTTPIncomingKey(name),tsArrayValue(&tsArray{values:values}))}else{a:=(*tsArray)(old.ref);a.values=append(a.values,value)};return};if kind==0||kind==2{if old.kind==tsStringKind{sep:=", ";if kind==2{sep="; "};tsSet(loop,dest,tsNodeHTTPIncomingKey(name),tsStringReference(tsStringUTF8(tsText(old)+sep+tsText(value))))}else{tsSet(loop,dest,tsNodeHTTPIncomingKey(name),value)};return};if tsTruthy(tsNodeHTTPIncomingGet(loop,msg,"joinDuplicateHeaders"))&&!tsIsUndefined(old){tsSet(loop,dest,tsNodeHTTPIncomingKey(name),tsStringReference(tsStringUTF8(tsText(old)+", "+tsText(value))))}else if tsIsUndefined(old){tsSet(loop,dest,tsNodeHTTPIncomingKey(name),value)}}
func tsNodeHTTPIncomingDistinct(loop *tsLoop,field,value,dest tsValue){name:=tsNodeHTTPIncomingKey(strings.ToLower(tsText(field)));old:=tsGet(loop,dest,name);if old.kind!=tsArrayKind{tsSet(loop,dest,name,tsArrayValue(&tsArray{values:[]tsValue{value}}))}else{a:=(*tsArray)(old.ref);a.values=append(a.values,value)}}
func tsNodeHTTPIncomingRaw(loop *tsLoop,msg tsValue)tsValue{return tsNodeHTTPIncomingRawHeaders(loop,msg)}
func tsNodeHTTPIncomingHeaders(loop *tsLoop, msg tsValue, name string, distinct bool) tsValue {
	keys := tsNodeHTTPIncomingSingleton.keys
	trailers := name == "trailers" || name == "trailersDistinct"
	cache, count := keys["headers"], keys["headerCount"]
	if trailers { cache, count = keys["trailers"], keys["trailerCount"] }
	if distinct {
		if trailers { cache = keys["trailersDistinct"] } else { cache = keys["headersDistinct"] }
	}
	if value := tsGet(loop, msg, cache); !tsNullish(value) { return value }
	// The Bun native handle can produce the same folded headers object without
	// allocating the intermediate raw header array.
	if !trailers && !distinct && tsIsNull(tsGet(loop, msg, keys["raw"])) {
		if headerSource := tsGet(loop, msg, keys["headerSource"]); !tsNullish(headerSource) {
			method := tsNodeHTTPIncomingGet(loop, headerSource, "headersObject")
			adder := tsNodeHTTPIncomingGet(loop, msg, "_addHeaderLine")
			defaultAdder := tsGet(loop, tsObjectValue(tsNodeHTTPIncomingSingleton.prototype), tsNodeHTTPIncomingKey("_addHeaderLine"))
			if method.kind == tsFunctionKind && tsStrictEqual(adder, defaultAdder) {
				if built := tsCallReceiver(loop, method, headerSource); !tsIsUndefined(built) {
					tsSet(loop, msg, cache, built)
					return built
				}
			}
		}
	}
	source := tsNodeHTTPIncomingRaw(loop, msg)
	if trailers { source = tsNodeHTTPIncomingGet(loop, msg, "rawTrailers") }
	out := tsNewObject()
	out.prototype = tsNull
	values := tsNodeHTTPIncomingArray(source).values
	n := int(tsNumber(tsGet(loop, msg, count)))
	if n > len(values) { n = len(values) }
	for i := 0; i+1 < n; i += 2 {
		if distinct { tsNodeHTTPIncomingDistinct(loop, values[i], values[i+1], tsObjectValue(out)) } else { tsNodeHTTPIncomingAdd(loop, msg, values[i], values[i+1], tsObjectValue(out)) }
	}
	result := tsObjectValue(out)
	tsSet(loop, msg, cache, result)
	return result
}
func tsNodeHTTPIncomingCreate(loop *tsLoop, args []tsValue) tsValue {
	tsNodeHTTPIncomingInit(loop)
	keys := tsNodeHTTPIncomingSingleton.keys
	internal := tsInternalHTTPModule(loop)
	handleKey := tsInternalHTTPGet(loop, internal, "kHandle")
	native := tsStrictEqual(tsArg(args, 0), handleKey)
	socket := tsArg(args, 0)
	if native { socket = tsArg(args, 7) }
	options := tsNewObject()
	if !tsNullish(socket) {
		if hwm := tsNodeHTTPIncomingGet(loop, socket, "readableHighWaterMark"); !tsIsUndefined(hwm) { options.set("highWaterMark", hwm) }
	}
	base := tsNewObject()
	base.prototype = tsObjectValue(tsNodeHTTPIncomingSingleton.prototype)
	stream := tsNodeReadableCreate(loop, tsObjectValue(base), tsObjectValue(options))
	(*tsObject)(stream.ref).prototype = tsObjectValue(tsNodeHTTPIncomingSingleton.prototype)
	for key, value := range map[string]tsValue{"aborted":tsBooleanValue(false), "eof":tsBooleanValue(false), "headers":tsNull, "headersDistinct":tsNull, "raw":tsNull, "headerCount":tsNumberValue(0), "trailers":tsNull, "trailersDistinct":tsNull, "trailerCount":tsNumberValue(0), "abortController":tsNull, "reqKeepAlive":tsU, "headerSource":tsNull} {
		tsSet(loop, stream, keys[key], value)
	}
	tsNodeHTTPIncomingSet(loop, stream, "_consuming", tsBooleanValue(false))
	tsNodeHTTPIncomingSet(loop, stream, "_dumped", tsBooleanValue(false))
	tsNodeHTTPIncomingSet(loop, stream, "complete", tsBooleanValue(false))
	tsNodeHTTPIncomingSet(loop, stream, "_closed", tsBooleanValue(false))
	tsNodeHTTPIncomingSet(loop, stream, "rawTrailers", tsArrayValue(&tsArray{}))
	if native {
		tsNodeHTTPIncomingSet(loop, stream, "url", tsArg(args, 1))
		tsNodeHTTPIncomingSet(loop, stream, "method", tsArg(args, 2))
		tsSet(loop, stream, keys["headerSource"], tsArg(args, 5))
		tsSet(loop, stream, keys["handle"], tsArg(args, 5))
		tsSet(loop, stream, tsNodeHTTPIncomingGet(loop, internal, "noBodySymbol"), tsBooleanValue(!tsTruthy(tsArg(args, 6))))
		tsNodeHTTPIncomingSet(loop, stream, "socket", socket)
		tsNodeHTTPIncomingSet(loop, stream, "client", socket)
		tsNodeHTTPIncomingSet(loop, stream, "upgrade", tsNull)
	} else {
		for _, name := range []string{"httpVersionMajor", "httpVersionMinor", "httpVersion", "statusCode", "statusMessage"} { tsNodeHTTPIncomingSet(loop, stream, name, tsNull) }
		tsNodeHTTPIncomingSet(loop, stream, "socket", socket)
		tsNodeHTTPIncomingSet(loop, stream, "client", socket)
		tsNodeHTTPIncomingSet(loop, stream, "url", tsNodeHTTPIncomingKey(""))
		tsNodeHTTPIncomingSet(loop, stream, "method", tsNull)
		tsNodeHTTPIncomingSet(loop, stream, "upgrade", tsNull)
		tsNodeHTTPIncomingSet(loop, stream, "joinDuplicateHeaders", tsBooleanValue(false))
		tsNodeHTTPIncomingSet(loop, stream, "rawHeaders", tsArrayValue(&tsArray{}))
		tsSet(loop, stream, tsNodeHTTPIncomingGet(loop, internal, "noBodySymbol"), tsBooleanValue(false))
	}
	state := tsNodeHTTPIncomingGet(loop, stream, "_readableState")
	tsNodeHTTPIncomingSet(loop, state, "readingMore", tsBooleanValue(true))
	return stream
}
func tsNodeHTTPIncomingAddLines(loop *tsLoop, receiver, headers tsValue, count int) tsValue {
	if headers.kind != tsArrayKind || len((*tsArray)(headers.ref).values) == 0 { return tsU }
	keys := tsNodeHTTPIncomingSingleton.keys
	complete := tsTruthy(tsNodeHTTPIncomingGet(loop, receiver, "complete"))
	cache, countKey := keys["headers"], keys["headerCount"]
	if complete {
		tsNodeHTTPIncomingSet(loop, receiver, "rawTrailers", headers)
		cache, countKey = keys["trailers"], keys["trailerCount"]
	} else {
		tsSet(loop, receiver, keys["raw"], headers)
	}
	tsSet(loop, receiver, countKey, tsNumberValue(float64(count)))
	dest := tsGet(loop, receiver, cache)
	if tsTruthy(dest) {
		values := (*tsArray)(headers.ref).values
		for i := 0; i < count && i+1 < len(values); i += 2 {
			tsNodeHTTPIncomingCall(loop, receiver, "_addHeaderLine", values[i], values[i+1], dest)
		}
	}
	return tsU
}

func tsNodeHTTPIncomingDumpAndClose(loop *tsLoop, receiver tsValue) tsValue {
	tsNodeHTTPIncomingSet(loop, receiver, "_dumped", tsBooleanValue(true))
	state := tsNodeReadableStateFor(receiver)
	if state != nil {
		state.ended, state.endEmitted, state.destroyed, state.closed = true, true, true, true
		tsNodeReadableSyncState(state)
	}
	stateObject := tsNodeHTTPIncomingGet(loop, receiver, "_readableState")
	for _, name := range []string{"ended", "endEmitted", "destroyed", "closed", "closeEmitted"} { tsNodeHTTPIncomingSet(loop, stateObject, name, tsBooleanValue(true)) }
	return tsU
}

func tsNodeHTTPIncomingRead(loop *tsLoop, receiver tsValue) tsValue {
	internal := tsInternalHTTPModule(loop)
	keys := tsNodeHTTPIncomingSingleton.keys
	if !tsTruthy(tsNodeHTTPIncomingGet(loop, receiver, "_consuming")) {
		tsNodeHTTPIncomingSet(loop, tsNodeHTTPIncomingGet(loop, receiver, "_readableState"), "readingMore", tsBooleanValue(false))
		tsNodeHTTPIncomingSet(loop, receiver, "_consuming", tsBooleanValue(true))
	}
	socket := tsNodeHTTPIncomingGet(loop, receiver, "socket")
	handle := tsGet(loop, receiver, keys["handle"])
	if tsNullish(handle) {
		if !tsNullish(socket) && tsTruthy(tsNodeHTTPIncomingGet(loop, socket, "readable")) { tsNodeHTTPIncomingReadStart(loop, socket) }
		return tsU
	}
	if !tsNullish(socket) && tsTruthy(tsNodeHTTPIncomingGet(loop, socket, "readable")) {
		handoff := tsGet(loop, socket, keys["handoff"])
		if tsTruthy(tsNodeHTTPIncomingGet(loop, receiver, "upgrade")) || !tsIsUndefined(handoff) {
			tsNodeHTTPIncomingCall(loop, receiver, "onIncomingMessageResumeNodeHTTPResponse")
		} else { tsNodeHTTPIncomingReadStart(loop, socket) }
	}
	if tsTruthy(tsGet(loop, receiver, keys["eof"])) { return tsU }
	if tsTruthy(tsGet(loop, receiver, keys["noBody"])) {
		tsCall(loop, tsInternalHTTPGet(loop, internal, "emitEOFIncomingMessage"), receiver)
		return tsU
	}
	body := int(tsNumber(tsNodeHTTPIncomingGet(loop, handle, "hasBody")))
	states := tsInternalHTTPGet(loop, internal, "NodeHTTPBodyReadState")
	done := int(tsNumber(tsNodeHTTPIncomingGet(loop, states, "done")))
	none := int(tsNumber(tsNodeHTTPIncomingGet(loop, states, "none")))
	if body&done != 0 || body == none { tsCall(loop, tsInternalHTTPGet(loop, internal, "emitEOFIncomingMessage"), receiver) }
	if tsIsUndefined(tsNodeHTTPIncomingGet(loop, handle, "ondata")) {
		callback := tsFunctionValue(tsFunc(func(loop *tsLoop, args ...tsValue) tsValue {
			return tsCall(loop, tsInternalHTTPGet(loop, internal, "onDataIncomingMessage"), append([]tsValue{receiver}, args...)...)
		}))
		tsNodeHTTPIncomingSet(loop, handle, "ondata", callback)
	}
	return tsU
}

func tsNodeHTTPIncomingNextTickError(loop *tsLoop, receiver, err, callback tsValue) {
	next := tsFunctionValue(tsFunc(func(loop *tsLoop, _ ...tsValue) tsValue {
		return tsNodeHTTPIncomingOnError(loop, receiver, err, callback)
	}))
	tsInternalHTTPNextTick(loop, next)
}

func tsNodeHTTPIncomingOnError(loop *tsLoop, receiver, err, callback tsValue) tsValue {
	if callback.kind != tsFunctionKind { return tsU }
	count := tsNodeHTTPIncomingCall(loop, receiver, "listenerCount", tsNodeHTTPIncomingKey("error"))
	if tsNumber(count) == 0 { tsCall(loop, callback) } else { tsCall(loop, callback, err) }
	return tsU
}

func tsNodeHTTPIncomingFinishAfterDestroy(loop *tsLoop, receiver, cleanup, err, callback, finishedErr tsValue) tsValue {
	if !tsNullish(finishedErr) && tsNodeHTTPIncomingGet(loop, finishedErr, "code").kind == tsStringKind && tsText(tsNodeHTTPIncomingGet(loop, finishedErr, "code")) == "ERR_STREAM_PREMATURE_CLOSE" { finishedErr = tsNull }
	if cleanup.kind == tsFunctionKind { tsCall(loop, cleanup) }
	if tsTruthy(finishedErr) { err = finishedErr }
	tsNodeHTTPIncomingNextTickError(loop, receiver, err, callback)
	return tsU
}

func tsNodeHTTPIncomingDestroy(loop *tsLoop, receiver, err, callback tsValue) tsValue {
	keys := tsNodeHTTPIncomingSingleton.keys
	handle := tsGet(loop, receiver, keys["handle"])
	socket := tsNodeHTTPIncomingGet(loop, receiver, "socket")
	shouldAbort := !tsTruthy(tsNodeHTTPIncomingGet(loop, receiver, "readableEnded")) || !tsTruthy(tsNodeHTTPIncomingGet(loop, receiver, "complete"))
	if shouldAbort {
		tsSet(loop, receiver, keys["aborted"], tsBooleanValue(true))
		tsNodeHTTPIncomingCall(loop, receiver, "emit", tsNodeHTTPIncomingKey("aborted"))
	}
	if !tsNullish(handle) {
		tsSet(loop, receiver, keys["handle"], tsU)
		tsNodeHTTPIncomingSet(loop, handle, "onabort", tsU)
		tsNodeHTTPIncomingSet(loop, handle, "ondata", tsU)
		if shouldAbort && !tsTruthy(tsNodeHTTPIncomingGet(loop, handle, "finished")) && !tsNullish(socket) { tsNodeHTTPIncomingCall(loop, handle, "abort") }
		if shouldAbort && !tsNullish(socket) && !tsTruthy(tsNodeHTTPIncomingGet(loop, socket, "destroyed")) { tsNodeHTTPIncomingCall(loop, socket, "destroy", err) }
		if callback.kind == tsFunctionKind { tsCall(loop, tsInternalHTTPGet(loop, tsInternalHTTPModule(loop), "emitErrorNextTickIfErrorListenerNT"), receiver, err, callback) }
		return tsU
	}
	if !tsNullish(socket) && !tsTruthy(tsNodeHTTPIncomingGet(loop, socket, "destroyed")) && tsTruthy(tsNodeHTTPIncomingGet(loop, receiver, "aborted")) {
		tsNodeHTTPIncomingCall(loop, socket, "destroy", err)
		var cleanup tsValue
		finishedCallback := tsFunctionValue(tsFunc(func(loop *tsLoop, args ...tsValue) tsValue {
			return tsNodeHTTPIncomingFinishAfterDestroy(loop, receiver, cleanup, err, callback, tsArg(args, 0))
		}))
		cleanup = tsCall(loop, tsInternalStreamEndOfStreamModule(loop), socket, finishedCallback)
	} else { tsNodeHTTPIncomingNextTickError(loop, receiver, err, callback) }
	return tsU
}

func tsNodeHTTPIncomingInit(loop *tsLoop){if tsNodeHTTPIncomingSingleton.prototype!=nil{return};internal:=tsInternalHTTPModule(loop);keys:=map[string]tsValue{"aborted":tsNodeHTTPIncomingGet(loop,internal,"abortedSymbol"),"eof":tsNodeHTTPIncomingGet(loop,internal,"eofInProgress"),"handle":tsNodeHTTPIncomingGet(loop,internal,"kHandle"),"handoff":tsNodeHTTPIncomingGet(loop,internal,"kHandoffResponse"),"noBody":tsNodeHTTPIncomingGet(loop,internal,"noBodySymbol"),"abortController":tsNodeHTTPIncomingGet(loop,internal,"kAbortController")};for _,name:=range []string{"headers","headersDistinct","raw","headerCount","trailers","trailersDistinct","trailerCount","headerSource","reqKeepAlive"}{keys[name]=tsNewSymbol(tsStringUTF8(name))};tsNodeHTTPIncomingSingleton.keys=keys
	readable:=tsInternalStreamReadableModule(loop);ctor:=tsNodeHTTPIncomingGet(loop,readable,"Readable");p:=tsNewObject();p.prototype=tsGet(loop,ctor,tsNodeHTTPIncomingKey("prototype"));tsNodeHTTPIncomingSingleton.prototype=p
	fn:=tsFunc(func(_ *tsLoop,_ ...tsValue)tsValue{tsPropertyFailure("IncomingMessage must be constructed with new");return tsU});fn.name="IncomingMessage";fn.length=1;fn.constructible=true;fn.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeHTTPIncomingCreate(loop,args)};fn.receiverCall=func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeHTTPIncomingCreate(loop,append([]tsValue{receiver},args...))};fn.properties=tsNewObject();fn.properties.set("prototype",tsObjectValue(p));tsNodeHTTPIncomingSingleton.constructor=tsFunctionValue(fn);p.set("constructor",tsNodeHTTPIncomingSingleton.constructor)
	for name,value:=range map[string]tsValue{"httpVersion":tsNodeHTTPIncomingKey("1.1"),"httpVersionMajor":tsNumberValue(1),"httpVersionMinor":tsNumberValue(1),"url":tsNodeHTTPIncomingKey(""),"method":tsNull,"statusCode":tsNull,"statusMessage":tsNull,"upgrade":tsNull,"joinDuplicateHeaders":tsBooleanValue(false)}{p.set(name,value)}
	p.set("setTimeout",tsNamedNativeMethod("setTimeout",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{if cb:=tsArg(args,1);cb.kind==tsFunctionKind{tsNodeHTTPIncomingCall(loop,receiver,"on",tsNodeHTTPIncomingKey("timeout"),cb)};tsNodeHTTPIncomingCall(loop,tsNodeHTTPIncomingGet(loop,receiver,"socket"),"setTimeout",tsArg(args,0));return receiver}))
	p.set("onIncomingMessageResumeNodeHTTPResponse", tsNamedNativeMethod("onIncomingMessageResumeNodeHTTPResponse", 0, func(loop *tsLoop, receiver tsValue, _ ...tsValue) tsValue { handle := tsGet(loop, receiver, keys["handle"]); if !tsNullish(handle) && !tsTruthy(tsNodeHTTPIncomingGet(loop, receiver, "destroyed")) { tsNodeHTTPIncomingCall(loop, handle, "resume") }; return tsU }))
	p.set("_read", tsNamedNativeMethod("_read", 1, func(loop *tsLoop, receiver tsValue, _ ...tsValue) tsValue { return tsNodeHTTPIncomingRead(loop, receiver) }))
	p.set("_destroy", tsNamedNativeMethod("_destroy", 2, func(loop *tsLoop, receiver tsValue, args ...tsValue) tsValue { return tsNodeHTTPIncomingDestroy(loop, receiver, tsArg(args, 0), tsArg(args, 1)) }))
	p.set("_addHeaderLines", tsNamedNativeMethod("_addHeaderLines", 2, func(loop *tsLoop, receiver tsValue, args ...tsValue) tsValue { return tsNodeHTTPIncomingAddLines(loop, receiver, tsArg(args, 0), int(tsNumber(tsArg(args, 1)))) }))
	p.set("_addHeaderLine",tsNamedNativeMethod("_addHeaderLine",3,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{tsNodeHTTPIncomingAdd(loop,receiver,tsArg(args,0),tsArg(args,1),tsArg(args,2));return tsU}))
	p.set("_addHeaderLineDistinct",tsNamedNativeMethod("_addHeaderLineDistinct",3,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{tsNodeHTTPIncomingDistinct(loop,tsArg(args,0),tsArg(args,1),tsArg(args,2));return tsU}))
	p.set("_dumpAndCloseReadable", tsNamedNativeMethod("_dumpAndCloseReadable", 0, func(loop *tsLoop, receiver tsValue, _ ...tsValue) tsValue { return tsNodeHTTPIncomingDumpAndClose(loop, receiver) }))
	p.set("_dump",tsNamedNativeMethod("_dump",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{if !tsTruthy(tsNodeHTTPIncomingGet(loop,receiver,"_dumped")){tsNodeHTTPIncomingSet(loop,receiver,"_dumped",tsBooleanValue(true));tsNodeHTTPIncomingCall(loop,receiver,"removeAllListeners",tsNodeHTTPIncomingKey("data"));tsNodeHTTPIncomingCall(loop,receiver,"resume")};return tsU}))
	p.descriptors=map[string]*tsDescriptor{}
	p.descriptors["connection"]=&tsDescriptor{accessor:true,configurable:true,getter:tsNamedNativeMethod("get connection",0,func(loop *tsLoop,r tsValue,_ ...tsValue)tsValue{return tsNodeHTTPIncomingGet(loop,r,"socket")}),setter:tsNamedNativeMethod("set connection",1,func(loop *tsLoop,r tsValue,a ...tsValue)tsValue{tsNodeHTTPIncomingSet(loop,r,"socket",tsArg(a,0));return tsU})}
	p.descriptors["aborted"]=&tsDescriptor{accessor:true,configurable:true,getter:tsNamedNativeMethod("get aborted",0,func(loop *tsLoop,r tsValue,_ ...tsValue)tsValue{return tsGet(loop,r,keys["aborted"])}),setter:tsNamedNativeMethod("set aborted",1,func(loop *tsLoop,r tsValue,a ...tsValue)tsValue{return tsSet(loop,r,keys["aborted"],tsArg(a,0))})}
	p.descriptors["rawHeaders"]=&tsDescriptor{accessor:true,enumerable:true,configurable:true,getter:tsNamedNativeMethod("get rawHeaders",0,func(loop *tsLoop,r tsValue,_ ...tsValue)tsValue{return tsNodeHTTPIncomingRaw(loop,r)}),setter:tsNamedNativeMethod("set rawHeaders",1,func(loop *tsLoop,r tsValue,a ...tsValue)tsValue{v:=tsArg(a,0);tsSet(loop,r,keys["raw"],v);n:=0;if v.kind==tsArrayKind{n=len((*tsArray)(v.ref).values)};tsSet(loop,r,keys["headerCount"],tsNumberValue(float64(n)));return tsU})}
	for _,item:=range []struct{name string;distinct bool}{{"headers",false},{"headersDistinct",true},{"trailers",false},{"trailersDistinct",true}}{entry:=item;p.descriptors[entry.name]=&tsDescriptor{accessor:true,configurable:true,getter:tsNamedNativeMethod("get "+entry.name,0,func(loop *tsLoop,r tsValue,_ ...tsValue)tsValue{return tsNodeHTTPIncomingHeaders(loop,r,entry.name,entry.distinct)}),setter:tsNamedNativeMethod("set "+entry.name,1,func(loop *tsLoop,r tsValue,a ...tsValue)tsValue{cache:=keys["headers"];if strings.HasPrefix(entry.name,"trailers"){cache=keys["trailers"]};if strings.HasSuffix(entry.name,"Distinct"){if strings.HasPrefix(entry.name,"trailers"){cache=keys["trailersDistinct"]}else{cache=keys["headersDistinct"]}};return tsSet(loop,r,cache,tsArg(a,0))})}}
	p.descriptors["signal"]=&tsDescriptor{accessor:true,configurable:true,getter:tsNamedNativeMethod("get signal",0,func(loop *tsLoop,r tsValue,_ ...tsValue)tsValue{controller:=tsGet(loop,r,keys["abortController"]);if tsIsNull(controller){controller=tsAbortController(loop);tsSet(loop,r,keys["abortController"],controller);signal:=tsNodeHTTPIncomingGet(loop,controller,"signal");if tsTruthy(tsNodeHTTPIncomingGet(loop,r,"destroyed")){tsNodeHTTPIncomingCall(loop,controller,"abort")}else{cb:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{tsNodeHTTPIncomingCall(loop,controller,"abort");return tsU}));tsNodeHTTPIncomingCall(loop,r,"once",tsNodeHTTPIncomingKey("close"),cb)};return signal};return tsNodeHTTPIncomingGet(loop,controller,"signal")})}
	tsNodeHTTPIncomingSingleton.module=tsObjectValue(tsNewObject())
}
func tsNodeHTTPIncomingRawHeaders(loop *tsLoop, msg tsValue) tsValue {
	keys := tsNodeHTTPIncomingSingleton.keys
	raw := tsGet(loop, msg, keys["raw"])
	if !tsIsNull(raw) { return raw }
	if source := tsGet(loop, msg, keys["headerSource"]); !tsNullish(source) {
		if take := tsNodeHTTPIncomingGet(loop, source, "takeRawHeaders"); take.kind == tsFunctionKind { raw = tsCallReceiver(loop, take, source, tsGet(loop, msg, keys["headers"])) }
	}
	if tsIsUndefined(raw) { raw = tsArrayValue(&tsArray{}) }
	tsSet(loop, msg, keys["raw"], raw)
	count := 0
	if raw.kind == tsArrayKind { count = len((*tsArray)(raw.ref).values) }
	tsSet(loop, msg, keys["headerCount"], tsNumberValue(float64(count)))
	return raw
}
func tsNodeHTTPIncomingReadStart(loop *tsLoop,socket tsValue)tsValue{if !tsNullish(socket)&&!tsTruthy(tsNodeHTTPIncomingGet(loop,socket,"_paused"))&&tsTruthy(tsNodeHTTPIncomingGet(loop,socket,"readable")){tsNodeHTTPIncomingCall(loop,socket,"resume")};return tsU}
func tsNodeHTTPIncomingReadStop(loop *tsLoop,socket tsValue)tsValue{if !tsNullish(socket){tsNodeHTTPIncomingCall(loop,socket,"pause")};return tsU}
func tsNodeHTTPIncomingModule(loop *tsLoop)tsValue{tsNodeHTTPIncomingInit(loop);module:=(*tsObject)(tsNodeHTTPIncomingSingleton.module.ref);if !tsIsUndefined(module.get("IncomingMessage")){return tsObjectValue(module)};module.set("IncomingMessage",tsNodeHTTPIncomingSingleton.constructor);module.set("kReqShouldKeepAlive",tsNodeHTTPIncomingSingleton.keys["reqKeepAlive"]);module.set("readStart",tsNamedNativeMethod("readStart",1,func(loop *tsLoop,_ tsValue,a ...tsValue)tsValue{return tsNodeHTTPIncomingReadStart(loop,tsArg(a,0))}));module.set("readStop",tsNamedNativeMethod("readStop",1,func(loop *tsLoop,_ tsValue,a ...tsValue)tsValue{return tsNodeHTTPIncomingReadStop(loop,tsArg(a,0))}));module.set("default",tsObjectValue(module));return tsObjectValue(module)}
`
