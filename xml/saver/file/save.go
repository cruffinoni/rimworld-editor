package file

import (
	"reflect"
	"strconv"

	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/xml/saver"
	"github.com/cruffinoni/rimworld-editor/xml/types/primary"
)

func SaveWithBuffer(val any, tag string) error {
	b := saver.NewBuffer()
	return Save(val, b, tag)
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
	vInterface := v.Interface()

	var attr attributes.Attributes
	if attributeAssigner, ok := castToInterface[xml.AttributeAssigner](val); ok {
		attr = attributeAssigner.GetAttributes()
	}

	if _, ok := val.(primary.Empty); ok {
		b.WriteEmptyTag(tag, attr)
		return nil
	}
	b.OpenTag(tag, attr)
	kind := v.Kind()
	switch kind {
	case reflect.Slice:
		j := v.Len()
		// This is a special case in case the type has a custom
		// implementation of the TransformToXML() method.
		if transformer, ok := castToInterface[saver.Transformer](vInterface); ok {
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
	case reflect.Int:
		b.Write([]byte(strconv.FormatInt(v.Int(), 10)))
	case reflect.Float64:
		b.Write([]byte(strconv.FormatFloat(v.Float(), 'f', -1, 64)))
	case reflect.Struct:
		if vInterface == nil {
			return nil
		}
		b.Write([]byte("\n"))
		// If `vInterface` implements Transformer interface, retrieve it
		// as a Transformer and call TransformToXML() method.
		if transformer, ok := castToInterface[saver.Transformer](vInterface); ok {
			if err := transformer.TransformToXML(b); err != nil {
				return err
			}
			b.Write([]byte("\n"))
			b.CloseTagWithIndent(tag)
			return nil
		}
		for i := 0; i < v.NumField(); i++ {
			f := t.Field(i)
			vTag, ok := f.Tag.Lookup("xml")
			if !ok {
				continue
			}
			vField := v.Field(i)
			// Ignore zero value fields.
			if vField.IsZero() {
				continue
			}
			if err := Save(vField.Interface(), b, vTag); err != nil {
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
