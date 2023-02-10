package generator

import (
	"fmt"
	"log"
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
		log.Printf("'%s' > %T", k, s.Members[k])
	}
	fmt.Printf("\n")
}

func fixTypeMismatch(a, b *member) {
	switch va := a.T.(type) {
	case *CustomType:
		if ctB, okB := b.T.(*CustomType); okB {
			fixCustomType(va, ctB)
		} else if structType, okStruct := b.T.(*StructInfo); okStruct {
			if isRelevantType(va) {
				log.Panicf("fixTypeMismatch: double relevant type => %T & %T", va, structType)
			}
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
	case reflect.Kind:
		bt := b.T.(reflect.Kind)
		if va == reflect.Int64 && bt == reflect.Float64 {
			a.T = reflect.Float64
		} else if va == reflect.Float64 && bt == reflect.Int64 {
			b.T = reflect.Float64
		}
	}
}

func isSameType(a, b any) bool {
	if a == nil || b == nil {
		return a == b
	}
	switch va := a.(type) {
	case *CustomType:
		if bType, ok := b.(*CustomType); ok {
			return va.Name == bType.Name && va.Pkg == bType.Pkg &&
				isSameType(va.Type1, bType.Type1) && isSameType(va.Type2, bType.Type2)
		} else {
			return false
		}
	case *StructInfo:
		if bType, ok := b.(*StructInfo); ok {
			return va.Name == bType.Name && hasSameMembers(va, bType)
		} else {
			return false
		}
	default:
		return reflect.TypeOf(a) == reflect.TypeOf(b)
	}
}

func fixMembers(a, b *StructInfo) {
	for name, m := range a.Members {
		if _, ok := b.Members[name]; !ok {
			b.Members[name] = m
		}
	}
	for name, m := range b.Members {
		if _, ok := a.Members[name]; !ok {
			a.Members[name] = m
		}
	}
	for i := range a.Members {
		if _, ok := b.Members[i]; !ok {
			a.printOrderedMembers()
			b.printOrderedMembers()
			log.Panicf("fixMembers: '%v' doesn't exist in b", i)
		}
		if !isSameType(a.Members[i].T, b.Members[i].T) {
			fixTypeMismatch(a.Members[i], b.Members[i])
		}
	}
}
