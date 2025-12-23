package domain

import (
	"bytes"
	_xml "encoding/xml"
	"fmt"
	"strings"

	"github.com/cruffinoni/rimworld-editor/internal/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/internal/xml/encoder"
)

type Element struct {
	StartElement _xml.StartElement
	EndElement   _xml.EndElement
	Attr         attributes.Attributes
	Data         *Data
	index        int

	Next   *Element
	Prev   *Element
	Child  *Element
	Parent *Element
}

const DefaultSpacing = 4

func (e *Element) ToXML(spacing int) string {
	var sb strings.Builder
	n := e
	spaces := strings.Repeat(" ", spacing)
	for n != nil {
		sb.WriteString("\n" + spaces)
		if n.IsEmpty() {
			sb.WriteString("<" + n.GetName())
			if !n.Attr.Empty() {
				sb.WriteString(" " + n.Attr.Join(" "))
			}
			sb.WriteString("/>")
			n = n.Next
			continue
		}
		sb.WriteString("<" + n.GetName())
		if !n.Attr.Empty() {
			sb.WriteString(" " + n.Attr.Join(" "))
		}
		sb.WriteString(">")
		if n.Child != nil {
			sb.WriteString(n.Child.ToXML(spacing + DefaultSpacing))
		}
		if n.Data != nil {
			sb.WriteString(n.Data.GetString())
			sb.WriteString("</" + n.GetName() + ">")
		} else {
			sb.WriteString("\n" + spaces + "</" + n.GetName() + ">")
		}
		n = n.Next
	}
	return sb.String()
}

func (e *Element) TransformToXML(buffer *encoder.Buffer) error {
	buffer.WriteString(e.ToXML((buffer.GetDepth() + 1) * DefaultSpacing))
	return nil
}

func (e *Element) GetXMLTag() []byte {
	return []byte(e.GetName())
}

func (e *Element) IsEmpty() bool {
	return e.Child == nil && e.Data == nil
}

func (e *Element) DisplayDebug() string {
	var sb strings.Builder
	n := e
	for n != nil {
		sb.WriteString(fmt.Sprintf("Node %p (%v) [parent: %p] ", n, n.StartElement.Name.Local, n.Parent))
		if n.Child != nil {
			sb.WriteString(fmt.Sprintf("[child: %p] ", n.Child))
		}
		n = n.Next
	}
	return sb.String()
}

// GetIndex returns the index of the element in the list of elements or 0 if the element is not in a list
func (e *Element) GetIndex() int {
	return e.index
}

func (e *Element) Pretty(spacing int) string {
	var sb strings.Builder
	n := e
	for n != nil {
		sb.WriteString(strings.Repeat(" ", spacing) + "> " + n.xmlPath().String())
		if !n.Attr.Empty() {
			sb.WriteString(" [" + n.Attr.Join(", ") + "]")
		}
		if n.Child != nil {
			//sb.WriteString("\n")
			sb.WriteString(n.Child.Pretty(spacing + 2))
		}
		n = n.Next
	}
	return sb.String()
}

func (e *Element) GetName() string {
	return e.StartElement.Name.Local
}

func (e *Element) DisplayAllXMLPaths() string {
	var (
		sb strings.Builder
		n  = e
	)
	for n != nil {
		sb.WriteString(">" + n.xmlPath().String())
		if n.Child != nil {
			sb.WriteString(n.Child.DisplayAllXMLPaths())
		}
		sb.WriteString("\n")
		n = n.Next
	}
	return sb.String()
}

func (e *Element) xmlPath() *bytes.Buffer {
	b := &bytes.Buffer{}
	b.WriteString(e.StartElement.Name.Local)
	if e.index > 0 {
		b.WriteString(fmt.Sprintf("[%d]", e.index))
	}
	return b
}

func (e *Element) XMLPath() string {
	var (
		n  = e
		sb []byte
	)
	for n != nil {
		buffer := []byte{'>'}
		buffer = append(buffer, n.xmlPath().Bytes()...)
		sb = append(buffer, sb...)
		n = n.Parent
	}
	return string(sb[1:])
}

func (e *Element) FindTagFromData(data string) []*Element {
	var (
		result = make([]*Element, 0)
		n      = e
	)
	for n != nil {
		if n.Data != nil {
			if n.Data.GetString() == data {
				return []*Element{n}
			}
		}
		if n.Child != nil {
			if r := n.Child.FindTagFromData(data); r != nil {
				result = append(result, r...)
			}
		}
		n = n.Next
	}
	return result
}

func (e *Element) SetAttributes(_ attributes.Attributes) {
	// We ignore the attribution because the structure has already set the attributes
}

func (e *Element) GetAttributes() attributes.Attributes {
	return e.Attr
}

// interface.Assigner interface implementation

func (e *Element) GetPath() string {
	// We avoid to use GetPath in generics type
	return ""
}

func (e *Element) Assign(elem *Element) error {
	if elem == nil {
		return nil
	}
	*e = *elem
	return nil
}
