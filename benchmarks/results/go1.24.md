# Benchmark Results

**Generated:** 2026-10-18 10:34 UTC  
**Go:** 1.24  
**Platform:** linux/amd64  
**CPU:** AMD EPYC 7763 64-Core Processor                  

## Field Read

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Reflect | 59.2 | 0 | 0 | — |
| ReflectFast | 3.13 | 0 | 0 | 18.9× faster |
| SafeRefl | 54.8 | 0 | 0 | 1.1× faster |
| Offset | 10.3 | 0 | 0 | 5.7× faster |
| Accessor | 0.938 | 0 | 0 | 63.1× faster |
| Native | 0.625 | 0 | 0 | 94.7× faster |

## Slice At

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 1.08 | 0 | 0 | — |
| Direct | 0.999 | 0 | 0 | 1.1× faster |
| Reflect | 3.75 | 0 | 0 | 3.5× slower |

## Map Len

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.629 | 0 | 0 | — |
| Builtin | 0.625 | 0 | 0 | 1.0× faster |
| Reflect | 4.69 | 0 | 0 | 7.5× slower |

## Kind Of

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.936 | 0 | 0 | — |
| Reflect | 1.87 | 0 | 0 | 2.0× slower |

## Is Nil_ptr

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 2.81 | 0 | 0 | — |
| Reflect | 2.81 | 0 | 0 | 1.0× faster |

## Each Field

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 43.0 | 0 | 0 | — |
| Reflect | 140 | 56 | 5 | 3.2× slower |
| ReflectFull | 152 | 56 | 5 | 3.5× slower |
| Reflect2 | 23.3 | 0 | 0 | 1.8× faster |

## Copy Fields

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 249 | 0 | 0 | — |
| Manual | 0.312 | 0 | 0 | 798.0× faster |
| Reflect | 504 | 0 | 0 | 2.0× slower |

## To Map

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 214 | 336 | 2 | — |
| Reflect | 441 | 392 | 7 | 2.1× slower |
| JSON | 2602 | 776 | 25 | 12.2× slower |

## Map For Each

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 850 | 0 | 0 | — |
| Range | 879 | 0 | 0 | 1.0× slower |

## Get_int

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 54.8 | 0 | 0 | — |
| Reflect | 59.0 | 0 | 0 | 1.1× slower |
| ReflectFast | 3.13 | 0 | 0 | 17.5× faster |
| Offset | 10.3 | 0 | 0 | 5.3× faster |
| Accessor | 0.624 | 0 | 0 | 87.8× faster |
| AccFrom | 3.12 | 0 | 0 | 17.6× faster |
| Direct | 0.625 | 0 | 0 | 87.6× faster |

## Set_string

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 46.9 | 0 | 0 | — |
| Reflect | 62.7 | 0 | 0 | 1.3× slower |
| ReflectFast | 4.99 | 0 | 0 | 9.4× faster |
| Offset | 12.8 | 0 | 0 | 3.7× faster |
| Accessor | 1.25 | 0 | 0 | 37.6× faster |
| AccOn | 3.43 | 0 | 0 | 13.7× faster |
| Direct | 0.625 | 0 | 0 | 75.1× faster |

## DI Resolve

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 0.629 | 0 | 0 | — |
| SafeRefl | 143 | 0 | 0 | 227.3× slower |
| Accessor | 2.81 | 0 | 0 | 4.5× slower |
| Reflect2 | 20.1 | 0 | 0 | 32.0× slower |
| Reflect | 226 | 0 | 0 | 359.2× slower |

## JSON Decode

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibJSON | 2788 | 392 | 10 | — |
| SafeRefl | 468 | 0 | 0 | 6.0× faster |
| Reflect | 1527 | 0 | 0 | 1.8× faster |
| Reflect2 | 86.3 | 0 | 0 | 32.3× faster |
| Accessor | 9.80 | 0 | 0 | 284.3× faster |

## ORM Scan

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 14.1 | 0 | 0 | — |
| SafeRefl | 473 | 0 | 0 | 33.5× slower |
| Accessor | 12.2 | 0 | 0 | 1.2× faster |
| Reflect2 | 78.1 | 0 | 0 | 5.5× slower |
| Reflect | 918 | 0 | 0 | 65.0× slower |
| Copier | 6831 | 640 | 28 | 483.3× slower |

## Struct Copy

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 9.70 | 0 | 0 | — |
| Accessor | 8.47 | 0 | 0 | 1.1× faster |
| Reflect2 | 37.0 | 0 | 0 | 3.8× slower |
| Offset | 167 | 0 | 0 | 17.3× slower |
| Reflect | 689 | 0 | 0 | 71.1× slower |
| SafeRefl | 496 | 0 | 0 | 51.2× slower |
| Copier | 2974 | 432 | 17 | 306.7× slower |

