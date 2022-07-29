package generator

type buffer struct {
	writtenHeaders map[string]bool
	header         string
	body           string
	footer         string
}

const (
	headerXmlTypes = "xml/types"
	headerXml      = "xml"
)

func (b *buffer) writeImport(s string) {
	if v, ok := b.writtenHeaders[s]; ok && v {
		return
	}
	b.writtenHeaders[s] = true
	b.header += `import "github.com/cruffinoni/rimworld-editor/` + s + `"` + "\n"
}

func (b *buffer) writeToHeader(s string) {
	b.header += s
}

func (b *buffer) writeToBody(s string) {
	b.body += s
}

func (b *buffer) writeToFooter(s string) {
	b.footer += s
}

func (b *buffer) bytes() []byte {
	return []byte(b.header + b.body + b.footer)
}
