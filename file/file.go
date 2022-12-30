package file

import (
	"bytes"
	_xml "encoding/xml"
	"os"

	"golang.org/x/net/html/charset"

	"github.com/cruffinoni/rimworld-editor/xml"
)

type Opening struct {
	fileName string
	XML      *xml.Tree
}

func Open(fileName string) (*Opening, error) {
	fileOpening := &Opening{fileName: fileName}
	content, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(content)
	decoder := _xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel
	if err = decoder.Decode(&fileOpening.XML); err != nil {
		return nil, err
	}
	return fileOpening, nil
}
