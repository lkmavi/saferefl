// Package main demonstrates field inspection using FieldByName, Fields, and FieldsOf.
package main

import (
	"fmt"

	"github.com/lkmavi/saferefl"
)

type Product struct {
	ID       int    `json:"id" db:"product_id"`
	Name     string `json:"name"`
	Price    float64
	inStock  bool // unexported
}

func main() {
	// FieldsOf[T] — no instance needed, compile-time type
	fields, err := saferefl.FieldsOf[Product]()
	if err != nil {
		panic(err)
	}
	fmt.Println("=== FieldsOf[Product] ===")
	for _, f := range fields {
		fmt.Printf("  %-10s  exported=%-5v  tag=%q\n", f.Name, f.IsExported(), string(f.Tag))
	}

	// Fields — from an instance (value or pointer)
	p := &Product{ID: 1, Name: "Widget", Price: 9.99}
	fields, err = saferefl.Fields(p)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n=== Fields(&product) — %d direct fields ===\n", len(fields))
	for _, f := range fields {
		fmt.Printf("  %-10s  kind=%v\n", f.Name, f.Type.Kind())
	}

	// FieldByName[T] — look up a single field descriptor by name
	sf, ok := saferefl.FieldByName[Product]("Name")
	fmt.Printf("\n=== FieldByName[Product](\"Name\") ===\n")
	fmt.Printf("  found=%v  name=%q  type=%v  json=%q\n",
		ok, sf.Name, sf.Type, sf.Tag.Get("json"))

	// Not found
	_, ok = saferefl.FieldByName[Product]("Missing")
	fmt.Printf("\nFieldByName[Product](\"Missing\")  found=%v\n", ok)
}
