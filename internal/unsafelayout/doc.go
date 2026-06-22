// Package unsafelayout is the Layer 3 unsafe accelerator for saferefl.
//
// # Invariants
//
//  1. Every offset passed to [UnsafeFieldPtr] was produced by reflect.StructField.Offset
//     during TypeDescriptor construction (Layer 2) and is therefore reflect-verified.
//
//  2. [runSelfTest] has completed successfully before any Layer 3 function is invoked
//     at the call site. [AccelAvailable] reports the outcome.
//
//  3. No Layer 3 function is used if [AccelAvailable] returns false.
//
// # Build tags
//
//   - (no tag)            Layer 3 compiled; self-test runs at package init.
//   - reflectx_strict     Package compiles to no-op stubs; zero unsafe code included.
//   - saferefl_strict_panic  Self-test failure causes panic (for CI / security builds).
//
// # Map backend compatibility
//
//   - !go1.24 (hmap)      hmap.count is at offset 0; stable since Go 1.
//   - go1.24  (Swiss)     Map.used  is at offset 0; verified by self-test at startup.
//
// A self-test failure is not fatal by default: [AccelAvailable] returns false and
// all callers transparently fall back to the Layer 2 reflect path.
package unsafelayout
