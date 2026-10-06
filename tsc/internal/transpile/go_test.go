package transpile

import (
	"context"
	"strings"
	"testing"

	"github.com/microsoft/TypeScript/tsc/internal/core"
)

func TestTranspileGoTarget(t *testing.T) {
	options := Options{CompilerOptions: &core.CompilerOptions{Target: core.ScriptTargetGo}, ReportDiagnostics: true}
	output := TranspileModule(context.Background(), `async function f(){let x=await Promise.resolve(10);console.log(x);}f();`, options)
	if output == nil || len(output.Diagnostics) != 0 || !strings.Contains(output.OutputText, "m.await(") {
		t.Fatalf("Go transpilation failed: %+v", output)
	}
	output = TranspileModule(context.Background(), `class Unsupported { get x(){return 1;} }`, options)
	if output == nil || len(output.Diagnostics) == 0 || output.OutputText != "" {
		t.Fatalf("unsupported Go transpilation should return diagnostics: %+v", output)
	}
}
