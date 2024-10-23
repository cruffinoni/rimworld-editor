package generator

import (
	"github.com/cruffinoni/rimworld-editor/helper"

	"github.com/cruffinoni/rimworld-editor/internal/generator/paths"

	"github.com/cruffinoni/rimworld-editor/internal/xml"
)

func createArrayOrSlice(e *xml.Element, flag uint) any {
	k := e.Child
	count := 0
	for k != nil {
		count++
		if k.Data == nil && k.Child == nil && (count > 0 || k.Next != nil && k.Next.Next == nil) {
			// Count must be > 0 because empty slice/array must be considered as slice
			return createFixedArray(e, flag, &offset{
				el:   k,
				size: count - 1, // -1 maybe??
			})
		}
		k = k.Next
	}
	return createCustomSlice(e, flag)
}

const BasicStructName = "GeneratedStructStarter"

func createTypeFromElement(n *xml.Element, flag uint) any {
	childName := n.Child.GetName()
	if helper.IsListTag(childName) {
		return createArrayOrSlice(n, flag)
	} else if childName == "keys" {
		return createCustomTypeForMap(n, flag)
	} else if n.Child.Next != nil && n.Child.Next.GetName() == childName {
		return createArrayOrSlice(n, flag|forceChild)
	} else {
		return createStructure(n, flag)
	}
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
	//if n != nil && n.GetName() == "li" {
	//	log.Printf("n: %v", n.GetName())
	//}
	if st.Name == "" {
		*st = StructInfo{
			Name:    addUniqueNumber(BasicStructName),
			Members: make(map[string]*Member),
			Order:   make([]*Member, 0),
		}
	}
	for n != nil {
		if n.Child != nil {
			if helper.IsListTag(n.GetName()) {
				if err := handleElement(n.Child, st, flag); err != nil {
					return err
				}
			} else {
				st.addMember(n.GetName(), n.Attr, createTypeFromElement(n, flag))
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
