# Benchmark Results

**Generated:** 2026-08-23 06:43 UTC  
**Go:** 1.22  
**Platform:** linux/amd64  
**CPU:** AMD EPYC 7763 64-Core Processor                  

## Field Read

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Reflect | 70.1 | 8 | 1 | — |
| ReflectFast | 3.18 | 0 | 0 | 22.0× faster |
| SafeRefl | 45.5 | 0 | 0 | 1.5× faster |
| Offset | 11.5 | 0 | 0 | 6.1× faster |
| Accessor | 0.938 | 0 | 0 | 74.8× faster |
| Native | 0.315 | 0 | 0 | 222.5× faster |

## Slice At

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 1.01 | 0 | 0 | — |
| Direct | 1.00 | 0 | 0 | 1.0× faster |
| Reflect | 3.74 | 0 | 0 | 3.7× slower |

## Map Len

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.624 | 0 | 0 | — |
| Builtin | 0.313 | 0 | 0 | 2.0× faster |
| Reflect | 4.69 | 0 | 0 | 7.5× slower |

## Kind Of

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.936 | 0 | 0 | — |
| Reflect | 2.18 | 0 | 0 | 2.3× slower |

## Is Nil_ptr

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 3.45 | 0 | 0 | — |
| Reflect | 2.55 | 0 | 0 | 1.4× faster |

## Each Field

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 51.3 | 0 | 0 | — |
| Reflect | 158 | 56 | 5 | 3.1× slower |
| ReflectFull | 163 | 56 | 5 | 3.2× slower |
| Reflect2 | 24.9 | 0 | 0 | 2.1× faster |

## Copy Fields

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 285 | 0 | 0 | — |
| Manual | 0.329 | 0 | 0 | 867.1× faster |
| Reflect | 595 | 80 | 10 | 2.1× slower |

## To Map

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 234 | 336 | 2 | — |
| Reflect | 526 | 432 | 12 | 2.3× slower |
| JSON | 2718 | 776 | 25 | 11.6× slower |

## Map For Each

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 977 | 0 | 0 | — |
| Range | 993 | 0 | 0 | 1.0× slower |

## Get_int

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 45.4 | 0 | 0 | — |
| Reflect | 71.9 | 8 | 1 | 1.6× slower |
| ReflectFast | 3.12 | 0 | 0 | 14.5× faster |
| Offset | 11.7 | 0 | 0 | 3.9× faster |
| Accessor | 0.937 | 0 | 0 | 48.4× faster |
| AccFrom | 3.12 | 0 | 0 | 14.5× faster |
| Direct | 0.312 | 0 | 0 | 145.3× faster |

## Set_string

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 46.0 | 0 | 0 | — |
| Reflect | 72.8 | 8 | 1 | 1.6× slower |
| ReflectFast | 5.00 | 0 | 0 | 9.2× faster |
| Offset | 13.7 | 0 | 0 | 3.3× faster |
| Accessor | 1.25 | 0 | 0 | 36.8× faster |
| AccOn | 3.43 | 0 | 0 | 13.4× faster |
| Direct | 0.625 | 0 | 0 | 73.5× faster |

## DI Resolve

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 0.626 | 0 | 0 | — |
| SafeRefl | 152 | 0 | 0 | 242.8× slower |
| Accessor | 3.12 | 0 | 0 | 5.0× slower |
| Reflect2 | 20.9 | 0 | 0 | 33.4× slower |
| Reflect | 243 | 24 | 3 | 387.5× slower |

## JSON Decode

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibJSON | 2680 | 392 | 10 | — |
| SafeRefl | 525 | 0 | 0 | 5.1× faster |
| Reflect | 1735 | 160 | 20 | 1.5× faster |
| Reflect2 | 80.9 | 0 | 0 | 33.1× faster |
| Accessor | 9.54 | 0 | 0 | 280.9× faster |

## ORM Scan

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 14.1 | 0 | 0 | — |
| SafeRefl | 542 | 0 | 0 | 38.4× slower |
| Accessor | 12.3 | 0 | 0 | 1.1× faster |
| Reflect2 | 83.3 | 0 | 0 | 5.9× slower |
| Reflect | 988 | 80 | 10 | 70.0× slower |
| Copier | 7189 | 880 | 58 | 510.0× slower |

## Struct Copy

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 9.67 | 0 | 0 | — |
| Accessor | 10.5 | 0 | 0 | 1.1× slower |
| Reflect2 | 40.7 | 0 | 0 | 4.2× slower |
| Offset | 184 | 0 | 0 | 19.1× slower |
| Reflect | 785 | 80 | 10 | 81.2× slower |
| SafeRefl | 519 | 0 | 0 | 53.7× slower |
| Copier | 3204 | 552 | 32 | 331.2× slower |

