package types

import (
	"fmt"
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/xml/saver"
	"github.com/cruffinoni/rimworld-editor/xml/saver/file"
	"github.com/cruffinoni/rimworld-editor/xml/types/iterator"
	"github.com/cruffinoni/rimworld-editor/xml/unmarshal"
	"reflect"
	"strings"
)

type sliceData[T any] struct {
	data T
	// str is a string representation of data
	str  string
	attr attributes.Attributes
	fmt.Stringer
	tag string
}

func (s *sliceData[T]) Assign(e *xml.Element) error {
	var (
		err   error
		tKind = reflect.TypeOf(s.data).Kind()
	)
	if tKind == reflect.Ptr {
		err = unmarshal.Element(e, s.data)
	} else {
		err = unmarshal.Element(e, &s.data)
	}
	if err != nil {
		return err
	}
	s.UpdateStringRepresentation(s.data)
	return err
}

func (s *sliceData[T]) GetPath() string {
	return ""
}

func (s *sliceData[T]) SetAttributes(attributes attributes.Attributes) {
	s.attr = attributes
}

func (s *sliceData[T]) GetAttributes() attributes.Attributes {
	return s.attr
}

func (s *sliceData[T]) UpdateStringRepresentation(v T) {
	t := reflect.TypeOf(v)
	// We check if the type T implements the interface fmt.Stringer and has a
	// custom String() method.
	if ok := t.Implements(reflect.TypeOf(new(fmt.Stringer)).Elem()); ok {
		var m reflect.Method
		// If it's the case, we get the method String() of the type T and
		// call it.
		if m, ok = t.MethodByName("String"); ok {
			s.str = m.Func.Call([]reflect.Value{reflect.ValueOf(v)})[0].String()
		}
	} else {
		// Otherwise we use a basic string representation.
		s.str = fmt.Sprintf("%v", v)
	}
}

func (s sliceData[T]) String() string {
	return s.str
}

// Slice is a slice of data that is represented by sliceData.
// The main difference with a normal slice is that it can get and set attributes.
// Like a normal slice, the slice is a slice of T and you can iterate over it.
type Slice[T any] struct {
	data         []sliceData[T]
	attr         attributes.Attributes
	repeatingTag string
	cap          int
	iterator.SliceIndexer[T]
	saver.Transformer
	fmt.Stringer
}

func (s Slice[T]) TransformToXML(b *saver.Buffer) error {
	if s.repeatingTag == "" {
		return nil
	}
	lastElement := s.cap - 1
	for i, v := range s.data {
		if err := file.Save(v.data, b, s.repeatingTag); err != nil {
			return err
		}
		if i != lastElement {
			b.Write([]byte("\n"))
		}
	}
	return nil
}

func (s *Slice[T]) Capacity() int {
	return s.cap
}

func (s *Slice[T]) GetFromIndex(idx int) T {
	for i, d := range s.data {
		if i == idx {
			return d.data
		}
	}
	return *new(T)
}

func (s *Slice[T]) Assign(e *xml.Element) error {
	s.data = make([]sliceData[T], 0)
	defer func() {
		s.cap = len(s.data)
	}()
	n := e
	if n == nil {
		return nil
	}
	s.repeatingTag = n.GetName()
	for n != nil {
		d := sliceData[T]{
			tag:  n.GetName(),
			data: reflect.New(reflect.TypeOf(*new(T)).Elem()).Interface().(T),
		}
		// The child element inherits the attributes of the parent element
		// because we don't unmarshal the element directly but the children
		// since it's a slice.
		if n.Child != nil {
			n.Child.Attr = n.Attr
		}
		if err := unmarshal.Element(n.Child, &d); err != nil {
			return err
		}
		s.data = append(s.data, d)
		n = n.Next
	}
	return nil
}

func (s Slice[T]) String() string {
	b := strings.Builder{}
	b.WriteString("[")
	for i, d := range s.data {
		if i > 0 {
			b.WriteString(", " + d.String())
		} else {
			b.WriteString(d.String())
		}
	}
	b.WriteString("]")
	return b.String()
}

func (s *Slice[T]) GetPath() string {
	// We avoid to use GetPath in generics functions like Slice or Map
	return ""
}

func (s *Slice[T]) SetAttributes(attributes attributes.Attributes) {
	s.attr = attributes
}

func (s *Slice[T]) GetAttributes() attributes.Attributes {
	return s.attr
}
