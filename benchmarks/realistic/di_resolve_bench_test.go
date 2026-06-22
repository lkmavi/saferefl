package realistic

import (
	"reflect"
	"testing"

	"github.com/lkmavi/saferefl"
)

var sinkServices AppServices

// registry simulates a DI container's resolved-instance store.
var registry = map[string]interface{}{
	"A": &ServiceA{Name: "svcA"},
	"B": &ServiceB{Name: "svcB"},
	"C": &ServiceC{Name: "svcC"},
}

// pre-computed reflect.Values for the reflect benchmark.
var registryValues = map[string]reflect.Value{
	"A": reflect.ValueOf(registry["A"]),
	"B": reflect.ValueOf(registry["B"]),
	"C": reflect.ValueOf(registry["C"]),
}

// pre-bound Accessor bindings — resolved once at startup like a real DI container.
var (
	svcAAcc = mustMakeAccessor[*ServiceA](&AppServices{}, "A")
	svcBAcc = mustMakeAccessor[*ServiceB](&AppServices{}, "B")
	svcCAcc = mustMakeAccessor[*ServiceC](&AppServices{}, "C")
)

// BenchmarkDIResolve_Manual — direct pointer assignment, native baseline.
func BenchmarkDIResolve_Manual(b *testing.B) {
	a := registry["A"].(*ServiceA)
	bc := registry["B"].(*ServiceB)
	c := registry["C"].(*ServiceC)
	b.ResetTimer()
	for range b.N {
		sinkServices = AppServices{A: a, B: bc, C: c}
	}
}

// BenchmarkDIResolve_Saferefl — Layer 1: Set[*ServiceX] per dependency.
func BenchmarkDIResolve_Saferefl(b *testing.B) {
	a := registry["A"].(*ServiceA)
	bc := registry["B"].(*ServiceB)
	c := registry["C"].(*ServiceC)
	svc := &AppServices{}
	b.ResetTimer()
	for range b.N {
		_ = saferefl.Set[*ServiceA](svc, "A", a)
		_ = saferefl.Set[*ServiceB](svc, "B", bc)
		_ = saferefl.Set[*ServiceC](svc, "C", c)
		sinkServices = *svc
	}
}

// BenchmarkDIResolve_Accessor — Layer 3: pre-bound Accessor, inject per request.
// Simulates real DI: resolve bindings once at startup, inject on every request.
func BenchmarkDIResolve_Accessor(b *testing.B) {
	a := registry["A"].(*ServiceA)
	bc := registry["B"].(*ServiceB)
	c := registry["C"].(*ServiceC)
	svc := &AppServices{}
	ptr := saferefl.UnsafePtrOf(svc)
	b.ResetTimer()
	for range b.N {
		svcAAcc.Set(ptr, a)
		svcBAcc.Set(ptr, bc)
		svcCAcc.Set(ptr, c)
		sinkServices = *svc
	}
}

// BenchmarkDIResolve_Reflect — stdlib reflect field-by-name injection.
func BenchmarkDIResolve_Reflect(b *testing.B) {
	svc := AppServices{}
	rv := reflect.ValueOf(&svc).Elem()
	b.ResetTimer()
	for range b.N {
		rv.FieldByName("A").Set(registryValues["A"])
		rv.FieldByName("B").Set(registryValues["B"])
		rv.FieldByName("C").Set(registryValues["C"])
		sinkServices = svc
	}
}
