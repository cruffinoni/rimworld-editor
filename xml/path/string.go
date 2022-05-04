package path

import "github.com/cruffinoni/rimworld-editor/xml"

type StringMatcher struct {
}

func (s *StringMatcher) RawMatch(pattern string) bool {
	return true
}

func (s *StringMatcher) StrictMatch(node *xml.Tag, input string) bool {
	return input == node.GetName()
}
