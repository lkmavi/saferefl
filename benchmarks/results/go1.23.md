# Benchmark Results

**Generated:** 2026-06-22 18:16 UTC  
**Go:** 1.23  
**Platform:** linux/arm64  

## Get_int

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| L1 | 29.3 | 0 | 0 | — |
| Reflect | 1.67 | 0 | 0 | 17.5× faster |
| L2 | 5.37 | 0 | 0 | 5.5× faster |
| L3 | 0.535 | 0 | 0 | 54.8× faster |
| L3From | 1.61 | 0 | 0 | 18.2× faster |
| Direct | 0.267 | 0 | 0 | 109.8× faster |

## Set_string

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| L1 | 29.1 | 0 | 0 | — |
| Reflect | 3.03 | 0 | 0 | 9.6× faster |
| L2 | 6.15 | 0 | 0 | 4.7× faster |
| L3 | 0.802 | 0 | 0 | 36.4× faster |
| L3On | 2.15 | 0 | 0 | 13.6× faster |
| Direct | 0.267 | 0 | 0 | 109.1× faster |

## Field Read

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Reflect | 1.67 | 0 | 0 | — |
| L1 | 28.6 | 0 | 0 | 17.2× slower |
| L2 | 5.09 | 0 | 0 | 3.0× slower |
| L3 | 0.535 | 0 | 0 | 3.1× faster |
| Native | 0.267 | 0 | 0 | 6.3× faster |

## Slice At

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| L3 | 0.538 | 0 | 0 | — |
| Direct | 0.588 | 0 | 0 | 1.1× slower |
| Reflect | 1.85 | 0 | 0 | 3.4× slower |

## Map Len

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| L3 | 0.471 | 0 | 0 | — |
| Builtin | 0.279 | 0 | 0 | 1.7× faster |
| Reflect | 2.14 | 0 | 0 | 4.5× slower |

## DI Resolve

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 0.585 | 0 | 0 | — |
| L1 | 99.8 | 0 | 0 | 170.5× slower |
| L3 | 1.87 | 0 | 0 | 3.2× slower |
| Reflect2 | 8.57 | 0 | 0 | 14.6× slower |
| Reflect | 147 | 24 | 3 | 250.3× slower |

## JSON Decode

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibJSON | 1377 | 392 | 10 | — |
| L1 | 340 | 0 | 0 | 4.1× faster |
| Reflect | 985 | 160 | 20 | 1.4× faster |
| Reflect2 | 36.5 | 0 | 0 | 37.7× faster |
| L3 | 4.28 | 0 | 0 | 321.8× faster |

## ORM Scan

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 9.37 | 0 | 0 | — |
| L1 | 333 | 0 | 0 | 35.6× slower |
| L3 | 6.83 | 0 | 0 | 1.4× faster |
| Reflect2 | 37.4 | 0 | 0 | 4.0× slower |
| Reflect | 604 | 80 | 10 | 64.5× slower |
| Copier | 3796 | 880 | 58 | 405.2× slower |

## Struct Copy

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 5.62 | 0 | 0 | — |
| L3 | 5.40 | 0 | 0 | 1.0× faster |
| Reflect2 | 18.6 | 0 | 0 | 3.3× slower |
| L2 | 87.1 | 0 | 0 | 15.5× slower |
| Reflect | 469 | 80 | 10 | 83.4× slower |
| L1 | 336 | 0 | 0 | 59.8× slower |
| Copier | 1674 | 552 | 32 | 297.7× slower |

