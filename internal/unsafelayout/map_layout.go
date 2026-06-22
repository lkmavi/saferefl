//go:build !reflectx_strict

package unsafelayout

import "unsafe"

// mapLayout abstracts over the Go runtime's internal map representation.
// Exactly one implementation is registered via init() in the version-specific files.
type mapLayout interface {
	// MapLen returns the live element count for the map whose header pointer is m.
	MapLen(m unsafe.Pointer) int
}

var activeMapLayout mapLayout
