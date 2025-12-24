package faction

import (
	"github.com/spf13/cobra"

	"github.com/cruffinoni/rimworld-editor/generated"
	"github.com/cruffinoni/rimworld-editor/pkg/logging"
)

func NewCommand(logger logging.Logger, sg *generated.Savegame, rf Registerer) *cobra.Command {
	fl := NewList(logger, sg, rf)
	cmd := &cobra.Command{
		Use:   "faction",
		Short: "Faction commands",
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all factions",
		Run: func(cmd *cobra.Command, args []string) {
			fl.ListAllFactions()
		},
	})
	return cmd
}
