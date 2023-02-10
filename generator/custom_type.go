package generator

import (
	"log"
	"reflect"

	"github.com/cruffinoni/rimworld-editor/generator/paths"
	"github.com/cruffinoni/rimworld-editor/xml"
)

type CustomType struct {
	Name string
	Pkg  string
	// Type1 is the first type in the custom type
	// It can be anything from reflect.Kind, *CustomType, *xml.Element, etc.
	Type1      any
	Type2      any
	ImportPath string
}

func createEmptyType() any {
	return &CustomType{
		Name:       "Empty",
		Pkg:        "*primary",
		Type1:      nil,
		ImportPath: paths.PrimaryTypesPath,
	}
}

func createCustomSlice(e *xml.Element, flag uint) any {
	return &CustomType{
		Name:       "Slice",
		Pkg:        "*types",
		Type1:      createSubtype(e, flag, getTypeFromArray(e)),
		ImportPath: paths.CustomTypesPath,
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
		if ct.Name == "Empty" {
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
			Name:       "Map",
			Pkg:        "types",
			Type1:      reflect.String,
			Type2:      reflect.String,
			ImportPath: paths.CustomTypesPath,
		}
	}
	return &CustomType{
		Name:       "Map",
		Pkg:        "types",
		Type1:      k,
		Type2:      v,
		ImportPath: paths.CustomTypesPath,
	}
}
