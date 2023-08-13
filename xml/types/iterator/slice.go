package iterator

import "github.com/cruffinoni/rimworld-editor/xml/attributes"

type SliceIndexer[V any] interface {
	GetFromIndex(idx int) V
	Capacity() int
	Set(value V, attr attributes.Attributes, idx int)
	Add(value V, attr attributes.Attributes)
	Remove(idx int)
}

type SliceIterator[V any] struct {
	m   SliceIndexer[V]
	idx int
	cap int
}

func NewSliceIterator[V any](v SliceIndexer[V]) *SliceIterator[V] {
	return &SliceIterator[V]{m: v, idx: 0, cap: v.Capacity()}
}

func (si *SliceIterator[V]) Next() *SliceIterator[V] {
	si.idx++
	if si.idx > si.cap {
		panic("iterator overflow")
	}
	return &SliceIterator[V]{m: si.m, idx: si.idx, cap: si.cap}
}

func (si *SliceIterator[V]) Prev() *SliceIterator[V] {
	si.idx--
	if si.idx < 0 {
		panic("iterator underflow")
	}
	return &SliceIterator[V]{m: si.m, idx: si.idx, cap: si.cap}
}

func (si *SliceIterator[V]) Value() V {
	if si.HasNext() {
		return si.m.GetFromIndex(si.idx)
	}
	panic("iterator overflow")
}

func (si *SliceIterator[V]) HasNext() bool {
	return si.idx < si.cap
}

func (si *SliceIterator[V]) Capacity() int {
	return si.m.Capacity()
}

func (si *SliceIterator[V]) Index() int {
	return si.idx
}
