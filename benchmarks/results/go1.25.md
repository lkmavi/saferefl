# Benchmark Results

**Generated:** 2026-06-22 18:26 UTC  
**Go:** 1.25  
**Platform:** linux/arm64  

## Get_int

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| L1 | 28.0 | 0 | 0 | — |
| Reflect | 1.67 | 0 | 0 | 16.7× faster |
| L2 | 4.81 | 0 | 0 | 5.8× faster |
| L3 | 0.534 | 0 | 0 | 52.4× faster |
| L3From | 1.61 | 0 | 0 | 17.4× faster |
| Direct | 0.267 | 0 | 0 | 104.9× faster |

## Set_string

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| L1 | 27.6 | 0 | 0 | — |
| Reflect | 2.46 | 0 | 0 | 11.2× faster |
| L2 | 5.62 | 0 | 0 | 4.9× faster |
| L3 | 0.535 | 0 | 0 | 51.5× faster |
| L3On | 1.88 | 0 | 0 | 14.7× faster |
| Direct | 0.267 | 0 | 0 | 103.3× faster |

## Field Read

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Reflect | 1.67 | 0 | 0 | — |
| L1 | 27.4 | 0 | 0 | 16.4× slower |
| L2 | 4.82 | 0 | 0 | 2.9× slower |
| L3 | 0.536 | 0 | 0 | 3.1× faster |
| Native | 0.269 | 0 | 0 | 6.2× faster |

## Slice At

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| L3 | 0.538 | 0 | 0 | — |
| Direct | 0.585 | 0 | 0 | 1.1× slower |
| Reflect | 1.91 | 0 | 0 | 3.6× slower |

## Map Len

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| L3 | 0.278 | 0 | 0 | — |
| Builtin | 0.276 | 0 | 0 | 1.0× faster |
| Reflect | 2.14 | 0 | 0 | 7.7× slower |

## DI Resolve

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 0.496 | 0 | 0 | — |
| L1 | 95.2 | 0 | 0 | 192.0× slower |
| L3 | 1.40 | 0 | 0 | 2.8× slower |
| Reflect2 | 8.55 | 0 | 0 | 17.2× slower |
| Reflect | 128 | 0 | 0 | 258.4× slower |

## JSON Decode

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibJSON | 1447 | 392 | 10 | — |
| L1 | 285 | 0 | 0 | 5.1× faster |
| Reflect | 859 | 0 | 0 | 1.7× faster |
| Reflect2 | 36.5 | 0 | 0 | 39.6× faster |
| L3 | 2.95 | 0 | 0 | 491.2× faster |

## ORM Scan

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 8.40 | 0 | 0 | — |
| L1 | 284 | 0 | 0 | 33.8× slower |
| L3 | 6.61 | 0 | 0 | 1.3× faster |
| Reflect2 | 37.8 | 0 | 0 | 4.5× slower |
| Reflect | 504 | 0 | 0 | 60.1× slower |
| Copier | 3349 | 640 | 28 | 398.7× slower |

## Struct Copy

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 5.99 | 0 | 0 | — |
| L3 | 5.72 | 0 | 0 | 1.0× faster |
| Reflect2 | 21.1 | 0 | 0 | 3.5× slower |
| L2 | 87.8 | 0 | 0 | 14.7× slower |
| Reflect | 410 | 0 | 0 | 68.4× slower |
| L1 | 306 | 0 | 0 | 51.2× slower |
| Copier | 1403 | 432 | 17 | 234.2× slower |

