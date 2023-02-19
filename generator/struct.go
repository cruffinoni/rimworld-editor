package generator

import (
	"strconv"
	"strings"

	"github.com/cruffinoni/rimworld-editor/helper"
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
)

type member struct {
	T    any
	Attr attributes.Attributes
	Name string
}

type StructInfo struct {
	Name    string
	Members map[string]*member
	Order   []*member
}

const (
	flagNone = 0 << iota
	// skipChild indicates that the child of the current element should be skipped
	// and directly handled by the function handleElement.
	forceChildApplied = 1 << iota
	// forceChild is a flag that forces the child of the current child to be used
	// A.K.A., skip the current child and use the child of the current child
	// Useful for the case of list with custom tag
	forceChild

	// Ignore the slice tag (li) because we don't need it there. Used for maps where
	// only values inside the li tag is interesting
	ignoreSlice

	// Sometimes, there is no name attributed to a multiple grouped data and most
	// likely happens in lists. It's not possible for us to do the same thing.
	forceRandomName

	// Force to make a full check of all values in a list. This is persistent for lists
	// because a structure may vary from a one to another.
	forceFullCheck

	// InnerKeyword is the keyword for cases when the name of the element is
	// the same as the name of the parent.
	InnerKeyword = "_Inner"
)

var RegisteredMembers = make(map[string]*StructInfo)

// GenerateGoFiles generates the Go files (with the corresponding structs)
// for the given XML file, but it doesn't write anything.
// To do that, call WriteGoFile.
func GenerateGoFiles(root *xml.Element) *StructInfo {
	s := &StructInfo{
		Members: make(map[string]*member),
	}
	RegisteredMembers = make(map[string]*StructInfo)
	if err := handleElement(root, s, flagNone); err != nil {
		panic(err)
	}
	return s
}

var UniqueNumber = 0

// createStructure creates a new structure from the given element.
// Then the function will recursively call handleElement on the children of the element.
// It removes the duplicates from the members of the struct.
func createStructure(e *xml.Element, flag uint) any {
	// forceChild is a flag that forces the child of the current child to be used
	// It is useful for the case of lists
	if flag&forceChild > 0 {
		flag &^= forceChild
		// Quick way to determine if the child is a structure
		if e.Child != nil && e.Child.Child != nil {
			return createStructure(e.Child, flag|forceChildApplied)
		} else {
			panic("generate.createStructure|forceChild: missing child")
		}
	}
	if e.Child == nil {
		panic("generate.createStructure: missing child")
	}
	name := e.GetName()

	// This case comes when the tag is an innerList of a list which can happen multiple times
	// in the file, so we need to set it a random name
	if helper.IsListTag(name) {
		//log.Printf("generate.createStructure: '%s' & child name: %v", name, e.Child.GetName())
		return createStructure(e.Parent, flag|forceRandomName)
	}

	// In this case, the child has the same name as his parent which
	// is very confusing for structure names.
	if e.Parent != nil && name == e.Parent.GetName() {
		name += InnerKeyword
	}
	// vals is a special case where it serves as a transversal tag
	if (name == "vals" || strings.Contains(strings.ToLower(name), "inner")) && e.Parent != nil {
		//log.Printf("Special case for: %v = %v", name, e.Parent.GetName()+"_"+name)
		name = e.Parent.GetName() + "_" + name
	}
	if (flag & forceRandomName) > 0 {
		flag &^= forceRandomName
		name += strconv.Itoa(UniqueNumber)
		UniqueNumber++
	}
	s := &StructInfo{
		Name:    name,
		Members: make(map[string]*member),
	}
	// The forceFullCheck check apply only to this structure, not to the children
	if err := handleElement(e.Child, s, flag&^forceFullCheck); err != nil {
		panic(err)
	}
	// If "forceFullCheck" is asked, it means we are in a slice/map, and we want
	// to check all nodes to have all members possible
	if (flag & forceFullCheck) > 0 {
		n := e
		//log.Printf("Forcefullcheck on %s & %p", e.XMLPath(), n.Child)

		// forceChildApplied has been applied and so, we are in the children level and not
		// in the main structure level
		if n.Child != nil && forceChildApplied&flag == 0 {
			n = n.Child
		}
		flag &^= forceFullCheck | forceChildApplied
		for n != nil {
			if err := handleElement(n.Child, s, flag); err != nil {
				panic(err)
			}
			n = n.Next
		}
	}
	flag &^= forceChildApplied
	return s
}
