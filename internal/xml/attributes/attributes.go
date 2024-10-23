package attributes

import (
	"strings"
)

// Attributes is a map of XML attributes.
type Attributes map[string]string

func (m *Attributes) Join(sep string) string {
	if m.Empty() {
		return ""
	}
	var (
		sb strings.Builder
		i  = 0
	)
	for k, v := range *m {
		if i > 0 {
			sb.WriteString(sep)
		}
		sb.WriteString(k + "=" + "\"" + v + "\"")
		i++
	}
	return sb.String()
}

func (m *Attributes) Empty() bool {
	return len(*m) == 0
}

func (m *Attributes) Get(key string) string {
	if v, ok := (*m)[key]; ok {
		return v
	}
	return ""
}
