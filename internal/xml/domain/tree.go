package domain

import (
	"bytes"
	_xml "encoding/xml"
	"log"

	"github.com/cruffinoni/rimworld-editor/internal/xml/decoder"
)

type Tree struct {
	_xml.Unmarshaler
	Root *Element
}

func (t *Tree) Debug() string {
	return t.Root.DisplayDebug()
}

func (t *Tree) Pretty(space ...int) string {
	spacing := 3
	if len(space) > 0 {
		spacing = space[0]
	}
	return t.Root.Pretty(spacing)
}

func (t *Tree) ToXML() string {
	return t.Root.ToXML(0)
}

func (t *Tree) XMLPath() string {
	return t.Root.XMLPath()
}

func (t *Tree) FindElementFromData(data string) []*Element {
	tags := make([]*Element, 0)
	// From root, finds all tags that match the data
	if t.Root != nil {
		if validTag := t.Root.FindTagFromData(data); validTag != nil {
			tags = append(tags, validTag...)
		}
	}
	return tags
}

func (t *Tree) UnmarshalXML(xmlDecoder *_xml.Decoder, s _xml.StartElement) error {
	var (
		lastNode = &Element{StartElement: s}
		depth    = 0
	)
	t.Root = lastNode

	return decoder.UnmarshalEmbed(xmlDecoder,
		func(e *_xml.StartElement, ctx *decoder.Context) {
			if ctx.Depth > depth {
				if lastNode.Child != nil {
					// The last node has already a child, so we retrieve it
					lastNode = lastNode.Child
				} else {
					// No child and this is a new element, so we create a new node

					idx := decoder.InvalidIdx
					// If the index has a registered depth, we use it.
					// It adds additional information to the node
					if v, ok := ctx.Index[ctx.Depth]; ok {
						idx = v
					}
					n := &Element{
						StartElement: *e,
						EndElement:   _xml.EndElement{Name: e.Name},
						Parent:       lastNode,
						index:        idx,
						Attr:         ctx.Attr,
					}
					lastNode.Child = n
					lastNode = n
					depth = ctx.Depth
				}
			} else if ctx.Depth < depth {
				//printer.Debugf("Depth decreased: %v", depth)
				// The depth is smaller than the last one, so we go back to the parent
				if lastNode.Parent != nil {
					// We don't detect the end of a XML element, so we have to
					// go back where the new element is.
					// Because if we are at the end of multiple element, the depth
					// will decrease multiple times.
					for depth > ctx.Depth {
						lastNode = lastNode.Parent
						depth--
					}
					idx := decoder.InvalidIdx

					// We do the same thing as previously explained
					if v, ok := ctx.Index[ctx.Depth]; ok {
						idx = v
					}
					// TODO: Code factorization ?
					n := &Element{
						Parent:       lastNode.Parent,
						Prev:         lastNode.Prev,
						index:        idx,
						StartElement: *e,
						EndElement:   _xml.EndElement{Name: e.Name},
						Attr:         ctx.Attr,
					}
					lastNode.Next = n
					lastNode = n
				} else {
					// This case should not happen because every child must have a parent
					log.Fatal("no prev but depth decreased")
				}
			} else {
				// The depth is the same as the last one, so we create a new node
				// which is basically a sibling of the current node
				//printer.Debugf("Depth is the same: %v", depth)
				idx := decoder.InvalidIdx
				if v, ok := ctx.Index[ctx.Depth]; ok {
					idx = v
				}
				n := &Element{
					Prev:         lastNode,
					StartElement: *e,
					index:        idx,
					EndElement:   _xml.EndElement{Name: e.Name},
					Attr:         ctx.Attr,
				}
				// All children must have the same parent because
				// they are siblings
				n.Parent = lastNode.Parent
				lastNode.Next = n
				lastNode = n
			}
		},
		func(b []byte, ctx *decoder.Context) {
			// This is a small trick to ignore data with only spaces
			// It occurs when an XML element close
			s := string(bytes.TrimSpace(b))
			if s != "" {
				lastNode.Data = CreateDataType(s)
				//printer.Debugf("Data: '%v' from %s", s, lastNode.XMLPath())
			}
		})
}
