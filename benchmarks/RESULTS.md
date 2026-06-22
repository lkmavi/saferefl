# Benchmark Results

**Generated:** 2026-06-22 18:03 UTC  
**Go:** 1.26  
**Platform:** darwin/arm64  
**CPU:** Apple M3 Max  

## Get_int

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| L1 | 26.1 | 0 | 0 | — |
| Reflect | 1.69 | 0 | 0 | 15.5× faster |
| L2 | 4.80 | 0 | 0 | 5.4× faster |
| L3 | 0.532 | 0 | 0 | 49.1× faster |
| L3From | 1.42 | 0 | 0 | 18.4× faster |
| Direct | 0.267 | 0 | 0 | 97.9× faster |

## Set_string

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| L1 | 26.5 | 0 | 0 | — |
| Reflect | 2.67 | 0 | 0 | 9.9× faster |
| L2 | 5.59 | 0 | 0 | 4.7× faster |
| L3 | 0.534 | 0 | 0 | 49.7× faster |
| L3On | 1.86 | 0 | 0 | 14.2× faster |
| Direct | 0.265 | 0 | 0 | 100.1× faster |

## Field Read

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Reflect | 1.69 | 0 | 0 | — |
| L1 | 26.1 | 0 | 0 | 15.5× slower |
| L2 | 5.06 | 0 | 0 | 3.0× slower |
| L3 | 0.532 | 0 | 0 | 3.2× faster |
| Native | 0.266 | 0 | 0 | 6.3× faster |

## Slice At

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| L3 | 0.535 | 0 | 0 | — |
| Direct | 0.536 | 0 | 0 | 1.0× slower |
| Reflect | 1.87 | 0 | 0 | 3.5× slower |

## Map Len

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| L3 | 0.273 | 0 | 0 | — |
| Builtin | 0.271 | 0 | 0 | 1.0× faster |
| Reflect | 2.13 | 0 | 0 | 7.8× slower |

## DI Resolve

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 0.474 | 0 | 0 | — |
| L1 | 78.6 | 0 | 0 | 165.8× slower |
| L3 | 3.22 | 0 | 0 | 6.8× slower |
| Reflect2 | 8.53 | 0 | 0 | 18.0× slower |
| Reflect | 103 | 0 | 0 | 217.6× slower |

## JSON Decode

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibJSON | 1434 | 392 | 10 | — |
| L1 | 270 | 0 | 0 | 5.3× faster |
| Reflect | 687 | 0 | 0 | 2.1× faster |
| Reflect2 | 36.3 | 0 | 0 | 39.5× faster |
| L3 | 2.93 | 0 | 0 | 489.6× faster |

## ORM Scan

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 5.28 | 0 | 0 | — |
| L1 | 271 | 0 | 0 | 51.3× slower |
| L3 | 5.33 | 0 | 0 | 1.0× slower |
| Reflect2 | 37.0 | 0 | 0 | 7.0× slower |
| Reflect | 431 | 0 | 0 | 81.6× slower |
| Copier | 3109 | 640 | 28 | 588.3× slower |

## Struct Copy

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 2.49 | 0 | 0 | — |
| L3 | 3.92 | 0 | 0 | 1.6× slower |
| Reflect2 | 21.0 | 0 | 0 | 8.4× slower |
| L2 | 84.2 | 0 | 0 | 33.8× slower |
| Reflect | 305 | 0 | 0 | 122.2× slower |
| L1 | 286 | 0 | 0 | 114.7× slower |
| Copier | 1322 | 432 | 17 | 530.3× slower |

