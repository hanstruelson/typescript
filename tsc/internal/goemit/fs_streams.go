package goemit

// File streams keep their mutable state on one loop and submit one native I/O
// operation at a time. Buffer storage remains native bytes across this boundary.
const fsStreamRuntime = `
type tsFSWriteChunk struct { data []byte; callback tsValue }
type tsFSStream struct {
    nativeClose func(*tsLoop)*tsPromise;loop *tsLoop; events *tsFSEvents; descriptor tsValue; options tsValue
    path,encoding string; reading,opened,busy,flowing,ending,finished,destroyed,closing bool
    autoClose,emitClose,flush,backpressure bool; position,end int64; highWaterMark,queued int
    pendingRead tsValue; decoderTail []byte; count int64; chunks []tsFSWriteChunk; waiters []*tsPromise; finalError tsValue
}
func tsFSCreateStream(loop *tsLoop,path tsValue,options tsValue,reading bool)tsValue{
    options=tsArg(tsFSPrepare(loop,[]tsValue{options}),0)
    stream:=&tsFSStream{loop:loop,events:tsFSEventObject(loop),options:options,reading:reading,descriptor:tsNull,position:-1,end:math.MaxInt64,autoClose:tsTruthy(tsFSOption(options,"autoClose",tsBooleanValue(true))),emitClose:tsTruthy(tsFSOption(options,"emitClose",tsBooleanValue(true))),flush:tsTruthy(tsFSProperty(options,"flush"))}
    if options.kind==tsObjectKind{stream.nativeClose=(*tsObject)(options.ref).streamClose}
    stream.highWaterMark=int(tsFSInteger(tsFSProperty(options,"highWaterMark"),"highWaterMark",64*1024,1))
    if !tsIsUndefined(tsFSProperty(options,"start")){stream.position=tsFSInteger(tsFSProperty(options,"start"),"start",0,0)}
    stream.end=tsFSInteger(tsFSProperty(options,"end"),"end",math.MaxInt64,0)
    if stream.position>=0&&stream.end<stream.position{panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","end must be at least start")})}
    if !tsNullish(tsFSProperty(options,"encoding")){stream.encoding=tsFSEncoding(options,"utf8")}
    object:=stream.events.object;class:="WriteStream";if reading{class="ReadStream"};object.nativeClass=tsNativeClass(class);object.set("fd",tsNull);object.set("pending",tsBooleanValue(true));object.set("destroyed",tsBooleanValue(false));object.set("closed",tsBooleanValue(false));object.set("bytesRead",tsNumberValue(0));object.set("bytesWritten",tsNumberValue(0));object.set("readable",tsBooleanValue(reading));object.set("writable",tsBooleanValue(!reading));object.set("writableLength",tsNumberValue(0));object.set("readableHighWaterMark",tsNumberValue(float64(stream.highWaterMark)));object.set("writableHighWaterMark",tsNumberValue(float64(stream.highWaterMark)))
    signal:=tsFSAbortSignal(options)
    if signal!=nil{cancel:=make(chan struct{});var once sync.Once;stop:=func(){once.Do(func(){close(cancel)})};task:=loop.submit(func()tsResult{select{case <-signal.done:return tsResult{tsBooleanValue(true),false};case <-cancel:return tsResult{tsBooleanValue(false),false}}},func(result tsResult){if tsTruthy(result.value){stream.destroyStream(tsFSAbortError(signal))}});task.unref=true;loop.resourceClosers=append(loop.resourceClosers,stop);stream.events.listeners["close"]=append(stream.events.listeners["close"],tsFSListener{callback:tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{stop();return tsU})),once:true})}
    descriptor:=tsFSProperty(options,"fd")
    if descriptor.kind==tsObjectKind{descriptor=tsFSProperty(descriptor,"fd")}
    if tsNullish(descriptor){stream.path=tsFSPath(path);object.set("path",path)}else{stream.descriptor=descriptor}
    object.set("destroy",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{stream.destroyStream(tsArg(args,0));return tsObjectValue(object)})))
    object.set("close",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{callback:=tsArg(args,0);if callback.kind==tsFunctionKind{stream.events.listeners["close"]=append(stream.events.listeners["close"],tsFSListener{callback:callback,once:true})};stream.autoClose=true;stream.destroyStream(tsU);return tsU})))
    if reading{
        originalOn:=object.get("on")
        on:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{result:=tsCall(loop,originalOn,args...);if tsText(tsArg(args,0))=="data"{stream.flowing=true;stream.pumpRead()};return result}))
        object.set("on",on);object.set("addListener",on)
        object.set("pause",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{stream.flowing=false;return tsObjectValue(object)})))
        object.set("resume",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{stream.flowing=true;stream.pumpRead();return tsObjectValue(object)})))
        object.set("isPaused",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsBooleanValue(!stream.flowing)})))
        object.set("setEncoding",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{stream.encoding=tsFSEncoding(tsArg(args,0),"utf8");return tsObjectValue(object)})))
        object.set("pipe",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{
            destination:=tsArg(args,0);pipeOptions:=tsArg(args,1)
            stream.events.listeners["data"]=append(stream.events.listeners["data"],tsFSListener{callback:tsFunctionValue(tsFunc(func(loop *tsLoop,values ...tsValue)tsValue{if !tsTruthy(tsCall(loop,tsFSProperty(destination,"write"),tsArg(values,0))){stream.flowing=false;tsCall(loop,tsFSProperty(destination,"once"),tsStringReference(tsStringUTF8("drain")),tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{stream.flowing=true;stream.pumpRead();return tsU})))};return tsU}))})
            if tsTruthy(tsFSOption(pipeOptions,"end",tsBooleanValue(true))){stream.events.listeners["end"]=append(stream.events.listeners["end"],tsFSListener{callback:tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{tsCall(loop,tsFSProperty(destination,"end"));return tsU})),once:true})}
            stream.flowing=true;stream.pumpRead();return destination
        })))
        object.asyncIterator=func()*tsAsyncIterator{return &tsAsyncIterator{
            next:func(loop *tsLoop)*tsPromise{if stream.destroyed&&!tsIsUndefined(stream.finalError){return loop.resolved(stream.finalError,true)};if stream.finished||stream.destroyed{return loop.resolved(tsIteratorResult(tsU,true),false)};promise:=loop.promise();stream.waiters=append(stream.waiters,promise);stream.pumpRead();return promise},
            close:func(loop *tsLoop)*tsPromise{stream.destroyStream(tsU);return loop.resolved(tsIteratorResult(tsU,true),false)},
        }}
    }else{
        object.set("write",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsBooleanValue(stream.enqueueWrite(args,false))})))
        object.set("end",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{stream.enqueueWrite(args,true);return tsObjectValue(object)})))
    }
    completeOpen:=func(result tsResult){object.set("pending",tsBooleanValue(false));if result.rejected{stream.destroyStream(result.value);return};stream.descriptor=result.value;stream.opened=true;object.set("fd",result.value);if stream.path!=""{stream.events.emitEvent("open",result.value)};stream.events.emitEvent("ready");if stream.destroyed{stream.closeStream();return};if reading{stream.pumpRead()}else{stream.pumpWrite()}}
    if tsNullish(descriptor){flag:="w";if reading{flag="r"};flagValue:=tsFSOption(options,"flags",tsStringReference(tsStringUTF8(flag)));mode:=tsFSProperty(options,"mode");loop.submit(func()tsResult{return tsFSExecute(loop,"open",[]tsValue{path,flagValue,mode},true)},completeOpen)}else{loop.ready=append(loop.ready,func(){completeOpen(tsResult{descriptor,false})})}
    return tsObjectValue(object)
}
func(stream *tsFSStream)destroyStream(failure tsValue){
    if stream.destroyed{return};stream.destroyed=true;stream.finalError=failure;stream.events.object.set("destroyed",tsBooleanValue(true));stream.events.object.set("readable",tsBooleanValue(false));stream.events.object.set("writable",tsBooleanValue(false))
    loop:=stream.loop
    for _,promise:=range stream.waiters{if tsIsUndefined(failure){promise.settle(tsResult{tsIteratorResult(tsU,true),false})}else{promise.settle(tsResult{failure,true})}};stream.waiters=nil
    for _,chunk:=range stream.chunks{if chunk.callback.kind==tsFunctionKind{err:=failure;if tsIsUndefined(err){err=tsFSException("Error","ERR_STREAM_DESTROYED","Stream was destroyed")};tsCall(loop,chunk.callback,err)}};stream.chunks=nil
    if !tsIsUndefined(failure){loop.ready=append(loop.ready,func(){stream.events.emitEvent("error",failure)})}
    if !stream.busy{stream.closeStream()}
}
func(stream *tsFSStream)closeStream(){
    if !stream.opened||stream.closing{return};stream.closing=true;loop:=stream.loop
    finish:=func(result tsResult){stream.events.object.set("closed",tsBooleanValue(true));stream.events.object.set("fd",tsNull);if result.rejected{stream.events.emitEvent("error",result.value)};if stream.emitClose{stream.events.emitEvent("close")}}
    if stream.autoClose{if stream.nativeClose!=nil{promise:=stream.nativeClose(loop);loop.await(tsPromiseValue(promise),finish)}else{loop.submit(func()tsResult{return tsFSExecute(loop,"close",[]tsValue{stream.descriptor},true)},finish)}}else{loop.ready=append(loop.ready,func(){finish(tsResult{tsU,false})})}
}
func(stream *tsFSStream)pumpRead(){
    if !stream.opened||stream.busy||stream.finished||stream.destroyed||(!stream.flowing&&len(stream.waiters)==0){return}
    loop:=stream.loop;
    if !tsIsUndefined(stream.pendingRead){data:=stream.pendingRead;stream.pendingRead=tsU;if len(stream.waiters)>0{promise:=stream.waiters[0];stream.waiters=stream.waiters[1:];promise.settle(tsResult{tsIteratorResult(data,false),false})}else{stream.events.emitEvent("data",data)};stream.pumpRead();return}
    size:=stream.highWaterMark;if stream.position>=0&&stream.end!=math.MaxInt64{remaining:=stream.end-stream.position+1;if remaining<int64(size){size=int(remaining)}}
    if size<=0{stream.finishRead();return};stream.busy=true
    buffer:=tsNodeBuffer(make([]byte,size));position:=tsNull;if stream.position>=0{position=tsNumberValue(float64(stream.position))}
    loop.submit(func()tsResult{return tsFSDescriptorIO("read",[]tsValue{stream.descriptor,buffer,tsNumberValue(0),tsNumberValue(float64(size)),position})},func(result tsResult){
        stream.busy=false;if stream.destroyed{stream.closeStream();return};if result.rejected{stream.destroyStream(result.value);return}
        n:=int(tsNumber(tsFSProperty(result.value,"bytesRead")));if n==0{stream.finishRead();return};stream.count+=int64(n);if stream.position>=0{stream.position+=int64(n)};stream.events.object.set("bytesRead",tsNumberValue(float64(stream.count)))
        data:=tsNodeBuffer(tsFSBytes(buffer)[:n]);if stream.encoding!=""{bytes:=append(stream.decoderTail,tsFSBytes(data)...);cut:=len(bytes);switch stream.encoding{case "utf8":for start:=len(bytes)-1;start>=0&&start>=len(bytes)-4;start--{if utf8.RuneStart(bytes[start]){if !utf8.FullRune(bytes[start:]){cut=start};break}};case "utf16le":cut-=cut%2;if cut>=2{last:=binary.LittleEndian.Uint16(bytes[cut-2:]);if last>=0xd800&&last<=0xdbff{cut-=2}};case "base64","base64url":cut-=cut%3};stream.decoderTail=append([]byte{},bytes[cut:]...);if cut==0{stream.pumpRead();return};data=tsFSDecode(bytes[:cut],stream.encoding)}
        if len(stream.waiters)>0{promise:=stream.waiters[0];stream.waiters=stream.waiters[1:];promise.settle(tsResult{tsIteratorResult(data,false),false})}else if stream.flowing{stream.events.emitEvent("data",data)}else{stream.pendingRead=data}
        stream.pumpRead()
    })
}
func(stream *tsFSStream)finishRead(){if stream.finished{return};if len(stream.decoderTail)>0{tail:=tsFSDecode(stream.decoderTail,stream.encoding);stream.decoderTail=nil;if len(stream.waiters)>0{promise:=stream.waiters[0];stream.waiters=stream.waiters[1:];promise.settle(tsResult{tsIteratorResult(tail,false),false})}else if stream.flowing{stream.events.emitEvent("data",tail)}};stream.finished=true;stream.events.object.set("readableEnded",tsBooleanValue(true));stream.events.object.set("readable",tsBooleanValue(false));for _,promise:=range stream.waiters{promise.settle(tsResult{tsIteratorResult(tsU,true),false})};stream.waiters=nil;stream.events.emitEvent("end");if stream.autoClose{stream.destroyStream(tsU)}}
func(stream *tsFSStream)enqueueWrite(args []tsValue,end bool)bool{
    loop:=stream.loop;data:=tsArg(args,0);encoding:=tsArg(args,1);callback:=tsArg(args,2);if encoding.kind==tsFunctionKind{callback=encoding;encoding=tsU};if data.kind==tsFunctionKind&&end{callback=data;data=tsU}
    if !tsIsUndefined(callback)&&callback.kind!=tsFunctionKind{tsFSInvalid("callback","a function")}
    if stream.ending||stream.destroyed{failure:=tsFSException("Error","ERR_STREAM_WRITE_AFTER_END","write after end");if callback.kind==tsFunctionKind{loop.ready=append(loop.ready,func(){tsCall(loop,callback,failure)})};stream.destroyStream(failure);return false}
    if !tsIsUndefined(data){options:=tsNewObject();if tsIsUndefined(encoding){encoding=tsStringReference(tsStringUTF8("utf8"))};options.set("encoding",encoding);bytes:=tsFSData(data,tsObjectValue(options));stream.chunks=append(stream.chunks,tsFSWriteChunk{data:bytes,callback:callback});stream.queued+=len(bytes)}
    if end{stream.ending=true;stream.events.object.set("writableEnded",tsBooleanValue(true));if callback.kind==tsFunctionKind{stream.events.listeners["finish"]=append(stream.events.listeners["finish"],tsFSListener{callback:callback,once:true});if !tsIsUndefined(data){stream.chunks[len(stream.chunks)-1].callback=tsU}}}
    stream.events.object.set("writableLength",tsNumberValue(float64(stream.queued)));ok:=stream.queued<stream.highWaterMark;if !ok{stream.backpressure=true};stream.pumpWrite();return ok
}
func(stream *tsFSStream)pumpWrite(){
    if !stream.opened||stream.busy||stream.finished||stream.destroyed{return};loop:=stream.loop
    if len(stream.chunks)==0{
        if stream.backpressure{stream.backpressure=false;stream.events.emitEvent("drain")}
        if stream.ending{stream.finished=true;finish:=func(result tsResult){if result.rejected{stream.destroyStream(result.value);return};stream.events.object.set("writableFinished",tsBooleanValue(true));stream.events.object.set("writable",tsBooleanValue(false));stream.events.emitEvent("finish");if stream.autoClose{stream.destroyStream(tsU)}};if stream.flush{loop.submit(func()tsResult{return tsFSExecute(loop,"fsync",[]tsValue{stream.descriptor},true)},finish)}else{finish(tsResult{tsU,false})}};return
    }
    chunk:=stream.chunks[0];stream.chunks=stream.chunks[1:];stream.busy=true;position:=stream.position
    loop.submit(func()tsResult{record:=tsFSDescriptor(stream.descriptor);record.Lock();defer record.Unlock();written:=0;for written<len(chunk.data){var n int;var err error;if position<0{n,err=record.file.Write(chunk.data[written:])}else{n,err=record.file.WriteAt(chunk.data[written:],position+int64(written))};written+=n;if err!=nil{return tsResult{tsFSError(err,"write",record.path,""),true}};if n==0{return tsResult{tsFSError(io.ErrShortWrite,"write",record.path,""),true}}};return tsResult{tsNumberValue(float64(written)),false}},func(result tsResult){
        stream.busy=false;stream.queued-=len(chunk.data);stream.events.object.set("writableLength",tsNumberValue(float64(stream.queued)))
        if chunk.callback.kind==tsFunctionKind{err:=tsNull;if result.rejected{err=result.value};tsCall(loop,chunk.callback,err)}
        if stream.destroyed{stream.closeStream();return};if result.rejected{stream.destroyStream(result.value);return};n:=int64(tsNumber(result.value));stream.count+=n;if stream.position>=0{stream.position+=n};stream.events.object.set("bytesWritten",tsNumberValue(float64(stream.count)));stream.pumpWrite()
    })
}
func tsFSInstallStreams(fs *tsObject){
    fs.set("createReadStream",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsFSCreateStream(loop,tsArg(args,0),tsArg(args,1),true)})))
    fs.set("createWriteStream",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsFSCreateStream(loop,tsArg(args,0),tsArg(args,1),false)})))
}
`
