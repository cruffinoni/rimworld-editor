package game

import (
	"github.com/cruffinoni/rimworld-editor/xml"
)

type TickManager struct {
	StartingYear int32 `xml:"startingYear"`
}

func (t *TickManager) Assign(_ *xml.Element) error {
	return nil
}

func (t *TickManager) GetPath() string {
	return ""
}
