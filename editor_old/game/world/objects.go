package world

import (
	"github.com/cruffinoni/rimworld-editor/xml/types"
	"github.com/cruffinoni/rimworld-editor/xml/types/primary"
)

type Trader struct {
	TmpSavedPawns            string                     `xml:"tmpSavedPawns"`
	Stock                    types.Slice[primary.Empty] `xml:"stock"`
	LastStockGenerationTicks int64                      `xml:"lastStockGenerationTicks"`
}

type Object struct {
	Def                            string        `xml:"def"`
	ID                             string        `xml:"ID"`
	Tile                           string        `xml:"tile"`
	Faction                        string        `xml:"faction"`
	QuestTags                      primary.Empty `xml:"questTags"`
	Expiration                     string        `xml:"expiration"`
	PreviouslyGeneratedInhabitants string        `xml:"previouslyGeneratedInhabitants"`
	Trader                         *Trader       `xml:"trader"`
	NameInt                        string        `xml:"nameInt"`
}

type Objects struct {
	Objects types.Slice[*Object] `xml:"object"`
}
