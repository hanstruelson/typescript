package goemit

// The OS adapter owns descriptors and cancellation. Listener delivery stays on
// the calling event loop; no JavaScript executes on an inotify worker.
const fsWatchRuntime = `
type tsFSListener struct { callback tsValue; once bool }
type tsFSEvents struct { loop *tsLoop; object *tsObject; emitter *tsNodeEventEmitter; listeners map[string][]tsFSListener }
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
	remove:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{events.removeEvent(tsText(tsArg(args,0)),tsArg(args,1));return tsObjectValue(events.object)}))
	events.object.set("removeListener",remove);events.object.set("off",remove)
	events.object.set("removeAllListeners",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{if len(args)==0{events.listeners=map[string][]tsFSListener{}}else{events.removeAllEvent(tsText(args[0]))};return tsObjectValue(events.object)})))
	events.object.set("listenerCount",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNumberValue(float64(events.listenerCount(tsText(tsArg(args,0)))))})))
	events.object.set("listeners",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{out:=&tsArray{};for _,listener:=range events.listeners[tsText(tsArg(args,0))]{out.values=append(out.values,listener.callback)};return tsArrayValue(out)})))
	events.object.set("emit",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{event:=tsText(tsArg(args,0));return tsBooleanValue(events.emitEvent(event,args[1:]...))})))
    return events
}
func tsFSNodeEventObject(loop *tsLoop)*tsFSEvents{tsNodeEventsInit(loop);value:=tsNodeEventCreate(loop,tsU,tsU);object:=(*tsObject)(value.ref);return &tsFSEvents{loop:loop,object:object,emitter:tsNodeEventState(value),listeners:map[string][]tsFSListener{}}}
func(events *tsFSEvents)addEvent(event string,callback tsValue,once bool){if events.emitter!=nil{tsNodeEventAdd(events.loop,tsObjectValue(events.object),tsStringReference(tsStringUTF8(event)),callback,false,once);return};listener:=tsFSListener{callback:callback,once:once};events.listeners[event]=append(events.listeners[event],listener)}
func(events *tsFSEvents)removeEvent(event string,callback tsValue){if events.emitter!=nil{tsNodeEventRemove(events.loop,events.emitter,tsStringReference(tsStringUTF8(event)),callback);return};items:=events.listeners[event];for i:=len(items)-1;i>=0;i--{if tsStrictEqual(items[i].callback,callback){events.listeners[event]=append(items[:i],items[i+1:]...);break}}}
func(events *tsFSEvents)removeAllEvent(event string){if events.emitter!=nil{key:=tsNodeEventKey(events.loop,tsStringReference(tsStringUTF8(event)));delete(events.emitter.events,key);delete(events.emitter.names,key)};delete(events.listeners,event)}
func(events *tsFSEvents)listenerCount(event string)int{count:=len(events.listeners[event]);if events.emitter!=nil{count+=len(events.emitter.events[event])};return count}
func(events *tsFSEvents)emitEvent(event string,args ...tsValue)bool{
	emitted:=false
	if events.emitter!=nil{emitted=tsNodeEventEmit(events.loop,events.emitter,tsStringReference(tsStringUTF8(event)),args)}
	loop:=events.loop;snapshot:=append([]tsFSListener{},events.listeners[event]...)
	if len(snapshot)==0&&event=="error"&&len(args)>0&&events.emitter==nil{panic(tsThrown{args[0]})}
	for _,listener:=range snapshot{if listener.once{events.removeEvent(event,listener.callback)};tsCall(loop,listener.callback,args...);emitted=true}
	return emitted
}

type tsFSWatchState struct { once sync.Once; cancel chan struct{}; events *tsFSEvents; task *tsTask; closed, silentClose bool; listener, ignoreMatcher tsValue; signal *tsAbortState; handle *tsObject; lastEventKey string; lastEventAt time.Time }
func(state *tsFSWatchState)cancelNative(){state.once.Do(func(){close(state.cancel)})}
func(state *tsFSWatchState)close(){if state.closed{return};state.closed=true;state.cancelNative()}
func tsFSWatchEvent(state *tsFSWatchState,event string,args ...tsValue){
	if state==nil||state.events==nil{return}
	if event=="close"{
		state.events.loop.ready=append(state.events.loop.ready,func(){state.events.emitEvent("close",tsU)})
		return
	}
	if event=="error"{
		fromAbort:=len(args)>1&&tsTruthy(args[1])
		if fromAbort&&state.events.listenerCount("error")==0{return}
		if !fromAbort{
			err:=tsArg(args,0);code:=tsFSProperty(err,"code")
			if tsText(code)=="EACCES"{tsSet(state.events.loop,err,tsStringReference(tsStringUTF8("code")),tsStringReference(tsStringUTF8("EPERM")))}
		}
		state.events.emitEvent("error",tsArg(args,0));return
	}
	filename:=tsArg(args,0)
	if !tsNullish(filename)&&state.ignoreMatcher.kind==tsFunctionKind&&tsTruthy(tsCall(state.events.loop,state.ignoreMatcher,filename)){return}
	filenameKey:="";if !tsNullish(filename){filenameKey=tsStringValueWithLoop(state.events.loop,filename).String()};eventKey:=event+"\x00"+filenameKey;now:=time.Now();if eventKey==state.lastEventKey&&!state.lastEventAt.IsZero()&&now.Sub(state.lastEventAt)<=time.Millisecond{return};state.lastEventKey=eventKey;state.lastEventAt=now
	state.events.emitEvent("change",tsStringReference(tsStringUTF8(event)),filename)
	if state.listener.kind==tsFunctionKind{tsCallReceiver(state.events.loop,state.listener,tsObjectValue(state.events.object),tsStringReference(tsStringUTF8(event)),filename)}
}
func tsFSWatch(loop *tsLoop,args []tsValue)tsValue{
    inputPath:=tsFSPath(tsArg(args,0));options:=tsArg(args,1);callback:=tsArg(args,2)
    if options.kind==tsFunctionKind{callback=options;options=tsU}
	if options.kind==tsStringKind{normalized:=tsNewObject();normalized.set("encoding",options);options=tsObjectValue(normalized)}
	if !tsNullish(options)&&options.kind!=tsObjectKind&&options.kind!=tsInstanceKind{tsFSInvalid("options","an object or string")}
    if tsNullish(options){options=tsObjectValue(tsNewObject())}
    if callback.kind!=tsFunctionKind{callback=tsFunctionValue(tsFunc(func(*tsLoop,...tsValue)tsValue{return tsU}))}
    absolute,absoluteErr:=filepath.Abs(inputPath);if absoluteErr!=nil{panic(tsThrown{tsFSError(absoluteErr,"watch",inputPath,"")})};path:=absolute
    persistent:=tsFSBooleanOption(options,"persistent",true);recursive:=tsFSBooleanOption(options,"recursive",false);_=tsFSBooleanOption(options,"verbose",false)
	var signal *tsAbortState;if tsTruthy(tsFSProperty(options,"signal")){signal=tsFSAbortSignal(options)};var aborted <-chan struct{};if signal!=nil{aborted=signal.done}
    ignoreMatcher:=tsInternalFSWatchIgnoreMatcher(loop,tsFSProperty(options,"ignore"))
    state:=&tsFSWatchState{cancel:make(chan struct{}),events:tsFSNodeEventObject(loop),listener:callback,ignoreMatcher:ignoreMatcher,signal:signal}
    state.events.object.nativeClass=tsNativeClass("FSWatcher")
	handle:=tsNewObject();handle.nativeClass=tsNativeClass("FSEvent");state.handle=handle
    handle.set("close",tsNamedNativeMethod("close",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{state.close();return tsU}))
    handle.set("ref",tsNamedNativeMethod("ref",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{if state.task!=nil{state.task.unref=false};return tsU}))
    handle.set("unref",tsNamedNativeMethod("unref",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{if state.task!=nil{state.task.unref=true};return tsU}))
    handle.set("onchange",tsNamedNativeMethod("onchange",3,func(loop *tsLoop,_ tsValue,values ...tsValue)tsValue{
        status,eventType,filename:=tsArg(values,0),tsArg(values,1),tsArg(values,2)
		if tsIsNumeric(status)&&tsNumber(status)<0{state.silentClose=true;state.close();state.events.object.set("_handle",tsNull);errno:=syscall.Errno(-int64(tsNumber(status)));err:=tsFSError(errno,"watch",inputPath,"");tsSet(loop,err,tsStringReference(tsStringUTF8("syscall")),tsStringReference(tsStringUTF8("watch")));filenameText:=inputPath;if !tsNullish(filename){filenameText=tsStringValueWithLoop(loop,filename).String()};tsSet(loop,err,tsStringReference(tsStringUTF8("filename")),tsStringReference(tsStringUTF8(filenameText)));tsFSWatchEvent(state,"error",err,tsBooleanValue(false));return tsU}
        tsFSWatchEvent(state,tsText(eventType),filename);return tsU
    }))
    state.events.object.set("_handle",tsObjectValue(handle))
    state.events.object.set("start",tsNamedNativeMethod("start",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsU}))
	state.events.object.set("close",tsNamedNativeMethod("close",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{handle:=tsFSProperty(receiver,"_handle");if handle.kind==tsNullKind{return tsU};valid:=handle.kind==tsObjectKind&&(*tsObject)(handle.ref).nativeClass!=nil&&(*tsObject)(handle.ref).nativeClass.builtin=="FSEvent";if !valid{panic(tsThrown{tsFSException("Error","ERR_INTERNAL_ASSERTION","handle must be a FSEvent"+" (see internal assertion details)")})};method:=tsFSProperty(handle,"close");if method.kind==tsFunctionKind{tsCallReceiver(loop,method,handle)};return tsU}))
	state.events.object.set("ref",tsNamedNativeMethod("ref",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{handle:=tsFSProperty(receiver,"_handle");if tsTruthy(handle){method:=tsFSProperty(handle,"ref");if method.kind==tsFunctionKind{tsCallReceiver(loop,method,handle)}};return receiver}))
	state.events.object.set("unref",tsNamedNativeMethod("unref",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{handle:=tsFSProperty(receiver,"_handle");if tsTruthy(handle){method:=tsFSProperty(handle,"unref");if method.kind==tsFunctionKind{tsCallReceiver(loop,method,handle)}};return receiver}))
	if signal!=nil{if reason,already:=tsAbortReason(signal);already{state.closed=true;loop.ready=append(loop.ready,func(){tsFSWatchEvent(state,"close");tsFSWatchEvent(state,"error",reason,tsBooleanValue(true))});return tsObjectValue(state.events.object)}}
    encoding:=tsFSEncoding(options,"utf8")
    fd,err:=unix.InotifyInit1(unix.IN_NONBLOCK|unix.IN_CLOEXEC);if err!=nil{panic(tsThrown{tsFSError(err,"watch",inputPath,"")})}
    directories:=map[int]string{}
    mask:=uint32(unix.IN_MODIFY|unix.IN_ATTRIB|unix.IN_CREATE|unix.IN_DELETE|unix.IN_DELETE_SELF|unix.IN_MOVE_SELF|unix.IN_MOVED_FROM|unix.IN_MOVED_TO)
    add:=func(directory,relative string)error{wd,err:=unix.InotifyAddWatch(fd,directory,mask);if err==nil{directories[wd]=relative};return err}
    err=add(path,"");initialWalkError:=error(nil)
    if err==nil&&recursive{_=filepath.WalkDir(path,func(p string,d os.DirEntry,e error)error{if e!=nil{if initialWalkError==nil{initialWalkError=e};return nil};if d.IsDir()&&p!=path{relative,pathErr:=filepath.Rel(path,p);if pathErr!=nil{if initialWalkError==nil{initialWalkError=pathErr};return nil};if addErr:=add(p,relative);addErr!=nil&&initialWalkError==nil{initialWalkError=addErr}};return nil})}
    if err!=nil{
        unix.Close(fd)
        watchErr:=tsFSError(err,"watch",inputPath,"")
        tsSet(loop,watchErr,tsStringReference(tsStringUTF8("syscall")),tsStringReference(tsStringUTF8("watch")))
        tsSet(loop,watchErr,tsStringReference(tsStringUTF8("filename")),tsStringReference(tsStringUTF8(inputPath)))
        throwIfNoEntry:=tsFSProperty(options,"throwIfNoEntry")
        if errors.Is(err,syscall.ENOENT)&&!tsIsUndefined(throwIfNoEntry)&&!tsTruthy(throwIfNoEntry){state.closed=true;state.events.object.set("_handle",tsObjectValue(handle));return tsObjectValue(state.events.object)}
        panic(tsThrown{watchErr})
    }
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
				if flags&unix.IN_Q_OVERFLOW!=0{loop.post(func(){if !state.closed{tsFSWatchEvent(state,"change",tsNull)}});continue}
                if flags&unix.IN_IGNORED!=0{if relative==""&&( !isDirectory||!recursive){base:=filepath.Base(path);loop.post(func(){if !state.closed{tsFSWatchEvent(state,"rename",tsFSDecode([]byte(base),encoding))}})};delete(directories,wd);continue}
				newEntries:=[]string{}
				if recursive&&flags&unix.IN_ISDIR!=0&&flags&(unix.IN_CREATE|unix.IN_MOVED_TO)!=0{
					full:=filepath.Join(path,relative,filename);var addError error
					_=filepath.WalkDir(full,func(p string,d os.DirEntry,e error)error{
						if e!=nil{if addError==nil{addError=e};return nil}
						if d.IsDir(){rel,pathErr:=filepath.Rel(path,p);if pathErr!=nil{if addError==nil{addError=pathErr};return nil};if addErr:=add(p,rel);addErr!=nil&&addError==nil{addError=addErr}}
						if p!=full{rel,pathErr:=filepath.Rel(path,p);if pathErr==nil{newEntries=append(newEntries,rel)}else if addError==nil{addError=pathErr}}
						return nil
					})
					if addError!=nil{watchError:=tsFSError(addError,"watch",path,"");tsSet(loop,watchError,tsStringReference(tsStringUTF8("syscall")),tsStringReference(tsStringUTF8("watch")));loop.post(func(){if !state.closed{tsFSWatchEvent(state,"error",watchError,tsBooleanValue(false))}})}
				}
                if relative!=""{filename=filepath.Join(relative,filename)};if filename==""&&(!isDirectory||!recursive){filename=filepath.Base(path)}
                event:="change";if flags&(unix.IN_CREATE|unix.IN_DELETE|unix.IN_DELETE_SELF|unix.IN_MOVE_SELF|unix.IN_MOVED_FROM|unix.IN_MOVED_TO)!=0{event="rename"}
                eventValue:=tsStringReference(tsStringUTF8(event));fileValue:=tsU;if filename!=""{fileValue=tsFSDecode([]byte(filename),encoding)}
                loop.post(func(){if !state.closed{tsFSWatchEvent(state,tsText(eventValue),fileValue)}})
				for _,entry:=range newEntries{entryName:=entry;loop.post(func(){if !state.closed{tsFSWatchEvent(state,"rename",tsFSDecode([]byte(entryName),encoding))}})}
            }
        }
	},func(result tsResult){wasClosed:=state.closed;state.closed=true;if state.silentClose{return};if result.rejected{tsFSWatchEvent(state,"close");if !wasClosed{tsFSWatchEvent(state,"error",result.value,tsBooleanValue(false))};return};if !wasClosed&&signal!=nil{reason,wasAborted:=tsAbortReason(signal);if wasAborted{tsFSWatchEvent(state,"close");tsFSWatchEvent(state,"error",reason,tsBooleanValue(true));return}};tsFSWatchEvent(state,"close")})
    if initialWalkError!=nil{walkError:=tsFSError(initialWalkError,"watch",path,"");tsSet(loop,walkError,tsStringReference(tsStringUTF8("syscall")),tsStringReference(tsStringUTF8("watch")));loop.post(func(){if !state.closed{tsFSWatchEvent(state,"error",walkError,tsBooleanValue(false))}})}
    state.task.unref=!persistent
    loop.resourceClosers=append(loop.resourceClosers,state.cancelNative)
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
func tsFSWatchMetadata(info os.FileInfo)reflect.Value{
    if info==nil{return reflect.Value{}}
    metadata:=reflect.ValueOf(info.Sys());if metadata.Kind()==reflect.Pointer{metadata=metadata.Elem()};return metadata
}
func tsFSWatchTime(info os.FileInfo,field string,fallback time.Time)time.Time{
    metadata:=tsFSWatchMetadata(info);if metadata.Kind()!=reflect.Struct{return fallback};value:=metadata.FieldByName(field);if !value.IsValid()||value.Kind()!=reflect.Struct{return fallback};seconds,nanos:=value.FieldByName("Sec"),value.FieldByName("Nsec");if !seconds.IsValid()||!nanos.IsValid(){return fallback};return time.Unix(seconds.Int(),nanos.Int())
}
func tsFSWatchBirthtime(path string,info os.FileInfo)time.Time{
    var stat unix.Statx_t
    if unix.Statx(unix.AT_FDCWD,path,0,unix.STATX_BTIME,&stat)==nil&&stat.Mask&unix.STATX_BTIME!=0{return time.Unix(stat.Btime.Sec,int64(stat.Btime.Nsec))}
    return tsFSWatchTime(info,"Birthtimespec",time.Unix(0,0))
}
func tsFSWatchSameStat(path string,left,right os.FileInfo)bool{
    if left==nil||right==nil{return left==nil&&right==nil}
    if !os.SameFile(left,right)||left.Size()!=right.Size()||left.Mode()!=right.Mode(){return false}
    leftMeta,rightMeta:=tsFSWatchMetadata(left),tsFSWatchMetadata(right)
    for _,name:=range []string{"Dev","Ino","Mode","Nlink","Uid","Gid","Rdev","Blksize","Blocks"}{a,b:=leftMeta.FieldByName(name),rightMeta.FieldByName(name);if a.IsValid()!=b.IsValid(){return false};if a.IsValid()&&!reflect.DeepEqual(a.Interface(),b.Interface()){return false}}
    if !left.ModTime().Equal(right.ModTime()){return false}
    if !tsFSWatchTime(left,"Ctim",left.ModTime()).Equal(tsFSWatchTime(right,"Ctim",right.ModTime())){return false}
    return tsFSWatchBirthtime(path,left).Equal(tsFSWatchBirthtime(path,right))
}
func tsFSWatchStats(path string,info os.FileInfo,options tsValue)tsValue{
    if info==nil{return tsFSZeroStats(options)};value:=tsFSStats(info,options);tsFSBirthtime(value,unix.AT_FDCWD,path,0,options);return value
}
func tsFSStatStop(loop *tsLoop,state *tsFSWatchState){
    if state==nil||state.closed{return};state.close()
    process:=tsNodeProcessModule(loop);nextTick:=tsFSProperty(process,"nextTick")
    callback:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{state.events.emitEvent("stop");return tsU}))
    if nextTick.kind==tsFunctionKind{tsCallReceiver(loop,nextTick,process,callback)}else{loop.post(func(){state.events.emitEvent("stop")})}
}
func tsFSWatchFile(loop *tsLoop,args []tsValue)tsValue{
    path:=tsFSPath(tsArg(args,0));watchKey:=path;absolute,err:=filepath.Abs(path);if err!=nil{panic(tsThrown{tsFSError(err,"watchFile",path,"")})};path=absolute
    options:=tsArg(args,1);callback:=tsArg(args,2);if options.kind==tsFunctionKind{callback=options;options=tsU};if callback.kind!=tsFunctionKind{tsFSInvalid("listener","a function")}
    options=tsArg(tsFSPrepare(loop,[]tsValue{options}),0);persistent,bigint:=true,false;interval:=int32(5007)
    if options.kind==tsObjectKind||options.kind==tsInstanceKind{
        persistent=tsFSBooleanOption(options,"persistent",true);bigint=tsFSBooleanOption(options,"bigint",false)
        intervalValue:=tsFSProperty(options,"interval");if !tsIsUndefined(intervalValue){if !tsIsNumeric(intervalValue){tsFSInvalid("interval","a number")};interval=int32(tsNumberUint32(tsNumber(intervalValue)))}
    }
    if interval<5{interval=5}
    tsFSStatWatchers.Lock();if tsFSStatWatchers.states==nil{tsFSStatWatchers.states=map[*tsLoop]map[string]*tsFSWatchState{}};if tsFSStatWatchers.states[loop]==nil{tsFSStatWatchers.states[loop]=map[string]*tsFSWatchState{}};state:=tsFSStatWatchers.states[loop][watchKey];tsFSStatWatchers.Unlock()
    if state!=nil&&!state.closed{state.events.addEvent("change",callback,false);return tsObjectValue(state.events.object)}
    events:=tsFSNodeEventObject(loop);object:=events.object;object.nativeClass=tsNativeClass("StatWatcher")
    state=&tsFSWatchState{cancel:make(chan struct{}),events:events};events.addEvent("change",callback,false)
    handle:=tsNewObject();handle.set("ref",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{if state.task!=nil{state.task.unref=false};return tsU})));handle.set("unref",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{if state.task!=nil{state.task.unref=true};return tsU})));handle.set("close",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{state.close();return tsU})));object.set("_handle",tsObjectValue(handle))
    object.set("start",tsNamedNativeMethod("start",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsU}))
    object.set("stop",tsNamedNativeMethod("stop",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{if !state.closed{object.set("_handle",tsNull);tsFSStatStop(loop,state)};return tsU}))
    object.set("ref",tsNamedNativeMethod("ref",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{if !state.closed&&state.task!=nil{state.task.unref=false};return tsU}))
    object.set("unref",tsNamedNativeMethod("unref",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{if !state.closed&&state.task!=nil{state.task.unref=true};return tsU}))
    tsFSStatWatchers.Lock();tsFSStatWatchers.states[loop][watchKey]=state;tsFSStatWatchers.Unlock()
    state.task=loop.submit(func()tsResult{
        previous,previousErr:=os.Stat(path);if previousErr!=nil{previous=nil};previousValue:=tsFSWatchStats(path,previous,options)
        if previous==nil{old:=previousValue;loop.post(func(){if !state.closed{state.events.emitEvent("change",old,old)}})}
        ticker:=time.NewTicker(time.Duration(interval)*time.Millisecond);defer ticker.Stop()
        for{select{case <-state.cancel:return tsResult{tsU,false};case <-ticker.C:}
            current,currentErr:=os.Stat(path);if currentErr!=nil{current=nil};if !tsFSWatchSameStat(path,previous,current){value:=tsFSWatchStats(path,current,options);old:=previousValue;loop.post(func(){if !state.closed{state.events.emitEvent("change",value,old)}});previousValue=value};previous=current
        }
    },func(result tsResult){state.closed=true;if result.rejected{state.events.emitEvent("error",result.value)}})
    state.task.unref=!persistent
    loop.resourceClosers=append(loop.resourceClosers,func(){state.cancelNative();tsFSStatWatchers.Lock();delete(tsFSStatWatchers.states[loop],watchKey);if len(tsFSStatWatchers.states[loop])==0{delete(tsFSStatWatchers.states,loop)};tsFSStatWatchers.Unlock()})
    return tsObjectValue(object)
}
func tsFSUnwatchFile(loop *tsLoop,args []tsValue)tsValue{
    path:=tsFSPath(tsArg(args,0));callback:=tsArg(args,1)
    tsFSStatWatchers.Lock();state:=tsFSStatWatchers.states[loop][path];tsFSStatWatchers.Unlock();if state==nil{return tsU}
    if !tsTruthy(callback){state.events.removeAllEvent("change")}else{state.events.removeEvent("change",callback)}
    if state.events.listenerCount("change")==0{object:=state.events.object;object.set("_handle",tsNull);tsFSStatStop(loop,state);tsFSStatWatchers.Lock();delete(tsFSStatWatchers.states[loop],path);tsFSStatWatchers.Unlock()};return tsU
}
func tsFSZeroStats(options tsValue)tsValue{
    bigint:=tsTruthy(tsFSProperty(options,"bigint"));out:=tsNewObject();class:="Stats";if bigint{class="BigIntStats"};out.nativeClass=tsNativeClass(class);for _,key:=range []string{"dev","ino","mode","nlink","uid","gid","rdev","size","blksize","blocks","atimeMs","mtimeMs","ctimeMs","birthtimeMs"}{zero:=tsNumberValue(0);if bigint{zero=tsBigIntSigned(0)};out.set(key,zero)};for _,key:=range []string{"atime","mtime","ctime","birthtime"}{tsFSTimestamp(out,key,time.Unix(0,0),bigint)}
    return tsObjectValue(out)
}
`
