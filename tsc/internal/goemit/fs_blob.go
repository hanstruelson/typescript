package goemit

const fsBlobRuntime = `
type tsBlobFile struct {sync.Mutex;file *os.File;stat unix.Stat_t}
type tsBlob struct {data []byte;file *tsBlobFile;offset,length int64;mediaType string}
func tsBlobError()tsValue{return tsErrorValue(&tsRuntimeError{name:"NotReadableError",message:"The blob could not be read",fields:map[string]tsValue{"code":tsNumberValue(0)}})}
func(file *tsBlobFile)validate()bool{var stat unix.Stat_t;if unix.Fstat(int(file.file.Fd()),&stat)!=nil{return false};return stat.Size==file.stat.Size&&stat.Mtim==file.stat.Mtim&&stat.Ctim==file.stat.Ctim}
func(blob *tsBlob)readRange(offset,length int64)(result tsResult){
    defer func(){if failure:=recover();failure!=nil{result=tsResult{tsUnwrap(failure),true}}}()
    if blob.file==nil{data:=append([]byte{},blob.data[blob.offset+offset:blob.offset+offset+length]...);return tsResult{tsNodeBuffer(data),false}}
    file:=blob.file;file.Lock();defer file.Unlock();if !file.validate(){return tsResult{tsBlobError(),true}}
    data:=make([]byte,int(length));done:=0;for done<len(data){if !file.validate(){return tsResult{tsBlobError(),true}};size:=len(data)-done;if size>512*1024{size=512*1024};n,err:=file.file.ReadAt(data[done:done+size],blob.offset+offset+int64(done));done+=n;if err!=nil{if !errors.Is(err,io.EOF)||done<len(data){return tsResult{tsBlobError(),true}}};if n==0&&done<len(data){return tsResult{tsBlobError(),true}}};if !file.validate(){return tsResult{tsBlobError(),true}};return tsResult{tsNodeBuffer(data),false}
}
func tsBlobMediaType(value tsValue)string{if tsIsUndefined(value){return ""};text:=tsText(value);for _,r:=range text{if r<0x20||r>0x7e{return ""}};return strings.ToLower(text)}
func tsBlobObject(blob *tsBlob)tsValue{
    out:=tsNewObject();out.nativeClass=tsNativeClass("Blob");out.blob=blob;out.set("size",tsNumberValue(float64(blob.length)));out.set("type",tsStringReference(tsStringUTF8(blob.mediaType)))
    for _,method:=range []string{"text","arrayBuffer","bytes"}{name:=method;out.set(name,tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsPromiseValue(loop.start(func()tsResult{result:=blob.readRange(0,blob.length);if result.rejected{return result};data:=tsFSBytes(result.value);switch name{case "text":if bytes.HasPrefix(data,[]byte{0xef,0xbb,0xbf}){data=data[3:]};result.value=tsFSDecode(data,"utf8");case "arrayBuffer":result.value=tsArrayBufferValue(&tsArrayBuffer{data:data});case "bytes":(*tsTypedArray)(result.value.ref).nodeBuffer=false};return result}))})))}
    out.set("slice",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{index:=func(value tsValue,fallback int64)int64{if tsIsUndefined(value){return fallback};n:=tsNumber(value);if math.IsNaN(n){return 0};if n<0{n+=float64(blob.length)};if n<0{return 0};if n>float64(blob.length){return blob.length};return int64(math.Trunc(n))};first,last:=index(tsArg(args,0),0),index(tsArg(args,1),blob.length);if last<first{last=first};return tsBlobObject(&tsBlob{data:blob.data,file:blob.file,offset:blob.offset+first,length:last-first,mediaType:tsBlobMediaType(tsArg(args,2))})})))
    out.set("stream",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{
        position:=int64(0);closed:=false
        iterator:=&tsAsyncIterator{
            next:func(loop *tsLoop)*tsPromise{if closed||position>=blob.length{closed=true;return loop.resolved(tsIteratorResult(tsU,true),false)};offset:=position;size:=blob.length-position;if size>64*1024{size=64*1024};position+=size;return loop.start(func()tsResult{result:=blob.readRange(offset,size);if !result.rejected{(*tsTypedArray)(result.value.ref).nodeBuffer=false;result.value=tsIteratorResult(result.value,false)};return result})},
            close:func(loop *tsLoop)*tsPromise{closed=true;return loop.resolved(tsU,false)},
        };return tsWebStreamValue(iterator)
    })))
    return tsObjectValue(out)
}
func tsFSOpenAsBlob(loop *tsLoop,args []tsValue)tsValue{
    path:=tsFSPath(tsArg(args,0));mediaType:=tsBlobMediaType(tsFSProperty(tsArg(args,1),"type"))
    return tsPromiseValue(loop.start(func()tsResult{
        file,err:=os.Open(path);if err!=nil{return tsResult{tsFSError(err,"open",path,""),true}};state:=&tsBlobFile{file:file};if err:=unix.Fstat(int(file.Fd()),&state.stat);err!=nil{file.Close();return tsResult{tsFSError(err,"fstat",path,""),true}};if state.stat.Mode&syscall.S_IFMT!=syscall.S_IFREG{file.Close();return tsResult{tsFSError(syscall.EISDIR,"open",path,""),true}}
        runtime.SetFinalizer(state,func(file *tsBlobFile){file.file.Close()});return tsResult{tsBlobObject(&tsBlob{file:state,length:state.stat.Size,mediaType:mediaType}),false}
    }))
}
func tsBlobConstruct(loop *tsLoop,args []tsValue)tsValue{
    data:=[]byte{};parts:=tsArg(args,0);options:=tsArg(args,1);native:=tsText(tsFSProperty(options,"endings"))=="native"
    if !tsIsUndefined(parts){iterator:=tsIterate(parts);for iterator.next(loop){value:=iterator.value;switch value.kind{case tsTypedArrayKind,tsArrayBufferKind:data=append(data,tsFSBytes(value)...);case tsObjectKind:object:=(*tsObject)(value.ref);if object.blob!=nil{result:=object.blob.readRange(0,object.blob.length);if result.rejected{panic(tsThrown{result.value})};data=append(data,tsFSBytes(result.value)...)}else{text:=tsStringValue(value);data=append(data,tsFSEncode(text,"utf8")...)};default:text:=tsStringValue(value);encoded:=tsFSEncode(text,"utf8");if native{encoded=[]byte(strings.ReplaceAll(strings.ReplaceAll(string(encoded),"\r\n","\n"),"\r","\n"))};data=append(data,encoded...)}}}
    return tsBlobObject(&tsBlob{data:data,length:int64(len(data)),mediaType:tsBlobMediaType(tsFSProperty(options,"type"))})
}
`
