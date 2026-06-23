//go:build go1.24 && !reflectx_strict

package saferefl

import "unsafe"

// MapLenFast returns the number of elements in m without reflect interface dispatch.
// Reads Map.used (uint64 at offset 0) from the Go 1.24+ Swiss Tables map header.
// The layout assumption is verified by the self-test in internal/unsafelayout at package init.
// Call [EnableAccel] once at startup to confirm; absence of a call is safe.
func MapLenFast[K comparable, V any](m map[K]V) int {
	if m == nil {
		return 0
	}
	return int(*(*uint64)(*(*unsafe.Pointer)(unsafe.Pointer(&m))))
}
