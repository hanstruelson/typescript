package compiler_test

import (
	"github.com/microsoft/TypeScript/tsc/internal/core"
	"testing"
)

func TestGoClassDecorators(t *testing.T) {
	cases := []struct{ name, source string }{
		{"application and initializer order", `function decorate(label:string){console.log("evaluate",label);return (value:any,context:any)=>{console.log("apply",label,context.kind,context.name,value.field);context.addInitializer(()=>console.log("initialize",label));return value;};}function field(){console.log("static field");return 3;}@decorate("outer") @decorate("inner") class Decorated{static field=field();value=7;read(){return this.value;}}console.log(Decorated.field,new Decorated().read());`},
		{"replacement", `class Replacement{value=9;read(){return this.value;}}function replace(value:any,context:any){return Replacement;}@replace class Original{value=1;read(){return this.value;}}const x=new Original();console.log(x.value,x.read(),x instanceof Original,x instanceof Replacement);`},
		{"late initializer rejected", `let later:any;function dec(value:any,context:any){later=context.addInitializer;}@dec class Decorated{}try{later(()=>console.log("bad"));}catch(e){console.log("late rejected");}console.log(new Decorated() instanceof Decorated);`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			want := javascriptOutputTarget(t, tc.source, core.ScriptTargetES2022)
			text := emitGoOptions(t, tc.source, core.CompilerOptions{Strict: core.TSTrue})
			got, err := runGoProgram(t, text, true)
			if err != nil || got != want {
				t.Fatalf("got %q, want %q: %v", got, want, err)
			}
		})
	}
}
func TestGoLegacyClassDecorator(t *testing.T) {
	text := emitGoOptions(t, `function decorate(value:any){value.extra=7;}@decorate class Legacy{value=3;}console.log((Legacy as any).extra,new Legacy().value);`, core.CompilerOptions{Strict: core.TSTrue, ExperimentalDecorators: core.TSTrue})
	got, err := runGoProgram(t, text, true)
	if err != nil || got != "7 3\n" {
		t.Fatalf("%q: %v", got, err)
	}
}
