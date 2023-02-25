package generator

import (
	"log"

	"github.com/cruffinoni/rimworld-editor/xml"
)

var RegisteredMembers = make(map[string][]*StructInfo)

// GenerateGoFiles generates the Go files (with the corresponding structs)
// for the given XML file, but it doesn't write anything.
// To do that, call WriteGoFile.
func GenerateGoFiles(root *xml.Element) *StructInfo {
	s := &StructInfo{
		Members: make(map[string]*member),
	}
	log.Printf("Generating Go files for %s", root.XMLPath())
	RegisteredMembers = make(map[string][]*StructInfo)
	if err := handleElement(root, s, flagNone); err != nil {
		panic(err)
	}
	log.Printf("%d members registered. Fixing type mismatch.", len(RegisteredMembers))
	for i := range RegisteredMembers {
		l := len(RegisteredMembers[i])
		log.Printf("Fixing %s (%d fix to do)...", i, l)
		for j := 0; j < l; j++ {
			fixMembers(RegisteredMembers[i][0], RegisteredMembers[i][j])
		}
		deleteTitleDuplicate(RegisteredMembers[i][0])
	}
	return s
}
