package generator

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/cruffinoni/rimworld-editor/helper"
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
)

type Member struct {
	T    any
	Attr attributes.Attributes
	Name string
}

type StructInfo struct {
	Name    string
	Members map[string]*Member
	Order   []*Member
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
	// Sometimes, there is no name attributed to a multiple grouped data and most
	// likely happens in lists. It's not possible for us to do the same thing.
	forceRandomName
	// Force to make a full check of all values in a list. This is persistent for lists
	// because a structure may vary from a one to another.
	forceFullCheck

	// InnerKeyword is the keyword for cases when the name of the element is
	// the same as the name of the parent.
	InnerKeyword = "_Inner"

	// Complex is a custom kind that is used when the type is too complex to be handled.
	// For example, *types.Slice[*types.Slice[*types.Slice[[3]*Struct]]] is too complex and is not detected now.
	Complex reflect.Kind = 100
)

var UniqueNumber = int64(0)

func addUniqueNumber(name string) string {
	name += strconv.FormatInt(UniqueNumber, 10)
	UniqueNumber++
	return name
}

// It seems sometimes "node", "nodes" and "subNodes" are used multiple times at multiple levels inside the same file
// so let's consider those names like non-unique and local to their context
func needToForceRandomName(name string) bool {
	switch name {
	case "node":
	case "nodes":
	case "subnodes":
		return true
	}
	return false
}

const maxTransversalDepth = 10

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
			// The array authorize cells to be empty
			//panic("generate.createStructure|forceChild: missing child")
		}
	}
	if e.Child == nil {
		panic("generate.createStructure: missing child")
	}
	name := e.GetName()
	lowerName := strings.ToLower(name)

	if needToForceRandomName(lowerName) && (flag&forceRandomName) == 0 {
		flag |= forceRandomName
	} else if helper.IsListTag(name) {
		// This case comes when the tag is an innerList of a list which can happen multiple times
		// in the file, so we need to set it a random name

		//log.Printf("generate.createStructure: '%s' & child name: %v & e %v & %v", name, e.Child.GetName(), lowerName, e.Parent.GetName())
		return createStructure(e.Parent, flag|forceRandomName)
	}

	// In this case, the child has the same name as his parent which
	// is very confusing for structure names.
	p := e.Parent
	if p != nil && name == p.GetName() {
		for p != nil {
			name += "_" + p.GetName()
			p = p.Parent
		}
		name += InnerKeyword
	}
	// vals is a special case where it serves as a transversal tag
	if (name == "vals" || name == "values" || strings.Contains(lowerName, "inner")) && e.Parent != nil {
		//log.Printf("Special case for: %v = %v", name, e.Parent.GetName()+"_"+name)
		depth := 0
		parent := e.Parent
		for parent != nil && depth < maxTransversalDepth {
			depth++
			name = parent.GetName() + "_" + name
			parent = parent.Parent
		}
	}
	if (flag & forceRandomName) > 0 {
		flag &^= forceRandomName
		name = addUniqueNumber(name)
	}
	s := &StructInfo{
		Name:    name,
		Members: make(map[string]*Member),
	}
	// The forceFullCheck check apply only to this structure, not to the children

	// If "forceFullCheck" is asked, it means we are in a slice/map, and we want
	// to check all nodes to have all members possible
	if (flag & forceFullCheck) > 0 {
		n := e

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
	} else {
		if err := handleElement(e.Child, s, flag&^forceFullCheck); err != nil {
			panic(err)
		}
	}
	flag &^= forceChildApplied
	return s
}

// addMember adds a new Member to the StructInfo map.
// If the Member already exists, the function checks if the type of the existing Member and the new Member are the same.
// If they are not, the function fixes the type mismatch.
func (s *StructInfo) addMember(name string, attr attributes.Attributes, t any) {
	// If there is no existing Member with the same name, add the new Member to the map
	if _, ok := s.Members[name]; !ok {
		s.Members[name] = &Member{
			T:    t,
			Attr: attr,
			Name: name,
		}
		s.Order = append(s.Order, s.Members[name])
	} else {
		// Check if the existing Member and the new Member are of the same type
		if !IsSameType(t, s.Members[name].T, 0) {
			// log.Printf("Type mismatch: %v > %v | %v", name, s.Members[name].T, t)
			// If the types are different, fix the type mismatch
			err := fixTypeMismatch(s.Members[name], &Member{
				Name: name,
				T:    t,
				Attr: attr,
			})
			if err != nil {
				panic(err)
			}
		}
	}
}
