# Benchmark Results

**Generated:** 2026-07-19 08:22 UTC  
**Go:** 1.22  
**Platform:** linux/amd64  
**CPU:** AMD EPYC 7763 64-Core Processor                  

## Field Read

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Reflect | 69.0 | 8 | 1 | — |
| ReflectFast | 3.24 | 0 | 0 | 21.3× faster |
| SafeRefl | 45.3 | 0 | 0 | 1.5× faster |
| Offset | 11.6 | 0 | 0 | 5.9× faster |
| Accessor | 0.937 | 0 | 0 | 73.6× faster |
| Native | 0.312 | 0 | 0 | 220.8× faster |

## Slice At

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 1.01 | 0 | 0 | — |
| Direct | 0.998 | 0 | 0 | 1.0× faster |
| Reflect | 3.75 | 0 | 0 | 3.7× slower |

## Map Len

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.624 | 0 | 0 | — |
| Builtin | 0.313 | 0 | 0 | 2.0× faster |
| Reflect | 4.68 | 0 | 0 | 7.5× slower |

## Kind Of

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.935 | 0 | 0 | — |
| Reflect | 2.18 | 0 | 0 | 2.3× slower |

## Is Nil_ptr

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 3.43 | 0 | 0 | — |
| Reflect | 2.50 | 0 | 0 | 1.4× faster |

## Each Field

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 51.3 | 0 | 0 | — |
| Reflect | 157 | 56 | 5 | 3.1× slower |
| ReflectFull | 160 | 56 | 5 | 3.1× slower |
| Reflect2 | 24.9 | 0 | 0 | 2.1× faster |

## Copy Fields

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 286 | 0 | 0 | — |
| Manual | 0.312 | 0 | 0 | 916.5× faster |
| Reflect | 590 | 80 | 10 | 2.1× slower |

## To Map

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 226 | 336 | 2 | — |
| Reflect | 516 | 432 | 12 | 2.3× slower |
| JSON | 2783 | 776 | 25 | 12.3× slower |

## Map For Each

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 950 | 0 | 0 | — |
| Range | 973 | 0 | 0 | 1.0× slower |

## Get_int

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 45.3 | 0 | 0 | — |
| Reflect | 67.7 | 8 | 1 | 1.5× slower |
| ReflectFast | 3.12 | 0 | 0 | 14.5× faster |
| Offset | 11.6 | 0 | 0 | 3.9× faster |
| Accessor | 0.942 | 0 | 0 | 48.1× faster |
| AccFrom | 3.12 | 0 | 0 | 14.5× faster |
| Direct | 0.312 | 0 | 0 | 145.2× faster |

## Set_string

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 46.1 | 0 | 0 | — |
| Reflect | 71.5 | 8 | 1 | 1.5× slower |
| ReflectFast | 5.02 | 0 | 0 | 9.2× faster |
| Offset | 13.7 | 0 | 0 | 3.4× faster |
| Accessor | 1.25 | 0 | 0 | 37.0× faster |
| AccOn | 3.43 | 0 | 0 | 13.5× faster |
| Direct | 0.625 | 0 | 0 | 73.9× faster |

## DI Resolve

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 0.634 | 0 | 0 | — |
| SafeRefl | 153 | 0 | 0 | 241.6× slower |
| Accessor | 3.14 | 0 | 0 | 4.9× slower |
| Reflect2 | 21.0 | 0 | 0 | 33.1× slower |
| Reflect | 242 | 24 | 3 | 381.8× slower |

## JSON Decode

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibJSON | 2699 | 392 | 10 | — |
| SafeRefl | 525 | 0 | 0 | 5.1× faster |
| Reflect | 1731 | 160 | 20 | 1.6× faster |
| Reflect2 | 80.1 | 0 | 0 | 33.7× faster |
| Accessor | 9.55 | 0 | 0 | 282.7× faster |

## ORM Scan

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 14.1 | 0 | 0 | — |
| SafeRefl | 541 | 0 | 0 | 38.4× slower |
| Accessor | 12.4 | 0 | 0 | 1.1× faster |
| Reflect2 | 83.0 | 0 | 0 | 5.9× slower |
| Reflect | 983 | 80 | 10 | 69.7× slower |
| Copier | 7263 | 880 | 58 | 515.1× slower |

## Struct Copy

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 9.67 | 0 | 0 | — |
| Accessor | 10.5 | 0 | 0 | 1.1× slower |
| Reflect2 | 40.3 | 0 | 0 | 4.2× slower |
| Offset | 185 | 0 | 0 | 19.1× slower |
| Reflect | 781 | 80 | 10 | 80.8× slower |
| SafeRefl | 522 | 0 | 0 | 54.0× slower |
| Copier | 3240 | 552 | 32 | 335.1× slower |

