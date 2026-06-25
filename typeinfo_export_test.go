package saferefl_test

import (
	"reflect"
	"testing"

	"github.com/lkmavi/saferefl"
)

func TestTypeDescriptorOf_basic(t *testing.T) {
	rt := reflect.TypeOf(person{})
	desc := saferefl.TypeDescriptorOf(rt)
	if desc == nil {
		t.Fatal("TypeDescriptorOf returned nil")
	}
	if desc.Type != rt {
		t.Errorf("desc.Type = %v, want %v", desc.Type, rt)
	}
	if desc.Kind != reflect.Struct {
		t.Errorf("desc.Kind = %v, want Struct", desc.Kind)
	}
	if desc.Size != rt.Size() {
		t.Errorf("desc.Size = %d, want %d", desc.Size, rt.Size())
	}
}

func TestTypeDescriptorOf_fieldsByName(t *testing.T) {
	desc := saferefl.TypeDescriptorOf(reflect.TypeOf(person{}))
	fm, ok := desc.FieldsByName["Name"]
	if !ok {
		t.Fatal("FieldsByName missing 'Name'")
	}
	if fm.Type != reflect.TypeOf("") {
		t.Errorf("Name field type = %v, want string", fm.Type)
	}
	if !fm.Exported {
		t.Error("Name field should be exported")
	}
	sf, _ := reflect.TypeOf(person{}).FieldByName("Name")
	if fm.Offset != sf.Offset {
		t.Errorf("Offset = %d, want %d", fm.Offset, sf.Offset)
	}
}

func TestTypeDescriptorOf_unexportedField(t *testing.T) {
	desc := saferefl.TypeDescriptorOf(reflect.TypeOf(person{}))
	fm, ok := desc.FieldsByName["private"]
	if !ok {
		t.Fatal("FieldsByName missing 'private'")
	}
	if fm.Exported {
		t.Error("'private' field should not be exported")
	}
}

func TestTypeDescriptorOf_fieldsByTag(t *testing.T) {
	type S struct {
		Name string `json:"name" db:"user_name"`
		Age  int    `json:"age"`
	}
	desc := saferefl.TypeDescriptorOf(reflect.TypeOf(S{}))

	jsonTags := desc.FieldsByTag["json"]
	if jsonTags == nil {
		t.Fatal("FieldsByTag missing 'json' key")
	}
	if jsonTags["name"] == nil {
		t.Error("FieldsByTag[json][name] is nil")
	}
	if jsonTags["age"] == nil {
		t.Error("FieldsByTag[json][age] is nil")
	}

	dbTags := desc.FieldsByTag["db"]
	if dbTags == nil || dbTags["user_name"] == nil {
		t.Error("FieldsByTag[db][user_name] missing")
	}
}

func TestTypeDescriptorOf_cachedIdentity(t *testing.T) {
	rt := reflect.TypeOf(person{})
	d1 := saferefl.TypeDescriptorOf(rt)
	d2 := saferefl.TypeDescriptorOf(rt)
	if d1 != d2 {
		t.Error("TypeDescriptorOf must return the same pointer on repeated calls")
	}
}

func TestTypeDescriptorOf_promoted_offsets(t *testing.T) {
	type base struct{ X int }
	type outer struct {
		base
		Y int
	}
	desc := saferefl.TypeDescriptorOf(reflect.TypeOf(outer{}))
	fmX, ok := desc.FieldsByName["X"]
	if !ok {
		t.Fatal("promoted field X not in FieldsByName")
	}
	rt := reflect.TypeOf(outer{})
	sf, _ := rt.FieldByName("X")
	if fmX.Offset != sf.Offset {
		t.Errorf("X offset = %d, reflect says %d", fmX.Offset, sf.Offset)
	}
}

func TestTypeDescriptorOf_panicOnNonStruct(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for non-struct type, got none")
		}
	}()
	saferefl.TypeDescriptorOf(reflect.TypeOf(42))
}
