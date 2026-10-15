package goemit

const fsHandleRuntime = `
type tsFSFileRecord struct {
    sync.Mutex
    file *os.File
    fd int
    path string
    closed bool
}

var tsFSFiles struct {
    sync.Mutex
    files map[int]*tsFSFileRecord
}

func tsFSInitFiles() {
    if tsFSFiles.files!=nil{return}
    tsFSFiles.files=map[int]*tsFSFileRecord{
        0:{file:os.Stdin,fd:0},1:{file:os.Stdout,fd:1},2:{file:os.Stderr,fd:2},
    }
}

func tsFSRegister(file *os.File, path string) int {
    fd:=int(file.Fd());tsFSFiles.Lock();defer tsFSFiles.Unlock();tsFSInitFiles()
    tsFSFiles.files[fd]=&tsFSFileRecord{file:file,fd:fd,path:path};return fd
}

func tsFSDescriptor(value tsValue) *tsFSFileRecord {
    if value.kind==tsObjectKind {value=tsFSProperty(value,"fd")}
    fd:=int(tsFSInteger(value,"fd",-1,0))
    tsFSFiles.Lock();tsFSInitFiles();record:=tsFSFiles.files[fd];tsFSFiles.Unlock()
    if record==nil{panic(tsThrown{tsFSError(syscall.EBADF,"fd","","")})}
    return record
}

func tsFSFile(value tsValue, flags int, mode os.FileMode) (*tsFSFileRecord,bool,error) {
    if tsIsNumeric(value)||value.kind==tsObjectKind&&!tsIsUndefined(tsFSProperty(value,"fd")){return tsFSDescriptor(value),false,nil}
    path:=tsFSPath(value);file,err:=os.OpenFile(path,flags,mode)
    if err!=nil{return nil,false,err}
    return &tsFSFileRecord{file:file,fd:int(file.Fd()),path:path},true,nil
}

func tsFSPosition(value tsValue) int64 {
    if tsNullish(value){return -1}
    return tsFSInteger(value,"position",-1,-1)
}

func tsFSDescriptorIO(operation string, args []tsValue) tsResult {
    record:=tsFSDescriptor(tsArg(args,0));record.Lock();defer record.Unlock()
    if record.closed{return tsResult{tsFSError(syscall.EBADF,operation,"",""),true}}
    source:=tsArg(args,1);original:=source;position:=int64(-1)
    n:=0;var err error
    out:=tsNewObject()
    if operation=="readv"||operation=="writev"{
        buffers:=[][]byte{};it:=tsIterate(source);for it.next(){buffers=append(buffers,tsFSBytes(it.value))}
        position=tsFSPosition(tsArg(args,2))
        if operation=="readv"{if position<0{n,err=unix.Readv(record.fd,buffers)}else{n,err=unix.Preadv(record.fd,buffers,position)}}else{if position<0{n,err=unix.Writev(record.fd,buffers)}else{n,err=unix.Pwritev(record.fd,buffers,position)}}
        out.set("buffers",source)
    }else{
        var data []byte
        if operation=="write"&&source.kind==tsStringKind{
            encoding:=tsArg(args,3);positionValue:=tsArg(args,2)
            if positionValue.kind==tsStringKind{encoding=positionValue;positionValue=tsNull}
            position=tsFSPosition(positionValue);data=tsFSData(source,encoding)
        }else{
            offsetValue,lengthValue,positionValue:=tsArg(args,2),tsArg(args,3),tsArg(args,4)
            if source.kind==tsObjectKind||tsIsUndefined(source){
                options:=source;source=tsFSProperty(options,"buffer")
                if tsIsUndefined(source){source=tsNodeBuffer(make([]byte,16384))}
                original=source;offsetValue=tsFSProperty(options,"offset");lengthValue=tsFSProperty(options,"length");positionValue=tsFSProperty(options,"position")
            }else if offsetValue.kind==tsObjectKind{
                options:=offsetValue;offsetValue=tsFSProperty(options,"offset");lengthValue=tsFSProperty(options,"length");positionValue=tsFSProperty(options,"position")
            }
            data=tsFSBytes(source)
            offset:=int(tsFSInteger(offsetValue,"offset",0,0))
            if offset>len(data){panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","offset exceeds buffer length")})}
            length:=int(tsFSInteger(lengthValue,"length",int64(len(data)-offset),0))
            if length>len(data)-offset{panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","length exceeds buffer length")})}
            data=data[offset:offset+length];position=tsFSPosition(positionValue)
        }
        if operation=="read"{if position<0{n,err=record.file.Read(data)}else{n,err=record.file.ReadAt(data,position)}}else{if position<0{n,err=record.file.Write(data)}else{n,err=record.file.WriteAt(data,position)}}
        out.set("buffer",original)
    }
    if errors.Is(err,io.EOF){err=nil}
    if err!=nil{return tsResult{tsFSError(err,operation,record.path,""),true}}
    key:="bytesRead";if strings.HasPrefix(operation,"write"){key="bytesWritten"}
    out.set(key,tsNumberValue(float64(n)));return tsResult{tsObjectValue(out),false}
}

func tsFSFileHandle(loop *tsLoop, descriptor tsValue) tsValue {
    events:=tsFSEventObject(loop);out:=events.object;out.nativeClass=tsNativeClass("FileHandle");out.set("fd",descriptor)
    pending:=0;closeSubmitted:=false;webLocked:=false;var webClose func(*tsLoop)*tsPromise;var closing *tsPromise;var closeNow func(*tsLoop)
    closeNow=func(loop *tsLoop){if closing==nil||pending!=0||closeSubmitted{return};closeSubmitted=true;value:=tsFSInvoke(loop,"close","promise",[]tsValue{descriptor});inner:=(*tsPromise)(value.ref);inner.then(func(result tsResult){closing.settle(result);if !result.rejected{events.emitEvent("close")}})}
    methods:=map[string]string{
        "appendFile":"appendFile","chmod":"fchmod","chown":"fchown","datasync":"fdatasync",
        "read":"read","readv":"readv","readFile":"readFile","stat":"fstat","sync":"fsync",
        "truncate":"ftruncate","utimes":"futimes","write":"write","writev":"writev","writeFile":"writeFile",
    }
    for name,operation:=range methods{
        op:=operation
        out.set(name,tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{
            if closing!=nil||tsNumber(out.get("fd"))<0{return tsPromiseValue(loop.resolved(tsFSError(syscall.EBADF,op,"",""),true))}
            values:=append([]tsValue{descriptor},args...);pending++;value:=tsFSInvoke(loop,op,"promise",values);inner:=(*tsPromise)(value.ref);promise:=loop.promise();inner.then(func(result tsResult){pending--;promise.settle(result);if pending==0&&closing!=nil{closeNow(loop)}});return tsPromiseValue(promise)
        })))
    }
    out.set("close",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{
        if closing!=nil{return tsPromiseValue(closing)}
        if tsNumber(out.get("fd"))<0{return tsPromiseValue(loop.resolved(tsU,false))}
        closing=loop.promise();out.set("fd",tsNumberValue(-1));if webClose!=nil{cancel:=webClose;webClose=nil;cancel(loop)};closeNow(loop);return tsPromiseValue(closing)
    })))
    for _,method:=range []string{"createReadStream","createWriteStream"}{name:=method;out.set(name,tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{
        if closing!=nil{panic(tsThrown{tsFSError(syscall.EBADF,name,"","")})};options:=tsNewObject();if !tsIsUndefined(tsArg(args,0)){tsObjectSpread(options,tsArg(args,0))};options.set("fd",descriptor)
        pending++;released:=false;autoClose:=tsTruthy(tsFSOption(tsObjectValue(options),"autoClose",tsBooleanValue(true)))
        if autoClose{options.streamClose=func(loop *tsLoop)*tsPromise{if !released{released=true;pending--};if closing!=nil&&pending==0{closeNow(loop)};value:=tsCall(loop,out.get("close"));return (*tsPromise)(value.ref)}}
        stream:=tsFSCreateStream(loop,tsU,tsObjectValue(options),name=="createReadStream")
        release:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{if released{return tsU};released=true;pending--;if autoClose{out.set("fd",tsNumberValue(-1));if closing!=nil&&pending==0{closing.settle(tsResult{tsU,false})};events.emitEvent("close")}else if closing!=nil&&pending==0{closeNow(loop)};return tsU}))
        event:="close";if !autoClose{event="finish";if name=="createReadStream"{event="end"}};tsCall(loop,tsFSProperty(stream,"once"),tsStringReference(tsStringUTF8(event)),release)
        if !autoClose{tsCall(loop,tsFSProperty(stream,"once"),tsStringReference(tsStringUTF8("close")),release)};return stream
    })))}
    out.set("readableWebStream",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{
        if closing!=nil||tsNumber(out.get("fd"))<0||webLocked{panic(tsThrown{tsFSException("Error","ERR_INVALID_STATE","FileHandle is closed, closing, or already locked")})};webLocked=true;options:=tsArg(args,0);autoClose:=tsTruthy(tsFSOption(options,"autoClose",tsBooleanValue(false)));pending++;ended:=false;released:=false
        release:=func(loop *tsLoop)*tsPromise{ended=true;webClose=nil;if !released{released=true;pending--};if closing!=nil&&pending==0{closeNow(loop)};if autoClose{value:=tsCall(loop,out.get("close"));return (*tsPromise)(value.ref)};return loop.resolved(tsU,false)}
        iterator:=&tsAsyncIterator{
            next:func(loop *tsLoop)*tsPromise{
                if ended{return loop.resolved(tsIteratorResult(tsU,true),false)};buffer:=tsNodeBuffer(make([]byte,16384));(*tsTypedArray)(buffer.ref).nodeBuffer=false;promise:=loop.promise();read:=tsCall(loop,out.get("read"),buffer,tsNumberValue(0),tsNumberValue(16384),tsNull)
                loop.await(read,func(result tsResult){if ended{promise.settle(tsResult{tsIteratorResult(tsU,true),false});return};if result.rejected{release(loop);promise.settle(result);return};n:=int(tsNumber(tsFSProperty(result.value,"bytesRead")));if n==0{closing:=release(loop);loop.await(tsPromiseValue(closing),func(result tsResult){if !result.rejected{result.value=tsIteratorResult(tsU,true)};promise.settle(result)});return};value:=tsNodeBuffer(tsFSBytes(buffer)[:n]);(*tsTypedArray)(value.ref).nodeBuffer=false;promise.settle(tsResult{tsIteratorResult(value,false),false})});return promise
            },
            close:func(loop *tsLoop)*tsPromise{return release(loop)},
        }
        value:=tsWebStreamValue(iterator);state:=(*tsObject)(value.ref).webStream;webClose=func(loop *tsLoop)*tsPromise{return state.cancelStream(loop)};return value
    })))
    out.set("readLines",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{options:=tsNewObject();if !tsIsUndefined(tsArg(args,0)){tsObjectSpread(options,tsArg(args,0))};if tsIsUndefined(options.get("encoding")){options.set("encoding",tsStringReference(tsStringUTF8("utf8")))};input:=tsCall(loop,out.get("createReadStream"),tsObjectValue(options));return tsFSReadLines(loop,input)})))

    // FileHandle pull APIs share the iterator pipeline runtime. Reading stays
    // incremental: each next() performs one positional or sequential fs read.
    out.set("pull",tsNamedNativeMethod("pull",0,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{
        utils:=tsInternalStreamIterUtilsModule(loop)
        parsed:=tsCall(loop,tsGet(loop,utils,tsStringReference(tsStringUTF8("parsePullArgs"))),tsArrayValue(&tsArray{values:args}))
        transforms:=(*tsArray)(tsGet(loop,parsed,tsStringReference(tsStringUTF8("transforms"))).ref).values
        options:=tsGet(loop,parsed,tsStringReference(tsStringUTF8("options")))
        autoClose:=tsTruthy(tsFSOption(options,"autoClose",tsBooleanValue(false)))
        readSize:=int(tsFSInteger(tsFSProperty(options,"chunkSize"),"options.chunkSize",64*1024,1))
        pos:=tsFSInteger(tsFSProperty(options,"start"),"options.start",-1,0)
        remaining:=tsFSInteger(tsFSProperty(options,"limit"),"options.limit",-1,1)
        signal:=tsFSProperty(options,"signal")
        if !tsIsUndefined(signal)&&!tsNullish(signal){validate:=tsGet(loop,tsInternalValidatorsModule(loop),tsStringReference(tsStringUTF8("validateAbortSignal")));tsCall(loop,validate,signal,tsStringReference(tsStringUTF8("options.signal")));if tsTruthy(tsFSProperty(signal,"aborted")){source:=tsNodeIterFromIteratorObject([]tsValue{tsFSProperty(signal,"reason")},true);return tsCall(loop,tsGet(loop,tsInternalStreamIterPullModule(loop),tsStringReference(tsStringUTF8("pull"))),source,tsArrayValue(&tsArray{values:[]tsValue{}}))}}
        if tsNumber(out.get("fd"))<0||closing!=nil{panic(tsThrown{tsFSError(syscall.EBADF,"read","","")})}
        locked:=false;done:=false
        cleanup:=func(loop *tsLoop)*tsPromise{if done{return loop.resolved(tsU,false)};done=true;if locked{locked=false;pending--};if autoClose{value:=tsCall(loop,out.get("close"));if value.kind==tsPromiseKind{return (*tsPromise)(value.ref)};return loop.resolved(value,false)};if closing!=nil&&pending==0{closeNow(loop)};return loop.resolved(tsU,false)}
        iterator:=&tsAsyncIterator{}
        iterator.next=func(loop *tsLoop)*tsPromise{promise:=loop.promise();if done||remaining==0{closingPromise:=cleanup(loop);loop.await(tsPromiseValue(closingPromise),func(result tsResult){if result.rejected{promise.settle(result)}else{promise.settle(tsResult{tsIteratorResult(tsU,true),false})}});return promise};if !locked{locked=true;pending++};size:=readSize;if remaining>0&&remaining<int64(size){size=int(remaining)};buffer:=tsNodeBuffer(make([]byte,size));position:=tsNull;if pos>=0{position=tsNumberValue(float64(pos))};read:=tsCall(loop,out.get("read"),buffer,tsNumberValue(0),tsNumberValue(float64(size)),position);loop.await(read,func(result tsResult){if result.rejected{cleanup(loop);promise.settle(result);return};count:=int(tsNumber(tsFSProperty(result.value,"bytesRead")));if count<=0{closingPromise:=cleanup(loop);loop.await(tsPromiseValue(closingPromise),func(closed tsResult){if closed.rejected{promise.settle(closed)}else{promise.settle(tsResult{tsIteratorResult(tsU,true),false})}});return};if pos>=0{pos+=int64(count)};if remaining>0{remaining-=int64(count)};chunk:=tsNodeBuffer(tsFSBytes(buffer)[:count]);promise.settle(tsResult{tsIteratorResult(tsArrayValue(&tsArray{values:[]tsValue{chunk}}),false),false})});return promise}
        iterator.close=func(loop *tsLoop)*tsPromise{return cleanup(loop)}
        iterator.resume=func(loop *tsLoop,action string,value tsValue)*tsPromise{if action=="return"{return iterator.close(loop)};return iterator.next(loop)}
        source:=tsAsyncIteratorObject(iterator);pullArgs:=append([]tsValue{},transforms...);if !tsNullish(options)&&!tsIsUndefined(options){pullArgs=append(pullArgs,options)}
        values:=append([]tsValue{source},pullArgs...);return tsCall(loop,tsGet(loop,tsInternalStreamIterPullModule(loop),tsStringReference(tsStringUTF8("pull"))),values...)
    }))

    out.set("pullSync",tsNamedNativeMethod("pullSync",0,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{
        utils:=tsInternalStreamIterUtilsModule(loop);parsed:=tsCall(loop,tsGet(loop,utils,tsStringReference(tsStringUTF8("parsePullArgs"))),tsArrayValue(&tsArray{values:args}));transforms:=(*tsArray)(tsGet(loop,parsed,tsStringReference(tsStringUTF8("transforms"))).ref).values;options:=tsGet(loop,parsed,tsStringReference(tsStringUTF8("options")))
        autoClose:=tsTruthy(tsFSOption(options,"autoClose",tsBooleanValue(false)));readSize:=int(tsFSInteger(tsFSProperty(options,"chunkSize"),"options.chunkSize",64*1024,1));pos:=tsFSInteger(tsFSProperty(options,"start"),"options.start",-1,0);remaining:=tsFSInteger(tsFSProperty(options,"limit"),"options.limit",-1,1)
        if tsNumber(out.get("fd"))<0||closing!=nil{panic(tsThrown{tsFSError(syscall.EBADF,"read","","")})};done:=false;locked:=false
        cleanup:=func(){if done{return};done=true;if locked{locked=false;pending--};if autoClose{tsFSInvoke(loop,"close","sync",[]tsValue{descriptor});out.set("fd",tsNumberValue(-1));events.emitEvent("close")}else if closing!=nil&&pending==0{closeNow(loop)}}
        iterator:=tsNewObject();iterator.set("next",tsNamedNativeMethod("next",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{if done||remaining==0{cleanup();return tsIteratorResult(tsU,true)};if !locked{locked=true;pending++};size:=readSize;if remaining>0&&remaining<int64(size){size=int(remaining)};buffer:=tsNodeBuffer(make([]byte,size));position:=tsNull;if pos>=0{position=tsNumberValue(float64(pos))};result:=tsFSInvoke(loop,"read","sync",[]tsValue{descriptor,buffer,tsNumberValue(0),tsNumberValue(float64(size)),position});count:=int(tsNumber(result));if count<=0{cleanup();return tsIteratorResult(tsU,true)};if pos>=0{pos+=int64(count)};if remaining>0{remaining-=int64(count)};chunk:=tsNodeBuffer(tsFSBytes(buffer)[:count]);return tsIteratorResult(tsArrayValue(&tsArray{values:[]tsValue{chunk}}),false)}))
        iterator.set("return",tsNamedNativeMethod("return",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{cleanup();return tsIteratorResult(tsU,true)}));iterator.set(tsPropertyKey(tsWellKnownSymbol("iterator")),tsNamedNativeMethod("[Symbol.iterator]",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{return receiver}));source:=tsObjectValue(iterator)
        if len(transforms)==0{return source};values:=append([]tsValue{source},transforms...);return tsCall(loop,tsGet(loop,tsInternalStreamIterPullModule(loop),tsStringReference(tsStringUTF8("pullSync"))),values...)
    }))

    out.set("writer",tsNamedNativeMethod("writer",0,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{
        options:=tsArg(args,0);if tsNullish(options)||tsIsUndefined(options){options=tsObjectValue(tsNewObject())};if tsNumber(out.get("fd"))<0||closing!=nil{panic(tsThrown{tsFSError(syscall.EBADF,"write","","")})};autoClose:=tsTruthy(tsFSOption(options,"autoClose",tsBooleanValue(false)));threshold:=int(tsFSInteger(tsFSProperty(options,"chunkSize"),"options.chunkSize",64*1024,1));pos:=tsFSInteger(tsFSProperty(options,"start"),"options.start",-1,0);remaining:=tsFSInteger(tsFSProperty(options,"limit"),"options.limit",-1,1);total:=0;closed:=false;var failure tsValue
        writer:=tsNewObject();normalize:=func(chunk tsValue)tsValue{return tsInternalStreamIterUtilsToUint8Array(loop,chunk)}
        writer.set("write",tsNamedNativeMethod("write",1,func(loop *tsLoop,_ tsValue,writeArgs ...tsValue)tsValue{if !tsIsUndefined(failure){return tsPromiseValue(loop.resolved(failure,true))};if closed{return tsPromiseValue(loop.resolved(tsFSException("TypeError","ERR_INVALID_STATE","The writer is closed"),true))};chunk:=normalize(tsArg(writeArgs,0));length:=len(tsFSBytes(chunk));if remaining>=0&&int64(length)>remaining{return tsPromiseValue(loop.resolved(tsFSException("RangeError","ERR_OUT_OF_RANGE","write exceeds options.limit"),true))};position:=tsNull;if pos>=0{position=tsNumberValue(float64(pos));pos+=int64(length)};if remaining>0{remaining-=int64(length)};value:=tsCall(loop,out.get("write"),chunk,tsNumberValue(0),tsNumberValue(float64(length)),position);promise:=loop.promise();loop.await(value,func(result tsResult){if result.rejected{failure=result.value;promise.settle(result);return};total+=int(tsNumber(tsFSProperty(result.value,"bytesWritten")));promise.settle(tsResult{tsU,false})});return tsPromiseValue(promise)}))
        writer.set("writev",tsNamedNativeMethod("writev",1,func(loop *tsLoop,_ tsValue,writeArgs ...tsValue)tsValue{chunks:=tsInternalStreamIterUtilsConvertChunks(loop,tsArg(writeArgs,0));size:=0;for _,chunk:=range (*tsArray)(chunks.ref).values{size+=len(tsFSBytes(chunk))};if remaining>=0&&int64(size)>remaining{return tsPromiseValue(loop.resolved(tsFSException("RangeError","ERR_OUT_OF_RANGE","writev exceeds options.limit"),true))};position:=tsNull;if pos>=0{position=tsNumberValue(float64(pos));pos+=int64(size)};if remaining>0{remaining-=int64(size)};value:=tsCall(loop,out.get("writev"),chunks,position);promise:=loop.promise();loop.await(value,func(result tsResult){if result.rejected{failure=result.value;promise.settle(result);return};total+=int(tsNumber(tsFSProperty(result.value,"bytesWritten")));promise.settle(tsResult{tsU,false})});return tsPromiseValue(promise)}))
        writer.set("writeSync",tsNamedNativeMethod("writeSync",1,func(loop *tsLoop,_ tsValue,writeArgs ...tsValue)tsValue{if closed||!tsIsUndefined(failure){return tsBooleanValue(false)};chunk:=normalize(tsArg(writeArgs,0));data:=tsFSBytes(chunk);if len(data)>threshold||remaining>=0&&int64(len(data))>remaining{return tsBooleanValue(false)};position:=tsNull;if pos>=0{position=tsNumberValue(float64(pos))};written:=tsFSInvoke(loop,"write","sync",[]tsValue{descriptor,chunk,tsNumberValue(0),tsNumberValue(float64(len(data))),position});count:=int(tsNumber(written));total+=count;if pos>=0{pos+=int64(count)};if remaining>0{remaining-=int64(count)};return tsBooleanValue(count==len(data))}))
        writer.set("writevSync",tsNamedNativeMethod("writevSync",1,func(loop *tsLoop,_ tsValue,writeArgs ...tsValue)tsValue{if closed||!tsIsUndefined(failure){return tsBooleanValue(false)};chunks:=tsInternalStreamIterUtilsConvertChunks(loop,tsArg(writeArgs,0));size:=0;for _,chunk:=range (*tsArray)(chunks.ref).values{size+=len(tsFSBytes(chunk))};if size>threshold||remaining>=0&&int64(size)>remaining{return tsBooleanValue(false)};position:=tsNull;if pos>=0{position=tsNumberValue(float64(pos))};result:=tsFSInvoke(loop,"writev","sync",[]tsValue{descriptor,chunks,position});count:=int(tsNumber(result));total+=count;if pos>=0{pos+=int64(count)};if remaining>0{remaining-=int64(count)};return tsBooleanValue(count==size)}))
        writer.set("endSync",tsNamedNativeMethod("endSync",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{if closed{return tsNumberValue(float64(total))};closed=true;if autoClose{tsFSInvoke(loop,"close","sync",[]tsValue{descriptor});out.set("fd",tsNumberValue(-1));events.emitEvent("close")};return tsNumberValue(float64(total))}))
        writer.set("end",tsNamedNativeMethod("end",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{if closed{return tsPromiseValue(loop.resolved(tsNumberValue(float64(total)),false))};closed=true;if autoClose{value:=tsCall(loop,out.get("close"));promise:=loop.promise();loop.await(value,func(result tsResult){if result.rejected{promise.settle(result)}else{promise.settle(tsResult{tsNumberValue(float64(total)),false})}});return tsPromiseValue(promise)};return tsPromiseValue(loop.resolved(tsNumberValue(float64(total)),false))}))
        writer.set("fail",tsNamedNativeMethod("fail",1,func(loop *tsLoop,_ tsValue,failArgs ...tsValue)tsValue{failure=tsArg(failArgs,0);closed=true;if autoClose&&!tsIsUndefined(out.get("fd"))&&tsNumber(out.get("fd"))>=0{tsCall(loop,out.get("close"))};return tsU}))
        return tsObjectValue(writer)
    }))
    return tsObjectValue(out)
}
`
