package compiler_test

import (
	"strings"
	"testing"
)

func TestGoTimers(t *testing.T) {
	for _, tc := range []struct{ name, source, want string }{
		{"async ordering", `let x=0;setTimeout(()=>{x=10;console.log("timer",x);},0);Promise.resolve(0).then(()=>console.log("promise",x));console.log("sync",x);`, "sync 0\npromise 0\ntimer 10\n"},
		{"await timer", `async function f(){let x=await new Promise(resolve=>setTimeout(()=>resolve(10),0));console.log(x);}f();console.log("sync");`, "sync\n10\n"},
		{"timer callback arguments", `setTimeout((a,b)=>console.log(a,b),0,10,"value");`, "10 value\n"},
		{"timer reactions between tasks", `setTimeout(()=>{console.log("first");Promise.resolve(0).then(()=>console.log("reaction"));},0);setTimeout(()=>console.log("second"),0);`, "first\nreaction\nsecond\n"},
		{"canceled timer", `let handle=setTimeout(()=>console.log("bad"),60000);clearTimeout(handle);clearTimeout(handle);setTimeout(()=>console.log("done"));`, "done\n"},
		{"nested timer", `setTimeout(()=>{console.log("first");setTimeout(()=>console.log("nested"),0);},0);`, "first\nnested\n"},
		{"many immediate timers", `let count=0;for(let i=0;i<100;i++){setTimeout(()=>{count++;if(count===100)console.log(count);},0);}`, "100\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			output, err := runGoProgram(t, emitGoProgram(t, tc.source), true)
			if err != nil || output != tc.want {
				t.Fatalf("timer output %q, want %q, error %v", output, tc.want, err)
			}
		})
	}
}
func TestGoTimerErrors(t *testing.T) {
	source := `setTimeout(()=>{throw "timer failure";},0);setTimeout(()=>console.log("drained"),0);`
	output, err := runGoProgram(t, emitGoProgram(t, source), true)
	if err == nil || !strings.Contains(output, "timer failure") || !strings.Contains(output, "drained") || strings.Contains(output, "panic:") {
		t.Fatalf("timer failure was not safely reported: %q, %v", output, err)
	}
}
