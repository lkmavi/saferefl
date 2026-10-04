# Benchmark Results

**Generated:** 2026-10-04 11:56 UTC  
**Go:** 1.24  
**Platform:** linux/amd64  
**CPU:** AMD EPYC 9V74 80-Core Processor                  

## Field Read

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Reflect | 49.7 | 0 | 0 | — |
| ReflectFast | 2.50 | 0 | 0 | 19.9× faster |
| SafeRefl | 36.0 | 0 | 0 | 1.4× faster |
| Offset | 8.03 | 0 | 0 | 6.2× faster |
| Accessor | 0.821 | 0 | 0 | 60.5× faster |
| Native | 0.547 | 0 | 0 | 90.9× faster |

## Slice At

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.902 | 0 | 0 | — |
| Direct | 0.921 | 0 | 0 | 1.0× slower |
| Reflect | 3.01 | 0 | 0 | 3.3× slower |

## Map Len

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.547 | 0 | 0 | — |
| Builtin | 0.546 | 0 | 0 | 1.0× faster |
| Reflect | 3.83 | 0 | 0 | 7.0× slower |

## Kind Of

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.819 | 0 | 0 | — |
| Reflect | 1.64 | 0 | 0 | 2.0× slower |

## Is Nil_ptr

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 1.91 | 0 | 0 | — |
| Reflect | 2.46 | 0 | 0 | 1.3× slower |

## Each Field

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 33.8 | 0 | 0 | — |
| Reflect | 106 | 56 | 5 | 3.1× slower |
| ReflectFull | 109 | 56 | 5 | 3.2× slower |
| Reflect2 | 19.4 | 0 | 0 | 1.7× faster |

## Copy Fields

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 190 | 0 | 0 | — |
| Manual | 0.273 | 0 | 0 | 695.4× faster |
| Reflect | 410 | 0 | 0 | 2.2× slower |

## To Map

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 154 | 336 | 2 | — |
| Reflect | 337 | 392 | 7 | 2.2× slower |
| JSON | 1848 | 776 | 25 | 12.0× slower |

## Map For Each

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 596 | 0 | 0 | — |
| Range | 604 | 0 | 0 | 1.0× slower |

## Get_int

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 35.1 | 0 | 0 | — |
| Reflect | 49.7 | 0 | 0 | 1.4× slower |
| ReflectFast | 2.48 | 0 | 0 | 14.1× faster |
| Offset | 7.94 | 0 | 0 | 4.4× faster |
| Accessor | 0.546 | 0 | 0 | 64.2× faster |
| AccFrom | 2.19 | 0 | 0 | 16.0× faster |
| Direct | 0.546 | 0 | 0 | 64.2× faster |

## Set_string

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 36.0 | 0 | 0 | — |
| Reflect | 53.0 | 0 | 0 | 1.5× slower |
| ReflectFast | 3.83 | 0 | 0 | 9.4× faster |
| Offset | 9.93 | 0 | 0 | 3.6× faster |
| Accessor | 1.09 | 0 | 0 | 33.0× faster |
| AccOn | 2.47 | 0 | 0 | 14.6× faster |
| Direct | 0.546 | 0 | 0 | 65.9× faster |

## DI Resolve

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 0.558 | 0 | 0 | — |
| SafeRefl | 111 | 0 | 0 | 199.2× slower |
| Accessor | 2.46 | 0 | 0 | 4.4× slower |
| Reflect2 | 16.6 | 0 | 0 | 29.8× slower |
| Reflect | 183 | 0 | 0 | 328.4× slower |

## JSON Decode

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibJSON | 1837 | 392 | 10 | — |
| SafeRefl | 362 | 0 | 0 | 5.1× faster |
| Reflect | 1273 | 0 | 0 | 1.4× faster |
| Reflect2 | 64.6 | 0 | 0 | 28.4× faster |
| Accessor | 6.31 | 0 | 0 | 291.2× faster |

## ORM Scan

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 11.1 | 0 | 0 | — |
| SafeRefl | 370 | 0 | 0 | 33.4× slower |
| Accessor | 9.97 | 0 | 0 | 1.1× faster |
| Reflect2 | 64.5 | 0 | 0 | 5.8× slower |
| Reflect | 724 | 0 | 0 | 65.3× slower |
| Copier | 4740 | 640 | 28 | 427.4× slower |

## Struct Copy

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 7.24 | 0 | 0 | — |
| Accessor | 7.83 | 0 | 0 | 1.1× slower |
| Reflect2 | 29.8 | 0 | 0 | 4.1× slower |
| Offset | 128 | 0 | 0 | 17.6× slower |
| Reflect | 573 | 0 | 0 | 79.2× slower |
| SafeRefl | 373 | 0 | 0 | 51.5× slower |
| Copier | 2086 | 432 | 17 | 288.1× slower |

