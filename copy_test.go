package saferefl_test

import (
	"testing"

	"github.com/lkmavi/saferefl"
)

func TestCopyFields_sameType(t *testing.T) {
	src := &person{Name: "Alice", Age: 30, Score: 9.5, Active: true}
	dst := &person{}
	if err := saferefl.CopyFields(src, dst); err != nil {
		t.Fatalf("CopyFields: %v", err)
	}
	if dst.Name != "Alice" || dst.Age != 30 || dst.Score != 9.5 || !dst.Active {
		t.Errorf("CopyFields result = %+v, want %+v", dst, src)
	}
}

func TestCopyFields_differentTypes_matchingFields(t *testing.T) {
	type srcT struct {
		Name  string
		Age   int
		Extra string // not in dst
	}
	type dstT struct {
		Name    string
		Age     int
		Missing int // not in src
	}
	src := &srcT{Name: "Bob", Age: 25, Extra: "ignored"}
	dst := &dstT{}
	if err := saferefl.CopyFields(src, dst); err != nil {
		t.Fatalf("CopyFields: %v", err)
	}
	if dst.Name != "Bob" || dst.Age != 25 {
		t.Errorf("dst = %+v", dst)
	}
	if dst.Missing != 0 {
		t.Error("Missing should be unchanged (zero)")
	}
}

func TestCopyFields_assignableType(t *testing.T) {
	type srcT struct{ N int32 }
	type dstT struct{ N int64 } // int32 ConvertibleTo int64
	src := &srcT{N: 42}
	dst := &dstT{}
	if err := saferefl.CopyFields(src, dst); err != nil {
		t.Fatalf("CopyFields: %v", err)
	}
	if dst.N != 42 {
		t.Errorf("N = %d, want 42", dst.N)
	}
}

func TestCopyFields_skipIncompatible(t *testing.T) {
	type srcT struct{ V string }
	type dstT struct{ V []int } // string not convertible to []int
	src := &srcT{V: "hello"}
	dst := &dstT{}
	// Should silently skip, not error.
	if err := saferefl.CopyFields(src, dst); err != nil {
		t.Fatalf("CopyFields unexpected error: %v", err)
	}
	if dst.V != nil {
		t.Error("incompatible field should be left unchanged")
	}
}

func TestCopyFields_promoted(t *testing.T) {
	type base struct{ Name string }
	type srcT struct {
		base
		Age int
	}
	type dstT struct {
		Name string // promoted in src, direct in dst
		Age  int
	}
	src := &srcT{base: base{Name: "Carol"}, Age: 22}
	dst := &dstT{}
	if err := saferefl.CopyFields(src, dst); err != nil {
		t.Fatalf("CopyFields: %v", err)
	}
	if dst.Name != "Carol" || dst.Age != 22 {
		t.Errorf("dst = %+v", dst)
	}
}

func TestCopyFields_nilSrc(t *testing.T) {
	if err := saferefl.CopyFields(nil, &person{}); err == nil {
		t.Error("expected error for nil src")
	}
}

func TestCopyFields_nilDst(t *testing.T) {
	if err := saferefl.CopyFields(&person{}, nil); err == nil {
		t.Error("expected error for nil dst")
	}
}

func TestCopyFields_nonPtrSrc(t *testing.T) {
	if err := saferefl.CopyFields(person{}, &person{}); err == nil {
		t.Error("expected error for non-pointer src")
	}
}
