package typeinfo

import (
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
)

func TestPtrCacheLoadMiss(t *testing.T) {
	// Key 1 is never stored by any real code (abi.Type pointers are never that small).
	const key = uintptr(1)
	if _, ok := PtrCacheLoad(key); ok {
		t.Error("expected cache miss for synthetic key")
	}
}

func TestPtrCacheLoadStore(t *testing.T) {
	desc := TypeDescriptorOf(reflect.TypeOf(basicStruct{}))
	const key = uintptr(2) // synthetic key, never a real abi.Type address
	PtrCacheStore(key, desc)
	got, ok := PtrCacheLoad(key)
	if !ok {
		t.Fatal("expected cache hit after store")
	}
	if got != desc {
		t.Error("PtrCacheLoad returned wrong descriptor")
	}
}

func TestTypeDescriptorOf_panicsOnNonStruct(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for non-struct type, got none")
		}
	}()
	TypeDescriptorOf(reflect.TypeOf(42))
}

func TestTypeDescriptorOf_returnsSamePointer(t *testing.T) {
	rt := reflect.TypeOf(basicStruct{})
	a := TypeDescriptorOf(rt)
	b := TypeDescriptorOf(rt)
	if a != b {
		t.Error("expected same *TypeDescriptor pointer on repeated calls")
	}
}

func TestCacheConcurrency(t *testing.T) {
	// Use a locally-defined type so only this test touches its cache entry.
	type concurrentTarget struct{}
	rt := reflect.TypeOf(concurrentTarget{})

	// Clear any entry from a previous run (e.g. -count=2).
	globalCache.Delete(rt)

	var count atomic.Int64
	prev := testHookBuildDescriptor
	testHookBuildDescriptor = func() { count.Add(1) }
	defer func() { testHookBuildDescriptor = prev }()

	const n = 100
	results := make([]*TypeDescriptor, n)
	var wg sync.WaitGroup
	for i := range n {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			results[idx] = TypeDescriptorOf(rt)
		}(i)
	}
	wg.Wait()

	for i, desc := range results {
		if desc != results[0] {
			t.Errorf("goroutine %d got a different *TypeDescriptor", i)
		}
	}

	if got := count.Load(); got != 1 {
		t.Errorf("buildDescriptor called %d times, want exactly 1", got)
	}
}
