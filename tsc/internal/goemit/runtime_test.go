package goemit

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// Test the exact runtime embedded into programs, including worker failure modes
// that cannot be triggered by TypeScript syntax. The race detector checks that
// pointer-based completion publication does not share mutable loop state.
func TestRuntimeWorkers(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "runtime.go"), []byte("package main\n"+Runtime+ValueRuntime+ModuleRuntime+ClassRuntime+TypeRuntime+StringRuntime+RegexRuntime+EqualityRuntime+ObjectRuntime), 0600); err != nil {
		t.Fatal(err)
	}
	tests := `package main
import("testing";"runtime";"strings")
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
 for i:=0;i<1000;i++ {value:=i;loop.start(func()tsResult{return tsResult{value,false}}).then(func(result tsResult){count++;sum+=result.value.(int)})}
 if err:=loop.run();err!=nil {t.Fatal(err)}
 if count!=1000 || sum!=499500 || len(loop.pending)!=0 {t.Fatalf("lost results: %d %d",count,sum)}
}
func TestFinishFailureDrains(t *testing.T){
 loop:=tsNewLoop();called:=false
 loop.submit(func()tsResult{return tsResult{nil,false}},func(tsResult){panic("finish failure")})
 loop.submit(func()tsResult{return tsResult{nil,false}},func(tsResult){called=true})
 err:=loop.run();if err==nil || !strings.Contains(err.Error(),"finish failure") || !called || len(loop.pending)!=0 {t.Fatalf("drain failed: %v",err)}
}
func TestArbitraryValue(t *testing.T){
 type payload struct{Text string;Number int}
 loop:=tsNewLoop();var got payload
 loop.start(func()tsResult{return tsResult{payload{"value",10},false}}).then(func(result tsResult){got=result.value.(payload)})
 if err:=loop.run();err!=nil || got!=(payload{"value",10}) {t.Fatalf("payload lost: %v %v",got,err)}
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
