package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cruffinoni/rimworld-editor/internal/xml/attributes"
)

func Test_createFixedArray(t *testing.T) {
	tests := map[string]tests{
		"struct": {
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
					Members: map[string]*Member{
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
					Members: map[string]*Member{},
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
					Members: map[string]*Member{
						"technology": {
							T: createCustomSliceForTest(&emptyStructWithAttr{}),
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
					Members: map[string]*Member{
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
			root := resetVarsAndReadBuffer(t, tt.args)
			got := createFixedArray(root.Child, tt.args.flag, tt.args.o)
			require.IsType(t, got, tt.want)
			gotCasted := got.(*FixedArray)
			wanted := got.(*FixedArray)
			assert.Equal(t, wanted.Size, gotCasted.Size)
			require.IsTypef(t, wanted.PrimaryType, gotCasted.PrimaryType, "expected %+v (%T), got %+v (%T)", wanted.PrimaryType, wanted.PrimaryType, gotCasted.PrimaryType, gotCasted.PrimaryType)
			assert.Equal(t, wanted, got)
		})
	}
}

func createFixedArrayForTest(size int, pt any) *FixedArray {
	return &FixedArray{
		Size:        size,
		PrimaryType: pt,
	}
}
