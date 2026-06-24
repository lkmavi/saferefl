package saferefl

import (
	"reflect"
	"unsafe"

	"github.com/lkmavi/saferefl/internal/typeinfo"
)

// EachField calls fn for every exported field of obj in declaration order,
// including fields promoted from embedded structs. Recursion into embedded
// value-struct fields is transparent; embedded pointer-struct fields are
// dereferenced and skipped if nil. Return false from fn to stop early.
// obj must be a non-nil pointer to a struct.
//
// On the non-strict build (default), the val passed to fn aliases the struct's
// field memory — type-asserting copies the value out, so reads are safe.
// Do not write to the struct concurrently with the callback.
func EachField(obj any, fn func(name string, val any) bool) error {
	desc, objPtr, err := structPtrOf(obj)
	if err != nil {
		return err
	}
	if len(desc.IterPlan) > 0 {
		iterFlat(objPtr, desc.IterPlan, fn)
	} else {
		eachFieldRec(objPtr, desc, 0, fn)
	}
	return nil
}

// iterFlat is the fast path for EachField/ToMap: walks the pre-computed IterPlan
// in order, following EmbedChain pointer-deref steps per entry when needed.
// Returns false if fn stopped early, true otherwise.
func iterFlat(objPtr unsafe.Pointer, plan []typeinfo.IterEntry, fn func(string, any) bool) bool {
	for i := range plan {
		e := &plan[i]
		base := objPtr

		// Follow the chain of embedded *Struct pointer fields (usually 0 steps).
		skipped := false
		for _, chainOff := range e.EmbedChain {
			base = *(*unsafe.Pointer)(unsafe.Pointer(uintptr(base) + chainOff)) //nolint:gosec
			if base == nil {
				skipped = true
				break
			}
		}
		if skipped {
			continue
		}

		fieldPtr := unsafe.Pointer(uintptr(base) + e.Offset) //nolint:gosec
		if !fn(e.Name, fieldAny(e.AbiType, e.IfaceDirect, fieldPtr)) {
			return false
		}
	}
	return true
}

// eachFieldRec is the fallback path used in reflectx_strict mode (IterPlan is nil).
// It recurses into embedded structs at runtime, building interface values via reflect.
func eachFieldRec(objPtr unsafe.Pointer, desc *typeinfo.TypeDescriptor, baseOffset uintptr, fn func(string, any) bool) bool {
	for i := range desc.Fields {
		fm := &desc.Fields[i]
		absOffset := baseOffset + fm.Offset

		if fm.Anonymous {
			switch fm.Kind {
			case reflect.Struct:
				sub := typeinfo.TypeDescriptorOf(fm.Type)
				if !eachFieldRec(objPtr, sub, absOffset, fn) {
					return false
				}
				continue
			case reflect.Pointer:
				if fm.Type.Elem().Kind() == reflect.Struct {
					inner := *(*unsafe.Pointer)(unsafe.Pointer(uintptr(objPtr) + absOffset)) //nolint:gosec
					if inner == nil {
						continue
					}
					sub := typeinfo.TypeDescriptorOf(fm.Type.Elem())
					if !eachFieldRec(inner, sub, 0, fn) {
						return false
					}
					continue
				}
			}
		}

		if !fm.Exported {
			continue
		}

		fieldPtr := unsafe.Pointer(uintptr(objPtr) + absOffset) //nolint:gosec
		val := reflect.NewAt(fm.Type, fieldPtr).Elem().Interface()
		if !fn(fm.Name, val) {
			return false
		}
	}
	return true
}

// MapForEach calls fn for each key-value pair in m, stopping early if fn returns false.
// Performance is identical to a plain range loop — MapForEach IS range with a callback.
// The ~1 ns difference vs bare range in benchmarks is function-call overhead, which is
// within measurement noise for any non-trivial map size.
func MapForEach[K comparable, V any](m map[K]V, fn func(K, V) bool) {
	for k, v := range m {
		if !fn(k, v) {
			return
		}
	}
}
