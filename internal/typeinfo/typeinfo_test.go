package typeinfo

import (
	"reflect"
	"testing"
	"unsafe"
)

// --- test fixtures ---

type basicStruct struct {
	Name  string
	Age   int
	Score float64
}

type innerStruct struct {
	City string
	Zip  int
}

type outerStruct struct {
	innerStruct
	Name string
}

type taggedStruct struct {
	Name  string `json:"name" db:"name_col"`
	Age   int    `json:"age,omitempty"`
	Email string `json:"email" db:"email_col"`
	score float64
}

type rwStruct struct {
	ID   int64
	Name string
	Rate float64
}

// --- TypeDescriptor tests ---

func TestBuildDescriptor_basicStruct(t *testing.T) {
	rt := reflect.TypeOf(basicStruct{})
	desc := buildDescriptor(rt)

	if desc.Kind != reflect.Struct {
		t.Errorf("Kind = %v, want Struct", desc.Kind)
	}
	if desc.Type != rt {
		t.Error("Type mismatch")
	}
	if desc.Size != rt.Size() {
		t.Errorf("Size = %d, want %d", desc.Size, rt.Size())
	}
	if len(desc.Fields) != rt.NumField() {
		t.Errorf("Fields len = %d, want %d", len(desc.Fields), rt.NumField())
	}

	for i := range rt.NumField() {
		sf := rt.Field(i)
		fm := desc.Fields[i]
		if fm.Name != sf.Name {
			t.Errorf("Fields[%d].Name = %q, want %q", i, fm.Name, sf.Name)
		}
		if fm.Offset != sf.Offset {
			t.Errorf("Fields[%d].Offset = %d, want %d", i, fm.Offset, sf.Offset)
		}
		if fm.Type != sf.Type {
			t.Errorf("Fields[%d].Type mismatch", i)
		}
		if fm.Kind != sf.Type.Kind() {
			t.Errorf("Fields[%d].Kind = %v, want %v", i, fm.Kind, sf.Type.Kind())
		}
		if fm.Exported != sf.IsExported() {
			t.Errorf("Fields[%d].Exported = %v, want %v", i, fm.Exported, sf.IsExported())
		}
	}
}

func TestBuildDescriptor_embeddedFields(t *testing.T) {
	outerType := reflect.TypeOf(outerStruct{})
	innerType := reflect.TypeOf(innerStruct{})
	desc := buildDescriptor(outerType)

	// Direct field: offset matches reflect.
	nameFM, ok := desc.FieldsByName["Name"]
	if !ok {
		t.Fatal("FieldsByName missing 'Name'")
	}
	wantNameOffset := mustField(outerType, "Name").Offset
	if nameFM.Offset != wantNameOffset {
		t.Errorf("Name.Offset = %d, want %d", nameFM.Offset, wantNameOffset)
	}

	// Promoted field: offset = embeddingOffset + fieldOffset.
	cityFM, ok := desc.FieldsByName["City"]
	if !ok {
		t.Fatal("FieldsByName missing promoted 'City'")
	}
	innerOffset := mustField(outerType, "innerStruct").Offset
	cityOffset := mustField(innerType, "City").Offset
	if cityFM.Offset != innerOffset+cityOffset {
		t.Errorf("City.Offset = %d, want %d", cityFM.Offset, innerOffset+cityOffset)
	}

	zipFM, ok := desc.FieldsByName["Zip"]
	if !ok {
		t.Fatal("FieldsByName missing promoted 'Zip'")
	}
	zipOffset := mustField(innerType, "Zip").Offset
	if zipFM.Offset != innerOffset+zipOffset {
		t.Errorf("Zip.Offset = %d, want %d", zipFM.Offset, innerOffset+zipOffset)
	}
}

func TestBuildDescriptor_tags(t *testing.T) {
	rt := reflect.TypeOf(taggedStruct{})
	desc := buildDescriptor(rt)

	if desc.FieldsByTag["json"]["name"] == nil {
		t.Fatal("FieldsByTag[json][name] missing")
	}
	if desc.FieldsByTag["json"]["name"].Name != "Name" {
		t.Errorf("json:name → field %q, want Name", desc.FieldsByTag["json"]["name"].Name)
	}

	// `json:"age,omitempty"` → name component is "age".
	if desc.FieldsByTag["json"]["age"] == nil {
		t.Fatal("FieldsByTag[json][age] missing")
	}
	if desc.FieldsByTag["json"]["age"].Name != "Age" {
		t.Errorf("json:age → field %q, want Age", desc.FieldsByTag["json"]["age"].Name)
	}

	if desc.FieldsByTag["db"]["email_col"] == nil {
		t.Fatal("FieldsByTag[db][email_col] missing")
	}
	if desc.FieldsByTag["db"]["email_col"].Name != "Email" {
		t.Errorf("db:email_col → field %q, want Email", desc.FieldsByTag["db"]["email_col"].Name)
	}
}

func TestBuildDescriptor_unexportedField(t *testing.T) {
	rt := reflect.TypeOf(taggedStruct{})
	desc := buildDescriptor(rt)

	// Unexported fields appear in FieldsByName with Exported=false.
	scoreFM, ok := desc.FieldsByName["score"]
	if !ok {
		t.Fatal("FieldsByName missing 'score'")
	}
	if scoreFM.Exported {
		t.Error("'score' should have Exported=false")
	}
}

// --- GetFieldPtr / SetField tests ---

func TestGetSetField_roundtrip(t *testing.T) {
	rt := reflect.TypeOf(rwStruct{})
	desc := buildDescriptor(rt)
	s := &rwStruct{ID: 42, Name: "hello", Rate: 3.14}
	objPtr := unsafe.Pointer(s)

	idFM := desc.FieldsByName["ID"]
	if id := *(*int64)(GetFieldPtr(objPtr, idFM)); id != 42 {
		t.Errorf("GetFieldPtr ID = %d, want 42", id)
	}

	nameFM := desc.FieldsByName["Name"]
	if name := *(*string)(GetFieldPtr(objPtr, nameFM)); name != "hello" {
		t.Errorf("GetFieldPtr Name = %q, want hello", name)
	}

	rateFM := desc.FieldsByName["Rate"]
	if rate := *(*float64)(GetFieldPtr(objPtr, rateFM)); rate != 3.14 {
		t.Errorf("GetFieldPtr Rate = %v, want 3.14", rate)
	}

	if err := SetField(objPtr, idFM, reflect.ValueOf(int64(99))); err != nil {
		t.Fatalf("SetField ID: %v", err)
	}
	if s.ID != 99 {
		t.Errorf("after SetField ID = %d, want 99", s.ID)
	}

	if err := SetField(objPtr, nameFM, reflect.ValueOf("world")); err != nil {
		t.Fatalf("SetField Name: %v", err)
	}
	if s.Name != "world" {
		t.Errorf("after SetField Name = %q, want world", s.Name)
	}
}

func TestSetField_unexportedReturnsError(t *testing.T) {
	rt := reflect.TypeOf(taggedStruct{})
	desc := buildDescriptor(rt)
	s := &taggedStruct{}
	err := SetField(unsafe.Pointer(s), desc.FieldsByName["score"], reflect.ValueOf(1.0))
	if err == nil {
		t.Error("expected error setting unexported field, got nil")
	}
}

func TestSetField_typeMismatchReturnsError(t *testing.T) {
	rt := reflect.TypeOf(rwStruct{})
	desc := buildDescriptor(rt)
	s := &rwStruct{}
	// Attempt to assign string to int64 field.
	err := SetField(unsafe.Pointer(s), desc.FieldsByName["ID"], reflect.ValueOf("oops"))
	if err == nil {
		t.Error("expected type-mismatch error, got nil")
	}
}

// --- helpers ---

func mustField(t reflect.Type, name string) reflect.StructField {
	sf, ok := t.FieldByName(name)
	if !ok {
		panic("field not found: " + name)
	}
	return sf
}
