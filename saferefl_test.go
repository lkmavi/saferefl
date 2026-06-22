package saferefl_test

import (
	"errors"
	"testing"

	"github.com/lkmavi/saferefl"
)

// --- fixtures ---

type person struct {
	Name    string
	Age     int
	Score   float64
	Active  bool
	private string //nolint:unused
}

type address struct {
	City    string
	Country string
}

type employee struct {
	person  // embedded (promotes Name, Age, Score, Active)
	Company string
	Office  address  // value — dot-path without deref
	Contact *address // pointer — dot-path with deref
}

// --- Get tests ---

func TestGet_primitives(t *testing.T) {
	p := &person{Name: "Alice", Age: 30, Score: 9.5, Active: true}

	if v, err := saferefl.Get[string](p, "Name"); err != nil || v != "Alice" {
		t.Errorf("Get Name = %q, err=%v", v, err)
	}
	if v, err := saferefl.Get[int](p, "Age"); err != nil || v != 30 {
		t.Errorf("Get Age = %d, err=%v", v, err)
	}
	if v, err := saferefl.Get[float64](p, "Score"); err != nil || v != 9.5 {
		t.Errorf("Get Score = %v, err=%v", v, err)
	}
	if v, err := saferefl.Get[bool](p, "Active"); err != nil || !v {
		t.Errorf("Get Active = %v, err=%v", v, err)
	}
}

func TestGet_promoted_embedded(t *testing.T) {
	e := &employee{person: person{Name: "Bob", Age: 25}}
	if v, err := saferefl.Get[string](e, "Name"); err != nil || v != "Bob" {
		t.Errorf("Get promoted Name = %q, err=%v", v, err)
	}
}

func TestGet_nested_dotpath_value(t *testing.T) {
	e := &employee{Office: address{City: "Berlin", Country: "DE"}}
	if v, err := saferefl.Get[string](e, "Office.City"); err != nil || v != "Berlin" {
		t.Errorf("Get Office.City = %q, err=%v", v, err)
	}
}

func TestGet_nested_dotpath_pointer(t *testing.T) {
	e := &employee{Contact: &address{City: "NYC", Country: "US"}}
	if v, err := saferefl.Get[string](e, "Contact.City"); err != nil || v != "NYC" {
		t.Errorf("Get Contact.City = %q, err=%v", v, err)
	}
}

func TestGet_type_mismatch(t *testing.T) {
	p := &person{Name: "Alice"}
	_, err := saferefl.Get[int](p, "Name")
	if err == nil {
		t.Fatal("expected TypeMismatchError, got nil")
	}
	var tme *saferefl.TypeMismatchError
	if !errors.As(err, &tme) {
		t.Errorf("want TypeMismatchError, got %T: %v", err, err)
	}
}

func TestGet_field_not_found(t *testing.T) {
	p := &person{}
	_, err := saferefl.Get[string](p, "NonExistent")
	if err == nil {
		t.Fatal("expected FieldNotFoundError, got nil")
	}
	var fnf *saferefl.FieldNotFoundError
	if !errors.As(err, &fnf) {
		t.Errorf("want FieldNotFoundError, got %T: %v", err, err)
	}
}

func TestGet_nil_obj(t *testing.T) {
	_, err := saferefl.Get[string](nil, "Name")
	if err == nil {
		t.Error("expected error for nil obj")
	}
}

func TestGet_non_ptr_obj(t *testing.T) {
	_, err := saferefl.Get[string](person{}, "Name")
	if err == nil {
		t.Error("expected error for non-pointer obj")
	}
}

func TestGet_ptrToNonStruct(t *testing.T) {
	_, err := saferefl.Get[int](new(int), "Field")
	if err == nil {
		t.Error("expected error for pointer to non-struct")
	}
}

// --- MustGet tests ---

func TestMustGet_panic_on_error(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic from MustGet, got none")
		}
	}()
	p := &person{}
	saferefl.MustGet[string](p, "NonExistent")
}

func TestMustGet_returns_value(t *testing.T) {
	p := &person{Name: "Carol"}
	if v := saferefl.MustGet[string](p, "Name"); v != "Carol" {
		t.Errorf("MustGet Name = %q, want Carol", v)
	}
}

// --- Set tests ---

func TestSet_basic(t *testing.T) {
	p := &person{}

	if err := saferefl.Set[string](p, "Name", "Dave"); err != nil {
		t.Fatalf("Set Name: %v", err)
	}
	if p.Name != "Dave" {
		t.Errorf("after Set, Name = %q, want Dave", p.Name)
	}

	if err := saferefl.Set[int](p, "Age", 42); err != nil {
		t.Fatalf("Set Age: %v", err)
	}
	if p.Age != 42 {
		t.Errorf("after Set, Age = %d, want 42", p.Age)
	}
}

func TestSet_nested_dotpath(t *testing.T) {
	e := &employee{Office: address{}}
	if err := saferefl.Set[string](e, "Office.City", "Paris"); err != nil {
		t.Fatalf("Set Office.City: %v", err)
	}
	if e.Office.City != "Paris" {
		t.Errorf("after Set, Office.City = %q, want Paris", e.Office.City)
	}
}

func TestSet_readonly(t *testing.T) {
	p := &person{}
	err := saferefl.Set[string](p, "private", "x")
	if err == nil {
		t.Fatal("expected ReadOnlyError, got nil")
	}
	var roe *saferefl.ReadOnlyError
	if !errors.As(err, &roe) {
		t.Errorf("want ReadOnlyError, got %T: %v", err, err)
	}
}

func TestSet_nil_ptr_path(t *testing.T) {
	e := &employee{Contact: nil}
	err := saferefl.Set[string](e, "Contact.City", "Tokyo")
	if err == nil {
		t.Fatal("expected error traversing nil pointer, got nil")
	}
}

func TestSet_type_mismatch(t *testing.T) {
	p := &person{}
	err := saferefl.Set[int](p, "Name", 42)
	if err == nil {
		t.Fatal("expected TypeMismatchError, got nil")
	}
	var tme *saferefl.TypeMismatchError
	if !errors.As(err, &tme) {
		t.Errorf("want TypeMismatchError, got %T: %v", err, err)
	}
}

// --- MustSet tests ---

func TestMustSet_panic_on_error(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic from MustSet, got none")
		}
	}()
	p := &person{}
	saferefl.MustSet[string](p, "NonExistent", "x")
}

// --- Field utility tests ---

func TestFieldByName(t *testing.T) {
	sf, ok := saferefl.FieldByName[person]("Name")
	if !ok {
		t.Fatal("FieldByName Name: not found")
	}
	if sf.Name != "Name" {
		t.Errorf("got field %q, want Name", sf.Name)
	}
}

func TestFieldByName_nonStruct(t *testing.T) {
	_, ok := saferefl.FieldByName[int]("anything")
	if ok {
		t.Error("expected false for non-struct type T")
	}
}

func TestFields(t *testing.T) {
	fields, err := saferefl.Fields(&person{})
	if err != nil {
		t.Fatalf("Fields: %v", err)
	}
	if len(fields) == 0 {
		t.Error("Fields returned empty slice")
	}
	if fields[0].Name != "Name" {
		t.Errorf("first field = %q, want Name", fields[0].Name)
	}
}

func TestFields_structValue(t *testing.T) {
	fields, err := saferefl.Fields(person{})
	if err != nil {
		t.Fatalf("Fields on struct value: %v", err)
	}
	if len(fields) == 0 {
		t.Error("Fields returned empty slice for struct value")
	}
}

func TestFieldsOf(t *testing.T) {
	fields, err := saferefl.FieldsOf[person]()
	if err != nil {
		t.Fatalf("FieldsOf: %v", err)
	}
	if len(fields) == 0 {
		t.Error("FieldsOf returned empty slice")
	}
}

func TestFieldsOf_nonStruct(t *testing.T) {
	_, err := saferefl.FieldsOf[int]()
	if err == nil {
		t.Error("expected error for non-struct type")
	}
}

// --- Error message tests ---

func TestErrorMessages(t *testing.T) {
	e1 := &saferefl.FieldNotFoundError{Type: "pkg.Foo", FieldPath: "Bar"}
	if e1.Error() == "" {
		t.Error("FieldNotFoundError.Error() returned empty string")
	}
	e2 := &saferefl.TypeMismatchError{FieldPath: "Age", FieldType: "int", WantType: "string"}
	if e2.Error() == "" {
		t.Error("TypeMismatchError.Error() returned empty string")
	}
	e3 := &saferefl.ReadOnlyError{FieldPath: "secret"}
	if e3.Error() == "" {
		t.Error("ReadOnlyError.Error() returned empty string")
	}
}

// --- Additional edge-case tests ---

func TestFields_nilObj(t *testing.T) {
	_, err := saferefl.Fields(nil)
	if err == nil {
		t.Error("expected error for nil obj")
	}
}

func TestFields_nonStruct(t *testing.T) {
	_, err := saferefl.Fields(42)
	if err == nil {
		t.Error("expected error for non-struct obj")
	}
}

func TestSet_nonPtrObj(t *testing.T) {
	err := saferefl.Set[string](person{}, "Name", "x")
	if err == nil {
		t.Error("expected error for non-pointer obj")
	}
}

func TestSet_nilPtrObj(t *testing.T) {
	var p *person
	err := saferefl.Set[string](p, "Name", "x")
	if err == nil {
		t.Error("expected error for nil pointer obj")
	}
}

func TestGet_emptyPath(t *testing.T) {
	p := &person{}
	_, err := saferefl.Get[string](p, "")
	if err == nil {
		t.Error("expected error for empty field path")
	}
}

// --- Accessor tests ---

func TestMakeAccessor_basic(t *testing.T) {
	p := &person{Name: "Alice", Age: 30}

	nameAcc, err := saferefl.MakeAccessor[string](p, "Name")
	if err != nil {
		t.Fatalf("MakeAccessor Name: %v", err)
	}
	if v := nameAcc.Get(saferefl.UnsafePtrOf(p)); v != "Alice" {
		t.Errorf("Get Name = %q, want Alice", v)
	}
	nameAcc.Set(saferefl.UnsafePtrOf(p), "Bob")
	if p.Name != "Bob" {
		t.Errorf("after Set, Name = %q, want Bob", p.Name)
	}

	ageAcc, err := saferefl.MakeAccessor[int](p, "Age")
	if err != nil {
		t.Fatalf("MakeAccessor Age: %v", err)
	}
	if v := ageAcc.Get(saferefl.UnsafePtrOf(p)); v != 30 {
		t.Errorf("Get Age = %d, want 30", v)
	}
}

func TestMakeAccessor_getFrom_setOn(t *testing.T) {
	p := &person{Name: "Carol"}
	acc, _ := saferefl.MakeAccessor[string](p, "Name")

	v, err := acc.GetFrom(p)
	if err != nil || v != "Carol" {
		t.Errorf("GetFrom = %q, err=%v", v, err)
	}
	if err := acc.SetOn(p, "Dave"); err != nil {
		t.Fatalf("SetOn: %v", err)
	}
	if p.Name != "Dave" {
		t.Errorf("after SetOn, Name = %q, want Dave", p.Name)
	}
}

func TestMakeAccessor_dotpath_value(t *testing.T) {
	e := &employee{Office: address{City: "Berlin"}}
	acc, err := saferefl.MakeAccessor[string](e, "Office.City")
	if err != nil {
		t.Fatalf("MakeAccessor Office.City: %v", err)
	}
	if v := acc.Get(saferefl.UnsafePtrOf(e)); v != "Berlin" {
		t.Errorf("Get Office.City = %q, want Berlin", v)
	}
	acc.Set(saferefl.UnsafePtrOf(e), "Paris")
	if e.Office.City != "Paris" {
		t.Errorf("after Set, Office.City = %q, want Paris", e.Office.City)
	}
}

func TestMakeAccessor_dotpath_pointer(t *testing.T) {
	e := &employee{Contact: &address{City: "NYC"}}
	acc, err := saferefl.MakeAccessor[string](e, "Contact.City")
	if err != nil {
		t.Fatalf("MakeAccessor Contact.City: %v", err)
	}
	if v := acc.Get(saferefl.UnsafePtrOf(e)); v != "NYC" {
		t.Errorf("Get Contact.City = %q, want NYC", v)
	}
	acc.Set(saferefl.UnsafePtrOf(e), "London")
	if e.Contact.City != "London" {
		t.Errorf("after Set via chain, Contact.City = %q, want London", e.Contact.City)
	}
}

func TestMakeAccessor_typeMismatch(t *testing.T) {
	p := &person{}
	_, err := saferefl.MakeAccessor[int](p, "Name")
	if err == nil {
		t.Fatal("expected TypeMismatchError, got nil")
	}
	var tme *saferefl.TypeMismatchError
	if !errors.As(err, &tme) {
		t.Errorf("want TypeMismatchError, got %T", err)
	}
}

func TestMakeAccessor_fieldNotFound(t *testing.T) {
	p := &person{}
	_, err := saferefl.MakeAccessor[string](p, "NoSuchField")
	if err == nil {
		t.Fatal("expected FieldNotFoundError, got nil")
	}
	var fnf *saferefl.FieldNotFoundError
	if !errors.As(err, &fnf) {
		t.Errorf("want FieldNotFoundError, got %T", err)
	}
}

func TestMakeAccessor_nilObj(t *testing.T) {
	if _, err := saferefl.MakeAccessor[string](nil, "Name"); err == nil {
		t.Fatal("expected error for nil obj, got nil")
	}
}

func TestMakeAccessor_nonPtrObj(t *testing.T) {
	if _, err := saferefl.MakeAccessor[string](person{}, "Name"); err == nil {
		t.Fatal("expected error for non-pointer obj, got nil")
	}
}

func TestMakeAccessor_ptrToNonStruct(t *testing.T) {
	if _, err := saferefl.MakeAccessor[int](new(int), "Field"); err == nil {
		t.Fatal("expected error for pointer to non-struct, got nil")
	}
}

func TestMakeAccessor_intermediateNotStruct(t *testing.T) {
	// "Name.Sub" where Name is a string — intermediate segment is not a struct.
	if _, err := saferefl.MakeAccessor[string](&person{}, "Name.Sub"); err == nil {
		t.Fatal("expected error for non-struct intermediate segment, got nil")
	}
}

func TestMakeAccessor_getFrom_setOn_errors(t *testing.T) {
	acc, _ := saferefl.MakeAccessor[string](&person{}, "Name")

	if _, err := acc.GetFrom(nil); err == nil {
		t.Error("GetFrom(nil): expected error, got nil")
	}
	if _, err := acc.GetFrom(42); err == nil {
		t.Error("GetFrom(non-ptr): expected error, got nil")
	}
	if err := acc.SetOn(nil, "x"); err == nil {
		t.Error("SetOn(nil): expected error, got nil")
	}
}

func TestSet_interfaceField(t *testing.T) {
	type withAny struct{ V any }
	s := &withAny{}
	if err := saferefl.Set[string](s, "V", "hello"); err != nil {
		t.Fatalf("Set interface field: %v", err)
	}
	if s.V != "hello" {
		t.Errorf("V = %v, want hello", s.V)
	}
}

// --- Fuzz ---

func FuzzGet(f *testing.F) {
	type seed struct {
		Name   string
		Age    int
		Active bool
	}
	s := &seed{Name: "Alice", Age: 30, Active: true}

	f.Add("Name")
	f.Add("Age")
	f.Add("")
	f.Add(".")
	f.Add("NonExistent")
	f.Add("Name.Sub")
	f.Add("Age.Sub.Deep")

	f.Fuzz(func(t *testing.T, fieldPath string) {
		// Must never panic — only return errors.
		_, _ = saferefl.Get[any](s, fieldPath)
	})
}
