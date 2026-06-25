//go:build reflectx_strict

package saferefl

import "testing"

// TestFieldAny_panicInStrictMode verifies that fieldAny panics in reflectx_strict builds.
// IterPlan is always nil in strict mode so this path is never reached in production;
// the panic exists to catch future regressions where the invariant is violated.
func TestFieldAny_panicInStrictMode(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("fieldAny must panic in reflectx_strict build")
		}
	}()
	fieldAny(nil, nil)
}
