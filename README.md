# saferefl

[![CI](https://github.com/lkmavi/saferefl/actions/workflows/ci.yml/badge.svg)](https://github.com/lkmavi/saferefl/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/lkmavi/saferefl/branch/main/graph/badge.svg)](https://codecov.io/gh/lkmavi/saferefl)
[![Go Reference](https://pkg.go.dev/badge/github.com/lkmavi/saferefl.svg)](https://pkg.go.dev/github.com/lkmavi/saferefl)
[![Go 1.22+](https://img.shields.io/badge/Go-1.22+-blue)](https://go.dev/dl/)

Fast, safe reflection for Go — a generic-first alternative to [`reflect2`](https://github.com/modern-go/reflect2).

## Why

`reflect2` trades correctness for speed by reverse-engineering Go's internal runtime layout. This causes silent data corruption when Go internals change (e.g. the map rewrite to Swiss Tables in Go 1.24). `saferefl` gets comparable speed through a different route: generics, cached offsets, and a self-verifying unsafe layer that falls back gracefully instead of corrupting memory.

The result: `Accessor[T]` lands within **1.1–1.3× of hand-written code** for ORM scan and struct copy — with zero allocations. `Get[T]`/`Set[T]` are **faster than `reflect.FieldByName`** on the warm-cache path (~20 ns vs ~26 ns) with zero allocations. The low-level generics path beats `reflect2` on raw field reads and is 12× faster than reflect2 on the field-setting phase of JSON decode.

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

`Get[T]`/`Set[T]` resolve the field path on every call (~21 ns, faster than `reflect.FieldByName`). When you need to access the same field in a tight loop — ORM row scanning, DI injection, struct copying — pre-bind the path once with `Accessor[T]` and pay only ~0.55 ns per access:

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

## General-purpose API

Higher-level functions for struct iteration, mapping, and tag-based access. All use the prebuilt `IterPlan` cache — zero allocs on the hot path, embedded structs flattened automatically, nil pointer embeds skipped.

### Field iteration

```go
type Article struct {
    Title  string
    Author string
    Views  int
}

a := &Article{Title: "Go internals", Author: "Alice", Views: 1200}

// Iterate all exported fields in declaration order; return false to stop early.
_ = saferefl.EachField(a, func(name string, val any) bool {
    fmt.Printf("%s = %v\n", name, val)
    return true
})
// Title = Go internals
// Author = Alice
// Views = 1200
```

### Struct copying

```go
type Src struct { Name string; Score float64 }
type Dst struct { Name string; Score float64; Extra int }

src := &Src{Name: "Alice", Score: 9.5}
dst := &Dst{}

_ = saferefl.CopyFields(src, dst)
// dst.Name == "Alice", dst.Score == 9.5, dst.Extra == 0 (no matching field)
```

`CopyFields` uses `AssignableTo` first, then `ConvertibleTo` for compatible-but-different types (e.g. `float32` → `float64`). Unmatched and incompatible fields are silently skipped.

### Struct ↔ map

```go
// Struct → map (field names as keys)
m, _ := saferefl.ToMap(a)       // map["Title":"Go internals" "Author":"Alice" "Views":1200]

// Struct → map (tag values as keys; skips "-" and empty tag names; strips omitempty)
type Tagged struct {
    Name  string `json:"name"`
    Score float64 `json:"score,omitempty"`
    Skip  string  `json:"-"`
}
t := &Tagged{Name: "Bob", Score: 8.0}
m, _ = saferefl.ToMapByTag(t, "json")   // map["name":"Bob" "score":8.0]

// Map → struct (tries AssignableTo, then ConvertibleTo; skips unknown keys and nil values)
dst2 := &Tagged{}
_ = saferefl.FromMap(map[string]any{"Name": "Carol", "Score": float64(7.5)}, dst2)
```

### Tag-based field access

```go
type Record struct {
    ID   int    `db:"id"`
    Name string `db:"name"`
}

r := &Record{ID: 1, Name: "Alice"}

id, err := saferefl.GetByTag[int](r, "db", "id")       // 1, nil
_ = saferefl.SetByTag[string](r, "db", "name", "Bob")  // r.Name == "Bob"
```

Sentinel errors work with `errors.Is`:

```go
_, err = saferefl.GetByTag[int](r, "db", "missing")
errors.Is(err, saferefl.ErrFieldNotFound) // true
```

### Map and kind utilities

```go
// Typed map iteration — identical speed to plain range, early-stop on false
saferefl.MapForEach(map[string]int{"a": 1, "b": 2}, func(k string, v int) bool {
    fmt.Println(k, v)
    return true
})

// Fast kind/nil checks via raw abi.Type read (~0.28 ns, 0 allocs)
saferefl.KindOf("hello")  // reflect.String
saferefl.KindOf(42)       // reflect.Int
saferefl.IsNil((*int)(nil)) // true
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

All errors are typed and work with `errors.Is` / `errors.As`:

| Sentinel | Type | When |
|---|---|---|
| `ErrFieldNotFound` | `*FieldNotFoundError` | field path or tag value does not exist |
| `ErrTypeMismatch` | `*TypeMismatchError` | field type is not assignable/convertible to T |
| `ErrReadOnly` | `*ReadOnlyError` | attempted to write an unexported field |

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
│  General-purpose API                                      │
│  EachField · CopyFields · ToMap · ToMapByTag · FromMap    │
│  GetByTag[T] · SetByTag[T] · MapForEach · KindOf · IsNil │
│  - embedded structs flattened; nil pointer embeds skipped │
│  - zero allocs on warm cache (IterPlan fast path)         │
├──────────────────────────────────────────────────────────┤
│  Generic API                                              │
│  Get[T], Set[T], MustGet[T], MustSet[T]                  │
│  Fields, FieldsOf[T], FieldByName[T]                      │
│  - type-safe by construction, zero allocs on warm cache  │
│  - dot-path traversal through nested / pointer structs    │
├──────────────────────────────────────────────────────────┤
│  TypeInfo Cache                       (internal)          │
│  TypeDescriptorOf · IterPlan · sync.Map · atomic.Pointer  │
│  - struct metadata + flat IterPlan built once via reflect │
│  - zero-alloc reads after the first call per type         │
├──────────────────────────────────────────────────────────┤
│  Accessor API + Unsafe Primitives                         │
│  Accessor[T].Get/.Set · UnsafeSliceAt[T] · MapLenFast    │
│  - bind field path once, pay pointer arithmetic per call  │
│  - self-test at init(); graceful fallback on mismatch     │
│  - two map backends: hmap (< Go 1.24), Swiss (≥ Go 1.24) │
│  - disable with build tag: reflectx_strict                │
└──────────────────────────────────────────────────────────┘
```

## Performance

<!-- bench:start -->
# Benchmark Results

**Generated:** 2026-06-23 08:26 UTC  
**Go:** 1.26  
**Platform:** darwin/arm64  
**CPU:** Apple M3 Max  

## Field Read

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Reflect | 26.8 | 0 | 0 | — |
| ReflectFast | 1.75 | 0 | 0 | 15.3× faster |
| SafeRefl | 21.2 | 0 | 0 | 1.3× faster |
| Offset | 5.00 | 0 | 0 | 5.4× faster |
| Accessor | 0.555 | 0 | 0 | 48.2× faster |
| Native | 0.278 | 0 | 0 | 96.3× faster |

## Slice At

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.564 | 0 | 0 | — |
| Direct | 0.548 | 0 | 0 | 1.0× faster |
| Reflect | 1.91 | 0 | 0 | 3.4× slower |

## Map Len

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.277 | 0 | 0 | — |
| Builtin | 0.274 | 0 | 0 | 1.0× faster |
| Reflect | 2.20 | 0 | 0 | 7.9× slower |

## Get_int

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 21.2 | 0 | 0 | — |
| Reflect | 26.1 | 0 | 0 | 1.2× slower |
| ReflectFast | 1.75 | 0 | 0 | 12.1× faster |
| Offset | 5.03 | 0 | 0 | 4.2× faster |
| Accessor | 0.546 | 0 | 0 | 38.8× faster |
| AccFrom | 1.45 | 0 | 0 | 14.6× faster |
| Direct | 0.283 | 0 | 0 | 75.0× faster |

## Set_string

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 21.9 | 0 | 0 | — |
| Reflect | 29.3 | 0 | 0 | 1.3× slower |
| ReflectFast | 2.82 | 0 | 0 | 7.8× faster |
| Offset | 5.79 | 0 | 0 | 3.8× faster |
| Accessor | 0.551 | 0 | 0 | 39.6× faster |
| AccOn | 1.94 | 0 | 0 | 11.2× faster |
| Direct | 0.275 | 0 | 0 | 79.4× faster |

## DI Resolve

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 0.497 | 0 | 0 | — |
| SafeRefl | 68.7 | 0 | 0 | 138.3× slower |
| Accessor | 3.09 | 0 | 0 | 6.2× slower |
| Reflect2 | 8.86 | 0 | 0 | 17.8× slower |
| Reflect | 108 | 0 | 0 | 218.2× slower |

## JSON Decode

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibJSON | 1537 | 392 | 10 | — |
| SafeRefl | 232 | 0 | 0 | 6.6× faster |
| Reflect | 719 | 0 | 0 | 2.1× faster |
| Reflect2 | 38.0 | 0 | 0 | 40.4× faster |
| Accessor | 3.04 | 0 | 0 | 505.1× faster |

## ORM Scan

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 5.58 | 0 | 0 | — |
| SafeRefl | 236 | 0 | 0 | 42.3× slower |
| Accessor | 5.51 | 0 | 0 | 1.0× faster |
| Reflect2 | 38.6 | 0 | 0 | 6.9× slower |
| Reflect | 452 | 0 | 0 | 81.0× slower |
| Copier | 3308 | 640 | 28 | 592.3× slower |

## Struct Copy

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 2.64 | 0 | 0 | — |
| Accessor | 4.26 | 0 | 0 | 1.6× slower |
| Reflect2 | 21.5 | 0 | 0 | 8.1× slower |
| Offset | 87.4 | 0 | 0 | 33.1× slower |
| Reflect | 314 | 0 | 0 | 118.8× slower |
| SafeRefl | 236 | 0 | 0 | 89.2× slower |
| Copier | 1359 | 432 | 17 | 514.3× slower |

<!-- bench:end -->

> Run `./scripts/bench.sh` to regenerate. Full results: [benchmarks/RESULTS.md](benchmarks/RESULTS.md).

Key numbers to read: **Accessor** = `saferefl.Accessor[T]` hot-path (pre-bound once, pointer arithmetic only); **AccFrom/AccOn** = Accessor with interface→pointer conversion per call; **SafeRefl** = `saferefl.Get[T]`/`Set[T]` (named field access per call, 0 allocs); **Offset** = pre-computed offset + `reflect.NewAt` (mechanism used internally by saferefl); **ReflectFast** = stdlib reflect with pre-cached `reflect.Value` (best possible reflect). In all benchmarks **Reflect** uses `FieldByName` per call — the standard reflect usage baseline. **Reflect2** uses pre-compiled field descriptors from [reflect2](https://github.com/modern-go/reflect2), representing a well-optimised codec that caches metadata at startup.

### Results by Go version

Per-version results are generated automatically by the [Cross-version Benchmarks](.github/workflows/bench-matrix.yml) workflow (runs weekly, or trigger manually).

| Go version | Results | Map backend |
|---|---|---|
| 1.22 | [benchmarks/results/go1.22.md](benchmarks/results/go1.22.md) | hmap |
| 1.24 | [benchmarks/results/go1.24.md](benchmarks/results/go1.24.md) | Swiss Tables |

To generate locally (requires Docker):

```
make bench-docker        # all versions
make bench-docker-1.24   # single version
```

## Status

| Layer | Status | Description |
|---|---|---|
| TypeInfo Cache | ✅ Done | `internal/typeinfo`: struct metadata + flat `IterPlan`, `sync.Map` + `atomic.Pointer` cache |
| Generic API | ✅ Done | `Get[T]`, `Set[T]`, `MustGet[T]`, `MustSet[T]`, dot-path, `FieldByName[T]`, `Fields`, `FieldsOf[T]` |
| Accessor API + Unsafe Primitives | ✅ Done | `Accessor[T]`, `UnsafeSliceAt[T]`, `MapLenFast` — self-test at init, hmap/Swiss Tables backends |
| General-purpose API | ✅ Done | `EachField`, `CopyFields`, `ToMap`, `ToMapByTag`, `FromMap`, `GetByTag[T]`, `SetByTag[T]`, `MapForEach`, `KindOf`, `IsNil` |

## Examples

Runnable examples are in [`examples/`](examples/):

| Directory | Demonstrates |
|---|---|
| [`examples/basic/`](examples/basic/) | `Get`/`Set` primitive fields |
| [`examples/dotpath/`](examples/dotpath/) | dot-path traversal through nested structs |
| [`examples/fields/`](examples/fields/) | field inspection without an instance |
| [`examples/accessor/`](examples/accessor/) | `Accessor[T]` hot-path pre-binding |
| [`examples/primitives/`](examples/primitives/) | `MapLenFast`, `UnsafeSliceAt` |
| [`examples/iteration/`](examples/iteration/) | `EachField`, `MapForEach` |
| [`examples/mapping/`](examples/mapping/) | `ToMap`, `ToMapByTag`, `FromMap` |
| [`examples/copyfields/`](examples/copyfields/) | `CopyFields` |
| [`examples/tags/`](examples/tags/) | `GetByTag`, `SetByTag` |
| [`examples/introspect/`](examples/introspect/) | `KindOf`, `IsNil` |

```
go run ./examples/iteration/
go run ./examples/mapping/
go run ./examples/tags/
```

## Go version support

Floor: **Go 1.22**. CI matrix: 1.22 (hmap backend), 1.24 (Swiss Tables backend), stable. Go tip runs weekly as a non-blocking early-warning job.
