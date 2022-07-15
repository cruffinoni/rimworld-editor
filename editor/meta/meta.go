package meta

import (
	"errors"
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/xml/path"
)

type Mod struct {
	SteamId int64
	Id      string
}

type Meta struct {
	GameVersion string `xml:"gameVersion"`
	Mods        map[string]*Mod
}

func (m *Meta) SetAttributes(_ attributes.Attributes) {
	// No attributes need to be set.
}

func (m *Meta) GetAttributes() attributes.Attributes {
	return nil
}

func (m *Meta) Assign(e *xml.Element) error {
	elems := []path.Elements{
		path.FindWithPath("modIds>li[...]", e),
		path.FindWithPath("modSteamIds>li[...]", e),
		path.FindWithPath("modNames>li[...]", e),
	}
	for _, elem := range elems {
		if len(elem) == 0 {
			return errors.New("meta: no mod found")
		}
	}
	m.Mods = make(map[string]*Mod)
	for i, elem := range elems[2] {
		m.Mods[elem.Data.GetString()] = &Mod{
			SteamId: elems[1][i].Data.GetInt64(),
			Id:      elems[0][i].Data.GetString(),
		}
	}
	return nil
}

func (m *Meta) GetPath() string {
	return ""
}

func (m *Meta) GetGameVersion() string {
	return m.GameVersion
}

func (m *Meta) GetMods() map[string]*Mod {
	return m.Mods
}

func (m *Meta) GetMod(id string) *Mod {
	if mod, ok := m.Mods[id]; ok {
		return mod
	}
	return nil
}
