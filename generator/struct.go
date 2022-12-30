package generator

import (
	"bytes"
	"os"
	"strconv"

	"github.com/cruffinoni/rimworld-editor/helper"
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
)

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

const (
	flagNone = 0 << iota
	// skipChild indicates that the child of the current element should be skipped
	// and directly handled by the function handleElement.
	skipChild = 1 << iota
	// forceChild is a flag that forces the child of the current child to be used
	// A.K.A., skip the current child and use the child of the current child
	// Useful for the case of list with custom tag
	forceChild = 2 << iota

	ignoreSlice = 3 << iota

	forceRandomName = 4 << iota

	// InnerKeyword is the keyword for cases when the name of the element is the same as the name of the parent.
	InnerKeyword = "_Inner"
)

// GenerateGoFiles generates the Go files (with the corresponding structs)
// for the given XML file, but it doesn't write anything.
// To do that, call WriteGoFile.
func GenerateGoFiles(root *xml.Element) *StructInfo {
	s := &StructInfo{}
	if err := handleElement(root, s, flagNone); err != nil {
		panic(err)
	}
	return s
}

// WriteGoFile writes the struct Go code to the given path.
// It writes recursively the members of the struct. If a member is a struct,
// it will call WriteGoFile on it.
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
	return s.generateStructToPath(path)
}

// removeDuplicates removes the duplicates from the members of the struct.
func (s *StructInfo) removeDuplicates() {
	for i := 0; i < len(s.members); i++ {
		for j := i + 1; j < len(s.members); j++ {
			if s.members[i].name == s.members[j].name {
				s.members = append(s.members[:i], s.members[i+1:]...)
				i--
				break
			}
		}
	}
}

var uniqueNumber = 0

// createStructure creates a new structure from the given element.
// Then the function will recursively call handleElement on the children of the element.
// It removes the duplicates from the members of the struct.
func createStructure(e *xml.Element, flag uint) any {
	// forceChild is a flag that forces the child of the current child to be used
	// It is useful for the case of lists
	if flag&forceChild == forceChild {
		flag &^= forceChild
		//log.Println("Forcing child flag")
		// Quick way to determine if the child is a structure
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

	// TODO: Update doc for this line of code
	if helper.IsListTag(name) {
		//log.Printf("generate.createStructure: '%s' & child name: %v", name, e.Child.GetName())
		return createStructure(e.Parent, flag|forceRandomName)
	}

	// In this case, the child has the same name as his parent which
	// is very confusing for structure names.
	if e.Parent != nil && name == e.Parent.GetName() {
		name += InnerKeyword
	}
	if (flag & forceRandomName) > 0 {
		name += strconv.Itoa(uniqueNumber)
		uniqueNumber++
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
