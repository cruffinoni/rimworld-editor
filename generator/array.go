package generator

import (
	"log"
	"reflect"

	"github.com/cruffinoni/rimworld-editor/xml"
)

type FixedArray struct {
	size        int
	primaryType any
}

func createSubtype(e *xml.Element, flag uint, t any) any {
	switch t {
	case reflect.Invalid:
		// With an invalid type and no data, we can assume that the slice is empty
		if e.Data == nil {
			return createEmptyType()
		} else {
			return e
		}
	case reflect.Slice:
		return createCustomSlice(e.Child, flag)
	case reflect.Array:
		return createFixedArray(e.Child, flag, nil)
	case reflect.Struct:
		return createStructure(e, flag|forceFullCheck)
	}
	return nil
}

type offset struct {
	el   *xml.Element
	size int
}

func createFixedArray(e *xml.Element, flag uint, o *offset) any {
	f := &FixedArray{
		primaryType: createSubtype(e, flag, getTypeFromArray(e)),
		size:        o.size,
	}
	if o == nil {
		o = &offset{el: e}
	}
	k := o.el
	for k != nil {
		f.size++
		k = k.Next
	}
	log.Printf("Fixed array of size %d w/ pt %T at %s", f.size, f.primaryType, e.XMLPath())
	return f
}
