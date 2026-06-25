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
// # Tag-based access
//
// Access fields by struct tag value instead of field name. Useful for ORM-style
// field mapping where the tag is the canonical identifier:
//
//	name, _ := saferefl.GetByTag[string](&user, "json", "name")   // field with `json:"name"`
//	_        = saferefl.SetByTag[string](&user, "db", "col_name", "Alice")
//
// # Iteration and mapping
//
// Iterate over all exported fields of a struct in declaration order:
//
//	saferefl.EachField(&user, func(name string, val any) bool {
//	    fmt.Println(name, val)
//	    return true // false to stop early
//	})
//
// Convert a struct to a map (by field name or struct tag):
//
//	m, _  := saferefl.ToMap(&user)
//	m, _  := saferefl.ToMapByTag(&user, "json")   // keys are tag values
//	_      = saferefl.FromMap(m, &dst)             // populate struct from map
//
// Copy matching exported fields between two structs (DTO-to-entity mapping):
//
//	saferefl.CopyFields(&src, &dst)
//
// Iterate a typed map with early-stop support:
//
//	saferefl.MapForEach(m, func(k string, v int) bool { return true })
//
// # Introspection
//
// Fast kind and nil checks without reflect.Value boxing:
//
//	saferefl.KindOf(v)   // reflect.Kind, ~0.28 ns, 0 allocs
//	saferefl.IsNil(v)    // true for nil pointer/map/chan/func/slice/interface
//
// # TypeDescriptor — low-level cache
//
// Direct access to the prebuilt type metadata used internally by all APIs.
// Useful for plugin authors or code that needs repeated field access without
// the overhead of the generic Get/Set surface:
//
//	desc := saferefl.TypeDescriptorOf(reflect.TypeOf(user{}))
//	fm   := desc.FieldsByName["Name"]   // *FieldMeta with Offset, Type, Tag, …
//	fm   = desc.FieldsByTag["json"]["name"]
//
// # Error types
//
// All errors are typed and work with [errors.As] for detail inspection, and with
// [errors.Is] via sentinel values for simple checks:
//
//	errors.Is(err, saferefl.ErrFieldNotFound)
//	errors.Is(err, saferefl.ErrTypeMismatch)
//	errors.Is(err, saferefl.ErrReadOnly)
//
// Typed variants: [FieldNotFoundError], [TypeMismatchError], [ReadOnlyError].
package saferefl
