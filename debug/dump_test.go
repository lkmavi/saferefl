package debug

import (
	"bytes"
	"errors"
	"reflect"
	"strings"
	"testing"
	"unsafe"
)

type dumpSample struct {
	ID     int64
	Name   string
	Score  float64
	Active bool
	Tags   []string
}

func TestStructDump_basic(t *testing.T) {
	s := &dumpSample{ID: 42, Name: "Alice", Score: 9.5, Active: true}
	var buf bytes.Buffer
	if err := StructDump(s, &buf); err != nil {
		t.Fatalf("StructDump: %v", err)
	}
	out := buf.String()
	// Only the first field that starts in a given 16-byte row gets annotated.
	// ID (offset 0) is always first in row 0.
	// Score is always first in row 1: on 64-bit string is 16 bytes so Score lands at offset 24;
	// on 32-bit string is 8 bytes so Score lands at offset 16 — both are the first field in
	// the 16–31 byte row.
	// Active (offset 32 on 64-bit) gets its own row only on 64-bit; on 32-bit it shares
	// the 16–31 row with Score and is therefore not annotated.
	must := []string{"dumpSample", "ID", "Score"}
	if unsafe.Sizeof(uintptr(0)) == 8 {
		must = append(must, "Active")
	}
	for _, want := range must {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q:\n%s", want, out)
		}
	}
}

func TestStructDump_kindLabels(t *testing.T) {
	type withVariousKinds struct {
		Str   string
		Sli   []int
		Ptr   *int
		Iface any
	}
	s := &withVariousKinds{Str: "x"}
	var buf bytes.Buffer
	if err := StructDump(s, &buf); err != nil {
		t.Fatalf("StructDump: %v", err)
	}
	out := buf.String()
	// Str (offset 0) and Ptr are always first in their respective rows on all architectures.
	// On 64-bit: string=16 B, slice=24 B, *int=8 B, any=16 B — each field starts a new row.
	// On 32-bit: string=8 B and slice=12 B share row 0, *int=4 B is first in row 1,
	//   and any=8 B shares that row with *int — so "[]int" and "interface" are not annotated.
	must := []string{"string", "*int"}
	if unsafe.Sizeof(uintptr(0)) == 8 {
		must = append(must, "[]int", "interface")
	}
	for _, want := range must {
		if !strings.Contains(out, want) {
			t.Errorf("output missing kind label %q:\n%s", want, out)
		}
	}
}

func TestKindLabel_mapAndStruct(t *testing.T) {
	type nested struct{ V int }

	mapType := reflect.TypeOf(map[string]int{})
	if got := kindLabel(reflect.Map, mapType); got != "map[string]int" {
		t.Errorf("kindLabel(Map) = %q, want %q", got, "map[string]int")
	}

	nestType := reflect.TypeOf(nested{})
	if got := kindLabel(reflect.Struct, nestType); got != "struct nested" {
		t.Errorf("kindLabel(Struct) = %q, want %q", got, "struct nested")
	}
}

func TestStructDump_nilObj(t *testing.T) {
	if err := StructDump(nil, &bytes.Buffer{}); err == nil {
		t.Fatal("expected error for nil obj")
	}
}

func TestStructDump_notPointer(t *testing.T) {
	if err := StructDump(42, &bytes.Buffer{}); err == nil {
		t.Fatal("expected error for non-pointer value")
	}
}

func TestStructDump_nilPointer(t *testing.T) {
	var p *dumpSample
	if err := StructDump(p, &bytes.Buffer{}); err == nil {
		t.Fatal("expected error for nil pointer")
	}
}

func TestStructDump_notStruct(t *testing.T) {
	n := 42
	if err := StructDump(&n, &bytes.Buffer{}); err == nil {
		t.Fatal("expected error for pointer to non-struct")
	}
}

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) { return 0, errors.New("write error") }

func TestStructDump_writeError(t *testing.T) {
	s := &dumpSample{ID: 1}
	if err := StructDump(s, failingWriter{}); err == nil {
		t.Error("expected error from failing writer, got nil")
	}
}
