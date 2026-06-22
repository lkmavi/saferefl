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
// conversion automatically. Slower than Get but more convenient.
func (a Accessor[T]) GetFrom(obj any) (T, error) {
	p, _, err := structPtrOf(obj)
	if err != nil {
		var zero T
		return zero, err
	}
	return a.Get(p), nil
}

// SetOn writes val to the field in obj, handling the interface-to-pointer
// conversion automatically. Slower than Set but more convenient.
func (a Accessor[T]) SetOn(obj any, val T) error {
	p, _, err := structPtrOf(obj)
	if err != nil {
		return err
	}
	a.Set(p, val)
	return nil
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
