package compiler_test

import (
	"testing"

	"github.com/microsoft/TypeScript/tsc/internal/core"
)

func TestGoAdditionalSyntax(t *testing.T) {
	cases := []struct{ name, source, want string }{
		{"enums", `enum Mode{First,Next,Explicit=7,Last}enum TextEnum{A="a",B="b"}enum Derived{A=2,B=A+3}const enum Fixed{A=9}enum Merged{A=1}enum Merged{B=2}console.log(Mode.First,Mode.Next,Mode.Last,Mode[7],TextEnum.B,Derived.B,Fixed.A,Merged.A,Merged.B);`, "0 1 8 Explicit b 5 9 1 2\n"},
		{"optional chains", `let calls=0;function effect(){calls++;return 3;}const missing:any=null;const value:any={child:{read:(x:number)=>x+1}};console.log(missing?.child.read(effect()),value?.child.read(effect()),calls);console.log(missing?.[effect()],value.child.read?.(effect()),calls);try{console.log((missing?.child).read);}catch(e){console.log("group threw");}`, "undefined 4 1\nundefined 4 2\ngroup threw\n"},
		{"optional chain await", `let count=0;function effect(){count++;return 4;}async function read(x:any){return x?.run(await Promise.resolve(effect()));}async function main(){console.log(await read(null),count);console.log(await read({run:(n:number)=>n+1}),count);}main();`, "undefined 0\n5 1\n"},
		{"direct optional chain", `class Reader{read(value:any):number{return value?.child?.number??0;}}const r=new Reader();console.log(r.read(null),r.read({child:{number:7}}));`, "0 7\n"},
		{"named function expression", `let factorial=function fact(n:number):number{if(n<2)return 1;return n*fact(n-1);};const saved=factorial;factorial=(n:number)=>-1;console.log(saved(5),factorial(3));`, "120 -1\n"},
		{"logical assignments", `let count=0;function effect(){count++;return 7;}let a:any=0;let b:any=2;let c:any=null;a||=effect();b||=effect();c??=effect();a&&=effect();console.log(a,b,c,count);const target:any={key:0};function receiver(){count++;return target;}receiver().key ||= effect();console.log(target.key,count);`, "7 2 7 3\n7 5\n"},
		{"operators", `let n=5;n/=2;n%=2;let b=3;b**=3;b|=4;b&=31;b^=2;b<<=1;b>>=1;b>>>=1;console.log(n,b,2**3,~1,1<<33,-1>>>0);let x:any=0;let y:any="";console.log(x==y,null==undefined);function read(v:{text:string}|{size:number}){if("text" in v)return v.text;return v.size;}console.log(read({text:"yes"}),read({size:8}),"length" in []);`, "0.5 14 8 -2 2 4294967295\ntrue true\nyes 8 true\n"},
		{"number formatting", `console.log(4294967295,1e20,1e21,1e-6,1e-7,Infinity,-Infinity,NaN,-0);`, "4294967295 100000000000000000000 1e+21 0.000001 1e-7 Infinity -Infinity NaN 0\n"},
		{"native bitwise", `let n:int64=9007199254740993;n|=2;console.log(n,~n);n<<=1;console.log(n);`, "9007199254740995 -9007199254740996\n18014398509481990\n"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			text := emitGoOptions(t, tc.source, core.CompilerOptions{Strict: core.TSTrue})
			got, err := runGoProgram(t, text, true)
			if err != nil {
				t.Fatalf("%v\n%s", err, got)
			}
			if got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}
