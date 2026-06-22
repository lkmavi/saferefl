//go:build go1.24 && !reflectx_strict

package saferefl

import "unsafe"

// MapLenFast returns the number of elements in m using a direct read of
// Map.used (uint64 at offset 0 in Go 1.24+ Swiss Tables), bypassing interface dispatch.
// Layout verified by unsafelayout self-test at program startup.
// Returns 0 for a nil map without reading any memory.
func MapLenFast[K comparable, V any](m map[K]V) int {
	if m == nil {
		return 0
	}
	// A map variable IS a *Map pointer. Read Map.used at offset 0.
	return int(*(*uint64)(*(*unsafe.Pointer)(unsafe.Pointer(&m))))
}
