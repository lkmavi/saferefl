// Package main demonstrates basic Get/Set usage of saferefl.
package main

import (
	"errors"
	"fmt"

	"github.com/lkmavi/saferefl"
)

type User struct {
	Name  string
	Age   int
	Score float64
	role  string // unexported
}

func main() {
	u := &User{Name: "Alice", Age: 30, Score: 9.5, role: "admin"}

	// --- Get ---
	name, err := saferefl.Get[string](u, "Name")
	fmt.Printf("Get Name:  %q  (err=%v)\n", name, err)

	age, err := saferefl.Get[int](u, "Age")
	fmt.Printf("Get Age:   %d  (err=%v)\n", age, err)

	score, err := saferefl.Get[float64](u, "Score")
	fmt.Printf("Get Score: %.1f  (err=%v)\n", score, err)

	// --- Set ---
	_ = saferefl.Set[string](u, "Name", "Bob")
	_ = saferefl.Set[int](u, "Age", 31)
	fmt.Printf("\nAfter Set: Name=%q Age=%d\n", u.Name, u.Age)

	// --- MustGet / MustSet ---
	saferefl.MustSet[float64](u, "Score", 10.0)
	fmt.Printf("MustGet Score: %.1f\n", saferefl.MustGet[float64](u, "Score"))

	// --- Error: type mismatch ---
	_, err = saferefl.Get[int](u, "Name") // Name is string, not int
	var tme *saferefl.TypeMismatchError
	if errors.As(err, &tme) {
		fmt.Printf("\nTypeMismatchError: field=%q got=%s want=%s\n",
			tme.FieldPath, tme.FieldType, tme.WantType)
	}

	// --- Error: field not found ---
	_, err = saferefl.Get[string](u, "Email")
	var fnf *saferefl.FieldNotFoundError
	if errors.As(err, &fnf) {
		fmt.Printf("FieldNotFoundError: field=%q on type=%s\n", fnf.FieldPath, fnf.Type)
	}

	// --- Error: unexported field ---
	err = saferefl.Set[string](u, "role", "guest")
	var roe *saferefl.ReadOnlyError
	if errors.As(err, &roe) {
		fmt.Printf("ReadOnlyError: field=%q\n", roe.FieldPath)
	}
}
