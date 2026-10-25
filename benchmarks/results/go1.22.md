# Benchmark Results

**Generated:** 2026-10-25 10:56 UTC  
**Go:** 1.22  
**Platform:** linux/amd64  
**CPU:** AMD EPYC 7763 64-Core Processor                  

## Field Read

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Reflect | 69.4 | 8 | 1 | — |
| ReflectFast | 3.15 | 0 | 0 | 22.1× faster |
| SafeRefl | 45.3 | 0 | 0 | 1.5× faster |
| Offset | 11.6 | 0 | 0 | 6.0× faster |
| Accessor | 0.938 | 0 | 0 | 74.0× faster |
| Native | 0.313 | 0 | 0 | 221.8× faster |

## Slice At

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 1.02 | 0 | 0 | — |
| Direct | 1.00 | 0 | 0 | 1.0× faster |
| Reflect | 3.76 | 0 | 0 | 3.7× slower |

## Map Len

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.625 | 0 | 0 | — |
| Builtin | 0.313 | 0 | 0 | 2.0× faster |
| Reflect | 4.68 | 0 | 0 | 7.5× slower |

## Kind Of

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.939 | 0 | 0 | — |
| Reflect | 2.19 | 0 | 0 | 2.3× slower |

## Is Nil_ptr

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 3.45 | 0 | 0 | — |
| Reflect | 2.51 | 0 | 0 | 1.4× faster |

## Each Field

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 51.6 | 0 | 0 | — |
| Reflect | 161 | 56 | 5 | 3.1× slower |
| ReflectFull | 164 | 56 | 5 | 3.2× slower |
| Reflect2 | 25.1 | 0 | 0 | 2.1× faster |

## Copy Fields

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 285 | 0 | 0 | — |
| Manual | 0.312 | 0 | 0 | 912.0× faster |
| Reflect | 601 | 80 | 10 | 2.1× slower |

## To Map

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 233 | 336 | 2 | — |
| Reflect | 536 | 432 | 12 | 2.3× slower |
| JSON | 2853 | 776 | 25 | 12.2× slower |

## Map For Each

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 967 | 0 | 0 | — |
| Range | 988 | 0 | 0 | 1.0× slower |

## Get_int

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 45.1 | 0 | 0 | — |
| Reflect | 68.7 | 8 | 1 | 1.5× slower |
| ReflectFast | 3.13 | 0 | 0 | 14.4× faster |
| Offset | 11.6 | 0 | 0 | 3.9× faster |
| Accessor | 0.937 | 0 | 0 | 48.2× faster |
| AccFrom | 3.12 | 0 | 0 | 14.5× faster |
| Direct | 0.314 | 0 | 0 | 143.9× faster |

## Set_string

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 46.0 | 0 | 0 | — |
| Reflect | 72.4 | 8 | 1 | 1.6× slower |
| ReflectFast | 4.99 | 0 | 0 | 9.2× faster |
| Offset | 13.7 | 0 | 0 | 3.3× faster |
| Accessor | 1.25 | 0 | 0 | 36.7× faster |
| AccOn | 3.44 | 0 | 0 | 13.3× faster |
| Direct | 0.625 | 0 | 0 | 73.5× faster |

## DI Resolve

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 0.637 | 0 | 0 | — |
| SafeRefl | 152 | 0 | 0 | 238.8× slower |
| Accessor | 3.12 | 0 | 0 | 4.9× slower |
| Reflect2 | 20.9 | 0 | 0 | 32.9× slower |
| Reflect | 244 | 24 | 3 | 382.5× slower |

## JSON Decode

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibJSON | 2718 | 392 | 10 | — |
| SafeRefl | 525 | 0 | 0 | 5.2× faster |
| Reflect | 1736 | 160 | 20 | 1.6× faster |
| Reflect2 | 80.2 | 0 | 0 | 33.9× faster |
| Accessor | 9.54 | 0 | 0 | 284.8× faster |

## ORM Scan

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 14.1 | 0 | 0 | — |
| SafeRefl | 541 | 0 | 0 | 38.4× slower |
| Accessor | 12.4 | 0 | 0 | 1.1× faster |
| Reflect2 | 82.8 | 0 | 0 | 5.9× slower |
| Reflect | 986 | 80 | 10 | 69.9× slower |
| Copier | 7243 | 880 | 58 | 513.7× slower |

## Struct Copy

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 9.68 | 0 | 0 | — |
| Accessor | 10.5 | 0 | 0 | 1.1× slower |
| Reflect2 | 40.4 | 0 | 0 | 4.2× slower |
| Offset | 185 | 0 | 0 | 19.2× slower |
| Reflect | 786 | 80 | 10 | 81.1× slower |
| SafeRefl | 519 | 0 | 0 | 53.6× slower |
| Copier | 3252 | 552 | 32 | 335.8× slower |

