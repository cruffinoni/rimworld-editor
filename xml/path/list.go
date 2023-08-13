package path

import (
	"regexp"
	"strings"

	"github.com/cruffinoni/rimworld-editor/xml"
)

var regexListDetection = regexp.MustCompile(`\[[.]{3}\]`)

type ListMatch struct {
	Matcher
}

func (l *ListMatch) Build(pattern string) ComputedMatcher {
	idx := strings.Index(pattern, "[")
	return &ComputedListMatch{
		tags:          make(Elements, 0),
		pattern:       pattern[:idx],
		patternLength: len(pattern[:idx]),
	}
}

func (l *ListMatch) RawMatch(pattern string) bool {
	return regexListDetection.MatchString(pattern)
}

type ComputedListMatch struct {
	tags          Elements
	tagsCount     int
	pattern       string
	patternLength int
	ComputedMatcher
}

func (c *ComputedListMatch) StrictMatch(node *xml.Element, _ string) Elements {
	if node.GetName()[:c.patternLength] == c.pattern {
		c.tagsCount++
		c.tags = append(c.tags, node)
	} else if c.tagsCount > 0 {
		return c.tags
	}
	return nil
}

func (c *ComputedListMatch) TrailingMatch() Elements {
	if c.tagsCount > 0 {
		return c.tags
	}
	return nil
}
