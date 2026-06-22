package realistic

import "github.com/lkmavi/saferefl"

func mustMakeAccessor[T any](obj any, path string) saferefl.Accessor[T] {
	acc, err := saferefl.MakeAccessor[T](obj, path)
	if err != nil {
		panic(err)
	}
	return acc
}
