package saferefl

import "github.com/lkmavi/saferefl/internal/unsafelayout"

// EnableAccel checks whether the Layer 3 unsafe accelerator passed its self-test
// on the current Go version and architecture.
//
// It is safe to call at any time; the self-test runs once during package init.
// A non-nil error means Layer 3 is disabled and all operations use Layer 2 (reflect path).
// The absence of a call to EnableAccel is safe: degradation is automatic and silent.
func EnableAccel() error {
	return unsafelayout.EnableAccel()
}
