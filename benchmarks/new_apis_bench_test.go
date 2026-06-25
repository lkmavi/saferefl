package benchmarks

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/lkmavi/saferefl"
	reflect2 "github.com/modern-go/reflect2"
)

// Benchmark struct with 5 exported fields — representative for DTO mapping.
type userA struct {
	Name   string
	Age    int
	Email  string
	Score  float64
	Active bool
}

type userB struct {
	Name   string
	Age    int
	Email  string
	Score  float64
	Active bool
}

// tagged version for ToMapByTag benchmarks
type userTagged struct {
	Name   string  `json:"name"`
	Age    int     `json:"age"`
	Email  string  `json:"email"`
	Score  float64 `json:"score"`
	Active bool    `json:"active"`
}

var (
	sinkMap  map[string]any
	sinkBool bool
	sinkKind reflect.Kind
)

// ── KindOf ──────────────────────────────────────────────────────────────────

// BenchmarkKindOf_SafeRefl — saferefl.KindOf: abi.Type byte read, no alloc.
func BenchmarkKindOf_SafeRefl(b *testing.B) {
	u := &userA{Name: "Alice"}
	b.ResetTimer()
	for range b.N {
		sinkKind = saferefl.KindOf(u)
	}
}

// BenchmarkKindOf_Reflect — reflect.TypeOf(v).Kind(): full type lookup.
func BenchmarkKindOf_Reflect(b *testing.B) {
	u := &userA{Name: "Alice"}
	b.ResetTimer()
	for range b.N {
		sinkKind = reflect.TypeOf(u).Kind()
	}
}

// ── IsNil ────────────────────────────────────────────────────────────────────

// BenchmarkIsNil_ptr_SafeRefl — saferefl.IsNil on a non-nil pointer.
func BenchmarkIsNil_ptr_SafeRefl(b *testing.B) {
	u := &userA{}
	b.ResetTimer()
	for range b.N {
		sinkBool = saferefl.IsNil(u)
	}
}

// BenchmarkIsNil_ptr_Reflect — reflect.ValueOf(v).IsNil() on pointer.
func BenchmarkIsNil_ptr_Reflect(b *testing.B) {
	u := &userA{}
	b.ResetTimer()
	for range b.N {
		sinkBool = reflect.ValueOf(u).IsNil()
	}
}

// ── EachField ────────────────────────────────────────────────────────────────

// BenchmarkEachField_SafeRefl — saferefl.EachField: cached offsets, direct reads.
func BenchmarkEachField_SafeRefl(b *testing.B) {
	u := &userA{Name: "Alice", Age: 30, Email: "a@b.com", Score: 9.5, Active: true}
	b.ResetTimer()
	for range b.N {
		_ = saferefl.EachField(u, func(_ string, v any) bool {
			_ = v
			return true
		})
	}
}

// BenchmarkEachField_Reflect — reflect with pre-computed reflect.Value (best-case reflect;
// not apples-to-apples with EachField which resolves the struct pointer every call).
func BenchmarkEachField_Reflect(b *testing.B) {
	u := &userA{Name: "Alice", Age: 30, Email: "a@b.com", Score: 9.5, Active: true}
	rt := reflect.TypeOf(*u)
	rv := reflect.ValueOf(u).Elem()
	b.ResetTimer()
	for range b.N {
		for i := range rt.NumField() {
			_ = rv.Field(i).Interface()
		}
	}
}

// BenchmarkEachField_ReflectFull — reflect with ValueOf inside the loop:
// apples-to-apples with saferefl.EachField which resolves the pointer every call.
func BenchmarkEachField_ReflectFull(b *testing.B) {
	u := &userA{Name: "Alice", Age: 30, Email: "a@b.com", Score: 9.5, Active: true}
	rt := reflect.TypeOf(*u)
	b.ResetTimer()
	for range b.N {
		rv := reflect.ValueOf(u).Elem()
		for i := range rt.NumField() {
			_ = rv.Field(i).Interface()
		}
	}
}

// BenchmarkEachField_Reflect2 — reflect2: pre-compiled field list, UnsafeGet per field.
func BenchmarkEachField_Reflect2(b *testing.B) {
	u := &userA{Name: "Alice", Age: 30, Email: "a@b.com", Score: 9.5, Active: true}
	t2 := reflect2.TypeOf(*u).(*reflect2.UnsafeStructType)
	fields := make([]reflect2.StructField, t2.NumField())
	for i := range fields {
		fields[i] = t2.Field(i)
	}
	b.ResetTimer()
	for range b.N {
		for _, f := range fields {
			_ = f.Get(u)
		}
	}
}

// ── CopyFields ───────────────────────────────────────────────────────────────

// BenchmarkCopyFields_SafeRefl — saferefl.CopyFields: cached FieldsByName + reflect.Set.
func BenchmarkCopyFields_SafeRefl(b *testing.B) {
	src := &userA{Name: "Alice", Age: 30, Email: "a@b.com", Score: 9.5, Active: true}
	dst := &userB{}
	b.ResetTimer()
	for range b.N {
		_ = saferefl.CopyFields(src, dst)
	}
}

// BenchmarkCopyFields_Manual — direct field assignment (theoretical minimum).
func BenchmarkCopyFields_Manual(b *testing.B) {
	src := &userA{Name: "Alice", Age: 30, Email: "a@b.com", Score: 9.5, Active: true}
	dst := &userB{}
	b.ResetTimer()
	for range b.N {
		dst.Name = src.Name
		dst.Age = src.Age
		dst.Email = src.Email
		dst.Score = src.Score
		dst.Active = src.Active
	}
}

// BenchmarkCopyFields_Reflect — reflect: FieldByName per field, Value.Set.
func BenchmarkCopyFields_Reflect(b *testing.B) {
	src := &userA{Name: "Alice", Age: 30, Email: "a@b.com", Score: 9.5, Active: true}
	dst := &userB{}
	srcRT := reflect.TypeOf(*src)
	b.ResetTimer()
	for range b.N {
		srcRV := reflect.ValueOf(src).Elem()
		dstRV := reflect.ValueOf(dst).Elem()
		for i := range srcRT.NumField() {
			sf := srcRT.Field(i)
			dstRV.FieldByName(sf.Name).Set(srcRV.Field(i))
		}
	}
}

// ── ToMap ────────────────────────────────────────────────────────────────────

// BenchmarkToMap_SafeRefl — saferefl.ToMap: EachField + map insert.
func BenchmarkToMap_SafeRefl(b *testing.B) {
	u := &userA{Name: "Alice", Age: 30, Email: "a@b.com", Score: 9.5, Active: true}
	b.ResetTimer()
	for range b.N {
		m, _ := saferefl.ToMap(u)
		sinkMap = m
	}
}

// BenchmarkToMap_Reflect — reflect: NumField loop + Field(i).Interface().
func BenchmarkToMap_Reflect(b *testing.B) {
	u := &userA{Name: "Alice", Age: 30, Email: "a@b.com", Score: 9.5, Active: true}
	rt := reflect.TypeOf(*u)
	b.ResetTimer()
	for range b.N {
		rv := reflect.ValueOf(u).Elem()
		m := make(map[string]any, rt.NumField())
		for i := range rt.NumField() {
			m[rt.Field(i).Name] = rv.Field(i).Interface()
		}
		sinkMap = m
	}
}

// BenchmarkToMap_JSON — json.Marshal + json.Unmarshal (full-serialize baseline).
func BenchmarkToMap_JSON(b *testing.B) {
	u := &userTagged{Name: "Alice", Age: 30, Email: "a@b.com", Score: 9.5, Active: true}
	b.ResetTimer()
	for range b.N {
		data, _ := json.Marshal(u)
		m := map[string]any{}
		_ = json.Unmarshal(data, &m)
		sinkMap = m
	}
}

// ── MapForEach ───────────────────────────────────────────────────────────────

// BenchmarkMapForEach_SafeRefl — saferefl.MapForEach: range with fn callback.
func BenchmarkMapForEach_SafeRefl(b *testing.B) {
	m := make(map[string]int, 100)
	for i := range 100 {
		m[benchKey(i)] = i
	}
	b.ResetTimer()
	for range b.N {
		saferefl.MapForEach(m, func(_ string, v int) bool {
			sinkInt = v
			return true
		})
	}
}

// BenchmarkMapForEach_Range — plain range loop (theoretical minimum).
func BenchmarkMapForEach_Range(b *testing.B) {
	m := make(map[string]int, 100)
	for i := range 100 {
		m[benchKey(i)] = i
	}
	b.ResetTimer()
	for range b.N {
		for _, v := range m {
			sinkInt = v
		}
	}
}

func benchKey(i int) string {
	// cheap key generation without fmt.Sprintf alloc
	const prefix = "key"
	buf := [8]byte{}
	copy(buf[:], prefix)
	n := len(prefix)
	if i == 0 {
		buf[n] = '0'
		return string(buf[:n+1])
	}
	d := i
	digits := 0
	for d > 0 {
		digits++
		d /= 10
	}
	for j := n + digits - 1; j >= n; j-- {
		buf[j] = byte('0' + i%10)
		i /= 10
	}
	return string(buf[:n+digits])
}
