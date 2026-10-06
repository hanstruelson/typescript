package compiler_test

import (
	"github.com/microsoft/TypeScript/tsc/internal/core"
	"testing"
)

func TestGoSparseArrays(t *testing.T) {
	source := `const a:any[]=[,undefined,,4];console.log(a.length,0 in a,1 in a,2 in a,Object.hasOwn(a,"0"),Object.hasOwn(a,"1"));console.log(Object.keys(a).length);for(const x of a)console.log(x);const b=[...a];console.log(0 in b,2 in b,Object.keys(b).length);a[0]=7;a[6]=9;console.log(a.length,4 in a,5 in a,6 in a,Object.keys(a).length);a.length=2;a.length=7;console.log(6 in a,a[6],Object.keys(a).length);a.push(8);console.log(7 in a,a[7]);a.extra=3;console.log(Object.hasOwn(a,"extra"),Object.hasOwn(a,"01"));a[4294967295]=6;console.log(a.length,a[4294967295],Object.hasOwn(a,"4294967295"));const c=[1,,...[2,3],,5];console.log(c.length,1 in c,4 in c,c[5]);let keys="";for(const key in c)keys+=key;console.log(keys);async function read(){const sparse=[,await Promise.resolve(7),,];console.log(sparse.length,0 in sparse,1 in sparse,2 in sparse,sparse[1]);}read();`
	want := javascriptOutput(t, source)
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	got, err := runGoProgram(t, text, true)
	if err != nil || got != want {
		t.Fatalf("got %q want %q: %v", got, want, err)
	}
}
