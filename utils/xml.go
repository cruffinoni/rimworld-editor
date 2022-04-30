package utils

import (
	"encoding/xml"
	"io"
)

type GenericSimpleFormat struct {
	xml.Unmarshaler
	StartElement xml.StartElement
	EndElement   xml.EndElement
	Data         any
}

//type GenericSimpleFormat struct {
//	StartElement xml.StartElement
//	EndElement   xml.EndElement
//	Data         any
//}

func (s *GenericSimpleFormat) UnmarshalXML(decoder *xml.Decoder, _ xml.StartElement) error {
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
			s.StartElement = t
		case xml.EndElement:
			s.EndElement = t
		case xml.CharData:
			s.Data = t.Copy()
		}
	}
}
