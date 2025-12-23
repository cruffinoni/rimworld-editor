package decoder

import (
	_xml "encoding/xml"
	"io"

	"github.com/cruffinoni/rimworld-editor/internal/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/internal/xml/support"
)

// event is a type that represents a function to
// handle a event. T is a user-defined type.
// See Context for ctx
type event[T any] func(e T, ctx *Context)

type indexRemembering map[int]int

// Context give a context to an event
type Context struct {
	Index indexRemembering
	Attr  attributes.Attributes
	Depth int
}

func transformAttrToMap(attr *[]_xml.Attr) attributes.Attributes {
	attrMap := make(attributes.Attributes)
	for _, a := range *attr {
		attrMap[a.Name.Local] = a.Value
	}
	return attrMap
}

const InvalidIdx = -1

func UnmarshalEmbed(decoder *_xml.Decoder,
	onStartElement event[*_xml.StartElement],
	onCharByte event[[]byte]) error {
	ctx := &Context{
		Index: make(indexRemembering),
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
			ctx.Depth++
			ctx.Attr = transformAttrToMap(&t.Attr)
			if support.IsListTag(t.Name.Local) {
				ctx.Index[ctx.Depth]++
			}
			if onStartElement != nil {
				onStartElement(&t, ctx)
			}
		case _xml.EndElement:
			if ctx.Depth == 0 {
				continue
			}
			ctx.Attr = nil

			previousIdx := ctx.Depth + 1
			if !support.IsListTag(t.Name.Local) && ctx.Index[previousIdx] > 0 {
				delete(ctx.Index, previousIdx)
			}
			ctx.Depth--
		case _xml.CharData:
			if ctx.Depth == 0 {
				continue
			}
			if onCharByte != nil {
				onCharByte(t, ctx)
			}
		}
	}
}
