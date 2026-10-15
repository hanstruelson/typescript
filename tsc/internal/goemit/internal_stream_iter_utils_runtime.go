package goemit

// internalStreamIterUtilsRuntime adapts the common stream/iter validation,
// chunk conversion, signal, cursor, and transform protocol helpers.
const internalStreamIterUtilsRuntime = `
var tsInternalStreamIterUtils struct { sync.Once; module *tsObject }

func tsInternalStreamIterUtilsClampHWM(value tsValue) tsValue {
	n:=math.Trunc(tsNumber(value));if math.IsNaN(n)||n<1 { n=1 };if n>9007199254740991 { n=9007199254740991 };return tsNumberValue(n)
}
func tsInternalStreamIterUtilsOnAbort(loop *tsLoop,args []tsValue) tsValue {
	signal,handler:=tsArg(args,0),tsArg(args,1);if tsTruthy(tsGet(loop,signal,tsStringReference(tsStringUTF8("aborted")))) { tsCall(loop,handler);return tsU }
	if handler.kind!=tsFunctionKind { tsInternalValidatorInvalidType("handler","a function",handler) }
	add:=tsGet(loop,signal,tsStringReference(tsStringUTF8("addEventListener")));options:=tsNewObject();options.set("once",tsBooleanValue(true));tsCallReceiver(loop,add,signal,tsStringReference(tsStringUTF8("abort")),handler,tsObjectValue(options));return tsU
}
func tsInternalStreamIterUtilsMinCursor(loop *tsLoop,args []tsValue) tsValue {
	consumers,fallback:=tsArg(args,0),tsArg(args,1);minimum:=tsNumber(fallback);count:=0;iterator:=tsIterate(consumers)
	for iterator.next(){cursor:=tsNumber(tsGet(loop,iterator.value,tsStringReference(tsStringUTF8("cursor"))));if cursor<minimum{minimum=cursor;count=1}else if cursor==minimum{count++}}
	out:=tsNewObject();out.set("minCursor",tsNumberValue(minimum));out.set("minCursorConsumers",tsNumberValue(float64(count)));return tsObjectValue(out)
}
func tsInternalStreamIterUtilsToUint8Array(loop *tsLoop,chunk tsValue) tsValue {
	if chunk.kind==tsStringKind { return tsNodeBuffer(tsFSEncode((*tsString)(chunk.ref),"utf8")) }
	if chunk.kind!=tsTypedArrayKind||(*tsTypedArray)(chunk.ref).name!="Uint8Array" { tsInternalValidatorInvalidType("chunk","a string or Uint8Array",chunk) }
	return chunk
}
func tsInternalStreamIterUtilsAllUint8Array(chunks tsValue) bool {
	if chunks.kind!=tsArrayKind { return false };for _,chunk:=range (*tsArray)(chunks.ref).values { if chunk.kind==tsStringKind{return false} };return true
}
func tsInternalStreamIterUtilsConvertChunks(loop *tsLoop,chunks tsValue) tsValue {
	if chunks.kind!=tsArrayKind { tsInternalValidatorInvalidType("chunks","an array",chunks) }
	out:=&tsArray{};for _,chunk:=range (*tsArray)(chunks.ref).values {out.values=append(out.values,tsInternalStreamIterUtilsToUint8Array(loop,chunk))};return tsArrayValue(out)
}
func tsInternalStreamIterUtilsConcatBytes(chunks tsValue) tsValue {
	if chunks.kind!=tsArrayKind { tsInternalValidatorInvalidType("chunks","an array",chunks) };array:=(*tsArray)(chunks.ref);if array.length()==0{return tsNodeBuffer(nil)}
	if array.length()==1 {item:=array.values[0];if item.kind!=tsTypedArrayKind||(*tsTypedArray)(item.ref).name!="Uint8Array" {tsInternalValidatorInvalidType("chunks[0]","a Uint8Array",item)};data:=tsFSBytes(item);return tsNodeBuffer(append([]byte(nil),data...))}
	total:=0;for _,item:=range array.values {if item.kind!=tsTypedArrayKind||(*tsTypedArray)(item.ref).name!="Uint8Array" {tsInternalValidatorInvalidType("chunks","an array of Uint8Array values",item)};total+=len(tsFSBytes(item))}
	data:=make([]byte,0,total);for _,item:=range array.values{data=append(data,tsFSBytes(item)...)};return tsNodeBuffer(data)
}
func tsInternalStreamIterUtilsWrapError(value tsValue) tsValue {
	if value.kind==tsErrorKind||value.kind==tsObjectKind&&(*tsObject)(value.ref).nativeError{return value}
	return tsFSException("Error","ERR_OPERATION_FAILED",tsText(value))
}
func tsInternalStreamIterUtilsHasProtocol(loop *tsLoop,value,symbol tsValue) bool {
	if tsNullish(value)||value.kind!=tsObjectKind&&value.kind!=tsInstanceKind{return false};return tsGet(loop,value,symbol).kind==tsFunctionKind
}
func tsInternalStreamIterUtilsPullOptions(loop *tsLoop,value tsValue) bool {
	if tsNullish(value)||value.kind!=tsObjectKind&&value.kind!=tsInstanceKind{return false};return !tsHas(value,tsStringReference(tsStringUTF8("transform")))&&!tsHas(value,tsStringReference(tsStringUTF8("write")))
}
func tsInternalStreamIterUtilsIsTransform(loop *tsLoop,value tsValue) bool {
	return value.kind==tsFunctionKind||!tsNullish(value)&&tsGet(loop,value,tsStringReference(tsStringUTF8("transform"))).kind==tsFunctionKind
}
func tsInternalStreamIterUtilsParsePullArgs(loop *tsLoop,args []tsValue) tsValue {
	transforms:=append([]tsValue{},args...);options:=tsU
	if len(transforms)>0&&tsInternalStreamIterUtilsPullOptions(loop,transforms[len(transforms)-1]) {options=transforms[len(transforms)-1];transforms=transforms[:len(transforms)-1]}
	for index,value:=range transforms {if !tsInternalStreamIterUtilsIsTransform(loop,value){tsInternalValidatorInvalidType("transforms["+strconv.Itoa(index)+"]","a function or object with transform()",value)}}
	out:=tsNewObject();out.set("transforms",tsArrayValue(&tsArray{values:transforms}));out.set("options",options);return tsObjectValue(out)
}
func tsInternalStreamIterUtilsValidateBackpressure(loop *tsLoop,value tsValue) tsValue {
	choices:=tsArrayValue(&tsArray{values:[]tsValue{tsStringReference(tsStringUTF8("strict")),tsStringReference(tsStringUTF8("block")),tsStringReference(tsStringUTF8("drop-oldest")),tsStringReference(tsStringUTF8("drop-newest"))}})
	tsInternalValidatorOneOf(value,tsStringReference(tsStringUTF8("options.backpressure")),choices);return tsU
}
func tsInternalStreamIterUtilsModule(loop *tsLoop) tsValue {
	tsInternalStreamIterUtils.Do(func(){
		module:=tsNewObject();module.set("kPushDefaultHWM",tsNumberValue(4));module.set("kMultiConsumerDefaultHWM",tsNumberValue(16));module.set("kResolvedPromise",tsPromiseValue(loop.resolved(tsU,false)))
		methods:=[]string{"clampHWM","concatBytes","convertChunks","getMinCursor","hasProtocol","isPullOptions","isTransform","isTransformObject","onSignalAbort","parsePullArgs","toUint8Array","validateBackpressure","allUint8Array","wrapError"}
		for _,name:=range methods {method:=name;module.set(name,tsNamedNativeMethod(name,1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{switch method {
		case "clampHWM":return tsInternalStreamIterUtilsClampHWM(tsArg(args,0))
		case "concatBytes":return tsInternalStreamIterUtilsConcatBytes(tsArg(args,0))
		case "convertChunks":return tsInternalStreamIterUtilsConvertChunks(loop,tsArg(args,0))
		case "getMinCursor":return tsInternalStreamIterUtilsMinCursor(loop,args)
		case "hasProtocol":return tsBooleanValue(tsInternalStreamIterUtilsHasProtocol(loop,tsArg(args,0),tsArg(args,1)))
		case "isPullOptions":return tsBooleanValue(tsInternalStreamIterUtilsPullOptions(loop,tsArg(args,0)))
		case "isTransform","isTransformObject":return tsBooleanValue(tsInternalStreamIterUtilsIsTransform(loop,tsArg(args,0)))
		case "onSignalAbort":return tsInternalStreamIterUtilsOnAbort(loop,args)
		case "parsePullArgs":value:=tsArg(args,0);if value.kind!=tsArrayKind{tsInternalValidatorInvalidType("args","an array",value)};return tsInternalStreamIterUtilsParsePullArgs(loop,(*tsArray)(value.ref).values)
		case "toUint8Array":return tsInternalStreamIterUtilsToUint8Array(loop,tsArg(args,0))
		case "validateBackpressure":return tsInternalStreamIterUtilsValidateBackpressure(loop,tsArg(args,0))
		case "allUint8Array":return tsBooleanValue(tsInternalStreamIterUtilsAllUint8Array(tsArg(args,0)))
		default:return tsInternalStreamIterUtilsWrapError(tsArg(args,0))
		}}))}
		module.set("default",tsObjectValue(module));tsInternalStreamIterUtils.module=module
	})
	return tsObjectValue(tsInternalStreamIterUtils.module)
}
`
