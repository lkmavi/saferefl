//go:build reflectx_strict

package saferefl

import (
	"unsafe"

	"github.com/lkmavi/saferefl/internal/typeinfo"
)

// fieldAny is unreachable in reflectx_strict because IterPlan is always nil,
// so iterFlat is never called. Panics to surface any future regression.
func fieldAny(_ *typeinfo.IterEntry, _ unsafe.Pointer) any {
	panic("saferefl: fieldAny reached in reflectx_strict build — IterPlan must be nil")
}
