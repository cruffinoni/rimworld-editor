package xml

type ListIndex struct {
	Valid bool
	Idx   int
}

func (l *ListIndex) FromIndex(idx int) {
	l.Idx = idx
	l.Valid = true
}

func (l *ListIndex) ToIndex() int {
	return l.Idx
}

func (l *ListIndex) IsValid() bool {
	return l.Valid
}

func (l *ListIndex) Invalidate() {
	l.Valid = false
	l.Idx = 0
}

func (l *ListIndex) Increment() {
	l.Idx++
	l.Valid = true
}

func (l *ListIndex) Decrement() {
	l.Idx--
	l.Valid = true
}
