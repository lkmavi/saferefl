# Benchmark Results

**Generated:** 2026-06-22 18:30 UTC  
**Go:** 1.26  
**Platform:** linux/arm64  

## Get_int

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| L1 | 27.2 | 0 | 0 | — |
| Reflect | 1.72 | 0 | 0 | 15.9× faster |
| L2 | 4.84 | 0 | 0 | 5.6× faster |
| L3 | 0.548 | 0 | 0 | 49.7× faster |
| L3From | 1.46 | 0 | 0 | 18.6× faster |
| Direct | 0.277 | 0 | 0 | 98.2× faster |

## Set_string

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| L1 | 27.0 | 0 | 0 | — |
| Reflect | 2.74 | 0 | 0 | 9.8× faster |
| L2 | 5.70 | 0 | 0 | 4.7× faster |
| L3 | 0.549 | 0 | 0 | 49.1× faster |
| L3On | 1.92 | 0 | 0 | 14.1× faster |
| Direct | 0.274 | 0 | 0 | 98.6× faster |

## Field Read

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Reflect | 1.73 | 0 | 0 | — |
| L1 | 26.8 | 0 | 0 | 15.5× slower |
| L2 | 4.93 | 0 | 0 | 2.9× slower |
| L3 | 0.539 | 0 | 0 | 3.2× faster |
| Native | 0.269 | 0 | 0 | 6.4× faster |

## Slice At

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| L3 | 0.562 | 0 | 0 | — |
| Direct | 0.551 | 0 | 0 | 1.0× faster |
| Reflect | 1.93 | 0 | 0 | 3.4× slower |

## Map Len

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| L3 | 0.279 | 0 | 0 | — |
| Builtin | 0.279 | 0 | 0 | 1.0× faster |
| Reflect | 2.15 | 0 | 0 | 7.7× slower |

## DI Resolve

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 0.497 | 0 | 0 | — |
| L1 | 79.2 | 0 | 0 | 159.5× slower |
| L3 | 1.46 | 0 | 0 | 2.9× slower |
| Reflect2 | 12.0 | 0 | 0 | 24.2× slower |
| Reflect | 138 | 0 | 0 | 277.1× slower |

## JSON Decode

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibJSON | 1556 | 392 | 10 | — |
| L1 | 282 | 0 | 0 | 5.5× faster |
| Reflect | 931 | 0 | 0 | 1.7× faster |
| Reflect2 | 38.5 | 0 | 0 | 40.4× faster |
| L3 | 5.31 | 0 | 0 | 292.9× faster |

## ORM Scan

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 9.06 | 0 | 0 | — |
| L1 | 272 | 0 | 0 | 30.0× slower |
| L3 | 8.25 | 0 | 0 | 1.1× faster |
| Reflect2 | 45.5 | 0 | 0 | 5.0× slower |
| Reflect | 519 | 0 | 0 | 57.3× slower |
| Copier | 3386 | 640 | 28 | 373.7× slower |

## Struct Copy

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 7.03 | 0 | 0 | — |
| L3 | 6.91 | 0 | 0 | 1.0× faster |
| Reflect2 | 25.3 | 0 | 0 | 3.6× slower |
| L2 | 88.1 | 0 | 0 | 12.5× slower |
| Reflect | 421 | 0 | 0 | 59.8× slower |
| L1 | 282 | 0 | 0 | 40.1× slower |
| Copier | 1424 | 432 | 17 | 202.6× slower |

