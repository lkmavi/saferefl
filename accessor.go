package saferefl

import (
	"fmt"
	"reflect"
	"unsafe"
)

// Accessor holds a pre-resolved field binding for repeated typed access.
// Create once with [MakeAccessor]; Get and Set are then allocation-free.
//
// For hot loops that require maximum throughput, extract the struct pointer
// once with [UnsafePtrOf] and pass it to Get/Set directly.
type Accessor[T any] struct {
	// Fast path (chain == nil): direct offset from root struct pointer.
	// Covers simple fields and dot-paths through value-type intermediates.
	offset uintptr
	// Slow path: chain of steps for paths that cross pointer-to-struct fields.
	chain []accessorStep
}

// MakeAccessor resolves fieldPath for the struct type of obj once and returns
// a reusable Accessor[T]. obj is only used for its type; it need not be the
// actual object you will access at runtime.
//
// Returns [FieldNotFoundError] or [TypeMismatchError] on invalid paths/types.
func MakeAccessor[T any](obj any, fieldPath string) (Accessor[T], error) {
	if obj == nil {
		return Accessor[T]{}, fmt.Errorf("saferefl: obj must not be nil")
	}
	t := reflect.TypeOf(obj)
	if t.Kind() != reflect.Pointer {
		return Accessor[T]{}, fmt.Errorf("saferefl: obj must be a pointer to struct, got %v", t.Kind())
	}
	elem := t.Elem()
	if elem.Kind() != reflect.Struct {
		return Accessor[T]{}, fmt.Errorf("saferefl: obj must point to a struct, got pointer to %v", elem.Kind())
	}

	offset, chain, fieldType, err := resolveMetaPath(elem, fieldPath)
	if err != nil {
		return Accessor[T]{}, err
	}

	wantType := reflect.TypeOf((*T)(nil)).Elem()
	if !fieldType.AssignableTo(wantType) {
		return Accessor[T]{}, &TypeMismatchError{
			FieldPath: fieldPath,
			FieldType: fieldType.String(),
			WantType:  wantType.String(),
		}
	}

	return Accessor[T]{offset: offset, chain: chain}, nil
}

// Get reads the field from the struct pointed to by objPtr.
// objPtr must be the pointer to the same struct type used in [MakeAccessor].
// Use [UnsafePtrOf] to obtain objPtr from an interface value.
func (a Accessor[T]) Get(objPtr unsafe.Pointer) T {
	if a.chain == nil {
		return *(*T)(unsafe.Pointer(uintptr(objPtr) + a.offset))
	}
	return *(*T)(walkChain(objPtr, a.chain))
}

// Set writes val to the field in the struct pointed to by objPtr.
// objPtr must be the pointer to the same struct type used in [MakeAccessor].
// Use [UnsafePtrOf] to obtain objPtr from an interface value.
func (a Accessor[T]) Set(objPtr unsafe.Pointer, val T) {
	if a.chain == nil {
		*(*T)(unsafe.Pointer(uintptr(objPtr) + a.offset)) = val
		return
	}
	*(*T)(walkChain(objPtr, a.chain)) = val
}

// GetFrom reads the field from obj, handling the interface-to-pointer
// conversion automatically. Use Get with a pre-extracted [UnsafePtrOf]
// pointer for maximum throughput in tight loops.
func (a Accessor[T]) GetFrom(obj any) (T, error) {
	e := (*eface)(unsafe.Pointer(&obj))
	// Validate without reflect.TypeOf: read the kind byte directly from the
	// eface type word and compare. This avoids the full reflect.Type wrapping
	// and the Elem().Kind() check (the Accessor was already validated for a
	// pointer-to-struct type at MakeAccessor time).
	if e._typ == nil || efaceKind(e._typ) != reflect.Pointer || e.data == nil {
		var zero T
		return zero, fromError(obj)
	}
	return a.Get(e.data), nil
}

// SetOn writes val to the field in obj, handling the interface-to-pointer
// conversion automatically. Use Set with a pre-extracted [UnsafePtrOf]
// pointer for maximum throughput in tight loops.
func (a Accessor[T]) SetOn(obj any, val T) error {
	e := (*eface)(unsafe.Pointer(&obj))
	if e._typ == nil || efaceKind(e._typ) != reflect.Pointer || e.data == nil {
		return fromError(obj)
	}
	a.Set(e.data, val)
	return nil
}

// efaceKind returns the reflect.Kind of the dynamic type stored in the eface
// type word without constructing a full reflect.Type. It reads the Kind_ byte
// from Go's internal abi.Type struct at byte offset 23 (64-bit platforms).
// This offset has been stable since Go 1.18 and is the same byte reflect.Kind()
// returns, so any future change would break reflect itself first.
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
	return reflect.Kind(*(*uint8)(unsafe.Pointer(uintptr(typ) + kindOffset)) & 0x1f)
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

// UnsafePtrOf returns the raw pointer to the struct pointed to by obj.
// obj must be a non-nil pointer to a struct — same precondition as [Get]/[Set].
// Pass the result to [Accessor.Get]/[Accessor.Set] for zero-overhead field access.
func UnsafePtrOf(obj any) unsafe.Pointer {
	return (*eface)(unsafe.Pointer(&obj)).data
}

func walkChain(ptr unsafe.Pointer, chain []accessorStep) unsafe.Pointer {
	for _, step := range chain {
		ptr = unsafe.Pointer(uintptr(ptr) + step.offset)
		if step.ptr {
			ptr = *(*unsafe.Pointer)(ptr)
		}
	}
	return ptr
}
