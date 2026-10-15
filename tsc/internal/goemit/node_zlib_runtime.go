package goemit

// nodeZlibRuntime maps Node's zlib call and stream APIs onto Go's compression
// packages. Brotli and Zstandard retain the API but report an unsupported-codec
// error because the Go runtime has no matching standard-library codecs.
const nodeZlibRuntime = `
type tsNodeZlibState struct { mode string; options tsValue; input []byte; stream tsValue; level, strategy int; closed bool }
var tsNodeZlibSingleton struct { sync.Once; module tsValue; states map[uintptr]*tsNodeZlibState; transform tsValue }

func tsNodeZlibKey(name string) tsValue { return tsStringReference(tsStringUTF8(name)) }
func tsNodeZlibOption(loop *tsLoop,options tsValue,name string,fallback tsValue) tsValue { value:=tsGet(loop,options,tsNodeZlibKey(name));if tsIsUndefined(value){return fallback};return value }
func tsNodeZlibError(code,message string) tsValue { return tsErrorValue(&tsRuntimeError{name:"Error",message:message,fields:map[string]tsValue{"code":tsNodeZlibKey(code)}}) }
func tsNodeZlibInput(loop *tsLoop,value tsValue) []byte {
	if value.kind==tsStringKind { return []byte(tsText(value)) }
	if value.kind==tsArrayBufferKind||value.kind==tsTypedArrayKind||value.kind==tsDataViewKind { return append([]byte(nil),tsFSBytes(value)...) }
	if value.kind==tsArrayKind||value.kind==tsObjectKind { return append([]byte(nil),tsFSBytes(tsNodeBufferFrom(loop,[]tsValue{value}))...) }
	panic(tsThrown{tsNodeZlibError("ERR_INVALID_ARG_TYPE","The \"buffer\" argument must be of type string or an instance of Buffer, TypedArray, DataView, or ArrayBuffer")})
}
func tsNodeZlibLevel(loop *tsLoop,options tsValue) (int,error) {
	value:=tsNodeZlibOption(loop,options,"level",tsNumberValue(-1));if !tsIsNumeric(value){return 0,fmt.Errorf("options.level must be a number")};level:=int(tsNumber(value));if level < -1 || level > 9{return 0,fmt.Errorf("options.level is out of range")};return level,nil
}
func tsNodeZlibDictionary(loop *tsLoop,options tsValue) []byte {
	value:=tsGet(loop,options,tsNodeZlibKey("dictionary"));if tsNullish(value){return nil};switch value.kind{case tsTypedArrayKind,tsDataViewKind,tsArrayBufferKind:return append([]byte(nil),tsFSBytes(value)...);default:panic(tsThrown{tsNodeZlibError("ERR_INVALID_ARG_TYPE","The \"options.dictionary\" argument must be a Buffer, TypedArray, DataView, or ArrayBuffer")})}
}
func tsNodeZlibProcess(loop *tsLoop,mode string,input []byte,options tsValue) ([]byte,error) {
	if mode=="brotliCompress"||mode=="brotliDecompress"{return nil,fmt.Errorf("Brotli codec is not available in the Go standard library")}
	if mode=="zstdCompress"||mode=="zstdDecompress"{return nil,fmt.Errorf("Zstandard codec is not available in the Go standard library")}
	level,err:=tsNodeZlibLevel(loop,options);if err!=nil{return nil,err};dictionary:=tsNodeZlibDictionary(loop,options)
	var output bytes.Buffer
	switch mode {
	case "deflate","gzip","deflateRaw":
		var writer io.WriteCloser
		switch mode {
		case "deflate":if len(dictionary)>0{writer,err=zlib.NewWriterLevelDict(&output,level,dictionary)}else{writer,err=zlib.NewWriterLevel(&output,level)}
		case "gzip":writer,err=gzip.NewWriterLevel(&output,level)
		case "deflateRaw":writer,err=flate.NewWriterDict(&output,level,dictionary);if len(dictionary)==0{writer,err=flate.NewWriter(&output,level)}
		}
		if err!=nil{return nil,err};if _,err=writer.Write(input);err!=nil{_ = writer.Close();return nil,err};if err=writer.Close();err!=nil{return nil,err}
	case "inflate","gunzip","inflateRaw","unzip":
		var reader io.ReadCloser
		switch mode {
		case "inflate":if len(dictionary)>0{reader,err=zlib.NewReaderDict(bytes.NewReader(input),dictionary)}else{reader,err=zlib.NewReader(bytes.NewReader(input))}
		case "gunzip":reader,err=gzip.NewReader(bytes.NewReader(input))
		case "inflateRaw":if len(dictionary)>0{reader=flate.NewReaderDict(bytes.NewReader(input),dictionary)}else{reader=flate.NewReader(bytes.NewReader(input))}
		case "unzip":if len(input)>=2&&input[0]==0x1f&&input[1]==0x8b{reader,err=gzip.NewReader(bytes.NewReader(input))}else if len(dictionary)>0{reader,err=zlib.NewReaderDict(bytes.NewReader(input),dictionary)}else{reader,err=zlib.NewReader(bytes.NewReader(input))}
		}
		if err!=nil{return nil,err};defer reader.Close();maxOutput:=int64(1<<31-1);maxValue:=tsGet(loop,options,tsNodeZlibKey("maxOutputLength"));if !tsIsUndefined(maxValue){n:=tsNumber(maxValue);if n<0||math.IsNaN(n){return nil,fmt.Errorf("maxOutputLength is out of range")};maxOutput=int64(n)};limited:=io.LimitReader(reader,maxOutput+1);var data []byte;data,err=io.ReadAll(limited);if err!=nil{return nil,err};if int64(len(data))>maxOutput{return nil,fmt.Errorf("output exceeds maxOutputLength")};return data,nil
	default:return nil,fmt.Errorf("Unknown compression method %s",mode)
	}
	maxValue:=tsGet(loop,options,tsNodeZlibKey("maxOutputLength"));if !tsIsUndefined(maxValue)&&int64(output.Len())>int64(tsNumber(maxValue)){return nil,fmt.Errorf("output exceeds maxOutputLength")};return output.Bytes(),nil
}
func tsNodeZlibErrorValue(err error) tsValue {code:="Z_DATA_ERROR";message:=err.Error();if strings.Contains(message,"maxOutputLength"){code="ERR_BUFFER_TOO_LARGE"};return tsErrorValue(&tsRuntimeError{name:"Error",message:message,fields:map[string]tsValue{"code":tsNodeZlibKey(code),"errno":tsNumberValue(-3)}})}
func tsNodeZlibConvenience(loop *tsLoop,mode string,synchronous bool,args []tsValue) tsValue {
	input:=tsNodeZlibInput(loop,tsArg(args,0));options:=tsU;var callback tsValue
	if synchronous { options=tsArg(args,1) } else {
		if len(args)==0||args[len(args)-1].kind!=tsFunctionKind{panic(tsThrown{tsNodeZlibError("ERR_INVALID_ARG_TYPE","The \"callback\" argument must be of type function")})};callback=args[len(args)-1]
		if len(args)>2 {options=args[1]} else if len(args)==2&&args[1].kind!=tsFunctionKind {options=args[1]}
	}
	if !tsNullish(options)&&options.kind!=tsObjectKind&&options.kind!=tsInstanceKind{panic(tsThrown{tsNodeZlibError("ERR_INVALID_ARG_TYPE","The \"options\" argument must be an object")})}
	output,err:=tsNodeZlibProcess(loop,mode,input,options)
	if synchronous {if err!=nil{panic(tsThrown{tsNodeZlibErrorValue(err)})};buffer:=tsNodeBuffer(output);if tsTruthy(tsGet(loop,options,tsNodeZlibKey("info"))){result:=tsNewObject();result.set("buffer",buffer);result.set("engine",tsNull);return tsObjectValue(result)};return buffer}
	loop.ready=append(loop.ready,func(){if err!=nil{tsCall(loop,callback,tsNodeZlibErrorValue(err));return};buffer:=tsNodeBuffer(output);if tsTruthy(tsGet(loop,options,tsNodeZlibKey("info"))){result:=tsNewObject();result.set("buffer",buffer);result.set("engine",tsNull);buffer=tsObjectValue(result)};tsCall(loop,callback,tsNull,buffer)})
	return tsU
}
func tsNodeZlibStream(loop *tsLoop,mode string,options tsValue) tsValue {
	if tsNullish(options){options=tsObjectValue(tsNewObject())};if options.kind!=tsObjectKind&&options.kind!=tsInstanceKind{panic(tsThrown{tsNodeZlibError("ERR_INVALID_ARG_TYPE","The \"options\" argument must be an object")})};level,err:=tsNodeZlibLevel(loop,options);if err!=nil{panic(tsThrown{tsNodeZlibError("ERR_OUT_OF_RANGE",err.Error())})}
	state:=&tsNodeZlibState{mode:mode,options:options,level:level};streamOptions:=tsNewObject();tsObjectSpread(streamOptions,options)
	streamOptions.set("transform",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{chunk,callback:=tsArg(args,0),tsArg(args,2);if chunk.kind==tsStringKind{state.input=append(state.input,[]byte(tsText(chunk))...)}else{state.input=append(state.input,tsFSBytes(chunk)...)};if state.stream.kind==tsObjectKind{(*tsObject)(state.stream.ref).set("bytesWritten",tsNumberValue(float64(len(state.input))))};return tsCall(loop,callback,tsNull)})))
	streamOptions.set("flush",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{callback:=tsArg(args,0);data,err:=tsNodeZlibProcess(loop,state.mode,state.input,state.options);if err!=nil{return tsCall(loop,callback,tsNodeZlibErrorValue(err))};buffer:=tsNodeBuffer(data);if state.stream.kind==tsObjectKind{(*tsObject)(state.stream.ref).set("bytesRead",tsNumberValue(float64(len(data))))};return tsCall(loop,callback,tsNull,buffer)})))
	module:=tsInternalStreamTransformModule(loop);constructor:=tsGet(loop,module,tsNodeZlibKey("Transform"));stream:=tsConstruct(constructor,tsObjectValue(streamOptions));state.stream=stream;tsNodeZlibSingleton.states[uintptr((*tsObject)(stream.ref))]=state
	object:=(*tsObject)(stream.ref);object.set("bytesWritten",tsNumberValue(0));object.set("bytesRead",tsNumberValue(0));object.set("_mode",tsNodeZlibKey(mode))
	object.set("reset",tsNamedNativeMethod("reset",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{if current:=tsNodeZlibSingleton.states[uintptr(receiver.ref)];current!=nil{current.input=nil;receiverObj:=(*tsObject)(receiver.ref);receiverObj.set("bytesWritten",tsNumberValue(0));receiverObj.set("bytesRead",tsNumberValue(0))};return tsU}))
	object.set("close",tsNamedNativeMethod("close",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{if current:=tsNodeZlibSingleton.states[uintptr(receiver.ref)];current!=nil{current.closed=true};if callback:=tsArg(args,0);callback.kind==tsFunctionKind{tsNodeWritableDefer(loop,func(){tsCall(loop,callback)})};tsNodeClassicMethod(loop,receiver,"destroy",tsU);return tsU}))
	object.set("flush",tsNamedNativeMethod("flush",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{callback:=tsArg(args,1);if tsArg(args,0).kind==tsFunctionKind{callback=tsArg(args,0)};if callback.kind==tsFunctionKind{tsNodeWritableDefer(loop,func(){tsCall(loop,callback,tsNull)})};return tsU}))
	object.set("params",tsNamedNativeMethod("params",3,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{if current:=tsNodeZlibSingleton.states[uintptr(receiver.ref)];current!=nil{current.level=int(tsNumber(tsArg(args,0)));current.strategy=int(tsNumber(tsArg(args,1)))};if callback:=tsArg(args,2);callback.kind==tsFunctionKind{tsNodeWritableDefer(loop,func(){tsCall(loop,callback)})};return tsU}))
	return stream
}
func tsNodeZlibCtor(mode,name string) tsValue {
	function:=tsFunc(func(_ *tsLoop,_ ...tsValue)tsValue{tsPropertyFailure("Class constructor "+name+" cannot be invoked without 'new'");return tsU});function.name=name;function.length=1;function.constructible=true;function.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeZlibStream(loop,mode,tsArg(args,0))};function.properties=tsNewObject();return tsFunctionValue(function)
}
func tsNodeZlibCRC(loop *tsLoop,args []tsValue) tsValue {data:=tsNodeZlibInput(loop,tsArg(args,0));seed:=uint32(0);if !tsIsUndefined(tsArg(args,1)){seed=uint32(tsNumber(tsArg(args,1)))};return tsNumberValue(float64(crc32.Update(seed,crc32.IEEETable,data)))}
func tsNodeZlibConstants() *tsObject {
	values:=map[string]float64{"Z_NO_FLUSH":0,"Z_PARTIAL_FLUSH":1,"Z_SYNC_FLUSH":2,"Z_FULL_FLUSH":3,"Z_FINISH":4,"Z_BLOCK":5,"Z_TREES":6,"Z_OK":0,"Z_STREAM_END":1,"Z_NEED_DICT":2,"Z_ERRNO":-1,"Z_STREAM_ERROR":-2,"Z_DATA_ERROR":-3,"Z_MEM_ERROR":-4,"Z_BUF_ERROR":-5,"Z_VERSION_ERROR":-6,"Z_NO_COMPRESSION":0,"Z_BEST_SPEED":1,"Z_BEST_COMPRESSION":9,"Z_DEFAULT_COMPRESSION":-1,"Z_FILTERED":1,"Z_HUFFMAN_ONLY":2,"Z_RLE":3,"Z_FIXED":4,"Z_DEFAULT_STRATEGY":0,"ZLIB_VERNUM":0x12d0,"DEFLATE":1,"INFLATE":2,"GZIP":3,"GUNZIP":4,"DEFLATERAW":5,"INFLATERAW":6,"UNZIP":7,"BROTLI_DECODE":8,"BROTLI_ENCODE":9,"ZSTD_COMPRESS":10,"ZSTD_DECOMPRESS":11,"Z_MIN_WINDOWBITS":8,"Z_MAX_WINDOWBITS":15,"Z_DEFAULT_WINDOWBITS":15,"Z_MIN_CHUNK":64,"Z_MAX_CHUNK":math.MaxUint32,"Z_DEFAULT_CHUNK":16384,"Z_MIN_MEMLEVEL":1,"Z_MAX_MEMLEVEL":9,"Z_DEFAULT_MEMLEVEL":8,"Z_MIN_LEVEL":-1,"Z_MAX_LEVEL":9,"Z_DEFAULT_LEVEL":-1,"BROTLI_OPERATION_PROCESS":0,"BROTLI_OPERATION_FLUSH":1,"BROTLI_OPERATION_FINISH":2,"BROTLI_OPERATION_EMIT_METADATA":3,"ZSTD_e_continue":0,"ZSTD_e_flush":1,"ZSTD_e_end":2}
	out:=tsNewObject();for name,value:=range values{out.set(name,tsNumberValue(value))};return out
}
func tsNodeZlibModule(loop *tsLoop) tsValue {
	tsNodeZlibSingleton.Do(func(){tsNodeZlibSingleton.states=map[uintptr]*tsNodeZlibState{};tsNodeZlibSingleton.transform=tsGet(loop,tsInternalStreamTransformModule(loop),tsNodeZlibKey("Transform"));module:=tsNewObject();classes:=map[string]string{"Deflate":"deflate","Inflate":"inflate","Gzip":"gzip","Gunzip":"gunzip","DeflateRaw":"deflateRaw","InflateRaw":"inflateRaw","Unzip":"unzip","BrotliCompress":"brotliCompress","BrotliDecompress":"brotliDecompress","ZstdCompress":"zstdCompress","ZstdDecompress":"zstdDecompress"};for name,mode:=range classes{module.set(name,tsNodeZlibCtor(mode,name))}
		for _,method:=range []string{"deflate","gzip","deflateRaw","inflate","gunzip","inflateRaw","unzip","brotliCompress","brotliDecompress","zstdCompress","zstdDecompress"}{mode:=method;module.set(method,tsNamedNativeMethod(method,3,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeZlibConvenience(loop,mode,false,args)}));module.set(method+"Sync",tsNamedNativeMethod(method+"Sync",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeZlibConvenience(loop,mode,true,args)}))}
		for name,mode:=range classes{ctor:=tsGet(loop,tsObjectValue(module),tsNodeZlibKey(name));module.set("create"+name,tsNamedNativeMethod("create"+name,1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeZlibStream(loop,mode,tsArg(args,0))}));_ = ctor}
		constants:=tsNodeZlibConstants();codes:=tsNewObject();for _,name:=range []string{"Z_OK","Z_STREAM_END","Z_NEED_DICT","Z_ERRNO","Z_STREAM_ERROR","Z_DATA_ERROR","Z_MEM_ERROR","Z_BUF_ERROR","Z_VERSION_ERROR"}{value:=constants.get(name);codes.set(name,value);codes.set(strconv.Itoa(int(tsNumber(value))),tsNodeZlibKey(name))};module.set("constants",tsObjectValue(constants));module.set("codes",tsObjectValue(codes));module.set("crc32",tsNamedNativeMethod("crc32",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeZlibCRC(loop,args)}))
		for _,name:=range []string{"Z_NO_FLUSH","Z_PARTIAL_FLUSH","Z_SYNC_FLUSH","Z_FULL_FLUSH","Z_FINISH","Z_BLOCK","Z_OK","Z_STREAM_END","Z_NEED_DICT","Z_ERRNO","Z_STREAM_ERROR","Z_DATA_ERROR","Z_MEM_ERROR","Z_BUF_ERROR","Z_VERSION_ERROR","Z_NO_COMPRESSION","Z_BEST_SPEED","Z_BEST_COMPRESSION","Z_DEFAULT_COMPRESSION","Z_FILTERED","Z_HUFFMAN_ONLY","Z_RLE","Z_FIXED","Z_DEFAULT_STRATEGY","ZLIB_VERNUM","DEFLATE","INFLATE","GZIP","GUNZIP","DEFLATERAW","INFLATERAW","UNZIP","Z_MIN_WINDOWBITS","Z_MAX_WINDOWBITS","Z_DEFAULT_WINDOWBITS","Z_MIN_CHUNK","Z_MAX_CHUNK","Z_DEFAULT_CHUNK","Z_MIN_MEMLEVEL","Z_MAX_MEMLEVEL","Z_DEFAULT_MEMLEVEL","Z_MIN_LEVEL","Z_MAX_LEVEL","Z_DEFAULT_LEVEL"}{module.set(name,constants.get(name))}
		module.set("default",tsObjectValue(module));tsNodeZlibSingleton.module=tsObjectValue(module)
	});return tsNodeZlibSingleton.module
}
`
