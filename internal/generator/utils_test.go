package generator

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cruffinoni/rimworld-editor/file"
	"github.com/cruffinoni/rimworld-editor/internal/xml"
	"github.com/cruffinoni/rimworld-editor/internal/xml/attributes"
)

type args struct {
	e          *xml.Element
	flag       uint
	o          *offset
	xmlContent string
}

type tests struct {
	args args
	want any
}

func resetVarsAndReadBuffer(t *testing.T, args args) *xml.Element {
	UniqueNumber = 0
	RegisteredMembers = make(MemberVersioning)
	root, err := file.ReadFromBuffer(args.xmlContent)
	require.Nil(t, err)
	require.NotNil(t, root)
	return root.XML.Root
}

type emptyStructWithAttr struct {
	Attr attributes.Attributes
}

func createStructForTest(name string, m map[string]*Member) *StructInfo {
	for _, v := range m {
		if v.Attr == nil {
			v.Attr = make(attributes.Attributes)
		}
	}
	return &StructInfo{
		Name:    name,
		Members: m,
	}
}
