# Benchmark Results

**Generated:** 2026-11-01 12:51 UTC  
**Go:** 1.24  
**Platform:** linux/amd64  
**CPU:** AMD EPYC 7763 64-Core Processor                  

## Field Read

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Reflect | 59.2 | 0 | 0 | — |
| ReflectFast | 3.13 | 0 | 0 | 18.9× faster |
| SafeRefl | 55.0 | 0 | 0 | 1.1× faster |
| Offset | 10.3 | 0 | 0 | 5.7× faster |
| Accessor | 0.939 | 0 | 0 | 63.0× faster |
| Native | 0.628 | 0 | 0 | 94.2× faster |

## Slice At

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 1.08 | 0 | 0 | — |
| Direct | 1.00 | 0 | 0 | 1.1× faster |
| Reflect | 3.75 | 0 | 0 | 3.5× slower |

## Map Len

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.629 | 0 | 0 | — |
| Builtin | 0.626 | 0 | 0 | 1.0× faster |
| Reflect | 4.69 | 0 | 0 | 7.5× slower |

## Kind Of

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.938 | 0 | 0 | — |
| Reflect | 1.88 | 0 | 0 | 2.0× slower |

## Is Nil_ptr

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 2.82 | 0 | 0 | — |
| Reflect | 2.81 | 0 | 0 | 1.0× faster |

## Each Field

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 43.0 | 0 | 0 | — |
| Reflect | 141 | 56 | 5 | 3.3× slower |
| ReflectFull | 153 | 56 | 5 | 3.6× slower |
| Reflect2 | 23.4 | 0 | 0 | 1.8× faster |

## Copy Fields

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 249 | 0 | 0 | — |
| Manual | 0.312 | 0 | 0 | 795.7× faster |
| Reflect | 490 | 0 | 0 | 2.0× slower |

## To Map

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 213 | 336 | 2 | — |
| Reflect | 446 | 392 | 7 | 2.1× slower |
| JSON | 2617 | 776 | 25 | 12.3× slower |

## Map For Each

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 833 | 0 | 0 | — |
| Range | 866 | 0 | 0 | 1.0× slower |

## Get_int

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 54.8 | 0 | 0 | — |
| Reflect | 59.1 | 0 | 0 | 1.1× slower |
| ReflectFast | 3.14 | 0 | 0 | 17.4× faster |
| Offset | 10.3 | 0 | 0 | 5.3× faster |
| Accessor | 0.628 | 0 | 0 | 87.2× faster |
| AccFrom | 3.12 | 0 | 0 | 17.5× faster |
| Direct | 0.628 | 0 | 0 | 87.2× faster |

## Set_string

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 46.9 | 0 | 0 | — |
| Reflect | 62.7 | 0 | 0 | 1.3× slower |
| ReflectFast | 5.00 | 0 | 0 | 9.4× faster |
| Offset | 12.8 | 0 | 0 | 3.7× faster |
| Accessor | 1.25 | 0 | 0 | 37.6× faster |
| AccOn | 3.46 | 0 | 0 | 13.6× faster |
| Direct | 0.625 | 0 | 0 | 75.1× faster |

## DI Resolve

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 0.637 | 0 | 0 | — |
| SafeRefl | 143 | 0 | 0 | 224.3× slower |
| Accessor | 2.87 | 0 | 0 | 4.5× slower |
| Reflect2 | 20.2 | 0 | 0 | 31.6× slower |
| Reflect | 227 | 0 | 0 | 356.3× slower |

## JSON Decode

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibJSON | 2857 | 392 | 10 | — |
| SafeRefl | 469 | 0 | 0 | 6.1× faster |
| Reflect | 1528 | 0 | 0 | 1.9× faster |
| Reflect2 | 85.4 | 0 | 0 | 33.5× faster |
| Accessor | 9.82 | 0 | 0 | 290.8× faster |

## ORM Scan

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 14.1 | 0 | 0 | — |
| SafeRefl | 474 | 0 | 0 | 33.6× slower |
| Accessor | 12.1 | 0 | 0 | 1.2× faster |
| Reflect2 | 78.0 | 0 | 0 | 5.5× slower |
| Reflect | 909 | 0 | 0 | 64.3× slower |
| Copier | 6925 | 640 | 28 | 490.3× slower |

## Struct Copy

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 9.71 | 0 | 0 | — |
| Accessor | 8.49 | 0 | 0 | 1.1× faster |
| Reflect2 | 37.2 | 0 | 0 | 3.8× slower |
| Offset | 167 | 0 | 0 | 17.2× slower |
| Reflect | 686 | 0 | 0 | 70.7× slower |
| SafeRefl | 510 | 0 | 0 | 52.6× slower |
| Copier | 3016 | 432 | 17 | 310.7× slower |

