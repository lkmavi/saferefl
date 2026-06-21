# saferefl

Fast, safe reflection for Go — a generic-first alternative to [`reflect2`](https://github.com/modern-go/reflect2).

## Why

`reflect2` trades correctness for speed by reverse-engineering Go's internal runtime layout. This causes silent data corruption when Go internals change (e.g. the map rewrite to Swiss Tables in Go 1.24). `saferefl` gets comparable speed through a different route: generics, cached offsets, and a self-verifying unsafe layer that falls back gracefully instead of corrupting memory.

See [ADR-01](https://github.com/lkmavi/saferefl/discussions/3) for the full analysis and decision.

## Architecture

```
┌─────────────────────────────────────────────┐
│  Layer 1: Public API (generics-first)        │  Get[T], Set[T]
│  - type-safe by construction                 │
│  - zero unsafe for statically known types    │
├─────────────────────────────────────────────┤
│  Layer 2: Cached reflective layer            │  TypeDescriptorOf(t)
│  - built once via stdlib reflect             │
│  - sync.Map cache on reflect.Type            │
├─────────────────────────────────────────────┤
│  Layer 3: Optional unsafe accelerator        │  build tag: unsafe_accel
│  - self-test at init, auto-fallback on fail  │
│  - two map backends: hmap (<1.24), swiss (≥1.24) │
└─────────────────────────────────────────────┘
```

## Performance

<!-- bench:start -->
# Benchmark Results

**Generated:** 2026-06-21 12:40 UTC  
**Platform:** darwin/arm64  
**CPU:** Apple M3 Max  

## JSON Like

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibReflect | 193 | 64 | 1 | — |
| Reflect2 | 30.2 | 64 | 1 | 6.4× faster |
| CachedOffset | 80.1 | 64 | 1 | 2.4× faster |
| Generics | 26.2 | 64 | 1 | 7.4× faster |

## Read Int64

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibReflect | 25.9 | 0 | 0 | — |
| Reflect2 | 0.848 | 0 | 0 | 30.5× faster |
| CachedOffset | 5.71 | 0 | 0 | 4.5× faster |
| Generics | 0.286 | 0 | 0 | 90.4× faster |

## Write Int64

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibReflect | 27.1 | 0 | 0 | — |
| Reflect2 | 3.00 | 0 | 0 | 9.0× faster |
| CachedOffset | 10.9 | 0 | 0 | 2.5× faster |
| Generics | 0.264 | 0 | 0 | 102.6× faster |

## Read String

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibReflect | 27.1 | 0 | 0 | — |
| Reflect2 | 0.892 | 0 | 0 | 30.4× faster |
| CachedOffset | 5.25 | 0 | 0 | 5.2× faster |
| Generics | 0.546 | 0 | 0 | 49.7× faster |

<!-- bench:end -->

> Run `./scripts/bench.sh` to regenerate. Full results: [benchmarks/RESULTS.md](benchmarks/RESULTS.md).

## Status

Work in progress. See the [implementation plan](_local/IMPL-PLAN-01.md).

## Go version support

Floor: **Go 1.22**. CI matrix: 1.22, 1.23, 1.24, stable, tip.
