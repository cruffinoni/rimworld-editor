package generator

import (
	"log"
	"reflect"

	"github.com/cruffinoni/rimworld-editor/xml"
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
		pkg:        "*primary",
		type1:      nil,
		importPath: primaryTypesPath,
	}
}

func createCustomSlice(e *xml.Element, flag uint) any {
	return &CustomType{
		name:       "Slice",
		pkg:        "*types",
		type1:      createSubtype(e, flag, getTypeFromArray(e)),
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
		k = determineTypeFromData(c, flag|ignoreSlice|forceFullCheck)
		v any
	)
	if ct, ok := k.(*CustomType); ok {
		// primary.Empty does not implement comparable
		// Might be deleted when the types.Map type implements any as Key and not comparable anymore
		if ct.name == "Empty" {
			k = reflect.String
		}
	}
	//log.Printf("Key type: %T", k)
	c = c.Next
	//log.Printf("Determining value type from '%v'", c.XMLPath())
	v = determineTypeFromData(c, flag|ignoreSlice|forceFullCheck)
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
