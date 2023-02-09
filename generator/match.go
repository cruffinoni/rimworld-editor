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
		if ct.name == "Empty" && ct.pkg == "*primary" {
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
	if a.type1 != b.type1 {
		if isRelevantType(a.type1) && !isRelevantType(b.type1) {
			b.type1 = a.type1
			// If the main type is changed, update the package name, type name, and import path
			b.name = a.name
			b.pkg = a.pkg
			b.importPath = a.importPath
		} else if !isRelevantType(a.type1) && isRelevantType(b.type1) {
			a.type1 = b.type1
			// Perform the same update here
			a.name = b.name
			a.pkg = b.pkg
			a.importPath = b.importPath
		}
	}
	// Keep the relevant type between both `type2` fields.
	if a.type2 != nil && a.type2 != b.type2 {
		if isRelevantType(a.type2) && !isRelevantType(b.type2) {
			b.type2 = a.type2
		} else if !isRelevantType(a.type2) && isRelevantType(b.type2) {
			a.type2 = b.type2
		}
	}
}

// printOrderedMembers print members of s in an alphabetic order
func (s *StructInfo) printOrderedMembers() {
	if len(s.members) == 0 {
		return
	}
	m := make([]string, 0, len(s.members))
	for k := range s.members {
		m = append(m, k)
	}
	sort.Strings(m)
	log.Printf("Struct %v", s.name)
	for _, k := range m {
		log.Printf("'%s' > %T", k, s.members[k])
	}
	fmt.Printf("\n")
}

func fixTypeMismatch(a, b *member) {
	switch va := a.t.(type) {
	case *CustomType:
		if ctB, okB := b.t.(*CustomType); okB {
			fixCustomType(va, ctB)
		} else if structType, okStruct := b.t.(*StructInfo); okStruct {
			if isRelevantType(va) {
				log.Panicf("fixTypeMismatch: double relevant type => %T & %T", va, structType)
			}
			a.t = b.t
		} else {
			b.t = a.t
		}
	case *StructInfo:
		if bStruct, okStruct := b.t.(*StructInfo); okStruct {
			fixMembers(va, bStruct)
		} else {
			b.t = a.t
		}
	case reflect.Kind:
		bt := b.t.(reflect.Kind)
		if va == reflect.Int64 && bt == reflect.Float64 {
			a.t = reflect.Float64
		} else if va == reflect.Float64 && bt == reflect.Int64 {
			b.t = reflect.Float64
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
			return va.name == bType.name && va.pkg == bType.pkg &&
				isSameType(va.type1, bType.type1) && isSameType(va.type2, bType.type2)
		} else {
			return false
		}
	case *StructInfo:
		if bType, ok := b.(*StructInfo); ok {
			return va.name == bType.name && hasSameMembers(va, bType)
		} else {
			return false
		}
	default:
		return reflect.TypeOf(a) == reflect.TypeOf(b)
	}
}

func fixMembers(a, b *StructInfo) {
	for name, m := range a.members {
		if _, ok := b.members[name]; !ok {
			b.members[name] = m
		}
	}
	for name, m := range b.members {
		if _, ok := a.members[name]; !ok {
			a.members[name] = m
		}
	}
	for i := range a.members {
		if _, ok := b.members[i]; !ok {
			a.printOrderedMembers()
			b.printOrderedMembers()
			log.Panicf("fixMembers: '%v' doesn't exist in b", i)
		}
		if !isSameType(a.members[i].t, b.members[i].t) {
			fixTypeMismatch(a.members[i], b.members[i])
		}
	}
}
