//go:build reflectx_strict

package typeinfo

import "reflect"

// In reflectx_strict mode all unsafe layout reads are disabled.
// IterPlan is left nil so EachField/ToMap fall back to the reflect path.

func buildIterPlan(_ reflect.Type) []IterEntry { return nil }
