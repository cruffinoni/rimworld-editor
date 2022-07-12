package world

import "github.com/cruffinoni/rimworld-editor/xml"

type FactionRelation struct {
	Other    string `xml:"other"`
	Kind     string `xml:"kind"`
	Goodwill int64  `xml:"goodwill"`
}

func (f *FactionRelation) Assign(e *xml.Element) error {
	return nil
}

func (f *FactionRelation) GetPath() string {
	return "li"
}

type Faction struct {
	Leader            string             `xml:"leader"`
	Def               string             `xml:"def"`
	Name              string             `xml:"name"`
	LoadID            string             `xml:"loadID"`
	RandomKey         int64              `xml:"randomKey"`
	ColorFromSpectrum float64            `xml:"colorFromSpectrum"`
	CentralMelanin    float64            `xml:"centralMelanin"`
	Relations         []*FactionRelation `xml:"relations"`
	Ideos             *xml.Element       `xml:"ideos"`
	Kidnapped         *xml.Element       `xml:"kidnapped"`
}

func (f *Faction) Assign(e *xml.Element) error {
	return nil
}

func (f *Faction) GetPath() string {
	return "li"
}

type FactionManager struct {
	AllFactions []*Faction `xml:"allFactions"`
}

func (f *FactionManager) Assign(e *xml.Element) error {
	return nil
}

func (f *FactionManager) GetPath() string {
	return ""
}
