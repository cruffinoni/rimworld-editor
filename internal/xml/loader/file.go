package loader

import (
	"bytes"
	_xml "encoding/xml"
	"os"

	"golang.org/x/net/html/charset"

	"github.com/cruffinoni/rimworld-editor/internal/xml/domain"
)

type Opening struct {
	fileName string
	XML      *domain.Tree
}

func Open(fileName string) (*Opening, error) {
	fileOpening := &Opening{fileName: fileName}
	if err := fileOpening.ReOpen(); err != nil {
		return nil, err
	}
	return fileOpening, nil
}

func (o *Opening) ReOpen() error {
	content, err := os.ReadFile(o.fileName)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(content)
	decoder := _xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel
	if err := decoder.Decode(&o.XML); err != nil {
		return err
	}
	return nil
}

func ReadFromBuffer(buffer string) (*Opening, error) {
	fileOpening := &Opening{fileName: "localbuffer"}
	reader := bytes.NewReader([]byte(buffer))
	decoder := _xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel
	if err := decoder.Decode(&fileOpening.XML); err != nil {
		return nil, err
	}
	return fileOpening, nil
}
