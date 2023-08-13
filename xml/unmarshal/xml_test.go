package unmarshal

import (
	"testing"
)

func TestSimpleXML(t *testing.T) {
	s := struct {
		Foo []int `xml:"foo"`
	}{}
}
