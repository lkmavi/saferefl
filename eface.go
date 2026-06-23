package saferefl

import (
	"fmt"
	"reflect"
	"unsafe"
)

// eface is the memory layout of an empty interface (any).
// For pointer types stored in an interface, the data word holds the pointer value directly.
type eface struct {
	_typ unsafe.Pointer
	data unsafe.Pointer
}

// efaceKind returns the reflect.Kind of the dynamic type stored in the eface type word
// without constructing a full reflect.Type. It reads the Kind_ byte from Go's internal
// abi.Type struct at byte offset 23 (64-bit platforms).
//
// This offset has been stable since Go 1.18 and is the same byte reflect.Kind() returns,
// so any future change would break reflect itself first.
//
//go:nosplit
func efaceKind(typ unsafe.Pointer) reflect.Kind {
	// abi.Type layout (64-bit):
	//   Size_       uintptr  (0)
	//   PtrBytes    uintptr  (8)
	//   Hash        uint32   (16)
	//   Tflag       uint8    (20)
	//   Align_      uint8    (21)
	//   FieldAlign_ uint8    (22)
	//   Kind_       uint8    (23) ← this field
	const kindOffset = 23
	return reflect.Kind(*(*uint8)(unsafe.Pointer(uintptr(typ) + kindOffset)) & 0x1f) //nolint:gosec
}

// fromError builds a descriptive error for GetFrom/SetOn misuse.
// Kept out-of-line to avoid polluting the fast path's register allocation.
//
//go:noinline
func fromError(obj any) error {
	if obj == nil {
		return fmt.Errorf("saferefl: obj must not be nil")
	}
	t := reflect.TypeOf(obj)
	if t.Kind() != reflect.Pointer {
		return fmt.Errorf("saferefl: obj must be a pointer to struct, got %v", t.Kind())
	}
	return fmt.Errorf("saferefl: obj pointer must not be nil")
}
