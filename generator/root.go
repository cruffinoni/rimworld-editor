package generator

import (
	"log"

	"github.com/cruffinoni/rimworld-editor/cmd/app/ui/term/printer"
	"github.com/cruffinoni/rimworld-editor/xml"
)

type MemberVersioning map[string][]*StructInfo

var RegisteredMembers = make(MemberVersioning)

func cleanUpMVPtrs(mv MemberVersioning) {
	for i := range mv {
		l := len(mv[i])
		uniquePtr := make(map[*StructInfo]bool)
		if l >= 1 {
			for j := 0; j < l; j++ {
				if _, ok := uniquePtr[mv[i][j]]; !ok {
					uniquePtr[mv[i][j]] = true
				}
			}
			mv[i] = make([]*StructInfo, len(uniquePtr))
			j := 0
			for k := range uniquePtr {
				mv[i][j] = k
				j++
			}
		}
	}
}

// GenerateGoFiles generates the Go files (with the corresponding structs)
// for the given XML file, but it doesn't write anything.
// To do that, call WriteGoFile.
func GenerateGoFiles(root *xml.Element, withMVFix bool) *StructInfo {
	s := &StructInfo{
		Members: make(map[string]*Member),
	}
	//log.Printf("Generating Go files for %s", root.XMLPath())
	RegisteredMembers = make(MemberVersioning)
	UniqueNumber = 0
	if err := handleElement(root, s, flagNone); err != nil {
		panic(err)
	}
	if withMVFix {
		printer.Print("Cleaning up the MemberVersioning pointers")
		cleanUpMVPtrs(RegisteredMembers)
		printer.Printf("{-BOLD}%d{-RESET} members registered. Fixing type mismatch.", len(RegisteredMembers))
		FixRegisteredMembers(RegisteredMembers)
	}
	return s
}

func FixRegisteredMembers(mv MemberVersioning) {
	for i := range mv {
		l := len(mv[i])
		if l >= 1 {
			printer.Printf("Fixing %s ({-BOLD,F_RED}%d{-RESET} fix to do)...", i, l)
			for j := 1; j < l; j++ {
				//log.Printf("Name: %v (0) & %v (%d)", mv[i][0].Name, mv[i][j].Name, j)
				if mv[i][0] == mv[i][j] {
					log.Fatalf("Identical pointers: %p & %p / Probable infinite recursion - %v", mv[i][0], mv[i][j], i)
				}
				FixMembers(mv[i][0], mv[i][j])
				//log.Printf("Done")
			}
		}
		deleteDuplicateTitle(mv[i][0])
	}
}
