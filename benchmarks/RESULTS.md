# Benchmark Results

**Generated:** 2026-06-22 11:22 UTC  
**Platform:** darwin/arm64  
**CPU:** Apple M3 Max  

## Get_int

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Saferefl | 26.1 | 0 | 0 | — |
| Reflect | 1.75 | 0 | 0 | 14.9× faster |
| Accessor | 0.550 | 0 | 0 | 47.4× faster |
| AccessorFrom | 4.38 | 0 | 0 | 6.0× faster |
| Direct | 0.269 | 0 | 0 | 97.1× faster |

## Set_string

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Saferefl | 38.6 | 16 | 1 | — |
| Reflect | 2.72 | 0 | 0 | 14.2× faster |
| Accessor | 0.538 | 0 | 0 | 71.8× faster |
| AccessorOn | 4.86 | 0 | 0 | 8.0× faster |

## DI Resolve

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 0.474 | 0 | 0 | — |
| Saferefl | 104 | 24 | 3 | 219.5× slower |
| Reflect | 104 | 0 | 0 | 220.2× slower |

## JSON Decode

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibJSON | 1495 | 392 | 10 | — |
| Saferefl | 400 | 120 | 10 | 3.7× faster |
| Reflect | 105 | 0 | 0 | 14.2× faster |
| Reflect2 | 33.6 | 0 | 0 | 44.5× faster |

## ORM Scan

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 5.48 | 0 | 0 | — |
| Saferefl | 397 | 112 | 10 | 72.4× slower |
| Reflect | 448 | 0 | 0 | 81.7× slower |
| Copier | 3254 | 640 | 28 | 593.6× slower |

## Struct Copy

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 2.60 | 0 | 0 | — |
| Saferefl | 357 | 56 | 5 | 137.1× slower |
| Reflect | 58.0 | 0 | 0 | 22.3× slower |
| Copier | 1345 | 432 | 17 | 516.5× slower |

## JSON Like

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibReflect | 196 | 64 | 1 | — |
| Reflect2 | 32.0 | 64 | 1 | 6.1× faster |
| CachedOffset | 80.6 | 64 | 1 | 2.4× faster |
| Generics | 26.4 | 64 | 1 | 7.4× faster |

## Read Int64

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibReflect | 26.1 | 0 | 0 | — |
| Reflect2 | 0.844 | 0 | 0 | 30.9× faster |
| CachedOffset | 5.73 | 0 | 0 | 4.6× faster |
| Generics | 0.283 | 0 | 0 | 92.2× faster |

## Write Int64

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibReflect | 27.3 | 0 | 0 | — |
| Reflect2 | 2.99 | 0 | 0 | 9.1× faster |
| CachedOffset | 10.8 | 0 | 0 | 2.5× faster |
| Generics | 0.269 | 0 | 0 | 101.5× faster |

## Read String

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibReflect | 27.2 | 0 | 0 | — |
| Reflect2 | 0.893 | 0 | 0 | 30.5× faster |
| CachedOffset | 5.23 | 0 | 0 | 5.2× faster |
| Generics | 0.553 | 0 | 0 | 49.2× faster |

