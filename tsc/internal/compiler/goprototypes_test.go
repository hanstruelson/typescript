package compiler_test

import (
	"os/exec"
	"strings"
	"testing"
)

func TestGoPrototypeDifferential(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("Node unavailable")
	}
	source := `function F(x){this.x=x;} const a=new F(4);const b=new F(8);
 F.prototype.read=function(){return this.x;};console.log(a.read(),b.read(),a instanceof F,Object.getPrototypeOf(a)===F.prototype);
 const alias=F.prototype;alias.extra=7;console.log(a.extra,b.extra);a.extra=9;console.log(a.extra,b.extra);
 const old=F.prototype;F.prototype={read:function(){return this.x+1;}};const c=new F(2);console.log(a.read(),c.read(),a instanceof F,c instanceof F,Object.getPrototypeOf(a)===old);
 const Bound=F.bind(null,5);const bound=new Bound();console.log(bound.read(),bound instanceof Bound,bound instanceof F,Bound.prototype===undefined);
 function ObjectResult(){return {answer:42};}function PrimitiveResult(){this.answer=3;return 9;}console.log(new ObjectResult().answer,new PrimitiveResult().answer);
 console.log(F.call===F.call,F.apply===F.apply,F.bind===F.bind,Object.getPrototypeOf(F)===Function.prototype,Object.getPrototypeOf(F.prototype)===Object.prototype);
 const protoSource={get doubled(){return this.value*2;},set doubled(v){this.value=v/2;}};const child=Object.create(protoSource);child.value=3;console.log(child.doubled);child.doubled=10;console.log(child.value,protoSource.value===undefined);
 const proto={x:11};Object.setPrototypeOf(child,proto);console.log(child.x,Object.getPrototypeOf(child)===proto);try{Object.setPrototypeOf(proto,child);}catch(e){console.log(e.name);}
 console.log("read" in a,"toString" in a,"prototype" in F,"call" in F);console.log(Object.keys(F).length,Object.getOwnPropertyDescriptor(F,"prototype").enumerable);let inheritedKeys="";for(const key in c){inheritedKeys+=key+",";}console.log(inheritedKeys);const legacy={};legacy.__proto__=proto;console.log(legacy.x,legacy.__proto__===proto);
 const literal={__proto__:proto};const computed={["__proto__"]:proto};console.log(literal.x,computed.x===undefined,Object.hasOwn(computed,"__proto__"));const nil=Object.create(null);console.log(Object.getPrototypeOf(nil)===null,nil.toString===undefined);
 console.log(Object.prototype.isPrototypeOf.call(F.prototype,c),Object.prototype.hasOwnProperty.call(a,"x"),Object.prototype.hasOwnProperty.call(a,"read"));
 const arrow=()=>1;console.log(arrow.prototype===undefined);try{new arrow();}catch(e){console.log(e.name);}
 class Base {x=2;read(){return this.x;}}class Child extends Base {read(){return super.read()+1;}}
 class Mutable {read(){return 1;}replace(){this.read=function(){return 7;};}}const mutable=new Mutable();mutable.replace();console.log(mutable.read());
 const other=new Child();const item=new Child();console.log(item.read===other.read);console.log(item.read(),Object.getPrototypeOf(item)===Child.prototype,Object.getPrototypeOf(Child.prototype)===Base.prototype,item instanceof Base);
 Base.prototype.read=function(){return this.x+10;};console.log(item.read(),Base.prototype.read.call({x:20}));Child.prototype.read=function(){return this.x+30;};console.log(item.read());`
	want, err := exec.Command(node, "-e", source).CombinedOutput()
	if err != nil {
		t.Fatalf("Node: %v\n%s", err, want)
	}
	text := emitGoProgram(t, source)
	got, err := runGoProgram(t, text, false)
	if err != nil {
		t.Fatalf("Go: %v\n%s", err, got)
	}
	if got != string(want) {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}
func TestGoPrototypeFreeClassUsesDirectAccess(t *testing.T) {
	text := emitGoProgram(t, `class Point{x:number=3;read(){return this.x;}} const p=new Point();console.log(p.read());`)
	if strings.Contains(text, "prototypeFactory = func") {
		t.Fatal("Closed class allocated prototype factory")
	}
	if !strings.Contains(text, ".CallP72656164(loop") {
		t.Fatal("Closed class lost direct method call")
	}
	got, err := runGoProgram(t, text, false)
	if err != nil || got != "3\n" {
		t.Fatalf("%v %q", err, got)
	}
}
func TestGoPrototypeSharedInitializationRace(t *testing.T) {
	source := `function F(){this.x=1;}function inspect(){return Object.getOwnPropertyDescriptor(F,"prototype").value;}async function main(){const values=await Promise.all([go inspect(),go inspect()]);console.log(values[0]===values[1]);}main();`
	got, err := runGoProgram(t, emitGoProgram(t, source), true)
	if err != nil || got != "true\n" {
		t.Fatalf("%v %s", err, got)
	}
}
