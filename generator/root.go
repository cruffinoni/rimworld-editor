package generator

import (
	"github.com/cruffinoni/rimworld-editor/xml"
)

type MemberVersioning map[string][]*StructInfo

var RegisteredMembers = make(MemberVersioning)

// GenerateGoFiles generates the Go files (with the corresponding structs)
// for the given XML file, but it doesn't write anything.
// To do that, call WriteGoFile.
func GenerateGoFiles(root *xml.Element, withMVFix bool) *StructInfo {
	s := &StructInfo{
		Members: make(map[string]*member),
	}
	//log.Printf("Generating Go files for %s", root.XMLPath())
	RegisteredMembers = make(map[string][]*StructInfo)
	UniqueNumber = 0
	if err := handleElement(root, s, flagNone); err != nil {
		panic(err)
	}
	//printer.PrintSf("{-BOLD}%d{-RESET} members registered. Fixing type mismatch.", len(RegisteredMembers))
	if withMVFix {
		FixRegisteredMembers(RegisteredMembers)
	}
	return s
}

func FixRegisteredMembers(mv MemberVersioning) {
	for i := range mv {
		l := len(mv[i])
		//printer.PrintSf("Fixing %s (%d fix to do)...", i, l)
		for j := 0; j < l; j++ {
			fixMembers(mv[i][0], mv[i][j])
		}
		deleteTitleDuplicate(mv[i][0])
	}
}
