# Benchmark Results

**Generated:** 2026-06-23 08:35 UTC  
**Go:** 1.22  
**Platform:** linux/arm64  

## Field Read

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Reflect | 46.1 | 8 | 1 | — |
| ReflectFast | 1.70 | 0 | 0 | 27.1× faster |
| SafeRefl | 22.6 | 0 | 0 | 2.0× faster |
| Offset | 5.24 | 0 | 0 | 8.8× faster |
| Accessor | 0.544 | 0 | 0 | 84.6× faster |
| Native | 0.269 | 0 | 0 | 171.2× faster |

## Slice At

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.544 | 0 | 0 | — |
| Direct | 0.599 | 0 | 0 | 1.1× slower |
| Reflect | 1.90 | 0 | 0 | 3.5× slower |

## Map Len

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.479 | 0 | 0 | — |
| Builtin | 0.279 | 0 | 0 | 1.7× faster |
| Reflect | 2.18 | 0 | 0 | 4.6× slower |

## Get_int

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 22.5 | 0 | 0 | — |
| Reflect | 43.8 | 8 | 1 | 1.9× slower |
| ReflectFast | 1.72 | 0 | 0 | 13.1× faster |
| Offset | 5.13 | 0 | 0 | 4.4× faster |
| Accessor | 0.555 | 0 | 0 | 40.6× faster |
| AccFrom | 1.68 | 0 | 0 | 13.4× faster |
| Direct | 0.277 | 0 | 0 | 81.3× faster |

## Set_string

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 22.6 | 0 | 0 | — |
| Reflect | 45.7 | 8 | 1 | 2.0× slower |
| ReflectFast | 2.92 | 0 | 0 | 7.8× faster |
| Offset | 6.59 | 0 | 0 | 3.4× faster |
| Accessor | 0.834 | 0 | 0 | 27.2× faster |
| AccOn | 2.21 | 0 | 0 | 10.2× faster |
| Direct | 0.276 | 0 | 0 | 81.9× faster |

## DI Resolve

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 0.597 | 0 | 0 | — |
| SafeRefl | 78.6 | 0 | 0 | 131.7× slower |
| Accessor | 1.95 | 0 | 0 | 3.3× slower |
| Reflect2 | 8.83 | 0 | 0 | 14.8× slower |
| Reflect | 153 | 24 | 3 | 257.0× slower |

## JSON Decode

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibJSON | 1473 | 392 | 10 | — |
| SafeRefl | 293 | 0 | 0 | 5.0× faster |
| Reflect | 1046 | 160 | 20 | 1.4× faster |
| Reflect2 | 40.2 | 0 | 0 | 36.6× faster |
| Accessor | 4.43 | 0 | 0 | 332.6× faster |

## ORM Scan

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 9.50 | 0 | 0 | — |
| SafeRefl | 289 | 0 | 0 | 30.4× slower |
| Accessor | 7.00 | 0 | 0 | 1.4× faster |
| Reflect2 | 38.5 | 0 | 0 | 4.0× slower |
| Reflect | 610 | 80 | 10 | 64.2× slower |
| Copier | 3762 | 880 | 58 | 395.8× slower |

## Struct Copy

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 6.12 | 0 | 0 | — |
| Accessor | 5.60 | 0 | 0 | 1.1× faster |
| Reflect2 | 18.7 | 0 | 0 | 3.1× slower |
| Offset | 88.8 | 0 | 0 | 14.5× slower |
| Reflect | 487 | 80 | 10 | 79.6× slower |
| SafeRefl | 273 | 0 | 0 | 44.6× slower |
| Copier | 1621 | 552 | 32 | 264.7× slower |

