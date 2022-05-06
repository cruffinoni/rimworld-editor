package path

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"log"
	"strings"
)

type pattern struct {
	path    string
	matcher ComputedMatcher
}

type XMLTags []*xml.Tag

// ResultType is the type of the result of a match.
//type ResultType interface {
//	*xml.Tag | []*xml.Tag
//}

// Path is a path to a node in the XML tree.
type Path struct {
	patterns []*pattern
	discover *xml.Discover
}

type Matcher interface {
	Build(pattern string) ComputedMatcher
	RawMatch(pattern string) bool
}

type ComputedMatcher interface {
	StrictMatch(node *xml.Tag, input string) XMLTags
	TrailingMatch() XMLTags
}

type DefaultMatcher = StringMatch

var matchers = []Matcher{
	&WildcardMatch{},
	&ArrayMatch{},
	&ListMatch{},
}

func NewPathing(rawPattern string) *Path {
	split := strings.Split(rawPattern, ">")
	p := &Path{
		patterns: make([]*pattern, 0, len(split)),
	}
	for _, s := range split {
		pm := &pattern{
			path:    s,
			matcher: &DefaultMatcher{},
		}
		p.patterns = append(p.patterns, pm)
		for _, m := range matchers {
			if m.RawMatch(s) {
				pm.matcher = m.Build(s)
				log.Printf("Found matcher: %T", m)
				break
			}
		}
	}
	return p
}

func FindWithPath(rawPattern string, root *xml.Discover) XMLTags {
	p := NewPathing(rawPattern)
	return p.Find(root)
}

func (p *Path) Find(root *xml.Discover) XMLTags {
	var (
		r          XMLTags
		n          = root.Tag
		patternIdx = 0
	)
	cpyPatterns := make([]*pattern, len(p.patterns))
	copy(cpyPatterns, p.patterns)
	for n != nil {
		if r = p.patterns[patternIdx].matcher.StrictMatch(n, cpyPatterns[0].path); r == nil {
			n = n.Next
			continue
		}
		patternIdx++
		cpyPatterns = cpyPatterns[1:]
		if len(cpyPatterns) == 0 {
			return r
		} else {
			n = n.Child
		}
	}
	if r = p.patterns[patternIdx].matcher.TrailingMatch(); r != nil {
		return r
	}
	return nil
}
