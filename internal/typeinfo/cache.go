package typeinfo

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
)

type cacheEntry struct {
	once sync.Once
	desc atomic.Pointer[TypeDescriptor]
}

var globalCache sync.Map // reflect.Type → *cacheEntry

// testHookBuildDescriptor is called at the start of buildDescriptor if non-nil.
// Used in tests to count invocations and verify single construction.
var testHookBuildDescriptor func()

// TypeDescriptorOf returns the precomputed TypeDescriptor for t.
// Panics if t is not a struct type.
// Safe for concurrent use; the descriptor is built at most once per type.
func TypeDescriptorOf(t reflect.Type) *TypeDescriptor {
	if t.Kind() != reflect.Struct {
		panic(fmt.Sprintf("saferefl: TypeDescriptorOf requires struct type, got %v", t.Kind()))
	}
	// Fast path: type already cached — single atomic load, zero allocation.
	if v, ok := globalCache.Load(t); ok {
		e := v.(*cacheEntry)
		if d := e.desc.Load(); d != nil {
			return d
		}
		// Concurrent first build in progress: wait via once.Do.
		e.once.Do(func() { e.desc.Store(buildDescriptor(t)) })
		return e.desc.Load()
	}
	// Slow path: first time seeing this type — allocate entry and build.
	entry := &cacheEntry{}
	v, _ := globalCache.LoadOrStore(t, entry)
	e := v.(*cacheEntry)
	e.once.Do(func() { e.desc.Store(buildDescriptor(t)) })
	return e.desc.Load()
}

func buildDescriptor(t reflect.Type) *TypeDescriptor {
	if testHookBuildDescriptor != nil {
		testHookBuildDescriptor()
	}
	n := t.NumField()
	desc := &TypeDescriptor{
		Type:         t,
		Kind:         t.Kind(),
		Size:         t.Size(),
		Fields:       make([]FieldMeta, n),
		FieldsByName: make(map[string]*FieldMeta),
		FieldsByTag:  make(map[string]map[string]*FieldMeta),
	}
	for i := range n {
		sf := t.Field(i)
		desc.Fields[i] = FieldMeta{
			Name:      sf.Name,
			Index:     i,
			Offset:    sf.Offset,
			Type:      sf.Type,
			Kind:      sf.Type.Kind(),
			Tag:       sf.Tag,
			Anonymous: sf.Anonymous,
			Exported:  sf.IsExported(),
		}
	}
	collectNamed(desc, t, 0)
	return desc
}

// collectNamed populates FieldsByName and FieldsByTag, recursing into embedded structs.
// baseOffset accumulates the byte offset from the root struct.
// Outer fields shadow inner promoted fields (first write wins).
func collectNamed(desc *TypeDescriptor, t reflect.Type, baseOffset uintptr) {
	n := t.NumField()
	for i := range n {
		sf := t.Field(i)
		offset := baseOffset + sf.Offset

		if _, exists := desc.FieldsByName[sf.Name]; !exists {
			fm := &FieldMeta{
				Name:      sf.Name,
				Index:     i,
				Offset:    offset,
				Type:      sf.Type,
				Kind:      sf.Type.Kind(),
				Tag:       sf.Tag,
				Anonymous: sf.Anonymous,
				Exported:  sf.IsExported(),
			}
			desc.FieldsByName[sf.Name] = fm
			addTagEntries(desc, fm)
		}

		if sf.Anonymous && sf.Type.Kind() == reflect.Struct {
			collectNamed(desc, sf.Type, offset)
		}
	}
}

// addTagEntries parses the struct tag and adds key→name→field entries to FieldsByTag.
// Uses the same tag format as reflect.StructTag: `key:"name[,options]"`.
// First occurrence of a (key, name) pair wins (outer fields shadow embedded ones).
func addTagEntries(desc *TypeDescriptor, fm *FieldMeta) {
	s := string(fm.Tag)
	for s != "" {
		// Skip leading whitespace.
		for s != "" && s[0] == ' ' {
			s = s[1:]
		}
		if s == "" {
			break
		}

		// Parse the key: characters that are not space, colon, quote, or DEL.
		i := 0
		for i < len(s) && s[i] > ' ' && s[i] != ':' && s[i] != '"' && s[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(s) || s[i] != ':' || s[i+1] != '"' {
			break
		}
		key := s[:i]
		s = s[i+2:] // skip key and opening quote

		// Find closing quote.
		end := strings.IndexByte(s, '"')
		if end < 0 {
			break
		}
		value := s[:end]
		s = s[end+1:]

		// The field name is the component before the first comma.
		name, _, _ := strings.Cut(value, ",")
		if name == "" || name == "-" {
			continue
		}
		if desc.FieldsByTag[key] == nil {
			desc.FieldsByTag[key] = make(map[string]*FieldMeta)
		}
		if _, exists := desc.FieldsByTag[key][name]; !exists {
			desc.FieldsByTag[key][name] = fm
		}
	}
}
