package goemit

const fsDirectoryRuntime = `
func tsFSDirent(entry os.DirEntry, parent string, encoding string) tsValue {
    out:=tsNewObject();out.set("name",tsFSDecode([]byte(entry.Name()),encoding))
    out.set("parentPath",tsStringReference(tsStringUTF8(parent)));out.set("path",tsStringReference(tsStringUTF8(parent)))
    mode:=entry.Type();if mode==0{if info,err:=entry.Info();err==nil{mode=info.Mode()}}
    tsFSSetTypePredicates(out,mode);return tsObjectValue(out)
}

func tsFSReadDirectory(path string, options tsValue) (tsValue,error) {
    out:=&tsArray{};encoding:=tsFSEncoding(options,"utf8");withTypes:=tsTruthy(tsFSProperty(options,"withFileTypes"));recursive:=tsTruthy(tsFSProperty(options,"recursive"))
    var walk func(string,string)error
    walk=func(directory,prefix string)error{
        entries,err:=os.ReadDir(directory);if err!=nil{return err}
        for _,entry:=range entries{
            name:=filepath.Join(prefix,entry.Name())
            if withTypes{out.values=append(out.values,tsFSDirent(entry,directory,encoding))}else{out.values=append(out.values,tsFSDecode([]byte(name),encoding))}
            if recursive&&entry.IsDir(){if err:=walk(filepath.Join(directory,entry.Name()),name);err!=nil{return err}}
        };return nil
    }
    err:=walk(path,"");return tsArrayValue(out),err
}

type tsFSDirectory struct {
    sync.Mutex
    file *os.File
    path,encoding string
    closed,recursive bool
    pending []string
}

func (directory *tsFSDirectory) readEntry(loop *tsLoop) tsResult {
    directory.Lock();defer directory.Unlock()
    if directory.closed{return tsResult{tsFSException("Error","ERR_DIR_CLOSED","Directory handle was closed"),true}}
    for{
        if directory.file==nil{return tsResult{tsNull,false}}
        entries,err:=directory.file.ReadDir(1)
        if err!=nil&&!errors.Is(err,io.EOF){return tsResult{tsFSError(err,"readdir",directory.path,""),true}}
        if len(entries)>0{entry:=entries[0];if directory.recursive&&entry.IsDir(){directory.pending=append(directory.pending,filepath.Join(directory.path,entry.Name()))};return tsResult{tsFSDirent(entry,directory.path,directory.encoding),false}}
        if !directory.recursive{return tsResult{tsNull,false}}
        _ = directory.file.Close();directory.file=nil
        if len(directory.pending)==0{return tsResult{tsNull,false}}
        directory.path=directory.pending[0];directory.pending=directory.pending[1:];directory.file,err=os.Open(directory.path);if err!=nil{return tsResult{tsFSError(err,"opendir",directory.path,""),true}}
    }

}

func (directory *tsFSDirectory) close() tsResult {
    directory.Lock();defer directory.Unlock()
    if directory.closed{return tsResult{tsFSException("Error","ERR_DIR_CLOSED","Directory handle was closed"),true}}
    directory.closed=true
    if directory.file==nil{return tsResult{tsU,false}}
    if err:=directory.file.Close();err!=nil{return tsResult{tsFSError(err,"closedir",directory.path,""),true}}
    return tsResult{tsU,false}
}

func tsFSDir(file *os.File, path string, options tsValue) tsValue {
    directory:=&tsFSDirectory{file:file,path:path,encoding:tsFSEncoding(options,"utf8"),recursive:tsTruthy(tsFSProperty(options,"recursive"))}
    out:=tsNewObject();out.nativeClass=tsNativeClass("Dir");out.set("path",tsStringReference(tsStringUTF8(path)))
    for _,name:=range []string{"read","close"}{
        operation:=name
        var work func()tsResult=func()tsResult{return directory.readEntry(loop)};if name=="close"{work=directory.close}
        action:=work
        out.set(name+"Sync",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{result:=action();if result.rejected{panic(tsThrown{result.value})};return result.value})))
        out.set(name,tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{
            callback:=tsArg(args,0)
            if tsIsUndefined(callback){return tsPromiseValue(loop.start(action))}
            if callback.kind!=tsFunctionKind{tsFSInvalid("callback","a function")}
            loop.submit(action,func(result tsResult){if result.rejected{tsCall(loop,callback,result.value)}else if operation=="close"{tsCall(loop,callback,tsNull)}else{tsCall(loop,callback,tsNull,result.value)}});return tsU
        })))
    }
    out.asyncIterator=func()*tsAsyncIterator{
        finished:=false
        return &tsAsyncIterator{
            next:func(loop *tsLoop)*tsPromise{
                if finished{return loop.resolved(tsIteratorResult(tsU,true),false)}
                promise:=loop.promise()
                loop.submit(func()tsResult{return directory.readEntry(loop)},func(result tsResult){
                    if result.rejected{finished=true;directory.close();promise.settle(result);return}
                    if result.value.kind==tsNullKind{finished=true;closed:=directory.close();if closed.rejected{promise.settle(closed);return};promise.settle(tsResult{tsIteratorResult(tsU,true),false})}else{promise.settle(tsResult{tsIteratorResult(result.value,false),false})}
                });return promise
            },
            close:func(loop *tsLoop)*tsPromise{if finished{return loop.resolved(tsIteratorResult(tsU,true),false)};finished=true;promise:=loop.promise();loop.submit(directory.close,func(result tsResult){if !result.rejected{result.value=tsIteratorResult(tsU,true)};promise.settle(result)});return promise},
        }
    }
    return tsObjectValue(out)
}
`
