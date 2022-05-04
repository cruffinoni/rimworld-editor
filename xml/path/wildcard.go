package path

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"strings"
)

type WildcardMatcher struct{}

type ComputedWildcardMatcher struct {
	requiredPattern string
	length          int
}

func (w *WildcardMatcher) RawMatch(pattern string) bool {
	return strings.Contains(pattern, "*")
}

// Compute the required pattern and length of the wildcard
func computeMatcher[E ComputerMatcher, B Matcher](base B) E {
	return E{}
}

func (w *ComputedWildcardMatcher) StrictMatch(node *xml.Tag, input string) bool {

}
