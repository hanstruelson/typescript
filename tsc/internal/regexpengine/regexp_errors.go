// Copyright 2020 the V8 project authors. All rights reserved.
// This file adapts V8's regexp-error message catalog under LICENSE.v8.
package regexpengine

// regexpError follows the order of V8's regexp::Error enum. Keeping the enum
// and message table together makes diagnostic mappings stable as the parser is
// brought closer to V8's grammar.
type regexpError uint8

const (
	errorNone regexpError = iota
	errorStackOverflow
	errorAnalysisStackOverflow
	errorTooLarge
	errorUnterminatedGroup
	errorUnmatchedParen
	errorEscapeAtEndOfPattern
	errorInvalidPropertyName
	errorInvalidEscape
	errorInvalidDecimalEscape
	errorInvalidUnicodeEscape
	errorNothingToRepeat
	errorLoneQuantifierBrackets
	errorRangeOutOfOrder
	errorIncompleteQuantifier
	errorInvalidQuantifier
	errorInvalidGroup
	errorMultipleFlagDashes
	errorNotLinear
	errorRepeatedFlag
	errorInvalidFlagGroup
	errorTooManyCaptures
	errorInvalidCaptureGroupName
	errorDuplicateCaptureGroupName
	errorInvalidNamedReference
	errorInvalidNamedCaptureReference
	errorInvalidClassPropertyName
	errorInvalidCharacterClass
	errorUnterminatedCharacterClass
	errorOutOfOrderCharacterClass
	errorInvalidClassSetOperation
	errorInvalidCharacterInClass
	errorNegatedCharacterClassWithStrings
	errorUnsupportedBytecode
	errorNumErrors
)

var regexpErrorMessages = [...]string{
	errorNone:                             "",
	errorStackOverflow:                    "Maximum call stack size exceeded",
	errorAnalysisStackOverflow:            "Stack overflow",
	errorTooLarge:                         "Regular expression too large",
	errorUnterminatedGroup:                "Unterminated group",
	errorUnmatchedParen:                   "Unmatched ')'",
	errorEscapeAtEndOfPattern:             `\ at end of pattern`,
	errorInvalidPropertyName:              "Invalid property name",
	errorInvalidEscape:                    "Invalid escape",
	errorInvalidDecimalEscape:             "Invalid decimal escape",
	errorInvalidUnicodeEscape:             "Invalid Unicode escape",
	errorNothingToRepeat:                  "Nothing to repeat",
	errorLoneQuantifierBrackets:           "Lone quantifier brackets",
	errorRangeOutOfOrder:                  "numbers out of order in {} quantifier",
	errorIncompleteQuantifier:             "Incomplete quantifier",
	errorInvalidQuantifier:                "Invalid quantifier",
	errorInvalidGroup:                     "Invalid group",
	errorMultipleFlagDashes:               "Multiple dashes in flag group",
	errorNotLinear:                        "Cannot be executed in linear time",
	errorRepeatedFlag:                     "Repeated flag in flag group",
	errorInvalidFlagGroup:                 "Invalid flag group",
	errorTooManyCaptures:                  "Too many captures",
	errorInvalidCaptureGroupName:          "Invalid capture group name",
	errorDuplicateCaptureGroupName:        "Duplicate capture group name",
	errorInvalidNamedReference:            "Invalid named reference",
	errorInvalidNamedCaptureReference:     "Invalid named capture referenced",
	errorInvalidClassPropertyName:         "Invalid property name in character class",
	errorInvalidCharacterClass:            "Invalid character class",
	errorUnterminatedCharacterClass:       "Unterminated character class",
	errorOutOfOrderCharacterClass:         "Range out of order in character class",
	errorInvalidClassSetOperation:         "Invalid set operation in character class",
	errorInvalidCharacterInClass:          "Invalid character in character class",
	errorNegatedCharacterClassWithStrings: "Negated character class may contain strings",
	errorUnsupportedBytecode:              "Unsupported Bytecode",
}

func errorString(code regexpError) string {
	if code >= errorNumErrors {
		panic("invalid regexp error code")
	}
	return regexpErrorMessages[code]
}

func errorIsStackOverflow(code regexpError) bool {
	return code == errorStackOverflow || code == errorAnalysisStackOverflow || code == errorTooLarge
}
