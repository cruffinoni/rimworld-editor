package xml

import (
	_xml "encoding/xml"
	"io"

	"github.com/cruffinoni/rimworld-editor/helper"
	"github.com/cruffinoni/rimworld-editor/internal/xml/attributes"
)

// event is a type that represents a function to
// handle a event. T is a user-defined type.
// See Context for ctx
type event[T any] func(e T, ctx *Context)

type indexRemembering map[int]int

// Context give a context to an event
type Context struct {
	index indexRemembering
	attr  attributes.Attributes
	depth int
}

func transformAttrToMap(attr *[]_xml.Attr) attributes.Attributes {
	attrMap := make(attributes.Attributes)
	for _, a := range *attr {
		attrMap[a.Name.Local] = a.Value
	}
	return attrMap
}

const InvalidIdx = -1

func unmarshalEmbed(decoder *_xml.Decoder,
	onStartElement event[*_xml.StartElement],
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
			if helper.IsListTag(t.Name.Local) {
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

			previousIdx := ctx.depth + 1
			if !helper.IsListTag(t.Name.Local) && ctx.index[previousIdx] > 0 {
				delete(ctx.index, previousIdx)
			}
			ctx.depth--
		case _xml.CharData:
			if ctx.depth == 0 {
				continue
			}
			if onCharByte != nil {
				onCharByte(t, ctx)
			}
		}
	}
}
