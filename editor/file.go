package editor

import (
	"encoding/xml"
	"os"
)

type FileOpening struct {
	fileName string
	content  string
	xml      Savegame
}

func Open(fileName string) (*FileOpening, error) {
	fileOpening := &FileOpening{fileName: fileName}
	content, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	if err = xml.Unmarshal(content, &fileOpening.xml); err != nil {
		return nil, err
	}
	return fileOpening, nil
}

func (fileOpening *FileOpening) GetXML() Savegame {
	return fileOpening.xml
}
