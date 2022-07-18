package types

import (
	"errors"
	"fmt"
	"github.com/cruffinoni/rimworld-editor/helper"
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/xml/path"
	"github.com/cruffinoni/rimworld-editor/xml/saver"
	"github.com/cruffinoni/rimworld-editor/xml/types/iterator"
	"log"
	"reflect"
	"sort"
)

// Map is a map of K to V and require the XML file to have a "keys" and "values"
// element.

type MapComparable[T any] interface {
	Less(key reflect.Value, other T) bool
	Equal(key reflect.Value, other T) bool
	Great(key reflect.Value, other T) bool
}

// Map is a map of K to V.
// We don't restrict the type K to MapComparable[Map[K, V]] because K might be
// type of string, int or any primary type.
type Map[K, V comparable] struct {
	MapComparable[Map[K, V]]
	xml.Assigner
	iterator.MapIndexer[K, V]
	m          map[K]V
	sortedKeys []reflect.Value
}

func (m Map[K, V]) TransformToXML(buffer *saver.Buffer) error {
	if m.m == nil {
		buffer.WriteEmptyTag("keys", nil)
		buffer.WriteEmptyTag("values", nil)
		return nil
	}
	buffer.IncreaseDepth()
	buffer.WriteStringWithIndent("<keys>\n")
	buffer.IncreaseDepth()
	for k := range m.m {
		buffer.WriteStringWithIndent("<li>")
		buffer.IncreaseDepth()
		buffer.WriteString(fmt.Sprintf("%v", k))
		buffer.DecreaseDepth()
		buffer.WriteString("</li>\n")
	}
	buffer.DecreaseDepth()
	buffer.WriteStringWithIndent("</keys>\n")
	buffer.WriteStringWithIndent("<values>\n")
	buffer.IncreaseDepth()
	for _, v := range m.m {
		buffer.WriteStringWithIndent("<li>")
		buffer.IncreaseDepth()
		buffer.WriteString(fmt.Sprintf("%v", v))
		buffer.DecreaseDepth()
		buffer.WriteString("</li>\n")
	}
	buffer.DecreaseDepth()
	buffer.WriteStringWithIndent("</values>")
	buffer.DecreaseDepth()
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

func castDataFromKind[T comparable](kind reflect.Kind, d *xml.Data) T {
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
	log.Panicf("Map/castDataFromKind: cannot cast %v to %v", d, zero[T]())
	// Never reached
	return zero[T]()
}

func (m *Map[K, V]) Assign(e *xml.Element) error {
	m.m = make(map[K]V)
	n := e.Child
	if n == nil {
		return nil
	}
	keys := path.FindWithPath("keys>[...]", e)
	if len(keys) == 0 {
		return errors.New("Map/Assign: no key")
	}
	values := path.FindWithPath("values>[...]", e)
	if len(values) == 0 {
		return errors.New("Map/Assign: no value")
	}
	kKind := reflect.TypeOf(zero[K]()).Kind()
	vKind := reflect.TypeOf(zero[V]()).Kind()
	for i, key := range keys {
		if key.Data == nil || values[i].Data == nil {
			log.Panicf("Map/Assign: no data for %s or %s", key.StartElement.Name.Local, values[i].StartElement.Name.Local)
		}
		m.m[castDataFromKind[K](kKind, key.Data)] = castDataFromKind[V](vKind, values[i].Data)
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

func (m Map[K, V]) GetFromIndex(idx int) V {
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
