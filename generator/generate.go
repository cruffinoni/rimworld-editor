package generator

import (
	"log"
	"reflect"

	"github.com/cruffinoni/rimworld-editor/generator/paths"
	"github.com/cruffinoni/rimworld-editor/helper"
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
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
		if k.Data == nil && k.Child == nil && (count > 0 || k.Next != nil && k.Next.Next == nil) {
			// Count must be > 0 because empty slice/array must be considered as
			// slice
			return createFixedArray(e, flag, &offset{
				el:   k,
				size: count,
			})
		}
		count++
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
		if (flag & ignoreSlice) > 0 {
			flag &^= ignoreSlice
			return determineTypeFromData(c, flag)
		}
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
				Name:       "EmbeddedType",
				Pkg:        "primary",
				Type1:      t,
				ImportPath: paths.PrimaryTypesPath,
			}
		}
	}
	return t
}

func hasSameMembers(b, a *StructInfo, depth uint32) bool {
	if len(a.Members) != len(b.Members) {
		return false
	}
	for i := range a.Members {
		if !isSameType(b.Members[i], a.Members[i], depth+1) {
			return false
		}
	}
	return true
}

func handleElement(e *xml.Element, st *StructInfo, flag uint) error {
	n := e
	for n != nil {
		var t any
		if n.Child != nil {
			// Skip the "li" tag (or any custom type) since it's a slice and should not be a member of the struct
			if helper.IsListTag(n.GetName()) {
				if err := handleElement(n.Child, st, flag); err != nil {
					return err
				}
			} else {
				childName := n.Child.GetName()
				if helper.IsListTag(childName) {
					t = createArrayOrSlice(n, flag)
				} else if childName == "keys" {
					// Maps are constant in terms of naming, and that's how we recognize them
					t = createCustomTypeForMap(n, flag)
				} else {
					// Sometimes, slice are not marked as "li" so we need to check
					// if the next sibling has the same name.
					// If so, we consider it as a slice
					if n.Child.Next != nil && n.Child.Next.GetName() == childName {
						t = createArrayOrSlice(n, flag|forceChild)
					} else {
						t = createStructure(n, flag)
					}
				}
				// This is a special case where the root node has been created outside the process.
				// To recognize this special node, we don't set any name to it, but it refers as the root node.
				if st.Name == "" {
					*st = *t.(*StructInfo)
				} else {
					st.addMember(n.GetName(), n.Attr, t)
				}
			}
		} else if !helper.IsListTag(n.GetName()) {
			if n.Data != nil {
				t = n.Data.Kind()
				if !n.Attr.Empty() {
					t = &CustomType{
						Name:       "EmbeddedType",
						Pkg:        "*primary",
						Type1:      t,
						ImportPath: paths.PrimaryTypesPath,
					}
				}
			} else if n.Next != nil && n.Next.GetName() == n.GetName() {
				// This condition must not apply to list tags because empty elements in a list are valid
				t = createArrayOrSlice(n, flag)
				// Skip the next element since it's already handled
				for n.Next != nil && n.Next.GetName() == n.GetName() {
					n = n.Next
				}
			} else {
				t = createEmptyType()
			}
			//log.Printf("Add member %T w/ %v (%s) to %v", t, n.XMLPath(), n.GetName(), st.Name)
			st.addMember(n.GetName(), n.Attr, t)
		} else {
			// If we reach this code section, it means that we may have a list with no child, an empty list.
			t = createArrayOrSlice(n, flag)
			st.addMember(n.GetName(), n.Attr, t)
		}
		n = n.Next
	}
	RegisteredMembers[st.Name] = append(RegisteredMembers[st.Name], st)
	//if _, ok := RegisteredMembers[st.Name]; ok {
	//if !hasSameMembers(m, st) {
	//log.Printf("WARNING: struct %s (length %d - %p) is different from %s (length %d - %p)", m.Name, len(m.Members), m, st.Name, len(st.Members), st)
	//fixMembers(m, st)
	//}
	//} else {
	//	RegisteredMembers[st.Name] = st
	//}
	return nil
}

// addMember adds a new member to the StructInfo map.
// If the member already exists, the function checks if the type of the existing member and the new member are the same.
// If they are not, the function fixes the type mismatch.
func (s *StructInfo) addMember(name string, attr attributes.Attributes, t any) {
	// If there is no existing member with the same name, add the new member to the map
	if _, ok := s.Members[name]; !ok {
		s.Members[name] = &member{
			T:    t,
			Attr: attr,
			Name: name,
		}
		s.Order = append(s.Order, s.Members[name])
	} else {
		// Check if the existing member and the new member are of the same type
		if !isSameType(t, s.Members[name].T, 0) {
			//log.Printf("Type mismatch: %v > %v", name, s.Members[name])
			// If the types are different, fix the type mismatch
			fixTypeMismatch(s.Members[name], &member{
				Name: name,
				T:    t,
				Attr: attr,
			})
		}
	}
}
