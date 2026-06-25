package saferefl_test

import (
	"errors"
	"testing"

	"github.com/lkmavi/saferefl"
)

// shared fixture for conv tests
type convUser struct {
	Name   string  `json:"name"   db:"user_name"`
	Age    int     `json:"age"`
	Score  float64 `json:"-"` // skipped in ToMapByTag
	Hidden string  `json:""`  // empty tag name → skipped
}

// --- ToMap ---

func TestToMap_basic(t *testing.T) {
	u := &convUser{Name: "Alice", Age: 30, Score: 9.5}
	m, err := saferefl.ToMap(u)
	if err != nil {
		t.Fatalf("ToMap: %v", err)
	}
	if m["Name"] != "Alice" {
		t.Errorf("Name = %v, want Alice", m["Name"])
	}
	if m["Age"] != 30 {
		t.Errorf("Age = %v, want 30", m["Age"])
	}
	if m["Score"] != 9.5 {
		t.Errorf("Score = %v, want 9.5", m["Score"])
	}
}

func TestToMap_promoted(t *testing.T) {
	e := &employee{
		person:  person{Name: "Dave"},
		Company: "Corp",
	}
	m, err := saferefl.ToMap(e)
	if err != nil {
		t.Fatalf("ToMap: %v", err)
	}
	if m["Name"] != "Dave" {
		t.Errorf("promoted Name = %v, want Dave", m["Name"])
	}
	if m["Company"] != "Corp" {
		t.Errorf("Company = %v, want Corp", m["Company"])
	}
}

func TestToMap_pointerField(t *testing.T) {
	n := 99
	type S struct {
		Ptr *int
		Age int
	}
	s := &S{Ptr: &n, Age: 5}
	m, err := saferefl.ToMap(s)
	if err != nil {
		t.Fatalf("ToMap: %v", err)
	}
	gotPtr, ok := m["Ptr"].(*int)
	if !ok || gotPtr != &n || *gotPtr != 99 {
		t.Errorf("Ptr field: got %v (%T), want *int → 99", m["Ptr"], m["Ptr"])
	}
	if m["Age"] != 5 {
		t.Errorf("Age = %v, want 5", m["Age"])
	}
}

func TestToMap_nilObj(t *testing.T) {
	if _, err := saferefl.ToMap(nil); err == nil {
		t.Error("expected error for nil obj")
	}
}

// --- ToMapByTag ---

func TestToMapByTag_basic(t *testing.T) {
	u := &convUser{Name: "Alice", Age: 30, Score: 9.5}
	m, err := saferefl.ToMapByTag(u, "json")
	if err != nil {
		t.Fatalf("ToMapByTag: %v", err)
	}
	if m["name"] != "Alice" {
		t.Errorf("name = %v, want Alice", m["name"])
	}
	if m["age"] != 30 {
		t.Errorf("age = %v, want 30", m["age"])
	}
	// Score has json:"-" → must be absent
	if _, ok := m["Score"]; ok {
		t.Error("Score with json:\"-\" should be skipped")
	}
	// Hidden has json:"" (empty name) → must be absent
	if _, ok := m[""]; ok {
		t.Error("field with empty json name should be skipped")
	}
}

func TestToMapByTag_omitempty_key_stripped(t *testing.T) {
	type S struct {
		Value int `json:"value,omitempty"`
	}
	s := &S{Value: 42}
	m, _ := saferefl.ToMapByTag(s, "json")
	if m["value"] != 42 {
		t.Errorf("value = %v, want 42 (omitempty must be stripped from key)", m["value"])
	}
}

func TestToMapByTag_alternate_key(t *testing.T) {
	u := &convUser{Name: "Bob", Age: 20}
	m, err := saferefl.ToMapByTag(u, "db")
	if err != nil {
		t.Fatalf("ToMapByTag: %v", err)
	}
	if m["user_name"] != "Bob" {
		t.Errorf("user_name = %v, want Bob", m["user_name"])
	}
}

func TestToMapByTag_nilObj(t *testing.T) {
	if _, err := saferefl.ToMapByTag(nil, "json"); err == nil {
		t.Error("expected error for nil obj")
	}
}

// --- FromMap ---

func TestFromMap_basic(t *testing.T) {
	m := map[string]any{"Name": "Eve", "Age": 28, "Score": 7.7}
	dst := &convUser{}
	if err := saferefl.FromMap(m, dst); err != nil {
		t.Fatalf("FromMap: %v", err)
	}
	if dst.Name != "Eve" || dst.Age != 28 || dst.Score != 7.7 {
		t.Errorf("dst = %+v", dst)
	}
}

func TestFromMap_unknownKeys_skipped(t *testing.T) {
	m := map[string]any{"Name": "Eve", "NoSuchField": 99}
	dst := &convUser{}
	if err := saferefl.FromMap(m, dst); err != nil {
		t.Fatalf("FromMap with unknown key: %v", err)
	}
	if dst.Name != "Eve" {
		t.Errorf("Name = %v, want Eve", dst.Name)
	}
}

func TestFromMap_nilValues_skipped(t *testing.T) {
	m := map[string]any{"Name": nil, "Age": 5}
	dst := &convUser{Name: "original"}
	if err := saferefl.FromMap(m, dst); err != nil {
		t.Fatalf("FromMap: %v", err)
	}
	if dst.Name != "original" {
		t.Error("nil value should not overwrite existing field")
	}
	if dst.Age != 5 {
		t.Errorf("Age = %d, want 5", dst.Age)
	}
}

func TestFromMap_typeConversion_float64ToInt(t *testing.T) {
	// JSON unmarshal produces float64 for all numbers — FromMap must convert.
	m := map[string]any{"Age": float64(42)}
	dst := &convUser{}
	if err := saferefl.FromMap(m, dst); err != nil {
		t.Fatalf("FromMap float64→int: %v", err)
	}
	if dst.Age != 42 {
		t.Errorf("Age = %d, want 42", dst.Age)
	}
}

func TestFromMap_typeMismatch_error(t *testing.T) {
	m := map[string]any{"Name": []int{1, 2, 3}} // []int not assignable/convertible to string
	dst := &convUser{}
	err := saferefl.FromMap(m, dst)
	if err == nil {
		t.Fatal("expected TypeMismatchError, got nil")
	}
	if !errors.Is(err, saferefl.ErrTypeMismatch) {
		t.Errorf("expected ErrTypeMismatch, got %T: %v", err, err)
	}
}

func TestFromMap_nilDst(t *testing.T) {
	if err := saferefl.FromMap(map[string]any{}, nil); err == nil {
		t.Error("expected error for nil dst")
	}
}

// --- Pointer-embedded struct: EmbedChain coverage for flatToMap / flatToMapByTag ---

type ptrEmbedOuter struct {
	*address // pointer-embedded with exported fields → EmbedChain in IterPlan
	Name     string
}

func TestToMap_ptrEmbedNil(t *testing.T) {
	s := &ptrEmbedOuter{Name: "Eve"} // address pointer is nil
	m, err := saferefl.ToMap(s)
	if err != nil {
		t.Fatalf("ToMap: %v", err)
	}
	if m["Name"] != "Eve" {
		t.Errorf("Name = %v, want Eve", m["Name"])
	}
	if _, ok := m["City"]; ok {
		t.Error("City should be absent when embedded pointer is nil")
	}
}

func TestToMap_ptrEmbedNonNil(t *testing.T) {
	s := &ptrEmbedOuter{Name: "Eve", address: &address{City: "Berlin"}}
	m, err := saferefl.ToMap(s)
	if err != nil {
		t.Fatalf("ToMap: %v", err)
	}
	if m["City"] != "Berlin" {
		t.Errorf("City = %v, want Berlin", m["City"])
	}
}

type taggedPtrEmbed struct {
	*convUser        // pointer-embedded tagged struct
	Extra     string `json:"extra"`
}

func TestToMapByTag_ptrEmbedNil(t *testing.T) {
	s := &taggedPtrEmbed{Extra: "x"} // convUser pointer is nil
	m, err := saferefl.ToMapByTag(s, "json")
	if err != nil {
		t.Fatalf("ToMapByTag: %v", err)
	}
	if m["extra"] != "x" {
		t.Errorf("extra = %v, want x", m["extra"])
	}
	if _, ok := m["name"]; ok {
		t.Error("name should be absent when embedded pointer is nil")
	}
}

func TestToMapByTag_ptrEmbedNonNil(t *testing.T) {
	s := &taggedPtrEmbed{convUser: &convUser{Name: "Eve"}, Extra: "x"}
	m, err := saferefl.ToMapByTag(s, "json")
	if err != nil {
		t.Fatalf("ToMapByTag: %v", err)
	}
	if m["name"] != "Eve" {
		t.Errorf("name = %v, want Eve", m["name"])
	}
}

// --- toMapByTagRec fallback path (nil IterPlan, no exported fields) ---

func TestToMapByTag_fallback_noExported(t *testing.T) {
	m, err := saferefl.ToMapByTag(&withValueEmbed{}, "json")
	if err != nil {
		t.Fatalf("ToMapByTag: %v", err)
	}
	if len(m) != 0 {
		t.Errorf("expected empty map, got %v", m)
	}
}

func TestToMapByTag_fallback_ptrEmbedNil(t *testing.T) {
	m, err := saferefl.ToMapByTag(&withPtrEmbed{}, "json")
	if err != nil {
		t.Fatalf("ToMapByTag: %v", err)
	}
	if len(m) != 0 {
		t.Errorf("expected empty map, got %v", m)
	}
}

func TestToMapByTag_fallback_ptrEmbedNonNil(t *testing.T) {
	m, err := saferefl.ToMapByTag(&withPtrEmbed{hiddenInner: &hiddenInner{z: 7}}, "json")
	if err != nil {
		t.Fatalf("ToMapByTag: %v", err)
	}
	if len(m) != 0 {
		t.Errorf("expected empty map, got %v", m)
	}
}

func TestToMap_FromMap_roundtrip(t *testing.T) {
	src := &convUser{Name: "Frank", Age: 33, Score: 8.0}
	m, err := saferefl.ToMap(src)
	if err != nil {
		t.Fatalf("ToMap: %v", err)
	}
	dst := &convUser{}
	if err = saferefl.FromMap(m, dst); err != nil {
		t.Fatalf("FromMap: %v", err)
	}
	if *src != *dst {
		t.Errorf("roundtrip: got %+v, want %+v", dst, src)
	}
}
