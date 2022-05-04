package xml

import (
	_xml "encoding/xml"
	"fmt"
	"strings"
)

type Tag struct {
	StartElement _xml.StartElement
	EndElement   _xml.EndElement
	Attr         []_xml.Attr
	Data         any

	Next   *Tag
	Prev   *Tag
	Child  *Tag
	Parent *Tag
}

const DefaultSpacing = 3

func (t *Tag) Raw(spacing ...int) string {
	if len(spacing) == 0 {
		spacing = []int{DefaultSpacing}
	}
	var sb strings.Builder
	n := t
	for n != nil {
		d := n.Data
		if d == nil {
			d = "null"
		}
		sb.WriteString(fmt.Sprintf("%s<%v>%v</%v>\n", strings.Repeat(" ", spacing[0]), n.StartElement.Name.Local, d, t.EndElement.Name.Local))
		if n.Child != nil {
			sb.WriteString(n.Child.Raw(spacing[0] + 2))
		}
		n = n.Next
	}
	return sb.String()
}

func (t *Tag) DisplayDebug() string {
	n := t
	k := ""
	for n != nil {
		k += fmt.Sprintf("Node %p (%v) [parent: %p] ", n, n.StartElement.Name.Local, n.Parent)
		if n.Child != nil {
			k += n.Child.DisplayDebug()
		}
		n = n.Next
	}
	return k
}

func (t *Tag) Pretty(spacing int) string {
	var sb strings.Builder
	n := t
	for n != nil {
		sb.WriteString(fmt.Sprintf("%s> %s\n", strings.Repeat(" ", spacing), n.StartElement.Name.Local))
		if n.Child != nil {
			sb.WriteString(n.Child.Pretty(spacing + 2))
		}
		n = n.Next
	}
	return sb.String()
}

func (t *Tag) String() string {
	var s string
	n := t
	for n != nil {
		s = fmt.Sprintf("%v[%v] ", s, n.StartElement.Name.Local)
		if n.Child != nil {
			s += "(" + n.Child.String() + ") "
		}
		s = fmt.Sprintf("%v'%v' ", s, n.Data)
		n = n.Next
	}
	l := len(s)
	if l > 0 && s[l-1] == ' ' {
		s = s[:l-1]
	}
	return s
}
