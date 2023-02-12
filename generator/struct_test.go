package generator

import (
	"reflect"
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_createStructure(t *testing.T) {
	tests := map[string]tests{
		"simple": {
			args: args{
				xmlContent: `
<?xml version="1.0" encoding="utf-8"?>
<savegame>
	<quests>
		<completed>False</completed>
	</quests>
</savegame>
`,
			},
			want: createStructForTest("quests", map[string]*member{
				"completed": {
					T: reflect.String,
				},
			}),
		},

		"empty": {
			args: args{
				xmlContent: `
<?xml version="1.0" encoding="utf-8"?>
<savegame>
	<quests>
		<completed />
	</quests>
</savegame>
`,
			},
			want: createStructForTest("quests", map[string]*member{
				"completed": {
					T: createEmptyType(),
				},
			}),
		},

		"disguised structure": {
			args: args{
				xmlContent: `
		<?xml version="1.0" encoding="utf-8"?>
		<savegame>
			<skills>
				<skills>
					<li>
						<def>Shooting</def>
						<level>9</level>
						<passion>Major</passion>
					</li>
					<li>
						<def>Melee</def>
						<level>4</level>
						<passion>Minor</passion>
					</li>
					<li>
						<def>Construction</def>
					</li>
					<li>
						<def>Mining</def>
					</li>
					<li>
						<def>Cooking</def>
						<level>1</level>
					</li>
				</skills>
				<lastXpSinceMidnightResetTimestamp>6127575
				</lastXpSinceMidnightResetTimestamp>
			</skills>
		</savegame>
		`,
			},
			want: createStructForTest("skills", map[string]*member{
				"skills": {
					T: createCustomSliceForTest(createStructForTest("skills_Inner", map[string]*member{
						"def": {
							T: reflect.String,
						},
						"level": {
							T: reflect.Int64,
						},
						"passion": {
							T: reflect.String,
						},
					})),
				},
				"lastXpSinceMidnightResetTimestamp": {
					T: reflect.Int64,
				},
			}),
		},

		"with attr": {
			args: args{
				xmlContent: `
<?xml version="1.0" encoding="utf-8"?>
<savegame>
	<quests>
		<completed Class="Need_Mood" />
	</quests>
</savegame>
`,
			},
			want: createStructForTest("quests", map[string]*member{
				"completed": {
					T: createEmptyType(),
					Attr: map[string]string{
						"Class": "Need_Mood",
					},
				},
			}),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			root := resetVarsAndReadBuffer(t, tt.args)
			res := createStructure(root.Child, tt.args.flag)
			require.IsType(t, res, tt.want)
			got := res.(*StructInfo)
			wanted := tt.want.(*StructInfo)
			assert.Equal(t, wanted, got)
		})
	}
}

func TestGenerateGoFiles(t *testing.T) {
	tests := map[string]tests{
		"large cover (slice, array, empty)": {
			args: args{
				xmlContent: `
<?xml version="1.0" encoding="utf-8"?>
<savegame>
	<type>
		<li>
			<createdFromNoExpansionGame>True</createdFromNoExpansionGame>
			<foundation />
			<name>Rustican 2</name>
			<culture>Rustican</culture>
			<memes/>
			<precepts>
				<li>
					<name>cannibalisme</name>
					<def>Cannibalism_Classic</def>
					<ID>12</ID>
					<randomSeed>-2137262173</randomSeed>
					<usesDefiniteArticle>True</usesDefiniteArticle>
				</li>
				<li>
					<name>cadavres</name>
					<def>Corpses_Ugly</def>
					<ID>13</ID>
					<randomSeed>950454885</randomSeed>
					<usesDefiniteArticle>True</usesDefiniteArticle>
				</li>
				<li>
					<name>exécution</name>
					<def>Execution_Classic</def>
					<ID>18</ID>
					<randomSeed>502274313</randomSeed>
					<usesDefiniteArticle>True</usesDefiniteArticle>
				</li>
				<li>
					<name>viande d'insecte</name>
					<def>InsectMeatEating_Despised_Classic</def>
					<ID>14</ID>
					<randomSeed>-837163784</randomSeed>
					<usesDefiniteArticle>True</usesDefiniteArticle>
				</li>
				<li>
					<name>nom marital</name>
					<def>MarriageName_UsuallyMans</def>
					<ID>16</ID>
					<randomSeed>-855175996</randomSeed>
					<usesDefiniteArticle>True</usesDefiniteArticle>
				</li>
				<li>
					<name>manger la pâte nutritive</name>
					<def>NutrientPasteEating_Disgusting</def>
					<ID>11</ID>
					<randomSeed>-1141618061</randomSeed>
					<usesDefiniteArticle>True</usesDefiniteArticle>
				</li>
				<li>
					<name>utilisation des organes</name>
					<def>OrganUse_Classic</def>
					<ID>17</ID>
					<randomSeed>-1204706330</randomSeed>
					<usesDefiniteArticle>True</usesDefiniteArticle>
				</li>
				<li>
					<name>traite des esclaves</name>
					<def>Slavery_Classic</def>
					<ID>19</ID>
					<randomSeed>-301498631</randomSeed>
					<usesDefiniteArticle>True</usesDefiniteArticle>
				</li>
				<li>
					<name>amour physique</name>
					<def>Lovin_Free</def>
					<ID>15</ID>
					<randomSeed>-196120713</randomSeed>
					<usesDefiniteArticle>True</usesDefiniteArticle>
				</li>
				<li>
					<name>Époux</name>
					<def>SpouseCount_Female_MaxOne</def>
					<ID>21</ID>
					<randomSeed>-505669973</randomSeed>
					<usesDefiniteArticle>True</usesDefiniteArticle>
				</li>
				<li>
					<name>Épouses</name>
					<def>SpouseCount_Male_MaxOne</def>
					<ID>20</ID>
					<randomSeed>1950680083</randomSeed>
					<usesDefiniteArticle>True</usesDefiniteArticle>
				</li>
			</precepts>
			<thingStyleCategories/>
			<usedSymbols/>
			<usedSymbolPacks/>
			<style>
				<hairFrequencies>
					<vals>
						<li>
							<frequency>Normal</frequency>
						</li>
						<li>
							<frequency>Normal</frequency>
							<gender>Female</gender>
						</li>
						<li/>
						<li/>
						<li/>
						<li>
							<gender>Female</gender>
						</li>
						<li>
							<gender>Any</gender>
						</li>
						<li>
							<gender>Female</gender>
						</li>
						<li/>
						<li/>
						<li/>
						<li>
							<gender>Female</gender>
						</li>
						<li>
							<gender>Female</gender>
						</li>
						<li>
							<gender>Female</gender>
						</li>
						<li>
							<gender>Female</gender>
						</li>
						<li>
							<frequency>Normal</frequency>
							<gender>Female</gender>
						</li>
						<li>
							<frequency>Normal</frequency>
							<gender>Female</gender>
						</li>
						<li>
							<frequency>Normal</frequency>
							<gender>Female</gender>
						</li>
						<li>
							<frequency>Normal</frequency>
							<gender>MaleUsually</gender>
						</li>
						<li>
							<frequency>Normal</frequency>
							<gender>Any</gender>
						</li>
						<li>
							<frequency>Normal</frequency>
						</li>
						<li>
							<frequency>Normal</frequency>
							<gender>Female</gender>
						</li>
						<li>
							<frequency>Normal</frequency>
							<gender>Female</gender>
						</li>
						<li>
							<gender>Female</gender>
						</li>
						<li>
							<frequency>Normal</frequency>
							<gender>Female</gender>
						</li>
						<li>
							<frequency>Normal</frequency>
						</li>
						<li>
							<frequency>Normal</frequency>
						</li>
						<li/>
						<li>
							<frequency>Normal</frequency>
						</li>
						<li>
							<frequency>Normal</frequency>
						</li>
						<li>
							<gender>MaleUsually</gender>
						</li>
						<li>
							<gender>FemaleUsually</gender>
						</li>
						<li>
							<gender>FemaleUsually</gender>
						</li>
						<li>
							<gender>MaleUsually</gender>
						</li>
						<li>
							<gender>MaleUsually</gender>
						</li>
						<li>
							<gender>MaleUsually</gender>
						</li>
						<li>
							<gender>Any</gender>
						</li>
						<li>
							<gender>Any</gender>
						</li>
						<li>
							<gender>FemaleUsually</gender>
						</li>
						<li>
							<gender>FemaleUsually</gender>
						</li>
						<li>
							<gender>FemaleUsually</gender>
						</li>
						<li/>
						<li>
							<gender>MaleUsually</gender>
						</li>
						<li>
							<gender>MaleUsually</gender>
						</li>
						<li>
							<gender>Any</gender>
						</li>
						<li>
							<frequency>Normal</frequency>
							<gender>Any</gender>
						</li>
						<li>
							<gender>MaleUsually</gender>
						</li>
						<li>
							<gender>Any</gender>
						</li>
						<li/>
						<li>
							<gender>MaleUsually</gender>
						</li>
					</vals>
				</hairFrequencies>
				<styleForThingDef>
					<keys/>
					<values/>
				</styleForThingDef>
			</style>
			<id>1</id>
			<development />
		</li>

		<li>
			<createdFromNoExpansionGame>True</createdFromNoExpansionGame>
			<foundation />
			<name>Corunan</name>
			<culture>Corunan</culture>
			<memes/>
			<precepts>
				<li>
					<name>cannibalisme</name>
					<def>Cannibalism_Classic</def>
					<ID>23</ID>
					<randomSeed>-606671294</randomSeed>
					<usesDefiniteArticle>True</usesDefiniteArticle>
				</li>
				<li>
					<name>cadavres</name>
					<def>Corpses_Ugly</def>
					<ID>24</ID>
					<randomSeed>1865187567</randomSeed>
					<usesDefiniteArticle>True</usesDefiniteArticle>
				</li>
				<li>
					<name>exécution</name>
					<def>Execution_Classic</def>
					<ID>29</ID>
					<randomSeed>-422599061</randomSeed>
					<usesDefiniteArticle>True</usesDefiniteArticle>
				</li>
				<li>
					<name>viande d'insecte</name>
					<def>InsectMeatEating_Despised_Classic</def>
					<ID>25</ID>
					<randomSeed>968103858</randomSeed>
					<usesDefiniteArticle>True</usesDefiniteArticle>
				</li>
				<li>
					<name>nom marital</name>
					<def>MarriageName_UsuallyMans</def>
					<ID>27</ID>
					<randomSeed>1553466286</randomSeed>
					<usesDefiniteArticle>True</usesDefiniteArticle>
				</li>
				<li>
					<name>manger la pâte nutritive</name>
					<def>NutrientPasteEating_Disgusting</def>
					<ID>22</ID>
					<randomSeed>-1696315431</randomSeed>
					<usesDefiniteArticle>True</usesDefiniteArticle>
				</li>
				<li>
					<name>utilisation des organes</name>
					<def>OrganUse_Classic</def>
					<ID>28</ID>
					<randomSeed>2028129588</randomSeed>
					<usesDefiniteArticle>True</usesDefiniteArticle>
				</li>
				<li>
					<name>traite des esclaves</name>
					<def>Slavery_Classic</def>
					<ID>30</ID>
					<randomSeed>-470866308</randomSeed>
					<usesDefiniteArticle>True</usesDefiniteArticle>
				</li>
				<li>
					<name>amour physique</name>
					<def>Lovin_Free</def>
					<ID>26</ID>
					<randomSeed>-1442272627</randomSeed>
					<usesDefiniteArticle>True</usesDefiniteArticle>
				</li>
				<li>
					<name>Époux</name>
					<def>SpouseCount_Female_MaxOne</def>
					<ID>32</ID>
					<randomSeed>-954237569</randomSeed>
					<usesDefiniteArticle>True</usesDefiniteArticle>
				</li>
				<li>
					<name>Épouses</name>
					<def>SpouseCount_Male_MaxOne</def>
					<ID>31</ID>
					<randomSeed>246528037</randomSeed>
					<usesDefiniteArticle>True</usesDefiniteArticle>
				</li>
			</precepts>
			<thingStyleCategories/>
			<usedSymbols/>
			<usedSymbolPacks/>
			<style>
				<hairFrequencies>
					<vals>
						<li/>
						<li>
							<gender>Female</gender>
						</li>
						<li/>
						<li/>
						<li/>
						<li>
							<gender>Female</gender>
						</li>
						<li>
							<gender>Any</gender>
						</li>
						<li>
							<gender>Female</gender>
						</li>
						<li/>
						<li/>
						<li/>
						<li>
							<gender>Female</gender>
						</li>
						<li>
							<gender>Female</gender>
						</li>
						<li>
							<gender>Female</gender>
						</li>
						<li>
							<gender>Female</gender>
						</li>
						<li>
							<gender>Female</gender>
						</li>
						<li>
							<gender>Female</gender>
						</li>
						<li>
							<gender>Female</gender>
						</li>
						<li>
							<gender>MaleUsually</gender>
						</li>
						<li>
							<gender>Any</gender>
						</li>
						<li/>
						<li>
							<gender>Female</gender>
						</li>
						<li>
							<gender>Female</gender>
						</li>
						<li>
							<gender>Female</gender>
						</li>
						<li>
							<gender>Female</gender>
						</li>
						<li/>
						<li/>
						<li/>
						<li/>
						<li/>
						<li>
							<gender>MaleUsually</gender>
						</li>
						<li>
							<gender>FemaleUsually</gender>
						</li>
						<li>
							<gender>FemaleUsually</gender>
						</li>
						<li>
							<gender>MaleUsually</gender>
						</li>
						<li>
							<gender>MaleUsually</gender>
						</li>
						<li>
							<gender>MaleUsually</gender>
						</li>
						<li>
							<frequency>Normal</frequency>
							<gender>Any</gender>
						</li>
						<li>
							<frequency>Normal</frequency>
							<gender>Any</gender>
						</li>
						<li>
							<frequency>Normal</frequency>
							<gender>FemaleUsually</gender>
						</li>
						<li>
							<frequency>Normal</frequency>
							<gender>FemaleUsually</gender>
						</li>
						<li>
							<frequency>Normal</frequency>
							<gender>FemaleUsually</gender>
						</li>
						<li>
							<frequency>Normal</frequency>
						</li>
						<li>
							<frequency>Normal</frequency>
							<gender>MaleUsually</gender>
						</li>
						<li>
							<frequency>Normal</frequency>
							<gender>MaleUsually</gender>
						</li>
						<li>
							<frequency>Normal</frequency>
							<gender>Any</gender>
						</li>
						<li>
							<frequency>Normal</frequency>
							<gender>Any</gender>
						</li>
						<li>
							<frequency>Normal</frequency>
							<gender>MaleUsually</gender>
						</li>
						<li>
							<frequency>Normal</frequency>
							<gender>Any</gender>
						</li>
						<li/>
						<li>
							<gender>MaleUsually</gender>
						</li>
					</vals>
				</hairFrequencies>
				<styleForThingDef>
					<keys/>
					<values/>
				</styleForThingDef>
			</style>
			<id>2</id>
			<development />
		</li>
	</type>
</savegame>
`,
			},
			want: createStructForTest("savegame", map[string]*member{
				"type": {
					T: createCustomSliceForTest(createStructForTest("type", map[string]*member{
						"foundation":                 {T: createEmptyType()},
						"id":                         {T: reflect.Int64},
						"culture":                    {T: reflect.String},
						"usedSymbols":                {T: createEmptyType()},
						"usedSymbolPacks":            {T: createEmptyType()},
						"development":                {T: createEmptyType()},
						"createdFromNoExpansionGame": {T: reflect.String},
						"name":                       {T: reflect.String},
						"memes":                      {T: createEmptyType()},
						"precepts": {
							T: createCustomSliceForTest(createStructForTest("precepts", map[string]*member{
								"name":                {T: reflect.String},
								"def":                 {T: reflect.String},
								"ID":                  {T: reflect.Int64},
								"randomSeed":          {T: reflect.Int64},
								"usesDefiniteArticle": {T: reflect.String},
							})),
						},
						"thingStyleCategories": {T: createEmptyType()},
						"style": {
							T: createStructForTest("style", map[string]*member{
								"hairFrequencies": {
									T: createStructForTest("hairFrequencies", map[string]*member{
										"vals": {
											T: createFixedArrayForTest(50, createStructForTest("vals", map[string]*member{
												"frequency": {T: reflect.String},
												"gender":    {T: reflect.String},
											})),
										},
									}),
								},
								"styleForThingDef": {T: createCustomMapForTest(reflect.String, createEmptyType())},
							}),
						},
					})),
				},
			}),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			root := resetVarsAndReadBuffer(t, tt.args)
			if diff := deep.Equal(tt.want, GenerateGoFiles(root)); diff != nil {
				assert.FailNow(t, strings.Join(diff, "\n"))
			}
		})
	}
}
