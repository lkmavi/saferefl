package saferefl

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	"github.com/lkmavi/saferefl/internal/typeinfo"
)

// eface is the memory layout of an empty interface (any).
// For pointer types, the data word holds the pointer value directly.
type eface struct {
	_typ unsafe.Pointer
	data unsafe.Pointer
}

// structPtrOf returns a pointer to obj's underlying struct and the struct's reflect.Type.
// obj must be a non-nil pointer to a struct.
func structPtrOf(obj any) (unsafe.Pointer, reflect.Type, error) {
	if obj == nil {
		return nil, nil, fmt.Errorf("saferefl: obj must not be nil")
	}
	// reflect.TypeOf is cheaper than reflect.ValueOf: it inspects only the type word.
	t := reflect.TypeOf(obj)
	if t.Kind() != reflect.Pointer {
		return nil, nil, fmt.Errorf("saferefl: obj must be a pointer to struct, got %v", t.Kind())
	}
	elem := t.Elem()
	if elem.Kind() != reflect.Struct {
		return nil, nil, fmt.Errorf("saferefl: obj must point to a struct, got pointer to %v", elem.Kind())
	}
	// For pointer types, the interface data word holds the pointer value itself.
	// This avoids the more expensive reflect.Value path.
	p := (*eface)(unsafe.Pointer(&obj)).data
	if p == nil {
		return nil, nil, fmt.Errorf("saferefl: obj pointer must not be nil")
	}
	return p, elem, nil
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

// accessorStep is one segment in a pre-resolved field path.
// ptr=true means: after applying offset, dereference the resulting pointer.
type accessorStep struct {
	offset uintptr
	ptr    bool
}

// resolveMetaPath walks fieldPath through type descriptors without touching
// actual memory. Returns the field type of the final segment.
// If no step requires a pointer dereference, chain is nil and offset holds
// the total byte offset from the root struct (fast path for Accessor).
func resolveMetaPath(structType reflect.Type, fieldPath string) (offset uintptr, chain []accessorStep, fieldType reflect.Type, err error) {
	if fieldPath == "" {
		return 0, nil, nil, fmt.Errorf("saferefl: field path must not be empty")
	}

	currentType := structType
	remaining := fieldPath
	var steps []accessorStep
	var accOffset uintptr

	for {
		segment, rest, hasMore := strings.Cut(remaining, ".")

		desc := typeinfo.TypeDescriptorOf(currentType)
		fm, ok := desc.FieldsByName[segment]
		if !ok {
			return 0, nil, nil, &FieldNotFoundError{Type: currentType.String(), FieldPath: fieldPath}
		}

		if !hasMore {
			if len(steps) == 0 {
				// Pure value path: single summed offset, no chain needed.
				return accOffset + fm.Offset, nil, fm.Type, nil
			}
			steps = append(steps, accessorStep{offset: accOffset + fm.Offset})
			return 0, steps, fm.Type, nil
		}

		nextType := fm.Type
		isPtr := nextType.Kind() == reflect.Pointer
		if isPtr {
			nextType = nextType.Elem()
			// Flush accumulated non-ptr offset into this pointer step.
			steps = append(steps, accessorStep{offset: accOffset + fm.Offset, ptr: true})
			accOffset = 0
		} else {
			accOffset += fm.Offset
		}

		if nextType.Kind() != reflect.Struct {
			return 0, nil, nil, fmt.Errorf("saferefl: %q in path %q is not a struct or pointer to struct", segment, fieldPath)
		}
		currentType = nextType
		remaining = rest
	}
}
