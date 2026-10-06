package compiler_test

import (
	"strings"
	"testing"

	"github.com/microsoft/TypeScript/tsc/internal/core"
)

func TestGoTypeSystemCoverage(t *testing.T) {
	cases := []struct{ name, source, want string }{
		{"primitive unions", `function describe(value:number|string){if(typeof value==="number"){return value+1;}return value.length;}console.log(describe(4),describe("hello"));`, "5 5\n"},
		{"discriminated unions", `type Shape={kind:"square";size:number}|{kind:"circle";radius:number};function area(s:Shape){switch(s.kind){case "square":return s.size*s.size;case "circle":return s.radius*s.radius*3;}}console.log(area({kind:"square",size:4}),area({kind:"circle",radius:2}));`, "16 12\n"},
		{"flat intersections", `interface Named{name:string;}interface Count{count:number;}type Combined=Named&Count;function bump<T extends Combined>(x:T):T{x.count++;return x;}const original:Combined={name:"same",count:1};const result=bump(original);console.log(result===original,result.name,original.count,Object.keys(result).length);`, "true same 2 2\n"},
		{"generic functions and classes", `function identity<T>(x:T):T{return x;}function pluck<T,K extends keyof T>(x:T,key:K):T[K]{return x[key];}class Box<T>{value:T;constructor(value:T){this.value=value;}getValue():T{return this.value;}}const box=new Box<string>(identity("hello"));console.log(box.getValue(),identity(7),pluck({name:"world"},"name"));`, "hello 7 world\n"},
		{"mapped conditional and indexed types", `type Read<T>={readonly [K in keyof T]:T[K]};type Choice<T>=T extends number?{number:T}:{text:T};type Keys<T>=keyof T;function take<T,K extends Keys<T>>(x:Read<T>,key:K):T[K]{return x[key];}const record:Read<{name:string;count:number}>={name:"ok",count:3};const choice:Choice<number>={number:8};console.log(take(record,"name"),choice.number);`, "ok 8\n"},
		{"generic alias scope", `type T=number;type Alias<T>=T;type Maybe<T=string>=T|null;const text:Alias<string>="correct";const number:Alias<number>=7;function read(x:Maybe){if(x===null)return "none";return x;}console.log(text,number,read("ok"),read(null));`, "correct 7 ok none\n"},
		{"satisfies and const assertions", `const value={kind:"ready",count:2} as const satisfies {kind:string;count:number};console.log(value.kind,value.count);`, "ready 2\n"},
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
			if strings.Contains(text, "interface{}") {
				t.Fatal("type erasure used a Go interface")
			}
		})
	}
}
