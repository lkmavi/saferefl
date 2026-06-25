package saferefl

import (
	"reflect"
	"sync"
	"unsafe"

	"github.com/lkmavi/saferefl/internal/typeinfo"
)

// CopyFields copies exported fields from src to dst where field names match
// and types are assignable or convertible. Fields present only in src or only
// in dst are silently skipped. Both src and dst must be non-nil pointers to structs.
//
// The copy plan (which fields match and whether conversion is needed) is computed
// once per (srcType, dstType) pair and cached; repeated calls with the same types
// pay only the slice-iteration + reflect.Value.Set cost.
func CopyFields(src, dst any) error {
	srcDesc, srcPtr, err := structPtrOf(src)
	if err != nil {
		return err
	}
	dstDesc, dstPtr, err := structPtrOf(dst)
	if err != nil {
		return err
	}

	plan := loadCopyPlan(srcDesc, dstDesc)
	for i := range plan {
		e := &plan[i]
		srcFieldPtr := unsafe.Pointer(uintptr(srcPtr) + e.srcOffset) //nolint:gosec
		dstFieldPtr := unsafe.Pointer(uintptr(dstPtr) + e.dstOffset) //nolint:gosec
		srcVal := reflect.NewAt(e.srcType, srcFieldPtr).Elem()
		dstVal := reflect.NewAt(e.dstType, dstFieldPtr).Elem()
		if e.convert {
			dstVal.Set(srcVal.Convert(e.dstType))
		} else {
			dstVal.Set(srcVal)
		}
	}
	return nil
}

// copyEntry holds the pre-resolved copy instruction for one matching field pair.
type copyEntry struct {
	srcOffset uintptr
	dstOffset uintptr
	srcType   reflect.Type
	dstType   reflect.Type
	convert   bool // false = direct Set (assignable); true = Convert then Set
}

type copyPlanKey struct {
	src reflect.Type
	dst reflect.Type
}

var copyPlanCache sync.Map // copyPlanKey → []copyEntry

// loadCopyPlan returns the cached plan for (srcDesc.Type, dstDesc.Type), building it once.
func loadCopyPlan(srcDesc, dstDesc *typeinfo.TypeDescriptor) []copyEntry {
	key := copyPlanKey{srcDesc.Type, dstDesc.Type}
	if v, ok := copyPlanCache.Load(key); ok {
		return v.([]copyEntry)
	}

	// Build: iterate src exported fields, find matching dst field, resolve assignability.
	var plan []copyEntry
	for name, srcFm := range srcDesc.FieldsByName {
		if !srcFm.Exported {
			continue
		}
		dstFm, ok := dstDesc.FieldsByName[name]
		if !ok || !dstFm.Exported {
			continue
		}
		e := copyEntry{
			srcOffset: srcFm.Offset,
			dstOffset: dstFm.Offset,
			srcType:   srcFm.Type,
			dstType:   dstFm.Type,
		}
		if !srcFm.Type.AssignableTo(dstFm.Type) {
			if srcFm.Type.ConvertibleTo(dstFm.Type) {
				e.convert = true
			} else {
				continue // incompatible types: silently skip (DTO mapping pattern)
			}
		}
		plan = append(plan, e)
	}

	v, _ := copyPlanCache.LoadOrStore(key, plan)
	return v.([]copyEntry)
}
