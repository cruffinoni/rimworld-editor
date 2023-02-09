package xmlFile

import (
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/cruffinoni/rimworld-editor/helper"
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/xml/saver"
	"github.com/cruffinoni/rimworld-editor/xml/types/primary"
)

// SaveWithBuffer takes in a value of any type, and returns a saver.Buffer and any error that occurs during the saving process.
func SaveWithBuffer(val any) (*saver.Buffer, error) {
	b := saver.NewBuffer()
	valType := reflect.TypeOf(val)
	if valType.Kind() == reflect.Ptr {
		valType = valType.Elem()
	}
	err := Save(val, b, strings.ToLower(valType.Name()))
	if err == nil {
		b.RemoveEmptyLine()
	}
	return b, err
}

// castToInterface attempts to cast the given value to the specified type and returns the result and a boolean indicating whether the cast was successful.
func castToInterface[T any](val any) (T, bool) {
	if v, ok := val.(T); ok {
		return v, true
	}
	return *new(T), false
}

// Save recursively saves the given value to the provided buffer with the given tag.
func Save(val any, b *saver.Buffer, tag string) error {
	if val == nil {
		return nil
	}
	t := reflect.TypeOf(val)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	v := reflect.ValueOf(val)
	if v.IsZero() && v.Kind() == reflect.Ptr {
		return nil
	}
	transformer, implTransformer := castToInterface[saver.Transformer](v.Interface())
	var attr attributes.Attributes
	//y, z := castToInterface[saver.Transformer](v.Interface())
	//log.Printf("Val: %T & %v/%v", val, y, z)
	if attributeAssigner, ok := castToInterface[xml.AttributeAssigner](v.Interface()); ok {
		attr = attributeAssigner.GetAttributes()
	}
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	vi := v.Interface()
	// If the value is of type primary.Empty, write an empty tag with the given attributes and return.
	if _, ok := val.(*primary.Empty); ok {
		b.WriteEmptyTag(tag, attr)
		return nil
	}
	kind := v.Kind()
	//if kind == reflect.Int64 {
	//	log.Printf("Debug: => %v & %T", val, val)
	//}
	log.Printf("Content: '%v' (%T)", val, val)
	if helper.IsReflectPrimaryType(kind) && (v.IsZero() || v.Kind() == reflect.String && val == "") && !(kind == reflect.Int64 || kind == reflect.Float64) {
		log.Printf("Skipping: %v & %T", val, val)
		return nil
	}
	b.OpenTag(tag, attr)
	switch kind {
	// The `reflect.Slice` case may not be supported in later versions to give way to custom types: `primary.Empty`, `types.Slice`, etc.
	case reflect.Slice:
		j := v.Len()
		// This is a special case in case the type has a custom
		// implementation of the TransformToXML() method.
		if implTransformer {
			if err := transformer.TransformToXML(b); err != nil {
				return err
			}
		}
		for i := 0; i < j; i++ {
			b.Write([]byte("\n"))
			if err := Save(v.Index(i).Interface(), b, "li"); err != nil {
				return err
			}
		}
		if tag != "" {
			b.Write([]byte("\n"))
		}
		b.CloseTagWithIndent(tag)
		return nil
	case reflect.String:
		multipleLineTxt := strings.Contains(v.String(), "\n")
		if multipleLineTxt {
			b.Write([]byte{'\n'})
			b.IncreaseDepth()
		}
		b.Write([]byte(v.String()))
		if multipleLineTxt {
			b.Write([]byte{'\n'})
			b.DecreaseDepth()
		}
	case reflect.Int64:
		b.Write([]byte(strconv.FormatInt(v.Int(), 10)))
	case reflect.Float64:
		b.Write([]byte(strconv.FormatFloat(v.Float(), 'f', -1, 64)))
	case reflect.Struct:
		if vi == nil {
			return nil
		}
		b.Write([]byte("\n"))
		if implTransformer {
			if err := transformer.TransformToXML(b); err != nil {
				return err
			}
			b.Write([]byte("\n"))
			b.CloseTagWithIndent(tag)
			return nil
		}
		for i := 0; i < v.NumField(); i++ {
			f := t.Field(i)
			vf := v.Field(i)
			if !vf.CanInterface() {
				continue
			}
			// The field might be a custom type that implements saver.Transformer.
			// If so, we let the type handle the transformation from Type -> XML
			if transformer, ok := castToInterface[saver.Transformer](vf.Addr().Interface()); ok {
				if err := transformer.TransformToXML(b); err != nil {
					return err
				}
				if tag != "" {
					b.Write([]byte("\n"))
				}
				// We don't close the tag since it's only a MEMBER of the struct, so it can't decide
				// whenever the struct is completely parsed.
				continue
			}
			xmlTag, ok := f.Tag.Lookup("xml")
			if !ok {
				//log.Printf("No XML tag present for %v", f.Name)
				continue
			}
			// Structure that have empty value into their fields are ignored.
			if helper.IsReflectPrimaryType(vf.Kind()) && vf.IsZero() && !(vf.Kind() == reflect.Int64 || vf.Kind() == reflect.Float64) {
				continue
			}
			if err := Save(vf.Interface(), b, xmlTag); err != nil {
				return err
			}
			b.Write([]byte("\n"))
		}
		b.CloseTagWithIndent(tag)
		return nil
	}
	b.CloseTag(tag)
	return nil
}
