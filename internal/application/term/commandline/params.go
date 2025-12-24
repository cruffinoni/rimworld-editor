package commandline

import "reflect"

type Params struct {
	Name          string
	alias         []string
	Description   string
	Type          reflect.Kind
	OptionalValue any
}

func (p *Params) isNamed(name string) bool {
	for _, n := range p.alias {
		if n == name {
			return true
		}
	}
	return p.Name == name
}
