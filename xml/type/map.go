package _type

import (
	"errors"
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/path"
	"github.com/cruffinoni/rimworld-editor/xml/type/iterator"
	"log"
	"reflect"
)

// Map is a map of K to V and require the XML file to have a "keys" and "values"
// element.
type Map[K, V comparable] struct {
	xml.Assigner
	iterator.SliceIndexer[V]
	m map[K]V
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

func (m *Map[K, V]) GetFromKey(idx int) K {
	if m.m == nil {
		return zero[K]()
	}
	if idx < 0 || idx >= len(m.m) {
		log.Panic("Map/GetFromIndex: index out of range")
		return zero[K]()
	}
	i := 0
	for k := range m.m {
		if i == idx {
			return k
		}
		i++
	}
	log.Panicf("Map/GetFromKey: index %d not found", idx)
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
