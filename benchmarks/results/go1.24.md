# Benchmark Results

**Generated:** 2026-07-12 08:21 UTC  
**Go:** 1.24  
**Platform:** linux/amd64  
**CPU:** AMD EPYC 9V74 80-Core Processor                  

## Field Read

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Reflect | 49.9 | 0 | 0 | — |
| ReflectFast | 2.48 | 0 | 0 | 20.1× faster |
| SafeRefl | 35.0 | 0 | 0 | 1.4× faster |
| Offset | 8.01 | 0 | 0 | 6.2× faster |
| Accessor | 0.822 | 0 | 0 | 60.7× faster |
| Native | 0.548 | 0 | 0 | 91.0× faster |

## Slice At

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.902 | 0 | 0 | — |
| Direct | 0.921 | 0 | 0 | 1.0× slower |
| Reflect | 3.00 | 0 | 0 | 3.3× slower |

## Map Len

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.551 | 0 | 0 | — |
| Builtin | 0.549 | 0 | 0 | 1.0× faster |
| Reflect | 3.83 | 0 | 0 | 6.9× slower |

## Kind Of

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.822 | 0 | 0 | — |
| Reflect | 1.64 | 0 | 0 | 2.0× slower |

## Is Nil_ptr

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 1.97 | 0 | 0 | — |
| Reflect | 2.47 | 0 | 0 | 1.3× slower |

## Each Field

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 33.6 | 0 | 0 | — |
| Reflect | 105 | 56 | 5 | 3.1× slower |
| ReflectFull | 108 | 56 | 5 | 3.2× slower |
| Reflect2 | 19.4 | 0 | 0 | 1.7× faster |

## Copy Fields

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 191 | 0 | 0 | — |
| Manual | 0.273 | 0 | 0 | 698.0× faster |
| Reflect | 410 | 0 | 0 | 2.2× slower |

## To Map

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 151 | 336 | 2 | — |
| Reflect | 340 | 392 | 7 | 2.3× slower |
| JSON | 1787 | 776 | 25 | 11.8× slower |

## Map For Each

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 593 | 0 | 0 | — |
| Range | 593 | 0 | 0 | 1.0× faster |

## Get_int

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 35.1 | 0 | 0 | — |
| Reflect | 49.8 | 0 | 0 | 1.4× slower |
| ReflectFast | 2.48 | 0 | 0 | 14.2× faster |
| Offset | 8.07 | 0 | 0 | 4.3× faster |
| Accessor | 0.551 | 0 | 0 | 63.7× faster |
| AccFrom | 2.19 | 0 | 0 | 16.0× faster |
| Direct | 0.546 | 0 | 0 | 64.2× faster |

## Set_string

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 36.4 | 0 | 0 | — |
| Reflect | 53.4 | 0 | 0 | 1.5× slower |
| ReflectFast | 3.82 | 0 | 0 | 9.5× faster |
| Offset | 10.1 | 0 | 0 | 3.6× faster |
| Accessor | 1.09 | 0 | 0 | 33.3× faster |
| AccOn | 2.46 | 0 | 0 | 14.8× faster |
| Direct | 0.547 | 0 | 0 | 66.5× faster |

## DI Resolve

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 0.557 | 0 | 0 | — |
| SafeRefl | 110 | 0 | 0 | 197.7× slower |
| Accessor | 2.46 | 0 | 0 | 4.4× slower |
| Reflect2 | 16.0 | 0 | 0 | 28.7× slower |
| Reflect | 183 | 0 | 0 | 329.4× slower |

## JSON Decode

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibJSON | 1825 | 392 | 10 | — |
| SafeRefl | 362 | 0 | 0 | 5.0× faster |
| Reflect | 1273 | 0 | 0 | 1.4× faster |
| Reflect2 | 63.2 | 0 | 0 | 28.9× faster |
| Accessor | 6.28 | 0 | 0 | 290.5× faster |

## ORM Scan

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 11.1 | 0 | 0 | — |
| SafeRefl | 386 | 0 | 0 | 34.8× slower |
| Accessor | 9.96 | 0 | 0 | 1.1× faster |
| Reflect2 | 66.3 | 0 | 0 | 6.0× slower |
| Reflect | 725 | 0 | 0 | 65.3× slower |
| Copier | 4787 | 640 | 28 | 431.0× slower |

## Struct Copy

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 7.25 | 0 | 0 | — |
| Accessor | 7.84 | 0 | 0 | 1.1× slower |
| Reflect2 | 29.6 | 0 | 0 | 4.1× slower |
| Offset | 125 | 0 | 0 | 17.2× slower |
| Reflect | 573 | 0 | 0 | 79.1× slower |
| SafeRefl | 375 | 0 | 0 | 51.7× slower |
| Copier | 2105 | 432 | 17 | 290.2× slower |

