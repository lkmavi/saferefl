package saferefl

import (
	"unsafe"

	"github.com/lkmavi/saferefl/internal/unsafelayout"
)

// AccelAvailable reports whether the unsafe accelerator passed its self-test on
// the current Go version and architecture.
// When false, unsafe primitives degrade gracefully to their safe equivalents.
func AccelAvailable() bool {
	return unsafelayout.AccelAvailable()
}

// EnableAccel confirms that the unsafe accelerator is active.
// It runs a brief sanity check at package init; this function only reports the result.
//
// A non-nil error means the accelerator is disabled and operations fall back to
// stdlib reflect. Calling EnableAccel is optional — degradation is automatic and silent.
func EnableAccel() error {
	return unsafelayout.EnableAccel()
}

// UnsafeSliceAt returns a typed pointer to element at index in s without bounds
// checking. The caller must ensure 0 <= index < len(s).
//
// The slice data layout is a Go language guarantee (first word of the slice header
// is always the element array pointer), so no self-test is needed here.
// This is strictly faster than bounds-checked &s[index] when the caller can prove
// safety externally.
func UnsafeSliceAt[T any](s []T, index int) *T {
	data := *(*unsafe.Pointer)(unsafe.Pointer(&s)) //nolint:gosec
	return (*T)(unsafe.Add(data, uintptr(index)*unsafe.Sizeof(*new(T))))
}
