//go:build !gojit || !linux || !amd64

package tsruntime

// TsCompileDynamicRegExp is available only in generated programs built with
// the custom Go fork on Linux/amd64 and the gojit build tag.
func init() {
	TsDynamicRegExpCompiler = TsCompileDynamicRegExp
}

func TsCompileDynamicRegExp(_ *TsLoop, _ *TsString, _ string) TsValue {
	TsPropertyFailure("Dynamic RegExp compilation requires the custom Go fork on linux/amd64 with -tags gojit")
	return TsU
}
