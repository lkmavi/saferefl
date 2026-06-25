package typeinfo

import (
	"reflect"
	"unsafe"
)

// TypeDescriptor holds precomputed metadata for a struct type.
// Built once per type and cached; safe for concurrent reads after construction.
type TypeDescriptor struct {
	Type reflect.Type
	Kind reflect.Kind
	Size uintptr

	// Fields contains direct struct fields in declaration order.
	Fields []FieldMeta

	// FieldsByName contains all accessible fields including promoted fields
	// from value-embedded structs. Outer fields shadow inner ones (Go promotion rules).
	//
	// Note: fields promoted from pointer-embedded structs (e.g. type Outer struct{ *Inner })
	// are NOT included here — their accessibility depends on a runtime nil check.
	// EachField and ToMap do handle pointer-embedded structs (via IterPlan.EmbedChain),
	// but Get/Set/GetByTag/SetByTag require an explicit dot-path (e.g. "Inner.Name").
	FieldsByName map[string]*FieldMeta

	// FieldsByTag maps tag key → tag name → field, for ORM/JSON-style lookups.
	// Tag names are the first comma-separated component (e.g. "name" from `json:"name,omitempty"`).
	FieldsByTag map[string]map[string]*FieldMeta

	// IterPlan is a flat, pre-ordered list of exported fields for EachField/ToMap.
	// Embedded value-struct fields are expanded with root-relative offsets (EmbedChain nil).
	// Embedded *Struct fields are expanded with struct-relative offsets plus a chain of
	// pointer-field offsets to dereference at runtime. Nil in reflectx_strict builds.
	IterPlan []IterEntry

	// IterPlanIndex maps field name to index in IterPlan for O(1) FromMap writes.
	// Nil when IterPlan is nil.
	IterPlanIndex map[string]int
}

// IterEntry describes one exported field in the pre-computed flat iteration plan.
type IterEntry struct {
	Name        string
	Tag         reflect.StructTag // cached struct tag; used by ToMapByTag to avoid reflect.StructField lookup
	Type        reflect.Type      // reflect.Type of the field; used by fieldAny for interface-kinded fields
	AbiType     unsafe.Pointer    // raw *abi.Type for the field's concrete type; nil in reflectx_strict
	IfaceDirect bool              // true when the value fits directly in an interface data word (pointer/map/chan); determined by OR of TFlag[+20] and Kind_[+23] bit-5 for cross-version compatibility
	Offset      uintptr           // byte offset of the field from the base (last chain dereference, or root if EmbedChain is nil)
	EmbedChain  []uintptr         // nil → field is root-relative; else: ordered offsets of *Struct pointer fields to deref before adding Offset
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
