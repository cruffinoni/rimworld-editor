package main

import (
	"github.com/cruffinoni/rimworld-editor/cmd/ui"
	"github.com/cruffinoni/rimworld-editor/editor"
	"github.com/cruffinoni/rimworld-editor/generator"
	"github.com/jawher/mow.cli"
	"log"
	"os"
)

const (
	cliVersion = "0.0.1"
)

type modeType string

const (
	modeConsole = modeType("console")
	modeGUI     = modeType("gui")
)

func isValidMode(mode string) bool {
	return mode == string(modeConsole) || mode == string(modeGUI)
}

// Application is the main application.
type Application struct {
	ui.Options
	*cli.Cli

	fileOpening *editor.FileOpening
}

/*
	TODO:
		- Console mode: saveXML, generateGoFiles
		- GUI mode: edit data
*/

func CreateApplication() *Application {
	app := &Application{}
	app.Cli = cli.App("rimworld-editor", "Rimworld save game editor")
	app.Version("version", cliVersion)
	app.BoolOptPtr(&app.Verbose, "v verbose", false, "Verbose mode")
	app.BoolOptPtr(&app.Generate, "g generate", false, "Generate go files from xml")
	app.StringOptPtr(&app.Output, "o output", "generated", "Output file")
	app.StringArgPtr(&app.Mode, "MODE", string(modeConsole), "The mode to run the application in")
	app.StringArgPtr(&app.Input, "INPUT", "", "Save game file to explore") // TODO: Later use StringOptPtr and discover the file automatically
	app.Before = app.beforeExecution
	app.Action = func() {

	}
	return app
}

func (app *Application) beforeExecution() {
	if !isValidMode(app.Mode) {
		app.PrintHelp()
		cli.Exit(1)
	}
	var err error
	if err = app.ReadFile(); err != nil {
		log.Fatal(err)
	}
	if app.Generate {
		root := generator.GenerateGoFiles(app.fileOpening.XML.Root)
		if err = root.WriteGoFile(app.Output); err != nil {
			log.Fatal(err)
		}
	}
	//if app.Mode == string(modeConsole) {
	//	app.Action = app.consoleMode
	//} else if app.Mode == string(modeGUI) {
	//	app.Action = app.guiMode
	//}
}

func (app *Application) Run() error {
	return app.Cli.Run(os.Args)
}

func (app *Application) RunWithArgs(args []string) error {
	return app.Cli.Run(args)
}

func (app *Application) ReadFile() error {
	var err error
	app.fileOpening, err = editor.Open(app.Input)
	return err
}
