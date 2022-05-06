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
		sb.WriteString(fmt.Sprintf("%s<%v>%v</%v>\n", strings.Repeat(" ", spacing[0]), n.StartElement.Name.Local, d, n.EndElement.Name.Local))
		if n.Child != nil {
			sb.WriteString(n.Child.Raw(spacing[0] + 2))
		}
		n = n.Next
	}
	return sb.String()
}

func (t *Tag) DisplayDebug() string {
	var sb strings.Builder
	n := t
	for n != nil {
		sb.WriteString(fmt.Sprintf("Node %p (%v) [parent: %p] ", n, n.StartElement.Name.Local, n.Parent))
		if n.Child != nil {
			sb.WriteString(fmt.Sprintf("[child: %p] ", n.Child))
		}
		n = n.Next
	}
	return sb.String()
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

func (t *Tag) GetName() string {
	return t.StartElement.Name.Local
}

func (t *Tag) DisplayAllXMLPaths() string {
	var (
		sb strings.Builder
		n  = t
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

func (t *Tag) XMLPath() string {
	var (
		buffer []byte
		n      = t
		sb     []byte
	)
	for n != nil {
		buffer = make([]byte, len(n.StartElement.Name.Local)+2)
		buffer = append(buffer, n.StartElement.Name.Local+">"...)

		sb = append(buffer, sb...)

		n = n.Parent
	}
	var (
		s = string(sb)
		l = len(s)
	)
	if s[l-1] == '>' {
		s = s[:l-1]
	}
	return s
}

func (t *Tag) FindTagFromData(data string) []*Tag {
	var (
		result = make([]*Tag, 0)
		d      string
		ok     bool
		n      = t
	)
	//log.Printf("going with %v", n)
	for n != nil {
		//log.Printf("[%v] w/ '%v'", n.StartElement.Name.Local, n.Data)
		if d, ok = n.Data.(string); ok {
			if d == data {
				result = append(result, n)
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
