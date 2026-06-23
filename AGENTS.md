# saferefl — AI Agent Guide

This file helps AI coding assistants (Claude, Copilot, Cursor, etc.) understand the project structure, key concepts, and how to work with the codebase effectively.

## What this library does

`saferefl` is a Go library for fast, type-safe struct field access via reflection. It provides:

- `Get[T](obj, "Field")` / `Set[T](obj, "Field", val)` — named field access, faster than `reflect.FieldByName`, zero allocs on warm cache
- `MakeAccessor[T](obj, "Field")` → `Accessor[T]` — bind once, read/write at ~0.5 ns per call (pointer arithmetic only)
- `MapLenFast[K,V](m)` — map length without reflect (same speed as builtin `len`)
- `UnsafeSliceAt[T](s, i)` — unchecked slice element pointer

Use this instead of `reflect.FieldByName` in hot loops, ORM row scanners, DI containers, and struct-copy routines.

## Key files

| File | Purpose |
|------|---------|
| `get.go` | `Get[T]` — eface fast path → ptrCache → resolvePath |
| `set.go` | `Set[T]` — same structure as get.go |
| `accessor.go` | `Accessor[T]`, `MakeAccessor`, `GetFrom`, `SetOn`, `UnsafePtrOf` |
| `eface.go` | `eface` struct, `efaceKind` — reads `reflect.Kind` from raw `*abi.Type` |
| `resolve.go` | `resolvePath` — dot-path traversal through nested/pointer structs |
| `primitives.go` | `AccelAvailable`, `EnableAccel`, `UnsafeSliceAt` |
| `maplen_swiss.go` | `MapLenFast` for Go 1.24+ Swiss Tables (reads `Map.used` at offset 0) |
| `maplen_legacy.go` | `MapLenFast` for Go < 1.24 hmap (reads `hmap.count` at offset 0) |
| `maplen_strict.go` | `MapLenFast` stub when `reflectx_strict` tag is set — returns `len(m)` |
| `errors.go` | `FieldNotFoundError`, `TypeMismatchError`, `ReadOnlyError` |
| `field.go` | `Fields`, `FieldsOf[T]`, `FieldByName[T]` |
| `internal/typeinfo/` | `TypeDescriptor`, `ptrCache` — struct metadata built once via stdlib reflect |
| `internal/unsafelayout/` | `UnsafeFieldPtr`, `UnsafeSliceElemPtr`, `MapLen` — unsafe primitives with self-test |

## Architecture in one paragraph

`Get[T]` reads the raw `*abi.Type` from the `any` eface header, uses it as a key into `ptrCache` (a `sync.Map`), and on hit calls `setWithDesc`/`getWithDesc` which does offset arithmetic. First call per struct type falls to `setSlowPathDesc` which calls `reflect.TypeOf` once and stores the descriptor. `Accessor[T]` pre-computes the chain at construction time so each `Get`/`Set` call is 1-3 pointer dereferences. `MapLenFast` reads the first word of the runtime map header directly — verified by a self-test in `internal/unsafelayout` at `init()`.

## Build tags

| Tag | Effect |
|-----|--------|
| *(none)* | Normal build — unsafe accelerator compiled and self-tested at init |
| `reflectx_strict` | All `internal/unsafelayout` functions become no-ops; `MapLenFast` uses builtin `len` |

## How to run tests

```bash
go test ./...                          # default
go test -tags reflectx_strict ./...    # strict mode
go test -race ./...                    # race detector
```

## How to run benchmarks

```bash
cd benchmarks && go test -bench=. -benchmem -count=3
./scripts/bench.sh          # regenerates benchmarks/RESULTS.md and injects into README
```

## Unsafe assumptions (self-tested at init)

1. `efaceKind`: `*abi.Type` byte at offset 23 holds `Kind_` (stable since Go 1.18)
2. `MapLen` (hmap): `hmap.count` is the first field (stable since Go 1)
3. `MapLen` (Swiss Tables, Go 1.24+): `Map.used` is the first field (verified per Go version)

If any assumption fails, `AccelAvailable()` returns false and `init()` logs a warning.  
Call `EnableAccel()` once at program startup to surface the error explicitly.

## Do not do

- Do not add `go:linkname` to private runtime symbols
- Do not break the `ptrCache` keying logic (key is `uintptr` of `*abi.Type` of the *pointer* type, not the elem type)
- Do not add per-call allocations to `Get[T]`/`Set[T]` on the warm-cache path
- Do not skip the self-test when adding new `internal/unsafelayout` functions
- Do not remove `//go:noinline` from `setSlowPathDesc`/`getSlowPath` — it prevents `val`/result from heap-escaping

## Related resources

- [ADR-01 discussion](https://github.com/lkmavi/saferefl/discussions/3) — why this library exists and design decisions
- [Benchmark results](benchmarks/RESULTS.md) — detailed per-scenario numbers
- [CONTRIBUTING.md](CONTRIBUTING.md) — commit style, scopes, PR checklist
