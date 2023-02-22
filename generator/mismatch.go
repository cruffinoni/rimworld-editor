package generator

import (
	"fmt"
	"log"
	"math"
	"reflect"
	"sort"
	"strings"
)

func isRelevantType(t1 any) bool {
	if t1 == nil {
		return false
	}
	if ct, ok := t1.(*CustomType); ok {
		if ct.Name == "Empty" && ct.Pkg == "*primary" {
			return false
		}
	}
	// Int64 is not relevant compared to other type making him erased by float64 in
	// almost all cases of comparison that implicates an int64.
	if reflect.TypeOf(t1).Kind() == reflect.Int64 {
		return false
	}
	return true
}

// fixCustomType compares two CustomType values and updates the values to ensure they are consistent.
// If either `a` or `b` is nil, it is updated to the non-nil value.
// The function checks the `type1` and `type2` fields and updates them if they are not the same.
// The function ensures that the relevant type is kept between both CustomType values.
func fixCustomType(a, b *CustomType) {
	if a == nil {
		a = b
	}
	if b == nil {
		b = a
	}
	if a.Type1 != b.Type1 {
		if isRelevantType(a.Type1) && !isRelevantType(b.Type1) {
			b.Type1 = a.Type1
			// If the main type is changed, update the package name, type name, and import path
			b.Name = a.Name
			b.Pkg = a.Pkg
			b.ImportPath = a.ImportPath
		} else if !isRelevantType(a.Type1) && isRelevantType(b.Type1) {
			a.Type1 = b.Type1
			// Perform the same update here
			a.Name = b.Name
			a.Pkg = b.Pkg
			a.ImportPath = b.ImportPath
		}
	}
	// Keep the relevant type between both `type2` fields.
	if a.Type2 != nil && a.Type2 != b.Type2 {
		if isRelevantType(a.Type2) && !isRelevantType(b.Type2) {
			b.Type2 = a.Type2
		} else if !isRelevantType(a.Type2) && isRelevantType(b.Type2) {
			a.Type2 = b.Type2
		}
	}
}

// printOrderedMembers print Members of s in an alphabetic order
func (s *StructInfo) printOrderedMembers() {
	if len(s.Members) == 0 {
		return
	}
	m := make([]string, 0, len(s.Members))
	for k := range s.Members {
		m = append(m, k)
	}
	sort.Strings(m)
	log.Printf("Struct %v", s.Name)
	for _, k := range m {
		log.Printf("'%s' > %T (%v)", k, s.Members[k].T, s.Members[k].T)
	}
	fmt.Printf("\n")
}

func fixTypeMismatch(a, b *member) {
	switch va := a.T.(type) {
	case *CustomType:
		if ctB, okB := b.T.(*CustomType); okB {
			fixCustomType(va, ctB)
		} else if _, okStruct := b.T.(*StructInfo); okStruct {
			//if isRelevantType(va) {
			//	structType.printOrderedMembers()
			//	log.Printf("fixTypeMismatch: double relevant type => %T (%v w/ t %T) & %T (%v w/ %d members)", va, va.Name, va.Type1, structType, structType.Name, len(structType.Order))
			//}
			a.T = b.T
		} else {
			b.T = a.T
		}
	case *StructInfo:
		if bStruct, okStruct := b.T.(*StructInfo); okStruct {
			fixMembers(va, bStruct)
		} else {
			b.T = a.T
		}
	case *FixedArray:
		if bFArr, okStruct := b.T.(*FixedArray); okStruct {
			if va.Size != bFArr.Size {
				va.Size = int(math.Max(float64(va.Size), float64(bFArr.Size)))
				bFArr.Size = va.Size
			}
			if !isSameType(va.PrimaryType, bFArr.PrimaryType) {
				log.Printf("'%v' | '%v'", a.Name, a.Attr)
				log.Printf("mismatch type in fixed array w/ %T (len %d) & %T (len %d)", va.PrimaryType, va.Size, bFArr.PrimaryType, bFArr.Size)
				log.Printf("%+v", va.PrimaryType)
				log.Printf("%+v", bFArr.PrimaryType)
				if ex := explainIsSameType(va.PrimaryType, bFArr.PrimaryType, &explainations{content: make([]string, 0)}); len(ex.content) > 0 {
					log.Printf("%s", strings.Join(ex.content, "\n"))
				}
				va.PrimaryType = bFArr.PrimaryType
			}
		} else {
			b.T = a.T
		}
	case reflect.Kind:
		bt, ok := b.T.(reflect.Kind)
		if !ok {
			// We have completely 2 different types with same name. Example of tag <name> which might be a structure representing the name, forename and surname
			// of a pawn but can be also a string for "feature" tag.
			if isRelevantType(b.T) {
				log.Printf("b type ('%v' | %+v) is not reflect.Kind type w/ %T", b.Name, b, b.T)
				addUniqueNumber(b.Name)
			} else {
				b.T = a.T
			}
		}
		if va == reflect.Int64 && bt == reflect.Float64 {
			a.T = reflect.Float64
		} else if va == reflect.Float64 && bt == reflect.Int64 {
			b.T = reflect.Float64
		}
	}
}

type explainations struct {
	content []string
}

func explainIsSameType(a, b any, e *explainations) *explainations {
	if a == nil || b == nil {
		if a == nil && b != nil {
			return &explainations{content: append([]string{fmt.Sprintf("a is nil, b is %T", b)}, e.content...)}
		} else if a != nil && b == nil {
			return &explainations{content: append([]string{fmt.Sprintf("b is nil, a is %T", a)}, e.content...)}
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
			if !isSameType(va.Type1, bType.Type1) {
				e.content = append(e.content, fmt.Sprintf("[CustomType] a type1 is diff from b (%T != %T)", va.Type1, bType.Type1))
			}
			if !isSameType(va.Type2, bType.Type2) {
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
			if !hasSameMembers(va, bType) {
				e.content = append(e.content, fmt.Sprintf("[StructInfo] a has not the same members of b (len: %d <> %d)", len(va.Members), len(bType.Members)))
				va.printOrderedMembers()
				bType.printOrderedMembers()
			}
		} else {
			e.content = append(e.content, "[StructInfo] b is not type StructInfo but a is")
			return e
		}
	case *FixedArray:
		if bFixArr, ok := b.(*FixedArray); ok {
			if !isSameType(bFixArr.PrimaryType, va.PrimaryType) {
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

func isSameType(a, b any) bool {
	if a == nil || b == nil {
		return a == b
	}
	log.Printf("checking %v & %v", a, b)
	switch va := a.(type) {
	case *CustomType:
		if bType, ok := b.(*CustomType); ok && bType != nil {
			return va.Name == bType.Name && va.Pkg == bType.Pkg &&
				isSameType(va.Type1, bType.Type1) && isSameType(va.Type2, bType.Type2)
		} else {
			return false
		}
	case *StructInfo:
		if bType, ok := b.(*StructInfo); ok && bType != nil {
			//log.Printf("So: %v/%v => %v", va.Name, bType.Name, hasSameMembers(va, bType))
			return va.Name == bType.Name && hasSameMembers(va, bType)
		} else {
			return false
		}
	case *FixedArray:
		if bFixArr, ok := b.(*FixedArray); ok && bFixArr != nil {
			//log.Printf("%T & %T && %d & %d", bFixArr.PrimaryType, va.PrimaryType, va.Size, bFixArr.Size)
			return isSameType(bFixArr.PrimaryType, va.PrimaryType) && va.Size == bFixArr.Size
		} else {
			return false
		}
	case *member:
		if bMember, ok := b.(*member); ok && bMember != nil {
			return va.Name == bMember.Name && isSameType(va.T, bMember.T)
		} else {
			return false
		}
	case reflect.Kind:
		if bKind, ok := b.(reflect.Kind); ok {
			return bKind == va
		} else {
			return false
		}
	default:
		return reflect.TypeOf(a) == reflect.TypeOf(b)
	}
}

func updateOrderedMembers(a *StructInfo) {
	for i := range a.Order {
		a.Order[i] = a.Members[a.Order[i].Name]
	}
}

func fixMembers(a, b *StructInfo) {
	for name, m := range a.Members {
		if _, ok := b.Members[name]; !ok {
			b.Members[name] = m
			b.Order = append(b.Order, m)
		}
	}
	for name, m := range b.Members {
		if _, ok := a.Members[name]; !ok {
			a.Members[name] = m
			a.Order = append(a.Order, m)
		}
	}
	/*if a.Name == "hediffs" {
		log.Printf(">>Hediffs:")
		a.printOrderedMembers()
		b.printOrderedMembers()
	}*/
	for i := range a.Members {
		if _, ok := b.Members[i]; !ok {
			a.printOrderedMembers()
			b.printOrderedMembers()
			log.Panicf("fixMembers: '%v' doesn't exist in b", i)
		}
		if !isSameType(a.Members[i].T, b.Members[i].T) {
			fixTypeMismatch(a.Members[i], b.Members[i])
			updateOrderedMembers(a)
			updateOrderedMembers(b)
		}
	}
}
