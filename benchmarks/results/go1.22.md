# Benchmark Results

**Generated:** 2026-08-16 06:41 UTC  
**Go:** 1.22  
**Platform:** linux/amd64  
**CPU:** AMD EPYC 9V74 80-Core Processor                  

## Field Read

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Reflect | 57.6 | 8 | 1 | — |
| ReflectFast | 2.48 | 0 | 0 | 23.2× faster |
| SafeRefl | 32.9 | 0 | 0 | 1.8× faster |
| Offset | 9.02 | 0 | 0 | 6.4× faster |
| Accessor | 0.877 | 0 | 0 | 65.7× faster |
| Native | 0.273 | 0 | 0 | 211.1× faster |

## Slice At

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.905 | 0 | 0 | — |
| Direct | 0.930 | 0 | 0 | 1.0× slower |
| Reflect | 3.01 | 0 | 0 | 3.3× slower |

## Map Len

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.547 | 0 | 0 | — |
| Builtin | 0.274 | 0 | 0 | 2.0× faster |
| Reflect | 3.55 | 0 | 0 | 6.5× slower |

## Kind Of

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.819 | 0 | 0 | — |
| Reflect | 1.64 | 0 | 0 | 2.0× slower |

## Is Nil_ptr

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 2.46 | 0 | 0 | — |
| Reflect | 2.19 | 0 | 0 | 1.1× faster |

## Each Field

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 40.8 | 0 | 0 | — |
| Reflect | 121 | 56 | 5 | 3.0× slower |
| ReflectFull | 124 | 56 | 5 | 3.0× slower |
| Reflect2 | 20.8 | 0 | 0 | 2.0× faster |

## Copy Fields

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 205 | 0 | 0 | — |
| Manual | 0.273 | 0 | 0 | 750.2× faster |
| Reflect | 499 | 80 | 10 | 2.4× slower |

## To Map

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 171 | 336 | 2 | — |
| Reflect | 384 | 432 | 12 | 2.2× slower |
| JSON | 1921 | 776 | 25 | 11.2× slower |

## Map For Each

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 692 | 0 | 0 | — |
| Range | 723 | 0 | 0 | 1.0× slower |

## Get_int

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 32.8 | 0 | 0 | — |
| Reflect | 57.1 | 8 | 1 | 1.7× slower |
| ReflectFast | 2.48 | 0 | 0 | 13.2× faster |
| Offset | 9.24 | 0 | 0 | 3.5× faster |
| Accessor | 0.821 | 0 | 0 | 39.9× faster |
| AccFrom | 2.19 | 0 | 0 | 15.0× faster |
| Direct | 0.273 | 0 | 0 | 119.9× faster |

## Set_string

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 33.9 | 0 | 0 | — |
| Reflect | 60.5 | 8 | 1 | 1.8× slower |
| ReflectFast | 3.83 | 0 | 0 | 8.9× faster |
| Offset | 10.7 | 0 | 0 | 3.2× faster |
| Accessor | 1.10 | 0 | 0 | 30.8× faster |
| AccOn | 2.46 | 0 | 0 | 13.8× faster |
| Direct | 0.546 | 0 | 0 | 62.0× faster |

## DI Resolve

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 0.570 | 0 | 0 | — |
| SafeRefl | 110 | 0 | 0 | 193.2× slower |
| Accessor | 2.73 | 0 | 0 | 4.8× slower |
| Reflect2 | 16.7 | 0 | 0 | 29.3× slower |
| Reflect | 200 | 24 | 3 | 350.2× slower |

## JSON Decode

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibJSON | 1902 | 392 | 10 | — |
| SafeRefl | 384 | 0 | 0 | 5.0× faster |
| Reflect | 1429 | 160 | 20 | 1.3× faster |
| Reflect2 | 67.3 | 0 | 0 | 28.3× faster |
| Accessor | 6.08 | 0 | 0 | 312.9× faster |

## ORM Scan

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 11.0 | 0 | 0 | — |
| SafeRefl | 386 | 0 | 0 | 34.9× slower |
| Accessor | 10.0 | 0 | 0 | 1.1× faster |
| Reflect2 | 65.8 | 0 | 0 | 6.0× slower |
| Reflect | 789 | 80 | 10 | 71.4× slower |
| Copier | 5223 | 880 | 58 | 472.8× slower |

## Struct Copy

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 7.21 | 0 | 0 | — |
| Accessor | 7.83 | 0 | 0 | 1.1× slower |
| Reflect2 | 33.4 | 0 | 0 | 4.6× slower |
| Offset | 145 | 0 | 0 | 20.1× slower |
| Reflect | 655 | 80 | 10 | 90.9× slower |
| SafeRefl | 372 | 0 | 0 | 51.6× slower |
| Copier | 2387 | 552 | 32 | 331.0× slower |

