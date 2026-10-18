//go:build !cgo

// Package openssl reports an unavailable OpenSSL backend when cgo is disabled.
package openssl

import (
	"errors"
	"hash"
)

func Version() string     { return "" }
func VersionNumber() uint64 { return 0 }
func FIPSMode() bool      { return false }
func SupportsDigest(string) bool { return false }
func NewHash(string) (hash.Hash, error) {
	return nil, errors.New("OpenSSL support requires cgo")
}
