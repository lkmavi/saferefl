package typeinfo

import (
	"fmt"
	"reflect"
	"unsafe"
)

// GetFieldPtr returns a pointer to the field within the struct pointed to by objPtr.
// Adding a reflect-verified offset to a pointer of a live object is explicitly permitted
// by the unsafe.Pointer conversion rules.
func GetFieldPtr(objPtr unsafe.Pointer, f *FieldMeta) unsafe.Pointer {
	return unsafe.Pointer(uintptr(objPtr) + f.Offset)
}

// SetField sets the field value using reflect.NewAt — the documented safe path for
// pointer arithmetic within a single live object.
func SetField(objPtr unsafe.Pointer, f *FieldMeta, val reflect.Value) error {
	if !f.Exported {
		return fmt.Errorf("saferefl: cannot set unexported field %q", f.Name)
	}
	if !val.Type().AssignableTo(f.Type) {
		return fmt.Errorf("saferefl: cannot assign %v to field %q of type %v", val.Type(), f.Name, f.Type)
	}
	ptr := GetFieldPtr(objPtr, f)
	reflect.NewAt(f.Type, ptr).Elem().Set(val)
	return nil
}
