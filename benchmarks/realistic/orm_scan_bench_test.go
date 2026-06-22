package realistic

import (
	"reflect"
	"testing"

	"github.com/jinzhu/copier"
	"github.com/lkmavi/saferefl"
)

var sinkRow ORMRow

// columns simulates the ordered list of column names returned by the driver.
var columns = []string{
	"Col1", "Col2", "Col3", "Col4", "Col5",
	"Col6", "Col7", "Col8", "Col9", "Col10",
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

// BenchmarkORMScan_Saferefl — Layer 1: Set[T] per column using pre-mapped field names.
func BenchmarkORMScan_Saferefl(b *testing.B) {
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
