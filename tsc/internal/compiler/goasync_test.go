package compiler_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/microsoft/TypeScript/tsc/internal/bundled"
	"github.com/microsoft/TypeScript/tsc/internal/compiler"
	"github.com/microsoft/TypeScript/tsc/internal/core"
	"github.com/microsoft/TypeScript/tsc/internal/tsoptions"
	"github.com/microsoft/TypeScript/tsc/internal/tspath"
	"github.com/microsoft/TypeScript/tsc/internal/vfs/vfstest"
)

func emitGoProgram(t *testing.T, source string) string {
	t.Helper()
	fs := bundled.WrapFS(vfstest.FromMap(map[string]string{"/src/input.ts": source}, tspath.CaseSensitive))
	options := core.CompilerOptions{Target: core.ScriptTargetGo, OutDir: "/out"}
	p := compiler.NewProgram(compiler.ProgramOptions{Config: tsoptions.NewParsedCommandLine(&options, []tspath.RootedFilePath{"/src/input.ts"}, nil, "/src", fs.CaseSensitivity()), Host: compiler.NewCompilerHost(fs, bundled.LibPath(), nil, nil, nil)})
	var text string
	result := p.Emit(context.Background(), compiler.EmitOptions{WriteFile: func(path tspath.RootedFilePath, output string, data *compiler.WriteFileData) error {
		if path != "/out/input.go" {
			t.Errorf("unexpected output: %s", path)
		}
		text = output
		return nil
	}})
	if result.EmitSkipped || len(result.Diagnostics) != 0 {
		for _, diag := range result.Diagnostics {
			t.Log(diag.String())
		}
		t.Fatalf("Go emit failed")
	}
	return text
}
func runGoProgram(t *testing.T, text string, race bool) (string, error) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "main.go")
	if err := os.WriteFile(path, []byte(text), 0600); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	args := []string{"run"}
	if race {
		args = append(args, "-race")
	}
	args = append(args, path)
	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Env = append(os.Environ(), "GOWORK=off")
	output, err := cmd.CombinedOutput()
	if ctx.Err() != nil {
		t.Fatalf("generated program hung: %v\n%s", ctx.Err(), output)
	}
	return string(output), err
}
func TestGoAsyncStateMachines(t *testing.T) {
	for _, tc := range []struct{ name, source, want string }{
		{"promise executor and all", `async function f(){let values=await Promise.all([new Promise(resolve=>resolve(10)),Promise.resolve(20)]);console.log(values[0],values[1]);}f();console.log("root");`, "root\n10 20\n"},
		{"nested finally propagation", `async function f(){try{try{await 0;throw "failure";}finally{await 0;console.log("inner");}}catch(e){console.log(e);}finally{await 0;console.log("outer");}}f();`, "inner\nfailure\nouter\n"},
		{"synchronous null throw", `function fail(){throw null;}async function f(){try{fail();}catch(e){console.log(e);}await 0;}f();`, "null\n"},
		{"live array iteration", `async function f(){let values=[1];for(const x of values){await 0;console.log(x);if(x===1)values.push(2);}}f();`, "1\n2\n"},
		{"default parameters", `async function f(a=10,b=a+1){await 0;console.log(a,b);}f();`, "10 11\n"},
		{"var redeclaration", `async function f(){var x=1;await 0;var x=2;console.log(x);}f();`, "2\n"},
		{"synchronous prefix and queued await", `async function f(){ console.log("prefix"); await 10; console.log("resume"); } f(); console.log("root");`, "prefix\nroot\nresume\n"},
		{"sequential awaits and shared mutation", `async function f(){let a=10;let contents=await Promise.resolve(2);a+=contents;await Promise.resolve(3);console.log(a);} f();`, "12\n"},
		{"branch joins", `async function f(flag){let value;if(flag){value=await Promise.resolve(10);}else{value=await Promise.resolve(20);}console.log(value);}f(true);f(false);`, "10\n20\n"},
		{"lexical shadowing", `async function f(){let x=10;{let x=20;await 0;console.log(x);}console.log(x);}f();`, "20\n10\n"},
		{"for loop fresh captures", `async function f(){let saved=[];for(let i=0;i<3;i++){saved.push(()=>i);await 0;}console.log(saved[0](),saved[1](),saved[2]());}f();`, "0 1 2\n"},
		{"var loop shared captures", `async function f(){let saved=[];for(var i=0;i<3;i++){saved.push(()=>i);await 0;}console.log(saved[0](),saved[1](),saved[2]());}f();`, "3 3 3\n"},
		{"while do and loop exits", `async function f(){let i=0;while(i<4){i++;await 0;if(i===2)continue;if(i===4)break;console.log(i);}do{await 0;i--;}while(i>2);console.log(i);}f();`, "1\n3\n2\n"},
		{"labeled loops", `async function f(){outer:for(let i=0;i<3;i++){for(let j=0;j<3;j++){await 0;if(j===1)continue outer;console.log(i,j);}}}f();`, "0 0\n1 0\n2 0\n"},
		{"for of", `async function f(){let saved=[];for(const x of [10,20]){await 0;saved.push(()=>x);}console.log(saved[0](),saved[1]());}f();`, "10 20\n"},
		{"catch rejection and finally", `async function f(){try{await Promise.reject("failure");console.log("bad");}catch(e){console.log(e);}finally{await 0;console.log("finally");}console.log("after");}f();`, "failure\nfinally\nafter\n"},
		{"return through finally", `async function f(){try{return await Promise.resolve(10);}finally{await 0;console.log("finally");}}async function g(){console.log(await f());}g();`, "finally\n10\n"},
		{"finally overrides return", `async function f(){try{return 10;}finally{await 0;return 20;}}async function g(){console.log(await f());}g();`, "20\n"},
		{"loop exit through finally", `async function f(){for(let i=0;i<3;i++){try{await 0;if(i===0)continue;break;}finally{await 0;console.log(i);}}console.log("end");}f();`, "0\n1\nend\n"},
		{"expression evaluation order", `function first(){console.log("first");return 1;}function third(){console.log("third");return 3;}function show(a,b,c){console.log(a,b,c);}async function f(){show(first(),await Promise.resolve(2),third());}f();console.log("root");`, "first\nroot\nthird\n1 2 3\n"},
		{"short circuit", `async function f(){let x=false && await Promise.reject("bad");let y=true ? await 10 : await Promise.reject("bad");console.log(x,y);}f();`, "false 10\n"},
		{"temporal dead zone and const", `async function f(){try{let read=()=>x;read();let x=10;}catch(e){console.log("tdz");}try{const y=1;await 0;y=2;}catch(e){console.log("const");}}f();`, "tdz\nconst\n"},
		{"var initialization timing", `async function f(){console.log(x);await 0;var x=10;console.log(x);}f();`, "undefined\n10\n"},
		{"async return adoption", `async function f(){return Promise.resolve(10);}async function g(){console.log(await f());}g();`, "10\n"},
		{"promise handlers", `async function f(){await 0;throw "failure";}f().catch(e=>console.log(e));Promise.resolve(10).then(x=>console.log(x));console.log("root");`, "root\n10\nfailure\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			text := emitGoProgram(t, tc.source)
			output, err := runGoProgram(t, text, false)
			if err != nil || output != tc.want {
				t.Fatalf("output %q, want %q, error %v\n%s", output, tc.want, err, text)
			}
		})
	}
}
func TestGoAsyncWorkers(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "input.txt")
	if err := os.WriteFile(file, []byte("contents"), 0600); err != nil {
		t.Fatal(err)
	}
	source := `declare function readFile(path:string):Promise<string>;declare function delay(ms:number):Promise<void>;
 async function reader(){let a=10;let contents=await readFile(` + strconv.Quote(file) + `);console.log(a,contents);await delay(0);console.log("done");}reader();console.log("root");`
	text := emitGoProgram(t, source)
	output, err := runGoProgram(t, text, true)
	if err != nil || output != "root\n10 contents\ndone\n" {
		t.Fatalf("worker output: %q, %v", output, err)
	}
	source = `declare function delay(ms:number):Promise<void>; async function f(){for(let i=0;i<100;i++){await delay(0);}console.log("finished");}f();`
	output, err = runGoProgram(t, emitGoProgram(t, source), true)
	if err != nil || output != "finished\n" {
		t.Fatalf("fast completion lost: %q, %v", output, err)
	}
}
func TestGoAsyncUnhandledRejection(t *testing.T) {
	output, err := runGoProgram(t, emitGoProgram(t, `async function f(){await 0;throw "failure";}f();`), false)
	if err == nil || !strings.Contains(output, "Unhandled rejection: failure") {
		t.Fatalf("rejection lost: %q, %v", output, err)
	}
}
