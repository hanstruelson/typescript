//go:build cgo

// Package openssl exposes the small OpenSSL surface used by generated Node
// compatibility code. Hash operations use the OpenSSL EVP digest provider;
// the rest of the runtime can continue to use Go's standard crypto packages.
package openssl

/*
#cgo pkg-config: openssl
#include <openssl/crypto.h>
#include <openssl/evp.h>
#include <openssl/opensslv.h>
#include <stdlib.h>

static EVP_MD *ts_openssl_fetch_digest(const char *name) {
#if OPENSSL_VERSION_NUMBER >= 0x30000000L
	return EVP_MD_fetch(NULL, name, NULL);
#else
	return (EVP_MD *)EVP_get_digestbyname(name);
#endif
}

static void ts_openssl_free_digest(EVP_MD *digest) {
#if OPENSSL_VERSION_NUMBER >= 0x30000000L
	EVP_MD_free(digest);
#else
	(void)digest;
#endif
}

static int ts_openssl_fips_enabled(void) {
#if OPENSSL_VERSION_NUMBER >= 0x30000000L
	return EVP_default_properties_is_fips_enabled(NULL);
#else
	return 0;
#endif
}
*/
import "C"

import (
	"errors"
	"fmt"
	"hash"
	"runtime"
	"strings"
	"unsafe"
)

// Version returns the OpenSSL library version in the compact form Node exposes
// as process.versions.openssl.
func Version() string {
	text := C.GoString(C.OpenSSL_version(C.OPENSSL_VERSION))
	parts := strings.Fields(text)
	if len(parts) > 1 {
		return parts[1]
	}
	return text
}

// VersionNumber returns OpenSSL's numeric version constant.
func VersionNumber() uint64 { return uint64(C.OpenSSL_version_num()) }

// FIPSMode reports whether the OpenSSL default library context has FIPS
// properties enabled. It does not claim that Go's independent crypto
// implementations are operating in a validated FIPS mode.
func FIPSMode() bool { return C.ts_openssl_fips_enabled() == 1 }

// SupportsDigest reports whether the loaded OpenSSL providers expose name.
func SupportsDigest(name string) bool {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	digest := C.ts_openssl_fetch_digest(cname)
	if digest == nil {
		return false
	}
	C.ts_openssl_free_digest(digest)
	return true
}

// NewHash creates an incremental EVP digest backed by the loaded OpenSSL
// providers.
func NewHash(name string) (hash.Hash, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	digest := C.ts_openssl_fetch_digest(cname)
	if digest == nil {
		return nil, fmt.Errorf("OpenSSL digest %q is unavailable", name)
	}
	ctx := C.EVP_MD_CTX_new()
	if ctx == nil {
		C.ts_openssl_free_digest(digest)
		return nil, errors.New("OpenSSL could not allocate a digest context")
	}
	if C.EVP_DigestInit_ex(ctx, digest, nil) != 1 {
		C.EVP_MD_CTX_free(ctx)
		C.ts_openssl_free_digest(digest)
		return nil, fmt.Errorf("OpenSSL could not initialize digest %q", name)
	}
	value := &evpHash{ctx: ctx, digest: digest, size: int(C.EVP_MD_size(digest)), blockSize: int(C.EVP_MD_block_size(digest))}
	runtime.SetFinalizer(value, (*evpHash).Close)
	return value, nil
}

type evpHash struct {
	ctx       *C.EVP_MD_CTX
	digest    *C.EVP_MD
	size      int
	blockSize int
}

func (h *evpHash) Write(data []byte) (int, error) {
	if h.ctx == nil {
		return 0, errors.New("OpenSSL digest is closed")
	}
	if len(data) == 0 {
		return 0, nil
	}
	if C.EVP_DigestUpdate(h.ctx, unsafe.Pointer(unsafe.SliceData(data)), C.size_t(len(data))) != 1 {
		return 0, errors.New("OpenSSL digest update failed")
	}
	return len(data), nil
}

func (h *evpHash) Sum(destination []byte) []byte {
	if h.ctx == nil {
		panic("OpenSSL digest is closed")
	}
	copyContext := C.EVP_MD_CTX_new()
	if copyContext == nil {
		panic("OpenSSL could not allocate a digest copy")
	}
	defer C.EVP_MD_CTX_free(copyContext)
	if C.EVP_MD_CTX_copy_ex(copyContext, h.ctx) != 1 {
		panic("OpenSSL could not copy a digest context")
	}
	var output [C.EVP_MAX_MD_SIZE]C.uchar
	var length C.uint
	if C.EVP_DigestFinal_ex(copyContext, &output[0], &length) != 1 {
		panic("OpenSSL digest finalization failed")
	}
	return append(destination, C.GoBytes(unsafe.Pointer(&output[0]), C.int(length))...)
}

func (h *evpHash) Reset() {
	if h.ctx == nil {
		return
	}
	if C.EVP_DigestInit_ex(h.ctx, h.digest, nil) != 1 {
		panic("OpenSSL digest reset failed")
	}
}

func (h *evpHash) Size() int      { return h.size }
func (h *evpHash) BlockSize() int { return h.blockSize }

// Close releases the OpenSSL digest state. It is also called by the finalizer.
func (h *evpHash) Close() {
	if h.ctx == nil {
		return
	}
	C.EVP_MD_CTX_free(h.ctx)
	C.ts_openssl_free_digest(h.digest)
	h.ctx = nil
	h.digest = nil
	runtime.SetFinalizer(h, nil)
}
