package generator

import (
	"bytes"
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
	"log"
	"os"
	"reflect"
)

type flag uint

const (
	flagNone      flag = 0 << iota
	hasAttributes flag = 1 << iota
)

type metadata struct {
	dir   os.DirEntry
	flags flag
}

type StructInfo struct {
	name    string
	members []*member
	buf     bytes.Buffer
}

type member struct {
	name string
	t    any
	attr attributes.Attributes
}

type customType struct {
	name   string
	pkg    string
	types1 any
	types2 any
}

func GenerateGoFiles(root *xml.Element) *StructInfo {
	s := StructInfo{
		name:    "save",
		members: make([]*member, 0),
	}
	if err := handleElement(root, &s); err != nil {
		panic(err)
	}
	return &s
}

func (s *StructInfo) WriteGoFile(path string) error {
	path = "./" + path
	if _, err := os.Stat(path); err == nil {
		if err = os.RemoveAll(path); err != nil {
			return err
		}
	}
	if err := os.Mkdir(path, os.ModePerm); err != nil {
		return err
	}
	return s.generateStructTo(path)
}

func getTypeFromArray(e *xml.Element) reflect.Kind {
	k := e.Child
	if k != nil && k.Child != nil {
		return reflect.Struct
	}
	kt := reflect.Invalid
	for k != nil {
		if k.Data != nil {
			kdk := k.Data.Kind()
			if kt != reflect.Invalid && kdk != kt &&
				// Float64 and Int64 are interchangeable
				!(kdk == reflect.Float64 && kt == reflect.Int64) &&
				!(kdk == reflect.Int64 && kt == reflect.Float64) {
				log.Panicf("primary.EmbeddedType: found type %v, expected %v", kdk, kt)
			}
			// Float64 and Int64 are interchangeable, but we prefer to keep Float64
			if !(kdk == reflect.Int64 && kt == reflect.Float64) {
				kt = kdk
			}
		}
		k = k.Next
	}
	return kt
}

func generateStructure(e *xml.Element) *StructInfo {
	log.Printf("Generating structure for %v", e.GetName())
	s := &StructInfo{
		name:    e.GetName(),
		members: make([]*member, 0),
	}
	if err := handleElement(e.Child, s); err != nil {
		panic(err)
	}
	s.removeDuplicates()
	return s
}

func createCustomTypeForMap(e *xml.Element) any {
	if e.Child == nil {
		log.Panic("primary.EmbeddedType: missing child")
	}
	var (
		c     = e.Child
		k any = getTypeFromArray(c)
		v any
	)
	// We need to define this struct with of this all members
	if k == reflect.Struct {
		// This is a slice of structs, we need to define this struct will all members
		k = generateStructure(c)
	} else if k == reflect.Invalid {
		k = e
	}
	c = c.Next
	v = getTypeFromArray(c)
	// Repeat the operation for the value
	if v == reflect.Struct {
		v = generateStructure(c)
	} else if v == reflect.Invalid {
		v = e
	}
	// By default, maps are strings to strings
	if k == reflect.Invalid || v == reflect.Invalid {
		return &customType{
			name:   "Map",
			pkg:    "types",
			types1: reflect.String,
			types2: reflect.String,
		}
	}
	return &customType{
		name:   "Map",
		pkg:    "types",
		types1: k,
		types2: v,
	}
}

func (s *StructInfo) removeDuplicates() {
	for i := 0; i < len(s.members); i++ {
		for j := i + 1; j < len(s.members); j++ {
			if s.members[i].name == s.members[j].name {
				//if s.members[i].t != s.members[j].t {
				//	log.Panicf("primary.EmbeddedType: duplicate member %v with different types ; expected %T, got %T", s.members[i].name, s.members[i].t, s.members[j].t)
				//}
				s.members = append(s.members[:i], s.members[i+1:]...)
				i--
				break
			}
		}
	}
}

func createCustomSlice(e *xml.Element) any {
	c := e
	var t any
	t = getTypeFromArray(c)
	if t == reflect.Struct {
		t = generateStructure(c)
	} else if t == reflect.Invalid {
		t = e
	}
	return &customType{
		name:   "Slice",
		pkg:    "types",
		types1: t,
	}
}

func handleElement(e *xml.Element, st *StructInfo) error {
	n := e
	for n != nil {
		var t any
		if n.Child != nil {
			// Skip the "li" tag since it's a slice and should not be a member of the struct
			if n.GetName() == "li" {
				if err := handleElement(n.Child, st); err != nil {
					return err
				}
			} else {
				switch n.Child.GetName() {
				case "li":
					t = createCustomSlice(n)
				case "keys":
					t = createCustomTypeForMap(n)
				default:
					t = &StructInfo{
						name:    n.GetName(),
						members: make([]*member, 0),
					}
					if err := handleElement(n.Child, t.(*StructInfo)); err != nil {
						return err
					}
				}
				st.members = append(st.members, &member{
					name: n.GetName(),
					t:    t,
					attr: n.Attr,
				})
			}
		} else {
			if n.Data != nil {
				t = n.Data.Kind()
			} else {
				t = e
			}
			st.members = append(st.members, &member{
				name: n.GetName(),
				t:    t,
				attr: n.Attr,
			})
		}
		n = n.Next
	}
	return nil
}
