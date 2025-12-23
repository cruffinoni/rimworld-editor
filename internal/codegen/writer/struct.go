package writer

import (
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/cruffinoni/rimworld-editor/internal/codegen"
	"github.com/cruffinoni/rimworld-editor/internal/codegen/importpath"
	"github.com/cruffinoni/rimworld-editor/internal/xml/algorithm"
	"github.com/cruffinoni/rimworld-editor/internal/xml/interfaces"
	"github.com/cruffinoni/rimworld-editor/pkg/logging"
)

type GoWriter struct {
	forcedPackageName string
	registeredMember  codegen.MemberVersioning
	deleteFolder      bool
	logger            logging.Logger
}

func NewGoWriter(logger logging.Logger, registeredMember codegen.MemberVersioning, deleteFolder bool, forcedPackageName string) *GoWriter {
	return &GoWriter{
		forcedPackageName: forcedPackageName,
		registeredMember:  registeredMember,
		deleteFolder:      deleteFolder,
		logger:            logger,
	}
}

var DefaultGoWriter = NewGoWriter(logging.NewLogger("codegen", nil), nil, true, "")

// WriteGoFile writes the struct Go code to the given path.
// It writes recursively the members of the struct. If a member is a struct,
// it will call WriteGoFile on it.
func (gw *GoWriter) WriteGoFile(path string, s *codegen.StructInfo) error {
	if gw.deleteFolder {
		if _, err := os.Stat(path); err == nil {
			if err = os.RemoveAll(path); err != nil {
				return err
			}
		}
	}
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}
	if gw.registeredMember == nil {
		gw.registeredMember = codegen.RegisteredMembers
	}
	return gw.generateStructToPath(path, s)
}

type generic struct{}

type require interface {
	interfaces.Assigner
	interfaces.FieldValidator
	algorithm.Comparable[generic]
}

var (
	tRequired        = reflect.TypeOf((*require)(nil)).Elem()
	nbRequiredMethod = tRequired.NumMethod()

	localGenericName = reflect.TypeOf(generic{}).Name()
)

func writeRequiredInterfaces(b *buffer, structName string) {
	b.writeImport(importpath.XmlAttributes, importpath.HeaderXml)
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
					switch m.Name {
					case "SetAttributes":
						b.writeToFooter("attr ")
					case "ValidateField", "IsValidField":
						b.writeToFooter("field ")
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

		// Function's content
		switch m.Name {
		case "SetAttributes":
			b.writeToFooter("\t" + structIdentifier + ".Attr = attr\n")
		case "ValidateField":
			b.writeToFooter("\t if " + structIdentifier + ".FieldValidated == nil {" +
				"\n\t\t " + structIdentifier + ".FieldValidated = make(map[string]bool)\n\t}\n" +
				"\t" + structIdentifier + ".FieldValidated[field] = true\n",
			)
		case "CountValidatedField":
			b.writeToFooter("\t if " + structIdentifier + ".FieldValidated == nil {" +
				"\n\t\t return 0\n\t}\n",
			)
		}
		b.writeToFooter("\treturn ")
		if numReturnedValue > 0 {
			for c, rt := range returnedValue {
				if c > 0 {
					b.writeToFooter(", ")
				}
				switch rt.Kind() {
				case reflect.Bool:
					switch m.Name {
					case "IsValidField":
						b.writeToFooter(structIdentifier + ".FieldValidated[field]")
					default:
						b.writeToFooter("false")
					}
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					switch m.Name {
					case "CountValidatedField":
						b.writeToFooter("len(" + structIdentifier + ".FieldValidated)")
					default:
						b.writeToFooter("0")
					}
				case reflect.Float32, reflect.Float64:
					b.writeToFooter("0.0")
				case reflect.String:
					b.writeToFooter(`""`)
				case reflect.Pointer, reflect.Interface, reflect.Slice, reflect.Array, reflect.Map:
					switch m.Name {
					case "GetAttributes":
						b.writeToFooter(structIdentifier + ".Attr")
					default:
						b.writeToFooter("nil")
					}
				default:
					log.Panicf("codegen.StructInfo.writeRequiredInterfaces: unknown type %v", rt.Kind())
				}
			}
		}
		b.writeToFooter("}\n")
	}
}
