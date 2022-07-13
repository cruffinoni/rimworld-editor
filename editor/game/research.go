package game

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/type"
)

type ResearchManager struct {
	Progress _type.Map[string, float64] `xml:"progress"`
}

func (r *ResearchManager) Assign(_ *xml.Element) error {
	return nil
}

func (r *ResearchManager) GetPath() string {
	return ""
}

func (r *ResearchManager) SetAttributes(_ xml.Attributes) {
	// No attributes need to be set.
}

func (r *ResearchManager) GetAttributes() xml.Attributes {
	return nil
}

func (r *ResearchManager) GetProgress(key string) float64 {
	return r.Progress.Get(key)
}

func (r *ResearchManager) SetProgress(key string, value float64) {
	r.Progress.Set(key, value)
}

const ResearchMaxValue = 10000

func (r *ResearchManager) SetAllProgress(value float64) {
	for it := r.Progress.Iterator(); it != nil; it = it.Next() {
		r.Progress.Set(it.Key(), value)
	}
}
