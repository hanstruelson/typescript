package goemit

// Runtime is emitted into the single output file; its dependencies use the compiler module.
// Only the event-loop goroutine touches promises, bindings, and pending tasks.
// Workers receive the completion channel and never run application code.
const Runtime = `import (
 "unsafe"
 "fmt"
 "math"
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

type tsUndefined struct{}

type tsCell struct { value tsValue; initialized, constant bool }
func tsBinding(initialized, constant bool) *tsCell { return &tsCell{tsU, initialized, constant} }
func (c *tsCell) get() tsValue { if !c.initialized { panic("Cannot access binding before initialization") }; if ref,ok:=c.value.(tsImportRef);ok {return ref.module.read(ref.name)};return c.value }
func (c *tsCell) init(value tsValue) tsValue { c.value = value; c.initialized = true; return value }
func (c *tsCell) set(value tsValue) tsValue { if !c.initialized { panic("Cannot access binding before initialization") }; if c.constant { panic("Assignment to constant variable") }; c.value = value; return value }
func tsClone(c *tsCell) *tsCell { value := *c; return &value }
type tsFunction struct {call func(...tsValue) tsValue}
func tsFunc(call func(...tsValue)tsValue)*tsFunction{return &tsFunction{call:call}}
func tsCall(value tsValue, args ...tsValue) tsValue { f, ok := value.(*tsFunction); if !ok { panic("Value is not callable") }; return f.call(args...) }
func tsArg(args []tsValue, index int) tsValue { if index < len(args) { return tsUnbox(args[index]) }; return tsU }
func tsText(value tsValue) string { switch v := value.(type) { case tsUndefined: return "undefined"; case nil: return "null"; default: return fmt.Sprint(v) } }
func tsConsoleLog(values ...tsValue) tsValue { for i, value := range values { if i != 0 { fmt.Print(" ") }; fmt.Print(tsText(value)) }; fmt.Println(); return tsU }
func tsIsUndefined(value tsValue) bool {_,ok:=value.(tsUndefined);return ok}
func tsNullish(value tsValue) bool { if value==nil {return true}; _,ok:=value.(tsUndefined);return ok }
type tsArray struct { values []tsValue;properties map[string]tsValue }
type tsIterator struct {array *tsArray; text *tsString; ecma *goja.Object; index int; value tsValue}
func tsIterate(value tsValue) *tsIterator {switch v:=value.(type) {case *tsArray:return &tsIterator{array:v};case *tsString:return &tsIterator{text:v};case *tsECMAObject:if v.object.ClassName()=="Array"{values:=[]tsValue{};for i:=int64(0);i<v.object.Get("length").ToInteger();i++ {values=append(values,tsFromECMA(v.object.Get(strconv.FormatInt(i,10))))};return &tsIterator{array:&tsArray{values:values}}};return &tsIterator{ecma:v.object};default:panic("Value is not iterable")}}
func (it *tsIterator) next() bool {if it.ecma!=nil {call,_:=goja.AssertFunction(it.ecma.Get("next"));result,err:=call(it.ecma);if err!=nil {tsECMAFailure(err)};object:=result.ToObject(tsEngine());if object.Get("done").ToBoolean(){return false};it.value=tsFromECMA(object.Get("value"));return true};if it.array!=nil {if it.index>=len(it.array.values){return false};it.value=it.array.values[it.index];it.index++;return true};if it.index>=len(it.text.units){return false};start:=it.index;it.index++;c:=it.text.units[start];if c>=0xd800&&c<=0xdbff&&it.index<len(it.text.units)&&it.text.units[it.index]>=0xdc00&&it.text.units[it.index]<=0xdfff {it.index++};it.value=it.text.slice(start,it.index);return true}
func tsGet(value,key tsValue) tsValue { switch v:=value.(type) {
 case *tsClass:return v.static.tsProperties().get(tsPropertyKey(key))
 case tsDynamicObject:return v.tsProperties().get(tsPropertyKey(key))
 case *tsObject:return v.get(tsPropertyKey(key))
 case *tsArray:
  if value,ok:=v.properties[tsPropertyKey(key)];ok {return value}
  if tsText(key)=="length" {return float64(len(v.values))}
  if tsText(key)=="push" {return tsFunc(func(args ...tsValue) tsValue {v.values=append(v.values,args...);return float64(len(v.values))})}
  index,err:=strconv.Atoi(tsPropertyKey(key));if err!=nil||strconv.Itoa(index)!=tsPropertyKey(key)||index<0 || index>=len(v.values) {return tsU};return v.values[index]
 case *tsString:return tsStringGet(v,key)
 case string:return tsStringGet(tsStringUTF8(v),key)
 case *tsRegExp:return tsECMAGet(v.object,key)
 case *tsECMAObject:return tsECMAGet(v.object,key)
 case *tsRuntimeError:if tsText(key)=="message" {return tsStringUTF8(v.message)};if tsText(key)=="name" {return tsStringUTF8(v.name)};return tsU
 case tsNamespace:return v.module.read(tsPropertyKey(key))
 case *tsPromise:
  if tsText(key)=="then" || tsText(key)=="catch" {return tsFunc(func(args ...tsValue) tsValue {
   next:=v.loop.promise();success,failure:=tsArg(args,0),tsArg(args,1);if tsText(key)=="catch" {failure=success;success=tsU}
   v.then(func(result tsResult){callback:=success;if result.rejected {callback=failure};if _,ok:=callback.(*tsFunction);!ok {next.settle(result);return};func(){defer func(){if err:=recover();err!=nil {next.settle(tsResult{tsUnwrap(err),true})}}();next.settle(tsResult{tsCall(callback,result.value),false})}()});return next
  })}
 }
 panic("Unsupported property access")
}
func tsSet(value,key,item tsValue) tsValue {switch v:=value.(type){case *tsClass:return v.static.tsProperties().set(tsPropertyKey(key),item);case *tsObject:return v.set(tsPropertyKey(key),item);case *tsRegExp:return tsECMASet(v.object,key,item);case *tsECMAObject:return tsECMASet(v.object,key,item)};if object,ok:=value.(tsDynamicObject);ok {return object.tsProperties().set(tsPropertyKey(key),item)}; array,ok:=value.(*tsArray);if !ok {panic("Expected an array")};name:=tsPropertyKey(key);if name=="length" {length:=tsNumber(item);if length<0||length>4294967295||math.Trunc(length)!=length {panic("Invalid array length")};for len(array.values)<int(length){array.values=append(array.values,tsU)};array.values=array.values[:int(length)];return item};index,err:=strconv.Atoi(name);if err!=nil || index<0 || strconv.Itoa(index)!=name {if array.properties==nil{array.properties=map[string]tsValue{}};array.properties[name]=item;return item};for len(array.values)<=index {array.values=append(array.values,tsU)};array.values[index]=item;return item }
func tsNumber(value tsValue) float64 { switch v := value.(type) { case *tsString:return goja.StringFromUTF16(v.units).ToFloat();case string:return tsNumber(tsStringUTF8(v));case tsUndefined:return math.NaN();case float64: return v; case bool: if v { return 1 }; return 0; case nil: return 0; default: panic("Expected a number") } }
func tsTruthy(value tsValue) bool { switch v := value.(type) { case tsUndefined: return false; case nil: return false; case bool: return v; case float64: return v != 0 && !math.IsNaN(v); case *tsString:return len(v.units)!=0;case string: return v != ""; default: return true } }
func tsBinary(op string, left, right tsValue) tsValue {
 switch op {
 case "+": if _,ok:=left.(*tsString);ok{return tsStringConcat(tsStringValue(left),tsStringValue(right))};if _,ok:=right.(*tsString);ok{return tsStringConcat(tsStringValue(left),tsStringValue(right))}; return tsNumber(left)+tsNumber(right)
 case "-": return tsNumber(left)-tsNumber(right)
 case "*": return tsNumber(left)*tsNumber(right)
 case "/": return tsNumber(left)/tsNumber(right)
 case "%": return math.Mod(tsNumber(left),tsNumber(right))
 case "<":if a,ok:=left.(*tsString);ok {if b,ok:=right.(*tsString);ok{return tsStringCompare(a,b)<0}}; return tsNumber(left)<tsNumber(right)
 case "<=":if a,ok:=left.(*tsString);ok {if b,ok:=right.(*tsString);ok{return tsStringCompare(a,b)<=0}}; return tsNumber(left)<=tsNumber(right)
 case ">":if a,ok:=left.(*tsString);ok {if b,ok:=right.(*tsString);ok{return tsStringCompare(a,b)>0}}; return tsNumber(left)>tsNumber(right)
 case ">=":if a,ok:=left.(*tsString);ok {if b,ok:=right.(*tsString);ok{return tsStringCompare(a,b)>=0}}; return tsNumber(left)>=tsNumber(right)
 case "===", "!==": equal:=tsStrictEqual(left,right); if op=="!==" { return !equal }; return equal
 default: panic("Unsupported operator")
 }
}

type tsThrown struct {value tsValue}
func tsUnwrap(value any) tsValue {if thrown,ok:=value.(tsThrown);ok {return thrown.value};return value}
type tsResult struct { value tsValue; rejected bool }
type tsPromise struct { loop *tsLoop; settled, handled bool; result tsResult; listeners []func(tsResult) }
// A task is registered and finished by the loop. The worker exclusively writes
// result until publishing this pointer through the completion channel.
type tsTask struct { work func() tsResult; finish func(tsResult); result tsResult; canceled bool; cancel func(); timer, completed bool; deadline time.Time }
type tsLoop struct { pending map[*tsTask]struct{}; completions chan *tsTask; ready []func(); promises []*tsPromise; errors []tsValue; timers []*tsTask; modules []*tsModule }
func tsNewLoop() *tsLoop { return &tsLoop{pending:make(map[*tsTask]struct{}), completions:make(chan *tsTask)} }
func (l *tsLoop) promise() *tsPromise { p:= &tsPromise{loop:l}; l.promises=append(l.promises,p); return p }
func (p *tsPromise) settle(result tsResult) {
 if p.settled { return }
 if other, ok := result.value.(*tsPromise); ok && !result.rejected {
  if other==p { p.settle(tsResult{"Promise cannot resolve to itself",true}); return }
  other.then(p.settle); return
 }
 p.settled=true; p.result=result
 for _, callback := range p.listeners { cb:=callback; p.loop.ready=append(p.loop.ready,func(){cb(result)}) }; p.listeners=nil
}
func (p *tsPromise) then(callback func(tsResult)) { p.handled=true; if p.settled { result:=p.result; p.loop.ready=append(p.loop.ready,func(){callback(result)}) } else { p.listeners=append(p.listeners,callback) } }
func (l *tsLoop) resolved(value tsValue, rejected bool) *tsPromise { p:=l.promise(); p.settle(tsResult{value,rejected}); return p }
func (l *tsLoop) construct(executor tsValue) *tsPromise {
 p:=l.promise();finished:=false
 resolve:=tsFunc(func(args ...tsValue) tsValue {if !finished {finished=true;p.settle(tsResult{tsArg(args,0),false})};return tsU})
 reject:=tsFunc(func(args ...tsValue) tsValue {if !finished {finished=true;p.settle(tsResult{tsArg(args,0),true})};return tsU})
 func(){defer func(){if err:=recover();err!=nil && !finished {finished=true;p.settle(tsResult{tsUnwrap(err),true})}}();tsCall(executor,resolve,reject)}()
 return p
}
func (l *tsLoop) all(values tsValue) *tsPromise {
 p:=l.promise();array,ok:=values.(*tsArray);if !ok {p.settle(tsResult{"Promise.all expects an array",true});return p}
 output:= &tsArray{values:make([]tsValue,len(array.values))};remaining:=len(array.values)
 if remaining==0 {p.settle(tsResult{output,false});return p}
 for index,value:=range array.values {i:=index;l.await(value,func(result tsResult){if result.rejected {p.settle(result);return};output.values[i]=result.value;remaining--;if remaining==0 {p.settle(tsResult{output,false})}})}
 return p
}
func (l *tsLoop) await(value tsValue, callback func(tsResult)) { if p,ok:=value.(*tsPromise); ok { p.then(callback) } else { l.ready=append(l.ready,func(){callback(tsResult{value,false})}) } }
// submit is called only on the loop. Native adapters must snapshot inputs and
// must not capture application cells or mutate application-owned values.
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
 if _,ok:=callback.(*tsFunction);!ok {panic("setTimeout expects a function")}
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
func (l *tsLoop) clearTimeout(handle tsValue) tsValue { if task,ok:=handle.(*tsTask);ok && !task.canceled {task.canceled=true;if task.cancel!=nil {task.cancel()};if task.completed {l.removeTimer(task);delete(l.pending,task)}};return tsU }
func (l *tsLoop) readFile(path tsValue) *tsPromise { text,ok:=path.(*tsString); if !ok { return l.resolved("readFile expects a string",true) }; return l.start(func() tsResult { bytes,err:=os.ReadFile(text.String()); if err!=nil {return tsResult{err.Error(),true}}; return tsResult{tsStringUTF8(string(bytes)),false} }) }
func (l *tsLoop) delay(milliseconds tsValue) *tsPromise { duration:=tsTimerDuration(milliseconds); return l.start(func() tsResult {time.Sleep(duration); return tsResult{tsU,false}}) }
func (l *tsLoop) invoke(callback func()) {
 defer func(){if failure:=recover();failure!=nil {l.errors=append(l.errors,tsUnwrap(failure))}}()
 callback()
}
func (l *tsLoop) removeTimer(task *tsTask) {for i,timer:=range l.timers {if timer==task {copy(l.timers[i:],l.timers[i+1:]);l.timers[len(l.timers)-1]=nil;l.timers=l.timers[:len(l.timers)-1];return}}}
func (l *tsLoop) run() error {
 for {
  for len(l.ready)>0 {callback:=l.ready[0];l.ready[0]=nil;l.ready=l.ready[1:];l.invoke(callback)}
  // Process one timer at a time, with a reaction checkpoint between callbacks.
  // Deadline ordering prevents goroutine scheduling from reordering equal timers.
  if len(l.timers)>0 && l.timers[0].completed {
   task:=l.timers[0];l.removeTimer(task);delete(l.pending,task);l.invoke(func(){task.finish(task.result)});continue
  }
  if len(l.pending)==0 {break}
  task:=<-l.completions
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
 for _,p:=range l.promises {if p.settled && p.result.rejected && !p.handled {l.errors=append(l.errors,"Unhandled rejection: "+tsText(p.result.value))}}
 if len(l.errors)>0 {messages:=make([]string,len(l.errors));for i,value:=range l.errors {messages[i]=tsText(value)};return fmt.Errorf("Async runtime errors: %s",strings.Join(messages,"; "))}
 return nil
}

// Protected regions and pending abrupt completions mirror ES5's trys/ops stacks.
type tsHandler struct { catch, finally, end int; phase int }
type tsAbrupt struct { kind string; value tsValue; target, depth int }
type tsMachine struct { loop *tsLoop; pc int; step func(); output *tsPromise; async, done, suspended bool; result tsValue; handlers []tsHandler; abrupt []tsAbrupt; thrown tsValue; module *tsModule; blocked bool }
func (m *tsMachine) resume() {
 for !m.done && !m.suspended {
  func(){ defer func(){if err:=recover(); err!=nil {m.raise(err)}}(); m.step() }()
 }
}
func (m *tsMachine) await(value tsValue, next int) {
 m.suspended=true;m.blocked=true
 m.loop.await(value,func(result tsResult){m.blocked=false;m.suspended=false; if result.rejected {m.raise(result.value)} else {m.result=result.value; m.pc=next}; m.resume()})
}
func (m *tsMachine) complete(value tsValue) { m.done=true; m.result=value; if m.module!=nil {m.module.finish(tsResult{value,false})}; if m.async {m.output.settle(tsResult{value,false})} }
func (m *tsMachine) raise(value any) { m.transfer(tsAbrupt{kind:"throw",value:tsUnwrap(value)}) }
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
`

const EqualityRuntime = `
func tsStrictEqual(a,b tsValue)bool {switch left:=a.(type){case *tsString:right,ok:=b.(*tsString);return ok&&tsStringEqual(left,right);case float64:right,ok:=b.(float64);return ok&&left==right;case bool:right,ok:=b.(bool);return ok&&left==right;case nil:return b==nil;case tsUndefined:_,ok:=b.(tsUndefined);return ok;case *tsFunction:right,ok:=b.(*tsFunction);return ok&&left==right;case *tsObject:right,ok:=b.(*tsObject);return ok&&left==right;case *tsRegExp:right,ok:=b.(*tsRegExp);return ok&&left.object==right.object;case *tsECMAObject:right,ok:=b.(*tsECMAObject);return ok&&left.object==right.object;case *tsArray:right,ok:=b.(*tsArray);return ok&&left==right;case *tsClass:right,ok:=b.(*tsClass);return ok&&left==right;case *tsPromise:right,ok:=b.(*tsPromise);return ok&&left==right;case tsDynamicObject:right,ok:=b.(tsDynamicObject);return ok&&left==right};return false}
`
