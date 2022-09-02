package algorithm

import (
	"github.com/cruffinoni/rimworld-editor/xml/types/iterator"
)

func FindInSliceIf[T any](arr iterator.SliceIndexer[T], pred func(T) bool) (T, bool) {
	for i := 0; i < arr.Capacity(); i++ {
		if pred(arr.GetFromIndex(i)) {
			return arr.GetFromIndex(i), true
		}
	}
	var t T
	return t, false
}

func SliceForeach[T any](arr iterator.SliceIndexer[T], f func(T)) {
	for i := 0; i < arr.Capacity(); i++ {
		f(arr.GetFromIndex(i))
	}
}

// FindInSlice finds the first element in the slice that satisfies the predicate.
// T is the type of the element to find.
// C is the value to compare
// I is the iterator to use to iterate over the slice.
func FindInSlice[C Comparable[T], T any](ref iterator.SliceIndexer[C], comp T) (T, bool) {
	var (
		found = false
		value T
	)
	SliceForeach(ref, func(v C) {
		if !found && v.Equal(comp) {
			found = true
			value = v.Value()
		}
	})
	return value, found
}
