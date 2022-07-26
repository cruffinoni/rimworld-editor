package generator

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/iancoleman/strcase"
	"go/format"
	"log"
	"os"
	"reflect"
)

type bufferTemplate struct {
	writtenHeaders map[string]bool
	header         string
	body           string
	footer         string
}

const (
	headerXmlTypes = "xml/types"
	headerXml      = "xml"
)

func (t *bufferTemplate) writeToHeader(s string) {
	if v, ok := t.writtenHeaders[s]; ok && v {
		return
	}
	t.writtenHeaders[s] = true
	t.header += `import "github.com/cruffinoni/rimworld-editor/` + s + `"` + "\n"
}

func (t *bufferTemplate) writeToBody(s string) {
	t.body += s
}

func (t *bufferTemplate) writeToFooter(s string) {
	t.footer += s
}

func (t *bufferTemplate) bytes() []byte {
	return []byte(t.header + t.body + t.footer)
}

func getTypeName(t any) string {
	if t == nil {
		panic("nil type")
	}
	switch va := t.(type) {
	case reflect.Kind:
		return va.String()
	case *customType:
		return "*" + va.pkg + "." + va.name
	case *StructInfo:
		return "*" + strcase.ToCamel(va.name)
	case *xml.Element:
		return "*xml.Element"
	default:
		panic("unknown type")
	}
}

func checkTypeAndApply(t any, buffer *bufferTemplate, path string) error {
	switch va := t.(type) {
	case *StructInfo:
		return va.generateStructTo(path)
	case *xml.Element:
		buffer.writeToHeader(headerXml)
	}
	return nil
}

func (s *StructInfo) generateStructTo(path string) error {
	f, err := os.Create(path + "/" + strcase.ToSnake(s.name) + ".go")
	if err != nil {
		return err
	}
	var (
		buffer = &bufferTemplate{
			writtenHeaders: make(map[string]bool),
			header:         "package generated\n\n", // Might be edited later
		}
	)
	defer f.Close()
	buffer.writeToBody("type " + strcase.ToCamel(s.name) + " struct {\n")
	for _, m := range s.members {
		buffer.writeToBody("\t" + strcase.ToCamel(m.name) + " ")
		switch va := m.t.(type) {
		case *customType:
			buffer.writeToHeader(headerXmlTypes)
			buffer.writeToBody(va.pkg + "." + va.name + "[" + getTypeName(va.types1))
			if err = checkTypeAndApply(va.types1, buffer, path); err != nil {
				return err
			}
			if va.types2 != nil {
				if err = checkTypeAndApply(va.types2, buffer, path); err != nil {
					return err
				}
				buffer.writeToBody(", " + getTypeName(va.types2))
			}
			buffer.writeToBody("]")
		case reflect.Kind:
			buffer.writeToBody(va.String())
		case *StructInfo:
			buffer.writeToBody("*" + strcase.ToCamel(va.name))
			if err = va.generateStructTo(path); err != nil {
				return err
			}
		case *xml.Element:
			buffer.writeToHeader(headerXml)
			buffer.writeToBody("*xml.Element")
		}
		buffer.writeToBody(" `xml:\"" + m.name + "\"`\n")
	}
	buffer.writeToFooter("}\n")
	var b []byte = buffer.bytes()
	//log.Printf("Formatting %s...", path+"/"+strcase.ToSnake(s.name)+".go")
	b, err = format.Source(buffer.bytes())
	if err != nil {
		log.Printf("Bytes buffer:\n%s", buffer.bytes())
		return err
	}
	if _, err = f.Write(b); err != nil {
		return err
	}
	return nil
}
