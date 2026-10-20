package goemit

// nodeHTTP2StreamRuntime ports the HTTP/2 PUSH, extension-frame, and outbound
// stream operations from Bun's connection engine.
const nodeHTTP2StreamRuntime = `
func (connection *tsHTTP2Connection)handlePushPromise(sink tsHTTP2Sink,header tsHTTP2FrameHeader,payload []byte)bool{if connection.isServer{connection.sendGoAway(sink,tsHTTP2ErrorProtocol,[]byte("server received PUSH_PROMISE"));return true};if connection.localSettings.enablePush==0{connection.sendGoAway(sink,tsHTTP2ErrorProtocol,[]byte("PUSH_PROMISE with push disabled"));return true};offset,end:=0,len(payload);if tsHTTP2HasFrameFlag(header.flags,8){if len(payload)==0{connection.sendGoAway(sink,tsHTTP2ErrorFrameSize,[]byte("PUSH_PROMISE padded but empty"));return true};padding:=int(payload[0]);offset=1;if offset+padding>end{connection.sendGoAway(sink,tsHTTP2ErrorProtocol,[]byte("PUSH_PROMISE padding too large"));return true};end-=padding};if offset+4>end{connection.sendGoAway(sink,tsHTTP2ErrorFrameSize,[]byte("PUSH_PROMISE missing promised id"));return true};promised:=binary.BigEndian.Uint32(payload[offset:offset+4])&0x7fffffff;offset+=4;if promised==0||promised&1==1||promised<=connection.lastStreamID||connection.streams[promised]!=nil{connection.sendGoAway(sink,tsHTTP2ErrorProtocol,[]byte("PUSH_PROMISE invalid promised stream id"));return true};stream:=tsHTTP2NewProtocolStream(connection.remoteSettings.initialWindowSize,connection.localSettings.initialWindowSize);stream.state=tsHTTP2StreamReservedRemote;connection.streams[promised]=stream;connection.lastStreamID=promised;connection.headerBlock=append(connection.headerBlock[:0],payload[offset:end]...);meta:=tsHTTP2HeaderBlockMeta{target:promised,pushParent:header.streamID,hasPushParent:true,isRequest:true,disposition:tsHTTP2BlockDeliver};return connection.finishOrParkHeaderBlock(sink,header,meta)}
func (connection *tsHTTP2Connection)handleAltSvc(sink tsHTTP2Sink,header tsHTTP2FrameHeader,payload []byte)bool{if len(payload)<2{return false};originLength:=int(binary.BigEndian.Uint16(payload[:2]));if 2+originLength>len(payload){return false};origin,value:=payload[2:2+originLength],payload[2+originLength:];if connection.isServer||(header.streamID==0&&len(origin)==0)||(header.streamID!=0&&len(origin)!=0){return false};sink.onAltSvc(header.streamID,origin,value);return false}
func (connection *tsHTTP2Connection)handleOrigin(sink tsHTTP2Sink,header tsHTTP2FrameHeader,payload []byte)bool{if header.streamID!=0||connection.isServer{return false};sink.onOrigin(payload);return false}
func (connection *tsHTTP2Connection)beginHeaderBlock(){connection.encBuf=connection.encBuf[:0];var update [6]byte;if count:=connection.hpack.takePendingSizeUpdate(update[:],0);count>0{connection.encBuf=append(connection.encBuf,update[:count]...)}}
func (connection *tsHTTP2Connection)encodeHeader(name,value []byte,neverIndexed bool)bool{old:=len(connection.encBuf);connection.encBuf=append(connection.encBuf,make([]byte,len(name)+len(value)+64)...);written,errCode:=connection.hpack.encode(name,value,neverIndexed,connection.encBuf,old);if errCode!=0{connection.encBuf=connection.encBuf[:old];return false};connection.encBuf=connection.encBuf[:old+written];return true}
type tsHTTP2Priority struct{ parent uint32; weight uint8; exclusive bool }
func (connection *tsHTTP2Connection)sendHeaderBlock(sink tsHTTP2Sink,streamID uint32,endStream bool,priority *tsHTTP2Priority,paddingStrategy uint8){
	block:=connection.encBuf
	connection.encBuf=nil
	maxFrame:=max(int(connection.remoteSettings.maxFrameSize),1)
	priorityBytes:=0
	flags:=uint8(0)
	prefix:=[]byte(nil)
	if priority!=nil{
		priorityBytes=tsHTTP2StreamPrioritySize
		flags|=tsHTTP2FlagPriorityPresent
		wire:=tsHTTP2StreamPriorityWire{streamIdentifier:tsHTTP2InitUInt31(priority.parent,priority.exclusive),weight:priority.weight-1}
		prefix=make([]byte,tsHTTP2StreamPrioritySize)
		wire.write(prefix)
	}
	available:=max(maxFrame-priorityBytes,0)
	padding:=uint8(0)
	if len(block)<available{padding=tsHTTP2PaddingFor(paddingStrategy,len(block),available-1)}
	paddingOverhead:=0
	if padding>0{paddingOverhead=1+int(padding);flags|=tsHTTP2FlagPadded}
	firstLength:=min(len(block),max(available-paddingOverhead,0))
	if firstLength==len(block){flags|=tsHTTP2FlagEndHeaders}else{padding=0;paddingOverhead=0;flags &^= tsHTTP2FlagPadded;firstLength=available}
	if endStream{flags|=tsHTTP2FlagEndStream}
	payload:=make([]byte,0,priorityBytes+paddingOverhead+firstLength)
	if padding>0{payload=append(payload,padding)}
	payload=append(payload,prefix...)
	payload=append(payload,block[:firstLength]...)
	if padding>0{payload=append(payload,make([]byte,int(padding))...)}
	frames:=make([]byte,0,tsHTTP2FrameHeaderSize+len(payload)+(len(block)-firstLength)+((len(block)-firstLength)/maxFrame+1)*tsHTTP2FrameHeaderSize)
	header:=tsHTTP2WriteFrameHeader(tsHTTP2FrameHeader{length:uint32(len(payload)),frameType:tsHTTP2FrameHeaders,flags:flags,streamID:streamID})
	frames=append(frames,header...)
	frames=append(frames,payload...)
	frameCount:=uint64(1)
	offset:=firstLength
	for offset<len(block){length:=min(len(block)-offset,maxFrame);continuationFlags:=uint8(0);if offset+length==len(block){continuationFlags=tsHTTP2FlagEndHeaders};header:=tsHTTP2WriteFrameHeader(tsHTTP2FrameHeader{length:uint32(length),frameType:tsHTTP2FrameContinuation,flags:continuationFlags,streamID:streamID});frames=append(frames,header...);frames=append(frames,block[offset:offset+length]...);frameCount++;offset+=length}
	connection.writeFrameBatch(sink,frames,frameCount)
	connection.encBuf=block[:0]
	stream:=connection.streams[streamID]
	if stream==nil{stream=tsHTTP2NewProtocolStream(connection.remoteSettings.initialWindowSize,connection.localSettings.initialWindowSize);connection.streams[streamID]=stream}
	event:=tsHTTP2EventSendHeaders
	if endStream{event=tsHTTP2EventSendHeadersEndStream}
	if next,err:=tsHTTP2Transition(stream.state,event);err==0{stream.state=next}
	if connection.isLocalStreamID(streamID)&&streamID>connection.lastLocalStreamID{connection.lastLocalStreamID=streamID}
}
func tsHTTP2PaddingFor(strategy uint8,frameLen,maxLen int)uint8{if maxLen<frameLen{return 0};switch strategy{case 1:diff:=(frameLen+tsHTTP2FrameHeaderSize)%8;if diff==0{return 0};padded:=min(frameLen+(8-diff),maxLen);return uint8(min(max(padded-frameLen,0),255));case 2:return uint8(min(maxLen-frameLen,255));default:return 0}}
func (connection *tsHTTP2Connection)sendData(sink tsHTTP2Sink,streamID uint32,data []byte,endStream bool,paddingStrategy uint8)int{
	stream:=connection.streams[streamID]
	if stream==nil||!tsHTTP2CanSendData(stream.state){return 0}
	if len(data)==0{
		flags:=uint8(0)
		if endStream{flags|=tsHTTP2FlagEndStream}
		connection.writeFrame(sink,tsHTTP2FrameData,flags,streamID,nil)
		if endStream{if next,err:=tsHTTP2Transition(stream.state,tsHTTP2EventSendEndStream);err==0{stream.state=next}}
		return 0
	}
	available:=max(connection.sendWindow.available(),int64(0))
	available=min(available,max(stream.sendWindow.available(),int64(0)))
	if available==0{return 0}
	maxFrame:=int64(max(connection.remoteSettings.maxFrameSize,1))
	frames:=make([]byte,0)
	framesSent:=uint64(0)
	consumed:=int64(0)
	dataSent:=0
	for dataSent<len(data)&&available>0{
		maxSize:=min(maxFrame,available)
		size:=min(len(data)-dataSent,int(maxSize))
		padding:=uint8(0)
		if paddingStrategy!=0{padding=tsHTTP2PaddingFor(paddingStrategy,size,int(maxSize)-1)}
		payloadSize:=size
		if padding>0{payloadSize+=1+int(padding)}
		frameFlags:=uint8(0)
		last:=dataSent+size==len(data)
		if endStream&&last{frameFlags|=tsHTTP2FlagEndStream}
		if padding>0{frameFlags|=tsHTTP2FlagPadded}
		payload:=make([]byte,0,payloadSize)
		if padding>0{payload=append(payload,padding)}
		payload=append(payload,data[dataSent:dataSent+size]...)
		if padding>0{payload=append(payload,make([]byte,int(padding))...)}
		header:=tsHTTP2WriteFrameHeader(tsHTTP2FrameHeader{length:uint32(payloadSize),frameType:tsHTTP2FrameData,flags:frameFlags,streamID:streamID})
		frames=append(frames,header...)
		frames=append(frames,payload...)
		framesSent++
		dataSent+=size
		consumed+=int64(payloadSize)
		available-=int64(payloadSize)
	}
	if framesSent==0{return 0}
	connection.writeFrameBatch(sink,frames,framesSent)
	connection.sendWindow.consume(consumed)
	stream.sendWindow.consume(consumed)
	if endStream&&dataSent==len(data){if next,err:=tsHTTP2Transition(stream.state,tsHTTP2EventSendEndStream);err==0{stream.state=next}}
	return dataSent
}
func (connection *tsHTTP2Connection)closeStream(streamID uint32){delete(connection.streams,streamID)}
func (connection *tsHTTP2Connection)replenishStream(sink tsHTTP2Sink,streamID uint32){stream:=connection.streams[streamID];if stream==nil||stream.state==tsHTTP2StreamClosed||!stream.recvWindow.needsUpdate(){return};if increment:=stream.recvWindow.takeUpdate();increment>0{connection.sendWindowUpdate(sink,streamID,increment)}}
func (connection *tsHTTP2Connection)sendPushPromise(sink tsHTTP2Sink,parentID,promisedID uint32){block:=connection.encBuf;connection.encBuf=nil;maxFrame:=max(int(connection.remoteSettings.maxFrameSize),5);firstCapacity:=maxFrame-4;firstLength:=min(len(block),firstCapacity);first:=make([]byte,4+firstLength);binary.BigEndian.PutUint32(first,promisedID&tsHTTP2MaxStreamID);copy(first[4:],block[:firstLength]);flags:=uint8(0);if firstLength==len(block){flags=tsHTTP2FlagEndHeaders};frames:=make([]byte,0,tsHTTP2FrameHeaderSize+len(first)+(len(block)-firstLength)+((len(block)-firstLength)/maxFrame+1)*tsHTTP2FrameHeaderSize);header:=tsHTTP2WriteFrameHeader(tsHTTP2FrameHeader{length:uint32(len(first)),frameType:tsHTTP2FramePushPromise,flags:flags,streamID:parentID});frames=append(frames,header...);frames=append(frames,first...);frameCount:=uint64(1);offset:=firstLength;for offset<len(block){length:=min(len(block)-offset,maxFrame);continuationFlags:=uint8(0);if offset+length==len(block){continuationFlags=tsHTTP2FlagEndHeaders};header:=tsHTTP2WriteFrameHeader(tsHTTP2FrameHeader{length:uint32(length),frameType:tsHTTP2FrameContinuation,flags:continuationFlags,streamID:parentID});frames=append(frames,header...);frames=append(frames,block[offset:offset+length]...);frameCount++;offset+=length};connection.writeFrameBatch(sink,frames,frameCount);connection.encBuf=block[:0];stream:=connection.streams[promisedID];if stream==nil{stream=tsHTTP2NewProtocolStream(connection.remoteSettings.initialWindowSize,connection.localSettings.initialWindowSize);connection.streams[promisedID]=stream};stream.state=tsHTTP2StreamReservedLocal;if promisedID>connection.lastLocalStreamID{connection.lastLocalStreamID=promisedID}}
`
