package saferefl

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/lkmavi/saferefl/internal/typeinfo"
)

// GetByTag returns the value of the struct field whose tag key matches tagValue.
// Example: GetByTag[string](&user, "json", "name") reads the field tagged `json:"name"` or `json:"name,omitempty"`.
// obj must be a non-nil pointer to a struct.
func GetByTag[T any](obj any, tagKey, tagValue string) (T, error) {
	var zero T
	e := (*eface)(unsafe.Pointer(&obj)) //nolint:gosec
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
	if desc, ok := typeinfo.PtrCacheLoad(uintptr(e._typ)); ok {
		return getByTagWithDesc[T](e.data, desc, tagKey, tagValue, wantType)
	}
	return getByTagSlowPath[T](obj, e, tagKey, tagValue, wantType)
}

func getByTagWithDesc[T any](objPtr unsafe.Pointer, desc *typeinfo.TypeDescriptor, tagKey, tagValue string, wantType reflect.Type) (T, error) {
	var zero T
	fm, err := lookupByTag(desc, tagKey, tagValue)
	if err != nil {
		return zero, err
	}
	if !fm.Type.AssignableTo(wantType) {
		return zero, &TypeMismatchError{FieldPath: tagPath(tagKey, tagValue), FieldType: fm.Type.String(), WantType: wantType.String()}
	}
	fieldPtr := unsafe.Pointer(uintptr(objPtr) + fm.Offset) //nolint:gosec
	if fm.Type == wantType {
		return *(*T)(fieldPtr), nil
	}
	val, ok := reflect.NewAt(fm.Type, fieldPtr).Elem().Interface().(T)
	if !ok {
		return zero, &TypeMismatchError{FieldPath: tagPath(tagKey, tagValue), FieldType: fm.Type.String(), WantType: wantType.String()}
	}
	return val, nil
}

//go:noinline
func getByTagSlowPath[T any](obj any, e *eface, tagKey, tagValue string, wantType reflect.Type) (T, error) {
	var zero T
	elem := reflect.TypeOf(obj).Elem()
	if elem.Kind() != reflect.Struct {
		return zero, fmt.Errorf("saferefl: obj must point to a struct, got pointer to %v", elem.Kind())
	}
	desc := typeinfo.TypeDescriptorOf(elem)
	typeinfo.PtrCacheStore(uintptr(e._typ), desc)
	return getByTagWithDesc[T](e.data, desc, tagKey, tagValue, wantType)
}

// SetByTag sets the value of the struct field whose tag key matches tagValue.
// obj must be a non-nil pointer to a struct.
func SetByTag[T any](obj any, tagKey, tagValue string, val T) error {
	e := (*eface)(unsafe.Pointer(&obj)) //nolint:gosec
	if e._typ == nil {
		return fmt.Errorf("saferefl: obj must not be nil")
	}
	if efaceKind(e._typ) != reflect.Pointer {
		return fmt.Errorf("saferefl: obj must be a pointer to struct, got %v", efaceKind(e._typ))
	}
	if e.data == nil {
		return fmt.Errorf("saferefl: obj pointer must not be nil")
	}
	wantType := reflect.TypeOf((*T)(nil)).Elem()
	if desc, ok := typeinfo.PtrCacheLoad(uintptr(e._typ)); ok {
		return setByTagWithDesc[T](e.data, desc, tagKey, tagValue, wantType, val)
	}
	return setByTagSlowPath[T](obj, e, tagKey, tagValue, wantType, val)
}

func setByTagWithDesc[T any](objPtr unsafe.Pointer, desc *typeinfo.TypeDescriptor, tagKey, tagValue string, wantType reflect.Type, val T) error {
	fm, err := lookupByTag(desc, tagKey, tagValue)
	if err != nil {
		return err
	}
	if !wantType.AssignableTo(fm.Type) {
		return &TypeMismatchError{FieldPath: tagPath(tagKey, tagValue), FieldType: fm.Type.String(), WantType: wantType.String()}
	}
	fieldPtr := unsafe.Pointer(uintptr(objPtr) + fm.Offset) //nolint:gosec
	if fm.Type == wantType {
		*(*T)(fieldPtr) = val
		return nil
	}
	setSlowPath(fieldPtr, fm.Type, wantType, val)
	return nil
}

//go:noinline
func setByTagSlowPath[T any](obj any, e *eface, tagKey, tagValue string, wantType reflect.Type, val T) error {
	elem := reflect.TypeOf(obj).Elem()
	if elem.Kind() != reflect.Struct {
		return fmt.Errorf("saferefl: obj must point to a struct, got pointer to %v", elem.Kind())
	}
	desc := typeinfo.TypeDescriptorOf(elem)
	typeinfo.PtrCacheStore(uintptr(e._typ), desc)
	return setByTagWithDesc[T](e.data, desc, tagKey, tagValue, wantType, val)
}

// lookupByTag resolves a FieldMeta by tag key and value from a TypeDescriptor.
func lookupByTag(desc *typeinfo.TypeDescriptor, tagKey, tagValue string) (*typeinfo.FieldMeta, error) {
	tags, ok := desc.FieldsByTag[tagKey]
	if !ok {
		return nil, &FieldNotFoundError{Type: desc.Type.String(), FieldPath: tagPath(tagKey, tagValue)}
	}
	fm, ok := tags[tagValue]
	if !ok {
		return nil, &FieldNotFoundError{Type: desc.Type.String(), FieldPath: tagPath(tagKey, tagValue)}
	}
	return fm, nil
}

func tagPath(key, value string) string { return key + `:"` + value + `"` }
