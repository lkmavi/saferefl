# Benchmark Results

**Generated:** 2026-08-09 06:59 UTC  
**Go:** 1.22  
**Platform:** linux/amd64  
**CPU:** AMD EPYC 7763 64-Core Processor                  

## Field Read

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Reflect | 68.6 | 8 | 1 | — |
| ReflectFast | 3.32 | 0 | 0 | 20.6× faster |
| SafeRefl | 45.3 | 0 | 0 | 1.5× faster |
| Offset | 11.6 | 0 | 0 | 5.9× faster |
| Accessor | 0.940 | 0 | 0 | 72.9× faster |
| Native | 0.313 | 0 | 0 | 219.0× faster |

## Slice At

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 1.01 | 0 | 0 | — |
| Direct | 1.000 | 0 | 0 | 1.0× faster |
| Reflect | 3.75 | 0 | 0 | 3.7× slower |

## Map Len

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.625 | 0 | 0 | — |
| Builtin | 0.312 | 0 | 0 | 2.0× faster |
| Reflect | 4.69 | 0 | 0 | 7.5× slower |

## Kind Of

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.935 | 0 | 0 | — |
| Reflect | 2.19 | 0 | 0 | 2.3× slower |

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
| Manual | 0.312 | 0 | 0 | 916.2× faster |
| Reflect | 592 | 80 | 10 | 2.1× slower |

## To Map

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 226 | 336 | 2 | — |
| Reflect | 516 | 432 | 12 | 2.3× slower |
| JSON | 2770 | 776 | 25 | 12.3× slower |

## Map For Each

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 965 | 0 | 0 | — |
| Range | 982 | 0 | 0 | 1.0× slower |

## Get_int

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 45.3 | 0 | 0 | — |
| Reflect | 67.6 | 8 | 1 | 1.5× slower |
| ReflectFast | 3.12 | 0 | 0 | 14.5× faster |
| Offset | 11.5 | 0 | 0 | 3.9× faster |
| Accessor | 0.936 | 0 | 0 | 48.4× faster |
| AccFrom | 3.12 | 0 | 0 | 14.5× faster |
| Direct | 0.312 | 0 | 0 | 144.9× faster |

## Set_string

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 45.9 | 0 | 0 | — |
| Reflect | 71.6 | 8 | 1 | 1.6× slower |
| ReflectFast | 5.00 | 0 | 0 | 9.2× faster |
| Offset | 13.7 | 0 | 0 | 3.3× faster |
| Accessor | 1.25 | 0 | 0 | 36.8× faster |
| AccOn | 3.43 | 0 | 0 | 13.4× faster |
| Direct | 0.626 | 0 | 0 | 73.3× faster |

## DI Resolve

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 0.627 | 0 | 0 | — |
| SafeRefl | 153 | 0 | 0 | 244.2× slower |
| Accessor | 3.14 | 0 | 0 | 5.0× slower |
| Reflect2 | 21.3 | 0 | 0 | 33.9× slower |
| Reflect | 242 | 24 | 3 | 386.0× slower |

## JSON Decode

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibJSON | 2696 | 392 | 10 | — |
| SafeRefl | 522 | 0 | 0 | 5.2× faster |
| Reflect | 1733 | 160 | 20 | 1.6× faster |
| Reflect2 | 82.2 | 0 | 0 | 32.8× faster |
| Accessor | 9.56 | 0 | 0 | 282.1× faster |

## ORM Scan

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 14.1 | 0 | 0 | — |
| SafeRefl | 541 | 0 | 0 | 38.4× slower |
| Accessor | 12.4 | 0 | 0 | 1.1× faster |
| Reflect2 | 84.5 | 0 | 0 | 6.0× slower |
| Reflect | 986 | 80 | 10 | 70.0× slower |
| Copier | 7311 | 880 | 58 | 518.7× slower |

## Struct Copy

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 9.68 | 0 | 0 | — |
| Accessor | 10.5 | 0 | 0 | 1.1× slower |
| Reflect2 | 41.9 | 0 | 0 | 4.3× slower |
| Offset | 185 | 0 | 0 | 19.1× slower |
| Reflect | 781 | 80 | 10 | 80.7× slower |
| SafeRefl | 520 | 0 | 0 | 53.7× slower |
| Copier | 3247 | 552 | 32 | 335.4× slower |

