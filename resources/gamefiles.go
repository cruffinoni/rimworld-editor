package resources

import (
	"fmt"
	"github.com/cruffinoni/rimworld-editor/cmd/app/ui/term/printer"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/cruffinoni/rimworld-editor/file"
	"github.com/cruffinoni/rimworld-editor/generator"
	"github.com/cruffinoni/rimworld-editor/generator/files"
	"github.com/cruffinoni/rimworld-editor/resources/discover"
	"github.com/cruffinoni/rimworld-editor/xml"
)

type GameData struct {
	fileData       map[string]*GroupedThematic
	accessibleData any
}

type GroupedThematic struct {
	elements map[string]*xml.Element
	cap      int
	mv       generator.MemberVersioning
	roots    []*generator.StructInfo
}

func NewGameData() *GameData {
	return &GameData{
		fileData: make(map[string]*GroupedThematic),
	}
}

func (g *GameData) PrintThemes() {
	log.Println("Game fileData themes:")
	for theme, elements := range g.fileData {
		log.Printf("  %s: %d elements", theme, elements.cap)
	}
}

func (g *GameData) DiscoverGameData() error {
	gp, err := discover.GetGamePath()
	if err != nil {
		return err
	}
	gp = filepath.Join(gp, "Data")
	log.Println(gp)
	err = filepath.Walk(gp, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking %s: %w", path, err)
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(info.Name(), ".xml") {
			return nil
		}
		var f *file.Opening
		f, err = file.Open(path)
		if err != nil {
			return fmt.Errorf("error opening %s: %w", path, err)
		}
		dir := filepath.Base(filepath.Dir(path))
		//log.Printf("Full dir of %s: %s", dir, filepath.Dir(path))
		if g.fileData[dir] == nil {
			g.fileData[dir] = &GroupedThematic{elements: make(map[string]*xml.Element)}
		}
		g.fileData[dir].elements[path] = f.XML.Root
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (gp *GroupedThematic) mergeMemberVersioning(from generator.MemberVersioning) {
	if gp.mv == nil {
		gp.mv = make(generator.MemberVersioning)
	}
	for k, v := range from {
		gp.mv[k] = append(gp.mv[k], v...)
	}
}

func (g *GameData) GenerateGoFiles() error {
	basePath := "./resources/generated/"
	if _, err := os.Stat(basePath); err == nil {
		if err = os.RemoveAll(basePath); err != nil {
			return err
		}
	}
	for p := range g.fileData {
		printer.PrintSf("[%s] Processing {-BOLD,F_RED}%d{-RESET} element(s) ...", p, len(g.fileData[p].elements))
		for _, element := range g.fileData[p].elements {
			root := generator.GenerateGoFiles(element, false)
			g.fileData[p].mergeMemberVersioning(generator.RegisteredMembers)
			g.fileData[p].roots = append(g.fileData[p].roots, root)
		}
		generator.FixRegisteredMembers(g.fileData[p].mv)
		for _, r := range g.fileData[p].roots {
			if err := files.WriteGoFile(basePath+p, r, false, g.fileData[p].mv); err != nil {
				log.Fatal(err)
			}
		}
	}
	return nil
}

func (g *GameData) ReadGameFiles() {
	// TODO: to do
}
