# saferefl

[![CI](https://github.com/lkmavi/saferefl/actions/workflows/ci.yml/badge.svg)](https://github.com/lkmavi/saferefl/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/lkmavi/saferefl/branch/main/graph/badge.svg)](https://codecov.io/gh/lkmavi/saferefl)
[![Go Reference](https://pkg.go.dev/badge/github.com/lkmavi/saferefl.svg)](https://pkg.go.dev/github.com/lkmavi/saferefl)
[![Go 1.22+](https://img.shields.io/badge/Go-1.22+-blue)](https://go.dev/dl/)

Fast, safe reflection for Go — a generic-first alternative to [`reflect2`](https://github.com/modern-go/reflect2).

## Why

`reflect2` trades correctness for speed by reverse-engineering Go's internal runtime layout. This causes silent data corruption when Go internals change (e.g. the map rewrite to Swiss Tables in Go 1.24). `saferefl` gets comparable speed through a different route: generics, cached offsets, and a self-verifying unsafe layer that falls back gracefully instead of corrupting memory.

The result: `Accessor[T]` lands within **1.1–1.3× of hand-written code** for ORM scan and struct copy — with zero allocations. `Get[T]`/`Set[T]` also achieve zero allocations on the fast path. The low-level generics path beats `reflect2` on raw field reads and is 12× faster than reflect2 on the field-setting phase of JSON decode.

See [ADR-01](https://github.com/lkmavi/saferefl/discussions/3) for the full analysis and decision.

## Quick Start

```go
import "github.com/lkmavi/saferefl"

type User struct {
    Name string
    Age  int
    Score float64
}

u := &User{Name: "Alice", Age: 30, Score: 9.5}

// Read any field by name — type-safe, zero boxing on the fast path
name, err := saferefl.Get[string](u, "Name")   // "Alice", nil
age,  err := saferefl.Get[int](u, "Age")        // 30, nil

// Write
_ = saferefl.Set[string](u, "Name", "Bob")
_ = saferefl.Set[int](u, "Age", 31)

// Panic variants for statically-known valid paths
name = saferefl.MustGet[string](u, "Name")
saferefl.MustSet[float64](u, "Score", 10.0)
```

### Hot-path: Accessor[T]

`Get[T]`/`Set[T]` resolve the field path on every call (~27 ns). When you need to access the same field in a tight loop — ORM row scanning, DI injection, struct copying — pre-bind the path once with `Accessor[T]` and pay only ~0.55 ns per access:

```go
// Build once (e.g. at program startup or statement-prepare time)
ageAcc, _ := saferefl.MakeAccessor[int](u, "Age")

// Use many times — 0 allocations, pointer arithmetic only
ptr := saferefl.UnsafePtrOf(u)
age := ageAcc.Get(ptr)   // 0.55 ns, 0 allocs
ageAcc.Set(ptr, 31)      // 0.55 ns, 0 allocs

// Convenience form when you have an interface value, not a raw pointer
age, _ = ageAcc.GetFrom(u)    // 1.4 ns — eface extraction + field read
_ = ageAcc.SetOn(u, 31)       // 1.9 ns
```

### Dot-path traversal

Intermediate struct fields and pointer-to-struct fields are transparently traversed:

```go
type Address struct {
    City    string
    Country string
}

type Employee struct {
    User
    Office  Address
    Contact *Address
}

e := &Employee{
    User:    User{Name: "Carol"},
    Office:  Address{City: "Berlin"},
    Contact: &Address{City: "NYC"},
}

// Promoted field from embedded User
name, _ := saferefl.Get[string](e, "Name")          // "Carol"

// Value intermediate field
city, _ := saferefl.Get[string](e, "Office.City")    // "Berlin"

// Pointer intermediate field — nil pointer returns an error, never panics
city, _ =  saferefl.Get[string](e, "Contact.City")   // "NYC"
```

### Field inspection

```go
// Direct fields of a type (no instance needed)
fields, _ := saferefl.FieldsOf[User]()
for _, f := range fields {
    fmt.Println(f.Name, f.Type)
}

// From an instance (struct value or pointer)
fields, _ = saferefl.Fields(u)

// Single field lookup
sf, ok := saferefl.FieldByName[User]("Name")
```

## Error types

All errors are typed and work with `errors.As`:

| Type | When |
|---|---|
| `*FieldNotFoundError` | field path does not exist on the type |
| `*TypeMismatchError` | field type is not assignable to T |
| `*ReadOnlyError` | attempted to Set an unexported field |

```go
_, err := saferefl.Get[int](u, "Name")   // Name is string, not int

var tme *saferefl.TypeMismatchError
if errors.As(err, &tme) {
    fmt.Println(tme.FieldType, "vs", tme.WantType) // "string vs int"
}
```

## Architecture

```
┌──────────────────────────────────────────────────────────┐
│  Layer 1: Public API (generics-first)                     │
│  Get[T], Set[T], MustGet[T], MustSet[T], Accessor[T]     │
│  - type-safe by construction                              │
│  - dot-path traversal through nested / pointer structs    │
├──────────────────────────────────────────────────────────┤
│  Layer 2: Cached reflective layer       TypeDescriptorOf  │
│  - struct metadata built once via stdlib reflect          │
│  - sync.Map + atomic.Pointer[T] cache, zero-alloc reads  │
├──────────────────────────────────────────────────────────┤
│  Layer 3: Self-verified unsafe accelerator                │
│  Accessor[T].Get/.Set · UnsafeSliceAt[T] · MapLenFast    │
│  - self-test at init(); graceful fallback on mismatch     │
│  - two map backends: hmap (< Go 1.24), Swiss (≥ Go 1.24) │
│  - disable with build tag: reflectx_strict                │
└──────────────────────────────────────────────────────────┘
```

## Performance

<!-- bench:start -->
# Benchmark Results

**Generated:** 2026-06-22 18:03 UTC  
**Go:** 1.26  
**Platform:** darwin/arm64  
**CPU:** Apple M3 Max  

## Get_int

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| L1 | 26.1 | 0 | 0 | — |
| Reflect | 1.69 | 0 | 0 | 15.5× faster |
| L2 | 4.80 | 0 | 0 | 5.4× faster |
| L3 | 0.532 | 0 | 0 | 49.1× faster |
| L3From | 1.42 | 0 | 0 | 18.4× faster |
| Direct | 0.267 | 0 | 0 | 97.9× faster |

## Set_string

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| L1 | 26.5 | 0 | 0 | — |
| Reflect | 2.67 | 0 | 0 | 9.9× faster |
| L2 | 5.59 | 0 | 0 | 4.7× faster |
| L3 | 0.534 | 0 | 0 | 49.7× faster |
| L3On | 1.86 | 0 | 0 | 14.2× faster |
| Direct | 0.265 | 0 | 0 | 100.1× faster |

## Field Read

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Reflect | 1.69 | 0 | 0 | — |
| L1 | 26.1 | 0 | 0 | 15.5× slower |
| L2 | 5.06 | 0 | 0 | 3.0× slower |
| L3 | 0.532 | 0 | 0 | 3.2× faster |
| Native | 0.266 | 0 | 0 | 6.3× faster |

## Slice At

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| L3 | 0.535 | 0 | 0 | — |
| Direct | 0.536 | 0 | 0 | 1.0× slower |
| Reflect | 1.87 | 0 | 0 | 3.5× slower |

## Map Len

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| L3 | 0.273 | 0 | 0 | — |
| Builtin | 0.271 | 0 | 0 | 1.0× faster |
| Reflect | 2.13 | 0 | 0 | 7.8× slower |

## DI Resolve

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 0.474 | 0 | 0 | — |
| L1 | 78.6 | 0 | 0 | 165.8× slower |
| L3 | 3.22 | 0 | 0 | 6.8× slower |
| Reflect2 | 8.53 | 0 | 0 | 18.0× slower |
| Reflect | 103 | 0 | 0 | 217.6× slower |

## JSON Decode

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibJSON | 1434 | 392 | 10 | — |
| L1 | 270 | 0 | 0 | 5.3× faster |
| Reflect | 687 | 0 | 0 | 2.1× faster |
| Reflect2 | 36.3 | 0 | 0 | 39.5× faster |
| L3 | 2.93 | 0 | 0 | 489.6× faster |

## ORM Scan

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 5.28 | 0 | 0 | — |
| L1 | 271 | 0 | 0 | 51.3× slower |
| L3 | 5.33 | 0 | 0 | 1.0× slower |
| Reflect2 | 37.0 | 0 | 0 | 7.0× slower |
| Reflect | 431 | 0 | 0 | 81.6× slower |
| Copier | 3109 | 640 | 28 | 588.3× slower |

## Struct Copy

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 2.49 | 0 | 0 | — |
| L3 | 3.92 | 0 | 0 | 1.6× slower |
| Reflect2 | 21.0 | 0 | 0 | 8.4× slower |
| L2 | 84.2 | 0 | 0 | 33.8× slower |
| Reflect | 305 | 0 | 0 | 122.2× slower |
| L1 | 286 | 0 | 0 | 114.7× slower |
| Copier | 1322 | 432 | 17 | 530.3× slower |

<!-- bench:end -->

> Run `./scripts/bench.sh` to regenerate. Full results: [benchmarks/RESULTS.md](benchmarks/RESULTS.md).

Key numbers to read: **L3** = `saferefl.Accessor[T]` hot-path (pre-bound once, pointer arithmetic only); **L3From/L3On** = Accessor with interface→pointer conversion per call; **L1** = `saferefl.Get[T]`/`Set[T]` (path resolution per call, all 0 allocs); **L2** = cached-offset + `reflect.NewAt` (intermediate layer, used internally by the library). In all realistic benchmarks **Reflect** uses `FieldByName` per call — the common usage baseline. **Reflect2** uses pre-compiled field descriptors from [reflect2](https://github.com/modern-go/reflect2), representing a well-optimised codec that caches metadata at startup.

### Results by Go version

Per-version results are generated automatically by the [Cross-version Benchmarks](.github/workflows/bench-matrix.yml) workflow (runs weekly, or trigger manually).

| Go version | Results |
|---|---|
| 1.22 | [benchmarks/results/go1.22.md](benchmarks/results/go1.22.md) |
| 1.23 | [benchmarks/results/go1.23.md](benchmarks/results/go1.23.md) |
| 1.24 | [benchmarks/results/go1.24.md](benchmarks/results/go1.24.md) |
| 1.25 | [benchmarks/results/go1.25.md](benchmarks/results/go1.25.md) |
| 1.26 | [benchmarks/results/go1.26.md](benchmarks/results/go1.26.md) |

To generate locally (requires Docker):

```
make bench-docker        # all versions
make bench-docker-1.24   # single version
```

## Status

| Layer | Status | Description |
|---|---|---|
| Layer 2 — TypeInfo Cache | ✅ Done | `internal/typeinfo`: struct metadata, `sync.Map` + `atomic.Pointer` cache, direct pointer arithmetic |
| Layer 1 — Generics API | ✅ Done | `Get[T]`, `Set[T]`, `MustGet[T]`, `MustSet[T]`, `Accessor[T]`, dot-path, `FieldByName[T]`, `Fields`, `FieldsOf[T]` |
| Layer 3 — Unsafe Accelerator | ✅ Done | `internal/unsafelayout`: self-test at init, hmap/Swiss Tables backends, `UnsafeSliceAt[T]`, `MapLenFast` |

## Examples

Runnable examples are in [`examples/`](examples/):

- [`examples/basic/`](examples/basic/) — Get/Set primitive fields
- [`examples/dotpath/`](examples/dotpath/) — dot-path traversal through nested structs
- [`examples/fields/`](examples/fields/) — field inspection without an instance

```
go run ./examples/basic/
go run ./examples/dotpath/
go run ./examples/fields/
```

## Go version support

Floor: **Go 1.22**. CI matrix: 1.22, 1.23, 1.24, 1.25, 1.26, stable, tip.
