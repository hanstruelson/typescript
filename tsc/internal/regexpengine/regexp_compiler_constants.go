package regexpengine

// V8 regexp compiler range tables use half-open endpoints, with a terminal
// marker one past the largest Unicode code point.
const v8CompilerRangeEndMarker = 0x110000

// v8CompilerMaxCodeUnit and v8CompilerCharMask port the compiler's code-unit
// width helpers. Character masks are all low bits because both supported
// widths are powers of two.
func v8CompilerMaxCodeUnit(oneByte bool) rune {
	if oneByte {
		return 0xff
	}
	return 0xffff
}

func v8CompilerCharMask(oneByte bool) uint32 {
	return uint32(v8CompilerMaxCodeUnit(oneByte))
}

var v8CompilerSpaceRanges = [...]int{
	'\t', '\r' + 1, ' ', ' ' + 1, 0x00a0, 0x00a1, 0x1680, 0x1681,
	0x2000, 0x200b, 0x2028, 0x202a, 0x202f, 0x2030, 0x205f, 0x2060,
	0x3000, 0x3001, 0xfeff, 0xff00, v8CompilerRangeEndMarker,
}

var v8CompilerWordRanges = [...]int{
	'0', '9' + 1, 'A', 'Z' + 1, '_', '_' + 1, 'a', 'z' + 1,
	v8CompilerRangeEndMarker,
}

var v8CompilerDigitRanges = [...]int{'0', '9' + 1, v8CompilerRangeEndMarker}

var v8CompilerSurrogateRanges = [...]int{0xd800, 0xd801, v8CompilerRangeEndMarker}

var v8CompilerLineTerminatorRanges = [...]int{
	0x000a, 0x000b, 0x000d, 0x000e, 0x2028, 0x202a,
	v8CompilerRangeEndMarker,
}

const (
	v8CompilerSpaceRangeCount              = len(v8CompilerSpaceRanges)
	v8CompilerWordRangeCount               = len(v8CompilerWordRanges)
	v8CompilerDigitRangeCount              = len(v8CompilerDigitRanges)
	v8CompilerSurrogateRangeCount          = len(v8CompilerSurrogateRanges)
	v8CompilerLineTerminatorRangeCount     = len(v8CompilerLineTerminatorRanges)
	v8CompilerMaxLookaheadForBoyerMoore    = 8
	v8CompilerPatternTooShortForBoyerMoore = 2
)
