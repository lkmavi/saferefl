//go:build reflectx_strict

package unsafelayout

import (
	"errors"
	"unsafe"
)

// All unsafe primitive functions are compiled out when reflectx_strict is set.
// Callers must check AccelAvailable() before calling any other function.

func AccelAvailable() bool { return false }

func EnableAccel() error {
	return errors.New("saferefl: unsafe accelerator disabled by reflectx_strict build tag")
}

func UnsafeFieldPtr(_ unsafe.Pointer, _ uintptr) unsafe.Pointer { return nil }

func UnsafeSliceElemPtr(_ unsafe.Pointer, _ int, _ uintptr) unsafe.Pointer { return nil }

func MapLen(_ unsafe.Pointer) int { return 0 }
