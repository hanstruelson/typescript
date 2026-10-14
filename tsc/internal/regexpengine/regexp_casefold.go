// Copyright 2019 the V8 project authors. All rights reserved.
// This file adapts the ECMAScript/V8 regexp case-folding rules under LICENSE.v8.
package regexpengine

import "sort"

// v8ASTCaseFoldingEquivalenceKey mirrors CaseFolding::EquivalenceKey. The
// generated simple-fold target is a stable representative of each ICU fold
// class; IgnoreSet members remain distinct in non-Unicode mode.
func v8ASTCaseFoldingEquivalenceKey(character rune, mode v8TonodeCaseFoldingMode) rune {
	switch mode {
	case v8TonodeNonUnicodeCaseFolding:
		if character < 0 || character > 0xffff {
			panic("V8 non-Unicode case-fold key requires a UTF-16 code unit")
		}
		if regexpLegacyIgnoreSet[character] {
			return character
		}
	case v8TonodeUnicodeCaseFolding:
		if character < 0 || character > 0x10ffff {
			panic("V8 Unicode case-fold key requires a Unicode code point")
		}
	default:
		panic("unknown V8 case-folding mode")
	}
	return regexpSimpleCaseFold(character)
}

// v8ASTCaseFoldingCloseOver ports CaseFolding::CloseOver for the code-point
// ranges used by the regexp frontend. It preserves input ranges and adds each
// eligible simple-fold class. Non-Unicode mode leaves IgnoreSet members only
// when they were already present in the input.
func v8ASTCaseFoldingCloseOver(ranges []v8ASTCharacterRange, mode v8TonodeCaseFoldingMode) []v8ASTCharacterRange {
	original := v8ASTCanonicalizeRanges(ranges)
	if mode == v8TonodeNonUnicodeCaseFolding {
		for _, current := range original {
			if current.to > 0xffff {
				panic("V8 non-Unicode case closure received supplementary characters")
			}
		}
	} else if mode != v8TonodeUnicodeCaseFolding {
		panic("unknown V8 case-folding mode")
	}

	result := append([]v8ASTCharacterRange(nil), original...)
	contains := func(character rune) bool {
		index := sort.Search(len(original), func(index int) bool {
			return original[index].to >= character
		})
		return index < len(original) && original[index].from <= character
	}
	for _, class := range regexpCaseFoldClasses {
		seeded := false
		for _, character := range class {
			if mode == v8TonodeNonUnicodeCaseFolding &&
				(character > 0xffff || regexpLegacyIgnoreSet[character]) {
				continue
			}
			if contains(character) {
				seeded = true
				break
			}
		}
		if !seeded {
			continue
		}
		for _, character := range class {
			if mode == v8TonodeNonUnicodeCaseFolding &&
				(character > 0xffff || regexpLegacyIgnoreSet[character]) {
				continue
			}
			result = append(result, v8ASTRangeSingleton(character))
		}
	}
	return v8ASTCanonicalizeRanges(result)
}

func v8ASTCaseFoldingHasSupplementaryEquivalents(character uint16) bool {
	return regexpSupplementaryEquivalentsSet[character]
}

func v8ASTCaseFoldingHasOtherUnicodeEquivalent(character rune) bool {
	index, ok := regexpCaseFoldClassByRune[character]
	return ok && len(regexpCaseFoldClasses[index]) > 1
}

// forEachUnicodeCaseEquivalent visits the pinned simple-fold class for c.
func forEachUnicodeCaseEquivalent(c rune, visit func(rune) bool) bool {
	return forEachSimpleCaseEquivalent(c, visit)
}

// forEachLegacyCaseEquivalent visits only characters in c's ECMAScript
// non-Unicode Canonicalize class. Pinned fold data supplies candidates;
// Canonicalize excludes Unicode-only pairs such as ASCII s and long s.
func forEachLegacyCaseEquivalent(c rune, visit func(rune) bool) bool {
	key := canonicalizeLegacy(c)
	return forEachSimpleCaseEquivalent(c, func(folded rune) bool {
		return canonicalizeLegacy(folded) == key && visit(folded)
	})
}

func forEachSimpleCaseEquivalent(c rune, visit func(rune) bool) bool {
	if visit(c) {
		return true
	}
	if index, ok := regexpCaseFoldClassByRune[c]; ok {
		for _, equivalent := range regexpCaseFoldClasses[index] {
			if equivalent != c && visit(equivalent) {
				return true
			}
		}
	}
	return false
}

// regexpSimpleCaseFold applies the pinned one-code-point Unicode default fold
// used by V8's Unicode-set string alternatives. Missing entries fold to self.
func regexpSimpleCaseFold(c rune) rune {
	if folded, ok := regexpSimpleCaseFoldTo[c]; ok {
		return folded
	}
	return c
}

// canonicalizeLegacy implements ECMAScript's non-Unicode Canonicalize rule.
// It is used by the V8-derived AST rewrite before the Go source backend.
func canonicalizeLegacy(c rune) rune {
	upper := c
	if c < 0x80 {
		if c >= 'a' && c <= 'z' {
			upper = c - ('a' - 'A')
		}
	} else if mapped, ok := regexpLegacyUppercase[c]; ok {
		upper = mapped
	}
	if c >= 0x80 && upper < 0x80 {
		return c
	}
	return upper
}
