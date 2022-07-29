package generator

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"log"
	"reflect"
)

type CustomType struct {
	name string
	pkg  string
	// type1 is the first type in the custom type
	// It can be anything from reflect.Kind, *CustomType, *xml.Element, etc.
	type1 any
	type2 any
}

func createCustomSlice(e *xml.Element, flag uint) any {
	c := e
	var t any
	t = getTypeFromArray(c)
	if t == reflect.Struct {
		t = createStructure(c, flag|skipChild)
		//if c.Child.GetName() == "li" {
		//	t = createCustomSlice(c.Child, flag)
		//} else {
		//	t = createStructure(c, flag)
		//}
	} else if t == reflect.Invalid {
		t = e
	}
	return &CustomType{
		name:  "Slice",
		pkg:   "types",
		type1: t,
	}
}

func createCustomTypeForMap(e *xml.Element, flag uint) any {
	if e.Child == nil {
		log.Panic("generate.createCustomTypeForMap: missing child")
	}

	var (
		c = e.Child
		k = determineTypeFromData(c, flag)
		v any
	)
	c = c.Next
	v = determineTypeFromData(c, flag)
	// By default, maps are strings to strings
	if k == reflect.Invalid || v == reflect.Invalid {
		return &CustomType{
			name:  "Map",
			pkg:   "types",
			type1: reflect.String,
			type2: reflect.String,
		}
	}
	return &CustomType{
		name:  "Map",
		pkg:   "types",
		type1: k,
		type2: v,
	}
}
