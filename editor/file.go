package editor

import (
	_xml "encoding/xml"
	"github.com/cruffinoni/rimworld-editor/xml"
	"os"
)

type FileOpening struct {
	fileName string
	content  string
	xml      *xml.Tree
}

func Open(fileName string) (*FileOpening, error) {
	fileOpening := &FileOpening{fileName: fileName}
	content, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	if err = _xml.Unmarshal(content, &fileOpening.xml); err != nil {
		return nil, err
	}
	return fileOpening, nil
}

func (fileOpening *FileOpening) GetXML() *xml.Tree {
	return fileOpening.xml
}
