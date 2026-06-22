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

// ---- FieldRead: layer comparison for a single int field ----
//
// Reflect = stdlib reflect with pre-cached reflect.Value (fastest reflect path)
// L1      = saferefl.Get[int] — full path resolution per call
// L3      = Accessor.Get — pre-bound, pointer arithmetic only
// Native  = direct struct field access (theoretical minimum)

func BenchmarkFieldRead_Reflect(b *testing.B) {
	u := &benchUser{ID: 42}
	rv := reflect.ValueOf(u).Elem()
	b.ResetTimer()
	for i := range b.N {
		sinkInt = int(rv.Field(0).Int()) + i
	}
}

func BenchmarkFieldRead_L1(b *testing.B) {
	u := &benchUser{ID: 42}
	b.ResetTimer()
	for i := range b.N {
		v, _ := saferefl.Get[int](u, "ID")
		sinkInt = v + i
	}
}

func BenchmarkFieldRead_L2(b *testing.B) {
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

func BenchmarkFieldRead_L3(b *testing.B) {
	u := &benchUser{ID: 42}
	ptr := saferefl.UnsafePtrOf(u)
	b.ResetTimer()
	for i := range b.N {
		sinkInt = benchIDAccessor.Get(ptr) + i
	}
}

func BenchmarkFieldRead_Native(b *testing.B) {
	u := &benchUser{ID: 42}
	b.ResetTimer()
	for i := range b.N {
		sinkInt = u.ID + i
	}
}

// ---- SliceAt: Layer 3 vs direct vs reflect ----

func BenchmarkSliceAt_L3(b *testing.B) {
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

// ---- MapLen: Layer 3 vs builtin vs reflect ----

func newBenchMap() map[string]int {
	m := make(map[string]int, 100)
	for i := range 100 {
		m[fmt.Sprintf("key%d", i)] = i
	}
	return m
}

func BenchmarkMapLen_L3(b *testing.B) {
	m := newBenchMap()
	_ = saferefl.EnableAccel()
	b.ResetTimer()
	for range b.N {
		sinkInt = saferefl.MapLenFast(m)
	}
}

func BenchmarkMapLen_Builtin(b *testing.B) {
	m := newBenchMap()
	b.ResetTimer()
	for range b.N {
		sinkInt = len(m)
	}
}

func BenchmarkMapLen_Reflect(b *testing.B) {
	m := newBenchMap()
	rv := reflect.ValueOf(m)
	b.ResetTimer()
	for range b.N {
		sinkInt = rv.Len()
	}
}
