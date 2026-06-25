package saferefl_test

import (
	"errors"
	"testing"

	"github.com/lkmavi/saferefl"
)

type tagged struct {
	Name    string  `json:"name"          db:"user_name"`
	Age     int     `json:"age,omitempty" db:"user_age"`
	Score   float64 `json:"score"`
	private string  `db:"priv"` //nolint:unused
}

// --- GetByTag ---

func TestGetByTag_basic(t *testing.T) {
	u := &tagged{Name: "Alice", Age: 30, Score: 9.5}

	name, err := saferefl.GetByTag[string](u, "json", "name")
	if err != nil || name != "Alice" {
		t.Fatalf("GetByTag json:name = %q, %v", name, err)
	}

	age, err := saferefl.GetByTag[int](u, "json", "age")
	if err != nil || age != 30 {
		t.Fatalf("GetByTag json:age = %d, %v", age, err)
	}
}

func TestGetByTag_alternate_key(t *testing.T) {
	u := &tagged{Name: "Bob"}
	name, err := saferefl.GetByTag[string](u, "db", "user_name")
	if err != nil || name != "Bob" {
		t.Fatalf("GetByTag db:user_name = %q, %v", name, err)
	}
}

func TestGetByTag_omitempty_stripped(t *testing.T) {
	// "age,omitempty" → tagValue is "age", not "age,omitempty"
	u := &tagged{Age: 42}
	age, err := saferefl.GetByTag[int](u, "json", "age")
	if err != nil || age != 42 {
		t.Fatalf("GetByTag json:age (omitempty) = %d, %v", age, err)
	}
}

func TestGetByTag_tag_key_not_found(t *testing.T) {
	u := &tagged{}
	_, err := saferefl.GetByTag[string](u, "xml", "name")
	if err == nil {
		t.Fatal("expected error for missing tag key, got nil")
	}
	if !errors.Is(err, saferefl.ErrFieldNotFound) {
		t.Errorf("expected ErrFieldNotFound, got %T: %v", err, err)
	}
}

func TestGetByTag_tag_value_not_found(t *testing.T) {
	u := &tagged{}
	_, err := saferefl.GetByTag[string](u, "json", "nonexistent")
	if err == nil {
		t.Fatal("expected error for missing tag value, got nil")
	}
	if !errors.Is(err, saferefl.ErrFieldNotFound) {
		t.Errorf("expected ErrFieldNotFound, got %T: %v", err, err)
	}
}

func TestGetByTag_type_mismatch(t *testing.T) {
	u := &tagged{Name: "X"}
	_, err := saferefl.GetByTag[int](u, "json", "name") // name is string
	if err == nil {
		t.Fatal("expected TypeMismatchError, got nil")
	}
	if !errors.Is(err, saferefl.ErrTypeMismatch) {
		t.Errorf("expected ErrTypeMismatch, got %T: %v", err, err)
	}
	var tme *saferefl.TypeMismatchError
	if !errors.As(err, &tme) {
		t.Errorf("expected *TypeMismatchError via errors.As, got %T", err)
	}
}

func TestGetByTag_nil_obj(t *testing.T) {
	_, err := saferefl.GetByTag[string](nil, "json", "name")
	if err == nil {
		t.Fatal("expected error for nil obj")
	}
}

func TestGetByTag_non_ptr_obj(t *testing.T) {
	_, err := saferefl.GetByTag[string](tagged{}, "json", "name")
	if err == nil {
		t.Fatal("expected error for non-pointer obj")
	}
}

// --- SetByTag ---

func TestSetByTag_basic(t *testing.T) {
	u := &tagged{}
	if err := saferefl.SetByTag[string](u, "json", "name", "Carol"); err != nil {
		t.Fatalf("SetByTag json:name: %v", err)
	}
	if u.Name != "Carol" {
		t.Errorf("Name = %q, want Carol", u.Name)
	}
}

func TestSetByTag_alternate_key(t *testing.T) {
	u := &tagged{}
	if err := saferefl.SetByTag[int](u, "db", "user_age", 25); err != nil {
		t.Fatalf("SetByTag db:user_age: %v", err)
	}
	if u.Age != 25 {
		t.Errorf("Age = %d, want 25", u.Age)
	}
}

func TestSetByTag_readonly(t *testing.T) {
	u := &tagged{}
	err := saferefl.SetByTag[string](u, "db", "priv", "x")
	if err == nil {
		t.Fatal("expected ReadOnlyError, got nil")
	}
	if !errors.Is(err, saferefl.ErrReadOnly) {
		t.Errorf("expected ErrReadOnly, got %T: %v", err, err)
	}
}

func TestSetByTag_type_mismatch(t *testing.T) {
	u := &tagged{}
	err := saferefl.SetByTag[int](u, "json", "name", 42) // name is string
	if err == nil {
		t.Fatal("expected TypeMismatchError, got nil")
	}
	if !errors.Is(err, saferefl.ErrTypeMismatch) {
		t.Errorf("expected ErrTypeMismatch, got %T: %v", err, err)
	}
}

func TestSetByTag_tag_not_found(t *testing.T) {
	u := &tagged{}
	err := saferefl.SetByTag[string](u, "json", "no_such_field", "x")
	if err == nil {
		t.Fatal("expected FieldNotFoundError, got nil")
	}
	if !errors.Is(err, saferefl.ErrFieldNotFound) {
		t.Errorf("expected ErrFieldNotFound, got %T: %v", err, err)
	}
}

// --- SetByTag: nil / non-ptr / nil-ptr guard ---

func TestSetByTag_nil_obj(t *testing.T) {
	if err := saferefl.SetByTag[string](nil, "json", "name", "x"); err == nil {
		t.Error("expected error for nil obj")
	}
}

func TestSetByTag_non_ptr_obj(t *testing.T) {
	if err := saferefl.SetByTag[string](tagged{}, "json", "name", "x"); err == nil {
		t.Error("expected error for non-pointer obj")
	}
}

func TestSetByTag_nil_ptr(t *testing.T) {
	if err := saferefl.SetByTag[string]((*tagged)(nil), "json", "name", "x"); err == nil {
		t.Error("expected error for nil pointer")
	}
}

// --- setByTagSlowPath coverage ---

// freshTagged is only used in TestSetByTag_slowPath, so SetByTag is called
// before GetByTag for this type — the ptr cache is cold and setByTagSlowPath runs.
type freshTagged struct {
	Value string `json:"value"`
}

func TestSetByTag_slowPath(t *testing.T) {
	u := &freshTagged{}
	if err := saferefl.SetByTag[string](u, "json", "value", "hello"); err != nil {
		t.Fatalf("SetByTag slow path: %v", err)
	}
	if u.Value != "hello" {
		t.Errorf("Value = %q, want hello", u.Value)
	}
}

// --- errors.Is / errors.As ---

func TestErrors_Is(t *testing.T) {
	u := &tagged{}

	_, err := saferefl.Get[string](u, "NoSuchField")
	if !errors.Is(err, saferefl.ErrFieldNotFound) {
		t.Errorf("Get missing field: expected ErrFieldNotFound, got %v", err)
	}

	_, err = saferefl.Get[int](u, "Name")
	if !errors.Is(err, saferefl.ErrTypeMismatch) {
		t.Errorf("Get type mismatch: expected ErrTypeMismatch, got %v", err)
	}

	err = saferefl.Set[string](u, "private", "x")
	if !errors.Is(err, saferefl.ErrReadOnly) {
		t.Errorf("Set unexported: expected ErrReadOnly, got %v", err)
	}
}
