# Changelog

All notable changes to saferefl will be documented in this file.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

**General-purpose reflection API**
- `KindOf(v any) reflect.Kind` — fast kind check via raw `abi.Type` byte read (~0.28 ns, 0 allocs)
- `IsNil(v any) bool` — nil check for pointer, map, chan, func, slice, and interface values
- `EachField(obj any, fn func(name string, val any) bool) error` — iterate exported fields in declaration order; early-stop on `fn` returning false; handles embedded structs and nil embedded pointers
- `CopyFields(src, dst any) error` — copy matching exported fields between any two struct pointers; uses `AssignableTo`/`ConvertibleTo` for type flexibility; silently skips incompatible fields
- `ToMap(obj any) (map[string]any, error)` — convert struct to `map[string]any` keyed by field name
- `ToMapByTag(obj any, tagKey string) (map[string]any, error)` — convert struct to map keyed by tag value; skips `"-"` and empty tag names; strips `omitempty` and other options from keys
- `FromMap(m map[string]any, dst any) error` — populate struct from map; tries `AssignableTo` then `ConvertibleTo` (e.g. `float64`→`int` from JSON); returns `TypeMismatchError` on failure; skips unknown keys and nil map values
- `MapForEach[K comparable, V any](m map[K]V, fn func(K, V) bool)` — typed map iteration with early-stop; identical speed to plain `range` (~0 allocs)
- `TypeDescriptorOf(t reflect.Type) *TypeDescriptor` — expose the prebuilt type-metadata cache publicly
- `TypeDescriptor`, `FieldMeta` — public type aliases for `internal/typeinfo` types

**Tag-based field access**
- `GetByTag[T]`, `SetByTag[T]` — field access by struct tag value (e.g. `json:"name"`, `db:"col"`)
- Sentinel errors `ErrFieldNotFound`, `ErrTypeMismatch`, `ErrReadOnly` for use with `errors.Is`
- `Unwrap()` on all typed errors so `errors.Is` and `errors.As` both work

**Testing and CI**
- `efaceKind` self-test at package init — catches `abi.Type.Kind_` layout changes at startup
- `FuzzUnsafeFieldPtr` — fuzz test for unsafe field pointer correctness
- `TestConsistency` — verifies `Get[any]` and `UnsafeFieldPtr` agree with reflect baseline
- CI: `fuzz` job (60s `FuzzGet` + `FuzzUnsafeFieldPtr`) and `consistency` job (`-race`)
- CI: `test-32bit` job — runs tests on `linux/386` (native) and `linux/arm` (QEMU) with `SAFEREFL_STRICT=1` to validate portable `kindOffset` on 32-bit platforms
- CI: `test-arm64` job — native macOS arm64 run with race detector and `SAFEREFL_STRICT=1`
- Benchmarks: `KindOf`, `IsNil`, `EachField`, `CopyFields`, `ToMap`, `MapForEach` vs stdlib reflect, reflect2, JSON, and manual baselines

**Documentation**
- `docs/invariants.md` — formal Layer 3 invariants and Go-release checklist
- `docs/migrate-from-reflect2.md` — migration guide from reflect2 (no code dependency)

### Fixed
- `efaceKind` now computes the `abi.Type.Kind_` offset as `2*unsafe.Sizeof(uintptr(0))+7`
  instead of the hardcoded `23`, making it correct on 32-bit platforms (386, arm)

---

## [0.1.0-beta.1] - 2026-06-23

Initial beta release. All three API layers are shipped and self-verified.

### Added

**Generic API**
- `Get[T](obj any, fieldPath string) (T, error)` — type-safe struct field read, zero allocs on warm cache
- `Set[T](obj any, fieldPath string, val T) error` — type-safe struct field write
- `MustGet[T]`, `MustSet[T]` — panic variants for statically-known-valid paths
- Dot-path traversal for nested structs and pointer-to-struct fields (`"Address.City"`)
- Promoted fields from embedded structs accessible by promoted name

**Field inspection**
- `Fields(obj any) ([]reflect.StructField, error)` — direct fields of a struct value or pointer
- `FieldsOf[T]() ([]reflect.StructField, error)` — direct fields of a type without an instance
- `FieldByName[T](name string) (reflect.StructField, bool)` — single field lookup

**Accessor API**
- `MakeAccessor[T](obj any, fieldPath string) (Accessor[T], error)` — pre-resolve a field path once
- `Accessor[T].Get(objPtr unsafe.Pointer) T` — ~0.55 ns, zero allocs, pure pointer arithmetic
- `Accessor[T].Set(objPtr unsafe.Pointer, val T)` — same
- `Accessor[T].GetFrom(obj any) (T, error)` — convenience form with interface→pointer extraction
- `Accessor[T].SetOn(obj any, val T) error` — convenience form
- `UnsafePtrOf(obj any) unsafe.Pointer` — extract raw pointer for Accessor hot path

**Unsafe Primitives**
- `MapLenFast[K, V](m map[K]V) int` — map length without reflect dispatch
  - hmap backend (`!go1.24`): reads `hmap.count` at offset 0
  - Swiss Tables backend (`go1.24+`): reads `Map.used` at offset 0
- `UnsafeSliceAt[T](s []T, index int) *T` — unchecked slice element pointer (Go language guarantee)
- `AccelAvailable() bool` — reports whether self-test passed at init
- `EnableAccel() error` — returns descriptive error if accelerator is unavailable
- `reflectx_strict` build tag — compiles all unsafe code to no-op stubs
- `SAFEREFL_STRICT` env var — upgrades self-test failure to `panic` (for CI / security builds)

**Errors**
- `FieldNotFoundError` — field path does not exist on the type
- `TypeMismatchError` — field type not assignable to T; carries `FieldType` and `WantType`
- `ReadOnlyError` — attempted write to an unexported field

**Internals**
- `internal/typeinfo`: TypeDescriptor cache built once per struct type via stdlib reflect;
  `FieldsByName` for fast name lookup; `FieldsByTag` for tag-keyed ORM lookups;
  `PtrCache` (`atomic.Pointer`) eliminates `reflect.Type.Elem()` on the hot path
- `internal/unsafelayout`: self-test framework (`selfTestStruct`, `selfTestSlice`, `selfTestMap`)
  verifies layout assumptions against reflect at package init

**Tooling**
- `debug.StructDump(obj any, w io.Writer) error` — annotated hex dump of struct memory
- `scripts/pre-release-check.sh` — full local CI mirror with coverage gate (≥85%)
- `scripts/bench.sh`, `scripts/bench-versions.sh` — benchmark helpers

**CI**
- GitHub Actions: Go 1.22, 1.24, stable matrix; race detector; all build-tag combinations
- Weekly cross-version benchmark workflow with automated PR for results
- Weekly go-tip compatibility workflow; opens/updates GitHub issue on failure

**Documentation**
- `README.md` with quick-start, performance tables, architecture diagram
- `benchmarks/RESULTS.md` — baseline numbers (darwin/arm64, Apple M3 Max)
- `benchmarks/TARGETS.md` — acceptance criteria per API tier
- `CONTRIBUTING.md` with checklist for new Go version support
- `AGENTS.md` — AI agent coding guidelines for this repo

[Unreleased]: https://github.com/lkmavi/saferefl/compare/v0.1.0-beta.1...HEAD
[0.1.0-beta.1]: https://github.com/lkmavi/saferefl/releases/tag/v0.1.0-beta.1
