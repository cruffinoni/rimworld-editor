package xml

import (
	_xml "encoding/xml"
	"io"
)

type event[T any] func(e T, ctx *Context)

type Context struct {
	index int
	depth int
}

func localXMLUnmarshal(decoder *_xml.Decoder,
	onStartElement event[*_xml.StartElement],
	onEndElement event[*_xml.EndElement],
	onCharByte event[[]byte]) error {
	ctx := &Context{}
	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		switch t := token.(type) {
		case _xml.StartElement:
			t.End()
			//log.Printf("StartElement: %s\n", t.Name.Local)
			////log.Printf("Attribute: %v\n", t.Attr)
			//startAcquired = true
			ctx.depth++
			onStartElement(&t, ctx)
		case _xml.EndElement:
			if ctx.depth == 0 {
				continue
			}
			//log.Printf("EndElement: %s\n", t.Name.Local)
			ctx.depth--
			//log.Printf("Depth: %v", ctx.depth)
			ctx.index++
			onEndElement(&t, ctx)
		case _xml.CharData:
			////log.Printf("CharData: '%s'\n", string(t))
			//log.Print("CharData called")
			if ctx.depth == 0 {
				continue
			}
			onCharByte(t, ctx)
		}
	}
}
