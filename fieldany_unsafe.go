//go:build !reflectx_strict

package saferefl

import (
	"reflect"
	"unsafe"

	"github.com/lkmavi/saferefl/internal/typeinfo"
)

// fieldAny constructs an interface (any) value that points directly into the struct's
// field memory — zero allocations for all non-interface field kinds.
//
// For IfaceDirect types (pointer, map, chan, func):
//
//	eface.data = *fieldPtr  (the pointer/handle value itself)
//
// For interface-kinded fields (any, io.Reader, etc.):
//
//	the field memory is already an interface; use reflect.NewAt to unwrap correctly.
//	This path allocates but is correct — re-boxing via eface would nest the interface.
//
// For all other types (int, string, bool, slice, struct, …):
//
//	eface.data = fieldPtr   (a pointer to the value inside the struct)
//
// # GC safety
//
// The struct must have escaped to the heap before reaching this call, which is
// guaranteed because it was passed as any to EachField/ToMap (Go escape analysis
// always heap-allocates values whose address is taken and passed as interface).
// The GC traces eface.data as a pointer, so the struct stays alive as long as
// any returned interface value is reachable.
//
// # Aliasing
//
// For non-direct types the returned any aliases the struct field. Type-asserting
// the value copies the field value out, so reads are safe. The struct must not be
// written to concurrently while EachField/ToMap callbacks execute.
func fieldAny(entry *typeinfo.IterEntry, fieldPtr unsafe.Pointer) any {
	if entry.IfaceDirect {
		var e eface
		e._typ = entry.AbiType
		e.data = *(*unsafe.Pointer)(fieldPtr)
		return *(*any)(unsafe.Pointer(&e)) //nolint:gosec
	}
	// Interface-kinded fields: the field memory is itself an interface value.
	// Setting eface.data = fieldPtr would produce a nested interface (any-in-any).
	// Use reflect to read and return the concrete value stored in the field.
	if efaceKind(entry.AbiType) == reflect.Interface {
		return reflect.NewAt(entry.Type, fieldPtr).Elem().Interface()
	}
	var e eface
	e._typ = entry.AbiType
	e.data = fieldPtr
	return *(*any)(unsafe.Pointer(&e)) //nolint:gosec
}
