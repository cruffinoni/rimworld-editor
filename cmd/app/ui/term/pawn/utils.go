package pawn

import (
	"fmt"
	"log"

	"github.com/cruffinoni/rimworld-editor/generated"
)

func getPawnFullNameColorFormatted(p *generated.Thing) string {
	if p.Name == nil {
		return "<nil pawn name>"
	}
	log.Printf("=> %+v", p.Name)
	if p.Name.Nick != "" {
		return fmt.Sprintf("{{{-F_GREEN}}}%s {{{-F_MAGENTA}}}%s {{{-F_CYAN}}}%s{{{-RESET}}}", p.Name.First, p.Name.Nick, p.Name.Last)
	} else {
		return fmt.Sprintf("{{{-F_GREEN}}}%s {{{-F_CYAN}}}%s{{{-RESET}}}", p.Name.First, p.Name.Last)
	}
}
