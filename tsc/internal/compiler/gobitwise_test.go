package compiler_test

import (
	"github.com/microsoft/TypeScript/tsc/internal/core"
	"strings"
	"testing"
)

func TestGoNumberBitwise(t *testing.T) {
	cases := []struct{ name, source string }{
		{"boundaries", `const values=[0,-0,1,-1,1.9,-1.9,2147483647,2147483648,4294967295,4294967296,4294967297,-4294967297,9007199254740991,NaN,Infinity,-Infinity];for(const x of values){console.log(x|0,x&255,x^123,~x);for(const n of [-33,-32,-1,0,1,31,32,33,65,1.9,NaN,Infinity])console.log(x<<n,x>>n,x>>>n);}`},
		{"coercion", `const values:any[]=[null,undefined,true,false,"4294967297","-3.9","bad",[],[3],[1,2],{}, {valueOf:()=>7},{valueOf:()=>({}),toString:()=>"9"}];for(const x of values)console.log(x|0,x<<1,x>>1,x>>>1,~x);`},
		{"compounds", `let log="";function left(){log+="L";return 4294967295;}function right(){log+="R";return 33;}console.log(left()<<right(),log);let n=4294967295;n>>>=0;console.log(n);n<<=33;console.log(n);n>>=1;n|=3;n&=255;n^=17;console.log(n);const obj:any={n:-1};obj.n>>>=1;console.log(obj.n);`},
		{"coercion order", `let order="";const a:any={valueOf:()=>{order+="a";return -1;}};const b:any={valueOf:()=>{order+="b";return 33;}};console.log(a>>>b,order);`},
		{"exponentiation", `let x:number=3;x**=4;console.log(x,2**-3,(-1)**Infinity,(-0)**-3,NaN**0);const base:any={valueOf:()=>2};console.log(base**3);`},
		{"async", `function bits<T extends number>(x:T){return (x<<33)^(x>>>31);}async function main(){let x:number=await Promise.resolve(-1);x>>>=1;console.log(x,~x,bits<number>(x));}main();`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			want := javascriptOutput(t, tc.source)
			text := emitGoOptions(t, tc.source, core.CompilerOptions{Strict: core.TSTrue})
			got, err := runGoProgram(t, text, true)
			if err != nil || got != want {
				t.Fatalf("got %q, want %q: %v", got, want, err)
			}
		})
	}
}
func TestGoNumberBitwiseNativeEmission(t *testing.T) {
	text := emitGoOptions(t, `function shift(x:number,n:number):number{return x>>>n;}let x:number=-1;x<<=33;console.log(shift(x,1),~x);`, core.CompilerOptions{Strict: core.TSTrue})
	for _, fragment := range []string{`tsNumberBitwise(">>>",`, `tsNumberBitwise("<<",`, `tsNumberBitwiseNot(`} {
		if !strings.Contains(text, fragment) {
			t.Fatalf("missing native call %s", fragment)
		}
	}
	got, err := runGoProgram(t, text, true)
	if err != nil || got != "2147483647 1\n" {
		t.Fatalf("%q %v", got, err)
	}
}
