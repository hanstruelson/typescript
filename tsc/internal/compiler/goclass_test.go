package compiler_test

import (
	"context"
	"os/exec"
	"strings"
	"testing"

	"github.com/microsoft/TypeScript/tsc/internal/bundled"
	"github.com/microsoft/TypeScript/tsc/internal/compiler"
	"github.com/microsoft/TypeScript/tsc/internal/core"
	"github.com/microsoft/TypeScript/tsc/internal/tsoptions"
	"github.com/microsoft/TypeScript/tsc/internal/tspath"
	"github.com/microsoft/TypeScript/tsc/internal/vfs/vfstest"
)

// Compare the executable Go output with the compiler's normal JavaScript output.
// This catches lexical-super mistakes that hand-written expected strings miss.
func javascriptOutput(t *testing.T, source string) string {
	return javascriptOutputTarget(t, source, core.ScriptTargetESNext)
}
func javascriptOutputTarget(t *testing.T, source string, target core.ScriptTarget) string {
	t.Helper()
	fs := bundled.WrapFS(vfstest.FromMap(map[string]string{"/src/input.ts": source}, tspath.CaseSensitive))
	options := core.CompilerOptions{Target: target, OutDir: "/out"}
	p := compiler.NewProgram(compiler.ProgramOptions{Config: tsoptions.NewParsedCommandLine(&options, []tspath.RootedFilePath{"/src/input.ts"}, nil, "/src", fs.CaseSensitivity()), Host: compiler.NewCompilerHost(fs, bundled.LibPath(), nil, nil, nil)})
	var js string
	result := p.Emit(context.Background(), compiler.EmitOptions{WriteFile: func(path tspath.RootedFilePath, text string, data *compiler.WriteFileData) error {
		js = text
		return nil
	}})
	if result.EmitSkipped || len(result.Diagnostics) > 0 {
		t.Fatalf("JavaScript emit failed: %v", result.Diagnostics)
	}
	output, err := exec.Command("node", "--input-type=module", "-e", js).CombinedOutput()
	if err != nil {
		t.Fatalf("JavaScript failed: %v\n%s", err, output)
	}
	return string(output)
}
func TestGoClasses(t *testing.T) {
	for _, tc := range []struct{ name, source string }{
		{"ancestry survives flattening", `
 class A {age:number=1;}class B extends A {}class C extends B {}class D extends C {}class Other {age:number=1;}
 const d=new D();console.log(d instanceof A,d instanceof B,d instanceof C,d instanceof D,d instanceof Other);console.log(new A() instanceof D);`},
		{"constructor errors are catchable", `
 class A {age:number=1;constructor(){console.log("base constructor");}}
 class Early extends A {constructor(){console.log(this.age);super();}}
 class Twice extends A {constructor(){super();super();}}
 class Missing extends A {constructor(){}}
 try{new Early();}catch(e){console.log("early caught");}
 try{new Twice();}catch(e){console.log("twice caught");}
 try{new Missing();}catch(e){console.log("missing caught");}`},
		{"base fields precede default arguments", `
 class A {age:number=7;constructor(n:number=this.age){console.log("default",n,this.age);this.age=n+1;}}
 class B extends A {age:number=this.age+10;}
 const b=new B();console.log(b.age);`},
		{"structured synchronous branches", `
 class A {age:number=1;f(n:number){let result=n>0?n:10;if(result>5){this.age+=result;}else{this.age=2;}const choose=n&&this.age;console.log(this.age,choose);return this.age;}}
 class B extends A {f(n:number){return super.f(n)+1;}}
 const b=new B();console.log(b.f(2));console.log(b.f(0));`},
		{"class method loop fallback", `
 class A {age:number=0;grow(n:number){for(let i=0;i<n;i++){this.age++;}return this.age;}}
 class B extends A {grow(n:number){return super.grow(n+1);}}
 const b=new B();console.log(b.grow(3));`},
		{"four level lexical super", `
 class A {age:number=1;speak(){console.log("A",this.age);} report(){this.speak();}}
 class B extends A {age:number=2;speak(){console.log("B",this.age);super.speak();}}
 class C extends B {age:number=3;speak(){console.log("C",this.age);super.speak();}}
 class D extends C {age:number=4;speak(){console.log("D",this.age);super.speak();}}
 const d=new D();d.speak();d.report();`},
		{"skipped ancestor and independent instances", `
 class A {age:number=1;speak(){console.log("A",this.age);}setAge(n:number){this.age=n;}}
 class B extends A {speak(){console.log("B",this.age);super.speak();}}
 class C extends B {}
 class D extends C {speak(){console.log("D",this.age);super.speak();}change(){super.setAge(10);}}
 const d=new D();const a=new A();d.change();d.speak();a.speak();`},
		{"constructor and initializer ordering", `
 class A {age:number=1;constructor(n:number){console.log("A before",this.age);this.age=n;this.report();}report(){console.log("A report",this.age);}}
 class B extends A {age:number=2;constructor(n:number){console.log("B before");super(n+1);console.log("B after",this.age);this.age+=10;}report(){console.log("B report",this.age);}}
 class C extends B {age:number=3;constructor(n:number){super(n+1);console.log("C after",this.age);}}
 class D extends C {}
 const d=new D(5);console.log("final",d.age);`},
		{"heterogeneous animal array", `
 class Animal {age:number=1;speak(){console.log("animal",this.age);}grow(){this.age++;this.speak();}}
 class Person extends Animal {age:number=2;speak(){console.log("person",this.age);}}
 class Dog extends Animal {age:number=3;speak(){console.log("dog",this.age);}}
 class Cat extends Animal {age:number=4;speak(){console.log("cat",this.age);}}
 const animals:Animal[]=[new Person(),new Dog(),new Cat()];
 animals[0].grow();animals[1].grow();animals[2].grow();console.log(animals[0].age,animals[1].age,animals[2].age);
 function show(animal:Animal){animal.speak();animal.age=10;animal.speak();}show(animals[0]);`},
		{"bracket and direct fields share storage", `
 class A {age:number=1;set(n:number){this.age=n;}get(){return this.age;}}
 class B extends A {age:number=2;set(n:number){super.set(n+1);}}
 const b=new B();let key="age";b[key]=10;console.log(b.age,b.get());b.set(20);console.log(b[key]);b.age=30;console.log(b[key]);b[key]+=2;console.log(b.age);b["extra"]=99;console.log(b["extra"],b["missing"]);b["set"](40);console.log(b.age);`},
		{"await and lexical arrow super", `
 class A {age:number=1;speak(){console.log("A",this.age);}}
 class B extends A {age:number=2;async speak(){await 0;const f=()=>super.speak();f();console.log("B",this.age);}}
 const b=new B();b.speak();console.log("sync");`},
		{"uninitialized numeric field", `
 class A {age!:number;read(){console.log(this.age,this["age"]);}set(){this.age=10;}}
 const a=new A();a.read();a.set();a.read();`},
		{"super keeps declaring class", `
 class A {age:number=1;f(){console.log("A.f",this.age);this.g();}g(){console.log("A.g",this.age);}}
 class B extends A {f(){console.log("B.f",this.age);super.f();}g(){console.log("B.g",this.age);super.g();}}
 class C extends B {f(){console.log("C.f",this.age);super.f();}}
 class D extends C {age:number=10;g(){console.log("D.g",this.age);super.g();}}
 const d=new D();d.f();`},
		{"initializer side effects and defaults", `
 class A {age:number=1;constructor(n:number=4){this.age=n;}}
 class B extends A {age:number=this.age+10;constructor(n:number=5){super(n);}}
 const b=new B();console.log(b.age);`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			text := emitGoProgram(t, tc.source)
			if !strings.Contains(text, "float64") || !strings.Contains(text, "struct {") {
				t.Fatal("missing concrete numeric struct")
			}
			got, err := runGoProgram(t, text, true)
			if err != nil {
				t.Fatalf("generated Go failed: %v\n%s", err, got)
			}
			want := javascriptOutput(t, tc.source)
			if got != want {
				t.Fatalf("Go output:\n%s\nJavaScript output:\n%s", got, want)
			}
		})
	}
}

func TestGoClassLayout(t *testing.T) {
	text := emitGoProgram(t, `class Animal {age:number=1;grow(){this.age++;}}class Person extends Animal {age:number=2;}const p=new Person();p.grow();console.log(p.age);`)
	// A redeclaration has one physical field per concrete class. Its initializer
	// runs separately; the inherited body accesses that same descendant field.
	if strings.Count(text, "P616765    float64") != 2 && strings.Count(text, "P616765 float64") != 2 {
		// gofmt pads according to the metadata fields; parse the declaration instead.
		count := 0
		for _, line := range strings.Split(text, "\n") {
			parts := strings.Fields(line)
			if len(parts) == 2 && parts[0] == "P616765" && parts[1] == "float64" {
				count++
			}
		}
		if count != 2 {
			t.Fatalf("expected two concrete structs with one age field each, got %d", count)
		}
	}
	start := strings.Index(text, "func (self *tsClass")
	end := strings.Index(text, "func main()")
	if start < 0 || end < start {
		t.Fatal("missing class implementations")
	}
	methods := text[start:end]
	if strings.Contains(methods, "&tsMachine") || strings.Contains(methods, "tsGet(") || strings.Contains(methods, "tsSet(") {
		t.Fatal("simple methods used a machine or dynamic property lookup")
	}
	if !strings.Contains(text, "View struct") || !strings.Contains(text, "Invalid concrete class layout") {
		t.Fatal("missing concrete class views and checked layout casts")
	}
}

func TestGoClassDiagnostics(t *testing.T) {
	for _, source := range []string{
		`class A {age:bigint=1n;}`,
		`class A {get age(){return 1;}}`,
		`class A {#age:number=1;}`,
		`class A {*f(){}}`,
		`function base(){return 0;}class A extends base(){}`,
		`class A {constructor(public age:number){}}`,
		`class A extends B {}class B extends A {}`,
	} {
		t.Run(source, func(t *testing.T) {
			fs := bundled.WrapFS(vfstest.FromMap(map[string]string{"/src/input.ts": source}, tspath.CaseSensitive))
			options := core.CompilerOptions{Target: core.ScriptTargetGo, OutDir: "/out"}
			p := compiler.NewProgram(compiler.ProgramOptions{Config: tsoptions.NewParsedCommandLine(&options, []tspath.RootedFilePath{"/src/input.ts"}, nil, "/src", fs.CaseSensitivity()), Host: compiler.NewCompilerHost(fs, bundled.LibPath(), nil, nil, nil)})
			written := false
			result := p.Emit(context.Background(), compiler.EmitOptions{WriteFile: func(path tspath.RootedFilePath, text string, data *compiler.WriteFileData) error {
				written = true
				return nil
			}})
			if !result.EmitSkipped || len(result.Diagnostics) == 0 || written {
				t.Fatalf("expected diagnostic and no output: %+v", result)
			}
		})
	}
}

func TestGoClassModuleBundle(t *testing.T) {
	fs := bundled.WrapFS(vfstest.FromMap(map[string]string{
		"/src/main.ts":   `import {Person} from "./helper.js";const p=new Person();p.speak();p["age"]=10;p.speak();console.log(Person.total);Person.bump();console.log(Person.total);`,
		"/src/helper.ts": `class Animal {static total:number=1;static bump(){this.total++;}age:number=1;speak(){console.log("animal",this.age);}}export class Person extends Animal {age:number=2;speak(){console.log("person",this.age);super.speak();}}`,
	}, tspath.CaseSensitive))
	options := core.CompilerOptions{Target: core.ScriptTargetGo, OutDir: "/out"}
	p := compiler.NewProgram(compiler.ProgramOptions{Config: tsoptions.NewParsedCommandLine(&options, []tspath.RootedFilePath{"/src/main.ts"}, nil, "/src", fs.CaseSensitivity()), Host: compiler.NewCompilerHost(fs, bundled.LibPath(), nil, nil, nil)})
	var text string
	written := 0
	result := p.Emit(context.Background(), compiler.EmitOptions{WriteFile: func(path tspath.RootedFilePath, output string, data *compiler.WriteFileData) error {
		text = output
		written++
		return nil
	}})
	if result.EmitSkipped || len(result.Diagnostics) > 0 || written != 1 {
		for _, diag := range result.Diagnostics {
			t.Log(diag.String())
		}
		t.Fatalf("class module emission failed: %+v", result)
	}
	got, err := runGoProgram(t, text, true)
	if err != nil {
		t.Fatalf("class bundle failed: %v\n%s", err, got)
	}
	if want := "person 2\nanimal 2\nperson 10\nanimal 10\n1\n2\n"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestGoStaticClasses(t *testing.T) {
	for _, source := range []string{
		`class A {static value:number=1;}class B extends A {static set(){super.value=8;console.log(super.value,this.value);}}B.set();console.log(A.value,B.value);A["extra"]=10;console.log(B["extra"]);B["extra"]=20;console.log(A["extra"],B["extra"]);`,
		`class A {static x:number=1;static {let value=2;this.x+=value;console.log("block",this.x);}static y:number=this.x+1;}class B extends A {static {for(let i=0;i<2;i++){this.x++;}console.log("child",this.x,super.x);}}console.log(A.x,A.y,B.x,B.y,Object.keys(A).length,Object.hasOwn(B,"y"));`,
		`class Counter {static count:number=1;static add(n:number){this.count+=n;return this.count;} value:number=10;}console.log(Counter.count,Counter.add(2),Counter.count);const a=new Counter();const b=new Counter();console.log(a.value,b.value,Counter["count"]);Counter["count"]=8;console.log(Counter.add(1));`,
		`class A {static value:number=1;static who(){console.log("A",this.value,this===B);}static grow(){this.value++;}}class B extends A {static who(){console.log("B",this.value);super.who();}}class C extends B {static who(){console.log("C",this.value);super.who();}}A.value=4;console.log(B.value,C.value);B.grow();console.log(A.value,B.value,C.value);C.who();console.log(typeof A);`,
		`class A {static value:number=1;static speak(){console.log("A",this.value);}}class B extends A {static value:number=2;static speak(){console.log("B",this.value,super.value);super.speak();}}class C extends B {static value:number=3;static speak(){console.log("C",this.value,super.value);super.speak();}}class D extends C {static value:number=4;static speak(){console.log("D",this.value,super.value);super.speak();}}D.speak();A.speak();`,
		`let count=0;class A {static first:number=(count+=1);static second:number=A.first+1;static owner:any=this;static name:string="😀";static async run(){await 0;console.log(this.name,this===A);}}console.log(count,A.first,A.second,A.owner===A);new A();new A();console.log(count);A.run();console.log("sync");`,
	} {
		t.Run(source, func(t *testing.T) {
			text := emitGoProgram(t, source)
			got, err := runGoProgram(t, text, true)
			if err != nil {
				t.Fatalf("%v\n%s", err, got)
			}
			want := javascriptOutput(t, source)
			if got != want {
				t.Fatalf("Go:\n%s\nJavaScript:\n%s", got, want)
			}
		})
	}
}
