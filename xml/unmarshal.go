package xml

import (
	_xml "encoding/xml"
	"io"
)

type event[T any] func(e T, ctx *Context)

type indexRemembering map[int]int

type Context struct {
	index indexRemembering
	attr  Attributes
	depth int
}

func transformAttrToMap(attr *[]_xml.Attr) Attributes {
	attrMap := make(Attributes)
	for _, a := range *attr {
		attrMap[a.Name.Local] = a.Value
	}
	return attrMap
}

const InvalidIdx = -1

func unmarshalEmbed(decoder *_xml.Decoder,
	onStartElement event[*_xml.StartElement],
	onEndElement event[*_xml.EndElement],
	onCharByte event[[]byte]) error {
	ctx := &Context{
		index: make(indexRemembering),
	}
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
			ctx.depth++
			ctx.attr = transformAttrToMap(&t.Attr)
			if t.Name.Local == "li" {
				ctx.index[ctx.depth]++
			}
			if onStartElement != nil {
				onStartElement(&t, ctx)
			}
		case _xml.EndElement:
			if ctx.depth == 0 {
				continue
			}
			ctx.attr = nil
			//log.Printf("EndElement: %s at %v => %v\n", t.Name.Local, ctx.depth, ctx.index[ctx.depth])
			ctx.depth--
			//log.Printf("Depth: %v - %v - %v", ctx.depth, ctx.index[ctx.depth+1], ctx.index[ctx.depth])
			if onEndElement != nil {
				onEndElement(&t, ctx)
			}

			previousIdx := ctx.depth + 2
			//log.Printf("Same idx: %v", ctx.index[previousIdx])
			if t.Name.Local != "li" && ctx.index[previousIdx] > 0 {
				//log.Printf("debug: %v - %v - %v", ctx.depth, ctx.index[previousIdx], ctx.index[ctx.depth])
				delete(ctx.index, previousIdx)
			}
			//log.Printf("ctx.index[ctx.depth+1]: %v", ctx.index[ctx.depth+1])
		case _xml.CharData:
			//log.Print("CharData called")
			if ctx.depth == 0 {
				continue
			}
			if onCharByte != nil {
				onCharByte(t, ctx)
			}
		}
	}
}
