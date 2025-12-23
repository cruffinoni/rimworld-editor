package encoder

type Transformer interface {
	// TransformToXML transforms the type to XML. It returns ErrEmptyValue when the
	// type is empty and can't be transformed to XML. The tag must be open BEFORE calling this method.
	TransformToXML(buffer *Buffer) error
	// GetXMLTag give the XML tag. It returns nil when the type can't define a unique
	// XML tag.
	GetXMLTag() []byte
}
