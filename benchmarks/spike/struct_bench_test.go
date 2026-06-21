package spike

import (
	"reflect"
	"testing"
	"unsafe"

	reflect2 "github.com/modern-go/reflect2"
)

// Sinks prevent the compiler from eliminating benchmark bodies as dead code.
var (
	sinkInt64  int64
	sinkString string
)

// --- Pre-computed metadata (simulates what Layer 2 TypeInfo cache would hold) ---

var (
	userReflType = reflect.TypeOf(User{})
	userType2    = reflect2.TypeOf(User{}).(reflect2.StructType)

	idSF     = mustStructField(userReflType, "ID")
	idOffset = idSF.Offset
	idRType  = idSF.Type
	idField2 = userType2.FieldByName("ID")

	nameSF     = mustStructField(userReflType, "Name")
	nameOffset = nameSF.Offset
	nameRType  = nameSF.Type
	nameField2 = userType2.FieldByName("Name")
)

func mustStructField(t reflect.Type, name string) reflect.StructField {
	f, ok := t.FieldByName(name)
	if !ok {
		panic("field not found: " + name)
	}
	return f
}

// readField simulates the Layer 1 / Layer 3 path: type known at compile time,
// offset pre-computed — compiler emits a direct load with no reflect overhead.
func readField[T any](objPtr unsafe.Pointer, offset uintptr) T {
	return *(*T)(unsafe.Pointer(uintptr(objPtr) + offset))
}

// writeField is the write counterpart of readField.
func writeField[T any](objPtr unsafe.Pointer, offset uintptr, val T) {
	*(*T)(unsafe.Pointer(uintptr(objPtr) + offset)) = val
}

// ============================================================
// READ: int64 field
// ============================================================

// BenchmarkReadInt64_StdlibReflect — full dynamic path: ValueOf + FieldByName on every call.
func BenchmarkReadInt64_StdlibReflect(b *testing.B) {
	u := newUser()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkInt64 = reflect.ValueOf(u).Elem().FieldByName("ID").Int()
	}
}

// BenchmarkReadInt64_Reflect2 — reflect2 UnsafeGet: field descriptor cached, direct pointer read.
func BenchmarkReadInt64_Reflect2(b *testing.B) {
	u := newUser()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkInt64 = *(*int64)(idField2.UnsafeGet(unsafe.Pointer(u)))
	}
}

// BenchmarkReadInt64_CachedOffset — Layer 2 path: pre-computed offset + reflect.NewAt (documented, safe).
func BenchmarkReadInt64_CachedOffset(b *testing.B) {
	u := newUser()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fPtr := unsafe.Pointer(uintptr(unsafe.Pointer(u)) + idOffset)
		sinkInt64 = reflect.NewAt(idRType, fPtr).Elem().Int()
	}
}

// BenchmarkReadInt64_Generics — Layer 1/3 path: pre-computed offset + direct cast, zero boxing.
func BenchmarkReadInt64_Generics(b *testing.B) {
	u := newUser()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkInt64 = readField[int64](unsafe.Pointer(u), idOffset)
	}
}

// ============================================================
// WRITE: int64 field
// ============================================================

func BenchmarkWriteInt64_StdlibReflect(b *testing.B) {
	u := newUser()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reflect.ValueOf(u).Elem().FieldByName("ID").SetInt(int64(i))
	}
}

func BenchmarkWriteInt64_Reflect2(b *testing.B) {
	u := newUser()
	newVal := int64(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		idField2.UnsafeSet(unsafe.Pointer(u), unsafe.Pointer(&newVal))
	}
}

func BenchmarkWriteInt64_CachedOffset(b *testing.B) {
	u := newUser()
	val := reflect.ValueOf(int64(100))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fPtr := unsafe.Pointer(uintptr(unsafe.Pointer(u)) + idOffset)
		reflect.NewAt(idRType, fPtr).Elem().Set(val)
	}
}

func BenchmarkWriteInt64_Generics(b *testing.B) {
	u := newUser()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		writeField[int64](unsafe.Pointer(u), idOffset, int64(i))
	}
}

// ============================================================
// READ: string field (string header is 2 words — different profile than int64)
// ============================================================

func BenchmarkReadString_StdlibReflect(b *testing.B) {
	u := newUser()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkString = reflect.ValueOf(u).Elem().FieldByName("Name").String()
	}
}

func BenchmarkReadString_Reflect2(b *testing.B) {
	u := newUser()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkString = *(*string)(nameField2.UnsafeGet(unsafe.Pointer(u)))
	}
}

func BenchmarkReadString_CachedOffset(b *testing.B) {
	u := newUser()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fPtr := unsafe.Pointer(uintptr(unsafe.Pointer(u)) + nameOffset)
		sinkString = reflect.NewAt(nameRType, fPtr).Elem().String()
	}
}

func BenchmarkReadString_Generics(b *testing.B) {
	u := newUser()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkString = readField[string](unsafe.Pointer(u), nameOffset)
	}
}
