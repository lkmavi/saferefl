//go:build !reflectx_strict

package unsafelayout

import (
	"errors"
	"log"
	"os"
	"reflect"
	"unsafe"
)

var accelOK bool

func init() {
	accelOK = runSelfTest()
	if !accelOK {
		msg := "[saferefl] Layer 3 unsafe accelerator self-test FAILED — falling back to Layer 2"
		if _, strict := os.LookupEnv("SAFEREFL_STRICT"); strict {
			panic(msg)
		}
		log.Println(msg)
	}
}

// AccelAvailable reports whether the self-test passed on this Go version and arch.
// When false, all Layer 3 operations degrade gracefully to the Layer 2 reflect path.
func AccelAvailable() bool { return accelOK }

// EnableAccel returns nil if the self-test passed, or a descriptive error.
// Call at program startup to confirm Layer 3 is active; absence of a call is safe.
func EnableAccel() error {
	if !accelOK {
		return errors.New("saferefl: Layer 3 unsafe accelerator failed self-test on this Go version/arch")
	}
	return nil
}

// runSelfTest verifies each Layer 3 assumption against the reflect baseline.
// Returns true only if every check passes.
func runSelfTest() bool {
	return selfTestStruct() && selfTestSlice() && selfTestMap()
}

// --- struct test ---

type stFields struct {
	A int64
	B string
	C float64
	D bool
}

func selfTestStruct() bool {
	rt := reflect.TypeOf(stFields{})
	s := stFields{A: 0xDEADBEEF, B: "saferefl", C: 2.718, D: true}
	ptr := unsafe.Pointer(&s)

	for _, name := range []string{"A", "B", "C", "D"} {
		sf, _ := rt.FieldByName(name)
		fPtr := UnsafeFieldPtr(ptr, sf.Offset)
		rfPtr := unsafe.Pointer(reflect.ValueOf(&s).Elem().FieldByName(name).UnsafeAddr())
		if fPtr != rfPtr {
			return false
		}
	}
	return true
}

// --- slice test ---

func selfTestSlice() bool {
	s := []int64{10, 20, 30, 40}
	// reflect.SliceHeader.Data is the first word of the slice header.
	sliceData := *(*unsafe.Pointer)(unsafe.Pointer(&s))
	elemSize := unsafe.Sizeof(int64(0))

	for i, want := range s {
		got := *(*int64)(UnsafeSliceElemPtr(sliceData, i, uintptr(elemSize)))
		if got != want {
			return false
		}
	}
	return true
}

// --- map test ---

func selfTestMap() bool {
	m := map[string]int{"x": 1, "y": 2, "z": 3}
	mapPtr := unsafe.Pointer(reflect.ValueOf(m).Pointer())
	got := activeMapLayout.MapLen(mapPtr)
	return got == len(m)
}
