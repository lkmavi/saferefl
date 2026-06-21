package saferefl

import "fmt"

// FieldNotFoundError is returned when the field path does not resolve to a field on the type.
type FieldNotFoundError struct {
	Type      string
	FieldPath string
}

func (e *FieldNotFoundError) Error() string {
	return fmt.Sprintf("saferefl: field %q not found on type %s", e.FieldPath, e.Type)
}

// TypeMismatchError is returned when the field's type is not assignable to the requested type T.
type TypeMismatchError struct {
	FieldPath string
	FieldType string // actual type stored in the field
	WantType  string // requested type parameter T
}

func (e *TypeMismatchError) Error() string {
	return fmt.Sprintf("saferefl: field %q has type %s, cannot use as %s", e.FieldPath, e.FieldType, e.WantType)
}

// ReadOnlyError is returned when attempting to set an unexported (read-only) field.
type ReadOnlyError struct {
	FieldPath string
}

func (e *ReadOnlyError) Error() string {
	return fmt.Sprintf("saferefl: field %q is unexported (read-only)", e.FieldPath)
}
