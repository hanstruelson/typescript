package compiler_test

import (
	"github.com/microsoft/TypeScript/tsc/internal/core"
	"strings"
	"testing"
)

func TestGoCheckedConversions(t *testing.T) {
	source := `let input:any="9223372036854775807";let a:int64=input;let u:uint64="18446744073709551615" as any;console.log(a,u);for(const text of ["oops","","9223372036854775808","1.5"]){try{input=text;let n:int64=input;console.log("bad",n);}catch(e){console.log(e.name);}}let n:int16=-1;try{let x:uint8=n;console.log("bad",x);}catch(e){console.log(e.name);}n=255;console.log(n.uint8(),"1e3".int64(),true.int32(),u.toString(),u.string());let wide:number=1e39;try{let f:float32=wide;console.log("bad",f);}catch(e){console.log(e.name);}const precise:int64=9007199254740993;const f:number=9007199254740992;console.log(precise>f,precise==f,precise!=f);`
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	got, err := runGoProgram(t, text, true)
	want := "9223372036854775807 18446744073709551615\nRangeError\nRangeError\nRangeError\nRangeError\nRangeError\n255 1000 1 18446744073709551615 18446744073709551615\nRangeError\ntrue false true\n"
	if err != nil || got != want {
		t.Fatalf("got %q want %q: %v", got, want, err)
	}
	if !strings.Contains(text, "tsConvertUint8FromInt16(") {
		t.Fatal("known conversion missed its native helper")
	}
}
func TestGoExplicitStrictConversions(t *testing.T) {
	source := `let input:any="123";try{let n:int64=input;console.log("bad",n);}catch(e){console.log(e.name,e.message);}const n:int64=input.int64();console.log(n,n.int32());let wide:int64=300;try{let x:int8=wide;console.log("bad",x);}catch(e){console.log(e.name);}try{console.log(wide.int8());}catch(e){console.log(e.name);}try{console.log("oops".number());}catch(e){console.log(e.name);}`
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue, CoerceAny: core.TSFalse})
	got, err := runGoProgram(t, text, true)
	if err != nil || !strings.Contains(got, "Cannot implicitly convert string to int64") || !strings.HasSuffix(got, "123 123\nTypeError\nRangeError\nRangeError\n") {
		t.Fatalf("%q: %v", got, err)
	}
}

func TestGoShortConversionNames(t *testing.T) {
	source := `const value:any="42";console.log(value.number(),value.float32(),value.int8(),value.int16(),value.int32(),value.int64(),value.uint8(),value.uint16(),value.uint32(),value.uint64());const n:int64=value.int64();console.log(n.string(),n.toString(),n.boolean());try{console.log("256".uint8());}catch(e){console.log(e.name);}`
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue, CoerceAny: core.TSFalse})
	got, err := runGoProgram(t, text, true)
	want := "42 42 42 42 42 42 42 42 42 42\n42 42 true\nRangeError\n"
	if err != nil || got != want {
		t.Fatalf("got %q want %q: %v", got, want, err)
	}
}

func TestGoGrowablePrimitiveArrays(t *testing.T) {
	source := `const a:int32[]=[1,2];a.push(3);console.log(a.join(","),a.pop());a.unshift(0);console.log(a.shift(),a.join(","));a.splice(0,1,9);console.log(a.join(","));const alias:any=a;alias.push("10");console.log(a.join(","));try{alias.push("oops");}catch(e){console.log(e.name);}console.log(a.length);const b:any[]=[1,,3];b.reverse();console.log(b.join(","),1 in b);const s:string[]=["a"];s.push("b");console.log(s.pop(),s.join(","));`
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	got, err := runGoProgram(t, text, true)
	want := "1,2,3 3\n0 1,2\n9,2\n9,2,10\nRangeError\n3\n3,,1 false\nb a\n"
	if err != nil || got != want {
		t.Fatalf("got %q want %q: %v", got, want, err)
	}
}

func TestGoGrowableStrictArrays(t *testing.T) {
	source := `function identity(a:int32[]):int32[]{a.push(2);return a;}class Box{values:int32[]=[];}const a:int32[]=[1];a[0]=3;const b=identity(a);const alias:any=b;try{alias.push("4");}catch(e){console.log(e.name);}alias.push("4".int32());const box=new Box();box.values=b;box.values.push(5);console.log(a.join(","),a===b,box.values===a);const u:uint64[]=[18446744073709551615];console.log(u.pop()!.string());`
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue, CoerceAny: core.TSFalse})
	got, err := runGoProgram(t, text, true)
	want := "TypeError\n3,2,4,5 true true\n18446744073709551615\n"
	if err != nil || got != want {
		t.Fatalf("got %q want %q: %v", got, want, err)
	}
}

func TestGoLosslessWidening(t *testing.T) {
	source := `const small:int32=10;let wide:int64=small;const unsigned:uint16=65535;let signed:int32=unsigned;const f:float32=1.5;let n:number=f;let payload:any=small;wide=payload;payload=f;n=payload;const a:int64=10;console.log(a==small,small==a,a===small,small===a,f==n,n==f,f<=n,n>=f,wide,signed,n);function accept(x:int64):int64{return x;}console.log(accept(small));class Box{value:int64=small;}console.log(new Box().value);const values:int64[]=[];values.push(small);console.log(values[0]);let tiny:int8=10;let single:float32=tiny;console.log(single);const precise:int64=9007199254740993;const rounded:number=9007199254740992;payload=precise;console.log(precise>rounded,rounded<precise,payload>rounded,rounded<payload,payload==rounded);try{let rejected:number=precise;console.log("bad",rejected);}catch(e){console.log(e.name);}try{let rejected:int16=small;console.log("bad",rejected);}catch(e){console.log(e.name);}try{let rejected:uint32=small;console.log("bad",rejected);}catch(e){console.log(e.name);}try{let rejected:float32=small;console.log("bad",rejected);}catch(e){console.log(e.name);}`
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue, CoerceAny: core.TSFalse})
	got, err := runGoProgram(t, text, true)
	want := "true true true true true true true true 10 65535 1.5\n10\n10\n10\n10\ntrue true true true false\nTypeError\nTypeError\nTypeError\nTypeError\n"
	if err != nil || got != want {
		t.Fatalf("got %q want %q: %v", got, want, err)
	}
	if !strings.Contains(text, "tsConvertInt64FromInt32(") || !strings.Contains(text, "tsConvertNumberFromFloat32(") {
		t.Fatal("known widening missed native helpers")
	}
}

func TestGoScalarStringComparisonPolicy(t *testing.T) {
	source := `const n:number=10;const text:string="10";const exact:int64=9223372036854775807;let a:any="9223372036854775807";let b:any=exact;for(const test of [()=>n==text,()=>text==n,()=>n<=text,()=>text>=n,()=>a==b,()=>b==a]){try{console.log(test());}catch(e){console.log(e.name);}}let bad:any="oops";try{console.log(n==bad);}catch(e){console.log(e.name);}try{console.log(n=="oops");}catch(e){console.log(e.name);}console.log(n===text,text===n);`
	for _, test := range []struct {
		policy core.Tristate
		want   string
	}{{core.TSTrue, "true\ntrue\ntrue\ntrue\ntrue\ntrue\nRangeError\nRangeError\nfalse false\n"}, {core.TSFalse, "TypeError\nTypeError\nTypeError\nTypeError\nTypeError\nTypeError\nTypeError\nTypeError\nfalse false\n"}} {
		text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue, CoerceAny: test.policy})
		got, err := runGoProgram(t, text, true)
		if err != nil || got != test.want {
			t.Fatalf("policy %v got %q want %q: %v", test.policy, got, test.want, err)
		}
	}
}

func TestGoGrowableNativeAccess(t *testing.T) {
	source := `const a:int64[]=[];const alias:any=a;for(let i=0;i<100;i++){a.push(i.int64());}a[0]=9223372036854775807;a[1]=a[0];const index:int32=1;console.log(a[0],a[index],alias[1],a.length,alias.length);try{a.length=101;}catch(e){console.log(e.name);}try{a[-1]=7;}catch(e){console.log(e.name);}const u:uint64[]=[18446744073709551615];u.push(u[0]);console.log(u[0],u[1],u.pop());const s:string[]=["a"];s.push("b");s[0]=s[1];console.log(s.join(","));`
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue, CoerceAny: core.TSFalse})
	got, err := runGoProgram(t, text, true)
	want := "9223372036854775807 9223372036854775807 9223372036854775807 100 100\nRangeError\nRangeError\n18446744073709551615 18446744073709551615 18446744073709551615\nb,b\n"
	if err != nil || got != want {
		t.Fatalf("got %q want %q: %v", got, want, err)
	}
	for _, helper := range []string{"tsDenseRead[int64](", "tsDenseWrite[int64](", "tsDensePush[int64](", "tsDenseRead[uint64](", "tsDenseRead[*tsString]("} {
		if !strings.Contains(text, helper) {
			t.Fatal("missing direct native helper", helper)
		}
	}
}

func TestGoDenseNativeOperations(t *testing.T) {
	source := `function sum(a:int64[]):int64{let result:int64=0;for(let i:int32=0;i<a.length;i++){result+=a[i];}return result;}const x:int64=10;const y:int32=20;const a:int64[]=[];a.push(x,y);a[0]=y;a[1]+=x;console.log(sum(a),a.join(","));const b:int32[]=[1,2];a.push(...b);a.unshift(y);console.log(a.pop()!.string(),a.shift()!.string(),sum(a));const c:int64[]=[...b,...a];console.log(c.join(","));c.fill(y,0,2).reverse().copyWithin(1,0,2);console.log(c.join(","));const removed=c.splice(1,2,x,y);console.log(removed.join(","),c.join(","),c.slice(1,3).join(","));const empty:int64[]=[];console.log(empty.pop(),empty.shift());const sized=new Array<int64>(3);const strings=new Array<string>(2);console.log(sized.join(","),strings.join("|"));const anyArray:any[]=[];anyArray.push("text",x);console.log(anyArray.join(","));try{c[99]=x;}catch(e){console.log(e.name);}try{console.log(c[99]);}catch(e){console.log(e.name);}const sparse:any[]=[1,,3];try{const dense:int64[]=sparse;console.log(dense);}catch(e){console.log(e.name);}`
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue, CoerceAny: core.TSFalse})
	got, err := runGoProgram(t, text, true)
	want := "50 20,30\n2 20 51\n1,2,20,30,1\n1,1,30,20,20\n1,30 1,10,20,20,20 10,20\nundefined undefined\n0,0,0 |\ntext,10\nRangeError\nRangeError\nRangeError\n"
	if err != nil || got != want {
		t.Fatalf("got %q want %q: %v", got, want, err)
	}
	for _, helper := range []string{"tsDensePush[int64](", "tsDensePop[int64](", "tsDenseSplice[int64](", "tsDenseFill[int64](", "tsDenseArrayCell[int64]", "tsConvertInt64FromInt32("} {
		if !strings.Contains(text, helper) {
			t.Fatal("missing native operation", helper)
		}
	}
}

func TestGoDenseNativeNarrowing(t *testing.T) {
	source := `let source:int64=123;const a:int32[]=[];try{a.push(source);console.log(a[0]);}catch(e){console.log(e.name);}source=9223372036854775807;try{a.push(source);}catch(e){console.log(e.name);}let input:any="42";try{a.push(input);console.log(a[a.length-1]);}catch(e){console.log(e.name);}input="bad";try{a.push(input);}catch(e){console.log(e.name);}const exact:int32=7;a.push(exact);console.log(a.pop());`
	for _, test := range []struct {
		flag core.Tristate
		want string
	}{{core.TSTrue, "123\nRangeError\n42\nRangeError\n7\n"}, {core.TSFalse, "TypeError\nTypeError\nTypeError\nTypeError\n7\n"}} {
		text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue, CoerceAny: test.flag})
		got, err := runGoProgram(t, text, true)
		if err != nil || got != test.want {
			t.Fatalf("mode %v got %q want %q: %v", test.flag, got, test.want, err)
		}
		if !strings.Contains(text, "tsConvertInt32FromInt64(") {
			t.Fatal("narrowing missed native helper")
		}
	}
}

func TestGoDenseDynamicIndexes(t *testing.T) {
	source := `const a:int32[]=[1];let index:any=0;a[index]=2;console.log(a[index]++,++a[index],a[0]);const n:int64[]=[];const text:string="9223372036854775807";n.push(text);console.log(n[0]);const flags:boolean[]=[];const number:int32=1;flags.push(number);console.log(flags[0]);let calls=0;a.slice(0,1,++calls);console.log(calls);let size:any=2;console.log(new Array<int32>(size).join(","));size="7";console.log(new Array<int32>(size).join(","));`
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	got, err := runGoProgram(t, text, true)
	if err != nil || got != "2 4 4\n9223372036854775807\ntrue\n1\n0,0\n7\n" {
		t.Fatalf("got %q: %v", got, err)
	}
}

func TestGoNativeNullableArrays(t *testing.T) {
	source := `const a:(int64|null|undefined)[]=[1,null,undefined,4];const alias:any=a;console.log(a.length,a[0],a[1],a[2]);a.push(5,null);a.unshift(undefined,2);console.log(a.join("|"),a.pop(),a.shift());a[0]=null;alias.push("7");a.fill(null,1,3).reverse().copyWithin(1,0,2);console.log(a.join("|"));const removed=a.splice(1,2,undefined,9);console.log(removed.join("|"),a.join("|"),a.slice(1,4).join("|"));const b:(int32|null)[]=[1,null];b.push(2);console.log(b.pop(),b[1]);try{b.push(undefined);}catch(e){console.log(e.name);}const union:(number|string)[]=[1,"x"];union.push("y");console.log(union.join(","));`
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	got, err := runGoProgram(t, text, true)
	want := "4 1 null undefined\n|2|1|||4|5| null undefined\n7|7|5||||\n7|5 7||9|||| |9|\n2 null\nTypeError\n1,x,y\n"
	if err != nil || got != want {
		t.Fatalf("got %q want %q: %v", got, want, err)
	}
	if !strings.Contains(text, "tsNullablePush[int64](") || !strings.Contains(text, "tsNullableRead[int64](") {
		t.Fatal("missing native nullable operations")
	}
}
func TestGoNativeArrayCallbacks(t *testing.T) {
	source := `const a:int32[]=[3,1,2];let total:int32=0;a.forEach((x,i,arr)=>{total+=x;console.log(i,arr===a);});console.log(total,a.map(x=>x*2).join(","),a.filter(x=>x>1).join(","));console.log(a.every(x=>x>0),a.some(x=>x===2),a.find(x=>x===2),a.findLast(x=>x<3),a.findIndex(x=>x===9),a.findLastIndex(x=>x===1));console.log(a.reduce((sum,x)=>sum+x,0),a.reduceRight((sum,x)=>sum-x),a.toSorted((x,y)=>x-y).join(","),a.join(","));a.sort((x,y)=>x-y);console.log(a.join(","));const b:(number|null|undefined)[]=[1,null,undefined,4];console.log(b.map(x=>x==null?null:x*2).join("|"),b.filter(x=>x!=null).join("|"));console.log(b.reduce((sum,x)=>sum+(x??0),0));const c:number[]=[1,2,3];console.log(c.find(x=>{c[0]=10;return true;}));`
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	got, err := runGoProgram(t, text, true)
	want := "0 true\n1 true\n2 true\n6 6,2,4 3,2\ntrue true 2 2 -1 1\n6 -2 1,2,3 3,1,2\n1,2,3\n2|||8 1|4\n5\n1\n"
	if err != nil || got != want {
		t.Fatalf("got %q want %q: %v", got, want, err)
	}
	for _, helper := range []string{"tsDenseMap[", "tsNullableMap[float64,float64](", "tsDenseReduce[", "tsDenseSort[", "tsDenseFilter["} {
		if !strings.Contains(strings.ReplaceAll(text, " ", ""), helper) {
			t.Fatal("missing native callback", helper)
		}
	}
}

func TestGoNativeCallbackMutation(t *testing.T) {
	source := `function twice(x:number):number{return x*2;}const increment=(x:number)=>x+1;const a:number[]=[1,2,3];console.log(a.map(twice).join(","),a.map(increment).join(","));let calls=0;const mapped=a.map(x=>{calls++;a.push(9);return x+1;});console.log(mapped.join(","),calls,a.length);const b:number[]=[1,2,3];console.log(b.filter(x=>{b[0]=10;return true;}).join(","));const c:(number|null|undefined)[]=[3,null,undefined,1,null];c.sort((x,y)=>(x??0)-(y??0));console.log(c.join("|"),c[0],c[1],c[4]);const d:(number|null|undefined)[]=[1,null,undefined,2];const clean:number[]=d.filter(x=>x!=null);console.log(clean.join(","));const widened:(number|null)[]=a.map(x=>x);widened.push(null);console.log(widened.length,widened.pop());const empty:number[]=[];try{empty.reduce((x,y)=>x+y);}catch(e){console.log(e.name);}console.log(empty.reduce((x,y)=>x+y,10));let count=0;console.log(a.some(x=>{count++;return true;}),count);`
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	got, err := runGoProgram(t, text, true)
	want := "2,4,6 2,3,4\n2,3,4 3 6\n1,2,3\n||1|3| null null undefined\n1,2\n7 null\nTypeError\n10\ntrue 1\n"
	if err != nil || got != want {
		t.Fatalf("got %q want %q: %v", got, want, err)
	}
}

func TestGoNativeCallbackControlFlow(t *testing.T) {
	source := `function loop(x:number):number{let n=0;for(let i=0;i<x;i++){n+=i;}return n;}const a:number[]=[1,2,3];console.log(a.map(loop).join(","));let count=0;console.log(a.map(x=>{try{if(x===2)throw "two";return x*2;}catch(e){return 20;}finally{count++;}}).join(","),count);console.log(a.map(x=>{try{return x;}finally{if(x===2)return 9;}}).join(","));console.log(a.map(x=>{try{return x;}finally{try{try{return 99;}finally{throw "cancel";}}catch(e){}}}).join(","));const b:(number|null)[]=[1,null];console.log(b.map(x=>{if(x===null)return null;let result=0;for(let i=0;i<x;i++){result++;}return result;}).join("|"));`
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	got, err := runGoProgram(t, text, true)
	want := "0,1,3\n2,20,6 3\n1,9,3\n1,2,3\n1|\n"
	if err != nil || got != want {
		t.Fatalf("got %q want %q: %v", got, want, err)
	}
}

func TestGoNativeCallbackFlatMapAndDefaults(t *testing.T) {
	source := `const a:number[]=[1,2,3];console.log(a.flatMap(x=>[x,x+10]).join(","));const twice=(x:number)=>[x,x];console.log(a.flatMap(twice).join(","));console.log(a.flatMap(x=>{if(x===2)return [];return [x];}).join(","));const b:(number|undefined)[]=[1,undefined,3];console.log(b.map((x=5)=>x*2).join(","));console.log(b.map((x=5)=>{try{return x*2;}finally{}}).join(","));function local(){const add=(x:number)=>x+2;return a.map(add).join(",");}console.log(local());const c:(int64|null)[]=[9223372036854775807,null];const d:(int64|null)[]=[...c];d.push(...c);console.log(d.join("|"));const sized=new Array<string|null>(2);console.log(sized.length,sized.join("|"));`
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	got, err := runGoProgram(t, text, true)
	want := "1,11,2,12,3,13\n1,1,2,2,3,3\n1,3\n2,10,6\n2,10,6\n3,4,5\n9223372036854775807||9223372036854775807|\n2 |\n"
	if err != nil || got != want {
		t.Fatalf("got %q want %q: %v", got, want, err)
	}
	if !strings.Contains(text, "tsDenseFlatMap[") {
		t.Fatal("flatMap did not select native storage")
	}
}

func TestGoNativeNamedNullableCallbacks(t *testing.T) {
	source := `function present(x:number|null):boolean{return x!==null;}function pair(x:number):(number|null)[]{return [x,null];}const a:(number|null)[]=[1,null,3];console.log(a.filter(present).join("|"));const b:number[]=[1,2];console.log(b.flatMap(pair).join("|"));let total=0;function add(x:number):void{for(let i=0;i<x;i++){total++;}return;}b.forEach(add);b.forEach(x=>{for(let i=0;i<x;i++){total++;}});console.log(total);function sum(values:(number|null)[]):number{let result=0;for(let i=0;i<values.length;i++){result+=values[i]??0;}return result;}console.log(sum(a));`
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	got, err := runGoProgram(t, text, true)
	if err != nil || got != "1|3\n1||2|\n6\n4\n" {
		t.Fatalf("got %q: %v", got, err)
	}
	for _, helper := range []string{"tsDenseFilter[", "tsDenseFlatMap[", "tsDenseForEach["} {
		if !strings.Contains(text, helper) {
			t.Fatal("missing native callback", helper)
		}
	}
}
