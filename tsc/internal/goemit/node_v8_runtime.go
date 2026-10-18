package goemit

// nodeV8Runtime maps Bun's node:v8 wrapper onto Go heap counters and a
// versioned structured-value format. GCProfiler observes Go runtime GC cycles,
// not JSC heap collections; snapshots and serializer data are local formats.
const nodeV8Runtime = `
type tsNodeV8Edge struct { Key string; Target int; Group int; Role string }
type tsNodeV8Node struct { Kind string; Number uint64; Nonce uint64; Boolean bool; Units []uint16; Big string; Text string; Bytes []byte; Name string; Length int; Offset int; Backing int; HasBacking bool; Transfer uint32; HasTransfer bool; HasName bool; HasMessage bool; Buffer bool; Set bool; Edges []tsNodeV8Edge }
type tsNodeV8Document struct { Version int; Root int; Tagged bool; Nodes []tsNodeV8Node }
type tsNodeV8ProfilerState struct { started bool; start time.Time; gcCount uint32; startUsed, startCapacity uint64 }
type tsNodeV8SerializerState struct { buffer []byte; transfers map[unsafe.Pointer]uint32 }
type tsNodeV8DeserializerState struct { buffer []byte; offset int; transfers map[uint32]tsValue }
type tsNodeV8PromiseHookHandle struct { id uint64; stopped bool }
var tsNodeV8Singleton struct { sync.Once; module tsValue; mu sync.Mutex; profilers map[uintptr]*tsNodeV8ProfilerState; flags string; versionTag uint32; tagValid bool }
var tsNodeV8Magic=[]byte{0xff,'T','S','V',2}
var tsNodeV8MagicLegacy=[]byte{0xff,'T','S','V',1}
func tsNodeV8HeaderVersion(data []byte)(int,int){if bytes.HasPrefix(data,tsNodeV8Magic){return 2,len(tsNodeV8Magic)};if bytes.HasPrefix(data,tsNodeV8MagicLegacy){return 1,len(tsNodeV8MagicLegacy)};return 0,0}
func tsNodeV8Key(name string)tsValue{return tsStringReference(tsStringUTF8(name))}
func tsNodeV8BuiltinErrorName(name string)bool{switch name{case "Error","TypeError","RangeError","ReferenceError","SyntaxError","URIError","EvalError","AggregateError","SuppressedError":return true};return false}
func tsNodeV8Error(code,message string){panic(tsThrown{tsFSException("Error",code,message)})}
func tsNodeV8NotImplemented(name string){tsNodeV8Error("ERR_METHOD_NOT_IMPLEMENTED","node:v8 "+name+" is not implemented")}
func tsNodeV8HeapSnapshot()string{var memory runtime.MemStats;runtime.ReadMemStats(&memory);meta:=map[string]any{"node_fields":[]string{"type","name","id","self_size","edge_count","trace_node_id","detachedness"},"node_types":[][]string{{"hidden","array","string","object","code","closure","regexp","number","native","synthetic","concatenated string","sliced string","symbol","bigint"},{"string"},{"number"},{"number"},{"number"},{"number"},{"number"}},"edge_fields":[]string{"type","name_or_index","to_node"},"edge_types":[][]string{{"context","element","property","internal","hidden","shortcut","weak"},{"string_or_number"},{"node"}},"trace_function_info_fields":[]string{"function_id","name","script_name","script_id","line","column"},"trace_node_fields":[]string{"id","function_info_index","count","size","children"},"sample_fields":[]string{"timestamp_us","last_assigned_id"},"location_fields":[]string{"object_index","script_id","line","column"}};snapshot:=map[string]any{"snapshot":map[string]any{"meta":meta,"node_count":1,"edge_count":0,"trace_function_count":0,"trace_node_count":0,"sample_count":0,"location_count":0,"root_index":0},"nodes":[]uint64{9,0,1,uint64(memory.HeapAlloc),0,0,0},"edges":[]uint64{},"strings":[]string{"synthetic Go heap snapshot"},"trace_function_infos":[]uint64{},"trace_tree":[]uint64{},"samples":[]uint64{},"locations":[]uint64{}};data,_:=json.Marshal(snapshot);return string(data)}
func tsNodeV8HeapStats(loop *tsLoop)tsValue{var mem runtime.MemStats;runtime.ReadMemStats(&mem);total:=uint64(mem.Sys);if osinfo:=tsNodeOSModule();!tsIsUndefined(tsGet(osinfo,tsNodeV8Key("totalmem"))){fn:=tsGet(loop,osinfo,tsNodeV8Key("totalmem"));if fn.kind==tsFunctionKind{candidate:=tsCall(loop,fn);if tsIsNumeric(candidate)&&tsNumber(candidate)>0{total=uint64(tsNumber(candidate))}}};out:=tsNewObject();values:=[]struct{name string;value uint64}{{"total_heap_size",mem.HeapSys},{"total_heap_size_executable",0},{"total_physical_size",mem.Sys},{"total_available_size",max(total-uint64(mem.HeapAlloc),0)},{"used_heap_size",mem.HeapAlloc},{"total_allocated_bytes",mem.TotalAlloc},{"heap_size_limit",total},{"malloced_memory",mem.HeapAlloc},{"peak_malloced_memory",mem.TotalAlloc},{"does_zap_garbage",0},{"number_of_native_contexts",1},{"number_of_detached_contexts",0},{"total_global_handles_size",8192},{"used_global_handles_size",2208},{"external_memory",0}};for _,entry:=range values{out.set(entry.name,tsNumberValue(float64(entry.value)))};return tsObjectValue(out)}
func tsNodeV8HeapSpaceStats(loop *tsLoop)tsValue{names:=[]string{"read_only_space","new_space","old_space","code_space","shared_space","trusted_space","new_large_object_space","large_object_space","code_large_object_space","shared_large_object_space","shared_trusted_space","shared_trusted_large_object_space","trusted_large_object_space"};var mem runtime.MemStats;runtime.ReadMemStats(&mem);out:=&tsArray{};for _,name:=range names{size,used:=uint64(0),uint64(0);if name=="old_space"{size=mem.HeapSys;used=mem.HeapAlloc};item:=tsNewObject();item.set("space_name",tsNodeV8Key(name));item.set("space_size",tsNumberValue(float64(size)));item.set("space_used_size",tsNumberValue(float64(used)));item.set("space_available_size",tsNumberValue(float64(max(size-used,0))));item.set("physical_space_size",tsNumberValue(float64(size)));out.values=append(out.values,tsObjectValue(item))};return tsArrayValue(out)}
func tsNodeV8Encode(loop *tsLoop,value tsValue) []byte {return tsNodeV8EncodeWithTransfers(loop,value,nil)}
func tsNodeV8EncodeWithTransfers(loop *tsLoop,value tsValue,transferIDs map[unsafe.Pointer]uint32) []byte {
	tagged:=tsInternalSerializationBuffersTag(value);root:=value;hasBuffers:=tagged.kind==tsArrayKind;if hasBuffers{root=tagged}
	doc:=tsNodeV8Document{Version:2,Tagged:hasBuffers};refs:=map[unsafe.Pointer]int{}
	var visit func(tsValue)(int,error)
	visit=func(item tsValue)(int,error){node:=tsNodeV8Node{};switch item.kind{
	case tsUndefinedKind:node.Kind="undefined"
	case tsNullKind:node.Kind="null"
	case tsBooleanKind:node.Kind="boolean";node.Boolean=item.number!=0
	case tsNumberKind,tsFloat32Kind,tsIntKind,tsInt8Kind,tsInt16Kind,tsInt32Kind,tsInt64Kind,tsUintKind,tsUint8Kind,tsUint16Kind,tsUint32Kind,tsUint64Kind:node.Kind="number";node.Number=math.Float64bits(tsNumber(item))
	case tsStringKind:node.Kind="string";node.Units=append([]uint16(nil),(*tsString)(item.ref).units...)
	case tsBigIntKind:node.Kind="bigint";node.Big=(*big.Int)(item.ref).String()
		case tsRegExpKind:if id,ok:=refs[item.ref];ok{return id,nil};regexp:=(*tsRegExp)(item.ref);node.Kind="regexp";node.Units=append([]uint16(nil),regexp.source.units...);node.Text=regexp.flags
	case tsSymbolKind,tsFunctionKind,tsClassKind:return 0,fmt.Errorf("Value of type %s cannot be serialized",tsValueTypeName(item))
	case tsArrayBufferKind:if id,ok:=refs[item.ref];ok{return id,nil};node.Kind="arraybuffer";if transferID,ok:=transferIDs[item.ref];ok{node.HasTransfer=true;node.Transfer=transferID}else{node.Bytes=append([]byte(nil),(*tsArrayBuffer)(item.ref).data...)}
		case tsTypedArrayKind:if id,ok:=refs[item.ref];ok{return id,nil};array:=(*tsTypedArray)(item.ref);node.Kind="typedarray";node.Name=array.name;node.Buffer=array.nodeBuffer;node.Length=array.length;node.Offset=array.offset;backing,err:=visit(tsValue{ref:unsafe.Pointer(array.buffer),kind:tsArrayBufferKind});if err!=nil{return 0,err};node.Backing=backing;node.HasBacking=true
		case tsDataViewKind:if id,ok:=refs[item.ref];ok{return id,nil};view:=(*tsDataView)(item.ref);node.Kind="dataview";node.Length=view.length;node.Offset=view.offset;backing,err:=visit(tsValue{ref:unsafe.Pointer(view.buffer),kind:tsArrayBufferKind});if err!=nil{return 0,err};node.Backing=backing;node.HasBacking=true
	case tsCollectionKind:if id,ok:=refs[item.ref];ok{return id,nil};node.Kind="collection";node.Set=(*tsCollection)(item.ref).set
	case tsArrayKind:if id,ok:=refs[item.ref];ok{return id,nil};node.Kind="array";node.Length=(*tsArray)(item.ref).length()
		case tsErrorKind:if id,ok:=refs[item.ref];ok{return id,nil};error:=(*tsRuntimeError)(item.ref);node.Kind="error";node.Name=error.name;node.Text=error.message;node.HasName=!tsNodeV8BuiltinErrorName(error.name);node.HasMessage=true
		case tsObjectKind,tsInstanceKind:if id,ok:=refs[item.ref];ok{return id,nil};if port:=tsNodeWorkerPortFor(item);port!=nil{if transferID,ok:=transferIDs[item.ref];ok{node.Kind="messageport";node.HasTransfer=true;node.Transfer=transferID}else{return 0,fmt.Errorf("MessagePort is not in the transfer list")}}else if nonce,ok:=tsNodeNetBlockListSerializationNonce(item);ok{node.Kind="blocklist";node.Nonce=nonce}else{node.Kind="object";if item.kind==tsObjectKind&&(*tsObject)(item.ref).nativeDate{node.Kind="date";node.Number=math.Float64bits((*tsObject)(item.ref).nativeTime)}else if item.kind==tsObjectKind&&(*tsObject)(item.ref).nativeError{node.Kind="error";name:=tsGet(loop,item,tsNodeV8Key("name"));message:=tsGet(loop,item,tsNodeV8Key("message"));if tsIsUndefined(name){node.Name="Error"}else{node.Name=tsText(name)};node.HasName=tsObjectHasOwn(item,"name");if !tsIsUndefined(message){node.Text=tsText(message);node.HasMessage=tsObjectHasOwn(item,"message")}}}
	default:return 0,fmt.Errorf("Value of type %s cannot be serialized",tsValueTypeName(item))
	};id:=len(doc.Nodes);doc.Nodes=append(doc.Nodes,node);if item.kind==tsArrayBufferKind||item.kind==tsTypedArrayKind||item.kind==tsDataViewKind||item.kind==tsCollectionKind||item.kind==tsArrayKind||item.kind==tsObjectKind||item.kind==tsInstanceKind||item.kind==tsErrorKind||item.kind==tsRegExpKind{refs[item.ref]=id}
		edges:=[]tsNodeV8Edge{};switch item.kind{
		case tsArrayKind:a:=(*tsArray)(item.ref);for index:=0;index<a.length();index++{if a.holes[index]{continue};child,err:=visit(a.at(index));if err!=nil{return 0,err};edges=append(edges,tsNodeV8Edge{Key:strconv.Itoa(index),Target:child,Role:"property"})};for _,name:=range a.order{childValue,ok:=a.properties[name];if !ok{continue};if d:=a.descriptors[name];d!=nil&&(d.accessor||!d.enumerable){continue};child,err:=visit(childValue);if err!=nil{return 0,err};edges=append(edges,tsNodeV8Edge{Key:name,Target:child,Role:"property"})}
		case tsCollectionKind:c:=(*tsCollection)(item.ref);group:=0;for _,entry:=range c.entries{if !entry.live{continue};key,err:=visit(entry.key);if err!=nil{return 0,err};edges=append(edges,tsNodeV8Edge{Target:key,Group:group,Role:"key"});if !c.set{child,err:=visit(entry.value);if err!=nil{return 0,err};edges=append(edges,tsNodeV8Edge{Target:child,Group:group,Role:"value"})};group++}
		case tsObjectKind,tsInstanceKind,tsErrorKind:if node.Kind!="date"&&node.Kind!="blocklist"&&node.Kind!="messageport"{keys:=tsEnumerableOwnKeys(loop,item);for _,key:=range keys.values{desc:=tsOwnDescriptor(loop,item,tsText(key));if desc.kind==tsObjectKind{d:=(*tsObject)(desc.ref);if _,accessor:=d.values["get"];accessor{continue};childValue:=d.values["value"];child,err:=visit(childValue);if err!=nil{return 0,err};edges=append(edges,tsNodeV8Edge{Key:tsText(key),Target:child,Role:"property"})}};if node.Kind=="error"{for _,name:=range []string{"cause","errors"}{if !tsIsUndefined(tsOwnDescriptor(loop,item,name)){child,err:=visit(tsGet(loop,item,tsNodeV8Key(name)));if err!=nil{return 0,err};edges=append(edges,tsNodeV8Edge{Key:name,Target:child,Role:"property"})}}}}
		};doc.Nodes[id].Edges=edges;return id,nil
	}
	rootID,err:=visit(root);if err!=nil{tsNodeV8Error("ERR_DATA_CLONE_ERROR",err.Error())};doc.Root=rootID;body,err:=json.Marshal(doc);if err!=nil{tsNodeV8Error("ERR_BUFFER_TOO_LARGE",err.Error())};out:=append([]byte(nil),tsNodeV8Magic...);out=append(out,body...);return out
}
func tsNodeV8Decode(loop *tsLoop,data []byte)tsValue{return tsNodeV8DecodeWithTransfers(loop,data,nil)}
func tsNodeV8DecodeWithTransfers(loop *tsLoop,data []byte,transfers map[uint32]tsValue)tsValue{
	headerVersion,headerLength:=tsNodeV8HeaderVersion(data);if headerVersion==0{tsNodeV8Error("ERR_BUFFER_UNSUPPORTED","invalid node:v8 serialization header");return tsU}
	var doc tsNodeV8Document
	if err:=json.Unmarshal(data[headerLength:],&doc);err!=nil{tsNodeV8Error("ERR_BUFFER_UNSUPPORTED",err.Error());return tsU}
	if doc.Version!=headerVersion||doc.Root<0||doc.Root>=len(doc.Nodes)||len(doc.Nodes)>1<<20{tsNodeV8Error("ERR_BUFFER_UNSUPPORTED","invalid serialized value graph");return tsU}
	values:=make([]tsValue,len(doc.Nodes))
	for index,node:=range doc.Nodes{
		switch node.Kind{
		case "undefined":values[index]=tsU
		case "null":values[index]=tsNull
		case "boolean":values[index]=tsBooleanValue(node.Boolean)
		case "number":values[index]=tsNumberValue(math.Float64frombits(node.Number))
		case "string":values[index]=tsStringReference(tsStringUnits(node.Units))
		case "bigint":number,ok:=new(big.Int).SetString(node.Big,10);if !ok{tsNodeV8Error("ERR_BUFFER_UNSUPPORTED","invalid bigint");return tsU};values[index]=tsBigIntValue(number)
		case "regexp":value:=tsConstruct(tsRegExpConstructor(),tsStringReference(tsStringUnits(node.Units)),tsNodeV8Key(node.Text));if value.kind!=tsRegExpKind{tsNodeV8Error("ERR_BUFFER_UNSUPPORTED","invalid serialized RegExp");return tsU};values[index]=value
		case "object":values[index]=tsObjectValue(tsNewObject())
		case "error":message:=tsU;if node.HasMessage||node.Text!=""{message=tsNodeV8Key(node.Text)};if node.Name=="AggregateError"{values[index]=tsNewError(loop,node.Name,[]tsValue{tsArrayValue(&tsArray{}),message})}else{values[index]=tsNewError(loop,node.Name,[]tsValue{message})};if node.HasName||!tsNodeV8BuiltinErrorName(node.Name){tsSet(loop,values[index],tsNodeV8Key("name"),tsNodeV8Key(node.Name))}
		case "blocklist":if value,ok:=tsNodeNetBlockListFromSerializationNonce(node.Nonce);ok{values[index]=value}else{tsNodeV8Error("ERR_BUFFER_UNSUPPORTED","BlockList structured-clone reference is no longer available");return tsU}
		case "messageport":value,ok:=transfers[node.Transfer];if !ok||tsNodeWorkerPortFor(value)==nil{tsNodeV8Error("ERR_BUFFER_UNSUPPORTED","missing transferred MessagePort");return tsU};values[index]=value
		case "date":values[index]=tsFSDateMillis(math.Float64frombits(node.Number))
		case "array":if node.Length<0||node.Length>len(data){tsNodeV8Error("ERR_BUFFER_UNSUPPORTED","invalid serialized array length");return tsU};values[index]=tsArrayValue(&tsArray{values:make([]tsValue,node.Length),holes:map[int]bool{}})
		case "collection":values[index]=tsCollectionValue(&tsCollection{set:node.Set,index:map[tsCollectionKey]int{}})
		case "arraybuffer":
			if node.HasTransfer{value,ok:=transfers[node.Transfer];if !ok||value.kind!=tsArrayBufferKind{tsNodeV8Error("ERR_BUFFER_UNSUPPORTED","missing transferred ArrayBuffer");return tsU};values[index]=value}else{values[index]=tsArrayBufferValue(&tsArrayBuffer{data:append([]byte(nil),node.Bytes...),maxByteLength:len(node.Bytes)})}
			case "dataview":
				var buffer *tsArrayBuffer
				if node.HasBacking{if node.Backing<0||node.Backing>=len(values)||values[node.Backing].kind!=tsArrayBufferKind{tsNodeV8Error("ERR_BUFFER_UNSUPPORTED","invalid serialized DataView backing buffer");return tsU};buffer=(*tsArrayBuffer)(values[node.Backing].ref)}else if node.HasTransfer{value,ok:=transfers[node.Transfer];if !ok||value.kind!=tsArrayBufferKind{tsNodeV8Error("ERR_BUFFER_UNSUPPORTED","missing transferred ArrayBuffer");return tsU};buffer=(*tsArrayBuffer)(value.ref)}else{buffer=&tsArrayBuffer{data:append([]byte(nil),node.Bytes...),maxByteLength:len(node.Bytes)}}
			offset,length:=node.Offset,node.Length;if offset<0||length<0||offset>len(buffer.data)||length>len(buffer.data)-offset{tsNodeV8Error("ERR_BUFFER_UNSUPPORTED","invalid serialized DataView bounds");return tsU};values[index]=tsDataViewValue(&tsDataView{buffer:buffer,offset:offset,length:length,properties:map[string]tsValue{}})
		case "typedarray":
				bufferValue:=tsU;offset:=node.Offset
				if node.HasBacking{if node.Backing<0||node.Backing>=len(values)||values[node.Backing].kind!=tsArrayBufferKind{tsNodeV8Error("ERR_BUFFER_UNSUPPORTED","invalid serialized TypedArray backing buffer");return tsU};bufferValue=values[node.Backing]}else if node.HasTransfer{bufferValue=transfers[node.Transfer];if tsIsUndefined(bufferValue)||bufferValue.kind!=tsArrayBufferKind{tsNodeV8Error("ERR_BUFFER_UNSUPPORTED","missing transferred ArrayBuffer");return tsU}}else{bufferValue=tsArrayBufferValue(&tsArrayBuffer{data:append([]byte(nil),node.Bytes...),maxByteLength:len(node.Bytes)})}
			name:=node.Name;if name==""{name="Uint8Array"};if node.Length<0{tsNodeV8Error("ERR_BUFFER_UNSUPPORTED","invalid serialized TypedArray length");return tsU};arrayValue:=tsTypedArrayBufferExtended(name,[]tsValue{bufferValue,tsNumberValue(float64(offset)),tsNumberValue(float64(node.Length))});if arrayValue.kind!=tsTypedArrayKind{tsNodeV8Error("ERR_BUFFER_UNSUPPORTED","invalid typed array");return tsU};(*tsTypedArray)(arrayValue.ref).nodeBuffer=node.Buffer;values[index]=arrayValue
		default:tsNodeV8Error("ERR_BUFFER_UNSUPPORTED","unknown serialized value kind");return tsU
		}
	}
	for index,node:=range doc.Nodes{
		target:=values[index]
		for _,edge:=range node.Edges{if edge.Target<0||edge.Target>=len(values){tsNodeV8Error("ERR_BUFFER_UNSUPPORTED","invalid serialized value reference");return tsU}}
		switch node.Kind{
		case "object","error":object:=(*tsObject)(target.ref);for _,edge:=range node.Edges{object.set(edge.Key,values[edge.Target]);if node.Kind=="error"&&(edge.Key=="cause"||edge.Key=="errors"){if object.descriptors==nil{object.descriptors=map[string]*tsDescriptor{}};object.descriptors[edge.Key]=&tsDescriptor{writable:true,configurable:true}}}
		case "array":array:=(*tsArray)(target.ref);for i:=0;i<array.length();i++{array.holes[i]=true};for _,edge:=range node.Edges{if edge.Role!="property"{continue};if offset,err:=strconv.Atoi(edge.Key);err==nil&&strconv.Itoa(offset)==edge.Key&&offset>=0&&offset<array.length(){array.values[offset]=values[edge.Target];delete(array.holes,offset)}else{if array.properties==nil{array.properties=map[string]tsValue{}};array.order=append(array.order,edge.Key);array.properties[edge.Key]=values[edge.Target]}}
		case "collection":collection:=(*tsCollection)(target.ref);pending:=map[int]tsValue{};for _,edge:=range node.Edges{if edge.Role=="key"{pending[edge.Group]=values[edge.Target]}else if edge.Role=="value"{key:=pending[edge.Group];collection.put(loop,key,values[edge.Target]);delete(pending,edge.Group)}};for _,key:=range pending{collection.put(loop,key,key)}
		}
	}
	if doc.Tagged{return tsInternalSerializationBuffersRestore(values[doc.Root])};return values[doc.Root]
}
func tsNodeV8ByteView(value tsValue)[]byte{switch value.kind{case tsArrayBufferKind,tsTypedArrayKind,tsDataViewKind:return tsNodeBufferByteView(value);default:tsNodeV8Error("ERR_INVALID_ARG_TYPE","The value argument must be an ArrayBuffer or an ArrayBufferView");return nil}}
func tsNodeV8SerializerReceiver(receiver tsValue)*tsNodeV8SerializerState{if receiver.kind!=tsObjectKind{tsPropertyFailure("Serializer method called on an incompatible receiver")};state,ok:=(*tsObject)(receiver.ref).nativeData.(*tsNodeV8SerializerState);if !ok||state==nil{tsPropertyFailure("Serializer method called on an incompatible receiver")};return state}
func tsNodeV8DeserializerReceiver(receiver tsValue)*tsNodeV8DeserializerState{if receiver.kind!=tsObjectKind{tsPropertyFailure("Deserializer method called on an incompatible receiver")};state,ok:=(*tsObject)(receiver.ref).nativeData.(*tsNodeV8DeserializerState);if !ok||state==nil{tsPropertyFailure("Deserializer method called on an incompatible receiver")};return state}
func tsNodeV8AppendUint32(buffer []byte,value uint32)[]byte{return append(buffer,byte(value>>24),byte(value>>16),byte(value>>8),byte(value))}
func tsNodeV8ReadFixed(state *tsNodeV8DeserializerState,count int)[]byte{if count<0||state.offset<0||state.offset>len(state.buffer)||count>len(state.buffer)-state.offset{tsNodeV8Error("ERR_BUFFER_OUT_OF_BOUNDS","Attempt to access memory outside buffer bounds");return nil};start:=state.offset;state.offset+=count;return state.buffer[start:state.offset]}
func tsNodeV8Uint32(value tsValue,name string)uint32{if !tsIsNumeric(value){tsInternalValidatorInvalidType(name,"a number",value)};number:=tsNumber(value);if math.IsNaN(number)||math.IsInf(number,0)||math.Trunc(number)!=number||number<0||number>4294967295{tsInternalValidatorOutOfRange(name,"an unsigned 32-bit integer",value)};return uint32(number)}
func tsNodeV8SerializerMethods(prototype *tsObject){
	prototype.set("writeHeader",tsNamedNativeMethod("writeHeader",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeV8SerializerReceiver(receiver);state.buffer=append(state.buffer,tsNodeV8Magic...);return tsU}))
	prototype.set("writeValue",tsNamedNativeMethod("writeValue",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeV8SerializerReceiver(receiver);payload:=tsNodeV8EncodeWithTransfers(loop,tsArg(args,0),state.transfers);if uint64(len(payload))>uint64(^uint32(0)){tsNodeV8Error("ERR_BUFFER_TOO_LARGE","Serialized value exceeds the maximum size");return tsU};state.buffer=tsNodeV8AppendUint32(state.buffer,uint32(len(payload)));state.buffer=append(state.buffer,payload...);return tsBooleanValue(true)}))
	prototype.set("releaseBuffer",tsNamedNativeMethod("releaseBuffer",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeV8SerializerReceiver(receiver);out:=append([]byte(nil),state.buffer...);state.buffer=nil;return tsNodeBuffer(out)}))
	prototype.set("transferArrayBuffer",tsNamedNativeMethod("transferArrayBuffer",2,func(_ *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeV8SerializerReceiver(receiver);id:=tsNodeV8Uint32(tsArg(args,0),"id");value:=tsArg(args,1);if value.kind!=tsArrayBufferKind{tsInternalValidatorInvalidType("arrayBuffer","an ArrayBuffer",value)};if state.transfers==nil{state.transfers=map[unsafe.Pointer]uint32{}};state.transfers[value.ref]=id;return tsU}))
	prototype.set("writeUint32",tsNamedNativeMethod("writeUint32",1,func(_ *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeV8SerializerReceiver(receiver);state.buffer=tsNodeV8AppendUint32(state.buffer,tsNodeV8Uint32(tsArg(args,0),"value"));return tsU}))
	prototype.set("writeUint64",tsNamedNativeMethod("writeUint64",1,func(_ *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeV8SerializerReceiver(receiver);value:=tsArg(args,0);if value.kind!=tsBigIntKind{tsInternalValidatorInvalidType("value","a bigint",value)};number:=(*big.Int)(value.ref);if number.Sign()<0||number.BitLen()>64{tsInternalValidatorOutOfRange("value","an unsigned 64-bit integer",value)};n:=number.Uint64();for shift:=56;shift>=0;shift-=8{state.buffer=append(state.buffer,byte(n>>uint(shift)))};return tsU}))
	prototype.set("writeDouble",tsNamedNativeMethod("writeDouble",1,func(_ *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeV8SerializerReceiver(receiver);value:=tsArg(args,0);if !tsIsNumeric(value){tsInternalValidatorInvalidType("value","a number",value)};bits:=math.Float64bits(tsNumber(value));for shift:=56;shift>=0;shift-=8{state.buffer=append(state.buffer,byte(bits>>uint(shift)))};return tsU}))
	prototype.set("writeRawBytes",tsNamedNativeMethod("writeRawBytes",1,func(_ *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeV8SerializerReceiver(receiver);state.buffer=append(state.buffer,tsNodeV8ByteView(tsArg(args,0))...);return tsU}))
}
func tsNodeV8DeserializerMethods(prototype *tsObject){
	prototype.set("readHeader",tsNamedNativeMethod("readHeader",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeV8DeserializerReceiver(receiver);if state.offset<0||state.offset>len(state.buffer){return tsBooleanValue(false)};_,length:=tsNodeV8HeaderVersion(state.buffer[state.offset:]);if length==0{return tsBooleanValue(false)};state.offset+=length;return tsBooleanValue(true)}))
	prototype.set("readValue",tsNamedNativeMethod("readValue",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeV8DeserializerReceiver(receiver);prefix:=tsNodeV8ReadFixed(state,4);length:=int(prefix[0])<<24|int(prefix[1])<<16|int(prefix[2])<<8|int(prefix[3]);payload:=tsNodeV8ReadFixed(state,length);return tsNodeV8DecodeWithTransfers(loop,payload,state.transfers)}))
	prototype.set("transferArrayBuffer",tsNamedNativeMethod("transferArrayBuffer",2,func(_ *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeV8DeserializerReceiver(receiver);id:=tsNodeV8Uint32(tsArg(args,0),"id");value:=tsArg(args,1);if value.kind!=tsArrayBufferKind{tsInternalValidatorInvalidType("arrayBuffer","an ArrayBuffer",value)};if state.transfers==nil{state.transfers=map[uint32]tsValue{}};state.transfers[id]=value;return tsU}))
	prototype.set("getWireFormatVersion",tsNamedNativeMethod("getWireFormatVersion",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{_ = tsNodeV8DeserializerReceiver(receiver);return tsNumberValue(1)}))
	prototype.set("readUint32",tsNamedNativeMethod("readUint32",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeV8DeserializerReceiver(receiver);data:=tsNodeV8ReadFixed(state,4);value:=uint32(data[0])<<24|uint32(data[1])<<16|uint32(data[2])<<8|uint32(data[3]);return tsNumberValue(float64(value))}))
	prototype.set("readUint64",tsNamedNativeMethod("readUint64",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeV8DeserializerReceiver(receiver);data:=tsNodeV8ReadFixed(state,8);var value uint64;for _,b:=range data{value=value<<8|uint64(b)};return tsBigIntValue(new(big.Int).SetUint64(value))}))
	prototype.set("readDouble",tsNamedNativeMethod("readDouble",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeV8DeserializerReceiver(receiver);data:=tsNodeV8ReadFixed(state,8);var value uint64;for _,b:=range data{value=value<<8|uint64(b)};return tsNumberValue(math.Float64frombits(value))}))
	prototype.set("readRawBytes",tsNamedNativeMethod("readRawBytes",1,func(_ *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeV8DeserializerReceiver(receiver);length:=tsNodeV8Uint32(tsArg(args,0),"length");return tsNodeBuffer(tsNodeV8ReadFixed(state,int(length)))}))
}
func tsNodeV8Ctor(name string,prototype *tsObject,deserializer bool)tsValue{constructor:=tsFunc(func(_ *tsLoop,_ ...tsValue)tsValue{tsPropertyFailure("Class constructor "+name+" cannot be invoked without 'new'");return tsU});constructor.name=name;constructor.length=0;constructor.constructible=true;constructor.construct=func(_ *tsLoop,args ...tsValue)tsValue{object:=tsNewObject();object.prototype=tsObjectValue(prototype);if deserializer{input:=tsNodeV8ByteView(tsArg(args,0));object.nativeData=&tsNodeV8DeserializerState{buffer:append([]byte(nil),input...)}}else{object.nativeData=&tsNodeV8SerializerState{}};return tsObjectValue(object)};constructor.properties=tsNewObject();constructor.properties.set("prototype",tsObjectValue(prototype));prototype.set("constructor",tsFunctionValue(constructor));return tsFunctionValue(constructor)}
func tsNodeV8PromiseHookStop(id uint64){tsPromiseHookStop(id)}
func tsNodeV8PromiseHookStopFunction(handle *tsNodeV8PromiseHookHandle)tsValue{stop:=tsNamedNativeMethod("stop",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{if !handle.stopped{handle.stopped=true;tsNodeV8PromiseHookStop(handle.id)};return tsU});return stop}
func tsNodeV8PromiseHookCallbacks(loop *tsLoop,value tsValue)tsPromiseHookCallbacks{if value.kind!=tsObjectKind&&value.kind!=tsInstanceKind{tsInternalValidatorInvalidType("callbacks","an object",value)};callbacks:=tsPromiseHookCallbacks{};for _,entry:=range []struct{name string;target *tsValue}{{"init",&callbacks.init},{"before",&callbacks.before},{"after",&callbacks.after},{"settled",&callbacks.settled}}{item:=tsGet(loop,value,tsNodeV8Key(entry.name));if !tsIsUndefined(item)&&item.kind!=tsFunctionKind&&item.kind!=tsClassKind{tsInternalValidatorInvalidType(entry.name,"a function",item)};if item.kind==tsClassKind{tsInternalValidatorInvalidType(entry.name,"a function",item)};*entry.target=item};return callbacks}
func tsNodeV8PromiseHooks(loop *tsLoop)tsValue{hooks:=tsNewObject();hooks.set("createHook",tsNamedNativeMethod("createHook",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{callbacks:=tsNodeV8PromiseHookCallbacks(loop,tsArg(args,0));handle:=&tsNodeV8PromiseHookHandle{id:tsPromiseHookRegister(callbacks)};object:=tsNewObject();object.set("stop",tsNodeV8PromiseHookStopFunction(handle));return tsObjectValue(object)}));for _,entry:=range []struct{name string;field string}{{"onInit","init"},{"onBefore","before"},{"onAfter","after"},{"onSettled","settled"}}{name,field:=entry.name,entry.field;hooks.set(name,tsNamedNativeMethod(name,1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{callback:=tsArg(args,0);if callback.kind!=tsFunctionKind{tsInternalValidatorInvalidType("callback","a function",callback)};callbacks:=tsPromiseHookCallbacks{};switch field{case "init":callbacks.init=callback;case "before":callbacks.before=callback;case "after":callbacks.after=callback;case "settled":callbacks.settled=callback};handle:=&tsNodeV8PromiseHookHandle{id:tsPromiseHookRegister(callbacks)};return tsNodeV8PromiseHookStopFunction(handle)}))};return tsObjectValue(hooks)}
func tsNodeV8Module(loop *tsLoop)tsValue{tsNodeV8Singleton.Do(func(){module:=tsNewObject();
	module.set("cachedDataVersionTag",tsNamedNativeMethod("cachedDataVersionTag",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{if !tsNodeV8Singleton.tagValid{tsNodeV8Singleton.versionTag=crc32.ChecksumIEEE([]byte(runtime.Version()+"/"+runtime.GOOS+"/"+runtime.GOARCH+tsNodeV8Singleton.flags));tsNodeV8Singleton.tagValid=true};return tsNumberValue(float64(tsNodeV8Singleton.versionTag))}))
	serializerPrototype:=tsNewObject();deserializerPrototype:=tsNewObject();tsNodeV8SerializerMethods(serializerPrototype);tsNodeV8DeserializerMethods(deserializerPrototype);serializer:=tsNodeV8Ctor("Serializer",serializerPrototype,false);deserializer:=tsNodeV8Ctor("Deserializer",deserializerPrototype,true);defaultSerializerPrototype:=tsNewObject();defaultSerializerPrototype.prototype=tsObjectValue(serializerPrototype);defaultDeserializerPrototype:=tsNewObject();defaultDeserializerPrototype.prototype=tsObjectValue(deserializerPrototype);module.set("Deserializer",deserializer);module.set("Serializer",serializer);module.set("DefaultDeserializer",tsNodeV8Ctor("DefaultDeserializer",defaultDeserializerPrototype,true));module.set("DefaultSerializer",tsNodeV8Ctor("DefaultSerializer",defaultSerializerPrototype,false))
	isOneByte:=tsNamedNativeMethod("isStringOneByteRepresentation",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{value:=tsArg(args,0);if value.kind!=tsStringKind{tsFSInvalid("content","a string")};for _,unit:=range (*tsString)(value.ref).units{if unit>255{return tsBooleanValue(false)}};return tsBooleanValue(true)});module.set("isStringOneByteRepresentation",isOneByte)
	module.set("getHeapStatistics",tsNamedNativeMethod("getHeapStatistics",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNodeV8HeapStats(loop)}));module.set("getHeapSpaceStatistics",tsNamedNativeMethod("getHeapSpaceStatistics",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNodeV8HeapSpaceStats(loop)}))
	codeStats:=tsNewObject();for _,name:=range []string{"code_and_metadata_size","bytecode_and_metadata_size","external_script_source_size","cpu_profiler_metadata_size"}{codeStats.set(name,tsNumberValue(0))};module.set("getHeapCodeStatistics",tsNamedNativeMethod("getHeapCodeStatistics",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsObjectValue(codeStats)}))
	module.set("getCppHeapStatistics",tsNamedNativeMethod("getCppHeapStatistics",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{kind:=tsArg(args,0);if tsIsUndefined(kind){kind=tsNodeV8Key("detailed")};name:=tsInternalValidatorsModule(loop);validate:=tsGet(loop,name,tsNodeV8Key("validateOneOf"));if validate.kind==tsFunctionKind{values:=tsArrayValue(&tsArray{values:[]tsValue{tsNodeV8Key("brief"),tsNodeV8Key("detailed")}});tsCall(loop,validate,kind,tsNodeV8Key("type"),values)};out:=tsNewObject();for _,key:=range []string{"committed_size_bytes","resident_size_bytes","used_size_bytes"}{out.set(key,tsNumberValue(0))};out.set("space_statistics",tsArrayValue(&tsArray{}));out.set("type_names",tsArrayValue(&tsArray{}));out.set("detail_level",kind);return tsObjectValue(out)}))
	module.set("setFlagsFromString",tsNamedNativeMethod("setFlagsFromString",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{value:=tsArg(args,0);if value.kind!=tsStringKind{tsFSInvalid("flags","a string")};tsNodeV8Singleton.flags+=" "+tsText(value);tsNodeV8Singleton.tagValid=false;return tsU}))
	module.set("serialize",tsNamedNativeMethod("serialize",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeBuffer(tsNodeV8Encode(loop,tsArg(args,0)))}))
	module.set("deserialize",tsNamedNativeMethod("deserialize",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{value:=tsArg(args,0);switch value.kind{case tsTypedArrayKind,tsDataViewKind,tsArrayBufferKind:return tsNodeV8Decode(loop,tsNodeBufferByteView(value));default:tsNodeV8Error("ERR_INVALID_ARG_TYPE","The value argument must be a Buffer, TypedArray, DataView, or ArrayBuffer");return tsU}}))
	for _,name:=range []string{"takeCoverage","stopCoverage","setHeapSnapshotNearHeapLimit"}{method:=name;module.set(method,tsNamedNativeMethod(method,0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{tsNodeV8NotImplemented(method);return tsU}))}
	module.set("promiseHooks",tsNodeV8PromiseHooks(loop))
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
		tsNodeV8Singleton.profilers[uintptr(unsafe.Pointer(object))] = &tsNodeV8ProfilerState{}
		tsNodeV8Singleton.mu.Unlock()
		runtime.SetFinalizer(object, func(collected *tsObject) {
			tsNodeV8Singleton.mu.Lock()
			delete(tsNodeV8Singleton.profilers, uintptr(unsafe.Pointer(collected)))
			tsNodeV8Singleton.mu.Unlock()
		})
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
		state.start = time.Time{}
		state.startUsed, state.startCapacity, state.gcCount = 0, 0, 0
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
