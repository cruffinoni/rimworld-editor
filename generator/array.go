package generator

import (
	"reflect"

	"github.com/cruffinoni/rimworld-editor/xml"
)

type FixedArray struct {
	Size        int
	PrimaryType any
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
	return t
}

type offset struct {
	el   *xml.Element
	size int
}

func createFixedArray(e *xml.Element, flag uint, o *offset) any {
	f := &FixedArray{
		PrimaryType: createSubtype(e, flag, getTypeFromArray(e)),
		Size:        0,
	}
	if o == nil {
		o = &offset{el: e.Child}
	} else {
		f.Size = o.size
	}
	k := o.el
	for k != nil {
		f.Size++
		k = k.Next
	}
	return f
}

func (a *FixedArray) ValidateField(_ string) {
}

func (a *FixedArray) IsValidField(_ string) bool {
	return true
}
