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
    return tsObjectValue(out)
}
`
