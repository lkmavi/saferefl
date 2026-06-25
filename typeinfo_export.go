package saferefl

import (
	"reflect"

	"github.com/lkmavi/saferefl/internal/typeinfo"
)

// TypeDescriptor holds precomputed metadata for a struct type.
// Built once per type and cached globally; safe for concurrent reads after construction.
// Use [TypeDescriptorOf] to obtain one; do not construct directly.
type TypeDescriptor = typeinfo.TypeDescriptor

// FieldMeta holds precomputed metadata for a single struct field.
// Offset is always relative to the root struct, even for promoted fields.
type FieldMeta = typeinfo.FieldMeta

// TypeDescriptorOf returns the precomputed [TypeDescriptor] for t.
// Panics if t is not a struct type.
// The descriptor is built at most once per type and cached globally.
func TypeDescriptorOf(t reflect.Type) *TypeDescriptor {
	return typeinfo.TypeDescriptorOf(t)
}
