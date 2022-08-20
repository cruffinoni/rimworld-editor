package algorithm

import (
	"github.com/cruffinoni/rimworld-editor/xml/types/iterator"
)

type Finder[A iterator.SliceIndexer[T], T any] interface {
	FindIf(arr A, f func(T) bool) (T, bool)
}

func FindIf[T any, A iterator.SliceIndexer[T]](arr A, pred func(T) bool) (T, bool) {
	for i := 0; i < arr.Capacity(); i++ {
		if pred(arr.GetFromIndex(i)) {
			return arr.GetFromIndex(i), true
		}
	}
	var t T
	return t, false
}

func Foreach[T any, S iterator.SliceIndexer[T]](arr S, f func(T)) {
	for i := 0; i < arr.Capacity(); i++ {
		f(arr.GetFromIndex(i))
	}
}

func FindInSlice[T any, C Comparable[T], I iterator.SliceIndexer[C]](ref I, comp C) (T, bool) {
	var (
		found = false
		value T
	)
	Foreach(ref, func(v C) {
		if v.Equal(comp) {
			found = true
		}
	})
	return value, found
}
