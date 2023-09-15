package generator

import (
	"log"
	"reflect"

	"github.com/cruffinoni/rimworld-editor/generator/paths"
	"github.com/cruffinoni/rimworld-editor/helper"
	"github.com/cruffinoni/rimworld-editor/xml"
)

func determineArrayOrSliceKind(e *xml.Element) reflect.Kind {
	k := e
	for k != nil {
		if k.Data == nil && k.Child == nil {
			return reflect.Array
		}
		k = k.Next
	}
	return reflect.Slice
}

// getTypeFromArray returns the type of the element as a reflect.Kind.
// If the element is not a valid type, it returns reflect.Invalid.
func getTypeFromArray(e *xml.Element) reflect.Kind {
	// Return Invalid if the element has no child
	if e.Child == nil {
		if e.Data == nil {
			return reflect.Invalid
		} else {
			return e.Data.Kind()
		}
	}

	k := e.Child
	kt := reflect.Invalid

	// Determine if the element is a structure or slice
	// In some cases, the first structure may have no data, so we check its siblings for a child
	// Note to run the check until we find a valid sibling with a child.
	// Think about the case where only the last value has a child of a huge array
	if k.Next != nil && k.Next.GetName() == k.GetName() {
		// This part of code is aimed to list with multiple elements
		n := k.Next
		for n != nil && n.GetName() == k.GetName() {
			if n.Child != nil {
				if helper.IsListTag(n.Child.GetName()) {
					return determineArrayOrSliceKind(n)
				}
				return reflect.Struct
			}
			n = n.Next
		}
	}
	if k.Child != nil {
		// On the other hand, this part only check the first element
		if helper.IsListTag(k.Child.GetName()) {
			return determineArrayOrSliceKind(k)
		}
		return reflect.Struct
	}

	for k != nil {
		if k.Data != nil {
			kdk := k.Data.Kind()
			if kt != reflect.Invalid && kdk != kt &&
				// Float64 and Int64 can be interchangeable
				!(kdk == reflect.Float64 && kt == reflect.Int64) &&
				!(kdk == reflect.Int64 && kt == reflect.Float64) {
				log.Printf("Data: '%v' & kind %s", k.Data, k.Data.Kind())
				log.Panicf("getTypeFromArray: found type %v, expected %v on path %v ('%v')", kdk, kt, k.XMLPath(), k.Data.GetData())
			}
			// Float64 and Int64 can be interchangeable, but we prefer to keep Float64
			if !(kt == reflect.Float64 && kdk == reflect.Int64) {
				kt = kdk
			}
		}
		k = k.Next
	}
	return kt
}

func createArrayOrSlice(e *xml.Element, flag uint) any {
	k := e.Child
	count := 0
	for k != nil {
		count++
		if k.Data == nil && k.Child == nil && (count > 0 || k.Next != nil && k.Next.Next == nil) {
			// Count must be > 0 because empty slice/array must be considered as slice
			return createFixedArray(e, flag, &offset{
				el:   k,
				size: count,
			})
		}
		k = k.Next
	}
	return createCustomSlice(e, flag)
}

// determineTypeFromData returns the type of data from the element.
// If the element is not a primitive type, it returns either a
// StructInfo or a CustomType. Otherwise, it returns the type of the
// element as a reflect.Kind
// If the type is invalid, we consider it as a *xml.Element.
// Is it useful for empty tags.
func determineTypeFromData(e *xml.Element, flag uint) any {
	t := any(getTypeFromArray(e))
	//log.Printf("Type of %v is %v", e.XMLPath(), t)
	// We need to define this struct with of this all members
	if t == reflect.Struct || t == reflect.Slice {
		c := e.Child
		// If the child is a list, let's create a slice from it
		if helper.IsListTag(c.Child.GetName()) {
			// We set the forceChild flag to true to force the function createStructure
			// to take the children of the list and not the list itself.
			t = createArrayOrSlice(c, flag|forceChild)
		} else {
			// Otherwise, a basic struct is created
			// We pass 'e' instead of 'c' because createStructure will take the children of 'e'
			t = createStructure(e, flag)
		}
	} else if t == reflect.Invalid {
		if e.Data == nil {
			t = createEmptyType()
		} else {
			t = e
		}
	} else {
		if !e.Attr.Empty() {
			//log.Println("primary.EmbeddedType: found attributes on path", e.XMLPath())
			return &CustomType{
				Name:       "Type",
				Pkg:        "embedded",
				Type1:      t,
				ImportPath: paths.EmbeddedTypePath,
			}
		}
	}
	return t
}

const BasicStructName = "GeneratedStructStarter"

func processSpecialCase(st *StructInfo) {
	*st = StructInfo{
		Name:    addUniqueNumber(BasicStructName),
		Members: make(map[string]*Member),
		Order:   make([]*Member, 0),
	}
}

func processChildNode(n *xml.Element, st *StructInfo, flag uint) error {
	var t any
	childName := n.Child.GetName()
	if helper.IsListTag(childName) {
		t = createArrayOrSlice(n, flag)
	} else if childName == "keys" {
		t = createCustomTypeForMap(n, flag)
	} else if n.Child.Next != nil && n.Child.Next.GetName() == childName {
		t = createArrayOrSlice(n, flag|forceChild)
	} else {
		t = createStructure(n, flag)
	}
	st.addMember(n.GetName(), n.Attr, t)
	return nil
}

func processLeafNode(n *xml.Element, st *StructInfo, flag uint) {
	var t any
	if n.Data != nil {
		t = n.Data.Kind()
		if !n.Attr.Empty() {
			t = &CustomType{
				Name:       "Type",
				Pkg:        "*embedded",
				Type1:      t,
				ImportPath: paths.EmbeddedTypePath,
			}
		}
	} else if n.Next != nil && n.Next.GetName() == n.GetName() {
		t = createArrayOrSlice(n, flag)
		for n.Next != nil && n.Next.GetName() == n.GetName() {
			n = n.Next
		}
	} else {
		t = createEmptyType()
	}
	st.addMember(n.GetName(), n.Attr, t)
}

func handleElement(e *xml.Element, st *StructInfo, flag uint) error {
	n := e
	//log.Printf("ST: %v", st.Name)
	//if n != nil && n.GetName() == "li" {
	//	log.Printf("n: %v", n.GetName())
	//}
	if st.Name == "" {
		processSpecialCase(st)
	}
	for n != nil {
		if n.Child != nil {
			if helper.IsListTag(n.GetName()) {
				if err := handleElement(n.Child, st, flag); err != nil {
					return err
				}
			} else {
				if err := processChildNode(n, st, flag); err != nil {
					return err
				}
			}
		} else if !helper.IsListTag(n.GetName()) {
			processLeafNode(n, st, flag)
		} else {
			t := createArrayOrSlice(n, flag)
			st.addMember(n.GetName(), n.Attr, t)
		}
		n = n.Next
	}
	RegisteredMembers[st.Name] = append(RegisteredMembers[st.Name], st)
	return nil
}
