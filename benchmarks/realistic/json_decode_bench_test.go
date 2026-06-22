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

// fieldNames is the ordered list of field names — simulates JSON key mapping built at codec-init time.
var productFieldNames = []string{
	"ID", "SKU", "Title", "Description", "Price",
	"Stock", "Weight", "Active", "Category", "Tags",
}

// pre-computed reflect2 field descriptors — built once at startup, like a real codec.
var (
	productType2    = reflect2.TypeOf(Product{}).(reflect2.StructType)
	productFields2  = buildProductFields2()
	productOffsets  = buildProductOffsets()
)

func buildProductFields2() []reflect2.StructField {
	fields := make([]reflect2.StructField, len(productFieldNames))
	for i, name := range productFieldNames {
		fields[i] = productType2.FieldByName(name)
	}
	return fields
}

func buildProductOffsets() []uintptr {
	rt := reflect.TypeOf(Product{})
	offsets := make([]uintptr, len(productFieldNames))
	for i, name := range productFieldNames {
		f, _ := rt.FieldByName(name)
		offsets[i] = f.Offset
	}
	return offsets
}

// pre-parsed source values matching newProduct().
var productSrc = newProduct()

// pre-bound Accessors for the L3 path — resolved once at startup, like a real codec.
var (
	pdIDAccJD       = mustMakeAccessor[int64](&Product{}, "ID")
	pdSKUAccJD      = mustMakeAccessor[string](&Product{}, "SKU")
	pdTitleAccJD    = mustMakeAccessor[string](&Product{}, "Title")
	pdDescAccJD     = mustMakeAccessor[string](&Product{}, "Description")
	pdPriceAccJD    = mustMakeAccessor[float64](&Product{}, "Price")
	pdStockAccJD    = mustMakeAccessor[int](&Product{}, "Stock")
	pdWeightAccJD   = mustMakeAccessor[float64](&Product{}, "Weight")
	pdActiveAccJD   = mustMakeAccessor[bool](&Product{}, "Active")
	pdCategoryAccJD = mustMakeAccessor[string](&Product{}, "Category")
	pdTagsAccJD     = mustMakeAccessor[string](&Product{}, "Tags")
)

// BenchmarkJSONDecode_StdlibJSON — full encoding/json.Unmarshal including JSON parsing.
func BenchmarkJSONDecode_StdlibJSON(b *testing.B) {
	b.ResetTimer()
	for range b.N {
		var p Product
		_ = json.Unmarshal(jsonPayload, &p)
		sinkInt64 = p.ID
	}
}

// BenchmarkJSONDecode_L1 — Layer 1: Set[T] per field with per-call name resolution.
// Equivalent in cost to Reflect: both resolve field names dynamically per call.
func BenchmarkJSONDecode_L1(b *testing.B) {
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

// BenchmarkJSONDecode_Reflect — stdlib reflect: FieldByName per field, per-call resolution.
// Same cost model as L1 — both resolve field names on every call with no pre-binding.
func BenchmarkJSONDecode_Reflect(b *testing.B) {
	dst := Product{}
	dstV := reflect.ValueOf(&dst).Elem()
	src := productSrc
	srcV := reflect.ValueOf(src)
	b.ResetTimer()
	for range b.N {
		for _, name := range productFieldNames {
			dstV.FieldByName(name).Set(srcV.FieldByName(name))
		}
		sinkInt64 = dst.ID
	}
}

// BenchmarkJSONDecode_Reflect2 — reflect2: pre-compiled field descriptors + UnsafeSet.
// Represents a well-optimised codec that caches reflect2 metadata at startup.
func BenchmarkJSONDecode_Reflect2(b *testing.B) {
	dst := Product{}
	dstPtr := unsafe.Pointer(&dst)
	src := productSrc
	srcPtr := unsafe.Pointer(&src)
	b.ResetTimer()
	for range b.N {
		for i, f := range productFields2 {
			f.UnsafeSet(dstPtr, unsafe.Pointer(uintptr(srcPtr)+productOffsets[i]))
		}
		sinkInt64 = dst.ID
	}
}

// BenchmarkJSONDecode_L3 — Layer 3: pre-bound Accessor per field, pointer arithmetic only.
// Represents a generated/pre-compiled codec where field bindings are resolved once at startup.
func BenchmarkJSONDecode_L3(b *testing.B) {
	dst := &Product{}
	ptr := saferefl.UnsafePtrOf(dst)
	src := productSrc
	b.ResetTimer()
	for range b.N {
		pdIDAccJD.Set(ptr, src.ID)
		pdSKUAccJD.Set(ptr, src.SKU)
		pdTitleAccJD.Set(ptr, src.Title)
		pdDescAccJD.Set(ptr, src.Description)
		pdPriceAccJD.Set(ptr, src.Price)
		pdStockAccJD.Set(ptr, src.Stock)
		pdWeightAccJD.Set(ptr, src.Weight)
		pdActiveAccJD.Set(ptr, src.Active)
		pdCategoryAccJD.Set(ptr, src.Category)
		pdTagsAccJD.Set(ptr, src.Tags)
		sinkInt64 = dst.ID
	}
}
