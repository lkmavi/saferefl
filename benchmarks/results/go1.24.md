# Benchmark Results

**Generated:** 2026-09-20 10:52 UTC  
**Go:** 1.24  
**Platform:** linux/amd64  
**CPU:** AMD EPYC 9V74 80-Core Processor                  

## Field Read

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Reflect | 49.8 | 0 | 0 | — |
| ReflectFast | 2.50 | 0 | 0 | 19.9× faster |
| SafeRefl | 35.3 | 0 | 0 | 1.4× faster |
| Offset | 7.99 | 0 | 0 | 6.2× faster |
| Accessor | 0.820 | 0 | 0 | 60.7× faster |
| Native | 0.547 | 0 | 0 | 91.1× faster |

## Slice At

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.902 | 0 | 0 | — |
| Direct | 0.920 | 0 | 0 | 1.0× slower |
| Reflect | 3.01 | 0 | 0 | 3.3× slower |

## Map Len

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.546 | 0 | 0 | — |
| Builtin | 0.547 | 0 | 0 | 1.0× slower |
| Reflect | 3.83 | 0 | 0 | 7.0× slower |

## Kind Of

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.819 | 0 | 0 | — |
| Reflect | 1.64 | 0 | 0 | 2.0× slower |

## Is Nil_ptr

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 1.97 | 0 | 0 | — |
| Reflect | 2.46 | 0 | 0 | 1.3× slower |

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
| SafeRefl | 189 | 0 | 0 | — |
| Manual | 0.273 | 0 | 0 | 692.5× faster |
| Reflect | 416 | 0 | 0 | 2.2× slower |

## To Map

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 153 | 336 | 2 | — |
| Reflect | 339 | 392 | 7 | 2.2× slower |
| JSON | 1795 | 776 | 25 | 11.7× slower |

## Map For Each

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 595 | 0 | 0 | — |
| Range | 610 | 0 | 0 | 1.0× slower |

## Get_int

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 35.2 | 0 | 0 | — |
| Reflect | 49.6 | 0 | 0 | 1.4× slower |
| ReflectFast | 2.49 | 0 | 0 | 14.2× faster |
| Offset | 7.95 | 0 | 0 | 4.4× faster |
| Accessor | 0.546 | 0 | 0 | 64.5× faster |
| AccFrom | 2.19 | 0 | 0 | 16.1× faster |
| Direct | 0.547 | 0 | 0 | 64.5× faster |

## Set_string

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 36.0 | 0 | 0 | — |
| Reflect | 52.5 | 0 | 0 | 1.5× slower |
| ReflectFast | 3.83 | 0 | 0 | 9.4× faster |
| Offset | 9.89 | 0 | 0 | 3.6× faster |
| Accessor | 1.09 | 0 | 0 | 33.0× faster |
| AccOn | 2.46 | 0 | 0 | 14.7× faster |
| Direct | 0.547 | 0 | 0 | 65.9× faster |

## DI Resolve

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 0.559 | 0 | 0 | — |
| SafeRefl | 110 | 0 | 0 | 197.4× slower |
| Accessor | 2.46 | 0 | 0 | 4.4× slower |
| Reflect2 | 15.9 | 0 | 0 | 28.4× slower |
| Reflect | 183 | 0 | 0 | 328.0× slower |

## JSON Decode

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibJSON | 1828 | 392 | 10 | — |
| SafeRefl | 362 | 0 | 0 | 5.0× faster |
| Reflect | 1273 | 0 | 0 | 1.4× faster |
| Reflect2 | 64.5 | 0 | 0 | 28.3× faster |
| Accessor | 6.29 | 0 | 0 | 290.8× faster |

## ORM Scan

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 11.1 | 0 | 0 | — |
| SafeRefl | 370 | 0 | 0 | 33.5× slower |
| Accessor | 9.96 | 0 | 0 | 1.1× faster |
| Reflect2 | 65.4 | 0 | 0 | 5.9× slower |
| Reflect | 724 | 0 | 0 | 65.5× slower |
| Copier | 4787 | 640 | 28 | 432.9× slower |

## Struct Copy

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 7.26 | 0 | 0 | — |
| Accessor | 7.84 | 0 | 0 | 1.1× slower |
| Reflect2 | 30.1 | 0 | 0 | 4.1× slower |
| Offset | 124 | 0 | 0 | 17.1× slower |
| Reflect | 573 | 0 | 0 | 78.9× slower |
| SafeRefl | 371 | 0 | 0 | 51.1× slower |
| Copier | 2114 | 432 | 17 | 291.1× slower |

