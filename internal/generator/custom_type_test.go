package generator

import (
	"reflect"
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cruffinoni/rimworld-editor/internal/generator/paths"
)

func createCustomSliceForTest(type1 any) *CustomType {
	return &CustomType{
		Name:       "Slice",
		Pkg:        "*types",
		Type1:      type1,
		ImportPath: paths.CustomTypesPath,
	}
}

func createCustomMapForTest(t1, t2 any) *CustomType {
	return &CustomType{
		Name:       "Map",
		Pkg:        "types",
		Type1:      t1,
		Type2:      t2,
		ImportPath: paths.CustomTypesPath,
	}
}

func Test_createCustomType(t *testing.T) {
	tests := map[string]tests{
		"slice_primary": {
			args: args{
				xmlContent: `
<?xml version="1.0" encoding="utf-8"?>
<savegame>
	<vals>
		<li>hello</li>
		<li>world</li>
	</vals>
</savegame>
`,
			},
			want: createCustomSliceForTest(reflect.String),
		},

		"slice_struct": {
			args: args{
				xmlContent: `
<?xml version="1.0" encoding="utf-8"?>
<savegame>
	<vals>
		<li>
			<quests>
				<completed>True</completed>
			</quests>
		</li>
		<li>
			<quests>
				<completed>False</completed>
			</quests>
		</li>
	</vals>
</savegame>
`,
			},
			want: createCustomSliceForTest(createStructForTest("vals", map[string]*Member{
				"quests": {
					T: createStructForTest("quests", map[string]*Member{
						"completed": {
							T: reflect.String,
						},
					}),
				},
			})),
		},

		"slice_multiple struct": {
			args: args{
				xmlContent: `
<?xml version="1.0" encoding="utf-8"?>
<savegame>
	<vals>
		<li>
			<quests>
				<completed>True</completed>
				<progression>100</progression>
			</quests>
		</li>
		<li>
			<quests>
				<completed>False</completed>
			</quests>
		</li>
	</vals>
</savegame>
`,
			},
			want: createCustomSliceForTest(createStructForTest("vals", map[string]*Member{
				"quests": {
					T: createStructForTest("quests", map[string]*Member{
						"completed": {
							T: reflect.String,
						},
						"progression": {
							T: reflect.Int64,
						},
					}),
				},
			})),
		},

		"slice_struct start at second elem": {
			args: args{
				xmlContent: `
<?xml version="1.0" encoding="utf-8"?>
<savegame>
	<vals>
		<li>
			<quests />
		</li>
		<li>
			<quests>
				<completed>False</completed>
			</quests>
		</li>
	</vals>
</savegame>
`,
			},
			want: createCustomSliceForTest(createStructForTest("vals", map[string]*Member{
				"quests": {
					T: createStructForTest("quests", map[string]*Member{
						"completed": {
							T: reflect.String,
						},
					}),
				},
			})),
		},

		"slice_with empty": {
			args: args{
				xmlContent: `
<?xml version="1.0" encoding="utf-8"?>
<savegame>
	<vals>
		<li>
			<quests>
				<completed>True</completed>
				<progression>100</progression>
			</quests>
		</li>
		<li>
			<quests>
				<completed>False</completed>
				<guests Null="true" />
			</quests>
		</li>
	</vals>
</savegame>
`,
			},
			want: createCustomSliceForTest(createStructForTest("vals", map[string]*Member{
				"quests": {
					T: createStructForTest("quests", map[string]*Member{
						"completed": {
							T: reflect.String,
						},
						"progression": {
							T: reflect.Int64,
						},
						"guests": {
							T: createEmptyType(),
							Attr: map[string]string{
								"Null": "true",
							},
						},
					}),
				},
			})),
		},

		"slice_struct_variable": {
			args: args{
				xmlContent: `
<?xml version="1.0" encoding="utf-8"?>
<savegame>
	<things>
		<thing>
			<def>Filth_RubbleRock</def>
			<id>Filth_RubbleRock3236</id>
			<map>0</map>
			<pos>(184, 0, 3)</pos>
			<questTags />
			<disappearAfterTicks>2931683</disappearAfterTicks>
		</thing>
		<thing>
			<def>Filth_RubbleRock</def>
			<id>Filth_RubbleRock3237</id>
			<map>0</map>
			<pos>(185, 0, 3)</pos>
			<questTags />
			<disappearAfterTicks>2944394</disappearAfterTicks>
		</thing>
		<thing>
			<def>Plant_Grass</def>
			<id>Plant_Grass8968</id>
			<map>0</map>
			<pos>(95, 0, 245)</pos>
			<health>85</health>
			<questTags />
			<growth>0.649532</growth>
			<age>821805</age>
		</thing>
	</things>
</savegame>
`,
			},
			want: createCustomSliceForTest(createStructForTest("things", map[string]*Member{
				"id":                  {T: reflect.String},
				"pos":                 {T: reflect.String},
				"def":                 {T: reflect.String},
				"map":                 {T: reflect.Int64},
				"disappearAfterTicks": {T: reflect.Int64},
				"questTags":           {T: createEmptyType()},
				"growth":              {T: reflect.Float64},
				"age":                 {T: reflect.Int64},
				"health":              {T: reflect.Int64},
				"thing": {T: createStructForTest("thing", map[string]*Member{
					"id":                  {T: reflect.String},
					"pos":                 {T: reflect.String},
					"def":                 {T: reflect.String},
					"map":                 {T: reflect.Int64},
					"disappearAfterTicks": {T: reflect.Int64},
					"questTags":           {T: createEmptyType()},
					"growth":              {T: reflect.Float64},
					"age":                 {T: reflect.Int64},
					"health":              {T: reflect.Int64},
				})},
			})),
		},

		"slice_mix float64/int64": {
			args: args{
				xmlContent: `
<?xml version="1.0" encoding="utf-8"?>
<savegame>
	<vals>
		<li>
			<progression>65</progression>
		</li>
		<li>
			<progression>35.3</progression>
		</li>
		<li>
			<progression>23</progression>
		</li>
	</vals>
</savegame>
`,
			},
			want: createCustomSliceForTest(createStructForTest("vals", map[string]*Member{
				"progression": {
					T: reflect.Float64,
				},
			})),
		},

		"slice_only int64": {
			args: args{
				xmlContent: `
<?xml version="1.0" encoding="utf-8"?>
<savegame>
	<vals>
		<li>
			<count>65</count>
		</li>
		<li>
			<count>353</count>
		</li>
		<li>
			<count>23</count>
		</li>
	</vals>
</savegame>
`,
			},
			want: createCustomSliceForTest(createStructForTest("vals", map[string]*Member{
				"count": {
					T: reflect.Int64,
				},
			})),
		},

		"map_simple": {
			args: args{
				xmlContent: `
<?xml version="1.0" encoding="utf-8"?>
<savegame>
	<progress>
		<keys>
			<li>CataphractArmor</li>
			<li>JumpPack</li>
			<li>BrainWiring</li>
		</keys>
		<values>
			<li>0</li>
			<li>0</li>
			<li>0</li>
		</values>
	</progress>
</savegame>
`,
			},
			want: createCustomMapForTest(reflect.String, reflect.Int64),
		},

		"map_max_int64&float64": {
			args: args{
				xmlContent: `
<?xml version="1.0" encoding="utf-8"?>
<savegame>
	<progress>
		<keys>
			<li>CataphractArmor</li>
			<li>JumpPack</li>
			<li>BrainWiring</li>
		</keys>
		<values>
			<li>4</li>
			<li>0</li>
			<li>0.3</li>
		</values>
	</progress>
</savegame>
`,
			},
			want: createCustomMapForTest(reflect.String, reflect.Float64),
		},

		"map_empty": {
			args: args{
				xmlContent: `
<?xml version="1.0" encoding="utf-8"?>
<savegame>
	<progress>
		<keys />
		<values />
	</progress>
</savegame>
`,
			},
			want: createCustomMapForTest(reflect.String, createEmptyType()),
		},

		"map_complex_type": {
			args: args{
				xmlContent: `
<?xml version="1.0" encoding="utf-8"?>
<savegame>
	<progress>
		<keys>
			<li>
				<assignedPawns>
					<li>Thing_Human313</li>
				</assignedPawns>
			</li>
		</keys>
		<values>
			<li>68574</li>
		</values>
	</progress>
</savegame>
`,
			},
			want: createCustomMapForTest(createCustomSliceForTest(reflect.String), reflect.Int64),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			root := resetVarsAndReadBuffer(t, tt.args)
			var res any
			if strings.Contains(name, "map") {
				res = createCustomTypeForMap(root.Child, tt.args.flag)
			} else {
				res = createCustomSlice(root.Child, tt.args.flag)
			}
			require.IsType(t, res, tt.want)
			got := res.(*CustomType)
			wanted := tt.want.(*CustomType)
			require.Equal(t, wanted.Name, got.Name)
			if diff := deep.Equal(wanted.Type1, got.Type1); diff != nil {
				assert.FailNow(t, strings.Join(diff, "\n"))
			}
			if diff := deep.Equal(wanted.Type2, got.Type2); diff != nil {
				assert.FailNow(t, strings.Join(diff, "\n"))
			}
		})
	}
}
