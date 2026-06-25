package saferefl

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"unsafe"
)

// eface is the memory layout of an empty interface (any).
// For pointer types stored in an interface, the data word holds the pointer value directly.
type eface struct {
	_typ unsafe.Pointer
	data unsafe.Pointer
}

// kindOffset is the byte offset of Kind_ in abi.Type, valid on all platforms:
//
//	abi.Type layout:
//	  Size_       uintptr  (+0)
//	  PtrBytes    uintptr  (+ptrSize)
//	  Hash        uint32   (+2*ptrSize)
//	  Tflag       uint8    (+2*ptrSize+4)
//	  Align_      uint8    (+2*ptrSize+5)
//	  FieldAlign_ uint8    (+2*ptrSize+6)
//	  Kind_       uint8    (+2*ptrSize+7)  ← this field
//
// 64-bit: 2×8+7 = 23. 32-bit: 2×4+7 = 15. Evaluated at compile time.
const kindOffset = 2*unsafe.Sizeof(uintptr(0)) + 7

func init() {
	// Verify that kindOffset points to the Kind_ byte of abi.Type.
	// If the layout ever changes, Get/Set return wrong errors instead of corrupting data,
	// but we want to catch this early rather than emit confusing diagnostics.
	var x *int
	iface := any(x)
	e := (*eface)(unsafe.Pointer(&iface)) //nolint:gosec
	if efaceKind(e._typ) != reflect.Pointer {
		msg := "[saferefl] efaceKind self-test FAILED — abi.Type.Kind_ offset is wrong for this Go version; Get/Set will return incorrect errors"
		if _, strict := os.LookupEnv("SAFEREFL_STRICT"); strict {
			panic(msg)
		}
		log.Println(msg)
	}
}

// efaceKind returns the reflect.Kind of the dynamic type stored in the eface type word
// without constructing a full reflect.Type. It reads the Kind_ byte from Go's internal
// abi.Type struct at [kindOffset] and masks the low 5 bits (the kind, without flags).
//
// Note: abi.Type.Kind_ stores both the kind (low 5 bits) and, on Go 1.22–1.25, the
// directiface flag (bit 5). On Go 1.26+ the directiface flag moved to TFlag[+20].
// efaceKind masks the low 5 bits only, so it remains correct on all versions.
// Verified at package init.
//
//go:nosplit
func efaceKind(typ unsafe.Pointer) reflect.Kind {
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
