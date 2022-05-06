package main

import (
	"github.com/cruffinoni/rimworld-editor/editor"
	"os"
)

func main() {
	if f, err := editor.Open("part.xml"); err != nil {
		panic(err)
	} else {
		_ = f
		//fmt.Printf(f.GetXML().Raw())
		//pawns := f.GetXML().FindTagsFromData("Nikolai")
		//log.Printf("%+v", pawns)
		//if len(pawns) > 0 {
		//	log.Printf("Pawn path: %s", pawns[0].XMLPath())
		//
		os.WriteFile("test.xml", []byte(f.GetXML().Pretty()), 0644)
		//fmt.Println(f.GetXML().Pretty())
		//p := path.NewPathing("meta>game>scenario>parts>storyWatcher>researchManager>progress>storyteller>history>archive>autoRecorderGroups>li>recorders>li>recorders>li>recorders>li>recorders>historyEventsManager>colonistEvents>vals>taleManager>tales>li>pawnData")
		//
		//if n := p.Find(f.GetXML()); n != nil {
		//	log.Printf("Data of node: '%v' (%d)", n, len(n))
		//} else {
		//	log.Println("NODE NOT FOUND")
		//}
	}
}
