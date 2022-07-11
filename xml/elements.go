package xml

type Elements []*Element

func (e Elements) FindElementFromClass(class string) *Element {
	for _, el := range e {
		if el.Attr.Get("Class") == class {
			return el
		}
	}
	return nil
}
