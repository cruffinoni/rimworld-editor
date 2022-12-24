package path

import (
	"regexp"
	"strings"

	"github.com/cruffinoni/rimworld-editor/xml"
)

var regexAttributeDetection = regexp.MustCompile(`({[a-zA-Z]+})|({[a-zA-Z]+="[a-zA-Z]+"})`)

type AttributeMatch struct {
	Matcher
}

func (a *AttributeMatch) Build(pattern string) ComputedMatcher {
	res := regexAttributeDetection.Find([]byte(pattern))
	if res == nil {
		return nil
	}
	resStr := string(res[1 : len(res)-1])
	if idx := strings.Index(resStr, "="); idx != -1 {
		return &ComputedAttributeMatch{
			key:   resStr[:idx],
			value: resStr[idx+2 : len(resStr)-1],
			query: strings.ReplaceAll(pattern, string(res), ""),
		}
	} else {
		return &ComputedAttributeMatch{
			value: resStr,
			query: strings.ReplaceAll(pattern, string(res), ""),
		}
	}
}

func (a *AttributeMatch) RawMatch(pattern string) bool {
	return regexAttributeDetection.MatchString(pattern)
}

type ComputedAttributeMatch struct {
	ComputedMatcher
	key   string
	value string
	query string
}

func (c *ComputedAttributeMatch) StrictMatch(node *xml.Element, _ string) Elements {
	if node.GetName() != c.query {
		return nil
	}
	if c.key != "" {
		if node.Attr[c.key] == c.value {
			return Elements{node}
		}
	} else {
		for _, attr := range node.Attr {
			if attr == c.value {
				return Elements{node}
			}
		}
	}
	return nil
}

func (c *ComputedAttributeMatch) TrailingMatch() Elements {
	return nil
}
