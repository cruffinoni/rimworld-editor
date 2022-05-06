package xml

import (
	"bytes"
	_xml "encoding/xml"
	"log"
	"os"
)

type Discover struct {
	_xml.Unmarshaler
	Tag *Tag
}

func (d *Discover) Debug() {
	d.Tag.DisplayDebug()
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

func (d *Discover) FindTagsFromData(data string) []*Tag {
	var (
		validTag *Tag
		tags     = make([]*Tag, 0)
		n        = d.Tag
	)
	for n != nil {
		if validTag = n.FindTagFromData(data); validTag != nil {
			tags = append(tags, validTag)
		}
		n = n.Next
	}
	return tags
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
					//log.Printf("(ctx.depth > depth - %v) creating a new node: %v", e.Name.Local, ctx.depth)
					n := &Tag{
						StartElement: *e,
						EndElement:   _xml.EndElement{Name: e.Name},
						Parent:       lastNode,
						Attr:         e.Attr,
					}
					lastNode.Child = n
					lastNode = n
					//log.Printf("(ctx.depth > depth) created a new node / Parent: %p", lastNode)
					depth = ctx.depth
				}
			} else if ctx.depth < depth {
				if lastNode.Parent != nil {
					lastNode = lastNode.Parent
					nn := &Tag{
						Parent:       lastNode.Parent,
						Prev:         lastNode.Prev,
						StartElement: *e,
						EndElement:   _xml.EndElement{Name: e.Name},
						Attr:         e.Attr,
					}
					//log.Println("(ctx.depth < depth) retrieving prev node")
					lastNode.Next = nn
					lastNode = nn
					depth = ctx.depth
				} else {
					// This case should not happen because every child must have a parent
					log.Println("no prev but depth decreased")
					log.Printf("Last node: %v", lastNode.XMLPath())
					//fmt.Printf(d.Pretty())
					os.Exit(0)
				}
			} else {
				//log.Println("(ctx.depth == depth) creating a new node")
				nn := &Tag{
					Prev:         lastNode,
					StartElement: *e,
					EndElement:   _xml.EndElement{Name: e.Name},
					Attr:         e.Attr,
				}
				// First node is null because it is the root node
				if lastNode == nil {
					//log.Println("Root node created")
					lastNode = nn
					d.Tag = lastNode
				} else {
					// All children must have the same parent
					nn.Parent = lastNode.Parent
					lastNode.Next = nn
					lastNode = nn
					//log.Printf("Basic node created w/ parent: %p", lastNode.Parent)
				}
			}
		},
		func(e *_xml.EndElement, ctx *Context) {
		},
		func(b []byte, ctx *Context) {
			s := string(bytes.Trim(b, "\n\t"))
			if s != "" {
				lastNode.Data = s
				//log.Printf("Data: '%v'", s)
			}
		})
}
