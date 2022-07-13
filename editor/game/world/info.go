package world

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	_type "github.com/cruffinoni/rimworld-editor/xml/type"
)

type Info struct {
	Name                  string                   `xml:"name"`
	PlanetCoverage        float64                  `xml:"planetCoverage"`
	PersistentRandomValue int64                    `xml:"persistentRandomValue"`
	OverallRainfall       string                   `xml:"overallRainfall"`
	OverallTemperature    string                   `xml:"overallTemperature"`
	InitialMapSize        string                   `xml:"initialMapSize"`
	FactionCounts         _type.Map[string, int64] `xml:"factionCounts"`
}

func (i *Info) Assign(_ *xml.Element) error {
	return nil
}

func (i *Info) GetPath() string {
	return ""
}

func (i *Info) SetAttributes(_ xml.Attributes) {
	// No attributes need to be set.
}

func (i *Info) GetAttributes() xml.Attributes {
	return nil
}
