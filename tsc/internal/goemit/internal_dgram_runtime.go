package goemit

// internalDgramRuntime supplies the libuv-shaped handle used by node:dgram.
// Go's UDPConn replaces Bun's raw UDP Rust binding; synthetic descriptors work
// for handles created and shared inside this runtime, while foreign descriptor
// adoption is not modeled.
const internalDgramRuntime = `
type tsInternalDgramHandle struct { id int; conn *net.UDPConn; ipv6 bool; closed,adopted bool }
var tsInternalDgramSingleton struct { module tsValue; key,fdKey tsValue; prototype *tsObject; constructor tsValue; next int; handles map[int]*tsInternalDgramHandle }

func tsInternalDgramInit(loop *tsLoop){
	if tsInternalDgramSingleton.prototype!=nil{return}
	tsInternalDgramSingleton.key=tsNewSymbol(tsStringUTF8("state symbol"));tsInternalDgramSingleton.fdKey=tsNewSymbol(tsStringUTF8("internal UDP descriptor"));tsInternalDgramSingleton.handles=map[int]*tsInternalDgramHandle{};tsInternalDgramSingleton.next=10000
	proto:=tsNewObject();tsInternalDgramSingleton.prototype=proto
	newHandle:=func(ipv6 bool)*tsInternalDgramHandle{tsInternalDgramSingleton.next++;handle:=&tsInternalDgramHandle{id:tsInternalDgramSingleton.next,ipv6:ipv6};tsInternalDgramSingleton.handles[handle.id]=handle;return handle}
	getHandle:=func(loop *tsLoop,value tsValue)*tsInternalDgramHandle{id:=int(tsNumber(tsGet(loop,value,tsInternalDgramSingleton.fdKey)));return tsInternalDgramSingleton.handles[id]}
	proto.set("bind",tsNamedNativeMethod("bind",3,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsInternalDgramBind(loop,receiver,args,false)}))
	proto.set("bind6",tsNamedNativeMethod("bind6",3,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsInternalDgramBind(loop,receiver,args,true)}))
	proto.set("open",tsNamedNativeMethod("open",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsInternalDgramOpen(loop,receiver,tsArg(args,0))}))
	proto.set("getsockname",tsNamedNativeMethod("getsockname",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsInternalDgramSockName(loop,receiver,tsArg(args,0))}))
	proto.set("close",tsNamedNativeMethod("close",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsInternalDgramClose(loop,receiver,tsArg(args,0))}))
	proto.set("ref",tsNamedNativeMethod("ref",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{return receiver}))
	proto.set("unref",tsNamedNativeMethod("unref",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{return receiver}))
	proto.set("hasRef",tsNamedNativeMethod("hasRef",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{return tsBooleanValue(true)}))
	proto.descriptors=map[string]*tsDescriptor{"fd":{accessor:true,configurable:true,getter:tsNamedNativeMethod("get fd",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{handle:=getHandle(loop,receiver);if handle==nil||handle.closed{return tsNumberValue(-1)};return tsNumberValue(float64(handle.id))})}}
	ctor:=tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{return tsInternalDgramCreate(loop,newHandle(false))});ctor.name="UDP";ctor.constructible=true;ctor.construct=func(loop *tsLoop,_ ...tsValue)tsValue{return tsInternalDgramCreate(loop,newHandle(false))};ctor.properties=tsNewObject();ctor.properties.set("prototype",tsObjectValue(proto));tsInternalDgramSingleton.constructor=tsFunctionValue(ctor)
}
func tsInternalDgramCreate(loop *tsLoop,handle *tsInternalDgramHandle)tsValue{
	object:=tsNewObject();object.prototype=tsObjectValue(tsInternalDgramSingleton.prototype);object.nativeClass=tsNativeClass("UDP");tsSet(loop,tsObjectValue(object),tsInternalDgramSingleton.fdKey,tsNumberValue(float64(handle.id)));tsSet(loop,tsObjectValue(object),tsInternalDgramSingleton.key,tsObjectValue(tsNewObject()));return tsObjectValue(object)
}
func tsInternalDgramBind(loop *tsLoop,receiver tsValue,args []tsValue,ipv6 bool)tsValue{
	id:=int(tsNumber(tsGet(loop,receiver,tsInternalDgramSingleton.fdKey)));handle:=tsInternalDgramSingleton.handles[id];if handle==nil||handle.closed{return tsNumberValue(-9)}
	address:=tsText(tsArg(args,0));if address=="undefined"||address=="null"||address==""{if ipv6{address="::"}else{address="0.0.0.0"}}
	port:=int(tsNumber(tsArg(args,1)));network:="udp4";if ipv6{network="udp6"};resolved,err:=net.ResolveUDPAddr(network,net.JoinHostPort(address,strconv.Itoa(port)));if err!=nil{return tsNumberValue(-22)};conn,err:=net.ListenUDP(network,resolved);if err!=nil{return tsNumberValue(-98)}
	if handle.conn!=nil{handle.conn.Close()};handle.conn=conn;handle.ipv6=ipv6;return tsNumberValue(0)
}
func tsInternalDgramOpen(loop *tsLoop,receiver,descriptor tsValue)tsValue{
	fd:=int(tsNumber(descriptor));if !tsIsNumeric(descriptor)||float64(fd)!=tsNumber(descriptor)||fd<0{return tsNumberValue(-22)};handle:=tsInternalDgramSingleton.handles[fd];if handle==nil||handle.closed{return tsNumberValue(-22)};if handle.adopted{return tsNumberValue(-17)}
	oldID:=int(tsNumber(tsGet(loop,receiver,tsInternalDgramSingleton.fdKey)));old:=tsInternalDgramSingleton.handles[oldID];if old!=nil&&old!=handle&&old.conn!=nil{old.conn.Close();old.closed=true;delete(tsInternalDgramSingleton.handles,oldID)}
	handle.adopted=true;tsSet(loop,receiver,tsInternalDgramSingleton.fdKey,tsNumberValue(float64(fd)));return tsNumberValue(0)
}
func tsInternalDgramSockName(loop *tsLoop,receiver,out tsValue)tsValue{
	id:=int(tsNumber(tsGet(loop,receiver,tsInternalDgramSingleton.fdKey)));handle:=tsInternalDgramSingleton.handles[id];if handle==nil||handle.closed||handle.conn==nil{return tsNumberValue(-9)};address,ok:=handle.conn.LocalAddr().(*net.UDPAddr);if !ok{return tsNumberValue(-9)};family:="IPv4";if address.IP.To4()==nil{family="IPv6"};tsSet(loop,out,tsInternalDgramKey(loop,"address"),tsInternalDgramKey(loop,address.IP.String()));tsSet(loop,out,tsInternalDgramKey(loop,"port"),tsNumberValue(float64(address.Port)));tsSet(loop,out,tsInternalDgramKey(loop,"family"),tsInternalDgramKey(loop,family));return tsNumberValue(0)
}
func tsInternalDgramClose(loop *tsLoop,receiver,callback tsValue)tsValue{
	id:=int(tsNumber(tsGet(loop,receiver,tsInternalDgramSingleton.fdKey)));handle:=tsInternalDgramSingleton.handles[id];if handle!=nil&&!handle.closed{handle.closed=true;if handle.conn!=nil{handle.conn.Close()};delete(tsInternalDgramSingleton.handles,id)};tsSet(loop,receiver,tsInternalDgramSingleton.fdKey,tsNumberValue(-1));if callback.kind==tsFunctionKind{process:=tsNodeProcessModule(loop);next:=tsGet(loop,process,tsInternalDgramKey(loop,"nextTick"));if next.kind==tsFunctionKind{tsCallReceiver(loop,next,process,callback)}};return tsNumberValue(0)
}
func tsInternalDgramConn(loop *tsLoop,receiver tsValue)*net.UDPConn{id:=int(tsNumber(tsGet(loop,receiver,tsInternalDgramSingleton.fdKey)));if handle:=tsInternalDgramSingleton.handles[id];handle!=nil&&!handle.closed{return handle.conn};return nil}
func tsInternalDgramSetConn(loop *tsLoop,receiver tsValue,conn *net.UDPConn) {
	id:=int(tsNumber(tsGet(loop,receiver,tsInternalDgramSingleton.fdKey)));if handle:=tsInternalDgramSingleton.handles[id];handle!=nil{handle.conn=conn;handle.closed=false}
}
func tsInternalDgramKey(loop *tsLoop,name string)tsValue{return tsStringReference(tsStringUTF8(name))}
func tsInternalDgramIsInt32(loop *tsLoop,value tsValue)bool{if !tsIsNumeric(value){return false};number:=tsNumber(value);return math.Trunc(number)==number&&number>=-2147483648&&number<=2147483647}
func tsInternalDgramGuessHandleType(loop *tsLoop,value tsValue)tsValue{if !tsInternalDgramIsInt32(loop,value)||tsNumber(value)<0{return tsInternalDgramKey(loop,"UNKNOWN")};handle:=tsInternalDgramSingleton.handles[int(tsNumber(value))];if handle==nil||handle.closed{return tsInternalDgramKey(loop,"UNKNOWN")};return tsInternalDgramKey(loop,"UDP")}
func tsInternalDgramCreateSocketHandle(loop *tsLoop,args []tsValue)tsValue{
	address,port,addressType,fd,flags:=tsArg(args,0),tsArg(args,1),tsArg(args,2),tsArg(args,3),tsArg(args,4);ctor:=tsInternalDgramSingleton.constructor;handle:=tsConstruct(ctor);if !tsNullish(fd)&&tsInternalDgramIsInt32(loop,fd)&&tsNumber(fd)>0{if tsText(tsInternalDgramGuessHandleType(loop,fd))!="UDP"{tsInternalDgramClose(loop,handle,tsU);return tsNumberValue(-22)};err:=tsInternalDgramOpen(loop,handle,fd);if tsNumber(err)!=0{tsInternalDgramClose(loop,handle,tsU);return err}}else if tsTruthy(port)||tsTruthy(address){method:="bind";if tsText(addressType)=="udp6"{method="bind6"};err:=tsInternalDgramBind(loop,handle,[]tsValue{address,port,flags},method=="bind6");if tsNumber(err)!=0{tsInternalDgramClose(loop,handle,tsU);return err}};return handle
}
func tsInternalDgramModule(loop *tsLoop)tsValue{
	tsInternalDgramInit(loop);if !tsIsUndefined(tsInternalDgramSingleton.module){return tsInternalDgramSingleton.module};module:=tsNewObject();module.set("kStateSymbol",tsInternalDgramSingleton.key);module.set("UDP",tsInternalDgramSingleton.constructor);module.set("_createSocketHandle",tsNamedNativeMethod("_createSocketHandle",5,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsInternalDgramCreateSocketHandle(loop,args)}));module.set("guessHandleType",tsNamedNativeMethod("guessHandleType",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsInternalDgramGuessHandleType(loop,tsArg(args,0))}));module.set("default",tsObjectValue(module));tsInternalDgramSingleton.module=tsObjectValue(module);return tsInternalDgramSingleton.module
}
`
