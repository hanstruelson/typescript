package goemit

// Native class methods take an explicit JavaScript receiver. Property extraction
// leaves them unbound; method calls, .call, .apply and .bind pass the receiver.
const nativeClassRuntime = `
func tsCallReceiver(loop *tsLoop,value,receiver tsValue,args ...tsValue)tsValue{
    if value.kind==tsFunctionKind{function:=(*tsFunction)(value.ref);if function.receiverCall!=nil{return function.receiverCall(loop,receiver,args...)}}
    return tsCall(loop,value,args...)
}
func tsFunctionGet(loop *tsLoop,value tsValue,name string)tsValue{
    function:=(*tsFunction)(value.ref)
    if object:=function.readProperties();object!=nil{if tsOwnObjectProperty(object,name){if descriptor:=object.descriptors[name];descriptor!=nil&&descriptor.accessor{if tsIsUndefined(descriptor.getter){return tsU};return tsCallReceiver(loop,descriptor.getter,value)};return object.values[name]}}
    if function.metadataRef.Load()!=nil{return tsPrototypeLookup(loop,value,tsPrototypeOf(value),name)}
    if name=="prototype"{return tsFunctionPrototype(value)};if name=="name"{return tsStringReference(tsStringUTF8(function.name))};if name=="length"{return tsNumberValue(float64(function.length))}
    return tsPrototypeLookup(loop,value,tsPrototypeOf(value),name)
}
func tsNamedNativeMethod(name string,length int,call func(*tsLoop,tsValue,...tsValue)tsValue)tsValue{function:=&tsFunction{receiverCall:call,properties:tsNewObject()};function.properties.set("name",tsStringReference(tsStringUTF8(name)));function.properties.set("length",tsNumberValue(float64(length)));function.properties.descriptors=map[string]*tsDescriptor{"name":{configurable:true},"length":{configurable:true}};return tsFunctionValue(function)}
func tsNativeMethod(call func(*tsLoop,tsValue,...tsValue)tsValue)tsValue{return tsFunctionValue(&tsFunction{receiverCall:call})}
var tsNativeClasses struct {sync.Once;classes map[string]*tsClass}
func tsNativeClass(name string)*tsClass{
    tsNativeClasses.Do(func(){
        tsNativeClasses.classes=map[string]*tsClass{}
        for _,name:=range []string{"Stats","BigIntStats","StatFs","BigIntStatFs","Dirent","Dir","FSWatcher","StatWatcher","ReadStream","WriteStream","Buffer","Date","AbortSignal","AbortController","FileHandle","Blob","ReadableStream","ReadableStreamDefaultReader","ReadableStreamBYOBReader"}{
            class:=&tsClass{builtin:name,static:tsInstanceValue(tsNewProperties()),nativePrototype:tsNewObject()};tsNativeClasses.classes[name]=class;class.nativePrototype.set("constructor",tsClassValue(class));properties:=tsInstanceProperties(class.static);properties.define("prototype");properties.extra["prototype"]=tsObjectValue(class.nativePrototype)
            class.construct=func(loop *tsLoop,args ...tsValue)tsValue{
                switch name{
                case "Blob":return tsBlobConstruct(loop,args)
                case "Date":return tsDateConstruct(args)
                case "Stats":return tsFSConstructStats(args)
                case "Dirent":return tsFSConstructDirent(args)
                case "ReadStream","WriteStream":return tsFSCreateStream(loop,tsArg(args,0),tsArg(args,1),name=="ReadStream")
                case "Buffer":module:=tsNodeBufferModule();method:="from";if tsIsNumeric(tsArg(args,0)){method="alloc"};return tsCall(loop,tsFSProperty(module,method),args...)
                case "AbortController":return tsAbortController(loop)
                default:panic(tsThrown{tsFSException("TypeError","ERR_ILLEGAL_CONSTRUCTOR","Illegal constructor")})
                }
            }
        }
        tsInstallDateClass(tsNativeClasses.classes["Date"])
        for _,name:=range []string{"Stats","BigIntStats","Dirent"}{
            class:=tsNativeClasses.classes[name]
            for _,predicate:=range []string{"isFile","isDirectory","isSymbolicLink","isBlockDevice","isCharacterDevice","isFIFO","isSocket"}{
                method:=predicate;class.nativePrototype.set(method,tsNativeMethod(func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{
                    if receiver.kind!=tsObjectKind{tsFSInvalid("this","a filesystem metadata object")};object:=(*tsObject)(receiver.ref);mode:=object.nativeMode;if class.builtin=="Dirent"&&object.nativeClass!=class{mode=os.ModeIrregular}
                    if class.builtin=="Stats"||class.builtin=="BigIntStats"{value:=object.get("mode");bits:=uint32(0);if value.kind==tsBigIntKind{bits=uint32((*big.Int)(value.ref).Uint64())}else{bits=tsNumberUint32(tsNumber(value))};switch bits&syscall.S_IFMT{case syscall.S_IFDIR:mode=os.ModeDir;case syscall.S_IFLNK:mode=os.ModeSymlink;case syscall.S_IFBLK:mode=os.ModeDevice;case syscall.S_IFCHR:mode=os.ModeDevice|os.ModeCharDevice;case syscall.S_IFIFO:mode=os.ModeNamedPipe;case syscall.S_IFSOCK:mode=os.ModeSocket;case syscall.S_IFREG:mode=0;default:mode=os.ModeIrregular}}
                    answer:=false;switch method{case "isFile":answer=mode.IsRegular();case "isDirectory":answer=mode.IsDir();case "isSymbolicLink":answer=mode&os.ModeSymlink!=0;case "isBlockDevice":answer=mode&os.ModeDevice!=0&&mode&os.ModeCharDevice==0;case "isCharacterDevice":answer=mode&os.ModeCharDevice!=0;case "isFIFO":answer=mode&os.ModeNamedPipe!=0;case "isSocket":answer=mode&os.ModeSocket!=0};return tsBooleanValue(answer)
                }))
            }
        }
        buffer:=tsNativeClasses.classes["Buffer"]
        methods:=[]string{};for _,operation:=range []string{"read","write"}{for _,spec:=range []string{"UInt","Uint","Int","UInt16","Uint16","Int16","UInt32","Uint32","Int32","Float","Double","BigUInt64","BigUint64","BigInt64"}{for _,endian:=range []string{"LE","BE"}{methods=append(methods,operation+spec+endian)}};for _,spec:=range []string{"UInt8","Uint8","Int8"}{methods=append(methods,operation+spec)}}
        for _,name:=range methods{method:=name;buffer.nativePrototype.set(method,tsNativeMethod(func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{if receiver.kind!=tsTypedArrayKind||!(*tsTypedArray)(receiver.ref).nodeBuffer{tsFSInvalid("this","a Buffer")};return tsCall(loop,tsNodeBufferGet(receiver,method),args...)}))}
        for _,name:=range []string{"toString","toLocaleString","toJSON","equals","compare","copy","fill","write","slice","subarray","includes","indexOf","lastIndexOf","swap16","swap32","swap64"}{method:=name;buffer.nativePrototype.set(method,tsNativeMethod(func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{if receiver.kind!=tsTypedArrayKind||!(*tsTypedArray)(receiver.ref).nodeBuffer{tsFSInvalid("this","a Buffer")};return tsCall(loop,tsNodeBufferGet(receiver,method),args...)}))}
        for _,class:=range tsNativeClasses.classes{class.nativePrototype.descriptors=map[string]*tsDescriptor{};for _,key:=range class.nativePrototype.order{class.nativePrototype.descriptors[key]=&tsDescriptor{writable:true,configurable:true}}}
    })
    return tsNativeClasses.classes[name]
}
func tsFSConstructStats(args []tsValue)tsValue{
    out:=tsNewObject();out.nativeClass=tsNativeClass("Stats");for index,name:=range []string{"dev","mode","nlink","uid","gid","rdev","blksize","ino","size","blocks","atimeMs","mtimeMs","ctimeMs","birthtimeMs"}{out.set(name,tsArg(args,index))};return tsObjectValue(out)
}
func tsFSConstructDirent(args []tsValue)tsValue{
    out:=tsNewObject();out.nativeClass=tsNativeClass("Dirent");out.set("name",tsArg(args,0));out.set("parentPath",tsArg(args,2));out.set("path",tsArg(args,2));mode:=os.ModeIrregular;switch int(tsNumber(tsArg(args,1))){case 1:mode=0;case 2:mode=os.ModeDir;case 3:mode=os.ModeSymlink;case 4:mode=os.ModeNamedPipe;case 5:mode=os.ModeSocket;case 6:mode=os.ModeDevice|os.ModeCharDevice;case 7:mode=os.ModeDevice};out.nativeMode=mode;return tsObjectValue(out)
}
`
