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
	reading, upgrading bool
	upgradeData []byte
	upgradeContinue func()
	timeoutTask *tsTask
	server *tsNodeNetServerState
	localAddress, localFamily, remoteAddress, remoteFamily string
	localPort, remotePort int
}
type tsNodeNetServerState struct {
	loop *tsLoop
	server tsValue
	listener net.Listener
	connections int
	closing, closed, unref bool
	path string
	options tsValue
}
type tsNodeNetBlockRule struct { network *net.IPNet; ip, end net.IP; family, kind, text string }
var tsNodeNetSingleton struct {
	sync.Once
	module, socketCtor, serverCtor, blockListCtor, socketAddressCtor tsValue
	socketProto, serverProto, blockListProto, socketAddressProto *tsObject
	sockets map[uintptr]*tsNodeNetSocketState
	servers map[uintptr]*tsNodeNetServerState
	blocks map[uintptr][]tsNodeNetBlockRule
	defaultAutoSelectFamily bool
	defaultAutoSelectFamilyAttemptTimeout int
	tlsKeylogPath string
	tlsKeylogConfigured, tlsKeylogWarned bool
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
func tsNodeNetError(err error, operation string) tsValue {
	if err==nil{return tsNull}
	code:="ERR_SOCKET_ERROR";message:=err.Error()
	var op *net.OpError;if errors.As(err,&op){operation=op.Op;if op.Err!=nil{message=op.Err.Error()}}
	var errno syscall.Errno
	if errors.As(err,&errno){code=errno.Error()}else{switch {case strings.Contains(message,"refused"):code="ECONNREFUSED";case strings.Contains(message,"reset"):code="ECONNRESET";case strings.Contains(message,"timed out"):code="ETIMEDOUT";case strings.Contains(message,"address already in use"):code="EADDRINUSE";case strings.Contains(message,"permission denied"):code="EACCES"}}
	error:=&tsRuntimeError{name:"Error",message:operation+" "+code+": "+message,fields:map[string]tsValue{"code":tsNodeNetKey(code),"syscall":tsNodeNetKey(operation)}}
	return tsErrorValue(error)
}
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
	if state==nil||state.timeout<=0||state.closed{return};loop:=state.loop
	if state.timeoutTask!=nil{loop.clearTimeout(tsTaskValue(state.timeoutTask))}
	var task *tsTask
	task=loop.setTimeout(tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{if state.timeoutTask==task&&!state.closed{state.timeoutTask=nil;tsNodeNetEmit(loop,state.stream,"timeout")};return tsU})),tsNumberValue(float64(state.timeout/time.Millisecond)))
	task.unref=state.unref;state.timeoutTask=task
}
func tsNodeNetSetConn(state *tsNodeNetSocketState, conn net.Conn) {
	state.conn=conn;state.connecting=false;state.connected=true
	if tcp,ok:=conn.(*net.TCPConn);ok{_ = tcp.SetNoDelay(state.noDelay);_ = tcp.SetKeepAlive(state.keepAlive);if state.keepAlive&&state.keepAliveDelay>0{_ = tcp.SetKeepAlivePeriod(state.keepAliveDelay)}}
	state.localAddress,state.localPort,state.localFamily=tsNodeNetAddress(conn.LocalAddr());state.remoteAddress,state.remotePort,state.remoteFamily=tsNodeNetAddress(conn.RemoteAddr());tsNodeNetSyncSocket(state)
	state.loop.resourceClosers=append(state.loop.resourceClosers,func(){_ = conn.Close()})
}
type tsNodeNetPrefixConn struct{net.Conn;prefix []byte}
func (conn *tsNodeNetPrefixConn)Read(buffer []byte)(int,error){if len(conn.prefix)>0{n:=copy(buffer,conn.prefix);conn.prefix=conn.prefix[n:];return n,nil};return conn.Conn.Read(buffer)}
func tsNodeNetStartRead(state *tsNodeNetSocketState) {
	if state==nil||state.conn==nil||state.closed||state.reading{return};loop,conn:=state.loop,state.conn;state.reading=true
	var data []byte;var readErr error;var count int
	task:=loop.submit(func()tsResult{buffer:=make([]byte,32*1024);count,readErr=conn.Read(buffer);if count>0{data=append([]byte(nil),buffer[:count]...)};return tsResult{tsU,false}},func(_ tsResult){
		state.reading=false
		if state.upgrading{state.upgrading=false;if count>0{state.upgradeData=append(state.upgradeData,data...)};_ = state.conn.SetReadDeadline(time.Time{});resume:=state.upgradeContinue;state.upgradeContinue=nil;if resume!=nil{resume()};return}
		if state.closed{return};if count>0{state.bytesRead+=int64(count);tsNodeNetSyncSocket(state);tsNodeNetResetTimeout(state);if readable:=tsNodeReadableStateFor(state.stream);readable!=nil{tsNodeReadablePush(loop,readable,tsNodeBuffer(data),tsU,false)}}
		if readErr!=nil{if errors.Is(readErr,io.EOF){if readable:=tsNodeReadableStateFor(state.stream);readable!=nil{tsNodeReadablePush(loop,readable,tsNull,tsU,false)};if !tsTruthy(tsNodeNetGet(loop,state.stream,"allowHalfOpen")){if writable:=tsNodeWritableStateFor(state.stream);writable!=nil{tsNodeWritableEnd(loop,writable,nil)}};return};if !state.closed{tsNodeNetEmit(loop,state.stream,"error",tsNodeNetError(readErr,"read"));destroy:=tsNodeNetGet(loop,state.stream,"destroy");if destroy.kind==tsFunctionKind{tsCallReceiver(loop,destroy,state.stream)}};return};tsNodeNetStartRead(state)
	});task.unref=state.unref
}
func tsNodeNetSocketWrite(loop *tsLoop, receiver tsValue, args ...tsValue) tsValue {
	state:=tsNodeNetSocketFor(receiver);callback:=tsArg(args,2);if state==nil||state.conn==nil||state.closed{err:=tsFSException("Error","ERR_SOCKET_CLOSED","Socket is not connected");if callback.kind==tsFunctionKind{tsCall(loop,callback,err)};return tsU}
	chunk:=tsArg(args,0);var data []byte;if chunk.kind==tsStringKind{data=[]byte(tsText(chunk))}else{data=append([]byte(nil),tsFSBytes(chunk)...)};conn:=state.conn
	task:=loop.submit(func()tsResult{_,err:=conn.Write(data);if err!=nil{return tsResult{tsNodeNetError(err,"write"),true}};return tsResult{tsU,false}},func(result tsResult){if result.rejected{if callback.kind==tsFunctionKind{tsCall(loop,callback,result.value)};if !state.closed{destroy:=tsNodeNetGet(loop,receiver,"destroy");if destroy.kind==tsFunctionKind{tsCallReceiver(loop,destroy,receiver,result.value)}};return};state.bytesWritten+=int64(len(data));tsNodeNetSyncSocket(state);tsNodeNetResetTimeout(state);if callback.kind==tsFunctionKind{tsCall(loop,callback,tsU)}});task.unref=state.unref;return tsU
}
func tsNodeNetSocketFinal(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue {
	state:=tsNodeNetSocketFor(receiver);callback:=tsArg(args,0);if state!=nil&&state.conn!=nil{if closer,ok:=state.conn.(interface{CloseWrite()error});ok{if err:=closer.CloseWrite();err!=nil&&callback.kind==tsFunctionKind{tsCall(loop,callback,tsNodeNetError(err,"shutdown"));return tsU}}};if callback.kind==tsFunctionKind{tsCall(loop,callback,tsU)};return tsU
}
func tsNodeNetSocketDestroy(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue {
	state:=tsNodeNetSocketFor(receiver);err,callback:=tsArg(args,0),tsArg(args,1);if state!=nil&&!state.closed{state.closed=true;state.connecting=false;if state.timeoutTask!=nil{loop.clearTimeout(tsTaskValue(state.timeoutTask));state.timeoutTask=nil};if state.conn!=nil{_ = state.conn.Close()};tsNodeNetSyncSocket(state);if state.server!=nil{state.server.connections--;if state.server.connections<0{state.server.connections=0};tsNodeNetSet(state.server.loop,state.server.server,"_connections",tsNumberValue(float64(state.server.connections)));tsNodeNetServerMaybeClose(state.server)}};if callback.kind==tsFunctionKind{tsCall(loop,callback,err)};return tsU
}
func tsNodeNetDialSocket(loop *tsLoop,receiver tsValue,state *tsNodeNetSocketState,network,address string,local tsValue,localPort int){
	var dialed net.Conn;task:=loop.submit(func()tsResult{dialer:=net.Dialer{};if network=="tcp"&&!tsNullish(local){dialer.LocalAddr=&net.TCPAddr{IP:net.ParseIP(tsText(local)),Port:localPort}};conn,err:=dialer.Dial(network,address);if err!=nil{return tsResult{tsNodeNetError(err,"connect"),true}};dialed=conn;return tsResult{tsU,false}},func(result tsResult){if result.rejected{state.connecting=false;tsNodeNetSyncSocket(state);tsNodeNetEmit(loop,receiver,"error",result.value);return};tsNodeNetSetConn(state,dialed);tsNodeNetStartRead(state);tsNodeNetResetTimeout(state);tsNodeNetEmit(loop,receiver,"connect")});task.unref=state.unref
}
func tsNodeNetSocketConnect(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue {
	state:=tsNodeNetSocketFor(receiver);if state==nil{return receiver};if state.connected||state.connecting{panic(tsThrown{tsFSException("Error","ERR_SOCKET_DGRAM_IS_CONNECTED","Socket is already connecting")})}
	arg:=tsArg(args,0);callback:=tsArg(args,1);options:=arg
	if arg.kind==tsArrayKind{items:=(*tsArray)(arg.ref);options=tsArg(items.values,0);if tsIsUndefined(callback){callback=tsArg(items.values,1)}}else if arg.kind!=tsObjectKind&&arg.kind!=tsInstanceKind&&!tsNullish(arg){normalized:=tsNodeNetNormalizeArgs(args);items:=(*tsArray)(normalized.ref);options=tsArg(items.values,0);callback=tsArg(items.values,1)}
	if callback.kind==tsFunctionKind{tsCallReceiver(loop,tsNodeNetGet(loop,receiver,"once"),receiver,tsNodeNetKey("connect"),callback)}
	path:=tsNodeNetGet(loop,options,"path");network:="tcp";host:=tsText(tsNodeNetGet(loop,options,"host"));port:=tsNodeNetGet(loop,options,"port")
	if !tsNullish(path)&&tsText(path)!=""{network="unix";host=tsText(path)}else{if host==""||host=="undefined"{host="localhost"};if tsIsUndefined(port){port=tsNodeNetGet(loop,options,"port")}}
	portText:="0";if !tsNullish(port){portText=tsText(port)};local:=tsNodeNetGet(loop,options,"localAddress");localPort:=int(tsNumber(tsNodeNetGet(loop,options,"localPort")));state.connecting=true;tsNodeNetSyncSocket(state)
	if network=="unix"{tsNodeNetDialSocket(loop,receiver,state,network,host,local,localPort);return receiver}
	if net.ParseIP(strings.Split(host,"%")[0])!=nil{tsNodeNetDialSocket(loop,receiver,state,network,net.JoinHostPort(host,portText),local,localPort);return receiver}
	lookup:=tsNodeNetGet(loop,options,"lookup");if lookup.kind!=tsFunctionKind{lookup=tsNodeNetGet(loop,tsNodeDNSModule(loop,false),"lookup")};lookupOptions:=tsNewObject();family:=tsNodeNetGet(loop,options,"family");if family.kind==tsStringKind{switch tsText(family){case "IPv4":family=tsNumberValue(4);case "IPv6":family=tsNumberValue(6)}};if !tsIsUndefined(family){lookupOptions.set("family",family)};hints:=tsNodeNetGet(loop,options,"hints");if !tsIsUndefined(hints){lookupOptions.set("hints",hints)}
	resolved:=tsFunctionValue(tsFunc(func(loop *tsLoop,values ...tsValue)tsValue{err,addressValue,familyValue:=tsArg(values,0),tsArg(values,1),tsArg(values,2);if !tsNullish(err){state.connecting=false;tsNodeNetSyncSocket(state);tsNodeNetEmit(loop,receiver,"lookup",err,tsU,tsU,tsNodeNetKey(host));tsNodeNetEmit(loop,receiver,"error",err);return tsU};if addressValue.kind!=tsStringKind||(!tsInternalNetIPv4(addressValue)&&!tsInternalNetIPv6(addressValue)){state.connecting=false;tsNodeNetSyncSocket(state);bad:=tsFSException("TypeError","ERR_INVALID_IP_ADDRESS","Invalid IP address");tsNodeNetEmit(loop,receiver,"error",bad);return tsU};tsNodeNetEmit(loop,receiver,"lookup",tsNull,addressValue,familyValue,tsNodeNetKey(host));tsNodeNetDialSocket(loop,receiver,state,network,net.JoinHostPort(tsText(addressValue),portText),local,localPort);return tsU}))
	tsCall(loop,lookup,tsNodeNetKey(host),tsObjectValue(lookupOptions),resolved);return receiver
}
func tsNodeNetSocketCreate(loop *tsLoop,options tsValue)tsValue {
	tsNodeDuplexInit(loop);if tsNullish(options){options=tsObjectValue(tsNewObject())};stream:=tsNodeDuplexCreate(loop,tsU,options);object:=(*tsObject)(stream.ref);object.prototype=tsObjectValue(tsNodeNetSingleton.socketProto)
	state:=&tsNodeNetSocketState{loop:loop,stream:stream,connecting:false,unref:false,remoteFamily:"",noDelay:tsTruthy(tsNodeNetGet(loop,options,"noDelay")),keepAlive:tsTruthy(tsNodeNetGet(loop,options,"keepAlive")),typeOfService:0};if delay:=tsNodeNetGet(loop,options,"keepAliveInitialDelay");tsIsNumeric(delay)&&tsNumber(delay)>0{state.keepAliveDelay=time.Duration(tsNumber(delay))*time.Millisecond};tsNodeNetSingleton.sockets[uintptr(stream.ref)]=state;tsNodeEventCreate(loop,stream,tsU);tsNodeNetSyncSocket(state)
	if timeout:=tsNodeNetGet(loop,options,"timeout");!tsNullish(timeout)&&tsNumber(timeout)>0{state.timeout=time.Duration(tsNumber(timeout))*time.Millisecond}
	return stream
}
func tsNodeNetSocketCtor(loop *tsLoop)tsValue {
	constructor:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeNetSocketCreate(loop,tsArg(args,0))});constructor.name="Socket";constructor.length=1;constructor.constructible=true;constructor.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeNetSocketCreate(loop,tsArg(args,0))};constructor.receiverCall=func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeNetSocketCreate(loop,tsArg(args,0))};constructor.properties=tsNewObject();constructor.properties.set("prototype",tsObjectValue(tsNodeNetSingleton.socketProto));constructor.properties.descriptors=map[string]*tsDescriptor{"prototype":{writable:false}};return tsFunctionValue(constructor)
}
func tsNodeNetSocketInstall(loop *tsLoop) {
	p:=tsNodeNetSingleton.socketProto
	methods:=map[string]tsValue{
		"connect":tsNamedNativeMethod("connect",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeNetSocketConnect(loop,receiver,args...)}),
		"_write":tsNamedNativeMethod("_write",3,tsNodeNetSocketWrite),
		"_final":tsNamedNativeMethod("_final",1,tsNodeNetSocketFinal),
		"_destroy":tsNamedNativeMethod("_destroy",2,tsNodeNetSocketDestroy),
		"address":tsNamedNativeMethod("address",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeNetSocketFor(receiver);if state==nil||!state.connected{return tsNull};out:=tsNewObject();out.set("address",tsStringReference(tsStringUTF8(state.localAddress)));out.set("family",tsStringReference(tsStringUTF8(state.localFamily)));out.set("port",tsNumberValue(float64(state.localPort)));return tsObjectValue(out)}),
		"setNoDelay":tsNamedNativeMethod("setNoDelay",1,func(_ *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeNetSocketFor(receiver);if state!=nil{state.noDelay=!tsStrictEqual(tsArg(args,0),tsBooleanValue(false));if tcp,ok:=state.conn.(*net.TCPConn);ok{_ = tcp.SetNoDelay(state.noDelay)}};return receiver}),
		"setKeepAlive":tsNamedNativeMethod("setKeepAlive",2,func(_ *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeNetSocketFor(receiver);if state!=nil{state.keepAlive=!tsStrictEqual(tsArg(args,0),tsBooleanValue(false));delay:=tsArg(args,1);if !tsIsUndefined(delay)&&tsIsNumeric(delay)&&tsNumber(delay)>0{state.keepAliveDelay=time.Duration(tsNumber(delay))*time.Millisecond};if tcp,ok:=state.conn.(*net.TCPConn);ok{_ = tcp.SetKeepAlive(state.keepAlive);if state.keepAlive&&state.keepAliveDelay>0{_ = tcp.SetKeepAlivePeriod(state.keepAliveDelay)}}};return receiver}),
		"setTypeOfService":tsNamedNativeMethod("setTypeOfService",1,func(_ *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeNetSocketFor(receiver);value:=tsArg(args,0);if !tsIsNumeric(value)||tsNumber(value)!=math.Trunc(tsNumber(value))||tsNumber(value)<0||tsNumber(value)>255{tsFSInvalid("tos","an integer between 0 and 255")};if state!=nil{state.typeOfService=int(tsNumber(value))};return receiver}),
		"getTypeOfService":tsNamedNativeMethod("getTypeOfService",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{if state:=tsNodeNetSocketFor(receiver);state!=nil{return tsNumberValue(float64(state.typeOfService))};return tsNumberValue(0)}),
		"setTimeout":tsNamedNativeMethod("setTimeout",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeNetSocketFor(receiver);if state!=nil{timeout:=tsArg(args,0);if !tsNullish(timeout){if !tsIsNumeric(timeout)||tsNumber(timeout)<0{tsFSInvalid("timeout","a non-negative number")};state.timeout=time.Duration(tsNumber(timeout))*time.Millisecond;tsNodeNetSet(loop,receiver,"timeout",timeout);tsNodeNetResetTimeout(state)};if listener:=tsArg(args,1);listener.kind==tsFunctionKind{tsCallReceiver(loop,tsNodeNetGet(loop,receiver,"once"),receiver,tsNodeNetKey("timeout"),listener)}};return receiver}),
		"ref":tsNamedNativeMethod("ref",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{if state:=tsNodeNetSocketFor(receiver);state!=nil{state.unref=false};return receiver}),
		"unref":tsNamedNativeMethod("unref",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{if state:=tsNodeNetSocketFor(receiver);state!=nil{state.unref=true};return receiver}),
		"resetAndDestroy":tsNamedNativeMethod("resetAndDestroy",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{if state:=tsNodeNetSocketFor(receiver);state!=nil&&state.conn!=nil{_ = state.conn.Close()};return tsCallReceiver(loop,tsNodeNetGet(loop,receiver,"destroy"),receiver)}),
		"_reset":tsNamedNativeMethod("_reset",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeNetSocketFor(receiver);if state!=nil&&state.conn!=nil{_ = state.conn.Close();state.conn=nil};return tsCallReceiver(loop,tsNodeNetGet(loop,receiver,"destroy"),receiver)}),
		"_getsockname":tsNamedNativeMethod("_getsockname",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeNetSocketFor(receiver);if state==nil||state.localAddress==""{return tsNull};out:=tsNewObject();out.set("address",tsNodeNetKey(state.localAddress));out.set("port",tsNumberValue(float64(state.localPort)));out.set("family",tsNodeNetKey(state.localFamily));return tsObjectValue(out)}),
		"_getpeername":tsNamedNativeMethod("_getpeername",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeNetSocketFor(receiver);if state==nil||state.remoteAddress==""{return tsNull};out:=tsNewObject();out.set("address",tsNodeNetKey(state.remoteAddress));out.set("port",tsNumberValue(float64(state.remotePort)));out.set("family",tsNodeNetKey(state.remoteFamily));return tsObjectValue(out)}),
		"destroySoon":tsNamedNativeMethod("destroySoon",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{end:=tsNodeNetGet(loop,receiver,"end");if end.kind==tsFunctionKind{tsCallReceiver(loop,end,receiver)};return receiver}),
		"_unrefTimer":tsNamedNativeMethod("_unrefTimer",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{if state:=tsNodeNetSocketFor(receiver);state!=nil&&state.timeoutTask!=nil{state.timeoutTask.unref=state.unref};return tsU}),
	}
	for name,method:=range methods{p.set(name,method)}
	if p.descriptors==nil{p.descriptors=map[string]*tsDescriptor{}}
	for _,name:=range []string{"bytesRead","bytesWritten","_bytesDispatched","bufferSize","connecting","_connecting","pending","readyState","timeout","localAddress","localPort","localFamily","remoteAddress","remotePort","remoteFamily"}{key:=name;p.descriptors[name]=&tsDescriptor{accessor:true,configurable:true,getter:tsNamedNativeMethod("get "+key,0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeNetSocketFor(receiver);if state==nil{return tsU};switch key{case "bytesRead":return tsNumberValue(float64(state.bytesRead));case "bytesWritten":value:=float64(state.bytesWritten);if writable:=tsNodeWritableStateFor(receiver);writable!=nil{value+=float64(writable.length)};return tsNumberValue(value);case "_bytesDispatched":return tsNumberValue(float64(state.bytesWritten));case "bufferSize":if state.conn==nil||state.closed{return tsU};if writable:=tsNodeWritableStateFor(receiver);writable!=nil{return tsNumberValue(float64(writable.length))};return tsNumberValue(0);case "connecting","_connecting":return tsBooleanValue(state.connecting);case "pending":return tsBooleanValue(state.connecting||state.conn==nil);case "timeout":return tsNumberValue(float64(state.timeout/time.Millisecond));case "readyState":if state.connecting{return tsStringReference(tsStringUTF8("opening"))};readable,writable:=tsTruthy(tsNodeNetGet(state.loop,receiver,"readable")),tsTruthy(tsNodeNetGet(state.loop,receiver,"writable"));if readable&&writable{return tsStringReference(tsStringUTF8("open"))};if readable{return tsStringReference(tsStringUTF8("readOnly"))};if writable{return tsStringReference(tsStringUTF8("writeOnly"))};return tsStringReference(tsStringUTF8("closed"));case "localAddress":if state.localAddress==""{return tsU};return tsStringReference(tsStringUTF8(state.localAddress));case "localPort":if state.localAddress==""{return tsU};return tsNumberValue(float64(state.localPort));case "localFamily":if state.localAddress==""{return tsU};return tsStringReference(tsStringUTF8(state.localFamily));case "remoteAddress":if state.remoteAddress==""{return tsU};return tsStringReference(tsStringUTF8(state.remoteAddress));case "remotePort":if state.remoteAddress==""{return tsU};return tsNumberValue(float64(state.remotePort));default:if state.remoteFamily==""{return tsU};return tsStringReference(tsStringUTF8(state.remoteFamily))}})}}
}
func tsNodeNetMakeAddr(address string,port int)tsValue {
	out:=tsNewObject();host,family,_:=tsNodeNetAddress(&net.TCPAddr{IP:net.ParseIP(address),Port:port});if host==""{host=address};out.set("address",tsStringReference(tsStringUTF8(host)));out.set("port",tsNumberValue(float64(port)));out.set("family",tsStringReference(tsStringUTF8(family)));return tsObjectValue(out)
}
func tsNodeNetServerMaybeClose(state *tsNodeNetServerState) {
	if state==nil||!state.closing||state.closed||state.connections>0{return};state.closed=true
	if state.listener!=nil{_ = state.listener.Close();state.listener=nil}
	loop:=state.loop;loop.ready=append(loop.ready,func(){tsNodeNetEmit(loop,state.server,"close")})
}
func tsNodeNetServerAccept(state *tsNodeNetServerState) {
	if state==nil||state.listener==nil||state.closing{return};listener:=state.listener;var conn net.Conn;var acceptErr error
	task:=state.loop.submit(func()tsResult{conn,acceptErr=listener.Accept();if acceptErr!=nil{return tsResult{tsNodeNetError(acceptErr,"accept"),true}};return tsResult{tsU,false}},func(result tsResult){
		if state.closing||state.listener==nil{if conn!=nil{_ = conn.Close()};tsNodeNetServerMaybeClose(state);return}
		if result.rejected{state.closing=true;state.listener=nil;tsNodeNetEmit(state.loop,state.server,"error",result.value);tsNodeNetServerMaybeClose(state);return}
		options:=tsNodeNetGet(state.loop,state.server,"_socketOptions");socket:=tsNodeNetSocketCreate(state.loop,options);child:=tsNodeNetSocketFor(socket);child.server=state;tsNodeNetSetConn(child,conn)
		remote:=child.remoteAddress;family:=child.remoteFamily;if family=="IPv4"{family="ipv4"}else if family=="IPv6"{family="ipv6"};blockList:=tsNodeNetGet(state.loop,state.server,"blockList")
		if blockList.kind==tsObjectKind&&tsNodeNetBlockListCheck(state.loop,blockList,tsNodeNetKey(remote),tsNodeNetKey(family)){data:=tsNewObject();data.set("localAddress",tsNodeNetKey(child.localAddress));data.set("localPort",tsNumberValue(float64(child.localPort)));data.set("localFamily",tsNodeNetKey(child.localFamily));data.set("remoteAddress",tsNodeNetKey(child.remoteAddress));data.set("remotePort",tsNumberValue(float64(child.remotePort)));data.set("remoteFamily",tsNodeNetKey(child.remoteFamily));_ = conn.Close();child.closed=true;tsNodeNetEmit(state.loop,state.server,"drop",tsObjectValue(data));tsNodeNetServerAccept(state);return}
		maxConnections:=tsNodeNetGet(state.loop,state.server,"maxConnections");if tsIsNumeric(maxConnections)&&tsNumber(maxConnections)>0&&float64(state.connections)>=tsNumber(maxConnections){data:=tsNewObject();data.set("localAddress",tsNodeNetKey(child.localAddress));data.set("localPort",tsNumberValue(float64(child.localPort)));data.set("localFamily",tsNodeNetKey(child.localFamily));data.set("remoteAddress",tsNodeNetKey(child.remoteAddress));data.set("remotePort",tsNumberValue(float64(child.remotePort)));data.set("remoteFamily",tsNodeNetKey(child.remoteFamily));_ = conn.Close();child.closed=true;tsNodeNetEmit(state.loop,state.server,"drop",tsObjectValue(data));tsNodeNetServerAccept(state);return}
		state.connections++;tsNodeNetSet(state.loop,state.server,"_connections",tsNumberValue(float64(state.connections)));tsNodeNetSyncSocket(child);tsNodeNetEmit(state.loop,state.server,"connection",socket);tsNodeNetStartRead(child);tsNodeNetServerAccept(state)
	});task.unref=state.unref
}
func tsNodeNetListen(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue {
	state:=tsNodeNetServerFor(receiver);if state==nil{return receiver};if state.listener!=nil{panic(tsThrown{tsFSException("Error","ERR_SERVER_ALREADY_LISTEN","Listen method has been called more than once without closing.")})}
	arg:=tsArg(args,0);options:=arg;callback:=tsArg(args,1);host:="";port:=0;network:="tcp";address:=""
	if arg.kind==tsObjectKind||arg.kind==tsInstanceKind{options=arg;path:=tsNodeNetGet(loop,arg,"path");if !tsNullish(path)&&tsText(path)!=""{network="unix";address=tsText(path)}else{port=int(tsNumber(tsNodeNetGet(loop,arg,"port")));host=tsText(tsNodeNetGet(loop,arg,"host"));if host=="undefined"{host=""};address=net.JoinHostPort(host,strconv.Itoa(port))};if listener:=tsNodeNetGet(loop,arg,"callback");listener.kind==tsFunctionKind{callback=listener}}
	if arg.kind==tsStringKind {text:=tsText(arg);if parsed,err:=strconv.Atoi(text);err==nil{port=parsed;host=tsText(tsArg(args,1));if host=="undefined"||host==""{host=""};address=net.JoinHostPort(host,strconv.Itoa(port))}else{network="unix";address=text}}
	if arg.kind==tsNumberKind||tsIsSigned(arg)||tsIsUnsigned(arg){port=int(tsNumber(arg));host=tsText(tsArg(args,1));if host=="undefined"||host==""{host=""};address=net.JoinHostPort(host,strconv.Itoa(port));if tsArg(args,1).kind==tsFunctionKind{callback=tsArg(args,1)}}
	if arg.kind==tsFunctionKind{callback=arg;address=":0"}
	if address==""{address=":0"};if callback.kind==tsFunctionKind{tsCallReceiver(loop,tsNodeNetGet(loop,receiver,"once"),receiver,tsNodeNetKey("listening"),callback)}
	listener,err:=net.Listen(network,address);if err!=nil{tsNodeNetEmit(loop,receiver,"error",tsNodeNetError(err,"listen"));return receiver};state.listener=listener;state.closing=false;state.closed=false;state.path=address;tsNodeNetSet(loop,receiver,"_socketOptions",state.options)
	if network=="unix"{tsNodeNetSet(loop,receiver,"_address",tsNodeNetKey(listener.Addr().String()))}else if host,listenPort,family:=tsNodeNetAddress(listener.Addr());host!=""{tsNodeNetSet(loop,receiver,"_address",tsNodeNetMakeAddr(host,listenPort));tsNodeNetSet(loop,receiver,"_connectionKey",tsNodeNetKey(family+":"+host+":"+strconv.Itoa(listenPort)))}
	state.loop.resourceClosers=append(state.loop.resourceClosers,func(){_ = listener.Close();if network=="unix"{_ = os.Remove(address)}});tsNodeNetEmit(loop,receiver,"listening");tsNodeNetServerAccept(state);return receiver
}
func tsNodeNetServerClose(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue {
	state:=tsNodeNetServerFor(receiver);callback:=tsArg(args,0)
	if state==nil||state.listener==nil{if callback.kind==tsFunctionKind{failure:=tsFSException("Error","ERR_SERVER_NOT_RUNNING","Server is not running.");wrapped:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{tsCall(loop,callback,failure);return tsU}));tsCallReceiver(loop,tsNodeNetGet(loop,receiver,"once"),receiver,tsNodeNetKey("close"),wrapped)};if state!=nil{if state.closed&&state.connections==0{loop.ready=append(loop.ready,func(){tsNodeNetEmit(loop,receiver,"close")})}else{state.closing=true;tsNodeNetServerMaybeClose(state)}}else{loop.ready=append(loop.ready,func(){tsNodeNetEmit(loop,receiver,"close")})};return receiver}
	if callback.kind==tsFunctionKind{wrapped:=tsFunctionValue(tsFunc(func(loop *tsLoop,values ...tsValue)tsValue{tsCall(loop,callback,tsNull);return tsU}));tsCallReceiver(loop,tsNodeNetGet(loop,receiver,"once"),receiver,tsNodeNetKey("close"),wrapped)}
	state.closing=true;_ = state.listener.Close();state.listener=nil;tsNodeNetServerMaybeClose(state);return receiver
}
func tsNodeNetServerCreate(loop *tsLoop,args ...tsValue)tsValue {
	options,listener:=tsArg(args,0),tsArg(args,1);if options.kind==tsFunctionKind{listener=options;options=tsObjectValue(tsNewObject())};if tsNullish(options){options=tsObjectValue(tsNewObject())};object:=tsNewObject();object.prototype=tsObjectValue(tsNodeNetSingleton.serverProto);server:=tsNodeEventCreate(loop,tsObjectValue(object),tsU);state:=&tsNodeNetServerState{loop:loop,server:server,options:options};tsNodeNetSingleton.servers[uintptr(server.ref)]=state
	object.set("_connections",tsNumberValue(0));object.set("_socketOptions",options);object.set("allowHalfOpen",tsBooleanValue(tsTruthy(tsNodeNetGet(loop,options,"allowHalfOpen"))));object.set("pauseOnConnect",tsBooleanValue(tsTruthy(tsNodeNetGet(loop,options,"pauseOnConnect"))));object.set("noDelay",tsBooleanValue(tsTruthy(tsNodeNetGet(loop,options,"noDelay"))));object.set("keepAlive",tsBooleanValue(tsTruthy(tsNodeNetGet(loop,options,"keepAlive"))));object.set("keepAliveInitialDelay",tsNodeNetGet(loop,options,"keepAliveInitialDelay"));object.set("maxConnections",tsNodeNetGet(loop,options,"maxConnections"));if listener.kind==tsFunctionKind{tsCallReceiver(loop,tsNodeNetGet(loop,server,"on"),server,tsNodeNetKey("connection"),listener)};return server
}
func tsNodeNetServerCtor(loop *tsLoop)tsValue {constructor:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeNetServerCreate(loop,args...)});constructor.name="Server";constructor.length=2;constructor.constructible=true;constructor.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeNetServerCreate(loop,args...)};constructor.receiverCall=func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeNetServerCreate(loop,args...)};constructor.properties=tsNewObject();constructor.properties.set("prototype",tsObjectValue(tsNodeNetSingleton.serverProto));constructor.properties.descriptors=map[string]*tsDescriptor{"prototype":{writable:false}};return tsFunctionValue(constructor)}
func tsNodeNetServerInstall(loop *tsLoop) {
	p:=tsNodeNetSingleton.serverProto;p.set("listen",tsNamedNativeMethod("listen",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeNetListen(loop,receiver,args...)}));p.set("close",tsNamedNativeMethod("close",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeNetServerClose(loop,receiver,args...)}));p.set("address",tsNamedNativeMethod("address",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeNetServerFor(receiver);if state==nil||state.listener==nil{return tsNull};return tsNodeNetGet(loop,receiver,"_address")}));p.set("getConnections",tsNamedNativeMethod("getConnections",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeNetServerFor(receiver);count:=0;if state!=nil&&state.listener!=nil{count=state.connections};if callback:=tsArg(args,0);callback.kind==tsFunctionKind{tsCall(loop,callback,tsNull,tsNumberValue(float64(count)))};return receiver}));for name,unref:=range map[string]bool{"ref":false,"unref":true}{method:=name;p.set(method,tsNamedNativeMethod(method,0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{if state:=tsNodeNetServerFor(receiver);state!=nil{state.unref=unref};return receiver}))}
	p.set(tsPropertyKey(tsWellKnownSymbol("asyncDispose")),tsNamedNativeMethod("[Symbol.asyncDispose]",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeNetServerFor(receiver);if state==nil||state.listener==nil{return tsPromiseValue(loop.resolved(tsU,false))};promise:=loop.promise();callback:=tsFunctionValue(tsFunc(func(_ *tsLoop,args ...tsValue)tsValue{err:=tsArg(args,0);if tsTruthy(err){promise.settle(tsResult{err,true})}else{promise.settle(tsResult{tsU,false})};return tsU}));tsNodeNetServerClose(loop,receiver,[]tsValue{callback}...);return tsPromiseValue(promise)}))
	if p.descriptors==nil{p.descriptors=map[string]*tsDescriptor{}};p.descriptors["listening"]=&tsDescriptor{accessor:true,configurable:true,getter:tsNamedNativeMethod("get listening",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeNetServerFor(receiver);return tsBooleanValue(state!=nil&&state.listener!=nil)})};p.descriptors["connections"]=&tsDescriptor{accessor:true,configurable:true,getter:tsNamedNativeMethod("get connections",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeNetServerFor(receiver);if state==nil{return tsNull};return tsNumberValue(float64(state.connections))})}
}
func tsNodeNetBlockListAdd(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue {
	key:=tsNodeNetPointer(receiver);method:=tsText(tsArg(args,0));firstValue,secondValue,familyValue:=tsArg(args,1),tsArg(args,2),tsArg(args,3)
	rule:=tsNodeNetBlockRule{kind:method}
	switch method {
	case "addAddress":
		rule.ip,rule.family=tsNodeNetBlockIP(loop,firstValue,"address",familyValue)
		rule.text=fmt.Sprintf("Address: %s %s",tsNodeNetBlockFamilyLabel(rule.family),rule.ip.String())
	case "addRange":
		start,startFamily:=tsNodeNetBlockIP(loop,firstValue,"start",familyValue);end,endFamily:=tsNodeNetBlockIP(loop,secondValue,"end",familyValue)
		if startFamily!=endFamily||len(start)!=len(end)||tsNodeNetIPCompare(start,end)>0{panic(tsThrown{tsFSException("RangeError","ERR_INVALID_ARG_VALUE","The argument 'start' must come before end")})}
		rule.ip,rule.end,rule.family=start,end,startFamily
		rule.text=fmt.Sprintf("Range: %s %s-%s",tsNodeNetBlockFamilyLabel(rule.family),start.String(),end.String())
	case "addSubnet":
		ip,family:=tsNodeNetBlockIP(loop,firstValue,"network",familyValue);prefixValue:=secondValue
		bits:=128;if family=="ipv4"{bits=32};if !tsIsNumeric(prefixValue)||tsNumber(prefixValue)!=math.Trunc(tsNumber(prefixValue))||tsNumber(prefixValue)<0||tsNumber(prefixValue)>float64(bits){tsFSInvalid("prefix",fmt.Sprintf("an integer between 0 and %d",bits))}
		prefix:=int(tsNumber(prefixValue));rule.family,rule.ip=family,ip
		mask:=net.CIDRMask(prefix,bits);rule.network=&net.IPNet{IP:ip.Mask(mask),Mask:mask}
		rule.text=fmt.Sprintf("Subnet: %s %s/%d",tsNodeNetBlockFamilyLabel(family),ip.String(),prefix)
	default:return tsU
	}
	// Bun inserts the newest entry at index zero; the public rules getter uses
	// this same order.
	tsNodeNetSingleton.blocks[key]=append([]tsNodeNetBlockRule{rule},tsNodeNetSingleton.blocks[key]...);return tsU
}
func tsNodeNetBlockFamilyLabel(family string)string{if family=="ipv6"{return "IPv6"};return "IPv4"}
func tsNodeNetIPCompare(a,b net.IP)int{if len(a)!=len(b){if len(a)<len(b){return -1};return 1};for i:=range a{if a[i]<b[i]{return -1};if a[i]>b[i]{return 1}};return 0}
func tsNodeNetBlockIP(loop *tsLoop,value tsValue,argName string,familyValue tsValue)(net.IP,string){
	family:="ipv4"
	if value.kind==tsObjectKind||value.kind==tsInstanceKind{
		address:=tsNodeNetGet(loop,value,"address");providedFamily:=tsNodeNetGet(loop,value,"family")
		if address.kind!=tsStringKind||providedFamily.kind!=tsStringKind{tsFSInvalid(argName,"a SocketAddress or string")}
		family=strings.ToLower(tsText(providedFamily));ip:=net.ParseIP(tsText(address));return tsNodeNetNormalizeBlockIP(ip,family,argName),family
	}
	if value.kind!=tsStringKind{tsFSInvalid(argName,"a string")}
	if !tsIsUndefined(familyValue){if familyValue.kind!=tsStringKind{tsFSInvalid("family","a string")};family=strings.ToLower(tsText(familyValue))}
	ip:=net.ParseIP(tsText(value));return tsNodeNetNormalizeBlockIP(ip,family,argName),family
}
func tsNodeNetNormalizeBlockIP(ip net.IP,family,argName string)net.IP{
	switch family{
	case "ipv4":if ip!=nil&&ip.To4()!=nil{return append(net.IP(nil),ip.To4()...)}
	case "ipv6":if ip!=nil&&ip.To4()==nil&&ip.To16()!=nil{return append(net.IP(nil),ip.To16()...)}
	default:tsFSInvalid("family","'ipv4' or 'ipv6'")
	}
	panic(tsThrown{tsFSException("TypeError","ERR_INVALID_IP_ADDRESS","Invalid IP address for "+argName)})
}
func tsNodeNetBlockListCheck(loop *tsLoop,receiver,address,family tsValue)bool {
	fam:="ipv4";var ip net.IP
	if address.kind==tsObjectKind||address.kind==tsInstanceKind{addressValue:=tsNodeNetGet(loop,address,"address");familyValue:=tsNodeNetGet(loop,address,"family");if addressValue.kind!=tsStringKind||familyValue.kind!=tsStringKind{return false};fam=strings.ToLower(tsText(familyValue));ip=tsNodeNetNormalizeBlockIPForCheck(net.ParseIP(tsText(addressValue)),fam)}else{if address.kind!=tsStringKind{tsFSInvalid("address","a string or SocketAddress")};if !tsIsUndefined(family){if family.kind!=tsStringKind{tsFSInvalid("family","a string")};fam=strings.ToLower(tsText(family))};if fam!="ipv4"&&fam!="ipv6"{tsFSInvalid("family","'ipv4' or 'ipv6'")};ip=tsNodeNetNormalizeBlockIPForCheck(net.ParseIP(tsText(address)),fam)};if ip==nil{return false}
	for _,rule:=range tsNodeNetSingleton.blocks[tsNodeNetPointer(receiver)]{
		if rule.kind=="addAddress"&&rule.family==fam&&tsNodeNetIPCompare(rule.ip,ip)==0{return true}
		if rule.kind=="addRange"&&rule.family==fam&&tsNodeNetIPCompare(ip,rule.ip)>=0&&tsNodeNetIPCompare(ip,rule.end)<=0{return true}
		if rule.kind=="addSubnet"&&rule.network!=nil{candidate:=ip;if rule.family=="ipv6"&&fam=="ipv4"{candidate=net.ParseIP("::ffff:"+ip.String())};if rule.network.Contains(candidate){return true}}
	};return false
}
func tsNodeNetNormalizeBlockIPForCheck(ip net.IP,family string)net.IP{if ip==nil{return nil};if family=="ipv4"{if v4:=ip.To4();v4!=nil{return append(net.IP(nil),v4...) };return nil};if family=="ipv6"&&ip.To4()==nil&&ip.To16()!=nil{return append(net.IP(nil),ip.To16()...)};return nil}
func tsNodeNetBlockListCtor(loop *tsLoop)tsValue {
	p:=tsNodeNetSingleton.blockListProto;p.set("addAddress",tsNamedNativeMethod("addAddress",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeNetBlockListAdd(loop,receiver,append([]tsValue{tsNodeNetKey("addAddress")},args...)...)}));p.set("addRange",tsNamedNativeMethod("addRange",3,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeNetBlockListAdd(loop,receiver,append([]tsValue{tsNodeNetKey("addRange")},args...)...)}));p.set("addSubnet",tsNamedNativeMethod("addSubnet",3,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeNetBlockListAdd(loop,receiver,append([]tsValue{tsNodeNetKey("addSubnet")},args...)...)}));p.set("check",tsNamedNativeMethod("check",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsBooleanValue(tsNodeNetBlockListCheck(loop,receiver,tsArg(args,0),tsArg(args,1)))}))
	if p.descriptors==nil{p.descriptors=map[string]*tsDescriptor{}};p.descriptors["rules"]=&tsDescriptor{accessor:true,configurable:true,getter:tsNamedNativeMethod("get rules",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{out:=&tsArray{};for _,rule:=range tsNodeNetSingleton.blocks[tsNodeNetPointer(receiver)]{out.values=append(out.values,tsNodeNetKey(rule.text))};return tsArrayValue(out)})}
	var constructor *tsFunction;constructor=tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{object:=tsNewObject();object.prototype=tsObjectValue(p);tsNodeNetSingleton.blocks[uintptr(unsafe.Pointer(object))]=nil;return tsObjectValue(object)});constructor.name="BlockList";constructor.length=0;constructor.constructible=true;constructor.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsCall(loop,tsFunctionValue(constructor))};constructor.properties=tsNewObject();constructor.properties.set("prototype",tsObjectValue(p));constructor.properties.set("isBlockList",tsNamedNativeMethod("isBlockList",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{value:=tsArg(args,0);_,ok:=tsNodeNetSingleton.blocks[tsNodeNetPointer(value)];return tsBooleanValue(ok)}));tsNodeNetSingleton.blockListProto=p;return tsFunctionValue(constructor)
}
func tsNodeNetSocketAddressCtor(loop *tsLoop)tsValue {
	p:=tsNodeNetSingleton.socketAddressProto;var constructor *tsFunction;constructor=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{options:=tsArg(args,0);address:=tsText(tsNodeNetGet(loop,options,"address"));if address=="undefined"||address==""{address="127.0.0.1"};port:=int(tsNumber(tsNodeNetGet(loop,options,"port")));family:=tsText(tsNodeNetGet(loop,options,"family"));if family=="undefined"||family==""{if strings.Contains(address,":"){family="ipv6"}else{family="ipv4"}};object:=tsNewObject();object.prototype=tsObjectValue(p);object.set("address",tsNodeNetKey(address));object.set("port",tsNumberValue(float64(port)));object.set("family",tsNodeNetKey(family));object.set("flowlabel",tsNumberValue(tsNumber(tsNodeNetGet(loop,options,"flowlabel"))));return tsObjectValue(object)});constructor.name="SocketAddress";constructor.length=1;constructor.constructible=true;constructor.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsCall(loop,tsFunctionValue(constructor),args...)};constructor.properties=tsNewObject();constructor.properties.set("prototype",tsObjectValue(p));constructor.properties.set("isSocketAddress",tsNamedNativeMethod("isSocketAddress",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{value:=tsArg(args,0);return tsBooleanValue(value.kind==tsObjectKind&&tsGet(loop,value,tsNodeNetKey("address")).kind==tsStringKind&&tsGet(loop,value,tsNodeNetKey("port")).kind==tsNumberKind)}));constructor.properties.set("parse",tsNamedNativeMethod("parse",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{host,portText,err:=net.SplitHostPort(tsText(tsArg(args,0)));if err!=nil{return tsNull};port,err:=strconv.Atoi(portText);if err!=nil{return tsNull};options:=tsNewObject();options.set("address",tsNodeNetKey(host));options.set("port",tsNumberValue(float64(port)));return tsCall(loop,tsFunctionValue(constructor),tsObjectValue(options))}));tsNodeNetSingleton.socketAddressProto=p;return tsFunctionValue(constructor)
}
func tsNodeNetNormalizeArgs(args []tsValue)tsValue {
	options:=tsNewObject();callback:=tsNull;if len(args)==0{return tsArrayValue(&tsArray{values:[]tsValue{tsObjectValue(options),callback}})};first:=args[0];if first.kind==tsObjectKind||first.kind==tsInstanceKind{options=(*tsObject)(first.ref)}else if first.kind==tsStringKind{if port,err:=strconv.Atoi(tsText(first));err==nil{options.set("port",tsNumberValue(float64(port)))}else{options.set("path",first)}}else{options.set("port",first);if len(args)>1&&args[1].kind==tsStringKind{options.set("host",args[1])}};if len(args)>0&&args[len(args)-1].kind==tsFunctionKind{callback=args[len(args)-1]};return tsArrayValue(&tsArray{values:[]tsValue{tsObjectValue(options),callback}})
}
type tsNodeNetKeyLogFile struct{path string}
func (writer tsNodeNetKeyLogFile)Write(line []byte)(int,error){file,err:=os.OpenFile(writer.path,os.O_WRONLY|os.O_CREATE|os.O_APPEND,0600);if err!=nil{return len(line),nil};_,_ = file.Write(line);_ = file.Close();return len(line),nil}
func tsNodeNetKeyLogWriter(loop *tsLoop)io.Writer{tsNodeNetModule(loop);if !tsNodeNetSingleton.tlsKeylogConfigured{return nil};if !tsNodeNetSingleton.tlsKeylogWarned{tsNodeNetSingleton.tlsKeylogWarned=true;process:=tsNodeProcessModule(loop);emit:=tsNodeNetGet(loop,process,"emitWarning");warning:=tsNodeNetKey("Using --tls-keylog makes TLS connections insecure by writing secret key material to file "+tsNodeNetSingleton.tlsKeylogPath);if emit.kind==tsFunctionKind{tsCall(loop,emit,warning)}};return tsNodeNetKeyLogFile{path:tsNodeNetSingleton.tlsKeylogPath}}
func tsNodeNetReadExecArgv(loop *tsLoop){process:=tsNodeProcessModule(loop);argv:=tsNodeNetGet(loop,process,"execArgv");if argv.kind!=tsArrayKind{return};values:=(*tsArray)(argv.ref).values;for index:=0;index<len(values);index++{arg:=tsText(values[index]);switch arg{case "--no-network-family-autoselection","--no-enable-network-family-autoselection":tsNodeNetSingleton.defaultAutoSelectFamily=false;case "--network-family-autoselection","--enable-network-family-autoselection":tsNodeNetSingleton.defaultAutoSelectFamily=true;case "--network-family-autoselection-attempt-timeout":if index+1<len(values){if n,err:=strconv.Atoi(tsText(values[index+1]));err==nil&&n>=1{if n<10{n=10};tsNodeNetSingleton.defaultAutoSelectFamilyAttemptTimeout=n};index++};case "--tls-keylog":if index+1<len(values){tsNodeNetSingleton.tlsKeylogConfigured=true;tsNodeNetSingleton.tlsKeylogPath=tsText(values[index+1]);index++}};if strings.HasPrefix(arg,"--network-family-autoselection-attempt-timeout="){if n,err:=strconv.Atoi(strings.TrimPrefix(arg,"--network-family-autoselection-attempt-timeout="));err==nil&&n>=1{if n<10{n=10};tsNodeNetSingleton.defaultAutoSelectFamilyAttemptTimeout=n}};if strings.HasPrefix(arg,"--tls-keylog="){tsNodeNetSingleton.tlsKeylogConfigured=true;tsNodeNetSingleton.tlsKeylogPath=strings.TrimPrefix(arg,"--tls-keylog=")}}}
func tsNodeNetCreateConnection(loop *tsLoop,args []tsValue)tsValue {
	normalized:=tsNodeNetNormalizeArgs(args);values:=(*tsArray)(normalized.ref).values;socket:=tsNodeNetSocketCreate(loop,tsArg(values,0));timeout:=tsNodeNetGet(loop,tsArg(values,0),"timeout");if !tsNullish(timeout)&&tsNumber(timeout)>0{state:=tsNodeNetSocketFor(socket);state.timeout=time.Duration(tsNumber(timeout))*time.Millisecond};tsNodeNetSocketConnect(loop,socket,[]tsValue{normalized}...);return socket
}
func tsNodeNetModule(loop *tsLoop)tsValue {
	tsNodeNetSingleton.Do(func(){tsNodeNetSingleton.sockets=map[uintptr]*tsNodeNetSocketState{};tsNodeNetSingleton.servers=map[uintptr]*tsNodeNetServerState{};tsNodeNetSingleton.blocks=map[uintptr][]tsNodeNetBlockRule{};tsNodeNetSingleton.defaultAutoSelectFamily=true;tsNodeNetSingleton.defaultAutoSelectFamilyAttemptTimeout=500;tsNodeNetReadExecArgv(loop)
		events:=tsNodeEventsModule(loop);eventCtor:=tsNodeNetGet(loop,events,"EventEmitter");eventProto:=tsNodeNetGet(loop,eventCtor,"prototype");duplex:=tsInternalStreamDuplexModule(loop);duplexCtor:=tsNodeNetGet(loop,duplex,"Duplex");duplexProto:=tsNodeNetGet(loop,duplexCtor,"prototype")
		tsNodeNetSingleton.socketProto=tsNewObject();tsNodeNetSingleton.socketProto.prototype=duplexProto;tsNodeNetSingleton.serverProto=tsNewObject();tsNodeNetSingleton.serverProto.prototype=eventProto;tsNodeNetSingleton.blockListProto=tsNewObject();tsNodeNetSingleton.socketAddressProto=tsNewObject()
		tsNodeNetSocketInstall(loop);tsNodeNetServerInstall(loop);tsNodeNetSingleton.socketCtor=tsNodeNetSocketCtor(loop);tsNodeNetSingleton.serverCtor=tsNodeNetServerCtor(loop);tsNodeNetSingleton.blockListCtor=tsNodeNetBlockListCtor(loop);tsNodeNetSingleton.socketAddressCtor=tsNodeNetSocketAddressCtor(loop)
		module:=tsNewObject();module.set("Socket",tsNodeNetSingleton.socketCtor);module.set("Server",tsNodeNetSingleton.serverCtor);module.set("Stream",tsNodeNetSingleton.socketCtor);module.set("BlockList",tsNodeNetSingleton.blockListCtor);module.set("SocketAddress",tsNodeNetSingleton.socketAddressCtor);module.set("createServer",tsNamedNativeMethod("createServer",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeNetServerCreate(loop,args...)}));module.set("createConnection",tsNamedNativeMethod("createConnection",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeNetCreateConnection(loop,args)}));module.set("connect",module.get("createConnection"));module.set("isIP",tsNamedNativeMethod("isIP",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{if tsInternalNetIPv4(tsArg(args,0)){return tsNumberValue(4)};if tsInternalNetIPv6(tsArg(args,0)){return tsNumberValue(6)};return tsNumberValue(0)}));module.set("isIPv4",tsNodeNetGet(loop,tsInternalNetIPModule(),"isIPv4"));module.set("isIPv6",tsNodeNetGet(loop,tsInternalNetIPModule(),"isIPv6"));module.set("_normalizeArgs",tsNamedNativeMethod("_normalizeArgs",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeNetNormalizeArgs(args)}));module.set("_setSimultaneousAccepts",tsNamedNativeMethod("_setSimultaneousAccepts",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{if !tsNodeNetSingleton.simultaneousAcceptsWarned{tsNodeNetSingleton.simultaneousAcceptsWarned=true;process:=tsNodeProcessModule(loop);emit:=tsNodeNetGet(loop,process,"emitWarning");if emit.kind==tsFunctionKind{tsCall(loop,emit,tsNodeNetKey("net._setSimultaneousAccepts() is deprecated and will be removed."),tsNodeNetKey("DeprecationWarning"),tsNodeNetKey("DEP0121"))}};return tsU}));module.set("getDefaultAutoSelectFamily",tsNamedNativeMethod("getDefaultAutoSelectFamily",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsBooleanValue(tsNodeNetSingleton.defaultAutoSelectFamily)}));module.set("setDefaultAutoSelectFamily",tsNamedNativeMethod("setDefaultAutoSelectFamily",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{value:=tsArg(args,0);if value.kind!=tsBooleanKind{tsFSInvalid("autoSelectFamilyDefault","a boolean")};tsNodeNetSingleton.defaultAutoSelectFamily=tsTruthy(value);return tsBooleanValue(tsNodeNetSingleton.defaultAutoSelectFamily)}));module.set("getDefaultAutoSelectFamilyAttemptTimeout",tsNamedNativeMethod("getDefaultAutoSelectFamilyAttemptTimeout",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNumberValue(float64(tsNodeNetSingleton.defaultAutoSelectFamilyAttemptTimeout))}));module.set("setDefaultAutoSelectFamilyAttemptTimeout",tsNamedNativeMethod("setDefaultAutoSelectFamilyAttemptTimeout",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{value:=tsArg(args,0);if !tsIsNumeric(value)||tsNumber(value)!=math.Trunc(tsNumber(value))||tsNumber(value)<1||tsNumber(value)>math.MaxInt32{tsFSInvalid("value","an integer greater than or equal to 1")};number:=int(tsNumber(value));if number<10{number=10};tsNodeNetSingleton.defaultAutoSelectFamilyAttemptTimeout=number;return tsNumberValue(float64(number))}));module.set("default",tsObjectValue(module));tsNodeNetSingleton.module=tsObjectValue(module)
	});return tsNodeNetSingleton.module
}
`
