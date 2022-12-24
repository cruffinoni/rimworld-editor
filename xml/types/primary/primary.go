package primary

import (
	"fmt"
	"reflect"

	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/xml/saver"
)

var (
	TypeNames []string
)

func init() {
	TypeNames = []string{
		reflect.TypeOf((*EmbeddedType[int])(nil)).Elem().Name(),
		reflect.TypeOf((*EmbeddedType[int64])(nil)).Elem().Name(),
		reflect.TypeOf((*EmbeddedType[uint])(nil)).Elem().Name(),
		reflect.TypeOf((*EmbeddedType[uint64])(nil)).Elem().Name(),
		reflect.TypeOf((*EmbeddedType[float64])(nil)).Elem().Name(),
		reflect.TypeOf((*EmbeddedType[float32])(nil)).Elem().Name(),
		reflect.TypeOf((*EmbeddedType[bool])(nil)).Elem().Name(),
		reflect.TypeOf((*EmbeddedType[string])(nil)).Elem().Name(),
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

type EmbeddedType[T comparable] struct {
	fmt.Stringer
	saver.Transformer

	data T
	tag  string
	// str is the string representation of the data.
	str  string
	attr attributes.Attributes
}

func (pt EmbeddedType[T]) TransformToXML(buffer *saver.Buffer) error {
	buffer.IncreaseDepth()
	buffer.WriteWithIndent([]byte(pt.str))
	buffer.DecreaseDepth()
	return nil
}

func lazyCheck(data any) {
	switch reflect.TypeOf(data).Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func, reflect.Struct:
		panic("primary.EmbeddedType: data must be a primitive type")
	}
}

func (pt *EmbeddedType[T]) Assign(e *xml.Element) error {
	// The type T must be a primitive type.
	lazyCheck(pt.data)
	if v, ok := e.Data.GetData().(T); ok {
		pt.tag = e.GetName()
		pt.data = v
		pt.str = fmt.Sprintf("%v", v)
		return nil
	} else {
		return fmt.Errorf("xml: cannot assign %T to %T", e.Data.GetData(), pt)
	}
}

func (pt *EmbeddedType[T]) GetPath() string {
	return ""
}

func (pt *EmbeddedType[T]) SetAttributes(attributes attributes.Attributes) {
	pt.attr = attributes
}

func (pt *EmbeddedType[T]) GetAttributes() attributes.Attributes {
	return pt.attr
}

func (pt EmbeddedType[T]) String() string {
	return pt.str
}
