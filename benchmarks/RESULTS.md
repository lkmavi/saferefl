# Benchmark Results

**Generated:** 2026-06-21 12:40 UTC  
**Platform:** darwin/arm64  
**CPU:** Apple M3 Max  

## JSON Like

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibReflect | 193 | 64 | 1 | — |
| Reflect2 | 30.2 | 64 | 1 | 6.4× faster |
| CachedOffset | 80.1 | 64 | 1 | 2.4× faster |
| Generics | 26.2 | 64 | 1 | 7.4× faster |

## Read Int64

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibReflect | 25.9 | 0 | 0 | — |
| Reflect2 | 0.848 | 0 | 0 | 30.5× faster |
| CachedOffset | 5.71 | 0 | 0 | 4.5× faster |
| Generics | 0.286 | 0 | 0 | 90.4× faster |

## Write Int64

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibReflect | 27.1 | 0 | 0 | — |
| Reflect2 | 3.00 | 0 | 0 | 9.0× faster |
| CachedOffset | 10.9 | 0 | 0 | 2.5× faster |
| Generics | 0.264 | 0 | 0 | 102.6× faster |

## Read String

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibReflect | 27.1 | 0 | 0 | — |
| Reflect2 | 0.892 | 0 | 0 | 30.4× faster |
| CachedOffset | 5.25 | 0 | 0 | 5.2× faster |
| Generics | 0.546 | 0 | 0 | 49.7× faster |

