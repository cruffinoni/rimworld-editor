package utils

import (
	"encoding/xml"
	"io"
	"log"
)

type XMLTag struct {
	StartElement xml.StartElement
	EndElement   xml.EndElement
	Data         any
}

type GenericSimpleFormat struct {
	xml.Unmarshaler
	Data []XMLTag
}

type onStartElement func(e *xml.StartElement, index int)
type onEndElement func(e *xml.EndElement, index int)
type onCharData func(c []byte, index int)

func localXMLUnmarshaler(decoder *xml.Decoder, os onStartElement, oe onEndElement, oc onCharData) error {
	i := 0
	startAcquired := false
	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		log.Printf("Token type: %T\n", token)
		switch t := token.(type) {
		case xml.StartElement:
			log.Printf("Attribute: %v\n", t.Attr)
			startAcquired = true
			os(&t, i)
		case xml.EndElement:
			if !startAcquired {
				continue
			}
			oe(&t, i)
			startAcquired = false
			i++
		case xml.CharData:
			log.Printf("CharData: '%s'\n", string(t))
			if !startAcquired {
				continue
			}
			oc(t, i)
		}
	}
}

func (f *GenericSimpleFormat) UnmarshalXML(decoder *xml.Decoder, _ xml.StartElement) error {
	f.Data = make([]XMLTag, 0)
	return localXMLUnmarshaler(decoder,
		func(e *xml.StartElement, index int) {
			f.Data = append(f.Data, XMLTag{StartElement: *e})
		},
		func(e *xml.EndElement, index int) {
			f.Data[index].EndElement = *e
		},
		func(c []byte, index int) {
			f.Data[index].Data = string(c)
		})
}

type GenericFormatByMap struct {
	xml.Unmarshaler
	Data map[string]*XMLTag
}

func (m *GenericFormatByMap) UnmarshalXML(decoder *xml.Decoder, _ xml.StartElement) error {
	m.Data = make(map[string]*XMLTag)
	var lastIdx string
	return localXMLUnmarshaler(decoder,
		func(e *xml.StartElement, index int) {
			log.Printf("StartElement: %s\n", e.Name.Local)
			m.Data[e.Name.Local] = &XMLTag{StartElement: *e}
			lastIdx = e.Name.Local
		},
		func(e *xml.EndElement, index int) {
			m.Data[e.Name.Local].EndElement = *e
		},
		func(c []byte, index int) {
			m.Data[lastIdx].Data = string(c)
		})
}
