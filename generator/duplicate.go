package generator

func toUpperIfNeeded(r []byte) string {
	if r[0] >= 'a' && r[0] <= 'z' {
		r[0] -= 32
		return string(r)
	}
	return string(r)
}

func deleteDuplicateTitle(a *StructInfo) {
	edited := false
	for i := range a.Members {
		toTitle := toUpperIfNeeded([]byte(i))
		if _, okBasic := a.Members[toTitle]; okBasic && toTitle != i {
			delete(a.Members, toTitle)
			edited = true
		}
	}
	if edited {
		updateOrderedMembers(a)
	}
}
