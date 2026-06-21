package spike

import (
	"reflect"
	"testing"
	"unsafe"

	reflect2 "github.com/modern-go/reflect2"
)

// fieldEntry holds pre-computed metadata for one struct field.
// In a real codec, this is built once during codec initialization — not per call.
type fieldEntry struct {
	name   string
	idx    int
	offset uintptr
	rtype  reflect.Type
	field2 reflect2.StructField
}

var userFieldEntries = buildUserFieldEntries()

func buildUserFieldEntries() []fieldEntry {
	t := reflect.TypeOf(User{})
	t2 := reflect2.TypeOf(User{}).(reflect2.StructType)
	names := []string{"ID", "Name", "Email", "Score", "Active"}
	entries := make([]fieldEntry, len(names))
	for i, name := range names {
		sf, ok := t.FieldByName(name)
		if !ok {
			panic("field not found: " + name)
		}
		entries[i] = fieldEntry{
			name:   name,
			idx:    sf.Index[0],
			offset: sf.Offset,
			rtype:  sf.Type,
			field2: t2.FieldByName(name),
		}
	}
	return entries
}

// Pre-parsed values — simulates tokens already decoded from the wire format.
var preReflValues = []reflect.Value{
	reflect.ValueOf(int64(42)),
	reflect.ValueOf("Alice"),
	reflect.ValueOf("alice@example.com"),
	reflect.ValueOf(float64(9.5)),
	reflect.ValueOf(true),
}

// src holds the pre-parsed typed values for unsafe/reflect2 paths.
var src = struct {
	id     int64
	name   string
	email  string
	score  float64
	active bool
}{42, "Alice", "alice@example.com", 9.5, true}

// fieldSetters are type-specific write closures for the generics path.
// In a generated codec, each setter is specialised — no runtime type dispatch.
var fieldSetters = []func(structPtr unsafe.Pointer){
	func(p unsafe.Pointer) { writeField[int64](p, userFieldEntries[0].offset, src.id) },
	func(p unsafe.Pointer) { writeField[string](p, userFieldEntries[1].offset, src.name) },
	func(p unsafe.Pointer) { writeField[string](p, userFieldEntries[2].offset, src.email) },
	func(p unsafe.Pointer) { writeField[float64](p, userFieldEntries[3].offset, src.score) },
	func(p unsafe.Pointer) { writeField[bool](p, userFieldEntries[4].offset, src.active) },
}

// ============================================================
// JSON-like decode: set 5 fields of a freshly allocated struct.
// Each iteration allocates a new User to simulate per-object decoding.
// ============================================================

// BenchmarkJSONLike_StdlibReflect — naive path: FieldByName per field, no metadata cache.
func BenchmarkJSONLike_StdlibReflect(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u := User{}
		rv := reflect.ValueOf(&u).Elem()
		for j, entry := range userFieldEntries {
			rv.FieldByName(entry.name).Set(preReflValues[j])
		}
		sinkInt64 = u.ID
	}
}

// BenchmarkJSONLike_Reflect2 — reflect2 UnsafeSet with pre-cached field descriptors.
func BenchmarkJSONLike_Reflect2(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u := User{}
		uPtr := unsafe.Pointer(&u)
		userFieldEntries[0].field2.UnsafeSet(uPtr, unsafe.Pointer(&src.id))
		userFieldEntries[1].field2.UnsafeSet(uPtr, unsafe.Pointer(&src.name))
		userFieldEntries[2].field2.UnsafeSet(uPtr, unsafe.Pointer(&src.email))
		userFieldEntries[3].field2.UnsafeSet(uPtr, unsafe.Pointer(&src.score))
		userFieldEntries[4].field2.UnsafeSet(uPtr, unsafe.Pointer(&src.active))
		sinkInt64 = u.ID
	}
}

// BenchmarkJSONLike_CachedOffset — Layer 2 path: pre-computed offsets + reflect.NewAt (safe).
func BenchmarkJSONLike_CachedOffset(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u := User{}
		uPtr := unsafe.Pointer(&u)
		for j, entry := range userFieldEntries {
			fPtr := unsafe.Pointer(uintptr(uPtr) + entry.offset)
			reflect.NewAt(entry.rtype, fPtr).Elem().Set(preReflValues[j])
		}
		sinkInt64 = u.ID
	}
}

// BenchmarkJSONLike_Generics — Layer 3 path: pre-computed offsets + type-specific setters, zero boxing.
func BenchmarkJSONLike_Generics(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u := User{}
		uPtr := unsafe.Pointer(&u)
		for _, setter := range fieldSetters {
			setter(uPtr)
		}
		sinkInt64 = u.ID
	}
}
