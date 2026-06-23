//go:build !go1.24 && !reflectx_strict

package unsafelayout

import "unsafe"

func init() { activeMapLayout = legacyMapLayout{} }

// legacyMapLayout implements mapLayout for Go's hmap (pre-1.24).
type legacyMapLayout struct{}

// MapLen reads hmap.count at byte offset 0.
// The count field has been the first field of hmap since Go 1 by design —
// the len() built-in relies on this layout. Self-test verifies it at startup.
func (legacyMapLayout) MapLen(m unsafe.Pointer) int { return *(*int)(m) }

// MapLen is the concrete, inlineable package-level function.
// Uses the same direct read as legacyMapLayout so the compiler can inline it across packages,
// unlike the interface-dispatch path used only by selfTestMap.
func MapLen(m unsafe.Pointer) int { return *(*int)(m) }
