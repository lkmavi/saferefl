package saferefl

import "reflect"

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

	// Fast path: identical types — zero-copy direct write.
	// Safety: ptr was obtained via reflect-verified offset arithmetic on a live object.
	if fm.Type == wantType {
		*(*T)(ptr) = val
		return nil
	}
	// Slow path: T is a concrete type assignable to an interface field.
	// reflect.ValueOf(&val).Elem() is safe even for nil interface values.
	reflect.NewAt(fm.Type, ptr).Elem().Set(reflect.ValueOf(&val).Elem())
	return nil
}

// MustSet is like [Set] but panics on any error.
func MustSet[T any](obj any, fieldPath string, val T) {
	if err := Set(obj, fieldPath, val); err != nil {
		panic("saferefl.MustSet: " + err.Error())
	}
}
