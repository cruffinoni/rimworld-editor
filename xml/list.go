package xml

import _xml "encoding/xml"

type List struct {
	_xml.Unmarshaler
	Data []Tag
}

func (l *List) UnmarshalXML(decoder *_xml.Decoder, _ _xml.StartElement) error {
	l.Data = make([]Tag, 0)
	return localXMLUnmarshal(decoder,
		func(e *_xml.StartElement, ctx *Context) {
			l.Data = append(l.Data, Tag{StartElement: *e})
		},
		func(e *_xml.EndElement, ctx *Context) {
			l.Data[ctx.index].EndElement = *e
		},
		func(c []byte, ctx *Context) {
			l.Data[ctx.index].Data = string(c)
		})
}
