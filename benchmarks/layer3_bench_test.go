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

// ---- FieldRead: all layers compared for a single int field ----
//
// Layer2 (ReflectNewAt) = cached offset + reflect.NewAt per read — original design
// Layer1 (Get)          = saferefl.Get[int] with full path resolution each call
// Layer3 (Accessor)     = Accessor.Get with pre-extracted unsafe.Pointer
// Native                = direct struct field access (theoretical minimum)

func BenchmarkFieldRead_Reflect(b *testing.B) {
	u := &benchUser{ID: 42}
	rv := reflect.ValueOf(u).Elem()
	b.ResetTimer()
	for i := range b.N {
		sinkInt = int(rv.Field(0).Int()) + i
	}
}

func BenchmarkFieldRead_ReflectNewAt(b *testing.B) {
	u := &benchUser{ID: 42}
	rt := reflect.TypeOf(u.ID)
	offset := reflect.TypeOf(*u).Field(0).Offset
	ptr := saferefl.UnsafePtrOf(u)
	b.ResetTimer()
	for i := range b.N {
		sinkInt = int(reflect.NewAt(rt, unsafe.Add(ptr, offset)).Elem().Int()) + i
	}
}

func BenchmarkFieldRead_Get(b *testing.B) {
	u := &benchUser{ID: 42}
	b.ResetTimer()
	for i := range b.N {
		v, _ := saferefl.Get[int](u, "ID")
		sinkInt = v + i
	}
}

func BenchmarkFieldRead_Accessor(b *testing.B) {
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

func BenchmarkSliceAt_Layer3(b *testing.B) {
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

func BenchmarkMapLen_Layer3(b *testing.B) {
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
