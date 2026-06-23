package realistic

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/jinzhu/copier"
	"github.com/lkmavi/saferefl"
	reflect2 "github.com/modern-go/reflect2"
)

var sinkRow ORMRow

// columns simulates the ordered list of column names returned by the driver.
var columns = []string{
	"Col1", "Col2", "Col3", "Col4", "Col5",
	"Col6", "Col7", "Col8", "Col9", "Col10",
}

// pre-bound Accessor per column — resolved at statement-prepare time, like a real ORM.
var (
	ormCol1  = mustMakeAccessor[int64](&ORMRow{}, "Col1")
	ormCol2  = mustMakeAccessor[string](&ORMRow{}, "Col2")
	ormCol3  = mustMakeAccessor[string](&ORMRow{}, "Col3")
	ormCol4  = mustMakeAccessor[float64](&ORMRow{}, "Col4")
	ormCol5  = mustMakeAccessor[int64](&ORMRow{}, "Col5")
	ormCol6  = mustMakeAccessor[bool](&ORMRow{}, "Col6")
	ormCol7  = mustMakeAccessor[string](&ORMRow{}, "Col7")
	ormCol8  = mustMakeAccessor[float64](&ORMRow{}, "Col8")
	ormCol9  = mustMakeAccessor[int64](&ORMRow{}, "Col9")
	ormCol10 = mustMakeAccessor[string](&ORMRow{}, "Col10")
)

// pre-compiled reflect2 field descriptors and source offsets — built once at query-prepare time.
var (
	ormType2   = reflect2.TypeOf(ORMRow{}).(reflect2.StructType)
	ormFields2 = buildORMFields2()
	ormSrcOff  = buildORMSrcOff()
)

func buildORMFields2() []reflect2.StructField {
	fields := make([]reflect2.StructField, len(columns))
	for i, name := range columns {
		fields[i] = ormType2.FieldByName(name)
	}
	return fields
}

func buildORMSrcOff() []uintptr {
	rt := reflect.TypeOf(ORMRow{})
	off := make([]uintptr, len(columns))
	for i, name := range columns {
		f, _ := rt.FieldByName(name)
		off[i] = f.Offset
	}
	return off
}

// pre-parsed row values (simulates sql.Rows.Scan having already decoded the wire data).
var rowValues = newORMValues()

// rowSrc is used as the "source" struct for the copier benchmark.
var rowSrc = ORMRow{
	Col1: 1, Col2: "Alice", Col3: "alice@example.com", Col4: 9.5, Col5: 100,
	Col6: true, Col7: "admin", Col8: 0.75, Col9: 42, Col10: "active",
}

// BenchmarkORMScan_Manual — direct field assignment, native baseline.
func BenchmarkORMScan_Manual(b *testing.B) {
	vals := rowValues
	b.ResetTimer()
	for range b.N {
		sinkRow = ORMRow{
			Col1:  vals[0].(int64),
			Col2:  vals[1].(string),
			Col3:  vals[2].(string),
			Col4:  vals[3].(float64),
			Col5:  vals[4].(int64),
			Col6:  vals[5].(bool),
			Col7:  vals[6].(string),
			Col8:  vals[7].(float64),
			Col9:  vals[8].(int64),
			Col10: vals[9].(string),
		}
	}
}

// BenchmarkORMScan_SafeRefl — Layer 1: Set[T] per column using pre-mapped field names.
func BenchmarkORMScan_SafeRefl(b *testing.B) {
	row := &ORMRow{}
	vals := rowValues
	b.ResetTimer()
	for range b.N {
		_ = saferefl.Set[int64](row, "Col1", vals[0].(int64))
		_ = saferefl.Set[string](row, "Col2", vals[1].(string))
		_ = saferefl.Set[string](row, "Col3", vals[2].(string))
		_ = saferefl.Set[float64](row, "Col4", vals[3].(float64))
		_ = saferefl.Set[int64](row, "Col5", vals[4].(int64))
		_ = saferefl.Set[bool](row, "Col6", vals[5].(bool))
		_ = saferefl.Set[string](row, "Col7", vals[6].(string))
		_ = saferefl.Set[float64](row, "Col8", vals[7].(float64))
		_ = saferefl.Set[int64](row, "Col9", vals[8].(int64))
		_ = saferefl.Set[string](row, "Col10", vals[9].(string))
		sinkRow = *row
	}
}

// BenchmarkORMScan_Accessor — Layer 3: pre-bound Accessor per column.
// Simulates real ORM: prepare bindings once per statement, scan every row in the hot loop.
func BenchmarkORMScan_Accessor(b *testing.B) {
	row := &ORMRow{}
	ptr := saferefl.UnsafePtrOf(row)
	vals := rowValues
	b.ResetTimer()
	for range b.N {
		ormCol1.Set(ptr, vals[0].(int64))
		ormCol2.Set(ptr, vals[1].(string))
		ormCol3.Set(ptr, vals[2].(string))
		ormCol4.Set(ptr, vals[3].(float64))
		ormCol5.Set(ptr, vals[4].(int64))
		ormCol6.Set(ptr, vals[5].(bool))
		ormCol7.Set(ptr, vals[6].(string))
		ormCol8.Set(ptr, vals[7].(float64))
		ormCol9.Set(ptr, vals[8].(int64))
		ormCol10.Set(ptr, vals[9].(string))
		sinkRow = *row
	}
}

// BenchmarkORMScan_Reflect2 — reflect2: pre-compiled field descriptors + UnsafeSet.
// Simulates a well-optimised ORM that caches reflect2 metadata at statement-prepare time.
func BenchmarkORMScan_Reflect2(b *testing.B) {
	src := rowSrc
	dst := &ORMRow{}
	srcPtr := unsafe.Pointer(&src)
	dstPtr := unsafe.Pointer(dst)
	b.ResetTimer()
	for range b.N {
		for i, f := range ormFields2 {
			f.UnsafeSet(dstPtr, unsafe.Pointer(uintptr(srcPtr)+ormSrcOff[i]))
		}
		sinkRow = *dst
	}
}

// BenchmarkORMScan_Reflect — stdlib reflect: FieldByName + Set per column.
func BenchmarkORMScan_Reflect(b *testing.B) {
	row := ORMRow{}
	rv := reflect.ValueOf(&row).Elem()
	vals := rowValues
	b.ResetTimer()
	for range b.N {
		for i, col := range columns {
			rv.FieldByName(col).Set(reflect.ValueOf(vals[i]))
		}
		sinkRow = row
	}
}

// BenchmarkORMScan_Copier — github.com/jinzhu/copier struct-to-struct.
func BenchmarkORMScan_Copier(b *testing.B) {
	dst := ORMRow{}
	b.ResetTimer()
	for range b.N {
		_ = copier.Copy(&dst, &rowSrc)
		sinkRow = dst
	}
}
