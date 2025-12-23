package algorithm

import (
	"github.com/cruffinoni/rimworld-editor/internal/xml/collection"
	"github.com/cruffinoni/rimworld-editor/internal/xml/iterator"
)

func FindInMapIf[A iterator.MapIndexer[K, V], K comparable, V any](arr A, pred func(*collection.Pair[K, V]) bool) (*collection.Pair[K, V], bool) {
	for i := 0; i < arr.Capacity(); i++ {
		p := &collection.Pair[K, V]{
			Key:   arr.GetKeyFromIndex(i),
			Value: arr.GetFromIndex(i),
		}
		if pred(p) {
			return p, true
		}
	}
	t := &collection.Pair[K, V]{}
	return t, false
}

func MapForeach[S iterator.MapIndexer[K, V], K comparable, V any](arr S, f func(*collection.Pair[K, V])) {
	for i := 0; i < arr.Capacity(); i++ {
		f(&collection.Pair[K, V]{
			Key:   arr.GetKeyFromIndex(i),
			Value: arr.GetFromIndex(i),
		})
	}
}

// FindInMap finds the first element in the slice that satisfies the predicate.
// I is the iterator to use to iterate over the slice.
func FindInMap[I iterator.MapIndexer[K, V], K comparable, V any](ref I, comp *collection.Pair[K, V]) (*collection.Pair[K, V], bool) {
	var (
		found = false
		value *collection.Pair[K, V]
	)
	MapForeach(ref, func(v *collection.Pair[K, V]) {
		if !found && v.Equal(comp) {
			found = true
			value = v
		}
	})
	return value, found
}
