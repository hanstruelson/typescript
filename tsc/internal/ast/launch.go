package ast

// Launch operator overrides are ordinary methods in the shared TypeScript AST.
// These names are compiler-internal and are never spelled in Launch syntax.
const (
	launchOperatorPrefix = "__launch_operator_"

	LaunchOperatorPlusMethodName      = launchOperatorPrefix + "plus"
	LaunchOperatorAssignMethodName    = launchOperatorPrefix + "assign"
	LaunchOperatorSliceMethodName     = launchOperatorPrefix + "slice"
	LaunchOperatorMinusMethodName     = launchOperatorPrefix + "minus"
	LaunchOperatorMultiplyMethodName  = launchOperatorPrefix + "multiply"
	LaunchOperatorDivideMethodName    = launchOperatorPrefix + "divide"
	LaunchOperatorRemainderMethodName = launchOperatorPrefix + "remainder"
	LaunchOperatorPowerMethodName     = launchOperatorPrefix + "power"
	LaunchOperatorShiftLeftMethodName = launchOperatorPrefix + "shift_left"
	LaunchOperatorShiftRightMethodName = launchOperatorPrefix + "shift_right"
	LaunchOperatorUnsignedShiftRightMethodName = launchOperatorPrefix + "unsigned_shift_right"
	LaunchOperatorBitAndMethodName    = launchOperatorPrefix + "bit_and"
	LaunchOperatorBitOrMethodName     = launchOperatorPrefix + "bit_or"
	LaunchOperatorBitXorMethodName    = launchOperatorPrefix + "bit_xor"
	LaunchOperatorNotMethodName       = launchOperatorPrefix + "not"
	LaunchOperatorBitNotMethodName    = launchOperatorPrefix + "bit_not"
	LaunchOperatorLogicalAndMethodName = launchOperatorPrefix + "logical_and"
	LaunchOperatorLogicalOrMethodName = launchOperatorPrefix + "logical_or"
	LaunchOperatorNullishMethodName   = launchOperatorPrefix + "nullish"
	LaunchOperatorLessThanMethodName  = launchOperatorPrefix + "less_than"
	LaunchOperatorLessEqualMethodName = launchOperatorPrefix + "less_equal"
	LaunchOperatorGreaterThanMethodName = launchOperatorPrefix + "greater_than"
	LaunchOperatorGreaterEqualMethodName = launchOperatorPrefix + "greater_equal"
	LaunchOperatorEqualMethodName     = launchOperatorPrefix + "equal"
	LaunchOperatorNotEqualMethodName  = launchOperatorPrefix + "not_equal"
	LaunchOperatorStrictEqualMethodName = launchOperatorPrefix + "strict_equal"
	LaunchOperatorStrictNotEqualMethodName = launchOperatorPrefix + "strict_not_equal"
	LaunchOperatorInMethodName        = launchOperatorPrefix + "in"
	LaunchOperatorInstanceOfMethodName = launchOperatorPrefix + "instanceof"
	LaunchOperatorCommaMethodName     = launchOperatorPrefix + "comma"
	LaunchOperatorIncrementMethodName = launchOperatorPrefix + "increment"
	LaunchOperatorDecrementMethodName = launchOperatorPrefix + "decrement"

	LaunchOperatorPlusAssignMethodName = launchOperatorPrefix + "plus_assign"
	LaunchOperatorMinusAssignMethodName = launchOperatorPrefix + "minus_assign"
	LaunchOperatorMultiplyAssignMethodName = launchOperatorPrefix + "multiply_assign"
	LaunchOperatorDivideAssignMethodName = launchOperatorPrefix + "divide_assign"
	LaunchOperatorRemainderAssignMethodName = launchOperatorPrefix + "remainder_assign"
	LaunchOperatorPowerAssignMethodName = launchOperatorPrefix + "power_assign"
	LaunchOperatorShiftLeftAssignMethodName = launchOperatorPrefix + "shift_left_assign"
	LaunchOperatorShiftRightAssignMethodName = launchOperatorPrefix + "shift_right_assign"
	LaunchOperatorUnsignedShiftRightAssignMethodName = launchOperatorPrefix + "unsigned_shift_right_assign"
	LaunchOperatorBitAndAssignMethodName = launchOperatorPrefix + "bit_and_assign"
	LaunchOperatorBitOrAssignMethodName = launchOperatorPrefix + "bit_or_assign"
	LaunchOperatorBitXorAssignMethodName = launchOperatorPrefix + "bit_xor_assign"
	LaunchOperatorLogicalAndAssignMethodName = launchOperatorPrefix + "logical_and_assign"
	LaunchOperatorLogicalOrAssignMethodName = launchOperatorPrefix + "logical_or_assign"
	LaunchOperatorNullishAssignMethodName = launchOperatorPrefix + "nullish_assign"
)

// LaunchOperatorMethodName gives every overloadable binary operator its own
// hidden member. Unary +, -, !, and ~ share the corresponding member name and
// are distinguished from binary overloads by their zero-argument signature.
func LaunchOperatorMethodName(operator Kind) string {
	switch operator {
	case KindPlusToken:
		return LaunchOperatorPlusMethodName
	case KindMinusToken:
		return LaunchOperatorMinusMethodName
	case KindAsteriskToken:
		return LaunchOperatorMultiplyMethodName
	case KindSlashToken:
		return LaunchOperatorDivideMethodName
	case KindPercentToken:
		return LaunchOperatorRemainderMethodName
	case KindAsteriskAsteriskToken:
		return LaunchOperatorPowerMethodName
	case KindLessThanLessThanToken:
		return LaunchOperatorShiftLeftMethodName
	case KindGreaterThanGreaterThanToken:
		return LaunchOperatorShiftRightMethodName
	case KindGreaterThanGreaterThanGreaterThanToken:
		return LaunchOperatorUnsignedShiftRightMethodName
	case KindAmpersandToken:
		return LaunchOperatorBitAndMethodName
	case KindBarToken:
		return LaunchOperatorBitOrMethodName
	case KindCaretToken:
		return LaunchOperatorBitXorMethodName
	case KindExclamationToken:
		return LaunchOperatorNotMethodName
	case KindTildeToken:
		return LaunchOperatorBitNotMethodName
	case KindAmpersandAmpersandToken:
		return LaunchOperatorLogicalAndMethodName
	case KindBarBarToken:
		return LaunchOperatorLogicalOrMethodName
	case KindQuestionQuestionToken:
		return LaunchOperatorNullishMethodName
	case KindLessThanToken:
		return LaunchOperatorLessThanMethodName
	case KindLessThanEqualsToken:
		return LaunchOperatorLessEqualMethodName
	case KindGreaterThanToken:
		return LaunchOperatorGreaterThanMethodName
	case KindGreaterThanEqualsToken:
		return LaunchOperatorGreaterEqualMethodName
	case KindEqualsEqualsToken:
		return LaunchOperatorEqualMethodName
	case KindExclamationEqualsToken:
		return LaunchOperatorNotEqualMethodName
	case KindEqualsEqualsEqualsToken:
		return LaunchOperatorStrictEqualMethodName
	case KindExclamationEqualsEqualsToken:
		return LaunchOperatorStrictNotEqualMethodName
	case KindInKeyword:
		return LaunchOperatorInMethodName
	case KindInstanceOfKeyword:
		return LaunchOperatorInstanceOfMethodName
	case KindCommaToken:
		return LaunchOperatorCommaMethodName
	case KindEqualsToken:
		return LaunchOperatorAssignMethodName
	case KindPlusEqualsToken:
		return LaunchOperatorPlusAssignMethodName
	case KindMinusEqualsToken:
		return LaunchOperatorMinusAssignMethodName
	case KindAsteriskEqualsToken:
		return LaunchOperatorMultiplyAssignMethodName
	case KindSlashEqualsToken:
		return LaunchOperatorDivideAssignMethodName
	case KindPercentEqualsToken:
		return LaunchOperatorRemainderAssignMethodName
	case KindAsteriskAsteriskEqualsToken:
		return LaunchOperatorPowerAssignMethodName
	case KindLessThanLessThanEqualsToken:
		return LaunchOperatorShiftLeftAssignMethodName
	case KindGreaterThanGreaterThanEqualsToken:
		return LaunchOperatorShiftRightAssignMethodName
	case KindGreaterThanGreaterThanGreaterThanEqualsToken:
		return LaunchOperatorUnsignedShiftRightAssignMethodName
	case KindAmpersandEqualsToken:
		return LaunchOperatorBitAndAssignMethodName
	case KindBarEqualsToken:
		return LaunchOperatorBitOrAssignMethodName
	case KindCaretEqualsToken:
		return LaunchOperatorBitXorAssignMethodName
	case KindAmpersandAmpersandEqualsToken:
		return LaunchOperatorLogicalAndAssignMethodName
	case KindBarBarEqualsToken:
		return LaunchOperatorLogicalOrAssignMethodName
	case KindQuestionQuestionEqualsToken:
		return LaunchOperatorNullishAssignMethodName
	case KindPlusPlusToken:
		return LaunchOperatorIncrementMethodName
	case KindMinusMinusToken:
		return LaunchOperatorDecrementMethodName
	default:
		return ""
	}
}

// LaunchOperatorBaseMethodName maps a compound assignment to its binary
// operator. A dedicated compound override takes precedence when present.
func LaunchOperatorBaseMethodName(operator Kind) string {
	switch operator {
	case KindPlusEqualsToken:
		return LaunchOperatorPlusMethodName
	case KindMinusEqualsToken:
		return LaunchOperatorMinusMethodName
	case KindAsteriskEqualsToken:
		return LaunchOperatorMultiplyMethodName
	case KindSlashEqualsToken:
		return LaunchOperatorDivideMethodName
	case KindPercentEqualsToken:
		return LaunchOperatorRemainderMethodName
	case KindAsteriskAsteriskEqualsToken:
		return LaunchOperatorPowerMethodName
	case KindLessThanLessThanEqualsToken:
		return LaunchOperatorShiftLeftMethodName
	case KindGreaterThanGreaterThanEqualsToken:
		return LaunchOperatorShiftRightMethodName
	case KindGreaterThanGreaterThanGreaterThanEqualsToken:
		return LaunchOperatorUnsignedShiftRightMethodName
	case KindAmpersandEqualsToken:
		return LaunchOperatorBitAndMethodName
	case KindBarEqualsToken:
		return LaunchOperatorBitOrMethodName
	case KindCaretEqualsToken:
		return LaunchOperatorBitXorMethodName
	case KindAmpersandAmpersandEqualsToken:
		return LaunchOperatorLogicalAndMethodName
	case KindBarBarEqualsToken:
		return LaunchOperatorLogicalOrMethodName
	case KindQuestionQuestionEqualsToken:
		return LaunchOperatorNullishMethodName
	default:
		return ""
	}
}

// LaunchOperatorSpelling provides the source spelling for printing a hidden
// operator method declaration back as `operator override <op>`.
func LaunchOperatorSpelling(methodName string) (string, bool) {
	switch methodName {
	case LaunchOperatorPlusMethodName: return "+", false
	case LaunchOperatorMinusMethodName: return "-", false
	case LaunchOperatorMultiplyMethodName: return "*", false
	case LaunchOperatorDivideMethodName: return "/", false
	case LaunchOperatorRemainderMethodName: return "%", false
	case LaunchOperatorPowerMethodName: return "**", false
	case LaunchOperatorShiftLeftMethodName: return "<<", false
	case LaunchOperatorShiftRightMethodName: return ">>", false
	case LaunchOperatorUnsignedShiftRightMethodName: return ">>>", false
	case LaunchOperatorBitAndMethodName: return "&", false
	case LaunchOperatorBitOrMethodName: return "|", false
	case LaunchOperatorBitXorMethodName: return "^", false
	case LaunchOperatorNotMethodName: return "!", false
	case LaunchOperatorBitNotMethodName: return "~", false
	case LaunchOperatorLogicalAndMethodName: return "&&", false
	case LaunchOperatorLogicalOrMethodName: return "||", false
	case LaunchOperatorNullishMethodName: return "??", false
	case LaunchOperatorLessThanMethodName: return "<", false
	case LaunchOperatorLessEqualMethodName: return "<=", false
	case LaunchOperatorGreaterThanMethodName: return ">", false
	case LaunchOperatorGreaterEqualMethodName: return ">=", false
	case LaunchOperatorEqualMethodName: return "==", false
	case LaunchOperatorNotEqualMethodName: return "!=", false
	case LaunchOperatorStrictEqualMethodName: return "===", false
	case LaunchOperatorStrictNotEqualMethodName: return "!==", false
	case LaunchOperatorInMethodName: return "in", true
	case LaunchOperatorInstanceOfMethodName: return "instanceof", true
	case LaunchOperatorCommaMethodName: return ",", false
	case LaunchOperatorAssignMethodName: return "=", false
	case LaunchOperatorPlusAssignMethodName: return "+=", false
	case LaunchOperatorMinusAssignMethodName: return "-=", false
	case LaunchOperatorMultiplyAssignMethodName: return "*=", false
	case LaunchOperatorDivideAssignMethodName: return "/=", false
	case LaunchOperatorRemainderAssignMethodName: return "%=", false
	case LaunchOperatorPowerAssignMethodName: return "**=", false
	case LaunchOperatorShiftLeftAssignMethodName: return "<<=", false
	case LaunchOperatorShiftRightAssignMethodName: return ">>=", false
	case LaunchOperatorUnsignedShiftRightAssignMethodName: return ">>>=", false
	case LaunchOperatorBitAndAssignMethodName: return "&=", false
	case LaunchOperatorBitOrAssignMethodName: return "|=", false
	case LaunchOperatorBitXorAssignMethodName: return "^=", false
	case LaunchOperatorLogicalAndAssignMethodName: return "&&=", false
	case LaunchOperatorLogicalOrAssignMethodName: return "||=", false
	case LaunchOperatorNullishAssignMethodName: return "??=", false
	case LaunchOperatorIncrementMethodName: return "++", false
	case LaunchOperatorDecrementMethodName: return "--", false
	case LaunchOperatorSliceMethodName: return "slice", true
	default: return "", false
	}
}
