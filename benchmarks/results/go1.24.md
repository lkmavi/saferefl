# Benchmark Results

**Generated:** 2026-10-11 12:18 UTC  
**Go:** 1.24  
**Platform:** linux/amd64  
**CPU:** AMD EPYC 7763 64-Core Processor                  

## Field Read

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Reflect | 59.1 | 0 | 0 | — |
| ReflectFast | 3.12 | 0 | 0 | 18.9× faster |
| SafeRefl | 54.8 | 0 | 0 | 1.1× faster |
| Offset | 10.3 | 0 | 0 | 5.7× faster |
| Accessor | 0.940 | 0 | 0 | 62.9× faster |
| Native | 0.624 | 0 | 0 | 94.7× faster |

## Slice At

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 1.07 | 0 | 0 | — |
| Direct | 0.999 | 0 | 0 | 1.1× faster |
| Reflect | 3.75 | 0 | 0 | 3.5× slower |

## Map Len

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.625 | 0 | 0 | — |
| Builtin | 0.624 | 0 | 0 | 1.0× faster |
| Reflect | 4.68 | 0 | 0 | 7.5× slower |

## Kind Of

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.936 | 0 | 0 | — |
| Reflect | 1.87 | 0 | 0 | 2.0× slower |

## Is Nil_ptr

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 2.81 | 0 | 0 | — |
| Reflect | 2.81 | 0 | 0 | 1.0× slower |

## Each Field

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 42.8 | 0 | 0 | — |
| Reflect | 139 | 56 | 5 | 3.3× slower |
| ReflectFull | 152 | 56 | 5 | 3.6× slower |
| Reflect2 | 23.3 | 0 | 0 | 1.8× faster |

## Copy Fields

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 249 | 0 | 0 | — |
| Manual | 0.312 | 0 | 0 | 797.4× faster |
| Reflect | 490 | 0 | 0 | 2.0× slower |

## To Map

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 214 | 336 | 2 | — |
| Reflect | 441 | 392 | 7 | 2.1× slower |
| JSON | 2655 | 776 | 25 | 12.4× slower |

## Map For Each

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 835 | 0 | 0 | — |
| Range | 864 | 0 | 0 | 1.0× slower |

## Get_int

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 54.7 | 0 | 0 | — |
| Reflect | 59.0 | 0 | 0 | 1.1× slower |
| ReflectFast | 3.12 | 0 | 0 | 17.5× faster |
| Offset | 10.3 | 0 | 0 | 5.3× faster |
| Accessor | 0.624 | 0 | 0 | 87.7× faster |
| AccFrom | 3.12 | 0 | 0 | 17.6× faster |
| Direct | 0.625 | 0 | 0 | 87.6× faster |

## Set_string

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 46.8 | 0 | 0 | — |
| Reflect | 62.6 | 0 | 0 | 1.3× slower |
| ReflectFast | 4.99 | 0 | 0 | 9.4× faster |
| Offset | 12.8 | 0 | 0 | 3.7× faster |
| Accessor | 1.25 | 0 | 0 | 37.5× faster |
| AccOn | 3.43 | 0 | 0 | 13.6× faster |
| Direct | 0.625 | 0 | 0 | 75.0× faster |

## DI Resolve

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 0.626 | 0 | 0 | — |
| SafeRefl | 143 | 0 | 0 | 228.4× slower |
| Accessor | 2.82 | 0 | 0 | 4.5× slower |
| Reflect2 | 19.6 | 0 | 0 | 31.4× slower |
| Reflect | 225 | 0 | 0 | 358.9× slower |

## JSON Decode

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibJSON | 2789 | 392 | 10 | — |
| SafeRefl | 467 | 0 | 0 | 6.0× faster |
| Reflect | 1525 | 0 | 0 | 1.8× faster |
| Reflect2 | 77.5 | 0 | 0 | 36.0× faster |
| Accessor | 9.79 | 0 | 0 | 284.8× faster |

## ORM Scan

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 14.1 | 0 | 0 | — |
| SafeRefl | 475 | 0 | 0 | 33.7× slower |
| Accessor | 12.1 | 0 | 0 | 1.2× faster |
| Reflect2 | 78.3 | 0 | 0 | 5.6× slower |
| Reflect | 912 | 0 | 0 | 64.7× slower |
| Copier | 6784 | 640 | 28 | 481.2× slower |

## Struct Copy

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 9.68 | 0 | 0 | — |
| Accessor | 8.49 | 0 | 0 | 1.1× faster |
| Reflect2 | 36.2 | 0 | 0 | 3.7× slower |
| Offset | 167 | 0 | 0 | 17.3× slower |
| Reflect | 686 | 0 | 0 | 70.9× slower |
| SafeRefl | 499 | 0 | 0 | 51.5× slower |
| Copier | 2954 | 432 | 17 | 305.2× slower |

