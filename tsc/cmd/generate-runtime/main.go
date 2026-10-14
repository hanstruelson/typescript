package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/microsoft/TypeScript/tsc/internal/goemit"
)

func main() {
	source, err := goemit.RuntimePackageSource()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	_, sourceFile, _, _ := runtime.Caller(0)
	output := filepath.Join(filepath.Dir(sourceFile), "..", "..", "runtime", "runtime_generated.go")
	if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := os.WriteFile(output, source, 0644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
