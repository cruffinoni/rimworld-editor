package collection

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"sort"
	"unicode"

	"github.com/cruffinoni/rimworld-editor/internal/xml/support"

	"github.com/cruffinoni/rimworld-editor/internal/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/internal/xml/binder"
	"github.com/cruffinoni/rimworld-editor/internal/xml/domain"
	"github.com/cruffinoni/rimworld-editor/internal/xml/encoder"
	"github.com/cruffinoni/rimworld-editor/internal/xml/encoder/reflection"
	"github.com/cruffinoni/rimworld-editor/internal/xml/interfaces"
	"github.com/cruffinoni/rimworld-editor/internal/xml/iterator"
	"github.com/cruffinoni/rimworld-editor/internal/xml/query"
	"github.com/cruffinoni/rimworld-editor/internal/xml/scalar"
	"github.com/cruffinoni/rimworld-editor/pkg/logging"
)

type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

func (p *Pair[K, V]) Equal(rhs *Pair[K, V]) bool {
	//return p.Key == rhs.Key && p.Val == rhs.Val
	return true
}

type MapComparable[T any] interface {
	Less(key reflect.Value, other T) bool
	Equal(key reflect.Value, other T) bool
	Great(key reflect.Value, other T) bool
}

// Map is a map of K to V.
// We don't restrict the type K to MapComparable[Map[K, V]] because K might be
// type of string, int or multiple primary type.
type Map[K comparable, V any] struct {
	MapComparable[Map[K, V]]
	interfaces.Assigner
	iterator.MapIndexer[K, V]
	m          map[K]V
	sortedKeys []reflect.Value

	tag  string
	attr attributes.Attributes

	logger logging.Logger
}

func castToInterface[T any](val any) (T, bool) {
	if v, ok := val.(T); ok {
		return v, true
	}
	return *new(T), false
}

func (m *Map[K, V]) TransformToXML(b *encoder.Buffer) error {
	b.OpenTag(m.tag, m.attr)
	b.WriteString("\n")
	defer func() {
		b.WriteString("\n")
		b.CloseTagWithIndent(m.tag)
		b.WriteString("\n")
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
		if err := reflection.Save(k, b, ""); err != nil {
			return err
		}
		if unicode.IsSpace(rune(b.Bytes()[b.Len()-1])) {
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
		if err := reflection.Save(v, b, ""); err != nil {
			return err
		}
		if unicode.IsSpace(rune(b.Bytes()[b.Len()-1])) {
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

func castDataFromKind[T any](kind reflect.Kind, d *domain.Data) T {
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

func (m *Map[K, V]) Assign(e *domain.Element) error {
	if m.Capacity() == 0 {
		m.m = make(map[K]V)
		if e.Parent != nil {
			m.tag = e.Parent.GetName()
			m.attr = e.Parent.GetAttributes()
		} else {
			return fmt.Errorf("map.Assign: map's parent is nil")
		}
	}
	if e.Child == nil {
		return nil
	}
	//printer.Debugf("Tag: %v", m.tag)
	keys := query.FindWithPath("keys>[...]", e, m.logger)
	if len(keys) == 0 {
		return errors.New("Map/Assign: no key")
	}
	//printer.Debugf("e=%v", e.GetName())
	values := query.FindWithPath("values>[...]", e, m.logger)
	if len(values) == 0 {
		return errors.New("Map/Assign: no value")
	}
	if len(keys) != len(values) {
		return errors.New("Map/Assign: keys length differs from values length")
	}
	//printer.Debugf("Keys: %v, Val: %v", keys[0].XMLPath(), values[0].XMLPath())
	//printer.Debugf("Keys: %+v, Val: %+v", keys[0].last, values[0].last)
	kKind := reflect.TypeOf(zero[K]()).Kind()
	vKind := reflect.TypeOf(zero[V]()).Kind()
	_, isEmpty := any(zero[V]()).(*scalar.Empty)
	//printer.Debugf("%T is empty ? %v", zero[V](), isEmpty)
	for i, key := range keys {
		if key.Data == nil {
			log.Panicf("Map/Assign: no data for key %s", key.GetName())
		}

		_, okAssigner := any(zero[V]()).(interfaces.Assigner)
		// Special case with array because we need to check if the type implements domain.Assigner interface
		// and not the array itself
		if vKind == reflect.Array {
			okAssigner = reflect.TypeOf(zero[V]()).Elem().AssignableTo(reflect.TypeOf((*interfaces.Assigner)(nil)).Elem())
		}
		// This might be a custom type that implements domain.Assigner interface
		if okAssigner && !isEmpty {
			if values[i].Child == nil {
				m.logger.WithFields(logging.Fields{
					"path":  e.XMLPath(),
					"index": i,
				}).Debug("Map.Assign: no child")
			}
			var (
				subValue    = new(V)
				subValueVal = reflect.ValueOf(subValue)
			)
			//printer.Debugf("! > %v & %T + '%v' & POSSIBLE? %v", subValue, subValue, subValueVal.Kind(), subValueVal.Elem().CanAddr())
			subValueVal = subValueVal.Elem()
			if subValueVal.Kind() == reflect.Ptr {
				// Initialize the pointer
				subValueVal.Set(reflect.New(subValueVal.Type().Elem()))
				if err := binder.Element(m.logger, values[i].Child, subValueVal.Interface()); err != nil {
					return err
				}
			} else if subValueVal.Kind() == reflect.Array { // A pointer to an array is not assignable to an array
				for j := 0; j < subValueVal.Len(); j++ {
					if subValueVal.Index(j).Kind() == reflect.Ptr {
						subValueVal.Index(j).Set(reflect.New(subValueVal.Index(j).Type().Elem()))
						if err := binder.Element(m.logger, values[i].Child, subValueVal.Index(j).Interface()); err != nil {
							return err
						}
					} else {
						panic("Map/Assign: array element must be a pointer") // TODO: Handle this case
					}
				}
			} else {
				panic("Map/Assign: value must be a pointer") // TODO: Handle this case
			}
			m.m[castDataFromKind[K](kKind, key.Data)] = subValueVal.Interface().(V)
			//printer.Debugf("!!=> %v > %v", castDataFromKind[K](kKind, key.last), m.m[castDataFromKind[K](kKind, key.last)])
		} else if values[i].Data == nil || isEmpty {
			// There is a key with no data
			m.m[castDataFromKind[K](kKind, key.Data)] = zero[V]()
		} else if _, isElement := any(zero[V]()).(*domain.Element); isElement {
			// Special if V is a domain.Element because we pass a pointer to the data for castDataFromKind
			// so, we don't use this function but assign directly to the map
			m.m[castDataFromKind[K](kKind, key.Data)] = castTemplate[V](values[i])
		} else {
			m.m[castDataFromKind[K](kKind, key.Data)] = castDataFromKind[V](vKind, values[i].Data)
		}
		//printer.Debugf("=> %v > %v", castDataFromKind[K](kKind, key.last), m.m[castDataFromKind[K](kKind, key.last)])
	}

	v := reflect.ValueOf(m.m)
	k := reflect.ValueOf(zero[K]()).Kind()
	m.sortedKeys = v.MapKeys()
	// Primary type implements natively operator<
	if support.IsReflectPrimaryType(k) {
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
		// Custom type must implement operator< with the function "func (MapComparable) Less(key reflect.Value, other T) bool"
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
		log.Panic("Map/At: index out of range")
		return zero[V]()
	}
	i := 0
	for _, v := range m.m {
		if i == idx {
			return v
		}
		i++
	}
	log.Panicf("Map/At: index %d not found", idx)
	return zero[V]()
}

func (m *Map[K, V]) GetKeyFromIndex(idx int) K {
	if m.m == nil {
		return zero[K]()
	}
	if idx < 0 || idx >= len(m.m) {
		log.Panic("Map/At: index out of range")
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

func (m *Map[K, V]) SetLogger(logger logging.Logger) {
	m.logger = logger
}

func (m *Map[K, V]) ValidateField(_ string) {
}

func (m *Map[K, V]) IsValidField(_ string) bool {
	return true
}

func (m *Map[K, V]) CountValidatedField() int {
	return m.Capacity()
}
