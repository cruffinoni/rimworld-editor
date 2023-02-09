package generator

import (
	"log"
	"reflect"

	"github.com/cruffinoni/rimworld-editor/helper"
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
)

// getTypeFromArray returns the type of the element as a reflect.Kind.
// If the element is not a valid type, it returns reflect.Invalid.
func getTypeFromArray(e *xml.Element) reflect.Kind {
	// Return Invalid if the element has no child
	if e.Child == nil {
		return reflect.Invalid
	}

	k := e.Child
	kt := reflect.Invalid

	// Determine if the element is a structure or slice
	// In some cases, the first structure may have no data, so we check its siblings for a child
	// We only need to do this once
	if k.Next != nil && k.Next.GetName() == k.GetName() {
		// This part of code is aimed to list with multiple elements
		if k.Next.Child != nil {
			if helper.IsListTag(k.Next.Child.GetName()) {
				return reflect.Slice
			}
			return reflect.Struct
		}
	} else if k.Child != nil {
		// On the other hand, this part only check the first element
		if helper.IsListTag(k.Child.GetName()) {
			return reflect.Slice
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
				log.Panicf("primary.EmbeddedType: found type %v, expected %v on path %v ('%v')", kdk, kt, k.XMLPath(), k.Data.GetData())
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
			//log.Println("Ignoring slice")
			return determineTypeFromData(c, flag)
		}
		// If the child is a list, let's create a slice from it
		if helper.IsListTag(c.Child.GetName()) {
			// We set the forceChild flag to true to force the function createStructure
			// to take the children of the list and not the list itself.
			t = createCustomSlice(c, flag|forceChild)
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
				name:       "EmbeddedType",
				pkg:        "primary",
				type1:      t,
				importPath: primaryTypesPath,
			}
		}
	}
	return t
}

func hasSameMembers(a, b *StructInfo) bool {
	if len(a.members) != len(b.members) {
		return false
	}
	for i := range a.members {
		if b.members[i] == nil {
			return false
		}
		if reflect.TypeOf(a.members[i]) != reflect.TypeOf(b.members[i]) {
			return false
		}
		if ct, ok := a.members[i].t.(*CustomType); ok {
			if ctB, okB := b.members[i].t.(*CustomType); !okB || ct.type1 != ctB.type1 || ct.type2 != ctB.type2 {
				return false
			}
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
			if flag&skipChild != 0 || helper.IsListTag(n.GetName()) {
				flag &^= skipChild
				if err := handleElement(n.Child, st, flag); err != nil {
					return err
				}
			} else {
				childName := n.Child.GetName()
				if helper.IsListTag(childName) {
					t = createCustomSlice(n, flag|skipChild)
				} else if childName == "keys" {
					// Maps are constant in terms of naming, and that's how we recognize them
					t = createCustomTypeForMap(n, flag)
				} else {
					// Sometimes, slice are not marked as "li" so we need to check
					// if the next sibling has the same name.
					// If so, we consider it as a slice
					if n.Child.Next != nil && n.Child.Next.GetName() == childName {
						t = createCustomSlice(n, flag|forceChild)
					} else {
						t = createStructure(n, flag)
					}
				}
				// This is a special case where the root node has been created outside the process.
				// To recognize this special node, we don't set any name to it, but it refers as the root node.
				if st.name == "" {
					*st = *t.(*StructInfo)
				} else {
					st.addMember(n.GetName(), n.Attr, t)
				}
			}
		} else {
			if n.Data != nil {
				t = n.Data.Kind()
				if !n.Attr.Empty() {
					t = &CustomType{
						name:       "EmbeddedType",
						pkg:        "primary",
						type1:      t,
						importPath: primaryTypesPath,
					}
				}
			} else if n.Next != nil && n.Next.GetName() == n.GetName() {
				t = createCustomSlice(n, flag)
				// Skip the next element since it's already handled
				for n.Next != nil && n.Next.GetName() == n.GetName() {
					n = n.Next
				}
			} else {
				t = createEmptyType()
			}
			st.addMember(n.GetName(), n.Attr, t)
		}
		n = n.Next
	}
	if m, ok := registeredMembers[st.name]; ok && !hasSameMembers(m, st) {
		//if m.name == "thing" {
		//	log.Printf("WARNING: struct %s (length %d - %p) is different from %s (length %d - %p)", m.name, len(m.members), m, st.name, len(st.members), st)
		//}
		fixMembers(m, st)
		//if m.name == "thing" {
		//	log.Printf("WARNING: struct %s (length %d - %p) is different from %s (length %d - %p)", m.name, len(m.members), m, st.name, len(st.members), st)
		//}
	} else {
		registeredMembers[st.name] = st
	}
	return nil
}

// addMember adds a new member to the StructInfo map.
// If the member already exists, the function checks if the type of the existing member and the new member are the same.
// If they are not, the function fixes the type mismatch.
func (s *StructInfo) addMember(name string, attr attributes.Attributes, t any) {
	// If there is no existing member with the same name, add the new member to the map
	if _, ok := s.members[name]; !ok {
		s.members[name] = &member{
			t:    t,
			attr: attr,
		}
	} else {
		// Check if the existing member and the new member are of the same type
		if kind, okKind := s.members[name].t.(reflect.Kind); !isSameType(s.members[name].t, t) || (okKind && kind != t.(reflect.Kind)) {
			//log.Printf("Type mismatch: %v > %v", name, s.members[name])
			// If the types are different, fix the type mismatch
			fixTypeMismatch(s.members[name], &member{
				t:    t,
				attr: attr,
			})
		}
	}
}
