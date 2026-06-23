// Package main demonstrates the Layer 3 low-level primitives:
// UnsafeSliceAt, MapLenFast, AccelAvailable, and EnableAccel.
//
// These functions bypass the reflect package entirely — they read runtime
// struct layout directly from memory. Call EnableAccel once at startup to
// verify that the layout assumptions hold on the current Go runtime.
package main

import (
	"fmt"

	"github.com/lkmavi/saferefl"
)

func main() {
	// --- EnableAccel: optional startup self-test ---
	// EnableAccel runs a brief sanity check and enables the fast paths.
	// If it returns an error the library falls back to safe alternatives
	// automatically, so calling it is recommended but not required.
	if err := saferefl.EnableAccel(); err != nil {
		fmt.Printf("accel unavailable, using fallback: %v\n", err)
	} else {
		fmt.Println("Layer 3 accelerator: OK")
	}

	// AccelAvailable reports whether EnableAccel succeeded (or the
	// built-in self-test passed without an explicit call).
	fmt.Printf("AccelAvailable: %v\n\n", saferefl.AccelAvailable())

	// --- UnsafeSliceAt: direct pointer to a slice element ---
	// Returns *T without any bounds-check overhead from the reflect package.
	// The caller is still responsible for staying within slice bounds.
	scores := []float64{9.1, 8.5, 7.8, 9.9, 6.0}

	// Read via pointer — same as &scores[2] but bypasses reflect overhead.
	p := saferefl.UnsafeSliceAt(scores, 2)
	fmt.Printf("scores[2] via UnsafeSliceAt: %.1f\n", *p)

	// Write through the returned pointer — modifies the original slice.
	*saferefl.UnsafeSliceAt(scores, 4) = 10.0
	fmt.Printf("scores after write: %v\n\n", scores)

	// Typical pattern: scan a large slice without allocating reflect.Values.
	sum := 0.0
	for i := range scores {
		sum += *saferefl.UnsafeSliceAt(scores, i)
	}
	fmt.Printf("sum (via UnsafeSliceAt): %.1f\n\n", sum)

	// --- MapLenFast: read map length without reflect ---
	// Falls back to len(m) under the reflectx_strict build tag or when
	// the runtime layout doesn't match expectations.
	inventory := map[string]int{
		"apple":  120,
		"banana": 45,
		"cherry": 300,
	}
	fmt.Printf("MapLenFast: %d items\n", saferefl.MapLenFast(inventory))

	// Works with any key/value types.
	ids := map[int]bool{1: true, 2: true, 3: false}
	fmt.Printf("MapLenFast (int→bool): %d items\n\n", saferefl.MapLenFast(ids))

	// --- Combining with Accessor in a hot loop ---
	// Build accessor once, extract pointer per object, loop at full speed.
	type Record struct {
		ID    int
		Score float64
	}

	scoreAcc, _ := saferefl.MakeAccessor[float64](&Record{}, "Score")
	records := make([]*Record, 5)
	for i := range records {
		records[i] = &Record{ID: i + 1, Score: float64(i) * 1.5}
	}

	fmt.Println("Records via Accessor + UnsafePtrOf:")
	for _, r := range records {
		ptr := saferefl.UnsafePtrOf(r)
		fmt.Printf("  ID=%d  Score=%.1f\n", r.ID, scoreAcc.Get(ptr))
	}
}
