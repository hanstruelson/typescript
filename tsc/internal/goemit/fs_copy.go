package goemit

const fsCopyRuntime = `
func tsFSCopyFile(source, destination string, mode int) (err error) {
    if mode<0||mode>7{panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","copyFile mode must be between 0 and 7")})}
    input,err:=os.Open(source);if err!=nil{return err};defer input.Close()
    info,err:=input.Stat();if err!=nil{return err};if info.IsDir(){return syscall.EISDIR}
    if target,e:=os.Stat(destination);e==nil&&os.SameFile(info,target){return syscall.EINVAL}
    flags:=os.O_WRONLY|os.O_CREATE|os.O_TRUNC;if mode&1!=0{flags|=os.O_EXCL}
    output,err:=os.OpenFile(destination,flags,info.Mode().Perm());if err!=nil{return err}
    defer func(){closeErr:=output.Close();if err==nil{err=closeErr};if err!=nil{os.Remove(destination)}}()
    if mode&6!=0{err=unix.IoctlFileClone(int(output.Fd()),int(input.Fd()));if err==nil{return nil};if mode&4!=0{return err};err=nil}
    _,err=io.Copy(output,input)
    if err==nil{err=output.Chmod(info.Mode().Perm())}
    return err
}

func tsFSCopyFailure(code, message, source, destination string) {
    value:=tsFSException("SystemError",code,message)
    object:=(*tsRuntimeError)(value.ref)
    object.fields["path"]=tsStringReference(tsStringUTF8(source));object.fields["dest"]=tsStringReference(tsStringUTF8(destination));object.fields["syscall"]=tsStringReference(tsStringUTF8("cp"))
    panic(tsThrown{value})
}

func tsFSCopyFilter(loop *tsLoop, filter tsValue, source, destination string, asynchronous bool) bool {
    if tsIsUndefined(filter){return true};if filter.kind!=tsFunctionKind{tsFSInvalid("filter","a function")}
    args:=[]tsValue{tsStringReference(tsStringUTF8(source)),tsStringReference(tsStringUTF8(destination))}
    var result tsResult
    if !asynchronous{result.value=tsCall(loop,filter,args...)}else{
        complete:=make(chan tsResult,1)
        loop.post(func(){
            defer func(){if failure:=recover();failure!=nil{complete<-tsResult{tsUnwrap(failure),true}}}()
            loop.await(tsCall(loop,filter,args...),func(value tsResult){complete<-value})
        })
        result=<-complete
    }
    if result.rejected{panic(tsThrown{result.value})}
    if result.value.kind!=tsBooleanKind{panic(tsThrown{tsFSException("TypeError","ERR_INVALID_RETURN_VALUE","cp filter must return a boolean")})}
    return tsTruthy(result.value)
}

func tsFSCopy(loop *tsLoop, source, destination string, options tsValue, asynchronous bool) error {
    recursive:=tsTruthy(tsFSProperty(options,"recursive"));dereference:=tsTruthy(tsFSProperty(options,"dereference"))
    force:=tsTruthy(tsFSOption(options,"force",tsBooleanValue(true)));errorOnExist:=tsTruthy(tsFSProperty(options,"errorOnExist"))
    preserve:=tsTruthy(tsFSProperty(options,"preserveTimestamps"));verbatim:=tsTruthy(tsFSProperty(options,"verbatimSymlinks"))
    filter:=tsFSProperty(options,"filter");mode:=int(tsFSInteger(tsFSProperty(options,"mode"),"mode",0,0))
    absoluteSource,err:=filepath.Abs(source);if err!=nil{return err};absoluteDestination,err:=filepath.Abs(destination);if err!=nil{return err}
    if absoluteSource==absoluteDestination||strings.HasPrefix(absoluteDestination,absoluteSource+string(filepath.Separator)){
        tsFSCopyFailure("ERR_FS_CP_EINVAL","Cannot copy a path into itself",source,destination)
    }
    var copyPath func(string,string)error
    copyPath=func(src,dst string)error{
        if !tsFSCopyFilter(loop,filter,src,dst,asynchronous){return nil}
        var info os.FileInfo;var err error
        if dereference{info,err=os.Stat(src)}else{info,err=os.Lstat(src)};if err!=nil{return err}
        target,targetErr:=os.Lstat(dst)
        if targetErr!=nil&&!errors.Is(targetErr,os.ErrNotExist){return targetErr}
        if targetErr==nil&&os.SameFile(info,target){tsFSCopyFailure("ERR_FS_CP_EINVAL","Source and destination are the same file",src,dst)}
        if info.IsDir(){
            if !recursive{tsFSCopyFailure("ERR_FS_EISDIR","Recursive option is required to copy directories",src,dst)}
            if targetErr==nil&&!target.IsDir(){tsFSCopyFailure("ERR_FS_CP_DIR_TO_NON_DIR","Cannot overwrite a non-directory with a directory",src,dst)}
            if targetErr!=nil{if err:=os.MkdirAll(dst,info.Mode().Perm());err!=nil{return err}}
            entries,err:=os.ReadDir(src);if err!=nil{return err}
            for _,entry:=range entries{if err:=copyPath(filepath.Join(src,entry.Name()),filepath.Join(dst,entry.Name()));err!=nil{return err}}
            if preserve{return os.Chtimes(dst,info.ModTime(),info.ModTime())};return nil
        }
        if targetErr==nil&&target.IsDir(){tsFSCopyFailure("ERR_FS_CP_NON_DIR_TO_DIR","Cannot overwrite a directory with a non-directory",src,dst)}
        if targetErr==nil&&!force{if errorOnExist{tsFSCopyFailure("ERR_FS_CP_EEXIST","Destination already exists",src,dst)};return nil}
        if err:=os.MkdirAll(filepath.Dir(dst),0777);err!=nil{return err}
        if info.Mode()&os.ModeSymlink!=0{
            link,err:=os.Readlink(src);if err!=nil{return err}
            if !verbatim&&!filepath.IsAbs(link){link,err=filepath.Abs(filepath.Join(filepath.Dir(src),link));if err!=nil{return err}}
            if targetErr==nil{if err:=os.Remove(dst);err!=nil{return err}}
            return os.Symlink(link,dst)
        }
        if info.Mode()&os.ModeNamedPipe!=0{tsFSCopyFailure("ERR_FS_CP_FIFO_PIPE","Cannot copy a FIFO",src,dst)}
        if info.Mode()&os.ModeSocket!=0{tsFSCopyFailure("ERR_FS_CP_SOCKET","Cannot copy a socket",src,dst)}
        if !info.Mode().IsRegular(){tsFSCopyFailure("ERR_FS_CP_UNKNOWN","Cannot copy this file type",src,dst)}
        if err:=tsFSCopyFile(src,dst,mode);err!=nil{return err}
        if preserve{return os.Chtimes(dst,info.ModTime(),info.ModTime())};return nil
    }
    return copyPath(source,destination)
}
`
