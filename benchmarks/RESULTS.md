# Benchmark Results

**Generated:** 2026-06-23 08:26 UTC  
**Go:** 1.26  
**Platform:** darwin/arm64  
**CPU:** Apple M3 Max  

## Field Read

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Reflect | 26.8 | 0 | 0 | — |
| ReflectFast | 1.75 | 0 | 0 | 15.3× faster |
| SafeRefl | 21.2 | 0 | 0 | 1.3× faster |
| Offset | 5.00 | 0 | 0 | 5.4× faster |
| Accessor | 0.555 | 0 | 0 | 48.2× faster |
| Native | 0.278 | 0 | 0 | 96.3× faster |

## Slice At

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.564 | 0 | 0 | — |
| Direct | 0.548 | 0 | 0 | 1.0× faster |
| Reflect | 1.91 | 0 | 0 | 3.4× slower |

## Map Len

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.277 | 0 | 0 | — |
| Builtin | 0.274 | 0 | 0 | 1.0× faster |
| Reflect | 2.20 | 0 | 0 | 7.9× slower |

## Get_int

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 21.2 | 0 | 0 | — |
| Reflect | 26.1 | 0 | 0 | 1.2× slower |
| ReflectFast | 1.75 | 0 | 0 | 12.1× faster |
| Offset | 5.03 | 0 | 0 | 4.2× faster |
| Accessor | 0.546 | 0 | 0 | 38.8× faster |
| AccFrom | 1.45 | 0 | 0 | 14.6× faster |
| Direct | 0.283 | 0 | 0 | 75.0× faster |

## Set_string

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 21.9 | 0 | 0 | — |
| Reflect | 29.3 | 0 | 0 | 1.3× slower |
| ReflectFast | 2.82 | 0 | 0 | 7.8× faster |
| Offset | 5.79 | 0 | 0 | 3.8× faster |
| Accessor | 0.551 | 0 | 0 | 39.6× faster |
| AccOn | 1.94 | 0 | 0 | 11.2× faster |
| Direct | 0.275 | 0 | 0 | 79.4× faster |

## DI Resolve

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 0.497 | 0 | 0 | — |
| SafeRefl | 68.7 | 0 | 0 | 138.3× slower |
| Accessor | 3.09 | 0 | 0 | 6.2× slower |
| Reflect2 | 8.86 | 0 | 0 | 17.8× slower |
| Reflect | 108 | 0 | 0 | 218.2× slower |

## JSON Decode

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibJSON | 1537 | 392 | 10 | — |
| SafeRefl | 232 | 0 | 0 | 6.6× faster |
| Reflect | 719 | 0 | 0 | 2.1× faster |
| Reflect2 | 38.0 | 0 | 0 | 40.4× faster |
| Accessor | 3.04 | 0 | 0 | 505.1× faster |

## ORM Scan

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 5.58 | 0 | 0 | — |
| SafeRefl | 236 | 0 | 0 | 42.3× slower |
| Accessor | 5.51 | 0 | 0 | 1.0× faster |
| Reflect2 | 38.6 | 0 | 0 | 6.9× slower |
| Reflect | 452 | 0 | 0 | 81.0× slower |
| Copier | 3308 | 640 | 28 | 592.3× slower |

## Struct Copy

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 2.64 | 0 | 0 | — |
| Accessor | 4.26 | 0 | 0 | 1.6× slower |
| Reflect2 | 21.5 | 0 | 0 | 8.1× slower |
| Offset | 87.4 | 0 | 0 | 33.1× slower |
| Reflect | 314 | 0 | 0 | 118.8× slower |
| SafeRefl | 236 | 0 | 0 | 89.2× slower |
| Copier | 1359 | 432 | 17 | 514.3× slower |

