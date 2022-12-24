package types

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/xml/saver"
	"github.com/cruffinoni/rimworld-editor/xml/saver/file"
	"github.com/cruffinoni/rimworld-editor/xml/types/iterator"
	"github.com/cruffinoni/rimworld-editor/xml/unmarshal"
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
	log.Println("Assign sliceData")
	var (
		err   error
		tKind = reflect.TypeOf(s.data).Kind()
	)
	if tKind == reflect.Ptr {
		err = unmarshal.Element(e, s.data)
	} else {
		err = unmarshal.Element(e, &s.data)
	}
	// log.Printf("Data '%v' / %T", s.data, s.data)
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
		s.str = fmt.Sprintf("'%v'", v)
	}
}

func (s *sliceData[T]) String() string {
	log.Println("String called")
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
		log.Print("Slice.TransformToXML: No repeating tag specified.")
		return nil
	}
	lastElement := s.cap - 1
	for i, v := range s.data {
		log.Printf("Slice.TransformToXML: %v at %v // %v", v.data, i, s.repeatingTag)
		if err := file.Save(v.data, b, s.repeatingTag); err != nil {
			return err
		}
		if i != lastElement {
			b.Write([]byte("\n"))
		}
	}
	return nil
}

func (s Slice[T]) Capacity() int {
	return s.cap
}

func (s Slice[T]) GetFromIndex(idx int) T {
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
	log.Printf("Slice.Assign: Repeating tag: %v", s.repeatingTag)
	for n != nil {
		d := sliceData[T]{
			tag: n.GetName(),
		}
		// Set d.data to zero depending on the type of T. Either a pointer or a
		// value.
		switch tType := reflect.TypeOf(d.data).Kind(); tType {
		case reflect.Ptr, reflect.Interface, reflect.Struct, reflect.Map, reflect.Slice:
			d.data = reflect.New(reflect.TypeOf(*new(T)).Elem()).Interface().(T)
		case reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
			d.data = zero[T]()
		}
		// The child element inherits the attributes of the parent element
		// because we don't unmarshal the element directly but the children
		// since it's a slice.
		if n.Child != nil {
			n.Child.Attr = n.Attr
			if err := unmarshal.Element(n.Child, &d); err != nil {
				return err
			}
		} else {
			// TODO: Might be reworked to something more elegant.
			if err := unmarshal.Element(n, &d); err != nil {
				return err
			}
		}
		log.Printf("Slice.Assign: Data: %v", n.Data.GetData())
		log.Printf("Data of d: '%v' (%T)", d.data, d.data)
		s.data = append(s.data, d)
		n = n.Next
	}
	log.Println("Slice.Assign: end")
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
