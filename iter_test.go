package saferefl_test

import (
	"testing"

	"github.com/lkmavi/saferefl"
)

// Fixtures for fallback-path (eachFieldRec / recurseEmbedded) coverage.
// These structs have no exported fields, so buildIterPlan returns nil
// and EachField/ToMap fall back to the reflect-based recursive path.

type unexportedOnly struct{ x, y int }

type hiddenInner struct{ z int }

type withValueEmbed struct{ hiddenInner } // anonymous value-embedded, all-unexported inner
type withPtrEmbed struct{ *hiddenInner }  // anonymous pointer-embedded, all-unexported inner

// --- EachField ---

func TestEachField_basic(t *testing.T) {
	p := &person{Name: "Alice", Age: 30, Score: 9.5, Active: true}
	got := map[string]any{}
	if err := saferefl.EachField(p, func(name string, val any) bool {
		got[name] = val
		return true
	}); err != nil {
		t.Fatalf("EachField: %v", err)
	}
	// private field must be excluded
	if _, ok := got["private"]; ok {
		t.Error("EachField included unexported field 'private'")
	}
	checks := map[string]any{"Name": "Alice", "Age": 30, "Score": 9.5, "Active": true}
	for k, want := range checks {
		if got[k] != want {
			t.Errorf("field %q = %v, want %v", k, got[k], want)
		}
	}
}

func TestEachField_promoted(t *testing.T) {
	e := &employee{
		person:  person{Name: "Bob", Age: 25},
		Company: "Acme",
		Office:  address{City: "Berlin"},
	}
	got := map[string]any{}
	_ = saferefl.EachField(e, func(name string, val any) bool {
		got[name] = val
		return true
	})
	// Promoted fields from embedded person must appear at the top level.
	if got["Name"] != "Bob" {
		t.Errorf("promoted Name = %v, want Bob", got["Name"])
	}
	if got["Age"] != 25 {
		t.Errorf("promoted Age = %v, want 25", got["Age"])
	}
	if got["Company"] != "Acme" {
		t.Errorf("Company = %v, want Acme", got["Company"])
	}
}

func TestEachField_stopEarly(t *testing.T) {
	p := &person{Name: "Alice", Age: 30, Score: 9.5, Active: true}
	count := 0
	_ = saferefl.EachField(p, func(_ string, _ any) bool {
		count++
		return count < 2 // stop after second field
	})
	if count != 2 {
		t.Errorf("expected early stop after 2 fields, got %d", count)
	}
}

func TestEachField_embeddedNilPtr(t *testing.T) {
	type S struct {
		*address // nil embedded pointer → its fields are skipped
		Extra    string
	}
	s := &S{Extra: "hello"}
	got := map[string]any{}
	_ = saferefl.EachField(s, func(name string, val any) bool {
		got[name] = val
		return true
	})
	if _, ok := got["City"]; ok {
		t.Error("EachField visited fields of nil embedded pointer")
	}
	if got["Extra"] != "hello" {
		t.Errorf("Extra = %v, want hello", got["Extra"])
	}
}

func TestEachField_pointerField(t *testing.T) {
	n := 42
	m := map[string]int{"a": 1}
	type S struct {
		Ptr *int
		M   map[string]int
		Age int
	}
	s := &S{Ptr: &n, M: m, Age: 7}
	got := map[string]any{}
	_ = saferefl.EachField(s, func(name string, val any) bool {
		got[name] = val
		return true
	})
	if gotPtr, ok := got["Ptr"].(*int); !ok || gotPtr != &n || *gotPtr != 42 {
		t.Errorf("Ptr field: got %v (%T), want *int pointing to 42", got["Ptr"], got["Ptr"])
	}
	if gotMap, ok := got["M"].(map[string]int); !ok || gotMap["a"] != 1 {
		t.Errorf("M field: got %v (%T), want map[string]int{a:1}", got["M"], got["M"])
	}
	if got["Age"] != 7 {
		t.Errorf("Age = %v, want 7", got["Age"])
	}
}

func TestEachField_nilEmbeddedPtrSkipped_nonNilVisited(t *testing.T) {
	type S struct {
		*address
		Extra string
	}
	addr := &address{City: "Paris", Country: "FR"}
	s := &S{address: addr, Extra: "x"}
	got := map[string]any{}
	_ = saferefl.EachField(s, func(name string, val any) bool {
		got[name] = val
		return true
	})
	if got["City"] != "Paris" {
		t.Errorf("City = %v, want Paris", got["City"])
	}
	if got["Country"] != "FR" {
		t.Errorf("Country = %v, want FR", got["Country"])
	}
}

func TestEachField_nilObj(t *testing.T) {
	if err := saferefl.EachField(nil, func(_ string, _ any) bool { return true }); err == nil {
		t.Error("expected error for nil obj, got nil")
	}
}

func TestEachField_nonPtr(t *testing.T) {
	if err := saferefl.EachField(person{}, func(_ string, _ any) bool { return true }); err == nil {
		t.Error("expected error for non-pointer obj, got nil")
	}
}

// --- MapForEach ---

func TestMapForEach_basic(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	got := map[string]int{}
	saferefl.MapForEach(m, func(k string, v int) bool {
		got[k] = v
		return true
	})
	for k, want := range m {
		if got[k] != want {
			t.Errorf("key %q: got %d, want %d", k, got[k], want)
		}
	}
}

func TestMapForEach_stopEarly(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	count := 0
	saferefl.MapForEach(m, func(_ string, _ int) bool {
		count++
		return count < 2
	})
	if count != 2 {
		t.Errorf("expected early stop after 2, got %d", count)
	}
}

// --- eachFieldRec / recurseEmbedded fallback path ---

func TestEachField_fallback_noExportedFields(t *testing.T) {
	s := &unexportedOnly{x: 1, y: 2}
	var count int
	if err := saferefl.EachField(s, func(_ string, _ any) bool { count++; return true }); err != nil {
		t.Fatalf("EachField: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 fields for all-unexported struct, got %d", count)
	}
}

func TestEachField_fallback_valueEmbedUnexported(t *testing.T) {
	// Anonymous value-embedded struct with no exported fields → recurseEmbedded reflect.Struct case.
	s := &withValueEmbed{hiddenInner: hiddenInner{z: 1}}
	var count int
	if err := saferefl.EachField(s, func(_ string, _ any) bool { count++; return true }); err != nil {
		t.Fatalf("EachField: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 fields, got %d", count)
	}
}

func TestEachField_fallback_ptrEmbedNilUnexported(t *testing.T) {
	// Anonymous pointer-embedded nil → recurseEmbedded reflect.Pointer nil case.
	s := &withPtrEmbed{}
	var count int
	if err := saferefl.EachField(s, func(_ string, _ any) bool { count++; return true }); err != nil {
		t.Fatalf("EachField: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 fields, got %d", count)
	}
}

func TestEachField_fallback_ptrEmbedNonNilUnexported(t *testing.T) {
	// Anonymous pointer-embedded non-nil → recurseEmbedded reflect.Pointer non-nil case.
	s := &withPtrEmbed{hiddenInner: &hiddenInner{z: 5}}
	var count int
	if err := saferefl.EachField(s, func(_ string, _ any) bool { count++; return true }); err != nil {
		t.Fatalf("EachField: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 fields, got %d", count)
	}
}

func TestMapForEach_nilMap(t *testing.T) {
	var m map[string]int
	count := 0
	saferefl.MapForEach(m, func(_ string, _ int) bool { count++; return true })
	if count != 0 {
		t.Errorf("nil map: expected 0 iterations, got %d", count)
	}
}
