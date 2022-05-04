package xml

import (
	_xml "encoding/xml"
)

type Discover struct {
	_xml.Unmarshaler
	Tag *Tag
}

func (d *Discover) Pretty(space ...int) string {
	spacing := 3
	if len(space) > 0 {
		spacing = space[0]
	}
	return d.Tag.Pretty(spacing)
}

func (d *Discover) Raw() string {
	return d.Tag.Raw()
}

func (d *Discover) UnmarshalXML(decoder *_xml.Decoder, _ _xml.StartElement) error {
	var (
		lastNode *Tag
		depth    = 1
	)

	return localXMLUnmarshal(decoder,
		func(e *_xml.StartElement, ctx *Context) {
			if ctx.depth > depth {
				if lastNode.Child != nil {
					//log.Printf("(ctx.depth > depth) retrieving next node: %v", ctx.depth)
					lastNode = lastNode.Child
				} else {
					//log.Printf("(ctx.depth > depth) creating a new node: %v", ctx.depth)
					n := &Tag{
						StartElement: *e,
						Parent:       lastNode,
						Attr:         e.Attr,
					}
					lastNode.Child = n
					lastNode = n
					depth = ctx.depth
				}
			} else if ctx.depth < depth {
				if lastNode.Parent != nil {
					nn := &Tag{
						Prev:         lastNode.Parent,
						StartElement: *e,
						Attr:         e.Attr,
					}
					//log.Println("(ctx.depth < depth) retrieving prev node")
					lastNode.Parent.Next = nn
					lastNode = nn
					depth = ctx.depth
				} else {
					//log.Printf("no prev but depth decreased")
				}
			} else {
				//log.Println("(ctx.depth == depth) creating a new node")
				nn := &Tag{
					Prev:         lastNode,
					StartElement: *e,
					Attr:         e.Attr,
				}
				if lastNode == nil {
					lastNode = nn
					d.Tag = lastNode
				} else {
					lastNode.Next = nn
					lastNode = lastNode.Next
				}
			}
		},
		func(e *_xml.EndElement, ctx *Context) {
		},
		func(b []byte, ctx *Context) {
		})
}
