# Benchmark Results

**Generated:** 2026-06-22 12:17 UTC  
**Platform:** darwin/arm64  
**CPU:** Apple M3 Max  

## Get_int

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Saferefl | 27.5 | 0 | 0 | — |
| Reflect | 1.73 | 0 | 0 | 15.9× faster |
| Accessor | 0.546 | 0 | 0 | 50.4× faster |
| AccessorFrom | 4.63 | 0 | 0 | 5.9× faster |
| Direct | 0.272 | 0 | 0 | 100.9× faster |

## Set_string

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Saferefl | 38.7 | 16 | 1 | — |
| Reflect | 2.73 | 0 | 0 | 14.1× faster |
| Accessor | 0.543 | 0 | 0 | 71.2× faster |
| AccessorOn | 5.17 | 0 | 0 | 7.5× faster |

## Field Read

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Reflect | 1.73 | 0 | 0 | — |
| ReflectNewAt | 4.90 | 0 | 0 | 2.8× slower |
| Get | 27.0 | 0 | 0 | 15.6× slower |
| Accessor | 0.545 | 0 | 0 | 3.2× faster |
| Native | 0.273 | 0 | 0 | 6.4× faster |

## Slice At

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Layer3 | 0.549 | 0 | 0 | — |
| Direct | 0.548 | 0 | 0 | 1.0× faster |
| Reflect | 1.92 | 0 | 0 | 3.5× slower |

## Map Len

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Layer3 | 0.278 | 0 | 0 | — |
| Builtin | 0.274 | 0 | 0 | 1.0× faster |
| Reflect | 2.17 | 0 | 0 | 7.8× slower |

## DI Resolve

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 0.483 | 0 | 0 | — |
| Saferefl | 105 | 24 | 3 | 217.6× slower |
| Accessor | 1.09 | 0 | 0 | 2.3× slower |
| Reflect | 106 | 0 | 0 | 218.4× slower |

## JSON Decode

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibJSON | 1492 | 392 | 10 | — |
| Saferefl | 389 | 120 | 10 | 3.8× faster |
| Reflect | 105 | 0 | 0 | 14.2× faster |
| Reflect2 | 33.6 | 0 | 0 | 44.4× faster |

## ORM Scan

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 5.54 | 0 | 0 | — |
| Saferefl | 383 | 112 | 10 | 69.2× slower |
| Accessor | 6.70 | 0 | 0 | 1.2× slower |
| Reflect | 440 | 0 | 0 | 79.5× slower |
| Copier | 3241 | 640 | 28 | 585.4× slower |

## Struct Copy

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 2.57 | 0 | 0 | — |
| Saferefl | 350 | 56 | 5 | 136.0× slower |
| Reflect | 57.2 | 0 | 0 | 22.2× slower |
| Accessor | 3.31 | 0 | 0 | 1.3× slower |
| Copier | 1340 | 432 | 17 | 521.1× slower |

## JSON Like

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibReflect | 192 | 64 | 1 | — |
| Reflect2 | 30.3 | 64 | 1 | 6.3× faster |
| CachedOffset | 80.2 | 64 | 1 | 2.4× faster |
| Generics | 26.0 | 64 | 1 | 7.4× faster |

## Read Int64

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibReflect | 26.1 | 0 | 0 | — |
| Reflect2 | 0.848 | 0 | 0 | 30.7× faster |
| CachedOffset | 5.73 | 0 | 0 | 4.5× faster |
| Generics | 0.287 | 0 | 0 | 90.7× faster |

## Write Int64

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibReflect | 27.4 | 0 | 0 | — |
| Reflect2 | 2.99 | 0 | 0 | 9.2× faster |
| CachedOffset | 10.9 | 0 | 0 | 2.5× faster |
| Generics | 0.273 | 0 | 0 | 100.5× faster |

## Read String

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibReflect | 27.4 | 0 | 0 | — |
| Reflect2 | 0.889 | 0 | 0 | 30.8× faster |
| CachedOffset | 5.18 | 0 | 0 | 5.3× faster |
| Generics | 0.546 | 0 | 0 | 50.1× faster |

