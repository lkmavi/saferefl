package typeinfo

import (
	"reflect"
	"testing"
	"unsafe"
)

type accessSample struct {
	X int64
	Y string
	Z float64
}

func TestGetFieldFast_matchesGetFieldPtr(t *testing.T) {
	s := accessSample{X: 7, Y: "ok", Z: 3.14}
	ptr := unsafe.Pointer(&s)
	desc := TypeDescriptorOf(reflect.TypeOf(s))

	for _, name := range []string{"X", "Y", "Z"} {
		fm := desc.FieldsByName[name]
		got := GetFieldFast(ptr, fm)
		want := GetFieldPtr(ptr, fm)
		if got != want {
			t.Errorf("field %q: GetFieldFast=%p, GetFieldPtr=%p", name, got, want)
		}
	}
}

func TestGetFieldFast_readValue(t *testing.T) {
	s := accessSample{X: 42}
	ptr := unsafe.Pointer(&s)
	desc := TypeDescriptorOf(reflect.TypeOf(s))
	fm := desc.FieldsByName["X"]

	got := *(*int64)(GetFieldFast(ptr, fm))
	if got != 42 {
		t.Errorf("GetFieldFast read X = %d, want 42", got)
	}
}
