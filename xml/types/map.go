package types

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"sort"
	"unicode"

	"github.com/cruffinoni/rimworld-editor/helper"
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/xml/path"
	"github.com/cruffinoni/rimworld-editor/xml/saver"
	"github.com/cruffinoni/rimworld-editor/xml/types/iterator"
	"github.com/cruffinoni/rimworld-editor/xml/types/primary"
)

// Map is a map of K to V and require the XML file to have a "keys" and "values"
// element.

type Pair[K, V any] struct {
	Key   K
	Value V
}

func (p *Pair[K, V]) Equal(rhs *Pair[K, V]) bool {
	//return p.Key == rhs.Key && p.Value == rhs.Value
	return true
}

type MapComparable[T any] interface {
	Less(key reflect.Value, other T) bool
	Equal(key reflect.Value, other T) bool
	Great(key reflect.Value, other T) bool
}

// Map is a map of K to V.
// We don't restrict the type K to MapComparable[Map[K, V]] because K might be
// type of string, int or any primary type.
type Map[K comparable, V any] struct {
	MapComparable[Map[K, V]]
	xml.Assigner
	iterator.MapIndexer[K, V]
	m          map[K]V
	sortedKeys []reflect.Value

	tag  string
	attr attributes.Attributes
}

func castToInterface[T any](val any) (T, bool) {
	if v, ok := val.(T); ok {
		return v, true
	}
	return *new(T), false
}

func (m *Map[K, V]) TransformToXML(b *saver.Buffer) error {
	b.OpenTag(m.tag, m.attr)
	b.Write([]byte("\n"))
	defer func() {
		b.Write([]byte("\n"))
		b.CloseTagWithIndent(m.tag)
	}()
	if m.m == nil || m.Capacity() == 0 {
		b.WriteEmptyTag("keys", nil)
		b.WriteEmptyTag("values", nil)
		return nil
	}
	b.IncreaseDepth()
	b.WriteStringWithIndent("<keys>\n")
	b.IncreaseDepth()
	for k := range m.m {
		b.WriteStringWithIndent("<li>")
		if transformer, ok := castToInterface[saver.Transformer](&k); ok {
			tag := transformer.GetXMLTag()
			if tag != nil {
				b.OpenTag(string(tag), nil)
			}
			if err := transformer.TransformToXML(b); err != nil {
				return err
			}
			if tag != nil {
				b.CloseTagWithIndent(string(tag))
				b.Write([]byte{'\n'})
			}
		} else {
			b.WriteString(fmt.Sprintf("%v", k))
		}
		if unicode.IsSpace(rune(b.Buffer()[b.Len()-1])) {
			b.WriteStringWithIndent("</li>\n")
		} else {
			b.WriteString("</li>\n")
		}
	}
	b.DecreaseDepth()
	b.WriteStringWithIndent("</keys>\n")
	b.WriteStringWithIndent("<values>\n")
	b.IncreaseDepth()
	for _, v := range m.m {
		b.WriteStringWithIndent("<li>")
		if transformer, ok := castToInterface[saver.Transformer](v); ok {
			tag := transformer.GetXMLTag()
			if tag != nil {
				b.OpenTag(string(tag), nil)
			}
			if err := transformer.TransformToXML(b); err != nil {
				return err
			}
			if tag != nil {
				b.CloseTagWithIndent(string(tag))
				b.Write([]byte{'\n'})
			}
		} else {
			b.WriteString(fmt.Sprintf("%v", v))
		}
		if unicode.IsSpace(rune(b.Buffer()[b.Len()-1])) {
			b.WriteStringWithIndent("</li>\n")
		} else {
			b.WriteString("</li>\n")
		}
	}
	b.DecreaseDepth()
	b.WriteStringWithIndent("</values>")
	b.DecreaseDepth()
	return nil
}

func zero[T any]() T {
	return *new(T)
}

func castTemplate[T any](value any) T {
	if v, ok := value.(T); ok {
		return v
	}
	log.Panicf("Map/castTemplate: cannot cast %T to %T", value, zero[T]())
	// Never reached
	return zero[T]()
}

func (m *Map[K, V]) GetXMLTag() []byte {
	return nil
}

func castDataFromKind[T any](kind reflect.Kind, d *xml.Data) T {
	switch kind {
	case reflect.String:
		return castTemplate[T](d.GetString())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return castTemplate[T](d.GetInt64())
	case reflect.Bool:
		return castTemplate[T](d.GetBool())
	case reflect.Float32, reflect.Float64:
		return castTemplate[T](d.GetFloat64())
	}
	log.Panicf("Map/castDataFromKind: cannot cast %T to %T", d.GetData(), zero[T]())
	// Never reached
	return zero[T]()
}

func (m *Map[K, V]) Assign(e *xml.Element) error {
	m.m = make(map[K]V)
	if e.Parent != nil {
		m.tag = e.Parent.GetName()
		m.attr = e.Parent.GetAttributes()
	} else {
		return fmt.Errorf("map.Assign: map's parent is nil")
	}
	if e.Child == nil {
		return nil
	}
	//log.Printf("Tag: %v", m.tag)
	keys := path.FindWithPath("keys>[...]", e)
	if len(keys) == 0 {
		return errors.New("Map/Assign: no key")
	}
	//log.Printf("e=%v", e.GetName())
	values := path.FindWithPath("values>[...]", e)
	if len(values) == 0 {
		return errors.New("Map/Assign: no value")
	}
	if len(keys) != len(values) {
		return errors.New("Map/Assign: keys length differs from values length")
	}
	//log.Printf("Keys: %v, Value: %v", keys[0].XMLPath(), values[0].XMLPath())
	//log.Printf("Keys: %+v, Value: %+v", keys[0].Data, values[0].Data)
	kKind := reflect.TypeOf(zero[K]()).Kind()
	vKind := reflect.TypeOf(zero[V]()).Kind()
	_, isEmpty := any(zero[V]()).(*primary.Empty)
	//log.Printf("%T is empty ? %v", zero[V](), isEmpty)
	for i, key := range keys {
		if key.Data == nil {
			log.Panicf("Map/Assign: no data for key %s", key.StartElement.Name.Local)
		}
		////log.Printf("There is data? %+v\n> %+v\n>> %+v\n>>>%+v", values[i], values[i].Child, values[i].Child.Child, values[i].Child.Child.Child)
		//log.Printf("Type ok? '%T'", zero[V]())
		// This might be a custom type that implements xml.Assigner interface
		if _, ok := any(zero[V]()).(xml.Assigner); ok && !isEmpty {
			if values[i].Child == nil {
				log.Printf("Map/Assign: no child at %s | Index: %d", e.XMLPath(), i)
			}
			var (
				subValue    = new(V)
				subValueVal = reflect.ValueOf(subValue)
			)
			//log.Printf("! > %v & %T + '%v' & POSSIBLE? %v", subValue, subValue, subValueVal.Kind(), subValueVal.Elem().CanAddr())
			subValueVal = subValueVal.Elem()
			if subValueVal.Kind() == reflect.Ptr {
				subValueVal.Set(reflect.New(subValueVal.Type().Elem()))
			}
			if err := subValueVal.Interface().(xml.Assigner).Assign(values[i].Child); err != nil { // Use the child because we don't to assign the <li> tag element
				return err
			}
			//log.Printf("Res: %v (%T)", subValueVal.Interface().(V), subValueVal.Interface().(V))
			m.m[castDataFromKind[K](kKind, key.Data)] = subValueVal.Interface().(V)
			//log.Printf("!!=> %v > %v", castDataFromKind[K](kKind, key.Data), m.m[castDataFromKind[K](kKind, key.Data)])
		} else if values[i].Data == nil || isEmpty {
			// There is a key with no data
			m.m[castDataFromKind[K](kKind, key.Data)] = zero[V]()
		} else if _, isElement := any(zero[V]()).(*xml.Element); isElement {
			// Special if V is a xml.Element because we pass a pointer to the data for castDataFromKind
			// so, we don't use this function but assign directly to the map
			m.m[castDataFromKind[K](kKind, key.Data)] = castTemplate[V](values[i])
		} else {
			m.m[castDataFromKind[K](kKind, key.Data)] = castDataFromKind[V](vKind, values[i].Data)
		}
		//log.Printf("=> %v > %v", castDataFromKind[K](kKind, key.Data), m.m[castDataFromKind[K](kKind, key.Data)])
	}
	v := reflect.ValueOf(m.m)
	k := reflect.ValueOf(zero[K]()).Kind()
	m.sortedKeys = v.MapKeys()
	// Primary type implements natively operator<
	if helper.IsReflectPrimaryType(k) {
		switch k {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			sort.Slice(m.sortedKeys, func(i, j int) bool {
				return m.sortedKeys[i].Int() < m.sortedKeys[j].Int()
			})
		case reflect.String:
			sort.Slice(m.sortedKeys, func(i, j int) bool {
				return m.sortedKeys[i].String() < m.sortedKeys[j].String()
			})
		case reflect.Bool:
			panic("Map/Assign: cannot sort bool")
		case reflect.Float32, reflect.Float64:
			sort.Slice(m.sortedKeys, func(i, j int) bool {
				return m.sortedKeys[i].Float() < m.sortedKeys[j].Float()
			})
		}
	} else {
		// Custom type must implement operator<
		if !reflect.TypeOf(zero[K]()).Implements(reflect.TypeOf((*MapComparable[Map[K, V]])(nil)).Elem()) {
			panic("Map/Assign: custom type must implement MapComparable interface")
		}
		sort.Slice(m.sortedKeys, func(i, j int) bool {
			return m.sortedKeys[i].Interface().(MapComparable[K]).Less(m.sortedKeys[i], m.sortedKeys[j].Interface().(K))
		})
	}
	return nil
}

func (m *Map[K, V]) GetPath() string {
	return ""
}

func (m *Map[K, V]) Get(key K) V {
	if m.m == nil {
		return zero[V]()
	}
	return m.m[key]
}

func (m *Map[K, V]) GetFromIndex(idx int) V {
	if m.m == nil {
		return zero[V]()
	}
	if idx < 0 || idx >= len(m.m) {
		log.Panic("Map/GetFromIndex: index out of range")
		return zero[V]()
	}
	i := 0
	for _, v := range m.m {
		if i == idx {
			return v
		}
		i++
	}
	log.Panicf("Map/GetFromIndex: index %d not found", idx)
	return zero[V]()
}

func (m *Map[K, V]) GetKeyFromIndex(idx int) K {
	if m.m == nil {
		return zero[K]()
	}
	if idx < 0 || idx >= len(m.m) {
		log.Panic("Map/GetFromIndex: index out of range")
		return zero[K]()
	}
	for i, k := range m.sortedKeys {
		if i == idx {
			return k.Interface().(K)
		}
		i++
	}
	log.Panicf("Map/GetKeyFromIndex: index %d not found", idx)
	// Never reached
	return zero[K]()
}

func (m *Map[K, V]) Set(key K, value V) {
	if m.m == nil {
		m.m = make(map[K]V)
	}
	m.m[key] = value
}

func (m *Map[K, V]) Capacity() int {
	return len(m.m)
}

func (m *Map[K, V]) Iterator() *iterator.MapIterator[K, V] {
	return iterator.NewMapIterator[K, V](m)
}

func (m *Map[K, V]) SetAttributes(_ attributes.Attributes) {
	// No attributes need to be set.
}

func (m *Map[K, V]) GetAttributes() attributes.Attributes {
	return nil
}
