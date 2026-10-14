package goemit

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestGeneratedRuntimePackageIsCurrent(t *testing.T) {
	want, err := RuntimePackageSource()
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join("..", "..", "runtime", "runtime_generated.go")
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("%s is stale; run go generate ./internal/goemit", path)
	}
}
