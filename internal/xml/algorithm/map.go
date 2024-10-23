package algorithm

import (
	"github.com/cruffinoni/rimworld-editor/internal/xml/types"
	"github.com/cruffinoni/rimworld-editor/internal/xml/types/iterator"
)

func FindInMapIf[A iterator.MapIndexer[K, V], K comparable, V any](arr A, pred func(*types.Pair[K, V]) bool) (*types.Pair[K, V], bool) {
	for i := 0; i < arr.Capacity(); i++ {
		p := &types.Pair[K, V]{
			Key:   arr.GetKeyFromIndex(i),
			Value: arr.GetFromIndex(i),
		}
		if pred(p) {
			return p, true
		}
	}
	t := &types.Pair[K, V]{}
	return t, false
}

func MapForeach[S iterator.MapIndexer[K, V], K comparable, V any](arr S, f func(*types.Pair[K, V])) {
	for i := 0; i < arr.Capacity(); i++ {
		f(&types.Pair[K, V]{
			Key:   arr.GetKeyFromIndex(i),
			Value: arr.GetFromIndex(i),
		})
	}
}

// FindInMap finds the first element in the slice that satisfies the predicate.
// I is the iterator to use to iterate over the slice.
func FindInMap[I iterator.MapIndexer[K, V], K comparable, V any](ref I, comp *types.Pair[K, V]) (*types.Pair[K, V], bool) {
	var (
		found = false
		value *types.Pair[K, V]
	)
	MapForeach(ref, func(v *types.Pair[K, V]) {
		if !found && v.Equal(comp) {
			found = true
			value = v
		}
	})
	return value, found
}
