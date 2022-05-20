package meta

type Meta struct {
	GameVersion string   `xml:"gameVersion"`
	ModsIds     []string `xml:"modIds>li"`
	ModsSteamId []int    `xml:"modSteamIds>li"`
	ModsName    []string `xml:"modNames>li"`
}

func (m *Meta) GetGameVersion() string {
	return m.GameVersion
}

//
//func (m *Meta) GetModByName(name string) *Mod {
//	for i, mod := range m.ModsName.Names {
//		if mod == name {
//			return &Mod{
//				Id:      m.ModsIds.Ids[i],
//				Name:    name,
//				SteamId: m.ModsSteamId.ToString(i),
//			}
//		}
//	}
//	return nil
//}
//
//func (m *Meta) GetModById(id string) *Mod {
//	for i, mod := range m.ModsIds.Ids {
//		if mod == id {
//			return &Mod{
//				Id:      id,
//				Name:    m.ModsName.Names[i],
//				SteamId: m.ModsSteamId.ToString(i),
//			}
//		}
//	}
//	return nil
//}
//
//func (m *Meta) GetModBySteamId(id int) *Mod {
//	for i, mod := range m.ModsSteamId.SteamIds {
//		if mod == id {
//			return &Mod{
//				Id:      m.ModsIds.Ids[i],
//				Name:    m.ModsName.Names[i],
//				SteamId: m.ModsSteamId.ToString(mod),
//			}
//		}
//	}
//	return nil
//}
