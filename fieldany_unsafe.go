//go:build !reflectx_strict

package saferefl

import "unsafe"

// fieldAny constructs an interface (any) value that points directly into the struct's
// field memory — zero allocations for all field kinds.
//
// For IfaceDirect types (pointer, map, chan, func):
//
//	eface.data = *fieldPtr  (the pointer/handle value itself)
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
func fieldAny(abiType unsafe.Pointer, ifaceDirect bool, fieldPtr unsafe.Pointer) any {
	var e eface
	e._typ = abiType
	if ifaceDirect {
		e.data = *(*unsafe.Pointer)(fieldPtr) //nolint:gosec
	} else {
		e.data = fieldPtr //nolint:gosec
	}
	return *(*any)(unsafe.Pointer(&e)) //nolint:gosec
}
