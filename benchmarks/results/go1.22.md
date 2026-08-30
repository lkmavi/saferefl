# Benchmark Results

**Generated:** 2026-08-30 11:29 UTC  
**Go:** 1.22  
**Platform:** linux/amd64  
**CPU:** AMD EPYC 7763 64-Core Processor                  

## Field Read

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Reflect | 68.2 | 8 | 1 | — |
| ReflectFast | 3.12 | 0 | 0 | 21.8× faster |
| SafeRefl | 45.6 | 0 | 0 | 1.5× faster |
| Offset | 11.6 | 0 | 0 | 5.9× faster |
| Accessor | 0.938 | 0 | 0 | 72.7× faster |
| Native | 0.313 | 0 | 0 | 218.0× faster |

## Slice At

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 1.01 | 0 | 0 | — |
| Direct | 1.00 | 0 | 0 | 1.0× faster |
| Reflect | 3.78 | 0 | 0 | 3.7× slower |

## Map Len

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.625 | 0 | 0 | — |
| Builtin | 0.313 | 0 | 0 | 2.0× faster |
| Reflect | 4.68 | 0 | 0 | 7.5× slower |

## Kind Of

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.937 | 0 | 0 | — |
| Reflect | 2.19 | 0 | 0 | 2.3× slower |

## Is Nil_ptr

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 3.48 | 0 | 0 | — |
| Reflect | 2.50 | 0 | 0 | 1.4× faster |

## Each Field

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 51.2 | 0 | 0 | — |
| Reflect | 157 | 56 | 5 | 3.1× slower |
| ReflectFull | 160 | 56 | 5 | 3.1× slower |
| Reflect2 | 24.9 | 0 | 0 | 2.1× faster |

## Copy Fields

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 285 | 0 | 0 | — |
| Manual | 0.312 | 0 | 0 | 913.3× faster |
| Reflect | 592 | 80 | 10 | 2.1× slower |

## To Map

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 226 | 336 | 2 | — |
| Reflect | 523 | 432 | 12 | 2.3× slower |
| JSON | 2782 | 776 | 25 | 12.3× slower |

## Map For Each

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 936 | 0 | 0 | — |
| Range | 979 | 0 | 0 | 1.0× slower |

## Get_int

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 45.7 | 0 | 0 | — |
| Reflect | 68.4 | 8 | 1 | 1.5× slower |
| ReflectFast | 3.12 | 0 | 0 | 14.7× faster |
| Offset | 11.7 | 0 | 0 | 3.9× faster |
| Accessor | 0.935 | 0 | 0 | 48.9× faster |
| AccFrom | 3.12 | 0 | 0 | 14.6× faster |
| Direct | 0.312 | 0 | 0 | 146.4× faster |

## Set_string

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 45.9 | 0 | 0 | — |
| Reflect | 73.9 | 8 | 1 | 1.6× slower |
| ReflectFast | 4.99 | 0 | 0 | 9.2× faster |
| Offset | 13.7 | 0 | 0 | 3.3× faster |
| Accessor | 1.25 | 0 | 0 | 36.8× faster |
| AccOn | 3.43 | 0 | 0 | 13.4× faster |
| Direct | 0.624 | 0 | 0 | 73.6× faster |

## DI Resolve

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 0.628 | 0 | 0 | — |
| SafeRefl | 154 | 0 | 0 | 245.2× slower |
| Accessor | 3.13 | 0 | 0 | 5.0× slower |
| Reflect2 | 21.1 | 0 | 0 | 33.6× slower |
| Reflect | 242 | 24 | 3 | 385.0× slower |

## JSON Decode

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibJSON | 2695 | 392 | 10 | — |
| SafeRefl | 524 | 0 | 0 | 5.1× faster |
| Reflect | 1728 | 160 | 20 | 1.6× faster |
| Reflect2 | 80.2 | 0 | 0 | 33.6× faster |
| Accessor | 9.54 | 0 | 0 | 282.4× faster |

## ORM Scan

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 14.1 | 0 | 0 | — |
| SafeRefl | 539 | 0 | 0 | 38.3× slower |
| Accessor | 14.9 | 0 | 0 | 1.1× slower |
| Reflect2 | 82.8 | 0 | 0 | 5.9× slower |
| Reflect | 982 | 80 | 10 | 69.7× slower |
| Copier | 7265 | 880 | 58 | 515.4× slower |

## Struct Copy

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 9.68 | 0 | 0 | — |
| Accessor | 10.5 | 0 | 0 | 1.1× slower |
| Reflect2 | 41.1 | 0 | 0 | 4.3× slower |
| Offset | 189 | 0 | 0 | 19.6× slower |
| Reflect | 780 | 80 | 10 | 80.6× slower |
| SafeRefl | 521 | 0 | 0 | 53.8× slower |
| Copier | 3260 | 552 | 32 | 336.9× slower |

