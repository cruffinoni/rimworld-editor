package path

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"log"
	"strings"
)

type pattern struct {
	path    string
	matcher Matcher
}

type Path struct {
	patterns []pattern
	discover *xml.Discover
}

type Matcher interface {
	RawMatch(pattern string) bool
}

// Template ComputerMatcher with any type and add a function to build the matcher
type ComputerMatcher interface {
	StrictMatch(node *xml.Tag, input string) bool
}

var (
	DefaultMatcher = &StringMatcher{}
)

// savegame>meta>gameVersion
// savegame>meta>game*

func NewPathing(rawPattern string) *Path {
	p := &Path{
		patterns: make([]pattern, 0),
	}
	for _, s := range strings.Split(rawPattern, ">") {
		p.patterns = append(p.patterns, pattern{s, &strictMatcher{}})
	}
}

/*
TODO:
- Implements wildcard for path
- Impl array for path
*/

func (p *Path) Find(root *xml.Discover) *xml.Tag {
	if len(p.pattern) == 0 {
		log.Printf("null pattern")
		return nil
	}
	cpyPattern := make([]string, len(p.pattern))
	copy(cpyPattern, p.pattern)
	n := root.Tag
	for n != nil {
		log.Printf("Comparing %s to %s", cpyPattern[0], n.GetName())
		if n.GetName() == cpyPattern[0] {
			cpyPattern = cpyPattern[1:]
			if len(cpyPattern) == 0 {
				return n
			}
			n = n.Child
		} else {
			n = n.Next
		}
	}
	return n
}
