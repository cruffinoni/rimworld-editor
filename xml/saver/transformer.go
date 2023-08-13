package saver

type Transformer interface {
	TransformToXML(buffer *Buffer) error
	// GetXMLTag give the XML tag. It returns nil when the type can't define a unique
	// XML tag.
	GetXMLTag() []byte
}
