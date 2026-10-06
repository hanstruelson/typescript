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
	source, err := formatValueSource([]byte("package main\n" + Runtime + ValueRuntime + ModuleRuntime + ClassRuntime + TypeRuntime + StringRuntime + RegexRuntime + EqualityRuntime + ObjectRuntime + CollectionRuntime + SetRuntime + ArrayRuntime + ArrayFlattenRuntime + ArrayLikeRuntime + ArrayBuiltinRuntime + BufferRuntime + TypedArrayBuiltinRuntime + NativeArrayAccessRuntime + DecoratorRuntime))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "runtime.go"), source, 0600); err != nil {
		t.Fatal(err)
	}
	tests := `package main
import("testing";"runtime";"strings";"unsafe";"math")
func TestIteratorUsesConsumerLoop(t *testing.T){
 wanted:=tsNewLoop();called:=false;it:=&tsIterator{pull:func(loop *tsLoop)(tsValue,bool){called=true;if loop!=wanted{t.Fatal("iterator used its creation loop")};return tsU,false}}
 if it.next(wanted)||!called{t.Fatal("iterator pull was not executed")}
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
}

// Only host APIs and generic static-type selection may inspect Go interfaces.
// Application values must never fall back to interface boxing or assertions.
func TestRuntimeInterfaceBoundaries(t *testing.T) {
	source, err := formatValueSource([]byte("package main\n" + Runtime + ValueRuntime + ModuleRuntime + ClassRuntime + TypeRuntime + StringRuntime + RegexRuntime + EqualityRuntime + ObjectRuntime + CollectionRuntime + SetRuntime + ArrayRuntime + ArrayFlattenRuntime + ArrayLikeRuntime + ArrayBuiltinRuntime + BufferRuntime + TypedArrayBuiltinRuntime + NativeArrayAccessRuntime + DecoratorRuntime))
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
