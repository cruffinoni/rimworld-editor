package main

import (
	"github.com/cruffinoni/rimworld-editor/editor"
	"github.com/cruffinoni/rimworld-editor/xml/unmarshal"
	"log"
)

func main() {
	f, err := editor.Open("test/alone.xml")
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
	var m editor.Savegame
	err = unmarshal.Element(f.XML.Root, &m)
	if err != nil {
		return
	}
	//log.Printf("%+v & err: %v", m.Meta.GetMods(), err)
	//for _, p := range path.FindWithPath(`meta>beta>li[2]>faction{Faction}`, f.XML.Root) {
	//	log.Printf("=> %v", p.XMLPath())
	//	fmt.Println(p.ToXML(0))
	//for _, d := range path.FindWithPath("li>things>thing>name>first", p) {
	//	log.Printf("d: %v", d)
	//}
	//}
}
