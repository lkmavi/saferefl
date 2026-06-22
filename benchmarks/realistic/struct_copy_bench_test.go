package realistic

import (
	"reflect"
	"testing"

	"github.com/jinzhu/copier"
	"github.com/lkmavi/saferefl"
)

var sinkDst UserDst

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

// BenchmarkStructCopy_Saferefl — Layer 1: Get[T]/Set[T] per field.
func BenchmarkStructCopy_Saferefl(b *testing.B) {
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

// BenchmarkStructCopy_Reflect — stdlib reflect field-by-field copy.
func BenchmarkStructCopy_Reflect(b *testing.B) {
	src := newUserSrc()
	dst := UserDst{}
	srcV := reflect.ValueOf(src)
	dstV := reflect.ValueOf(&dst).Elem()
	b.ResetTimer()
	for range b.N {
		for i := 0; i < dstV.NumField(); i++ {
			dstV.Field(i).Set(srcV.Field(i))
		}
		sinkDst = dst
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
