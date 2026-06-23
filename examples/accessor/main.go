// Package main demonstrates the Layer 3 Accessor API: pre-bound field access
// with zero-cost field resolution on the hot path.
//
// Use MakeAccessor once at startup (or statement-prepare time), then call
// Get/Set per value in a tight loop with no per-call reflection overhead.
package main

import (
	"fmt"

	"github.com/lkmavi/saferefl"
)

type Order struct {
	ID       int
	Customer string
	Total    float64
	Paid     bool
}

func main() {
	// --- Build accessors once ---
	// MakeAccessor resolves and validates the field path at construction time.
	// The returned Accessor holds a plain byte offset — no reflection at call time.
	idAcc, err := saferefl.MakeAccessor[int](&Order{}, "ID")
	if err != nil {
		panic(err)
	}
	customerAcc, err := saferefl.MakeAccessor[string](&Order{}, "Customer")
	if err != nil {
		panic(err)
	}
	totalAcc, err := saferefl.MakeAccessor[float64](&Order{}, "Total")
	if err != nil {
		panic(err)
	}
	paidAcc, err := saferefl.MakeAccessor[bool](&Order{}, "Paid")
	if err != nil {
		panic(err)
	}

	// --- Use UnsafePtrOf to extract a stable raw pointer once per object ---
	// Safe as long as the object is not moved (Go's GC does not move objects today).
	o := &Order{}
	ptr := saferefl.UnsafePtrOf(o)

	// --- Set via raw pointer (fastest path, zero allocations) ---
	idAcc.Set(ptr, 1001)
	customerAcc.Set(ptr, "Alice")
	totalAcc.Set(ptr, 149.99)
	paidAcc.Set(ptr, true)

	fmt.Printf("Order after Set: %+v\n", *o)

	// --- Get via raw pointer ---
	fmt.Printf("ID=%d  Customer=%q  Total=%.2f  Paid=%v\n",
		idAcc.Get(ptr),
		customerAcc.Get(ptr),
		totalAcc.Get(ptr),
		paidAcc.Get(ptr),
	)

	// --- Convenience methods: GetFrom / SetOn accept any (no UnsafePtrOf required) ---
	// Slightly slower than the raw-pointer path due to interface unwrapping, but still
	// avoids all per-call field resolution.
	o2 := &Order{ID: 2002, Customer: "Bob"}
	v, _ := idAcc.GetFrom(o2)
	fmt.Printf("\nGetFrom: ID=%d\n", v)

	_ = customerAcc.SetOn(o2, "Carol")
	fmt.Printf("SetOn: Customer=%q\n", o2.Customer)

	// --- Batch processing: the hot-path pattern ---
	// Build accessors once, then loop over many objects.
	orders := []*Order{
		{ID: 1, Customer: "X"},
		{ID: 2, Customer: "Y"},
		{ID: 3, Customer: "Z"},
	}
	fmt.Println("\nBatch read:")
	for _, ord := range orders {
		p := saferefl.UnsafePtrOf(ord)
		fmt.Printf("  id=%d  customer=%q\n", idAcc.Get(p), customerAcc.Get(p))
	}

	// --- Type mismatch at construction time, not at runtime ---
	_, err = saferefl.MakeAccessor[int](&Order{}, "Customer") // Customer is string
	fmt.Printf("\nType mismatch at MakeAccessor: %v\n", err)

	// --- Unknown field at construction time ---
	_, err = saferefl.MakeAccessor[string](&Order{}, "Address")
	fmt.Printf("Unknown field at MakeAccessor: %v\n", err)
}
