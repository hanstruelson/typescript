package compiler_test

import (
	"github.com/microsoft/TypeScript/tsc/internal/core"
	"strings"
	"testing"
)

func TestGoCollections(t *testing.T) {
	cases := []struct{ name, source string }{
		{"map keys", `const key={};const m=new Map<any,any>([["a",1],[key,2],[NaN,3],[-0,4]]);console.log(m.size,m.get(key),m.get({}),m.get(NaN),m.get(0),m.set("a",5)===m);for(const [k,v] of m)console.log(typeof k,v);console.log(m.has("a"),m.delete("a"),m.delete("a"),m.size);m.clear();console.log(m.size);`},
		{"set", `const key={};const s=new Set<any>([1,1,NaN,NaN,-0,0,key]);console.log(s.size,s.has(NaN),s.has(key),s.has({}),s.add(7)===s);for(const x of s.values())console.log(typeof x);for(const pair of s.entries())console.log(pair[0]===pair[1]);console.log(s.delete(1),s.delete(1));s.clear();console.log(s.size);`},
		{"live iterators", `const m=new Map([[1,"a"],[2,"b"]]);const it=m.keys();console.log(it.next().value);m.delete(2);m.set(3,"c");console.log(it.next().value,it.next().done);m.set(4,"d");console.log(it.next().done);const s=new Set([1,2]);s.forEach((v,k)=>{console.log(v,k);if(v===1){s.delete(2);s.add(3);}});m.forEach((v,k,own)=>console.log(v,k,own===m));`},
		{"clear during iteration", `const s=new Set([1,2]);const it=s.values();console.log(it.next().value);s.clear();s.add(3);console.log(it.next().value,it.next().done);const m=new Map([["a",1],["b",2]]);m.delete("a");m.set("a",3);console.log([...m.keys()].join(","));`},
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
func TestGoArrayMethods(t *testing.T) {
	source := `const a:any[]=[1,,3,undefined,NaN];console.log(a.includes(undefined),a.indexOf(undefined),a.includes(NaN),a.indexOf(NaN));console.log(a.map((x,i)=>x===undefined?i:x*2).join("|"),a.filter(x=>x!==undefined).join("|"));console.log(a.slice(1,4).length,0 in a.slice(1,4));console.log(a.reduce((sum,x)=>sum+(x||0),0),a.some(x=>x===3),a.every(x=>x!==2),a.findIndex(x=>x===undefined));console.log([...a.keys()].join(","),[...a.values()].join(","));const b=[1,2,3];console.log(b.pop(),b.shift(),b.unshift(7,8),b.join(","));console.log(b.splice(1,1,9,10).join(","),b.join(","));console.log(b.reverse().join(","),b.fill(4,1,-1).join(","));console.log([10,2,1].sort().join(","),[10,2,1].sort((a,b)=>a-b).join(","));const sparse:any[]=[1,,3];const tail:any[]=[,5];console.log(sparse.concat(tail,7).join("|"));let count=0;sparse.forEach(x=>count+=x);console.log(count,[1,2,3].reduceRight((a,b)=>a-b),[1,2,3].at(-1));`
	want := javascriptOutput(t, source)
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	got, err := runGoProgram(t, text, true)
	if err != nil || got != want {
		t.Fatalf("got %q, want %q: %v", got, want, err)
	}
}

func TestGoGenericArrayPrototypeReceivers(t *testing.T) {
	source := `const a:any={0:3,2:1,length:3};console.log(Array.prototype.map.call(a,(x:any)=>x===undefined?"u":x*2).join(","));console.log(Array.prototype.slice.call(a,0,3).join("|"));Array.prototype.reverse.call(a);console.log(a[0],a[1],a[2],a.length);Array.prototype.fill.call(a,9,1,3);console.log(a[0],a[1],a[2]);const b:any={0:"a",2:"c",length:3};console.log(Array.prototype.includes.call(b,undefined),Array.prototype.indexOf.call(b,"c"));`
	want := javascriptOutput(t, source)
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	got, err := runGoProgram(t, text, true)
	if err != nil || got != want {
		t.Fatalf("got %q, want %q: %v", got, want, err)
	}
}

func TestGoArrayPropertyDescriptors(t *testing.T) {
	source := `const a:any[]=[];Object.defineProperty(a,"0",{value:7,writable:false,enumerable:false,configurable:false});console.log(a.length,a[0],Object.keys(a).join(","),Object.getOwnPropertyDescriptor(a,"0").writable);try{a[0]=8;}catch(e){console.log(e.name)};Object.defineProperty(a,"1",{get(){return 9},enumerable:true,configurable:true});console.log(a.length,a[1],Object.keys(a).join(","),Object.getOwnPropertyDescriptor(a,"1").get!==undefined);delete a[1];console.log(a.length,1 in a);`
	want := javascriptOutput(t, source)
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	got, err := runGoProgram(t, text, true)
	if err != nil || got != want {
		t.Fatalf("got %q, want %q: %v", got, want, err)
	}
}

func TestGoSparseArrayReflection(t *testing.T) {
	source := `const a:any[]=[];a[3]=4;a.extra=9;Object.defineProperty(a,"1",{value:2,enumerable:false,configurable:true});console.log(a.length,Object.keys(a).join(","),Object.getOwnPropertyNames(a).join(","),Reflect.ownKeys(a).join(","));console.log(0 in a,1 in a,2 in a,3 in a,a.extra);delete a[1];console.log(a.length,1 in a,Object.getOwnPropertyNames(a).join(","));a.length=2;console.log(a.length,3 in a,Object.keys(a).join(","));`
	want := javascriptOutput(t, source)
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	got, err := runGoProgram(t, text, true)
	if err != nil || got != want {
		t.Fatalf("got %q, want %q: %v", got, want, err)
	}
}
func TestGoTypedArrays(t *testing.T) {
	source := `const a=new Float64Array([1,2,3]);a[1]=4.5;const view=a.subarray(1);view[0]=7.5;const copy=a.slice();copy[0]=99;console.log(a.length,a.byteLength,a.BYTES_PER_ELEMENT,a.join(","),copy[0]);a.set(a.subarray(0,2),1);console.log(a.join(","));console.log(a.map(x=>x*2).join(","),a.filter(x=>x>1).join(","),a.reduce((sum,x)=>sum+x,0));console.log(new Int8Array([128,-129,257,NaN]).join(","));console.log(new Uint8ClampedArray([-.5,.5,1.5,2.5,255.5,NaN]).join(","));console.log(new Int16Array([32768,-32769]).join(","),new Uint16Array([-1,65537]).join(","));console.log(new Int32Array([2147483648,4294967297]).join(","),new Uint32Array([-1,4294967297]).join(","));const f=new Float32Array([1/3]);console.log(f[0]);console.log([...new Uint8Array([3,1,2]).sort()].join(","));const z=new Float64Array(2);z[9]=7;console.log(z[9],z.length,z[0]);`
	want := javascriptOutput(t, source)
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	got, err := runGoProgram(t, text, true)
	if err != nil || got != want {
		t.Fatalf("got %q, want %q: %v", got, want, err)
	}
}
func TestGoArrayConstructors(t *testing.T) {
	source := `const a=new Array<any>(3);console.log(a.length,Object.keys(a).length,0 in a);const b=new Array(1,2,3);console.log(Array.isArray(b),Array.isArray(new Float64Array(2)),Array.of(4,5).join(","));console.log(Array.from(new Set([3,1,3]),(x,i)=>x+i).join(","));const M=Map;const m=new M([[1,2]]);console.log(m.get(1),m instanceof Map,Map===Map,m instanceof Set);const F=Float32Array;const f=new F([2,4]);console.log(f instanceof Float32Array,f instanceof Float64Array,Float32Array.BYTES_PER_ELEMENT,Float32Array.name,Object.keys(f).join(","),Object.hasOwn(f,"length"));const u=new Uint8Array(1);u["-0"]=7;u["1.5"]=9;console.log(u[0],u["-0"],u["1.5"],Object.keys(u).join(","));`
	want := javascriptOutput(t, source)
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	got, err := runGoProgram(t, text, true)
	if err != nil || got != want {
		t.Fatalf("got %q, want %q: %v", got, want, err)
	}
}
func TestGoArrayBufferViews(t *testing.T) {
	source := `const buffer=new ArrayBuffer(16);const a=new Uint32Array(buffer);const b=new Uint8Array(buffer);a[0]=258;console.log(buffer.byteLength,a.length,b.length,a.buffer===buffer,b.buffer===buffer,b[0],b[1]);const tail=new Uint32Array(buffer,4,2);tail[0]=7;console.log(a[1],tail.byteOffset,tail.byteLength,tail.subarray(1).byteOffset,tail.subarray(1).buffer===buffer);const copy=buffer.slice(0,4);const copied=new Uint32Array(copy);copied[0]=9;console.log(a[0],copied[0],copy===buffer,buffer instanceof ArrayBuffer);try{new Uint32Array(buffer,1);}catch(e){console.log("alignment");}try{new Float64Array(buffer,8,2);}catch(e){console.log("bounds");}const empty=new Uint8Array(buffer,16);console.log(empty.length,empty.byteOffset);`
	want := javascriptOutput(t, source)
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	got, err := runGoProgram(t, text, true)
	if err != nil || got != want {
		t.Fatalf("got %q, want %q: %v", got, want, err)
	}
}

func TestGoDataViewAndResizableArrayBuffer(t *testing.T) {
	source := `const buffer=new ArrayBuffer(16,{maxByteLength:32});const view=new DataView(buffer,2,12);view.setUint16(0,258);view.setInt32(2,-7,true);view.setFloat64(4,1.5,true);view.setBigUint64(0,123n,true);console.log(buffer.resizable,buffer.maxByteLength,buffer.byteLength,view.byteOffset,view.byteLength,view.getBigUint64(0,true),view.getUint16(0,true),view.getInt32(2,true),view.getFloat64(4,true),view instanceof DataView,ArrayBuffer.isView(view),Object.prototype.toString.call(view));buffer.resize(24);console.log(buffer.byteLength,view.byteLength,view.getUint8(0));buffer.resize(4);console.log(buffer.byteLength);try{view.getUint8(0)}catch(e){console.log(e.name)};try{buffer.resize(33)}catch(e){console.log(e.name)};const fixed=buffer.transferToFixedLength();console.log(buffer.byteLength,fixed.byteLength,fixed.resizable,fixed.maxByteLength);const detached=new DataView(new ArrayBuffer(2));detached.buffer.transfer();try{detached.getUint8(0)}catch(e){console.log(e.name)}`
	want := javascriptOutput(t, source)
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	got, err := runGoProgram(t, text, true)
	if err != nil || got != want {
		t.Fatalf("got %q, want %q: %v", got, want, err)
	}
}
func TestGoSetAlgebra(t *testing.T) {
	source := `const a=new Set([3,1,2]);const b=new Set([2,4]);console.log([...a.union(b)].join(","),[...a.intersection(b)].join(","),[...a.difference(b)].join(","),[...a.symmetricDifference(b)].join(","));console.log(a.isSubsetOf(b),a.isSupersetOf(new Set([1,2])),a.isDisjointFrom(new Set([7])));const m=new Map([[2,"x"],[5,"y"]]);console.log([...a.union(m)].join(","),[...a.intersection(m)].join(","));`
	want := javascriptOutput(t, source)
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	got, err := runGoProgram(t, text, true)
	if err != nil || got != want {
		t.Fatalf("got %q, want %q: %v", got, want, err)
	}
}
func TestGoArrayCopyMethods(t *testing.T) {
	source := `const a:any[]=[1,,3,4];a.copyWithin(1,0,3);console.log(a.join("|"),1 in a,2 in a,a.toReversed().join("|"),2 in a.toReversed());const x=new Float64Array([1,2,3,4]);x.copyWithin(1,0,3);console.log(x.join(","),x.toReversed().join(","),x.toSorted().join(","),x.join(","));console.log(Uint8Array.of(256,-1).join(","),Float64Array.from(new Set([2,3]),x=>x*2).join(","),ArrayBuffer.isView(x),ArrayBuffer.isView(x.buffer));`
	want := javascriptOutput(t, source)
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	got, err := runGoProgram(t, text, true)
	if err != nil || got != want {
		t.Fatalf("got %q, want %q: %v", got, want, err)
	}
}

func TestGoTypedArrayNativeAccess(t *testing.T) {
	source := `function first(a:Float64Array){return a[0];}const a=new Float64Array([1,2]);a[1]=7.5;console.log(a[1],first(a),first(new Uint8Array([255]) as unknown as Float64Array));async function worker(x:Float64Array){x[0]=9;return x[0];}async function main(){console.log(await launch worker(a),a[0]);}main();`
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	for _, fragment := range []string{"tsNativeArrayRead[float64](", "tsNativeArrayWrite[float64]("} {
		if !strings.Contains(text, fragment) {
			t.Fatalf("missing native access %s", fragment)
		}
	}
	got, err := runGoProgram(t, text, true)
	if err != nil || got != "7.5 1 255\n9 9\n" {
		t.Fatalf("%q: %v", got, err)
	}
}

func TestGoBigIntTypedArrays(t *testing.T) {
	source := `const a=new BigInt64Array([1n,-2n,18446744073709551615n]);const b=new BigUint64Array([1n,-1n]);a[1]=-3n;b[0]=2n;console.log(a.join(","),b.join(","),typeof a[0],a.BYTES_PER_ELEMENT,a.buffer.byteLength);const view=new BigInt64Array(a.buffer,8,2);console.log(view[0].toString(),view[1].toString());try{a[0]=1;}catch(e){console.log(e.name)};console.log(a.map(x=>x+1n).join(","));const copy=new BigInt64Array(a);copy.sort();copy.fill(-4n,0,1);console.log(copy.join(","));`
	want := javascriptOutput(t, source)
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	got, err := runGoProgram(t, text, true)
	if err != nil || got != want {
		t.Fatalf("got %q, want %q: %v", got, want, err)
	}
}

func TestGoFloat16TypedArray(t *testing.T) {
	source := `const a=new Float16Array([1.5,1/3,65504,Infinity]);const view=new Float16Array(a.buffer);console.log(a.BYTES_PER_ELEMENT,a.byteLength,a.join(","),view[0],view[1]);a[0]=2.25;console.log(view[0]);`
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	got, err := runGoProgram(t, text, true)
	if err != nil || got != "2 8 1.5,0.333251953125,65504,Infinity 1.5 0.333251953125\n2.25\n" {
		t.Fatalf("got %q: %v", got, err)
	}
}
func TestGoArrayFlattenAndReplace(t *testing.T) {
	source := `const a:any[]=[1,,[2,,[3]],4];console.log(a.flat().join("|"),a.flat(Infinity).join("|"),a.flat(0).length);const sparse:any[]=[1,,3];console.log(sparse.flatMap(x=>[x,x*2]).join(","));const b:any[]=[1,,3];console.log(b.with(-1,7).join("|"),1 in b.with(-1,7),b.toSpliced(1,1,8,9).join("|"),b.join("|"));const x=new Uint8Array([1,2]);console.log(x.with(1,257).join(","),x.join(","));try{b.with(9,0);}catch(e){console.log(e.name);}try{[].reduce((a,b)=>a+b);}catch(e){console.log(e.name);}try{new Uint32Array(new ArrayBuffer(3));}catch(e){console.log(e.name);}`
	want := javascriptOutput(t, source)
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	got, err := runGoProgram(t, text, true)
	if err != nil || got != want {
		t.Fatalf("got %q, want %q: %v", got, want, err)
	}
}
func TestGoArrayLikeConstructors(t *testing.T) {
	source := `const source={length:3,0:7,2:9};console.log(Array.from(source).join("|"),new Float64Array(source).join("|"));const x=new Uint8Array(4);x.set(source,1);console.log(x.join(","),Uint8Array.from("12").join(","),Array(3).length);`
	want := javascriptOutput(t, source)
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	got, err := runGoProgram(t, text, true)
	if err != nil || got != want {
		t.Fatalf("got %q, want %q: %v", got, want, err)
	}
}
