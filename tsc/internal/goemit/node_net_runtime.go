package goemit

// nodeNetRuntime adapts Node's Socket/Server API to Go's standard net package.
// The source's libuv/Bun socket handles are replaced at the host boundary.
const nodeNetRuntime = `
type tsNodeNetSocketState struct {
	loop *tsLoop
	stream tsValue
	conn net.Conn
	connecting, connected, closed, unref bool
	bytesRead, bytesWritten int64
	timeout time.Duration
	noDelay, keepAlive bool
	keepAliveDelay time.Duration
	typeOfService int
	typeOfServiceSet bool
	syncWriteFD int
	syncWrite, syncWriteAutoClose bool
	pendingData, pendingEncoding, pendingWriteCallback tsValue
	connectRace bool
	connectAddresses []tsNodeNetConnectAddress
	connectAttemptIndex int
	connectAttemptTimeout time.Duration
	connectErrors []tsValue
	connectCancel context.CancelFunc
	connectTrace *tsNodeNetSocketState
	reading, upgrading, tlsHandshaking bool
	traceConnectActive bool
	onreadEnabled bool
	onreadCallback, onreadBuffer, onreadBufferFactory tsValue
	onreadTail []byte
	onreadTailSet bool
	onreadEOF bool
	upgradeData []byte
	upgradeContinue func()
	abortDispose tsValue
	timeoutTask *tsTask
	server *tsNodeNetServerState
	localAddress, localFamily, remoteAddress, remoteFamily string
	localPort, remotePort int
}
type tsNodeNetConnectAddress struct { ip string; family int }
type tsNodeNetSocketAddressState struct { address, family string; port int; flowLabel uint32 }
type tsNodeNetServerState struct {
	loop *tsLoop
	server tsValue
	listener net.Listener
	clusterHandle tsValue
	clusterListenID uint64
	clusterListening bool
	connections int
	closing, closed, unref bool
	path, network string
	options tsValue
}
type tsNodeNetBlockRule struct { network *net.IPNet; ip, end net.IP; family, endFamily, kind, text string }
type tsNodeNetBlockListState struct { mu sync.RWMutex; rules []tsNodeNetBlockRule }
var tsNodeNetSingleton struct {
	sync.Once
	module, socketCtor, serverCtor, blockListCtor, socketAddressCtor tsValue
	socketProto, serverProto, blockListProto, socketAddressProto *tsObject
	sockets map[uintptr]*tsNodeNetSocketState
	servers map[uintptr]*tsNodeNetServerState
	blocks map[uintptr]*tsNodeNetBlockListState
	blockListMapMu sync.RWMutex
	socketAddresses map[uintptr]*tsNodeNetSocketAddressState
	socketAddressMapMu sync.RWMutex
	blockListCloneMu sync.Mutex
	blockListSerialized map[uint64]tsValue
	blockListObjectNonce map[uintptr]uint64
	blockListNextNonce atomic.Uint64
	defaultAutoSelectFamily bool
	defaultAutoSelectFamilyAttemptTimeout int
	normalizedArgsSymbol tsValue
	perfConnectKey tsValue
	tlsKeylogPath string
	tlsKeylogConfigured, tlsKeylogWarned bool
	traceEvents tsValue
	simultaneousAcceptsWarned bool
}
func tsNodeNetKey(name string) tsValue { return tsStringReference(tsStringUTF8(name)) }
func tsNodeNetGet(loop *tsLoop, value tsValue, name string) tsValue { return tsGet(loop,value,tsNodeNetKey(name)) }
func tsNodeNetSet(loop *tsLoop, value tsValue, name string, item tsValue) { tsSet(loop,value,tsNodeNetKey(name),item) }
func tsNodeNetPointer(value tsValue) uintptr { if value.kind!=tsObjectKind{return 0};return uintptr(value.ref) }
func tsNodeNetSocketFor(value tsValue) *tsNodeNetSocketState { return tsNodeNetSingleton.sockets[tsNodeNetPointer(value)] }
func tsNodeNetServerFor(value tsValue) *tsNodeNetServerState { return tsNodeNetSingleton.servers[tsNodeNetPointer(value)] }
func tsNodeNetEmit(loop *tsLoop, target tsValue, event string, args ...tsValue) tsValue {
	method:=tsNodeNetGet(loop,target,"emit");if method.kind!=tsFunctionKind{return tsBooleanValue(false)}
	values:=append([]tsValue{tsNodeNetKey(event)},args...);return tsCallReceiver(loop,method,target,values...)
}
func tsNodeNetReportError(loop *tsLoop,thrown any){
	errorValue:=tsUnwrap(thrown);process:=tsNodeProcessModule(loop);emit:=tsNodeNetGet(loop,process,"emit");handled:=emit.kind==tsFunctionKind&&tsTruthy(tsCallReceiver(loop,emit,process,tsNodeNetKey("uncaughtException"),errorValue,tsNodeNetKey("uncaughtException")));if !handled{loop.ready=append(loop.ready,func(){panic(tsThrown{errorValue})})}
}
func tsNodeNetTraceStart(loop *tsLoop,state *tsNodeNetSocketState,pipePath string){
	if state==nil||state.traceConnectActive{return};if tsIsUndefined(tsNodeNetSingleton.traceEvents){tsNodeNetSingleton.traceEvents=tsNodeModuleRequire(loop,[]tsValue{tsNodeNetKey("internal/trace_events")})};trace:=tsNodeNetSingleton.traceEvents;enabled:=tsNodeNetGet(loop,trace,"isCategoryGroupEnabled");if enabled.kind!=tsFunctionKind||!tsTruthy(tsCallReceiver(loop,enabled,trace,tsNodeNetKey("node,node.net,node.net.native"))){return};state.traceConnectActive=true
	if pipePath!=""{abstract:=strings.HasPrefix(pipePath,"\x00");if abstract{pipePath=strings.TrimPrefix(pipePath,"\x00")};data:=tsNewObject();kind:="file";if abstract{kind="abstract socket"};data.set("path_type",tsNodeNetKey(kind));data.set("pipe_path",tsNodeNetKey(pipePath));emit:=tsNodeNetGet(loop,trace,"emitEventWithArgs");if emit.kind==tsFunctionKind{tsCallReceiver(loop,emit,trace,tsNodeNetKey("b"),tsNodeNetKey("node,node.net,node.net.native"),tsNodeNetKey("connect"),tsU,tsObjectValue(data))}}else{emit:=tsNodeNetGet(loop,trace,"emitEvent");if emit.kind==tsFunctionKind{tsCallReceiver(loop,emit,trace,tsNodeNetKey("b"),tsNodeNetKey("node,node.net,node.net.native"),tsNodeNetKey("connect"))}}
}
func tsNodeNetTraceEnd(loop *tsLoop,state *tsNodeNetSocketState){if state==nil||!state.traceConnectActive{return};state.traceConnectActive=false;trace:=tsNodeNetSingleton.traceEvents;emit:=tsNodeNetGet(loop,trace,"emitEvent");if emit.kind==tsFunctionKind{tsCallReceiver(loop,emit,trace,tsNodeNetKey("e"),tsNodeNetKey("node,node.net,node.net.native"),tsNodeNetKey("connect"))}}
func tsNodeNetPerfConnectKey()tsValue{if tsIsUndefined(tsNodeNetSingleton.perfConnectKey){tsNodeNetSingleton.perfConnectKey=tsNewSymbol(tsStringUTF8("node.net.performance.connect"))};return tsNodeNetSingleton.perfConnectKey}
func tsNodeNetPerfStart(loop *tsLoop,receiver tsValue,host string,port int){shared:=tsInternalSharedModule(loop);hasObserver:=tsNodeNetGet(loop,shared,"hasObserver");if hasObserver.kind!=tsFunctionKind||!tsTruthy(tsCallReceiver(loop,hasObserver,shared,tsNodeNetKey("net"))){return};detail:=tsNewObject();detail.set("host",tsNodeNetKey(host));detail.set("port",tsNumberValue(float64(port)));context:=tsNewObject();context.set("type",tsNodeNetKey("net"));context.set("name",tsNodeNetKey("connect"));context.set("detail",tsObjectValue(detail));start:=tsNodeNetGet(loop,shared,"startPerf");if start.kind==tsFunctionKind{tsCallReceiver(loop,start,shared,receiver,tsNodeNetPerfConnectKey(),tsObjectValue(context))}}
func tsNodeNetPerfStop(loop *tsLoop,receiver tsValue){key:=tsNodeNetPerfConnectKey();if !tsTruthy(tsGet(loop,receiver,key)){return};shared:=tsInternalSharedModule(loop);hasObserver:=tsNodeNetGet(loop,shared,"hasObserver");if hasObserver.kind!=tsFunctionKind||!tsTruthy(tsCallReceiver(loop,hasObserver,shared,tsNodeNetKey("net"))){return};stop:=tsNodeNetGet(loop,shared,"stopPerf");if stop.kind==tsFunctionKind{tsCallReceiver(loop,stop,shared,receiver,key,tsU)}}
func tsNodeNetEmitSafe(loop *tsLoop,target tsValue,event string,args ...tsValue){
	defer func(){if thrown:=recover();thrown!=nil{tsNodeNetReportError(loop,thrown)}}()
	tsNodeNetEmit(loop,target,event,args...)
}
func tsNodeNetError(err error, operation string) tsValue {
	if err==nil{return tsNull}
	code:="ERR_SOCKET_ERROR";message:=err.Error()
	var op *net.OpError;if errors.As(err,&op){operation=op.Op;if op.Err!=nil{message=op.Err.Error()}}
	var errno syscall.Errno
	fields:=map[string]tsValue{"syscall":tsNodeNetKey(operation)}
	if errors.As(err,&errno){switch errno{case syscall.ECONNREFUSED:code="ECONNREFUSED";case syscall.ECONNRESET:code="ECONNRESET";case syscall.ECONNABORTED:code="ECONNABORTED";case syscall.EPIPE:code="EPIPE";case syscall.ETIMEDOUT:code="ETIMEDOUT";case syscall.EADDRINUSE:code="EADDRINUSE";case syscall.EADDRNOTAVAIL:code="EADDRNOTAVAIL";case syscall.EACCES:code="EACCES";case syscall.EINVAL:code="EINVAL";case syscall.ENOENT:code="ENOENT";case syscall.EBADF:code="EBADF";case syscall.ENOTCONN:code="ENOTCONN";case syscall.EHOSTUNREACH:code="EHOSTUNREACH";case syscall.ENETUNREACH:code="ENETUNREACH";case syscall.EAFNOSUPPORT:code="EAFNOSUPPORT";case syscall.EINTR:code="EINTR";case syscall.EIO:code="EIO"};fields["errno"]=tsNumberValue(-float64(errno))}else{var timeout net.Error;switch {case errors.As(err,&timeout)&&timeout.Timeout():code="ETIMEDOUT";case strings.Contains(message,"refused"):code="ECONNREFUSED";case strings.Contains(message,"reset"):code="ECONNRESET";case strings.Contains(message,"timed out"):code="ETIMEDOUT";case strings.Contains(message,"address already in use"):code="EADDRINUSE";case strings.Contains(message,"permission denied"):code="EACCES"}}
	fields["code"]=tsNodeNetKey(code);if code=="ETIMEDOUT"{if _,ok:=fields["errno"];!ok{fields["errno"]=tsNumberValue(-float64(syscall.ETIMEDOUT))}}
	return tsErrorValue(&tsRuntimeError{name:"Error",message:operation+" "+code+": "+message,fields:fields})
}
func tsNodeNetConnectError(loop *tsLoop,value tsValue,address string,port int,local tsValue,localPort int)tsValue{if tsNullish(value)||value.kind!=tsErrorKind{return value};tsNodeNetSet(loop,value,"address",tsNodeNetKey(address));tsNodeNetSet(loop,value,"port",tsNumberValue(float64(port)));runtimeError:=(*tsRuntimeError)(value.ref);message:="connect "+tsText(runtimeError.fields["code"])+" "+net.JoinHostPort(address,strconv.Itoa(port));if localPort>0{localAddress:=local;if tsNullish(localAddress){if ip:=net.ParseIP(strings.Split(address,"%")[0]);ip!=nil&&ip.To4()==nil{localAddress=tsNodeNetKey("::")}else{localAddress=tsNodeNetKey("0.0.0.0")}};tsNodeNetSet(loop,value,"localAddress",localAddress);tsNodeNetSet(loop,value,"localPort",tsNumberValue(float64(localPort)));message+=" - Local ("+net.JoinHostPort(tsText(localAddress),strconv.Itoa(localPort))+")"};runtimeError.message=message;return value}
func tsNodeNetAddress(addr net.Addr) (string,int,string) {
	if addr==nil{return "",0,""}
	switch value:=addr.(type){case *net.TCPAddr: family:="IPv6";if value.IP.To4()!=nil{family="IPv4"};return value.IP.String(),value.Port,family
	case *net.UDPAddr: family:="IPv6";if value.IP.To4()!=nil{family="IPv4"};return value.IP.String(),value.Port,family
	case *net.UnixAddr:return value.Name,0,"Unix"
	default:host,portText,err:=net.SplitHostPort(addr.String());if err==nil{port,_:=strconv.Atoi(portText);family:="IPv6";if ip:=net.ParseIP(host);ip==nil||ip.To4()!=nil{family="IPv4"};return host,port,family};return addr.String(),0,"Unix"}
}
func tsNodeNetSyncSocket(state *tsNodeNetSocketState) {
	if state==nil{return};stream:=state.stream
	tsNodeNetSet(state.loop,stream,"connecting",tsBooleanValue(state.connecting));tsNodeNetSet(state.loop,stream,"pending",tsBooleanValue(state.connecting));tsNodeNetSet(state.loop,stream,"destroyed",tsBooleanValue(state.closed));tsNodeNetSet(state.loop,stream,"bytesRead",tsNumberValue(float64(state.bytesRead)));tsNodeNetSet(state.loop,stream,"bytesWritten",tsNumberValue(float64(state.bytesWritten)))
	tsNodeNetSet(state.loop,stream,"localAddress",tsStringReference(tsStringUTF8(state.localAddress)));tsNodeNetSet(state.loop,stream,"localPort",tsNumberValue(float64(state.localPort)));tsNodeNetSet(state.loop,stream,"localFamily",tsStringReference(tsStringUTF8(state.localFamily)))
	tsNodeNetSet(state.loop,stream,"remoteAddress",tsStringReference(tsStringUTF8(state.remoteAddress)));tsNodeNetSet(state.loop,stream,"remotePort",tsNumberValue(float64(state.remotePort)));tsNodeNetSet(state.loop,stream,"remoteFamily",tsStringReference(tsStringUTF8(state.remoteFamily)))
	ready:="opening";if state.closed{ready="closed"}else if state.connected{ready="open"};tsNodeNetSet(state.loop,stream,"readyState",tsStringReference(tsStringUTF8(ready)))
}
func tsNodeNetResetTimeout(state *tsNodeNetSocketState) {
	if state==nil||state.closed{return};loop:=state.loop
	if state.timeoutTask!=nil{loop.clearTimeout(tsTaskValue(state.timeoutTask));state.timeoutTask=nil}
	if state.timeout<=0{return}
	var task *tsTask
	task=loop.setTimeout(tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{if state.timeoutTask==task&&!state.closed{state.timeoutTask=nil;if tsTruthy(tsNodeNetGet(loop,state.stream,"_pendingData")){return tsU};if writable:=tsNodeWritableStateFor(state.stream);writable!=nil&&writable.length>0{return tsU};tsNodeNetEmit(loop,state.stream,"timeout")};return tsU})),tsNumberValue(float64(state.timeout/time.Millisecond)))
	task.unref=state.unref;state.timeoutTask=task
}
func tsNodeNetUnderlyingConn(conn net.Conn)net.Conn{for{switch wrapped:=conn.(type){case interface{NetConn()net.Conn}:conn=wrapped.NetConn();case *tsNodeNetPrefixConn:conn=wrapped.Conn;default:return conn}}}
func tsNodeNetSetConn(state *tsNodeNetSocketState, conn net.Conn) {
	state.conn=conn;state.connecting=false;state.connected=true
	raw:=tsNodeNetUnderlyingConn(conn)
	if tcp,ok:=raw.(*net.TCPConn);ok{_ = tcp.SetNoDelay(state.noDelay);_ = tcp.SetKeepAlive(state.keepAlive);if state.keepAlive&&state.keepAliveDelay>0{_ = tcp.SetKeepAlivePeriod(state.keepAliveDelay)}};if state.typeOfServiceSet{_ = tsNodeNetSetSocketTypeOfService(conn,state.typeOfService)}
	state.localAddress,state.localPort,state.localFamily=tsNodeNetAddress(conn.LocalAddr());state.remoteAddress,state.remotePort,state.remoteFamily=tsNodeNetAddress(conn.RemoteAddr());tsNodeNetSyncSocket(state)
	state.loop.resourceClosers=append(state.loop.resourceClosers,func(){_ = conn.Close()})
	// Node's open handler arms the configured idle timeout when the handle becomes live.
	tsNodeNetResetTimeout(state)
}
func tsNodeNetSocketTypeOfServiceTarget(conn net.Conn)(*net.TCPConn,int,int,error){raw:=tsNodeNetUnderlyingConn(conn);tcp,ok:=raw.(*net.TCPConn);if !ok{return nil,0,0,syscall.ENOTCONN};remote,ok:=tcp.RemoteAddr().(*net.TCPAddr);if !ok{return nil,0,0,syscall.EAFNOSUPPORT};if remote.IP.To4()!=nil{return tcp,syscall.IPPROTO_IP,syscall.IP_TOS,nil};return tcp,syscall.IPPROTO_IPV6,syscall.IPV6_TCLASS,nil}
func tsNodeNetSetSocketTypeOfService(conn net.Conn,value int)error{tcp,level,option,err:=tsNodeNetSocketTypeOfServiceTarget(conn);if err!=nil{return err};raw,err:=tcp.SyscallConn();if err!=nil{return err};var socketErr error;controlErr:=raw.Control(func(fd uintptr){socketErr=syscall.SetsockoptInt(int(fd),level,option,value)});if controlErr!=nil{return controlErr};return socketErr}
func tsNodeNetGetSocketTypeOfService(conn net.Conn)(int,error){tcp,level,option,err:=tsNodeNetSocketTypeOfServiceTarget(conn);if err!=nil{return 0,err};raw,err:=tcp.SyscallConn();if err!=nil{return 0,err};value:=0;var socketErr error;controlErr:=raw.Control(func(fd uintptr){value,socketErr=syscall.GetsockoptInt(int(fd),level,option)});if controlErr!=nil{return 0,controlErr};return value,socketErr}
type tsNodeNetPrefixConn struct{net.Conn;prefix []byte}
func (conn *tsNodeNetPrefixConn)Read(buffer []byte)(int,error){if len(conn.prefix)>0{n:=copy(buffer,conn.prefix);conn.prefix=conn.prefix[n:];return n,nil};return conn.Conn.Read(buffer)}
func tsNodeNetDeliverData(loop *tsLoop,state *tsNodeNetSocketState,data []byte)(continueReading bool){
	continueReading=true;defer func(){if thrown:=recover();thrown!=nil{tsNodeNetReportError(loop,thrown);continueReading=true}}()
	if state==nil{return false};symbol:=tsGet(loop,tsInternalNetSymbolsModule(),tsNodeNetKey("kDestroyOnRead"));if tsTruthy(tsGet(loop,state.stream,symbol)){destroy:=tsNodeNetGet(loop,state.stream,"destroy");if destroy.kind==tsFunctionKind{tsCallReceiver(loop,destroy,state.stream)};return false}
	if readable:=tsNodeReadableStateFor(state.stream);readable!=nil{return tsTruthy(tsNodeReadablePush(loop,readable,tsNodeBuffer(data),tsU,false))};return true
}
func tsNodeNetDeliverOnRead(loop *tsLoop,state *tsNodeNetSocketState,data []byte)(continueReading bool){
	if state==nil||!state.onreadEnabled{return true};continueReading=true;offset:=0
	isPaused:=func()bool{method:=tsNodeNetGet(loop,state.stream,"isPaused");return method.kind==tsFunctionKind&&tsTruthy(tsCallReceiver(loop,method,state.stream))}
	refreshBuffer:=func(){if state.onreadBufferFactory.kind!=tsFunctionKind{return};next:=tsCall(loop,state.onreadBufferFactory);if next.kind==tsTypedArrayKind{array:=(*tsTypedArray)(next.ref);if array.name=="Uint8Array"||array.nodeBuffer{state.onreadBuffer=next}}}
	for offset<len(data){destination:=state.onreadBuffer
		if tsStrictEqual(destination,tsBooleanValue(true)){var result tsValue;func(){defer func(){if thrown:=recover();thrown!=nil{tsNodeNetReportError(loop,thrown)}}();result=tsCall(loop,state.onreadCallback,tsNumberValue(float64(len(data)-offset)),tsBooleanValue(true))}();func(){defer func(){if thrown:=recover();thrown!=nil{tsNodeNetReportError(loop,thrown)}}();refreshBuffer()}();if tsTruthy(tsNodeNetGet(loop,state.stream,"destroyed")){return false};if tsStrictEqual(result,tsBooleanValue(false))||isPaused(){state.onreadTail=nil;state.onreadTailSet=true;return false};return true}
		if destination.kind!=tsTypedArrayKind{return true};array:=(*tsTypedArray)(destination.ref);if array.name!="Uint8Array"&&!array.nodeBuffer{return true};target:=tsFSBytes(destination);if len(target)==0{failure:=tsFSException("Error","ENOBUFS","read ENOBUFS");tsNodeNetSet(loop,failure,"errno",tsNumberValue(105));tsNodeNetSet(loop,failure,"syscall",tsNodeNetKey("read"));destroy:=tsNodeNetGet(loop,state.stream,"destroy");if destroy.kind==tsFunctionKind{tsCallReceiver(loop,destroy,state.stream,failure)};return false}
		n:=len(target);if remaining:=len(data)-offset;remaining<n{n=remaining};copy(target[:n],data[offset:offset+n]);offset+=n;var result tsValue
		func(){defer func(){if thrown:=recover();thrown!=nil{tsNodeNetReportError(loop,thrown)}}();result=tsCall(loop,state.onreadCallback,tsNumberValue(float64(n)),destination)}()
		func(){defer func(){if thrown:=recover();thrown!=nil{tsNodeNetReportError(loop,thrown)}}();refreshBuffer()}()
		if tsTruthy(tsNodeNetGet(loop,state.stream,"destroyed")){return false};if tsStrictEqual(result,tsBooleanValue(false))||isPaused(){state.onreadTail=append([]byte(nil),data[offset:]...);state.onreadTailSet=true;return false}
	}
	return continueReading
}
func tsNodeNetReadFailure(loop *tsLoop,state *tsNodeNetSocketState,err error){
	if state==nil||state.closed{return}
	readError:=tsNodeTLSEstablishedError(loop,state.stream,err,"read")
	if !tsTruthy(tsNodeNetGet(loop,state.stream,"_hadError")){
		// The TLS wrapper listens for this private event to own teardown of its
		// upgraded socket. Go's TLS Conn combines the raw and encrypted handles,
		// so this is only needed by adapters that still install that hook.
		symbols:=tsInternalNetSymbolsModule();readErrorKey:=tsGet(loop,symbols,tsNodeNetKey("kReadError"));listenerCount:=tsNodeNetGet(loop,state.stream,"listenerCount")
		if readErrorKey.kind==tsSymbolKind&&listenerCount.kind==tsFunctionKind&&tsNumber(tsCallReceiver(loop,listenerCount,state.stream,readErrorKey))>0{
			emit:=tsNodeNetGet(loop,state.stream,"emit");if emit.kind==tsFunctionKind&&tsTruthy(tsCallReceiver(loop,emit,state.stream,readErrorKey,readError)){tsNodeNetSet(loop,state.stream,"_hadError",tsBooleanValue(true));destroy:=tsNodeNetGet(loop,state.stream,"destroy");if destroy.kind==tsFunctionKind{tsCallReceiver(loop,destroy,state.stream)};return}
		}
		if listenerCount.kind==tsFunctionKind&&tsNumber(tsCallReceiver(loop,listenerCount,state.stream,tsNodeNetKey("error")))>0{
			// Destroying with the read error lets the stream layer report it once
			// and mark close as an errored close, as SocketHandlers does.
			tsNodeNetSet(loop,state.stream,"_hadError",tsBooleanValue(true));destroy:=tsNodeNetGet(loop,state.stream,"destroy");if destroy.kind==tsFunctionKind{tsCallReceiver(loop,destroy,state.stream,readError)};return
		}
	}
	// With no error listener Node treats a closed read side like EOF. Preserve
	// buffered data and let allowHalfOpen decide whether writes stay available.
	if readable:=tsNodeReadableStateFor(state.stream);readable!=nil{tsNodeReadablePush(loop,readable,tsNull,tsU,false)}
	if !tsTruthy(tsNodeNetGet(loop,state.stream,"allowHalfOpen")){if writable:=tsNodeWritableStateFor(state.stream);writable!=nil{tsNodeWritableEnd(loop,writable,nil)}}
}
func tsNodeNetStartRead(state *tsNodeNetSocketState) {
	if state==nil||state.conn==nil||state.closed||state.reading||state.tlsHandshaking{return};loop,conn:=state.loop,state.conn
	if state.onreadEnabled&&state.onreadTailSet{tail:=append([]byte(nil),state.onreadTail...);state.onreadTail=nil;state.onreadTailSet=false;if len(tail)>0&&!tsNodeNetDeliverOnRead(loop,state,tail){return};if state.onreadTailSet{return};if state.onreadEOF{state.onreadEOF=false;tsNodeNetFinishReadEOF(loop,state);return}}
	state.reading=true
	var data []byte;var readErr error;var count int
	task:=loop.submit(func()tsResult{buffer:=make([]byte,32*1024);count,readErr=conn.Read(buffer);if count>0{data=append([]byte(nil),buffer[:count]...)};return tsResult{tsU,false}},func(_ tsResult){
		state.reading=false
		if state.upgrading{state.upgrading=false;if count>0{state.upgradeData=append(state.upgradeData,data...)};_ = state.conn.SetReadDeadline(time.Time{});resume:=state.upgradeContinue;state.upgradeContinue=nil;if resume!=nil{resume()};return}
		if state.closed{return};continueReading:=true;if count>0{state.bytesRead+=int64(count);tsNodeNetSyncSocket(state);tsNodeNetResetTimeout(state);if state.onreadEnabled{continueReading=tsNodeNetDeliverOnRead(loop,state,data)}else{continueReading=tsNodeNetDeliverData(loop,state,data)}}
		if readErr!=nil{if errors.Is(readErr,io.EOF){if state.onreadEnabled&&state.onreadTailSet{state.onreadEOF=true;return};tsNodeNetFinishReadEOF(loop,state);return};tsNodeNetReadFailure(loop,state,readErr);return};if continueReading{tsNodeNetStartRead(state)}
	});task.unref=state.unref
}
func tsNodeNetFinishReadEOF(loop *tsLoop,state *tsNodeNetSocketState){if readable:=tsNodeReadableStateFor(state.stream);readable!=nil{tsNodeReadablePush(loop,readable,tsNull,tsU,false)};if !tsTruthy(tsNodeNetGet(loop,state.stream,"allowHalfOpen")){if writable:=tsNodeWritableStateFor(state.stream);writable!=nil{tsNodeWritableEnd(loop,writable,nil)}}}
func tsNodeNetSocketWrite(loop *tsLoop, receiver tsValue, args ...tsValue) tsValue {
	state:=tsNodeNetSocketFor(receiver);callback:=tsArg(args,2);if state==nil||state.closed||(state.conn==nil&&!state.syncWrite){err:=tsFSException("Error","ERR_SOCKET_CLOSED","Socket is not connected");if callback.kind==tsFunctionKind{tsCall(loop,callback,err)};return tsU}
	chunk:=tsArg(args,0);encoding:=tsArg(args,1);if encoding.kind==tsFunctionKind{encoding=tsU};var data []byte;if chunk.kind==tsStringKind{buffer:=tsNodeBufferFrom(loop,[]tsValue{chunk,encoding});data=append([]byte(nil),tsNodeBufferVisibleBytes(buffer)...)}else{data=append([]byte(nil),tsFSBytes(chunk)...)}
	if state.connecting{state.pendingData=chunk;state.pendingEncoding=encoding;state.pendingWriteCallback=callback;tsNodeNetSet(loop,receiver,"_pendingData",chunk);resume:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{if tsStrictEqual(state.pendingWriteCallback,callback){state.pendingData=tsU;state.pendingEncoding=tsU;state.pendingWriteCallback=tsU;tsNodeNetSet(loop,receiver,"_pendingData",tsNull);tsNodeNetSocketWrite(loop,receiver,[]tsValue{chunk,encoding,callback}...)};return tsU}));closed:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{if tsStrictEqual(state.pendingWriteCallback,callback){state.pendingData=tsU;state.pendingEncoding=tsU;state.pendingWriteCallback=tsU;tsNodeNetSet(loop,receiver,"_pendingData",tsNull);if callback.kind==tsFunctionKind{tsCall(loop,callback,tsFSException("Error","ERR_SOCKET_CLOSED_BEFORE_CONNECTION","Socket closed before the connection was established"))}};return tsU}));tsCallReceiver(loop,tsNodeNetGet(loop,receiver,"once"),receiver,tsNodeNetKey("connect"),resume);tsCallReceiver(loop,tsNodeNetGet(loop,receiver,"once"),receiver,tsNodeNetKey("close"),closed);return tsU}
	if state.syncWrite&&state.conn==nil{tsNodeNetResetTimeout(state);offset:=0;for offset<len(data){written,err:=syscall.Write(state.syncWriteFD,data[offset:]);if err!=nil{if callback.kind==tsFunctionKind{tsCall(loop,callback,tsNodeNetError(err,"write"))};return tsU};if written==0{if callback.kind==tsFunctionKind{tsCall(loop,callback,tsFSException("Error","EIO","write returned no bytes"))};return tsU};offset+=written};state.bytesWritten+=int64(offset);tsNodeNetSyncSocket(state);if callback.kind==tsFunctionKind{tsCall(loop,callback,tsU)};return tsU}
	conn:=state.conn
	var writeErr error;task:=loop.submit(func()tsResult{_,writeErr=conn.Write(data);if writeErr!=nil{return tsResult{tsNodeNetError(writeErr,"write"),true}};return tsResult{tsU,false}},func(result tsResult){if result.rejected{writeError:=tsNodeTLSEstablishedError(loop,receiver,writeErr,"write");if callback.kind==tsFunctionKind{tsCall(loop,callback,writeError)};if !state.closed{destroy:=tsNodeNetGet(loop,receiver,"destroy");if destroy.kind==tsFunctionKind{tsCallReceiver(loop,destroy,receiver,writeError)}};return};state.bytesWritten+=int64(len(data));tsNodeNetSyncSocket(state);tsNodeNetResetTimeout(state);if callback.kind==tsFunctionKind{tsCall(loop,callback,tsU)}});task.unref=state.unref;return tsU
}
func tsNodeNetSocketWritev(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{entries:=tsArg(args,0);callback:=tsArg(args,1);if entries.kind!=tsArrayKind{return tsNodeNetSocketWrite(loop,receiver,[]tsValue{entries,tsNodeNetKey("buffer"),callback}...)};var data []byte;for _,entry:=range (*tsArray)(entries.ref).values{chunk,encoding:=tsNodeNetGet(loop,entry,"chunk"),tsNodeNetGet(loop,entry,"encoding");if chunk.kind==tsStringKind{data=append(data,tsNodeBufferVisibleBytes(tsNodeBufferFrom(loop,[]tsValue{chunk,encoding}))...)}else{data=append(data,tsFSBytes(chunk)...)} };return tsNodeNetSocketWrite(loop,receiver,[]tsValue{tsNodeBuffer(data),tsNodeNetKey("buffer"),callback}...)}
func tsNodeNetSocketFinal(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue {
	state:=tsNodeNetSocketFor(receiver);callback:=tsArg(args,0);if state!=nil&&state.connecting{resume:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{return tsNodeNetSocketFinal(loop,receiver,[]tsValue{callback}...)}));tsCallReceiver(loop,tsNodeNetGet(loop,receiver,"once"),receiver,tsNodeNetKey("connect"),resume);return tsU};if state!=nil&&state.conn!=nil{if closer,ok:=state.conn.(interface{CloseWrite()error});ok{if err:=closer.CloseWrite();err!=nil&&callback.kind==tsFunctionKind{tsCall(loop,callback,tsNodeNetError(err,"shutdown"));return tsU}}};if callback.kind==tsFunctionKind{tsCall(loop,callback,tsU)};return tsU
}
// Bun's closeAdoptedTLSRawNT defers closing a shared raw handle. Here a Go TLS
// connection owns one net.Conn; _destroy closes that connection and the local
// Duplex destroy path controls close-event delivery. There is no second handle.
func tsNodeNetSocketDestroy(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue {
	state:=tsNodeNetSocketFor(receiver);err,callback:=tsArg(args,0),tsArg(args,1);if state!=nil&&!state.closed{state.closed=true;if state.connectCancel!=nil{state.connectCancel();state.connectCancel=nil};tsNodeNetTraceEnd(loop,state);tsNodeNetTraceEnd(loop,state.connectTrace);state.connectTrace=nil;state.connecting=false;if state.timeoutTask!=nil{loop.clearTimeout(tsTaskValue(state.timeoutTask));state.timeoutTask=nil};if state.conn!=nil{_ = state.conn.Close()};if state.syncWrite&&state.syncWriteAutoClose{_ = syscall.Close(state.syncWriteFD);state.syncWrite=false};tsNodeNetSet(loop,receiver,"_handle",tsNull);tsNodeNetSyncSocket(state);if state.server!=nil{state.server.connections--;if state.server.connections<0{state.server.connections=0};tsNodeNetSet(state.server.loop,state.server.server,"_connections",tsNumberValue(float64(state.server.connections)));tsNodeNetServerMaybeClose(state.server)}};if callback.kind==tsFunctionKind{tsCall(loop,callback,err)};return tsU
}
func tsNodeNetSocketAbort(loop *tsLoop,socket tsValue,state *tsNodeNetSocketState,signal *tsAbortState){
	if state==nil||signal==nil{return}
	_,aborted:=tsAbortReason(signal)
	if aborted{loop.ready=append(loop.ready,func(){destroy:=tsNodeNetGet(loop,socket,"destroy");if destroy.kind==tsFunctionKind{tsCallReceiver(loop,destroy,socket,tsFSAbortError(signal))}});return}
	signalValue:=tsObjectValue(signal.events.object)
	onAbort:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{if !state.closed{destroy:=tsNodeNetGet(loop,socket,"destroy");if destroy.kind==tsFunctionKind{tsCallReceiver(loop,destroy,socket,tsFSAbortError(signal))}};return tsU}))
	state.abortDispose=tsInternalAbortListenerAdd(loop,signalValue,onAbort)
	onClose:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{dispose:=tsGet(loop,state.abortDispose,tsPropertyKey(tsWellKnownSymbol("dispose")));if dispose.kind==tsFunctionKind{tsCallReceiver(loop,dispose,state.abortDispose)};state.abortDispose=tsU;return tsU}))
	tsCallReceiver(loop,tsNodeNetGet(loop,socket,"once"),socket,tsNodeNetKey("close"),onClose)
}
func tsNodeNetSocketEnd(loop *tsLoop,socket tsValue){
	if tsTruthy(tsNodeNetGet(loop,socket,"allowHalfOpen")){return}
	originalWrite:=tsNodeNetGet(loop,socket,"write")
	writeAfterFIN:=tsNamedNativeMethod("write",3,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{
		if !tsTruthy(tsNodeNetGet(loop,receiver,"writableEnded")){return tsCallReceiver(loop,originalWrite,receiver,args...)}
		encoding,callback:=tsArg(args,1),tsArg(args,2);if encoding.kind==tsFunctionKind{callback=encoding}
		err:=tsFSException("Error","EPIPE","This socket has been ended by the other party")
		if callback.kind==tsFunctionKind{loop.ready=append(loop.ready,func(){tsCall(loop,callback,err)})}
		destroy:=tsNodeNetGet(loop,receiver,"destroy");if destroy.kind==tsFunctionKind{tsCallReceiver(loop,destroy,receiver,err)}
		return tsBooleanValue(false)
	})
	tsSet(loop,socket,tsNodeNetKey("write"),writeAfterFIN)
}
func tsNodeNetDialSocket(loop *tsLoop,receiver tsValue,state *tsNodeNetSocketState,network,address string,local tsValue,localPort int){
	var dialed net.Conn;var dialErr error
	if network=="tcp"&&!state.connectRace{remote,_,_:=net.SplitHostPort(address);family:=6;if ip:=net.ParseIP(strings.Split(remote,"%")[0]);ip!=nil&&ip.To4()!=nil{family=4};tsNodeNetEmitSafe(loop,receiver,"connectionAttempt",tsNodeNetKey(remote),tsNodeNetGet(loop,receiver,"_port"),tsNumberValue(float64(family)));if !state.connecting{return};if blockList:=tsNodeNetGet(loop,receiver,"blockList");tsTruthy(blockList){check:=tsNodeNetGet(loop,blockList,"check");if check.kind==tsFunctionKind&&tsTruthy(tsCallReceiver(loop,check,blockList,tsNodeNetKey(remote),tsNodeNetKey(fmt.Sprintf("ipv%d",family)))){tsNodeNetConnectionFailure(loop,receiver,state,tsFSException("Error","ERR_IP_BLOCKED","IP address is blocked: "+remote));return}}}
	if !state.connectRace{pipePath:="";if network=="unix"{pipePath=address};tsNodeNetTraceStart(loop,state,pipePath);if network=="tcp"{remote,_,_:=net.SplitHostPort(address);tsNodeNetPerfStart(loop,receiver,remote,int(tsNumber(tsNodeNetGet(loop,receiver,"_port"))))}}
	ctx,cancel:=context.WithCancel(context.Background());state.connectCancel=cancel
	task:=loop.submit(func()tsResult{
		dialer:=net.Dialer{};if network=="tcp"&&(!tsNullish(local)||localPort!=0){localIP:=net.IP(nil);if !tsNullish(local){localIP=net.ParseIP(tsText(local))}else if remote,_,err:=net.SplitHostPort(address);err==nil{if ip:=net.ParseIP(remote);ip!=nil&&ip.To4()==nil{localIP=net.ParseIP("::")}else{localIP=net.IPv4zero}};dialer.LocalAddr=&net.TCPAddr{IP:localIP,Port:localPort}};if state.connectRace&&state.connectAttemptTimeout>0&&state.connectAttemptIndex<len(state.connectAddresses){dialer.Timeout=state.connectAttemptTimeout};conn,err:=dialer.DialContext(ctx,network,address);if err!=nil{dialErr=err;return tsResult{tsNodeNetError(err,"connect"),true}};dialed=conn;return tsResult{tsU,false}
	},func(result tsResult){
		if state.connectCancel==cancel{state.connectCancel=nil};cancel()
		timedOut:=false;if state.connectRace&&result.rejected{if netErr,ok:=dialErr.(net.Error);ok{timedOut=netErr.Timeout()}}
		if !timedOut{if state.connectRace{tsNodeNetTraceEnd(loop,state.connectTrace);state.connectTrace=nil}else{tsNodeNetTraceEnd(loop,state)}}
		if state.closed||!state.connecting{if timedOut&&state.connectTrace!=nil{tsNodeNetTraceEnd(loop,state.connectTrace);state.connectTrace=nil};if dialed!=nil{_ = dialed.Close()};return}
		if state.connectRace&&result.rejected{
			addressIndex:=state.connectAttemptIndex-1;attempt:=state.connectAddresses[addressIndex];port:=tsNodeNetGet(loop,receiver,"_port");result.value=tsNodeNetConnectError(loop,result.value,attempt.ip,int(tsNumber(port)),local,localPort);if timedOut{tsNodeNetEmitSafe(loop,receiver,"connectionAttemptTimeout",tsNodeNetKey(attempt.ip),port,tsNumberValue(float64(attempt.family)));if state.connectTrace!=nil{tsNodeNetTraceEnd(loop,state.connectTrace);state.connectTrace=nil}}else{tsNodeNetEmitSafe(loop,receiver,"connectionAttemptFailed",tsNodeNetKey(attempt.ip),port,tsNumberValue(float64(attempt.family)),result.value)};if !state.connecting{return};state.connectErrors=append(state.connectErrors,result.value);if state.connectAttemptIndex<len(state.connectAddresses){tsNodeNetStartNextConnectAttempt(loop,receiver,state,network,local,localPort);return};failure:=tsNodeNetConnectAggregate(loop,state);tsNodeNetConnectionFailure(loop,receiver,state,failure);return
		}
		if result.rejected{if network=="tcp"{remote,_,_:=net.SplitHostPort(address);family:=6;if ip:=net.ParseIP(strings.Split(remote,"%")[0]);ip!=nil&&ip.To4()!=nil{family=4};port:=int(tsNumber(tsNodeNetGet(loop,receiver,"_port")));result.value=tsNodeNetConnectError(loop,result.value,remote,port,local,localPort);tsNodeNetEmitSafe(loop,receiver,"connectionAttemptFailed",tsNodeNetKey(remote),tsNodeNetGet(loop,receiver,"_port"),tsNumberValue(float64(family)),result.value);if !state.connecting{return}};tsNodeNetConnectionFailure(loop,receiver,state,result.value);return}
		if state.connectRace{state.connectRace=false;state.connectAddresses=nil;state.connectErrors=nil;state.connectAttemptIndex=0}
		tsNodeNetSetConn(state,dialed);tsNodeNetStartRead(state);tlsClient:=tsTruthy(tsNodeNetGet(loop,receiver,"secureConnecting"));tsNodeNetEmitSafe(loop,receiver,"connect");if !tlsClient{tsNodeNetEmitSafe(loop,receiver,"ready")};tsNodeNetPerfStop(loop,receiver)
	});task.unref=state.unref
}
func tsNodeNetConnectionFailure(loop *tsLoop,receiver tsValue,state *tsNodeNetSocketState,failure tsValue){if state.connectCancel!=nil{state.connectCancel();state.connectCancel=nil};state.connecting=false;state.connectRace=false;state.connectAddresses=nil;state.connectErrors=nil;tsNodeNetSyncSocket(state);if tsTruthy(tsNodeNetGet(loop,receiver,"secureConnecting")){tsNodeNetSet(loop,receiver,"secureConnecting",tsBooleanValue(false));tsNodeNetSet(loop,receiver,"_hadError",tsBooleanValue(true));tsNodeNetSet(loop,receiver,"_closeAfterHandlingError",tsBooleanValue(true))};destroy:=tsNodeNetGet(loop,receiver,"destroy");if destroy.kind==tsFunctionKind{tsCallReceiver(loop,destroy,receiver,failure)}}
func tsNodeNetConnectAggregate(loop *tsLoop,state *tsNodeNetSocketState)tsValue{failure:=tsFSException("AggregateError","ERR_SOCKET_CONNECTION_FAILED","All connection attempts failed");errorsArray:=&tsArray{values:append([]tsValue(nil),state.connectErrors...)};tsNodeNetSet(loop,failure,"errors",tsArrayValue(errorsArray));return failure}
func tsNodeNetStartNextConnectAttempt(loop *tsLoop,receiver tsValue,state *tsNodeNetSocketState,network string,local tsValue,localPort int){
	if state==nil||!state.connecting{return}
	index:=state.connectAttemptIndex;if index>=len(state.connectAddresses){return};address:=state.connectAddresses[index];state.connectAttemptIndex++;port:=tsNodeNetGet(loop,receiver,"_port")
	if blockList:=tsNodeNetGet(loop,receiver,"blockList");tsTruthy(blockList){check:=tsNodeNetGet(loop,blockList,"check");if check.kind==tsFunctionKind&&tsTruthy(tsCallReceiver(loop,check,blockList,tsNodeNetKey(address.ip),tsNodeNetKey(fmt.Sprintf("ipv%d",address.family)))){failure:=tsFSException("Error","ERR_IP_BLOCKED","IP address is blocked: "+address.ip);tsNodeNetEmitSafe(loop,receiver,"connectionAttemptFailed",tsNodeNetKey(address.ip),port,tsNumberValue(float64(address.family)),failure);state.connectErrors=append(state.connectErrors,failure);if !state.connecting{return};if state.connectAttemptIndex<len(state.connectAddresses){tsNodeNetStartNextConnectAttempt(loop,receiver,state,network,local,localPort)}else{tsNodeNetConnectionFailure(loop,receiver,state,tsNodeNetConnectAggregate(loop,state))};return}}
	tsNodeNetEmitSafe(loop,receiver,"connectionAttempt",tsNodeNetKey(address.ip),port,tsNumberValue(float64(address.family)))
	attempted:=tsNodeNetGet(loop,receiver,"autoSelectFamilyAttemptedAddresses");if attempted.kind!=tsArrayKind{array:=&tsArray{};tsNodeNetSet(loop,receiver,"autoSelectFamilyAttemptedAddresses",tsArrayValue(array));attempted=tsArrayValue(array)};(*tsArray)(attempted.ref).values=append((*tsArray)(attempted.ref).values,tsNodeNetKey(fmt.Sprintf("%s:%s",address.ip,tsText(port))))
	if !state.connecting{return};state.connectTrace=&tsNodeNetSocketState{};tsNodeNetTraceStart(loop,state.connectTrace,"");tsNodeNetPerfStart(loop,receiver,address.ip,int(tsNumber(port)));addressValue:=net.JoinHostPort(address.ip,tsText(port));localValue:=local;if tsNullish(localValue)&&localPort!=0{if address.family==6{localValue=tsNodeNetKey("::")}else{localValue=tsNodeNetKey("0.0.0.0")}};tsNodeNetDialSocket(loop,receiver,state,network,addressValue,localValue,localPort)
}
func tsNodeNetSocketConnect(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue {
	state:=tsNodeNetSocketFor(receiver);if state==nil{return receiver};if (state.connected&&!state.closed)||state.connecting{panic(tsThrown{tsFSException("Error","ERR_SOCKET_DGRAM_IS_CONNECTED","Socket is already connecting")})};if state.closed{undestroy:=tsNodeNetGet(loop,receiver,"_undestroy");if undestroy.kind==tsFunctionKind{tsCallReceiver(loop,undestroy,receiver)};state.closed=false;state.connected=false;state.connecting=false;state.conn=nil;state.localAddress="";state.localFamily="";state.localPort=0;state.remoteAddress="";state.remoteFamily="";state.remotePort=0;state.pendingData=tsU;state.pendingEncoding=tsU;state.pendingWriteCallback=tsU;tsNodeNetSet(loop,receiver,"_sockname",tsNull)}
	arg:=tsArg(args,0);callback:=tsArg(args,1);options:=arg
	if tsNodeNetIsNormalizedArgs(loop,arg){items:=(*tsArray)(arg.ref);options=tsArg(items.values,0);if tsIsUndefined(callback){callback=tsArg(items.values,1)}}else if arg.kind!=tsObjectKind&&arg.kind!=tsInstanceKind&&!tsNullish(arg){normalized:=tsNodeNetNormalizeArgs(args);items:=(*tsArray)(normalized.ref);options=tsArg(items.values,0);callback=tsArg(items.values,1)}
	path:=tsNodeNetGet(loop,options,"path");portValue:=tsNodeNetGet(loop,options,"port");hostValue:=tsNodeNetGet(loop,options,"host");host:="localhost";if !tsIsUndefined(hostValue)&&!tsIsNull(hostValue){if hostValue.kind!=tsStringKind{tsFSInvalid("options.host","a string")};if tsText(hostValue)!=""{host=tsText(hostValue)}}
	network:="tcp";if !tsNullish(path)&&tsText(path)!=""{if path.kind!=tsStringKind{tsFSInvalid("options.path","a string")};network="unix";host=tsText(path)}else if tsIsUndefined(portValue){panic(tsThrown{tsFSException("TypeError","ERR_MISSING_ARGS","The \"options\" or \"port\" or \"path\" argument must be specified")})}
	port:=0;if !tsNullish(portValue){if portValue.kind!=tsStringKind&&!tsIsNumeric(portValue){tsFSInvalid("options.port","a number or string")};number:=portValue.ToNumber();if math.IsNaN(number)||math.IsInf(number,0)||math.Trunc(number)!=number||number<0||number>65535{panic(tsThrown{tsFSException("RangeError","ERR_SOCKET_BAD_PORT","Port should be >= 0 and < 65536. Received "+tsText(portValue))})};port=int(number)}
	local:=tsNodeNetGet(loop,options,"localAddress");if !tsNullish(local)&&(!tsInternalNetIPv4(local)&&!tsInternalNetIPv6(local)){panic(tsThrown{tsFSException("TypeError","ERR_INVALID_IP_ADDRESS","Invalid IP address: "+tsText(local))})};localPortValue:=tsNodeNetGet(loop,options,"localPort");localPort:=0;if !tsNullish(localPortValue){if !tsIsNumeric(localPortValue){tsFSInvalid("options.localPort","a number")};localPort=int(tsNumber(localPortValue));if localPort<0||localPort>65535{tsFSInvalid("options.localPort","an integer between 0 and 65535")}}
	family:=tsNodeNetGet(loop,options,"family");if family.kind==tsStringKind{switch tsText(family){case "IPv4":family=tsNumberValue(4);case "IPv6":family=tsNumberValue(6)}};autoSelect:=tsNodeNetSingleton.defaultAutoSelectFamily;autoValue:=tsNodeNetGet(loop,options,"autoSelectFamily");if !tsNullish(autoValue){if autoValue.kind!=tsBooleanKind{tsFSInvalid("options.autoSelectFamily","a boolean")};autoSelect=tsTruthy(autoValue)}
	attemptTimeout:=tsNodeNetSingleton.defaultAutoSelectFamilyAttemptTimeout;attemptTimeoutValue:=tsNodeNetGet(loop,options,"autoSelectFamilyAttemptTimeout");if !tsNullish(attemptTimeoutValue){if !tsIsNumeric(attemptTimeoutValue)||tsNumber(attemptTimeoutValue)!=math.Trunc(tsNumber(attemptTimeoutValue))||tsNumber(attemptTimeoutValue)<1||tsNumber(attemptTimeoutValue)>math.MaxInt32{tsFSInvalid("options.autoSelectFamilyAttemptTimeout","an integer greater than or equal to 1")};attemptTimeout=int(tsNumber(attemptTimeoutValue));if attemptTimeout<10{attemptTimeout=10}}
	lookup:=tsNodeNetGet(loop,options,"lookup");if !tsNullish(lookup)&&lookup.kind!=tsFunctionKind{tsFSInvalid("options.lookup","a function")};if lookup.kind!=tsFunctionKind{lookup=tsNodeNetGet(loop,tsNodeDNSModule(loop,false),"lookup")}
	if callback.kind==tsFunctionKind{event:="connect";if tsTruthy(tsNodeNetGet(loop,receiver,"secureConnecting")){event="secureConnect"};tsCallReceiver(loop,tsNodeNetGet(loop,receiver,"once"),receiver,tsNodeNetKey(event),callback)}
	state.connecting=true;state.connected=false;tsNodeNetSyncSocket(state);tsNodeNetSet(loop,receiver,"_host",tsNodeNetKey(host));tsNodeNetSet(loop,receiver,"_port",tsNumberValue(float64(port)))
	if fd:=tsNodeNetGet(loop,options,"fd");!tsNullish(fd)&&!tsIsUndefined(fd){if state.syncWrite&&int(tsNumber(fd))==state.syncWriteFD&&!tsStrictEqual(tsNodeNetGet(loop,options,"readable"),tsBooleanValue(true)){state.connecting=false;state.connected=true;tsNodeNetSyncSocket(state);loop.ready=append(loop.ready,func(){tsNodeNetEmitSafe(loop,receiver,"connect");tsNodeNetEmitSafe(loop,receiver,"ready")});return receiver};state.connecting=false;tsNodeNetSyncSocket(state);failure:=tsFSException("TypeError","ERR_INVALID_FD_TYPE","Descriptor-backed network reads are not available for this fd");tsNodeNetConnectionFailure(loop,receiver,state,failure);return receiver}
	if network=="unix"{tsNodeNetDialSocket(loop,receiver,state,network,host,local,localPort);return receiver}
	if ip:=net.ParseIP(strings.Split(host,"%")[0]);ip!=nil{tsNodeNetDialSocket(loop,receiver,state,network,net.JoinHostPort(host,strconv.Itoa(port)),local,localPort);return receiver}
	lookupOptions:=tsNewObject();if !tsIsUndefined(family){lookupOptions.set("family",family)};hints:=tsNodeNetGet(loop,options,"hints");if !tsTruthy(hints){hints=tsNumberValue(0)};lookupOptions.set("hints",hints);familyConstrained:=tsIsNumeric(family)&&(tsNumber(family)==4||tsNumber(family)==6);multi:=autoSelect&&!familyConstrained&&tsNullish(local);if multi{lookupOptions.set("all",tsBooleanValue(true))}
	resolved:=tsFunctionValue(tsFunc(func(loop *tsLoop,values ...tsValue)tsValue{
		err,addressValue,familyValue:=tsArg(values,0),tsArg(values,1),tsArg(values,2);if !tsNullish(err){tsNodeNetEmitSafe(loop,receiver,"lookup",err,tsU,tsU,tsNodeNetKey(host));tsNodeNetConnectionFailure(loop,receiver,state,err);return tsU}
		if multi{if addressValue.kind!=tsArrayKind{bad:=tsFSException("Error","ERR_INVALID_IP_ADDRESS","DNS lookup did not return an address list");tsNodeNetConnectionFailure(loop,receiver,state,bad);return tsU};valid:=make([][]tsNodeNetConnectAddress,2);seen:=[]map[string]bool{{},{}};firstFamily:=0;for _,entry:=range (*tsArray)(addressValue.ref).values{ipValue,ipFamily:=tsNodeNetGet(loop,entry,"address"),int(tsNumber(tsNodeNetGet(loop,entry,"family")));ip:=tsText(ipValue);tsNodeNetEmitSafe(loop,receiver,"lookup",tsNull,ipValue,tsNumberValue(float64(ipFamily)),tsNodeNetKey(host));if !state.connecting{return tsU};if !tsInternalNetIPv4(ipValue)&&!tsInternalNetIPv6(ipValue){continue};if ipFamily!=4&&ipFamily!=6{continue};if firstFamily==0{firstFamily=ipFamily};slot:=1;if ipFamily==firstFamily{slot=0};if seen[slot]==nil{seen[slot]=map[string]bool{}};if !seen[slot][ip]{seen[slot][ip]=true;valid[slot]=append(valid[slot],tsNodeNetConnectAddress{ip:ip,family:ipFamily})}}
			if len(valid[0])==0&&len(valid[1])==0{failure:=tsFSException("TypeError","ERR_INVALID_IP_ADDRESS","No valid IP addresses were returned for "+host);tsNodeNetConnectionFailure(loop,receiver,state,failure);return tsU};addresses:=make([]tsNodeNetConnectAddress,0,len(valid[0])+len(valid[1]));for index:=0;index<len(valid[0])||index<len(valid[1]);index++{if index<len(valid[0]){addresses=append(addresses,valid[0][index])};if index<len(valid[1]){addresses=append(addresses,valid[1][index])}};if len(addresses)==1{tsNodeNetDialSocket(loop,receiver,state,network,net.JoinHostPort(addresses[0].ip,strconv.Itoa(port)),local,localPort);return tsU};state.connectRace=true;state.connectAddresses=addresses;state.connectAttemptIndex=0;state.connectAttemptTimeout=time.Duration(attemptTimeout)*time.Millisecond;state.connectErrors=nil;tsNodeNetSet(loop,receiver,"autoSelectFamilyAttemptedAddresses",tsArrayValue(&tsArray{}));tsNodeNetStartNextConnectAttempt(loop,receiver,state,network,local,localPort);return tsU}
		if addressValue.kind!=tsStringKind||(!tsInternalNetIPv4(addressValue)&&!tsInternalNetIPv6(addressValue)){failure:=tsFSException("TypeError","ERR_INVALID_IP_ADDRESS","Invalid IP address");tsNodeNetEmitSafe(loop,receiver,"lookup",tsNull,addressValue,familyValue,tsNodeNetKey(host));tsNodeNetConnectionFailure(loop,receiver,state,failure);return tsU};if tsNumber(familyValue)!=4&&tsNumber(familyValue)!=6{failure:=tsFSException("Error","ERR_INVALID_ADDRESS_FAMILY","Invalid address family");tsNodeNetEmitSafe(loop,receiver,"lookup",tsNull,addressValue,familyValue,tsNodeNetKey(host));tsNodeNetConnectionFailure(loop,receiver,state,failure);return tsU}
		tsNodeNetEmitSafe(loop,receiver,"lookup",tsNull,addressValue,familyValue,tsNodeNetKey(host));if !state.connecting{return tsU};tsNodeNetDialSocket(loop,receiver,state,network,net.JoinHostPort(tsText(addressValue),strconv.Itoa(port)),local,localPort);return tsU
	}))
	tsCall(loop,lookup,tsNodeNetKey(host),tsObjectValue(lookupOptions),resolved);return receiver
}
func tsNodeNetIsUint8Array(value tsValue)bool{if value.kind!=tsTypedArrayKind{return false};array:=(*tsTypedArray)(value.ref);return array.name=="Uint8Array"||array.nodeBuffer}
func tsNodeNetConfigureOnRead(loop *tsLoop,state *tsNodeNetSocketState,options tsValue){
	onread:=tsNodeNetGet(loop,options,"onread");if !tsTruthy(onread){return};if onread.kind!=tsObjectKind&&onread.kind!=tsInstanceKind{panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"TypeError",message:"onread must be an object"})})};callback:=tsNodeNetGet(loop,onread,"callback");if callback.kind!=tsFunctionKind{panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"TypeError",message:"onread.callback must be a function"})})};buffer:=tsNodeNetGet(loop,onread,"buffer");if !tsNodeNetIsUint8Array(buffer)&&buffer.kind!=tsFunctionKind{return};state.onreadEnabled=true;state.onreadCallback=callback
	if buffer.kind==tsFunctionKind{state.onreadBufferFactory=buffer;first:=tsCall(loop,buffer);if tsNodeNetIsUint8Array(first){state.onreadBuffer=first}else{state.onreadBuffer=tsBooleanValue(true)}}else{state.onreadBuffer=buffer}
}
func tsNodeNetSocketCreate(loop *tsLoop,options tsValue)tsValue {
	tsNodeDuplexInit(loop);if tsNullish(options){options=tsObjectValue(tsNewObject())}
	for _,name:=range []string{"objectMode","readableObjectMode","writableObjectMode"}{if tsTruthy(tsNodeNetGet(loop,options,name)){panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_VALUE","The property 'options."+name+"' is not supported")})}}
	if delay:=tsNodeNetGet(loop,options,"keepAliveInitialDelay");!tsIsUndefined(delay){if !tsIsNumeric(delay){tsFSInvalid("options.keepAliveInitialDelay","a number")}}
	if fd:=tsNodeNetGet(loop,options,"fd");!tsIsUndefined(fd)&&(!tsIsNumeric(fd)||tsNumber(fd)!=math.Trunc(tsNumber(fd))||tsNumber(fd)<0||tsNumber(fd)>math.MaxInt32){tsFSInvalid("options.fd","an integer between 0 and 2147483647")}
	syncWriteFD:=0;syncWrite,syncWriteAutoClose:=false,false
	if tsTruthy(tsNodeNetGet(loop,options,"writable")){fd:=tsNodeNetGet(loop,options,"fd");if !tsIsUndefined(fd){descriptor:=int(tsNumber(fd));var stat syscall.Stat_t;if err:=syscall.Fstat(descriptor,&stat);err!=nil{panic(tsThrown{tsFSException("TypeError","ERR_INVALID_FD_TYPE","Unsupported fd type: UNKNOWN")})};kind:=stat.Mode&syscall.S_IFMT;if kind==syscall.S_IFREG{panic(tsThrown{tsFSException("TypeError","ERR_INVALID_FD_TYPE","Unsupported fd type: FILE")})};readable:=tsStrictEqual(tsNodeNetGet(loop,options,"readable"),tsBooleanValue(true));if kind==syscall.S_IFIFO||kind==syscall.S_IFCHR||kind==syscall.S_IFSOCK&&!readable{syncWrite=true;syncWriteFD=descriptor;syncWriteAutoClose=descriptor>2}}}
	blockList:=tsNodeNetGet(loop,options,"blockList");if tsTruthy(blockList){isBlockList:=tsNodeNetGet(loop,tsNodeNetSingleton.blockListCtor,"isBlockList");if isBlockList.kind!=tsFunctionKind||!tsTruthy(tsCallReceiver(loop,isBlockList,tsNodeNetSingleton.blockListCtor,blockList)){tsFSInvalid("options.blockList","a net.BlockList")}}
	stream:=tsNodeDuplexCreate(loop,tsU,options);object:=(*tsObject)(stream.ref);object.prototype=tsObjectValue(tsNodeNetSingleton.socketProto);if tsTruthy(blockList){tsNodeNetSet(loop,stream,"blockList",blockList)}
	state:=&tsNodeNetSocketState{loop:loop,stream:stream,connecting:false,unref:false,remoteFamily:"",noDelay:tsTruthy(tsNodeNetGet(loop,options,"noDelay")),keepAlive:tsTruthy(tsNodeNetGet(loop,options,"keepAlive")),typeOfService:0,syncWriteFD:syncWriteFD,syncWrite:syncWrite,syncWriteAutoClose:syncWriteAutoClose};if delay:=tsNodeNetGet(loop,options,"keepAliveInitialDelay");tsIsNumeric(delay){milliseconds:=int32(tsNumber(delay));if milliseconds<0{milliseconds=0};state.keepAliveDelay=time.Duration(milliseconds)*time.Millisecond};tsNodeNetSingleton.sockets[uintptr(stream.ref)]=state;tsNodeEventCreate(loop,stream,tsU);tsNodeNetSyncSocket(state);tsNodeNetConfigureOnRead(loop,state,options)
	onEnd:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{tsNodeNetSocketEnd(loop,stream);return tsU}));tsCallReceiver(loop,tsNodeNetGet(loop,stream,"once"),stream,tsNodeNetKey("end"),onEnd)
	if syncWrite&&tsIsUndefined(tsNodeNetGet(loop,options,"readable")){tsNodeNetFinishReadEOF(loop,state)}
	if timeout:=tsNodeNetGet(loop,options,"timeout");!tsNullish(timeout)&&tsNumber(timeout)>0{state.timeout=time.Duration(tsNumber(timeout))*time.Millisecond}
	if tsTruthy(tsFSProperty(options,"signal")){if signal:=tsFSAbortSignal(options);signal!=nil{tsNodeNetSocketAbort(loop,stream,state,signal)}}
	return stream
}
func tsNodeNetSocketCtor(loop *tsLoop)tsValue {
	constructor:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeNetSocketCreate(loop,tsArg(args,0))});constructor.name="Socket";constructor.length=1;constructor.constructible=true;constructor.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeNetSocketCreate(loop,tsArg(args,0))};constructor.receiverCall=func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeNetSocketCreate(loop,tsArg(args,0))};constructor.properties=tsNewObject();constructor.properties.set("prototype",tsObjectValue(tsNodeNetSingleton.socketProto));constructor.properties.descriptors=map[string]*tsDescriptor{"prototype":{writable:false}};return tsFunctionValue(constructor)
}
func tsNodeNetSocketInstall(loop *tsLoop) {
	p:=tsNodeNetSingleton.socketProto
	methods:=map[string]tsValue{
		"connect":tsNamedNativeMethod("connect",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeNetSocketConnect(loop,receiver,args...)}),
		"_read":tsNamedNativeMethod("_read",1,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{tsNodeNetStartRead(tsNodeNetSocketFor(receiver));return tsU}),
		"_write":tsNamedNativeMethod("_write",3,tsNodeNetSocketWrite),
		"_writev":tsNamedNativeMethod("_writev",2,tsNodeNetSocketWritev),
		"_final":tsNamedNativeMethod("_final",1,tsNodeNetSocketFinal),
		"_destroy":tsNamedNativeMethod("_destroy",2,tsNodeNetSocketDestroy),
		"address":tsNamedNativeMethod("address",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{out:=tsNewObject();out.set("address",tsNodeNetGet(loop,receiver,"localAddress"));out.set("family",tsNodeNetGet(loop,receiver,"localFamily"));out.set("port",tsNodeNetGet(loop,receiver,"localPort"));return tsObjectValue(out)}),
		"setNoDelay":tsNamedNativeMethod("setNoDelay",1,func(_ *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeNetSocketFor(receiver);if state!=nil{state.noDelay=!tsStrictEqual(tsArg(args,0),tsBooleanValue(false));if tcp,ok:=tsNodeNetUnderlyingConn(state.conn).(*net.TCPConn);ok{_ = tcp.SetNoDelay(state.noDelay)}};return receiver}),
		"setKeepAlive":tsNamedNativeMethod("setKeepAlive",2,func(_ *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeNetSocketFor(receiver);if state!=nil{state.keepAlive=tsTruthy(tsArg(args,0));delay:=tsArg(args,1);milliseconds:=int32(0);if !tsIsUndefined(delay){milliseconds=int32(tsNumber(delay))};if milliseconds<0{milliseconds=0};state.keepAliveDelay=time.Duration(milliseconds)*time.Millisecond;if tcp,ok:=tsNodeNetUnderlyingConn(state.conn).(*net.TCPConn);ok{_ = tcp.SetKeepAlive(state.keepAlive);if state.keepAlive&&state.keepAliveDelay>0{_ = tcp.SetKeepAlivePeriod(state.keepAliveDelay)}}};return receiver}),
		"setTypeOfService":tsNamedNativeMethod("setTypeOfService",1,func(_ *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeNetSocketFor(receiver);value:=tsArg(args,0);if !tsIsNumeric(value)||tsNumber(value)!=math.Trunc(tsNumber(value))||tsNumber(value)<0||tsNumber(value)>255{tsFSInvalid("tos","an integer between 0 and 255")};if state!=nil{tos:=int(tsNumber(value));if !state.typeOfServiceSet||tos!=state.typeOfService{state.typeOfService=tos;state.typeOfServiceSet=true;if state.conn!=nil{if err:=tsNodeNetSetSocketTypeOfService(state.conn,tos);err!=nil&&runtime.GOOS!="windows"{panic(tsThrown{tsNodeNetError(err,"setTypeOfService")})}}}};return receiver}),
		"getTypeOfService":tsNamedNativeMethod("getTypeOfService",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeNetSocketFor(receiver);if state==nil||state.conn==nil{if state!=nil&&state.typeOfServiceSet{return tsNumberValue(float64(state.typeOfService))};return tsNumberValue(0)};value,err:=tsNodeNetGetSocketTypeOfService(state.conn);if err==nil{return tsNumberValue(float64(value))};if runtime.GOOS=="windows"{return tsNumberValue(float64(state.typeOfService))};panic(tsThrown{tsNodeNetError(err,"getTypeOfService")})}),
		"setTimeout":tsNamedNativeMethod("setTimeout",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeNetSocketFor(receiver);if state!=nil&&!state.closed{timeout:=tsArg(args,0);if !tsNullish(timeout){if !tsIsNumeric(timeout)||tsNumber(timeout)<0{tsFSInvalid("timeout","a non-negative number")};state.timeout=time.Duration(tsNumber(timeout))*time.Millisecond;tsNodeNetSet(loop,receiver,"timeout",timeout);tsNodeNetResetTimeout(state)};if listener:=tsArg(args,1);!tsIsUndefined(listener){if listener.kind!=tsFunctionKind{tsFSInvalid("callback","a function")};if tsIsNumeric(timeout)&&tsNumber(timeout)==0{tsCallReceiver(loop,tsNodeNetGet(loop,receiver,"removeListener"),receiver,tsNodeNetKey("timeout"),listener)}else{tsCallReceiver(loop,tsNodeNetGet(loop,receiver,"once"),receiver,tsNodeNetKey("timeout"),listener)}}};return receiver}),
		"ref":tsNamedNativeMethod("ref",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{if state:=tsNodeNetSocketFor(receiver);state!=nil{state.unref=false};return receiver}),
		"unref":tsNamedNativeMethod("unref",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{if state:=tsNodeNetSocketFor(receiver);state!=nil{state.unref=true};return receiver}),
		"resetAndDestroy":tsNamedNativeMethod("resetAndDestroy",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeNetSocketFor(receiver);if state==nil||state.closed||(!state.connecting&&state.conn==nil&&!state.syncWrite){return tsCallReceiver(loop,tsNodeNetGet(loop,receiver,"destroy"),receiver,tsFSException("Error","ERR_SOCKET_CLOSED","Socket is closed"))};if state.connecting{reset:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{return tsCallReceiver(loop,tsNodeNetGet(loop,receiver,"_reset"),receiver)}));tsCallReceiver(loop,tsNodeNetGet(loop,receiver,"once"),receiver,tsNodeNetKey("connect"),reset);return receiver};if state.conn!=nil{_ = state.conn.Close()};return tsCallReceiver(loop,tsNodeNetGet(loop,receiver,"destroy"),receiver)}),
		"_reset":tsNamedNativeMethod("_reset",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeNetSocketFor(receiver);if state!=nil&&state.conn!=nil{_ = state.conn.Close();state.conn=nil};return tsCallReceiver(loop,tsNodeNetGet(loop,receiver,"destroy"),receiver)}),
		"_getsockname":tsNamedNativeMethod("_getsockname",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeNetSocketFor(receiver);if state==nil||state.localAddress==""||state.closed||state.connecting{return tsObjectValue(tsNewObject())};out:=tsNewObject();out.set("address",tsNodeNetKey(state.localAddress));out.set("port",tsNumberValue(float64(state.localPort)));out.set("family",tsNodeNetKey(state.localFamily));return tsObjectValue(out)}),
		"_getpeername":tsNamedNativeMethod("_getpeername",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeNetSocketFor(receiver);if state==nil||state.remoteAddress==""||state.closed||state.connecting{return tsObjectValue(tsNewObject())};out:=tsNewObject();out.set("address",tsNodeNetKey(state.remoteAddress));out.set("port",tsNumberValue(float64(state.remotePort)));out.set("family",tsNodeNetKey(state.remoteFamily));return tsObjectValue(out)}),
		"destroySoon":tsNamedNativeMethod("destroySoon",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{destroy:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{return tsCallReceiver(loop,tsNodeNetGet(loop,receiver,"destroy"),receiver)}));if tsTruthy(tsNodeNetGet(loop,receiver,"writableFinished")){tsCall(loop,destroy);return receiver};end:=tsNodeNetGet(loop,receiver,"end");if end.kind==tsFunctionKind{tsCallReceiver(loop,end,receiver)};tsCallReceiver(loop,tsNodeNetGet(loop,receiver,"once"),receiver,tsNodeNetKey("finish"),destroy);return receiver}),
		"_unrefTimer":tsNamedNativeMethod("_unrefTimer",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{if state:=tsNodeNetSocketFor(receiver);state!=nil&&state.timeoutTask!=nil{state.timeoutTask.unref=state.unref};return tsU}),
	}
	for name,method:=range methods{p.set(name,method)}
	if p.descriptors==nil{p.descriptors=map[string]*tsDescriptor{}}
	for _,name:=range []string{"bytesRead","bytesWritten","_bytesDispatched","bufferSize","connecting","_connecting","pending","readyState","timeout","localAddress","localPort","localFamily","remoteAddress","remotePort","remoteFamily"}{key:=name;p.descriptors[name]=&tsDescriptor{accessor:true,configurable:true,getter:tsNamedNativeMethod("get "+key,0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeNetSocketFor(receiver);if state==nil{return tsU};switch key{case "bytesRead":return tsNumberValue(float64(state.bytesRead));case "bytesWritten":value:=float64(state.bytesWritten);if writable:=tsNodeWritableStateFor(receiver);writable!=nil{value+=float64(writable.length)};return tsNumberValue(value);case "_bytesDispatched":return tsNumberValue(float64(state.bytesWritten));case "bufferSize":if state.conn==nil||state.closed{return tsU};if writable:=tsNodeWritableStateFor(receiver);writable!=nil{return tsNumberValue(float64(writable.length))};return tsNumberValue(0);case "connecting","_connecting":return tsBooleanValue(state.connecting);case "pending":return tsBooleanValue(state.connecting||(!state.connected&&state.conn==nil)||state.closed);case "timeout":return tsNumberValue(float64(state.timeout/time.Millisecond));case "readyState":if state.connecting{return tsStringReference(tsStringUTF8("opening"))};readable,writable:=tsTruthy(tsNodeNetGet(state.loop,receiver,"readable")),tsTruthy(tsNodeNetGet(state.loop,receiver,"writable"));if readable&&writable{return tsStringReference(tsStringUTF8("open"))};if readable{return tsStringReference(tsStringUTF8("readOnly"))};if writable{return tsStringReference(tsStringUTF8("writeOnly"))};return tsStringReference(tsStringUTF8("closed"));case "localAddress":if state.localAddress==""||state.closed{return tsU};return tsStringReference(tsStringUTF8(state.localAddress));case "localPort":if state.localAddress==""||state.closed{return tsU};return tsNumberValue(float64(state.localPort));case "localFamily":if state.localAddress==""||state.closed{return tsU};return tsStringReference(tsStringUTF8(state.localFamily));case "remoteAddress":if state.remoteAddress==""||state.closed{return tsU};return tsStringReference(tsStringUTF8(state.remoteAddress));case "remotePort":if state.remoteAddress==""||state.closed{return tsU};return tsNumberValue(float64(state.remotePort));default:if state.remoteFamily==""||state.closed{return tsU};return tsStringReference(tsStringUTF8(state.remoteFamily))}})}}
}
func tsNodeNetMakeAddr(address string,port int)tsValue {
	out:=tsNewObject();host,family,_:=tsNodeNetAddress(&net.TCPAddr{IP:net.ParseIP(address),Port:port});if host==""{host=address};out.set("address",tsStringReference(tsStringUTF8(host)));out.set("port",tsNumberValue(float64(port)));out.set("family",tsStringReference(tsStringUTF8(family)));return tsObjectValue(out)
}
func tsNodeNetServerMaybeClose(state *tsNodeNetServerState) {
	if state==nil||!state.closing||state.closed||state.connections>0{return};state.closed=true;tsNodeNetSet(state.loop,state.server,"_handle",tsNull)
	if state.listener!=nil{_ = state.listener.Close();state.listener=nil}
	if state.network=="unix"&&!strings.HasPrefix(state.path,"\x00"){_ = os.Remove(state.path)}
	loop:=state.loop;loop.ready=append(loop.ready,func(){tsNodeNetEmit(loop,state.server,"close")})
}
func tsNodeNetServerAccept(state *tsNodeNetServerState) {
	if state==nil||state.listener==nil||state.closing{return}
	listener:=state.listener
	var conn net.Conn
	var acceptErr error
	task:=state.loop.submit(func()tsResult{
		conn,acceptErr=listener.Accept()
		if acceptErr!=nil{return tsResult{tsNodeNetError(acceptErr,"accept"),true}}
		return tsResult{tsU,false}
	},func(result tsResult){
		if state.closing||state.listener==nil{if conn!=nil{_ = conn.Close()};tsNodeNetServerMaybeClose(state);return}
		if result.rejected{if conn!=nil{_ = conn.Close()};state.closing=true;state.listener=nil;tsNodeNetEmit(state.loop,state.server,"error",result.value);tsNodeNetServerMaybeClose(state);return}
		options:=tsNodeNetGet(state.loop,state.server,"_socketOptions")
		socket:=tsNodeNetSocketCreate(state.loop,options);child:=tsNodeNetSocketFor(socket);child.server=state;tsNodeNetSetConn(child,conn)
		remote:=child.remoteAddress;family:=child.remoteFamily;if family=="IPv4"{family="ipv4"}else if family=="IPv6"{family="ipv6"};blockList:=tsNodeNetGet(state.loop,state.server,"blockList")
		if blockList.kind==tsObjectKind&&tsNodeNetBlockListCheck(state.loop,blockList,tsNodeNetKey(remote),tsNodeNetKey(family)){data:=tsNewObject();data.set("localAddress",tsNodeNetKey(child.localAddress));data.set("localPort",tsNumberValue(float64(child.localPort)));data.set("localFamily",tsNodeNetKey(child.localFamily));data.set("remoteAddress",tsNodeNetKey(child.remoteAddress));data.set("remotePort",tsNumberValue(float64(child.remotePort)));data.set("remoteFamily",tsNodeNetKey(child.remoteFamily));_ = conn.Close();child.closed=true;tsNodeNetEmit(state.loop,state.server,"drop",tsObjectValue(data));tsNodeNetServerAccept(state);return}
		maxConnections:=tsNodeNetGet(state.loop,state.server,"maxConnections");if tsIsNumeric(maxConnections)&&tsNumber(maxConnections)>0&&float64(state.connections)>=tsNumber(maxConnections){data:=tsNewObject();data.set("localAddress",tsNodeNetKey(child.localAddress));data.set("localPort",tsNumberValue(float64(child.localPort)));data.set("localFamily",tsNodeNetKey(child.localFamily));data.set("remoteAddress",tsNodeNetKey(child.remoteAddress));data.set("remotePort",tsNumberValue(float64(child.remotePort)));data.set("remoteFamily",tsNodeNetKey(child.remoteFamily));_ = conn.Close();child.closed=true;tsNodeNetEmit(state.loop,state.server,"drop",tsObjectValue(data));tsNodeNetServerAccept(state);return}
		state.connections++;tsNodeNetSet(state.loop,state.server,"_connections",tsNumberValue(float64(state.connections)));tsNodeNetSet(state.loop,socket,"server",state.server);tsNodeNetSet(state.loop,socket,"_server",state.server);tsNodeNetSet(state.loop,socket,"isServer",tsBooleanValue(true));tsNodeNetSyncSocket(child)
		if tsTruthy(tsNodeNetGet(state.loop,state.server,"pauseOnConnect")){if readable:=tsNodeReadableStateFor(socket);readable!=nil{readable.flowing=false;readable.flowingSet=true;tsNodeReadableSyncState(readable)}}
		tsNodeNetEmitSafe(state.loop,state.server,"connection",socket)
		_,tlsPending:=conn.(*cryptotls.Conn)
		if !tlsPending&&!tsTruthy(tsNodeNetGet(state.loop,state.server,"pauseOnConnect")){tsNodeNetStartRead(child)}
		tsNodeNetServerAccept(state)
	})
	task.unref=state.unref
}
func tsNodeNetListenPort(value tsValue, name string) int {
	if value.kind!=tsStringKind&&!tsIsNumeric(value){tsFSInvalid(name,"a number or string")}
	number:=value.ToNumber()
	if math.IsNaN(number)||math.IsInf(number,0)||math.Trunc(number)!=number||number<0||number>=65536 {
		panic(tsThrown{tsFSException("RangeError","ERR_SOCKET_BAD_PORT",name+" should be >= 0 and < 65536. Received "+tsText(value))})
	}
	return int(number)
}
func tsNodeNetFormatListenError(loop *tsLoop,failure tsValue,address string,port int)tsValue{
	code:=tsText(tsNodeNetGet(loop,failure,"code"));description:="";switch code{case "EADDRINUSE":description="address already in use";case "EACCES":description="permission denied";case "EADDRNOTAVAIL":description="address not available";case "EINVAL":description="invalid argument";default:return failure}
	tsNodeNetSet(loop,failure,"syscall",tsNodeNetKey("listen"));tsNodeNetSet(loop,failure,"address",tsNodeNetKey(address));where:=address;if port!=0{tsNodeNetSet(loop,failure,"port",tsNumberValue(float64(port)));where=address+":"+strconv.Itoa(port)}
	message:="listen "+code+": "+description;if where!=""{message+=" "+where};tsNodeNetSet(loop,failure,"message",tsNodeNetKey(message));return failure
}
func tsNodeNetTCPListenConfig(reusePort,ipv6Only bool)net.ListenConfig{config:=net.ListenConfig{};if reusePort||ipv6Only{config.Control=func(network,address string,raw syscall.RawConn)error{host,_,_:=net.SplitHostPort(address);var socketErr error;controlErr:=raw.Control(func(fd uintptr){if reusePort{socketErr=syscall.SetsockoptInt(int(fd),syscall.SOL_SOCKET,syscall.SO_REUSEPORT,1)};if socketErr==nil&&ipv6Only&&(strings.HasSuffix(network,"6")||strings.Contains(host,":")){socketErr=syscall.SetsockoptInt(int(fd),syscall.IPPROTO_IPV6,syscall.IPV6_V6ONLY,1)}});if controlErr!=nil{return controlErr};return socketErr}};return config}
func tsNodeNetListenNative(loop *tsLoop,receiver tsValue,state *tsNodeNetServerState,path,network,host string,port int,fd tsValue,reusePort,ipv6Only,readableAll,writableAll bool) {
	address:=path;if network=="tcp"{address=net.JoinHostPort(host,strconv.Itoa(port))}
	var listener net.Listener;var err error
	if !tsIsUndefined(fd) { descriptor:=int(tsNumber(fd));file:=os.NewFile(uintptr(descriptor),"node-net-listen-fd");if file==nil{err=syscall.EBADF}else{listener,err=net.FileListener(file);_ = file.Close()} } else if network=="tcp"&&(reusePort||ipv6Only) { listener,err=tsNodeNetTCPListenConfig(reusePort,ipv6Only).Listen(context.Background(),network,address) } else { listener,err=net.Listen(network,address) }
	if err!=nil { if !tsNullish(state.clusterHandle){tsNodeNetClusterHandleClose(loop,state.clusterHandle);state.clusterHandle=tsU};state.clusterListening=false;errorAddress:=host;if network=="unix"{errorAddress=path};failure:=tsNodeNetFormatListenError(loop,tsNodeNetError(err,"listen"),errorAddress,port);loop.ready=append(loop.ready,func(){tsNodeNetEmit(loop,receiver,"error",failure)});return }
	if secure:=tsNodeTLSServerFor(receiver);secure!=nil { secure.mu.RLock();tlsConfig:=secure.config;secure.mu.RUnlock();listener=&tsNodeTLSHandshakeListener{Listener:listener,config:tlsConfig} }
	state.listener=listener;state.closing=false;state.closed=false;state.path=address;state.network=network;tsNodeNetSet(loop,receiver,"_handle",tsObjectValue(tsNewObject()));tsNodeNetSet(loop,receiver,"_socketOptions",state.options)
	if network=="unix" { tsNodeNetSet(loop,receiver,"_address",tsNodeNetKey(listener.Addr().String())) } else if actualHost,listenPort,family:=tsNodeNetAddress(listener.Addr());actualHost!="" { tsNodeNetSet(loop,receiver,"_address",tsNodeNetMakeAddr(actualHost,listenPort));tsNodeNetSet(loop,receiver,"_connectionKey",tsNodeNetKey(family+":"+actualHost+":"+strconv.Itoa(port))) }
	if network=="unix"&&(readableAll||writableAll)&&runtime.GOOS!="windows"&&!strings.HasPrefix(path,"\x00") { if info,statErr:=os.Stat(path);statErr!=nil { _ = listener.Close();state.listener=nil;state.closed=true;tsNodeNetSet(loop,receiver,"_handle",tsNull);failure:=tsNodeNetFormatListenError(loop,tsNodeNetError(statErr,"chmod"),path,0);loop.ready=append(loop.ready,func(){tsNodeNetEmit(loop,receiver,"error",failure)});return } else { permissions:=os.FileMode(0);if readableAll{permissions|=0044};if writableAll{permissions|=0022};if chmodErr:=os.Chmod(path,info.Mode().Perm()|permissions);chmodErr!=nil { _ = listener.Close();state.listener=nil;state.closed=true;tsNodeNetSet(loop,receiver,"_handle",tsNull);failure:=tsNodeNetFormatListenError(loop,tsNodeNetError(chmodErr,"chmod"),path,0);loop.ready=append(loop.ready,func(){tsNodeNetEmit(loop,receiver,"error",failure)});return } } }
	state.loop.resourceClosers=append(state.loop.resourceClosers,func(){_ = listener.Close();if network=="unix"&&!strings.HasPrefix(path,"\x00"){_ = os.Remove(path)}})
	loop.ready=append(loop.ready,func(){if state.listener!=nil&&!state.closed{tsNodeNetEmit(loop,receiver,"listening")}});tsNodeNetServerAccept(state)
}
func tsNodeNetClusterHandleClose(loop *tsLoop,handle tsValue) { if !tsNullish(handle) { close:=tsNodeNetGet(loop,handle,"close");if close.kind==tsFunctionKind{tsCallReceiver(loop,close,handle)} } }
func tsNodeNetErrorValue(loop *tsLoop,value tsValue,operation string) tsValue {
	if value.kind==tsErrorKind{return value}
	code:=tsNodeNetGet(loop,value,"code");if code.kind!=tsStringKind{code=tsNodeNetKey("ERR_SOCKET_ERROR")}
	message:=tsNodeNetGet(loop,value,"message");if tsIsUndefined(message)||tsIsNull(message){message=tsNodeNetKey(operation+" "+tsText(code))}
	return tsErrorValue(&tsRuntimeError{name:"Error",message:tsText(message),fields:map[string]tsValue{"code":code,"syscall":tsNodeNetKey(operation)}})
}
func tsNodeNetCheckClusterBindError(loop *tsLoop,handle tsValue,port int,address string)tsValue {
	if port<=0||tsNullish(handle){return tsNull};getsockname:=tsNodeNetGet(loop,handle,"getsockname");if getsockname.kind!=tsFunctionKind{return tsNull};out:=tsObjectValue(tsNewObject());err:=tsCallReceiver(loop,getsockname,handle,out);if tsTruthy(err){return tsNodeNetErrorValue(loop,err,"listen")};actual:=tsNodeNetGet(loop,out,"port");if tsIsNumeric(actual)&&int(tsNumber(actual))!=port{failure:=tsFSException("Error","EADDRINUSE","listen EADDRINUSE "+address+":"+strconv.Itoa(port));tsNodeNetSet(loop,failure,"syscall",tsNodeNetKey("listen"));tsNodeNetSet(loop,failure,"address",tsNodeNetKey(address));tsNodeNetSet(loop,failure,"port",tsNumberValue(float64(port)));return failure};return tsNull
}
func tsNodeNetClusterConnection(loop *tsLoop,state *tsNodeNetServerState,handle,errValue,clientHandle tsValue) {
	if state==nil||state.closed||!tsStrictEqual(state.clusterHandle,handle) { tsNodeNetClusterHandleClose(loop,clientHandle);return }
	if tsTruthy(errValue) { tsNodeNetEmitSafe(loop,state.server,"error",tsNodeNetErrorValue(loop,errValue,"accept"));return }
	maximum:=tsNodeNetGet(loop,state.server,"maxConnections");if tsIsNumeric(maximum)&&tsNumber(maximum)>0&&float64(state.connections)>=tsNumber(maximum) { tsNodeNetEmitSafe(loop,state.server,"drop");tsNodeNetClusterHandleClose(loop,clientHandle);return }
	socket:=tsNodeNetGet(loop,clientHandle,"_nodeNetSocket")
	if socket.kind!=tsObjectKind||tsNodeNetSocketFor(socket)==nil {
		options:=tsNewObject();options.set("fd",tsNodeNetGet(loop,clientHandle,"fd"));options.set("fdIsRawSocket",tsBooleanValue(true));options.set("pauseOnConnect",tsNodeNetGet(loop,state.server,"pauseOnConnect"));socket=tsNodeNetSocketCreate(loop,tsObjectValue(options));tsNodeNetSocketConnect(loop,socket,[]tsValue{tsObjectValue(options)});return
	}
	child:=tsNodeNetSocketFor(socket);if child.closed { tsNodeNetClusterHandleClose(loop,clientHandle);return }
	if tsTruthy(tsNodeNetGet(loop,state.server,"noDelay")) { child.noDelay=true;if tcp,ok:=tsNodeNetUnderlyingConn(child.conn).(*net.TCPConn);ok{_ = tcp.SetNoDelay(true)} }
	if tsTruthy(tsNodeNetGet(loop,state.server,"keepAlive")) { child.keepAlive=true;delay:=tsNodeNetGet(loop,state.server,"keepAliveInitialDelay");if tsIsNumeric(delay){child.keepAliveDelay=time.Duration(tsNumber(delay))*time.Millisecond};if tcp,ok:=tsNodeNetUnderlyingConn(child.conn).(*net.TCPConn);ok{_ = tcp.SetKeepAlive(true);if child.keepAliveDelay>0{_ = tcp.SetKeepAlivePeriod(child.keepAliveDelay)}} }
	blockList:=tsNodeNetGet(loop,state.server,"blockList");if blockList.kind==tsObjectKind { family:=strings.ToLower(child.remoteFamily);if tsNodeNetBlockListCheck(loop,blockList,tsNodeNetKey(child.remoteAddress),tsNodeNetKey(family)){tsCallReceiver(loop,tsNodeNetGet(loop,socket,"destroy"),socket);return} }
	if child.server!=nil&&child.server!=state { previous:=child.server;previous.connections--;if previous.connections<0{previous.connections=0};tsNodeNetSet(previous.loop,previous.server,"_connections",tsNumberValue(float64(previous.connections)));child.server=nil;tsNodeNetServerMaybeClose(previous) }
	child.server=state;state.connections++;tsNodeNetSet(loop,state.server,"_connections",tsNumberValue(float64(state.connections)));tsNodeNetSet(loop,socket,"server",state.server);tsNodeNetSet(loop,socket,"_server",state.server);tsNodeNetSet(loop,socket,"isServer",tsBooleanValue(true));tsNodeNetSyncSocket(child)
	if tsTruthy(tsNodeNetGet(loop,state.server,"pauseOnConnect")) { if readable:=tsNodeReadableStateFor(socket);readable!=nil{readable.flowing=false;readable.flowingSet=true;tsNodeReadableSyncState(readable)} } else if !child.reading { tsNodeNetStartRead(child) }
	tsNodeNetEmitSafe(loop,state.server,"connection",socket)
}
func tsNodeNetListenInCluster(loop *tsLoop,receiver tsValue,state *tsNodeNetServerState,path,network,host string,clusterAddress tsValue,port int,backlog,fd tsValue,exclusive,ipv6Only,reusePort,readableAll,writableAll bool) bool {
	if exclusive||tsTruthy(tsGet(loop,tsInternalClusterIsPrimaryModule(loop),tsNodeNetKey("isPrimary"))){return false}
	cluster:=tsNodeModuleRequire(loop,[]tsValue{tsNodeNetKey("node:cluster")});getServer:=tsNodeNetGet(loop,cluster,"_getServer");if getServer.kind!=tsFunctionKind{return false}
	address:=host;clusterPort:=port;addressType:=4;addressValue:=clusterAddress;portValue:=tsNumberValue(float64(port));addressTypeValue:=tsNumberValue(4);if !tsNullish(clusterAddress){address=tsText(clusterAddress);if ip:=net.ParseIP(strings.Split(address,"%")[0]);ip!=nil{addressType=4;if ip.To4()==nil{addressType=6}};addressTypeValue=tsNumberValue(float64(addressType))}else{addressValue=tsNull}
	if !tsIsUndefined(fd)&&!tsIsNull(fd)&&tsIsNumeric(fd)&&tsNumber(fd)>=0{addressValue=tsNull;portValue=tsNull;addressTypeValue=tsNull}
	if network=="unix"{address=path;clusterPort=-1;addressValue=tsNodeNetKey(path);portValue=tsNumberValue(-1);addressType=-1;addressTypeValue=tsNumberValue(-1)}
	state.clusterListenID++;listenID:=state.clusterListenID
	flags:=0;if ipv6Only{flags|=1};if reusePort{flags|=2};query:=tsNewObject();query.set("address",addressValue);query.set("port",portValue);query.set("addressType",addressTypeValue);query.set("fd",fd);query.set("flags",tsNumberValue(float64(flags)));query.set("backlog",backlog);query.set("readableAll",tsBooleanValue(readableAll));query.set("writableAll",tsBooleanValue(writableAll));query.set("reusePort",tsBooleanValue(reusePort));if tsNodeTLSServerFor(receiver)!=nil{query.set("sharedOnly",tsBooleanValue(true))}
	callback:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{
		errValue,handle,reply:=tsArg(args,0),tsArg(args,1),tsArg(args,2)
		if listenID!=state.clusterListenID { tsNodeNetClusterHandleClose(loop,handle);return tsU }
		if tsTruthy(errValue) { state.clusterListening=false;state.clusterHandle=tsU;failure:=tsNodeNetErrorValue(loop,errValue,"listen");tsNodeNetFormatListenError(loop,failure,address,port);if hint:=tsNodeNetGet(loop,reply,"bunHint");hint.kind==tsStringKind { if failure.kind==tsErrorKind { (*tsRuntimeError)(failure.ref).message+="\n  note: "+tsText(hint) } };tsNodeNetEmitSafe(loop,receiver,"error",failure);return tsU }
		if bindFailure:=tsNodeNetCheckClusterBindError(loop,handle,port,address);!tsNullish(bindFailure){state.clusterListening=false;state.clusterHandle=tsU;tsNodeNetClusterHandleClose(loop,handle);tsNodeNetEmitSafe(loop,receiver,"error",tsNodeNetFormatListenError(loop,bindFailure,address,port));return tsU}
		sharedFD:=tsNodeNetGet(loop,handle,"sharedFd")
		if tsIsNumeric(sharedFD)&&tsNumber(sharedFD)>=0 { state.clusterHandle=handle;state.closed=false;tsNodeNetSet(loop,receiver,"_usingWorkers",tsBooleanValue(true));tsNodeNetListenNative(loop,receiver,state,"",network,host,port,sharedFD,false,false,readableAll,writableAll);if !state.closed{tsNodeNetSet(loop,handle,"adopted",tsBooleanValue(true))};return tsU }
		if tsNullish(handle) { state.clusterListening=false;tsNodeNetListenNative(loop,receiver,state,path,network,host,port,fd,reusePort,ipv6Only,readableAll,writableAll);return tsU }
		state.clusterHandle=handle;state.path=path;state.network=network;state.closed=false;state.closing=false;tsNodeNetSet(loop,receiver,"_handle",handle);tsNodeNetSet(loop,receiver,"_usingWorkers",tsBooleanValue(true));if path!=""{tsNodeNetSet(loop,handle,"unix",tsNodeNetKey(path));tsNodeNetSet(loop,receiver,"_address",tsNodeNetKey(path))}else if sockname:=tsNodeNetGet(loop,reply,"sockname");!tsNullish(sockname){tsNodeNetSet(loop,receiver,"_address",sockname)};tsNodeNetSet(loop,handle,"_clusterOwner",receiver)
		onconnection:=tsFunctionValue(tsFunc(func(loop *tsLoop,values ...tsValue)tsValue{return tsNodeNetClusterConnection(loop,state,handle,tsArg(values,0),tsArg(values,1))}));tsNodeNetSet(loop,handle,"onconnection",onconnection)
		listen:=tsNodeNetGet(loop,handle,"listen");if listen.kind==tsFunctionKind { backlogNumber:=511;if tsIsNumeric(backlog)&&tsNumber(backlog)!=0{backlogNumber=int(tsNumber(backlog))};tsCallReceiver(loop,listen,handle,tsNumberValue(float64(backlogNumber))) }
		if tsTruthy(tsNodeNetGet(loop,receiver,"_unref")){unref:=tsNodeNetGet(loop,handle,"unref");if unref.kind==tsFunctionKind{tsCallReceiver(loop,unref,handle)}}
		loop.ready=append(loop.ready,func(){if !tsNullish(tsNodeNetGet(loop,receiver,"_handle"))&&!state.closed{tsNodeNetEmit(loop,receiver,"listening")}});return tsU
	}))
	state.clusterHandle=tsU;state.clusterListening=true;tsNodeNetSet(loop,receiver,"_usingWorkers",tsBooleanValue(true));tsCallReceiver(loop,getServer,cluster,receiver,tsObjectValue(query),callback);return true
}
func tsNodeNetListen(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue {
	state:=tsNodeNetServerFor(receiver);if state==nil{return receiver}
	arg:=tsArg(args,0);hostname:=tsArg(args,1);onListen:=tsArg(args,2);backlog:=tsU;exclusive:=false;reusePort:=false;ipv6Only:=false;readableAll:=false;writableAll:=false;path:="";fd:=tsU;port:=0;network:="tcp";host:="::";clusterAddress:=tsNull;hasPort:=false
	if arg.kind==tsStringKind {number:=arg.ToNumber();if !math.IsNaN(number){arg=tsNumberValue(number)}}
	if arg.kind==tsStringKind {
		path=tsText(arg);network="unix"
		if tsIsNumeric(hostname)&&tsNumber(hostname)==math.Trunc(tsNumber(hostname))&&tsNumber(hostname)>0{backlog=hostname}else if hostname.kind==tsFunctionKind{onListen=hostname}
		hostname=tsU
	}else if arg.kind==tsObjectKind||arg.kind==tsInstanceKind||arg.kind==tsArrayKind {
		if pathValue:=tsNodeNetGet(loop,arg,"path");pathValue.kind==tsStringKind{path=tsText(pathValue)}
		if hostValue:=tsNodeNetGet(loop,arg,"host");hostValue.kind==tsStringKind&&tsText(hostValue)!=""{host=tsText(hostValue);clusterAddress=hostValue}
		if portValue:=tsNodeNetGet(loop,arg,"port");!tsIsUndefined(portValue)&&!tsIsNull(portValue){port=tsNodeNetListenPort(portValue,"options.port");hasPort=true}else if tsIsNull(portValue)||tsHas(arg,tsNodeNetKey("port")){port=0;hasPort=true}
		exclusive=tsTruthy(tsNodeNetGet(loop,arg,"exclusive"));reusePort=tsTruthy(tsNodeNetGet(loop,arg,"reusePort"));ipv6Only=tsTruthy(tsNodeNetGet(loop,arg,"ipv6Only"));readableAll=tsTruthy(tsNodeNetGet(loop,arg,"readableAll"));writableAll=tsTruthy(tsNodeNetGet(loop,arg,"writableAll"));backlog=tsNodeNetGet(loop,arg,"backlog")
		fdValue:=tsNodeNetGet(loop,arg,"fd");if tsIsNumeric(fdValue)&&tsNumber(fdValue)>=0{fd=fdValue;port=0;hasPort=true}
		if callback:=tsNodeNetGet(loop,arg,"callback");callback.kind==tsFunctionKind{onListen=callback}
		if signalValue:=tsFSProperty(arg,"signal");!tsIsUndefined(signalValue){signal:=tsFSAbortSignal(arg);_,aborted:=tsAbortReason(signal);if aborted{loop.ready=append(loop.ready,func(){tsNodeNetServerClose(loop,receiver)})}else{signalObject:=tsObjectValue(signal.events.object);onAbort:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{tsNodeNetServerClose(loop,receiver);return tsU}));dispose:=tsInternalAbortListenerAdd(loop,signalObject,onAbort);onClose:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{method:=tsGet(loop,dispose,tsPropertyKey(tsWellKnownSymbol("dispose")));if method.kind==tsFunctionKind{tsCallReceiver(loop,method,dispose)};return tsU}));tsCallReceiver(loop,tsNodeNetGet(loop,receiver,"once"),receiver,tsNodeNetKey("close"),onClose)}}
		if !hasPort {if arg.kind==tsStringKind{if path==""||!tsNodeNetIsPipeName(path){panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_VALUE","The argument 'options' is invalid")})}}else{if !tsHas(arg,tsNodeNetKey("path")){panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_VALUE","The argument 'options' must have the property \"port\" or \"path\"")})};pathValue:=tsNodeNetGet(loop,arg,"path");if pathValue.kind!=tsStringKind||!tsNodeNetIsPipeName(tsText(pathValue)){panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_VALUE","The argument 'options' is invalid")})};path=tsText(pathValue)};network="unix";host=path}
		if hasPort {path=""}
	}else if arg.kind==tsFunctionKind {
		onListen=arg;port=0;hasPort=true
	}else if tsIsUndefined(arg)||tsIsNull(arg) {
		port=0;hasPort=true
	}else if tsIsNumeric(arg) {
		port=tsNodeNetListenPort(arg,"options.port");hasPort=true
		if tsIsNumeric(hostname){backlog=hostname;hostname=tsU}else if hostname.kind==tsFunctionKind{onListen=hostname;hostname=tsU}else if hostname.kind==tsStringKind&&tsIsNumeric(onListen){backlog=onListen;onListen=tsArg(args,3)}
		if hostname.kind==tsStringKind&&tsText(hostname)!=""{host=tsText(hostname);clusterAddress=hostname}
	}else {
		panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_VALUE","The argument 'options' is invalid")})
	}
	if reusePort{exclusive=true};_ = exclusive;_ = backlog
	if hasPort{network="tcp";if host==""{host="::"}}
	if network=="unix"&&(readableAll||writableAll)&&(runtime.GOOS=="linux"||runtime.GOOS=="android")&&strings.HasPrefix(path,"\x00"){panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_VALUE","The argument 'options' cannot set readableAll or writableAll to true when path is an abstract unix socket")})}
	if state.listener!=nil||state.clusterListening||!tsNullish(state.clusterHandle){panic(tsThrown{tsFSException("Error","ERR_SERVER_ALREADY_LISTEN","Listen method has been called more than once without closing.")})}
	if onListen.kind==tsFunctionKind{tsCallReceiver(loop,tsNodeNetGet(loop,receiver,"once"),receiver,tsNodeNetKey("listening"),onListen)}
	if tsNodeNetListenInCluster(loop,receiver,state,path,network,host,clusterAddress,port,backlog,fd,exclusive,ipv6Only,reusePort,readableAll,writableAll){return receiver}
	tsNodeNetListenNative(loop,receiver,state,path,network,host,port,fd,reusePort,ipv6Only,readableAll,writableAll);return receiver
}
func tsNodeNetIsPipeName(path string)bool{number:=tsStringNumber(tsStringUTF8(path));return math.IsNaN(number)||number<0}
func tsNodeNetServerClose(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue {
	state:=tsNodeNetServerFor(receiver);callback:=tsArg(args,0)
	if state!=nil&&state.listener==nil&&(state.clusterListening||!tsNullish(state.clusterHandle)){if callback.kind==tsFunctionKind{tsCallReceiver(loop,tsNodeNetGet(loop,receiver,"once"),receiver,tsNodeNetKey("close"),callback)};state.clusterListenID++;state.clusterListening=false;state.closing=true;handle:=state.clusterHandle;state.clusterHandle=tsU;tsNodeNetClusterHandleClose(loop,handle);tsNodeNetSet(loop,receiver,"_handle",tsNull);tsNodeNetServerMaybeClose(state);return receiver}
	if state==nil||state.listener==nil{if callback.kind==tsFunctionKind{failure:=tsFSException("Error","ERR_SERVER_NOT_RUNNING","Server is not running.");wrapped:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{tsCall(loop,callback,failure);return tsU}));tsCallReceiver(loop,tsNodeNetGet(loop,receiver,"once"),receiver,tsNodeNetKey("close"),wrapped)};if state!=nil{if state.closed&&state.connections==0{loop.ready=append(loop.ready,func(){tsNodeNetEmit(loop,receiver,"close")})}else{state.closing=true;tsNodeNetServerMaybeClose(state)}}else{loop.ready=append(loop.ready,func(){tsNodeNetEmit(loop,receiver,"close")})};return receiver}
	if callback.kind==tsFunctionKind{wrapped:=tsFunctionValue(tsFunc(func(loop *tsLoop,values ...tsValue)tsValue{tsCall(loop,callback,tsNull);return tsU}));tsCallReceiver(loop,tsNodeNetGet(loop,receiver,"once"),receiver,tsNodeNetKey("close"),wrapped)}
	state.clusterListenID++;state.closing=true;state.clusterListening=false;_ = state.listener.Close();state.listener=nil;if !tsNullish(state.clusterHandle){tsNodeNetClusterHandleClose(loop,state.clusterHandle);state.clusterHandle=tsU};if state.network=="unix"&&!strings.HasPrefix(state.path,"\x00"){_ = os.Remove(state.path)};tsNodeNetSet(loop,receiver,"_handle",tsNull);tsNodeNetServerMaybeClose(state);return receiver
}
func tsNodeNetServerCreate(loop *tsLoop,args ...tsValue)tsValue {
	options,listener:=tsArg(args,0),tsArg(args,1)
	if options.kind==tsFunctionKind{listener=options;options=tsObjectValue(tsNewObject())}else if tsNullish(options)||options.kind==tsObjectKind||options.kind==tsInstanceKind{copy:=tsNewObject();tsObjectSpread(copy,options);options=tsObjectValue(copy)}else{panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_TYPE","The \"options\" argument must be of type object or function")})}
	allowHalfOpen:=tsNodeNetGet(loop,options,"allowHalfOpen");if tsIsUndefined(allowHalfOpen){allowHalfOpen=tsBooleanValue(false)}
	keepAlive:=tsTruthy(tsNodeNetGet(loop,options,"keepAlive"))
	keepAliveInitialDelay:=tsNodeNetGet(loop,options,"keepAliveInitialDelay");if !tsIsUndefined(keepAliveInitialDelay){if !tsIsNumeric(keepAliveInitialDelay){tsFSInvalid("options.keepAliveInitialDelay","a number")};if math.IsNaN(tsNumber(keepAliveInitialDelay))||math.IsInf(tsNumber(keepAliveInitialDelay),0){panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","options.keepAliveInitialDelay is out of range")})};if tsNumber(keepAliveInitialDelay)<0{keepAliveInitialDelay=tsNumberValue(0)}}else{keepAliveInitialDelay=tsNumberValue(0)}
	keepAliveInitialDelay=tsNumberValue(float64(int32(tsNumber(keepAliveInitialDelay))))
	highWaterMark:=tsNodeNetGet(loop,options,"highWaterMark");if tsIsUndefined(highWaterMark){streamState:=tsInternalStreamStateModule(loop);getDefault:=tsNodeNetGet(loop,streamState,"getDefaultHighWaterMark");highWaterMark=tsCallReceiver(loop,getDefault,streamState,tsBooleanValue(false))}
	pauseOnConnect:=tsBooleanValue(tsTruthy(tsNodeNetGet(loop,options,"pauseOnConnect")))
	noDelay:=tsBooleanValue(tsTruthy(tsNodeNetGet(loop,options,"noDelay")))
	blockList:=tsNodeNetGet(loop,options,"blockList");if tsTruthy(blockList){check:=tsNodeNetGet(loop,tsNodeNetSingleton.blockListCtor,"isBlockList");if check.kind!=tsFunctionKind||!tsTruthy(tsCallReceiver(loop,check,tsNodeNetSingleton.blockListCtor,blockList)){panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_TYPE","The \"options.blockList\" property must be an instance of net.BlockList")})}}
	object:=tsNewObject();object.prototype=tsObjectValue(tsNodeNetSingleton.serverProto);server:=tsNodeEventCreate(loop,tsObjectValue(object),tsU);state:=&tsNodeNetServerState{loop:loop,server:server,options:options};tsNodeNetSingleton.servers[uintptr(server.ref)]=state
	object.set("_connections",tsNumberValue(0));object.set("_handle",tsNull);object.set("_usingWorkers",tsBooleanValue(false));object.set("workers",tsArrayValue(&tsArray{}));object.set("_unref",tsBooleanValue(false));object.set("_socketOptions",options);object.set("allowHalfOpen",allowHalfOpen);object.set("pauseOnConnect",pauseOnConnect);object.set("noDelay",noDelay);object.set("keepAlive",tsBooleanValue(keepAlive));object.set("keepAliveInitialDelay",keepAliveInitialDelay);object.set("highWaterMark",highWaterMark);object.set("maxConnections",tsNodeNetGet(loop,options,"maxConnections"));if tsTruthy(blockList){tsNodeNetSet(loop,server,"blockList",blockList)}
	if listener.kind==tsFunctionKind{tsCallReceiver(loop,tsNodeNetGet(loop,server,"on"),server,tsNodeNetKey("connection"),listener)};return server
}
func tsNodeNetServerCtor(loop *tsLoop)tsValue {constructor:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeNetServerCreate(loop,args...)});constructor.name="Server";constructor.length=2;constructor.constructible=true;constructor.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeNetServerCreate(loop,args...)};constructor.receiverCall=func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeNetServerCreate(loop,args...)};constructor.properties=tsNewObject();constructor.properties.set("prototype",tsObjectValue(tsNodeNetSingleton.serverProto));constructor.properties.descriptors=map[string]*tsDescriptor{"prototype":{writable:false}};return tsFunctionValue(constructor)}
func tsNodeNetServerListening(loop *tsLoop,receiver tsValue)bool{return !tsNullish(tsNodeNetGet(loop,receiver,"_handle"))}
func tsNodeNetServerAddress(loop *tsLoop,receiver tsValue)tsValue{if !tsNodeNetServerListening(loop,receiver){return tsNull};address:=tsNodeNetGet(loop,receiver,"_address");if !tsNullish(address){return address};state:=tsNodeNetServerFor(receiver);if state==nil||tsNullish(state.clusterHandle){return tsNull};if unix:=tsNodeNetGet(loop,state.clusterHandle,"unix");unix.kind==tsStringKind{return unix};getsockname:=tsNodeNetGet(loop,state.clusterHandle,"getsockname");if getsockname.kind==tsFunctionKind{out:=tsObjectValue(tsNewObject());err:=tsCallReceiver(loop,getsockname,state.clusterHandle,out);if !tsTruthy(err){return out}};return tsNull}
func tsNodeNetServerInstall(loop *tsLoop) {
	p:=tsNodeNetSingleton.serverProto;p.set("listen",tsNamedNativeMethod("listen",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeNetListen(loop,receiver,args...)}));p.set("close",tsNamedNativeMethod("close",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeNetServerClose(loop,receiver,args...)}));p.set("address",tsNamedNativeMethod("address",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{return tsNodeNetServerAddress(loop,receiver)}));p.set("getConnections",tsNamedNativeMethod("getConnections",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeNetServerFor(receiver);count:=0;if state!=nil&&tsNodeNetServerListening(loop,receiver){count=state.connections};if callback:=tsArg(args,0);callback.kind==tsFunctionKind{tsCall(loop,callback,tsNull,tsNumberValue(float64(count)))};return receiver}));for name,unref:=range map[string]bool{"ref":false,"unref":true}{method:=name;p.set(method,tsNamedNativeMethod(method,0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{if state:=tsNodeNetServerFor(receiver);state!=nil{state.unref=unref;tsNodeNetSet(loop,receiver,"_unref",tsBooleanValue(unref));handle:=tsNodeNetGet(loop,receiver,"_handle");operation:=method;if unref{operation="unref"}else{operation="ref"};if fn:=tsNodeNetGet(loop,handle,operation);fn.kind==tsFunctionKind{tsCallReceiver(loop,fn,handle)}};return receiver}))}
	p.set("_emitCloseIfDrained",tsNamedNativeMethod("_emitCloseIfDrained",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeNetServerFor(receiver);if state!=nil&&!tsNodeNetServerListening(loop,receiver)&&state.connections==0{loop.ready=append(loop.ready,func(){tsNodeNetEmit(loop,receiver,"close")})};return tsU}))
	p.set("getsockname",tsNamedNativeMethod("getsockname",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{out:=tsArg(args,0);address:=tsNodeNetGet(loop,receiver,"_address");tsNodeNetSet(loop,out,"port",tsNodeNetGet(loop,address,"port"));return out}))
	events:=tsNodeEventsModule(loop);captureSymbol:=tsNodeNetGet(loop,events,"captureRejectionSymbol");if captureSymbol.kind==tsSymbolKind{p.set(tsPropertyKey(captureSymbol),tsNamedNativeMethod("[captureRejection]",3,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{failure,event,socket:=tsArg(args,0),tsArg(args,1),tsArg(args,2);eventName:=tsText(event);if eventName=="connection"||eventName=="secureConnection"{destroy:=tsNodeNetGet(loop,socket,"destroy");if destroy.kind==tsFunctionKind{tsCallReceiver(loop,destroy,socket,failure)}}else{tsNodeNetEmit(loop,receiver,"error",failure)};return tsU}))}
	p.set(tsPropertyKey(tsWellKnownSymbol("asyncDispose")),tsNamedNativeMethod("[Symbol.asyncDispose]",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{if !tsNodeNetServerListening(loop,receiver){return tsPromiseValue(loop.resolved(tsU,false))};promise:=loop.promise();callback:=tsFunctionValue(tsFunc(func(_ *tsLoop,args ...tsValue)tsValue{err:=tsArg(args,0);if tsTruthy(err){promise.settle(tsResult{err,true})}else{promise.settle(tsResult{tsU,false})};return tsU}));tsNodeNetServerClose(loop,receiver,[]tsValue{callback}...);return tsPromiseValue(promise)}))
	if p.descriptors==nil{p.descriptors=map[string]*tsDescriptor{}};p.descriptors["listening"]=&tsDescriptor{accessor:true,configurable:true,getter:tsNamedNativeMethod("get listening",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{return tsBooleanValue(tsNodeNetServerListening(loop,receiver))})};p.descriptors["connections"]=&tsDescriptor{accessor:true,configurable:true,getter:tsNamedNativeMethod("get connections",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeNetServerFor(receiver);if state==nil{return tsNull};return tsNumberValue(float64(state.connections))})}
}
func tsNodeNetBlockListFor(value tsValue)*tsNodeNetBlockListState{key:=tsNodeNetPointer(value);tsNodeNetSingleton.blockListMapMu.RLock();state:=tsNodeNetSingleton.blocks[key];tsNodeNetSingleton.blockListMapMu.RUnlock();return state}
func tsNodeNetBlockListInvalidThis(){panic(tsThrown{tsFSException("TypeError","ERR_INVALID_THIS","Value of 'this' must be of type net.BlockList")})}
func tsNodeNetBlockListAdd(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue {
	state:=tsNodeNetBlockListFor(receiver);if state==nil{tsNodeNetBlockListInvalidThis()}
	method:=tsText(tsArg(args,0));firstValue,secondValue,familyValue:=tsArg(args,1),tsArg(args,2),tsArg(args,3)
	rule:=tsNodeNetBlockRule{kind:method}
	switch method {
	case "addAddress":
		rule.ip,rule.family=tsNodeNetBlockIP(loop,firstValue,"address",familyValue)
		rule.text=fmt.Sprintf("Address: %s %s",tsNodeNetBlockFamilyLabel(rule.family),rule.ip.String())
	case "addRange":
		start,startFamily:=tsNodeNetBlockIP(loop,firstValue,"start",familyValue);end,endFamily:=tsNodeNetBlockIP(loop,secondValue,"end",familyValue)
		if startFamily==endFamily&&len(start)==len(end)&&tsNodeNetIPCompare(start,end)>0{panic(tsThrown{tsFSException("RangeError","ERR_INVALID_ARG_VALUE","The argument 'start' must come before end")})}
		rule.ip,rule.end,rule.family,rule.endFamily=start,end,startFamily,endFamily
		rule.text=fmt.Sprintf("Range: %s %s-%s",tsNodeNetBlockFamilyLabel(rule.family),start.String(),end.String())
	case "addSubnet":
		ip,family:=tsNodeNetBlockIP(loop,firstValue,"network",familyValue);prefixValue:=secondValue
		bits:=128;if family=="ipv4"{bits=32};if !tsIsNumeric(prefixValue){tsFSInvalid("prefix","a number")};if tsNumber(prefixValue)!=math.Trunc(tsNumber(prefixValue))||tsNumber(prefixValue)<0||tsNumber(prefixValue)>float64(bits){panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE",fmt.Sprintf("prefix is out of range; expected an integer between 0 and %d",bits))})}
		prefix:=int(tsNumber(prefixValue));rule.family,rule.ip=family,ip
		mask:=net.CIDRMask(prefix,bits);rule.network=&net.IPNet{IP:ip.Mask(mask),Mask:mask}
		rule.text=fmt.Sprintf("Subnet: %s %s/%d",tsNodeNetBlockFamilyLabel(family),ip.String(),prefix)
	default:return tsU
	}
	// Bun inserts the newest entry at index zero; the public rules getter uses
	// this same order.
	state.mu.Lock();state.rules=append([]tsNodeNetBlockRule{rule},state.rules...);state.mu.Unlock();return tsU
}
func tsNodeNetBlockFamilyLabel(family string)string{if family=="ipv6"{return "IPv6"};return "IPv4"}
func tsNodeNetIPCompare(a,b net.IP)int{if len(a)!=len(b){if len(a)<len(b){return -1};return 1};for i:=range a{if a[i]<b[i]{return -1};if a[i]>b[i]{return 1}};return 0}
func tsNodeNetBlockIP(loop *tsLoop,value tsValue,argName string,familyValue tsValue)(net.IP,string){
	family:="ipv4"
	if tsNodeNetIsSocketAddress(value){
		address:=tsNodeNetGet(loop,value,"address");providedFamily:=tsNodeNetGet(loop,value,"family")
		family=strings.ToLower(tsText(providedFamily));text:=tsText(address);ip:=net.ParseIP(text);return tsNodeNetNormalizeBlockIP(ip,family,argName,text),family
	}
	if value.kind!=tsStringKind{tsFSInvalid(argName,"a string or SocketAddress")}
	if !tsIsUndefined(familyValue){if familyValue.kind!=tsStringKind{tsFSInvalid("family","a string")};family=strings.ToLower(tsText(familyValue))}
	text:=tsText(value);ip:=net.ParseIP(text);return tsNodeNetNormalizeBlockIP(ip,family,argName,text),family
}
func tsNodeNetNormalizeBlockIP(ip net.IP,family,argName,original string)net.IP{
	switch family{
	case "ipv4":if ip!=nil&&ip.To4()!=nil&&!strings.Contains(original,":"){return append(net.IP(nil),ip.To4()...)}
	case "ipv6":if ip!=nil&&ip.To16()!=nil&&strings.Contains(original,":"){return append(net.IP(nil),ip.To16()...)}
	default:panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_VALUE","The property 'family' must be 'ipv4' or 'ipv6'. Received "+family)})
	}
	panic(tsThrown{tsFSException("TypeError","ERR_INVALID_IP_ADDRESS","Invalid IP address for "+argName)})
}
func tsNodeNetSocketAddressFor(value tsValue)*tsNodeNetSocketAddressState{key:=tsNodeNetPointer(value);tsNodeNetSingleton.socketAddressMapMu.RLock();state:=tsNodeNetSingleton.socketAddresses[key];tsNodeNetSingleton.socketAddressMapMu.RUnlock();return state}
func tsNodeNetIsSocketAddress(value tsValue)bool{return tsNodeNetSocketAddressFor(value)!=nil}
func tsNodeNetBlockListCheck(loop *tsLoop,receiver,address,family tsValue)bool {
	state:=tsNodeNetBlockListFor(receiver);if state==nil{tsNodeNetBlockListInvalidThis()}
	fam:="ipv4";var ip net.IP
	if tsNodeNetIsSocketAddress(address){addressValue:=tsNodeNetGet(loop,address,"address");familyValue:=tsNodeNetGet(loop,address,"family");fam=strings.ToLower(tsText(familyValue));text:=tsText(addressValue);ip=tsNodeNetNormalizeBlockIPForCheck(net.ParseIP(text),fam,text)}else{if address.kind!=tsStringKind{tsFSInvalid("address","a string or SocketAddress")};if !tsIsUndefined(family){if family.kind!=tsStringKind{tsFSInvalid("family","a string")};fam=strings.ToLower(tsText(family))};if fam!="ipv4"&&fam!="ipv6"{return false};ip=tsNodeNetNormalizeBlockIPForCheck(net.ParseIP(tsText(address)),fam,tsText(address))};if ip==nil{return false}
	state.mu.RLock();defer state.mu.RUnlock();for _,rule:=range state.rules{
		if rule.kind=="addAddress"&&rule.family==fam&&tsNodeNetIPCompare(rule.ip,ip)==0{return true}
		if rule.kind=="addRange"&&rule.family==fam&&rule.endFamily==fam&&tsNodeNetIPCompare(ip,rule.ip)>=0&&tsNodeNetIPCompare(ip,rule.end)<=0{return true}
		if rule.kind=="addSubnet"&&rule.network!=nil{candidate:=ip;if rule.family=="ipv6"&&fam=="ipv4"{candidate=net.ParseIP("::ffff:"+ip.String())};if rule.network.Contains(candidate){return true}}
	};return false
}
func tsNodeNetNormalizeBlockIPForCheck(ip net.IP,family,original string)net.IP{if ip==nil{return nil};if family=="ipv4"{if v4:=ip.To4();v4!=nil&&!strings.Contains(original,":"){return append(net.IP(nil),v4...) };return nil};if family=="ipv6"&&strings.Contains(original,":")&&ip.To16()!=nil{return append(net.IP(nil),ip.To16()...)};return nil}
func tsNodeNetBlockListSerializationNonce(value tsValue)(uint64,bool){key:=tsNodeNetPointer(value);if tsNodeNetBlockListFor(value)==nil{return 0,false};tsNodeNetSingleton.blockListCloneMu.Lock();defer tsNodeNetSingleton.blockListCloneMu.Unlock();if nonce,ok:=tsNodeNetSingleton.blockListObjectNonce[key];ok{return nonce,true};nonce:=tsNodeNetSingleton.blockListNextNonce.Add(1);if nonce==0{nonce=tsNodeNetSingleton.blockListNextNonce.Add(1)};tsNodeNetSingleton.blockListObjectNonce[key]=nonce;tsNodeNetSingleton.blockListSerialized[nonce]=value;return nonce,true}
func tsNodeNetBlockListFromSerializationNonce(nonce uint64)(tsValue,bool){tsNodeNetSingleton.blockListCloneMu.Lock();defer tsNodeNetSingleton.blockListCloneMu.Unlock();value,ok:=tsNodeNetSingleton.blockListSerialized[nonce];return value,ok}
func tsNodeNetBlockListCtor(loop *tsLoop)tsValue {
	p:=tsNodeNetSingleton.blockListProto;p.set("addAddress",tsNamedNativeMethod("addAddress",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeNetBlockListAdd(loop,receiver,append([]tsValue{tsNodeNetKey("addAddress")},args...)...)}));p.set("addRange",tsNamedNativeMethod("addRange",3,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeNetBlockListAdd(loop,receiver,append([]tsValue{tsNodeNetKey("addRange")},args...)...)}));p.set("addSubnet",tsNamedNativeMethod("addSubnet",3,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeNetBlockListAdd(loop,receiver,append([]tsValue{tsNodeNetKey("addSubnet")},args...)...)}));p.set("check",tsNamedNativeMethod("check",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsBooleanValue(tsNodeNetBlockListCheck(loop,receiver,tsArg(args,0),tsArg(args,1)))}))
	if p.descriptors==nil{p.descriptors=map[string]*tsDescriptor{}};p.descriptors["rules"]=&tsDescriptor{accessor:true,configurable:true,getter:tsNamedNativeMethod("get rules",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{out:=&tsArray{};state:=tsNodeNetBlockListFor(receiver);if state==nil{tsNodeNetBlockListInvalidThis()};state.mu.RLock();for _,rule:=range state.rules{out.values=append(out.values,tsNodeNetKey(rule.text))};state.mu.RUnlock();return tsArrayValue(out)})}
	var constructor *tsFunction;constructor=tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{object:=tsNewObject();object.prototype=tsObjectValue(p);tsNodeNetSingleton.blockListMapMu.Lock();tsNodeNetSingleton.blocks[uintptr(unsafe.Pointer(object))]=&tsNodeNetBlockListState{};tsNodeNetSingleton.blockListMapMu.Unlock();return tsObjectValue(object)});constructor.name="BlockList";constructor.length=0;constructor.constructible=true;constructor.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsCall(loop,tsFunctionValue(constructor))};constructor.properties=tsNewObject();constructor.properties.set("prototype",tsObjectValue(p));constructor.properties.set("isBlockList",tsNamedNativeMethod("isBlockList",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsBooleanValue(tsNodeNetBlockListFor(tsArg(args,0))!=nil)}));tsNodeNetSingleton.blockListProto=p;return tsFunctionValue(constructor)
}
func tsNodeNetSocketAddressInvalidThis(){panic(tsThrown{tsFSException("TypeError","ERR_INVALID_THIS","Value of 'this' must be of type net.SocketAddress")})}
func tsNodeNetSocketAddressFamily(value tsValue)(string,bool){
	if value.kind==tsStringKind { family:=tsText(value);if strings.EqualFold(family,"ipv4"){return "ipv4",true};if strings.EqualFold(family,"ipv6"){return "ipv6",true};return "",false }
	if tsIsNumeric(value) { number:=tsNumber(value);if math.Trunc(number)==number { if number==float64(syscall.AF_INET){return "ipv4",true};if number==float64(syscall.AF_INET6){return "ipv6",true} } }
	return "",false
}
func tsNodeNetSocketAddressBadFamily(value tsValue){
	message:="The \"options.family\" property must be 'ipv4' or 'ipv6'."
	if tsIsNumeric(value){message="The \"options.family\" property must be AF_INET or AF_INET6."}else if value.kind!=tsStringKind{message="The \"options.family\" property must be a string or number."}
	panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_VALUE",message)})
}
func tsNodeNetSocketAddressIPv6String(ip net.IP)string{
	bytes:=ip.To16();if bytes==nil{return ""};groups:=make([]uint16,8);for index:=range groups{groups[index]=uint16(bytes[index*2])<<8|uint16(bytes[index*2+1])}
	bestStart,bestLength:=-1,1;for start:=0;start<len(groups);{if groups[start]!=0{start++;continue};end:=start;for end<len(groups)&&groups[end]==0{end++};if end-start>bestLength{bestStart,bestLength=start,end-start};start=end}
	if bestStart<0{parts:=make([]string,8);for index,group:=range groups{parts[index]=strconv.FormatUint(uint64(group),16)};return strings.Join(parts,":")}
	left,right:=make([]string,0,bestStart),make([]string,0,len(groups)-bestStart-bestLength);for index:=0;index<bestStart;index++{left=append(left,strconv.FormatUint(uint64(groups[index]),16))};for index:=bestStart+bestLength;index<len(groups);index++{right=append(right,strconv.FormatUint(uint64(groups[index]),16))};return strings.Join(left,":")+"::"+strings.Join(right,":")
}
func tsNodeNetSocketAddressParseIPv4Number(text string)(uint64,bool){
	if text==""{return 0,false};base:=10
	if len(text)>=2&&text[0]=='0'&&(text[1]=='x'||text[1]=='X'){base=16;text=text[2:];if text==""{return 0,true}}else if len(text)>=2&&text[0]=='0'{base=8;text=text[1:];if text==""{return 0,true}}
	value,err:=strconv.ParseUint(text,base,64);return value,err==nil
}
func tsNodeNetSocketAddressCanonicalIPv4(host string)net.IP{
	parts:=strings.Split(host,".");if len(parts)>1&&parts[len(parts)-1]==""{parts=parts[:len(parts)-1]};if len(parts)==0||len(parts)>4{return nil};values:=make([]uint64,len(parts));for index,part:=range parts{value,ok:=tsNodeNetSocketAddressParseIPv4Number(part);if !ok{return nil};if index<len(parts)-1&&value>255{return nil};values[index]=value}
	remainingBits:=uint((5-len(parts))*8);maxLast:=uint64(1)<<remainingBits;if remainingBits==64{maxLast=^uint64(0)};if values[len(values)-1]>=maxLast{return nil}
	address:=values[len(values)-1];for index:=0;index<len(values)-1;index++{address|=values[index]<<uint(24-8*index)}
	return net.IPv4(byte(address>>24),byte(address>>16),byte(address>>8),byte(address))
}
func tsNodeNetSocketAddressNew(loop *tsLoop,options tsValue,parsed bool,addressOverride string,portOverride int)tsValue{
	family:="ipv4";addressValue:=tsU;portValue:=tsU;flowValue:=tsU
	if !parsed&&!tsIsUndefined(options){if options.kind!=tsObjectKind&&options.kind!=tsArrayKind&&options.kind!=tsFunctionKind&&options.kind!=tsInstanceKind&&options.kind!=tsClassKind{panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_TYPE","The \"options\" argument must be of type object")})};familyValue:=tsNodeNetGet(loop,options,"family");if !tsIsUndefined(familyValue){var ok bool;family,ok=tsNodeNetSocketAddressFamily(familyValue);if !ok{tsNodeNetSocketAddressBadFamily(familyValue)}};addressValue=tsNodeNetGet(loop,options,"address");portValue=tsNodeNetGet(loop,options,"port");flowValue=tsNodeNetGet(loop,options,"flowlabel")}
	if parsed{addressValue=tsNodeNetKey(addressOverride);portValue=tsNumberValue(float64(portOverride));if strings.Contains(addressOverride,":"){family="ipv6"}}
	port:=0;if !tsIsUndefined(portValue){if !tsIsNumeric(portValue)||math.IsNaN(tsNumber(portValue))||math.IsInf(tsNumber(portValue),0){message:="The \"options.port\" argument must be a valid IP port number.";if !tsIsNumeric(portValue){message+=" Got "+tsValueTypeName(portValue)+"."};panic(tsThrown{tsFSException("RangeError","ERR_SOCKET_BAD_PORT",message)})};number:=tsNumber(portValue);if math.Trunc(number)!=number||number<0||number>65535{panic(tsThrown{tsFSException("RangeError","ERR_SOCKET_BAD_PORT","The \"options.port\" argument must be a valid IP port number.")})};port=int(number)}
	flow:=uint32(0);if !tsIsUndefined(flowValue){if !tsIsNumeric(flowValue){tsFSInvalid("options.flowlabel","a number")};number:=tsNumber(flowValue);if math.IsNaN(number)||math.IsInf(number,0)||math.Trunc(number)!=number||number<0||number>float64(^uint32(0)){panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","options.flowlabel is out of range")})};flow=uint32(number)}
	address:="127.0.0.1";if family=="ipv6"{address="::"};zoneStripped:=false;if !tsIsUndefined(addressValue){if addressValue.kind!=tsStringKind{tsFSInvalid("options.address","a string")};address=tsText(addressValue);if family=="ipv6"{if zone:=strings.IndexByte(address,'%');zone>=0{address=address[:zone];zoneStripped=true}}}
	var ip net.IP;if family=="ipv4"{if !strings.Contains(address,":"){ip=net.ParseIP(address);if ip!=nil{ip=ip.To4()}}}else if strings.Contains(address,":"){ip=net.ParseIP(address);if ip!=nil{ip=ip.To16()}}
	if ip==nil{panic(tsThrown{tsFSException("TypeError","ERR_INVALID_IP_ADDRESS","Invalid socket address")})}
	if parsed{if family=="ipv4"{address=ip.String()}else{address=tsNodeNetSocketAddressIPv6String(ip)}}else if zoneStripped{address=tsNodeNetSocketAddressIPv6String(ip)}else if tsIsUndefined(addressValue){if family=="ipv4"{address="127.0.0.1"}else{address="::"}}
	if family=="ipv4"{flow=0};object:=tsNewObject();object.prototype=tsObjectValue(tsNodeNetSingleton.socketAddressProto);state:=&tsNodeNetSocketAddressState{address:address,family:family,port:port,flowLabel:flow};key:=uintptr(unsafe.Pointer(object));tsNodeNetSingleton.socketAddressMapMu.Lock();tsNodeNetSingleton.socketAddresses[key]=state;tsNodeNetSingleton.socketAddressMapMu.Unlock();return tsObjectValue(object)
}
func tsNodeNetSocketAddressParse(loop *tsLoop,input tsValue)tsValue{
	if input.kind!=tsStringKind{tsFSInvalid("input","a string")};parsed,err:=url.Parse("http://"+tsText(input));if err!=nil||parsed.Host==""{return tsU};host:=parsed.Hostname();if host==""{return tsU};port:=0;if portText:=parsed.Port();portText!=""{value,parseErr:=strconv.ParseUint(portText,10,32);if parseErr!=nil{return tsU};if value<=65535{port=int(value)}}
	if strings.Contains(host,":"){ip:=net.ParseIP(host);if ip==nil{return tsU};return tsNodeNetSocketAddressNew(loop,tsU,true,tsNodeNetSocketAddressIPv6String(ip),port)}
	ip:=tsNodeNetSocketAddressCanonicalIPv4(host);if ip==nil{return tsU};return tsNodeNetSocketAddressNew(loop,tsU,true,ip.String(),port)
}
func tsNodeNetSocketAddressCtor(loop *tsLoop)tsValue {
	p:=tsNodeNetSingleton.socketAddressProto;if p.descriptors==nil{p.descriptors=map[string]*tsDescriptor{}}
	for _,field:=range []string{"address","family","port","flowlabel","addrFamily","flowLabel"}{name:=field;getter:=tsNamedNativeMethod("get "+name,0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeNetSocketAddressFor(receiver);if state==nil{tsNodeNetSocketAddressInvalidThis()};switch name{case "address":return tsNodeNetKey(state.address);case "family":return tsNodeNetKey(state.family);case "port":return tsNumberValue(float64(state.port));case "addrFamily":if state.family=="ipv6"{return tsNumberValue(float64(syscall.AF_INET6))};return tsNumberValue(float64(syscall.AF_INET));default:return tsNumberValue(float64(state.flowLabel))}});p.descriptors[name]=&tsDescriptor{accessor:true,configurable:true,getter:getter}}
	p.set("toJSON",tsNamedNativeMethod("toJSON",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeNetSocketAddressFor(receiver);if state==nil{tsNodeNetSocketAddressInvalidThis()};object:=tsNewObject();object.set("address",tsNodeNetKey(state.address));object.set("family",tsNodeNetKey(state.family));object.set("port",tsNumberValue(float64(state.port)));object.set("flowlabel",tsNumberValue(float64(state.flowLabel)));return tsObjectValue(object)}))
	constructor:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeNetSocketAddressNew(loop,tsArg(args,0),false,"",0)});constructor.name="SocketAddress";constructor.length=1;constructor.constructible=true;constructor.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeNetSocketAddressNew(loop,tsArg(args,0),false,"",0)};constructor.properties=tsNewObject();constructor.properties.set("prototype",tsObjectValue(p));p.set("constructor",tsFunctionValue(constructor));constructor.properties.set("isSocketAddress",tsNamedNativeMethod("isSocketAddress",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsBooleanValue(tsNodeNetIsSocketAddress(tsArg(args,0)))}));constructor.properties.set("parse",tsNamedNativeMethod("parse",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeNetSocketAddressParse(loop,tsArg(args,0))}));return tsFunctionValue(constructor)
}
func tsNodeNetNormalizeArgs(args []tsValue)tsValue {
	options:=tsNewObject();callback:=tsNull;if len(args)==0{return tsNodeNetMarkNormalizedArgs([]tsValue{tsObjectValue(options),callback})};first:=args[0];if first.kind==tsObjectKind||first.kind==tsInstanceKind||first.kind==tsArrayKind{options=tsNodeNetObjectValue(first)}else if first.kind==tsStringKind{if tsNodeNetIsPipeName(tsText(first)){options.set("path",first)}else{options.set("port",first)}}else{options.set("port",first);if len(args)>1&&args[1].kind==tsStringKind{options.set("host",args[1])}};if len(args)>0&&args[len(args)-1].kind==tsFunctionKind{callback=args[len(args)-1]};return tsNodeNetMarkNormalizedArgs([]tsValue{tsObjectValue(options),callback})
}
func tsNodeNetNormalizedArgsKey()tsValue{if tsIsUndefined(tsNodeNetSingleton.normalizedArgsSymbol){tsNodeNetSingleton.normalizedArgsSymbol=tsNewSymbol(tsStringUTF8("node.net.normalizedArgs"))};return tsNodeNetSingleton.normalizedArgsSymbol}
func tsNodeNetMarkNormalizedArgs(values []tsValue)tsValue{array:=&tsArray{values:values,properties:map[string]tsValue{}};array.properties[tsPropertyKey(tsNodeNetNormalizedArgsKey())]=tsBooleanValue(true);return tsArrayValue(array)}
func tsNodeNetIsNormalizedArgs(loop *tsLoop,value tsValue)bool{return value.kind==tsArrayKind&&tsTruthy(tsGet(loop,value,tsNodeNetNormalizedArgsKey()))}
func tsNodeNetObjectValue(value tsValue)tsValue{if value.kind==tsObjectKind{return value};out:=tsNewObject();if value.kind==tsArrayKind{for index,item:=range (*tsArray)(value.ref).values{out.set(strconv.Itoa(index),item)}};return tsObjectValue(out)}
type tsNodeNetKeyLogFile struct{path string}
func (writer tsNodeNetKeyLogFile)Write(line []byte)(int,error){file,err:=os.OpenFile(writer.path,os.O_WRONLY|os.O_CREATE|os.O_APPEND,0600);if err!=nil{return len(line),nil};_,_ = file.Write(line);_ = file.Close();return len(line),nil}
func tsNodeNetKeyLogWriter(loop *tsLoop)io.Writer{tsNodeNetModule(loop);if !tsNodeNetSingleton.tlsKeylogConfigured{return nil};if !tsNodeNetSingleton.tlsKeylogWarned{tsNodeNetSingleton.tlsKeylogWarned=true;process:=tsNodeProcessModule(loop);emit:=tsNodeNetGet(loop,process,"emitWarning");warning:=tsNodeNetKey("Using --tls-keylog makes TLS connections insecure by writing secret key material to file "+tsNodeNetSingleton.tlsKeylogPath);if emit.kind==tsFunctionKind{tsCall(loop,emit,warning)}};return tsNodeNetKeyLogFile{path:tsNodeNetSingleton.tlsKeylogPath}}
func tsNodeNetReadExecArgv(loop *tsLoop){process:=tsNodeProcessModule(loop);argv:=tsNodeNetGet(loop,process,"execArgv");if argv.kind!=tsArrayKind{return};values:=(*tsArray)(argv.ref).values;for index:=0;index<len(values);index++{arg:=tsText(values[index]);switch arg{case "--no-network-family-autoselection","--no-enable-network-family-autoselection":tsNodeNetSingleton.defaultAutoSelectFamily=false;case "--network-family-autoselection","--enable-network-family-autoselection":tsNodeNetSingleton.defaultAutoSelectFamily=true;case "--network-family-autoselection-attempt-timeout":if index+1<len(values){if n,err:=strconv.Atoi(tsText(values[index+1]));err==nil&&n>=1{if n<10{n=10};tsNodeNetSingleton.defaultAutoSelectFamilyAttemptTimeout=n};index++};case "--tls-keylog":if index+1<len(values){tsNodeNetSingleton.tlsKeylogConfigured=true;tsNodeNetSingleton.tlsKeylogPath=tsText(values[index+1]);index++}};if strings.HasPrefix(arg,"--network-family-autoselection-attempt-timeout="){if n,err:=strconv.Atoi(strings.TrimPrefix(arg,"--network-family-autoselection-attempt-timeout="));err==nil&&n>=1{if n<10{n=10};tsNodeNetSingleton.defaultAutoSelectFamilyAttemptTimeout=n}};if strings.HasPrefix(arg,"--tls-keylog="){tsNodeNetSingleton.tlsKeylogConfigured=true;tsNodeNetSingleton.tlsKeylogPath=strings.TrimPrefix(arg,"--tls-keylog=")}}}
func tsNodeNetCreateConnection(loop *tsLoop,args []tsValue)tsValue {
	normalized:=tsNodeNetNormalizeArgs(args);values:=(*tsArray)(normalized.ref).values;socket:=tsNodeNetSocketCreate(loop,tsArg(values,0));timeout:=tsNodeNetGet(loop,tsArg(values,0),"timeout");if !tsNullish(timeout)&&tsNumber(timeout)>0{state:=tsNodeNetSocketFor(socket);state.timeout=time.Duration(tsNumber(timeout))*time.Millisecond};tsNodeNetSocketConnect(loop,socket,[]tsValue{normalized}...);return socket
}
func tsNodeNetModule(loop *tsLoop)tsValue {
	tsNodeNetSingleton.Do(func(){tsNodeNetSingleton.sockets=map[uintptr]*tsNodeNetSocketState{};tsNodeNetSingleton.servers=map[uintptr]*tsNodeNetServerState{};tsNodeNetSingleton.blocks=map[uintptr]*tsNodeNetBlockListState{};tsNodeNetSingleton.socketAddresses=map[uintptr]*tsNodeNetSocketAddressState{};tsNodeNetSingleton.blockListSerialized=map[uint64]tsValue{};tsNodeNetSingleton.blockListObjectNonce=map[uintptr]uint64{};tsNodeNetSingleton.defaultAutoSelectFamily=true;tsNodeNetSingleton.defaultAutoSelectFamilyAttemptTimeout=500;tsNodeNetReadExecArgv(loop)
		events:=tsNodeEventsModule(loop);eventCtor:=tsNodeNetGet(loop,events,"EventEmitter");eventProto:=tsNodeNetGet(loop,eventCtor,"prototype");duplex:=tsInternalStreamDuplexModule(loop);duplexCtor:=tsNodeNetGet(loop,duplex,"Duplex");duplexProto:=tsNodeNetGet(loop,duplexCtor,"prototype")
		tsNodeNetSingleton.socketProto=tsNewObject();tsNodeNetSingleton.socketProto.prototype=duplexProto;tsNodeNetSingleton.serverProto=tsNewObject();tsNodeNetSingleton.serverProto.prototype=eventProto;tsNodeNetSingleton.blockListProto=tsNewObject();tsNodeNetSingleton.socketAddressProto=tsNewObject()
		tsNodeNetSocketInstall(loop);tsNodeNetServerInstall(loop);tsNodeNetSingleton.socketCtor=tsNodeNetSocketCtor(loop);tsNodeNetSingleton.serverCtor=tsNodeNetServerCtor(loop);tsNodeNetSingleton.blockListCtor=tsNodeNetBlockListCtor(loop);tsNodeNetSingleton.socketAddressCtor=tsNodeNetSocketAddressCtor(loop)
		module:=tsNewObject();module.set("Socket",tsNodeNetSingleton.socketCtor);module.set("Server",tsNodeNetSingleton.serverCtor);module.set("Stream",tsNodeNetSingleton.socketCtor);module.set("BlockList",tsNodeNetSingleton.blockListCtor);module.set("SocketAddress",tsNodeNetSingleton.socketAddressCtor);module.set("createServer",tsNamedNativeMethod("createServer",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeNetServerCreate(loop,args...)}));module.set("createConnection",tsNamedNativeMethod("createConnection",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeNetCreateConnection(loop,args)}));module.set("connect",module.get("createConnection"));module.set("isIP",tsNamedNativeMethod("isIP",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{if tsInternalNetIPv4(tsArg(args,0)){return tsNumberValue(4)};if tsInternalNetIPv6(tsArg(args,0)){return tsNumberValue(6)};return tsNumberValue(0)}));module.set("isIPv4",tsNodeNetGet(loop,tsInternalNetIPModule(),"isIPv4"));module.set("isIPv6",tsNodeNetGet(loop,tsInternalNetIPModule(),"isIPv6"));module.set("_normalizeArgs",tsNamedNativeMethod("_normalizeArgs",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeNetNormalizeArgs(args)}));module.set("_setSimultaneousAccepts",tsNamedNativeMethod("_setSimultaneousAccepts",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{if !tsNodeNetSingleton.simultaneousAcceptsWarned{tsNodeNetSingleton.simultaneousAcceptsWarned=true;process:=tsNodeProcessModule(loop);emit:=tsNodeNetGet(loop,process,"emitWarning");if emit.kind==tsFunctionKind{tsCall(loop,emit,tsNodeNetKey("net._setSimultaneousAccepts() is deprecated and will be removed."),tsNodeNetKey("DeprecationWarning"),tsNodeNetKey("DEP0121"))}};return tsU}));module.set("getDefaultAutoSelectFamily",tsNamedNativeMethod("getDefaultAutoSelectFamily",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsBooleanValue(tsNodeNetSingleton.defaultAutoSelectFamily)}));module.set("setDefaultAutoSelectFamily",tsNamedNativeMethod("setDefaultAutoSelectFamily",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{value:=tsArg(args,0);if value.kind!=tsBooleanKind{tsFSInvalid("autoSelectFamilyDefault","a boolean")};tsNodeNetSingleton.defaultAutoSelectFamily=tsTruthy(value);return tsBooleanValue(tsNodeNetSingleton.defaultAutoSelectFamily)}));module.set("getDefaultAutoSelectFamilyAttemptTimeout",tsNamedNativeMethod("getDefaultAutoSelectFamilyAttemptTimeout",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNumberValue(float64(tsNodeNetSingleton.defaultAutoSelectFamilyAttemptTimeout))}));module.set("setDefaultAutoSelectFamilyAttemptTimeout",tsNamedNativeMethod("setDefaultAutoSelectFamilyAttemptTimeout",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{value:=tsArg(args,0);if !tsIsNumeric(value)||tsNumber(value)!=math.Trunc(tsNumber(value))||tsNumber(value)<1||tsNumber(value)>math.MaxInt32{tsFSInvalid("value","an integer greater than or equal to 1")};number:=int(tsNumber(value));if number<10{number=10};tsNodeNetSingleton.defaultAutoSelectFamilyAttemptTimeout=number;return tsNumberValue(float64(number))}));module.set("default",tsObjectValue(module));tsNodeNetSingleton.module=tsObjectValue(module)
	});return tsNodeNetSingleton.module
}
`
