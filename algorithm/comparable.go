package algorithm

import "github.com/cruffinoni/rimworld-editor/xml/types/iterator"

type Comparable[T any] interface {
	Less(rhs T) bool
	Greater(rhs T) bool
	Equal(rhs T) bool
	Val() T
}

type Findable[T any] interface {
	Comparable[T]
	iterator.SliceIndexer[T]
}
