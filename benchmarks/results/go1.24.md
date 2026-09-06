# Benchmark Results

**Generated:** 2026-09-06 10:30 UTC  
**Go:** 1.24  
**Platform:** linux/amd64  
**CPU:** AMD EPYC 9V74 80-Core Processor                  

## Field Read

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Reflect | 64.6 | 0 | 0 | — |
| ReflectFast | 3.21 | 0 | 0 | 20.1× faster |
| SafeRefl | 45.7 | 0 | 0 | 1.4× faster |
| Offset | 10.6 | 0 | 0 | 6.1× faster |
| Accessor | 1.07 | 0 | 0 | 60.6× faster |
| Native | 0.706 | 0 | 0 | 91.5× faster |

## Slice At

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 1.16 | 0 | 0 | — |
| Direct | 1.19 | 0 | 0 | 1.0× slower |
| Reflect | 3.88 | 0 | 0 | 3.3× slower |

## Map Len

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.705 | 0 | 0 | — |
| Builtin | 0.705 | 0 | 0 | 1.0× slower |
| Reflect | 4.93 | 0 | 0 | 7.0× slower |

## Kind Of

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 1.06 | 0 | 0 | — |
| Reflect | 2.12 | 0 | 0 | 2.0× slower |

## Is Nil_ptr

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 2.47 | 0 | 0 | — |
| Reflect | 3.17 | 0 | 0 | 1.3× slower |

## Each Field

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 43.3 | 0 | 0 | — |
| Reflect | 137 | 56 | 5 | 3.2× slower |
| ReflectFull | 141 | 56 | 5 | 3.3× slower |
| Reflect2 | 25.0 | 0 | 0 | 1.7× faster |

## Copy Fields

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 246 | 0 | 0 | — |
| Manual | 0.352 | 0 | 0 | 697.8× faster |
| Reflect | 528 | 0 | 0 | 2.1× slower |

## To Map

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 197 | 336 | 2 | — |
| Reflect | 436 | 392 | 7 | 2.2× slower |
| JSON | 2324 | 776 | 25 | 11.8× slower |

## Map For Each

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 760 | 0 | 0 | — |
| Range | 770 | 0 | 0 | 1.0× slower |

## Get_int

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 44.9 | 0 | 0 | — |
| Reflect | 64.1 | 0 | 0 | 1.4× slower |
| ReflectFast | 3.20 | 0 | 0 | 14.0× faster |
| Offset | 10.4 | 0 | 0 | 4.3× faster |
| Accessor | 0.705 | 0 | 0 | 63.7× faster |
| AccFrom | 2.84 | 0 | 0 | 15.8× faster |
| Direct | 0.705 | 0 | 0 | 63.7× faster |

## Set_string

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 46.5 | 0 | 0 | — |
| Reflect | 67.9 | 0 | 0 | 1.5× slower |
| ReflectFast | 4.95 | 0 | 0 | 9.4× faster |
| Offset | 13.0 | 0 | 0 | 3.6× faster |
| Accessor | 1.41 | 0 | 0 | 33.0× faster |
| AccOn | 3.17 | 0 | 0 | 14.7× faster |
| Direct | 0.708 | 0 | 0 | 65.6× faster |

## DI Resolve

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 0.728 | 0 | 0 | — |
| SafeRefl | 144 | 0 | 0 | 197.2× slower |
| Accessor | 3.22 | 0 | 0 | 4.4× slower |
| Reflect2 | 21.0 | 0 | 0 | 28.8× slower |
| Reflect | 250 | 0 | 0 | 343.9× slower |

## JSON Decode

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibJSON | 2361 | 392 | 10 | — |
| SafeRefl | 467 | 0 | 0 | 5.1× faster |
| Reflect | 1642 | 0 | 0 | 1.4× faster |
| Reflect2 | 85.1 | 0 | 0 | 27.8× faster |
| Accessor | 8.12 | 0 | 0 | 290.7× faster |

## ORM Scan

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 14.2 | 0 | 0 | — |
| SafeRefl | 480 | 0 | 0 | 33.8× slower |
| Accessor | 12.8 | 0 | 0 | 1.1× faster |
| Reflect2 | 83.6 | 0 | 0 | 5.9× slower |
| Reflect | 939 | 0 | 0 | 66.0× slower |
| Copier | 6142 | 640 | 28 | 431.6× slower |

## Struct Copy

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 9.34 | 0 | 0 | — |
| Accessor | 10.1 | 0 | 0 | 1.1× slower |
| Reflect2 | 38.8 | 0 | 0 | 4.2× slower |
| Offset | 161 | 0 | 0 | 17.2× slower |
| Reflect | 740 | 0 | 0 | 79.2× slower |
| SafeRefl | 487 | 0 | 0 | 52.1× slower |
| Copier | 2710 | 432 | 17 | 290.1× slower |

