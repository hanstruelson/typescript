// Copyright 2012 the V8 project authors. All rights reserved.
// This file adapts V8's RegExp source escaping behavior under LICENSE.v8.
package regexpengine

// EscapedSourceUTF16 returns the representation used when exposing a compiled
// pattern's source. It escapes line terminators and unescaped slashes outside
// character classes while preserving UTF-16 code units, including lone
// surrogates. The empty source is represented as (?:).
func (p *Program) EscapedSourceUTF16() []uint16 {
	if len(p.source) == 0 {
		return []uint16{'(', '?', ':', ')'}
	}
	result := make([]uint16, 0, len(p.source))
	inClass := false
	for i := 0; i < len(p.source); {
		c := p.source[i]
		if c == '\\' {
			if i+1 < len(p.source) && isUTF16LineTerminator(p.source[i+1]) {
				// The source backslash is ignored; escape the line terminator on
				// the following iteration.
				i++
				continue
			}
			result = append(result, c)
			i++
			if i < len(p.source) {
				result = append(result, p.source[i])
				i++
			}
			continue
		}
		if c == '/' && !inClass {
			result = append(result, '\\')
		}
		switch c {
		case '[':
			inClass = true
		case ']':
			inClass = false
		case '\n':
			result = append(result, '\\', 'n')
			i++
			continue
		case '\r':
			result = append(result, '\\', 'r')
			i++
			continue
		case 0x2028:
			result = append(result, '\\', 'u', '2', '0', '2', '8')
			i++
			continue
		case 0x2029:
			result = append(result, '\\', 'u', '2', '0', '2', '9')
			i++
			continue
		}
		result = append(result, c)
		i++
	}
	return result
}

// RegExpStringUTF16 returns a slash-delimited display form of the compiled
// source using this program's canonical flag order.
func (p *Program) RegExpStringUTF16() []uint16 {
	source := p.EscapedSourceUTF16()
	flags := p.Flags.String()
	result := make([]uint16, 0, len(source)+len(flags)+2)
	result = append(result, '/')
	result = append(result, source...)
	result = append(result, '/')
	for _, flag := range flags {
		result = append(result, uint16(flag))
	}
	return result
}

func isUTF16LineTerminator(c uint16) bool {
	return c == '\n' || c == '\r' || c == 0x2028 || c == 0x2029
}
