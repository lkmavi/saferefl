package saferefl_test

import (
	"testing"

	"github.com/lkmavi/saferefl"
)

func TestAccelAvailable(t *testing.T) {
	// Verify the function is callable and returns a consistent value.
	got := saferefl.AccelAvailable()
	if got != saferefl.AccelAvailable() {
		t.Error("AccelAvailable() is not idempotent")
	}
}

func TestEnableAccel(t *testing.T) {
	err := saferefl.EnableAccel()
	if saferefl.AccelAvailable() && err != nil {
		t.Errorf("EnableAccel() = %v, want nil when AccelAvailable=true", err)
	}
	if !saferefl.AccelAvailable() && err == nil {
		t.Error("EnableAccel() = nil, want error when AccelAvailable=false")
	}
}

func TestUnsafeSliceAt_read(t *testing.T) {
	s := []int{10, 20, 30}
	for i, want := range s {
		got := *saferefl.UnsafeSliceAt(s, i)
		if got != want {
			t.Errorf("UnsafeSliceAt[%d] = %d, want %d", i, got, want)
		}
	}
}

func TestUnsafeSliceAt_write(t *testing.T) {
	s := []int{1, 2, 3}
	*saferefl.UnsafeSliceAt(s, 1) = 99
	if s[1] != 99 {
		t.Error("write through UnsafeSliceAt did not persist in original slice")
	}
}

func TestUnsafeSliceAt_strings(t *testing.T) {
	s := []string{"a", "b", "c"}
	if got := *saferefl.UnsafeSliceAt(s, 2); got != "c" {
		t.Errorf("UnsafeSliceAt[2] = %q, want %q", got, "c")
	}
}

func TestMapLenFast_nil(t *testing.T) {
	var m map[string]int
	if got := saferefl.MapLenFast(m); got != 0 {
		t.Errorf("MapLenFast(nil) = %d, want 0", got)
	}
}

func TestMapLenFast_empty(t *testing.T) {
	m := make(map[string]int)
	if got := saferefl.MapLenFast(m); got != 0 {
		t.Errorf("MapLenFast(empty) = %d, want 0", got)
	}
}

func TestMapLenFast_matchesBuiltin(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	if got, want := saferefl.MapLenFast(m), len(m); got != want {
		t.Errorf("MapLenFast = %d, len = %d", got, want)
	}
}

func TestMapLenFast_intKey(t *testing.T) {
	m := map[int]string{1: "one", 2: "two"}
	if got := saferefl.MapLenFast(m); got != 2 {
		t.Errorf("MapLenFast = %d, want 2", got)
	}
}
