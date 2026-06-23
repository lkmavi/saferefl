# Migrating from reflect2

This guide shows the saferefl equivalent for each common `reflect2` operation.
saferefl has no runtime dependency on reflect2.

---

## Type lookup

| reflect2 | saferefl |
|---|---|
| `reflect2.TypeOf(v)` | `reflect.TypeOf(v)` (unchanged — use stdlib) |
| `reflect2.Type2(reflect.TypeOf(v))` | `typeinfo.TypeDescriptorOf(reflect.TypeOf(v))` |

`TypeDescriptor` is cached globally on first access, so repeated calls are O(1).

---

## Struct field access

| reflect2 | saferefl |
|---|---|
| `t.FieldByName("X")` | `saferefl.FieldByName[S]("X")` |
| `f.UnsafeGet(ptr)` | `saferefl.Get[T](ptr, "X")` |
| `f.UnsafeSet(ptr, val)` | `saferefl.Set[T](ptr, "X", val)` |

`Get[T]` and `Set[T]` are type-safe generics — no `interface{}` casts at the call site.

---

## Hot-path field access (zero per-call overhead)

reflect2 encourages caching the field object and calling `UnsafeGet`/`UnsafeSet` per object.
saferefl's `Accessor[T]` is the direct equivalent:

```go
// reflect2 style
t   := reflect2.TypeOf(MyStruct{}).(*reflect2.UnsafeStructType)
f   := t.FieldByName("Score")
val := f.UnsafeGet(ptr).(float64)
f.UnsafeSet(ptr, 3.14)

// saferefl Accessor style — built once, used many times
acc := saferefl.MakeAccessor[float64](&MyStruct{}, "Score")
val := acc.Get(&s)
acc.Set(&s, 3.14)
```

`Accessor[T]` compiles to a single pointer-arithmetic instruction on the hot path
(≈ 0.55 ns/op, within 2× of a native field access).

---

## Slice element access

| reflect2 | saferefl |
|---|---|
| `reflect2.PtrTo(elemType).UnsafeIndirect(...)` | `saferefl.SliceAt[T](s, i)` |

---

## Map length (no alloc)

| reflect2 | saferefl |
|---|---|
| n/a (reflect2 has no fast MapLen) | `saferefl.MapLenFast(m)` |

`MapLenFast` reads the runtime map header directly (≈ 0.28 ns/op, same as `len(m)`).
It requires `AccelAvailable() == true`; falls back to `reflect.ValueOf(m).Len()` otherwise.

---

## Error handling

reflect2 panics on type mismatches. saferefl returns typed errors:

```go
v, err := saferefl.Get[int](obj, "Name") // Name is string
// err is *saferefl.TypeMismatchError, not a panic
```

Use `errors.As` to inspect:

```go
var mismatch *saferefl.TypeMismatchError
if errors.As(err, &mismatch) {
    fmt.Println(mismatch.FieldType, mismatch.WantType)
}
```

---

## Build tags

| reflect2 concern | saferefl equivalent |
|---|---|
| No strict mode | `reflectx_strict` build tag — compiles out all unsafe code |
| No self-test | `saferefl.EnableAccel()` at startup — returns error if self-test failed |

---

## Dependency

reflect2 imports `github.com/modern-go/reflect2`. saferefl has zero non-stdlib dependencies
in its main module.
