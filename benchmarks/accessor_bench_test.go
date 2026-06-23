package benchmarks

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	"github.com/lkmavi/saferefl"
)

// sinkIntPtr prevents the compiler from optimising away slice-element reads.
var sinkIntPtr *int

// ---- FieldRead: variant comparison for a single int field ----
//
// Reflect     = stdlib reflect: FieldByName per call (standard usage, fair baseline)
// ReflectFast = stdlib reflect: pre-cached Value + Field(0) (best possible reflect)
// SafeRefl    = saferefl.Get[int]: named access, path resolution per call
// Offset      = pre-computed offset + reflect.NewAt (internal mechanism saferefl uses)
// Accessor    = Accessor.Get: pre-bound, pointer arithmetic only
// Native      = direct struct field access (theoretical minimum)

// BenchmarkFieldRead_Reflect — stdlib reflect: FieldByName per call.
func BenchmarkFieldRead_Reflect(b *testing.B) {
	u := &benchUser{ID: 42}
	b.ResetTimer()
	for i := range b.N {
		sinkInt = int(reflect.ValueOf(u).Elem().FieldByName("ID").Int()) + i
	}
}

// BenchmarkFieldRead_ReflectFast — stdlib reflect: pre-cached Value + Field(0).
func BenchmarkFieldRead_ReflectFast(b *testing.B) {
	u := &benchUser{ID: 42}
	rv := reflect.ValueOf(u).Elem()
	b.ResetTimer()
	for i := range b.N {
		sinkInt = int(rv.Field(0).Int()) + i
	}
}

// BenchmarkFieldRead_SafeRefl — saferefl.Get[int]: named access per call.
func BenchmarkFieldRead_SafeRefl(b *testing.B) {
	u := &benchUser{ID: 42}
	b.ResetTimer()
	for i := range b.N {
		v, _ := saferefl.Get[int](u, "ID")
		sinkInt = v + i
	}
}

// BenchmarkFieldRead_Offset — pre-computed offset + reflect.NewAt.
func BenchmarkFieldRead_Offset(b *testing.B) {
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

// BenchmarkFieldRead_Accessor — Accessor.Get: pre-bound, pointer arithmetic only.
func BenchmarkFieldRead_Accessor(b *testing.B) {
	u := &benchUser{ID: 42}
	ptr := saferefl.UnsafePtrOf(u)
	b.ResetTimer()
	for i := range b.N {
		sinkInt = benchIDAccessor.Get(ptr) + i
	}
}

// BenchmarkFieldRead_Native — direct struct field read (theoretical minimum).
func BenchmarkFieldRead_Native(b *testing.B) {
	u := &benchUser{ID: 42}
	b.ResetTimer()
	for i := range b.N {
		sinkInt = u.ID + i
	}
}

// ---- SliceAt: saferefl.UnsafeSliceAt vs direct vs reflect ----

// BenchmarkSliceAt_SafeRefl — saferefl.UnsafeSliceAt: direct element pointer, no bounds check.
func BenchmarkSliceAt_SafeRefl(b *testing.B) {
	s := make([]int, 1000)
	for i := range s {
		s[i] = i
	}
	_ = saferefl.EnableAccel()
	b.ResetTimer()
	for i := range b.N {
		sinkIntPtr = saferefl.UnsafeSliceAt(s, i%1000)
	}
}

// BenchmarkSliceAt_Direct — &s[i]: bounds-checked element pointer.
func BenchmarkSliceAt_Direct(b *testing.B) {
	s := make([]int, 1000)
	for i := range s {
		s[i] = i
	}
	b.ResetTimer()
	for i := range b.N {
		sinkIntPtr = &s[i%1000]
	}
}

// BenchmarkSliceAt_Reflect — reflect.Value.Index: reflect-based element access.
func BenchmarkSliceAt_Reflect(b *testing.B) {
	s := make([]int, 1000)
	for i := range s {
		s[i] = i
	}
	rv := reflect.ValueOf(s)
	b.ResetTimer()
	for i := range b.N {
		sinkInt = int(rv.Index(i % 1000).Int())
	}
}

// ---- MapLen: saferefl.MapLenFast vs builtin vs reflect ----

func newBenchMap() map[string]int {
	m := make(map[string]int, 100)
	for i := range 100 {
		m[fmt.Sprintf("key%d", i)] = i
	}
	return m
}

// BenchmarkMapLen_SafeRefl — saferefl.MapLenFast: direct count read, no reflect.
func BenchmarkMapLen_SafeRefl(b *testing.B) {
	m := newBenchMap()
	_ = saferefl.EnableAccel()
	b.ResetTimer()
	for range b.N {
		sinkInt = saferefl.MapLenFast(m)
	}
}

// BenchmarkMapLen_Builtin — len(m): builtin map length.
func BenchmarkMapLen_Builtin(b *testing.B) {
	m := newBenchMap()
	b.ResetTimer()
	for range b.N {
		sinkInt = len(m)
	}
}

// BenchmarkMapLen_Reflect — reflect.Value.Len(): reflect-based map length.
func BenchmarkMapLen_Reflect(b *testing.B) {
	m := newBenchMap()
	rv := reflect.ValueOf(m)
	b.ResetTimer()
	for range b.N {
		sinkInt = rv.Len()
	}
}
