package goemit

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf16"

	"github.com/microsoft/TypeScript/tsc/internal/ast"
	"github.com/microsoft/TypeScript/tsc/internal/regexpengine"
	"github.com/microsoft/TypeScript/tsc/internal/stringutil"
)

// compiledRegExp lowers one statically known source/flags pair through the
// V8-shaped parser and graph compiler, then embeds its direct Go matcher.
func (e *emitter) compiledRegExp(pattern, flags string, node *ast.Node) string {
	if e.regexMatchers == nil {
		e.regexMatchers = make(map[string]string)
	}
	key := pattern + "\x00" + flags
	matcher := e.regexMatchers[key]
	if matcher == "" {
		matcher = e.unique("RegexMatcher")
		code, err := regexpengine.CompileUTF16ToGo(regExpUTF16(pattern), flags, matcher)
		if err != nil {
			e.fail(node, "regular expression cannot be emitted as Go code: "+err.Error())
			return "tsU"
		}
		e.regexMatchers[key] = matcher
		e.regexText.WriteString(code)
		e.regexText.WriteByte('\n')
	}
	return fmt.Sprintf("tsNewCompiledRegExp(%s,%s,%s,%sCaptureSlots,%s)", stringLiteral(pattern), stringLiteral(flags), matcher+"Into", matcher, namedCaptureLiteral(pattern, flags))
}

func namedCaptureLiteral(pattern, flags string) string {
	captures, err := regexpengine.CaptureNamesOrderedUTF16(regExpUTF16(pattern), flags)
	if err != nil || len(captures) == 0 {
		return "nil"
	}
	entries := make([]string, 0, len(captures))
	for _, captureGroup := range captures {
		values := make([]string, len(captureGroup.Indices))
		for index, capture := range captureGroup.Indices {
			values[index] = strconv.Itoa(capture)
		}
		entries = append(entries, "{name:"+strconv.Quote(captureGroup.Name)+",indices:[]int{"+strings.Join(values, ",")+"}}")
	}
	return "[]tsRegExpNamedCapture{" + strings.Join(entries, ",") + "}"
}

func regExpUTF16(text string) []uint16 {
	units := make([]uint16, 0, len(text))
	for len(text) > 0 {
		character, size := stringutil.DecodeJSStringRune(text)
		text = text[size:]
		if character > 0xffff {
			leading, trailing := utf16.EncodeRune(character)
			units = append(units, uint16(leading), uint16(trailing))
		} else {
			units = append(units, uint16(character))
		}
	}
	return units
}

func (b *machineBuilder) literalRegExpConstructor(node *ast.Node, arguments []*ast.Node) string {
	pattern, flags, known := "", "", len(arguments) <= 2
	if len(arguments) > 0 {
		if arguments[0].Kind == ast.KindStringLiteral || arguments[0].Kind == ast.KindNoSubstitutionTemplateLiteral {
			pattern = arguments[0].Text()
		} else {
			known = false
		}
	}
	if len(arguments) > 1 {
		if arguments[1].Kind == ast.KindStringLiteral || arguments[1].Kind == ast.KindNoSubstitutionTemplateLiteral {
			flags = arguments[1].Text()
		} else {
			known = false
		}
	}
	spread := false
	for _, argument := range arguments {
		if argument.Kind == ast.KindSpreadElement {
			spread = true
		}
	}
	if known && !spread {
		return b.e.compiledRegExp(pattern, flags, node)
	}
	if spread {
		return "tsDynamicRegExp(loop," + b.spreadArguments(arguments) + "...)"
	}
	values := make([]string, len(arguments))
	for index, argument := range arguments {
		values[index] = b.expression(argument)
	}
	return "tsDynamicRegExp(loop," + strings.Join(values, ",") + ")"
}
