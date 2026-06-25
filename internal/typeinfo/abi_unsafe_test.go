//go:build !reflectx_strict

package typeinfo

import (
	"reflect"
	"testing"
)

// --- collectIter / buildIterPlan ---

type iterShadowInner struct{ X string }
type iterShadowOuter struct {
	iterShadowInner     // promotes X string — must be deduped
	X               int // direct X int — must win
}

// TestCollectIter_shadowingDedup verifies the seen-map dedup path: when an outer
// direct field has the same name as a promoted field, only the outer field appears
// in IterPlan.
func TestCollectIter_shadowingDedup(t *testing.T) {
	plan := buildIterPlan(reflect.TypeOf(iterShadowOuter{}))

	count := 0
	for _, e := range plan {
		if e.Name == "X" {
			count++
			if e.Type != reflect.TypeOf(0) {
				t.Errorf("IterPlan X type = %v, want int", e.Type)
			}
		}
	}
	if count != 1 {
		t.Errorf("IterPlan has %d entries named X, want 1", count)
	}
}

type iterPtrInner struct{ City string }
type iterPtrOuter struct {
	*iterPtrInner // pointer-embedded: should appear with non-nil EmbedChain
	Name          string
}

// TestCollectIter_ptrEmbedChain verifies the pointer-embedded struct path:
// EmbedChain must be non-nil and contain the offset of the pointer field.
func TestCollectIter_ptrEmbedChain(t *testing.T) {
	plan := buildIterPlan(reflect.TypeOf(iterPtrOuter{}))

	var cityEntry *IterEntry
	for i := range plan {
		if plan[i].Name == "City" {
			cityEntry = &plan[i]
		}
	}
	if cityEntry == nil {
		t.Fatal("IterPlan missing City entry from pointer-embedded struct")
	}
	if len(cityEntry.EmbedChain) == 0 {
		t.Error("City EmbedChain must be non-nil for pointer-embedded struct")
	}
	// EmbedChain[0] must equal the offset of the *iterPtrInner field in iterPtrOuter.
	sf, _ := reflect.TypeOf(iterPtrOuter{}).FieldByName("iterPtrInner")
	if cityEntry.EmbedChain[0] != sf.Offset {
		t.Errorf("EmbedChain[0] = %d, want %d (offset of *iterPtrInner)", cityEntry.EmbedChain[0], sf.Offset)
	}
}

// TestCollectIter_typeFieldSet verifies that IterEntry.Type is populated correctly.
func TestCollectIter_typeFieldSet(t *testing.T) {
	plan := buildIterPlan(reflect.TypeOf(basicStruct{}))
	for _, e := range plan {
		if e.Type == nil {
			t.Errorf("IterEntry %q has nil Type", e.Name)
		}
	}
}

// --- collectNamed shadowing ---

type namedShadowInner struct{ X string }
type namedShadowOuter struct {
	namedShadowInner
	X int // shadows namedShadowInner.X
}

// TestCollectNamed_shadowing verifies that the two-pass approach correctly registers
// the outer field first so promoted fields with the same name are skipped.
func TestCollectNamed_shadowing(t *testing.T) {
	desc := buildDescriptor(reflect.TypeOf(namedShadowOuter{}))

	fm, ok := desc.FieldsByName["X"]
	if !ok {
		t.Fatal("FieldsByName missing X")
	}
	if fm.Type != reflect.TypeOf(0) {
		t.Errorf("FieldsByName X type = %v, want int (outer must shadow inner)", fm.Type)
	}
}
