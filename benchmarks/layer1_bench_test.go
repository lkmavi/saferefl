package benchmarks

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/lkmavi/saferefl"
)

type benchUser struct {
	ID   int
	Name string
}

var (
	sinkInt    int
	sinkString string
)

// BenchmarkGet_int_L1 — Layer 1: Get[int] with full path resolution per call.
func BenchmarkGet_int_L1(b *testing.B) {
	u := &benchUser{ID: 42}
	b.ResetTimer()
	for i := range b.N {
		v, _ := saferefl.Get[int](u, "ID")
		sinkInt = v + i
	}
}

func BenchmarkGet_int_Reflect(b *testing.B) {
	u := &benchUser{ID: 42}
	rv := reflect.ValueOf(u).Elem()
	b.ResetTimer()
	for i := range b.N {
		sinkInt = int(rv.Field(0).Int()) + i
	}
}

// BenchmarkGet_int_L2 — Layer 2: pre-computed offset + reflect.NewAt (warm TypeInfo cache path).
func BenchmarkGet_int_L2(b *testing.B) {
	u := &benchUser{ID: 42}
	rt := reflect.TypeOf(benchUser{})
	f, _ := rt.FieldByName("ID")
	offset := f.Offset
	rtype := f.Type
	ptr := unsafe.Pointer(u)
	b.ResetTimer()
	for i := range b.N {
		sinkInt = int(reflect.NewAt(rtype, unsafe.Pointer(uintptr(ptr)+offset)).Elem().Int()) + i
	}
}

// BenchmarkSet_string_L1 — Layer 1: Set[string] with full path resolution per call.
func BenchmarkSet_string_L1(b *testing.B) {
	u := &benchUser{}
	b.ResetTimer()
	for range b.N {
		_ = saferefl.Set[string](u, "Name", "Alice")
	}
	sinkString = u.Name
}

func BenchmarkSet_string_Reflect(b *testing.B) {
	u := &benchUser{}
	rv := reflect.ValueOf(u).Elem()
	b.ResetTimer()
	for range b.N {
		rv.Field(1).SetString("Alice")
	}
	sinkString = u.Name
}

// BenchmarkSet_string_L2 — Layer 2: pre-computed offset + reflect.NewAt (warm TypeInfo cache path).
func BenchmarkSet_string_L2(b *testing.B) {
	u := &benchUser{}
	rt := reflect.TypeOf(benchUser{})
	f, _ := rt.FieldByName("Name")
	offset := f.Offset
	rtype := f.Type
	ptr := unsafe.Pointer(u)
	b.ResetTimer()
	for range b.N {
		reflect.NewAt(rtype, unsafe.Pointer(uintptr(ptr)+offset)).Elem().SetString("Alice")
	}
	sinkString = u.Name
}

// --- Accessor benchmarks: pre-bound path, per-call cost only ---

var (
	benchIDAccessor   = mustMakeAccessor[int](&benchUser{}, "ID")
	benchNameAccessor = mustMakeAccessor[string](&benchUser{}, "Name")
)

func mustMakeAccessor[T any](obj any, path string) saferefl.Accessor[T] {
	acc, err := saferefl.MakeAccessor[T](obj, path)
	if err != nil {
		panic(err)
	}
	return acc
}

// BenchmarkGet_int_L3 — Layer 3: Accessor with pre-extracted unsafe.Pointer (maximum throughput).
func BenchmarkGet_int_L3(b *testing.B) {
	u := &benchUser{ID: 42}
	ptr := saferefl.UnsafePtrOf(u)
	b.ResetTimer()
	for i := range b.N {
		sinkInt = benchIDAccessor.Get(ptr) + i
	}
}

// BenchmarkGet_int_L3From — Layer 3: Accessor.GetFrom (interface convenience path).
func BenchmarkGet_int_L3From(b *testing.B) {
	u := &benchUser{ID: 42}
	b.ResetTimer()
	for i := range b.N {
		v, _ := benchIDAccessor.GetFrom(u)
		sinkInt = v + i
	}
}

// BenchmarkSet_string_L3 — Layer 3: Accessor.Set with pre-extracted pointer.
func BenchmarkSet_string_L3(b *testing.B) {
	u := &benchUser{}
	ptr := saferefl.UnsafePtrOf(u)
	b.ResetTimer()
	for range b.N {
		benchNameAccessor.Set(ptr, "Alice")
	}
	sinkString = u.Name
}

// BenchmarkSet_string_L3On — Layer 3: Accessor.SetOn (interface convenience path).
func BenchmarkSet_string_L3On(b *testing.B) {
	u := &benchUser{}
	b.ResetTimer()
	for range b.N {
		_ = benchNameAccessor.SetOn(u, "Alice")
	}
	sinkString = u.Name
}

// BenchmarkSet_string_Direct — native pointer write for reference (theoretical minimum).
func BenchmarkSet_string_Direct(b *testing.B) {
	u := &benchUser{}
	ptr := unsafe.Pointer(u)
	const nameOffset = unsafe.Offsetof(benchUser{}.Name)
	b.ResetTimer()
	for range b.N {
		*(*string)(unsafe.Pointer(uintptr(ptr) + nameOffset)) = "Alice"
	}
	sinkString = u.Name
}

// BenchmarkGet_int_Direct — native pointer read for reference (theoretical minimum).
func BenchmarkGet_int_Direct(b *testing.B) {
	u := &benchUser{ID: 42}
	ptr := unsafe.Pointer(u)
	const idOffset = unsafe.Offsetof(benchUser{}.ID)
	b.ResetTimer()
	for i := range b.N {
		sinkInt = *(*int)(unsafe.Pointer(uintptr(ptr) + idOffset)) + i
	}
}
