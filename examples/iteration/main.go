// Package main demonstrates EachField and MapForEach —
// field iteration over structs and early-exit iteration over maps.
package main

import (
	"fmt"
	"strings"

	"github.com/lkmavi/saferefl"
)

type Event struct {
	ID       int
	Title    string
	Speaker  string
	Capacity int
	Online   bool
}

func main() {
	e := &Event{ID: 42, Title: "Go Conf", Speaker: "Alice", Capacity: 200, Online: true}

	// --- EachField: iterate all exported fields in declaration order ---
	fmt.Println("=== EachField ===")
	_ = saferefl.EachField(e, func(name string, val any) bool {
		fmt.Printf("  %-10s = %v\n", name, val)
		return true // return false to stop early
	})

	// --- EachField: stop early on first matching field ---
	fmt.Println("\n=== EachField early stop (find first string) ===")
	_ = saferefl.EachField(e, func(name string, val any) bool {
		if s, ok := val.(string); ok {
			fmt.Printf("  first string field: %s = %q\n", name, s)
			return false // stop
		}
		return true
	})

	// --- EachField: generic "non-zero" validator ---
	fmt.Println("\n=== Validate non-zero fields ===")
	blank := &Event{ID: 1, Title: "", Speaker: "Bob"} // Title is empty
	var missing []string
	_ = saferefl.EachField(blank, func(name string, val any) bool {
		switch v := val.(type) {
		case string:
			if v == "" {
				missing = append(missing, name)
			}
		case int:
			if v == 0 {
				missing = append(missing, name)
			}
		}
		return true
	})
	if len(missing) > 0 {
		fmt.Printf("  missing required fields: %s\n", strings.Join(missing, ", "))
	}

	// --- EachField: embedded structs — promoted fields appear at the top level ---
	type Location struct {
		City    string
		Country string
	}
	type Conference struct {
		Location // value-embedded: City and Country are promoted
		Name     string
		Year     int
	}
	c := &Conference{
		Location: Location{City: "Berlin", Country: "DE"},
		Name:     "GopherCon EU",
		Year:     2025,
	}
	fmt.Println("\n=== EachField with embedded struct ===")
	_ = saferefl.EachField(c, func(name string, val any) bool {
		fmt.Printf("  %-10s = %v\n", name, val)
		return true
	})

	// --- MapForEach: typed iteration with early-exit ---
	scores := map[string]int{
		"Alice": 95,
		"Bob":   72,
		"Carol": 88,
		"Dave":  60,
	}

	fmt.Println("\n=== MapForEach: collect scores ≥ 80 ===")
	var top []string
	saferefl.MapForEach(scores, func(name string, score int) bool {
		if score >= 80 {
			top = append(top, fmt.Sprintf("%s(%d)", name, score))
		}
		return true
	})
	fmt.Printf("  top scorers: %v\n", top)

	// --- MapForEach: find first match (early exit) ---
	fmt.Println("\n=== MapForEach early exit (first score < 70) ===")
	saferefl.MapForEach(scores, func(name string, score int) bool {
		if score < 70 {
			fmt.Printf("  first low scorer found: %s = %d\n", name, score)
			return false // stop iteration
		}
		return true
	})
}
