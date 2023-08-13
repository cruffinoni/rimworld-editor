package path

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/cruffinoni/rimworld-editor/xml"
)

var regexArrayDetection = regexp.MustCompile(`\[\d+\]`)

type ArrayMatch struct {
	Matcher
}

func (a *ArrayMatch) Build(pattern string) ComputedMatcher {
	res := regexArrayDetection.Find([]byte(pattern))
	if res == nil {
		return nil
	}
	res = res[1 : len(res)-1]
	if nb, err := strconv.Atoi(string(res)); err != nil {
		return nil
	} else {
		idx := strings.Index(pattern, "[")

		return &ComputedArrayMatch{
			listIndex:     nb,
			pattern:       pattern[:idx],
			patternLength: len(pattern[:idx]),
			matchedCount:  1, // Array starts at 1
		}
	}
}

func (a *ArrayMatch) RawMatch(pattern string) bool {
	return regexArrayDetection.MatchString(pattern)
}

type ComputedArrayMatch struct {
	listIndex     int
	pattern       string
	patternLength int
	matchedCount  int
	ComputedMatcher
}

func (c *ComputedArrayMatch) StrictMatch(node *xml.Element, _ string) Elements {
	if node.GetName()[:c.patternLength] == c.pattern {
		//log.Printf("Total of matched count: %d", c.matchedCount)
		if c.matchedCount == c.listIndex {
			return Elements{node}
		}
		c.matchedCount++
	}
	return nil
}

func (c *ComputedArrayMatch) TrailingMatch() Elements {
	return nil
}
