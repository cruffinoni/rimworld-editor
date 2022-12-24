package generator

import (
	"log"
	"reflect"

	"github.com/cruffinoni/rimworld-editor/helper"
	"github.com/cruffinoni/rimworld-editor/xml"
)

// getTypeFromArray returns the type of the element as a reflect.Kind
// It returns reflect.Invalid if the element is not a valid type
func getTypeFromArray(e *xml.Element) reflect.Kind {
	k := e.Child
	if k != nil && k.Child != nil {
		if helper.IsListTag(k.Child.GetName()) {
			return reflect.Slice
		}
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
			log.Printf("2: invalid type: %v => '%v' (data: '%v')", t, e.XMLPath(), e.Data)
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

func handleElement(e *xml.Element, st *StructInfo, flag uint) error {
	n := e
	log.Printf("St (%s) => %s & %v", st.name, e.GetName(), n.Child)
	for n != nil {
		var t any
		if n.Child != nil {
			// Skip the "li" tag since it's a slice and should not be a member of the struct
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
				if st == nil {
					st = t.(*StructInfo)
				} else {
					st.members = append(st.members, &member{
						name: n.GetName(),
						t:    t,
						attr: n.Attr,
					})
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
