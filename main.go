package main

import (
	"github.com/cruffinoni/rimworld-editor/editor"
	"log"
)

func main() {
	println("Hello, world.")
	if f, err := editor.Open("TEST.rws"); err != nil {
		panic(err)
	} else {
		log.Printf("'%+#v'", f.GetXML().Game.Info)
		//log.Printf("-> %v", f.GetXML().Meta.GetModById("ludeon.rimworld"))
	}
}
