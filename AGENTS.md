# saferefl — AI Agent Guide

This file helps AI coding assistants (Claude, Copilot, Cursor, etc.) understand the project structure, key concepts, and how to work with the codebase effectively.

## What this library does

`saferefl` is a Go library for fast, type-safe struct reflection. It provides four layers of API:

- `Get[T]` / `Set[T]` / `MustGet[T]` / `MustSet[T]` — named field access with dot-path traversal, faster than `reflect.FieldByName`, zero allocs on warm cache
- `MakeAccessor[T](obj, "Field")` → `Accessor[T]` — bind once, read/write at ~0.5 ns per call (pointer arithmetic only)
- `EachField` · `CopyFields` · `ToMap` · `ToMapByTag` · `FromMap` · `GetByTag[T]` · `SetByTag[T]` · `MapForEach` · `KindOf` · `IsNil` — general-purpose struct/map API built on the `IterPlan` cache
- `MapLenFast[K,V](m)` · `UnsafeSliceAt[T](s, i)` — unsafe primitives, self-tested at init

## Key files

| File | Purpose |
|------|---------|
| `get.go` | `Get[T]` — eface fast path → ptrCache → resolvePath |
| `set.go` | `Set[T]` — same structure as get.go |
| `accessor.go` | `Accessor[T]`, `MakeAccessor`, `GetFrom`, `SetOn`, `UnsafePtrOf` |
| `eface.go` | `eface` struct, `efaceKind` — reads `reflect.Kind` from raw `*abi.Type` |
| `introspect.go` | `KindOf`, `IsNil` — fast kind/nil checks via raw `abi.Type` byte read |
| `iter.go` | `EachField`, `MapForEach`, `eachFieldRec`, `recurseEmbedded` (fallback path for structs with no exported fields) |
| `conv.go` | `ToMap`, `ToMapByTag`, `FromMap`, `flatToMap`, `flatToMapByTag` |
| `copy.go` | `CopyFields` — matches fields by name, uses `AssignableTo`/`ConvertibleTo` |
| `tag.go` | `GetByTag[T]`, `SetByTag[T]` — field access by struct tag value |
| `fieldany_unsafe.go` | `fieldAny(entry, objPtr)` — boxes a struct field into `any` using the precomputed `IfaceDirect` flag |
| `fieldany_safe.go` | `fieldAny` stub for `reflectx_strict` build tag — uses `reflect.NewAt` |
| `resolve.go` | `resolvePath` — dot-path traversal through nested/pointer structs |
| `primitives.go` | `AccelAvailable`, `EnableAccel`, `UnsafeSliceAt` |
| `maplen_swiss.go` | `MapLenFast` for Go 1.24+ Swiss Tables (reads `Map.used` at offset 0) |
| `maplen_legacy.go` | `MapLenFast` for Go < 1.24 hmap (reads `hmap.count` at offset 0) |
| `maplen_strict.go` | `MapLenFast` stub for `reflectx_strict` — returns `len(m)` |
| `errors.go` | `FieldNotFoundError`, `TypeMismatchError`, `ReadOnlyError`; sentinel vars `ErrFieldNotFound`, `ErrTypeMismatch`, `ErrReadOnly` |
| `field.go` | `Fields`, `FieldsOf[T]`, `FieldByName[T]` |
| `internal/typeinfo/abi_unsafe.go` | `abiTypeOf`, `isIfaceDirect`, `buildIterPlan`, `collectIter` — builds the flat `IterPlan` for each struct type |
| `internal/typeinfo/typeinfo.go` | `TypeDescriptor`, `FieldMeta`, `IterEntry` — struct metadata types |
| `internal/typeinfo/cache.go` | `TypeDescriptorOf`, `ptrCache` (sync.Map keyed by `uintptr(*abi.Type)`) |
| `internal/unsafelayout/` | `UnsafeFieldPtr`, `UnsafeSliceElemPtr`, `MapLen` — unsafe primitives with self-test |
| `debug/dump.go` | `StructDump` — annotated hex dump of any struct's memory (dev/test use only) |

## Architecture in one paragraph

`Get[T]` reads the raw `*abi.Type` from the `any` eface header, uses it as a key into `ptrCache` (a `sync.Map`), and on hit calls `getWithDesc` which does offset arithmetic. The first call per struct type falls to `getSlowPath` which calls `reflect.TypeOf` once and stores the descriptor. `Accessor[T]` pre-computes the chain at construction time so each `Get`/`Set` call is 1–3 pointer dereferences. The general-purpose API (`EachField`, `ToMap`, etc.) uses `TypeDescriptor.IterPlan` — a flat `[]IterEntry` built once by `buildIterPlan`/`collectIter` in `internal/typeinfo/abi_unsafe.go`. Each `IterEntry` carries the field offset, `*abi.Type` pointer, precomputed `IfaceDirect` flag, tag, and embed chain. At call time, `fieldAny` boxes the field into `any` using `IfaceDirect` to choose between direct storage (pointer types) and pointer-to-copy indirection (value types). Structs with no exported fields fall back to a reflect-based recursive path (`eachFieldRec` / `recurseEmbedded`). `MapLenFast` reads the first word of the runtime map header directly — verified by a self-test in `internal/unsafelayout` at `init()`.

## Build tags

| Tag | Effect |
|-----|--------|
| *(none)* | Normal build — unsafe accelerator compiled and self-tested at init |
| `reflectx_strict` | All `internal/unsafelayout` functions become no-ops; `MapLenFast` uses builtin `len`; `fieldAny` uses `reflect.NewAt` |

## How to run tests

```bash
go test ./...                          # default
go test -tags reflectx_strict ./...    # strict mode
go test -race ./...                    # race detector
SAFEREFL_STRICT=1 go test ./...        # panic (not log) on self-test failure
```

## How to run benchmarks

```bash
cd benchmarks && go test -bench=. -benchmem -count=3
./scripts/bench.sh          # regenerates benchmarks/RESULTS.md and injects into README
```

## Unsafe assumptions (self-tested at init)

1. **`abiTypeOf`**: `reflect.Type` is a non-empty interface whose data word holds `*abi.Type` directly — pointer types are stored directly in the interface data word (stable since Go 1).
2. **`isIfaceDirect`**: the directiface flag (bit 5, `0x20`) lives in `Kind_` at byte `2*sizeof(uintptr)+7` on Go 1.22–1.25, and in `TFlag` at byte `2*sizeof(uintptr)+4` on Go 1.26+. OR-ing both bytes makes the check version-agnostic. Offset is computed dynamically (not hardcoded 23) so it is also correct on 32-bit platforms.
3. **`MapLen` (hmap)**: `hmap.count` is the first field (stable since Go 1).
4. **`MapLen` (Swiss Tables, Go 1.24+)**: `Map.used` is the first field (verified per Go version).

If any assumption fails, `AccelAvailable()` returns false and `init()` logs a warning.
Call `EnableAccel()` once at program startup to surface the error explicitly.
Set `SAFEREFL_STRICT=1` to panic instead of log.

## Do not do

- Do not add `go:linkname` to private runtime symbols.
- Do not break the `ptrCache` keying logic — the key is `uintptr` of `*abi.Type` of the *pointer* type, not the elem type.
- Do not add per-call allocations to `Get[T]`/`Set[T]` or `EachField`/`ToMap` on the warm-cache path.
- Do not skip the self-test when adding new `internal/unsafelayout` or `internal/typeinfo/abi_unsafe.go` functions.
- Do not remove `//go:noinline` from `setSlowPathDesc`/`getSlowPath` — it prevents `val`/result from heap-escaping.
- Do not use `unsafe.Pointer(uintptr_value)` for pointer arithmetic — always use `unsafe.Add(ptr, offset)` to keep the pattern GC-safe and avoid `go vet unsafeptr`.
- Do not mutate `IterEntry.IfaceDirect` after construction — it is computed once in `collectIter` and must not change.
- Do not call `buildIterPlan` outside of `buildDescriptor` — the plan is a cache artifact built exactly once per type.

## Related resources

- [ADR-01 discussion](https://github.com/lkmavi/saferefl/discussions/3) — why this library exists and design decisions
- [Benchmark results](benchmarks/RESULTS.md) — detailed per-scenario numbers
- [CONTRIBUTING.md](CONTRIBUTING.md) — commit style, scopes, PR checklist
