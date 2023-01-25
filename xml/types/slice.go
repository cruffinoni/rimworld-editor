package types

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/cruffinoni/rimworld-editor/helper"
	"github.com/cruffinoni/rimworld-editor/xml/saver/xmlFile"

	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/xml/saver"
	"github.com/cruffinoni/rimworld-editor/xml/unmarshal"
)

type sliceData[T any] struct {
	data T
	// str is a string representation of data
	str  string
	attr attributes.Attributes
	fmt.Stringer
	tag  string
	kind reflect.Kind
}

func (s *sliceData[T]) Assign(e *xml.Element) error {
	var err error
	s.kind = reflect.TypeOf(s.data).Kind()
	if s.kind == reflect.Ptr {
		err = unmarshal.Element(e, s.data)
	} else if helper.IsReflectPrimaryType(s.kind) {
		switch s.kind {
		case reflect.String:
			s.data = castTemplate[T](e.Data.GetString())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			s.data = castTemplate[T](e.Data.GetInt64())
		case reflect.Bool:
			s.data = castTemplate[T](e.Data.GetBool())
		case reflect.Float32, reflect.Float64:
			s.data = castTemplate[T](e.Data.GetFloat64())
		default:
			return fmt.Errorf("sliceData.Assign: can't assign primary type %T to %T", e.Data.GetData(), s.data)
		}
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
		s.str = fmt.Sprintf("'%+v'", v)
	}
}

func (s *sliceData[T]) String() string {
	return s.str
}

func (s *sliceData[T]) GetXMLTag() []byte {
	return nil
}

func (s *sliceData[T]) TransformToXML(b *saver.Buffer) error {
	//log.Printf("sliceData.TransformToXML => %v", b.GetDepth())
	b.OpenTag(s.tag, s.attr)
	if err := xmlFile.Save(s.data, b, ""); err != nil {
		return err
	}
	if helper.IsReflectPrimaryType(s.kind) {
		if s.kind == reflect.String && strings.Contains(reflect.ValueOf(s.data).String(), "\n") {
			b.CloseTagWithIndent(s.tag)
		}
		b.CloseTag(s.tag)
		return nil
	}
	b.CloseTagWithIndent(s.tag)
	return nil
}

// Slice is a slice of data that is represented by sliceData.
// The main difference with a normal slice is that it can get and set attributes.
// Like a normal slice, the slice is a slice of T and you can iterate over it.
type Slice[T any] struct {
	data         []sliceData[T]
	attr         attributes.Attributes
	repeatingTag string
	name         string
	cap          int
}

func (s *Slice[T]) TransformToXML(b *saver.Buffer) error {
	//if s.repeatingTag == "" {
	//	log.Print("Slice.TransformToXML: No repeating tag specified.")
	//	return nil
	//}
	//b.OpenTag(s.repeatingTag, s.attr)
	for _, v := range s.data {
		if err := v.TransformToXML(b); err != nil {
			return err
		}
		b.Write([]byte("\n"))
	}
	//b.WriteString("\n")
	//b.CloseTagWithIndent(s.repeatingTag)
	//b.WriteString("\n")
	return nil
}

func (s *Slice[T]) GetXMLTag() []byte {
	return []byte(s.repeatingTag)
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
	n := e
	if n == nil {
		log.Printf("n is nil")
		return nil
	}
	if n.Parent != nil {
		s.name = n.Parent.GetName()
		//log.Printf("Slice.Assign: Parent is %s", n.Parent.GetName())
	} else {
		//log.Printf("Slice.Assign: Assigning to slice without parent")
		s.name = "unknown"
	}
	s.repeatingTag = n.GetName()
	if !strings.Contains(reflect.TypeOf(zero[T]()).Name(), "types.Slice") {
		for n.Child != nil && n.Child.Child != nil {
			n = n.Child
		}
	}
	for n != nil {
		sd := sliceData[T]{
			tag: n.GetName(),
		}
		// Set sd.data to zero depending on the type of T. Either a pointer or a
		// value.
		switch tType := reflect.TypeOf(sd.data).Kind(); tType {
		case reflect.Ptr, reflect.Interface, reflect.Struct, reflect.Map, reflect.Slice:
			sd.data = reflect.New(reflect.TypeOf(*new(T)).Elem()).Interface().(T)
		case reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
			sd.data = zero[T]()
		}
		var attr attributes.Attributes
		if n.Child != nil {
			//log.Printf("Slice.Assign: child => %v", n.Child.GetName())
			attr = n.Child.Attr
			if err := unmarshal.Element(n.Child, &sd); err != nil {
				return err
			}
		} else {
			//log.Printf("Slice.Assign: no child=> %v", n.GetName())
			attr = n.Attr
			if err := unmarshal.Element(n, &sd); err != nil {
				return err
			}
		}
		s.data = append(s.data, sd)
		sd.SetAttributes(attr)
		n = n.Next
	}
	s.cap = len(s.data)
	//log.Printf("Slice.Assign: end => %s", s)
	return nil
}

func (s *Slice[T]) String() string {
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
