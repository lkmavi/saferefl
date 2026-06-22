# saferefl

[![CI](https://github.com/lkmavi/saferefl/actions/workflows/ci.yml/badge.svg)](https://github.com/lkmavi/saferefl/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/lkmavi/saferefl/branch/main/graph/badge.svg)](https://codecov.io/gh/lkmavi/saferefl)
[![Go Reference](https://pkg.go.dev/badge/github.com/lkmavi/saferefl.svg)](https://pkg.go.dev/github.com/lkmavi/saferefl)
[![Go 1.22+](https://img.shields.io/badge/Go-1.22+-blue)](https://go.dev/dl/)

Fast, safe reflection for Go — a generic-first alternative to [`reflect2`](https://github.com/modern-go/reflect2).

## Why

`reflect2` trades correctness for speed by reverse-engineering Go's internal runtime layout. This causes silent data corruption when Go internals change (e.g. the map rewrite to Swiss Tables in Go 1.24). `saferefl` gets comparable speed through a different route: generics, cached offsets, and a self-verifying unsafe layer that falls back gracefully instead of corrupting memory.

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
┌─────────────────────────────────────────────┐
│  Layer 1: Public API (generics-first)        │  Get[T], Set[T]
│  - type-safe by construction                 │
│  - zero-copy fast path for exact types       │
├─────────────────────────────────────────────┤
│  Layer 2: Cached reflective layer            │  TypeDescriptorOf(t)
│  - built once via stdlib reflect             │
│  - sync.Map cache keyed on reflect.Type      │
├─────────────────────────────────────────────┤
│  Layer 3: Optional unsafe accelerator        │  build tag: unsafe_accel
│  - self-test at init, auto-fallback on fail  │  (planned)
│  - two map backends: hmap (<1.24), swiss (≥1.24) │
└─────────────────────────────────────────────┘
```

## Performance

<!-- bench:start -->
# Benchmark Results

**Generated:** 2026-06-22 11:22 UTC  
**Platform:** darwin/arm64  
**CPU:** Apple M3 Max  

## Get_int

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Saferefl | 26.1 | 0 | 0 | — |
| Reflect | 1.75 | 0 | 0 | 14.9× faster |
| Accessor | 0.550 | 0 | 0 | 47.4× faster |
| AccessorFrom | 4.38 | 0 | 0 | 6.0× faster |
| Direct | 0.269 | 0 | 0 | 97.1× faster |

## Set_string

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Saferefl | 38.6 | 16 | 1 | — |
| Reflect | 2.72 | 0 | 0 | 14.2× faster |
| Accessor | 0.538 | 0 | 0 | 71.8× faster |
| AccessorOn | 4.86 | 0 | 0 | 8.0× faster |

## DI Resolve

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 0.474 | 0 | 0 | — |
| Saferefl | 104 | 24 | 3 | 219.5× slower |
| Reflect | 104 | 0 | 0 | 220.2× slower |

## JSON Decode

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibJSON | 1495 | 392 | 10 | — |
| Saferefl | 400 | 120 | 10 | 3.7× faster |
| Reflect | 105 | 0 | 0 | 14.2× faster |
| Reflect2 | 33.6 | 0 | 0 | 44.5× faster |

## ORM Scan

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 5.48 | 0 | 0 | — |
| Saferefl | 397 | 112 | 10 | 72.4× slower |
| Reflect | 448 | 0 | 0 | 81.7× slower |
| Copier | 3254 | 640 | 28 | 593.6× slower |

## Struct Copy

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 2.60 | 0 | 0 | — |
| Saferefl | 357 | 56 | 5 | 137.1× slower |
| Reflect | 58.0 | 0 | 0 | 22.3× slower |
| Copier | 1345 | 432 | 17 | 516.5× slower |

## JSON Like

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibReflect | 196 | 64 | 1 | — |
| Reflect2 | 32.0 | 64 | 1 | 6.1× faster |
| CachedOffset | 80.6 | 64 | 1 | 2.4× faster |
| Generics | 26.4 | 64 | 1 | 7.4× faster |

## Read Int64

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibReflect | 26.1 | 0 | 0 | — |
| Reflect2 | 0.844 | 0 | 0 | 30.9× faster |
| CachedOffset | 5.73 | 0 | 0 | 4.6× faster |
| Generics | 0.283 | 0 | 0 | 92.2× faster |

## Write Int64

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibReflect | 27.3 | 0 | 0 | — |
| Reflect2 | 2.99 | 0 | 0 | 9.1× faster |
| CachedOffset | 10.8 | 0 | 0 | 2.5× faster |
| Generics | 0.269 | 0 | 0 | 101.5× faster |

## Read String

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibReflect | 27.2 | 0 | 0 | — |
| Reflect2 | 0.893 | 0 | 0 | 30.5× faster |
| CachedOffset | 5.23 | 0 | 0 | 5.2× faster |
| Generics | 0.553 | 0 | 0 | 49.2× faster |

<!-- bench:end -->

> Run `./scripts/bench.sh` to regenerate. Full results: [benchmarks/RESULTS.md](benchmarks/RESULTS.md).

Numbers above are for the spike benchmarks (Layer 2 cached-offset path). Layer 1 `Get[T]`/`Set[T]` benchmarks are in `benchmarks/layer1_bench_test.go`.

## Status

| Layer | Status | Description |
|---|---|---|
| Layer 2 — TypeInfo Cache | ✅ Done | `internal/typeinfo`: struct metadata, offset cache, field access via `reflect.NewAt` |
| Layer 1 — Generics API | ✅ Done | `Get[T]`, `Set[T]`, `MustGet[T]`, `MustSet[T]`, dot-path, `FieldByName[T]`, `Fields`, `FieldsOf[T]` |
| Layer 3 — Unsafe Accelerator | 🔲 Planned | opt-in `unsafe_accel` build tag, self-test at init, hmap/Swiss Tables backends |

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
