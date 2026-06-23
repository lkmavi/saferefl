// Package saferefl provides fast, type-safe struct reflection for Go.
//
// # Generic API
//
// Access and modify struct fields by name with full type safety. Zero allocations
// on the warm-cache path and faster than [reflect.Value.FieldByName].
//
//	name, _ := saferefl.Get[string](&user, "Name")
//	_       = saferefl.Set[string](&user, "Name", "Alice")
//	name    = saferefl.MustGet[string](&user, "Name")
//
// Dot-separated paths traverse nested structs and pointer-to-struct fields:
//
//	city, _ := saferefl.Get[string](&employee, "Address.City")
//
// Struct inspection by name or tag:
//
//	fields, _ := saferefl.FieldsOf[User]()
//	sf, ok   := saferefl.FieldByName[User]("Name")
//
// # Accessor API — hot-path field access
//
// Bind a field path once at startup or statement-prepare time; subsequent
// Get/Set calls are pure pointer arithmetic with no per-call reflection.
// Ideal for ORM row scanning, DI injection, and struct copying in tight loops.
//
//	acc, _ := saferefl.MakeAccessor[int](&user, "Age")
//	ptr    := saferefl.UnsafePtrOf(&user)
//	age    := acc.Get(ptr)   // ~0.5 ns, 0 allocs
//	acc.Set(ptr, 31)
//
// # Unsafe Primitives — slice and map fast paths
//
// Direct element access without reflect interface dispatch.
// All functions self-test their runtime assumptions at package init; use
// [EnableAccel] at startup to surface any mismatch, or check [AccelAvailable].
//
//	// Element pointer from a slice — no bounds-check overhead.
//	p := saferefl.UnsafeSliceAt(slice, i)
//
//	// Map length without reflect.
//	n := saferefl.MapLenFast(m)
//
// # Error types
//
// All errors are typed and work with [errors.As]:
// [FieldNotFoundError], [TypeMismatchError], [ReadOnlyError].
package saferefl
