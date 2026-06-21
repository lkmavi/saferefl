# Spike Results

**Date:** 2026-06-21  
**Machine:** Apple M3 Max, darwin/arm64, Go 1.26.3  
**Command:** `go test ./benchmarks/spike/ -bench=. -benchmem -count=5 -benchtime=1s`

---

## Single-field Read — int64

| Approach | ns/op | allocs/op | vs Stdlib | vs Reflect2 |
|---|---|---|---|---|
| StdlibReflect (`FieldByName` + `.Int()`) | 25.5 | 0 | baseline | 31× slower |
| Reflect2 (`UnsafeGet` + cast) | 0.83 | 0 | 31× faster | baseline |
| CachedOffset (`reflect.NewAt` + `.Int()`) | 5.6 | 0 | 4.6× faster | 6.7× slower |
| **Generics** (direct `*(*int64)(ptr + offset)`) | **0.28** | **0** | **91× faster** | **3× faster** |

## Single-field Write — int64

| Approach | ns/op | allocs/op | vs Stdlib | vs Reflect2 |
|---|---|---|---|---|
| StdlibReflect (`FieldByName` + `.SetInt()`) | 26.9 | 0 | baseline | 9.2× slower |
| Reflect2 (`UnsafeSet`) | 2.93 | 0 | 9.2× faster | baseline |
| CachedOffset (`reflect.NewAt` + `.Set()`) | 10.7 | 0 | 2.5× faster | 3.7× slower |
| **Generics** (direct `*(*int64)(ptr + offset) = val`) | **0.26** | **0** | **103× faster** | **11× faster** |

## Single-field Read — string

| Approach | ns/op | allocs/op | vs Stdlib | vs Reflect2 |
|---|---|---|---|---|
| StdlibReflect (`FieldByName` + `.String()`) | 26.9 | 0 | baseline | 31× slower |
| Reflect2 (`UnsafeGet` + cast) | 0.87 | 0 | 31× faster | baseline |
| CachedOffset (`reflect.NewAt` + `.String()`) | 5.2 | 0 | 5.2× faster | 6× slower |
| **Generics** (direct `*(*string)(ptr + offset)`) | **0.53** | **0** | **51× faster** | **1.6× faster** |

## JSON-like Decode — 5 fields (ID int64, Name/Email string, Score float64, Active bool)

Each iteration allocates a fresh `User{}` (simulates per-object decoding).

| Approach | ns/op | allocs/op | B/op | vs Stdlib | vs Reflect2 |
|---|---|---|---|---|---|
| StdlibReflect (`FieldByName` × 5) | 187 | 1 | 64 | baseline | 6.3× slower |
| Reflect2 (`UnsafeSet` × 5) | 29.5 | 1 | 64 | 6.3× faster | baseline |
| CachedOffset (`reflect.NewAt` × 5) | 78.5 | 1 | 64 | 2.4× faster | 2.7× slower |
| **Generics** (type-specific setters × 5) | **25.0** | **1** | **64** | **7.5× faster** | **1.2× faster** |

*The single alloc in all JSON-like variants is the `User{}` struct being decoded into — equal across all approaches.*

---

## Key Findings

### 1. Generics path beats reflect2 across the board

On single-field ops, `readField[int64]` (direct cast via pre-computed offset) is **3× faster than reflect2.UnsafeGet** for reads and **11× faster for writes**. This is because reflect2's `UnsafeSet` goes through a virtual method dispatch on the internal `Type` hierarchy, while the generic helper gets inlined by the compiler to a single load/store instruction.

**Implication for Layer 1:** for statically-typed callers (`Get[T]/Set[T]`), we get better-than-reflect2 performance without any unsafe internal layout access.

### 2. CachedOffset (Layer 2) is the bottleneck, not the allocation

CachedOffset is **2.5–6.7× slower than reflect2** despite having 0 allocs. The overhead is `reflect.NewAt(typ, ptr).Elem().Int/String/Set()` — approximately 5–10 ns per field. This is primarily:
- Two `reflect.Value` struct constructions (stack, but not free on ARM)
- Flag-checking in `.Int()` / `.String()` / `.Set()`

**Implication for Layer 2:** acceptable for moderate-frequency use (DI, ORM row scan), but not for hot-path JSON decoding where Layer 3 is needed.

### 3. CachedOffset is ~2.4× faster than stdlib in JSON-like scenario

Even without unsafe, pre-computing offsets and using `reflect.NewAt` cuts the JSON-like time from 187 ns to 78.5 ns. This validates the "cache field metadata once" hypothesis from the ADR.

### 4. `reflect.Value.FieldByName` dominates stdlib cost

Single-field stdlib read: 25.5 ns. This is almost entirely `FieldByName` — the hash map lookup in `reflect.Type`. Replacing it with a cached index or pre-computed offset removes ~80% of the cost.

---

## Revised Performance Targets for TARGETS.md

Based on these numbers (Apple M3 Max, Go 1.26.3):

| Layer | Operation | Target (ns/op) | Basis |
|---|---|---|---|
| Layer 1 (Generics) | single-field read/write | ≤ 1 ns | measured 0.26–0.55 ns |
| Layer 2 (CachedOffset) | single-field read/write | ≤ 12 ns | measured 5.2–10.7 ns |
| Layer 2 (CachedOffset) | JSON-like 5-field decode | ≤ 100 ns | measured 78.5 ns |
| Layer 3 (Unsafe accel) | single-field read/write | ≤ 1 ns | same as Layer 1 (same code path) |
| Layer 3 (Unsafe accel) | JSON-like 5-field decode | ≤ 35 ns | measured 25 ns (generics) |

---

## Open Questions Raised by Spike

1. **reflect.NewAt overhead source**: 5.6 ns vs 0.28 ns — is this flag-checking, or does `reflect.NewAt` not get inlined on ARM? Worth profiling with `pprof` before optimizing Layer 2.

2. **reflect2.UnsafeSet write overhead (2.93 ns vs 0.26 ns)**: reflect2 goes through a `typeInfo.Set` virtual dispatch. Our Layer 3 `writeField[T]` avoids this entirely. This means **Layer 3 is not just "equivalent to reflect2" — it's faster**, and Layer 3 might even be overengineered for write-heavy workloads where Layer 1 alone suffices.

3. **JSON-like alloc dominance**: all 4 approaches show 1 alloc / 64 B — the `User{}` struct. In a real decoder with a pre-allocated target (like `json.Unmarshal(&u, data)`), the alloc goes away and the field-setting cost becomes more visible. Consider a no-alloc variant of the JSON-like benchmark for Phase 3.

4. **reflect2 on Go 1.26**: reflect2 v1.0.2 compiled and ran correctly on Go 1.26.3. The map layout change in 1.24 (Swiss Tables) did not affect struct field operations (only `map` built-in), confirming the ADR analysis.
