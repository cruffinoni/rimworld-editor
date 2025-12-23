package codegen

import (
	"reflect"

	"github.com/cruffinoni/rimworld-editor/internal/xml/domain"
	"github.com/cruffinoni/rimworld-editor/pkg/logging"
)

type FixedArray struct {
	Size        int
	PrimaryType any
}

func createSubtype(logger logging.Logger, e *domain.Element, flag uint, t any) any {
	switch t {
	case Complex:
		return e
	case reflect.Invalid:
		// With an invalid type and no data, we can assume that the slice is empty
		if e.Data == nil {
			return createEmptyType()
		}

		return e
	case reflect.Slice:
		return createCustomSlice(logger, e.Child, flag)
	case reflect.Array:
		return createFixedArray(logger, e.Child, flag, nil)
	case reflect.Struct:
		return createStructure(logger, e, flag|forceFullCheck)
	default:
		return t
	}
}

type offset struct {
	el   *domain.Element
	size int
}

func createFixedArray(logger logging.Logger, e *domain.Element, flag uint, o *offset) any {
	f := &FixedArray{
		PrimaryType: createSubtype(logger, e, flag, getTypeFromArrayOrSlice(logger, e)),
		Size:        1, // Minimum size is 1
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
