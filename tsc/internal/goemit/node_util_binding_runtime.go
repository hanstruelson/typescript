package goemit

// nodeUtilBindingRuntime translates the Rust node:util host bindings using
// standard Go errors, string processing, and the local Buffer encoder table.
const nodeUtilBindingRuntime = `
type tsNodeUtilSystemError struct { code int32; name string }
func tsNodeUtilSystemErrors() []tsNodeUtilSystemError {
    return []tsNodeUtilSystemError{
        {-int32(syscall.E2BIG),"E2BIG"},{-int32(syscall.EACCES),"EACCES"},{-int32(syscall.EADDRINUSE),"EADDRINUSE"},{-int32(syscall.EADDRNOTAVAIL),"EADDRNOTAVAIL"},{-int32(syscall.EAFNOSUPPORT),"EAFNOSUPPORT"},{-int32(syscall.EAGAIN),"EAGAIN"},{-int32(syscall.EALREADY),"EALREADY"},{-int32(syscall.EBADF),"EBADF"},{-int32(syscall.EBADMSG),"EBADMSG"},{-int32(syscall.EBUSY),"EBUSY"},{-int32(syscall.ECANCELED),"ECANCELED"},{-int32(syscall.ECHILD),"ECHILD"},{-int32(syscall.ECONNABORTED),"ECONNABORTED"},{-int32(syscall.ECONNREFUSED),"ECONNREFUSED"},{-int32(syscall.ECONNRESET),"ECONNRESET"},{-int32(syscall.EDEADLK),"EDEADLK"},{-int32(syscall.EDESTADDRREQ),"EDESTADDRREQ"},{-int32(syscall.EDOM),"EDOM"},{-int32(syscall.EEXIST),"EEXIST"},{-int32(syscall.EFAULT),"EFAULT"},{-int32(syscall.EFBIG),"EFBIG"},{-int32(syscall.EHOSTUNREACH),"EHOSTUNREACH"},{-int32(syscall.EIDRM),"EIDRM"},{-int32(syscall.EILSEQ),"EILSEQ"},{-int32(syscall.EINPROGRESS),"EINPROGRESS"},{-int32(syscall.EINTR),"EINTR"},{-int32(syscall.EINVAL),"EINVAL"},{-int32(syscall.EIO),"EIO"},{-int32(syscall.EISCONN),"EISCONN"},{-int32(syscall.EISDIR),"EISDIR"},{-int32(syscall.ELOOP),"ELOOP"},{-int32(syscall.EMFILE),"EMFILE"},{-int32(syscall.EMLINK),"EMLINK"},{-int32(syscall.EMSGSIZE),"EMSGSIZE"},{-int32(syscall.ENAMETOOLONG),"ENAMETOOLONG"},{-int32(syscall.ENETDOWN),"ENETDOWN"},{-int32(syscall.ENETRESET),"ENETRESET"},{-int32(syscall.ENETUNREACH),"ENETUNREACH"},{-int32(syscall.ENFILE),"ENFILE"},{-int32(syscall.ENOBUFS),"ENOBUFS"},{-int32(syscall.ENODEV),"ENODEV"},{-int32(syscall.ENOENT),"ENOENT"},{-int32(syscall.ENOEXEC),"ENOEXEC"},{-int32(syscall.ENOLCK),"ENOLCK"},{-int32(syscall.ENOMEM),"ENOMEM"},{-int32(syscall.ENOMSG),"ENOMSG"},{-int32(syscall.ENOPROTOOPT),"ENOPROTOOPT"},{-int32(syscall.ENOSPC),"ENOSPC"},{-int32(syscall.ENOSYS),"ENOSYS"},{-int32(syscall.ENOTCONN),"ENOTCONN"},{-int32(syscall.ENOTDIR),"ENOTDIR"},{-int32(syscall.ENOTEMPTY),"ENOTEMPTY"},{-int32(syscall.ENOTSOCK),"ENOTSOCK"},{-int32(syscall.ENOTSUP),"ENOTSUP"},{-int32(syscall.ENOTTY),"ENOTTY"},{-int32(syscall.ENXIO),"ENXIO"},{-int32(syscall.EOVERFLOW),"EOVERFLOW"},{-int32(syscall.EPERM),"EPERM"},{-int32(syscall.EPIPE),"EPIPE"},{-int32(syscall.EPROTO),"EPROTO"},{-int32(syscall.EPROTONOSUPPORT),"EPROTONOSUPPORT"},{-int32(syscall.EPROTOTYPE),"EPROTOTYPE"},{-int32(syscall.ERANGE),"ERANGE"},{-int32(syscall.EROFS),"EROFS"},{-int32(syscall.ESPIPE),"ESPIPE"},{-int32(syscall.ESRCH),"ESRCH"},{-int32(syscall.ETIMEDOUT),"ETIMEDOUT"},{-int32(syscall.ETXTBSY),"ETXTBSY"},{-int32(syscall.EXDEV),"EXDEV"},
    }
}
func tsNodeUtilBindingInternalErrorName(value tsValue) tsValue {
    number:=int32(tsNumber(value))
    for _,entry:=range tsNodeUtilSystemErrors(){if entry.code==number{return tsStringReference(tsStringUTF8(entry.name))}}
    return tsStringReference(tsStringUTF8("Unknown system error "+strconv.Itoa(int(number))))
}
func tsNodeUtilBindingInternalErrorEntries() tsValue {
    result:=&tsArray{}
    for _,entry:=range tsNodeUtilSystemErrors(){result.values=append(result.values,tsNumberValue(float64(entry.code)),tsStringReference(tsStringUTF8(entry.name)))}
    return tsArrayValue(result)
}
func tsNodeUtilBindingUVError(name string) tsValue {
    for _,entry:=range tsNodeUtilSystemErrors(){if entry.name==name{return tsNumberValue(float64(entry.code))}}
    return tsNumberValue(-float64(syscall.EINVAL))
}
func tsNodeUtilBindingTranslateSysError(value tsValue) tsValue {
    if !tsIsNumeric(value){return tsNumberValue(-float64(syscall.EINVAL))}
    number:=int32(tsNumber(value));if number<=0{return tsNumberValue(float64(number))}
    if runtime.GOOS=="windows" {return tsNumberValue(-float64(syscall.EINVAL))}
    return tsNumberValue(-float64(number))
}
func tsNodeUtilBindingSplitNewLines(loop *tsLoop,value tsValue) tsValue {
    if value.kind!=tsStringKind{tsPropertyFailure("extractedSplitNewLines expects a string")}
    text:=tsStringValueWithLoop(loop,value).String();out:=&tsArray{}
    if text==""{out.values=append(out.values,tsStringReference(tsStringUTF8("")));return tsArrayValue(out)}
    start:=0
    for index:=0;index<len(text);index++{if text[index]=='\n'{out.values=append(out.values,tsStringReference(tsStringUTF8(text[start:index+1])));start=index+1}}
    if start<len(text){out.values=append(out.values,tsStringReference(tsStringUTF8(text[start:]))) }
    return tsArrayValue(out)
}
func tsNodeUtilBindingNormalizeEncoding(loop *tsLoop,value tsValue) tsValue {
    name:=strings.ToLower(tsNodeUtilBindingText(loop,value))
    if name==""{return tsStringReference(tsStringUTF8("utf8"))}
    canonical:=map[string]string{"utf8":"utf8","utf-8":"utf8","ucs2":"ucs2","ucs-2":"ucs2","utf16le":"utf16le","utf16-le":"utf16le","latin1":"latin1","binary":"latin1","ascii":"ascii","base64":"base64","base64url":"base64url","hex":"hex","buffer":"buffer"}
    if normalized,ok:=canonical[name];ok{return tsStringReference(tsStringUTF8(normalized))}
    return tsU
}
func tsNodeUtilBindingText(loop *tsLoop,value tsValue) string {return tsStringValueWithLoop(loop,value).String()}
func tsNodeUtilBindingParseEnv(loop *tsLoop,value tsValue) tsValue {
    if value.kind!=tsStringKind{tsInternalValidatorInvalidType("content","a string",value)}
    text:=strings.TrimPrefix(tsText(value),"\ufeff");out:=tsNewObject()
    for _,rawLine:=range strings.Split(text,"\n") {
        line:=strings.TrimSpace(strings.TrimSuffix(rawLine,"\r"));if line==""||strings.HasPrefix(line,"#"){continue}
        if strings.HasPrefix(line,"export "){line=strings.TrimSpace(strings.TrimPrefix(line,"export "))}
        split:=strings.IndexByte(line,'=');if split<0{continue}
        key:=strings.TrimSpace(line[:split]);if key==""{continue}
        valid:=true;for index,char:=range key{if !(char=='_'||char>='A'&&char<='Z'||char>='a'&&char<='z'||index>0&&char>='0'&&char<='9'){valid=false;break}}
        if !valid{continue}
        raw:=strings.TrimSpace(line[split+1:]);valueText:=raw
        if len(raw)>0&&(raw[0]=='\''||raw[0]=='"'||raw[0]==96) {
            quote:=raw[0];end:=strings.LastIndexByte(raw[1:],quote);if end>=0{valueText=raw[1:1+end];if quote=='"'{valueText=strings.NewReplacer("\\n","\n","\\r","\r","\\t","\t","\\\"","\"","\\\\","\\").Replace(valueText)}} else {valueText=raw[1:]}
        } else {
            if comment:=strings.Index(raw," #");comment>=0{raw=raw[:comment]};valueText=strings.TrimSpace(raw)
        }
        out.set(key,tsStringReference(tsStringUTF8(valueText)))
    }
    return tsObjectValue(out)
}
func tsNodeUtilBindingModule(loop *tsLoop) tsValue {
    out:=tsNewObject()
    out.set("internalErrorName",tsNamedNativeMethod("internalErrorName",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeUtilBindingInternalErrorName(tsArg(args,0))}))
    out.set("internalErrorEntries",tsNamedNativeMethod("internalErrorEntries",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNodeUtilBindingInternalErrorEntries()}))
    out.set("etimedoutErrorCode",tsNamedNativeMethod("etimedoutErrorCode",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNodeUtilBindingUVError("ETIMEDOUT")}))
    out.set("enobufsErrorCode",tsNamedNativeMethod("enobufsErrorCode",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNodeUtilBindingUVError("ENOBUFS")}))
    out.set("ecanceledErrorCode",tsNamedNativeMethod("ecanceledErrorCode",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNodeUtilBindingUVError("ECANCELED")}))
    out.set("uvTranslateSysError",tsNamedNativeMethod("uvTranslateSysError",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeUtilBindingTranslateSysError(tsArg(args,0))}))
    out.set("extractedSplitNewLinesFastPathStringsOnly",tsNamedNativeMethod("extractedSplitNewLinesFastPathStringsOnly",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeUtilBindingSplitNewLines(loop,tsArg(args,0))}))
    out.set("normalizeEncoding",tsNamedNativeMethod("normalizeEncoding",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeUtilBindingNormalizeEncoding(loop,tsArg(args,0))}))
    out.set("parseEnv",tsNamedNativeMethod("parseEnv",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeUtilBindingParseEnv(loop,tsArg(args,0))}))
    return tsObjectValue(out)
}
`
