package world

import (
	"github.com/cruffinoni/rimworld-editor/xml/types"
	"github.com/cruffinoni/rimworld-editor/xml/types/primary"
)

type Object struct {
	Def                            string        `xml:"def"`
	ID                             string        `xml:"ID"`
	Tile                           string        `xml:"tile"`
	Faction                        string        `xml:"faction"`
	QuestTags                      primary.Empty `xml:"questTags"`
	Expiration                     string        `xml:"expiration"`
	PreviouslyGeneratedInhabitants string        `xml:"previouslyGeneratedInhabitants"`
	Trader                         struct {
		Text          string `xml:",chardata"`
		TmpSavedPawns string `xml:"tmpSavedPawns"`
		Stock         struct {
			Text   string `xml:",chardata"`
			IsNull string `xml:"IsNull,attr"`
		} `xml:"stock"`
		LastStockGenerationTicks string `xml:"lastStockGenerationTicks"`
	} `xml:"trader"`
	NameInt string `xml:"nameInt"`
}

type Objects struct {
	Objects types.Slice[*Object] `xml:"object"`
}
