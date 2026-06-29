package saferefl

import (
	"reflect"
	"unsafe"
)

// KindOf returns the reflect.Kind of v without constructing a full reflect.Type.
// Uses the same abi.Type fast-path as [Get]/[Set] — ~1 ns vs ~26 ns for reflect.TypeOf.
// Returns [reflect.Invalid] for nil interfaces.
func KindOf(v any) reflect.Kind {
	if v == nil {
		return reflect.Invalid
	}
	e := (*eface)(unsafe.Pointer(&v)) //nolint:gosec
	if e._typ == nil {
		return reflect.Invalid
	}
	return efaceKind(e._typ)
}

// TypeOf returns the [reflect.Type] for T without requiring a value of that type.
// Equivalent to reflect.TypeOf((*T)(nil)).Elem() but reads more clearly at call sites.
// Commonly used as a comparable type key in dependency-injection containers and registries.
func TypeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

// IsNil reports whether v is a nil interface, nil pointer, nil map, nil channel,
// nil function, or nil slice. Returns false for all non-nilable types (int, string, …).
//
// Pointer, map, channel, and UnsafePointer nils are detected via the interface data word.
// Func and slice fall back to reflect: slices use indirect interface storage (eface.data
// points to the slice header, not the header itself), and func's layout is more complex.
func IsNil(v any) bool {
	if v == nil {
		return true
	}
	e := (*eface)(unsafe.Pointer(&v)) //nolint:gosec
	if e._typ == nil {
		return true
	}
	switch efaceKind(e._typ) {
	case reflect.Pointer, reflect.Map, reflect.Chan, reflect.UnsafePointer:
		return e.data == nil
	case reflect.Func, reflect.Slice:
		return reflect.ValueOf(v).IsNil()
	default:
		return false
	}
}
