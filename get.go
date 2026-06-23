package saferefl

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	"github.com/lkmavi/saferefl/internal/typeinfo"
)

// Get returns the value of fieldPath on obj as type T.
// fieldPath supports dot-separated paths for nested access (e.g. "Address.City").
// obj must be a non-nil pointer to a struct.
//
// Returns [TypeMismatchError] if the field's type is not assignable to T,
// or [FieldNotFoundError] if the path does not exist on the type.
func Get[T any](obj any, fieldPath string) (T, error) {
	var zero T
	if fieldPath == "" {
		return zero, fmt.Errorf("saferefl: field path must not be empty")
	}
	e := (*eface)(unsafe.Pointer(&obj))
	if e._typ == nil {
		return zero, fmt.Errorf("saferefl: obj must not be nil")
	}
	if efaceKind(e._typ) != reflect.Pointer {
		return zero, fmt.Errorf("saferefl: obj must be a pointer to struct, got %v", efaceKind(e._typ))
	}
	if e.data == nil {
		return zero, fmt.Errorf("saferefl: obj pointer must not be nil")
	}

	wantType := reflect.TypeOf((*T)(nil)).Elem()

	// Hot path: descriptor cached by pointer-type key — no t.Elem() needed.
	if desc, ok := typeinfo.PtrCacheLoad(uintptr(e._typ)); ok {
		return getWithDesc[T](e.data, desc, fieldPath, wantType)
	}

	// First call for this struct type: resolve elem type and build the descriptor.
	return getSlowPath[T](obj, e, fieldPath, wantType)
}

// getWithDesc resolves fieldPath against desc and returns the typed value.
// Inlined by the caller when the descriptor is already cached (hot path).
func getWithDesc[T any](objPtr unsafe.Pointer, desc *typeinfo.TypeDescriptor, fieldPath string, wantType reflect.Type) (T, error) {
	var zero T

	// Fast path: single-segment field (no dot) — skip resolvePath entirely.
	if strings.IndexByte(fieldPath, '.') < 0 {
		fm, ok := desc.FieldsByName[fieldPath]
		if !ok {
			return zero, &FieldNotFoundError{Type: desc.Type.String(), FieldPath: fieldPath}
		}
		if !fm.Type.AssignableTo(wantType) {
			return zero, &TypeMismatchError{FieldPath: fieldPath, FieldType: fm.Type.String(), WantType: wantType.String()}
		}
		fieldPtr := unsafe.Pointer(uintptr(objPtr) + fm.Offset)
		if fm.Type == wantType {
			return *(*T)(fieldPtr), nil
		}
		return reflect.NewAt(fm.Type, fieldPtr).Elem().Interface().(T), nil
	}

	// Dot-path: delegate to the full resolver.
	ptr, fm, err := resolvePath(objPtr, desc.Type, fieldPath)
	if err != nil {
		return zero, err
	}
	if !fm.Type.AssignableTo(wantType) {
		return zero, &TypeMismatchError{FieldPath: fieldPath, FieldType: fm.Type.String(), WantType: wantType.String()}
	}
	if fm.Type == wantType {
		return *(*T)(ptr), nil
	}
	return reflect.NewAt(fm.Type, ptr).Elem().Interface().(T), nil
}

// getSlowPath handles the first call for a given struct type: it resolves the elem
// reflect.Type, builds the TypeDescriptor, populates the ptr cache, then delegates.
//
//go:noinline
func getSlowPath[T any](obj any, e *eface, fieldPath string, wantType reflect.Type) (T, error) {
	var zero T
	elem := reflect.TypeOf(obj).Elem()
	if elem.Kind() != reflect.Struct {
		return zero, fmt.Errorf("saferefl: obj must point to a struct, got pointer to %v", elem.Kind())
	}
	desc := typeinfo.TypeDescriptorOf(elem)
	typeinfo.PtrCacheStore(uintptr(e._typ), desc)
	return getWithDesc[T](e.data, desc, fieldPath, wantType)
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
