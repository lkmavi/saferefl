# saferefl

[![CI](https://github.com/lkmavi/saferefl/actions/workflows/ci.yml/badge.svg)](https://github.com/lkmavi/saferefl/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/lkmavi/saferefl/branch/main/graph/badge.svg)](https://codecov.io/gh/lkmavi/saferefl)
[![Go Reference](https://pkg.go.dev/badge/github.com/lkmavi/saferefl.svg)](https://pkg.go.dev/github.com/lkmavi/saferefl)
[![Go 1.22+](https://img.shields.io/badge/Go-1.22+-blue)](https://go.dev/dl/)

Fast, safe reflection for Go — a generic-first alternative to [`reflect2`](https://github.com/modern-go/reflect2).

## Why

`reflect2` trades correctness for speed by reverse-engineering Go's internal runtime layout. This causes silent data corruption when Go internals change (e.g. the map rewrite to Swiss Tables in Go 1.24). `saferefl` gets comparable speed through a different route: generics, cached offsets, and a self-verifying unsafe layer that falls back gracefully instead of corrupting memory.

The result: `Accessor[T]` lands within **1.2–2.3× of hand-written code** for ORM scan, DI injection, and struct copy — with zero allocations. The low-level generics path beats `reflect2` on raw field reads.

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

`Get[T]`/`Set[T]` resolve the field path on every call (27–39 ns). When you need to access the same field in a tight loop — ORM row scanning, DI injection, struct copying — pre-bind the path once with `Accessor[T]` and pay only ~0.55 ns per access:

```go
// Build once (e.g. at program startup or statement-prepare time)
ageAcc, _ := saferefl.MakeAccessor[int](u, "Age")

// Use many times — 0 allocations, pointer arithmetic only
ptr := saferefl.UnsafePtrOf(u)
age := ageAcc.Get(ptr)   // 0.55 ns, 0 allocs
ageAcc.Set(ptr, 31)      // 0.55 ns, 0 allocs

// Convenience form when you have an interface value, not a raw pointer
age, _ = ageAcc.GetFrom(u)    // 4.6 ns — eface extraction + field read
_ = ageAcc.SetOn(u, 31)       // 5.2 ns
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

**Generated:** 2026-06-22 12:17 UTC  
**Platform:** darwin/arm64  
**CPU:** Apple M3 Max  

## Get_int

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Saferefl | 27.5 | 0 | 0 | — |
| Reflect | 1.73 | 0 | 0 | 15.9× faster |
| Accessor | 0.546 | 0 | 0 | 50.4× faster |
| AccessorFrom | 4.63 | 0 | 0 | 5.9× faster |
| Direct | 0.272 | 0 | 0 | 100.9× faster |

## Set_string

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Saferefl | 38.7 | 16 | 1 | — |
| Reflect | 2.73 | 0 | 0 | 14.1× faster |
| Accessor | 0.543 | 0 | 0 | 71.2× faster |
| AccessorOn | 5.17 | 0 | 0 | 7.5× faster |

## Field Read

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Reflect | 1.73 | 0 | 0 | — |
| ReflectNewAt | 4.90 | 0 | 0 | 2.8× slower |
| Get | 27.0 | 0 | 0 | 15.6× slower |
| Accessor | 0.545 | 0 | 0 | 3.2× faster |
| Native | 0.273 | 0 | 0 | 6.4× faster |

## Slice At

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Layer3 | 0.549 | 0 | 0 | — |
| Direct | 0.548 | 0 | 0 | 1.0× faster |
| Reflect | 1.92 | 0 | 0 | 3.5× slower |

## Map Len

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Layer3 | 0.278 | 0 | 0 | — |
| Builtin | 0.274 | 0 | 0 | 1.0× faster |
| Reflect | 2.17 | 0 | 0 | 7.8× slower |

## DI Resolve

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 0.483 | 0 | 0 | — |
| Saferefl | 105 | 24 | 3 | 217.6× slower |
| Accessor | 1.09 | 0 | 0 | 2.3× slower |
| Reflect | 106 | 0 | 0 | 218.4× slower |

## JSON Decode

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibJSON | 1492 | 392 | 10 | — |
| Saferefl | 389 | 120 | 10 | 3.8× faster |
| Reflect | 105 | 0 | 0 | 14.2× faster |
| Reflect2 | 33.6 | 0 | 0 | 44.4× faster |

## ORM Scan

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 5.54 | 0 | 0 | — |
| Saferefl | 383 | 112 | 10 | 69.2× slower |
| Accessor | 6.70 | 0 | 0 | 1.2× slower |
| Reflect | 440 | 0 | 0 | 79.5× slower |
| Copier | 3241 | 640 | 28 | 585.4× slower |

## Struct Copy

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 2.57 | 0 | 0 | — |
| Saferefl | 350 | 56 | 5 | 136.0× slower |
| Reflect | 57.2 | 0 | 0 | 22.2× slower |
| Accessor | 3.31 | 0 | 0 | 1.3× slower |
| Copier | 1340 | 432 | 17 | 521.1× slower |

## JSON Like

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibReflect | 192 | 64 | 1 | — |
| Reflect2 | 30.3 | 64 | 1 | 6.3× faster |
| CachedOffset | 80.2 | 64 | 1 | 2.4× faster |
| Generics | 26.0 | 64 | 1 | 7.4× faster |

## Read Int64

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibReflect | 26.1 | 0 | 0 | — |
| Reflect2 | 0.848 | 0 | 0 | 30.7× faster |
| CachedOffset | 5.73 | 0 | 0 | 4.5× faster |
| Generics | 0.287 | 0 | 0 | 90.7× faster |

## Write Int64

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibReflect | 27.4 | 0 | 0 | — |
| Reflect2 | 2.99 | 0 | 0 | 9.2× faster |
| CachedOffset | 10.9 | 0 | 0 | 2.5× faster |
| Generics | 0.273 | 0 | 0 | 100.5× faster |

## Read String

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibReflect | 27.4 | 0 | 0 | — |
| Reflect2 | 0.889 | 0 | 0 | 30.8× faster |
| CachedOffset | 5.18 | 0 | 0 | 5.3× faster |
| Generics | 0.546 | 0 | 0 | 50.1× faster |

<!-- bench:end -->

> Run `./scripts/bench.sh` to regenerate. Full results: [benchmarks/RESULTS.md](benchmarks/RESULTS.md).

Key numbers to read: **Accessor** rows show Layer 3 hot-path performance; **Saferefl** rows show Layer 1 (`Get[T]`/`Set[T]`) with full path resolution per call; **Read/Write Int64/String** (spike) compare raw field access across all approaches including reflect2.

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

Floor: **Go 1.22**. CI matrix: 1.22, 1.23, 1.24, stable, tip.
