package path

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"log"
	"strings"
)

type WildcardMatch struct {
	Matcher
}

func (w *WildcardMatch) Build(pattern string) ComputedMatcher {
	idx := strings.Index(pattern, "*")
	if idx == -1 {
		return &DefaultMatcher{}
	}
	log.Printf("Required pattern found: %s", pattern[:idx])
	return &ComputedWildcardMatcher{
		requiredPattern: pattern[:idx],
		length:          len(pattern[:idx]),
	}
}

func (w *WildcardMatch) RawMatch(pattern string) bool {
	return strings.Contains(pattern, "*")
}

type ComputedWildcardMatcher struct {
	ComputedMatcher
	requiredPattern string
	length          int
}

func (w *ComputedWildcardMatcher) StrictMatch(node *xml.Tag, input string) XMLTags {
	if node.GetName()[:w.length] == w.requiredPattern {
		return XMLTags{node}
	}
	return nil
}

func (w *ComputedWildcardMatcher) TrailingMatch() XMLTags {
	return nil
}
