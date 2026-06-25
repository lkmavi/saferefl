package saferefl

import (
	"reflect"
	"strings"
	"unsafe"

	"github.com/lkmavi/saferefl/internal/typeinfo"
)

// ToMap converts obj to map[string]any using exported field names as keys.
// Promoted fields from embedded structs are included and appear at the top level.
// obj must be a non-nil pointer to a struct.
//
// On non-strict builds the map values alias the struct's field memory; type-asserting
// a value copies it out. See [EachField] for the aliasing details.
func ToMap(obj any) (map[string]any, error) {
	desc, objPtr, err := structPtrOf(obj)
	if err != nil {
		return nil, err
	}
	out := make(map[string]any, len(desc.FieldsByName))
	if len(desc.IterPlan) > 0 {
		flatToMap(objPtr, desc.IterPlan, out)
	} else {
		eachFieldRec(objPtr, desc, 0, func(name string, val any) bool {
			out[name] = val
			return true
		})
	}
	return out, nil
}

// ToMapByTag converts obj to map[string]any using the first component of the specified
// struct tag as keys. Fields without the tag or tagged "-" are skipped.
// Example: ToMapByTag(&user, "json") uses `json:"name"` → key "name".
// obj must be a non-nil pointer to a struct.
func ToMapByTag(obj any, tagKey string) (map[string]any, error) {
	desc, objPtr, err := structPtrOf(obj)
	if err != nil {
		return nil, err
	}
	out := make(map[string]any, len(desc.FieldsByName))
	if len(desc.IterPlan) > 0 {
		flatToMapByTag(objPtr, desc.IterPlan, tagKey, out)
	} else {
		toMapByTagRec(objPtr, desc, 0, tagKey, out)
	}
	return out, nil
}

// flatToMap is the 0-alloc fast path for ToMap: walks the pre-computed IterPlan,
// boxing each field value via direct eface construction.
func flatToMap(objPtr unsafe.Pointer, plan []typeinfo.IterEntry, out map[string]any) {
	for i := range plan {
		e := &plan[i]
		base := objPtr

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
		out[e.Name] = fieldAny(e.AbiType, e.IfaceDirect, fieldPtr)
	}
}

// flatToMapByTag is the fast path for ToMapByTag: same as flatToMap but
// keys come from the struct tag stored in each IterEntry.
func flatToMapByTag(objPtr unsafe.Pointer, plan []typeinfo.IterEntry, tagKey string, out map[string]any) {
	for i := range plan {
		e := &plan[i]

		raw := e.Tag.Get(tagKey)
		if raw == "" || raw == "-" {
			continue
		}
		key, _, _ := strings.Cut(raw, ",")
		if key == "" || key == "-" {
			continue
		}

		base := objPtr
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
		out[key] = fieldAny(e.AbiType, e.IfaceDirect, fieldPtr)
	}
}

// toMapByTagRec is the fallback recursive path used in reflectx_strict mode.
func toMapByTagRec(objPtr unsafe.Pointer, desc *typeinfo.TypeDescriptor, baseOffset uintptr, tagKey string, out map[string]any) {
	for i := range desc.Fields {
		fm := &desc.Fields[i]
		absOffset := baseOffset + fm.Offset

		if fm.Anonymous {
			switch fm.Kind {
			case reflect.Struct:
				sub := typeinfo.TypeDescriptorOf(fm.Type)
				toMapByTagRec(objPtr, sub, absOffset, tagKey, out)
				continue
			case reflect.Pointer:
				if fm.Type.Elem().Kind() == reflect.Struct {
					inner := *(*unsafe.Pointer)(unsafe.Pointer(uintptr(objPtr) + absOffset)) //nolint:gosec
					if inner != nil {
						sub := typeinfo.TypeDescriptorOf(fm.Type.Elem())
						toMapByTagRec(inner, sub, 0, tagKey, out)
					}
					continue
				}
			}
		}

		if !fm.Exported {
			continue
		}

		raw := fm.Tag.Get(tagKey)
		if raw == "" || raw == "-" {
			continue
		}
		key, _, _ := strings.Cut(raw, ",")
		if key == "" || key == "-" {
			continue
		}

		fieldPtr := unsafe.Pointer(uintptr(objPtr) + absOffset) //nolint:gosec
		out[key] = reflect.NewAt(fm.Type, fieldPtr).Elem().Interface()
	}
}

// FromMap sets fields of dst from m using field names as keys. Unknown keys are skipped.
// Values that are not directly assignable are converted when possible (e.g. float64 → int,
// useful after JSON unmarshaling into map[string]any). Returns [TypeMismatchError] if a
// matching field exists but the value is neither assignable nor convertible.
// dst must be a non-nil pointer to a struct.
func FromMap(m map[string]any, dst any) error {
	desc, dstPtr, err := structPtrOf(dst)
	if err != nil {
		return err
	}

	for key, val := range m {
		fm, ok := desc.FieldsByName[key]
		if !ok || !fm.Exported {
			continue
		}
		if val == nil {
			continue
		}

		srcVal := reflect.ValueOf(val)
		dstFieldPtr := unsafe.Pointer(uintptr(dstPtr) + fm.Offset) //nolint:gosec
		dstField := reflect.NewAt(fm.Type, dstFieldPtr).Elem()

		if srcVal.Type().AssignableTo(fm.Type) {
			dstField.Set(srcVal)
			continue
		}
		if srcVal.Type().ConvertibleTo(fm.Type) {
			dstField.Set(srcVal.Convert(fm.Type))
			continue
		}
		return &TypeMismatchError{
			FieldPath: key,
			FieldType: fm.Type.String(),
			WantType:  srcVal.Type().String(),
		}
	}
	return nil
}
