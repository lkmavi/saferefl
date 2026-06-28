# Benchmark Results

**Generated:** 2026-06-28 08:56 UTC  
**Go:** 1.24  
**Platform:** linux/amd64  
**CPU:** AMD EPYC 9V74 80-Core Processor                  

## Field Read

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Reflect | 64.3 | 0 | 0 | — |
| ReflectFast | 3.19 | 0 | 0 | 20.2× faster |
| SafeRefl | 45.4 | 0 | 0 | 1.4× faster |
| Offset | 10.5 | 0 | 0 | 6.1× faster |
| Accessor | 1.06 | 0 | 0 | 60.7× faster |
| Native | 0.705 | 0 | 0 | 91.2× faster |

## Slice At

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 1.16 | 0 | 0 | — |
| Direct | 1.19 | 0 | 0 | 1.0× slower |
| Reflect | 3.88 | 0 | 0 | 3.3× slower |

## Map Len

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.709 | 0 | 0 | — |
| Builtin | 0.706 | 0 | 0 | 1.0× faster |
| Reflect | 4.94 | 0 | 0 | 7.0× slower |

## Kind Of

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 1.06 | 0 | 0 | — |
| Reflect | 2.11 | 0 | 0 | 2.0× slower |

## Is Nil_ptr

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 2.48 | 0 | 0 | — |
| Reflect | 3.17 | 0 | 0 | 1.3× slower |

## Each Field

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 43.4 | 0 | 0 | — |
| Reflect | 136 | 56 | 5 | 3.1× slower |
| ReflectFull | 140 | 56 | 5 | 3.2× slower |
| Reflect2 | 25.0 | 0 | 0 | 1.7× faster |

## Copy Fields

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 247 | 0 | 0 | — |
| Manual | 0.352 | 0 | 0 | 701.8× faster |
| Reflect | 539 | 0 | 0 | 2.2× slower |

## To Map

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 197 | 336 | 2 | — |
| Reflect | 439 | 392 | 7 | 2.2× slower |
| JSON | 2316 | 776 | 25 | 11.8× slower |

## Map For Each

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 777 | 0 | 0 | — |
| Range | 794 | 0 | 0 | 1.0× slower |

## Get_int

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 45.3 | 0 | 0 | — |
| Reflect | 64.0 | 0 | 0 | 1.4× slower |
| ReflectFast | 3.20 | 0 | 0 | 14.2× faster |
| Offset | 10.3 | 0 | 0 | 4.4× faster |
| Accessor | 0.708 | 0 | 0 | 64.0× faster |
| AccFrom | 2.82 | 0 | 0 | 16.1× faster |
| Direct | 0.705 | 0 | 0 | 64.3× faster |

## Set_string

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 46.5 | 0 | 0 | — |
| Reflect | 67.9 | 0 | 0 | 1.5× slower |
| ReflectFast | 4.94 | 0 | 0 | 9.4× faster |
| Offset | 12.9 | 0 | 0 | 3.6× faster |
| Accessor | 1.41 | 0 | 0 | 32.9× faster |
| AccOn | 3.54 | 0 | 0 | 13.1× faster |
| Direct | 0.706 | 0 | 0 | 65.8× faster |

## DI Resolve

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 0.722 | 0 | 0 | — |
| SafeRefl | 143 | 0 | 0 | 198.4× slower |
| Accessor | 3.20 | 0 | 0 | 4.4× slower |
| Reflect2 | 21.8 | 0 | 0 | 30.3× slower |
| Reflect | 236 | 0 | 0 | 327.3× slower |

## JSON Decode

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibJSON | 2358 | 392 | 10 | — |
| SafeRefl | 468 | 0 | 0 | 5.0× faster |
| Reflect | 1642 | 0 | 0 | 1.4× faster |
| Reflect2 | 80.4 | 0 | 0 | 29.3× faster |
| Accessor | 8.11 | 0 | 0 | 290.9× faster |

## ORM Scan

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 14.2 | 0 | 0 | — |
| SafeRefl | 481 | 0 | 0 | 33.8× slower |
| Accessor | 12.8 | 0 | 0 | 1.1× faster |
| Reflect2 | 82.2 | 0 | 0 | 5.8× slower |
| Reflect | 937 | 0 | 0 | 65.8× slower |
| Copier | 6189 | 640 | 28 | 434.5× slower |

## Struct Copy

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 9.36 | 0 | 0 | — |
| Accessor | 10.1 | 0 | 0 | 1.1× slower |
| Reflect2 | 38.6 | 0 | 0 | 4.1× slower |
| Offset | 164 | 0 | 0 | 17.5× slower |
| Reflect | 740 | 0 | 0 | 79.0× slower |
| SafeRefl | 486 | 0 | 0 | 51.9× slower |
| Copier | 2711 | 432 | 17 | 289.5× slower |

