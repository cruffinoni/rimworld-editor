package unmarshal

import (
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/path"
	"github.com/cruffinoni/rimworld-editor/xml/types/primary"
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

var (
	elementStruct          = reflect.TypeOf((*xml.Element)(nil))
	elementStructName      = elementStruct.Elem().Name()
	elementEmptyStructName = reflect.TypeOf((*primary.Empty)(nil)).Elem().Name()
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

func makePointer(v reflect.Value) reflect.Value {
	v.Set(reflect.New(v.Type().Elem()))
	return v
}

func attributeDataToField(v reflect.Value, e *xml.Element) {
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

func createValueFromPrimaryType(t reflect.Type, e *xml.Element) reflect.Value {
	k := t.Kind()
	if e.Data == nil {
		log.Panicf("createValueFromPrimaryType: no data for %s", t.Name())
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

func skipPath(element *xml.Element, pathStr string) *xml.Element {
	p := path.FindWithPath(pathStr, element)

	if len(p) > 1 {
		panic("unmarshal: multiple elements found")
	}
	if len(p) == 0 {
		panic("unmarshal: no element found")
	}
	n := p[0].Child
	return n
}

func Element(element *xml.Element, dest any) error {
	// Do a copy of the element to avoid modifying the original
	n := element
	if n == nil {
		return nil
	}

	destAssigner, destIsAssigner := dest.(xml.Assigner)
	if destIsAssigner {
		skippingPath := destAssigner.GetPath()
		if skippingPath != "" {
			n = skipPath(n, skippingPath)
		}
	}

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
	for n != nil {
		f := findFieldFromName(t, v, n.GetName())
		if f != -1 {
			fieldValue := v.Field(f)
			fieldKind := fieldValue.Kind()
			var fieldPtr reflect.Value
			// If the field is a pointer, we need to allocate a new value
			if fieldKind == reflect.Ptr {
				fieldPtr = makePointer(fieldValue)
				fieldValue = fieldPtr.Elem()
				fieldKind = fieldValue.Kind()
			}
			switch fieldKind {
			case reflect.Ptr:
				// The function doesn't support multiple pointers (a.k.a. pointers to pointers)
				log.Fatalf("unmarshal: field %v has multiple pointer", fieldValue.Type())
			case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Bool, reflect.Float32, reflect.Float64:
				attributeDataToField(fieldValue, n)
			case reflect.Array:
				l := fieldValue.Len()
				// Create a slice
				fieldValue.Set(reflect.New(reflect.ArrayOf(l, fieldValue.Type().Elem())).Elem())
				// If there is no child element, we are done and left the slice empty
				if n.Child == nil {
					log.Println("unmarshal: array empty")
					continue
				}
				// ft is the type of the slice
				ft := fieldValue.Type().Elem()
				// Let's avoid to skip elements in our linked list
				nBefore := n.Child
				idx := 0
				for nBefore != nil {
					// Special case for xml.Element, set directly to the field
					if ft == elementStruct {
						fieldValue.Index(idx).Set(reflect.ValueOf(nBefore))
					} else if primary.IsEmbeddedPrimaryType(ft.Name()) {
						fieldValue.Index(idx).Set(createValueFromPrimaryType(ft, nBefore))
					} else {
						if ft.Kind() != reflect.Ptr {
							panic("unmarshal: array element type must be a pointer")
						}
						newEntry := reflect.New(ft.Elem())
						if nBefore.Child == nil {
							newEntry.Interface().(xml.Assigner).SetAttributes(nBefore.Attr)
						} else {
							if err := Element(nBefore.Child, newEntry.Interface().(xml.Assigner)); err != nil {
								panic(err)
							}
						}
						fieldValue.Index(idx).Set(newEntry)
					}
					idx++
					if idx >= l {
						log.Panicf("index out of range: %v | %v (%d len)", n.XMLPath(), ft.String(), l)
					}
					nBefore = nBefore.Next
				}
			case reflect.Struct:
				typeName := fieldValue.Type().Name()
				// Special case for xml.Element, set directly to the field
				if isXMLElement(fieldValue) {
					fieldPtr.Set(reflect.ValueOf(n))
					break
				} else if primary.IsEmbeddedPrimaryType(typeName) ||
					isEmptyType(typeName) {
					// it must be a safe cast because the structures are known
					cast := fieldValue.Addr().Interface().(xml.Assigner)
					// We know, for sure, that struct doesn't return any error
					_ = cast.Assign(n)
					cast.SetAttributes(n.Attr)
				} else {
					if fieldValue.Kind() == reflect.Ptr {
						fieldPtr = makePointer(fieldValue)
						fieldValue = fieldPtr.Elem()
					}
					// Otherwise, we need to call the unmarshal function recursively
					if err := Element(n.Child, fieldValue.Addr().Interface().(xml.Assigner)); err != nil {
						panic(err)
					}
				}
			}
		}
		n = n.Next
	}
	if destIsAssigner {
		// The variable element correspond to current element of the xml file
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
