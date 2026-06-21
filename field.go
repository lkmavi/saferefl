package saferefl

import (
	"fmt"
	"reflect"
)

// FieldByName returns the reflect.StructField for the named field of struct type T.
// Searches direct fields and promoted fields from embedded structs,
// identical to reflect.Type.FieldByName.
func FieldByName[T any](name string) (reflect.StructField, bool) {
	t := reflect.TypeOf((*T)(nil)).Elem()
	if t.Kind() != reflect.Struct {
		return reflect.StructField{}, false
	}
	return t.FieldByName(name)
}

// Fields returns the direct struct fields of obj's type in declaration order.
// obj may be a struct value or a pointer to struct.
func Fields(obj any) ([]reflect.StructField, error) {
	if obj == nil {
		return nil, fmt.Errorf("saferefl: obj must not be nil")
	}
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("saferefl: obj must be a struct or pointer to struct, got %v", t.Kind())
	}
	n := t.NumField()
	out := make([]reflect.StructField, n)
	for i := range n {
		out[i] = t.Field(i)
	}
	return out, nil
}

// FieldsOf returns the direct struct fields of type T in declaration order.
// Returns an error if T is not a struct type.
func FieldsOf[T any]() ([]reflect.StructField, error) {
	t := reflect.TypeOf((*T)(nil)).Elem()
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("saferefl: FieldsOf requires a struct type, got %v", t.Kind())
	}
	n := t.NumField()
	out := make([]reflect.StructField, n)
	for i := range n {
		out[i] = t.Field(i)
	}
	return out, nil
}
