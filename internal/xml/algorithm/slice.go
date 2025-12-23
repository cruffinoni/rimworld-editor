package algorithm

import (
	"github.com/cruffinoni/rimworld-editor/internal/xml/iterator"
)

func FindInSliceIf[T any](arr iterator.SliceIndexer[T], pred func(T) bool) (T, bool) {
	for i, j := 0, arr.Capacity(); i < j; i++ {
		val := arr.At(i)
		if pred(val) {
			return val, true
		}
	}
	return *new(T), false
}

func SliceForeach[T any](arr iterator.SliceIndexer[T], f func(T)) {
	for i := 0; i < arr.Capacity(); i++ {
		f(arr.At(i))
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
			value = v.Val()
		}
	})
	return value, found
}

//func SetInSliceIf[T any](arr iterator.SliceIndexer[T], pred func(T) bool, value T, attr attributes.Attributes) bool {
//	for i, j := 0, arr.Capacity(); i < j; i++ {
//		val := arr.At(i)
//		if pred(val) {
//			arr.Set(value, attr, i)
//			return true
//		}
//	}
//	return false
//}
