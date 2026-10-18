package goemit

// nodeQUICRuntime preserves Bun's experimental QUIC module and internal helper
// surfaces. The transport adapter uses quic-go while mapping the JS objects and
// callbacks onto this runtime's event loop.
const nodeQUICRuntime = `
type tsNodeQUICEndpointState struct{object *tsObject;options tsValue;closed,closing,listening,destroyed,busy,refed,clientHTTPKnown,clientHTTP,suppressCloseChannels bool;maxPerHost,maxTotal,activeSessions,creationOrder int;idleTimeout time.Duration;idleTimer *time.Timer;address,closeError tsValue;onSession tsValue;stats tsValue;closedPromise *tsPromise;task *tsTask;sniMu sync.RWMutex;sniCertificates map[string]cryptotls.Certificate}
type tsNodeQUICSessionState struct{object *tsObject;endpoint tsValue;closed,closing,destroyed,opened,client,graceful,gracefulClose,closeApplication bool;closeCode uint64;closeReason string;nextBidi,nextUni int;nextDatagram uint64;maxPendingDatagrams int;streamIdleTimeout time.Duration;streams []tsValue;callbacks map[string]tsValue;stats tsValue;path tsValue;options tsValue;localCertificate *x509.Certificate;openedPromise,closedPromise *tsPromise}
type tsNodeQUICStreamState struct{object *tsObject;session tsValue;id int;direction string;closed,destroyed,ended,early,opened,outboundSet,writeEnded,writerErrored bool;highWaterMark int;headers,pendingTrailers,writer,writerError,fileHandle tsValue;blob *tsBlob;blobOffset int64;outboundIterator *tsAsyncIterator;callbacks map[string]tsValue;priority tsValue;stats tsValue;closedPromise *tsPromise;idleTimer *time.Timer;drains []*tsPromise}
var tsNodeQUICSingleton struct{sync.Once;module,binding,symbols,diagnostics,http2util,states,stats tsValue;endpointProto,sessionProto,streamProto,errorProto *tsObject;endpointClass,sessionClass,streamClass,errorClass tsValue;endpoints map[*tsObject]*tsNodeQUICEndpointState;sessions map[*tsObject]*tsNodeQUICSessionState;streams map[*tsObject]*tsNodeQUICStreamState;quicErrors map[*tsObject]bool;channels map[string]tsValue;private,createDisconnected tsValue;symbolValues map[string]tsValue;nextEndpointOrder int;callbacksInitialized bool}
func tsNodeQUICKey(s string)tsValue{return tsStringReference(tsStringUTF8(s))}
func tsNodeQUICGet(loop *tsLoop,v tsValue,n string)tsValue{return tsGet(loop,v,tsNodeQUICKey(n))}
func tsNodeQUICSet(loop *tsLoop,v tsValue,n string,x tsValue){tsSet(loop,v,tsNodeQUICKey(n),x)}
func tsNodeQUICError(code,message string)tsValue{return tsFSException("Error",code,message)}
func tsNodeQUICArray(names ...string)tsValue{a:=&tsArray{};for _,n:=range names{a.values=append(a.values,tsNodeQUICKey(n))};return tsArrayValue(a)}
func tsNodeQUICPublish(loop *tsLoop,name string,payload tsValue){diagnostics:=tsInternalQUICDiagnosticsModule(loop);channel:=tsNodeQUICGet(loop,diagnostics,name);publish:=tsNodeQUICGet(loop,channel,"publish");if publish.kind==tsFunctionKind{tsCallReceiver(loop,publish,channel,payload)}}
func tsNodeQUICData(fields map[string]tsValue)tsValue{object:=tsNewObject();for name,value:=range fields{object.set(name,value)};return tsObjectValue(object)}
func tsNodeQUICMakeConstructor(name string,proto *tsObject,length int,construct func(*tsLoop,[]tsValue)tsValue)tsValue{fn:=tsFunc(func(_ *tsLoop,_ ...tsValue)tsValue{tsPropertyFailure("Class constructor "+name+" cannot be invoked without 'new'");return tsU});fn.name=name;fn.length=length;fn.constructible=true;fn.construct=func(loop *tsLoop,args ...tsValue)tsValue{return construct(loop,args)};fn.properties=tsNewObject();fn.properties.set("prototype",tsObjectValue(proto));fn.properties.descriptors=map[string]*tsDescriptor{"prototype":{writable:false}};proto.set("constructor",tsFunctionValue(fn));return tsFunctionValue(fn)}
func tsNodeQUICDefineGetter(proto *tsObject,name string,get func(*tsLoop,tsValue)tsValue,set func(*tsLoop,tsValue,tsValue)tsValue){if proto.descriptors==nil{proto.descriptors=map[string]*tsDescriptor{}};descriptor:=&tsDescriptor{accessor:true,configurable:true,getter:tsNamedNativeMethod("get "+name,0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{return get(loop,receiver)})};if set!=nil{descriptor.setter=tsNamedNativeMethod("set "+name,1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return set(loop,receiver,tsArg(args,0))})};proto.descriptors[name]=descriptor}
func tsNodeQUICZeroStats(loop *tsLoop,names []string)tsValue{_ = loop;o:=tsNewObject();for _,n:=range names{o.set(n,tsBigIntValue(big.NewInt(0)))};return tsObjectValue(o)}
func tsNodeQUICStatsSet(stats tsValue,name string,value int64){if stats.kind!=tsObjectKind{return};(*tsObject)(stats.ref).set(name,tsBigIntSigned(value))}
func tsNodeQUICStatsSetUnsigned(stats tsValue,name string,value uint64){if stats.kind!=tsObjectKind{return};(*tsObject)(stats.ref).set(name,tsBigIntValue(new(big.Int).SetUint64(value)))}
func tsNodeQUICStatsAdd(stats tsValue,name string,value uint64){if stats.kind!=tsObjectKind||value==0{return};record:=(*tsObject)(stats.ref);current:=record.values[name];next:=new(big.Int);if current.kind==tsBigIntKind{next.Set((*big.Int)(current.ref))}else if tsIsNumeric(current){next.SetInt64(int64(tsNumber(current)))};next.Add(next,new(big.Int).SetUint64(value));record.values[name]=tsBigIntValue(next)}
func tsNodeQUICRecord(loop *tsLoop,v tsValue,n string)tsValue{if v.kind!=tsObjectKind{return tsU};return (*tsObject)(v.ref).get(n)}
func tsNodeQUICCloneProperties(loop *tsLoop,v tsValue)tsValue{o:=tsNewObject();if v.kind==tsObjectKind||v.kind==tsInstanceKind{for _,k:=range tsEnumerableOwnKeys(loop,v).values{o.set(tsPropertyKey(k),tsGet(loop,v,k))}};return tsObjectValue(o)}
func tsNodeQUICInstallEmitter(loop *tsLoop,obj *tsObject,proto *tsObject)tsValue{events:=tsNodeEventsModule(loop);ctor:=tsNodeQUICGet(loop,events,"EventEmitter");eventProto:=tsNodeQUICGet(loop,ctor,"prototype");if proto!=nil{proto.prototype=eventProto;obj.prototype=tsObjectValue(proto)}else{obj.prototype=eventProto};return tsNodeEventCreate(loop,tsObjectValue(obj),tsU)}
func tsNodeQUICProcessEndpointOptions(loop *tsLoop, options tsValue) tsValue {
	processed := tsNodeQUICCloneProperties(loop, options)
	address := tsNodeQUICGet(loop, options, "address")
	if tsIsUndefined(address) {
		return processed
	}
	tsNodeNetModule(loop)
	if address.kind == tsStringKind {
		constructor := tsNodeQUICGet(loop, tsNodeNetModule(loop), "SocketAddress")
		parse := tsNodeQUICGet(loop, constructor, "parse")
		address = tsCall(loop, parse, address)
	} else if !tsNodeNetIsSocketAddress(address) {
		address = tsNodeNetSocketAddressNew(loop, address, false, "", 0)
	}
	tsSet(loop,processed,tsNodeQUICKey("address"),address)
	return processed
}
func tsNodeQUICCreateEndpoint(loop *tsLoop, options tsValue) tsValue {
	tsNodeQUICInit(loop)
	if tsIsUndefined(options) {
		options = tsObjectValue(tsNewObject())
	}
	tsNodeQUICValidateEndpointOptions(loop, options)
	config := options
	options = tsNodeQUICProcessEndpointOptions(loop, options)
	obj := tsNewObject()
	tsNodeQUICInstallEmitter(loop, obj, tsNodeQUICSingleton.endpointProto)
	tsNodeQUICSingleton.nextEndpointOrder++
	state := &tsNodeQUICEndpointState{
		object: obj, options: options, maxPerHost: 100, maxTotal: 10000,
		refed: true, creationOrder: tsNodeQUICSingleton.nextEndpointOrder,
		address: tsU, onSession: tsU,
		stats: tsNodeQUICZeroStats(loop, []string{"createdAt", "destroyedAt", "bytesReceived", "bytesSent", "packetsReceived", "packetsSent", "serverSessions", "clientSessions", "serverBusyCount", "retryCount", "retryRateLimited", "versionNegotiationCount", "versionNegotiationRateLimited", "statelessResetCount", "statelessResetRateLimited", "immediateCloseCount", "immediateCloseRateLimited", "sessionCreationRateLimited", "packetsBlocked"}),
		closedPromise: loop.promise(),
	}
	// Endpoint idleTimeout is expressed in seconds; the wrapper's idle endpoint
	// timer multiplies it by 1000 when calling setTimeout.
	if timeout, ok := tsNodeQUICOptionNumber(loop, options, "idleTimeout"); ok && timeout > 0 && timeout < float64(math.MaxInt64)/float64(time.Second) {
		state.idleTimeout = time.Duration(timeout * float64(time.Second))
	}
	if value := tsNodeQUICGet(loop, options, "maxConnectionsPerHost"); !tsIsUndefined(value) {
		state.maxPerHost = int(tsNodeQUICValidatedEndpointLimit(value, "maxConnectionsPerHost"))
	}
	if value := tsNodeQUICGet(loop, options, "maxConnectionsTotal"); !tsIsUndefined(value) {
		state.maxTotal = int(tsNodeQUICValidatedEndpointLimit(value, "maxConnectionsTotal"))
	}
	tsNodeQUICSingleton.endpoints[obj] = state
	tsNodeQUICStatsSet(state.stats, "createdAt", time.Now().UnixNano())
	obj.set("options", options)
	tsNodeQUICPublish(loop, "onEndpointCreatedChannel", tsNodeQUICData(map[string]tsValue{"endpoint": tsObjectValue(obj), "config": config}))
	return tsObjectValue(obj)
}
func tsNodeQUICValidatedEndpointLimit(value tsValue,name string)uint64{if tsIsNumeric(value){number:=tsNumber(value);if number>=0&&number<=65535&&math.Trunc(number)==number{return uint64(number)}};tsFSInvalid(name,"an integer from 0 to 65535");return 0}
func tsNodeQUICValidateEndpointOptions(loop *tsLoop, options tsValue) {
	if !tsNodeQUICObject(options) { tsFSInvalid("options", "an object") }
	for _, name := range []string{"maxConnectionsPerHost", "maxConnectionsTotal"} {
		if value := tsNodeQUICGet(loop, options, name); !tsIsUndefined(value) { tsNodeQUICValidatedEndpointLimit(value, name) }
	}
	if policy := tsNodeQUICGet(loop, options, "blockListPolicy"); !tsIsUndefined(policy) && (policy.kind != tsStringKind || (tsText(policy) != "deny" && tsText(policy) != "allow")) { tsFSInvalid("options.blockListPolicy", "'deny' or 'allow'") }
	if list := tsNodeQUICGet(loop, options, "blockList"); !tsIsUndefined(list) {
		netModule := tsNodeNetModule(loop)
		constructor := tsNodeQUICGet(loop, netModule, "BlockList")
		check := tsNodeQUICGet(loop, constructor, "isBlockList")
		if check.kind != tsFunctionKind || !tsTruthy(tsCall(loop, check, list)) { tsFSInvalid("options.blockList", "a net.BlockList") }
	}
	for _, name := range []string{"retryTokenExpiration", "tokenExpiration", "addressLRUSize"} {
		value := tsNodeQUICGet(loop, options, name)
		if tsIsUndefined(value) { continue }
		if value.kind == tsBigIntKind { if (*big.Int)(value.ref).Sign() >= 0 { continue } } else if tsIsNumeric(value) && tsNumber(value) >= 0 && !math.IsInf(tsNumber(value), 0) && !math.IsNaN(tsNumber(value)) && math.Trunc(tsNumber(value)) == tsNumber(value) { continue }
		tsFSInvalid("options."+name, "a non-negative integer")
	}
	for _, name := range []string{"retryRate", "retryBurst", "statelessResetRate", "statelessResetBurst", "versionNegotiationRate", "versionNegotiationBurst", "immediateCloseRate", "immediateCloseBurst", "sessionCreationRate", "sessionCreationBurst", "rxDiagnosticLoss", "txDiagnosticLoss"} {
		value := tsNodeQUICGet(loop, options, name)
		if !tsIsUndefined(value) && (!tsIsNumeric(value) || math.IsNaN(tsNumber(value)) || tsNumber(value) < 0) { tsFSInvalid("options."+name, "a non-negative number") }
	}
	for _, name := range []string{"udpReceiveBufferSize", "udpSendBufferSize"} {
		value := tsNodeQUICGet(loop, options, name)
		if !tsIsUndefined(value) && (!tsIsNumeric(value) || tsNumber(value) < 0 || math.IsInf(tsNumber(value), 0) || math.IsNaN(tsNumber(value)) || math.Trunc(tsNumber(value)) != tsNumber(value)) { tsFSInvalid("options."+name, "a non-negative integer") }
	}
	if value := tsNodeQUICGet(loop, options, "udpTTL"); !tsIsUndefined(value) && (!tsIsNumeric(value) || tsNumber(value) < 0 || tsNumber(value) > 255 || math.Trunc(tsNumber(value)) != tsNumber(value)) { tsFSInvalid("options.udpTTL", "an integer from 0 to 255") }
	for _, name := range []string{"resetTokenSecret", "tokenSecret"} {
		value := tsNodeQUICGet(loop, options, name)
		if tsIsUndefined(value) { continue }
		if value.kind != tsTypedArrayKind && value.kind != tsDataViewKind { tsFSInvalid("options."+name, "an ArrayBufferView") }
		if len(tsFSBytes(value)) != 16 { tsFSInvalid("options."+name, "a 16-byte ArrayBufferView") }
	}
	if address := tsNodeQUICGet(loop, options, "address"); !tsIsUndefined(address) && address.kind != tsStringKind && !tsNodeQUICObject(address) { tsFSInvalid("options.address", "a SocketAddress, string, or object") }
}
func tsNodeQUICCreateSession(loop *tsLoop, endpoint, options tsValue) tsValue {
	tsNodeQUICInit(loop)
	if tsNullish(options) { options = tsObjectValue(tsNewObject()) }
	obj := tsNewObject()
	tsNodeQUICInstallEmitter(loop, obj, tsNodeQUICSingleton.sessionProto)
	state := &tsNodeQUICSessionState{object: obj, endpoint: endpoint, nextBidi: 0, nextUni: 2, nextDatagram: 1, maxPendingDatagrams: 128, callbacks: map[string]tsValue{}, stats: tsNodeQUICZeroStats(loop, []string{"createdAt", "destroyedAt", "closingAt", "handshakeCompletedAt", "handshakeConfirmedAt", "bytesReceived", "bytesSent", "bidiInStreamCount", "bidiOutStreamCount", "uniInStreamCount", "uniOutStreamCount", "maxBytesInFlight", "bytesInFlight", "blockCount", "cwnd", "latestRtt", "minRtt", "rttVar", "smoothedRtt", "ssthresh", "pktSent", "bytesRecv", "pktRecv", "pktLost", "bytesLost", "pingRecv", "pktDiscarded", "datagramsReceived", "datagramsSent", "datagramsAcknowledged", "datagramsLost", "streamsIdleTimedOut"}), path: tsNull, options: options}
	state.openedPromise, state.closedPromise = loop.promise(), loop.promise()
	if timeout, ok := tsNodeQUICOptionNumber(loop, options, "streamIdleTimeout"); ok && timeout > 0 && timeout < float64(math.MaxInt64)/float64(time.Millisecond) { state.streamIdleTimeout = time.Duration(timeout * float64(time.Millisecond)) }
	state.localCertificate = tsNodeQUICLocalCertificate(loop, options)
	tsNodeQUICSingleton.sessions[obj] = state
	tsNodeQUICStatsSet(state.stats, "createdAt", time.Now().UnixNano())
	obj.set("endpoint", endpoint)
	obj.set("applicationOptions", tsNodeQUICApplicationOptions(loop, options))
	obj.set("localTransportParams", tsNodeQUICLocalTransportParams(loop, options))
	obj.set("remoteTransportParams", tsU)
	tsNodeQUICApplySessionCallbacks(loop, state, options)
	return tsObjectValue(obj)
}

func tsNodeQUICApplicationOptions(loop *tsLoop, options tsValue) tsValue {
	input := tsNodeQUICGet(loop, options, "application")
	result := tsNewObject()
	for _, field := range []struct{name string; source string; fallback uint64}{
		{"maxHeaderPairs", "maxHeaderPairs", 128},
		{"maxHeaderLength", "maxHeaderLength", 16384},
		{"maxFieldSectionSize", "maxFieldSectionSize", 0},
		{"qpackMaxDtableCapacity", "qpackMaxDTableCapacity", 0},
		{"qpackEncoderMaxDtableCapacity", "qpackEncoderMaxDTableCapacity", 0},
		{"qpackBlockedStreams", "qpackBlockedStreams", 0},
	} {
		value := tsNodeQUICGet(loop, input, field.source)
		if value.kind == tsBigIntKind && (*big.Int)(value.ref).Sign() >= 0 {
			result.set(field.name, tsBigIntValue(new(big.Int).Set((*big.Int)(value.ref))))
		} else if tsIsNumeric(value) && tsNumber(value) >= 0 && math.Trunc(tsNumber(value)) == tsNumber(value) {
			result.set(field.name, tsBigIntUnsigned(uint64(tsNumber(value))))
		} else if tsIsUndefined(value) {
			result.set(field.name, tsBigIntUnsigned(field.fallback))
		} else {
			tsFSInvalid("application."+field.source, "a non-negative integer")
		}
	}
	for _, field := range []struct{name string; fallback bool}{{"enableConnectProtocol", false}, {"enableDatagrams", true}} {
		value := tsNodeQUICGet(loop, input, field.name)
		if tsIsUndefined(value) { result.set(field.name, tsBooleanValue(field.fallback)) } else { result.set(field.name, tsBooleanValue(tsTruthy(value))) }
	}
	return tsObjectValue(result)
}

func tsNodeQUICLocalTransportParams(loop *tsLoop, options tsValue) tsValue {
	input := tsNodeQUICGet(loop, options, "transportParams")
	result := tsNewObject()
	for _, field := range []struct {
		name     string
		fallback uint64
	}{
		{"initialMaxStreamDataBidiLocal", 256 * 1024},
		{"initialMaxStreamDataBidiRemote", 256 * 1024},
		{"initialMaxStreamDataUni", 256 * 1024},
		{"initialMaxData", 1024 * 1024},
		{"initialMaxStreamsBidi", 100},
		{"initialMaxStreamsUni", 3},
		{"maxIdleTimeout", 10},
		{"maxUdpPayloadSize", 1200},
		{"ackDelayExponent", 3},
		{"maxAckDelay", 25},
		{"activeConnectionIDLimit", 8},
		{"maxDatagramFrameSize", 1200},
	} {
		value := tsNodeQUICGet(loop, input, field.name)
		if value.kind == tsBigIntKind {
			result.set(field.name, tsBigIntValue(new(big.Int).Set((*big.Int)(value.ref))))
		} else if tsIsNumeric(value) {
			result.set(field.name, tsBigIntUnsigned(uint64(tsNumber(value))))
		} else {
			result.set(field.name, tsBigIntUnsigned(field.fallback))
		}
	}
	result.set("disableActiveMigration", tsBooleanValue(tsTruthy(tsNodeQUICGet(loop, input, "disableActiveMigration"))))
	result.set("initialSCID", tsU)
	result.set("retrySCID", tsU)
	result.set("originalDCID", tsU)
	return tsObjectValue(result)
}

func tsNodeQUICCreateStream(loop *tsLoop, session tsValue, id int, direction string) tsValue {
	tsNodeDuplexInit(loop)
	var stream tsValue
	options := tsNewObject()
	options.set("write", tsFunctionValue(tsFunc(func(loop *tsLoop, args ...tsValue) tsValue {
		tsNodeQUICWriteTransport(loop, stream, tsArg(args, 0), tsArg(args, 1), tsArg(args, 2))
		return tsU
	})))
	options.set("final", tsFunctionValue(tsFunc(func(loop *tsLoop, args ...tsValue) tsValue {
		tsNodeQUICFinalizeTransportWrite(loop, stream, tsArg(args, 0))
		return tsU
	})))
	value := tsNodeDuplexCreate(loop, tsU, tsObjectValue(options))
	stream = value
	obj := (*tsObject)(value.ref)
	parent := obj.prototype
	proto := tsNodeQUICSingleton.streamProto
	if tsIsUndefined(proto.prototype) {
		proto.prototype = parent
	}
	obj.prototype = tsObjectValue(proto)
	publicDirection := "bidi"
	if direction == "unidirectional" { publicDirection = "uni" }
	defaultPriority := tsNewObject()
	defaultPriority.set("level", tsNodeQUICKey("default"))
	defaultPriority.set("incremental", tsBooleanValue(false))
	state := &tsNodeQUICStreamState{object: obj, session: session, id: id, direction: publicDirection, opened: true, highWaterMark: 16 * 1024, headers: tsU, pendingTrailers: tsU, callbacks: map[string]tsValue{}, priority: tsObjectValue(defaultPriority), stats: tsNodeQUICZeroStats(loop, []string{"createdAt", "openedAt", "receivedAt", "ackedAt", "destroyedAt", "bytesReceived", "bytesSent", "maxOffset", "maxOffsetAck", "maxOffsetRecv", "finalSize", "bytesAccumulated", "maxBytesAccumulated"})}
	state.closedPromise = loop.promise()
	tsNodeQUICSingleton.streams[obj] = state
	tsNodeQUICStatsSet(state.stats, "createdAt", time.Now().UnixNano())
	tsNodeQUICStatsSet(state.stats, "openedAt", time.Now().UnixNano())
	tsNodeQUICStreamTransportStateForInit()
	tsNodeQUICStreamTransportStates[uintptr(unsafe.Pointer(obj))] = &tsNodeQUICStreamWireState{}
	if sessionState := tsNodeQUICSessionFor(session); sessionState != nil {
		sessionState.streams = append(sessionState.streams, value)
		for _, name := range []string{"onheaders", "ontrailers", "oninfo", "onwanttrailers"} {
			if callback := sessionState.callbacks[name]; callback.kind == tsFunctionKind {
				state.callbacks[name] = callback
			}
		}
	}
	obj.set("id", tsNumberValue(float64(id)))
	obj.set("direction", tsNodeQUICKey(publicDirection))
	obj.set("session", session)
	return value
}
func tsNodeQUICEndpointFor(v tsValue)*tsNodeQUICEndpointState{if v.kind!=tsObjectKind{return nil};return tsNodeQUICSingleton.endpoints[(*tsObject)(v.ref)]}
func tsNodeQUICSessionFor(v tsValue)*tsNodeQUICSessionState{if v.kind!=tsObjectKind{return nil};return tsNodeQUICSingleton.sessions[(*tsObject)(v.ref)]}
func tsNodeQUICStreamFor(v tsValue)*tsNodeQUICStreamState{if v.kind!=tsObjectKind{return nil};return tsNodeQUICSingleton.streams[(*tsObject)(v.ref)]}
func tsNodeQUICEndStream(loop *tsLoop, receiver, errorValue tsValue, options ...tsValue) tsValue {
	state := tsNodeQUICStreamFor(receiver)
	if state == nil || state.closed { return receiver }
	session := state.session
	defaultCode := uint64(0)
	if !tsIsUndefined(errorValue) { defaultCode = 1 }
	code := tsNodeQUICStreamErrorCode(loop, errorValue, defaultCode)
	if len(options) != 0 {
		if optionCode := tsNodeQUICGet(loop, options[0], "code"); !tsIsUndefined(optionCode) { code = tsNodeQUICValidateWireCode(optionCode, "options.code") }
	}
	tsNodeQUICResetTransportStream(loop, receiver, errorValue, code)
	if state.idleTimer != nil { state.idleTimer.Stop(); state.idleTimer = nil }
	tsNodeQUICCloseFileHandle(loop, state)
	state.closed, state.destroyed = true, true
	tsNodeQUICStatsSet(state.stats, "destroyedAt", time.Now().UnixNano())
	rejected := !tsIsUndefined(errorValue)
	if rejected { if callback:=state.callbacks["onerror"];callback.kind==tsFunctionKind{tsCallReceiver(loop,callback,receiver,errorValue)};tsCallReceiver(loop, tsNodeQUICGet(loop, receiver, "emit"), receiver, tsNodeQUICKey("error"), errorValue) }
	if state.closedPromise != nil { state.closedPromise.settle(tsResult{errorValue, rejected}) }
	tsCallReceiver(loop, tsNodeQUICGet(loop, receiver, "emit"), receiver, tsNodeQUICKey("close"))
	tsNodeQUICPublish(loop, "onStreamClosedChannel", tsNodeQUICData(map[string]tsValue{"stream":receiver,"session":session,"error":errorValue,"stats":state.stats}))
	state.session = tsNull
	if sessionState := tsNodeQUICSessionFor(session); sessionState != nil && sessionState.graceful { tsNodeQUICMaybeFinishSession(loop, session) }
	return tsPromiseValue(state.closedPromise)
}

func tsNodeQUICGracefulCloseSession(loop *tsLoop, receiver, options tsValue) tsValue {
	state := tsNodeQUICSessionFor(receiver)
	if state == nil { return tsU }
	if state.closed || state.closing { return tsPromiseValue(state.closedPromise) }
	if tsIsUndefined(options) { options = tsObjectValue(tsNewObject()) }
	if !tsNodeQUICObject(options) { tsFSInvalid("options", "an object") }
	codeValue := tsNodeQUICGet(loop, options, "code")
	code := uint64(0)
	if !tsIsUndefined(codeValue) {
		if codeValue.kind != tsBigIntKind && !tsIsNumeric(codeValue) { tsFSInvalid("options.code", "a number or bigint") }
		code = tsNodeQUICValidateWireCode(codeValue, "options.code")
	}
	typeValue := tsNodeQUICGet(loop, options, "type")
	if tsIsUndefined(typeValue) { typeValue = tsNodeQUICKey("transport") }
	if typeValue.kind != tsStringKind || (tsText(typeValue) != "transport" && tsText(typeValue) != "application") { tsFSInvalid("options.type", "'transport' or 'application'") }
	reasonValue := tsNodeQUICGet(loop, options, "reason")
	reason := ""
	if !tsIsUndefined(reasonValue) {
		if reasonValue.kind != tsStringKind { tsFSInvalid("options.reason", "a string") }
		reason = tsText(reasonValue)
	}
	state.closing, state.graceful = true, true
	state.gracefulClose = true
	state.closeCode, state.closeReason = code, reason
	state.closeApplication = tsStrictEqual(typeValue, tsNodeQUICKey("application"))
	tsNodeQUICStatsSet(state.stats, "closingAt", time.Now().UnixNano())
	tsNodeQUICPublish(loop, "onSessionClosingChannel", tsNodeQUICData(map[string]tsValue{"session":receiver}))
	tsNodeQUICBeginGracefulTransportSession(receiver)
	tsNodeQUICMaybeFinishSession(loop, receiver)
	return tsPromiseValue(state.closedPromise)
}

func tsNodeQUICMaybeFinishSession(loop *tsLoop, receiver tsValue) {
	state := tsNodeQUICSessionFor(receiver)
	if state == nil || state.closed || !state.closing || !state.graceful { return }
	for _, stream := range state.streams {
		if streamState := tsNodeQUICStreamFor(stream); streamState != nil && !streamState.closed { return }
	}
	options := tsNewObject()
	options.set("code", tsBigIntUnsigned(state.closeCode))
	closeType := "transport"
	if state.closeApplication { closeType = "application" }
	options.set("type", tsNodeQUICKey(closeType))
	options.set("reason", tsNodeQUICKey(state.closeReason))
	tsNodeQUICCloseSession(loop, receiver, tsU, tsObjectValue(options))
}

func tsNodeQUICCloseSession(loop *tsLoop, receiver, errorValue tsValue, options ...tsValue) tsValue {
	state := tsNodeQUICSessionFor(receiver)
	if state == nil { return tsU }
	if state.closed { return tsPromiseValue(state.closedPromise) }
	endpoint := state.endpoint
	tsNodeQUICStatsSet(state.stats, "closingAt", time.Now().UnixNano())
	code := uint64(0)
	reason := ""
	applicationClose := false
	wireClose := false
	if len(options) != 0 {
		if !tsIsUndefined(options[0]) && !tsNodeQUICObject(options[0]) { tsFSInvalid("options", "an object") }
		if tsNodeQUICObject(options[0]) {
			wireClose = true
			if optionCode := tsNodeQUICGet(loop, options[0], "code"); !tsIsUndefined(optionCode) { code = tsNodeQUICValidateWireCode(optionCode, "options.code") }
			if optionType := tsNodeQUICGet(loop, options[0], "type"); !tsIsUndefined(optionType) {
				if optionType.kind != tsStringKind || (tsText(optionType) != "transport" && tsText(optionType) != "application") { tsFSInvalid("options.type", "'transport' or 'application'") }
				applicationClose = tsText(optionType) == "application"
			}
			if optionReason := tsNodeQUICGet(loop, options[0], "reason"); !tsIsUndefined(optionReason) { if optionReason.kind != tsStringKind { tsFSInvalid("options.reason", "a string") }; reason = tsText(optionReason) }
		}
	}
	if wireClose { tsNodeQUICCloseTransportSession(receiver, applicationClose, code, reason) }
	state.closing, state.graceful = true, false
	for _, stream := range state.streams { tsNodeQUICEndStream(loop, stream, errorValue) }
	state.closed, state.destroyed = true, true
	tsNodeQUICStatsSet(state.stats, "destroyedAt", time.Now().UnixNano())
	errorProvided := !tsIsUndefined(errorValue)
	rejected := tsTruthy(errorValue)
	if !state.opened && state.openedPromise != nil {
		openError := errorValue
		if tsNullish(openError) { openError = tsNodeQUICError("ERR_QUIC_CONNECTION_FAILED", "Session closed before the handshake completed") }
		state.openedPromise.settle(tsResult{openError, true})
	}
	state.closedPromise.settle(tsResult{errorValue, rejected})
	if errorProvided {
		tsNodeQUICPublish(loop, "onSessionErrorChannel", tsNodeQUICData(map[string]tsValue{"session":receiver,"error":errorValue}))
		if callback := state.callbacks["onerror"]; callback.kind == tsFunctionKind { tsCallReceiver(loop, callback, receiver, errorValue) }
	}
	tsCallReceiver(loop, tsNodeQUICGet(loop, receiver, "emit"), receiver, tsNodeQUICKey("close"), errorValue)
	tsNodeQUICPublish(loop, "onSessionClosedChannel", tsNodeQUICData(map[string]tsValue{"session":receiver,"error":errorValue,"stats":state.stats}))
	state.endpoint, state.path = tsNull, tsU
	if endpointState := tsNodeQUICEndpointFor(endpoint); endpointState != nil {
		if endpointState.activeSessions > 0 { endpointState.activeSessions-- }
		tsNodeQUICFinishEndpoint(loop, endpoint)
		tsNodeQUICScheduleEndpointIdleClose(loop, endpoint)
	}
	return tsPromiseValue(state.closedPromise)
}
func tsNodeQUICFinishEndpoint(loop *tsLoop, receiver tsValue) {
	state := tsNodeQUICEndpointFor(receiver)
	if state == nil || state.closed || !state.closing || state.activeSessions != 0 { return }
	// Established QUIC sessions still need the UDP socket during graceful
	// shutdown. Stop the listener only once the final session has completed.
	tsNodeQUICCloseTransportEndpoint(receiver)
	state.closed, state.destroyed, state.listening = true, true, false
	state.address = tsU
	tsNodeQUICStatsSet(state.stats, "destroyedAt", time.Now().UnixNano())
	rejected := !tsIsUndefined(state.closeError)
	state.closedPromise.settle(tsResult{state.closeError, rejected})
	tsCallReceiver(loop, tsNodeQUICGet(loop, receiver, "emit"), receiver, tsNodeQUICKey("close"))
	if !state.suppressCloseChannels {
		if rejected { tsNodeQUICPublish(loop, "onEndpointErrorChannel", tsNodeQUICData(map[string]tsValue{"endpoint":receiver,"error":state.closeError})) }
		tsNodeQUICPublish(loop, "onEndpointClosedChannel", tsNodeQUICData(map[string]tsValue{"endpoint":receiver,"stats":state.stats}))
	}
}

func tsNodeQUICCancelEndpointIdle(state *tsNodeQUICEndpointState) {
	if state != nil && state.idleTimer != nil { state.idleTimer.Stop(); state.idleTimer = nil }
}

func tsNodeQUICScheduleEndpointIdleClose(loop *tsLoop, receiver tsValue) {
	state := tsNodeQUICEndpointFor(receiver)
	if state == nil || state.closed || state.closing || state.listening || state.activeSessions != 0 || state.idleTimeout <= 0 { return }
	tsNodeQUICCancelEndpointIdle(state)
	var timer *time.Timer
	timer = time.AfterFunc(state.idleTimeout, func() {
		loop.post(func() {
			current := tsNodeQUICEndpointFor(receiver)
			if current == nil || current.idleTimer != timer { return }
			current.idleTimer = nil
			current.suppressCloseChannels = true
			tsNodeQUICCloseEndpoint(loop, receiver, tsU)
		})
	})
	state.idleTimer = timer
}

func tsNodeQUICGracefulCloseEndpoint(loop *tsLoop, receiver tsValue) tsValue {
	state := tsNodeQUICEndpointFor(receiver)
	if state == nil { return tsU }
	if !state.closed && !state.closing {
		tsNodeQUICCancelEndpointIdle(state)
		state.closing = true
		if !state.suppressCloseChannels { tsNodeQUICPublish(loop, "onEndpointClosingChannel", tsNodeQUICData(map[string]tsValue{"endpoint":receiver,"hasPendingError":tsBooleanValue(!tsIsUndefined(state.closeError))})) }
	}
	tsNodeQUICFinishEndpoint(loop, receiver)
	return tsPromiseValue(state.closedPromise)
}

func tsNodeQUICCloseEndpoint(loop *tsLoop, receiver, errorValue tsValue) tsValue {
	state := tsNodeQUICEndpointFor(receiver)
	if state == nil { return tsU }
	if state.closed { return tsPromiseValue(state.closedPromise) }
	alreadyClosing := state.closing
	if !state.closing {
		tsNodeQUICCancelEndpointIdle(state)
		state.closing = true
		if !state.suppressCloseChannels { tsNodeQUICPublish(loop, "onEndpointClosingChannel", tsNodeQUICData(map[string]tsValue{"endpoint":receiver,"hasPendingError":tsBooleanValue(!tsIsUndefined(errorValue))})) }
	}
	state.closeError = errorValue
	tsNodeQUICCloseTransportEndpoint(receiver)
	for object, sessionState := range tsNodeQUICSingleton.sessions {
		if sessionState != nil && !sessionState.closed && tsStrictEqual(sessionState.endpoint, receiver) {
			closeOptions := tsU
			if !alreadyClosing && sessionState.opened && !tsNullish(errorValue) {
				closeOptions = tsNodeQUICErrorCloseOptions(loop, errorValue)
			}
			tsNodeQUICCloseSession(loop, tsObjectValue(object), errorValue, closeOptions)
		}
	}
	tsNodeQUICFinishEndpoint(loop, receiver)
	return tsPromiseValue(state.closedPromise)
}
func tsNodeQUICSetCallback(proto *tsObject, name string) {
	tsNodeQUICDefineGetter(proto, name, func(_ *tsLoop, receiver tsValue) tsValue {
		if stream := tsNodeQUICStreamFor(receiver); stream != nil {
			if value, ok := stream.callbacks[name]; ok { return value }
		}
		if session := tsNodeQUICSessionFor(receiver); session != nil {
			if value, ok := session.callbacks[name]; ok { return value }
		}
		if endpoint := tsNodeQUICEndpointFor(receiver); endpoint != nil && name == "onstream" { return endpoint.onSession }
		return tsU
	}, func(loop *tsLoop, receiver, value tsValue) tsValue {
		if value.kind != tsFunctionKind && !tsIsUndefined(value) { tsFSInvalid(name, "a function") }
		if stream := tsNodeQUICStreamFor(receiver); stream != nil {
			if value.kind == tsFunctionKind && (name == "onheaders" || name == "ontrailers" || name == "oninfo" || name == "onwanttrailers") {
				if !tsNodeQUICHeadersSupported(stream.session) {
					panic(tsThrown{tsNodeQUICError("ERR_INVALID_STATE", "The negotiated QUIC application protocol does not support headers")})
				}
			}
			stream.callbacks[name] = value
		}
		if session := tsNodeQUICSessionFor(receiver); session != nil { session.callbacks[name] = value }
		if endpoint := tsNodeQUICEndpointFor(receiver); endpoint != nil && name == "onstream" { endpoint.onSession = value }
		return value
	})
}
func tsNodeQUICInstall(loop *tsLoop){ep:=tsNodeQUICSingleton.endpointProto;sp:=tsNodeQUICSingleton.sessionProto;qp:=tsNodeQUICSingleton.streamProto
 for _,name:=range []string{"onerror","onstream","ondatagram","ondatagramstatus","onpathvalidation","onkeylog","onqlog","onsessionticket","onversionnegotiation","onhandshake","onnewtoken","onearlyrejected","onorigin","ongoaway"}{tsNodeQUICSetCallback(sp,name)}
 for _,name:=range []string{"onerror","onblocked","onreset","onheaders","ontrailers","oninfo","onwanttrailers"}{tsNodeQUICSetCallback(qp,name)}
 for _,name:=range []string{"closed","closing","listening","destroyed","busy","stats"}{n:=name;tsNodeQUICDefineGetter(ep,n,func(_ *tsLoop,r tsValue)tsValue{s:=tsNodeQUICEndpointFor(r);if s==nil{return tsU};if n=="closed"{return tsPromiseValue(s.closedPromise)};if n=="stats"{return s.stats};var b bool;switch n{case "closing":b=s.closing;case "listening":b=s.listening;case "destroyed":b=s.destroyed;case "busy":b=s.busy};return tsBooleanValue(b)},func(loop *tsLoop,r,v tsValue)tsValue{s:=tsNodeQUICEndpointFor(r);if s!=nil&&n=="busy"{if s.closed||s.closing{panic(tsThrown{tsNodeQUICError("ERR_INVALID_STATE","Endpoint is closed")})};busy:=tsTruthy(v);if busy!=s.busy{s.busy=busy;tsNodeQUICPublish(loop,"onEndpointBusyChangeChannel",tsNodeQUICData(map[string]tsValue{"endpoint":r,"busy":tsBooleanValue(busy)}))}};return v})}
 for _,name:=range []string{"maxConnectionsPerHost","maxConnectionsTotal"}{n:=name;tsNodeQUICDefineGetter(ep,n,func(_ *tsLoop,r tsValue)tsValue{s:=tsNodeQUICEndpointFor(r);if s==nil{return tsU};if n=="maxConnectionsPerHost"{return tsNumberValue(float64(s.maxPerHost))};return tsNumberValue(float64(s.maxTotal))},func(_ *tsLoop,r,v tsValue)tsValue{s:=tsNodeQUICEndpointFor(r);if s==nil{return v};number:=float64(-1);if tsIsNumeric(v){number=tsNumber(v)};if number<0||number>65535||math.Trunc(number)!=number{tsFSInvalid(n,"an integer from 0 to 65535")};if n=="maxConnectionsPerHost"{s.maxPerHost=int(number)}else{s.maxTotal=int(number)};return v})}
 tsNodeQUICDefineGetter(ep,"address",func(_ *tsLoop,r tsValue)tsValue{s:=tsNodeQUICEndpointFor(r);if s==nil||s.closed||s.closing||tsIsUndefined(s.address){return tsU};return s.address},nil)
 for _,name:=range []string{"closing","opened","closed","destroyed","maxDatagramSize","maxPendingDatagrams","stats","endpoint","path","certificate","peerCertificate","ephemeralKeyInfo","applicationOptions","localTransportParams","remoteTransportParams"}{n:=name;tsNodeQUICDefineGetter(sp,n,func(loop *tsLoop,r tsValue)tsValue{s:=tsNodeQUICSessionFor(r);if s==nil{return tsU};switch n{case "closing":return tsBooleanValue(s.closing);case "opened":return tsPromiseValue(s.openedPromise);case "closed":return tsPromiseValue(s.closedPromise);case "destroyed":return tsBooleanValue(s.destroyed);case "maxDatagramSize":return tsNodeQUICSessionMaxDatagramSize(r);case "maxPendingDatagrams":return tsNumberValue(float64(s.maxPendingDatagrams));case "stats":return s.stats;case "endpoint":return s.endpoint;case "path":return s.path;case "certificate":return tsNodeQUICSessionCertificate(loop,r);case "peerCertificate":return tsNodeQUICSessionPeerCertificate(loop,r);case "ephemeralKeyInfo":return tsNodeQUICEphemeralKeyInfo(loop,r);case "applicationOptions","localTransportParams","remoteTransportParams":return tsNodeQUICRecord(nil,r,n)};return tsU},func(_ *tsLoop,r,v tsValue)tsValue{if n=="maxPendingDatagrams"{if !tsIsNumeric(v)||tsNumber(v)<0||tsNumber(v)>65535||math.Trunc(tsNumber(v))!=tsNumber(v){tsFSInvalid("maxPendingDatagrams","an integer from 0 to 65535")};if s:=tsNodeQUICSessionFor(r);s!=nil{s.maxPendingDatagrams=int(tsNumber(v))};return v};return v})}
 for _,name:=range []string{"applicationOptions","localTransportParams","remoteTransportParams"}{n:=name;tsNodeQUICDefineGetter(sp,n,func(_ *tsLoop,r tsValue)tsValue{s:=tsNodeQUICSessionFor(r);if s==nil||s.closed{return tsNull};return tsNodeQUICRecord(nil,r,n)},func(_ *tsLoop,_ tsValue,v tsValue)tsValue{return v})}
 tsNodeQUICDefineGetter(sp,"stats",func(_ *tsLoop,r tsValue)tsValue{s:=tsNodeQUICSessionFor(r);if s==nil{return tsU};tsNodeQUICRefreshSessionStats(r);return s.stats},nil)
 tsNodeQUICDefineGetter(sp,"endpoint",func(_ *tsLoop,r tsValue)tsValue{s:=tsNodeQUICSessionFor(r);if s==nil||s.closed{return tsNull};return s.endpoint},nil)
 tsNodeQUICDefineGetter(sp,"path",func(_ *tsLoop,r tsValue)tsValue{s:=tsNodeQUICSessionFor(r);if s==nil||s.closed{return tsU};return s.path},nil)
 tsNodeQUICDefineGetter(qp,"session",func(_ *tsLoop,r tsValue)tsValue{s:=tsNodeQUICStreamFor(r);if s==nil||s.closed{return tsNull};return s.session},nil)
 for _,name:=range []string{"pending","early","highWaterMark","headers","pendingTrailers","stats","session","id","direction","destroyed","closed"}{n:=name;tsNodeQUICDefineGetter(qp,n,func(_ *tsLoop,r tsValue)tsValue{s:=tsNodeQUICStreamFor(r);if s==nil{return tsU};switch n{case "pending":return tsBooleanValue(!s.opened && !s.closed);case "early":return tsBooleanValue(s.early);case "highWaterMark":return tsNumberValue(float64(s.highWaterMark));case "headers":if tsIsUndefined(s.headers){return tsU};return s.headers;case "pendingTrailers":if tsIsUndefined(s.pendingTrailers){return tsU};return s.pendingTrailers;case "stats":return s.stats;case "session":return s.session;case "id":return tsBigIntUnsigned(uint64(s.id));case "direction":return tsNodeQUICKey(s.direction);case "destroyed":return tsBooleanValue(s.destroyed);case "closed":return tsPromiseValue(s.closedPromise)};return tsU},func(_ *tsLoop,r,v tsValue)tsValue{s:=tsNodeQUICStreamFor(r);if s==nil{return v};if n=="highWaterMark"{if !tsIsNumeric(v)||tsNumber(v)<0||tsNumber(v)>4294967295||math.Trunc(tsNumber(v))!=tsNumber(v){tsFSInvalid("highWaterMark","an integer from 0 to 4294967295")};s.highWaterMark=int(tsNumber(v))};return v})}
 tsNodeQUICDefineGetter(qp,"headers",func(_ *tsLoop,r tsValue)tsValue{s:=tsNodeQUICStreamFor(r);if s==nil||tsIsUndefined(s.headers){return tsU};return s.headers},nil)
 tsNodeQUICDefineGetter(qp,"pendingTrailers",func(_ *tsLoop,r tsValue)tsValue{s:=tsNodeQUICStreamFor(r);if s==nil||tsIsUndefined(s.pendingTrailers){return tsU};return s.pendingTrailers},func(loop *tsLoop,r,v tsValue)tsValue{s:=tsNodeQUICStreamFor(r);if s==nil{return v};if !tsNodeQUICHeadersSupported(s.session){panic(tsThrown{tsNodeQUICError("ERR_INVALID_STATE","The negotiated QUIC application protocol does not support headers")})};if tsIsUndefined(v){s.pendingTrailers=tsU;return v};tsNodeQUICValidateHeaders(v);s.pendingTrailers=v;return v})
 tsNodeQUICDefineGetter(qp,"priority",func(loop *tsLoop,r tsValue)tsValue{s:=tsNodeQUICStreamFor(r);if s==nil||s.closed{return tsNull};if !tsTruthy(tsNodeQUICGet(loop,tsNodeQUICSessionStateObject(s.session),"isPrioritySupported")){return tsNull};return s.priority},nil)
 tsNodeQUICDefineGetter(qp,"writer",func(loop *tsLoop,r tsValue)tsValue{return tsNodeQUICWriter(loop,r)},nil)
 ep.set("close",tsNamedNativeMethod("close",0,func(loop *tsLoop,r tsValue,_ ...tsValue)tsValue{return tsNodeQUICGracefulCloseEndpoint(loop,r)}));ep.set("destroy",tsNamedNativeMethod("destroy",1,func(loop *tsLoop,r tsValue,args ...tsValue)tsValue{return tsNodeQUICCloseEndpoint(loop,r,tsArg(args,0))}));ep.set("ref",tsNamedNativeMethod("ref",0,func(_ *tsLoop,r tsValue,_ ...tsValue)tsValue{if s:=tsNodeQUICEndpointFor(r);s!=nil{s.refed=true;if s.task!=nil{s.task.unref=false}};return r}));ep.set("unref",tsNamedNativeMethod("unref",0,func(_ *tsLoop,r tsValue,_ ...tsValue)tsValue{if s:=tsNodeQUICEndpointFor(r);s!=nil{s.refed=false;if s.task!=nil{s.task.unref=true}};return r}));ep.set("hasRef",tsNamedNativeMethod("hasRef",0,func(_ *tsLoop,r tsValue,_ ...tsValue)tsValue{s:=tsNodeQUICEndpointFor(r);return tsBooleanValue(s!=nil&&s.refed)}));ep.set("setSNIContexts",tsNamedNativeMethod("setSNIContexts",2,func(loop *tsLoop,r tsValue,args ...tsValue)tsValue{return tsNodeQUICSetSNIContexts(loop,r,tsArg(args,0),tsArg(args,1))}))
 sp.set("createBidirectionalStream", tsNamedNativeMethod("createBidirectionalStream", 1, func(loop *tsLoop, receiver tsValue, args ...tsValue) tsValue {
		return tsNodeQUICOpenStream(loop, receiver, "bidirectional", tsArg(args, 0))
	}))
 sp.set("createUnidirectionalStream", tsNamedNativeMethod("createUnidirectionalStream", 1, func(loop *tsLoop, receiver tsValue, args ...tsValue) tsValue {
		return tsNodeQUICOpenStream(loop, receiver, "unidirectional", tsArg(args, 0))
	}))
 sp.set("sendDatagram", tsNamedNativeMethod("sendDatagram", 2, func(loop *tsLoop, receiver tsValue, args ...tsValue) tsValue {
		return tsNodeQUICSendDatagram(loop, receiver, tsArg(args, 0), tsArg(args, 1))
	}))
 sp.set("updateKey", tsNamedNativeMethod("updateKey", 0, func(loop *tsLoop, receiver tsValue, _ ...tsValue) tsValue {
		if state := tsNodeQUICSessionFor(receiver); state == nil || state.closed { panic(tsThrown{tsNodeQUICError("ERR_INVALID_STATE", "The QUIC session is closed")}) }
		return tsU
	}))
 sp.set("close", tsNamedNativeMethod("close", 1, func(loop *tsLoop, receiver tsValue, args ...tsValue) tsValue { return tsNodeQUICGracefulCloseSession(loop, receiver, tsArg(args, 0)) }))
 sp.set("destroy", tsNamedNativeMethod("destroy", 2, func(loop *tsLoop, receiver tsValue, args ...tsValue) tsValue { return tsNodeQUICCloseSession(loop, receiver, tsArg(args, 0), tsArg(args, 1)) }))
	qp.set("destroy",tsNamedNativeMethod("destroy",2,func(loop *tsLoop,r tsValue,args ...tsValue)tsValue{return tsNodeQUICEndStream(loop,r,tsArg(args,0),tsArg(args,1))}))
	qp.set("setOutbound",tsNamedNativeMethod("setOutbound",1,func(loop *tsLoop,r tsValue,args ...tsValue)tsValue{value:=tsArg(args,0);valid:=value.kind==tsArrayBufferKind||value.kind==tsTypedArrayKind||value.kind==tsDataViewKind;if value.kind==tsObjectKind&&(*tsObject)(value.ref).blob!=nil{valid=true};if !valid{tsFSInvalid("outbound","an ArrayBuffer, ArrayBufferView, or Blob")};return tsNodeQUICSetOutbound(loop,r,value)}))
	qp.set("stopSending",tsNamedNativeMethod("stopSending",1,func(loop *tsLoop,r tsValue,args ...tsValue)tsValue{return tsNodeQUICStopSending(loop,r,tsArg(args,0))}))
	qp.set("resetStream",tsNamedNativeMethod("resetStream",1,func(loop *tsLoop,r tsValue,args ...tsValue)tsValue{return tsNodeQUICResetWriteSide(loop,r,tsArg(args,0))}))
	qp.set("sendHeaders", tsNamedNativeMethod("sendHeaders", 2, func(loop *tsLoop, r tsValue, args ...tsValue) tsValue {
		state := tsNodeQUICStreamFor(r)
		if state == nil || state.closed { return tsBooleanValue(false) }
		headers := tsArg(args, 0)
		if !tsNodeQUICObject(headers) { tsFSInvalid("headers", "an object") }
		options := tsArg(args, 1)
		if tsIsUndefined(options) { options = tsObjectValue(tsNewObject()) }
		if !tsNodeQUICObject(options) { tsFSInvalid("options", "an object") }
		terminal := tsNodeQUICGet(loop, options, "terminal")
		if tsIsUndefined(terminal) { terminal = tsBooleanValue(false) } else if terminal.kind != tsBooleanKind { tsFSInvalid("options.terminal", "a boolean") }
		return tsNodeQUICSendHeaders(loop, r, headers, tsTruthy(terminal), "initial")
	}))
	qp.set("sendInformationalHeaders", tsNamedNativeMethod("sendInformationalHeaders", 1, func(loop *tsLoop, r tsValue, args ...tsValue) tsValue {
		state := tsNodeQUICStreamFor(r)
		if state == nil || state.closed { return tsBooleanValue(false) }
		if session := tsNodeQUICSessionFor(state.session); session == nil || session.client { panic(tsThrown{tsNodeQUICError("ERR_INVALID_STATE", "Informational headers can only be sent by a server")}) }
		headers := tsArg(args, 0)
		if !tsNodeQUICObject(headers) { tsFSInvalid("headers", "an object") }
		return tsNodeQUICSendHeaders(loop, r, headers, false, "informational")
	}))
	qp.set("sendTrailers", tsNamedNativeMethod("sendTrailers", 1, func(loop *tsLoop, r tsValue, args ...tsValue) tsValue {
		state := tsNodeQUICStreamFor(r)
		if state == nil || state.closed { return tsBooleanValue(false) }
		headers := tsArg(args, 0)
		if !tsNodeQUICObject(headers) { tsFSInvalid("headers", "an object") }
		return tsNodeQUICSendHeaders(loop, r, headers, false, "trailing")
	}))
	qp.set("setBody", tsNamedNativeMethod("setBody", 1, func(loop *tsLoop, r tsValue, args ...tsValue) tsValue {
		body := tsArg(args, 0)
		valid := body.kind == tsStringKind || body.kind == tsArrayBufferKind || body.kind == tsTypedArrayKind || body.kind == tsDataViewKind || body.kind == tsPromiseKind || tsIsNull(body) || tsNodeQUICIsFileHandle(body) || (body.kind == tsObjectKind && (*tsObject)(body.ref).blob != nil) || tsNodeIterFromIsSyncIterable(loop, body) || tsNodeIterFromIsAsyncIterable(loop, body)
		if !valid { tsNodeIterFromInvalid("body", "a string, ArrayBuffer, ArrayBufferView, Blob, Iterable, AsyncIterable, Promise, or null", body) }
		return tsNodeQUICSetOutbound(loop, r, body)
	}))
	qp.set("setPriority", tsNamedNativeMethod("setPriority", 1, func(loop *tsLoop, r tsValue, args ...tsValue) tsValue {
		state := tsNodeQUICStreamFor(r)
		if state == nil || state.closed { return tsU }
		if !tsTruthy(tsNodeQUICGet(loop, tsNodeQUICSessionStateObject(state.session), "isPrioritySupported")) { panic(tsThrown{tsNodeQUICError("ERR_INVALID_STATE", "The session does not support stream priority")}) }
		options := tsArg(args, 0)
		if tsIsUndefined(options) { options = tsObjectValue(tsNewObject()) }
		if !tsNodeQUICObject(options) { tsFSInvalid("options", "an object") }
		level := tsNodeQUICGet(loop, options, "level")
		if tsIsUndefined(level) { level = tsNodeQUICKey("default") }
		if level.kind != tsStringKind || (tsText(level) != "default" && tsText(level) != "low" && tsText(level) != "high") { tsFSInvalid("options.level", "'default', 'low', or 'high'") }
		incremental := tsNodeQUICGet(loop, options, "incremental")
		if tsIsUndefined(incremental) { incremental = tsBooleanValue(false) } else if incremental.kind != tsBooleanKind { tsFSInvalid("options.incremental", "a boolean") }
		previousLevel := tsNodeQUICGet(loop, state.priority, "level")
		previousIncremental := tsNodeQUICGet(loop, state.priority, "incremental")
		changed := !tsStrictEqual(previousLevel, level) || !tsStrictEqual(previousIncremental, incremental)
		priority := tsNewObject()
		priority.set("level", level)
		priority.set("incremental", incremental)
		state.priority = tsObjectValue(priority)
		urgency := uint64(3)
		switch tsText(level) { case "high": urgency = 0; case "low": urgency = 7 }
		if changed { tsNodeQUICSendPriorityUpdate(r, urgency, tsTruthy(incremental)) }
		return tsU
	}))
}
func tsNodeQUICValidateHeaders(value tsValue){if value.kind!=tsObjectKind&&value.kind!=tsArrayKind{tsFSInvalid("headers","an object or array")};if value.kind==tsObjectKind{for _,name:=range (*tsObject)(value.ref).order{if name=="connection"||name=="upgrade"||name=="keep-alive"||name=="transfer-encoding"||name=="proxy-connection"{panic(tsThrown{tsNodeQUICError("ERR_HTTP2_INVALID_CONNECTION_HEADERS","Connection-specific headers are forbidden")})};if strings.ContainsAny(name,"\r\n\x00"){panic(tsThrown{tsNodeQUICError("ERR_INVALID_HTTP_TOKEN","Invalid header name")})}}}}
func tsNodeQUICWriter(loop *tsLoop, stream tsValue) tsValue {
	state := tsNodeQUICStreamFor(stream)
	if state != nil && !tsIsUndefined(state.writer) {
		return state.writer
	}
	if state != nil && state.outboundSet { panic(tsThrown{tsNodeQUICError("ERR_INVALID_STATE", "Stream outbound already configured with a body source")}) }
	out := tsNewObject()
	for _, name := range []string{"write", "writev", "writeSync", "writevSync", "end", "endSync", "fail"} {
		methodName := name
		out.set(methodName, tsNamedNativeMethod(methodName, 2, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
			switch methodName {
			case "write":
				if state == nil || state.closed || state.writeEnded { return tsPromiseValue(loop.resolved(tsNodeQUICError("ERR_INVALID_STATE", "Writer is closed"), true)) }
				if state.writerErrored { return tsPromiseValue(loop.resolved(state.writerError, true)) }
				if wire := tsNodeQUICWireStateFor(stream); wire == nil || wire.transport == nil || tsNodeQUICWriterAvailable(state, wire.transport) == 0 { return tsPromiseValue(loop.resolved(tsNodeQUICError("ERR_INVALID_STATE", "Stream write buffer is full"), true)) }
				promise := loop.promise()
				completion := tsFunctionValue(tsFunc(func(_ *tsLoop, callbackArgs ...tsValue) tsValue {
					problem := tsArg(callbackArgs, 0)
					promise.settle(tsResult{problem, !tsNullish(problem) && !tsIsUndefined(problem)})
					return tsU
				}))
				tsNodeQUICWriteTransport(loop, stream, tsArg(args, 0), tsU, completion)
				return tsPromiseValue(promise)
			case "writeSync":
				return tsNodeQUICSyncWrite(loop, stream, tsNodeQUICInputBytes(loop, tsArg(args, 0)))
			case "writev", "writevSync":
				chunks := tsArg(args, 0)
				if chunks.kind != tsArrayKind {
					if methodName == "writevSync" { return tsBooleanValue(false) }
					return tsPromiseValue(loop.resolved(tsFSException("TypeError", "ERR_INVALID_ARG_TYPE", "The chunks argument must be an array"), true))
				}
				array := (*tsArray)(chunks.ref)
				data := make([]byte, 0)
				for index := 0; index < array.length(); index++ { chunk := array.at(index); if chunk.kind != tsStringKind && chunk.kind != tsTypedArrayKind && chunk.kind != tsDataViewKind && chunk.kind != tsArrayBufferKind { if methodName == "writevSync" { return tsBooleanValue(false) }; return tsPromiseValue(loop.resolved(tsFSException("TypeError", "ERR_INVALID_ARG_TYPE", "Each chunk must be a string or Uint8Array"), true)) }; data = append(data, tsNodeQUICInputBytes(loop, chunk)...) }
				if methodName == "writev" { if state == nil || state.closed || state.writeEnded { return tsPromiseValue(loop.resolved(tsNodeQUICError("ERR_INVALID_STATE", "Writer is closed"), true)) }; if state.writerErrored { return tsPromiseValue(loop.resolved(state.writerError, true)) }; if wire := tsNodeQUICWireStateFor(stream); wire == nil || wire.transport == nil || tsNodeQUICWriterAvailable(state, wire.transport) == 0 { return tsPromiseValue(loop.resolved(tsNodeQUICError("ERR_INVALID_STATE", "Stream write buffer is full"), true)) } }
				if methodName == "writevSync" { return tsNodeQUICSyncWrite(loop, stream, data) }
				promise := loop.promise()
				completion := tsFunctionValue(tsFunc(func(_ *tsLoop, callbackArgs ...tsValue) tsValue { problem := tsArg(callbackArgs, 0); promise.settle(tsResult{problem, !tsNullish(problem) && !tsIsUndefined(problem)}); return tsU }))
				tsNodeQUICWriteTransport(loop, stream, tsNodeBuffer(data), tsU, completion)
				return tsPromiseValue(promise)
			case "endSync":
				return tsNodeQUICSyncEnd(loop, stream)
			case "end":
				promise := loop.promise()
				completion := tsFunctionValue(tsFunc(func(_ *tsLoop, callbackArgs ...tsValue) tsValue {
					problem := tsArg(callbackArgs, 0)
					if !tsNullish(problem) && !tsIsUndefined(problem) { promise.settle(tsResult{problem, true}); return tsU }
					if current := tsNodeQUICWireStateFor(stream); current != nil && current.transport != nil { promise.settle(tsResult{tsNumberValue(float64(current.transport.totalBytes)), false}) } else { promise.settle(tsResult{tsNumberValue(0), false}) }
					return tsU
				}))
				tsNodeQUICFinalizeTransportWrite(loop, stream, completion)
				return tsPromiseValue(promise)
			case "fail":
				reason := tsArg(args, 0)
				if tsIsUndefined(reason) { reason = tsNodeQUICError("ERR_INVALID_STATE", "Failed") }
				if state != nil { state.writerErrored = true; state.writerError = reason; for _, pending := range state.drains { pending.settle(tsResult{reason, true}) }; state.drains = nil }
				tsNodeQUICResetTransportStream(loop, stream, reason)
				return tsNodeQUICEndStream(loop, stream, reason)
			default:
				return tsU
			}
		}))
	}
	tsNodeQUICDefineGetter(out, "desiredSize", func(_ *tsLoop, _ tsValue) tsValue {
		if state == nil || state.closed || state.writerErrored { return tsNull }
		current := tsNodeQUICWireStateFor(stream)
		if current == nil || current.transport == nil || current.writeEnded { return tsNull }
		return tsNumberValue(float64(tsNodeQUICWriterAvailable(state, current.transport)))
	}, nil)
	types := tsInternalStreamIterTypesModule(loop)
	drainable := tsNodeQUICGet(loop, types, "drainableProtocol")
	if !tsIsUndefined(drainable) {
		out.set(tsPropertyKey(drainable), tsNamedNativeMethod("[drainableProtocol]", 0, func(loop *tsLoop, _ tsValue, _ ...tsValue) tsValue {
			if state == nil || state.closed || state.writerErrored { return tsNull }
			if wire := tsNodeQUICWireStateFor(stream); wire == nil || wire.transport == nil { return tsNull } else if tsNodeQUICWriterAvailable(state, wire.transport) > 0 { return tsNull }
			pending := loop.promise()
			state.drains = append(state.drains, pending)
			return tsPromiseValue(pending)
		}))
	}
	asyncDispose := tsPropertyKey(tsWellKnownSymbol("asyncDispose"))
	out.set(asyncDispose, tsNamedNativeMethod("[Symbol.asyncDispose]", 0, func(loop *tsLoop, _ tsValue, _ ...tsValue) tsValue {
		if state != nil && !state.writeEnded && !state.writerErrored {
			reason := tsNodeQUICError("ERR_INVALID_STATE", "Failed")
			state.writerErrored, state.writerError = true, reason
			tsNodeQUICResetTransportStream(loop, stream, reason)
			tsNodeQUICEndStream(loop, stream, reason)
		}
		return tsPromiseValue(loop.resolved(tsU, false))
	}))
	dispose := tsPropertyKey(tsWellKnownSymbol("dispose"))
	out.set(dispose, tsNamedNativeMethod("[Symbol.dispose]", 0, func(loop *tsLoop, _ tsValue, _ ...tsValue) tsValue {
		if state != nil && !state.writeEnded && !state.writerErrored {
			reason := tsNodeQUICError("ERR_INVALID_STATE", "Failed")
			state.writerErrored, state.writerError = true, reason
			tsNodeQUICResetTransportStream(loop, stream, reason)
			tsNodeQUICEndStream(loop, stream, reason)
		}
		return tsU
	}))
	value := tsObjectValue(out)
	if state != nil { state.writer = value }
	return value
}
func tsNodeQUICInit(loop *tsLoop){tsNodeQUICTransportInit();tsNodeQUICSingleton.Do(func(){tsNodeQUICSingleton.endpoints=map[*tsObject]*tsNodeQUICEndpointState{};tsNodeQUICSingleton.sessions=map[*tsObject]*tsNodeQUICSessionState{};tsNodeQUICSingleton.streams=map[*tsObject]*tsNodeQUICStreamState{};tsNodeQUICSingleton.channels=map[string]tsValue{};tsNodeQUICSingleton.symbolValues=map[string]tsValue{};tsNodeQUICSingleton.private=tsNewSymbol(tsStringUTF8("kPrivateConstructor"));tsNodeQUICSingleton.createDisconnected=tsNewSymbol(tsStringUTF8("kCreateDisconnected"));tsNodeQUICSingleton.endpointProto=tsNewObject();tsNodeQUICSingleton.sessionProto=tsNewObject();tsNodeQUICSingleton.streamProto=tsNewObject();tsNodeQUICSingleton.errorProto=tsNewObject();tsNodeQUICSingleton.streamProto.prototype=tsObjectValue(func()*tsObject{tsNodeDuplexInit(loop);return tsNodeDuplexSingleton.prototype}());tsNodeQUICInstall(loop)
 ep:=tsNodeQUICMakeConstructor("QuicEndpoint",tsNodeQUICSingleton.endpointProto,1,func(loop *tsLoop,args []tsValue)tsValue{return tsNodeQUICCreateEndpoint(loop,tsArg(args,0))});tsNodeQUICSingleton.endpointClass=ep
 sc:=tsNodeQUICMakeConstructor("QuicSession",tsNodeQUICSingleton.sessionProto,3,func(_ *tsLoop,_ []tsValue)tsValue{panic(tsThrown{tsFSException("TypeError","ERR_ILLEGAL_CONSTRUCTOR","Illegal constructor")})});tsNodeQUICSingleton.sessionClass=sc
 qc:=tsNodeQUICMakeConstructor("QuicStream",tsNodeQUICSingleton.streamProto,5,func(_ *tsLoop,_ []tsValue)tsValue{panic(tsThrown{tsFSException("TypeError","ERR_ILLEGAL_CONSTRUCTOR","Illegal constructor")})});tsNodeQUICSingleton.streamClass=qc
 ec:=tsNodeQUICMakeConstructor("QuicError",tsNodeQUICSingleton.errorProto,2,func(loop *tsLoop,args []tsValue)tsValue{
		messageValue:=tsArg(args,0);if messageValue.kind!=tsStringKind{tsFSInvalid("message","a string")}
		options:=tsArg(args,1);if tsIsUndefined(options){options=tsObjectValue(tsNewObject())};if !tsNodeQUICObject(options){tsFSInvalid("options","an object")}
		errorCode:=tsNodeQUICGet(loop,options,"errorCode");if tsIsUndefined(errorCode){panic(tsThrown{tsFSException("TypeError","ERR_MISSING_ARGS","The \"options.errorCode\" argument must be specified")})}
		if errorCode.kind!=tsBigIntKind&&!tsIsNumeric(errorCode){tsFSInvalid("options.errorCode","a number or bigint")}
		numericCode:=new(big.Int);if errorCode.kind==tsBigIntKind{numericCode.Set((*big.Int)(errorCode.ref))}else{number:=tsNumber(errorCode);if math.IsNaN(number)||math.Trunc(number)!=number||number<0||number>=math.Ldexp(1,62){tsFSInvalid("options.errorCode","a non-negative 62-bit integer")};numericCode.SetUint64(uint64(number))}
		if numericCode.Sign()<0||numericCode.BitLen()>62{tsFSInvalid("options.errorCode","a non-negative 62-bit integer")}
		code:=tsNodeQUICGet(loop,options,"code");if tsIsUndefined(code){code=tsNodeQUICKey("ERR_QUIC_STREAM_ABORTED")};if code.kind!=tsStringKind{tsFSInvalid("options.code","a string")}
		typeValue:=tsNodeQUICGet(loop,options,"type");if tsIsUndefined(typeValue){typeValue=tsNodeQUICKey("application")};kind:=tsText(typeValue);if kind!="transport"&&kind!="application"{tsFSInvalid("options.type","'transport' or 'application'")}
		o:=tsNewObject();o.prototype=tsObjectValue(tsNodeQUICSingleton.errorProto);o.set("name",tsNodeQUICKey("QuicError"));o.set("message",messageValue);o.set("code",code);o.set("type",typeValue);o.set("errorCode",tsBigIntValue(numericCode));if tsNodeQUICSingleton.quicErrors==nil{tsNodeQUICSingleton.quicErrors=map[*tsObject]bool{}};tsNodeQUICSingleton.quicErrors[o]=true;return tsObjectValue(o)
	});tsNodeQUICSingleton.errorClass=ec
	tsSet(loop,ec,tsNodeQUICKey("isQuicError"),tsNamedNativeMethod("isQuicError",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{value:=tsArg(args,0);return tsBooleanValue(value.kind==tsObjectKind&&tsNodeQUICSingleton.quicErrors[(*tsObject)(value.ref)])}))
 for _,name:=range []string{"QuicEndpoint","QuicSession","QuicStream"}{var ctor tsValue;switch name{case "QuicEndpoint":ctor=ep;case "QuicSession":ctor=sc;default:ctor=qc};stats:=tsNodeQUICMakeStatsClass(name+"Stats");tsSet(loop,ctor,tsNodeQUICKey("Stats"),stats)}
 tsNodeQUICSingleton.module=tsNodeQUICBuildCore(loop)
})}
func tsNodeQUICMakeStatsClass(name string)tsValue{proto:=tsNewObject();ctor:=tsNodeQUICMakeConstructor(name,proto,2,func(loop *tsLoop,args []tsValue)tsValue{o:=tsNewObject();o.prototype=tsObjectValue(proto);if len(args)>1&&tsArg(args,1).kind==tsObjectKind{for _,n:=range (*tsObject)(args[1].ref).order{o.set(n,(*tsObject)(args[1].ref).get(n))}};for _,method:=range []string{"toJSON","toString"}{o.set(method,proto.get(method))};return tsObjectValue(o)});proto.set("toJSON",tsNamedNativeMethod("toJSON",0,func(loop *tsLoop,r tsValue,_ ...tsValue)tsValue{return tsNodeQUICCloneProperties(loop,r)}));proto.set("toString",tsNamedNativeMethod("toString",0,func(loop *tsLoop,r tsValue,_ ...tsValue)tsValue{return tsJSONBuiltin(loop,"stringify",[]tsValue{tsNodeQUICCloneProperties(loop,r)})}));return ctor}
func tsNodeQUICBuildCore(loop *tsLoop)tsValue{out:=tsNewObject();out.set("connect", tsNamedNativeMethod("connect", 2, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
		return tsNodeQUICConnect(loop, tsArg(args, 0), tsArg(args, 1))
	}))
	out.set("listen", tsNamedNativeMethod("listen", 2, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
		return tsNodeQUICListen(loop, tsArg(args, 0), tsArg(args, 1))
	}))
	out.set("QuicEndpoint",tsNodeQUICSingleton.endpointClass);out.set("QuicSession",tsNodeQUICSingleton.sessionClass);out.set("QuicStream",tsNodeQUICSingleton.streamClass);out.set("QuicError",tsNodeQUICSingleton.errorClass);out.set("CC_ALGO_RENO",tsNodeQUICKey("reno"));out.set("CC_ALGO_CUBIC",tsNodeQUICKey("cubic"));out.set("CC_ALGO_BBR",tsNodeQUICKey("bbr"));out.set("DEFAULT_CIPHERS",tsNodeQUICKey("TLS_AES_128_GCM_SHA256:TLS_AES_256_GCM_SHA384:TLS_CHACHA20_POLY1305_SHA256:TLS_AES_128_CCM_SHA256"));out.set("DEFAULT_GROUPS",tsNodeQUICKey("X25519:P-256:P-384:P-521"));out.set("getQuicStreamState",tsNamedNativeMethod("getQuicStreamState",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeQUICStreamStateObject(tsArg(args,0))}));out.set("getQuicSessionState",tsNamedNativeMethod("getQuicSessionState",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeQUICSessionStateObject(tsArg(args,0))}));out.set("getQuicEndpointState",tsNamedNativeMethod("getQuicEndpointState",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeQUICEndpointStateObject(tsArg(args,0))}));out.set("default",tsObjectValue(out));return tsObjectValue(out)}
func tsNodeQUICModule(loop *tsLoop)tsValue{tsNodeQUICInit(loop);return tsNodeQUICPublicModule(loop)}
func tsNodeQUICPublicModule(loop *tsLoop)tsValue{core:=tsNodeQUICSingleton.module;if core.kind!=tsObjectKind{tsNodeQUICInit(loop);core=tsNodeQUICSingleton.module};out:=tsNewObject();for _,name:=range []string{"connect","listen","QuicEndpoint","QuicError","QuicSession","QuicStream"}{out.set(name,tsNodeQUICGet(loop,core,name))};cc:=tsNewObject();cc.set("RENO",tsNodeQUICGet(loop,core,"CC_ALGO_RENO"));cc.set("CUBIC",tsNodeQUICGet(loop,core,"CC_ALGO_CUBIC"));cc.set("BBR",tsNodeQUICGet(loop,core,"CC_ALGO_BBR"));constants:=tsNewObject();constants.set("cc",tsObjectValue(cc));constants.set("DEFAULT_CIPHERS",tsNodeQUICGet(loop,core,"DEFAULT_CIPHERS"));constants.set("DEFAULT_GROUPS",tsNodeQUICGet(loop,core,"DEFAULT_GROUPS"));out.set("constants",tsObjectValue(constants));out.set("default",tsObjectValue(out));return tsObjectValue(out)}
func tsNodeQUICBindingStateOffsets(binding *tsObject) {
	values := map[string]float64{
		"IDX_STATE_ENDPOINT_BOUND": 0, "IDX_STATE_ENDPOINT_BOUND_SIZE": 1,
		"IDX_STATE_ENDPOINT_RECEIVING": 1, "IDX_STATE_ENDPOINT_RECEIVING_SIZE": 1,
		"IDX_STATE_ENDPOINT_LISTENING": 2, "IDX_STATE_ENDPOINT_LISTENING_SIZE": 1,
		"IDX_STATE_ENDPOINT_CLOSING": 3, "IDX_STATE_ENDPOINT_CLOSING_SIZE": 1,
		"IDX_STATE_ENDPOINT_BUSY": 4, "IDX_STATE_ENDPOINT_BUSY_SIZE": 1,
		"IDX_STATE_ENDPOINT_MAX_CONNECTIONS_PER_HOST": 6, "IDX_STATE_ENDPOINT_MAX_CONNECTIONS_PER_HOST_SIZE": 2,
		"IDX_STATE_ENDPOINT_MAX_CONNECTIONS_TOTAL": 8, "IDX_STATE_ENDPOINT_MAX_CONNECTIONS_TOTAL_SIZE": 2,
		"IDX_STATE_ENDPOINT_PENDING_CALLBACKS": 16, "IDX_STATE_ENDPOINT_PENDING_CALLBACKS_SIZE": 8,
		"IDX_STATE_SESSION_LISTENER_FLAGS": 0, "IDX_STATE_SESSION_CLOSING": 4,
		"IDX_STATE_SESSION_GRACEFUL_CLOSE": 5, "IDX_STATE_SESSION_SILENT_CLOSE": 6,
		"IDX_STATE_SESSION_STATELESS_RESET": 7, "IDX_STATE_SESSION_HANDSHAKE_COMPLETED": 8,
		"IDX_STATE_SESSION_HANDSHAKE_CONFIRMED": 9, "IDX_STATE_SESSION_STREAM_OPEN_ALLOWED": 10,
		"IDX_STATE_SESSION_PRIORITY_SUPPORTED": 11, "IDX_STATE_SESSION_HEADERS_SUPPORTED": 12,
		"IDX_STATE_SESSION_WRAPPED": 13, "IDX_STATE_SESSION_APPLICATION_TYPE": 14,
		"IDX_STATE_SESSION_NO_ERROR_CODE": 16, "IDX_STATE_SESSION_INTERNAL_ERROR_CODE": 24,
		"IDX_STATE_SESSION_MAX_DATAGRAM_SIZE": 32, "IDX_STATE_SESSION_LAST_DATAGRAM_ID": 40,
		"IDX_STATE_SESSION_MAX_PENDING_DATAGRAMS": 48,
		"IDX_STATE_STREAM_ID": 0, "IDX_STATE_STREAM_PENDING": 8,
		"IDX_STATE_STREAM_FIN_SENT": 9, "IDX_STATE_STREAM_FIN_RECEIVED": 10,
		"IDX_STATE_STREAM_READ_ENDED": 11, "IDX_STATE_STREAM_WRITE_ENDED": 12,
		"IDX_STATE_STREAM_RESET": 13, "IDX_STATE_STREAM_RESET_CODE": 16,
		"IDX_STATE_STREAM_HAS_OUTBOUND": 24, "IDX_STATE_STREAM_HAS_READER": 25,
		"IDX_STATE_STREAM_WANTS_BLOCK": 26, "IDX_STATE_STREAM_WANTS_HEADERS": 27,
		"IDX_STATE_STREAM_WANTS_RESET": 28, "IDX_STATE_STREAM_WANTS_TRAILERS": 29,
		"IDX_STATE_STREAM_RECEIVED_EARLY_DATA": 30, "IDX_STATE_STREAM_WRITE_DESIRED_SIZE": 32,
		"IDX_STATE_STREAM_HIGH_WATER_MARK": 36,
	}
	for name, value := range values {
		binding.set(name, tsNumberValue(value))
	}
}

func tsNodeQUICBindingSetCallbacks(loop *tsLoop, callbacks tsValue) tsValue {
	if tsNodeQUICSingleton.callbacksInitialized {
		return tsU
	}
	if !tsNodeQUICObject(callbacks) {
		panic(tsThrown{tsFSException("TypeError", "ERR_INVALID_ARG_TYPE", "The \"callbacks\" argument must be of type object")})
	}
	for _, name := range []string{
		"onEndpointClose", "onSessionClose", "onSessionEarlyDataRejected", "onSessionGoaway",
		"onSessionDatagram", "onSessionDatagramStatus", "onSessionHandshake", "onSessionKeyLog",
		"onSessionQlog", "onSessionNew", "onSessionNewToken", "onSessionOrigin",
		"onSessionPathValidation", "onSessionTicket", "onSessionVersionNegotiation",
		"onStreamBlocked", "onStreamClose", "onStreamCreated", "onStreamDrain",
		"onStreamHeaders", "onStreamReset", "onStreamTrailers",
	} {
		if tsNodeQUICGet(loop, callbacks, name).kind != tsFunctionKind {
			panic(tsThrown{tsFSException("TypeError", "ERR_MISSING_ARGS", "Missing Callback: " + name)})
		}
	}
	tsNodeQUICSingleton.callbacksInitialized = true
	return tsU
}

func tsNodeQUICBindingModule(loop *tsLoop)tsValue{tsNodeQUICInit(loop);if tsNodeQUICSingleton.binding.kind==tsObjectKind{return tsNodeQUICSingleton.binding};o:=tsNewObject();o.set("Endpoint",tsNodeQUICSingleton.endpointClass);o.set("setCallbacks",tsNamedNativeMethod("setCallbacks",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeQUICBindingSetCallbacks(loop,tsArg(args,0))}));constantValues:=map[string]float64{"CC_ALGO_RENO":0,"CC_ALGO_CUBIC":1,"CC_ALGO_BBR":2,"PREFERRED_ADDRESS_IGNORE":0,"PREFERRED_ADDRESS_USE":1,"DEFAULT_PREFERRED_ADDRESS_POLICY":0,"STREAM_DIRECTION_BIDIRECTIONAL":0,"STREAM_DIRECTION_UNIDIRECTIONAL":1,"CLOSECONTEXT_CLOSE":0,"CLOSECONTEXT_BIND_FAILURE":1,"CLOSECONTEXT_LISTEN_FAILURE":5,"CLOSECONTEXT_RECEIVE_FAILURE":3,"CLOSECONTEXT_SEND_FAILURE":4,"CLOSECONTEXT_START_FAILURE":2,"QUIC_STREAM_HEADERS_KIND_HINTS":0,"QUIC_STREAM_HEADERS_KIND_INITIAL":1,"QUIC_STREAM_HEADERS_KIND_TRAILING":2,"QUIC_STREAM_HEADERS_FLAGS_NONE":0,"QUIC_STREAM_HEADERS_FLAGS_TERMINAL":1,"DEFAULT_MAX_SOCKETADDRESS_LRU_SIZE":1024,"DEFAULT_RETRYTOKEN_EXPIRATION":10,"DEFAULT_REGULARTOKEN_EXPIRATION":10,"DEFAULT_MAX_PACKET_LENGTH":1200,"DEFAULT_MAX_STREAM_DATA":262144,"DEFAULT_MAX_DATA":1048576,"DEFAULT_MAX_IDLE_TIMEOUT":10,"DEFAULT_MAX_STREAMS_BIDI":100,"DEFAULT_MAX_STREAMS_UNI":3,"DEFAULT_ACTIVE_CONNECTION_ID_LIMIT":2,"DEFAULT_MAX_HEADER_LIST_PAIRS":128,"DEFAULT_MAX_HEADER_LENGTH":8192,"QUIC_PROTO_MAX":1,"QUIC_PROTO_MIN":1,"DEFAULT_HANDSHAKE_TIMEOUT":10000};for name,value:=range constantValues{o.set(name,tsNumberValue(value))};tsNodeQUICBindingStateOffsets(o);for _,name:=range []string{"CC_ALGO_RENO_STR","CC_ALGO_CUBIC_STR","CC_ALGO_BBR_STR","DEFAULT_CIPHERS","DEFAULT_GROUPS"}{value:=map[string]string{"CC_ALGO_RENO_STR":"reno","CC_ALGO_CUBIC_STR":"cubic","CC_ALGO_BBR_STR":"bbr","DEFAULT_CIPHERS":"TLS_AES_128_GCM_SHA256:TLS_AES_256_GCM_SHA384:TLS_CHACHA20_POLY1305_SHA256:TLS_AES_128_CCM_SHA256","DEFAULT_GROUPS":"X25519:P-256:P-384:P-521"}[name];o.set(name,tsNodeQUICKey(value))};for prefix,names:=range map[string][]string{"IDX_STATS_ENDPOINT_":{"CREATED_AT","DESTROYED_AT","BYTES_RECEIVED","BYTES_SENT","PACKETS_RECEIVED","PACKETS_SENT","SERVER_SESSIONS","CLIENT_SESSIONS","SERVER_BUSY_COUNT","RETRY_COUNT","RETRY_RATE_LIMITED","VERSION_NEGOTIATION_COUNT","VERSION_NEGOTIATION_RATE_LIMITED","STATELESS_RESET_COUNT","STATELESS_RESET_RATE_LIMITED","IMMEDIATE_CLOSE_COUNT","IMMEDIATE_CLOSE_RATE_LIMITED","SESSION_CREATION_RATE_LIMITED","PACKETS_BLOCKED"},"IDX_STATS_SESSION_":{"CREATED_AT","CLOSING_AT","DESTROYED_AT","HANDSHAKE_COMPLETED_AT","HANDSHAKE_CONFIRMED_AT","GRACEFUL_CLOSING_AT","BYTES_RECEIVED","BYTES_SENT","BIDI_IN_STREAM_COUNT","BIDI_OUT_STREAM_COUNT","UNI_IN_STREAM_COUNT","UNI_OUT_STREAM_COUNT","KEY_UPDATE_COUNT","LOSS_RETRANSMIT_COUNT","MAX_BYTES_IN_FLIGHT","BYTES_IN_FLIGHT","BLOCK_COUNT","CWND","LATEST_RTT","MIN_RTT","RTTVAR","SMOOTHED_RTT","SSTHRESH","DATAGRAMS_RECEIVED","DATAGRAMS_SENT","DATAGRAMS_ACKNOWLEDGED","DATAGRAMS_LOST","STREAMS_IDLE_TIMED_OUT","PKT_SENT","PKT_RECV","PKT_LOST","BYTES_RECV","BYTES_LOST","PING_RECV","PKT_DISCARDED"},"IDX_STATS_STREAM_":{"CREATED_AT","OPENED_AT","RECEIVED_AT","ACKED_AT","DESTROYED_AT","BYTES_RECEIVED","BYTES_SENT","MAX_OFFSET","MAX_OFFSET_ACK","MAX_OFFSET_RECV","FINAL_SIZE","BYTES_ACCUMULATED","MAX_BYTES_ACCUMULATED"}}{for i,n:=range names{o.set(prefix+n,tsNumberValue(float64(i)))};o.set(prefix+"COUNT",tsNumberValue(float64(len(names))))};o.set("default",tsObjectValue(o));tsNodeQUICSingleton.binding=tsObjectValue(o);return tsNodeQUICSingleton.binding}
func tsNodeQUICObject(v tsValue)bool{return v.kind==tsObjectKind||v.kind==tsInstanceKind}
func tsInternalQUICSymbolsModule(_ *tsLoop)tsValue{if tsNodeQUICSingleton.symbols.kind==tsObjectKind{return tsNodeQUICSingleton.symbols};o:=tsNewObject();if tsNodeQUICSingleton.symbolValues==nil{tsNodeQUICSingleton.symbolValues=map[string]tsValue{}};names:=[]string{"kAttachFileHandle","kBlocked","kConnect","kDrain","kDatagram","kDatagramStatus","kEarlyDataRejected","kFinishClose","kGoaway","kHandshake","kHandshakeCompleted","kVerifyPeer","kHeaders","kKeylog","kListen","kQlog","kNewSession","kNewStream","kNewToken","kStreamCallbacks","kStreamIdleTimeout","kOrigin","kOwner","kPathValidation","kPrivateConstructor","kRemoveSession","kRemoveStream","kReset","kSessionTicket","kTrailers","kVersionNegotiation"};for _,name:=range names{value:=tsNewSymbol(tsStringUTF8(name));o.set(name,value);tsNodeQUICSingleton.symbolValues[name]=value};o.set("kInspect",tsNewSymbol(tsStringUTF8("nodejs.util.inspect.custom")));o.set("default",tsObjectValue(o));tsNodeQUICSingleton.symbols=tsObjectValue(o);return tsNodeQUICSingleton.symbols}
func tsInternalQUICDiagnosticsModule(loop *tsLoop)tsValue{if tsNodeQUICSingleton.diagnostics.kind==tsObjectKind{return tsNodeQUICSingleton.diagnostics};dc:=tsNodeDiagnosticsModule(loop);out:=tsNewObject();names:=map[string]string{"onEndpointCreatedChannel":"quic.endpoint.created","onEndpointListeningChannel":"quic.endpoint.listen","onEndpointClosingChannel":"quic.endpoint.closing","onEndpointClosedChannel":"quic.endpoint.closed","onEndpointErrorChannel":"quic.endpoint.error","onEndpointBusyChangeChannel":"quic.endpoint.busy.change","onEndpointClientSessionChannel":"quic.session.created.client","onEndpointServerSessionChannel":"quic.session.created.server","onSessionOpenStreamChannel":"quic.session.open.stream","onSessionReceivedStreamChannel":"quic.session.received.stream","onSessionSendDatagramChannel":"quic.session.send.datagram","onSessionUpdateKeyChannel":"quic.session.update.key","onSessionClosingChannel":"quic.session.closing","onSessionClosedChannel":"quic.session.closed","onSessionReceiveDatagramChannel":"quic.session.receive.datagram","onSessionReceiveDatagramStatusChannel":"quic.session.receive.datagram.status","onSessionPathValidationChannel":"quic.session.path.validation","onSessionNewTokenChannel":"quic.session.new.token","onSessionTicketChannel":"quic.session.ticket","onSessionVersionNegotiationChannel":"quic.session.version.negotiation","onSessionOriginChannel":"quic.session.receive.origin","onSessionHandshakeChannel":"quic.session.handshake","onSessionGoawayChannel":"quic.session.goaway","onSessionEarlyRejectedChannel":"quic.session.early.rejected","onStreamClosedChannel":"quic.stream.closed","onStreamHeadersChannel":"quic.stream.headers","onStreamTrailersChannel":"quic.stream.trailers","onStreamInfoChannel":"quic.stream.info","onStreamResetChannel":"quic.stream.reset","onStreamBlockedChannel":"quic.stream.blocked","onSessionErrorChannel":"quic.session.error","onEndpointConnectChannel":"quic.endpoint.connect"};for key,name:=range names{channel:=tsCall(loop,tsNodeQUICGet(loop,dc,"channel"),tsNodeQUICKey(name));out.set(key,channel)};out.set("default",tsObjectValue(out));tsNodeQUICSingleton.diagnostics=tsObjectValue(out);return tsNodeQUICSingleton.diagnostics}
func tsInternalQUICHTTP2UtilModule(_ *tsLoop)tsValue{if tsNodeQUICSingleton.http2util.kind==tsObjectKind{return tsNodeQUICSingleton.http2util};out:=tsNewObject();valid:=map[string]bool{":status":true,":method":true,":authority":true,":scheme":true,":path":true,":protocol":true};out.set("assertValidPseudoHeader",tsNamedNativeMethod("assertValidPseudoHeader",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{if !valid[tsText(tsArg(args,0))]{panic(tsThrown{tsNodeQUICError("ERR_HTTP2_INVALID_PSEUDOHEADER","Invalid HTTP/2 pseudo-header")})};return tsU}));out.set("assertValidPseudoHeaderTrailer",tsNamedNativeMethod("assertValidPseudoHeaderTrailer",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{panic(tsThrown{tsNodeQUICError("ERR_HTTP2_INVALID_PSEUDOHEADER","Pseudo-headers are not permitted in trailers")})}));out.set("buildNgHeaderString",tsNamedNativeMethod("buildNgHeaderString",3,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeQUICBuildHeaders(loop,tsArg(args,0),tsArg(args,1),tsArg(args,2),valid)}));out.set("default",tsObjectValue(out));tsNodeQUICSingleton.http2util=tsObjectValue(out);return tsNodeQUICSingleton.http2util}
func tsNodeQUICBuildHeaders(loop *tsLoop,source,validator,strict tsValue,valid map[string]bool)tsValue{pairs:=[][2]string{};if source.kind==tsArrayKind{a:=(*tsArray)(source.ref);for i:=0;i+1<len(a.values);i+=2{pairs=append(pairs,[2]string{tsText(a.values[i]),tsText(a.values[i+1])})}}else if tsNodeQUICObject(source){o:=(*tsObject)(source.ref);for _,n:=range o.order{pairs=append(pairs,[2]string{n,tsText(o.values[n])})}};var regular,pseudo strings.Builder;count:=0;for _,pair:=range pairs{name,value:=strings.ToLower(pair[0]),pair[1];if name==""{continue};if strings.ContainsAny(name,"\r\n\x00")||strings.ContainsAny(value,"\r\n\x00"){panic(tsThrown{tsFSException("TypeError","ERR_INVALID_CHAR","Invalid HTTP/2 header content")})};if strings.HasPrefix(name,":"){if !valid[name]{panic(tsThrown{tsNodeQUICError("ERR_HTTP2_INVALID_PSEUDOHEADER","Invalid HTTP/2 pseudo-header")})};if validator.kind==tsFunctionKind{result:=tsCall(loop,validator,tsNodeQUICKey(name),tsNodeQUICKey(value));if !tsIsUndefined(result)&&!tsNullish(result){panic(tsThrown{result})}};pseudo.WriteString(name);pseudo.WriteByte(0);pseudo.WriteString(value);pseudo.WriteByte(0);pseudo.WriteByte(0)}else{if name=="connection"||name=="upgrade"||name=="http2-settings"||name=="keep-alive"||name=="proxy-connection"||name=="transfer-encoding"{panic(tsThrown{tsNodeQUICError("ERR_HTTP2_INVALID_CONNECTION_HEADERS","HTTP/1 connection-specific headers are forbidden")})};regular.WriteString(name);regular.WriteByte(0);regular.WriteString(value);regular.WriteByte(0);regular.WriteByte(0)};count++};out:=&tsArray{values:[]tsValue{tsNodeQUICKey(pseudo.String()+regular.String()),tsNumberValue(float64(count))}};_ = strict;return tsArrayValue(out)}
func tsInternalQUICStatsModule(loop *tsLoop)tsValue{tsNodeQUICInit(loop);if tsNodeQUICSingleton.stats.kind==tsObjectKind{return tsNodeQUICSingleton.stats};out:=tsNewObject();out.set("QuicEndpointStats",tsNodeQUICMakeStatsClass("QuicEndpointStats"));out.set("QuicSessionStats",tsNodeQUICMakeStatsClass("QuicSessionStats"));out.set("QuicStreamStats",tsNodeQUICMakeStatsClass("QuicStreamStats"));out.set("kCreateDisconnected",tsNodeQUICSingleton.createDisconnected);out.set("default",tsObjectValue(out));tsNodeQUICSingleton.stats=tsObjectValue(out);return tsNodeQUICSingleton.stats}
func tsNodeQUICStateConstructor(name string,fields []string)tsValue{proto:=tsNewObject();ctor:=tsNodeQUICMakeConstructor(name,proto,2,func(loop *tsLoop,args []tsValue)tsValue{o:=tsNewObject();o.prototype=tsObjectValue(proto);initial:=tsArg(args,1);for _,field:=range fields{v:=tsU;if tsNodeQUICObject(initial){v=tsNodeQUICGet(loop,initial,field)};if tsIsUndefined(v){v=tsBooleanValue(false)};o.set(field,v)};o.set("toJSON",proto.get("toJSON"));o.set("toString",proto.get("toString"));return tsObjectValue(o)});proto.set("toJSON",tsNamedNativeMethod("toJSON",0,func(loop *tsLoop,r tsValue,_ ...tsValue)tsValue{return tsNodeQUICCloneProperties(loop,r)}));proto.set("toString",tsNamedNativeMethod("toString",0,func(loop *tsLoop,r tsValue,_ ...tsValue)tsValue{return tsJSONBuiltin(loop,"stringify",[]tsValue{tsNodeQUICCloneProperties(loop,r)})}));return ctor}
func tsInternalQUICStateModule(_ *tsLoop)tsValue{if tsNodeQUICSingleton.states.kind==tsObjectKind{return tsNodeQUICSingleton.states};out:=tsNewObject();out.set("QuicEndpointState",tsNodeQUICStateConstructor("QuicEndpointState",[]string{"isBound","isReceiving","isListening","isClosing","isBusy","maxConnectionsPerHost","maxConnectionsTotal","pendingCallbacks"}));out.set("QuicSessionState",tsNodeQUICStateConstructor("QuicSessionState",[]string{"hasPathValidationListener","hasDatagramListener","hasDatagramStatusListener","hasSessionTicketListener","hasNewTokenListener","hasOriginListener","isClosing","isGracefulClose","isSilentClose","isStatelessReset","isHandshakeCompleted","isHandshakeConfirmed","isStreamOpenAllowed","isPrioritySupported","headersSupported","isWrapped","applicationType","noErrorCode","internalErrorCode","maxDatagramSize","lastDatagramId","maxPendingDatagrams"}));out.set("QuicStreamState",tsNodeQUICStateConstructor("QuicStreamState",[]string{"id","pending","finSent","finReceived","readEnded","writeEnded","reset","hasOutbound","hasReader","wantsBlock","wantsHeaders","wantsReset","wantsTrailers","early","resetCode","writeDesiredSize","highWaterMark"}));out.set("default",tsObjectValue(out));tsNodeQUICSingleton.states=tsObjectValue(out);return tsNodeQUICSingleton.states}
func tsInternalQUICCoreModule(loop *tsLoop)tsValue{tsNodeQUICInit(loop);return tsNodeQUICSingleton.module}
func tsInternalQUICModule(loop *tsLoop,name string)tsValue{switch name{case "internal/quic/binding":return tsInternalQUICBindingModule(loop);case "internal/quic/symbols":return tsInternalQUICSymbolsModule(loop);case "internal/quic/diagnostics":return tsInternalQUICDiagnosticsModule(loop);case "internal/quic/http2util":return tsInternalQUICHTTP2UtilModule(loop);case "internal/quic/stats":return tsInternalQUICStatsModule(loop);case "internal/quic/state":return tsInternalQUICStateModule(loop);case "internal/quic/quic":return tsInternalQUICCoreModule(loop)};return tsU}
func tsInternalQUICBindingModule(loop *tsLoop)tsValue{return tsNodeQUICBindingModule(loop)}
func tsNodeQUICAddon(loop *tsLoop,name string)(tsValue,bool){if name=="quic"{module:=tsNodeQUICModule(loop);process:=tsNodeProcessModule(loop);warning:=tsNodeQUICGet(loop,process,"emitWarning");if warning.kind==tsFunctionKind{tsCall(loop,warning,tsNodeQUICKey("quic is an experimental feature and might change at any time"),tsNodeQUICKey("ExperimentalWarning"))};return module,true};return tsU,false}
`
