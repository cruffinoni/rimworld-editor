package codegen

import (
	"github.com/cruffinoni/rimworld-editor/internal/codegen/importpath"
	"github.com/cruffinoni/rimworld-editor/internal/xml/domain"
	"github.com/cruffinoni/rimworld-editor/internal/xml/support"
	"github.com/cruffinoni/rimworld-editor/pkg/logging"
)

func createArrayOrSlice(logger logging.Logger, e *domain.Element, flag uint) any {
	k := e.Child
	count := 0
	for k != nil {
		count++
		if k.Data == nil && k.Child == nil && (count > 0 || k.Next != nil && k.Next.Next == nil) {
			// Count must be > 0 because empty slice/array must be considered as slice
			return createFixedArray(logger, e, flag, &offset{
				el:   k,
				size: count - 1, // -1 maybe??
			})
		}
		k = k.Next
	}
	return createCustomSlice(logger, e, flag)
}

const BasicStructName = "GeneratedStructStarter"

func createTypeFromElement(logger logging.Logger, n *domain.Element, flag uint) any {
	childName := n.Child.GetName()
	if support.IsListTag(childName) {
		return createArrayOrSlice(logger, n, flag)
	} else if childName == "keys" {
		return createCustomTypeForMap(logger, n, flag)
	} else if n.Child.Next != nil && n.Child.Next.GetName() == childName {
		return createArrayOrSlice(logger, n, flag|forceChild)
	} else {
		return createStructure(logger, n, flag)
	}
}

func processLeafNode(logger logging.Logger, n *domain.Element, st *StructInfo, flag uint) {
	var t any
	if n.Data != nil {
		t = n.Data.Kind()
		if !n.Attr.Empty() {
			t = &CustomType{
				Name:       "Type",
				Pkg:        "*embedded",
				Type1:      t,
				ImportPath: importpath.EmbeddedTypePath,
			}
		}
	} else if n.Next != nil && n.Next.GetName() == n.GetName() {
		t = createArrayOrSlice(logger, n, flag)
		for n.Next != nil && n.Next.GetName() == n.GetName() {
			n = n.Next
		}
	} else {
		t = createEmptyType()
	}
	st.addMember(logger, n.GetName(), n.Attr, t)
}

func handleElement(logger logging.Logger, e *domain.Element, st *StructInfo, flag uint) error {
	n := e
	//if n != nil && n.GetName() == "li" {
	//	printer.Debugf("n: %v", n.GetName())
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
			if support.IsListTag(n.GetName()) {
				if err := handleElement(logger, n.Child, st, flag); err != nil {
					return err
				}
			} else {
				st.addMember(logger, n.GetName(), n.Attr, createTypeFromElement(logger, n, flag))
			}
		} else if !support.IsListTag(n.GetName()) {
			processLeafNode(logger, n, st, flag)
		} else {
			t := createArrayOrSlice(logger, n, flag)
			st.addMember(logger, n.GetName(), n.Attr, t)
		}
		n = n.Next
	}
	RegisteredMembers[st.Name] = append(RegisteredMembers[st.Name], st)
	return nil
}
