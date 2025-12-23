package pawn

import (
	cli "github.com/jawher/mow.cli"

	"github.com/cruffinoni/rimworld-editor/cmd/app/ui/term/faction"
	"github.com/cruffinoni/rimworld-editor/generated"
	"github.com/cruffinoni/rimworld-editor/pkg/logging"
)

func RegisterPawnCommands(logger logging.Logger, cl *cli.Cli, rp PawnsRegisterer, rf faction.Registerer, savegame *generated.Savegame) {
	const pawnParamDescription = "The pawn ID (e.g. Thing_Human119753)"
	pl := NewList(logger, savegame, rp, rf)
	pi := NewInjury(logger, rp)
	cl.Command("pawn", "Pawn commands", func(cmd *cli.Cmd) {
		cmd.Command("injury i", "Commands to manipulate pawn's injury", func(cmd *cli.Cmd) {
			cmd.Command("remove-all", "Remove all injuries from a pawn", func(cmd *cli.Cmd) {
				pawnID := cmd.StringArg("PAWN_ID", "", pawnParamDescription)
				cmd.Action = func() {
					pi.RemoveInjuries(*pawnID)
				}
			})

			cmd.Command("list l", "List all injuries from a pawn", func(cmd *cli.Cmd) {
				pawnID := cmd.StringArg("PAWN_ID", "", pawnParamDescription)
				cmd.Action = func() {
					pi.List(*pawnID)
				}
			})

			cmd.Command("remove r", "Remove a specific injuries from a pawn", func(cmd *cli.Cmd) {
				cmd.Spec = "PAWN_ID INJURY_IDS..."
				pawnID := cmd.StringArg("PAWN_ID", "", pawnParamDescription)
				injuries := cmd.StringsArg("INJURY_IDS", nil, "The load ID of each injury to remove")
				cmd.Action = func() {
					pi.Remove(*pawnID, *injuries)
				}
			})
		})
		cmd.Command("list l", "Commands to lists pawn's related data", func(cmd *cli.Cmd) {
			cmd.Command("world w", "List all pawns that alive in the game including your pawns and the faction leaders", cli.ActionCommand(pl.ListAllPawns))
		})

		ps := NewSkills(logger, savegame, rp, rf)
		cmd.Command("skill s", "Commands to manage pawns skills", func(cmd *cli.Cmd) {
			cmd.Command("list l", "List all pawns skills or a specific pawn if PAWN_ID is set", func(cmd *cli.Cmd) {
				cmd.Spec = "[PAWN_ID]"
				pawnID := cmd.StringArg("PAWN_ID", "ALL", pawnParamDescription)
				cmd.Action = func() {
					ps.ListPawnsSkills(*pawnID)
				}
			})

			cmd.Command("set", "Set a pawn's skill", func(cmd *cli.Cmd) {
				cmd.Spec = "PAWN_ID SKILL LEVEL [PASSION]"
				pawnID := cmd.StringArg("PAWN_ID", "", pawnParamDescription)
				skill := cmd.StringArg("SKILL", "", "The skill to set")
				level := cmd.IntArg("LEVEL", 0, "The level of the skill to set")
				passion := cmd.StringArg("PASSION", PassionNone, "The passion to set, must be either 'Major' or 'Minor'")
				cmd.Action = func() {
					ps.Edit(*pawnID, *skill, *passion, *level)
				}
			})

			cmd.Command("forceGraduate fg", "Force all skills of a pawn to be graduate from 0 - nb of skills", func(cmd *cli.Cmd) {
				pawnID := cmd.StringArg("PAWN_ID", "", pawnParamDescription)
				cmd.Action = func() {
					ps.ForceGraduate(*pawnID)
				}
			})
		})
	})
}
