package compiler_test

import (
	"context"
	"regexp"
	"strings"
	"testing"

	"github.com/microsoft/TypeScript/tsc/internal/bundled"
	"github.com/microsoft/TypeScript/tsc/internal/compiler"
	"github.com/microsoft/TypeScript/tsc/internal/core"
	"github.com/microsoft/TypeScript/tsc/internal/tsoptions"
	"github.com/microsoft/TypeScript/tsc/internal/tspath"
	"github.com/microsoft/TypeScript/tsc/internal/vfs/vfstest"
)

func TestGoSpecializedGenerics(t *testing.T) {
	cases := []struct {
		name, source, want string
		native             bool
	}{
		{"identity", `function identity<T>(value:T):T{return value;}console.log(identity(7),identity("hi"),identity(true));`, "7 hi true\n", true},
		{"explicit any", `function identity<T>(value:T):T{return value;}const a:any={value:7};const b=identity<any>(a);console.log(b===a,b.value,identity<any>("text"));`, "true 7 text\n", true},
		{"defaults and multiple parameters", `function choose<T=string>(flag:boolean,a:T,b:T):T{if(flag)return a;return b;}function second<T,U>(a:T,b:U):U{return b;}console.log(choose(true,"a","b"),second<number,string>(3,"four"),choose(false,5,9));`, "a four 9\n", true},
		{"native numeric constraints", `function sum<T extends number>(a:T,b:T):T{const result=a+b;return result as T;}function min<T extends number>(a:T,b:T):T{if(a<b)return a;return b;}console.log(sum(3,4),min(9,2));`, "7 2\n", true},
		{"exact native integer", `function identity<T>(value:T):T{return value;}function sum<T extends int64>(a:T,b:T):T{return (a+b) as T;}const n:int64=9007199254740993;console.log(identity<int64>(n),sum<int64>(n,2));`, "9007199254740993 9007199254740995\n", true},
		{"indexed access", `function pluck<T,K extends keyof T>(value:T,key:K):T[K]{return value[key];}console.log(pluck({name:"ok"},"name"),pluck({count:7},"count"));`, "ok 7\n", true},
		{"conditional return", `type Result<T>=T extends number?number:string;function convert<T>(value:T):Result<T>{return value as any;}console.log(convert(7),convert("text"));`, "7 text\n", true},
		{"dependent local storage", `function read<T extends {value:any}>(x:T):any{const value:T["value"]=x.value;return value;}console.log(read({value:7}),read({value:"text"}),read({value:true}));`, "7 text true\n", true},
		{"recursive specialization", `function factorial<T extends number>(n:T):T{if(n<2)return 1 as T;return (n*factorial<T>((n-1) as T)) as T;}console.log(factorial(5));`, "120\n", true},
		{"nested specialization", `function identity<T>(x:T):T{return x;}function outer<T>(x:T):T{return identity<T>(x);}console.log(outer(7),outer("ok"));`, "7 ok\n", true},
		{"inferred result", `function identity<T>(x:T){return x;}console.log(identity(7),identity("ok"));`, "7 ok\n", true},

		{"string constraint", `function text<T extends string>(value:T):string{return value;}console.log(text("ok"));`, "ok\n", true},
		{"captured binding", `let offset=2;function add<T extends number>(value:T):number{return value+offset;}console.log(add(3));offset=7;console.log(add(3));`, "5\n10\n", true},
		{"reassignment", `function identity<T>(value:T):T{return value;}identity=<T>(value:T):T=>value;console.log(identity(8));`, "8\n", false},
		{"first class and extra arguments", `let count=0;function identity<T>(value:T):T{return value;}const ref=identity;function side(){count++;return 4;}console.log(ref(3),identity(8,side()),count);`, "3 8 1\n", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			source := emitGoOptions(t, tc.source, core.CompilerOptions{Strict: core.TSTrue})
			hasNative := regexp.MustCompile(`func tsSpecialized\d+_[0-9a-f]+\(`).MatchString(source)
			if hasNative != tc.native {
				t.Fatalf("concrete specialization = %v, want %v", hasNative, tc.native)
			}
			if strings.Contains(source, "[any]") {
				t.Fatal("TypeScript any emitted as Go any")
			}
			if tc.name == "explicit any" && !regexp.MustCompile(`func tsSpecialized\d+_[0-9a-f]+\(a0 tsValue\)`).MatchString(source) {
				t.Fatal("any instantiation did not use tsValue")
			}
			if tc.name == "exact native integer" && !regexp.MustCompile(`func tsSpecialized\d+_[0-9a-f]+\(a0 int64`).MatchString(source) {
				t.Fatal("native integer instantiation was erased")
			}
			got, err := runGoProgram(t, source, true)
			if err != nil || got != tc.want {
				t.Fatalf("got %q, want %q, error %v", got, tc.want, err)
			}
		})
	}
}

func TestGoSpecializationCache(t *testing.T) {
	source := emitGoOptions(t, `function identity<T>(x:T):T{return x;}console.log(identity(1),identity(2),identity<number>(3));`, core.CompilerOptions{Strict: core.TSTrue})
	if count := len(regexp.MustCompile(`func tsSpecialized\d+_[0-9a-f]+\(`).FindAllString(source, -1)); count != 1 {
		t.Fatalf("emitted %d specializations for the same numeric representation", count)
	}
	if !regexp.MustCompile(`func tsSpecialized\d+_[0-9a-f]+\(a0 float64\) float64 \{\s*return a0\s*\}`).MatchString(source) {
		t.Fatal("primitive specialization did not emit an unboxed direct return")
	}
	got, err := runGoProgram(t, source, true)
	if err != nil || got != "1 2 3\n" {
		t.Fatalf("%q %v", got, err)
	}
}
func TestGoSpecializationExpansionLimit(t *testing.T) {
	// The program deliberately grows its type at every recursive call. Emission
	// must terminate even though executing the program would never terminate.
	source := emitGoOptions(t, `function grow<T>(x:T):number{return grow<T[]>([x]);}if(false)grow(1);`, core.CompilerOptions{Strict: core.TSTrue})
	if count := len(regexp.MustCompile(`func tsSpecialized\d+_[0-9a-f]+\(`).FindAllString(source, -1)); count > 64 {
		t.Fatalf("unbounded recursive specialization: %d", count)
	}
	if got, err := runGoProgram(t, source, true); err != nil || got != "" {
		t.Fatalf("bounded emission did not compile: %q %v", got, err)
	}
}

func TestGoSpecializationsAcrossModules(t *testing.T) {
	fs := bundled.WrapFS(vfstest.FromMap(map[string]string{
		"/src/main.ts": `import "./a.js";import "./b.js";`,
		"/src/a.ts":    `function identity<T>(x:T):T{return x;}console.log(identity(1));export {};`,
		"/src/b.ts":    `function identity<T>(x:T):T{return x;}console.log(identity(2));export {};`,
	}, tspath.CaseSensitive))
	options := core.CompilerOptions{Target: core.ScriptTargetGo, OutDir: "/out", Strict: core.TSTrue}
	p := compiler.NewProgram(compiler.ProgramOptions{Config: tsoptions.NewParsedCommandLine(&options, []tspath.RootedFilePath{"/src/main.ts"}, nil, "/src", fs.CaseSensitivity()), Host: compiler.NewCompilerHost(fs, bundled.LibPath(), nil, nil, nil)})
	text := ""
	result := p.Emit(context.Background(), compiler.EmitOptions{WriteFile: func(path tspath.RootedFilePath, output string, data *compiler.WriteFileData) error {
		text = output
		return nil
	}})
	if result.EmitSkipped || len(result.Diagnostics) > 0 {
		t.Fatalf("module emission failed: %v", result.Diagnostics)
	}
	got, err := runGoProgram(t, text, true)
	if err != nil || got != "1\n2\n" {
		t.Fatalf("%q %v", got, err)
	}
}
