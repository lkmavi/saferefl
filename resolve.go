package saferefl

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	"github.com/lkmavi/saferefl/internal/typeinfo"
)

// structPtrOf returns a pointer to obj's underlying struct and the struct's reflect.Type.
// obj must be a non-nil pointer to a struct.
func structPtrOf(obj any) (unsafe.Pointer, reflect.Type, error) {
	if obj == nil {
		return nil, nil, fmt.Errorf("saferefl: obj must not be nil")
	}
	rv := reflect.ValueOf(obj)
	if rv.Kind() != reflect.Pointer {
		return nil, nil, fmt.Errorf("saferefl: obj must be a pointer to struct, got %v", rv.Kind())
	}
	if rv.IsNil() {
		return nil, nil, fmt.Errorf("saferefl: obj pointer must not be nil")
	}
	elem := rv.Type().Elem()
	if elem.Kind() != reflect.Struct {
		return nil, nil, fmt.Errorf("saferefl: obj must point to a struct, got pointer to %v", elem.Kind())
	}
	return unsafe.Pointer(rv.Pointer()), elem, nil
}

// resolvePath walks fieldPath (dot-separated segments) from objPtr of structType and
// returns a pointer to the target field along with its metadata.
//
// Intermediate fields that are pointers to structs are transparently dereferenced.
// A nil intermediate pointer is returned as an error rather than causing a panic.
func resolvePath(objPtr unsafe.Pointer, structType reflect.Type, fieldPath string) (unsafe.Pointer, *typeinfo.FieldMeta, error) {
	if fieldPath == "" {
		return nil, nil, fmt.Errorf("saferefl: field path must not be empty")
	}

	currentPtr := objPtr
	currentType := structType
	remaining := fieldPath

	for {
		segment, rest, hasMore := strings.Cut(remaining, ".")

		desc := typeinfo.TypeDescriptorOf(currentType)
		fm, ok := desc.FieldsByName[segment]
		if !ok {
			return nil, nil, &FieldNotFoundError{
				Type:      currentType.String(),
				FieldPath: fieldPath,
			}
		}

		fieldPtr := typeinfo.GetFieldPtr(currentPtr, fm)

		if !hasMore {
			return fieldPtr, fm, nil
		}

		// Intermediate segment: step into the struct or through a pointer-to-struct.
		nextType := fm.Type
		nextPtr := fieldPtr

		if nextType.Kind() == reflect.Pointer {
			nextPtr = *(*unsafe.Pointer)(fieldPtr)
			if nextPtr == nil {
				return nil, nil, fmt.Errorf("saferefl: nil pointer at field %q in path %q", segment, fieldPath)
			}
			nextType = nextType.Elem()
		}

		if nextType.Kind() != reflect.Struct {
			return nil, nil, fmt.Errorf("saferefl: %q in path %q is not a struct or pointer to struct", segment, fieldPath)
		}

		currentPtr = nextPtr
		currentType = nextType
		remaining = rest
	}
}
