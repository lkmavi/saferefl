# Benchmark Results

**Generated:** 2026-09-27 11:38 UTC  
**Go:** 1.24  
**Platform:** linux/amd64  
**CPU:** AMD EPYC 7763 64-Core Processor                  

## Field Read

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Reflect | 59.4 | 0 | 0 | — |
| ReflectFast | 3.15 | 0 | 0 | 18.8× faster |
| SafeRefl | 55.1 | 0 | 0 | 1.1× faster |
| Offset | 10.3 | 0 | 0 | 5.7× faster |
| Accessor | 0.944 | 0 | 0 | 62.9× faster |
| Native | 0.625 | 0 | 0 | 94.9× faster |

## Slice At

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 1.08 | 0 | 0 | — |
| Direct | 1.000 | 0 | 0 | 1.1× faster |
| Reflect | 3.75 | 0 | 0 | 3.5× slower |

## Map Len

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.625 | 0 | 0 | — |
| Builtin | 0.625 | 0 | 0 | 1.0× slower |
| Reflect | 4.69 | 0 | 0 | 7.5× slower |

## Kind Of

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.940 | 0 | 0 | — |
| Reflect | 1.87 | 0 | 0 | 2.0× slower |

## Is Nil_ptr

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 2.81 | 0 | 0 | — |
| Reflect | 2.81 | 0 | 0 | 1.0× faster |

## Each Field

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 43.8 | 0 | 0 | — |
| Reflect | 140 | 56 | 5 | 3.2× slower |
| ReflectFull | 152 | 56 | 5 | 3.5× slower |
| Reflect2 | 23.3 | 0 | 0 | 1.9× faster |

## Copy Fields

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 250 | 0 | 0 | — |
| Manual | 0.313 | 0 | 0 | 800.6× faster |
| Reflect | 490 | 0 | 0 | 2.0× slower |

## To Map

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 215 | 336 | 2 | — |
| Reflect | 443 | 392 | 7 | 2.1× slower |
| JSON | 2641 | 776 | 25 | 12.3× slower |

## Map For Each

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 836 | 0 | 0 | — |
| Range | 870 | 0 | 0 | 1.0× slower |

## Get_int

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 54.8 | 0 | 0 | — |
| Reflect | 59.0 | 0 | 0 | 1.1× slower |
| ReflectFast | 3.12 | 0 | 0 | 17.6× faster |
| Offset | 10.3 | 0 | 0 | 5.3× faster |
| Accessor | 0.624 | 0 | 0 | 87.9× faster |
| AccFrom | 3.13 | 0 | 0 | 17.5× faster |
| Direct | 0.625 | 0 | 0 | 87.7× faster |

## Set_string

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 47.5 | 0 | 0 | — |
| Reflect | 62.8 | 0 | 0 | 1.3× slower |
| ReflectFast | 4.99 | 0 | 0 | 9.5× faster |
| Offset | 12.8 | 0 | 0 | 3.7× faster |
| Accessor | 1.25 | 0 | 0 | 38.0× faster |
| AccOn | 3.43 | 0 | 0 | 13.9× faster |
| Direct | 0.630 | 0 | 0 | 75.5× faster |

## DI Resolve

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 0.626 | 0 | 0 | — |
| SafeRefl | 143 | 0 | 0 | 228.6× slower |
| Accessor | 2.89 | 0 | 0 | 4.6× slower |
| Reflect2 | 19.8 | 0 | 0 | 31.6× slower |
| Reflect | 227 | 0 | 0 | 362.2× slower |

## JSON Decode

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibJSON | 2786 | 392 | 10 | — |
| SafeRefl | 466 | 0 | 0 | 6.0× faster |
| Reflect | 1532 | 0 | 0 | 1.8× faster |
| Reflect2 | 85.6 | 0 | 0 | 32.6× faster |
| Accessor | 9.79 | 0 | 0 | 284.5× faster |

## ORM Scan

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 14.2 | 0 | 0 | — |
| SafeRefl | 485 | 0 | 0 | 34.2× slower |
| Accessor | 12.2 | 0 | 0 | 1.2× faster |
| Reflect2 | 78.0 | 0 | 0 | 5.5× slower |
| Reflect | 919 | 0 | 0 | 64.9× slower |
| Copier | 6910 | 640 | 28 | 488.0× slower |

## Struct Copy

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 9.71 | 0 | 0 | — |
| Accessor | 8.52 | 0 | 0 | 1.1× faster |
| Reflect2 | 36.1 | 0 | 0 | 3.7× slower |
| Offset | 167 | 0 | 0 | 17.2× slower |
| Reflect | 689 | 0 | 0 | 71.0× slower |
| SafeRefl | 496 | 0 | 0 | 51.1× slower |
| Copier | 3000 | 432 | 17 | 309.0× slower |

