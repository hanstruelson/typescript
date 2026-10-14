// Package regexpengine flag definitions and parsing, corresponding to V8 regexp-flags.h.
package regexpengine

import (
	"fmt"
	"strings"
)

type Flag uint16

const (
	// These bit positions match REGEXP_FLAG_LIST in V8's regexp-flags.h.
	// The declarations remain in flag-character order below, while the bit
	// positions are intentionally shuffled to match V8.
	FlagHasIndices  Flag = 1 << 7
	FlagGlobal      Flag = 1 << 0
	FlagIgnoreCase  Flag = 1 << 1
	FlagLinear      Flag = 1 << 6
	FlagMultiline   Flag = 1 << 2
	FlagDotAll      Flag = 1 << 5
	FlagUnicode     Flag = 1 << 4
	FlagUnicodeSets Flag = 1 << 8
	FlagSticky      Flag = 1 << 3
)

const (
	regexpFlagCount = 9
	regexpFlagMask  = FlagHasIndices | FlagGlobal | FlagIgnoreCase | FlagLinear |
		FlagMultiline | FlagDotAll | FlagUnicode | FlagUnicodeSets | FlagSticky
	regexpFlagOrder = "dgilmsuvy"
	regexpFlagRange = (1 << regexpFlagCount) - 1
)

// Keep the V8 flag bits contiguous even though their alphabetic declaration
// order differs from their bit positions.
var _ [int(regexpFlagMask) - regexpFlagRange]struct{}
var _ [regexpFlagRange - int(regexpFlagMask)]struct{}

type Flags struct {
	Global, HasIndices, IgnoreCase, Linear, Multiline, DotAll bool
	Unicode, UnicodeSets, Sticky                              bool
}

func (f Flags) Has(flag Flag) bool {
	switch flag {
	case FlagGlobal:
		return f.Global
	case FlagHasIndices:
		return f.HasIndices
	case FlagIgnoreCase:
		return f.IgnoreCase
	case FlagLinear:
		return f.Linear
	case FlagMultiline:
		return f.Multiline
	case FlagDotAll:
		return f.DotAll
	case FlagUnicode:
		return f.Unicode
	case FlagUnicodeSets:
		return f.UnicodeSets
	case FlagSticky:
		return f.Sticky
	default:
		return false
	}
}

func (f *Flags) Set(flag Flag, value bool) {
	switch flag {
	case FlagGlobal:
		f.Global = value
	case FlagHasIndices:
		f.HasIndices = value
	case FlagIgnoreCase:
		f.IgnoreCase = value
	case FlagLinear:
		f.Linear = value
	case FlagMultiline:
		f.Multiline = value
	case FlagDotAll:
		f.DotAll = value
	case FlagUnicode:
		f.Unicode = value
	case FlagUnicodeSets:
		f.UnicodeSets = value
	case FlagSticky:
		f.Sticky = value
	}
}

func IsHasIndices(f Flags) bool    { return f.HasIndices }
func IsGlobal(f Flags) bool        { return f.Global }
func IsIgnoreCase(f Flags) bool    { return f.IgnoreCase }
func IsLinear(f Flags) bool        { return f.Linear }
func IsMultiline(f Flags) bool     { return f.Multiline }
func IsDotAll(f Flags) bool        { return f.DotAll }
func IsUnicode(f Flags) bool       { return f.Unicode }
func IsUnicodeSets(f Flags) bool   { return f.UnicodeSets }
func IsSticky(f Flags) bool        { return f.Sticky }
func IsEitherUnicode(f Flags) bool { return f.Unicode || f.UnicodeSets }
func ShouldOptionallyStepBackToLeadSurrogate(f Flags) bool {
	return IsEitherUnicode(f) && (f.Global || f.Sticky)
}

func TryFlagFromChar(c rune) (Flag, bool) {
	switch c {
	case 'd':
		return FlagHasIndices, true
	case 'g':
		return FlagGlobal, true
	case 'i':
		return FlagIgnoreCase, true
	case 'l':
		return FlagLinear, true
	case 'm':
		return FlagMultiline, true
	case 's':
		return FlagDotAll, true
	case 'u':
		return FlagUnicode, true
	case 'v':
		return FlagUnicodeSets, true
	case 'y':
		return FlagSticky, true
	default:
		return 0, false
	}
}

func (f Flags) String() string {
	var b strings.Builder
	for _, c := range regexpFlagOrder {
		flag, _ := TryFlagFromChar(c)
		if f.Has(flag) {
			b.WriteRune(c)
		}
	}
	return b.String()
}

func ParseFlags(text string) (Flags, error) {
	var f, seen Flags
	for position, c := range text {
		flag, ok := TryFlagFromChar(c)
		if !ok {
			return Flags{}, fmt.Errorf("regular expression at %d: invalid flag %q", position, c)
		}
		if seen.Has(flag) {
			return Flags{}, fmt.Errorf("regular expression at %d: duplicate flag %q", position, c)
		}
		seen.Set(flag, true)
		f.Set(flag, true)
	}
	if strings.ContainsRune(text, 'u') && strings.ContainsRune(text, 'v') {
		return Flags{}, fmt.Errorf("regular expression at 0: flags u and v cannot be used together")
	}
	return f, nil
}
