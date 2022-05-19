package xml

import (
	"bytes"
	_xml "encoding/xml"
	"log"
	"os"
)

type Tree struct {
	_xml.Unmarshaler
	Root *Element
}

func (t *Tree) Debug() {
	t.Root.DisplayDebug()
}

func (t *Tree) Pretty(space ...int) string {
	spacing := 3
	if len(space) > 0 {
		spacing = space[0]
	}
	return t.Root.Pretty(spacing)
}

func (t *Tree) Raw() string {
	return t.Root.Raw()
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

func (t *Tree) UnmarshalXML(decoder *_xml.Decoder, _ _xml.StartElement) error {
	var (
		lastNode *Element
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
					idx := InvalidIdx
					if v, ok := ctx.index[ctx.depth]; ok {
						//log.Printf("(ctx.depth > depth) index: %v", v)
						idx = v
					}
					n := &Element{
						StartElement: *e,
						EndElement:   _xml.EndElement{Name: e.Name},
						Parent:       lastNode,
						Index:        idx,
						Attr:         ctx.attr,
					}
					lastNode.Child = n
					lastNode = n
					////log.Printf("(ctx.depth > depth) created a new node / Parent: %p", lastNode)
					depth = ctx.depth
				}
			} else if ctx.depth < depth {
				if lastNode.Parent != nil {
					for depth > ctx.depth {
						lastNode = lastNode.Parent
						depth--
					}
					idx := InvalidIdx
					if v, ok := ctx.index[ctx.depth]; ok {
						idx = v
					}
					//log.Printf("(ctx.depth < depth) creating a new node w/ idx: %v", idx)
					n := &Element{
						Parent:       lastNode.Parent,
						Left:         lastNode.Left,
						Index:        idx,
						StartElement: *e,
						EndElement:   _xml.EndElement{Name: e.Name},
						Attr:         ctx.attr,
					}
					//log.Println("(ctx.depth < depth) retrieving prev node")
					lastNode.Right = n
					lastNode = n
				} else {
					// This case should not happen because every child must have a parent
					log.Println("no prev but depth decreased")
					//log.Printf("Last node: %v", lastNode.XMLPath())
					//fmt.Printf(t.Pretty())
					os.Exit(1)
				}
			} else {
				//log.Printf("(ctx.depth == depth - %s) creating a new node", e.Name.Local)
				idx := InvalidIdx
				if v, ok := ctx.index[ctx.depth]; ok {
					//log.Printf("(ctx.depth == depth) index: %v", v)
					idx = v
				}
				n := &Element{
					Left:         lastNode,
					StartElement: *e,
					Index:        idx,
					EndElement:   _xml.EndElement{Name: e.Name},
					Attr:         ctx.attr,
				}
				// First node is null because it is the root node
				if lastNode == nil {
					//log.Println("Root node created")
					lastNode = n
					t.Root = lastNode
				} else {
					// All children must have the same parent
					n.Parent = lastNode.Parent
					lastNode.Right = n
					lastNode = n
					//log.Printf("Basic node created w/ parent: %p", lastNode.Parent)
				}
			}
		},
		nil,
		func(b []byte, ctx *Context) {
			s := string(bytes.Trim(b, "\n\t"))
			if s != "" {
				lastNode.Data = s
				//log.Printf("Data: '%v'", s)
			}
		})
}
