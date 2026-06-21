// Package saferefl provides fast, safe reflection for Go.
//
// Layer 1 — Generics API (this package):
//
//	name, _ := saferefl.Get[string](&user, "Name")
//	_       = saferefl.Set[string](&user, "Name", "Alice")
//	name    = saferefl.MustGet[string](&user, "Name")
//
// Dot-separated paths traverse nested structs and pointer-to-struct fields:
//
//	city, _ := saferefl.Get[string](&employee, "Address.City")
//
// Layer 2 — TypeInfo cache (internal/typeinfo):
// Precomputed field offsets used by Layer 1. Accessible via [internal/typeinfo.TypeDescriptorOf].
//
// Error types: [FieldNotFoundError], [TypeMismatchError], [ReadOnlyError].
package saferefl
