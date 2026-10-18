package goemit

// internalClusterSharedHandleRuntime adapts Bun's SharedHandle using Go's
// standard net listeners. OS-level descriptor transfer remains a lightweight
// placeholder, while handle lifetime and worker membership are represented.
const internalClusterSharedHandleRuntime = `
type tsInternalClusterSharedState struct { object *tsObject; key string; workers map[int]tsValue; handle tsValue; listener net.Listener; packet net.PacketConn; descriptor *os.File; errno int; sharedOnly bool }
var tsInternalClusterSharedHandles struct { sync.Once; prototype *tsObject; states map[*tsObject]*tsInternalClusterSharedState; constructor tsValue }
func tsInternalClusterSharedFor(receiver tsValue) *tsInternalClusterSharedState { if receiver.kind!=tsObjectKind{return nil};return tsInternalClusterSharedHandles.states[(*tsObject)(receiver.ref)] }
func tsInternalClusterSharedKey(name string) tsValue { return tsStringReference(tsStringUTF8(name)) }
func tsInternalClusterSharedErrno(err error) int {
	var errno syscall.Errno
	if errors.As(err,&errno) { return -int(errno) }
	return -int(syscall.EINVAL)
}
func tsInternalClusterSharedValidateFD(fd tsValue) int {
	if !tsIsNumeric(fd) { return -int(syscall.EINVAL) }
	descriptor:=int(int32(tsNumber(fd)))
	if descriptor<0 { return -int(syscall.EBADF) }
	socketType,err:=syscall.GetsockoptInt(descriptor,syscall.SOL_SOCKET,syscall.SO_TYPE)
	if err!=nil || (socketType!=syscall.SOCK_STREAM&&socketType!=syscall.SOCK_DGRAM) { return -int(syscall.EINVAL) }
	if socketType==syscall.SOCK_STREAM { if _,err:=syscall.Getpeername(descriptor);err==nil { return -int(syscall.EINVAL) } }
	return 0
}
func tsInternalClusterSharedListenConfig(flags int,udp bool) net.ListenConfig {
	ipv6Only:=flags&1!=0
	reusePort:=udp&&flags&4!=0
	config:=net.ListenConfig{}
	if ipv6Only||reusePort {
		config.Control=func(network,address string,raw syscall.RawConn)error {
			var socketErr error
			controlErr:=raw.Control(func(fd uintptr) {
				if reusePort {
					socketErr=syscall.SetsockoptInt(int(fd),syscall.SOL_SOCKET,syscall.SO_REUSEADDR,1)
					if socketErr==nil&&(runtime.GOOS=="darwin"||runtime.GOOS=="ios"||runtime.GOOS=="freebsd") { socketErr=syscall.SetsockoptInt(int(fd),syscall.SOL_SOCKET,0x0200,1) }
				}
				if socketErr==nil&&ipv6Only&&(strings.HasSuffix(network,"6")||strings.Contains(address,"[")) { socketErr=syscall.SetsockoptInt(int(fd),syscall.IPPROTO_IPV6,syscall.IPV6_V6ONLY,1) }
			})
			if controlErr!=nil { return controlErr }
			return socketErr
		}
	}
	return config
}
func tsInternalClusterSharedInit(loop *tsLoop) {
	tsInternalClusterSharedHandles.Do(func(){tsInternalClusterSharedHandles.states=map[*tsObject]*tsInternalClusterSharedState{};prototype:=tsNewObject()
		prototype.set("add",tsNamedNativeMethod("add",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsInternalClusterSharedFor(receiver);if state==nil{tsPropertyFailure("Invalid SharedHandle receiver")};worker,callback:=tsArg(args,0),tsArg(args,1);id:=int(tsNumber(tsGet(loop,worker,tsInternalClusterSharedKey("id"))));state.workers[id]=worker;if callback.kind==tsFunctionKind{tsCall(loop,callback,tsNumberValue(float64(state.errno)),tsNull,state.handle)};return tsU}))
		prototype.set("has",tsNamedNativeMethod("has",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsInternalClusterSharedFor(receiver);if state==nil{return tsBooleanValue(false)};id:=int(tsNumber(tsGet(loop,tsArg(args,0),tsInternalClusterSharedKey("id"))));_,ok:=state.workers[id];return tsBooleanValue(ok)}))
		prototype.set("remove",tsNamedNativeMethod("remove",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsInternalClusterSharedFor(receiver);if state==nil{return tsBooleanValue(false)};id:=int(tsNumber(tsGet(loop,tsArg(args,0),tsInternalClusterSharedKey("id"))));if _,ok:=state.workers[id];!ok{return tsBooleanValue(false)};delete(state.workers,id);if len(state.workers)>0{return tsBooleanValue(false)};if state.listener!=nil{_ = state.listener.Close();state.listener=nil};if state.packet!=nil{_ = state.packet.Close();state.packet=nil};if state.descriptor!=nil{_ = state.descriptor.Close();state.descriptor=nil};if state.handle.kind==tsObjectKind{if path:=tsGet(loop,state.handle,tsInternalClusterSharedKey("path"));path.kind==tsStringKind{_ = os.Remove(tsText(path))}};state.object.set("handle",tsNull);return tsBooleanValue(true)}))
		tsInternalClusterSharedHandles.prototype=prototype;constructor:=tsFunc(func(_ *tsLoop,_ ...tsValue)tsValue{tsPropertyFailure("Class constructor SharedHandle cannot be invoked without 'new'");return tsU});constructor.name="SharedHandle";constructor.constructible=true;constructor.properties=tsNewObject();constructor.properties.set("prototype",tsObjectValue(prototype));constructor.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsInternalClusterSharedCreate(loop,args)};tsInternalClusterSharedHandles.constructor=tsFunctionValue(constructor)
	})
}
func tsInternalClusterSharedCreate(loop *tsLoop,args []tsValue) tsValue {
	key,address,options:=tsArg(args,0),tsArg(args,1),tsArg(args,2);if options.kind!=tsObjectKind&&options.kind!=tsInstanceKind{options=tsObjectValue(tsNewObject())};state:=&tsInternalClusterSharedState{key:tsText(key),workers:map[int]tsValue{},errno:0,sharedOnly:tsTruthy(tsGet(loop,options,tsInternalClusterSharedKey("sharedOnly")))}
	state.handle=tsNull
	fd,port:=tsGet(loop,options,tsInternalClusterSharedKey("fd")),tsGet(loop,options,tsInternalClusterSharedKey("port"))
	if tsIsNumeric(fd)&&tsNumber(fd)>=0 {
		state.errno=tsInternalClusterSharedValidateFD(fd)
		if state.errno==0 { descriptor:=int(int32(tsNumber(fd)));state.descriptor=os.NewFile(uintptr(descriptor),"node-cluster-shared-fd");state.handle=tsNodeClusterObject(loop,map[string]tsValue{"fd":tsNumberValue(float64(descriptor)),"port":port}) }
	} else {
		portNumber:=0;if tsIsNumeric(port){portNumber=int(tsNumber(port))}
		host:="";if !tsNullish(address){host=tsText(address)}
		addressType:=tsText(tsGet(loop,options,tsInternalClusterSharedKey("addressType")))
		flags:=0;if value:=tsGet(loop,options,tsInternalClusterSharedKey("flags"));tsIsNumeric(value){flags=int(tsNumber(value))}
		network:="tcp";udp:=strings.HasPrefix(addressType,"udp")
		if udp { network="udp";if addressType=="udp4"{network="udp4"}else if addressType=="udp6"{network="udp6"} } else if addressType=="4" { network="tcp4" } else if addressType=="6" { network="tcp6" }
		addressText:=net.JoinHostPort(host,strconv.Itoa(portNumber));if portNumber<0||addressType=="-1"{network="unix";addressText=host}
		if udp {
			listenConfig:=tsInternalClusterSharedListenConfig(flags,true);packet,err:=listenConfig.ListenPacket(context.Background(),network,addressText)
			if err!=nil{state.errno=tsInternalClusterSharedErrno(err)}else{
				state.packet=packet;local:=packet.LocalAddr();actualPort:=portNumber
				if addr,ok:=local.(*net.UDPAddr);ok{actualPort=addr.Port;if udpConn,ok:=packet.(*net.UDPConn);ok{if file,fileErr:=udpConn.File();fileErr==nil{state.descriptor=file}}}
				state.handle=tsNodeClusterObject(loop,map[string]tsValue{"fd":tsNumberValue(tsInternalClusterSharedDescriptor(state.descriptor)),"port":tsNumberValue(float64(actualPort))})
			}
		} else {
			var listener net.Listener;var err error
			if network=="unix"{listener,err=net.Listen(network,addressText)}else{listenConfig:=tsInternalClusterSharedListenConfig(flags,false);listener,err=listenConfig.Listen(context.Background(),network,addressText)}
			if err!=nil{state.errno=tsInternalClusterSharedErrno(err)}else{
				state.listener=listener;var fdValue tsValue=tsNumberValue(-1);actualPort:=portNumber
				if tcp,ok:=listener.(*net.TCPListener);ok{if addr,ok:=tcp.Addr().(*net.TCPAddr);ok{actualPort=addr.Port};if file,fileErr:=tcp.File();fileErr==nil{state.descriptor=file;fdValue=tsNumberValue(float64(file.Fd()))}}
				handleFields:=map[string]tsValue{"fd":fdValue,"port":tsNumberValue(float64(actualPort))};if network=="unix"{handleFields["path"]=address}
				state.handle=tsNodeClusterObject(loop,handleFields)
			}
		}
	}
	value:=tsNewObject();value.prototype=tsObjectValue(tsInternalClusterSharedHandles.prototype);state.object=value;value.set("key",key);value.set("handle",state.handle);value.set("errno",tsNumberValue(float64(state.errno)));value.set("data",tsU);value.set("sharedOnly",tsBooleanValue(state.sharedOnly));tsInternalClusterSharedHandles.states[value]=state;return tsObjectValue(value)
}
func tsInternalClusterSharedDescriptor(file *os.File) float64 { if file==nil{return -1};return float64(file.Fd()) }
func tsInternalClusterSharedHandleModule(loop *tsLoop) tsValue { tsInternalClusterSharedInit(loop);return tsInternalClusterSharedHandles.constructor }
`
