//go:build !go1.24 && !reflectx_strict

package saferefl

import "unsafe"

// MapLenFast returns the number of elements in m using a direct read of
// hmap.count (int at offset 0), bypassing interface dispatch.
// For Go < 1.24 (classic hmap layout). Layout verified by unsafelayout self-test.
// Returns 0 for a nil map without reading any memory.
func MapLenFast[K comparable, V any](m map[K]V) int {
	if m == nil {
		return 0
	}
	// A map variable IS a *hmap pointer. Read hmap.count at offset 0.
	return *(*int)(*(*unsafe.Pointer)(unsafe.Pointer(&m)))
}
