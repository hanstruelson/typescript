package goemit

// nodeV8Runtime maps Bun's node:v8 wrapper onto Go heap counters and a
// versioned structured-value format. GCProfiler observes Go runtime GC cycles,
// not JSC heap collections; snapshots and serializer data are local formats.
const nodeV8Runtime = `
type tsNodeV8Edge struct { Key string; Target int; Group int; Role string }
type tsNodeV8Node struct { Kind string; Number uint64; Nonce uint64; Boolean bool; Units []uint16; Big string; Bytes []byte; Name string; Length int; Buffer bool; Set bool; Edges []tsNodeV8Edge }
type tsNodeV8Document struct { Version int; Root int; Tagged bool; Nodes []tsNodeV8Node }
type tsNodeV8ProfilerState struct { started bool; start time.Time; gcCount uint32; startUsed, startCapacity uint64 }
var tsNodeV8Singleton struct { sync.Once; module tsValue; mu sync.Mutex; profilers map[uintptr]*tsNodeV8ProfilerState; flags string; versionTag uint32; tagValid bool }
var tsNodeV8Magic=[]byte{0xff,'T','S','V',1}
func tsNodeV8Key(name string)tsValue{return tsStringReference(tsStringUTF8(name))}
func tsNodeV8Error(code,message string){panic(tsThrown{tsFSException("Error",code,message)})}
func tsNodeV8NotImplemented(name string){tsNodeV8Error("ERR_METHOD_NOT_IMPLEMENTED","node:v8 "+name+" is not implemented")}
func tsNodeV8HeapSnapshot()string{var memory runtime.MemStats;runtime.ReadMemStats(&memory);snapshot:=map[string]any{"snapshot":map[string]any{"meta":map[string]any{"node_fields":[]string{"type","name","id","self_size","edge_count","trace_node_id","detachedness"}},"node_count":1,"edge_count":0},"nodes":[]uint64{0,0,1,uint64(memory.HeapAlloc),0,0,0},"edges":[]uint64{},"strings":[]string{"synthetic Go heap snapshot"}};data,_:=json.Marshal(snapshot);return string(data)}
func tsNodeV8HeapStats(loop *tsLoop)tsValue{var mem runtime.MemStats;runtime.ReadMemStats(&mem);total:=uint64(mem.Sys);if osinfo:=tsNodeOSModule();!tsIsUndefined(tsGet(osinfo,tsNodeV8Key("totalmem"))){fn:=tsGet(loop,osinfo,tsNodeV8Key("totalmem"));if fn.kind==tsFunctionKind{candidate:=tsCall(loop,fn);if tsIsNumeric(candidate)&&tsNumber(candidate)>0{total=uint64(tsNumber(candidate))}}};out:=tsNewObject();values:=map[string]uint64{"total_heap_size":mem.HeapSys,"total_heap_size_executable":0,"total_physical_size":mem.Sys,"total_available_size":max(total-uint64(mem.HeapAlloc),0),"used_heap_size":mem.HeapAlloc,"total_allocated_bytes":mem.TotalAlloc,"heap_size_limit":total,"malloced_memory":mem.HeapAlloc,"peak_malloced_memory":mem.TotalAlloc,"does_zap_garbage":0,"number_of_native_contexts":1,"number_of_detached_contexts":0,"total_global_handles_size":8192,"used_global_handles_size":2208,"external_memory":0};for name,value:=range values{out.set(name,tsNumberValue(float64(value)))};return tsObjectValue(out)}
func tsNodeV8HeapSpaceStats(loop *tsLoop)tsValue{names:=[]string{"read_only_space","new_space","old_space","code_space","shared_space","trusted_space","new_large_object_space","large_object_space","code_large_object_space","shared_large_object_space","shared_trusted_space","shared_trusted_large_object_space","trusted_large_object_space"};var mem runtime.MemStats;runtime.ReadMemStats(&mem);out:=&tsArray{};for _,name:=range names{size,used:=uint64(0),uint64(0);if name=="old_space"{size=mem.HeapSys;used=mem.HeapAlloc};item:=tsNewObject();item.set("space_name",tsNodeV8Key(name));item.set("space_size",tsNumberValue(float64(size)));item.set("space_used_size",tsNumberValue(float64(used)));item.set("space_available_size",tsNumberValue(float64(max(size-used,0))));item.set("physical_space_size",tsNumberValue(float64(size)));out.values=append(out.values,tsObjectValue(item))};return tsArrayValue(out)}
func tsNodeV8Encode(loop *tsLoop,value tsValue) []byte {
	tagged:=tsInternalSerializationBuffersTag(value);root:=value;hasBuffers:=tagged.kind==tsArrayKind;if hasBuffers{root=tagged}
	doc:=tsNodeV8Document{Version:1,Tagged:hasBuffers};refs:=map[unsafe.Pointer]int{}
	var visit func(tsValue)(int,error)
	visit=func(item tsValue)(int,error){node:=tsNodeV8Node{};switch item.kind{
	case tsUndefinedKind:node.Kind="undefined"
	case tsNullKind:node.Kind="null"
	case tsBooleanKind:node.Kind="boolean";node.Boolean=item.number!=0
	case tsNumberKind,tsFloat32Kind,tsIntKind,tsInt8Kind,tsInt16Kind,tsInt32Kind,tsInt64Kind,tsUintKind,tsUint8Kind,tsUint16Kind,tsUint32Kind,tsUint64Kind:node.Kind="number";node.Number=math.Float64bits(tsNumber(item))
	case tsStringKind:node.Kind="string";node.Units=append([]uint16(nil),(*tsString)(item.ref).units...)
	case tsBigIntKind:node.Kind="bigint";node.Big=(*big.Int)(item.ref).String()
	case tsSymbolKind,tsFunctionKind,tsClassKind:return 0,fmt.Errorf("Value of type %s cannot be serialized",tsValueTypeName(item))
	case tsArrayBufferKind:if id,ok:=refs[item.ref];ok{return id,nil};node.Kind="arraybuffer";node.Bytes=append([]byte(nil),(*tsArrayBuffer)(item.ref).data...)
	case tsTypedArrayKind:if id,ok:=refs[item.ref];ok{return id,nil};array:=(*tsTypedArray)(item.ref);node.Kind="typedarray";node.Bytes=append([]byte(nil),tsFSBytes(item)...);node.Name=array.name;node.Buffer=array.nodeBuffer
	case tsDataViewKind:if id,ok:=refs[item.ref];ok{return id,nil};view:=(*tsDataView)(item.ref);node.Kind="dataview";node.Bytes=append([]byte(nil),view.buffer.data[view.offset:view.offset+view.length]...)
	case tsCollectionKind:if id,ok:=refs[item.ref];ok{return id,nil};node.Kind="collection";node.Set=(*tsCollection)(item.ref).set
	case tsArrayKind:if id,ok:=refs[item.ref];ok{return id,nil};node.Kind="array";node.Length=(*tsArray)(item.ref).length()
	case tsObjectKind,tsInstanceKind,tsErrorKind:if id,ok:=refs[item.ref];ok{return id,nil};if nonce,ok:=tsNodeNetBlockListSerializationNonce(item);ok{node.Kind="blocklist";node.Nonce=nonce}else{node.Kind="object";if item.kind==tsObjectKind&&(*tsObject)(item.ref).nativeDate{node.Kind="date";node.Number=math.Float64bits((*tsObject)(item.ref).nativeTime)}}
	default:return 0,fmt.Errorf("Value of type %s cannot be serialized",tsValueTypeName(item))
	};id:=len(doc.Nodes);doc.Nodes=append(doc.Nodes,node);if item.kind==tsArrayBufferKind||item.kind==tsTypedArrayKind||item.kind==tsDataViewKind||item.kind==tsCollectionKind||item.kind==tsArrayKind||item.kind==tsObjectKind||item.kind==tsInstanceKind||item.kind==tsErrorKind{refs[item.ref]=id}
		edges:=[]tsNodeV8Edge{};switch item.kind{
		case tsArrayKind:a:=(*tsArray)(item.ref);for index:=0;index<a.length();index++{if a.holes[index]{continue};child,err:=visit(a.at(index));if err!=nil{return 0,err};edges=append(edges,tsNodeV8Edge{Key:strconv.Itoa(index),Target:child,Role:"property"})};for name,childValue:=range a.properties{if d:=a.descriptors[name];d!=nil&&(d.accessor||!d.enumerable){continue};child,err:=visit(childValue);if err!=nil{return 0,err};edges=append(edges,tsNodeV8Edge{Key:name,Target:child,Role:"property"})}
		case tsCollectionKind:c:=(*tsCollection)(item.ref);group:=0;for _,entry:=range c.entries{if !entry.live{continue};key,err:=visit(entry.key);if err!=nil{return 0,err};edges=append(edges,tsNodeV8Edge{Target:key,Group:group,Role:"key"});if !c.set{child,err:=visit(entry.value);if err!=nil{return 0,err};edges=append(edges,tsNodeV8Edge{Target:child,Group:group,Role:"value"})};group++}
		case tsObjectKind,tsInstanceKind,tsErrorKind:if node.Kind!="date"&&node.Kind!="blocklist"{keys:=tsEnumerableOwnKeys(loop,item);for _,key:=range keys.values{desc:=tsOwnDescriptor(loop,item,tsText(key));if desc.kind==tsObjectKind{d:=(*tsObject)(desc.ref);if _,accessor:=d.values["get"];accessor{continue};childValue:=d.values["value"];child,err:=visit(childValue);if err!=nil{return 0,err};edges=append(edges,tsNodeV8Edge{Key:tsText(key),Target:child,Role:"property"})}}}
		};doc.Nodes[id].Edges=edges;return id,nil
	}
	rootID,err:=visit(root);if err!=nil{tsNodeV8Error("ERR_BUFFER_TOO_LARGE",err.Error())};doc.Root=rootID;body,err:=json.Marshal(doc);if err!=nil{tsNodeV8Error("ERR_BUFFER_TOO_LARGE",err.Error())};out:=append([]byte(nil),tsNodeV8Magic...);out=append(out,body...);return out
}
func tsNodeV8Decode(loop *tsLoop,data []byte)tsValue{if !bytes.HasPrefix(data,tsNodeV8Magic){tsNodeV8Error("ERR_BUFFER_UNSUPPORTED","invalid node:v8 serialization header");return tsU};var doc tsNodeV8Document;if err:=json.Unmarshal(data[len(tsNodeV8Magic):],&doc);err!=nil{tsNodeV8Error("ERR_BUFFER_UNSUPPORTED",err.Error());return tsU};if doc.Version!=1||doc.Root<0||doc.Root>=len(doc.Nodes){tsNodeV8Error("ERR_BUFFER_UNSUPPORTED","invalid serialized value graph");return tsU};values:=make([]tsValue,len(doc.Nodes));for index,node:=range doc.Nodes{switch node.Kind{case "undefined":values[index]=tsU;case "null":values[index]=tsNull;case "boolean":values[index]=tsBooleanValue(node.Boolean);case "number":values[index]=tsNumberValue(math.Float64frombits(node.Number));case "string":values[index]=tsStringReference(tsStringUnits(node.Units));case "bigint":number,ok:=new(big.Int).SetString(node.Big,10);if !ok{tsNodeV8Error("ERR_BUFFER_UNSUPPORTED","invalid bigint");return tsU};values[index]=tsBigIntValue(number);case "object":values[index]=tsObjectValue(tsNewObject());case "blocklist":if value,ok:=tsNodeNetBlockListFromSerializationNonce(node.Nonce);ok{values[index]=value}else{tsNodeV8Error("ERR_BUFFER_UNSUPPORTED","BlockList structured-clone reference is no longer available");return tsU};case "date":values[index]=tsFSDateMillis(math.Float64frombits(node.Number));case "array":values[index]=tsArrayValue(&tsArray{values:make([]tsValue,node.Length),holes:map[int]bool{}});case "collection":values[index]=tsCollectionValue(&tsCollection{set:node.Set,index:map[tsCollectionKey]int{}});case "arraybuffer":values[index]=tsArrayBufferValue(&tsArrayBuffer{data:append([]byte(nil),node.Bytes...),maxByteLength:len(node.Bytes)});case "dataview":buffer:=&tsArrayBuffer{data:append([]byte(nil),node.Bytes...),maxByteLength:len(node.Bytes)};values[index]=tsDataViewValue(&tsDataView{buffer:buffer,length:len(node.Bytes),properties:map[string]tsValue{}});case "typedarray":buffer:=tsArrayBufferValue(&tsArrayBuffer{data:append([]byte(nil),node.Bytes...),maxByteLength:len(node.Bytes)});name:=node.Name;if name==""{name="Uint8Array"};arrayValue:=tsTypedArrayBufferExtended(name,[]tsValue{buffer,tsNumberValue(0),tsU});if arrayValue.kind!=tsTypedArrayKind{tsNodeV8Error("ERR_BUFFER_UNSUPPORTED","invalid typed array");return tsU};(*tsTypedArray)(arrayValue.ref).nodeBuffer=node.Buffer;values[index]=arrayValue;default:tsNodeV8Error("ERR_BUFFER_UNSUPPORTED","unknown serialized value kind");return tsU}}
	for index,node:=range doc.Nodes{target:=values[index];switch node.Kind{case "object":object:=(*tsObject)(target.ref);for _,edge:=range node.Edges{object.set(edge.Key,values[edge.Target])};case "array":array:=(*tsArray)(target.ref);for i:=0;i<array.length();i++{array.holes[i]=true};for _,edge:=range node.Edges{if edge.Role!="property"{continue};if offset,err:=strconv.Atoi(edge.Key);err==nil&&strconv.Itoa(offset)==edge.Key&&offset>=0&&offset<array.length(){array.values[offset]=values[edge.Target];delete(array.holes,offset)}else{if array.properties==nil{array.properties=map[string]tsValue{}};array.properties[edge.Key]=values[edge.Target]}};case "collection":collection:=(*tsCollection)(target.ref);pending:=map[int]tsValue{};for _,edge:=range node.Edges{if edge.Role=="key"{pending[edge.Group]=values[edge.Target]}else if edge.Role=="value"{key:=pending[edge.Group];collection.put(key,values[edge.Target]);delete(pending,edge.Group)}};for _,key:=range pending{collection.put(key,key)}}}
	if doc.Tagged{return tsInternalSerializationBuffersRestore(values[doc.Root])};return values[doc.Root]
}
func tsNodeV8Ctor(loop *tsLoop,name,baseMessage string)tsValue{prototype:=tsNewObject();constructor:=tsFunc(func(_ *tsLoop,_ ...tsValue)tsValue{tsNodeV8NotImplemented(baseMessage);return tsU});constructor.name=name;constructor.length=0;constructor.constructible=true;constructor.construct=func(_ *tsLoop,_ ...tsValue)tsValue{tsNodeV8NotImplemented(baseMessage);return tsU};constructor.properties=tsNewObject();constructor.properties.set("prototype",tsObjectValue(prototype));prototype.set("constructor",tsFunctionValue(constructor));return tsFunctionValue(constructor)}
func tsNodeV8Module(loop *tsLoop)tsValue{tsNodeV8Singleton.Do(func(){module:=tsNewObject();
	module.set("cachedDataVersionTag",tsNamedNativeMethod("cachedDataVersionTag",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{if !tsNodeV8Singleton.tagValid{tsNodeV8Singleton.versionTag=crc32.ChecksumIEEE([]byte(runtime.Version()+"/"+runtime.GOOS+"/"+runtime.GOARCH+tsNodeV8Singleton.flags));tsNodeV8Singleton.tagValid=true};return tsNumberValue(float64(tsNodeV8Singleton.versionTag))}))
	deserializer:=tsNodeV8Ctor(loop,"Deserializer","Deserializer");serializer:=tsNodeV8Ctor(loop,"Serializer","Serializer");module.set("Deserializer",deserializer);module.set("Serializer",serializer);module.set("DefaultDeserializer",tsNodeV8Ctor(loop,"DefaultDeserializer","Deserializer"));module.set("DefaultSerializer",tsNodeV8Ctor(loop,"DefaultSerializer","Serializer"))
	isOneByte:=tsNamedNativeMethod("isStringOneByteRepresentation",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{value:=tsArg(args,0);if value.kind!=tsStringKind{tsFSInvalid("content","a string")};for _,unit:=range (*tsString)(value.ref).units{if unit>255{return tsBooleanValue(false)}};return tsBooleanValue(true)});module.set("isStringOneByteRepresentation",isOneByte)
	module.set("getHeapStatistics",tsNamedNativeMethod("getHeapStatistics",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNodeV8HeapStats(loop)}));module.set("getHeapSpaceStatistics",tsNamedNativeMethod("getHeapSpaceStatistics",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNodeV8HeapSpaceStats(loop)}))
	codeStats:=tsNewObject();for _,name:=range []string{"code_and_metadata_size","bytecode_and_metadata_size","external_script_source_size","cpu_profiler_metadata_size"}{codeStats.set(name,tsNumberValue(0))};module.set("getHeapCodeStatistics",tsNamedNativeMethod("getHeapCodeStatistics",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsObjectValue(codeStats)}))
	module.set("getCppHeapStatistics",tsNamedNativeMethod("getCppHeapStatistics",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{kind:=tsArg(args,0);if tsIsUndefined(kind){kind=tsNodeV8Key("detailed")};name:=tsInternalValidatorsModule(loop);validate:=tsGet(loop,name,tsNodeV8Key("validateOneOf"));if validate.kind==tsFunctionKind{values:=tsArrayValue(&tsArray{values:[]tsValue{tsNodeV8Key("brief"),tsNodeV8Key("detailed")}});tsCall(loop,validate,kind,tsNodeV8Key("type"),values)};out:=tsNewObject();for _,key:=range []string{"committed_size_bytes","resident_size_bytes","used_size_bytes"}{out.set(key,tsNumberValue(0))};out.set("space_statistics",tsArrayValue(&tsArray{}));out.set("type_names",tsArrayValue(&tsArray{}));out.set("detail_level",kind);return tsObjectValue(out)}))
	module.set("setFlagsFromString",tsNamedNativeMethod("setFlagsFromString",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{value:=tsArg(args,0);if value.kind!=tsStringKind{tsFSInvalid("flags","a string")};tsNodeV8Singleton.flags+=" "+tsText(value);tsNodeV8Singleton.tagValid=false;return tsU}))
	module.set("serialize",tsNamedNativeMethod("serialize",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeBuffer(tsNodeV8Encode(loop,tsArg(args,0)))}))
	module.set("deserialize",tsNamedNativeMethod("deserialize",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{value:=tsArg(args,0);switch value.kind{case tsTypedArrayKind,tsDataViewKind,tsArrayBufferKind:return tsNodeV8Decode(loop,tsNodeBufferByteView(value));default:tsNodeV8Error("ERR_INVALID_ARG_TYPE","The value argument must be a Buffer, TypedArray, DataView, or ArrayBuffer");return tsU}}))
	for _,name:=range []string{"takeCoverage","stopCoverage","setHeapSnapshotNearHeapLimit"}{method:=name;module.set(method,tsNamedNativeMethod(method,0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{tsNodeV8NotImplemented(method);return tsU}))}
	hooks:=tsNewObject();for _,name:=range []string{"createHook","onInit","onBefore","onAfter","onSettled"}{method:=name;hooks.set(method,tsNamedNativeMethod(method,0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{tsNodeV8NotImplemented(method);return tsU}))};module.set("promiseHooks",tsObjectValue(hooks))
	snapshot:=tsNewObject();for _,name:=range []string{"addDeserializeCallback","addSerializeCallback","setDeserializeMainFunction"}{method:=name;snapshot.set(method,tsNamedNativeMethod(method,1,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{tsNodeV8Error("ERR_NOT_BUILDING_SNAPSHOT","Operation cannot be invoked when not building startup snapshot");return tsU}))};snapshot.set("isBuildingSnapshot",tsNamedNativeMethod("isBuildingSnapshot",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsBooleanValue(false)}));module.set("startupSnapshot",tsObjectValue(snapshot))
	profilerPrototype := tsNewObject()
	profiler := tsFunc(func(_ *tsLoop, _ ...tsValue) tsValue {
		object := tsNewObject()
		object.prototype = tsObjectValue(profilerPrototype)
		value := tsObjectValue(object)
		tsNodeV8Singleton.mu.Lock()
		if tsNodeV8Singleton.profilers == nil {
			tsNodeV8Singleton.profilers = map[uintptr]*tsNodeV8ProfilerState{}
		}
		tsNodeV8Singleton.profilers[uintptr(object)] = &tsNodeV8ProfilerState{}
		tsNodeV8Singleton.mu.Unlock()
		return value
	})
	profiler.name = "GCProfiler"
	profiler.length = 0
	profiler.constructible = true
	profiler.construct = func(loop *tsLoop, _ ...tsValue) tsValue { return tsCall(loop, tsFunctionValue(profiler)) }
	profiler.properties = tsNewObject()
	profiler.properties.set("prototype", tsObjectValue(profilerPrototype))
	profilerPrototype.set("constructor", tsFunctionValue(profiler))
	profilerPrototype.set("start", tsNamedNativeMethod("start", 0, func(_ *tsLoop, receiver tsValue, _ ...tsValue) tsValue {
		if receiver.kind != tsObjectKind {
			panic(tsThrown{tsFSException("TypeError", "", "Illegal invocation")})
		}
		key := uintptr(receiver.ref)
		tsNodeV8Singleton.mu.Lock()
		state := tsNodeV8Singleton.profilers[key]
		if state == nil {
			tsNodeV8Singleton.mu.Unlock()
			panic(tsThrown{tsFSException("TypeError", "", "Illegal invocation")})
		}
		if !state.started {
			var memory runtime.MemStats
			runtime.ReadMemStats(&memory)
			state.started = true
			state.start = time.Now()
			state.gcCount = memory.NumGC
			state.startUsed = memory.HeapAlloc
			state.startCapacity = memory.HeapSys
		}
		tsNodeV8Singleton.mu.Unlock()
		return tsU
	}))
	stop := tsNamedNativeMethod("stop", 0, func(_ *tsLoop, receiver tsValue, _ ...tsValue) tsValue {
		if receiver.kind != tsObjectKind {
			panic(tsThrown{tsFSException("TypeError", "", "Illegal invocation")})
		}
		key := uintptr(receiver.ref)
		tsNodeV8Singleton.mu.Lock()
		state := tsNodeV8Singleton.profilers[key]
		if state == nil || !state.started {
			tsNodeV8Singleton.mu.Unlock()
			return tsU
		}
		start := state.start
		startUsed, startCapacity, startGC := state.startUsed, state.startCapacity, state.gcCount
		state.started = false
		delete(tsNodeV8Singleton.profilers, key)
		tsNodeV8Singleton.mu.Unlock()

		var memory runtime.MemStats
		runtime.ReadMemStats(&memory)
		statistics := &tsArray{}
		first := startGC + 1
		if memory.NumGC-startGC > 256 {
			first = memory.NumGC - 255
		}
		for cycle := first; cycle <= memory.NumGC && cycle > startGC; cycle++ {
			pause := memory.PauseNs[(cycle-1)%uint32(len(memory.PauseNs))]
			before := tsNewObject()
			before.set("totalHeapSize", tsNumberValue(float64(startCapacity)))
			before.set("totalHeapSizeExecutable", tsNumberValue(0))
			before.set("totalPhysicalSize", tsNumberValue(float64(startCapacity)))
			before.set("totalAvailableSize", tsNumberValue(float64(max(startCapacity-startUsed, 0))))
			before.set("usedHeapSize", tsNumberValue(float64(startUsed)))
			before.set("heapSizeLimit", tsNumberValue(0))
			before.set("mallocedMemory", tsNumberValue(0))
			before.set("peakMallocedMemory", tsNumberValue(0))
			before.set("externalMemory", tsNumberValue(0))
			before.set("totalGlobalHandlesSize", tsNumberValue(0))
			before.set("usedGlobalHandlesSize", tsNumberValue(0))
			after := tsNewObject()
			after.set("totalHeapSize", tsNumberValue(float64(memory.HeapSys)))
			after.set("totalHeapSizeExecutable", tsNumberValue(0))
			after.set("totalPhysicalSize", tsNumberValue(float64(memory.HeapSys)))
			after.set("totalAvailableSize", tsNumberValue(float64(max(memory.HeapSys-memory.HeapAlloc, 0))))
			after.set("usedHeapSize", tsNumberValue(float64(memory.HeapAlloc)))
			after.set("heapSizeLimit", tsNumberValue(0))
			after.set("mallocedMemory", tsNumberValue(0))
			after.set("peakMallocedMemory", tsNumberValue(0))
			after.set("externalMemory", tsNumberValue(0))
			after.set("totalGlobalHandlesSize", tsNumberValue(0))
			after.set("usedGlobalHandlesSize", tsNumberValue(0))
			event := tsNewObject()
			event.set("gcType", tsNodeV8Key("MarkSweepCompact"))
			event.set("cost", tsNumberValue(float64(pause)/1000))
			event.set("beforeGC", tsObjectValue(before))
			event.set("afterGC", tsObjectValue(after))
			statistics.values = append(statistics.values, tsObjectValue(event))
		}
		out := tsNewObject()
		out.set("version", tsNumberValue(1))
		out.set("startTime", tsNumberValue(float64(start.UnixMilli())))
		out.set("endTime", tsNumberValue(float64(time.Now().UnixMilli())))
		out.set("statistics", tsArrayValue(statistics))
		return tsObjectValue(out)
	})
	profilerPrototype.set("stop", stop)
	profilerPrototype.set(tsPropertyKey(tsWellKnownSymbol("dispose")), stop)
	module.set("GCProfiler", tsFunctionValue(profiler))
	module.set("getHeapSnapshot",tsNamedNativeMethod("getHeapSnapshot",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{stream:=tsNodeStreamModule(loop);readable:=tsGet(loop,stream,tsNodeV8Key("Readable"));from:=tsGet(loop,readable,tsNodeV8Key("from"));return tsCall(loop,from,tsArrayValue(&tsArray{values:[]tsValue{tsStringReference(tsStringUTF8(tsNodeV8HeapSnapshot()))}}))}))
	module.set("writeHeapSnapshot",tsNamedNativeMethod("writeHeapSnapshot",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{pathValue:=tsArg(args,0);if tsIsUndefined(pathValue){date:=time.Now();thread:=tsGet(loop,tsNodeWorkerThreadsModule(loop),tsNodeV8Key("threadId"));pathValue=tsNodeV8Key(fmt.Sprintf("Heap-%04d%02d%02d-%02d%02d%02d-%d-%s.heapsnapshot",date.Year(),int(date.Month())-1,date.Day(),date.Hour(),date.Minute(),date.Second(),os.Getpid(),tsText(thread)))};if pathValue.kind!=tsStringKind{tsFSInvalid("path","a string")};path:=tsText(pathValue);if path==""{tsFSInvalid("path","a non-empty string")};if err:=os.WriteFile(path,[]byte(tsNodeV8HeapSnapshot()),0644);err!=nil{tsNodeProcessError("Error",err)};return pathValue}))
	module.set("default",tsObjectValue(module));tsNodeV8Singleton.module=tsObjectValue(module)
	});return tsNodeV8Singleton.module
}
`
