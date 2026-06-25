//go:build !reflectx_strict

package saferefl_test

import (
	"reflect"
	"testing"

	"github.com/lkmavi/saferefl"
)

// TestTypeDescriptorOf_shadowingIterPlan verifies that collectIter's seen-map dedup
// (fix 2) correctly drops the promoted Base.X in favor of the outer Outer.X.
func TestTypeDescriptorOf_shadowingIterPlan(t *testing.T) {
	type shadowBase struct{ X string }
	type shadowOuter struct {
		shadowBase     // promoted X string — must be excluded from IterPlan
		X          int // direct X int   — must win
	}
	desc := saferefl.TypeDescriptorOf(reflect.TypeOf(shadowOuter{}))

	// IterPlan must contain X exactly once with the outer type (int).
	count := 0
	for _, e := range desc.IterPlan {
		if e.Name == "X" {
			count++
			if e.Type != reflect.TypeOf(0) {
				t.Errorf("IterPlan X type = %v, want int", e.Type)
			}
		}
	}
	if count != 1 {
		t.Errorf("IterPlan has %d entries named X, want 1", count)
	}

	// FieldsByName must also resolve to the outer int field (fix 1 — collectNamed shadowing).
	fm, ok := desc.FieldsByName["X"]
	if !ok {
		t.Fatal("FieldsByName missing X")
	}
	if fm.Type != reflect.TypeOf(0) {
		t.Errorf("FieldsByName X type = %v, want int", fm.Type)
	}
}
