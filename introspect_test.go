package saferefl_test

import (
	"io"
	"reflect"
	"testing"

	"github.com/lkmavi/saferefl"
)

// --- KindOf ---

func TestKindOf_basicTypes(t *testing.T) {
	cases := []struct {
		val  any
		want reflect.Kind
	}{
		{42, reflect.Int},
		{"hello", reflect.String},
		{3.14, reflect.Float64},
		{true, reflect.Bool},
		{(*int)(nil), reflect.Pointer},
		{[]int{}, reflect.Slice},
		{map[string]int{}, reflect.Map},
		{make(chan int), reflect.Chan},
	}
	for _, tc := range cases {
		if got := saferefl.KindOf(tc.val); got != tc.want {
			t.Errorf("KindOf(%T) = %v, want %v", tc.val, got, tc.want)
		}
	}
}

func TestKindOf_nil(t *testing.T) {
	if got := saferefl.KindOf(nil); got != reflect.Invalid {
		t.Errorf("KindOf(nil) = %v, want Invalid", got)
	}
}

func TestKindOf_struct(t *testing.T) {
	if got := saferefl.KindOf(person{}); got != reflect.Struct {
		t.Errorf("KindOf(person{}) = %v, want Struct", got)
	}
}

func TestKindOf_ptrToStruct(t *testing.T) {
	if got := saferefl.KindOf(&person{}); got != reflect.Pointer {
		t.Errorf("KindOf(&person{}) = %v, want Pointer", got)
	}
}

func TestKindOf_matches_reflect(t *testing.T) {
	vals := []any{42, "x", 3.14, true, (*int)(nil), []int{1}, map[string]int{"a": 1}, person{}}
	for _, v := range vals {
		want := reflect.TypeOf(v).Kind()
		if got := saferefl.KindOf(v); got != want {
			t.Errorf("KindOf(%T): got %v, reflect says %v", v, got, want)
		}
	}
}

// --- IsNil ---

func TestIsNil_nilInterface(t *testing.T) {
	if !saferefl.IsNil(nil) {
		t.Error("IsNil(nil) = false, want true")
	}
}

func TestIsNil_nilPointer(t *testing.T) {
	var p *int
	if !saferefl.IsNil(p) {
		t.Error("IsNil((*int)(nil)) = false, want true")
	}
}

func TestIsNil_nonNilPointer(t *testing.T) {
	x := 1
	if saferefl.IsNil(&x) {
		t.Error("IsNil(&x) = true, want false")
	}
}

func TestIsNil_nilMap(t *testing.T) {
	var m map[string]int
	if !saferefl.IsNil(m) {
		t.Error("IsNil(nil map) = false, want true")
	}
}

func TestIsNil_nonNilMap(t *testing.T) {
	m := map[string]int{}
	if saferefl.IsNil(m) {
		t.Error("IsNil(non-nil map) = true, want false")
	}
}

func TestIsNil_nilChan(t *testing.T) {
	var ch chan int
	if !saferefl.IsNil(ch) {
		t.Error("IsNil(nil chan) = false, want true")
	}
}

func TestIsNil_nilSlice(t *testing.T) {
	var s []int
	if !saferefl.IsNil(s) {
		t.Error("IsNil(nil slice) = false, want true")
	}
}

func TestIsNil_emptySlice(t *testing.T) {
	s := []int{}
	if saferefl.IsNil(s) {
		t.Error("IsNil(empty slice) = true, want false")
	}
}

func TestIsNil_nonNilableTypes(t *testing.T) {
	vals := []any{42, "hello", 3.14, true, person{}}
	for _, v := range vals {
		if saferefl.IsNil(v) {
			t.Errorf("IsNil(%T{}) = true, want false", v)
		}
	}
}

// --- TypeOf ---

func TestTypeOf_matchesReflect(t *testing.T) {
	cases := []struct {
		got  reflect.Type
		want reflect.Type
	}{
		{saferefl.TypeOf[int](), reflect.TypeOf(0)},
		{saferefl.TypeOf[string](), reflect.TypeOf("")},
		{saferefl.TypeOf[float64](), reflect.TypeOf(0.0)},
		{saferefl.TypeOf[bool](), reflect.TypeOf(false)},
		{saferefl.TypeOf[person](), reflect.TypeOf(person{})},
		{saferefl.TypeOf[*person](), reflect.TypeOf((*person)(nil))},
	}
	for _, tc := range cases {
		if tc.got != tc.want {
			t.Errorf("TypeOf[%v] = %v, want %v", tc.want, tc.got, tc.want)
		}
	}
}

func TestTypeOf_interfaceType(t *testing.T) {
	got := saferefl.TypeOf[error]()
	want := reflect.TypeOf((*error)(nil)).Elem()
	if got != want {
		t.Errorf("TypeOf[error] = %v, want %v", got, want)
	}
	if got.Kind() != reflect.Interface {
		t.Errorf("TypeOf[error].Kind() = %v, want Interface", got.Kind())
	}

	got = saferefl.TypeOf[io.Reader]()
	want = reflect.TypeOf((*io.Reader)(nil)).Elem()
	if got != want {
		t.Errorf("TypeOf[io.Reader] = %v, want %v", got, want)
	}
}

func TestTypeOf_usableAsMapKey(t *testing.T) {
	m := map[reflect.Type]string{
		saferefl.TypeOf[int]():    "int",
		saferefl.TypeOf[string](): "string",
		saferefl.TypeOf[person](): "person",
	}
	if m[saferefl.TypeOf[int]()] != "int" {
		t.Error("map lookup by TypeOf[int] failed")
	}
	if m[saferefl.TypeOf[string]()] != "string" {
		t.Error("map lookup by TypeOf[string] failed")
	}
	if m[saferefl.TypeOf[person]()] != "person" {
		t.Error("map lookup by TypeOf[person] failed")
	}
}

func TestTypeOf_sliceAndMap(t *testing.T) {
	if got := saferefl.TypeOf[[]int](); got.Kind() != reflect.Slice {
		t.Errorf("TypeOf[[]int].Kind() = %v, want Slice", got.Kind())
	}
	if got := saferefl.TypeOf[map[string]int](); got.Kind() != reflect.Map {
		t.Errorf("TypeOf[map[string]int].Kind() = %v, want Map", got.Kind())
	}
}
