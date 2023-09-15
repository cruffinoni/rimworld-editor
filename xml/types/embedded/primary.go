package embedded

import (
	"fmt"
	"reflect"

	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/xml/saver"
	"github.com/cruffinoni/rimworld-editor/xml/saver/xmlFile"
)

var (
	TypeNames []string
)

func init() {
	TypeNames = []string{
		reflect.TypeOf((*Type[int])(nil)).Elem().Name(),
		reflect.TypeOf((*Type[int64])(nil)).Elem().Name(),
		reflect.TypeOf((*Type[uint])(nil)).Elem().Name(),
		reflect.TypeOf((*Type[uint64])(nil)).Elem().Name(),
		reflect.TypeOf((*Type[float64])(nil)).Elem().Name(),
		reflect.TypeOf((*Type[float32])(nil)).Elem().Name(),
		reflect.TypeOf((*Type[bool])(nil)).Elem().Name(),
		reflect.TypeOf((*Type[string])(nil)).Elem().Name(),
	}
}

func IsEmbeddedPrimaryType(name string) bool {
	for _, n := range TypeNames {
		if name == n {
			return true
		}
	}
	return false
}

type Type[T comparable] struct {
	fmt.Stringer
	saver.Transformer

	data T
	tag  string
	// str is the string representation of the data.
	str  string
	attr attributes.Attributes
}

func (pt *Type[T]) TransformToXML(buffer *saver.Buffer) error {
	l := len(pt.str)
	if l == 0 {
		return xmlFile.ErrEmptyValue
	}
	buffer.WriteString(pt.str)
	return nil
}

func lazyCheck(data any) {
	switch reflect.TypeOf(data).Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func, reflect.Struct:
		panic("primary.EmbeddedType: data must be a primitive type")
	}
}

func (pt *Type[T]) Assign(e *xml.Element) error {
	// The type T must be a primitive type.
	lazyCheck(pt.data)
	if v, ok := e.Data.GetData().(T); ok {
		pt.tag = e.GetName()
		pt.data = v
		pt.str = fmt.Sprintf("%v", v)
		return nil
	} else {
		return fmt.Errorf("EmbeddedType.Assign: cannot assign %T to %T", e.Data.GetData(), pt)
	}
}

func (pt *Type[T]) GetXMLTag() []byte {
	return []byte(pt.tag)
}

func (pt *Type[T]) ValidateField(_ string) {
}

func (pt *Type[T]) IsValidField(_ string) bool {
	return len(pt.str) > 0
}

func (pt *Type[T]) CountValidatedField() int {
	if pt.IsValidField("") {
		return 1
	}
	return 0
}

func (pt *Type[T]) GetPath() string {
	return ""
}

func (pt *Type[T]) SetAttributes(attributes attributes.Attributes) {
	pt.attr = attributes
}

func (pt *Type[T]) GetAttributes() attributes.Attributes {
	return pt.attr
}

func (pt *Type[T]) String() string {
	return pt.str
}
