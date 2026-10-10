package compiler_test

import (
	"os/exec"
	"strings"
	"testing"
)

func TestGoMathErrorMilestoneDifferential(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("Node unavailable")
	}
	source := `const m=Math;const round=m.round;console.log(m===Math,round===Math.round,Math.abs(-4),Math.imul(4294967295,5),Math.clz32(1));
 console.log(Object.is(Math.round(-0.5),-0),Math.round(0.49999999999999994),Math.max(),Math.min(),Math.hypot(3,4),Math.pow(NaN,0),Math.pow(1,Infinity));
 console.log(Math.fround(16777217),Math.sign(-3),Math.trunc(-3.8),Math.atan2(0,-0));
 console.log(Object.keys(Math).length,Math.round.length,Object.getOwnPropertyDescriptor(Math,"PI").writable);
 for(const C of [Error,TypeError,RangeError,ReferenceError,SyntaxError,URIError,EvalError]){const e=new C("bad",{cause:7});console.log(e.name,e.message,e.cause,e instanceof C,e instanceof Error,e.toString(),Object.keys(e).length);}
 const a=AggregateError([1,2],"many");console.log(a.name,a.errors.length,a.toString(),a instanceof Error);
 console.log(Error.prototype.toString.call({name:"",message:"x"}),Error.prototype.toString.call({name:"X",message:""}));
 try{Math.abs(1n);}catch(e){console.log(e instanceof TypeError,e.name);}`
	want, err := exec.Command(node, "-e", source).CombinedOutput()
	if err != nil {
		t.Fatalf("Node: %v\n%s", err, want)
	}
	got, err := runGoProgram(t, emitGoProgram(t, source), false)
	if err != nil {
		t.Fatalf("Go: %v\n%s", err, got)
	}
	if got != string(want) {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestGoDeleteMilestoneDifferential(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("Node unavailable")
	}
	source := `const proto={x:8};const o=Object.create(proto);o.x=undefined;console.log(Object.hasOwn(o,"x"),"x" in o,delete o.x,o.x,Object.hasOwn(o,"x"));o.x=4;console.log(Object.keys(o).join(","));
 Object.defineProperty(o,"fixed",{value:2});console.log(delete o.fixed,o.fixed);function strictDelete(){"use strict";try{delete o.fixed;}catch(e){console.log(e instanceof TypeError);}}strictDelete();
 class Person{x:number=3;y:number=7;}const p=new Person();console.log(delete p.x,p.x,"x" in p,p.y,Object.keys(p).join(","));p.x=9;console.log(p.x,"x" in p,Object.keys(p).join(","));
 const array=[1,"two",3];console.log(delete array[1],array.length,1 in array,array[1]);console.log(delete array.length);console.log(delete "abc"[0],delete "abc"[3]);`
	want, err := exec.Command(node, "-e", strings.ReplaceAll(source, ":number", "")).CombinedOutput()
	if err != nil {
		t.Fatalf("Node: %v\n%s", err, want)
	}
	got, err := runGoProgram(t, emitGoProgram(t, source), false)
	if err != nil {
		t.Fatalf("Go: %v\n%s", err, got)
	}
	if got != string(want) {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestGoArgumentsMilestoneDifferential(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("Node unavailable")
	}
	source := `function mapped(a){console.log(arguments.length,arguments[0],arguments.callee===mapped);a=4;console.log(arguments[0]);arguments[0]=6;console.log(a);const capture=()=>arguments[0];console.log(capture());delete arguments[0];a=8;console.log(arguments[0],Object.hasOwn(arguments,"0"));}mapped(2,3);
 function strict(a){"use strict";a=4;console.log(arguments[0]);arguments[0]=6;console.log(a);try{arguments.callee;}catch(e){console.log(e instanceof TypeError);}}strict(2);
 function defaults(a=7){a=8;console.log(arguments.length,arguments[0]);}defaults();defaults(2);
 function redefine(a){Object.defineProperty(arguments,"0",{writable:false});a=9;console.log(arguments[0],Object.getOwnPropertyDescriptor(arguments,"0").value);}redefine(2);`
	want, err := exec.Command(node, "-e", source).CombinedOutput()
	if err != nil {
		t.Fatalf("Node: %v\n%s", err, want)
	}
	got, err := runGoProgram(t, emitGoProgram(t, source), false)
	if err != nil {
		t.Fatalf("Go: %v\n%s", err, got)
	}
	if got != string(want) {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestGoDeleteKeepsUnrelatedNativeField(t *testing.T) {
	source := `class Point{x:number=3;y:number=7;}const p=new Point();delete p.x;console.log(p.y);`
	text := emitGoProgram(t, source)
	if !strings.Contains(text, ".P79\n") && !strings.Contains(text, ".P79;") {
		t.Fatal("unaffected y field did not retain a native load")
	}
	got, err := runGoProgram(t, text, false)
	if err != nil || got != "7\n" {
		t.Fatalf("%v %s", err, got)
	}
}

func TestGoGeneratorMilestoneDifferential(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("Node unavailable")
	}
	source := `let called=0;function* gen(a){called++;try{const x=yield a;yield x+1;return 9;}finally{console.log("finally");}}const g=gen(3);console.log(called);let r=g.next();console.log(r.value,r.done,called);r=g.next(5);console.log(r.value,r.done);r=g.next();console.log(r.value,r.done);console.log(g.next().value,g.next().done);
 function* cleanup(){try{yield 1;}finally{yield 2;}}const c=cleanup();console.log(c.next().value);r=c.return(7);console.log(r.value,r.done);r=c.next();console.log(r.value,r.done);
 function* caught(){try{yield 1;}catch(e){yield e;}return 4;}const t=caught();t.next();console.log(t.throw("x").value,t.next().value);
 function* inner(){yield 2;return 8;}function* outer(){yield 1;const x=yield* inner();yield x;yield* [3,4];}let text="";for(const value of outer()){text+=value+",";}console.log(text);`
	want, err := exec.Command(node, "-e", source).CombinedOutput()
	if err != nil {
		t.Fatalf("Node: %v\n%s", err, want)
	}
	got, err := runGoProgram(t, emitGoProgram(t, source), false)
	if err != nil {
		t.Fatalf("Go: %v\n%s", err, got)
	}
	if got != string(want) {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestGoComputedClassesMilestoneDifferential(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("Node unavailable")
	}
	source := `let key="x";let count=0;function fieldName(){count++;return key;}class Base{[fieldName()]=3;["read"](){return this.x;}static ["answer"]=7;}key="y";const a=new Base();const b=new Base();console.log(count,a.x,b.x,a.read(),Base.answer,Object.keys(a).join(","));class Child extends Base{[fieldName()]=9;}const c=new Child();console.log(count,c.x,c.y,c.read(),Object.keys(c).join(","));class G{*values(){yield 1;yield 2;}}let text="";for(const v of new G().values()){text+=v;}console.log(text);`
	want, err := exec.Command(node, "-e", source).CombinedOutput()
	if err != nil {
		t.Fatalf("Node: %v\n%s", err, want)
	}
	got, err := runGoProgram(t, emitGoProgram(t, source), false)
	if err != nil {
		t.Fatalf("Go: %v\n%s", err, got)
	}
	if got != string(want) {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestGoSymbolMilestoneDifferential(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("Node unavailable")
	}
	source := `const a=Symbol("x"),b=Symbol("x");console.log(typeof a,a===b,a.description,Symbol().description,Symbol("").description,String(a),a.toString());console.log(Symbol.for("shared")===Symbol.for("shared"),Symbol.keyFor(Symbol.for("shared")),Symbol.keyFor(a));
 const o={[a]:3,[b]:4,normal:7};console.log(o[a],o[b],Object.keys(o).join(","),Object.getOwnPropertyNames(o).join(","),Object.getOwnPropertySymbols(o).length,Reflect.ownKeys(o).length,JSON.stringify(o));console.log(Object.getOwnPropertySymbols(o)[0]===a,Object.getOwnPropertySymbols(o)[1]===b);console.log(delete o[a],o[a],o[b]);
 Object.defineProperty(o,a,{value:9,configurable:true});console.log(o[a],Object.getOwnPropertyDescriptor(o,a).enumerable);console.log(Reflect.deleteProperty(o,a),o[a]);
 try{Math.abs(a);}catch(e){console.log(e instanceof TypeError);}try{""+a;}catch(e){console.log(e instanceof TypeError);}console.log(JSON.stringify([a]),JSON.stringify(a));
 class C{[a]=5;}const c=new C();console.log(c[a],Object.getOwnPropertySymbols(c)[0]===a,Object.keys(c).length);`
	want, err := exec.Command(node, "-e", source).CombinedOutput()
	if err != nil {
		t.Fatalf("Node: %v\n%s", err, want)
	}
	got, err := runGoProgram(t, emitGoProgram(t, source), false)
	if err != nil {
		t.Fatalf("Go: %v\n%s", err, got)
	}
	if got != string(want) {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestGoAsyncGeneratorMilestoneDifferential(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("Node unavailable")
	}
	source := `async function* gen(){yield Promise.resolve(1);await Promise.resolve(0);yield 2;return 3;}async function main(){const g=gen();const a=g.next(),b=g.next(),c=g.next();for(const r of await Promise.all([a,b,c])){console.log(r.value,r.done);}let text="";for await(const v of gen()){text+=v;}console.log(text);async function* bad(){yield 1;throw new TypeError("bad");}const t=bad();await t.next();try{await t.next();}catch(e){console.log(e.name,e.message);}const r=await t.next();console.log(r.done,r.value);}main();`
	want, err := exec.Command(node, "-e", source).CombinedOutput()
	if err != nil {
		t.Fatalf("Node: %v\n%s", err, want)
	}
	got, err := runGoProgram(t, emitGoProgram(t, source), false)
	if err != nil {
		t.Fatalf("Go: %v\n%s", err, got)
	}
	if got != string(want) {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestGoReceiverAndSymbolProtocolDifferential(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("Node unavailable")
	}
	source := `function sloppy(){return this;}function strict(){"use strict";return this;}console.log(sloppy()===globalThis,strict()===undefined,sloppy.call(null)===globalThis,typeof sloppy.call(3),sloppy.call(3).valueOf(),strict.call(3));const o={x:9,f:function(){return this.x;}};console.log(o.f(),o.f.call({x:4}));
 const value={[Symbol.toPrimitive](hint){return hint==="string"?"key":10;}};const object={[value]:3};console.log(object.key,value+2);
 const iterable={[Symbol.iterator]:function*(){yield 1;yield 2;}};let text="";for(const v of iterable){text+=v;}console.log(text);`
	want, err := exec.Command(node, "-e", source).CombinedOutput()
	if err != nil {
		t.Fatalf("Node: %v\n%s", err, want)
	}
	got, err := runGoProgram(t, emitGoProgram(t, source), false)
	if err != nil {
		t.Fatalf("Go: %v\n%s", err, got)
	}
	if got != string(want) {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestGoMathPrecisionAndNativeLowering(t *testing.T) {
	source := `console.log(Math.sumPrecise([1e30,0.1,-1e30]),Math.sumPrecise([1e308,1e308,-1e308,-1e308,1]));console.log(Object.is(Math.sumPrecise([]),-0),Object.is(Math.sumPrecise([-0]),-0),Object.is(Math.sumPrecise([-0,0]),0));console.log(Math.f16round(1.00048828125),Math.f16round(65520),Object.is(Math.f16round(-1e-20),-0));console.log(Math.max(NaN,Infinity),Math.min(NaN,-Infinity));`
	got, err := runGoProgram(t, emitGoProgram(t, source), false)
	want := "0.1 1\ntrue true true\n1 Infinity true\nNaN NaN\n"
	if err != nil || got != want {
		t.Fatalf("%v\ngot %s\nwant %s", err, got, want)
	}
	text := emitGoProgram(t, `function root(value:number):number{return Math.sqrt(value);}console.log(root(9));`)
	if !strings.Contains(text, "math.Sqrt(") {
		t.Fatal("known Math.sqrt lost native lowering")
	}
	got, err = runGoProgram(t, text, false)
	if err != nil || got != "3\n" {
		t.Fatalf("%v %s", err, got)
	}
}

func TestGoOptionalDeleteAndArgumentsDefaults(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("Node unavailable")
	}
	source := `let count=0;function key(){count++;return "x";}const empty=null;console.log(delete empty?.[key()],count);const value={get x(){throw 7;}};console.log(delete value?.x,Object.hasOwn(value,"x"));function args(a=arguments.length){console.log(a,arguments.length);arguments=7;console.log(arguments);}args();function capture(a){const f=()=>arguments;arguments=8;console.log(f());}capture(1);`
	want, err := exec.Command(node, "-e", source).CombinedOutput()
	if err != nil {
		t.Fatalf("Node: %v\n%s", err, want)
	}
	got, err := runGoProgram(t, emitGoProgram(t, source), false)
	if err != nil {
		t.Fatalf("Go: %v\n%s", err, got)
	}
	if got != string(want) {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestGoClassAccessorsAndDeleteReflectionDifferential(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("Node unavailable")
	}
	source := `const key="value";class Base{x=3;get [key](){return this.x;}set [key](v){this.x=v;}static get ["answer"](){return 7;}}const b=new Base();console.log(b.value);b.value=9;console.log(b.x,Base.answer,Object.keys(b).join(","));console.log(Object.getOwnPropertyDescriptor(Base.prototype,key).enumerable);console.log(Object.getOwnPropertyDescriptor(b,"x").value);Object.defineProperty(b,"x",{value:12,configurable:false,writable:true,enumerable:false});console.log(b.x,Object.keys(b).length,Reflect.deleteProperty(b,"x"));b.x=13;console.log(b.value);`
	want, err := exec.Command(node, "-e", source).CombinedOutput()
	if err != nil {
		t.Fatalf("Node: %v\n%s", err, want)
	}
	got, err := runGoProgram(t, emitGoProgram(t, source), false)
	if err != nil {
		t.Fatalf("Go: %v\n%s", err, got)
	}
	if got != string(want) {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestGoGeneratorAbruptClosingDifferential(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("Node unavailable")
	}
	source := `function* gen(){try{yield 1;yield 2;}finally{console.log("closed");}}for(const v of gen()){console.log(v);break;}function early(){for(const v of gen()){return v;}}console.log(early());async function* asyncGen(){try{yield Promise.reject("bad");}catch(e){yield e;}}async function main(){const r=await asyncGen().next();console.log(r.value,r.done);}main();`
	want, err := exec.Command(node, "-e", source).CombinedOutput()
	if err != nil {
		t.Fatalf("Node: %v\n%s", err, want)
	}
	got, err := runGoProgram(t, emitGoProgram(t, source), false)
	if err != nil {
		t.Fatalf("Go: %v\n%s", err, got)
	}
	if got != string(want) {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestGoGeneratorDelegationIdentityDifferential(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("Node unavailable")
	}
	source := `function* gen(){yield 1;}const a=gen(),b=gen();console.log(a.next===b.next,Object.getPrototypeOf(a)===gen.prototype);try{new gen();}catch(e){console.log(e instanceof TypeError);}const next=a.next;try{next();}catch(e){console.log(e instanceof TypeError);}console.log(next.call(b).value);async function* inner(){yield Promise.resolve(2);return 7;}async function* outer(){const value=yield* inner();yield value;}async function main(){let text="";for await(const value of outer()){text+=value+",";}console.log(text);}main();`
	want, err := exec.Command(node, "-e", source).CombinedOutput()
	if err != nil {
		t.Fatalf("Node: %v\n%s", err, want)
	}
	got, err := runGoProgram(t, emitGoProgram(t, source), false)
	if err != nil {
		t.Fatalf("Go: %v\n%s", err, got)
	}
	if got != string(want) {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestGoDeletedNumericFieldUsesUndefinedSemantics(t *testing.T) {
	got, err := runGoProgram(t, emitGoProgram(t, `class C{x=3;y=4;}const c=new C();delete c.x;console.log(c.x+1,c.y+1,c.x===undefined);`), false)
	if err != nil || got != "NaN 5 true\n" {
		t.Fatalf("%v %s", err, got)
	}
}

func TestGoErrorSubclassDifferential(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("Node unavailable")
	}
	source := `class Problem extends TypeError{code=7;constructor(message){super(message,{cause:3});}read(){return this.message;}}const p=new Problem("bad");console.log(p instanceof Problem,p instanceof TypeError,p instanceof Error,p.name,p.message,p.cause,p.code,p.read(),p.toString(),Object.keys(p).join(","));class Default extends Error{}const d=new Default("default");console.log(d instanceof Error,d.message,d.toString());console.log(Object.prototype.toString.call(new Error()));`
	want, err := exec.Command(node, "-e", source).CombinedOutput()
	if err != nil {
		t.Fatalf("Node: %v\n%s", err, want)
	}
	got, err := runGoProgram(t, emitGoProgram(t, source), false)
	if err != nil {
		t.Fatalf("Go: %v\n%s", err, got)
	}
	if got != string(want) {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestGoArgumentsAndSymbolCopyDifferential(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("Node unavailable")
	}
	source := `function f(a){a=7;console.log([...arguments].join(","),Object.getOwnPropertySymbols(arguments)[0]===Symbol.iterator,Object.prototype.toString.call(arguments));}f(1,2);const key=Symbol("key"),hidden=Symbol("hidden");const a={[key]:3,x:4};Object.defineProperty(a,hidden,{value:9});const b=Object.assign({},a);const c={...a};const {[key]:removed,...rest}=a;console.log(b[key],c[key],removed,rest[key],rest.x,Object.getOwnPropertySymbols(b).length,Object.getOwnPropertySymbols(c).length);const descriptors={[key]:{value:8,enumerable:true}};const d=Object.defineProperties({},descriptors);console.log(d[key]);`
	want, err := exec.Command(node, "-e", source).CombinedOutput()
	if err != nil {
		t.Fatalf("Node: %v\n%s", err, want)
	}
	got, err := runGoProgram(t, emitGoProgram(t, source), false)
	if err != nil {
		t.Fatalf("Go: %v\n%s", err, got)
	}
	if got != string(want) {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestGoAsyncFromSyncIterationDifferential(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("Node unavailable")
	}
	source := `function* items(){try{yield Promise.resolve(3);yield Promise.resolve(4);}finally{console.log("closed");}}async function main(){for await(const x of items()){console.log(x);break;}async function* outer(){yield* items();}for await(const x of outer()){console.log("delegated",x);}}main();`
	want, err := exec.Command(node, "-e", source).CombinedOutput()
	if err != nil {
		t.Fatalf("Node: %v\n%s", err, want)
	}
	got, err := runGoProgram(t, emitGoProgram(t, source), false)
	if err != nil {
		t.Fatalf("Go: %v\n%s", err, got)
	}
	if got != string(want) {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestGoAsyncGeneratorReturnAwaitDifferential(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("Node unavailable")
	}
	source := `async function* g(){try{yield 1;}catch(e){console.log("caught",e);yield 2;}finally{console.log("finally");}}async function main(){const a=g();await a.next();const result=await a.return(Promise.reject("bad"));console.log(result.value,result.done);console.log((await a.next()).done);const b=g();await b.next();console.log((await b.return(Promise.resolve(9))).value);}main();`
	want, err := exec.Command(node, "-e", source).CombinedOutput()
	if err != nil {
		t.Fatalf("Node: %v\n%s", err, want)
	}
	got, err := runGoProgram(t, emitGoProgram(t, source), false)
	if err != nil {
		t.Fatalf("Go: %v\n%s", err, got)
	}
	if got != string(want) {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestGoComputedFieldCollisionDifferential(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("Node unavailable")
	}
	source := `let key="x";class C{[key]=3;x=this.x+1;[key]=this.x+2;y=this.x;}const c=new C();console.log(c.x,c.y,Object.keys(c).join(","));class Base{[key]=5;}class Child extends Base{[key]=this.x+1;y=this.x;}const d=new Child();console.log(d.x,d.y);class S{static [key]=3;static x=this.x+1;static [key]=this.x+2;}console.log(S.x);class Access{static get answer(){return 7;}static method(){return 8;}}console.log(Object.keys(Access).length,Object.getOwnPropertyDescriptor(Access,"method").value===Access.method,Object.getOwnPropertyDescriptor(Access,"answer").enumerable);`
	want, err := exec.Command(node, "-e", source).CombinedOutput()
	if err != nil {
		t.Fatalf("Node: %v\n%s", err, want)
	}
	got, err := runGoProgram(t, emitGoProgram(t, source), false)
	if err != nil {
		t.Fatalf("Go: %v\n%s", err, got)
	}
	if got != string(want) {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestGoSymbolPropertyKeyRetainsAllocation(t *testing.T) {
	text := emitGoProgram(t, `console.log("done");`)
	text += `func init(){loop:=tsNewLoop();holder:=tsNewObject();func(){symbol:=tsNewSymbol(tsStringUTF8("retained"));holder.set(loop,(*tsSymbol)(symbol.ref).key,tsNumberValue(7))}();for i:=0;i<8;i++{runtime.GC()};symbol:=tsKeyValue(holder.order[0]);if (*tsSymbol)(symbol.ref).description.String()!="retained"{panic("symbol key lost its allocation")};runtime.KeepAlive(holder)}`
	got, err := runGoProgram(t, text, false)
	if err != nil || got != "done\n" {
		t.Fatalf("%v\n%s", err, got)
	}
}

func TestGoComputedClassNameTDZDifferential(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("Node unavailable")
	}
	source := `try{class C{[C]=1;}}catch(e){console.log(e instanceof ReferenceError);}class Big{value=1n;add(){this.value+=2n;return this.value;}}const b=new Big();console.log(String(b.add()),typeof b.value);`
	want, err := exec.Command(node, "-e", source).CombinedOutput()
	if err != nil {
		t.Fatalf("Node: %v\n%s", err, want)
	}
	got, err := runGoProgram(t, emitGoProgram(t, source), false)
	if err != nil {
		t.Fatalf("Go: %v\n%s", err, got)
	}
	if got != string(want) {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestGoNativeProtocolEdgesDifferential(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("Node unavailable")
	}
	source := `function loop(){var keys=["x","y"];for(var i=0;i<keys.length;i++){console.log(keys[i]);}}loop();function args(){console.log(delete arguments,arguments.length);}args(1);function* g(){}async function* a(){}const root=Object.getPrototypeOf(g),asyncRoot=Object.getPrototypeOf(a);console.log(typeof root,root.constructor.name,asyncRoot.constructor.name,root.prototype.next.name,asyncRoot.prototype.next.name,Object.prototype.toString.call(g),Object.prototype.toString.call(a));try{Reflect.construct(function(){},[],Math.max);}catch(e){console.log(e instanceof TypeError);}let reads=0;const result={done:false,get value(){reads++;return 7;}};const obj={[Symbol.iterator](){return {next(){return result;},return(){return result;}};}};function* delegate(){yield* obj;}const it=delegate();const first=it.next();console.log(first===result,reads);it.return();console.log(reads);result.done=true;it.return();console.log(reads);`
	want, err := exec.Command(node, "-e", source).CombinedOutput()
	if err != nil {
		t.Fatalf("Node: %v\n%s", err, want)
	}
	got, err := runGoProgram(t, emitGoProgram(t, source), false)
	if err != nil {
		t.Fatalf("Go: %v\n%s", err, got)
	}
	if got != string(want) {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestGoClassExpressionDifferential(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("Node unavailable")
	}
	source := `const C=class {x=3;read(){return this.x;}};const c=new C();console.log(c.read(),c instanceof C);const D=class Inner extends C{read(){return super.read()+1;}self(){return this instanceof Inner;}};const d=new D();console.log(d.read(),d.self(),typeof Inner);function factory(value){return class {x=value;read(){return this.x;}};}const A=factory(5),B=factory(7);console.log(new A().read(),new B().read());const G=class {*items(){yield 3;}};console.log(new G().items().next().value);`
	want, err := exec.Command(node, "-e", source).CombinedOutput()
	if err != nil {
		t.Fatalf("Node: %v\n%s", err, want)
	}
	got, err := runGoProgram(t, emitGoProgram(t, source), false)
	if err != nil {
		t.Fatalf("Go: %v\n%s", err, got)
	}
	if got != string(want) {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestGoReflectReceiverAndFailureDifferential(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("Node unavailable")
	}
	source := `const proto={get value(){return this.x;},set value(v){this.x=v;}};const receiver={x:3};console.log(Reflect.get(proto,"value",receiver),Reflect.set(proto,"value",7,receiver),receiver.x);const target={};Object.defineProperty(target,"fixed",{value:1});console.log(Reflect.set(target,"fixed",2),Reflect.defineProperty(target,"fixed",{value:2}),target.fixed);try{Object.defineProperty(target,"fixed",{value:2});}catch(e){console.log(e instanceof TypeError);}const other={};console.log(Reflect.set({data:4},"data",8,other),other.data);Object.preventExtensions(other);console.log(Reflect.set({},"missing",2,other),Reflect.defineProperty(other,"missing",{value:3}));const throwing={set value(v){throw "setter";}};try{Reflect.set(throwing,"value",1,receiver);}catch(e){console.log(e);}let called=0;const ownSetter={set data(v){called++;}};console.log(Reflect.set({data:1},"data",2,ownSetter),called);const locked={};const p={};Object.preventExtensions(locked);console.log(Reflect.setPrototypeOf(locked,p),Reflect.setPrototypeOf(locked,Object.getPrototypeOf(locked)));const cycle={};console.log(Reflect.setPrototypeOf(cycle,cycle));try{Reflect.setPrototypeOf({},42);}catch(e){console.log(e instanceof TypeError);}const f=function(){};Object.preventExtensions(f);console.log(Reflect.setPrototypeOf(f,p));const list=[];console.log(Reflect.setPrototypeOf(list,p),Object.getPrototypeOf(list)===p);console.log(Object.getPrototypeOf(Object.create(f))===f,Object.getPrototypeOf(Object.create(list))===list);try{Object.create(undefined);}catch(e){console.log(e instanceof TypeError);}`
	want, err := exec.Command(node, "-e", source).CombinedOutput()
	if err != nil {
		t.Fatalf("Node: %v\n%s", err, want)
	}
	got, err := runGoProgram(t, emitGoProgram(t, source), false)
	if err != nil {
		t.Fatalf("Go: %v\n%s", err, got)
	}
	if got != string(want) {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestGoAsyncDelegationFailuresDifferential(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("Node unavailable")
	}
	source := `function* sync(){}async function* asyncGen(){}const syncNext=Object.getPrototypeOf(sync).prototype.next,asyncNext=Object.getPrototypeOf(asyncGen).prototype.next;try{syncNext.call(asyncGen());}catch(e){console.log(e instanceof TypeError);}let closed=0;const obj={[Symbol.asyncIterator](){return {next(){return Promise.resolve({value:1,done:false});},return(){closed++;return Promise.resolve({done:true});}};}};async function* outer(){try{yield* obj;}catch(e){yield e instanceof TypeError;}}const throwing={[Symbol.asyncIterator](){return {next(){throw "next failure";}};}};async function* caught(){try{yield* throwing;}catch(e){yield e;}}async function main(){try{await asyncNext.call(sync());}catch(e){console.log(e instanceof TypeError);}const a=outer();console.log((await a.next()).value);console.log((await a.throw("x")).value,closed);console.log((await a.next()).done);console.log((await caught().next()).value);}main();`
	want, err := exec.Command(node, "-e", source).CombinedOutput()
	if err != nil {
		t.Fatalf("Node: %v\n%s", err, want)
	}
	got, err := runGoProgram(t, emitGoProgram(t, source), false)
	if err != nil {
		t.Fatalf("Go: %v\n%s", err, got)
	}
	if got != string(want) {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestGoArrayPrototypeAndFunctionDeletionDifferential(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("Node unavailable")
	}
	source := `const push=Function.prototype.call.bind(Array.prototype.push),join=Function.prototype.call.bind(Array.prototype.join),has=Function.prototype.call.bind(Object.prototype.hasOwnProperty),enumerable=Function.prototype.call.bind(Object.prototype.propertyIsEnumerable);const a=[1,2];console.log(Array.isArray(Array.prototype),Object.getPrototypeOf(a)===Array.prototype,a.join===Array.prototype.join,a.push===Array.prototype.push,Array.prototype.values===Array.prototype[Symbol.iterator]);console.log(push(a,3),join(a,"-"),has(a,"length"),enumerable(a,"length"),enumerable(a,"0"),Object.keys(Array.prototype).length);const object={0:"x",length:1};console.log(Array.prototype.push.call(object,"y"),Array.prototype.join.call(object,":"));function named(a,b){}console.log(named.name,named.length,delete named.name,delete named.length,Object.hasOwn(named,"name"),Object.hasOwn(named,"length"),named.name,named.length);`
	want, err := exec.Command(node, "-e", source).CombinedOutput()
	if err != nil {
		t.Fatalf("Node: %v\n%s", err, want)
	}
	got, err := runGoProgram(t, emitGoProgram(t, source), false)
	if err != nil {
		t.Fatalf("Go: %v\n%s", err, got)
	}
	if got != string(want) {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestGoNativeErrorStackAccessor(t *testing.T) {
	source := `const desc=Object.getOwnPropertyDescriptor(Error.prototype,"stack");console.log(desc.get.name,desc.get.length,desc.set.name,desc.set.length,desc.enumerable,desc.configurable);const e=new Error("oops");console.log(typeof e.stack,e.stack.includes("oops"),Object.hasOwn(e,"stack"));console.log(desc.get.call({})===undefined,Error.prototype.stack===undefined);desc.set.call(e,"custom");console.log(e.stack,desc.get.call(e).includes("oops"),Object.getOwnPropertyDescriptor(e,"stack").enumerable);try{desc.set.call(Error.prototype,"");}catch(e){console.log(e instanceof TypeError);}try{desc.set.call({},3);}catch(e){console.log(e instanceof TypeError);}class Child extends TypeError{}const c=new Child("child");console.log(typeof c.stack,c.stack.includes("child"));try{Math.abs(Symbol());}catch(e){console.log(typeof e.stack,e.stack.includes("TypeError"));}`
	got, err := runGoProgram(t, emitGoProgram(t, source), false)
	if err != nil {
		t.Fatalf("Go: %v\n%s", err, got)
	}
	want := "get stack 0 set stack 1 false true\nstring true false\ntrue true\ncustom true true\ntrue\ntrue\nstring true\nstring true\n"
	if got != want {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}
