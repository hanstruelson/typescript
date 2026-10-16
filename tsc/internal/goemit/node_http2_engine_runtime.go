package goemit

// nodeHTTP2EngineRuntime ports the protocol-side state and control-frame
// operations from Bun's src/runtime/api/bun/h2/connection.rs.
const nodeHTTP2EngineRuntime = `
const (
 tsHTTP2PseudoMethod uint8 = 1 << iota
 tsHTTP2PseudoScheme
 tsHTTP2PseudoAuthority
 tsHTTP2PseudoPath
 tsHTTP2PseudoStatus
 tsHTTP2PseudoProtocol
 tsHTTP2PseudoUnknown
)
const tsHTTP2MaxOutboundAckQueue uint32=1000
const tsHTTP2DefaultStreamResetBurst uint32=1000
const tsHTTP2DefaultStreamResetRate uint32=33
const (
 tsHTTP2WriteDropped = -1
 tsHTTP2WriteQueued = 0
 tsHTTP2WriteSent = 1
)
type tsHTTP2PendingLocalSettings struct { settings tsHTTP2SettingsState }
type tsHTTP2ProtocolStream struct {
 state uint8
 sendWindow tsHTTP2SendWindow
 recvWindow tsHTTP2RecvWindow
 recvFinalHeaders bool
 contentLength *uint64
 recvBodyBytes uint64
 openedAfterGoaway bool
}
func tsHTTP2NewProtocolStream(initialSend,initialRecv uint32)*tsHTTP2ProtocolStream{return &tsHTTP2ProtocolStream{state:tsHTTP2StreamIdle,sendWindow:tsHTTP2NewSendWindow(initialSend),recvWindow:tsHTTP2NewRecvWindow(initialRecv)}}
type tsHTTP2DataInFlight struct { streamID,payloadRemaining,dataRemaining uint32; endStream,discard bool }
type tsHTTP2BlockDisposition uint8
const (
 tsHTTP2BlockDeliver tsHTTP2BlockDisposition = iota
 tsHTTP2BlockStreamClosed
 tsHTTP2BlockRefused
)
type tsHTTP2HeaderBlockMeta struct { target,pushParent uint32; hasPushParent bool; flags uint8; endStream,isRequest bool; disposition tsHTTP2BlockDisposition }
type tsHTTP2HeaderBlockInFlight struct { continuationStream uint32; meta tsHTTP2HeaderBlockMeta }
type tsHTTP2Feed struct { consumed int; fatal bool }
type tsHTTP2RateLimit struct { burst,rate,tokens uint32; lastRefill time.Time }
func tsHTTP2NewRateLimit(burst,rate uint32)tsHTTP2RateLimit{return tsHTTP2RateLimit{burst:burst,rate:rate,tokens:burst,lastRefill:time.Now()}}
func (limit *tsHTTP2RateLimit)drain()bool{elapsed:=uint64(time.Since(limit.lastRefill)/time.Second);if elapsed>0{seconds:=elapsed;if seconds>math.MaxUint32{seconds=math.MaxUint32};gain:=seconds*uint64(limit.rate);if gain>math.MaxUint32{gain=math.MaxUint32};available:=uint64(limit.tokens)+gain;if available>uint64(limit.burst){available=uint64(limit.burst)};limit.tokens=uint32(available);limit.lastRefill=limit.lastRefill.Add(time.Duration(seconds)*time.Second)};if limit.tokens==0{return false};limit.tokens--;return true}
type tsHTTP2Sink interface {
 write([]byte)int
 onError(int32,uint32,[]byte)
 onLocalSettings(tsHTTP2SettingsState)
 onRemoteSettings(tsHTTP2SettingsState)
 onPing([]byte,bool)
 onGoAway(uint32,uint32,[]byte)
 onWindowUpdate(uint32,uint32)
 shouldStop()bool
 onStreamOpen(uint32)
 canOpenStream()bool
 goawaySent()bool
 onRemoteCustomSetting(uint16,uint32)
 onHeader(uint32,[]byte,[]byte,bool)
 onHeadersComplete(uint32,bool,uint8)
 onData(uint32,[]byte)
 onStreamEnd(uint32,uint8)
 onStreamReset(uint32,uint32)
 onStreamRejected(uint32)
 onPushPromise(uint32,uint32)
 onAltSvc(uint32,[]byte,[]byte)
 onOrigin([]byte)
 onTooManyInvalidFrames()
 onFrameCounters(uint64,uint64)
 isLocalStream(uint32)bool
 highestStartedStreamID()uint32
 isStreamReading(uint32)bool
}
type tsHTTP2Connection struct {
 isServer bool
 framesReceived,framesSent uint64
 localSettings,remoteSettings tsHTTP2SettingsState
 localSettingsAcked bool
 maxHeaderListPairs,maxInvalidFrames uint32
 ackedLocalInitialWindow,enforcedMaxHeaderListSize uint32
 pendingLocalSettingsAcks []tsHTTP2PendingLocalSettings
 invalidFrameCount,maxSettings uint32
 sendWindow tsHTTP2SendWindow
 recvWindow tsHTTP2RecvWindow
 hpack *tsHTTP2HPACKCoder
 streams map[uint32]*tsHTTP2ProtocolStream
 headerBlockInFlight *tsHTTP2HeaderBlockInFlight
 headerBlock []byte
 dataInFlight *tsHTTP2DataInFlight
 streamedDataConsumed int
 terminated bool
 outboundAckPending uint32
 streamResetLimit,sentResetLimit tsHTTP2RateLimit
 resetFlood bool
 encBuf []byte
 replenishBuf [][2]uint32
 evictBuf []uint32
 prefaceReceived uint32
 lastStreamID uint32
 goingAway bool
}
func tsHTTP2NewConnection(isServer bool,local tsHTTP2SettingsState)*tsHTTP2Connection{return &tsHTTP2Connection{isServer:isServer,localSettings:local,remoteSettings:tsHTTP2DefaultSettingsState(),maxHeaderListPairs:128,maxInvalidFrames:1000,ackedLocalInitialWindow:65535,enforcedMaxHeaderListSize:local.maxHeaderListSize,maxSettings:32,sendWindow:tsHTTP2NewSendWindow(65535),recvWindow:tsHTTP2NewRecvWindow(65535),hpack:tsHTTP2NewHPACKCoder(tsHTTP2DefaultSettingsState().headerTableSize,tsHTTP2NewLiteralHPACKEncoder()),streams:map[uint32]*tsHTTP2ProtocolStream{},streamResetLimit:tsHTTP2NewRateLimit(tsHTTP2DefaultStreamResetBurst,tsHTTP2DefaultStreamResetRate),sentResetLimit:tsHTTP2NewRateLimit(tsHTTP2DefaultStreamResetBurst,tsHTTP2DefaultStreamResetRate)}}
func (connection *tsHTTP2Connection)writeFrame(sink tsHTTP2Sink,frameType,flags uint8,streamID uint32,payload []byte){connection.framesSent++;sink.onFrameCounters(connection.framesReceived,connection.framesSent);header:=tsHTTP2WriteFrameHeader(tsHTTP2FrameHeader{length:uint32(len(payload)),frameType:frameType,flags:flags,streamID:streamID});sink.write(header);if len(payload)>0{sink.write(payload)}}
func (connection *tsHTTP2Connection)sendSettingsAck(sink tsHTTP2Sink){connection.writeFrame(sink,4,1,0,nil)}
func (connection *tsHTTP2Connection)sendPingAck(sink tsHTTP2Sink,payload []byte){connection.writeFrame(sink,6,1,0,payload)}
func (connection *tsHTTP2Connection)sendWindowUpdate(sink tsHTTP2Sink,streamID,increment uint32){payload:=make([]byte,4);binary.BigEndian.PutUint32(payload,increment);connection.writeFrame(sink,8,0,streamID,payload)}
func (connection *tsHTTP2Connection)localConnectionError(sink tsHTTP2Sink,code uint32,libCode int32,debug []byte){connection.goingAway=true;connection.terminated=true;payload:=make([]byte,8+len(debug));binary.BigEndian.PutUint32(payload,connection.lastStreamID&0x7fffffff);binary.BigEndian.PutUint32(payload[4:],code);copy(payload[8:],debug);connection.writeFrame(sink,7,0,0,payload);sink.onError(libCode,connection.lastStreamID,debug)}
func (connection *tsHTTP2Connection)sendGoAway(sink tsHTTP2Sink,code uint32,debug []byte){connection.localConnectionError(sink,code,-505,debug)}
func (connection *tsHTTP2Connection)noteReset(byPeer bool){if !connection.isServer||connection.resetFlood{return};limit:=&connection.sentResetLimit;if byPeer{limit=&connection.streamResetLimit};if !limit.drain(){connection.resetFlood=true}}
func (connection *tsHTTP2Connection)sendRstStream(sink tsHTTP2Sink,streamID,code uint32){payload:=make([]byte,4);binary.BigEndian.PutUint32(payload,code);connection.writeFrame(sink,3,0,streamID,payload);if stream:=connection.streams[streamID];stream!=nil&&stream.recvFinalHeaders&&stream.state!=tsHTTP2StreamClosed{connection.noteReset(false)}}
func (connection *tsHTTP2Connection)checkResetFlood(sink tsHTTP2Sink)bool{if !connection.resetFlood||connection.terminated{return false};connection.localConnectionError(sink,tsHTTP2ErrorEnhanceYourCalm,-904,[]byte("too many stream resets"));return true}
func (connection *tsHTTP2Connection)noteOutboundAck(sink tsHTTP2Sink)bool{if connection.outboundAckPending<math.MaxUint32{connection.outboundAckPending++};if connection.outboundAckPending<tsHTTP2MaxOutboundAckQueue{return false};connection.localConnectionError(sink,tsHTTP2ErrorEnhanceYourCalm,-904,[]byte("too many outbound control frames queued"));return true}
func (connection *tsHTTP2Connection)noteOutboundDrained(){connection.outboundAckPending=0}
func (connection *tsHTTP2Connection)setStreamResetLimit(burst,rate uint32){if connection.streamResetLimit.burst!=burst||connection.streamResetLimit.rate!=rate{connection.streamResetLimit=tsHTTP2NewRateLimit(burst,rate)}}
func (connection *tsHTTP2Connection)receive(sink tsHTTP2Sink,data []byte)tsHTTP2Feed{offset:=0;if connection.terminated{return tsHTTP2Feed{consumed:len(data),fatal:true}};preface:=[]byte("PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n");if connection.isServer&&connection.prefaceReceived<uint32(len(preface)){need:=len(preface)-int(connection.prefaceReceived);available:=min(need,len(data));for i:=0;i<available;i++{if data[i]!=preface[int(connection.prefaceReceived)+i]{connection.localConnectionError(sink,tsHTTP2ErrorProtocol,-903,[]byte("invalid connection preface"));return tsHTTP2Feed{consumed:available,fatal:true}}};connection.prefaceReceived+=uint32(available);offset+=available;if connection.prefaceReceived<uint32(len(preface)){return tsHTTP2Feed{consumed:offset}}};if connection.dataInFlight!=nil{inflight:=*connection.dataInFlight;connection.dataInFlight=nil;available:=min(len(data)-offset,int(inflight.payloadRemaining));dataNow:=min(available,int(inflight.dataRemaining));if dataNow>0&&!inflight.discard{sink.onData(inflight.streamID,data[offset:offset+dataNow])};inflight.dataRemaining-=uint32(dataNow);inflight.payloadRemaining-=uint32(available);offset+=available;if inflight.payloadRemaining==0{connection.finishStreamedData(sink,&inflight)}else{connection.dataInFlight=&inflight;connection.replenishWindows(sink);return tsHTTP2Feed{consumed:offset}}};for{remaining:=data[offset:];if len(remaining)<9||sink.shouldStop(){break};header,ok:=tsHTTP2ParseFrameHeader(remaining);if !ok{break};if header.length>connection.localSettings.maxFrameSize{connection.sendGoAway(sink,tsHTTP2ErrorFrameSize,[]byte("frame exceeds SETTINGS_MAX_FRAME_SIZE"));return tsHTTP2Feed{consumed:offset,fatal:true}};total:=9+int(header.length);if len(remaining)<total{if header.frameType==0&&connection.headerBlockInFlight==nil&&header.streamID!=0&&header.length>0{padded:=tsHTTP2HasFrameFlag(header.flags,8);availablePayload:=len(remaining)-9;if !padded||availablePayload>=1{if connection.beginStreamedData(sink,header,remaining,padded){return tsHTTP2Feed{consumed:offset,fatal:true}};offset+=connection.streamedDataConsumed}};break};payload:=remaining[9:total];if connection.dispatch(sink,header,payload){return tsHTTP2Feed{consumed:offset+total,fatal:true}};offset+=total;if connection.checkResetFlood(sink){return tsHTTP2Feed{consumed:offset,fatal:true}}};if connection.checkResetFlood(sink){return tsHTTP2Feed{consumed:offset,fatal:true}};connection.replenishWindows(sink);return tsHTTP2Feed{consumed:offset}}
func (connection *tsHTTP2Connection)replenishWindows(sink tsHTTP2Sink){if connection.recvWindow.needsUpdate(){increment:=connection.recvWindow.takeUpdate();if increment>0{connection.sendWindowUpdate(sink,0,increment)}};buffer:=connection.replenishBuf[:0];for id,stream:=range connection.streams{if stream.state!=tsHTTP2StreamClosed&&stream.recvWindow.needsUpdate()&&sink.isStreamReading(id){if increment:=stream.recvWindow.takeUpdate();increment>0{buffer=append(buffer,[2]uint32{id,increment})}}};for _,item:=range buffer{connection.sendWindowUpdate(sink,item[0],item[1])};connection.replenishBuf=buffer;evict:=connection.evictBuf[:0];for id,stream:=range connection.streams{if stream.state==tsHTTP2StreamClosed{evict=append(evict,id)}};for _,id:=range evict{delete(connection.streams,id)};connection.evictBuf=evict}
func (connection *tsHTTP2Connection)dispatch(sink tsHTTP2Sink,header tsHTTP2FrameHeader,payload []byte)bool{if header.frameType!=7{connection.framesReceived++;sink.onFrameCounters(connection.framesReceived,connection.framesSent)};if inflight:=connection.headerBlockInFlight;inflight!=nil&&(header.frameType!=9||header.streamID!=inflight.continuationStream){connection.sendGoAway(sink,tsHTTP2ErrorProtocol,[]byte("expected CONTINUATION frame"));return true};validation:=tsHTTP2ValidateFrameHeader(header,connection.localSettings.maxFrameSize);switch validation.kind{case tsHTTP2HeaderConnectionError:connection.sendGoAway(sink,uint32(validation.code),[]byte("frame validation failed"));return true;case tsHTTP2HeaderStreamError:connection.sendRstStream(sink,uint32(validation.streamID),uint32(validation.code));return false};switch header.frameType{case 4:return connection.handleSettings(sink,header,payload);case 6:return connection.handlePing(sink,header,payload);case 7:return connection.handleGoAway(sink,payload);case 8:return connection.handleWindowUpdate(sink,header,payload);case 1:return connection.handleHeaders(sink,header,payload);case 9:return connection.handleContinuation(sink,header,payload);case 0:return connection.handleData(sink,header,payload);case 3:return connection.handleRstStream(sink,header,payload);case 5:return connection.handlePushPromise(sink,header,payload);case 10:return connection.handleAltSvc(sink,header,payload);case 12:return connection.handleOrigin(sink,header,payload);default:return false}}
func (connection *tsHTTP2Connection)handleSettings(sink tsHTTP2Sink,header tsHTTP2FrameHeader,payload []byte)bool{if tsHTTP2HasFrameFlag(header.flags,1){if header.length!=0{connection.sendGoAway(sink,tsHTTP2ErrorFrameSize,[]byte("SETTINGS ACK with payload"));return true};if connection.localSettingsAcked&&len(connection.pendingLocalSettingsAcks)==0{return false};connection.localSettingsAcked=true;acked:=connection.localSettings;if len(connection.pendingLocalSettingsAcks)>0{acked=connection.pendingLocalSettingsAcks[0].settings;connection.pendingLocalSettingsAcks=connection.pendingLocalSettingsAcks[1:]};connection.ackedLocalInitialWindow=acked.initialWindowSize;connection.enforcedMaxHeaderListSize=acked.maxHeaderListSize;sink.onLocalSettings(acked);return false};if uint32(len(payload)/6)>connection.maxSettings{connection.sendGoAway(sink,tsHTTP2ErrorEnhanceYourCalm,[]byte("SETTINGS: too many settings entries"));return true};if code:=tsHTTP2ValidateSettingsPayload(payload);code!=tsHTTP2ErrorNoError{connection.sendGoAway(sink,uint32(code),[]byte("SETTINGS value out of range"));return true};oldTable:=connection.remoteSettings.headerTableSize;oldInitialWindow:=connection.remoteSettings.initialWindowSize;for offset:=0;offset+6<=len(payload);offset+=6{id:=binary.BigEndian.Uint16(payload[offset:]);value:=binary.BigEndian.Uint32(payload[offset+2:]);if id==8&&value==0&&connection.remoteSettings.enableConnectProtocol==1{connection.sendGoAway(sink,tsHTTP2ErrorProtocol,[]byte("SETTINGS: server attempted to disable enableConnectProtocol"));return true};switch id{case 1,2,3,4,5,6,8:connection.remoteSettings.apply(id,value);case 9: /* RFC 7540 priorities extension is recognized and ignored. */ default:sink.onRemoteCustomSetting(id,value)}};if connection.remoteSettings.headerTableSize!=oldTable{connection.hpack.queueEncoderCapacity(connection.remoteSettings.headerTableSize)};if connection.remoteSettings.initialWindowSize!=oldInitialWindow{delta:=int64(connection.remoteSettings.initialWindowSize)-int64(oldInitialWindow);for _,stream:=range connection.streams{if stream.state!=tsHTTP2StreamClosed{stream.sendWindow.applyInitialDelta(delta)}}};connection.sendSettingsAck(sink);if connection.noteOutboundAck(sink){return true};sink.onRemoteSettings(connection.remoteSettings);return false}
func (connection *tsHTTP2Connection)handlePing(sink tsHTTP2Sink,header tsHTTP2FrameHeader,payload []byte)bool{if tsHTTP2HasFrameFlag(header.flags,1){sink.onPing(payload,true);return false};echo:=append([]byte(nil),payload[:8]...);connection.sendPingAck(sink,echo);sink.onPing(echo,false);return connection.noteOutboundAck(sink)}
func (connection *tsHTTP2Connection)handleGoAway(sink tsHTTP2Sink,payload []byte)bool{if len(payload)<8{connection.sendGoAway(sink,tsHTTP2ErrorFrameSize,[]byte("GOAWAY too short"));return true};lastStreamID:=binary.BigEndian.Uint32(payload[:4])&0x7fffffff;code:=binary.BigEndian.Uint32(payload[4:8]);initiatedLocally:=lastStreamID%2==0;if !connection.isServer{initiatedLocally=!initiatedLocally};if lastStreamID>0&&!initiatedLocally{connection.sendGoAway(sink,tsHTTP2ErrorProtocol,[]byte("GOAWAY: invalid last_stream_id"));return true};connection.goingAway=true;sink.onGoAway(code,lastStreamID,payload[8:]);return false}
func (connection *tsHTTP2Connection)handleWindowUpdate(sink tsHTTP2Sink,header tsHTTP2FrameHeader,payload []byte)bool{increment:=binary.BigEndian.Uint32(payload[:4])&0x7fffffff;if increment==0{if header.streamID==0{connection.sendGoAway(sink,tsHTTP2ErrorProtocol,[]byte("WINDOW_UPDATE with 0 increment"));return true};connection.sendRstStream(sink,header.streamID,tsHTTP2ErrorProtocol);if stream:=connection.streams[header.streamID];stream!=nil{stream.state=tsHTTP2StreamClosed};sink.onStreamReset(header.streamID,tsHTTP2ErrorProtocol);return false};if header.streamID==0{if code:=connection.sendWindow.increase(increment);code!=tsHTTP2ErrorNoError{connection.sendGoAway(sink,uint32(code),[]byte("connection flow-control window overflow"));return true}}else if stream:=connection.streams[header.streamID];stream!=nil{if code:=stream.sendWindow.increase(increment);code!=tsHTTP2ErrorNoError{connection.sendRstStream(sink,header.streamID,uint32(code));if stream=connection.streams[header.streamID];stream!=nil{stream.state=tsHTTP2StreamClosed};sink.onStreamReset(header.streamID,uint32(code));return false}};sink.onWindowUpdate(header.streamID,increment);return false}
func (connection *tsHTTP2Connection)handleHeaders(sink tsHTTP2Sink,header tsHTTP2FrameHeader,payload []byte)bool{offset,end:=0,len(payload);if tsHTTP2HasFrameFlag(header.flags,8){if len(payload)==0{connection.sendGoAway(sink,tsHTTP2ErrorFrameSize,[]byte("HEADERS padded but empty"));return true};padding:=int(payload[0]);offset=1;if offset+padding>end{connection.sendGoAway(sink,tsHTTP2ErrorProtocol,[]byte("HEADERS padding too large"));return true};end-=padding};if tsHTTP2HasFrameFlag(header.flags,32){if offset+5>end{connection.sendGoAway(sink,tsHTTP2ErrorFrameSize,[]byte("HEADERS priority truncated"));return true};offset+=5};endStream:=tsHTTP2HasFrameFlag(header.flags,1);sendInitial:=connection.remoteSettings.initialWindowSize;recvInitial:=connection.localSettings.initialWindowSize;_,exists:=connection.streams[header.streamID];isNew:=!exists;if isNew&&connection.isServer&&header.streamID%2==0{connection.sendGoAway(sink,tsHTTP2ErrorProtocol,[]byte("invalid stream id for HEADERS"));return true};refused:=isNew&&connection.isServer&&!sink.canOpenStream();disposition:=tsHTTP2BlockDeliver;if refused{disposition=tsHTTP2BlockRefused}else{stream:=connection.streams[header.streamID];if stream==nil{stream=tsHTTP2NewProtocolStream(sendInitial,recvInitial);connection.streams[header.streamID]=stream};if isNew{stream.openedAfterGoaway=sink.goawaySent()};current:=stream.state;event:=tsHTTP2EventRecvHeaders;if endStream{event=tsHTTP2EventRecvHeadersEndStream};next,transitionError:=tsHTTP2Transition(current,event);if transitionError==0{stream.state=next}else if transitionError==tsHTTP2TransitionProtocolError{connection.sendGoAway(sink,tsHTTP2ErrorProtocol,[]byte("HEADERS in invalid stream state"));return true}else if current==tsHTTP2StreamHalfClosedRemote{connection.localConnectionError(sink,tsHTTP2ErrorStreamClosed,-510,[]byte("HEADERS: stream closed"));return true}else{disposition=tsHTTP2BlockStreamClosed}};if isNew{if header.streamID>connection.lastStreamID{connection.lastStreamID=header.streamID};if !refused{sink.onStreamOpen(header.streamID)}};connection.headerBlock=append(connection.headerBlock[:0],payload[offset:end]...);meta:=tsHTTP2HeaderBlockMeta{target:header.streamID,flags:header.flags,endStream:endStream,isRequest:connection.isServer&&isNew,disposition:disposition};return connection.finishOrParkHeaderBlock(sink,header,meta)}
func (connection *tsHTTP2Connection)finishOrParkHeaderBlock(sink tsHTTP2Sink,header tsHTTP2FrameHeader,meta tsHTTP2HeaderBlockMeta)bool{if !tsHTTP2HasFrameFlag(header.flags,4){connection.headerBlockInFlight=&tsHTTP2HeaderBlockInFlight{continuationStream:header.streamID,meta:meta};return false};return connection.finishHeaderBlock(sink,meta)}
func (connection *tsHTTP2Connection)handleContinuation(sink tsHTTP2Sink,header tsHTTP2FrameHeader,payload []byte)bool{inflight:=connection.headerBlockInFlight;if inflight==nil{connection.sendGoAway(sink,tsHTTP2ErrorProtocol,[]byte("unexpected CONTINUATION frame"));return true};connection.headerBlockInFlight=nil;capSize:=max(int(connection.enforcedMaxHeaderListSize),65536);if len(connection.headerBlock)+len(payload)>capSize{connection.sendGoAway(sink,tsHTTP2ErrorCompression,[]byte("header block too large"));return true};connection.headerBlock=append(connection.headerBlock,payload...);if !tsHTTP2HasFrameFlag(header.flags,4){connection.headerBlockInFlight=inflight;return false};return connection.finishHeaderBlock(sink,inflight.meta)}
func (connection *tsHTTP2Connection)finishHeaderBlock(sink tsHTTP2Sink,meta tsHTTP2HeaderBlockMeta)bool{target,pushParent,hasPushParent,flags,endStream,isRequest,disposition:=meta.target,meta.pushParent,meta.hasPushParent,meta.flags,meta.endStream,meta.isRequest,meta.disposition;if hasPushParent{sink.onPushPromise(pushParent,target)};block:=connection.headerBlock;connection.headerBlock=nil;offset:=0;rejected:=false;streams:=connection.streams;isTrailer:=!hasPushParent&&disposition==tsHTTP2BlockDeliver&&streams[target]!=nil&&streams[target].recvFinalHeaders;malformed:=isTrailer&&!endStream;seenRegular:=false;seenPseudo:=uint8(0);sawConnect,sawHost,informational:=false,false,false;var contentLength *uint64;listSize,fieldCount:=0,0;for offset<len(block){header,errorCode:=connection.hpack.decode(block[offset:]);if errorCode!=0{connection.sendGoAway(sink,tsHTTP2ErrorCompression,[]byte("HPACK decode error"));connection.headerBlock=block[:0];return true};offset+=header.next;listSize+=len(header.name)+len(header.value)+32;fieldCount++;if rejected||uint32(listSize)>connection.enforcedMaxHeaderListSize||uint32(fieldCount)>connection.maxHeaderListPairs{rejected=true;continue};if disposition!=tsHTTP2BlockDeliver{continue};if !malformed{name,value:=header.name,header.value;if tsHTTP2IsMalformedHeaderName(name)||tsHTTP2IsMalformedHeaderValue(value){malformed=true}else if len(name)>0&&name[0]==':'{var bit uint8;switch string(name[1:]){case "method":bit=tsHTTP2PseudoMethod;case "scheme":bit=tsHTTP2PseudoScheme;case "authority":bit=tsHTTP2PseudoAuthority;case "path":bit=tsHTTP2PseudoPath;case "status":bit=tsHTTP2PseudoStatus;case "protocol":bit=tsHTTP2PseudoProtocol;default:bit=tsHTTP2PseudoUnknown};wrongDirection:=connection.isServer&&string(name[1:])=="status";protocolDisabled:=connection.isServer&&string(name[1:])=="protocol"&&connection.localSettings.enableConnectProtocol==0;if seenRegular||bit==tsHTTP2PseudoUnknown||seenPseudo&bit!=0||wrongDirection||protocolDisabled||isTrailer||len(value)==0{malformed=true};if string(name[1:])=="status"&&len(value)==3&&value[0]=='1'{informational=true};seenPseudo|=bit;if string(name[1:])=="method"&&string(value)=="CONNECT"{sawConnect=true}}else{seenRegular=true;switch string(name){case "connection","keep-alive","proxy-connection","transfer-encoding","upgrade":malformed=true;case "host":if isRequest{sawHost=true};case "te":if !strings.EqualFold(string(value),"trailers"){malformed=true};case "content-length":parsed,ok:=tsHTTP2ParseContentLengthStrict(value);if !ok||contentLength!=nil{malformed=true}else{contentLength=&parsed}}};if malformed{continue};sink.onHeader(target,header.name,header.value,header.neverIndexed)}};connection.headerBlock=block[:0];switch disposition{case tsHTTP2BlockRefused:connection.sendRstStream(sink,target,tsHTTP2ErrorEnhanceYourCalm);sink.onStreamRejected(target);return false;case tsHTTP2BlockStreamClosed:connection.sendRstStream(sink,target,tsHTTP2ErrorStreamClosed);if stream:=connection.streams[target];stream!=nil{stream.state=tsHTTP2StreamClosed};sink.onStreamReset(target,tsHTTP2ErrorStreamClosed);return false};if isRequest&&!rejected&&!malformed{extendedConnect:=seenPseudo&tsHTTP2PseudoProtocol!=0;if sawConnect&&!extendedConnect{malformed=seenPseudo&(tsHTTP2PseudoScheme|tsHTTP2PseudoPath)!=0||seenPseudo&tsHTTP2PseudoAuthority==0}else{malformed=seenPseudo&(tsHTTP2PseudoMethod|tsHTTP2PseudoScheme|tsHTTP2PseudoPath)!=(tsHTTP2PseudoMethod|tsHTTP2PseudoScheme|tsHTTP2PseudoPath)||(seenPseudo&tsHTTP2PseudoAuthority==0&&!sawHost)||(extendedConnect&&(!sawConnect||seenPseudo&tsHTTP2PseudoAuthority==0))}};if !hasPushParent&&connection.isServer&&!malformed&&!rejected{if stream:=connection.streams[target];stream!=nil{if !sawConnect&&stream.contentLength==nil&&contentLength!=nil{value:=*contentLength;stream.contentLength=&value};if endStream&&stream.contentLength!=nil&&*stream.contentLength!=stream.recvBodyBytes{malformed=true}}};if malformed&&!rejected{count:=connection.invalidFrameCount;if connection.invalidFrameCount<math.MaxUint32{connection.invalidFrameCount++};if count>connection.maxInvalidFrames{connection.terminated=true;sink.onTooManyInvalidFrames();return true};connection.sendRstStream(sink,target,tsHTTP2ErrorProtocol);if stream:=connection.streams[target];stream!=nil{stream.state=tsHTTP2StreamClosed};sink.onStreamReset(target,tsHTTP2ErrorProtocol);sink.onStreamRejected(target);return false};if rejected{connection.sendRstStream(sink,target,tsHTTP2ErrorEnhanceYourCalm);if stream:=connection.streams[target];stream!=nil{stream.state=tsHTTP2StreamClosed};sink.onStreamReset(target,tsHTTP2ErrorEnhanceYourCalm);sink.onStreamRejected(target);return false};if !hasPushParent&&!informational{if stream:=connection.streams[target];stream!=nil{stream.recvFinalHeaders=true}};sink.onHeadersComplete(target,endStream,flags);if endStream{if stream:=connection.streams[target];stream!=nil{sink.onStreamEnd(target,stream.state)}};return false}
`
