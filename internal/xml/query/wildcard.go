package query

import (
	"strings"

	"github.com/cruffinoni/rimworld-editor/internal/xml/domain"
)

type WildcardMatch struct {
	Matcher
}

func (w *WildcardMatch) Build(pattern string) ComputedMatcher {
	idx := strings.Index(pattern, "*")
	if idx == -1 {
		return &DefaultMatcher{}
	}
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
	nodes           Elements
	requiredPattern string
	length          int
}

func (w *ComputedWildcardMatcher) StrictMatch(node *domain.Element, _ string) Elements {
	if w.length == 0 {
		w.nodes = append(w.nodes, node)
		return nil
	} else if node.GetName()[:w.length] == w.requiredPattern {
		return Elements{node}
	}
	return nil
}

func (w *ComputedWildcardMatcher) TrailingMatch() Elements {
	if w.length == 0 {
		return w.nodes
	}
	return nil
}
