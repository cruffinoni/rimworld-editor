package binder

import (
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/cruffinoni/rimworld-editor/internal/xml/domain"
	"github.com/cruffinoni/rimworld-editor/internal/xml/interfaces"
	"github.com/cruffinoni/rimworld-editor/internal/xml/query"
	"github.com/cruffinoni/rimworld-editor/internal/xml/scalar"
	"github.com/cruffinoni/rimworld-editor/internal/xml/support"
	"github.com/cruffinoni/rimworld-editor/pkg/logging"
)

func findFieldFromName(t reflect.Type, value reflect.Value, name string) int {
	for i := 0; i < value.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("xml")
		if tag == "" {
			tag = f.Name
		}
		if tag == name && f.IsExported() {
			return i
		}
	}
	return -1
}

var (
	elementStruct          = reflect.TypeOf((*domain.Element)(nil))
	elementStructName      = elementStruct.Elem().Name()
	elementEmptyStructName = reflect.TypeOf((*scalar.Empty)(nil)).Elem().Name()
)

func isEmptyType(name string) bool {
	return name == elementEmptyStructName
}

func isXMLElement(v reflect.Value) bool {
	if v.Type().Kind() == reflect.Ptr {
		return v.Type().Elem().Name() == elementStructName
	}
	return v.Type().Name() == elementStructName
}

func attributeDataToField(v reflect.Value, e *domain.Element) {
	if e.Data == nil {
		return
	}
	switch v.Kind() {
	case reflect.String:
		v.SetString(e.Data.GetString())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v.SetInt(e.Data.GetInt64())
	case reflect.Bool:
		v.SetBool(e.Data.GetBool())
	case reflect.Float32, reflect.Float64:
		v.SetFloat(e.Data.GetFloat64())
	}
}

func createValueFromPrimaryType(logger logging.Logger, t reflect.Type, e *domain.Element) reflect.Value {
	k := t.Kind()
	if e.Data == nil {
		logger.WithField("type", t.Name()).Debug("createValueFromPrimaryType: no data")
		return reflect.Zero(t)
	}
	d := e.Data
	if d.Kind() != k {
		log.Panicf("createValueFromPrimaryType: type mismatch: %v != %v", k, d.Kind())
	}
	switch k {
	case reflect.String:
		return reflect.ValueOf(d.GetString())
	case reflect.Int64:
		return reflect.ValueOf(d.GetInt64())
	case reflect.Bool:
		return reflect.ValueOf(d.GetBool())
	case reflect.Float64:
		return reflect.ValueOf(d.GetFloat64())
	default:
		log.Panicf("createValueFromPrimaryType: unsupported type %s", t.String())
	}
	// Never reached
	return reflect.Value{}
}

func skipPath(logger logging.Logger, element *domain.Element, pathStr string) *domain.Element {
	p := query.FindWithPath(pathStr, element, logger)

	if len(p) > 1 {
		panic("unmarshal: multiple elements found")
	}
	if len(p) == 0 {
		panic("unmarshal: no element found")
	}
	n := p[0].Child
	return n
}

func Element(logger logging.Logger, element *domain.Element, dest any) error {
	// Do a copy of the element to avoid modifying the original
	n := element
	if n == nil || n.GetName() == "history" {
		return nil
	}

	if setter, ok := dest.(interfaces.LoggerSetter); ok {
		setter.SetLogger(logger)
	}

	destAssigner, destIsAssigner := dest.(interfaces.Assigner)
	if destIsAssigner {
		skippingPath := destAssigner.GetPath()
		if skippingPath != "" {
			n = skipPath(logger, n, skippingPath)
		}
	}
	validator, canValidate := dest.(interfaces.FieldValidator)

	t := reflect.TypeOf(dest)
	if t.Kind() != reflect.Pointer {
		return errors.New("dest must be a pointer")
	}
	t = t.Elem()
	if t.Kind() == reflect.Ptr {
		return fmt.Errorf("multiple pointers found for type %s", reflect.TypeOf(dest).String())
	}
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("unmarshal: path %s: type %s is not a struct", n.XMLPath(), t.String())
	}
	v := reflect.ValueOf(dest).Elem()
	if v.Kind() == reflect.Invalid {
		return errors.New("value of dest is invalid")
	}

	if isXMLElement(v) {
		logger.WithField("type", v.Type().String()).Debug("special case: unmarshal: field is domain.Element")
		//reflect.ValueOf(dest).Set(reflect.ValueOf(&element))
		v.Set(reflect.ValueOf(reflect.ValueOf(element).Elem().Interface().(domain.Element)))
		return nil
	}
	//printer.Debugf("Doing unmarshal for type %s", n.XMLPath())
	for n != nil {
		f := findFieldFromName(t, v, n.GetName())
		//printer.Debugf("n: %v | %v & f: %v", n.GetName(), n.Attr, f)
		if f != -1 && n.GetName() != "history" {
			fieldValue := v.Field(f)
			if canValidate {
				validator.ValidateField(t.Field(f).Name)
			}
			fieldKind := fieldValue.Kind()
			// If the field is a pointer, we need to allocate a new value if it has not been done before
			if fieldKind == reflect.Ptr {
				if fieldValue.IsNil() {
					fieldValue.Set(reflect.New(fieldValue.Type().Elem())) // Allocate new value for the pointer
				}
				fieldValue = fieldValue.Elem()
				fieldKind = fieldValue.Kind()
			}
			switch fieldKind {
			case reflect.Ptr:
				// The function doesn't support multiple pointers (a.k.a. pointers to pointers)
				log.Fatalf("unmarshal: field %v has multiple pointer", fieldValue.Type())
			case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Bool, reflect.Float32, reflect.Float64:
				//printer.Debugf("Attributing data: %v / %s > '%v'", n.XMLPath(), fieldKind.String(), n.Data)
				attributeDataToField(fieldValue, n)
			case reflect.Array:
				l := fieldValue.Len()
				// Create a slice
				fieldValue.Set(reflect.New(reflect.ArrayOf(l, fieldValue.Type().Elem())).Elem())
				// If there is no child element, we are done and left the slice empty
				if n.Child == nil {
					logger.Debug("unmarshal: array empty")
					continue
				}
				// ft is the type of the slice
				ft := fieldValue.Type().Elem()
				// Avoid skipping elements in our linked list
				nChild := n.Child
				idx := 0
				for nChild != nil {
					if idx >= l {
						log.Panicf("index out of range: %v | %v (%d len)", n.XMLPath(), ft.String(), l)
					}
					// Special case for domain.Element, set directly to the field
					if ft == elementStruct {
						fieldValue.Index(idx).Set(reflect.ValueOf(nChild))
					} else if scalar.IsEmbeddedPrimaryType(ft.Name()) || support.IsReflectPrimaryType(ft.Kind()) {
						fieldValue.Index(idx).Set(createValueFromPrimaryType(logger, ft, nChild))
					} else {
						if ft.Kind() != reflect.Ptr {
							panic("unmarshal: array element type must be a pointer")
						}
						newEntry := reflect.New(ft.Elem())
						if nChild.Child == nil {
							newEntry.Interface().(interfaces.Assigner).SetAttributes(nChild.Attr)
						} else {
							if err := Element(logger, nChild.Child, newEntry.Interface().(interfaces.Assigner)); err != nil {
								panic(err)
							}
						}
						fieldValue.Index(idx).Set(newEntry)
					}
					idx++
					nChild = nChild.Next
				}
			case reflect.Struct:
				typeName := fieldValue.Type().Name()
				// Special case for domain.Element, set directly to the field
				//printer.Debugf("unmarshal: struct %v", typeName)
				if isXMLElement(fieldValue) {
					logger.WithField("type", fieldValue.Type().String()).Debug("unmarshal: field is domain.Element")
					v.Field(f).Set(reflect.ValueOf(n))
					break
				} else if scalar.IsEmbeddedPrimaryType(typeName) || isEmptyType(typeName) {
					// it must be a safe cast because the structures are known
					cast := fieldValue.Addr().Interface().(interfaces.Assigner)
					_ = cast.Assign(n)
					cast.SetAttributes(n.Attr)
				} else if assigner, ok := fieldValue.Addr().Interface().(interfaces.Assigner); ok {
					assigner.SetAttributes(n.Attr)
					// Otherwise, we need to call the unmarshal function recursively
					if err := Element(logger, n.Child, fieldValue.Addr().Interface().(interfaces.Assigner)); err != nil {
						panic(err)
					}
				}
			}
		}
		n = n.Next
	}
	if destIsAssigner {
		// The variable element corresponds to the current element of the XML file,
		// but dest might be the parent because we go through all the fields
		// of the struct.
		if t.Kind() == reflect.Struct && element.Parent != nil {
			destAssigner.SetAttributes(element.Parent.Attr)
		} else {
			destAssigner.SetAttributes(element.Attr)
		}
		return destAssigner.Assign(element)
	}
	return nil
}
