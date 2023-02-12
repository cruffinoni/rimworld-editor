package generator

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cruffinoni/rimworld-editor/file"
	"github.com/cruffinoni/rimworld-editor/generator/paths"
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
)

type args struct {
	e          *xml.Element
	flag       uint
	o          *offset
	xmlContent string
}

type tests struct {
	args args
	want *FixedArray
}

func sameMembers(a, b map[string]*member) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if m, ok := b[k]; !ok {
			return false
		} else {
			t1 := reflect.TypeOf(m)
			if t1.Kind() == reflect.Ptr {
				t1 = t1.Elem()
			}
			t2 := reflect.TypeOf(v)
			if t2.Kind() == reflect.Ptr {
				t2 = t2.Elem()
			}
			if t1.Name() != t2.Name() {
				return false
			}
		}
	}
	return true
}

func Test_createFixedArray(t *testing.T) {
	tests := map[string]tests{
		"struct of struct": {
			args: args{
				xmlContent: `
<?xml version="1.0" encoding="utf-8"?>
<savegame>
	<vals>
		<li>
			<frequency>Normal</frequency>
		</li>
		<li/>
		<li>
			<frequency>Normal</frequency>
			<gender>Female</gender>
		</li>
	</vals>
</savegame>
`,
			},
			want: &FixedArray{
				Size: 3,
				PrimaryType: &StructInfo{
					Name: "vals",
					Members: map[string]*member{
						"frequency": {},
						"gender":    {},
					},
				},
			},
		},

		"empty": {
			args: args{
				xmlContent: `
<?xml version="1.0" encoding="utf-8"?>
<savegame>
	<vals>
		<li/>
		<li/>
		<li/>
	</vals>
</savegame>
`,
			},
			want: &FixedArray{
				Size: 3,
				PrimaryType: &StructInfo{
					Name:    "vals",
					Members: map[string]*member{},
				},
			},
		},

		"nested slice": {
			args: args{
				xmlContent: `
<?xml version="1.0" encoding="utf-8"?>
<savegame>
	<vals>
		<li>
			<technology>
				<li>1</li>
				<li>2</li>
				<li>3</li>
				<li>4</li>
			</technology>
		</li>
		<li/>
	</vals>
</savegame>
`,
			},
			want: &FixedArray{
				Size: 2,
				PrimaryType: &StructInfo{
					Name: "vals",
					Members: map[string]*member{
						"technology": {
							T: &CustomType{
								Name: "types",
								Pkg:  "Slice",
								Type1: &struct {
									Attr attributes.Attributes
								}{
									Attr: nil,
								},
								ImportPath: paths.CustomTypesPath,
							},
						},
					},
				},
			},
		},

		"nested array": {
			args: args{
				xmlContent: `
<?xml version="1.0" encoding="utf-8"?>
<savegame>
	<vals>
		<li>
			<technology>
				<li/>
				<li/>
				<li>
					<progession>100</progession>
				</li>
				<li/>
			</technology>
		</li>
		<li/>
	</vals>
</savegame>
`,
			},
			want: &FixedArray{
				Size: 2,
				PrimaryType: &StructInfo{
					Name: "vals",
					Members: map[string]*member{
						"technology": {
							T: &FixedArray{
								Size: 4,
								PrimaryType: &struct {
									Attr       attributes.Attributes
									Progession int64
								}{},
							},
						},
					},
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			UniqueNumber = 0
			RegisteredMembers = make(map[string]*StructInfo)
			root, err := file.ReadFromBuffer(tt.args.xmlContent)
			require.Nil(t, err)
			require.NotNil(t, root)
			got := createFixedArray(root.XML.Root.Child, tt.args.flag, tt.args.o)
			require.IsType(t, got, tt.want)
			gotCasted := got.(*FixedArray)
			assert.Equal(t, tt.want.Size, gotCasted.Size)
			require.IsTypef(t, tt.want.PrimaryType, gotCasted.PrimaryType, "expected %+v (%T), got %+v (%T)", tt.want.PrimaryType, tt.want.PrimaryType, gotCasted.PrimaryType, gotCasted.PrimaryType)
			if si, ok := tt.want.PrimaryType.(*StructInfo); ok && !sameMembers(si.Members, gotCasted.PrimaryType.(*StructInfo).Members) {
				assert.Fail(t, fmt.Sprintf("expected members to be %+v (%T), got %+v (%T)", si.Members, si.Members, gotCasted.PrimaryType.(*StructInfo).Members, gotCasted.PrimaryType.(*StructInfo).Members))
			}
		})
	}
}
