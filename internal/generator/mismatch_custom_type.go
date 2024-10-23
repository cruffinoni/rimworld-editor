package generator

import (
	"errors"
	"reflect"

	"github.com/cruffinoni/printer"
)

var ErrUnsolvableMismatch = errors.New("unsolvable mismatch")

// fixCustomType main purpose is to reconcile the types of two CustomType values.
// If either of them is nil, it updates the nil value to match the other.
// If both are non-nil, it ensures that the Type1 and Type2 fields of both CustomType values are consistent by using helper functions.
func fixCustomType(a, b *CustomType) error {
	if a == nil {
		a = b
		return nil
	}
	if b == nil {
		b = a
		return nil
	}

	if err := reconcileTypes(&a.Type1, &b.Type1, a, b); err != nil {
		if !errors.Is(err, ErrUnsolvableMismatch) {
			return err
		}
		printer.Debugf("(t1) data a: %v -> %T (%v)", a.Name, a.Type1, a.Type1)
		printer.Debugf("(t1) data b: %v -> %T (%v)", b.Name, b.Type1, b.Type1)
		printer.Debugf("Multiple type detected")

		// a: *CustomType[*multiple.Type] / *CustomType[*xml.Element]
		// b: *CustomType[*multiple.Type] / *CustomType[*xml.Element]
		a.Type1 = createXMLElementType()
		b.Type1 = createXMLElementType()
		return nil
	}
	if err := reconcileTypes(&a.Type2, &b.Type2, a, b); err != nil {
		printer.Debugf("(t2) data a: %v -> %T (%+v)", a.Name, a.Type2, a.Type2)
		printer.Debugf("(t2) data b: %v -> %T (%+v)", b.Name, b.Type2, b.Type2)
		printer.Debugf("(t2) can't decide on which type work on between %T & %T", a.Type2, b.Type2)
		panic("Multiple type detected for type2")
		return err
	}
	return nil
}

func shouldUpdateType(fromType, toType reflect.Type, fromValue, toValue any) bool {
	return (isRelevantType(fromValue) && !isRelevantType(toValue)) ||
		(fromValue != nil && fromType.Kind() == reflect.Float64 && toType.Kind() == reflect.Int64)
}

func updateCustomType(dest, src *CustomType) {
	dest.Type1 = src.Type1
	dest.Name = src.Name
	dest.Pkg = src.Pkg
	dest.ImportPath = src.ImportPath
}

func reconcileTypes(aType, bType *any, a, b *CustomType) error {
	if *aType == nil && *bType == nil {
		return nil
	}
	aRefType := reflect.TypeOf(*aType)
	bRefType := reflect.TypeOf(*bType)

	if !IsSameType(*aType, *bType, 0) {
		if shouldUpdateType(aRefType, bRefType, *aType, *bType) {
			updateCustomType(b, a)
		} else if shouldUpdateType(bRefType, aRefType, *bType, *aType) {
			updateCustomType(a, b)
		} else {
			return handleMismatch(aType, bType, a, b)
		}
	}
	return nil
}

func handleMismatch(aType, bType *any, a, b *CustomType) error {
	switch va := (*aType).(type) {
	case reflect.Kind:
		if vb, ok := (*bType).(reflect.Kind); ok {
			if isIntFloatMismatch(va, vb) {
				*aType = reflect.Float64
			} else if isIntFloatMismatch(vb, va) {
				*bType = reflect.Float64
			}
			if isStringMapMismatch(va, a, vb) {
				*aType = b.Type2
			}
		} else {
			return ErrUnsolvableMismatch
		}
	case *StructInfo:
		if vb, ok := (*bType).(*StructInfo); ok {
			FixMembers(va, vb)
		} else {
			return ErrUnsolvableMismatch
		}
	default:
		return ErrUnsolvableMismatch
	}
	return nil
}

func isIntFloatMismatch(a, b reflect.Kind) bool {
	return a == reflect.Int64 && b == reflect.Float64
}

func isStringMapMismatch(kind reflect.Kind, a *CustomType, bKind reflect.Kind) bool {
	return kind == reflect.String && a.Name == "Map" && bKind != reflect.String
}
