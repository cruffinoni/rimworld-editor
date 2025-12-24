package pawn

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cruffinoni/rimworld-editor/generated"
	"github.com/cruffinoni/rimworld-editor/internal/application/term/faction"
	"github.com/cruffinoni/rimworld-editor/pkg/logging"
)

func NewCommand(logger logging.Logger, savegame *generated.Savegame, rp PawnsRegisterer, rf faction.Registerer) *cobra.Command {
	pl := NewList(logger, savegame, rp, rf)
	pi := NewInjury(logger, rp)
	ps := NewSkills(logger, savegame, rp, rf)

	pawnCmd := &cobra.Command{
		Use:   "pawn",
		Short: "Pawn commands",
	}

	injuryCmd := &cobra.Command{
		Use:     "injury",
		Short:   "Commands to manipulate pawn's injury",
		Aliases: []string{"i"},
	}
	injuryCmd.AddCommand(&cobra.Command{
		Use:   "remove-all <PAWN_ID>",
		Short: "Remove all injuries from a pawn",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			pi.RemoveInjuries(args[0])
		},
	})
	injuryCmd.AddCommand(&cobra.Command{
		Use:   "list <PAWN_ID>",
		Short: "List all injuries from a pawn",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			pi.List(args[0])
		},
	})
	injuryCmd.AddCommand(&cobra.Command{
		Use:   "remove <PAWN_ID> <INJURY_IDS...>",
		Short: "Remove a specific injuries from a pawn",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			pi.Remove(args[0], args[1:])
		},
	})
	pawnCmd.AddCommand(injuryCmd)

	listCmd := &cobra.Command{
		Use:     "list",
		Short:   "Commands to lists pawn's related data",
		Aliases: []string{"l"},
	}
	listCmd.AddCommand(&cobra.Command{
		Use:     "world",
		Short:   "List all pawns that alive in the game including your pawns and the faction leaders",
		Aliases: []string{"w"},
		Run: func(cmd *cobra.Command, args []string) {
			pl.ListAllPawns()
		},
	})
	pawnCmd.AddCommand(listCmd)

	skillCmd := &cobra.Command{
		Use:     "skill",
		Short:   "Commands to manage pawns skills",
		Aliases: []string{"s"},
	}
	skillCmd.AddCommand(&cobra.Command{
		Use:   "list [PAWN_ID]",
		Short: "List all pawns skills or a specific pawn if PAWN_ID is set",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			pawnID := "ALL"
			if len(args) > 0 {
				pawnID = args[0]
			}
			ps.ListPawnsSkills(pawnID)
		},
	})
	skillCmd.AddCommand(&cobra.Command{
		Use:   "set <PAWN_ID> <SKILL> <LEVEL> [PASSION]",
		Short: "Set a pawn's skill",
		Args:  cobra.RangeArgs(3, 4),
		RunE: func(cmd *cobra.Command, args []string) error {
			level, err := strconv.Atoi(args[2])
			if err != nil {
				return fmt.Errorf("invalid level: %w", err)
			}
			passion := PassionNone
			if len(args) > 3 {
				passion = args[3]
			}
			ps.Edit(args[0], args[1], passion, level)
			return nil
		},
	})
	skillCmd.AddCommand(&cobra.Command{
		Use:     "forceGraduate <PAWN_ID>",
		Short:   "Force all skills of a pawn to be graduate from 0 - nb of skills",
		Aliases: []string{"fg"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ps.ForceGraduate(args[0])
		},
	})
	pawnCmd.AddCommand(skillCmd)

	return pawnCmd
}
