package unmarshal

import (
	"errors"
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/path"
	"log"
	"reflect"
)

func findFieldFromName(t reflect.Type, value reflect.Value, name string) int {
	for i := 0; i < value.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("xml")
		if tag == "" {
			tag = f.Name
		}
		if tag == name {
			return i
		}
	}
	return -1
}

var elementStructName = reflect.TypeOf((*xml.Element)(nil)).Elem().Name()

func isXMLElement(v reflect.Value) bool {
	if v.Type().Kind() == reflect.Ptr {
		return v.Type().Elem().Name() == elementStructName
	}
	return v.Type().Name() == elementStructName
}

func Element(element *xml.Element, dest xml.Assigner) error {
	n := element
	if n == nil {
		return nil
	}
	skippingPath := dest.GetPath()
	if skippingPath != "" {
		p := path.FindWithPath(dest.GetPath(), element)

		if len(p) > 1 {
			return errors.New("unmarshal: multiple elements found")
		}
		if len(p) == 0 {
			return errors.New("unmarshal: no element found")
		}
		n = p[0].Child
		log.Printf("unmarshal: skipping path: %v", skippingPath)
	}

	t := reflect.TypeOf(dest)
	if t.Kind() != reflect.Pointer {
		return errors.New("dest must be a pointer")
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return errors.New("unmarshal: dest must be a struct")
	}
	v := reflect.ValueOf(dest).Elem()

	log.Printf("n: %v", n)
	nCpy := n
	for n != nil {
		f := findFieldFromName(t, v, n.GetName())
		log.Printf("Node name: %v & f: %v", n.GetName(), f)
		if f != -1 {
			fieldValue := v.Field(f)
			if fieldValue.Kind() == reflect.Ptr {
				fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
				fieldValue = fieldValue.Elem()
				log.Printf("Kind: %v", fieldValue.Kind())
			}
			switch fieldValue.Kind() {
			case reflect.Ptr:
				log.Printf("Ptr: %v", fieldValue.Type())
			case reflect.String:
				log.Printf("String: %v", n.Data)
				fieldValue.SetString(n.Data.GetString())
			case reflect.Int:
				log.Printf("Int: %v", n.Data)
				fieldValue.SetInt(int64(n.Data.GetInt()))
			case reflect.Bool:
				log.Printf("Bool: %v", n.Data)
				fieldValue.SetBool(n.Data.GetBool())
			case reflect.Slice:
				log.Fatal("unmarshal: slice not implemented")
			case reflect.Struct:
				if isXMLElement(fieldValue) {
					fieldValue.Set(reflect.ValueOf(n))
					log.Printf("Struct is type of xml.Element, setting to %v", n)
					continue
				}
				log.Printf("Struct: %v ; child: %v", fieldValue.String(), n.Child)
				if err := Element(n.Child, fieldValue.Addr().Interface().(xml.Assigner)); err != nil {
					return err
				}
			}
		}
		n = n.Next
	}
	return dest.Assign(nCpy)
}
