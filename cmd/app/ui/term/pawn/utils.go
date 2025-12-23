package pawn

import (
	"fmt"

	"github.com/cruffinoni/rimworld-editor/generated"
)

func getPawnFullName(p *generated.Thing) string {
	if p.Name == nil {
		return "<nil pawn name>"
	}
	if p.Name.Nick != "" {
		return fmt.Sprintf("%s \"%s\" %s", p.Name.First, p.Name.Nick, p.Name.Last)
	}
	return fmt.Sprintf("%s %s", p.Name.First, p.Name.Last)
}
