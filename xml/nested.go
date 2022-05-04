package xml

import (
	_xml "encoding/xml"
	"log"
	"strings"
)

type Nested struct {
	_xml.Unmarshaler
	Data  content
	Token string
}

type Index struct {
	Class string
	Def   string
}

type content map[*Index]*Tag

func (c *content) getIndexFromClass(class string) *Index {
	if class == "" {
		return nil
	}
	for k := range *c {
		if k.Class == class {
			return k
		}
	}
	return nil
}

func (c *content) getIndexFromDef(def string) *Index {
	if def == "" {
		return nil
	}
	for k := range *c {
		if k.Def == def {
			return k
		}
	}
	return nil
}

func findClassFromAttr(attr []_xml.Attr) string {
	for _, a := range attr {
		if a.Name.Local == "Class" {
			return a.Value
		}
	}
	return ""
}

func (l *Nested) getIndexFromClass(class string) *Index {
	return l.Data.getIndexFromClass(class)
}

func (l *Nested) getIndexFromDef(def string) *Index {
	return l.Data.getIndexFromDef(def)
}

func (l *Nested) indexExists(index *Index) bool {
	return l.getIndexFromClass(index.Class) != nil && l.getIndexFromDef(index.Def) != nil
}

func (l *Nested) UnmarshalXML(decoder *_xml.Decoder, _ _xml.StartElement) error {
	l.Data = make(content)
	if l.Token == "" {
		l.Token = "li"
	}
	var (
		lastNode *Tag
		idx      *Index
		isDef    = false
		depth    = 1
	)
	return localXMLUnmarshal(decoder,
		func(e *_xml.StartElement, ctx *Context) {
			c := findClassFromAttr(e.Attr)
			if c != "" && e.Name.Local == l.Token {
				idx = &Index{Class: c}
				l.Data[idx] = &Tag{StartElement: *e}
				lastNode = l.Data[idx]
				return
			} else if idx != nil {
				isDef = e.Name.Local == "def"
			}

			if ctx.depth > depth {
				if lastNode.Child != nil {
					log.Println("(ctx.depth > depth) retrieving next node")
					lastNode = lastNode.Child
				} else {
					log.Println("(ctx.depth > depth) creating a new node")
					n := &Tag{
						StartElement: *e,
						Parent:       lastNode,
					}
					lastNode.Child = n
					lastNode = n
					depth = ctx.depth
				}
			} else if ctx.depth < depth {
				if lastNode.Parent != nil {
					log.Println("(ctx.depth < depth) retrieving prev node")
					lastNode = lastNode.Parent
					depth = ctx.depth
				} else {
					log.Printf("no prev but depth decreased")
				}
			} else {
				log.Println("(ctx.depth == depth) creating a new node")
				lastNode.Next = &Tag{
					Prev:         lastNode,
					StartElement: *e,
				}
				lastNode = lastNode.Next
			}
		},
		func(e *_xml.EndElement, ctx *Context) {
		},
		func(b []byte, ctx *Context) {
			if isDef {
				idx.Def = string(b)
			}
			lastNode.Data = strings.Trim(string(b), "\n\t")
		},
	)
}
