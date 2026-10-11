// Copyright 2019 the V8 project authors. All rights reserved.
// This file adapts the ECMAScript/V8 regexp case-folding rules under LICENSE.v8.
package regexpengine

import "unicode"

// forEachUnicodeCaseEquivalent visits the simple-fold class for c. Go's
// Unicode tables currently predate the supplementary case mapping called out
// by V8's generator, so preserve that mapping explicitly as well.
func forEachUnicodeCaseEquivalent(c rune, visit func(rune) bool) bool {
	if visit(c) {
		return true
	}
	for folded := unicode.SimpleFold(c); folded != c; folded = unicode.SimpleFold(folded) {
		if visit(folded) {
			return true
		}
	}
	if c == 0x00df || c == 0x1e9e || c == 0x1df95 {
		for _, equivalent := range [...]rune{0x00df, 0x1e9e, 0x1df95} {
			if equivalent != c && visit(equivalent) {
				return true
			}
		}
	}
	return false
}

// forEachLegacyCaseEquivalent visits only characters in c's ECMAScript
// non-Unicode Canonicalize class. SimpleFold supplies candidate equivalents;
// Canonicalize excludes Unicode-only pairs such as ASCII s and long s.
func forEachLegacyCaseEquivalent(c rune, visit func(rune) bool) bool {
	key := canonicalizeLegacy(c)
	if visit(c) {
		return true
	}
	for folded := unicode.SimpleFold(c); folded != c; folded = unicode.SimpleFold(folded) {
		if canonicalizeLegacy(folded) == key && visit(folded) {
			return true
		}
	}
	return false
}
