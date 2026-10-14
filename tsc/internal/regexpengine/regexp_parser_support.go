package regexpengine

import "unicode/utf16"

const (
	maxParserNesting   = 2048
	maxCaptureGroups   = 32767
	quantifierInfinity = int(1<<31 - 1)
)

type choiceBranch struct{ id, branch int }

func utf16String(units []uint16) string { return string(utf16.Decode(units)) }

func isIdentifierStart(character rune) bool {
	return character == '$' || character == '_' || unicodeRangeContains(regexpUnicodeIdentifierStartRanges, character)
}

func isIdentifierContinue(character rune) bool {
	return isIdentifierStart(character) || unicodeRangeContains(regexpUnicodeIdentifierContinueRanges, character) || character == 0x200c || character == 0x200d
}

func validCaptureName(name string) bool {
	for index, character := range name {
		if index == 0 {
			if !isIdentifierStart(character) {
				return false
			}
		} else if !isIdentifierContinue(character) {
			return false
		}
	}
	return name != ""
}
