// Package main demonstrates ToMap, ToMapByTag, and FromMap —
// converting structs to/from maps for JSON-like serialization, API responses,
// and round-trip data transformation.
package main

import (
	"encoding/json"
	"fmt"

	"github.com/lkmavi/saferefl"
)

type Article struct {
	ID      int     `json:"id"      db:"article_id"`
	Title   string  `json:"title"   db:"title"`
	Author  string  `json:"author"  db:"author_name"`
	Score   float64 `json:"score"   db:"score"`
	Draft   bool    `json:"draft"   db:"-"`
	private string  // unexported — excluded from all maps
}

func main() {
	a := &Article{
		ID:      101,
		Title:   "Go Internals",
		Author:  "Alice",
		Score:   9.2,
		Draft:   false,
		private: "hidden",
	}

	// --- ToMap: field names as keys ---
	fmt.Println("=== ToMap (field names) ===")
	m, err := saferefl.ToMap(a)
	if err != nil {
		panic(err)
	}
	for k, v := range m {
		fmt.Printf("  %-8s = %v\n", k, v)
	}

	// --- ToMapByTag: use json tag names as keys ---
	fmt.Println("\n=== ToMapByTag(json) ===")
	jm, err := saferefl.ToMapByTag(a, "json")
	if err != nil {
		panic(err)
	}
	for k, v := range jm {
		fmt.Printf("  %-8s = %v\n", k, v)
	}

	// --- ToMapByTag: use db tag names (Draft has db:"-" — skipped) ---
	fmt.Println("\n=== ToMapByTag(db) — Draft omitted (tagged \"-\") ===")
	dm, err := saferefl.ToMapByTag(a, "db")
	if err != nil {
		panic(err)
	}
	for k, v := range dm {
		fmt.Printf("  %-14s = %v\n", k, v)
	}

	// --- Round-trip: struct → JSON → map[string]any → struct via FromMap ---
	// This mirrors the common pattern after json.Unmarshal into map[string]any.
	fmt.Println("\n=== JSON round-trip via FromMap ===")
	raw, _ := json.Marshal(jm)
	fmt.Printf("  JSON: %s\n", raw)

	var decoded map[string]any
	_ = json.Unmarshal(raw, &decoded)
	// json.Unmarshal produces float64 for all numbers — FromMap converts automatically.

	dst := &Article{}
	if err = saferefl.FromMap(decoded, dst); err != nil {
		panic(err)
	}
	// FromMap uses field names (not tag names), so map keys must match.
	// For tag-based population, combine ToMapByTag (reading) with explicit key mapping.
	fmt.Printf("  Decoded ID via FromMap on field names: not set (keys are json names)\n")

	// Correct round-trip: use field-name map for FromMap.
	dst2 := &Article{}
	if err = saferefl.FromMap(m, dst2); err != nil {
		panic(err)
	}
	fmt.Printf("  Decoded: ID=%d Title=%q Author=%q Score=%.1f Draft=%v\n",
		dst2.ID, dst2.Title, dst2.Author, dst2.Score, dst2.Draft)

	// --- FromMap: type conversion (float64 from JSON → int) ---
	fmt.Println("\n=== FromMap with float64→int conversion ===")
	fromJSON := map[string]any{
		"ID":    float64(202), // json.Unmarshal always gives float64
		"Title": "Quick Tips",
	}
	dst3 := &Article{}
	_ = saferefl.FromMap(fromJSON, dst3)
	fmt.Printf("  ID=%d Title=%q\n", dst3.ID, dst3.Title)

	// --- FromMap: unknown keys are silently skipped ---
	fmt.Println("\n=== FromMap: unknown keys skipped ===")
	dst4 := &Article{Title: "original"}
	_ = saferefl.FromMap(map[string]any{
		"Title":      "updated",
		"NoSuchKey":  "ignored",
		"ExtraField": 999,
	}, dst4)
	fmt.Printf("  Title=%q (NoSuchKey and ExtraField silently skipped)\n", dst4.Title)

	// --- Promoted fields from embedded structs appear in the map ---
	fmt.Println("\n=== ToMap with promoted embedded fields ===")
	type Address struct {
		City    string
		Country string
	}
	type Employee struct {
		Address        // value-embedded: City and Country promoted to top level
		Name    string `json:"name"`
		Dept    string `json:"dept"`
	}
	emp := &Employee{
		Address: Address{City: "Berlin", Country: "DE"},
		Name:    "Bob",
		Dept:    "Engineering",
	}
	em, _ := saferefl.ToMap(emp)
	for k, v := range em {
		fmt.Printf("  %-10s = %v\n", k, v)
	}
}
