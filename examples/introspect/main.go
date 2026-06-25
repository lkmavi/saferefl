// Package main demonstrates KindOf and IsNil — fast type introspection
// without full reflect.Type construction.
package main

import (
	"fmt"
	"reflect"

	"github.com/lkmavi/saferefl"
)

type Config struct {
	Host string
	Port int
}

func describe(label string, v any) {
	fmt.Printf("  %-30s  kind=%-12v  isNil=%v\n",
		label, saferefl.KindOf(v), saferefl.IsNil(v))
}

func main() { //nolint:funlen
	var nilPtr *Config
	var nilMap map[string]int
	var nilChan chan int
	var nilSlice []string
	var nilFunc func()

	nonNilPtr := &Config{Host: "localhost", Port: 8080}
	nonNilMap := map[string]int{"a": 1}
	nonNilSlice := []string{"x"}
	nonNilChan := make(chan int, 1)

	// --- KindOf: reflects.Kind without reflect.TypeOf ---
	fmt.Println("=== KindOf ===")
	describe("nil", nil)
	describe("42 (int)", 42)
	describe(`"hello" (string)`, "hello")
	describe("3.14 (float64)", 3.14)
	describe("true (bool)", true)
	describe("(*Config)(nil)", nilPtr)
	describe("&Config{...}", nonNilPtr)
	describe("map[string]int(nil)", nilMap)
	describe("map[string]int{...}", nonNilMap)
	describe("chan int (nil)", nilChan)
	describe("chan int (non-nil)", nonNilChan)
	describe("[]string(nil)", nilSlice)
	describe("[]string{...}", nonNilSlice)
	describe("Config{} (struct)", Config{})

	// KindOf agrees with reflect.TypeOf for non-nil values
	fmt.Println("\n=== KindOf matches reflect.TypeOf ===")
	vals := []any{42, "x", 3.14, true, nonNilPtr, nonNilMap, nonNilSlice}
	for _, v := range vals {
		want := reflect.TypeOf(v).Kind()
		got := saferefl.KindOf(v)
		match := got == want
		fmt.Printf("  %-20T  saferefl=%v  reflect=%v  match=%v\n", v, got, want, match)
	}

	// --- IsNil: covers all nilable types ---
	fmt.Println("\n=== IsNil ===")
	describe("nil interface", nil)
	describe("(*Config)(nil)", nilPtr)
	describe("&Config{...}", nonNilPtr)
	describe("map (nil)", nilMap)
	describe("map (non-nil)", nonNilMap)
	describe("chan (nil)", nilChan)
	describe("chan (non-nil)", nonNilChan)
	describe("[]string (nil)", nilSlice)
	describe("[]string{...}", nonNilSlice)
	describe("func (nil)", nilFunc)

	// Non-nilable types always return false
	fmt.Println("\n=== IsNil on non-nilable types ===")
	describe("42 (int)", 42)
	describe(`"hello" (string)`, "hello")
	describe("Config{} (struct)", Config{})

	// --- Practical use: nil-safe value dispatcher ---
	fmt.Println("\n=== Practical: nil-safe handler dispatch ===")
	inputs := []any{nil, nonNilPtr, nilPtr, nonNilMap, nilSlice, "hello", 0}
	for _, v := range inputs {
		k := saferefl.KindOf(v)
		isNil := saferefl.IsNil(v)
		switch {
		case v == nil || isNil:
			fmt.Printf("  nil (%v)\n", k)
		case k == reflect.Pointer:
			fmt.Printf("  pointer to %v\n", reflect.TypeOf(v).Elem())
		case k == reflect.Map:
			fmt.Printf("  map with %d entries\n", saferefl.MapLenFast(nonNilMap))
		case k == reflect.String:
			fmt.Printf("  string: %q\n", v)
		case k == reflect.Int:
			fmt.Printf("  int: %d\n", v)
		default:
			fmt.Printf("  other kind: %v\n", k)
		}
	}
}
