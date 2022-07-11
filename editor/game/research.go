package game

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/type"
	"log"
)

type ResearchManager struct {
	Progress _type.Map[string, int32] `xml:"progress"`
}

func (r *ResearchManager) Assign(_ *xml.Element) error {
	return nil
}

func (r *ResearchManager) GetPath() string {
	return ""
}

func (r *ResearchManager) GetProgress(key string) int32 {
	return r.Progress.Get(key)
}

func (r *ResearchManager) SetProgress(key string, value int32) {
	r.Progress.Set(key, value)
}

func (r *ResearchManager) SetAllProgress(value int32) {
	for it := r.Progress.Iterator(); it != nil; it = it.Next() {
		log.Printf("SetAllProgress: %s = %d", it.Key(), value)
		r.Progress.Set(it.Key(), value)
	}
}
