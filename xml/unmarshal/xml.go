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

var (
	elementStruct     = reflect.TypeOf((*xml.Element)(nil))
	elementStructName = elementStruct.Elem().Name()
)

func isXMLElement(v reflect.Value) bool {
	if v.Type().Kind() == reflect.Ptr {
		return v.Type().Elem().Name() == elementStructName
	}
	return v.Type().Name() == elementStructName
}

func makePointer(v reflect.Value) reflect.Value {
	v.Set(reflect.New(v.Type().Elem()))
	return v.Elem()
}

func attributeDataToField(v reflect.Value, e *xml.Element) {
	if e.Data == nil {
		return
	}
	switch v.Kind() {
	case reflect.String:
		v.SetString(e.Data.GetString())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v.SetInt(int64(e.Data.GetInt()))
	case reflect.Bool:
		v.SetBool(e.Data.GetBool())
	case reflect.Float32, reflect.Float64:
		v.SetFloat(e.Data.GetFloat())
	}
}

func createValueFromPrimaryType(t reflect.Type, e *xml.Element) reflect.Value {
	k := t.Kind()
	if e.Data == nil {
		log.Panicf("createValueFromPrimaryType: no data for %s", t.Name())
	}
	d := e.Data
	if d.GetKind() != k {
		log.Panicf("createValueFromPrimaryType: type mismatch: %v != %v", k, d.GetKind())
	}
	switch k {
	case reflect.String:
		log.Printf("=> '%v'", d.GetString())
		return reflect.ValueOf(d.GetString())
	case reflect.Int:
		return reflect.ValueOf(d.GetInt())
	case reflect.Int8:
		return reflect.ValueOf(d.GetInt8())
	case reflect.Int16:
		return reflect.ValueOf(d.GetInt16())
	case reflect.Int32:
		return reflect.ValueOf(d.GetInt32())
	case reflect.Int64:
		return reflect.ValueOf(d.GetInt64())
	case reflect.Bool:
		return reflect.ValueOf(d.GetBool())
	case reflect.Float32, reflect.Float64:
		return reflect.ValueOf(d.GetFloat())
	default:
		log.Panicf("createValueFromPrimaryType: unsupported type %s", t.String())
	}
	// Never reached
	return reflect.Value{}
}

func Element(element *xml.Element, dest xml.Assigner) error {
	// Do a copy of the element to avoid modifying the original
	n := element
	if n == nil {
		return nil
	}
	// TODO: maybe delete this part
	skippingPath := dest.GetPath()
	if skippingPath != "" {
		log.Printf("Skipping path: '%v'...", skippingPath)
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
	nCpy := n
	for n != nil {
		f := findFieldFromName(t, v, n.GetName())
		if f != -1 {
			fieldValue := v.Field(f)
			fieldKind := fieldValue.Kind()
			// If the field is a pointer, we need to allocate a new value
			if fieldKind == reflect.Ptr {
				fieldValue = makePointer(fieldValue)
				fieldKind = fieldValue.Kind()
			}
			switch fieldKind {
			case reflect.Ptr:
				// The function doesn't support multiple pointers (a.k.a. pointers to pointers)
				log.Fatalf("unmarshal: field %v has multiple pointer", fieldValue.Type())
			case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Bool, reflect.Float32, reflect.Float64:
				attributeDataToField(fieldValue, n)
			case reflect.Slice:
				// Create a slice
				fieldValue.Set(reflect.MakeSlice(fieldValue.Type(), 0, 0))
				// If there is no child element, we are done and left the slice empty
				if n.Child == nil {
					log.Printf("unmarshal: slice empty")
					continue
				}
				// ft is the type of the slice
				ft := fieldValue.Type().Elem()
				// Let's avoid to skip elements in our linked list
				nBefore := n.Child
				for nBefore != nil {
					// Special case for xml.Element, set directly to the field
					if ft == elementStruct {
						fieldValue.Set(reflect.Append(fieldValue, reflect.ValueOf(nBefore)))
					} else {
						fieldValue.Set(reflect.Append(fieldValue, createValueFromPrimaryType(ft, nBefore)))
					}
					nBefore = nBefore.Next
				}
			case reflect.Struct:
				// Special case for xml.Element, set directly to the field
				if isXMLElement(fieldValue) {
					fieldValue.Set(reflect.ValueOf(n))
					break
				}
				// Otherwise, we need to call the unmarshal function recursively
				if err := Element(n.Child, fieldValue.Addr().Interface().(xml.Assigner)); err != nil {
					return err
				}
			}
		}
		n = n.Next
	}
	return dest.Assign(nCpy)
}
