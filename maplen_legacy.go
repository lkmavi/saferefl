//go:build !go1.24 && !reflectx_strict

package saferefl

import "unsafe"

// MapLenFast returns the number of elements in m without reflect interface dispatch.
// Reads hmap.count (int at offset 0) from the pre-1.24 Go map header.
// The layout assumption is verified by the self-test in internal/unsafelayout at package init.
// Call [EnableAccel] once at startup to confirm; absence of a call is safe.
func MapLenFast[K comparable, V any](m map[K]V) int {
	if m == nil {
		return 0
	}
	return *(*int)(*(*unsafe.Pointer)(unsafe.Pointer(&m)))
}
