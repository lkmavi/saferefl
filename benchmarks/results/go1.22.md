# Benchmark Results

**Generated:** 2026-09-13 11:18 UTC  
**Go:** 1.22  
**Platform:** linux/amd64  
**CPU:** AMD EPYC 9V74 80-Core Processor                  

## Field Read

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Reflect | 57.0 | 8 | 1 | — |
| ReflectFast | 2.48 | 0 | 0 | 23.0× faster |
| SafeRefl | 33.1 | 0 | 0 | 1.7× faster |
| Offset | 9.05 | 0 | 0 | 6.3× faster |
| Accessor | 0.837 | 0 | 0 | 68.1× faster |
| Native | 0.273 | 0 | 0 | 208.9× faster |

## Slice At

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.901 | 0 | 0 | — |
| Direct | 0.921 | 0 | 0 | 1.0× slower |
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
| Reflect | 2.18 | 0 | 0 | 1.1× faster |

## Each Field

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 41.0 | 0 | 0 | — |
| Reflect | 121 | 56 | 5 | 3.0× slower |
| ReflectFull | 124 | 56 | 5 | 3.0× slower |
| Reflect2 | 20.8 | 0 | 0 | 2.0× faster |

## Copy Fields

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 208 | 0 | 0 | — |
| Manual | 0.274 | 0 | 0 | 758.1× faster |
| Reflect | 486 | 80 | 10 | 2.3× slower |

## To Map

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 170 | 336 | 2 | — |
| Reflect | 382 | 432 | 12 | 2.3× slower |
| JSON | 1958 | 776 | 25 | 11.5× slower |

## Map For Each

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 691 | 0 | 0 | — |
| Range | 717 | 0 | 0 | 1.0× slower |

## Get_int

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 32.8 | 0 | 0 | — |
| Reflect | 57.0 | 8 | 1 | 1.7× slower |
| ReflectFast | 2.47 | 0 | 0 | 13.3× faster |
| Offset | 9.37 | 0 | 0 | 3.5× faster |
| Accessor | 0.819 | 0 | 0 | 40.1× faster |
| AccFrom | 2.37 | 0 | 0 | 13.8× faster |
| Direct | 0.273 | 0 | 0 | 120.2× faster |

## Set_string

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 33.4 | 0 | 0 | — |
| Reflect | 60.5 | 8 | 1 | 1.8× slower |
| ReflectFast | 3.83 | 0 | 0 | 8.7× faster |
| Offset | 10.6 | 0 | 0 | 3.1× faster |
| Accessor | 1.09 | 0 | 0 | 30.5× faster |
| AccOn | 2.50 | 0 | 0 | 13.3× faster |
| Direct | 0.551 | 0 | 0 | 60.6× faster |

## DI Resolve

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 0.561 | 0 | 0 | — |
| SafeRefl | 110 | 0 | 0 | 195.3× slower |
| Accessor | 2.73 | 0 | 0 | 4.9× slower |
| Reflect2 | 16.7 | 0 | 0 | 29.8× slower |
| Reflect | 200 | 24 | 3 | 356.2× slower |

## JSON Decode

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibJSON | 1900 | 392 | 10 | — |
| SafeRefl | 380 | 0 | 0 | 5.0× faster |
| Reflect | 1434 | 160 | 20 | 1.3× faster |
| Reflect2 | 66.4 | 0 | 0 | 28.6× faster |
| Accessor | 6.09 | 0 | 0 | 311.9× faster |

## ORM Scan

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 11.1 | 0 | 0 | — |
| SafeRefl | 397 | 0 | 0 | 35.9× slower |
| Accessor | 10.0 | 0 | 0 | 1.1× faster |
| Reflect2 | 65.3 | 0 | 0 | 5.9× slower |
| Reflect | 789 | 80 | 10 | 71.3× slower |
| Copier | 5214 | 880 | 58 | 471.3× slower |

## Struct Copy

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 7.21 | 0 | 0 | — |
| Accessor | 7.82 | 0 | 0 | 1.1× slower |
| Reflect2 | 33.2 | 0 | 0 | 4.6× slower |
| Offset | 146 | 0 | 0 | 20.3× slower |
| Reflect | 661 | 80 | 10 | 91.7× slower |
| SafeRefl | 375 | 0 | 0 | 52.0× slower |
| Copier | 2360 | 552 | 32 | 327.2× slower |

