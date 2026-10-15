package core

import (
	"strconv"
	"testing"
)

func TestNativeNumericCommonType(t *testing.T) {
	tests := []struct {
		left, right, want string
	}{
		{"int32", "uint32", "int64"},
		{"uint32", "int32", "int64"},
		{"int32", "float32", "float64"},
		{"float32", "int32", "float64"},
		{"int16", "uint16", "int32"},
		{"float32", "int16", "float32"},
		{"uint8", "float32", "float32"},
		{"float32", "float64", "float64"},
		{"int64", "uint64", ""},
		{"int64", "float64", ""},
		{"uint64", "float64", ""},
		{"number", "float64", "float64"},
	}
	for _, test := range tests {
		t.Run(test.left+"_"+test.right, func(t *testing.T) {
			if got := NativeNumericCommonType(test.left, test.right); got != test.want {
				t.Fatalf("NativeNumericCommonType(%q, %q) = %q, want %q", test.left, test.right, got, test.want)
			}
		})
	}
}

func TestNativeNumericCommonTypeCoversEveryPair(t *testing.T) {
	kinds := []string{"number", "float32", "float64", "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64"}
	targets := append([]string(nil), kinds...)
	for _, left := range kinds {
		for _, right := range kinds {
			got := NativeNumericCommonType(left, right)
			if got == "" {
				for _, target := range targets {
					if NativeNumericLosslessConversion(target, left) && NativeNumericLosslessConversion(target, right) {
						t.Fatalf("pair %s/%s has safe target %s but planner selected exact comparison", left, right, target)
					}
				}
				continue
			}
			if !NativeNumericLosslessConversion(got, left) || !NativeNumericLosslessConversion(got, right) {
				t.Fatalf("planner selected lossy target %s for %s/%s", got, left, right)
			}
		}
	}

	if strconv.IntSize != 32 && strconv.IntSize != 64 {
		t.Fatalf("unexpected platform int width %d", strconv.IntSize)
	}
}
