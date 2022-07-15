package attributes

type Attributes map[string]string

func (m *Attributes) Join(sep string) string {
	var (
		b []byte
		i = 0
	)
	for k, v := range *m {
		if i > 0 {
			b = append(b, sep...)
		}
		b = append(b, k...)
		b = append(b, '=')
		b = append(b, '"')
		b = append(b, v...)
		b = append(b, '"')
		i++
	}
	return string(b)
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
