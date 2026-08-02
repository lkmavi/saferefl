# Benchmark Results

**Generated:** 2026-08-02 08:26 UTC  
**Go:** 1.22  
**Platform:** linux/amd64  
**CPU:** AMD EPYC 7763 64-Core Processor                  

## Field Read

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Reflect | 68.9 | 8 | 1 | — |
| ReflectFast | 3.12 | 0 | 0 | 22.1× faster |
| SafeRefl | 45.2 | 0 | 0 | 1.5× faster |
| Offset | 11.6 | 0 | 0 | 6.0× faster |
| Accessor | 0.938 | 0 | 0 | 73.5× faster |
| Native | 0.313 | 0 | 0 | 220.5× faster |

## Slice At

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 1.01 | 0 | 0 | — |
| Direct | 0.999 | 0 | 0 | 1.0× faster |
| Reflect | 3.76 | 0 | 0 | 3.7× slower |

## Map Len

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.625 | 0 | 0 | — |
| Builtin | 0.315 | 0 | 0 | 2.0× faster |
| Reflect | 4.68 | 0 | 0 | 7.5× slower |

## Kind Of

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 0.936 | 0 | 0 | — |
| Reflect | 2.18 | 0 | 0 | 2.3× slower |

## Is Nil_ptr

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 3.43 | 0 | 0 | — |
| Reflect | 2.50 | 0 | 0 | 1.4× faster |

## Each Field

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 51.2 | 0 | 0 | — |
| Reflect | 158 | 56 | 5 | 3.1× slower |
| ReflectFull | 161 | 56 | 5 | 3.1× slower |
| Reflect2 | 24.9 | 0 | 0 | 2.1× faster |

## Copy Fields

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 286 | 0 | 0 | — |
| Manual | 0.312 | 0 | 0 | 916.3× faster |
| Reflect | 595 | 80 | 10 | 2.1× slower |

## To Map

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 229 | 336 | 2 | — |
| Reflect | 520 | 432 | 12 | 2.3× slower |
| JSON | 2719 | 776 | 25 | 11.9× slower |

## Map For Each

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 952 | 0 | 0 | — |
| Range | 984 | 0 | 0 | 1.0× slower |

## Get_int

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 45.1 | 0 | 0 | — |
| Reflect | 70.4 | 8 | 1 | 1.6× slower |
| ReflectFast | 3.12 | 0 | 0 | 14.4× faster |
| Offset | 11.5 | 0 | 0 | 3.9× faster |
| Accessor | 0.936 | 0 | 0 | 48.2× faster |
| AccFrom | 3.12 | 0 | 0 | 14.5× faster |
| Direct | 0.312 | 0 | 0 | 144.3× faster |

## Set_string

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| SafeRefl | 45.9 | 0 | 0 | — |
| Reflect | 73.3 | 8 | 1 | 1.6× slower |
| ReflectFast | 5.00 | 0 | 0 | 9.2× faster |
| Offset | 13.8 | 0 | 0 | 3.3× faster |
| Accessor | 1.25 | 0 | 0 | 36.8× faster |
| AccOn | 3.43 | 0 | 0 | 13.4× faster |
| Direct | 0.624 | 0 | 0 | 73.6× faster |

## DI Resolve

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 0.630 | 0 | 0 | — |
| SafeRefl | 154 | 0 | 0 | 244.0× slower |
| Accessor | 3.12 | 0 | 0 | 5.0× slower |
| Reflect2 | 20.9 | 0 | 0 | 33.2× slower |
| Reflect | 242 | 24 | 3 | 384.1× slower |

## JSON Decode

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| StdlibJSON | 2697 | 392 | 10 | — |
| SafeRefl | 524 | 0 | 0 | 5.1× faster |
| Reflect | 1730 | 160 | 20 | 1.6× faster |
| Reflect2 | 80.8 | 0 | 0 | 33.4× faster |
| Accessor | 9.55 | 0 | 0 | 282.5× faster |

## ORM Scan

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 14.2 | 0 | 0 | — |
| SafeRefl | 542 | 0 | 0 | 38.3× slower |
| Accessor | 12.3 | 0 | 0 | 1.1× faster |
| Reflect2 | 82.1 | 0 | 0 | 5.8× slower |
| Reflect | 983 | 80 | 10 | 69.4× slower |
| Copier | 7194 | 880 | 58 | 508.2× slower |

## Struct Copy

| Variant | ns/op | B/op | allocs/op | vs first |
|---|---|---|---|---|
| Manual | 9.67 | 0 | 0 | — |
| Accessor | 10.6 | 0 | 0 | 1.1× slower |
| Reflect2 | 40.4 | 0 | 0 | 4.2× slower |
| Offset | 186 | 0 | 0 | 19.2× slower |
| Reflect | 782 | 80 | 10 | 80.8× slower |
| SafeRefl | 520 | 0 | 0 | 53.7× slower |
| Copier | 3210 | 552 | 32 | 331.9× slower |

