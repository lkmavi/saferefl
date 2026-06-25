package saferefl

import (
	"errors"
	"fmt"
)

// Sentinel errors for use with [errors.Is].
// Typed errors ([FieldNotFoundError], [TypeMismatchError], [ReadOnlyError]) wrap these
// so callers can use either errors.Is (simple check) or errors.As (access details).
var (
	ErrFieldNotFound = errors.New("saferefl: field not found")
	ErrTypeMismatch  = errors.New("saferefl: type mismatch")
	ErrReadOnly      = errors.New("saferefl: field is read-only")
)

// FieldNotFoundError is returned when the field path does not resolve to a field on the type.
type FieldNotFoundError struct {
	Type      string
	FieldPath string
}

func (e *FieldNotFoundError) Error() string {
	return fmt.Sprintf("saferefl: field %q not found on type %s", e.FieldPath, e.Type)
}

func (e *FieldNotFoundError) Unwrap() error { return ErrFieldNotFound }

// TypeMismatchError is returned when the field's type is not assignable to the requested type T.
type TypeMismatchError struct {
	FieldPath string
	FieldType string // actual type stored in the field
	WantType  string // requested type parameter T
}

func (e *TypeMismatchError) Error() string {
	return fmt.Sprintf("saferefl: field %q has type %s, cannot use as %s", e.FieldPath, e.FieldType, e.WantType)
}

func (e *TypeMismatchError) Unwrap() error { return ErrTypeMismatch }

// ReadOnlyError is returned when attempting to set an unexported (read-only) field.
type ReadOnlyError struct {
	FieldPath string
}

func (e *ReadOnlyError) Error() string {
	return fmt.Sprintf("saferefl: field %q is unexported (read-only)", e.FieldPath)
}

func (e *ReadOnlyError) Unwrap() error { return ErrReadOnly }
