//go:build !reflectx_strict

package unsafelayout

import "unsafe"

// UnsafeFieldPtr returns a pointer to a struct field at the given byte offset.
// offset must have been produced by reflect.StructField.Offset (Layer 2 guarantee).
// Safety: adding a reflect-verified offset to a pointer of a live object is explicitly
// permitted by Go's unsafe.Pointer conversion rules.
func UnsafeFieldPtr(objPtr unsafe.Pointer, offset uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(objPtr) + offset)
}

// UnsafeSliceElemPtr returns a pointer to the element at index inside a slice.
// slicePtr must point to the slice's underlying array (the Data field of reflect.SliceHeader).
// elemSize must equal reflect.Type.Size() for the element type, verified at registration.
func UnsafeSliceElemPtr(sliceData unsafe.Pointer, index int, elemSize uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(sliceData) + uintptr(index)*elemSize)
}

// MapLen returns the number of live elements in the map whose internal header is at m.
// m must be obtained via reflect.Value.Pointer() on a map value.
func MapLen(m unsafe.Pointer) int {
	return activeMapLayout.MapLen(m)
}
