package game

import (
	"github.com/cruffinoni/rimworld-editor/xml"
)

type TickManager struct {
	StartingYear int64 `xml:"startingYear"`
}

func (t *TickManager) SetAttributes(_ xml.Attributes) {
	// No attributes need to be set.
}

func (t *TickManager) GetAttributes() xml.Attributes {
	return nil
}

func (t *TickManager) Assign(_ *xml.Element) error {
	return nil
}

func (t *TickManager) GetPath() string {
	return ""
}
