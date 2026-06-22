package saferefl

import (
	"reflect"
	"unsafe"
)

// Set sets the value of fieldPath on obj to val.
// fieldPath supports dot-separated paths (e.g. "Address.City").
// obj must be a non-nil pointer to a struct.
//
// Returns [ReadOnlyError] for unexported fields,
// [TypeMismatchError] if T is not assignable to the field's type,
// or [FieldNotFoundError] if the path does not exist.
func Set[T any](obj any, fieldPath string, val T) error {
	objPtr, structType, err := structPtrOf(obj)
	if err != nil {
		return err
	}
	ptr, fm, err := resolvePath(objPtr, structType, fieldPath)
	if err != nil {
		return err
	}
	if !fm.Exported {
		return &ReadOnlyError{FieldPath: fieldPath}
	}

	wantType := reflect.TypeOf((*T)(nil)).Elem()
	if !wantType.AssignableTo(fm.Type) {
		return &TypeMismatchError{
			FieldPath: fieldPath,
			FieldType: fm.Type.String(),
			WantType:  wantType.String(),
		}
	}

	// Fast path: identical types — direct write, zero allocations.
	// Safety: ptr was obtained via reflect-verified offset arithmetic on a live object.
	if fm.Type == wantType {
		*(*T)(ptr) = val
		return nil
	}

	// Slow path: T is a concrete type assignable to an interface-typed field
	// (e.g. Set[*os.File](obj, "reader", f) where the field type is io.Reader).
	// Isolated into a noinline function so that taking &val there does NOT
	// cause val to escape to the heap in this (fast-path) function.
	setSlowPath(ptr, fm.Type, wantType, val)
	return nil
}

// setSlowPath writes a concrete value of type T into a field whose runtime type
// is an interface that T implements.
//
// The //go:noinline directive is load-bearing: it prevents escape analysis from
// seeing &val here and marking val as heap-escaping in the calling Set[T] function,
// keeping the common fast path at zero allocations.
//
//go:noinline
func setSlowPath[T any](fieldPtr unsafe.Pointer, dstType, srcType reflect.Type, val T) {
	reflect.NewAt(dstType, fieldPtr).Elem().Set(
		reflect.NewAt(srcType, unsafe.Pointer(&val)).Elem(),
	)
}

// MustSet is like [Set] but panics on any error.
func MustSet[T any](obj any, fieldPath string, val T) {
	if err := Set(obj, fieldPath, val); err != nil {
		panic("saferefl.MustSet: " + err.Error())
	}
}
