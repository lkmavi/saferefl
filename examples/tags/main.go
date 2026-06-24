// Package main demonstrates GetByTag and SetByTag — accessing struct fields
// by their tag value rather than their Go field name.
// Useful for ORM-style column mapping, JSON key access, and config binding.
package main

import (
	"errors"
	"fmt"

	"github.com/lkmavi/saferefl"
)

type Product struct {
	ID       int     `json:"id"        db:"product_id"`
	Name     string  `json:"name"      db:"product_name"`
	Price    float64 `json:"price"     db:"unit_price"`
	InStock  bool    `json:"in_stock"  db:"in_stock"`
	Internal string  `json:"-"`        // excluded from json tag lookup
}

func main() {
	p := &Product{
		ID:      7,
		Name:    "Widget Pro",
		Price:   29.99,
		InStock: true,
		Internal: "hidden",
	}

	// --- GetByTag: read a field by its json tag value ---
	fmt.Println("=== GetByTag (json) ===")
	name, err := saferefl.GetByTag[string](p, "json", "name")
	fmt.Printf("  json:name  = %q  (err=%v)\n", name, err)

	price, err := saferefl.GetByTag[float64](p, "json", "price")
	fmt.Printf("  json:price = %.2f  (err=%v)\n", price, err)

	inStock, err := saferefl.GetByTag[bool](p, "json", "in_stock")
	fmt.Printf("  json:in_stock = %v  (err=%v)\n", inStock, err)

	// --- GetByTag: read using db tag ---
	fmt.Println("\n=== GetByTag (db) ===")
	dbName, err := saferefl.GetByTag[string](p, "db", "product_name")
	fmt.Printf("  db:product_name = %q  (err=%v)\n", dbName, err)

	unitPrice, err := saferefl.GetByTag[float64](p, "db", "unit_price")
	fmt.Printf("  db:unit_price   = %.2f  (err=%v)\n", unitPrice, err)

	// --- SetByTag: write a field by tag ---
	fmt.Println("\n=== SetByTag ===")
	_ = saferefl.SetByTag[string](p, "json", "name", "Widget Ultra")
	_ = saferefl.SetByTag[float64](p, "json", "price", 49.99)
	fmt.Printf("  after SetByTag: Name=%q Price=%.2f\n", p.Name, p.Price)

	// --- Tag value not found → FieldNotFoundError ---
	fmt.Println("\n=== Tag not found ===")
	_, err = saferefl.GetByTag[string](p, "json", "no_such_tag")
	if errors.Is(err, saferefl.ErrFieldNotFound) {
		fmt.Printf("  ErrFieldNotFound: %v\n", err)
	}

	// --- Tag key not registered on the struct → FieldNotFoundError ---
	_, err = saferefl.GetByTag[string](p, "xml", "name")
	if errors.Is(err, saferefl.ErrFieldNotFound) {
		fmt.Printf("  ErrFieldNotFound (no xml tags): %v\n", err)
	}

	// --- Type mismatch → TypeMismatchError ---
	fmt.Println("\n=== Type mismatch ===")
	_, err = saferefl.GetByTag[int](p, "json", "name") // name is string, not int
	var tme *saferefl.TypeMismatchError
	if errors.As(err, &tme) {
		fmt.Printf("  TypeMismatchError: field=%s got=%s want=%s\n",
			tme.FieldPath, tme.FieldType, tme.WantType)
	}

	// --- omitempty in tag value: "price,omitempty" → key is "price" ---
	fmt.Println("\n=== omitempty stripped from tag key ===")
	type WithOmit struct {
		Value int `json:"value,omitempty"`
	}
	wo := &WithOmit{Value: 55}
	v, err := saferefl.GetByTag[int](wo, "json", "value") // NOT "value,omitempty"
	fmt.Printf("  json:\"value,omitempty\" → GetByTag key is \"value\": %d  (err=%v)\n", v, err)

	// --- Combining: build a generic JSON-key setter for dynamic configuration ---
	fmt.Println("\n=== Dynamic config setter via SetByTag ===")
	patches := map[string]any{
		"name":     "Widget Max",
		"price":    float64(59.99),
		"in_stock": false,
	}
	p2 := &Product{ID: 8, Name: "original", Price: 1.0, InStock: true}
	for key, val := range patches {
		switch v := val.(type) {
		case string:
			_ = saferefl.SetByTag[string](p2, "json", key, v)
		case float64:
			_ = saferefl.SetByTag[float64](p2, "json", key, v)
		case bool:
			_ = saferefl.SetByTag[bool](p2, "json", key, v)
		}
	}
	fmt.Printf("  p2 after patches: Name=%q Price=%.2f InStock=%v\n",
		p2.Name, p2.Price, p2.InStock)
}
