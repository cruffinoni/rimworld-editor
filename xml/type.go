package xml

import "fmt"

type EmbeddedPrimaryType[T comparable] struct {
	data T
	// str is the string representation of the data.
	str  string
	attr Attributes
}

func (pt *EmbeddedPrimaryType[T]) Assign(e *Element) error {
	if v, ok := e.Data.GetData().(T); ok {
		pt.data = v
		pt.str = fmt.Sprintf("%v", v)
		return nil
	} else {
		return fmt.Errorf("xml: cannot assign %T to %T", e.Data.GetData(), pt)
	}
}

func (pt *EmbeddedPrimaryType[T]) GetPath() string {
	return ""
}

func (pt *EmbeddedPrimaryType[T]) SetAttributes(attributes Attributes) {
	pt.attr = attributes
}

func (pt *EmbeddedPrimaryType[T]) GetAttributes() Attributes {
	return pt.attr
}

func (pt *EmbeddedPrimaryType[T]) String() string {
	return pt.str
}
