package goemit

// FSRuntime implements the Node filesystem boundary using Go's file APIs.
// Native workers do I/O; JavaScript callbacks always run on their owning loop.
// Bun is a behavior reference, not a dependency of generated executables.
const FSRuntime = fsModuleRuntime + fsEncodingRuntime + fsBufferMethodsRuntime + fsOperationRuntime + fsHandleRuntime + fsDirectoryRuntime + fsCopyRuntime + fsIteratorRuntime + fsGlobRuntime + fsWatchRuntime + fsStreamRuntime + abortRuntime + webStreamRuntime + fsBlobRuntime + fsReaderRuntime

const fsModuleRuntime = `
var tsFSModules struct {
    sync.Mutex
    fs, promises *tsObject
}

var tsFSOperations = []string{
    "access", "appendFile", "chmod", "chown", "close", "copyFile", "cp",
    "glob", "exists", "fchmod", "fchown", "fdatasync", "fstat", "fsync", "ftruncate",
    "futimes", "lchmod", "lchown", "link", "lstat", "lutimes", "mkdir",
    "mkdtemp", "mkdtempDisposable", "open", "opendir", "read", "readFile",
    "readdir", "readlink", "readv", "realpath", "rename", "rm", "rmdir",
    "stat", "statfs", "symlink", "truncate", "unlink", "utimes", "write",
    "writeFile", "writev",
}

func tsFSModule(name string) tsValue {
    tsFSModules.Lock()
    defer tsFSModules.Unlock()
    if tsFSModules.fs == nil {
        fs, promises := tsNewObject(), tsNewObject()
        tsFSModules.fs, tsFSModules.promises = fs, promises
        for _,name:=range []string{"Stats","Dirent","Dir","ReadStream","WriteStream"}{fs.set(name,tsClassValue(tsNativeClass(name)))};fs.set("FileReadStream",fs.get("ReadStream"));fs.set("FileWriteStream",fs.get("WriteStream"))
        constants := tsFSConstants()
        fs.set("constants", tsObjectValue(constants))
        promises.set("constants", tsObjectValue(constants))
        for _, key := range []string{"F_OK", "R_OK", "W_OK", "X_OK"} { fs.set(key, constants.get(key)) }
        for _, operation := range tsFSOperations {
            op := operation
            fs.set(op+"Sync", tsFunctionValue(tsFunc(func(loop *tsLoop, args ...tsValue) tsValue {
                return tsFSInvoke(loop, op, "sync", args)
            })))
            if op != "mkdtempDisposable" {
                fs.set(op, tsFunctionValue(tsFunc(func(loop *tsLoop, args ...tsValue) tsValue {
                    return tsFSInvoke(loop, op, "callback", args)
                })))
            }
            switch op {
            case "close", "exists", "fchmod", "fchown", "fdatasync", "fstat", "fsync", "ftruncate", "futimes", "read", "readv", "write", "writev":
                continue
            }
            promises.set(op, tsFunctionValue(tsFunc(func(loop *tsLoop, args ...tsValue) tsValue {
                return tsFSInvoke(loop, op, "promise", args)
            })))
        }
        promises.set("glob",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsFSGlobIterator(loop,args)})))
        for _,name:=range []string{"realpath","realpathSync"}{value:=fs.get(name);function:=(*tsFunction)(value.ref);function.properties=tsNewObject();function.properties.set("native",value)}
        tsFSInstallWatchers(fs,promises)
        tsFSInstallStreams(fs)
        fs.set("openAsBlob",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsFSOpenAsBlob(loop,args)})))
        fs.set("promises", tsObjectValue(promises))
    }
    if name == "fs/promises" { return tsObjectValue(tsFSModules.promises) }
    return tsObjectValue(tsFSModules.fs)
}

func tsFSProperty(value tsValue, name string) tsValue {
    if tsNullish(value) { return tsU }
    return tsGet(value, tsStringReference(tsStringUTF8(name)))
}

func tsFSOption(value tsValue, name string, fallback tsValue) tsValue {
    item := tsFSProperty(value, name)
    if tsIsUndefined(item) { return fallback }
    return item
}

func tsFSException(name, code, message string) tsValue {
    return tsErrorValue(&tsRuntimeError{name:name, message:message, fields:map[string]tsValue{
        "code":tsStringReference(tsStringUTF8(code)),
    }})
}

func tsFSInvalid(name, expected string) {
    panic(tsThrown{tsFSException("TypeError", "ERR_INVALID_ARG_TYPE", "The " + name + " argument must be " + expected)})
}

func tsFSInteger(value tsValue, name string, fallback, minimum int64) int64 {
    if tsIsUndefined(value) { return fallback }
    if !tsIsNumeric(value) { tsFSInvalid(name, "a number") }
    number := tsNumber(value)
    if math.IsNaN(number) || math.IsInf(number, 0) || math.Trunc(number) != number || number < float64(minimum) || number >= 9223372036854775808.0 {
        panic(tsThrown{tsFSException("RangeError", "ERR_OUT_OF_RANGE", name + " is out of range")})
    }
    return int64(number)
}

func tsFSPath(value tsValue) string {
    var path string
    switch value.kind {
    case tsStringKind:
        path = (*tsString)(value.ref).String()
    case tsTypedArrayKind:
        array := (*tsTypedArray)(value.ref)
        if !array.nodeBuffer { tsFSInvalid("path", "a string, Buffer, or file URL") }
        path = string(tsFSBytes(value))
    case tsObjectKind:
        protocol := tsFSProperty(value, "protocol")
        if tsText(protocol) != "file:" { tsFSInvalid("path", "a string, Buffer, or file URL") }
        host := tsText(tsFSProperty(value, "hostname"))
        if host != "" && host != "localhost" { panic(tsThrown{tsFSException("TypeError", "ERR_INVALID_FILE_URL_HOST", "File URL host must be empty or localhost")}) }
        encoded := tsText(tsFSProperty(value, "pathname"))
        if strings.Contains(strings.ToLower(encoded), "%2f") { panic(tsThrown{tsFSException("TypeError", "ERR_INVALID_FILE_URL_PATH", "File URL path must not contain encoded slashes")}) }
        var err error
        path, err = url.PathUnescape(encoded)
        if err != nil { panic(tsThrown{tsFSException("TypeError", "ERR_INVALID_FILE_URL_PATH", err.Error())}) }
    default:
        tsFSInvalid("path", "a string, Buffer, or file URL")
    }
    if strings.ContainsRune(path, 0) { panic(tsThrown{tsFSException("TypeError", "ERR_INVALID_ARG_VALUE", "Path must not contain null bytes")}) }
    return path
}

func tsFSError(err error, operation, path, dest string) tsValue {
    code, errno := "UNKNOWN", float64(-1)
    var number syscall.Errno
    if errors.As(err, &number) {
        errno = -float64(number)
        for name, value := range map[string]syscall.Errno{
            "ENOENT":syscall.ENOENT, "EACCES":syscall.EACCES, "EPERM":syscall.EPERM,
            "EEXIST":syscall.EEXIST, "ENOTDIR":syscall.ENOTDIR, "EISDIR":syscall.EISDIR,
            "ENOTEMPTY":syscall.ENOTEMPTY, "EBADF":syscall.EBADF, "EINVAL":syscall.EINVAL,
            "EXDEV":syscall.EXDEV, "ELOOP":syscall.ELOOP, "EMFILE":syscall.EMFILE,
            "ENFILE":syscall.ENFILE, "ENOSPC":syscall.ENOSPC, "EROFS":syscall.EROFS,
            "ENAMETOOLONG":syscall.ENAMETOOLONG, "EIO":syscall.EIO, "ENOSYS":syscall.ENOSYS,
            "ENOTSUP":syscall.ENOTSUP, "EFBIG":syscall.EFBIG, "EAGAIN":syscall.EAGAIN,
        } { if number == value { code = name; break } }
    } else if errors.Is(err, os.ErrNotExist) { code = "ENOENT" } else if errors.Is(err, os.ErrExist) { code = "EEXIST" } else if errors.Is(err, os.ErrPermission) { code = "EACCES" }
    var pe *os.PathError
    if errors.As(err, &pe) { operation = pe.Op }
    message := code + ": " + err.Error()
    fields := map[string]tsValue{
        "code":tsStringReference(tsStringUTF8(code)), "errno":tsNumberValue(errno),
        "syscall":tsStringReference(tsStringUTF8(operation)),
    }
    if path != "" { fields["path"] = tsStringReference(tsStringUTF8(path)) }
    if dest != "" { fields["dest"] = tsStringReference(tsStringUTF8(dest)) }
    return tsErrorValue(&tsRuntimeError{name:"Error", message:message, fields:fields})
}

func tsFSConstants() *tsObject {
    out := tsNewObject()
    values := map[string]int{
        "F_OK":0, "R_OK":4, "W_OK":2, "X_OK":1,
        "COPYFILE_EXCL":1, "COPYFILE_FICLONE":2, "COPYFILE_FICLONE_FORCE":4,
        "O_RDONLY":os.O_RDONLY, "O_WRONLY":os.O_WRONLY, "O_RDWR":os.O_RDWR,
        "O_CREAT":os.O_CREATE, "O_EXCL":os.O_EXCL, "O_TRUNC":os.O_TRUNC,
        "O_APPEND":os.O_APPEND, "O_SYNC":os.O_SYNC,
        "O_DIRECT":unix.O_DIRECT,"O_DIRECTORY":unix.O_DIRECTORY,"O_NOFOLLOW":unix.O_NOFOLLOW,"O_NONBLOCK":unix.O_NONBLOCK,"O_DSYNC":unix.O_DSYNC,"O_NOCTTY":unix.O_NOCTTY,
        "UV_DIRENT_UNKNOWN":0,"UV_DIRENT_FILE":1,"UV_DIRENT_DIR":2,"UV_DIRENT_LINK":3,"UV_DIRENT_FIFO":4,"UV_DIRENT_SOCKET":5,"UV_DIRENT_CHAR":6,"UV_DIRENT_BLOCK":7,
        "S_IFMT":syscall.S_IFMT, "S_IFREG":syscall.S_IFREG, "S_IFDIR":syscall.S_IFDIR,
        "S_IFLNK":syscall.S_IFLNK, "S_IFBLK":syscall.S_IFBLK, "S_IFCHR":syscall.S_IFCHR,
        "S_IFIFO":syscall.S_IFIFO, "S_IFSOCK":syscall.S_IFSOCK,
        "S_IRUSR":0400, "S_IWUSR":0200, "S_IXUSR":0100,
        "S_IRGRP":0040, "S_IWGRP":0020, "S_IXGRP":0010,
        "S_IROTH":0004, "S_IWOTH":0002, "S_IXOTH":0001,
    }
    for key, value := range values { out.set(key, tsNumberValue(float64(value))) }
    return out
}

// Snapshot configuration on the loop; native workers never evaluate user getters.
func tsFSPrepare(loop *tsLoop, args []tsValue) []tsValue {
    result := append([]tsValue{}, args...)
    for i, value := range result {
        if value.kind != tsObjectKind && value.kind != tsInstanceKind { continue }
        if value.kind==tsObjectKind&&(*tsObject)(value.ref).nativeClass!=nil{continue}
        if !tsIsUndefined(tsFSProperty(value, "fd")) || !tsIsUndefined(tsFSProperty(value, "protocol")) { continue }
        options := tsNewObject()
        for _, key := range []string{"encoding", "flag", "mode", "recursive", "force", "errorOnExist", "dereference", "preserveTimestamps", "verbatimSymlinks", "filter", "withFileTypes", "bigint", "throwIfNoEntry", "offset", "length", "position", "buffer", "flush", "maxRetries", "retryDelay", "bufferSize", "cwd", "exclude", "signal", "persistent", "interval", "autoClose", "emitClose", "start", "end", "highWaterMark", "flags"} {
            item := tsFSProperty(value, key)
            if !tsIsUndefined(item) { options.set(key, item) }
        }
        result[i] = tsObjectValue(options)
    }
    return result
}

func tsFSInvoke(loop *tsLoop, op, mode string, args []tsValue) tsValue {
    callback := tsU
    if mode == "callback" {
        if len(args) == 0 || args[len(args)-1].kind != tsFunctionKind { tsFSInvalid("callback", "a function") }
        callback = args[len(args)-1]
        args = args[:len(args)-1]
    }
    var snapshot []tsValue
    if mode == "promise" {
        var failure tsValue
        func(){ defer func(){if v:=recover();v!=nil{failure=tsUnwrap(v)}}(); snapshot=tsFSPrepare(loop,args) }()
        if !tsIsUndefined(failure) { return tsPromiseValue(loop.resolved(failure,true)) }
    } else { snapshot = tsFSPrepare(loop,args) }
    work := func() tsResult { return tsFSExecute(loop, op, snapshot, mode != "sync") }
    if mode == "sync" {
        result := work()
        if result.rejected { panic(tsThrown{result.value}) }
        if op == "read" || op == "write" || op == "readv" || op == "writev" {
            key := "bytesRead"; if strings.HasPrefix(op,"write") {key="bytesWritten"}
            return tsFSProperty(result.value,key)
        }
        return result.value
    }
    if mode == "promise" {
        promise := loop.promise()
        loop.submit(work,func(result tsResult){
            if op=="open"&&!result.rejected { result.value=tsFSFileHandle(loop,result.value) }
            promise.settle(result)
        })
        return tsPromiseValue(promise)
    }
    loop.submit(work,func(result tsResult){
        if op=="exists" { tsCall(loop,callback,result.value);return }
        if result.rejected { tsCall(loop,callback,result.value);return }
        switch op {
        case "read","write","readv","writev":
            key, buffer := "bytesRead", "buffer"
            if strings.HasPrefix(op,"write") { key="bytesWritten" }
            if strings.HasSuffix(op,"v") {buffer="buffers"}
            tsCall(loop,callback,tsNull,tsFSProperty(result.value,key),tsFSProperty(result.value,buffer))
        case "access","appendFile","chmod","chown","close","copyFile","cp","fchmod","fchown","fdatasync","fsync","ftruncate","futimes","lchmod","lchown","link","lutimes","rename","rm","rmdir","symlink","truncate","unlink","utimes","writeFile":
            tsCall(loop,callback,tsNull)
        default:
            tsCall(loop,callback,tsNull,result.value)
        }
    })
    return tsU
}
`
