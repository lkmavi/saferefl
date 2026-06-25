package saferefl

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	"github.com/lkmavi/saferefl/internal/typeinfo"
)

// structPtrOf validates that obj is a non-nil pointer to a struct, resolves its
// TypeDescriptor (from cache on the hot path, building it on first call), and returns both.
// Used by EachField, CopyFields, ToMap, etc. — these are not on the Get/Set hot path.
func structPtrOf(obj any) (*typeinfo.TypeDescriptor, unsafe.Pointer, error) {
	e := (*eface)(unsafe.Pointer(&obj))
	if e._typ == nil {
		return nil, nil, fmt.Errorf("saferefl: obj must not be nil")
	}
	if efaceKind(e._typ) != reflect.Pointer {
		return nil, nil, fmt.Errorf("saferefl: obj must be a pointer to struct, got %v", efaceKind(e._typ))
	}
	if e.data == nil {
		return nil, nil, fmt.Errorf("saferefl: obj pointer must not be nil")
	}
	if desc, ok := typeinfo.PtrCacheLoad(uintptr(e._typ)); ok {
		return desc, e.data, nil
	}
	elem := reflect.TypeOf(obj).Elem()
	if elem.Kind() != reflect.Struct {
		return nil, nil, fmt.Errorf("saferefl: obj must point to a struct, got pointer to %v", elem.Kind())
	}
	desc := typeinfo.TypeDescriptorOf(elem)
	typeinfo.PtrCacheStore(uintptr(e._typ), desc)
	return desc, e.data, nil
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

		fieldPtr := unsafe.Pointer(uintptr(currentPtr) + fm.Offset)

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
