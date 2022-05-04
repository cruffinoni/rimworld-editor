package main

import (
	"github.com/cruffinoni/rimworld-editor/editor"
	"log"
	"time"
)

func main() {
	println("Hello, world.")
	t := time.Now()
	if f, err := editor.Open("TEST.rws"); err != nil {
		panic(err)
	} else {
		log.Printf("Time elapsed: %v", time.Since(t))
		_ = f
		for i := 0; i < 100; i++ {
			f.GetXML().Game.Raw()
		}
		//log.Println(f.GetXML().Meta.Pretty(4))
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
