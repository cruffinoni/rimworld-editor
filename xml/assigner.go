package xml

type Assigner interface {
	Assign(e *Element) error
	GetPath() string
}
