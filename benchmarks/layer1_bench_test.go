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

// BenchmarkGet_int compares Get[int] against plain reflect field access.
func BenchmarkGet_int_Saferefl(b *testing.B) {
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

// BenchmarkSet_string compares Set[string] against plain reflect field write.
func BenchmarkSet_string_Saferefl(b *testing.B) {
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

// BenchmarkGet_int_Accessor — Accessor with pre-extracted unsafe.Pointer (maximum throughput).
func BenchmarkGet_int_Accessor(b *testing.B) {
	u := &benchUser{ID: 42}
	ptr := saferefl.UnsafePtrOf(u)
	b.ResetTimer()
	for i := range b.N {
		sinkInt = benchIDAccessor.Get(ptr) + i
	}
}

// BenchmarkGet_int_AccessorFrom — Accessor.GetFrom (interface convenience path).
func BenchmarkGet_int_AccessorFrom(b *testing.B) {
	u := &benchUser{ID: 42}
	b.ResetTimer()
	for i := range b.N {
		v, _ := benchIDAccessor.GetFrom(u)
		sinkInt = v + i
	}
}

// BenchmarkSet_string_Accessor — Accessor.Set with pre-extracted pointer.
func BenchmarkSet_string_Accessor(b *testing.B) {
	u := &benchUser{}
	ptr := saferefl.UnsafePtrOf(u)
	b.ResetTimer()
	for range b.N {
		benchNameAccessor.Set(ptr, "Alice")
	}
	sinkString = u.Name
}

// BenchmarkSet_string_AccessorOn — Accessor.SetOn (interface convenience path).
func BenchmarkSet_string_AccessorOn(b *testing.B) {
	u := &benchUser{}
	b.ResetTimer()
	for range b.N {
		_ = benchNameAccessor.SetOn(u, "Alice")
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
		sinkInt = *(*int)(unsafe.Pointer(uintptr(ptr)+idOffset)) + i
	}
}
