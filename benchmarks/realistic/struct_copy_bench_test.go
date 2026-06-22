package realistic

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/jinzhu/copier"
	"github.com/lkmavi/saferefl"
	reflect2 "github.com/modern-go/reflect2"
)

var sinkDst UserDst

// pre-bound Accessor pairs for the L3 path — resolved once at startup.
var (
	srcIDAccSC    = mustMakeAccessor[int](&UserSrc{}, "ID")
	srcNameAccSC  = mustMakeAccessor[string](&UserSrc{}, "Name")
	srcEmailAccSC = mustMakeAccessor[string](&UserSrc{}, "Email")
	srcScoreAccSC = mustMakeAccessor[float64](&UserSrc{}, "Score")
	srcActAccSC   = mustMakeAccessor[bool](&UserSrc{}, "Active")

	dstIDAccSC    = mustMakeAccessor[int](&UserDst{}, "ID")
	dstNameAccSC  = mustMakeAccessor[string](&UserDst{}, "Name")
	dstEmailAccSC = mustMakeAccessor[string](&UserDst{}, "Email")
	dstScoreAccSC = mustMakeAccessor[float64](&UserDst{}, "Score")
	dstActAccSC   = mustMakeAccessor[bool](&UserDst{}, "Active")
)

var copyFieldNames = []string{"ID", "Name", "Email", "Score", "Active"}

// pre-computed reflect2 field descriptors for the Reflect2 path.
var (
	srcType2     = reflect2.TypeOf(UserSrc{}).(reflect2.StructType)
	dstType2     = reflect2.TypeOf(UserDst{}).(reflect2.StructType)
	dstFields2SC = buildDstFields2SC()
	srcOffsets2  = buildSrcOffsets2()
)

func buildDstFields2SC() []reflect2.StructField {
	fields := make([]reflect2.StructField, len(copyFieldNames))
	for i, name := range copyFieldNames {
		fields[i] = dstType2.FieldByName(name)
	}
	return fields
}

func buildSrcOffsets2() []uintptr {
	rt := reflect.TypeOf(UserSrc{})
	offsets := make([]uintptr, len(copyFieldNames))
	for i, name := range copyFieldNames {
		f, _ := rt.FieldByName(name)
		offsets[i] = f.Offset
	}
	return offsets
}

// copyFieldMeta pairs the source and destination offsets for one field.
type copyFieldMeta struct {
	srcOffset uintptr
	dstOffset uintptr
	rtype     reflect.Type
}

// copyFieldMetas are built once — same strategy as the TypeInfo cache in Layer 2.
var copyFieldMetas = buildCopyFieldMetas()

func buildCopyFieldMetas() []copyFieldMeta {
	srcType := reflect.TypeOf(UserSrc{})
	dstType := reflect.TypeOf(UserDst{})
	metas := make([]copyFieldMeta, len(copyFieldNames))
	for i, name := range copyFieldNames {
		sf, _ := srcType.FieldByName(name)
		df, _ := dstType.FieldByName(name)
		metas[i] = copyFieldMeta{srcOffset: sf.Offset, dstOffset: df.Offset, rtype: sf.Type}
	}
	return metas
}

// BenchmarkStructCopy_Manual — direct struct literal assignment, the native baseline.
func BenchmarkStructCopy_Manual(b *testing.B) {
	src := newUserSrc()
	b.ResetTimer()
	for range b.N {
		sinkDst = UserDst{
			ID:     src.ID,
			Name:   src.Name,
			Email:  src.Email,
			Score:  src.Score,
			Active: src.Active,
		}
	}
}

// BenchmarkStructCopy_L3 — Layer 3: pre-bound Accessor pairs, direct pointer arithmetic.
// Simulates a real copy pipeline where bindings are resolved once at startup.
func BenchmarkStructCopy_L3(b *testing.B) {
	src := newUserSrc()
	dst := &UserDst{}
	srcPtr := saferefl.UnsafePtrOf(&src)
	dstPtr := saferefl.UnsafePtrOf(dst)
	b.ResetTimer()
	for range b.N {
		dstIDAccSC.Set(dstPtr, srcIDAccSC.Get(srcPtr))
		dstNameAccSC.Set(dstPtr, srcNameAccSC.Get(srcPtr))
		dstEmailAccSC.Set(dstPtr, srcEmailAccSC.Get(srcPtr))
		dstScoreAccSC.Set(dstPtr, srcScoreAccSC.Get(srcPtr))
		dstActAccSC.Set(dstPtr, srcActAccSC.Get(srcPtr))
		sinkDst = *dst
	}
}

// BenchmarkStructCopy_Reflect2 — reflect2: pre-compiled field descriptors + UnsafeSet.
// Represents a well-optimised copy pipeline that caches reflect2 metadata at startup.
func BenchmarkStructCopy_Reflect2(b *testing.B) {
	src := newUserSrc()
	dst := &UserDst{}
	srcPtr := unsafe.Pointer(&src)
	dstPtr := unsafe.Pointer(dst)
	b.ResetTimer()
	for range b.N {
		for i, f := range dstFields2SC {
			f.UnsafeSet(dstPtr, unsafe.Pointer(uintptr(srcPtr)+srcOffsets2[i]))
		}
		sinkDst = *dst
	}
}

// BenchmarkStructCopy_L2 — Layer 2 path: pre-computed offsets + reflect.NewAt per field.
// Slower than Reflect2 (reflect2.UnsafeSet avoids reflect.Value allocation), but faster
// than per-call FieldByName and entirely allocation-free.
func BenchmarkStructCopy_L2(b *testing.B) {
	src := newUserSrc()
	dst := &UserDst{}
	srcPtr := unsafe.Pointer(&src)
	dstPtr := unsafe.Pointer(dst)
	b.ResetTimer()
	for range b.N {
		for _, m := range copyFieldMetas {
			sp := unsafe.Pointer(uintptr(srcPtr) + m.srcOffset)
			dp := unsafe.Pointer(uintptr(dstPtr) + m.dstOffset)
			reflect.NewAt(m.rtype, dp).Elem().Set(reflect.NewAt(m.rtype, sp).Elem())
		}
		sinkDst = *dst
	}
}

// BenchmarkStructCopy_Reflect — stdlib reflect: FieldByName per field, per-call resolution.
// No pre-binding — simulates the common case where the caller resolves fields dynamically.
func BenchmarkStructCopy_Reflect(b *testing.B) {
	src := newUserSrc()
	dst := UserDst{}
	srcV := reflect.ValueOf(src)
	dstV := reflect.ValueOf(&dst).Elem()
	b.ResetTimer()
	for range b.N {
		for _, name := range copyFieldNames {
			dstV.FieldByName(name).Set(srcV.FieldByName(name))
		}
		sinkDst = dst
	}
}

// BenchmarkStructCopy_L1 — Layer 1: Get[T]/Set[T] per field with per-call name resolution.
func BenchmarkStructCopy_L1(b *testing.B) {
	src := newUserSrc()
	dst := &UserDst{}
	b.ResetTimer()
	for range b.N {
		_ = saferefl.Set[int](dst, "ID", saferefl.MustGet[int](&src, "ID"))
		_ = saferefl.Set[string](dst, "Name", saferefl.MustGet[string](&src, "Name"))
		_ = saferefl.Set[string](dst, "Email", saferefl.MustGet[string](&src, "Email"))
		_ = saferefl.Set[float64](dst, "Score", saferefl.MustGet[float64](&src, "Score"))
		_ = saferefl.Set[bool](dst, "Active", saferefl.MustGet[bool](&src, "Active"))
		sinkDst = *dst
	}
}

// BenchmarkStructCopy_Copier — github.com/jinzhu/copier.
func BenchmarkStructCopy_Copier(b *testing.B) {
	src := newUserSrc()
	dst := UserDst{}
	b.ResetTimer()
	for range b.N {
		_ = copier.Copy(&dst, &src)
		sinkDst = dst
	}
}
