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

func SaveWithBuffer(val any) (*saver.Buffer, error) {
	b := saver.NewBuffer()
	valType := reflect.TypeOf(val)
	if valType.Kind() == reflect.Ptr {
		valType = valType.Elem()
	}
	err := Save(val, b, strings.ToLower(valType.Name()))
	return b, err
}

func castToInterface[T any](val any) (T, bool) {
	if v, ok := val.(T); ok {
		return v, true
	}
	return *new(T), false
}

func Save(val any, b *saver.Buffer, tag string) error {
	if val == nil {
		return nil
	}
	t := reflect.TypeOf(val)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	v := reflect.ValueOf(val)
	if v.IsZero() {
		return nil
	}
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	vi := v.Interface()

	var attr attributes.Attributes
	if attributeAssigner, ok := castToInterface[xml.AttributeAssigner](val); ok {
		attr = attributeAssigner.GetAttributes()
	}
	if _, ok := val.(primary.Empty); ok {
		b.WriteEmptyTag(tag, attr)
		return nil
	}
	kind := v.Kind()
	if helper.IsReflectPrimaryType(kind) && v.IsZero() {
		return nil
	}
	b.OpenTag(tag, attr)
	switch kind {
	// The `reflect.Slice` case may not be supported in later versions to give way to custom types: `primary.Empty`, `types.Slice`, etc.
	case reflect.Slice:
		j := v.Len()
		// This is a special case in case the type has a custom
		// implementation of the TransformToXML() method.
		if transformer, ok := castToInterface[saver.Transformer](vi); ok {
			if err := transformer.TransformToXML(b); err != nil {
				return err
			}
		} else {
			for i := 0; i < j; i++ {
				b.Write([]byte("\n"))
				if err := Save(v.Index(i).Interface(), b, "li"); err != nil {
					return err
				}
			}
		}
		b.Write([]byte("\n"))
		b.CloseTagWithIndent(tag)
		return nil
	case reflect.String:
		b.Write([]byte(v.String()))
	case reflect.Int64:
		b.Write([]byte(strconv.FormatInt(v.Int(), 10)))
	case reflect.Float64:
		b.Write([]byte(strconv.FormatFloat(v.Float(), 'f', -1, 64)))
	case reflect.Struct:
		if vi == nil {
			return nil
		}
		b.Write([]byte("\n"))
		if transformer, ok := castToInterface[saver.Transformer](vi); ok {
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
			// The field might be a custom type that implements saver.Transformer.
			// If so, we let the type handle the transformation from Type -> XML
			if transformer, ok := castToInterface[saver.Transformer](vf.Interface()); ok {
				if err := transformer.TransformToXML(b); err != nil {
					return err
				}
				b.Write([]byte("\n"))
				// We don't close the tag since it's only a MEMBER of the struct, so it can't decide
				// whenever the struct is completely parsed.
				continue
			}
			xmlTag, ok := f.Tag.Lookup("xml")
			if !ok {
				log.Print("No XML tag present")
				continue
			}
			// Structure that have empty value into their fields are ignored.
			if helper.IsReflectPrimaryType(vf.Kind()) && vf.IsZero() {
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
