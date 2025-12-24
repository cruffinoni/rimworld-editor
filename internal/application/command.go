package application

import (
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/cruffinoni/rimworld-editor/generated"
	"github.com/cruffinoni/rimworld-editor/internal/application/term"
	"github.com/cruffinoni/rimworld-editor/internal/xml/binder"
)

func (app *Application) newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "rimworld-editor",
		Short:         "Rimworld save game editor",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.beforeExecution(cmd); err != nil {
				return err
			}
			return app.execute()
		},
	}
	cmd.Version = cliVersion
	cmd.SetVersionTemplate("{{.Version}}\n")
	cmd.Flags().BoolVarP(&app.Verbose, "verbose", "v", false, "Verbose mode")
	cmd.Flags().BoolVarP(&app.Generate, "generate", "g", false, "Generate go files from xml")
	cmd.Flags().BoolVarP(&app.Save, "save", "s", true, "Save your modifications when exiting the application")
	cmd.Flags().StringVarP(&app.Output, "output", "o", "generated", "Output folder for generated files")
	cmd.Flags().StringVarP(&app.Mode, "mode", "m", modeConsole, "The mode to run the application in")
	cmd.Flags().IntVarP(&app.MaxSaveGameFileDiscover, "maxnb", "x", 10, "Maximum number of save games to discover")
	cmd.Flags().StringVarP(&app.Input, "defaultsave", "d", "", "Default save game to load from your Rimworld saves game folder")
	cmd.Flags().StringVarP(&app.OperatingSystem, "operating-system", "O", "", "Force a operating system file path finding")
	cmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(cliVersion)
		},
	})
	return cmd
}

func (app *Application) execute() error {
	if app.Mode == modeConsole {
		app.ui = &term.Console{}
	} else if app.Mode == modeGUI {
		return fmt.Errorf("mode %q is not implemented", app.Mode)
	}
	app.ui.SetLogger(app.logger)
	structInit := &generated.GeneratedStructStarter0{}
	app.logger.Debug("Unmarshalling XML")
	s := spinner.New(spinner.CharSets[21], 100*time.Millisecond)
	s.FinalMSG = "XML file unmarshalled successfully\n"
	s.Start()
	if err := binder.Element(app.logger, app.fileOpening.XML.Root, structInit); err != nil {
		return err
	}
	s.Stop()
	structInit.ValidateField("Savegame")
	app.logger.Debug("Initializing UI")
	app.ui.Init(&app.Options, structInit.Savegame)
	app.logger.Debug("Running UI")
	if err := app.ui.Execute(os.Args); err != nil {
		return err
	}
	if app.Save {
		app.logger.Debug("End of execution, generating new file")
		if err := app.SaveGameFile(structInit.Savegame); err != nil {
			return err
		}
	}
	return nil
}
