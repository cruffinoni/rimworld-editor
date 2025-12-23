package reflection

import (
	_xml "encoding/xml"
	"errors"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/cruffinoni/rimworld-editor/internal/xml/support"

	"github.com/cruffinoni/rimworld-editor/internal/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/internal/xml/domain"
	"github.com/cruffinoni/rimworld-editor/internal/xml/encoder"
	"github.com/cruffinoni/rimworld-editor/internal/xml/interfaces"
	"github.com/cruffinoni/rimworld-editor/internal/xml/scalar"
	"github.com/cruffinoni/rimworld-editor/pkg/logging"
)

// SaveWithBuffer takes in a value of multiple type, and returns a encoder.Buffer and multiple error that occurs during the saving process.
func SaveWithBuffer(logger logging.Logger, val any) (*encoder.Buffer, error) {
	b := encoder.NewBuffer(logger)
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

// castTo attempts to cast the given value to the specified type and returns the result and a boolean indicating whether the cast was successful.
func castTo[T any](val any) (T, bool) {
	if v, ok := val.(T); ok {
		return v, true
	}
	return *new(T), false
}

func capitalize(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}

// Save recursively saves the given value to the provided buffer with the given tag.
func Save(val any, b *encoder.Buffer, tag string) error {
	if val == nil {
		return nil
	}
	logger := b.Logger()
	t := reflect.TypeOf(val)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	v := reflect.ValueOf(val)
	if v.IsZero() && v.Kind() == reflect.Ptr {
		return nil
	}
	validator, implValidator := castTo[interfaces.FieldValidator](v.Interface())
	transformer, implTransformer := castTo[encoder.Transformer](v.Interface())
	var attr attributes.Attributes
	if attributeAssigner, ok := castTo[interfaces.AttributeAssigner](v.Interface()); ok {
		attr = attributeAssigner.GetAttributes()
	}
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	vi := v.Interface()
	valKind := v.Kind()
	// If the value is of type scalar.Empty, write an empty tag with the given attributes and return.
	if _, ok := val.(*scalar.Empty); ok || valKind == reflect.Ptr && v.IsZero() {
		b.WriteEmptyTag(tag, attr)
		return nil
	}
	//if valKind == reflect.Int64 {
	//	printer.Debugf("Debug: => %v & %T", val, val)
	//}
	//printer.Debugf("Content: '%v' (%T)", val, val)
	if support.IsReflectPrimaryType(valKind) && (v.IsZero() || v.Kind() == reflect.String && val == "") && !(valKind == reflect.Int64 || valKind == reflect.Float64) {
		logger.WithFields(logging.Fields{
			"value": val,
			"type":  t.String(),
		}).Debug("Skipping empty value")
		return nil
	}
	_, isXMLElement := val.(*domain.Element)
	// Only non xmL.Element types can open a tag. The structure already does that by itself.
	if !isXMLElement {
		b.OpenTag(tag, attr)
	}
	//printer.Debugf("Kind: %v & %T", valKind, val)
	switch valKind {
	case reflect.Array:
		j := v.Len()
		if j == 0 {
			b.WriteEmptyTag(tag, attr)
			return nil
		}
		transformer, implTransformer = castTo[encoder.Transformer](v.Index(0).Interface())
		// This is a special case in case the type has a custom
		// implementation of the TransformToXML() method.
		if implTransformer {
			if err := transformer.TransformToXML(b); err != nil {
				if errors.Is(err, encoder.ErrEmptyValue) {
					b.RevertToLatestPoint()
					b.WriteEmptyTag(tag, attr)
					return nil
				}
			}
			if !isXMLElement {
				_, _ = b.Write([]byte("\n"))
				b.CloseTagWithIndent(tag)
			}
			return nil
		}
		for i := 0; i < j; i++ {
			_, _ = b.Write([]byte("\n"))
			vIndex := v.Index(i)
			// Sometimes it does not detect nil as it should
			if vIndex.Kind() == reflect.Ptr && vIndex.IsNil() {
				continue
			}
			idxInterface := v.Index(i).Interface()
			if fieldValidator, ok := idxInterface.(interfaces.FieldValidator); ok && fieldValidator != nil && fieldValidator.CountValidatedField() == 0 {
				if attributeAssigner, ok := castTo[interfaces.AttributeAssigner](idxInterface); ok {
					b.WriteEmptyTag("li", attributeAssigner.GetAttributes())
					continue
				} else {
					log.Fatal("can't cast child to attribute assigner")
				}
			}
			if err := Save(idxInterface, b, "li"); err != nil {
				return err
			}
		}
		if tag != "" {
			_, _ = b.Write([]byte("\n"))
		}
		b.CloseTagWithIndent(tag)
		return nil
	case reflect.String:
		multipleLineTxt := strings.Contains(v.String(), "\n")
		if multipleLineTxt {
			_, _ = b.Write([]byte{'\n'})
			b.IncreaseDepth()
		}
		_xml.Escape(b, []byte(v.String()))
		if multipleLineTxt {
			_, _ = b.Write([]byte{'\n'})
			b.DecreaseDepth()
		}
	case reflect.Bool:
		_, _ = b.Write([]byte(capitalize(strconv.FormatBool(v.Bool()))))
	case reflect.Int64:
		_, _ = b.Write([]byte(strconv.FormatInt(v.Int(), 10)))
	case reflect.Float64:
		_, _ = b.Write([]byte(strconv.FormatFloat(v.Float(), 'f', -1, 64)))
	case reflect.Struct:
		if vi == nil {
			return nil
		}
		if implTransformer {
			if err := transformer.TransformToXML(b); err != nil {
				if errors.Is(err, encoder.ErrEmptyValue) {
					b.RevertToLatestPoint()
					b.WriteEmptyTag(tag, attr)
					return nil
				}
				return err
			}
			if !isXMLElement {
				_, _ = b.Write([]byte("\n"))
				b.CloseTagWithIndent(tag)
			}
			return nil
		}
		if implValidator && validator.CountValidatedField() == 0 {
			b.RevertToLatestPoint()
			b.WriteEmptyTag(tag, attr)
			return nil
		}
		_, _ = b.Write([]byte("\n"))
		for i := 0; i < v.NumField(); i++ {
			f := t.Field(i)
			vf := v.Field(i)
			if !vf.CanInterface() {
				continue
			}
			xmlTag, ok := f.Tag.Lookup("xml")
			if !ok {
				continue
			}
			//if !implValidator {
			//	printer.Debugf("Validator implemented for %T (field %v): invalid field", v.Interface(), f.Name)
			//} else {
			//	printer.Debugf("Validator implemented for %T (field %v): > %v", v.Interface(), f.Name, validator.IsValidField(f.Name))
			//}
			if implValidator && !validator.IsValidField(f.Name) {
				//printer.Debugf("Ignoring field %v", f.Name)
				continue
			}
			// The field might be a custom type that implements encoder.Transformer.
			// If so, we let the type handle the transformation from Type -> XML
			if transformer, ok := castTo[encoder.Transformer](vf.Addr().Interface()); ok {
				if err := transformer.TransformToXML(b); err != nil {
					return err
				}
				if tag != "" {
					_, _ = b.Write([]byte("\n"))
				}
				// We don't close the tag since it's only a MEMBER of the struct, so it can't decide
				// whenever the struct is completely parsed.
				continue
			}
			// Structure that have empty value into their fields are ignored.
			if support.IsReflectPrimaryType(vf.Kind()) && vf.IsZero() && !(vf.Kind() == reflect.Int64 || vf.Kind() == reflect.Float64) {
				continue
			}
			if err := Save(vf.Interface(), b, xmlTag); err != nil {
				return err
			}
			_, _ = b.Write([]byte("\n"))
		}
		b.CloseTagWithIndent(tag)
		return nil
	default:
		log.Panicf("Unhandled type: %v", valKind)
	}
	b.CloseTag(tag)
	return nil
}
