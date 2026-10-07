package goemit

const fsReaderRuntime = `
type tsFSLineReader struct {
    events *tsFSEvents;source *tsAsyncIterator;lines []string;buffer strings.Builder
    waiters []*tsPromise;closed,pulling,flowing,skipLF bool;failure tsValue
}
func(reader *tsFSLineReader)deliverLines(){
    for len(reader.lines)>0&&(len(reader.waiters)>0||reader.flowing){line:=reader.lines[0];reader.lines=reader.lines[1:];value:=tsStringReference(tsStringUTF8(line));if len(reader.waiters)>0{promise:=reader.waiters[0];reader.waiters=reader.waiters[1:];promise.settle(tsResult{tsIteratorResult(value,false),false})};if reader.flowing{reader.events.emitEvent("line",value)}}
    if reader.closed&&len(reader.lines)==0{for _,promise:=range reader.waiters{if tsIsUndefined(reader.failure){promise.settle(tsResult{tsIteratorResult(tsU,true),false})}else{promise.settle(tsResult{reader.failure,true})}};reader.waiters=nil}
}
func(reader *tsFSLineReader)pumpLines(loop *tsLoop){
    reader.deliverLines();if reader.closed||reader.pulling||(!reader.flowing&&len(reader.waiters)==0){return};reader.pulling=true
    promise:=reader.source.next(loop)
    loop.await(tsPromiseValue(promise),func(result tsResult){reader.pulling=false
        if reader.closed{return};if result.rejected{reader.failure=result.value;reader.closed=true;reader.deliverLines();if len(reader.events.listeners["error"])>0{reader.events.emitEvent("error",result.value)};reader.events.emitEvent("close");return}
        if tsTruthy(tsFSProperty(result.value,"done")){if reader.buffer.Len()>0{reader.lines=append(reader.lines,reader.buffer.String());reader.buffer.Reset()};reader.closed=true;reader.deliverLines();reader.events.emitEvent("close");return}
        value:=tsFSProperty(result.value,"value");text:="";if value.kind==tsStringKind{text=tsText(value)}else{text=tsText(tsFSDecode(tsFSBytes(value),"utf8"))}
        for _,r:=range text{if reader.skipLF{reader.skipLF=false;if r=='\n'{continue}};if r=='\r'||r=='\n'{reader.lines=append(reader.lines,reader.buffer.String());reader.buffer.Reset();reader.skipLF=r=='\r'}else{reader.buffer.WriteRune(r)}};reader.pumpLines(loop)
    })
}
func tsFSReadLines(loop *tsLoop,input tsValue)tsValue{
    reader:=&tsFSLineReader{events:tsFSEventObject(loop),source:tsAsyncIterate(input)};out:=reader.events.object
    on:=out.get("on");out.set("on",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{value:=tsCall(loop,on,args...);if tsText(tsArg(args,0))=="line"{reader.flowing=true;reader.pumpLines(loop)};return value})))
    out.set("close",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{if !reader.closed{reader.closed=true;reader.lines=nil;reader.source.close(loop);reader.deliverLines();reader.events.emitEvent("close")};return tsU})))
    out.set("pause",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{reader.flowing=false;return tsObjectValue(out)})))
    out.set("resume",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{reader.flowing=true;reader.pumpLines(loop);return tsObjectValue(out)})))
    out.asyncIterator=func()*tsAsyncIterator{return &tsAsyncIterator{
        next:func(loop *tsLoop)*tsPromise{if reader.closed&&len(reader.lines)==0{if !tsIsUndefined(reader.failure){return loop.resolved(reader.failure,true)};return loop.resolved(tsIteratorResult(tsU,true),false)};promise:=loop.promise();reader.waiters=append(reader.waiters,promise);reader.pumpLines(loop);return promise},
        close:func(loop *tsLoop)*tsPromise{wasClosed:=reader.closed;reader.closed=true;reader.lines=nil;reader.deliverLines();promise:=loop.promise();loop.await(tsPromiseValue(reader.source.close(loop)),func(result tsResult){if !result.rejected{result.value=tsIteratorResult(tsU,true)};promise.settle(result)});if !wasClosed{reader.events.emitEvent("close")};return promise},
    }};return tsObjectValue(out)
}
`
