//go:build reflectx_strict

package saferefl

// MapLenFast returns len(m). Under the reflectx_strict build tag all unsafe
// accelerators are disabled, so this falls back to the builtin.
func MapLenFast[K comparable, V any](m map[K]V) int {
	return len(m)
}
