package compiler_test

import (
	"context"
	"strings"
	"testing"

	"github.com/microsoft/TypeScript/tsc/internal/bundled"
	"github.com/microsoft/TypeScript/tsc/internal/compiler"
	"github.com/microsoft/TypeScript/tsc/internal/core"
	"github.com/microsoft/TypeScript/tsc/internal/tsoptions"
	"github.com/microsoft/TypeScript/tsc/internal/tspath"
	"github.com/microsoft/TypeScript/tsc/internal/vfs/vfstest"
)

func emitGoOptions(t *testing.T, source string, options core.CompilerOptions) string {
	t.Helper()
	options.Target = core.ScriptTargetGo
	options.OutDir = "/out"
	fs := bundled.WrapFS(vfstest.FromMap(map[string]string{"/src/input.ts": source}, tspath.CaseSensitive))
	p := compiler.NewProgram(compiler.ProgramOptions{Config: tsoptions.NewParsedCommandLine(&options, []tspath.RootedFilePath{"/src/input.ts"}, nil, "/src", fs.CaseSensitivity()), Host: compiler.NewCompilerHost(fs, bundled.LibPath(), nil, nil, nil)})
	var text string
	result := p.Emit(context.Background(), compiler.EmitOptions{WriteFile: func(path tspath.RootedFilePath, output string, data *compiler.WriteFileData) error {
		text = output
		return nil
	}})
	if result.EmitSkipped || len(result.Diagnostics) > 0 {
		for _, diag := range result.Diagnostics {
			t.Log(diag.String())
		}
		t.Fatal("emit failed")
	}
	return text
}
func TestGoTypedValues(t *testing.T) {
	source := `let fast:number=10;let nullable:number|null|undefined=null;let dynamic:any=3;console.log(fast+2,nullable+2,dynamic+2);nullable=4;console.log(nullable+fast);nullable=undefined;console.log(nullable,nullable+fast,typeof fast);function sum(a:number,b:number):number{return a+b;}console.log(sum(2,3));function accept(a:string|null|undefined){console.log(a,typeof a);}accept("ok");accept(null);accept(undefined);function strict(a:string){console.log(a);}try{strict(dynamic);}catch(e){console.log(e.message);}try{sum("wrong" as any,1);}catch(e){console.log("number rejected");}`
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	if !strings.Contains(text, "tsTypedCell[float64]") || !strings.Contains(text, "tsOptional[float64]") {
		t.Fatal("missing native/nullable numeric paths")
	}
	got, err := runGoProgram(t, text, true)
	if err != nil {
		t.Fatalf("%v\n%s", err, got)
	}
	want := "12 2 5\n14\nundefined NaN number\n5\nok string\nnull object\nundefined undefined\nThis value only accepts a string. Set coerceAny to true to enable automatic conversion.\nnumber rejected\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
func TestGoCoerceAny(t *testing.T) {
	source := `function text(value:string){console.log(value,typeof value);}function number(value:number){console.log(value,typeof value);}let value:any=10;text(value);value="12";number(value);text(null as any);text(undefined as any);`
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue, CoerceAny: core.TSTrue})
	got, err := runGoProgram(t, text, true)
	if err != nil {
		t.Fatalf("%v\n%s", err, got)
	}
	if want := "10 string\n12 number\nnull string\nundefined string\n"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
func TestGoCoerceOption(t *testing.T) {
	fs := vfstest.FromMap(map[string]string{"/src/input.ts": ""}, tspath.CaseSensitive)
	config := tsoptions.ParseCommandLine([]string{"--target", "go", "--coerceAny", "/src/input.ts"}, fs, "/src")
	if len(config.Errors) > 0 || !config.CompilerOptions().CoerceAny.IsTrue() {
		t.Fatalf("CLI option failed: %+v", config)
	}
	config = tsoptions.ParseJsonConfigFileContent(map[string]any{"compilerOptions": map[string]any{"target": "go", "coerceAny": true}, "files": []any{"input.ts"}}, fs, "/src", nil, "/src/tsconfig.json", nil, nil)
	if len(config.Errors) > 0 || !config.CompilerOptions().CoerceAny.IsTrue() {
		t.Fatalf("config option failed: %+v", config)
	}
}
func TestGoUTF16Strings(t *testing.T) {
	source := `const s="A😀B";console.log(s.length,s.charCodeAt(1),s.charCodeAt(2),s.codePointAt(1));console.log(s.slice(1,3),s.substring(3,1),s.at(-1));console.log(s.indexOf("B"),s.includes("😀"),s.startsWith("A"),s.endsWith("B"));console.log("x".repeat(3),"x".padStart(3,"ab"),"x".padEnd(3,"ab"));console.log("\uFEFF \tx \n".trim());console.log("a,b,c".split(",")[1],s.split("").length);const lone="\ud800";console.log(lone.length,lone.charCodeAt(0),lone.isWellFormed(),lone.toWellFormed().charCodeAt(0));console.log(String.fromCharCode(0xd83d,0xde00),String.fromCodePoint(0x1f600));console.log("😀"==="\ud83d\ude00","é".normalize("NFD").length,"ß".toUpperCase(),"ΟΣ".toLowerCase());for(const c of s){console.log(c,c.length);}console.log("😀"<"\uE000");console.log("abcdef".substr(-3,2),"abc".charAt(9),"abc".charCodeAt(9));console.log("a\"b".anchor("x\"y"));console.log("x".bold(),"x".big(),"x".blink(),"x".fixed(),"x".fontcolor("red"),"x".fontsize(2),"x".italics(),"x".link("url"),"x".small(),"x".strike(),"x".sub(),"x".sup());console.log("foo".replace("o","X"),"foo".replaceAll("o","X"));`
	source += "console.log(`hello ${s} ${10}`);"
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	got, err := runGoProgram(t, text, true)
	if err != nil {
		t.Fatalf("%v\n%s", err, got)
	}
	want := javascriptOutput(t, source)
	if got != want {
		t.Fatalf("Go:\n%s\nJavaScript:\n%s", got, want)
	}
}
func TestGoECMAScriptRegex(t *testing.T) {
	source := `const regex=/(?<word>a+)(b)\2/g;const match=regex.exec("zaabb");console.log(match[0],match.index,match.groups.word,regex.lastIndex);console.log(regex.test("zaabb"),regex.lastIndex);const look=/(?<=a)b/;console.log(look.test("ab"));console.log("ab ab".replace(/(a)(b)/g,"$2$1"));console.log("ab".replace(/(a)/,(match,value,index)=>value.toUpperCase()+index));console.log("a1b2".split(/(\d)/).length,"abc".search(/b/));for(const m of "a1a2".matchAll(/a(\d)/g)){console.log(m[0],m[1],m.index);}console.log(new RegExp("a+","i").test("AA"));console.log("😀".match(/./gu)[0]);try{new RegExp("[","g");}catch(e){console.log(e.name);}console.log(/x/.source,/x/g.flags);`
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	got, err := runGoProgram(t, text, true)
	if err != nil {
		t.Fatalf("%v\n%s", err, got)
	}
	want := javascriptOutput(t, source)
	if got != want {
		t.Fatalf("Go:\n%s\nJavaScript:\n%s", got, want)
	}
}

func TestGoSyntaxLowering(t *testing.T) {
	source := `const original={z:3,a:1,"2":2,"1":1};const {a:renamed=10,z,...rest}=original;const [first,,third=30,...tail]=[1,2,undefined,4,5];console.log(renamed,z,rest["1"],first,third,tail[0],tail[1]);const copy={...original,z:9};for(const key in copy){console.log(key,copy[key]);}function unpack({value=10},{nested:{x}}={nested:{x:3}},...items){console.log(value,x,items.length,items[0]);}unpack({},undefined,4,5);function sum(a:number,b:number):number{return a+b;}console.log(sum(...[2,3]));console.log([..."A😀B"].length);let result=0;for(let i=0;i<3;i++){switch(i){case 0:result+=1;break;case 1:result+=2;continue;default:result+=3;}console.log("loop",i);}console.log(result);switch(2){case 1:console.log("wrong");break;default:console.log("default");case 2:console.log("two");case 3:console.log("fallthrough");}console.log(Object.keys(copy)[0],Object.values(copy)[0],Object.entries(copy)[0][0],Object.hasOwn(copy,"a"));const weird={"\ud800":1,"\udc00":2};console.log(weird["\ud800"],weird["\udc00"],Object.keys(weird)[0].charCodeAt(0));console.log(String.raw({raw:["a","b"]}),String.raw({raw:["a","b"]},"X"));`
	source += "console.log(String.raw`a\\nb${10}c`);"
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	got, err := runGoProgram(t, text, true)
	if err != nil {
		t.Fatalf("%v\n%s", err, got)
	}
	want := javascriptOutput(t, source)
	if got != want {
		t.Fatalf("Go:\n%s\nJavaScript:\n%s", got, want)
	}
}
func TestGoPrimitiveClassFields(t *testing.T) {
	source := `class A {name:string="😀";age:number=10;score:number|null=null;value:any=3;flag:boolean=true;print(text:string){console.log(this.name,this.name.length,this.age,this.score,this.value,this.flag,text);}}class B extends A {name:string="B";score:number|null=2;print(text:string){super.print(text+"!");}}const b=new B();b.print("hello");b["name"]="😀";b["score"]=null;b["value"]="dynamic";b.print("bye");console.log(Object.keys(b).length);try{b["age"]="wrong";}catch(e){console.log(e.name);}const f=()=>1;console.log(f===f,f===(()=>1));type Maybe=number|null|undefined;let n:Maybe=null;console.log(n,n+1);n=undefined;console.log(n,n+1);`
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	got, err := runGoProgram(t, text, true)
	if err != nil {
		t.Fatalf("%v\n%s", err, got)
	}
	want := "B 1 10 2 3 true hello!\n😀 2 10 null dynamic true bye!\n5\nTypeError\ntrue false\nnull 1\nundefined NaN\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestGoUnicodeAndConversions(t *testing.T) {
	source := `console.log("I".toLocaleLowerCase("tr"),"i".toLocaleUpperCase("tr"));console.log("e\u0301".normalize(),"\ud800a".toUpperCase().charCodeAt(0));console.log("a".localeCompare("A","en",{sensitivity:"base"}),"2".localeCompare("10","en",{numeric:true})<0);console.log(String(-0),String(Infinity),String(1e21),String(0.000001),String([1,null,undefined,3]),String({}));console.log(String.fromCharCode(Infinity).charCodeAt(0),"abc".split("",Infinity).length);function numeric(x:number){console.log(x);}let value:any="0x10";numeric(value);value="0b11";numeric(value);value="0o10";numeric(value);for(const [a,b] of [[1,2],[3,4]]){console.log(a,b);}for(const key in null){console.log("bad");}`
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue, CoerceAny: core.TSTrue})
	got, err := runGoProgram(t, text, true)
	if err != nil {
		t.Fatalf("%v\n%s", err, got)
	}
	// coerceAny deliberately converts values at typed calls; vanilla JS does not.
	want := strings.Replace(javascriptOutput(t, source), "0x10\n0b11\n0o10\n", "16\n3\n8\n", 1)
	if got != want {
		t.Fatalf("Go:\n%s\nJavaScript:\n%s", got, want)
	}
}

func TestGoNativeFunctionLayout(t *testing.T) {
	text := emitGoOptions(t, `function add(a:number,b:number):number{return a+b;}console.log(add(2,3));`, core.CompilerOptions{Strict: core.TSTrue})
	if !strings.Contains(text, "func(float64, float64) float64") {
		t.Fatal("missing native function signature")
	}
	if !strings.Contains(text, ".read()") || !strings.Contains(text, ".initNative(nativeArg0)") {
		t.Fatal("native function boxed primitive parameters")
	}
}

func TestGoDestructuredParameterInitialization(t *testing.T) {
	source := `function fail({a=b},b=2){console.log(a,b);}try{fail({});}catch(e){console.log("caught");}function good({a=1},{b=a}={}){console.log(a,b);}good({});`
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	got, err := runGoProgram(t, text, true)
	if err != nil {
		t.Fatalf("%v\n%s", err, got)
	}
	want := javascriptOutput(t, source)
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestGoTypedInitializerBoundary(t *testing.T) {
	source := `let value:any=null;try{let s:string=value;console.log("bad",s);}catch(e){console.log(e.name);}value=undefined;try{const s:string=value;console.log("bad",s);}catch(e){console.log(e.name);}`
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	got, err := runGoProgram(t, text, true)
	if err != nil {
		t.Fatalf("%v\n%s", err, got)
	}
	if got != "TypeError\nTypeError\n" {
		t.Fatalf("unexpected output %q", got)
	}
	text = emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue, CoerceAny: core.TSTrue})
	got, err = runGoProgram(t, text, true)
	if err != nil {
		t.Fatalf("%v\n%s", err, got)
	}
	if got != "bad null\nbad undefined\n" {
		t.Fatalf("unexpected coercion output %q", got)
	}
}
