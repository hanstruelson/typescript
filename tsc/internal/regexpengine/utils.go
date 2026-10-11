package regexpengine

// AdvanceStringIndexUTF16 implements ECMAScript AdvanceStringIndex over UTF-16
// code units. Callers supply a non-negative index, normally obtained through
// ToLength on RegExp.lastIndex.
func AdvanceStringIndexUTF16(input []uint16, index int, unicodeMode bool) int {
	return advanceStringIndex(input, index, unicodeMode)
}

// IsMatchedCapture reports whether capture is present in match. Capture zero
// is the full match; negative and out-of-range capture numbers are unmatched.
func IsMatchedCapture(match Match, capture int) bool {
	if capture < 0 || capture >= len(match.Captures) {
		return false
	}
	span := match.Captures[capture]
	return span.Start >= 0 && span.End >= 0
}

// CaptureUTF16 returns a copy of a matched capture's UTF-16 code units. A
// matched empty capture returns an empty slice and true; an unmatched capture
// or invalid span returns nil and false. Copying prevents callers from
// mutating the subject through the returned value.
func CaptureUTF16(input []uint16, match Match, capture int) ([]uint16, bool) {
	view, ok := match.CaptureUTF16(input, capture)
	if !ok {
		return nil, false
	}
	return append([]uint16(nil), view...), true
}
