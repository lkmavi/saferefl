package benchmarks

import (
	"reflect"
	"testing"

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
