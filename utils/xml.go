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
		switch t := token.(type) {
		case xml.StartElement:
			log.Printf("StartElement: %s\n", t.Name.Local)
			log.Printf("Attribute: %v\n", t.Attr)
			startAcquired = true
			os(&t, i)
		case xml.EndElement:
			if !startAcquired {
				continue
			}
			log.Printf("EndElement: %s\n", t.Name.Local)
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

type XMLListReader struct {
	xml.Unmarshaler
	Data map[string]map[string][]*XMLTag
}

func findClassFromAttr(attr []xml.Attr) string {
	for _, a := range attr {
		if a.Name.Local == "Class" {
			return a.Value
		}
	}
	return ""
}

func (l *XMLListReader) UnmarshalXML(decoder *xml.Decoder, _ xml.StartElement) error {
	l.Data = make(map[string][]*XMLTag)
	lastClassName := ""
	lastIdx := 0
	toIgnore := false
	return localXMLUnmarshaler(decoder,
		func(e *xml.StartElement, index int) {
			if tmpClassName := findClassFromAttr(e.Attr); tmpClassName == "" {
				l.Data[lastClassName] = append(l.Data[lastClassName], &XMLTag{StartElement: *e})
				lastIdx = len(l.Data[lastClassName]) - 1
				toIgnore = false
			} else {
				lastClassName = tmpClassName
				toIgnore = true
			}
		},
		func(e *xml.EndElement, index int) {
			if toIgnore {
				return
			}
			l.Data[lastClassName][lastIdx].EndElement = *e
		},
		func(c []byte, index int) {
			if toIgnore {
				return
			}
			l.Data[lastClassName][lastIdx].Data = string(c)
		})
}
