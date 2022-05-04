package main

import (
	"github.com/cruffinoni/rimworld-editor/editor"
	"github.com/cruffinoni/rimworld-editor/xml/path"
	"log"
)

func main() {
	println("Hello, world.")
	if f, err := editor.Open("TEST.rws"); err != nil {
		panic(err)
	} else {
		_ = f

		p := path.NewPathing("meta>gameVersion")

		log.Printf("%+v", p.Find(f.GetXML()).DisplayDebug())
	}
}
