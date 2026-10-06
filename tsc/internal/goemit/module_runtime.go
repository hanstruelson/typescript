package goemit

// ModuleRuntime implements demand-driven module evaluation. Export publication
// is a suspension point; each module owns one machine and never restarts it.
const ModuleRuntime = `
type tsExport struct { ready bool; value tsValue; binding tsBindingCell }
type tsImportRef struct { module *tsModule; name string }
type tsNamespace struct { module *tsModule }
type tsModuleRequest struct { missing map[string]struct{}; full bool; complete func(tsResult) }
type tsModule struct {
 name string
 loop *tsLoop
 initialize func(*tsModule)
 machine *tsMachine
 exports map[string]*tsExport
 requests []*tsModuleRequest
 queued, finished, failed bool
 failure tsValue
}
func tsNewModule(loop *tsLoop,name string) *tsModule {module:= &tsModule{name:name,loop:loop,exports:make(map[string]*tsExport)};loop.modules=append(loop.modules,module);return module}
func (module *tsModule) read(name string) tsValue {return module.readSeen(name,make(map[tsImportRef]bool))}
func (module *tsModule) readSeen(name string,seen map[tsImportRef]bool) tsValue {
 key:=tsImportRef{module,name};if seen[key] {panic("Circular export alias: "+module.name+":"+name)};seen[key]=true
 entry:=module.exports[name]
 if entry==nil || !entry.ready {panic("Export is not ready: "+module.name+":"+name)}
 value:=entry.value
 if entry.binding!=nil {if !entry.binding.initializedState() {panic("Export binding is uninitialized")};value=entry.binding.raw()}
 if ref,ok:=value.(tsImportRef);ok {return ref.module.readSeen(ref.name,seen)}
 return value
}
func (module *tsModule) cached(names []string,full bool)(tsResult,bool){
 if !full {all:=true;for _,name:=range names {entry:=module.exports[name];if entry==nil || !entry.ready {all=false;break}};if all {return tsResult{module,false},true}}
 if module.failed {return tsResult{module.failure,true},true}
 if full && !module.finished {return tsResult{},false}
 for _,name:=range names {entry:=module.exports[name];if entry==nil || !entry.ready {return tsResult{},false}}
 return tsResult{module,false},true
}
func (module *tsModule) request(names []string,full bool,complete func(tsResult)) {
 if result,ready:=module.cached(names,full);ready {complete(result);return}
 if module.finished {complete(tsResult{"Missing export in completed module: "+module.name,true});return}
 request:= &tsModuleRequest{missing:make(map[string]struct{}),full:full,complete:complete}
 for _,name:=range names {entry:=module.exports[name];if entry==nil || !entry.ready {request.missing[name]=struct{}{}}}
 module.requests=append(module.requests,request)
 module.schedule()
}
func (module *tsModule) schedule(){
 if module.queued || module.finished || module.failed || len(module.requests)==0 {return}
 // An import or await in progress resumes through its own completion callback.
 if module.machine!=nil && module.machine.blocked {return}
 module.queued=true
 module.loop.ready=append(module.loop.ready,func(){
  module.queued=false
  if module.finished || module.failed || len(module.requests)==0 {return}
  if module.machine==nil {module.initialize(module)}
  if !module.machine.blocked {module.machine.suspended=false;module.machine.resume()}
 })
}
func (module *tsModule) publish(name string,value tsValue,binding tsBindingCell){
 module.exports[name]= &tsExport{ready:true,value:value,binding:binding}
 remaining:=module.requests[:0]
 for _,request:=range module.requests {
  delete(request.missing,name)
  if !request.full && len(request.missing)==0 {callback:=request.complete;module.loop.ready=append(module.loop.ready,func(){callback(tsResult{module,false})})} else {remaining=append(remaining,request)}
 }
 module.requests=remaining
 // With no outstanding requests the remainder stays suspended until another
 // named import or a side-effect import asks for more execution.
 module.schedule()
}
func (module *tsModule) finish(result tsResult){
 if result.rejected {module.failed=true;module.failure=result.value} else {module.finished=true}
 for _,request:=range module.requests {
  callback:=request.complete;completion:=result
  if !result.rejected && len(request.missing)!=0 {completion=tsResult{"Missing export in completed module: "+module.name,true}}
  if !completion.rejected {completion.value=module}
  module.loop.ready=append(module.loop.ready,func(){callback(completion)})
 }
 module.requests=nil
}
func (m *tsMachine) importModule(module *tsModule,names []string,full bool,next int){
 if result,ready:=module.cached(names,full);ready {if result.rejected {m.raise(result.value)} else {m.result=result.value;m.pc=next};return}
 m.suspended=true;m.blocked=true
 module.request(names,full,func(result tsResult){
  // Queue even an immediate missing-export error to avoid reentrant resume.
  m.loop.ready=append(m.loop.ready,func(){m.blocked=false;m.suspended=false;if result.rejected {m.raise(result.value)} else {m.result=result.value;m.pc=next};m.resume()})
 })
}
func (m *tsMachine) exportValue(name string,value tsValue,binding tsBindingCell,next int){m.pc=next;m.suspended=true;m.module.publish(name,value,binding)}
`
