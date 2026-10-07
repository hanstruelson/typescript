package goemit

// The OS adapter owns descriptors and cancellation. Listener delivery stays on
// the calling event loop; no JavaScript executes on an inotify worker.
const fsWatchRuntime = `
type tsFSListener struct { callback tsValue; once bool }
type tsFSEvents struct { loop *tsLoop; object *tsObject; listeners map[string][]tsFSListener }
func tsFSEventObject(loop *tsLoop)*tsFSEvents{
    events:=&tsFSEvents{loop:loop,object:tsNewObject(),listeners:map[string][]tsFSListener{}}
    for _,name:=range []string{"on","addListener","once","prependListener","prependOnceListener"}{
        operation:=name
        events.object.set(name,tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{
            event:=tsText(tsArg(args,0));callback:=tsArg(args,1);if callback.kind!=tsFunctionKind{tsFSInvalid("listener","a function")}
            listener:=tsFSListener{callback:callback,once:operation=="once"||operation=="prependOnceListener"}
            if strings.HasPrefix(operation,"prepend"){events.listeners[event]=append([]tsFSListener{listener},events.listeners[event]...)}else{events.listeners[event]=append(events.listeners[event],listener)}
            return tsObjectValue(events.object)
        })))
    }
    remove:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{event:=tsText(tsArg(args,0));callback:=tsArg(args,1);items:=events.listeners[event];for i:=len(items)-1;i>=0;i--{if tsStrictEqual(items[i].callback,callback){events.listeners[event]=append(items[:i],items[i+1:]...);break}};return tsObjectValue(events.object)}))
    events.object.set("removeListener",remove);events.object.set("off",remove)
    events.object.set("removeAllListeners",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{if len(args)==0{events.listeners=map[string][]tsFSListener{}}else{delete(events.listeners,tsText(args[0]))};return tsObjectValue(events.object)})))
    events.object.set("listenerCount",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNumberValue(float64(len(events.listeners[tsText(tsArg(args,0))])))})))
    events.object.set("listeners",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{out:=&tsArray{};for _,listener:=range events.listeners[tsText(tsArg(args,0))]{out.values=append(out.values,listener.callback)};return tsArrayValue(out)})))
    events.object.set("emit",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{event:=tsText(tsArg(args,0));return tsBooleanValue(events.emitEvent(event,args[1:]...))})))
    return events
}
func(events *tsFSEvents)emitEvent(event string,args ...tsValue)bool{
    loop:=events.loop
    snapshot:=append([]tsFSListener{},events.listeners[event]...)
    if len(snapshot)==0{if event=="error"&&len(args)>0{panic(tsThrown{args[0]})};return false}
    for _,listener:=range snapshot{
        if listener.once{items:=events.listeners[event];for i,item:=range items{if tsStrictEqual(item.callback,listener.callback){events.listeners[event]=append(items[:i],items[i+1:]...);break}}}
        tsCall(loop,listener.callback,args...)
    };return true
}

type tsFSWatchState struct { once sync.Once; cancel chan struct{}; events *tsFSEvents; task *tsTask; closed bool }
func(state *tsFSWatchState)cancelNative(){state.once.Do(func(){close(state.cancel)})}
func(state *tsFSWatchState)close(){if state.closed{return};state.closed=true;state.cancelNative()}
func tsFSWatch(loop *tsLoop,args []tsValue)tsValue{
    path:=tsFSPath(tsArg(args,0));options:=tsArg(args,1);callback:=tsArg(args,2)
    if options.kind==tsFunctionKind{callback=options;options=tsU}
    if !tsIsUndefined(callback)&&callback.kind!=tsFunctionKind{tsFSInvalid("listener","a function")}
    signal:=tsFSAbortSignal(options);var aborted <-chan struct{};if signal!=nil{aborted=signal.done}
    if signal!=nil{_,already:=tsAbortReason(signal);if already{events:=tsFSEventObject(loop);events.object.set("close",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsU})));for _,name:=range []string{"ref","unref"}{events.object.set(name,tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsObjectValue(events.object)})))};loop.ready=append(loop.ready,func(){events.emitEvent("close")});return tsObjectValue(events.object)}}
    encoding:=tsFSEncoding(options,"utf8");recursive:=tsTruthy(tsFSProperty(options,"recursive"))
    fd,err:=unix.InotifyInit1(unix.IN_NONBLOCK|unix.IN_CLOEXEC);if err!=nil{panic(tsThrown{tsFSError(err,"watch",path,"")})}
    directories:=map[int]string{}
    mask:=uint32(unix.IN_MODIFY|unix.IN_ATTRIB|unix.IN_CREATE|unix.IN_DELETE|unix.IN_DELETE_SELF|unix.IN_MOVE_SELF|unix.IN_MOVED_FROM|unix.IN_MOVED_TO)
    add:=func(directory,relative string)error{wd,err:=unix.InotifyAddWatch(fd,directory,mask);if err==nil{directories[wd]=relative};return err}
    err=add(path,"")
    if err==nil&&recursive{err=filepath.WalkDir(path,func(p string,d os.DirEntry,e error)error{if e!=nil{return e};if d.IsDir()&&p!=path{relative,e:=filepath.Rel(path,p);if e!=nil{return e};return add(p,relative)};return nil})}
    if err!=nil{unix.Close(fd);panic(tsThrown{tsFSError(err,"watch",path,"")})}
    state:=&tsFSWatchState{cancel:make(chan struct{}),events:tsFSEventObject(loop)};state.events.object.nativeClass=tsNativeClass("FSWatcher")
    if callback.kind==tsFunctionKind{state.events.listeners["change"]=append(state.events.listeners["change"],tsFSListener{callback:callback})}
    isDirectory:=false;if info,e:=os.Stat(path);e==nil{isDirectory=info.IsDir()}
    state.task=loop.submit(func()tsResult{
        defer unix.Close(fd);buffer:=make([]byte,64*1024)
        for{
            select{case <-state.cancel:return tsResult{tsU,false};case <-aborted:return tsResult{tsU,false};default:}
            descriptors:=[]unix.PollFd{{Fd:int32(fd),Events:unix.POLLIN}};_,e:=unix.Poll(descriptors,50);if e==syscall.EINTR{continue};if e!=nil{return tsResult{tsFSError(e,"watch",path,""),true}}
            if descriptors[0].Revents==0{continue}
            n,e:=unix.Read(fd,buffer);if e==syscall.EAGAIN||e==syscall.EINTR{continue};if e!=nil{return tsResult{tsFSError(e,"watch",path,""),true}}
            for offset:=0;offset+16<=n;{
                wd:=int(int32(binary.NativeEndian.Uint32(buffer[offset:])));flags:=binary.NativeEndian.Uint32(buffer[offset+4:]);length:=int(binary.NativeEndian.Uint32(buffer[offset+12:]));if offset+16+length>n{break}
                raw:=buffer[offset+16:offset+16+length];if zero:=bytes.IndexByte(raw,0);zero>=0{raw=raw[:zero]};filename:=string(raw);relative:=directories[wd];offset+=16+length
                if flags&unix.IN_Q_OVERFLOW!=0{return tsResult{tsFSException("Error","ENOSPC","Filesystem watcher queue overflow"),true}}
                if flags&unix.IN_IGNORED!=0{delete(directories,wd);continue}
                if recursive&&flags&unix.IN_ISDIR!=0&&flags&(unix.IN_CREATE|unix.IN_MOVED_TO)!=0{full:=filepath.Join(path,relative,filename);_ = filepath.WalkDir(full,func(p string,d os.DirEntry,e error)error{if e==nil&&d.IsDir(){rel,_:=filepath.Rel(path,p);return add(p,rel)};return nil})}
                if relative!=""{filename=filepath.Join(relative,filename)};if filename==""&&!isDirectory{filename=filepath.Base(path)}
                event:="change";if flags&(unix.IN_CREATE|unix.IN_DELETE|unix.IN_DELETE_SELF|unix.IN_MOVE_SELF|unix.IN_MOVED_FROM|unix.IN_MOVED_TO)!=0{event="rename"}
                eventValue:=tsStringReference(tsStringUTF8(event));fileValue:=tsNull;if filename!=""{fileValue=tsFSDecode([]byte(filename),encoding)}
                loop.post(func(){if !state.closed{state.events.emitEvent("change",eventValue,fileValue)}})
            }
        }
    },func(result tsResult){state.closed=true;if result.rejected{state.events.emitEvent("error",result.value)};state.events.emitEvent("close")})
    state.task.unref=!tsTruthy(tsFSOption(options,"persistent",tsBooleanValue(true)))
    loop.resourceClosers=append(loop.resourceClosers,state.cancelNative)
    state.events.object.set("close",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{state.close();return tsU})))
    for _,name:=range []string{"ref","unref"}{operation:=name;state.events.object.set(name,tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{state.task.unref=operation=="unref";return tsObjectValue(state.events.object)})))}
    return tsObjectValue(state.events.object)
}
func tsFSWatchIterator(loop *tsLoop,args []tsValue)tsValue{
    signal:=tsFSAbortSignal(tsArg(args,1));if signal!=nil{_,aborted:=tsAbortReason(signal);if aborted{finished:=false;return tsAsyncIteratorObject(&tsAsyncIterator{next:func(loop *tsLoop)*tsPromise{if finished{return loop.resolved(tsIteratorResult(tsU,true),false)};finished=true;return loop.resolved(tsFSAbortError(signal),true)},close:func(loop *tsLoop)*tsPromise{finished=true;return loop.resolved(tsIteratorResult(tsU,true),false)}})}}
    watcher:=tsFSWatch(loop,args);queue:=[]tsValue{};waiters:=[]*tsPromise{};closed:=false;failure:=tsU
    deliver:=func(value tsValue){if len(waiters)>0{promise:=waiters[0];waiters=waiters[1:];promise.settle(tsResult{tsIteratorResult(value,false),false})}else{queue=append(queue,value)}}
    tsCall(loop,tsFSProperty(watcher,"on"),tsStringReference(tsStringUTF8("change")),tsFunctionValue(tsFunc(func(loop *tsLoop,values ...tsValue)tsValue{out:=tsNewObject();out.set("eventType",tsArg(values,0));out.set("filename",tsArg(values,1));deliver(tsObjectValue(out));return tsU})))
    tsCall(loop,tsFSProperty(watcher,"on"),tsStringReference(tsStringUTF8("error")),tsFunctionValue(tsFunc(func(loop *tsLoop,values ...tsValue)tsValue{failure=tsArg(values,0);closed=true;for _,promise:=range waiters{promise.settle(tsResult{failure,true})};waiters=nil;return tsU})))
    tsCall(loop,tsFSProperty(watcher,"on"),tsStringReference(tsStringUTF8("close")),tsFunctionValue(tsFunc(func(loop *tsLoop,values ...tsValue)tsValue{closed=true;if signal!=nil{_,aborted:=tsAbortReason(signal);if aborted{failure=tsFSAbortError(signal)}};for _,promise:=range waiters{if tsIsUndefined(failure){promise.settle(tsResult{tsIteratorResult(tsU,true),false})}else{promise.settle(tsResult{failure,true})}};waiters=nil;return tsU})))
    return tsAsyncIteratorObject(&tsAsyncIterator{
        next:func(loop *tsLoop)*tsPromise{if !tsIsUndefined(failure){return loop.resolved(failure,true)};if len(queue)>0{value:=queue[0];queue=queue[1:];return loop.resolved(tsIteratorResult(value,false),false)};if closed{return loop.resolved(tsIteratorResult(tsU,true),false)};promise:=loop.promise();waiters=append(waiters,promise);return promise},
        close:func(loop *tsLoop)*tsPromise{closed=true;tsCall(loop,tsFSProperty(watcher,"close"));for _,promise:=range waiters{promise.settle(tsResult{tsIteratorResult(tsU,true),false})};waiters=nil;return loop.resolved(tsIteratorResult(tsU,true),false)},
    })
}
func tsFSInstallWatchers(fs,promises *tsObject){
    fs.set("watchFile",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsFSWatchFile(loop,args)})))
    fs.set("unwatchFile",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsFSUnwatchFile(loop,args)})))
    fs.set("watch",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsFSWatch(loop,args)})))
    promises.set("watch",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsFSWatchIterator(loop,args)})))
}
var tsFSStatWatchers struct { sync.Mutex; states map[*tsLoop]map[string]*tsFSWatchState }
func tsFSWatchFile(loop *tsLoop,args []tsValue)tsValue{
    path:=tsFSPath(tsArg(args,0));absolute,err:=filepath.Abs(path);if err!=nil{panic(tsThrown{tsFSError(err,"watchFile",path,"")})};path=absolute
    options:=tsArg(args,1);callback:=tsArg(args,2);if options.kind==tsFunctionKind{callback=options;options=tsU};if callback.kind!=tsFunctionKind{tsFSInvalid("listener","a function")}
    options=tsArg(tsFSPrepare(loop,[]tsValue{options}),0)
    interval:=tsFSInteger(tsFSProperty(options,"interval"),"interval",5007,1)
    tsFSStatWatchers.Lock();if tsFSStatWatchers.states==nil{tsFSStatWatchers.states=map[*tsLoop]map[string]*tsFSWatchState{}};if tsFSStatWatchers.states[loop]==nil{tsFSStatWatchers.states[loop]=map[string]*tsFSWatchState{}};state:=tsFSStatWatchers.states[loop][path];tsFSStatWatchers.Unlock()
    if state!=nil&&!state.closed{state.events.listeners["change"]=append(state.events.listeners["change"],tsFSListener{callback:callback});return tsObjectValue(state.events.object)}
    state=&tsFSWatchState{cancel:make(chan struct{}),events:tsFSEventObject(loop)};state.events.object.nativeClass=tsNativeClass("StatWatcher");state.events.listeners["change"]=[]tsFSListener{{callback:callback}}
    tsFSStatWatchers.Lock();tsFSStatWatchers.states[loop][path]=state;tsFSStatWatchers.Unlock()
    state.task=loop.submit(func()tsResult{
        previous,previousErr:=os.Stat(path);previousValue:=tsFSZeroStats(options);if previousErr==nil{previousValue=tsFSStats(previous,options)}
        ticker:=time.NewTicker(time.Duration(interval)*time.Millisecond);defer ticker.Stop()
        for{select{case <-state.cancel:return tsResult{tsU,false};case <-ticker.C:}
            current,currentErr:=os.Stat(path);changed:=(currentErr==nil)!=(previousErr==nil)
            if currentErr==nil&&previousErr==nil{changed=current.Size()!=previous.Size()||!current.ModTime().Equal(previous.ModTime())||current.Mode()!=previous.Mode()||!os.SameFile(current,previous)}
            if changed{value:=tsFSZeroStats(options);if currentErr==nil{value=tsFSStats(current,options)};old:=previousValue;loop.post(func(){if !state.closed{state.events.emitEvent("change",value,old)}});previousValue=value}
            previous,previousErr=current,currentErr
        }
    },func(result tsResult){state.closed=true;if result.rejected{state.events.emitEvent("error",result.value)};state.events.emitEvent("stop")})
    state.task.unref=!tsTruthy(tsFSOption(options,"persistent",tsBooleanValue(true)));loop.resourceClosers=append(loop.resourceClosers,func(){state.cancelNative();tsFSStatWatchers.Lock();delete(tsFSStatWatchers.states[loop],path);if len(tsFSStatWatchers.states[loop])==0{delete(tsFSStatWatchers.states,loop)};tsFSStatWatchers.Unlock()})
    for _,name:=range []string{"ref","unref"}{operation:=name;state.events.object.set(name,tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{state.task.unref=operation=="unref";return tsObjectValue(state.events.object)})))}
    return tsObjectValue(state.events.object)
}
func tsFSUnwatchFile(loop *tsLoop,args []tsValue)tsValue{
    path:=tsFSPath(tsArg(args,0));path,_=filepath.Abs(path);callback:=tsArg(args,1)
    tsFSStatWatchers.Lock();state:=tsFSStatWatchers.states[loop][path];tsFSStatWatchers.Unlock();if state==nil{return tsU}
    if tsIsUndefined(callback){state.events.listeners["change"]=nil}else{var items []tsFSListener;items=state.events.listeners["change"];out:=[]tsFSListener{};for _,item:=range items{if !tsStrictEqual(item.callback,callback){out=append(out,item)}};state.events.listeners["change"]=out}
    if len(state.events.listeners["change"])==0{state.close();tsFSStatWatchers.Lock();delete(tsFSStatWatchers.states[loop],path);tsFSStatWatchers.Unlock()};return tsU
}
func tsFSZeroStats(options tsValue)tsValue{
    bigint:=tsTruthy(tsFSProperty(options,"bigint"));out:=tsNewObject();class:="Stats";if bigint{class="BigIntStats"};out.nativeClass=tsNativeClass(class);for _,key:=range []string{"dev","ino","mode","nlink","uid","gid","rdev","size","blksize","blocks","atimeMs","mtimeMs","ctimeMs","birthtimeMs"}{zero:=tsNumberValue(0);if bigint{zero=tsBigIntSigned(0)};out.set(key,zero)};for _,key:=range []string{"atime","mtime","ctime","birthtime"}{tsFSTimestamp(out,key,time.Unix(0,0),bigint)}
    return tsObjectValue(out)
}
`
