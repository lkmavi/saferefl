package realistic

import (
	"encoding/json"
	"reflect"
	"testing"
	"unsafe"

	"github.com/lkmavi/saferefl"
	reflect2 "github.com/modern-go/reflect2"
)

var (
	sinkProduct Product
	sinkInt64   int64
	jsonPayload []byte
)

func init() {
	var err error
	jsonPayload, err = json.Marshal(newProduct())
	if err != nil {
		panic(err)
	}
}

// pre-computed metadata for reflect2 path (simulates codec initialisation).
var (
	productType2  = reflect2.TypeOf(Product{}).(reflect2.StructType)
	productFields = buildProductFields()
)

type productFieldMeta struct {
	field2 reflect2.StructField
	offset uintptr
	rtype  reflect.Type
}

func buildProductFields() []productFieldMeta {
	names := []string{
		"ID", "SKU", "Title", "Description", "Price",
		"Stock", "Weight", "Active", "Category", "Tags",
	}
	rt := reflect.TypeOf(Product{})
	entries := make([]productFieldMeta, len(names))
	for i, name := range names {
		sf, _ := rt.FieldByName(name)
		entries[i] = productFieldMeta{
			field2: productType2.FieldByName(name),
			offset: sf.Offset,
			rtype:  sf.Type,
		}
	}
	return entries
}

// pre-parsed source values matching newProduct().
var productSrc = newProduct()

// BenchmarkJSONDecode_StdlibJSON — full encoding/json.Unmarshal including parsing.
func BenchmarkJSONDecode_StdlibJSON(b *testing.B) {
	b.ResetTimer()
	for range b.N {
		var p Product
		_ = json.Unmarshal(jsonPayload, &p)
		sinkInt64 = p.ID
	}
}

// BenchmarkJSONDecode_Saferefl — Layer 1 Set[T] per field (pre-parsed values).
func BenchmarkJSONDecode_Saferefl(b *testing.B) {
	dst := &Product{}
	b.ResetTimer()
	for range b.N {
		_ = saferefl.Set[int64](dst, "ID", productSrc.ID)
		_ = saferefl.Set[string](dst, "SKU", productSrc.SKU)
		_ = saferefl.Set[string](dst, "Title", productSrc.Title)
		_ = saferefl.Set[string](dst, "Description", productSrc.Description)
		_ = saferefl.Set[float64](dst, "Price", productSrc.Price)
		_ = saferefl.Set[int](dst, "Stock", productSrc.Stock)
		_ = saferefl.Set[float64](dst, "Weight", productSrc.Weight)
		_ = saferefl.Set[bool](dst, "Active", productSrc.Active)
		_ = saferefl.Set[string](dst, "Category", productSrc.Category)
		_ = saferefl.Set[string](dst, "Tags", productSrc.Tags)
		sinkInt64 = dst.ID
	}
}

// BenchmarkJSONDecode_Reflect — stdlib reflect SetField per field (pre-parsed values).
func BenchmarkJSONDecode_Reflect(b *testing.B) {
	srcV := reflect.ValueOf(productSrc)
	dst := Product{}
	dstV := reflect.ValueOf(&dst).Elem()
	b.ResetTimer()
	for range b.N {
		for i := 0; i < dstV.NumField(); i++ {
			dstV.Field(i).Set(srcV.Field(i))
		}
		sinkInt64 = dst.ID
	}
}

// BenchmarkJSONDecode_Reflect2 — reflect2 UnsafeSet per field with pre-cached descriptors (pre-parsed values).
func BenchmarkJSONDecode_Reflect2(b *testing.B) {
	dst := Product{}
	dstPtr := unsafe.Pointer(&dst)
	src := productSrc
	srcPtr := unsafe.Pointer(&src)
	b.ResetTimer()
	for range b.N {
		for _, f := range productFields {
			f.field2.UnsafeSet(dstPtr, unsafe.Pointer(uintptr(srcPtr)+f.offset))
		}
		sinkInt64 = dst.ID
	}
}
