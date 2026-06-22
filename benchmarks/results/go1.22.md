# Benchmark Results

**Generated:** 2026-06-22 18:11 UTC  
**Go:** 1.22  
**Platform:** linux/arm64  

## Get_int

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| L1 | 28.9 | 0 | 0 | — |
| Reflect | 1.71 | 0 | 0 | 16.9× faster |
| L2 | 5.19 | 0 | 0 | 5.6× faster |
| L3 | 0.554 | 0 | 0 | 52.2× faster |
| L3From | 1.67 | 0 | 0 | 17.3× faster |
| Direct | 0.271 | 0 | 0 | 106.7× faster |

## Set_string

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| L1 | 28.0 | 0 | 0 | — |
| Reflect | 2.81 | 0 | 0 | 10.0× faster |
| L2 | 6.58 | 0 | 0 | 4.3× faster |
| L3 | 0.818 | 0 | 0 | 34.2× faster |
| L3On | 2.19 | 0 | 0 | 12.8× faster |
| Direct | 0.274 | 0 | 0 | 102.3× faster |

## Field Read

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Reflect | 1.70 | 0 | 0 | — |
| L1 | 28.1 | 0 | 0 | 16.6× slower |
| L2 | 5.22 | 0 | 0 | 3.1× slower |
| L3 | 0.558 | 0 | 0 | 3.0× faster |
| Native | 0.276 | 0 | 0 | 6.1× faster |

## Slice At

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| L3 | 0.553 | 0 | 0 | — |
| Direct | 0.604 | 0 | 0 | 1.1× slower |
| Reflect | 1.88 | 0 | 0 | 3.4× slower |

## Map Len

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| L3 | 0.484 | 0 | 0 | — |
| Builtin | 0.286 | 0 | 0 | 1.7× faster |
| Reflect | 2.15 | 0 | 0 | 4.4× slower |

## DI Resolve

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 0.581 | 0 | 0 | — |
| L1 | 92.4 | 0 | 0 | 159.1× slower |
| L3 | 1.88 | 0 | 0 | 3.2× slower |
| Reflect2 | 8.59 | 0 | 0 | 14.8× slower |
| Reflect | 147 | 24 | 3 | 253.6× slower |

## JSON Decode

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibJSON | 1399 | 392 | 10 | — |
| L1 | 329 | 0 | 0 | 4.3× faster |
| Reflect | 1004 | 160 | 20 | 1.4× faster |
| Reflect2 | 36.5 | 0 | 0 | 38.3× faster |
| L3 | 4.28 | 0 | 0 | 326.6× faster |

## ORM Scan

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 9.37 | 0 | 0 | — |
| L1 | 333 | 0 | 0 | 35.5× slower |
| L3 | 6.81 | 0 | 0 | 1.4× faster |
| Reflect2 | 37.5 | 0 | 0 | 4.0× slower |
| Reflect | 606 | 80 | 10 | 64.7× slower |
| Copier | 3657 | 880 | 58 | 390.5× slower |

## Struct Copy

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 5.92 | 0 | 0 | — |
| L3 | 5.39 | 0 | 0 | 1.1× faster |
| Reflect2 | 18.6 | 0 | 0 | 3.1× slower |
| L2 | 87.9 | 0 | 0 | 14.9× slower |
| Reflect | 473 | 80 | 10 | 80.0× slower |
| L1 | 324 | 0 | 0 | 54.8× slower |
| Copier | 1589 | 552 | 32 | 268.6× slower |

