package path

import "github.com/cruffinoni/rimworld-editor/xml"

type StringMatch struct {
}

func (s *StringMatch) RawMatch(_ string) bool {
	return true
}

func (s *StringMatch) Build(_ string) ComputedMatcher {
	return &StringMatch{}
}

func (s *StringMatch) StrictMatch(node *xml.Element, input string) Elements {
	if input == node.GetName() {
		return Elements{node}
	}
	return nil
}

func (s *StringMatch) TrailingMatch() Elements {
	return nil
}
