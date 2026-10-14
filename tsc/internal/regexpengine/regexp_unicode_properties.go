// Package regexpengine Unicode property parsing and membership semantics.
package regexpengine

import (
	"sort"
	"strings"
	"unicode"
)

func ecmaWhitespaceRanges() []classRange {
	return []classRange{{'\t', '\r'}, {' ', ' '}, {0xa0, 0xa0}, {0x1680, 0x1680}, {0x2000, 0x200a}, {0x2028, 0x2029}, {0x202f, 0x202f}, {0x205f, 0x205f}, {0x3000, 0x3000}, {0xfeff, 0xfeff}}
}
func canonicalUnicodeProperty(s string) (string, bool) {
	parts := strings.SplitN(s, "=", 2)
	if len(parts) == 1 {
		name := parts[0]
		if property, ok := regexpUnicodeStringPropertyAliases[name]; ok {
			return "string-property:" + property, true
		}
		if category, ok := regexpUnicodeGeneralCategoryAliases[name]; ok {
			return "gc:" + category, true
		}
		if name == "Any" || name == "ASCII" || name == "Assigned" {
			return "binary:" + name, true
		}
		if property, ok := regexpUnicodeBinaryAliases[name]; ok {
			return "binary:" + property, true
		}
		return "", false
	}
	key, ok := regexpUnicodePropertyAliases[parts[0]]
	if !ok {
		return "", false
	}
	value := parts[1]
	switch key {
	case "gc":
		if category, ok := regexpUnicodeGeneralCategoryAliases[value]; ok {
			return "gc:" + category, true
		}
	case "sc", "scx":
		if script, ok := regexpUnicodeScriptAliases[value]; ok {
			if key == "scx" {
				if code, ok := regexpUnicodeScriptCodeByName[script]; ok {
					return "script-extension:" + code + "|" + script, true
				}
				return "", false
			}
			return "script:" + script, true
		}
	}
	return "", false
}
func unicodePropertyContains(c rune, property string) bool {
	parts := strings.SplitN(property, ":", 2)
	if len(parts) != 2 {
		return false
	}
	switch parts[0] {
	case "gc":
		return unicodeRangeContains(regexpUnicodeGeneralCategoryRanges[parts[1]], c)
	case "script":
		return unicodeRangeContains(regexpUnicodeScriptRanges[parts[1]], c)
	case "script-extension":
		values := strings.SplitN(parts[1], "|", 2)
		if len(values) != 2 {
			return false
		}
		return unicodeScriptExtensionContains(c, values[0], values[1])
	case "binary":
		switch parts[1] {
		case "ASCII":
			return c <= 0x7f
		case "Any":
			return c >= 0 && c <= unicode.MaxRune
		case "Assigned":
			return !unicodeRangeContains(regexpUnicodeUnassignedRanges, c)
		}
		return unicodeRangeContains(unicodeBinaryPropertyRanges(parts[1]), c)
	}
	return false
}

// unicodeBinaryPropertyRanges reuses the separately named identifier ranges
// because capture-name scans consume them directly. Both copies come from the
// same pinned ICU 78.3 / Unicode 17.0 data used by the property parser.
func unicodeBinaryPropertyRanges(name string) []classRange {
	switch name {
	case "ID_Start":
		return regexpUnicodeIdentifierStartRanges
	case "ID_Continue":
		return regexpUnicodeIdentifierContinueRanges
	default:
		return regexpUnicodeBinaryPropertyRanges[name]
	}
}

func unicodeRangeContains(ranges []classRange, c rune) bool {
	index := sort.Search(len(ranges), func(i int) bool { return ranges[i].hi >= c })
	return index < len(ranges) && ranges[index].lo <= c
}

func unicodeScriptExtensionContains(c rune, code, script string) bool {
	index := sort.Search(len(regexpScriptExtensionOverrides), func(i int) bool {
		return regexpScriptExtensionOverrides[i].hi >= c
	})
	if index < len(regexpScriptExtensionOverrides) && regexpScriptExtensionOverrides[index].lo <= c {
		return strings.Contains(regexpScriptExtensionOverrides[index].values, "|"+code+"|")
	}
	return unicodeRangeContains(regexpUnicodeScriptRanges[script], c)
}
