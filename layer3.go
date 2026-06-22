package saferefl

import (
	"unsafe"

	"github.com/lkmavi/saferefl/internal/unsafelayout"
)

// AccelAvailable reports whether the Layer 3 unsafe accelerator passed its
// self-test on the current Go version and architecture.
// When false, Layer 3 operations degrade gracefully to their safe equivalents.
func AccelAvailable() bool {
	return unsafelayout.AccelAvailable()
}

// UnsafeSliceAt returns a typed pointer to element at index in s without
// bounds checking. The caller must ensure 0 <= index < len(s).
//
// The slice data layout is a Go language guarantee (first word of the slice
// header is always the element array pointer), so no self-test is needed.
// This is strictly faster than the bounds-checked &s[index] when the caller
// can prove safety externally.
func UnsafeSliceAt[T any](s []T, index int) *T {
	// Read the data pointer from the slice header (offset 0, always).
	data := *(*unsafe.Pointer)(unsafe.Pointer(&s))
	return (*T)(unsafe.Add(data, uintptr(index)*unsafe.Sizeof(*new(T))))
}

// MapLenFast is defined in layer3_map_*.go (build-tag split) to eliminate
// interface dispatch overhead. See those files for implementation details.
