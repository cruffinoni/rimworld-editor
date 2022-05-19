package path

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"regexp"
	"strings"
)

var regexListDetection = regexp.MustCompile(`\[[.]{3}\]`)

type ListMatch struct {
	Matcher
}

func (l *ListMatch) Build(pattern string) ComputedMatcher {
	idx := strings.Index(pattern, "[")
	return &ComputedListMatch{
		tags:          make(XMLTags, 0),
		pattern:       pattern[:idx],
		patternLength: len(pattern[:idx]),
	}
}

func (l *ListMatch) RawMatch(pattern string) bool {
	return regexListDetection.MatchString(pattern)
}

type ComputedListMatch struct {
	tags          XMLTags
	tagsCount     int
	pattern       string
	patternLength int
	ComputedMatcher
}

func (c *ComputedListMatch) StrictMatch(node *xml.Element, _ string) XMLTags {
	if node.GetName()[:c.patternLength] == c.pattern {
		c.tagsCount++
		c.tags = append(c.tags, node)
	} else if c.tagsCount > 0 {
		return c.tags
	}
	return nil
}

func (c *ComputedListMatch) TrailingMatch() XMLTags {
	if c.tagsCount > 0 {
		return c.tags
	}
	return nil
}
