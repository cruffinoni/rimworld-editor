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
	type1      any
	type2      any
	importPath string
}

func createEmptyType() any {
	return &CustomType{
		name:       "Empty",
		pkg:        "primary",
		type1:      reflect.Invalid,
		importPath: primaryTypesPath,
	}
}

func createCustomSlice(e *xml.Element, flag uint) any {
	var t any
	t = getTypeFromArray(e)
	//log.Printf("Type of array: %v (%v) & %p", t, e.XMLPath(), e)
	switch t {
	case reflect.Invalid:
		// With an invalid type and no data, we can assume that the slice is empty
		if e.Data == nil {
			t = createEmptyType()
		} else {
			log.Printf("invalid type: %v => '%v' (data: '%v')", t, e.XMLPath(), e.Data)
			t = e
		}
	case reflect.Slice:
		//log.Printf("Creating custom slice from %s", e.Child.XMLPath())
		t = createCustomSlice(e.Child, flag)
	case reflect.Struct:
		//log.Printf("[Struct] creating a structure from %s", e.XMLPath())
		t = createStructure(e, flag|skipChild)
	}
	return &CustomType{
		name:       "Slice",
		pkg:        "types",
		type1:      t,
		importPath: customTypesPath,
	}
}

func createCustomTypeForMap(e *xml.Element, flag uint) any {
	if e.Child == nil {
		log.Panic("generate.createCustomTypeForMap: missing child")
	}

	//log.Printf("Determining key type from %s", e.Child.XMLPath())
	var (
		c = e.Child
		k = determineTypeFromData(c, flag|ignoreSlice)
		v any
	)
	if ct, ok := k.(*CustomType); ok {
		// primary.Empty does not implement comparable
		// TODO: Might be deleted when the types.Map type implements any as Key and not comparable anymore
		if ct.name == "Empty" {
			k = reflect.String
		}
	}
	//log.Printf("Key type: %T", k)
	c = c.Next
	//log.Printf("Determining value type from '%v'", c.XMLPath())
	v = determineTypeFromData(c, flag|ignoreSlice)
	//log.Printf("Value type: %T", v)
	// By default, maps are strings to strings
	if k == reflect.Invalid || v == reflect.Invalid {
		return &CustomType{
			name:       "Map",
			pkg:        "types",
			type1:      reflect.String,
			type2:      reflect.String,
			importPath: customTypesPath,
		}
	}
	return &CustomType{
		name:       "Map",
		pkg:        "types",
		type1:      k,
		type2:      v,
		importPath: customTypesPath,
	}
}
