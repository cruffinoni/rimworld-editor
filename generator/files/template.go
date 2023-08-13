package files

import (
	"errors"
	"fmt"
	"go/format"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/cruffinoni/rimworld-editor/generator"
	"github.com/cruffinoni/rimworld-editor/generator/paths"
	"github.com/cruffinoni/rimworld-editor/xml"
)

func getTypeName(t any) string {
	if t == nil {
		panic("nil type")
	}
	switch va := t.(type) {
	case reflect.Kind:
		return va.String()
	case *generator.CustomType:
		var s strings.Builder
		if va.Name == "Empty" {
			s.WriteString(va.Pkg + "." + va.Name)
			return s.String()
		}
		s.WriteString(va.Pkg + "." + va.Name + "[" + getTypeName(va.Type1))
		if va.Type2 != nil {
			s.WriteString(", " + getTypeName(va.Type2))
		}
		s.WriteString("]")
		return s.String()
	case *generator.StructInfo:
		return "*" + strcase.ToCamel(va.Name)
	case *generator.FixedArray:
		return fmt.Sprintf("[%d] %s", va.Size, getTypeName(va.PrimaryType))
	case *xml.Element:
		return "*xml.Element"
	default:
		panic("unknown type")
	}
}

func (gw *GoWriter) checkTypeAndApply(t any, buffer *buffer, path string) error {
	switch va := t.(type) {
	case *generator.StructInfo:
		return gw.generateStructToPath(path, va)
	case *xml.Element:
		buffer.writeImport(paths.HeaderXml)
	case *generator.CustomType:
		var err error
		buffer.writeImport(va.ImportPath)
		if err = gw.checkTypeAndApply(va.Type1, buffer, path); err != nil {
			return err
		}
		if va.Type2 != nil {
			if err = gw.checkTypeAndApply(va.Type2, buffer, path); err != nil {
				return err
			}
		}
	case *generator.FixedArray:
		buffer.writeToBody(fmt.Sprintf("[%d]", va.Size))
		if err := gw.checkTypeAndApply(va.PrimaryType, buffer, path); err != nil {
			return err
		}
	}
	return nil
}

func (gw *GoWriter) writeHeader(b *buffer, pkgName string) {
	b.writeToHeader("// Code generated by rimworld-editor. DO NOT EDIT.\n\n")
	b.writeToHeader("package " + pkgName + "\n\n")
}

func removeInnerKeyword(s string) string {
	return strings.Replace(s, generator.InnerKeyword, "", -1)
}

func (gw *GoWriter) writeCustomType(c *generator.CustomType, b *buffer, path string) error {
	var err error
	b.writeImport(c.ImportPath)
	// log.Printf("Custom type %+v", *c)
	b.writeToBody(c.Pkg + "." + c.Name)
	if c.Type1 == nil {
		// log.Printf("Types: %v & %v", c.type1, c.type2)
		return nil
	}
	b.writeToBody("[" + getTypeName(c.Type1))
	if err = gw.checkTypeAndApply(c.Type1, b, path); err != nil {
		return err
	}
	if c.Type2 != nil {
		if err = gw.checkTypeAndApply(c.Type2, b, path); err != nil {
			return err
		}
		b.writeToBody(", " + getTypeName(c.Type2))
	}
	b.writeToBody("]")
	return err
}

const basicPackageName = "generated"

func (gw *GoWriter) generateStructToPath(path string, s *generator.StructInfo) error {
	if _, err := os.Stat(path + "/" + strcase.ToSnake(s.Name) + ".go"); !errors.Is(err, os.ErrNotExist) {
		// log.Printf("generateStructToPath: file already exists at: %v", path+"/"+strcase.ToSnake(s.Name)+".go")
		// log.Printf("Size: %d from %p", len(s.Members), s)
		return nil
	}
	f, err := os.Create(path + "/" + strcase.ToSnake(s.Name) + ".go")
	if err != nil {
		return err
	}
	buf := &buffer{
		writtenHeaders: make(map[string]bool),
	}
	if gw.forcedPackageName != "" {
		gw.writeHeader(buf, gw.forcedPackageName)
	} else {
		gw.writeHeader(buf, basicPackageName)
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			log.Fatalf("generator.StructInfo.generateStructToPath: can't close the file %v", err)
		}
	}(f)
	structName := strcase.ToCamel(s.Name)
	if structName == "" {
		panic("empty struct name")
	}
	buf.writeToBody("type " + structName + " struct {\nAttr attributes.Attributes\nFieldValidated map[string]bool\n\n")
	// log.Printf("S: %s | %d", s.Name, len(registeredMembers[s.Name]))
	for _, m := range gw.registeredMember[s.Name][0].Order { // Use the best matched version of s.name
		buf.writeToBody("\t" + strcase.ToCamel(m.Name) + " ")
		switch va := m.T.(type) {
		case *generator.CustomType:
			if err = gw.writeCustomType(va, buf, path); err != nil {
				return err
			}
		case reflect.Kind:
			buf.writeToBody(va.String())
		case *generator.StructInfo:
			buf.writeToBody("*" + strcase.ToCamel(va.Name))
			if s.Name == va.Name {
				log.Panicf("duplicate name for %s & %s", s.Name, va.Name)
			}
			if err = gw.generateStructToPath(path, va); err != nil {
				return err
			}
		case *xml.Element:
			// headerXml will be imported in the buffer when we write the
			// required import statement.
			buf.writeToBody("*xml.Element")
		case *generator.FixedArray:
			buf.writeToBody(fmt.Sprintf("[%d] %s", va.Size, getTypeName(va.PrimaryType)))
			if err := gw.checkTypeAndApply(va.PrimaryType, buf, path); err != nil {
				return err
			}
		}
		buf.writeToBody(" `xml:\"" + removeInnerKeyword(m.Name) + "\"`\n")
	}
	buf.writeToFooter("}\n")
	writeRequiredInterfaces(buf, structName)
	var b []byte
	b, err = format.Source(buf.bytes())
	if err != nil {
		log.Printf("Err: Format buffer:\n%s", buf.bytes())
		return err
	}
	if _, err = f.Write(b); err != nil {
		return err
	}
	return nil
}
