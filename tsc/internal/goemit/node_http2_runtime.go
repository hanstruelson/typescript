package goemit

// nodeHTTP2Runtime maps Bun's public HTTP/2 wrapper onto stream, socket,
// TLS, event, and settings adapters. The HTTP/2 engine and native frame parser
// are represented by API-shaped session and stream objects.
const nodeHTTP2BaseRuntime = `
type tsNodeHTTP2SessionState struct {
 object tsValue
 server,connected,closed,closing,destroyed,closeEmitted bool
	transportAttached,protocolStarted,goawaySent,receiving,transportDetached,writeBackpressured,corkFlushScheduled bool
	strictSingleValueFields,strictFieldWhitespaceValidation,pendingHeaderCompressionError bool
 engine *tsHTTP2Connection
 socket tsValue
	serverObject tsValue
	transportDataListener,transportCloseListener,transportStartListener,transportDrainListener,transportErrorListener tsValue
	transportStartEvent string
	incoming []byte
 outgoing [][]byte
	corkedWriteBuffer []byte
	pendingTransportWrites [][]byte
	reentrantInput [][]byte
 localSettings,remoteSettings,initialSettings tsValue
 pendingSettingsCallbacks []tsValue
 pingCallbacks map[string][]tsValue
 outstandingPings int
 outstandingSettings,maxOutstandingSettings int
 maxOutstandingPings,maxHeaderListPairs,maxSettings,maxSessionInvalidFrames,maxSessionRejectedStreams,maxSessionMemory,maxSendHeaderBlockLength uint32
 streamResetBurst,streamResetRate uint32
 paddingStrategy uint8
 remoteCustomFilter map[uint16]bool
 remoteCustomValues map[uint16]uint32
 localCustomSettings tsValue
 streams map[int]tsValue
 pushParents map[int]int
	 rejectedStreams uint32
 nextStreamID int
 lastPeerStreamID int
 timeout *tsTask
 encrypted bool
 sendWindow tsHTTP2SendWindow
 recvWindow tsHTTP2RecvWindow
}
type tsNodeHTTP2StreamState struct {
 object,session,headers,rawHeaders,trailers tsValue
 id int
	exposed bool
	closed,destroyed,ended,headersSent,headersReceived,waitForTrailers,trailerReady,sentTrailers,readingPaused bool
 rstCode int
	weight uint16
 timeout *tsTask
 sendWindow tsHTTP2SendWindow
 recvWindow tsHTTP2RecvWindow
 protocolState uint8
	paddingStrategy uint8
 outbound []tsHTTP2OutboundWrite
}
const (
 tsHTTP2StreamIdle uint8 = iota+1
 tsHTTP2StreamOpen
 tsHTTP2StreamReservedLocal
 tsHTTP2StreamReservedRemote
 tsHTTP2StreamHalfClosedLocal
 tsHTTP2StreamHalfClosedRemote
 tsHTTP2StreamClosed
)
const (
 tsHTTP2EventSendHeaders uint8 = iota+1
 tsHTTP2EventRecvHeaders
 tsHTTP2EventSendHeadersEndStream
 tsHTTP2EventRecvHeadersEndStream
 tsHTTP2EventSendEndStream
 tsHTTP2EventRecvEndStream
 tsHTTP2EventSendPushPromise
 tsHTTP2EventRecvPushPromise
 tsHTTP2EventSendRst
 tsHTTP2EventRecvRst
)
const (
 tsHTTP2TransitionProtocolError = iota+1
 tsHTTP2TransitionStreamClosed
)
func tsHTTP2Transition(state,event uint8)(uint8,int){protocol:=tsHTTP2TransitionProtocolError;closed:=tsHTTP2TransitionStreamClosed;switch state{case tsHTTP2StreamIdle:switch event{case tsHTTP2EventSendHeaders,tsHTTP2EventRecvHeaders:return tsHTTP2StreamOpen,0;case tsHTTP2EventSendHeadersEndStream:return tsHTTP2StreamHalfClosedLocal,0;case tsHTTP2EventRecvHeadersEndStream:return tsHTTP2StreamHalfClosedRemote,0;case tsHTTP2EventSendPushPromise:return tsHTTP2StreamReservedLocal,0;case tsHTTP2EventRecvPushPromise:return tsHTTP2StreamReservedRemote,0;default:return state,protocol};case tsHTTP2StreamReservedLocal:switch event{case tsHTTP2EventSendHeaders:return tsHTTP2StreamHalfClosedRemote,0;case tsHTTP2EventSendHeadersEndStream,tsHTTP2EventSendRst,tsHTTP2EventRecvRst:return tsHTTP2StreamClosed,0;default:return state,protocol};case tsHTTP2StreamReservedRemote:switch event{case tsHTTP2EventRecvHeaders:return tsHTTP2StreamHalfClosedLocal,0;case tsHTTP2EventRecvHeadersEndStream,tsHTTP2EventSendRst,tsHTTP2EventRecvRst:return tsHTTP2StreamClosed,0;default:return state,protocol};case tsHTTP2StreamOpen:switch event{case tsHTTP2EventSendEndStream,tsHTTP2EventSendHeadersEndStream:return tsHTTP2StreamHalfClosedLocal,0;case tsHTTP2EventRecvEndStream,tsHTTP2EventRecvHeadersEndStream:return tsHTTP2StreamHalfClosedRemote,0;case tsHTTP2EventSendRst,tsHTTP2EventRecvRst:return tsHTTP2StreamClosed,0;case tsHTTP2EventSendHeaders,tsHTTP2EventRecvHeaders:return tsHTTP2StreamOpen,0;default:return state,protocol};case tsHTTP2StreamHalfClosedLocal:switch event{case tsHTTP2EventRecvEndStream,tsHTTP2EventRecvHeadersEndStream,tsHTTP2EventSendRst,tsHTTP2EventRecvRst:return tsHTTP2StreamClosed,0;case tsHTTP2EventRecvHeaders:return tsHTTP2StreamHalfClosedLocal,0;default:return state,protocol};case tsHTTP2StreamHalfClosedRemote:switch event{case tsHTTP2EventSendEndStream,tsHTTP2EventSendHeadersEndStream,tsHTTP2EventSendRst,tsHTTP2EventRecvRst:return tsHTTP2StreamClosed,0;case tsHTTP2EventSendHeaders:return tsHTTP2StreamHalfClosedRemote,0;case tsHTTP2EventRecvEndStream,tsHTTP2EventRecvHeaders,tsHTTP2EventRecvHeadersEndStream:return state,closed;default:return state,protocol};case tsHTTP2StreamClosed:if event==tsHTTP2EventSendRst||event==tsHTTP2EventRecvRst{return state,0};return state,closed;default:return state,protocol}}
func tsHTTP2CanReceiveData(state uint8)bool{return state==tsHTTP2StreamOpen||state==tsHTTP2StreamHalfClosedLocal}
func tsHTTP2CanSendData(state uint8)bool{return state==tsHTTP2StreamOpen||state==tsHTTP2StreamHalfClosedRemote}
type tsHTTP2SendWindow struct { remaining int64 }
type tsHTTP2RecvWindow struct { size,consumed int64 }
func tsHTTP2NewSendWindow(initial uint32)tsHTTP2SendWindow{return tsHTTP2SendWindow{remaining:int64(initial)}}
func tsHTTP2NewRecvWindow(initial uint32)tsHTTP2RecvWindow{return tsHTTP2RecvWindow{size:int64(initial)}}
func (window tsHTTP2SendWindow)available()int64{if window.remaining>0{return window.remaining};return 0}
func (window *tsHTTP2SendWindow)consume(n int64){window.remaining-=n}
func (window *tsHTTP2SendWindow)increase(increment uint32)int{next:=window.remaining+int64(increment);if next>0x7fffffff{return tsHTTP2ErrorFlowControl};window.remaining=next;return tsHTTP2ErrorNoError}
func (window *tsHTTP2SendWindow)applyInitialDelta(delta int64){window.remaining=min(window.remaining+delta,int64(0x7fffffff))}
func (window *tsHTTP2RecvWindow)onData(n int64){window.consumed+=n}
func (window tsHTTP2RecvWindow)isOverflowed()bool{return window.consumed>window.size}
func (window tsHTTP2RecvWindow)isOverflowedWith(limit int64)bool{return window.consumed>max(limit,window.size)}
func (window tsHTTP2RecvWindow)needsUpdate()bool{return window.consumed>0&&window.consumed>=window.size/2}
func (window *tsHTTP2RecvWindow)takeUpdate()uint32{if window.consumed<=0{return 0};increment:=min(window.consumed,int64(0x7fffffff));window.consumed-=increment;return uint32(increment)}
func (window *tsHTTP2RecvWindow)grow(delta int64){window.size+=delta}
// Wire constants and the 31-bit helper mirror the frame types, flags, error
// codes, and UInt31WithReserved definitions in Bun's h2_frame_parser.rs.
const (
 tsHTTP2FrameHeaderSize = 9
 tsHTTP2MaxWindowSize uint32 = 0x7fffffff
 tsHTTP2MaxStreamID uint32 = 0x7fffffff
 tsHTTP2MaxFrameSize uint32 = 0x00ffffff
 tsHTTP2DefaultFrameSize uint32 = 16384
 tsHTTP2DefaultWindowSize uint32 = 65535
 tsHTTP2DefaultHeaderTableSize uint32 = 4096
 tsHTTP2StreamPrioritySize = 5
)
const (
 tsHTTP2FrameData uint8 = 0x00
 tsHTTP2FrameHeaders uint8 = 0x01
 tsHTTP2FramePriority uint8 = 0x02
 tsHTTP2FrameRSTStream uint8 = 0x03
 tsHTTP2FrameSettings uint8 = 0x04
 tsHTTP2FramePushPromise uint8 = 0x05
 tsHTTP2FramePing uint8 = 0x06
 tsHTTP2FrameGoAway uint8 = 0x07
 tsHTTP2FrameWindowUpdate uint8 = 0x08
 tsHTTP2FrameContinuation uint8 = 0x09
 tsHTTP2FrameAltSvc uint8 = 0x0a
 tsHTTP2FrameOrigin uint8 = 0x0c
)
const (
 tsHTTP2FlagEndStream uint8 = 0x01
 tsHTTP2FlagAck uint8 = 0x01
 tsHTTP2FlagEndHeaders uint8 = 0x04
 tsHTTP2FlagPadded uint8 = 0x08
 tsHTTP2FlagPriorityPresent uint8 = 0x20
)
const (
 tsHTTP2SettingHeaderTableSize uint16 = 0x01
 tsHTTP2SettingEnablePush uint16 = 0x02
 tsHTTP2SettingMaxConcurrentStreams uint16 = 0x03
 tsHTTP2SettingInitialWindowSize uint16 = 0x04
 tsHTTP2SettingMaxFrameSize uint16 = 0x05
 tsHTTP2SettingMaxHeaderListSize uint16 = 0x06
 tsHTTP2SettingEnableConnectProtocol uint16 = 0x08
 tsHTTP2SettingNoRFC7540Priorities uint16 = 0x09
)
type tsHTTP2UInt31WithReserved uint32
func tsHTTP2InitUInt31(value uint32,reserved bool)tsHTTP2UInt31WithReserved{value&=tsHTTP2MaxStreamID;if reserved{value|=0x80000000};return tsHTTP2UInt31WithReserved(value)}
func (value tsHTTP2UInt31WithReserved)reserved()bool{return uint32(value)&0x80000000!=0}
func (value tsHTTP2UInt31WithReserved)uint31()uint32{return uint32(value)&tsHTTP2MaxStreamID}
func (value tsHTTP2UInt31WithReserved)toUint32()uint32{return uint32(value)}
func (value tsHTTP2UInt31WithReserved)write(dst []byte)bool{if len(dst)<4{return false};binary.BigEndian.PutUint32(dst[:4],uint32(value));return true}
type tsHTTP2StreamPriorityWire struct { streamIdentifier tsHTTP2UInt31WithReserved; weight uint8 }
func (priority tsHTTP2StreamPriorityWire)write(dst []byte)bool{if len(dst)<tsHTTP2StreamPrioritySize{return false};if !priority.streamIdentifier.write(dst[:4]){return false};dst[4]=priority.weight;return true}
type tsHTTP2FrameHeader struct { length uint32; frameType,flags uint8; streamID uint32 }
type tsHTTP2HeaderValidation struct { kind,streamID,code int }
func tsHTTP2HasFrameFlag(flags,mask uint8)bool{return flags&mask!=0}
const (
 tsHTTP2HeaderOK = iota
 tsHTTP2HeaderConnectionError
 tsHTTP2HeaderStreamError
)
const (
 tsHTTP2ErrorNoError = iota
 tsHTTP2ErrorProtocol
 tsHTTP2ErrorInternal
 tsHTTP2ErrorFlowControl
 tsHTTP2ErrorSettingsTimeout
 tsHTTP2ErrorStreamClosed
 tsHTTP2ErrorFrameSize
 tsHTTP2ErrorRefusedStream
 tsHTTP2ErrorCancel
 tsHTTP2ErrorCompression
 tsHTTP2ErrorConnect
 tsHTTP2ErrorEnhanceYourCalm
 tsHTTP2ErrorInadequateSecurity
 tsHTTP2ErrorHTTP11Required
 tsHTTP2ErrorMaxPendingSettingsAck
)
type tsNodeHTTP2ServerState struct {object tsValue; secure bool; sessions map[uintptr]tsValue; options tsValue; listener tsValue}
type tsHTTP2SettingEntry struct{id uint16;value uint32}
var tsNodeHTTP2Singleton struct{sync.Once; module tsValue; constants tsValue; sessionPrototype,streamPrototype,serverRequestPrototype,serverResponsePrototype,serverPrototype,secureServerPrototype *tsObject; sessions map[uintptr]*tsNodeHTTP2SessionState; streams map[uintptr]*tsNodeHTTP2StreamState; servers map[uintptr]*tsNodeHTTP2ServerState; sensitiveHeaders *tsSymbol}
func tsNodeHTTP2Key(name string)tsValue{return tsStringReference(tsStringUTF8(name))}
func tsNodeHTTP2Error(code,message string)tsValue{return tsFSException("Error",code,message)}
func tsHTTP2ParseFrameHeader(data []byte)(tsHTTP2FrameHeader,bool){if len(data)<tsHTTP2FrameHeaderSize{return tsHTTP2FrameHeader{},false};return tsHTTP2FrameHeader{length:uint32(data[0])<<16|uint32(data[1])<<8|uint32(data[2]),frameType:data[3],flags:data[4],streamID:tsHTTP2UInt31WithReserved(binary.BigEndian.Uint32(data[5:tsHTTP2FrameHeaderSize])).uint31()},true}
func tsHTTP2WriteFrameHeader(header tsHTTP2FrameHeader)[]byte{out:=make([]byte,tsHTTP2FrameHeaderSize);out[0]=byte(header.length>>16);out[1]=byte(header.length>>8);out[2]=byte(header.length);out[3]=header.frameType;out[4]=header.flags;tsHTTP2InitUInt31(header.streamID,false).write(out[5:]);return out}
func tsHTTP2FrameSizeError(header tsHTTP2FrameHeader)tsHTTP2HeaderValidation{if header.streamID!=0&&header.frameType==0{return tsHTTP2HeaderValidation{kind:tsHTTP2HeaderStreamError,streamID:int(header.streamID),code:tsHTTP2ErrorFrameSize}};return tsHTTP2HeaderValidation{kind:tsHTTP2HeaderConnectionError,code:tsHTTP2ErrorFrameSize}}
func tsHTTP2ValidateFrameHeader(header tsHTTP2FrameHeader,maxFrameSize uint32)tsHTTP2HeaderValidation{
	known:=false
	connectionScope:=false
	streamScope:=false
	fixedLength:=-1
	switch header.frameType{
	case tsHTTP2FrameData,tsHTTP2FrameHeaders,tsHTTP2FramePushPromise,tsHTTP2FrameContinuation:known=true;streamScope=true
	case tsHTTP2FramePriority:known=true;streamScope=true;fixedLength=tsHTTP2StreamPrioritySize
	case tsHTTP2FrameRSTStream:known=true;streamScope=true;fixedLength=4
	case tsHTTP2FrameSettings:known=true;connectionScope=true
	case tsHTTP2FramePing:known=true;connectionScope=true;fixedLength=8
	case tsHTTP2FrameGoAway:known=true;connectionScope=true
	case tsHTTP2FrameWindowUpdate:known=true;fixedLength=4
	case tsHTTP2FrameAltSvc,tsHTTP2FrameOrigin:known=true
	}
	if !known{return tsHTTP2HeaderValidation{kind:tsHTTP2HeaderOK}}
	if header.length>maxFrameSize{return tsHTTP2FrameSizeError(header)}
	if connectionScope&&header.streamID!=0{return tsHTTP2HeaderValidation{kind:tsHTTP2HeaderConnectionError,code:tsHTTP2ErrorProtocol}}
	if streamScope&&header.streamID==0{return tsHTTP2HeaderValidation{kind:tsHTTP2HeaderConnectionError,code:tsHTTP2ErrorProtocol}}
	if fixedLength>=0&&header.length!=uint32(fixedLength){return tsHTTP2FrameSizeError(header)}
	if header.frameType==tsHTTP2FrameSettings&&header.length%6!=0{return tsHTTP2FrameSizeError(header)}
	return tsHTTP2HeaderValidation{kind:tsHTTP2HeaderOK}
}
type tsHTTP2SettingsState struct {headerTableSize,enablePush,maxConcurrentStreams,initialWindowSize,maxFrameSize,maxHeaderListSize,enableConnectProtocol uint32}
func tsHTTP2DefaultSettingsState()tsHTTP2SettingsState{return tsHTTP2SettingsState{headerTableSize:tsHTTP2DefaultHeaderTableSize,enablePush:1,maxConcurrentStreams:math.MaxUint32,initialWindowSize:tsHTTP2DefaultWindowSize,maxFrameSize:tsHTTP2DefaultFrameSize,maxHeaderListSize:65535}}
func tsHTTP2SettingsFromValue(value tsValue)tsHTTP2SettingsState{settings:=tsHTTP2DefaultSettingsState();if !tsNodeHTTP2Object(value){return settings};object:=(*tsObject)(value.ref);for _,name:=range object.order{entry:=object.values[name];switch name{case "headerTableSize":settings.headerTableSize=uint32(tsNumber(entry));case "enablePush":if tsTruthy(entry){settings.enablePush=1}else{settings.enablePush=0};case "maxConcurrentStreams":settings.maxConcurrentStreams=uint32(tsNumber(entry));case "initialWindowSize":settings.initialWindowSize=uint32(tsNumber(entry));case "maxFrameSize":settings.maxFrameSize=uint32(tsNumber(entry));case "maxHeaderListSize","maxHeaderSize":settings.maxHeaderListSize=uint32(tsNumber(entry));case "enableConnectProtocol":if tsTruthy(entry){settings.enableConnectProtocol=1}else{settings.enableConnectProtocol=0}}};return settings}
func (settings *tsHTTP2SettingsState)apply(id uint16,value uint32){switch id{case 1:settings.headerTableSize=value;case 2:settings.enablePush=value;case 3:settings.maxConcurrentStreams=value;case 4:settings.initialWindowSize=value;case 5:settings.maxFrameSize=value;case 6:settings.maxHeaderListSize=value;case 8:settings.enableConnectProtocol=value}}
func tsHTTP2ValidateSetting(id uint16,value uint32)int{switch id{case tsHTTP2SettingEnablePush,tsHTTP2SettingEnableConnectProtocol:if value>1{return tsHTTP2ErrorProtocol};case tsHTTP2SettingInitialWindowSize:if value>tsHTTP2MaxWindowSize{return tsHTTP2ErrorFlowControl};case tsHTTP2SettingMaxFrameSize:if value<tsHTTP2DefaultFrameSize||value>tsHTTP2MaxFrameSize{return tsHTTP2ErrorProtocol}};return tsHTTP2ErrorNoError}
func tsHTTP2ValidateSettingsPayload(data []byte)int{for offset:=0;offset+6<=len(data);offset+=6{id:=binary.BigEndian.Uint16(data[offset:]);value:=binary.BigEndian.Uint32(data[offset+2:]);if code:=tsHTTP2ValidateSetting(id,value);code!=tsHTTP2ErrorNoError{return code}};return tsHTTP2ErrorNoError}
func tsHTTP2ParseContentLengthStrict(value []byte)(uint64,bool){if len(value)==0{return 0,false};var result uint64;for _,b:=range value{if b<'0'||b>'9'{return 0,false};digit:=uint64(b-'0');if result>(math.MaxUint64-digit)/10{return 0,false};result=result*10+digit};return result,true}
const (
 tsHTTP2HPACKUnableToDecode = iota+1
 tsHTTP2HPACKEmptyName
 tsHTTP2HPACKUnableToEncode
)
const tsHTTP2MaxHPACKHeaderSize = 65536
type tsHTTP2DecodedHeader struct { name,value []byte; neverIndexed bool; staticIndex uint16; next int }
type tsHTTP2HPACKBackend interface {
 encode(name,value []byte,neverIndexed bool,dst []byte,offset int)(int,bool)
	decodeBlock(src []byte,maxCapacity uint32)([]tsHTTP2DecodedHeader,bool)
 setEncoderMaxCapacity(uint32)
 close()
}
type tsHTTP2HPACKCoder struct { backend tsHTTP2HPACKBackend; encoderCapacity uint32; pendingCapacity *uint32 }
func tsHTTP2NewHPACKCoder(maxCapacity uint32,backend tsHTTP2HPACKBackend)*tsHTTP2HPACKCoder{return &tsHTTP2HPACKCoder{backend:backend,encoderCapacity:maxCapacity}}
func (coder *tsHTTP2HPACKCoder)queueEncoderCapacity(capacity uint32){if capacity==coder.encoderCapacity&&coder.pendingCapacity==nil{return};next:=capacity;coder.pendingCapacity=&next}
func tsHTTP2WriteHPACKTableSizeUpdate(dst []byte,offset int,value uint32)int{if offset<0||offset>len(dst)||len(dst)-offset<6{return 0};index:=offset;if value<31{dst[index]=0x20|byte(value);return 1};dst[index]=0x3f;index++;rest:=value-31;for rest>=128{dst[index]=byte(rest)|0x80;index++;rest>>=7};dst[index]=byte(rest);index++;return index-offset}
func (coder *tsHTTP2HPACKCoder)takePendingSizeUpdate(dst []byte,offset int)int{if coder.pendingCapacity==nil||offset<0||offset>len(dst)||len(dst)-offset<6{return 0};capacity:=*coder.pendingCapacity;coder.pendingCapacity=nil;if coder.backend!=nil{coder.backend.setEncoderMaxCapacity(capacity)};coder.encoderCapacity=capacity;return tsHTTP2WriteHPACKTableSizeUpdate(dst,offset,capacity)}
func (coder *tsHTTP2HPACKCoder)encode(name,value []byte,neverIndexed bool,dst []byte,offset int)(int,int){if len(name)+len(value)>tsHTTP2MaxHPACKHeaderSize||offset<0||offset>len(dst){return 0,tsHTTP2HPACKUnableToEncode};if coder.backend==nil{return 0,tsHTTP2HPACKUnableToEncode};written,ok:=coder.backend.encode(name,value,neverIndexed,dst,offset);if !ok{return 0,tsHTTP2HPACKUnableToEncode};return written,0}
func (coder *tsHTTP2HPACKCoder)decodeBlock(src []byte,maxCapacity uint32)([]tsHTTP2DecodedHeader,int){if coder.backend==nil{return nil,tsHTTP2HPACKUnableToDecode};headers,ok:=coder.backend.decodeBlock(src,maxCapacity);if !ok{return nil,tsHTTP2HPACKUnableToDecode};for _,header:=range headers{if len(header.name)==0{return nil,tsHTTP2HPACKEmptyName}};return headers,0}
func (coder *tsHTTP2HPACKCoder)close(){if coder.backend!=nil{coder.backend.close();coder.backend=nil}}
func tsNodeHTTP2State(value tsValue)*tsNodeHTTP2SessionState{if value.kind!=tsObjectKind{return nil};return tsNodeHTTP2Singleton.sessions[uintptr(value.ref)]}
func tsNodeHTTP2Stream(value tsValue)*tsNodeHTTP2StreamState{if value.kind!=tsObjectKind{return nil};return tsNodeHTTP2Singleton.streams[uintptr(value.ref)]}
func tsNodeHTTP2Get(loop *tsLoop,value tsValue,name string)tsValue{return tsGet(loop,value,tsNodeHTTP2Key(name))}
func tsNodeHTTP2Set(loop *tsLoop,value tsValue,name string,item tsValue){tsSet(loop,value,tsNodeHTTP2Key(name),item)}
func tsNodeHTTP2Object(value tsValue)bool{return value.kind==tsObjectKind||value.kind==tsInstanceKind}
func tsNodeHTTP2EmptyObject()tsValue{return tsObjectValue(tsNewObject())}
func tsNodeHTTP2SettingsDefaults()tsValue{return tsObjectValue(&tsObject{order:[]string{"headerTableSize","enablePush","maxConcurrentStreams","initialWindowSize","maxFrameSize","maxHeaderListSize","maxHeaderSize","enableConnectProtocol"},values:map[string]tsValue{"headerTableSize":tsNumberValue(4096),"enablePush":tsBooleanValue(true),"maxConcurrentStreams":tsNumberValue(4294967295),"initialWindowSize":tsNumberValue(65535),"maxFrameSize":tsNumberValue(16384),"maxHeaderListSize":tsNumberValue(65535),"maxHeaderSize":tsNumberValue(65535),"enableConnectProtocol":tsBooleanValue(false)}})}
func tsNodeHTTP2StrictSingleValueEnabled(loop *tsLoop,options tsValue)bool{_ = loop;value:=tsFSProperty(options,"strictSingleValueFields");return tsIsUndefined(value)||tsTruthy(value)}
func tsHTTP2IsMalformedHeaderName(name []byte)bool{if len(name)==0{return true};start:=0;if name[0]==':'{start=1};if start==len(name){return true};for _,c:=range name[start:]{if !((c>='a'&&c<='z')||(c>='0'&&c<='9')||strings.ContainsRune("!#$%&'*+-.^_\\x60|~",rune(c))){return true}};return false}
func tsHTTP2IsMalformedHeaderValue(value []byte)bool{for _,c:=range value{if c==0||c=='\r'||c=='\n'{return true}};return false}
func tsHTTP2IsRequestPseudoHeader(name []byte)bool{switch string(name){case ":path",":method",":scheme",":protocol",":authority":return true};return false}
func tsHTTP2IsResponsePseudoHeader(name []byte)bool{return string(name)==":status"}
var tsHTTP2SingleValueHeaders=map[string]bool{"tk":true,"age":true,"dnt":true,"date":true,"etag":true,"from":true,"host":true,":path":true,"range":true,":status":true,":method":true,":scheme":true,"expires":true,"referer":true,"if-match":true,"if-range":true,"location":true,":protocol":true,":authority":true,"user-agent":true,"content-md5":true,"retry-after":true,"content-type":true,"max-forwards":true,"authorization":true,"content-range":true,"if-none-match":true,"last-modified":true,"content-length":true,"content-encoding":true,"content-language":true,"content-location":true,"if-modified-since":true,"if-unmodified-since":true,"proxy-authorization":true,"access-control-max-age":true,"x-content-type-options":true,"upgrade-insecure-requests":true,"access-control-request-method":true,"access-control-allow-credentials":true}
func tsNodeHTTP2SettingError(name string,value tsValue,kind string){message:=fmt.Sprintf("Invalid value for setting %q: %s",name,tsText(value));panic(tsThrown{tsFSException(kind,"ERR_HTTP2_INVALID_SETTING_VALUE",message)})}
func tsNodeHTTP2ValidateSettings(loop *tsLoop,settings tsValue){_ = loop;if !tsNodeHTTP2Object(settings){panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_TYPE","The \"settings\" argument must be of type object")})};ranges:=map[string][2]float64{"headerTableSize":{0,4294967295},"maxConcurrentStreams":{0,4294967295},"initialWindowSize":{0,2147483647},"maxFrameSize":{16384,16777215},"maxHeaderListSize":{0,4294967295},"maxHeaderSize":{0,4294967295}};for name,bounds:=range ranges{value:=tsFSProperty(settings,name);if tsIsUndefined(value){continue};if !tsIsNumeric(value)||math.IsNaN(tsNumber(value))||math.IsInf(tsNumber(value),0)||tsNumber(value)<bounds[0]||tsNumber(value)>bounds[1]{tsNodeHTTP2SettingError(name,value,"RangeError")}};for _,name:=range []string{"enablePush","enableConnectProtocol"}{value:=tsFSProperty(settings,name);if !tsIsUndefined(value)&&value.kind!=tsBooleanKind{tsNodeHTTP2SettingError(name,value,"TypeError")}};custom:=tsFSProperty(settings,"customSettings");if !tsIsUndefined(custom){if !tsNodeHTTP2Object(custom){tsNodeHTTP2SettingError("customSettings",custom,"RangeError")};object:=(*tsObject)(custom.ref);if len(object.order)>10{panic(tsThrown{tsFSException("Error","ERR_HTTP2_TOO_MANY_CUSTOM_SETTINGS","Number of custom settings exceeds MAX_ADDITIONAL_SETTINGS")})};for _,key:=range object.order{id,err:=strconv.Atoi(key);value:=object.values[key];if err!=nil||id<0||id>65535||!tsIsNumeric(value)||math.IsNaN(tsNumber(value))||tsNumber(value)<0||tsNumber(value)>4294967295{tsNodeHTTP2SettingError(key,value,"RangeError")}}}}
func tsNodeHTTP2SettingsEntries(value tsValue,update bool)[]tsHTTP2SettingEntry{entries:=[]tsHTTP2SettingEntry{};if !tsNodeHTTP2Object(value){return entries};addNumber:=func(name string,id uint16){item:=tsFSProperty(value,name);if !tsIsUndefined(item)&&!tsIsNull(item){entries=append(entries,tsHTTP2SettingEntry{id,uint32(tsNumber(item))})}};addNumber("headerTableSize",1);if item:=tsFSProperty(value,"enablePush");!tsIsUndefined(item)&&!tsIsNull(item){n:=uint32(0);if tsTruthy(item){n=1};entries=append(entries,tsHTTP2SettingEntry{2,n})};addNumber("maxConcurrentStreams",3);addNumber("initialWindowSize",4);addNumber("maxFrameSize",5);if update{if item:=tsFSProperty(value,"maxHeaderSize");!tsIsUndefined(item)&&!tsIsNull(item){entries=append(entries,tsHTTP2SettingEntry{6,uint32(tsNumber(item))})}else{addNumber("maxHeaderListSize",6)}}else{if item:=tsFSProperty(value,"maxHeaderListSize");!tsIsUndefined(item)&&!tsIsNull(item){entries=append(entries,tsHTTP2SettingEntry{6,uint32(tsNumber(item))})}else{addNumber("maxHeaderSize",6)}};if item:=tsFSProperty(value,"enableConnectProtocol");!tsIsUndefined(item)&&!tsIsNull(item){n:=uint32(0);if tsTruthy(item){n=1};entries=append(entries,tsHTTP2SettingEntry{8,n})};custom:=tsFSProperty(value,"customSettings");if tsNodeHTTP2Object(custom){source:=(*tsObject)(custom.ref);keys:=append([]string(nil),source.order...);sort.Slice(keys,func(i,j int)bool{a,_:=strconv.Atoi(keys[i]);b,_:=strconv.Atoi(keys[j]);return a<b});for _,name:=range keys{id,_:=strconv.Atoi(name);entries=append(entries,tsHTTP2SettingEntry{uint16(id),uint32(tsNumber(source.values[name]))})}};return entries}
func tsNodeHTTP2PackSettings(loop *tsLoop,value tsValue)tsValue{if tsIsUndefined(value){return tsNodeBuffer(nil)};tsNodeHTTP2ValidateSettings(loop,value);entries:=tsNodeHTTP2SettingsEntries(value,false);data:=make([]byte,len(entries)*6);for i,entry:=range entries{binary.BigEndian.PutUint16(data[i*6:],entry.id);binary.BigEndian.PutUint32(data[i*6+2:],entry.value)};return tsNodeBuffer(data)}
func tsNodeHTTP2UnpackSettings(loop *tsLoop,value,options tsValue)tsValue{if tsIsUndefined(value){return tsNodeHTTP2SettingsDefaults()};if value.kind!=tsTypedArrayKind&&value.kind!=tsDataViewKind&&value.kind!=tsArrayBufferKind{panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_TYPE","The \"buf\" argument must be an instance of Buffer or TypedArray")})};data:=tsNodeBufferByteView(value);if len(data)%6!=0{panic(tsThrown{tsFSException("RangeError","ERR_HTTP2_INVALID_PACKED_SETTINGS_LENGTH","Packed settings length must be a multiple of six")})};out:=tsNewObject();custom:=tsNewObject();known:=map[int]string{1:"headerTableSize",2:"enablePush",3:"maxConcurrentStreams",4:"initialWindowSize",5:"maxFrameSize",6:"maxHeaderListSize",8:"enableConnectProtocol"};hasCustom:=false;for i:=0;i<len(data);i+=6{id:=int(binary.BigEndian.Uint16(data[i:]));number:=binary.BigEndian.Uint32(data[i+2:]);if name,ok:=known[id];ok{if name=="enablePush"||name=="enableConnectProtocol"{out.set(name,tsBooleanValue(number!=0))}else{out.set(name,tsNumberValue(float64(number)))};if name=="maxHeaderListSize"{out.set("maxHeaderSize",tsNumberValue(float64(number)))}}else{custom.set(strconv.Itoa(id),tsNumberValue(float64(number)));hasCustom=true}};if hasCustom{out.set("customSettings",tsObjectValue(custom))};if tsTruthy(tsFSProperty(options,"validate")){tsNodeHTTP2ValidateSettings(loop,tsObjectValue(out))};return tsObjectValue(out)}
func tsNodeHTTP2HeaderObject(loop *tsLoop,value tsValue)tsValue{out:=tsNewObject();switch value.kind{case tsArrayKind:array:=(*tsArray)(value.ref);for i:=0;i+1<len(array.values);i+=2{name:=strings.ToLower(tsText(array.values[i]));item:=array.values[i+1];existing:=out.get(name);if tsIsUndefined(existing){out.set(name,item)}else if existing.kind==tsArrayKind{(*tsArray)(existing.ref).values=append((*tsArray)(existing.ref).values,item)}else{out.set(name,tsArrayValue(&tsArray{values:[]tsValue{existing,item}}))}};case tsObjectKind,tsInstanceKind:source:=(*tsObject)(value.ref);for _,name:=range source.order{out.set(strings.ToLower(name),source.values[name])}};if tsNodeHTTP2Singleton.sensitiveHeaders!=nil{symbol:=tsValue{kind:tsSymbolKind,ref:unsafe.Pointer(tsNodeHTTP2Singleton.sensitiveHeaders)};sensitive:=tsGet(loop,value,symbol);if !tsIsUndefined(sensitive){tsSet(loop,tsObjectValue(out),symbol,sensitive)}};return tsObjectValue(out)}
func tsNodeHTTP2ValidateHeaders(loop *tsLoop,headers tsValue){tsNodeHTTP2ValidateHeadersStrict(loop,headers,true)}
func tsNodeHTTP2ValidateHeadersStrict(loop *tsLoop,headers tsValue,strictSingle bool){_ = loop;if !tsNodeHTTP2Object(headers){panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_TYPE","The \"headers\" argument must be of type object")})};object:=(*tsObject)(headers.ref);seenRegular:=false;seenPseudo:=map[string]bool{};allowedPseudo:=map[string]bool{":method":true,":scheme":true,":authority":true,":path":true,":status":true,":protocol":true};for _,name:=range object.order{if name==""||tsHTTP2IsMalformedHeaderName([]byte(name)){panic(tsThrown{tsFSException("TypeError","ERR_HTTP2_INVALID_HEADER_NAME","Invalid HTTP/2 header name")})};if strings.ToLower(name)!=name{panic(tsThrown{tsFSException("TypeError","ERR_HTTP2_INVALID_HEADER_NAME","Header names must be lowercase")})};if name=="connection"||name=="proxy-connection"||name=="keep-alive"||name=="transfer-encoding"||name=="upgrade"{panic(tsThrown{tsFSException("TypeError","ERR_HTTP2_INVALID_CONNECTION_HEADERS","HTTP/1 connection-specific headers are forbidden")})};if strings.HasPrefix(name,":"){if !allowedPseudo[name]||seenRegular||seenPseudo[name]{panic(tsThrown{tsFSException("TypeError","ERR_HTTP2_INVALID_PSEUDOHEADER","Invalid HTTP/2 pseudo-header")})};seenPseudo[name]=true}else{seenRegular=true};value:=object.values[name];values:=tsNodeHTTP2Values(value);if strictSingle&&tsHTTP2SingleValueHeaders[name]&&len(values)>1{panic(tsThrown{tsFSException("TypeError","ERR_HTTP2_HEADER_SINGLE_VALUE",fmt.Sprintf("Header field %q must only have a single value",name))})};for _,entry:=range values{if entry.kind!=tsStringKind&&entry.kind!=tsNumberKind&&entry.kind!=tsIntKind&&entry.kind!=tsBooleanKind{panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_TYPE","HTTP/2 header values must be string, number, or boolean")})};if tsHTTP2IsMalformedHeaderValue([]byte(tsText(entry))){panic(tsThrown{tsFSException("TypeError","ERR_INVALID_CHAR","Invalid character in header content")})}}}}
func tsNodeHTTP2ValidateHeadersFor(loop *tsLoop,owner,headers tsValue){strict:=true;if session:=tsNodeHTTP2State(owner);session!=nil{strict=session.strictSingleValueFields}else if stream:=tsNodeHTTP2Stream(owner);stream!=nil{if session:=tsNodeHTTP2State(stream.session);session!=nil{strict=session.strictSingleValueFields}}else if tsNodeHTTP2Object(owner){stream:=tsNodeHTTP2Get(loop,owner,"stream");if streamState:=tsNodeHTTP2Stream(stream);streamState!=nil{if session:=tsNodeHTTP2State(streamState.session);session!=nil{strict=session.strictSingleValueFields}}};tsNodeHTTP2ValidateHeadersStrict(loop,headers,strict)}
func tsNodeHTTP2Values(value tsValue)[]tsValue{if value.kind==tsArrayKind{return (*tsArray)(value.ref).values};return []tsValue{value}}
func tsNodeHTTP2HeaderArray(loop *tsLoop,headers tsValue)tsValue{out:=&tsArray{};if tsNodeHTTP2Object(headers){object:=(*tsObject)(headers.ref);for _,name:=range object.order{for _,value:=range tsNodeHTTP2Values(object.values[name]){out.values=append(out.values,tsNodeHTTP2Key(name),value)}}};return tsArrayValue(out)}
func tsNodeHTTP2BuildConstants(loop *tsLoop)tsValue{_ = loop;out:=tsNewObject();
out.set("NGHTTP2_ERR_FRAME_SIZE_ERROR",tsNumberValue(-522));
out.set("NGHTTP2_ERR_PROTO",tsNumberValue(-505));
out.set("NGHTTP2_ERR_STREAM_CLOSED",tsNumberValue(-510));
out.set("NGHTTP2_ERR_BAD_CLIENT_MAGIC",tsNumberValue(-903));
out.set("NGHTTP2_ERR_FLOODED",tsNumberValue(-904));
out.set("NGHTTP2_SESSION_SERVER",tsNumberValue(0));
out.set("NGHTTP2_SESSION_CLIENT",tsNumberValue(1));
out.set("NGHTTP2_STREAM_STATE_IDLE",tsNumberValue(1));
out.set("NGHTTP2_STREAM_STATE_OPEN",tsNumberValue(2));
out.set("NGHTTP2_STREAM_STATE_RESERVED_LOCAL",tsNumberValue(3));
out.set("NGHTTP2_STREAM_STATE_RESERVED_REMOTE",tsNumberValue(4));
out.set("NGHTTP2_STREAM_STATE_HALF_CLOSED_LOCAL",tsNumberValue(5));
out.set("NGHTTP2_STREAM_STATE_HALF_CLOSED_REMOTE",tsNumberValue(6));
out.set("NGHTTP2_STREAM_STATE_CLOSED",tsNumberValue(7));
out.set("NGHTTP2_FLAG_NONE",tsNumberValue(0));
out.set("NGHTTP2_FLAG_END_STREAM",tsNumberValue(1));
out.set("NGHTTP2_FLAG_END_HEADERS",tsNumberValue(4));
out.set("NGHTTP2_FLAG_ACK",tsNumberValue(1));
out.set("NGHTTP2_FLAG_PADDED",tsNumberValue(8));
out.set("NGHTTP2_FLAG_PRIORITY",tsNumberValue(32));
out.set("DEFAULT_SETTINGS_HEADER_TABLE_SIZE",tsNumberValue(4096));
out.set("DEFAULT_SETTINGS_ENABLE_PUSH",tsNumberValue(1));
out.set("DEFAULT_SETTINGS_MAX_CONCURRENT_STREAMS",tsNumberValue(4294967295));
out.set("DEFAULT_SETTINGS_INITIAL_WINDOW_SIZE",tsNumberValue(65535));
out.set("DEFAULT_SETTINGS_MAX_FRAME_SIZE",tsNumberValue(16384));
out.set("DEFAULT_SETTINGS_MAX_HEADER_LIST_SIZE",tsNumberValue(65535));
out.set("DEFAULT_SETTINGS_ENABLE_CONNECT_PROTOCOL",tsNumberValue(0));
out.set("MAX_MAX_FRAME_SIZE",tsNumberValue(16777215));
out.set("MIN_MAX_FRAME_SIZE",tsNumberValue(16384));
out.set("MAX_INITIAL_WINDOW_SIZE",tsNumberValue(2147483647));
out.set("NGHTTP2_SETTINGS_HEADER_TABLE_SIZE",tsNumberValue(1));
out.set("NGHTTP2_SETTINGS_ENABLE_PUSH",tsNumberValue(2));
out.set("NGHTTP2_SETTINGS_MAX_CONCURRENT_STREAMS",tsNumberValue(3));
out.set("NGHTTP2_SETTINGS_INITIAL_WINDOW_SIZE",tsNumberValue(4));
out.set("NGHTTP2_SETTINGS_MAX_FRAME_SIZE",tsNumberValue(5));
out.set("NGHTTP2_SETTINGS_MAX_HEADER_LIST_SIZE",tsNumberValue(6));
out.set("NGHTTP2_SETTINGS_ENABLE_CONNECT_PROTOCOL",tsNumberValue(8));
out.set("PADDING_STRATEGY_NONE",tsNumberValue(0));
out.set("PADDING_STRATEGY_ALIGNED",tsNumberValue(1));
out.set("PADDING_STRATEGY_MAX",tsNumberValue(2));
out.set("PADDING_STRATEGY_CALLBACK",tsNumberValue(1));
out.set("NGHTTP2_NO_ERROR",tsNumberValue(0));
out.set("NGHTTP2_PROTOCOL_ERROR",tsNumberValue(1));
out.set("NGHTTP2_INTERNAL_ERROR",tsNumberValue(2));
out.set("NGHTTP2_FLOW_CONTROL_ERROR",tsNumberValue(3));
out.set("NGHTTP2_SETTINGS_TIMEOUT",tsNumberValue(4));
out.set("NGHTTP2_STREAM_CLOSED",tsNumberValue(5));
out.set("NGHTTP2_FRAME_SIZE_ERROR",tsNumberValue(6));
out.set("NGHTTP2_REFUSED_STREAM",tsNumberValue(7));
out.set("NGHTTP2_CANCEL",tsNumberValue(8));
out.set("NGHTTP2_COMPRESSION_ERROR",tsNumberValue(9));
out.set("NGHTTP2_CONNECT_ERROR",tsNumberValue(10));
out.set("NGHTTP2_ENHANCE_YOUR_CALM",tsNumberValue(11));
out.set("NGHTTP2_INADEQUATE_SECURITY",tsNumberValue(12));
out.set("NGHTTP2_HTTP_1_1_REQUIRED",tsNumberValue(13));
out.set("NGHTTP2_DEFAULT_WEIGHT",tsNumberValue(16));
out.set("HTTP2_HEADER_STATUS",tsNodeHTTP2Key(":status"));
out.set("HTTP2_HEADER_METHOD",tsNodeHTTP2Key(":method"));
out.set("HTTP2_HEADER_AUTHORITY",tsNodeHTTP2Key(":authority"));
out.set("HTTP2_HEADER_SCHEME",tsNodeHTTP2Key(":scheme"));
out.set("HTTP2_HEADER_PATH",tsNodeHTTP2Key(":path"));
out.set("HTTP2_HEADER_PROTOCOL",tsNodeHTTP2Key(":protocol"));
out.set("HTTP2_HEADER_ACCEPT_ENCODING",tsNodeHTTP2Key("accept-encoding"));
out.set("HTTP2_HEADER_ACCEPT_LANGUAGE",tsNodeHTTP2Key("accept-language"));
out.set("HTTP2_HEADER_ACCEPT_RANGES",tsNodeHTTP2Key("accept-ranges"));
out.set("HTTP2_HEADER_ACCEPT",tsNodeHTTP2Key("accept"));
out.set("HTTP2_HEADER_ACCESS_CONTROL_ALLOW_CREDENTIALS",tsNodeHTTP2Key("access-control-allow-credentials"));
out.set("HTTP2_HEADER_ACCESS_CONTROL_ALLOW_HEADERS",tsNodeHTTP2Key("access-control-allow-headers"));
out.set("HTTP2_HEADER_ACCESS_CONTROL_ALLOW_METHODS",tsNodeHTTP2Key("access-control-allow-methods"));
out.set("HTTP2_HEADER_ACCESS_CONTROL_ALLOW_ORIGIN",tsNodeHTTP2Key("access-control-allow-origin"));
out.set("HTTP2_HEADER_ACCESS_CONTROL_EXPOSE_HEADERS",tsNodeHTTP2Key("access-control-expose-headers"));
out.set("HTTP2_HEADER_ACCESS_CONTROL_REQUEST_HEADERS",tsNodeHTTP2Key("access-control-request-headers"));
out.set("HTTP2_HEADER_ACCESS_CONTROL_REQUEST_METHOD",tsNodeHTTP2Key("access-control-request-method"));
out.set("HTTP2_HEADER_AGE",tsNodeHTTP2Key("age"));
out.set("HTTP2_HEADER_AUTHORIZATION",tsNodeHTTP2Key("authorization"));
out.set("HTTP2_HEADER_CACHE_CONTROL",tsNodeHTTP2Key("cache-control"));
out.set("HTTP2_HEADER_CONNECTION",tsNodeHTTP2Key("connection"));
out.set("HTTP2_HEADER_CONTENT_DISPOSITION",tsNodeHTTP2Key("content-disposition"));
out.set("HTTP2_HEADER_CONTENT_ENCODING",tsNodeHTTP2Key("content-encoding"));
out.set("HTTP2_HEADER_CONTENT_LENGTH",tsNodeHTTP2Key("content-length"));
out.set("HTTP2_HEADER_CONTENT_TYPE",tsNodeHTTP2Key("content-type"));
out.set("HTTP2_HEADER_COOKIE",tsNodeHTTP2Key("cookie"));
out.set("HTTP2_HEADER_DATE",tsNodeHTTP2Key("date"));
out.set("HTTP2_HEADER_ETAG",tsNodeHTTP2Key("etag"));
out.set("HTTP2_HEADER_FORWARDED",tsNodeHTTP2Key("forwarded"));
out.set("HTTP2_HEADER_HOST",tsNodeHTTP2Key("host"));
out.set("HTTP2_HEADER_IF_MODIFIED_SINCE",tsNodeHTTP2Key("if-modified-since"));
out.set("HTTP2_HEADER_IF_NONE_MATCH",tsNodeHTTP2Key("if-none-match"));
out.set("HTTP2_HEADER_IF_RANGE",tsNodeHTTP2Key("if-range"));
out.set("HTTP2_HEADER_LAST_MODIFIED",tsNodeHTTP2Key("last-modified"));
out.set("HTTP2_HEADER_LINK",tsNodeHTTP2Key("link"));
out.set("HTTP2_HEADER_LOCATION",tsNodeHTTP2Key("location"));
out.set("HTTP2_HEADER_RANGE",tsNodeHTTP2Key("range"));
out.set("HTTP2_HEADER_REFERER",tsNodeHTTP2Key("referer"));
out.set("HTTP2_HEADER_SERVER",tsNodeHTTP2Key("server"));
out.set("HTTP2_HEADER_SET_COOKIE",tsNodeHTTP2Key("set-cookie"));
out.set("HTTP2_HEADER_STRICT_TRANSPORT_SECURITY",tsNodeHTTP2Key("strict-transport-security"));
out.set("HTTP2_HEADER_TRANSFER_ENCODING",tsNodeHTTP2Key("transfer-encoding"));
out.set("HTTP2_HEADER_TE",tsNodeHTTP2Key("te"));
out.set("HTTP2_HEADER_UPGRADE_INSECURE_REQUESTS",tsNodeHTTP2Key("upgrade-insecure-requests"));
out.set("HTTP2_HEADER_UPGRADE",tsNodeHTTP2Key("upgrade"));
out.set("HTTP2_HEADER_USER_AGENT",tsNodeHTTP2Key("user-agent"));
out.set("HTTP2_HEADER_VARY",tsNodeHTTP2Key("vary"));
out.set("HTTP2_HEADER_X_CONTENT_TYPE_OPTIONS",tsNodeHTTP2Key("x-content-type-options"));
out.set("HTTP2_HEADER_X_FRAME_OPTIONS",tsNodeHTTP2Key("x-frame-options"));
out.set("HTTP2_HEADER_KEEP_ALIVE",tsNodeHTTP2Key("keep-alive"));
out.set("HTTP2_HEADER_PROXY_CONNECTION",tsNodeHTTP2Key("proxy-connection"));
out.set("HTTP2_HEADER_X_XSS_PROTECTION",tsNodeHTTP2Key("x-xss-protection"));
out.set("HTTP2_HEADER_ALT_SVC",tsNodeHTTP2Key("alt-svc"));
out.set("HTTP2_HEADER_CONTENT_SECURITY_POLICY",tsNodeHTTP2Key("content-security-policy"));
out.set("HTTP2_HEADER_EARLY_DATA",tsNodeHTTP2Key("early-data"));
out.set("HTTP2_HEADER_EXPECT_CT",tsNodeHTTP2Key("expect-ct"));
out.set("HTTP2_HEADER_ORIGIN",tsNodeHTTP2Key("origin"));
out.set("HTTP2_HEADER_PURPOSE",tsNodeHTTP2Key("purpose"));
out.set("HTTP2_HEADER_TIMING_ALLOW_ORIGIN",tsNodeHTTP2Key("timing-allow-origin"));
out.set("HTTP2_HEADER_X_FORWARDED_FOR",tsNodeHTTP2Key("x-forwarded-for"));
out.set("HTTP2_HEADER_PRIORITY",tsNodeHTTP2Key("priority"));
out.set("HTTP2_HEADER_ACCEPT_CHARSET",tsNodeHTTP2Key("accept-charset"));
out.set("HTTP2_HEADER_ACCESS_CONTROL_MAX_AGE",tsNodeHTTP2Key("access-control-max-age"));
out.set("HTTP2_HEADER_ALLOW",tsNodeHTTP2Key("allow"));
out.set("HTTP2_HEADER_CONTENT_LANGUAGE",tsNodeHTTP2Key("content-language"));
out.set("HTTP2_HEADER_CONTENT_LOCATION",tsNodeHTTP2Key("content-location"));
out.set("HTTP2_HEADER_CONTENT_MD5",tsNodeHTTP2Key("content-md5"));
out.set("HTTP2_HEADER_CONTENT_RANGE",tsNodeHTTP2Key("content-range"));
out.set("HTTP2_HEADER_DNT",tsNodeHTTP2Key("dnt"));
out.set("HTTP2_HEADER_EXPECT",tsNodeHTTP2Key("expect"));
out.set("HTTP2_HEADER_EXPIRES",tsNodeHTTP2Key("expires"));
out.set("HTTP2_HEADER_FROM",tsNodeHTTP2Key("from"));
out.set("HTTP2_HEADER_IF_MATCH",tsNodeHTTP2Key("if-match"));
out.set("HTTP2_HEADER_IF_UNMODIFIED_SINCE",tsNodeHTTP2Key("if-unmodified-since"));
out.set("HTTP2_HEADER_MAX_FORWARDS",tsNodeHTTP2Key("max-forwards"));
out.set("HTTP2_HEADER_PREFER",tsNodeHTTP2Key("prefer"));
out.set("HTTP2_HEADER_PROXY_AUTHENTICATE",tsNodeHTTP2Key("proxy-authenticate"));
out.set("HTTP2_HEADER_PROXY_AUTHORIZATION",tsNodeHTTP2Key("proxy-authorization"));
out.set("HTTP2_HEADER_REFRESH",tsNodeHTTP2Key("refresh"));
out.set("HTTP2_HEADER_RETRY_AFTER",tsNodeHTTP2Key("retry-after"));
out.set("HTTP2_HEADER_TRAILER",tsNodeHTTP2Key("trailer"));
out.set("HTTP2_HEADER_TK",tsNodeHTTP2Key("tk"));
out.set("HTTP2_HEADER_VIA",tsNodeHTTP2Key("via"));
out.set("HTTP2_HEADER_WARNING",tsNodeHTTP2Key("warning"));
out.set("HTTP2_HEADER_WWW_AUTHENTICATE",tsNodeHTTP2Key("www-authenticate"));
out.set("HTTP2_HEADER_HTTP2_SETTINGS",tsNodeHTTP2Key("http2-settings"));
out.set("HTTP2_METHOD_ACL",tsNodeHTTP2Key("ACL"));
out.set("HTTP2_METHOD_BASELINE_CONTROL",tsNodeHTTP2Key("BASELINE-CONTROL"));
out.set("HTTP2_METHOD_BIND",tsNodeHTTP2Key("BIND"));
out.set("HTTP2_METHOD_CHECKIN",tsNodeHTTP2Key("CHECKIN"));
out.set("HTTP2_METHOD_CHECKOUT",tsNodeHTTP2Key("CHECKOUT"));
out.set("HTTP2_METHOD_CONNECT",tsNodeHTTP2Key("CONNECT"));
out.set("HTTP2_METHOD_COPY",tsNodeHTTP2Key("COPY"));
out.set("HTTP2_METHOD_DELETE",tsNodeHTTP2Key("DELETE"));
out.set("HTTP2_METHOD_GET",tsNodeHTTP2Key("GET"));
out.set("HTTP2_METHOD_HEAD",tsNodeHTTP2Key("HEAD"));
out.set("HTTP2_METHOD_LABEL",tsNodeHTTP2Key("LABEL"));
out.set("HTTP2_METHOD_LINK",tsNodeHTTP2Key("LINK"));
out.set("HTTP2_METHOD_LOCK",tsNodeHTTP2Key("LOCK"));
out.set("HTTP2_METHOD_MERGE",tsNodeHTTP2Key("MERGE"));
out.set("HTTP2_METHOD_MKACTIVITY",tsNodeHTTP2Key("MKACTIVITY"));
out.set("HTTP2_METHOD_MKCALENDAR",tsNodeHTTP2Key("MKCALENDAR"));
out.set("HTTP2_METHOD_MKCOL",tsNodeHTTP2Key("MKCOL"));
out.set("HTTP2_METHOD_MKREDIRECTREF",tsNodeHTTP2Key("MKREDIRECTREF"));
out.set("HTTP2_METHOD_MKWORKSPACE",tsNodeHTTP2Key("MKWORKSPACE"));
out.set("HTTP2_METHOD_MOVE",tsNodeHTTP2Key("MOVE"));
out.set("HTTP2_METHOD_OPTIONS",tsNodeHTTP2Key("OPTIONS"));
out.set("HTTP2_METHOD_ORDERPATCH",tsNodeHTTP2Key("ORDERPATCH"));
out.set("HTTP2_METHOD_PATCH",tsNodeHTTP2Key("PATCH"));
out.set("HTTP2_METHOD_POST",tsNodeHTTP2Key("POST"));
out.set("HTTP2_METHOD_PRI",tsNodeHTTP2Key("PRI"));
out.set("HTTP2_METHOD_PROPFIND",tsNodeHTTP2Key("PROPFIND"));
out.set("HTTP2_METHOD_PROPPATCH",tsNodeHTTP2Key("PROPPATCH"));
out.set("HTTP2_METHOD_PUT",tsNodeHTTP2Key("PUT"));
out.set("HTTP2_METHOD_REBIND",tsNodeHTTP2Key("REBIND"));
out.set("HTTP2_METHOD_REPORT",tsNodeHTTP2Key("REPORT"));
out.set("HTTP2_METHOD_SEARCH",tsNodeHTTP2Key("SEARCH"));
out.set("HTTP2_METHOD_TRACE",tsNodeHTTP2Key("TRACE"));
out.set("HTTP2_METHOD_UNBIND",tsNodeHTTP2Key("UNBIND"));
out.set("HTTP2_METHOD_UNCHECKOUT",tsNodeHTTP2Key("UNCHECKOUT"));
out.set("HTTP2_METHOD_UNLINK",tsNodeHTTP2Key("UNLINK"));
out.set("HTTP2_METHOD_UNLOCK",tsNodeHTTP2Key("UNLOCK"));
out.set("HTTP2_METHOD_UPDATE",tsNodeHTTP2Key("UPDATE"));
out.set("HTTP2_METHOD_UPDATEREDIRECTREF",tsNodeHTTP2Key("UPDATEREDIRECTREF"));
out.set("HTTP2_METHOD_VERSION_CONTROL",tsNodeHTTP2Key("VERSION-CONTROL"));
out.set("HTTP_STATUS_CONTINUE",tsNumberValue(100));
out.set("HTTP_STATUS_SWITCHING_PROTOCOLS",tsNumberValue(101));
out.set("HTTP_STATUS_PROCESSING",tsNumberValue(102));
out.set("HTTP_STATUS_EARLY_HINTS",tsNumberValue(103));
out.set("HTTP_STATUS_OK",tsNumberValue(200));
out.set("HTTP_STATUS_CREATED",tsNumberValue(201));
out.set("HTTP_STATUS_ACCEPTED",tsNumberValue(202));
out.set("HTTP_STATUS_NON_AUTHORITATIVE_INFORMATION",tsNumberValue(203));
out.set("HTTP_STATUS_NO_CONTENT",tsNumberValue(204));
out.set("HTTP_STATUS_RESET_CONTENT",tsNumberValue(205));
out.set("HTTP_STATUS_PARTIAL_CONTENT",tsNumberValue(206));
out.set("HTTP_STATUS_MULTI_STATUS",tsNumberValue(207));
out.set("HTTP_STATUS_ALREADY_REPORTED",tsNumberValue(208));
out.set("HTTP_STATUS_IM_USED",tsNumberValue(226));
out.set("HTTP_STATUS_MULTIPLE_CHOICES",tsNumberValue(300));
out.set("HTTP_STATUS_MOVED_PERMANENTLY",tsNumberValue(301));
out.set("HTTP_STATUS_FOUND",tsNumberValue(302));
out.set("HTTP_STATUS_SEE_OTHER",tsNumberValue(303));
out.set("HTTP_STATUS_NOT_MODIFIED",tsNumberValue(304));
out.set("HTTP_STATUS_USE_PROXY",tsNumberValue(305));
out.set("HTTP_STATUS_TEMPORARY_REDIRECT",tsNumberValue(307));
out.set("HTTP_STATUS_PERMANENT_REDIRECT",tsNumberValue(308));
out.set("HTTP_STATUS_BAD_REQUEST",tsNumberValue(400));
out.set("HTTP_STATUS_UNAUTHORIZED",tsNumberValue(401));
out.set("HTTP_STATUS_PAYMENT_REQUIRED",tsNumberValue(402));
out.set("HTTP_STATUS_FORBIDDEN",tsNumberValue(403));
out.set("HTTP_STATUS_NOT_FOUND",tsNumberValue(404));
out.set("HTTP_STATUS_METHOD_NOT_ALLOWED",tsNumberValue(405));
out.set("HTTP_STATUS_NOT_ACCEPTABLE",tsNumberValue(406));
out.set("HTTP_STATUS_PROXY_AUTHENTICATION_REQUIRED",tsNumberValue(407));
out.set("HTTP_STATUS_REQUEST_TIMEOUT",tsNumberValue(408));
out.set("HTTP_STATUS_CONFLICT",tsNumberValue(409));
out.set("HTTP_STATUS_GONE",tsNumberValue(410));
out.set("HTTP_STATUS_LENGTH_REQUIRED",tsNumberValue(411));
out.set("HTTP_STATUS_PRECONDITION_FAILED",tsNumberValue(412));
out.set("HTTP_STATUS_PAYLOAD_TOO_LARGE",tsNumberValue(413));
out.set("HTTP_STATUS_URI_TOO_LONG",tsNumberValue(414));
out.set("HTTP_STATUS_UNSUPPORTED_MEDIA_TYPE",tsNumberValue(415));
out.set("HTTP_STATUS_RANGE_NOT_SATISFIABLE",tsNumberValue(416));
out.set("HTTP_STATUS_EXPECTATION_FAILED",tsNumberValue(417));
out.set("HTTP_STATUS_TEAPOT",tsNumberValue(418));
out.set("HTTP_STATUS_MISDIRECTED_REQUEST",tsNumberValue(421));
out.set("HTTP_STATUS_UNPROCESSABLE_ENTITY",tsNumberValue(422));
out.set("HTTP_STATUS_LOCKED",tsNumberValue(423));
out.set("HTTP_STATUS_FAILED_DEPENDENCY",tsNumberValue(424));
out.set("HTTP_STATUS_TOO_EARLY",tsNumberValue(425));
out.set("HTTP_STATUS_UPGRADE_REQUIRED",tsNumberValue(426));
out.set("HTTP_STATUS_PRECONDITION_REQUIRED",tsNumberValue(428));
out.set("HTTP_STATUS_TOO_MANY_REQUESTS",tsNumberValue(429));
out.set("HTTP_STATUS_REQUEST_HEADER_FIELDS_TOO_LARGE",tsNumberValue(431));
out.set("HTTP_STATUS_UNAVAILABLE_FOR_LEGAL_REASONS",tsNumberValue(451));
out.set("HTTP_STATUS_INTERNAL_SERVER_ERROR",tsNumberValue(500));
out.set("HTTP_STATUS_NOT_IMPLEMENTED",tsNumberValue(501));
out.set("HTTP_STATUS_BAD_GATEWAY",tsNumberValue(502));
out.set("HTTP_STATUS_SERVICE_UNAVAILABLE",tsNumberValue(503));
out.set("HTTP_STATUS_GATEWAY_TIMEOUT",tsNumberValue(504));
out.set("HTTP_STATUS_HTTP_VERSION_NOT_SUPPORTED",tsNumberValue(505));
out.set("HTTP_STATUS_VARIANT_ALSO_NEGOTIATES",tsNumberValue(506));
out.set("HTTP_STATUS_INSUFFICIENT_STORAGE",tsNumberValue(507));
out.set("HTTP_STATUS_LOOP_DETECTED",tsNumberValue(508));
out.set("HTTP_STATUS_BANDWIDTH_LIMIT_EXCEEDED",tsNumberValue(509));
out.set("HTTP_STATUS_NOT_EXTENDED",tsNumberValue(510));
out.set("HTTP_STATUS_NETWORK_AUTHENTICATION_REQUIRED",tsNumberValue(511));
return tsObjectValue(out)}

func tsNodeHTTP2SessionNew(loop *tsLoop,server bool,socket,options tsValue)tsValue{
 object:=tsNewObject();object.prototype=tsObjectValue(tsNodeHTTP2Singleton.sessionPrototype);value:=tsObjectValue(object);tsNodeEventCreate(loop,value,tsU)
 local:=tsNodeHTTP2SettingsDefaults();submitted:=tsFSProperty(options,"settings");if !tsIsUndefined(submitted){tsNodeHTTP2ValidateSettings(loop,submitted);copy:=tsNewObject();for _,name:=range (*tsObject)(local.ref).order{copy.set(name,(*tsObject)(local.ref).values[name])};if tsNodeHTTP2Object(submitted){source:=(*tsObject)(submitted.ref);for _,name:=range source.order{copy.set(name,source.values[name])}};local=tsObjectValue(copy)}
 state:=&tsNodeHTTP2SessionState{object:value,server:server,socket:socket,localSettings:local,remoteSettings:tsNodeHTTP2SettingsDefaults(),streams:map[int]tsValue{},pushParents:map[int]int{},nextStreamID:1,encrypted:tsTruthy(tsNodeHTTP2Get(loop,socket,"encrypted")),sendWindow:tsHTTP2NewSendWindow(65535),recvWindow:tsHTTP2NewRecvWindow(65535),strictSingleValueFields:tsNodeHTTP2StrictSingleValueEnabled(loop,options),outstandingSettings:1,engine:tsHTTP2NewConnection(server,tsHTTP2SettingsFromValue(local)),remoteCustomFilter:map[uint16]bool{},remoteCustomValues:map[uint16]uint32{},localCustomSettings:tsNodeHTTP2EmptyObject()};tsNodeHTTP2ConfigureSessionOptions(loop,state,submitted,options);state.initialSettings=submitted;if tsIsUndefined(state.initialSettings)||tsIsNull(state.initialSettings){state.initialSettings=tsNodeHTTP2EmptyObject()};if server{state.nextStreamID=2};if custom:=tsFSProperty(local,"customSettings");tsNodeHTTP2Object(custom){state.localCustomSettings=custom};if custom:=tsFSProperty(options,"remoteCustomSettings");custom.kind==tsArrayKind{for _,item:=range (*tsArray)(custom.ref).values{if !tsIsNumeric(item){continue};id:=tsNumber(item);if id<0||id>65535||id!=math.Trunc(id){continue};if len(state.remoteCustomFilter)<10{state.remoteCustomFilter[uint16(id)]=true}}};tsNodeHTTP2Singleton.sessions[uintptr(unsafe.Pointer(object))]=state
 for name,item:=range map[string]tsValue{"socket":socket,"connecting":tsBooleanValue(!server&&tsIsNull(socket)),"destroyed":tsBooleanValue(false),"closed":tsBooleanValue(false),"alpnProtocol":tsNodeHTTP2Key("h2"),"localSettings":local,"remoteSettings":state.remoteSettings,"type":tsNodeHTTP2Key(func()string{if server{return "server"};return "client"}()),"originSet":tsArrayValue(&tsArray{}),"state":tsNodeHTTP2EmptyObject()}{tsNodeHTTP2Set(loop,value,name,item)}
 tsNodeHTTP2AttachTransport(loop,state)
 tsNodeHTTP2RefreshPublicState(loop,state,state.engine.framesReceived,state.engine.framesSent)
 return value
}
func tsNodeHTTP2StreamNew(loop *tsLoop,session tsValue,id int,headers tsValue)tsValue{
 options:=tsNewObject();var stream tsValue;options.set("read",tsNamedNativeMethod("_read",1,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsU}));options.set("write",tsNamedNativeMethod("_write",3,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeHTTP2StreamWrite(loop,stream,args...)}));options.set("final",tsNamedNativeMethod("_final",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeHTTP2StreamFinal(loop,stream,args...)}));stream=tsNodeDuplexCreate(loop,tsU,tsObjectValue(options));object:=(*tsObject)(stream.ref);object.prototype=tsObjectValue(tsNodeHTTP2Singleton.streamPrototype);state:=&tsNodeHTTP2StreamState{object:stream,session:session,id:id,headers:headers,rawHeaders:tsNodeHTTP2HeaderArray(loop,headers),trailers:tsNodeHTTP2EmptyObject(),sendWindow:tsHTTP2NewSendWindow(65535),recvWindow:tsHTTP2NewRecvWindow(65535),protocolState:tsHTTP2StreamOpen,weight:16,paddingStrategy:func()uint8{if state:=tsNodeHTTP2State(session);state!=nil{return state.paddingStrategy};return 0}()};tsNodeHTTP2Singleton.streams[uintptr(unsafe.Pointer(object))]=state
 for name,value:=range map[string]tsValue{"id":tsNumberValue(float64(id)),"session":session,"sentHeaders":tsNodeHTTP2EmptyObject(),"pending":tsBooleanValue(false),"destroyed":tsBooleanValue(false),"rstCode":tsNumberValue(0),"headersSent":tsBooleanValue(false),"aborted":tsBooleanValue(false),"bufferSize":tsNumberValue(0),"state":tsNodeHTTP2EmptyObject(),"endAfterHeaders":tsBooleanValue(false),"headRequest":tsBooleanValue(false),"originSet":tsArrayValue(&tsArray{})}{tsNodeHTTP2Set(loop,stream,name,value)}
 tsNodeHTTP2AttachStreamReadControl(loop,stream)
 return stream
}
func tsNodeHTTP2SyncStreamLifecycle(stream tsValue,state *tsNodeHTTP2StreamState){if state==nil{return};if readable:=tsNodeReadableStateFor(stream);readable!=nil{readable.destroyed=state.destroyed;readable.closed=state.closed;tsNodeReadableSyncState(readable)};if writable:=tsNodeWritableStateFor(stream);writable!=nil{writable.destroyed=state.destroyed;writable.closed=state.closed;tsNodeWritableSync(writable)}}
func tsNodeHTTP2SessionSetState(loop *tsLoop,state *tsNodeHTTP2SessionState,name string,value bool){if state==nil{return};tsNodeHTTP2Set(loop,state.object,name,tsBooleanValue(value));switch name{case "closed":state.closed=value;case "destroyed":state.destroyed=value;case "connecting":state.connected=!value}}
func tsNodeHTTP2SessionClose(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeHTTP2State(receiver);callback:=tsArg(args,0);if state==nil||state.closing||state.closed||state.destroyed{return tsU};state.closing=true;tsNodeHTTP2Set(loop,receiver,"closed",tsBooleanValue(true));if callback.kind==tsFunctionKind{tsCallReceiver(loop,tsNodeHTTP2Get(loop,receiver,"once"),receiver,tsNodeHTTP2Key("close"),callback)};tsNodeHTTP2SessionGoaway(loop,receiver,tsNumberValue(0),tsNumberValue(0),tsNodeBuffer(nil));tsNodeHTTP2MaybeFinishSessionClose(loop,state);return tsU}
func tsNodeHTTP2FailStream(loop *tsLoop,session *tsNodeHTTP2SessionState,stream tsValue,code uint32){streamState:=tsNodeHTTP2Stream(stream);if streamState==nil{return};streamState.rstCode=int(code);if session!=nil&&session.engine!=nil{if protocol:=session.engine.streams[uint32(streamState.id)];protocol!=nil{protocol.state=tsHTTP2StreamClosed}};loop.post(func(){current:=tsNodeHTTP2Stream(stream);if current==nil||current.closed{return};if code!=tsHTTP2ErrorNoError&&code!=tsHTTP2ErrorCancel{tsNodeNetEmitSafe(loop,stream,"error",tsNodeHTTP2StreamResetError(code))};tsNodeHTTP2StreamClose(loop,stream,tsNumberValue(float64(code)));if session!=nil{tsNodeHTTP2MaybeFinishSessionClose(loop,session)}})}
func tsNodeHTTP2RequestOptionType(name,expected string){panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_TYPE","The \"options."+name+"\" argument must be of type "+expected)})}
func tsNodeHTTP2PendingStreamCancelError(loop *tsLoop,cause tsValue)tsValue{message:="The pending stream has been canceled";err:=tsFSException("Error","ERR_HTTP2_STREAM_CANCEL",message);if cause.kind==tsObjectKind||cause.kind==tsInstanceKind{if causeMessage:=tsNodeHTTP2Get(loop,cause,"message");causeMessage.kind==tsStringKind{message+=" (caused by: "+tsText(causeMessage)+")";tsNodeHTTP2Set(loop,err,"message",tsStringReference(tsStringUTF8(message)))};if err.kind==tsErrorKind{(*tsRuntimeError)(err.ref).fields["cause"]=cause}};return err}
func tsNodeHTTP2SessionDestroy(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeHTTP2State(receiver);errorValue,code:=tsArg(args,0),tsArg(args,1);if state==nil{return tsU};state.destroyed=true;state.closed=true;tsNodeHTTP2Set(loop,receiver,"destroyed",tsBooleanValue(true));tsNodeHTTP2Set(loop,receiver,"closed",tsBooleanValue(true));tsNodeHTTP2DetachTransport(loop,state);streamError:=errorValue;if tsNullish(streamError)||tsIsUndefined(streamError){streamError=tsNodeHTTP2PendingStreamCancelError(loop,tsU)};for _,id:=range tsNodeHTTP2StreamSnapshotIDs(state){stream,exists:=state.streams[id];if !exists{continue};if streamState:=tsNodeHTTP2Stream(stream);streamState!=nil&&!streamState.closed{tsNodeNetEmitSafe(loop,stream,"error",streamError);tsNodeHTTP2StreamClose(loop,stream,tsNumberValue(float64(tsNumber(code))))}};if !tsNullish(errorValue)&&!tsIsUndefined(errorValue){tsNodeNetEmit(loop,receiver,"error",errorValue)};if !tsNullish(state.socket){if method:=tsNodeHTTP2Get(loop,state.socket,"destroy");method.kind==tsFunctionKind{tsCallReceiver(loop,method,state.socket,errorValue)}};tsNodeHTTP2EmitSessionClose(loop,state);return tsU}
func tsNodeHTTP2SessionPing(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeHTTP2State(receiver);callback,payload:=tsArg(args,0),tsArg(args,1);if callback.kind!=tsFunctionKind&&tsArg(args,1).kind==tsFunctionKind{callback=tsArg(args,1);payload=tsArg(args,0)};if !tsIsUndefined(payload)&&!tsIsNull(payload)&&payload.kind!=tsTypedArrayKind&&payload.kind!=tsArrayBufferKind{panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_TYPE","The \"payload\" argument must be a Buffer or TypedArray")})};data:=append([]byte(nil),tsFSBytes(payload)...);if len(data)==0{data=make([]byte,8);binary.BigEndian.PutUint64(data,uint64(time.Now().UnixNano()))};if len(data)!=8{panic(tsThrown{tsFSException("RangeError","ERR_HTTP2_PING_LENGTH","The ping payload must be exactly 8 bytes")})};if state==nil||state.closed||state.closing||state.engine==nil{if callback.kind==tsFunctionKind{loop.post(func(){tsCall(loop,callback,tsFSException("Error","ERR_HTTP2_INVALID_SESSION","Session is closed"))})};return tsBooleanValue(false)};if state.outstandingPings>=int(state.maxOutstandingPings){if callback.kind==tsFunctionKind{loop.post(func(){tsCall(loop,callback,tsFSException("Error","ERR_HTTP2_PING_CANCEL","HTTP/2 ping was canceled"))})};return tsBooleanValue(false)};state.outstandingPings++;key:=string(data);if callback.kind==tsFunctionKind{if state.pingCallbacks==nil{state.pingCallbacks=map[string][]tsValue{}};state.pingCallbacks[key]=append(state.pingCallbacks[key],callback)};state.engine.writeFrame(tsNodeHTTP2EngineSinkFor(loop,state),6,0,0,data);return tsBooleanValue(true)}
func tsNodeHTTP2SessionSettings(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeHTTP2State(receiver);settings,callback:=tsArg(args,0),tsArg(args,1);if state==nil{return tsU};tsNodeHTTP2ValidateSettings(loop,settings);if state.outstandingSettings>=state.maxOutstandingSettings{err:=tsFSException("Error","ERR_HTTP2_MAX_PENDING_SETTINGS_ACK","Maximum number of pending settings acknowledgements");tsNodeHTTP2SessionDestroy(loop,receiver,err,tsNumberValue(float64(tsHTTP2ErrorInternal)));return tsU};copy:=tsObjectValue(tsObjectSpreadCopy((*tsObject)(state.localSettings.ref)));if tsNodeHTTP2Object(settings){source:=(*tsObject)(settings.ref);for _,name:=range source.order{value:=source.values[name];if tsIsNull(value){switch name{case "headerTableSize","maxConcurrentStreams","initialWindowSize","maxFrameSize","maxHeaderListSize","maxHeaderSize":continue}};tsNodeHTTP2Set(loop,copy,name,value)};custom:=tsFSProperty(settings,"customSettings");if tsNodeHTTP2Object(custom){merged:=tsObjectValue(tsObjectSpreadCopy((*tsObject)(state.localCustomSettings.ref)));for _,name:=range (*tsObject)(custom.ref).order{tsNodeHTTP2Set(loop,merged,name,(*tsObject)(custom.ref).values[name])};state.localCustomSettings=merged;tsNodeHTTP2Set(loop,copy,"customSettings",merged)}};state.localSettings=copy;tsNodeHTTP2Set(loop,receiver,"localSettings",copy);if state.engine!=nil{wire:=tsHTTP2SettingsFromValue(copy);state.engine.localSettings=wire;state.engine.pendingLocalSettingsAcks=append(state.engine.pendingLocalSettingsAcks,tsHTTP2PendingLocalSettings{settings:wire});state.pendingSettingsCallbacks=append(state.pendingSettingsCallbacks,callback);state.outstandingSettings++;payload:=tsNodeHTTP2SettingsUpdatePayload(loop,settings);state.engine.writeFrame(tsNodeHTTP2EngineSinkFor(loop,state),4,0,0,payload)}else if callback.kind==tsFunctionKind{loop.post(func(){tsCall(loop,callback,tsNull)})};return tsU}
func tsNodeHTTP2SessionRequest(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeHTTP2State(receiver);headers,options:=tsArg(args,0),tsArg(args,1);if state==nil||state.server{panic(tsThrown{tsFSException("Error","ERR_HTTP2_INVALID_SESSION","Client request requires a client HTTP/2 session")})};if state.destroyed||state.closed||state.closing{panic(tsThrown{tsFSException("Error","ERR_HTTP2_INVALID_SESSION","Cannot make a request on a destroyed session")})};if !tsNullish(headers){headers=tsNodeHTTP2HeaderObject(loop,headers);headers=tsNodeHTTP2OrderRequestPseudoHeaders(loop,headers);tsNodeHTTP2ValidateHeadersFor(loop,receiver,headers)}else{headers=tsNodeHTTP2EmptyObject()};signal:=tsNodeHTTP2Get(loop,options,"signal");if !tsNullish(signal)&&tsTruthy(signal){tsInternalValidatorAbortSignal(signal,tsStringReference(tsStringUTF8("options.signal")));if tsTruthy(tsNodeHTTP2Get(loop,signal,"aborted")){stream:=tsNodeHTTP2StreamNew(loop,receiver,0,headers);tsNodeHTTP2Set(loop,stream,"id",tsU);tsNodeHTTP2AttachRequestSignal(loop,stream,signal);return stream}};id:=tsNodeHTTP2AllocateStreamID(state);if id<0{stream:=tsNodeHTTP2StreamNew(loop,receiver,0,headers);tsNodeHTTP2Set(loop,stream,"id",tsU);loop.post(func(){tsNodeNetEmitSafe(loop,stream,"error",tsFSException("Error","ERR_HTTP2_OUT_OF_STREAMS","No more stream IDs are available"))});return stream};stream:=tsNodeHTTP2StreamNew(loop,receiver,id,headers);if streamState:=tsNodeHTTP2Stream(stream);streamState!=nil{streamState.exposed=true};state.streams[id]=stream;if !tsNullish(options)&&!tsNodeHTTP2Object(options){tsNodeHTTP2FailStream(loop,state,stream,tsHTTP2ErrorInternal);return stream};if tsNodeHTTP2AttachRequestSignal(loop,stream,signal){return stream};priority:=tsNodeHTTP2EmptyObject();tsNodeHTTP2Set(loop,stream,"priority",priority);endStream:=false;if value:=tsNodeHTTP2Get(loop,options,"endStream");!tsIsUndefined(value){if value.kind!=tsBooleanKind{tsNodeHTTP2RequestOptionType("endStream","boolean")};endStream=tsTruthy(value)};waitForTrailers:=false;if value:=tsNodeHTTP2Get(loop,options,"waitForTrailers");value.kind==tsBooleanKind{waitForTrailers=tsTruthy(value)};if value:=tsNodeHTTP2Get(loop,options,"silent");!tsIsUndefined(value)&&value.kind!=tsBooleanKind{tsNodeHTTP2RequestOptionType("silent","boolean")};streamState:=tsNodeHTTP2Stream(stream);if streamState!=nil{streamState.waitForTrailers=waitForTrailers;streamState.paddingStrategy=state.paddingStrategy;if value:=tsNodeHTTP2Get(loop,options,"paddingStrategy");tsIsNumeric(value){switch uint32(tsNumber(value)){case 1:streamState.paddingStrategy=1;case 2:streamState.paddingStrategy=2;default:streamState.paddingStrategy=0}}};weight:=tsNodeHTTP2Get(loop,options,"weight");parent:=tsNodeHTTP2Get(loop,options,"parent");exclusive:=tsNodeHTTP2Get(loop,options,"exclusive");hasPriority:=false;if !tsIsUndefined(weight){if !tsIsNumeric(weight){tsNodeHTTP2RequestOptionType("weight","number")};weightNumber:=int32(tsNumber(weight));if weightNumber<1||weightNumber>255{tsNodeHTTP2FailStream(loop,state,stream,tsHTTP2ErrorInternal);return stream};weight=tsNumberValue(float64(weightNumber));hasPriority=true;if streamState!=nil{streamState.weight=uint16(weightNumber)}}else{weight=tsNumberValue(16)};if !tsIsUndefined(parent){if !tsIsNumeric(parent){tsNodeHTTP2RequestOptionType("parent","number")};parentNumber:=int32(tsNumber(parent));if parentNumber<=0||uint32(parentNumber)>tsHTTP2MaxStreamID{tsNodeHTTP2FailStream(loop,state,stream,tsHTTP2ErrorInternal);return stream};parent=tsNumberValue(float64(parentNumber));hasPriority=true}else{parent=tsNumberValue(0)};if !tsIsUndefined(exclusive){if exclusive.kind!=tsBooleanKind{tsNodeHTTP2RequestOptionType("exclusive","boolean")};if tsTruthy(exclusive){hasPriority=true}}else{exclusive=tsBooleanValue(false)};for name,value:=range map[string]tsValue{"endStream":tsBooleanValue(endStream),"waitForTrailers":tsBooleanValue(waitForTrailers),"weight":weight,"parent":parent,"exclusive":exclusive,"hasPriority":tsBooleanValue(hasPriority)}{tsNodeHTTP2Set(loop,stream,name,value);tsNodeHTTP2Set(loop,priority,name,value)};tsNodeNetEmit(loop,receiver,"streamCreated",stream);if tsNodeHTTP2SessionMemoryBytes(state)>uint64(state.maxSessionMemory)*1024*1024{if state.rejectedStreams<math.MaxUint32{state.rejectedStreams++};tsNodeHTTP2FailStream(loop,state,stream,tsHTTP2ErrorEnhanceYourCalm);if state.rejectedStreams>=state.maxSessionRejectedStreams{state.engine.localConnectionError(tsNodeHTTP2EngineSinkFor(loop,state),tsHTTP2ErrorEnhanceYourCalm,-904,[]byte("too many rejected HTTP/2 streams"))};return stream};wireEndStream:=endStream&&(!waitForTrailers||state.server);tsNodeHTTP2Set(loop,stream,"endAfterHeaders",tsBooleanValue(wireEndStream));tsNodeHTTP2SendHeaders(loop,stream,headers,wireEndStream);return stream}
func tsNodeHTTP2SessionGoaway(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeHTTP2State(receiver);code,last:=tsArg(args,0),tsArg(args,1);if tsIsUndefined(code){code=tsNumberValue(0)};lastID:=0;if state!=nil{lastID=state.lastPeerStreamID};if !tsIsUndefined(last)&&!tsIsNull(last){explicit:=int32(tsNumber(last));if explicit>0{lastID=int(explicit)}};last=tsNumberValue(float64(lastID));tsNodeHTTP2Set(loop,receiver,"goawayCode",code);tsNodeHTTP2Set(loop,receiver,"goawayLastStreamID",last);debug:=append([]byte(nil),tsFSBytes(tsArg(args,2))...);if state!=nil&&state.engine!=nil{payload:=make([]byte,8+len(debug));binary.BigEndian.PutUint32(payload,uint32(lastID)&0x7fffffff);binary.BigEndian.PutUint32(payload[4:],uint32(int32(tsNumber(code))));copy(payload[8:],debug);state.goawaySent=true;state.engine.goingAway=true;state.engine.writeFrame(tsNodeHTTP2EngineSinkFor(loop,state),7,0,0,payload)};return tsU}
func tsNodeHTTP2SessionWindow(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{value:=tsArg(args,0);if !tsIsNumeric(value)||tsNumber(value)<0||tsNumber(value)>4294967295{panic(tsThrown{tsFSException("RangeError","ERR_HTTP2_INVALID_WINDOW_SIZE","Invalid HTTP/2 window size")})};window:=uint32(tsNumber(value));state:=tsNodeHTTP2State(receiver);if state!=nil&&state.engine!=nil{connection:=state.engine;if connection.recvWindow.consumed>int64(window){panic(tsThrown{tsFSException("RangeError","ERR_HTTP2_INVALID_WINDOW_SIZE","Expected windowSize to be greater than usedWindowSize")})};old:=connection.recvWindow.size;if int64(window)>old{increment:=uint32(int64(window)-old);connection.recvWindow.grow(int64(increment));connection.sendWindowUpdate(tsNodeHTTP2EngineSinkFor(loop,state),0,increment)};if uint32(connection.localSettings.initialWindowSize)<window{connection.localSettings.initialWindowSize=window;state.localSettings=tsNodeHTTP2SettingsValue(loop,connection.localSettings);tsNodeHTTP2Set(loop,receiver,"localSettings",state.localSettings)};for _,stream:=range connection.streams{if stream.recvWindow.consumed<=int64(window){stream.recvWindow.size=int64(window)}};state.recvWindow.size=int64(window);tsNodeHTTP2RefreshPublicState(loop,state,connection.framesReceived,connection.framesSent)};tsNodeHTTP2Set(loop,receiver,"localWindowSize",tsNumberValue(float64(window)));return tsU}
func tsNodeHTTP2StreamClose(loop *tsLoop,receiver,code tsValue)tsValue{state:=tsNodeHTTP2Stream(receiver);if state==nil||state.closed{return tsU};session:=tsNodeHTTP2State(state.session);rst:=uint32(tsNumber(code));if session!=nil&&session.engine!=nil&&state.protocolState!=tsHTTP2StreamClosed{session.engine.sendRstStream(tsNodeHTTP2EngineSinkFor(loop,session),uint32(state.id),rst);if protocol:=session.engine.streams[uint32(state.id)];protocol!=nil{protocol.state=tsHTTP2StreamClosed};session.engine.closeStream(uint32(state.id));state.protocolState=tsHTTP2StreamClosed};state.closed=true;state.destroyed=true;state.ended=true;tsNodeHTTP2DrainOutbound(loop,state);state.rstCode=int(rst);tsNodeHTTP2SyncStreamLifecycle(receiver,state);tsNodeHTTP2Set(loop,receiver,"destroyed",tsBooleanValue(true));tsNodeHTTP2Set(loop,receiver,"rstCode",tsNumberValue(float64(rst)));if session!=nil{tsNodeHTTP2RefreshStreamState(loop,session,receiver);delete(session.streams,state.id)};tsNodeNetEmit(loop,receiver,"close");delete(tsNodeHTTP2Singleton.streams,uintptr(receiver.ref));return tsU}
func tsNodeHTTP2StreamSendTrailers(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeHTTP2Stream(receiver);if state==nil{return tsU};if state.sentTrailers{panic(tsThrown{tsFSException("Error","ERR_HTTP2_TRAILERS_ALREADY_SENT","Trailers have already been sent")})};if !state.trailerReady{panic(tsThrown{tsFSException("Error","ERR_HTTP2_TRAILERS_NOT_READY","Trailers are not ready to be sent")})};trailers:=tsArg(args,0);if tsNullish(trailers){trailers=tsNodeHTTP2EmptyObject()};if !tsNodeHTTP2Object(trailers)||trailers.kind==tsArrayKind{panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_TYPE","The \"headers\" argument must be of type object")})};trailers=tsNodeHTTP2HeaderObject(loop,trailers);tsNodeHTTP2ValidateHeadersFor(loop,receiver,trailers);state.trailers=trailers;state.waitForTrailers=false;state.trailerReady=false;if len((*tsObject)(trailers.ref).order)==0{tsNodeHTTP2StreamNoTrailers(loop,receiver);return tsU};state.sentTrailers=true;tsNodeHTTP2Set(loop,receiver,"sentTrailers",trailers);tsNodeHTTP2SendTrailers(loop,receiver,trailers);return tsU}
func tsNodeHTTP2StreamSetTimeout(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{duration,callback:=tsArg(args,0),tsArg(args,1);if !tsIsNumeric(duration)||tsNumber(duration)<0{panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","The timeout must be a non-negative number")})};state:=tsNodeHTTP2Stream(receiver);if state!=nil{if state.timeout!=nil{loop.clearTimeout(tsTaskValue(state.timeout))};if tsNumber(duration)>0{state.timeout=loop.setTimeout(tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{tsNodeNetEmit(loop,receiver,"timeout");return tsU})),duration)}};if callback.kind==tsFunctionKind{tsCallReceiver(loop,tsNodeHTTP2Get(loop,receiver,"once"),receiver,tsNodeHTTP2Key("timeout"),callback)};return receiver}
func tsNodeHTTP2OrderRequestPseudoHeaders(loop *tsLoop,headers tsValue)tsValue{if !tsNodeHTTP2Object(headers){return headers};source:=(*tsObject)(headers.ref);ordered:=tsNewObject();for pass:=0;pass<2;pass++{for _,name:=range source.order{pseudo:=strings.HasPrefix(name,":");if pseudo&&name==":status"{panic(tsThrown{tsFSException("TypeError","ERR_HTTP2_INVALID_PSEUDOHEADER","Invalid HTTP/2 pseudo-header for a pushed request")})};if (pass==0)!=pseudo{continue};ordered.set(name,source.values[name])}};if tsNodeHTTP2Singleton.sensitiveHeaders!=nil{symbol:=tsValue{kind:tsSymbolKind,ref:unsafe.Pointer(tsNodeHTTP2Singleton.sensitiveHeaders)};sensitive:=tsGet(loop,headers,symbol);if !tsIsUndefined(sensitive){tsSet(loop,tsObjectValue(ordered),symbol,sensitive)}};return tsObjectValue(ordered)}
func tsNodeHTTP2StreamPush(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{parent:=tsNodeHTTP2Stream(receiver);headers,options,callback:=tsArg(args,0),tsArg(args,1),tsArg(args,2);if parent==nil{return tsU};session:=tsNodeHTTP2State(parent.session);if session==nil||!session.server||session.engine==nil{panic(tsThrown{tsFSException("Error","ERR_HTTP2_INVALID_SESSION","Push streams can only be created by a server")})};headers=tsNodeHTTP2HeaderObject(loop,headers);headers=tsNodeHTTP2OrderRequestPseudoHeaders(loop,headers);tsNodeHTTP2ValidateHeadersFor(loop,receiver,headers);id:=tsNodeHTTP2AllocateStreamID(session);if id<0{panic(tsThrown{tsFSException("Error","ERR_HTTP2_OUT_OF_STREAMS","No more stream IDs are available")})};stream:=tsNodeHTTP2StreamNew(loop,parent.session,id,headers);if streamState:=tsNodeHTTP2Stream(stream);streamState!=nil{streamState.exposed=true};session.streams[id]=stream;session.engine.beginHeaderBlock();if !tsNodeHTTP2EncodeHeaderFields(loop,session,session.engine,headers,uint32(id)){_ = options;return stream};session.engine.sendPushPromise(tsNodeHTTP2EngineSinkFor(loop,session),uint32(parent.id),uint32(id));if callback.kind==tsFunctionKind{loop.post(func(){tsCall(loop,callback,tsNull,stream)})};_ = options;return stream}
func tsNodeHTTP2SessionTimeout(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{duration,callback:=tsArg(args,0),tsArg(args,1);if !tsIsNumeric(duration)||tsNumber(duration)<0{panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","The timeout must be a non-negative number")})};state:=tsNodeHTTP2State(receiver);if state!=nil{if state.timeout!=nil{loop.clearTimeout(tsTaskValue(state.timeout))};if tsNumber(duration)>0{state.timeout=loop.setTimeout(tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{tsNodeNetEmit(loop,receiver,"timeout");return tsU})),duration)}};if callback.kind==tsFunctionKind{tsCallReceiver(loop,tsNodeHTTP2Get(loop,receiver,"once"),receiver,tsNodeHTTP2Key("timeout"),callback)};return receiver}
func tsNodeHTTP2SetMethods(loop *tsLoop,prototype *tsObject){
 for name,length:=range map[string]int{"unref":0,"ref":0,"setTimeout":2,"ping":2,"goaway":3,"setLocalWindowSize":1,"settings":2,"close":1,"destroy":2,"request":2,"setNextStreamID":1,"origin":1,"altsvc":3}{method:=name;prototype.set(method,tsNamedNativeMethod(method,length,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{switch method{case "unref","ref":if session:=tsNodeHTTP2State(receiver);session!=nil&&session.socket.kind==tsObjectKind{fn:=tsNodeHTTP2Get(loop,session.socket,method);if fn.kind==tsFunctionKind{tsCallReceiver(loop,fn,session.socket)}};return receiver;case "setTimeout":return tsNodeHTTP2SessionTimeout(loop,receiver,args...);case "ping":return tsNodeHTTP2SessionPing(loop,receiver,args...);case "goaway":return tsNodeHTTP2SessionGoaway(loop,receiver,args...);case "setLocalWindowSize":return tsNodeHTTP2SessionWindow(loop,receiver,args...);case "settings":return tsNodeHTTP2SessionSettings(loop,receiver,args...);case "close":return tsNodeHTTP2SessionClose(loop,receiver,args...);case "destroy":return tsNodeHTTP2SessionDestroy(loop,receiver,args...);case "request":return tsNodeHTTP2SessionRequest(loop,receiver,args...);case "setNextStreamID":return tsNodeHTTP2SetNextStreamID(loop,receiver,tsArg(args,0));case "origin":return tsNodeHTTP2SessionOrigin(loop,receiver,args...);case "altsvc":return tsNodeHTTP2SessionAltSvc(loop,receiver,args...)};return receiver}))}
 for name,length:=range map[string]int{"sendTrailers":1,"setTimeout":2,"priority":1,"close":2,"pushStream":3,"respond":2,"respondWithFile":3,"respondWithFD":3,"additionalHeaders":1}{method:=name;prototype.set(method,tsNamedNativeMethod(method,length,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{switch method{case "sendTrailers":return tsNodeHTTP2StreamSendTrailers(loop,receiver,args...);case "setTimeout":return tsNodeHTTP2StreamSetTimeout(loop,receiver,args...);case "close":return tsNodeHTTP2StreamClose(loop,receiver,tsArg(args,0));case "priority":out:=tsNodeHTTP2EmptyObject();if tsNodeHTTP2Object(tsArg(args,0)){for k,v:=range (*tsObject)(tsArg(args,0).ref).values{outObj:=(*tsObject)(out.ref);outObj.set(k,v)}};tsNodeHTTP2Set(loop,receiver,"priority",out);return tsU;case "pushStream":return tsNodeHTTP2StreamPush(loop,receiver,args...);case "respond":headers:=tsArg(args,0);if !tsNullish(headers){headers=tsNodeHTTP2HeaderObject(loop,headers);tsNodeHTTP2ValidateHeaders(loop,headers)};tsNodeHTTP2Set(loop,receiver,"sentHeaders",headers);tsNodeHTTP2Set(loop,receiver,"headersSent",tsBooleanValue(true));state:=tsNodeHTTP2Stream(receiver);if state!=nil{state.headersSent=true};if tsTruthy(tsNodeHTTP2Get(loop,tsArg(args,1),"endStream")){tsNodeHTTP2StreamClose(loop,receiver,tsNumberValue(0))};return tsU;case "respondWithFile","respondWithFD":return tsNodeHTTP2StreamRespondFile(loop,receiver,args...);case "additionalHeaders":headers:=tsNodeHTTP2HeaderObject(loop,tsArg(args,0));tsNodeHTTP2ValidateHeaders(loop,headers);tsNodeNetEmit(loop,receiver,"headers",headers);return tsU};return tsU}))}
}
func tsNodeHTTP2StreamRespondFile(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{pathValue:=tsArg(args,0);if tsIsNumeric(pathValue){fd:=int(tsNumber(pathValue));file:=os.NewFile(uintptr(fd),"http2-response");if file==nil{panic(tsThrown{tsFSException("Error","EBADF","Bad file descriptor")})};_ = file.Close()}else if pathValue.kind!=tsStringKind{tsFSInvalid("path","a string or file descriptor")};headers:=tsArg(args,1);if !tsNullish(headers){headers=tsNodeHTTP2HeaderObject(loop,headers);tsNodeHTTP2ValidateHeaders(loop,headers)};tsNodeHTTP2Set(loop,receiver,"sentHeaders",headers);tsNodeHTTP2Set(loop,receiver,"headersSent",tsBooleanValue(true));return tsU}
func tsNodeHTTP2SessionClass(loop *tsLoop,name string,server bool)tsValue{constructor:=tsFunc(func(_ *tsLoop,_ ...tsValue)tsValue{panic(tsThrown{tsFSException("TypeError","ERR_HTTP2_INVALID_SESSION","HTTP/2 session constructors are internal")})});constructor.name=name;constructor.length=0;constructor.constructible=true;constructor.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeHTTP2SessionNew(loop,server,tsArg(args,0),tsArg(args,1))};constructor.properties=tsNewObject();constructor.properties.set("prototype",tsObjectValue(tsNodeHTTP2Singleton.sessionPrototype));return tsFunctionValue(constructor)}
func tsNodeHTTP2StreamClass(loop *tsLoop,name string,server bool)tsValue{constructor:=tsFunc(func(_ *tsLoop,_ ...tsValue)tsValue{tsPropertyFailure("Class constructor "+name+" cannot be invoked without 'new'");return tsU});constructor.name=name;constructor.length=3;constructor.constructible=true;constructor.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeHTTP2StreamNew(loop,tsArg(args,1),int(tsNumber(tsArg(args,0))),tsArg(args,2))};constructor.properties=tsNewObject();constructor.properties.set("prototype",tsObjectValue(tsNodeHTTP2Singleton.streamPrototype));return tsFunctionValue(constructor)}
func tsNodeHTTP2ObjectClass(loop *tsLoop,name string,prototype *tsObject)tsValue{constructor:=tsFunc(func(_ *tsLoop,_ ...tsValue)tsValue{tsPropertyFailure("Class constructor "+name+" cannot be invoked without 'new'");return tsU});constructor.name=name;constructor.length=1;constructor.constructible=true;constructor.construct=func(loop *tsLoop,args ...tsValue)tsValue{object:=tsNewObject();object.prototype=tsObjectValue(prototype);value:=tsObjectValue(object);tsNodeEventCreate(loop,value,tsU);return value};constructor.properties=tsNewObject();constructor.properties.set("prototype",tsObjectValue(prototype));prototype.set("constructor",tsFunctionValue(constructor));return tsFunctionValue(constructor)}
func tsNodeHTTP2Connect(loop *tsLoop,args []tsValue)tsValue{
 authority,options,listener:=tsArg(args,0),tsArg(args,1),tsArg(args,2);if options.kind==tsFunctionKind{listener=options;options=tsNodeHTTP2EmptyObject()};if tsNullish(options){options=tsNodeHTTP2EmptyObject()};if authority.kind!=tsStringKind&&authority.kind!=tsObjectKind{tsFSInvalid("authority","a string or URL")};text:=tsText(authority);if authority.kind==tsObjectKind{text=tsText(tsNodeHTTP2Get(loop,authority,"href"))};parsed,err:=url.Parse(text);if err!=nil||parsed.Scheme==""{panic(tsThrown{tsFSException("TypeError","ERR_HTTP2_INVALID_AUTHORITY","Invalid HTTP/2 authority")})};secure:=parsed.Scheme=="https:";if !secure&&parsed.Scheme!="http:"{panic(tsThrown{tsFSException("TypeError","ERR_HTTP2_INVALID_PROTOCOL","Protocol must be http: or https:")})};host:=parsed.Hostname();port:=parsed.Port();if port==""{if secure{port="443"}else{port="80"}};socketOptions:=tsNewObject();socketOptions.set("host",tsNodeHTTP2Key(host));socketOptions.set("port",tsNumberValue(float64(func()int{n,_:=strconv.Atoi(port);return n}())));if secure{socketOptions.set("ALPNProtocols",tsArrayValue(&tsArray{values:[]tsValue{tsNodeHTTP2Key("h2")}}));if servername:=tsNodeHTTP2Get(loop,options,"servername");!tsIsUndefined(servername){socketOptions.set("servername",servername)};socket:=tsNodeTLSConnect(loop,[]tsValue{tsObjectValue(socketOptions)});session:=tsNodeHTTP2SessionNew(loop,false,socket,options);state:=tsNodeHTTP2State(session);state.encrypted=true;tsNodeHTTP2Set(loop,session,"connecting",tsBooleanValue(true));if listener.kind==tsFunctionKind{tsCallReceiver(loop,tsNodeHTTP2Get(loop,session,"once"),session,tsNodeHTTP2Key("connect"),listener)};tsCallReceiver(loop,tsNodeHTTP2Get(loop,socket,"once"),socket,tsNodeHTTP2Key("secureConnect"),tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{state.connected=true;tsNodeHTTP2Set(loop,session,"connecting",tsBooleanValue(false));tsNodeHTTP2Set(loop,session,"connected",tsBooleanValue(true));tsNodeNetEmit(loop,session,"connect",socket,tsNodeHTTP2SettingsDefaults());tsNodeNetEmit(loop,session,"remoteSettings",state.remoteSettings);return tsU})));return session};socket:=tsNodeNetSocketCreate(loop,options);session:=tsNodeHTTP2SessionNew(loop,false,socket,options);state:=tsNodeHTTP2State(session);state.connected=false;tsNodeHTTP2Set(loop,session,"connecting",tsBooleanValue(true));if listener.kind==tsFunctionKind{tsCallReceiver(loop,tsNodeHTTP2Get(loop,session,"once"),session,tsNodeHTTP2Key("connect"),listener)};tsCallReceiver(loop,tsNodeHTTP2Get(loop,socket,"once"),socket,tsNodeHTTP2Key("connect"),tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{state.connected=true;tsNodeHTTP2Set(loop,session,"connecting",tsBooleanValue(false));tsNodeHTTP2Set(loop,session,"connected",tsBooleanValue(true));tsNodeNetEmit(loop,session,"connect",socket,tsNodeHTTP2SettingsDefaults());tsNodeNetEmit(loop,session,"remoteSettings",state.remoteSettings);return tsU})));_ = tsNodeNetModule(loop);raw:=tsNodeNetSocketFor(socket);if raw!=nil{tsNodeNetSocketConnect(loop,socket,socketOptions)};return session
}

func tsNodeHTTP2RequestCreate(loop *tsLoop,stream,headers tsValue)tsValue{optionsObject:=tsNewObject();optionsObject.set("read",tsNamedNativeMethod("_read",1,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsU}));options:=tsObjectValue(optionsObject);request:=tsNodeReadableCreate(loop,tsU,options);object:=(*tsObject)(request.ref);object.prototype=tsObjectValue(tsNodeHTTP2Singleton.serverRequestPrototype);tsNodeEventCreate(loop,request,tsU);if tsNullish(headers){headers=tsNodeHTTP2EmptyObject()};raw:=tsNodeHTTP2HeaderArray(loop,headers);for name,value:=range map[string]tsValue{"stream":stream,"session":tsNodeHTTP2Get(loop,stream,"session"),"headers":headers,"rawHeaders":raw,"trailers":tsNodeHTTP2EmptyObject(),"rawTrailers":tsArrayValue(&tsArray{}),"httpVersion":tsNodeHTTP2Key("2.0"),"httpVersionMajor":tsNumberValue(2),"httpVersionMinor":tsNumberValue(0),"complete":tsBooleanValue(false),"aborted":tsBooleanValue(false),"url":tsNodeHTTP2Get(loop,headers,":path"),"method":tsNodeHTTP2Get(loop,headers,":method"),"socket":tsNodeHTTP2Get(loop,tsNodeHTTP2Get(loop,stream,"session"),"socket"),"rawSocket":tsNodeHTTP2Get(loop,tsNodeHTTP2Get(loop,stream,"session"),"socket"),"headersDistinct":headers,"trailersDistinct":tsNodeHTTP2EmptyObject()}{tsNodeHTTP2Set(loop,request,name,value)};return request}
func tsNodeHTTP2ResponseCreate(loop *tsLoop,stream tsValue)tsValue{object:=tsNewObject();object.prototype=tsObjectValue(tsNodeHTTP2Singleton.serverResponsePrototype);response:=tsObjectValue(object);tsNodeEventCreate(loop,response,tsU);for name,value:=range map[string]tsValue{"stream":stream,"session":tsNodeHTTP2Get(loop,stream,"session"),"headers":tsNodeHTTP2EmptyObject(),"trailers":tsNodeHTTP2EmptyObject(),"statusCode":tsNumberValue(200),"statusMessage":tsNodeHTTP2Key(""),"headersSent":tsBooleanValue(false),"finished":tsBooleanValue(false),"writableEnded":tsBooleanValue(false),"writableFinished":tsBooleanValue(false),"sendDate":tsBooleanValue(true),"socket":tsNodeHTTP2Get(loop,tsNodeHTTP2Get(loop,stream,"session"),"socket"),"req":tsNodeHTTP2RequestCreate(loop,stream,tsNodeHTTP2Get(loop,stream,"headers"))}{tsNodeHTTP2Set(loop,response,name,value)};return response}
func tsNodeHTTP2ResponseHeaders(loop *tsLoop,receiver tsValue)tsValue{value:=tsNodeHTTP2Get(loop,receiver,"headers");if tsNodeHTTP2Object(value){return tsObjectValue(tsObjectSpreadCopy((*tsObject)(value.ref)))};return tsNodeHTTP2EmptyObject()}
func tsObjectSpreadCopy(source *tsObject)*tsObject{out:=tsNewObject();if source!=nil{for _,name:=range source.order{out.set(name,source.values[name])}};return out}
func tsNodeHTTP2ResponseHeaderMethod(loop *tsLoop,receiver tsValue,name string,args []tsValue)tsValue{headers:=tsNodeHTTP2ResponseHeaders(loop,receiver);object:=(*tsObject)(headers.ref);key:=strings.ToLower(tsText(tsArg(args,0)));switch name{case "setHeader":if tsTruthy(tsNodeHTTP2Get(loop,receiver,"headersSent")){panic(tsThrown{tsFSException("Error","ERR_HTTP2_HEADERS_SENT","Cannot set headers after they are sent")})};object.set(key,tsArg(args,1));tsNodeHTTP2Set(loop,receiver,"headers",headers);return tsArg(args,1);case "appendHeader":value:=tsArg(args,1);old:=object.get(key);if tsIsUndefined(old){object.set(key,value)}else if old.kind==tsArrayKind{(*tsArray)(old.ref).values=append((*tsArray)(old.ref).values,value)}else{object.set(key,tsArrayValue(&tsArray{values:[]tsValue{old,value}}))};tsNodeHTTP2Set(loop,receiver,"headers",headers);return tsU;case "getHeader":return object.get(key);case "hasHeader":return tsBooleanValue(!tsIsUndefined(object.get(key)));case "removeHeader":delete(object.values,key);return tsU;case "getHeaders":return headers;case "getHeaderNames":out:=&tsArray{};for _,entry:=range object.order{out.values=append(out.values,tsNodeHTTP2Key(entry))};return tsArrayValue(out);case "addTrailers":trailers:=tsNodeHTTP2HeaderObject(loop,tsArg(args,0));tsNodeHTTP2ValidateHeaders(loop,trailers);tsNodeHTTP2Set(loop,receiver,"trailers",trailers);return tsU};return tsU}
func tsNodeHTTP2ResponseWriteHead(loop *tsLoop,receiver tsValue,args []tsValue)tsValue{status:=tsArg(args,0);if !tsIsNumeric(status)||tsNumber(status)<100||tsNumber(status)>999{panic(tsThrown{tsFSException("RangeError","ERR_HTTP2_STATUS_INVALID","Invalid status code")})};tsNodeHTTP2Set(loop,receiver,"statusCode",status);headers:=tsArg(args,2);if tsArg(args,1).kind==tsObjectKind{headers=tsArg(args,1)};if !tsNullish(headers){headers=tsNodeHTTP2HeaderObject(loop,headers);tsNodeHTTP2ValidateHeaders(loop,headers);current:=tsNodeHTTP2ResponseHeaders(loop,receiver);source:=(*tsObject)(headers.ref);for _,name:=range source.order{(*tsObject)(current.ref).set(name,source.values[name])};tsNodeHTTP2Set(loop,receiver,"headers",current)};return tsNodeHTTP2ResponseEndHeaders(loop,receiver)}
func tsNodeHTTP2ResponseEndHeaders(loop *tsLoop,receiver tsValue)tsValue{headers:=tsNodeHTTP2ResponseHeaders(loop,receiver);headersObj:=(*tsObject)(headers.ref);headersObj.set(":status",tsNodeHTTP2Get(loop,receiver,"statusCode"));stream:=tsNodeHTTP2Get(loop,receiver,"stream");tsNodeHTTP2StreamRespond(loop,stream,headers,tsNodeHTTP2EmptyObject());tsNodeHTTP2Set(loop,receiver,"headersSent",tsBooleanValue(true));return receiver}
func tsNodeHTTP2StreamRespond(loop *tsLoop,stream,headers,options tsValue)tsValue{if tsNullish(options){options=tsNodeHTTP2EmptyObject()};if !tsNullish(headers){headers=tsNodeHTTP2HeaderObject(loop,headers);tsNodeHTTP2ValidateHeadersFor(loop,stream,headers)}else{headers=tsNodeHTTP2EmptyObject()};if tsIsUndefined(tsNodeHTTP2Get(loop,headers,":status")){tsNodeHTTP2Set(loop,headers,":status",tsNumberValue(200))};endStream:=tsTruthy(tsNodeHTTP2Get(loop,options,"endStream"));tsNodeHTTP2SendHeaders(loop,stream,headers,endStream);if state:=tsNodeHTTP2Stream(stream);state!=nil{state.waitForTrailers=tsTruthy(tsNodeHTTP2Get(loop,options,"waitForTrailers"));tsNodeHTTP2Set(loop,stream,"waitForTrailers",tsBooleanValue(state.waitForTrailers))};return tsU}
func tsNodeHTTP2ResponseWrite(loop *tsLoop,receiver tsValue,args []tsValue)tsValue{if !tsTruthy(tsNodeHTTP2Get(loop,receiver,"headersSent")){tsNodeHTTP2ResponseEndHeaders(loop,receiver)};stream:=tsNodeHTTP2Get(loop,receiver,"stream");write:=tsNodeHTTP2Get(loop,stream,"write");if write.kind==tsFunctionKind{return tsCallReceiver(loop,write,stream,tsArg(args,0),tsArg(args,1),tsArg(args,2))};return tsBooleanValue(true)}
func tsNodeHTTP2ResponseEnd(loop *tsLoop,receiver tsValue,args []tsValue)tsValue{if !tsTruthy(tsNodeHTTP2Get(loop,receiver,"headersSent")){tsNodeHTTP2ResponseEndHeaders(loop,receiver)};stream:=tsNodeHTTP2Get(loop,receiver,"stream");end:=tsNodeHTTP2Get(loop,stream,"end");if end.kind==tsFunctionKind{tsCallReceiver(loop,end,stream,tsArg(args,0),tsArg(args,1),tsArg(args,2))};tsNodeHTTP2Set(loop,receiver,"writableEnded",tsBooleanValue(true));tsNodeHTTP2Set(loop,receiver,"finished",tsBooleanValue(true));tsNodeHTTP2Set(loop,receiver,"writableFinished",tsBooleanValue(true));return receiver}
func tsNodeHTTP2InstallResponse(loop *tsLoop){p:=tsNodeHTTP2Singleton.serverResponsePrototype;for name,length:=range map[string]int{"setHeader":2,"appendHeader":2,"getHeader":1,"getHeaderNames":0,"getHeaders":0,"hasHeader":1,"removeHeader":1,"addTrailers":1}{method:=name;p.set(method,tsNamedNativeMethod(method,length,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeHTTP2ResponseHeaderMethod(loop,receiver,method,args)}))};p.set("writeHead",tsNamedNativeMethod("writeHead",3,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeHTTP2ResponseWriteHead(loop,receiver,args)}));p.set("flushHeaders",tsNamedNativeMethod("flushHeaders",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{return tsNodeHTTP2ResponseEndHeaders(loop,receiver)}));p.set("write",tsNamedNativeMethod("write",3,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeHTTP2ResponseWrite(loop,receiver,args)}));p.set("end",tsNamedNativeMethod("end",3,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeHTTP2ResponseEnd(loop,receiver,args)}));p.set("destroy",tsNamedNativeMethod("destroy",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{stream:=tsNodeHTTP2Get(loop,receiver,"stream");if destroy:=tsNodeHTTP2Get(loop,stream,"destroy");destroy.kind==tsFunctionKind{tsCallReceiver(loop,destroy,stream,tsArg(args,0))};return receiver}));p.set("setTimeout",tsNamedNativeMethod("setTimeout",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeHTTP2StreamSetTimeout(loop,tsNodeHTTP2Get(loop,receiver,"stream"),args...)}));for name,status:=range map[string]int{"writeContinue":100,"writeInformation":103,"writeEarlyHints":103}{method,code:=name,status;p.set(method,tsNamedNativeMethod(method,2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{headers:=tsNodeHTTP2HeaderObject(loop,tsArg(args,1));headersObj:=(*tsObject)(headers.ref);headersObj.set(":status",tsNumberValue(float64(code)));tsNodeHTTP2StreamRespond(loop,tsNodeHTTP2Get(loop,receiver,"stream"),headers,tsNodeHTTP2EmptyObject());return tsU}))}}
func tsNodeHTTP2InstallStream(loop *tsLoop){p:=tsNodeHTTP2Singleton.streamPrototype;tsNodeHTTP2SetMethods(loop,p);p.set("sendTrailers",tsNamedNativeMethod("sendTrailers",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeHTTP2StreamSendTrailers(loop,receiver,args...)}));p.set("setTimeout",tsNamedNativeMethod("setTimeout",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeHTTP2StreamSetTimeout(loop,receiver,args...)}));p.set("close",tsNamedNativeMethod("close",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeHTTP2StreamClose(loop,receiver,tsArg(args,0))}));p.set("pushStream",tsNamedNativeMethod("pushStream",3,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeHTTP2StreamPush(loop,receiver,args...)}));p.set("respond",tsNamedNativeMethod("respond",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeHTTP2StreamRespond(loop,receiver,tsArg(args,0),tsArg(args,1))}));for _,name:=range []string{"respondWithFile","respondWithFD"}{method:=name;p.set(method,tsNamedNativeMethod(method,3,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeHTTP2StreamRespondFile(loop,receiver,args...)}))};p.set("additionalHeaders",tsNamedNativeMethod("additionalHeaders",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{headers:=tsNodeHTTP2HeaderObject(loop,tsArg(args,0));tsNodeHTTP2ValidateHeaders(loop,headers);tsNodeNetEmit(loop,receiver,"headers",headers);return tsU}))}
func tsNodeHTTP2AttachServerSessionErrorListener(loop *tsLoop,server,session tsValue){listener:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{tsNodeNetEmitSafe(loop,server,"sessionError",tsArg(args,0),session);return tsU}));tsCallReceiver(loop,tsNodeHTTP2Get(loop,session,"on"),session,tsNodeHTTP2Key("error"),listener)}
func tsNodeHTTP2AttachServer(loop *tsLoop,server tsValue,secure bool,options,listener tsValue)tsValue{state:=&tsNodeHTTP2ServerState{object:server,secure:secure,sessions:map[uintptr]tsValue{},options:options,listener:listener};tsNodeHTTP2Singleton.servers[uintptr(server.ref)]=state;event:="connection";if secure{event="secureConnection"};handler:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{socket:=tsArg(args,0);session:=tsNodeHTTP2SessionNew(loop,true,socket,options);sessionState:=tsNodeHTTP2State(session);sessionState.serverObject=server;sessionState.connected=true;sessionState.encrypted=secure;tsNodeHTTP2Set(loop,session,"connected",tsBooleanValue(true));tsNodeHTTP2Set(loop,session,"connecting",tsBooleanValue(false));tsNodeHTTP2Set(loop,session,"alpnProtocol",tsNodeHTTP2Key("h2"));state.sessions[uintptr(session.ref)]=session;tsNodeHTTP2AttachServerSessionErrorListener(loop,server,session);tsNodeNetEmit(loop,server,"session",session);tsNodeNetEmit(loop,session,"localSettings",sessionState.localSettings);return tsU}));tsCallReceiver(loop,tsNodeHTTP2Get(loop,server,"on"),server,tsNodeHTTP2Key(event),handler);if listener.kind==tsFunctionKind{tsCallReceiver(loop,tsNodeHTTP2Get(loop,server,"on"),server,tsNodeHTTP2Key("request"),listener)};return server}
func tsNodeHTTP2ServerSettings(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{settings:=tsArg(args,0);tsNodeHTTP2ValidateSettings(loop,settings);state:=tsNodeHTTP2Singleton.servers[uintptr(receiver.ref)];if state!=nil{state.options=tsObjectValue(tsObjectSpreadCopy((*tsObject)(state.options.ref)));tsNodeHTTP2Set(loop,state.options,"settings",settings);for _,session:=range state.sessions{tsNodeHTTP2SessionSettings(loop,session,settings)}};return receiver}
func tsNodeHTTP2ServerClose(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{close:=tsNodeHTTP2Get(loop,receiver,"close");if close.kind==tsFunctionKind{return tsCallReceiver(loop,close,receiver,tsArg(args,0))};return receiver}
func tsNodeHTTP2InstallServer(loop *tsLoop){netModule:=tsNodeNetModule(loop);base:=tsNodeNetSingleton.serverProto;serverProto:=tsNewObject();serverProto.prototype=tsObjectValue(base);serverProto.set("updateSettings",tsNamedNativeMethod("updateSettings",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeHTTP2ServerSettings(loop,receiver,args...)}));serverProto.set("closeIdleConnections",tsNamedNativeMethod("closeIdleConnections",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{return receiver}));serverProto.set("setTimeout",tsNodeNetGet(loop,base,"setTimeout"));serverProto.set("close",tsNodeNetGet(loop,base,"close"));tsNodeHTTP2Singleton.serverPrototype=serverProto;tlsModule:=tsNodeTLSModule(loop);tlsCtor:=tsNodeTLSGet(loop,tlsModule,"Server");tlsProto:=tsNodeTLSGet(loop,tlsCtor,"prototype");secureProto:=tsNewObject();secureProto.prototype=tlsProto;secureProto.set("updateSettings",tsNamedNativeMethod("updateSettings",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeHTTP2ServerSettings(loop,receiver,args...)}));secureProto.set("setTimeout",tsNodeTLSGet(loop,tlsProto,"setTimeout"));secureProto.set("close",tsNodeNetGet(loop,base,"close"));secureProto.set("closeIdleConnections",tsNamedNativeMethod("closeIdleConnections",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{return receiver}));tsNodeHTTP2Singleton.secureServerPrototype=secureProto;_ = netModule}
func tsNodeHTTP2CreateServer(loop *tsLoop,args []tsValue,secure bool)tsValue{options,listener:=tsArg(args,0),tsArg(args,1);if options.kind==tsFunctionKind{listener=options;options=tsNodeHTTP2EmptyObject()};if tsNullish(options){options=tsNodeHTTP2EmptyObject()};tsNodeHTTP2ValidateSettings(loop,func()tsValue{value:=tsNodeHTTP2Get(loop,options,"settings");if tsIsUndefined(value){return tsNodeHTTP2EmptyObject()};return value}());var server tsValue;if secure{server=tsNodeTLSCreateServer(loop,[]tsValue{options})}else{server=tsNodeNetServerCreate(loop,options)};(*tsObject)(server.ref).prototype=tsObjectValue(func()*tsObject{if secure{return tsNodeHTTP2Singleton.secureServerPrototype};return tsNodeHTTP2Singleton.serverPrototype}());return tsNodeHTTP2AttachServer(loop,server,secure,options,listener)}
func tsNodeHTTP2PerformServerHandshake(loop *tsLoop,args []tsValue)tsValue{socket,options:=tsArg(args,0),tsArg(args,1);if tsNullish(options){options=tsNodeHTTP2EmptyObject()};session:=tsNodeHTTP2SessionNew(loop,true,socket,options);state:=tsNodeHTTP2State(session);state.connected=true;tsNodeHTTP2Set(loop,session,"connected",tsBooleanValue(true));tsNodeNetEmit(loop,session,"localSettings",state.localSettings);return session}
func tsNodeHTTP2Module(loop *tsLoop)tsValue{tsNodeHTTP2Singleton.Do(func(){tsNodeHTTP2Singleton.sessions=map[uintptr]*tsNodeHTTP2SessionState{};tsNodeHTTP2Singleton.streams=map[uintptr]*tsNodeHTTP2StreamState{};tsNodeHTTP2Singleton.servers=map[uintptr]*tsNodeHTTP2ServerState{};tsNodeHTTP2Singleton.sensitiveHeaders=(*tsSymbol)(tsNewSymbol(tsStringUTF8("nodejs.http2.sensitiveHeaders")).ref);events:=tsNodeEventsModule(loop);eventCtor:=tsNodeHTTP2Get(loop,events,"EventEmitter");eventProto:=tsNodeHTTP2Get(loop,eventCtor,"prototype");tsNodeHTTP2Singleton.sessionPrototype=tsNewObject();tsNodeHTTP2Singleton.sessionPrototype.prototype=eventProto;tsNodeHTTP2Singleton.streamPrototype=tsNewObject();tsNodeDuplexInit(loop);tsNodeHTTP2Singleton.streamPrototype.prototype=tsObjectValue(tsNodeDuplexSingleton.prototype);tsNodeHTTP2Singleton.serverRequestPrototype=tsNewObject();tsNodeReadableInit(loop);tsNodeHTTP2Singleton.serverRequestPrototype.prototype=tsNodeReadableSingleton.prototype;tsNodeHTTP2Singleton.serverResponsePrototype=tsNewObject();tsNodeHTTP2Singleton.serverResponsePrototype.prototype=eventProto;tsNodeHTTP2SetMethods(loop,tsNodeHTTP2Singleton.sessionPrototype);tsNodeHTTP2InstallStream(loop);tsNodeHTTP2InstallResponse(loop);tsNodeHTTP2InstallServer(loop)
 module:=tsNewObject();constants:=tsNodeHTTP2BuildConstants(loop);tsNodeHTTP2Singleton.constants=constants;module.set("constants",constants);module.set("sensitiveHeaders",tsValue{kind:tsSymbolKind,ref:unsafe.Pointer(tsNodeHTTP2Singleton.sensitiveHeaders)});module.set("getDefaultSettings",tsNamedNativeMethod("getDefaultSettings",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNodeHTTP2SettingsDefaults()}));module.set("getPackedSettings",tsNamedNativeMethod("getPackedSettings",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeHTTP2PackSettings(loop,tsArg(args,0))}));module.set("getUnpackedSettings",tsNamedNativeMethod("getUnpackedSettings",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeHTTP2UnpackSettings(loop,tsArg(args,0),tsArg(args,1))}));module.set("connect",tsNamedNativeMethod("connect",3,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeHTTP2Connect(loop,args)}));module.set("performServerHandshake",tsNamedNativeMethod("performServerHandshake",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeHTTP2PerformServerHandshake(loop,args)}));module.set("createServer",tsNamedNativeMethod("createServer",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeHTTP2CreateServer(loop,args,false)}));module.set("createSecureServer",tsNamedNativeMethod("createSecureServer",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeHTTP2CreateServer(loop,args,true)}));module.set("Http2ServerRequest",tsNodeHTTP2ObjectClass(loop,"Http2ServerRequest",tsNodeHTTP2Singleton.serverRequestPrototype));module.set("Http2ServerResponse",tsNodeHTTP2ObjectClass(loop,"Http2ServerResponse",tsNodeHTTP2Singleton.serverResponsePrototype));module.set("ClientHttp2Session",tsNodeHTTP2SessionClass(loop,"ClientHttp2Session",false));module.set("Http2Session",tsNodeHTTP2SessionClass(loop,"Http2Session",false));module.set("ServerHttp2Session",tsNodeHTTP2SessionClass(loop,"ServerHttp2Session",true));module.set("Http2Stream",tsNodeHTTP2StreamClass(loop,"Http2Stream",false));module.set("ClientHttp2Stream",tsNodeHTTP2StreamClass(loop,"ClientHttp2Stream",false));module.set("ServerHttp2Stream",tsNodeHTTP2StreamClass(loop,"ServerHttp2Stream",true));module.set("default",tsObjectValue(module));tsNodeHTTP2Singleton.module=tsObjectValue(module)
 });return tsNodeHTTP2Singleton.module}
`

const nodeHTTP2Runtime = nodeHTTP2BaseRuntime + nodeHTTP2EngineRuntime + nodeHTTP2DataRuntime + nodeHTTP2StreamRuntime + nodeHTTP2TransportRuntime
