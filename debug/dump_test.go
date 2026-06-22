package debug

import (
	"bytes"
	"strings"
	"testing"
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
	// Only fields that start first in their 16-byte row get an annotation.
	// ID (offset 0) and Score (offset 24) are first in their rows; Name (offset 8)
	// shares row 0 with ID so it is not annotated — this is by design.
	for _, want := range []string{"dumpSample", "ID", "Score", "Active"} {
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
	for _, want := range []string{"string", "[]int", "*int", "interface"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing kind label %q:\n%s", want, out)
		}
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
