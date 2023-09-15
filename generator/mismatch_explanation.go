package generator

import (
	"fmt"
	"reflect"
)

type explanations struct {
	content []string
}

func explainIsSameType(a, b any, e *explanations) *explanations {
	if a == nil || b == nil {
		if a == nil && b != nil {
			return &explanations{content: append([]string{fmt.Sprintf("a is nil, b is %T", b)}, e.content...)}
		} else if a != nil && b == nil {
			return &explanations{content: append([]string{fmt.Sprintf("b is nil, a is %T", a)}, e.content...)}
		}
	}
	switch va := a.(type) {
	case *CustomType:
		if bType, ok := b.(*CustomType); ok {
			if va.Name != bType.Name {
				e.content = append(e.content, fmt.Sprintf("[CustomType] a name is diff from b ('%v' != '%v')", va.Name, bType.Name))
			}
			if va.Pkg != bType.Pkg {
				e.content = append(e.content, fmt.Sprintf("[CustomType] a pkg is diff from b ('%v' != '%v')", va.Pkg, bType.Pkg))
			}
			if !IsSameType(bType.Type1, va.Type1, 0) {
				e.content = append(e.content, fmt.Sprintf("[CustomType] a type1 is diff from b (%T != %T)", va.Type1, bType.Type1))
			}
			if !IsSameType(bType.Type2, va.Type2, 0) {
				e.content = append(e.content, fmt.Sprintf("[CustomType] a type2 is diff from b (%T != %T)", va.Type2, bType.Type2))
			}
			return e
		} else {
			e.content = append(e.content, "[CustomType] b is not type CustomType but a is")
			return e
		}
	case *StructInfo:
		if bType, ok := b.(*StructInfo); ok {
			if va.Name != bType.Name {
				e.content = append(e.content, fmt.Sprintf("[StructInfo] a name is diff from b ('%v' != '%v')", va.Name, bType.Name))
			}
			if !hasSameMembers(bType, va, 0) {
				e.content = append(e.content, fmt.Sprintf("[StructInfo] a has not the same members of b (len: %d <> %d)", len(va.Members), len(bType.Members)))
				va.PrintOrderedMembers()
				bType.PrintOrderedMembers()
			}
		} else {
			e.content = append(e.content, "[StructInfo] b is not type StructInfo but a is")
			return e
		}
	case *FixedArray:
		if bFixArr, ok := b.(*FixedArray); ok {
			if !IsSameType(va.PrimaryType, bFixArr.PrimaryType, 0) {
				e.content = append(e.content, fmt.Sprintf("[FixedArray] a is not same type w/ b (%T != %T)", bFixArr.PrimaryType, va.PrimaryType))
				return e
			}
			if va.Size != bFixArr.Size {
				e.content = append(e.content, fmt.Sprintf("[FixedArray] a size is not same as b (%d != %d)", va.Size, bFixArr.Size))
				return e
			}
		} else {
			e.content = append(e.content, "[FixedArray] b is not type FixedArray")
			return e
		}
	case reflect.Kind:
		if bKind, ok := b.(reflect.Kind); ok {
			if bKind != va {
				e.content = append(e.content, fmt.Sprintf("[Kind] a is not same as b (%s != %s)", bKind, va))
				return e
			}
		} else {
			e.content = append(e.content, "[Kind] b is not type reflect.Kind but a is")
			return e
		}
	default:
		if reflect.TypeOf(a) != reflect.TypeOf(b) {
			e.content = append(e.content, fmt.Sprintf("[Type] a is not same as b (%s!= %s)", reflect.TypeOf(a), reflect.TypeOf(b)))
			return e
		}
	}
	return e
}
