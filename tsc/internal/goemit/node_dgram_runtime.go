package goemit

// nodeDgramRuntime maps the UDP Socket API onto Go's standard net.UDPConn.
// Cluster descriptor handoff, kernel back-pressure queues, and native
// multicast socket options remain approximations at this prototype boundary.
const nodeDgramRuntime = `
const tsNodeDgramUnbound = 0
const tsNodeDgramBinding = 1
const tsNodeDgramBound = 2
const tsNodeDgramDisconnected = 0
const tsNodeDgramConnecting = 1
const tsNodeDgramConnected = 2
type tsNodeDgramState struct {
	loop *tsLoop
	socket, handle, stateObject tsValue
	conn *net.UDPConn
	remote *net.UDPAddr
	typeName string
	bindState, connectState int
	receiving, closed, unrefOnBind bool
	recvBufferSize, sendBufferSize int
	queueSize, queueCount int
	receiveBlockList, sendBlockList, queue, abortDispose tsValue
	reuseAddr, reusePort, ipv6Only tsValue
	joined map[string]bool
	broadcast, multicastLoopback bool
}
var tsNodeDgramSingleton struct { module,constructor tsValue; prototype *tsObject; key,asyncDispose,dispose tsValue; states map[*tsObject]*tsNodeDgramState; channel tsValue }

func tsNodeDgramKey(loop *tsLoop,name string)tsValue{return tsStringReference(tsStringUTF8(name))}
func tsNodeDgramGet(loop *tsLoop,value tsValue,name string)tsValue{return tsGet(loop,value,tsNodeDgramKey(loop,name))}
func tsNodeDgramSet(loop *tsLoop,value tsValue,name string,item tsValue){tsSet(loop,value,tsNodeDgramKey(loop,name),item)}
func tsNodeDgramStateFor(value tsValue)*tsNodeDgramState{if value.kind!=tsObjectKind{return nil};return tsNodeDgramSingleton.states[(*tsObject)(value.ref)]}
func tsNodeDgramError(loop *tsLoop,name,code,message string)tsValue{return tsFSException(name,code,message)}
func tsNodeDgramThrow(loop *tsLoop,name,code,message string){panic(tsThrown{tsNodeDgramError(loop,name,code,message)})}
func tsNodeDgramEmit(loop *tsLoop,socket tsValue,event string,args ...tsValue)tsValue{return tsCallReceiver(loop,tsNodeDgramGet(loop,socket,"emit"),socket,append([]tsValue{tsNodeDgramKey(loop,event)},args...)...)}
func tsNodeDgramPost(loop *tsLoop,callback func()){loop.post(callback)}
func tsNodeDgramCallback(loop *tsLoop,callback tsValue,args ...tsValue){if callback.kind!=tsFunctionKind{return};values:=append([]tsValue(nil),args...);tsNodeDgramPost(loop,func(){tsCall(loop,callback,values...)})}
func tsNodeDgramValidate(loop *tsLoop,name string,args ...tsValue)tsValue{return tsCall(loop,tsNodeDgramGet(loop,tsInternalValidatorsModule(loop),name),args...)}
func tsNodeDgramStateSync(loop *tsLoop,state *tsNodeDgramState){if state==nil{return};tsNodeDgramSet(loop,state.stateObject,"handle",state.handle);tsNodeDgramSet(loop,state.stateObject,"receiving",tsBooleanValue(state.receiving));tsNodeDgramSet(loop,state.stateObject,"bindState",tsNumberValue(float64(state.bindState)));tsNodeDgramSet(loop,state.stateObject,"connectState",tsNumberValue(float64(state.connectState)));tsNodeDgramSet(loop,state.stateObject,"queue",state.queue);tsNodeDgramSet(loop,state.stateObject,"unrefOnBind",tsBooleanValue(state.unrefOnBind));tsNodeDgramSet(loop,state.stateObject,"recvBufferSize",tsNumberValue(float64(state.recvBufferSize)));tsNodeDgramSet(loop,state.stateObject,"sendBufferSize",tsNumberValue(float64(state.sendBufferSize)));tsNodeDgramSet(loop,state.stateObject,"receiveBlockList",state.receiveBlockList);tsNodeDgramSet(loop,state.stateObject,"sendBlockList",state.sendBlockList);tsNodeDgramSet(loop,state.stateObject,"reuseAddr",state.reuseAddr);tsNodeDgramSet(loop,state.stateObject,"reusePort",state.reusePort);tsNodeDgramSet(loop,state.stateObject,"ipv6Only",state.ipv6Only)}
func tsNodeDgramHealth(loop *tsLoop,state *tsNodeDgramState){if state==nil||state.closed||tsNullish(state.handle){tsNodeDgramThrow(loop,"Error","ERR_SOCKET_DGRAM_NOT_RUNNING","Socket is closed")}}
func tsNodeDgramAddrValue(loop *tsLoop,address *net.UDPAddr)tsValue{if address==nil{return tsNull};family:="IPv6";if address.IP.To4()!=nil{family="IPv4"};out:=tsNewObject();out.set("address",tsNodeDgramKey(loop,address.IP.String()));out.set("family",tsNodeDgramKey(loop,family));out.set("port",tsNumberValue(float64(address.Port)));return tsObjectValue(out)}
func tsNodeDgramResolve(loop *tsLoop,state *tsNodeDgramState,address tsValue,port int,network string)(*net.UDPAddr,error){host:=tsText(address);if host=="undefined"||host=="null"||host==""{if network=="udp6"{host="::1"}else{host="127.0.0.1"}};return net.ResolveUDPAddr(network,net.JoinHostPort(host,strconv.Itoa(port)))}
func tsNodeDgramSocketControl(state *tsNodeDgramState,operation func(int)error)error{
	if state==nil||state.conn==nil{return errors.New("socket is not bound")};raw,err:=state.conn.SyscallConn();if err!=nil{return err};var operationErr error;if err=raw.Control(func(fd uintptr){operationErr=operation(int(fd))});err!=nil{return err};return operationErr
}
func tsNodeDgramSetSocketInt(state *tsNodeDgramState,level,option,value int)error{return tsNodeDgramSocketControl(state,func(fd int)error{return syscall.SetsockoptInt(fd,level,option,value)})}
func tsNodeDgramSocketIPv6(state *tsNodeDgramState)bool{if state==nil||state.conn==nil{return state!=nil&&state.typeName=="udp6"};address,ok:=state.conn.LocalAddr().(*net.UDPAddr);return ok&&address.IP.To4()==nil}
func tsNodeDgramInterfaceIndex(address string)(uint32,error){
	if address==""{return 0,nil};if iface,err:=net.InterfaceByName(address);err==nil{return uint32(iface.Index),nil};ip:=net.ParseIP(address);if ip==nil{return 0,errors.New("invalid multicast interface address")};interfaces,err:=net.Interfaces();if err!=nil{return 0,err};for _,iface:=range interfaces{addresses,_:=iface.Addrs();for _,entry:=range addresses{candidate:=entry.String();if slash:=strings.IndexByte(candidate,'/');slash>=0{candidate=candidate[:slash]};if net.ParseIP(candidate)!=nil&&net.ParseIP(candidate).Equal(ip){return uint32(iface.Index),nil}}};return 0,errors.New("multicast interface address is not assigned locally")
}
func tsNodeDgramIPv4Interface(address string)([4]byte,error){
	var result [4]byte;if address==""{return result,nil};ip:=net.ParseIP(address);if ip!=nil&&ip.To4()!=nil{copy(result[:],ip.To4());return result,nil};iface,err:=net.InterfaceByName(address);if err!=nil{return result,errors.New("invalid IPv4 multicast interface")};addresses,err:=iface.Addrs();if err!=nil{return result,err};for _,entry:=range addresses{candidate:=entry.String();if slash:=strings.IndexByte(candidate,'/');slash>=0{candidate=candidate[:slash]};if parsed:=net.ParseIP(candidate);parsed!=nil&&parsed.To4()!=nil{copy(result[:],parsed.To4());return result,nil}};return result,errors.New("interface has no IPv4 address")
}
func tsNodeDgramSetMulticastInterface(state *tsNodeDgramState,address string)error{
	if tsNodeDgramSocketIPv6(state){index,err:=tsNodeDgramInterfaceIndex(address);if err!=nil{return err};return tsNodeDgramSetSocketInt(state,syscall.IPPROTO_IPV6,syscall.IPV6_MULTICAST_IF,int(index))};iface,err:=tsNodeDgramIPv4Interface(address);if err!=nil{return err};return tsNodeDgramSocketControl(state,func(fd int)error{return syscall.SetsockoptInet4Addr(fd,syscall.IPPROTO_IP,syscall.IP_MULTICAST_IF,iface)})
}
func tsNodeDgramMembership(state *tsNodeDgramState,group,interfaceAddress string,join bool)error{
	groupIP:=net.ParseIP(group);if groupIP==nil{return errors.New("invalid multicast group")};option4,option6:=syscall.IP_DROP_MEMBERSHIP,syscall.IPV6_LEAVE_GROUP;if join{option4=syscall.IP_ADD_MEMBERSHIP;option6=syscall.IPV6_JOIN_GROUP}
	return tsNodeDgramSocketControl(state,func(fd int)error{if group4:=groupIP.To4();group4!=nil{iface,err:=tsNodeDgramIPv4Interface(interfaceAddress);if err!=nil{return err};var request syscall.IPMreq;copy(request.Multiaddr[:],group4);request.Interface=iface;return syscall.SetsockoptIPMreq(fd,syscall.IPPROTO_IP,option4,&request)};group16:=groupIP.To16();if group16==nil{return errors.New("invalid multicast group")};index,err:=tsNodeDgramInterfaceIndex(interfaceAddress);if err!=nil{return err};var request syscall.IPv6Mreq;copy(request.Multiaddr[:],group16);request.Interface=index;return syscall.SetsockoptIPv6Mreq(fd,syscall.IPPROTO_IPV6,option6,&request)})
}
func boolToInt(value bool)int{if value{return 1};return 0}
func tsNodeDgramSocketOptionFailure(loop *tsLoop,err error,operation string){if err!=nil{panic(tsThrown{tsNodeNetError(err,operation)})}}
func tsNodeDgramBlock(loop *tsLoop,list tsValue,address string)bool{if tsNullish(list){return false};check:=tsNodeDgramGet(loop,list,"check");if check.kind!=tsFunctionKind{return false};family:="ipv4";if strings.Contains(address,":"){family="ipv6"};return tsTruthy(tsCallReceiver(loop,check,list,tsNodeDgramKey(loop,address),tsNodeDgramKey(loop,family)))}
func tsNodeDgramReceive(loop *tsLoop,state *tsNodeDgramState){conn:=state.conn;if conn==nil{return};go func(){buffer:=make([]byte,65535);for{n,remote,err:=conn.ReadFromUDP(buffer);if err!=nil{if state.closed{return};current:=err;loop.post(func(){if !state.closed{tsNodeDgramEmit(loop,state.socket,"error",tsNodeNetError(current,"recvmsg"))}});return};data:=append([]byte(nil),buffer[:n]...);source:=*remote;loop.post(func(){if state.closed||!state.receiving{return};if state.connectState==tsNodeDgramConnected&&state.remote!=nil&&(state.remote.IP.String()!=source.IP.String()||state.remote.Port!=source.Port){return};if tsNodeDgramBlock(loop,state.receiveBlockList,source.IP.String()){return};family:="IPv6";if source.IP.To4()!=nil{family="IPv4"};info:=tsNewObject();info.set("address",tsNodeDgramKey(loop,source.IP.String()));info.set("port",tsNumberValue(float64(source.Port)));info.set("family",tsNodeDgramKey(loop,family));info.set("size",tsNumberValue(float64(len(data))));tsNodeDgramEmit(loop,state.socket,"message",tsNodeBuffer(data),tsObjectValue(info))})}}()}
func tsNodeDgramBind(loop *tsLoop,state *tsNodeDgramState,args []tsValue)tsValue{
	tsNodeDgramHealth(loop,state)
	if state.bindState!=tsNodeDgramUnbound{tsNodeDgramThrow(loop,"Error","ERR_SOCKET_ALREADY_BOUND","Socket is already bound")}
	state.bindState=tsNodeDgramBinding
	address:=tsNodeDgramKey(loop,"");port:=tsNumberValue(0);var callback tsValue
	first:=tsArg(args,0)
	if first.kind==tsFunctionKind{callback=first}else if first.kind==tsObjectKind||first.kind==tsInstanceKind{
		if tsNodeDgramGet(loop,first,"recvStart").kind==tsFunctionKind{tsNodeDgramThrow(loop,"Error","ERR_NOT_IMPLEMENTED","Socket.prototype.bind(handle) is not implemented")}
		address=tsNodeDgramGet(loop,first,"address");if tsNullish(address){address=tsNodeDgramKey(loop,"")};port=tsNodeDgramGet(loop,first,"port");if tsIsUndefined(port){port=tsNumberValue(0)}
	}else{port=first;if len(args)>1&&tsArg(args,1).kind==tsStringKind{address=tsArg(args,1)}}
	if len(args)>0&&tsArg(args,len(args)-1).kind==tsFunctionKind{callback=tsArg(args,len(args)-1)}
	if callback.kind==tsFunctionKind{once:=tsNodeDgramGet(loop,state.socket,"once");tsCallReceiver(loop,once,state.socket,tsNodeDgramKey(loop,"listening"),callback)}
	if first.kind==tsObjectKind||first.kind==tsInstanceKind{fd:=tsNodeDgramGet(loop,first,"fd");if tsInternalDgramIsInt32(loop,fd)&&tsNumber(fd)>0{
		internal:=tsInternalDgramModule(loop);guess:=tsCall(loop,tsNodeDgramGet(loop,internal,"guessHandleType"),fd);if tsText(guess)!="UDP"{tsNodeDgramThrow(loop,"TypeError","ERR_INVALID_FD_TYPE","Unsupported socket type: "+tsText(guess))};opened:=tsNodeDgramCall(loop,state.handle,"open",fd);if tsNumber(opened)!=0{code:="EINVAL";if tsNumber(opened)==-17{code="EEXIST"};tsNodeDgramThrow(loop,"Error",code,"open "+code)};state.conn=tsInternalDgramConn(loop,state.handle);if state.conn==nil{tsNodeDgramThrow(loop,"Error","EBADF","open EBADF")};state.closed=false;state.bindState=tsNodeDgramBound;state.receiving=true;tsNodeDgramStateSync(loop,state);tsNodeDgramReceive(loop,state);tsNodeDgramPost(loop,func(){if !state.closed{tsNodeDgramEmit(loop,state.socket,"listening")}});return state.socket
	}}
	if tsIsUndefined(port)||tsNullish(port){port=tsNumberValue(0)}
	portValue:=tsCall(loop,tsNodeDgramGet(loop,tsInternalValidatorsModule(loop),"validatePort"),port,tsNodeDgramKey(loop,"port"),tsBooleanValue(false));portNumber:=int(tsNumber(portValue))
	if tsNullish(address)||tsText(address)==""{if state.typeName=="udp6"{address=tsNodeDgramKey(loop,"::")}else{address=tsNodeDgramKey(loop,"0.0.0.0")}}
	tsNodeDgramValidate(loop,"validateString",address,tsNodeDgramKey(loop,"address"));host:=tsText(address);network:=state.typeName
	resolved,err:=net.ResolveUDPAddr(network,net.JoinHostPort(host,strconv.Itoa(portNumber)))
	if err==nil{var conn *net.UDPConn;conn,err=net.ListenUDP(network,resolved);if err==nil{state.conn=conn;state.remote=nil;state.closed=false;state.bindState=tsNodeDgramBound;state.receiving=true;if state.recvBufferSize>0{err=conn.SetReadBuffer(state.recvBufferSize)};if err==nil&&state.sendBufferSize>0{err=conn.SetWriteBuffer(state.sendBufferSize)};if err==nil{tsInternalDgramSetConn(loop,state.handle,conn);tsNodeDgramStateSync(loop,state);tsNodeDgramReceive(loop,state);if state.unrefOnBind{state.unrefOnBind=false};tsNodeDgramPost(loop,func(){if !state.closed{tsNodeDgramEmit(loop,state.socket,"listening")}});return state.socket};conn.Close();state.conn=nil}}
	state.bindState=tsNodeDgramUnbound;tsNodeDgramStateSync(loop,state);problem:=tsNodeNetError(err,"bind");tsNodeDgramPost(loop,func(){tsNodeDgramEmit(loop,state.socket,"error",problem)});return state.socket
}
func tsNodeDgramBytes(loop *tsLoop,value tsValue)[]byte{
	if value.kind==tsArrayKind{out:=[]byte{};for _,entry:=range (*tsArray)(value.ref).values{out=append(out,tsNodeDgramBytes(loop,entry)...)};return out}
	if value.kind==tsStringKind||value.kind==tsTypedArrayKind||value.kind==tsDataViewKind{converted:=tsNodeBufferFrom(loop,[]tsValue{value});return append([]byte(nil),tsNodeBufferByteView(converted)...)}
	tsNodeDgramThrow(loop,"TypeError","ERR_INVALID_ARG_TYPE","The \"buffer\" argument must be of type string or an instance of Buffer, TypedArray, or DataView");return nil
}
func tsNodeDgramSendFailure(loop *tsLoop,state *tsNodeDgramState,callback,problem tsValue,emitWithoutCallback bool){
	if callback.kind==tsFunctionKind{tsNodeDgramCallback(loop,callback,problem);return};if emitWithoutCallback{tsNodeDgramPost(loop,func(){if !state.closed{tsNodeDgramEmit(loop,state.socket,"error",problem)}})}
}
func tsNodeDgramSend(loop *tsLoop,state *tsNodeDgramState,args []tsValue)tsValue{
	tsNodeDgramHealth(loop,state);buffer:=tsArg(args,0);offset,length:=0,-1;port,address,callback:=tsU,tsU,tsU;connected:=state.connectState==tsNodeDgramConnected
	if len(args)>1&&tsIsNumeric(tsArg(args,1))&&len(args)>2&&tsIsNumeric(tsArg(args,2)){offset=int(tsNumber(tsArg(args,1)));length=int(tsNumber(tsArg(args,2)));if connected{callback=tsArg(args,3)}else{port=tsArg(args,3);address=tsArg(args,4);callback=tsArg(args,5)}}else if connected{callback=tsArg(args,1)}else{port=tsArg(args,1);address=tsArg(args,2);callback=tsArg(args,3)}
	if len(args)>0&&tsArg(args,len(args)-1).kind==tsFunctionKind{callback=tsArg(args,len(args)-1)}
	data:=tsNodeDgramBytes(loop,buffer);if length>=0{if offset<0||length<0||offset>len(data)||offset+length>len(data){tsNodeDgramThrow(loop,"RangeError","ERR_BUFFER_OUT_OF_BOUNDS","offset or length is outside the bounds of the Buffer")};data=append([]byte(nil),data[offset:offset+length]...)}
	remote:=state.remote
	if !connected{validated:=tsCall(loop,tsNodeDgramGet(loop,tsInternalValidatorsModule(loop),"validatePort"),port,tsNodeDgramKey(loop,"Port"),tsBooleanValue(false));portNumber:=int(tsNumber(validated));if address.kind==tsFunctionKind{callback=address;address=tsU};if tsNullish(address){address=tsNodeDgramKey(loop,"")};tsNodeDgramValidate(loop,"validateString",address,tsNodeDgramKey(loop,"address"));var err error;remote,err=tsNodeDgramResolve(loop,state,address,portNumber,state.typeName);if err!=nil{tsNodeDgramSendFailure(loop,state,callback,tsNodeNetError(err,"getaddrinfo"),true);return tsU};if tsNodeDgramBlock(loop,state.sendBlockList,remote.IP.String()){tsNodeDgramSendFailure(loop,state,callback,tsNodeDgramError(loop,"Error","ERR_IP_BLOCKED","IP address is blocked: "+remote.IP.String()),false);return tsU}}
	if state.bindState==tsNodeDgramUnbound{tsNodeDgramBind(loop,state,[]tsValue{tsNumberValue(0)});if state.conn==nil{return tsU}};if state.conn==nil{tsNodeDgramSendFailure(loop,state,callback,tsNodeDgramError(loop,"Error","ERR_SOCKET_DGRAM_NOT_RUNNING","Socket is not running"),true);return tsU}
	n,err:=state.conn.WriteToUDP(data,remote);if err!=nil{tsNodeDgramSendFailure(loop,state,callback,tsNodeNetError(err,"send"),true);return tsU};if callback.kind==tsFunctionKind{tsNodeDgramCallback(loop,callback,tsNull,tsNumberValue(float64(n)))};return tsU
}
func tsNodeDgramCreate(loop *tsLoop,args []tsValue)tsValue{
	typeValue,listener:=tsArg(args,0),tsArg(args,1);options:=tsObjectValue(tsNewObject());var lookup tsValue
	if typeValue.kind==tsObjectKind||typeValue.kind==tsInstanceKind{options=typeValue;typeValue=tsNodeDgramGet(loop,options,"type");lookup=tsNodeDgramGet(loop,options,"lookup")}
	typeName:=tsText(typeValue);if typeName!="udp4"&&typeName!="udp6"{tsNodeDgramThrow(loop,"TypeError","ERR_SOCKET_BAD_TYPE","Bad socket type specified. Valid types are: udp4, udp6")}
	validateSize:=func(name string)int{value:=tsNodeDgramGet(loop,options,name);if tsTruthy(value){tsNodeDgramValidate(loop,"validateUint32",value,tsNodeDgramKey(loop,"options."+name))};return int(tsNumber(value))};recvSize,sendSize:=validateSize("recvBufferSize"),validateSize("sendBufferSize")
	for _,name:=range []string{"receiveBlockList","sendBlockList"}{value:=tsNodeDgramGet(loop,options,name);if tsTruthy(value){netModule:=tsNodeNetModule(loop);block:=tsNodeDgramGet(loop,netModule,"BlockList");if !tsTruthy(tsCall(loop,tsNodeDgramGet(loop,block,"isBlockList"),value)){tsNodeDgramThrow(loop,"TypeError","ERR_INVALID_ARG_TYPE","options."+name+" must be a net.BlockList")}}}
	if !tsIsUndefined(lookup){tsNodeDgramValidate(loop,"validateFunction",lookup,tsNodeDgramKey(loop,"lookup"))}
	object:=tsNewObject();object.prototype=tsObjectValue(tsNodeDgramSingleton.prototype);socket:=tsNodeEventCreate(loop,tsObjectValue(object),tsU);stateObject:=tsNewObject();module:=tsInternalDgramModule(loop);handle:=tsConstruct(tsNodeDgramGet(loop,module,"UDP"));state:=&tsNodeDgramState{loop:loop,socket:socket,handle:handle,stateObject:tsObjectValue(stateObject),typeName:typeName,bindState:tsNodeDgramUnbound,connectState:tsNodeDgramDisconnected,recvBufferSize:recvSize,sendBufferSize:sendSize,receiveBlockList:tsNodeDgramGet(loop,options,"receiveBlockList"),sendBlockList:tsNodeDgramGet(loop,options,"sendBlockList"),reuseAddr:tsNodeDgramGet(loop,options,"reuseAddr"),reusePort:tsNodeDgramGet(loop,options,"reusePort"),ipv6Only:tsNodeDgramGet(loop,options,"ipv6Only"),joined:map[string]bool{}}
	tsNodeDgramSingleton.states[object]=state;tsSet(loop,tsObjectValue(object),tsNodeDgramSingleton.key,state.stateObject);tsNodeDgramSet(loop,socket,"type",tsNodeDgramKey(loop,typeName));tsNodeDgramStateSync(loop,state)
	if listener.kind==tsFunctionKind{tsNodeDgramCall(loop,socket,"on",tsNodeDgramKey(loop,"message"),listener)}
	signal:=tsNodeDgramGet(loop,options,"signal");if !tsIsUndefined(signal){tsNodeDgramValidate(loop,"validateAbortSignal",signal,tsNodeDgramKey(loop,"options.signal"));if signal.kind==tsObjectKind&&(*tsObject)(signal.ref).abort!=nil{abort:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{if !state.closed{tsNodeDgramCall(loop,socket,"close")};return tsU}));state.abortDispose=tsNodeDgramCall(loop,tsNodeEventsModule(loop),"addAbortListener",signal,abort)}}
	if tsNullish(tsNodeDgramSingleton.channel){channel:=tsNodeDgramCall(loop,tsNodeDiagnosticsModule(loop),"channel",tsNodeDgramKey(loop,"udp.socket"));tsNodeDgramSingleton.channel=channel};if tsTruthy(tsNodeDgramGet(loop,tsNodeDgramSingleton.channel,"hasSubscribers")){record:=tsNewObject();record.set("socket",socket);tsNodeDgramCall(loop,tsNodeDgramSingleton.channel,"publish",tsObjectValue(record))}
	return socket
}
func tsNodeDgramCall(loop *tsLoop,receiver tsValue,name string,args ...tsValue)tsValue{fn:=tsNodeDgramGet(loop,receiver,name);if fn.kind!=tsFunctionKind{return tsU};return tsCallReceiver(loop,fn,receiver,args...)}
func tsNodeDgramInstall(loop *tsLoop){
	if tsNodeDgramSingleton.prototype!=nil{return};tsNodeEventsInit(loop);tsNodeDgramSingleton.states=map[*tsObject]*tsNodeDgramState{};tsNodeDgramSingleton.key=tsGet(loop,tsInternalDgramModule(loop),"kStateSymbol");tsNodeDgramSingleton.asyncDispose=tsWellKnownSymbol("asyncDispose");tsNodeDgramSingleton.dispose=tsWellKnownSymbol("dispose");proto:=tsNewObject();proto.prototype=tsObjectValue(tsNodeEventsSingleton.prototype);tsNodeDgramSingleton.prototype=proto
	methods:=map[string]string{"bind":"bind","connect":"connect","disconnect":"disconnect","send":"send","sendto":"sendto","close":"close","address":"address","remoteAddress":"remoteAddress","setBroadcast":"setBroadcast","setTTL":"setTTL","setMulticastTTL":"setMulticastTTL","setMulticastLoopback":"setMulticastLoopback","setMulticastInterface":"setMulticastInterface","addMembership":"addMembership","dropMembership":"dropMembership","addSourceSpecificMembership":"addSourceSpecificMembership","dropSourceSpecificMembership":"dropSourceSpecificMembership","ref":"ref","unref":"unref","setRecvBufferSize":"setRecvBufferSize","setSendBufferSize":"setSendBufferSize","getRecvBufferSize":"getRecvBufferSize","getSendBufferSize":"getSendBufferSize","getSendQueueSize":"getSendQueueSize","getSendQueueCount":"getSendQueueCount","_healthCheck":"healthCheck","_stopReceiving":"stopReceiving"}
	for name,operation:=range methods{method:=operation;proto.set(name,tsNamedNativeMethod(name,6,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeDgramMethod(loop,receiver,method,args)}))}
	proto.set(tsPropertyKey(tsNodeDgramSingleton.asyncDispose),tsNamedNativeMethod("[Symbol.asyncDispose]",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{promise:=loop.promise();state:=tsNodeDgramStateFor(receiver);if state==nil||tsNullish(state.handle){promise.settle(tsResult{tsU,false});return tsPromiseValue(promise)};cb:=tsFunctionValue(tsFunc(func(_ *tsLoop,_ ...tsValue)tsValue{promise.settle(tsResult{tsU,false});return tsU}));tsNodeDgramMethod(loop,receiver,"close",[]tsValue{cb});return tsPromiseValue(promise)}))
	proto.set(tsPropertyKey(tsNodeDgramSingleton.dispose),tsNamedNativeMethod("[Symbol.dispose]",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{return tsNodeDgramMethod(loop,receiver,"close",nil)}))
	proto.descriptors=map[string]*tsDescriptor{}
	stateGetter:=func(field string)tsValue{return tsNamedNativeMethod("get _"+field,0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeDgramStateFor(receiver);if state==nil{return tsU};switch field{case "handle":return state.handle;case "receiving":return tsBooleanValue(state.receiving);case "bindState":return tsNumberValue(float64(state.bindState));case "queue":return state.queue;case "reuseAddr":return state.reuseAddr};return tsU})}
	stateSetter:=func(field string)tsValue{return tsNamedNativeMethod("set _"+field,1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeDgramStateFor(receiver);if state==nil{return tsU};value:=tsArg(args,0);switch field{case "handle":state.handle=value;case "receiving":state.receiving=tsTruthy(value);case "bindState":state.bindState=int(tsNumber(value));case "queue":state.queue=value;case "reuseAddr":state.reuseAddr=value};tsNodeDgramStateSync(loop,state);return tsU})}
	for _,name:=range []string{"handle","receiving","bindState","queue","reuseAddr"}{proto.descriptors["_"+name]=&tsDescriptor{accessor:true,configurable:true,getter:stateGetter(name),setter:stateSetter(name)}}
}
func tsNodeDgramModule(loop *tsLoop)tsValue{
	tsNodeDgramInstall(loop);if !tsIsUndefined(tsNodeDgramSingleton.module){return tsNodeDgramSingleton.module};constructor:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeDgramCreate(loop,args)});constructor.name="Socket";constructor.length=2;constructor.constructible=true;constructor.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeDgramCreate(loop,args)};constructor.receiverCall=func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeDgramCreate(loop,args)};constructor.properties=tsNewObject();constructor.properties.set("prototype",tsObjectValue(tsNodeDgramSingleton.prototype));tsNodeDgramSingleton.constructor=tsFunctionValue(constructor);tsNodeDgramSingleton.prototype.set("constructor",tsNodeDgramSingleton.constructor)
	module:=tsNewObject();module.set("Socket",tsNodeDgramSingleton.constructor);module.set("createSocket",tsNamedNativeMethod("createSocket",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeDgramCreate(loop,args)}));module.set("default",tsObjectValue(module));tsNodeDgramSingleton.module=tsObjectValue(module);return tsNodeDgramSingleton.module
}
func tsNodeDgramMethod(loop *tsLoop,receiver tsValue,name string,args []tsValue)tsValue{
	state:=tsNodeDgramStateFor(receiver);if state==nil{tsNodeDgramThrow(loop,"TypeError","ERR_INVALID_THIS","Socket method called on incompatible receiver")}
	arg:=func(index int)tsValue{return tsArg(args,index)}
	switch name{
	case "bind":return tsNodeDgramBind(loop,state,args)
	case "connect":
		tsNodeDgramHealth(loop,state)
		if state.connectState!=tsNodeDgramDisconnected{tsNodeDgramThrow(loop,"Error","ERR_SOCKET_DGRAM_IS_CONNECTED","Socket is already connected")}
		port:=tsCall(loop,tsNodeDgramGet(loop,tsInternalValidatorsModule(loop),"validatePort"),arg(0),tsNodeDgramKey(loop,"Port"),tsBooleanValue(false));address:=arg(1);callback:=arg(2)
		if address.kind==tsFunctionKind{callback=address;address=tsU};if tsNullish(address){address=tsNodeDgramKey(loop,"")};tsNodeDgramValidate(loop,"validateString",address,tsNodeDgramKey(loop,"address"))
		state.connectState=tsNodeDgramConnecting;tsNodeDgramStateSync(loop,state)
		if state.bindState==tsNodeDgramUnbound{tsNodeDgramBind(loop,state,[]tsValue{tsNumberValue(0)})}
		if state.conn==nil{state.connectState=tsNodeDgramDisconnected;tsNodeDgramStateSync(loop,state);problem:=tsNodeDgramError(loop,"Error","ERR_SOCKET_DGRAM_NOT_RUNNING","Socket is not running");if callback.kind==tsFunctionKind{tsNodeDgramCallback(loop,callback,problem)}else{tsNodeDgramPost(loop,func(){tsNodeDgramEmit(loop,receiver,"error",problem)})};return receiver}
		remote,err:=tsNodeDgramResolve(loop,state,address,int(tsNumber(port)),state.typeName)
		if err!=nil{state.connectState=tsNodeDgramDisconnected;tsNodeDgramStateSync(loop,state);problem:=tsNodeNetError(err,"getaddrinfo");if callback.kind==tsFunctionKind{tsNodeDgramCallback(loop,callback,problem)}else{tsNodeDgramPost(loop,func(){tsNodeDgramEmit(loop,receiver,"error",problem)})};return receiver}
		if tsNodeDgramBlock(loop,state.sendBlockList,remote.IP.String()){state.connectState=tsNodeDgramDisconnected;tsNodeDgramStateSync(loop,state);problem:=tsNodeDgramError(loop,"Error","ERR_IP_BLOCKED","IP address is blocked: "+remote.IP.String());if callback.kind==tsFunctionKind{tsNodeDgramCallback(loop,callback,problem)}else{tsNodeDgramPost(loop,func(){tsNodeDgramEmit(loop,receiver,"error",problem)})};return receiver}
		state.remote=remote;state.connectState=tsNodeDgramConnected;tsNodeDgramStateSync(loop,state);if callback.kind==tsFunctionKind{tsNodeDgramCall(loop,receiver,"once",tsNodeDgramKey(loop,"connect"),callback)};tsNodeDgramPost(loop,func(){if !state.closed{tsNodeDgramEmit(loop,receiver,"connect")}});return receiver
	case "disconnect":if state.connectState!=tsNodeDgramConnected{tsNodeDgramThrow(loop,"Error","ERR_SOCKET_DGRAM_NOT_CONNECTED","Socket is not connected")};state.remote=nil;state.connectState=tsNodeDgramDisconnected;tsNodeDgramStateSync(loop,state);return tsU
	case "send":return tsNodeDgramSend(loop,state,args)
	case "sendto":if len(args)<5{tsNodeDgramThrow(loop,"TypeError","ERR_MISSING_ARGS","The \"buffer\", \"offset\", \"length\", \"port\", and \"address\" arguments are required")};tsNodeDgramValidate(loop,"validateNumber",arg(1),tsNodeDgramKey(loop,"offset"));tsNodeDgramValidate(loop,"validateNumber",arg(2),tsNodeDgramKey(loop,"length"));tsNodeDgramValidate(loop,"validateNumber",arg(3),tsNodeDgramKey(loop,"port"));tsNodeDgramValidate(loop,"validateString",arg(4),tsNodeDgramKey(loop,"address"));return tsNodeDgramSend(loop,state,args)
	case "close":
		tsNodeDgramHealth(loop,state);callback:=arg(0);if callback.kind==tsFunctionKind{tsNodeDgramCall(loop,receiver,"once",tsNodeDgramKey(loop,"close"),callback)};state.receiving=false;state.closed=true;state.bindState=tsNodeDgramUnbound;if state.abortDispose.kind==tsObjectKind{dispose:=tsGet(loop,state.abortDispose,tsWellKnownSymbol("dispose"));if dispose.kind==tsFunctionKind{tsCallReceiver(loop,dispose,state.abortDispose)}};if state.conn!=nil{state.conn.Close();state.conn=nil};if !tsNullish(state.handle){tsNodeDgramCall(loop,state.handle,"close")};state.handle=tsNull;tsNodeDgramStateSync(loop,state);tsNodeDgramPost(loop,func(){tsNodeDgramEmit(loop,receiver,"close")});return receiver
	case "address":tsNodeDgramHealth(loop,state);if state.conn==nil{tsNodeDgramThrow(loop,"Error","EBADF","getsockname EBADF")};address,_:=state.conn.LocalAddr().(*net.UDPAddr);return tsNodeDgramAddrValue(loop,address)
	case "remoteAddress":tsNodeDgramHealth(loop,state);if state.conn==nil{tsNodeDgramThrow(loop,"Error","ERR_SOCKET_DGRAM_NOT_RUNNING","Socket is not running")};if state.connectState!=tsNodeDgramConnected||state.remote==nil{tsNodeDgramThrow(loop,"Error","ERR_SOCKET_DGRAM_NOT_CONNECTED","Socket is not connected")};return tsNodeDgramAddrValue(loop,state.remote)
	case "setBroadcast":tsNodeDgramHealth(loop,state);if state.conn==nil{tsNodeDgramThrow(loop,"Error","EBADF","setBroadcast EBADF")};state.broadcast=tsTruthy(arg(0));tsNodeDgramSocketOptionFailure(loop,tsNodeDgramSetSocketInt(state,syscall.SOL_SOCKET,syscall.SO_BROADCAST,boolToInt(state.broadcast)),"setBroadcast");return tsBooleanValue(state.broadcast)
	case "setTTL","setMulticastTTL":value:=arg(0);tsCall(loop,tsNodeDgramGet(loop,tsInternalValidatorsModule(loop),"validateNumber"),value,tsNodeDgramKey(loop,"ttl"));tsNodeDgramHealth(loop,state);if state.conn==nil{tsNodeDgramThrow(loop,"Error","EBADF",name+" EBADF")};ttl:=tsNumber(value);if math.Trunc(ttl)!=ttl||ttl<0||ttl>255||name=="setTTL"&&ttl==0{tsNodeDgramThrow(loop,"RangeError","ERR_SOCKET_BAD_TTL","The ttl argument must be an integer between 1 and 255")};level,option:=syscall.IPPROTO_IP,syscall.IP_TTL;if name=="setMulticastTTL"{option=syscall.IP_MULTICAST_TTL};if tsNodeDgramSocketIPv6(state){level=syscall.IPPROTO_IPV6;if name=="setMulticastTTL"{option=syscall.IPV6_MULTICAST_HOPS}else{option=syscall.IPV6_UNICAST_HOPS}};tsNodeDgramSocketOptionFailure(loop,tsNodeDgramSetSocketInt(state,level,option,int(ttl)),name);return value
	case "setMulticastLoopback":tsNodeDgramHealth(loop,state);if state.conn==nil{tsNodeDgramThrow(loop,"Error","EBADF","setMulticastLoopback EBADF")};state.multicastLoopback=tsTruthy(arg(0));level,option:=syscall.IPPROTO_IP,syscall.IP_MULTICAST_LOOP;if tsNodeDgramSocketIPv6(state){level=syscall.IPPROTO_IPV6;option=syscall.IPV6_MULTICAST_LOOP};tsNodeDgramSocketOptionFailure(loop,tsNodeDgramSetSocketInt(state,level,option,boolToInt(state.multicastLoopback)),"setMulticastLoopback");return tsBooleanValue(state.multicastLoopback)
	case "setMulticastInterface":tsNodeDgramValidate(loop,"validateString",arg(0),tsNodeDgramKey(loop,"interfaceAddress"));tsNodeDgramHealth(loop,state);if state.conn==nil{tsNodeDgramThrow(loop,"Error","ERR_SOCKET_DGRAM_NOT_RUNNING","Socket is not running")};err:=tsNodeDgramSetMulticastInterface(state,tsText(arg(0)));if err!=nil{tsNodeDgramThrow(loop,"Error","EINVAL","setMulticastInterface EINVAL")};return tsU
	case "addMembership","dropMembership":
		group,iface:=arg(0),arg(1);if tsNullish(group)||tsText(group)==""{tsNodeDgramThrow(loop,"TypeError","ERR_MISSING_ARGS","The \"multicastAddress\" argument must be specified")};tsNodeDgramValidate(loop,"validateString",group,tsNodeDgramKey(loop,"multicastAddress"));if !tsIsUndefined(iface){tsNodeDgramValidate(loop,"validateString",iface,tsNodeDgramKey(loop,"interfaceAddress"))};tsNodeDgramHealth(loop,state);if state.conn==nil{if !tsInternalNetIPv4(group)&&!tsInternalNetIPv6(group){tsNodeDgramThrow(loop,"Error","EINVAL",name+" EINVAL")};tsNodeDgramThrow(loop,"Error","ERR_SOCKET_DGRAM_NOT_RUNNING","Socket is not running")};ifaceText:="";if !tsNullish(iface){ifaceText=tsText(iface)};key:=tsText(group)+"%"+ifaceText;err:=tsNodeDgramMembership(state,tsText(group),ifaceText,name=="addMembership");if err!=nil{tsNodeDgramThrow(loop,"Error","EINVAL",name+" EINVAL")};if name=="addMembership"{state.joined[key]=true}else{delete(state.joined,key)};return tsU
	case "addSourceSpecificMembership","dropSourceSpecificMembership":source,group,iface:=arg(0),arg(1),arg(2);tsNodeDgramValidate(loop,"validateString",source,tsNodeDgramKey(loop,"sourceAddress"));tsNodeDgramValidate(loop,"validateString",group,tsNodeDgramKey(loop,"groupAddress"));if !tsIsUndefined(iface){tsNodeDgramValidate(loop,"validateString",iface,tsNodeDgramKey(loop,"interfaceAddress"))};tsNodeDgramHealth(loop,state);if state.conn==nil{if !tsInternalNetIPv4(source)&&!tsInternalNetIPv6(source)||!tsInternalNetIPv4(group)&&!tsInternalNetIPv6(group){tsNodeDgramThrow(loop,"Error","EINVAL",name+" EINVAL")};tsNodeDgramThrow(loop,"Error","ERR_SOCKET_DGRAM_NOT_RUNNING","Socket is not running")};tsNodeDgramThrow(loop,"Error","ENOTSUP",name+" ENOTSUP")
	case "ref":state.unrefOnBind=false;tsNodeDgramStateSync(loop,state);return receiver
	case "unref":state.unrefOnBind=true;tsNodeDgramStateSync(loop,state);return receiver
	case "setRecvBufferSize","setSendBufferSize":size:=arg(0);tsCall(loop,tsNodeDgramGet(loop,tsInternalValidatorsModule(loop),"validateUint32"),size,tsNodeDgramKey(loop,"size"));stateSize:=int(tsNumber(size));if state.conn==nil{syscall:="uv_recv_buffer_size";if name=="setSendBufferSize"{syscall="uv_send_buffer_size"};tsNodeDgramThrow(loop,"Error","ERR_SOCKET_BUFFER_SIZE","Could not get or set buffer size: "+syscall+" returned EBADF")};var err error;if name=="setRecvBufferSize"{err=state.conn.SetReadBuffer(stateSize);state.recvBufferSize=stateSize}else{err=state.conn.SetWriteBuffer(stateSize);state.sendBufferSize=stateSize};if err!=nil{panic(tsThrown{tsNodeNetError(err,"setsockopt")})};return tsU
	case "getRecvBufferSize":tsNodeDgramHealth(loop,state);if state.conn==nil{tsNodeDgramThrow(loop,"Error","ERR_SOCKET_BUFFER_SIZE","Could not get or set buffer size: uv_recv_buffer_size returned EBADF")};return tsNumberValue(float64(state.recvBufferSize))
	case "getSendBufferSize":tsNodeDgramHealth(loop,state);if state.conn==nil{tsNodeDgramThrow(loop,"Error","ERR_SOCKET_BUFFER_SIZE","Could not get or set buffer size: uv_send_buffer_size returned EBADF")};return tsNumberValue(float64(state.sendBufferSize))
	case "getSendQueueSize":return tsNumberValue(float64(state.queueSize))
	case "getSendQueueCount":return tsNumberValue(float64(state.queueCount))
	case "healthCheck":tsNodeDgramHealth(loop,state);return tsU
	case "stopReceiving":state.receiving=false;return tsU
	}
	return tsU
}
`
