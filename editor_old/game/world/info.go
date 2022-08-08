package world

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/xml/types"
)

type Info struct {
	Name                  string                   `xml:"name"`
	PlanetCoverage        float64                  `xml:"planetCoverage"`
	PersistentRandomValue int64                    `xml:"persistentRandomValue"`
	OverallRainfall       string                   `xml:"overallRainfall"`
	OverallTemperature    string                   `xml:"overallTemperature"`
	InitialMapSize        string                   `xml:"initialMapSize"`
	FactionCounts         types.Map[string, int64] `xml:"factionCounts"`
}

func (i *Info) Assign(_ *xml.Element) error {
	return nil
}

func (i *Info) GetPath() string {
	return ""
}

func (i *Info) SetAttributes(_ attributes.Attributes) {
	// No attributes need to be set.
}

func (i *Info) GetAttributes() attributes.Attributes {
	return nil
}
