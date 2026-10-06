package compiler_test

import (
	"regexp"
	"testing"

	"github.com/microsoft/TypeScript/tsc/internal/core"
)

func TestGoClassParameters(t *testing.T) {
	cases := []struct{ name, source string }{
		{"parameter properties", `class Item{constructor(public count:number=3,readonly name:string="item",private active:boolean=true){}read(){return this.active;}}const x=new Item();console.log(x.count,x.name,x.read(),Object.hasOwn(x,"count"));x.count++;console.log(x.count);`},
		{"parameter property order", `class Base{seen:any=this.x;present=Object.hasOwn(this,"x");constructor(public x=3){const keys=Object.keys(this);console.log(this.seen,this.present,this.x,keys[0],keys[1]);}}class Child extends Base{seenChild:any=this.x;constructor(public x=4){super(9);console.log(this.seenChild,this.x);}}new Base();new Child();`},
		{"destructured constructor and methods", `class DataRecord{value:number;constructor({value=3}:any={}){this.value=value;}read({a=2,...rest}:any,[first,...tail]:number[]){return this.value+a+rest.b+first+tail.length;}}const x=new DataRecord();console.log(x.read({b:4},[5,6,7]));`},
		{"constructor rest and spreads", `class Base{first:any;count:number;constructor(...values:any[]){this.first=values[0];this.count=values.length;}}class Child extends Base{constructor(...values:any[]){super(...values);}}const values=[7,8,9];const x=new Child(...values);console.log(x.first,x.count);`},
		{"method rest and async patterns", `class Reader{sum(start:number,...values:number[]){let n=start;for(const x of values)n+=x;return n;}async asyncSum({start=2}:any,...values:number[]){let n=start;for(const x of values)n+=await Promise.resolve(x);return n;}}const r=new Reader();console.log(r.sum(1,2,3));r.asyncSum({},3,4).then(value=>console.log(value));`},
		{"parameter default order", `class Values{constructor(public a:number=2,public b:number=a+3){}read(a:number=1,b:number=a+2){return a+b;}}const x=new Values();console.log(x.a,x.b,x.read());`},
		{"parameter temporal dead zones", `function selfDefault(x:number=x){return x;}class Reader{self(x:number=x){return x;}later(a:number=b,b:number=3){return a+b;}}const r=new Reader();try{selfDefault();}catch(e){console.log("function TDZ");}try{r.self();}catch(e){console.log("method TDZ");}try{r.later();}catch(e){console.log("later TDZ");}`},
		{"nullable and optional properties", `class NullableBox{constructor(public value:number|null=null,public label?:string){}}const x=new NullableBox();console.log(x.value===null,x.label===undefined);x.value=3;console.log(x.value);`},
		{"reference sharing", `class Box{constructor(public value:any){}}const original={count:1};const x=new Box(original);x.value.count++;console.log(x.value===original,original.count);`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			want := javascriptOutput(t, tc.source)
			text := emitGoOptions(t, tc.source, core.CompilerOptions{Strict: core.TSTrue})
			got, err := runGoProgram(t, text, true)
			if err != nil || got != want {
				t.Fatalf("got %q, want %q, error %v", got, want, err)
			}
		})
	}
}
func TestGoNativeParameterProperties(t *testing.T) {
	source := `class Base{constructor(public value:int64){}}class Child extends Base{constructor(public other:int64=9007199254740993){super(9223372036854775807);}}const x=new Child();console.log(x.value,x.other);class Direct{constructor(public value:int64){}}console.log(new Direct(9223372036854775807).value);class DestructNative{read({value=9007199254740993}:{value?:int64}){return value;}}console.log(new DestructNative().read({}));const AliasDirect=Direct;console.log(new AliasDirect(9223372036854775807).value);class OptionalNative{constructor(public value?:int64){}}console.log(new OptionalNative(9223372036854775807).value,new OptionalNative().value===undefined);`
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	if !regexp.MustCompile(`P76616c7565\s+int64`).MatchString(text) {
		t.Fatal("parameter property did not use native int64 storage")
	}
	got, err := runGoProgram(t, text, true)
	if want := "9223372036854775807 9007199254740993\n9223372036854775807\n9007199254740993\n9223372036854775807\n9223372036854775807 true\n"; err != nil || got != want {
		t.Fatalf("got %q, want %q, error %v", got, want, err)
	}
}

func TestGoNativeRestArguments(t *testing.T) {
	source := `function first(...values:int64[]){return values[0];}class NativeRest{value:int64;constructor(...values:int64[]){this.value=values[0];}read(...values:int64[]){return values[1];}}const x=new NativeRest(9223372036854775807);console.log(x.value,x.read(1,9007199254740993),first(9223372036854775807));`
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	got, err := runGoProgram(t, text, true)
	if want := "9223372036854775807 9007199254740993 9223372036854775807\n"; err != nil || got != want {
		t.Fatalf("got %q, want %q, error %v", got, want, err)
	}
}
