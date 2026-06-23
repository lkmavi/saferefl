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

// BenchmarkGet_int_SafeRefl — saferefl.Get[int]: named field access per call.
func BenchmarkGet_int_SafeRefl(b *testing.B) {
	u := &benchUser{ID: 42}
	b.ResetTimer()
	for i := range b.N {
		v, _ := saferefl.Get[int](u, "ID")
		sinkInt = v + i
	}
}

// BenchmarkGet_int_Reflect — stdlib reflect: FieldByName per call (fair per-call cost).
func BenchmarkGet_int_Reflect(b *testing.B) {
	u := &benchUser{ID: 42}
	b.ResetTimer()
	for i := range b.N {
		sinkInt = int(reflect.ValueOf(u).Elem().FieldByName("ID").Int()) + i
	}
}

// BenchmarkGet_int_ReflectFast — stdlib reflect: pre-cached Value + Field(0) (best possible reflect).
func BenchmarkGet_int_ReflectFast(b *testing.B) {
	u := &benchUser{ID: 42}
	rv := reflect.ValueOf(u).Elem()
	b.ResetTimer()
	for i := range b.N {
		sinkInt = int(rv.Field(0).Int()) + i
	}
}

// BenchmarkGet_int_Offset — pre-computed offset + reflect.NewAt (internal cache mechanism).
func BenchmarkGet_int_Offset(b *testing.B) {
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

// BenchmarkSet_string_SafeRefl — saferefl.Set[string]: named field write per call.
func BenchmarkSet_string_SafeRefl(b *testing.B) {
	u := &benchUser{}
	b.ResetTimer()
	for range b.N {
		_ = saferefl.Set[string](u, "Name", "Alice")
	}
	sinkString = u.Name
}

// BenchmarkSet_string_Reflect — stdlib reflect: FieldByName per call (fair per-call cost).
func BenchmarkSet_string_Reflect(b *testing.B) {
	u := &benchUser{}
	b.ResetTimer()
	for range b.N {
		reflect.ValueOf(u).Elem().FieldByName("Name").SetString("Alice")
	}
	sinkString = u.Name
}

// BenchmarkSet_string_ReflectFast — stdlib reflect: pre-cached Value + Field(1) (best possible reflect).
func BenchmarkSet_string_ReflectFast(b *testing.B) {
	u := &benchUser{}
	rv := reflect.ValueOf(u).Elem()
	b.ResetTimer()
	for range b.N {
		rv.Field(1).SetString("Alice")
	}
	sinkString = u.Name
}

// BenchmarkSet_string_Offset — pre-computed offset + reflect.NewAt (internal cache mechanism).
func BenchmarkSet_string_Offset(b *testing.B) {
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

// BenchmarkGet_int_Accessor — Accessor.Get with pre-extracted pointer (maximum throughput).
func BenchmarkGet_int_Accessor(b *testing.B) {
	u := &benchUser{ID: 42}
	ptr := saferefl.UnsafePtrOf(u)
	b.ResetTimer()
	for i := range b.N {
		sinkInt = benchIDAccessor.Get(ptr) + i
	}
}

// BenchmarkGet_int_AccFrom — Accessor.GetFrom: interface convenience path.
func BenchmarkGet_int_AccFrom(b *testing.B) {
	u := &benchUser{ID: 42}
	b.ResetTimer()
	for i := range b.N {
		v, _ := benchIDAccessor.GetFrom(u)
		sinkInt = v + i
	}
}

// BenchmarkSet_string_Accessor — Accessor.Set with pre-extracted pointer (maximum throughput).
func BenchmarkSet_string_Accessor(b *testing.B) {
	u := &benchUser{}
	ptr := saferefl.UnsafePtrOf(u)
	b.ResetTimer()
	for range b.N {
		benchNameAccessor.Set(ptr, "Alice")
	}
	sinkString = u.Name
}

// BenchmarkSet_string_AccOn — Accessor.SetOn: interface convenience path.
func BenchmarkSet_string_AccOn(b *testing.B) {
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
