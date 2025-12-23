package codegen

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cruffinoni/rimworld-editor/internal/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/internal/xml/domain"
	"github.com/cruffinoni/rimworld-editor/internal/xml/loader"
)

type args struct {
	e          *domain.Element
	flag       uint
	o          *offset
	xmlContent string
}

type tests struct {
	args args
	want any
}

func resetVarsAndReadBuffer(t *testing.T, args args) *domain.Element {
	UniqueNumber = 0
	RegisteredMembers = make(MemberVersioning)
	root, err := loader.ReadFromBuffer(args.xmlContent)
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
