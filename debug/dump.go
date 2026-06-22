// Package debug provides developer utilities built on saferefl's TypeDescriptor.
// Import it only in debug/test builds — it is not needed for production use of saferefl.
package debug

import (
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"
	"unicode"
	"unsafe"

	"github.com/lkmavi/saferefl/internal/typeinfo"
)

const bytesPerRow = 16

// StructDump writes an annotated hex dump of obj's memory to w.
// Each row shows the byte offset, hex values, ASCII representation, and the
// name+type of the struct field that starts at that offset.
//
// obj must be a non-nil pointer to a struct.
//
// Example output:
//
//	User (32 bytes)
//	+0000  01 00 00 00 00 00 00 00  │ ........  ID      int64
//	+0008  05 00 00 00 00 00 00 00  │ ........  Name    string (ptr+len)
//	+0018  9a 99 19 40              │ ...@      Score   float64
func StructDump(obj any, w io.Writer) error {
	if obj == nil {
		return fmt.Errorf("debug: obj must not be nil")
	}
	rv := reflect.ValueOf(obj)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return fmt.Errorf("debug: obj must be a non-nil pointer to a struct")
	}
	rt := rv.Type().Elem()
	if rt.Kind() != reflect.Struct {
		return fmt.Errorf("debug: obj must point to a struct, got pointer to %v", rt.Kind())
	}

	desc := typeinfo.TypeDescriptorOf(rt)
	size := int(rt.Size())
	rawMem := unsafe.Slice((*byte)(unsafe.Pointer(rv.Pointer())), size)

	// Build offset → field name map for annotation.
	fieldAt := buildFieldAt(desc)

	fmt.Fprintf(w, "%s (%d bytes)\n", rt.Name(), size)
	fmt.Fprintf(w, "%-6s  %-47s  │  %-16s  %s\n", "offset", "hex", "ascii", "field")
	fmt.Fprintf(w, "%s\n", strings.Repeat("─", 85))

	for row := 0; row*bytesPerRow < size; row++ {
		start := row * bytesPerRow
		end := start + bytesPerRow
		if end > size {
			end = size
		}
		chunk := rawMem[start:end]

		// Hex column.
		var hexBuf strings.Builder
		for i, b := range chunk {
			if i > 0 {
				hexBuf.WriteByte(' ')
			}
			fmt.Fprintf(&hexBuf, "%02x", b)
		}

		// ASCII column (printable only, '.' for control/non-ASCII).
		var asciiBuf strings.Builder
		for _, b := range chunk {
			r := rune(b)
			if b < 0x20 || b > 0x7e || !unicode.IsPrint(r) {
				asciiBuf.WriteByte('.')
			} else {
				asciiBuf.WriteByte(b)
			}
		}

		// Field annotation: first field that starts in this row.
		annotation := ""
		for off := start; off < end; off++ {
			if fa, ok := fieldAt[uintptr(off)]; ok {
				annotation = fa
				break
			}
		}

		fmt.Fprintf(w, "+%04x  %-47s  │  %-16s  %s\n",
			start, hexBuf.String(), asciiBuf.String(), annotation)
	}
	return nil
}

// fieldAnnotation holds the display string for a field at a given offset.
type fieldAnnotation struct {
	offset uintptr
	label  string
}

func buildFieldAt(desc *typeinfo.TypeDescriptor) map[uintptr]string {
	// Collect only direct (non-promoted) fields to avoid duplicate annotations.
	type entry struct {
		offset uintptr
		label  string
	}
	entries := make([]entry, 0, len(desc.Fields))
	for i := range desc.Fields {
		fm := &desc.Fields[i]
		entries = append(entries, entry{
			offset: fm.Offset,
			label:  fmt.Sprintf("%-12s %s", fm.Name, kindLabel(fm.Kind, fm.Type)),
		})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].offset < entries[j].offset })

	m := make(map[uintptr]string, len(entries))
	for _, e := range entries {
		m[e.offset] = e.label
	}
	return m
}

func kindLabel(k reflect.Kind, t reflect.Type) string {
	switch k {
	case reflect.String:
		return "string"
	case reflect.Slice:
		return fmt.Sprintf("[]%s", t.Elem().Name())
	case reflect.Pointer:
		return fmt.Sprintf("*%s", t.Elem().Name())
	case reflect.Map:
		return fmt.Sprintf("map[%s]%s", t.Key().Name(), t.Elem().Name())
	case reflect.Struct:
		return fmt.Sprintf("struct %s", t.Name())
	case reflect.Interface:
		return "interface"
	default:
		return k.String()
	}
}
