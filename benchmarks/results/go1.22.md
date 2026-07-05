# Benchmark Results

**Generated:** 2026-07-05 08:55 UTC  
**Go:** 1.22  
**Platform:** linux/amd64  
**CPU:** AMD EPYC 7763 64-Core Processor                  

## Field Read

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Reflect | 68.6 | 8 | 1 | — |
| ReflectFast | 3.12 | 0 | 0 | 22.0× faster |
| SafeRefl | 45.2 | 0 | 0 | 1.5× faster |
| Offset | 11.5 | 0 | 0 | 5.9× faster |
| Accessor | 0.940 | 0 | 0 | 72.9× faster |
| Native | 0.314 | 0 | 0 | 218.7× faster |

## Slice At

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 1.01 | 0 | 0 | — |
| Direct | 0.998 | 0 | 0 | 1.0× faster |
| Reflect | 3.75 | 0 | 0 | 3.7× slower |

## Map Len

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.625 | 0 | 0 | — |
| Builtin | 0.314 | 0 | 0 | 2.0× faster |
| Reflect | 4.68 | 0 | 0 | 7.5× slower |

## Kind Of

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.936 | 0 | 0 | — |
| Reflect | 2.18 | 0 | 0 | 2.3× slower |

## Is Nil_ptr

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 3.43 | 0 | 0 | — |
| Reflect | 2.50 | 0 | 0 | 1.4× faster |

## Each Field

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 51.5 | 0 | 0 | — |
| Reflect | 157 | 56 | 5 | 3.0× slower |
| ReflectFull | 161 | 56 | 5 | 3.1× slower |
| Reflect2 | 24.9 | 0 | 0 | 2.1× faster |

## Copy Fields

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 286 | 0 | 0 | — |
| Manual | 0.312 | 0 | 0 | 916.2× faster |
| Reflect | 595 | 80 | 10 | 2.1× slower |

## To Map

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 228 | 336 | 2 | — |
| Reflect | 522 | 432 | 12 | 2.3× slower |
| JSON | 2754 | 776 | 25 | 12.1× slower |

## Map For Each

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 952 | 0 | 0 | — |
| Range | 973 | 0 | 0 | 1.0× slower |

## Get_int

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 44.8 | 0 | 0 | — |
| Reflect | 68.0 | 8 | 1 | 1.5× slower |
| ReflectFast | 3.12 | 0 | 0 | 14.4× faster |
| Offset | 11.7 | 0 | 0 | 3.8× faster |
| Accessor | 0.936 | 0 | 0 | 47.8× faster |
| AccFrom | 3.15 | 0 | 0 | 14.2× faster |
| Direct | 0.312 | 0 | 0 | 143.5× faster |

## Set_string

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 45.6 | 0 | 0 | — |
| Reflect | 72.6 | 8 | 1 | 1.6× slower |
| ReflectFast | 5.00 | 0 | 0 | 9.1× faster |
| Offset | 13.7 | 0 | 0 | 3.3× faster |
| Accessor | 1.25 | 0 | 0 | 36.5× faster |
| AccOn | 4.77 | 0 | 0 | 9.6× faster |
| Direct | 0.624 | 0 | 0 | 73.1× faster |

## DI Resolve

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 0.627 | 0 | 0 | — |
| SafeRefl | 152 | 0 | 0 | 242.1× slower |
| Accessor | 3.12 | 0 | 0 | 5.0× slower |
| Reflect2 | 20.9 | 0 | 0 | 33.4× slower |
| Reflect | 242 | 24 | 3 | 386.6× slower |

## JSON Decode

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibJSON | 2676 | 392 | 10 | — |
| SafeRefl | 520 | 0 | 0 | 5.1× faster |
| Reflect | 1727 | 160 | 20 | 1.5× faster |
| Reflect2 | 80.8 | 0 | 0 | 33.1× faster |
| Accessor | 9.54 | 0 | 0 | 280.5× faster |

## ORM Scan

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 14.1 | 0 | 0 | — |
| SafeRefl | 543 | 0 | 0 | 38.5× slower |
| Accessor | 12.3 | 0 | 0 | 1.1× faster |
| Reflect2 | 82.9 | 0 | 0 | 5.9× slower |
| Reflect | 982 | 80 | 10 | 69.6× slower |
| Copier | 7131 | 880 | 58 | 505.9× slower |

## Struct Copy

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 9.67 | 0 | 0 | — |
| Accessor | 10.5 | 0 | 0 | 1.1× slower |
| Reflect2 | 40.6 | 0 | 0 | 4.2× slower |
| Offset | 185 | 0 | 0 | 19.1× slower |
| Reflect | 783 | 80 | 10 | 81.0× slower |
| SafeRefl | 520 | 0 | 0 | 53.8× slower |
| Copier | 3160 | 552 | 32 | 326.6× slower |

