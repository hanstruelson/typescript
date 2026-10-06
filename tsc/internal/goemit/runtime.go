package goemit

// Runtime is emitted into the single output file; its dependencies use the compiler module.
// Each event loop owns its task registry and runs callbacks on one goroutine.
// User workers have their own loop and share application references explicitly.
const Runtime = `import (
 "unsafe"
 "sync"
 "fmt"
 "math"
 "math/big"
 "os"
 "strings"
 "sort"
 "time"
 "unicode/utf16"
 "unicode/utf8"
 "strconv"
 "github.com/dop251/goja"
 "golang.org/x/text/cases"
 "golang.org/x/text/language"
 "golang.org/x/text/unicode/norm"
 "golang.org/x/text/collate"
)

type tsCell struct { value tsValue; initialized, constant bool }
func tsBinding(initialized, constant bool) *tsCell { return &tsCell{tsU, initialized, constant} }
func (c *tsCell) get() tsValue {if !c.initialized {panic("Cannot access binding before initialization")};if c.value.kind==tsImportKind{ref:=(*tsImportRef)(c.value.ref);return ref.module.read(ref.name)};return c.value}
func (c *tsCell) init(value tsValue) tsValue { c.value = value; c.initialized = true; return value }
func (c *tsCell) set(value tsValue) tsValue { if !c.initialized { panic("Cannot access binding before initialization") }; if c.constant { panic("Assignment to constant variable") }; c.value = value; return value }
func tsClone(c *tsCell) *tsCell { value := *c; return &value }
type tsFunction struct {call func(*tsLoop,...tsValue) tsValue}
func tsFunc(call func(*tsLoop,...tsValue)tsValue)*tsFunction{return &tsFunction{call:call}}
func tsCall(loop *tsLoop,value tsValue,args ...tsValue)tsValue{if value.kind!=tsFunctionKind{panic("Value is not callable")};return (*tsFunction)(value.ref).call(loop,args...)}
func tsArg(args []tsValue, index int) tsValue { if index < len(args) { return args[index] }; return tsU }
func tsNumberText(number float64)string{
 if math.IsNaN(number){return "NaN"};if math.IsInf(number,1){return "Infinity"};if math.IsInf(number,-1){return "-Infinity"};if number==0{return "0"}
 magnitude:=math.Abs(number);if magnitude>=1e-6&&magnitude<1e21{return strconv.FormatFloat(number,'f',-1,64)}
 text:=strconv.FormatFloat(number,'e',-1,64);index:=strings.LastIndexByte(text,'e');exponent,_:=strconv.Atoi(text[index+1:]);sign:="";if exponent>=0{sign="+"};return text[:index+1]+sign+strconv.Itoa(exponent)
}
func tsText(value tsValue)string{switch value.kind{
 case tsUndefinedKind:return "undefined"
 case tsNullKind:return "null"
 case tsIntKind,tsInt8Kind,tsInt16Kind,tsInt32Kind,tsInt64Kind:return strconv.FormatInt(int64(math.Float64bits(value.number)),10)
 case tsUintKind,tsUint8Kind,tsUint16Kind,tsUint32Kind,tsUint64Kind:return strconv.FormatUint(math.Float64bits(value.number),10)
 case tsFloat32Kind:return strconv.FormatFloat(value.number,'g',-1,32)
 case tsNumberKind:return tsNumberText(value.number)
 case tsBooleanKind:return strconv.FormatBool(value.number!=0)
 case tsStringKind:return (*tsString)(value.ref).String()
 case tsErrorKind:return (*tsRuntimeError)(value.ref).String()
 case tsObjectKind:return fmt.Sprint((*tsObject)(value.ref))
 case tsArrayKind:return fmt.Sprint((*tsArray)(value.ref))
 case tsFunctionKind:return fmt.Sprint((*tsFunction)(value.ref))
 case tsClassKind:return fmt.Sprint((*tsClass)(value.ref))
 case tsInstanceKind:return "[object Object]"
 case tsRegExpKind:return fmt.Sprint((*tsRegExp)(value.ref))
 case tsECMAKind:return fmt.Sprint((*tsECMAObject)(value.ref))
 case tsPromiseKind:return fmt.Sprint((*tsPromise)(value.ref))
 case tsTaskKind:return fmt.Sprint((*tsTask)(value.ref))
 case tsNamespaceKind:return fmt.Sprint(*(*tsNamespace)(value.ref))
 case tsImportKind:return fmt.Sprint(*(*tsImportRef)(value.ref))
 case tsModuleKind:return fmt.Sprint((*tsModule)(value.ref))
 case tsIteratorKind:return fmt.Sprint((*tsIterator)(value.ref))
 default:panic("Invalid JavaScript value tag")}}
func(value tsValue)String()string{return tsText(value)}
func tsConsoleLog(values ...tsValue) tsValue { for i, value := range values { if i != 0 { fmt.Print(" ") }; fmt.Print(tsText(value)) }; fmt.Println(); return tsU }
func tsIsUndefined(value tsValue) bool {return value.kind==tsUndefinedKind}
func tsNullish(value tsValue) bool {return value.kind==tsNullKind || value.kind==tsUndefinedKind}
type tsArray struct { values []tsValue;properties map[string]tsValue;holes map[int]bool }
type tsIterator struct {array *tsArray; text *tsString; ecma *goja.Object; index int; value tsValue}
func tsIterate(value tsValue)*tsIterator{switch value.kind{case tsArrayKind:return &tsIterator{array:(*tsArray)(value.ref)};case tsStringKind:return &tsIterator{text:(*tsString)(value.ref)};case tsECMAKind:v:=(*tsECMAObject)(value.ref);if v.object.ClassName()=="Array"{values:=[]tsValue{};for i:=int64(0);i<v.object.Get("length").ToInteger();i++{values=append(values,tsFromECMA(v.object.Get(strconv.FormatInt(i,10))))};return &tsIterator{array:&tsArray{values:values}}};return &tsIterator{ecma:v.object};default:panic("Value is not iterable")}}
func (it *tsIterator) next() bool {if it.ecma!=nil {call,_:=goja.AssertFunction(it.ecma.Get("next"));result,err:=call(it.ecma);if err!=nil {tsECMAFailure(err)};object:=result.ToObject(tsEngine());if object.Get("done").ToBoolean(){return false};it.value=tsFromECMA(object.Get("value"));return true};if it.array!=nil {if it.index>=len(it.array.values){return false};it.value=it.array.values[it.index];it.index++;return true};if it.index>=len(it.text.units){return false};start:=it.index;it.index++;c:=it.text.units[start];if c>=0xd800&&c<=0xdbff&&it.index<len(it.text.units)&&it.text.units[it.index]>=0xdc00&&it.text.units[it.index]<=0xdfff {it.index++};it.value=it.text.slice(start,it.index);return true}
func tsGet(value,key tsValue) tsValue { switch value.kind {
 case tsClassKind:v:=(*tsClass)(value.ref);return tsInstanceProperties(v.static).get(tsPropertyKey(key))
 case tsInstanceKind:v:=(*tsProperties)(value.ref);return v.get(tsPropertyKey(key))
 case tsObjectKind:v:=(*tsObject)(value.ref);return v.get(tsPropertyKey(key))
 case tsArrayKind:v:=(*tsArray)(value.ref);
  if value,ok:=v.properties[tsPropertyKey(key)];ok {return value}
  if tsText(key)=="length" {return float64(len(v.values))}
  if tsText(key)=="push" {return tsFunc(func(args ...tsValue) tsValue {v.values=append(v.values,args...);return float64(len(v.values))})}
  index,err:=strconv.Atoi(tsPropertyKey(key));if err!=nil||strconv.Itoa(index)!=tsPropertyKey(key)||index<0 || index>=len(v.values) {return tsU};return v.values[index]
 case tsStringKind:v:=(*tsString)(value.ref);return tsStringGet(v,key)
 case tsRegExpKind:v:=(*tsRegExp)(value.ref);return tsECMAGet(v.object,key)
 case tsECMAKind:v:=(*tsECMAObject)(value.ref);return tsECMAGet(v.object,key)
 case tsErrorKind:v:=(*tsRuntimeError)(value.ref);if tsText(key)=="message" {return tsStringUTF8(v.message)};if tsText(key)=="name" {return tsStringUTF8(v.name)};return tsU
 case tsNamespaceKind:v:=(*tsNamespace)(value.ref);return v.module.read(tsPropertyKey(key))
 case tsPromiseKind:v:=(*tsPromise)(value.ref);
  if tsText(key)=="then" || tsText(key)=="catch" {return tsFunc(func(args ...tsValue) tsValue {
   next:=loop.promise();success,failure:=tsArg(args,0),tsArg(args,1);if tsText(key)=="catch" {failure=success;success=tsU}
   loop.await(tsPromiseValue(v),func(result tsResult){callback:=success;if result.rejected {callback=failure};if callback.kind!=tsFunctionKind {next.settle(result);return};func(){defer func(){if err:=recover();err!=nil {next.settle(tsResult{tsUnwrap(err),true})}}();next.settle(tsResult{tsCall(callback,result.value),false})}()});return next
  })}
 }
 panic("Unsupported property access")
}
func tsSet(value,key,item tsValue)tsValue{switch value.kind{case tsClassKind:return tsInstanceProperties((*tsClass)(value.ref).static).set(tsPropertyKey(key),item);case tsInstanceKind:return (*tsProperties)(value.ref).set(tsPropertyKey(key),item);case tsObjectKind:return (*tsObject)(value.ref).set(tsPropertyKey(key),item);case tsRegExpKind:return tsECMASet((*tsRegExp)(value.ref).object,key,item);case tsECMAKind:return tsECMASet((*tsECMAObject)(value.ref).object,key,item)};if value.kind!=tsArrayKind{panic("Expected an array")};array:=(*tsArray)(value.ref);name:=tsPropertyKey(key);if name=="length"{length:=tsNumber(item);if length<0||length>4294967295||math.Trunc(length)!=length{panic("Invalid array length")};for len(array.values)<int(length){tsArrayHole(array)};clear(array.values[int(length):]);array.values=array.values[:int(length)];for index:=range array.holes{if index>=int(length){delete(array.holes,index)}};return item};index,err:=strconv.Atoi(name);if err!=nil||index<0||uint64(index)>=4294967295||strconv.Itoa(index)!=name{if array.properties==nil{array.properties=map[string]tsValue{}};array.properties[name]=item;return item};for len(array.values)<=index{tsArrayHole(array)};delete(array.holes,index);array.values[index]=item;return item}
func tsNumber(value tsValue) float64 {switch value.kind {case tsIntKind,tsInt8Kind,tsInt16Kind,tsInt32Kind,tsInt64Kind:return float64(int64(math.Float64bits(value.number)));case tsUintKind,tsUint8Kind,tsUint16Kind,tsUint32Kind,tsUint64Kind:return float64(math.Float64bits(value.number));case tsFloat32Kind,tsNumberKind,tsBooleanKind:return value.number;case tsNullKind:return 0;case tsUndefinedKind:return math.NaN();case tsStringKind:return goja.StringFromUTF16((*tsString)(value.ref).units).ToFloat();default:panic("Expected a number")}}
func tsTruthy(value tsValue) bool {switch value.kind {case tsUndefinedKind,tsNullKind:return false;case tsIntKind,tsInt8Kind,tsInt16Kind,tsInt32Kind,tsInt64Kind,tsUintKind,tsUint8Kind,tsUint16Kind,tsUint32Kind,tsUint64Kind:return math.Float64bits(value.number)!=0;case tsFloat32Kind,tsNumberKind:return value.number!=0 && !math.IsNaN(value.number);case tsBooleanKind:return value.number!=0;case tsStringKind:return len((*tsString)(value.ref).units)!=0;default:return true}}
func tsBinary(op string, left, right tsValue) tsValue {
 if left.kind==right.kind {switch left.kind {
 case tsIntKind:if op=="+"||op=="-"||op=="*"||op=="/"||op=="%"||op=="**"||op=="&"||op=="|"||op=="^"||op=="<<"||op==">>"||op==">>>"{return tsIntegerBinary[int](op,left,right)}
 case tsInt8Kind:if op=="+"||op=="-"||op=="*"||op=="/"||op=="%"||op=="**"||op=="&"||op=="|"||op=="^"||op=="<<"||op==">>"||op==">>>"{return tsIntegerBinary[int8](op,left,right)}
 case tsInt16Kind:if op=="+"||op=="-"||op=="*"||op=="/"||op=="%"||op=="**"||op=="&"||op=="|"||op=="^"||op=="<<"||op==">>"||op==">>>"{return tsIntegerBinary[int16](op,left,right)}
 case tsInt32Kind:if op=="+"||op=="-"||op=="*"||op=="/"||op=="%"||op=="**"||op=="&"||op=="|"||op=="^"||op=="<<"||op==">>"||op==">>>"{return tsIntegerBinary[int32](op,left,right)}
 case tsInt64Kind:if op=="+"||op=="-"||op=="*"||op=="/"||op=="%"||op=="**"||op=="&"||op=="|"||op=="^"||op=="<<"||op==">>"||op==">>>"{return tsIntegerBinary[int64](op,left,right)}
 case tsUintKind:if op=="+"||op=="-"||op=="*"||op=="/"||op=="%"||op=="**"||op=="&"||op=="|"||op=="^"||op=="<<"||op==">>"||op==">>>"{return tsIntegerBinary[uint](op,left,right)}
 case tsUint8Kind:if op=="+"||op=="-"||op=="*"||op=="/"||op=="%"||op=="**"||op=="&"||op=="|"||op=="^"||op=="<<"||op==">>"||op==">>>"{return tsIntegerBinary[uint8](op,left,right)}
 case tsUint16Kind:if op=="+"||op=="-"||op=="*"||op=="/"||op=="%"||op=="**"||op=="&"||op=="|"||op=="^"||op=="<<"||op==">>"||op==">>>"{return tsIntegerBinary[uint16](op,left,right)}
 case tsUint32Kind:if op=="+"||op=="-"||op=="*"||op=="/"||op=="%"||op=="**"||op=="&"||op=="|"||op=="^"||op=="<<"||op==">>"||op==">>>"{return tsIntegerBinary[uint32](op,left,right)}
 case tsUint64Kind:if op=="+"||op=="-"||op=="*"||op=="/"||op=="%"||op=="**"||op=="&"||op=="|"||op=="^"||op=="<<"||op==">>"||op==">>>"{return tsIntegerBinary[uint64](op,left,right)}
} }
 if tsIsNumeric(left)&&tsIsNumeric(right) {switch op{case "<","<=",">",">=":comparison,valid:=tsNumericCompare(left,right);if !valid{return false};switch op{case "<":return comparison<0;case "<=":return comparison<=0;case ">":return comparison>0;default:return comparison>=0}}}

 switch op {
 case "in":return tsHas(right,left)
 case "==","!=":equal:=tsLooseEqual(left,right);if op=="!="{return !equal};return equal
 case "**":return tsPow(tsNumber(tsToPrimitive(left)),tsNumber(tsToPrimitive(right)))
 case "&","|","^","<<",">>",">>>":return tsUint32Bitwise(op,tsToUint32(left),tsToUint32(right))
 case "+": if left.kind==tsStringKind{return tsStringConcat(tsStringValue(left),tsStringValue(right))};if right.kind==tsStringKind{return tsStringConcat(tsStringValue(left),tsStringValue(right))}; return tsNumber(left)+tsNumber(right)
 case "-": return tsNumber(left)-tsNumber(right)
 case "*": return tsNumber(left)*tsNumber(right)
 case "/": return tsNumber(left)/tsNumber(right)
 case "%": return math.Mod(tsNumber(left),tsNumber(right))
 case "<":if left.kind==tsStringKind && right.kind==tsStringKind{return tsStringCompare((*tsString)(left.ref),(*tsString)(right.ref))<0}; return tsNumber(left)<tsNumber(right)
 case "<=":if left.kind==tsStringKind && right.kind==tsStringKind{return tsStringCompare((*tsString)(left.ref),(*tsString)(right.ref))<=0}; return tsNumber(left)<=tsNumber(right)
 case ">":if left.kind==tsStringKind && right.kind==tsStringKind{return tsStringCompare((*tsString)(left.ref),(*tsString)(right.ref))>0}; return tsNumber(left)>tsNumber(right)
 case ">=":if left.kind==tsStringKind && right.kind==tsStringKind{return tsStringCompare((*tsString)(left.ref),(*tsString)(right.ref))>=0}; return tsNumber(left)>=tsNumber(right)
 case "===", "!==": equal:=tsStrictEqual(left,right); if op=="!==" { return !equal }; return equal
 default: panic("Unsupported operator")
 }
}

type tsThrown struct {value tsValue}
// recover() returns Go's any; translate it once at the panic boundary.
func tsUnwrap(value any)tsValue{switch v:=value.(type){case tsThrown:return v.value;case tsValue:return v;case error:return tsStringReference(tsStringUTF8(v.Error()));case string:return tsStringReference(tsStringUTF8(v));case nil:return tsNull;default:return tsStringReference(tsStringUTF8(fmt.Sprint(v)))}}
type tsResult struct { value tsValue; rejected bool }
type tsPromise struct {mu sync.Mutex; loop *tsLoop; settled, handled bool; result tsResult; listeners []func(tsResult) }
// A task is registered and finished by the loop. The worker exclusively writes
// result until publishing this pointer through the completion channel.
type tsTask struct { work func() tsResult; finish func(tsResult); result tsResult; canceled bool; cancel func(); timer, completed bool; deadline time.Time }
type tsLoop struct { queueMu sync.Mutex;incoming []func();wake chan struct{}; ecma *goja.Runtime; pending map[*tsTask]struct{}; completions chan *tsTask; ready []func(); promises []*tsPromise; errors []tsValue; timers []*tsTask; modules []*tsModule }
func tsNewLoop() *tsLoop { return &tsLoop{pending:make(map[*tsTask]struct{}), completions:make(chan *tsTask),wake:make(chan struct{},1)} }
func (l *tsLoop) promise() *tsPromise { p:= &tsPromise{loop:l}; l.promises=append(l.promises,p); return p }
func(l *tsLoop)post(callback func()){l.queueMu.Lock();l.incoming=append(l.incoming,callback);l.queueMu.Unlock();select{case l.wake<-struct{}{}:default:}}
func(l *tsLoop)drain(){l.queueMu.Lock();l.ready=append(l.ready,l.incoming...);l.incoming=nil;l.queueMu.Unlock()}
func (p *tsPromise) settle(result tsResult) {
 if result.value.kind==tsPromiseKind&&!result.rejected {
  other:=(*tsPromise)(result.value.ref)
  if other==p{p.settle(tsResult{"Promise cannot resolve to itself",true});return}
  other.then(p.settle);return
 }
 p.mu.Lock();if p.settled {p.mu.Unlock();return};p.settled=true;p.result=result;listeners:=p.listeners;p.listeners=nil;p.mu.Unlock()
 for _,callback:=range listeners{cb:=callback;p.loop.post(func(){cb(result)})}
}
func(p *tsPromise)then(callback func(tsResult)){p.mu.Lock();p.handled=true;if !p.settled{p.listeners=append(p.listeners,callback);p.mu.Unlock();return};result:=p.result;p.mu.Unlock();p.loop.post(func(){callback(result)})}
func (l *tsLoop) resolved(value tsValue, rejected bool) *tsPromise { p:=l.promise(); p.settle(tsResult{value,rejected}); return p }
func (l *tsLoop) constructPromise(executor tsValue) *tsPromise {
 p:=l.promise();finished:=false
 resolve:=tsFunc(func(args ...tsValue) tsValue {if !finished {finished=true;p.settle(tsResult{tsArg(args,0),false})};return tsU})
 reject:=tsFunc(func(args ...tsValue) tsValue {if !finished {finished=true;p.settle(tsResult{tsArg(args,0),true})};return tsU})
 func(){defer func(){if err:=recover();err!=nil && !finished {finished=true;p.settle(tsResult{tsUnwrap(err),true})}}();tsCall(executor,resolve,reject)}()
 return p
}
func (l *tsLoop) all(values tsValue) *tsPromise {
 p:=l.promise();if values.kind!=tsArrayKind{p.settle(tsResult{"Promise.all expects an array",true});return p};array:=(*tsArray)(values.ref)
 output:= &tsArray{values:make([]tsValue,len(array.values))};remaining:=len(array.values)
 if remaining==0 {p.settle(tsResult{output,false});return p}
 for index,value:=range array.values {i:=index;l.await(value,func(result tsResult){if result.rejected {p.settle(result);return};output.values[i]=result.value;remaining--;if remaining==0 {p.settle(tsResult{output,false})}})}
 return p
}
func(l *tsLoop)await(value tsValue,callback func(tsResult)){
 if value.kind!=tsPromiseKind {l.ready=append(l.ready,func(){callback(tsResult{value,false})});return}
 p:=(*tsPromise)(value.ref);if p.loop==l{p.then(callback);return}
 done:=make(chan tsResult,1);p.then(func(result tsResult){done<-result});l.submit(func()tsResult{return <-done},callback)
}
// submit is called only on the owning loop. Native adapters snapshot inputs;
// user workers run application code with their own loop and shared references.
func (l *tsLoop) submit(work func() tsResult, finish func(tsResult)) *tsTask {
 task:= &tsTask{work:work,finish:finish}
 l.pending[task]=struct{}{}
 go func(task *tsTask,done chan<- *tsTask) {
  result:=tsResult{"Worker terminated without a result",true}
  defer func(){if failure:=recover();failure!=nil {result=tsResult{tsUnwrap(failure),true}};task.result=result;done<-task}()
  result=task.work()
 }(task,l.completions)
 return task
}
// spawn reuses task registration and publishes only completion to the caller.
// The tagged argument slice copies scalar slots; references remain shared.
func(l *tsLoop)spawn(function tsValue,args ...tsValue)*tsPromise{
 if function.kind!=tsFunctionKind{panic("go requires a function call")}
 return l.start(func()tsResult{
  child:=tsNewLoop();loop:=child;_ = loop;result:=tsResult{tsU,false};completed:=false
  func(){defer func(){if failure:=recover();failure!=nil{result=tsResult{tsUnwrap(failure),true};completed=true}}();value:=tsCall(function,args...);child.await(value,func(r tsResult){result=r;completed=true})}()
  if err:=child.run();err!=nil{return tsResult{tsErrorValue(&tsRuntimeError{"Error",err.Error()}),true}}
  if !completed{return tsResult{"Worker did not complete",true}};return result
 })
}
func (l *tsLoop) start(work func() tsResult) *tsPromise {p:=l.promise();l.submit(work,p.settle);return p}
// Browser-style timers return a handle, not a Promise. The worker only sleeps;
// callback invocation and captured-variable access happen on the event loop.
func tsTimerDuration(value tsValue) time.Duration {
 milliseconds:=float64(0)
 if !tsIsUndefined(value) {milliseconds=tsNumber(value)}
 if math.IsNaN(milliseconds) || math.IsInf(milliseconds,0) || milliseconds<0 || milliseconds>2147483647 {milliseconds=0}
 return time.Duration(math.Trunc(milliseconds))*time.Millisecond
}
func (l *tsLoop) setTimeout(callback, milliseconds tsValue, args ...tsValue) *tsTask {
 if callback.kind!=tsFunctionKind {panic("setTimeout expects a function")}
 duration:=tsTimerDuration(milliseconds)
 deadline:=time.Now().Add(duration)
 canceled:=make(chan struct{})
 // Only the finish function captures the callback and its arguments.
 var task *tsTask
 task=l.submit(func()tsResult {timer:=time.NewTimer(time.Until(deadline));defer timer.Stop();select {case <-timer.C:case <-canceled:};return tsResult{tsU,false}},func(result tsResult){
  if task.canceled {return}
  if result.rejected {l.errors=append(l.errors,result.value);return}
  tsCall(callback,args...)
 })
 task.cancel=func(){close(canceled)}
 task.timer=true;task.deadline=deadline
 l.timers=append(l.timers,task)
 sort.SliceStable(l.timers,func(i,j int)bool{return l.timers[i].deadline.Before(l.timers[j].deadline)})
 return task
}
func(l *tsLoop)clearTimeout(handle tsValue)tsValue{if handle.kind==tsTaskKind{task:=(*tsTask)(handle.ref);if !task.canceled{task.canceled=true;if task.cancel!=nil{task.cancel()};if task.completed{l.removeTimer(task);delete(l.pending,task)}}};return tsU}
func(l *tsLoop)readFile(path tsValue)*tsPromise{if path.kind!=tsStringKind{return l.resolved("readFile expects a string",true)};text:=(*tsString)(path.ref);return l.start(func()tsResult{bytes,err:=os.ReadFile(text.String());if err!=nil{return tsResult{err.Error(),true}};return tsResult{tsStringUTF8(string(bytes)),false}})}
func (l *tsLoop) delay(milliseconds tsValue) *tsPromise { duration:=tsTimerDuration(milliseconds); return l.start(func() tsResult {time.Sleep(duration); return tsResult{tsU,false}}) }
func (l *tsLoop) invoke(callback func()) {
 defer func(){if failure:=recover();failure!=nil {l.errors=append(l.errors,tsUnwrap(failure))}}()
 callback()
}
func (l *tsLoop) removeTimer(task *tsTask) {for i,timer:=range l.timers {if timer==task {copy(l.timers[i:],l.timers[i+1:]);l.timers[len(l.timers)-1]=nil;l.timers=l.timers[:len(l.timers)-1];return}}}
func (l *tsLoop) run() error {
 for {
  l.drain()
  for len(l.ready)>0 {callback:=l.ready[0];l.ready[0]=nil;l.ready=l.ready[1:];l.invoke(callback);l.drain()}
  // Process one timer at a time, with a reaction checkpoint between callbacks.
  // Deadline ordering prevents goroutine scheduling from reordering equal timers.
  if len(l.timers)>0 && l.timers[0].completed {
   task:=l.timers[0];l.removeTimer(task);delete(l.pending,task);l.invoke(func(){task.finish(task.result)});continue
  }
  l.drain();if len(l.ready)>0{continue}
  if len(l.pending)==0 {break}
  var task *tsTask;select{case task=<-l.completions:case <-l.wake:continue}
  if _,ok:=l.pending[task];!ok {l.errors=append(l.errors,"Unknown completed task");continue}
  task.completed=true
  if task.timer {
   if task.canceled {l.removeTimer(task);delete(l.pending,task)}
   continue
  }
  delete(l.pending,task)
  l.invoke(func(){task.finish(task.result)})
 }
 for _,module:=range l.modules {if len(module.requests)>0 && !module.failed {missing:=[]string{};for _,request:=range module.requests {if request.full {missing=append(missing,"<module completion>")};for name:=range request.missing {missing=append(missing,name)}};sort.Strings(missing);l.errors=append(l.errors,"Unresolved module dependency: "+module.name+" needs "+strings.Join(missing,", "))}}
 for _,p:=range l.promises {p.mu.Lock();if p.settled && p.result.rejected && !p.handled {l.errors=append(l.errors,"Unhandled rejection: "+tsText(p.result.value))};p.mu.Unlock()}
 if len(l.errors)>0 {messages:=make([]string,len(l.errors));for i,value:=range l.errors {messages[i]=tsText(value)};return fmt.Errorf("Async runtime errors: %s",strings.Join(messages,"; "))}
 return nil
}

// Protected regions and pending abrupt completions mirror ES5's trys/ops stacks.
type tsHandler struct { catch, finally, end int; phase int }
type tsAbrupt struct { kind string; value tsValue; target, depth int }
type tsMachine struct { loop *tsLoop; pc int; step func(); output *tsPromise; async, done, suspended bool; result tsValue; handlers []tsHandler; abrupt []tsAbrupt; thrown tsValue; module *tsModule; blocked bool }
func (m *tsMachine) resume() {
 for !m.done && !m.suspended {
  func(){ defer func(){if err:=recover(); err!=nil {m.raise(tsUnwrap(err))}}(); m.step() }()
 }
}
func (m *tsMachine) await(value tsValue, next int) {
 m.suspended=true;m.blocked=true
 m.loop.await(value,func(result tsResult){m.blocked=false;m.suspended=false; if result.rejected {m.raise(result.value)} else {m.result=result.value; m.pc=next}; m.resume()})
}
func (m *tsMachine) complete(value tsValue) { m.done=true; m.result=value; if m.module!=nil {m.module.finish(tsResult{value,false})}; if m.async {m.output.settle(tsResult{value,false})} }
func (m *tsMachine) raise(value tsValue) {m.transfer(tsAbrupt{kind:"throw",value:value})}
func (m *tsMachine) transfer(op tsAbrupt) {
 for len(m.handlers)>op.depth {
  index:=len(m.handlers)-1; handler:= &m.handlers[index]
  if op.kind=="throw" && handler.phase==0 && handler.catch>=0 { handler.phase=1; m.thrown=op.value; m.pc=handler.catch; return }
  if handler.phase<2 && handler.finally>=0 { handler.phase=2; m.abrupt=append(m.abrupt,op); m.pc=handler.finally; return }
  if handler.phase==2 { m.abrupt=m.abrupt[:len(m.abrupt)-1] }
  m.handlers=m.handlers[:index]
 }
 switch op.kind {
 case "jump": m.pc=op.target
 case "return": m.complete(op.value)
 case "throw": m.done=true; if m.module!=nil {m.module.finish(tsResult{op.value,true})} else if m.async {m.output.settle(tsResult{op.value,true})} else {panic(tsThrown{op.value})}
 }
}
func (m *tsMachine) endFinally() { index:=len(m.handlers)-1; m.handlers=m.handlers[:index]; last:=len(m.abrupt)-1; op:=m.abrupt[last]; m.abrupt=m.abrupt[:last]; m.transfer(op) }
func tsPow(base,exponent float64)float64{if math.IsInf(exponent,0)&&math.Abs(base)==1{return math.NaN()};return math.Pow(base,exponent)}
func tsNumberUint32(number float64)uint32{if number==0||math.IsNaN(number)||math.IsInf(number,0){return 0};number=math.Mod(math.Trunc(number),4294967296);if number<0{number+=4294967296};return uint32(number)}
func tsToUint32(value tsValue)uint32{return tsNumberUint32(tsNumber(tsToPrimitive(value)))}
func tsNumberBitwise(op string,left,right float64)float64{return tsUint32Bitwise(op,tsNumberUint32(left),tsNumberUint32(right))}
func tsUint32Bitwise(op string,a,b uint32)float64{switch op{case "&":return float64(int32(a&b));case "|":return float64(int32(a|b));case "^":return float64(int32(a^b));case "<<":return float64(int32(a<<(b&31)));case ">>":return float64(int32(a)>>(b&31));default:return float64(a>>(b&31))}}
func tsNumberBitwiseNot(value float64)float64{return float64(^int32(tsNumberUint32(value)))}
func tsBitwiseNot(value tsValue)tsValue{
 if tsIsSigned(value)||tsIsUnsigned(value){bits:=^math.Float64bits(value.number);return tsIntegerBits(value.kind,bits)}
 return tsNumberValue(float64(^int32(tsToUint32(value))))
}
func tsIsPrimitiveValue(value tsValue)bool{return tsIsNumeric(value)||value.kind==tsStringKind||value.kind==tsBooleanKind||tsNullish(value)}
func tsToPrimitive(value tsValue)tsValue{
 if tsIsPrimitiveValue(value){return value}
 for _,name:=range []string{"valueOf","toString"}{method:=tsGet(value,tsStringReference(tsStringUTF8(name)));if method.kind==tsFunctionKind{result:=tsCall(method);if tsIsPrimitiveValue(result){return result}}}
 switch value.kind{case tsObjectKind,tsInstanceKind:return tsStringReference(tsStringUTF8("[object Object]"));case tsArrayKind:return tsStringReference(tsStringValue(value));case tsECMAKind:return tsFromECMA((*tsECMAObject)(value.ref).object.ToString());case tsRegExpKind:return tsFromECMA((*tsRegExp)(value.ref).object.ToString())}
 panic(tsThrown{tsErrorValue(&tsRuntimeError{"TypeError","Cannot convert object to primitive value"})})
}
func tsLooseEqual(a,b tsValue)bool{
 if tsStrictEqual(a,b){return true};if tsNullish(a)&&tsNullish(b){return true}
 if a.kind==tsBooleanKind{return tsLooseEqual(tsNumberValue(a.number),b)};if b.kind==tsBooleanKind{return tsLooseEqual(a,tsNumberValue(b.number))}
 if tsIsNumeric(a)&&b.kind==tsStringKind{return tsStrictEqual(a,tsNumberValue(tsNumber(b)))};if tsIsNumeric(b)&&a.kind==tsStringKind{return tsLooseEqual(b,a)}
 if tsIsPrimitiveValue(a)&&!tsNullish(a)&&!tsIsPrimitiveValue(b){return tsLooseEqual(a,tsToPrimitive(b))};if tsIsPrimitiveValue(b)&&!tsNullish(b)&&!tsIsPrimitiveValue(a){return tsLooseEqual(tsToPrimitive(a),b)}
 return false
}
func tsHas(value,key tsValue)bool{
 name:=tsPropertyKey(key)
 switch value.kind{
 case tsObjectKind:_,ok:=(*tsObject)(value.ref).values[name];return ok
 case tsClassKind:return tsHas(tsInstanceValue(tsInstanceProperties((*tsClass)(value.ref).static)),key)
 case tsInstanceKind:for p:=(*tsProperties)(value.ref);p!=nil;p=p.prototype{if _,ok:=p.declared[name];ok{return true};if _,ok:=p.extra[name];ok{return true};if _,ok:=p.methods[name];ok{return true}};return false
 case tsArrayKind:array:=(*tsArray)(value.ref);if name=="length"||name=="push"{return true};if _,ok:=array.properties[name];ok{return true};index,err:=strconv.Atoi(name);return err==nil&&index>=0&&index<len(array.values)&&strconv.Itoa(index)==name&&!array.holes[index]
 case tsECMAKind,tsRegExpKind:var object *goja.Object;if value.kind==tsECMAKind{object=(*tsECMAObject)(value.ref).object}else{object=(*tsRegExp)(value.ref).object};for ;object!=nil;object=object.Prototype(){for _,property:=range object.GetOwnPropertyNames(){if property==name{return true}}};return false
 }
 panic(tsThrown{tsErrorValue(&tsRuntimeError{"TypeError","Right operand of in must be an object"})})
}
`

const EqualityRuntime = `
func tsStrictEqual(a,b tsValue)bool {if tsIsNumeric(a)&&tsIsNumeric(b){return tsNumericEqual(a,b)};if a.kind!=b.kind{return false};switch a.kind {case tsUndefinedKind,tsNullKind:return true;case tsNumberKind,tsBooleanKind:return a.number==b.number;case tsStringKind:return tsStringEqual((*tsString)(a.ref),(*tsString)(b.ref));case tsRegExpKind:return (*tsRegExp)(a.ref).object==(*tsRegExp)(b.ref).object;case tsECMAKind:return (*tsECMAObject)(a.ref).object==(*tsECMAObject)(b.ref).object;case tsFunctionKind,tsObjectKind,tsArrayKind,tsClassKind,tsPromiseKind,tsInstanceKind:return a.ref==b.ref;default:return false}}

`
