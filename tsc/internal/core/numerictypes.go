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
