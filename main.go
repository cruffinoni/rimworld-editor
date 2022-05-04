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
		_ = f
		log.Println(f.GetXML().Meta.Pretty())
		//log.Printf("Discover: '%v'", f.GetXML().Meta.Pretty(2))
		//log.Printf("->>: '%+#v'", f.GetXML().Meta.Tag.String())
		//log.Printf("->>: '%+#v'", *f.GetXML().Meta.Tag.Next)
		//log.Printf("->>: '%+#v'", *f.GetXML().Meta.Tag.Next.Next)
		//log.Printf("'%+#v'", f)
		//for _, v := range f.GetXML().Game.Scenario.Name {
		//	log.Printf("=> %+#v", v)
		//}
		//for _, t := range f.GetXML().Game.Scenario.Parts.Data {
		//	log.Printf("=> %v", t.String())
		//}
		//log.Printf("-> %v", f.GetXML().Game.Scenario.Parts)
	}
}
