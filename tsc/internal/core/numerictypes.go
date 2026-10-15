package core

import (
	"strconv"
	"strings"
)

// NativeNumericPromotion selects a native arithmetic result. Mixed integer
// operands are converted with range checks before the operation is executed.
func NativeNumericPromotion(a, b string) string {
	if a == b {
		return a
	}
	if a == "number" || b == "number" || a == "float64" || b == "float64" {
		return "float64"
	}
	if a == "float32" || b == "float32" {
		return "float32"
	}
	width := func(t string) int {
		switch t {
		case "int8", "uint8":
			return 8
		case "int16", "uint16":
			return 16
		case "int32", "uint32":
			return 32
		case "int64", "uint64":
			return 64
		default:
			return strconv.IntSize
		}
	}
	bits := max(width(a), width(b))
	unsigned := strings.HasPrefix(a, "uint") && strings.HasPrefix(b, "uint")
	if a == "uint64" || b == "uint64" || (strconv.IntSize == 64 && (a == "uint" || b == "uint")) {
		unsigned = true
	}
	prefix := "int"
	if unsigned {
		prefix = "uint"
	}
	return prefix + strconv.Itoa(bits)
}

// ConversionMethodTarget names explicit primitive conversion intrinsics.
func ConversionMethodTarget(name string) string {
	switch name {
	case "toNumber", "number", "toFloat64", "float64":
		return "number"
	case "toFloat32", "float32":
		return "float32"
	case "toInt", "int":
		return "int"
	case "toInt8", "int8":
		return "int8"
	case "toInt16", "int16":
		return "int16"
	case "toInt32", "int32":
		return "int32"
	case "toInt64", "int64":
		return "int64"
	case "toUint", "uint":
		return "uint"
	case "toUint8", "uint8":
		return "uint8"
	case "toUint16", "uint16":
		return "uint16"
	case "toUint32", "uint32":
		return "uint32"
	case "toUint64", "uint64":
		return "uint64"
	case "toString", "string":
		return "string"
	case "toBoolean", "boolean":
		return "boolean"
	}
	return ""
}

// NativeNumericLosslessConversion reports whether every value of source is
// exactly representable in target. This is a type guarantee, not a value check.
func NativeNumericLosslessConversion(target, source string) bool {
	if target == "float64" {
		target = "number"
	}
	if source == "float64" {
		source = "number"
	}
	width := func(kind string) int {
		switch kind {
		case "int8", "uint8":
			return 8
		case "int16", "uint16":
			return 16
		case "int32", "uint32":
			return 32
		case "int64", "uint64":
			return 64
		case "int", "uint":
			return strconv.IntSize
		}
		return 0
	}
	sourceBits, targetBits := width(source), width(target)
	if target == source {
		return sourceBits != 0 || target == "number" || target == "float32"
	}
	if source == "float32" && target == "number" {
		return true
	}
	if sourceBits == 0 {
		return false
	}
	if target == "number" || target == "float32" {
		precision := 53
		if target == "float32" {
			precision = 24
		}
		return sourceBits <= precision
	}
	if targetBits == 0 {
		return false
	}
	sourceUnsigned, targetUnsigned := strings.HasPrefix(source, "uint"), strings.HasPrefix(target, "uint")
	if sourceUnsigned == targetUnsigned {
		return targetBits >= sourceBits
	}
	return sourceUnsigned && !targetUnsigned && targetBits > sourceBits
}

// NativeNumericCommonType returns the least-cost native type that exactly
// represents every value of both source types. An empty target means the pair
// needs an exact comparison path instead of conversion to a common type.
//
// The selector considers third types as well as either input type. For
// example, int32 and uint32 share int64, while int32 and float32 share
// float64. It prefers fewer conversions, then a smaller representation, then
// integer targets over floating targets when both plans have the same cost.
func NativeNumericCommonType(a, b string) string {
	canonical := func(kind string) string {
		if kind == "number" {
			return "float64"
		}
		return kind
	}
	a, b = canonical(a), canonical(b)
	if !isNativeNumericType(a) || !isNativeNumericType(b) {
		return ""
	}
	if a == b {
		return a
	}

	targets := [...]string{
		"int8", "uint8", "int16", "uint16", "int32", "uint32",
		"int64", "uint64", "int", "uint", "float32", "float64",
	}
	width := func(kind string) int {
		switch kind {
		case "int8", "uint8":
			return 8
		case "int16", "uint16":
			return 16
		case "int32", "uint32", "float32":
			return 32
		case "int64", "uint64", "float64":
			return 64
		case "int", "uint":
			return strconv.IntSize
		default:
			return 0
		}
	}
	integerTarget := func(kind string) bool { return kind != "float32" && kind != "float64" }

	best := ""
	bestConversions, bestWidth := int(^uint(0)>>1), int(^uint(0)>>1)
	bestInteger := false
	for _, target := range targets {
		if !NativeNumericLosslessConversion(target, a) || !NativeNumericLosslessConversion(target, b) {
			continue
		}
		conversions := 0
		if canonical(target) != a {
			conversions++
		}
		if canonical(target) != b {
			conversions++
		}
		targetWidth := width(target)
		targetInteger := integerTarget(target)
		if conversions < bestConversions ||
			conversions == bestConversions && targetWidth < bestWidth ||
			conversions == bestConversions && targetWidth == bestWidth && targetInteger && !bestInteger {
			best = target
			bestConversions = conversions
			bestWidth = targetWidth
			bestInteger = targetInteger
		}
	}
	return best
}

func isNativeNumericType(kind string) bool {
	switch kind {
	case "number", "float32", "float64", "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return true
	default:
		return false
	}
}
