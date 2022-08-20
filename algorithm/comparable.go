package algorithm

type Comparable[T any] interface {
	~string | ~int | ~float32
	Less(rhs T)
	Greater(rhs T)
	Equal(rhs T)
}

type Findable[T any] interface {
}
