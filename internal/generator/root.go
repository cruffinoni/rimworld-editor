package generator

import (
	"github.com/cruffinoni/rimworld-editor/internal/xml"
	"github.com/cruffinoni/rimworld-editor/pkg/logging"
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
func GenerateGoFiles(logger logging.Logger, root *xml.Element, withMVFix bool) *StructInfo {
	s := &StructInfo{
		Members: make(map[string]*Member),
	}
	//printer.Debugf("Generating Go files for %s", root.XMLPath())
	RegisteredMembers = make(MemberVersioning)
	UniqueNumber = 0
	if err := handleElement(logger, root, s, flagNone); err != nil {
		panic(err)
	}
	if withMVFix {
		logger.Debug("Cleaning up MemberVersioning pointers")
		cleanUpMVPtrs(RegisteredMembers)
		logger.WithField("members", len(RegisteredMembers)).Debug("Fixing type mismatch for registered members")
		FixRegisteredMembers(logger, RegisteredMembers)
	}
	return s
}

func FixRegisteredMembers(logger logging.Logger, mv MemberVersioning) {
	for i := range mv {
		l := len(mv[i])
		if l >= 1 {
			logger.WithFields(logging.Fields{
				"struct": i,
				"count":  l,
			}).Debug("Fixing registered members")
			for j := 1; j < l; j++ {
				//printer.Debugf("Name: %v (0) & %v (%d)", mv[i][0].Name, mv[i][j].Name, j)
				if mv[i][0] == mv[i][j] {
					logger.WithFields(logging.Fields{
						"struct":     i,
						"pointer_a":  mv[i][0],
						"pointer_b":  mv[i][j],
						"suspicious": true,
					}).Debug("Identical pointers detected")
					continue
				}
				FixMembers(logger, mv[i][0], mv[i][j])
				//printer.Debugf("Done")
			}
		}
		deleteDuplicateTitle(mv[i][0])
	}
}
