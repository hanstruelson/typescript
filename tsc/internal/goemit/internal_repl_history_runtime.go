package goemit

// internalReplHistoryRuntime adapts readline history management to local file
// IO and the compiler event loop.
const internalReplHistoryRuntime = `
type tsInternalReplHistoryState struct {
	context, history tsValue
	size, index int
	filePath string
	file *os.File
	timer *tsTask
	isFlushing, writing, pending, removeDuplicates bool
	originalHistoryPrev tsValue
}
var tsInternalReplHistory struct { sync.Once; module tsValue; constructor tsValue; prototype *tsObject; states map[uintptr]*tsInternalReplHistoryState }

func tsInternalReplHistoryStateFor(receiver tsValue) *tsInternalReplHistoryState {
	if receiver.kind!=tsObjectKind{return nil};return tsInternalReplHistory.states[uintptr(receiver.ref)]
}
func tsInternalReplHistoryItems(state *tsInternalReplHistoryState) *tsArray {
	if state.history.kind!=tsArrayKind{state.history=tsArrayValue(&tsArray{})};return (*tsArray)(state.history.ref)
}
func tsInternalReplHistoryGet(loop *tsLoop, state *tsInternalReplHistoryState, name string) tsValue { return tsGet(loop,state.context,tsStringReference(tsStringUTF8(name))) }
func tsInternalReplHistoryCall(loop *tsLoop, receiver tsValue, name string,args ...tsValue) tsValue {
	method:=tsGet(loop,receiver,tsStringReference(tsStringUTF8(name)));if method.kind!=tsFunctionKind{return tsU};return tsCallReceiver(loop,method,receiver,args...)
}
func tsInternalReplHistoryOutput(loop *tsLoop,state *tsInternalReplHistoryState,message string) {
	write:=tsGet(loop,state.context,tsStringReference(tsStringUTF8("_writeToOutput")));if write.kind==tsFunctionKind{tsCallReceiver(loop,write,state.context,tsStringReference(tsStringUTF8(message)));refresh:=tsGet(loop,state.context,tsStringReference(tsStringUTF8("_refreshLine")));if refresh.kind==tsFunctionKind{tsCallReceiver(loop,refresh,state.context)}}
}
func tsInternalReplHistoryReady(loop *tsLoop,callback,context tsValue) tsValue {
	if callback.kind==tsFunctionKind{return tsCall(loop,callback,tsNull,context)};return tsU
}
func tsInternalReplHistorySetPrevious(loop *tsLoop, receiver tsValue,state *tsInternalReplHistoryState) {
	state.originalHistoryPrev=tsGet(loop,state.context,tsStringReference(tsStringUTF8("_historyPrev")))
	previous:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{
		if len(tsInternalReplHistoryItems(state).values)==0{tsInternalReplHistoryOutput(loop,state,"\nPersistent history support disabled. Set the NODE_REPL_HISTORY environment\nvariable to a valid, user-writable path to enable.\n")}
		tsSet(loop,state.context,tsStringReference(tsStringUTF8("_historyPrev")),state.originalHistoryPrev)
		if state.originalHistoryPrev.kind==tsFunctionKind{return tsCallReceiver(loop,state.originalHistoryPrev,state.context)};return tsU
	}))
	tsSet(loop,state.context,tsStringReference(tsStringUTF8("_historyPrev")),previous)
}
func tsInternalReplHistoryPath(state *tsInternalReplHistoryState) string {
	if state.filePath!=""{return state.filePath};home,err:=os.UserHomeDir();if err!=nil{return ""};state.filePath=filepath.Join(home,".node_repl_history");return state.filePath
}
func tsInternalReplHistoryFlush(loop *tsLoop,state *tsInternalReplHistoryState) {
	state.timer=nil
	if state.file==nil{return}
	if state.writing{state.pending=true;return}
	state.writing=true
	items:=tsInternalReplHistoryItems(state);lines:=make([]string,len(items.values));for index,item:=range items.values{lines[index]=tsText(item)}
	_,err:=state.file.WriteAt([]byte(strings.Join(lines,"\n")),0);state.writing=false
	if err!=nil{return}
	if state.pending{state.pending=false;tsInternalReplHistoryOnLine(loop,state);return}
	state.isFlushing=state.timer!=nil
	if !state.isFlushing{tsInternalReplHistoryCall(loop,state.context,"emit",tsStringReference(tsStringUTF8("flushHistory")))}
}
func tsInternalReplHistoryOnLine(loop *tsLoop,state *tsInternalReplHistoryState) {
	state.isFlushing=true
	if state.timer!=nil{loop.clearTimeout(tsTaskValue(state.timer))}
	state.timer=loop.setTimeout(tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{tsInternalReplHistoryFlush(loop,state);return tsU})),tsNumberValue(15))
}
func tsInternalReplHistoryInitialize(loop *tsLoop,receiver tsValue,state *tsInternalReplHistoryState,callback tsValue) tsValue {
	if state.filePath=="" {
		// An explicit empty filePath disables disk persistence.
		if tsGet(loop,tsGet(loop,receiver,tsStringReference(tsStringUTF8("options"))),tsStringReference(tsStringUTF8("filePath"))).kind==tsStringKind{tsInternalReplHistorySetPrevious(loop,receiver,state);return tsInternalReplHistoryReady(loop,callback,state.context)}
	}
	path:=tsInternalReplHistoryPath(state);if path==""{tsInternalReplHistoryOutput(loop,state,"\nError: Could not get the home directory.\nREPL session history will not be persisted.\n");tsInternalReplHistorySetPrevious(loop,receiver,state);return tsInternalReplHistoryReady(loop,callback,state.context)}
	if tsTruthy(tsGet(loop,tsInternalReplNodeShimsModule(loop),tsStringReference(tsStringUTF8("isEnabled")))) { /* filesystem permission is represented by node-shims.has */ }
	_ = os.MkdirAll(filepath.Dir(path),0700)
	file,err:=os.OpenFile(path,os.O_CREATE|os.O_RDWR,0600)
	if err!=nil{tsInternalReplHistoryOutput(loop,state,"\nError: Could not open history file.\nREPL session history will not be persisted.\n");tsInternalReplHistorySetPrevious(loop,receiver,state);return tsInternalReplHistoryReady(loop,callback,state.context)}
	_ = file.Chmod(0600)
	data,_:=os.ReadFile(path);_ = file.Truncate(0);if len(data)>0{parts:=strings.FieldsFunc(string(data),func(r rune)bool{return r=='\n'||r=='\r'});items:=&tsArray{};for index:=0;index<len(parts)&&index<state.size;index++{items.values=append(items.values,tsStringReference(tsStringUTF8(parts[index])))};state.history=tsArrayValue(items)}
	state.file=file
	tsInternalReplHistoryCall(loop,state.context,"on",tsStringReference(tsStringUTF8("line")),tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{tsInternalReplHistoryOnLine(loop,state);return tsU})))
	tsInternalReplHistoryCall(loop,state.context,"once",tsStringReference(tsStringUTF8("exit")),tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{if state.isFlushing{tsInternalReplHistoryCall(loop,state.context,"once",tsStringReference(tsStringUTF8("flushHistory")),tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{tsInternalReplHistoryClose(loop,state);return tsU})));return tsU};tsInternalReplHistoryClose(loop,state);return tsU})))
	tsInternalReplHistoryFlush(loop,state);return tsInternalReplHistoryReady(loop,callback,state.context)
}
func tsInternalReplHistoryClose(loop *tsLoop,state *tsInternalReplHistoryState) { if state.timer!=nil{loop.clearTimeout(tsTaskValue(state.timer));state.timer=nil};if state.file!=nil{_ = state.file.Close();state.file=nil} }
func tsInternalReplHistoryConstructor(loop *tsLoop,context,options tsValue) tsValue {
	if tsNullish(options){options=tsObjectValue(tsNewObject())};if options.kind!=tsObjectKind&&options.kind!=tsInstanceKind{tsFSInvalid("options","an object")}
	if history:=tsGet(loop,options,tsStringReference(tsStringUTF8("history")));!tsIsUndefined(history)&&history.kind!=tsArrayKind{tsFSInvalid("history","an array")}
	size:=30;if value:=tsGet(loop,options,tsStringReference(tsStringUTF8("size")));!tsIsUndefined(value){if !tsIsNumeric(value)||tsNumber(value)<0{tsFSInvalid("size","a non-negative number")};size=int(tsNumber(value))}else if value:=tsGet(loop,context,tsStringReference(tsStringUTF8("historySize")));tsIsNumeric(value){size=int(tsNumber(value))}
	history:=tsGet(loop,options,tsStringReference(tsStringUTF8("history")));if tsIsUndefined(history){history=tsArrayValue(&tsArray{})}
	filePath:=tsGet(loop,options,tsStringReference(tsStringUTF8("filePath")));path:="";if filePath.kind==tsStringKind{path=strings.TrimSpace(tsText(filePath))}
	state:=&tsInternalReplHistoryState{context:context,history:history,size:size,index:-1,filePath:path,removeDuplicates:tsTruthy(tsGet(loop,options,tsStringReference(tsStringUTF8("removeHistoryDuplicates"))))}
	object:=tsNewObject();object.prototype=tsObjectValue(tsInternalReplHistory.prototype);value:=tsObjectValue(object);tsInternalReplHistory.states[uintptr(value.ref)]=state;object.set("context",context);object.set("options",options);return value
}
func tsInternalReplHistoryInstall(loop *tsLoop,prototype *tsObject) {
	prototype.descriptors=map[string]*tsDescriptor{}
	for _,property:=range []string{"size","isFlushing","history","index"}{name:=property;prototype.descriptors[name]=&tsDescriptor{accessor:true,enumerable:true,configurable:true,getter:tsNamedNativeMethod("get "+name,0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsInternalReplHistoryStateFor(receiver);switch name{case "size":return tsNumberValue(float64(state.size));case "isFlushing":return tsBooleanValue(state.isFlushing);case "history":return state.history;case "index":return tsNumberValue(float64(state.index))};return tsU}),setter:tsNamedNativeMethod("set "+name,1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsInternalReplHistoryStateFor(receiver);switch name{case "history":if tsArg(args,0).kind!=tsArrayKind{tsFSInvalid("history","an array")};state.history=tsArg(args,0);case "index":state.index=int(tsArg(args,0).ToInt32())};return tsArg(args,0)})}
	}
	prototype.set("initialize",tsNamedNativeMethod("initialize",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsInternalReplHistoryInitialize(loop,receiver,tsInternalReplHistoryStateFor(receiver),tsArg(args,0))}))
	prototype.set("addHistory",tsNamedNativeMethod("addHistory",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsInternalReplHistoryStateFor(receiver);line:=tsInternalReplHistoryGet(loop,state,"line");text:=tsText(line);if text==""||state.size==0||strings.TrimSpace(text)==""{return line};history:=tsInternalReplHistoryItems(state);if tsTruthy(tsArg(args,0))&&state.index==-1&&len(history.values)>0{history.values=history.values[1:]}else if tsTruthy(tsArg(args,1))&&len(history.values)>0{history.values=history.values[1:]};normalized:=strings.ReplaceAll(text,"\n","\r");if len(history.values)==0||tsText(history.values[0])!=normalized{if state.removeDuplicates{for index,value:=range history.values{if tsText(value)==normalized{history.values=append(history.values[:index],history.values[index+1:]...);break}}};history.values=append([]tsValue{tsStringReference(tsStringUTF8(normalized))},history.values...);if len(history.values)>state.size{history.values=history.values[:state.size]}};state.index=-1;value:=history.values[0];if tsTruthy(tsArg(args,0)){value=tsInternalReadlineReverse(loop,value,tsStringReference(tsStringUTF8("\r")),tsStringReference(tsStringUTF8("\r")))};tsInternalReplHistoryCall(loop,state.context,"emit",tsStringReference(tsStringUTF8("history")),state.history);return value}))
	prototype.set("canNavigateToNext",tsNamedNativeMethod("canNavigateToNext",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsInternalReplHistoryStateFor(receiver);return tsBooleanValue(state.index>-1&&len(tsInternalReplHistoryItems(state).values)>0)}))
	prototype.set("navigateToNext",tsNamedNativeMethod("navigateToNext",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsInternalReplHistoryStateFor(receiver);items:=tsInternalReplHistoryItems(state);if state.index<=-1||len(items.values)==0{return tsNull};search:=tsText(tsArg(args,0));index:=state.index-1;for index>=0&&(!strings.HasPrefix(tsText(items.values[index]),search)||tsStrictEqual(tsInternalReplHistoryGet(loop,state,"line"),items.values[index])){index--};state.index=index;if index<0{return tsStringReference(tsStringUTF8(search))};return tsStringReference(tsStringUTF8(strings.ReplaceAll(tsText(items.values[index]),"\r","\n")))}))
	prototype.set("canNavigateToPrevious",tsNamedNativeMethod("canNavigateToPrevious",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsInternalReplHistoryStateFor(receiver);return tsBooleanValue(state.index!=len(tsInternalReplHistoryItems(state).values)&&len(tsInternalReplHistoryItems(state).values)>0)}))
	prototype.set("navigateToPrevious",tsNamedNativeMethod("navigateToPrevious",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsInternalReplHistoryStateFor(receiver);items:=tsInternalReplHistoryItems(state);if state.index==len(items.values)||len(items.values)==0{return tsNull};search:=tsText(tsArg(args,0));index:=state.index+1;for index<len(items.values)&&(!strings.HasPrefix(tsText(items.values[index]),search)||tsStrictEqual(tsInternalReplHistoryGet(loop,state,"line"),items.values[index])){index++};state.index=index;if index>=len(items.values){return tsStringReference(tsStringUTF8(search))};return tsStringReference(tsStringUTF8(strings.ReplaceAll(tsText(items.values[index]),"\r","\n")))}))
	prototype.set("closeHandle",tsNamedNativeMethod("closeHandle",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{tsInternalReplHistoryClose(tsInternalReplHistoryStateFor(receiver));return tsPromiseValue(loop.resolved(tsU,false))}))
}
func tsInternalReplHistoryModule(loop *tsLoop) tsValue {
	tsInternalReplHistory.Do(func(){tsInternalReplHistory.states=map[uintptr]*tsInternalReplHistoryState{};tsInternalReplHistory.prototype=tsNewObject();tsInternalReplHistory.prototype.prototype=tsObjectPrototype();tsInternalReplHistoryInstall(loop,tsInternalReplHistory.prototype);constructor:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsInternalReplHistoryConstructor(loop,tsArg(args,0),tsArg(args,1))});constructor.name="ReplHistory";constructor.length=2;constructor.constructible=true;constructor.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsInternalReplHistoryConstructor(loop,tsArg(args,0),tsArg(args,1))};constructor.properties=tsNewObject();constructor.properties.set("prototype",tsObjectValue(tsInternalReplHistory.prototype));tsInternalReplHistory.prototype.set("constructor",tsFunctionValue(constructor));module:=tsNewObject();module.set("ReplHistory",tsFunctionValue(constructor));module.set("default",tsObjectValue(module));tsInternalReplHistory.constructor=tsFunctionValue(constructor);tsInternalReplHistory.module=tsObjectValue(module)})
	return tsInternalReplHistory.module
}
`
