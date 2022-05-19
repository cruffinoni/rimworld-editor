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

type XMLTags []*xml.Element

// ResultType is the type of the result of a match.
//type ResultType interface {
//	*xml.Element | []*xml.Element
//}

// Path is a path to a node in the XML tree.
type Path struct {
	patterns []*pattern
	discover *xml.Tree
}

type Matcher interface {
	Build(pattern string) ComputedMatcher
	RawMatch(pattern string) bool
}

type ComputedMatcher interface {
	StrictMatch(node *xml.Element, input string) XMLTags
	TrailingMatch() XMLTags
}

type DefaultMatcher = StringMatch

var matchers = []Matcher{
	&WildcardMatch{},
	&ArrayMatch{},
	&ListMatch{},
	&AttributeMatch{},
}

func NewPathing(rawPattern string) *Path {
	split := strings.Split(rawPattern, ">")
	p := &Path{
		patterns: make([]*pattern, 0, len(split)),
	}
	for _, s := range split {
		//log.Printf("s: %s", s)
		pm := &pattern{
			path:    s,
			matcher: &DefaultMatcher{},
		}
		for _, m := range matchers {
			if m.RawMatch(s) {
				pm.matcher = m.Build(s)
				if pm.matcher == nil {
					log.Fatalf("failed to build matcher for %s", s)
				}
				//log.Printf("Found matcher: %T", m)
				break
			}
		}
		//log.Printf("%v: Matcher: %T", pm.path, pm.matcher)
		p.patterns = append(p.patterns, pm)
	}
	return p
}

func FindWithPath(rawPattern string, root *xml.Element) XMLTags {
	p := NewPathing(rawPattern)
	return p.Find(root)
}

func (p *Path) Find(root *xml.Element) XMLTags {
	var (
		r          XMLTags
		n          = root
		patternIdx = 0
	)
	cpyPatterns := make([]*pattern, len(p.patterns))
	copy(cpyPatterns, p.patterns)
	for n != nil {
		if r = p.patterns[patternIdx].matcher.StrictMatch(n, cpyPatterns[0].path); r == nil {
			n = n.Right
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
	log.Printf("Not found at %s (%T)", cpyPatterns[0].path, cpyPatterns[0].matcher)
	if r = p.patterns[patternIdx].matcher.TrailingMatch(); r != nil {
		return r
	}
	return nil
}
