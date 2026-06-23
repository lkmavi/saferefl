# Layer 3 Formal Invariants

This document states the invariants that must hold for `internal/unsafelayout` to be correct.
Every invariant is enforced at runtime by `runSelfTest()` before any unsafe primitive is invoked.

---

## Invariant 1 — Offset provenance

> Every `offset` argument passed to `UnsafeFieldPtr` was produced by `reflect.StructField.Offset`
> during `TypeDescriptor` construction.

**Why it matters:** `unsafe.Pointer` arithmetic is only safe when the offset is within a live
object. Using an offset that was not reflect-verified could point outside the object's allocation.

**Enforcement:** `TypeDescriptorOf` records `reflect.StructField.Offset` for every field at
descriptor-build time. No other path produces offsets consumed by `UnsafeFieldPtr`.

---

## Invariant 2 — Self-test precedes use

> `runSelfTest()` completes successfully (returns `true`) before any unsafe primitive is called
> at a call site in the same process.

**Why it matters:** The Go runtime's internal struct layouts (especially `hmap` and Swiss Tables
map headers) are not part of the public API. They can change between Go versions. The self-test
verifies the assumed layouts against the reflect baseline at `package init` time.

**Enforcement:** `AccelAvailable()` returns `false` if the self-test failed. Callers are required
to check `AccelAvailable()` or call `EnableAccel()` at startup. A `false` result means the
unsafe layer is silently bypassed in favour of the reflect path.

---

## Invariant 3 — No unsafe call when accel unavailable

> No unsafe primitive (`UnsafeFieldPtr`, `UnsafeSliceElemPtr`, `MapLen`) is called if
> `AccelAvailable()` returns `false`.

**Why it matters:** If the self-test failed, layout assumptions have been violated. Calling
unsafe primitives in that state produces incorrect results without any error signal.

**Enforcement:** Build-time (`reflectx_strict` tag) or runtime (`AccelAvailable()` guard).
The `reflectx_strict` build tag compiles all primitives to explicit no-op stubs so the
compiler can verify the guarantee statically.

---

## Invariant 4 — Map pointer validity

> `MapLen(m)` is only called with a pointer obtained from `reflect.Value.Pointer()` on a
> non-nil, non-empty map value.

**Why it matters:** An empty map may have a `nil` or sentinel internal pointer. Reading the
count field of a nil pointer is undefined behaviour.

**Enforcement:** Call sites must check `len(m) != 0` (or equivalent) before calling `MapLen`.
`selfTestMap` and `TestMapLen_matches_builtin` skip the zero-length case explicitly.

---

## Invariant 5 — Slice data pointer validity

> `UnsafeSliceElemPtr(sliceData, index, elemSize)` is only called with:
> - `sliceData` = the `Data` field of a live slice header (not a stale copy)
> - `index` in `[0, len(slice))`
> - `elemSize` = `reflect.Type.Size()` for the element type, verified at registration

**Why it matters:** Out-of-bounds indexing or a stale pointer produces a dangling-pointer read.

---

## What the self-test does NOT verify

- Correctness of `UnsafeFieldPtr` for types registered after `package init` — new types are
  verified implicitly because their offsets come from `reflect.StructField.Offset` (Invariant 1).
- GOARCH-specific struct padding beyond what `reflect` reports — reflect already accounts for
  alignment, so this is covered transitively.
- Map key/value reads (`MapGet`, `MapSet`) — not yet implemented; the interface exists in the
  plan but only `MapLen` is shipped in v0.x.

---

## Breaking-change checklist (new Go release)

When a new Go version ships:

1. Run the full CI matrix including `go tip`.
2. If `selfTestMap()` fails → update `map_layout_swiss.go` or `map_layout_legacy.go` to match
   the new runtime layout and add a `go1.XX` build tag.
3. If `selfTestStruct()` or `selfTestSlice()` fails → file a critical issue; these layouts have
   been stable since Go 1 and a failure indicates a fundamental ABI break.
4. If the failure cannot be fixed safely → document the degradation: `AccelAvailable()` returns
   `false` and callers transparently fall back to the reflect path.
