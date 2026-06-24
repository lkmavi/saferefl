//go:build !reflectx_strict

package typeinfo

import (
	"log"
	"os"
	"reflect"
	"unsafe"
)

// abiKindOffset is the byte offset of abi.Type.Kind_ within the abi.Type struct.
// Identical to saferefl.kindOffset; duplicated here to avoid package-cycle.
//
//	abi.Type layout (struct fields and offsets are identical across all Go versions, 64-bit):
//	  Size_       uintptr  +0
//	  PtrBytes    uintptr  +8
//	  Hash        uint32   +16
//	  TFlag       uint8    +20  ← abiFlagOffset
//	  Align_      uint8    +21
//	  FieldAlign_ uint8    +22
//	  Kind_       uint8    +23  ← abiKindOffset
//
// The "directiface" flag (bit 5, 0x20) indicates values stored directly in the interface
// data word. Its location within abi.Type changed in Go 1.26:
//   - Go 1.22–1.25: bit 5 is set in Kind_[+23]  (KindDirectIface)
//   - Go 1.26+:     bit 5 is set in TFlag[+20]  (TFlagDirectIface)
//
// Verified empirically with Docker images for each version (1.22–1.26).
const abiKindOffset = 2*unsafe.Sizeof(uintptr(0)) + 7 // Kind_ offset (+23 on 64-bit)
const abiFlagOffset = 2*unsafe.Sizeof(uintptr(0)) + 4 // TFlag offset (+20 on 64-bit)

// directIfaceMask is bit 5 (0x20) — the "directiface" flag position in both TFlag and Kind_.
// OR-ing both bytes covers all supported Go versions (1.22–1.25 use Kind_, 1.26+ use TFlag).
const directIfaceMask = uint8(1 << 5)

func init() {
	_, strict := os.LookupEnv("SAFEREFL_STRICT")
	fail := func(msg string) {
		if strict {
			panic(msg)
		}
		log.Println(msg)
	}

	// Verify abiTypeOf: boxing (*int)(nil) as any produces an eface whose _typ
	// we can extract both via our iface cast and by boxing through reflect.
	var xptr *int
	iface := (any)(xptr)
	type efaceWords struct{ _typ, data unsafe.Pointer }
	e := (*efaceWords)(unsafe.Pointer(&iface))
	if got := abiTypeOf(reflect.TypeOf(xptr)); got != e._typ {
		fail("[saferefl/typeinfo] abiTypeOf self-test FAILED — reflect.Type iface layout changed")
	}

	// Verify isIfaceDirect: *int is a pointer type (stored directly in iface data word),
	// int is a value type (stored via pointer to heap copy).
	// The directiface flag is in Kind_[+23] on Go 1.22–1.25 and in TFlag[+20] on Go 1.26+;
	// isIfaceDirect ORs both bytes so it works across the supported range (Go 1.22+).
	ptrAbi := abiTypeOf(reflect.TypeOf(xptr))
	intAbi := abiTypeOf(reflect.TypeOf(0))
	if !isIfaceDirect(ptrAbi) {
		fail("[saferefl/typeinfo] isIfaceDirect self-test FAILED — *int should have directiface flag")
	}
	if isIfaceDirect(intAbi) {
		fail("[saferefl/typeinfo] isIfaceDirect self-test FAILED — int should not have directiface flag")
	}
}

// abiTypeOf extracts the raw *abi.Type pointer from a reflect.Type interface value.
//
// reflect.Type's only concrete implementation is *reflect.rtype, which is the same
// struct as internal/abi.Type. Because *rtype is a pointer type it is stored directly
// in the interface data word (directiface flag is set for pointer types).
// The layout of a non-empty interface is {itab *itab; data unsafe.Pointer}, so
// the second word is the *rtype value itself.
func abiTypeOf(t reflect.Type) unsafe.Pointer {
	type iface struct {
		_ unsafe.Pointer // itab
		p unsafe.Pointer // *abi.Type (pointer types are stored directly, not via indirection)
	}
	return (*iface)(unsafe.Pointer(&t)).p //nolint:gosec
}

// isIfaceDirect reports whether values of this type are stored directly in the
// interface data word (not via a pointer to a heap-allocated copy).
//
// The directiface flag (bit 5, 0x20) moved in Go 1.26:
//   - Go 1.22–1.25: set in Kind_[+23]
//   - Go 1.26+:     set in TFlag[+20]
//
// OR-ing both bytes makes the check version-agnostic across the supported range (1.22+).
func isIfaceDirect(abiType unsafe.Pointer) bool {
	base := uintptr(abiType)
	tflag := *(*uint8)(unsafe.Pointer(base + abiFlagOffset)) //nolint:gosec
	kind  := *(*uint8)(unsafe.Pointer(base + abiKindOffset)) //nolint:gosec
	return (tflag|kind)&directIfaceMask != 0
}

// buildIterPlan constructs the flat IterEntry slice for the given struct type.
// Called once per type inside buildDescriptor.
func buildIterPlan(t reflect.Type) []IterEntry {
	var entries []IterEntry
	collectIter(t, 0, nil, &entries)
	return entries
}

// collectIter recursively walks t's fields, accumulating baseOffset for value-embedded
// structs (same chain) and appending to chain for pointer-embedded structs (new chain entry,
// reset baseOffset to 0).
func collectIter(t reflect.Type, base uintptr, chain []uintptr, out *[]IterEntry) {
	for i := range t.NumField() {
		sf := t.Field(i)
		off := base + sf.Offset

		if sf.Anonymous {
			k := sf.Type.Kind()
			if k == reflect.Struct {
				// Value-embedded: flatten with accumulated offset, same chain.
				collectIter(sf.Type, off, chain, out)
				continue
			}
			if k == reflect.Pointer && sf.Type.Elem().Kind() == reflect.Struct {
				// Pointer-embedded: the pointer itself (at offset off) is the next chain step.
				// Fields inside the pointed-to struct are relative to that struct (base resets to 0).
				newChain := make([]uintptr, len(chain)+1)
				copy(newChain, chain)
				newChain[len(chain)] = off
				collectIter(sf.Type.Elem(), 0, newChain, out)
				continue
			}
		}

		if !sf.IsExported() {
			continue
		}

		abiType := abiTypeOf(sf.Type)

		// Share the chain slice for all entries at the same embedding depth.
		var ch []uintptr
		if len(chain) > 0 {
			ch = chain
		}
		*out = append(*out, IterEntry{
			Name:        sf.Name,
			Tag:         sf.Tag,
			AbiType:     abiType,
			IfaceDirect: isIfaceDirect(abiType),
			Offset:      off,
			EmbedChain:  ch,
		})
	}
}
