// Package main demonstrates dot-path traversal through nested and pointer fields.
package main

import (
	"fmt"

	"github.com/lkmavi/saferefl"
)

type Address struct {
	City    string
	Country string
}

type Person struct {
	Name string
	Age  int
}

type Employee struct {
	Person          // embedded — Name and Age are promoted
	Company string
	Office  Address  // value field
	Home    *Address // pointer field — nil is safe
}

func main() {
	e := &Employee{
		Person:  Person{Name: "Carol", Age: 28},
		Company: "Acme",
		Office:  Address{City: "Berlin", Country: "DE"},
		Home:    &Address{City: "Munich", Country: "DE"},
	}

	// Promoted field from embedded Person
	name, _ := saferefl.Get[string](e, "Name")
	fmt.Printf("Name (promoted):    %q\n", name)

	// Direct field
	company, _ := saferefl.Get[string](e, "Company")
	fmt.Printf("Company:            %q\n", company)

	// Dot-path through a value struct field
	city, _ := saferefl.Get[string](e, "Office.City")
	fmt.Printf("Office.City:        %q\n", city)

	// Dot-path through a pointer field
	homeCity, _ := saferefl.Get[string](e, "Home.City")
	fmt.Printf("Home.City:          %q\n", homeCity)

	// Set through a value struct field
	_ = saferefl.Set[string](e, "Office.City", "Hamburg")
	fmt.Printf("Office.City after Set: %q\n", e.Office.City)

	// Nil pointer in path — returns error, never panics
	e.Home = nil
	_, err := saferefl.Get[string](e, "Home.City")
	fmt.Printf("\nNil pointer path: err=%v\n", err)

	// Non-existent intermediate segment — FieldNotFoundError
	_, err = saferefl.Get[string](e, "Address.City")
	fmt.Printf("Bad segment:      err=%v\n", err)
}
