//go:build reflectx_strict

package unsafelayout

import "testing"

func TestDisabled_AccelAvailable(t *testing.T) {
	if AccelAvailable() {
		t.Error("AccelAvailable() = true, want false with reflectx_strict")
	}
}

func TestDisabled_EnableAccel(t *testing.T) {
	err := EnableAccel()
	if err == nil {
		t.Error("EnableAccel() = nil, want error with reflectx_strict")
	}
}

func TestDisabled_MapLen(t *testing.T) {
	if n := MapLen(nil); n != 0 {
		t.Errorf("MapLen(nil) = %d, want 0 with reflectx_strict", n)
	}
}

func TestDisabled_UnsafeFieldPtr(t *testing.T) {
	if p := UnsafeFieldPtr(nil, 0); p != nil {
		t.Errorf("UnsafeFieldPtr returned non-nil %p with reflectx_strict", p)
	}
}

func TestDisabled_UnsafeSliceElemPtr(t *testing.T) {
	if p := UnsafeSliceElemPtr(nil, 0, 1); p != nil {
		t.Errorf("UnsafeSliceElemPtr returned non-nil %p with reflectx_strict", p)
	}
}
