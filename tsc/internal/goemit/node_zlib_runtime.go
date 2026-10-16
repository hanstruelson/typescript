package goemit

// nodeZlibRuntime maps Node's zlib call and stream APIs onto Go's compression
// packages. Brotli and Zstandard use third-party Go codecs.
const nodeZlibRuntime = `
type tsNodeZlibState struct { mode string; options tsValue; input []byte; output *bytes.Buffer; encoder io.WriteCloser; encoderFlush func() error; stream tsValue; level, strategy int; bytesWritten int64; closed bool }
var tsNodeZlibSingleton struct { sync.Once; module tsValue; states map[uintptr]*tsNodeZlibState; transform tsValue }

func tsNodeZlibKey(name string) tsValue { return tsStringReference(tsStringUTF8(name)) }
func tsNodeZlibOption(loop *tsLoop,options tsValue,name string,fallback tsValue) tsValue { value:=tsGet(loop,options,tsNodeZlibKey(name));if tsIsUndefined(value){return fallback};return value }
func tsNodeZlibError(code,message string) tsValue { return tsErrorValue(&tsRuntimeError{name:"Error",message:message,fields:map[string]tsValue{"code":tsNodeZlibKey(code)}}) }
func tsNodeZlibInput(loop *tsLoop,value tsValue,name string) []byte {
	if value.kind==tsStringKind { return []byte(tsText(value)) }
	if value.kind==tsArrayBufferKind||value.kind==tsTypedArrayKind||value.kind==tsDataViewKind { return append([]byte(nil),tsFSBytes(value)...) }
	panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_TYPE","The \""+name+"\" argument must be of type string, Buffer, TypedArray, DataView, or ArrayBuffer")})
}
func tsNodeZlibLevel(loop *tsLoop,options tsValue) (int,error) {
	value:=tsNodeZlibOption(loop,options,"level",tsNumberValue(-1));if !tsIsNumeric(value){return 0,fmt.Errorf("options.level must be a number")};number:=tsNumber(value);if math.IsNaN(number)||math.IsInf(number,0)||math.Trunc(number)!=number{return 0,fmt.Errorf("The value of \"options.level\" must be an integer")};level:=int(number);if level < -1 || level > 9{return 0,fmt.Errorf("The value of \"options.level\" is out of range")};return level,nil
}
func tsNodeZlibDictionary(loop *tsLoop,options tsValue) []byte {
	value:=tsGet(loop,options,tsNodeZlibKey("dictionary"));if tsIsUndefined(value){return nil};switch value.kind{case tsTypedArrayKind,tsDataViewKind,tsArrayBufferKind:return append([]byte(nil),tsFSBytes(value)...);default:panic(tsThrown{tsNodeZlibError("ERR_INVALID_ARG_TYPE","The \"options.dictionary\" argument must be a Buffer, TypedArray, DataView, or ArrayBuffer")})}
}
func tsNodeZlibCheckInt(loop *tsLoop,options tsValue,name string,minimum,maximum,def int)int{
	value:=tsGet(loop,options,tsNodeZlibKey(name));if tsIsUndefined(value){return def};tsInternalValidatorInt32(value,tsNodeZlibKey("options."+name),tsNumberValue(float64(minimum)),tsNumberValue(float64(maximum)));return int(tsNumber(value))
}
func tsNodeZlibParams(loop *tsLoop,options tsValue,mode string){
	params:=tsGet(loop,options,tsNodeZlibKey("params"));if tsIsUndefined(params)||tsNullish(params){return}
	maximum:=9;code:="ERR_BROTLI_INVALID_PARAM";if strings.HasPrefix(mode,"zstd"){code="ERR_ZSTD_INVALID_PARAM";if strings.HasPrefix(mode,"zstdCompress"){maximum=402}else{maximum=100}}
	seen:=map[int]bool{};for _,key:=range tsEnumerableOwnKeys(loop,params).values{text:=tsText(key);number:=tsStringNumber(tsStringKey(text));if math.IsNaN(number)||math.IsInf(number,0)||number<0||number>float64(maximum)||math.Trunc(number)!=number{panic(tsThrown{tsNodeZlibError(code,"Invalid compression parameter: "+text)})};index:=int(number);if seen[index]{panic(tsThrown{tsNodeZlibError(code,"Invalid compression parameter: "+text)})};seen[index]=true;value:=tsGet(loop,params,key);if !tsIsNumeric(value)&&value.kind!=tsBooleanKind{panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_TYPE","The \"options.params[key]\" argument must be of type number")})}}
}
func tsNodeZlibBrotliOptions(loop *tsLoop,options tsValue) brotli.WriterOptions {
	quality,window:=11,22
	params:=tsGet(loop,options,tsNodeZlibKey("params"))
	if !tsNullish(params){
		if value:=tsGet(loop,params,tsNodeZlibKey("1"));!tsIsUndefined(value){quality=int(tsToUint32(value))}
		if value:=tsGet(loop,params,tsNodeZlibKey("2"));!tsIsUndefined(value){window=int(tsToUint32(value))}
	}
	if quality<0||quality>11{panic(tsThrown{tsNodeZlibError("ERR_BROTLI_PARAM_SET_FAILED","Brotli quality must be between 0 and 11")})}
	if window<10||window>24{panic(tsThrown{tsNodeZlibError("ERR_BROTLI_PARAM_SET_FAILED","Brotli window must be between 10 and 24")})}
	return brotli.WriterOptions{Quality:quality,LGWin:window}
}
func tsNodeZlibValidateOptions(loop *tsLoop,mode string,options tsValue){
	if tsNullish(options){return}
	brotli:=strings.HasPrefix(mode,"brotli");zstd:=strings.HasPrefix(mode,"zstd");flushMaximum:=5;finishDefault:=4;if brotli{flushMaximum=3;finishDefault=2};if zstd{flushMaximum=2;finishDefault=2}
	tsInternalValidatorCheckRanges(tsGet(loop,options,tsNodeZlibKey("flush")),tsNodeZlibKey("options.flush"),tsNumberValue(0),tsNumberValue(float64(flushMaximum)),tsNumberValue(0))
	tsInternalValidatorCheckRanges(tsGet(loop,options,tsNodeZlibKey("finishFlush")),tsNodeZlibKey("options.finishFlush"),tsNumberValue(0),tsNumberValue(float64(flushMaximum)),tsNumberValue(float64(finishDefault)))
	tsInternalValidatorCheckRanges(tsGet(loop,options,tsNodeZlibKey("maxOutputLength")),tsNodeZlibKey("options.maxOutputLength"),tsNumberValue(1),tsNumberValue(4294967296),tsNumberValue(4294967296))
	if chunk:=tsGet(loop,options,tsNodeZlibKey("chunkSize"));!tsIsUndefined(chunk){if !tsIsNumeric(chunk){panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_TYPE","The \"options.chunkSize\" argument must be of type number")})};number:=tsNumber(chunk);if math.IsInf(number,0){panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","The value of \"options.chunkSize\" is out of range. It must be a finite number")})};if math.IsNaN(number){chunk=tsNumberValue(16384)}else if number<64{panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","The value of \"options.chunkSize\" is out of range. It must be >= 64")})}}
	if brotli||zstd{tsNodeZlibParams(loop,options,mode);if brotli&& !tsIsUndefined(tsGet(loop,options,tsNodeZlibKey("dictionary"))){_ = tsNodeZlibDictionary(loop,options)};if mode=="zstdCompress"{pledged:=tsGet(loop,options,tsNodeZlibKey("pledgedSrcSize"));if tsIsNumeric(pledged){tsInternalValidatorUint32(pledged,tsNodeZlibKey("pledgedSrcSize"),tsBooleanValue(false))}};return}
	windowDefault:=15;decompressMode:=mode=="inflate"||mode=="gunzip"||mode=="unzip";if decompressMode{windowDefault=0};windowMinimum:=8;if mode=="gzip"{windowMinimum=9};window:=tsGet(loop,options,tsNodeZlibKey("windowBits"));if decompressMode&&(tsNullish(window)||tsIsNumeric(window)&&tsNumber(window)==0){window=tsNumberValue(0)}else{window=tsInternalValidatorCheckRanges(window,tsNodeZlibKey("options.windowBits"),tsNumberValue(float64(windowMinimum)),tsNumberValue(15),tsNumberValue(float64(windowDefault)))};if !(decompressMode&&tsIsNumeric(window)&&tsNumber(window)==0){tsInternalValidatorInt32(window,tsNodeZlibKey("options.windowBits"),tsU,tsU)}
	level:=tsInternalValidatorCheckRanges(tsGet(loop,options,tsNodeZlibKey("level")),tsNodeZlibKey("options.level"),tsNumberValue(-1),tsNumberValue(9),tsNumberValue(-1));tsInternalValidatorInt32(level,tsNodeZlibKey("options.level"),tsU,tsU)
	mem:=tsInternalValidatorCheckRanges(tsGet(loop,options,tsNodeZlibKey("memLevel")),tsNodeZlibKey("options.memLevel"),tsNumberValue(1),tsNumberValue(9),tsNumberValue(8));tsInternalValidatorInt32(mem,tsNodeZlibKey("options.memLevel"),tsU,tsU)
	strategy:=tsInternalValidatorCheckRanges(tsGet(loop,options,tsNodeZlibKey("strategy")),tsNodeZlibKey("options.strategy"),tsNumberValue(0),tsNumberValue(4),tsNumberValue(0));tsInternalValidatorInt32(strategy,tsNodeZlibKey("options.strategy"),tsU,tsU)
	if dictionary:=tsGet(loop,options,tsNodeZlibKey("dictionary"));!tsIsUndefined(dictionary){_ = tsNodeZlibDictionary(loop,options)}
}
func tsNodeZlibMaxOutput(loop *tsLoop,options tsValue)int64{
	value:=tsInternalValidatorCheckRanges(tsGet(loop,options,tsNodeZlibKey("maxOutputLength")),tsNodeZlibKey("options.maxOutputLength"),tsNumberValue(1),tsNumberValue(4294967296),tsNumberValue(4294967296));return int64(tsNumber(value))
}
func tsNodeZlibChunkSize(loop *tsLoop,options tsValue)tsValue{
	value:=tsGet(loop,options,tsNodeZlibKey("chunkSize"));if tsIsUndefined(value)||tsIsNumeric(value)&&math.IsNaN(tsNumber(value)){return tsNumberValue(16384)};return value
}
func tsNodeZlibNewEncoder(loop *tsLoop,mode string,output *bytes.Buffer,level int,dictionary []byte,optionArgs ...tsValue)(io.WriteCloser,func()error,error){
	switch mode{
	case "deflate":
		if len(dictionary)>0{writer,err:=zlib.NewWriterLevelDict(output,level,dictionary);return writer,writer.Flush,err};writer,err:=zlib.NewWriterLevel(output,level);return writer,writer.Flush,err
	case "gzip":
		writer,err:=gzip.NewWriterLevel(output,level);return writer,writer.Flush,err
	case "deflateRaw":
		if len(dictionary)>0{writer,err:=flate.NewWriterDict(output,level,dictionary);return writer,writer.Flush,err};writer,err:=flate.NewWriter(output,level);return writer,writer.Flush,err
	case "brotliCompress":
		writer:=brotli.NewWriterOptions(output,tsNodeZlibBrotliOptions(loop,tsArg(optionArgs,0)));return writer,writer.Flush,nil
	case "zstdCompress":
		writer,err:=zstd.NewWriter(output,tsNodeZlibZstdEncoderOptions(nil,tsArg(optionArgs,0),dictionary)...);return writer,writer.Flush,err
	default:return nil,nil,fmt.Errorf("Unknown compression method %s",mode)
	}
}
func tsNodeZlibZstdDictionary(loop *tsLoop,options tsValue)[]byte{
	value:=tsGet(loop,options,tsNodeZlibKey("dictionary"));if !tsTruthy(value){return nil};if value.kind==tsTypedArrayKind||value.kind==tsDataViewKind{return append([]byte(nil),tsFSBytes(value)...)};return nil
}
func tsNodeZlibZstdUint32(value tsValue)uint32{
	number:=tsNumber(value);if math.IsNaN(number)||math.IsInf(number,0)||number==0{return 0};number=math.Mod(math.Trunc(number),4294967296);if number<0{number+=4294967296};return uint32(number)
}
func tsNodeZlibZstdEncoderOptions(loop *tsLoop,options tsValue,dictionary []byte)[]zstd.EOption{
	result:=[]zstd.EOption{zstd.WithEncoderConcurrency(1)};if len(dictionary)>0{result=append(result,zstd.WithEncoderDictRaw(0,dictionary))};if !tsNullish(options){params:=tsGet(loop,options,tsNodeZlibKey("params"));if !tsNullish(params){level:=tsGet(loop,params,tsNodeZlibKey("100"));if !tsIsUndefined(level){result=append(result,zstd.WithEncoderLevel(zstd.EncoderLevelFromZstd(int(int32(tsNodeZlibZstdUint32(level)))))};checksum:=tsGet(loop,params,tsNodeZlibKey("201"));if !tsIsUndefined(checksum){result=append(result,zstd.WithEncoderCRC(tsTruthy(checksum)))}}};return result
}
func tsNodeZlibZstdDecoderOptions(dictionary []byte,maxOutput int64)[]zstd.DOption{
	memory:=uint64(maxOutput);if memory<128*1024*1024{memory=128*1024*1024};result:=[]zstd.DOption{zstd.WithDecoderConcurrency(1),zstd.WithDecoderMaxMemory(memory)};if len(dictionary)>0{result=append(result,zstd.WithDecoderDictRaw(0,dictionary))};return result
}
func tsNodeZlibProcess(loop *tsLoop,mode string,input []byte,options tsValue) ([]byte,error) {
	if mode=="brotliCompress"{
		var output bytes.Buffer
		writer:=brotli.NewWriterOptions(&output,tsNodeZlibBrotliOptions(loop,options))
		if _,err:=writer.Write(input);err!=nil{_ = writer.Close();return nil,err}
		if err:=writer.Close();err!=nil{return nil,err}
		if int64(output.Len())>tsNodeZlibMaxOutput(loop,options){return nil,fmt.Errorf("output exceeds maxOutputLength")}
		return output.Bytes(),nil
	}
	if mode=="brotliDecompress"{
		maxOutput:=tsNodeZlibMaxOutput(loop,options)
		reader:=brotli.NewReader(bytes.NewReader(input))
		data,err:=io.ReadAll(io.LimitReader(reader,maxOutput+1));if err!=nil{return nil,err}
		if int64(len(data))>maxOutput{return nil,fmt.Errorf("output exceeds maxOutputLength")}
		return data,nil
	}
	if mode=="zstdCompress"{
		writer,err:=zstd.NewWriter(nil,tsNodeZlibZstdEncoderOptions(loop,options,tsNodeZlibZstdDictionary(loop,options))...);if err!=nil{return nil,err};data:=writer.EncodeAll(input,nil);writer.Close();if int64(len(data))>tsNodeZlibMaxOutput(loop,options){return nil,fmt.Errorf("output exceeds maxOutputLength")};return data,nil
	}
	if mode=="zstdDecompress"{
		maxOutput:=tsNodeZlibMaxOutput(loop,options);reader,err:=zstd.NewReader(bytes.NewReader(input),tsNodeZlibZstdDecoderOptions(tsNodeZlibZstdDictionary(loop,options),maxOutput)...);if err!=nil{return nil,err};defer reader.Close();data,err:=io.ReadAll(io.LimitReader(reader,maxOutput+1));if err!=nil{return nil,err};if int64(len(data))>maxOutput{return nil,fmt.Errorf("output exceeds maxOutputLength")};return data,nil
	}
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
		if err!=nil{return nil,err};defer reader.Close();maxOutput:=tsNodeZlibMaxOutput(loop,options);limited:=io.LimitReader(reader,maxOutput+1);var data []byte;data,err=io.ReadAll(limited);if err!=nil{return nil,err};if int64(len(data))>maxOutput{return nil,fmt.Errorf("output exceeds maxOutputLength")};return data,nil
	default:return nil,fmt.Errorf("Unknown compression method %s",mode)
	}
	if int64(output.Len())>tsNodeZlibMaxOutput(loop,options){return nil,fmt.Errorf("output exceeds maxOutputLength")};return output.Bytes(),nil
}
func tsNodeZlibErrorValue(err error) tsValue {code:="Z_DATA_ERROR";message:=err.Error();if strings.Contains(message,"maxOutputLength"){code="ERR_BUFFER_TOO_LARGE"};return tsErrorValue(&tsRuntimeError{name:"Error",message:message,fields:map[string]tsValue{"code":tsNodeZlibKey(code),"errno":tsNumberValue(-3)}})}
func tsNodeZlibInfoEngine(loop *tsLoop,mode string,options tsValue,inputLength int)tsValue{
	engine:=tsNodeZlibStream(loop,mode,options);if engine.kind==tsObjectKind{object:=(*tsObject)(engine.ref);count:=tsNumberValue(float64(inputLength));object.set("bytesWritten",count);object.set("bytesRead",count);if state:=tsNodeZlibSingleton.states[uintptr(engine.ref)];state!=nil{state.bytesWritten=int64(inputLength);state.closed=true;if state.encoder!=nil{_ = state.encoder.Close();state.encoder=nil}}};return engine
}
func tsNodeZlibInputLength(value tsValue)int{
	switch value.kind{case tsStringKind:return len(tsText(value));case tsTypedArrayKind,tsDataViewKind,tsArrayBufferKind:return len(tsFSBytes(value));default:return 0}
}
func tsNodeZlibConvenience(loop *tsLoop,mode string,synchronous bool,args []tsValue) tsValue {
	options:=tsArg(args,1);var callback tsValue
	if !synchronous{
		if len(args)>1&&args[1].kind==tsFunctionKind{callback=args[1];options=tsObjectValue(tsNewObject())}else{callback=tsArg(args,2)}
	}
	if tsNullish(options){options=tsObjectValue(tsNewObject())};if mode=="zstdCompress"&&tsNodeZlibInputLength(tsArg(args,0))>0&&!tsTruthy(tsGet(loop,options,tsNodeZlibKey("pledgedSrcSize"))){snapshot:=tsNewObject();tsObjectSpread(snapshot,options);snapshot.set("pledgedSrcSize",tsNumberValue(float64(tsNodeZlibInputLength(tsArg(args,0)))));options=tsObjectValue(snapshot)};tsNodeZlibValidateOptions(loop,mode,options)
	if !synchronous&&callback.kind!=tsFunctionKind{panic(tsThrown{tsNodeZlibError("ERR_INVALID_ARG_TYPE","The \"callback\" argument must be of type function")})}
	input:=tsNodeZlibInput(loop,tsArg(args,0),"buffer");output,err:=tsNodeZlibProcess(loop,mode,input,options)
	if synchronous{if err!=nil{panic(tsThrown{tsNodeZlibErrorValue(err)})};buffer:=tsNodeBuffer(output);if tsTruthy(tsGet(loop,options,tsNodeZlibKey("info"))){result:=tsNewObject();result.set("buffer",buffer);result.set("engine",tsNodeZlibInfoEngine(loop,mode,options,len(input)));return tsObjectValue(result)};return buffer}
	loop.ready=append(loop.ready,func(){if err!=nil{tsCall(loop,callback,tsNodeZlibErrorValue(err));return};buffer:=tsNodeBuffer(output);if tsTruthy(tsGet(loop,options,tsNodeZlibKey("info"))){result:=tsNewObject();result.set("buffer",buffer);result.set("engine",tsNodeZlibInfoEngine(loop,mode,options,len(input)));buffer=tsObjectValue(result)};tsCall(loop,callback,tsNull,buffer)})
	return tsU
}
func tsNodeZlibStream(loop *tsLoop,mode string,options tsValue) tsValue {
	if tsNullish(options){options=tsObjectValue(tsNewObject())};tsNodeZlibValidateOptions(loop,mode,options)
	snapshot:=tsNewObject();tsObjectSpread(snapshot,options);if !strings.HasPrefix(mode,"zstd")&&!tsIsUndefined(tsGet(loop,options,tsNodeZlibKey("dictionary"))){snapshot.set("dictionary",tsNodeBuffer(tsNodeZlibDictionary(loop,options)))};if strings.HasPrefix(mode,"zstd")&&len(tsNodeZlibZstdDictionary(loop,options))>0{snapshot.set("dictionary",tsNodeBuffer(tsNodeZlibZstdDictionary(loop,options)))};optionSnapshot:=tsObjectValue(snapshot)
	level,strategy:=-1,0;if !strings.HasPrefix(mode,"brotli")&&!strings.HasPrefix(mode,"zstd"){var levelErr error;level,levelErr=tsNodeZlibLevel(loop,optionSnapshot);if levelErr!=nil{panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE",levelErr.Error())})};strategy=tsNodeZlibCheckInt(loop,optionSnapshot,"strategy",0,4,0)};var err error
	state:=&tsNodeZlibState{mode:mode,options:optionSnapshot,level:level,strategy:strategy,output:&bytes.Buffer{}}
	if mode=="deflate"||mode=="gzip"||mode=="deflateRaw"||mode=="brotliCompress"||mode=="zstdCompress"{dictionary:=tsNodeZlibDictionary(loop,optionSnapshot);if mode=="zstdCompress"{dictionary=tsNodeZlibZstdDictionary(loop,optionSnapshot)};state.encoder,state.encoderFlush,err=tsNodeZlibNewEncoder(loop,mode,state.output,level,dictionary,optionSnapshot);if err!=nil{panic(tsThrown{tsNodeZlibError("ERR_INVALID_ARG_VALUE",err.Error())})}}
	streamOptions:=tsNewObject();streamOptions.set("autoDestroy",tsBooleanValue(true));tsObjectSpread(streamOptions,optionSnapshot)
	if tsTruthy(tsGet(loop,optionSnapshot,tsNodeZlibKey("encoding")))||tsTruthy(tsGet(loop,optionSnapshot,tsNodeZlibKey("objectMode")))||tsTruthy(tsGet(loop,optionSnapshot,tsNodeZlibKey("writableObjectMode"))){streamOptions.set("encoding",tsNull);streamOptions.set("objectMode",tsBooleanValue(false));streamOptions.set("writableObjectMode",tsBooleanValue(false))}
	streamOptions.set("transform",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{
		chunk,callback:=tsArg(args,0),tsArg(args,2);if state.closed{return tsCall(loop,callback,tsNodeZlibError("ERR_INVALID_STATE","zlib binding closed"))};input:=tsNodeZlibInput(loop,chunk,"chunk");state.bytesWritten+=int64(len(input));if state.stream.kind==tsObjectKind{obj:=(*tsObject)(state.stream.ref);obj.set("bytesWritten",tsNumberValue(float64(state.bytesWritten)));obj.set("bytesRead",tsNumberValue(float64(state.bytesWritten)))}
		if state.encoder!=nil{if _,writeErr:=state.encoder.Write(input);writeErr!=nil{return tsCall(loop,callback,tsNodeZlibErrorValue(writeErr))};if state.encoderFlush!=nil{if flushErr:=state.encoderFlush();flushErr!=nil{return tsCall(loop,callback,tsNodeZlibErrorValue(flushErr))}};data:=append([]byte(nil),state.output.Bytes()...);state.output.Reset();if len(data)>0{return tsCall(loop,callback,tsNull,tsNodeBuffer(data))};return tsCall(loop,callback,tsNull)}
		state.input=append(state.input,input...);return tsCall(loop,callback,tsNull)
	})))
	streamOptions.set("flush",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{
		callback:=tsArg(args,0);if state.closed{return tsCall(loop,callback,tsNodeZlibError("ERR_INVALID_STATE","zlib binding closed"))}
		if state.encoder!=nil{if closeErr:=state.encoder.Close();closeErr!=nil{return tsCall(loop,callback,tsNodeZlibErrorValue(closeErr))};state.encoder=nil;data:=append([]byte(nil),state.output.Bytes()...);state.output.Reset();return tsCall(loop,callback,tsNull,tsNodeBuffer(data))}
		data,processErr:=tsNodeZlibProcess(loop,state.mode,state.input,state.options);if processErr!=nil{return tsCall(loop,callback,tsNodeZlibErrorValue(processErr))};return tsCall(loop,callback,tsNull,tsNodeBuffer(data))
	})))
	module:=tsInternalStreamTransformModule(loop);constructor:=tsGet(loop,module,tsNodeZlibKey("Transform"));stream:=tsConstruct(constructor,tsObjectValue(streamOptions));state.stream=stream;tsNodeZlibSingleton.states[uintptr((*tsObject)(stream.ref))]=state
	object:=(*tsObject)(stream.ref);object.set("_destroy",tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{cause,callback:=tsArg(args,0),tsArg(args,1);if state.encoder!=nil{_ = state.encoder.Close();state.encoder=nil};state.closed=true;if callback.kind==tsFunctionKind{tsCall(loop,callback,cause)};return tsU})));object.set("bytesWritten",tsNumberValue(0));if object.descriptors==nil{object.descriptors=map[string]*tsDescriptor{}};object.descriptors["bytesRead"]=&tsDescriptor{accessor:true,configurable:true,getter:tsNamedNativeMethod("get bytesRead",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{return tsGet(receiver,tsNodeZlibKey("bytesWritten"))}),setter:tsNamedNativeMethod("set bytesRead",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{value:=tsArg(args,0);(*tsObject)(receiver.ref).set("bytesWritten",value);if current:=tsNodeZlibSingleton.states[uintptr(receiver.ref)];current!=nil&&tsIsNumeric(value){current.bytesWritten=int64(tsNumber(value))};return tsU})};object.descriptors["_closed"]=&tsDescriptor{accessor:true,enumerable:true,configurable:true,getter:tsNamedNativeMethod("get _closed",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{current:=tsNodeZlibSingleton.states[uintptr(receiver.ref)];return tsBooleanValue(current==nil||current.closed)})};object.set("_mode",tsNodeZlibKey(mode));chunkSize:=tsNodeZlibChunkSize(loop,optionSnapshot);object.set("_chunkSize",chunkSize);object.set("_maxOutputLength",tsNumberValue(float64(tsNodeZlibMaxOutput(loop,optionSnapshot))));object.set("_info",tsGet(loop,optionSnapshot,tsNodeZlibKey("info")))
	flushDefault,finishDefault,fullFlushDefault:=float64(0),float64(4),float64(3);if strings.HasPrefix(mode,"brotli"){finishDefault,fullFlushDefault=2,1};if strings.HasPrefix(mode,"zstd"){finishDefault,fullFlushDefault=2,1};object.set("_defaultFlushFlag",tsNodeZlibOption(loop,optionSnapshot,"flush",tsNumberValue(flushDefault)));object.set("_finishFlushFlag",tsNodeZlibOption(loop,optionSnapshot,"finishFlush",tsNumberValue(finishDefault)));object.set("_defaultFullFlushFlag",tsNodeZlibOption(loop,optionSnapshot,"fullFlush",tsNumberValue(fullFlushDefault)))
	object.set("reset",tsNamedNativeMethod("reset",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{current:=tsNodeZlibSingleton.states[uintptr(receiver.ref)];if current==nil||current.closed{tsPropertyFailure("zlib binding closed")};if current.encoder!=nil{_ = current.encoder.Close()};current.input=nil;current.bytesWritten=0;current.output.Reset();if current.mode=="deflate"||current.mode=="gzip"||current.mode=="deflateRaw"||current.mode=="brotliCompress"||current.mode=="zstdCompress"{dictionary:=tsNodeZlibDictionary(loop,current.options);if current.mode=="zstdCompress"{dictionary=tsNodeZlibZstdDictionary(loop,current.options)};current.encoder,current.encoderFlush,_=tsNodeZlibNewEncoder(loop,current.mode,current.output,current.level,dictionary,current.options)};receiverObj:=(*tsObject)(receiver.ref);receiverObj.set("bytesWritten",tsNumberValue(0));receiverObj.set("bytesRead",tsNumberValue(0));return tsU}))
	object.set("close",tsNamedNativeMethod("close",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{if current:=tsNodeZlibSingleton.states[uintptr(receiver.ref)];current!=nil{if current.encoder!=nil{_ = current.encoder.Close()};current.closed=true};if callback:=tsArg(args,0);callback.kind==tsFunctionKind{tsNodeWritableDefer(loop,func(){tsCall(loop,callback)})};tsNodeClassicMethod(loop,receiver,"destroy",tsU);return tsU}))
	object.set("flush",tsNamedNativeMethod("flush",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{callback:=tsArg(args,1);if tsArg(args,0).kind==tsFunctionKind{callback=tsArg(args,0)};if current:=tsNodeZlibSingleton.states[uintptr(receiver.ref)];current!=nil&&current.encoder!=nil&&!current.closed{if err:=current.encoderFlush();err!=nil{return tsCall(loop,callback,tsNodeZlibErrorValue(err))};data:=append([]byte(nil),current.output.Bytes()...);current.output.Reset();if len(data)>0{tsNodeClassicMethod(loop,receiver,"push",tsNodeBuffer(data))}};if callback.kind==tsFunctionKind{tsNodeWritableDefer(loop,func(){tsCall(loop,callback)})};return tsU}))
	object.set("params",tsNamedNativeMethod("params",3,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{level,strategy:=tsArg(args,0),tsArg(args,1);tsInternalValidatorInt32(level,tsNodeZlibKey("level"),tsNumberValue(-1),tsNumberValue(9));tsInternalValidatorInt32(strategy,tsNodeZlibKey("strategy"),tsNumberValue(0),tsNumberValue(4));if current:=tsNodeZlibSingleton.states[uintptr(receiver.ref)];current!=nil{current.level=int(tsNumber(level));current.strategy=int(tsNumber(strategy))};if callback:=tsArg(args,2);callback.kind==tsFunctionKind{tsNodeWritableDefer(loop,func(){tsCall(loop,callback)})};return tsU}))
	return stream
}
func tsNodeZlibCtor(mode,name string) tsValue {
	function:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{if strings.HasPrefix(mode,"zstd"){tsPropertyFailure("Class constructor "+name+" cannot be invoked without 'new'")};return tsNodeZlibStream(loop,mode,tsArg(args,0))});function.name=name;function.length=1;function.constructible=true;function.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeZlibStream(loop,mode,tsArg(args,0))};function.properties=tsNewObject();return tsFunctionValue(function)
}
func tsNodeZlibCRC(loop *tsLoop,args []tsValue) tsValue {
	if len(args)==0{panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_TYPE","The \"data\" argument must be of type string or an instance of Buffer, TypedArray, DataView, or ArrayBuffer")})};data:=tsNodeZlibInput(loop,tsArg(args,0),"data");seed:=uint32(0)
	if value:=tsArg(args,1);!tsIsUndefined(value){if !tsIsNumeric(value){panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_TYPE","The \"value\" argument must be of type number")})};number:=tsNumber(value);if math.IsNaN(number)||math.IsInf(number,0)||math.Trunc(number)!=number{panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","The value of \"value\" must be an integer")})};if number<0||number>math.MaxUint32{panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","The value of \"value\" is out of range. It must be >= 0 and <= 4294967295")})};seed=uint32(number)}
	return tsNumberValue(float64(crc32.Update(seed,crc32.IEEETable,data)))}
func tsNodeZlibConstants() *tsObject {
	values:=map[string]float64{
		"Z_NO_FLUSH":0,"Z_PARTIAL_FLUSH":1,"Z_SYNC_FLUSH":2,"Z_FULL_FLUSH":3,"Z_FINISH":4,"Z_BLOCK":5,"Z_TREES":6,"Z_OK":0,"Z_STREAM_END":1,"Z_NEED_DICT":2,"Z_ERRNO":-1,"Z_STREAM_ERROR":-2,"Z_DATA_ERROR":-3,"Z_MEM_ERROR":-4,"Z_BUF_ERROR":-5,"Z_VERSION_ERROR":-6,
		"Z_NO_COMPRESSION":0,"Z_BEST_SPEED":1,"Z_BEST_COMPRESSION":9,"Z_DEFAULT_COMPRESSION":-1,"Z_FILTERED":1,"Z_HUFFMAN_ONLY":2,"Z_RLE":3,"Z_FIXED":4,"Z_DEFAULT_STRATEGY":0,"ZLIB_VERNUM":0x12d0,
		"DEFLATE":1,"INFLATE":2,"GZIP":3,"GUNZIP":4,"DEFLATERAW":5,"INFLATERAW":6,"UNZIP":7,"BROTLI_DECODE":8,"BROTLI_ENCODE":9,"ZSTD_COMPRESS":10,"ZSTD_DECOMPRESS":11,
		"Z_MIN_WINDOWBITS":8,"Z_MAX_WINDOWBITS":15,"Z_DEFAULT_WINDOWBITS":15,"Z_MIN_CHUNK":64,"Z_MAX_CHUNK":math.Inf(1),"Z_DEFAULT_CHUNK":16384,"Z_MIN_MEMLEVEL":1,"Z_MAX_MEMLEVEL":9,"Z_DEFAULT_MEMLEVEL":8,"Z_MIN_LEVEL":-1,"Z_MAX_LEVEL":9,"Z_DEFAULT_LEVEL":-1,
		"BROTLI_OPERATION_PROCESS":0,"BROTLI_OPERATION_FLUSH":1,"BROTLI_OPERATION_FINISH":2,"BROTLI_OPERATION_EMIT_METADATA":3,"BROTLI_PARAM_MODE":0,"BROTLI_MODE_GENERIC":0,"BROTLI_MODE_TEXT":1,"BROTLI_MODE_FONT":2,"BROTLI_DEFAULT_MODE":0,"BROTLI_PARAM_QUALITY":1,"BROTLI_MIN_QUALITY":0,"BROTLI_MAX_QUALITY":11,"BROTLI_DEFAULT_QUALITY":11,"BROTLI_PARAM_LGWIN":2,"BROTLI_MIN_WINDOW_BITS":10,"BROTLI_MAX_WINDOW_BITS":24,"BROTLI_LARGE_MAX_WINDOW_BITS":30,"BROTLI_DEFAULT_WINDOW":22,"BROTLI_PARAM_LGBLOCK":3,"BROTLI_MIN_INPUT_BLOCK_BITS":16,"BROTLI_MAX_INPUT_BLOCK_BITS":24,"BROTLI_PARAM_DISABLE_LITERAL_CONTEXT_MODELING":4,"BROTLI_PARAM_SIZE_HINT":5,"BROTLI_PARAM_LARGE_WINDOW":6,"BROTLI_PARAM_NPOSTFIX":7,"BROTLI_PARAM_NDIRECT":8,
		"BROTLI_DECODER_RESULT_ERROR":0,"BROTLI_DECODER_RESULT_SUCCESS":1,"BROTLI_DECODER_RESULT_NEEDS_MORE_INPUT":2,"BROTLI_DECODER_RESULT_NEEDS_MORE_OUTPUT":3,"BROTLI_DECODER_PARAM_DISABLE_RING_BUFFER_REALLOCATION":0,"BROTLI_DECODER_PARAM_LARGE_WINDOW":1,"BROTLI_DECODER_NO_ERROR":0,"BROTLI_DECODER_SUCCESS":1,"BROTLI_DECODER_NEEDS_MORE_INPUT":2,"BROTLI_DECODER_NEEDS_MORE_OUTPUT":3,
		"BROTLI_DECODER_ERROR_FORMAT_EXUBERANT_NIBBLE":-1,"BROTLI_DECODER_ERROR_FORMAT_RESERVED":-2,"BROTLI_DECODER_ERROR_FORMAT_EXUBERANT_META_NIBBLE":-3,"BROTLI_DECODER_ERROR_FORMAT_SIMPLE_HUFFMAN_ALPHABET":-4,"BROTLI_DECODER_ERROR_FORMAT_SIMPLE_HUFFMAN_SAME":-5,"BROTLI_DECODER_ERROR_FORMAT_CL_SPACE":-6,"BROTLI_DECODER_ERROR_FORMAT_HUFFMAN_SPACE":-7,"BROTLI_DECODER_ERROR_FORMAT_CONTEXT_MAP_REPEAT":-8,"BROTLI_DECODER_ERROR_FORMAT_BLOCK_LENGTH_1":-9,"BROTLI_DECODER_ERROR_FORMAT_BLOCK_LENGTH_2":-10,"BROTLI_DECODER_ERROR_FORMAT_TRANSFORM":-11,"BROTLI_DECODER_ERROR_FORMAT_DICTIONARY":-12,"BROTLI_DECODER_ERROR_FORMAT_WINDOW_BITS":-13,"BROTLI_DECODER_ERROR_FORMAT_PADDING_1":-14,"BROTLI_DECODER_ERROR_FORMAT_PADDING_2":-15,"BROTLI_DECODER_ERROR_FORMAT_DISTANCE":-16,"BROTLI_DECODER_ERROR_DICTIONARY_NOT_SET":-19,"BROTLI_DECODER_ERROR_INVALID_ARGUMENTS":-20,"BROTLI_DECODER_ERROR_ALLOC_CONTEXT_MODES":-21,"BROTLI_DECODER_ERROR_ALLOC_TREE_GROUPS":-22,"BROTLI_DECODER_ERROR_ALLOC_CONTEXT_MAP":-25,"BROTLI_DECODER_ERROR_ALLOC_RING_BUFFER_1":-26,"BROTLI_DECODER_ERROR_ALLOC_RING_BUFFER_2":-27,"BROTLI_DECODER_ERROR_ALLOC_BLOCK_TYPE_TREES":-30,"BROTLI_DECODER_ERROR_UNREACHABLE":-31,
		"ZSTD_e_continue":0,"ZSTD_e_flush":1,"ZSTD_e_end":2,"ZSTD_fast":1,"ZSTD_dfast":2,"ZSTD_greedy":3,"ZSTD_lazy":4,"ZSTD_lazy2":5,"ZSTD_btlazy2":6,"ZSTD_btopt":7,"ZSTD_btultra":8,"ZSTD_btultra2":9,
		"ZSTD_c_compressionLevel":100,"ZSTD_c_windowLog":101,"ZSTD_c_hashLog":102,"ZSTD_c_chainLog":103,"ZSTD_c_searchLog":104,"ZSTD_c_minMatch":105,"ZSTD_c_targetLength":106,"ZSTD_c_strategy":107,"ZSTD_c_enableLongDistanceMatching":160,"ZSTD_c_ldmHashLog":161,"ZSTD_c_ldmMinMatch":162,"ZSTD_c_ldmBucketSizeLog":163,"ZSTD_c_ldmHashRateLog":164,"ZSTD_c_contentSizeFlag":200,"ZSTD_c_checksumFlag":201,"ZSTD_c_dictIDFlag":202,"ZSTD_c_nbWorkers":400,"ZSTD_c_jobSize":401,"ZSTD_c_overlapLog":402,"ZSTD_d_windowLogMax":100,"ZSTD_CLEVEL_DEFAULT":3,
		"ZSTD_error_no_error":0,"ZSTD_error_GENERIC":1,"ZSTD_error_prefix_unknown":10,"ZSTD_error_version_unsupported":12,"ZSTD_error_frameParameter_unsupported":14,"ZSTD_error_frameParameter_windowTooLarge":16,"ZSTD_error_corruption_detected":20,"ZSTD_error_checksum_wrong":22,"ZSTD_error_literals_headerWrong":24,"ZSTD_error_dictionary_corrupted":30,"ZSTD_error_dictionary_wrong":32,"ZSTD_error_dictionaryCreation_failed":34,"ZSTD_error_parameter_unsupported":40,"ZSTD_error_parameter_combination_unsupported":41,"ZSTD_error_parameter_outOfBound":42,"ZSTD_error_tableLog_tooLarge":44,"ZSTD_error_maxSymbolValue_tooLarge":46,"ZSTD_error_maxSymbolValue_tooSmall":48,"ZSTD_error_stabilityCondition_notRespected":50,"ZSTD_error_stage_wrong":60,"ZSTD_error_init_missing":62,"ZSTD_error_memory_allocation":64,"ZSTD_error_workSpace_tooSmall":66,"ZSTD_error_dstSize_tooSmall":70,"ZSTD_error_srcSize_wrong":72,"ZSTD_error_dstBuffer_null":74,"ZSTD_error_noForwardProgress_destFull":80,"ZSTD_error_noForwardProgress_inputEmpty":82,
	}
	out:=tsNewObject();names:=make([]string,0,len(values));for name:=range values{names=append(names,name)};sort.Strings(names);for _,name:=range names{out.set(name,tsNumberValue(values[name]))};return out
}
func tsNodeZlibFreeze(object *tsObject){if object.descriptors==nil{object.descriptors=map[string]*tsDescriptor{}};for _,name:=range object.order{object.descriptors[name]=&tsDescriptor{writable:false,enumerable:true,configurable:false}};object.nonExtensible=true}
func tsNodeZlibDefine(object *tsObject,name string,value tsValue,enumerable,writable,configurable bool){object.set(name,value);if object.descriptors==nil{object.descriptors=map[string]*tsDescriptor{}};object.descriptors[name]=&tsDescriptor{writable:writable,enumerable:enumerable,configurable:configurable}}
func tsNodeZlibModule(loop *tsLoop) tsValue {
	tsNodeZlibSingleton.Do(func(){
		tsNodeZlibSingleton.states=map[uintptr]*tsNodeZlibState{};tsNodeZlibSingleton.transform=tsGet(loop,tsInternalStreamTransformModule(loop),tsNodeZlibKey("Transform"));module:=tsNewObject()
		classes:=[]struct{name,mode string}{{"Deflate","deflate"},{"Inflate","inflate"},{"Gzip","gzip"},{"Gunzip","gunzip"},{"DeflateRaw","deflateRaw"},{"InflateRaw","inflateRaw"},{"Unzip","unzip"},{"BrotliCompress","brotliCompress"},{"BrotliDecompress","brotliDecompress"},{"ZstdCompress","zstdCompress"},{"ZstdDecompress","zstdDecompress"}}
		module.set("crc32",tsNamedNativeMethod("crc32",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeZlibCRC(loop,args)}))
		for _,class:=range classes{module.set(class.name,tsNodeZlibCtor(class.mode,class.name))}
		conveniences:=[]struct{name,mode string;sync bool}{{"deflate","deflate",false},{"deflateSync","deflate",true},{"gzip","gzip",false},{"gzipSync","gzip",true},{"deflateRaw","deflateRaw",false},{"deflateRawSync","deflateRaw",true},{"unzip","unzip",false},{"unzipSync","unzip",true},{"inflate","inflate",false},{"inflateSync","inflate",true},{"gunzip","gunzip",false},{"gunzipSync","gunzip",true},{"inflateRaw","inflateRaw",false},{"inflateRawSync","inflateRaw",true},{"brotliCompress","brotliCompress",false},{"brotliCompressSync","brotliCompress",true},{"brotliDecompress","brotliDecompress",false},{"brotliDecompressSync","brotliDecompress",true},{"zstdCompress","zstdCompress",false},{"zstdCompressSync","zstdCompress",true},{"zstdDecompress","zstdDecompress",false},{"zstdDecompressSync","zstdDecompress",true}}
		for _,method:=range conveniences{name,mode,syncMode:=method.name,method.mode,method.sync;length:=3;if syncMode{length=2};module.set(name,tsNamedNativeMethod(name,length,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeZlibConvenience(loop,mode,syncMode,args)}))}
		for _,class:=range []struct{name,mode string}{{"Deflate","deflate"},{"Inflate","inflate"},{"DeflateRaw","deflateRaw"},{"InflateRaw","inflateRaw"},{"Gzip","gzip"},{"Gunzip","gunzip"},{"Unzip","unzip"},{"BrotliCompress","brotliCompress"},{"BrotliDecompress","brotliDecompress"},{"ZstdCompress","zstdCompress"},{"ZstdDecompress","zstdDecompress"}}{name,mode:=class.name,class.mode;module.set("create"+name,tsNamedNativeMethod("create"+name,1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeZlibStream(loop,mode,tsArg(args,0))}))}
		constants:=tsNodeZlibConstants();codes:=tsNewObject();codeNames:=[]string{"Z_OK","Z_STREAM_END","Z_NEED_DICT","Z_ERRNO","Z_STREAM_ERROR","Z_DATA_ERROR","Z_MEM_ERROR","Z_BUF_ERROR","Z_VERSION_ERROR"};for _,name:=range codeNames{codes.set(name,constants.get(name))};for _,name:=range codeNames{codes.set(strconv.Itoa(int(tsNumber(constants.get(name)))),tsNodeZlibKey(name))};tsNodeZlibFreeze(constants);tsNodeZlibFreeze(codes);module.set("constants",tsObjectValue(constants));module.set("codes",tsObjectValue(codes))
		for _,name:=range []string{"Z_NO_FLUSH","Z_PARTIAL_FLUSH","Z_SYNC_FLUSH","Z_FULL_FLUSH","Z_FINISH","Z_BLOCK","Z_OK","Z_STREAM_END","Z_NEED_DICT","Z_ERRNO","Z_STREAM_ERROR","Z_DATA_ERROR","Z_MEM_ERROR","Z_BUF_ERROR","Z_VERSION_ERROR","Z_NO_COMPRESSION","Z_BEST_SPEED","Z_BEST_COMPRESSION","Z_DEFAULT_COMPRESSION","Z_FILTERED","Z_HUFFMAN_ONLY","Z_RLE","Z_FIXED","Z_DEFAULT_STRATEGY","ZLIB_VERNUM","DEFLATE","INFLATE","GZIP","GUNZIP","DEFLATERAW","INFLATERAW","UNZIP","Z_MIN_WINDOWBITS","Z_MAX_WINDOWBITS","Z_DEFAULT_WINDOWBITS","Z_MIN_CHUNK","Z_MAX_CHUNK","Z_DEFAULT_CHUNK","Z_MIN_MEMLEVEL","Z_MAX_MEMLEVEL","Z_DEFAULT_MEMLEVEL","Z_MIN_LEVEL","Z_MAX_LEVEL","Z_DEFAULT_LEVEL"}{tsNodeZlibDefine(module,name,constants.get(name),false,false,false)}
		module.set("default",tsObjectValue(module));tsNodeZlibSingleton.module=tsObjectValue(module)
	});return tsNodeZlibSingleton.module
}
`
