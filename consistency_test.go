//go:build !reflectx_strict

package saferefl_test

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/lkmavi/saferefl"
	"github.com/lkmavi/saferefl/internal/unsafelayout"
)

// TestConsistency verifies that the Get[any] pipeline (Layer 1 + Layer 2) and
// UnsafeFieldPtr (Layer 3) both produce values that match the reflect baseline
// for every exported field of a representative set of struct types.
func TestConsistency(t *testing.T) {
	if !unsafelayout.AccelAvailable() {
		t.Skip("unsafe accelerator not available on this Go version/arch")
	}

	t.Run("int_variants", func(t *testing.T) {
		type S struct {
			A int
			B int32
			C int64
			D uint
			E uint64
		}
		checkConsistency(t, &S{A: 1, B: 2, C: 3, D: 4, E: 5})
	})

	t.Run("string_bool_float", func(t *testing.T) {
		type S struct {
			Name   string
			Active bool
			Score  float64
			Rate   float32
		}
		checkConsistency(t, &S{Name: "hello", Active: true, Score: 9.99, Rate: 1.5})
	})

	t.Run("mixed_types", func(t *testing.T) {
		type S struct {
			X int64
			Y string
			Z float32
			W bool
			V uint32
		}
		checkConsistency(t, &S{X: 42, Y: "world", Z: 1.5, W: true, V: 99})
	})

	t.Run("pointer_field", func(t *testing.T) {
		n := 7
		type S struct {
			N *int
			S string
		}
		checkConsistency(t, &S{N: &n, S: "ptr"})
	})
}

// checkConsistency verifies field-by-field that Get[any] and UnsafeFieldPtr
// both agree with the reflect baseline for every exported field of obj.
func checkConsistency(t *testing.T, obj any) {
	t.Helper()
	rv := reflect.ValueOf(obj).Elem()
	rt := rv.Type()
	objPtr := unsafe.Pointer(reflect.ValueOf(obj).Pointer())

	for i := 0; i < rt.NumField(); i++ {
		sf := rt.Field(i)
		if !sf.IsExported() {
			continue
		}
		want := rv.Field(i).Interface()

		got, err := saferefl.Get[any](obj, sf.Name)
		if err != nil {
			t.Errorf("field %q: Get[any] error: %v", sf.Name, err)
			continue
		}
		if got != want {
			t.Errorf("field %q: Get[any]=%v, reflect=%v", sf.Name, got, want)
		}

		fptr := unsafelayout.UnsafeFieldPtr(objPtr, sf.Offset)
		gotUnsafe := reflect.NewAt(sf.Type, fptr).Elem().Interface()
		if gotUnsafe != want {
			t.Errorf("field %q: UnsafeFieldPtr=%v, reflect=%v", sf.Name, gotUnsafe, want)
		}
	}
}
