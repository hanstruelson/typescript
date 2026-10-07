package goemit

import (
	"context"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Test the exact runtime embedded into programs, including worker failure modes
// that cannot be triggered by TypeScript syntax. The race detector checks that
// pointer-based completion publication does not share mutable loop state.
func TestRuntimeWorkers(t *testing.T) {
	dir := t.TempDir()
	source, err := formatValueSource([]byte("package main\n" + Runtime + ValueRuntime + ModuleRuntime + ClassRuntime + TypeRuntime + ConversionRuntime + GrowableArrayRuntime + StringRuntime + RegexRuntime + EqualityRuntime + ObjectRuntime + CollectionRuntime + SetRuntime + ArrayRuntime + ArrayFlattenRuntime + ArrayLikeRuntime + ArrayBuiltinRuntime + BufferRuntime + TypedArrayBuiltinRuntime + NativeArrayAccessRuntime + DecoratorRuntime))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "runtime.go"), source, 0600); err != nil {
		t.Fatal(err)
	}
	tests := `package main
import("testing";"runtime";"strings";"unsafe";"math")
func TestNativeDateSetters(t *testing.T){
 loop:=tsNewLoop();date:=tsDateConstruct(loop,[]tsValue{tsStringReference(tsStringUTF8("2024-01-31T12:34:56.789Z"))})
 invoke:=func(name string,args ...tsValue)tsValue{return tsCallReceiver(loop,tsGet(loop,date,tsStringReference(tsStringUTF8(name))),date,args...)}
 invoke("setUTCMonth",tsNumberValue(1));if tsText(invoke("toISOString"))!="2024-03-02T12:34:56.789Z"{t.Fatal("month overflow")}
 invoke("setUTCHours",tsNumberValue(25),tsNumberValue(2));if tsText(invoke("toISOString"))!="2024-03-03T01:02:56.789Z"{t.Fatal("hours normalization/defaults")}
 invoke("setUTCSeconds",tsNumberValue(math.NaN()));if !math.IsNaN(tsNumber(invoke("getTime"))){t.Fatal("NaN setter")}
 invoke("setUTCFullYear",tsNumberValue(2000));if tsText(invoke("toISOString"))!="2000-01-01T00:00:00.000Z"{t.Fatal("invalid-date recovery")}
 if tsText(invoke("toGMTString"))!=tsText(invoke("toUTCString")){t.Fatal("GMT alias")}
 if loop.ecma!=nil{t.Fatal("Date initialized regexp engine")}
}
func TestNativeStringNumber(t *testing.T){
 cases:=[]struct{text string;want float64}{{"",0},{"\ufeff\u00a0 42 \u2029",42},{"0xFF",255},{"0o17",15},{"0b11",3},{".5",.5},{"1e309",math.Inf(1)},{"-0",math.Copysign(0,-1)},{"+Infinity",math.Inf(1)},{"1_000",math.NaN()},{"0x1p2",math.NaN()},{"Inf",math.NaN()},{"-0xff",math.NaN()},{"1e",math.NaN()},{"\u00851",math.NaN()}}
 for _,item:=range cases{got:=tsStringNumber(tsStringUTF8(item.text));if math.IsNaN(item.want){if !math.IsNaN(got){t.Errorf("%q: wanted NaN, got %v",item.text,got)}}else if got!=item.want||got==0&&math.Signbit(got)!=math.Signbit(item.want){t.Errorf("%q: wanted %v, got %v",item.text,item.want,got)}}
 loop:=tsNewLoop();if tsStringValue(loop,tsNumberValue(1e21)).String()!="1e+21"{t.Fatal("number formatting")};if loop.ecma!=nil{t.Fatal("numeric conversion initialized regexp engine")}
}
func BenchmarkDenseInt64Push(b *testing.B){storage:=&tsGrowableStorage[int64]{values:make([]int64,0,1)};b.ReportAllocs();b.ResetTimer();for i:=0;i<b.N;i++{storage.values=storage.values[:0];tsDensePush(storage,int64(i))};runtime.KeepAlive(storage)}
func BenchmarkAnyInt64Push(b *testing.B){loop:=tsNewLoop();array:=&tsArray{values:make([]tsValue,0,1)};b.ReportAllocs();b.ResetTimer();for i:=0;i<b.N;i++{array.values=array.values[:0];array.appendItems(loop,tsInt64Value(int64(i)))};runtime.KeepAlive(array)}
func BenchmarkDenseInt64Read(b *testing.B){storage:=&tsGrowableStorage[int64]{values:[]int64{7}};var sum int64;b.ReportAllocs();b.ResetTimer();for i:=0;i<b.N;i++{sum+=tsDenseRead(storage,0)};runtime.KeepAlive(sum)}
func BenchmarkAnyInt64Read(b *testing.B){loop:=tsNewLoop();value:=tsArrayValue(&tsArray{values:[]tsValue{tsInt64Value(7)}});key:=tsNumberValue(0);var sum int64;b.ReportAllocs();b.ResetTimer();for i:=0;i<b.N;i++{sum+=tsNative[int64](tsGet(loop,value,key))};runtime.KeepAlive(sum)}
func BenchmarkDenseInt64NullablePush(b *testing.B){loop:=tsNewLoop();array:=tsNewGrowableArray(loop,"int64",3,false);storage:=(*tsGrowableStorage[int64])(array.native);storage.values=make([]int64,0,1);storage.tags=make([]uint8,0,1);b.ReportAllocs();b.ResetTimer();for i:=0;i<b.N;i++{storage.values=storage.values[:0];storage.tags=storage.tags[:0];tsNullablePush(storage,tsOptional[int64]{value:int64(i)})};runtime.KeepAlive(storage)}
func BenchmarkDenseInt64ForEach(b *testing.B){loop:=tsNewLoop();array:=tsNewGrowableArray(loop,"int64",0,false);storage:=(*tsGrowableStorage[int64])(array.native);storage.values=make([]int64,64);for i:=range storage.values{storage.values[i]=1};var sum int64;callback:=func(i int){sum+=storage.values[i]};b.ReportAllocs();b.ResetTimer();for i:=0;i<b.N;i++{tsDenseForEach(storage,callback)};runtime.KeepAlive(sum)}
func BenchmarkAnyInt64ForEach(b *testing.B){loop:=tsNewLoop();array:=&tsArray{values:make([]tsValue,64)};for i:=range array.values{array.values[i]=tsInt64Value(1)};method:=tsSequenceMethod(tsArrayValue(array),"forEach");var sum int64;callback:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{sum+=tsNative[int64](tsArg(args,0));return tsU});b.ReportAllocs();b.ResetTimer();for i:=0;i<b.N;i++{tsCall(loop,method,tsFunctionValue(callback))};runtime.KeepAlive(sum)}
func TestIteratorUsesConsumerLoop(t *testing.T){
 wanted:=tsNewLoop();called:=false;it:=&tsIterator{pull:func(loop *tsLoop)(tsValue,bool){called=true;if loop!=wanted{t.Fatal("iterator used its creation loop")};return tsU,false}}
 if it.next(wanted)||!called{t.Fatal("iterator pull was not executed")}
}
func TestNullablePackedStorage(t *testing.T){
 loop:=tsNewLoop();array:=tsNewGrowableArray(loop,"int64",3,true);storage:=(*tsGrowableStorage[int64])(array.native)
 tsNullablePush(storage,tsOptional[int64]{value:math.MaxInt64},tsOptional[int64]{tag:1},tsOptional[int64]{tag:2})
 if array.values!=nil||len(storage.values)!=3||len(storage.tags)!=3||unsafe.Sizeof(storage.values[0])+unsafe.Sizeof(storage.tags[0])!=9{t.Fatal("nullable slots are not packed native values and bytes")}
 array.appendItems(loop,tsInt64Value(7));if len(storage.tags)!=4||storage.tags[3]!=0{t.Fatal("dynamic alias lost tag alignment")}
 tsDenseReverse(storage);tsDenseCopyWithin(storage,1,0,2);if tsNullableAt(storage,1).value!=7||tsNullableAt(storage,2).tag!=2{t.Fatal("bulk operations lost tags")}
 array.resize(1);if len(storage.values)!=1||len(storage.tags)!=1{t.Fatal("shrink lost tag alignment")}
 stringsArray:=tsNewGrowableArray(loop,"string",1,false);refs:=(*tsGrowableStorage[*tsString])(stringsArray.native);tsNullablePush(refs,tsOptional[*tsString]{value:tsStringUTF8("retained")},tsOptional[*tsString]{tag:1});runtime.GC();if refs.values[0].String()!="retained"||refs.values[1]!=nil{t.Fatal("GC or absent pointer storage is wrong")};tsNullableFill(refs,tsOptional[*tsString]{tag:1},0,math.Inf(1));if refs.values[0]!=nil{t.Fatal("null retained a pointer")}
}
func TestGrowableNativeStorage(t *testing.T){
 loop:=tsNewLoop();array:=tsNewGrowableArray(loop,"int32",0,true);array.appendItems(loop,tsInt32Value(1),tsInt32Value(2));storage:=(*tsGrowableStorage[int32])(array.native)
 if unsafe.Sizeof(storage.values[0])!=4||array.values!=nil{t.Fatal("array is not stored in native slots")}
 array.appendItems(loop,tsStringReference(tsStringUTF8("3")));if array.length()!=3||array.at(2).kind!=tsInt32Kind{t.Fatal("push failed checked conversion")}
 runtime.GC();if tsNative[int32](array.at(0))!=1{t.Fatal("GC lost native storage")}
 tsCall(loop,tsSequenceMethod(tsArrayValue(array),"reverse"));if tsNative[int32](array.at(0))!=3{t.Fatal("reverse failed")}
}
func checkGrowableSlots[T tsPrimitive](t *testing.T,kind string,width uintptr,initial T){
 loop:=tsNewLoop();array:=tsNewGrowableArray(loop,kind,0,false);value:=tsArrayValue(array)
 for i:=0;i<1024;i++ {tsGrowableArrayPush[T](loop,value,kind,tsPrimitiveValue(initial))}
 storage:=(*tsGrowableStorage[T])(array.native);if unsafe.Sizeof(storage.values[0])!=width||len(storage.values)!=1024||array.values!=nil{t.Fatal("wrong native storage",kind)}
 runtime.GC();got:=tsGrowableArrayRead[T](loop,value,tsNumberValue(1023),kind);if tsNative[T](got)!=initial||tsNative[T](array.at(1023))!=initial{t.Fatal("growth or GC lost the value",kind)}
 tsGrowableArrayWrite[T](loop,value,tsNumberValue(1024),tsPrimitiveValue(initial),kind);if array.length()!=1025||tsNative[T](array.at(1024))!=initial{t.Fatal("indexed append broke shared storage",kind)}
}
func TestPrimitiveSlotWidths(t *testing.T){
 checkGrowableSlots[int8](t,"int8",1,-128);checkGrowableSlots[uint8](t,"uint8",1,255)
 checkGrowableSlots[int16](t,"int16",2,-32768);checkGrowableSlots[uint16](t,"uint16",2,65535)
 checkGrowableSlots[int32](t,"int32",4,-2147483648);checkGrowableSlots[uint32](t,"uint32",4,4294967295)
 checkGrowableSlots[int64](t,"int64",8,math.MaxInt64);checkGrowableSlots[uint64](t,"uint64",8,math.MaxUint64)
 checkGrowableSlots[float32](t,"float32",4,1.5);checkGrowableSlots[float64](t,"number",8,1.5)
 checkGrowableSlots[bool](t,"boolean",1,true);checkGrowableSlots[*tsString](t,"string",unsafe.Sizeof((*tsString)(nil)),tsStringUTF8("retained"))
}
func TestNativeArrayStorageAndGC(t *testing.T){
 array:=tsNativeArray(make([]float64,1024),"Float64Array",func(n float64)float64{return n})
 if unsafe.Sizeof((*tsNativeArrayStorage[float64])(array.storage).values[0])!=8{t.Fatal("wrong float64 slot width")}
 value:=tsTypedArrayValue(array);runtime.GC();view:=array.view(1,3);view.write(0,7.5)
 if array.read(1)!=7.5||view.buffer!=array.buffer||tsNativeArrayRead[float64](tsNewLoop(),value,tsNumberValue(1),"Float64Array").number!=7.5{t.Fatal("native view lost identity or storage")}
 small:=tsNativeArray(make([]int8,4),"Int8Array",func(n float64)int8{return int8(n)})
 if unsafe.Sizeof((*tsNativeArrayStorage[int8])(small.storage).values[0])!=1{t.Fatal("wrong int8 slot width")}
}
func TestFailures(t *testing.T){
 for _,work:=range []func()tsResult{
  func()tsResult{panic("worker failure")},
  func()tsResult{panic(nil)},
  func()tsResult{runtime.Goexit();return tsResult{}},
 }{
  loop:=tsNewLoop();called:=0
  loop.start(work).then(func(result tsResult){called++;if !result.rejected {t.Error("failure became success")}})
  if err:=loop.run();err!=nil {t.Fatal(err)}
  if called!=1 || len(loop.pending)!=0 {t.Fatalf("lost completion: %d",called)}
 }
}
func TestManyCompletions(t *testing.T){
 loop:=tsNewLoop();count,sum:=0,0
 for i:=0;i<1000;i++ {value:=i;loop.start(func()tsResult{return tsResult{tsNumberValue(float64(value)),false}}).then(func(result tsResult){count++;sum+=int(result.value.number)})}
 if err:=loop.run();err!=nil {t.Fatal(err)}
 if count!=1000 || sum!=499500 || len(loop.pending)!=0 {t.Fatalf("lost results: %d %d",count,sum)}
}
func TestFinishFailureDrains(t *testing.T){
 loop:=tsNewLoop();called:=false
 loop.submit(func()tsResult{return tsResult{tsNull,false}},func(tsResult){panic("finish failure")})
 loop.submit(func()tsResult{return tsResult{tsNull,false}},func(tsResult){called=true})
 err:=loop.run();if err==nil || !strings.Contains(err.Error(),"finish failure") || !called || len(loop.pending)!=0 {t.Fatalf("drain failed: %v",err)}
}
func TestValueLayoutAndGC(t *testing.T){
 if unsafe.Sizeof(tsValue{})!=24 || unsafe.Offsetof(tsValue{}.kind)!=0 || unsafe.Offsetof(tsValue{}.number)!=8 || unsafe.Offsetof(tsValue{}.ref)!=16 {t.Fatal("incorrect value ABI")}
 values:=make([]tsValue,1000)
 for i:=range values {values[i]=tsStringReference(tsStringUTF8(strings.Repeat("x",i+1)))}
 runtime.GC()
 for i,value:=range values {if len((*tsString)(value.ref).units)!=i+1 {t.Fatal("lost GC reference")}}
}
func TestInstanceGC(t *testing.T){
 type instance struct{properties *tsProperties;number float64}
 values:=make([]tsValue,1000)
 for i:=range values{self:=&instance{properties:tsNewProperties(),number:float64(i)};self.properties.self=unsafe.Pointer(self);values[i]=tsInstanceValue(self.properties)}
 runtime.GC()
 for i,value:=range values{if (*instance)((*tsProperties)(value.ref).self).number!=float64(i){t.Fatal("lost concrete instance")}}
}
func TestReferenceDispatch(t *testing.T){
 loop:=tsNewLoop()
 object:=tsNewObject();array:=&tsArray{values:[]tsValue{tsNumberValue(3)}}
 object.set(loop,"array",tsArrayValue(array));runtime.GC()
 if tsGet(loop,tsGet(loop,tsObjectValue(object),tsStringReference(tsStringUTF8("array"))),tsNumberValue(0)).number!=3{t.Fatal("nested reference lost")}
 tsSet(loop,tsObjectValue(object),tsStringReference(tsStringUTF8("value")),tsNumberValue(7))
 if tsGet(loop,tsObjectValue(object),tsStringReference(tsStringUTF8("value"))).number!=7{t.Fatal("object mutation lost")}
 for _,value:=range []tsValue{tsU,tsNull,tsNumberValue(1),tsBooleanValue(true),tsStringReference(tsStringUTF8("x")),tsArrayValue(array)}{
  back:=tsFromECMA(tsToECMA(loop,value))
  if value.kind==tsArrayKind{if tsGet(loop,back,tsNumberValue(0)).number!=3{t.Fatal("Goja array boundary")}}else if !tsStrictEqual(value,back){t.Fatal("Goja primitive boundary")}
 }
}
func TestWorkerContext(t *testing.T){
 parent:=tsNewLoop();shared:=&tsObject{values:map[string]tsValue{}};var workerLoop *tsLoop
 function:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{
  if loop==parent{panic("worker reused parent loop")};workerLoop=loop
  object:=(*tsObject)(args[0].ref);object.values["answer"]=tsInt64Value(9223372036854775807)
  return tsPromiseValue(loop.delay(tsNumberValue(1)))
 })
 parent.spawn(tsFunctionValue(function),tsObjectValue(shared)).then(func(r tsResult){if r.rejected{t.Error(tsText(r.value))}})
 if err:=parent.run();err!=nil{t.Fatal(err)}
 if workerLoop==nil||tsText(shared.values["answer"])!="9223372036854775807"{t.Fatal("worker lost shared reference")}
}
func TestPrimitiveTags(t *testing.T){
 loop:=tsNewLoop()
 if tsStrictEqual(tsNumberValue(math.NaN()),tsNumberValue(math.NaN())) || !tsStrictEqual(tsNumberValue(0),tsNumberValue(math.Copysign(0,-1))) {t.Fatal("numeric equality")}
 if tsStrictEqual(tsNumberValue(1),tsBooleanValue(true)) || tsStrictEqual(tsNull,tsU) {t.Fatal("different tags compare equal")}
 if !tsStrictEqual(tsStringReference(tsStringUTF8("abc")),tsStringReference(tsStringUTF8("abc"))) {t.Fatal("string equality uses identity")}
 if tsTruthy(tsU)||tsTruthy(tsNull)||tsTruthy(tsNumberValue(math.NaN()))||tsTruthy(tsBooleanValue(false))||tsTruthy(tsStringReference(tsStringUTF8(""))) {t.Fatal("truthiness")}
 if tsBoundary(loop,tsStringReference(tsStringUTF8("12")),"number",0,true).number!=12 {t.Fatal("coercion")}
 if tsOptionalFrom[float64](tsNull).tag!=1 || tsOptionalFrom[float64](tsU).tag!=2 || tsOptionalFrom[float64](tsNumberValue(2)).value!=2 {t.Fatal("nullable tags")}
}

`
	if err := os.WriteFile(filepath.Join(dir, "runtime_test.go"), []byte(tests), 0600); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"go.mod", "go.sum"} {
		data, err := os.ReadFile(filepath.Join("..", "..", name))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, name), data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", "test", "-race", "-timeout=30s", ".")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOWORK=off", "GO111MODULE=on")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("runtime tests: %v\n%s", err, output)
	}
	if os.Getenv("TSC_GO_ARRAY_BENCH") == "1" {
		bench := exec.CommandContext(ctx, "go", "test", "-run=^$", "-bench=Benchmark(Dense|Any)Int64", "-benchmem", "-benchtime=100ms", "-count=3", ".")
		bench.Dir = dir
		bench.Env = cmd.Env
		out, err := bench.CombinedOutput()
		if err != nil {
			t.Fatalf("array benchmarks: %v\n%s", err, out)
		}
		t.Logf("array benchmarks:\n%s", out)
	}
}

// Only host APIs and generic static-type selection may inspect Go interfaces.
// Application values must never fall back to interface boxing or assertions.
func TestRuntimeInterfaceBoundaries(t *testing.T) {
	source, err := formatValueSource([]byte("package main\n" + Runtime + ValueRuntime + ModuleRuntime + ClassRuntime + TypeRuntime + ConversionRuntime + GrowableArrayRuntime + StringRuntime + RegexRuntime + EqualityRuntime + ObjectRuntime + CollectionRuntime + SetRuntime + ArrayRuntime + ArrayFlattenRuntime + ArrayLikeRuntime + ArrayBuiltinRuntime + BufferRuntime + TypedArrayBuiltinRuntime + NativeArrayAccessRuntime + DecoratorRuntime))
	if err != nil {
		t.Fatal(err)
	}
	for _, legacy := range []string{"tsUnbox", "tsValueOf", "tsDynamicObject"} {
		if strings.Contains(string(source), legacy) {
			t.Fatalf("legacy adapter remains: %s", legacy)
		}
	}
	file, err := parser.ParseFile(token.NewFileSet(), "runtime.go", source, 0)
	if err != nil {
		t.Fatal(err)
	}
	allowed := map[string]bool{"tsUnwrap": true, "tsFromECMA": true, "tsECMAFailure": true, "tsNative": true, "tsPrimitiveValue": true}
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		ast.Inspect(fn.Body, func(node ast.Node) bool {
			if _, ok := node.(*ast.TypeAssertExpr); ok && !allowed[fn.Name.Name] {
				t.Errorf("interface assertion in %s", fn.Name.Name)
			}
			if call, ok := node.(*ast.CallExpr); ok {
				if name, ok := call.Fun.(*ast.Ident); ok && name.Name == "any" {
					if (fn.Name.Name != "tsNative" && fn.Name.Name != "tsPrimitiveValue") || len(call.Args) != 1 || valueType(call.Args[0]) != "zero" {
						t.Errorf("value boxed in %s", fn.Name.Name)
					}
				}
			}
			return true
		})
	}
}
func TestRejectInterfaceValueFallback(t *testing.T) {
	for _, source := range []string{
		"package main;func f(value any)tsValue{return value}",
		"package main;func f(value tsValue){_ = value.(*tsObject)}",
	} {
		if _, err := formatValueSource([]byte(source)); err == nil {
			t.Fatal("interface fallback accepted")
		}
	}
}

func TestMultipleReturnValueABI(t *testing.T) {
	source, err := formatValueSource([]byte("package main\ntype tsValue struct{}\nfunc pair()(tsValue,bool){return tsValue{},false}\n"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(source), "tsBooleanValue") {
		t.Fatal("native boolean return was boxed using the first result type")
	}
}
