# Benchmark Results

**Generated:** 2026-07-26 08:29 UTC  
**Go:** 1.22  
**Platform:** linux/amd64  
**CPU:** Intel(R) Xeon(R) Platinum 8370C CPU @ 2.80GHz  

## Field Read

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Reflect | 58.3 | 8 | 1 | — |
| ReflectFast | 3.05 | 0 | 0 | 19.1× faster |
| SafeRefl | 37.2 | 0 | 0 | 1.6× faster |
| Offset | 9.92 | 0 | 0 | 5.9× faster |
| Accessor | 0.881 | 0 | 0 | 66.1× faster |
| Native | 0.529 | 0 | 0 | 110.2× faster |

## Slice At

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 1.16 | 0 | 0 | — |
| Direct | 1.16 | 0 | 0 | 1.0× faster |
| Reflect | 3.68 | 0 | 0 | 3.2× slower |

## Map Len

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.585 | 0 | 0 | — |
| Builtin | 0.443 | 0 | 0 | 1.3× faster |
| Reflect | 3.20 | 0 | 0 | 5.5× slower |

## Kind Of

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.752 | 0 | 0 | — |
| Reflect | 1.74 | 0 | 0 | 2.3× slower |

## Is Nil_ptr

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 2.68 | 0 | 0 | — |
| Reflect | 2.04 | 0 | 0 | 1.3× faster |

## Each Field

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 41.4 | 0 | 0 | — |
| Reflect | 148 | 56 | 5 | 3.6× slower |
| ReflectFull | 150 | 56 | 5 | 3.6× slower |
| Reflect2 | 21.5 | 0 | 0 | 1.9× faster |

## Copy Fields

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 234 | 0 | 0 | — |
| Manual | 0.289 | 0 | 0 | 810.2× faster |
| Reflect | 574 | 80 | 10 | 2.5× slower |

## To Map

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 234 | 336 | 2 | — |
| Reflect | 541 | 432 | 12 | 2.3× slower |
| JSON | 2692 | 776 | 25 | 11.5× slower |

## Map For Each

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 827 | 0 | 0 | — |
| Range | 859 | 0 | 0 | 1.0× slower |

## Get_int

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 37.1 | 0 | 0 | — |
| Reflect | 59.7 | 8 | 1 | 1.6× slower |
| ReflectFast | 3.05 | 0 | 0 | 12.2× faster |
| Offset | 9.95 | 0 | 0 | 3.7× faster |
| Accessor | 0.869 | 0 | 0 | 42.7× faster |
| AccFrom | 2.62 | 0 | 0 | 14.2× faster |
| Direct | 0.527 | 0 | 0 | 70.3× faster |

## Set_string

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 40.5 | 0 | 0 | — |
| Reflect | 66.2 | 8 | 1 | 1.6× slower |
| ReflectFast | 4.17 | 0 | 0 | 9.7× faster |
| Offset | 11.0 | 0 | 0 | 3.7× faster |
| Accessor | 1.16 | 0 | 0 | 35.0× faster |
| AccOn | 2.85 | 0 | 0 | 14.2× faster |
| Direct | 0.596 | 0 | 0 | 67.8× faster |

## DI Resolve

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 0.583 | 0 | 0 | — |
| SafeRefl | 124 | 0 | 0 | 211.9× slower |
| Accessor | 2.93 | 0 | 0 | 5.0× slower |
| Reflect2 | 14.5 | 0 | 0 | 24.9× slower |
| Reflect | 234 | 24 | 3 | 401.5× slower |

## JSON Decode

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibJSON | 2553 | 392 | 10 | — |
| SafeRefl | 476 | 0 | 0 | 5.4× faster |
| Reflect | 1609 | 160 | 20 | 1.6× faster |
| Reflect2 | 59.3 | 0 | 0 | 43.0× faster |
| Accessor | 6.37 | 0 | 0 | 400.7× faster |

## ORM Scan

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 11.6 | 0 | 0 | — |
| SafeRefl | 480 | 0 | 0 | 41.5× slower |
| Accessor | 11.5 | 0 | 0 | 1.0× faster |
| Reflect2 | 61.3 | 0 | 0 | 5.3× slower |
| Reflect | 946 | 80 | 10 | 81.7× slower |
| Copier | 7076 | 880 | 58 | 610.7× slower |

## Struct Copy

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 7.24 | 0 | 0 | — |
| Accessor | 7.37 | 0 | 0 | 1.0× slower |
| Reflect2 | 29.8 | 0 | 0 | 4.1× slower |
| Offset | 148 | 0 | 0 | 20.4× slower |
| Reflect | 699 | 80 | 10 | 96.6× slower |
| SafeRefl | 473 | 0 | 0 | 65.4× slower |
| Copier | 3202 | 552 | 32 | 442.3× slower |

