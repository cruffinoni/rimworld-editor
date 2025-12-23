package domain

import (
	"github.com/cruffinoni/rimworld-editor/internal/xml/encoder"
)

type Elements []*Element

func (e Elements) FindElementFromClass(class string) *Element {
	for _, el := range e {
		if el.Attr.Get("Class") == class {
			return el
		}
	}
	return nil
}

func (e Elements) TransformToXML(buffer *encoder.Buffer) error {
	for _, el := range e {
		buffer.Write([]byte("\n"))
		if err := el.TransformToXML(buffer); err != nil {
			return err
		}
	}
	return nil
}
