//go:build gojit && linux && amd64

package regexpengine

import "sync"

var jitMatchers sync.Map

// StoreJITMatcher makes a generated matcher available to the Go runtime
// bridge that requested its compilation.
func StoreJITMatcher(id uint64, matcher func([]uint16, int, []int32) (int, int, bool)) {
	if id == 0 || matcher == nil {
		panic("regexpengine: invalid JIT matcher registration")
	}
	if _, loaded := jitMatchers.LoadOrStore(id, matcher); loaded {
		panic("regexpengine: duplicate JIT matcher registration")
	}
}

// LoadJITMatcher retrieves the matcher installed by the generated JIT snippet.
func LoadJITMatcher(id uint64) (func([]uint16, int, []int32) (int, int, bool), bool) {
	matcher, ok := jitMatchers.Load(id)
	if !ok {
		return nil, false
	}
	return matcher.(func([]uint16, int, []int32) (int, int, bool)), true
}

// DeleteJITMatcher releases the temporary registration after compilation.
func DeleteJITMatcher(id uint64) {
	jitMatchers.Delete(id)
}
