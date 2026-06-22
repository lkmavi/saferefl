# Benchmark Results

**Generated:** 2026-06-22 18:21 UTC  
**Go:** 1.24  
**Platform:** linux/arm64  

## Get_int

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| L1 | 29.2 | 0 | 0 | — |
| Reflect | 1.67 | 0 | 0 | 17.4× faster |
| L2 | 4.81 | 0 | 0 | 6.1× faster |
| L3 | 0.534 | 0 | 0 | 54.6× faster |
| L3From | 1.61 | 0 | 0 | 18.2× faster |
| Direct | 0.267 | 0 | 0 | 109.3× faster |

## Set_string

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| L1 | 28.8 | 0 | 0 | — |
| Reflect | 2.46 | 0 | 0 | 11.7× faster |
| L2 | 5.99 | 0 | 0 | 4.8× faster |
| L3 | 0.802 | 0 | 0 | 35.9× faster |
| L3On | 2.14 | 0 | 0 | 13.4× faster |
| Direct | 0.266 | 0 | 0 | 108.0× faster |

## Field Read

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Reflect | 1.66 | 0 | 0 | — |
| L1 | 28.5 | 0 | 0 | 17.1× slower |
| L2 | 4.82 | 0 | 0 | 2.9× slower |
| L3 | 0.536 | 0 | 0 | 3.1× faster |
| Native | 0.267 | 0 | 0 | 6.2× faster |

## Slice At

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| L3 | 0.537 | 0 | 0 | — |
| Direct | 0.588 | 0 | 0 | 1.1× slower |
| Reflect | 1.85 | 0 | 0 | 3.4× slower |

## Map Len

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| L3 | 0.471 | 0 | 0 | — |
| Builtin | 0.282 | 0 | 0 | 1.7× faster |
| Reflect | 2.14 | 0 | 0 | 4.5× slower |

## DI Resolve

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 0.583 | 0 | 0 | — |
| L1 | 81.9 | 0 | 0 | 140.4× slower |
| L3 | 1.87 | 0 | 0 | 3.2× slower |
| Reflect2 | 8.57 | 0 | 0 | 14.7× slower |
| Reflect | 133 | 0 | 0 | 228.0× slower |

## JSON Decode

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibJSON | 1456 | 392 | 10 | — |
| L1 | 318 | 0 | 0 | 4.6× faster |
| Reflect | 888 | 0 | 0 | 1.6× faster |
| Reflect2 | 36.6 | 0 | 0 | 39.8× faster |
| L3 | 4.29 | 0 | 0 | 339.4× faster |

## ORM Scan

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 9.29 | 0 | 0 | — |
| L1 | 342 | 0 | 0 | 36.8× slower |
| L3 | 6.79 | 0 | 0 | 1.4× faster |
| Reflect2 | 37.8 | 0 | 0 | 4.1× slower |
| Reflect | 518 | 0 | 0 | 55.8× slower |
| Copier | 3348 | 640 | 28 | 360.3× slower |

## Struct Copy

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 5.91 | 0 | 0 | — |
| L3 | 5.40 | 0 | 0 | 1.1× faster |
| Reflect2 | 18.6 | 0 | 0 | 3.1× slower |
| L2 | 91.3 | 0 | 0 | 15.4× slower |
| Reflect | 419 | 0 | 0 | 70.8× slower |
| L1 | 326 | 0 | 0 | 55.2× slower |
| Copier | 1407 | 432 | 17 | 238.0× slower |

