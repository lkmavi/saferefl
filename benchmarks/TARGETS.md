# Performance Targets

**Baseline date:** 2026-06-22  
**Status:** All APIs shipped

These targets define the performance envelope for each part of the API.

---

## Acceptance Criteria

| API | Comparison | Target |
|---|---|---|
| `Get[T]`/`Set[T]` (warm cache) | vs `reflect.FieldByName` per call | ≤ 1× slower (must be faster) |
| `Accessor[T]` hot path | vs native direct field access | ≤ 2× slower |
| `MapLenFast` | vs builtin `len(m)` | ≤ 1.5× slower |

**Generic API** (`Get[T]`/`Set[T]`) resolves the field path on every call but uses a cached TypeDescriptor.  
**Accessor API** (`MakeAccessor` once, `.Get`/`.Set` per call) is pure pointer arithmetic with no per-call reflection.  
**Unsafe Primitives** (`MapLenFast`, `UnsafeSliceAt`) read runtime internals directly; self-tested at package init.

---

## Realistic Scenario Targets

| Scenario | Metric | Target (Accessor API) |
|---|---|---|
| Struct-to-struct copy (5 fields) | ns/op | ≤ 5× vs Manual |
| JSON-like field assignment (10 fields) | ns/op | ≤ 3× vs reflect2 |
| DI resolve (inject 3 deps) | ns/op | ≤ 5× vs Manual |
| ORM row scan (10 columns) | ns/op | ≤ 3× vs reflect2 |

All scenarios: `allocs/op` must not exceed the reflect2 baseline for the same scenario.

---

## How to Verify

```bash
make bench
```

To compare two runs statistically:

```bash
make bench-stat OLD=bench-main.json NEW=bench-pr.json
```

---

## Notes

- Targets apply with TypeInfo cache warm (first-call overhead is excluded by design).
- `encoding/json` in the JSON-decode scenario is a full parse+assign baseline, not a fair
  direct competitor — it is included only for context.
- Allocation targets may be relaxed for interface-typed fields (boxing is unavoidable there).
