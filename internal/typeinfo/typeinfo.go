package typeinfo

import "reflect"

// TypeDescriptor holds precomputed metadata for a struct type.
// Built once per type and cached; safe for concurrent reads after construction.
type TypeDescriptor struct {
	Type reflect.Type
	Kind reflect.Kind
	Size uintptr

	// Fields contains direct struct fields in declaration order.
	Fields []FieldMeta

	// FieldsByName contains all accessible fields including promoted fields
	// from embedded structs. Outer fields shadow inner ones (Go promotion rules).
	FieldsByName map[string]*FieldMeta

	// FieldsByTag maps tag key → tag name → field, for ORM/JSON-style lookups.
	// Tag names are the first comma-separated component (e.g. "name" from `json:"name,omitempty"`).
	FieldsByTag map[string]map[string]*FieldMeta
}

// FieldMeta holds precomputed metadata for a single struct field.
type FieldMeta struct {
	Name      string
	Index     int     // position within the declaring struct (not the root struct for promoted fields)
	Offset    uintptr // byte offset from the start of the root containing struct
	Type      reflect.Type
	Kind      reflect.Kind
	Tag       reflect.StructTag
	Anonymous bool
	Exported  bool
}
