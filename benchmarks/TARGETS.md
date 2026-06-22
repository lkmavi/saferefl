# Performance Targets

**Baseline date:** 2026-06-22  
**Status:** pre-Layer 3 (Layer 1 + Layer 2 only)

These targets must be met before the Layer 3 unsafe accelerator is considered production-ready.
CI enforces them via the `bench-regression` job on every PR.

---

## Acceptance Criteria

| Layer | Comparison | Target |
|---|---|---|
| Layer 1 (statically known type) | vs native direct field access | ≤ 1.2× slower |
| Layer 2 (dynamic, no unsafe) | vs reflect2 (cached offsets) | ≤ 2.0× slower |
| Layer 3 (unsafe accel, opt-in) | vs reflect2 (cached offsets) | ≤ 1.1× slower |

**Layer 1** is the `Get[T]`/`Set[T]` public API with TypeInfo cache warm.  
**Layer 2** is the internal `typeinfo.GetFieldPtr` + `reflect.NewAt` path.  
**Layer 3** is the optional `internal/unsafelayout` accelerator (Phase 4).

---

## Realistic Scenario Targets

| Scenario | Metric | Target (Layer 1) |
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
