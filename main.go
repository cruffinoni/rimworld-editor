package main

import (
	"github.com/cruffinoni/rimworld-editor/editor"
	"github.com/cruffinoni/rimworld-editor/xml/path"
	"log"
)

func main() {
	f, err := editor.Open("test/part.xml")
	if err != nil {
		log.Fatal(err)
	}

	//p := &path.AttributeMatch{}

	//fmt.Println(f.GetXML().Pretty())
	//for _, p := range f.GetXML().FindElementFromData("Freyja") {
	//	log.Println(p.XMLPath())
	//}

	//log.Printf("Retrieving pawns...")
	//for _, p := range path.FindWithPath("game>maps>li[1]>things", f.GetXML().Root) {
	//	log.Printf("=> %v", p.XMLPath())
	//	for _, d := range path.FindWithPath("li>things>thing>name>first", p) {
	//		log.Printf("d: %v", d)
	//	}
	//}
	for _, p := range path.FindWithPath(`meta>beta>li[2]>faction{Faction}`, f.GetXML().Root) {
		log.Printf("=> %v", p.XMLPath())
		//for _, d := range path.FindWithPath("li>things>thing>name>first", p) {
		//	log.Printf("d: %v", d)
		//}
	}
}
