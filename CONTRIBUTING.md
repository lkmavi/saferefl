# Contributing to saferefl

Thank you for your interest in contributing!

---

## Requirements

- **Go 1.22+** (minimum supported version per [ADR-01](/_local/ADR-01-fast-safe-reflection.md))
- **golangci-lint** for code quality checks

---

## Quick Start

```bash
git clone https://github.com/lkmavi/saferefl
cd saferefl

go build ./...
go test -race ./...
golangci-lint run --timeout=5m
```

---

## Development Workflow

### 1. Fork & Clone

```bash
git clone https://github.com/YOUR_USERNAME/saferefl
cd saferefl
git remote add upstream https://github.com/lkmavi/saferefl
```

### 2. Create Feature Branch

```bash
git checkout -b feat/your-feature
# or
git checkout -b fix/issue-number-description
```

### 3. Make Changes

- Follow code style guidelines below
- Add tests for all new functionality — the library must maintain correctness guarantees across Go versions
- If touching `internal/unsafelayout`: document the invariant being relied upon and update the self-test

### 4. Validate Before Commit

```bash
go fmt ./...
bash scripts/pre-release-check.sh
```

### 5. Create Pull Request

```bash
git add .
git commit -m "feat(cache): add FieldOffset caching for embedded structs"
git push origin feat/your-feature
```

Then open a PR on GitHub.

---

## Pull Request Guidelines

### Requirements

- [ ] All tests pass: `go test -race ./...`
- [ ] All build-tag variants build: `-tags reflectx_strict`
- [ ] Linter passes: `golangci-lint run`
- [ ] Code is formatted: `go fmt ./...`
- [ ] CHANGELOG.md updated (for features/fixes)
- [ ] Unsafe code changes include a self-test in `internal/unsafelayout`

### PR Title Format

```
feat(api): add Set[T] generic setter with field path
feat(cache): cache TypeInfo offsets for embedded structs
feat(accel): add Swiss Tables map backend for Go 1.24+
fix(cache): race condition in sync.Map lookup under high concurrency
docs: add migration guide from reflect2
test(accel): fuzz test for reflect vs unsafe result divergence
perf(api): reduce allocations in Get[T] fast path
ci: add Go 1.24 to test matrix
chore: bump golangci-lint to v2.x
```

### PR Description Template

```markdown
## Summary
Brief description of changes.

## Changes
- Change 1
- Change 2

## Testing
How was this tested? Include benchmark results for perf changes.

## Area impact
Which area does this touch (Generic API / TypeInfo Cache / Unsafe Primitives)?
If Unsafe Primitives: was the self-test updated?

## Related Issues
Closes #123
```

---

## Code Style

### Go Conventions

- `gofmt` for formatting (tabs, not spaces)
- Follow [Effective Go](https://go.dev/doc/effective_go)
- Generics-first: prefer `Get[T]/Set[T]` over `interface{}` wherever the type is known at compile time

### Naming

| Type | Convention | Example |
|------|------------|---------|
| Exported | PascalCase | `TypeInfo`, `FieldOffset`, `EnableAccel` |
| Unexported | camelCase | `buildTypeInfo`, `lookupOffset` |
| Acronyms | Uppercase | `GetRType`, `UnsafePtr` |
| Constants | PascalCase | `MaxCacheSize` |

### Error Handling

```go
if err != nil {
    return fmt.Errorf("saferefl: operation failed: %w", err)
}
```

---

## Architecture

This library has three concerns. Contributions must respect their boundaries:

| Area | Location | What it does | Unsafe? |
|------|----------|--------------|---------|
| Generic API | `*.go` (public) | `Get[T]`, `Set[T]`, `Accessor[T]` | No |
| TypeInfo Cache | `internal/typeinfo/` | Struct metadata via stdlib `reflect`, `sync.Map` | No |
| Unsafe Primitives | `internal/unsafelayout/` | Optional accelerator: direct map/slice reads | Yes, self-tested |

**Rules for Unsafe Primitives:**
- Every `unsafe.Pointer` offset access must be validated by a synthetic self-test in `init()`
- No `go:linkname` to private runtime symbols
- Files must use build constraints: `//go:build go1.24` for Swiss Tables backend, `//go:build !go1.24` for legacy
- If the self-test fails at init, log a warning — the caller should check [AccelAvailable]

---

## Commit Messages

We use [Conventional Commits](https://www.conventionalcommits.org/):

```
type(scope): description
```

### Types

| Type | Description |
|------|-------------|
| `feat` | New feature |
| `fix` | Bug fix |
| `docs` | Documentation |
| `test` | Tests |
| `refactor` | Refactoring |
| `perf` | Performance improvement |
| `ci` | CI/CD changes |
| `chore` | Maintenance |

### Scopes

| Scope | Description |
|-------|-------------|
| `api` | Generic public API (`Get[T]`, `Set[T]`, `Accessor[T]`) |
| `cache` | TypeInfo cache internals |
| `accel` | Unsafe accelerator (`internal/unsafelayout`) |
| `selftest` | Self-test framework |
| `bench` | Benchmarks |
| `ci` | CI configuration |
| `deps` | Dependencies |

---

## Testing

### Run all tests

```bash
go test -race ./...
```

### Run strict mode (unsafe layer excluded from binary)

```bash
go test -race -tags reflectx_strict ./...
```

### Run benchmarks

```bash
go test -bench=. -benchmem ./...
```

### Run fuzz tests (Go 1.22+)

```bash
go test -fuzz=FuzzUnsafeVsReflect -fuzztime=60s ./internal/unsafelayout/...
```

### Pre-Release Validation

```bash
bash scripts/pre-release-check.sh
```

---

## Where We Need Help

- **Benchmarks** — realistic profiles for JSON codec, ORM mapping, DI container scenarios
- **Go version testing** — run the test suite on Go 1.22, 1.23, 1.24, and Go tip and report results
- **Map backends** — Swiss Tables map layout for Go 1.24+ (`internal/unsafelayout`)
- **Fuzz tests** — catching divergence between reflect-based and unsafe results
- **Documentation** — usage examples, migration guide from `reflect2` / `json-iterator`

---

## Questions?

- Open a [GitHub Issue](https://github.com/lkmavi/saferefl/issues)
- Check existing [Discussions](https://github.com/lkmavi/saferefl/discussions)

---

*Thank you for contributing to saferefl!*
