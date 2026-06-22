//go:build reflectx_strict

package saferefl

// MapLenFast falls back to len(m) under reflectx_strict (all unsafe paths disabled).
func MapLenFast[K comparable, V any](m map[K]V) int { return len(m) }
