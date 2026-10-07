package goemit

const fsOperationRuntime = `
func tsFSExecute(loop *tsLoop, operation string, args []tsValue, asynchronous bool) (result tsResult) {
    result.value=tsU
    defer func(){if failure:=recover();failure!=nil{result=tsResult{tsUnwrap(failure),true}}}()
    value:=tsArg(args,0)
    options:=tsArg(args,1)
    var err error
    var path,dest string
    switch operation {
    case "glob":return tsFSGlob(loop,args,asynchronous)
    case "read","write","readv","writev":
        return tsFSDescriptorIO(operation,args)
    case "close","fchmod","fchown","fdatasync","fstat","fsync","ftruncate","futimes":
        record:=tsFSDescriptor(value)
        record.Lock()
        defer record.Unlock()
        if record.closed {return tsResult{tsFSError(syscall.EBADF,operation,"",""),true}}
        switch operation {
        case "close":
            err=record.file.Close();record.closed=true
            tsFSFiles.Lock();delete(tsFSFiles.files,record.fd);tsFSFiles.Unlock()
        case "fchmod":err=record.file.Chmod(tsFSMode(options,0))
        case "fchown":err=record.file.Chown(int(tsFSInteger(options,"uid",0,-1)),int(tsFSInteger(tsArg(args,2),"gid",0,-1)))
        case "fdatasync":err=unix.Fdatasync(record.fd)
        case "fsync":err=record.file.Sync()
        case "ftruncate":err=record.file.Truncate(tsFSInteger(options,"length",0,0))
        case "futimes":times:=[]unix.Timeval{unix.NsecToTimeval(tsFSTime(options).UnixNano()),unix.NsecToTimeval(tsFSTime(tsArg(args,2)).UnixNano())};err=unix.Futimes(record.fd,times)
        case "fstat":var info os.FileInfo;info,err=record.file.Stat();if err==nil{result.value=tsFSStats(info,options);tsFSBirthtime(result.value,record.fd,"",unix.AT_EMPTY_PATH,options)}
        }
    case "readFile","writeFile","appendFile":
        opts:=options;if operation!="readFile"{opts=tsArg(args,2)}
        var signal *tsAbortState;if asynchronous{signal=tsFSAbortSignal(opts);tsFSCheckAbort(signal)}
        defaultFlag:="r";if operation=="writeFile"{defaultFlag="w"};if operation=="appendFile"{defaultFlag="a"}
        flag:=tsFSOption(opts,"flag",tsStringReference(tsStringUTF8(defaultFlag)))
        record,owned,openErr:=tsFSFile(value,tsFSFlags(flag),tsFSMode(tsFSProperty(opts,"mode"),0666))
        if openErr!=nil{return tsResult{tsFSError(openErr,"open",tsFSPath(value),""),true}}
        if owned{defer record.file.Close()}
        record.Lock();defer record.Unlock()
        if record.closed{return tsResult{tsFSError(syscall.EBADF,operation,"",""),true}}
        if operation=="readFile" {
            var data []byte
            var output bytes.Buffer;chunk:=make([]byte,512*1024);for{tsFSCheckAbort(signal);n,e:=record.file.Read(chunk);if n>0{output.Write(chunk[:n])};if e!=nil{if !errors.Is(e,io.EOF){err=e};break};if n==0{break}};tsFSCheckAbort(signal);data=output.Bytes()
            if err==nil{result.value=tsFSDecode(data,tsFSEncoding(opts,"buffer"))}
        }else{
            data:=tsFSData(options,opts)
            for len(data)>0 {tsFSCheckAbort(signal);var n int;piece:=data;if len(piece)>512*1024{piece=piece[:512*1024]};n,err=record.file.Write(piece);if n>0{data=data[n:]};if err!=nil{break};if n==0{err=io.ErrShortWrite;break}}
            if err==nil&&tsTruthy(tsFSProperty(opts,"flush")){err=record.file.Sync()}
        }
        path=record.path
    case "open":
        path=tsFSPath(value)
        flags:=options;if tsIsUndefined(flags){flags=tsStringReference(tsStringUTF8("r"))}
        var file *os.File
        file,err=os.OpenFile(path,tsFSFlags(flags),tsFSMode(tsArg(args,2),0666))
        if err==nil{result.value=tsNumberValue(float64(tsFSRegister(file,path)))}
    case "exists":
        path=tsFSPath(value);_,err=os.Stat(path);result.value=tsBooleanValue(err==nil);err=nil
    case "access":
        path=tsFSPath(value);mode:=tsFSInteger(options,"mode",0,0)
        if mode>7{panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","mode must be between 0 and 7")})}
        err=unix.Access(path,uint32(mode))
    case "stat","lstat":
        path=tsFSPath(value)
        var info os.FileInfo
        if operation=="stat"{info,err=os.Stat(path)}else{info,err=os.Lstat(path)}
        if err!=nil&&tsFSProperty(options,"throwIfNoEntry").kind==tsBooleanKind&&!tsTruthy(tsFSProperty(options,"throwIfNoEntry"))&&errors.Is(err,os.ErrNotExist){err=nil;break}
        if err==nil{result.value=tsFSStats(info,options);flags:=0;if operation=="lstat"{flags=unix.AT_SYMLINK_NOFOLLOW};tsFSBirthtime(result.value,unix.AT_FDCWD,path,flags,options)}
    case "statfs":
        path=tsFSPath(value);var info unix.Statfs_t;err=unix.Statfs(path,&info)
        if err==nil{out:=tsNewObject();metadata:=reflect.ValueOf(info);bigint:=tsTruthy(tsFSProperty(options,"bigint"));class:="StatFs";if bigint{class="BigIntStatFs"};out.nativeClass=tsNativeClass(class);for key,field:=range map[string]string{"type":"Type","bsize":"Bsize","blocks":"Blocks","bfree":"Bfree","bavail":"Bavail","files":"Files","ffree":"Ffree"}{out.set(key,tsFSMetadataField(metadata.FieldByName(field),bigint))};result.value=tsObjectValue(out)}

    case "mkdir":
        path=tsFSPath(value);modeValue:=tsFSProperty(options,"mode")
        if tsIsNumeric(options)||options.kind==tsStringKind{modeValue=options}
        mode:=tsFSMode(modeValue,0777)
        if tsTruthy(tsFSProperty(options,"recursive")){
            var first string;first,err=tsFSMkdirAll(path,mode);if err==nil&&first!=""{result.value=tsStringReference(tsStringUTF8(first))}
        }else{err=os.Mkdir(path,mode)}
    case "mkdtemp","mkdtempDisposable":
        path=tsFSPath(value);var created string
        created,err=tsFSMkdtemp(path)
        if err==nil{result.value=tsFSDecode([]byte(created),tsFSEncoding(options,"utf8"))}
        if err==nil&&operation=="mkdtempDisposable"{
            absolute,_:=filepath.Abs(created);out:=tsNewObject();out.set("path",result.value)
            remove:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{
                opts:=tsNewObject();opts.set("recursive",tsBooleanValue(true));opts.set("force",tsBooleanValue(true))
                mode:="sync";if asynchronous{mode="promise"};return tsFSInvoke(loop,"rm",mode,[]tsValue{tsStringReference(tsStringUTF8(absolute)),tsObjectValue(opts)})
            }))
            out.set("remove",remove);result.value=tsObjectValue(out)
        }
    case "readdir":
        path=tsFSPath(value);result.value,err=tsFSReadDirectory(path,options)
    case "opendir":
        path=tsFSPath(value);var file *os.File;file,err=os.Open(path)
        if err==nil{var info os.FileInfo;info,err=file.Stat();if err!=nil||!info.IsDir(){file.Close();if err==nil{err=syscall.ENOTDIR}}}
        if err==nil{result.value=tsFSDir(file,path,options)}
    case "readlink":path=tsFSPath(value);var text string;text,err=os.Readlink(path);if err==nil{result.value=tsFSDecode([]byte(text),tsFSEncoding(options,"utf8"))}
    case "realpath":path=tsFSPath(value);var text string;text,err=filepath.EvalSymlinks(path);if err==nil{text,err=filepath.Abs(text)};if err==nil{result.value=tsFSDecode([]byte(text),tsFSEncoding(options,"utf8"))}
    case "rename","link","symlink":
        path,dest=tsFSPath(value),tsFSPath(options)
        switch operation{case "rename":err=os.Rename(path,dest);case "link":err=os.Link(path,dest);case "symlink":err=os.Symlink(path,dest)}
    case "unlink":
        path=tsFSPath(value);var info os.FileInfo;info,err=os.Lstat(path);if err==nil&&info.IsDir(){err=syscall.EISDIR};if err==nil{err=os.Remove(path)}
    case "rmdir","rm":
        path=tsFSPath(value);recursive:=tsTruthy(tsFSProperty(options,"recursive"));force:=tsTruthy(tsFSProperty(options,"force"))
        var info os.FileInfo;info,err=os.Lstat(path)
        if err!=nil&&operation=="rm"&&force&&errors.Is(err,os.ErrNotExist){err=nil;break}
        if err==nil{
            if operation=="rmdir"&&!info.IsDir(){err=syscall.ENOTDIR
            }else if operation=="rm"&&info.IsDir()&&!recursive{
                failure:=tsFSError(syscall.EISDIR,"rm",path,"");object:=(*tsRuntimeError)(failure.ref);object.name="SystemError";object.fields["code"]=tsStringReference(tsStringUTF8("ERR_FS_EISDIR"));return tsResult{failure,true}
            }else if recursive{err=tsFSRemoveAll(path,options)}else{err=os.Remove(path)}
        }
    case "truncate":path=tsFSPath(value);err=os.Truncate(path,tsFSInteger(options,"length",0,0))
    case "chmod","lchmod":path=tsFSPath(value);if operation=="lchmod"{err=unix.Fchmodat(unix.AT_FDCWD,path,uint32(tsFSModeBits(options,0)),unix.AT_SYMLINK_NOFOLLOW)}else{err=os.Chmod(path,tsFSMode(options,0))}
    case "chown","lchown":path=tsFSPath(value);uid,gid:=int(tsFSInteger(options,"uid",0,-1)),int(tsFSInteger(tsArg(args,2),"gid",0,-1));if operation=="chown"{err=os.Chown(path,uid,gid)}else{err=os.Lchown(path,uid,gid)}
    case "utimes","lutimes":path=tsFSPath(value);atime,mtime:=tsFSTime(options),tsFSTime(tsArg(args,2));if operation=="utimes"{err=os.Chtimes(path,atime,mtime)}else{err=unix.UtimesNanoAt(unix.AT_FDCWD,path,[]unix.Timespec{unix.NsecToTimespec(atime.UnixNano()),unix.NsecToTimespec(mtime.UnixNano())},unix.AT_SYMLINK_NOFOLLOW)}
    case "copyFile":path,dest=tsFSPath(value),tsFSPath(options);err=tsFSCopyFile(path,dest,int(tsFSInteger(tsArg(args,2),"mode",0,0)))
    case "cp":path,dest=tsFSPath(value),tsFSPath(options);err=tsFSCopy(loop,path,dest,tsArg(args,2),asynchronous)
    default:panic(tsThrown{tsFSException("Error","ERR_METHOD_NOT_IMPLEMENTED","Filesystem operation is not implemented: "+operation)})
    }
    if err!=nil{return tsResult{tsFSError(err,operation,path,dest),true}}
    return result
}

func tsFSModeBits(value tsValue, fallback int64) int64 {
    if value.kind==tsStringKind {number,err:=strconv.ParseInt((*tsString)(value.ref).String(),8,32);if err!=nil{tsFSInvalid("mode","an octal string or integer")};value=tsNumberValue(float64(number))}
    number:=tsFSInteger(value,"mode",fallback,0)
    if number>07777{panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","mode is out of range")})};return number
}
func tsFSMode(value tsValue, fallback int64) os.FileMode {
    n:=tsFSModeBits(value,fallback);mode:=os.FileMode(n&0777)
    if n&04000!=0{mode|=os.ModeSetuid};if n&02000!=0{mode|=os.ModeSetgid};if n&01000!=0{mode|=os.ModeSticky};return mode
}

func tsFSFlags(value tsValue) int {
    if tsIsNumeric(value){return int(tsFSInteger(value,"flags",0,0))}
    if value.kind!=tsStringKind{tsFSInvalid("flags","a string or number")}
    switch tsText(value){
    case "r":return os.O_RDONLY
    case "r+":return os.O_RDWR
    case "rs","sr":return os.O_RDONLY|os.O_SYNC
    case "rs+","sr+":return os.O_RDWR|os.O_SYNC
    case "w":return os.O_WRONLY|os.O_CREATE|os.O_TRUNC
    case "wx","xw":return os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_EXCL
    case "w+":return os.O_RDWR|os.O_CREATE|os.O_TRUNC
    case "wx+","xw+":return os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_EXCL
    case "a":return os.O_WRONLY|os.O_CREATE|os.O_APPEND
    case "ax","xa":return os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_EXCL
    case "a+":return os.O_RDWR|os.O_CREATE|os.O_APPEND
    case "ax+","xa+":return os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_EXCL
    case "as","sa":return os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_SYNC
    case "as+","sa+":return os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_SYNC
    }
    panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_VALUE","Unknown file flag: "+tsText(value))})
}

func tsFSTime(value tsValue) time.Time {
    if value.kind==tsObjectKind{object:=(*tsObject)(value.ref);if object.nativeClass!=nil&&object.nativeClass.builtin=="Date"{value=tsNumberValue(object.nativeTime);if math.IsNaN(value.number){tsFSInvalid("time","a valid Date")};return tsDateInstant(value.number)};value=tsFSProperty(value,"_milliseconds");if !tsIsUndefined(value){return time.Unix(0,int64(tsNumber(value)*1000000))}}
    if value.kind==tsStringKind{number,err:=strconv.ParseFloat(tsText(value),64);if err!=nil{tsFSInvalid("time","a number, numeric string, or Date")};value=tsNumberValue(number)}
    if !tsIsNumeric(value){tsFSInvalid("time","a number, numeric string, or Date")}
    seconds:=tsNumber(value);if math.IsNaN(seconds)||math.IsInf(seconds,0){tsFSInvalid("time","a finite number")}
    if seconds<0{seconds=float64(time.Now().UnixNano())/1e9};whole,fraction:=math.Modf(seconds);return time.Unix(int64(whole),int64(fraction*1e9))
}




func tsFSMetadataField(field reflect.Value, bigint bool)tsValue{
    if !field.IsValid(){if bigint{return tsBigIntSigned(0)};return tsNumberValue(0)}
    switch field.Kind(){case reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64:if bigint{return tsBigIntSigned(field.Int())};return tsNumberValue(float64(field.Int()));case reflect.Uint,reflect.Uint8,reflect.Uint16,reflect.Uint32,reflect.Uint64:if bigint{return tsBigIntUnsigned(field.Uint())};return tsNumberValue(float64(field.Uint()))}
    panic("Invalid native filesystem metadata field")
}
func tsFSTimestamp(out *tsObject,name string,stamp time.Time,bigint bool){
    if bigint{nanoseconds:=new(big.Int).Mul(big.NewInt(stamp.Unix()),big.NewInt(1000000000));nanoseconds.Add(nanoseconds,big.NewInt(int64(stamp.Nanosecond())));out.set(name+"Ns",tsBigIntValue(nanoseconds));out.set(name+"Ms",tsBigIntValue(new(big.Int).Div(new(big.Int).Set(nanoseconds),big.NewInt(1000000))))}else{out.set(name+"Ms",tsNumberValue(float64(stamp.Unix())*1000+float64(stamp.Nanosecond())/1e6))};if out.nativeDates!=nil{delete(out.nativeDates,name)}
}
func tsFSStats(info os.FileInfo, options tsValue) tsValue {
    bigint:=tsTruthy(tsFSProperty(options,"bigint"));out:=tsNewObject();class:="Stats";if bigint{class="BigIntStats"};out.nativeClass=tsNativeClass(class);metadata:=reflect.ValueOf(info.Sys());if metadata.Kind()==reflect.Pointer{metadata=metadata.Elem()}
    for name,field:=range map[string]string{"dev":"Dev","ino":"Ino","mode":"Mode","nlink":"Nlink","uid":"Uid","gid":"Gid","rdev":"Rdev","blksize":"Blksize","blocks":"Blocks"}{out.set(name,tsFSMetadataField(metadata.FieldByName(field),bigint))}
    if bigint{out.set("size",tsBigIntSigned(info.Size()))}else{out.set("size",tsNumberValue(float64(info.Size())))}
    stamp:=func(name string,fallback time.Time)time.Time{if metadata.Kind()!=reflect.Struct{return fallback};field:=metadata.FieldByName(name);if !field.IsValid()||field.Kind()!=reflect.Struct{return fallback};seconds,nanos:=field.FieldByName("Sec"),field.FieldByName("Nsec");if !seconds.IsValid()||!nanos.IsValid(){return fallback};return time.Unix(seconds.Int(),nanos.Int())}
    times:=map[string]time.Time{"mtime":info.ModTime(),"atime":stamp("Atim",info.ModTime()),"ctime":stamp("Ctim",info.ModTime()),"birthtime":stamp("Birthtimespec",time.Unix(0,0))}
    for name,t:=range times{tsFSTimestamp(out,name,t,bigint)}
    tsFSSetTypePredicates(out,info.Mode());return tsObjectValue(out)
}
func tsFSBirthtime(value tsValue,fd int,path string,flags int,options tsValue){
    var metadata unix.Statx_t
    if unix.Statx(fd,path,flags,unix.STATX_BTIME,&metadata)==nil&&metadata.Mask&unix.STATX_BTIME!=0{stamp:=time.Unix(metadata.Btime.Sec,int64(metadata.Btime.Nsec));tsFSTimestamp((*tsObject)(value.ref),"birthtime",stamp,tsTruthy(tsFSProperty(options,"bigint")))}
}

func tsFSSetTypePredicates(out *tsObject,mode os.FileMode){out.nativeMode=mode;if out.nativeClass==nil{out.nativeClass=tsNativeClass("Dirent")}}

func tsFSMkdirAll(path string, mode os.FileMode) (string,error) {
    if info,err:=os.Stat(path);err==nil{if !info.IsDir(){return "",syscall.ENOTDIR};return "",nil}else if !errors.Is(err,os.ErrNotExist){return "",err}
    parent:=filepath.Dir(path);first:=""
    if parent!=path{var err error;first,err=tsFSMkdirAll(parent,mode);if err!=nil{return "",err}}
    if err:=os.Mkdir(path,mode);err!=nil{if info,e:=os.Stat(path);e==nil&&info.IsDir(){return first,nil};return "",err}
    if first==""{first=path};return first,nil
}

func tsFSMkdtemp(prefix string) (string,error) {
    const alphabet="ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
    for i:=0;i<100;i++{var suffix [6]byte;if _,err:=cryptorand.Read(suffix[:]);err!=nil{return "",err};for j,b:=range suffix{suffix[j]=alphabet[int(b)%len(alphabet)]};path:=prefix+string(suffix[:]);err:=os.Mkdir(path,0700);if err==nil{return path,nil};if !errors.Is(err,os.ErrExist){return "",err}}
    return "",syscall.EEXIST
}
func tsFSRemoveAll(path string,options tsValue)error{
    retries:=tsFSInteger(tsFSProperty(options,"maxRetries"),"maxRetries",0,0);delay:=tsFSInteger(tsFSProperty(options,"retryDelay"),"retryDelay",100,0)
    for attempt:=int64(0);;attempt++{err:=os.RemoveAll(path);if err==nil{return nil};retry:=errors.Is(err,syscall.EBUSY)||errors.Is(err,syscall.EMFILE)||errors.Is(err,syscall.ENFILE)||errors.Is(err,syscall.ENOTEMPTY)||errors.Is(err,syscall.EPERM);if !retry||attempt>=retries{return err};time.Sleep(time.Duration(attempt+1)*time.Duration(delay)*time.Millisecond)}
}
`
