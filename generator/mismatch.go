package generator

import (
	"errors"
	"log"
	"math"
	"reflect"
	"sort"
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
	return true
}

// PrintOrderedMembers print Members of s in an alphabetic order
func (s *StructInfo) PrintOrderedMembers() {
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
		if ct, ok := s.Members[k].T.(*CustomType); ok {
			log.Printf("  > Custom types: %T & %T", ct.Type1, ct.Type2)
		}
	}
}

func containsParentMember(s *StructInfo, a *Member) bool {
	for n := range s.Members {
		if n == a.Name {
			return true
		}
	}
	return false
}

func consolidateSubStructures(newMember, existingStruct *StructInfo) error {
	if _, ok := existingStruct.Members[newMember.Name]; !ok {
		m := &Member{
			T:    newMember,
			Attr: nil,
			Name: newMember.Name,
		}
		existingStruct.Members[newMember.Name] = m
		existingStruct.Order = append(existingStruct.Order, m)
	}
	return fixTypeMismatch(&Member{
		T:    existingStruct.Members[newMember.Name],
		Attr: nil,
		Name: newMember.Name,
	}, &Member{
		T:    newMember,
		Attr: nil,
		Name: newMember.Name,
	})
}

func fixTypeMismatch(a, b *Member) error {
	//log.Printf("Types mismatch: %v (%T) & %v (%T)", getTypeName(a.T), a.T, getTypeName(b.T), b.T)
	switch va := a.T.(type) {
	// a: *CustomType
	// b: ?
	case *CustomType:
		// The type of "a" is primary empty, so it's a non-relevant type
		// a: *CustomType[*primary.Empty]
		// b: ?
		if IsEmptyType(va) {
			a.T = b.T
			return nil
		}
		switch vb := b.T.(type) {
		// a: *CustomType[?]
		// b: *CustomType[?]
		case *CustomType:
			//printer.Printf("Fixing 2 custom types: %+v & %+v", va.Type1, vb.Type1)
			if IsEmptyType(vb) {
				b.T = a.T
				return nil
			}
			return fixCustomType(&va, &vb)

		// a: *CustomType[?]
		// b: *StructInfo
		case *StructInfo:
			// This case happens when "a.T" is a slice and "b.T" is only a structure and not a slice of structure (happens when there is only 1 element, and it's detected as a structure).
			// Consolidate both structure to make sure the types match
			// a: CustomType[*StructInfo]
			// b: *StructInfo
			if sliceStructType, ok := va.Type1.(*StructInfo); ok {
				return consolidateSubStructures(sliceStructType, vb)
			} else {
				log.Printf("sub type A not handled: %T (%+v)", va.Type1, va)
				a.T = b.T
			}

		// a: *CustomType[?]
		// b: *FixedArray
		case *FixedArray:
			if IsEmptyType(va) {
				a.T = b.T
				return nil
			}
			switch vbPrimaryType := vb.PrimaryType.(type) {

			// a: *CustomType[?]
			// b: *FixedArray[*StructInfo]
			case *StructInfo:

				// a: *CustomType[*StructInfo]
				// b: *FixedArray[*StructInfo]
				sliceStructType, ok := va.Type1.(*StructInfo)
				if !ok {
					log.Fatalf("SliceStruct in FixedArray detected but the primary type is not supported: %T (%v)", va.Type1, va.Type1)
				}
				//log.Printf("Arr: %v & %v", vbPrimaryType.Name, sliceStructType.Name)
				return fixTypeMismatch(&Member{
					T:    vbPrimaryType,
					Attr: nil,
					Name: vbPrimaryType.Name,
				}, &Member{
					T:    sliceStructType,
					Attr: nil,
					Name: sliceStructType.Name,
				})

			// a: *CustomType[?]
			// b: *FixedArray[*CustomType[?]]
			case *CustomType:
				if IsEmptyType(vbPrimaryType) {
					vb.PrimaryType = va.Type1
					// FixedArray take priority over CustomType because slice are detected by default but array are a specialized slice
					a.T = b.T
					return nil
				}
				return fixTypeMismatch(&Member{
					T:    vbPrimaryType,
					Attr: nil,
					Name: vbPrimaryType.Name,
				}, &Member{
					T:    va.Type1,
					Attr: nil,
					Name: va.Name,
				})
			}

		// a: ?
		// b: ?
		default:
			log.Panicf("type B not handled: %+v (%T) | %+v (%T)", a.T, a.T, b.T, b.T)
		}
	case *StructInfo:
		if bStruct, okStruct := b.T.(*StructInfo); okStruct {
			if !containsParentMember(va, a) {
				FixMembers(va, bStruct)
			}
			return nil
		} else {
			// If "b" is not a "StructInfo", the relevant type might be on the other side.
			// To avoid code duplication, we just swap those 2 and run through the actual code
			switch b.T.(type) {
			case *FixedArray, *CustomType, reflect.Kind:
				if err := fixTypeMismatch(b, a); err != nil {
					return err
				}
				a.T = b.T
			}
		}
	case *FixedArray:
		if bFArr, okStruct := b.T.(*FixedArray); okStruct {
			if va.Size != bFArr.Size {
				va.Size = int(math.Max(float64(va.Size), float64(bFArr.Size)))
				bFArr.Size = va.Size
			}
			if !IsSameType(bFArr.PrimaryType, va.PrimaryType, 0) {
				return fixTypeMismatch(&Member{
					T:    bFArr.PrimaryType,
					Attr: nil,
					Name: "",
				}, &Member{
					T:    va.PrimaryType,
					Attr: nil,
					Name: "",
				})
				//log.Printf("mismatch type in fixed array w/ %T (len %d) & %T (len %d)", va.PrimaryType, va.Size, bFArr.PrimaryType, bFArr.Size)
				//return nil
				//switch b.T.(type) {
				//case *FixedArray, *CustomType, reflect.Kind:
				//	if err := fixTypeMismatch(b, a); err != nil {
				//		return err
				//	}
				//	a.T = b.T
				//}
			}
		} else {
			log.Printf("Mismatch between FixedArray (%v) & %T (%v)", getTypeName(a), b.T, getTypeName(b.T))
			switch b.T.(type) {
			case *FixedArray, *CustomType, reflect.Kind:
				if err := fixTypeMismatch(b, a); err != nil {
					return err
				}
				a.T = b.T
			}
		}
	case reflect.Kind:
		bt, ok := b.T.(reflect.Kind)
		if !ok {
			// We have completely 2 different types with same name. Example of tag <name> which might be a structure representing the name, forename and surname
			// of a pawn but can be also a string for "feature" tag.
			if isRelevantType(b.T) {
				b.Name = addUniqueNumber(b.Name)
			} else {
				b.T = a.T
			}
		}

		// Previous type is an int64 and the next type might an int overflow considered as a string
		if va == reflect.Int64 && bt == reflect.String {
			log.Printf("Previous type is an int64 and the next type might an int overflow considered as a string")
			a.T = reflect.String
		} else if va == reflect.Int64 && bt == reflect.Float64 {
			a.T = reflect.Float64
		} else if va == reflect.Float64 && bt == reflect.Int64 {
			b.T = reflect.Float64
		}
	}
	return nil
}

const MaxDepth = 50

func hasSameMembers(a, b *StructInfo, depth uint32) bool {
	//log.Printf("A: %v & b %v => %d", a.Name, b.Name, depth)
	if depth > MaxDepth {
		panic("max depth reached")
		return false
	}
	if len(a.Members) != len(b.Members) {
		return false
	}
	for i := range a.Members {
		if _, ok := b.Members[i]; !ok {
			return false
		}
		if getTypeName(a.Name) == getTypeName(a.Members[i].Name) ||
			getTypeName(a.Name) == getTypeName(b.Members[i].Name) {
			return false
		}
		if !IsSameType(a.Members[i], b.Members[i], depth+1) {
			return false
		}
	}
	return true
}

func getTypeName(a any) string {
	switch va := a.(type) {
	case *CustomType:
		return getTypeName(va.Type1)
	case *StructInfo:
		return va.Name
	case *FixedArray:
		return getTypeName(va.PrimaryType)
	case *Member:
		return va.Name
	case reflect.Kind:
		return va.String()
	case nil:
		return "null"
	default:
		return ""
	}
}

// IsSameType compares the types of two objects a and b and returns true if they are the same type.
// depth is used to prevent infinite recursion in case of nested types.
func IsSameType(a, b any, depth uint32) bool {
	// If depth has exceeded MaxDepth, log a message and return false to stop further recursion.
	if depth > MaxDepth {
		log.Printf("Max depth reached: %v", depth)
		return false
	}
	// If either a or b is nil, compare their equality and return the result.
	if a == nil || b == nil {
		return a == b
	}
	// Check the type of "a" using type switch.
	switch va := a.(type) {
	case *CustomType:
		// If "a" is a pointer to CustomType, compare its properties with b.
		if bType, ok := b.(*CustomType); ok && bType != nil {
			// If the type's name is the same as the parent's name, return false to avoid infinite recursion.
			if getTypeName(va.Type1) == "" || getTypeName(bType.Type1) == "" {
				log.Panicf("type name is empty: %T > %+v | %T > %+v", va.Type1, va, bType.Type1, bType)
			}
			//if getTypeName(va.Type1) == getTypeName(bType.Type1) || getTypeName(va.Type2) == getTypeName(bType.Type2) {
			//	return false
			//}
			// Compare the properties of CustomType a with b recursively using IsSameType.
			return va.Name == bType.Name && va.Pkg == bType.Pkg &&
				IsSameType(bType.Type1, va.Type1, depth+1) && IsSameType(bType.Type2, va.Type2, depth+1)
		} else {
			return false
		}
	case *StructInfo:
		// If "a" is a pointer to StructInfo, compare its properties with b.
		if bType, ok := b.(*StructInfo); ok && bType != nil {
			// Compare the properties of StructInfo a with b recursively using hasSameMembers.
			return va.Name == bType.Name && hasSameMembers(bType, va, depth+1)
		} else {
			return false
		}
	case *FixedArray:
		// If "a" is a pointer to FixedArray, compare its properties with b.
		if bFixArr, ok := b.(*FixedArray); ok && bFixArr != nil {
			// Compare the properties of FixedArray a with b recursively using IsSameType.
			return IsSameType(va.PrimaryType, bFixArr.PrimaryType, depth+1) && va.Size == bFixArr.Size
		} else {
			return false
		}
	case *Member:
		// If "a" is a pointer to Member, compare its properties with b.
		if bMember, ok := b.(*Member); ok && bMember != nil {
			// Compare the properties of Member a with b recursively using IsSameType.
			return va.Name == bMember.Name && IsSameType(bMember.T, va.T, depth+1)
		} else {
			return false
		}
	case reflect.Kind:
		// If "a" is a pointer to Member, compare its properties with b.
		if bKind, ok := b.(reflect.Kind); ok {
			return bKind == va
		} else {
			return false
		}
	default:
		// For all other types, compare their types using reflect.TypeOf.
		return reflect.TypeOf(a) == reflect.TypeOf(b)
	}
}

func updateOrderedMembers(a *StructInfo) {
	l := len(a.Order)
	for i := 0; i < l; i++ {
		if _, ok := a.Members[a.Order[i].Name]; !ok {
			if i+1 >= len(a.Order) {
				a.Order = a.Order[:i]
			} else {
				a.Order = append(a.Order[:i], a.Order[i+1:]...)
			}
			l--
			continue
		}
		a.Order[i] = a.Members[a.Order[i].Name]
	}
}

func detectForRecursiveStructure(a *StructInfo, detected bool) bool {
	for _, i := range a.Order {
		if i.Name == a.Name {
			log.Printf("Struct has a member with the same name of the struct")

			//os.Exit(0)
			return true
		}
	}
	return false
}

var TmpFixMemberCall = 0

func FixMembers(a, b *StructInfo) {
	TmpFixMemberCall++
	//log.Printf("Tmp: %d", TmpFixMemberCall)
	if TmpFixMemberCall > 100000 {
		log.Panicf("Probable infinite recursion w/ %v/%v (a/b) calling himself.", a.Name, b.Name)
	}
	//if detectForRecursiveStructure(a, false) || detectForRecursiveStructure(b, false) {
	//	return
	//}
	//printer.Printf("Fixing structures (%v & %v): %+v (%p)", a.Name, b.Name, a, a)
	//if a == b {
	//	log.Printf("Pointer is the same: %v (%p) = %v (%p)", a.Name, a, b.Name, b)
	//	return
	//}
	for name, m := range a.Members {
		if _, ok := b.Members[name]; !ok {
			b.Members[name] = m
			b.Order = append(b.Order, m)
		}
		/*
			 else if !IsSameType(m, mb, 0) {
				if strings.Contains(a.Name, BasicStructName) {
					log.Printf("Here")
				}
				// Member exist in "a" and "b" with the same name, but it has not the same type which leads to ignoring one of them
				// because the fix below assume both type are identical but doesn't contain the same thing
				err := fixTypeMismatch(a.Members[name], b.Members[name])
				if errors.Is(err, ErrUnsolvableMismatch) {
					b.Members[name].Name = addUniqueNumber(b.Members[name].Name)
				} else if err != nil {
					panic(err.Error())
				}
			}
		*/
	}
	for name, m := range b.Members {
		if _, ok := a.Members[name]; !ok {
			a.Members[name] = m
			a.Order = append(a.Order, m)
		}
	}
	for i := range a.Members {
		if _, ok := b.Members[i]; !ok {
			a.PrintOrderedMembers()
			b.PrintOrderedMembers()
			log.Panicf("FixMembers: '%v' doesn't exist in b", i)
		}
		if !IsSameType(a.Members[i].T, b.Members[i].T, 0) {
			if a.Members[i] == b.Members[i] {
				//log.Panicf("Identical pointer called in 'FixMembers': %v", a.Members[i].Name)
				continue
			}
			// ????
			//if getTypeName(a.Members[i].T) == i || getTypeName(b.Members[i].T) == i {
			//	continue
			//}

			//log.Printf("%T (%v) | %T (%v) is not same type", a.Members[i].T, a.Members[i].Name, b.Members[i].T, b.Members[i].Name)
			err := fixTypeMismatch(a.Members[i], b.Members[i])
			if errors.Is(err, ErrUnsolvableMismatch) {
				b.Members[i].Name = addUniqueNumber(b.Members[i].Name)
			} else if err != nil {
				panic(err.Error())
			}
			updateOrderedMembers(a)
			updateOrderedMembers(b)
		}
	}
	//log.Printf("End")
}
