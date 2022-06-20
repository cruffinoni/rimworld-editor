package xml

import (
	"bytes"
	_xml "encoding/xml"
	"fmt"
	"strings"
)

type Element struct {
	StartElement _xml.StartElement
	EndElement   _xml.EndElement
	Attr         Attributes
	Data         *Data
	Index        int

	Next   *Element
	Prev   *Element
	Child  *Element
	Parent *Element
}

const DefaultSpacing = 2

func (e *Element) ToXML(spacing int) string {
	var sb strings.Builder
	n := e
	spaces := strings.Repeat(" ", spacing)
	for n != nil {
		sb.WriteString("\n" + spaces)
		sb.WriteString("<" + n.StartElement.Name.Local)
		if !n.Attr.Empty() {
			sb.WriteString(" " + n.Attr.Join(" "))
		}
		sb.WriteString(">")
		if n.Child != nil {
			sb.WriteString(n.Child.ToXML(spacing + DefaultSpacing))
		}
		if n.Data != nil {
			sb.WriteString(n.Data.GetString())
			sb.WriteString("</" + n.StartElement.Name.Local + ">")
		} else {
			sb.WriteString("\n" + spaces + "</" + n.StartElement.Name.Local + ">")
		}
		n = n.Next
	}
	return sb.String()
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

func (e *Element) Pretty(spacing int) string {
	var sb strings.Builder
	n := e
	for n != nil {
		sb.WriteString(strings.Repeat(" ", spacing) + "> " + n.StartElement.Name.Local)
		if !n.Attr.Empty() {
			sb.WriteString(" [" + n.Attr.Join(", ") + "]")
		}
		sb.WriteString("\n")
		if n.Child != nil {
			sb.WriteString(n.Child.Pretty(spacing + 2))
		}
		n = n.Next
	}
	return sb.String()
}

//func (e *Element) String() string {
//	var s string
//	s = fmt.Sprintf("%v[%v/%d] ", s, e.StartElement.Name.Local, e.Index)
//	if e.Child != nil {
//		s += "(" + e.Child.String() + ") "
//	}
//	s = fmt.Sprintf("%v'%v' ", s, e.Data)
//	l := len(s)
//	if l > 0 && s[l-1] == ' ' {
//		s = s[:l-1]
//	}
//	return s
//}

func (e *Element) GetName() string {
	return e.StartElement.Name.Local
}

func (e *Element) DisplayAllXMLPaths() string {
	var (
		sb strings.Builder
		n  = e
	)
	for n != nil {
		sb.WriteString(">" + n.StartElement.Name.Local)
		if n.Child != nil {
			sb.WriteString(n.Child.DisplayAllXMLPaths())
		}
		n = n.Next
	}
	return sb.String()
}

func (e *Element) xmlPath() *bytes.Buffer {
	b := &bytes.Buffer{}
	b.WriteString(e.StartElement.Name.Local)
	if e.Index > 0 {
		b.WriteString(fmt.Sprintf("[%d]", e.Index))
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
