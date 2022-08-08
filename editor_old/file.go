package editor_old

import (
	_xml "encoding/xml"
	"github.com/cruffinoni/rimworld-editor/xml"
	"os"
)

type FileOpening struct {
	fileName string
	content  string
	XML      *xml.Tree
}

func Open(fileName string) (*FileOpening, error) {
	fileOpening := &FileOpening{fileName: fileName}
	content, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	if err = _xml.Unmarshal(content, &fileOpening.XML); err != nil {
		return nil, err
	}
	return fileOpening, nil
}
