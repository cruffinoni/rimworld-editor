package files

import (
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/cruffinoni/rimworld-editor/algorithm"
	"github.com/cruffinoni/rimworld-editor/generator"
	"github.com/cruffinoni/rimworld-editor/generator/paths"
	"github.com/cruffinoni/rimworld-editor/xml"
)

// WriteGoFile writes the struct Go code to the given path.
// It writes recursively the members of the struct. If a member is a struct,
// it will call WriteGoFile on it.
func WriteGoFile(path string, s *generator.StructInfo) error {
	path = "./" + path
	if _, err := os.Stat(path); err == nil {
		if err = os.RemoveAll(path); err != nil {
			return err
		}
	}
	if err := os.Mkdir(path, os.ModePerm); err != nil {
		return err
	}
	return generateStructToPath(path, s)
}

type generic struct{}

type require interface {
	xml.Assigner
	algorithm.Comparable[generic]
}

var (
	tRequired        = reflect.TypeOf((*require)(nil)).Elem()
	nbRequiredMethod = tRequired.NumMethod()

	localGenericName = reflect.TypeOf(generic{}).Name()
)

func writeRequiredInterfaces(b *buffer, structName string) {
	b.writeImport(paths.XmlAttributes, paths.HeaderXml)
	for i := 0; i < nbRequiredMethod; i++ {
		m := tRequired.Method(i)
		structIdentifier := strings.ToLower(structName[:1])
		b.writeToFooter("\n" +
			"func (" + structIdentifier + " *" + structName + ") ")
		b.writeToFooter(m.Name + "(")
		if m.Type.NumIn() > 0 {
			totalIn := m.Type.NumIn()
			for j := 0; j < totalIn; j++ {
				if j > 0 {
					b.writeToFooter(", ")
				}
				if localGenericName == m.Type.In(j).Name() {
					b.writeToFooter("*" + structName)
				} else {
					if m.Name == "SetAttributes" {
						b.writeToFooter("attr ")
					}
					b.writeToFooter(m.Type.In(j).String())
				}
			}
		}
		b.writeToFooter(")")
		numReturnedValue := m.Type.NumOut()
		returnedValue := make([]reflect.Type, 0, numReturnedValue)
		if numReturnedValue > 0 {
			if numReturnedValue > 1 {
				b.writeToFooter(" (")
			}
			for j := 0; j < numReturnedValue; j++ {
				if j > 0 {
					b.writeToFooter(", ")
				}
				o := m.Type.Out(j)
				if o.Name() == localGenericName {
					returnedValue = append(returnedValue, reflect.TypeOf((*require)(nil)))
					b.writeToFooter("*" + structName)
				} else {
					returnedValue = append(returnedValue, o)
					b.writeToFooter(o.String())
				}
			}
			if numReturnedValue > 1 {
				b.writeToFooter(")")
			}
		}
		b.writeToFooter(" {\n")
		if m.Name == "SetAttributes" {
			b.writeToFooter("\t" + structIdentifier + ".Attr = attr\n")
		}
		b.writeToFooter("\treturn ")
		if numReturnedValue > 0 {
			for c, rt := range returnedValue {
				if c > 0 {
					b.writeToFooter(", ")
				}
				switch rt.Kind() {
				case reflect.Bool:
					b.writeToFooter("false")
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					b.writeToFooter("0")
				case reflect.Float32, reflect.Float64:
					b.writeToFooter("0.0")
				case reflect.String:
					b.writeToFooter(`""`)
				case reflect.Pointer, reflect.Interface, reflect.Slice, reflect.Array, reflect.Map:
					if m.Name == "GetAttributes" {
						b.writeToFooter(structIdentifier + ".Attr")
					} else {
						b.writeToFooter("nil")
					}
				default:
					log.Panicf("generator.StructInfo.writeRequiredInterfaces: unknown type %v", rt.Kind())
				}
			}
		}
		b.writeToFooter("}\n")
	}
}
