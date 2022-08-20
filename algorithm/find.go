package algorithm

import (
	"github.com/cruffinoni/rimworld-editor/xml/types/iterator"
	"reflect"
)

type Finder[A iterator.SliceIndexer[T], T any] interface {
	FindIf(arr A, f func(T) bool) (T, bool)
}

func newFinder[A iterator.SliceIndexer[T], T any](f func(T) bool) Finder[A, T] {
	return &finder[A, T]{f: f}
}

type finder[A iterator.SliceIndexer[T], T any] struct {
	f func(T) bool
}

func (f *finder[A, T]) FindIf(ref A, pred func(T) bool) (T, bool) {
	for i := 0; i < ref.Capacity(); i++ {
		if pred(ref.GetFromIndex(i)) {
			return ref.GetFromIndex(i), true
		}
	}
	var t T
	return t, false
}

func FindIf[A iterator.SliceIndexer[T], T any](arr A, f func(T) bool) (T, bool) {
	return newFinder[A, T](f).FindIf(arr, f)
}

func FindInSlice[K any, T Comparable[K]](ref T, comp K) (K, bool) {
	for k := range ref {
		if comp.Equal(reflect.ValueOf(k), reflect.ValueOf(comp)) {
			return k, true
		}
	}
	var k K
	return k, false
}
