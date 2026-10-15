package goemit

// internalFSStreamsRuntime exposes the filesystem stream constructors already
// backed by the shared fs stream state machine, plus the FileSink wrapper hook.
const internalFSStreamsRuntime = `
var tsInternalFSStreamsSingleton struct { sync.Once; value tsValue; fastPath tsValue }

func tsInternalFSStreamsModule(loop *tsLoop) tsValue {
	tsInternalFSStreamsSingleton.Do(func(){
		module:=tsNewObject()
		fs:=tsFSModule("fs")
		module.set("ReadStream",tsFSProperty(fs,"ReadStream"))
		module.set("WriteStream",tsFSProperty(fs,"WriteStream"))
		fastPath:=tsNewSymbol(tsStringUTF8("kWriteStreamFastPath"))
		module.set("kWriteStreamFastPath",fastPath)
		module.set("writableFromFileSink",tsNamedNativeMethod("writableFromFileSink",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{
			sink:=tsArg(args,0)
			if sink.kind!=tsObjectKind { tsFSInvalid("fileSink","an object") }
			if tsFSProperty(sink,"write").kind!=tsFunctionKind||tsFSProperty(sink,"end").kind!=tsFunctionKind { tsFSInvalid("fileSink","an object with write() and end() functions") }
			events:=tsFSEventObject(loop);object:=events.object;object.nativeClass=tsNativeClass("WriteStream")
			object.set("fd",tsNull);object.set("path",tsU);object.set("closed",tsBooleanValue(false));object.set("destroyed",tsBooleanValue(false));object.set("writable",tsBooleanValue(true));object.set(tsPropertyKey(fastPath),sink)
			ended:=false
			object.set("write",tsNamedNativeMethod("write",1,func(loop *tsLoop,receiver tsValue,values ...tsValue)tsValue{
				if ended { failure:=tsFSException("Error","ERR_STREAM_WRITE_AFTER_END","write after end");callback:=tsArg(values,2);if callback.kind==tsFunctionKind{loop.ready=append(loop.ready,func(){tsCall(loop,callback,failure)})};return tsBooleanValue(false) }
				data,encoding,callback:=tsArg(values,0),tsArg(values,1),tsArg(values,2);if encoding.kind==tsFunctionKind{callback=encoding}
				result:=tsCallReceiver(loop,tsFSProperty(sink,"write"),sink,data)
				if callback.kind==tsFunctionKind { complete:=func(value tsResult){err:=tsNull;if value.rejected{err=value.value};tsCall(loop,callback,err)};if result.kind==tsPromiseKind{loop.await(result,complete)}else{loop.ready=append(loop.ready,func(){complete(tsResult{tsU,false})})} }
				return tsBooleanValue(true)
			}))
			object.set("end",tsNamedNativeMethod("end",0,func(loop *tsLoop,receiver tsValue,values ...tsValue)tsValue{
				if ended{return tsObjectValue(object)};ended=true;object.set("writableEnded",tsBooleanValue(true));data,callback:=tsArg(values,0),tsArg(values,1);if data.kind==tsFunctionKind{callback=data;data=tsU};if !tsIsUndefined(data){_ = tsCallReceiver(loop,tsFSProperty(sink,"write"),sink,data)}
				result:=tsCallReceiver(loop,tsFSProperty(sink,"end"),sink);finish:=func(value tsResult){if value.rejected{events.emitEvent("error",value.value);return};object.set("writable",tsBooleanValue(false));object.set("writableFinished",tsBooleanValue(true));events.emitEvent("finish");if callback.kind==tsFunctionKind{tsCall(loop,callback,tsNull)}}
				if result.kind==tsPromiseKind{loop.await(result,finish)}else{loop.ready=append(loop.ready,func(){finish(tsResult{tsU,false})})};return tsObjectValue(object)
			}))
			object.set("destroy",tsNamedNativeMethod("destroy",0,func(loop *tsLoop,receiver tsValue,values ...tsValue)tsValue{object.set("destroyed",tsBooleanValue(true));object.set("writable",tsBooleanValue(false));events.emitEvent("close");return tsObjectValue(object)}))
			object.set("close",object.get("destroy"))
			return tsObjectValue(object)
		}))
		module.set("default",tsObjectValue(module))
		tsInternalFSStreamsSingleton.fastPath=fastPath
		tsInternalFSStreamsSingleton.value=tsObjectValue(module)
	})
	return tsInternalFSStreamsSingleton.value
}
`
