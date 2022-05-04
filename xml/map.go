package xml

import (
	_xml "encoding/xml"
	"log"
)

type Map struct {
	_xml.Unmarshaler
	Data map[string]*Tag
}

func (m *Map) UnmarshalXML(decoder *_xml.Decoder, _ _xml.StartElement) error {
	m.Data = make(map[string]*Tag)
	var lastIdx string
	return localXMLUnmarshal(decoder,
		func(e *_xml.StartElement, _ *Context) {
			log.Printf("StartElement: %s\n", e.Name.Local)
			m.Data[e.Name.Local] = &Tag{StartElement: *e}
			lastIdx = e.Name.Local
		},
		func(e *_xml.EndElement, _ *Context) {
			m.Data[e.Name.Local].EndElement = *e
		},
		func(c []byte, _ *Context) {
			m.Data[lastIdx].Data = string(c)
		})
}
