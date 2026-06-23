package saferefl

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	"github.com/lkmavi/saferefl/internal/typeinfo"
)

// Set sets the value of fieldPath on obj to val.
// fieldPath supports dot-separated paths (e.g. "Address.City").
// obj must be a non-nil pointer to a struct.
//
// Returns [ReadOnlyError] for unexported fields,
// [TypeMismatchError] if T is not assignable to the field's type,
// or [FieldNotFoundError] if the path does not exist.
func Set[T any](obj any, fieldPath string, val T) error {
	if fieldPath == "" {
		return fmt.Errorf("saferefl: field path must not be empty")
	}
	e := (*eface)(unsafe.Pointer(&obj))
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

	// Hot path: descriptor cached by pointer-type key — no t.Elem() needed.
	if desc, ok := typeinfo.PtrCacheLoad(uintptr(e._typ)); ok {
		return setWithDesc[T](e.data, desc, fieldPath, wantType, val)
	}

	// First call for this struct type: resolve elem type and build the descriptor.
	return setSlowPathDesc[T](obj, e, fieldPath, wantType, val)
}

// setWithDesc resolves fieldPath against desc and writes val.
func setWithDesc[T any](objPtr unsafe.Pointer, desc *typeinfo.TypeDescriptor, fieldPath string, wantType reflect.Type, val T) error {
	// Fast path: single-segment field (no dot) — skip resolvePath entirely.
	if strings.IndexByte(fieldPath, '.') < 0 {
		fm, ok := desc.FieldsByName[fieldPath]
		if !ok {
			return &FieldNotFoundError{Type: desc.Type.String(), FieldPath: fieldPath}
		}
		if !fm.Exported {
			return &ReadOnlyError{FieldPath: fieldPath}
		}
		if !wantType.AssignableTo(fm.Type) {
			return &TypeMismatchError{FieldPath: fieldPath, FieldType: fm.Type.String(), WantType: wantType.String()}
		}
		fieldPtr := unsafe.Pointer(uintptr(objPtr) + fm.Offset)
		if fm.Type == wantType {
			*(*T)(fieldPtr) = val
			return nil
		}
		setSlowPath(fieldPtr, fm.Type, wantType, val)
		return nil
	}

	// Dot-path: delegate to the full resolver.
	ptr, fm, err := resolvePath(objPtr, desc.Type, fieldPath)
	if err != nil {
		return err
	}
	if !fm.Exported {
		return &ReadOnlyError{FieldPath: fieldPath}
	}
	if !wantType.AssignableTo(fm.Type) {
		return &TypeMismatchError{FieldPath: fieldPath, FieldType: fm.Type.String(), WantType: wantType.String()}
	}
	if fm.Type == wantType {
		*(*T)(ptr) = val
		return nil
	}
	setSlowPath(ptr, fm.Type, wantType, val)
	return nil
}

// setSlowPathDesc handles the first call for a given struct type.
//
//go:noinline
func setSlowPathDesc[T any](obj any, e *eface, fieldPath string, wantType reflect.Type, val T) error {
	elem := reflect.TypeOf(obj).Elem()
	if elem.Kind() != reflect.Struct {
		return fmt.Errorf("saferefl: obj must point to a struct, got pointer to %v", elem.Kind())
	}
	desc := typeinfo.TypeDescriptorOf(elem)
	typeinfo.PtrCacheStore(uintptr(e._typ), desc)
	return setWithDesc[T](e.data, desc, fieldPath, wantType, val)
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
