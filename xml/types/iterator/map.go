package iterator

import "log"

type MapIndexer[K, V any] interface {
	SliceIndexer[V]
	GetFromKey(idx int) K
}

type MapIterator[K, V any] struct {
	m   MapIndexer[K, V]
	idx int
	cap int
}

func NewMapIterator[K, V any](v MapIndexer[K, V]) *MapIterator[K, V] {
	return &MapIterator[K, V]{m: v, idx: 0, cap: v.Capacity()}
}

func (mi *MapIterator[K, V]) Next() *MapIterator[K, V] {
	mi.idx++
	if mi.idx > mi.cap {
		log.Panic("iterator overflow")
	} else if mi.idx == mi.cap {
		return nil
	}
	return &MapIterator[K, V]{m: mi.m, idx: mi.idx, cap: mi.cap}
}

func (mi *MapIterator[K, V]) Prev() *MapIterator[K, V] {
	mi.idx--
	if mi.idx < 0 {
		log.Panic("iterator underflow")
	}
	return &MapIterator[K, V]{m: mi.m, idx: mi.idx, cap: mi.cap}
}

func (mi *MapIterator[K, V]) Key() K {
	if mi.HasNext() {
		return mi.m.GetFromKey(mi.idx)
	}
	panic("iterator overflow")
}

func (mi *MapIterator[K, V]) Value() V {
	if mi.HasNext() {
		return mi.m.GetFromIndex(mi.idx)
	}
	panic("iterator overflow")
}

func (mi *MapIterator[K, V]) HasNext() bool {
	return mi.idx < mi.cap
}

func (mi *MapIterator[K, V]) Capacity() int {
	return mi.m.Capacity()
}

func (mi *MapIterator[K, V]) Index() int {
	return mi.idx
}
