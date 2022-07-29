package generator

import (
	"bytes"
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
	"log"
	"math/rand"
	"os"
	"reflect"
)

type flag uint

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
	if err := handleElement(root, &s, flagNone); err != nil {
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

func isListTag(tag string) bool {
	return tag == "li" || tag == "list"
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
				log.Panicf("primary.EmbeddedType: found type %v, expected %v on path %v ('%v')", kdk, kt, k.XMLPath(), k.Data.GetData())
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

const (
	stringFixedSize = 5
	letterBytes     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func generateRandomString() string {
	b := make([]byte, stringFixedSize)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func createStructure(e *xml.Element, flag uint) *StructInfo {
	if flag&forceChild == forceChild {
		flag &^= forceChild
		if e.Child != nil && e.Child.Child != nil {
			return createStructure(e.Child, flag)
		} else {
			panic("generate.createStructure|forceChild: missing child")
		}
	}
	if e.Child == nil {
		panic("generate.createStructure: missing child")
	}
	name := e.GetName()
	if isListTag(name) {
		name += generateRandomString()
		log.Printf("debug: name '%v' & xmlPath: %v", name, e.XMLPath())
	}
	s := &StructInfo{
		name:    name,
		members: make([]*member, 0),
	}
	if err := handleElement(e.Child, s, flag); err != nil {
		panic(err)
	}
	s.removeDuplicates()
	return s
}

func getTypeFromAny(e *xml.Element, flag uint) any {
	t := any(getTypeFromArray(e))
	// We need to define this struct with of this all members
	if t == reflect.Struct {
		c := e.Child
		if isListTag(c.Child.GetName()) {
			t = createCustomSlice(c, flag|forceChild)
		} else {
			t = createStructure(e, flag)
		}
	} else if t == reflect.Invalid {
		t = e
	}
	return t
}

func createCustomTypeForMap(e *xml.Element, flag uint) any {
	if e.Child == nil {
		log.Panic("generate.createCustomTypeForMap: missing child")
	}

	var (
		c = e.Child
		k = getTypeFromAny(c, flag)
		v any
	)
	c = c.Next
	v = getTypeFromAny(c, flag)
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

func createCustomSlice(e *xml.Element, flag uint) any {
	c := e
	var t any
	t = getTypeFromArray(c)
	if t == reflect.Struct {
		t = createStructure(c, flag|skipChild)
		//if c.Child.GetName() == "li" {
		//	t = createCustomSlice(c.Child, flag)
		//} else {
		//	t = createStructure(c, flag)
		//}
	} else if t == reflect.Invalid {
		t = e
	}
	return &customType{
		name:   "Slice",
		pkg:    "types",
		types1: t,
	}
}

const (
	flagNone = 0 << iota
	// skipChild indicates that the child of the current element should be skipped
	// and directly handled by the function handleElement.
	skipChild = 1 << iota
	// forceChild is a flag that forces the child of the current child to be used
	// A.K.A., skip the current child and use the child of the current child
	// Useful for the case of list with custom tag
	forceChild = 2 << iota
)

func handleElement(e *xml.Element, st *StructInfo, flag uint) error {
	n := e
	for n != nil {
		var t any
		if n.Child != nil {
			// Skip the "li" tag since it's a slice and should not be a member of the struct
			if flag&skipChild != 0 || isListTag(n.GetName()) {
				flag &^= skipChild
				if err := handleElement(n.Child, st, flag); err != nil {
					return err
				}
			} else {
				switch n.Child.GetName() {
				case "li":
					t = createCustomSlice(n, flag|skipChild)
				case "keys":
					t = createCustomTypeForMap(n, flag)
				default:
					// Sometimes, slice are not marked as "li" so we need to check
					// if the next children has the same name.
					// If so, we consider it as a slice
					if n.Child.Next != nil && n.Child.Next.GetName() == n.Child.GetName() {
						t = createCustomSlice(n, flag|forceChild)
					} else {
						t = createStructure(n, flag)
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
