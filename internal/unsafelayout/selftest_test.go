//go:build !reflectx_strict

package unsafelayout

import (
	"reflect"
	"testing"
	"unsafe"
)

func TestSelfTest_passes(t *testing.T) {
	if !runSelfTest() {
		t.Fatal("runSelfTest() = false on current Go version/arch")
	}
}

func TestSelfTest_struct_offset_match(t *testing.T) {
	if !selfTestStruct() {
		t.Fatal("struct offset self-test failed")
	}
}

func TestSelfTest_slice_elem(t *testing.T) {
	if !selfTestSlice() {
		t.Fatal("slice element self-test failed")
	}
}

func TestSelfTest_map_len(t *testing.T) {
	if !selfTestMap() {
		t.Fatal("map length self-test failed")
	}
}

func TestAccelAvailable(t *testing.T) {
	// After package init, AccelAvailable must match runSelfTest().
	if AccelAvailable() != runSelfTest() {
		t.Error("AccelAvailable() disagrees with runSelfTest()")
	}
}

func TestEnableAccel(t *testing.T) {
	err := EnableAccel()
	if AccelAvailable() && err != nil {
		t.Errorf("EnableAccel() = %v, want nil when AccelAvailable=true", err)
	}
	if !AccelAvailable() && err == nil {
		t.Error("EnableAccel() = nil, want error when AccelAvailable=false")
	}
}

func TestEnableAccel_whenFailed(t *testing.T) {
	prev := accelOK
	accelOK = false
	defer func() { accelOK = prev }()

	if err := EnableAccel(); err == nil {
		t.Error("EnableAccel() = nil, want error when accelOK=false")
	}
}

func TestUnsafeFieldPtr_matches_reflect(t *testing.T) {
	type sample struct {
		X int
		Y string
		Z float64
	}
	rt := reflect.TypeOf(sample{})
	s := sample{X: 7, Y: "hi", Z: 1.5}
	ptr := unsafe.Pointer(&s)

	for _, name := range []string{"X", "Y", "Z"} {
		sf, _ := rt.FieldByName(name)
		got := UnsafeFieldPtr(ptr, sf.Offset)
		want := unsafe.Pointer(reflect.ValueOf(&s).Elem().FieldByName(name).UnsafeAddr())
		if got != want {
			t.Errorf("UnsafeFieldPtr %q: got %p, want %p", name, got, want)
		}
	}
}

func TestUnsafeSliceElemPtr_roundtrip(t *testing.T) {
	s := []int32{100, 200, 300}
	sliceData := *(*unsafe.Pointer)(unsafe.Pointer(&s))
	size := unsafe.Sizeof(int32(0))

	for i, want := range s {
		ptr := UnsafeSliceElemPtr(sliceData, i, size)
		got := *(*int32)(ptr)
		if got != want {
			t.Errorf("index %d: got %d, want %d", i, got, want)
		}
	}
}

func TestMapLen_matches_builtin(t *testing.T) {
	maps := []map[string]int{
		{},
		{"a": 1},
		{"a": 1, "b": 2, "c": 3},
	}
	for _, m := range maps {
		mapPtr := unsafe.Pointer(reflect.ValueOf(m).Pointer())
		if len(m) == 0 {
			// reflect.Value.Pointer() on nil/empty map may return 0; skip.
			continue
		}
		got := MapLen(mapPtr)
		if got != len(m) {
			t.Errorf("MapLen = %d, want %d", got, len(m))
		}
	}
}

func FuzzUnsafeFieldPtr(f *testing.F) {
	type sample struct {
		A int64
		B string
		C float64
		D bool
		E uint32
	}
	rt := reflect.TypeOf(sample{})
	for i := 0; i < rt.NumField(); i++ {
		f.Add(i)
	}
	f.Fuzz(func(t *testing.T, fieldIdx int) {
		if fieldIdx < 0 || fieldIdx >= rt.NumField() {
			return
		}
		sf := rt.Field(fieldIdx)
		s := sample{A: 0xDEAD, B: "test", C: 3.14, D: true, E: 42}
		ptr := unsafe.Pointer(&s)
		got := UnsafeFieldPtr(ptr, sf.Offset)
		want := unsafe.Pointer(reflect.ValueOf(&s).Elem().Field(fieldIdx).UnsafeAddr())
		if got != want {
			t.Errorf("field %q: got %p, want %p", sf.Name, got, want)
		}
	})
}
