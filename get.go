package saferefl

import "reflect"

// Get returns the value of fieldPath on obj as type T.
// fieldPath supports dot-separated paths for nested access (e.g. "Address.City").
// obj must be a non-nil pointer to a struct.
//
// Returns [TypeMismatchError] if the field's type is not assignable to T,
// or [FieldNotFoundError] if the path does not exist on the type.
func Get[T any](obj any, fieldPath string) (T, error) {
	var zero T
	objPtr, structType, err := structPtrOf(obj)
	if err != nil {
		return zero, err
	}
	ptr, fm, err := resolvePath(objPtr, structType, fieldPath)
	if err != nil {
		return zero, err
	}

	wantType := reflect.TypeOf((*T)(nil)).Elem()
	if !fm.Type.AssignableTo(wantType) {
		return zero, &TypeMismatchError{
			FieldPath: fieldPath,
			FieldType: fm.Type.String(),
			WantType:  wantType.String(),
		}
	}

	// Fast path: identical types — zero-copy dereference, no boxing.
	// Safety: ptr was obtained via reflect-verified offset arithmetic on a live object.
	if fm.Type == wantType {
		return *(*T)(ptr), nil
	}
	// Slow path: field type implements interface T (e.g. Get[io.Reader] on a *os.File field).
	return reflect.NewAt(fm.Type, ptr).Elem().Interface().(T), nil
}

// MustGet is like [Get] but panics on any error.
// Use it when fieldPath is statically known to be valid and the type is correct.
func MustGet[T any](obj any, fieldPath string) T {
	v, err := Get[T](obj, fieldPath)
	if err != nil {
		panic("saferefl.MustGet: " + err.Error())
	}
	return v
}
