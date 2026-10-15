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
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue, CoerceAny: core.TSFalse})
	if !strings.Contains(text, "tsTypedCell[float64]") || !strings.Contains(text, "tsOptional[float64]") {
		t.Fatal("missing native/nullable numeric paths")
	}
	got, err := runGoProgram(t, text, true)
	if err != nil {
		t.Fatalf("%v\n%s", err, got)
	}
	want := "12 2 5\n14\nundefined NaN number\n5\nok string\nnull object\nundefined undefined\nCannot implicitly convert number to string with coerceAny disabled; use .string()\nnumber rejected\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
func TestGoCoerceAny(t *testing.T) {
	source := `function text(value:string){console.log(value,typeof value);}function number(value:number){console.log(value,typeof value);}let value:any=10;text(value);value="12";number(value);text(null as any);text(undefined as any);`
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
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
func TestGoComparisonEvaluationOrder(t *testing.T) {
	source := `let calls=0;function left():number{calls++;console.log("left");return 12;}function right():any{calls++;console.log("right");return "12";}try{console.log(left()===right());}catch(e){console.log(e.name);}console.log(calls);`
	for _, tc := range []struct {
		policy core.Tristate
		result string
	}{{core.TSUnknown, "true"}, {core.TSFalse, "TypeError"}} {
		text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue, CoerceAny: tc.policy})
		got, err := runGoProgram(t, text, true)
		want := "left\nright\n" + tc.result + "\n2\n"
		if err != nil || got != want {
			t.Fatalf("got %q (%v), want %q", got, err, want)
		}
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

func TestGoRegExpSymbolProtocols(t *testing.T) {
	source := `const r=/a/g;const custom:any={lastIndex:0,[Symbol.match](s:string){return ["custom",s]},[Symbol.replace](s:string,replacement:string){return s+replacement}};console.log("ba".match(custom).join("|"),"x".replace(custom,"!"));r.lastIndex=1;console.log(r.test("ba"),r.lastIndex);const sticky=/a/y;sticky.lastIndex=1;console.log(sticky.exec("ba")[0],sticky.lastIndex);try{"x".match({[Symbol.match]:1} as any);}catch(e){console.log(e.name)};`
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
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue, CoerceAny: core.TSFalse})
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
	if !strings.Contains(text, "func(*tsLoop, float64, float64) float64") {
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
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue, CoerceAny: core.TSFalse})
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

func TestGoMixedComparisonPolicy(t *testing.T) {
	source := `let n:number=12;let v:any="12";try{console.log(n===v,v===n,n<=v,v>=n);}catch(e){console.log(e.name);}let flag:boolean=true;v=1;try{console.log(flag===v,v!==flag);}catch(e){console.log(e.name);}let text:string="12";v=12;try{console.log(text===v,v===text);}catch(e){console.log(e.name);}let maybe:number|null=null;v=null;console.log(maybe===v,maybe===null);let dynamic:any=12;v="12";console.log(dynamic===v);`
	for _, tc := range []struct {
		name   string
		policy core.Tristate
		want   string
	}{
		{"default", core.TSUnknown, "true true true true\ntrue false\ntrue true\ntrue true\nfalse\n"},
		{"enabled", core.TSTrue, "true true true true\ntrue false\ntrue true\ntrue true\nfalse\n"},
		{"disabled", core.TSFalse, "TypeError\nTypeError\nTypeError\ntrue true\nfalse\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue, CoerceAny: tc.policy})
			got, err := runGoProgram(t, text, true)
			if err != nil || got != tc.want {
				t.Fatalf("got %q (%v), want %q", got, err, tc.want)
			}
		})
	}
}

func TestGoNativeEqualityPaths(t *testing.T) {
	source := `const fastA:number=5;const fastB:number=5;let nullableA:number|null|undefined=5;let nullableB:number|null|undefined=5;let dynamic:any="5";let flagA:boolean=true;let flagB:boolean=false;let textA:string="same";let textB:string="same";console.log(fastA==fastB,fastA!=fastB,fastA===fastB,fastA!==fastB,nullableA==nullableB,nullableA!=nullableB,nullableA===nullableB,nullableA!==nullableB,fastA==dynamic,fastA!=dynamic,fastA===dynamic,fastA!==dynamic);console.log(flagA==flagB,flagA!=flagB,flagA===flagB,flagA!==flagB,flagA<flagB,flagA<=flagB,flagA>flagB,flagA>=flagB,flagB<flagA,flagB<=flagA,flagB>flagA,flagB>=flagA,textA==textB,textA!=textB,textA===textB,textA!==textB);nullableA=null;dynamic=undefined;console.log(nullableA==dynamic,nullableA!=dynamic,nullableA===dynamic,nullableA!==dynamic);dynamic=null;console.log(nullableA==dynamic,nullableA!=dynamic,nullableA===dynamic,nullableA!==dynamic);nullableB=undefined;console.log(nullableA==nullableB,nullableA!=nullableB,nullableA===nullableB,nullableA!==nullableB);`
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	for _, op := range []string{"==", "!=", "===", "!=="} {
		if strings.Contains(text, `tsruntime.TsBinary(loop, "`+op+`",`) || strings.Contains(text, `TsPolicyCompare("`+op+`",`) {
			t.Fatalf("typed equality %s fell back to the generic comparison path", op)
		}
	}
	if !strings.Contains(text, "TsPolicyEqualNumberAny(loop,") || !strings.Contains(text, "TsPolicyEqualOptionalNumberAny(loop,") {
		t.Fatalf("number-versus-any equality did not use its specialized runtime helpers")
	}
	got, err := runGoProgram(t, text, false)
	if err != nil {
		t.Fatalf("%v\n%s", err, got)
	}
	want := "true false true false true false true false true false true false\nfalse true false true false false true true true true false false true false true false\ntrue false false true\ntrue false true false\ntrue false false true\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestGoNativeNonnumericComparisonPaths(t *testing.T) {
	source := `let textA:string|null|undefined="same";let textB:string|null|undefined="same";let flagA:boolean|null|undefined=false;let flagB:boolean|null|undefined=null;const flag:boolean=true;const number:int32=1;let maybeNumber:int64|null|undefined=null;console.log(textA==textB,textA===textB,textA!==textB);console.log(textA<textB,textA<=textB,textA>textB,textA>=textB);textA=null;textB="2";console.log(textA<textB,textA<=textB,textA>textB,textA>=textB);textB=undefined;console.log(textA==textB,textA===textB,textA!==textB);console.log(flagA<flagB,flagA<=flagB,flagA>flagB,flagA>=flagB);flagA=true;flagB=false;console.log(flagA<flagB,flagA<=flagB,flagA>flagB,flagA>=flagB);flagA=undefined;console.log(flagA<flagB,flagA<=flagB,flagA>flagB,flagA>=flagB);console.log(flag==number,flag!=number,flag<number,flag<=number,flag>number,flag>=number);console.log(flagB==maybeNumber,flagB===maybeNumber,flagB<maybeNumber,flagB<=maybeNumber,flagB>maybeNumber,flagB>=maybeNumber);flagA=true;maybeNumber=1;console.log(flagA==maybeNumber,flagA===maybeNumber,flagA<maybeNumber,flagA<=maybeNumber,flagA>maybeNumber,flagA>=maybeNumber);maybeNumber=undefined;console.log(flagA<maybeNumber,flagA<=maybeNumber,flagA>maybeNumber,flagA>=maybeNumber);console.log(number===flag,number!==flag,textA===flag,textA!==flag);`
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	for _, op := range []string{"==", "!=", "===", "!==", "<", "<=", ">", ">="} {
		if strings.Contains(text, `TsPolicyCompare("`+op+`"`) || strings.Contains(text, `tsruntime.TsBinary(loop, "`+op+`",`) {
			t.Fatalf("statically known nonnumeric comparison %s used the generic path", op)
		}
	}
	if !strings.Contains(text, "TsStringEqual(") || !strings.Contains(text, ".Tag == 0") {
		t.Fatal("nullable string equality did not use content and tag fast paths")
	}
	got, err := runGoProgram(t, text, true)
	if err != nil {
		t.Fatalf("%v\n%s", err, got)
	}
	want := "true true false\nfalse true false true\ntrue true false false\ntrue false true\nfalse true false true\nfalse false true true\nfalse false false false\ntrue false false true false true\nfalse false false true false true\ntrue false false true false true\nfalse false false false\nfalse true false true\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestGoTypedDynamicComparisonDispatch(t *testing.T) {
	source := `const small:int32=7;let dynamic:any=7;const flag:boolean=true;let maybe:string|null|undefined="x";let dynamicText:any="x";console.log(small==dynamic,small===dynamic,small>=dynamic);console.log(flag==dynamic,flag===dynamic);maybe=undefined;dynamicText=null;console.log(maybe==dynamicText,maybe===dynamicText);`
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue, CoerceAny: core.TSTrue})
	for _, helper := range []string{"TsPolicyCompareNative[int32]", "TsPolicyCompareNative[bool]", "TsPolicyCompareOptionalNative[*tsruntime.TsString]"} {
		if !strings.Contains(text, helper) {
			for _, line := range strings.Split(text, "\n") {
				if strings.Contains(line, "Compare") || strings.Contains(line, "Binary") {
					t.Log(line)
				}
			}
			t.Fatalf("typed/dynamic comparison did not use shared native dispatch %s", helper)
		}
	}
	got, err := runGoProgram(t, text, true)
	if err != nil || got != "true true true\ntrue true\ntrue false\n" {
		t.Fatalf("got %q: %v", got, err)
	}
}

func TestGoTypedDynamicComparisonPreservesObjectCoercion(t *testing.T) {
	source := `let log:string="";const text:string="1";let dynamic:any={valueOf:()=>{log+="V";return 1;}};console.log(text==dynamic,log);log="";const number:number=1;console.log(dynamic>=number,log);`
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue, CoerceAny: core.TSTrue})
	got, err := runGoProgram(t, text, true)
	if err != nil || got != "true V\ntrue V\n" {
		t.Fatalf("got %q: %v", got, err)
	}
}

func TestGoTypedBitwiseLowering(t *testing.T) {
	source := `const a:int8=-2;const b:int8=3;const minusOne:int8=-1;const shift:int8=1;const f:float32=3.5;const n:number=1.5;const flag:boolean=true;const off:boolean=false;console.log(a&b,a|b,a^b,a<<shift,a>>shift,minusOne>>>shift);console.log(f&n,~f,flag&off,flag|off);try{console.log(a<<minusOne);}catch(e){console.log(e.name);}`
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	if !strings.Contains(text, "TsNativeIntegerBitwise[int8]") || !strings.Contains(text, "TsNumberBitwise(") {
		for _, line := range strings.Split(text, "\n") {
			if strings.Contains(line, "tsNativeIntegerBitwise") || strings.Contains(line, "tsNumberBitwise") || strings.Contains(line, "tsBinary(\"&\"") {
				t.Log(line)
			}
		}
		t.Fatal("typed integer and Number bitwise paths were not selected")
	}
	got, err := runGoProgram(t, text, true)
	want := "2 -1 -3 -4 -1 127\n1 -4 0 1\nRangeError\n"
	if err != nil || got != want {
		t.Fatalf("got %q want %q: %v", got, want, err)
	}
}

func TestGoNativeTruthinessPaths(t *testing.T) {
	source := `const zero:int32=0;const nan:number=NaN;const empty:string="";let maybe:boolean|null|undefined=undefined;console.log(!zero,!nan,!empty,!maybe);console.log(zero||2,empty||"fallback",maybe??true);`
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	if !strings.Contains(text, "!= 0") || !strings.Contains(text, "math.IsNaN(float64(") || !strings.Contains(text, "len(") || !strings.Contains(text, ".Tag == 0") || strings.Contains(text, "TsTruthy(") {
		for _, line := range strings.Split(text, "\n") {
			if strings.Contains(line, "tsTruthy") || strings.Contains(line, "tsNullish") || strings.Contains(line, ".tag!=0") || strings.Contains(line, "len(") || strings.Contains(line, "math.IsNaN(float64(") {
				t.Log(line)
			}
		}
		t.Fatal("known primitive truthiness/nullish paths were not emitted directly")
	}
	got, err := runGoProgram(t, text, true)
	if err != nil || got != "true true true true\n2 fallback true\n" {
		t.Fatalf("got %q: %v", got, err)
	}
}

func TestGoNumericComparisonPlans(t *testing.T) {
	source := `const signed:int32=2147483647;const unsigned:uint32=4294967295;const rounded:float32=0.1;const exact:float64=0.1;const wide:int64=9007199254740993;const nearby:float64=9007199254740992;const minimum:int64=-9223372036854775808;const minimumFloat:float64=-9223372036854775808;const signedMax:int64=9223372036854775807;const signedLimit:float64=9223372036854775808;const unsignedMax:uint64=18446744073709551615;const unsignedLimit:float64=18446744073709551616;let optional:int64|null|undefined=null;console.log(signed<unsigned,signed<=unsigned,signed>unsigned,signed>=unsigned,signed===unsigned,signed!=unsigned);console.log(rounded>exact,rounded===exact);console.log(wide<nearby,wide<=nearby,wide>nearby,wide>=nearby,wide==nearby,wide!=nearby,wide===nearby,wide!==nearby);console.log(nearby<wide,nearby<=wide,nearby>wide,nearby>=wide,nearby==wide,nearby!=wide,nearby===wide,nearby!==wide);console.log(minimum<minimumFloat,minimum<=minimumFloat,minimum>minimumFloat,minimum>=minimumFloat,minimum==minimumFloat,minimum!=minimumFloat,minimum===minimumFloat,minimum!==minimumFloat);console.log(signedMax<signedLimit,signedMax<=signedLimit,signedMax>signedLimit,signedMax>=signedLimit,signedMax==signedLimit,signedMax!=signedLimit);console.log(unsignedMax<unsignedLimit,unsignedMax<=unsignedLimit,unsignedMax>unsignedLimit,unsignedMax>=unsignedLimit,unsignedMax==unsignedLimit,unsignedMax!=unsignedLimit);console.log(unsignedLimit<unsignedMax,unsignedLimit<=unsignedMax,unsignedLimit>unsignedMax,unsignedLimit>=unsignedMax,unsignedLimit==unsignedMax,unsignedLimit!=unsignedMax);console.log(optional<=0,optional==undefined,optional===null);optional=undefined;console.log(optional<0,optional==null,optional===undefined);`
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	if !strings.Contains(text, `TsNumericComparison(">"`) {
		t.Fatal("int64/float64 comparison did not use the exact numeric fallback")
	}
	if !strings.Contains(text, "int64(") || !strings.Contains(text, "float64(") {
		t.Fatal("safe third-target numeric widening was not emitted")
	}
	got, err := runGoProgram(t, text, true)
	if err != nil {
		t.Fatalf("%v\n%s", err, got)
	}
	want := "true true false false false true\ntrue false\nfalse false true true false true false true\ntrue true false false false true false true\nfalse true false true true false true false\ntrue true false false false true\ntrue true false false false true\nfalse false true true false true\ntrue true true\nfalse true true\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestGoDynamicEqualityDispatch(t *testing.T) {
	source := `let left:any="5";let right:any=5;console.log(left==right,left!=right,left===right,left!==right);console.log(right==left,right!=left,right===left,right!==left);`
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	if !strings.Contains(text, `tsruntime.TsBinary(loop, "==",`) {
		t.Fatal("dynamic equality did not use the factored runtime dispatch")
	}
	got, err := runGoProgram(t, text, false)
	if err != nil {
		t.Fatalf("%v\n%s", err, got)
	}
	want := "true false false true\ntrue false false true\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestGoNativeNumericConversions(t *testing.T) {
	source := `let big:int64=9223372036854775807;let unsigned:uint64=18446744073709551615;let small:int32=7;let f:float32=1.25;let dynamic:any=big;console.log(big,unsigned,small,f,dynamic,typeof big);try{small=big;}catch(e){console.log(e.name);}try{small=1.5 as any;}catch(e){console.log(e.name);}small=12 as any;console.log(small);let exact:int64=9007199254740993;exact++;let tiny:int32=1;let inferred=exact+tiny;console.log(exact,inferred);`
	for _, coerce := range []core.Tristate{core.TSTrue, core.TSFalse} {
		program := source
		rangeError := "RangeError"
		if coerce == core.TSFalse {
			// Widening remains implicit; narrowing requires an explicit method.
			program = strings.ReplaceAll(program, "small=12 as any", "small=(12 as any).int32()")
			rangeError = "TypeError"
		}
		text := emitGoOptions(t, program, core.CompilerOptions{Strict: core.TSTrue, CoerceAny: coerce})
		got, err := runGoProgram(t, text, true)
		if err != nil {
			t.Fatalf("%v\n%s", err, got)
		}
		if want := "9223372036854775807 18446744073709551615 7 1.25 9223372036854775807 number\n" + rangeError + "\n" + rangeError + "\n12\n9007199254740994 9007199254740995\n"; got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	}
}

func TestGoWorkerEventLoop(t *testing.T) {
	source := `declare function delay(ms:number):Promise<void>;let shared={value:0};function child(x:int32):int32{return x;}async function worker(x:int32){await delay(1);shared.value=x;return child(x);}async function main(){const result=await go worker(7);console.log(result,shared.value);try{await go fail();}catch(e){console.log("rejected");}}function fail(){throw "worker failed";}main();`
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	got, err := runGoProgram(t, text, true)
	if err != nil {
		t.Fatalf("%v\n%s", err, got)
	}
	if got != "7 7\nrejected\n" {
		t.Fatalf("got %q", got)
	}
}

func TestGoConcurrentWorkerImports(t *testing.T) {
	fs := bundled.WrapFS(vfstest.FromMap(map[string]string{
		"/src/main.ts":   `async function worker(){const m=await import("./helper.js");return m.value;}async function main(){const values=await Promise.all([go worker(),go worker()]);console.log(values[0],values[1]);}main();`,
		"/src/helper.ts": `console.log("initialized");export const value:int64=9223372036854775807;`,
	}, tspath.CaseSensitive))
	options := core.CompilerOptions{Target: core.ScriptTargetGo, OutDir: "/out"}
	p := compiler.NewProgram(compiler.ProgramOptions{Config: tsoptions.NewParsedCommandLine(&options, []tspath.RootedFilePath{"/src/main.ts"}, nil, "/src", fs.CaseSensitivity()), Host: compiler.NewCompilerHost(fs, bundled.LibPath(), nil, nil, nil)})
	var text string
	result := p.Emit(context.Background(), compiler.EmitOptions{WriteFile: func(path tspath.RootedFilePath, output string, data *compiler.WriteFileData) error {
		text = output
		return nil
	}})
	if result.EmitSkipped || len(result.Diagnostics) > 0 {
		for _, diag := range result.Diagnostics {
			t.Log(diag.String())
		}
		t.Fatal("worker import emission failed")
	}
	got, err := runGoProgram(t, text, true)
	if err != nil {
		t.Fatalf("%v\n%s", err, got)
	}
	if want := "initialized\n9223372036854775807 9223372036854775807\n"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestGoNativeNumericWidths(t *testing.T) {
	source := `type Short=int16;let a:int8=-128;let b:Short=32767;let c:uint8=255;let d:uint16=65535;let e:uint32=4294967295;let f:int=7;let g:uint=8;let h:float64=0.125;let rounded:float32=0.1;class Box{value:uint64=18446744073709551615;unset:int32;take(x:uint64):uint64{return x;}}const box=new Box();console.log(a,b,c,d,e,f,g,h,rounded as float64);console.log(box.value,box.unset,box.take(18446744073709551615));`
	text := emitGoOptions(t, source, core.CompilerOptions{Strict: core.TSTrue})
	got, err := runGoProgram(t, text, true)
	if err != nil {
		t.Fatalf("%v\n%s", err, got)
	}
	if want := "-128 32767 255 65535 4294967295 7 8 0.125 0.10000000149011612\n18446744073709551615 undefined 18446744073709551615\n"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
func TestGoNativeConstantRange(t *testing.T) {
	for _, source := range []string{`let x:int8=128;`, `let x:uint64=-1;`, `let x:int32=1.5;`, `let x:int32=2147483647+1;`, `let x:int64=9223372036854775808;`} {
		fs := bundled.WrapFS(vfstest.FromMap(map[string]string{"/src/input.ts": source}, tspath.CaseSensitive))
		options := core.CompilerOptions{Target: core.ScriptTargetGo, OutDir: "/out"}
		p := compiler.NewProgram(compiler.ProgramOptions{Config: tsoptions.NewParsedCommandLine(&options, []tspath.RootedFilePath{"/src/input.ts"}, nil, "/src", fs.CaseSensitivity()), Host: compiler.NewCompilerHost(fs, bundled.LibPath(), nil, nil, nil)})
		written := false
		result := p.Emit(context.Background(), compiler.EmitOptions{WriteFile: func(path tspath.RootedFilePath, text string, data *compiler.WriteFileData) error {
			written = true
			return nil
		}})
		if !result.EmitSkipped || len(result.Diagnostics) == 0 || written {
			t.Fatalf("invalid constant was accepted: %s", source)
		}
	}
}
