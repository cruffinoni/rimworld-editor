package main

import (
	"github.com/cruffinoni/rimworld-editor/editor"
	"github.com/cruffinoni/rimworld-editor/xml/path"
	"log"
	"os"
)

func main() {
	if f, err := editor.Open("part.xml"); err != nil {
		panic(err)
	} else {
		_ = f
		// Nikolai
		pawns := f.GetXML().FindTagsFromData("0")
		if len(pawns) > 0 {
			for _, p := range pawns {
				log.Printf("Pawn xmlPath: %s", p.XMLPath())
				xmlPath := path.FindWithPath(p.XMLPath(), f.GetXML())
				log.Printf("Data of node: '%v' (%d)", xmlPath, len(xmlPath))
			}
		} else {
			log.Println("NODE NOT FOUND")
		}

		os.WriteFile("test.txt", []byte(f.GetXML().Pretty()), 0644)

		//fmt.Println(f.GetXML().Pretty())
		//p := path.NewPathing("game>taleManager>tales>li>pawnData>name>first")
		//
		//if n := p.Find(f.GetXML()); n != nil {
		//	log.Printf("Data of node: '%v' (%d)", n, len(n))
		//} else {
		//	log.Println("NODE NOT FOUND")
		//}
	}
}
