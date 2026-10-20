//go:build gojit && linux && amd64

package tsruntime

import (
	"strconv"
	"sync/atomic"

	"github.com/microsoft/TypeScript/tsc/internal/regexpengine"
	"jit"
)

var tsJITRegexRequestID atomic.Uint64

func init() {
	TsDynamicRegExpCompiler = TsCompileDynamicRegExp
}

// TsCompileDynamicRegExp parses and compiles a runtime pattern inside the process,
// then installs its generated matcher on the new RegExp. The Code handle is
// retained with the RegExp and deliberately never closed automatically: this
// bridge predates a proven GC owner/retirement protocol, so each successful
// dynamic compilation leaks its executable mapping for process lifetime.
func TsCompileDynamicRegExp(loop *TsLoop, source *TsString, flags string) TsValue {
	expression, captureSlots, err := regexpengine.CompileUTF16ToGoMatcherLiteral(source.Units, flags)
	if err != nil {
		TsRegExpSyntaxFailure("Invalid regular expression: " + err.Error())
		return TsU
	}

	id := tsJITRegexRequestID.Add(1)
	if id == 0 {
		TsPropertyFailure("Dynamic RegExp compilation ID space exhausted")
		return TsU
	}
	defer regexpengine.DeleteJITMatcher(id)

	snippet := "regexpengine.StoreJITMatcher(" + strconv.FormatUint(id, 10) + ", " + expression + ")"
	code, err := jit.Eval(snippet)
	if err != nil {
		if code != nil {
			_ = code.Close()
		}
		TsPropertyFailure("Dynamic RegExp JIT compilation failed: " + err.Error())
		return TsU
	}
	matcher, ok := regexpengine.LoadJITMatcher(id)
	if !ok {
		_ = code.Close()
		TsPropertyFailure("Dynamic RegExp JIT did not install its matcher")
		return TsU
	}

	names, err := regexpengine.CaptureNamesOrderedUTF16(source.Units, flags)
	if err != nil {
		_ = code.Close()
		TsRegExpSyntaxFailure("Invalid regular expression: " + err.Error())
		return TsU
	}
	namedCaptures := make([]TsRegExpNamedCapture, 0, len(names))
	for _, named := range names {
		namedCaptures = append(namedCaptures, TsRegExpNamedCapture{
			Name:    named.Name,
			Indices: append([]int(nil), named.Indices...),
		})
	}
	value := TsNewCompiledRegExp(loop, source, TsStringUTF8(flags), matcher, captureSlots, namedCaptures)
	TsRegExpFromValue(value).JitCode = code
	return value
}
