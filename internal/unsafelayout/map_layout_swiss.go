//go:build go1.24 && !reflectx_strict

package unsafelayout

import "unsafe"

func init() { activeMapLayout = swissMapLayout{} }

// swissMapLayout implements mapLayout for Go 1.24+ Swiss Tables maps.
type swissMapLayout struct{}

// MapLen reads Map.used at byte offset 0.
// In Go 1.24 internal/runtime/maps, Map.used (uint64) is the first field.
// Self-test verifies this assumption at startup; on mismatch AccelAvailable returns false.
func (swissMapLayout) MapLen(m unsafe.Pointer) int { return int(*(*uint64)(m)) }

// MapLen is the concrete, inlineable package-level function.
// Uses the same direct read as swissMapLayout so the compiler can inline it across packages,
// unlike the interface-dispatch path used only by selfTestMap.
func MapLen(m unsafe.Pointer) int { return int(*(*uint64)(m)) }
